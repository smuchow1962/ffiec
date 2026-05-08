---
status: informative
alignment-reference: FFIEC IT Examination Handbook (Booklet 12 — Operations; Booklet 2 — Information Security); 12 CFR Part 30 Appendix B; 12 CFR §53 (Computer-Security Incident Notification Rule); OCC Bulletin 2020-15; OCC Bulletin 2013-29 (three-lines-of-defense); OCC Bulletin 2023 TPRM Interagency Guidance; AT-C §105 (Attestation Standards); URSIT (Uniform Rating System for Information Technology); SR 11-7
companion-docs:
  - docs/audit-procedures.md (institution-side P-1..P-N procedures)
  - docs/regulator-pack/examination-response-workflow.md (PBC list and response workflow)
  - docs/regulator-pack/examiner-approval-template.md (examiner sign-off template)
  - docs/regulator-pack/examiner-training.md (examiner orientation to the chain)
  - docs/regulator-pack/finding-language.md (supervisory-letter phrasing)
  - docs/incident-response-playbook.md (Scenarios 3, 4, 9, 11)
  - spec/chain-of-custody-v1.md §1.4, §4.1, §4.3, §7, §10.1, §10.5, §10.7
date: 2026-05-07
version: 1.0.0
---

# FDIC / OCC Examination Overlay — Chain-of-Custody Control Testing

> **What this doc is.** Examination overlay for the FFIEC chain-of-custody control during FDIC, OCC, and Federal Reserve IT examinations. The overlay maps the chain spec and the institution-side audit procedures to the FFIEC examination workflow: pre-exam scoping, PBC (Provided By Client) request lists, control-test procedures, deficiency-severity thresholds, URSIT rating impact, supervisory-letter language, and follow-up examination workflow.
>
> **Audience.** Examiner-in-Charge (EIC), supporting examiners, the institution's IT Risk Officer and Compliance Office during examination preparation, and the FFIEC working group when standardizing examination guidance for AI-decision logging controls.
>
> **Scope discipline.** This overlay does not modify the spec or institution-side controls. It is examination-side guidance: what an examiner asks for, what they test, and how they cite findings. Institutions reading the overlay see the examination expectations they should prepare against.

## How to use this overlay

1. **Pre-exam scoping (60 days before exam open).** Send the institution the Pre-Exam Scoping Checklist (§1) so the institution prepares the PBC artifacts. Confirm posture, vendor relationships, and any in-period control changes during the scoping call.
2. **PBC delivery (5 business days before exam open).** Verify the institution produced every artifact in §2. A missing artifact at exam open is an early signal of a control-program-maturity issue.
3. **On-site testing (4-week exam, 20–60 hours of chain-specific testing depending on institution size).** Execute the procedures in §3 against the institution-side P-N audit procedures from `docs/audit-procedures.md`. Sample sizes follow §4.
4. **Deficiency disposition (during exam and at exam close).** Apply the severity matrix in §5 to each finding. Distinguish operational findings (workpaper-documented), MRA (Matters Requiring Attention, 90-day remediation), and MRIA (Matters Requiring Immediate Attention, may restrict further AI deployment).
5. **URSIT rating impact (post-exam).** Use the URSIT integration in §6 to determine whether the chain's operating effectiveness warrants an Operations rating change.
6. **Supervisory-letter drafting.** Use the language templates in §7. Pair with `docs/regulator-pack/finding-language.md` for general finding-phrasing standards.
7. **Follow-up examination cycle.** Use §8 for remediation tracking and historical-entry handling between examination cycles.

---

## §1. Pre-exam scoping — Examiner conversation guide

Sixty days before exam open, the EIC schedules a 60–90 minute scoping call with the institution's Chief Risk Officer or Chief Information Security Officer. The call covers the following:

### §1.1 Posture confirmation

- Which posture does the institution operate per spec §4.1.2 — FFIEC-conformant constants, vendor-namespaced constants, or a documented hybrid? The institution's CC8.1 control description names the choice.
- Has the posture changed during the examination period? If yes, request the change-management record at PBC delivery.
- Does the institution use a vendor-hosted implementation (Herald.Compliance or equivalent)? If yes, expect to apply the vendor-management questions in §3.10.

