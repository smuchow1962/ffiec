package verify

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
)

// §7 verification procedure over the audit-file shape.
//
// This is the §7 step 1-10 walk the spec normates (chain-of-custody §7),
// operating on the on-disk `audit_file` object the negative-vector corpus
// materializes: a header, an ordered entry list, and a daily seal. It is
// DISTINCT from the older Ledger/ChainEntry walk in pipeline.go — that
// shape binds the per-event MAC over event_payload_jcs with an
// entry_id-derived key; the audit-file shape binds it per §4.1 over
// (prev_hash || canonical) under the tenant-derived session key. The two
// chain shapes are not interchangeable, so each carries its own walk.
//
// The walk is verifier capability, not a test-only construct: a future
// `verifier verify --audit-file` CLI mode drives WalkAuditFile and prints
// the Outcome's normative triple. The conformance gate consumes the same
// function so the corpus asserts live verifier behavior, not a re-derived
// answer key.
//
// Scope: steps 1-11. Steps 1-10 are the crypto-walkable core (no key
// needed). Step 11 (the §4.3 Ed25519 seal signature) runs when the caller
// supplies the institution's public key via WalkAuditFileWithKey; the
// no-key WalkAuditFile path stops after step 10. Step 12 (cadence/dev-
// mode) is bound INTO the step-11 sign_payload — a tampered cadence or
// dev_mode changes the reconstructed bytes and fails the signature — so
// it needs no separate walk function. The structural slice of step 11
// that needs no key — the §4.3.2 algorithm/key-type field compare — stays
// available via CheckAlgorithmKeyType for callers that have only the
// resolved key-type descriptor.

// AuditFile is the on-disk shape the corpus materializes under the
// top-level `audit_file` key. Field tags mirror the corpus JSON.
type AuditFile struct {
	Header  AuditHeader  `json:"header"`
	Entries []AuditEntry `json:"entries"`
	Seal    AuditSeal    `json:"seal"`

	// NDJSONTruncated, when non-nil, carries the mid-write-truncation
	// fixture marker (§7 implementation note). The corpus represents a
	// truncated file structurally rather than as raw bytes, so the
	// pre-flight check reads this marker instead of seeking the last byte.
	NDJSONTruncated *string `json:"_ndjson_truncated"`
}

// AuditHeader is the §7 file-header block (steps 1-3a inputs).
type AuditHeader struct {
	FormatVersion       string `json:"format_version"`
	TenantID            string `json:"tenant_id"`
	SealDate            string `json:"seal_date"`
	GenesisHashHex      string `json:"genesis_hash_hex"`
	HKDFInputsDigestHex string `json:"hkdf_inputs_digest_hex"`
}

// AuditEntry is one chain entry (steps 4-9 inputs). The event object is
// kept as a decoded map only for the binding check; the MAC is computed
// over event_canonical_hex (the pinned canonical bytes), never over a
// re-canonicalization, so the walk does not depend on re-deriving JCS.
type AuditEntry struct {
	Seq               int                    `json:"seq"`
	RunID             string                 `json:"run_id"`
	KeyVersion        int                    `json:"key_version"`
	KeyFingerprintHex string                 `json:"key_fingerprint_hex"`
	FormatVersion     string                 `json:"format_version"`
	Event             map[string]interface{} `json:"event"`
	EventCanonicalHex string                 `json:"event_canonical_hex"`
	PrevHashHex       string                 `json:"prev_hash_hex"`
	PayloadHashHex    string                 `json:"payload_hash_hex"`
}

// AuditSeal is the daily seal block (step 10 input + the §4.3 step-11
// signature fields). The sign_payload_* fields carry the §4.3 byte form
// the seal was signed over; the walk reconstructs that form from the
// structured fields and confirms it equals the published bytes before
// verifying the signature, so a published payload that does not bind the
// real merkle_root/tenant is itself a step-11 failure.
type AuditSeal struct {
	MerkleRootHex string `json:"merkle_root_hex"`
	Algorithm     string `json:"algorithm"`
	// ResolvedPublicKeyType is the descriptor a §4.3.2 algorithm/key-type
	// mismatch fixture carries (e.g. "rsa-3072" when algorithm is
	// "ed25519"). Absent on conformant seals.
	ResolvedPublicKeyType string `json:"resolved_public_key_type"`

	// §4.3 sign_payload + signature fields (step 11 inputs). The seal
	// records both the reconstructable structured fields AND the published
	// sign_payload_hex so the verifier can cross-check the two.
	SealDate            string `json:"seal_date"`
	FormatVersion       string `json:"format_version"`
	SignPayloadVersion  string `json:"sign_payload_version"`
	Cadence             string `json:"cadence"`
	DevMode             bool   `json:"dev_mode"`
	TenantID            string `json:"tenant_id"`
	HKDFInputsDigestHex string `json:"hkdf_inputs_digest_hex"`
	SignPayloadHex      string `json:"sign_payload_hex"`
	SignatureB64        string `json:"signature_b64"`

	// v1.0b / v1.0c sign_payload additions, threaded through to the
	// reconstruction when the seal's version binds them.
	KeyVersionsCanon         string `json:"key_versions_canon"`
	KMSHandleURIsDigestHex   string `json:"kms_handle_uris_digest_hex"`
	OperationalEventsLogRoot string `json:"operational_events_log_root_hex"`
}

