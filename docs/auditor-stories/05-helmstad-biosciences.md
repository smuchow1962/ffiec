# 🧾 Diary of an Audit Day — Helmstad BioSciences

**Engagement:** FDA BIMO Pre-Inspection Readiness — Phase II NSCLC Program
**Client:** Helmstad BioSciences — mid-size oncology biopharma, Cambridge, MA, ~$1.2B revenue, ~380 employees, two Phase III trials in progress, four Phase II, plus a discovery pipeline
**Status:** TesseraSeal live on the AI clinical-trial-eligibility screening tool for 4 months. Everything else is legacy.
**Date:** Tuesday, the week after Atrio
**Audit team lead:** Dawn
**Client liaison:** Dr. Matti Østergaard, VP of Quality and Regulatory Affairs
**Trial in scope:** Phase II second-line therapy in metastatic NSCLC, ~280 enrolled patients across 12 sites in the US and EU

---

### Context

Helmstad BioSciences turned on TesseraSeal four months ago. Not on the QMS. Not on the lab systems. Not on the safety database. They put it on one thing: the AI clinical-trial-eligibility screening tool that decides whether a candidate identified in CRO-supplied data lakes is eligible for the Phase II NSCLC trial. The chain captures every model call, every classification, every confidence score, every reviewer accept-or-reject. The model card hash, the protocol document hash, the de-identified patient hash, the reviewer SSO subject. All of it, sealed.

The chain itself is the v1.0b wire form — per-event HMAC-SHA-256 chain at capture per spec §4.1, daily Merkle seal under RFC 6962 per spec §4.2, HSM-rooted Ed25519 root signature per spec §4.3 with `sign_payload_version = "v1.0b"` binding `key_versions_canon` and `kms_handle_uris_digest` at the day-aggregate level. The reviewer-facing entries carry `chain_kind = "audit"` per spec §3 enumeration; the model-call entries carry `chain_kind = "model_call"` and `gen_ai.request.model` plus `gen_ai.response.model` per spec §4.4 and §7 step 12a. Verifier exit codes follow the §10.12 contract.

Everything else at Helmstad runs on the same plumbing every other mid-size oncology sponsor runs on. Veeva Vault for QMS. A Medidata Rave EDC operated by the CRO. Argus for safety. Salesforce-based CTMS for site management. Email and SharePoint for everyday documents. Two CROs supply the data feeds — a global CRO for the Phase III programs and a specialty oncology CRO for the Phase II work. Where the chain doesn't extend, the regulatory record is built on tooling Helmstad does not own end-to-end. Spec §10.19 calls this surface a chain-coverage map, and Helmstad is supposed to publish one.

Dr. Østergaard knows this. He has prepped for FDA inspections at three prior companies and has read the 483 letters from a half-dozen others. He picked the AI side first because the FDA's draft guidance on AI/ML in drug development told him to, the eligibility tool was the smallest blast radius that mattered, and four months was a realistic timeline. He wants Dawn's team to run the BIMO readiness review like the inspector will run it in six weeks. Find the gaps. Name them. He will take the report into his pre-inspection prep meeting and decide which gaps get closed before the inspector arrives and which get explained.

The Daubert grounding for the chain's evidentiary posture — testability of the §7 verification procedure, peer review of the spec, known error rate of the cryptographic primitives, general acceptance under FIPS standards — is in spec §1.1. The epistemic scope is in §1.2: the chain proves what the AI said and that the record was not tampered with after capture; the chain does NOT prove the AI's statement was clinically accurate, that the statement complied with FDA Quality System Regulation, or that the statement was free of bias. For an FDA Bioresearch Monitoring inspector reading 21 CFR Part 11 against the chain, the §1.2 line is load-bearing — the chain is the integrity foundation, not the truth foundation.

This is the fourth audit Dawn's team has done in four weeks. Last week was a multi-tenant BaaS platform — forty-seven tenants, fourteen hundred verifier runs, zero failures. The week before that was a specialty steel mill in northwest Indiana with three zones and three different posture grades. The week before that was a top-twenty health system that had sealed exactly one inference path and asked the audit to fund the rest. And the week before that was the gold standard.

Today is the third bifurcated audit in four weeks. Different industry. Different inspector. Same shape.

---

### What FDA brings versus what FFIEC brings

The chain-of-custody specification was written under FFIEC working-group authority for AI-driven decisions in regulated banking. The FDA's BIMO inspectorate operates under 21 CFR Part 11 (electronic records and signatures), 21 CFR Part 820 (Quality System Regulation), ICH GCP, ICH E6(R3), and the December 2024 FDA AI/ML draft guidance. The two regulatory regimes overlap on cryptographic discipline — both require integrity-bound electronic records under documented controls — and diverge on procedural language. The chain composes alongside Part 11 the way it composes alongside SR 11-7: same primitives, different regulator, different inspector questions.

For Helmstad, the chain's spec sections that mapped most directly to FDA Part 11 expectations were: §10.5 (FIPS 140-2 Level 3 HSM custody) for Part 11 electronic-signature integrity, §4.1 (per-event HMAC chain) for the record-tamper-detection requirement, §4.3 (HSM-rooted Ed25519 signature with `sign_payload_version` dispatch) for the time-stamping integrity Part 11 §11.10 names, §10.3 (append-only enforcement) for the non-modification requirement, §10.13 (evidentiary artifacts retention) for the documentation Part 11 §11.10(e) requires, §10.17 (HSM partition ceremony attestation) for the dual-control electronic-signature ceremony, and §10.22 (redaction discipline) for the PHI-handling overlay HIPAA Privacy Rule contributes alongside Part 11.

The §10.21 cross-vendor model-handover schema was not directly named in 21 CFR Part 11 — Part 11 predates AI/ML SaMD by decades — but it composed alongside the FDA's December 2024 PCCP guidance and the Software-as-a-Medical-Device framework. The §10.21 attribute family let Helmstad's regulatory-affairs team prove the model-card-and-validation-report lineage that FDA's draft AI/ML guidance asks sponsors to produce. The §10.20 training-data retention discipline composed alongside FDA's post-market-surveillance expectations. The two together — §10.20 plus §10.21 — gave Helmstad an FDA-defensible model-lifecycle-and-retention posture that mapped cleanly onto the inspector's Part 11 framework even though the spec did not name FDA explicitly.

Where the chain reached, the integrity foundation came from FFIEC-grade cryptographic primitives. Where the chain did not reach — the legacy stack the BIMO inspector would also test — the institution's HITRUST-equivalent attestation regime and the trial's GCP / GMP processes carried the evidence load. The §10.19 chain-coverage map made the boundary between the two regimes explicit. That was the engagement's structural framing.

---

### Audit Team

| Name | Role |
|---|---|
| Dawn | Lead Auditor — governance and narrative |
| Raj | Database specialist |
| Elena | CRM systems |
| Mike | Application / API layer |
| Diana | IAM and access control |
| Luis | DevOps / logs / pipelines |
| Chen | Data engineering / ETL |
| Tom | Internal-audit liaison specialist (visiting team; partners with the client CAE) |

---

### 🌅 8:30 AM — Kickoff

The drive in was twenty-five minutes from the hotel through Cambridge morning traffic. Dawn had her coffee in the cup holder and the engagement brief on her tablet.

Dawn had stopped expecting another Northbridge. The cleanest engagement she had run in years was now five engagements back, and the team had quietly retired the assumption that it would repeat. Helmstad layered FDA Part 11 strict, clinical-trial-duration retention, and a CRO-supplied predictive model on top of a chain primitive the team had been working with for a month. The expectation on the drive in was *this will not be Northbridge*.

*Four bifurcated audits in five weeks*, she thought. *Different industries. Same architecture.*

Northbridge had been the gold standard. TesseraSeal everywhere — every credit decision, every wire, every IAM change, every ETL job. The team had spent four days trying to find a gap and found a stale comment in a YAML file. Dawn had driven home that day with the report half-written in her head and the rest of it dictated into her phone.

Mercator had been the bifurcation introduced. AI imaging on the chain. Claims, billing, and the EHR off the chain. Patricia Okonkwo had asked the audit to draw the line cleanly so she could fund the next phase. The CAE had asked Dawn to commit to language that would survive the board read. Dawn had given him the language and watched him write it down word for word.

Stelvio had been the three-zone version. AI vision QC on the chain. OT on the floor. IT business systems in the back office. Maria Costanza had wanted the report to triage the gaps so the CFO could pick which zone to fund first. The audit had made the case for the OT zone. Maria had taken it to the CFO Friday.

Atrio had been last week. A mid-size BaaS platform with forty-seven fintech-tenant programs running on the same TesseraSeal deployment. Multi-tenant chain with per-tenant Merkle roots and a cross-tenant exposure index that the audit had stressed-tested for fourteen hundred verifier runs across two-hundred-and-eighty days of randomly selected dates and not produced a single FAIL. Dawn had driven home from Atrio with a quiet confidence she wasn't sure she had felt before. The platform-vendor pattern wasn't theoretical anymore. It worked.

Today was Helmstad. Mid-size oncology biopharma. AI eligibility tool sealed. Everything else legacy. FDA BIMO inspection in six weeks.

*It never is*, Dawn thought. *But sometimes the part that's chained is the part that decides, and that's enough — until the FDA inspector asks how the inputs got to the chain.*

She elaborated to herself as she pulled into the visitor lot off Kendall Square. The Mercator pattern was the closest precedent. Different regulator, different stakes, but the same architectural shape: AI side sealed, source side legacy, boundary at the vendor handoff. The difference was that the FDA inspector who would walk Helmstad in six weeks was a Bioresearch Monitoring inspector with a checklist that included 21 CFR Part 11, ALCOA+, ICH E6(R3), and the new AI/ML draft guidance. The HIPAA auditor had been the regulator at Mercator. The FDA inspector at Helmstad would ask sharper questions about source-of-truth. *Plaintiff bar versus FDA inspector*, she thought. *Different motivations, similar reads.*

The §1.1 Daubert lens applied to FDA hearings the same way it applied to federal-court authentication: testability of the procedure, peer review of the standard, known error rate of the primitives, general acceptance of the cryptographic foundation. The chain's response to each was named in the spec — §7 byte-exact procedure for testability, the public test-vector corpus and Apache-2.0 reference implementation for peer review, §1.3 EUF-CMA / second-preimage / EUF-CMA composition for known error rate, FIPS 180-4 / FIPS 186-5 / FIPS 198-1 for general acceptance. An FDA inspector raising a Part 11 challenge would not name Daubert by name, but the structure of their question — can the integrity claim be falsified, has it been peer-reviewed, what is the failure rate, is the construction generally accepted — was the same four-question structure. The chain was built so the answer to each question came from the shipped artifacts, not from Helmstad's internal documents.

The lobby was small. Glass-walled conference rooms on the right. A receptionist who recognized Tom's name on the calendar and waved them through.

The conference room was on the fourth floor with a view of the Charles River. Dr. Østergaard was already there with two of his direct reports and a single sheet of paper at each chair. No deck. The sheet was a system map: a green box on the left labeled "TesseraSeal scope: nsclc-phase2 eligibility classifier," a red box on the right labeled "Legacy scope: Veeva QMS, Argus, CTMS, EDC (CRO), CRO data feeds, lab/LIMS, email/SharePoint." A dotted line down the middle. Two arrows crossing the line — one labeled "CRO ingestion," one labeled "EDC extract."

"Good morning," Dr. Østergaard said. He had a soft Danish accent and a careful, deliberate way of speaking. "Before we start. The FDA inspectors arrive in six weeks. I have prepped for inspections at three prior companies. I have seen what happens when the AI side is good and the source side is not. Last year a mid-cap sponsor I will not name received a 483 because their AI eligibility tool was excellent and their CRO data feed was a black box. Two months ago another sponsor — a company called Vellisar Therapeutics, hypothetically — received a 483 because their EDC audit trail was on the CRO's SOC 2 and the inspector wanted Helmstad's word for it, not the CRO's. And six months ago there was a 483 over an Argus audit-trail-table being mutable by DBAs. I have read all three."

He tapped the system map.

"This is what I want from you. The AI side I am reasonably confident about — it has been live for four months and the engineers have been disciplined. The legacy side I want you to look at the way the FDA inspector will look at it. Find what they will find. Tell me what to fix in the next six weeks and what to explain. I would rather hear it from you in May than from the inspector in June."

Dawn put her coffee down. "Thank you for being direct. It saves us a day."

Tom — the visiting team's internal-audit liaison — had been on a call with Helmstad's Chief Quality Officer the day before. He nodded. "We agreed yesterday on the bifurcation framing. The CQO is supportive. He wants the report to come out as one assessment with the AI-side and legacy-side findings drawn in separate sections so the inspector can see the boundary."

"Yes," Dr. Østergaard said. "That is what I want. Two sections. One report."

Dawn looked around the table at her team. "Okay. Morning is the AI side. Mike and Chen on point — Helmstad's engineers will walk you through the eligibility classifier. Diana, you will do the IAM split — both sides of the line. The line is what we are mapping. Afternoon is legacy. Raj on Argus and the lab database. Elena on the Salesforce CTMS. Luis on the CRO ingestion pipeline. Chen on the EDC extract handoff. We reconvene at three for the reconciliation test. Five-thirty debrief."

Dr. Østergaard nodded along. "The eligibility-tool team is expecting Mike and Chen at nine. Dr. Hannah Reisch — the lead clinical informaticist — built most of the chain integration herself with the platform team. She will not waste your time."

"Good," Mike said.

"You will like her."

Elena had been listening quietly. "Dr. Østergaard. Quick orienting question. The Salesforce CTMS — is it just for trial-site management or does it also hold KOL relationships?"

"Both. The CTMS is a Salesforce instance. Site management for the active trials. The commercial team also uses it as a CRM for KOL outreach. Two business lines. None of it is in the chain. It is the same Salesforce setup we had before any of the AI work. I want you to look at it because I want it documented in the report alongside everything else."

"That is clear."

Luis, half to himself: "Site monitoring notes and KOL outreach in the same Salesforce footprint. That is going to be interesting."

Dr. Østergaard smiled faintly. "It is interesting. It is also typical. Mid-size sponsors do this everywhere."

He stood up and gestured at the door. "Mike, Chen — let me walk you to the engineering floor. Dr. Reisch is expecting you."

---

### 🧩 9:15 AM — Walking the Eligibility Classifier in Production

Mike had expected the AI walkthrough to follow the usual pharma pattern. Engineering team puts up slides. Talks about "validated systems." Shows a quality-managed Splunk dashboard. Hand-waves around the boundary where the chain stops.

It did not go that way.

Dr. Hannah Reisch was a former oncology research associate who had moved into clinical informatics six years ago and into AI tooling three years after that. She had a laptop open and one terminal window. No deck.

"You want to see a sealed eligibility decision," she said. "Pick a date. Pick a site."

Mike picked March 22. Site 04 — Mass General Cancer Center.

Dr. Reisch typed for ten seconds. "Okay. March 22, Site 04. The classifier ran 87 candidates that day. 14 came back eligible, 38 ineligible, 35 human-review-required. Pick one."

"Pick the one with the highest confidence score that came back eligible and was accepted by a reviewer."

She typed again. The terminal showed a JSON entry — structured, with fields Mike recognized. `model_id`, `model_version`, `patient_hash`, `eligibility_classification`, `confidence`, `criteria_doc_hash`, `reviewer_subject`, `reviewer_decision`, `reviewer_reason_code`. Each field had a hash next to it. A tenant binding. A sequence number. A signature reference.

"This is the entry for decision `nsclc-2026-03-22-mgh-00041`," Dr. Reisch said. "De-identified patient — the hash is HIPAA Safe-Harbor — was assessed against protocol version 4.2 of the NSCLC eligibility criteria. Model returned `eligible` at 0.93 confidence. The CRC reviewed and accepted at 14:18 local. Per spec §4.4, the entry carries `chain_kind = 'model_call'`, `gen_ai.request.model`, `gen_ai.response.model`, `ffiec.chain.payload_hash` as 32 raw bytes, `ffiec.chain.key_fingerprint` as 16 raw bytes per spec §3, and `ffiec.chain.captured_at` with nanosecond precision. Here." She rotated the screen.

```
{
  "entry_id": "nsclc-2026-03-22-mgh-00041",
  "tenant": "helmstad-trial-screen",
  "service": "nsclc-phase2",
  "seq": 184729,
  "ts": "2026-03-22T18:14:03.211Z",
  "model_id": "eligibility-classifier-v3.7",
  "model_version": "sha256:8f2a91...c4d3",
  "patient_hash": "sha256:a91f8b...ee27",
  "criteria_doc_hash": "sha256:protocol-4.2:b71c...92a4",
  "site_id": "site-04-mgh",
  "classification": "eligible",
  "confidence": 0.93,
  "reviewer_subject": "h.tan@partners.org",
  "reviewer_decision": "accept",
  "reviewer_reason_code": "criteria-match-confirmed",
  "reviewer_ts": "2026-03-22T18:18:47.092Z",
  "prev_hmac": "sha256:...",
  "this_hmac": "sha256:...",
  "merkle_path": [...],
  "seal_ref": "nsclc-phase2-2026-03-22-eod"
}
```

"Day-boundary semantics?" Mike asked. "Per spec §4.2.2 the day partition uses `received_at`, not `captured_at`. How does that interact with late-arriving events?"

"`received_at` per §4.2.2 is stamped by the ledger server upon ingest, not by the SDK. We log clock-skew detection events when `|received_at − captured_at|` exceeds five minutes per §4.2.2's configurable threshold. Application hosts are NTP-synchronized per §10.4. If a chain entry arrives at the ledger after the daily seal for its `received_at` UTC date is sealed, we stamp `ffiec.chain.late_binding = true` per §4.2.2 normative requirement and include the entry in the next day's seal — the original day's seal MUST NOT be altered. The verifier reports late-binding entries explicitly as `late-binding entries: N` under `Status: PASS` per §7. We've had four late-binding entries in 124 days — all from one Friday-evening site network outage in February. The verifier's working-paper output shows them. None are integrity violations; all are normal-operations behavior."

"And every tenant-day produces a seal record, even with zero events?"

"Every tenant-day. Per §4.2 empty-day seal continuity. The Merkle root for an empty day is `SHA-256(b"")` per RFC 6962 §2.1 — the constant `e3b0c44...b855`. The verifier reproduces it mechanically. A missing empty-day seal would surface as `missing seal for tenant-day {D}` per §4.2 control-completeness failure. We've had seven empty days in 124 — pre-trial-launch days when no candidates were screened. All sealed."

"Run the verifier on it," Mike said.

Dr. Reisch typed:

```
herald-verify --tenant=helmstad-trial-screen --service=nsclc-phase2 \
  --date=2026-03-22 --entry-id=nsclc-2026-03-22-mgh-00041
```

