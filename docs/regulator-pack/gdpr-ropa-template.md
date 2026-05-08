# GDPR Article 30 Records of Processing Activities — chain template

> **What this doc is.** A populated Records of Processing Activities (RoPA) entry for the chain, in the format required by GDPR Article 30. The template is institution-ready: drop the institution's name into the controller field, fill in the vendor name if SaaS-hosted, and the rest is reusable. Article 30 requires every processing activity to be documented; without a current RoPA entry the institution cannot answer "what processing are you doing and why?" when a Data Protection Authority asks. This is the answer for the chain.

## Why the RoPA matters

Article 30(1) requires controllers to maintain a record of processing activities under their responsibility. The record must be available to the supervisory authority on request. It is not a one-time deliverable — it is a living document that must be current. An out-of-date RoPA is a documented Article 30 finding.

The RoPA is also the institution's internal map. Privacy reviewers, auditors, the DPO, and the legal team all consult the RoPA to understand what data is processed for what purpose. When a data subject asks "what processing affects me?" the RoPA is the index that lets the privacy team answer comprehensively.

The chain is one processing activity among many. This template is its entry. The institution may have dozens or hundreds of other entries covering core banking, marketing, HR, vendor management, and so on. The chain's entry sits alongside them in the institution's RoPA database.

## The template