// Outcome is the §7 normative verifier output triple plus the §10.12
// exit code. Status is "PASS" or "FAIL"; Step is the bare step number
// ("1".."10") or "pre-flight"; Reason is the rendered reason string with
// position tokens substituted; ExitCode is the §10.12 code.
type Outcome struct {
	Status   string
	Step     string
	Reason   string
	ExitCode int
}

// pass is the clean-walk outcome (status PASS, exit 0, no step/reason).
func pass() Outcome { return Outcome{Status: "PASS", ExitCode: 0} }

// fail builds a FAIL outcome at a named step with exit code 1 (a §7
// integrity failure per §10.12).
func fail(step, reason string) Outcome {
	return Outcome{Status: "FAIL", Step: step, Reason: reason, ExitCode: 1}
}

// structuralFail builds a FAIL outcome for a pre-flight / structural
// input error (exit code 2 per §10.12).
func structuralFail(step, reason string) Outcome {
	return Outcome{Status: "FAIL", Step: step, Reason: reason, ExitCode: 2}
}

// IKMRegistry resolves the institution IKM for a key_version. The corpus
// pins two generations (v1, v2) in chain_vectors.json; an absent version
// (a tampered or unknown key_version) returns ok=false, driving the §7
// step-7 lookup-miss reason.
type IKMRegistry map[int][]byte

// WalkAuditFile executes the §7 step 1-10 procedure over af, returning
// the first failing step's normative Outcome, or a clean PASS when every
// executed step holds. Step 11 (seal signature) is SKIPPED — this is the
// no-key path, kept for callers that drive only the crypto-walkable core
// (the contract-only negative fixtures whose seals carry no real
// signature, and any caller without a configured public key).
//
// For the full §7 walk including the live Ed25519 step-11 verification,
// call WalkAuditFileWithKey with the institution's public key.
func WalkAuditFile(af *AuditFile, ikms IKMRegistry) Outcome {
	return WalkAuditFileWithKey(af, ikms, nil)
}

// WalkAuditFileWithKey executes the §7 procedure over af. When pub is
// non-nil it runs the full steps 1-11; when pub is nil it runs steps 1-10
// and stops (the WalkAuditFile path). The walk stops at the first failure
// per §7 ("reports the most specific reason and stops processing the
// affected unit").
//
// The verifier consumes pub as PUBLIC-key material only. It never reads
// the private seed; signing is the corpus generator's job, not the
// verifier's.
//
// Cognitive-complexity note: the body is a flat sequence of guarded step
// calls, each returning an Outcome whose non-empty Status==FAIL short-
// circuits. A reviewer reads the step list top-to-bottom against the §7
// numbered procedure. The per-step logic lives in the small functions
// below (and in auditsign.go for step 11) so this orchestrator stays a
// legible step ledger.
func WalkAuditFileWithKey(af *AuditFile, ikms IKMRegistry, pub ed25519.PublicKey) Outcome {
	if out := checkTruncation(af); out.Status == "FAIL" {
		return out
	}
	if out := checkFormatVersion(af.Header); out.Status == "FAIL" {
		return out
	}
	if out := checkHKDFInputsDigest(af.Header); out.Status == "FAIL" {
		return out
	}
	if out := checkGenesis(af.Header); out.Status == "FAIL" {
		return out
	}
	if out := walkEntries(af, ikms); out.Status == "FAIL" {
		return out
	}
	if out := checkMerkleRootAudit(af); out.Status == "FAIL" {
		return out
	}
	// §7 step 11: seal signature, only when a public key is configured.
	if pub != nil {
		if out := CheckSealSignatureV1(af.Seal, pub); out.Status == "FAIL" {
			return out
		}
	}
	return pass()
}

