# AI policy alignment

> **What this doc is.** How the chain aligns with U.S. and international AI policy frameworks. Closes the AI policy alignment partial.

## U.S. federal AI policy

### NIST AI Risk Management Framework (AI RMF)

The chain supports AI RMF Functions:

| Function | How chain supports |
|---|---|
| **Govern (GV)** | GV.OC (organizational context); GV.SC (cybersecurity supply chain): supply-chain doc + cosign + reproducible builds |
| **Map (MP)** | MP.5 (impacts of AI systems): chain captures decisions; downstream analysis maps to outcomes |
| **Measure (MS)** | MS.4 (validate AI system performance): chain captures the inputs and outputs that performance metrics analyze |
| **Manage (MG)** | MG.4.1 (validity and reliability) — chain provides integrity-bearing records; MG.4.3 (accountability and transparency) — chain output is the transparency artifact; MG.5.1 (incident response) — IR playbook |

### Executive Order 14110 (Safe, Secure, and Trustworthy AI)

The EO directs federal agencies to develop AI risk-management practices. For financial regulators, this manifests as expectations around AI governance. The chain supports the EO's themes:

- **Trust and verification** — the chain's integrity property supports trust claims about AI decisions
- **Algorithmic accountability** — the chain captures the AI's behavior with traceability
- **Transparency** — the chain's output is reviewable by regulators

### U.S. Treasury AI Risk Management Framework (Feb 2026)

The Treasury framework (referenced in `spec/chain-of-custody-v1.md` §10) addresses AI risk in financial services. Logging integrity is a specific expectation; the chain satisfies it. The framework's broader expectations align with SR 11-7 model-risk lifecycle controls; the chain composes with the institution's existing MRM program at four specific points:

| Treasury RMF expectation | Chain composition |
|---|---|
| **Model inventory.** The institution maintains an inventory of AI models in use. | The chain captures every decision; the resulting per-tenant per-day record set supports model-inventory completeness checks (every model that produced a decision is enumerated by the `gen_ai.request.model` field on the integrity-bound canonical bytes). |
| **Model lifecycle controls.** SR 11-7-style validation, ongoing monitoring, retirement evidence. | Chain output supports validation (re-run decisions under captured `gen_ai_parameters`), ongoing monitoring (drift analysis on chain-bearing data), and retirement evidence (the chain shows the model's final decisions with full provenance). |
| **Third-party model governance.** Vendor-supplied models in production. | Chain output is independent of vendor cooperation; the institution can produce integrity-bearing records of vendor-model decisions even after vendor change or vendor failure. |
| **Ongoing monitoring.** Drift, accuracy, bias evaluation. | Drift analysis runs on chain-bearing data with confidence the data is what the model actually produced. The integrity property removes a class of "did the vendor's logging system drop or alter events" doubts. |

The MRM director defending the institution's chain adoption in front of an OCC examiner cites these compositions explicitly; the Treasury framework's expectations are not invented in v1.0 but are recently consolidated, and the chain's evidence base directly supports them.

### CFPB AI guidance

CFPB's positions on AI in consumer-facing decisions emphasize:

- Adverse-action notice requirements when AI denies credit
- UDAAP analysis for AI-driven decisions
- Disparate-impact testing for fair-lending compliance

The chain provides the integrity-bearing record CFPB expects when reviewing AI decisions; the institution's broader compliance procedures handle the consumer-facing requirements (`docs/customer-dispute-procedures.md`).

### OCC emerging supervisory posture on AI

The OCC has been signalling — through speeches, supervisory letters to specific institutions, and joint statements with the Fed and FDIC — increasing focus on AI model-risk governance. As of v1.0-rework publication, the OCC has not issued formal new AI-specific guidance; instead, it operationalises AI model-risk supervision under the existing SR 11-7 / OCC Bulletin 2011-12 model-risk-management framework.

The chain is increasingly pertinent to OCC AI examinations under the existing framework — no new rule required. Specifically:

- **OCC Bulletin 2011-12 §III (Model Validation).** The chain's integrity-bearing records support the validator's effective challenge of AI-agent decisions. The validator has confidence the data being analysed is what the model actually produced.
- **OCC Bulletin 2011-12 §IV (Governance).** The chain's verifier output is institution-side evidence the model-governance committee can review (per `MRM-COMMITTEE-BRIEF.md`).
- **OCC supervisory letters on AI risk (2024+).** Specific institutions have received supervisory letters requesting documented AI model-risk frameworks; the chain provides one integrity-control component the institution cites.

For an MRM director defending chain adoption to the OCC, the framing is: "The chain is the technical control implementing the integrity-of-AI-decision-evidence requirement under the existing SR 11-7 / 2011-12 framework. We adopted the chain to position the institution for current and emerging OCC expectations, including expectations the OCC has signalled through supervisory letters to peer institutions." The framing avoids over-claiming the OCC has explicitly required the chain (it has not) while recognising the chain's relevance to the OCC's evolving supervisory posture.

A formal new OCC AI-supervisory issuance is a v1.x candidate trigger; if it lands, the spec working group reviews and updates this section accordingly.

## International AI policy

### EU AI Act

The EU AI Act (entered into force 2024; full applicability by 2026) imposes obligations on high-risk AI systems used in financial services:

| Article | Chain support |
|---|---|
| **Article 9 (Risk management system)** | Chain composes with the institution's broader risk-management framework |
| **Article 10 (Data governance)** | Chain captures data with integrity; institution's broader data-governance program manages |
| **Article 12 (Record-keeping)** | **The chain is the record-keeping mechanism for high-risk AI systems** |
| **Article 13 (Transparency)** | Chain output supports transparency reporting |
| **Article 14 (Human oversight)** | Chain captures AI decisions for human review with sufficient information for the human reviewer to second-guess the AI. The Act's Article 14 specifically requires *effective* human oversight: the chain's `gen_ai_parameters` (decoding parameters), `audit.*` (institutional payload including system prompt and retrieval context), and `gen_ai.request/response.model` (model identifier) are the substrate that makes Article 14 oversight effective. The institution configures `audit.*` schema to ensure the human reviewer has the input set the AI actually used. |
| **Article 15 (Accuracy, robustness, cybersecurity)** | Chain provides the cybersecurity layer for AI-decision records |
| **Article 26 (Obligations of deployers)** | For financial-services deployers of high-risk AI, Article 26 imposes specific obligations including monitoring (Art. 26(5)), logging (Art. 26(6) — logs retained at least 6 months unless longer required by Union or national law), and informing affected persons (Art. 26(11)). The chain composes: the per-event MAC + daily seal + HSM-signed root provide the integrity-bearing log for Art. 26(6); the operational events catalog (spec §10.2) provides the monitoring substrate for Art. 26(5); the chain's audit-trail availability supports Art. 26(11) institution-side communications to affected persons. |

The headline alignment is **Article 12** — record-keeping is exactly what the chain does for AI systems. Articles 14 and 26 are the secondary alignments where the chain is a load-bearing component of the deployer's compliance posture, not just a record-keeping artifact.

### EBA Guidelines on AI in Financial Services (2024+)

The EBA guidelines extend EU AI Act to banking specifics. Section 5 (logging and audit trail) requires institutions to maintain comprehensive logging of AI decisions. The chain satisfies this.

### DORA (Digital Operational Resilience Act)

DORA imposes ICT risk-management obligations. The chain composes at five articles:

- **Article 5 (ICT risk management)** — the chain is one ICT control within the institution's broader framework
- **Article 6 (ICT risk-management framework)** — DORA Art. 6 requires institutions to have a sound, comprehensive, well-documented ICT risk-management framework. The chain is one component of that framework; the institution's ICT risk function references the chain's threat model (`docs/design/09-threat-model.md`), the IR playbook, and the operational events catalog as the chain-specific contributions to the institution's documented framework. The chain is not the framework itself; it composes within it.
- **Article 8 (Identification of ICT-supported business functions)** — DORA Art. 8 requires institutions to identify ICT-supported business functions and the underlying ICT systems, including third-party providers. AI agents in production are ICT-supported business functions; the chain's tenant-and-deployment topology (per-tenant ledger, BYOC and vendor-hosted variants) is the integrity-bearing record of which ICT systems support which business functions. The institution's Art. 8 register references the chain's per-tenant scope.
- **Article 17 (Incident reporting)** — 24-hour reporting for major incidents; the chain's IR playbook accommodates DORA timing
- **Article 28 (ICT third-party risk)** — vendor SOC reports for chain implementations satisfy DORA's third-party-risk-management framework

The MRM director coordinating with the institution's ICT risk function uses the Art. 6 + Art. 8 + Art. 17 + Art. 28 composition to articulate the chain's role in the broader ICT framework, not as a standalone control.

### NIS2

NIS2 expands cybersecurity obligations across the EU. The chain supports:

- **Article 21 (Cybersecurity risk-management measures)** — chain is one risk-management measure
- **Article 23 (Reporting obligations)** — 24-hour early warning + 72-hour notification + 1-month final report; the IR playbook adapts

### UK Financial Conduct Authority (FCA) and Prudential Regulation Authority (PRA)

UK regulators have published guidance on AI in financial services. The chain supports:

- FCA's principles on AI use in financial markets
- PRA's expectations on operational resilience, including AI-related risks

### Singapore MAS Veritas / FEAT principles

Singapore's MAS has published the Veritas framework and FEAT (Fairness, Ethics, Accountability, Transparency) principles. The chain supports the Accountability and Transparency principles by providing integrity-bearing records.

### Japan Financial Services Agency (JFSA)

JFSA has published principles for AI use in financial services. The chain supports the related accountability and audit-trail expectations.

### Brazilian Central Bank (BCB)

BCB has emerging AI guidance. The chain supports the audit-trail expectations.

## State-level U.S. AI policy

State-level AI laws are emerging:

- **California (CCPA + emerging AI legislation)** — consumer rights regarding AI-driven decisions
- **Texas (TDPSA)** — data-privacy alignment
- **Virginia (VCDPA)** — consumer rights
- **Colorado (CPA)** — privacy plus AI-specific expectations
- **Connecticut (CTDPA)** — consumer rights
- **New York (proposed AI legislation)** — algorithmic accountability

State-specific compliance is the institution's compliance team's responsibility; the chain provides the substrate.

## Cross-jurisdictional consideration for institutions operating internationally

For institutions operating across multiple jurisdictions, the chain provides one integrity-bearing record consumed by multiple regulators in their respective frameworks. The institution maintains:

- Per-jurisdiction control descriptions mapping the chain to local frameworks
- Per-jurisdiction notification procedures aligned with local rules
- Per-jurisdiction privacy-store residency where required

The chain itself is jurisdiction-neutral. Multi-jurisdictional adaptation is the institution's compliance program's responsibility, with the chain providing the consistent substrate.

## Tracking AI policy evolution

AI policy is rapidly evolving. The institution's compliance program tracks:

- Federal regulator publications (NIST, Treasury, CFPB, FFIEC member agencies)
- State legislation
- International developments (EU, UK, Singapore, Japan, Brazil)
- Industry standards (ISO/IEC 42001 AI management system, etc.)

The chain's primitive set is policy-stable: integrity, tamper evidence, independent verifiability are durable concepts across policy evolution. The chain may need to add new attributes or extensions (see candidate attributes in `05-otlp-wire.md`) to align with new policy specifics; the spec evolves with the policy landscape.

## Adoption-decision considerations under emerging policy

Institutions deciding whether to adopt the chain consider:

- Current regulatory expectations (FFIEC II.C.10, EBA Section 5, etc.)
- Anticipated regulatory expectations (EU AI Act full applicability, emerging state laws)
- The institution's risk profile and AI use case
- The cost of adoption vs the cost of not adopting (audit defensibility, vendor independence, customer dispute handling)

For most institutions with material AI use, adoption is increasingly the default. The cost model articulates the operating cost; the policy landscape articulates the regulatory expectation.

## Multi-jurisdiction retention coordination

Institutions operating across multiple jurisdictions face overlapping retention rules. The chain's spec defaults and operator guide assume the typical U.S. FFIEC posture of seven-year retention of audit-trail records, which is the dominant retention floor for financial-services audit logs in U.S. supervision. The EU AI Act introduces a different floor for high-risk AI deployers under Article 26(6) — logs retained "at least 6 months" unless a longer Union or national-law retention requirement applies (Regulation (EU) 2024/1689, Article 26(6)). For an institution subject to both U.S. and EU supervision, the retention rules compose; for an institution subject to one or the other, the applicable rule is whichever the institution's supervisor invokes. The chain itself is jurisdiction-neutral — it produces the integrity-bearing record at any retention horizon — so the question is institution-side: how does the institution determine the operating retention period when multiple jurisdictional minimums interact?

The operating principle is that retention runs to the longer of the applicable minimums. The institution's control description names the rationale and the applicable jurisdictional minimums explicitly, so an examiner working any one jurisdiction sees both the institution's choice and the reasoning. Specifically the institution operates the longer of:

- (a) the U.S. FFIEC seven-year retention floor when the institution is supervised by an FFIEC member agency (OCC, Fed, FDIC, NCUA, CFPB);
- (b) the EU AI Act Article 26(6) six-month minimum when the institution is a deployer of a high-risk AI system in scope of the AI Act, with the actual retention extended to the longer period any applicable Union or national law mandates;
- (c) the institution's commercial retention policy where its broader information-governance program imposes a retention floor on AI-decision evidence;
- (d) any other applicable jurisdictional minimum — UK FCA / PRA records-retention guidance, Singapore MAS Veritas evidence-retention expectations, Japan JFSA audit-trail retention principles, Brazil BCB record-keeping expectations, or any state-level or local-supervisory retention requirement the institution has identified during its compliance scoping.

The institution's control description states the operating retention period as a single value, names each contributing minimum, and identifies which minimum is the binding floor. The named retention period is the basis of the chain's storage capacity planning, the cost-model retention horizon, and the institution's data-lifecycle program for AI-decision records.

### Worked examples

**Multinational bank with U.S.-supervised and EU-supervised entities.** The institution operates U.S. national-bank entities under OCC supervision and EU subsidiaries under EBA / national-competent-authority supervision, with at least one EU subsidiary deploying high-risk AI in scope of the AI Act. Applicable minimums: seven years (U.S. FFIEC, contributing minimum (a)); six months (Article 26(6), contributing minimum (b)); the institution's commercial retention policy (contributing minimum (c), assumed not to exceed seven years for AI-decision evidence). The binding floor is seven years; the U.S. FFIEC retention dominates. The institution's control description says "the institution retains AI-decision audit-trail records for seven years from the date of the originating decision, satisfying the U.S. FFIEC seven-year retention floor and concurrently exceeding the EU AI Act Article 26(6) six-month minimum applicable to the institution's EU-supervised entities." A single retention period operates across all entities; the chain's storage planning is sized to seven years.

**Non-U.S. institution operating in an EU-only regime.** The institution is an EU-headquartered bank without U.S.-supervised entities, deploying high-risk AI in scope of the AI Act. Applicable minimums: Article 26(6) six-month minimum (contributing minimum (b)); national-law retention if any (contributing minimum (b) extension, examples being German GoBD seven-year retention for accounting-relevant records, French Code de commerce ten-year retention for commercial records); the institution's commercial retention policy (contributing minimum (c)). If the institution's national law imposes a seven-year accounting-related retention on AI-decision records that touch financial bookings, that is the binding floor; if the AI-decision records are scoped narrowly enough to be outside any extended national-law retention, the floor is the Article 26(6) six-month minimum or the institution's commercial policy, whichever is longer. The institution's control description states the binding floor and names the rationale.

**Bank in a third jurisdiction (APAC, LATAM).** The institution is supervised in a jurisdiction without a U.S. FFIEC seven-year norm and without explicit EU AI Act applicability — for example, a Singapore-supervised bank operating under MAS Veritas / FEAT, a Japan-supervised bank under JFSA principles, or a Brazil-supervised bank under BCB emerging guidance. Applicable minimums: the local supervisor's records-retention guidance (contributing minimum (d)); the institution's commercial retention policy (contributing minimum (c)); any cross-border arrangement that imports a minimum from another jurisdiction (a Singapore branch of a U.S. bank, for example, may inherit the U.S. seven-year floor through its parent-entity supervision). The institution's compliance team consults the applicable local law, names the binding floor in the control description, and sizes chain storage accordingly.

The chain itself does not change across these scenarios. The retention horizon for the institution's audit-trail records does change, and the institution's control description carries the rationale so a working examiner sees the choice and the reasoning at examination time without needing to reconstruct it from underlying regulation.

## Template institution AI policy framework

The sections above describe which controls the chain satisfies under each external framework. Institutions adopting the chain need a complementary artifact: an institution-side AI policy framework that names the chain as one technical control among many, anchors the chain's role within the institution's broader AI governance, and presents the policy hierarchy to MRM, IT-Risk, Legal, and Compliance in a single readable document.

The template below is a one-page outline an institution adapts to its house style. The sections name the framework hierarchy explicitly so the chain's role is visible in context. The institution edits the placeholders, removes the parenthetical guidance, and presents the result as the institution's AI policy framework. The template is obviously customisable — the institution's AI policy is institution-specific, and a template cannot substitute for the institution's own deliberation about scope, roles, and reporting cadence — but the outline accelerates the conversation between IT-Risk, Legal, and Compliance during chain adoption, especially in institutions where AI policy does not yet exist or does not yet reference logging-integrity controls.

> **`<institution name>` — AI policy framework**
>
> **Version.** `<version number>`. **Effective date.** `<date>`. **Owner.** `<CISO or Chief AI Officer>`. **Approval body.** `<institution's risk committee or AI governance committee>`.
>
> **(a) Scope.** This framework governs the institution's use of AI systems in production. It applies to: AI-driven customer-facing decisions (credit, deposit, payments, fraud); AI-driven internal decisions affecting customers indirectly (collections strategy, marketing eligibility); AI-assisted employee decisions where the AI's output materially shapes the human decision (relationship-management recommendations, investment advice generation). It does not apply to: pure productivity tools (general-purpose summarisation, drafting assistance) where the AI output does not directly or indirectly produce a customer-facing decision. *(The institution scopes "production" against its own AI inventory; pilot and pre-production deployments are governed by the framework once they cross a threshold the institution defines.)*
>
> **(b) Frameworks named.** This framework references and aligns with the external frameworks applicable to the institution's regulatory perimeter. *(The institution names the relevant set; not every institution is subject to every framework. Listed here is the union the chain's documentation supports.)* NIST AI Risk Management Framework, AI RMF 1.0 (2023). Executive Order 14110, Safe, Secure, and Trustworthy Development and Use of Artificial Intelligence (2023). U.S. Treasury Financial Services AI Risk Management Framework (Feb 2026). Federal Reserve SR 11-7, Guidance on Model Risk Management (2011). OCC Bulletin 2011-12, Supervisory Guidance on Model Risk Management. CFPB AI guidance on adverse-action notices and UDAAP analysis. Regulation (EU) 2024/1689, the EU AI Act. European Banking Authority Guidelines on AI in Financial Services (2024+). Regulation (EU) 2022/2554, DORA. Directive (EU) 2022/2555, NIS2. *(Plus jurisdictional frameworks the institution is subject to — UK FCA / PRA, Singapore MAS Veritas / FEAT, Japan JFSA, Brazil BCB, applicable U.S. state laws.)*
>
> **(c) Chain of custody role.** The institution operates the FFIEC chain-of-custody control (`spec/chain-of-custody-v1.md`) as one technical control among many in the institution's AI governance program. The chain provides integrity-bearing records of AI decisions: the per-event HMAC chain, per-day Merkle seal, and HSM-signed root collectively defy retroactive alteration of AI-decision audit-trail data. The chain's coverage area is: integrity-of-decision-evidence at examination time and at customer-dispute time. The chain's coverage area is NOT: AI model validation (the institution's MRM program covers this), AI model performance monitoring (the institution's MRM program covers this), AI training-data integrity (out of scope for chain v1.x), AI fairness or bias testing (the institution's broader compliance program covers this), AI policy compliance (this framework covers it). The chain composes with the institution's other AI-governance controls; it does not substitute for them.
>
> **(d) Roles and responsibilities.** *(The institution names its specific roles; the template below is the typical bank shape.)*
>
> - The MRM committee owns AI model validation, periodic re-validation, and model-inventory completeness. The committee receives chain-derived evidence (verifier output, per-model decision-count distribution) on a documented cadence and uses it to support effective challenge of AI models in scope of SR 11-7.
> - The IT-Risk function owns AI-system operational risk, including chain-control operational health (verifier execution cadence, reconciliation cadence, IR playbook execution evidence). The IT-Risk function reviews the chain's residual-risk register `<reference>` and either accepts the residual posture or documents a deviation with mitigating controls.
> - The Legal function owns regulatory interpretation, including the institution's posture on each external framework named in section (b). Legal also owns the institution's response to regulator-published guidance updates that affect the chain's role in the framework hierarchy.
> - The Compliance function owns the institution's compliance evidence production, including the per-jurisdiction control descriptions referenced in this framework. Compliance prepares the chain-derived evidence for examination engagements (FFIEC IT, CSF, EBA, DORA Art. 28 third-party).
> - The chain-operations team operates the chain control day-to-day: SDK deployment in AI-call paths, ledger-server operation, HSM custody, IKM custody, verifier execution, reconciliation execution, IR playbook execution. The team is typically a subset of the institution's broader cyber-ops or platform-engineering function with the chain-specific authority defined in the institution's `<roles charter>`.
> - The Internal Audit function owns independent assessment, periodically running the verifier against an institution-side sample and reporting findings to the audit committee. Internal Audit's chain-related work is scoped through the institution's audit plan.
>
> **(e) Reporting cadence and escalation.** The MRM committee receives a quarterly chain-evidence summary (verifier output highlights, per-model decision-count distribution, drift indicators). The audit committee receives a semi-annual chain-evidence summary (verifier-output verbatim, reconciliation summary, IR-execution evidence) per the audit-committee-summary cadence. The IT-Risk function receives operational-event monitoring on the cadence the operational-events catalog suggests (continuous for `chain.verification_failure`, weekly for `master.reconciliation_completed`). The chain-operations team escalates per the IR playbook to the institution's incident-management program; major incidents follow the institution's regulator-notification procedures including DORA Art. 17 24-hour and FFIEC 36-hour clocks where applicable.
>
> **(f) Annual policy review.** This framework is reviewed annually by the named owner with input from MRM, IT-Risk, Legal, Compliance, and the chain-operations team. The review covers: the framework hierarchy in section (b) (regulator-published updates, new frameworks the institution has come into scope of); the chain's role in section (c) (spec version updates, new chain primitives, retired primitives); the roles and responsibilities in section (d) (organisational changes, role re-assignments); the reporting cadence in section (e) (committee charter changes, regulator expectations on cadence). The reviewed framework is presented to the approval body for re-adoption.

The template is intentionally short. A one-page artifact is consumed by the people doing the institution's work; a thirty-page artifact is consumed by the auditor confirming the artifact exists. The institution may extend the template — adding sections on training-program for AI-using employees, on AI-vendor governance specifics, on the institution's customer-communication posture for AI-driven decisions — but the core sections above are the load-bearing items an examiner expects to see when asking "where is your AI policy framework, and what role does the chain play in it?"

The chain's role anchored in the framework hierarchy is the load-bearing presentation point. The institution that says "the chain is one technical control among many; here is what it covers, here is what it does not, here is who owns it, here is how often we report it, here is when we review the framework" presents a defensible governance posture. The institution that has adopted the chain without naming the chain's role in a written framework is partially exposed — the chain's evidence holds, but the policy framing around it is missing, and an MRM examiner working SR 11-7 governance discipline asks for the framing first and the evidence second.
