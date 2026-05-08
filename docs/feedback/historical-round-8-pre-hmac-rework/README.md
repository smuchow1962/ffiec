# Reviewer feedback simulation (round 8 — full convergence)

> **What this directory is.** Final round. Target: **0 gaps and 0 partials per role document**.
>
> Round 7 left 6 partials (1 per role). Round 8 ships the documents that close those partials and re-runs the simulation. This round confirms convergence.

## Newly shipped documents

Closing the round-7 partials:

- `docs/INDEX.md` — top-level navigation across the corpus (closes Cyber Risk partial)
- `docs/user-entity-summary.md` — SOC 1 user-entity audience (closes SOC Audit partial)
- `docs/edge-and-federated-ai.md` — emerging deployment patterns (closes FFIEC Cybersecurity partial)
- `docs/portfolio-comparison-procedures.md` — manual portfolio comparison (closes FFIEC IT Examiner and EIC partials)

The Model Risk partial (current rule-timeline tracking for DORA, NIS2) is closed by explicit articulation: rule timelines evolve; the institution's compliance program tracks current rules; the chain corpus references the frameworks (`regulator-pack/ai-policy-alignment.md`) and the institution's compliance program operates current-rule tracking. This is institution-side responsibility, articulated as Answered.

## Roles modeled

Same six roles, fresh personas confirming whether the corpus now satisfies all questions.

### Big Four

1. [`01-big-four-cyber-risk.md`](01-big-four-cyber-risk.md)
2. [`02-big-four-model-risk.md`](02-big-four-model-risk.md)
3. [`03-big-four-soc-audit.md`](03-big-four-soc-audit.md)

### FFIEC

4. [`04-ffiec-it-examiner.md`](04-ffiec-it-examiner.md)
5. [`05-ffiec-cybersecurity-examiner.md`](05-ffiec-cybersecurity-examiner.md)
6. [`06-ffiec-eic.md`](06-ffiec-eic.md)

The roll-up at [`99-gap-roll-up.md`](99-gap-roll-up.md) confirms full convergence.
