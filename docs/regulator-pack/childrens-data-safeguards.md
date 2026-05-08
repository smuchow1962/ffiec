# Children's Data Safeguards — Heightened Protections

> **What this doc is.** Procedures and controls for chain deployments that capture data about minors — under-13 (US, COPPA) or under-16 (EU, GDPR Article 8). Written so the institution can deploy the chain on AI systems serving school-funded loan products, youth banking, parental-consent-gated services, or any AI flow that touches a minor's personal data.

> **Why this matters.** Children are a special category of data subject. Multiple regulations layer extra requirements: parental-consent traceability, stricter minimisation, behavioral-inference restrictions, opt-out rights specific to minors, and shorter retention questions. Operating the chain on children's data without these protections is a regulatory and reputational risk that compounds across regimes.

---

## 1. The applicable regimes

| Regulation | Citation | Scope | Threshold |
|---|---|---|---|
| **COPPA** (US) | 15 U.S.C. §§6501–6506; 16 CFR Part 312 | Online services directed to children; collection of personal information from children | Under 13 |
| **GDPR** (EU) | Articles 6(1)(a), 8, Recital 38 | All processing where consent is the legal basis for an information-society service offered directly to a child | Under 16 (member states may lower to 13) |
| **CPRA** (California) | §1798.120(c), §1798.140(z) opt-out for minors | Sale or sharing of personal information about a consumer the business has actual knowledge is under 16 | Under 16 (under 13 requires affirmative authorization from parent or guardian) |
| **State-level extensions** | e.g., Florida HB 3 (online protections for minors), Utah Social Media Regulation Act | Various; commonly affect online services and social media | Varies, typically under 18 |

The institution applies the **strictest applicable rule** to any data subject in scope.

---

## 2. When children's data appears in the chain

Common deployment scenarios:

| Scenario | Likely data flows |
|---|---|
| School-funded lending (AI assesses repayment-likelihood for student-loan products) | Student name and contact, parent/guardian name and financial data, school identifier, academic standing |
| Youth banking products (AI determines eligibility, fraud-flagging, fee structure) | Minor's account info, parent/guardian co-signer info, transaction patterns |
| Custodial accounts (AI manages investment recommendations) | Beneficiary minor's identity, custodian's identity, investment behavior |
| Family-plan financial products (AI cross-recommends across household) | Family-member identity including minors |
| Healthcare-financed lending for pediatric care | Minor patient identity, diagnosis codes, treatment history (PHI per `hipaa-minimum-necessary.md`) |

Each flow triggers the heightened-protection regime described below.

---

## 3. Parental consent traceability

### 3.1 The consent record

For services that require parental consent (COPPA verifiable parental consent; GDPR Article 8 holder-of-parental-responsibility consent), the consent itself is a discrete record that must be:

| Property | Requirement |
|---|---|
| Verifiable | The consent was obtained from a person with parental responsibility, verified using a method appropriate to the data sensitivity (signed form, government-ID match, video verification, micro-deposit verification, etc.) |
| Dated | Date of consent recorded |
| Scoped | The consent specifies which processing activities and data categories it covers |
| Revocable | The consent can be withdrawn; the institution operates the withdrawal as a privacy event |
| Retained | The consent record is retained for the entire processing duration plus appropriate post-processing retention |

The consent record lives in the **privacy-store**, not the chain ledger. The chain may carry a token referencing the consent record (e.g., `audit.consent.parental_consent_id` field) but does not duplicate the consent details.

### 3.2 Linking consent to chain entries

When the AI processes a minor's data and the chain captures the decision:

| Chain field | Content |
|---|---|
| `audit.subject.is_minor` | Boolean; true if the data subject is under the applicable age threshold |
| `audit.subject.parental_consent_token` | Token referencing the privacy-store consent record (HMAC of consent ID under privacy-store key) |
| `audit.subject.minor_age_band` | Coarse age band (e.g., "under_13", "13_to_15", "16_to_17") rather than exact age, to limit identifiability |

This binding lets the institution prove, on a per-decision basis, that a verifiable parental consent was on file at the time of processing.

### 3.3 Consent withdrawal

If a parent or guardian withdraws consent:

1. The institution stops processing prospectively.
2. The consent record is updated in the privacy-store with a `withdrawn_at` timestamp.
3. An `audit.consent.withdrawn` chain entry is appended to the chain, parent-linked to the consent's prior chain references.
4. Erasure of personal information proceeds per `ccpa-cpra-rights.md` Right to Delete and GDPR Article 17, subject to legal-obligation retention exceptions.
5. Future AI decisions for that subject are blocked or routed to a non-AI pathway.

---

## 4. Stricter data minimisation

### 4.1 The narrower scope rule

For minors, the institution narrows the SDK's capture scope below what it would capture for adults under the same AI system:

| Field type | Adult scope | Minor scope |
|---|---|---|
| Direct identifiers (name, account number) | Tokenized | Tokenized — identical |
| Contact (address, phone) | Tokenized | Tokenized — identical |
| Free-text reasoning | Tokenized or redacted | Tokenized AND length-limited; longer narratives are summarized into category codes |
| Behavioral inferences (creditworthiness, risk scores) | Captured | Captured only if directly supporting the documented purpose; speculative inferences not captured |
| Lifestyle / spending patterns | Captured (subject to Article 9 special-category analysis) | **Not captured** unless explicitly required by purpose and consent covers it |
| Social/relational data | Captured if relevant | **Not captured** for under-13; tokenized for 13–16 |
| Geolocation | Tokenized to coarse granularity | Tokenized to coarsest possible granularity (e.g., ZIP code level rather than street) |

The narrower scope is enforced by SDK configuration. Chain entries for minors carry an `audit.subject.is_minor` flag that the SDK honors at capture time.

### 4.2 Documenting the narrower scope

The institution maintains a written minimisation memo for each AI system processing minors. The memo names the documented purpose, the minimum data set required for the purpose, and the configuration that enforces it. The memo is reviewed annually and on any system change.

---

## 5. Behavioral-inference restrictions

Behavioral inference — using observed behavior to derive likely traits, preferences, or future actions — is more strictly limited for minors. The institution's defaults:

| Inference category | Adult posture | Minor posture |
|---|---|---|
| Creditworthiness inference (for lending) | Permitted under documented purpose | Permitted only with parental consent and only for the credit decision itself; not retained for cross-product use |
| Marketing / preference inference | Permitted under purpose limitation | **Prohibited** for under-13; restricted for 13–16 (no behavioral targeting under CPRA §1798.120(c) and forthcoming CPPA regulations) |
| Risk-scoring for fraud | Permitted | Permitted, with shorter score retention |
| Affinity grouping (clustering by behavior) | Permitted | **Prohibited** for under-13; restricted for 13–16 |
| Mood / emotion inference (from text or interaction) | Subject to Article 9 analysis | **Prohibited** outright for minors |
| Predictive lifecycle inference (life-event forecasting) | Subject to Article 9 analysis | **Prohibited** outright for minors |

When the AI system technically supports an inference category prohibited for minors, the SDK or upstream service blocks the inference for any subject flagged as a minor. The chain entry for a minor captures only the inferences permitted for that age band.

---

## 6. Shorter retention consideration

The default retention for chain entries is 7 years per `retention-justification.md`. For minors, the institution evaluates whether 7 years is necessary:

| Question | Decision rule |
|---|---|
| Does the FFIEC 7-year supervisory floor apply to this AI's decisions about a minor? | If yes (e.g., a credit decision under ECOA), retain 7 years |
| Does the minor reach majority within 7 years? | If yes, the institution considers whether retention beyond majority is necessary or defensible |
| Is there a state-law shorter-retention rule for children? | Some states (varies) require shorter retention for minors' data; honor the shorter rule |
| Does the parent / guardian / now-adult subject request earlier deletion? | Honor the request to the extent consistent with legal-obligation retention; documented exception applies until end of 7 years |

The institution's default for minor data: retain 7 years if a regulatory floor requires it; otherwise consider retention to age-of-majority plus 2 years (covering young-adult dispute window). The actual decision is recorded in the minimisation memo.

For data outside any regulatory retention floor (e.g., behavioral inference data, marketing-related data), retention is shorter — typically 1 to 2 years for minors, deleted promptly when no longer needed.

---

## 7. CPRA opt-out rights for minors

