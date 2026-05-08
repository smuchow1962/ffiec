---
title: EU Articulation Extension — AI Act Article 12, EU-QTSA Time Stamps, Multilingual Adverse-Action Schema
status: informative
aligned-with:
  - Regulation (EU) 2024/1689 (AI Act)
  - Regulation (EU) 910/2014 (eIDAS) and Regulation (EU) 2024/1183 (eIDAS 2.0)
  - Regulation (EU) 2022/2554 (DORA)
  - National anti-discrimination statutes named below
extends: dora-articulation-overlay.md
date: 2026-05-07
version: 1.0.0
---

# EU Articulation Extension

> **What this doc is.** A focused extension to `dora-articulation-overlay.md` covering three areas the DORA overlay does not articulate explicitly: (1) AI Act Regulation (EU) 2024/1689 Article 12 logging conformance, (2) EU-QTSA-listed Qualified Electronic Time Stamp guidance for v1.0a deployments under spec §10.14, and (3) generalisation of the spec §10.11 adverse-action notice schema for non-ECOA EU jurisdictions. The DORA overlay is the load-bearing source for DORA, eIDAS Article 25/26 signature qualification, EBA outsourcing flow-down, NIS2 incident notification, and Schrems II supplementary measures; this extension does not duplicate those sections.

> **What this doc is NOT.** Not a normative extension to the v1.0a specification. Not a replacement for the DORA overlay. The integrity primitives in v1.0a are framework-neutral; the supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) governs: the chain proves what was said and that the record was not tampered with after capture; it does not prove the substantive correctness of the captured content.

---

## 1. Reading order and cross-reference to the DORA overlay

This extension assumes the reader has consumed `dora-articulation-overlay.md`. The two documents partition the EU articulation space as follows.

| Topic | Document | Section |
|---|---|---|
| DORA Articles 28-30 (ICT third-party register) | `dora-articulation-overlay.md` | §2-§3 |
| DORA Article 19 (incident classification and notification clocks) | `dora-articulation-overlay.md` | §4 |
| DORA Articles 26-27 (TLPT under TIBER-EU / CBEST) | `dora-articulation-overlay.md` | §5 |
| DORA Articles 6-9, 12 (RTO/RPO and BCP) | `dora-articulation-overlay.md` | §6 |
| EBA/GL/2019/02 outsourcing audit rights | `dora-articulation-overlay.md` | §7 |
| eIDAS Article 25/26 signature qualification (AdES vs QES) | `dora-articulation-overlay.md` | §8 |
| NIS2 timing and CSIRT coordination | `dora-articulation-overlay.md` | §9 |
| DORA Article 29 concentration risk | `dora-articulation-overlay.md` | §10 |
| Schrems II supplementary measures | `dora-articulation-overlay.md` | §11 |
| GDPR Article 22 / AI Act Article 86 right to explanation | `dora-articulation-overlay.md` | §12 |
| ECB SSM SREP integration | `dora-articulation-overlay.md` | §13 |
| **AI Act Article 12 logging conformance** | **this extension** | **§2** |
| **EU-QTSA-listed time-stamp authorities for §10.14** | **this extension** | **§3** |
| **Multilingual / non-ECOA adverse-action schema** | **this extension** | **§4** |
| **EU jurisdiction tenant naming under Pattern B** | **this extension** | **§5** |

Three reading orders.

| Reader | Reading order |
|---|---|
| AI Act compliance counsel for a high-risk system | §2 (AI Act Article 12 mapping) → §4 (multilingual adverse-action schema) → DORA overlay §12 (right to explanation) |
| Institution selecting a Trusted-Service Provider for time stamps | §3 (EU-QTSA TSAs) → DORA overlay §8.4 (forward-note, now resolved by §3 of this document) |
| Multi-jurisdiction EU bank deploying Pattern B | §5 (EU jurisdiction tenant naming) → DORA overlay §11 (Schrems II) → DORA overlay §13 (SSM) |

