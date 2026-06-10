package vectors

import (
	"errors"
	"testing"
)

// TestConformanceGate_NegativeCorpus wires the negative corpus
// (N001-N038) into the gate. For each row of the authoritative
// negative INDEX.md:
//
//   - materialized (input.json + expected_output.txt on disk): assert
//     the verifier emits the pinned Status/Step/Reason and the §10.12
//     exit code the target step implies.
//   - stub (description.md only): SKIP, recording the expected reason
//     so the conformance bar re-arms automatically when the fixture
//     lands.
//
// Today every negative is a stub (0/38 materialized), so this gate
// SKIPs all of them with their recorded expectations. The runner is
// built so the moment a fixture lands the gate asserts it with no code
// change — materialization state drives the behavior. The negative
// fixtures are a spec-side authoring artifact (Heather / spec-author
// lane), not a Go gap.
func TestConformanceGate_NegativeCorpus(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("negative corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverNegativeVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverNegativeVectors: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatal("no negative vectors parsed from INDEX.md — index present but empty?")
	}

	var asserted, skipped, requiredStubs int
	for _, v := range vectors {
		v := v
		t.Run(v.Slot, func(t *testing.T) {
			result := RunNegativeVector(v)
			if result.Skipped {
				skipped++
				if v.Required {
					requiredStubs++
				}
				t.Skipf("%s", result.SkipReason)
			}
			asserted++
			for _, c := range result.Report.Checks {
				if !c.Pass {
					t.Errorf("[FAIL] %s — %s", c.Name, c.Detail)
				}
			}
			if !result.Report.Pass() {
				t.Fatalf("negative vector %s: assertion failed", v.Slot)
			}
		})
	}

	t.Logf("negative corpus gate: %d asserted, %d skipped (%d required-but-stub) across %d INDEX rows",
		asserted, skipped, requiredStubs, len(vectors))
}

// TestNegativeIndex_ParsesAllRows is a guard on the INDEX parser: the
// shipped corpus enumerates N001-N038, so the parser must find a
// reasonable number of vector rows. A regression that breaks the table
// parse (e.g. an INDEX format change) would silently empty the gate;
// this catches it.
func TestNegativeIndex_ParsesAllRows(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("negative corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverNegativeVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverNegativeVectors: %v", err)
	}

	// The corpus ships N001-N038. Allow for index evolution but catch a
	// parse that finds almost nothing.
	const minExpected = 30
	if len(vectors) < minExpected {
		t.Errorf("parsed %d negative rows, expected at least %d — INDEX format may have drifted",
			len(vectors), minExpected)
	}

	// Every parsed row must carry a non-empty expected reason and a
	// recognizable target — a row with blanks means the column mapping
	// drifted.
	for _, v := range vectors {
		if v.ExpectedReason == "" {
			t.Errorf("vector %s parsed with empty expected reason", v.Slot)
		}
		if v.Target == "" {
			t.Errorf("vector %s parsed with empty target", v.Slot)
		}
	}
}
