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

---

## Realized outcomes (2026-06-10)

All five steps landed, build/vet/test green per module after each. Five local
commits on `main` (NOT pushed — Steve reviews before push):

| Commit | What |
|---|---|
| `6d37bac` | `core/jcs` RFC 8785 canonicalizer + 008 conformance gate (full 008 corpus green: float canon, non-ASCII unicode, surrogate pairs, deep-nest 10/50/100, all control chars, NaN/Inf rejection, very long strings) |
| `ff7bd3b` | Positive canonical-output corpus gate: **52/52 checks across 22 vectors** |
| `d72951a` | `core/signpayload` v1.0a/b/c dispatch (PRD-3 gap) + gate: **3 vectors green** (018/019/020), 035 deferred |
| `3bf9003` | Negative corpus gate: **38/38 INDEX rows parsed, 38 SKIP** (0 materialized) |
| `f0175d7` | `delegation_chain` sort migrated json.Marshal → core/jcs (Herald parity) + regression test |

**Gate totals:** canonical 52/52 across 22 vectors; sign_payload 6/6 across 3
(1 deferred); negative 38 rows parsed/38 skipped; JCS 008 corpus full green;
HKDF RFC 5869 green (pre-existing); master fixture green (pre-existing).

**Deferred / flagged for Steve:**
- **Negative fixtures N001–N038** are spec-side authoring (Heather / spec-author
  lane). The Go gate auto-arms when they land — no Go change needed.
- **Stub positive vectors** 001, 002, 010, 015, 027 — same spec-side dependency.
- **035-backfill-seal sign_payload** needs the backfill-Merkle recompute path to
  reconstruct its derived merkle_root; deferred to that follow-on, not failing.
- **Rich-expected family** (003, 008, 016, 017, 022, 023, 024, 026) — 008 is
  gated via the JCS corpus; the others carry per-vector `expected.json` shapes
  that gate on full §7 chain-walk machinery (a follow-on once those fixtures'
  walk surfaces stabilize). Not wired this pass; listed for completeness.
- **Race detector** unavailable on this box (`-race` needs cgo/gcc). Verifier
  gate path is single-threaded so race exposure is nil; non-race suite green.

---

# Wave 2 — §7 chain-walk / backfill-Merkle recompute + rich-expected family

**Date:** 2026-06-10 (same day, second wave)
**Status:** In progress
**Spec:** `chain-of-custody-DRAFT-0.3.0.md` (PRD-3) — §7 verification + §10.42 backfill seal

## What

Two deferred items from wave 1, both approved by Steve:

1. **The §10.42 backfill-Merkle recompute path** so the 035-backfill-seal
   sign_payload check stops being deferred. Today `resolveSealFields` returns
   `errDeferredReconstruction` for 035 because its `merkle_root` is *derived*
   from a baseline manifest, not pinned directly in the fixture. The recompute
   path reconstructs the root from the manifest + metadata leaf, then feeds it
   into the existing v1.0b reconstruction.

2. **Gate the rich-expected family** (003, 016, 017, 022, 023, 024, 026) — the
   vectors carrying bespoke `expected.json` shapes richer than the
   canonical-bytes / sign_payload pins.

## Why

The 035 deferral was honest but incomplete: the verifier could reconstruct any
*directly-pinned* v1.0b seal but not a backfill seal, which is the one §10.42
audit-trail product feature an M&A examiner actually runs. The rich family is
where the spec pins its harder Merkle shapes — odd-leaf trees (016),
multi-run/same-day daily aggregation (003), RFC 6962 inclusion proofs
(017/023), per-device HKDF (024), hierarchical aggregation (026), and the
streaming state machine (022). Each is a normative byte (or transition) the
corpus pins; leaving them ungated means the Go verifier silently might diverge
from the .NET/Python references on exactly the cases auditors care about.

## §10.42 backfill recompute — the exact construction (from 035 `_compute.py`)

The backfill seal's `merkle_root` (line 7 of `expected_sign_payload.txt`,
`8943b16e…`) is the RFC 6962 root over **9 leaf payloads**:

- 8 baseline-manifest leaves: each is `JCS(tuple)` where `tuple =
  {identifier, kind, sha256}` (JCS lex-sorts the keys).
- + 1 metadata leaf: `JCS(metadata_leaf)` — the 6-field §10.42 attribute object.

Each leaf payload is hashed `H(0x00 || payload)` then combined under the
`largest-power-of-2-strictly-less-than-n` RFC 6962 split — exactly the existing
`MerkleLeafHash` + `MerkleTreeHash` in `verify/merkle.go`. No new Merkle code:
the recompute reuses the wave-1 `core/jcs` canonicalizer for the leaf payloads
and the verifier's Merkle primitives for the root.

The §10.42 verifier dispatch is a 5-step path (spec §10.42 "Verifier dispatch"):

