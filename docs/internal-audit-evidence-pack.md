# Internal audit evidence pack

> **What this doc is.** A Chief Audit Executive-facing evidence pack covering how chain-of-custody artifacts integrate into the IIA's International Professional Practices Framework. Specifically: how the third line of defense plans audits of the chain (IIA 1100/1200/2000/2200), gathers evidence (IIA 2300), reports findings (IIA 2400), and works in combined assurance with the second-line risk and compliance functions and with external audit. The document also covers chain-derived Key Risk Indicators feeding the quarterly audit-committee reporting cadence.

> **Audience.** Chief Audit Executive (CAE), Director of IT/Operations Audit, IT Audit Manager assigned to chain-control testing, the audit-committee liaison, and the institution's external audit partner where the AI is in SOX 404 ICFR scope.

> **What this pack does NOT do.** It does not change the spec's normative core. It does not replace the institution's existing audit charter, Risk-Based Internal Audit plan, or Quality Assurance and Improvement Program. It is the audit-governance layer that sits above the spec — the playbook the third line uses to apply IIA Standards to chain evidence with confidence.

---

## §A1. Purpose and IIA-Standard grounding

The chain-of-custody specification is audit-ready at the third-line level. The verifier procedure (spec §7) is independently testable by internal audit without auditee knowledge of sampling. The chain-entry sequence is deterministic and unchangeable; audit can spot-check entries without the operations team knowing which entries are being tested until the audit report is issued. The operational-events schema provides the hook for audit-data-analytics queries of the aggregated ledger. The incident-response playbook integrates with the audit committee's quarterly reporting cadence.

What the spec does not explicitly address is the audit-governance layer — how internal audit plans, scopes, samples, evaluates, documents, and reports under IIA Standards. This pack closes that layer. The grounding standards are:

- **IIA 1100 Independence and Objectivity.** §A4 covers ledger-access independence and the segregation-of-duties controls audit operates.
- **IIA 1220 Due Professional Care.** §A4, §A8, §A11 cover sampling rigor, documentation, and contemporaneous workpaper practice.
- **IIA 2100 Nature of Work / 2120 Risk Assessment.** §A2, §A3 cover risk-based audit planning for the chain.
- **IIA 2200 Engagement Planning.** §A6 covers pre-implementation review (PIR).
- **IIA 2300 Performing the Engagement.** §A5 control-identification, §A8 sampling, §A9 root-cause analysis, §A11 workpapers.
- **IIA 2310 Identifying Information.** Sufficiency, reliability, relevance, usefulness — covered throughout.
- **IIA 2400 Communicating Results / 2410 Criteria for Communicating.** §A12 covers audit-committee reporting and the annual opinion.
- **IIA 2500 Monitoring Progress.** §A13 covers repeat-finding and remediation-effectiveness testing.
- **IIA 2600 Communicating the Acceptance of Risks / Quality Assessment Review.** §A14 covers QAR-defensibility of chain-audit procedures.
- **PCAOB AS 2201 + AS 1301.** §A7 covers SOX 404 ICFR admissibility when AI is in financial-reporting scope.

The pack is organized so a CAE can hand individual sections to specialty leads (IT auditor, model-risk auditor, external-audit liaison) without requiring the recipient to read the whole document.

---

## §A2. Risk-based audit planning for the chain