### §1.2 Scope confirmation for the examination period

- Examination period: typically 12 months (rolling) for tier-1, 18 months for community banks.
- Chain-period coverage: confirm the institution's seal-publication cadence has been continuous for the period. Gaps in seal publication ≥72 hours during the period are pre-flagged for §3.5 testing.
- Tenant scope: which tenants are in the institution's deployment? Multi-tenant institutions confirm fingerprint isolation per spec §4.1 (P-26 cross-check).

### §1.3 Control-change events during the period

- Master-key rotations: how many, dates, change-management record IDs.
- Algorithm rotations (spec §4.3.2 / §10.10.2): planned or executed during the period?
- Posture changes (FFIEC ↔ vendor-namespaced): planned or executed during the period?
- HSM changes (vendor change, hardware refresh): planned or executed during the period?
- Each control change is examined separately to confirm change-management discipline and continuity of integrity.

### §1.4 Examination logistics

- The EIC names whether they will run the verifier directly (auditor-run posture per the SOC 2 attestation overlay), or rely on institution-run verifier output with re-performance sampling, or rely on a current SOC report. The institution's IT Operations team prepares the corresponding access path.
- The EIC names sample-testing dates so the institution's IT Operations and IR teams are aware of when the examiner pulls evidence.
- Multi-agency coordination: if OCC and Federal Reserve are co-examining, the lead-agency designation per §3.10 is confirmed.

### §1.5 Coordination with active incident response

If the institution has an active chain-related incident (Scenario 3, 4, 9, or 11 per the IR playbook) at the time of pre-exam scoping, the EIC and the institution's IR coordinator agree on the daily-sync mechanic per `docs/incident-response-playbook.md` and the examination-versus-investigation role split.

---

## §2. PBC list — what the institution produces 5 business days before exam open

The institution delivers the artifacts below. Each artifact has an expected format and a coverage period.

### §2.1 Governance and policy artifacts

| Artifact | Format | Coverage period |
|---|---|---|
| Board-approved AI-decision-logging policy naming the chain as the control | PDF, dated | Current |
| CC8.1 control description for the chain (posture, configuration, change-management discipline, three-lines-of-defense alignment) | PDF, dated | Current |
| Three-lines-of-defense assignment matrix (see §3.9) | PDF or spreadsheet | Current |
| Anomaly-review committee charter and most recent 12 months of meeting minutes | PDF | Examination period |
| Last URSIT rating from prior examination (chain-related findings if any) | PDF | Trailing |

### §2.2 Operational evidence

| Artifact | Format | Coverage period |
|---|---|---|
| Daily seal records with HSM `signed_at` timestamps | CSV, columns: `seal_date, signed_at, merkle_root, signatures_count, key_version` | Examination period |
| Master-key rotation log (dates, change-management record IDs, attesting officers) | JSON or CSV | Examination period |
| Verifier output from internal-audit runs (P-13 evidence) | CLI output + summary report | Examination period |
| Operational events log (`master.*`, `seal.*`, `ledger.*`, `chain.*`, `hsm.*` per spec §10.2) | JSON or NDJSON | Examination period |
| `master.reconciliation_completed` events showing fingerprint reconciliation results (P-6) | JSON | Examination period |
| Any chain-detected incident records and the institution's IR-playbook execution evidence | PDF + JSON | Examination period |

### §2.3 HSM custody evidence

| Artifact | Format | Coverage period |
|---|---|---|
| HSM vendor audit log (key-load events, extraction attempts, access events) | Vendor-format export | Examination period |
| Change-management record of each IKM load with operator identity | PDF, signed | Per IKM operation |
| HSM attestation certificate (if vendor provides) confirming non-extractability | PDF | Most recent issued |
| HSM operator role list with seal-job operator separated from administrator | PDF | Current |

### §2.4 Vendor-management evidence (if vendor-hosted)