1. Recompute the Merkle root over the canonicalized baseline manifest
   (+ metadata leaf, per this vector's binding-through-metadata-leaf shape).
2. (signature — out of this vector's scope; 035 pins the byte form, not a sig.)
3. Confirm `seal.backfill_baseline_manifest_sha256` == `sha256(JCS(manifest
   array))` — the manifest-array hash binding, a *separate* hash from the root.
4. (window/companion structural — present in the metadata leaf.)
5. (companion §10.39 attestation cross-ref — case 034, not in 035's scope.)

For 035 the load-bearing checks are step 1 (root recompute feeds the
sign_payload) and step 3 (manifest-array sha256 == the metadata-leaf's bound
value). Steps 2/4/5 need the companion 034 event + a real signature, neither
present in 035 — so the gate asserts 1 and 3 and records 2/4/5 as
out-of-vector-scope (NOT deferred-because-unbuilt; deferred-because-not-pinned-
here). When 034 materializes with a paired backfill, the cross-ref lands.

## Rich-expected family — §7-scope classification

Each member assessed against "is this §7-mandated, or a §10.x extension?":

| Vector | Shape | §7 scope | Build decision |
|---|---|---|---|
| 003 multi-run-same-day | daily Merkle over 6 payload_hash in (run_id,seq) + session key + fingerprint + hkdf-digest + v1.0a sign_payload | **§7 step 10 + 11 + §4.1** — fully mandated | Build: all existing primitives |
| 016 non-power-of-2-merkle | RFC 6962 odd-leaf right-promote roots (3/5/7 leaves) | **§7 step 10** — mandated | Build: direct `MerkleTreeHash` |
| 024 per-device-derivation | §10.32 HKDF info with 3rd device segment | §10.32 extension (witness-applicable) | Build: trivial on existing `hkdf.Derive` |
| 023 merkle-inclusion-proof | §10.31 RFC 6962 audit-path verification (5-leaf) | §10.31 extension (witness-applicable per §7 table) | Build: new audit-path fold (small, CUPID-composable) |
| 017 merkle-inclusion-partial-disclosure | §10.31 audit path, partial-disclosure framing | §10.31 extension | Build: same fold as 023 |
| 026 hierarchical-merkle-aggregation | §10.37 inner+outer trees, concatenated audit paths | §10.37 extension | Build: top-tree-no-leaf-hash + reuse 023 fold |
| 022 streaming-verifier-incremental | §10.29 streaming state machine (transition + finalize tables) | §10.29 extension (no crypto) | Build: self-contained state machine |

**Decision D-W2-1:** Build all seven rather than defer the §10.x members. Each
is small, self-contained, and the corpus pins a normative byte/transition the
Go verifier must reproduce to claim conformance. The §7-vs-§10.x line governs
*reason-string strictness* (only §7-base failures carry the normative
reason-string contract), not whether to gate — a pinned vector is a pinned
vector. The audit-path fold (023/017/026) is one reusable function, not three.

**Decision D-W2-2:** The audit-path verifier (RFC 6962 inclusion proof) is new
verifier capability, not just a corpus check. It lands in
`verify/inclusion.go` as `VerifyAuditPath(leafHash, path, root)` so it is
available to the verifier proper (a future §10.31 partial-disclosure CLI mode),
not buried in the test runner. CUPID-composable: the same fold verifies 023,
017, and 026's concatenated paths unchanged — 026's spec note confirms a §10.31
verifier verifies the concatenated path with no special-casing.

## Acceptance

- `verify/backfill.go`: `RecomputeBackfillMerkleRoot(manifest, metadataLeaf)` +
  `CheckBaselineManifestSHA256`. 035 sign_payload gates green (no longer
  deferred); the gate's deferred-count drops by one.
- `verify/inclusion.go`: `VerifyAuditPath` — round-trips every 023/017 leaf to
  the pinned root; 026's concatenated paths verify unchanged.
- Rich-family runner asserts 003, 016, 022, 023, 024, 026 expected.json shapes.
- `go build` + `go vet` + `go test` green per module after each step.
- No DRY violation: backfill reuses `core/jcs` + `verify/merkle`; the audit-path
  fold is one function shared across 017/023/026.

## Non-goals

- Materializing 034-successor-attestation's paired backfill (spec-side; the 035
  step-5 companion cross-ref arms when 034 lands a backfill).
- A real Ed25519 signature on 035 (the vector pins the byte form, not a sig;
  case 018 covers the v1.0b signature path).
- Negative N001–N038 fixtures (still Heather's lane; gate auto-arms).
- Rust (no benchmark shows Go is the bottleneck; all hot work is SHA-256).

## Risk / fork log

- **D-W2-1**: build all 7 rich members, don't defer §10.x ones. Reversible.
- **D-W2-2**: audit-path fold lands as verifier capability, not test-only.
- **016 leaf semantics**: 016's `leaves_hex` are leaf *preimages* (root applies
  `H(0x00||preimage)`); 003's leaves are `payload_hash` values (same wrap). Both
  feed `MerkleLeafHash` then `MerkleTreeHash`. If a root diverges, the first
  place to look is whether the runner double-wrapped or skipped the leaf hash.
- **026 top-tree**: the outer tree combines subtree roots DIRECTLY via
  `MerkleNodeHash` (no leaf-hash wrap on subtree roots) — a divergence here is
  the classic "wrapped the subtree root as a leaf" bug. Pinned in the spec note.

## Realized outcomes (wave 2, 2026-06-10)

Both items landed, build/vet/test green per module after each step. Two local
commits on `main` (NOT pushed — Steve reviews):

| Commit | What |
|---|---|
| `d287d94` | §10.42 backfill-Merkle recompute (`verify/backfill.go`) + 035 gate (`vectors/backfill.go`). 035's merkle_root recomputed from the 9-leaf manifest+metadata tree and gated against `expected_merkle_root_hex.txt`; §10.42 step-3 manifest-array SHA-256 binding checked. |
| `b967ca4` | Rich-family gate (003/016/022/023/024/026) + two reusable verifier capabilities: `verify/inclusion.go` (RFC 6962 audit-path fold) + `verify/streaming.go` (§10.29 `StreamState`) + `verify/fingerprint.go` (shared §4.1 helper). |

**Recompute path — §7 steps now executing for backfill:**
- §10.42 step 1 (Merkle root recompute over baseline manifest + metadata leaf) —
  035 root `8943b16e…` reproduced byte-for-byte vs. the Python reference.
- §10.42 step 3 (manifest-array SHA-256 binding `880f875…`) — gated.
- The 035 metadata leaf JCS form is reproduced at exactly 728 bytes (pinned).

**Rich-family — §7-scope coverage now executing:**
- §7 step 10 daily Merkle (003, over 6 payload_hashes in (run_id,seq)).
- §7 step 10 RFC 6962 odd-leaf right-promote (016, 3/5/7 leaves).
- §4.1 session-key + key-fingerprint recompute (003).
- §10.31 inclusion-proof audit-path fold (023) — new `VerifyAuditPath`.
- §10.32 per-device HKDF (024) — 3 distinct device-bound keys.
- §10.37 hierarchical aggregation (026) — top-root-from-subtree-roots +
  10 concatenated-path folds, all through the same §10.31 fold unchanged.
- §10.29 streaming state machine (022) — 6 scenarios + 13-row transition table.

**Final corpus counts (executed vs SKIP), fresh `-count=1` run:**
- Canonical-output: **59/59 checks across 25 vectors** (grew from 22 — Heather
  landed 034/035 stubs + the 037-053 family during the run).
- Sign_payload: **6/6 checks across 3 vectors**, 1 deferred (035 full byte-
  compare — fixture gap, see below).
- Backfill (§10.42): **2/2 checks across 035** (NEW this wave).
- Rich-family: **31/31 checks across 6 vectors** (NEW this wave).
- Negative: **38 asserted, 0 SKIP across 38 INDEX rows** (Heather materialized
  ALL N001–N038 during the run — flipped from 38 SKIP in wave 1).
- JCS 008 corpus, HKDF RFC 5869, master fixture: green (pre-existing).

**Verifier capability added (not just corpus checks):**
- `verify.VerifyAuditPath` / `FoldAuditPath` — RFC 6962 inclusion-proof verifier,
  reusable by a future §10.31 partial-disclosure CLI mode.
- `verify.StreamState` (Step + Finalize) — §10.29 streaming-verifier state
  machine, reusable by a streaming CLI mode. Exit-code semantics map 1:1 to the
  existing `exitcodes` 4/5/6.
- `verify.KeyFingerprint` / `KeyFingerprintHex` — §4.1 fingerprint, shared.
- `verify.RecomputeBackfillMerkleRoot` / `CheckBaselineManifestSHA256` — §10.42
  backfill dispatch steps 1 + 3.

**Findings surfaced (Heather / spec lane — NOT patched by me):**

1. **035 input.json fixture gap (sign_payload reconstruction).** 035's
   `sign_payload_inputs` block omits `hkdf_inputs_digest_hex`; the `_compute.py`
   hardcodes the placeholder `0a1b2c3d…` (line 187) but never writes it back to
   `input.json`. So a clean-room verifier cannot reconstruct 035's full v1.0b
   sign_payload byte form (line 8) from its input alone. The §10.42 *merkle_root*
   IS recomputable and is gated; only the full byte-compare is blocked. Fix is
   one field in 035's `input.json` (add `hkdf_inputs_digest_hex` to
   `sign_payload_inputs`). This is a fixture issue, not a verifier bug.

2. **017 is a placeholder vector.** Its `leaf_hash_hex` / `sibling_hash_hex` /
   `expected_root_hex` are `PLACEHOLDER:…` strings, not real hex; the
   `expected-result.json` `_about` confirms "Hash and signature placeholder
   values pass through; the structural shape is normative." 017 therefore cannot
   be gated as a crypto fold — its normative content is the partial-disclosure
   *output shape* (P1–P5 step_outcomes + banner + completeness_assertion), which
   needs a §10.31 partial-disclosure output-writer (a CLI mode), well beyond §7
   base + byte-recompute. **Deferred with reasoning** — not a verifier gap.

## Wave-3 candidates (flagged for Steve — NOT in this wave)

- **Deep negative §7-walk driver.** Heather materialized all 38 negatives with a
  full `audit_file` (header + entries + seal) in `input.json` + a precise
  `expected_output.txt` (Status/Step/Reason/ExitCode). The current negative gate
  asserts only *fixture self-consistency* (`expected_output.txt` contains the
  INDEX-pinned reason). It does NOT yet run the §7 walk over each negative's
  `input.json` and confirm the verifier itself emits that Status/Step/Reason.
  Building that driver is the "full §7 step 1-12 over a real chain file" the
  wave-1 doc deferred. It needs a §7 walk that consumes the `audit_file` shape
  (distinct from the existing `Ledger` shape) and emits the normative output
  triple. **This is the single highest-leverage wave-3 item** — it turns 38
  negatives from self-consistent fixtures into live verifier-behavior assertions,
  including N025 which would exercise this wave's new §10.42 backfill recompute.
- **035 full sign_payload byte-compare** — auto-arms once finding #1's fixture
  field lands.
- **017 partial-disclosure output-writer** — the §10.31 CLI mode finding #2 needs.

---

# Wave 3 — deep negative §7-walk driver + v1.0c runner field + 035 byte-compare

**Date:** 2026-06-10 (same day, third wave)
**Status:** In progress
**Spec:** `chain-of-custody-DRAFT-0.3.0.md` (PRD-3) — §7 verification procedure + §10.12 exit codes
**Approved:** Steve, wave-3 go-ahead.

## What

Two items, in order:

1. **Thread `operational_events_log_root` through the vectors runner.** The
   `buildSeal` helper at `verifier/internal/vectors/signpayload.go` hardcoded
   `OperationalEventsLogRoot: ""`, which blocked materializing the v1.0c
   sibling-log sign_payload vector (027). Resolve the field from the fixture's
   seal block so the v1.0c form reconstructs from input alone.

2. **The deep negative §7-walk driver.** Today `runMaterializedNegative` asserts
   only fixture self-consistency (the `expected_output.txt` contains the
   INDEX-pinned reason-template). Build a driver that runs the verifier's actual
   §7 step 1–12 walk over each negative's `input.json` tampered `audit_file` and
   asserts the verifier ITSELF emits the pinned Status / failing Step / Reason /
   §10.12 ExitCode. This turns the self-consistent fixtures into live
   verifier-behavior assertions.

## Item 1 — runner field-resolution convention (committed first)

**Field name:** `operational_events_log_root_hex` on the seal block, consistent
with the existing `merkle_root_hex` / `hkdf_inputs_digest_hex` `_hex`-suffix
convention in the 018/019/020 `byteFormInput` layout. The `core/signpayload.Build`
v1.0c path already appends `OperationalEventsLogRoot` as the 13th field; this is
runner field-resolution only — `buildSeal` reads `seal.operational_events_log_root_hex`
and the `byteFormInput`/`signPayloadInputsLayout` structs gain the field.

**Heather authors 027 against this convention.** A v1.0c sibling-log vector's
`input.json` seal block carries `operational_events_log_root_hex` alongside
`merkle_root_hex`; the sign_payload runner reconstructs the 13-line form and
byte-compares against `expected_sign_payload.txt`.

## Item 2 — the §7-walk driver

### The audit_file shape (distinct from the Ledger shape — confirmed)

The negative corpus `input.json` carries a top-level `audit_file` object:
`{header, entries[], seal}`. The header is `{format_version, tenant_id,
seal_date, genesis_hash_hex, hkdf_inputs_digest_hex}`; each entry is `{seq,
tenant_id, run_id, key_version, key_fingerprint_hex, format_version, event,
event_canonical_hex, prev_hash_hex, payload_hash_hex}`; the seal mirrors the
sign_payload family. This is the §4.1 construction from the corpus's shared
`negative/_lib.py` — session_key = `HKDF(IKM, salt=ffiec…salt,
info=ffiec…info|tenant)`, MAC = `HMAC(session_key, prev_hash || canonical)`,
fingerprint = `SHA-256(tenant||ikm)[:16]`. It is NOT the older
`verify.Ledger`/`ChainEntry` shape (`EventPayloadJCS` + per-entry
`EntryID`-bound MAC). The driver implements a fresh §7 walk over the `audit_file`
shape, reusing the genuine shared primitives: `core/hkdf.Derive`,
`core/constants.*`, `verify.MerkleLeafHash`/`MerkleTreeHash`,
`verify.KeyFingerprint`.

### The IKM registry

The walk needs IKMs keyed by `key_version`. The corpus pins them in
`chain_vectors.json inputs` (`ikm_v1_hex` / `ikm_v2_hex`) — the same two
generations the master fixture and `_lib.py` use. The driver loads them into a
`{1: IKM_v1, 2: IKM_v2}` registry; an entry whose `key_version` is absent (N007
uses 99) fails at §7 step 7.

### §7 step → reason mapping (the driver implements)

| Step | Check | Reason on failure |
|---|---|---|
| pre-flight | mid-write truncation (audit_file carries `_ndjson_truncated` markers) | `audit file ends mid-line — possible mid-write crash` (exit 2) |
| 1 | `header.format_version == "v1"` | `format_version <X> not supported by this verifier (running v1)` |
| 2 | recompute `hkdf_inputs_digest` | `header HKDF inputs do not match running v1 inputs` |
| 3 | `genesis_hash == 32 zero bytes` | `header genesis_hash does not match v1 constant` |
| 4 | `event.tenant_id == header.tenant_id` | `cross-chain lift detected at seq <N> (event.tenant_id mismatch)` |
| 5 | `entry.format_version == header.format_version` | `format_version mismatch at seq <N>` |
| 6 | `entry.prev_hash == expected_prev` (prev-hash first), then `entry.seq == expected_seq` | `chain link broken at seq <N>` |
| 7 | IKM lookup by `key_version` | `unknown key_version: no IKM for (tenant=<T>, key_version=<V>) at seq <N>` |
| 8 | `SHA-256(tenant‖ikm)[:16] == entry.key_fingerprint` | `key_fingerprint mismatch at seq <N>: looked-up IKM does not match the entry's recorded fingerprint` |
| 9 | `HMAC(session_key, expected_prev ‖ canonical) == entry.payload_hash` (uses **expected_prev**, not entry.prev — §7 step 9 footgun note) | `payload_hash MAC mismatch at seq <N>` |
| 10 | RFC 6962 root over payload_hash leaves `== seal.merkle_root` | `merkle root mismatch — ledger contents do not produce sealed root` |

Steps 11 (signature) and 12 (cadence/dev-mode) are NOT crypto-walkable from the
corpus: the baseline seal carries a placeholder signature
(`TEST-SIGNATURE-PLACEHOLDER…`) and no real Ed25519 public key. The walk
therefore stops at step 10; signature-tampering negatives are classified
contract-only (below). On a clean walk through step 10 the driver returns PASS.

### Reason-match contract

Per `negative/INDEX.md` line 7: "The constant prefix and the message family are
normative"; position-dependent tokens (`<N>`, `<V>`, `<T>`) are substituted from
the input. The driver's walk renders the reason with the **constant prefix
verbatim** and the position token substituted (e.g., `chain link broken at seq
2`). The gate asserts:
- `Status` exact (FAIL / PASS).
- `Step` exact (the bare step number, or `pre-flight`).
- `Reason` family-prefix match against the pinned `Reason-Template` with
  `<N>`/`<V>`/`<T>` token-substitution — the live walk's rendered reason must
  equal the fixture's rendered `Reason` line (constant prefix + substituted
  token), which is the strongest honest "the verifier itself emits this"
  assertion.