The terminal hesitated for four seconds and then printed:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key nsclc-phase2-2026-q1
```

The output format matched the spec §7 normative form — `Status:` / `Step:` / `Reason:` with single-space separators, line terminator `0x0A`. Exit code 0 per spec §10.12. The verifier had executed all twelve §7 steps in order: format-version pre-flight, HKDF-inputs digest, genesis-hash, tenant-id character class per §3, per-entry binding, per-entry format, structural walk, IKM lookup, fingerprint check, MAC recompute (using `expected_prev_hash` from the structural walk per §4.1 inviolate property 8, not the entry's claimed `prev_hash`), Merkle recomputation under RFC 6962, signature verification with `sign_payload_version = "v1.0b"` dispatch, and the §7 step 11 `key_versions` cross-check. The §7 step 12a check on `gen_ai.request.model` and `gen_ai.response.model` had also fired and passed.

Mike leaned back.

"Now pick one from ninety days ago," he said.

Dr. Reisch picked an entry from December 22, 2025 — same site, same protocol version. Ran the verifier. Same four seconds. Same PASS. Same twelve steps. That entry had been sealed under `sign_payload_version = "v1.0a"` — pre-Round-17 close-out — and the v1.0b verifier had dispatched on the seal record's version field per spec §7 step 11 to reconstruct the 10-line form. Monotonic dispatch was working as documented.

Mike asked one more question. "DR posture. If your SDK loses local persistence in the middle of a run — disk corruption, container crash without state replication — what happens to chain continuity?"

Dr. Reisch answered without pausing. "Spec §10.25 run-resume contract. The SDK has three places to look for the chain tail — in-memory state, local persistence sidecar, and the ledger's chain-tail endpoint. We use all three in priority order. The sidecar is a per-`(tenant_id, run_id)` SQLite row holding `(latest_seq, latest_payload_hash, key_version, last_commit_utc)`. We file-lock the row so a single writer per run is enforced — that's the §10.25 single-writer-per-run rule. If the sidecar is gone, the SDK queries the ledger's chain-tail endpoint and rejoins from the returned tail. If the ledger is unreachable AND local persistence is missing, the SDK refuses to emit until the tail is acquired or an operator authorizes a fresh genesis under a NEW `run_id` per §10.25 DR rejoin discipline. We don't silently re-genesis under the same `run_id`. The §4.4 genesis-block-uniqueness rule plus the ledger ingestion cross-check at §10.25 catches anything that slips through."

"What about adverse-action-style notification — when a candidate is screened out and informed?"

"Spec §10.11 by analogy. The §10.11 schema is normated for ECOA adverse-action notice translation, but the spec body explicitly applies the discipline to state-insurance-law adverse-action notices and other regulator-named declination scenarios. For clinical-trial screening, the FDA does not require notification to declined candidates by spec — the trial protocol governs — but Helmstad's CC8.1 names the candidate-notification timing as an institution-side process and the chain entry for the decline carries `audit.ecoa.translation.target_language` and `audit.ecoa.translation.output_hash` per §10.11 schema, treating the §10.11 attribute set as the chain-of-custody record of the institution's notification commitment. The §10.11.2 FCRA reinvestigation timing pattern doesn't directly apply, but the architectural shape — a parent decision plus a downstream notification with bound timestamps — is the same. We use §10.11 as a reusable schema where the institution-side notification timing matters."

"And entity succession? If Helmstad is acquired."

"Spec §10.24. We'd emit a `chain.entity_succession` operational event per §10.2 with `from_entity_legal_name`, `to_entity_legal_name`, LEI fields per RFC 9101, `effective_utc`, `kind` (e.g., `acquisition`), the regulator filing identifier when applicable, and the dual-signatures array per §10.17 schema with `entity_affiliation` per Round-17 M&A-P1. The chain stays under the same `(tenant_id, run_id)` keying unless the acquirer explicitly renames `tenant_id`. Chain entries before the succession remain verifiable under the original entity's binding; entries after the succession are bound under the acquirer's signature on the transfer-day seal. We don't break chain continuity — we record the legal-entity change. That's relevant for FDA inspections in the M&A context because a sponsor might be acquired between an inspection date and a 483 response."

Mike wrote both down. *§10.25 DR posture — three-place tail acquisition + single-writer enforcement + DR refuse-to-emit. §10.24 succession — chain stays continuous; legal-entity change recorded. Same architectural pattern as Northbridge described.*

> **✓ Confirmation #1 — Chain integrity holds at four months and ninety days, across the v1.0a → v1.0b dispatch boundary**
> The eligibility classifier has been emitting sealed entries for 124 days. The verifier resolves a recent entry (sealed under v1.0b 12-line `sign_payload`) in 4 seconds and a 90-day-old entry (sealed under v1.0a 10-line `sign_payload`) in 4 seconds. Twelve verification steps per spec §7 including HMAC recomputation under HMAC-SHA-256 per §4.1, RFC 6962 Merkle path resolution per §4.2 (leaf-prefix `0x00`, internal-prefix `0x01`), and Ed25519 strict-canonical signature verification per spec §4.3 against the published quarterly public key. The verifier dispatches on `sign_payload_version` per §7 step 11 — pre-amendment chains under their original 6-line form, v1.0a chains under their 10-line form, v1.0b chains under their 12-line form. The chain endures across the seal boundary and across the wire-form amendment.

Mike asked: "What is sealing the seal?"

"Daily Ed25519 signature on AWS CloudHSM in `us-east-1`. CloudHSM Classic — FIPS 140-2 Level 3 per spec §10.5. We chose CloudHSM after our SOC 2 engagement said an on-prem HSM was overkill for our footprint. The IQ/OQ documentation for the CloudHSM configuration is in Veeva — I can show you if you want. The FDA accepts cloud HSM for Part 11 provided the sponsor's IQ/OQ is documented. Ours is. The seal-job operator role grants `sign` only — `extract`, `delete`, and `import` require separate authorization per §10.5 separation of duties. Spec §10.5 fault-injection-attack residual-risk posture is documented in our CC8.1 — FIPS 140-2 Level 3 does not certify DFA resistance; we accept FIA as a residual risk under the v1.0 baseline and our incident-response procedure covers detection and recovery if a fault-injection attack succeeds in the field. No nation-state-attacker threat model in our scope; if we expand into geopolitically sensitive regions later, we'd revisit per §10.5's institution-side guidance on Common Criteria EAL5+ or FIPS 140-3 Level 4."

"And the IKM generation?"

"Inside CloudHSM. HSM internal RNG per spec §10.6.1 — first of the three conformant patterns. The IKM never leaves the HSM in cleartext. We're running Model B from §4.1.1 — session-key-delivered, with HSM-resident PRK and SDK-side Expand. The PRK transits the SDK process briefly under the same memory-protection posture as Model A's IKM. We document the PRK-handling posture in CC8.1. The IKM is 32 bytes — minimum per §10.6 — and never enters application memory."

"Show me the IQ/OQ."

She pulled up the IQ/OQ document in Veeva. It was 47 pages. Configuration captures, key-attribute definitions, partition isolation, cross-region failover behavior, key-ceremony minutes from January 2026. Dr. Østergaard's signature on page 47. The minutes named the signatories with `role`, `name`, and `entity_affiliation` per spec §10.17 schema — the `entity_affiliation` field had been added in Round-17 M&A-P1, and Helmstad had retrofitted it three weeks ago.

> **✓ Confirmation #2 — IQ/OQ for CloudHSM signing infrastructure is documented to Part 11 standard, and the partition ceremony is chain-coupled per §10.17**
> CloudHSM configuration for the daily Ed25519 seal is documented in a 47-page IQ/OQ in the Veeva QMS. The HSM is FIPS 140-2 Level 3 per spec §10.5 (CloudHSM Classic; AWS KMS multi-tenant managed-key service is explicitly NOT conformant). Key ceremony minutes are recorded with two-of-three approval (Dr. Østergaard, the platform engineering lead, and the SecOps lead) under the §10.5 separation-of-duties posture. The document is signed and dated. The configuration is reproducible from the IQ/OQ. The chain emits `chain.partition_ceremony_attended` per spec §10.17 for partition creation, IKM rotation, and partition-PIN reset, with `signatories[].entity_affiliation` per the Round-17 M&A-P1 close-out. The paper-and-PDF attendance log remains the dispute-resolution record for ink-signed authenticity per §10.17; the chain event is the integrity-bound attestation that the ceremony occurred. The two records together form the §10.17 evidence pair the FDA accepts for Part 11 electronic-signature integrity.

Chen had been quiet, watching the chain entries scroll. He spoke up.

"Show me the model-card binding. And the cross-vendor handover — who supplied the weights?"

Dr. Reisch pulled up the model registry. Each model version had a SHA-256 of the weights file, a SHA-256 of the model card document, a SHA-256 of the validation report, and the date the version was promoted to production. The chain entry's `model_version` field was a hash of that bundle. Below the bundle were three more attributes — `audit.model_handover.provider`, `audit.model_handover.model_artifact_sha256`, `audit.model_handover.model_card_sha256` — per spec §10.21.

"The model came from our internal-ML team — `audit.model_handover.provider = 'helmstad-internal-ml'`. The handover entry on January 14 carries `model_version`, `model_artifact_sha256`, `model_card_sha256`, and `fairness_audit_report_sha256` per §10.21. We don't carry `contract_id` because it's an in-house handover under a single corporate authority — the §10.21 contract-binding clause permits omission with CC8.1 documentation, which is in Veeva. The training-data retention floor is 24 months under our DPIA, longer than the model's 18-month deployment window plus 90-day investigation buffer per §10.20. We're conformant on §10.20 floor."

"Every chain entry's `model_version` references a specific tuple of weights, model card, and validation report. If any of the three changes, the hash changes. We can prove which model card was in force the day a decision was made."

> **✓ Confirmation #3 — Model card, validation report, and provider are hash-bound to every chain entry under §10.21 and §10.20**
> The `model_version` field in each chain entry is a SHA-256 over a tuple of (model weights, model card document, validation report). Any of the three changes, the hash changes. The §10.21 `audit.model_handover.*` family captures the provider identity, the artifact hash, the model-card hash, and the fairness-audit-report hash on the handover entry; downstream model-call entries reference the handover by `parent_run_id` / `parent_seq` per §4.4. The §10.20 training-data retention floor is set at the deployment window plus a 90-day investigation buffer, with the GDPR Article 5(1)(c) data-minimization tension resolved through Article 6(1)(f) legitimate interests tied to FDA Quality System submissions. The FDA's draft AI/ML guidance asks sponsors to prove which model card governed a specific decision; Helmstad answers from the chain alone, with §10.21 cross-anchor support if a future external provider is introduced.

Mike asked: "Deployment intent? You're running a single model in production today. But if you A/B-tested two models, or ran a validation environment alongside production, how would the chain disambiguate?"

Dr. Reisch nodded. "Spec §4.4.2 deployment-intent capture. Today every model-call entry carries `audit.deployment.intent = 'production'` and `audit.deployment.policy_version`. For Helmstad, that's effectively a no-op — single model, single region, no canary, no A/B test. But the schema is in place. If we ever run a clinical-decisioning validation environment alongside production — call it `validation` per institution-named value documented in CC8.1 — the chain would carry `audit.deployment.intent = 'validation'` on those entries and the SOC team's per-model decision-count distribution analysis under audit-procedures P-26 would disambiguate validation runs from real-applicant runs. The §4.4.2 enum has `regulatory_sandbox` for sandbox programs and `disparate_impact_test_run` for quarterly DI testing — neither directly maps to FDA today, but the architectural shape is there for sponsors who run an FDA-Sponsored sandbox under the AI Safety Institute's pilot programs. We'd document the institution-named value in CC8.1 alongside the FDA enum."

"And the protocol — the eligibility criteria document?"

"Same pattern. The `criteria_doc_hash` references the protocol version and amendment in force on the decision date. Protocol 4.2 was promoted on January 14, 2026. Every entry from January 14 forward references the 4.2 hash. If we promote 4.3, the hash changes the moment the new criteria are in force."

Mike wrote that down. *That is the ALCOA+ Original attribute on the inputs side.*

> **✓ Confirmation #4 — Protocol document hash binds the eligibility criteria to the decision**
> Every chain entry references the SHA-256 of the protocol document version and amendment in force on the decision date. The classifier cannot decide against a version of the criteria that is not in the registry. The criteria document itself is stored in Veeva with version control. The hash is the bridge.

"Reviewer decisions," Chen said. "Walk me through the human-in-the-loop capture."

Dr. Reisch pulled up a record where the classifier had returned `human-review-required` at 0.62 confidence. The reviewer — a clinical research coordinator at the Memorial Sloan Kettering site — had reviewed the candidate's de-identified profile, made a manual determination of `ineligible`, and entered a reason code from a controlled vocabulary: `prior-systemic-therapy-exclusion`. The reviewer's SSO subject was captured. The decision timestamp was captured. The reason code was captured. The free-text justification — also captured, with a hash, into the chain.

"The reviewer authenticates with SSO," Dr. Reisch said. "The SSO subject becomes part of the chain entry. The reason code is from a controlled vocabulary. The free-text justification is captured — but redacted at the SDK boundary per §10.22, not post-MAC. The captured JSON IS the redacted form. We carry `audit.redaction.policy_id`, `audit.redaction.policy_version`, `audit.redaction.redacted_field_paths`, `audit.redaction.redaction_method`, and `audit.redaction.disposition = 'redacted_at_sdk'` so a future CFPB-style examiner — or, in our case, an FDA inspector — can confirm the redaction posture mechanically. The hash of the unredacted text lives in our compliance vault under court-controlled IKM escrow. All of it is sealed. We can prove who decided what at what time and why."

> **✓ Confirmation #5 — Human reviewer decisions are sealed end to end with §10.22 redaction discipline**
> Clinical research coordinators authenticate via SSO. Their accept/reject decision, the reason code from the controlled vocabulary, the free-text justification (redacted pre-MAC at the SDK boundary per §10.22 normative posture), and the timestamp are all captured into the chain entry. The reviewer's SSO subject is the attribution. The `audit.redaction.*` family per §10.22 is bound under the per-event MAC so the FDA inspector can read what was redacted, where, and how. The chain shows who decided what, when, and why — with the redaction discipline named on the chain itself rather than inferred from architecture.

Mike looked at Chen. Chen looked at Mike.

"That is Part 11 audit-trail discipline," Mike said. "On the AI side."

"On the AI side," Chen agreed.

Mike thought a moment longer. "And the §4.1.3 per-event MAC algorithm-agility posture? You're a long-retention sponsor — 25 years of trial-record retention is a long horizon for HMAC-SHA-256."

Dr. Reisch nodded. "We're emitting `payload_hash_alt` per spec §4.1.3 — RECOMMENDED at v1.0b, candidate-normative for v1.x. Second-algorithm MAC over the same canonical bytes. We chose HMAC-SHA-3-256 as the alt algorithm; it's documented in CC8.1 as part of our cryptographic-agility posture. The verifier under v1.0b reports `payload_hash_alt: PASS-INFORMATIVE` when both MACs verify — informative-only at v1.0b per §4.1.3 informative posture. If a SHA-256 break is announced, the §4.3.2 emergency-patch SLA gives the spec working group 30 days to publish; institutions migrate within 90 days for HMAC/SHA-256 breaks. We have a 25-year retention horizon so we're betting on the alt MAC carrying us through any near-term cryptographic surprises. The §4.1.3 candidate-normative AND-security posture is what we'd shift to if a SHA-256 break required the alt MAC to become the primary integrity evidence."

Mike wrote that down. *§4.1.3 — Helmstad emits `payload_hash_alt` as informative safety margin. Long-retention sponsor reasoning. Same posture pattern that NetiVa and Atrio used for the same reason.*

> **✓ Confirmation #5a — `payload_hash_alt` per §4.1.3 emitted as v1.0b RECOMMENDED safety margin**
> Helmstad emits the per-entry `payload_hash_alt` second-algorithm MAC per spec §4.1.3 (v1.0b RECOMMENDED, candidate-normative for v1.x). Second algorithm is HMAC-SHA-3-256 documented in CC8.1. Verifier reports `payload_hash_alt: PASS-INFORMATIVE` on each entry. If a SHA-256 break is announced, the §4.3.2 30-day emergency-patch SLA plus the 90-day institutional migration window gives Helmstad a graceful migration path; the alt MAC is the in-band early-warning capability. 25-year retention horizon makes the safety margin worth the per-entry compute and storage cost.

Mike turned back to Dr. Reisch. "One more — best-evidence posture under FRE 1001-1004. Spec §5.2."

Dr. Reisch nodded. "We produce both forms in any litigation hold. The captured JSON is the content-bearing form — the human-readable record of the inference. The canonical bytes per RFC 8785 JCS are the integrity-bearing form — what the MAC covers. Both are originals under FRE 1001(d). FRE 1003 admits duplicates of either form unless authenticity is challenged, and the §7 verification procedure is the procedural answer to such a challenge. We coded the discovery production around that posture from day one."

Mike wrote that down. *§5.2 — captured JSON for content, canonical bytes for integrity. Both originals. Same framing as Northbridge.*

Mike turned to a different angle. "OTLP transport identification per spec §4.4.3. Walk me through the Resource attributes."

Dr. Reisch typed a query against the OTLP receiver's traffic logs. The Resource block came up:

```
service.name           = "nsclc-phase2-classifier"
service.version        = "3.7.2"
ffiec.chain.spec       = "v1.0"
ffiec.chain.posture    = "ffiec"
ffiec.chain.format_version = "v1"
```

"All five required Resource attributes per §4.4.3. The OTLP traffic is dispatched on `ffiec.chain.spec` at the receiver — traffic without that Resource attribute is regular OTel telemetry and is NOT routed to the chain-of-custody pipeline. Traffic with the attribute IS routed to the chain pipeline. The dispatch decision is integrity-bound through §5 because the OTLP Resource is part of the canonical bytes — a tampered Resource attribute would surface as a MAC mismatch at §7 step 9. We use OTLP/HTTP for the production export with `X-FFIEC-Chain-Spec: v1.0` and `X-FFIEC-Chain-Posture: ffiec` HTTP headers per §4.4.3 recommended discipline; the headers and the Resource attribute MUST agree, and the receiver cross-checks them on body decode."

"And severity for chain-of-custody traffic per §4.4.4?"

"Receiver stamps `SeverityNumber` in the 9..20 range per §4.4.4 normative requirement. Our QuickLogBuilder analog positions the value dynamically per ingest. Stamped `SeverityText = 'OTLP'` matching the TesseraSeal reference convention. SDK-side severity is `SEVERITY_NUMBER_UNSPECIFIED = 0` because the chain entries ship as OTLP `Span` records, not `LogRecord`, so SDK-side severity is moot. Collector pass-through per §4.4 prohibits dropping or downsampling chain traffic; severity-floor filters in our collector configuration explicitly exempt chain traffic. We tested the exemption with a synthetic-canary control during commissioning."

Mike wrote that down. *§4.4.3 + §4.4.4 — Resource attributes, header agreement, receiver-side severity stamping in 9..20, collector pass-through verified. All to spec.*

---

### 🧠 10:00 AM — The Pipeline That Feeds the Classifier

Chen had been waiting for this part. The chain at the model boundary is one thing. The chain at the input boundary is the other thing. They are not the same thing.

Helmstad's clinical-informatics platform engineer was a person named Devansh Ramaswamy. Eight years at Helmstad, six of them on data-pipeline engineering. He had a Jupyter notebook open with the pipeline lineage rendered as a DAG.

"Walk me through the input data for the entry Dr. Reisch just pulled," Chen said.

Devansh pulled up the lineage for `nsclc-2026-03-22-mgh-00041`. The candidate's de-identified profile had been ingested at 06:14 UTC that morning from a SFTP delivery from Quintessa Research — Helmstad's global CRO. The SFTP delivery was a tarball of de-identified candidate records. Each tarball was PGP-signed by Quintessa's signing key. The Helmstad ingestion service verified the PGP signature, computed a SHA-256 of the tarball, and recorded the hash as a referenced artifact in a chain entry under `chain_kind = "audit"` (per spec §3 enumeration) carrying the `audit.external_artifact.*` family per spec §10.19 — `kind = "cro_sftp_tarball"`, `identifier` = the Quintessa delivery ID, `sha256` = the tarball hash, `received_at_utc`, `source_party = "quintessa-research"`, `evidentiary_role = "chain_of_custody_handoff"`. The candidate records were then unpacked, parsed, and written into the inference-input warehouse. The inference-input warehouse row carried a foreign key to the ingestion chain entry.

"So when the classifier ran on this candidate at 18:14 UTC," Chen said, "the input had been in the warehouse for twelve hours, and the warehouse row was bound by foreign key to a chain entry that references the SHA-256 of the tarball Quintessa sent us this morning."

"Yes."

"Show me the ingestion chain entry."

Devansh did. The entry recorded the SFTP source, the PGP signature verification result, the tarball SHA-256, the file count, the byte count, and the ingestion-service identity. All sealed.

"Good," Chen said. "Now — what about Quintessa's side?"

Devansh's expression changed slightly. Not unhappy. Realistic.

"Quintessa's source-side history is outside our chain. We can prove that the tarball we received was signed by Quintessa's PGP key and that the SHA-256 we recorded matches the tarball. We cannot prove that the contents of the tarball reflect the source EHR records as they existed in the originating site's clinical system. That is Quintessa's responsibility. The chain-coverage map per spec §10.19 names this exactly — Quintessa is in category 4, third-party systems out of contractual inspection reach for source-EHR integrity, with SOC 2 Type II reliance as the institutional substitute. We're conformant on §10.19 because the boundary is named."

Chen wrote that down.

"Show me the chain-coverage map."

Devansh pulled up the CC8.1 control description. There it was — a five-category map per spec §10.19 enumerating chain-instrumented institutional systems (the eligibility classifier, the inference-input warehouse, the CRO ingestion service), institutional systems not yet chain-instrumented (Veeva QMS, Argus, the LIMS, the trial-visit warehouse, the CTMS), third-party systems under contractual inspection (a contracted lab partner that provides hash-bound delivery), third-party systems out of contractual reach (Quintessa for source-EHR integrity, Medidata for the CRO-operated EDC), and external evidentiary artifacts hash-anchored under §10.19 (the SOC 2 report PDF, FDA correspondence). Each row carried a `coverage_map_version` and an `effective_utc`. The chain anchor was the `chain.coverage_map_published` operational event per §10.2, emitted on every map change plus monthly re-emission per §10.19's anchor-cadence guidance.

"Their SOC 2 covers the extraction process?"

"Yes. Their SOC 2 Type II report covers the extraction-and-de-identification pipeline. We have the report. It is in Veeva. The SOC 2 covers the period through December 31, 2025. The renewal report is expected in July. The PDF hash is anchored on chain via `audit.external_artifact.*` per §10.19 the day we received it."

"Have you reviewed the SOC 2 yourself?"

"Quality and Reg Affairs reviewed it. I've read the executive summary."

> **⚠️ Boundary-by-design #1 — CRO source-side history is OUT-of-contractual-reach in the §10.19 map**
> The chain captures what Helmstad ingested. It does not capture what was in the source EHR before Quintessa extracted it, or what changes Quintessa applied during extraction and de-identification. The provenance line stops at the SFTP boundary. **Per spec §10.19, this is not a hidden boundary — it is named in Helmstad's chain-coverage map under category 4 (third-party out-of-reach) with SOC 2 Type II as the institutional substitute, the SOC 2 PDF hash-anchored via `audit.external_artifact.*`, and a `chain.coverage_map_published` event per §10.2 binding the version of the map in force on each seal day.** The FDA inspector will ask about this; Helmstad's answer is the chain-coverage map plus the §10.19 worked-example shape. This was not a finding when we looked; it is a §10.19 evidentiary artifact handled to spec.

Chen circled "FDA will ask" twice in his notebook.

"What about the EDC?" he asked. "Medidata Rave. The trial-visit data."

"That is a one-way data feed from the CRO to us. Quintessa operates the EDC for this trial. We get extracts. The EDC's audit trail is the CRO's; we get extracts but the audit trail itself is not extractable in a way that we can hash and bind to the chain."

"And the runbook for the EDC extract — does it cite spec §10.16 SaaS-edge connector lag bounds?"

Devansh paused. "The EDC extract is a daily file delivery, not a streaming SaaS-edge mirror. So §10.16 doesn't strictly apply — §10.16 is for change-stream subscription connectors. But our daily file pipeline doesn't have the four numbers either — median lag, 95th-percentile SLO, alerting threshold, RTO. We use 'next-day delivery' as the wording."

Chen's pen hovered over the page. He looked at the runbook.

"'Next-day delivery' without the four numbers. If this were a §10.16 streaming mirror, that would be Finding-001 non-conformance per the §10.16 severity-classification clause — imprecise lag wording is non-conformance, not a Nit. Since it's a batch file pipeline rather than a streaming mirror, the strict §10.16 wording test doesn't apply by-the-letter, but the spirit does — and the §4.4.6 `audit.connector_source.*` family for SaaS-edge connectors is the natural shape for an EDC extract that ever moved to a streaming model. We will note that as a forward-readiness item in the report. Different from Northbridge's Salesforce-mirror finding, but the same architectural lesson — wording IS the testable claim."

"So the EDC is a vendor-managed Part 11 audit trail that lives in the CRO's environment."

"Yes."

> **⚠️ Boundary-by-design #2 — EDC audit trail is in the CRO's SOC 2 scope, not Helmstad's; named in the §10.19 chain-coverage map**
> Medidata Rave is operated by Quintessa for the NSCLC trial. The Part 11 audit trail of trial-visit data lives in Quintessa's environment. Helmstad receives extracts but does not operate the EDC and does not hold the audit trail. Vendor-management dependency. **Per spec §10.19, Quintessa-operated EDC is named in Helmstad's chain-coverage map under category 4 (third-party out-of-reach for the audit-trail integrity claim), with SOC 2 Type II as the institutional substitute.** The FDA will ask Helmstad to demonstrate that the CRO's controls meet Part 11; Helmstad answers from the §10.19 map plus the SOC 2 plus the contract. **Forward-readiness note (informative):** when the EDC extract moves from a daily file pipeline to a streaming change-stream mirror, the institution will fall under §10.16 (SaaS-edge capture connectors) and MUST quantify the four numbers; today's "next-day delivery" wording would become non-conformant under §10.16's severity-classification clause if the connector type changes.

Chen closed his notebook on the pipeline section. *The chain is honest about where it stops. That is the right answer. The question is whether the FDA inspector will accept it.*

---

### 🔐 11:00 AM — Diana on IAM, Both Sides of the Line

Diana started the IAM session by sitting down with Helmstad's identity-platform lead — a person named Rohan Patel — and asking the question she always asked first.

"Show me a credential rotation for the AI eligibility service."

Rohan pulled up the chain entry for the most recent rotation — May 1, 2026, at 03:00 UTC. The eligibility-classifier service account `svc-nsclc-warehouse-reader` had its credentials rotated automatically. The rotation event emitted a sealed chain entry under `chain_kind = "operational"` per spec §3 enumeration with the old credential fingerprint, the new credential fingerprint, the rotation reason (scheduled), the approver (Dr. Østergaard, with a two-of-three approval), the approval ticket, and the timestamp. The rotation event was a `master_key.rotated` operational event per spec §10.2; its sibling `master_key.rotation_observed` event landed in the ledger when the first chain entry under the new `key_version` arrived. The day-after seal record's `key_versions` field listed both versions per spec §10.10 — the within-day rotation handled mechanically through per-entry `key_version` lookup at §7 step 7, no special verifier path required. Constant-time comparison of `key_fingerprint` per spec §10.8 happened at §7 step 8 BEFORE any MAC compute, so a botched rotation reusing the same `key_version` for a different IKM would have surfaced at the fingerprint check, not buried in a MAC-mismatch storm.

> **✓ Confirmation #6 — AI service-account IAM is chain-coupled end to end with §10.10 rotation discipline**
> Every credential lifecycle event — issuance, rotation, revocation, scope change — for the eligibility-classifier service produces a sealed chain entry. Policy-level changes (who can rotate, who can approve, what the rotation interval is) require two-of-three approval per spec §10.5 separation-of-duties and are themselves chain entries. There is no path to change a service credential without producing a chain record. Rotation crossing the seal boundary is handled per spec §10.10 with the day-after seal's `key_versions` field listing both generations; the `master_key.rotated` and `master_key.rotation_observed` operational events per §10.2 are the SOC team's evidence trail. The §7 step 8 fingerprint check (constant-time per §10.8) catches botched rotations before MAC compute. The §10.7 software-key-adapter exclusion is in force — production builds do NOT carry the dev adapter on disk; CC8.1 names the packaging-exclusion mechanism per §10.7's "unreachable in production" requirement.

Diana had one more question on the AI side before moving to the legacy systems. "Multi-region posture? CloudHSM is in `us-east-1`. What's your DR posture if `us-east-1` becomes unavailable?"

Rohan pulled up the architecture diagram. "Single region for now. The classifier runs in `us-east-1`, the ledger is in `us-east-1`, CloudHSM is in `us-east-1`. We're not running spec §10.15 multi-region resilience yet. CC8.1 names the single-region posture explicitly. If we expand to multi-region — likely when we open the EU sites for the Phase III program next year — we'd run Pattern A per §10.15 with `us-east-1` as the seal region and `eu-west-1` as the replication region. Per-event MAC is region-agnostic because the IKM derives the same `session_key` in any region. Run-locality per §10.15 invariant 2 is enforced at the SDK process boundary — one SDK process per region, each pinned to its region's IKM custody endpoint and ledger endpoint per spec §4.4 SDK per-process region binding. The OPTIONAL `ffiec.chain.region` attribute records the region under MAC binding per §4.4 attribute table for incident-response reconstruction. Per-region event-count reconciliation per §10.15 invariant 5 with `master.cross_region_replication_completed` operational events per §10.2."

"And the operational discipline — are the staff trained on it?"

"Not yet. Training planned for Q3 when the EU expansion ramps. Today everything is in `us-east-1`. The §10.15 posture is on the architecture roadmap, not on the production posture."

Diana wrote that down. *Single region today; §10.15 Pattern A planned. Document the gap explicitly so the FDA inspector reads single-region as deliberate, not as an oversight.*

Diana asked: "Software-key adapter exclusion. Spec §10.7. How do you confirm production builds do not carry the dev adapter?"

Rohan pulled up the build pipeline. "Compile-time exclusion per §10.7. Production release builds use a `--no-dev-adapter` flag that compiles out the dev-adapter source code entirely. The packaging exclusion is the second layer — production images don't carry the dev-adapter package on disk; the registry-side check asserts the artifact is absent. CC8.1 names the packaging-exclusion mechanism per §10.7's "unreachable in production" requirement. Run-time-only environment-variable gating is NOT what we do — §10.7 explicitly names that as non-conformant. We use compile-time + packaging double-protection. The verifier under `--strict` mode also rejects any seal whose `dev_mode = true` or whose `kms_handle_uri` begins with `'plaintext-'` per §10.7 — that's the regulator-visible-line guarantee. We've never seen `dev_mode = true` in production, by construction."

> **✓ Confirmation #6a — Software-key adapter excluded per §10.7 with double-protection**
> Production builds compile out the dev adapter (compile-time exclusion). Production images do not ship the dev-adapter package (packaging exclusion). Run-time environment-variable gating is explicitly NOT used — that pattern is non-conformant per §10.7. The verifier under `--strict` mode rejects any seal whose `dev_mode = true` or whose `kms_handle_uri` begins with `"plaintext-"`. Per the §10.7 regulator-visible-line guarantee, development-only key material cannot accidentally ship to production via configuration drift, AND a chain that the dev adapter produced is structurally rejected at verification time even if it somehow reached production storage.

Diana asked: "And the clinical research coordinators? The CRCs at the trial sites who accept or reject eligibility decisions?"

"They authenticate via SSO. The SSO subject is part of every chain entry where they make a decision. The SSO upstream is Helmstad's federated identity provider — Okta. Okta is federated to each site's institutional identity provider via SAML."

"Walk me through it."

He did. A CRC at Mass General authenticates with their Partners HealthCare credentials, the Partners SAML IdP asserts to Okta, Okta asserts to the eligibility-classifier UI, the UI captures the SSO subject and writes it into the chain entry when the CRC accepts or rejects a candidate.

"How many institutional IdPs do you federate with?"

"Twelve sites in the trial. Twelve institutional IdPs, give or take — some of the smaller sites use a shared community IdP."

"And each institutional IdP is the upstream root of trust for the SSO subject."

"Yes."

Diana wrote that down carefully.

> **✓ Confirmation #7 — CRC reviewer decisions are sealed under SSO subject attribution**
> Every clinical research coordinator's accept/reject decision is sealed with the CRC's SSO subject as the attribution. The chain entry captures who decided. The SSO upstream resolves through Okta to twelve institutional IdPs. The chain trusts the SAML subject as far as the institutional IdP is trustworthy.

"What about the Argus DBAs?"

Rohan's expression shifted.

"Argus is on a SQL Server database. The DBA team has admin access to the database. There is a Helmstad Active Directory group that grants Argus DBA membership. Twelve people are in the group right now. Six are over thirty days old in their grant. Two are over a year. The Argus application has its own Part 11 audit trail, which writes to a database table. The DBAs have UPDATE permission on the audit-trail table. So technically yes, they could modify it. We have not seen evidence of that."

Diana asked, very calmly, "Has anyone reviewed whether the audit-trail table was modified in the past twelve months?"

Rohan paused. "I do not know. I would have to ask the DBA team."

"Please ask them this afternoon."

Diana paused. "And the application-level enforcement? Spec §10.3 names two layers — application level (no UPDATE / DELETE statements on event tables) AND database role level (writer role grants INSERT and SELECT only; UPDATE / DELETE / TRUNCATE revoked). For Argus to be brought under §10.3, both layers would need remediation."

Rohan: "Application level is upstream — that's an Oracle Health Sciences product, not Helmstad code. We don't have access to modify Argus's source. Database role level is on us, and it's where the DBAs get UPDATE today. The remediation path under §10.3 would be either (a) work with Oracle to bring Argus's audit-trail discipline up to chain-grade — unlikely on our timeline; or (b) move the Argus audit trail to a separate chain-instrumented store, with Argus writing through a connector under §4.4.6 SaaS-edge connector source attribution. That's a multi-quarter project. Today, Argus is in §10.19 chain-coverage map category 2."

"Understood. The CAPA for Argus has to acknowledge both layers. Don't write it as a database-role-only change — that closes half the gap and leaves the application-level vulnerability."

Rohan wrote that down.

> **⚠️ Gap-001 — Argus Part 11 audit-trail table is mutable by twelve DBAs; would be non-conformance with §10.3 if Argus were chain-instrumented**
> Twelve DBAs hold admin access to the Argus database. The Part 11 audit trail is a SQL Server table with no out-of-band integrity binding. The DBAs have UPDATE permission. **Spec §10.3 requires append-only enforcement at two layers — application level (no UPDATE / DELETE statements on event tables) and database role level (the writer role grants INSERT and SELECT only; UPDATE / DELETE / TRUNCATE revoked).** Argus is NOT chain-instrumented today; the §10.3 requirement does not bind the FDA inspection by spec. But the architectural shape — a mutable audit-trail table whose integrity rests on twelve people's discipline — is exactly the shape §10.3 was written to remove. The FDA has issued 483s on this exact configuration in the past year. Per the §10.19 chain-coverage map, Argus is named under category 2 (institutional system not yet chain-instrumented); the §10.19 evidentiary substitute today is the quarterly access review, which Raj will surface as Gap-005. The remediation path is to bring Argus under the chain — application-level append-only plus database role-level INSERT-and-SELECT-only per §10.3. The CAPA is what we will recommend.

"And the rest of the legacy AD? The non-AI privileged users — the application admins, the integration accounts, the service accounts that talk to the safety database, the data-engineering pipelines?"

Rohan pulled up the AD groups list. Forty-three privileged groups. Two hundred and seventy-eight unique users. Quarterly access reviews. The most recent one — March — had been signed off by the application owners with the Q1 attestation language Helmstad used for HITRUST.

"The chain reaches the AI service-account boundary. Every other AD group is reviewed quarterly per the institution's HITRUST-equivalent attestation, but the access reviews are shape-checks — they confirm group membership, not the privileges the group grants on individual tables or files. Per spec §10.1 key-fingerprint reconciliation, our chain-instrumented systems run weekly fingerprint reconciliation against the IKM roster — that's a one-week compromise-detection window. The legacy non-chain systems run quarterly access reviews — that's a 90-day window, and the review only checks group membership, not table-level privileges. The two regimes are not commensurable."

Diana wrote that down. "And per §10.1 multi-deployment uniqueness enforcement — single global IKM registry. You're single-region today, single-deployment, so the multi-deployment posture isn't tested yet. When you go multi-region for the EU expansion, you'll need a single global IKM registry covering every deployment. Acceptable patterns are AWS Secrets Manager with cross-region replication, Azure Key Vault with multi-region failover, or a centralized registry service. Document the choice in CC8.1 before EU launch."

"Yes. CC8.1 will name the chosen pattern when we go multi-region."

> **✓ Confirmation #7a — Single-region IKM registry today; §10.1 multi-deployment posture documented for EU expansion**
> Helmstad operates a single-region single-deployment posture today. Per spec §10.1 weekly fingerprint reconciliation, the chain-instrumented systems run weekly fingerprint reconciliation against the IKM roster, bounding the master-compromise detection window to at most one week. The legacy non-chain systems run quarterly access reviews per the HITRUST-equivalent attestation — a 90-day window, group-membership-shape only. When Helmstad expands to multi-region for the EU Phase III sites, the institution's CC8.1 will name a single global IKM registry per §10.1 multi-deployment uniqueness enforcement. Today's posture is conformant; the EU expansion will require an explicit CC8.1 update.

"And the lab DBAs? The LIMS?"

"Same architecture. Lab DBAs have admin access. The LIMS audit trail is on the same database. No chain. Different team, different group, similar exposure."

> **⚠️ Gap-002 — LIMS audit trail has the same DBA-mutable shape as Argus**
> The lab information management system audit trail is a database table mutable by the LIMS DBA team. No out-of-band integrity binding. No chain. The lab data feeds into the trial through the EDC and through case report forms. The integrity of the lab audit trail is operator discipline. Same §10.3 architectural lesson, same §10.19 chain-coverage-map category 2 disposition. Same CAPA shape.

Diana made a note in her book that the AI-side IAM finding was clean within its boundary, the legacy-side IAM finding was as bad as Helmstad's peers, and the relationship between the two was the actually-interesting finding — exactly the same shape she had written down at Mercator three weeks ago.

She had one more question before lunch.

"Per spec §10.23, when the FDA needs subject-keyed retrieval — `produce all eligibility decisions for subjects enrolled at site 04 during March` — how do you respond?"

Rohan nodded. "Shape 1 chain-anchored. Each CUEC entry is itself a chain entry under `chain_kind = 'operational'`. The subject identifier is hashed per institution policy — lowercased subject ID hashed under SHA-256, that's the canonicalization documented in CC8.1. The chain entry carries `consumer_index.consumer_id_hash`, `consumer_index.run_id`, `consumer_index.seq`, `consumer_index.relationship` — for us, `relationship` is `'screening_subject'`, an institution-named value per §10.23 enumeration. The chain alone is the integrity-bound retrieval substrate. Append-only per §10.3. The FDA's CID-equivalent — the inspector's data request — produces a verifiable subject-keyed list from the chain itself, no separate index attestation step the inspector has to trust."

"Why Shape 1 and not Shape 2?"

"Volume. Phase II screening is hundreds of decisions per site per month, not millions. Shape 1 is cost-effective at our volume. Shape 2 is for high-volume institutions where per-consumer chain entries don't scale economically. The institution's CC8.1 names the chosen shape and the rationale — that's all spec §10.23 normates."

Diana wrote that down. *§10.23 Shape 1 — chain-anchored CUEC. Same posture pattern Atrio used last week with per-tenant Shape 1; same shape Mercator wanted but didn't have funded.*

> **✓ Confirmation #8a — Subject-keyed retrieval is chain-anchored per §10.23 Shape 1**
> Per spec §10.23, the consumer-correlation index for FDA subject-keyed CID-equivalent retrieval operates as Shape 1 (chain-anchored). Each subject-binding entry is its own chain entry under `chain_kind = "operational"` per §3 enumeration, append-only per §10.3 by chain construction. The chain alone is the integrity-bound retrieval substrate; no separate institution-controlled index that the FDA must trust. The institution's CC8.1 names the shape and the canonicalization mechanism for `consumer_index.consumer_id_hash` per §10.23 normative requirement.

---

### 🧪 12:00 PM — Lunch and the Argument About ALCOA+ "Original"

The team gathered in a side conference room with sandwiches from the building's cafeteria. Tom was on the phone with the CQO. Dr. Østergaard had ducked out to a regulatory affairs huddle.

Dawn put her sandwich down before she'd taken a bite.

"Let's talk about the morning."

Mike: "AI side is real. The verifier resolves a 90-day-old entry in four seconds, executing all twelve §7 steps. Model card and protocol document are hash-bound through `audit.model_handover.*` per §10.21. Reviewer decisions are sealed under SSO with `audit.redaction.*` discipline per §10.22. This is the Part 11 audit trail the FDA's draft AI/ML guidance asks for."

Chen: "The pipeline-to-chain reconciliation works at the SFTP boundary. We can prove that the tarball Quintessa sent us is the tarball we processed — `audit.external_artifact.*` per §10.19. We cannot prove what was in Quintessa's source systems before extraction. That is the right boundary to draw and §10.19 names it explicitly in our chain-coverage map. The FDA inspector is going to push on it; the §10.19 map is the answer."

Diana: "AI service IAM is sealed under §10.10 rotation discipline. CRC IAM is SSO-coupled — the chain trusts the institutional SAML subject as far as the SAML subject is trustworthy. Argus DBAs and LIMS DBAs are not chained — Gap-001 and Gap-002 against the §10.3 append-only architectural shape. Same legacy posture as the diary baseline."

Raj: "I have the full Argus database review after lunch. Going in expecting what Diana described — DBA-mutable audit table, no integrity binding, no reconciliation. Same §10.3 architectural lesson."

Dawn: "Okay. Question. ALCOA+ has nine attributes. Walk me through the AI side. Attributable — yes, SSO subject bound under per-event MAC per §4.1. Legible — yes, captured JSON per §5.2 best-evidence posture. Contemporaneous — yes, the chain entry is within 200 ms of the decision; `captured_at` carries nanosecond precision per §4.4. Original — that's where I want to argue. Original on the AI side means the chain entry is the original record of what the model said, with §5.2 naming both the captured JSON (content-bearing form) and the canonical bytes (integrity-bearing form) as originals under FRE 1001(d). But the FDA cares about Original on the inputs. The candidate's lab values, the tumor staging, the prior-therapy history. Where does Original live on the inputs side?"

The room got quiet for a beat.

Mike: "The chain is the original record of what the AI saw. It is not the original record of what the EHR said before Quintessa touched it."

Dawn: "Exactly. Two different definitions of Original. Original-as-input-to-the-model and Original-as-source-of-truth. The chain handles the first cleanly. The CRO handles the second through their SOC 2 and their PGP signature on the SFTP delivery. The FDA inspector is going to ask which Original we are claiming."

Tom, off the phone: "The CQO is asking me what I think the inspector will weight. I told him my read — the inspector will accept the model-side Original because the chain is convincing. They will push on the source-side Original, especially if there is a discrepancy in the trace."

Chen: "There is going to be a discrepancy. There always is. The question is whether we find it before the inspector does."

Dawn tapped her pen on the table. "Tom — tell the CQO that the report will frame ALCOA+ as a per-attribute walk-through with the boundary called out for each attribute. Original gets two paragraphs. One for AI-output Original, which the chain handles. One for source-data Original, which Quintessa's SOC 2 handles. We are honest about the boundary."

Tom relayed it. He came off the call after a minute. "He is good with that framing."

Diana: "Three more attributes I want to flag. Complete — the AI side is complete on what the chain captures, but the chain does not capture pre-screen exclusions that happen in the CRO's pipeline before the candidate ever reaches our classifier. So Complete on the AI side is complete-for-classified-candidates, not complete-for-all-screened-candidates — and again, §10.19 names that. Consistent — yes for the chain; per-event canonical bytes per RFC 8785 JCS plus chain-stamp preservation per §6 storage rules. Enduring — yes, the daily seal cadence per §4.2.1 and the daily Ed25519 signature on CloudHSM under FIPS 140-2 Level 3 per §10.5 mean the seal endures across rotations and key generations under §10.10. Trusted-time integration per §10.14 is RECOMMENDED but not REQUIRED for v1.0; we run NTP discipline per §10.4 and document it. Available — yes, TesseraSeal retrieval works in four seconds for a 90-day-old entry. We confirmed this morning."

Dawn: "Good. Write up the ALCOA+ walk for both sides. Boundary called out per attribute."

Luis had been quiet. He looked up from his laptop.

"I was reading their CRO ingestion runbook while you were all talking. The runbook section title is `Multi-Tenant Operations` — it doesn't cite spec §10.1, which is a §10.18 cross-referencing nit; I flagged it. The SFTP delivery from Quintessa lands in an S3 bucket. The bucket has versioning enabled and a 7-year retention lock. CloudTrail is enabled. The CloudTrail logs are written to a separate AWS account and the cross-account permissions are configured with an IAM role that the Helmstad SecOps team controls. The cross-account write path looks clean. The S3 bucket itself — the one that holds the Quintessa tarballs — is owned by the Helmstad data-engineering team. They have `s3:PutBucketLogging` permission. Three engineers."

Dawn put her sandwich down again.

"Three engineers can disable CloudTrail logging on the ingestion bucket."

"Three."

"And the chain entry references a SHA-256 of the tarball at ingestion time. So even if someone disabled CloudTrail and replaced the tarball, the chain entry would still verify against the recorded SHA-256 — the §10.19 `audit.external_artifact.sha256` is bound under the per-event MAC per §4.1, which is bound under the daily Merkle seal per §4.2, which is bound under the HSM-rooted Ed25519 signature per §4.3. Three independent layers per §1.4 compositional security."

"Yes. The chain catches the tampering at the recorded-hash layer. But the FDA inspector will ask why three engineers can disable CloudTrail at all. That is a separate finding — not a chain integrity finding, but an institutional-control completeness finding the §10.19 chain-coverage map should describe explicitly."

Dawn wrote it down. "Good catch. Nit-001 — coverage-map should name the CloudTrail-on-ingestion-bucket dependency under §10.19 category 1 (chain-instrumented institutional systems) with the three-engineer-bypass risk explicit."

Tom, off the phone again: "The CQO asked one more question. He wants to know — between us — whether this assessment is going to recommend that Helmstad delay the FDA inspection."

Dawn looked at him.

"That is his question, not mine," Tom said.

"That is a question for Dr. Østergaard," Dawn said. "Not for us. Our job is the assessment."

The team finished lunch. Elena had already wandered off at 12:25 with the CTMS admin's calendar invite on her laptop. The rest of them rinsed coffee cups and walked back out to the engineering floor.

---

### 🔄 1:00 PM — Mike on the API Layer

Mike's afternoon was the API surface. Helmstad runs the eligibility classifier behind an internal API gateway with a custom Lambda authorizer and a request-signing layer. Every call to the classifier — from the CRC review UI, from the batch-screening daily job, from the CRO data-feed ingestion side — goes through this gateway.

Every call emits a chain entry. Request, headers (redacted pre-MAC at the SDK boundary per §10.22 — there is no PHI in the request anyway, since the inputs are all de-identified, but the `audit.redaction.*` family is emitted for completeness), authorizer decision, downstream service, response code, response hash, latency. Sealed. The OTLP transport identification per spec §4.4.3 — `service.name = "nsclc-phase2-classifier"`, `service.version = "3.7"`, `ffiec.chain.spec = "v1.0"`, `ffiec.chain.posture = "ffiec"`, `ffiec.chain.format_version = "v1"` — is set on every Resource. The collector pass-through per §4.4 is verified — no transformation rewrites the chain attributes between SDK and ledger; severity-based filtering per §4.4.4 is exempted for chain traffic; receiver stamping is `SeverityNumber` in the 9–20 range per §4.4.4.

"Show me a call that returned a 5xx," Mike said.

The API engineer pulled up an entry from April 8. 503. The authorizer had failed because CloudHSM was briefly unreachable during a planned maintenance window. The chain entry recorded the failure. The downstream classifier had served a fallback `model-unavailable` response. The fallback was itself sealed as a separate inference-shaped entry with `model_id: fallback-noop`, `classification: unavailable`, and a reference to the upstream API failure. The chain captured the routing decision through the §4.4.1 routing schema — `audit.routing.attempt`, `audit.routing.failover` with `failover_reason = "transport_error"` (CloudHSM unreachable), and a terminating `audit.routing.refused` event with `refusal_reason = "all_circuits_open"` because no fallback model had been configured for HSM-down conditions and the institution's policy version `nsclc-router-v2.1.0` had no alternate route. The §10.10 rotation-observed event for the day showed no `key_versions = [old, new]` mixed seal — the rotation that morning had completed before the HSM maintenance window. HSM unavailability is named in spec §4.3.1 — captured events continue to be ingested and chained because per-event MAC is independent of the HSM; the seal job retries with exponential backoff, and Helmstad's CC8.1 covers the 72-hour notification posture per §4.3.1.

> **✓ Confirmation #8 — API gateway and fallback paths are sealed**
> Every classifier API call — including failed authorizations, planned-maintenance outages, and fallback responses — produces a sealed chain entry. The April 8 CloudHSM maintenance window shows up cleanly: 14 classifier API calls during the window, all routed to the fallback responder, all sealed. No silent gaps.

Mike turned to the lead API engineer. "What about the CRO ingestion API? The endpoint Devansh's team uses to trigger the SFTP-to-warehouse pipeline?"

The API engineer pulled it up. The ingestion endpoint emits a chain entry of type `cro-ingestion` with the SFTP source, the PGP-signature verification result, the tarball SHA-256, and the file/byte counts. That part is chained.

"What about the EDC extract endpoint? The one that pulls trial-visit data from Quintessa's Medidata Rave?"

The API engineer paused.

"That is a one-way file transfer. SFTP delivery from Quintessa, lands in a separate S3 bucket, gets unpacked into the trial-visit warehouse. The unpacking is not chained. The trial-visit warehouse is not chained. We do not call the classifier on trial-visit data — that is post-enrollment, not screening."

"So the EDC extract has no chain."

"Correct."

> **⚠️ Surprise #5 — EDC extract pipeline has no chain**
> The Medidata Rave extract from Quintessa is delivered via SFTP, unpacked into the trial-visit warehouse, and processed without chain entries. No HMAC. No Merkle binding. No verifier. The trial-visit data is in scope for the FDA's BIMO inspection because it is part of the GCP record. The EDC's own audit trail lives in Quintessa's environment. Helmstad does not bind the extract to its own integrity record.

Mike wrote that down.

"What about Veeva QMS? The QMS is the source of truth for the protocol document — and the protocol-document hash is what `criteria_doc_hash` resolves to on every chain entry. How does Veeva's audit trail compose with the chain?"

The API engineer pulled up the Veeva QMS audit-trail browser. It was deep — every change to every controlled document was recorded with a timestamp, an actor, a before/after comparison, and a reason code. But the standard audit-trail report defaulted to filtering on `EffectiveDateChangeReason` not in `('Late Effective Date',)`.

"What is `Late Effective Date`?"

"A Veeva workflow that allows up to 30 days of effective-date back-dating. If a document was authored on March 1 but went through review and approval until March 28, the user can set effective date to March 1 retroactively. That's the workflow Veeva ships."

Mike's pen paused. "Does that show up in the audit trail?"

"Yes. With a reason-code annotation. But the standard audit-trail report doesn't include rows where the reason is `Late Effective Date` unless the user explicitly toggles a filter."

"So a document that was effective-back-dated by 28 days does not appear in the standard audit-trail report unless you know the filter exists."

"Correct."

> **⚠️ Gap-003 — Veeva QMS Late Effective Date workflow is invisible in the standard audit-trail report unless the filter is explicitly toggled**
> Veeva's `Late Effective Date` workflow allows up to 30 days of effective-date back-dating. The audit trail records the change, but the standard report filters those rows out by default. A user (or auditor, or inspector) running the standard report sees a clean trail. Toggling the filter shows the back-dated changes. **Per spec §10.18 runbook cross-referencing, the institution's CC8.1 control description for Veeva audit-trail access SHOULD cross-reference §10.3 append-only enforcement and name the Late-Effective-Date filter discoverability requirement explicitly — without it, an examiner reading the runbook does not know the standard report is incomplete.** This is exactly the discoverability-deficiency shape §10.18 was written to remediate. The Veeva audit trail itself is real; the standard report is incomplete by default. Helmstad's CC8.1 today does not cross-reference §10.18 or name the Late-Effective-Date filter as a load-bearing discovery item.

"What about the safety database? Argus. Does anything talk to Argus over an API?"

"The serious-adverse-event reporting integration. When a SAE is reported at a site, it gets entered into the EDC, the EDC pushes it to Quintessa's pharmacovigilance team, and they enter it into Argus. The Helmstad-side pull from Argus is a daily ETL into our reporting warehouse. That ETL is not chained."

> **⚠️ Surprise #6 — Argus-to-warehouse ETL has no chain**
> Adverse-event data flows from Argus to Helmstad's reporting warehouse via a daily ETL. The ETL is not chained. No HMAC. No reconciliation. The warehouse is what the QA team uses to prepare safety summaries for the DSMB. Same diary-baseline pattern: critical regulatory data on an unchained pipeline.

Mike turned to the API engineer. "How does the SAE reporting flow get audited today?"

"Standard validation. SOPs. Periodic reconciliation between Argus and the warehouse. The reconciliation is run by a person, with a checklist."

"How long does the reconciliation take?"

"A day. Quarterly. They sample."

"Sample what?"

"They check that the count of SAEs in Argus matches the count in the warehouse for the quarter. If counts match, they sign off. If counts don't match, they investigate the difference."

"What if the counts match but the contents are different?"

The API engineer paused.

"They don't check that."

> **⚠️ Surprise #7 — Quarterly Argus-warehouse reconciliation checks counts, not contents**
> The reconciliation between Argus and the safety-reporting warehouse is a quarterly count check. If the count of SAEs matches, the reconciliation passes. There is no per-row hash comparison. A row whose contents have been modified after Argus-side capture would not be detected by the reconciliation. The chain at the AI side has no analogue here. The reconciliation is the only integrity gate, and it is a count gate.

Mike closed his notebook on the API section. *Two API surfaces. One chained. One not. The line is exactly where we expected it to be.*

He flipped back through to the morning's notes. The chain on the AI side was conformant on every property he had tested: per-event MAC integrity per §4.1, daily Merkle seal under RFC 6962 per §4.2, HSM-rooted Ed25519 signature per §4.3 with `sign_payload_version = "v1.0b"` dispatch, OTLP transport identification per §4.4.3 with all five required Resource attributes, severity-resilient receiver stamping per §4.4.4, model-handover attribution per §10.21, training-data retention floor per §10.20, redaction discipline pre-MAC at the SDK boundary per §10.22, software-key adapter excluded per §10.7 with double-protection, IKM generation inside CloudHSM per §10.6.1's highest-assurance pattern, partition ceremony attestation per §10.17 with `entity_affiliation`, run-resume contract per §10.25 with three-place tail acquisition, fork detection per §10.25 ledger ingestion cross-check, genesis-block uniqueness per §4.4 and §10.25, key-fingerprint reconciliation per §10.1 weekly cadence, append-only enforcement per §10.3 at both application and database role layers, constant-time comparison per §10.8 at §7 step 8 and step 9, evidentiary artifacts retention per §10.13, NTP-based trusted-time foundation per §10.4 with §10.14 RFC 3161 noted as an institution-side option for the future, exit-code contract per §10.12, cadence-and-dev-mode check per §7 step 12, gen_ai-completeness check per §7 step 12a, and the wire-bound observation rule per §4 (observation and verification are valid only on wire-or-on-disk artifacts; in-process state is NOT a spec-conformant view).

Every one of those properties had been tested today, on the AI side. Every one had passed. The §1.4 compositional security argument held: three independent layers, each defending against a different attack class, each FIPS-standardized, each bound to the institution's CC8.1 control description by spec section number per §10.18 cross-referencing.

The legacy side was a different story.

---

### 🧬 2:00 PM — Chen on the Pipeline, Raj on the Database

Chen and Raj split at 2 PM. Chen took the CRO ingestion pipeline and the EDC extract pipeline. Raj took Argus and the LIMS database.

Chen started with the daily Quintessa SFTP delivery. He walked the pipeline end to end — SFTP receiver, PGP-signature verifier, tarball-hash recorder, unpacking, schema validation, warehouse load, foreign-key binding to the chain entry. The chain entry was always written before the warehouse load committed. If the chain write failed, the warehouse load was rolled back.

"That is a transaction boundary," Chen said. "Chain-or-roll-back."

"Yes. We made that choice deliberately. The platform team decided that the warehouse should never hold a row that was not chain-bound."

Chen wrote that down. *Transaction boundary at the chain. That is a design choice that survives audit.*

He moved to the EDC extract pipeline. Different shape. The Medidata Rave extract lands in S3 as a daily file. The unpacking job is an Airflow DAG that reads the file, transforms it, and writes to the trial-visit warehouse. No chain. No transaction boundary.

"Why is the EDC pipeline different?"

"It was built before TesseraSeal was deployed. The eligibility classifier was the pilot. The EDC extract was on the legacy stack. We did not extend the chain to the trial-visit warehouse because the eligibility tool does not consume trial-visit data — that is post-enrollment."

"But the trial-visit data is what the FDA cares about for GCP compliance."

"Yes."

"So the chain is on the part that was the pilot. The part that the FDA cares about most for GCP is on the legacy stack."

Devansh — the platform engineer — was honest. "Yes. That is the situation. Dr. Østergaard is aware. The phase-2 plan is to extend the chain to the trial-visit warehouse before the next BIMO inspection. We did not get there in time for this one."

Chen wrote that down. *Honest. Same shape as Mercator's lab pipeline.*

Meanwhile, Raj had walked across the engineering floor to the database team's bullpen. The Argus DBA lead — a person named Karthik Sharma — pulled up the SQL Server management console.

"Show me the audit-trail table for the past twelve months."

Karthik did. The table had 387,000 rows. Each row had a timestamp, an actor, an action, an entity-id, and a before/after JSON. The actor field was a username from the Argus authentication layer. The timestamps were monotonic. Raj noted that none of the rows carried any out-of-band integrity binding — no HMAC, no Merkle anchor, no signature. The integrity claim rested entirely on the trust that no DBA had modified the table.

If Argus had been chain-instrumented under the spec, every row would carry `payload_hash` (32 raw bytes per spec §4.1), `prev_hash` (32 raw bytes linking to the previous entry), `key_fingerprint` (16 raw bytes from `SHA-256(utf8(tenant_id) || ikm)[:16]` per spec §3), `key_version` (the IKM generation), `format_version`, and `mac_computed_at_utc`. The MAC would cover the canonical bytes of the row content per RFC 8785 JCS, excluding the chain-stamp fields per spec §4.1 inviolate property 7. The day's rows would aggregate into a Merkle root per spec §4.2 with `(run_id, seq)` ordering, signed by an HSM under FIPS 140-2 Level 3 per spec §10.5. A DBA UPDATE on any row would be detected at chain verification time — the recomputed `payload_hash` would not match the persisted value at §7 step 9 with reason `payload_hash MAC mismatch at seq N`.

"Has anyone modified rows in this table in the past twelve months?"

Karthik paused. "The DBA team has UPDATE permission. I do not believe anyone has modified rows. I would have to query the SQL Server transaction log to be sure, and the transaction log retention is 30 days."

"So you can confirm no modifications in the past 30 days but not before that."

"Correct."

"And the SQL Server transaction log retention — who controls it?"

"The DBA team. Same group of twelve."

Raj tapped his pen against the page. "So the integrity-detection window for Argus is 30 days, controlled by the same group that has UPDATE permission on the audit-trail table. If a DBA modified a row 31 days ago, you cannot prove they did or did not. The transaction log is gone."

"Correct."

"And the transaction log itself — does it have integrity binding? An integrity-bearing transaction log would let you detect modifications even after the 30-day rotation window."

"It does not. SQL Server transaction logs are append-only by construction at the page level, but they are not cryptographically integrity-bound. A DBA with sufficient privileges can re-initialize the log file. The DBA team has those privileges."

"So the integrity claim on the Argus audit trail is — entirely — that twelve people have not done what they have permission to do."

"Correct."

Raj wrote that down a third time. He drew a circle around it. *Twelve-person discipline as integrity foundation. Compare to chain's three-layer §1.4 compositional security with FIPS-standardized primitives. Different security model. Same regulatory expectation.*

If Argus had been brought under §10.3 — application-level append-only plus database-role-level INSERT-and-SELECT-only — the DBA UPDATE permission would have been revoked at the role layer. If Argus had been brought under §10.1 key-fingerprint reconciliation, every observed `(tenant_id, key_version, key_fingerprint)` triple would have been matched weekly against the IKM roster, and any drift would have been a high-priority alert. If Argus had been chain-instrumented, the compromise-detection window would have been bounded to one week per §10.1 reconciliation cadence rather than the 30 days the SQL Server transaction log allowed.

Raj thought about the gap. The chain's compromise-detection window per §10.1 was at most one week plus the time between observation and remediation. The Argus compromise-detection window was 30 days, then forever. The two security models were not commensurable; the chain's model was a cryptographic-detection-with-bounded-window pattern, and the Argus model was a trust-the-DBA-discipline pattern. The FDA had been issuing 483s on the Argus pattern for years.



Raj wrote that down. "Who has UPDATE permission on the audit-trail table specifically?"

Karthik queried it. "Twelve DBAs and the application service account. The application writes new rows. The DBAs are technically able to modify."

"Has anyone reviewed who has UPDATE permission on this table in the past year?"

"Not that I know of. It would be in the access-review records that QA maintains."

"Can you pull those?"

Karthik pulled them up. Quarterly access reviews. The most recent was Q1 2026 — March. The review confirmed that twelve DBAs were members of the Argus-DBA group. It did not specifically call out UPDATE-on-audit-trail-table as a privilege, because the group grants admin and admin includes UPDATE on the table by default.

"So the access review confirms that twelve people are admins. It does not specifically attest that those twelve have not modified the audit-trail table."

"Correct."

Raj wrote that down twice.

> **⚠️ Gap-005 — Argus access review attests group membership, not table-level modification; would be Partial under §10.3 if Argus were chain-instrumented**
> The quarterly access review confirms that twelve DBAs hold the Argus-DBA group, which grants admin. The review does not specifically attest that the audit-trail table has not been modified. The review is shape-checking the IAM, not integrity-checking the audit trail. The FDA inspector who looks at this finding will ask the same question. **Spec §10.3 requires append-only enforcement at both the application layer and the database role layer; the database role for the chain ledger writer SHOULD be granted INSERT and SELECT only with UPDATE / DELETE / TRUNCATE revoked.** Argus is NOT chain-instrumented, so the §10.3 binding does not apply by spec — but the architectural shape would fail §10.3 if it were. Remediation path: bring Argus under the chain or, at minimum, revoke UPDATE on the audit-trail table at the database-role layer.

Raj walked over to the LIMS database next. Same architecture. Same DBA admin group. Same access-review pattern. Same 30-day SQL Server transaction-log retention. Same finding.

He closed his notebook on the database section and walked back to the conference room.

---

### 🧪 2:15 PM — Elena on the CTMS (and the KOL CRM)

Elena had a 90-minute slot for the CTMS walkthrough. The CTMS admin was a person named Jordan Beck. Friendly. Eight years at Helmstad. Had migrated the CTMS from Veeva CTMS to Salesforce three years ago because the commercial team wanted KOL relationship management on the same instance.

"Walk me through how a site monitoring visit gets recorded," Elena said.

Jordan: "The CRA — clinical research associate — visits a site. They take notes during the visit. They come back and enter a Site Monitoring Visit Report into the CTMS. The Report object is the canonical record. They attach the visit notes to the Report object as a Long Text field — Visit Summary."

"Visit Summary. Long Text. Salesforce field history?"

Jordan paused. "Field history is enabled for some fields. Let me check Visit Summary specifically."

They checked. Field history was enabled for the Report's status field, the visit-date field, and the CRA-name field. It was not enabled for Visit Summary or for any of the comment fields.

"So if a CRA edits the Visit Summary after submitting the report, it is not in field history."

"Correct."

"Why is field history not enabled for Visit Summary?"

"Storage cost. Salesforce charges per field-history row at scale, and Visit Summary is a high-volume free-text field. We turned it off for the long-text fields years ago to control the bill."

Elena had heard this answer four times in five weeks. She wrote it down anyway.

> **⚠️ Gap-007 — CTMS Site Monitoring Visit Summary field history is disabled; CTMS is not chain-instrumented**
> The Long Text field that holds CRA visit notes does not have Salesforce field history enabled. Edits after submission are not recorded. The field is the canonical record of what the CRA observed at the site. The decision to disable field history was made for storage-cost reasons. The HITRUST-style attestation pattern from peer engagements — where this is renewed as a Partial finding cycle after cycle without the remediation being funded — is the same pattern here. **Per spec §10.19, CTMS is named in the chain-coverage map under category 2 (institutional system not yet chain-instrumented), and the §10.19 evidentiary substitute is the Salesforce field-history shape-check Elena just performed.** Remediation under spec discipline would be a §4.4.6 SaaS-edge connector mirroring CTMS records into a chain-instrumented store, with `audit.connector_source.system = "salesforce-cdc"` and the four §10.16 numbers — median lag, 95th-percentile SLO, alerting threshold, RTO. The §10.16 wording test would apply at that point; today it does not because CTMS is not yet chain-instrumented and the §10.19 map describes the boundary explicitly.

"What about KOL outreach? Same Salesforce instance, separate org?"

"Same instance. Same Salesforce org-within-the-tenant. The commercial team uses a Contact-and-Activity model for KOL relationships. Field history on the Activity comment field is also off, for the same reason."

"PHI in the KOL comments? Or only commercial context?"

"Only commercial context. The commercial team is not supposed to log PHI. We have a DLP rule that scans the field for known PHI patterns and flags it for review. The DLP rule has caught fewer than ten instances in three years."

"Has the DLP rule caught anything in the past quarter?"

"Two cases. Both reviewed by Compliance. Both turned out to be false positives — names of physician contacts that the rule confused with patient names. No PHI in the field."

> **⚠️ Partial-002 — CTMS-as-KOL-CRM has DLP-only protection on PHI in free-text fields**
> The same Salesforce instance that holds Site Monitoring Visit Reports also holds KOL outreach activities for the commercial team. PHI is policy-prohibited but technically possible. The DLP rule has caught two cases in the past quarter, both false positives. The control is operator discipline plus DLP, not chain-grade. The CTMS is the same instance the diary baseline saw at peer companies. **Per spec §10.22, redaction discipline lives at the SDK boundary pre-MAC; a DLP-on-Salesforce field control is post-capture and outside the §10.22 conformant posture entirely.** This is the right architectural lesson — the DLP control catches what it catches, and the field is therefore not chain-grade evidence. The remediation under §10.19 is to bring the field under a chain-instrumented mirror with §10.22 redaction at the SDK boundary; today the field is in §10.19 category 2 (not yet chain-instrumented) with DLP as the §10.19 evidentiary substitute.

Elena closed her notebook on the CTMS. "Jordan. Thank you. I appreciate the directness."

"I have been waiting for someone to ask the right questions about Visit Summary. The platform team has been trying to fund the field-history uplift for years. The remediation plan has been the same for four cycles."

"Of course."

She walked back to the conference room. *Same answer. Different industry. Fourth time in five weeks.*

---

### 📊 3:00 PM — The Reconciliation Test

Dawn called the team together at 3 PM in the main conference room. Dr. Østergaard was back from his regulatory-affairs huddle. The CQO was on Zoom from the Boston office.

"We are going to pick five eligibility decisions at random and trace them end to end," Dawn said. "Backwards from the chain entry to the CRO source. Forwards from the chain entry to enrollment and to the trial-visit data. The AI side is the chained side. Both ends are the legacy side. We are testing the boundary."

Dr. Østergaard nodded. "Pick the five. I will not interfere."

Dawn turned to Mike. Mike pulled up the chain database in the conference room's projector and used a deterministic pseudo-random sampler — based on the date, seeded so the choice was reproducible — to select five decisions from the past 60 days. The sampler returned five entry IDs.

```
nsclc-2026-04-02-mgh-00018
nsclc-2026-04-15-mskcc-00074
nsclc-2026-04-21-uchicago-00049
nsclc-2026-04-28-msd-00031
nsclc-2026-05-03-stanford-00027
```

"Run the verifier on all five," Dawn said.

Mike ran the verifier on all five. Four seconds each. Five PASS results. Five times twelve verification steps each. Sixty steps. All clean. The §7 step 12a `gen_ai_model_identifier_missing` check fired and passed on each of the five — every entry carried both `gen_ai.request.model` and `gen_ai.response.model` non-empty per spec §4.4 and §7 step 12a normative requirement. Constant-time compare per §10.8 visible in the verifier source. Verifier exit code 0 on each per §10.12. The §10.25 run-resume contract had been exercised correctly — every entry's `prev_hash` linked to the previous entry's `payload_hash` walked structurally, no genesis-form entries at any `seq > 1` per §4.4 genesis-block uniqueness rule.

> **The AI side: 5 of 5 PASS. All twelve §7 steps clean on each. Including §7 step 12a gen_ai-completeness checks.**

Mike turned to the team. "Witness-mode test, too. The FDA inspector runs the verifier without the master key — they're not Helmstad's IKM custodian. Witness mode per §7 executes steps 1, 2, 3, 3a, 4, 5, 6, 10, 11, 12, and 12a; it skips steps 7, 8, and 9 because those require IKM access. The output is `Status: PASS-STRUCTURALLY, key-bound verification skipped`."

He typed:

```
herald-verify --tenant=helmstad-trial-screen --service=nsclc-phase2 \
  --date=2026-04-02 --entry-id=nsclc-2026-04-02-mgh-00018
