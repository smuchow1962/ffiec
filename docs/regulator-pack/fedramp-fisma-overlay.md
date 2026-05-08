---
status: informative
alignment-reference: FISMA (44 USC §3551 et seq.); NIST SP 800-37 Rev 2 (RMF); NIST SP 800-53 Rev 5; NIST SP 800-53A Rev 5 (assessment procedures); NIST SP 800-171 Rev 3 (CUI in non-federal systems); NIST SP 800-172 (enhanced for high-value assets); NIST AI RMF 1.0 (AI 100-1); NIST AI 600-1 (Generative AI Profile); FedRAMP Rev 5 baselines; OMB M-22-09 (zero-trust); OMB M-23-22 (AI use); OMB M-24-04 (zero-trust progress); OMB M-24-10 (AI use-case inventories); OMB M-24-15 (AI governance); OMB Circular A-130; CMMC 2.0; CNSSP-15 / CNSA Suite 1.0 / 2.0; CISA BOD 22-01 (KEV); CISA BOD 23-01 (asset visibility); CISA Zero Trust Maturity Model 2.0; FedRAMP Authorization Act (FY2023 NDAA); FIPS 199 (security categorization); SP 800-218 / 218A SSDF; NARA records-retention schedules; M-23-07 (electronic records); CNSSI 1253 (Cross-Domain Solutions)
companion-docs:
  - docs/audit-procedures.md (institution-side P-1..P-N procedures consumed by 3PAO)
  - docs/regulator-pack/fdic-occ-examination-overlay.md (parallel banking-examiner overlay)
  - docs/templates/soc2-control-matrix.md (parallel commercial-attestation matrix)
  - docs/regulator-pack/deployment-package.md (deployment topology and authorization boundary)
  - docs/regulator-pack/handbook-mapping.md (FFIEC handbook mapping)
  - docs/regulator-pack/CSF-2.0.md (NIST CSF 2.0 mapping)
  - docs/incident-response-playbook.md (Scenarios 1–11)
  - docs/supply-chain.md (CISA BOD posture; SSDF alignment)
  - docs/edge-and-federated-ai.md (Cross-Domain Solutions composition)
  - spec/chain-of-custody-v1.md §1.4, §4.1, §4.3, §4.4, §7, §10.1, §10.5–10.13, §10.15
date: 2026-05-07
version: 1.0.0
---

# FedRAMP / FISMA Overlay — Federal Authorization Substrate

> **What this doc is.** Federal authorization-package overlay for the FFIEC chain-of-custody control. The overlay maps the chain to NIST SP 800-53 Rev 5 controls, OMB AI memoranda, NIST AI RMF, CUI handling under SP 800-171 Rev 3, CMMC 2.0 scoping for defense / IC use cases, CNSSP-15 / CNSA Suite for National Security Systems eligibility, zero-trust composition under M-22-09 / M-24-04, NARA records-retention for federal records, CISA BOD posture, FedRAMP 3PAO procedures, and reciprocity across agencies under the FedRAMP Authorization Act.
>
> **Audience.** Federal civilian agency CISOs, Authorizing Officials (AOs), Chief AI Officers (CAIOs), AI Governance Boards, FedRAMP 3PAOs, Federal CISO Council members, CISA agency-coordination staff, Defense / IC partner-sharing reviewers.
>
> **Scope discipline.** None of the content modifies the spec normative core. The overlay is an authorization-package layer that an adopting agency or 3PAO consumes once and reuses across agencies under FedRAMP reciprocity intent. The cryptographic substrate (FIPS-validated SHA-256, HMAC, Ed25519, HKDF) is sound; this overlay closes the authorization-package gaps an agency would otherwise write from scratch.

## How to use this overlay

1. **SP 800-53 Rev 5 control mapping (§A).** The agency's SSP cites this mapping rather than rewriting it per agency.
2. **OMB AI memo alignment (§B).** Agencies satisfying M-23-22, M-24-10, M-24-15 cite the chain primitives by what each memo calls them.
3. **NIST AI RMF subcategory mapping (§C).** AI Risk Management Plans cite the chain in AI RMF terms.
4. **CUI handling (§D).** Chain artifacts flowing to non-federal systems satisfy SP 800-171 Rev 3.
5. **CMMC 2.0 scoping (§E).** Defense / IC partner sharing scopes the chain components under CMMC Level 2 or Level 3.
6. **CNSSP-15 / CNSA Suite (§F).** v1.0a is civilian-CUI conformant; NSS deployment is a future variant track.
7. **Zero-trust composition (§G).** The chain's contribution to M-22-09 / M-24-04 zero-trust pillars and the CISA Zero Trust Maturity Model 2.0.
8. **NARA records-retention (§H).** Permanent-record handling under OMB Circular A-130 and M-23-07.
9. **CISA BOD posture (§I).** Vulnerability management and asset visibility for chain components.
10. **FedRAMP 3PAO procedures (§J).** Per-control assessment objectives in SP 800-53A Rev 5 format.
11. **FedRAMP reciprocity (§K).** Authorization-package shape that supports cross-agency reuse.
12. **Authorization boundary and FIPS 199 (§L).** Boundary scoping and security categorization.
13. **Cross-Domain Solutions (§M).** IC partner-sharing composition.

---

## §A. NIST SP 800-53 Rev 5 control-mapping table

The chain prescribes controls; each maps to NIST SP 800-53 Rev 5 controls for FedRAMP Moderate and High baselines. The mapping is the SSP-baseline. Each control has organization-defined values (ODVs) the agency parameterizes per risk tier.

### §A.1 Cryptographic and key-management controls

