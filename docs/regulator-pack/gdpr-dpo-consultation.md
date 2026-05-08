# GDPR Article 38 — DPO consultation procedure

> **What this doc is.** The institution's procedure for consulting the Data Protection Officer (DPO) before deploying the chain on a new AI system or expanding it to a new data category. Article 38 requires the DPO to be involved properly and in a timely manner in all issues relating to the protection of personal data. For chain deployment, that means a defined consultation trigger, a defined review scope, a written DPO opinion within a defined SLA, and a documented override path when the institution decides to proceed despite DPO concerns. This doc names all four.

## Why DPO consultation is a procedure, not an ad-hoc conversation

Article 38(1) requires the institution to involve the DPO "properly and in a timely manner in all issues which relate to the protection of personal data." Article 39(1)(c) names the DPO's specific tasks, including providing advice on DPIAs and monitoring their performance. Article 36(4) caps the DPO's authority: the DPO's opinion is not binding, but if the institution overrides it, the override must be documented in writing.

Without a defined procedure, "consulting the DPO" can mean anything from a hallway conversation to a formal written opinion. An examiner reviewing an institution's privacy program will ask: at what points was the DPO consulted on the chain deployment, what did they say, and what did the institution do with the input? The answer needs to be specific and documentable.

The procedure below names: when the consultation is triggered, what the DPO reviews, how long the DPO has to respond, what the response looks like, and how an override is handled if needed.

## Consultation triggers

```
Trigger 1 — Initial deployment       → before the chain captures its first event for a new AI system
Trigger 2 — New data category        → before a new category enters the chain (e.g., health data added)
Trigger 3 — Cross-border change      → before storage / processing topology crosses a new jurisdiction
Trigger 4 — Material configuration   → before a privacy-by-design configuration change with material impact
Trigger 5 — Annual review            → annually, on the DPIA anniversary
Trigger 6 — Incident-driven          → after a chain or privacy-store breach, before resuming operations
Trigger 7 — Article 36 prior-consult → when residual risk is high despite mitigation
```

Each trigger is operationalized: a project team or operations team identifies the situation matches a trigger and submits a consultation package. The institution's privacy program documents which triggers fired in which year and which DPO opinions were issued.

### Trigger 1 — Initial deployment

Before the chain captures its first event for a new AI system, the project team submits a consultation package containing:

- The AI program description (purpose, data subjects, decision types)
- The proposed privacy-by-design configuration (which fields are tokenized, redacted, omitted)
- The DPIA (per `gdpr-dpia-template.md`)
- The lawful-basis analysis (per `gdpr-lawful-basis.md`)
- The RoPA entry draft (per `gdpr-ropa-template.md`)
- The cross-border-transfer assessment (if applicable)
- The vendor relationship structure (per `gdpr-controller-vs-processor.md`, if vendor-hosted)

The DPO reviews and issues a written opinion within 20 business days.

### Trigger 2 — New data category

When the institution adds a new data category to the chain (e.g., health data is captured for the first time, or behavioral inferences are expanded to a new behavior class), the project team submits an updated DPIA section, an updated configuration, and the lawful-basis update. The DPO reviews per the new-category-specific risk and the Article 9 implications if applicable.

### Trigger 3 — Cross-border change

If the institution moves storage or processing across jurisdictions (e.g., EU-tenant data is replicated to a US backup location, or a US institution adds an EU subsidiary), the cross-border posture changes. The DPO reviews the updated transfer assessment, Schrems II posture, supplementary safeguards, and SCC currency.

### Trigger 4 — Material configuration change

Configuration changes that meaningfully affect privacy posture trigger consultation. Examples: expanding tokenization scope, adding regex patterns that cover a new sensitive-data class, removing redaction from a field that previously carried it, extending retention beyond the FFIEC mandate. Cosmetic changes (renaming a field, updating non-PII metadata) do not trigger consultation.

The institution's change-management process flags configuration changes for DPO review based on a documented impact-classification rubric. Material changes route to the DPO; non-material changes proceed via standard change management with the DPO informed via the routine reporting cycle.

### Trigger 5 — Annual review

Annually, the DPO reviews the institution's chain-related privacy program: DPIA, RoPA entry, lawful-basis analysis, demonstrability artifacts, operational-event review, incident summary. The DPO issues an annual opinion on the program's continued adequacy.

### Trigger 6 — Incident-driven

After a chain or privacy-store breach, before operations resume, the DPO reviews the incident analysis, the remediation, and the post-incident control updates. The DPO confirms or conditions the resumption.

### Trigger 7 — Article 36 prior-consultation