Internal audit maintains a risk-based audit plan per IIA 2120. The plan allocates audit effort to high-risk, high-impact areas. When the institution deploys the chain, internal audit decides whether the chain is a standalone audit-plan entry or a complementary control supporting other primary controls (the AI Governance Committee's oversight of AI decision-making, for example).

The plan-entry decision depends on three factors:

**Materiality of AI decisions covered.** If the chain covers high-materiality decisions (model output feeds into ICFR, capital models, or lending decisions), the chain audit warrant is annual detailed testing plus continuous-auditing monitoring. Medium-materiality decisions (customer-facing decision logging, policy decisions) warrant biennial detailed testing with annual monitoring. Low-materiality decisions (operational logging) warrant triennial detailed testing.

**First-year deployment premium.** In the chain's first year post-implementation, internal audit increases the annual testing scope to confirm the control is operating as designed and the auditee's control description matches reality. Scope includes the pre-implementation review (§A6), operational readiness testing, and a full year-1 verifier run. After year 1, the plan reverts to the materiality-driven cadence.

**Governance maturity.** Mature AI Governance Committee (meets monthly, escalates findings, tracks remediation per the three-lines framework) lowers the chain-audit risk and may shift cadence to biennial. Immature governance keeps the warrant annual. The CAE assesses governance maturity at plan-development time.

### §A2.1 Audit-plan entry template

The plan entry reads:

> Internal Audit will test the chain-of-custody control [annually | biennially | triennially] per `docs/audit-procedures.md` and `docs/internal-audit-evidence-pack.md`. Objectives: (a) confirm the verifier operates per spec §7; (b) sample chain entries to confirm integrity; (c) test the incident-response playbook has been exercised; (d) assess compliance with IIA Standards 2310 evidence requirements. Continuous-auditing queries (§A10) operate between annual engagements for real-time anomaly monitoring.

### §A2.2 ADA integration in the plan

The chain's aggregated ledger exposes query APIs (per `docs/vendor-hosted-controls.md` audit-rights template and Compliance design §10.3) enabling audit to run continuous queries. Internal audit's ADA tool-kit includes templates for: daily seal-age queries, weekly verifier-failure trend queries, monthly key-rotation reconciliation queries, ad-hoc incident-window forensic queries.

The plan entry names which ADA queries operate continuously and which are run on annual cadence. A representative split: daily seal-age and weekly verifier-failure trend queries operate continuously; monthly key-rotation reconciliation queries operate continuously; deep-dive sampling and root-cause analysis operate within the annual engagement.

---

## §A3. Audit universe segmentation and the chain's place

The chain lives in the IT Audit / Operations Audit segment of the audit universe. It also intersects:

- **SOX 404 ICFR audits** when AI is in financial-reporting scope (§A7).
- **Model risk audits** under SR 11-7 framing — model-risk auditor confirms the scope of AI decisions logged matches the institution's model inventory; chain-audit confirms the captured records are unaltered.
- **Vendor-management audits** when the chain is vendor-hosted — vendor-audit lead consumes the vendor's SOC report and the vendor-conformance attestation; chain-audit confirms the institution's CUECs operate.
- **Cybersecurity audits** — chain's cryptographic foundation and incident-response posture are inputs to the cybersecurity audit's annual work plan.

Cross-functional coordination (§A10) prevents duplication and ensures findings are not siloed by specialty.

---

## §A4. Independence and ledger access

Internal audit's independence is foundational under IIA 1100. The chain ledger (per Compliance design §10.3) is a database table held by the institution. The first-line operations team controls write-access. When internal audit samples chain entries to test integrity, audit must query the ledger — and audit must do so independently of the operations team.

### §A4.1 Segregation-of-duties posture

The institution operates two database roles:

1. **`chain_ledger_writer` (or equivalent).** Held by the seal-job and ingestion pipeline. INSERT/SELECT only on ledger tables and operational-events tables. No DELETE, no UPDATE.
2. **`internal_audit_chain_viewer`.** Held by the audit team. SELECT only on ledger tables, operational-events tables, and daily-seals tables.

The audit-query role is provisioned by the DBA or security team (second or third line), NOT by the chain-operations team (first line). Audit may rotate the audit-query role's credentials at least quarterly per the institution's access-control policy. Audit confirms the credential rotation log to ensure the operations team did not interfere.

### §A4.2 Pre-engagement independence test (P-A1)

Before executing chain-audit procedures, audit confirms:

1. The audit-query role is active and audit-team members can query the ledger independently.
2. The role's creation record (in the database-access-control log) shows it was provisioned by DBA, not by the chain-operations team.
3. There is no shared login between the audit-query role and the operations-team role.
4. The audit-query role's last credential rotation was within the institution's documented cadence.

Failures escalate as a segregation-of-duties deficiency requiring remediation before audit testing proceeds. Audit may be unable to test independently per IIA 1220 if the operations team has administrative access to the audit-query role or to the query-audit logs.

### §A4.3 Sample-selection independence

When audit selects a sample for verifier testing, the sample IDs are recorded in audit's workpapers before any ledger query is run. The sample is not communicated to the operations team until after verifier results are obtained. This prevents operations from pre-computing verifications for the sampled entries (a form of cooperation that, while well-intentioned, would compromise the audit's independence assertion).

---

## §A5. Control identification and design assessment

Per IIA 2130, audit identifies the control, understands its design, confirms it is operating, and assesses its effectiveness. For the chain:

### §A5.1 Control objective

The institution captures all AI agent decisions in a tamper-evident, independently verifiable ledger such that any unauthorized modification to a captured decision is detectable by an independent party (audit, regulator, or customer) without trust in the institution's infrastructure or processes.

### §A5.2 Key design elements (must all be present)

1. **Per-event integrity.** Every AI decision is captured in real time with a HMAC-SHA-256 MAC computed by the SDK.
2. **Per-run chaining.** Entries within a run form a chain (`prev_hash` field) such that any missing or reordered entry breaks the chain.
3. **Daily sealing.** Every day's entries are summarized in a Merkle tree, and the root is signed by the HSM on behalf of the institution.
4. **Independent verifiability.** The verifier procedure (spec §7) uses only the institution's published public key (no private-key access required) to validate all prior assertions.

### §A5.3 Operational evidence audit confirms

- Daily seal records with timestamps confirming seals are published within 60 minutes of UTC midnight.
- Institutional policy (board-approved) naming the chain as a control for AI-decision integrity.
- Control description (CC8.1 or equivalent) documenting the chain's design, the IKM custody location, the seal-job SLA, and the incident-response procedure.
- Change-management log showing any changes to the control's configuration are documented and approved.

