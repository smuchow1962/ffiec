# Outside reviewer — US federal civilian agency CISO

**Reviewer.** Marcus Patel-Williams. Deputy Chief Information Security Officer at a large US federal civilian agency operating high-volume citizen-facing decisions at a scale comparable to IRS, SSA, CMS, VA, and Department of Education. Reports to the agency CISO with direct authority over agency-wide AI security posture and Authorization to Operate (ATO) decisions for AI systems. BS Computer Science, MS Cybersecurity Policy and Risk Management. Started at NSA's Information Assurance Directorate (now Cybersecurity Directorate), moved to a civilian agency CIO office, eight years federal CISO leadership. Holds CISSP, CISM, CGRC (formerly CAP), and the FedRAMP-certified 3PAO Independent Assessor designation. Member of the Federal CISO Council. Has issued ATOs for cloud-deployed AI under FedRAMP Moderate and High. Currently evaluating chain-of-custody infrastructure for (a) the agency's AI-driven citizen-facing decisions, (b) internal IT operational AI (tier-1 helpdesk, SOC analyst augmentation), and (c) potential expansion to defense and Intelligence Community partner sharing under additional CMMC and IC-CD authorization processes.

**Why this review matters.** Marcus's agency operates under FISMA + NIST RMF (SP 800-37 Rev 2), with cloud workloads under FedRAMP and AI workloads under multiple OMB memoranda — M-22-09 (zero-trust), M-23-22 (agency AI use), M-24-04 (zero-trust progress), M-24-10 (AI use-case inventories), M-24-15 (AI governance roles). The chain is proposed as the agency's AI-decision audit-trail substrate. Marcus's ATO sign-off is required before any AI deployment for citizen-facing decisions can proceed. For defense and IC partner-sharing, additional CMMC 2.0 (DoD) and IC-CD reviews follow. The questions below identify the gaps Marcus would mark in his ATO package and in the System Security Plan (SSP) submitted to the Authorizing Official (AO).

**Reading angle.** US federal information-security regulatory landscape — FISMA + NIST RMF (SP 800-37 Rev 2), FedRAMP Rev 5 baselines, NIST SP 800-53 Rev 5 controls, SP 800-171 Rev 3 + SP 800-172 for CUI, NIST AI RMF 1.0 + the Generative AI Profile (AI 600-1), the OMB AI memoranda, NIST CSF 2.0, CISA Binding Operational Directives (BOD 22-01, BOD 23-01), SP 800-218 / 218A SSDF, CMMC 2.0 for defense / IC, CNSSP-15 / CNSA Suite for National Security Systems, Cross-Domain Solutions for IC partner sharing, OMB Circular A-130, NARA records-retention. The OMB memos persist regardless of the EO 14110 rescission.

**Distinct from existing reviewers.** Reuven (drop 04) covers formal cryptographic security; the substrate is sound. Sophia (drop 10, FDIC/OCC) covers FFIEC bank-examiner consumption from the regulatory-supervisor side. Marcus Brennan (drop 10, SOC 2 / ISAE 3000) covers AICPA TSC mapping for commercial attestation — FedRAMP uses SP 800-53 Rev 5, not TSCs, so that mapping doesn't satisfy a FedRAMP SSP or SAP. Henrik (drop 08) covers EU regulatory, not US federal authorization. My threshold is an ATO-authorization gap, an SP 800-53 mapping gap, or an OMB / AI RMF alignment gap that forces the agency's authorization team to write substantial substrate themselves before the SSP can be assembled. I stay clear of cryptographic evaluation, attestation engagement procedure, and bank-examiner consumption.

**Stopping criterion declared up front.** 13 findings clustered around (a) NIST SP 800-53 Rev 5 control mapping precision, (b) OMB AI memorandum vocabulary alignment, (c) CUI handling under SP 800-171 Rev 3, (d) CMMC 2.0 and CNSSP-15 / CNSA Suite scoping for defense / IC use cases, (e) zero-trust architecture composition under M-22-09 and M-24-04, (f) NARA records-retention scheduling, (g) CISA BOD vulnerability-management alignment, and (h) NIST AI RMF subcategory mapping for the agency's AI Risk Management Plan. The cryptographic substrate is sound; this drop is about the authorization-package surfaces a federal agency assembles before deployment.

---

## Headline

The chain's design is unusually clean for federal authorization: byte-reproducible verifier output (§7), prescribed FIPS 140-2 Level 3 HSM custody (§10.5), explicit FIPS-validated RNG provenance (§10.6.1), an exit-code contract for the verifier CLI (§10.12), and separation-of-duties language at the seal-job operator level. The cryptographic substrate is FIPS-current — SHA-256 (FIPS 180-4), HMAC (FIPS 198-1), Ed25519 (FIPS 186-5), HKDF (RFC 5869) all sit on FedRAMP-acceptable algorithm lists. **The substrate is sound. The gaps cluster at the authorization-package layer — the surfaces an agency assembles before issuing the ATO.** None of the findings below affect the integrity claim; all affect the authorization-package work each adopting agency would otherwise write from scratch on its own ATO timeline.