---

## 2. AI Act Article 12 logging conformance — chain attribute mapping

### 2.1 What Article 12 requires

Regulation (EU) 2024/1689 Article 12 imposes three requirements on high-risk AI systems.

- Article 12(1). High-risk AI systems shall technically allow for the automatic recording of events ('logs') over the lifetime of the system.
- Article 12(2). The logging capabilities shall ensure a level of traceability of the AI system's functioning that is appropriate to the intended purpose of the system, and shall facilitate the post-market monitoring referred to in Article 72.
- Article 12(3). For high-risk AI systems referred to in point 1(a) of Annex III (remote biometric identification systems), and for other high-risk systems where the technical specifications so provide, logs shall record at minimum (a) the period of use, (b) the reference databases used, (c) the input data for which the search led to a match, and (d) the identification of natural persons involved in result verification.

The chain captures substantively the right content for any high-risk financial-services AI system in scope of Annex III point 5 (creditworthiness, credit-scoring, life and health insurance pricing, emergency services dispatch). The articulation gap is vocabulary, not substance: an institution adopting the chain for an AI Act high-risk-system needs to name which chain attributes satisfy which Article 12 requirement. This section provides that mapping for the most common financial-services high-risk shapes.

### 2.2 Article 12(2) lifetime traceability — chain attribute mapping

| Article 12(2) requirement | Chain attribute or event satisfying the requirement |
|---|---|
| Log records the period of use of the AI system | The chain entry's `received_at` field (spec §4.4 normative) records the wall-clock arrival time of each event under MAC binding; the first and last entries of a `chain_id` delineate the per-run usage period; the institution's CC8.1 names the chain-id-to-deployment mapping that aggregates per-run periods into a per-system lifetime record |
| Log facilitates post-market monitoring under Article 72 | The chain's verifier output (spec §7) provides byte-reproducible per-event verification suitable for the post-market monitoring system's evidence intake; the institution's CC8.1 names the cadence at which verified chain extracts are fed to the post-market monitoring repository |
| Log produces evidence of model identity | The `gen_ai.response.model` attribute (spec §4.4 normative) records the model identifier the vendor's API answered with; combined with `gen_ai.request.model` (the requested identifier), the chain detects vendor-side silent rerouting between requested and answered models |
| Log produces evidence of decision input | RFC 8785-canonicalised JSON of the application content per spec §4.1 inviolate property 7 — the institution's payload-redaction posture (spec §4.4 informative) names which input fields are integrity-bound and which are redacted with a fingerprint pointer |
| Log produces evidence of decision output | The `gen_ai.response.text` (or model-specific output attribute) is integrity-bound; for multi-step agent decisions, the chain's parent-linkage (`audit.parent.event_id`) reconstructs the decision path |
| Log produces evidence of human-in-the-loop verification | The `audit.routing.reviewer_identity` attribute (spec §4.4.1 normative routing schema) records the natural person who reviewed the decision; for systems without human-in-the-loop, the field is absent and the institution's policy documents the absence |

The mapping is informative — institutions extend it with operational-context attributes — but provides the baseline mapping every EU institution would otherwise write itself.

### 2.3 Article 12(3) Annex III point 1(a) — remote biometric identification

The chain v1.0a is not designed for biometric identification systems and does not provide the reference-database integration Article 12(3)(b) requires. EU institutions deploying biometric identification systems use a separate audit primitive for the biometric-match record; the chain may bind the institution's downstream decision (e.g., whether to act on the match) but does not bind the match itself. The institution's CC8.1 names the biometric-match audit primitive separately from the chain.

### 2.4 Annex III point 5 — financial-services high-risk shapes

