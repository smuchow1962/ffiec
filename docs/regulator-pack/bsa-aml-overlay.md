---
status: informative
alignment-reference: Bank Secrecy Act (31 USC §5311 et seq.); USA PATRIOT Act §314(a) and §314(b); 31 CFR Chapter X (FinCEN regulations); 31 CFR Part 1010 (general); 31 CFR Part 1020 (banks); OFAC sanctions and SDN screening (31 CFR Part 501); FinCEN AML Program Effectiveness Rule (2024); FinCEN AML Program Modernization Guidance (2023); FFIEC BSA/AML Examination Manual (2021 ed.); AML Act of 2020 (AMLA); Beneficial Ownership Information (BOI) reporting (31 CFR §1010.380); 2019 SAR Narrative Rule amendments; SR 11-7 (model risk management); OCC Bulletin 2021-53 (third-party model risk management); ECOA / Reg B adverse-action notice rules
companion-docs:
  - docs/audit-procedures.md (institution-side P-1..P-N procedures; §C extends with AML-specific procedures)
  - docs/regulator-pack/fdic-occ-examination-overlay.md (FDIC/OCC examination workflow)
  - docs/regulator-pack/examination-response-workflow.md (general PBC list)
  - docs/incident-response-playbook.md (Scenario 4 master-key compromise notification)
  - spec/chain-of-custody-v1.md §4.4 (audit.* event namespace; routing.* events)
  - design/12-aml-event-schema.md (AML event-schema reference, planned)
date: 2026-05-07
version: 1.0.0
---

# BSA / AML Overlay — Chain-of-Custody for AML Decision Provenance

> **What this doc is.** Operational overlay for the FFIEC chain-of-custody control when the chain is used by AML compliance programs to capture transaction-monitoring (TM) alerts, sanctions-screening dispositions, SAR decisions, KYC/CDD enrichment, beneficial-ownership inferences, and lookback-investigation findings. The overlay normalizes the AML-specific event schema, audit procedures, FinCEN/OCC examiner expectations, look-back-order response procedures, and customer-dispute disclosure shape so AML compliance teams can use the chain as the primary evidence source for FinCEN, OCC, FDIC, and state examinations.
>
> **Audience.** Chief BSA/AML Officer, Chief Compliance Officer, AML compliance program staff, FinCEN examiners, OCC/FDIC AML specialists, Compliance Committee members, AML model governance teams (SR 11-7).
>
> **Why this overlay.** The chain captures and integrity-binds whatever event the institution emits. AML compliance teams need a normative event taxonomy so events are consistent across institutions, examiner-interpretable, and ready for FinCEN look-back response. Without normalization, each institution home-grows its own AML audit-event taxonomy, and examiners see incompatible evidence trails across institutions. This overlay folds the AML-specific schema and procedures into v1.0a as a companion document — the spec body remains unchanged.

## How to use this overlay

1. **AML event schema (§A).** Use the schema as the normative event taxonomy for AML decisions. Compliance teams emit events with the named fields; the chain integrity-binds whatever is emitted.
2. **AML-specific audit procedures (§B).** Extend the institution's audit procedure document with the AML-N procedures.
3. **FinCEN look-back response (§C).** Use the response procedure when responding to a FinCEN look-back order or RFI.
4. **Customer-dispute response (§D).** Use the customer-disclosure shape when a customer disputes a TM block or SAR decision.
5. **Examiner orientation (§E).** Use the examiner-orientation reference for FinCEN, OCC, and FDIC AML specialists.
6. **Privacy-by-design composition (§F).** Compose with the privacy-by-design tokenization story so customer-identity tracking remains consistent across regions and over time.
7. **Model-risk governance (§G).** Connect the chain's AML events to the institution's SR 11-7 model-risk governance for AML models.

---

## §A. AML event schema — normative taxonomy

The events below are emitted under the spec's `audit.*` namespace per spec §4.4. Each event carries the standard chain attributes (event_id, run_id, seq, captured_at, tenant_id, key_fingerprint, MAC) plus the AML-specific fields named below.

### §A.1 Transaction-monitoring (TM) events

#### `audit.tm.alert_generated`

Emitted by the TM system at the moment an alert fires.

