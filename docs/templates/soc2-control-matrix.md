---
status: normative-template
alignment-reference: AICPA SSAE 18 / AT-C §105 / AT-C §205 / TSP 100 (2017) — Trust Services Criteria; AICPA SOC for Cybersecurity; AICPA SOC for AI Systems (2025 draft); ISAE 3000 / ISAE 3402 (international); AT-C §205 §A21 (management-prepared evidence reliability); AT-C §205 §A47 (sampling); AT-C §205 §A52 (deficiency disposition)
companion-templates:
  - docs/templates/soc2-section3-description-of-system.md (Section 3 description; consume per-control mapping in §C.1)
  - docs/audit-procedures.md (institution-side P-1..P-N procedures; auditor re-performs on a sub-sample)
date: 2026-05-07
version: 1.0.0
---

# SOC 2 Type II — Control Matrix and Auditor Procedures

> **What this template is.** A control-by-control matrix the SOC firm uses on a chain-of-custody attestation engagement. The matrix complements `soc2-section3-description-of-system.md` (the institution-side Section 3 narrative) by giving the SOC firm a per-Trust-Service-Criterion (TSC) mapping, an auditor-side procedure for each control, evidence-sampling methodology, and disposition guidance for findings. The template is structured for direct paste into the SOC firm's engagement workpapers and the SAR (Security Assessment Report) where applicable.
>
> **Audience.** SOC 2 / ISAE 3000 / ISAE 3402 attestation partners and engagement teams. The institution-side equivalent is `audit-procedures.md`.
>
> **Scope discipline.** The matrix does not modify the spec or institution-side controls. It is auditor-engagement substrate: per-control mapping, sampling methodology, and disposition rubric the engagement consumes once and reuses across institutions.

## How to use this template

1. **Per-TSC mapping (§A).** The chain prescribes controls that map to multiple Trust Service Criteria. The SOC firm's engagement workpapers cite this mapping rather than rewriting it per engagement.
2. **Auditor-side procedures (§B).** The institution operates the P-N procedures in `docs/audit-procedures.md`. The SOC firm re-performs on a sub-sample using the procedures in §B and draws conclusions per AT-C §205 §A21.
3. **Sampling methodology (§C).** §C names the sampling unit, frequency, and stratification for each test type.
4. **Subservice-organization scoping (§D).** §D names the carve-out vs inclusive method for cloud HSM and ledger services.
5. **Disposition rubric (§E).** §E names the disposition (control performed effectively, control deficiency, significant deficiency, material weakness) per finding pattern.
6. **Multi-framework cross-walks (§F).** §F maps each chain-prescribed control to ISAE 3000 / 3402, FedRAMP, HITRUST CSF, and PCI DSS for engagements that combine frameworks.
7. **SOC for AI Systems alignment (§G).** §G maps chain capabilities to the AICPA SOC for AI Systems criteria for engagements consuming the chain as substrate for AI-system attestation.

---

## §A. Trust Service Criteria — control-by-control mapping

The chain's normative controls map to multiple TSCs. The mapping below is the engagement-baseline. Engagements may refine for specific scope, but the baseline is consistent across institutions.

### §A.1 Cryptographic and key-management controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §10.5 | HSM custody (FIPS 140-2 Level 3 or higher) | CC6.1 (logical access) | CC6.2 (system resources protection); CC6.3 (network communications); CC6.6 (encryption); CC6.7 (physical access) |
| §10.6.1 | RNG provenance attestation | CC6.1 (cryptographic key management) | CC6.6 (encryption controls) |
| §10.6 | IKM length and entropy | CC6.6 | CC6.1 |
| §4.1 | Per-event HMAC binding | CC6.6 | CC4.1 (monitoring) |
| §4.3 | Daily Merkle seal + HSM signature | CC6.6 | CC4.1; AU-9 (federal cross-walk) |
| §10.10 | Algorithm rotation | CC8.1 (change management) | CC6.6; CC4.2 (deficiency remediation) |
| §10.7 | Software-key adapter exclusion in production | CC8.1 | CC6.1 (logical access — preventing dev-key path); SI-7 (integrity, federal cross-walk) |