### §A5.4 Operating-effectiveness tests audit runs

- P-11: Internal audit ran the verifier at least once during the period on a representative sample of entries.
- P-30: Daily seals meet the 60-minute SLA (stratified sample of 5+ seals per month across the period).
- P-6: IKM fingerprint reconciliation operates at the documented cadence with zero unresolved fingerprint mismatches.
- P-17: Any vendor-provided chain implementation has a current SOC 2 report covering the controls named in the spec.

### §A5.5 Effectiveness conclusion

The control is effective if: 100% of sampled chain entries pass the verifier procedure §7 (no HMAC mismatches, no chain-link breaks); 100% of sampled daily seals have signed_at timestamps within 60 minutes of the UTC day boundary; reconciliation procedures run at the documented cadence with zero unresolved fingerprint mismatches; any incidents had documented IR responses per the playbook.

---

## §A6. Pre-implementation review (PIR)

IIA 2200 requires audit to scope the engagement. When the institution deploys the chain for the first time, the CAE conducts a pre-implementation review 4-6 weeks before the chain goes live in production. The PIR is a gate; the chain does not go live without internal audit's PIR sign-off.

### §A6.1 PIR checklist

**Specification conformance.**

- The institution has selected a chain implementation conformant to FFIEC chain-of-custody-v1.0.
- The institution has documented which spec sections (§4.1.1 posture selection, §10.6 IKM key length, §10.10 key-rotation procedure) apply.
- Audit has obtained and reviewed the chosen implementation's source code or the vendor's attestation that the implementation is conformant.

**IKM custody and key management.**

- The institution has selected the HSM (cloud or on-premises) and confirmed FIPS 140-2 L3 certification.
- The IKM derivation procedure (HKDF per spec §3) is documented and tested in a non-production environment.
- The institution has run the spec's per-tenant determinism test: derive the session key twice with the same IKM and confirm byte-identical output.
- The IKM rotation procedure is documented and has been dry-run against a non-production tenant.

**Per-event capture and MAC computation.**

- The SDK is configured to emit all required attributes per spec §4.4.
- A sample AI decision has been captured and the entry's `payload_hash` has been independently recomputed by audit (using spec §5 canonical-form rules) to confirm byte-level accuracy.
- The OTLP transport (HTTP or gRPC) is configured with TLS 1.3 and tested in a staging environment.

**Daily sealing and Merkle root signature.**

- The seal-job is deployed and has run in staging for at least one week, producing daily seal records with Merkle roots.
- Audit has verified that the seal-job runs within 60 minutes of UTC midnight per spec §4.3.
- The HSM is configured to sign the Merkle root, and the signature can be independently validated using the institution's published public key.
- The daily seal record includes the required fields: `seal_date`, `merkle_root`, `root_signature`, `signed_at`, `entries_count`.

**Ledger aggregation and query API.**

- The ledger storage is deployed and audit has been granted a SELECT-only query role per §A4.
- Audit has tested the query API and confirmed entries are returned with all canonical-form fields intact.
- The ledger is confirmed to be append-only (UPDATE and DELETE roles are denied).

**Verifier binary.**

- The institution has obtained the verifier binary per spec §7 or a vendor implementation.
- Audit has confirmed the binary is reproducible-build compatible: the institution has rebuilt from source and confirmed byte-level equivalence.
- Audit has run the verifier against a sample of chain entries from staging and obtained PASS results.
- The verifier's output format matches spec §7.

**Control description and policies.**

- The institution has drafted a CC8.1 (or equivalent) control description naming the chain.
- The institution's board or audit committee has approved the control description in writing.
- The incident-response playbook has been reviewed and approved.

**Audit readiness.**

- Audit has drafted the chain-audit procedures specific to the institution's implementation.
- Audit has confirmed it has the tools (verifier binary, SQL client, ADA query templates) to execute procedures independently.
- Audit's charter or engagement letter has been updated to include chain testing in the annual audit plan.

### §A6.2 PIR sign-off memo

At the conclusion of the PIR, the CAE prepares a memorandum to the CFO and CRO:

> Internal Audit has completed the pre-implementation review of the chain-of-custody control. Design review confirmed the implementation aligns with the FFIEC spec. Operational readiness was confirmed through staging-environment testing. The control is approved for production deployment. Audit will execute the control-testing procedures as planned in the [date] annual audit engagement. Outstanding items: [list any].

Any design or operational gaps identified in the PIR are recorded as pre-implementation audit findings and tracked for remediation before go-live.

---

## §A7. SOX 404 ICFR admissibility — when chain evidence supports an ICFR opinion

When AI is used for high-risk financial decisions (credit-scoring, allowance-for-loan-loss estimation, economic-capital calculations, ALM), those AI decisions are part of ICFR scope. PCAOB AS 2201 (An Audit of Internal Control over Financial Reporting) requires the external auditor to obtain sufficient evidence about ICFR design and operating effectiveness. The chain is a control over decision-capture integrity; it is not a control over the decision's correctness.