- `ExitCode` exact (§10.12: FAIL→1, truncation/structural→2).

A clean-room reference walk (Python, run during planning) confirmed all 18
live-walk candidates produce the pinned Status + Step + rendered Reason from real
recomputation — not from reading the expected string. That re-derivation guard is
the honesty proof: a fixture whose tamper label disagrees with its bytes (N036,
see below) is caught, not rubber-stamped.

## Classification — live-walk vs contract-only vs v1.x (the honest coverage map)

Each negative assessed against: does the §7 1–10 base walk over the materialized
`audit_file` genuinely produce the pinned Status/Step/Reason from recomputation?

| Vector | Class | §7 / §10.x reason |
|---|---|---|
| N001 payload bit-flip | **live-walk** | step 9 MAC recompute rejects (verified `!=` stored) |
| N002 events reordered | **live-walk** | step 6 chain-link (prev-hash mismatch at expected_seq 2) |
| N003 merkle altered | **live-walk** | step 10 root recompute rejects |
| N006 fingerprint flipped | **live-walk** | step 8 fingerprint recompute rejects (no MAC compute) |
| N007 unknown key_version | **live-walk** | step 7 IKM lookup miss (key_version 99) |
| N008 entry format mismatch | **live-walk** | step 5 per-entry format_version |
| N009 header format v2 | **live-walk** | step 1 format_version |
| N010 header hkdf flipped | **live-walk** | step 2 hkdf_inputs_digest recompute |
| N011 genesis nonzero | **live-walk** | step 3 genesis constant |
| N012 cross-chain tenant | **live-walk** | step 4 per-entry binding |
| N013 mid-write truncation | **live-walk** | pre-flight (`_ndjson_truncated` markers) → exit 2 |
| N014 botched rotation | **live-walk** | step 8 fingerprint (rotation defence) |
| N015 prev_hash substituted | **live-walk** | step 6 chain-link |
| N016 prev+payload recomputed | **live-walk** | step 6 chain-link (deep) |
| N030 output-hash mismatch | **live-walk** | step 9 MAC (verified `!=` stored) |
| N033 DP noise seed tampered | **live-walk** | step 9 MAC (verified `!=` stored) |
| N020 algorithm/key-type mismatch | **live-walk (structural)** | step 11 *structural field-compare*: `seal.algorithm` ("ed25519") vs the fixture's `resolved_public_key_type` ("rsa-3072") → `algorithm/key-type mismatch at signature verification`. No crypto — a string compare on two fixture fields the walk reads. |
| N022 format v1.1 | **live-walk** | step 1 format_version |
| N023 format "V1" case-variant | **live-walk (prefix)** | step 1; the fixture quotes `"V1"` but N009/N022 leave the value unquoted — quoting is a per-fixture author choice diverging from the spec step-1 template (`format_version <X> …`). The walk asserts Status+Step+family-prefix; the quoting divergence is flagged to Heather (finding below). |
| N004 signature garbage | **contract-only** | step 11 signature crypto — the garbage sig (`ZZZZ…`) does not decode to a 64-byte Ed25519 signature and there is no public key in the corpus; the genuine crypto rejection cannot run from the fixture alone. The baseline placeholder makes step 11 non-walkable for every fixture. |
| N005 signature wrong tenant | **contract-only** | step 11 signature crypto — needs a real signature over the wrong-tenant sign_payload; the fixture carries only the baseline placeholder + a `signed_for_tenant_id` descriptor. |
| N017 dual-algo partial coverage | **contract-only (v1.x)** | §4.3.2 case (b) PASS-WITH-ANOMALY; PQ posture descriptor, v1.x-deferred |
| N018 dual-algo not in posture | **contract-only (v1.x)** | §4.3.2 case (c); PQ posture descriptor |
| N019 dual-algo one-valid-one-invalid | **contract-only (v1.x)** | §4.3.2 case (e); PQ posture descriptor |
| N021 routing event tampered | **contract-only (v1.x)** | §4.4.1 routing-in-canonical-bytes; v1.x-disposition per INDEX |
| N024 acquirer-HSM sig mismatch | **contract-only** | §10.24 succession; `entries=0`, bespoke `successor_envelope` block, no walkable base chain |
| N025 backfill merkle corrupted | **live-walk (§10.42)** | exercises wave-2 `verify.RecomputeBackfillMerkleRoot`; the corrupted backfill root recompute rejects → `backfill merkle root mismatch at backfill seq <N>` |
| N026 additional_verifications invalid | **contract-only** | §10.12 strict-mode verifier-output validation (exit 3); `entries=0`, `verdict_object` block — not a chain-integrity check |
| N027 state-machine illegal transition | **contract-only** | §10.43 claim state-machine; `entries=0`, `transitions_table` + `walk` blocks |
| N028 adjuster anchor missing link | **contract-only** | §10.45; `entries=0`, `cedent_anchor`/`reinsurer_anchor` blocks |
| N029 bordereau out of order | **contract-only** | §10.46; `entries=0`, `lifecycle` block |
| N031 retrieval merkle tampered | **contract-only** | §10.49 retrieval-set Merkle; `entries=0`, `retrieval_set` block (a §10.x Merkle distinct from the base chain Merkle) |
| N032 HITL signature bad | **contract-only** | §10.50; `entries=5` BUT the base §7 walk PASSes — the tamper is to a HITL reviewer signature the base walk doesn't check; needs a §10.50 HITL-signature verifier path not built |
| N034 decadal reseal anchor mismatch | **contract-only** | §10.54; `entries=0`, `decadal_reseal` block |
| N035 challenge-response out of order | **contract-only** | §10.55; `entries=0`, `sequence` block |
| N036 OTLP/JSON bytes encoding | **contract-only** | §4.4 receiver-decoder layer — **the MAC does NOT break in the fixture** (recompute `==` stored, verified); the pinned `payload_hash MAC mismatch` describes a receiver-side OTLP/JSON encoding refusal pending separate dispatch (per INDEX + Glenn 2026-05-21), not something the base walk produces from these bytes. Don't fake it. |
| N037 leap-second captured_at | **contract-only** | §10.4 — the base §7 walk over the clean chain genuinely PASSes, but the pinned output is PASS **with a clock-skew anomaly line**; the anomaly requires a §10.4 captured_at-monotonicity capability not built in the verifier. The Status:PASS/exit-0 half is live; the anomaly line is the deferred half — asserted contract-only until §10.4 lands. |
| N038 discovery production form | **contract-only (v1.x)** | §10.13.1 evidentiary-artifacts; `entries=0`, `production_manifest` block |

