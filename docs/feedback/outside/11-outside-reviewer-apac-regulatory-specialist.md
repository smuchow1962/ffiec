# Outside reviewer — APAC regulatory technology specialist (ex-Japan)

**Reviewer.** Mei Chen. Singapore-based managing director at a Big Four APAC regulatory-technology practice. Background: six years at the Monetary Authority of Singapore (MAS) Technology Risk Supervision Department where she co-authored the 2021 revision of the MAS Technology Risk Management Guidelines and contributed to the FEAT Principles working group on Fairness, Ethics, Accountability, and Transparency in AI/data analytics; three years at the Hong Kong Monetary Authority (HKMA) Banking Supervision Department on TM-G-1 implementation and the 2024 HKMA AI principles circular; eight years in Big Four APAC regtech advising banks and insurers across Singapore, Hong Kong, Sydney, Mumbai, and Seoul. Holds CISA, FRM, and IBF Advanced (Institute of Banking and Finance Singapore). Active in the MAS Veritas initiative for responsible AI assessment in financial services. Three Tier-1 APAC bank clients are evaluating chain-of-custody infrastructure: a Singapore-headquartered regional bank with operations in eight ASEAN markets, an Australian Tier-1 with operations in NZ/HK/Singapore, and an Indian private-sector bank under RBI supervision.

**Why this review matters.** Mei's three Tier-1 APAC bank clients are evaluating the FFIEC chain-of-custody v1.0 spec as a candidate technical standard. Her sign-off would unblock pilot deployments under their MAS / HKMA / APRA / RBI compliance programs respectively. The institutions are not US-regulated; they read the spec as a technical-merit document and the question is whether it composes with APAC regulatory architecture without forcing each institution to write a substantial APAC-conformance overlay. Mei's lens is explicitly ex-Japan — the Japan-specific regulatory architecture (FSA, BoJ, PIPC, J-SOX, APPI, CRYPTREC) is covered in drop #06 by Hiroshi Nakamura. Mei's territory is the rest of APAC: Singapore, Hong Kong, India, Australia, Korea, Malaysia, Indonesia, Thailand, the Philippines.

**Reading angle.** Principle-based regulatory frameworks (vs. rule-based) — APAC supervisors lean toward principle-based regulation more than US/EU. MAS FEAT Principles, HKMA AI principles, RBI guidelines on responsible AI, APRA CPS 230 (Operational Risk Management) and CPS 234 (Information Security), Bank Negara Malaysia AI guidelines, Bank of Thailand digital-services guidelines, FSC Korea financial-IT supervision. Cross-border data flows under the Indian Digital Personal Data Protection Act 2023, Singapore PDPA, Hong Kong PDPO, Australian Privacy Act 1988 (with Privacy and Other Legislation Amendment Bill 2024), Korean PIPA. Multi-jurisdiction supervision across APAC trading hubs (Singapore-Hong Kong axis), the ASEAN Banking Integration Framework, the Asia-Pacific Group Payments and Securities Settlement Standards. Mei reads the spec as a technical-merit document; her question is whether an APAC institution adopting it can satisfy the regional regulatory framework without writing a substantial conformance overlay each.