### §A7.1 What the chain IS admissible for under PCAOB AS 2201

- **Operating effectiveness of decision-capture integrity controls.** Audit samples AI-decision records, runs the verifier, and concludes whether the captured records are tamper-evident.
- **IT general controls assessment.** Audit cites the chain's key-rotation procedures, seal-job change-management, and incident-response logs as evidence supporting IT general controls.
- **Walk-through testing of the decision-making process.** When the external auditor traces a financial transaction through the AI-driven decision process, audit provides the chain record as the tamper-evident audit trail.

### §A7.2 What the chain is NOT admissible for

- **Correctness or appropriateness of the AI decision itself.** That is model-risk-management scope per SR 11-7.
- **Model training-data quality or bias properties.** Model-risk controls, not IT general controls.
- **Model-validation requirements.** Backtesting, sensitivity analysis, governance review per SR 11-7.

### §A7.3 External-audit coordination

Before year-end SOX 404 testing, internal audit produces a summary for the external auditor:

1. Chain-control testing results from the period (verifier pass rates, seal-age SLA compliance, reconciliation results).
2. Any chain-detected incidents and their IR responses.
3. Confirmation that the chain's design and operation align with the spec the external auditor will rely on.

### §A7.4 Audit-opinion language for AI-in-ICFR scope

> We assessed the operating effectiveness of the chain-of-custody control for AI decision-capture integrity. We sampled [N] decisions across [date range], ran the independent verifier procedure on each, and obtained PASS results for all sampled entries. No ICFR weaknesses related to decision-capture integrity were identified in our testing. This opinion covers the integrity of the capture audit trail; it does not cover the correctness of the AI decisions themselves, which are subject to model-risk-management controls under SR 11-7.

---

## §A8. Sampling methodology

IIA 2330 (documentation) and 2310 (evidence sufficiency, reliability, relevance) require audit to document the sampling design, the population, and the basis for concluding the sample is sufficient. The chain contains thousands or millions of entries per year. Audit's sample must support a defensible conclusion.

### §A8.1 Population definition

The population is all captured AI decisions within the audit period (typically 12 months). Audit defines scope based on risk:

- High-materiality decisions (model output feeds ICFR): population is all instances of that decision type within the period.
- Medium-materiality decisions: population is all decisions within the period.
- Low-materiality decisions: population may be stratified by tenant or business line.

### §A8.2 Sample-size determination

Use attribute sampling with:

- **Confidence level:** 90% (standard for audit evidence).
- **Acceptable error rate:** 1% (escalation threshold if sampled error rate is at or above this).
- **Expected population error rate:** 0%.
- **Sample size:** for a population of 1 million entries, ~300 entries. For populations under 100k, typically 50-100. CAE consults IIA or AICPA sampling tables per the institution's audit methodology.

### §A8.3 Stratification

- **By day:** divide the audit period into weeks or months; stratify the sample so at least 3-5 weeks are represented.
- **By business line:** if the chain covers multiple AI use cases, allocate the sample across use cases per their materiality.
- **By source:** if the chain covers multiple AI platforms, stratify by source to ensure the control operates consistently.

### §A8.4 Sampling execution

1. Audit uses a random-number seed (documented and retained in workpapers) to select the sample. Simple random or systematic (every nth entry).
2. For each sampled entry, audit records the entry's `(tenant_id, run_id, seq)` identifier in workpapers BEFORE running the verifier. This prevents the operations team from knowing which entries are being tested.
3. Audit runs the verifier on each sampled entry per spec §7 and records the PASS/FAIL result.

### §A8.5 Evaluation

- **0 failures.** Control is operating effectively for the sampled population. Conclude with 90% confidence that the error rate in the full population is below 1%.
- **1 failure.** Review the root cause (§A9). Isolated incident with documented cause may leave the control effective; systemic cause escalates to a control-design assessment.
- **2 or more failures.** Escalate to management. Expand the sample to determine the population error rate. If expanded error rate is at or above 1%, issue an MRA (control-design deficiency).

### §A8.6 Documentation

Workpapers document: population definition and period; sample size and stratification; random-number seed or sampling method; list of sampled entry IDs; verifier results; root-cause analysis for any failures; conclusion and recommendation.

---

## §A9. Root-cause analysis decision tree for verifier failures

When the verifier returns a failure, audit determines whether the failure is a control-design deficiency, a control-operation failure, or a security incident. IIA 2320 requires analysis and evaluation; counting failures is not sufficient.

### §A9.1 Verifier failure types (per spec §7)