**What the spec does well for federal authorization.** The cryptographic primitives map cleanly to SP 800-53 Rev 5 control families: HSM custody (§10.5) → SC-12, SC-13, SC-28; append-only (§4.1 + §10.3) → SI-7, AU-9; reconciliation (§10.1) → CA-7, AU-12; verifier (§7) gives ConMon a deterministic, replayable signal under AU-6; §10.12 exit-code contract is exactly the disposition signal ConMon tooling consumes; §10.13 retention → AU-11, SI-12. The chain is broadly FedRAMP-friendly at the technical layer. The §1.4 compositional-security argument — three independent custody layers under separation of duties — is precisely the defense-in-depth posture an AO looks for when authorizing citizen-facing decisions.

**What the spec needs for the federal-authorization layer.** An SP 800-53 Rev 5 control-mapping table for FedRAMP Moderate and High. M-23-22 vocabulary alignment. CUI marking and handling under SP 800-171 Rev 3 for chain artifacts that touch non-federal systems. A CMMC 2.0 scoping note for defense / IC. A CNSSP-15 / CNSA Suite acknowledgment that v1.0 Ed25519 + SHA-256 is not currently CNSA Suite and an NSS-conformant variant is a v1.x scope item. A zero-trust composition note covering identity (PIV/CAC), device trust, workload trust, and data classification. A NIST AI RMF subcategory mapping. A CISA BOD vulnerability-management note. A NARA records-retention alignment note. None integrity-bearing; all authorization-package quality items.

Each finding below carries a Status (Partial, Nit, or Gap) and a concrete closing-language proposal. The proposals are the kind that land cleanly as informative annexes or supplementary documents in the regulator pack — they do not require changes to the spec's normative core.

---

## Supporting questions

### FQ-1. NIST SP 800-53 Rev 5 control-mapping table — FedRAMP baselines

**Question.** FedRAMP authorization (Joint Authorization Board or Agency authorization under FedRAMP Rev 5 baselines) requires the System Security Plan (SSP) to map each implemented control to the baseline NIST SP 800-53 Rev 5 control identifier. The chain prescribes controls — HSM custody (§10.5), append-only enforcement (§4.1 + §10.3), reconciliation (§10.1), software-key adapter exclusion (§10.7), constant-time comparison (§10.8), IKM retention (§10.9), evidentiary-artifact retention (§10.13), algorithm rotation (§10.10) — but does not map them to NIST SP 800-53 controls explicitly. Each adopting agency must do its own mapping for the SSP, and without a baseline in the spec the result is variation in precision across agencies and inconsistent cross-agency reciprocity (FedRAMP authorizations are designed to be reusable across agencies — reuse depends on consistent SSP control mapping). The mapping I would expect to see at minimum:

- §10.5 HSM custody → SC-12, SC-13, SC-28 / SC-28(1), PE-3 where on-prem HSM applies.
- §10.6.1 RNG → SC-12 / SC-12(2) / SC-12(3), SC-13.
- §10.7 software-key adapter exclusion → CM-7 / CM-7(1), CM-2, SI-7.
- §4.1 inviolate + §10.3 append-only → SI-7 / SI-7(1), AU-9 / AU-9(2) / AU-9(3).
- §10.1 reconciliation → CA-7 / CA-7(1), AU-12 / AU-12(1).
- §7 verifier + §10.12 exit codes → AU-6 / AU-6(1) / AU-6(3).
- §10.13 retention → AU-11, SI-12.
- §10.10 algorithm rotation → CM-3 / CM-3(1), SC-12, SC-13.
- §10.4 time sync → AU-8 / AU-8(1).
- §4.4 OTLP + §5.1 TLS 1.3 → SC-8 / SC-8(1), SC-23.
- §10.15 multi-region → CP-7, CP-9, SC-5 for seal-region failover.

Each control has organization-defined values (ODVs) the agency parameterizes per risk tier (AU-11 retention period being a clear example); §10.13 + §10.9 give a cross-walk but the ODV interaction is not explicit.

**Status.** Partial. The cryptographic substrate is FedRAMP-friendly; the mapping table is missing.

**Closing language proposal.** A new informative annex `regulator-pack/nist-800-53-rev5-mapping.md` listing each chain-prescribed control alongside its NIST SP 800-53 Rev 5 control identifier(s) for FedRAMP Moderate and High baselines, with a column for the organization-defined values (ODVs) each agency parameterizes (retention period, scan frequency). Pair with a worked SSP-fragment example showing how a Moderate-baseline agency would cite the chain in its SC-13 control-implementation description.

### FQ-2. OMB M-23-22 alignment — "automated logging of inputs, outputs, and material changes"