| Field | Type | Description |
|---|---|---|
| `alert_id` | string | Unique identifier in the bank's TM ledger |
| `transaction_id` | string | The transaction that triggered the alert |
| `model_name` | string | The TM model or rules-engine identifier |
| `model_version` | string | Version tag for reproducibility |
| `score_or_rule_match` | numeric or string | Numeric score (ML) or rule name (rules-based) |
| `flags` | array of strings | Detected factors: `velocity, geography, customer_typology, amount_threshold, network_pattern, time_pattern` |
| `alert_severity` | enum | `high, medium, low` |
| `explanation_factors` | array of objects | Top 3–5 SHAP/LIME contribution factors when ML; `{factor, contribution}` |

#### `audit.tm.alert_suppressed`

Emitted when the TM system suppresses a potential alert. Captures false-negative provenance.

| Field | Type | Description |
|---|---|---|
| `would_be_alert_id` | string | Hypothetical ID if the alert were emitted |
| `transaction_id` | string | The transaction the model evaluated |
| `model_score` | numeric | What score the TM system computed |
| `suppression_reason` | enum | `below_threshold, deduplication_recent, policy_exception, cooldown, whitelist_corridor` |
| `policy_or_threshold_ref` | string | Reference to the policy or threshold that drove suppression |
| `emitted_at` | timestamp | When the suppression decision was made |

Volume note: institutions emit this for every suppression. Storage cost is non-trivial; the institution's CC8.1 may name a sampling posture for low-confidence suppressions (e.g., "below threshold by > 50%"), with the SOC team confirming the sampling is documented and not discriminatory.

#### `audit.tm.alert_triage`

Emitted when Compliance triages a TM alert.

| Field | Type | Description |
|---|---|---|
| `alert_id` | string | Cross-references the `audit.tm.alert_generated` event |
| `triage_decision` | enum | `true_positive, false_positive, policy_exception` |
| `triage_reason` | string | Free-text rationale |
| `triaged_by` | string | Compliance officer ID |
| `escalated_to_sar` | boolean | Whether escalated to SAR-decision workflow |
| `triaged_at` | timestamp | UTC timestamp |

#### `audit.tm.threshold_tuned`

Emitted when Compliance makes a deliberate threshold change.

| Field | Type | Description |
|---|---|---|
| `model_name_or_rule_name` | string | The TM model or rule affected |
| `parameter_name` | string | E.g., `velocity_threshold_usd_per_day_wire` |
| `old_value` | numeric or string | Pre-change value |
| `new_value` | numeric or string | Post-change value |
| `justification` | string | Free-text rationale |
| `approved_by` | string | Compliance officer who approved |
| `tested_in_sandbox` | boolean | Whether tested in non-production first |
| `test_results_summary` | string | Pre-change vs post-change false-positive and true-positive rates |
| `deployed_at` | timestamp | When the change went live |
| `expected_impact` | string | Projected change in alert volume |
| `jurisdiction_scope` | enum | `US, EU, APAC, global` (per regional regime) |

### §A.2 OFAC and sanctions events

#### `audit.ofac.screening_result`

Emitted when the bank screens a customer or transaction against sanctions lists.

| Field | Type | Description |
|---|---|---|
| `screening_id` | string | Unique ID for the screening run |
| `customer_id` | string | Tokenized per privacy-by-design (§F) |
| `screening_list` | enum | `SDN, SDGT, DPL, EU_consolidated, UN_1267, OFAC_sectoral, FATF_high_risk` |
| `hits_found` | integer | Number of potential name matches |
| `match_scores` | array of numeric | Similarity scores for each hit |
| `disposition` | enum | `true_hit, false_positive, waived_per_policy` |
| `waiver_authority` | string | If waived: e.g., "Treasury OFAC General License GL-101" or "internal policy exception approved by CCO" |
| `additional_names_screened` | array of strings | Related names (BO, guarantor, authorized users) |
| `resolved_by` | string | Compliance officer ID |
| `resolved_at` | timestamp | UTC timestamp |

#### `audit.ofac.advisory_received`

Emitted when FinCEN or another regulator issues an OFAC advisory.

| Field | Type | Description |
|---|---|---|
| `advisory_id` | string | E.g., "FIN-2024-A101" |
| `advisory_date` | date | Issue date |
| `advisory_subject` | string | Brief subject line |
| `institution_response_required` | boolean | Whether response is required |
| `response_due_date` | date | Required response deadline |

