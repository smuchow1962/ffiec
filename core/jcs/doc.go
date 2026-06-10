// Package jcs implements RFC 8785 JSON Canonicalization Scheme — the
// byte-exact canonical encoding the spec's per-event MAC, daily seal,
// and verdict-object compute over.
//
// Spec authority: §4.1 (canonical bytes are the MAC input), §4.2
// (canonical bytes are the Merkle leaf input), §10.12 verdict-object
// (canonical JCS bytes are the trailing line of every verifier run).
//
// Status: implemented (jcs.go + number.go). The 008-jcs-edge-cases
// corpus is the conformance gate (conformance_test.go) — byte-identical
// to the .NET reference's JcsRfc8785ConformanceTests answer key. The
// §7 pre-flight JCS self-test (baked-in vector 008 byte comparison)
// can wrap Canonicalize so a verifier built without conformant JCS
// refuses at startup with exit code 3 per §10.12.
//
// Implementation decision (D-1, locked 2026-05-21): stdlib-only,
// ~150 lines, mirrors the existing in-repo posture of writing HKDF
// inline rather than reaching for a third-party library. Vector 008
// is the conformance pin.
//
// Commit-2 implementation gotchas (per Richard's cross-consult
// 2026-05-21, lifting the .NET CanonicalJson.cs algorithm
// line-by-line):
//
//   - Key ordering MUST be byte-ordinal (`sort.Strings` over the
//     UTF-8 byte sequence of each key). Do NOT reach for
//     text/collate or any locale-aware sort — that would diverge
//     from .NET's StringComparer.Ordinal and Python's sorted(set(...)).
//   - Unicode escapes in string values MUST be lowercase `\u00XX`.
//     The encoding/json stdlib emits UPPERCASE `\u00XX` and is
//     therefore non-canonical; do NOT delegate to json.Marshal.
//   - Absent / nil fields MUST be omitted entirely from the output.
//     Emitting `"key":null` for an unset optional field is
//     non-canonical and breaks the v1.0c-optional
//     `operational_events_log_root` byte-pin.
//   - Integer formatting MUST go through strconv.FormatInt; the
//     stdlib's json.Marshal treats every JSON number as a float64
//     and that round-trip can lose precision for int64 values.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/CanonicalJson.cs and the Python
// reference share the byte-form contract. Any divergence between the
// three implementations on the same input is a conformance break.
package jcs