| Artifact | Format | Coverage period |
|---|---|---|
| Vendor risk-assessment completion form (institution-owned) | PDF | Most recent annual |
| Service agreement with seal-publication SLA (≤60 min), incident-response SLA (≤36 hours), audit-evidence retention SLA (≥7 years) | PDF | Current |
| Vendor's most recent SOC 2 Type II report or equivalent third-party attestation | PDF | Most recent issued |
| Vendor incident records during the examination period (if any) | PDF | Examination period |

### §2.5 Privilege handling

The institution's legal counsel may flag specific incident-response analysis as work-product or attorney-client privileged. The EIC and counsel agree on scope before exam open. Privilege scope is documented in the workpapers; disputes are escalated to the agency's regional counsel. The standard examination posture is that the chain's primary integrity evidence (seal records, verifier output, HSM logs, operational events) is operational and non-privileged; legal-analysis memoranda may carry privilege.

---

## §3. Examination procedures — overlay onto institution-side audit procedures

The institution operates the P-N procedures in `docs/audit-procedures.md`. The EIC's procedures cross-walk to those, with examination-side sampling and disposition. The EIC may either re-perform the institution's tests on the examiner's sample, OR sample the institution-run output and re-perform a sub-sample under AT-C §205 §A21 for reliability assessment.

### §3.1 Verifier sample-testing (overlays P-13)

The examiner runs the verifier (or relies on institution-run output with re-performance) over a sample of chain entries per §4.

- For each entry, verify PASS / FAIL with failure-reason granularity per spec §10.12 exit-code contract.
- Failure rates feed the severity matrix in §5.
- Random-seed, sample design, and exit-code distribution are documented in the workpapers.

### §3.2 IKM custody verification (overlays P-5)

- Confirm the HSM audit log shows IKM load events with documented operator identities.
- Confirm the change-management record for each IKM operation matches the HSM log.
- Confirm the institution's HSM operator role list shows seal-job operator separated from HSM administrator (spec §10.5 separation-of-duties; small institutions with documented dual-control compensating control are conformant).
- Confirm production builds exclude the software-key adapter (P-5 sub-procedure; spec §10.7).

### §3.3 Master-key rotation review (overlays P-6)

For each rotation in the period:
- Confirm the change-management record with approval signatures.
- Confirm the institution's `master.reconciliation_completed` event for the rotation period shows `fingerprint_unmatched_count = 0` after the rotation (or, if non-zero, that the institution's investigation per P-6 closed the discrepancy).
- Confirm the verifier produces PASS for entries on both sides of the rotation boundary.

### §3.4 Seal-age testing — 60-minute SLA (overlays P-30)

- Pull the daily seal records for the period.
- Per the §4 sampling table, sample at minimum 5 seals per month per tenant (stratified across day-of-week and the days following any HSM-availability event).
- For each sample: compute `signed_at - day_boundary`. The delta MUST be < 60 minutes.
- Apply the severity matrix in §5.4 (seal-age failure-rate disposition).

### §3.5 72-hour notification posture (overlays P-9)

- For any seal delay exceeding 72 hours during the period, confirm the institution's notification record matches the 12 CFR §53 Computer-Security Incident Notification Rule timing.
- Confirm the notification was filed via the agency's incident-reporting portal (not just an internal memo).
- Confirm the institution's IR playbook Scenario 3 was activated.

### §3.6 Anomaly evidence completeness (overlays P-22 through P-25, P-31, P-33)

For each anomaly category surfaced during the period, the EIC samples the institution's evaluation records per the relevant P-N procedure. Anomalies missing one or more of the four required evidence pieces are deficient regardless of root cause.

### §3.7 Incident-response playbook integration (overlays P-14, P-15, P-16)

- Confirm the institution's IR playbook addresses every scenario in `docs/incident-response-playbook.md` and is reviewed annually.
- For each chain-detected incident in the period: confirm the IR-playbook activation timing matched playbook targets (36-hour Computer-Security Incident Notification Rule per Scenario 4; 72-hour seal-delay notification per Scenario 3).
- Confirm post-incident review was conducted with documented root-cause analysis.

### §3.8 Operational events log retention (overlays P-10)