### §A.2 Integrity and append-only controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §4.1 inviolate property | Append-only chain-storage enforcement | CC6.1 (logical access — preventing modification) | CC8.1 (change management); CC7.2 (system monitoring) |
| §10.3 | Database role restrictions (P-1) | CC6.1 | CC8.1 |
| §10.8 | Constant-time comparison | CC6.6 | CC7.2 |
| §10.4 | Time synchronization (NTP) | CC7.5 (system operations) | CC7.2 |
| §4.2 | Daily seal coverage (every tenant-day) | CC4.1 | CC7.2 |

### §A.3 Reconciliation and monitoring controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §10.1 (P-6) | Key-fingerprint reconciliation | CC4.1 (monitoring) | CC4.2 (deficiency remediation); CC7.2 (anomaly detection) |
| §10.1 | Count reconciliation | CC4.1 | CC7.2 |
| §10.1 | Anomaly-detection cadence | CC7.2 | CC4.1 |
| §10.2 (P-13) | Internal-audit verifier-run cadence | CC4.1 | CC4.2 |

### §A.4 Verifier and assessment controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §7 | Verifier procedure (12-step) | CC4.1 | CC4.2 |
| §10.12 | Verifier exit-code contract | CC4.1 | A1.1 (availability — verifier reproducibility) |
| §10.13 | Verifier-build trust path (cosign + GPG) | CC8.1 | CC6.6 |
| §7 | Reproducible-build verification (P-12) | CC8.1 | CC4.1 |

### §A.5 Incident-response controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| IR Playbook Scenarios 1–11 | IR playbook documentation and execution | CC7.4 (incident response) | CC4.2 (remediation) |
| 12 CFR §53 / 36-hour rule | Cyber-incident notification posture | CC2.3 (communication and operations) | CC7.4 |
| IR Scenario 3 | Seal-publication delay notification (72 hours) | CC2.3 | CC7.4 |

### §A.6 Configuration and change-management controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §4.1.2 | Posture choice (FFIEC-conformant vs vendor-namespaced) | CC8.1 | CC2.2 (communication of objectives) |
| §10.10.2 | Within-day algorithm rotation | CC8.1 | CC6.6; CC4.2 |
| §4.3.2 | Algorithm rotation across the day boundary | CC8.1 | CC6.6 |

### §A.7 Retention and confidentiality controls

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §10.13 | Evidentiary artifact retention (7-year floor) | CC9.1 (risk mitigation through information retention) | C1.1 (confidentiality categorization); A1.2 (availability) |
| §10.9 | IKM retention coupling | C1.1 | CC9.1 |
| §1.2 (epistemic scope) | Privacy-by-design composition with chain artifacts | P1.1 (privacy notice) | P5.1 (privacy retention); C1.2 (confidentiality access) |

### §A.8 Vendor-management controls (when vendor-hosted)

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §10.5 vendor scoping | Vendor SOC report consumption (P-17) | CC9.2 (vendor risk management) | CC8.1 |
| §10.5 vendor scoping | Vendor IAM boundary (P-18) | CC6.1 | CC9.2 |
| §10.13 vendor scoping | Container image trust path (P-19) | CC8.1 | CC6.1 |
| §10.16 | SaaS-edge connector lag SLO + alerting threshold + RTO | CC4.1 (monitoring) | CC9.2; CC7.2 (anomaly detection) |
| §10.17 | HSM partition ceremony attestation (signatory `entity_affiliation` per Round-17 M&A-P1) | CC6.1 | CC8.1; CC2.3 |
| §10.18 | Runbook ↔ spec-section cross-referencing | CC2.2 (communication) | CC8.1 |
| §10.19 | Chain-coverage map version-stamped + chain-anchored (Round-17 M&A-P3) | CC2.2 | CC9.2; CC8.1 |
| §10.20 | Training-data retention floor ≥ longest deployment window + investigation buffer | CC9.1 | C1.1; CC9.2 |
| §10.21 | Cross-vendor model-handover schema (`audit.model_handover.*`, `training_shard_manifest_sha256` per Round-17 M&A-P2) | CC9.2 | CC8.1 |
| §10.22 | Pre-MAC redaction discipline (`audit.redaction.*` family per Round-17 CFPB-P2) | CC6.6 | C1.1; CC4.1 |