// checkTruncation is the §7 pre-flight mid-write-truncation refusal. The
// corpus marks a truncated file with the _ndjson_truncated marker rather
// than raw bytes, so the verifier rejects on the marker's presence. Exit
// code 2 (structural input error).
func checkTruncation(af *AuditFile) Outcome {
	if af.NDJSONTruncated != nil {
		return structuralFail("pre-flight", "audit file ends mid-line — possible mid-write crash")
	}
	return pass()
}

// checkFormatVersion is §7 step 1. A v1 verifier accepts only "v1";
// every variant ("v2", "v1.1", "V1") is refused here with the value
// rendered into the reason.
func checkFormatVersion(h AuditHeader) Outcome {
	if h.FormatVersion != constants.FormatVersion {
		return fail("1", fmt.Sprintf(
			"format_version %s not supported by this verifier (running v1)", h.FormatVersion))
	}
	return pass()
}

// checkHKDFInputsDigest is §7 step 2. Recompute SHA-256(salt ||
// info_for_tenant || length_LE32) and compare against the header pin.
func checkHKDFInputsDigest(h AuditHeader) Outcome {
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + h.TenantID)
	lengthLE32 := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthLE32, uint32(constants.HKDFOutputLength))

	sum := sha256.New()
	sum.Write([]byte(constants.HKDFSalt))
	sum.Write(info)
	sum.Write(lengthLE32)
	got := hex.EncodeToString(sum.Sum(nil))

	if got != h.HKDFInputsDigestHex {
		return fail("2", "header HKDF inputs do not match running v1 inputs")
	}
	return pass()
}

// genesisZeroHex is the §4.1 genesis prev_hash: 32 zero bytes, hex.
var genesisZeroHex = hex.EncodeToString(make([]byte, constants.GenesisPrevHashLen))

// checkGenesis is §7 step 3. The header genesis_hash MUST be 32 zero bytes.
func checkGenesis(h AuditHeader) Outcome {
	if h.GenesisHashHex != genesisZeroHex {
		return fail("3", "header genesis_hash does not match v1 constant")
	}
	return pass()
}

// walkEntries runs steps 4-9 in (run_id, seq) file order, threading the
// structurally-walked expected_prev_hash forward. The MAC at step 9 uses
// expected_prev (the walked value), NOT entry.prev_hash, per the §7 step
// 9 footgun note. The first failing entry's Outcome short-circuits the
// walk.
func walkEntries(af *AuditFile, ikms IKMRegistry) Outcome {
	expectedPrev := genesisZeroHex
	for i := range af.Entries {
		e := &af.Entries[i]
		expectedSeq := i + 1
		if out := walkEntry(af.Header, e, expectedSeq, expectedPrev, ikms); out.Status == "FAIL" {
			return out
		}
		expectedPrev = e.PayloadHashHex
	}
	return pass()
}

// walkEntry runs steps 4-9 for one entry. The step order is normative
// (§7 "Step ordering"): 7 (IKM lookup) before 8 (fingerprint) before 9
// (MAC), because step 9's expected MAC depends on the IKM step 7 resolved
// and step 8 confirmed.
func walkEntry(h AuditHeader, e *AuditEntry, expectedSeq int, expectedPrev string, ikms IKMRegistry) Outcome {
	if out := checkBinding(h, e, expectedSeq); out.Status == "FAIL" {
		return out
	}
	if out := checkEntryFormat(h, e, expectedSeq); out.Status == "FAIL" {
		return out
	}
	if out := checkStructuralLink(e, expectedSeq, expectedPrev); out.Status == "FAIL" {
		return out
	}
	ikm, ok := ikms[e.KeyVersion]
	if !ok {
		return fail("7", fmt.Sprintf(
			"unknown key_version: no IKM for (tenant=%s, key_version=%d) at seq %d",
			h.TenantID, e.KeyVersion, expectedSeq))
	}
	if out := checkFingerprintEntry(h, e, ikm, expectedSeq); out.Status == "FAIL" {
		return out
	}
	return checkMAC(h, e, ikm, expectedPrev, expectedSeq)
}

// checkBinding is §7 step 4: the per-entry cross-chain-lift defence. The
// event's tenant_id MUST match the header tenant_id (the corpus's
// run_id binding is checked structurally via seq at step 6).
func checkBinding(h AuditHeader, e *AuditEntry, expectedSeq int) Outcome {
	if eventTenantID(e) != h.TenantID {
		return fail("4", fmt.Sprintf(
			"cross-chain lift detected at seq %d (event.tenant_id mismatch)", expectedSeq))
	}
	return pass()
}

// eventTenantID reads the event object's tenant_id, returning "" when
// absent or non-string so a missing binding surfaces as a step-4 mismatch
// rather than a panic.
func eventTenantID(e *AuditEntry) string {
	v, ok := e.Event["tenant_id"].(string)
	if !ok {
		return ""
	}
	return v
}