The financial-services high-risk shapes (creditworthiness, credit-scoring, life and health insurance pricing) are the natural deployment targets for the chain. The Article 12 mapping for these shapes uses the standard chain attribute schema without extension. The institution's CC8.1 names the high-risk classification per AI Act Annex III, the corresponding model identifier captured in `gen_ai.response.model`, and the post-market monitoring procedure under Article 72.

### 2.5 Provider-vs-deployer responsibilities under Articles 9, 13, 14

A practical articulation: the chain serves both the AI Act provider (model vendor) and the AI Act deployer (Union financial entity) — two distinct legal categories with overlapping but non-identical obligations.

| Party | AI Act category | Chain artifact supporting compliance |
|---|---|---|
| Model vendor (Anthropic, OpenAI, Google, Mistral) | Provider under Article 9 (risk-management system); Article 13 (transparency); Article 14 (human oversight design) | Vendor's model card and technical documentation, fed by the chain's `gen_ai.response.model` field as observed in production. The chain does not produce the model card, but it observes the deployed model identifier, which the vendor's documentation must match. Discrepancy is detectable from the chain |
| Union financial entity | Deployer under Article 26 (deployer obligations); Article 86 (right to explanation, see DORA overlay §12) | Chain entries plus the institution's deployer-side documentation: Article 26(5) automatic-logging recordkeeping is satisfied by the chain; Article 26(7) human oversight by the institution's `audit.routing.reviewer_identity` capture; Article 86 right to explanation by the procedure named in DORA overlay §12 |

The AI Act provider obligations stay with the vendor; the chain helps the deployer detect when the vendor's published model identifier diverges from what the deployment actually invokes.

---

## 3. EU-QTSA-listed time stamps for spec §10.14 — resolving the v1.0a posture

### 3.1 The Article 41 presumption

eIDAS Article 41 confers on a Qualified Electronic Time Stamp a presumption of accuracy of the date and time it indicates and of the integrity of the data to which the date and time are bound. A non-qualified RFC 3161 time stamp is technically valid as an integrity binding but does NOT carry the Article 41 presumption — the burden of proof on the date and time stays with the institution producing the time stamp.

The DORA overlay §8.4 names RFC 3161 trusted-time integration as a v1.x extension candidate. **This extension resolves the v1.0a posture explicitly: institutions deploy RFC 3161 trusted time stamps under spec §10.14 today, choosing an EU-Trusted-List-listed TSA where the Article 41 presumption is required.** No spec-text change is needed; the resolution is at the institution's procurement and CC8.1 layer.

### 3.2 EU-QTSA TSA selection guidance

The European Commission maintains the EU Trusted List of Qualified Trust Service Providers under eIDAS Article 22. TSPs offering Qualified Time Stamp services appear on the list under the service type `QTimestamp`. The institution selects a TSP from the Trusted List for its jurisdiction or any Member State. Common choices for Union financial entities include the QTSPs operated by DigiCert (EU QTSA), GlobalSign (EU QTSA), Sectigo (EU QTSA), Atos QTSA, and the national QTSPs operated by Member-State infrastructures (e.g., InfoCert in Italy, Buypass in Norway/EEA, Trustpro in Czechia).

The institution's CC8.1 names:

- The TSP selected from the Trusted List, by name and Trust List entry identifier.
- The TSA endpoint URL the institution's RFC 3161 client invokes.
- The TSA certificate chain stored in the institution's tenant configuration for time-stamp verification.
- The cadence at which the institution requests time stamps (per-event, per-batch, per-seal, or per-day).

### 3.3 RFC 3161 client integration with the chain

The RFC 3161 client integrates with the chain at one of three points, depending on the institution's evidence requirements.

| Integration point | Frequency | Article 41 presumption attaches to |
|---|---|---|
| Per-event time stamp | Every chain entry | Each entry's MAC payload |
| Per-seal time stamp | Each daily Merkle seal | The seal's signed root and the events covered by the seal |
| Per-day time stamp | Once per UTC day | The day's accumulated seal records |