### §A.9 ECOA / FCRA adverse-action workflow controls (Round-17 CFPB)

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §10.11 | ECOA / state-DOI translation entry schema (`audit.ecoa.translation.*`) | CC2.3 (communication) | CC4.1 |
| §10.11.1 | ECOA adverse-action reasons schema (`audit.ecoa.adverse_action.*` per Round-17 CFPB-P1) | CC2.3 | CC8.1; CC4.1 |
| §10.11.2 | FCRA reinvestigation event family (`audit.fcra.reinvestigation.*` per Round-17 CFPB-G1) | CC2.3 | CC4.1; CC7.4 |
| §10.23 | Consumer-correlation index integrity (Shape 1 chain-anchored OR Shape 2 daily attestation; Round-17 CFPB-G2) | CC4.1 | CC9.1; C1.1 |

### §A.10 Underwriting / disparate-impact controls (Round-17 NAIC; state-DOI engagements)

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §4.4.5 | `audit.underwriting.features.*` REQUIRED on model-driven underwriting / triage / pricing decisions (Round-17 NAIC-P1) | CC4.1 | CC8.1 |
| §4.4.5 | `audit.disparate_impact.*` per-period DI test artifacts (Round-17 NAIC-P2; RECOMMENDED at v1.0b) | CC4.1 | CC4.2 |

### §A.11 MAC algorithm agility and entity succession (Round-17 NIST-P1, M&A-G1)

| Spec section | Control | Primary TSC | Secondary TSCs |
|---|---|---|---|
| §4.1.3 | OPTIONAL `payload_hash_alt` second-algorithm MAC + `ffiec.chain.algorithm_alt` (RECOMMENDED at v1.0b for quantum-readiness; institution names alt algorithm in CC8.1) | CC6.6 | CC8.1; CC9.1 |
| §10.24 | Entity succession event (`chain.entity_succession` with dual-signatory binding per Round-17 M&A-G1) | CC8.1 | CC2.3 |

---

## §B. Auditor-side procedures — re-performance and re-verification

The institution operates the P-N procedures in `docs/audit-procedures.md`. The SOC firm executes the procedures below to produce attestation evidence.

### §B.1 RNG-source attestation procedure (auditor-side)

For the §10.6.1 RNG-provenance control, the SOC firm:

1. **Sample the most-recent `master_key.generated` operational event for each tenant in scope.** Confirm the event records the RNG source identifier per spec §10.6.1.
2. **Verify the institution's CC8.1 documentation names the same RNG source.** A mismatch between the operational event and CC8.1 is a control-design deficiency (not just an evidence gap).
3. **Request the vendor's FIPS validation certificate.** For HSM-internal RNG: the HSM's FIPS 140-2 Level 3 (or higher) validation certificate. For OS-level CSPRNG (rare in production; not conformant for spec §10.5 deployments): the OS distribution's CSPRNG documentation.
4. **For HSM-internal RNG:** request the HSM's RNG self-test logs from the previous quarter.
5. **For OS-level CSPRNG:** request the institution's host-hardening procedure documentation.
6. **Reconcile across (1)–(5).** Discrepancies are control deficiencies the institution remediates per CC4.2.

