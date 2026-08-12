# Current State
_as of 2026-08-12_

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
- Local branch `main` is 1 commit ahead of `origin/main` (the gitleaks-hooks commit) — confirm this local-only commit is intentional to keep unpushed, or push it.
- `docs/regulator-pack/` in this repo overlaps in name/content with `ffiec-public/docs/regulator-pack/` (which the public repo's README lists as canonical). Verify this repo's copy is genuinely stale/historical and not accidentally diverging content that should be reconciled or deleted.
- Whether the wave-4 PRD work (crypto/CRLF hardening, N023 tightening) referenced in recent commit messages has a corresponding update needed in `ffiec-public`'s spec text, or whether it's already reflected there.

## Next action
Read the tip of `docs/design/finra-gap-analysis-2026-07-02.md` and the regulator-pack `finra-overlay.md`/`sec-overlay.md` to confirm which FINRA/SEC content is current vs. superseded by `ffiec-public`, then resolve the docs/regulator-pack duplication question above before doing further work in this repo's `docs/` tree.

## Stop condition
Halt and return to Steve before pushing the unpushed `gitleaks` commit to `origin/main`, and before any change that touches spec-shaped content (test-vector semantics, chain-of-custody field definitions) without first checking whether the same change needs to land in `ffiec-public` first (spec is canonical there).

## Needs approval
- Pushing local commits to `origin/main` — confirm intent first (this repo is PRIVATE per project memory; verify the intended audience of any push).
- Any edit to `docs/regulator-pack/` content — check against `ffiec-public` canonical copy first; don't let the two silently diverge.
- Architectural changes to `core/`, `ledger/`, or `verifier/` — Jared owns this Go/Rust lane.