| Spec section | NIST SP 800-53 Rev 5 controls | ODVs (agency-parameterized) |
|---|---|---|
| §10.5 HSM custody | SC-12, SC-13, SC-28, SC-28(1), PE-3 (where on-prem HSM applies) | HSM FIPS validation level (SC-13) — ≥ FIPS 140-2 Level 3 |
| §10.6.1 RNG provenance | SC-12, SC-12(2), SC-12(3), SC-13 | RNG source identifier, FIPS validation certificate ID |
| §10.6 IKM length and entropy | SC-12, SC-13 | IKM minimum length (≥ 32 bytes) |
| §10.7 software-key adapter exclusion | CM-7, CM-7(1), CM-2, SI-7 | Production exclusion verification frequency |
| §4.1 inviolate property + §10.3 append-only | SI-7, SI-7(1), AU-9, AU-9(2), AU-9(3) | Storage-layer immutability mechanism (S3 Object Lock, Azure immutable blob, on-prem WORM) |
| §10.1 reconciliation | CA-7, CA-7(1), AU-12, AU-12(1) | Reconciliation cadence (default weekly per spec §10.1) |
| §7 verifier + §10.12 exit codes | AU-6, AU-6(1), AU-6(3) | Verifier-run cadence |
| §10.13 retention | AU-11, SI-12 | Retention period (≥ 7 years floor; see §H for permanent records) |
| §10.10 algorithm rotation | CM-3, CM-3(1), SC-12, SC-13 | Rotation cadence; rotation triggers |
| §10.4 time synchronization | AU-8, AU-8(1) | NTP server set; tolerance threshold |
| §4.4 OTLP + §5.1 TLS 1.3 | SC-8, SC-8(1), SC-23 | TLS minimum version (1.3); cipher suite set |
| §10.15 multi-region | CP-7, CP-9, SC-5 | Failover RTO/RPO |
| §10.8 constant-time comparison | SC-13 | (No ODV) |
| IR Playbook Scenarios 1–11 | IR-4, IR-5, IR-6, IR-8 | Notification timeline (36 hours per 12 CFR §53; 72 hours per IR Scenario 3) |
| Vendor management | SA-9, SA-9(1), SA-9(2) | Vendor SLAs; complementary user-entity controls |
| Source-code build trust | SA-11, SA-11(1), SR-3, SR-4 | Cosign signing; GPG verification; reproducible-build verification cadence |

### §A.2 Worked SSP fragment — SC-13 example

A Moderate-baseline agency cites the chain in its SC-13 control-implementation description:

> **SC-13 — Cryptographic Protection (Implementation).** The system uses FIPS-validated cryptographic primitives for all chain-of-custody integrity operations: SHA-256 (FIPS 180-4) for hashing, HMAC-SHA-256 (FIPS 198-1) for per-event MAC, Ed25519 (FIPS 186-5) for daily-seal signatures, HKDF (RFC 5869) for key derivation. All cryptographic operations execute in a FIPS 140-2 Level 3 HSM (`[AGENCY:fillin — HSM product and certificate ID]`), per FFIEC chain-of-custody specification §10.5. RNG provenance is recorded on the `master_key.generated` operational event per §10.6.1 and reflected in the agency's CC8.1 control description (per §A.1 RNG-source attestation).

### §A.3 Worked SSP fragment — SI-7 example

> **SI-7 — Software, Firmware, and Information Integrity (Implementation).** The chain-of-custody system enforces append-only integrity on the chain ledger per spec §4.1 and §10.3. Database role restrictions deny UPDATE, DELETE, and TRUNCATE on `events` and `daily_seals` tables. Storage-layer immutability is enforced via `[AGENCY:fillin — S3 Object Lock / Azure immutable blob / on-prem WORM]`. Tamper detection is provided by per-event HMAC binding and daily Merkle seals signed by the agency's HSM. Verifier runs (per spec §7) re-execute the integrity chain and produce PASS/FAIL with failure-reason granularity per spec §10.12 exit codes.

### §A.4 Continuous monitoring (CA-7) substrate

The chain supports continuous monitoring (CA-7) with deterministic, replayable signals:
- Verifier output is the operational-effectiveness signal.
- §10.1 reconciliation produces per-tenant-day evidence on a documented cadence.
- §10.2 operational events provide telemetry for the agency's CDM (Continuous Diagnostics and Mitigation) program.
- §10.12 exit codes integrate cleanly with ConMon tooling.

---

## §B. OMB AI memorandum alignment

OMB M-22-09, M-23-22, M-24-04, M-24-10, M-24-15 persist regardless of administration changes. The chain provides the technical substrate for several minimum risk-management practices.

### §B.1 OMB M-23-22 minimum risk-management practices

M-23-22 requires "automated logging of AI system inputs, outputs, and material changes" for rights-impacting and safety-impacting AI use cases. The chain delivers this:

| M-23-22 requirement | Chain primitive |
|---|---|
| Automated logging of AI inputs | Captured prompt content; spec §4.4 RECOMMENDED `gen_ai.request.*` attributes |
| Automated logging of AI outputs | Captured model response; `gen_ai.response.*` attributes |
| Tool calls | `chain_kind = "tool_call"` events |
| Routing decisions | `audit.routing.*` events per spec §4.4.1 |
| Material changes | `master.*`, `seal.*`, `master_key.generated` operational events; algorithm rotation events per spec §10.10 |
| Material changes — model versions | `gen_ai.response.model` per call; `audit.deployment.*` per spec §4.4.2 |

### §B.2 OMB M-24-10 — Annual AI Use Case Inventory

M-24-10 requires agencies to inventory AI use cases annually. Inventory fields populated from chain attributes:

| M-24-10 inventory field | Chain attribute |
|---|---|
| AI system identifier | `tenant_id` |
| Models in use | Distinct `gen_ai.response.model` values |
| Decision categories | `chain_kind` enum + `audit.*` event categories |
| Volume of decisions | Chain-entry counts per period |
| Rights / safety impact | Institution-defined customer-impact tier; documented in CC8.1 |
| Operational status | `master.reconciliation_completed` events; verifier output |

### §B.3 OMB M-24-15 — AI Governance Roles and AI Risk Management Plan

M-24-15 designates the Chief AI Officer (CAIO) and the AI Governance Board. The chain provides substrate for the AI Risk Management Plan:
- **Operational-effectiveness evidence:** verifier output across the period.
- **Monitoring evidence:** §10.1 reconciliation events.
- **Incident-response substrate:** IR playbook + chain-detected anomalies.
- **Accountability:** each AI decision has a chain entry; the entry's HMAC + Merkle + HSM signature ensures non-repudiation.

### §B.4 OMB M-22-09 / M-24-04 — Zero Trust Architecture

See §G for the full zero-trust composition.

### §B.5 OMB Circular A-130 — Information as a Strategic Resource

Chain artifacts are federal records subject to Circular A-130. See §H for records-management composition.

---

## §C. NIST AI RMF 1.0 subcategory mapping

NIST AI 100-1 (AI RMF 1.0, January 2023) and AI 600-1 (Generative AI Profile, July 2024) define functions and subcategories. The chain provides substrate for several.

### §C.1 AI RMF 1.0 mapping

| Function | Subcategory | Chain primitive |
|---|---|---|
| GOVERN | GOVERN 1.4 (legal and regulatory requirements satisfied) | Verifier output as integrity attestation; `audit.aml.remediation_milestone_tracked` events |
| GOVERN | GOVERN 4.1 (organizational practices) | CC8.1 control description; three-lines-of-defense alignment per `fdic-occ-examination-overlay.md` §3.9 |
| MAP | MAP 4.1 (risk identification) | `audit.kyc.customer_risk_score_computed` events; institution risk-tier assignment |
| MEASURE | MEASURE 2.7 (AI system security and resilience) | Verifier output; §10.1 reconciliation; incident-response activation events |
| MEASURE | MEASURE 4.2 (incidents reported) | IR Scenario 4 / Scenario 9 notification events |
| MANAGE | MANAGE 4.1 (post-deployment monitoring) | `audit.aml.performance_report_generated`; threshold-tuning events; lookback events |
| MANAGE | MANAGE 4.3 (incident remediation) | `audit.aml.lookback_remediation`; IR playbook completion events |

### §C.2 NIST AI 600-1 Generative AI Profile mapping

The chain captures GenAI-specific attributes:

| GenAI Profile element | Chain attribute |
|---|---|
| Model identification at inference | `gen_ai.request.model`, `gen_ai.response.model` |
| Prompt provenance | Captured prompt content; system prompt content-hash |
| Output integrity | Per-event MAC + daily Merkle seal |
| Hallucination tracking | `audit.aml.sar_narrative_drafted.narrative_integrity_check` (and analogous fields in other AML/KYC events) |
| Retrieval context | RAG document IDs and content-hashes per spec §4.4 |
| Decoding parameters | `gen_ai_parameters` per spec §4.4 (P-25 reproducibility surface) |

### §C.3 Citation in the agency AI Risk Management Plan

The CAIO cites the chain in the agency's AI Risk Management Plan (M-24-15):

> The agency operates a chain-of-custody control per FFIEC v1.0a specification for all rights-impacting and safety-impacting AI decisions. The chain captures inputs, outputs, tool calls, routing decisions, and material changes per OMB M-23-22 minimum practices. Independent verification (§7 verifier procedure) confirms operational effectiveness. The chain is mapped to NIST SP 800-53 Rev 5 controls (per FedRAMP-FISMA overlay §A) and to NIST AI RMF subcategories (per overlay §C).

---

## §D. CUI handling under NIST SP 800-171 Rev 3

Chain artifacts in federal-agency deployments frequently include Controlled Unclassified Information (CUI) — agency decisions about citizens are CUI Privacy under the CUI Registry (32 CFR Part 2002), and depending on the agency, also CUI Tax (IRS), CUI Banking (Treasury), CUI Health Information (HHS).

### §D.1 CUI categorization for chain artifacts

| Artifact | CUI category (typical) |
|---|---|
| Captured prompt and model response | CUI Privacy (always); CUI Tax / Banking / Health (agency-specific) |
| Operational events | CUI Privacy (when references citizen-identifying data) |
| Daily seal records | Public (the seal record itself; the underlying chain is CUI) |
| Master-key fingerprint roster | CUI (agency-internal cryptographic material) |
| HSM audit logs | CUI (agency-internal access logs) |
| Verifier output | CUI Privacy (when references citizen-identifying chain entries) |

### §D.2 SP 800-171 Rev 3 controls for non-federal systems

When chain artifacts flow to non-federal systems (cloud HSM, vendor-hosted ledger, supplier-side support), SP 800-171 Rev 3 (May 2024 revision) controls apply:

| SP 800-171 Rev 3 family | Chain control |
|---|---|
| 03.01 — Access Control | HSM access restriction; ledger writer role separation (P-1, P-2) |
| 03.03 — Audit and Accountability | Operational events; verifier output; retention per §10.13 |
| 03.04 — Configuration Management | CC8.1 posture documentation; algorithm rotation |
| 03.05 — Identification and Authentication | mTLS / SPIFFE / HSM-issued tokens per spec §4.1.1 |
| 03.06 — Incident Response | IR Playbook Scenarios 1–11 |
| 03.08 — Media Protection | Storage-layer immutability; backup integrity |
| 03.09 — Personnel Security | HSM operator screening; seal-job operator separation |
| 03.10 — Physical Protection | HSM physical security; data center controls |
| 03.11 — Risk Assessment | Threat model documentation |
| 03.12 — Security Assessment | 3PAO procedures (§J) |
| 03.13 — System and Communications Protection | TLS 1.3; cryptographic substrate |
| 03.14 — System and Information Integrity | Verifier output; tamper detection |