**Question.** OMB Memorandum M-23-22 ("Advancing Governance, Innovation, and Risk Management for Agency Use of Artificial Intelligence," 28 March 2024) imposes specific minimum risk-management practices on federal agencies for "rights-impacting" and "safety-impacting" AI use cases. Among them: an obligation to "automatically log AI system inputs, outputs, and material changes." The chain literally captures inputs (prompt), outputs (model response), tool calls, routing decisions, and operational state — the M-23-22 obligation is precisely what the chain delivers. But the spec writes in FFIEC vocabulary and does not cross-walk to the M-23-22 vocabulary. An agency CIO assembling its M-23-22 compliance package — required for every rights-impacting and safety-impacting AI use case — has to translate manually. M-23-22 has also been operationalized into the Annual AI Use Case Inventory (M-24-10) and the AI Governance Board's risk register (M-24-15); both need to cite the chain by what M-23-22 calls it.

The same gap applies to the EO 14110 transition. The Biden EO has been rescinded under the new administration, but the OMB memos (M-22-09, M-23-22, M-24-04, M-24-10, M-24-15) persist because they predate or stand independent of the EO. Federal AI policy is volatile; the chain should frame itself as agency-policy-neutral infrastructure satisfying the underlying minimum-practice obligations regardless of which administration's AI EO is in force.

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/omb-m-23-22-alignment.md` mapping each M-23-22 minimum risk-management practice to the chain attribute(s) and operational event(s) that satisfy it. At minimum: M-23-22 "automated logging of inputs, outputs, and material changes" maps to the per-event canonical bytes (captured prompt, model response, tool calls per `chain_kind = "model_call" | "tool_call"`, routing decisions per §4.4.1, master-key rotation events per §10.2). Add parallel mappings for M-24-10 (Annual AI Use Case Inventory fields populated from chain attributes) and M-24-15 (AI Risk Management Plan substrate citing the verifier output as operational-effectiveness evidence). Frame the annex as agency-policy-neutral so it survives administration changes.

### FQ-3. NIST AI RMF 1.0 subcategory mapping for the agency AI Risk Management Plan

**Question.** NIST AI Risk Management Framework 1.0 (NIST AI 100-1, January 2023) and the Generative AI Profile (NIST AI 600-1, July 2024) define functions and subcategories agencies are encouraged to align with — Govern, Map, Measure, Manage. The chain provides the technical substrate for several subcategories explicitly: GOVERN 1.4 (legal and regulatory requirements satisfied), MEASURE 2.7 (AI system security and resilience monitored), MANAGE 4.1 (post-deployment monitoring), and others. M-24-15 names the AI RMF as a reference framework agencies use to structure their AI Risk Management Plan; an agency citing the chain in its AI RMF subcategory implementations needs the mapping to be defensible at the AI Governance Board's review. The chain's spec does not currently include this mapping; the agency's CAIO (Chief AI Officer, designated under M-24-15) writes one for the agency's plan, and each agency writes its own — variation across agencies undermines cross-agency reciprocity.

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/nist-ai-rmf-mapping.md` listing each NIST AI RMF 1.0 function and subcategory and the chain primitive, attribute, or operational event that supports it. Pair with a similar table for the NIST AI 600-1 Generative AI Profile, where the chain's GenAI attributes (`gen_ai.request.model`, `gen_ai.response.model`, the §7 step 12a completeness check) provide direct substrate. Makes the chain citable in the agency's AI Risk Management Plan under M-24-15.

### FQ-4. CUI handling under NIST SP 800-171 Rev 3 — chain artifacts that flow to non-federal systems

**Question.** Chain artifacts in a federal-agency deployment frequently include Controlled Unclassified Information (CUI) — agency decisions about citizens count as CUI Privacy under the CUI Registry (32 CFR Part 2002), and depending on the agency, also CUI Tax (IRS), CUI Banking (Treasury), or CUI Health Information (HHS). When chain artifacts flow to non-federal systems — the cloud HSM endpoint, vendor-hosted ledger storage, and any supplier-side support tooling — NIST SP 800-171 Rev 3 controls apply (May 2024 revision), covering CUI marking, handling, transmission protections, and access controls.

