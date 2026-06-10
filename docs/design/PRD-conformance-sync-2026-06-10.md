# PRD — TesseraSeal Go verifier: conformance-corpus sync + PRD-3 re-target

**Date:** 2026-06-10
**Author:** Jared (Go/Rust systems engineer)
**Status:** In progress
**Repo:** `E:\dev\ffiec` (Go verifier) against spec `E:\dev\ffiec-public\spec\chain-of-custody-DRAFT-0.3.0.md`

---

## What

Bring the Go verifier current with the rest of the TesseraSeal chain-of-custody
ecosystem (the Herald.Compliance C# reference + the Herald.Py Python reference)
and wire the materialized conformance corpus into the per-commit test gate.

Concretely, four coherent steps:

1. Add a real **RFC 8785 JCS canonicalizer** to `core/jcs` (today a doc-only
   stub; the verifier uses stdlib `json.Encoder` which is *not* JCS). This is
   the load-bearing primitive every canonical-output vector pins against.
2. Wire the **positive canonical-output corpus** (`expected_canonical.txt` +
   `expected_canonical_sha256.txt` family, ~23 vectors) into the conformance
   gate as a corpus-walk runner.
3. Add **sign_payload_version dispatch** (`core/signpayload`): the 6-line
   pre-amendment / 10-line v1.0a / 12-line v1.0b / 13-line v1.0c forms, and
   wire the sign_payload vectors (018, 019, 020, 035).
4. Wire the **negative corpus** as exit-code assertions where vectors are
   materialized; emit SKIP (not FAIL) for the stub majority, listed explicitly.

## Why

The morning research spike confirmed: repo builds/vets/tests green at main
`@ 64ba771`, but the **conformance corpus is not wired to the gate** and the
verifier predates spec PRD-3 / 0.3.0. Steve's directive: pull all Herald.OSS /
Herald.Compliance currency into TesseraSeal, accommodate all the vectors, use
the new vector DSL.

The verifier is the artifact auditors run and regulators read. "Conformant"
means it produces byte-identical output to the C# and Python references for the
same input — and the **only** thing that enforces that three-language agreement
is the shared test-vector corpus. Wiring the corpus into the gate is what turns
"we wrote a Go verifier" into "the Go verifier is a conformant third
implementation."

## The sync relationship (discovered, documented — this was the architectural fork)

**There is no automated mirror/sync mechanism between Herald.Compliance and the
Go verifier, and there should not be one.** The ecosystem's "byte-identical
Core↔OSS mirror" pattern is a *.NET-internal* device (two C# assemblies kept in
lockstep). It does not — and must not — extend across the language boundary to Go.

The sync contract across C# / Python / Go is **the shared spec + the shared
conformance corpus**:

- All three implement the same `chain-of-custody-DRAFT-*.md`.
- All three must produce byte-identical output against
  `E:\dev\ffiec-public\spec\test-vectors\`.
- The corpus IS the contract. A divergence between any two implementations on
  the same vector is a conformance break, caught by the corpus, not by a sync script.

This is the correct CUPID-Composable shape: the corpus is the seam; each
implementation is independently substitutable behind it. A code-level sync would
couple the Go verifier to C# internals it has no business knowing.

**Decision D-SYNC-1 (reversible):** Do not build a Herald→ffiec code-sync
mechanism. Treat the corpus as the contract. If a future need for shared
*fixtures* (not code) emerges, add a fixture-vendoring manifest then — not now.

## Recent Herald.Compliance changes since the last ffiec commit (2026-05-23)

Surveyed `Herald/Modules/Herald.Compliance/src/Audit/Chain/` git log since
2026-05-20. The only chain commit is `c682f79` (LicenseStateTransitionChainEntry
audit-chain sink + T-FFI-7) — a *new event-type sink*, not a change to canonical
output, MAC, Merkle, or sign_payload behavior. The C# `CanonicalJson.cs` byte
contract is stable and is the line-by-line reference for the Go JCS port.

Net: no canonical-output drift to chase from the C# side. The Go-side gap is
purely **the verifier never had real JCS or sign_payload dispatch wired**, plus
the corpus never being gated.

## The vector DSL (found)

Steve said "we now have a DSL to handle new vector work." **The DSL is the
per-vector `_compute.py` generator + the structured `input.json` /
`expected_canonical.txt` / `expected_canonical_sha256.txt` materialization
convention**, indexed by `spec/test-vectors/PRD-4-INDEX.md` and the negative
`spec/test-vectors/negative/INDEX.md`.

Each materialized vector is authored by writing a `_compute.py` that emits the
canonical bytes deterministically; the `expected_canonical_sha256.txt` first row
is `sha256(expected_canonical.txt)`, and embedded `*_canonical_bytes_utf8` values
inside `input.json` carry their own sha256 pins. This is a
generator-and-pin convention, not a bespoke parser-grammar DSL. It is the
canonical authoring path; the Go runner *consumes* it, it does not re-author it.

(If a richer grammar-style DSL exists elsewhere, I did not find it. Reported
honestly per directive.)

## Corpus schema map (what the gate must consume)

| Family | Expected file(s) | Vectors | Gate assertion |
|---|---|---|---|
| Master fixture | `chain_vectors.json` | (1) | Already gated by `RunMasterFixture` (HKDF/fingerprint/digest) |
| Canonical-output | `expected_canonical.txt` + `_sha256.txt` | 021, 025, 034, 036–053 (~23) | JCS idempotence + sha256 self-consistency + embedded-bytes pins |
| Sign-payload | `expected_sign_payload.txt` + `_sha256.txt` | 018, 019, 020, 035 | sign_payload byte-form reconstruction == pinned bytes |
| Rich-expected | `expected.json` / `expected-result.json` | 003, 008, 016, 017, 022, 023, 024, 026 | Per-vector (008 = JCS edge-case categories) |
| Stub | (none; description.md only) | 001, 002, 010, 015, 027 | SKIP, listed |
| Negative | (none materialized; description.md only) | N001–N038 | SKIP all, listed; runner ready for when they materialize |

The **canonical-output family is the highest-leverage wiring** because it is
schema-uniform across 23 vectors and self-pinning: the runner does not need to
parse each vector's bespoke embedded object. It asserts:

1. `sha256(expected_canonical.txt)` == first row of `expected_canonical_sha256.txt`
   (self-consistency — catches corpus corruption).
2. Re-canonicalizing the parsed `expected_canonical.txt` through the Go JCS
   implementation reproduces it byte-for-byte (**JCS idempotence — this is the
   verifier-side proof** that Go agrees with whatever produced the pin).
3. Each `*_canonical_bytes_utf8` in `input.json` hashes to its declared
   `*_canonical_sha256` and is JCS-idempotent.

## Acceptance

- `core/jcs` implements RFC 8785; passes the 008-jcs-edge-cases corpus
  (byte-identical to the C# `JcsRfc8785ConformanceTests` answer key).
- A positive corpus-walk runner asserts the canonical-output family (≥23
  materialized vectors green); SKIPs stubs with an explicit list.
- `core/signpayload` reconstructs v1.0a/b/c forms; sign_payload vectors green.
- A negative-corpus runner is in place; asserts materialized negatives (0 today)
  and SKIPs N001–N038 with the exit-code expectations recorded for when they land.
- Both runners wired into `go test` so every commit runs every materialized vector.
- `go build` + `go vet` + `go test` green per module after each step.
- No DRY violation: the verifier's existing ad-hoc `marshalNoNewline` /
  `json.Marshal`-for-sort paths converge on `core/jcs` where they need JCS.

## Non-goals

- Building a Herald→ffiec code-sync script (D-SYNC-1: the corpus is the contract).
- Materializing negative vectors N001–N038 (spec-side artifact authoring; out of
  the Go lane — flagged for Heather/spec-author coordination).
- Materializing stub positive vectors (001, 002, 010, 015, 027 — same).
- Full §7 step 1–12 walk over a real multi-event chain file for every primitive
  (the canonical-output + sign_payload families gate the byte foundation; the
  full chain-walk over `input.json` chain entries is a follow-on once the
  rich-expected family fixtures stabilize).
- Rust. No benchmark shows Go is the bottleneck here; JCS is allocation-light.

## Risk / fork log

- **D-SYNC-1** (above): no code-sync; corpus is the contract. Reversible.
- **JCS number formatting**: ECMA-262 NumberToString is the subtle part. The Go
  port mirrors the C# `FormatEcma262Number` algorithm. The 008 corpus float
  cases are the witness; if any diverge, that is the first place to look.
- **Stub-heavy negative corpus**: 0/38 negatives materialized. The runner is
  built to gate them the moment fixtures land, but today the negative
  conformance bar is "SKIP with recorded expectations," not "assert." This is a
  spec-side materialization dependency, not a Go gap — flagged for Steve.