### §D.3 CUI marking guidance for chain artifacts

#### §D.3.1 OTLP transport markings

Header attributes naming the CUI category:
- `cui_category`: `"PRIV"` (Privacy), `"TAX"`, `"BANK"`, `"HLTH"`
- `cui_dissemination`: `"FEDONLY"` (federal-only), `"NOFORN"` (no foreign nationals), `"BASIC"` (basic CUI)

The transport layer carries the marking; the chain entry's content remains tokenized per privacy-by-design.

#### §D.3.2 Storage markings

Chain audit-file headers and seal records carry CUI markings per the CUI Registry banner conventions:

```
//CUI//PRIV/FEDONLY
```

The marking is institution-managed; the spec is marking-agnostic.

### §D.4 SP 800-172 enhanced controls for high-value assets

For rights-impacting AI under M-23-22 (typically high-value assets), SP 800-172 enhanced controls supplement SP 800-171. Relevant enhanced controls:

| SP 800-172 control | Chain composition |
|---|---|
| 03.04.5e — Configuration Change Management — Identifying and Authenticating | Cosign + GPG verification of verifier and SDK builds |
| 03.06.5e — Incident Handling — Specialized Capabilities | Chain-detected anomaly + IR playbook |
| 03.13.2e — Communications Authenticity | Per-event MAC + daily seal HSM signature |

---

## §E. CMMC 2.0 scoping for defense / Intelligence Community use

When chain artifacts are shared with Department of Defense or Intelligence Community partners, the Cybersecurity Maturity Model Certification (CMMC 2.0) framework applies to the vendor and the cloud-service stack.

### §E.1 CMMC 2.0 levels

| Level | Scope | Chain composition |
|---|---|---|
| Level 1 | Foundational; FCI (Federal Contract Information) only | Not applicable for chain (chain handles CUI) |
| Level 2 | Advanced; CUI; SP 800-171 Rev 3 controls + 3rd-party assessment | Default for civilian-to-defense partner sharing |
| Level 3 | Expert; high-value assets; SP 800-171 + SP 800-172 + DoD-led assessment | Required for APT-resistant deployments |

### §E.2 Chain components inside the CMMC boundary

For defense / IC use, the CMMC boundary includes:
- Vendor SDK distributed to the contractor.
- Cloud HSM (when used).
- Ledger storage (vendor-hosted or self-hosted).
- Verifier-build pipeline.
- Cosign / GPG signing infrastructure.

### §E.3 CMMC Level 2 control coverage

Level 2 maps directly to SP 800-171 Rev 3. The CUI-handling annex in §D carries forward. Specific Level 2 considerations:
- Vendor's SDK build pipeline meets SSDF practices (SP 800-218).
- Cloud HSM holds DoD Impact Level 4 (IL4) authorization or higher.
- Vendor's SOC 2 Type II report covers the chain control area.

### §E.4 CMMC Level 3 supplemental controls

Level 3 adds SP 800-172 enhanced controls. Specific Level 3 considerations:
- Vendor performs continuous monitoring per SP 800-172 §03.06.5e.
- HSM holds DoD Impact Level 5 (IL5) or higher.
- Cosign keys held in HSM with attestation chain to a DoD-approved root.

### §E.5 DoD CIO Cybersecurity Service Provider scoping

When cloud HSM or ledger operates as a managed service for defense / IC use, the DoD CIO's Cybersecurity Service Provider (CSSP) scoping applies. The agency's contracting officer engages the CSSP for scoping; the chain's substrate (operational events, verifier output, IR playbook) supports the CSSP's continuous-monitoring obligations.

### §E.6 Authorization-track separation

Defense-side authorization is led by DoD, not the civilian Authorizing Official. The CMMC scoping is substrate; the authorization decision is the DoD program office's. The agency expanding from civilian-only to civilian-plus-defense maintains a continuous evidence path via this overlay; the authorization tracks remain separate.

---

## §F. CNSSP-15 / CNSA Suite — National Security Systems eligibility

When the chain is used for a National Security System (NSS) — touching classified information or intelligence operations — Committee on National Security Systems Policy 15 (CNSSP-15) governs cryptographic algorithm selection.

### §F.1 v1.0a posture

The chain's v1.0a algorithm pair (Ed25519 + SHA-256) is FIPS-validated under FIPS 186-5 and FIPS 180-4 respectively. For FedRAMP Moderate/High and civilian-CUI workloads, v1.0a is conformant.

CNSA Suite 1.0 prefers ECDSA P-384 + SHA-384 + AES-256-GCM + RSA-3072+. Ed25519 is not currently in CNSA Suite 1.0. **v1.0a is therefore NOT eligible for NSS deployment without an algorithm extension.**

### §F.2 Path to NSS eligibility

The cleaner forward path is a future variant under CNSA Suite 1.0 (and a future variant under CNSA Suite 2.0 when post-quantum migration lands). The variant would be a `posture = nss-cnsa-1.0` option at SDK construct time, dispatched via the spec §4.1.2 vendor-namespaced constants mechanism, with §4.3 sign_payload extended to bind the algorithm selector.

### §F.3 NSS-conformant variant — informative sketch

| Component | v1.0a algorithm | NSS-conformant variant |
|---|---|---|
| Per-event MAC | HMAC-SHA-256 | HMAC-SHA-384 |
| Daily Merkle | RFC 6962 over SHA-256 | RFC 6962 over SHA-384 |
| Daily-seal signature | Ed25519 | ECDSA P-384 (Pattern A cosign with Ed25519 + ECDSA P-384) |
| Symmetric encryption (where used) | AES-256-GCM | AES-256-GCM (already conformant) |
| Key derivation | HKDF-SHA-256 | HKDF-SHA-384 |