- Confirm operational events for the period are retained per spec §10.13 (7-year floor).
- Sample-test by querying the log store for old events; confirm immutability by attempting an UPDATE under the ledger writer role (P-1) and confirming rejection.

### §3.9 Three-lines-of-defense alignment

The FFIEC framework expects:
- **First line (IT Operations)** owns the seal job, HSM operations, IKM rotation schedule, daily verifier run. Escalates anomalies to the IT Risk Officer within 1 business day.
- **Second line (IT Risk or Model Risk Management Committee)** reviews chain-operation metrics monthly: seal-age trend, verifier failure rate, fingerprint reconciliation results. Assesses control-design vs operational-incident.
- **Third line (Internal Audit or SOC)** tests control operating effectiveness quarterly via P-13 sample verification. Issues an annual attestation to the audit committee.

The EIC confirms each line is named with an explicit owner and reviewer in the institution's CC8.1 control description. Vague or unassigned lines are documented as a control-program-maturity observation; absence of an explicit owner at any line is an MRA.

### §3.10 Vendor-management testing (TPRM Interagency Guidance 2023)

For institutions using vendor-hosted implementations:
- Confirm the vendor risk-assessment completion form rates the vendor based on criticality, security, financial stability.
- Confirm the service agreement carries the SLAs in §2.4.
- Confirm the institution received and reviewed the vendor's most recent SOC 2 Type II report.
- Confirm the vendor's HSM is FIPS 140-2 Level 3 or higher per spec §10.5.
- Confirm the vendor's IR plan includes FFIEC notification within 36 hours per spec §10.5 / 12 CFR §53.
- Confirm the institution has a failover or second-vendor arrangement if seal-publication delay exceeds 72 hours (matching the IR Scenario 3 escalation).

### §3.11 Multi-agency examination coordination

When OCC and Federal Reserve (or FDIC and CFPB) are co-examining:
- Lead agency is named at scoping (typically the primary supervisor).
- Supporting examiners reference the lead agency's PBC list rather than requesting duplicate data.
- The supervisory letter is issued by the lead agency; supporting agencies may issue supplemental jurisdiction-specific findings.
- The institution's remediation plan responds to all agencies' findings in a single coordinated response.

### §3.12 SOC report incorporation (AT-C §205 §A19)

If the institution holds a current (issued within 12 months) SOC 2 Type II report covering the chain control:
- Review the SOC report's scope (period, in-scope controls, attest opinion).
- Assess the SOC auditor's testing procedures per AT-C §105 sufficiency.
- Request the SOC auditor's observations or control-design recommendations.
- Re-test only (a) deficiencies the SOC report noted, (b) subsequent-period controls (post-SOC), and (c) management remediation evidence.
- Document in the workpapers which tests relied on the SOC report and which the EIC performed independently.

A clean SOC report typically narrows EIC re-testing to current-state evidence. A qualified SOC report typically expands EIC re-testing in the affected control areas.

---

## §4. Sampling methodology — examiner side

### §4.1 Sample-size guidance by institution size

| Institution category | Sample size per quarter | Sample size per year |
|---|---:|---:|
| Community bank (< $10B assets) | 50 entries | 200 entries |
| Regional bank ($10B–$250B assets) | 100 entries | 400 entries |
| Tier-1 institution (> $250B assets) | 300 entries | 1200 entries |

### §4.2 Sample design

For each quarter in the period:
- Identify 3–5 random dates.
- For each date, select 25–50 entries at random (per the institution-size table).
- Run the verifier on every sampled entry.
- Document sample design, random-number seed, sampling tool, and exit-code distribution in the workpapers.

### §4.3 Stratification

The sample is stratified across:
- **Tenants.** At least 5 entries per tenant.
- **Time.** At least one sample date per month.
- **Day-of-week.** At least 3 weekdays observed across the sample.
- **Around control-change events.** Additional entries from days adjacent to master-key rotations, posture changes, HSM swaps.
- **Around incident-identified periods.** 50–100 supplementary entries from any period flagged by the institution's IR program.

### §4.4 Confidence target

