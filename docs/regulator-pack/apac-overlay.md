---
title: APAC Articulation Overlay — MAS / HKMA / APRA / RBI / FSC and the ASEAN Banking Integration Framework
status: informative
aligned-with:
  - MAS Technology Risk Management Guidelines (revised 2021)
  - MAS FEAT Principles (2018; 2024 revision under public consultation)
  - MAS Veritas assessment toolkit (phases 1, 2, 3)
  - MAS Notice 644 (Banking Act §60) cyber incident notification
  - HKMA TM-G-1 General Principles for Technology Risk Management
  - HKMA AI principles circular (2024, retail banking)
  - HKMA Cyber-Resilience Assessment Framework (CRAF) 2.0
  - APRA CPS 234 Information Security
  - APRA CPS 230 Operational Risk Management
  - RBI Master Direction on Outsourcing of Financial Services (updated 2024)
  - RBI Cyber Security Framework (2016, with periodic updates)
  - RBI letter on AI/ML governance expectations (November 2023)
  - Indian Digital Personal Data Protection Act 2023 (DPDP Act)
  - Singapore Personal Data Protection Act (PDPA)
  - Hong Kong Personal Data (Privacy) Ordinance (PDPO)
  - Australian Privacy Act 1988 + Privacy and Other Legislation Amendment Bill 2024
  - Korean Personal Information Protection Act (PIPA, March 2023 amendment)
  - APEC Cross-Border Privacy Rules (CBPR)
  - ASEAN Banking Integration Framework (ABIF)
  - SEACEN coordination forum
date: 2026-05-07
version: 1.0.0
---

# APAC Articulation Overlay

> **What this doc is.** A single articulation overlay that maps the FFIEC chain-of-custody v1.0a specification onto the APAC ex-Japan supervisory framework. Written so a competent authority — a MAS Technology Risk supervisor, an HKMA Banking Supervision examiner, an APRA Cyber and Operational Risk reviewer, an RBI Cyber Security and IT Examination team member, or a FSC Korea financial-IT supervisor — can read this document alongside the spec and confirm what the chain delivers, what it does not, and where supplementary measures are required for application in the relevant jurisdiction.

> **What this doc is NOT.** Not a normative extension to the v1.0a specification. Not a translation of the spec into APAC-framework language. Not an APAC-specific spec fork. The integrity primitives in v1.0a are framework-neutral cryptographic constructs; the supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) governs: the chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content.

> **Scope: APAC ex-Japan.** Japan-specific articulation (FSA, BoJ, PIPC, J-SOX, APPI, CRYPTREC) is covered separately in the Japan regulator-pack documents. This overlay covers Singapore, Hong Kong, Australia, India, Korea, Malaysia, Indonesia, Thailand, the Philippines, and the ASEAN Banking Integration Framework. For Korean-supervised institutions deploying under FSC and the Korean AI Basic Act 2025, the `korea-overlay.md` provides the deeper Korea-specific articulation; this overlay names Korea at the regional cross-jurisdiction level.

---

## 1. Scope and reading order

This overlay is consumed in five reading orders.

| Reader | Reading order |
|---|---|
| MAS Technology Risk supervisor (Singapore) | §2 (MAS FEAT mapping) → §3 (MAS Veritas composition) → §10 (incident timing) → §11 (cross-jurisdiction Pattern B) → §13 (translation table) |
| HKMA Banking Supervision examiner (Hong Kong) | §4 (HKMA TM-G-1 alignment) → §10 (incident timing) → §11 (cross-jurisdiction Pattern B) → §13 |
| APRA Cyber and Operational Risk reviewer (Australia) | §5 (APRA CPS 234 / CPS 230 alignment) → §10 → §11 → §13 |
| RBI Cyber Security and IT Examination team (India) | §6 (RBI guidelines + DPDP Act 2023) → §10 → §11 → §13 |
| FSC Korea financial-IT supervisor | §7 (FSC Korea cross-reference) → `korea-overlay.md` (full Korea articulation) |
| ASEAN multi-jurisdiction Tier-1 institution | §11 (Pattern B per-jurisdiction) → §12 (APEC CBPR) → §10 (incident timing matrix) |

The overlay does not duplicate material that already lives in adjacent regulator-pack documents. Cross-references are exact — if a section sends the reader to `international-transfers.md` §2.2, that section is the load-bearing source and this overlay names where it sits in the APAC-framework picture.

---

## 2. MAS FEAT Principles — chain-attribute composition

### 2.1 What FEAT requires

The Monetary Authority of Singapore published the Principles to Promote Fairness, Ethics, Accountability and Transparency (FEAT) for the use of artificial intelligence and data analytics in 2018; a 2024 revision is under public consultation. The four principles govern AI/data-analytics use in Singapore financial services.

