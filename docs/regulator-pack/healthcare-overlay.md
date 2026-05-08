# Healthcare regulatory overlay

> **What this doc is.** A healthcare-specific overlay covering HIPAA Privacy and Security Rules, the HIPAA Breach Notification Rule, FDA software-as-a-medical-device guidance (including the 2024 Predetermined Change Control Plan final guidance), HITRUST CSF v11, the ONC 21st Century Cures Act audit-log certification criterion at 45 CFR 170.315(d)(2), CMS-0057-F prior-authorization timeliness, state AI laws (Colorado AI Act, California SB 1120 and similar), Joint Commission patient-safety standards, and the Patient Safety and Quality Improvement Act privilege framework. The chain spec is sector-neutral; this overlay is what a US health system needs alongside the spec to deploy the chain in clinical operations.

> **Audience.** Health-system CISO, HIPAA Privacy Officer, HITRUST assessor, FDA regulatory-affairs lead, ONC compliance lead, CMS-0057-F program owner, AI Governance Committee chair, Patient Safety Evaluation System administrator, and the institution's General Counsel. Hospital and integrated-delivery-network deployments only — small physician-practice and ambulatory-surgical-center deployments will scale this down to fit their HITRUST scope and HIPAA risk profile.

> **What this overlay does NOT do.** It does not change the spec's normative core. It does not replace the institution's HITRUST CSF v11 control description, FDA Quality System submission, ONC certification posture, CMS-0057-F operational program, or PSE policy. It composes alongside those existing instruments and names where the chain is the integrity foundation for the AI-decision audit trail those instruments rely on.

---

## How to read this document

The healthcare regulatory landscape touches the chain at six surfaces. Each surface gets its own section. The sections are designed to read independently — the FDA section is for the regulatory-affairs lead, the HIPAA sections are for the Privacy Officer, the HITRUST section is for the assessor, and so on. Cross-references between sections are explicit so a reader chasing a particular thread does not have to read the whole document to see how the pieces fit together.

The chain captures AI-decision integrity. Healthcare regulation cares about more than integrity — it cares about minimum-necessary use, breach-notification timing, FDA-device-status preservation, EHR-audit-log certification posture, prior-authorization timeliness, patient-safety reporting privilege, and clinician-override review. The integrity foundation the chain provides is necessary for several of those concerns and irrelevant to others. This overlay names which is which so the institution does not over-claim or under-claim what the chain does for healthcare regulation.

---

## §H1. PHI in chain artifacts — three operational postures

Clinical AI input vectors are PHI. Sepsis prediction takes lactate, MAP, white-cell count, ICD-10 codes for current admission, and recent vital signs. Every one of those is identifying information when combined with admission timing and the small per-day per-institution ICU sepsis population. Tokenizing the lactate value defeats the audit purpose; preserving it puts PHI in the chain artifact and inherits the full HIPAA Security Rule control set on the chain ledger.

The institution chooses one of three postures per deployment and names the choice in the institution's HITRUST CSF v11 control description (closer mapping is HITRUST 06.d Change Management or 09.aa Audit Logging). The choice is made by the AI Governance Committee in consultation with the HIPAA Privacy Officer.

### §H1.1 Posture (a) — full-content capture

The chain entry records the clinical input vector and the AI output as captured. The entry inherits PHI status under HIPAA. The chain ledger inherits HIPAA Security Rule §164.312 controls — encryption at rest, access logging, workforce training, audit reviews. The institution names the ledger as a designated record set under 45 CFR 164.501 if the institution's policy assigns that designation.

This posture closes the auditability question cleanly. The verifier output reads exactly the bytes the AI consumed and produced. Forensic investigation into a clinical-decision incident can reconstruct the AI's input/output without consulting the EHR. The cost is that the chain ledger becomes a full HIPAA-controlled environment with the full retention and breach-notification surface that implies. The institution's vendor-management evidence under HITRUST 05.k (Service Provider Compliance) covers any vendor operating any part of the ledger.

Use this posture when the AI's input/output is not stored elsewhere with sufficient integrity controls, when forensic-reconstruction needs are high (the prior-authorization-decision case where a regulator may ask "show me exactly what the AI saw and said"), or when the institution's policy treats the chain as the primary AI-decision record set.

### §H1.2 Posture (b) — hash-only capture

This is the recommended default. The chain entry records `audit.clinical.input_hash` and `audit.clinical.output_hash` as lowercase hex SHA-256 of the canonicalized clinical input and AI output. The input/output bytes themselves live in the EHR under existing HIPAA-compliant EHR controls (Epic Hyperspace, Cerner Powerchart, Meditech, Allscripts/Veradigm, athenahealth — all certified under 45 CFR 170.315 for audit-log content). The chain proves what the AI saw and said without holding the PHI bytes; integrity-binding rests on SHA-256's collision-resistance assumption.

This is structurally the same posture spec §10.11 took for ECOA translation evidence: capture the hash under the per-event MAC binding, keep the source bytes under the institution's existing controlled environment. The chain inherits no PHI status under this posture, the chain ledger remains a non-HIPAA-controlled environment, and the institution's HITRUST scope on the chain is bounded to the ledger-as-non-PHI-store.

The verifier reads the chain entry's hash. To reconstruct a specific decision, the institution pulls the EHR-stored bytes by the chain entry's `audit.ehr.encounter_id` (or equivalent linkage attribute), recomputes the SHA-256 over the canonical bytes, and compares against the chain entry's recorded hash. A match proves the chain is unaltered; a mismatch identifies either a chain-integrity break or an EHR-side modification of the source bytes.

Use this posture for the sepsis-prediction case, the clinical-documentation-AI case, and any deployment where the EHR is the natural source of truth for the AI's input/output.

### §H1.3 Posture (c) — limited data set

Under 45 CFR 164.514(e) the chain operates on a Limited Data Set with a Data Use Agreement. The chain entry holds non-direct identifiers (no name, no address, no contact information, no MRN) but may hold dates and ZIP codes. The institution accepts that 45 CFR 164.524 (HIPAA Right of Access) requests flow to the EHR rather than the chain. The DUA names the chain's purpose (decision-integrity audit), the institution's use restriction (no marketing, no fundraising, no contact, no re-identification), and the chain's data-handling controls.

This posture is operationally awkward — the institution maintains both an LDS-scoped chain and a full-PHI EHR, and reconciling the two during forensic investigation requires the institution to maintain the linkage. Use this only when the institution's privacy posture explicitly favors LDS scoping (some research-heavy academic medical centers prefer it; most operational-clinical environments do not).

