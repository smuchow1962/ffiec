package jcs

import (
	"bytes"
	"errors"
)

// ErrSelfTestFailed is returned by SelfTest when the
// canonicalizer's output for the baked-in fixture diverges from the
// baked-in expected bytes. Per spec §7 pre-flight, a verifier built
// with a non-conformant JCS implementation MUST refuse all chain
// processing for the run with exit code 3 (configuration error).
//
// Wrap-aware: callers check `errors.Is(err, ErrSelfTestFailed)` to
// branch on the self-test failure without parsing the message.
var ErrSelfTestFailed = errors.New("verifier JCS self-test failed: canonicalization implementation does not conform to RFC 8785")

// selfTestInput is the smallest non-trivial JCS input that exercises
// the three areas most likely to silently disagree across
// implementations:
//
//   - Object key ordering (digits before uppercase, uppercase before
//     underscore, underscore before lowercase).
//   - String escape of one control character ( — no short
//     escape; verifies lowercase-hex output).
//   - Non-ASCII UTF-8 passthrough (no \uXXXX escape).
//
// The input is constructed in Go form (not parsed from JSON bytes)
// so the self-test is independent of the JSON parser; only the
// canonicalizer's output is under test.
//
// Per spec §7: "The verifier compiles both the fixture JSON and
// the expected canonical bytes into the verifier binary at build
// time so the self-test is independent of any runtime configuration
// or external file." This realization satisfies the discipline —
// the fixture is a Go literal initialized at first call, the
// expected bytes are a Go string literal next to it.
var selfTestInput = map[string]interface{}{
	"2":     int64(1),    // sorts AFTER "10" in UTF-16 code-unit order
	"10":    int64(2),    // sorts BEFORE "2"
	"a":     "café",      // non-ASCII passthrough
	"ctrl":  "\x0b",      // U+000B vertical tab — hex-escape, no short form
	"Upper": int64(3),    // uppercase sorts before lowercase
	"_":     int64(4),    // underscore sorts after uppercase, before lowercase
}

// selfTestExpected is the canonical output of selfTestInput under a
// conformant JCS implementation. Computed by writing the input
// through this package's Canonicalize on 2026-05-21 and confirmed
// against the vector 008 conformance corpus's edge cases.
//
// If this package's writeFloat64 / writeString / writeObject paths
// drift from RFC 8785 in any way the selfTestInput surfaces, the
// SelfTest call returns ErrSelfTestFailed and the verifier exits
// with code 3 at startup — auditor-readable failure mode.
const selfTestExpected = "{\"10\":2,\"2\":1,\"Upper\":3,\"_\":4,\"a\":\"café\",\"ctrl\":\"\\u000b\"}"

// SelfTest runs the baked-in fixture through Canonicalize and
// confirms byte-equality with the baked-in expected bytes. Returns
// nil on success; returns an error wrapping ErrSelfTestFailed on
// any divergence.
//
// Per spec §7 pre-flight: the verifier MUST run this once per
// process invocation before any chain verification. Callers that
// wrap this (the verifier's CLI in Commit 5) translate the error
// into exit code 3 (configuration error per §10.12) with the
// spec-mandated reason string.
//
// The cost is microseconds (one Canonicalize call against a
// 6-key fixture); the self-test runs at startup, not per file or
// per chain entry.
func SelfTest() error {
	got, err := Canonicalize(selfTestInput)
	if err != nil {
		return &selfTestError{inner: err, gotBytes: nil}
	}
	if !bytes.Equal(got, []byte(selfTestExpected)) {
		return &selfTestError{gotBytes: got}
	}
	return nil
}

// selfTestError wraps a self-test divergence with the actual bytes
// the canonicalizer produced; the verifier's startup error log can
// surface them so a developer debugging the divergence sees what
// changed without re-running the canonicalizer in isolation.
type selfTestError struct {
	inner    error  // non-nil when the divergence was an underlying canonicalizer error
	gotBytes []byte // populated when Canonicalize succeeded but produced wrong bytes
}

func (e *selfTestError) Error() string {
	if e.inner != nil {
		return ErrSelfTestFailed.Error() + ": canonicalizer error during self-test: " + e.inner.Error()
	}
	return ErrSelfTestFailed.Error() + ": got " + string(e.gotBytes) + " want " + selfTestExpected
}

// Unwrap returns the sentinel ErrSelfTestFailed so callers can
// errors.Is against it without depending on the wrapping type.
// The inner canonicalizer error is preserved for inspection via
// (*selfTestError).inner when callers need it.
func (e *selfTestError) Unwrap() error {
	return ErrSelfTestFailed
}
