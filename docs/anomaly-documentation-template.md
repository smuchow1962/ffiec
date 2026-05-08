# Anomaly documentation template

> **What this doc is.** Template the institution and the SOC team use to document anomalies the verifier reports during a SOC reporting period. Required for SOC engagements where the verifier is run without `--strict` and anomalies must be evaluated in context.

## Why this matters

The verifier reports two classes of output:

- **Pass / Fail per day** — the integrity verdict
- **Anomalies** — operational observations that don't fail the day but warrant attention

For SOC engagements, the anomaly evaluation is the institution's documented assessment that each anomaly has an acceptable operational explanation. Without documented evaluation, the SOC opinion has unaccounted anomalies and the reader cannot evaluate.

This template is the form the institution fills in for each anomaly, signed by the institution's control owner.

## Template

For each anomaly the verifier reports during the SOC reporting period:

```
========================================
ANOMALY EVALUATION RECORD
========================================

Anomaly identifier:     [verifier-assigned ID, e.g., ANM-2026-04-15-001]
Tenant:                 [tenant_id]
Affected day(s):        [YYYY-MM-DD or range]
Anomaly type:           [Sealing delay |
                         Late-binding rate elevated |
                         Clock-skew |
                         Missing seal continuity |
                         Software-key in production |
                         key_fingerprint mismatch (spec §7 step 8) |
                         unknown_key_version (spec §7 step 7) |
                         format_version mismatch (spec §7 step 5) |
                         hkdf_inputs_digest mismatch (spec §7 step 2) |
                         master_key_rotation_observed |
                         audit_file_truncation_detected |
                         Other]
Severity assessment:    [Low | Medium | High | Critical]

DESCRIPTION
  [What the verifier reported, in the verifier's own language]

OPERATIONAL EXPLANATION
  [The institution's explanation of what caused the anomaly]

EVIDENCE REFERENCES
  - [Incident log reference, ticket ID, change-management record]
  - [Operational events from the institution's monitoring stack, with timestamps]
  - [Network/infrastructure logs supporting the explanation]

ROOT CAUSE
  [Confirmed | Probable | Suspected | Under investigation]
  [Summary of the root cause]

REMEDIATION
  [Already remediated | Remediation in progress | No remediation needed]
  [Description of remediation, with target completion date if in progress]

INSTITUTION'S DETERMINATION
  [The anomaly is operationally explained and does not affect chain integrity. |
   The anomaly is operationally explained but reflects a control weakness that
   warrants attention in the next quarterly control review. |
   The anomaly was not adequately explained; escalation under the IR playbook.]

CONTROL OWNER SIGN-OFF
  Owner:            [Name and role]
  Sign-off date:    [YYYY-MM-DD]
  Signature:        [Wet/electronic signature per institution's policy]

REVIEWING ENGAGEMENT (if applicable)
  Engagement:       [SOC 1 / SOC 2 / Internal audit / Other]
  Reviewer:         [Name and firm]
  Disposition:      [Accepted / Accepted with note / Not accepted]
  Reviewer notes:   [Optional]
========================================
```

## When to fill in this record

- During the institution's standard ongoing-monitoring of the chain (typically monthly review of the verifier output)
- During SOC engagement testing when the SOC team identifies the anomaly
- During post-incident review for medium-or-higher severity anomalies

## Severity assessment guidance

Match to the institution's existing severity framework. Reference mapping:

| Anomaly | Default severity |
|---|---|
| Late-binding rate elevated, single day | Low |
| Late-binding rate elevated, sustained | Medium |
| Sealing delay 1–24 hours | Low |
| Sealing delay 24–72 hours | Medium |
| Sealing delay >72 hours without notification | High |
| Missing seal continuity | High |
| Clock-skew anomalies | Low |
| Sealing delay associated with HSM cluster outage | Medium (covered by IR playbook scenario 5) |
| Software-key (dev-mode) in production | **Critical** (covered by IR playbook scenario 6) |
| **`key_fingerprint mismatch`** (spec §7 step 8) | **Critical** (covered by IR playbook Scenario 7; the load-bearing rework primitive — investigate against IKM roster, NOT chain content) |
| **`unknown_key_version`** (spec §7 step 7) | **High** (covered by IR playbook Scenario 8; investigate against IKM-roster retention or provisioning gaps) |
| **`format_version mismatch`** at entry (spec §7 step 5) | **Medium** (covered by IR Scenario 1 follow-up; SDK version-handling defect) |
| **`hkdf_inputs_digest mismatch`** at file pre-flight (spec §7 step 2) | **Critical** (format-construction defect or constants drift) |
| **`master_key_rotation_observed`** (normal-operations during documented rotation) | **Low** when paired with documented `master_key.rotated` event; **Medium** when undocumented (route to Scenario 7 triage) |
| **`audit_file_truncation_detected`** (spec §4.1 mid-write truncation refusal) | **Medium** (writer-side crash-recovery; covered by IR Scenario 9) |
| **`gen_ai_model_identifier_missing`** (spec §7 step 12a) | **Medium** (control-completeness for SR 11-7 reproducibility, NOT chain-integrity; the institution's MRM program loses reproduction surface for the affected entries) |

The institution's risk function approves the severity assignment.

## What "operationally explained" means

An anomaly is operationally explained when:

1. The root cause is identified
2. The cause is consistent with operational reality (not a tampering signal)
3. Evidence supports the explanation (logs, tickets, change records)
4. The explanation does not require speculation about adversary behavior

If any of these is missing, the anomaly is NOT operationally explained and the institution treats it under the IR playbook rather than as a routine evaluation.

### Higher bar for `key_fingerprint mismatch`

For `key_fingerprint mismatch` (spec §7 step 8) anomalies specifically, "operationally explained" carries a stricter bar than for sealing delays or clock-skew. The institution's evidence MUST include all four:

1. **Identification of the IKM-roster row that was wrong.** Specifically: which `(tenant_id, key_version)` pair had an IKM that did not produce the recorded `key_fingerprint`?
2. **Documented change-management approval for the IKM-roster correction.** The correction itself is a change to a security-bearing artifact and requires the institution's standard change-management approval.
3. **Re-run of the affected period.** The verifier MUST be re-run against the corrected roster, demonstrating PASS on the affected period (or a documented gap if PASS cannot be reached because the original IKM is unrecoverable, in which case the gap is filed as an integrity-control failure under IR Scenario 7).
4. **Reconciliation cross-check.** The institution's most recent `master.reconciliation_completed` event MUST show the corrected `(tenant_id, key_version, key_fingerprint)` triple matching the IKM roster (`fingerprint_unmatched_count=0` for the affected pair).

Without all four, the `key_fingerprint mismatch` is NOT operationally explained and is escalated to IR Scenario 7 triage case 4 (suspected unauthorized key substitution; 36-hour clock starts).

### Higher bar for `hkdf_inputs_digest mismatch` and `audit_file_truncation_detected`

For `hkdf_inputs_digest mismatch` (spec §7 step 2): operationally explained when ONE of two acceptable root causes applies:

1. **Known SDK / verifier constants change** documented in the change-management log (e.g., upgrade across a `format_version` boundary).
2. **Multi-region deployment-drift case.** Two production SDKs at the same `format_version` produce different `hkdf_inputs_digest` for the same tenant — typically because the institution's multi-region deployment has drifted (one region's SDK was patched, the other was not, or one region's HKDF constants were misconfigured). The institution's response is re-deployment to consistent constants across all regions, followed by re-verification across the affected period.

Mismatches without one of these two explanations are escalated as format-construction defects under IR Scenario 1.

For `audit_file_truncation_detected` (spec §4.1): operationally explained when the writer process's crash is identified in the institution's host-level monitoring AND the recovery outcome is documented. The `recovery_outcome` field on the operational event takes one of three values:

- **`complete`** — all dropped events are recovered from SDK-local SQLite buffer or upstream OTLP retention; the recovered events are restored to the chain via the institution's documented recovery procedure. Anomaly is operationally explained.
- **`partial`** — some dropped events are recovered, some are not. The recovered events are restored to the chain. The un-recovered subset is documented as an unrecoverable gap for those specific events with IR Scenario 9 disposition. The anomaly record names which `(run_id, seq)` events were recovered and which were not.
- **`unrecoverable`** — no events recovered. The institution treats the gap as an integrity-control failure and starts the 36-hour clock per IR Scenario 9 notification guidance.

### Worked example — `key_fingerprint mismatch` anomaly record

Friction-reducer for the first time an institution operates the four-piece evidence package. Sample completed record:

```
========================================
ANOMALY EVALUATION RECORD
========================================

Anomaly identifier:     ANM-2026-04-15-001
Tenant:                 tenant_acme_prod_us_east_1
Affected day(s):        2026-04-12 through 2026-04-14
Anomaly type:           key_fingerprint mismatch (spec §7 step 8)
Severity assessment:    Critical

DESCRIPTION
  Verifier reported `key_fingerprint mismatch at seq 17:
  looked-up IKM does not match the entry's recorded fingerprint`
  on 4 events across 3 days. Recorded fingerprint:
  2eede65f0f764c97eaf3b3f306a48537. Looked-up IKM produced:
  b94c1a77b40bf5106c66ca6c1c1b4989. (tenant_id, key_version) was
  (tenant_acme_prod_us_east_1, 1).

OPERATIONAL EXPLANATION
  Botched key rotation on 2026-04-12. The HSM operator created
  a new IKM under key_version=1 (replacing the prior IKM), instead
  of incrementing to key_version=2. Chain entries captured between
  the rotation and the discovery on 2026-04-15 stamped the OLD
  fingerprint (b94c...) but the IKM-roster entry for key_version=1
  was the NEW IKM (2eed...). The verifier correctly detected the
  identity mismatch at lookup time before any MAC compute.

EVIDENCE REFERENCES
  - HSM admin change-management ticket CM-2026-04-12-039 documenting
    the (incorrect) rotation procedure
  - master_key.rotated event for the rotation (custodian audit log)
  - master.reconciliation_completed for week 14 with
    fingerprint_unmatched_count=4

ROOT CAUSE
  Confirmed.
  HSM operator misread the procedure and overwrote the existing
  key_version=1 IKM instead of provisioning key_version=2. Procedure
  has been updated to require explicit version-increment confirmation.

REMEDIATION
  Already remediated.
  - Step 1: IKM-roster row for (tenant_acme_prod_us_east_1, key_version=1)
    restored to original IKM bytes from KMS history (2026-04-15 14:30 UTC).
  - Step 2: HSM operator change-management approval for the restoration
    recorded as CM-2026-04-15-002.
  - Step 3: Verifier re-run against 2026-04-12 through 2026-04-14
    period under restored roster — PASS for all affected events.
  - Step 4: Reconciliation cross-check on 2026-04-21 (week 15) shows
    fingerprint_unmatched_count=0 for (tenant_acme_prod_us_east_1, 1).

INSTITUTION'S DETERMINATION
  The anomaly is operationally explained (botched rotation correctly
  remediated) and does not affect chain integrity. The four-piece
  evidence package is complete (see EVIDENCE REFERENCES + REMEDIATION).
  Institution updated the rotation procedure to prevent recurrence.

CONTROL OWNER SIGN-OFF
  Owner:            Jane Doe, Chain Operations Lead
  Sign-off date:    2026-04-22
  Signature:        [electronic signature per institution policy]

REVIEWING ENGAGEMENT
  Engagement:       SOC 2 Type II (period: 2026-Q2)
  Reviewer:         [SOC firm representative]
  Disposition:      Accepted
  Reviewer notes:   Four-piece evidence package complete; cross-checked
                    against P-22 procedure; no engagement finding.
========================================
```

The worked example is template-quality, not normative-quality, but reduces the friction the first time an institution stands up the evidence pipeline.

## SOC engagement use

SOC teams evaluating anomalies:

1. Pull the institution's anomaly-evaluation records for the reporting period
2. Confirm each anomaly has a complete record
3. Sample-test the evidence references (institution's logs, tickets)
4. Confirm severity assessment is consistent with the institution's framework
5. Confirm the institution's determination matches the evidence
6. Document the sampling and the disposition

The SOC team may accept the institution's determination, accept with a noted concern, or not accept (escalate to a finding). The disposition is recorded in the engagement working papers.

## Examination use

FFIEC examiners reviewing anomalies:

- For routine anomalies that the institution has explained, accept the institution's determination
- For elevated-severity anomalies, sample-test the evidence
- For unexplained anomalies, treat as a finding under the examination's standard finding framework

The examination uses the same records the SOC team uses; both consume the same artifacts.

## Retention

Anomaly-evaluation records are retained at least as long as the chain events for the period (typically 7 years; institution-defined). Records are part of the institution's control-evidence repository.

## Review and update

This template is reviewed annually. Updates triggered by:

- New anomaly types the verifier reports (driven by spec version updates)
- Changes to the institution's severity framework
- Lessons learned from SOC or examination engagements

The template is institution-flexible; institutions adapt the structure but preserve the substance (anomaly identification, explanation, evidence, determination, sign-off).
