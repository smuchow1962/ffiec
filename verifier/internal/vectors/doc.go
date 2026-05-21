// Package vectors is the conformance-gate runner for the FFIEC
// chain-of-custody test-vector corpus at
// E:/dev/ffiec-public/spec/test-vectors/.
//
// What this package does:
//
//   - Locates the corpus via the FFIEC_PUBLIC_VECTORS_DIR env var
//     (default "../ffiec-public/spec/test-vectors", chosen for the
//     side-by-side repo layout that ffiec + ffiec-public ship under).
//     Decision D-2 locked 2026-05-21.
//   - Loads the master byte-level fixture (chain_vectors.json) and
//     decodes its inputs + expected outputs into typed values.
//   - Runs the verifier's primitives against the inputs and asserts
//     byte-equivalence against the expected outputs.
//
// What this package is for:
//
//   - A conformance gate. The Go test-vector consumer climbs from
//     0/N to N/N pass count as the verifier upgrade lands (Commits
//     1-6 per the Phase 2 plan). Every PR runs the gate; a green
//     bar is the merge condition.
//   - A drift detector. When a spec amendment lands new pinned
//     bytes in chain_vectors.json, the gate catches the drift on
//     the next CI run.
//   - A cross-implementation witness. The same JSON fixture is
//     consumed by the .NET reference at Herald.Compliance and the
//     Python reference at Herald.Py; three implementations
//     producing byte-identical output against the same fixture is
//     the byte-equivalence-vector discipline this spec ships under.
//
// What this package is NOT:
//
//   - A unit-test framework. Use Go's testing package; this is the
//     fixture loader + assertion shape that test files call.
//   - A negative-vector runner. Negative vectors (N001-N046) test
//     the verifier's failure reason strings; they land in this
//     package alongside the positive runner in Commit 4-5 as the
//     §7 step-by-step procedure lights up.
//   - A spec source. The spec lives at ffiec-public; this package
//     reads its byte-pinned fixtures, never authors them.
//
// Status: Commit 1 lands the loader + a runner skeleton that
// asserts the master-fixture bytes. The current verifier does not
// yet produce conformant output (it uses struct-order JSON, not
// RFC 8785 JCS; spec-conformance lands in Commit 2). The failing
// asserts in the conformance-gate test surface that gap exactly —
// the test is intentionally a red bar today and climbs to green as
// Commits 2-5 land.
package vectors