The sample design SHOULD achieve 90%+ confidence of detecting a 1% failure rate in the population. Smaller institutions with smaller volumes may run a census of high-risk strata plus a stratified sample of the broader population.

### §4.5 Community bank scoping flexibility

For community banks (< $10B), the EIC may scope the examination to:
- A single business line (e.g., the auto-lending AI using the chain).
- A pilot period of one quarter (with commitment to expand in the next cycle).
- A 30–50 entry year-1 sample, expanding to 200 entries in subsequent examinations.

The institution and the EIC agree on scope at pre-exam scoping. The supervisory letter notes the scoping decision.

### §4.6 Tier-1 institution scoping flexibility

For tier-1 institutions, the EIC may segment the examination by business line or region:
- Sampling and testing delegated to regional examination teams under EIC supervision.
- Focused exam of one critical business line as a pilot, with commitment to expand.
- Sample-of-samples approach: institution's internal audit performs verifier runs weekly; the EIC spot-checks the internal-audit results.

The 4-week examination window for tier-1 institutions typically allocates 20–40 hours of EIC time across the four weeks for chain-specific testing, plus 4 hours pre-exam scoping and 8 hours post-exam analysis.

---

## §5. Deficiency-severity matrix — examination-side disposition

Apply the matrix to each finding. The matrix is binding for examination consistency across the FFIEC system.

### §5.1 Verifier failure rate

| Failure rate (sampled entries) | Disposition |
|---|---|
| 0% | Control operating effectively; no finding |
| > 0% but < 0.1% | Operational finding; documented in workpapers; no MRA unless remediation was inadequate |
| 0.1% – 1% | Control-design deficiency; MRA with 90-day remediation plan |
| > 1%, OR multiple failures on the same entry | MRIA; immediate action plan; may restrict further AI deployment pending fix |
| Structural failure (missing seal record, negative-test-vector match) | MRIA; immediate escalation to RFI; possible consent-order discussion |

### §5.2 Fingerprint mismatch (P-22 / P-6)

| Pattern | Disposition |
|---|---|
| Procedural cause (IKM rotation procedure did not update the fingerprint registry); institution's investigation closed the discrepancy in real time | Operational finding |
| Procedural cause; investigation lagged or was inadequate | MRA |
| Cryptographic cause (HMAC computation defect, key-derivation defect) | MRIA |
| Cross-tenant key swap from vendor-side incident | MRIA + vendor-management escalation |

### §5.3 Seal-age SLA breach (P-30 / P-8)

| Breach pattern | Disposition |
|---|---|
| 0 breaches | Control operating effectively |
| 1–2 seal ages exceeding 60 min, with documented incident response and remediation | Operational finding |
| 3+ seal-age failures in the period | MRA (HSM capacity or seal-job timing requires re-balance) |
| Seal-age failures correlating with a documented HSM event (e.g., scheduled patching) and the institution's change-management documented the downtime and recovery time | Operational finding (no MRA) |
| Seal delay exceeding 72 hours without 12 CFR §53 notification | MRIA |

### §5.4 Three-lines-of-defense alignment

| Pattern | Disposition |
|---|---|
| All three lines explicitly named with owner / reviewer; documented escalation procedures | Control operating effectively |
| Lines named but escalation procedures vague or undocumented | Operational finding |
| One or more lines lack an explicit owner | MRA |
| Two or more lines unassigned, OR independence between first/third line broken | MRIA |

### §5.5 Vendor-management findings

| Pattern | Disposition |
|---|---|
| Vendor risk assessment current; SLAs documented; SOC report received and reviewed | Control operating effectively |
| One missing element (SLA, SOC report, or risk assessment) | Operational finding |
| Two+ missing elements | MRA |
| No vendor risk-assessment process; no failover plan; no SLA | MRIA |

### §5.6 Posture documentation (CC8.1)

| Pattern | Disposition |
|---|---|
| CC8.1 names posture, configuration mechanism, change-management discipline | Control operating effectively |
| One element missing | Operational finding |
| CC8.1 missing entirely, OR posture in operation does not match CC8.1 | MRA |
| Multiple postures observed concurrently in production without documented hybrid posture | MRIA |