// checkEntryFormat is §7 step 5: entry.format_version == header.format_version.
func checkEntryFormat(h AuditHeader, e *AuditEntry, expectedSeq int) Outcome {
	if e.FormatVersion != h.FormatVersion {
		return fail("5", fmt.Sprintf("format_version mismatch at seq %d", expectedSeq))
	}
	return pass()
}

// checkStructuralLink is §7 step 6. prev_hash is checked first (the
// chain-link reason); seq ordering is the secondary check. Both surface
// at the same step with the chain-link reason because the corpus's
// reorder/substitution tampers break the prev_hash linkage.
func checkStructuralLink(e *AuditEntry, expectedSeq int, expectedPrev string) Outcome {
	if e.PrevHashHex != expectedPrev {
		return fail("6", fmt.Sprintf("chain link broken at seq %d", expectedSeq))
	}
	if e.Seq != expectedSeq {
		return fail("6", fmt.Sprintf("chain link broken at seq %d", expectedSeq))
	}
	return pass()
}

// checkFingerprintEntry is §7 step 8: recompute SHA-256(tenant || ikm)[:16]
// and compare against the entry's recorded fingerprint. No MAC compute
// happens on a mismatch — this catches the botched-rotation failure mode
// before the expensive MAC.
func checkFingerprintEntry(h AuditHeader, e *AuditEntry, ikm []byte, expectedSeq int) Outcome {
	got := KeyFingerprintHex(h.TenantID, ikm)
	if got != e.KeyFingerprintHex {
		return fail("8", fmt.Sprintf(
			"key_fingerprint mismatch at seq %d: looked-up IKM does not match the entry's recorded fingerprint",
			expectedSeq))
	}
	return pass()
}

// checkMAC is §7 step 9. Derive the session key, recompute
// HMAC-SHA-256(session_key, expected_prev || canonical_bytes), and
// compare against the entry's payload_hash. The MAC input uses
// expectedPrev (the structurally-walked value), never entry.prev_hash.
func checkMAC(h AuditHeader, e *AuditEntry, ikm []byte, expectedPrev string, expectedSeq int) Outcome {
	prev, err := hex.DecodeString(expectedPrev)
	if err != nil {
		return fail("9", fmt.Sprintf("payload_hash MAC mismatch at seq %d", expectedSeq))
	}
	canonical, err := hex.DecodeString(e.EventCanonicalHex)
	if err != nil {
		return fail("9", fmt.Sprintf("payload_hash MAC mismatch at seq %d", expectedSeq))
	}
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + h.TenantID)
	sessionKey := hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength)
	gotMAC := hex.EncodeToString(hmacSHA256(sessionKey, append(prev, canonical...)))
	if gotMAC != e.PayloadHashHex {
		return fail("9", fmt.Sprintf("payload_hash MAC mismatch at seq %d", expectedSeq))
	}
	return pass()
}

// checkMerkleRootAudit is §7 step 10: recompute the RFC 6962 root over the
// entries' payload_hash leaves and compare against the seal's merkle_root.
func checkMerkleRootAudit(af *AuditFile) Outcome {
	leaves := make([][]byte, 0, len(af.Entries))
	for i := range af.Entries {
		ph, err := hex.DecodeString(af.Entries[i].PayloadHashHex)
		if err != nil {
			return fail("10", "merkle root mismatch — ledger contents do not produce sealed root")
		}
		leaves = append(leaves, MerkleLeafHash(ph))
	}
	got := hex.EncodeToString(MerkleTreeHash(leaves))
	if got != af.Seal.MerkleRootHex {
		return fail("10", "merkle root mismatch — ledger contents do not produce sealed root")
	}
	return pass()
}

// CheckAlgorithmKeyType is the structural slice of §7 step 11 that needs
// no public key: when the seal's algorithm string disagrees with the
// resolved public-key type (e.g. seal claims "ed25519" but the resolved
// key is "rsa-3072"), the verifier reports the specific mismatch rather
// than the generic signature failure. Callers with a resolved-key-type
// descriptor (the §4.3.2 mismatch corpus) drive this; the base walk does
// not, because the conformant corpus carries no resolved-key-type field.
func CheckAlgorithmKeyType(seal AuditSeal) (Outcome, bool) {
	if seal.ResolvedPublicKeyType == "" || seal.Algorithm == seal.ResolvedPublicKeyType {
		return pass(), false
	}
	return fail("11", "algorithm/key-type mismatch at signature verification"), true
}