Evidence retention: the SOC firm's working paper records the RNG source identifier, the FIPS certificate ID, the self-test log range, and the reconciliation result.

### §B.2 Verifier-output evidence — three postures

§7 verifier output is byte-for-byte reproducible (§10.12 exit-code contract). The SOC firm chooses one of three postures based on engagement scope, evidence-reliability needs, and cost.

#### §B.2.1 Institution-run posture

The institution runs the verifier and provides the output to the SOC firm. AT-C §205 §A21 requires the auditor to evaluate management-prepared evidence reliability:

1. The SOC firm samples a subset of institution-run outputs (typical: 20% of the institution's verifier runs during the period, OR 30 outputs minimum).
2. The SOC firm re-performs the verifier on the same input ledger artifacts.
3. The SOC firm compares the re-performed exit-code and failure-reason granularity against the institution's reported output.
4. Mismatches are control deficiencies (the institution's verifier output is unreliable; the engagement's reliance on management-prepared evidence is undermined).

Posture cost: low. Posture reliability: lowest of the three (relies on institution's verifier-run integrity).

#### §B.2.2 Auditor-run posture

The SOC firm runs the verifier directly against the chain artifacts the institution produces.

1. The institution provides the verifier binary (with cosign + GPG verification per `verifier-validate.sh`).
2. The institution provides the tenant's public key.
3. The institution provides read access to the chain ledger for the engagement period.
4. The SOC firm runs the verifier on the agreed sample.

Posture cost: medium. Posture reliability: high (verifier is the auditor's controlled tool).

#### §B.2.3 Independent-third-party-run posture

A neutral third party (per spec §7 witness verifier) runs the verifier and produces an attestation.

Posture cost: high. Posture reliability: highest. Used in highest-stakes engagements (combined regulatory exam + SOC 2 opinion).

#### §B.2.4 Posture documentation

The institution's CC8.1 control description names which posture the institution operates under. The SOC firm's evaluation procedure is fixed by this template per §B.2.1–§B.2.3 above.

### §B.3 §10.1 reconciliation testing — sampling methodology

The institution operates P-6 (key-fingerprint reconciliation), count reconciliation, and anomaly-detection per spec §10.1.

#### §B.3.1 Sampling unit

The tenant-day reconciliation record. Each unit includes:
- The `master.reconciliation_completed` operational event.
- The ledger-side count for that tenant-day.
- The fingerprint-roster snapshot at reconciliation time.
- Any anomalies detected and their disposition.

#### §B.3.2 Sampling frequency

Tiered by tenant risk:

| Risk tier | Frequency |
|---|---|
| Census-tier (high-risk per institution risk register; new tenants in first 90 days) | Every reconciliation in the period |
| Standard | Statistical sample of 60 tenant-days per period (95% confidence; 5% tolerable misstatement) |
| Backstop | All reconciliations from the most recent 30 days before report-issuance |

The combined sample produces both period-coverage evidence (statistical) and currency evidence (backstop).

### §B.4 §10.7 software-key adapter exclusion — auditor test procedure

The institution's CC8.1 names the exclusion pattern. The SOC firm tests:

1. **Preferred procedure (highest reliability):** sample 10% of production-image pulls from the institution's registry. Inspect each image's content manifest (cosign-verified) and confirm absence of the adapter assembly, package, or module.
2. **Alternative procedure (medium reliability):** inspect the institution's deployment-pipeline output and the cosign-verified manifest. Less reliable than direct image inspection because the manifest signing is upstream of registry storage.
3. **Fallback procedure (low reliability):** code-search against the production source tree. Only acceptable when the institution does not store production images centrally.
4. **Inadequate procedure:** rely on CC8.1 documentation alone. Insufficient under AT-C §105.

Evidence retention: working paper records the procedure used, sample size, and result.

### §B.5 §10.5 HSM custody — auditor procedure

1. Confirm the HSM is FIPS 140-2 Level 3 or higher.
2. Confirm the HSM operator role list shows seal-job operator separated from HSM administrator (or, for institutions with documented dual-control compensating control, that the dual-control is operational).
3. Sample HSM audit logs for the period; confirm key-load events match change-management records.
4. Confirm no extraction events appear in the HSM audit log (the IKM is non-extractable).
5. For cloud HSM (AWS CloudHSM, Azure Managed HSM, Google Cloud HSM): scope per §D (subservice-organization).

### §B.6 §10.13 evidentiary artifact retention — auditor test procedure

1. **Sampling unit:** one tenant-day's full evidence package. Each package includes: chain entries for the day, daily seal record, fingerprint-roster snapshot, KMS access log, RNG attestation reference.
2. **Frequency:** sample of 30 tenant-days per period, biased toward recent days (10 from the trailing 30 days, 10 from the trailing 90 days, 10 from earlier in the period).
3. **For each sampled package:** confirm each artifact is present, retained per §10.13 (7-year floor or institution's longer policy), and read-only / append-only at the storage layer.
4. **Disposition:** a missing artifact is a control deficiency in CC9.1; multiple missing artifacts are a significant deficiency.

### §B.7 §4.1 inviolate property — append-only enforcement test

1. Confirm database role restrictions per P-1.
2. Sample-test by attempting an UPDATE against `events` and `daily_seals` tables under the ledger writer role; confirm rejection.
3. Confirm the institution's storage-layer immutability flag (e.g., S3 Object Lock, Azure immutable blob, on-prem WORM volume) is set on retained artifacts.
4. Sample-test by attempting a DELETE against retained artifacts; confirm rejection.

### §B.8 §10.4 time synchronization — auditor procedure

1. Pull the institution's NTP-sync evidence per P-7.
2. Sample-test by checking time drift on a representative host.
3. Confirm host clock-drift alerts fire when drift exceeds the institution's documented threshold.
4. Sample 10 chain entries; confirm `captured_at` timestamps are within tolerance of the host clock at capture time.

---

## §C. Sampling methodology — engagement-wide

### §C.1 Sampling table

The SOC firm consumes the same sampling table the institution uses (`audit-procedures.md` §Sampling), with the auditor-side floor of 30 entries per stratum for testability.

| Population size | Sampling unit | Sample size |
|---|---|---:|
| < 100 events | Per event | All (census) |
| 100–500 | Per event | 30 |
| 500–2,500 | Per event | 50 |
| > 2,500 | Per event | 75 |
| Reconciliation records | Per tenant-day | See §B.3.2 |

### §C.2 Stratification

The auditor stratifies by:
- **Tenant.** At least 5 entries per tenant.
- **Time.** At least one sample date per month of the period.
- **Day-of-week.** At least 3 weekdays observed.
- **Around control changes.** Additional entries from days adjacent to master-key rotations, posture changes.
- **Around incidents.** 50–100 supplementary entries from any IR-program-flagged period.

### §C.3 Auditor sample-design documentation

Per AT-C §205 §A47, the working paper records:
- Random-number seed.
- Sampling tool (e.g., R `sample()`, Python `random.sample()`, statistical software).
- Distribution of exit codes and failure reasons across the sample.
- Stratification proof (e.g., the count of entries per stratum).

### §C.4 Census handling

For high-risk strata (rare event types, recent control changes, incident-flagged periods): test the entire stratum rather than sampling. Document the census determination in the workpaper.

---

## §D. Subservice-organization scoping

### §D.1 Default posture — carve-out method

When the institution uses a cloud HSM (AWS CloudHSM, Azure Managed HSM, Google Cloud HSM) or vendor-hosted ledger service, the cloud provider is a subservice organization under AICPA TSP. The default scoping posture is the **carve-out method**:

- The cloud provider's controls are explicitly excluded from the SOC 2 scope.
- The SOC report names the cloud provider's complementary user-entity controls (CUECs).
- The SOC report points to the cloud provider's own SOC 2 / ISAE 3402 / FedRAMP attestation.

#### §D.1.1 Required complementary user-entity controls (CUECs)

When carve-out applies, the institution operates the following CUECs:
- **Cryptographic key management (CKM) policy:** the institution's policy for IKM provisioning, custody, rotation, and retirement.
- **Key-rotation procedure:** the institution's procedure for executing IKM rotation per spec §10.10.
- **HSM access logs review:** the institution's monthly review of HSM access logs for unexpected access patterns.
- **Subservice-organization SOC report review on receipt:** the institution's documented review procedure for the cloud provider's SOC 2 / ISAE 3402 report.
- **CUEC enumeration:** explicit list of CUECs in Section 3 of the SOC 2 report (per `soc2-section3-description-of-system.md`).

### §D.2 Inclusive method (less common)

The cloud provider's controls are included in the institution's SOC 2 scope by explicit reference. Used when:
- The institution has contractual visibility into the cloud provider's chain-related controls.
- The engagement is sized to include cloud-provider control testing (substantial scope expansion).

The inclusive method requires the cloud provider's cooperation and substantially expands engagement scope and cost. It is rare for cloud HSM in commercial SOC engagements.

### §D.3 Documenting the choice

The institution's CC8.1 names which method applies. The SOC firm's working paper records the choice and the supporting rationale.

---

## §E. Disposition rubric — chain-related findings

Per AT-C §205 §A52, each finding is dispositioned as one of:
- Control performed effectively.
- Control deficiency (inconsequential).
- Significant deficiency.
- Material weakness.

### §E.1 Single MAC-mismatch on a tenant-day

| Pattern | Disposition |
|---|---|
| Surfaced by institution's §10.1 reconciliation; remediated within the period | Control performed effectively (the reconciliation control caught the issue) |
| Not surfaced by institution's reconciliation (auditor caught during sampling) | Control deficiency in §10.1 reconciliation |
| Multiple MAC-mismatches across multiple tenant-days; not surfaced | Significant deficiency |

### §E.2 Missing seal record for a tenant-day with zero events

Per §4.2, every tenant-day requires a seal even with zero events. Missing empty-day seal is reported as `missing seal for tenant-day {D}` per §10.12.

| Pattern | Disposition |
|---|---|
| Single occurrence; institution's IR Scenario 3 activated | Control deficiency (control-completeness, not chain-integrity) |
| Recurring occurrences | Significant deficiency |
| Recurring with no IR Scenario 3 activation | Material weakness |

### §E.3 `dev_mode=true` seal in production

Per §10.7, this is the "double-protection regulator-visible line."

| Pattern | Disposition |
|---|---|
| Single occurrence | Significant deficiency |
| Recurring | Material weakness |
| Recurring without remediation plan | Material weakness + adverse opinion |

### §E.4 Algorithm rotation crossing tenant-day boundary without §10.10 procedure compliance

| Pattern | Disposition |
|---|---|
| Operational impact bounded; remediation in progress | Control deficiency |
| Operational impact significant; remediation incomplete | Significant deficiency |
| Operational impact significant + chain integrity compromised | Material weakness |

### §E.5 Late-binding entries above threshold

Per §4.2.2, late-binding is normal. High-volume late-binding may indicate operational deficiency.

| Pattern | Disposition |
|---|---|
| Within institution's documented threshold | No finding |
| Above threshold; institution's IR Scenario 3 activated | Control deficiency (operational-only) |
| Above threshold without IR activation | Control deficiency in IR-program operating effectiveness |

### §E.6 Verifier failure rate

| Failure rate | Disposition |
|---|---|
| 0% | Control performed effectively |
| > 0% but < 0.1% | Control deficiency |
| 0.1% – 1% | Significant deficiency |
| > 1% | Material weakness |
| Structural (negative-test-vector match) | Material weakness + adverse opinion |

### §E.7 Three-lines-of-defense alignment

| Pattern | Disposition |
|---|---|
| All three lines explicitly documented and operating | Control performed effectively |
| Documentation gaps but operational effectiveness shown | Control deficiency |
| One or more lines unassigned | Significant deficiency |
| Independence between first/third line broken | Material weakness |

---

## §F. Multi-framework cross-walks

Each chain-prescribed control maps to control objectives in adjacent frameworks. Engagements combining multiple frameworks consume the cross-walk to avoid duplicate evidence collection.

### §F.1 SOC 2 (TSC) ↔ ISAE 3000 / 3402

ISAE 3000 (assurance engagements other than audits or reviews of historical financial information) and ISAE 3402 (assurance reports on controls at a service organization) use control objectives rather than TSC. The mapping:

| TSC | ISAE 3000 / 3402 control objective |
|---|---|
| CC6.1 (logical access) | Access control objective |
| CC6.6 (encryption) | Cryptographic key management objective |
| CC4.1 (monitoring) | Monitoring and detection objective |
| CC8.1 (change management) | Change management objective |
| CC7.4 (incident response) | Incident response objective |
| CC9.1 (information retention) | Records retention objective |
| C1.1 (confidentiality) | Confidentiality objective |
| A1.1 (availability) | Availability objective |
| P1.1 (privacy notice) | Privacy notice objective |

### §F.2 SOC 2 ↔ FedRAMP (NIST SP 800-53 Rev 5)

FedRAMP uses NIST SP 800-53 Rev 5 controls for federal authorization. The mapping appears in `docs/regulator-pack/fedramp-fisma-overlay.md`. SOC 2 + FedRAMP combined engagements consume the same evidence from `audit-procedures.md`; the auditor's working paper documents both attestations.

### §F.3 SOC 2 ↔ HITRUST CSF

HITRUST CSF (healthcare) consumes SOC 2 evidence with healthcare-specific overlays (HIPAA Security Rule, HITECH). The mapping:

| TSC | HITRUST CSF objective |
|---|---|
| CC6.1 / CC6.6 | 01 — Access Control; 06 — Information Security Incident Management |
| CC4.1 | 06 — Audit Logging and Monitoring |
| CC8.1 | 09 — Change Management |
| CC9.1 | 09 — Records and Information Management |

### §F.4 SOC 2 ↔ PCI DSS

PCI DSS (payment card industry) consumes SOC 2 evidence with cardholder-data-specific overlays. The mapping:

| TSC | PCI DSS Requirement |
|---|---|
| CC6.1 / CC6.6 | Req 3 (Protect stored cardholder data); Req 4 (Encrypt transmission) |
| CC4.1 | Req 10 (Track and monitor access) |
| CC8.1 | Req 6 (Develop and maintain secure systems) |
| CC9.1 | Req 10.5 (Retention) |

### §F.5 Multi-framework evidence reuse

The institution produces a single body of evidence (per `audit-procedures.md`). The SOC firm consumes it through the framework-specific lens. The chain's verifier output, operational events, reconciliation records, and HSM logs are framework-neutral substrate.

---

## §G. SOC for AI Systems alignment (AICPA 2025 draft)

AICPA's SOC for AI Systems (currently in 2025 field-test under AICPA Audit and Attest Standards Board) introduces AI-specific control criteria. The chain is designed precisely as AI logging infrastructure, so it maps cleanly to several SOC for AI Systems criteria.

### §G.1 Mapping (informative; subject to AICPA final publication)

| SOC for AI Systems criterion | Chain attribute / event |
|---|---|
| Data quality at inference | Captured prompt content, retrieval context, RAG document IDs and content-hashes (per spec §4.4) |
| Model performance monitoring | `audit.model.performance_report_generated` (institution-emitted; spec §4.4 audit.* namespace) |
| Model bias monitoring | Model-call chain entries plus institution-side bias-evaluation events |
| Model drift detection | `audit.deployment.intent` distribution per spec §4.4.2 (canary, A/B test, multi-region drift) |
| Model lifecycle management | `audit.deployment.policy_version`; spec §10.10 algorithm rotation |
| Model-related access controls | `gen_ai.request.model`, `gen_ai.response.model`, model-call provenance |

### §G.2 Posture for SOC for AI Systems engagements

The chain provides the integrity substrate. The AI-system-specific controls (data quality, bias, drift) are operated by the institution and feed events into the chain. The SOC for AI Systems engagement consumes the chain entries as evidence; the engagement's procedures verify the institution's AI-system controls, with the chain providing tamper-evident attestation that the events were not altered after capture.

### §G.3 Update cadence

This section evolves with AICPA's SOC for AI Systems final publication. The mapping is informative; engagements consuming the SOC for AI Systems framework refine the mapping per the published criteria when AICPA finalizes the standard.

---

## §H. Engagement workpaper template

The following workpaper sections are produced for each chain attestation engagement.

### §H.1 Engagement scope

- Period covered.
- Tenants in scope.
- Topology (self-hosted, BYOC, vendor-hosted).
- Subservice-organization scoping (carve-out vs inclusive).
- Verifier-output posture (institution-run, auditor-run, third-party-run).

### §H.2 Per-control test results

For each control in §A:
- Procedure executed (referencing §B).
- Sample size and stratification.
- Re-performance results (PASS/FAIL distribution).
- Findings and disposition (referencing §E).

### §H.3 Subservice organization controls

- Carve-out CUECs (if carve-out method).
- Inclusive controls tested (if inclusive method).
- Vendor SOC reports reviewed.

### §H.4 Findings register

For each finding:
- Description.
- Disposition per §E.
- Institution's remediation plan.
- Reporting in SOC report Section 5 (auditor's tests and results).

### §H.5 Opinion conclusion

- Type 1 vs Type 2 distinction.
- Unqualified, qualified, adverse, or disclaimer.
- Subsequent events through report date.

---

## §I. Disposition guide for the SOC report

### §I.1 Section 1 — Independent service auditor's opinion

Standard SOC 2 Type II opinion language, modified to reference the chain controls. The institution's CC8.1 attributes named in §A are cited in the engagement's description.

### §I.2 Section 2 — Management assertion

The institution asserts the chain operates per the spec and per the controls in §A. The CC8.1 attributes are part of the assertion.

### §I.3 Section 3 — Description of the system

Use `soc2-section3-description-of-system.md` template. The description references this matrix for the per-control mapping.

### §I.4 Section 4 — Trust Services Criteria, controls, and results

For each TSC in §A, list the controls and the auditor's tests with results. The matrix in §A is the foundation; the §B procedures are the tests.

### §I.5 Section 5 — Other information

Any subsequent events, scope limitations, and management responses to findings.

---

## §J. Engagement-cost guidance

| Institution category | Engagement effort (hours) |
|---:|---:|
| Community bank (< $10B) | 80–120 |
| Regional bank ($10B–$250B) | 200–300 |
| Tier-1 (> $250B) | 400–800 (multi-region) |

The matrix and procedure templates reduce per-engagement effort by 20–30% vs writing engagement procedures from scratch on first chain engagement. Subsequent engagements at the same institution amortize the matrix further.

---

## §K. AICPA peer review readiness

The matrix and procedures support AICPA peer review:
- Per-TSC mapping (§A) is consistent across engagements at different institutions.
- Sampling methodology (§C) is documented and reproducible.
- Disposition rubric (§E) produces consistent finding severity across engagements.
- Subservice-organization scoping (§D) follows AICPA TSP defaults.

Engagement teams using this matrix and the institution-side `audit-procedures.md` produce SOC reports that read consistently across firms and institutions, supporting peer review and downstream consumer reliance.
