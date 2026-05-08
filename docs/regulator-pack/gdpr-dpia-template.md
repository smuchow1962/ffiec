# GDPR Article 35 Data Protection Impact Assessment — chain template

> **What this doc is.** A template for the Data Protection Impact Assessment (DPIA) the institution must complete before deploying the chain on AI decisions affecting credit, fraud, employment, healthcare, or other regulated activities. Article 35 requires a DPIA for processing likely to result in high risk; chain-of-custody capture of AI decisions on those subjects is high-risk by definition. Without a DPIA, the institution cannot demonstrate that it considered the privacy impact before deploying. The six sections below are the institution's worked answer to that consideration.

## When a DPIA is required

Article 35(1) requires a DPIA when processing is "likely to result in a high risk to the rights and freedoms of natural persons." Article 35(3) lists three trigger cases that always require a DPIA:

1. Systematic and extensive evaluation of personal aspects, including profiling, where decisions produce legal or similarly significant effects.
2. Processing on a large scale of special-category data (Article 9) or criminal-conviction data (Article 10).
3. Systematic monitoring of publicly accessible areas on a large scale.

The chain captures decisions from systems that, by design, do (1) — they evaluate personal aspects (creditworthiness, fraud risk, eligibility) and produce decisions with legal or similarly significant effect (loan approval/denial, account closure, pricing). Trigger (1) applies on its face. If the AI system also processes special-category data (healthcare-financing, behavioral inferences proxying for protected classes), trigger (2) applies as well.

A DPIA is required before deployment. Continued processing without one is an Article 35 violation. The DPIA is also reviewed annually and whenever material change occurs.

## How the template is used

The institution's privacy team completes the six sections in dialog with the chain-operations team and the AI program owner. The DPO reviews and provides written opinion (per `gdpr-dpo-consultation.md`). The completed DPIA is signed by the responsible executive (typically the Chief Risk Officer or Chief Information Security Officer) and stored in the institution's GRC tool. It is referenced from the RoPA entry (`gdpr-ropa-template.md`).

## The template

### Section 1 — Necessity & proportionality

| Element | Content |
|---|---|
| **Purpose** | Integrity-bearing record of AI-driven decisions to support FFIEC banking-regulation retention, fairness audit, model-risk management oversight, customer-dispute defense, and litigation evidence. |
| **Necessity** | FFIEC banking regulations (12 CFR 1002 ECOA, OCC/Fed/FDIC examination retention, CFPB supervisory) legally require retention of AI decision records. The chain is the institution's chosen mechanism for satisfying that legal obligation while also enabling fairness audit and dispute defense. |
| **Proportionality** | Tokenization (privacy-by-design.md) limits personal-data exposure to a small mapping set in the privacy-store. Retention is bounded to seven years, aligned with the regulator's mandate (`gdpr-lawful-basis.md`). Special-category fields receive stricter handling per the configuration. |
| **Less invasive alternatives considered** | (1) Non-integrity-bearing logs: rejected — does not satisfy the integrity requirement that makes records defensible against tampering claims. (2) Retention shorter than 7 years: rejected — conflicts with FFIEC mandate. (3) Reduced field capture (decision + timestamp only): rejected — does not support fairness audit, MRM oversight, or detailed dispute defense. (4) On-device-only logging without server-side custody: rejected — does not survive device loss / customer churn / device tampering. |
| **Why the chosen approach is proportional** | Tokenization keeps personal data out of the chain itself; only tokens are captured. The privacy-store mapping is in separate custody under separate access controls. The seven-year retention matches the regulatory necessity, not an institutional preference. The data-subject-rights procedures (DSAR, rectification, post-retention erasure) ensure the data subject retains meaningful control. |

### Section 2 — Data categories

| Category | Examples | Privacy-by-design treatment |
|---|---|---|
| Customer identifiers | customer_id, account_number, application_id | Tokenized to chain; original retained in privacy-store |
| Decision-relevant attributes | income, employment status, credit history, debt-to-income | Tokenized or omitted per configuration; aggregate-derived fields (FICO score) may be captured directly where they are themselves the AI input |
| Free-text model inputs | prompt content, supporting-document excerpts | Regex-redacted for inadvertent PII; tokenized fields replaced with tokens before capture |
| Free-text model outputs | reasoning chain, explanation text | Regex-redacted; tokenized at capture if customer-identifying terms appear |
| Decision outputs | approve/deny, risk score, pricing tier | Captured directly; not personal data per se but combined with customer ID becomes personal data |
| Model metadata | model version, prompt template ID, parameter set | Captured directly; not personal data |
| Audit context | timestamp, run_id, tenant_id, key_version | Captured directly; not personal data |
| Special-category data (when in scope) | health condition codes (healthcare-financing), behavioral signals proxying for protected classes | Tokenized with stricter scope; some fields omitted entirely; Article 9(2) basis named per `gdpr-lawful-basis.md` |

