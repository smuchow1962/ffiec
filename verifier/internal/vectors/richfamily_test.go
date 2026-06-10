package vectors

import (
	"errors"
	"sort"
	"testing"
)

// TestConformanceGate_RichFamily wires the rich-expected family (003,
// 016, 022, 023, 024, 026) into the gate. Each member carries a bespoke
// expected.json shape — daily Merkle aggregation, RFC 6962 odd-leaf
// roots, inclusion proofs, per-device HKDF, hierarchical aggregation,
// the streaming state machine — and the runner dispatches per slot to
// the matching check function. 017 is NOT here: it is a placeholder
// vector (PLACEHOLDER:… hex strings, no real bytes) whose normative
// content is the partial-disclosure output shape, deferred with spec
// reasoning in the wave-2 PRD section.
func TestConformanceGate_RichFamily(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("rich-family corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors := DiscoverRichVectors(dir)
	if len(vectors) == 0 {
		t.Skip("no rich-family vectors materialized")
	}
	sort.Slice(vectors, func(i, j int) bool { return vectors[i].Slot < vectors[j].Slot })

	t.Logf("rich-family corpus: %d materialized vector(s)", len(vectors))

	totalChecks, totalPass := 0, 0
	for _, v := range vectors {
		v := v
		t.Run(v.Slot, func(t *testing.T) {
			report := RunRichVector(v)
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
			t.Logf("%s: %d/%d checks", v.Slot, report.PassCount(), len(report.Checks))
		})
	}

	t.Logf("rich-family corpus gate: %d/%d checks passed across %d vector(s)",
		totalPass, totalChecks, len(vectors))
}
