---
status: normative-template
alignment-reference: AICPA SSAE 18 / AT-C §205 / TSP 100 (2017) — Trust Services Criteria
companion-templates:
  - soc-pack/section-4-template.md (control-activities mapping)
  - soc-pack/sox404-icfr-composition.md (SOX 404 boundary)
date: 2026-05-07
version: 1.0.0
---

# SOC 2 Type II Section 3 — Description of the System (template)

> **What this template is.** A drop-in Section 3 description-of-system template the institution uses when its SOC 2 Type II report covers an FFIEC chain-of-custody implementation under the v1.0a normative specification. The template is structured for direct paste into the institution's SOC 2 report with marked fields the institution fills in. The template addresses TSP 100 description criteria DC1 through DC9 and supports a Trust Services opinion over Security, Availability, Processing Integrity, Confidentiality, and (where applicable) Privacy.

## How to use this template

1. **Pick the topology.** The chain is topology-agnostic per spec §4. The institution selects one of three named topology variants — self-hosted, BYOC, or vendor-hosted — and uses the corresponding subsection in §3.2 (Infrastructure). Subsections for the other variants are deleted in the final report.
2. **Fill in the marked fields.** Every `[INSTITUTION:fillin]` token names a value the institution provides from its own deployment. Tokens are explicit so a reviewer or service auditor can scan for unfilled gaps before issuance.
3. **Pair with Section 4.** Section 4 (control activities) consumes the boundary, scope, and CUEC enumeration this template establishes. Section 4 is templated separately at `soc-pack/section-4-template.md`.
4. **Tie to ICFR if SOX is in scope.** When the chain feeds disclosure-relevant outputs (credit-loss provisioning, valuation models, fair-lending statistical analyses), pair this Section 3 with `soc-pack/sox404-icfr-composition.md` to scope the integrity-of-recording claim against SR 11-7 model validation.
5. **Period-end cutoff.** Use §3.7 (Period-end cutoff procedure) verbatim — the cutoff procedure is part of the description so the service auditor's coverage of the period is unambiguous.

## §3.1 Services provided (DC1 — services provided)

[INSTITUTION:fillin — institution legal name] (the "Service Organization") operates a chain-of-custody system that captures AI-driven decisions in regulated workflows. The system produces an integrity-bound, examiner-verifiable record of:

- The model's prompts and responses at inference time.
- The tools and external services the AI agent invoked.
- The routing decisions that selected providers and policy versions.
- The institution's operational state at the moment of decision.

The captured record is per-event MAC-bound under HMAC-SHA-256 (FIPS 198-1), aggregated daily into Merkle roots per RFC 6962 over SHA-256 (FIPS 180-4), and the daily root is signed by the institution's HSM under Ed25519 (FIPS 186-5). The composition is documented in the FFIEC chain-of-custody specification v1.0a §1.4 (compositional security).

The system is provided to [INSTITUTION:fillin — internal customer scope: institution business units / external customer scope: regulated tenants / both]. Each customer is termed a **tenant** for purposes of this description.

The services in scope of this report:

| Service | Description |
|---|---|
| Chain capture (SDK) | Per-event capture of AI activity into the institution's audit ledger under per-event MAC. |
| Daily seal job | Aggregation of each tenant-day's events into a Merkle root and HSM-signed seal record. |
| Verifier (offline) | Re-execution of the §7 procedure by examiners or SOC engagements against retained ledger artifacts. |
| Tenant key management | Provisioning, rotation, and retention of tenant IKM material under [INSTITUTION:fillin — HSM/KMS product]. |
| Operational events | Emission of `master.*`, `seal.*`, `ledger.*`, `chain.*`, `hsm.*` events per spec §10.2 for control evidence. |

Services explicitly out of scope (per spec §1.2 epistemic scope):

- Factual accuracy of the AI's statements.
- Policy compliance of the AI's decisions.
- Statistical bias evaluation of decision populations.
- Model-training integrity (training data provenance, training-pipeline reproducibility).

## §3.2 Infrastructure (DC2 — infrastructure)

The institution operates one of three topology variants. Use the subsection that matches the deployment; delete the others.

### §3.2.A Self-hosted topology (DC2 — infrastructure)

The institution operates the SDK, the audit ledger, and the HSM under its own administrative control.

