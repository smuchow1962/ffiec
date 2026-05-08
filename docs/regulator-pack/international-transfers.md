# International Transfers — Schrems II Compliance for the Chain

> **What this doc is.** Procedures and topology guidance for institutions that operate the chain across borders. Written to answer the GDPR Articles 44–50 question — when is the institution lawfully transferring personal data, what supplementary safeguards does Schrems II require, and how does the chain's tokenization interact with the transfer rules.

> **Why this matters.** *Schrems II* (CJEU C-311/18, July 2020) invalidated the EU-US Privacy Shield and forced controllers to verify, on a per-transfer basis, that supplementary safeguards bring third-country protection up to EU standards. The European Data Protection Board (EDPB) Recommendations 01/2020 lay out the methodology. An institution that cannot show a Transfer Impact Assessment (TIA) for each non-adequate-country flow is non-compliant the moment a customer or DPA asks.

---

## 1. The transfer rules at a glance

GDPR Article 44 sets the general principle: personal data transferred outside the EEA must benefit from an adequate level of protection. Articles 45–49 enumerate the lawful grounds:

| Mechanism | Article | When usable |
|---|---|---|
| Adequacy decision | 45 | Country has been formally found adequate (UK, Switzerland, Japan, South Korea, EU-US Data Privacy Framework signatories) |
| Standard Contractual Clauses (SCCs) | 46(2)(c) | Most common; controller and processor sign EU-Commission-approved clauses |
| Binding Corporate Rules (BCRs) | 47 | Multinational with internal data-flow rules approved by lead DPA |
| Derogations | 49 | Limited, case-by-case (consent, contract necessity, important reasons of public interest) |

For US-bound transfers, the post-Schrems-II expectation is **SCCs + Transfer Impact Assessment + supplementary safeguards** unless the recipient is certified under the EU-US Data Privacy Framework (DPF) and the data category is in scope. The institution maintains a TIA for each non-adequate-country flow regardless of the legal mechanism chosen.

---

## 2. Four deployment scenarios

The chain can be deployed in topologies that range from "no transfer" to "vendor-hosted multi-tenant across regions." Each scenario has a different transfer profile.

### 2.1 Scenario A — EU on-prem

The institution is an EU entity. All chain components (SDK, ledger, privacy-store, HSM) run in EU data centers under the institution's direct control. No personal data leaves the EEA.

| Element | Status |
|---|---|
| Transfer occurring? | No |
| SCCs required? | No |
| TIA required? | No |
| Adequacy considerations | n/a |
| Cloud-vendor sub-processor risk | n/a if no cloud is used; if EU-region cloud is used, see Scenario B caveats below |

This is the simplest topology. No further transfer paperwork is needed. The privacy-store and chain ledger are within EU jurisdiction; data-subject rights, supervisory authority access, and incident response all run under EU law without conflict.

### 2.2 Scenario B — EU data, US cloud

The institution stores chain data on a US-headquartered cloud provider (AWS, Azure, GCP) even when the storage region selected is "EU." US cloud providers remain subject to US legal process (FISA §702, Executive Order 12333, the CLOUD Act). Schrems II identifies this category as the central risk.

| Element | Status |
|---|---|
| Transfer occurring? | Yes (the US cloud provider is a US entity, even when the storage region is EU) |
| SCCs required? | Yes — between the EU institution (controller) and the US cloud provider (processor) |
| TIA required? | Yes — must address FISA §702 and EO 12333 surveillance risk per EDPB Recommendations 01/2020 |
| Supplementary safeguards | Encryption with EU-held keys (the cloud provider cannot access plaintext); strict access logging; confidential computing where supported |
| Sub-processor authorization | Required per Article 28(2) and (4); cloud provider must publish sub-processor list and notify on changes |

The institution's TIA documents:

1. The category of data transferred (chain entries containing pseudonymized personal data; privacy-store mappings if also hosted in US cloud).
2. The legal basis for the transfer (typically Article 46(2)(c) SCCs).
3. The third-country legal environment (US surveillance law).
4. The likelihood that the recipient's data is accessed by US authorities for the data category in question.
5. The supplementary safeguards in place.
6. The institution's conclusion on whether the transfer offers an essentially equivalent level of protection.

