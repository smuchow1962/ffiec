// Package verdict implements the §10.12 + §10.29 verifier verdict
// object — the JCS-canonical machine-readable supplement that every
// verifier run emits as its trailing line per the §7
// `Verdict-Object: <jcs-bytes>` discipline.
//
// Spec authority:
//
//   - §10.12 — verifier CLI exit-code contract (0=PASS, 1=FAIL,
//     2=structural/input error, 3=configuration error) plus the
//     closed `additional_verifications` enumeration of bonus-
//     verification markers.
//   - §10.29 — streaming-mode verifier state codes (4=streaming
//     all-pass, 5=streaming anomaly-detected, 6=streaming
//     key-rotation pending). Non-terminal; only emitted by
//     streaming-mode verifiers.
//   - §7 — the `Verdict-Object: <jcs-bytes>` trailing-line
//     discipline; the verdict object schema's six required fields
//     plus the optional operational_events_log_root.
//
// Status: stub. The verdict-object types + JCS-canonical writer
// land in Commit 5 of the verifier upgrade. The closed-enumeration
// lookup table for `additional_verifications` markers lives here
// alongside the verdict struct.
//
// Commit-5 default shape (per Richard's cross-consult 2026-05-21):
// 2-field verdict object — `additional_verifications` + `exit_code`
// — matching what the .NET + Python references emit today and
// gating against vector 036 sub-cases (036a/b/c) without breaking
// cross-implementation byte-pin tests. The four spec-normative
// non-optional fields (`posture`, `verifier_version`,
// `verifier_spec_version_supported`, `trust_anchor_manifest_sha256`)
// and the v1.0c-optional `operational_events_log_root` are scaffold
// TODOs with explicit comment markers; Steve's open call on
// 2-vs-6-field shape decides whether they ship in Commit 5 or wait
// for fixture 036d to materialize cross-impl convergence.
//
// The 18-marker `additional_verifications` enumeration ships in
// full at Commit 5 even though .NET + Python carry only 1 today.
// The enum is a validator (`IsKnownMarker`), not a serializer, so
// the larger Go-side set does not affect byte-equivalence — it
// just lets a future Go verifier dispatching one of the other 17
// markers do so without an enum-table backfill.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/VerifierVerdict.cs explicitly
// states it "mirrors the Python reference at
// Herald.Py/src/herald/_verdict.py; both implementations agree
// byte-for-byte on the integer exit codes, the streaming/terminal
// partition, and the structured verdict's JCS-canonical form."
// Three-way symmetry is the conformance bar.
//
// Cross-consult: Steve dispatched Richard 2026-05-21 for byte-form
// review of this surface against the .NET + Python references.
// Findings shape Commit 5 specifically.
package verdict
