# Round 15 — Big Four SOC Audit Engagement Partner review

**Reviewer:** Khadija Ramirez, Engagement Partner, Big Four SOC practice (Madrid office). Twenty-two years on SOC engagements; current portfolio is European bank-tech, including several institutions that ship AI into customer-facing journeys under EBA and Banco de España oversight. I read this material the way I read any first-look package — the description in front of me, the procedures in front of me, and the question "if my team had to staff this on Monday, would the work-papers stand up at concurring-partner review and at peer review."

**Materials read for this pass:** `spec/chain-of-custody-v1.md` §7 step 12a; `audit-procedures.md` P-22..P-29; `soc-pack/control-evidence-events.md` (regulator-fingerprint and cold-DR events); `soc-pack/section-4-template.md`; `control-map/CUECs.md` (CUEC-CRY-06 + CUEC-IR-05); `control-map/TSC-mapping.md` (CC8.1); `anomaly-documentation-template.md` (severity table); `user-entity-summary.md`; `regulator-pack/sample-report.md`.

**Stopping criterion:** zero High, zero Medium remaining.

---

## Headline

I expected to come in with at least one Medium on this round. I did not find one. The package is at the level where my team would staff against it without needing to invent compensating procedures around gaps in the spec, the description template, or the procedure catalog. The three things I was watching for on this read — (a) does the new §7 step 12a actually reach the working paper through a procedure with a sample-size floor, (b) does the regulator-fingerprint reception path have the events, the CUEC, the procedure, and the description-template line, and (c) does the cold-DR-key dry-run consumption have the same four-part trail — all three landed.

The 0/0 result is below.

## What I tested against

For a SOC 2 Type II engagement on a chain-of-custody implementation, the package has to give me four artifacts that compose:

1. **A description template that names every control the opinion will attach to.** Section 4 has to be specific enough that the reader (a regulator, a downstream user entity, an external financial-statement auditor) can identify the control and the boundary.
2. **A procedure catalog that produces evidence for every control claim, with sample sizes that withstand concurring-partner challenge.** Underpowered sampling is the most common reason a SOC report gets returned to engagement at peer review.
3. **A control-map that ties each procedure into TSC criteria, so the opinion language is defensible.** A procedure that does not map to a criterion is procedure-without-purpose; a criterion that has no procedure is opinion-without-evidence.
4. **An anomaly-handling pattern that the reviewer can follow at three levels — what the verifier emits, what the institution does with it, what the SOC team tests.** Anomaly handling is where engagements get into trouble; if the institution and the SOC team disagree on what counts as "operationally explained," the engagement stalls.

I went through the package looking for breaks in any of those four artifacts.

## §7 step 12a — GenAI model identifier completeness

This is the new normative check that fires per-event when a chain entry carries the OTel `gen_ai.*` namespace prefix and either `gen_ai.request.model` or `gen_ai.response.model` is missing. The discriminator is the literal namespace `gen_ai.` — a chain entry with `tool.*` or `audit.*` only does not trigger the check. Under `--strict` it FAILs; under non-strict it is PASS-WITH-ANOMALY. Severity in the anomaly table is Medium and explicitly framed as control-completeness for SR 11-7 reproducibility, not chain-integrity.

What I was looking for: does the check actually reach my working paper, or does it stop at the verifier?

It reaches the working paper. The path is:

- **Verifier emits the anomaly.** Spec §7 step 12a is normative; an implementation that omits the check is non-conformant.
- **Anomaly-documentation-template carries the severity (Medium).** The institution's anomaly-evaluation record has to land on this row. The Medium framing is the right level — chain integrity is unaffected; what is affected is the institution's ability to reproduce the model call later, which is the SR 11-7 question.
- **Audit-procedure P-25 reaches it through stratified sampling.** The 3 × 5 × 3 = 45-entry minimum (model-version × decision-class × customer-impact-tier) gives my team enough surface to find a missing-identifier pattern if one exists. The fallback for institutions without a customer-impact-tier taxonomy — 3 × 7 = 21-entry minimum, with the absence noted as a control-program-maturity observation rather than a finding — is exactly right for the smaller / regional banks my Madrid practice serves. They typically do not yet have a customer-impact-tier framework; treating its absence as a maturity observation rather than a finding lets the SOC report ship without artificially inflating the institution's deficiency profile.
- **TSC-mapping puts it under CC7.2 (system monitoring) and CC8.1 (change management for SDK constants).** That is where the opinion language lives.

The defer-prohibition language in §7 step 12a — "implementations MUST NOT defer it to a second pass; deferral makes the check observable-at-scale (memory-overhead and latency for the deferred queue) and is non-conformant" — is the right call. A deferred-queue implementation would be observable in the institution's host metrics in a way that would force my team to ask "is your verifier the conformant one?" The inline requirement removes that question.