| Field | Content |
|---|---|
| **Processing activity name** | FFIEC Chain-of-Custody for AI Decision Logging |
| **Description** | Tamper-evident audit trail of AI-driven decisions affecting credit, fraud determination, and pricing. The chain captures decision inputs, model parameters, and outputs in an integrity-bearing append-only log to support fairness audits, customer disputes, regulatory examination, and litigation defense. |
| **Purpose(s)** | (1) Compliance with banking-regulation retention requirements (FFIEC, ECOA, CFPB supervisory). (2) Fairness audit and disparate-impact detection. (3) Model-risk management oversight per SR 11-7 / OCC 2011-12. (4) Customer-dispute response per `customer-dispute-procedures.md`. (5) Litigation evidence per `litigation-support.md`. |
| **Legal basis (Article 6)** | Article 6(1)(c) legal obligation (FFIEC banking regulations 12 CFR 1002 and equivalent OCC/Fed/FDIC/CFPB retention authorities) for retention and primary use. Article 6(1)(f) legitimate interests (fairness audit, MRM oversight, dispute defense) for supplementary uses, with balancing test documented in `gdpr-lawful-basis.md`. |
| **Special-category basis (Article 9)** | When special-category data is in scope: Article 9(2)(b) for lending under ECOA's social-protection function, 9(2)(g) for AML/fraud public interest, 9(2)(h) for healthcare-financing decisions. Field-by-field mapping in the institution's DPIA (`gdpr-dpia-template.md`). |
| **Data controller** | `[Institution legal name and registered address]`. The controller determines the purposes and means of the chain processing. |
| **Joint controller (if any)** | Typically none. If the institution operates a joint AI program with another regulated entity, name the joint controller and reference the Article 26 joint-controller agreement. See `gdpr-controller-vs-processor.md`. |
| **Data Protection Officer** | `[DPO name, contact email, phone, registered address]`. |
| **Data processor (if any)** | If the chain is vendor-hosted: `[Vendor name, registered address, lead contact]`. The vendor processes on behalf of the controller under a Data Processing Agreement (DPA) per Article 28. If the chain is on-premises and operated internally, this field is "internal — no external processor." See `gdpr-controller-vs-processor.md`. |
| **Sub-processors (if vendor-hosted)** | Listed in the vendor's sub-processor schedule, with the institution's authorization recorded per Article 28(2) and 28(4). Typical sub-processors: cloud infrastructure provider (AWS / Azure / GCP), HSM service, observability backend. |
| **Categories of data subjects** | Customers and applicants whose AI-driven decisions are captured. May include account holders, loan applicants, fraud-flagged transaction parties, and (in healthcare-financing contexts) patients whose payment products are AI-priced. |
| **Categories of personal data** | Customer identifiers (tokenized in the chain; full identifier in the privacy-store). Decision-relevant attributes (income, employment, credit history, behavioral signals — tokenized or redacted per the institution's privacy-by-design configuration). Model inputs (free-text fields tokenized or regex-redacted; structured fields tokenized where they carry PII). Decision outputs (approve/deny, risk score, pricing). Model metadata (version, parameters, prompt template). Audit-context attributes (timestamp, run identifier, tenant identifier, key version). |
| **Special-category data** | When in scope: health data (healthcare-financing), inferred attributes that may proxy for protected classes (age, marital status, ethnicity proxies), biometric authentication signals. Tokenization or omission is applied per the privacy-by-design configuration. |
| **Recipients (internal)** | Compliance team (audit and reporting), Model Risk Management (MRM oversight), Customer Service (dispute handling), Litigation Counsel (under privilege, when active), Internal Audit (control testing). Access controlled by the institution's role-based access policy. |
| **Recipients (external)** | Banking regulators (CFPB, OCC, Federal Reserve, FDIC, NCUA) under examination authority. Courts and opposing counsel under discovery (when applicable). External auditors under the institution's audit-engagement scope. Data Protection Authorities under Article 31 cooperation. |
| **Retention period** | 7 years from event capture per 12 CFR 1002 (ECOA records) and equivalent federal-regulator retention standards. Extended by litigation hold per the institution's legal-hold procedures (`incident-response-playbook.md` litigation-hold sections). Privacy-store mappings deleted at end of retention period unless documented hold is in place; chain entries become anonymized at that point per `gdpr-pseudonymization-vs-anonymization.md`. |
| **Cross-border transfers** | Address each storage and processing location. If processing remains within EU/EEA: no transfer mechanism required. If chain is stored or processed in a third country (e.g., US cloud region): transfer mechanism is Standard Contractual Clauses (SCCs) per the EU Commission's 2021/914 implementing decision, plus a Transfer Impact Assessment (TIA) addressing Schrems II concerns. Tokens-in-chain (post-de-mapping) are not personal data and are not subject to transfer restrictions. Privacy-store mappings remain in EU jurisdiction where the data subject resides in EU/EEA. See the institution's cross-border-transfer assessment for full detail. |
| **Security measures** | Technical: HMAC-SHA-256 per-entry integrity (spec §4.1), daily Merkle seal (spec §4.2), HSM-backed Ed25519 signature (spec §4.3), encryption at rest with separate confidentiality key, TLS 1.3 in transit, HSM-protected master keys per spec §10. Organizational: separation of duties (chain-operations and privacy-store custody are different teams), access logging, key-fingerprint reconciliation per spec §10.1 P-6, incident response per `incident-response-playbook.md`, annual penetration testing, annual DR exercise, annual DPIA review. |
| **Data-subject rights procedures** | Article 15 (access): customer-correlation index plus privacy-store extract per `gdpr-dsar-fulfillment.md`. Article 16 (rectification): correction-record chain entry plus privacy-store update per `gdpr-article-16-rectification.md`. Article 17 (erasure): privacy-store mapping deletion per `gdpr-article-17-procedures.md`. Article 18 (restriction): processing flag in the privacy-store; chain capture continues but downstream use is restricted. Article 20 (portability): privacy-store extract in machine-readable format; chain entries excluded as institution-internal evidence. Article 21 (objection): legitimate-interests balancing review by the privacy team. Article 22 (automated decision-making): human review per `customer-dispute-procedures.md`. |
| **Data Protection Impact Assessment** | Completed per Article 35; current version dated `[date]`. Stored at `[location]`. Reviewed annually. See `gdpr-dpia-template.md` for the template content. |
| **Date created** | `[date the RoPA entry was first created]` |
| **Date last reviewed** | `[date of most recent review by the DPO or privacy team]` |
| **Next review due** | `[date]` (annual review cadence; sooner if material change occurs) |

## How to use the template

1. **First-time deployment.** Before the chain goes into production, the privacy team populates this template, the DPO reviews it, and the entry is added to the institution's RoPA database. This is part of the deployment-approval pack the AI Governance Committee receives.

2. **Material change triggers review.** If any of the following occurs, the privacy team reviews and updates the entry within 30 days: new AI use case enrolled in the chain, new data category captured, change in vendor or sub-processor, change in cross-border topology, change in retention period (e.g., a new regulator-mandated extension), incident-response update affecting the security-measures field.

3. **Annual review.** The DPO confirms the entry is current. The review timestamp updates the "Date last reviewed" field.

4. **Authority response.** If a Data Protection Authority requests the institution's RoPA, the chain entry is exported in the institution's standard RoPA format (typically the institution's GRC tool's CSV or JSON export). The institution responds within the authority's deadline (typically 30 days under Article 31).

## Worked example — populated entry for a hypothetical institution

| Field | Content |
|---|---|
| **Processing activity name** | FFIEC Chain-of-Custody for AI Decision Logging |
| **Description** | (As above) |
| **Data controller** | Acme Bank N.A., 100 Corporate Way, Pittsburgh PA 15222, US (EU operations: Acme Bank Europe GmbH, Mainzer Landstrasse 17, 60329 Frankfurt am Main, Germany) |
| **Joint controller** | None |
| **Data Protection Officer** | Maria Garcia, dpo@acmebank.eu, +49 69 1234 5678, Mainzer Landstrasse 17, 60329 Frankfurt am Main, Germany |
| **Data processor** | Vidimus Inc., 200 Vendor Plaza, San Francisco CA 94107, US (DPA executed 2025-11-12, sub-processor schedule version 4) |
| **Categories of data subjects** | Loan applicants and account holders in Acme Bank N.A. and Acme Bank Europe GmbH; healthcare-financing patients enrolled in MedPay program |
| **Categories of personal data** | (As above; institution adds specific fields per its privacy-by-design configuration) |
| **Recipients** | Acme Compliance, Acme MRM, Acme Customer Service, Acme Litigation Counsel; CFPB, OCC, BaFin, ECB; external auditor `[name]`; courts under discovery |
| **Retention period** | 7 years from event capture, extended by active litigation holds (currently 2 active matters) |
| **Cross-border transfers** | EU customer data: stored and processed within Frankfurt region; no transfer. US customer data: stored in US-East-1 region, no transfer. Vendor backups encrypted with EU-held keys per Schrems II supplementary safeguard. TIA dated 2025-09-15. |
| **DPIA** | Version 2.1, dated 2026-02-14. Stored in GRC tool. Next review 2027-02-14. |
| **Date created** | 2025-08-22 |
| **Date last reviewed** | 2026-04-30 |
| **Next review due** | 2027-04-30 |

The example shows how the template absorbs an institution's specific facts. Each field is concrete; nothing is speculative.

## Common errors

- **Listing the chain as a security measure rather than a processing activity.** The chain is a processing activity in its own right (it captures and stores personal data). The integrity-bearing controls are the security measures the activity implements. The RoPA entry must treat them at the right level.
- **Naming "consent" as the basis.** Consent is not the basis (see `gdpr-lawful-basis.md`). Naming consent here will trigger withdrawal claims that the institution cannot honor under the legal-obligation retention.
- **Omitting the vendor.** If a vendor hosts the chain, the vendor is a processor under Article 28. The RoPA entry must name the vendor and reference the DPA.
- **Stale retention.** If the institution extends retention by litigation hold, the RoPA entry's retention field must reflect that. "7 years" alone does not capture an active hold.
- **No DPIA reference.** Article 30 and Article 35 are linked: an Article 35 DPIA exists for high-risk processing, and the RoPA references it. Omitting the reference signals incomplete documentation.

## Cross-references

- `gdpr-lawful-basis.md` — the Article 6 and 9 analysis cited in the legal-basis fields
- `gdpr-dpia-template.md` — the DPIA whose existence is recorded in the DPIA field
- `gdpr-controller-vs-processor.md` — the controller / processor / joint-controller determination
- `gdpr-dsar-fulfillment.md`, `gdpr-article-16-rectification.md`, `gdpr-article-17-procedures.md` — the data-subject-rights procedures cited
- `gdpr-article-25-demonstrability.md` — the demonstrability artifacts that support the security-measures field
- `privacy-by-design.md` — the tokenization configuration cited in the categories-of-data field
- spec §10 — the technical security measures the entry references
