# GDPR controller vs processor — vendor-hosted multi-tenant determination

> **What this doc is.** The institution's framework for determining controller / processor / joint-controller status when the chain is vendor-hosted, and the contractual and operational obligations that follow. The default recommended structure is "vendor as processor under Article 28," with the institution as sole controller. The joint-controller structure under Article 26 is available but adds complexity and is rarely the right fit. This doc names both, recommends one, and lists the DPA contents the institution must require either way.

## Why the determination matters

GDPR distinguishes controllers (Article 4(7)) from processors (Article 4(8)). The controller determines purposes and means of processing; the processor processes on the controller's behalf under documented instructions. The two roles carry different obligations:

| Obligation | Controller | Processor |
|---|---|---|
| Lawful basis (Article 6) | Yes — must establish | No — relies on controller's basis |
| Article 13/14 transparency | Yes — must inform data subjects | No — controller informs |
| Data-subject rights (Articles 15-22) | Yes — primary respondent | Assists controller |
| Article 30 RoPA | Article 30(1) record | Article 30(2) record (different format) |
| Article 33-34 breach notification | Yes — to DPA and data subjects | Notifies controller within 24-48 hours |
| Article 35 DPIA | Yes — required for high-risk | Assists controller |

When two parties jointly determine purposes and means, they are joint controllers under Article 26 and share controller obligations per a joint-controller arrangement document.

The wrong determination at the start creates downstream problems. A vendor structured as "processor" but acting as a joint controller exposes both parties to liability mismatch. A vendor structured as "joint controller" but acting only as processor over-burdens the vendor with controller obligations it has no basis to satisfy. Both parties need to name the role honestly and document the determination.

## The two structures

### Recommended — vendor as processor (Article 28)

In this structure:

- The institution (bank) is the sole controller.
- The vendor is the processor, processing chain data on the institution's documented instructions.
- The relationship is governed by a Data Processing Agreement (DPA) per Article 28(3).
- Each tenant on the vendor's platform is under the controller-instructions of the bank that owns the tenant.

The institution determines what gets captured (configuration), what gets tokenized (privacy-by-design.md SDK configuration), how long data is retained (subject to FFIEC mandate), who has access. The vendor implements the chain per the spec and the institution's configuration. The vendor does not determine purposes; it executes.

This structure is recommended because it matches the operational reality of a vendor-hosted chain platform. The vendor provides the platform; each bank-customer determines what to do with it. The vendor's commitments are about service-delivery, not data-purpose.

### Alternative — joint controllers (Article 26)

In this structure:

- The institution and the vendor are joint controllers.
- They jointly determine the purposes and means of processing (the vendor proposes tokenization scope, the bank approves; they jointly determine retention; they jointly determine analytic features that operate on the chain).
- Their relationship is governed by a joint-controller arrangement under Article 26(1).

Joint-controller arrangements name the responsibilities of each party for: data-subject-rights fulfillment, Article 13/14 transparency, breach notification, data-subject contact point. The arrangement is summarized publicly per Article 26(2) so data subjects know whom to contact.

Joint-controller structure is rare for chain-hosting because it requires the vendor to participate in purpose-determination. Most vendors do not want this role; it expands their regulatory exposure beyond service-delivery. The structure is appropriate when the vendor offers analytic services that go beyond the spec — for instance, cross-customer fairness benchmarking — and the vendor is named in those analytic outputs. The recommendation is to keep the vendor in processor role and place any cross-customer analytic services in a separate, opt-in arrangement.

## The Article 28 DPA — required contents

When the structure is vendor-as-processor, Article 28(3) requires the DPA to set out in writing:

| Article 28(3) requirement | DPA content |
|---|---|
| Subject matter and duration | Chain-hosting service; for the duration of the master service agreement |
| Nature and purpose of processing | Capture, store, and provide access to AI-decision audit-trail records per the chain spec; institution-internal use plus regulator response |
| Type of personal data | Customer identifiers (tokenized in the chain), decision-relevant attributes (tokenized or omitted per institution configuration), free-text model inputs and outputs (regex-redacted), decision outputs, model metadata, audit context |
| Categories of data subjects | Customers and applicants of the institution whose AI-driven decisions are captured |
| Obligations and rights of the controller | Institution determines configuration, instructions, and disposition of data; institution responds to data-subject requests; institution informs vendor of regulatory inquiries |
| Article 28(3)(a) — process only on documented instructions | Vendor commits to processing only per institution's instructions, with named exceptions for compliance with EU/Member State law (vendor must inform institution before exceptional processing) |
| Article 28(3)(b) — confidentiality of personnel | Vendor's personnel access chain data only under written confidentiality obligation; access logged |
| Article 28(3)(c) — security measures | Per Article 32: encryption at rest, encryption in transit, access controls, audit logging, key management. Vendor commits to specific controls inventoried in the DPA security schedule |
| Article 28(3)(d) — sub-processor authorization | Vendor maintains a sub-processor list, notifies institution of intended changes with 30-day advance notice; institution may object |
| Article 28(3)(e) — assist with data-subject rights | Vendor provides technical means to extract, modify (correction-record), or delete privacy-store records on institution's request |
| Article 28(3)(f) — assist with controller's obligations | Vendor assists with Article 32 security, Article 33-34 breach notification, Article 35 DPIA, Article 36 prior consultation |
| Article 28(3)(g) — return or delete on termination | Upon termination of the service agreement, vendor returns or deletes chain data per institution's choice; vendor provides deletion certification |
| Article 28(3)(h) — make available compliance information | Vendor provides SOC 2 Type II reports, audit attestations, and compliance information on request |

The DPA also addresses operational specifics:

| DPA section | Content |
|---|---|
| Sub-processor list | Specific sub-processors named with their geographic location and processing scope |
| Sub-processor authorization procedure | How institution authorizes new sub-processors; 30-day window; institution's right to object; vendor's obligation to find alternative if objection sustained |
| Security measures schedule | Specific controls (encryption ciphers, key management, access controls, monitoring) inventoried with control-objective references |
| Audit rights | Institution's right to audit vendor's compliance, frequency, scope, vendor's commitment to provide reasonable cooperation |
| Breach notification | Vendor commits to notify institution within 24 hours of detection, providing the information the institution needs to comply with Article 33-34 timelines |
| International transfers | Standard Contractual Clauses (2021/914 Module Two — controller to processor) attached as schedule; Transfer Impact Assessment referenced |
| Data location | Specific cloud regions or data centers where chain data is processed and stored |
| Data isolation | Vendor commits to multi-tenant isolation: no co-mingling of data across institutions; tenant-level access controls; tenant-level encryption-key separation |
| Termination assistance | Vendor commits to data return / deletion within named timelines; format of returned data; deletion certification |

The DPA is reviewed annually and updated when sub-processors change, when regulatory guidance updates, or when material service changes occur.

## Sub-processor authorization — the operational mechanics

Article 28(2) requires the controller to authorize the use of sub-processors. Article 28(4) requires the processor to impose the same data-protection obligations on sub-processors as it has accepted from the controller. The institution operates this in practice:

1. **Initial sub-processor list.** At DPA signing, the vendor's current sub-processor list is attached as a schedule. The institution authorizes the listed sub-processors in signing the DPA.

2. **Change notification.** When the vendor proposes a new sub-processor, it notifies the institution at least 30 days in advance. The notification names the sub-processor, the processing scope, the location, and the security measures the sub-processor commits to.

3. **Institution review.** The institution's privacy team reviews the proposed sub-processor against:
   - Service necessity (is the sub-processor needed for the chain service?)
   - Data location (does the location create cross-border transfer concerns?)
   - Security measures (does the sub-processor's commitment match the institution's standards?)
   - Reputation (any documented privacy or security incidents?)

4. **Decision.** The institution accepts or objects within the 30-day window. Acceptance is automatic if no objection is raised; the institution typically responds explicitly so the DPA's authorization log is complete.

5. **Objection handling.** If the institution objects, the vendor either finds an alternative sub-processor or accepts that the institution's tenant will not use the new sub-processor's services. A material objection that the vendor cannot accommodate may trigger early termination of the service agreement.

The institution maintains a log of sub-processor authorizations: each approval, each rejection, the rationale. The log is part of the institution's vendor-management records and is available to the DPO and to examiners.

## SOC 2 Type II as compliance evidence