- (a) Format-version mismatch.
- (b) HKDF-inputs mismatch.
- (c) Genesis-hash mismatch.
- (d) Cross-chain-lift (entry for wrong tenant).
- (e) Format-version mismatch in entry.
- (f) Chain-link broken (`prev_hash` does not match prior entry's `payload_hash`).
- (g) Unknown `key_version` (IKM not found).
- (h) `key_fingerprint` mismatch.
- (i) MAC mismatch (`payload_hash` does not match recomputed HMAC).

### §A9.2 Decision tree

1. **Expected Scenario 9 negative-test entry?** If yes (the entry is documented in the institution's test matrix and test-result log), the failure is intentional. Not a control failure.
2. **Cross-tenant-lift (type d)?** If yes, investigate: SDK misconfigured the tenant, or cross-tenant data copy occurred. Escalate to IR Scenario 4. If tenant matches, failure is a different type.
3. **Key-version or fingerprint issue (type g or h)?** Cross-reference IKM rotation log. If a documented rotation event coincides with the entry's `captured_at`, audit traces the rotation procedure. If `key_version` is stale or fingerprint is wrong, this is a procedural gap (rotation executed but registry not updated). Not a crypto defect; a process gap.
4. **MAC mismatch (type i) with correct key and fingerprint?** The entry's `payload_hash` does not match the recomputed HMAC. Possibilities: entry bytes modified after MAC computation (tampering — escalate to IR Scenario 1 or 2); canonical-form field modified in transit; IKM mismatch between SDK and ledger sides. Escalate to forensic team.

### §A9.3 Audit conclusion

After root-cause analysis, the failure is one of:

- **Expected.** Documented in the test matrix; not a control deficiency.
- **Procedural.** Operational root cause (rotation not recorded, entry misconfigured, transformation at the collector). Issue MRA against operations to update procedures.
- **Design.** Cryptographic or architectural root cause (spec logic not implemented correctly). Issue MRA against development/architecture team and may restrict further AI deployment pending fix.
- **Incident.** Security compromise. Escalate to IR and issue MRIA (Matters Requiring Immediate Attention) to the audit committee.

---

## §A10. Continuous-auditing queries and ADA integration

IIA's Technology Audit Guidance encourages audit to move from annual testing to continuous monitoring where feasible. The chain's aggregated ledger exposes query APIs enabling audit to run continuous or near-real-time queries.

### §A10.1 Continuous-auditing objectives

- **Real-time anomaly detection.** Alert audit if verifier failure rate exceeds the threshold, if a seal delay exceeds 60 minutes, or if a fingerprint mismatch occurs.
- **Near-real-time forensic readiness.** When an incident is reported, audit can immediately query the chain to identify the time window and the entries affected.
- **Preventive escalation.** Flag operational trends (increasing seal delays, rising fingerprint mismatch frequency) before they become control failures.

### §A10.2 Query templates

**Daily seal-age query.**

```sql
SELECT seal_date, signed_at,
  EXTRACT(EPOCH FROM (signed_at - CAST(seal_date AS TIMESTAMP))) AS seal_age_seconds
FROM daily_seals
WHERE signed_at > NOW() - INTERVAL '7 days'
ORDER BY seal_age_seconds DESC;
```

Audit monitors: any `seal_age_seconds > 3600` (60 minutes) is flagged. Threshold alert: more than 2 seals per week exceeding 60 minutes escalates to management.

**Verifier-failure query.**

```sql
SELECT DATE(captured_at) AS day,
  COUNT(*) AS entries_count,
  COUNT(CASE WHEN verification_result = 'FAIL' THEN 1 END) AS fail_count
FROM chain_entries_verified
WHERE captured_at > NOW() - INTERVAL '7 days'
GROUP BY DATE(captured_at)
ORDER BY fail_count DESC;
```

Audit monitors: any day with `fail_count` greater than 0.1% of `entries_count` is flagged.

**Key-fingerprint anomaly query.**

```sql
SELECT tenant_id, key_version, COUNT(DISTINCT key_fingerprint) AS distinct_fingerprints
FROM chain_entries
WHERE key_fingerprint IS NOT NULL
GROUP BY tenant_id, key_version
HAVING COUNT(DISTINCT key_fingerprint) > 1;
```

Audit monitors: any tenant-key_version pair with multiple distinct fingerprints indicates potential key-rotation errors or cross-tenant contamination.

**Incident-window forensic query.**

```sql
SELECT tenant_id, run_id, seq, captured_at, payload_hash
FROM chain_entries
WHERE captured_at BETWEEN [start] AND [end]
ORDER BY seq;
```

Used when an incident is reported. Audit can independently verify these entries or compare to a backup copy.

### §A10.3 Integration with the institution's monitoring stack

Audit integrates queries into the institution's audit platform (KPMG Audit Analytics, Deloitte LucidEm, ACL, Alteryx, custom Tableau or Python). Daily jobs run the seal-age query and load results into the SIEM; threshold alerts trigger if exceeded; dashboards show trends over time.

---

## §A11. Workpapers and QAR-defensibility

IIA 2600 requires every five years an external Quality Assurance Review (QAR) by a peer audit firm. The QAR team will assess sampling design, independence, evidence sufficiency, findings rigor, cross-functional coordination, and prior-year remediation re-test. Workpapers must be contemporaneous and complete to defend the work.

### §A11.1 Workpaper contents per chain-audit engagement

1. **Procedure documentation.** The written procedure (per `docs/audit-procedures.md` and this evidence pack) with engagement-period dates, objectives, and steps.
2. **Population and sampling.** Population definition, sampling methodology, sample size, random-number seed.
3. **Execution evidence.** Per sampled item: item ID (`tenant_id`, `run_id`, `seq`), test performed (verifier run), result (PASS or FAIL), notes.
4. **Root-cause analysis.** Investigation notes for any failures, per §A9.
5. **Independence confirmation.** Note that the audit sample was selected without knowledge of the operations team and the sample was not communicated until verifier results were obtained.
6. **Sufficiency assessment.** Conclusion statement supporting the sample's sufficiency under IIA 2310.

### §A11.2 Continuous-auditing documentation

Workpapers also include the query logic and frequency, threshold settings, log of queries run and results over the period, alerts triggered and the investigation/remediation that followed.

### §A11.3 QAR-defensibility checklist

Before the QAR engagement, audit confirms workpapers answer:

- Was the sample size adequate? (Sample-size calculation and rationale.)
- Was the sample selected appropriately? (Random-number seed, selection method, evidence operations did not interfere.)
- Were failures investigated thoroughly? (Root-cause-analysis workpapers.)
- Were independence requirements met? (Query-role provisioning evidence and sample-communication timeline.)
- Was prior-year remediation re-tested? (§A13.)

Workpapers are contemporaneous, not reconstructed after the fact. Key decisions (sample size, threshold settings, root-cause conclusions) are documented with the reasoning at the time of decision.

---

## §A12. Audit-committee reporting and the annual opinion

The audit-committee-summary names quarterly reporting of verifier-run results, anomalies, seal-age metrics, and master-key reconciliation status. This pack defines the report template.

### §A12.1 Quarterly chain-audit report (3-5 pages plus exhibits)

**Status summary (one paragraph).** "Internal Audit completed testing of the chain-of-custody control for [date range]. Verifier validation was performed on [N] sampled AI decisions across [date range]. All sampled entries passed integrity verification. The 60-minute seal-publication SLA was met for [X] of [Y] days (compliance rate [%]). [N] operational anomalies were investigated and resolved. Overall control status: OPERATING EFFECTIVELY."

**Verifier-testing results.** Chart showing verifier-pass rate by month (target 100%); table naming any failures (date, entry ID, failure type, root cause, remediation).

**Seal-age compliance.** Histogram or line chart of seal-age distribution (target 100% within 60 minutes); narrative explaining any delays.

**Key-rotation reconciliation.** Confirmation that monthly reconciliation occurred per procedure; any fingerprint mismatches and remediation.

**Incident-response playbook status.** Confirmation that the playbook has been maintained and either a real incident was handled per the playbook or a tabletop exercise was conducted.

**Audit findings and remediation tracking.** Table of any deficiencies identified, severity, management response, remediation status.

**Third-line independence statement.** "Audit maintained independence throughout the testing period. Audit's query-role access to the chain ledger is provisioned by the DBA. No management interference in sample selection or evidence gathering was encountered."

### §A12.2 Annual audit opinion language

Included in the CAE's annual internal-audit opinion to the audit committee:

> We assessed the operating effectiveness of the chain-of-custody control for AI decision-capture integrity during the audit period. We obtained sufficient and reliable evidence through verifier-based testing of a stratified random sample of [N] entries, key-rotation reconciliation testing, and seal-age SLA compliance testing. All sampled entries passed integrity verification, and the SLA was met for [X]% of days. We identified [number] deficiencies, [severity levels], all with management remediation plans tracking to [dates]. Based on our testing, we conclude the chain-of-custody control is operating effectively and provides integrity-bearing evidence of captured AI decisions. This opinion covers the integrity of decision capture; it does not cover the correctness of the AI decisions themselves, which are subject to model-risk-management controls.

### §A12.3 First-year education session

In the first year post-deployment, the CAE schedules a dedicated audit-committee education session (30-45 minutes outside the regular reporting cadence): walkthrough of the verifier procedure, explanation of the cryptographic foundation (HMAC, Merkle tree, HSM signing), sample audit-test results, Q&A.

---

## §A13. Repeat-finding and remediation-effectiveness testing

When audit identifies a deficiency in year 1, management responds with a remediation plan. In year 2, audit re-tests to confirm remediation. The spec does not address how audit tracks repeat findings; this section closes that gap.

### §A13.1 Remediation-effectiveness testing procedure

1. **Timing.** At the start of each annual engagement, the audit team identifies any prior-year chain-related findings.
2. **Remediation-plan review.** Audit pulls the management response from the prior-year audit file: deficiency description, root cause, corrective action, responsible party, target remediation date, residual-risk acknowledgment.
3. **Current-year re-test.** A test procedure designed to confirm remediation. For documentation gaps, audit confirms the procedure is documented and approved and has been followed (last three rotations match the documented schedule). For operational gaps, audit re-tests the metric (seal-age SLA, fingerprint reconciliation cadence). For control-design gaps, audit re-confirms access controls and that no new overlaps have emerged.
4. **Severity assessment.** If remediation is confirmed complete: close with a note. If partially complete or ineffective: determine whether this is a repeat finding or a new finding. If not executed: escalate from deficiency to significant deficiency or material weakness.

### §A13.2 Repeat-finding escalation

- **Same finding repeats in two consecutive audits.** CAE escalates to CFO and CRO in writing requesting a revised remediation plan with specific accountability and timeline.
- **Same finding in a third consecutive audit.** Escalates to the audit committee. Triggers a broader assessment of management's commitment to the control.

### §A13.3 Tracking document

Audit maintains a separate tracking document for prior-year findings and remediation status. Reviewed at the start of each engagement; status reported to the audit committee in the annual opinion.

---

## §A14. Cross-functional audit coordination

The chain involves three audit specialties: cryptography/security, IT general controls, business/AI model-risk. The CAE assigns roles and coordinates findings.

### §A14.1 Role assignments

- **Cryptography/IT security auditor.** Confirms HMAC computation, Merkle-seal construction, HSM signing per spec. May engage external crypto expert if in-house team lacks expertise.
- **IT general-controls auditor.** Confirms ledger database, seal-job automation, access-control roles, backup/recovery procedures.
- **Model-risk auditor.** Confirms scope of AI decisions logged, business-rules logic determining which decisions are logged, model-version tracking.
- **Business-process auditor.** Confirms incident-response procedures are understood by first-line operations and operationally feasible.

### §A14.2 Coordination touchpoints

- **Pre-engagement kickoff.** Team aligns on scope, objectives, evidence-gathering plan. CAE confirms which decisions the chain captures, which teams own each component, and the primary contact for each question.
- **Evidence-gathering.** Each team member gathers per specialty. Procedures P-1 through P-17 (and the additions) are assigned across the team to avoid duplication.
- **Findings consolidation.** Team meets to discuss findings and synthesize. CAE classifies findings as design issues (affects control's inherent design) or operating issues (affects execution).

### §A14.3 External-expert engagement

For first-year deployments, external crypto and IT expertise is typical. CAE defines scope narrowly (e.g., "validate the verifier binary and HSM signing procedure; do not conduct a full cryptographic implementation audit"). The external expert's findings are incorporated into workpapers with clear attribution.

---

## §A15. Whistleblower-investigation procedure (P-WB1)

When an employee whistleblower or external regulator alleges misconduct related to an AI decision, audit may investigate. Chain evidence supports the investigation by confirming decision time, decision-bytes integrity, and whether related decisions were deleted from the ledger.

### §A15.1 Procedure

1. **Trigger.** A whistleblower allegation or regulatory inquiry names a specific AI decision (by tenant, run_id, approximate timestamp, or customer ID).
2. **Scope definition.** CAE defines the evidence scope in consultation with General Counsel: date range, tenant, business context.
3. **Evidence preservation.** Audit immediately locks the relevant ledger segment. Audit sends a preservation notice to the DBA directing no bulk deletions until the investigation concludes. Audit creates a read-only snapshot of the ledger for the relevant range, stored in audit's evidence locker with audit-trail logging.
4. **Independent evidence extraction.** Audit (NOT the accused party or their subordinates) queries the ledger snapshot using the audit-query role.
5. **Verifier-based authenticity check.** For each extracted entry, audit runs the verifier per spec §7. Failures are documented as evidence of possible tampering.
6. **Forensic analysis.** If the allegation is that a decision was altered, audit compares the extracted entry to a backup or prior copy.
7. **Expert input.** For complex allegations, audit may engage the model-risk team or an external ML expert to assess whether the model version on the alleged date matches the documented version.

### §A15.2 Confidentiality

The investigation file is highly confidential. Access limited to: CAE, investigation lead, General Counsel, Audit Committee Chair (when briefing). Access logs to the ledger snapshot are maintained and reviewed.

### §A15.3 Litigation readiness

Upon completion, the CAE determines: will this matter result in litigation? If yes, evidence is transferred to General Counsel's litigation hold. The chain's tamper-evidence property (cryptographic verification) is admissible as a business record under FRE 803(6) and as a digital signature under UCITA or state equivalents. Audit retains workpapers per the institution's retention policy (typically 7 years for fraud investigations).

---

## §A16. Severity rating thresholds for chain-related findings

Findings are rated using the PCAOB/AICPA hierarchy.

### §A16.1 Deficiency (low-risk, no committee escalation)

- Single verifier failure with documented isolated cause; remediated within 48 hours.
- Single seal-age delay 1-60 minutes beyond the SLA with documented cause.
- Single fingerprint mismatch from documented key-rotation procedural gap; corrected within hours.
- Isolated configuration issue (e.g., verifier output log not retained per spec §10.9).

### §A16.2 Significant deficiency (medium-risk, escalate to committee within 90 days)

- Verifier failure rate 0.1%-1% of sampled entries without documented root cause.
- Multiple seals (3+) per month exceed the 60-minute SLA, indicating systematic issue.
- Fingerprint mismatches on multiple entries from the same `key_version`.
- Audit-query role not properly segregated from operations role.
- Incident-response playbook exists but not exercised or tested during the period.
- Multiple days of missing seal records.
- Chain ledger not accessible to audit via query API; audit must request data from operations team.

### §A16.3 Material weakness (high-risk, escalate to audit committee immediately)

- Verifier failure rate above 1% of sampled entries (systematic).
- Evidence of tampering (entry's bytes were modified in the ledger after ingestion; discovered by re-running the verifier on prior result vs current ledger state).
- IKM compromise or suspected compromise (HSM audit trail shows unauthorized access or extraction attempt).
- Chain infrastructure has no backup or recovery procedure (single failure results in permanent loss).
- More than 5 consecutive days of missing seal records.
- Institution has been unable to produce the chain's aggregated ledger despite audit requests.

### §A16.4 Escalation timing

- Deficiencies tracked by management; monitored by audit in follow-up.
- Significant deficiencies reported to the audit committee in the quarterly read-out; 90-day remediation.
- Material weaknesses reported within 5 business days of discovery; 1-week immediate action plan; included in any SOX 404 ICFR disclosure if AI is in ICFR scope.

---

## §A17. Chain-derived KRIs for the audit committee

Audit reports the following KRIs to the audit committee quarterly:

| KRI | Source | Threshold |
|---|---|---|
| Verifier pass rate | Daily verifier runs on sampled entries | Target 100%; alert below 99.9% |
| Seal-age SLA compliance | Daily seal records | Target 100% within 60 minutes; alert at any seal above 80 minutes |
| Fingerprint mismatch rate | Reconciliation log | Target 0%; alert above 0.1% |
| IKM custody anomalies | HSM audit log | Target 0; alert at 1 or more per quarter |
| Incident-response time (vendor-hosted) | Vendor escalation log | Target under 1 hour; alert above 4 hours |
| Vendor SOC report currency | Vendor evidence | Target current within 12 months; alert if expired |
| First-line + second-line monthly review | Meeting minutes | Target attendance per RACI; alert if missed |

The KRIs feed the audit-committee dashboard. The committee reviews quarterly and directs management actions when thresholds are breached.

---

## §A18. Three-lines governance for the chain

The chain has roles distributed across the three lines.

### §A18.1 First line (AI Operations / Chain-Operations Team)

Executes the seal job daily within the 60-minute SLA. Maintains the HSM in operational state. Performs IKM rotation per the documented cadence. Logs all operational events. Detects and initially responds to operational anomalies. Escalates any suspected security incidents to the IR coordinator within 30 minutes. Provides daily and weekly operational metrics to the second line.

### §A18.2 Second line (IT Risk / Compliance / MRM)

Reviews first-line operational metrics monthly. Assesses whether observed anomalies indicate a control-design issue or an operational incident. Escalates suspected material incidents to the CRO and General Counsel within 4 hours. Provides a monthly summary to the audit committee of chain-related risks.

### §A18.3 Third line (Internal Audit)

Plans and executes the chain-control-testing procedures per the annual audit plan. Tests design and operating effectiveness. Verifies first-line operational documentation and second-line risk assessment. Investigates control failures and assesses severity per §A16. Provides the annual audit opinion. Reports findings to the audit committee quarterly.

### §A18.4 Role-separation controls

The second line's risk assessment is independent of the first line's IR documentation. The third line's audit findings do not supersede first or second-line operational decisions; audit's role is to test and report; management decides remediation. The three lines report to different organizational leaders: first line to CIO; second line to CRO; third line to the audit committee functionally and CEO administratively.

### §A18.5 Governance touchpoints

- Monthly chain-operations review: first line presents; second line assesses trends; third line listens and notes audit-focus areas.
- Annual control-description review: first line documents; second line approves; third line reviews in PIR or annual engagement.
- Incident-response playbook exercise: first and second line conduct the tabletop; third line may attend as observer.
- Quarterly audit-committee reporting: second line presents risk metrics; third line presents audit findings and opinions; committee makes governance decisions.

---

## §A19. Stopping criterion

This evidence pack closes the audit-governance layer above the spec. Internal audit can plan engagements (§A2, §A3), maintain independence (§A4), identify and assess controls (§A5), conduct PIRs (§A6), align with SOX 404 ICFR (§A7), sample defensibly (§A8), root-cause failures rigorously (§A9), run continuous queries (§A10), document workpapers for QAR-defensibility (§A11), report to the audit committee (§A12), re-test prior-year remediation (§A13), coordinate cross-functionally (§A14), use chain evidence in whistleblower investigations (§A15), rate findings consistently (§A16), feed quarterly KRIs to the committee (§A17), and operate within the three-lines framework (§A18).

The chain-of-custody specification is audit-ready. With this pack, the third line can apply IIA Standards to chain evidence with confidence, issue an opinion the QAR will defend, and integrate cleanly with external audit on SOX 404 ICFR engagements.