#### `audit.ofac.advisory_response_completed`

Emitted when the institution responds to an OFAC advisory.

| Field | Type | Description |
|---|---|---|
| `advisory_id` | string | Cross-references the advisory |
| `response_taken` | enum | `enhanced_screening, customer_update, threshold_adjustment, policy_clarification, no_action_required` |
| `response_details` | string | Free-text description |
| `completed_at` | timestamp | UTC timestamp |

### §A.3 SAR and CTR events

#### `audit.aml.sar_decision`

Emitted when Compliance commits to filing or not filing a SAR.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `decision` | enum | `will_file, will_not_file, undecided_pending_investigation` |
| `decision_rationale` | string | Free-text rationale |
| `supporting_alert_ids` | array of strings | TM alerts or investigation activities that drove the decision |
| `filed_by_officer` | string | Compliance officer ID |
| `decided_at` | timestamp | UTC timestamp |

#### `audit.aml.sar_filed`

Emitted when the SAR is submitted to FinCEN.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `sar_id_internal` | string | Bank's internal SAR ID |
| `sar_id_fincen` | string | FinCEN acknowledgment ID |
| `filed_at` | timestamp | When submitted to FinCEN |
| `decision_event_id` | string | Cross-references the `audit.aml.sar_decision` event |

The gap between `sar_decision.decided_at` and `sar_filed.filed_at` is the investigation window; examiners audit whether timing was appropriate.

#### `audit.aml.sar_narrative_drafted`

Emitted when the SAR narrative is drafted (with or without AI assistance).

| Field | Type | Description |
|---|---|---|
| `sar_id` | string | Internal SAR ID |
| `drafted_by` | enum | `human_only, ai_assisted, ai_generated_human_reviewed` |
| `model_if_ai` | string | If AI-assisted: e.g., `gpt-4-turbo`, `claude-sonnet` |
| `model_version` | string | Version tag for reproducibility |
| `human_reviewer_id` | string | Reviewing Compliance officer |
| `reviewed_at` | timestamp | UTC timestamp |
| `facts_verified_by` | string | Source of facts: TM system, manual investigation, SAR case file |
| `narrative_integrity_check` | string | Hallucination scan, factual accuracy verification, FinCEN guidance alignment check |

#### `audit.aml.ctr_filed`

Emitted when the bank files a Currency Transaction Report.

| Field | Type | Description |
|---|---|---|
| `transaction_id` | string | Cash-transaction ID |
| `customer_id` | string | Tokenized |
| `amount_usd` | numeric | USD amount |
| `filing_reason` | enum | `meets_reporting_threshold_10k, meets_reporting_threshold_aggregated, filed_at_customer_request` |
| `filed_by_system` | string | Automated system identifier or "manual" |
| `filed_at` | timestamp | UTC timestamp |
| `ctr_id_fincen` | string | FinCEN acknowledgment ID |
| `automation_model_name_if_any` | string | Model identifier if filed by automation |
| `automation_model_version` | string | Version tag |

### §A.4 KYC / CDD and beneficial-ownership events

#### `audit.kyc.customer_risk_score_computed`

Emitted when an ML model computes a customer risk score.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `model_name` | string | E.g., "customer-risk-scorer-v2" |
| `model_version` | string | Build hash or version tag |
| `score` | numeric | Numeric risk score |
| `risk_category` | enum | `low, medium, high, critical` |
| `primary_factors` | array of strings | Top 3–5 risk factors |
| `computed_at` | timestamp | UTC timestamp |
| `inputs_used` | array of strings | Data sources: `transaction_history_12m, network_analysis, sanctions_screening_history` |

#### `audit.kyc.customer_tier_assigned`

Emitted when Compliance assigns a customer to a risk tier.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `tier` | enum | `low, medium, high, critical` |
| `assigned_on_date` | date | Effective date of tier assignment |
| `assigned_by` | string | Compliance officer ID |
| `rationale` | string | Free-text rationale referencing the risk score event |
| `risk_score_event_id` | string | Cross-references `audit.kyc.customer_risk_score_computed` |

#### `audit.kyc.beneficial_owner_inferred`

