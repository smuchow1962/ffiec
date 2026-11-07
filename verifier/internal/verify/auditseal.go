package verify

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/mmpworks/ffiec/verifier/internal/exitcodes"
)

// §7 steps 10 and 11 over an on-disk NDJSON audit file. The chain sink writes the
// file during capture, and the seal job signs the day's root later, so the seal
// arrives as its own input. The public key is a third input, and no field of the
// audit file or the seal can supply it.

// NoSealAnomaly is the anomaly line for a clean walk that had no seal to check.
const NoSealAnomaly = "no daily seal supplied: steps 10 and 11 did not run, so a cut at the end of the file cannot be detected"

// sealAlgorithmEd25519 is the only signature this verifier checks. --root-key loads
// an Ed25519 key, so a seal naming another algorithm cannot verify under it.
const sealAlgorithmEd25519 = "ed25519"

const (
	signPayloadPreAmendment = ""
	signPayloadV1_0b        = "v1.0b"
	signPayloadV1_0c        = "v1.0c"
)

// LoadAuditSeal reads one seal record. Its field names are those of the audit_file.seal
// block in the conformance vectors. A seal the verifier cannot read, or one that names
// no root, means the §7 procedure cannot begin, so it comes back as a FAIL Outcome with
// exit code 2.
func LoadAuditSeal(raw []byte) (AuditSeal, Outcome) {
	var seal AuditSeal
	if err := json.Unmarshal(raw, &seal); err != nil {
		return AuditSeal{}, structuralFail("pre-flight", fmt.Sprintf("malformed seal: %v", err))
	}
	if seal.MerkleRootHex == "" {
		return AuditSeal{}, structuralFail("pre-flight", "seal names no merkle_root_hex")
	}
	return seal, pass()
}

// WalkAuditChainSealed runs the keyed walk and then the seal steps: the pre-flight
// and steps 1-11. af.Seal holds the seal.
func WalkAuditChainSealed(af *AuditFile, ikms IKMRegistry, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	if out := WalkAuditChain(af, ikms); out.Status == "FAIL" {
		return out
	}
	if out := checkSeal(af, pub, policy); out.Status == "FAIL" {
		return out
	}
	return judgeSeal(af, policy)
}

// WalkAuditChainWitnessSealed runs the no-key walk and then the seal steps: the
// pre-flight and steps 1-6, 10 and 11. Steps 10 and 11 need no master key. Steps 7-9
// do, so a pass here stays a witness pass: a posture warning adds its anomaly lines and
// leaves the witness status, because PASS-WITH-ANOMALY would imply that the MACs were checked.
func WalkAuditChainWitnessSealed(af *AuditFile, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	witness := WalkAuditChainWitness(af)
	if witness.Status == "FAIL" {
		return witness
	}
	if out := checkSeal(af, pub, policy); out.Status == "FAIL" {
		return out
	}
	judged := judgeSeal(af, policy)
	if judged.Status == "FAIL" {
		return judged
	}
	witness.Anomalies = judged.Anomalies
	return witness
}

// checkSeal stops at the first failure, as §7 does, and SealAnomalies names what it
// did not check. A list of checks in §7 order keeps this auditor-readable.
func checkSeal(af *AuditFile, pub ed25519.PublicKey, policy SealPolicy) Outcome {
	for _, check := range []func() Outcome{
		func() Outcome { return checkMerkleRootAudit(af) },
		func() Outcome { return checkSealCoversFile(af) },
		func() Outcome { return checkSealFloor(af.Seal, policy) },
		func() Outcome { return checkSealAlgorithm(af.Seal) },
		func() Outcome { return checkSealSingleSignature(af.Seal, pub) },
		func() Outcome { return checkSealKeyVersions(af) },
		func() Outcome { return checkSealKMSDigest(af) },
	} {
		if out := check(); out.Status == "FAIL" {
			return out
		}
	}
	return pass()
}

// checkSealSingleSignature checks the one signature this verifier can judge: signature_b64,
// or else the list's ed25519 entry. A seal with neither fails here, so a forged root stays an
// integrity finding, whatever an editor of the unsigned file does to the list.
func checkSealSingleSignature(seal AuditSeal, pub ed25519.PublicKey) Outcome {
	if seal.SignatureB64 == "" {
		seal.SignatureB64 = ed25519EntryOf(seal.SignaturesList)
	}
	return CheckSealSignatureV1(seal, pub)
}

// ed25519EntryOf returns the list's ed25519 signature, or "" when it holds none.
func ed25519EntryOf(list json.RawMessage) string {
	var entries []struct {
		Algorithm string `json:"algorithm"`
		Signature string `json:"signature"`
	}
	if err := json.Unmarshal(list, &entries); err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.Algorithm == sealAlgorithmEd25519 {
			return entry.Signature
		}
	}
	return ""
}