The per-seal pattern is the operational sweet spot for most deployments: one TSA round-trip per UTC day, Article 41 presumption attaches to the seal's signed root, and the chain's per-event MAC binds individual events to the sealed root via the standard verification procedure (spec §7 step 10). The institution's CC8.1 names the chosen integration point.

### 3.4 Multi-jurisdiction TSA choice

Institutions operating across Member States may use multiple TSAs (e.g., a German QTSA for tenants under BaFin supervision, a French QTSA for tenants under ACPR). The chain accommodates this through the per-tenant TSA configuration in the institution's tenant-configuration store; the verifier consumes the per-tenant TSA certificate chain at verification time. The institution's CC8.1 names the per-tenant TSA mapping and the cross-jurisdiction correlation procedure for tenants whose operations span Member States.

### 3.5 Non-EU operations

Institutions operating outside the Union may use a non-Trusted-List TSA (e.g., a US-based RFC 3161 TSA) for operations not subject to the Article 41 presumption. The chain's integrity binding is unchanged; what differs is the legal weight of the time stamp in EU proceedings. The institution's CC8.1 names which tenants operate under the EU-QTSA TSA and which operate under a non-Trusted-List TSA, with the rationale documented in the institution's record of supervisory exposure.

---

## 4. Multilingual / non-ECOA adverse-action schema

### 4.1 The §10.11 schema and its US-jurisdiction scope

Spec §10.11 defines the `audit.ecoa.adverse_action.*` schema for ECOA adverse-action notices required under 15 U.S.C. §1691(d) and 12 CFR Part 1002. The schema names the language of the notice (`audit.ecoa.adverse_action.notice_language` ∈ {`en`, `es`, `vi`, `ko`, `zh`, `tl`}), the jurisdiction (`audit.ecoa.adverse_action.jurisdiction` = `us`), and the regulatory basis (`audit.ecoa.adverse_action.regulatory_basis` = `ecoa`).

The schema's namespace is jurisdictional. EU institutions issuing equivalent adverse-action notices under national anti-discrimination law need a parallel schema that covers the same operational ground (notice language, jurisdiction, regulatory basis, decision-factor enumeration) without misnaming the regulatory basis.

### 4.2 Parallel `audit.adverse_action.*` schema for non-US jurisdictions

This extension defines a parallel `audit.adverse_action.*` schema with the same eight attributes as `audit.ecoa.adverse_action.*`, generalised to any jurisdiction. The two schemas coexist: US institutions continue to use `audit.ecoa.adverse_action.*`; EU and other non-US institutions use `audit.adverse_action.*`. Both are valid v1.0a chain attributes; the institution's CC8.1 names which schema applies for which jurisdiction.

| Attribute | Type | Description |
|---|---|---|
| `audit.adverse_action.event_id` | string (UUID) | The unique identifier of the adverse-action notice event |
| `audit.adverse_action.parent_decision_event_id` | string (UUID) | The parent chain event that produced the adverse decision |
| `audit.adverse_action.notice_language` | string (BCP 47) | The language code of the notice as delivered to the consumer (e.g., `de-DE`, `fr-FR`, `it-IT`, `nl-NL`, `pl-PL`, `sv-SE`) |
| `audit.adverse_action.jurisdiction` | string | The jurisdiction whose statute requires the notice (e.g., `de`, `fr`, `it`, `nl`, `pl`, `se`) |
| `audit.adverse_action.regulatory_basis` | string | The statutory or regulatory citation (e.g., `de-agg-19` for German Allgemeines Gleichbehandlungsgesetz §19; `fr-code-conso` for the French Code de la consommation; `nl-awgb` for the Dutch Algemene wet gelijke behandeling; `it-codice-pari-opportunita` for the Italian Codice delle pari opportunità) |
| `audit.adverse_action.delivery_timestamp` | string (RFC 3339 UTC) | When the notice was delivered to the consumer |
| `audit.adverse_action.decision_factors` | array of strings | The statutory-required enumeration of decision factors, in plain language |
| `audit.adverse_action.consumer_rights_disclosed` | array of strings | The consumer rights named in the notice (e.g., `right_to_human_review`, `right_to_object`, `right_to_explanation`) |

