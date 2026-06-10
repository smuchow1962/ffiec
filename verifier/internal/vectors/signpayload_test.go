package vectors

import (
	"errors"
	"testing"
)

// TestConformanceGate_SignPayloadCorpus wires the sign_payload vectors
// (018, 019, 020, 035 — all v1.0b) into the gate. For each, it
// reconstructs the sign_payload from input.json via core/signpayload
// and asserts byte-identity with the pinned expected_sign_payload.txt
// (and sha256 when present). This is the byte-equivalence proof for the
// sign_payload byte form against the .NET / Python references.
func TestConformanceGate_SignPayloadCorpus(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("sign_payload corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverSignPayloadVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverSignPayloadVectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatal("no sign_payload vectors discovered — corpus present but empty?")
	}

	t.Logf("sign_payload corpus: %d materialized vectors", len(vectors))

	totalChecks, totalPass, deferred := 0, 0, 0
	for _, v := range vectors {
		v := v
		t.Run(v.Slot, func(t *testing.T) {
			result := RunSignPayloadVector(v)
			if result.Deferred {
				deferred++
				t.Skipf("deferred: %s", result.Reason)
			}
			report := result.Report
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

	t.Logf("sign_payload corpus gate: %d/%d checks passed across %d vectors (%d deferred to backfill-Merkle tier)",
		totalPass, totalChecks, len(vectors)-deferred, deferred)
}