### §5.7 Incident-response timing

| Pattern | Disposition |
|---|---|
| IR Scenario 4 (master-key compromise) — 36-hour notification met | Conformant |
| IR Scenario 4 — notification beyond 36 hours but within reasonable extension | Operational finding |
| IR Scenario 4 — notification beyond 72 hours, OR notification not filed | MRIA + 12 CFR §53 violation |
| IR Scenario 3 (seal not published within 72 hours) — notification met | Conformant |
| IR Scenario 3 — notification beyond 72 hours | MRIA |

---

## §6. URSIT rating integration

The Uniform Rating System for Information Technology rates each IT function on a 1–5 scale. The chain's operating effectiveness affects the **Operations** function rating (the chain is part of the audit-trail control set). The four URSIT factors — Audit, Management, Development & Acquisition, Support & Delivery — each absorb chain evidence differently.

### §6.1 URSIT impact assessment criteria

The EIC evaluates whether the chain's presence and operating effectiveness warrant an Operations rating change:

1. **Risk reduction.** Does the chain materially reduce the likelihood of undetected AI-decision tampering?
2. **Operational integration.** Does the chain operationally integrate with the institution's existing incident-response program?
3. **Evidence quality.** Does the chain produce evidence that supports the SOC team's third-line testing?

If all three answers are yes, the chain warrants a rating upgrade for Operations:
- From 3 (Satisfactory with residual risk) to 2 (Strong) — typical for first-time deployments operating effectively in their first year.
- From 2 (Strong) to 1 (Outstanding) — for institutions with mature AI governance, no chain-related findings during the period, third-line attestation in place, and full integration with the IR program.

### §6.2 Negative URSIT impact

If the chain has chain-related MRA or MRIA findings during the period, the Operations rating may be downgraded:
- One MRA + adequate remediation plan → no rating change (MRA is the supervisory action).
- Multiple MRAs OR an unremediated MRA from the prior cycle → downgrade by one notch.
- An MRIA OR a 12 CFR §53 violation → downgrade by one or two notches; possible consent-order trigger.

### §6.3 URSIT documentation in the supervisory letter

The supervisory letter names the chain's contribution to the Operations rating explicitly. The presence of the chain alone does not automatically change the rating; the integration with the IR program, demonstrated SOC testing, and absence of control-failure findings during the period are the basis for the upgrade.

---

## §7. Supervisory-letter language templates

The following templates produce consistent language across examiners. Pair with `docs/regulator-pack/finding-language.md` for general finding-phrasing standards.

### §7.1 Control operating effectively

> We confirmed that the institution's AI-decision logging control (chain-of-custody per FFIEC guidance) operates with appropriate technical controls, management oversight, and third-party assurance. We performed sample testing of [N] chain entries using the published verifier and confirmed [outcome — typical: 100% passed integrity verification]. The institution's incident-response procedures for chain-detected anomalies are appropriately documented and tested. The control will be revisited in the next examination.

### §7.2 Control deficiency — procedural

> We noted that the institution's control description for the chain-of-custody system did not include [specific gap, e.g., the change-management procedure for HKDF posture changes]. Corrective action is required within 90 days. The institution should [required action].

### §7.3 Control deficiency — operational

> During our sample testing of the chain-of-custody entries, we identified [N] entries with [failure mode, e.g., fingerprint mismatches]. Investigation confirmed the defect was a procedural issue [description] that the institution remediated on [date]. The institution's incident-response procedures appropriately detected and escalated this issue. We will monitor the institution's remediation in the follow-up examination.

### §7.4 Control deficiency — design (MRA)

> We identified a design defect in the institution's chain-of-custody control [description]. This is a Matters Requiring Attention (MRA) and requires corrective action within 90 days. The remediation plan should include [required steps].

### §7.5 Critical deficiency (MRIA)

> We identified [description of structural or critical chain-integrity defect, e.g., verifier failure rate exceeding 1% of the sample, OR a `dev_mode=true` seal in production, OR a seal delay exceeding 72 hours without 12 CFR §53 notification]. This is a Matter Requiring Immediate Attention (MRIA) and requires immediate action. The institution must submit an action plan within [period, typically 30 days] and may be restricted from further AI deployment pending fix.