### 4.3 Per-jurisdiction `regulatory_basis` enumeration

The institution's CC8.1 names the per-jurisdiction regulatory basis it uses. The following table is illustrative for the most common EU jurisdictions; institutions extend it as their operations span additional Member States.

| Member State | Statute | `audit.adverse_action.regulatory_basis` value |
|---|---|---|
| Germany | Allgemeines Gleichbehandlungsgesetz §19 | `de-agg-19` |
| France | Loi n° 2008-496 du 27 mai 2008; Code de la consommation L312-16 | `fr-loi-2008-496` |
| Italy | Decreto legislativo 198/2006 (Codice delle pari opportunità) | `it-codice-pari-opportunita` |
| Netherlands | Algemene wet gelijke behandeling | `nl-awgb` |
| Spain | Ley Orgánica 3/2007 | `es-lo-3-2007` |
| Belgium | Loi du 10 mai 2007 / Wet van 10 mei 2007 | `be-loi-2007-05-10` |
| Sweden | Diskrimineringslagen (2008:567) | `se-dl-2008-567` |
| Poland | Ustawa o równym traktowaniu (3 December 2010) | `pl-uort-2010` |
| Austria | Gleichbehandlungsgesetz | `at-gbg` |
| Ireland | Equal Status Acts 2000-2018 | `ie-esa-2000-2018` |

The `regulatory_basis` enumeration is institution-extensible. The institution's CC8.1 documents the values it uses; the spec verifier accepts any string-typed value.

### 4.4 Notice language under BCP 47

The `audit.adverse_action.notice_language` field uses BCP 47 language tags (e.g., `de-DE`, `fr-FR`, `nl-NL`). For Member States with multiple official languages, the institution names the specific language tag used for the notice (e.g., `nl-BE` for Dutch in Belgium, `fr-BE` for French in Belgium, `de-BE` for German in Belgium; `fi-FI` for Finnish in Finland, `sv-FI` for Swedish in Finland; `mt-MT` for Maltese in Malta, `en-MT` for English in Malta).

The chain's canonical bytes are UTF-8 per RFC 8785; the encoding accommodates the full Unicode range without modification. The institution's notice-generation tooling produces the notice in the recipient's preferred language; the chain captures the tag, the delivery timestamp, and the decision-factor enumeration. The notice content itself may be redacted with a fingerprint pointer per the institution's payload-redaction posture (spec §4.4 informative).

### 4.5 Cross-reference to GDPR Article 22 / AI Act Article 86 right to explanation

The `audit.adverse_action.consumer_rights_disclosed` field names which consumer rights the notice disclosed. Under GDPR Article 22 and AI Act Article 86, the deployer must disclose the right to obtain human review of an automated decision and the right to contest the decision. The DORA overlay §12 names the procedure that operationalises these rights using the chain as evidence; this extension's schema captures the disclosure event in the chain itself.

---

## 5. EU jurisdiction tenant naming under spec §10.15 Pattern B

### 5.1 Pattern B as the recommended posture for multi-jurisdiction EU institutions

Spec §10.15 defines Pattern B (per-region tenant_id) as the resilience pattern for institutions with per-region or per-jurisdiction event-isolation requirements. For EU institutions operating across Member States under ECB SSM consolidated supervision (for significant institutions) or under direct national supervision (for less significant institutions), Pattern B is the recommended posture: per-jurisdiction tenant_ids align the chain with the SSM/national-supervisor topology and support examination-data-sovereignty rules under each national procedural code.

The DORA overlay §13 names the SSM SREP integration; this extension names the per-jurisdiction tenant naming and HSM custody pinning that operationalise Pattern B for EU deployments.