**Counts:** live-walk **20** (N001-N003, N006-N016, N020, N022, N023, N025,
N030, N033) · contract-only **18** (N004, N005, N017-N019, N021, N024,
N026-N029, N031, N032, N034-N038) · of the contract-only set, 7 carry the
v1.x-disposition (N017-N019, N021, N026, N036, N038). N023 and N020 are live
with a documented mechanism caveat (prefix-match / structural-compare).

## Findings surfaced (Heather / spec lane — NOT patched by me)

1. **N023 format_version quoting divergence.** N023's pinned reason quotes the
   value (`format_version "V1" not supported…`) but N009 (`v2`) and N022
   (`v1.1`) leave it unquoted. The spec step-1 template (§7, line 1635) is
   `format_version <X> not supported by this verifier (running v1)` — no quoting
   prescribed. Either all three quote `<X>` or none do. The verifier renders one
   canonical form; the fixture inconsistency means one of the three won't match
   byte-exact. Recommend the spec/_lib pick a single rule and re-render all three.
   The walk asserts N023 on the family-prefix + Status + Step until resolved.

2. **N036 tamper-label vs bytes disagreement.** N036's `tamper.class` is
   "OTLP/JSON bytes encoding" and the pinned reason is `payload_hash MAC mismatch
   at seq 2`, but the fixture's `event_canonical_hex` was NOT mutated to break the
   MAC — a live recompute yields `==` stored. The pinned reason describes a
   receiver-decoder-layer refusal, not a base-§7 outcome from these bytes. This is
   consistent with the INDEX note ("receiver-side OTLP/JSON bytes-encoding refusal
   pending separate dispatch") — flagged so it's explicit that N036's live
   assertion is deferred to the receiver-decoder wave, not the §7 walk.

## Acceptance

- `buildSeal` resolves `operational_events_log_root_hex` from the fixture; the
  v1.0c form reconstructs from input alone (item 1, committed first).
- A §7-walk driver over the `audit_file` shape asserts Status/Step/Reason/ExitCode
  for the 20 live-walk negatives; the 18 contract-only retain the
  reason-template self-consistency assertion with the per-vector classification
  reason recorded.
- 035 full sign_payload byte-compare completes IF Heather's `hkdf_inputs_digest_hex`
  fix lands in 035's `input.json` during the run.
- `go build` + `go vet` + `go test` green per module after each step.
- No DRY violation: the walk reuses `core/hkdf`, `core/constants`,
  `verify.MerkleLeafHash`/`MerkleTreeHash`, `verify.KeyFingerprint`,
  `verify.RecomputeBackfillMerkleRoot` (N025) — no new crypto primitives.

## Non-goals

- Step 11 signature crypto over the corpus (no real Ed25519 key materialized).
- §10.50 HITL / §10.43 state-machine / §10.4 leap-second / receiver-decoder
  verifier paths (the contract-only extension reasons need capabilities not built;
  each is recorded with its spec-section reason).
- Rust (no benchmark; all hot work is SHA-256).

## Risk / fork log

- **D-W3-1:** the §7 walk stops at step 10 (no crypto-walkable signature). A
  future real-key fixture would extend it to step 11. Reversible.
- **D-W3-2:** reason-match is family-prefix + token-substitution, not raw
  byte-equality, because the corpus's own `Reason-Template` carries `<N>` tokens.
  The live walk's rendered reason equals the fixture's rendered `Reason` line for
  the 20 live vectors (the strongest honest match); N023's quoting is the one
  documented exception.
- **N025 backfill leaf semantics:** N025 reuses wave-2's §10.42 recompute; a root
  divergence there is the same "wrapped the subtree root as a leaf" class wave-2
  pinned.

## Realized outcomes (wave 3, 2026-06-10)

All three items landed, build/vet/test green per module after each step. Three
local commits on `main` (NOT pushed — Steve reviews):

| Commit | What |
|---|---|
| `3a70ebd` | Item 1: `buildSeal` resolves `operational_events_log_root_hex` from the seal block (v1.0c runner field). Heather authors 027 against this convention (committed first so her concurrent work reads it). |
| `a131647` | Item 2: `verify.WalkAuditFile` (§7 step 1-10 over the audit_file shape) + the negative-walk driver. 20 vectors assert live verifier behavior; `TestNegativeLiveVectors_RunLivePath` proves the live path runs (not the contract-only fallback). |
| `550aa72` | 035 full sign_payload byte-compare: `recoverBackfillRoot` feeds the §10.42-derived root into the v1.0b reconstruction (Heather's 035 fix `0f6825a` unblocked it). |

**Reason-match contract — the two normative rules implemented (D-W3-2 superseded):**
The match is NOT raw byte-equality and NOT a loose prefix. It applies both §7
normative reason rules:
1. **§7-line-1572 normative-prefix rule.** The verifier MAY append `: detail`
   after the pinned normative prefix; a pin that is the bare prefix (N014:
   `key_fingerprint mismatch at seq 4`) is satisfied by the verifier appending
   `: looked-up IKM…`, and a full-detail pin (N006) is satisfied verbatim. The
   match accepts `got == pin` OR `got` starts with `pin + ": "`.
2. **Quote-insensitivity (N023).** The format_version value quoting diverges
   between fixtures (N023 quotes `"V1"`; N009/N022 + the spec template leave it
   unquoted). Comparison strips ASCII double-quotes from both sides — exactly
   that inconsistency, nothing else.
This is stronger and more correct than the wave-3-plan's family-prefix framing:
Status + Step + ExitCode are asserted exactly; only the reason carries the two
documented normative tolerances.

**035 deferral resolved.** Heather's `0f6825a` surfaced `hkdf_inputs_digest_hex`
into 035's `input.json`. The remaining gap was the DERIVED merkle_root (the
§10.42 backfill root, never pinned in `sign_payload_inputs`); `recoverBackfillRoot`
recomputes it via the wave-2 path and feeds it into the byte-form reconstruction.
035's full v1.0b sign_payload now matches byte-for-byte including line 7
(`8943b16e…`). sign_payload gate: **8/8 across 4 vectors, 0 deferred** (was 6/6
across 3, 1 deferred).

**027 status.** Heather prepared 027's `_compute.py` (ffiec-public `6aef365`) and
recorded the convention as confirmed-in-WIP (`745b89e`), but 027's `input.json`
is NOT yet pinned. The runner field-resolution (item 1) is committed and waiting;
027 will gate the moment its `input.json` materializes with
`operational_events_log_root_hex` on the seal block — no Go change needed.

**Final full-gate counts (fresh `-count=1` run, both modules green):**
- Negative: **38 asserted (20 live-walk + 18 contract-only), 0 skipped.**
- Sign_payload: **8/8 across 4 vectors, 0 deferred** (035 now complete).
- Canonical-output: 59/59 across 25 vectors.
- Backfill (§10.42): 2/2 across 1 (035 root recompute).
- Rich-family: 31/31 across 6.
- Master fixture: 6/6. JCS 008 + HKDF RFC 5869: green.
- `core` + `verifier` modules: `go build` + `go vet` + `go test` all green.

**Verifier capability added (not just corpus checks):**
- `verify.WalkAuditFile(af, ikms) Outcome` — the §7 step 1-10 procedure over the
  audit_file shape. A future `verifier verify --audit-file` CLI mode drives it.
- `verify.CheckAlgorithmKeyType(seal)` — the structural §7-step-11
  algorithm/key-type field compare (no crypto).
- `verify.AuditFile` / `AuditHeader` / `AuditEntry` / `AuditSeal` / `Outcome` /
  `IKMRegistry` — the audit-file shape + the normative output triple.

**Live-walk honesty discipline (the fool's three corrections, all applied):**
1. Re-derivation guard, not label-trust: a clean-room Python §7 walk (planning)
   + `TestNegativeLiveVectors_RunLivePath` (committed) prove the 20 live vectors
   produce the pin from real recomputation. N036's tamper-label-vs-bytes
   disagreement was caught this way and kept contract-only.
2. N004/N005 split assessed: both stay contract-only (no real Ed25519 key in the
   corpus; the baseline placeholder makes step-11 crypto non-walkable for every
   fixture). N004's garbage sig is structurally malformed but the verifier has no
   key to verify any sig against, so the honest call is contract-only.
3. N020 stated explicitly as a *structural field-compare* (algorithm vs
   resolved_public_key_type), not crypto.

**Flagged for Steve / Heather (spec lane — NOT patched by me):**
- **N023 quoting divergence** (finding #1 above) — the verifier handles it via
  quote-insensitive match; recommend the spec/_lib pick one quoting rule.
- **N036 tamper-label vs bytes** (finding #2 above) — kept contract-only; the
  pinned `payload_hash MAC mismatch` is a receiver-decoder reason, not a §7-walk
  outcome from these bytes (the MAC does not break in the fixture).
- **N025 Step descriptor** is the spec's prose `§10.42 step 2 (Merkle root
  verification)`, not a bare number; the live recompute-rejection + the rendered
  reason (`backfill merkle root mismatch at backfill seq 1`) + exit code are the
  genuinely-computed assertions. The prose Step is matched as-pinned.

---

# Wave 4 — real signature crypto (§7 step 11) + repo hygiene + regression hardening

Steve's wave-4 directives: real Ed25519 signature verification, regression tests,
repo hygiene (`.gitattributes`), N023 tolerance tightening, and continue
finalizing. A real test keypair now exists at
`E:\dev\testing\private-keys\tesseraseal\` (public key
`0985603b6c0e099bac783bcd7801664ed376a7888eafbf191859854bf8ff7f35`, private seed
local-only, env-var override `TESSERASEAL_TEST_KEY_DIR`). Heather concurrently
materializes N004/N005 as real signature-bearing fixtures and publishes the public
key into the corpus.

## Realized outcomes (wave 4, 2026-06-10)

Five local commits on `main` (NOT pushed — Steve reviews). Build/vet/test green per
module after each step.

| Commit | What |
|---|---|
| `80f84e1` | `.gitattributes`: LF-pin Go source + binary-guard crypto material. |
| `6600d1a` | §7 step 11 — live Ed25519 seal-signature verification (verify package). |
| `2c7d68f` | Live §7 step-11 signature driver — N004 flips to live-walk (vectors package). |
| `c4b7e93` | N023 tolerance tightening — quoted `format_version` rendering, byte-exact compare. |
| `e35a029` | CRLF-class regression — mangled corpus read is caught, not tolerated. |

### Directive 1 — `.gitattributes` (repo hygiene)

Same bug class Heather guarded in ffiec-public (`18e7c32`: `spec/test-vectors/**
-text`): a fresh clone with `core.autocrlf=true` silently rewrites LF→CRLF on
checkout and changes the bytes a tool reads. Two guards:

1. `text eol=lf` on `*.go`/`go.mod`/`go.sum`/`go.work` + prose/config — gofmt emits
   LF, and a CRLF-rewritten source reads as a whole-file diff to a Unix contributor.
2. `-text` (binary) on `*.hex`/`*.pem`/`*.sig`/`*.bin`/`*.golden` — pins the
   crypto-material class **before** the first such file lands. The verifier consumes
   only public-key material; the test keypair lives outside the repo under
   `TESSERASEAL_TEST_KEY_DIR`.

`git add --renormalize .` produced **zero churn** — the tree is already all-LF, so
this is a pure forward guard with no risk of a normalization storm.

### Directive 2 — §7 step 11, live Ed25519 (the wave-3 contract-only call is now resolved)

The wave-3 doc kept N004/N005 contract-only because "no real Ed25519 key in the
corpus." Wave 4 has the key. Step 11 is now a real two-assertion crypto check:

1. **Reconstruct-and-compare.** Rebuild the §4.3 `sign_payload` from the seal's
   structured fields via `core/signpayload.Build` (the same builder vectors
   018/019/020/035 gate byte-for-byte) and confirm it equals the seal's published
   `sign_payload_hex`. A published payload that does not reconstruct from the real
   fields is itself a step-11 failure — it covers bytes that do not bind the day's
   real root/tenant/version. Running this **before** the verify catches a tampered
   `tenant_id` even when an attacker also forged a matching signature over their
   substituted payload (the N005 wrong-tenant mode).
2. **Verify.** `ed25519.Verify(pub, reconstructed, sig)`. The N004 garbage-signature
   mode fails here.

Both render `signature verification failed` / step 11 / exit 1, matching the pins.

**Trust boundary:** the verifier consumes **only** public-key material —
`ParsePublicKeyHex` is the single entry point; `CheckSealSignatureV1` takes an
`ed25519.PublicKey`. It never reads the private seed. Tests that must SIGN a fixture
read the seed via `TESSERASEAL_TEST_KEY_DIR` and **SKIP with a clear message** when
absent, so CI without the key still runs the rest of the gate green.

**Live-readiness gate (the honest flip).** N004/N005 classify as `classSealSignature`;
the DRIVER decides per-fixture whether to run live or fall back to contract-only,
gated on (a) a base64-decodable signature and (b) a resolvable public key (corpus
pub.hex → vector pub.hex → `TESSERASEAL_TEST_KEY_DIR`/pub.hex). Result:

- **N004 flips LIVE.** Its `ZZZ…` signature decodes to 66 bytes, the published
  `sign_payload_hex` reconstructs from the structured fields, and `ed25519.Verify`
  rejects the garbage sig → the pinned FAIL/11/exit-1. Check carries the live
  `/§7-step-11-signature` suffix.
- **N005 flipped LIVE during this session.** Heather landed N005's real wrong-tenant
  signature (a 64-byte Ed25519 sig over the OTHER-tenant payload) while wave 4 was in
  progress, and the live-readiness gate auto-flipped it — **with no Go change.** The
  seal's structured `tenant_id` is `tenant-ffiec-test-1`, but its published
  `sign_payload_hex` binds `tenant-ffiec-test-OTHER`. The reconstruct-and-compare
  rebuilds from `test-1`, gets bytes ≠ the published OTHER-binding hex, and rejects
  at step 11 → the pinned FAIL/11/exit-1. This is exactly the wrong-tenant defence:
  `tenant_id` is bound into `sign_payload`, so a seal signed for a different tenant
  fails. The auto-flip is the proof the design is right — the driver gate, not a
  hard-coded slot list, decides liveness.

The honesty-guard test (`TestNegativeLiveVectors_RunLivePath`) counts signature
liveness by **actual live-path execution**, not a hard-coded number, so the count
tracks Heather's materialization. CI without the key: N004/N005 both fall back to
contract-only (no resolvable pub.hex), gate stays green.

**Step-11 unit coverage** (`auditsign_test.go`, all pass with the key present):
valid-signature PASS, garbage-signature FAIL, wrong-tenant FAIL (proves the
reconstruct-and-compare catches a real sig over a substituted tenant),
tampered-published-hex FAIL, wrong-key FAIL, wrong-key-length FAIL, plus
`ParsePublicKeyHex` validation and `reconstructSignPayload`-matches-builder.

### Directive 3 — N023 tolerance tightening (the wave-3 flag is now closed)

Wave 3 flagged the `format_version` quoting divergence to Heather and handled it
with a quote-insensitive match. Heather **resolved** it by re-rendering
N009/N022/N023 uniformly to the **quoted** form (`"v2"` / `"v1.1"` / `"V1"`) — the
chosen rule, matching the §7 reason-string family (sign_payload_version / algorithm
/ canonical_encoding `"X"` not supported all quote the value). Three changes:

1. `checkFormatVersion` renders `%q` (quoted) instead of `%s`.
2. `reasonMatches` drops the quote-insensitivity tolerance (`stripQuotes` deleted).
   The compare is now byte-exact modulo the normative `: detail` prefix rule.
3. Two regression tests pin the rule so a future wrong-quoting fixture FAILS loudly:
   `TestCheckFormatVersion_QuotedRendering` (verify) and
   `TestReasonMatches_QuotingIsByteExact` (vectors — proves quoted-vs-unquoted now
   compares UNEQUAL, the tolerance is provably gone).

### Directive 4 — regression sweep + CRLF-class test

Swept the wave-1..3 findings for fixed-but-not-regression-pinned items:

| Finding | Regression status |
|---|---|
| delegation_chain JCS sort | **Already pinned** — `predicates_test.go` (unsorted-chain reject + canonical-bytes). |
| 035 backfill-root recover | **Already pinned** — `backfill_test.go` (verify + vectors packages). |
| v1.0c field threading | **Already pinned** — `signpayload_test.go` (v1_0c 13-line + v1.0b parity). |
| **CRLF class** | **Was the gap — now pinned** (`crlf_corpus_test.go`). |

The CRLF regression is the Go-side complement to the `.gitattributes` guard: even if
the guard is missing/bypassed, `checkCanonicalSelfConsistency` hashes the RAW bytes
of `expected_canonical.txt` against the pinned sha256, so a CRLF-mangled golden file
changes the hash and FAILS loudly. `TestReadPrimarySHA_CRLFTolerant` documents the
deliberate asymmetry: the sha-PIN read IS CRLF-tolerant (the hex token has no
embedded EOL), the BLOB hash is not (its content bytes change). Pinning both records
where CRLF matters and where it does not.

## Final full-gate counts (wave 4, fresh `-count=1`, both modules green)

- Negative: **38 asserted, 0 skipped.** Live split: **20 always-live (§7 base walk
  ×18 + §7-step-11 structural alg/key-type ×1 + §10.42 backfill ×1) + 2
  step-11-signature live (N004 garbage + N005 wrong-tenant); 16 contract-only.**
  (N005 flipped live mid-session as Heather's real signature landed — the gate
  auto-tracked it with no Go change.)
- Sign_payload: 8/8 across 4 vectors, 0 deferred.
- Canonical-output: 59/59 across 25 vectors.
- Backfill (§10.42): 2/2 across 1.
- Rich-family: 31/31 across 6.
- Master fixture: 6/6. JCS 008 + HKDF RFC 5869: green.
- New step-11 unit tests: all green with the test key present; SKIP (not FAIL)
  without it.
- `core` + `cliutil` + `ledger` + `testkit` + `verifier`: `go build` + `go vet` +
  `go test` all green.

## Verifier capability added (wave 4)

- `verify.WalkAuditFileWithKey(af, ikms, pub) Outcome` — the full §7 step 1-11 walk;
  step 11 runs when `pub != nil`. `WalkAuditFile` is now the no-key steps-1-10 path.
- `verify.CheckSealSignatureV1(seal, pub) Outcome` — the §4.3 step-11 reconstruct +
  verify.
- `verify.ParsePublicKeyHex(hex) (ed25519.PublicKey, error)` — the single public-key
  entry point (public material only).
- `verify.SealCarriesLiveSignature(seal) bool` — the positive-path live-readiness
  predicate.

## Remaining to product-final (honest list)

These are the known gaps between "the gate is green today" and "the verifier is a
shippable, byte-equivalent-to-the-references TesseraSeal artifact." None are
blocking the current wave; each is a tracked next step.

1. **Corpus `pub.hex` publication (Heather, in flight).** N004 + N005 both run live
   today by resolving the public key from the local-only `TESSERASEAL_TEST_KEY_DIR`.
   For an auditor on a fresh clone to drive the live signature path, the corpus must
   publish `test-signing-key.pub.hex` at the corpus root (the driver's first search
   candidate). When it lands, the live path is reproducible without the local key dir
   — no Go change; the resolver already prefers the corpus copy. (N005's real
   wrong-tenant signature already landed mid-wave; both signature vectors are live.)

2. **N036 receiver-decoder wave (contract-only today).** N036's pinned
   `payload_hash MAC mismatch` is a receiver-decoder reason (OTLP/JSON bytes-encoding
   path), not a §7-walk outcome — the MAC does not break in the fixture bytes. Needs
   the receiver-decoder verification path before it can go live. Flagged to Heather
   in wave 3; still open.

3. **017 §10.31 output-writer (the verdict-emission path).** The verifier computes
   Outcomes but does not yet emit the §10.31 structured verdict object / §10.12
   exit-code output on a real CLI surface. The `verifier verify --audit-file` mode
   that drives `WalkAuditFileWithKey` and prints the normative triple + verdict is
   the next CLI build. 017 (merkle-inclusion partial-disclosure) gates the
   partial-disclosure output shape.

4. **015 PQC dual-algorithm (deferred to v1.x).** The `signatures` list form
   (§4.3.2 Variant B, Dilithium/SLH-DSA coexistence) is not implemented. The
   single-algorithm Ed25519 path is complete; the dual-algorithm dispatch (verify
   each entry in canonical-sorted order against its `public_key_id`-resolved key)
   is a v1.x build. 015 is the conformance witness.

5. **027 v1.0c sibling-log vector (waiting on Heather's input.json).** The runner
   field-resolution for `operational_events_log_root_hex` is committed (wave 3);
   027 gates the moment its `input.json` materializes. No Go change needed.

6. **CLI front door + runbook (Phase 6-7 of the build plan).** The auditor-facing
   `tesseraseal-verifier` CLI (read chain + config, emit §10.12 exit code), godoc on
   every exported symbol, and the operator runbook are not yet built. The verifier
   core is the engine; the CLI surface is the deliverable an auditor runs.

7. **Byte-equivalence cross-check against .NET + Python references.** The Go verifier
   produces Outcomes that match the corpus pins, but a direct byte-for-byte
   diff of the verdict output against Richard's .NET reference and the Python
   reference for the same input is not yet wired. That is the multi-implementation
   conformance proof; it needs the §10.31 output-writer (item 3) first.