The spec's §10.5 names AWS CloudHSM, Azure Managed HSM, Google Cloud HSM — all FedRAMP-authorized, with AWS GovCloud / Azure Government / Google Cloud for Government variants at FedRAMP High. But the spec does not address CUI marking or handling for chain artifacts. Specifically: (a) chain entries containing CUI flow over OTLP without explicit CUI marking; (b) §10.13's retention list does not name the CUI-handling regime that applies to retained artifacts; (c) §1.2's epistemic-scope language is silent on the CUI sensitivity of "what the AI said" content (which can include citizen PII, tax, health data).

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/cui-handling.md` covering: (1) CUI categorization for chain artifacts (citizen PII at minimum; agency-specific deployments may capture CUI Tax, Banking, Health); (2) SP 800-171 Rev 3 controls when chain artifacts flow to non-federal systems (cloud HSM, vendor-hosted ledger, supplier-side support), with worked control-implementation descriptions; (3) CUI marking guidance for OTLP transport (header attributes naming the CUI category) and storage (markings on audit-file headers and seal records); (4) supplemental SP 800-172 enhanced controls for high-value assets and APT-resistant deployments (relevant for rights-impacting AI under M-23-22).

### FQ-5. CMMC 2.0 scoping for defense / Intelligence Community partner sharing

**Question.** When chain artifacts are shared with Department of Defense or Intelligence Community partners — explicit in my agency's expansion roadmap — the Cybersecurity Maturity Model Certification (CMMC 2.0) framework applies to the vendor and the cloud-service stack. CMMC 2.0 Level 2 (default for CUI-handling) maps to SP 800-171 Rev 3 + a third-party assessment; Level 3 (high-value or advanced-persistent-threat scenarios) adds SP 800-172 enhanced controls + DoD-led assessment. The vendor's SDK build pipeline, the cloud HSM, the ledger storage, and the supplier's verifier-build trust path all need CMMC compliance for defense / IC use. The spec is silent on CMMC scoping. An agency expanding into defense / IC partner sharing has to assemble the CMMC evidence package itself.

**Status.** Gap.

**Closing language proposal.** A new informative annex `regulator-pack/cmmc-2-0-scoping.md` covering: (1) chain components inside the CMMC boundary for defense / IC use (vendor SDK, cloud HSM, ledger storage, verifier-build pipeline); (2) CMMC 2.0 Level 2 control coverage (SP 800-171 Rev 3 controls map directly — much of the CUI-handling annex from FQ-4 carries forward); (3) CMMC 2.0 Level 3 supplemental controls (SP 800-172 enhanced); (4) DoD CIO Cybersecurity Service Provider scoping when cloud HSM or ledger operates as a managed service. Pairs with the CUI-handling annex so an agency expanding from civilian-only to civilian-plus-defense has a continuous evidence path. Defense-side authorization is led by DoD, not the agency's civilian AO — the annex is substrate, not the authorization decision.

### FQ-6. CNSSP-15 / CNSA Suite — National Security Systems eligibility for v1.0 Ed25519 + SHA-256

**Question.** When the chain is used for a National Security System (NSS) — the threshold a federal agency crosses when the AI use case touches classified information or intelligence operations — Committee on National Security Systems Policy 15 (CNSSP-15) governs cryptographic algorithm selection. CNSSP-15 mandates the Commercial National Security Algorithm (CNSA) Suite. CNSA Suite 1.0 names ECDSA P-384 + AES-256-GCM + SHA-384 + RSA-3072+; CNSA Suite 2.0 (2022 update) adds CRYSTALS-Kyber and CRYSTALS-Dilithium with a 2030–2033 PQ transition timeline.

The chain's v1.0 algorithm pair — Ed25519 + SHA-256 — is not currently in CNSA Suite 1.0 (which prefers ECDSA P-384 and SHA-384). Ed25519 is FIPS-validated under FIPS 186-5, but the CNSA Suite preference is explicit. A v1.0 chain is acceptable for FedRAMP Moderate/High and for civilian-agency CUI workloads, but NOT eligible for NSS deployment without algorithm extension. The §4.3.2 algorithm-rotation commitment + §10.10.2 within-day algorithm rotation provide substrate for an NSS-conformant variant — Pattern A cosign with Ed25519 + ECDSA P-384 + SHA-384 — but the spec does not currently acknowledge the NSS gap or define the CNSA-conformant variant.

The cleaner forward path: define an NSS-conformant variant in v1.x under CNSA Suite 1.0 (and a future variant under CNSA Suite 2.0 when PQ migration lands). The variant would be a "posture = nss-cnsa-1.0" option at SDK construct time, dispatched via the §4.1.2 vendor-namespaced constants mechanism, with §4.3 sign_payload extended to bind the algorithm selector.

**Status.** Gap (forward).

**Closing language proposal.** Add a new spec section §10.16 "National Security System eligibility (informative for v1.0; forward-scope for v1.x)" acknowledging that v1.0's Ed25519 + SHA-256 substrate is not currently in CNSA Suite 1.0 and naming the path to an NSS-conformant variant: a v1.x extension defining "nss-cnsa-1.0" (ECDSA P-384 + SHA-384 + AES-256-GCM) and a future "nss-cnsa-2.0" (PQ migration). For v1.0, name the limitation explicitly: civilian-CUI deployments are conformant; NSS deployments require the v1.x extension. Frames the gap as a forward-scope item rather than a v1.0 defect.

### FQ-7. Cross-Domain Solutions for IC partner sharing — JCS canonicalization + binary signatures

**Question.** When chain artifacts cross security domains for IC partner sharing — for example, an agency-internal U//CUI artifact crossing into a TS//SI//NF compartment, or an unclassified summary crossing back to U//CUI — they pass through Cross-Domain Solutions (CDS) as defined under CNSSI 1253 and the IC's Cross-Domain Architecture. CDS systems perform content sanitization, dirty-word filtering, and format validation. The chain's RFC 8785 JCS canonical bytes and binary Ed25519 signatures interact with CDS in two ways: (a) CDS dirty-word filters may flag canonical-JSON content (citizen PII, internal identifiers, model response text) requiring redaction before crossing; (b) CDS format validators may not natively recognize the Ed25519 binary signature shape and may strip or alter it, breaking verification on the receiving side.

**Status.** Gap.

**Closing language proposal.** Add a new section to `docs/edge-and-federated-ai.md` (or a new `docs/cross-domain-sharing.md`) titled "Cross-Domain Solutions composition" covering: (1) CDS interaction with canonical bytes — the agency assesses CDS dirty-word filtering and applies redaction at a layer above the chain (the chain holds the original; the CDS-crossing form is a documented derived product); (2) CDS interaction with binary signatures — Base64-encoding for CDS-compatible transport, with receiving-side decoding for verification; (3) integrity story across the CDS boundary — the original signature does NOT verify on the redacted form; the receiving side either verifies on the redacted form under a re-signed seal in the receiving agency's authority, or on the original chain (which requires receiving-agency access to the original IKM, typically prevented by CDS architecture). Frame the chain as composing alongside CDS rather than passing through it; document the CDS-crossing posture in CC8.1.

### FQ-8. OMB M-22-09 + M-24-04 zero-trust architecture composition

**Question.** OMB M-22-09 ("Moving the U.S. Government Toward Zero Trust Cybersecurity Principles," 26 January 2022) and M-24-04 (the FY2024 progress update, 4 December 2023) require federal agencies to deploy zero-trust architectures. M-22-09 names five pillars: Identity, Devices, Networks, Applications and Workloads, and Data. The chain interacts with each — Identity (SDK authenticates to IKM custodian via SPIFFE / mTLS / HSM-issued tokens per §4.1.1), Devices (SDK runs on agency-managed compliant endpoints), Networks (TLS 1.3 per §5.1; mTLS to ledger), Applications (SDK + ledger + verifier as zero-trust workloads), and Data (classification-aware per FQ-4).

§4.1.1 names mTLS, SPIFFE, and HSM-issued bearer tokens — aligned with zero-trust Identity. But the spec does not compose the chain against the full M-22-09 / M-24-04 architecture: (a) the identity model addresses service-to-service auth but not PIV/CAC user-identity flows for SOC-team or examiner verifier access; (b) the device-trust pillar is silent on whether the SDK host must be on the agency's compliant-endpoint roster; (c) the data-pillar interaction is implicit — the agency's zero-trust data-classification engine needs to know how to read chain content sensitivity.

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/zero-trust-composition.md` covering: (1) the chain's M-22-09 / M-24-04 pillar composition — Identity (§4.1.1 mechanisms + PIV/CAC for verifier access), Devices (compliant-endpoint posture for SDK and verifier hosts), Networks (TLS 1.3 per §5.1 + OTLP auth per §4.4.3), Applications (SDK / ledger / verifier as zero-trust workloads with attestation), Data (CUI classification per FQ-4); (2) cross-walk to the CISA Zero Trust Maturity Model 2.0 stages (Traditional / Initial / Advanced / Optimal) so M-24-04 progress reporting can cite the chain's contribution; (3) CC8.1 language naming which pillar each chain control supports. Makes the chain named substrate rather than ambient infrastructure.

