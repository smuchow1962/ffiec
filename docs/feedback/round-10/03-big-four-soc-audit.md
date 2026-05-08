# Round 10 — Big Four SOC Audit Engagement Partner review

**Reviewer.** Diane Mehrotra. Big Four SOC Audit Engagement Partner. 23 years issuing SOC 1 Type II and SOC 2 Type II opinions; last 12 focused on AI-adjacent fintech and bank-tech vendors. Recent engagement-quality reviewer on three SOC reports for chain-of-custody-style cryptographic controls.

**Reading angle.** Auditability. Every control claim must be testable mechanically. Operational evidence trails must be one-to-one with claims. The user-entity audience must consume the report cleanly.

**Scope of this read.** The chain-of-custody v1.0-final spec (with focus on §7 verification and §10 operational requirements), the audit procedures (especially the updated P-3 and P-6), the operational-events schema (the new `audit_file.truncation_detected` and `master_key.retired` events plus the `master_version` → `key_version` vocabulary clarification), the Section 4 description starter, the CUECs and TSC mapping, the anomaly template (with the new failure modes and the higher "operationally explained" bar), the user-entity summary, and the sample verifier report (with the rotation day, the FAIL day, the `--strict` invocation, and the new TSC criterion appendix).

I did not read prior round-9 feedback or the round-8-pre-rework material. This is a first-look read.

---

## Headline assessment

The corpus is unusual for its market. SOC engagement teams normally spend the first quarter of an engagement reverse-engineering the service organization's evidence model into something testable. Here that work is already done. The verifier's twelve-step procedure (spec §7) is structured so each named failure mode resolves to a specific operational event, a specific audit procedure, and a specific TSC sub-criterion. The new sample-report appendix that maps verifier output lines to TSC criteria (page 285+) is the artifact I would normally have to construct myself in working papers — having it shipped as an examiner-distributable means the SOC team's cross-walk is the institution's cross-walk, which is good for audit consistency across institutions.

I would issue an unmodified opinion against this control set if the institution operated the CUECs as documented and the verifier ran clean across the period. The SOC-readiness of the documentation is at the high end of what I have seen in the AI-adjacent fintech segment.

The remainder of this review is the small set of gaps and friction points that would surface during engagement planning or fieldwork.

---

## What is strong, by SOC dimension

### Testability of control claims

Each numbered audit procedure (P-1 through P-21) names the population, the evidence source, and the test. P-6 (key-fingerprint reconciliation) is exemplary: it specifies the population (`(tenant_id, key_version, key_fingerprint)` triples observed in chain entries), the evidence source (`master.reconciliation_completed` event with `fingerprint_unmatched_count` field), the test (cross-recompute against the IKM roster), and the failure-mode disposition (the three named root causes with §7 step 8 as the catch-point).

P-3's per-tenant determinism replay is the procedural detail SOC teams have been asking for in this segment for two years. More on this below — it is good but has one open item.

### Mechanical evidence consumption

The operational-events schema is canonical-form JSON with a stable namespace. The events I would query during fieldwork — `seal.job_started`, `seal.job_completed`, `chain.verification_failure`, `master.reconciliation_completed`, the new `master_key.retired`, the new `audit_file.truncation_detected` — all carry a `correlation_id` that ties events to their operational story. The `chain.verification_failure` event's `step` field references the spec §7 procedure step, which means the SOC team can SQL-query the institution's log store for `step=8` failures and immediately know to investigate the IKM-roster row rather than the chain content. That is the difference between a one-day analysis and a one-week analysis.

The "How SOC teams use this" table at the bottom of `control-evidence-events.md` is the table I would normally hand-author at the start of an engagement.

### Sample-report TSC appendix

The appendix added to `sample-report.md` (lines 285+) is the most engagement-ready artifact in the corpus. Verifier output line on the left, TSC criterion (and ancillary CC sub-criteria) on the right. When I close out fieldwork, I attach the verifier bundle to the working papers and reference this appendix for the evidence allocation. The appendix is short enough to inspect in one read and complete enough that I do not have to author my own.