**Distinct from Hiroshi (drop #06, Japan).** Hiroshi's drop covers Japan-specific regulatory architecture (FSA, BoJ, PIPC, J-SOX, APPI, CRYPTREC, JIS X 0208 character handling, Reiwa era calendar). Mei's drop covers APAC ex-Japan; some structural parallels exist (operational-resilience frameworks across regions), but the substantive content of APAC ex-Japan regulation is distinct from Japan-specific regulation. The two drops are complementary; an APAC institution operating in both Japan and ex-Japan (e.g., a Tier-1 with HQ in Singapore, ops in Tokyo and Hong Kong) needs both perspectives.

**Distinct from Henrik (drop #08, EU).** Henrik covers EU regulatory architecture (DORA, AI Act, eIDAS, EBA Guidelines, ECB SSM). The EU's regulatory style is rule-based and statute-driven; APAC ex-Japan is more principle-based with detailed implementation guidance. Where Henrik finds gaps in EU statute mapping, Mei finds gaps in APAC principle-based assessment frameworks (FEAT, Veritas, HKMA AI principles). The gap shapes are different.

**Stopping criterion.** Eight to twelve questions covering the APAC regulatory-technical convergence: MAS FEAT Principles mapping; MAS Veritas assessment framework composition; HKMA TM-G-1 and AI principles alignment; APRA CPS 234 information security alignment; APRA CPS 230 operational risk alignment; RBI guidelines on AI in financial services; Indian DPDP Act 2023 cross-border restrictions; APAC cloud HSM regional listings; multi-jurisdiction supervision across MAS / HKMA / APRA / RBI; APAC incident-notification timing matrix. Mei does not file Partials for issues already covered by Hiroshi (Japan) or Henrik (EU). Her threshold for a Partial is unmanageable APAC regulatory exposure — a gap that would force the institution to write a substantial APAC conformance overlay that ought to be in the spec or the supplementary documents.

---

## Headline

The chain's substantive cryptographic and procedural choices are NIST/IETF-standardised and APAC-conformance-friendly — APAC supervisors typically accept FIPS 140-2 Level 3, NIST-standardised cryptography, and RFC-standardised wire formats. **My read pass found nine places where APAC regulatory architecture creates obligations the spec does not address — clustering around principle-based assessment frameworks (which APAC supervisors prefer over rule-based statute mapping), regional supervisory architecture, and APAC-specific cross-border data restrictions.** A pragmatic close-out path exists for each — most are informative annexes or supplementary regulator-pack documents, not normative spec changes. The spec's strongest item from an APAC perspective is the byte-reproducible verifier output and the test-vector corpus; APAC supervisors are receptive to demonstrable test artifacts.

**What the spec does well for APAC institutions.** The cryptographic primitives are NIST/IETF-standardised — APAC supervisors do not maintain regional cryptographic standards equivalent to CRYPTREC (which is Japan-specific). The §10.5 cloud HSM products have APAC region availability. The §10.15 multi-region resilience patterns (Pattern A active-active with seal-region pinning; Pattern B per-region tenant_id) align with APAC institutions' multi-jurisdiction operations across SG/HK/TH/SG/MY/ID/IN/AU/KR. The verifier output is byte-for-byte reproducible — APAC supervisors are receptive to this and it composes well with MAS Veritas assessment artifacts.

**What the spec needs for APAC regulatory conformance.** Annex documenting MAS FEAT Principles mapping (Singapore + ASEAN convergence); MAS Veritas assessment-framework composition (the chain provides Accountability and Transparency primitives that compose with Veritas); HKMA TM-G-1 and AI principles alignment; APRA CPS 234 / CPS 230 alignment for Australian institutions; RBI guidelines on AI in financial services for Indian institutions; APAC cross-border data restriction matrix; APAC cloud HSM regional listing; APAC supervisory architecture mapping; APAC incident-notification timing matrix.

Each question below is a one-paragraph gap. The closing is an annex document, a supplementary regulator-pack file, an informative spec paragraph, or a supplementary APAC-conformance overlay the working group could endorse without disrupting the spec's normative core.

---

## Supporting questions

### AQ-1. MAS FEAT Principles mapping — chain-attribute composition with Fairness / Ethics / Accountability / Transparency

**Question.** MAS FEAT Principles (2018 original; 2024 revision under public consultation) govern AI/data analytics use in Singapore financial services. The four principles are: Fairness (justifiable, with regular review), Ethics (aligned with ethical standards, values, and codes of conduct), Accountability (internal accountability for use of AIDA, external accountability for use's impact), and Transparency (proactive communication, opportunity for review). The chain provides direct primitives for Accountability (every decision is integrity-bound to the deciding model + parameters + inputs) and Transparency (the chain produces queryable evidence of decisions for review). Fairness and Ethics are application-domain decisions, not chain primitives, but the chain's `audit.routing.*` and decision-rationale attributes are the substrate institutions use to evidence their fairness/ethics assessments. MAS Veritas — the responsible-AI assessment toolkit — produces fairness assessments that should compose with chain artifacts. The spec doesn't speak FEAT vocabulary, so each Singapore institution writes its own FEAT mapping. This duplicates effort across institutions and produces inconsistent FEAT posture.

**Status:** Partial.

**Closing language proposal.** A new informative annex (parallel to Henrik's proposed AI Act Article 12 mapping) titled "MAS FEAT Principles mapping (informative for Singapore-supervised institutions)" naming each FEAT principle and the chain attribute(s)/event(s) that fulfill the principle's evidentiary requirement. For example: Accountability Principle 4 ("AIDA decisions are explainable, transparent, and fair") → `audit.routing.decision_rationale` (institution's own decision-rationale attribute on chain entries) + `audit.deployment.intent` schema (per §4.4.2) + `gen_ai_parameters` capture; Transparency Principle 2 ("Customers, when affected by AIDA decisions, are provided clear explanations") → ECOA-style adverse-action notice schema (currently in §10.11, generalises per Henrik's EQ-8 to non-ECOA contexts). The annex is informative — institutions extend it with operational-context attributes — but provides the baseline mapping every Singapore institution would otherwise write itself.

### AQ-2. MAS Veritas assessment-framework composition

**Question.** MAS Veritas (2020 launch, 2022 phase 2 with broader scope, 2024 phase 3 in pilot) is MAS's responsible-AI assessment toolkit consisting of Fairness, Ethics, Accountability, and Transparency assessment methodologies (the FEAT-named toolkit per AQ-1). Veritas phase 2 published the Fairness assessment methodology and phase 3 is publishing the Ethics assessment methodology. Singapore institutions deploying high-stakes AI undergo Veritas assessment (formally voluntary; effectively required for MAS approval of new AIDA deployments). Veritas assessment requires evidence the chain naturally produces (decision history, parameter capture, input/output binding, fairness-cohort analysis substrate). The chain SHOULD compose well with Veritas, but the composition isn't spelled out anywhere — institutions and Veritas-assessing firms write their own composition.

**Status:** Partial.

**Closing language proposal.** A new supplementary document `docs/apac-regulator-pack/mas-veritas-composition.md` (or similar) naming how the chain artifacts feed the Veritas assessment workflow. Specifically: chain entries provide the per-decision substrate Veritas's Fairness assessment runs against; chain's `audit.routing.*` decision-rationale attributes provide the substrate Veritas's Accountability assessment runs against; chain's `audit.deployment.intent` (per §4.4.2) provides the substrate Veritas's Transparency assessment runs against. The composition is informative and the Veritas methodology evolves; the document references the current Veritas phase and tracks updates as Veritas evolves. This is an APAC-conformance overlay; not a normative spec change.

### AQ-3. HKMA TM-G-1 and AI principles alignment

**Question.** HKMA TM-G-1 (General Principles for Technology Risk Management) is the foundational HKMA technology-risk supervisory module; the 2024 HKMA circular on AI principles in retail banking adds AI-specific governance expectations. HKMA-supervised institutions (most of HK's 150+ Authorized Institutions) must demonstrate compliance with both. The chain's HSM-custody requirements (§10.5) align with TM-G-1's cryptographic-key-management expectations. The chain's append-only enforcement (§4.1) and reconciliation procedures (§10.1) align with TM-G-1's audit-log-integrity expectations. The chain's `audit.routing.*` and decision-rationale schemas align with the HKMA AI principles' transparency expectations. But the alignment isn't documented anywhere — each HK institution writes its own TM-G-1 mapping.

**Status:** Partial.

**Closing language proposal.** A new supplementary document `docs/apac-regulator-pack/hkma-tm-g-1-mapping.md` naming each TM-G-1 module section and the chain prescription that fulfills it. The mapping is informative and HKMA-specific; institutions in other APAC jurisdictions read it as illustrative.

### AQ-4. APRA CPS 234 and CPS 230 alignment for Australian institutions

**Question.** Australian Prudential Regulation Authority (APRA) operates two foundational prudential standards: CPS 234 (Information Security; effective July 2019, with periodic enhancement) and CPS 230 (Operational Risk Management; effective July 2025, replacing CPS 231 + CPS 232 + CPS 233). CPS 234 requires APRA-regulated entities to maintain information security capability commensurate with the size and extent of threats; CPS 230 requires comprehensive operational risk management including third-party arrangements and business continuity. The chain's controls (HSM custody §10.5, append-only §4.1, reconciliation §10.1) directly support CPS 234 expectations (information-asset confidentiality/integrity/availability) and CPS 230 expectations (operational-risk identification, third-party-arrangement controls, business-continuity testing). Australian institutions need a mapping document.

**Status:** Partial.

**Closing language proposal.** A new supplementary document `docs/apac-regulator-pack/apra-cps-234-230-mapping.md` naming each CPS 234 / CPS 230 paragraph and the chain prescription that fulfills it. Australian institutions adopting the chain cite this document in their APRA Prudential Standards Compliance Self-Assessment.

### AQ-5. RBI guidelines on AI in financial services + Indian DPDP Act 2023 cross-border restrictions

**Question.** Reserve Bank of India (RBI) has issued multiple guidelines around AI in financial services: the RBI Master Direction on Outsourcing of Financial Services (updated 2024), the RBI guidelines on Cyber Security Framework (2016, with periodic updates), and the November 2023 letter to banks on AI/ML governance expectations. Indian Digital Personal Data Protection Act 2023 introduces cross-border-transfer restrictions on personal data: the central government will notify "trusted countries" to which transfer is permitted; transfer to non-notified countries requires standard contractual clauses or adequacy. The chain's cross-border data implications are significant for Indian institutions: an Indian bank operating chain artifacts in a non-Indian region (for backup, replication, or KMS access) must comply with DPDP cross-border restrictions. §10.15 Pattern B (per-region tenant_id) addresses the data-locality posture but DPDP-specific guidance is missing. RBI cyber-security guidelines also impose specific encryption-key-management expectations (RBI master circular on Customer Service in Banks, plus the cyber-security framework) that the chain's §10.5 HSM custody addresses but doesn't map specifically.

**Status:** Partial.

**Closing language proposal.** A new supplementary document `docs/apac-regulator-pack/rbi-india-mapping.md` naming each RBI guideline and the chain prescription that fulfills it, plus an explicit DPDP Act 2023 cross-border posture: Indian institutions deploy Pattern B with `bank-in` tenant pinned to AWS / Azure / Google APAC India region (Mumbai, Delhi, Hyderabad), with chain artifacts not crossing the Indian border. Institutions requiring cross-border processing (e.g., for cross-border-correspondent-banking AI) operate under DPDP's standard contractual clauses with the destination country.

### AQ-6. APAC cloud HSM regional listings — explicit guidance

**Question.** §10.5 names cloud HSM products (AWS CloudHSM, Azure Managed HSM, Google Cloud HSM, Thales Luna, Entrust nShield) without naming APAC regional availability explicitly. APAC institutions deploying must determine APAC region availability themselves: AWS CloudHSM (`ap-northeast-1` Tokyo, `ap-northeast-2` Seoul, `ap-northeast-3` Osaka, `ap-southeast-1` Singapore, `ap-southeast-2` Sydney, `ap-southeast-3` Jakarta, `ap-south-1` Mumbai, `ap-east-1` Hong Kong); Azure Managed HSM (Japan East/West, Korea Central, Southeast Asia, East Asia, Australia East/Central/Southeast, India Central/West/South); Google Cloud HSM (`asia-east1`, `asia-southeast1`, `asia-northeast1-3`, `asia-south1-2`, `australia-southeast1-2`). The mapping is research each institution does independently, producing varied posture. APAC supervisors increasingly require local-region data residency; explicit listing helps institutions choose conformant regions.

**Status:** Nit.

**Closing language proposal.** §10.5 adds a paragraph: "**APAC region availability (informative).** APAC region availability of the cloud HSM products listed above includes: AWS CloudHSM (Tokyo `ap-northeast-1`, Seoul `ap-northeast-2`, Osaka `ap-northeast-3`, Singapore `ap-southeast-1`, Sydney `ap-southeast-2`, Jakarta `ap-southeast-3`, Mumbai `ap-south-1`, Hong Kong `ap-east-1`); Azure Managed HSM (Japan East/West, Korea Central, Southeast Asia, East Asia, Australia East/Central/Southeast, India Central/West/South); Google Cloud HSM (`asia-east1`, `asia-southeast1`, `asia-northeast1-3`, `asia-south1-2`, `australia-southeast1-2`). APAC institutions select the region matching their data-residency requirements; cross-border processing requires the institution's CC8.1 documentation of the cross-border posture under the relevant national law (Singapore PDPA, HK PDPO, Indian DPDP Act 2023, Australian Privacy Act 1988, Korean PIPA)."

### AQ-7. APAC supervisory architecture — single tenant_id under multi-jurisdiction supervision

**Question.** An APAC bank operating in Singapore + Hong Kong + Tokyo + Sydney + Mumbai faces five national supervisors plus the ASEAN Banking Integration Framework (ABIF) for cross-border ASEAN operations and the SEACEN (South-East Asian Central Banks) coordination forum. Each supervisor may demand chain access during examination — examination-data sovereignty rules typically require artifacts accessible from within the supervisor's jurisdiction. §10.15 Pattern B (per-region tenant_id) covers this when "regions" map to "jurisdictions," but APAC institutions read "region" as cloud-region rather than supervisory-jurisdiction. Without explicit APAC mapping, institutions may default to Pattern A (single seal-region pinning across APAC) and find that Singapore operations' chain artifacts are stored under a Hong Kong seal-region tenant — defensible technically, awkward when MAS asks for Singapore-jurisdiction-only chain artifacts and the institution has to extract the Singapore subset.

**Status:** Partial.

**Closing language proposal.** §10.15 Pattern B adds a paragraph (parallel to Henrik's proposed EU-jurisdiction posture): "**APAC-jurisdiction posture (informative for APAC institutions).** Pattern B (per-region tenant_id) is the recommended posture for APAC institutions where per-jurisdiction event isolation supports per-supervisor examination access. Mapping per-jurisdiction tenant_ids onto APAC national supervisors (e.g., `bank-sg` / `bank-hk` / `bank-au` / `bank-in` / `bank-kr` corresponding to MAS, HKMA, APRA, RBI, FSC respectively) aligns the chain with the APAC supervisory topology and supports examination-data-sovereignty rules under each national framework. Cross-jurisdiction correlation (for institutions reporting to a consolidated supervisor — e.g., Australian Tier-1 banks reporting consolidated under APRA, Singapore-headquartered regional banks reporting consolidated under MAS Group Supervision) reads cross-tenant correlation as institution-side evidence; the institution's CC8.1 names the cross-jurisdiction correlation procedure. Per-jurisdiction tenant_id is also compatible with the §10.5 HSM custody requirement when each jurisdiction's tenant is custody-pinned to a HSM in the corresponding APAC region (per AQ-6 above)."

### AQ-8. APAC incident-notification timing matrix

**Question.** APAC institutions face simultaneous incident-notification clocks under multiple regulations:

- **MAS Notice 644 (Banking Act §60)**: significant cyber incidents reported within 1 hour of awareness; followed by detailed report within 14 days
- **HKMA TM-G-1**: incidents reported as soon as practicable, with HKMA Cyber-Resilience-Assessment-Framework (CRAF) 2.0 introducing tiered timing
- **APRA CPS 234**: information security incidents reported "no later than 72 hours after the regulated entity becomes aware of the incident"
- **RBI Cyber Security Framework**: incidents reported within 6 hours of detection; specific operational-risk incidents have shorter 2-hour reporting under RBI Master Direction on IT Outsourcing
- **FSC Korea**: cyber incidents reported within 24 hours under Electronic Financial Supervisory Regulation
- Plus national personal-data breach notification clocks (Singapore PDPA: as soon as practicable, generally 72 hours; HK PDPO: not specified statutorily but operational expectations; Indian DPDP Act 2023: 72 hours; Australian Privacy Act: as soon as practicable, generally 30 days; Korean PIPA: 72 hours)

The spec's incident-response-playbook (referenced in the v1.0-final-amendment change-log) mentions FFIEC's 36-hour clock but doesn't provide an APAC regulatory timing matrix. APAC institutions writing their own runbooks duplicate this effort and risk inconsistent posture during an actual incident.

**Status:** Partial.

**Closing language proposal.** The supplementary `incident-response-playbook.md` (parallel to Henrik's EQ-6 EU timing matrix) adds a new section: "APAC regulatory incident-notification timing matrix" with a table showing each APAC regulation, its trigger event, its notification clock, its reporting authority, and the chain artifacts that support the notification. This supplementary document is an APAC-conformance overlay layered on the existing FFIEC-side playbook; it does not change the spec's normative core.

### AQ-9. APEC CBPR cross-border data-flow framework — chain artifact handling

**Question.** APAC ex-Japan institutions operating across multiple APEC member economies face the APEC Cross-Border Privacy Rules (CBPR) framework: a voluntary certification system administered by APEC Data Privacy Subgroup. APEC member economies include Singapore, Hong Kong, Australia, Korea, Mexico, Canada, USA, Japan, Taiwan, Philippines, the Philippines. The CBPR framework provides a transfer mechanism for personal data flows between participating economies, parallel to the EU's adequacy framework. Chain artifacts containing personal data (or pseudonymised personal data) flowing between APEC economies SHOULD compose with CBPR; the spec doesn't address. APEC institutions evaluating composition write their own analysis; inconsistent posture results.

**Status:** Nit.

**Closing language proposal.** §10.15 Pattern B (or the supplementary `apac-regulator-pack/`) adds a paragraph: "**APEC CBPR posture (informative for APEC institutions).** Chain artifacts flowing between APEC member economies operate under the APEC Cross-Border Privacy Rules framework when the institution maintains APEC CBPR certification. The chain's per-region tenant_id pattern (Pattern B) supports CBPR's transfer-mechanism evidence requirements: the institution's CC8.1 documentation includes the institution's CBPR certification status and the chain artifacts subject to cross-border flow. Institutions without CBPR certification operate under bilateral national-law transfer mechanisms (Singapore PDPA's notification-of-transfer requirement, HK PDPO's voluntary code, etc.). The spec is data-flow-mechanism-neutral; the institution's chosen mechanism is a CC8.1 attribute."

---

## Per-role roll-up

**For the spec working group.** Seven Partials and two Nits. None affect the cryptographic-integrity claim or the implementation contract. All affect the documentation overlay needed for APAC ex-Japan institutions to deploy without writing substantial conformance documents themselves. The pragmatic close-out path is supplementary `docs/apac-regulator-pack/` documents (parallel to the existing FFIEC regulator-pack and Henrik's proposed EU regulator-pack from drop #08) plus targeted informative paragraphs in the spec at §10.5 and §10.15.

**For institutions evaluating chain deployment in APAC ex-Japan.** APAC ex-Japan institutions deploying the chain operate against the v1.0-final-amendment normative core (which is sound and APAC-friendly) plus a written-by-the-institution APAC-conformance overlay covering MAS FEAT Principles, MAS Veritas composition, HKMA TM-G-1, APRA CPS 234 / CPS 230, RBI guidelines, APAC cloud HSM regional listings, multi-jurisdiction supervisory architecture, incident-notification timing, and APEC CBPR. The overlay is substantial — multiple person-months for a Tier-1 institution operating across multiple APAC jurisdictions — and would benefit from being authored once at the spec layer.

**For the FFIEC working group's outreach to APAC regulatory bodies.** The spec's substantive choices are APAC-conformance-friendly; the gap is documentation, not substance. A working-group outreach to MAS Technology Risk Supervision Department, HKMA Banking Supervision Department, APRA Cyber and Operational Risk Division, RBI Cyber Security and Information Technology Examination Department, and FSC Korea Financial Information Technology Supervision could establish a supplementary APAC regulator-pack jointly, accelerating APAC adoption while preserving FFIEC-anchored spec governance. The MAS / HKMA bilateral working relationship is particularly receptive to FFIEC-equivalent technical standards because the Singapore-Hong Kong financial corridor consumes US-financial-regulatory technical artifacts directly.

---

## Where I would prioritize

1. **AQ-1 (MAS FEAT Principles mapping).** Highest leverage for Singapore institutions; establishes the FEAT mapping pattern that informs the broader APAC principle-based regulatory mapping. One annex.

2. **AQ-7 (APAC multi-jurisdiction Pattern B mapping)** and **AQ-8 (APAC incident-notification timing matrix).** Operationally pressing for any APAC institution operating cross-border. Two paragraph additions plus one timing-matrix table.

3. **AQ-3, AQ-4, AQ-5 (HKMA TM-G-1, APRA CPS 234/230, RBI India mappings).** Three supplementary `apac-regulator-pack/` documents covering Singapore-Hong Kong-Australia-India institutional supervisory landscape.

4. **AQ-2 (MAS Veritas composition).** Singapore-specific but informative for the broader APAC principle-based regulatory pattern. One supplementary document.

5. **AQ-6 (APAC cloud HSM regional listings).** Editorial; one paragraph in §10.5.

6. **AQ-9 (APEC CBPR posture).** Forward-scope; informational.

---

## Stopping criterion

This drop carries 7 Partial + 2 Nit findings against the v1.0-final-amendment spec, all clustered around APAC ex-Japan regulatory-architecture conformance overlays. None affect the spec's cryptographic-integrity claim or its normative core; all affect the documentation overlay an APAC institution would otherwise write for itself. The substantive cryptographic-strength choices (HMAC-SHA-256, HKDF-SHA-256, Ed25519, RFC 8785, RFC 6962 Merkle, FIPS 140-2 Level 3 HSM custody) are NIST-aligned and APAC-conformance-friendly; the conformance gaps are at documentation and operational-procedure layers, not at the cryptographic-primitive layer.

Three of my Tier-1 APAC bank clients are in evaluation. Their reading is consistent: the technical primitives are at parity with or above what they would expect from a regional technical standard, but the APAC-conformance overlay is missing and they don't want to author it independently. A working-group cycle producing an APAC regulator-pack — informative annexes for MAS FEAT mapping, MAS Veritas composition, supplementary documents for HKMA / APRA / RBI mappings, plus targeted §10.5 / §10.15 APAC-region clarifications — would unblock pilot deployments. I would recommend this work proceed in parallel with the existing FFIEC regulator-pack and Henrik's proposed EU regulator-pack, rather than waiting for v1.x.

The chain's design has unusual APAC affinities — the test-vector corpus, the byte-reproducible verifier, the principle-style invariant articulation in §4.1 — that compose well with APAC supervisors' expectation patterns. Capitalising on these affinities through an APAC regulator-pack would expand the spec's adoption footprint substantially and produce a multi-region adoption pattern (US-FFIEC, EU, APAC) that no single jurisdiction's conformance work alone would achieve.