The vendor's annual SOC 2 Type II report is the institution's primary evidence of vendor compliance with the DPA. The report's scope must cover:

- Trust Services Criteria for Security, Availability, Confidentiality, and Privacy
- Specific control objectives mapped to Article 32 security requirements
- The vendor's chain-hosting platform (not a parent company's general infrastructure)
- The reporting period (annually)

The institution's vendor-management process reviews the SOC 2 report on receipt:

| Review element | Action |
|---|---|
| Audit period | Confirm the report covers the most recent 12 months |
| Auditor identity | Confirm auditor is independent and reputable |
| Scope | Confirm chain-hosting platform is in scope |
| Controls operating effectively | Review the auditor's opinion; any qualifications addressed |
| Control deviations | Review any noted deviations; vendor remediation plan obtained |
| Sub-processors | Confirm sub-processors are within the SOC 2 scope or have their own attestations |
| Privacy-specific controls | Confirm privacy-store custody, key management, access controls are evidenced |

A SOC 2 Type II report with material qualifications or scope gaps triggers institution-side remediation: vendor commits to remediation plan; institution may impose additional controls (additional audit, additional monitoring, contract penalty) until the next clean report.

## When healthcare data is in scope — BAA-equivalent coverage

For institutions processing healthcare-related lending or healthcare-financing decisions, HIPAA imposes additional requirements. The vendor becomes a Business Associate under HIPAA; a Business Associate Agreement (BAA) is required.

The BAA covers HIPAA-specific obligations: minimum necessary scope, breach notification timelines (60-day under HIPAA Breach Notification Rule), administrative safeguards under the Security Rule, audit rights. The BAA is typically a separate document from the DPA, but they reference each other and align on overlapping obligations.

The vendor's SOC 2 Type II report scope includes HIPAA-specific controls when healthcare data is in scope. The institution may also require the vendor's HITRUST certification as additional assurance.

## Determination flowchart

```
Is the vendor making decisions about purposes? → Yes → Joint controller (Article 26)
                                                ↓ No
Is the vendor making decisions about means      → Yes → Joint controller (Article 26)
beyond service-delivery configuration?            ↓ No
                                                Processor (Article 28) — recommended
```

In practice, vendors offering a chain-hosting platform that the institution configures fall into the processor box. Vendors offering cross-tenant analytic services that they design and operate fall into the joint-controller box for those services. The institution may have both relationships with the same vendor — processor for the chain platform, joint controller for an analytic service — and the agreements name each scope distinctly.

## Worked example — Vidimus serves Acme Bank EU

Acme Bank Europe GmbH contracts with Vidimus Inc. (US-headquartered, EU subsidiary) for chain-hosting. The structure:

- **Roles.** Acme Bank Europe is sole controller. Vidimus EU is the processor.
- **DPA.** Executed between Acme Bank Europe and Vidimus EU. Includes SCCs Module Two for the US-headquartered parent's incidental access (e.g., support escalation).
- **Sub-processors.** Vidimus's sub-processor list includes AWS EU (Frankfurt region for Acme's tenant), an HSM service in EU, and Vidimus's US support function for incident escalation. Acme authorizes all three at DPA signing.
- **Data location.** Acme's tenant data is stored exclusively in AWS Frankfurt; Vidimus EU operates the platform; US parent has no production access.
- **SOC 2.** Vidimus provides annual SOC 2 Type II covering the Frankfurt platform. Acme reviews annually.
- **Schrems II safeguards.** EU-held encryption keys; US parent cannot decrypt EU-tenant data.

Acme's RoPA entry (`gdpr-ropa-template.md`) names Vidimus EU as processor and references the DPA. Acme's DPIA documents the cross-border-transfer risk and the supplementary safeguards.

## Cross-references

- `gdpr-ropa-template.md` — the controller / processor designation appears in the RoPA
- `gdpr-lawful-basis.md` — the controller is the party that establishes the basis
- `gdpr-dpia-template.md` — the DPIA is the controller's responsibility, with processor assistance
- `gdpr-dpo-consultation.md` — the DPO reviews vendor-relationship structure at deployment
- `vendor-hosted-controls.md` — the broader vendor-controls framework
- `vendor-conformance-attestation.md` — the conformance attestation a vendor provides
- spec §10 — security measures the DPA security schedule references