When the DPIA's residual-risk analysis (`gdpr-dpia-template.md` Section 5) identifies risks that remain high despite mitigation, Article 36(1) requires prior consultation with the supervisory authority. The DPO is the institution's interface for this consultation. The DPO's opinion in this case includes whether the residual risk warrants Article 36 escalation.

## DPO review scope

For each consultation trigger, the DPO reviews the relevant elements. The full scope spans:

| Review element | What the DPO confirms |
|---|---|
| DPIA completeness (Article 35) | All six sections completed; risks named; mitigations specific; residual risk articulated; review cadence set |
| Privacy by design (Article 25) | Tokenization configured for all PII; stricter handling for special-category fields; demonstrability pack assembled |
| Lawful basis (Article 6 + 9) | Legal-obligation basis documented for retention; legitimate-interests balancing test on file; Article 9 exception named for special categories |
| Data-subject rights | Procedures in place for Article 15 access, Article 16 rectification, Article 17 erasure, Article 18 restriction, Article 20 portability, Article 21 objection, Article 22 review |
| International transfers | Transfer mechanism named (SCC, adequacy, BCR); TIA completed; supplementary safeguards documented |
| Vendor relationship | Controller / processor / joint-controller determination correct; DPA executed; sub-processor authorization current; SOC 2 Type II reviewed |
| Children's data | If in scope: parental consent procedure; heightened tokenization; shorter retention where lawful |
| HIPAA interactions | If healthcare data is in scope: BAA executed; HIPAA Security Rule administrative safeguards mapped; HIPAA Breach Notification timeline aligned |
| CCPA / state-law interactions | If California or other state-law residents are in scope: state-specific procedures available |
| Multi-jurisdictional conflict | Conflict-resolution posture documented; legal team's role named |
| Demonstrability | Five artifact types per `gdpr-article-25-demonstrability.md` are assembled |
| Operational readiness | Incident response procedures cover privacy-relevant scenarios; access logging operates; reconciliation runs |

The DPO does not redo each analysis; the project team submits the analysis and the DPO confirms the analysis is sound. Where the DPO finds gaps or weaknesses, the opinion names them and proposes specific remediation.

## DPO opinion format

The DPO's written opinion has one of three dispositions:

### Approved

The DPO concludes the deployment (or change) meets the institution's privacy program requirements and complies with applicable law. No conditions; the project team may proceed.

### Approved with conditions

The DPO concludes the deployment is sound subject to specified improvements. The opinion names each condition with a deadline (e.g., "implement quarterly privacy-store key rotation within 60 days," "update DPIA to address [specific risk] within 30 days," "restrict deployment to credit decisions only — do not expand to employment without re-consultation").

The project team may proceed and tracks the conditions to closure. Failure to close a condition by its deadline triggers re-consultation.

### Not approved

The DPO concludes the deployment should not proceed without remediation. The opinion names the gaps and the required remediation. The project team must address the gaps and resubmit before deployment.

A "not approved" opinion is rare in practice. Most issues surface during preparation and are addressed before formal consultation. When "not approved" is issued, it typically reflects a fundamental gap (no DPIA, no lawful basis named, vendor relationship undocumented, cross-border transfer mechanism absent).

## SLA — 20 business days

The DPO commits to issuing a written opinion within 20 business days of receiving a complete consultation package. The clock starts at the date the package is received and the DPO confirms it is complete; if the package is incomplete, the DPO returns it with specific gaps named, and the clock starts when the corrected package arrives.

For Trigger 6 (incident-driven), the SLA is shorter — typically 5 business days — to support timely operational resumption. The DPO commits to expedited review.

For Trigger 7 (Article 36 prior consultation), the SLA may extend if the DPO determines a supervisory-authority engagement is needed. The institution's project team is informed of the extended timeline and the rationale.

If the DPO cannot meet the SLA (vacation, workload, complex consultation), the DPO notifies the project team within the first 10 days and proposes a revised deadline. The institution's executive sponsor is copied so SLA pressure does not result in an inadequate opinion.

## Override procedure (Article 36(4))

Article 36(4) caps the DPO's authority: the opinion is advisory, not binding. If the institution decides to proceed against an "approved with conditions" opinion (without addressing the conditions) or a "not approved" opinion, the override must be documented in writing.

The override procedure:

1. **Override decision authority.** The institution's senior leadership (typically the Chief Risk Officer, Chief Privacy Officer, or a privacy committee that includes both) decides whether to override. Mid-level project teams cannot override; the decision rises to executive level.

2. **Override documentation.** A written decision is produced naming:
   - The DPO opinion being overridden
   - The specific conditions or non-approval grounds
   - The institution's rationale for proceeding (business urgency, alternative-mitigation, disagreement with risk assessment)
   - The compensating controls the institution will operate to address the DPO's concerns
   - The review cadence for re-evaluation
   - The signatories of the override decision