The variant is a future-scope item, not part of v1.0a. v1.0a institutions reading this overlay see the limitation explicitly: civilian-CUI deployments are conformant; NSS deployments require the future variant.

### §F.4 CNSA Suite 2.0 — post-quantum migration

CNSA Suite 2.0 (2022 update) adds CRYSTALS-Kyber and CRYSTALS-Dilithium with a 2030–2033 PQ transition timeline. Future spec evolution will define a v1.x.y CNSA-Suite-2.0 variant with PQ algorithms when the algorithms are FIPS-validated and CNSA-approved for production NSS use.

### §F.5 Documenting the limitation

The agency's SSP names v1.0a as the deployed variant for civilian-CUI use. The agency explicitly notes the v1.0a algorithm set is not CNSA Suite 1.0 conformant; NSS deployment is deferred pending the future variant.

---

## §G. Zero Trust Architecture composition (M-22-09 / M-24-04)

OMB M-22-09 ("Moving the U.S. Government Toward Zero Trust Cybersecurity Principles," 26 January 2022) and M-24-04 (FY2024 progress update, 4 December 2023) require federal agencies to deploy zero-trust architectures.

### §G.1 The five M-22-09 pillars and chain composition

| Pillar | Chain composition |
|---|---|
| **Identity** | Spec §4.1.1: SDK authenticates to IKM custodian via SPIFFE / mTLS / HSM-issued tokens. Verifier access uses PIV/CAC for human operators. |
| **Devices** | The SDK host must be on the agency's compliant-endpoint roster. The verifier host similarly. |
| **Networks** | TLS 1.3 per spec §5.1; mTLS to ledger; OTLP authentication per spec §4.4.3. |
| **Applications and Workloads** | SDK / ledger / verifier as zero-trust workloads with attestation. The spec §1.4 compositional-security argument is precisely zero-trust defense-in-depth. |
| **Data** | Classification-aware per §D (CUI handling). Data-classification engine reads chain content sensitivity from CUI markings. |

### §G.2 CISA Zero Trust Maturity Model 2.0 cross-walk

| ZTMM stage | Chain contribution |
|---|---|
| Traditional | Not applicable (chain is zero-trust by design) |
| Initial | mTLS / SPIFFE in §4.1.1; basic operational-events monitoring |
| Advanced | Full zero-trust pillar coverage; CC8.1 documents the pillar mapping |
| Optimal | Continuous attestation of SDK and verifier hosts; HSM attestation chain |

### §G.3 CC8.1 language for zero-trust pillar mapping

> The chain-of-custody control supports the agency's zero-trust architecture per OMB M-22-09 / M-24-04. Identity (§4.1.1 mechanisms + PIV/CAC verifier access). Devices (compliant-endpoint posture for SDK and verifier hosts). Networks (TLS 1.3 per §5.1 + OTLP authentication per §4.4.3). Applications (SDK / ledger / verifier as zero-trust workloads with attestation). Data (CUI classification per overlay §D; data-classification engine reads CUI markings). The chain's contribution to the agency's M-24-04 zero-trust progress reporting is the named substrate, not ambient infrastructure.

---

## §H. NARA records-retention and OMB Circular A-130

Federal records are subject to NARA records-retention schedules. The General Records Schedule (GRS) covers most administrative records; agency-specific schedules cover mission records.

### §H.1 Records categorization for chain artifacts

| Artifact | Typical record category |
|---|---|
| Operational events (`master.*`, `seal.*`, `ledger.*`) | Temporary; GRS 4.1 (general administrative) |
| AI-decision events (citizen-facing) | Often permanent under agency-specific NARA schedule |
| Daily seal records | Temporary; retained for permanent records' integrity-anchor lifetime |
| HSM audit logs | Temporary; security-records retention |
| Verifier output | Temporary; audit-records retention |

Permanent records examples:
- SSA disability decisions.
- IRS audit decisions.
- VA benefits adjudication.
- CMS coverage determinations.
- Department of Education student-aid decisions.

### §H.2 Retention floor vs permanent-record obligations

The spec §10.13 7-year floor is the FFIEC banking horizon. Federal-agency permanent records must be preserved for decades, eventually transferred to NARA. During agency custody:
- The chain entries for permanent records are retained per the agency's NARA schedule.
- The 7-year floor is the minimum; the actual retention extends to the permanent-record horizon.

### §H.3 Format-preservation for permanent records

For multi-decade preservation:
- **Algorithm rotation (§4.3.2 / §10.10.2) becomes load-bearing** for permanent records. As SHA-256 weakens (anticipated PQ-era cryptographic erosion), the agency rotates to SHA-384 or post-quantum hash families per the variant track in §F.
- **Key rotation** preserves verifiability. The IKM-retention coupling (§10.9) extends from 7 years to the permanent-record horizon.
- **Verifier preservation.** The verifier binary used to validate permanent records must remain executable for the record's lifetime. The agency stores verifier binaries with reproducible-build provenance.

### §H.4 OMB Circular A-130 information-resources management

The chain is named in the agency's Information Resources Management posture under Circular A-130. The chain entries are managed information; the verifier output is decision-support information.

### §H.5 M-23-07 transition to electronic records

OMB M-23-07 requires federal records to be transitioned to electronic format. Chain entries are already electronic and NARA-ingestible. Specific compliance:
- Chain entries are retained in NARA-compliant electronic formats (JSON, CSV with documented schemas).
- Metadata is preserved alongside the records.
- The verifier and verifier-build provenance is preserved with the records.