Emitted when an ML model or data-enrichment service infers beneficial ownership.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `legal_entity_name` | string | The legal-entity customer name |
| `inferred_bo_id` | string | Tokenized identifier for the inferred owner |
| `confidence_score` | numeric | Model confidence (0.0–1.0) |
| `inference_method` | enum | `ml_model, public_filings_match, transaction_flow_analysis, third_party_data_provider` |
| `supporting_evidence` | string | Summary of evidence |
| `inferred_at` | timestamp | UTC timestamp |
| `confirmed_by_customer_disclosure` | boolean | Whether customer subsequently confirmed |

### §A.5 Lookback and remediation events

#### `audit.aml.lookback_initiated`

Emitted when Compliance starts a lookback investigation.

| Field | Type | Description |
|---|---|---|
| `lookback_id` | string | Unique ID for the lookback |
| `cohort_description` | string | E.g., "all SARs filed Jan–Mar 2026" |
| `initiated_by` | string | Compliance officer ID |
| `initiated_at` | timestamp | UTC timestamp |

#### `audit.aml.lookback_finding`

Emitted when a finding is discovered.

| Field | Type | Description |
|---|---|---|
| `lookback_id` | string | Cross-references the lookback |
| `finding_id` | string | Unique ID for this finding |
| `finding_type` | enum | `inappropriate_sar_filed, sar_should_have_been_filed_false_negative, threshold_drift, model_drift_detected` |
| `alert_or_sar_id_affected` | string | Cross-references the affected alert or SAR |
| `description` | string | Free-text description |
| `discovered_at` | timestamp | UTC timestamp |

#### `audit.aml.lookback_remediation`

Emitted when Compliance commits to remediation.

| Field | Type | Description |
|---|---|---|
| `lookback_id` | string | Cross-references the lookback |
| `finding_id` | string | Cross-references the finding |
| `remediation_action` | enum | `model_retrain_with_new_data, threshold_adjustment, policy_exception_removed, enhanced_manual_review, sar_filed_late` |
| `action_details` | string | Free-text details and metrics |
| `action_start_date` | date | When remediation began |
| `action_completion_date` | date | When remediation completed |
| `responsible_officer` | string | Compliance officer ID |

### §A.6 Model-risk governance events (SR 11-7 alignment)

#### `audit.model_risk.validation_test_completed`

Emitted when an AML model is validated.

| Field | Type | Description |
|---|---|---|
| `model_name` | string | Model identifier |
| `test_date` | date | Date of validation |
| `test_type` | enum | `annual_backtesting, fairness_audit, false_positive_analysis, threshold_sensitivity_analysis, drift_detection` |
| `test_results_summary` | string | Free-text summary of results |
| `passed_validation` | boolean | Whether the model passed |
| `findings_and_remediation` | string | If failed: what remediation was taken |
| `validator_id` | string | Internal audit or external validator ID |

### §A.7 Performance-monitoring events

#### `audit.aml.performance_report_generated`

Emitted monthly by Compliance.

| Field | Type | Description |
|---|---|---|
| `reporting_month` | date | Month covered |
| `model_name` | string | TM model identifier |
| `alerts_generated` | integer | Count of alerts in the month |
| `alerts_triaged_true_positive` | integer | True-positive count |
| `alerts_triaged_false_positive` | integer | False-positive count |
| `false_positive_rate` | numeric | Percentage |
| `sars_filed` | integer | SAR count in the month |
| `lookback_findings` | integer | Missed-alert count from prior-month lookback |
| `model_threshold_changes_this_month` | integer | Count of `audit.tm.threshold_tuned` events |
| `threshold_rationale_summary` | string | Summary of why changes were made |
| `reported_by` | string | Compliance officer ID |
| `reported_at` | timestamp | UTC timestamp |

### §A.8 MOU and consent-order events

#### `audit.aml.remediation_milestone_tracked`

Emitted when a milestone under an MOU or consent order is tracked.

| Field | Type | Description |
|---|---|---|
| `mou_or_consent_order_id` | string | Reference to the MOU or consent order |
| `milestone_description` | string | E.g., "implement TM threshold tuning with documented governance" |
| `target_completion_date` | date | Required deadline |
| `actual_completion_date` | date | When completed |
| `completion_evidence` | string | References to chain events: `"deployed model v3 on 2026-09-15, see audit.tm.threshold_tuned entries for Sept 2026"` |
| `verified_by` | string | Internal audit or external examiner who confirmed |