| FEAT principle | What MAS expects |
|---|---|
| Fairness | The use of AIDA in decision-making is justifiable, with regular review for unintended bias and discriminatory effect |
| Ethics | The use of AIDA is aligned with the institution's ethical standards, values, and codes of conduct |
| Accountability | Internal accountability for the institution's use of AIDA; external accountability for the impact of AIDA on customers and counterparties |
| Transparency | Proactive communication about AIDA use; opportunity for customers to seek review of decisions |

### 2.2 FEAT-to-chain mapping

The chain provides direct primitives for Accountability and Transparency. Fairness and Ethics are application-domain decisions; the chain produces the substrate institutions use to evidence their fairness and ethics assessments, but does not itself decide whether the AIDA use is fair or ethical.

| FEAT principle | Chain attribute or event satisfying the principle's evidentiary requirement |
|---|---|
| Accountability — AIDA decisions are integrity-bound to model + parameters + inputs | The three-layer integrity construction (per-event MAC §4.1, daily Merkle seal §4.2, HSM Ed25519 signature §4.3); `gen_ai.response.model` (spec §4.4 normative) records the model identifier the vendor's API answered with; `audit.routing.decision_rationale` (spec §4.4.1) records the institution's internal decision rationale |
| Accountability — AIDA decisions are explainable, transparent, and fair | `audit.routing.decision_rationale` plus the institution's `audit.deployment.intent` schema (spec §4.4.2); the chain captures whether the decision was steady-state production, A/B test, canary, or vendor-side reroute |
| Transparency — Customers, when affected by AIDA decisions, are provided clear explanations | The institution's adverse-action notice schema (spec §10.11 for ECOA contexts; `audit.adverse_action.*` schema per `eu-articulation-extension.md` §4.2 for non-ECOA) — for Singapore institutions, `audit.adverse_action.regulatory_basis = sg-feat` is the recommended value |
| Fairness — substantive fairness assessment | The chain's integrity-bound decision data is the substrate; the institution's fairness-testing program (separate institutional artifact, not on the chain) consumes chain-derived data and produces the fairness assessment |
| Ethics — alignment with institution's ethical standards | Institution-side; the chain does not itself decide ethics. The institution's AIDA governance framework names the ethical standards; the chain's evidence supports compliance review |

### 2.3 The 2024 FEAT revision