### §H.6 CC8.1 language for permanent-record handling

> For chain entries categorized as permanent records under the agency's NARA schedule, the institution's IKM retention extends from the 7-year floor (§10.9) to the permanent-record horizon. Algorithm rotation becomes a permanent-records preservation control, not just a cryptanalysis-defense control. Verifier binaries used to validate permanent records are retained with reproducible-build provenance for the record's lifetime. The records-management team coordinates with the IT team on format-preservation milestones aligned with the agency's NARA schedule.

---

## §I. CISA BOD posture (vulnerability management and asset visibility)

CISA Binding Operational Directives apply to federal civilian executive-branch agencies.

### §I.1 BOD 22-01 — Known Exploited Vulnerabilities

BOD 22-01 ("Reducing the Significant Risk of Known Exploited Vulnerabilities," 3 November 2021) requires remediation of vulnerabilities in the CISA KEV catalog within specific timelines:
- 15 days for critical-severity KEV entries.
- 30 days for high-severity KEV entries.

### §I.1.1 Chain components subject to BOD 22-01

| Component | KEV exposure |
|---|---|
| Vendor SDK | Transitive dependencies in the SDK; vendor publishes CVE notifications |
| Verifier binary | Transitive dependencies in the verifier; vendor publishes CVE notifications |
| Ledger server | Database engine, application runtime, transitive dependencies |
| Cloud HSM endpoint | Cloud-provider responsibility; institution's CC8.1 names the responsibility split |

### §I.1.2 KEV remediation procedure