### §A.9 Adverse-action events

#### `audit.aml.adverse_action_notice_issued`

Emitted when AML factors drive a credit denial.

| Field | Type | Description |
|---|---|---|
| `customer_id` | string | Tokenized |
| `adverse_action_reason` | string | E.g., "failed_aml_customer_risk_assessment" or "transactions_inconsistent_with_stated_business" |
| `supporting_aml_decision_ids` | array of strings | Cross-references to TM alert, triage, or KYC risk score |
| `notice_issued_at` | timestamp | UTC timestamp |
| `issued_to_customer` | date | When customer received the notice |

---

## §B. AML-specific audit procedures (overlay onto `audit-procedures.md`)

The institution's audit procedure document includes the following AML-specific procedures. See `docs/audit-procedures.md` §AML procedures for the canonical list (P-41 through P-48 below).

### §B.1 AML event-completeness sample

For each `audit.tm.alert_generated` in the period, confirm the corresponding `audit.tm.alert_triage` exists within the institution's documented triage SLA. Missing triage events are control-completeness gaps; the SOC team escalates.

### §B.2 Suppression-volume baseline

The institution maintains a documented rolling 90-day baseline of `audit.tm.alert_suppressed` volume per model and per jurisdiction. The SOC team confirms the baseline is documented and that suppression-rate excursions trigger the institution's documented review procedure.

### §B.3 OFAC waiver authority verification

For each `audit.ofac.screening_result` with `disposition = waived_per_policy`, confirm the `waiver_authority` references either (a) a Treasury OFAC General License or (b) an internal policy exception with documented CCO approval. Waivers without documented authority are a significant deficiency.

### §B.4 SAR-decision-to-filing latency

For each `audit.aml.sar_decision` with `decision = will_file`, confirm a corresponding `audit.aml.sar_filed` exists within the institution's documented investigation window (typically 30–60 days). Outliers are reviewed for appropriateness.

### §B.5 SAR-narrative-AI provenance

For each `audit.aml.sar_narrative_drafted` with `drafted_by = ai_assisted` or `ai_generated_human_reviewed`, confirm `human_reviewer_id` and `narrative_integrity_check` are populated. Missing review attestation is a model-governance gap.

### §B.6 Lookback closure

For each `audit.aml.lookback_initiated`, confirm the corresponding `audit.aml.lookback_finding` and `audit.aml.lookback_remediation` events close out the lookback within the institution's documented window. Open lookbacks beyond the window are reviewed for cause.

### §B.7 Model-validation cadence

For each AML model in the institution's MRM inventory, confirm at least one `audit.model_risk.validation_test_completed` event per year. Missing annual validation is a SR 11-7 finding.

### §B.8 Performance-report continuity

For each month in the period, confirm an `audit.aml.performance_report_generated` event was emitted. Missing months are an effectiveness-rule (FinCEN 2024) gap.

### §B.9 Multi-region customer continuity

For multi-region institutions, sample 10 customers with cross-region activity. Verify the chain contains TM alerts from all regions where the customer was active. Confirm no region was silently dropped from monitoring. Missing regional coverage is a control-completeness gap.

---

## §C. FinCEN look-back order response procedure

When FinCEN issues a look-back order ("provide all records and decisions relating to customer X for the period Y–Z"), Compliance responds using the chain as the primary evidence source.

### §C.1 Scope the request

1. Identify all chain entries matching the customer_id (tokenized per §F) and the date range.
2. Query the ledger for `(customer_id_token, event_date >= Y, event_date <= Z)`.
3. Pull all event types in scope: `audit.tm.alert_*`, `audit.ofac.*`, `audit.aml.*`, `audit.kyc.*`.

### §C.2 Completeness assertion

Confirm all relevant AML events are in scope. If any events are outside the chain (e.g., TM alerts captured but triage decisions in a separate system not yet integrated), Compliance explicitly states what is in the chain and what is outside scope. This is the "load-bearing scope assertion" — examiners rely on it.

### §C.3 Verify the ledger span

1. Run the verifier on the requested period.
2. Produce a signed verifier report showing PASS or FAIL.
3. If FAIL, investigate (usually indicates a missing key or configuration mismatch).
4. If PASS, include the report in the response.