The 2024 FEAT revision (under public consultation as of this overlay's publication date) is expected to expand the operational specificity of the principles, particularly around fairness-cohort analysis and post-deployment monitoring. The chain's design is forward-compatible with the revision: the underlying integrity primitives do not change; what changes is the institution's CC8.1 mapping. The institution updates its CC8.1 once the revised principles are finalised.

---

## 3. MAS Veritas assessment-framework composition

### 3.1 What Veritas is

MAS Veritas is the responsible-AI assessment toolkit consisting of Fairness, Ethics, Accountability, and Transparency assessment methodologies. Veritas phase 1 (2020) launched the framework; phase 2 (2022) published the Fairness assessment methodology; phase 3 (2024 pilot) is publishing the Ethics, Accountability, and Transparency assessment methodologies. Singapore institutions deploying high-stakes AI undergo Veritas assessment — formally voluntary, effectively required for MAS approval of new AIDA deployments.

### 3.2 Chain artifacts feeding Veritas assessment

| Veritas methodology phase | Chain artifact feeding the assessment |
|---|---|
| Veritas Fairness assessment (phase 2) | The chain's per-decision substrate: `gen_ai.request.messages`, `gen_ai.response.text`, model parameters, decision factors. The institution's fairness-cohort analysis runs against this substrate, computing demographic disparity per Veritas's Fairness methodology |
| Veritas Accountability assessment (phase 3) | The chain's `audit.routing.*` decision-rationale attributes; the institution's three-line-of-defense governance framework names the chain as the operational evidence of accountability |
| Veritas Transparency assessment (phase 3) | The chain's `audit.deployment.intent` (spec §4.4.2) plus the institution's adverse-action schema (`audit.adverse_action.*` with `regulatory_basis = sg-feat`); the institution's transparency posture is reviewed against the chain's evidence of customer-facing disclosure |
| Veritas Ethics assessment (phase 3) | The institution's AIDA governance framework; the chain provides the audit substrate for ethics-review proceedings if the institution's framework ties ethical review to specific decision events |

### 3.3 Veritas evolution and chain forward-compatibility

The Veritas methodology evolves — phase 4 is anticipated for 2027 and may extend the assessment to additional AI shapes (foundation models, multi-modal AI, multi-agent systems). The chain's design is forward-compatible: the integrity primitives do not change; the institution's CC8.1 names the Veritas phase and the corresponding chain-attribute mapping. As Veritas evolves, the institution updates its CC8.1; the chain artefacts captured under earlier Veritas phases remain valid evidence for those phases (per spec §10.9 retention).

---

## 4. HKMA TM-G-1 and 2024 AI principles alignment

### 4.1 TM-G-1 General Principles for Technology Risk Management

The Hong Kong Monetary Authority's TM-G-1 module is the foundational technology-risk supervisory module for Authorized Institutions (AIs) in Hong Kong. The 2024 HKMA circular on AI principles in retail banking adds AI-specific governance expectations.

| TM-G-1 module section | Chain prescription satisfying the requirement |
|---|---|
| TM-G-1 §3 (Governance and oversight) | The chain is one technical control within the institution's AI governance framework; the institution's CC8.1 names the chain's role |
| TM-G-1 §4.2 (Cryptographic key management) | Spec §10.5 HSM custody under FIPS 140-2 Level 3; the institution's CC8.1 names the HSM model and certification |
| TM-G-1 §5.1 (Audit log integrity) | The three-layer integrity construction (per-event MAC, daily Merkle seal, HSM signature); the verifier (spec §7) provides byte-reproducible verification |
| TM-G-1 §5.3 (Audit log retention) | Spec §10.9 retention coupling between the IKM and the chain entries; the institution's CC8.1 names the retention period |
| TM-G-1 §6 (Reconciliation procedures) | Spec §10.1 reconciliation procedures; audit procedures P-1 through P-10 |
| TM-G-1 §7 (Incident response) | The institution's IR runbook with HKMA notification timelines (per §10 of this overlay) |
| 2024 AI principles §3 (Transparency) | The institution's adverse-action schema (`audit.adverse_action.regulatory_basis = hk-hkma-ai-principles`); chain-captured decision rationale |
| 2024 AI principles §5 (Human oversight) | The chain's `audit.routing.reviewer_identity` records the natural person who reviewed the decision; for systems without human-in-the-loop, the field is absent and the institution's policy documents the absence |

### 4.2 HKMA Cyber-Resilience Assessment Framework (CRAF) 2.0

CRAF 2.0 introduces tiered cyber-resilience expectations for Authorized Institutions. The chain's integrity primitives compose with CRAF's tier-2 and tier-3 expectations on cryptographic-key management, audit-log integrity, and incident detection. The institution's CRAF self-assessment names the chain as the operational evidence for the relevant tier-2/tier-3 controls.

---

## 5. APRA CPS 234 and CPS 230 alignment

### 5.1 CPS 234 (Information Security)

Australian Prudential Regulation Authority Prudential Standard CPS 234 (effective 1 July 2019) requires APRA-regulated entities to maintain information security capability commensurate with the size and extent of threats.

| CPS 234 paragraph | Chain prescription satisfying the requirement |
|---|---|
| CPS 234 §15(a) (Confidentiality of information assets) | Institution's encryption-at-rest posture documented in `article-32-security-mapping.md`; per-tenant HKDF binding (spec §4.1) prevents cross-tenant lift |
| CPS 234 §15(b) (Integrity of information assets) | The three-layer integrity construction (per-event MAC, daily Merkle seal, HSM signature) |
| CPS 234 §15(c) (Availability of information assets) | Multi-region resilience patterns (spec §10.15 Pattern A active-active; Pattern B per-region) |
| CPS 234 §16 (Information security incidents) | The institution's IR runbook (per §10 of this overlay); incident-detection via spec §10.2 operational events |
| CPS 234 §17 (Notification of incidents) | 72-hour notification clock (per §10 of this overlay); the chain's verifier output (spec §7) is evidence supporting the notification |
| CPS 234 §18 (Internal audit review) | Audit procedures P-1 through P-10 (institutional adoption) |

### 5.2 CPS 230 (Operational Risk Management)

CPS 230 (effective 1 July 2025, replacing CPS 231/232/233) requires comprehensive operational risk management including third-party arrangements and business continuity.

| CPS 230 paragraph | Chain prescription satisfying the requirement |
|---|---|
| CPS 230 §3 (Operational risk identification) | The chain detects integrity failures via spec §10.2 operational events; the institution's risk register names AI-decision-integrity risk and the chain's mitigation |
| CPS 230 §5 (Third-party arrangements) | DORA-overlay §2 lattice (the institution + SDK provider + receiver/ledger provider + HSM provider + LLM provider); the institution's third-party risk register names each provider |
| CPS 230 §7 (Business continuity) | Multi-region resilience patterns (spec §10.15); the institution's BCP names the chain's recovery-time and recovery-point objectives |
| CPS 230 §9 (Operational resilience testing) | Annual operational-resilience drills; the institution's CC8.1 names the chosen seal cadence (daily default; hourly for tier-1 with regulator approval) and the recovery-time target |

### 5.3 Cross-reference to APRA CPS 232 (Business Continuity Management)

CPS 232 was replaced by CPS 230 effective 1 July 2025; institutions transitioning from CPS 232 to CPS 230 update their CC8.1 to name CPS 230 §7 (business continuity) rather than the legacy CPS 232 reference.

---

## 6. RBI guidelines + Indian DPDP Act 2023

### 6.1 RBI Master Direction on Outsourcing of Financial Services (2024 update)

The RBI Master Direction on Outsourcing of Financial Services (updated 2024) governs outsourcing arrangements between Indian banks and third-party providers. The chain's typical deployment involves an SDK from vendor A, ledger storage from vendor B, KMS / HSM from cloud provider C; the Master Direction §6 requires the bank's contract with vendor A to flow audit rights through to vendors B and C.

| Master Direction paragraph | Chain prescription satisfying the requirement |
|---|---|
| §3 (Material outsourcing approval) | The institution's procurement names the chain provider as material outsourcing if the chain supports a critical function; Board approval per §3.1 |
| §6 (Outsourcing contract requirements) | Audit-rights clause; supervisor-access clause; sub-outsourcing notification; exit strategy; business continuity. The DORA overlay §3 Article 30(2) checklist is structurally parallel and serves as the institutional template |
| §7 (Supervisor access) | The verifier (spec §7) is the offline, no-network-call inspection mechanism; the institution's CC8.1 names how RBI examiners access verifier output |
| §12 (Material risk events) | Incident notification per RBI Cyber Security Framework (per §10 of this overlay) |

### 6.2 RBI Cyber Security Framework (2016, updated)

The RBI Cyber Security Framework imposes specific encryption-key-management expectations on banks. The chain's spec §10.5 HSM custody satisfies the framework's tier-1 cryptographic-key-management expectations. The institution's CC8.1 names the HSM model, certification, and the rotation cadence.

### 6.3 Indian DPDP Act 2023 — cross-border transfer restrictions

The Digital Personal Data Protection Act 2023 introduces cross-border-transfer restrictions on personal data: the central government will notify "trusted countries" to which transfer is permitted; transfer to non-notified countries requires standard contractual clauses or a transfer-impact assessment.

The chain's typical multi-region deployment under spec §10.15 Pattern B is the recommended posture for Indian institutions: the `bank-in` tenant pins to AWS / Azure / Google APAC India region (Mumbai, Delhi, Hyderabad), with chain artifacts not crossing the Indian border. Institutions requiring cross-border processing (e.g., for cross-border-correspondent-banking AI) operate under DPDP's standard contractual clauses with the destination country.

| Deployment shape | DPDP Act 2023 conformance |
|---|---|
| Single-jurisdiction Indian-region tenant (`bank-in` pinned to APAC India region) | Conformant by default; no cross-border transfer occurs |
| Multi-jurisdiction with India operations (`bank-in` + non-India tenants under Pattern B) | Conformant if `bank-in` does not replicate to non-India regions; the institution's CC8.1 names the per-tenant residency posture |
| Cross-border processing (e.g., correspondent-banking AI flows India data to Singapore for processing) | Requires SCC + transfer-impact assessment; the institution's RoPA names the transfer mechanism |

### 6.4 RBI letter on AI/ML governance (November 2023)

The November 2023 RBI letter to banks on AI/ML governance expectations articulates governance, transparency, and explainability requirements. The chain provides the operational evidence; the institution's AI governance framework names how chain evidence supports each expectation.

---

## 7. FSC Korea — cross-reference to korea-overlay.md

The Financial Services Commission (Korea) and its examination arm, the Financial Supervisory Service (금융감독원), operate under the Personal Information Protection Act, the Credit Information Use and Protection Act, the Electronic Financial Transactions Act, the Electronic Financial Supervisory Regulation, and the AI Basic Act 2025.

For Korean-supervised institutions, the deeper Korea-specific articulation lives in `korea-overlay.md`. This APAC overlay covers Korea at the regional cross-jurisdiction level — the per-jurisdiction tenant naming under §11.2 includes the `bank-kr` mapping, and the incident-notification timing matrix in §10 includes the FSC Korea 24-hour clock.

Institutions operating across Korea and other APAC jurisdictions read both documents: this overlay for the regional cross-jurisdiction posture, and `korea-overlay.md` for the FSS examination procedures, the AI Basic Act 2025 mapping, the K-ISMS-P control alignment, the 망분리 architecture, and the KISA incident-notification clock.

---

## 8. Cross-reference to other APAC supervisors not enumerated above

For supervisors not enumerated in §2-§7 (Bank Negara Malaysia, Bank of Thailand, Bangko Sentral ng Pilipinas, Bank Indonesia, Reserve Bank of New Zealand, Bank of the Lao PDR), the institution's CC8.1 names the relevant supervisory framework and the chain-prescription mapping using the structural template established in §2-§6. The chain's primitives are framework-neutral; what the institution writes in CC8.1 differs per supervisor, not the chain construction.

The ASEAN Banking Integration Framework (ABIF) and SEACEN coordination forum operate at the regional supervisory layer; cross-border ASEAN operations under ABIF compose with Pattern B per-jurisdiction tenant naming (§11) and the APEC CBPR cross-jurisdiction posture (§12).

---

## 9. APAC cloud HSM region availability (informative)

Spec §10.5 names cloud HSM products without enumerating APAC region availability. The following is the APAC region availability for the named products as of this overlay's publication date.

| Cloud HSM product | APAC region availability |
|---|---|
| AWS CloudHSM | Tokyo `ap-northeast-1`, Seoul `ap-northeast-2`, Osaka `ap-northeast-3`, Singapore `ap-southeast-1`, Sydney `ap-southeast-2`, Jakarta `ap-southeast-3`, Mumbai `ap-south-1`, Hong Kong `ap-east-1` |
| Azure Managed HSM | Japan East/West, Korea Central, Southeast Asia, East Asia, Australia East/Central/Southeast, India Central/West/South |
| Google Cloud HSM | `asia-east1`, `asia-east2`, `asia-southeast1`, `asia-southeast2`, `asia-northeast1`, `asia-northeast2`, `asia-northeast3`, `asia-south1`, `asia-south2`, `australia-southeast1`, `australia-southeast2` |
| Thales Luna | On-prem (any APAC institution data center) |
| Entrust nShield | On-prem (any APAC institution data center); certain SaaS variants in Singapore and Tokyo |
| Utimaco SecurityServer | On-prem |

APAC institutions select the region matching their data-residency requirements; cross-border processing requires the institution's CC8.1 documentation of the cross-border posture under the relevant national law.

---

## 10. APAC regulatory incident-notification timing matrix

APAC institutions face simultaneous incident-notification clocks under multiple regulations. The institution's IR commander operates the **tightest applicable clock first** under the simultaneous-notification rule.

| Regulation / supervisor | Trigger event | Notification clock | Reporting authority | Chain artifact supporting the notification |
|---|---|---|---|---|
| MAS Notice 644 (Banking Act §60), Singapore | Significant cyber incident | 1 hour from awareness | MAS | Verifier output (spec §7); spec §10.2 operational events |
| MAS Notice 644 follow-up | Detailed report | 14 days | MAS | Full incident dossier including chain evidence |
| HKMA TM-G-1 | Cyber incident | "As soon as practicable" with CRAF 2.0 tiered timing | HKMA | Verifier output; spec §10.2 operational events |
| APRA CPS 234 §17, Australia | Information security incident | 72 hours from awareness | APRA | Verifier output; spec §10.2 operational events |
| RBI Cyber Security Framework, India | Cyber incident | 6 hours from detection | RBI | Verifier output; spec §10.2 operational events |
| RBI Master Direction on IT Outsourcing | Operational-risk incident affecting outsourcing | 2 hours | RBI | Verifier output; institution's vendor-management procedure |
| FSC Korea Electronic Financial Supervisory Regulation | Cyber incident | 24 hours | FSC / FSS | Verifier output; spec §10.2 operational events; see `korea-overlay.md` for KISA additional clock |
| Bank Negara Malaysia Risk Management in Technology (RMiT) | Operational incident affecting material services | "Immediately" with formal report within 24 hours | BNM | Verifier output; spec §10.2 operational events |
| Bank of Thailand IT/Cybersecurity Notification | Cyber incident | 24 hours | BoT | Verifier output |
| Bangko Sentral ng Pilipinas Circular 982 | Cyber incident | 2 hours from confirmation | BSP | Verifier output |
| Bank Indonesia IT Risk Management Regulation | Critical system failure | 24 hours | Bank Indonesia | Verifier output |
| Reserve Bank of New Zealand BS11 | Material technology incident | "As soon as practicable" with formal report within 72 hours | RBNZ | Verifier output |
| Singapore PDPA personal-data breach | Personal-data breach affecting >500 individuals | 72 hours | PDPC | Verifier output; institution's RoPA |
| Hong Kong PDPO personal-data breach | Personal-data breach (operational expectation) | "As soon as practicable" | PCPD | Verifier output |
| Indian DPDP Act 2023 personal-data breach | Personal-data breach | 72 hours | Data Protection Board of India | Verifier output |
| Australian Privacy Act personal-data breach | Eligible data breach | "As soon as practicable" — generally 30 days | OAIC | Verifier output |
| Korean PIPA personal-data breach | Personal-data breach | 72 hours | PIPC | Verifier output; see `korea-overlay.md` |

The matrix is an operational tool. The institution's IR runbook names the regulators applicable to its specific footprint; the runbook's ordered activation list reflects the tightest applicable clock first.

---

## 11. Pattern B per-jurisdiction tenant naming for APAC institutions

### 11.1 Pattern B as the recommended posture for multi-jurisdiction APAC institutions

Spec §10.15 Pattern B (per-region tenant_id) is the resilience pattern for institutions with per-region or per-jurisdiction event-isolation requirements. For APAC institutions operating across multiple national supervisors, Pattern B is the recommended posture: per-jurisdiction tenant_ids align the chain with the APAC supervisory topology and support examination-data-sovereignty rules under each national framework.

### 11.2 Per-jurisdiction tenant naming convention

The institution adopts a per-jurisdiction tenant_id naming convention that maps each national supervisor to a stable tenant identifier. The convention below is illustrative; institutions adopt their own naming as long as it is documented in CC8.1 and stable across the chain's retention period.

| Jurisdiction | National supervisor | Recommended tenant_id pattern | HSM region |
|---|---|---|---|
| Singapore | MAS | `bank-sg` | AWS CloudHSM `ap-southeast-1`, Azure Managed HSM Southeast Asia, Google Cloud HSM `asia-southeast1` |
| Hong Kong | HKMA | `bank-hk` | AWS CloudHSM `ap-east-1`, Azure Managed HSM East Asia, Google Cloud HSM `asia-east2` |
| Australia | APRA | `bank-au` | AWS CloudHSM `ap-southeast-2`, Azure Managed HSM Australia East / Central / Southeast, Google Cloud HSM `australia-southeast1` |
| New Zealand | RBNZ | `bank-nz` | AWS CloudHSM `ap-southeast-2` (or local variant), Azure Managed HSM Australia East, Google Cloud HSM `australia-southeast1` |
| India | RBI | `bank-in` | AWS CloudHSM `ap-south-1`, Azure Managed HSM India Central / West / South, Google Cloud HSM `asia-south1` |
| Korea | FSC / FSS | `bank-kr` | AWS CloudHSM `ap-northeast-2`, Azure Managed HSM Korea Central, Google Cloud HSM `asia-northeast3` |
| Malaysia | BNM | `bank-my` | AWS CloudHSM `ap-southeast-1` (closest), Azure Managed HSM Southeast Asia, Google Cloud HSM `asia-southeast1` |
| Indonesia | Bank Indonesia | `bank-id` | AWS CloudHSM `ap-southeast-3`, Azure Managed HSM Southeast Asia, Google Cloud HSM `asia-southeast2` |
| Thailand | BoT | `bank-th` | AWS CloudHSM `ap-southeast-1` (closest), Azure Managed HSM Southeast Asia, Google Cloud HSM `asia-southeast1` |
| Philippines | BSP | `bank-ph` | AWS CloudHSM `ap-southeast-1` (closest), Azure Managed HSM Southeast Asia, Google Cloud HSM `asia-southeast1` |

### 11.3 Cross-jurisdiction correlation under group consolidated supervision

For institutions operating under group consolidated supervision (Australian Tier-1 banks under APRA, Singapore-headquartered regional banks under MAS Group Supervision, Indian conglomerate-affiliated banks under RBI consolidated supervision), cross-tenant correlation reads as institution-side evidence. The institution's CC8.1 names the cross-jurisdiction correlation procedure that the institution operates internally; the chain's per-tenant HKDF binding (spec §4.1) means a tenant's events cannot mechanically be lifted into another tenant's chain.

### 11.4 ASEAN Banking Integration Framework

Cross-border ASEAN operations under the ABIF compose with Pattern B per-jurisdiction tenant naming. An institution operating under ABIF's Qualified ASEAN Banks framework deploys per-ASEAN-jurisdiction tenant_ids (`bank-sg`, `bank-my`, `bank-id`, `bank-th`, `bank-ph`); cross-jurisdiction correlation supports the ABIF's multi-jurisdiction-supervisor coordination expectations.

---

## 12. APEC Cross-Border Privacy Rules

### 12.1 The CBPR framework

The APEC Cross-Border Privacy Rules framework is a voluntary certification system administered by the APEC Data Privacy Subgroup. APEC member economies include Singapore, Hong Kong, Australia, Korea, the Philippines, Mexico, Canada, USA, Japan, and Taiwan. The CBPR framework provides a transfer mechanism for personal data flows between participating economies, parallel to the EU's adequacy framework.

### 12.2 Chain artifacts under CBPR

Chain artifacts containing personal data (or pseudonymised personal data) flowing between APEC economies operate under CBPR when the institution maintains APEC CBPR certification. The chain's per-region tenant_id pattern (Pattern B) supports CBPR's transfer-mechanism evidence requirements.

| CBPR requirement | Chain artifact supporting compliance |
|---|---|
| Identification of personal data subject to cross-border transfer | The institution's RoPA + the chain's `audit.adverse_action.*` schema (where applicable) |
| Notice and choice for individuals | Institution-side notice and consent management; the chain captures the consent reference if the institution's notice procedure binds consent IDs to decisions |
| Transfer-mechanism evidence | The institution's CBPR certification; the chain provides operational evidence the certification's posture is observed |
| Onward-transfer accountability | The DORA-overlay §2 lattice (the institution + SDK + receiver + HSM + LLM); the institution's third-party register names each downstream party |

### 12.3 Institutions without CBPR certification

Institutions without CBPR certification operate under bilateral national-law transfer mechanisms (Singapore PDPA's notification-of-transfer requirement, Hong Kong PDPO's voluntary code, Indian DPDP Act 2023's standard contractual clauses, Korean PIPA's Article 28 cross-border transfer rules). The spec is data-flow-mechanism-neutral; the institution's chosen mechanism is a CC8.1 attribute.

---

## 13. Translation table — APAC supervisory expectations to chain artifacts

The translation table is consumed by the institution's compliance officer mapping the chain's deliverables to each APAC supervisor's specific examination expectations.

| Spec section | MAS | HKMA | APRA | RBI | FSC Korea | Conformance status |
|---|---|---|---|---|---|---|
| §4.1 Per-event MAC | TRM Guidelines §6.4 (cryptographic controls) | TM-G-1 §4.2 (cryptographic key management) | CPS 234 §15(b) (integrity) | Cyber Security Framework §6.1 | EFSR §17 (key management); see `korea-overlay.md` | CONFORMANT |
| §4.2 Daily Merkle seal | TRM Guidelines §6.5 (audit logging) | TM-G-1 §5.1 (audit log integrity) | CPS 234 §15(b) | Cyber Security Framework §7.3 | EFSR §17 (logging) | CONFORMANT |
| §4.3 HSM signature | TRM Guidelines §6.4 | TM-G-1 §4.2 | CPS 234 §15(b) | Cyber Security Framework §6.1 | EFSR §17; see `korea-overlay.md` | CONFORMANT |
| §4.4.1 Routing decisions | FEAT Accountability principle | 2024 AI principles §3 (Transparency) | CPS 230 §3 (operational risk identification) | November 2023 RBI letter §3 (AI governance) | AI Basic Act 2025 Article 17; see `korea-overlay.md` | PARTIAL (institution's CC8.1 names routing-policy versioning) |
| §4.4.2 Deployment intent | FEAT Accountability + Transparency | 2024 AI principles §5 (Human oversight) | CPS 230 §3 | November 2023 RBI letter §3 | AI Basic Act 2025 Article 18 | PARTIAL (institution's CC8.1 names deployment-policy versioning) |
| §6 Deployment topologies | TRM Guidelines §3 (governance) | TM-G-1 §3 (governance) | CPS 234 §17 | Master Direction Outsourcing §3 | EFSR §17 + 망분리 (see `korea-overlay.md`) | CONFORMANT |
| §7 Verifier | TRM Guidelines §6.5 | TM-G-1 §5.3 | CPS 234 §18 (internal audit) | Master Direction Outsourcing §7 (supervisor access) | EFSR §17 | CONFORMANT |
| §10.5 HSM custody | TRM Guidelines §6.4 | TM-G-1 §4.2 | CPS 234 §15(b) | Cyber Security Framework §6.1 | EFSR §17 + K-ISMS-P; see `korea-overlay.md` | CONFORMANT |
| §10.9 Retention | TRM Guidelines §6.5 | TM-G-1 §5.3 | CPS 234 §17 | Cyber Security Framework §7.3 | FEFTA + FIEA; see `korea-overlay.md` | CONFORMANT |
| §10.11 Adverse-action notice | FEAT Transparency principle | 2024 AI principles §3 | CPS 230 §3 | RBI Master Direction on Customer Service | AI Basic Act 2025 Article 17 + adverse-action notice; see `korea-overlay.md` | PARTIAL (institution uses `audit.adverse_action.regulatory_basis = sg-feat / hk-hkma-ai-principles / au-cps-230 / in-rbi-2023 / kr-fsc` per jurisdiction) |
| §10.14 Trusted-time integration | TRM Guidelines §6.5 | TM-G-1 §5.1 | CPS 234 §15(b) | Cyber Security Framework §7.3 | EFSR §17 | PARTIAL (institution adopts RFC 3161 with a TSA in the relevant jurisdiction; APAC equivalent of EU-Trusted-List is institution-selected) |
| §10.15 Multi-region resilience | TRM Guidelines §3 + §6.6 | TM-G-1 §6 (reconciliation) | CPS 230 §7 (BCP) | Master Direction Outsourcing §6 | EFSR §17 + BCP; see `korea-overlay.md` | CONFORMANT (Pattern B per-jurisdiction recommended) |

---

## 14. Operational checklist for an APAC institution deploying the chain

The institution working through MAS / HKMA / APRA / RBI / FSC compliance produces the following artefacts. This checklist names the minimum set; the institution's compliance department extends it as their specific exposure requires.

| Artefact | Source | Reviewed by |
|---|---|---|
| MAS FEAT Principles mapping (per §2 of this overlay) | Institution's CC8.1 | MAS Technology Risk supervisor |
| MAS Veritas assessment composition (per §3 of this overlay) | Institution's CC8.1 + Veritas-assessment-firm engagement | MAS during AIDA approval |
| HKMA TM-G-1 mapping (per §4 of this overlay) | Institution's CC8.1 | HKMA Banking Supervision examiner |
| APRA CPS 234 / CPS 230 mapping (per §5 of this overlay) | Institution's CC8.1 | APRA Cyber and Operational Risk reviewer |
| RBI guidelines mapping + DPDP Act 2023 cross-border posture (per §6 of this overlay) | Institution's CC8.1 + RoPA | RBI Cyber Security and IT Examination team |
| FSC Korea mapping (cross-reference to `korea-overlay.md`) | Institution's CC8.1 | FSS comprehensive examination (종합검사) |
| APAC cloud HSM region selection (per §9) | Institution's procurement + CC8.1 | Institution's IT security |
| APAC incident-notification timing matrix (per §10) | Institution's IR runbook | National supervisor for each jurisdiction the institution operates in |
| Pattern B per-jurisdiction tenant naming (per §11) | Institution's CC8.1 | National competent authorities for each jurisdiction |
| APEC CBPR posture (per §12) | Institution's CBPR certification + RoPA | Institution's DPO + APEC Data Privacy Subgroup |
| Per-jurisdiction `audit.adverse_action.regulatory_basis` enumeration | Institution's CC8.1 | Institution's consumer-protection counsel |

---

## 15. Bottom line

The chain v1.0a is conformant with APAC supervisory architecture for the jurisdictions enumerated in this overlay. The conformance is operational: the institution's CC8.1 names the per-supervisor mapping, the per-jurisdiction tenant naming, the cloud HSM region, and the incident-notification timing per the matrix in §10; the spec produces the integrity-bearing audit trail under any of these institutional choices.

The chain's strongest item from an APAC perspective is the byte-reproducible verifier output and the test-vector corpus; APAC supervisors are receptive to demonstrable test artefacts and to integrity-bearing audit trails over assertion-based control evidence. The principle-based regulatory style of MAS, HKMA, APRA, RBI, and FSC composes well with the chain's principle-style invariant articulation in spec §4.1; the institution's CC8.1 names the per-supervisor terminology and the chain produces the operational evidence under any terminology.

For multi-jurisdiction APAC Tier-1 institutions, Pattern B per-jurisdiction tenant naming with HSM custody pinned to the corresponding APAC region is the recommended posture. This composes with the ASEAN Banking Integration Framework, the SEACEN coordination forum, and the APEC CBPR cross-jurisdiction posture without forcing the institution to operate dual chain implementations or write per-supervisor conformance overlays from scratch.

For Korean-supervised institutions, this overlay is the regional cross-jurisdiction articulation; the deeper Korea-specific articulation is in `korea-overlay.md`. For Japanese-supervised institutions, the Japan regulator-pack documents (separate from this APAC overlay) provide the FSA / BoJ / APPI / J-SOX / CRYPTREC articulation.

---

## Document control

| Field | Value |
|---|---|
| Document | APAC Articulation Overlay |
| Version | 1.0.0 |
| Status | Informative — institution-side articulation; no normative spec change |
| Aligned with | MAS / HKMA / APRA / RBI / FSC Korea / ASEAN ABIF / APEC CBPR |
| Date | 2026-05-07 |
