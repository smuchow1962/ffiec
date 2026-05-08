---
status: normative-template
alignment-reference: EBA/GL/2019/02 (Guidelines on outsourcing arrangements) §§9, 13, 14; DORA Article 30(2)(d) and (h); DORA Articles 28-30 (ICT third-party risk management); GDPR Article 28(3); EBA/GL/2019/02 §§4.7, 4.12-4.13
companion-templates:
  - docs/regulator-pack/gdpr-controller-vs-processor.md (GDPR DPA template)
  - docs/vendor-conformance-attestation.md (vendor-conformance attestation procedure)
date: 2026-05-07
version: 1.0.0
---

# EBA outsourcing and DORA audit-rights clause (template)

> **What this template is.** A drop-in contractual clause set the institution attaches to vendor agreements covering chain-of-custody hosting, audit-ledger services, HSM custody, or any other ICT third-party arrangement under EBA/GL/2019/02 and DORA. The template is structured for direct paste into a Master Services Agreement, Statement of Work, or stand-alone Data Processing and ICT Services Agreement. The template establishes competent-authority audit rights, on-site inspection rights, delegated-audit rights, sub-outsourcing pre-approval, exit strategies, replication of EU-resident data, and register-of-information attestations.

## How to use this template

1. **Identify the regulatory perimeter.** The institution operates under one or both of:
   - **EBA/GL/2019/02** (Guidelines on outsourcing arrangements) — applies to credit institutions, investment firms, payment institutions, and electronic-money institutions in the Union for outsourcing of "critical or important functions."
   - **DORA** (Regulation (EU) 2022/2554) — applies to financial entities (per Article 2(1)) for ICT third-party services arrangements (per Articles 28-30).
   The institution's legal team confirms which framework applies; in most cases for chain-of-custody arrangements, both apply.