### Anomaly template

The template's higher-bar treatment of `key_fingerprint mismatch` (the four-element evidence requirement) is exactly the kind of structured assertion I want to see in a control-evidence repository. "Operationally explained" is no longer an institution's free-text claim — it is a four-piece evidence package that either is in the file or is not. That converts a judgment call into a checklist test, which is what audit fieldwork is supposed to be.

The same higher-bar treatment for `hkdf_inputs_digest mismatch` and `audit_file_truncation_detected` is welcome and consistent.

### Vocabulary clarification

The `master_version` → `key_version` (per-entry) and `key_versions` (seal-record list, plural) clarification is a small piece of writing but it removes a real footgun. In my prior review of this segment I have seen institutions confuse "the custodian's KMS-side key alias" with "the chain's per-entry generation integer" and produce control descriptions that named the wrong field. The note at line 97 of `control-evidence-events.md` ("Custodian-side operational labels remain whatever the institution's KMS / HSM operator framework uses and are not directly consumed by the chain") closes that confusion cleanly and lets the institution's control description stay accurate when the KMS-side label and the chain-side `key_version` diverge.

---

## Gaps and friction points

These are the items I would raise during engagement planning. None are blockers; all should be addressed in the next round of documentation hardening.

### 1. P-3 self-test shape — sampling cadence is unstated

The new P-3 procedural detail is good. The two acceptable shapes (self-test and sample-comparison) are clear, the comparison target is unambiguous ("any `key_fingerprint` previously stamped on a chain entry for the same `(tenant_id, key_version)`"), and the working-paper requirement is named.

What is missing: the **cadence** at which the institution runs the self-test. The SOC team needs to test that the self-test was operated during the reporting period, not just that a self-test exists. P-3 should specify either:

- a minimum cadence (e.g. weekly, aligned with P-6 reconciliation), OR
- a rule that the self-test runs once per IKM rotation (so the cadence floats with the institution's rotation policy), OR
- explicit deference to the institution's documented cadence with a SOC-test that the cadence is reasonable.

Without a cadence rule, an institution could run the self-test once at SDK provisioning, never re-run it, and pass P-3. That is technically conformant but defeats the per-tenant determinism assurance over the period.

**Recommended wording.** Add to P-3:

> The institution's documentation MUST specify the cadence at which the self-test (or sample-comparison) is operated. The minimum acceptable cadence is the lesser of (a) the institution's IKM rotation cadence, or (b) the institution's P-6 reconciliation cadence. The SOC team confirms the documented cadence was operated during the reporting period.

### 2. P-3 sample-comparison shape — non-production IKM provisioning is unstated

The sample-comparison shape says "the institution provides a non-production IKM the SOC team derives against directly." That is fine, but it leaves open: how does the institution provision the non-production IKM to the SOC team safely? IKMs are the highest-sensitivity material in the system. A loose process here would have the SOC team handling production-equivalent secrets without a documented chain-of-custody.

**Recommended wording.** Add to P-3:

> When the sample-comparison shape is operated, the institution's procedure MUST document: (a) the non-production IKM's provenance (it MUST NOT be reused from production, even from a different tenant); (b) the secure transport mechanism to the SOC team (typical: a sealed envelope from the HSM custodian, opened in the SOC team's evidence room and shredded after the test); (c) the working-paper retention policy for the IKM (the SOC team's working paper records the fingerprint and the test result; the IKM bytes themselves are NOT retained).

### 3. CUECs.md — vocabulary drift between CUEC-CRY-04 and the spec

CUEC-CRY-04 (line 28) reads:

> The institution operates session-key-id reconciliation at a documented cadence (no more than weekly per spec recommendation).

The spec §10.1 and audit-procedures P-6 use "key-fingerprint reconciliation," not "session-key-id reconciliation." The earlier `session_key_id` field was removed in the v1.0-rework (per the change log) in favor of `key_fingerprint`. The CUEC text should be updated to match.

**Recommended wording.** "The institution operates **key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1)." This is a minor edit but it matters for audit-evidence cross-references — when the SOC team queries for "session-key-id" in the institution's procedures and finds nothing, the CUEC has lied. When the institution's procedures use "key-fingerprint" and the CUEC uses "session-key-id," the SOC team has to reconcile the vocabulary mid-engagement.

### 4. Section 4 template — the `[bracketed]` items for new operational events

The Section 4 starter at `section-4-template.md` predates the new operational events. Specifically, "Procedures" (lines 53-58) lists the daily seal job, HSM PIN rotation, master key rotation, session-key-id reconciliation (with the same vocabulary drift as point 3), verifier validation, and IR. It does not mention:

- The institution's **IKM retirement procedure** that the new `master_key.retired` event documents (spec §10.9 requires the institution's procedure to be documented in the control description).
- The institution's **mid-write-truncation recovery procedure** that the new `audit_file.truncation_detected` event references (the SDK-local SQLite buffer plus IR Scenario 9).
- The institution's **`--strict` verifier-run cadence** vs the operational-anomaly verifier-run cadence (the sample report shows both invocations; the Section 4 description should name which the institution operates and on what schedule).