---

## §8. Follow-up examination workflow

### §8.1 MRA remediation plan submission

When the institution responds to an MRA, the submission includes:
- Description of the root cause.
- Corrective action taken or planned.
- Timeline for completion (typically 90 days).
- Evidence of completion (updated procedure document with approval date; first-execution evidence of the new procedure).

### §8.2 Evidence retention for follow-up

The institution retains remediation evidence for the follow-up examination:
- Updated control procedure document with approval signatures.
- First month of execution evidence (e.g., weekly reconciliation reports for the month after implementation).
- Training or communication to staff about the updated procedure.

### §8.3 Follow-up testing

The follow-up examiner samples the institution's ongoing execution of the remediated control (typically pulls 3 months of weekly evidence from the follow-up exam period) to confirm the procedure is still operating effectively.

### §8.4 Follow-up disposition

| Pattern | Disposition |
|---|---|
| Remediation evidence shows completion; follow-up sampling shows ongoing compliance | Close the MRA |
| Follow-up sampling shows the defect has recurred | Re-issue as new MRA, OR escalate to MRIA if it is a second repeat |
| Remediation incomplete at follow-up | MRA remains open with extended deadline; possible URSIT downgrade |

### §8.5 Historical-entry handling after remediation

When corrective action addresses a chain-integrity defect (fingerprint mismatch, seal-age SLA breach, missing seal record), the institution may treat affected entries as historical artifacts only if the root cause was procedural, not cryptographic.

- **Procedural root cause.** The fingerprint registry was not updated when the IKM was rotated. The institution updates the registry and documents the correction in its IR log. Historical entries remain in the ledger with the old fingerprints — the verifier still rejects them, correctly identifying the procedural-error period — and the institution's incident response explains the defect and the remediation. The verifier's failure on those entries is the audit trail of the procedural error.
- **Cryptographic root cause.** The HMAC computation was buggy. The verifier output cannot be trusted for historical entries. The institution engages the auditor and the examiner to determine whether the defect scope can be bounded. The follow-up examination focuses on confirming the fix, not on re-hashing the historical entries (the spec and verifier do not require historical re-hashing).

### §8.6 Cross-examination evidence reference

The 7-year retention requirement (spec §10.13) is on the institution. The examiner's workpapers retain evidence relating to the current period plus cross-referenced evidence from prior periods that bears on current findings. The examiner does not transfer original chain evidence; the institution retains custody throughout.

---

## §9. Examination scheduling and incident coordination

### §9.1 Active incident during examination

