# Audit-procedures — Cluster D addenda

> **What this doc is.** Sidecar audit procedures (P-50 through P-55) added during the Cluster D close-out covering healthcare regulatory-timeliness, state AI laws, healthcare 24/7 readiness, clinician-override capture, RNG-source registry cross-check, and verifier amendment-awareness. The procedures are written in the same shape as `docs/audit-procedures.md` P-1 through P-40 and integrate into the same SOC and FFIEC examination evidence repository.

> **Why a sidecar.** Cluster B owns the next round of additions to `docs/audit-procedures.md` directly. Cluster D's additions land in this sidecar to avoid collision; future merging into the main audit-procedures.md is editorial.

> **Audience.** SOC engagement leads, FFIEC examiners, and the institution's internal-audit team running healthcare and chain-amendment sample procedures.

---

## P-50. CMS-0057-F prior-authorization timeliness sample

**Control claim.** AI-driven prior-authorization decisions affecting Medicare Advantage, state Medicaid managed-care, CHIP managed-care, and Qualified Health Plan members meet the CMS-0057-F timeliness window (7 calendar days standard; 72 hours expedited).

**Evidence consumed.** Chain entries with `audit.healthcare.regulatory_clock_basis` containing `cms-0057-f`. Chain entries' `audit.healthcare.regulatory_clock_start` and `audit.healthcare.decision_returned_at` attributes.

**Procedure.**

1. Pull chain entries from the audit period where `audit.healthcare.regulatory_clock_basis` includes `cms-0057-f`.
2. Stratify by request type (standard vs expedited; the expedited population is small and warrants a higher sampling rate).
3. Sample 50-100 entries per stratum.
4. For each sampled entry, compute `decision_returned_at - regulatory_clock_start` and confirm the result falls within the rule's window (7 days standard; 72 hours expedited).
5. Record any failures (window exceeded) as findings.
6. Re-run the verifier on each sampled entry per spec §7 to confirm chain integrity for the timeliness evidence.

**Pass.** All sampled entries within the window. Verifier PASS for every sampled entry.

**Escalation.** Failures escalate to the CMS-0057-F program owner for remediation; chronic failures (above 0.5% of sampled entries missing the window) escalate to the AI Governance Committee. Verifier failures follow the standard P-11 escalation path.

**Reference.** `docs/regulator-pack/healthcare-overlay.md` §H7.

---

## P-51. State AI disclosure and oversight sample

**Control claim.** AI decisions affecting residents of states with AI-specific requirements (Colorado AI Act, California SB 1120, others as enacted) carry the required state-specific attributes and the institution's downstream disclosure or oversight workflow operated as required.

**Evidence consumed.** Chain entries with state-specific attributes: `audit.regulatory.disclosure.colorado_aia.*`, `audit.california.sb1120.physician_review_id`, `audit.california.sb1120.physician_review_at`, `audit.california.sb1120.physician_concur`, and analogous attributes for other states as legislation is enacted.

**Procedure.**

1. Identify the resident-state distribution of the institution's AI-decision population for the audit period.
2. Map applicable state-AI-law disclosure or oversight requirements per `docs/regulator-pack/healthcare-overlay.md` §H10 (and any successor matrix).
3. Sample 50 chain entries per applicable state.
4. For each sampled entry, confirm the required state-specific attributes are present and well-formed.
5. For California SB 1120-scope entries, confirm a physician review entry exists in the chain (`audit.california.sb1120.physician_review_id` populated with a valid NPI or NPI-hash per §H1 posture; `audit.california.sb1120.physician_review_at` present and timestamped).
6. Cross-check institution's downstream disclosure or oversight workflow against a sample of entries to confirm the workflow operated.

**Pass.** Required attributes present for every sampled entry; downstream workflow confirmed for every sampled entry.

**Escalation.** Missing or malformed attributes escalate to General Counsel and the AI Governance Committee. Failed downstream workflow steps escalate to the responsible operational owner.

**Reference.** `docs/regulator-pack/healthcare-overlay.md` §H10.

---

## P-52. Healthcare 24/7 operational readiness

**Control claim.** The institution's chain operations remain available through EHR planned downtime, code-blue / mass-casualty events, and ransomware scenarios. The institution's WORM-backup of chain artifacts is current and accessible only by the IR team.

**Evidence consumed.** EHR planned-downtime procedure document; institution's IR-tabletop exercise records; WORM-backup access logs and configuration; ransomware-response runbook.

**Procedure.**

1. Pull the EHR planned-downtime procedure currently in force. Confirm it references chain operations and names the chain's role during downtime.
2. Confirm the procedure has been exercised in the audit period (operational evidence: at least one planned downtime occurred and the chain continued capture; the seal job operated through the window or the cadence-relaxation under §4.3.1 was applied with documented cause).
3. Pull the institution's IR-tabletop exercise records. Confirm at least one exercise in the audit period covered a healthcare-specific scenario (mass casualty or ransomware) and the chain's role was documented in the exercise scenario.
4. Pull the WORM-backup configuration. Confirm cross-region replication operates at the documented cadence (daily minimum). Confirm the IR-team-only access role is provisioned and a sample access test has been performed in the audit period.
5. For ransomware-specific evidence, confirm the chain's append-only storage posture per spec §10.3 ledger-aggregation requirements remains intact (no UPDATE / DELETE roles on the ledger tables).

**Pass.** All four conditions met (downtime procedure exists and has been exercised; tabletop covered healthcare scenario; WORM-backup current with IR-team access tested; append-only posture intact).

**Escalation.** Missing downtime procedure or unexercised tabletop is a deficiency. Stale WORM-backup or absent IR-team access test is a significant deficiency. Loss of append-only posture is a material weakness (immediate audit-committee escalation).