CPRA §1798.120(c) gives minors specific opt-out rights:

- **Under 16**: Sale or sharing of personal information requires opt-in. The default is no sale or sharing.
- **Under 13**: Opt-in must come from a parent or guardian (verifiable parental consent per §1798.120(c)(1)).
- **13 to 15**: Opt-in must come from the minor themselves (§1798.120(c)(2)).

The institution operates these defaults. The chain captures decisions; "sale" or "sharing" of chain-derived data is governed by these opt-in rules. The institution's privacy notice for minors and their parents discloses the opt-in defaults and the mechanism to grant or revoke consent.

For automated decision-making opt-out (CPRA §1798.185(a)(16) and forthcoming regulations), minors' opt-out rights track adult rights with the addition that, for under-13, the parent or guardian exercises the right.

---

## 8. Audit procedure for children's-data-specific protections

The institution's audit procedures (`audit-procedures.md` operates as the parent document) include a dedicated procedure for minor-data deployments. The procedure samples chain entries flagged with `audit.subject.is_minor = true` and confirms:

| Audit check | Expected evidence |
|---|---|
| Parental consent record exists for the data subject | Privacy-store consent record retrievable by `audit.subject.parental_consent_token` |
| Consent is verifiable per applicable rule | Record carries verification method consistent with COPPA/GDPR/state requirement |
| SDK configuration applies the minor-scope rules | Configuration file shows minor-specific tokenization and field-restriction settings; Chain entry fields conform |
| No prohibited behavioral inferences are present | No inference fields for the subject outside the permitted set |
| Retention timer is correct | Privacy-store record carries the applicable retention boundary (typically 7 years or age-of-majority+2, whichever is documented) |
| CPRA opt-in defaults are honored | No "sale" or "share" downstream flow recorded for the subject without explicit opt-in |
| Withdrawal events are honored | If consent was withdrawn, no AI decisions recorded after the withdrawal date |

The audit runs annually as a dedicated engagement; minor-data findings escalate to the privacy office and the AI Governance Committee.

---

## 9. Escalation and governance

Decisions involving minor data are subject to elevated governance:

| Decision | Approver |
|---|---|
| Initial deployment of an AI system that processes minor data | AI Governance Committee + DPO |
| New data-category capture for an existing minor-data AI | DPO + privacy office |
| Retention-policy variation from default (e.g., shortening or extending for a specific minor-data product) | DPO + legal team |
| Cross-jurisdiction transfer of minor data | DPO + legal team + privacy office (per `international-transfers.md`) |
| Incident affecting minor data | IR commander + DPO + AI Governance Committee chair |

The DPO Consultation Procedure (per the institution's separate DPO consultation document) is mandatory for any minor-data deployment regardless of risk-tier classification.

---

## 10. Cross-references

- COPPA 15 U.S.C. §§6501–6506; 16 CFR Part 312 — operationalized here.
- GDPR Article 8, Recital 38 — operationalized here.
- CPRA §1798.120(c), §1798.140(z) — operationalized here.
- `hipaa-minimum-necessary.md` — applies in addition where minor data is also PHI (pediatric care).
- `retention-justification.md` — default retention; this document narrows for minor data where appropriate.
- `pseudonymization-vs-anonymization.md` — tokenization analysis; minor data requires the same pseudonymization analysis with tighter scope.
- `ccpa-cpra-rights.md` — consumer rights; this document specifies minor-specific opt-in defaults.
- `international-transfers.md` — cross-border transfer of minor data.
- `breach-notification-matrix.md` — breach notification for minor-data incidents follows the matrix; HHS guidance for COPPA-covered breaches may add notification recipients.
- `dpia-template.md` — DPIA includes a minor-data section when applicable.
- `audit-procedures.md` — parent procedure document referenced in §8.

---

## 11. Review cadence

This document is reviewed annually by the DPO, privacy office, legal counsel, and AI Governance Committee. Triggers for early review:

- A new state law extends minor-data protections.
- The CPPA publishes final regulations on minor-data automated decision-making.
- A regulator (FTC for COPPA, DPA for GDPR) issues new guidance.
- The institution adds or removes a product line touching minor data.
- An incident reveals a control gap specific to minor data.
