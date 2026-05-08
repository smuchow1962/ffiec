# GDPR Article 15 — Data Subject Access Request fulfillment

> **What this doc is.** The institution's workflow for responding to a Data Subject Access Request (DSAR) under GDPR Article 15 when chain-captured AI decisions are in scope. The chain is indexed by (tenant_id, run_id, seq); customers do not know those identifiers. The institution's customer-correlation index (CUI) maps customer_id to the affected (tenant_id, run_id) pairs, and the privacy-store extract is the response artifact. Chain entries themselves are institution-internal evidence and supplement the response only when the institution's customer-communication policy calls for transparency. This doc names the workflow, the timeline, and the redaction rules.

## What Article 15 grants

GDPR Article 15(1) grants data subjects the right to obtain from the controller confirmation of whether personal data concerning them is being processed and, if so, access to the personal data. Article 15(3) requires the controller to provide a copy of the data being processed. Article 15(4) constrains the copy in cases where it would adversely affect the rights and freedoms of others.

The chain processing affects the data subject; Article 15 applies. The data subject is entitled to see what personal data the institution holds about them and what processing is occurring. The institution must produce the response within 30 days of receipt (Article 12(3)), extensible by 60 days for complex requests with notification within the first 30.

What the data subject is *not* entitled to is the institution's full audit-trail evidence. Article 15 is about the personal data, not about every internal artifact that mentions the data. The chain is institution-internal evidence — the integrity-bearing audit log that supports examination, dispute defense, and litigation. The personal data within the chain is the data subject's; the full audit-trail context is not. The workflow below honors that distinction.

## The workflow

```
1. Receive    → DSAR intake; identity verified
2. Map        → CUI lookup: customer_id → list of (tenant_id, run_id) pairs
3. Extract    → privacy-store records as primary response
4. Supplement → chain entries when customer-comm policy or specific request requires
5. Redact    → remove other data subjects' data from supplementary content
6. Deliver    → response within 30 days (extendable to 90)
7. Log       → record fulfillment in privacy log
```

## Step 1 — Receive and verify

The customer (or authorized representative) submits the request through the institution's standard DSAR intake. The privacy team:

- Logs the request with date, requester identity, scope ("everything," or specific accounts / decisions / time periods).
- Verifies the requester's identity using the institution's standard customer-authentication procedures. Identity verification is a precondition; producing personal data to the wrong person is itself a breach.
- Verifies authorization if a representative submitted the request (power of attorney, court order, parental authority for a minor).
- Acknowledges receipt within a reasonable time (typically within five business days).

The 30-day Article 12(3) clock starts at receipt of a verified request.

## Step 2 — Customer-correlation index lookup

The chain is indexed by (tenant_id, run_id, seq). The customer's identifier — visible to the customer — is the customer_id used in the institution's account systems, applications, or interaction records. To find the chain entries that touch this customer, the institution operates a customer-correlation index:

| CUI structure | Content |
|---|---|
| Lookup key | customer_id (the institution's customer master identifier, possibly augmented with account_number, application_id) |
| Mapping | List of (tenant_id, run_id) pairs where the customer's identity was tokenized into chain entries |
| Maintenance | Updated at chain-capture time: when a tokenized customer identifier enters the chain, the CUI records the (tenant_id, run_id) pair against the customer's master identifier |
| Custody | Co-located with the privacy-store; under the same access controls as the privacy-store; encrypted at rest |
| Retention | Same as the privacy-store mapping (7-year retention or until erasure under `gdpr-article-17-procedures.md`) |

The CUI is what makes DSAR fulfillment operationally feasible. Without it, the institution would need to scan every chain entry across every run looking for the customer's tokens — a query that does not scale. With the CUI, lookup is direct: customer_id → (tenant_id, run_id) list.

The CUI lookup is itself a privacy operation; access is logged per `gdpr-article-25-demonstrability.md` Artifact 4.

## Step 3 — Extract privacy-store records as primary response

The privacy-store is the source of truth for the customer's personal data. The DSAR response's primary artifact is the privacy-store extract: every record the institution holds about this customer in the privacy-store, presented in human-readable form.

The extract carries:

| Element | Content |
|---|---|
| Customer identity records | Name, address, contact information, account numbers — the personal data the privacy-store holds |
| Decision-relevant attributes | Decision factors the institution captured for the AI decisions: income, employment, credit history, behavioral signals — to the extent these are in the privacy-store rather than the chain |
| Decision outputs | The decisions made (approve/deny, pricing, risk score) tied to the customer |
| Source dates | When each piece of data was captured |
| Processing purposes | What the institution did or is doing with each piece of data |
| Recipients | Who has received the data (internal recipients per the RoPA entry; external recipients on a case-by-case basis) |
| Retention | The retention period for each category |

The extract is generated from the privacy-store by the privacy team's standard tooling. The format is the institution's standard DSAR response format — typically a PDF or a structured JSON depending on the customer's preference.

## Step 4 — Supplement with chain entries (conditional)

The chain entries are institution-internal evidence. They are not the primary response. The institution supplements the response with chain content only when:

- The customer's request specifically asks for the AI-decision audit trail ("show me what the AI saw and decided").
- The institution's customer-communication policy commits to transparency about AI decisions, and the AI-decision summary is part of the institution's standard DSAR response.
- The institution's legal team determines that producing the chain entries supports an articulated transparency or dispute-defense interest.

When supplementation is appropriate, the institution extracts the chain entries for the (tenant_id, run_id) pairs from the CUI lookup. The extract contains:

- Chain entries' decision-summary fields (what the AI decided, on what factors, when)
- Audit-context attributes (timestamp, run identifier — the institution may rename these to friendlier names like "decision date" and "internal reference")
- The chain's integrity-attestation note ("this record is part of an integrity-bearing audit trail per FFIEC banking standards")

The chain entries presented are *de-tokenized* for fields where the tokenized data is the customer's own (the customer is entitled to see their own personal data). For fields where the data is another party's (e.g., free-text reasoning that mentions another applicant in a joint application), the redaction in step 5 applies.

## Step 5 — Redact other data subjects' data

Some chain entries reference multiple parties. A joint-application decision references two applicants. A guarantor decision references the applicant and the guarantor. A model's reasoning chain may refer to comparable accounts or peers. These references contain other data subjects' personal data; producing them in this customer's DSAR response would breach the other parties' privacy.

The redaction rules:

| Content | Treatment |
|---|---|
| The requesting customer's own personal data | Produced in clear (the customer is entitled to see their own data) |
| Another natural person's identifying data | Redacted (with placeholder marker showing redaction occurred) |
| Aggregate references that don't identify specific others ("compared against 100,000 similar applications") | Produced in clear |
| Internal staff identifiers (the loan officer who reviewed) | Redacted unless the institution's transparency policy specifically allows disclosure |
| Vendor-internal references | Redacted unless they are necessary for the customer's understanding |
| Free-text reasoning mentioning specific parties | Redacted at the mention level; surrounding context preserved |

Article 15(4) explicitly permits redaction "where it would adversely affect the rights and freedoms of others." The redaction must be proportionate — redact only what is necessary to protect the others' rights, not the entire reasoning passage.

The redaction process is documented: who reviewed, what was redacted, on what basis. The documentation supports the institution's response if the customer challenges the redaction extent.

## Step 6 — Deliver

The institution delivers the response to the customer within 30 days of receipt. For complex requests (multiple runs, multiple tenants, extensive chain supplementation, complex redaction), the institution may extend by 60 days per Article 12(3); the customer is notified within the first 30 days that the extension is being used and why.

Delivery method matches the institution's standard customer-communication channels: secure web portal, registered mail, encrypted email. The customer receives the response and acknowledgment instructions; the institution retains the delivery record.

## Step 7 — Privacy log

Every DSAR fulfillment is logged in the institution's privacy log. The log entry carries:

| Field | Content |
|---|---|
| Request ID | Institution-internal request identifier |
| Date received | When the request arrived |
| Date verified | When the customer's identity was confirmed |
| Date completed | When the response was delivered |
| Customer token | Customer's ID token (NOT the original ID — the privacy log itself does not reintroduce PII) |
| (tenant_id, run_id) pairs | Pairs touched by the response |
| Response scope | Privacy-store only / privacy-store + chain supplement |
| Redaction summary | Categories of redaction applied (or "none") |
| Extension used | Yes / No; reason if yes |
| Officer | Identity of the privacy officer who completed the response |

## Worked example — typical DSAR response

A customer submits "send me everything you have about me" on 2026-05-07. The privacy team:

- **Step 1.** Verifies identity via the institution's standard authentication (knowledge-based questions plus photo ID upload). Logs receipt.
- **Step 2.** CUI lookup returns three (tenant_id, run_id) pairs: a 2024 loan application, a 2025 credit-line increase, a 2026 fraud-flag review.
- **Step 3.** Privacy-store extract produces 14 pages: identity records, contact data, the loan application data, the credit-line review data, the fraud-flag investigation outcome.
- **Step 4.** Customer's request was "everything," and the institution's policy includes AI-decision summary in DSAR responses. Chain entries for the three runs are extracted; AI-decision summary is added (4 additional pages).
- **Step 5.** The fraud-flag run's chain entry contains references to comparable suspicious accounts. Those references are redacted. The customer's own data in the entry is preserved.
- **Step 6.** Response delivered via the institution's secure customer portal on 2026-05-22 (15 days after receipt).
- **Step 7.** Privacy log records the request and the response.

Total response: 18 pages, delivered within 30 days, with documented redaction. The customer can review their data and exercise their other rights (rectification, restriction, objection, erasure outside retention) based on what they receive.

## Worked example — request requiring extension

A customer submits a DSAR involving 12 years of relationship, multiple AI products, and a complex investigative history with several runs flagged as fraud-related. The CUI lookup returns 47 (tenant_id, run_id) pairs across four tenants.

The privacy team determines the request is "complex" within the meaning of Article 12(3): the volume and the redaction work exceed the 30-day window. On day 15, the institution notifies the customer:

> Your access request received on 2026-05-07 is complex due to the volume of records and the multi-tenant scope. We are using the 60-day extension permitted by GDPR Article 12(3). Your response will be delivered no later than 2026-08-05. We will notify you if it is ready earlier.

The institution proceeds with extraction and redaction. The response is delivered on 2026-07-12 (66 days after receipt), within the extended deadline.

## What a DSAR response does NOT include

| Excluded content | Reason |
|---|---|
| The institution's full chain — all entries from all runs ever | Not personal data of the requester; institution-internal audit trail |
| Other data subjects' chain entries | Article 15(4) — would adversely affect rights of others |
| Internal MRM analysis or model-development artifacts | Not personal data of the requester; trade secret if applicable |
| Predictions or risk scores about the customer's future behavior derived after the period the customer asked about | Subject to separate evaluation per Article 22 if the customer asks about automated decision-making specifically |
| Litigation work-product | Privileged; produced only under formal litigation discovery, not via DSAR |
| Other customers' contact information mentioned in shared documents | Article 15(4) |

The response is comprehensive about the customer's personal data and bounded against scope creep into institution-internal materials.

## Special cases

### When the request comes during active litigation

If the customer (or their counsel) submits a DSAR while litigation against the institution is active, the institution's legal team is consulted before response. The DSAR is honored for the personal-data scope; chain entries that would be produced under litigation discovery anyway may be excluded from the DSAR response and produced under the litigation framework instead. The customer is notified of the dual track.

### When the request involves a deceased customer's data

GDPR (and most member-state law) does not extend Article 15 rights to deceased persons. A relative submitting a request on behalf of a deceased customer is evaluated under the institution's relevant law (probate authority, estate administration). The DSAR workflow does not automatically apply; the legal team handles.

### When the request involves a child's data

A child's request, or a parent's request on behalf of a child, is verified for parental authority and processed with heightened protections. The institution's children's-data procedures apply (referenced in `gdpr-dpia-template.md` Section 3 if children's data is in scope).

### When the request asks for chain-integrity verification

A customer asking "prove this AI decision wasn't tampered with" is asking about chain integrity, not Article 15 access. The institution responds via the customer-side verification path documented in `customer-dispute-procedures.md` ("How do I know the bank didn't tamper with this record?"). The DSAR workflow continues for the personal-data scope; the verification path is handled separately and may produce the institution's verifier output as additional evidence.

## Cross-references

- `gdpr-article-17-procedures.md` — companion procedure for erasure
- `gdpr-article-16-rectification.md` — companion procedure for correction
- `customer-dispute-procedures.md` — customer-side verification path for integrity questions
- `gdpr-article-25-demonstrability.md` — privacy-store access logging that supports CUI lookups
- `gdpr-ropa-template.md` — RoPA entry that names DSAR fulfillment in the data-subject-rights field
- `privacy-by-design.md` — tokenization design that the privacy-store extract reverses for the requesting customer's own data
- spec §10 — chain-entry structure that the CUI maps customer_id back to
