# Margarethe Lindqvist — DORA / EBA First-Look Review of FFIEC chain-of-custody v1.0a
**Date:** 2026-05-07
**Reviewer:** Margarethe Lindqvist, Senior Supervisor (FI / EBA secondment)
**Posture:** First encounter with v1.0a; no prior iteration history. American specification reviewed with European supervisory lens.

## Initial supervisory impression

The specification is technically competent and unusually well-disciplined for an American sectoral standard. The three-layer integrity model (per-event HMAC, daily Merkle aggregation, HSM-bound Ed25519 over a v1.0a `sign_payload`) maps cleanly onto what the EBA expects an institution to demonstrate when an examiner asks "show me how you would prove the AI agent's record of decisions was not silently rewritten by your vendor or by your own DBA." The cryptographic substrate — FIPS 180-4, FIPS 186-5, FIPS 198-1, RFC 5869, RFC 6962, RFC 8785 — is the substrate I would accept without argument from an institution under SREP review. The verifier's offline, single-binary, no-network-calls posture is closer to what the European Court of Auditors expects of a third-party-verifiable artefact than what most regulated institutions currently produce.

The specification is, however, written from a US sectoral horizon. DORA Articles 5–28, the EBA Guidelines on outsourcing (EBA/GL/2019/02), the Schrems II supplementary-measures jurisprudence, the eIDAS Regulation (910/2014), TIBER-EU, and NIS2 are not referenced anywhere in the normative text. The regulator-pack documents (international-transfers, multi-jurisdiction-conflict, gdpr-controller-vs-processor, article-32-security-mapping, breach-notification-matrix) carry the GDPR and the Schrems II picture credibly, but DORA — which is now the supervisory framework binding every credit institution, payment institution, investment firm, and CSDR entity that runs an ICT third-party arrangement in the Union — appears only once, in the breach-notification matrix, as a 24-hour clock reference. That is insufficient. An EU supervisor reviewing this evidence today would accept the integrity claim, accept the GDPR mapping, and immediately ask the institution to produce a DORA-articulated overlay before any joint examination team would close the file.

The specification is therefore close — but not yet sufficient on its own — for a competent-authority dialogue under the European framework. The integrity primitives carry their weight; the institutional-articulation layer between the primitives and the European supervisory expectations is where the work remains.

## Findings

### Gaps (DORA/EBA expectations the spec does not address)

**G-1. No DORA Article 28-30 controller / processor / sub-processor positioning of the SDK provider, the receiver provider, and the LLM provider.**
**Section/file:** `spec/chain-of-custody-v1.md` §4 (primitives), §10 (security overview); `docs/regulator-pack/gdpr-controller-vs-processor.md`.
**EU framework reference:** DORA Articles 28-30 (ICT third-party risk management, contractual provisions); EBA/GL/2019/02 §§4.12-4.13 (sub-outsourcing).
**Issue:** The GDPR processor analysis is competent, but DORA imposes a parallel and stricter regime. An ICT third-party service provider supplying integrity-bearing audit infrastructure to an EU credit institution is in scope of Article 28(2) and the contractual register under Article 28(3). The LLM provider — OpenAI, Anthropic, Google — is a sub-contractor in the chain (Article 30(2)(a)) and the institution must articulate this in its register. The specification does not name the LLM provider's position in the lattice at all, despite this being the central novelty of the workload.
**What an EU supervisor would expect to see:** A DORA-overlay document mapping each role (SDK author, ledger operator, HSM operator, LLM provider) onto Article 28-30 categories; a contractual-clause checklist citing Article 30(2)(a) through (h) (description of services, location of processing, sub-contracting permissions, exit strategy, audit rights, monitoring rights, incident-cooperation, termination rights); explicit articulation that the LLM provider is a "supporting ICT services" provider in the meaning of Article 3(21) and is therefore in scope of Article 30 contractual provisions even when accessed by API.

**G-2. No DORA Article 29 ICT concentration-risk articulation when a single receiver provider serves N institutions across the EU.**
**Section/file:** `docs/design/00-overview.md` §6.3 (vendor-hosted topology); `spec/chain-of-custody-v1.md` §10.15.
**EU framework reference:** DORA Article 29 (concentration risk at the level of ICT third-party service providers); EBA/GL/2019/02 §4.7.
**Issue:** The vendor-hosted topology contemplates a multi-tenant ledger and HSM in a vendor's cloud. If that vendor accumulates N EU-supervised institutions on the same infrastructure, the failure of that vendor becomes a Union-level systemic event. The specification does not require the vendor to publish or the institution to track concentration metrics, nor does it acknowledge that a "critical ICT third-party service provider" designation under DORA Article 31 would change the supervisory regime materially.
**What an EU supervisor would expect to see:** A normative requirement (or at minimum a strong recommendation) that vendor-hosted deployments publish per-institution-segregation evidence; a documented expectation that institutions running on a single vendor coordinate exit-strategy testing under Article 28(8); explicit acknowledgement that the spec's "vendor-hosted multi-tenant" topology may trigger Article 31 designation if the vendor's market share crosses thresholds.

