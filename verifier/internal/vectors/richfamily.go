package vectors

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmpworks/ffiec/core/constants"
	"github.com/mmpworks/ffiec/core/hkdf"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// RichVector is one rich-expected-family vector: a corpus directory
// whose expected.json (or expected-result.json) carries assertions
// richer than the canonical-bytes / sign_payload pins. Each member has
// a bespoke shape, so the runner dispatches per slot to a check
// function rather than parsing a uniform schema.
type RichVector struct {
	Slot string
	Dir  string
}

// richFamilyDispatch maps a vector slot to its check function. Adding a
// member is one table row + one function — the runner loop stays
// schema-agnostic. Slots NOT in this table are not gated by the rich
// runner (they are either gated elsewhere or deferred with reasoning in
// the wave-2 PRD section).
var richFamilyDispatch = map[string]func(RichVector) *Report{
	"003-multi-run-same-day":              runMultiRunSameDay,
	"016-non-power-of-2-merkle":           runNonPowerOf2Merkle,
	"023-merkle-inclusion-proof-rfc6962":  runInclusionProof,
	"024-per-device-derivation":           runPerDeviceDerivation,
	"026-hierarchical-merkle-aggregation": runHierarchicalMerkle,
	"022-streaming-verifier-incremental":  runStreamingVerifier,
}

// DiscoverRichVectors returns the rich-family vectors present in the
// corpus that the runner knows how to gate (the richFamilyDispatch
// keys), each as a RichVector. Vectors in the table but absent from the
// corpus are silently skipped (the corpus is a moving target).
func DiscoverRichVectors(corpusDir string) []RichVector {
	var found []RichVector
	for slot := range richFamilyDispatch {
		dir := filepath.Join(corpusDir, slot)
		if dirExists(dir) {
			found = append(found, RichVector{Slot: slot, Dir: dir})
		}
	}
	return found
}

// RunRichVector dispatches one rich-family vector to its check function.
func RunRichVector(v RichVector) *Report {
	fn, ok := richFamilyDispatch[v.Slot]
	if !ok {
		r := &Report{}
		r.Checks = append(r.Checks, failCheck(v.Slot+"/dispatch", "no rich-family check registered for slot"))
		return r
	}
	return fn(v)
}

// ---- shared decode helpers ----