### §C.4 Produce the response package

The response includes:
1. **Verifier report (PASS).** Cover letter naming the verifier version, the period, the tenant_id, and the customer_id_token.
2. **Chain-entry exports.** CSV or JSON, one line per entry, with human-readable field labels. The export includes all event types in scope.
3. **Compliance Officer certification statement.**

   > These entries are complete and accurate copies from our chain-of-custody ledger, independently verified by [verifier version] on [date]. The chain's integrity proves no entries were inserted, deleted, or modified after capture. The redactions below are applied per [policy reference]; no material facts have been redacted.

4. **Redaction legend.** Privacy-by-design redactions are listed; the redaction policy is cited.

### §C.5 Audit trail of the response

Compliance logs the response preparation:
- Verifier was run.
- Certification was signed.
- Package was delivered to FinCEN.

The log uses the operational-events schema (e.g., `audit.aml.fincen_response_prepared`) so the response itself is integrity-bound.

### §C.6 Follow-up disposition

FinCEN may issue a follow-up RFI for additional periods or related customers. The institution repeats §C.1–§C.5 for each scope expansion.

---

## §D. Customer-dispute response — chain evidence production

When a customer disputes a TM block or SAR decision, Compliance produces evidence to the customer (and potentially to the customer's lawyer or a court).

### §D.1 Production format

#### §D.1.1 Verifier report

Run the standalone verifier on the ledger span covering the customer's activity and the dispute window. Produce the verifier's PASS/FAIL report with a cover letter:

> This report proves our AML system's records on this matter were not altered after capture.

#### §D.1.2 Chain-entry extracts

Extract the specific chain entries relevant to the customer (alert generated, triage decision, SAR decision) in human-readable form (canonical JSON with field labels). Include them in the dispute response package.

#### §D.1.3 Certification

The Compliance Officer signs:

> These entries are true and accurate copies from our ledger, verified by [verifier version] on [date].

#### §D.1.4 Redaction

Redact:
- Other customers' data.
- Enforcement actions and active investigations.
- Privileged legal-analysis memoranda.

Preserve the customer-facing findings (why the decision was made about this customer).

### §D.2 Disposition of customer challenges

The verifier report is the integrity anchor; the certified extracts are the readable evidence. The customer's lawyer cannot challenge "how do we know this chain entry is authentic?" because the verifier proves authenticity. Disputes shift to questions of decision-rationale (covered by `triage_reason`, `decision_rationale`, `explanation_factors`) rather than evidence-authenticity.

### §D.3 Composition with `customer-dispute-procedures.md`

This section composes with `docs/customer-dispute-procedures.md` (general customer-dispute response). The §D shape is the AML-specific overlay; the general procedure handles non-AML disputes.

---

## §E. Examiner orientation reference

FinCEN, OCC, and FDIC examiners reviewing AML compliance programs that use the chain consume the following reference material.

### §E.1 What the chain proves

The chain proves no alert or decision was inserted, deleted, or modified after capture. It does not prove:
- The model's underlying score was accurate (model-validation evidence is separate).
- The Compliance officer's triage decision was correct (decision-quality is separate).
- The SAR narrative is factually accurate (narrative-quality is separate).

The chain is the integrity foundation, not the truth foundation. Examiners assess decision quality separately.

### §E.2 What the chain enables

Examiners can:
- Sample TM alerts and cross-check to triage decisions.
- Sample SAR decisions and trace to TM alerts and investigation activities.
- Verify that the institution's threshold-tuning was deliberate and documented (not silent code pushes).
- Verify that lookback investigations identified missed alerts and were remediated.
- Verify that AI-assisted SAR narratives went through human review.
- Verify model-validation cadence per SR 11-7.

### §E.3 Sample-testing procedure for AML examination

1. **Sample TM alerts.** 100 alerts across the period (stratified by model, jurisdiction, customer-impact tier).
2. **Cross-check triage.** For each alert, confirm an `audit.tm.alert_triage` exists within the SLA.
3. **Sample SAR decisions.** 30 SARs across the period.
4. **Trace SAR decisions to alerts.** For each SAR, verify the `supporting_alert_ids` reference valid alerts.
5. **Sample OFAC waivers.** All waived hits in the period (typically census).
6. **Verify waiver authority.** Each waiver references documented authority.
7. **Sample lookback findings.** All lookback findings in the period.
8. **Verify remediation.** Each finding has a remediation event.
9. **Run the verifier.** PASS over the sampled cohort confirms no alteration.

### §E.4 Examiner-orientation training module

Pair this overlay with `docs/regulator-pack/examiner-training.md` for general chain orientation. AML-specific orientation adds:
- AML event schema walk-through (§A) — 30 minutes.
- AML audit procedures (§B) — 30 minutes.
- Lookback-response procedure (§C) — 20 minutes.
- Customer-dispute response (§D) — 20 minutes.
- Hands-on lab: run the verifier on a sample AML cohort and interpret the output — 1 hour.

### §E.5 Findings dispositions

| Finding | Disposition |
|---|---|
| Missing triage event for a TM alert | Control-completeness gap; severity Medium |
| Missing OFAC waiver authority | Significant deficiency |
| Missing SAR-narrative AI-review attestation | Model-governance gap; severity Medium |
| Missing model-validation event | SR 11-7 finding; severity High |
| Lookback finding without remediation event | Control-program-maturity finding; severity High |
| Verifier failure on AML sample | Integrity finding; severity High; route to `fdic-occ-examination-overlay.md` §5 |

---

## §F. Privacy-by-design composition — customer-identity tokenization

The chain never sees the original customer name or PII. Instead:

1. The privacy-store (per `docs/regulator-pack/privacy-by-design.md`) maintains a token-to-original mapping.
2. Compliance emits chain events with the `customer_id` token.
3. The token is stable within a tenant_id so "customer X is flagged as high-risk" can be tracked across a year of transactions.

### §F.1 Cross-region consistency

For multi-region institutions, the customer_id token is consistent across regions for the same customer. The privacy-store coordinates token allocation; the chain references the token.

### §F.2 Customer-dispute disclosure

When a customer disputes a decision, the institution decodes the token to the customer's identity (in the privacy-store, not in the chain) and produces the decoded extracts to the customer. The chain remains tokenized.

### §F.3 FinCEN look-back response

FinCEN look-back orders typically reference customers by name. The institution decodes the customer name to the token, queries the chain for the token, and produces the chain entries. The response cover letter cross-references customer name to token internally; the chain entries themselves remain tokenized in the response package.

### §F.4 Operational complexity

The privacy-by-design composition is operationally complex. The institution's CC8.1 names:
- The privacy-store technology (HSM-backed token store, vault, etc.).
- The token-allocation procedure.
- The token-to-original mapping retention.
- The customer-disclosure decoding procedure.

---

## §G. Model-risk governance composition (SR 11-7)

OCC Bulletin 2021-53 (Model Risk Management for Third-Party Models) and SR 11-7 require institutions to monitor, validate, and govern AI models. AML models are subject to these requirements.

### §G.1 Chain events as model-governance evidence

The chain captures:
- `audit.model_risk.validation_test_completed` — documented validation cadence.
- `audit.tm.threshold_tuned` — documented threshold changes with sandbox testing.
- `audit.aml.lookback_finding` and `audit.aml.lookback_remediation` — documented model-drift detection and corrective action.
- `audit.aml.performance_report_generated` — monthly model-performance metrics.

### §G.2 SR 11-7 examination procedure

The OCC's Model Risk Management examination consumes the chain events:
1. Pull all `audit.model_risk.validation_test_completed` events for AML models in the period.
2. Confirm at least annual validation per model.
3. Pull all threshold-tuning and remediation events.
4. Confirm threshold changes were tested in sandbox before deployment.
5. Confirm lookback-driven retraining was deployed and post-retraining validation was performed.

### §G.3 Effective-challenge composition

The chain's `gen_ai_parameters` and `explanation_factors` fields support effective challenge: the institution can audit a sample of AML decisions and confirm the model's reasoning was reasonable. P-25 in `docs/audit-procedures.md` covers the SR 11-7 reproducibility surface for AI-decision events; AML decisions inherit the same reproducibility shape.

---

## §H. Composition with FFIEC BSA/AML Examination Manual

The FFIEC BSA/AML Examination Manual (2021 ed.) names the four pillars of an effective BSA/AML program:
1. **Governance** — documented BSA officer, board approval, policies.
2. **Risk assessment** — documented customer, product, geography risk.
3. **Internal controls** — TM, sanctions screening, KYC, SAR filing.
4. **Independent testing** — internal audit and SOC.

The chain supports each pillar with integrity-bound evidence:
- Pillar 1 (Governance): `audit.aml.remediation_milestone_tracked` events, threshold-tuning approvals, MOU compliance.
- Pillar 2 (Risk assessment): `audit.kyc.customer_risk_score_computed`, `audit.kyc.customer_tier_assigned` events.
- Pillar 3 (Internal controls): all TM, OFAC, SAR, CTR, BO events.
- Pillar 4 (Independent testing): `audit.model_risk.validation_test_completed`, `audit.aml.performance_report_generated`, lookback events.

The chain is the consolidated audit trail across all four pillars.

---

## §I. FinCEN AML Program Effectiveness Rule (2024) composition

The 2024 Rule emphasizes effectiveness, not just compliance. Effectiveness includes:
- Alert volume trends.
- False-positive rates.
- SAR filing rates.
- Model-performance metrics.
- Lookback-finding rates.
- Remediation timeliness.

### §I.1 Effectiveness reporting from the chain

The institution can produce a one-page Effectiveness Report each month from chain events:
- Alerts generated, triaged-true-positive, triaged-false-positive (from `audit.tm.alert_*`).
- SARs filed, SAR-decision-to-filing latency (from `audit.aml.sar_*`).
- OFAC hits true-vs-waived (from `audit.ofac.*`).
- Model-validation passes / fails (from `audit.model_risk.*`).
- Lookback findings discovered and remediated (from `audit.aml.lookback_*`).

The Effectiveness Report is itself emitted as `audit.aml.performance_report_generated` so the report is integrity-bound.

### §I.2 Chief BSA Officer annual certification

The Chief BSA Officer's annual certification to the board cites the chain:

> Our AML program operates with continuous chain-of-custody-backed audit evidence. Independent verification confirms our records on AML decisions are not altered after capture. The chain supports our four-pillar effectiveness framework with integrity-bound evidence on alerts, triage, SAR filings, OFAC dispositions, KYC risk scoring, model validation, and lookback remediation.

---

## §J. Composition summary

| Concern | Chain event | Audit procedure | Examiner reference |
|---|---|---|---|
| TM alert generation | `audit.tm.alert_generated` | §B.1 | §E.3 |
| TM alert suppression | `audit.tm.alert_suppressed` | §B.2 | §E.3 |
| OFAC screening | `audit.ofac.screening_result` | §B.3 | §E.3 |
| SAR decision | `audit.aml.sar_decision` / `audit.aml.sar_filed` | §B.4 | §E.3 |
| SAR narrative AI provenance | `audit.aml.sar_narrative_drafted` | §B.5 | §E.3 |
| KYC risk scoring | `audit.kyc.customer_risk_score_computed` | §B.7 | §E.3 |
| Lookback investigation | `audit.aml.lookback_*` | §B.6 | §E.3 |
| Model validation | `audit.model_risk.validation_test_completed` | §B.7 | §G.2 |
| Performance monitoring | `audit.aml.performance_report_generated` | §B.8 | §I.1 |
| Multi-region continuity | All AML events | §B.9 | §E.3 |
| Customer dispute | All relevant events | §D | §D |
| FinCEN look-back | All AML events for cohort | §C | §C |

---

## §K. Future scope (informative; not v1.0a normative)

The following items are noted for future evolution. They do not modify v1.0a; the spec body and overlay remain stable.

### §K.1 FinCEN e-filing wire format

Once the AML schema stabilizes, a future spec version may define a FinCEN e-filing format (analogous to current SAR XML) directly from chain entries. The schema in §A is designed to anticipate this without modification.

### §K.2 Cross-institution lookback (314(b))

USA PATRIOT Act §314(b) allows institutions to share information about suspected money laundering or terrorist financing. A future spec version may define cross-institution chain-event sharing under §314(b) protections. The chain's tokenization in §F is designed to support this without exposing original customer identities to receiving institutions.

### §K.3 OFAC list-update integration

Future evolution may integrate OFAC list updates as chain events (`audit.ofac.list_update_received`) so screening-list provenance is integrity-bound. Currently treated as institution-side reference data outside the chain scope.
