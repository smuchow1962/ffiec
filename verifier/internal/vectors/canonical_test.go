package vectors

import (
	"errors"
	"testing"
)

// TestConformanceGate_CanonicalCorpus is the positive corpus-walk gate.
// It discovers every materialized canonical-output vector (those
// carrying expected_canonical.txt + expected_canonical_sha256.txt) and
// runs the three-tier conformance assertion against each:
//
//  1. self-consistency: sha256(expected_canonical.txt) == primary pin.
//  2. JCS idempotence: re-canonicalizing through core/jcs reproduces
//     the bytes — the verifier-side proof of agreement with the .NET /
//     Python references.
//  3. embedded-bytes pins: each *_canonical_bytes_utf8 in input.json
//     hashes to its sibling sha and is JCS-idempotent.
//
// This is the wiring that turns the corpus into a per-commit gate: any
// divergence between the Go canonicalizer and the pinned reference
// bytes fails the build. When the corpus is absent the test skips with
// the resolution diagnostic (a CI-config gap, not a test gap).
func TestConformanceGate_CanonicalCorpus(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("canonical corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverCanonicalVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverCanonicalVectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatal("no canonical-output vectors discovered — corpus present but empty?")
	}

	t.Logf("canonical corpus: %d materialized vectors", len(vectors))

	totalChecks, totalPass := 0, 0
	for _, v := range vectors {
		v := v
		t.Run(v.Slot, func(t *testing.T) {
			report := RunCanonicalVector(v)
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

	t.Logf("canonical corpus gate: %d/%d checks passed across %d vectors",
		totalPass, totalChecks, len(vectors))
}