```

No `--master-key` argument. The terminal printed:

```
Status: PASS-STRUCTURALLY, key-bound verification skipped
```

"Witness mode confirms structural integrity — format-version, HKDF-inputs digest, genesis hash, tenant-id character class, per-entry binding, structural walk, Merkle recomputation, signature verification, cadence/dev-mode check, gen_ai-completeness — without requiring IKM access. The FDA inspector can run this on a laptop you provide them, with the verifier binary downloaded from the §10.26 Cosign-signed release artifact channel under reproducible-build discipline. They never see your IKM. They get a `PASS-STRUCTURALLY` result that confirms everything except the per-entry MAC chain, and they have no need to verify the per-entry MAC chain because the structural integrity plus the signed daily seal plus the published public key plus the IQ/OQ for CloudHSM is enough for FRE 901(b)(9) authentication of the process."

Mike ran witness-mode on each of the five sample entries. Five PASS-STRUCTURALLY results. Five exit-code-0 outputs per §10.12.

> **The AI side under witness mode: 5 of 5 PASS-STRUCTURALLY. The inspector-facing verification path works.**

"Now backwards," Dawn said. "For each decision, find the CRO ingestion entry that brought the source data into our warehouse. Then ask Quintessa whether the source records are still on their side."

Chen took over. For each of the five entries, the CRO ingestion entry was straightforward to find — foreign key in the warehouse, sealed in a chain entry of type `cro-ingestion`. The PGP signature verified. The tarball SHA-256 matched. All five.

Then he pulled out his laptop and emailed Quintessa's audit-liaison address. The team had pre-arranged the Quintessa contact with Dr. Østergaard's office. The email asked Quintessa to confirm whether the source EHR records that fed each of the five decisions were still in their environment.

The reply came back in 22 minutes. Quintessa's audit liaison was named Marisol Vega. She confirmed that two of the five decisions had source records still in Quintessa's archive (the April 28 and May 3 decisions). One had source records that were partially available — the patient's lab values were in the archive but the tumor staging assessment had been corrected after Quintessa's snapshot, and the corrected version was in the source EHR but not the snapshot Quintessa had sent Helmstad. Two of the five had source records that had been removed from Quintessa's staging environment 90 days ago, per Quintessa's retention policy.

> **The backward trace: 2 of 5 clean. 1 of 5 with a corrected-source discrepancy. 2 of 5 blocked by Quintessa's retention.**

Chen read the Quintessa email out loud to the room.

"The April 15 decision," Chen said. "MSKCC. Patient's tumor staging was corrected after Quintessa's snapshot. The classifier ran on staging T3N2M0. The corrected staging in the source EHR is T3N2M1. The corrected staging would have changed the eligibility outcome — M1 disease is excluded from this trial."

The room was very quiet.

Dr. Østergaard, calmly: "What did the CRC do?"

Mike pulled up the chain entry. The classifier had returned `eligible` at 0.91 confidence. The CRC had reviewed and accepted at 16:44 local on April 15. The reviewer reason code was `criteria-match-confirmed`.

Dr. Østergaard: "Was the patient enrolled?"

Mike traced forward. The patient had been enrolled in the trial on April 19. They had completed Cycle 1 of the study drug. Cycle 2 was scheduled for May 14.

The room was very quiet.

Dr. Østergaard: "What does the EDC say about the patient's staging at enrollment?"

Mike could not access the EDC directly — the EDC was Quintessa's. He emailed Marisol Vega again. She replied in 9 minutes. The EDC at enrollment had recorded the patient's staging as T3N2M0, matching the classifier's input. The corrected staging at T3N2M1 was an EHR-side correction that had been applied at MSKCC on April 17 — two days after enrollment.

"So the EHR was corrected after enrollment," Mike said. "The trial-visit data in the EDC reflects the staging at the time of enrollment. The classifier saw the staging at the time of screening. The discrepancy is the EHR correction that happened in the two-day window between screening and enrollment."

Dr. Østergaard: "Was the EHR correction propagated to Quintessa?"

Marisol's third email, 14 minutes later: the correction had been propagated to Quintessa on April 18 but not back-fed to the eligibility-screening pipeline. Quintessa's screening pipeline is one-way. Once a candidate is marked eligible and forwarded to Helmstad, the screening pipeline does not re-process them. The correction would have shown up if the patient had been re-screened, but they were not re-screened — they were already enrolled.

Dr. Østergaard, very calmly: "We need to follow up on this patient. Today."

Dawn: "Yes. That is a clinical-quality issue, not an audit finding. The audit finding is that the screening pipeline does not have a re-screening loop for source-data corrections in the window between screening and enrollment. The chain proves what the classifier saw; the chain does not prove the data the classifier saw was current as of enrollment. That is exactly the §1.2 epistemic-scope distinction — chain proves what the AI said at time T; chain does not prove the statement remained accurate at any later time T'."

> **⚠️ Gap-008 (audit finding) plus Clinical-Quality-001 (medical-monitor follow-up) — Source-data correction in the screening-to-enrollment window is not re-evaluated**
> The classifier sees the candidate's data at the moment of screening. If the source EHR is corrected between screening and enrollment, the correction is not re-evaluated by the classifier. The CRC's accept decision is bound to the data the classifier saw, not to the corrected data. **The chain operates exactly as spec §1.2 describes — it proves what the AI said at the moment of screening; it does NOT prove that what the AI said remained accurate as the source EHR evolved between screening and enrollment.** The patient in question may have been enrolled into a trial they would have been excluded from under the corrected staging. The chain itself is conformant; the institutional control gap is the absence of a re-screening loop. This is a process-design finding, not a chain-integrity finding. The Daubert testability of the chain (§1.1) is unaffected — the chain's PASS still proves chain integrity. The §1.2 epistemic line is what the report needs to make explicit so the inspector reads chain-proves-X-not-Y correctly.

Dr. Østergaard wrote it down. "We will pursue this with the medical monitor today. The audit finding is what we asked you to find. The clinical follow-up is mine."

Dawn turned to the team. "What does the chain prove and not prove here, in the §1.2 frame?"

Mike read out the §1.2 line. "Chain proves what the model said at the moment of screening (a). Chain proves the record was not tampered with after capture (b). Chain does not prove the input was clinically accurate at the moment of screening (c). Chain does not prove the input remained accurate through enrollment (the screening-to-enrollment window). Chain does not prove the decision complied with the trial's clinical-quality SOPs (d). Chain does not prove the decision was free of bias (e). The April 15 finding sits at (c) — the input was not accurate at the moment of screening, because the source EHR was corrected two days later but the correction had been pending in the EHR's review queue at the moment of screening and the snapshot Quintessa sent us had captured the pre-correction state."

Dawn wrote that down. *§1.2 (c) — input accuracy. Not a chain finding. A clinical-quality finding with chain-derived evidence.*

The §10.13 evidentiary artifacts package would carry this exchange as part of the institution's litigation-support documentation. The chain would reproduce the full trace; the §1.2 line would distinguish what the chain proved from what the chain did not. The §5.2 best-evidence posture — captured JSON content-bearing form, canonical bytes integrity-bearing form, both originals under FRE 1001(d) — would carry the institution's discovery-production framework if the case ever became contested.

"Forward trace," Dawn said. "Of the five decisions, how many trace forward to enrollment?"

Mike pulled up the enrollment records. Four of the five had been enrolled. The fifth — the April 21 University of Chicago decision — had been classified as eligible but the patient had declined to participate. So the forward trace to enrollment was 4 of 5.

"Of the four enrolled, how many trace forward to clinical-visit data in the EDC?"

Mike worked with Marisol to confirm. Two of the four had completed Cycle 1 visits with EDC data captured. One had been enrolled but had not yet had their first visit. One had been enrolled, had Cycle 1 visit data captured, but the EDC extract for that site had been delayed because of a CRO-side processing backlog.

> **The forward trace: 4 of 5 to enrollment. 2 of 5 to clinical-visit data in the EDC.**

Dawn wrote the summary on the whiteboard:

```
AI side reconciliation:        5/5 PASS
Backward to CRO source:        2/5 clean
                               1/5 corrected-source discrepancy
                               2/5 blocked by CRO retention