| Component | Description |
|---|---|
| SDK runtime | [INSTITUTION:fillin — language and version, e.g., Python 3.11 / .NET 8 / Java 17] hosting the chain SDK in the institution's application processes. |
| Audit ledger | [INSTITUTION:fillin — DB product and version, e.g., PostgreSQL 15 with append-only role grants per spec §10.3]. |
| HSM | [INSTITUTION:fillin — model and FIPS validation certificate, e.g., AWS CloudHSM v2 (Cert. #3718, FIPS 140-2 Level 3)]. The seal-signing key is held non-extractable on this device. |
| Network | [INSTITUTION:fillin — VPC/subnet description; OTLP transport between SDK and ledger over TLS 1.3 per spec §5.1]. |
| Identity | [INSTITUTION:fillin — IAM/SSO platform, e.g., Okta + AWS IAM]. The seal-job operator role grants `sign` only per spec §10.5. |
| Monitoring | [INSTITUTION:fillin — observability platform, e.g., Datadog / New Relic / OpenTelemetry collector]. |

### §3.2.B BYOC topology — Bring-Your-Own-Cloud (DC2 — infrastructure)

The institution operates the SDK and ledger in the institution's cloud account; a vendor provides the SDK and seal-job software but does not hold the institution's HSM credentials or ledger access.

| Component | Description |
|---|---|
| SDK runtime | [INSTITUTION:fillin — language and version]. SDK code provided by [INSTITUTION:fillin — vendor name and product version]. |
| Audit ledger | [INSTITUTION:fillin — DB product]; institution holds DB credentials. Vendor has no read or write access. |
| HSM | [INSTITUTION:fillin — model and FIPS certificate]. Held in institution's cloud account. Vendor has no operator role. |
| Vendor-provided software supply chain | [INSTITUTION:fillin — SDK distribution mechanism: signed container image / signed wheel / verified Maven artifact]. SLSA attestation evidence retained per spec §10.13. |
| Network | [INSTITUTION:fillin — VPC topology]. Vendor connectivity restricted to documented support paths under [INSTITUTION:fillin — break-glass procedure]. |

### §3.2.C Vendor-hosted topology (DC2 — infrastructure)

A vendor hosts the audit ledger and HSM. The institution operates the SDK in its application processes; the vendor operates the ledger, the seal job, and the HSM custody. The vendor is a **subservice organization** for SOC 2 purposes (see §3.6 below).

| Component | Description |
|---|---|
| SDK runtime | [INSTITUTION:fillin — language and version], hosted by the institution. |
| Audit ledger | Hosted by [INSTITUTION:fillin — vendor name]. Connectivity over OTLP/gRPC TLS 1.3. |
| HSM | Operated by [INSTITUTION:fillin — vendor]. Institution holds public-key registry entry; vendor holds private-key custody under FIPS 140-2 Level 3 (or higher). |
| Vendor SOC report | [INSTITUTION:fillin — vendor SOC 2 Type II report, period and date]. Carve-in or carve-out treatment per §3.6. |
| Vendor-conformance attestation | [INSTITUTION:fillin — current attestation ID from the working-group registry per `docs/vendor-conformance-attestation.md`]. |

## §3.3 Software (DC3 — software)

The chain implementation comprises software components governed by spec §1.3 (security definitions), §4 (the four primitives), §5 (wire format), §6 (storage), and §7 (verification).

| Software | Function | Reference |
|---|---|---|
| Chain SDK | Captures events into the per-event MAC-bound chain. | Spec §4.1 |
| Seal-job worker | Computes daily Merkle roots and obtains HSM signatures. | Spec §4.2, §4.3 |
| Verifier (CLI) | Executes the §7 verification procedure offline. | Spec §7, §10.12 |
| OTLP encoder/decoder | Wire-format binding for chain extension fields. | Spec §4.4, §5 |
| Operational-event emitter | Emits `master.*`, `seal.*`, `chain.*`, `ledger.*`, `hsm.*` events for control evidence. | Spec §10.2 |
| Conformance test harness | Re-runs the FFIEC conformance corpus against the institution's deployment. | Spec §8, `spec/test-vectors/` |

The institution's deployed versions, by service:

| Service | Version | Source-of-truth identifier |
|---|---|---|
| Chain SDK | [INSTITUTION:fillin — semver] | [INSTITUTION:fillin — Git commit hash, container image digest, or signed artifact reference] |
| Seal-job worker | [INSTITUTION:fillin — semver] | [INSTITUTION:fillin — source-of-truth identifier] |
| Verifier | [INSTITUTION:fillin — semver] | [INSTITUTION:fillin — binary SHA-256, SLSA attestation reference] |
| Audit ledger schema | [INSTITUTION:fillin — schema version per institution's release log] | [INSTITUTION:fillin — migration manifest reference] |

The verifier binary's SHA-256 and the SLSA attestation reference are part of the evidentiary artifacts retained per spec §10.13 and are cited from the `verifier.run_completed` operational event for each verifier run during the period.

## §3.4 People (DC4 — people)

The institution operates the system under a documented role separation. Roles and responsibilities:

| Role | Responsibility | Headcount |
|---|---|---|
| Chain operations team | Operates the SDK, ledger, seal job, and verifier. Responds to chain anomalies surfaced by §7 verification. | [INSTITUTION:fillin — count] |
| HSM administrator | Provisions HSM, manages the public-key registry, executes documented key-rotation procedures. **Does NOT have operator-role `sign` access** per spec §10.5 separation of duties. | [INSTITUTION:fillin — count] |
| Seal-job operator | Holds `sign`-only authorization on the HSM. Triggers the daily seal job; cannot extract or modify keys. | [INSTITUTION:fillin — count] |
| Tenant key registrar | Provisions tenant IKMs into the registry per spec §10.1 (uniqueness enforcement). | [INSTITUTION:fillin — count] |
| Internal audit / SOC engagement | Tests controls per `docs/audit-procedures.md` and consumes verifier output as control evidence. | [INSTITUTION:fillin — count] |
| Incident response | Owns disposition of chain-detected anomalies per `docs/incident-response-playbook.md`. | [INSTITUTION:fillin — count] |
| Model risk management | Reviews algorithm-rotation evidence and dual-algorithm posture per spec §4.3.2 and §10.10.2. | [INSTITUTION:fillin — count] |

The institution's hiring, screening, and termination procedures applicable to these roles are documented in [INSTITUTION:fillin — HR policy reference]. Background checks are performed at hire per [INSTITUTION:fillin — institution policy reference].

## §3.5 Data (DC5 — data)

The data captured, processed, and retained by the system:

| Data class | Description | Retention |
|---|---|---|
| Chain events | Per-event JSON records under per-event MAC. Cover prompts, responses, tool calls, routing decisions, operational state. | [INSTITUTION:fillin — typically 7 years for FFIEC chain-of-custody data] |
| Daily seal records | One record per tenant-day with Merkle root and HSM signature. | Same as chain events. |
| Tenant IKM material | 32-byte secret per tenant. Held in HSM/KMS. | Per spec §10.9: at least as long as any chain entry stamped with the IKM's `key_version` is retained. |
| Public-key registry | Tenant public-key fingerprints and HSM-signed public keys. | Indefinite (historical anchors required for past chain verification). |
| Operational events | `master.*`, `seal.*`, `chain.*`, `ledger.*`, `hsm.*` events per spec §10.2. | Same as chain events per spec §10.2. |
| Verifier output | Three-line normative format per spec §7. Retained as working-paper evidence. | Same as chain events. |
| Audit-file headers | One per persisted audit file. Carry `format_version`, `genesis_hash`, `tenant_id`, `hkdf_inputs_digest`. | Same as chain events. |

Data classification:

- Chain entries MAY contain PII when the AI's prompts or responses include customer information. The institution's data-classification policy applies; the chain does NOT introduce a separate classification regime.
- Tenant IKM material is institution-secret; cleartext export is prohibited.
- Public keys, fingerprints, and seal signatures are institution-confidential but designed to be safely shared with examiners under witness-mode verification per spec §7.

## §3.6 Subservice organizations (DC8 — subservice organizations and complementary subservice-organization controls)

Use the subsection that matches the topology; delete the others.

### §3.6.A No subservice organizations (self-hosted only)

The institution operates the entire chain-of-custody system. There are no subservice organizations relevant to this report.

### §3.6.B Subservice organization — vendor-hosted ledger and HSM (carve-out treatment)

[INSTITUTION:fillin — vendor name] is a subservice organization for the audit-ledger hosting and HSM custody services. The carve-out method is used in this description: the vendor's controls are NOT included in the description of the system, and the service auditor does NOT test those controls.

The institution relies on the following Complementary Subservice Organization Controls (CSOCs):

| CSOC | Description |
|---|---|
| CSOC-VND-01 | Vendor maintains FIPS 140-2 Level 3 (or higher) HSM custody for the institution's signing key with non-extractable private-key storage per spec §10.5. |
| CSOC-VND-02 | Vendor enforces seal-job operator role with `sign`-only authorization; separation of duties from HSM administrator. |
| CSOC-VND-03 | Vendor operates the audit ledger under append-only role grants per spec §10.3. |
| CSOC-VND-04 | Vendor's SDK / seal-job software passes the FFIEC conformance corpus per `docs/vendor-conformance-attestation.md`; current attestation in the registry. |
| CSOC-VND-05 | Vendor maintains TLS 1.3 (or §5.1-compliant) transport between the institution's SDK and the vendor-hosted ledger. |
| CSOC-VND-06 | Vendor's SOC 2 Type II report covers the period and is reviewed by the institution per CUEC-VND-06 (see §3.8 below). |
| CSOC-VND-07 | Vendor publishes a public-key URL for vendor-conformance-attestation validation per `docs/vendor-conformance-attestation.md`. |

The institution's procedures for monitoring CSOC effectiveness are described in §3.8 (CUECs) and §3.9 (control activities).

### §3.6.C Subservice organization — vendor-hosted ledger and HSM (carve-in treatment)

[INSTITUTION:fillin — vendor name]'s controls over audit-ledger hosting and HSM custody are included in this description and tested by the service auditor. The vendor's relevant controls appear in §3.9 (control activities) alongside the institution's controls. The vendor's SOC 2 Type II report is referenced in [INSTITUTION:fillin — appendix reference] and supports the carve-in scope.

## §3.7 Period-end cutoff procedure (DC6 — period end and changes during the period)

The Type II report covers the period from [INSTITUTION:fillin — period start date] through [INSTITUTION:fillin — period end date].

The chain operates continuously with seal cadence per spec §4.2.1 (daily, hourly, or weekly). Because seal cadence is not aligned to report-period boundaries, the period-end cutoff requires explicit treatment so coverage is unambiguous.

The institution applies the following cutoff procedure:

1. **Latest verifiable seal date.** The latest seal date for which the verifier produces a PASS within the report-issuance window is the **cutoff seal date**. The cutoff seal date is on or before [INSTITUTION:fillin — period end date].
2. **Seal cadence.** The institution operates [INSTITUTION:fillin — daily / hourly / weekly] cadence per spec §4.2.1. Under [daily / hourly] cadence, the cutoff seal date is typically equal to the period end date. Under weekly cadence, the cutoff seal date may be up to seven days before the period end date.
3. **Events captured after the cutoff seal date but before the period end date.** Events in this window are captured under per-event MAC and are retained at the ledger but are NOT yet bound under a signed Merkle root. The Type II opinion's coverage of these events is conditional on the next seal completing as part of normal operations after the report-issuance date. The institution discloses this window's event count in [INSTITUTION:fillin — supplementary schedule].
4. **Late-binding entries crossing the cutoff.** Late-binding entries (`ffiec.chain.late_binding = true` per spec §4.2.2) captured after the cutoff seal date but related to events from before the cutoff are reported as anomalies in the verifier output. The institution's normal late-binding-rate baseline applies (see audit procedure P-31 baseline).
5. **Cutoff verifier output.** A verifier run executed on or after [INSTITUTION:fillin — period end + 7 days] over the period's tenant-days produces the cutoff verifier output. The output is retained as working-paper evidence per spec §10.13.

The cutoff procedure is invoked by the institution's chain operations team and the SOC engagement team jointly. The procedure's output is the cutoff verifier output, the supplementary schedule of post-cutoff captured events, and a reconciliation of the period's seal records against expected cadence (one seal per tenant-day for daily cadence; 24 seals per tenant-day for hourly cadence; one seal per tenant-week for weekly cadence; empty-day seals included where no events occurred per spec §4.2).

## §3.8 Complementary user-entity controls (DC9 — complementary user-entity controls)

The Service Organization's controls were designed with the assumption that certain controls are operated by user entities. The following CUECs are necessary to achieve the control objectives:

| CUEC | Description | Cross-reference |
|---|---|---|
| CUEC-01 | User entity (tenant) maintains the confidentiality of any tenant-side application credentials used to write events through the SDK. | spec §4.1 — session-key custody |
| CUEC-02 | User entity provides the SDK with the correct `tenant_id` per spec §3 character class. The SDK rejects non-conforming `tenant_id` values; the user entity is responsible for the value's correctness. | spec §3, §3.1 |
| CUEC-03 | User entity defines and operates the data-classification policy applicable to chain entries that may contain PII or other regulated content. | spec §3.5 (this document) |
| CUEC-04 | User entity reviews verifier output for PASS-with-anomaly findings (late-binding entries, cadence anomalies, structural-only PASS under witness mode) and dispositions them per its IR procedure. | spec §7, §10.10 |
| CUEC-05 | User entity selects a seal cadence per spec §4.2.1 that aligns with the report period. | §3.7 (this document) |
| CUEC-VND-06 | User entity validates the vendor-conformance attestation annually per `docs/vendor-conformance-attestation.md` and archives the validation log. | `docs/vendor-conformance-attestation.md` |
| CUEC-07 | User entity operates incident response per its standard procedure when chain-detected anomalies trigger the institution's notification clocks (FFIEC 36-hour, GDPR 72-hour, DORA 4-hour, NIS2 24-hour, applicable as the institution's regulatory perimeter requires). | `docs/incident-response-playbook.md` |
| CUEC-08 | User entity composes the chain's integrity-of-recording claim with SR 11-7 model validation when chain outputs feed ICFR or financial-disclosure processes. | `soc-pack/sox404-icfr-composition.md` |

The institution's procedures for confirming user-entity controls are operating effectively are described in §3.9 (control activities).

## §3.9 Control activities — Trust Services Criteria mapping (DC7 — control activities)

The institution's controls map to the Trust Services Criteria (TSP 100) as follows. Each control is implemented by the chain primitives, the operational-event emission, the verifier, or the institution's surrounding processes; the table summarizes the mapping for the description-of-system reader and points to the Section 4 control activities for testing detail.

### Common Criteria (CC1 — Control Environment)

| Criterion | Implementation |
|---|---|
| CC1.1 — COSO principle 1 | Institution's code of ethics covers chain operators; named in [INSTITUTION:fillin — ethics policy reference]. |
| CC1.2 — Board independence | [INSTITUTION:fillin — board structure reference]. |
| CC1.3 — Reporting structure | Chain operations report to [INSTITUTION:fillin — CISO / CTO / CRO]. Separation from HSM administration documented in §3.4. |
| CC1.4 — Workforce competence | [INSTITUTION:fillin — training program reference]. Operators of chain components hold [INSTITUTION:fillin — required certifications]. |
| CC1.5 — Accountability | Annual performance review covers chain-related responsibilities for in-scope roles. |

### CC2 — Communication and Information

| Criterion | Implementation |
|---|---|
| CC2.1 — Information quality | The chain's verifier output (three-line normative format per spec §7) is the load-bearing communication of chain integrity status. |
| CC2.2 — Internal communication | Chain anomalies surface via the operational events listed in spec §10.2 to [INSTITUTION:fillin — observability platform / SIEM]. |
| CC2.3 — External communication | Verifier output and `audit-procedures.md` procedures govern external (examiner) communication. |

### CC3 — Risk Assessment

| Criterion | Implementation |
|---|---|
| CC3.1 — Risk objectives | Chain integrity is a documented risk objective in [INSTITUTION:fillin — risk register reference]. |
| CC3.2 — Risk identification | Threat model documented in `docs/design/09-threat-model.md`; reviewed annually. |
| CC3.3 — Fraud risk | Adversaries A through I per `09-threat-model.md` cover the named fraud-relevant scenarios. |
| CC3.4 — Significant change | Spec amendments and SDK / verifier version changes trigger control reassessment; spec §12 change log feeds the institution's change-management calendar. |

### CC4 — Monitoring Activities

| Criterion | Implementation |
|---|---|
| CC4.1 — Ongoing evaluation | Daily seal-job logs (§10.13), `master.reconciliation_completed` events (§10.1), and verifier runs constitute ongoing evaluation. |
| CC4.2 — Deficiency communication | Chain-detected anomalies escalate per `incident-response-playbook.md`; deficiencies communicated to [INSTITUTION:fillin — committee]. |

### CC5 — Control Activities

Implemented through the spec §4 primitives and §10 operational requirements. See Section 4 for testing.

### CC6 — Logical and Physical Access

| Criterion | Implementation |
|---|---|
| CC6.1 — Logical access | Append-only enforcement at two layers per spec §10.3. Database role grants INSERT and SELECT only on chain tables. |
| CC6.2 — Account provisioning | [INSTITUTION:fillin — IAM platform] provisions chain operator accounts. Annual recertification per [INSTITUTION:fillin — recertification calendar]. |
| CC6.3 — Logical access modification | Append-only at the role level per spec §10.3. UPDATE / DELETE / TRUNCATE permissions REVOKE'd. |
| CC6.4 — Physical access | HSM physical access controls per [INSTITUTION:fillin — data-center policy / cloud-provider physical-security attestation]. |
| CC6.5 — Logical access removal | Termination procedure removes chain operator accounts within [INSTITUTION:fillin — SLA]. |
| CC6.6 — External-party access | Examiner witness-mode verification per spec §7 supports external access without exposing the IKM. |
| CC6.7 — Data transmission | OTLP transport over TLS 1.3 per spec §5.1. |
| CC6.8 — Malicious software | Software-key adapter exclusion per spec §10.7 (compile-time / packaging exclusion). |

### CC7 — System Operations

| Criterion | Implementation |
|---|---|
| CC7.1 — Detection of anomalies | Chain `verification_failure` events per spec §10.2; `master.reconciliation_completed` per §10.1. |
| CC7.2 — Monitoring of system components | Operational events listed in §10.2 surface to [INSTITUTION:fillin — observability platform]. |
| CC7.3 — Evaluation of security events | IR playbook scenarios cover chain-detected anomalies per `incident-response-playbook.md`. |
| CC7.4 — Incident response | 36-hour FFIEC cyber-incident clock applicable to chain-detected anomalies per `regulator-pack/breach-notification-matrix.md`; institution's IR commander operates the tightest applicable clock first. |
| CC7.5 — Recovery from incidents | DR procedures per `docs/dr-and-resilience.md`; multi-region resilience patterns per spec §10.15. |

### CC8 — Change Management

| Criterion | Implementation |
|---|---|
| CC8.1 — Change management | Changes to chain SDK, seal-job worker, verifier, or HSM configuration follow [INSTITUTION:fillin — change-management procedure]. The institution's CC8.1 control description names: (a) the IKM-registry layer, (b) the software-key adapter exclusion mechanism per spec §10.7, (c) the vendor-conformance attestation consumption per `docs/vendor-conformance-attestation.md`, (d) the dual-algorithm posture per spec §4.3.2 and §10.10.2 if applicable. |

### CC9 — Risk Mitigation

| Criterion | Implementation |
|---|---|
| CC9.1 — Risk mitigation | Retention per spec §10.9 (IKM retention coupling) and §10.13 (evidentiary artifacts). |
| CC9.2 — Vendor-management controls | Vendor-conformance attestation procedure per `docs/vendor-conformance-attestation.md`. Subservice organization treatment per §3.6 above. |

### Processing Integrity (PI1)

| Criterion | Implementation |
|---|---|
| PI1.1 — Inputs complete and accurate | Per-event MAC at capture (spec §4.1) covers the canonical-form input bytes; cross-tenant lift detected at §7 step 4. |
| PI1.2 — Inputs processed completely | Daily Merkle seal (spec §4.2) ensures no insertion, deletion, or reordering passes verification. Empty-day seals close the seal-sequence gap per spec §4.2. |
| PI1.3 — Outputs delivered as intended | Verifier output (spec §7) is the load-bearing delivery confirmation; three-line normative format. |
| PI1.4 — Outputs traceability | Cross-run linkage fields (`parent_run_id`, `parent_seq`, `dag_parents`) are integrity-bound per spec §5. |
| PI1.5 — Authority and responsibility | Operator-role separation per spec §10.5; CC1.3 reporting structure. |

### Availability (A1)

| Criterion | Implementation |
|---|---|
| A1.1 — Capacity | [INSTITUTION:fillin — capacity-planning reference]. Spec §10.15 multi-region resilience supports availability commitments. |
| A1.2 — Recovery | Multi-region pattern per spec §10.15 (Pattern A or Pattern B); DR procedures per `docs/dr-and-resilience.md`. |
| A1.3 — Recovery testing | Annual DR test per [INSTITUTION:fillin — DR-testing calendar]; chain-integrity verification post-recovery is the load-bearing recovery confirmation. |

### Confidentiality (C1)

| Criterion | Implementation |
|---|---|
| C1.1 — Confidential information identified | Tenant IKM material classified institution-secret; chain entries inherit the institution's data-classification policy per §3.5. |
| C1.2 — Disposal | IKM retirement per spec §10.9; chain-entry retention per §3.5. |

### Privacy (P)

When the institution claims Privacy:

| Criterion | Implementation |
|---|---|
| P1 — Notice | [INSTITUTION:fillin — privacy-notice URL]. The chain itself does not surface privacy notices; institution's surrounding process does. |
| P2 — Choice | [INSTITUTION:fillin — consent-management platform reference]. |
| P3 — Collection | Chain entries collected only when the AI's prompts or responses contain regulated content; collection is incidental to the AI workflow. |
| P4 — Use, retention, disposal | Retention per §3.5; disposal per IKM retirement and chain-entry retention exit. |
| P5 — Access | Data-subject access procedures per [INSTITUTION:fillin — DSAR procedure]. Witness-mode verification supports access without exposing the IKM. |
| P6 — Disclosure to third parties | Examiner disclosure per `docs/litigation-support.md`; vendor-attested chain entries do not transit to third parties without institution consent. |
| P7 — Quality | Chain-integrity verification supports data-quality assertions about recorded form, NOT about substantive accuracy (epistemic scope per spec §1.2). |
| P8 — Monitoring and enforcement | `master.reconciliation_completed` events per spec §10.1; verifier runs per §7. |

When the institution does NOT claim Privacy: this subsection is deleted in the final report; institution's separate privacy-program attestation covers the criterion.

## §3.10 Sampling population definition (DC6 — population coverage)

For each control activity the SOC engagement tests, the population is unambiguously defined. The table below is the institution's population definition by procedure family per `docs/audit-procedures.md`.

| Procedure | Population universe | Period coverage |
|---|---|---|
| P-6 — IKM-roster reconciliation | Every `(tenant_id, key_version, key_fingerprint)` triple observed during the period. | Period start through cutoff seal date. |
| P-22 — Anomaly disposition sample | Every `chain.verification_failure` event during the period. | Period start through cutoff seal date. |
| P-25 — Stratified chain-entry sample | Every chain entry written during the period, stratified by `chain_kind` per spec §3 enumeration. | Period start through cutoff seal date. |
| P-30 — Seal-record sample | (Tenant-day count × in-scope tenant count) for the period. Empty-day seals included. | Period start through cutoff seal date. |
| P-31 — Truncation-baseline sample | Every chain entry where `ffiec.chain.late_binding = true` during the period; institution's documented baseline is the comparison. | Period start through cutoff seal date. |
| P-37 — Multi-region replication sample | Every `master.cross_region_replication_completed` event during the period; cross-checked against seal-region event counts. | Period start through cutoff seal date. |

For Type II reports, "period start" and "cutoff seal date" replace "as of" semantics. The period coverage rule is stated explicitly so the auditor's testing is unambiguous and comparable across institutions.

## §3.11 Boundaries of the system (DC1 — boundaries)

The system's boundary, for purposes of this Section 3 description:

**In-scope.** The chain SDK, the audit ledger, the seal-job worker, the verifier, the tenant key registry, the public-key registry, the operational-event emission for the events listed in spec §10.2, and the institution's processes for change management, incident response, vendor management (where applicable), and SOC engagement consumption of verifier output.

**Out-of-scope.**

- LLM provider operations (the substantive AI inference). The chain captures what the institution received from the provider; it does NOT attest the provider's training-data integrity, model-card claims, or provider-side controls.
- Application logic that consumes the AI's response for downstream business decisions. The chain proves what the AI said and that the record is intact; it does NOT prove the application's interpretation or use of the AI's response complies with policy.
- Customer-facing channels that deliver the AI-derived output (web, mobile, contact center). Those channels operate under separate controls; chain entries provide the authoritative record of the AI's contribution.
- Statistical evaluation of decision populations for fair-lending, ECOA, or other compliance questions. Those are model-validation activities under SR 11-7 and are scoped separately.

The institution's IT witness, expert witness, examiner, and SOC engagement team rely on the boundary statement for scope determination. The boundary follows the epistemic-scope discipline of spec §1.2.

## §3.12 Significant changes during the period (DC6 — significant changes)

The institution discloses any significant change to the system during the period:

| Change | Date | Description | Impact on controls |
|---|---|---|---|
| [INSTITUTION:fillin — change description] | [INSTITUTION:fillin — date] | [INSTITUTION:fillin — description] | [INSTITUTION:fillin — control impact: e.g., "no impact on PI1 controls; CC8.1 change-management evidence retained"] |

If no significant changes occurred, this subsection is replaced with: "No significant changes occurred during the period."

## §3.13 Forward-looking elements (DC6 — subsequent events)

The institution discloses any condition that may affect the system's continued operation in a manner the SOC report's reader would consider material:

- [INSTITUTION:fillin — pending corpus-version transitions per `docs/vendor-conformance-attestation.md`].
- [INSTITUTION:fillin — pending algorithm rotations under spec §4.3.2].
- [INSTITUTION:fillin — vendor-conformance attestation expiring during the report-issuance window].
- [INSTITUTION:fillin — other forward-looking element].

If none, this subsection is replaced with: "No forward-looking conditions material to the report's reader."

## §3.14 Management's assertion (boilerplate)

The institution's management asserts that:

- The description of the system fairly presents the chain-of-custody system that was designed and operated during the period.
- The controls within the system were suitably designed throughout the period to meet the applicable Trust Services Criteria.
- The controls within the system operated effectively throughout the period to meet the applicable Trust Services Criteria — specifically, [INSTITUTION:fillin — list TSC categories: Security, Availability, Processing Integrity, Confidentiality, Privacy].
- The complementary user-entity controls and complementary subservice-organization controls assumed in the design are necessary for the controls to meet the applicable criteria.

The assertion is signed by [INSTITUTION:fillin — name], [INSTITUTION:fillin — title], on behalf of the institution.

```
Signed: ___________________________________
        [INSTITUTION:fillin — name]
        [INSTITUTION:fillin — title]
        [INSTITUTION:fillin — date]
```

## Appendix A — Cross-reference to spec sections

| Section in this template | Spec reference |
|---|---|
| §3.1 services in scope | spec §1.1, §1.2, §1.4 |
| §3.2 infrastructure | spec §4, §5.1, §10.5, §10.16 (SaaS-edge connectors when applicable) |
| §3.3 software | spec §4.1, §4.1.3 (Round-17 NIST-P1 algorithm agility), §4.2, §4.3, §4.4, §4.4.5 (Round-17 NAIC underwriting), §5, §6, §7, §10.2, §10.22 (Round-17 CFPB-P2 redaction discipline) |
| §3.4 people | spec §10.5, §10.6.1, §10.17 (HSM partition ceremony attestation, Round-17) |
| §3.5 data | spec §3, §10.9, §10.13, §10.20 (training-data retention floor, Round-17), §10.22 (pre-MAC redaction) |
| §3.6 subservice organizations | `docs/vendor-conformance-attestation.md`, spec §10.21 (cross-vendor model-handover), §10.24 (entity succession, Round-17 M&A-G1) |
| §3.7 period-end cutoff | spec §4.2.1, §4.2.2, §10.15 |
| §3.8 CUECs | spec §3, §4.1, §10.10, §10.15, §10.18 (CC8.1 / runbook cross-referencing), §10.19 (chain-coverage map) |
| §3.9 control activities | spec §1.4, §10.1 through §10.24 |
| §3.10 sampling populations | `docs/audit-procedures.md` |
| ECOA / FCRA workflow controls | spec §10.11, §10.11.1 (ECOA reasons schema), §10.11.2 (FCRA reinvestigation), §10.23 (consumer-correlation index) |

## Appendix B — How this Section 3 supports specific TSC opinions

| Opinion | Load-bearing references |
|---|---|
| Processing Integrity (PI1.1, PI1.2) | Per-event MAC (spec §4.1), daily Merkle seal (spec §4.2), HSM signature on root (spec §4.3); cross-tenant binding via HKDF info parameter; constant-time comparison discipline (spec §10.8); OPTIONAL second-algorithm MAC for quantum-readiness (spec §4.1.3, RECOMMENDED at v1.0b). |
| Security (CC6) | Append-only enforcement at two layers (spec §10.3); software-key adapter exclusion (spec §10.7); HSM custody at FIPS 140-2 Level 3 (spec §10.5); IKM minimum length and CSPRNG provenance (spec §10.6, §10.6.1); HSM partition ceremony attestation with signatory `entity_affiliation` (spec §10.17). |
| Availability (A1) | Multi-region resilience patterns (spec §10.15); operational events for ledger and seal-job lifecycle (spec §10.2); SaaS-edge connector lag SLO + RTO + alerting threshold (spec §10.16, when applicable). |
| Confidentiality (C1) | IKM custody (spec §10.5); chain-entry retention coupling (spec §10.9); pre-MAC redaction discipline with `audit.redaction.*` family (spec §10.22). |
| Privacy (when claimed) | Witness-mode verification supports DSAR without IKM exposure (spec §7); ECOA translation entry attribute schema (spec §10.11) supports preferred-language disclosure; ECOA adverse-action reasons schema (spec §10.11.1); FCRA reinvestigation event family (spec §10.11.2); pre-MAC redaction discipline (spec §10.22). |
| Vendor management (CC9.2) | Cross-vendor model-handover schema with `audit.model_handover.*` (spec §10.21); training-data retention floor (spec §10.20); chain-coverage map version-stamped and chain-anchored (spec §10.19); CC8.1 / runbook cross-referencing (spec §10.18); entity succession (spec §10.24). |

## Appendix C — Aoife G-1 / G-3 / G-5 / G-7 closure summary

The template is structured so the institution issues a SOC 2 Type II report without manufacturing description language. Specifically:

- **G-1 (Section 3 description).** §3.1 through §3.13 supply the description-of-system content required by AT-C §205. Three topology variants in §3.2 cover self-hosted, BYOC, and vendor-hosted deployments.
- **G-3 (materiality framework).** Population definitions in §3.10 are the load-bearing replacement for institution-discretion thresholds. Auditor evaluates findings against the defined populations.
- **G-5 (sampling population definition).** §3.10 names the population universe per procedure unambiguously.
- **G-7 (period-end cutoff).** §3.7 describes the cutoff procedure so Type II coverage of post-cutoff captured events is unambiguous.

## Summary

This Section 3 template gives the institution a SOC 2 Type II description-of-system that addresses the Trust Services Criteria, defines the system's boundary, names the subservice organization treatment, identifies the CUECs the user entity must operate, and stipulates the period-end cutoff procedure so coverage is unambiguous. The template pairs with `soc-pack/section-4-template.md` for control-activity testing language and `soc-pack/sox404-icfr-composition.md` for SOX 404 composition when the chain feeds disclosure-relevant outputs.

The institution paste-completes the marked fields, deletes the topology variants that don't apply, and pairs with the Section 4 template before issuance. The result is a description that a Big-4 attestation partner can defend at peer review without negotiating description language per engagement.