The institution's preferred safeguard is **EU-held keys**: the chain confidentiality key (AES-GCM at rest) and the privacy-store key (token reversibility) are held in an HSM physically located in the EU and operated by an EU subsidiary. The US cloud provider stores ciphertext only and cannot decrypt without access to the EU-held keys. Under this configuration, even if US authorities serve legal process on the cloud provider, the provider cannot produce plaintext — a posture explicitly recommended by EDPB Recommendations 01/2020 §85 (Use Case 1: Data storage for backup and other purposes that do not require access to data in the clear).

### 2.3 Scenario C — Multi-region active-active

The institution operates the chain across multiple regions (e.g., EU + US) for resilience and latency. Each region has its own ledger, privacy-store, and HSM. Cross-region replication is encrypted in transit and at rest.

| Element | Status |
|---|---|
| Transfer occurring? | Yes — cross-region replication transfers data from one jurisdiction to another |
| SCCs required? | Yes — between the EU controller and the US subsidiary (or other non-EU entity); intra-group transfers can use BCRs if approved |
| TIA required? | Yes for each non-adequate-country leg |
| Supplementary safeguards | Encryption at rest with region-specific keys held by the regional entity; replication ciphertext-only; data-subject rights honored at the originating region (EU rights at EU region) |
| DPA Article 28 | Required between EU controller and US-subsidiary processor; DPA names sub-processors, security measures, breach notification |

The institution's preferred design is to **store originating-region data only in the originating region**: EU customer data lives only in the EU region; US customer data lives only in the US region. The active-active topology then operates per-region, with no cross-region replication of personal data. Where cross-region replication is operationally necessary (e.g., for global anti-fraud detection), it is encrypted with regional keys such that the receiving region cannot decrypt without explicit cooperation.

### 2.4 Scenario D — Vendor-hosted multi-tenant

A SaaS vendor hosts the chain for multiple institutional customers. Each customer is a controller; the vendor is a processor (per `controller-processor-determination.md` if the institution has produced one, or per Article 28 default).

| Element | Status |
|---|---|
| Transfer occurring? | Depends on the vendor's data-residency commitments and the location of the controller |
| SCCs required? | Yes if the vendor or its sub-processors are outside the EEA |
| TIA required? | Yes per non-adequate-country leg |
| DPA | Required per Article 28(3); vendor commits to compliance, security measures, sub-processor authorization, breach notification, return/deletion at end of service |
| SOC 2 Type II | Vendor's SOC 2 Type II covers Article 32 controls; supplements DPA |
| Tenant isolation | Each customer's IKM is custody-separate per spec §4.1; vendor must commit to no co-mingling |

For multi-tenant vendor-hosted deployments, the vendor's DPF certification (if certified) covers transfers to the vendor for in-scope categories. The TIA still must address sub-processors and back-end infrastructure that may not be DPF-covered.

---

## 3. Standard Contractual Clauses — operational use

The 2021 SCCs (Commission Implementing Decision (EU) 2021/914) are modular: the institution selects the module matching the relationship.

| Module | Use case |
|---|---|
| Module 1 | Controller to controller |
| Module 2 | Controller to processor |
| Module 3 | Processor to processor (sub-processor) |
| Module 4 | Processor to controller |

For chain deployments:

- **EU institution to US cloud provider** — Module 2 (controller to processor).
- **EU institution to US subsidiary running the chain** — Module 2 if the subsidiary is processor; Module 1 if the subsidiary is co-controller (uncommon).
- **EU institution to vendor-hosted multi-tenant SaaS** — Module 2.
- **SaaS vendor to its US sub-processors** — Module 3 (processor to processor).

The institution's legal team executes the SCCs at the start of the relationship and re-executes at any material change. The signed SCCs are part of the institution's records of processing activities (RoPA) entry for the chain.

---

## 4. Transfer Impact Assessment — six-step methodology

