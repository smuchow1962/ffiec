# GDPR Article 6 lawful-basis analysis

> **What this doc is.** The institution's documented lawful basis for chain processing under GDPR Article 6. Two bases combine: legal obligation (Article 6(1)(c)) for the FFIEC-mandated retention, and legitimate interests (Article 6(1)(f)) for fairness audit, model risk management oversight, and dispute defense. The two are not interchangeable — the institution names which basis covers which use, performs the legitimate-interests balancing test where it applies, and documents the result. This doc is what the institution puts in front of a Data Protection Authority that asks "on what lawful basis are you processing this data?"

## Why naming the basis matters

GDPR Article 6(1) requires at least one valid lawful basis for any processing activity. Without one, the processing is unlawful regardless of how good the technical controls are. The chain is a processing activity: it captures, stores, and analyzes personal data (customer identifiers, decision factors, model inputs, decision outputs). The institution must name the basis or the processing fails Article 6 on its face.

Naming the basis also drives downstream obligations. Different bases trigger different data-subject rights. Consent (6(1)(a)) carries withdrawal rights. Legitimate interests (6(1)(f)) carries Article 21 objection rights. Legal obligation (6(1)(c)) carries the strongest defense against erasure requests but the weakest discretion to expand processing scope. The institution's choice of basis frames how it must respond to data-subject requests for years afterward.

The chain processing rests on two bases. Each covers a distinct slice of the activity. The institution must understand which basis covers what, and document both.

## Primary basis — Article 6(1)(c) legal obligation

The FFIEC banking regulations require integrity-bearing custody and retention of AI decisions affecting credit, fraud determination, employment lending, and adverse-action notices. The relevant authorities:

| Authority | Coverage |
|---|---|
| 12 CFR 1002.12 (Regulation B / ECOA) | 25-month minimum retention of credit records, extended to seven years for examination and supervisory purposes |
| 12 CFR 1024 (Regulation X / RESPA) | Mortgage-servicing records retention; chain captures AI decisions in mortgage workflows |
| OCC examination standards | Multi-year retention for examination evidence, typically 7 years |
| Fed and FDIC supervisory standards | Parallel multi-year retention for supervised institutions |
| CFPB supervisory authority | Consumer-financial-product records for the duration of the supervisory relationship plus retention period |

These authorities, taken together, mandate retention of the records the chain captures. The chain is the institution's mechanism for satisfying that mandate. Article 6(1)(c) — processing "necessary for compliance with a legal obligation to which the controller is subject" — covers the chain's capture and retention activities under those regulations.

The legal-obligation basis is strong: it does not require a balancing test, it carries an Article 17(3)(b) refusal right against most erasure requests during the retention window, and it satisfies Article 14 transparency without requiring affirmative consent. The institution discloses the processing in its privacy notice; consent is not asked for, because none is required.

## Supplementary basis — Article 6(1)(f) legitimate interests

Some uses of the chain extend beyond the strict legal-obligation scope. Three in particular:

1. **Fairness audit.** The institution uses the chain to detect disparate-impact patterns in AI decisions, supporting ECOA fair-lending obligations and the institution's own equal-credit-opportunity program. This is enabled by the chain's comprehensive capture but is not strictly required by the retention mandate alone.

2. **Model-risk management (MRM).** The institution's MRM program (per SR 11-7 / OCC 2011-12) uses the chain to monitor AI-model behavior over time, detect drift, and inform challenger-model evaluation. This is regulatory-aligned but operates as the institution's own oversight, not as a direct legal-obligation output.

3. **Customer-dispute defense and litigation evidence.** When a customer disputes an AI decision (or when the institution faces ECOA, FCRA, or contract-based litigation), the chain provides tamper-evident proof of what the AI decided. This is the institution's interest in defending its decisions, not a discrete legal mandate.

These three uses sit on Article 6(1)(f) — processing "necessary for the purposes of the legitimate interests pursued by the controller... except where such interests are overridden by the interests or fundamental rights and freedoms of the data subject."

