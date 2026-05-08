# Audit committee summary

> **What this doc is.** Three-page brief for the bank's audit committee chair. Parallel to the management summary (CEO-facing) and the MRM-committee brief; this is for audit-committee-level oversight.

## What the audit committee is overseeing

The chain-of-custody implementation is a control supporting the integrity of AI agent decisions captured by the bank. The audit committee oversees the institution's controls; the chain is one such control. This brief informs the committee about what the chain does, what it requires of the institution, and what the committee should expect to see in oversight materials.

## Why the audit committee should care

Three reasons:

1. **AI agent decisions affect customers.** The bank's AI is making customer-impacting decisions; the audit committee's oversight role extends to controls over those decisions.
2. **Regulatory expectations are growing.** FFIEC, federal banking agencies, and the CFPB all express increasing expectations around AI control. Logging integrity (FFIEC IT Handbook II.C.10) is a specific expectation; the chain satisfies it.
3. **Audit independence depends on integrity-bearing records.** Internal and external audit's effectiveness depends on the records they audit being authentic. The chain provides that authenticity property.

## What the chain does, in committee-relevant terms

The chain captures every AI agent decision the bank makes and produces a daily integrity-bearing seal. Three properties:

- **Authentic capture.** A captured event was produced by the institution's legitimate AI agent process at the claimed time and run.
- **Tamper evidence.** Any modification to a captured event — insertion, deletion, alteration — is detectable by an independent verifier.
- **Independent verifiability.** A regulator with no access to the institution beyond a public key can verify both properties without trusting the institution.

The committee's oversight of the chain is satisfied when the institution operates the controls that support these properties and when independent verification confirms the controls operate.

## What the committee should expect to see

### Quarterly reporting

The chain-operations team or internal audit reports:

- **Verifier-run results.** Pass / pass-with-anomalies / fail per quarter, with quarter-over-quarter trends.
- **Anomalies summary.** What anomalies appeared, what was their root cause, what was remediated.
- **Seal-cadence and notification status.** Confirmation that the seal cadence is operating per the institution's declared cadence and any notification thresholds (72-hour seal-delay) were not exceeded.
- **Master-key reconciliation status.** Confirmation that key-fingerprint reconciliation (spec §10.1) is operating; any unmatched fingerprints (recorded as `fingerprint_unmatched_count > 0` on the `master.reconciliation_completed` operational event) are explained.
- **Operational events summary.** High-severity operational events (chain-detected events, HSM operation failures, configuration changes) with status.

### Annual review

The committee receives an annual summary:

- **Control description update.** What changed in the chain configuration over the year (cadence, vendor, deployment topology).
- **CUEC operation.** The institution's complementary user-entity controls (`docs/control-map/CUECs.md`); confirmation each operates.
- **Independent assessment.** Internal audit's independent assessment of the chain controls; external auditor's view from the SOC engagement.
- **Cost trend.** Year-over-year cost trajectory.
- **Forward-looking.** Anticipated changes (cadence relaxation requests, vendor changes, regulatory engagement).

### Incident reporting

For chain-detected incidents:

- **Severity-assessed events.** What happened, what severity, what response.
- **Time-to-resolution.** How quickly the institution detected, contained, remediated.
- **Lessons learned.** What process or control changes resulted.
- **Notification status.** Whether regulators were notified, when, how the notifications were received.

The IR playbook (`docs/incident-response-playbook.md`) defines the committee's notification thresholds.

## Oversight questions the committee should ask

Standard questions for the chain owner:

1. Has the chain operated continuously over the period? Were there any sealing delays beyond the 72-hour threshold?
2. Are there any chain-detected events under investigation? What is the timeline for remediation?
3. What is the institution's posture on master-key compromise? Is reconciliation operating at the documented cadence?
4. Has the institution validated the verifier binary against the reproducible-build property at least once during the period?
5. Has the institution exercised its IR playbook for chain-detected events at least once during the period?
6. Has the institution's cost for the chain operated within budget?
7. Has the institution informed the regulator of any material chain-related events during the period?

The chain operations team's answers, supported by the operational events log and verifier output, satisfy the oversight role.

## Where chain meets the rest of the committee's oversight

| Committee oversight area | How chain composes |
|---|---|
| Internal control over financial reporting (ICFR) | Chain provides integrity-bearing records of AI decisions affecting financials; supports the ICFR program |
| Cybersecurity oversight | Chain is one cybersecurity control; the committee receives chain-detected event reports as part of broader cyber reporting |
| Vendor management | Chain composes with vendor management; the committee reviews vendor SOC reports for chain implementations |
| Risk management oversight | Chain is one risk control; risk committee reports on chain-related risks; audit committee reviews effectiveness |
| Regulatory examination response | Chain output is part of the response package; the committee receives the regulator's findings as part of broader examination summary |

## Decisions the committee may take

The committee may, based on chain-related reporting:

- Approve the institution's chain configuration changes (cadence relaxation, master-key custody changes)
- Direct the institution's chain-ops team to investigate specific chain-detected events further
- Recommend changes to the IR playbook or other procedures
- Direct internal audit to focus on specific chain controls
- Approve the institution's cost budget for chain operations

The chain doesn't dictate committee decisions; it provides the input.

## Reference materials available to the committee

- The full chain spec at `spec/chain-of-custody-v1.md`
- The chain's threat model at `docs/design/09-threat-model.md`
- The CUECs at `docs/control-map/CUECs.md`
- The TSC mapping at `docs/control-map/TSC-mapping.md`
- The IR playbook at `docs/incident-response-playbook.md`
- The cost model at `docs/cost-model.md`
- The most recent verifier report (provided by the chain-ops team)
- The most recent SOC report covering the chain (when available)

## Glossary (one paragraph)

The chain uses standard cryptographic terminology: HMAC (a keyed hash), Merkle tree (a hash-of-hashes structure), HSM (a hardware-protected signing device), Ed25519 (the digital signature algorithm). The full glossary is at `docs/design/10-glossary.md`. The committee does not need to memorize cryptographic detail; the chain operations team reports in plain terms.

## Update cadence

This brief is reviewed and updated:

- Annually, with the institution's control-description refresh
- After any material chain configuration change (cadence relaxation, vendor change, master-key rotation event)
- After any spec version update affecting the institution

The committee chair receives the updated brief at the next scheduled meeting after each update.
