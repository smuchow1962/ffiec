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
// Commit-5 shape (locked 2026-05-21 per Steve's J-1 answer): full
// 6-field Verdict-Object — `additional_verifications`, `exit_code`,
// `posture`, `trust_anchor_manifest_sha256`,
// `verifier_spec_version_supported`, `verifier_version` — plus the
// v1.0c-optional `operational_events_log_root` field per §10.79.
// Six fields ship at v1.0 because the spec normates them as
// required and shipping less than spec-required is a conformance
// gap we won't accept. Critical path: Heather's case 036 extension
// + §10.12 normative-text lift land first (her amendment queue
// priority); Glenn's .NET + Python Verdict 6-field catchup lands
// the cross-impl byte-equivalence pin Commit 5 gates against.
//
// Per the locked release plan, the `verifier_version` field carries
// the implementation identifier shape "<binary-name>-v1.0.0-YYYY-MM-DD"
// — version + build date so an examiner reading the verdict can
// confirm vintage at a glance. The binary name lands at Commit 6
// when Steve picks the final verifier product name.
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
// Three-way symmetry across the Go reference verifier + Visus
// (Python; renamed from Vidimus 2026-05-21) + Herald.Compliance
// (.NET embedded) is the conformance bar.
//
// Cross-consult: Steve dispatched Richard 2026-05-21 for byte-form
// review of this surface against the .NET + Python references.
// Findings shape Commit 5 specifically.
package verdict
