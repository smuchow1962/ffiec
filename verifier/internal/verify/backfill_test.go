package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
)

// case035Tuples is the 8-tuple synthetic baseline manifest pinned by
// test vector 035-backfill-seal. Each sha256 is SHA-256(identifier);
// the test recomputes them so the fixture's recipe is exercised, not
// just copied.
func case035Tuples(t *testing.T) []BaselineTuple {
	t.Helper()
	tuples := make([]BaselineTuple, 0, 8)
	for n := 0; n < 8; n++ {
		id := identifier035(n)
		sum := sha256.Sum256([]byte(id))
		tuples = append(tuples, BaselineTuple{
			Identifier: id,
			Kind:       "baseline_diary_record",
			SHA256:     hex.EncodeToString(sum[:]),
		})
	}
	return tuples
}

func identifier035(n int) string {
	// cape-madeline-loan-orig-%04d
	const prefix = "cape-madeline-loan-orig-"
	return prefix + zeroPad4(n)
}

func zeroPad4(n int) string {
	s := []byte{'0', '0', '0', '0'}
	i := len(s) - 1
	for n > 0 && i >= 0 {
		s[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(s)
}

// case035MetadataLeafJCS canonicalizes the §10.42 metadata leaf to the
// exact 728-byte form pinned in 035's expected_canonical.txt. The
// baseline-manifest sha256 inside the leaf is recomputed from the tuple
// recipe so the leaf is derived, not transcribed.
func case035MetadataLeafJCS(t *testing.T) []byte {
	t.Helper()
	manifestSHA := case035ManifestSHA(t)
	leaf := map[string]any{
		"seal.backfill_at_close":                     true,
		"seal.backfill_baseline_manifest_sha256":     manifestSHA,
		"seal.backfill_companion_attestation_run_id": "northbridge-cape-madeline-close-2026-04-15",
		"seal.backfill_window_end_utc":               "2026-04-15T14:00:00Z",
		"seal.backfill_window_start_utc":             "2024-10-15T00:00:00Z",
		"seal.dual_signatures": []any{
			map[string]any{
				"entity_affiliation": "from_entity",
				"name":               "Helena R. Vasquez",
				"role":               "CISO",
				"signature_b64":      "VqK9NvDpFmK0Dc/oM7lY30+JsgcDEx8/UMDuFY3hjkJK/yNfjhMCk1ZQqPHt8FgLKp5dzGzCyc6f3Qd0XpLq8w==",
			},
			map[string]any{
				"entity_affiliation": "to_entity",
				"name":               "Marcus K. Tan",
				"role":               "CISO",
				"signature_b64":      "Qq2L4HtIfBp9KrM6w0aXyEcBkM8RzfL2VbCnDjY1nXh/oTpFxKsYbWmUrLkPvNcGdHgQzSjOaWeRtUiAsBfDqE==",
			},
		},
	}
	b, err := jcs.Canonicalize(leaf)
	if err != nil {
		t.Fatalf("canonicalize metadata leaf: %v", err)
	}
	return b
}

func case035ManifestSHA(t *testing.T) string {
	t.Helper()
	canon, err := canonicalManifestArray(case035Tuples(t))
	if err != nil {
		t.Fatalf("canonicalManifestArray: %v", err)
	}
	sum := sha256.Sum256(canon)
	return hex.EncodeToString(sum[:])
}

// TestRecomputeBackfillMerkleRoot pins the §10.42 backfill root recompute
// against 035's expected_merkle_root_hex (the value bound at line 7 of
// the v1.0b sign_payload). A divergence here means the Go Merkle / JCS
// path disagrees with the Python reference on the 9-leaf backfill tree.
func TestRecomputeBackfillMerkleRoot(t *testing.T) {
	const wantRoot = "8943b16ee4fdb413e849c6909c5d79b71fa96962a5710cdd6a763bf09342c340"

	got, err := RecomputeBackfillMerkleRoot(case035Tuples(t), case035MetadataLeafJCS(t))
	if err != nil {
		t.Fatalf("RecomputeBackfillMerkleRoot: %v", err)
	}
	if got != wantRoot {
		t.Errorf("backfill merkle root mismatch:\n got=%s\nwant=%s", got, wantRoot)
	}
}

// TestBackfillMetadataLeafCanonicalLength guards the metadata-leaf JCS
// form against the 728-byte pin. The leaf length is load-bearing: a
// drift in JCS key ordering or escaping shifts the leaf bytes, which
// shifts the Merkle root, which breaks every downstream backfill check.
func TestBackfillMetadataLeafCanonicalLength(t *testing.T) {
	leaf := case035MetadataLeafJCS(t)
	if len(leaf) != 728 {
		t.Errorf("metadata leaf canonical length = %d, want 728", len(leaf))
	}
}

// TestCheckBaselineManifestSHA256 confirms §10.42 step 3: the manifest-
// array SHA-256 binding matches the value 035 pins inside the metadata
// leaf. The happy path and a tampered-claim path are both exercised.
func TestCheckBaselineManifestSHA256(t *testing.T) {
	const wantManifestSHA = "880f875178fce4c3b55ee5503c755457d1b820e24a5a90d7ad2245797f9c8488"
	tuples := case035Tuples(t)

	if err := CheckBaselineManifestSHA256(tuples, wantManifestSHA); err != nil {
		t.Errorf("CheckBaselineManifestSHA256 (valid): unexpected error: %v", err)
	}

	bad := "deadbeef" + wantManifestSHA[8:]
	if err := CheckBaselineManifestSHA256(tuples, bad); err == nil {
		t.Error("CheckBaselineManifestSHA256 (tampered): expected mismatch error, got nil")
	}
}