**G-3. No DORA Article 18 incident-classification fields — impact, duration, geographic spread, data-loss criteria.**
**Section/file:** `spec/chain-of-custody-v1.md` §4.3.1 (HSM unavailability), §10.10 (rotation); `docs/incident-response-playbook.md` (referenced).
**EU framework reference:** DORA Article 18 and the RTS on incident classification (Commission Delegated Regulation (EU) 2024/1772).
**Issue:** The specification produces excellent integrity evidence but does not produce the structured fields the RTS requires — geographic spread (per-Member-State affected counts), duration (in continuous hours), criticality of services affected, data-loss extent, economic impact bands. The institution will therefore have to construct the DORA Article 18 classification dossier from external sources at incident time.
**What an EU supervisor would expect to see:** An operational-event extension naming `incident.classification.*` attributes mapping to the RTS criteria — affected-jurisdictions list, continuous-duration counter, criticality enum (per the institution's CRITICAL/HIGH/MEDIUM/LOW classification), data-loss-extent quantification, recovery-objective state. Alternatively, an explicit cross-reference document showing how the chain's `master.*` operational events feed the RTS classification template.

**G-4. No DORA Articles 11-14 incident-reporting timeline articulation (4-hour initial / 72-hour intermediate / 1-month final).**
**Section/file:** `docs/regulator-pack/breach-notification-matrix.md` §1.
**EU framework reference:** DORA Articles 11-14; Commission Delegated Regulation 2024/1774 RTS on reporting; Implementing Regulation 2024/2956 templates.
**Issue:** The matrix mentions DORA at "24 hours initial / 72 hours intermediate / 1 month final" but the RTS finalised in 2024 sets the initial notification at **4 hours** from major-incident classification, not 24. The 24-hour figure is a legacy from the proposal stage. An EU supervisor reading "DORA gives 24 hours" in the institution's own evidence pack would treat the document as not current.
**What an EU supervisor would expect to see:** Correction to 4-hour initial notification per the final RTS; explicit reference to Implementing Regulation 2024/2956 templates; an institution-side runbook ensuring the chain's evidence can be assembled within the 4-hour window for the initial notification's "categorisation criteria met" question.

**G-5. No DORA Articles 26-27 TLPT (threat-led penetration testing) recoverability of red-team activity as audit-events without leaking attack TTPs.**
**Section/file:** `docs/design/09-threat-model.md` §2 (adversaries); spec §10 (security overview).
**EU framework reference:** DORA Articles 26-27; TIBER-EU framework; CBEST methodology.
**Issue:** TLPT is mandatory for significant institutions every three years. The chain captures AI agent activity but the spec does not address how a legitimate red-team engagement is recorded — the test team's prompt-injection attempts, their tool-call probing, their attempts to subvert the agent's policy guardrails — without those captures becoming an attacker's playbook in the ledger. A red-team engagement against an AI agent will produce thousands of chain entries that document attack patterns; if those entries are retained at the same retention as production, they become an information-leak path.
**What an EU supervisor would expect to see:** A TLPT mode where chain entries are tagged `audit.tlpt.engagement_id` and given a separate retention bucket; explicit articulation of how the TLPT lead and the institution's blue team can both verify integrity without the institution's broader operations team gaining read access; alignment with TIBER-EU §3.4 (test-data handling) and the RTS on TLPT under Article 26(11).

**G-6. No eIDAS qualification of the Ed25519 signature; the spec does not assert whether the seal signature is simple, advanced electronic, or qualified.**
**Section/file:** `spec/chain-of-custody-v1.md` §4.3 (HSM signature); §10.5 (HSM tier).
**EU framework reference:** eIDAS Regulation (EU) 910/2014 Articles 25-29; eIDAS 2.0 (Regulation (EU) 2024/1183).
**Issue:** Under eIDAS Article 25, a qualified electronic signature has the equivalent legal effect of a handwritten signature throughout the Union. An advanced electronic signature retains evidentiary value but a court may apply discretion. The seal signature in the chain is produced by an HSM at FIPS 140-2 Level 3 (or higher) and binds an Ed25519 key uniquely to the institution. That meets the four Article 26 criteria for an advanced electronic signature on its face, but eIDAS qualification additionally requires the signing device to be a qualified signature creation device (QSCD) on the EU Trust List and the certificate to be issued by a qualified trust service provider. The spec is silent on this.
**What an EU supervisor would expect to see:** Explicit declaration that the seal signature is an advanced electronic signature within the meaning of eIDAS Article 26 unless the HSM is on the EU Trust List of QSCDs and the public-key certificate is issued by a QTSP, in which case the qualified-electronic-signature regime applies; institution-side guidance on which posture they have selected and why.

**G-7. No NIS2 transposition awareness for essential entities operating the chain.**
**Section/file:** `docs/regulator-pack/breach-notification-matrix.md`; `docs/incident-response-playbook.md`.
**EU framework reference:** Directive (EU) 2022/2555 (NIS2), transposed into national law in 2024; Commission Implementing Regulation 2024/2690.
**Issue:** Many institutions running this chain in the EU are also "essential entities" under NIS2, which adds an additional incident-notification clock — early warning within 24 hours, incident notification within 72 hours, final report within one month — to the national CSIRT. The matrix does not list NIS2.
**What an EU supervisor would expect to see:** NIS2 added to the four-clock matrix; reconciliation rules between DORA Article 18 reporting and NIS2 reporting (DORA generally takes precedence for financial entities under Article 1(2) lex specialis, but the institution must document this).

**G-8. No EBA Guidelines on outsourcing (EBA/GL/2019/02) audit-rights clause for competent authorities.**
**Section/file:** `docs/regulator-pack/gdpr-controller-vs-processor.md` §"Audit rights"; `docs/design/00-overview.md` §6.
**EU framework reference:** EBA/GL/2019/02 §13 (rights of access and audit); §14 (security of data and systems).
**Issue:** The DPA template names "audit rights" generically. EBA/GL/2019/02 §13.2 specifically requires the contract to grant the competent authority and any auditor appointed by it unrestricted access rights to the service provider's premises and to all relevant business premises, devices, systems, networks, information, and data. The receiver-policy discovery endpoint (referenced at spec §4 and explored in the design overview) is the obvious place to surface this for an examiner. A vendor-hosted ledger that does not contractually grant the EBA-specified access right to the institution's competent authority is non-conformant under the Guidelines.
**What an EU supervisor would expect to see:** An explicit clause in the DPA template (or its DORA-overlay counterpart) granting the institution's competent authority and any appointed auditor the EBA/GL/2019/02 §13.2 access right; institution-side reception procedure for examiner inspection visits to vendor premises.

**G-9. No TARGET2/T2S analogy or ECB SREP framing on what happens when a chain that documents a clearing-relevant decision is corrupted.**
**Section/file:** `docs/dr-and-resilience.md`; `docs/design/09-threat-model.md` §3.1.1.
**EU framework reference:** ECB SREP guide §5.5 (operational risk, ICT risk); SSM expectations on operational resilience for clearing-relevant systems.
**Issue:** The chain's denial-of-service-against-integrity scenarios (backup tampering, WAL corruption, HSM unavailability) are documented at the institution level. They are not framed in terms of what an SSM joint supervisory team would view as the systemic-risk dimension — an integrity failure on a chain documenting AI-driven decisions in correspondent-banking, payment-screening, or KYC would propagate to TARGET2 reconciliation and would attract ECB attention well above the institution's primary supervisor.
**What an EU supervisor would expect to see:** An ICAAP/ILAAP residual-risk articulation that names chain-integrity failures as an operational-risk scenario in the institution's Pillar 2 capital allocation; cross-reference to the SSM's expectations on operational resilience.

**G-10. No data-residency-by-default in the spec for EU operations.**
**Section/file:** `spec/chain-of-custody-v1.md` §6 (audit-file format), §10.15 (multi-region).
**EU framework reference:** GDPR Articles 44-50; Schrems II (CJEU C-311/18); EDPB Recommendations 01/2020.
**Issue:** Pattern A (active-active with seal-region pinning) does not normate that the seal region for an EU tenant must be in the EEA. Pattern B is described as "appropriate when EU banking jurisdictions mandate in-region key custody" but is not the default. An EU-based controller deploying under Pattern A with a US-based seal region would have a transfer event that triggers the full Schrems II supplementary-measures regime, and that combination is the riskiest deployment shape. The spec presents Pattern A as RECOMMENDED, which an EU supervisor would read as the spec author's preference — and that preference is misaligned for EU controllers.
**What an EU supervisor would expect to see:** An explicit normative statement that for EU controllers, the seal region MUST be in an adequate jurisdiction (EEA, UK, or other adequacy decision) absent a documented Article 49 derogation or an EDPB-approved supplementary-measure architecture; Pattern B as the RECOMMENDED default for EU tenants whose risk posture treats US-cloud-provider access as an unacceptable trust boundary.

**G-11. No articulation of "right to explanation" interaction with the chain under GDPR Article 22 / EU AI Act Article 86.**
**Section/file:** `docs/litigation-support.md` §5.2; `docs/regulator-pack/article-32-security-mapping.md`.
**EU framework reference:** GDPR Article 22; EU AI Act (Regulation (EU) 2024/1689) Article 86 (right to explanation of individual decision-making).
**Issue:** The chain proves what the AI said. Article 22 and AI Act Article 86 give a data subject the right to obtain an explanation of an automated decision affecting them. The chain is the obvious source of that explanation — and the institution's ability to surface a chain-derived explanation to the data subject within reasonable time is itself an Article 22 / Article 86 control. The spec does not address how the chain composes with the data-subject explanation right.
**What an EU supervisor would expect to see:** A `regulator-pack/ai-act-article-86.md` companion document showing how a chain entry resolves into an Article 86 explanation, what the explanation must contain (the role of the AI in the decision, the main parameters, and the consequences), and what is redacted before disclosure to the data subject.

### Partials (addressed but with gaps for European application)

**P-1. Schrems II supplementary measures are well-articulated but EU-held-keys is "preferred," not normative.**
**Section/file:** `docs/regulator-pack/international-transfers.md` §2.2.
**EU framework reference:** Schrems II; EDPB Recommendations 01/2020 §85.
**Issue:** EU-held keys should be the floor for any vendor-hosted deployment serving an EU controller, not the preference.
**What would close it:** Promote EU-held keys from "preferred" to a normative MUST absent a documented EDPB-approved alternative.

**P-2. The GDPR Article 28 DPA template is complete; the DORA Article 30(2) overlay is missing.**
**Section/file:** `docs/regulator-pack/gdpr-controller-vs-processor.md`.
**Issue:** Article 28(3) elements are covered. DORA Article 30(2) adds full service description, processing locations, exit strategies and transition periods, service-performance monitoring, ICT-security awareness obligations, sub-contracting chain notification, termination for supervisory breach, and TLPT participation rights.
**What would close it:** A DORA-Article-30(2) supplement to the DPA template, or a parallel ICT-services-agreement template aligned with the EBA outsourcing register.

**P-3. The 36-hour FFIEC cyber-incident clock is referenced without DORA lex-specialis reconciliation.**
**Section/file:** `spec/chain-of-custody-v1.md` §4.3.1; `docs/regulator-pack/breach-notification-matrix.md`.
**Issue:** DORA Article 1(2) makes DORA the lex specialis for Union financial entities. For dual-supervised institutions, the precedence rule is not articulated.
**What would close it:** A precedence-rules table covering FFIEC 36h vs DORA 4h vs GDPR 72h vs HIPAA 60d vs NIS2 24h+72h; the IR commander operates the tightest applicable clock first.

**P-4. FIPS 140-2 Level 3 is named; the EU CC equivalence regime is asserted but not defended.**
**Section/file:** `spec/chain-of-custody-v1.md` §10.5.
**Issue:** EU institutions often hold HSMs evaluated under BSI AIS 31 or ANSSI CSPN. The FIPS / CC EAL4+ equivalence is asserted without detail.
**What would close it:** An equivalence table referencing EN 419 221-5 (server-signing QSCDs) and EN 419 211 (SSCDs).

**P-5. Pattern B is conformant; the EU jurisdictions that mandate it are not named.**
**Section/file:** `docs/design/00-overview.md` §6.4; spec §10.15.
**Issue:** "Some EU banking jurisdictions" is too vague. BaFin BAIT, ACPR's outsourcing notice, and Banca d'Italia Circular 285 have explicit articulation worth citing.
**What would close it:** An informative annex citing the specific national prudential frameworks; cross-reference to the EBA's cloud-outsourcing recommendations §4.6.

**P-6. The receiver-policy discovery endpoint is "implementation-specific."**
**Section/file:** `spec/chain-of-custody-v1.md` §4.
**Issue:** The endpoint is the natural surface for EBA/GL/2019/02 §13 audit-rights and for an examiner inspection visit. A per-vendor protocol forces the examiner to learn each implementation.
**What would close it:** A normative minimum surface so an examiner performs a baseline competent-authority inspection without per-vendor training.

### Nits (clarification asks)

**N-1. `mac_computed_at_utc` is wall-clock and forensic-only.** Under eIDAS Article 41, a qualified electronic time stamp shifts the burden of proof on date and time. The current posture is fine but a note acknowledging the upgrade path to qualified time stamps would be useful.

**N-2. `kms_handle_uri` examples are AWS-canonical.** Add `azure-managed-hsm:`, `gcp-cloud-hsm:`, and one on-prem `pkcs11://` form so the example set reads as cross-cloud.

**N-3. OTLP transport binding.** Some EU operations centres mandate FAPI 2.0 + mTLS for inter-zone transport. Acknowledging that OTLP can be wrapped under FAPI 2.0 mTLS helps institutions whose CSIRT requires it.

**N-4. Empty-day Merkle root collision.** A note that `SHA-256(b"")` is identical across tenants on empty days but discriminated by `tenant_id` and `seal_date` in the seal record would short-circuit a predictable supervisor question.

**N-5. Daubert framing.** Daubert is US-specific. For a European court, the analogue is eIDAS Articles 35-37 and the national rules transposing the e-Evidence Regulation (EU) 2023/1543. A short EU-evidence parallel makes the document portable.

### Confirmations (decisions I would defend in an SREP review)

**C-1. The HMAC + Merkle + Ed25519 + HSM composition.** Three-layer integrity with separation-of-duties between IKM custody, ledger storage, and HSM signing is a textbook implementation of EBA/GL/2019/02 §14 and Article 32 GDPR. The verifier feeding `expected_prev_hash` from the structural walk — not the entry's claimed `prev_hash` — into the MAC recompute closes a footgun I have seen in two prior audit-trail designs.

**C-2. The verifier's no-network-calls discipline.** A verifier that makes no outbound calls cannot be steered toward a forged ledger or a malicious public key at examination time. This is what the European Court of Auditors expects of a third-party verification artefact and is the reason the verifier is admissible as direct evidence.

**C-3. The HSM tier requirement.** FIPS 140-2 L3 / CC EAL4+ matches EBA expectations for cryptographic devices in regulated finance. Software-key fallback excluded from production builds at compile time is stronger than a runtime configuration check.

**C-4. The cross-tenant binding via HKDF info parameter.** Two tenants whose IKMs are accidentally swapped derive verifiably different session keys; the per-entry `key_fingerprint` check happens before any MAC compute. The answer to cross-tenant confusion is mathematical, not procedural.

**C-5. The v1.0a `sign_payload` algorithm-binding.** Binding `sign_payload_version` and `algorithm` into the signed payload closes the algorithm-confusion attack class at v1.0 without waiting on a v1.x amendment. The right posture for the dual-algorithm transitional period when Dilithium / SLH-DSA ship in EU-jurisdictional HSMs.

**C-6. The deterministic-build, twice-signed verifier supply chain.** Three independent compromises (cosign, GPG manifest, regulator-held fingerprint) would have to align for silent failure. This is the supply-chain floor EBA/GL/2019/02 §14.4 asks for.

**C-7. Multi-region resilience normated rather than left to implementations.** Pattern A and Pattern B normated, with per-tenant pattern switching forbidden absent a documented chain-discontinuity event, is the operational discipline DORA Article 12 (ICT business continuity policy) expects.

**C-8. The wire-or-on-disk observation rule.** Integrity claims must be cited from a persisted-or-wire artefact, not from in-process state. The correct boundary between operations tooling and supervisory observation.

## Bottom line

The specification is technically sufficient and would be accepted by an EU competent authority as the integrity layer in an ICT third-party assessment under DORA Article 28-30 — provided the institution layers a DORA-articulated overlay on top. The absent overlay is the work of weeks, not months: the integrity primitives are the hard part and they are correct. The DORA-mapping document, the eIDAS qualification declaration, the NIS2 reconciliation in the breach-notification matrix, the corrected DORA 4-hour clock, and the EU-held-keys normative posture for vendor-hosted EU controllers are the principal items that move the file from "credible American sectoral standard" to "directly referenceable in a Union joint examination team's working file."

Subject to the G-1 through G-11 remediation, I would accept this evidence in an SREP dialogue under Article 28 ICT third-party risk; without that remediation, I would request the institution produce a complementary DORA / eIDAS / NIS2 articulation before closing the assessment. The cryptographic substrate is not the obstacle; the supervisory translation layer is.
