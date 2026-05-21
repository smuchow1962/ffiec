package jcs

import (
	"errors"
	"testing"
)

// TestSelfTest_PassesAgainstBakedInFixture is the meta-test that
// guards the §7 pre-flight: if SelfTest() returns nil here, the
// verifier built from this commit's source passes its own startup
// gate. If this test fails the verifier would refuse to start on
// any auditor's machine — surface the divergence here before
// shipping.
func TestSelfTest_PassesAgainstBakedInFixture(t *testing.T) {
	if err := SelfTest(); err != nil {
		t.Fatalf("SelfTest() returned %v; verifier would refuse to start at §7 pre-flight", err)
	}
}

// TestSelfTest_ErrorWrapsSentinel confirms a divergence error
// wraps ErrSelfTestFailed so callers can errors.Is against the
// sentinel. The verifier's CLI translates the sentinel into exit
// code 3 per §10.12 ConfigurationError; the wrapping discipline
// is what lets a future caller match on the sentinel without
// depending on the wrapping type.
//
// We construct a faked error rather than corrupt the package's
// internal selfTestExpected; the wrapping shape is what's under
// test, not the divergence-detection itself.
func TestSelfTest_ErrorWrapsSentinel(t *testing.T) {
	e := &selfTestError{gotBytes: []byte("different bytes")}
	if !errors.Is(e, ErrSelfTestFailed) {
		t.Fatal("selfTestError does not wrap ErrSelfTestFailed; CLI exit-code 3 dispatch would miss")
	}
}

// TestSelfTest_ErrorMessageNamesBothBytes confirms the divergence
// error message carries both the got and want byte sequences so a
// developer debugging the divergence sees the gap without re-
// running the canonicalizer in isolation.
func TestSelfTest_ErrorMessageNamesBothBytes(t *testing.T) {
	e := &selfTestError{gotBytes: []byte("XXXX")}
	msg := e.Error()
	if !contains(msg, "got XXXX") {
		t.Errorf("error message missing got-bytes diagnostic: %q", msg)
	}
	if !contains(msg, "want") {
		t.Errorf("error message missing want-bytes diagnostic: %q", msg)
	}
}

// contains is a tiny local substring check; reaching for strings.Contains
// here would be the right call too, but the file otherwise has no
// imports and the helper keeps the test file's import list to just
// `errors` and `testing`.
func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