The institution's privacy-by-design.md SDK configuration is the authoritative field-by-field mapping. This DPIA section summarizes; the configuration is the source of truth.

### Section 3 — Risks

The DPIA names the privacy risks the institution has identified, the likelihood of each, the impact, and the mitigations in place.

| Risk | Likelihood | Impact | Mitigation | Residual likelihood / impact |
|---|---|---|---|---|
| **Privacy-store compromise** — token-to-PII mapping exposed | Low (separate physical custody, HSM-protected key, access-restricted, monitored) | High (customer identifiable; potentially reversible against external data) | Separate physical / logical custody, HSM-protected encryption key, access logging, quarterly key rotation, incident response per `incident-response-playbook.md` | Low / High — accepted on assumption of the controls remaining in force |
| **IKM compromise** — past chain entries forged | Low (HSM custody, FIPS 140-2 L3 tamper-detection, key rotation, reconciliation) | High (chain-integrity defense lost for affected period) | HSM-protected master keys (spec §10), incident-response procedures (Scenario 4), key-fingerprint reconciliation per spec §10.1 P-6 | Low / High — accepted on assumption of HSM physical-tamper-detection effectiveness |
| **Retention beyond necessity** — institution holds records past the 7-year mandate | Low (calendar-based deletion of privacy-store mappings, automated scripts, monitored) | Medium (continued processing without legal basis) | Calendar-based deletion automation, documented legal-hold procedures (`incident-response-playbook.md` litigation-hold), DPO annual review | Low / Medium — accepted with annual review |
| **Inadvertent PII in free-text fields** — regex-redaction misses an unusual format | Medium (regex coverage is broad but not exhaustive) | Low to Medium (one or a few records affected; remediable) | Regex-redaction at capture, post-capture sampling for review, configurable pattern updates per privacy-by-design.md | Medium / Low to Medium — accepted with quarterly review |
| **Data-subject access scope** — chain entries surface in DSAR responses creating unintended disclosure | Low (DSAR workflow uses privacy-store as primary; chain entries are institution-internal evidence) | Medium (release of internal MRM evidence) | DSAR workflow (`gdpr-dsar-fulfillment.md`) explicitly excludes chain entries from primary response; chain is supplementary only when customer-communication policy requires it | Low / Medium |
| **Cross-border transfer compromise** — chain or privacy-store data accessed under foreign legal process | Low to Medium (depends on storage topology and Schrems II posture) | High (third-country government access to EU personal data) | Standard Contractual Clauses, Transfer Impact Assessment, supplementary technical safeguards (EU-held encryption keys), storage-topology alignment with data-subject jurisdiction | Low / High — accepted only with the supplementary safeguards in place |
| **Children's data captured without heightened protections** | Low (not all institutions process children's data; opt-in feature where applicable) | High (if processed: COPPA / GDPR Article 8 violation risk) | Children's-data scope identified, parental consent verified, stricter tokenization, shorter retention where lawful | Low / High when in scope, otherwise n/a |

The institution evaluates each risk against the specific controls in place. Residual ratings are not aspirational; they reflect what the institution actually accepts.

### Section 4 — Rights & safeguards

The DPIA maps each relevant data-subject right to the institution's procedure for honoring it.

| Right | Procedure | Reference |
|---|---|---|
| Article 15 — access | Customer-correlation index plus privacy-store extract; chain entries excluded from primary response | `gdpr-dsar-fulfillment.md` |
| Article 16 — rectification | Investigation, fact-verification, correction-record chain entry (parent-linked), privacy-store update | `gdpr-article-16-rectification.md` |
| Article 17 — erasure | Privacy-store mapping deletion outside retention mandate; documented refusal during retention | `gdpr-article-17-procedures.md` |
| Article 18 — restriction | Privacy-store flag; downstream-use restriction; chain capture continues | Privacy team standard procedure |
| Article 20 — portability | Privacy-store extract in machine-readable format | Privacy team standard procedure |
| Article 21 — objection | Legitimate-interests balance review by privacy team; outcome documented | Privacy team standard procedure |
| Article 22 — automated decision-making | Human review per `customer-dispute-procedures.md`; chain provides the evidence the human reviewer reads | `customer-dispute-procedures.md` |
| Articles 13/14 — transparency | Privacy notice text covering chain processing | Institution's privacy notice |

Each procedure is documented separately. The DPIA does not duplicate the procedure content; it references the procedure and confirms it exists.

### Section 5 — Residual risk

The institution names the risks it accepts and the assumptions on which the acceptance rests.