**Reference.** `docs/regulator-pack/healthcare-overlay.md` §H11.

---

## P-53. Clinician-override capture audit

**Control claim.** AI-recommendation chain entries triggered by deployments under Joint Commission MM.05.01.01 or LD.04.01.07 scope are followed by clinician-override entries capturing override-action and reason-code where the clinician took an override action in the EHR.

**Evidence consumed.** AI-recommendation chain entries; clinician-override chain entries (`audit.clinical.override.*` attribute set per `docs/regulator-pack/healthcare-overlay.md` §H12); EHR override records for the sampled deployments.

**Procedure.**

1. Identify the chain entries representing AI recommendations from deployments in Joint Commission MM.05.01.01 or LD.04.01.07 scope.
2. Sample 50-100 entries.
3. For each sampled entry, query the EHR for the corresponding clinician-action record. If the EHR shows an override action, locate the override entry in the chain (linked via `audit.clinical.override.original_recommendation_run_id` and `audit.clinical.override.original_recommendation_seq`).
4. Confirm the override entry's `audit.clinical.override.override_reason_code` matches the institution's controlled vocabulary.
5. Confirm `audit.clinical.override.override_action` is one of the declared values (`accept_with_modification` | `reject_recommendation` | `defer_decision` | `escalate_to_attending`).
6. Confirm the override entry's `audit.clinical.override.timestamp` is consistent with the EHR's recorded action timestamp within the institution's NTP-synchronization budget.

**Pass.** Override entry exists for every sampled entry where the EHR shows an override action; reason codes and actions match institution declarations; timestamps are consistent.

**Escalation.** Missing override entries (override action in EHR but no override chain entry) escalate to the AI Governance Committee. Reason codes outside the controlled vocabulary escalate to the institution's MRM and the Joint Commission survey response owner.

**Reference.** `docs/regulator-pack/healthcare-overlay.md` §H12.

---

## P-54. RNG-source registry cross-check

**Control claim.** The institution's tenant-public-key registry's declared `rng_source` attribute matches the RNG source named in the institution's CC8.1 control description and matches the RNG source recorded on the `master_key.generated` operational event.

**Evidence consumed.** Tenant-public-key registry's `rng_source` attribute; institution's CC8.1 control description; `master_key.generated` operational event log.

**Procedure.**

1. Pull the tenant-public-key registry's `rng_source` declaration for each tenant under audit.
2. Pull the institution's CC8.1 control description for the audit period and locate the RNG-source claim.
3. Pull the `master_key.generated` operational event log entries for each tenant under audit.
4. Cross-check the three sources for consistency:
   - Registry `rng_source` matches CC8.1 RNG-source claim.
   - Registry `rng_source` matches `master_key.generated` event's recorded RNG source.
   - CC8.1 RNG-source claim matches `master_key.generated` event's recorded RNG source.

**Pass.** All three sources name the same RNG source per tenant.

**Escalation.** Mismatch among the three sources is a significant deficiency. The discrepancy may indicate (a) registry not updated after a key generation under a different RNG, (b) CC8.1 stale, or (c) `master_key.generated` event recorded incorrectly. Audit pulls the IKM-generation procedure for the period and traces the actual RNG source operationally.

**Reference.** `docs/herald-vendor-conformance-round-5-response.md` §1 (Q42 closure).

---

## P-55. Verifier amendment-awareness

**Control claim.** The institution's deployed verifier supports the spec amendment level the institution's seals are produced under. The institution's CC8.1 names the verifier amendment-awareness as a tracked attribute.

**Evidence consumed.** Institution's verifier deployment manifest (which release identifier is in production); institution's CC8.1 control description; sample of recent seal records the institution actually produces.

**Procedure.**

1. Pull the institution's verifier deployment manifest. Confirm the release identifier is recorded in the institution's CC8.1.
2. Pull a sample of recent seal records (at least 5 seals from each month of the audit period).
3. For each sampled seal, run the deployed verifier against the seal and confirm PASS.
4. Cross-check the seal records' `sign_payload_version` attribute (where present per the v1.0a amendment) against the verifier's declared amendment level. The verifier's amendment level must support the `sign_payload_version` values present in the institution's seals.
5. If the institution has progressed to a new spec amendment level during the audit period, confirm the verifier deployment was refreshed at the cutover and the institution's CC8.1 was updated.

**Pass.** Verifier supports the amendment level of all sampled seals. CC8.1 names the verifier amendment-awareness. Cutover (if any) is documented in change-management.

**Escalation.** A verifier producing `signature verification failed` against current-amendment seals indicates a verifier-staleness gap. The institution refreshes the verifier deployment and updates CC8.1. The audit finding is procedural, not cryptographic — the chain's integrity is unaffected; the institution's posture against current-amendment seals is the issue.

**Reference.** `docs/herald-vendor-conformance-round-5-response.md` §6 (Q47 closure).

---

## Composition with `docs/audit-procedures.md`

The procedures in this sidecar use the same evidence repository and the same control-evidence-events shape as the procedures in `docs/audit-procedures.md`. SOC engagements and FFIEC examinations consume both files in parallel. Future editorial passes may merge this sidecar into the main `audit-procedures.md` at the working group's discretion; the procedures themselves are stable and operational from the date of this addenda's publication.

The procedures align with:

- **`docs/regulator-pack/healthcare-overlay.md`** for the healthcare-specific procedures (P-50, P-51, P-52, P-53).
- **`docs/herald-vendor-conformance-round-5-response.md`** for the chain-amendment procedures (P-54, P-55).
- **`docs/internal-audit-evidence-pack.md`** for the IIA-Standard grounding the procedures inherit when run by the third line.
