package vectors

import (
	"errors"
	"testing"
)

// TestCorpusDir_ResolvesOrSkips proves the loader either resolves
// the corpus (developer machine has ffiec-public side-by-side) or
// returns ErrCorpusNotFound with a useful message (CI without the
// sibling repo, or a contributor who hasn't cloned it yet). Either
// outcome is acceptable for the Commit-1 baseline; the test asserts
// the error shape rather than the resolution itself.
func TestCorpusDir_ResolvesOrSkips(t *testing.T) {
	dir, err := CorpusDir()
	switch {
	case err == nil:
		if dir == "" {
			t.Fatal("CorpusDir returned nil error but empty dir")
		}
		t.Logf("conformance corpus resolved at %s", dir)
	case errors.Is(err, ErrCorpusNotFound):
		t.Logf("conformance corpus not found: %v", err)
		t.Logf("set %s to override; default is %s relative to test cwd",
			VectorsDirEnvVar, defaultVectorsDir)
	default:
		t.Fatalf("CorpusDir returned unexpected error: %v", err)
	}
}

// TestConformanceGate_MasterFixture is the load-bearing test of
// Commit 1. It loads the master fixture, runs every Commit-1-level
// conformance check, and asserts every check passes.
//
// What "every check passes" means at Commit 1:
//   - The fixture's HKDF constants match the in-repo constants.
//   - The HKDF-inputs digest recomputes byte-identical to the pin.
//   - Both IKM generations produce the expected session keys.
//   - Both IKM generations produce the expected key fingerprints.
//
// These checks do NOT depend on RFC 8785 JCS (which lands in Commit
// 2) so they are expected GREEN at Commit 1. The test thus
// establishes the conformance-gate discipline before any spec-
// dependent byte form is at risk; later commits extend the runner
// with the JCS-dependent checks and the gate climbs accordingly.
//
// When the corpus is not present on the build machine, the test
// skips with the resolution diagnostic. CI machines that intend to
// gate conformance MUST have ffiec-public checked out side-by-side
// or set FFIEC_PUBLIC_VECTORS_DIR; a skip on such a machine is a
// CI-config gap, not a test gap.
func TestConformanceGate_MasterFixture(t *testing.T) {
	fx, err := LoadMasterFixture()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("master fixture not available: %v", err)
	}
	if err != nil {
		t.Fatalf("LoadMasterFixture: %v", err)
	}

	report := RunMasterFixture(fx)
	t.Logf("conformance gate: %d/%d checks passed", report.PassCount(), len(report.Checks))

	for _, c := range report.Checks {
		if c.Pass {
			t.Logf("  [PASS] %s", c.Name)
		} else {
			t.Errorf("  [FAIL] %s — %s", c.Name, c.Detail)
		}
	}

	if !report.Pass() {
		t.Fatalf("conformance gate failed: %d of %d checks did not pass",
			len(report.Checks)-report.PassCount(), len(report.Checks))
	}
}