2. **Pick the clause variants.** Each clause has alternatives marked `[VARIANT-A]`, `[VARIANT-B]` where the institution selects per its posture. Unselected variants are deleted before contract execution.
3. **Fill in the marked fields.** Every `[INSTITUTION:fillin]` token names a value the institution provides. Tokens are explicit so a contracts reviewer can scan for unfilled gaps before execution.
4. **Normative vs explanatory language.** Each clause has a **Normative text** block (the contract language) and an **Explanatory note** block (the institution's understanding of why the clause exists). Only the normative text appears in the executed contract; the explanatory note supports the institution's contracts team and the regulator's joint examination team in confirming the clause's intent.
5. **Pair with the GDPR DPA template.** GDPR Article 28(3) elements (purpose, duration, types of personal data, obligations of the processor) are covered by the DPA template at `docs/regulator-pack/gdpr-controller-vs-processor.md`. This template adds the EBA and DORA elements that GDPR Article 28(3) does not address.

## Preamble — definitions and scope (use once, at the top of the clause set)

### Normative text

> **Definitions.** For the purposes of this Agreement and the clauses below:
>
> - "Critical or Important Function" has the meaning given in EBA/GL/2019/02 §1.6.
> - "ICT Services" has the meaning given in DORA Article 3(21).
> - "Sub-Outsourcing" has the meaning given in EBA/GL/2019/02 §1.7 and includes any further outsourcing of services by the Service Provider to a third party.
> - "Competent Authority" means the supervisory authority responsible for the Institution under EU or national financial-services regulation, including the European Banking Authority, the European Central Bank under the Single Supervisory Mechanism, the national prudential regulator, the data-protection authority, and any other authority with statutory access rights to the Institution's records under applicable law.
> - "Sub-Contractor" means any third party engaged by the Service Provider to perform any part of the Services.
> - "EEA" means the European Economic Area.
> - "Member State" means a member state of the European Union.
> - "Adequacy Decision" means a decision under GDPR Article 45 by which the European Commission has determined a third country provides an adequate level of data protection.

### Explanatory note

The definitions above import the EBA and DORA terms unchanged. The institution's contracts team uses these definitions for clauses §1 through §10 below; the institution's GDPR DPA template uses GDPR-specific definitions (Controller, Processor, Personal Data) and the two definition sets do not overlap.

## §1 — Description of services and dependencies (DORA Article 30(2)(a))

### Normative text

> 1.1 The Service Provider shall describe the Services in writing, including the nature of the Services, the Critical or Important Function (if any) that the Services support, and the dependencies of the Services on other parties (Sub-Contractors, software vendors, hardware vendors, telecommunications providers).
>
> 1.2 The description of the Services shall include the locations from which the Services are performed (the "Service Locations"), enumerated by Member State or third country, by data centre or processing facility, by jurisdiction of incorporation of the Service Provider's legal entity providing the Services. The description shall be updated within thirty (30) calendar days of any change in Service Location.
>
> 1.3 The description of the Services shall identify, for each Service Location, the legal regime governing the data and operations at that location, and any local-law obligation (including obligations of disclosure to non-EEA authorities) that may affect the Institution or the Service Provider's compliance with this Agreement.

### Explanatory note

DORA Article 30(2)(a) requires the contract to specify a "complete description of the services and functions to be provided." EBA/GL/2019/02 §13.1 requires the contract to specify the locations from which the services are provided. The two requirements are unified in §1.

The local-law disclosure requirement in §1.3 is the institution's hook for Schrems II supplementary-measures assessment: a Service Location in a non-adequate jurisdiction obliges the institution to evaluate the risk that local-law disclosure obligations conflict with EU data-protection requirements.

## §2 — Performance and service-level monitoring (DORA Article 30(2)(b) and (e))

### Normative text

> 2.1 The Service Provider shall provide the Services in accordance with documented service levels, including availability, integrity, confidentiality, latency, and recovery objectives. The service levels are specified in [INSTITUTION:fillin — SLA appendix reference].
>
> 2.2 The Service Provider shall report to the Institution monthly on actual service-level performance against the documented service levels. Reports shall be delivered within ten (10) business days of the end of each month and shall include:
>    (a) measured availability per Service;
>    (b) measured integrity verification outcomes (if applicable to the Services), including any failures of the chain-of-custody verification procedure under FFIEC chain-of-custody Specification §7;
>    (c) any incidents that affected service levels during the reporting month, with root-cause analysis;
>    (d) any outstanding remediation actions and their target completion dates.
>
> 2.3 The Institution may at any time, on reasonable notice, request additional service-level information beyond the monthly report, including operational telemetry, log records, and configuration evidence relevant to the Institution's supervisory or regulatory obligations. The Service Provider shall respond within five (5) business days.

### Explanatory note

DORA Article 30(2)(b) requires the contract to specify "the agreed service levels … in clear, complete, and enforceable terms." Article 30(2)(e) requires the Service Provider to "monitor the performance of contractual obligations … and provide the financial entity with reports."

The integrity-verification reporting in §2.2(b) is the chain-specific instrumentation: when the chain-of-custody Specification §7 verifier reports anomalies, those findings are the load-bearing service-level evidence.

## §3 — Information security (DORA Article 30(2)(c); EBA/GL/2019/02 §14)

### Normative text

> 3.1 The Service Provider shall implement and maintain technical and organizational measures appropriate to the risk of the Services, in accordance with:
>    (a) Article 32 of Regulation (EU) 2016/679 (the General Data Protection Regulation), where the Services involve the processing of personal data;
>    (b) DORA Articles 5-13 (ICT risk-management framework) and the Regulatory Technical Standards adopted thereunder;
>    (c) EBA/GL/2019/02 §14 (security of data and systems);
>    (d) for HSM custody, FIPS 140-2 Level 3 or higher (or the EU Common Criteria or eIDAS QSCD equivalent named in [INSTITUTION:fillin — institution security-policy reference]); and
>    (e) the FFIEC chain-of-custody Specification v1.0a/b §10 operational requirements, where the Services include hosting or operating chain-of-custody components, including without limitation: §10.16 SaaS-edge capture connectors (lag SLO, alerting threshold, RTO), §10.17 HSM partition ceremony attestation (signatory `entity_affiliation` per Round-17 M&A-P1), §10.18 CC8.1 and runbook cross-referencing, §10.19 chain-coverage boundary documentation (version-stamped and chain-anchored per Round-17 M&A-P3), §10.20 training-data retention vs deployment-window discipline, §10.21 cross-vendor model-handover schema, and §10.22 pre-MAC redaction discipline (Round-17 CFPB-P2). [honors §10.16, §10.17, §10.18, §10.19, §10.20, §10.21, §10.22]
>
> 3.2 The Service Provider shall provide the Institution, on the Institution's request, with documented evidence of the technical and organizational measures, including (i) certifications under ISO/IEC 27001, ISO/IEC 27017, ISO/IEC 27018, SOC 2 Type II, or equivalent; (ii) penetration-testing reports performed within the prior twelve (12) months; (iii) configuration baselines for systems supporting the Services; and (iv) any vendor-conformance attestation issued under the FFIEC chain-of-custody attestation procedure.
>
> 3.3 The Service Provider shall maintain ICT security awareness for its personnel performing the Services, including documented training programmes, role-based access reviews on at least an annual cadence, and incident-response training.

### Explanatory note

DORA Article 30(2)(c) requires the contract to specify the security measures the Service Provider implements. EBA/GL/2019/02 §14.4 requires the institution to ensure the Service Provider's information security is "appropriate to the assessed level of risk."

The chain-of-custody-specific requirement in §3.1(e) imports the spec's normative §10 operational requirements (HSM custody, append-only enforcement, IKM length, constant-time comparison, software-key adapter exclusion) as contractual obligations.

### Technical Notes — round-17 §10 obligations the Service Provider inherits

- **§10.16 SaaS-edge connectors.** When the Service Provider operates a mirror connector against a SaaS platform (CRM, helpdesk, ticketing, email), the contract obliges the Service Provider to quantify the median lag, the 95th-percentile lag (the lag SLO), the alerting threshold, and the recovery-time objective. Imprecise wording — "near real-time," "low-latency mirror" — is non-conformant. The Service Provider's runbook MUST cite the quantified bounds by number.
- **§10.17 HSM partition ceremony attestation.** Signatories on each ceremony record their `entity_affiliation` (per Round-17 M&A-P1) so the chain can distinguish a signer's authority across an entity-change boundary. The PDF attendance log retains its dispute-resolution role; the chain entry is the integrity-bound attestation that the ceremony occurred.
- **§10.18 CC8.1 and runbook cross-referencing.** Every runbook section that supports a normative spec requirement carries the spec section number that obligates it. A runbook section on multi-region failover names §10.15; a runbook section on HSM custody names §10.5; absence is a discoverability gap surfaced as a Nit by the SOC engagement team.
- **§10.19 chain-coverage map.** The Service Provider's coverage map is version-stamped (`coverage_map_version`, `effective_utc`) and chain-anchored via `chain.coverage_map_published` operational events, so an 18-month-lookback acquirer-side IT due-diligence team can determine which version was in force on a given date.
- **§10.20 training-data retention floor.** Where the Service Provider acts as a model provider, training-data shards are retained for the longest active deployment window of any model trained on the data, plus an investigation buffer (typically 60-90 days). A 90-day retention against an 18-month deployment window is a partial conformance.
- **§10.21 model-handover schema.** Where the Service Provider delivers a model to the Institution, the handover entry binds `model_artifact_sha256`, `model_card_sha256`, `fairness_audit_report_sha256`, and (where applicable) `training_shard_manifest_sha256` per Round-17 M&A-P2.
- **§10.22 pre-MAC redaction discipline.** Where the Service Provider's Services include redaction at the SDK boundary, redaction runs before the per-event MAC seals the canonical bytes. Post-MAC sidecar redaction is non-conformant unless paired with a §10.21-style cross-anchor.

## §4 — Audit rights and on-site inspections (EBA/GL/2019/02 §13.2; DORA Article 30(3))

### Normative text

> 4.1 The Institution and any auditor appointed by the Institution (including any external auditor and any auditor appointed by a Competent Authority) shall have the right, on reasonable notice, to:
>    (a) audit the Service Provider's compliance with this Agreement, the Service Provider's policies and procedures, and the Service Provider's technical and organizational measures;
>    (b) review the Service Provider's records, including transaction logs, audit trails, configuration records, change-management records, incident-response records, and the chain-of-custody records produced under the FFIEC chain-of-custody Specification (where applicable);
>    (c) conduct on-site inspections at any Service Location, including the Service Provider's offices, data centres, and any premises from which the Services are performed;
>    (d) interview the Service Provider's personnel performing the Services;
>    (e) test the Service Provider's controls, including testing of the Service Provider's response to chain-of-custody anomalies surfaced by the verification procedure under FFIEC Specification §7; and
>    (f) receive copies of any relevant document, record, or report prepared by or on behalf of the Service Provider.
>
> 4.2 The audit and inspection rights in §4.1 are unrestricted as to scope, frequency, location, and method, subject only to (i) reasonable advance notice (not less than ten (10) business days for routine audits; immediate access where the Competent Authority requires it for the exercise of its supervisory functions or where a security-incident-response action requires it); (ii) compliance by the auditor with reasonable confidentiality obligations; and (iii) reasonable accommodation by the Service Provider of operational continuity (the audit shall not unreasonably disrupt the Service Provider's operations).
>
> 4.3 The Service Provider shall provide the Institution and any auditor appointed by the Institution with full cooperation, including access to the Service Provider's premises, devices, systems, networks, information, and data, and shall make available any documentation, records, or personnel reasonably requested.
>
> 4.4 The Competent Authority of the Institution, and any auditor appointed by such Competent Authority, shall have the same rights as set out in §4.1 and §4.3, exercisable directly against the Service Provider without requiring the Institution's prior consent. The Service Provider acknowledges that the Competent Authority may request such access at any time and shall not require the Institution to obtain consent before granting it.
>
> 4.5 The audit and inspection rights in §4.1 through §4.4 apply equally to any Sub-Outsourcing arrangement under §6 below. The Service Provider shall ensure all Sub-Contractors agree to the same audit and inspection obligations, and the Service Provider shall facilitate the Institution's and the Competent Authority's exercise of those rights against any Sub-Contractor.
>
> 4.6 The Service Provider shall not charge the Institution or the Competent Authority any fee for the exercise of audit or inspection rights under §4.1 through §4.5.

### Explanatory note

EBA/GL/2019/02 §13.2 specifically requires the contract to grant the institution and "any auditor appointed by it, as well as the competent authority, all the rights of access … and audit rights" with respect to "the service provider's premises … and any other relevant devices, systems, networks and information used for providing the services outsourced." DORA Article 30(3) further requires the contract to ensure the Competent Authority has "the right of access to relevant business premises … and to be able to exercise that right effectively."

The Competent Authority direct-access provision in §4.4 is the load-bearing element. EBA/GL/2019/02 §13.2 requires the Service Provider to grant the Competent Authority access without requiring the Institution to mediate; a contract that requires Institution consent for Competent Authority access is non-conformant.

The chain-of-custody-specific testing reference in §4.1(e) supports the regulator's joint examination team in independently testing the chain's verification procedure during an inspection visit. The chain-of-custody verifier's offline, no-network-call discipline (Specification §7) supports this testing without Service Provider cooperation beyond providing the ledger artifacts.

## §5 — Delegated-audit and pooled-audit rights (EBA/GL/2019/02 §13.3)

### Normative text

> 5.1 The Institution may exercise its audit and inspection rights under §4 through:
>    (a) the Institution's internal audit function;
>    (b) an external auditor appointed by the Institution;
>    (c) a third-party auditor jointly appointed by the Institution and other financial entities receiving services from the Service Provider (a "Pooled Audit");
>    (d) a third-party auditor appointed by a Competent Authority; or
>    (e) the Competent Authority directly.
>
> 5.2 The Service Provider shall accept Pooled Audits and shall not require duplicative audits of the same matter within any twelve (12) month period unless the Institution has documented grounds for an additional audit (including findings from a prior audit, an incident, a regulatory request, or a material change in the Services).
>
> 5.3 The Service Provider shall accept the report or output of a Pooled Audit, an external-auditor audit, or a Competent-Authority-appointed audit as satisfying the Institution's contractual audit right for the matter covered by the audit; the Service Provider shall not require the Institution to perform an additional, duplicative audit of the same matter.
>
> 5.4 Where the Service Provider's existing third-party reporting (such as a SOC 2 Type II report, an ISO/IEC 27001 certificate, or a vendor-conformance attestation under the FFIEC chain-of-custody attestation procedure) covers the matter the Institution would otherwise audit, the Institution may rely on the existing reporting in lieu of conducting a separate audit, provided the existing reporting is current, the Institution has reviewed it, and the Institution's records of the review are maintained per the Institution's standard control-evidence retention.

### Explanatory note

EBA/GL/2019/02 §13.3 endorses pooled audits and reliance on third-party reporting as a means of avoiding duplicative auditing of the same matter across multiple financial entities served by a single Service Provider. The vendor-conformance attestation procedure at `docs/vendor-conformance-attestation.md` is one of the third-party reporting paths §5.4 contemplates.

The 12-month limitation in §5.2 is calibrated so the Institution can audit annually as a baseline cadence while reserving the right to perform incident-driven or finding-driven audits more frequently.

## §6 — Sub-outsourcing (EBA/GL/2019/02 §§13.4, 4.12-4.13; DORA Article 30(2)(a))

### Normative text

> 6.1 The Service Provider shall not Sub-Outsource any part of a Critical or Important Function or any ICT Service supporting a Critical or Important Function without the Institution's prior written consent. The Institution's consent shall not be unreasonably withheld but may be conditioned on the Service Provider's demonstration of:
>    (a) the proposed Sub-Contractor's compliance with the security and operational requirements of this Agreement;
>    (b) the proposed Sub-Contractor's acceptance of the audit and inspection rights of §4 (including direct rights of the Competent Authority);
>    (c) the proposed Sub-Contractor's location of processing and any local-law obligations (per §1.3); and
>    (d) the impact of the Sub-Outsourcing on the Service Provider's exit strategy under §9 below.
>
> 6.2 The Service Provider shall maintain a register of all Sub-Outsourcing arrangements and Sub-Contractors. The register shall include, for each Sub-Contractor, (i) the legal name and jurisdiction of incorporation; (ii) the Service Locations from which the Sub-Contractor performs the Services; (iii) the scope of the Sub-Outsourcing; (iv) the date of the Sub-Outsourcing; and (v) any material change to the Sub-Outsourcing arrangement.
>
> 6.3 The Service Provider shall notify the Institution of any proposed change to the Sub-Outsourcing arrangements (addition of a Sub-Contractor, change of Sub-Contractor's Service Location, change of scope, termination of a Sub-Outsourcing arrangement) at least thirty (30) calendar days before the change takes effect. The Institution may, within those thirty (30) calendar days, object to the change on reasonable grounds; if the Institution objects, the Service Provider shall not implement the change without resolving the Institution's concern, including (where appropriate) by exercising the termination right of §10 below.
>
> 6.4 The Service Provider shall ensure that all Sub-Outsourcing arrangements include provisions equivalent to §3 (information security), §4 (audit and inspection rights, including direct Competent-Authority access), §6.1 through §6.3 (further Sub-Outsourcing), §8 (incident notification and cooperation), and §9 (exit strategy). The Service Provider shall provide the Institution with copies of those provisions on request.
>
> 6.5 Where a Sub-Contractor's Service Location is in a third country (a country outside the EEA), the Service Provider shall ensure the Sub-Outsourcing arrangement satisfies the Institution's applicable cross-border-transfer obligations under GDPR Articles 44-50 and any successor framework. Where the third country is not the subject of an Adequacy Decision, the Service Provider shall implement the supplementary measures the Institution requires (per [INSTITUTION:fillin — institution Schrems II supplementary-measures policy reference]).

### Explanatory note

EBA/GL/2019/02 §4.12 provides that financial institutions must ensure the service provider does not sub-outsource the function without prior approval. §4.13 requires sub-outsourcing arrangements to maintain the same level of compliance as the original outsourcing.

DORA Article 30(2)(a) lists "the description of the services and functions to be provided" as a contract requirement that, as a practical matter, must include the sub-contracting chain so the financial entity has visibility into who actually performs the function.

The third-country cross-border-transfer linkage in §6.5 is the institution's hook for Schrems II supplementary-measures. A Sub-Contractor in a non-adequate jurisdiction triggers the institution's standard supplementary-measures procedure.

## §7 — Replication and resilience of EU-resident data (DORA Article 12; spec §10.15 alignment)

### Normative text

> 7.1 Where the Services include the storage or processing of personal data of EU data subjects or the operation of any critical infrastructure component supporting the Institution's Critical or Important Functions, the Service Provider shall maintain a primary processing location within the EEA or within a jurisdiction subject to an Adequacy Decision (the "EEA-Resident Storage Requirement"). Where the Services involve the operation of FFIEC chain-of-custody components for an EU-resident tenant, the seal region shall be within the EEA or within a jurisdiction subject to an Adequacy Decision.
>
> 7.2 The Service Provider shall replicate Institution data across at least two geographically separated Service Locations within the EEA (or within an Adequacy-Decision jurisdiction) to ensure resilience under DORA Article 12. The Service Provider shall document the replication topology, including the primary location, the secondary location(s), the replication latency, and the recovery-point and recovery-time objectives.
>
> 7.3 Where the Service Provider operates the FFIEC chain-of-custody under spec §10.15 multi-region resilience, the Service Provider shall:
>    (a) declare the chosen pattern (Pattern A active-active with seal-region pinning, or Pattern B per-region tenant_id) in [INSTITUTION:fillin — operational-architecture appendix reference];
>    (b) emit `master.cross_region_replication_completed` operational events per spec §10.2 for each tenant-day and source region;
>    (c) ensure the seal region for each EU-resident tenant remains within the EEA or an Adequacy-Decision jurisdiction throughout the contract term; and
>    (d) notify the Institution within twenty-four (24) hours of any unplanned change to the seal region (including failover events).
>
> 7.4 The Service Provider shall not transfer Institution data outside the EEA or an Adequacy-Decision jurisdiction without the Institution's prior written consent. Where the Institution consents to such a transfer, the transfer shall be subject to the supplementary measures the Institution requires per §6.5.

### Explanatory note

DORA Article 12 (ICT business continuity policy) requires financial entities to have ICT business-continuity arrangements proportionate to the risk. EU-resident data replication is the most direct contractual support for the institution's Article 12 obligations.

The chain-of-custody-specific provisions in §7.3 import the spec's §10.15 multi-region resilience patterns. Pattern A and Pattern B are mutually exclusive per tenant; the contract names which pattern the Service Provider operates so the Institution's CC8.1 control description and the Competent Authority's joint examination team have unambiguous understanding.

The EU-held-keys discipline (which a strict reading of EBA/GL/2019/02 §13 and Schrems II requires for vendor-hosted EU controllers) is supported by §7.1's EEA-Resident Storage Requirement.

## §8 — Incident notification and cooperation (DORA Articles 18-23; EBA/GL/2019/02 §13.5)

### Normative text

> 8.1 The Service Provider shall notify the Institution of any ICT-related incident affecting the Services without undue delay and in any event:
>    (a) within four (4) hours of the Service Provider's classification of the incident as "major" under DORA Article 18 (the "DORA Initial Notification Window"); and
>    (b) within twenty-four (24) hours of any incident affecting the integrity, availability, or confidentiality of the Services, regardless of major-incident classification.
>
> 8.2 The notification under §8.1 shall include:
>    (a) a description of the incident;
>    (b) the Services affected;
>    (c) the data and customers (or categories of data and customers) potentially affected;
>    (d) the duration of the incident (or the duration so far, where ongoing);
>    (e) the geographic spread of the incident (per Member State);
>    (f) any chain-of-custody integrity findings surfaced by the FFIEC verification procedure (per spec §7), including specific failure-reason strings;
>    (g) the Service Provider's assessment of severity and criticality;
>    (h) the immediate response actions the Service Provider is taking;
>    (i) the proposed remediation timeline; and
>    (j) any further notification the Service Provider proposes to make to its own competent authority or to law enforcement.
>
> 8.3 The Service Provider shall provide updates to the Institution on the incident at intervals not exceeding twenty-four (24) hours, until the incident is resolved and remediation is documented.
>
> 8.4 The Service Provider shall cooperate fully with the Institution's response to the incident, including by:
>    (a) providing forensic evidence relevant to the incident;
>    (b) preserving relevant logs, records, and configurations under a legal hold;
>    (c) supporting the Institution's notification obligations to its own Competent Authority, customers, and any other authority;
>    (d) supporting any digital-forensic investigation conducted by the Institution, the Competent Authority, or law enforcement; and
>    (e) participating in any post-incident review the Institution conducts.
>
> 8.5 The Service Provider shall maintain records of all incidents affecting the Services for the longer of (a) seven (7) years from the date of the incident or (b) the period required by DORA Article 19 (incident reporting record-keeping). The Service Provider shall provide such records to the Institution on the Institution's request.

### Explanatory note

DORA Article 19 sets the incident-reporting timelines. The 4-hour initial notification in §8.1(a) tracks the Regulatory Technical Standards adopted under DORA Article 18 (Commission Delegated Regulation (EU) 2024/1772). The institution's contracts team confirms the current RTS at execution and updates the clock in the contract if a subsequent regulatory amendment changes it.

The chain-of-custody-specific reporting in §8.2(f) imports the spec §7 verification procedure's normative failure-reason strings. When the chain's verifier reports a `payload_hash MAC mismatch at seq N` or any other named §7 failure, the Service Provider's incident notification includes that specific reason.

## §9 — Exit strategy and transition (DORA Article 28(7) and (8); EBA/GL/2019/02 §15)

### Normative text

> 9.1 The Service Provider shall maintain an exit strategy enabling the Institution to either bring the Services in-house or migrate them to another service provider, without disruption to the Institution's Critical or Important Functions. The exit strategy shall be documented in [INSTITUTION:fillin — exit-strategy appendix reference] and shall be:
>    (a) tested on a regular basis (not less than annually);
>    (b) reviewed and updated whenever a material change to the Services occurs; and
>    (c) provided to the Institution and to the Competent Authority on request.
>
> 9.2 On termination of this Agreement (for any reason), the Service Provider shall, during the Transition Period defined in §9.5:
>    (a) continue to provide the Services in accordance with this Agreement;
>    (b) cooperate with the Institution and any successor service provider to migrate the Services and the Institution's data;
>    (c) provide the Institution with a complete copy of all Institution data, in a format and on media specified by the Institution;
>    (d) for chain-of-custody Services, provide complete and verifiable copies of all chain entries, seal records, public keys, and operational events under the FFIEC chain-of-custody Specification, in the spec's wire format (per §5) so the Institution's successor service provider can continue verification using the institution's existing public-key registry and verifier;
>    (e) preserve cryptographic keys (public keys, public-key fingerprints, signed certificates) for the Institution's continuing verification of chain entries produced during the contract term;
>    (f) destroy or return the Institution's IKM material per the Institution's instructions, with documented evidence of destruction;
>    (g) destroy or return all other Institution data per the Institution's instructions; and
>    (h) provide the Institution and the Competent Authority with documented evidence of (f) and (g).
>
> 9.3 The Service Provider shall not, on termination, retain Institution data beyond what is required to comply with applicable law (including record-retention obligations imposed by financial-services regulation). Any data the Service Provider retains under such an obligation shall be (i) limited to the minimum necessary, (ii) securely stored, (iii) accessible to the Institution and the Competent Authority on request, and (iv) destroyed at the end of the regulatory retention period.
>
> 9.4 The Service Provider acknowledges that the Institution may require, under DORA Article 28(7), a "stress-test" of the exit strategy in advance of any contract amendment, change of Service Location, or material change to Sub-Outsourcing. The Service Provider shall cooperate with such stress-tests on the Institution's reasonable request.
>
> 9.5 The "Transition Period" is [INSTITUTION:fillin — duration; typically 12 months for Critical or Important Functions, 6 months for non-critical]. The Transition Period shall not be reduced without the Institution's written consent, regardless of the reason for termination.

### Explanatory note

DORA Article 28(7) requires the institution to ensure the contract provides for "the orderly termination of the contract in the event that the financial entity, the financial entity's Competent Authority or the [Lead Overseer] determines that the contract is unsuitable, illegal or otherwise unacceptable." Article 28(8) requires the institution to maintain "exit plans" that are "comprehensive, documented, and sufficiently tested."

The chain-of-custody-specific provisions in §9.2(d) and (e) ensure that on transition, the Institution can continue to verify chain entries produced during the contract term using its existing verifier and public-key registry. Without §9.2(d)(e), the Institution's chain-integrity record could become un-verifiable after Service Provider transition — a control-completeness failure the Competent Authority would surface in any subsequent supervisory review.

## §10 — Termination rights (DORA Article 28(7); EBA/GL/2019/02 §13.6)

### Normative text

> 10.1 The Institution may terminate this Agreement, on written notice and without penalty, in any of the following cases:
>    (a) the Service Provider commits a material breach of this Agreement that is not cured within thirty (30) days of notice;
>    (b) the Service Provider becomes insolvent, enters into administration, or any equivalent insolvency proceeding;
>    (c) the Service Provider fails to maintain the security and operational requirements of §3, or fails on a material chain-of-custody verification under FFIEC Specification §7 that is not remediated within the timeline the Institution requires;
>    (d) the Service Provider makes a material change to its Sub-Outsourcing arrangements that the Institution has objected to under §6.3;
>    (e) the Competent Authority of the Institution determines, on documented grounds, that the Services or the Service Provider are unsuitable, illegal, or otherwise unacceptable, and the Institution is required by the Competent Authority to terminate the Agreement;
>    (f) the Competent Authority's exercise of its rights under §4 is materially obstructed by the Service Provider, and the obstruction is not cured within fifteen (15) days of notice;
>    (g) the Service Provider fails to make the Initial Notification within the DORA Initial Notification Window in §8.1(a), and the failure is repeated;
>    (h) the Service Provider's vendor-conformance attestation under the FFIEC chain-of-custody attestation procedure is revoked by the project-side working group; or
>    (i) any other ground specified in EBA/GL/2019/02 §13.6 or DORA Article 28(7).
>
> 10.2 On termination under §10.1, the Service Provider shall comply with the exit-strategy obligations of §9.

### Explanatory note

The termination grounds in §10.1(c), (g), and (h) are the chain-of-custody-specific termination triggers. A persistent failure on chain-of-custody verification, a repeated failure to meet the DORA initial notification window for incidents, or a vendor-conformance attestation revocation each independently support termination — and the institution's contracts team uses these explicit grounds rather than relying solely on the general material-breach provision.

## §11 — Register of information (DORA Article 28(3))

### Normative text

> 11.1 The Service Provider shall provide to the Institution, in a form and on a cadence the Institution requires, the data necessary for the Institution to maintain its register of information under DORA Article 28(3) and the Implementing Technical Standards adopted thereunder.
>
> 11.2 The data the Service Provider provides shall include:
>    (a) the Service Provider's legal name, registration number, and jurisdiction of incorporation;
>    (b) the contractual identifiers (contract reference number, effective date, termination or expiration date);
>    (c) the Services provided and the Critical or Important Function (if any) they support;
>    (d) the Service Locations enumerated by Member State or third country;
>    (e) the Sub-Outsourcing chain (per §6.2 register);
>    (f) the security certifications maintained by the Service Provider (ISO/IEC 27001, SOC 2, etc.);
>    (g) the data-protection certifications and the location of personal-data processing;
>    (h) for chain-of-custody Services, the vendor-conformance attestation reference under the FFIEC chain-of-custody attestation procedure (per `docs/vendor-conformance-attestation.md`); and
>    (i) any other data the Institution reasonably requires for its register.
>
> 11.3 The Service Provider shall update the data it provides without undue delay on any material change.
>
> 11.4 The Service Provider shall provide the data in machine-readable form (such as JSON, XML, or CSV) compatible with the Institution's register-of-information system.

### Explanatory note

DORA Article 28(3) and the Implementing Technical Standards on the register of information (Commission Implementing Regulation (EU) 2024/2956) require the financial entity to maintain a register of all ICT third-party arrangements supporting Critical or Important Functions. The Service Provider's contractual obligation to provide the register data is the load-bearing operational support for the Institution's Article 28(3) compliance.

The chain-of-custody-specific reference in §11.2(h) imports the vendor-conformance attestation procedure as the trust path the Institution cites in its register entries for chain-of-custody Services.

## §12 — Threat-led penetration testing cooperation (DORA Articles 26-27)

### Normative text

> 12.1 The Service Provider shall cooperate fully with any threat-led penetration testing ("TLPT") the Institution conducts under DORA Articles 26-27 and the TIBER-EU framework, including:
>    (a) granting the TLPT lead and the TLPT testers reasonable access to the Service Provider's infrastructure for the duration of the TLPT engagement;
>    (b) providing the TLPT testers with documentation and configuration evidence necessary for the TLPT to be effective;
>    (c) refraining from disclosing the existence or scope of the TLPT to the Service Provider's personnel beyond the small group authorized to know (typically: senior security leadership, compliance, and legal — but NOT the Service Provider's day-to-day blue team unless the TLPT is a "purple-team" engagement);
>    (d) preserving any chain-of-custody records of TLPT activity (under the spec's `audit.tlpt.engagement_id` attribute) under the same retention and the same separation-of-duties controls as production data, and confirming such records are not exposed to the Service Provider's personnel beyond the authorized group; and
>    (e) post-TLPT, providing the Institution with the Service Provider's documented response to any findings within the timeline the Institution and the TLPT lead specify.

### Explanatory note

DORA Articles 26-27 mandate TLPT for significant institutions on a three-year cadence under the TIBER-EU framework. The Service Provider's cooperation is necessary because the TLPT engagement typically spans the institution's perimeter and the Service Provider's perimeter together; without contractual cooperation, the TLPT lead cannot effectively test the integrated attack surface.

The chain-of-custody-specific provision in §12.1(d) closes a gap the spec deliberately leaves to institution-side procedure: TLPT activity captured in the chain MUST not become an attacker's playbook in production. The retention-and-separation discipline ensures TLPT chain entries are tagged, separated, and not surfaced to the Service Provider's broader operations team.

## §13 — Liability and indemnity (general; institution-specific)

### Normative text

> 13.1 The Service Provider shall indemnify the Institution against losses, damages, and expenses arising from:
>    (a) the Service Provider's breach of this Agreement;
>    (b) the Service Provider's failure to comply with applicable law (including DORA, GDPR, EBA/GL/2019/02, the FFIEC chain-of-custody Specification's normative requirements where applicable);
>    (c) the Service Provider's negligent or wilful act or omission in the performance of the Services; and
>    (d) any chain-of-custody integrity finding traced to the Service Provider's act or omission, including any cost of incident response, customer notification, regulatory penalty, or litigation defence.
>
> 13.2 The liability cap under this Agreement shall not apply to liability arising from:
>    (a) the Service Provider's wilful misconduct or gross negligence;
>    (b) the Service Provider's breach of confidentiality obligations;
>    (c) the Service Provider's breach of data-protection obligations under GDPR; or
>    (d) any chain-of-custody integrity finding traced to the Service Provider's wilful misconduct or gross negligence.

### Explanatory note

EBA/GL/2019/02 §13.7 requires the contract to provide for adequate liability arrangements. The institution's contracts team negotiates the liability cap and the cap-exclusions per its standard policy.

The chain-of-custody-specific liability allocation in §13.1(d) and §13.2(d) is the institution's recourse for an integrity finding traced to a Service Provider error (e.g., a Sub-Outsourcing chain that introduced a non-conformant HSM, a configuration drift that bypassed the software-key adapter exclusion).

## §14 — General provisions

### Normative text

> 14.1 This Agreement shall be governed by [INSTITUTION:fillin — governing law].
>
> 14.2 Disputes shall be resolved in [INSTITUTION:fillin — forum].
>
> 14.3 Notices shall be served on:
>    (a) Institution: [INSTITUTION:fillin — institution notice address];
>    (b) Service Provider: [INSTITUTION:fillin — service provider notice address].
>
> 14.4 The Service Provider acknowledges that the Institution's regulatory perimeter (DORA, GDPR, EBA/GL/2019/02, NIS2 where applicable, national prudential legislation, FFIEC chain-of-custody Specification where applicable) is part of the basis on which the Institution entered into this Agreement. Material changes to the regulatory perimeter that affect the Service Provider's obligations under this Agreement shall be addressed by good-faith amendment of this Agreement, without unreasonable delay.

### Explanatory note

The regulatory-perimeter acknowledgment in §14.4 is the institution's hook for adjusting the contract when DORA RTSs are amended, when an Adequacy Decision changes, when a national prudential framework adopts the chain-of-custody specification as a normative reference, or when any other regulatory change affects the contract's operational substrate.

## §15 — Variant clauses (use as needed)

### §15.A — Critical ICT third-party service provider designation (DORA Article 31)

Use this variant when the Service Provider has been or is likely to be designated as a "critical ICT third-party service provider" under DORA Article 31 (subjecting the Service Provider to direct oversight by the Lead Overseer).

#### Normative text

> 15.A.1 If the Service Provider is designated as a critical ICT third-party service provider under DORA Article 31, the Service Provider shall:
>    (a) notify the Institution within five (5) business days of receiving notice of designation;
>    (b) cooperate with the Lead Overseer's exercise of its powers under DORA Articles 32-44, including providing the Lead Overseer with access to information, premises, systems, and personnel;
>    (c) make available to the Institution any reports, recommendations, or findings the Lead Overseer issues that affect the Services or the Institution; and
>    (d) implement any remedy the Lead Overseer requires, on the timeline the Lead Overseer requires.

#### Explanatory note

The Article 31 designation transforms the Service Provider's regulatory posture: the Lead Overseer exercises direct oversight, and the institution's contractual relationship is overlaid by the Lead Overseer's powers. The contractual provision ensures the Service Provider's cooperation with the Lead Overseer is not a contractual gap.

### §15.B — Concentration-risk monitoring (DORA Article 29)

Use this variant when the Service Provider serves a significant number of EU-supervised financial entities and DORA Article 29 concentration-risk monitoring is relevant to the institution's posture.

#### Normative text

> 15.B.1 The Service Provider shall provide the Institution, on the Institution's quarterly cadence, with anonymized statistics on the Service Provider's customer base relevant to DORA Article 29 concentration-risk monitoring, including:
>    (a) the number of EU-supervised financial entities the Service Provider serves;
>    (b) the breakdown of those entities by Member State;
>    (c) the proportion of the Service Provider's revenue derived from EU-supervised financial entities; and
>    (d) any other concentration-risk indicator the Institution requires.
>
> 15.B.2 The data provided under §15.B.1 shall be sufficient for the Institution to assess concentration risk under its DORA Article 29 obligations without identifying specific customers of the Service Provider, and the Institution shall hold the data in confidence.

#### Explanatory note

DORA Article 29 requires the financial entity to assess concentration risk arising from ICT third-party arrangements. The Service Provider's customer-base statistics are the load-bearing input the institution needs to make the assessment.

### §15.C — eIDAS qualified-electronic-signature posture (eIDAS Articles 25-29)

Use this variant when the chain seal signature is intended to qualify as a "qualified electronic signature" under eIDAS Regulation (EU) 910/2014.

#### Normative text

> 15.C.1 The Service Provider shall ensure the HSM used to sign chain-of-custody seal records is:
>    (a) a Qualified Signature Creation Device (QSCD) as defined in eIDAS Article 3(23) and listed on the EU Trust List for the Service Provider's Member State (or a Member State the Institution accepts); and
>    (b) operated under a Qualified Trust Service Provider (QTSP) under Article 24, with a current qualified-status notification.
>
> 15.C.2 The public key supporting the chain-of-custody seal-signature verification shall be issued under a qualified electronic-signature certificate satisfying eIDAS Article 28.
>
> 15.C.3 The Service Provider shall notify the Institution within five (5) business days of any change to the QSCD, the QTSP, or the qualified-status notification that affects the qualified-signature posture.

#### Explanatory note

eIDAS Article 25 grants a qualified electronic signature the equivalent legal effect of a handwritten signature throughout the Union. For institutions whose chain-of-custody outputs may enter Union legal proceedings, the qualified-signature posture is materially stronger than the advanced electronic signature default.

The institution's chain-of-custody seal-signature is, by spec §10.5, an Ed25519 signature produced by an HSM at FIPS 140-2 Level 3 (or higher). That construction satisfies the four Article 26 criteria for an advanced electronic signature without further qualification. To elevate to qualified electronic signature, the institution selects an HSM that is a QSCD on the EU Trust List and a public-key certificate issued by a QTSP — a posture worth contracting for in jurisdictions or proceedings where the qualified-signature default is the operational floor.

## Cross-reference table — clause to regulatory anchor

| Clause | Regulatory anchor |
|---|---|
| §1 — services description | DORA Article 30(2)(a); EBA/GL/2019/02 §13.1 |
| §2 — performance and SLA | DORA Article 30(2)(b)(e); EBA/GL/2019/02 §13.5 |
| §3 — information security | DORA Article 30(2)(c); EBA/GL/2019/02 §14; GDPR Article 32 |
| §4 — audit rights | EBA/GL/2019/02 §13.2; DORA Article 30(3) |
| §5 — pooled audits | EBA/GL/2019/02 §13.3 |
| §6 — sub-outsourcing | EBA/GL/2019/02 §§4.12, 4.13, 13.4; DORA Article 30(2)(a) |
| §7 — replication / EEA residency | DORA Article 12; spec §10.15; GDPR Articles 44-50 |
| §8 — incident notification | DORA Articles 18-23; EBA/GL/2019/02 §13.5 |
| §9 — exit strategy | DORA Article 28(7)(8); EBA/GL/2019/02 §15 |
| §10 — termination rights | DORA Article 28(7); EBA/GL/2019/02 §13.6 |
| §11 — register of information | DORA Article 28(3); Commission Implementing Regulation (EU) 2024/2956 |
| §12 — TLPT cooperation | DORA Articles 26-27; TIBER-EU |
| §13 — liability and indemnity | EBA/GL/2019/02 §13.7 |
| §14 — general provisions | (general) |
| §15.A — Article 31 designation | DORA Article 31, Articles 32-44 |
| §15.B — Article 29 concentration | DORA Article 29 |
| §15.C — eIDAS qualified signature | eIDAS Articles 25-29 |

## Cross-reference table — clause to FFIEC chain-of-custody spec section

| Clause | Spec reference |
|---|---|
| §3 — information security (chain-specific) | §10.3 (append-only), §10.5 (HSM custody), §10.6 / §10.6.1 (IKM length + CSPRNG provenance), §10.7 (software-key adapter), §10.8 (constant-time), §10.16 (SaaS-edge connectors), §10.17 (HSM partition ceremony attestation), §10.18 (CC8.1 cross-referencing), §10.19 (chain-coverage map), §10.20 (training-data retention floor), §10.21 (model-handover schema), §10.22 (pre-MAC redaction discipline) |
| §4 — audit rights (verifier testing) | §7 (verification procedure), §10.12 (exit-code contract) |
| §5 — vendor-conformance reliance | `docs/vendor-conformance-attestation.md` (project-side procedure) |
| §6 — sub-outsourcing | §10.5 (HSM operator role separation), §10.17 (per-customer-bank partition handle), §10.24 (entity succession) |
| §7 — replication | §10.15 (Pattern A and Pattern B), §10.16 (SaaS-edge mirror lag) |
| §8 — incident notification (chain failures) | §7 (named failure modes), §10.2 (operational events including `connector.outage` per §10.16) |
| §9 — exit (chain-data portability) | §5 (wire format), §6 (storage), §10.19 (coverage-map handover), §10.20 (training-data retention floor), §10.21 (model-handover schema), §10.24 (entity succession) |

## Summary

This template gives the institution a contractual clause set to attach to any vendor agreement covering chain-of-custody hosting, audit-ledger services, HSM custody, or related ICT third-party Services. The clauses cover: services description with locations; performance and service-level monitoring; information security with chain-specific requirements; full audit and inspection rights including direct Competent-Authority access; pooled-audit and third-party-reporting reliance; sub-outsourcing pre-approval and change notice; replication and EU-resident-data residency; incident notification on the DORA 4-hour clock; exit strategy with chain-data portability; termination grounds including chain-conformance failures; register-of-information data provision under DORA Article 28(3); TLPT cooperation; liability and indemnity; general provisions; and three variant clauses (DORA Article 31 designation, DORA Article 29 concentration-risk monitoring, eIDAS qualified-signature posture).

The template separates normative text (which appears in the executed contract) from explanatory notes (which support the institution's contracts team and the regulator's joint examination team). Marked fields are explicit so a reviewer can scan for unfilled gaps. The result is a contract clause set the Institution drops into a vendor agreement for an EBA/DORA-aligned outsourcing arrangement, satisfying EBA/GL/2019/02 §13 audit-rights requirements and DORA Article 30(2) contractual elements.

Margarethe G-1 / G-8 / G-10 closure: the template addresses DORA Articles 28-30 contractual provisions (clauses §1, §2, §3, §4, §6, §11), EBA/GL/2019/02 §13.2 audit-rights (clause §4), EU-held-keys discipline for vendor-hosted EU controllers (clause §7), and the qualified-signature posture (variant §15.C). The institution's contracts team paste-completes the marked fields, deletes the variants that don't apply, and the result is directly referenceable in a Union joint examination team's working file.