**No finding.**

## P-22..P-25 — anomaly evidence completeness

Reading these four procedures cold, the evidence pieces are specific enough that my staff seniors do not need a partner to interpret them.

- **P-22 (`key_fingerprint mismatch`)** — the four-piece evidence package (IKM-roster row identification + change-management approval + re-verification result + reconciliation cross-check) is what I would write if I were drafting this from scratch. The fourth piece — the cross-check against the most recent `master.reconciliation_completed` event with `fingerprint_unmatched_count=0` for the affected pair — is the one a junior auditor would miss; having it spelled out in the procedure forecloses that.
- **P-23 (`hkdf_inputs_digest mismatch`)** — the two acceptable explanations (known SDK / verifier constants change OR multi-region deployment-drift case) are exhaustive in the practical sense. Anything else routes to IR Scenario 1 as a format-construction defect, which is the right disposition. My team does not have to invent a third bucket.
- **P-24 (`audit_file_truncation_detected`)** — three recovery outcomes (`complete`, `partial`, `unrecoverable`) with explicit handling for each. The `unrecoverable` outcome triggering the 36-hour clock per IR Scenario 9 is the connective tissue I want to see; the SOC team is not making up its own incident-notification timing.
- **P-25 (`gen_ai_parameters` schema completeness)** — covered above under §7 step 12a; the procedure also confirms the institution's `gen_ai_parameters` schema covers the documented reproducibility surface (decoding parameters, sampler implementation, system prompt content or hash, retrieval context, intra-run dependencies). The disposition language — "schema gaps are not chain-integrity findings; they are MRM-program findings" — is the correct framing. The SOC team does not get pulled into MRM-substantive territory.

**No findings on P-22..P-25.**

## P-26 — model-inventory composition cross-check

This is the procedure that would have been my Medium last round if the sample-size floor were absent. It is not absent.

The per-stratum floor — "at least 10 chain entries per `(decision-class, customer-impact-tier)` pair, OR all chain entries in the stratum if the stratum's total is < 10 in the period" — gives my team statistical confidence in the per-stratum coverage assessment. The smaller-strata exhaustive-test fallback is the detail that protects the work-paper. A junior who pulls 5 entries out of a stratum that has 3 in it has tested zero entries; the floor language closes that.

The two failure modes ("models in chain but NOT in inventory" — High severity model-governance finding; "models in inventory but NOT in chain" — two root causes with different severities) are the correct decomposition. The "decommissioned but still listed" case is housekeeping; the "operating outside chain coverage" case is a chain-integration gap and is the higher-severity finding. The SOC team does not have to invent that bifurcation at engagement time.

**No finding.**

## P-27 — customer-dispute reproduction evidence completeness

The threshold-reasonableness review carve-out is the right shape. The SOC team confirms the institution's `decision_equivalent_threshold` exists, was committee-approved before the period, and is referenced consistently across the period's reproductions. That is process-of-governance, which is the SOC team's job. Substantive model-validation — whether the threshold is the right threshold — is the MRM committee's substantive work.

The threshold-reasonableness language itself — "a threshold that was set ad-hoc at dispute time, OR that varies across reproductions of the same decision-class without committee-approved cause, is a process-of-governance finding" — gives my team the test in one sentence. I can hand this to a senior and they can execute it without consulting me.

The recommended schema fields under `audit.reproduction.*` (`audit.reproduction.original_run_id`, `original_seq`, `iterations`, `decision_equivalent_count`, `decision_equivalent_threshold`, `outcome`) are concrete. My working paper carries those field names; my evidence pull is against those fields. The schema is RECOMMENDED rather than required, which is correct — the chain integrity-binds whatever the institution puts in, so making this MUST in the spec would be over-reach. As a SOC procedure, the recommendation is binding for the engagement; that is the right place for the obligation to live.

**No finding.**

## P-28 — regulator-fingerprint reception procedure operability

This procedure is new and it is the one I read first because regulator-fingerprint reception is the kind of failure that does not show up in the chain itself. If the institution gets a forged rotation notice and installs the forged fingerprint, the chain keeps verifying fine — but the regulator's view of the chain is wrong, which is the worse failure.

The procedure has the four pulls and the four checks I would expect:

1. Pull the three events (`regulator_fingerprint.rotation_received`, `.rotation_validated`, `.installed`) for the period.
2. Confirm any rotation has all three events with matching `correlation_id`. The correlation_id linkage is the detail; without it, the SOC team cannot prove the three events go together.
3. Confirm `overall_validation: PASS` for installed rotations.
4. For any FAIL: confirm IR Scenario 11 sub-variant was activated and the institution did NOT install the fingerprint.
5. Sample-test the change-management record cited in `regulator_fingerprint.installed.change_management_record_id`.

