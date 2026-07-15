---
status: plot-beats-stub (Story 12 decision committed 2026-05-11; full beats doc to follow)
audience: Steve (build out before drafting)
purpose: Story 12 seed — credit union marketing-AI vendor swap; foresight-cluster opener; plants Dawn's nagging question that pays off in Story 21
last updated: 2026-05-11
---

# Story 12 — Hill Country Federal Credit Union (working title)

## The decision (2026-05-11)

Story 12 is a **credit union marketing-AI vendor handover** under NCUA AIRES pressure. Steve's read on why it works: it broadens TesseraSeal's demonstrated scope (no credit union in Stories 01-11; NCUA unrepresented across the existing universe).

## Premise

**Hill Country Federal Credit Union** — large multi-state federally-insured FCU, ~$8B assets, Texas + neighboring states, NCUA-supervised. **TesseraSeal in production for 11 months** across the full member-experience surface.

Six months into a marketing-AI vendor handover: leaving a legacy banking-CRM platform (e.g., Total Expert), migrating to HubSpot Marketing Hub + an in-house ML scoring layer. **NCUA AIRES examination is 3 weeks out.** The CAE has been carrying anxiety about the vendor-handover boundary.

## TesseraSeal posture

- AWS-resident; Herald.Core.Aws runs the chain
- §10.21 (cross-vendor model-handover) shipped in a Herald release **7 months ago** — exactly when the CU started this transition. **Foresight pattern operating cleanly: the working group shipped the primitive before the institution needed it.**
- §10.69 (per-customer audit-trail subset disclosure) shipped earlier; member-disclosure surface is in production

## What this story establishes for the Stories 12-17 foresight cluster

- **Single substrate** (AWS only)
- **Single regulator dialog** (NCUA AIRES + CFPB §1033 cross-cut, but one supervisor in the room)
- **Single vendor handover** (one marketing platform → another)
- Cross-vendor anchor in its simplest form — hash old vendor's dump at cutover, anchor inside new chain, demonstrate byte-equality
- **Establishes the foresight pattern signature** that Stories 13-17 will continue

## Dawn's nagging question — the seed

During the reconciliation test, the byte-equality demonstration works cleanly. Dawn writes a note in the engagement file: *"works on one substrate. What happens when the substrate moves?"* She files it. The question sits unanswered for the rest of her audit-firm tenure.

**Two years later (Story 21), she gets to answer it from the vendor side.** Hill Country is named explicitly in Story 21's 4:30 PM scene as "a Story-12 era engagement."

## The signature audit moment

During reconciliation, Diana surfaces a record where a member received a marketing offer under the legacy vendor's ML scoring that the new vendor's ML scoring would never produce (the weights differ). Chen traces forward: the member subsequently took an auto loan; the loan rationale references the legacy marketing event. The CMO asks whether that's ECOA-defensible.

Dawn's answer: *"The decision is in the chain; the marketing event that informed it is in the chain; the rationale linkage is in the chain. Whether the policy is fair is a different question — but you can answer it now, with data."*

This is the first audit-side instance of marketing-data appearing as load-bearing chain content. It will become structural in Story 21.

## Regulatory framework

- **NCUA AIRES** (lead) — examination 3 weeks out; chain-of-custody artifacts will be presented in AIRES workpapers; §10.13 evidentiary artifacts compose with the AIRES workpaper model
- **CFPB §1033** (cross-cut) — Personal Financial Data Rights; member-disclosure packets must remain producible across the vendor handover
- **ECOA / Reg B** (cross-cut) — marketing-to-credit-decision linkage; the auto-loan case is the load-bearing example

## 4:30 PM client question

CU's CAE: *"If NCUA asks us in three weeks how the marketing AI changed across this transition, can we prove the change was tracked?"*

Dawn's answer: yes — the cross-vendor anchor proves exactly what changed and what didn't; the policy-vs-implementation distinction is preserved; the audit team can walk an examiner through the deltas live.

## Implementation prerequisites

- The legacy banking-CRM vendor name in the story (Total Expert proposed; could be Bonzo, BankBound, or another real banking-CRM product; pick before draft)
- The new vendor stack (HubSpot Marketing Hub + in-house ML proposed; comfort level with HubSpot as the named product)
- §10.21 must contain the field family the story exercises (`audit.vendor.handover.*` or similar — TBD; spec-hardening session's lane)

## Held-line principles (per Stories 14-17 doctrine)

- Foresight is operational, not editorial — the §10.21 shipping date is named matter-of-factly, never as "we anticipated this perfectly"
- The CU's CAE is calm and capable; the story doesn't position foresight as rescue
- Dawn's nagging-question note is filed quietly; she does not announce it at the time; the payoff in Story 21 is what gives it weight retroactively
- Steve appears only by reference (recusal protocol still active in the Story 12 timeline — he is the principal designer of the spec extension being exercised; Dawn discloses to Tom; Tom logs)

## Open questions for Steve

1. **Hill Country FCU as the name** — keep, or pick something different? (Texas-multi-state-FCU plausibility; Hill Country reads vaguely Hill Country.)
2. **Legacy vendor name** — Total Expert / Bonzo / something else?
3. **New vendor stack** — HubSpot Marketing Hub + in-house ML, or different combo?
4. **Where in the engagement does Dawn file the nagging-question note?** Suggested: at the end of the reconciliation test, in the engagement-file note section. (Could also be a drive-out monologue with Raj after Day 1.)
5. **Story 12's relationship to Stories 13-17.** Story 12 is the cluster opener; do Stories 13-17 build out as planned in `memory/story_beats_14_17.md`, or does Story 12 reshape any of them?
