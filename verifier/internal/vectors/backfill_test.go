package vectors

import (
	"errors"
	"testing"
)

// TestConformanceGate_BackfillCorpus wires the §10.42 backfill-seal
// vectors (035 today) into the gate. For each, it recomputes the seal's
// Merkle root from the baseline manifest + metadata leaf and asserts it
// matches expected_merkle_root_hex.txt, then recomputes the manifest-
// array SHA-256 and asserts the §10.42 step-3 binding. This is the
// load-bearing proof that the Go §10.42 recompute path reproduces the
// Python reference byte-for-byte.
func TestConformanceGate_BackfillCorpus(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("backfill corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverBackfillVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverBackfillVectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Skip("no backfill vectors materialized")
	}

	t.Logf("backfill corpus: %d materialized vector(s)", len(vectors))

	totalChecks, totalPass := 0, 0
	for _, v := range vectors {
		v := v
		t.Run(v.Slot, func(t *testing.T) {
			report := RunBackfillVector(v)
			totalChecks += len(report.Checks)
			totalPass += report.PassCount()
			for _, c := range report.Checks {
				if !c.Pass {
					t.Errorf("[FAIL] %s — %s", c.Name, c.Detail)
				}
			}
			if !report.Pass() {
				t.Fatalf("vector %s: %d of %d checks failed",
					v.Slot, len(report.Checks)-report.PassCount(), len(report.Checks))
			}
		})
	}

	t.Logf("backfill corpus gate: %d/%d checks passed across %d vector(s)",
		totalPass, totalChecks, len(vectors))
}