### §H1.4 Posture-selection record

The institution's HITRUST CSF v11 control description names the choice and the rationale. The AI Governance Committee meeting minute documenting the choice is retained per HITRUST 06.b (Information Backup) for the chain's full retention horizon. The HIPAA Privacy Officer signs the minute. The institution's SOC engagement tests for the existence of the minute and the consistency of the chain's actual capture with the declared posture.

---

## §H2. HIPAA minimum-necessary analysis applied to the chain

HIPAA §164.502(b) limits PHI use and disclosure to the minimum necessary for the stated purpose. HHS-OCR's enforcement posture treats "comprehensive logging for compliance purposes" as a category subject to documented purpose-versus-breadth analysis, not a category that enjoys an automatic minimum-necessary safe harbor. The chain's spec §1.2 epistemic-scope statement (the chain proves what the AI said and that the record was not tampered with; it does not prove the AI's statement is factually accurate, policy-compliant, or free of bias) is the right answer for the integrity claim, but it does not function as a minimum-necessary justification.

The institution conducts a written minimum-necessary analysis at deployment and reviews it annually under HITRUST CSF v11 control 19.b (Information Privacy and Protection — Minimum Necessary). The analysis answers four questions.

### §H2.1 Purpose

The chain's defined purpose under clinical deployment names the regulatory consumers the institution is satisfying. The standard set: integrity audit of AI decisions for FDA SaMD post-market surveillance (per §H5), Joint Commission patient-safety reporting (per §H12), ONC information-blocking compliance (per §H6), CMS-0057-F prior-authorization timeliness evidence where applicable (per §H7), state-AI-law disclosure where applicable (per §H10), and the institution's internal AI-governance program. The institution names which of these apply per deployment; not every deployment hits every consumer.

### §H2.2 Field inventory

Every field the chain captures is classified as one of:

- **Identifier under 45 CFR 164.514(b)(2).** Names, MRN, account numbers, dates more granular than year, ZIP codes more granular than three-digit, biometric identifiers, full-face photographs, and the residual catch-all. The chain captures none of these directly; identifiers flow through the institution's privacy-store / token-vault per `docs/privacy-by-design.md`.
- **Clinical finding.** Lactate, MAP, ICD-10 codes, vitals, lab values, medication doses. Under HHS-OCR's broader interpretation these are PHI when combined with admission timing and a small per-day per-institution population. Captured under one of the §H1 postures.
- **Decision metadata.** Model identifier, model version, routing-decision attributes, deployment-intent attributes, AI-recommendation result codes, confidence scores. Not PHI in the strict identifier sense but operationally sensitive.
- **Non-PHI infrastructure metadata.** `tenant_id`, `run_id`, `seq`, `prev_hash`, `payload_hash`, `key_version`, `key_fingerprint`, `captured_at`, `received_at`. These are the chain's structural fields and carry no PHI.

### §H2.3 Per-field necessity

For each PHI-bearing field, the institution documents the reason capture is necessary, the alternatives considered (hash-only per §H1.2, exclusion, partial redaction), and the rejection rationale. The analysis cites the regulatory consumer the field serves: a lactate value captured under §H1.2 hash-only posture is necessary because the FDA post-market surveillance program for the sepsis device requires the institution to be able to reproduce the AI's input/output on demand; a hash bound under per-event MAC plus EHR-side byte-storage satisfies that requirement without putting the lactate value in the chain artifact.

### §H2.4 Annual review

The analysis is reviewed when any of: HHS-OCR enforcement guidance updates, HITRUST CSF revision (the latest is v11; v12 is in working-draft as of 2026), FDA SaMD guidance revision, ONC certification criterion update, or AI Governance Committee determination triggered by a new deployment or a change in regulatory posture. Annual is the floor; events trigger ad hoc review.

The HIPAA Privacy Officer and the AI Governance Committee chair both sign. The institution's CC8.1 (HITRUST 06.d / 09.aa equivalent) cites the analysis by document identifier and review date.

---

## §H3. HIPAA Breach Notification — chain anomalies as §164.404(b) discovery events

Under 45 CFR 164.404(b), the breach-notification 60-day clock starts at the time the breach is "discovered or, by exercising reasonable diligence, would have been known." The chain itself is what makes "would have been known" enforceable for chain-detected anomalies. When the verifier outputs `key_fingerprint mismatch at seq N` and the affected entries reference PHI access (i.e., the deployment is under §H1.1 or §H1.2 posture), HHS-OCR's enforcement position is that the verifier output is the discovery moment under §164.404(b). The institution cannot argue "we did not know" once the verifier has flagged the entry.

This composes with the FFIEC 36-hour clock per `docs/regulator-pack/breach-notification-matrix.md`. The two clocks run concurrently when both apply. Which row of the matrix governs is determined by what the affected entries reference (PHI access only → HIPAA only; financial-services-regulated activity → FFIEC; both → both).

### §H3.1 The four-factor risk assessment

Under 45 CFR 164.402(2), an unauthorized acquisition, access, use, or disclosure of PHI is presumed a breach unless the institution demonstrates low probability of compromise via a four-factor risk assessment:

1. The nature and extent of the PHI involved (what was the data, was it identifiable, how sensitive).
2. The unauthorized person who used or to whom the disclosure was made (role, intent if knowable, training/sanctions exposure).
3. Whether the PHI was actually acquired or viewed (was the data only theoretically exposed, or was there evidence of access).
4. The extent to which risk has been mitigated (was the data recovered, were assurances obtained, were systems hardened).

The institution's IR playbook ties to this. The first step after a chain-detected anomaly that touches a deployment under §H1.1 or §H1.2 posture is the four-factor assessment. The HIPAA Privacy Officer signs the assessment contemporaneously. If the assessment concludes low probability of PHI compromise (the §164.402(2) safe harbor), the institution documents the conclusion and HHS-OCR notification is not triggered. If inconclusive or non-low, the 60-day clock applies and the institution proceeds to notification under §164.404.

### §H3.2 Concurrency with FFIEC and other clocks

When a chain-detected anomaly affects a deployment that touches both PHI access and FFIEC-regulated activity (the bank's healthcare-financing arm doing prior authorization on behalf of a Medicare Advantage plan, for example), both clocks run concurrently. The institution's IR coordinator runs both notification tracks in parallel. The four-factor assessment serves as the bridging document — its conclusion drives the HIPAA notification decision; the FFIEC 36-hour clock proceeds independently.

DORA's 24-hour clock (when the institution is in DORA scope), GDPR's 72-hour clock (when the affected data subjects include EU residents), and any state breach-notification clock (typically 30-60 days, varies by state) compose alongside. The institution's `docs/regulator-pack/breach-notification-matrix.md` is the source of truth for the full matrix; this section adds the HIPAA-specific composition rules.

### §H3.3 Contemporaneous documentation

The four-factor assessment is documented contemporaneously, not reconstructed after the fact. HHS-OCR enforcement records show that institutions producing reconstructed assessments after a complaint is filed routinely lose at the resolution-agreement stage. The IR playbook's standard scenario file for a chain-detected anomaly under PHI scope includes a four-factor assessment template that the on-call HIPAA Privacy Officer fills out at the moment of escalation. The signed template is retained in the chain's evidentiary-artifact set under §10.13 retention coupling.

---

## §H4. HIPAA 6-year retention floor and state-law overlays

45 CFR 164.530(j)(2) requires 6-year retention of documentation of policies, procedures, communications, and actions required by Subparts A and E of 45 CFR 164. AI-decision logs are within scope when the deployment is under §H1.1 or §H1.2 posture. State law often imposes longer floors that govern when they exceed HIPAA.

### §H4.1 The state-law matrix

The retention floor is the longest applicable: federal HIPAA, state medical-records law, state pediatric-records law, applicable insurance-records law, applicable utilization-management-records law. For a multi-state institution, the matrix runs to 50+ rows. A representative sample:

| State | Adult records | Pediatric records | Notes |
|---|---|---|---|
| California | 7 years (Cal. Health & Safety Code §123145) | 7 years past age 18 | Adult floor exceeds HIPAA. |
| New York | 6 years (NY Public Health Law §18) | 6 years past age 18 | Pediatric floor extends to ~24 years. |
| Illinois | 10 years past last visit | 10 years past age 18 | Both floors significantly exceed HIPAA. |
| Texas | 7 years from last visit | 7 years past age 18 | Adult floor exceeds HIPAA. |
| Florida | 5 years (HCFA-licensed facilities); 7 years (physician records) | Generally 7 years past age 18 | Provider-type-specific. |
| Pennsylvania | 7 years from last visit | 7 years past age 18 | Adult floor exceeds HIPAA. |
| Ohio | 6 years (FFIEC alignment) | 21 years (10 years past age of majority, conservative reading) | Pediatric floor very long. |
| Massachusetts | 7 years | 7 years past age 18 | Adult floor exceeds HIPAA. |
| Michigan | 7 years | 7 years past age 18 | Adult floor exceeds HIPAA. |
| Georgia | 10 years | 10 years past age 18 | Both floors exceed HIPAA. |

The institution's HITRUST CSF v11 control 06.c (Records Retention) names the matrix and the per-record disposition rules. The institution's general counsel maintains the matrix as a living document; updates flow through the institution's standard change-management process.

### §H4.2 IKM retention coupling

The chain's §10.9 retention coupling extends the institution's IKM-availability obligation to the longest applicable floor on any chain entry retained under that key_version. For pediatric chains in a state with 21-year floors, the institution plans for 25+ year IKM availability — well beyond the cloud HSM key-lifecycle defaults at AWS CloudHSM, Azure Managed HSM, and Google Cloud HSM, which typically default to 5-10 year customer-key-lifecycle horizons.

The institution's vendor-management evidence under HITRUST 05.k (Service Provider Compliance) includes the HSM vendor's commercial commitment to multi-decade key retention. Most cloud providers will commit to longer horizons under a custom enterprise agreement; the institution negotiates the commitment at contracting time and renews annually as the matrix evolves. The institution's chain-vendor MSA includes the same commitment cascading downward (per `docs/vendor-hosted-controls.md`).

### §H4.3 Re-keying as a retention-management option

Institutions facing 25+ year IKM-availability requirements may operate a re-keying program: at year N, re-encrypt or re-MAC the historical chain under a new IKM, retire the old IKM after a documented overlap window, and continue the retention horizon under the new key. This is a non-trivial operational program — re-keying changes the chain's per-event MAC values and breaks any verifier output produced under the old IKM unless the institution archives the old verifier output as a fixed artifact.

The recommended posture for institutions not already operating a re-keying program is to negotiate multi-decade HSM key-lifecycle commitments with the cloud HSM vendor and monitor the vendor's commitment evolution annually. Re-keying is a forward-scope workstream the institution may adopt later if the multi-decade commitment becomes infeasible.

---

## §H5. FDA software-as-a-medical-device evidence capture

FDA's December 2024 final guidance on Predetermined Change Control Plans (PCCPs) for AI/ML-enabled device software functions changed the post-market posture for clinical AI. Manufacturers can submit a PCCP at the original 510(k)/De Novo clearance/authorization that pre-authorizes specific categories of model changes; changes within the PCCP do not require new submission. Changes outside the PCCP do require new 510(k) or De Novo submission.

The chain captures decision history. Spec §4.4.1 routing-decision capture and §4.4.2 deployment-intent capture cover parts of the model-change story. The FDA-relevant question is whether the chain provides evidence supporting the manufacturer's PCCP-boundary distinction between pre-authorized changes and changes requiring new submission. Today the spec does not name FDA SaMD as a consumer of chain evidence and does not normate any attribute set tied to the PCCP boundary.

### §H5.1 Recommended attribute set (informative)

For chain entries representing model invocations of FDA-regulated AI, the institution emits the following attributes alongside the standard chain entry:

- `audit.fda.samd.device_identifier` — the device's UDI (Unique Device Identifier) or 510(k)/De Novo identifier as registered with FDA.
- `audit.fda.samd.model_version` — the algorithm version, expected to match the manufacturer's PCCP-bound version registry entry.
- `audit.fda.samd.pccp_boundary` — one of `pre_authorized` (the invocation operates within the PCCP-pre-authorized version) or `out_of_pccp` (the invocation operates outside the pre-authorized boundary, triggering a new-submission obligation). The institution's regulatory-affairs team owns this classification.
- `audit.fda.samd.intended_use_population` — the indications-for-use population the invocation served (e.g., `adult_icu_sepsis_screening`, `pediatric_radiology_chest_xray_initial_read`).

These are institution-emitted under the canonical bytes per spec §5; the chain's per-event MAC binds them; the verifier treats them as ordinary `audit.*` attributes (no PASS/FAIL impact on chain integrity).

### §H5.2 Operational consumer

The manufacturer's regulatory-affairs team (which may be the institution itself for institution-developed devices, or a vendor for vendor-developed devices) reads the chain at quarterly cadence for the PCCP-boundary distribution: how many invocations were `pre_authorized` versus `out_of_pccp`. Entries with `pccp_boundary = "out_of_pccp"` are an FDA-submission trigger. The institution's HITRUST CSF v11 control 02.b (Authorized Software Use) names the consumption procedure.

For institution-developed devices where the institution is the manufacturer of record, the institution's regulatory-affairs team consumes the chain directly. For vendor-developed devices, the vendor's regulatory-affairs team consumes via the institution-shared chain extract; the contract names the cadence and the data-protection terms.

### §H5.3 Forward-scope

FDA's SaMD landscape continues to evolve through 2026-2027. The spec working group does not normate the §H5.1 attribute set in the chain spec body until FDA publishes specific audit-log guidance — a candidate work item in FDA's 2025-2026 AI-software workstream. The institution adopting these attributes today is operating under informative guidance; the attribute names may shift if FDA publishes specific naming conventions.

---

## §H6. ONC EHR audit-log integration

Every ONC-certified EHR — Epic, Cerner/Oracle Health, Meditech, Allscripts/Veradigm, athenahealth, NextGen, eClinicalWorks, and the rest — implements 45 CFR 170.315(d)(2), the audit-record content certification criterion: date/time of action, type of action, patient identification, user identification, originating IP. The 21st Century Cures Act information-blocking rule at 45 CFR 171 makes the audit log itself patient-accessible under the right of access where PHI is disclosed.

The institution operating clinical AI deploys the chain alongside the EHR's native audit log. Three integration patterns are available; the institution selects per deployment and names the choice in the institution's HITRUST 09.aa control description.

### §H6.1 Pattern (a) — parallel-stream

The chain operates beside the EHR audit log without integration. The EHR's audit log retains its §170.315(d)(2) certification posture. The chain captures AI-decision-specific evidence the EHR's native audit log does not: model version, prompt, parameters, routing decisions, deployment-intent attributes. The two streams are reconciled at examination time via a documented mapping in the institution's HITRUST 09.aa control description: for any AI-driven event that touches the EHR, the institution can produce both the EHR audit-log entry and the corresponding chain entry, time-aligned to within the institution's NTP-synchronization budget.

This is the recommended default for most deployments. It preserves ONC certification, keeps the chain's evidentiary openness independent of the EHR's audit-log-as-patient-accessible-document posture, and allows the institution to evolve the chain's capture independently of the EHR vendor's product roadmap.

### §H6.2 Pattern (b) — composed-log

The chain's wire form is consumed by the EHR's audit-log subsystem (Epic Caboodle, Cerner CCL audit views, Meditech audit-log import paths) so the EHR's §170.315(d)(2) report includes the AI-decision evidence inline. Operationally heavy. Some EHR vendors expose suitable hooks (Epic's HL7 v2 audit gateway is the most mature); others do not. The integration adds the EHR vendor as a consumer of the chain's wire format and creates a dependency between chain evolution and EHR-vendor product cycles.

Recommended only where the institution has an existing EHR-vendor integration commitment that makes the operational cost worthwhile (large academic medical centers with mature Epic Caboodle deployments are the typical fit). The institution's vendor-management evidence under HITRUST 05.k must cover the EHR vendor's audit-log subsystem in addition to the chain vendor.

### §H6.3 Pattern (c) — chain-replaces-audit-log

Not recommended for ONC-certified EHRs. The chain's wire form does not natively conform to §170.315(d)(2)'s data-element set; substituting the chain for the EHR's audit log would jeopardize ONC certification and break the patient-access surface under 45 CFR 171. Workable only for AI systems operating outside the certified EHR boundary — a standalone clinical-decision-support tool that does not write back to the EHR is the typical fit.

### §H6.4 Information-blocking interaction

Chain entries are NOT patient-accessible under 45 CFR 164.524 by default. They are institution-internal evidence and are not part of the designated record set unless the institution explicitly designates them under 45 CFR 164.501. The HIPAA Privacy Officer makes the designation decision in consultation with the AI Governance Committee. The default posture (chain entries not in the designated record set) is the recommended default; institutions designating chain entries in the designated record set should expect a heavier patient-access workflow and longer DSAR-equivalent fulfillment cycles.

---

## §H7. CMS-0057-F prior-authorization timeliness

CMS's final rule "Advancing Interoperability and Improving Prior Authorization Processes" (CMS-0057-F, effective January 1, 2026) requires Medicare Advantage organizations, state Medicaid managed-care plans, CHIP managed-care plans, and Qualified Health Plans on the federal exchange to make prior-authorization decisions within 7 calendar days for standard requests and 72 hours for expedited requests. The rule also requires implementing a Prior Authorization API conforming to HL7 FHIR.

The chain captures the AI decision; the regulatory clock starts when the prior-authorization request is electronically received by the plan and ends when the decision is electronically returned. The chain's `captured_at` is when the AI evaluated the request; the chain's `received_at` is when the ledger ingested the entry. Neither matches the CMS-clock-start moment, which is when the FHIR request hit the plan's intake endpoint, often seconds or minutes before the AI was invoked.

### §H7.1 Recommended attribute set (informative)

The institution's AI captures alongside the standard chain entry:

- `audit.healthcare.regulatory_clock_start` — RFC 3339 UTC timestamp of the regulatory-clock-starting event (the FHIR API request receipt for CMS-0057-F).
- `audit.healthcare.regulatory_clock_deadline` — RFC 3339 UTC timestamp of the regulatory deadline, computed from clock start plus the rule's window (7 calendar days for standard, 72 hours for expedited).
- `audit.healthcare.regulatory_clock_basis` — the regulation imposing the clock. Examples: `cms-0057-f` (federal), `ncqa-um` (NCQA Utilization Management accreditation), `state-pa-rule:CA` (California prior-authorization rule), `state-pa-rule:NY` (New York). Multiple bases compose; the institution captures all that apply.
- `audit.healthcare.decision_returned_at` — RFC 3339 UTC timestamp of the electronic-return moment (when the FHIR response left the plan's edge).

These are institution-emitted under the canonical bytes per spec §5. They receive the chain's per-event MAC binding. A coordinated rewrite of any of these surfaces as a MAC mismatch at §7 step 9.

### §H7.2 Operational consumers

- **NCQA UM compliance team.** Consumes at monthly cadence for the timeliness-percentile distribution: P50, P95, P99 of `decision_returned_at - regulatory_clock_start`.
- **CMS-0057-F examination response.** When CMS examines the institution's compliance, the institution's response file references chain entries for sampled decisions. CMS's auditor independently verifies the chain entries against the institution's FHIR API logs to confirm timestamps are consistent.
- **Internal audit.** Sample-based testing of regulatory-clock compliance. See `docs/audit-procedures.md` P-50 (added to the audit-procedures-cluster-d-addenda, see §H7.3).

### §H7.3 Audit procedure P-50 (clinical regulatory-timeliness sample)

The institution's auditor samples N entries (N typically 50-100 per audit period) where `regulatory_clock_basis` includes `cms-0057-f`, computes `decision_returned_at - regulatory_clock_start` per entry, and confirms each entry meets the rule's window (7 days standard, 72 hours expedited; the institution's auditor distinguishes the two via a separate sampling plan because the expedited population is small and warrants a higher sampling rate). Failures are escalated to the CMS-0057-F program owner for remediation; chronic failures (>0.5% of sampled entries missing the window) escalate to the AI Governance Committee.

### §H7.4 Composition with §4.4.1 routing and §4.4.2 deployment-intent

The three §4.4 schemas compose orthogonally for prior-authorization AI: routing shows which AI made the call (which model in the routing decision), deployment-intent shows production/canary/A-B context (was this a production decision or a shadow-mode decision), and timeliness shows the call was made on time. The institution can combine all three to answer multi-faceted regulator questions: "for the production-mode standard prior-auth requests routed to our primary clinical-decision-support model, what was the timeliness distribution last quarter."

---

## §H8. HITRUST CSF v11 control mapping

HITRUST CSF v11 is the de-facto healthcare assurance framework. Most US health systems carry HITRUST attestation; many health-tech vendors require HITRUST-certified cloud and SaaS partners. The chain's prescriptions map to multiple HITRUST controls. This section lists the mapping; the standalone HITRUST mapping document at `docs/control-map/HITRUST-CSF-v11-mapping.md` carries the per-control evidence and sample-procedure detail.

### §H8.1 Cryptographic and key-management mapping

| Chain prescription | HITRUST CSF v11 control | Evidence consumed |
|---|---|---|
| §10.5 HSM custody (FIPS 140-2 Level 3) | 06.d (Cryptographic Key Management) | HSM vendor's FIPS validation certificate; HSM operational logs. |
| §10.5 HSM custody — transport security | 09.s (Secure Communications) | TLS configuration for SDK-to-ledger; HSM-host network-isolation evidence. |
| §10.6.1 IKM generation (RNG patterns) | 06.d | RNG-source declaration; institution's CC8.1 entry naming the RNG source. |
| §10.7 software-key adapter exclusion | 06.c (Information Backup) + 10.k (Change Control) | Build-artifact attestation that the production build excludes the software-key adapter; change-management log of any exception. |
| §10.10 key-rotation procedure | 06.d | Key-rotation log; institution's rotation-cadence policy. |

### §H8.2 Audit-logging and verifier mapping

| Chain prescription | HITRUST CSF v11 control | Evidence consumed |
|---|---|---|
| §10.1 reconciliation procedure | 09.aa (Audit Logging) + 12.b (Information Security Reviews) | Reconciliation log; review meeting minutes. |
| §7 verifier procedure | 09.aa + 12.b | Verifier output log; sample-based independent verification. |
| §10.13 evidentiary artifact retention | 06.c + 09.b (Audit Trail Retention) | Retention-policy doc; storage-location evidence. |
| §4.4.1 routing-decision capture | 09.aa | Routing-event chain entries; routing-policy doc. |
| §4.4.2 deployment-intent capture | 09.aa + 10.k | Deployment-intent chain entries; change-management log. |

### §H8.3 Vendor-management mapping

| Chain prescription | HITRUST CSF v11 control | Evidence consumed |
|---|---|---|
| Vendor-hosted controls | 05.k (Service Provider Compliance) | Vendor SOC 2 report; vendor-conformance attestation; MSA. |
| Vendor incident notification | 11.c (Incident Response) | Vendor incident-notification logs; institution's IR-coordination notes. |
| Fourth-party governance | 05.k + 11.c | Subprocessor list; flow-down contract terms; subprocessor SOC reports. |

### §H8.4 Healthcare-specific overlay mapping

| Healthcare overlay section | HITRUST CSF v11 control | Evidence consumed |
|---|---|---|
| §H1 PHI capture posture | 19.b (Information Privacy and Protection — Minimum Necessary) | Posture-selection minute; PHI-handling policy. |
| §H2 minimum-necessary analysis | 19.b | Written analysis; HIPAA Privacy Officer signature. |
| §H3 breach-notification four-factor assessment | 11.c + 19.f (Incident Reporting) | Four-factor assessment template (signed); breach-notification log. |
| §H4 retention floor + state matrix | 06.c (Records Retention) | State-by-state matrix; HSM key-lifecycle commitment evidence. |
| §H5 FDA SaMD evidence capture | 02.b (Authorized Software Use) | PCCP-boundary chain attributes; manufacturer's regulatory-affairs consumption log. |
| §H6 ONC EHR audit-log integration | 09.aa | Pattern selection; integration-mapping doc. |
| §H7 CMS-0057-F timeliness | 09.aa + 19.f | Regulatory-clock chain attributes; NCQA UM consumption log. |
| §H11 24/7 operational requirements | 12.c (Business Continuity and DR) + 11.c | Downtime procedure; DR runbook; ransomware-response runbook. |
| §H12 clinician-override capture | 09.aa | Override chain entries; override-distribution review minutes. |

### §H8.5 Composition with the SOC 2 TSC mapping

An institution carrying both SOC 2 and HITRUST attestations produces the same chain evidence; the two assessment frameworks read it through different control-objective lenses. The TSC-mapping at `docs/control-map/TSC-mapping.md` and the HITRUST mapping at `docs/control-map/HITRUST-CSF-v11-mapping.md` cite the same evidence repository. This is the unified-evidence-collection posture HITRUST itself recommends in its CMMI mapping guidance.

---

## §H9. Patient Safety Work Product privilege under PSQIA

The Patient Safety and Quality Improvement Act (42 USC 299b-21 to 299b-26) creates a federal evidentiary privilege for Patient Safety Work Product (PSWP). Information assembled within the institution's Patient Safety Evaluation System (PSE) for the purpose of reporting to a federally-listed Patient Safety Organization (PSO) is privileged and confidential under 42 CFR Part 3 — not admissible in malpractice litigation, professional disciplinary proceedings, or as evidence in administrative proceedings.

AI-driven patient-safety events — a sepsis-prediction false negative leading to delayed sepsis recognition, a clinical-documentation AI drafting an inaccurate problem list that propagates into the medical record, a prior-authorization AI denial that contributes to a delayed treatment outcome — are reportable patient-safety events. The institution's PSO-reporting workflow assembles evidence about the event for transmission to the PSO; chain entries are exactly the evidence the PSO workflow consumes.

### §H9.1 The two-track posture

The defensible posture separates chain-entries-as-evidence (open) from chain-derived-PSE-analysis (privileged):

1. **Chain entries themselves are NOT inherently PSWP.** They are institution-internal records of AI activity, not records assembled within the PSE for PSO reporting. The chain's evidentiary openness under `docs/legal-disclosure.md` (FRE 901(b)(9), 902(13), 902(14) framing) applies to chain entries in their default state. A discovery request for "all chain entries related to patient X" gets the chain entries.

2. **Chain-derived analysis assembled within the PSE for PSO reporting IS PSWP.** The institution's patient-safety committee analysis of "we observed N false-negatives in sepsis predictions over Q3, here is the root-cause analysis, here are the recommended model-tuning changes" is PSWP because it was assembled within the PSE for PSO transmission. A discovery request for "all PSE analyses derived from chain entries" gets the institution asserting PSQIA privilege.

The institution maintains the separation in PSE policy: the chain ledger is NOT a PSE component; the PSE consumes chain entries as input but the analytical work-product produced within the PSE is the PSWP, not the chain entries themselves. This is settled doctrine for non-AI patient-safety events; the chain composes alongside it without changing the structure.

### §H9.2 Operational implementation

The institution's PSE policy lists explicitly which artifacts are within the PSE (the analytical reports, the root-cause memos, the recommendation documents, the PSO-bound data extracts) and which are not (the chain ledger, the EHR record, the underlying clinical data sources). The PSE administrator owns the policy. The litigation-support team and the PSO contract review the policy annually.

When a patient-safety event triggers PSE analysis: the PSE-analyst pulls the relevant chain entries (using the institution's audit-query role per `docs/audit-procedures.md` independence procedures), copies them into the PSE-analytical workspace, performs the analysis, produces the PSWP-protected report, and transmits to the PSO. The original chain entries remain in the chain ledger, unprivileged. The PSE-resident copy and the analysis are PSWP.

### §H9.3 Discovery responses

When the institution receives a discovery request:

- "All chain entries related to patient X" → produced (chain entries are not PSWP).
- "All PSE analyses derived from chain entries" → privilege asserted under PSQIA.
- "All documents the institution prepared regarding the AI's decision in case Y" → response distinguishes chain entries (produced) from PSWP (privilege asserted with privilege log).

The institution's General Counsel's privilege-log practice covers both PSQIA and attorney-client privilege; the PSE-administrator and the privilege-log preparer coordinate so PSWP claims are accurate.

---

## §H10. State AI laws — Colorado, California, and the matrix

State-level AI legislation has accelerated since 2024. The matrix continues to evolve; this section lists the major shapers as of 2026-05 and the institution's posture.

### §H10.1 Colorado AI Act (SB 24-205)

Effective February 2026 for high-risk AI in employment, with precedent for AI broadly. Imposes consumer-disclosure obligations on developers and deployers of high-risk AI systems. The chain provides the audit-trail substrate; the disclosure language and timing are state-specific.

The institution's HITRUST 12.a-equivalent control description names the Colorado disclosure procedure: when an AI decision affects a Colorado consumer, the institution discloses (a) that AI was used, (b) the AI's decision summary, (c) the consumer's right to request human review. The disclosure timing is at decision-issuance for adverse decisions. The chain captures the disclosure event under the standard chain-entry mechanism (`audit.regulatory.disclosure.colorado_aia` attribute set, institution-defined).

### §H10.2 California SB 1120 (Physicians Make Decisions Act)

Effective January 2025. For utilization-management AI in California, the AI-assisted decision must be reviewed by a licensed California physician before becoming final. The chain's routing-decision capture (§4.4.1) records the AI-assist nature; the institution adds:

- `audit.california.sb1120.physician_review_id` — the reviewing physician's NPI (PHI under §H1.1; hash under §H1.2).
- `audit.california.sb1120.physician_review_at` — RFC 3339 UTC timestamp of the physician review.
- `audit.california.sb1120.physician_concur` — boolean, whether the physician concurred with the AI's recommendation.

Sample-based audit (P-51, see §H10.5) confirms physician review occurred for every AI-assisted UM decision affecting California residents.

### §H10.3 New York pending legislation

Multiple bills in 2025-2026 impose disclosure and bias-audit requirements. The institution monitors via regulatory-affairs newsletter; the chain's existing capture is sufficient to feed any plausible disclosure or bias-audit requirement that emerges.

### §H10.4 Other states

Texas, Florida, Illinois, Washington, Massachusetts, and Connecticut all have AI-related bills in various stages. The institution's regulatory-affairs team maintains a watch-list and updates the state-AI overlay annually or on significant legislative event.

### §H10.5 Audit procedure P-51 (state-AI disclosure and oversight sample)

The institution's auditor samples N AI-decision chain entries affecting residents of states with AI-specific requirements (Colorado AIA, California SB 1120, others as enacted), confirms the required state-specific attributes are present, and confirms the institution's downstream disclosure or oversight workflow operated as required. Failures escalate to the General Counsel and the AI Governance Committee.

### §H10.6 General principle

The chain composes alongside state-specific disclosure, oversight, and admissibility regimes without subsuming them. State medical-board disciplinary proceedings, state attorney-general enforcement actions, and state-court proceedings each have their own evidentiary and admissibility regimes that may or may not compose with the federal-court framing in `docs/legal-disclosure.md`. The institution's litigation-support team validates per state at the moment of dispute.

---

## §H11. Healthcare 24/7 operational requirements

Hospital operations are 24/7 with life-safety implications the FFIEC banking framework does not face. Three operational scenarios stress the chain's availability and degraded-mode posture beyond what spec §10.15 (multi-region resilience) covers directly.

### §H11.1 EHR scheduled downtime

Every certified EHR has a planned-downtime procedure. Epic's planned downtime is typically 4-8 hours quarterly for upgrades; downtime procedures route clinical activity through paper "downtime forms" and a read-only EHR mirror (Epic's Business Continuity Access; similar functions at Cerner/Meditech).

During EHR downtime, AI systems either go offline or operate in degraded mode. The chain is decoupled from the EHR by design. Chain operations continue through EHR planned downtime: chain SDKs continue capture; chain ledgers continue ingest; seal jobs continue. The institution's HITRUST 12.c (Business Continuity and Disaster Recovery) names the chain's downtime-procedure participation.

The IR posture treats planned EHR downtime as a non-anomaly event provided the chain's components remain operational. If chain components are co-located with EHR components and share the downtime window, the institution operates the chain's seal-job under §4.3.1 cadence relaxation — captured events continue to be ingested and chained (the per-event HMAC is independent of the HSM); seal publication is deferred until the maintenance window closes. The 72-hour SHOULD upper bound applies.

### §H11.2 Code blue, mass casualty, disaster

Clinical AI may be invoked under emergency conditions where the institution's risk-tolerance for chain-availability latency is much lower than planned-operation default. The §4.3.1 72-hour seal-publication SHOULD is operationally appropriate for FFIEC-scope deployments. Institutions with elevated incident posture may tighten to a 6-hour or 12-hour publication SLA per the institution's risk-tolerance statement, named in the institution's HITRUST 11.c control description.

In a clinical disaster, the institution's HSM may be physically affected (regional cloud outage, on-premises facility damage). The §10.15 multi-region resilience patterns apply; the institution's IR coordinator activates the documented failover procedure. The healthcare-specific extension is that the clinical-AI consumer (the bedside team) may need a degraded-mode signal: "the AI is operating, but seal publication is deferred; rely on the AI's recommendation knowing the audit trail has a temporary gap that will close within X hours." The institution names the degraded-mode signal protocol in the operator-guide.

### §H11.3 Ransomware

Healthcare ransomware is a 2024-2025 escalation pattern (Change Healthcare February 2024 was the most visible incident). The chain's posture during institution-wide ransomware is critical: the chain is forensic evidence for the post-incident investigation but only if the chain itself was not encrypted by the ransomware.

The chain's defenses against ransomware:

- **Append-only storage.** The ledger storage is append-only by design (per spec §10.3 ledger-aggregation requirements). A ransomware actor encrypting the storage destroys the chain bytes the same way it destroys any other data, but cannot retroactively alter prior chain entries to hide their actions because the daily seals are independent artifacts.
- **Separate key custody.** The IKM lives in HSM custody; the HSM-signing key for daily seals lives in HSM custody. A ransomware actor compromising the institution's general workforce identity does not obtain HSM key access (separate identity, separate authentication, separate audit trail).
- **Immutable WORM retention.** Cross-region replication to a write-once-read-many store at minimum daily cadence, accessible only to a separate IR-team identity. The institution's HITRUST 11.b (Backup) names the artifact backup posture.

The wave-5 close-out added an operator-guide section on ransomware DR. The healthcare-specific extension covers PHI handling during recovery (the recovered chain may contain PHI under §H1.1 posture; the recovery process itself is HIPAA-controlled), HHS-OCR notification timing (a ransomware event affecting PHI is a §164.404 breach unless the institution rebuts under the four-factor assessment), and Joint Commission patient-safety-event notification (if ransomware affects clinical operations, the event is reportable to the Joint Commission via standard patient-safety reporting channels).

### §H11.4 Audit procedure P-52 (healthcare 24/7 readiness)

The institution's auditor confirms (a) the EHR planned-downtime procedure names chain operations and the procedure has been exercised in the audit period, (b) the institution's IR-tabletop exercise covered at least one healthcare-specific scenario (mass casualty or ransomware) and the chain's role in the exercise was documented, (c) the institution's WORM-backup of chain artifacts is current and an IR-team-only access test has been performed in the audit period.

---

## §H12. Clinician-override capture

Joint Commission accreditation standards require systematic monitoring of high-risk medication-ordering systems (MM.05.01.01) and oversight of clinical-decision-support tools (LD.04.01.07). AI-assisted ordering systems and clinical-decision-support tools fall within both. A core requirement is that clinicians can override the AI's recommendation and that the override is captured and reviewed for systematic patterns indicating either the AI is mis-tuned (frequent overrides on a specific recommendation type) or the clinician workflow has a usability gap.

The chain captures decisions; it does not capture the override event explicitly in any normative attribute set. The override has its own discriminator parallel to the routing-decision discriminator. Without a capture schema, institutions emit override evidence under ad hoc attributes and the Joint Commission survey team variably accepts the institution's evidence quality.

### §H12.1 Override entry shape (informative)

The override is a chain entry of its own (not attributes attached to the AI-recommendation entry), parallel to spec §4.4.1's treatment of routing decisions. The override entry's `chain_kind` is `"audit"` (the v1 enumeration is closed and override does not warrant a new value at v1.x). The entry carries:

- `audit.clinical.override.original_recommendation_run_id` — parent-linkage to the AI recommendation being overridden.
- `audit.clinical.override.original_recommendation_seq` — paired with run_id for full linkage.
- `audit.clinical.override.clinician_npi` — the National Provider Identifier of the overriding clinician. PHI under §H1.1; hash-only under §H1.2 recommended.
- `audit.clinical.override.override_action` — one of `accept_with_modification` | `reject_recommendation` | `defer_decision` | `escalate_to_attending`.
- `audit.clinical.override.override_reason_code` — institution-defined reason code from a controlled vocabulary. The vocabulary is part of the institution's HITRUST 09.aa control description.
- `audit.clinical.override.override_reason_freetext_hash` — lowercase hex SHA-256 of any free-text rationale (the rationale itself lives in the EHR under standard EHR controls).
- `audit.clinical.override.timestamp` — RFC 3339 UTC of the override event.

### §H12.2 Operational consumers

- **Joint Commission survey response.** The institution's survey-response file consumes override-event chain entries for systematic-pattern analysis at quarterly cadence.
- **AI Governance Committee.** Reviews override distribution as a model-tuning signal. High override rates on a specific recommendation type are a re-tune trigger.
- **MRM (Model Risk Management) under SR 11-7 framing where applicable.** Override-rate trending is a model-performance signal.

### §H12.3 Audit procedure P-53 (clinician-override capture)

The institution's auditor samples N AI-recommendation chain entries (N typically 50-100), confirms an override entry exists for each entry where the EHR shows a clinician took an override action, and confirms the override entry's `override_reason_code` matches the institution's controlled vocabulary. Failures (override action taken in EHR but no override entry in the chain; reason code outside vocabulary; reason code missing) escalate to the AI Governance Committee.

### §H12.4 Composition

The §4.4.1 routing-decision capture, §4.4.2 deployment-intent capture, and §H12.1 override capture compose orthogonally. Three kinds of meta-decision, each captured as its own chain-entry kind, all under the canonical bytes per spec §5. The institution's auditor and the institution's regulatory consumers can query any combination.

---

## §H13. Deployment authorization checklist

The AI Governance Committee uses this checklist when reviewing a clinical AI deployment seeking authorization. Each item maps to a section of this overlay or to the spec/companion-doc cited.

- [ ] **§H1 PHI capture posture selected.** The deployment names one of (a) full-content, (b) hash-only, (c) limited data set. The HIPAA Privacy Officer and the AI Governance Committee chair have signed.
- [ ] **§H2 minimum-necessary analysis written.** The analysis covers purpose, field inventory, per-field necessity, and review trigger. Signed by the HIPAA Privacy Officer.
- [ ] **§H3 four-factor assessment template installed.** The IR playbook's chain-detected-anomaly scenario for this deployment includes the four-factor template. The on-call HIPAA Privacy Officer is named in the playbook.
- [ ] **§H4 retention floor confirmed.** The deployment's retention floor is named (HIPAA 6-year, longest applicable state floor for the populations served, IKM-availability commitment from the HSM vendor).
- [ ] **§H5 FDA SaMD evidence capture configured (if applicable).** If the deployment is FDA-regulated, the §H5.1 attribute set is captured. The manufacturer's regulatory-affairs team has confirmed consumption.
- [ ] **§H6 ONC EHR audit-log integration pattern selected.** The deployment names one of (a) parallel-stream, (b) composed-log, (c) chain-replaces-audit-log. The institution's HITRUST 09.aa control description carries the pattern.
- [ ] **§H7 CMS-0057-F timeliness capture configured (if applicable).** If the deployment is in CMS-0057-F scope, the §H7.1 attribute set is captured. The NCQA UM team has confirmed consumption.
- [ ] **§H8 HITRUST mapping current.** The institution's HITRUST CSF v11 control description references the chain controls and the healthcare-overlay sections that govern.
- [ ] **§H9 PSWP separation policy current.** The institution's PSE policy distinguishes chain-entries-as-evidence from chain-derived-PSE-analysis. The PSE administrator has signed.
- [ ] **§H10 state AI laws assessed.** The deployment's resident-state distribution has been mapped to applicable state-AI-law disclosure and oversight requirements; required attributes are captured.
- [ ] **§H11 24/7 readiness confirmed.** The deployment's chain operations have been tested through an EHR planned-downtime exercise and a ransomware-tabletop. The WORM-backup is current.
- [ ] **§H12 override capture configured (if applicable).** If the deployment is an AI-assisted ordering system or clinical-decision-support tool, override entries are captured per §H12.1. The Joint Commission survey response references the capture.
- [ ] **AI Governance Committee minute.** The committee's deployment-authorization minute names the deployment, cites this checklist's completion, and assigns operational ownership.

---

## §H14. Per-role roll-up

**Health-system CISO and AI Governance Committee.** This overlay closes the deployment-authorization gap for clinical AI under the chain. The committee uses §H13 as the gate. The chain composes with the institution's existing HITRUST/HIPAA/FDA/ONC/CMS posture without rewriting any of those instruments from scratch. Sign-off authority remains with the committee and the CISO; the chain raises the floor on integrity evidence, it does not replace the governance work that makes clinical AI deployable.

**HIPAA Privacy Officer.** §H1, §H2, §H3 are the Privacy Officer's primary surfaces. The minimum-necessary analysis (§H2) and the four-factor breach assessment template (§H3) are work products the Officer produces and signs. The retention-floor matrix (§H4) is co-owned with General Counsel. The PSE-separation policy (§H9) is co-owned with the PSE administrator.

**HITRUST assessor.** §H8 is the assessor's primary surface. The mapping at `docs/control-map/HITRUST-CSF-v11-mapping.md` carries the per-control evidence and sample-procedure detail; this overlay's §H8 is the orientation document. The assessor reads the institution's CC8.1-equivalent (HITRUST 06.d / 09.aa entries), identifies the chain-related controls, and consumes the corresponding evidence per the mapping table.

**FDA regulatory-affairs lead.** §H5 is the FDA lead's primary surface. The PCCP-boundary chain attributes (§H5.1) are the institution's ongoing capture; consumption at quarterly cadence (§H5.2) is the standard cycle. Forward-scope evolution (§H5.3) tracks FDA's 2025-2026 AI-software workstream.

**ONC compliance lead.** §H6 is the ONC lead's primary surface. The pattern selection (§H6.1-§H6.3) is the deployment-time decision; the information-blocking interaction (§H6.4) is the ongoing posture. The chain's parallel-stream pattern (§H6.1) is the recommended default; institutions adopting composed-log pattern (§H6.2) take on additional EHR-vendor coordination.

**CMS-0057-F program owner.** §H7 is the program owner's primary surface. The recommended attribute set (§H7.1) is captured for every CMS-0057-F-scope decision; consumption at monthly NCQA cadence (§H7.2) feeds the timeliness percentile. Audit procedure P-50 confirms compliance.

**PSE administrator.** §H9 is the PSE administrator's primary surface. The two-track posture (§H9.1) is the policy foundation; operational implementation (§H9.2) is the day-to-day workflow. Discovery responses (§H9.3) coordinate with General Counsel.

**Operations and IR.** §H11 covers operational realities the FFIEC framework does not naturally address. EHR scheduled downtime, code blue / mass casualty, and ransomware compose with the spec's IR playbook. The institution's operator-guide carries the day-to-day runbook detail; this overlay is the regulatory framing.

---

## §H15. Stopping criterion

This overlay covers HIPAA Privacy and Security Rules, the HIPAA Breach Notification Rule, FDA SaMD guidance including 2024 PCCP final guidance, HITRUST CSF v11, ONC §170.315(d)(2) and Cures Act information-blocking, CMS-0057-F prior-authorization timeliness, state AI laws (Colorado AIA, California SB 1120, others as enacted), Joint Commission patient-safety standards (MM.05.01.01, LD.04.01.07), and the Patient Safety and Quality Improvement Act privilege framework. None of the overlay content changes the spec's normative core. All of it composes alongside the existing privacy-by-design.md, audit-procedures.md, incident-response-playbook.md, legal-disclosure.md, operator-guide.md, and control-map artifacts.

The chain itself is sector-neutral. The healthcare regulatory composition is necessarily sector-specific. Other sector overlays (federal-civilian per OMB M-24-10, energy per NERC CIP, education per FERPA + IDEA, defense per CMMC + DoD AI ethics) compose alongside this one without changing the chain's substrate. Institutions deploying the chain across multiple regulated sectors maintain one chain implementation and multiple overlays — one per sector the institution operates in.