## The legitimate-interests balancing test

Article 6(1)(f) requires the institution to balance its legitimate interest against the data subject's rights. The institution documents the balance in three parts:

### Part 1 — Legitimate interest

The institution identifies the interest specifically. For each use:

- Fairness audit: detection of disparate-impact patterns, enabling early correction of biased models. Beneficial to data subjects collectively (other customers receive fairer decisions) and aligned with ECOA's fair-lending purpose.
- MRM oversight: detection of model drift and degradation, supporting safe and accurate AI decisions. Beneficial to data subjects whose decisions the model produces.
- Dispute defense: providing the institution the means to respond to disputes and litigation accurately. Beneficial to data subjects who dispute, in that the response is grounded in evidence rather than reconstruction.

### Part 2 — Necessity

The institution evaluates whether the interest could be satisfied with less data. For each use:

- Fairness audit: requires capture of decision factors and outcomes across the population. Aggregate-only logging would not enable per-decision audit; the chain's per-decision capture is necessary.
- MRM oversight: requires capture of model inputs, parameters, and outputs to compare actual against expected behavior. Reduced capture would prevent reproducibility of decisions for review.
- Dispute defense: requires capture of the specific decision and its inputs to support defense. Aggregate logging is insufficient for individual disputes.

For each use, the institution concludes the chain's capture scope is necessary. Less invasive alternatives (aggregate logging, reduced field set) would not satisfy the interest.

### Part 3 — Balancing against data-subject interests

The institution evaluates the data subject's interests and rights against the institution's interest:

| Data-subject interest | How the chain affects it | Mitigation |
|---|---|---|
| Privacy of personal information | Chain captures decision-relevant data, potentially including sensitive attributes | Tokenization (privacy-by-design.md) keeps personal data out of the chain itself; only tokens are captured |
| Control over data (Articles 15-22) | Data subjects can access (Article 15), rectify (Article 16), erase (Article 17), object (Article 21) | Procedures for each: `gdpr-dsar-fulfillment.md`, `gdpr-article-16-rectification.md`, `gdpr-article-17-procedures.md`, plus Article 21 review per the privacy team |
| Reasonable expectations | Data subjects expect AI decisions to be made and recorded, but may not expect 7-year retention | Disclosed in the privacy notice per Articles 13-14 |
| Risk of harm from processing | Tamper-evident audit trail does not create new exposure beyond what the underlying AI decision already creates | Tokenization plus privacy-store custody (separate physical and logical custody) limits exposure |

The institution concludes that with tokenization and the data-subject rights procedures in place, the data subject's interests are not overridden by the chain's processing. The legitimate-interests basis is valid for the supplementary uses.

The balancing-test record is itself documentation. It is retained, reviewed annually by the DPO (per `gdpr-dpo-consultation.md`), and updated whenever the chain's capture scope changes or the privacy-by-design configuration is revised.

## Why consent is not the basis

The institution does not request consent for chain processing. Three reasons:

1. **Consent must be freely given (Article 7).** A customer applying for credit cannot meaningfully consent to a processing activity that is required by law for the institution to make the decision. The "consent" would be coerced — the customer cannot decline and still receive the decision.

2. **The legal-obligation basis already covers retention.** Asking for consent in addition would create the impression that the customer can withdraw consent and trigger erasure. They cannot, because the retention is mandatory. Asking creates a false expectation.

3. **Consent withdrawal would conflict with Article 17(3)(b).** If consent were the basis, withdrawal would force erasure. The legal-obligation defense would be weakened by the institution's earlier claim that consent was the basis. The institution would be in the awkward position of arguing two contradictory things.

Consent is the wrong basis for this processing. The institution discloses the processing under Articles 13-14 (transparency) but does not solicit consent.

## Article 13/14 transparency without consent

The institution's privacy notice describes the chain processing in plain language:

> When an AI system makes a decision affecting your account or application, we record the decision and its supporting data in an integrity-bearing audit trail. We do this to comply with banking-regulation retention requirements, to detect and correct potential bias in the AI system, to manage model risk, and to defend the decision if you or a regulator later disputes it. The audit trail uses tokenization to limit personal-data exposure, and we retain it for the period required by banking regulators (typically seven years). You have rights to access, correct, or — outside the regulatory retention period — request erasure of your personal data; see [link to data-subject-rights page]. The lawful basis for this processing is legal obligation (FFIEC banking regulations) and our legitimate interests in fairness audit, model oversight, and dispute defense.

This satisfies Article 14 disclosure obligations without asking for consent. The data subject is informed; the institution has the lawful basis named.

## Special-category data — Article 9 interaction

Some AI decisions touch special-category data (Article 9): health-related underwriting, behavioral inferences that imply sexual orientation or religious belief, biometric authentication. Article 9(1) prohibits special-category processing except under enumerated 9(2) exceptions.

When special-category data is in scope, the institution names a 9(2) exception in addition to the 6(1) basis:

| 9(2) basis | When it applies |
|---|---|
| 9(2)(b) — employment, social security, social protection law | Lending decisions touching protected classes where ECOA fair-lending creates a parallel social-protection requirement |
| 9(2)(g) — substantial public interest | Fraud detection systems where AML/CFT public interest is documented |
| 9(2)(h) — preventive medicine, healthcare | Healthcare-financing decisions touching health data |
| 9(2)(a) — explicit consent | Reserved for cases where no other 9(2) exception applies; rare in this chain context |

The institution's tokenization configuration (privacy-by-design.md SDK configuration) applies stricter handling to special-category fields. The DPIA (`gdpr-dpia-template.md`) names the special-category fields in scope and the 9(2) exception that justifies their processing.

## Documenting the basis in the RoPA

The Records of Processing Activities entry for the chain (see `gdpr-ropa-template.md`) names the lawful basis explicitly:

> **Legal basis.** Article 6(1)(c) legal obligation (FFIEC banking regulations 12 CFR 1002 and equivalent OCC/Fed/FDIC/CFPB retention authorities) for retention. Article 6(1)(f) legitimate interests (fairness audit, MRM oversight, dispute defense) for supplementary uses, with balancing test on file dated `[date]`. For special-category data per Article 9, the institution applies Article 9(2)(b) (lending under ECOA), 9(2)(g) (fraud / AML), or 9(2)(h) (healthcare financing) per the field-specific configuration documented in the DPIA.

This is the documented basis a Data Protection Authority will see in an audit, an examiner will see in a regulator-pack package, and a litigant's counsel will see in discovery. Naming it once, here, with the citations and the balancing-test reference, is what defends the processing across all three audiences.

## Annual review

The lawful-basis analysis is reviewed annually as part of the DPIA review cadence (`gdpr-dpia-template.md` Section 6). The review confirms:

- The legal-obligation citations are current (FFIEC and CFR references unchanged or correctly updated)
- The legitimate-interests balance still stands (no change in capture scope, tokenization, or risk profile that would shift the balance)
- The Article 9(2) exceptions in scope are still applicable
- The privacy notice text still reflects the actual processing

If any element shifts, the institution updates the analysis, refreshes the privacy notice, and notifies the DPO. Material shifts (new data categories, expanded retention, new use cases) trigger DPO consultation per `gdpr-dpo-consultation.md`.

## Cross-references

- `gdpr-article-17-procedures.md` — the 17(3)(b) refusal that the legal-obligation basis enables
- `gdpr-ropa-template.md` — where the basis is recorded as Article 30 documentation
- `gdpr-dpia-template.md` — where the balancing test and special-category analysis are documented
- `gdpr-dpo-consultation.md` — DPO review of the basis at deployment and annually
- `privacy-by-design.md` — the tokenization design that supports the balancing-test mitigation
- spec §10.9 — the seven-year IKM retention that aligns with the legal-obligation basis
