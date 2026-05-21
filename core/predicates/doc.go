// Package predicates implements the per-event and per-day verifier
// predicates that compose into the §7 procedure plus the §10.x
// bonus-verification dispatches.
//
// Spec authority: §7 normates the 12-step procedure (with §7 step
// 3a sub-numbered for tenant_id character class). §10.x sections
// contribute additional predicates that emit
// `additional_verifications` markers under PASS (e.g., §10.42
// backfill seal, §10.21 cross-vendor handover, §10.58 component
// identity binding-walk). The Kognitos-12pt wave (per Richard's
// implementation plan) brings P-49 / P-50 / P-51 in PRD-3.
//
// Status: stub. The predicate dispatcher lands in Commits 3-5 of
// the verifier upgrade, one tranche per §7 step group:
//
//   - Commit 3: predicates for steps 1-3a (header pre-flight).
//   - Commit 4: predicates for steps 4-9 (per-event walk).
//   - Commit 5: predicates for steps 10-12a (per-day + verdict).
//
// Each predicate is a named function with a stable spec-mapped
// step ID (a string, NOT an integer — §7's sub-numbered "3a" / "12a"
// are first-class step IDs per Richard's cross-consult 2026-05-21;
// the §7 + §10.12 normative text is silent on integer-vs-string but
// the casual "§7 step 3a" usage in negative-vector INDEX.md confirms
// strings are the discriminator), a byte-exact failure reason string
// from §7 / negative-vector INDEX.md, and an explicit witness-mode
// applicability flag per §7's witness-mode-applicability table.
//
// Step IDs are typed as `type StepID string` so callers can't pass
// arbitrary strings, and so a future spec amendment adding "5b" or
// "11.1" doesn't require a type change — only a new const. The Go
// side is the first implementation to emit `Step: <id>` on the CLI's
// FAIL line; the .NET + Python references carry step semantics
// inline but emit nothing, so the Go choice of string-form step IDs
// becomes the de facto reference until Glenn migrates them up.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/ChainVerifier.cs is the orchestrator
// shape; its inline step methods are the per-step predicates.
package predicates