The CUEC-CRY-06 description in `CUECs.md` matches the procedure. The Section 4 template line ("Regulator-fingerprint reception procedure per `09-threat-model.md` §2.9 + audit-procedures P-28; emits `regulator_fingerprint.rotation_received` / `.rotation_validated` / `.installed` operational events; forged-notice failures route to IR Scenario 11 sub-variant") names the procedure, the events, and the failure-routing in one paragraph. The TSC-mapping CC8.1 row names the `regulator_fingerprint.installed` operational event as documenting trust-anchor-cache change events.

The four-part trail (procedure → CUEC → events → description-template line → TSC-map) is intact. My team does not have to reconstruct any of the four pieces at engagement time.

**No finding.**

## P-29 — cold-DR-key dry-run attestation consumption

Same shape as P-28; same four-part trail; same disposition.

- **Procedure (P-29):** pull the `cold_dr.dryrun_attestation_consumed` event(s); confirm at least one per project annual dry-run window in the period; confirm `consumption_outcome: PASS`; for PASS, confirm the archived attestation is retrievable; for FAIL or missing, confirm IR Scenario 11 was activated.
- **CUEC-IR-05:** institution consumes the project's annual `KEY-DR-DRYRUN-{year}.asc` cold-DR-key dry-run attestation within 30 days of publication; consumption recorded as the operational event; missed dry-run window triggers IR Scenario 11.
- **Operational event:** `cold_dr.dryrun_attestation_consumed` — the field schema (attestation_year, artifact_sha256, signing_identity_validated, consumption_outcome, consumed_by, archived_at) is what the procedure tests against. The `archived_at` field is the one that closes the loop on retrievability — without it the procedure could not test the "PASS events: confirm the archived attestation is retrievable" check.
- **Section 4 template line:** named explicitly with the path to the procedure.
- **TSC-mapping CC8.1:** `cold_dr.dryrun_attestation_consumed` operational event documenting institution's annual cold-DR-key dry-run consumption.

The "tenant_id: (institution-wide; not per-tenant)" note in the event schema is a small but correct detail. A multi-tenant institution operates one cold-DR fallback path against the project's release-management signature; tying the event to a specific tenant would mis-shape the procedure. The note tells my team to expect one event per institution per year, not one per tenant per year.

**No finding.**

## CC8.1 in the TSC-mapping

The CC8.1 row I was watching reads:

> Spec version change-control; `format_version` per chain entry (spec §4.4) and per file header bound under signed `sign_payload` (§4.3); **`master_key.retired` and `master_key.rotated` operational events documenting IKM-lifecycle change events**; institution's IKM-retirement procedure documented per spec §10.9; **`regulator_fingerprint.installed` operational event documenting trust-anchor-cache change events**; **`cold_dr.dryrun_attestation_consumed` operational event documenting institution's annual cold-DR-key dry-run consumption per `supply-chain.md`**

The line carries four classes of CC8.1 evidence (spec version control; IKM lifecycle; trust-anchor-cache change; cold-DR consumption) under one criterion. That is the right consolidation. CC8.1 is "change management procedures" and each of the four is a change-management event the SOC team tests against. A reviewer who pulls this line for the opinion-drafting workpaper has the four event names and can pull the four event streams in parallel.

**No finding.**

## Anomaly-documentation severity table

The `gen_ai_model_identifier_missing` row (Medium, control-completeness for SR 11-7 reproducibility, NOT chain-integrity) is the right severity assignment.

The framing matters: a Critical or High would over-call the failure (chain integrity is intact; the verifier is reporting a control-completeness gap, not a tampering signal). A Low would under-call it (the institution loses the ability to reproduce the model call, which is a load-bearing SR 11-7 capability for any AI-driven decision in a regulated journey). Medium is correct.

The "the institution's MRM program loses reproduction surface for the affected entries" framing tells the institution's risk function exactly what is at stake. That is the conversation the institution needs to have at its risk committee.

**No finding.**

## Section 4 description template

The procedures section of the template carries every operational procedure I expected, including the two new lines:

- "Regulator-fingerprint reception procedure per `09-threat-model.md` §2.9 + audit-procedures P-28; emits `regulator_fingerprint.rotation_received` / `.rotation_validated` / `.installed` operational events; forged-notice failures route to IR Scenario 11 sub-variant"
- "Annual cold-DR-key dry-run attestation consumption per `supply-chain.md` §"Institution-side consumption of cold-DR-key dry-run attestation" + audit-procedures P-29; emits `cold_dr.dryrun_attestation_consumed` operational event"

