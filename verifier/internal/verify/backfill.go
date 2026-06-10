package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/core/jcs"
)

// BaselineTuple is one entry in a §10.42 backfill baseline manifest: a
// {kind, identifier, sha256} record describing one inherited
// pre-acquisition artifact. The manifest is a JCS-canonical array of
// these, ordered lexicographically ascending by (kind, identifier) per
// §10.39 / §10.42.
type BaselineTuple struct {
	Identifier string `json:"identifier"`
	Kind       string `json:"kind"`
	SHA256     string `json:"sha256"`
}

// jcsBytes canonicalizes one tuple to its RFC 8785 bytes. JCS lex-sorts
// the keys, so the canonical bytes are {"identifier":…,"kind":…,
// "sha256":…} regardless of struct field order — the spec normates this
// lex-sort outcome explicitly in §10.42 step 1.
func (t BaselineTuple) jcsBytes() ([]byte, error) {
	return jcs.Canonicalize(map[string]any{
		"identifier": t.Identifier,
		"kind":       t.Kind,
		"sha256":     t.SHA256,
	})
}

// RecomputeBackfillMerkleRoot reconstructs the §10.42 backfill seal's
// Merkle root (the value bound at line 7 of the v1.0b sign_payload).
//
// The root covers two leaf classes, in order:
//
//  1. one leaf per baseline-manifest tuple — leaf payload is JCS(tuple);
//  2. one final metadata leaf — leaf payload is the already-canonical
//     JCS bytes of the §10.42 metadata object (the caller supplies these
//     bytes because the metadata leaf is a vector-pinned canonical form).
//
// Each leaf payload is hashed H(0x00 || payload) and combined under the
// RFC 6962 largest-power-of-2 split — the same construction §4.2 uses
// for the daily seal, so the backfill path reuses MerkleLeafHash +
// MerkleTreeHash with no Merkle-specific code of its own.
//
// Returns the root as lowercase hex (the form line 7 of the
// sign_payload carries).
func RecomputeBackfillMerkleRoot(tuples []BaselineTuple, metadataLeafJCS []byte) (string, error) {
	leaves := make([][]byte, 0, len(tuples)+1)
	for i := range tuples {
		payload, err := tuples[i].jcsBytes()
		if err != nil {
			return "", fmt.Errorf("canonicalize baseline tuple %d (%s): %w", i, tuples[i].Identifier, err)
		}
		leaves = append(leaves, MerkleLeafHash(payload))
	}
	leaves = append(leaves, MerkleLeafHash(metadataLeafJCS))

	root := MerkleTreeHash(leaves)
	return hex.EncodeToString(root), nil
}

// CheckBaselineManifestSHA256 implements §10.42 verifier-dispatch step 3:
// confirm the seal's seal.backfill_baseline_manifest_sha256 equals the
// SHA-256 of the JCS-canonical baseline-manifest ARRAY.
//
// This is a distinct hash from the Merkle root: the root is over the
// per-tuple leaves; this binding is over the whole canonicalized array
// as one document. A backfill seal whose bound manifest hash does not
// match the manifest content is a control-completeness anomaly (§10.42
// "mismatch is a control-completeness anomaly the verifier surfaces").
func CheckBaselineManifestSHA256(tuples []BaselineTuple, claimedHex string) error {
	canon, err := canonicalManifestArray(tuples)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(canon)
	got := hex.EncodeToString(sum[:])
	if got != claimedHex {
		return fmt.Errorf("baseline manifest sha256 mismatch: recomputed=%s claimed=%s", got, claimedHex)
	}
	return nil
}

// canonicalManifestArray builds the JCS-canonical bytes of the manifest
// as a single JSON array of tuple objects. The tuples MUST already be in
// (kind, identifier) lexicographic order per §10.39; this function does
// not re-sort — a mis-ordered manifest is itself a conformance break the
// caller's SHA-256 check surfaces.
func canonicalManifestArray(tuples []BaselineTuple) ([]byte, error) {
	arr := make([]any, 0, len(tuples))
	for i := range tuples {
		arr = append(arr, map[string]any{
			"identifier": tuples[i].Identifier,
			"kind":       tuples[i].Kind,
			"sha256":     tuples[i].SHA256,
		})
	}
	return jcs.Canonicalize(arr)
}