| Accepted risk | Assumption |
|---|---|
| IKM compromise via simultaneous physical HSM breach | FIPS 140-2 Level 3 physical tamper-detection is effective against attackers without HSM-vendor cooperation; nation-state-with-vendor-cooperation is out-of-scope for this DPIA |
| Privacy-store compromise via insider threat with elevated access | Access controls, separation of duties, monitoring, and quarterly access review reduce probability to low; high-residual-impact accepted on the basis of the controls |
| Cross-border transfer access via foreign legal process | SCCs plus EU-held keys reduce US-government-access scope; institution accepts that absolute prevention is not possible and relies on the supplementary safeguards plus the institution's legal team's response readiness |
| Inadvertent PII in free-text fields beyond regex coverage | Quarterly sample review surfaces patterns the regex misses; configuration is updated to cover them; institution accepts that some inadvertent capture occurs between updates |
| Retention beyond necessity due to litigation hold | Litigation holds are documented, reviewed by counsel, and lifted when litigation closes; institution accepts that some records remain longer than the seven-year mandate during active litigation |

Acceptances are signed off by the responsible executive. New risks discovered in the annual review are added to the table; existing acceptances are reaffirmed or revised.

### Section 6 — Review cadence

The DPIA is reviewed:

- **Annually**, on the anniversary of the most recent completion. The DPO confirms the analysis is current; the responsible executive resigns the acceptance section.
- **On material change**, including: new data category captured, expanded tokenization scope, new AI use case enrolled in the chain, change in vendor or sub-processor, change in cross-border topology, change in retention period, change in regulatory-guidance authority that affects the necessity / proportionality analysis.
- **On incident**, including: any chain or privacy-store breach; any chain-detected verification failure that surfaces an unanticipated risk; any data-subject complaint that surfaces an unanticipated rights gap.
- **On DPO request**, when the DPO judges that residual risk has shifted.

The DPIA's date of last review and date of next scheduled review are recorded. The DPO's most recent written opinion (per `gdpr-dpo-consultation.md`) is attached or referenced.

## Worked example — DPIA for a healthcare-financing AI

A bank's healthcare-financing arm deploys an AI that prices payment plans for medical procedures based on patient credit history and clinical-attribute inputs. The DPIA's six sections are populated:

- **Section 1.** Purpose: chain-of-custody for healthcare-financing decisions. Necessity: ECOA + state insurance / healthcare-finance retention + HIPAA Security Rule audit-control requirement. Proportionality: tokenization plus stricter health-data tokenization; 7-year retention.
- **Section 2.** Data categories include health-related attributes (CPT codes, procedure type, attending-provider type). Special-category data is in scope; Article 9(2)(h) basis named.
- **Section 3.** Privacy-store compromise risk evaluated against the elevated impact of health-data exposure; mitigations include heightened access controls and HIPAA-aligned safeguards (per regulator-pack HIPAA-specific docs not yet drafted).
- **Section 4.** DSAR workflow includes HIPAA "right to amend" alongside Article 16 rectification.
- **Section 5.** Accepted residual risks include healthcare-specific risks (provider-data inferences leaking patient diagnosis); mitigation includes provider-name redaction.
- **Section 6.** Annual review by the institution's HIPAA Privacy Officer and the DPO jointly.

The example shows the template absorbing healthcare-specific facts. Other regulated-activity contexts (employment, insurance, education) populate similarly.

## What completing this DPIA does not do

The DPIA is the institution's analysis. It is not a regulatory pre-approval. The institution does not file the DPIA with the Data Protection Authority unless Article 36 (prior consultation) is triggered — typically when the residual risk remains high despite mitigation. The DPO's role is to flag Article 36 cases per `gdpr-dpo-consultation.md`.

The DPIA is also not a one-time deliverable. The annual review and material-change triggers keep it current. A DPIA dated 2024 is not the institution's DPIA in 2026 if material change has occurred since.

## Cross-references

- `gdpr-lawful-basis.md` — the Article 6 and 9 analysis the DPIA's Section 1 references
- `gdpr-ropa-template.md` — the RoPA entry that records the DPIA's existence
- `gdpr-dpo-consultation.md` — the DPO review and Article 36 trigger
- `gdpr-article-25-demonstrability.md` — the demonstrability artifacts that support the DPIA's claims
- `gdpr-dsar-fulfillment.md`, `gdpr-article-16-rectification.md`, `gdpr-article-17-procedures.md` — the data-subject-rights procedures referenced in Section 4
- `privacy-by-design.md` — the SDK configuration referenced in Section 2
- `incident-response-playbook.md` — the IR procedures referenced in Section 3 mitigations
- spec §10 — the technical security measures referenced in Sections 3 and 5