**Recommended wording.** Add three bullets under "Procedures":

> - IKM retirement on `[institution-defined]` cadence per spec §10.9; retirement requires `chain_entries_referencing_remaining = 0` and is recorded as a `master_key.retired` operational event signed off by `[role]`.
> - Mid-write-truncation recovery per IR Scenario 9; the SDK-local SQLite buffer and OTLP retention are the documented recovery sources.
> - Verifier runs on `[institution-defined]` cadence; substantive SOC-evidence runs are invoked with `--strict`; anomaly-evaluation runs are invoked without.

### 5. TSC mapping — `audit_file.truncation_detected` and `master_key.retired` need a TSC home

The TSC mapping table covers the original event set but does not explicitly map the two new events to their TSC sub-criteria. The mapping is implied (`audit_file.truncation_detected` belongs under CC7.2 system monitoring; `master_key.retired` belongs under CC6.1 logical access and CC8.1 change management) but a SOC engagement team is going to want it stated.

**Recommended wording.** Add to the TSC mapping:

| TSC | Mapping |
|---|---|
| CC7.2 (Monitor system components) | ... `audit_file.truncation_detected` operational event; verifier-emitted on §4.1 mid-write refusal |
| CC6.1 (Logical access controls — administer) + CC8.1 (Change management) | ... `master_key.retired` operational event; institution's IKM retirement procedure under spec §10.9 |

### 6. Sample-report appendix — two minor coverage gaps

The new TSC criterion appendix at `sample-report.md` lines 285-305 is excellent, but two verifier output lines are present in the per-day detail and not mapped:

- `Late-binding count: N` — this is a Processing Integrity signal (events that arrived after the seal but before the next day's seal). It belongs under PI1.2 (output completeness — the institution did not lose the event, it sealed it in the next day). The appendix mentions late-binding rate anomalies under A1.2 but does not map the per-day count.
- `KMS handle URI` with the `plaintext-` prefix case — the appendix maps non-`plaintext-` URIs to CC6.7. It does not map the `plaintext-` URI case explicitly, which is the failure-mode lookup the SOC team would use during a `--strict` engagement to confirm the verifier refused the dev-mode seal. CC6.8 (prevent unauthorized software) is the right home — the spec §10.7 compile-time exclusion plus the verifier refusal is the defense-in-depth pattern.

**Recommended wording.** Add two rows:

| Verifier output line | TSC criterion / SOC procedure |
|---|---|
| `Late-binding count: N` | PI1.2 (processing integrity — output completeness); the event is sealed in the next day's seal, not lost |
| `KMS handle URI` with `plaintext-` prefix under `--strict` invocation | CC6.8 (prevent unauthorized software); verifier-side refusal of dev-mode seals; pairs with spec §10.7 compile-time exclusion |

### 7. Anomaly template — the four-piece evidence package needs a reciprocal SOC test

The anomaly template's higher bar for `key_fingerprint mismatch` ("the institution's evidence MUST include all four") is excellent for the institution. It does not have a reciprocal SOC test in `audit-procedures.md`. The SOC team needs a procedure that says "for each `key_fingerprint mismatch` anomaly in the period, confirm the institution's anomaly-evaluation record contains all four required evidence pieces." Without that procedure, the SOC team has to invent the test on the fly.

**Recommended wording.** Add to `audit-procedures.md`:

> **P-22. `key_fingerprint mismatch` anomaly evidence completeness.**
>
> - Pull all `key_fingerprint mismatch` (spec §7 step 8) anomalies from the verifier output for the period
> - For each anomaly, pull the institution's anomaly-evaluation record per `docs/anomaly-documentation-template.md`
> - Confirm the record contains all four required evidence pieces: (1) IKM-roster row identification; (2) change-management approval for the roster correction; (3) re-verification result on the corrected roster; (4) reconciliation cross-check showing `fingerprint_unmatched_count=0` for the affected pair
> - Any anomaly missing one or more of the four pieces is documented as an unmitigated control gap; the SOC team escalates per the engagement's standard finding workflow

A parallel P-23 for `hkdf_inputs_digest mismatch` and a P-24 for `audit_file_truncation_detected` would close the symmetry.

### 8. User-entity summary — the new events deserve a one-line mention

The user-entity summary is appropriately high-level for its audience. The downstream business partner reading it does not need to know the chain-event JSON schema. They probably should know that the verifier and the institution have evidence-events for two specific failure modes (audit-file truncation and IKM retirement) — those are the two scenarios where a downstream partner asking "how would I be told if something went wrong" deserves a concrete answer.

**Recommended wording.** Add to "What the SOC report attests to" or "How the user entity uses the SOC report":

> The institution's operational evidence stream includes events for low-frequency, high-impact scenarios: a mid-write file truncation detected at verification (`audit_file.truncation_detected`) and master-key retirement (`master_key.retired`). The SOC report attests that the institution emits these events when the corresponding action occurs and that the institution's response procedure is documented (IR Scenario 9 and spec §10.9 respectively). A downstream user entity relying on the chain has assurance through the SOC opinion that these scenarios are operationally accounted for.

---

## On the higher-bar "operationally explained" wording

Three specific things on this. I think they are right but I want to call out the second-order effects.

### `key_fingerprint mismatch` — the four-piece package

The four-piece evidence package is the right shape. In particular, requirement (4) — the reconciliation cross-check showing `fingerprint_unmatched_count=0` for the affected pair on the most recent reconciliation — closes the loop between the verifier finding and the operational evidence trail. Without (4), an institution could remediate the IKM roster, re-verify successfully, and have no operational evidence that the remediation actually persisted. With (4), the operational evidence is the institution's standing reconciliation event, which is independently produced and stored.

**Caveat.** The four-piece package is going to be hard for institutions to operate the first time. SOC engagement teams should expect to spend time in the first year of the post-rework period helping institutions stand up the evidence pipeline. The template should probably carry an example completed record so institutions have a worked example to copy. That is template-quality, not normative-quality, but it would reduce friction.

### `hkdf_inputs_digest mismatch` — the binding to format-version migrations

The wording "operationally explained ONLY when the root cause is identified as a known SDK / verifier constants change (e.g., upgrade across a `format_version` boundary documented in the change-management log)" is correct and conservative. It correctly forces the institution to either (a) have a change-management record, or (b) treat the mismatch as a format-construction defect.

**Caveat.** The wording does not address the case where the mismatch is between two *production* deployments at the same `format_version` — i.e., where two independent SDK installs produced different `hkdf_inputs_digest` for the same tenant. That would be a deployment-drift case, not a defect, and the institution would need a third disposition (re-deploy to consistent constants, then re-verify). The template should contemplate this case. It is rare but not impossible at institutions running multi-region deployments.

### `audit_file_truncation_detected` — the recovery-source enumeration

The wording is clean: the writer-process crash must be identified in host-level monitoring AND the dropped event must be recovered from SDK-local SQLite or upstream OTLP retention. The `recovery_outcome` field on the operational event (`complete | partial | unrecoverable`) gives the SOC team a clean three-way disposition.

**Caveat.** The `partial` outcome is not addressed in the anomaly template. A partial recovery means some of the dropped events are recovered and some are not. The institution's disposition of the unrecovered subset is the operational question. The template should specify that `partial` outcomes are documented as an unrecoverable gap for the un-recovered events specifically, with IR Scenario 9 disposition for those events, and the recovered events are restored to the chain via the institution's documented recovery procedure.

---

## Items I would have flagged but the corpus already addressed

These are points I drafted while reading and crossed out before submission, because the corpus addresses them. I list them so the working group knows the corpus has good coverage on the engagement-quality dimensions.

- The `key_fingerprint` check happens BEFORE any MAC compute (spec §7 step 8). This is the right ordering for a SOC-tester's mental model. A botched rotation produces one fingerprint-mismatch failure, not 142,087 MAC mismatches. The error message points at the IKM roster, not the chain. Audit-procedures P-6 and the sample-report FAIL-day example both call this out.
- The `--strict` invocation is shown in the sample report alongside the non-strict invocation, with a clear note ("Use `--strict` for substantive SOC testing"). Engagement teams have explicit guidance on which to use when. Many SOC engagements in this segment have suffered from auditor confusion about whether to run with or without strict mode; this resolves it.
- The seal record's `dev_mode` field plus the `--strict` refusal of dev-mode seals plus the spec §10.7 compile-time exclusion is the three-layer defense the segment has been moving toward. Most institutions only do one or two layers; this corpus does all three and the verifier-side refusal is the layer the SOC team can mechanically test.
- The cross-chain-lift defense at §7 step 4 (the `event.tenant_id == header.tenant_id` and `event.run_id == header.chain_id` assertions) is mentioned in the audit procedures' verifier methodology section in the sample report. SOC teams testing "an event from tenant A cannot be lifted into tenant B's chain" have a named procedure step to point at.
- The fail-closed semantics in §7 (an absent IKM is a `--strict` FAIL; otherwise PASS-WITH-ANOMALY) is clear and the sample report demonstrates the PASS-WITH-ANOMALY case. SOC teams running without `--strict` in operational-anomaly mode get a clear disposition rule.
- The HSM unavailability 72-hour notification SHOULD (spec §4.3.1) is wired through to CUEC-OPS-05, P-9, and the anomaly-template severity table. The control claim, the audit procedure, and the anomaly disposition are linked end-to-end.
- The retention coupling between IKMs and chain entries (spec §10.9) is addressed by the new `master_key.retired` event with the `chain_entries_referencing_remaining` field. The verifier reports `unknown_key_version` if the institution prematurely retires; the operational event documents the institution's intent. The two together close the retention-control evidence trail.

---

## Engagement-readiness summary

| Dimension | Assessment |
|---|---|
| Section 4 description starter | Engagement-ready with the three additions in point 4 |
| Audit procedures (control-claim → test) | Engagement-ready with the additions in points 1, 2, 7 |
| Operational evidence schema | Engagement-ready as-is; the vocabulary clarification is welcome |
| TSC mapping | Engagement-ready with the additions in points 5, 6 |
| Anomaly template | Engagement-ready with the additions in the "On the higher-bar" section |
| User-entity summary | Engagement-ready with the addition in point 8 |
| Sample report | Engagement-ready with the two appendix additions in point 6 |
| CUECs document | Engagement-ready with the vocabulary fix in point 3 |

I would issue an unmodified opinion against this control set. The post-rework documentation is at a level that lets the SOC team focus on testing the institution's operation of the controls, not on reverse-engineering what the controls are.

The eight gaps above are the items I would raise in the engagement-quality review of my own working papers, not blockers. Address them in the next documentation pass and the corpus is engagement-ready at the highest tier I have rated in this segment.

— Diane Mehrotra
Engagement Partner, SOC Attest Services