// refuseSealSignaturesList runs from JudgeSealPosture, once the ed25519 signature has verified, so
// no finding can hide behind it: anyone who can edit the seal file can add a list, and §10.12 routes
// exit 1 and exit 3 to different people.
// Passing would claim a dual-algorithm check that never ran. Failing with exit 1 would call
// a genuine dual-signed seal a tamper. So this is exit 3: the verifier cannot judge it.
func refuseSealSignaturesList(seal AuditSeal) Outcome {
	if seal.SignaturesList == nil {
		return pass()
	}
	return Outcome{Status: "FAIL", Step: "11", ExitCode: int(exitcodes.ConfigurationError),
		Reason: "seal carries a signatures list, and this verifier does not check one yet: dual-algorithm checks are not built"}
}

// checkSealCoversFile refuses a seal signed for another tenant or other HKDF inputs.
// The signature binds the seal's own fields, and only this comparison binds those
// fields to this file. It uses the step 11 reason because the signature is then
// genuine for some other seal, the same failure negative vector N005 pins.
func checkSealCoversFile(af *AuditFile) Outcome {
	if af.Seal.TenantID != af.Header.TenantID || af.Seal.HKDFInputsDigestHex != af.Header.HKDFInputsDigestHex {
		return fail("11", "signature verification failed")
	}
	return pass()
}

// checkSealAlgorithm uses the §4.3.2 reason. This is exit 1: the field is read before any
// signature is verified, so it must not move a forged seal out of exit 1.
func checkSealAlgorithm(seal AuditSeal) Outcome {
	if seal.Algorithm != sealAlgorithmEd25519 {
		return fail("11", fmt.Sprintf("algorithm %q not supported by this verifier (running v1.0c)", seal.Algorithm))
	}
	return pass()
}

// checkSealKeyVersions is the §7 step 11 cross-check. The signature can be genuine
// over a false list, so only the entries can show the list false. v1.0b and v1.0c
// carry the list as key_versions_canon. Any seal may also carry key_versions.
func checkSealKeyVersions(af *AuditFile) Outcome {
	const reason = "seal.key_versions does not match per-event key_version distribution"
	observed := observedKeyVersions(af)
	if bindsEntryLists(af.Seal) && af.Seal.KeyVersionsCanon != canonKeyVersions(observed) {
		return fail("11", reason)
	}
	if af.Seal.KeyVersions != nil && !slices.Equal(af.Seal.KeyVersions, observed) {
		return fail("11", reason)
	}
	return pass()
}

// checkSealKMSDigest applies §4.3: a flip in the day's kms_handle_uri values "produces
// a signature failure". The signature binds the seal's digest, so the digest must be
// the one these entries produce.
func checkSealKMSDigest(af *AuditFile) Outcome {
	if bindsEntryLists(af.Seal) && af.Seal.KMSHandleURIsDigestHex != kmsHandleURIsDigest(af) {
		return fail("11", "signature verification failed")
	}
	return pass()
}

func bindsEntryLists(seal AuditSeal) bool {
	return seal.SignPayloadVersion == signPayloadV1_0b || seal.SignPayloadVersion == signPayloadV1_0c
}

func observedKeyVersions(af *AuditFile) []int {
	versions := make([]int, 0, 1)
	for i := range af.Entries {
		versions = append(versions, af.Entries[i].KeyVersion)
	}
	sort.Ints(versions)
	return slices.Compact(versions)
}

func canonKeyVersions(versions []int) string {
	parts := make([]string, len(versions))
	for i, v := range versions {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

// kmsHandleURIsDigest is the §4.3 form: SHA-256 over the day's distinct URIs in byte
// order, joined by single newlines, with no trailing newline.
func kmsHandleURIsDigest(af *AuditFile) string {
	uris := make([]string, 0, 1)
	for i := range af.Entries {
		uris = append(uris, af.Entries[i].KMSHandleURI)
	}
	sort.Strings(uris)
	sum := sha256.Sum256([]byte(strings.Join(slices.Compact(uris), "\n")))
	return hex.EncodeToString(sum[:])
}

// SealAnomalies lists what a clean seal check did not establish.
func SealAnomalies(seal AuditSeal) []string {
	var anomalies []string
	if seal.DevMode {
		anomalies = append(anomalies, "seal is dev_mode: a development key signed it, so it is not a production seal: see §10.7")
	}
	if seal.SignPayloadVersion == signPayloadPreAmendment {
		anomalies = append(anomalies, "seal has no sign_payload_version, so its signature does not bind cadence or dev_mode: see §4.3")
	}
	if !bindsEntryLists(seal) && seal.KeyVersions == nil {
		anomalies = append(anomalies, "seal names no key_versions, so that cross-check did not run: see §7 step 11")
	}
	if seal.SignPayloadVersion == signPayloadV1_0c {
		anomalies = append(anomalies, fmt.Sprintf("operational_events_log_root %s is bound by the signature and was not walked: see §10.79", seal.OperationalEventsLogRoot))
	}
	return anomalies
}