### 5.2 Per-jurisdiction tenant naming convention

The institution adopts a per-jurisdiction tenant_id naming convention that maps each national supervisor to a stable tenant identifier. The convention below is illustrative; institutions adopt their own naming as long as it is documented in CC8.1 and stable across the chain's retention period.

| Jurisdiction | National supervisor | Recommended tenant_id pattern | HSM region |
|---|---|---|---|
| Germany | BaFin | `bank-de` | AWS CloudHSM `eu-central-1`, Azure Managed HSM Germany West Central, Google Cloud HSM `europe-west3` |
| France | ACPR | `bank-fr` | AWS CloudHSM `eu-west-3`, Azure Managed HSM France Central, Google Cloud HSM `europe-west9` |
| Italy | Banca d'Italia | `bank-it` | AWS CloudHSM `eu-south-1`, Azure Managed HSM Italy North, Google Cloud HSM `europe-west8` |
| Netherlands | DNB | `bank-nl` | AWS CloudHSM `eu-west-1` (Netherlands operations), Azure Managed HSM West Europe, Google Cloud HSM `europe-west4` |
| Spain | Banco de España | `bank-es` | AWS CloudHSM `eu-south-2`, Azure Managed HSM Spain Central, Google Cloud HSM `europe-southwest1` |
| Belgium | NBB | `bank-be` | AWS CloudHSM `eu-west-1` (Belgium proximity), Azure Managed HSM West Europe, Google Cloud HSM `europe-west1` |
| Sweden | Finansinspektionen | `bank-se` | AWS CloudHSM `eu-north-1`, Azure Managed HSM Sweden Central, Google Cloud HSM `europe-north1` |
| Poland | KNF | `bank-pl` | AWS CloudHSM `eu-central-1` (closest), Azure Managed HSM Poland Central, Google Cloud HSM `europe-central2` |
| Ireland | Central Bank of Ireland | `bank-ie` | AWS CloudHSM `eu-west-1`, Azure Managed HSM North Europe, Google Cloud HSM `europe-west2` |
| Finland | FIN-FSA | `bank-fi` | AWS CloudHSM `eu-north-1`, Azure Managed HSM North Europe, Google Cloud HSM `europe-north1` |

### 5.3 ECB SSM consolidated cross-tenant correlation

For significant institutions consolidated under the ECB SSM, cross-tenant correlation reads as institution-side evidence. The institution's CC8.1 names the cross-jurisdiction correlation procedure that the institution operates internally to satisfy SSM consolidated review while preserving per-jurisdiction examination access. The chain's per-tenant HKDF binding (spec §4.1) means a tenant's events cannot mechanically be lifted into another tenant's chain — the cross-chain-lift detection at spec §7 step 4 is the cryptographic floor that makes the per-jurisdiction segregation robust under nominal SSM cross-tenant analysis.

### 5.4 Cross-reference to Schrems II supplementary measures

The DORA overlay §11 names the Schrems II supplementary measures that apply to vendor-hosted EU controllers. The per-jurisdiction tenant naming in this section composes with §11 of the DORA overlay: each tenant's HSM is custody-pinned to an EU region, the seal-signing key remains under EU jurisdiction, and the institution's encryption-at-rest custody is in the same Member State as the tenant's national supervisor. The vendor cannot decrypt the tenant's chain content without the institution's CMEK or the institution's bring-your-own-key arrangement, even where the vendor's compute or storage transits non-Union infrastructure.

---

## 6. Operational checklist for an EU institution deploying the chain

The institution working through DORA + AI Act + GDPR + national-supervisor compliance produces the following artefacts. This checklist names the minimum set; the institution's compliance department extends it as their specific exposure requires.