Both lines name the procedure, the events, and the failure-routing. The institution's drafter does not have to invent the language; they fill in `[institution-defined]` brackets where they exist and strike sections that don't apply.

The drafting notes ("the description must be true; the description must be complete in scope; the description must be clear") and the review checklist before issuance are at the level of detail I would expect on a SOC report description template. The check "All `[bracketed placeholders]` are filled in or stricken" is the one that catches the most common pre-issuance failure mode (a placeholder shipped in production text).

**No finding.**

## CUECs CRY-06 and IR-05

Both CUECs name the procedure they pair with (P-28 for CRY-06; P-29 for IR-05), name the operational events they generate, and name the failure-routing scenario (IR Scenario 11 / IR Scenario 11 sub-variant). The descriptions are at the level a user-entity reading the SOC report can use to confirm whether the institution operates the control.

CRY-06 lives under "Cryptographic controls"; IR-05 lives under "Cryptographic controls" as well. I notice IR-05 is named with the IR prefix but tabled under cryptographic controls. That is a small categorization choice; functionally the cold-DR fallback is a cryptographic-controls-and-IR-readiness control, so the table placement does not produce a working-paper problem. Mentioning it for completeness rather than as a finding.

**No finding.**

## Sample report — the SOC team appendix

The verifier-output-line → TSC-criterion mapping table at the end of `regulator-pack/sample-report.md` is the workpaper allocation table I want my engagement team to use. It has every step's failure record mapped to a TSC criterion, including step 7 (`unknown key_version`) → CC6.1 + ID.AM-08, step 8 (`key_fingerprint mismatch`) → CC6.1 + CC6.7, step 12 dev-mode-refused → CC6.8 + CC8.1.

The pass-line entries (`Spec §7 steps executed: 1..12` → PI1.1; `Merkle match: PASS` → PI1.2; `HMAC chain walk: PASS` → PI1.1) give my team the positive-evidence allocation that supports the unmodified opinion. The negative-evidence allocations (the failure-record rows) give my team the language for the modified-opinion case.

**No finding.**

## What I tested for and did not find

To be explicit about what I went looking for so the next reviewer knows my coverage:

- **Underpowered sampling.** Looked at P-25 stratification (3 × 5 × 3 = 45 with fallback), P-26 per-stratum floor (10 entries or all-if-less), generic sample-size table at the bottom of `audit-procedures.md` (25 / 50 / 75 by population). All three are at concurring-partner-defensible levels. Not a finding.
- **Procedure-without-evidence.** Looked for procedures that test against an event the events catalog does not actually define. P-28 → three regulator-fingerprint events present and field-schemas defined; P-29 → cold-DR event present and field-schema defined. Both close.
- **Evidence-without-procedure.** Looked for events in `control-evidence-events.md` that no procedure tests against. The two new event groups are tested by P-28 and P-29 respectively. Not a finding.
- **Severity-table holes.** Looked at the anomaly severity table for the new `gen_ai_model_identifier_missing` row. Present at Medium with the correct framing. Not a finding.
- **TSC-mapping holes.** Looked at the CC8.1 row to confirm the new events are listed. Both regulator-fingerprint and cold-DR are named explicitly. Not a finding.
- **Description-template holes.** Looked at the procedures section of the Section 4 template for the two new procedures. Both present. Not a finding.
- **CUEC-without-procedure.** Looked at CUEC-CRY-06 and CUEC-IR-05 for procedure references. Both name P-28 and P-29 respectively. Not a finding.
- **Customer-impact-tier-stratification absence handling.** Looked at how P-25 and P-26 treat institutions without a customer-impact-tier taxonomy. P-25 has the explicit fallback (21-entry minimum, control-program-maturity observation, NOT a finding). P-26 has the per-stratum exhaustive-test fallback (test all entries in strata < 10). Both close the smaller / regional bank case. Not a finding.

## Stopping criterion

**0 High, 0 Medium.**

This package, on this read, would clear my engagement-acceptance review. My team would staff the work against this catalog without needing to write supplementary procedures or compensating controls. The opinion would attach to a description that is specific, complete, and testable.

I would expect the next round to focus on what happens at scale (multi-region institutions; institutions with materially different decision-class distributions across regions; institutions where the customer-impact-tier framework varies across business lines). Those are operational-shape questions rather than spec or procedure-catalog questions, so they belong in the engagement's scoping conversation rather than as findings against the package as written.

— Khadija Ramirez
Engagement Partner, Big Four SOC practice (Madrid)