If a chain-related incident is detected during an active examination:
- The institution's IR coordinator notifies the EIC immediately (phone + email).
- The EIC notifies the examination supervisor (the EIC's backup) by the same channel.
- The institution's IR team simultaneously initiates the standard 12 CFR §53 36-hour notification per Scenario 4 of the IR playbook, submitted via the agency's incident-reporting portal so it enters the regulatory record properly.
- The EIC and institution's IR coordinator establish a daily 8:00 AM sync to coordinate examination-team and IR-team roles.

### §9.2 Control changes during examination

The institution avoids scheduling control changes (master-key rotation, HSM replacement, posture change, seal-job re-platform) during an active examination. If unavoidable:
- Notify the EIC at pre-exam scoping (60 days before the exam).
- Provide planned change procedure and timeline.
- The EIC adjusts the sample-testing procedure: pre-change and post-change cohorts are sampled separately; the institution provides both old and new IKM for verification (per HSM custody protocols).

### §9.3 Emergency control change during examination

If an emergency change is required (e.g., HSM failure):
- The institution notifies the EIC within 4 hours.
- The EIC suspends examination testing for the day to allow the institution to execute the emergency change.
- Examination testing resumes the following day.
- The emergency change is documented in the workpapers as a control-change event and reported in the supervisory letter.

---

## §10. Confidential Supervisory Information (CSI) and disclosure guidance

The supervisory letter is CSI under FOIA exemption (b)(8). The institution may use it internally for risk management and board reporting, but cannot publish it.

### §10.1 Permissible disclosure language

The institution MAY cite the chain in public disclosures (annual report, proxy statement, investor call) with language like:

> We operate a cryptographically verified audit trail for AI decisions (chain-of-custody control per FFIEC guidance) to ensure the integrity of AI-driven decision-making. This control is subject to internal audit, external audit, and examination by our primary regulator.

The institution MAY cite the existence of an examination and the high-level conclusion:

> The Federal Reserve examination in [year] confirmed our AI-decision logging and governance controls are operating effectively.

### §10.2 Impermissible disclosure

The institution CANNOT cite specific findings, control deficiencies, or MRA language from the supervisory letter.

### §10.3 Remediation disclosure

If a major finding is remediated and the remediation is verified in a follow-up exam or SOC report, the institution may cite the remediation:

> We identified and remediated a procedural deficiency in [area] in [year]; the remediation was confirmed by [auditor / examiner] in [year].

This allows institutions to tell a positive story about AI governance while respecting CSI confidentiality.

---

## §11. Examination cost and feasibility

### §11.1 Time budget

| Institution category | EIC chain-specific time | Institution support time |
|---|---|---|
| Community bank (< $10B) | 50–60 hours | 20–30 hours |
| Regional bank ($10B–$250B) | 100–150 hours | 40–60 hours |
| Tier-1 (> $250B) | 200–400 hours (distributed) | 80–160 hours |

### §11.2 Pre-exam scoping (60 days before)

| Activity | EIC hours | Institution hours |
|---|---:|---:|
| Scoping call | 2 | 2 |
| PBC list assembly and review | 2 | 8 |

### §11.3 Post-exam analysis

| Activity | EIC hours |
|---|---:|
| Workpaper finalization | 4 |
| Supervisory-letter drafting | 4 |
| Follow-up scheduling | 1 |

---

## §12. Examiner training reference

The EIC's foundation in the chain is established in `docs/regulator-pack/examiner-training.md`. Examiners new to the chain complete the training module before participating in chain-specific examination procedures. The training covers:

1. Chain primitives (per-event MAC, daily Merkle, HSM signature, wire format) — 30 minutes overview + 2 hours spec walk-through (§1–§4).
2. Verifier procedure and failure modes (negative test vectors) — included in spec walk-through.
3. Audit procedures P-1 through P-N — 1 hour.
4. Examination toolkit (verifier CLI, flags, input/output, workpaper templates) — 1 hour.
5. Hands-on lab — 2 hours running the verifier on test chains.
6. Common findings and remediation case studies — 1 hour.

The training is delivered by the Federal Reserve's examination schools, OCC training, and the FDIC's examiner-development program. Refresher training is annual.

---

## §13. Inter-agency coordination summary

When OCC + Federal Reserve, OR FDIC + Federal Reserve, OR FDIC + CFPB are co-examining:

- Lead agency designation at scoping.
- Single combined PBC list to the institution.
- Single set of sample tests; supporting examiners reference lead-agency workpapers.
- Single supervisory letter from the lead agency; supporting agencies issue jurisdiction-specific supplemental findings if they have unique authority (e.g., CFPB on consumer-facing AI-decision disclosure).
- Single coordinated remediation plan from the institution responding to all agencies' findings.

The coordination prevents the institution from being examined twice on the same control and makes regulatory burden predictable.

---

## §14. Recommended FFIEC working-group enhancements (informative)

Enhancements that the FFIEC working group may publish to standardize examination across the agencies:

1. **Examiner Guidance for AI Decision Logging** — a 6–10 page addendum to the spec adapting it to the FFIEC examination framework. This document is the substrate.
2. **Examiner Training Module** — half-day curriculum delivered through Federal Reserve, OCC, and FDIC examiner schools.
3. **Examination-tool automation** — packaged toolkit integrating verifier execution, sample selection, and workpaper documentation. A future evolution (2028+) for tier-1 institutions.

These are working-group recommendations, not normative requirements. v1.0a institutions and examiners operate under the existing spec + this overlay; the working-group enhancements layer on top without changes to the normative core.