Forward to enrollment:         4/5 enrolled
                               1/5 patient declined
Forward to clinical-visit EDC: 2/5 captured
                               1/5 not yet visited
                               1/5 delayed in CRO processing
```

The room was quiet.

Dr. Østergaard, after a beat: "That is the picture. The AI side is solid. The boundary at the CRO source is the inspection risk. The forward trace into the EDC is fine when it works but it is not under our chain. Dawn — that is the framing for the report."

"Yes. That is the framing."

The CQO, on Zoom: "I want it on record that the AI-side reconciliation is 100% and that the backward and forward traces are mixed. The mixed-ness is not a TesseraSeal failure. It is the boundary."

Tom was writing the language down on his iPad.

---

### 😬 3:45 PM — Friction Between the AI Team and Clinical Operations

Dr. Østergaard had asked the clinical operations lead — a person named Sandra Mendelsson — to join the team for the friction conversation he knew was coming. Sandra had run clinical operations at Helmstad for five years and at two other sponsors before that. She had been nominally supportive of the eligibility-classifier project but had pushed back on extending the chain into the CRO data feed. Her position was that Quintessa would not accept it and that pushing for it would damage the relationship.

She had been listening to the reconciliation results from her office on Zoom and had walked over for the live conversation.

"I want to be honest about something," Sandra said. "The AI team has been disciplined about the chain on the model side. I think the four-month run is a real accomplishment. But the ask to extend the chain into Quintessa is unrealistic. Quintessa is not going to instrument their pipeline for our chain. That is not in our contract. That is not what they sell. We use Quintessa because they are the global Phase III CRO and they handle 14 of our trials. Asking them to instrument is not a conversation that ends in yes."

The AI team lead — Dr. Reisch — pushed back.

"I am not asking Quintessa to instrument. I am asking that the contract require Quintessa to provide hash-bound source records, with their own integrity attestation, that we can store as referenced artifacts in our chain — `audit.external_artifact.*` per spec §10.19, with `kind`, `identifier`, `sha256`, `received_at_utc`, `source_party = 'quintessa-research'`, and `evidentiary_role = 'chain_of_custody_handoff'`. If we ever move to a model-handover relationship — say, Quintessa starts supplying a CRO-side AI screening assist — we'd add `audit.model_handover.contract_id`, `contract_version`, and `contract_hash_sha256` per §10.21's Round-17 M&A-G2 close-out. That is a contract clause, not an engineering integration. The Northbridge Bank pattern from last quarter — the bank required their KYC vendor to deliver hash-bound records. The vendor agreed because the bank made it a contract requirement. We can do the same."

Sandra: "The bank is a customer of the vendor's retail product. We are a customer of the CRO's full-service trial operations. Different relationship. Different leverage."

Dr. Reisch: "I understand the leverage difference. I am still asking that we make the contract change for the next renewal."

Dawn looked at Dr. Østergaard.

Dr. Østergaard, calmly: "The right answer is a contract update, not an engineering integration. Sandra is right that Quintessa will not instrument. Hannah is right that a hash-bound delivery clause is a contract change, not an engineering ask. Both can be true. The next contract renewal is in November. We will draft the clause this quarter and put it on the November agenda."

Sandra, quieter: "If we draft it carefully — with their language — they may agree. Hash-bound delivery is not a heavy ask. The integrity attestation language is."

Dr. Østergaard: "We will draft it carefully. Sandra, you and Hannah will lead the drafting. I want a clause we can ship to legal by end of June."

Sandra nodded. Hannah nodded.

Dawn turned to Hannah. "What does the clause look like, drafted to spec?"

Hannah pulled out her notebook. "The framing is `audit.external_artifact.*` per spec §10.19. Quintessa delivers a tarball of de-identified candidate records, signed by their PGP key. Today our chain entry records `kind = 'cro_sftp_tarball'`, `identifier`, `sha256`, `received_at_utc`, `source_party = 'quintessa-research'`, `evidentiary_role = 'chain_of_custody_handoff'`. The contract clause we want adds Quintessa-side commitments: (a) a Quintessa-issued integrity-attestation accompanying each tarball, signed by Quintessa's HSM under a documented signature key; (b) a Quintessa-side retention floor for the source records the tarball was built from, set to longer than the §10.20 deployment-window-plus-investigation-buffer pattern — for clinical trials, that's the trial duration plus the post-trial-close investigation window, which for our Phase II NSCLC trial is at least 3 years; (c) a notification clause if Quintessa's source-record retention is exceeded — they tell us before they purge so we can request specific records before the deadline. None of that is a chain change on Quintessa's side; it's all institution-side discipline. We just want it bound contractually."

Sandra: "And the Quintessa attestation — what does that bind? They're not chain-instrumented."

Hannah: "Per §10.21 cross-vendor model-handover schema's bidirectional cross-anchor pattern, applied by analogy to data delivery rather than model delivery. The Quintessa attestation is a SHA-256 over the canonicalized source-record set the tarball was built from, plus a signature under Quintessa's documented key. Helmstad's chain entry binds the tarball hash; the Quintessa attestation binds the source-record-set hash. The two together let us verify, post-delivery, that the records Quintessa delivered match the records they had on their side. The §10.21 contract-binding sub-attributes — `audit.model_handover.contract_id`, `contract_version`, `contract_hash_sha256` per Round-17 M&A-G2 — would by analogy be `audit.external_artifact.contract_*` if we extend the schema, or live as institution-named attributes alongside `audit.external_artifact.*`. Either way, the contract clause is what makes it work."

Dawn wrote that down. *Hannah understands the spec at the level of cross-anchor reasoning. Sandra now sees the contract path and supports it.*

Dr. Østergaard: "Hannah, draft the contract clause language using §10.19 and §10.21 by-analogy framing. Sandra, draft the negotiation strategy. End of June."

Sandra: "Yes. And I will draft it with Quintessa's language as much as I can. Hash-bound delivery they will accept; integrity-attestation language they will need to think about. We will frame the integrity-attestation as `data-quality-attestation` in their language and map it to §10.19 / §10.21 in ours."

Hannah: "Good. That works."

Dawn wrote that down too. *Friction between AI discipline and clinical-operations pragmatism. Mediated by the VP. Resolved as a contract action with named owners, a date, and a spec-cited drafting framework. That is what mediation looks like when the principals are senior enough to make the decision and the technical lead understands the spec at the level of cross-anchor reasoning.*

The friction did not return for the rest of the day.

---

### 🔍 4:30 PM — The Inspector's Question

Dr. Østergaard had asked Dawn the question on the kickoff call three weeks ago. He asked it again now, in the conference room, with Tom and the CQO on Zoom and the team gathered around.

"If the FDA inspector picks a random patient and asks me to demonstrate that the AI eligibility decision was correct — what can I show?"

Dawn had been writing the answer in her head all day. She gave it now.

"You can show seven things, in this order.

"One — the chain entry for the eligibility decision. The wire-form attributes per spec §4.4 attribute table — `ffiec.chain.payload_hash` as 32 raw bytes, `ffiec.chain.key_fingerprint` as 16 raw bytes per §3, `ffiec.chain.tenant_id`, `ffiec.chain.captured_at` with nanosecond precision, `ffiec.chain.chain_kind = 'model_call'`, `gen_ai.request.model`, `gen_ai.response.model`. Sealed under per-event HMAC-SHA-256 per §4.1. The verifier output resolves in four seconds and produces twelve §7 verification steps. The inspector can run the verifier on a laptop you provide them — and per spec §10.26 reference verifier distribution, you cite the implementation, version, and verification key the institution uses to authenticate the binary at the moment it runs. CC8.1 names all three.

"Two — the verifier output itself, as a printable artifact. The §7 normative output format — `Status: PASS`, `Step: 12`, `Reason: <text>`. Field labels are exact, line terminator is `0x0A`. Exit code 0 per spec §10.12. Witness-mode output (without `--master-key`) reads `Status: PASS-STRUCTURALLY, key-bound verification skipped` per §7 — for the FDA inspector who does not hold Helmstad's IKM but wants to confirm structural integrity, witness mode is the conformant path.

"Three — the model-card hash, which resolves to a tuple of (weights, model card document, validation report). Plus the `audit.model_handover.*` family per spec §10.21 — provider identity, model_id, model_version, `model_artifact_sha256`, `model_card_sha256`, `fairness_audit_report_sha256`, `audit_report_languages`. The inspector can ask to see the model card document; you produce it from the registry; you recompute the SHA-256 in front of them; the hash matches the chain entry's `model_version` field. The §10.20 training-data retention floor is named in CC8.1.

"Four — the protocol document hash. Same drill. The inspector asks to see protocol version 4.2. You produce it from Veeva. You recompute the SHA-256. It matches the `criteria_doc_hash` field, which is bound under the per-event MAC per §4.1 inviolate property 7 (canonical form covers application content; chain stamp fields excluded).

"Five — the reviewer's accept decision, with the SSO subject, the reason code, the timestamp, and the redacted free-text justification (`audit.redaction.disposition = 'redacted_at_sdk'` per §10.22). All sealed. The inspector can ask who the reviewer is — the SSO subject resolves through Okta to the CRC's institutional identity. The §5.2 best-evidence posture: the captured JSON is the content-bearing form; the canonical bytes per RFC 8785 JCS are the integrity-bearing form; both originals under FRE 1001(d).

"Six — the daily Ed25519 seal from CloudHSM, with the public key and the IQ/OQ document. The seal record under spec §4.2 schema — `merkle_root`, `algorithm`, `public_key_id`, `key_versions`, `hkdf_inputs_digest`, `signature`, `cadence`, `signed_at`, `sign_payload_version = 'v1.0b'`. The HSM is FIPS 140-2 Level 3 per §10.5. Strict canonical Ed25519 per §4.3 and RFC 8032 §8.4. The §10.17 partition-ceremony attestation is on chain — `chain.partition_ceremony_attended` events with signatories' `entity_affiliation` per Round-17 M&A-P1. The inspector can verify the seal independently from the published public key.

"Seven — the §10.13 evidentiary artifacts retention package. Helmstad's policy is 25-year retention for trial records under Part 11. The chain entries are stored in a managed retention envelope. Per spec §10.13, this is the documentary evidence the institution retains for the chain-data retention period — SDK version manifest, SDK source-code hash and SLSA build attestation, HSM configuration including FIPS level and rotation history, daily seal-job logs with HSM-signed `signed_at` values, change-management records, verifier output for the period showing PASS for each tenant-day. Per spec §10.14, NTP-discipline timestamp foundation; RFC 3161 trusted-timestamp integration is RECOMMENDED but not REQUIRED for v1.0. CC8.1 cross-references each evidentiary item to the §10.13 list per §10.18.

"That is a Part 11 defensible evidence pack. Seven artifacts, each tied to a normative spec section by number. The inspector can independently verify each one in under twenty minutes. Per spec §1.1 Daubert posture — testability through §7, peer review through Apache-2.0 reference implementation, known error rate through the §1.3 EUF-CMA / second-preimage / EUF-CMA composition, general acceptance through FIPS standards.

"If the inspector asks where the source EHR data came from — you point to the §10.19 chain-coverage map, the SFTP delivery captured under `audit.external_artifact.*`, the PGP signature, the tarball SHA-256, and the Quintessa SOC 2 Type II report (also `audit.external_artifact.*`-anchored on receipt). The provenance line stops at the SFTP boundary. Quintessa is named in the §10.19 map under category 4 (third-party out-of-reach). Quintessa is the responsible party upstream of that boundary.

"If the inspector asks whether the source EHR was correct in the originating site's clinical system — you say that is outside Helmstad's chain of custody, per spec §1.2 epistemic scope. The chain proves what the AI saw at time T; the chain does not prove the AI's input was clinically accurate. The originating site is the responsible party for source EHR integrity. Quintessa is the responsible party for extraction integrity. Helmstad is the responsible party from the SFTP boundary forward — and the §10.19 chain-coverage map names that boundary in writing.

"If the inspector pushes — and they will — on the April 15 patient with the corrected staging that we found this afternoon, you say: we identified that case during a pre-inspection audit. We have notified the medical monitor. The corrective action is in flight. The clinical-quality follow-up is documented. The audit finding has been raised: the screening pipeline does not have a re-screening loop for source corrections in the screening-to-enrollment window. **This is a process-design finding, not a chain-integrity finding** — the §1.2 epistemic line distinguishes them. The remediation plan is in the CAPA system.

"If the inspector asks what about the EDC — you say the EDC is operated by Quintessa under their SOC 2 and is named in the §10.19 chain-coverage map under category 4 (third-party out-of-reach). You provide the SOC 2 report. You explain the vendor-management dependency. If the EDC migrates to a streaming change-stream mirror later, §10.16 SaaS-edge connector lag bounds will apply and the four numbers will be on the runbook.

"If the inspector asks what about Argus — you say the Argus audit trail is mutable by DBAs and that you have a CAPA in flight to bring it under the chain. You explain the §10.3 architectural target — application-level append-only plus database role-level INSERT-and-SELECT-only — and the timeline. You do not pretend the gap is not there. Per §10.18 runbook cross-referencing, when the CAPA lands, the runbook for the new chain-instrumented Argus pipeline will cite §10.3 by section number.

"That is the inspector-facing posture."

Dr. Østergaard tilted his head once. "If the inspector asks about post-quantum readiness — the FDA inspectorate has been asking some sponsors about quantum-resistant cryptography, especially for long-retention regulatory records?"

Dawn had this one too. "Spec §4.3.2 algorithm-rotation and quantum-readiness commitment. Spec working group commits to publishing an emergency spec patch within 30 days of credible demonstration of a practical attack on Ed25519, HMAC-SHA-256, or SHA-256. Institutions migrate to post-attack algorithms within 180 days for signature breaks and within 90 days for HMAC/SHA-256 breaks. Per spec §4.3.2 dual-algorithm posture, when v1.x ships a post-quantum signature algorithm — Dilithium per FIPS 204 or SLH-DSA per FIPS 205 — Helmstad will operate dual-algorithm posture transitionally. The seal record's `signatures` list per §4.2 schema co-signs under both algorithms; Variant B per-algorithm `sign_payload` per §4.3.2; AND-security per §4.3.2 — both signatures must verify for the seal to be valid. Today our chain runs single-algorithm Ed25519 as the v1.0 default; we monitor the cryptographic-agility roadmap published in `docs/regulator-pack/cryptographic-agility-roadmap.md` for the dated-dual-algorithm-seal mandate effective 2030-01-01. We have time. The §10.21 cross-vendor model-handover schema does not yet require post-quantum binding; that's a v1.x candidate-normative item per §4.1.3 RECOMMENDED for v1.0b."

Dr. Østergaard wrote that down. "Good. The inspector probably will not ask, but if they do."

"And §10.6.1 IKM generation — every IKM was generated inside the HSM by the HSM's internal CSPRNG, the highest-assurance posture per §10.6.1's three conformant patterns. The `master_key.generated` operational event per §10.2 records `'hsm.cloudhsm-classic'` as the RNG type. SOC engagements consume this evidence under our CC8.1 procedure. The 32-byte IKM minimum per §10.6 is enforced at IKM-provisioning time and at SDK-configure time."

"Six on the seven artifacts is the daily Ed25519 seal. The inspector independently verifies the seal — what binary do they run?"

"Spec §10.26 reference verifier distribution discipline. The reference verifier ships in a separate repository under Apache 2.0 with reproducible builds, Cosign-signed release artifacts, per-platform binaries (Linux x86_64 / ARM64, Windows x86_64, macOS x86_64 / ARM64), SHA-256 and SHA-512 manifests, CycloneDX SBOM. Our CC8.1 names the implementation, the version, and the verification key the institution uses to authenticate the binary at the moment we run the verifier — that's the §10.26 three-name citation discipline. The FDA inspector arrives with a laptop running the same binary or one of the conformant clean-room implementations per the Q-28 vendor-conformance attestation procedure. Either is acceptable. The conformance contract is the §7 procedure, not any specific binary."

"What if the inspector wants to use a binary they brought themselves?"

"Then they cite their binary in their working paper — implementation, version, verification key — and run it against our chain. The §7 procedure is byte-exact and the test-vector corpus is public, so any conformant implementation produces the same PASS. We trust the verifier on cosign-signed reproducible-build provenance, not on origin. That's what §10.26 is for."

Dr. Østergaard tilted his head. "If the inspector challenges the Daubert framing — they will not name Daubert, but they will ask about evidentiary weight under FDA hearing standards — what is the structured answer?"

Dawn had this one ready. "Spec §1.1 names the four factors and the chain's response to each. Testability — §7 is byte-exact, the test-vector corpus is public, any third party can falsify a verifier's PASS by producing a tampered chain that the §7 procedure does not reject. Peer review — the spec is developed under the FFIEC working-group process with periodic outside-reviewer drops, the reference implementation is Apache 2.0, the corpus is public, independent experts review both spec text and reference code. Known error rate — per §1.3, per-event MAC has EUF-CMA security under HMAC-SHA-256 / FIPS 198-1, daily Merkle seal has second-preimage resistance under SHA-256 / FIPS 180-4 / RFC 6962, HSM signature has EUF-CMA under Ed25519 / FIPS 186-5; a verifying false-negative requires simultaneous compromise of three independent custody layers per §1.4, plus the residual SDK-process scenario §1.2 names. General acceptance — every primitive is NIST-standardized, and the combination of HMAC-chained event records under a tenant-bound key plus a daily Merkle root signed by an HSM is standard in deployed audit-log systems including Certificate Transparency and Trillian. The §1.1 grounding is informative, not normative — it does not add requirements; it exists so a witness laying foundation under FRE 702 / FDA-hearing analog can answer the four questions from the shipped artifacts rather than from internal documents."

"And if the inspector asks what we know we don't prove?"

"Spec §1.2 epistemic scope. Three explicit non-claims: factual accuracy, policy compliance, freedom from bias. We hold the line that the chain is the integrity foundation, not the truth foundation, and Helmstad's other evidence regimes — the medical monitor, the Quality System, the SOPs — answer the truth questions. The §1.2 line reduces witness-stand confusion: the chain says the model said X at time T, full stop."

Dr. Østergaard wrote both down. "That is the answer."

"And the second layer of the §1.4 compositional security argument — what does the inspector ask if they want to know about the residual risks?"

"The §1.2 fourth class — application-process compromise. Adversary F per the threat model. A compromised process holding the live session key produces chain entries that verify as PASS for as long as the compromise persists. Per-event MAC is intact — same IKM, same session key. Daily Merkle seal signs those entries normally. The institution's IKM is NOT compromised, the ledger is NOT compromised, the HSM is NOT compromised — yet the chain produces a verifying record of events the legitimate AI agent did not generate. This is a forward-only attack window: past chain entries are unaffected. The window is bounded by the institution's host-hardening, intrusion-detection, and master-key-rotation controls. The chain composes alongside these controls; the chain alone does not defend against an attacker who has root on the SDK's host. Helmstad's CC8.1 names the host-hardening posture, the IDS subscriptions, and the rotation cadence as the compensating controls."

"And SDK-side enforcement of `gen_ai.{request,response}.model`?"

"Per spec §4.4 normative requirement, our SDK refuses to emit a chain entry whose attribute set includes any `gen_ai.*` attribute AND lacks either `gen_ai.request.model` or `gen_ai.response.model` non-empty. The refusal is at SDK-write time — entry rejected before MAC compute and before any wire emission. Raises a `GenAIModelIdentifierMissing` exception in our Python idiom so the operator's error path surfaces the misconfiguration immediately. The §7 step 12a verifier check is defense-in-depth; the SDK-side refusal closes the source so a misconfigured pipeline cannot silently produce chains that fail at audit time. We've never seen the exception in production — by construction, but the test suite exercises it on every CI run."

Dr. Østergaard had been listening with both hands flat on the table. He nodded once when Dawn finished.

"That is the posture," he said. "Tom — please make sure that articulation is in the report verbatim."

Tom was already typing it.

The CQO on Zoom: "Dawn — that articulation. Use exactly that language. The seven-artifact evidence pack is the inspection-day playbook."

Dawn: "Use it."

---

### 🌆 5:30 PM — Debrief

The conference room was quieter than it had been at kickoff. Dr. Østergaard had ordered coffee and pastries. The CQO was on Zoom from his car — he had a hard stop at 6:30. Sandra was on Zoom from her office. Hannah and Devansh were in the room.

Dawn stood up at the whiteboard.

"Two findings sections," she said. "AI side. Legacy side."

She wrote on the whiteboard.

#### ✅ What They Confirmed on the AI Side

| # | Finding | Spec section |
|---|---|---|
| 1 | Chain integrity holds at 4 months and at 90 days across the v1.0a → v1.0b sign_payload dispatch boundary. Verifier resolves any entry in ~4 seconds in 12 §7 steps. Exit code 0 per §10.12. | §4.1, §4.2, §4.3, §7, §10.12 |
| 2 | CloudHSM signing infrastructure is documented to Part 11 standard. FIPS 140-2 Level 3 per §10.5. IQ/OQ in Veeva, signed and dated. Partition ceremony chain-coupled per §10.17 with `entity_affiliation` per Round-17 M&A-P1. | §10.5, §10.17 |
| 3 | Model card, validation report, and provider hash-bound to every entry via `audit.model_handover.*`. Training-data retention floor at deployment-window plus 90-day investigation buffer. | §10.20, §10.21 |
| 4 | Protocol document hash binds eligibility criteria to the decision. ALCOA+ Original on the criteria side. Bound under §4.1 inviolate property 7. | §4.1, §5.2 |
| 5 | Reviewer decisions are sealed end to end under SSO subject. `audit.redaction.*` family bound under per-event MAC per §10.22 normative posture (pre-MAC at the SDK boundary). | §10.22, §4.4 |
| 6 | AI service-account IAM is chain-coupled. Two-of-three approval per §10.5 separation of duties. Rotation crossing seal boundary handled per §10.10. `master_key.rotated` and `master_key.rotation_observed` per §10.2. | §10.5, §10.10, §10.2 |
| 7 | CRC reviewer decisions are sealed under SSO subject attribution. Constant-time fingerprint compare per §10.8 at §7 step 8 before MAC compute. | §10.8, §7 |
| 8 | API gateway and fallback paths are sealed. Routing decisions captured per §4.4.1 — `audit.routing.attempt`, `failover`, `refused`. CloudHSM maintenance windows handled per §4.3.1. No silent gaps. OTLP transport identification per §4.4.3; collector pass-through per §4.4 and §4.4.4. | §4.4.1, §4.3.1, §4.4.3, §4.4.4 |

**AI side: 0 Gaps. 0 Partials. ALCOA+ defensible end to end on the surface area the chain covers.** Per spec §1.2 epistemic scope, "covers" means the chain proves what the AI said and that the record was not tampered with after capture; the chain does not prove the AI's statement was clinically accurate, that the statement complied with FDA Quality System Regulation, or that the statement was free of bias. Those are separate evidence regimes that compose alongside the chain.

#### ❌ What They Found on the Legacy Side

Each row carries the spec section that *would* remediate the gap if the system were brought under the chain. Today, every legacy-side row is named in the §10.19 chain-coverage map under either category 2 (institutional system not yet chain-instrumented) or category 4 (third-party out-of-reach), so the boundaries are explicit per §10.19 even where the controls themselves are not chain-grade.

| # | Finding | Spec section that would remediate |
|---|---|---|
| 1 | Argus DBAs have UPDATE permission on the Part 11 audit-trail table. No out-of-band integrity binding. | §10.3 append-only enforcement (would apply if Argus were chain-instrumented); §10.19 §2 today |
| 2 | LIMS audit trail has the same DBA-mutable shape as Argus. | §10.3, §10.19 §2 |
| 3 | Veeva Vault QMS has a "Late Effective Date" workflow that allows up to 30 days of effective-date back-dating without showing in the standard audit-trail report unless the user knows the filter. | §10.18 runbook cross-referencing (would surface the filter discoverability); §10.19 §2 |
| 4 | CTMS Site Monitoring Visit Summary field history is disabled. Salesforce CRM-pattern. | §10.22 redaction discipline; §4.4.6 connector_source if mirrored; §10.19 §2 |
| 5 | EDC extract pipeline has no chain. Trial-visit warehouse is loaded from CRO SFTP without HMAC binding. | §4.4.6 connector_source; §10.19 §4 today |
| 6 | Argus-to-warehouse ETL has no chain. Quarterly reconciliation is a count check, not a contents check. | §10.3, §10.19 §2 |
| 7 | CRO source-side history is outside the chain. Provenance stops at the SFTP boundary. | §10.19 §4 (named in the chain-coverage map; not a hidden boundary) |
| 8 | EDC audit trail is in the CRO's SOC 2 scope, not Helmstad's. Vendor-management dependency. | §10.19 §4 |
| 9 | Source-data correction in the screening-to-enrollment window is not re-evaluated by the classifier. April 15 patient. Clinical follow-up in flight. | §1.2 epistemic-scope distinction (chain proves what AI saw; not whether the input was current); process-design CAPA |
| 10 | Three engineers can disable CloudTrail logging on the CRO ingestion S3 bucket. | §10.19 §1 (chain-instrumented system; coverage-map should name the dependency) |
| 11 | CTMS-as-KOL-CRM has DLP-only protection on PHI-prohibited free-text fields. | §10.22 redaction discipline; §10.19 §2 |
| 12 | Quarterly access review attests group membership, not table-level modification. | §10.3 (would close the architectural gap); §10.19 §2 |
| 13 | Quintessa's retention policy removed source records for 2 of 5 reconciliation samples (90-day rotation). | §10.19 §4 (boundary named; vendor-side retention floor is contractual) |

**Legacy side: 4 Gaps. 6 Partials. 3 boundary-by-design items handled per §10.19.** The Partials are items where there is *some* protection — Veeva's audit trail is real if you know the filter (§10.18 cross-referencing would surface the filter's discoverability); Salesforce field history is on for some fields; the SFTP signature does verify upstream provenance to a point — but the protection is at the discretion of the system operator or the vendor and is therefore not chain-grade evidence. The remaining items are Gaps in the strict sense. The boundary-by-design items are not Gaps and not Partials — they are §10.19 chain-coverage-map disclosures handled to spec.

Dawn paused at the table for a moment before she began the Gap walk.

She wanted to be precise about the §10.16 framing. The EDC extract today was a daily file pipeline — not a streaming SaaS-edge change-stream connector. §10.16 strictly applies to mirror connectors that subscribe to a SaaS platform's change stream and replicate each captured record into a chain-instrumented store. Helmstad's EDC pipeline did not match that pattern. So the §10.16 four-number wording test did not apply by the letter.

But the spirit of §10.16 — that connector lag has to be quantified, that imprecise wording is non-conformance not Nit, that an examiner cannot test what an institution has not committed to — applied to any pipeline crossing the chain boundary. The §4.4.6 `audit.connector_source.*` family and the §10.16 four-number discipline were what the EDC pipeline would need when (not if) it migrated to a streaming model. The November contract renewal with Quintessa was the right moment to raise the §4.4.6 / §10.16 framing.

For the report, she would frame the EDC finding under §4.4.6 (the family that names the connector-source attribution) and §10.16 forward-readiness (the framework the EDC pipeline would need to meet under streaming). Today's status was Gap-006 — no chain — handled per §10.19 §4 (third-party out-of-reach in the chain-coverage map). Forward-readiness status was a Partial against §10.16 in the spirit-rather-than-letter sense, with the contract renewal as the named action.

Dawn ticked off the Gaps on her fingers as she summarized.

"Argus DBA write access is a Gap against the §10.3 architectural target. Veeva late-effective-date workflow is a Gap. EDC extract has no chain — Gap. CTMS overwrites — Gap.

"The Partials: CRO source-side history is a Partial in the layperson sense, but per §10.19 it's a boundary-by-design — the SOC 2 is anchored, the SFTP signature is bound under `audit.external_artifact.*`, and the tarball hash is on chain, but the source-side itself is outside our perimeter and the §10.19 chain-coverage map names that explicitly. The lab values discrepancy investigation is Partial pending the medical monitor's follow-up; the chain itself is conformant per §1.2 epistemic scope. Quintessa rotation policy is a Partial because their 90-day retention is contractually permitted but limits backward traceability. IAM bifurcation for non-AI users is a Partial — the legacy AD is not chain-coupled but the access reviews exist. Retention policy variance between Helmstad's 25-year requirement and Quintessa's 90-day staging — Partial; bound to the §10.20-style retention-floor reasoning if Quintessa ever supplies a model. Audit-trail filter discoverability in Veeva — Partial because the filter exists but is not the default; §10.18 cross-referencing would surface it."

#### 🔁 Side-by-side comparison

| Dimension | AI side (eligibility classifier) | Legacy side (Veeva, Argus, CTMS, EDC, CROs, LIMS) |
|---|---|---|
| Record integrity | Sealed via HMAC chain + daily Ed25519 seal on CloudHSM | DBA-mutable in Argus and LIMS; Veeva back-datable; CTMS overwrite-able |
| Identity coupling | Service IAM chain-coupled; CRC SSO chain-coupled | Argus and LIMS DBAs not chained; commercial CTMS users not chained |
| Reconciliation | Cross-checkpoint sealed events; verifier in 4 sec | Quarterly count checks in Argus; manual reconciliation elsewhere |
| Retention | Helmstad 25-year, sealed | Quintessa 90-day staging; SQL Server 30-day txlog |
| Verifier | `herald-verify` resolves in 12 §7 steps; cited per §10.26 (implementation, version, verification key) | No verifier; audit by inspection |
| FDA Part 11 audit-trail | Defensible per draft AI/ML guidance, §1.1 Daubert framing, §1.3 security definitions, §1.4 compositional security | Defensible only on shape, not on integrity |
| ALCOA+ Attributable | Yes (chained SSO subject under §4.1) | Yes (Argus authentication; CTMS Salesforce subject) |
| ALCOA+ Original (model output) | Yes (per §5.2: captured JSON content-bearing form + canonical bytes integrity-bearing form, both originals under FRE 1001(d)) | n/a |
| ALCOA+ Original (source data) | Boundary named in §10.19 chain-coverage map at CRO SFTP | Quintessa SOC 2 anchored under `audit.external_artifact.*` per §10.19 |
| ALCOA+ Contemporaneous | Yes (200 ms chain latency; nanosecond `captured_at` per §4.4; NTP discipline per §10.4) | Mostly yes; CTMS allows post-hoc edits |
| BIMO inspection posture | Strong — 7-artifact evidence pack tied to §7, §10.5, §10.13, §10.14, §10.17, §10.21, §10.26 | Mixed — vendor-management dependencies named in §10.19 |

Dawn paused on the table for a beat longer than the rest.

"The chain proves what the model said. The chain proves what the criteria document was. The chain proves who the reviewer was. The chain proves the reviewer's decision and reason. That is the FDA's question, on the AI side. Spec §1.2 names this as what the chain proves — (a) what the AI said at a specific time, and (b) that the record was not tampered with after capture.

"The chain does not prove what the source EHR said before Quintessa touched it. The chain does not prove what the EDC's audit trail says about a trial visit. The chain does not prove that an Argus row was not modified by a DBA with UPDATE permission. Spec §1.2 also names what the chain does not prove — (c) that the AI's statement is factually accurate, (d) that the AI's statement complied with policy, (e) that the AI's statement is free of bias. The same epistemic line applies, with FDA Quality System Regulation in the place of fair-lending.

"That is the inspection risk. Six weeks. Argus first — bring it under the chain or, at minimum, revoke UPDATE on the audit-trail table per the §10.3 architectural target. EDC contract clause for November renewal — `audit.external_artifact.*` plus `audit.model_handover.*` if the relationship deepens. Veeva late-effective-date filter as a CAPA with §10.18 runbook cross-referencing. CTMS field-history funding for the next budget cycle, with the §4.4.6 / §10.16 connector path as the design target. The April 15 patient is medical-monitor work and is in flight as of this afternoon — chain-integrity finding versus process-design finding distinguished per §1.2."

The CQO on Zoom: "That is the line for the inspection-prep memo. Dawn — Tom — thank you. The seven-artifact evidence pack is going to be how we run the first morning of the inspection. The legacy-side findings are going to be the second morning. We will know what is coming."

Dr. Østergaard: "Permission to share the bifurcation framing in the BIMO prep documents to the executive committee on Friday."

Dawn: "Permission granted. We will send the formal report by end of week. The bifurcation framing is the framing. Use it."

The CQO on Zoom: "One more thing. I want it on record that this assessment found the AI-side controls to meet or exceed best-known practices observed in deployed AI clinical-decision-support systems audited in the last twelve months. And that the legacy-side controls are at-or-below the median for mid-size oncology sponsors. That is the comparative posture we are taking to the executive committee. Dawn — you good with that characterization?"

Dawn took a beat. "I am, with the same caveat I gave at Mercator. The population of sponsors running chain-grade AI controls in production is small. Helmstad is in the top quartile of a fairly small group. The 'meets or exceeds best-known practices' phrasing is more accurate. I would also add — for the executive committee read — that the legacy-side gaps include items where the gap is known and the remediation has been deferred for budget or contractual reasons. The CTMS field-history-disabled finding has been on the HITRUST-equivalent attestation for prior audits. The Argus DBA finding has been a known industry concern for years. Helmstad is not surprised by these findings. The committee should not be surprised either."

The CQO: "Better. Use that."

Sandra, on Zoom: "Agreed."

Tom was writing the language down verbatim.

Hannah: "One more for the record. The April 15 patient. I want it documented that the audit team found the discrepancy during a pre-inspection readiness audit, not during the inspection. The fact that we found it before the inspector matters for the BIMO read."

Dawn: "Documented. The reconciliation methodology — five-decision random sample with backward and forward trace — is in the report. The April 15 finding is one of the five. It will read as a discrepancy detected by the audit, with corrective action in flight."

Dr. Østergaard stood up, shook Dawn's hand, then Tom's, then went around the room and shook each team member's hand individually. It took ninety seconds. He thanked each of them by name.

"Drive safely. The Massachusetts Pike is bad in the rain."

Dawn looked out the window. It had started raining at some point during the afternoon.

The team packed up. There was the usual quiet shuffle of laptops closing and notebooks going into bags. Diana paused on the way out and asked Rohan a private question about extending the AD access-review process to attest table-level privileges, not just group membership. Rohan answered it. Mike traded contact details with Hannah, who had stayed for the debrief tail. Chen and Devansh exchanged GitHub handles. Luis caught Karthik in the hallway and gave him an unsolicited recommendation about a SQL Server Always Encrypted pattern that would not have prevented DBA write access but would have made it visible faster. Karthik thanked him. Elena thanked Jordan for being honest about the field-history-disabled situation.

Sandra was the last one out. She stopped Dawn by the door.

"I want to say — I was wrong about the Quintessa contract clause. The audit gave me the framing I needed. The November renewal is going to be a conversation, not a fight. Hannah and I will get the clause to legal in June."

Dawn: "Good. The contract change is the right move. The chain into the CRO is not."

Sandra: "Understood."

Dawn sat for a moment after the room had mostly cleared. She thought about the day in the §10.18 cross-referencing frame.

Each finding tied to a spec section. Each spec section tied to an audit procedure. Each audit procedure tied to a SOC engagement test. Each SOC engagement test tied to an examiner workpaper. The chain of evidence — runbook → spec → design doc → audit procedure → SOC → examiner — was load-bearing, and the §10.18 normative cross-referencing requirement was what kept it walkable. A finding that did not cite its spec section was a finding that the institution could not test against; a runbook that did not cross-reference the spec was a runbook that the SOC team could not anchor to the conformance bar.

The Helmstad findings register would meet that bar. Argus DBA write access cited §10.3. LIMS the same. Veeva Late-Effective-Date cited §10.18 itself by self-reference — the runbook discoverability deficiency would be remediated by adopting the §10.18 cross-referencing discipline. CTMS overwrites cited §10.22 redaction discipline plus §4.4.6 connector_source. EDC extract cited §4.4.6 plus §10.16. The April 15 patient cited §1.2 epistemic scope. Each one had a normative section by number; each one would walk the evidence chain to an examiner workpaper.

The chain itself was conformant on every property she had tested. The chain-coverage map per §10.19 named the boundary explicitly. The §1.2 epistemic-scope distinction kept the witness-stand language clean: the chain proved what the AI said and that the record was not tampered with; the chain did not claim more.

She closed her notebook.

Dawn watched the room empty. She gathered her notes, clipped her pen back to the cover, and slung her bag over her shoulder.

Tom held the door for her on the way out. "Round four is in the rear-view."

"Round four is the second bifurcation in five weeks," Dawn said. "Mercator was the first. Helmstad is the same shape with a different inspector. The pattern is real."

"Different inspector. Same architecture."

"Same architecture. That is the line."

---

### 🧪 Post-debrief — What the report will actually say about each finding

Tom and Dawn stayed in the conference room for another twenty minutes after the rest of the team headed for their cars. Tom was drafting the report's findings register on his iPad; Dawn was reviewing it line by line.

"Argus DBA write access," Tom said. "Findings register draft: Gap-001."

Dawn: "Cite §10.3 for the architectural target. Cite §10.19 for the chain-coverage map disposition today. The remediation language is `bring Argus under the chain, or revoke UPDATE on the audit-trail table at the database role layer per §10.3`. Don't call it non-conformance — Argus is not chain-instrumented today, so §10.3 doesn't bind by spec. Call it a Gap against the architectural target."

"LIMS audit trail. Same shape. Gap-002."

"Same citations. Same language."

"Veeva Late-Effective-Date filter. Gap-003."

"Cite §10.18 for the cross-referencing requirement. The remediation is to update CC8.1 to name the filter as a load-bearing audit-trail-discovery item per §10.18. The runbook is real; the discoverability is not."

"CTMS field-history disabled. Gap-004."

"Cite §10.22 for the redaction-discipline architectural target. Cite §10.19 §2 for the chain-coverage map disposition today. Remediation is the §4.4.6 connector path with §10.16 lag bounds when the connector is built."

"Argus access review. Gap-005."

"Cite §10.3 for the architectural target. The Q1 2026 review attests group membership, not table-level modification. Remediation is to extend the access-review protocol to attest `UPDATE-on-audit-trail-table` separately."

"EDC extract pipeline. Gap-006."

"Cite §4.4.6 for the connector_source family. Cite §10.19 §4 for the chain-coverage map disposition today. The boundary at Quintessa is named; the §10.19 disclosure is conformant."

"CTMS Site Monitoring Visit Summary. Gap-007. Same citations as Gap-004 but on the CRA-side workflow rather than the KOL-side."

"Source-data correction in screening-to-enrollment window. Gap-008."

"Cite §1.2 for the epistemic-scope distinction. The chain operates correctly; the gap is in the institution's process design — the screening pipeline does not have a re-screening loop. Remediation is process-design CAPA, not a chain change. The clinical-quality follow-up on the April 15 patient is in flight separately under the medical monitor's responsibility."

"And the boundary-by-design items?"

"Three of them. CRO source-side history out-of-reach per §10.19 §4. EDC audit trail in CRO's SOC 2 scope per §10.19 §4. Quintessa's 90-day retention policy per §10.19 §4 with vendor-side retention-floor reasoning that would come under §10.20 if Quintessa supplied a model. None are Gaps. None are Partials. All are §10.19 chain-coverage map disclosures handled to spec."

Tom wrote each line carefully. The report's findings register would tie every Gap and Partial back to a spec section by number, the way Northbridge's report had. The §10.18 cross-referencing discipline ran in both directions — the runbook should cross-reference the spec, and the audit report should cross-reference the spec.

The final line in Tom's draft was the executive summary's framing.

```
Helmstad's screening pipeline meets the FFIEC chain-of-custody v1.0b
conformance bar end-to-end on the surface area the chain covers.
The legacy stack does not, and does not claim to. The §10.19 chain-
coverage map names the boundary explicitly. The remediation roadmap
is to extend the chain into the legacy stack on a phase basis,
beginning with Argus and the LIMS post-FDA inspection, with the EDC
extract following at the November contract renewal under §4.4.6
discipline, and CTMS following in the next budget cycle.
```

Dawn read it twice. "Use it."

---

### 🧾 Final Assessment Theme

The drive back to the hotel was twelve minutes through Cambridge rain. Dawn had her coffee, refilled, in the cup holder. The wipers were on intermittent.

She thought about the day. About the bifurcation. About Hannah's terminal showing a clean verifier resolution at 9:45 in the morning, all twelve §7 steps green, exit code 0 per §10.12. About Luis pointing at the S3 bucket policy at lunch and saying "three engineers can disable CloudTrail." About the April 15 patient and the staging correction that had moved between screening and enrollment — a §1.2 epistemic-scope finding, not a chain-integrity finding. About Dr. Østergaard's question at 4:30 and the seven-artifact answer she had given, each artifact tied to a normative spec section by number.

She thought about Northbridge, four weeks ago, where everything had been sealed and the audit had been almost boring — until Elena had read §10.16's severity-classification clause aloud and the engagement had landed Finding-001 non-conformance for imprecise lag wording on a Salesforce mirror. The discipline had survived: the wording IS the testable claim.

She thought about Mercator, three weeks ago, where the chain had been on the imaging path and the EHR had been mutable and Patricia had asked the audit to draw the line so she could fund the rest. The §10.19 chain-coverage map framing came from that engagement.

She thought about Stelvio, two weeks ago, where the chain had been on the QC vision and the OT had been on PLCs and Maria had asked the audit to triage three zones. Each zone got its own §10.19 boundary disclosure.

She thought about Atrio, last week, where forty-seven tenants had run on one chain and the verifier had not failed once across fourteen hundred runs. Per-tenant Shape 1 CUEC integrity per §10.23. Multi-tenant SaaS partition ceremony attestation per §10.17 with `entity_affiliation` per Round-17 M&A-P1.

Helmstad was the second bifurcated audit. AI side sealed. Legacy side legacy. Different industry. Different regulator. Same architectural pattern. Same line drawn between the part that decides and the part that supplies the inputs. The §10.19 chain-coverage map was the boundary's normative shape; the §1.2 epistemic-scope distinction was the line between integrity and clinical truth.

The seven-artifact evidence pack would carry the FDA inspection. Each artifact had a normative spec section behind it; each spec section had a verifiable property the chain-of-custody primitives delivered:

1. **Chain entry for the eligibility decision** (§4.4 attribute table; §3 chain_kind enumeration; §4.1 per-event MAC). The wire-form record itself, byte-for-byte preserved per §6 storage rules. Any tampering surfaces as `payload_hash MAC mismatch at seq N` per §7 step 9.

2. **Verifier output as a printable artifact** (§7 normative output format; §10.12 exit-code contract; §10.26 reference verifier distribution discipline). The §7 procedure walks twelve ordered steps; failure produces a specific named reason string with byte-for-byte normative wording. The §10.26 three-name CC8.1 citation lets an FDA inspector reproduce the verifier invocation byte-identically.

3. **Model-card hash** (§10.21 cross-vendor model-handover schema; §10.20 training-data retention vs deployment-window discipline). The `audit.model_handover.*` family captures the provider, model_id, model_version, artifact hash, model-card hash, fairness-audit-report hash. The §10.20 retention floor is at deployment-window plus 90-day investigation buffer; GDPR Article 6(1)(f) legitimate-interests resolution.

4. **Protocol document hash** (§4.1 inviolate property 7 — canonical form covers application content; chain stamp fields excluded). The `criteria_doc_hash` is bound under the per-event MAC; recomputing the SHA-256 from the registry document and comparing to the chain entry's field is the inspector-side verification path.

5. **Reviewer accept decision with §10.22 redaction discipline** (§10.22 normative posture — pre-MAC at the SDK boundary; `audit.redaction.*` family bound under per-event MAC). The §5.2 best-evidence posture: captured JSON content-bearing form, canonical bytes integrity-bearing form, both originals under FRE 1001(d). FRE 1003 admits duplicates of either form unless authenticity is challenged.

6. **Daily Ed25519 seal from CloudHSM** (§4.2 Daily Merkle seal under RFC 6962; §4.3 HSM-rooted root signature with `sign_payload_version = "v1.0b"`; §10.5 FIPS 140-2 Level 3 HSM custody; §10.17 partition ceremony attestation with `entity_affiliation` per Round-17 M&A-P1). Three-layer compositional security per §1.4: per-event HMAC + daily Merkle + HSM Ed25519 signature, each independent, each defending against a different attack class.

7. **§10.13 evidentiary artifacts retention package** (§10.13 documentary evidence list; §10.14 trusted-time integration via NTP discipline per §10.4; §10.18 CC8.1 cross-referencing discipline). SDK version manifest, source-code hash, SLSA build attestation, HSM configuration, daily seal-job logs, change-management records, verifier output for the period.

The thing that made Helmstad interesting was not that the chain held. The chain has held in every audit since Northbridge. The §1.4 compositional security argument — three independent layers (per-event MAC, daily Merkle seal, HSM-rooted root signature) under FIPS-standardized primitives — was load-bearing in every engagement. The thing that made Helmstad interesting was that the FDA inspection in six weeks would test the *boundary* — and the boundary was where the next class of audit findings would live, in this and every other bifurcated deployment for the rest of Dawn's career.

The thing that made Helmstad interesting was not that the chain held. The chain has held in every audit since Northbridge. The §1.4 compositional security argument — three independent layers (per-event MAC, daily Merkle seal, HSM-rooted root signature) under FIPS-standardized primitives — was load-bearing in every engagement. The thing that made Helmstad interesting was that the FDA inspection in six weeks would test the *boundary* — and the boundary was where the next class of audit findings would live, in this and every other bifurcated deployment for the rest of Dawn's career.

The §10.19 chain-coverage map was Helmstad's normative answer to the boundary question. Five categories: chain-instrumented institutional systems, institutional systems not yet chain-instrumented, third-party systems under contractual inspection, third-party systems out-of-reach, external evidentiary artifacts hash-anchored. Each row carried a `coverage_map_version` and an `effective_utc`; each version change emitted `chain.coverage_map_published` operational event per §10.2 with `coverage_map_sha256` so an 18-month-lookback auditor could verify which map version was in force on a given date. The Round-17 M&A-P3 close-out — the version-stamp-and-anchor requirement on the chain-coverage map — was load-bearing for FDA's pre-inspection lookback window the same way it was load-bearing for an acquirer's IT due-diligence lookback.

The §1.2 epistemic-scope discipline distinguished what the chain proved from what the chain did not prove. (a) what the AI said at time T. (b) the record was not tampered with after capture. NOT (c) factual accuracy. NOT (d) policy compliance. NOT (e) freedom from bias. The April 15 patient finding sat at (c) — the input was not accurate at the moment of screening, because a correction was pending in the EHR's review queue and the snapshot Quintessa had sent captured pre-correction state. The chain proved exactly what the chain claimed: the classifier saw staging T3N2M0; the chain did not claim staging T3N2M0 was clinically true. The CRC's accept decision was bound to what the chain saw, not to what was clinically true. The remediation was a process-design CAPA — add a re-screening loop for source corrections in the screening-to-enrollment window — not a chain change. The chain itself was conformant.

The chain proves what the model decided.
The contract has to prove what the model was given.

Dawn picked up her phone at a red light and dictated a one-line note for the report's executive summary.

> *"On the screening path, Helmstad can demonstrate what the classifier said and that the record was not tampered with — per spec §1.2 (b). On every other path, integrity is still vendor-managed — per the §10.19 chain-coverage map's category-2 and category-4 disclosures."*

The light turned green. She put the phone down and drove the rest of the way to the hotel through the rain.

---

### 📎 Appendix — Spec sections cited in this report

The report Dawn and Tom would file by end-of-week tied every determination back to the chain-of-custody specification by section number, so a future reader could walk from a finding to the binding requirement to the audit procedure to the SOC engagement to the examiner workpaper, per spec §10.18 runbook cross-referencing. The full citation list:

- **§1.1 Daubert four-factor grounding** — testability, peer review, known error rate, general acceptance; the FDA hearings analog.
- **§1.2 Epistemic scope** — what the chain proves (a, b) and does not prove (c, d, e); load-bearing for the April 15 patient framing.
- **§1.3 Security definitions** — EUF-CMA on per-event MAC, second-preimage on the Merkle seal, EUF-CMA on the HSM signature.
- **§1.4 Compositional security** — three independent layers; defense-in-depth.
- **§3 chain_kind enumeration** — `audit`, `model_call`, `tool_call`, `routing`, `translation`, `operational`.
- **§3 tenant_id character class and HKDF info parameter** — per-tenant session-key isolation.
- **§4.1 HMAC chain at capture** — per-event construction, fixed-width prev_hash, MAC-IS-payload_hash, canonical form excludes chain stamp, expected_prev_hash in MAC recompute.
- **§4.1.1 Session-key handshake** — Model A IKM-delivered vs Model B session-key-delivered; HSM-resident PRK pattern Helmstad uses.
- **§4.2 Daily Merkle seal** — RFC 6962 with leaf/internal prefix, `(run_id, seq)` ordering, empty-day root.
- **§4.2.1 Cadence** — daily by default; Helmstad's posture.
- **§4.3 HSM-rooted root signature** — Ed25519 strict canonical per RFC 8032 §8.4; sign_payload_version dispatch.
- **§4.3.1 HSM unavailability and notification** — 72-hour SHOULD; Helmstad's CC8.1 covers it.
- **§4.4 OpenTelemetry-native wire** — `ffiec.chain.*` attribute table; `gen_ai.request.model` and `gen_ai.response.model` REQUIRED on model-call entries.
- **§4.4.1 AI routing decisions** — `audit.routing.*` family; required event types per call shape.
- **§4.4.3 OTLP transport identification** — Resource attributes; HTTP headers and gRPC metadata.
- **§4.4.4 Severity for chain-of-custody traffic** — collector pass-through; receiver stamping in 9–20 range.
- **§4.4.6 SaaS-edge connector source attribution** — `audit.connector_source.*` family for any future EDC mirror.
- **§5.2 Best-evidence posture** — captured JSON content-bearing form; canonical bytes integrity-bearing form; both originals under FRE 1001(d).
- **§7 Verification procedure** — twelve ordered steps; normative output format; witness-mode behavior; step 12a gen_ai-completeness.
- **§10.2 Operational events** — `master_key.rotated`, `master_key.rotation_observed`, `chain.partition_ceremony_attended`, `chain.coverage_map_published`, `chain.entity_succession`.
- **§10.3 Append-only enforcement** — application + database role; the architectural target Argus and LIMS would meet under remediation.
- **§10.4 Time synchronization** — NTP discipline; Helmstad documents it.
- **§10.5 HSM custody** — FIPS 140-2 Level 3; CloudHSM Classic conformant; AWS KMS multi-tenant not conformant.
- **§10.6 IKM minimum length** — 32 bytes minimum.
- **§10.6.1 IKM generation requirements** — FIPS-validated CSPRNG; HSM internal RNG is the highest-assurance pattern.
- **§10.7 Software-key adapter exclusion** — production builds do not carry the dev adapter on disk.
- **§10.8 Constant-time comparison** — for fingerprint and MAC compare; `hmac.compare_digest` in Python.
- **§10.10 Rotation crossing the seal boundary** — `key_versions = [old, new]` in the day-after seal record.
- **§10.12 Verifier CLI exit-code contract** — 0/1/2/3.
- **§10.13 Evidentiary artifacts retention** — SDK manifest, source-code hash, HSM configuration, seal-job logs, change-management records, verifier output.
- **§10.14 Trusted-time integration** — RFC 3161 RECOMMENDED for v1.0; NTP discipline today.
- **§10.16 SaaS-edge capture connectors** — four-number quantification; severity-classification clause; non-conformance for imprecise lag wording.
- **§10.17 HSM partition ceremony attestation** — `chain.partition_ceremony_attended`; `entity_affiliation` on signatories per Round-17 M&A-P1.
- **§10.18 CC8.1 and runbook cross-referencing** — runbooks supporting normative spec requirements MUST cross-reference the spec section.
- **§10.19 Chain-coverage boundary documentation** — five-category map; `chain.coverage_map_published` operational event; `audit.external_artifact.*` family.
- **§10.20 Training-data retention vs deployment-window discipline** — retention floor at deployment-window plus 90-day investigation buffer; GDPR Article 6(1)(f) resolution.
- **§10.21 Cross-vendor model-handover schema** — `audit.model_handover.*` family; contract-binding sub-attributes per Round-17 M&A-G2.
- **§10.22 Redaction discipline** — pre-MAC at the SDK boundary; `audit.redaction.*` family.
- **§10.23 Consumer-correlation index integrity** — Shape 1 chain-anchored vs Shape 2 attestation.
- **§10.24 Entity succession** — `chain.entity_succession` operational event; dual-signatures with `entity_affiliation`.
- **§10.25 Run resume and chain-tail acquisition** — three-place tail acquisition; single-writer-per-run rule; DR rejoin discipline.
- **§10.26 Reference verifier distribution** — three-name CC8.1 citation (implementation, version, verification key).

Forty-plus normative sections from the chain-of-custody spec composed alongside Helmstad's existing FDA Quality System Regulation framework, ICH GCP / GMP processes, 21 CFR Part 11 electronic-records-and-signatures controls, and the institution's HITRUST-equivalent attestation regime. Where the chain reached, the integrity foundation came from the spec; where the chain did not reach, the §10.19 chain-coverage map was the boundary's normative shape and the regulatory framework on the other side carried the evidence load.

The §10.13 evidentiary-artifacts retention package would carry the engagement's trace through any future litigation or regulator-supervised dispute. Spec §10.13 names the documentary evidence the institution retains for the chain-data retention period or longer if a litigation hold extends it: SDK version manifest (the build identifier of the SDK in production during the period), SDK source-code hash and SLSA build-reproducibility attestation, HSM configuration (model, FIPS level, signing-key rotation history, separation-of-duties roster), daily seal-job logs (success/failure for each seal, timestamps, the HSM-signed `signed_at` value), change-management records covering any configuration change to the SDK / ledger / HSM during the period, and verifier output for the period showing PASS for each tenant-day. These artifacts substantiate FRE 901(b)(9) authentication of the process — the chain's verifier output proves the result; these artifacts prove the process. The institution's IT witness lays foundation from them at deposition without re-engineering the system. For Helmstad's BIMO inspection, the §10.13 package was equivalent: the inspector would read the seven artifacts of the evidence pack against the §10.13 retention list, and the IT lead's deposition-equivalent (the inspector's documented inquiry) would draw foundation from the same artifacts.

The §1.1 Daubert framing applied transitively. Each of the seven artifacts answered one or more of the four factors. Testability: artifact 2 (the verifier output) and artifact 6 (the daily seal); both reproducible from the §7 procedure plus the published public keys plus the test-vector corpus. Peer review: the spec itself plus the Apache-2.0 reference implementation plus the SOC 2 report; all three composed as the peer-review evidence. Known error rate: artifact 6 carried the §1.3 EUF-CMA / second-preimage / EUF-CMA composition by reference. General acceptance: every primitive in the chain — HMAC-SHA-256, RFC 6962 Merkle, Ed25519 — was NIST-standardized, and the combination matched deployed audit-log systems including Certificate Transparency. The FDA inspector did not need to know these facts to reach a PASS finding; the §1.1 grounding was informative, not normative. But if the inspector ever asked, every fact had a citation in the shipped artifacts.

---