func readVectorJSON(dir, name string, into any) error {
	raw, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return fmt.Errorf("read %s: %w", name, err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("decode %s: %w", name, err)
	}
	return nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// decodeHexLeaves applies MerkleLeafHash to each hex-decoded preimage,
// returning the leaf-hashed bytes ready for MerkleTreeHash. This is the
// shared "wrap leaf preimages" step 016 and 003 both need.
func decodeHexLeaves(preimagesHex []string) ([][]byte, error) {
	leaves := make([][]byte, 0, len(preimagesHex))
	for i, h := range preimagesHex {
		b, err := hex.DecodeString(h)
		if err != nil {
			return nil, fmt.Errorf("decode leaf %d: %w", i, err)
		}
		leaves = append(leaves, verify.MerkleLeafHash(b))
	}
	return leaves, nil
}

// ---- 016: non-power-of-2 Merkle (§7 step 10, RFC 6962 odd-leaf) ----

func runNonPowerOf2Merkle(v RichVector) *Report {
	r := &Report{}

	var fixture struct {
		SubFixtures map[string]struct {
			LeavesHex     []string `json:"leaves_hex"`
			MerkleRootHex string   `json:"merkle_root_hex"`
		} `json:"sub_fixtures"`
	}
	if err := readVectorJSON(v.Dir, "fixture.json", &fixture); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	// Deterministic order: 3, 5, 7 leaves.
	for _, key := range []string{"3_leaves", "5_leaves", "7_leaves"} {
		sub, ok := fixture.SubFixtures[key]
		if !ok {
			r.Checks = append(r.Checks, failCheck(v.Slot+"/"+key, "sub-fixture absent"))
			continue
		}
		r.Checks = append(r.Checks, checkMerkleRootOverPreimages(v.Slot+"/"+key, sub.LeavesHex, sub.MerkleRootHex))
	}
	return r
}

// checkMerkleRootOverPreimages wraps each preimage as a leaf, builds the
// RFC 6962 root, and asserts it matches wantRootHex.
func checkMerkleRootOverPreimages(name string, preimagesHex []string, wantRootHex string) Check {
	leaves, err := decodeHexLeaves(preimagesHex)
	if err != nil {
		return failCheck(name, "%v", err)
	}
	got := hex.EncodeToString(verify.MerkleTreeHash(leaves))
	if got != wantRootHex {
		return failCheck(name, "root mismatch: got=%s want=%s", got, wantRootHex)
	}
	return passCheck(name)
}

// ---- 003: multi-run same-day (§7 step 10 + 11 + §4.1) ----

func runMultiRunSameDay(v RichVector) *Report {
	r := &Report{}

	var fx struct {
		Inputs struct {
			TenantID string `json:"tenant_id"`
			IKMV1Hex string `json:"ikm_v1_hex"`
		} `json:"inputs"`
		Expected struct {
			SessionKeyV1Hex     string `json:"session_key_v1_hex"`
			KeyFingerprintV1Hex string `json:"key_fingerprint_v1_hex"`
			MerkleRootDailyHex  string `json:"merkle_root_daily_hex"`
		} `json:"expected"`
		MerkleLeaves []struct {
			PayloadHashHex string `json:"payload_hash_hex"`
		} `json:"merkle_leaves_in_order"`
	}
	if err := readVectorJSON(v.Dir, "fixture.json", &fx); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	r.Checks = append(r.Checks, checkSessionKeyHex(v.Slot+"/session-key", fx.Inputs.TenantID, fx.Inputs.IKMV1Hex, fx.Expected.SessionKeyV1Hex))
	r.Checks = append(r.Checks, checkFingerprintHex(v.Slot+"/key-fingerprint", fx.Inputs.TenantID, fx.Inputs.IKMV1Hex, fx.Expected.KeyFingerprintV1Hex))

	// Daily Merkle root over all payload_hash values in (run_id, seq)
	// order — the leaves are pre-ordered in merkle_leaves_in_order.
	payloadHashes := make([]string, 0, len(fx.MerkleLeaves))
	for _, l := range fx.MerkleLeaves {
		payloadHashes = append(payloadHashes, l.PayloadHashHex)
	}
	r.Checks = append(r.Checks, checkMerkleRootOverPreimages(v.Slot+"/merkle-root-daily", payloadHashes, fx.Expected.MerkleRootDailyHex))
	return r
}

// checkSessionKeyHex recomputes HKDF-SHA-256(IKM, salt, info_for_tenant,
// 32) and asserts it matches the pinned hex.
func checkSessionKeyHex(name, tenantID, ikmHex, wantHex string) Check {
	ikm, err := hex.DecodeString(ikmHex)
	if err != nil {
		return failCheck(name, "decode ikm: %v", err)
	}
	info := []byte(constants.HKDFInfoBase + constants.InfoTenantSeparator + tenantID)
	got := hex.EncodeToString(hkdf.Derive(ikm, []byte(constants.HKDFSalt), info, constants.HKDFOutputLength))
	if got != wantHex {
		return failCheck(name, "session key mismatch: got=%s want=%s", got, wantHex)
	}
	return passCheck(name)
}

// checkFingerprintHex recomputes SHA-256(utf8(tenant_id) || ikm)[:16]
// and asserts it matches the pinned hex.
func checkFingerprintHex(name, tenantID, ikmHex, wantHex string) Check {
	ikm, err := hex.DecodeString(ikmHex)
	if err != nil {
		return failCheck(name, "decode ikm: %v", err)
	}
	got := verify.KeyFingerprintHex(tenantID, ikm)
	if got != wantHex {
		return failCheck(name, "fingerprint mismatch: got=%s want=%s", got, wantHex)
	}
	return passCheck(name)
}
