# Current State
_as of 2026-09-22_

## PUBLIC AS OF 2026-09-22

**This repository is PUBLIC.** Steve's reasoning, and it is the project's own premise: the
README already promises "a chain of custody an examiner can independently walk, without
trusting the institution's vendor or the institution's ops team." **An examiner cannot walk
a chain with a closed-source verifier; they would only be trusting the vendor's word about
the vendor's own chain.** The verifier has to be inspectable or the custody claim is
self-report.

All five modules ship: `core/`, `verifier/`, `ledger/`, `cliutil/`, `testkit/`. Shipping
four of five was considered and rejected on 2026-09-22 because the README documents the
ledger in its module table and quick start, and `go.work` builds it.

**Pre-release gate, run 2026-09-22.** gitleaks reported six findings. All six were reviewed
one by one and accepted as false positives: four are key FINGERPRINTS, which are truncated
hashes over a PUBLIC key and are meant to be published, and the rest are synthetic values
in `docs/auditor-stories/`, which is fiction about a made-up institution. They are recorded
by fingerprint in `.gitleaksignore` with the reasoning, rather than rewritten out of
history, so a reader who runs gitleaks finds the same six and the explanation together.
The four carrying files moved to the local-only `_wip/` tree.

**Note for anyone repeating this:** moving those files did NOT clear the scan. gitleaks
reads history, not the working tree. Only the ignore file cleared it.

## What we're building right now
The Go reference implementation (three modules: `core/` primitives, `ledger/` ingest server, `verifier/` offline CLI) of the FFIEC AI chain-of-custody standard, plus the conformance-vector verifier work that tracks the spec as it evolves in the sibling `ffiec-public` repo. This repo's own `docs/` tree is deprecated (historical snapshot only) — the authoritative spec, test vectors, and governance now live in `ffiec-public`. Recent work has concentrated on the §7 verifier-walk logic: live Ed25519 signature verification, byte-exact reason-string rendering, and FINRA/SEC regulator-overlay drafting under `docs/regulator-pack/` (superseded status pending confirmation against ffiec-public — see Open questions).

## Active decisions
- README states this repo is implementation-only; `ffiec-public` is canonical for spec/vectors/governance — do not treat this repo's `docs/` as authoritative for anything spec-shaped.
- 2026-08-07: gitleaks pre-commit secret scanning added (Phase 1 rollout), matching the house-wide security push.
- Live crypto over mocked crypto: step-11 §7 signature verification moved from a stubbed/mocked check to a live Ed25519 walk against real key material (`2c7d68f`, `6600d1a`).
- Byte-exact comparison discipline for verifier output: N023 tolerance tightened to quoted `format_version` rendering with byte-exact compare rather than loose string matching (`c4b7e93`) — this mirrors a decision later formalized as a normative spec rule in `ffiec-public`'s CHANGELOG (§7 quoting convention).
- CRLF-class regression test added so a mangled corpus read is caught, not silently tolerated (`e35a029`).
- `.gitattributes` added to LF-pin Go source while binary-guarding crypto material (`80f84e1`) — a correctness/security guard, not a style choice.

## Open questions
- RESOLVED 2026-09-22: the gitleaks-hooks commit was pushed with the public release.
- `docs/regulator-pack/` in this repo overlaps in name/content with `ffiec-public/docs/regulator-pack/` (which the public repo's README lists as canonical). Verify this repo's copy is genuinely stale/historical and not accidentally diverging content that should be reconciled or deleted.
- Whether the wave-4 PRD work (crypto/CRLF hardening, N023 tightening) referenced in recent commit messages has a corresponding update needed in `ffiec-public`'s spec text, or whether it's already reflected there.

## Next action
Read the tip of `docs/design/finra-gap-analysis-2026-07-02.md` and the regulator-pack `finra-overlay.md`/`sec-overlay.md` to confirm which FINRA/SEC content is current vs. superseded by `ffiec-public`, then resolve the docs/regulator-pack duplication question above before doing further work in this repo's `docs/` tree.

## Stop condition
Halt and return to Steve before pushing the unpushed `gitleaks` commit to `origin/main`, and before any change that touches spec-shaped content (test-vector semantics, chain-of-custody field definitions) without first checking whether the same change needs to land in `ffiec-public` first (spec is canonical there).

## Needs approval
- **The repo is PUBLIC as of 2026-09-22.** Every push is a publication. Run gitleaks before any push, and treat the audience as regulators and examiners.
- Any edit to `docs/regulator-pack/` content — check against `ffiec-public` canonical copy first; don't let the two silently diverge.
- Architectural changes to `core/`, `ledger/`, or `verifier/` — Jared owns this Go/Rust lane.