The EDPB Recommendations 01/2020 lay out a six-step TIA. The institution operates them as a checklist:

| Step | Question | Output |
|---|---|---|
| 1. Know your transfers | What personal data is transferred, to which third country, to which recipient, for what purpose? | Transfer inventory entry |
| 2. Identify transfer tool | Adequacy decision? SCCs? BCRs? Derogation? | Article 46/47/49 mechanism named |
| 3. Assess third-country law | Does the third country's law (especially surveillance) ensure essentially equivalent protection? | Country-law analysis (institution may rely on EDPB country reports for major countries) |
| 4. Identify supplementary measures | If protection is not essentially equivalent, what technical/organizational/contractual measures bring it up? | Supplementary measures list (encryption with EU-held keys, pseudonymization, contractual restrictions on government access) |
| 5. Procedural steps | Adopt the supplementary measures; document the analysis | TIA document signed by DPO |
| 6. Re-evaluate | Periodically re-assess (annually or on triggering event) | Re-assessment record |

The institution's TIA template lives with the privacy office; this document references it but does not duplicate it.

---

## 5. DPA — Data Processing Agreement under Article 28

Article 28 governs the controller-processor relationship. The DPA must include the elements listed in Article 28(3):

| Element | What the DPA says |
|---|---|
| Subject matter | Hosting and operating the chain on behalf of the controller |
| Duration | Same as service contract |
| Nature and purpose | Integrity-bearing audit trail of automated decisions |
| Type of personal data | Pseudonymized customer identifiers, decision metadata; possibly PHI in healthcare deployments |
| Categories of data subjects | Customers, applicants, beneficiaries |
| Obligations of controller | Configure SDK; manage privacy-store; provide lawful instructions |
| Obligations of processor | Process only on documented instructions; ensure confidentiality of personnel; security per Article 32; sub-processor authorization; assist with data-subject rights; assist with breach notification; return/delete at end; allow audits |
| Sub-processor authorization | General or specific authorization; processor lists sub-processors and notifies of changes |
| International transfers | SCCs incorporated by reference for any onward transfer |

For vendor-hosted multi-tenant deployments, the vendor's standard DPA is reviewed by the institution's privacy office before signature. Vendor-proposed deviations from Article 28 minimums are documented and approved or rejected.

---

## 6. Tokenization and the transfer question

A frequent question: "Do chain tokens count as personal data for the transfer rules?" The answer turns on `pseudonymization-vs-anonymization.md`:

- During the retention period, tokens are pseudonymized data — still GDPR personal data — because the privacy-store mapping is retained and the institution can reverse them. Transfer rules apply.
- After the privacy-store mapping is deleted (end of retention or post-erasure), the tokens are anonymized — no longer GDPR personal data — and transfer rules do not apply.

The institution's transfer paperwork therefore covers the chain as personal-data transfer for the entire 7-year retention period. Post-de-mapping, transfers of the chain tokens are unrestricted under GDPR (other regulations may still apply).

---

## 7. Cross-references

- GDPR Articles 44–50 — operationalized here.
- `pseudonymization-vs-anonymization.md` — explains why chain tokens are personal data during retention.
- `multi-jurisdiction-conflict.md` — when transfers are demanded by US legal process and forbidden by GDPR.
- `dpia-template.md` §3 — risk analysis incorporates cross-border risk.
- `ropa-template.md` "Cross-border transfers" column references this document.
- `token-vault-architecture.md` — separate custody for the privacy-store key supports the EU-held-keys safeguard.
- `breach-notification-matrix.md` — cross-border incidents implicate multiple supervisory authorities.

---

## 8. Review cadence

This document is reviewed annually by the privacy office, legal team, and security architecture. Triggers for early review:

- A new adequacy decision (e.g., a country gains or loses adequacy).
- A CJEU decision affecting transfer mechanisms (e.g., a Schrems III).
- The institution adds or removes a deployment region.
- A US legal-process event (a CLOUD Act or FISA order) demonstrates that a supplementary safeguard is or is not effective.
- A vendor changes sub-processor list in a way that implicates new third countries.