The agency's CC8.1 names:
- KEV-catalog monitoring procedure (typically daily).
- Vendor CVE-disclosure relationship (vendor's CC8.1-equivalent commitment).
- Patch-deployment cadence aligned with the 15-day / 30-day timelines.
- Compensating controls when patch deployment exceeds the timeline.

### §I.2 BOD 23-01 — Asset Visibility

BOD 23-01 ("Improving Asset Visibility and Vulnerability Detection on Federal Networks," 3 October 2022) requires:
- Asset inventories at 14-day discovery cadence.
- Vulnerability scanning at 7-day cadence.

### §I.2.1 Chain components in CDM asset inventory

| Component | CDM tracking |
|---|---|
| Vendor SDK on agency hosts | Endpoint asset; tracked by host-based agent |
| Ledger server | Server asset; tracked by network discovery + host-based agent |
| Cloud HSM endpoint | Cloud asset; tracked by cloud asset management |
| Verifier binary on agency hosts | Endpoint asset; tracked by host-based agent |

### §I.2.2 Chain-event-to-CDM-telemetry mapping

The agency's CDM program consumes chain operational events as telemetry:
- `master.*` events → key-management asset events.
- `seal.*` events → cryptographic-operation asset events.
- `ledger.*` events → storage asset events.
- `chain.*` events → application asset events.
- `hsm.*` events → cryptographic-hardware asset events.

### §I.3 Worst-case integration with IR Scenario 4

If a chain component itself enters the KEV catalog (e.g., a critical CVE in the SDK's transitive dependency that allows MAC bypass):
- The agency activates IR Scenario 4 (master-key-compromise-equivalent posture).
- 36-hour notification per 12 CFR §53 and federal incident-reporting requirements.
- The seal-job continuity discipline applies to ledger-component patching — the agency patches without disrupting the seal-publication SLA.

### §I.4 Vendor-responsibility composition

The agency's vendor-management procedure (per `docs/regulator-pack/fdic-occ-examination-overlay.md` §3.10 + this overlay's vendor mapping) requires:
- Vendor's CC8.1-equivalent commitment to upstream CVE disclosure timing.
- Vendor's commitment to KEV-catalog monitoring of their SDK and dependencies.
- Vendor's patch-deployment commitment aligned with BOD 22-01 timelines.

### §I.5 Composition with `docs/supply-chain.md`

The agency's supply-chain documentation (extending `docs/supply-chain.md`) carries:
- KEV-catalog monitoring procedure.
- BOD 22-01 alignment with vendor CVE disclosures.
- BOD 23-01 alignment with CDM asset inventory.
- SSDF (SP 800-218 / 218A) alignment for the SDK and verifier build pipelines.

---

## §J. FedRAMP 3PAO procedures (NIST SP 800-53A Rev 5)

FedRAMP authorizations are supported by Third-Party Assessment Organization (3PAO) assessments performed under SP 800-53A Rev 5. Each control has corresponding assessment objectives + methods (examine / interview / test).

### §J.1 Per-control 3PAO procedures

The 3PAO consumes the institution's evidence (per `audit-procedures.md`) and applies SP 800-53A Rev 5 assessment objectives.

#### §J.1.1 SC-13 (Cryptographic Protection) assessment

| Objective | Method | Evidence |
|---|---|---|
| FIPS-validated cryptographic primitives are used | Examine | HSM FIPS validation certificate; SDK release notes naming FIPS-validated libraries |
| Cryptographic operations execute in HSM | Examine | HSM audit logs; configuration records |
| Algorithms are current FIPS-approved | Examine | Spec §10.5–10.13 algorithm declarations |

#### §J.1.2 SI-7 (Software, Firmware, and Information Integrity) assessment

| Objective | Method | Evidence |
|---|---|---|
| Append-only enforcement at storage layer | Test | Attempt UPDATE/DELETE under ledger writer role; confirm rejection |
| Tamper detection at chain layer | Test | Run verifier; PASS confirms tamper-free |
| Failure-mode disposition | Examine | Spec §10.12 exit codes; institution's IR playbook |

#### §J.1.3 AU-9 (Protection of Audit Information) assessment

| Objective | Method | Evidence |
|---|---|---|
| Audit information protected from unauthorized modification | Test | Append-only enforcement test; HSM signature verification |
| Audit information protected from unauthorized deletion | Test | Storage-layer immutability test |
| Cryptographic protection of audit information | Examine | Per-event MAC + daily Merkle + HSM signature substrate |

#### §J.1.4 AU-12 (Audit Record Generation) assessment

| Objective | Method | Evidence |
|---|---|---|
| Audit records generated for each AI decision | Examine | Chain entries for the period |
| Audit record content includes required attributes | Examine | Chain attribute schema per spec §4.4 |
| Audit record generation is reliable | Test | Run verifier; confirm completeness |

#### §J.1.5 CA-7 (Continuous Monitoring) assessment

| Objective | Method | Evidence |
|---|---|---|
| Continuous monitoring strategy defined | Examine | Institution's ConMon plan citing the chain |
| Monitoring metrics defined | Examine | Verifier-output cadence; reconciliation cadence |
| Monitoring results reported | Examine | `audit.aml.performance_report_generated` (or analogous) events |

#### §J.1.6 IR-4 (Incident Handling) assessment

| Objective | Method | Evidence |
|---|---|---|
| IR plan addresses chain-related incidents | Examine | IR Playbook Scenarios 1–11 |
| IR plan tested | Test | Tabletop exercise records; chain-detected anomaly handling |
| Notification timelines met | Examine | 36-hour and 72-hour notification records |

### §J.2 SAR (Security Assessment Report) disposition

For each control, the 3PAO records:
- **Satisfied** — control implementation is consistent with the SSP and assessment objectives.
- **Other Than Satisfied (OTS)** — control implementation has deficiencies; deficiency severity per §J.3.
- **Not Applicable** — control does not apply to the system.

### §J.3 Deficiency severity

| Severity | Pattern |
|---|---|
| Low | Minor procedural gaps; no chain-integrity impact |
| Moderate | Documented control deficiency; remediation plan required |
| High | Significant control deficiency; remediation timeline shortened; possible POA&M with frequent reporting |
| Critical | Material weakness or chain-integrity defect; possible authorization withdrawal |

### §J.4 POA&M (Plan of Action and Milestones)

For each OTS finding, the institution submits a POA&M with:
- Description of the deficiency.
- Corrective action.
- Resource allocation.
- Completion timeline.
- Milestones.

### §J.5 Continuous monitoring after authorization

The 3PAO's annual reassessment consumes:
- The SSP (updated for control changes).
- The SAR (current findings).
- The POA&M (remediation tracking).
- Operational evidence (verifier output; reconciliation events; IR activations).

---

## §K. FedRAMP reciprocity (FedRAMP Authorization Act, FY2023 NDAA)

The FedRAMP Authorization Act (signed December 2022) codifies reciprocity intent: an authorization issued by one agency or by the FedRAMP Joint Authorization Board is intended to be reciprocally accepted by other agencies.

### §K.1 Reciprocity-friendly authorization-package shape

Cross-agency reciprocity depends on:
- **Consistent SSP control mapping** (§A) — the upstream agency assembles the mapping once.
- **Consistent 3PAO assessment procedures** (§J) — the upstream agency's 3PAO operates the same procedures the downstream agency expects.
- **Consistent CUI handling** (§D) — chain artifacts flowing across agencies are CUI-marked consistently.
- **Consistent OMB AI memo alignment** (§B) — M-23-22 / M-24-10 / M-24-15 evidence is cited consistently.
- **Consistent AI RMF mapping** (§C) — AI Risk Management Plans cite the chain in AI RMF terms consistently.

### §K.2 Authorization-package shape

The upstream agency assembles:
1. SSP with §A mapping.
2. 3PAO SAR with §J procedures.
3. CUI handling annex per §D.
4. OMB M-23-22 alignment per §B.
5. NIST AI RMF mapping per §C.
6. POA&M (if any).

The downstream agency inherits the package and customizes:
- Agency-specific ODVs (retention period, scan frequency).
- Agency-specific identity-management integration.
- Agency-specific records-management posture per §H.

### §K.3 CC8.1 language for authorization-boundary scoping

> The agency operates the chain-of-custody control under FedRAMP [Moderate / High] authorization. Components inside the authorization boundary: the SDK on agency hosts; the IKM registry; the verifier execution environment. Components outside the authorization boundary: the cloud HSM service (carved-out per FedRAMP subservice scoping); the vendor-distributed SDK build artifacts (vendor's FedRAMP authorization referenced). The agency's complementary user-entity controls (CUECs) are documented in §[N] of the SSP.

### §K.4 Cross-agency package transfer

When the downstream agency adopts the upstream agency's authorization:
1. The downstream agency reads the upstream SSP and SAR.
2. The downstream agency confirms the agency-specific ODVs match its risk tolerance (or amends with documented rationale).
3. The downstream agency's AO accepts the upstream authorization or issues a downstream authorization referencing the upstream package.
4. Continuous-monitoring obligations transfer to the downstream agency for its share of the system.

---

## §L. Authorization Boundary and FIPS 199 categorization

### §L.1 FIPS 199 categorization

Under FISMA, every federal information system has a security categorization derived from the worst-case impact across the Confidentiality / Integrity / Availability triad.

The chain inherits the categorization of the AI system it serves:

| AI use case | Typical FIPS 199 categorization |
|---|---|
| Internal IT operational AI (helpdesk, SOC analyst augmentation) | Low to Moderate |
| Citizen-facing eligibility decision (rights-impacting per M-23-22) | Moderate, often High |
| Safety-impacting AI (medical, public-safety) | High |
| Defense / IC use cases | Outside FIPS 199 scope; CNSSP and IC-CD apply |

### §L.2 Authorization boundary scoping

The chain's authorization boundary defines components inside (subject to the agency's full SP 800-53 Rev 5 control set) and outside (subject to vendor SLAs, third-party attestation, complementary user-entity controls).

#### §L.2.1 Components inside the boundary (typical)

- SDK on agency hosts.
- IKM registry.
- Verifier execution environment (when run on agency hosts).
- Operational-event logging on agency hosts.

#### §L.2.2 Components outside the boundary (typical)

- Cloud HSM service (FedRAMP-authorized; carved-out per §D.1).
- Vendor-hosted ledger (when applicable; FedRAMP-authorized; carved-out).
- Vendor-distributed SDK build artifacts (vendor's SOC 2 / ISAE 3402 / FedRAMP attestation; carved-out).

### §L.3 Complementary User-Entity Controls (CUECs)

For components outside the boundary, the agency operates CUECs:
- HSM access reviews (monthly).
- Key-rotation discipline per spec §10.10.
- Vendor SOC 2 / ISAE 3402 / FedRAMP attestation review on receipt.
- Supply-chain trust path validation per `docs/supply-chain.md`.

### §L.4 CC8.1 language for authorization boundary

> The agency's authorization boundary for the chain-of-custody control includes the SDK on agency hosts, the IKM registry, and the verifier execution environment. The cloud HSM service is carved out under FedRAMP subservice-organization scoping; the agency operates complementary user-entity controls (CUECs) per §L.3 of the FedRAMP-FISMA overlay.

---

## §M. Cross-Domain Solutions composition (CNSSI 1253; IC partner sharing)

When chain artifacts cross security domains for IC partner sharing — for example, agency-internal U//CUI artifact crossing into a TS//SI//NF compartment — they pass through Cross-Domain Solutions (CDS) under CNSSI 1253 and the IC's Cross-Domain Architecture.

### §M.1 CDS interaction with canonical bytes

The chain's RFC 8785 JCS canonical bytes interact with CDS dirty-word filtering:
- CDS dirty-word filters may flag canonical-JSON content (citizen PII, internal identifiers, model response text).
- The agency redacts at a layer above the chain — the original chain holds the unredacted form; the CDS-crossing form is a documented derived product.

### §M.2 CDS interaction with binary signatures

The chain's binary Ed25519 signatures interact with CDS format validators:
- CDS format validators may not recognize the Ed25519 binary signature shape.
- The agency Base64-encodes signatures for CDS-compatible transport.
- The receiving side decodes and verifies.

### §M.3 Integrity story across the CDS boundary

The original signature does NOT verify on the redacted form. The agency operates one of two postures:
- **Re-signed posture:** the redacted form is re-signed under the receiving agency's authority. The receiving side verifies under the receiving agency's IKM.
- **Original-chain posture:** the receiving side accesses the original chain (which requires receiving-agency access to the original IKM, typically prevented by CDS architecture).

In practice, the re-signed posture is operationally simpler. The CDS-crossing form is a separate signed artifact; the original chain remains in the sending agency's custody.

### §M.4 CC8.1 language for CDS posture

> The chain composes alongside the agency's Cross-Domain Solutions per CNSSI 1253. Chain artifacts crossing security domains are redacted at a layer above the chain; the CDS-crossing form is a documented derived product. Binary Ed25519 signatures are Base64-encoded for CDS-compatible transport. The receiving agency verifies the redacted form under its own IKM (re-signed posture); the original chain remains in the sending agency's custody.

---

## §N. Composition summary

| Authorization concern | Overlay section | Companion doc |
|---|---|---|
| SP 800-53 Rev 5 mapping | §A | `audit-procedures.md` |
| OMB M-23-22 / M-24-10 / M-24-15 alignment | §B | Spec §4.4 audit.* namespace |
| NIST AI RMF mapping | §C | Spec §4.4; AI Risk Management Plan template |
| CUI handling | §D | `docs/regulator-pack/privacy-by-design.md` |
| CMMC 2.0 scoping | §E | `docs/supply-chain.md` |
| CNSA Suite NSS eligibility | §F | Spec §10.10 algorithm rotation |
| Zero-trust composition | §G | Spec §4.1.1; CISA ZTMM 2.0 |
| NARA records-retention | §H | Spec §10.13 retention |
| CISA BOD posture | §I | `docs/supply-chain.md` |
| FedRAMP 3PAO procedures | §J | `audit-procedures.md`; `soc2-control-matrix.md` |
| FedRAMP reciprocity | §K | `docs/regulator-pack/deployment-package.md` |
| Authorization boundary | §L | `docs/regulator-pack/deployment-package.md` |
| Cross-Domain Solutions | §M | `docs/edge-and-federated-ai.md` |

---

## §O. Future scope (informative; not v1.0a normative)

Items noted for future evolution:

### §O.1 NSS-conformant variant

A future spec version will define a `posture = nss-cnsa-1.0` variant (ECDSA P-384 + SHA-384 + AES-256-GCM) per §F. Defense / IC partner sharing requires this variant.

### §O.2 CNSA Suite 2.0 PQ variant

When CNSA Suite 2.0 PQ algorithms are FIPS-validated and CNSA-approved for production NSS use, a future spec version will define a `posture = nss-cnsa-2.0` variant. The transition timeline aligns with the 2030–2033 federal PQ migration.

### §O.3 OMB AI memo updates

OMB AI memos evolve with administration changes and AI-policy updates. The chain's framing in §B (agency-policy-neutral substrate satisfying minimum-practice obligations) is designed to survive memo updates. New memos may require new mapping sections; the chain's primitives are unchanged.

### §O.4 Federal records-management evolution

NARA's electronic-records-management framework continues to evolve. The chain's format-preservation discipline (§H.3) is designed to survive NARA-format updates with minimal change; only the storage-format wrapper evolves.