### FQ-9. CISA BOD 22-01 + BOD 23-01 — vulnerability management and asset visibility for chain components

**Question.** CISA Binding Operational Directives apply to federal civilian executive-branch agencies. BOD 22-01 ("Reducing the Significant Risk of Known Exploited Vulnerabilities," 3 November 2021) requires remediation of vulnerabilities in the CISA Known Exploited Vulnerabilities (KEV) catalog within specific timelines. BOD 23-01 ("Improving Asset Visibility and Vulnerability Detection on Federal Networks," 3 October 2022) requires comprehensive asset inventories and vulnerability scanning at specified frequencies. The chain SDK, ledger server, verifier binary, and cloud HSM endpoints are all federal IT assets subject to both BODs.

The spec is largely silent on vulnerability management. §10.13 names SDK version manifest and source-code hash (good for forensic reproduction) but does not address CVE / KEV posture — when a KEV-listed vulnerability hits a transitive dependency in the SDK or verifier, what is the agency's remediation path and timeline? `docs/supply-chain.md` provides some substrate but does not connect to BOD vocabulary explicitly.

**Status.** Partial.

**Closing language proposal.** Extend `docs/supply-chain.md` with a "CISA BOD posture" section covering: (1) BOD 22-01 alignment — KEV-catalog monitoring procedure, remediation timelines (15 days critical / 30 days high), CC8.1 patch-management documentation; (2) BOD 23-01 alignment — chain components in the agency CDM asset inventory, vulnerability-scanning cadence (7-day scan / 14-day asset discovery default), chain-event-to-CDM-telemetry mapping; (3) vendor responsibility for upstream CVE disclosure timing and the agency's CC8.1 procedure evaluating vendor patches against BOD 22-01 timelines; (4) worst case (a chain component itself enters KEV) connects to IR Scenario 4 (HSM key compromise) and the seal-job continuity discipline for ledger-component patching. Makes the chain's vulnerability-management posture BOD-aligned and CDM-citable.