| Artefact | Source | Reviewed by |
|---|---|---|
| AI Act high-risk classification per Annex III | Institution's AI policy team | Risk & compliance committee |
| AI Act Article 12 logging mapping (per §2 of this extension) | Institution's CC8.1 | National competent authority during examination |
| Article 26 deployer-obligations checklist | Institution's CC8.1 + DORA overlay §3 | Joint Examination Team during SSM SREP |
| eIDAS signature qualification declaration (AdES / AdES-QC / QES per DORA overlay §8) | Institution's CC8.1 | Institution's legal counsel |
| EU-QTSA TSA selection (per §3 of this extension) | Institution's procurement + CC8.1 | Institution's IT security |
| `audit.adverse_action.regulatory_basis` per-jurisdiction enumeration (per §4 of this extension) | Institution's CC8.1 | Institution's consumer-protection counsel |
| Pattern B per-jurisdiction tenant_id naming (per §5 of this extension) | Institution's CC8.1 | National competent authorities for each jurisdiction the institution operates in |
| DORA Article 28(3) register of information | Institution's procurement + CC8.1 | National competent authority during DORA examination |
| EBA/GL/2019/02 outsourcing audit-rights flow-down | Institution's procurement + DORA overlay §7 | National competent authority during outsourcing examination |
| NIS2 incident-notification timing matrix per DORA overlay §9 | Institution's IR runbook | National CSIRT |
| DORA Article 19 incident-classification procedure per DORA overlay §4 | Institution's IR runbook | National competent authority |
| Schrems II supplementary measures per DORA overlay §11 | Institution's CC8.1 + RoPA | Institution's DPO + national DPA |

The checklist is operational. The chain itself produces the integrity-bearing audit trail; the institution's CC8.1 declares the supervisory posture; the supervisory authorities consume the chain's verifier output and the institution's CC8.1 jointly during examination.

---

## 7. Bottom line

The chain v1.0a is conformant with EU regulatory architecture for the three areas this extension articulates — AI Act Article 12 logging, EU-QTSA-listed Qualified Electronic Time Stamps under §10.14, and non-ECOA adverse-action schema for EU jurisdictions. The conformance is operational: the institution's CC8.1 names the AI Act mapping, the TSA selection, the per-jurisdiction `regulatory_basis` enumeration, and the per-jurisdiction tenant naming; the spec produces the integrity-bearing audit trail under any of these institutional choices.

The DORA overlay is the load-bearing source for DORA, eIDAS Article 25/26, EBA outsourcing flow-down, NIS2, and Schrems II. This extension closes the three remaining EU articulation areas the DORA overlay does not articulate explicitly. Together, the two documents are the complete EU articulation set for v1.0a deployments by Union financial entities.

For institutions seeking the AI Act provider-side mapping (model vendors reading the chain to detect deployment-time discrepancy with their published model cards), the chain's `gen_ai.response.model` field is the load-bearing observation. The institution's CC8.1 names the cadence at which the institution reconciles observed model identifiers with vendor-published model cards; discrepancies are flagged through the institution's vendor-management procedure under DORA Article 28(2).

For institutions whose national supervisor has issued specific guidance beyond the SSM SREP framework (e.g., a national-language disclosure rule, a national-cybersecurity baseline beyond NIS2 transposition), the institution's CC8.1 names the national overlay separately. This extension does not enumerate every national supervisor's specific guidance; the per-jurisdiction `regulatory_basis` mechanism (§4.3) and the per-jurisdiction tenant naming (§5.2) compose to accommodate any national overlay the institution faces.

---

## Document control

| Field | Value |
|---|---|
| Document | EU Articulation Extension |
| Version | 1.0.0 |
| Status | Informative — institution-side articulation; no normative spec change |
| Aligned with | AI Act Reg (EU) 2024/1689; eIDAS Reg (EU) 910/2014 + Reg (EU) 2024/1183; DORA Reg (EU) 2022/2554 |
| Extends | `dora-articulation-overlay.md` |
| Date | 2026-05-07 |
