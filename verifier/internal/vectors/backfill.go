package vectors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mmpworks/ffiec/core/jcs"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// backfillRootFile is the §10.42 backfill seal's pinned Merkle root —
// the value bound at line 7 of the v1.0b sign_payload, recomputed from
// the baseline manifest + metadata leaf.
const (
	backfillRootFile     = "expected_merkle_root_hex.txt"
	syntheticManifestKey = "synthetic_baseline_manifest"
	metadataLeafKey      = "metadata_leaf"
)

// BackfillVector is one materialized §10.42 backfill-seal vector: a
// corpus directory carrying a synthetic_baseline_manifest + metadata_leaf
// in input.json plus expected_merkle_root_hex.txt. The vector pins the
// recompute of the seal's Merkle root from the baseline manifest, which
// is the §10.42 verifier-dispatch step-1 proof.
type BackfillVector struct {
	Slot string
	Dir  string
}

// DiscoverBackfillVectors walks the corpus and returns every vector
// carrying both a synthetic_baseline_manifest in input.json and the
// expected_merkle_root_hex.txt pin (the §10.42 backfill family).
func DiscoverBackfillVectors(corpusDir string) ([]BackfillVector, error) {
	entries, err := os.ReadDir(corpusDir)
	if err != nil {
		return nil, fmt.Errorf("read corpus dir %s: %w", corpusDir, err)
	}

	var found []BackfillVector
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(corpusDir, e.Name())
		if !fileExists(filepath.Join(dir, backfillRootFile)) {
			continue
		}
		if !fileExists(filepath.Join(dir, "input.json")) {
			continue
		}
		found = append(found, BackfillVector{Slot: e.Name(), Dir: dir})
	}
	return found, nil
}

// backfillInput is the part of a backfill vector's input.json the
// §10.42 recompute consumes: the baseline manifest tuples (and their
// bound array hash) plus the metadata leaf the recompute appends as the
// final Merkle leaf.
type backfillInput struct {
	MetadataLeaf         json.RawMessage `json:"metadata_leaf"`
	SyntheticBaselineMan struct {
		Tuples []verify.BaselineTuple `json:"tuples"`
		SHA256 string                 `json:"sha256"`
	} `json:"synthetic_baseline_manifest"`
}

// RunBackfillVector executes the §10.42 backfill verifier-dispatch path
// against one vector and reports the recompute outcomes:
//
//   - step 1: recompute the Merkle root over (baseline tuples + metadata
//     leaf) and assert it matches expected_merkle_root_hex.txt;
//   - step 3: recompute the baseline-manifest array SHA-256 and assert
//     it matches the value bound inside the metadata leaf.
//
// This is the load-bearing §10.42 conformance proof: a clean-room
// verifier must reproduce both the root and the manifest hash byte-for-
// byte from the manifest content alone.
func RunBackfillVector(v BackfillVector) *Report {
	r := &Report{}

	in, err := loadBackfillInput(v.Dir)
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read-input", "%v", err))
		return r
	}

	metadataLeafJCS, err := canonicalizeMetadataLeaf(in.MetadataLeaf)
	if err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/canonicalize-metadata-leaf", "%v", err))
		return r
	}

	r.Checks = append(r.Checks, checkBackfillRoot(v, in.SyntheticBaselineMan.Tuples, metadataLeafJCS))
	r.Checks = append(r.Checks, checkBackfillManifestSHA(v, in.SyntheticBaselineMan.Tuples))
	return r
}

func loadBackfillInput(dir string) (backfillInput, error) {
	raw, err := os.ReadFile(filepath.Join(dir, "input.json"))
	if err != nil {
		return backfillInput{}, fmt.Errorf("read input.json: %w", err)
	}
	var in backfillInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return backfillInput{}, fmt.Errorf("decode input.json: %w", err)
	}
	if len(in.SyntheticBaselineMan.Tuples) == 0 {
		return backfillInput{}, fmt.Errorf("input.json has no %s.tuples", syntheticManifestKey)
	}
	if len(in.MetadataLeaf) == 0 {
		return backfillInput{}, fmt.Errorf("input.json has no %s", metadataLeafKey)
	}
	return in, nil
}

// canonicalizeMetadataLeaf decodes the metadata_leaf JSON object and
// re-canonicalizes it through core/jcs so the leaf bytes the recompute
// hashes are the verifier's own JCS output, not a transcription of the
// fixture's pretty-printed form. This is the verifier-side proof that
// the Go JCS path reproduces the §10.42 metadata-leaf bytes.
func canonicalizeMetadataLeaf(raw json.RawMessage) ([]byte, error) {
	var leaf any
	if err := json.Unmarshal(raw, &leaf); err != nil {
		return nil, fmt.Errorf("decode metadata_leaf: %w", err)
	}
	return jcs.Canonicalize(leaf)
}

func checkBackfillRoot(v BackfillVector, tuples []verify.BaselineTuple, metadataLeafJCS []byte) Check {
	name := v.Slot + "/backfill-merkle-root"
	want, err := readBackfillRoot(v.Dir)
	if err != nil {
		return failCheck(name, "%v", err)
	}
	got, err := verify.RecomputeBackfillMerkleRoot(tuples, metadataLeafJCS)
	if err != nil {
		return failCheck(name, "recompute root: %v", err)
	}
	if got != want {
		return failCheck(name, "root mismatch: recomputed=%s want=%s", got, want)
	}
	return passCheck(name)
}

// checkBackfillManifestSHA asserts §10.42 step 3: the manifest-array
// SHA-256 matches the value bound inside the metadata leaf. The bound
// value is read from the metadata leaf's
// seal.backfill_baseline_manifest_sha256, surfaced redundantly as
// synthetic_baseline_manifest.sha256 in the fixture.
func checkBackfillManifestSHA(v BackfillVector, tuples []verify.BaselineTuple) Check {
	name := v.Slot + "/backfill-manifest-sha256"
	in, err := loadBackfillInput(v.Dir)
	if err != nil {
		return failCheck(name, "%v", err)
	}
	if in.SyntheticBaselineMan.SHA256 == "" {
		return failCheck(name, "fixture omits synthetic_baseline_manifest.sha256")
	}
	if err := verify.CheckBaselineManifestSHA256(tuples, in.SyntheticBaselineMan.SHA256); err != nil {
		return failCheck(name, "%v", err)
	}
	return passCheck(name)
}

func readBackfillRoot(dir string) (string, error) {
	raw, err := os.ReadFile(filepath.Join(dir, backfillRootFile))
	if err != nil {
		return "", fmt.Errorf("read %s: %w", backfillRootFile, err)
	}
	return strings.TrimSpace(string(raw)), nil
}