3. **Notification to the DPA (when applicable).** If the override touches Article 36 prior-consultation triggers, the institution notifies the supervisory authority of the override. The supervisory authority may intervene under its Article 58 powers.

4. **Audit trail.** The override decision is retained in the institution's privacy-program records, available to examiners, auditors, and the DPO for ongoing monitoring. The DPO's tasks under Article 39(1)(b) (monitoring compliance) include monitoring the override's compensating controls.

5. **DPO independence preserved.** Per Article 38(3), the DPO does not face penalty for issuing a "not approved" opinion that is overridden. The institution's HR and management practices recognize the DPO's independence; the DPO continues to perform their tasks regardless of override outcomes.

Overrides are rare and signal organizational tension. Most institutions never use the override path; the consultation procedure is designed so DPO concerns are addressed in the preparation phase before consultation.

## Worked example — initial deployment, approved with conditions

A bank's project team prepares to deploy the chain on a new fraud-detection AI. The consultation package is submitted on 2026-04-10. The DPO confirms the package is complete on 2026-04-12 and begins review.

The DPO issues a written opinion on 2026-05-05 (17 business days from package completion):

> **Opinion: Approved with conditions.**
>
> The deployment plan is sound. The DPIA, lawful-basis analysis, and configuration are complete. The vendor relationship is correctly structured as processor under Article 28. Three conditions apply:
>
> 1. The behavioral-signal field captured by the fraud-detection AI may include inferences that proxy for protected classes. The configuration must add tokenization to that field within 30 days of deployment. Currently the field is captured in clear.
>
> 2. The cross-border posture relies on the vendor's commitment to keep EU-tenant data in the Frankfurt region. The DPA security schedule must include this commitment as a contractual term. Currently it is in vendor's marketing material but not the DPA.
>
> 3. The annual review cadence in the DPIA Section 6 names the privacy team but not the DPO. The DPIA must explicitly include DPO sign-off in the annual-review record. Update the DPIA within 14 days.
>
> Subject to closure of these conditions, the deployment may proceed. I will follow up at the 30-day mark to confirm condition 1 is closed.

The project team:
- Updates the DPIA within 14 days (condition 3 closed)
- Negotiates the DPA addendum with the vendor and executes it within 21 days (condition 2 closed)
- Updates the configuration within 30 days; tokenization is applied to the behavioral-signal field; tests confirm the change (condition 1 closed)

All three conditions closed; the DPO confirms closure and the deployment proceeds with full DPO approval.

## Worked example — annual review, approved

A bank's existing chain deployment has been operating for 18 months. The annual review consultation package is submitted on 2026-08-15. The DPO reviews the DPIA refresh, the operational-event summary, the incident summary (none material), the demonstrability pack rehearsal results, and the privacy-program metrics.

The DPO issues a written opinion on 2026-09-05 (15 business days):

> **Opinion: Approved.**
>
> The chain deployment continues to meet the institution's privacy program requirements. The DPIA refresh accurately reflects the current state. No new risks identified. The privacy-program metrics show steady DSAR fulfillment within timeline, no privacy-store breaches, and successful key-fingerprint reconciliation throughout the period. The next annual review is due 2027-08-15.

The institution records the opinion and continues operations.

## DPO independence and reporting line

Article 38(3) requires the DPO to report directly to the highest level of management. In practice, the DPO reports to the Chief Privacy Officer (or in their absence, directly to the Chief Risk Officer or General Counsel). The DPO's compensation, performance evaluation, and tenure decisions are insulated from the project teams whose work the DPO reviews.

The DPO does not receive instructions on the substance of opinions. Project teams may dispute an opinion through factual clarification (presenting additional information that the DPO did not have); they cannot direct the opinion's outcome.

This independence is what gives DPO opinions their weight. An institution that pressures the DPO to approve undermines its own privacy program; the override path exists precisely so the institution can disagree with the DPO without compromising the DPO's independence.

## Cross-references

- `gdpr-dpia-template.md` — the DPIA the DPO reviews; Section 6 names DPO sign-off
- `gdpr-lawful-basis.md` — the lawful-basis analysis the DPO confirms
- `gdpr-article-25-demonstrability.md` — the demonstrability pack the DPO reviews
- `gdpr-controller-vs-processor.md` — the vendor relationship structure the DPO confirms
- `gdpr-ropa-template.md` — the RoPA entry the DPO sees as a draft during initial consultation
- `gdpr-article-17-procedures.md`, `gdpr-article-16-rectification.md`, `gdpr-dsar-fulfillment.md` — the data-subject-rights procedures the DPO confirms exist
- `incident-response-playbook.md` — the IR procedures the DPO references for Trigger 6