### FQ-10. NARA records-retention scheduling and OMB Circular A-130

**Question.** OMB Circular A-130 ("Managing Information as a Strategic Resource," 28 July 2016) requires federal agencies to manage information — including audit trails of agency decisions — as a strategic resource. Federal records are subject to National Archives and Records Administration (NARA) records-retention schedules. The General Records Schedule (GRS) covers most administrative records; agency-specific schedules cover mission records. Chain artifacts are clearly federal records: the captured AI decision is a record of agency action.

§10.9 and §10.13 set retention floors of 7 years (the FFIEC banking horizon). Federal agencies have different obligations: some mission records — the AI decisions of a Social Security disability claim, an IRS audit, a VA benefits adjudication — are permanent records under their NARA schedule. Permanent records are eventually transferred to NARA; during the agency-custody period (often decades), they must be preserved in a NARA-ingestible format (M-23-07 transition-to-electronic-records).

The spec does not address federal records-management. An adopting agency has to figure out: (a) which chain artifacts are temporary vs permanent records under their NARA schedule; (b) how the 7-year floor composes with multi-decade permanent-record obligations; (c) the format-preservation discipline — Ed25519 key rotation, SHA-256 weakening, §4.3.2 algorithm rotation, §10.9 IKM-retention coupling — all need to compose with permanent-record stewardship.

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/federal-records-management.md` covering: (1) records categorization for chain artifacts (operational events typically temporary; AI-decision events often permanent under agency-specific NARA schedules); (2) OMB Circular A-130 information-resources-management posture, with a worked example citing the chain in the agency's plan; (3) NARA Electronic Records Management composition — format-preservation across multi-decade horizons, algorithm-rotation discipline (§4.3.2 becomes load-bearing for permanent records), M-23-07 compliance; (4) CC8.1 language for permanent-record handling — the IKM-retention coupling extends from 7 years (FFIEC floor) to the permanent-record horizon, and algorithm-rotation becomes a permanent-records preservation control rather than a cryptanalysis-defense control.

### FQ-11. FedRAMP 3PAO assessment procedures — NIST SP 800-53A Rev 5

**Question.** FedRAMP authorizations are supported by Third-Party Assessment Organization (3PAO) assessments performed under NIST SP 800-53A Rev 5. Each control in the SP 800-53 Rev 5 catalog has a corresponding SP 800-53A Rev 5 assessment objective set (objectives + methods: examine / interview / test). The 3PAO's Security Assessment Plan (SAP) lays out procedures; the Security Assessment Report (SAR) documents results; the Plan of Action and Milestones (POA&M) tracks remediation.

Marcus Brennan's SOC 2 drop covers AICPA TSC mapping and engagement procedures. The federal equivalent under FedRAMP uses different procedures (SP 800-53A Rev 5 assessment objectives, FedRAMP-specific SAP / SAR / POA&M templates) and is performed by FedRAMP-recognized 3PAOs, not AICPA firms. The chain's testability story is strong — byte-reproducible verifier output, normative audit procedures, explicit control prescriptions — but the spec lacks a 3PAO procedure template parallel to what Brennan asked for under SOC 2.

**Status.** Partial.

**Closing language proposal.** A new `docs/regulator-pack/fedramp-3pao-procedures.md` (parallel to Brennan's MQ-2 SOC 2 doc) covering 3PAO procedures for each chain-prescribed control in SP 800-53A Rev 5 assessment-objective format: objectives per control, examine / interview / test methods, expected evidence, FedRAMP SAR-scoring disposition. Anchor to §10.1 reconciliation (continuous-monitoring substrate), §7 verifier output (deterministic test substrate), §10.2 operational events (interview + examine substrate for HSM custody, RNG, append-only), §10.13 retention (examine substrate). Gives 3PAOs a starting point that isn't written from scratch on first engagement.

### FQ-12. Reciprocity across agency authorizations under FedRAMP and the FedRAMP Authorization Act

**Question.** FedRAMP authorizations are designed to be reusable — an authorization issued by one agency or by the FedRAMP JAB is intended to be reciprocally accepted by other agencies, reducing duplicate authorization work. The FedRAMP Authorization Act (signed December 2022 as part of the FY2023 NDAA) codified this reciprocity intent. Reciprocity depends on: (a) consistent SSP control mapping (FQ-1), (b) consistent 3PAO assessment procedures (FQ-11), (c) consistent CUI handling (FQ-4), and (d) a common authorization-package shape downstream agencies can consume.

If each agency writes its own mapping, runs its own 3PAO under its own templates, and assembles its own package, the reciprocity intent is lost — the chain becomes a bespoke per-agency authorization rather than a common federal substrate. The spec does not address reciprocity. An agency leveraging an upstream agency's FedRAMP authorization needs the upstream package in a shape the downstream RMF recognizes, and the upstream agency needed the spec's substrate to assemble it consistently.

**Status.** Partial.

**Closing language proposal.** A new informative annex `regulator-pack/fedramp-reciprocity.md` covering: (1) the reciprocity intent under the FedRAMP Authorization Act and the AO's acceptance posture; (2) the authorization-package shape that supports reciprocity — the SSP mapping (FQ-1), 3PAO SAR template (FQ-11), CUI handling (FQ-4), M-23-22 mapping (FQ-2), AI RMF mapping (FQ-3) — as reusable substrate the upstream agency assembles once and downstream agencies inherit; (3) CC8.1 posture for the authorization boundary — components inside (ledger, HSM, verifier-build pipeline) vs outside (agency host running the SDK, agency identity-management). Frames the chain as reciprocity-friendly substrate rather than per-agency authorization burden.

### FQ-13. Authorization Boundary and System Categorization under FIPS 199

**Question.** Under FISMA, every federal information system has a security categorization under FIPS 199 — Low / Moderate / High — derived from the worst-case impact across the Confidentiality / Integrity / Availability triad. The categorization determines the FedRAMP baseline and the SP 800-53 Rev 5 control set. The Authorizing Official approves the categorization as part of the ATO.

For an AI system using the chain, the categorization is driven by the AI use case's impact level — a citizen-facing eligibility decision (rights-impacting under M-23-22) is at minimum Moderate, often High. The chain inherits the categorization of the AI system it serves. But the chain's authorization boundary — components inside vs outside the boundary — needs explicit definition for the SSP. §4's implementation topology (monolithic vs distributed) is informative on construction but not on boundary scoping. An agency CISO needs to know which components are inside (subject to the agency's full SP 800-53 Rev 5 control set) and which are outside (subject to vendor SLAs, third-party attestation, complementary user-entity controls).

**Status.** Partial.

**Closing language proposal.** A new section in `docs/regulator-pack/deployment-package.md` (or a new `docs/regulator-pack/fedramp-authorization-boundary.md`) covering: (1) FIPS 199 categorization for chain-of-custody — the chain inherits the AI use case's CIA impact level (rights-impacting typically Moderate / High; safety-impacting typically High); (2) authorization-boundary scoping — components inside (SDK on agency hosts, IKM registry, verifier execution environment) vs outside (cloud HSM service, vendor-hosted ledger, vendor-distributed SDK build artifacts), with corresponding SP 800-53 coverage per split; (3) complementary user-entity controls (CUECs) for components outside the boundary (HSM access reviews, key-rotation discipline, vendor SOC 2 / ISAE 3402 / FedRAMP attestation review, supply-chain trust path validation); (4) CC8.1 documentation of the inside/outside split. Makes the authorization boundary explicit and consistent across adopting agencies.

---

## Per-role roll-up

**For the spec working group.** Eleven Partials, one Gap (CMMC), one forward-Gap (CNSSP-15 / CNSA Suite). None affect the cryptographic-integrity claim. All affect the federal-authorization layer — SSP, 3PAO assessment, M-23-22 / M-24-15 alignment, records-retention, zero-trust composition. The pragmatic close-out path is a set of regulator-pack annexes (`nist-800-53-rev5-mapping.md`, `omb-m-23-22-alignment.md`, `nist-ai-rmf-mapping.md`, `cui-handling.md`, `cmmc-2-0-scoping.md`, `zero-trust-composition.md`, `federal-records-management.md`, `fedramp-3pao-procedures.md`, `fedramp-reciprocity.md`, `fedramp-authorization-boundary.md`), an extension to `docs/supply-chain.md` for CISA BOD posture, and a forward §10.16 in the spec acknowledging NSS eligibility as v1.x scope. Each annex layers on the existing normative core; none requires changes to normative text.

**For federal civilian agencies adopting the chain.** Technically FedRAMP-friendly — FIPS-current cryptography, deterministic ConMon-friendly verifier output, operational events that map to AU-12 / CA-7 control families. The authorization-package layer is where each adopting agency currently writes substantial substrate themselves. Closing the gaps reduces per-agency burden and enables FedRAMP reciprocity. The OMB AI memos (M-23-22, M-24-10, M-24-15) persist across administration changes; framing the chain's role around the persistent memos rather than the volatile EO of the day is the right defensive posture.

**For agencies expanding into defense / IC partner sharing.** v1.0 is civilian-CUI ready. Defense / IC use requires FQ-5 (CMMC scoping), FQ-6 (CNSSP-15 / CNSA Suite NSS eligibility, v1.x scope), and FQ-7 (Cross-Domain Solutions). Best framed as a separate authorization track — DoD-led for CMMC, IC-CIO-led for IC-CD — with the v1.x NSS variant as a deliberate working-group scope item rather than a v1.0 patch.

**For the agency CAIO and AI Governance Board.** The chain is the technical substrate for several M-23-22 minimum risk-management practices, M-24-10 Annual AI Use Case Inventory fields, and M-24-15 AI Risk Management Plan components. The board can cite the verifier output as evidence of operational effectiveness, §10.1 reconciliation as monitoring evidence, and §10.13 retention as substrate for incident-response and accountability. The §1.2 epistemic-scope discipline ("the chain proves what the AI said and that the record was not tampered with; not that the AI's statement is factually accurate, policy-compliant, or unbiased") is the right framing — the chain is the integrity foundation, not the truth foundation, and the agency's other AI-governance controls compose alongside it.

---

## Where I would prioritize

1. **FQ-1 (SP 800-53 Rev 5 control-mapping table).** Highest leverage. Every adopting agency's SSP starts here. One annex closes the largest per-agency authorization-burden item and enables FedRAMP reciprocity.

2. **FQ-2 (OMB M-23-22 alignment).** Time-sensitive. M-23-22, M-24-10, M-24-15 are operationally active — agencies are filing AI Use Case Inventory submissions and AI Risk Management Plans now.

3. **FQ-4 (CUI handling under SP 800-171 Rev 3).** Every adopting agency confronts CUI flow to non-federal systems. One annex.

4. **FQ-3 (NIST AI RMF 1.0 subcategory mapping).** Pairs with FQ-2; the framework agencies cite in their AI Risk Management Plan under M-24-15.

5. **FQ-11 + FQ-12 (3PAO procedures + reciprocity).** Pair tightly. 3PAO procedures are per-agency assessment substrate; reciprocity is cross-agency leveraging. Closing both enables upstream packages to be inherited downstream.

6. **FQ-8 (zero-trust composition).** M-22-09 / M-24-04 progress reporting needs the chain's contribution in zero-trust terms.

7. **FQ-13 (authorization boundary + FIPS 199 categorization).** Per-agency but consistent scoping saves each agency the boundary-design decision.

8. **FQ-9 (CISA BOD posture).** Extension to `docs/supply-chain.md`; lower urgency because the supply-chain doc already provides substrate.

9. **FQ-10 (NARA records-retention).** Permanent-record handling pushes §4.3.2 algorithm-rotation into preservation territory — surface it before agencies discover it the hard way.

10. **FQ-5 (CMMC scoping).** Forward-scope for defense / IC expansion; pairs with FQ-4.

11. **FQ-7 (Cross-Domain Solutions).** Forward-scope; IC partner sharing only.

12. **FQ-6 (CNSSP-15 / CNSA Suite NSS eligibility).** Forward-scope v1.x track. Frame in v1.0 as acknowledged limitation; defer the variant.

---

## Stopping criterion

This drop carries 11 Partial findings, 1 Gap (CMMC scoping for defense / IC use cases), and 1 forward-Gap (CNSSP-15 / CNSA Suite NSS eligibility, framed as v1.x scope) against the v1.0-final-amendment spec — all clustered around the federal-authorization layer. None affect the spec's cryptographic-integrity claim. None affect the implementation contract or the verifier procedure.

The chain is technically clean for federal authorization; the gap is the authorization-package layer that today is written by each adopting agency from scratch on its own ATO timeline. Closing the gaps above would produce consistent SSP control attribution across agencies, enable FedRAMP reciprocity, and give agency CAIOs and AI Governance Boards a citable substrate in their AI Risk Management Plans under M-24-15.

The pragmatic close-out path is a set of supplementary regulator-pack annexes, an extension to the supply-chain doc, and a forward §10.16 in the spec body acknowledging NSS eligibility as v1.x scope. None of the recommended changes touch the normative core.

For the ATO decision: v1.0 is a defensible substrate for a civilian-CUI AI-decision audit trail under FedRAMP Moderate or High, conditional on the adopting agency assembling the authorization-package layer themselves. With the annexes in place, v1.0 becomes a reciprocity-friendly substrate that materially reduces per-agency burden. v1.0 is NOT yet NSS-eligible; the v1.x NSS-conformant variant is the path for IC partner sharing, and framing it as a deliberate v1.x scope item rather than a v1.0 deficiency is the right posture.

The chain is technically ready. The authorization-package work is what stands between the technical readiness and the ATO signature.
