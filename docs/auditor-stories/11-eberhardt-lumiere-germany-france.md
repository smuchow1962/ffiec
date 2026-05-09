# 🧾 Diary of an Audit Day — Eberhardt Werkstoffe × Lumière AI

**Engagement:** Joint pre-audit ahead of EU AI Act enforcement and the BMW joint-supplier audit
**Clients (joint):**
- **Eberhardt Werkstoffe GmbH** — German automotive-electronics Mittelstand, ~€1.4B revenue, ~3,800 employees, Stuttgart HQ, Sindelfingen primary plant, third-generation family-owned. Sensors, ECUs, AI-driven predictive-maintenance modules for OEM customers (BMW, VW, Mercedes-Benz, Stellantis, Renault). Recent push into EV-specific modules (battery health, motor-controller predictive maintenance).
- **Lumière AI Sàrl** — French AI consultancy, Paris HQ in the 1er arrondissement, ~80 employees, ex-INRIA founders (Institut national de recherche en informatique et en automatique), strong fairness-and-explainability practice. Builds custom ML models for industrial customers including Eberhardt.
**Posture:** TesseraSeal in production at both companies. **The chain extends across the partnership boundary.**
- Eberhardt: 8 months on the OEM-facing inference path (every predictive-maintenance alert is sealed in Eberhardt's chain).
- Lumière: 4 months on the model-development pipeline (training data hashes, model artifact hashes, fairness-audit reports — sealed in Lumière's chain).
- The handoff: when Lumière delivers a new model, Lumière's model-artifact-hash + model-card-hash + fairness-audit-report-hash are recorded into an Eberhardt chain entry of `service.name = "model-handover"` as cross-vendor anchors — the schema spec §10.21 normates as `audit.model_handover.*`. Two chains compose at the handover event.
**Date:** Tuesday, the week after the engagement-before-this-one
**Audit team lead:** Dawn
**Stuttgart liaison:** Klaus Eberhardt, third-generation owner-CEO. 64. Engineering background — Maschinenbau from RWTH Aachen.
**Paris liaison:** Hélène Lefebvre, COO. 47. Ex-INRIA researcher, PhD in computational learning theory.
**Cities:** Stuttgart and Paris, both CET. About one hour apart by TGV-Eurostar. The team works in parallel and joins by video bridge midday and at debrief.

> **Reading note for this story.** The Eberhardt-Lumière engagement is the case that drove the spec's §10.20 (training-data retention vs deployment-window discipline) and §10.21 (cross-vendor model-handover schema) amendments — see §12 change log, Wave-6 fourth errata. Three findings surfaced in the field — a 90-day retention vs 9-18-month deployment-window asymmetry, an implicit cross-border-transfer attribute, and a singular `audit_report_language` field where multilingual reports would benefit from a plural array — are now closed by spec text. The story keeps the moments where the team caught them in the field, then reads them through the post-amendment spec so the reader sees the closure mechanism without losing the field history.

---

## Context

Eberhardt Werkstoffe builds the boxes the OEMs hide behind the dashboard. Sensors, ECUs (electronic control units), and the AI-driven predictive-maintenance modules that look at battery cells and motor controllers and decide whether the car needs to send the driver a warning or schedule a service appointment. Family business, third generation. Klaus Eberhardt's grandfather started it after the war making relays for Daimler-Benz. Klaus runs it now from a glass-fronted office in a Stuttgart industrial park that overlooks the Sindelfingen plant in the distance. The family still owns the company outright. The Mittelstand model — the German mid-size family-business sector — is a culture as much as an economic category. Klaus thinks of the chain the way his father thought of the relay-test bench: a thing you build once, build properly, and trust for thirty years.

Lumière AI is a different shape entirely. Eighty people in a renovated Haussmann building near the Louvre, founded six years ago by two ex-INRIA researchers who decided industrial ML had a fairness problem nobody was solving in production. They publish papers. They speak at NeurIPS. They also ship models that go into cars. Their customer list is short and exclusive — Eberhardt is the largest. Hélène Lefebvre runs the operations side; the founders run research. The two cultures meet at the model-handover boundary, which is where the chain meets the chain.

Eight months ago Eberhardt deployed TesseraSeal across the OEM-facing inference path. Every predictive-maintenance alert that goes from an Eberhardt module out to a BMW or VW vehicle is sealed in Eberhardt's chain. The hash of the model artifact that produced the alert is in the chain entry. The hash of the model card is in the chain entry. The vehicle ID is hashed into the chain entry. The OEM customer's deployment-side data is what BMW or VW logs on their side. Three chains, conceptually — Eberhardt's chain, Lumière's chain, the OEM's vehicle-side data — and the cross-vendor anchor (spec §10.21) is the bridge between the first two.

Four months ago Lumière deployed TesseraSeal on the model-development pipeline. Training data hashes, model artifact hashes, fairness-audit reports — all sealed in Lumière's chain under their own IKM (per §3 / §4.1, each tenant has its own per-tenant HKDF-bound session key). The deployment was driven by a 2024 incident. A battery-health prediction model regressed silently after a training-data update. Lumière delivered the regressed model to Eberhardt. Eberhardt deployed it to BMW. Two weeks of degraded predictions before someone noticed. Lumière did not have the chain at the time. They could not retrace which deployments had received the regressed model. Hélène had been the one who put the case together for the founders that the chain was no longer optional.

The handoff is the load-bearing thing. When Lumière delivers a new model to Eberhardt, an Eberhardt chain entry of `service.name = "model-handover"` records a JCS-canonical attribute set (per §5 RFC 8785 canonicalization). The shape the field team will recognize as the §10.21 `audit.model_handover.*` family:

```
{
  "audit.model_handover.provider": "lumiere-ai",
  "audit.model_handover.model_id": "battery-health",
  "audit.model_handover.model_version": "v4.2.1",
  "audit.model_handover.model_artifact_sha256": "abc123...",
  "audit.model_handover.model_card_sha256": "def456...",
  "audit.model_handover.fairness_audit_report_sha256": "ghi789...",
  "audit.model_handover.provider_chain_entry_id": "lum-2026-04-08-..."
}
```

The hash of Lumière's chain entry is the cross-vendor anchor. Verifying Eberhardt's chain confirms what Eberhardt received. Verifying Lumière's chain (with credentials Lumière grants on a request basis) confirms what Lumière handed over. Both must match for end-to-end verification. Spec §10.21 names this bidirectional cross-anchor verification by hash-equality between deployer and provider chain entries — independent verifier output on each side, plus the byte-equal hash join. That mechanism is the centre of gravity of this engagement.

Klaus and Hélène commissioned the joint audit because three things had converged. The EU AI Act enforcement deadlines for high-risk AI systems were approaching — predictive-maintenance for safety-critical automotive systems sits inside Annex III's high-risk category, and Lumière as the model provider had Article 16 obligations while Eberhardt as the deployer had Article 11 logging and Article 12 conformity-assessment obligations. (The chain composes onto Article 12 logging cleanly — §1.2 epistemic scope makes the proves/does-not-prove split explicit, which is exactly the language an EU AI Act conformity-assessment file needs.) BMW's vendor-management team had asked Eberhardt and Lumière for a joint-supplier audit that demonstrated end-to-end chain coverage at the OEM-supplier boundary. And the 2024 model-drift incident was in the rear view but not far enough that BMW had stopped asking about it.

The deliverable will be read by both companies' boards. By BMW's vendor-management team. By the German BSI (Bundesamt für Sicherheit in der Informationstechnik) for cybersecurity sign-off. And, for the cross-border GDPR aspect, by both the German LfDI Baden-Württemberg and the French CNIL. Five readers. Two countries. One report.

---

## Audit Team

The same eight-person team that walked Northbridge a quarter ago and Olmstead and the others. This week they split.

### Stuttgart (Eberhardt HQ)

- **Dawn** — Lead Auditor (governance + narrative)
- **Mike** — Application / API layer
- **Diana** — IAM and access control
- **Luis** — DevOps / logs / pipelines
- **Tom** — Internal-audit liaison specialist (visiting team — partners with the joint-engagement liaison from both companies)

### Paris (Lumière HQ)

- **Raj** — Database specialist
- **Elena** — CRM systems
- **Chen** — Data engineering / ETL

The split is deliberate. Eberhardt is the larger surface area and has the cross-vendor anchor on its side of the chain, so the bulk of the team takes Stuttgart. Lumière is denser per system and the model-development pipeline is the most concentrated piece of Paris work, so Raj, Elena, and Chen take Paris. The two sides are about an hour apart by train, both on CET, no time-zone offset. The team works in parallel and joins by video bridge midday and at debrief.

---

### 🌅 8:30 AM CET — Kickoff (Stuttgart)

Dawn had landed at Stuttgart Airport the night before and taken the S-Bahn into Mitte. The hotel was a small one near the Schlossplatz. The drive out to the Eberhardt office in the morning was twenty-five minutes through the Stuttgart hills. She had her coffee in the cup holder and the engagement brief on her tablet.

*Two countries*, she thought as she came up the autobahn ramp. *Two companies. One chain that crosses the boundary. EU AI Act for both. BMW watching. 2024 model-drift in the rear view. We test whether the cross-vendor anchor actually composes — what one chain says, the other confirms.*

She took stock of the prior engagements. Northbridge had been the gold standard — full single-tenant deployment, eighteen months mature, the team had spent four days trying to find a gap and found a stale comment in a YAML file. Ten engagements back now. One §10.16 non-conformance on the books, the chain otherwise byte-for-byte clean. Dawn had not seen its match in the ten that followed. Mercator had been the bifurcation — sepsis CDS on the chain, the EHR off the chain, Patricia Okonkwo's funding-roadmap framing. Stelvio had been the three-zone version — AI side sealed, OT mutable, IT business legacy, Maria Costanza's triage. Atrio had been the multi-tenant test — forty-seven tenants under twelve sponsor-bank IKMs, fourteen hundred verifier runs, zero failures, Naomi Reisinger's coordinated examiner room.

Then Helmstad. The biopharma. The CRO data feed where Quintessa had PGP-signed the SFTP delivery and Helmstad had recorded the SHA-256 at the boundary, and the source side beyond the boundary lived on Quintessa's SOC 2. Pacific Crescent — the utility, AI gas-pipeline leak detection on the chain, OT historian off the chain, the Brentwood alert that turned out to be a real small leak. Olmstead — the university, AI admissions screening on the chain, Slate free-text rationale-fields off the chain, two override-down decisions where the rationale was gone.

After Olmstead there had been two more. The Korean engagement at Sun-Won where the §4.4.1 cross-border-transfer attribute family had been the central finding — the chain had not carried an explicit `audit.cross_border_transfer.lawful_basis_type` and the PIPC examiner had pointed that out specifically. And Salt Pond Toys where the §10.19 chain-coverage map became the framing device for boundary discipline. Nine prior engagements before today. Today is the tenth-and-eleventh — joint, parallel, two countries.

By the time the team rolled into Stuttgart, TesseraSeal had been audited across the previous nine engagements — banking, healthcare, BaaS, industrial, biopharma, utility, higher-ed, Tel Aviv vendor, K-beauty retail across Korea + Taiwan, toy supply chain across RI + Shenzhen + LA. Cross-vendor automotive AI under EU AI Act Article 12 was a new composition — but the chain primitive was familiar. The team's posture coming into Eberhardt-Lumière was confidence in the primitive and confidence in many of the regulatory compositions it had already absorbed. The new questions were specific: training-data retention floor against a long deployment window, a cross-vendor model-handover schema with a plural-array `audit_report_languages` shape, and the within-EU cross-border attribute composition for a German-French Article 28 channel.

*It never is*, Dawn thought. *Today the question is whether the joint chain holds at the handoff. Each side's chain looks fine alone. The question is the seam.*

The Eberhardt building was in a Stuttgart industrial park, glass front, the Eberhardt name in plain stainless steel on the wall. Klaus Eberhardt was waiting in the lobby in a navy jumper and grey trousers. He was 64, lean, with engineer's hands and a directness Dawn recognized in the first thirty seconds.

"Dawn. Welcome to Stuttgart." His English was almost unaccented. "You came in last night?"

"I did. Hotel near the Schlossplatz."

"Good. The airport hotel is a tragedy." He shook hands with each of the Stuttgart team in turn — Mike, Diana, Luis, Tom — as they came through the revolving door from the parking lot.

The conference room was on the second floor with a view of a lawn that ran down to a pond. There was a long table with a Polycom video unit at the head and a screen on the far wall. Hélène Lefebvre was already on the screen from Paris, in a small office with a window behind her that showed grey Parisian morning. The Paris team — Raj, Elena, Chen — were visible in a corner of the screen, just arrived themselves.

Klaus did not waste time.

"Good morning. Let me say what I want from this. Eberhardt has had the chain for eight months on the OEM-facing inference path. I am reasonably confident in it. Lumière deployed four months ago on the model-development side. Hélène is reasonably confident in hers. Neither of us has tested the seam in front of an audit team. The EU AI Act enforcement deadline is in seven months. BMW's vendor-management team has asked us for a joint-supplier audit that demonstrates the seam works. We have a 2024 incident in our recent history where Lumière delivered a regressed battery-health model and we both saw it late. I want today's audit to test the seam. I want to know — if a BMW customer's car has a false-positive predictive-maintenance alert tomorrow — whether we can together identify the root cause: was it Lumière's model, or our integration, or BMW's vehicle-side data."

Hélène, on the screen: "I will say the same thing from Paris. Lumière has obligations as a provider under Article 16 of the EU AI Act — fairness audit, transparency to deployers. The chain entries support those obligations on our side. Klaus's chain on the deployer side covers Article 11 logging and Article 12 conformity assessment. The cross-vendor anchor is what makes the two halves a whole. I want it tested today."

Dawn put her coffee down on the table.

"Thank you both. That is exactly the right framing." She looked around the Stuttgart side of the table, then at the screen for the Paris side. "Stuttgart morning is the inference side and the Eberhardt IAM. Paris morning is the model-development chain and the Lumière database. We meet by video bridge at noon to walk the cross-vendor anchor live. Afternoon is the API layer in Stuttgart, the training-data pipeline in Paris. Three o'clock is the joint reconciliation test by video bridge — pick a model deployment, trace it end to end through both chains. Five-thirty is the joint debrief."

Klaus nodded. "The engineering team at Eberhardt is expecting Mike at nine. Our predictive-maintenance lead is Maximilian Brenner — twelve years here, started in firmware, moved into AI integration four years ago. He will not waste your time."

Hélène, on the screen: "Lumière's model-development lead is Aurélien Marchand. Six years at Lumière, started as a research engineer, runs the production pipeline now. He is expecting Raj at nine."

Tom — the internal-audit liaison — had been quiet. He looked up.

"I had a call yesterday with the BMW vendor-management lead — a person named Stefan Kuhn. He confirmed the deliverable will be read by his team alongside the BSI sign-off. He had three questions he asked us to address explicitly. One — does the cross-vendor anchor compose end to end without trust in either party's claim. Two — can a false-positive predictive-maintenance alert be root-caused across the seam. Three — what is the forensic gap if the 2024 incident happened today instead of two years ago."

Klaus and Hélène exchanged a look across the video bridge.

"Those are the right questions," Klaus said. "That is what we will answer."

---

### 🧩 9:15 AM CET — Mike on the Predictive-Maintenance Service (Stuttgart)

Maximilian Brenner had a laptop open and one terminal window. No deck. He was a person who clearly preferred a terminal to a slide.

"You want to see a sealed predictive-maintenance alert," he said. "Pick a date. Pick an OEM."

Mike picked April 8. BMW iX. The high-volume EV.

Maximilian typed for ten seconds. "April 8, BMW iX battery-health module. The classifier ran 14,287 inferences across the BMW iX fleet that day. 47 came back as `service-recommended`. 12 came back as `urgent-service-required`. Pick one."

"Pick an `urgent-service-required` with the highest confidence."

He typed again. The terminal showed a JSON entry — structured, with fields Mike recognized as the §4.4 `ffiec.chain.*` and §4.4.2 `audit.deployment.*` shape composed with OTel GenAI semconv:

```
{
  "ffiec.chain.spec": "v1.0",
  "ffiec.chain.format_version": "v1",
  "ffiec.chain.chain_kind": "model_call",
  "ffiec.chain.run_id": "eb-bh-2026-04-08-bmwix",
  "ffiec.chain.seq": 4827193,
  "ffiec.chain.tenant_id": "eberhardt-battery-health",
  "ffiec.chain.captured_at": "2026-04-08T11:42:17.083Z",
  "ffiec.chain.key_version": 6,
  "service.name": "pred-maint-inference",
  "gen_ai.request.model": "battery-health",
  "gen_ai.response.model": "battery-health-v4.2.1",
  "audit.deployment.intent": "production",
  "audit.deployment.policy_version": "eb-deploy-2026-q2",
  "audit.model_handover.model_artifact_sha256": "sha256:7f2a91...c4d3",
  "audit.model_handover.model_card_sha256": "sha256:8e1b82...a4f1",
  "vehicle_vin_hash": "sha256:a91f8b...ee27",
  "module_serial_hash": "sha256:c12e4f...8a91",
  "input_telemetry_hash": "sha256:d83a01...b3ef",
  "classification": "urgent-service-required",
  "confidence": 0.94,
  "downstream_oem": "bmw-ag",
  "model_handover_entry_ref": "eb-mh-2026-04-01-lum-00012",
  "ffiec.chain.prev_hash": "...",
  "ffiec.chain.payload_hash": "...",
  "merkle_path": [...],
  "seal_ref": "eberhardt-battery-health-2026-04-08-eod"
}
```

Mike noted the `audit.deployment.intent = "production"` (per §4.4.2 — emission required because Eberhardt operates a canary on a parallel pipeline) and the `audit.deployment.policy_version` paired with it (the §4.4.2 conditional requirement: any `audit.deployment.*` attribute requires the policy version).

"Run the verifier on it," Mike said.

Maximilian typed:

```
herald-verify --tenant=eberhardt-battery-health \
  --service=pred-maint-inference --date=2026-04-08 \
  --entry-id=eb-bh-2026-04-08-bmwix-00874
```

The terminal hesitated for four seconds and printed:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key eberhardt-battery-health-2026-q2
```

The shape Mike was looking for. Spec §7 names the verifier output format normatively (`Status:` / `Step:` / `Reason:`), and §10.12 names the exit-code contract (0 on PASS). Mike leaned back.

"Pick one from a hundred and eighty days ago."

Maximilian picked an entry from October 12, 2025. Same fleet, same model family — though running on `battery-health-v3.9.0`, the predecessor. Ran the verifier. Same four seconds. Same PASS. Same twelve steps. Twelve steps because §7's procedure is twelve steps in order — file-header pre-flight, fingerprint check, MAC recompute, Merkle resolution, signature verification — and the verifier walks them all on a 180-day-old entry the same way it walks them on an entry from this morning.

> **✓ Confirmation #1 — Eberhardt chain integrity holds at eight months and 180 days under the §7 twelve-step procedure**
> The predictive-maintenance service has been emitting sealed entries for 244 days. The verifier resolves a recent entry in 4 seconds and a 180-day-old entry in 4 seconds. Twelve verification steps per §7 including HMAC recomputation (§4.1), Merkle path resolution against the daily-seal Merkle root (§4.2), and Ed25519 signature verification against the published quarterly public key (§4.3). The chain endures across the seal boundary and across model-version transitions. The §1.3 security definitions hold across both sides — EUF-CMA for HMAC-SHA-256 per FIPS 198-1, second-preimage resistance for the RFC 6962 Merkle construction, EUF-CMA for Ed25519 per FIPS 186-5. The §4.1.1 session-key handshake is per-tenant deterministic; constant-time comparison per §10.8 is enforced at the verifier. IKM length conforms to §10.6 (32 bytes minimum) with §10.6.1 generation requirements (cryptographic-strength RNG, FIPS 140-3 attestation available on request). Output format conforms to §7 normative form; exit code 0 per §10.12. Per §10.26, Eberhardt's CC8.1 names the reference verifier implementation, the version, and the verification key it uses to authenticate the verifier binary at run time.

Mike asked: "What is sealing the seal?"

"Daily Ed25519 signature on a Thales Luna network HSM at the Sindelfingen plant. On-prem, FIPS 140-2 Level 3, satisfying §10.5 HSM custody. We are Mittelstand — for things this load-bearing we prefer on-prem. The IQ/OQ documentation for the HSM configuration is in our QMS. BSI IT-Grundschutz baseline — that is the German Federal Office for Information Security's compliance baseline. We are ISO 27001 certified, ISO/IEC 27017 certified for the cloud-touching parts, and TISAX certified for the auto-supply community — that is the Trusted Information Security Assessment Exchange the German automotive industry uses for coordinated supplier security assessments. The chain artifacts feed the TISAX assessment. Trusted-time integration is currently NTP-synchronized per §10.4 with §10.14 RFC 3161 trusted-timestamp tokens RECOMMENDED but not yet adopted — we may turn that on for the BMW-facing inference path next year for higher-stakes timestamp credibility."

"Show me the IQ/OQ."

Maximilian pulled it up. 62 pages. Configuration captures, key-attribute definitions, partition isolation, key-ceremony minutes from August 2025 when the daily seal had been first stood up. Klaus's signature on page 62.

Mike noted that the partition-creation ceremony fell squarely under §10.17 (HSM partition ceremony attestation) — `chain.partition_ceremony_attended` for `ceremony_type = "partition_created"` with the IQ/OQ minutes carrying the §10.17 signatory schema (role, name, and the Round-17 M&A-P1 `entity_affiliation` field that discriminates Eberhardt-authorized signers from any vendor-on-site witness).

"Show me the chain entry for that ceremony."

Maximilian pulled it up. `chain.partition_ceremony_attended` with `partition_handle`, `ceremony_started_at_utc`, `ceremony_completed_at_utc`, the `signatories` array (Klaus, the platform-engineering lead, and the SecOps lead, all with `entity_affiliation = "eberhardt-werkstoffe"`), the `witness` (Bettina Hofer, separate party from the signatories per the §10.17 separation requirement), and `attendance_pdf_sha256` binding the scanned IQ/OQ document.

> **✓ Confirmation #2 — IQ/OQ for the on-prem HSM signing infrastructure is documented to BSI IT-Grundschutz and ISO 27001 standards AND chain-coupled per §10.17**
> The Thales Luna network HSM at Sindelfingen has a 62-page IQ/OQ in the Eberhardt QMS. Key ceremony minutes are recorded with two-of-three approval (Klaus Eberhardt, the platform engineering lead, and the SecOps lead). The configuration is reproducible from the IQ/OQ. The artifacts feed into Eberhardt's TISAX assessment. ISO 26262 functional-safety auditors accept the chain artifacts as part of the audit trail for AI-derived safety-related outputs. The §10.17 `chain.partition_ceremony_attended` event was emitted for partition creation, IKM rotation, and the controlling-person rotation that happened when Eberhardt added Bettina to SecOps in November — three events, all with witness signatures per §10.17, all with `attendance_pdf_sha256` bound. The `entity_affiliation` field is recorded on every signatory per the Round-17 M&A-P1 fix.

Mike asked: "Show me the model-handover binding."

Maximilian pulled up the entry referenced by `model_handover_entry_ref` in the inference entry — `eb-mh-2026-04-01-lum-00012`.

```
{
  "ffiec.chain.chain_kind": "audit",
  "ffiec.chain.run_id": "eb-mh-2026-04-01-lum",
  "ffiec.chain.seq": 4801237,
  "service.name": "model-handover",
  "ffiec.chain.tenant_id": "eberhardt-battery-health",
  "ffiec.chain.captured_at": "2026-04-01T09:15:42.001Z",
  "audit.model_handover.provider": "lumiere-ai",
  "audit.model_handover.model_id": "battery-health",
  "audit.model_handover.model_version": "v4.2.1",
  "audit.model_handover.model_artifact_sha256": "sha256:7f2a91...c4d3",
  "audit.model_handover.model_card_sha256": "sha256:8e1b82...a4f1",
  "audit.model_handover.fairness_audit_report_sha256": "sha256:9d3c72...e5f2",
  "audit.model_handover.audit_report_languages": ["fr", "en"],
  "audit.model_handover.provider_chain_entry_id": "lum-mb-2026-04-01-bh-00031",
  "audit.model_handover.training_data_retention_floor_days": 720,
  "audit.model_handover.contract_id": "EB-LUM-MSA-2024-001",
  "audit.model_handover.contract_version": "v3.2",
  "audit.model_handover.contract_hash_sha256": "sha256:5e8c14...a7d2",
  "received_by": "max.brenner@eberhardt.de",
  "approver": "klaus.eberhardt@eberhardt.de",
  "seal_ref": "eberhardt-battery-health-2026-04-01-eod"
}
```

"That is the §10.21 schema landing cleanly," Mike said. "Provider, model_id, model_version, the three artifact hashes, the cross-anchor entry ID, the retention floor, and the contract triple." (The §10.21 contract triple — `contract_id` + `contract_version` + `contract_hash_sha256` — closes Round-17 M&A-G2: post-close acquirer reads the contract version that governed each delivery from the chain alone, no dependence on a seller's contract binder.)

"Each model-handover entry references the Lumière chain entry ID," Maximilian said. "When we receive a model from Lumière, the engineering lead approves the artifact bundle, the chain entry is written, and the inference service starts using the new model only after the chain entry is sealed. The inference entry's `model_handover_entry_ref` field points back to that handover entry. So every inference is bound to the model-handover record that admitted the model."

> **✓ Confirmation #3 — Model artifact, model card, and fairness-audit-report hashes are bound to every Eberhardt inference via the model-handover entry per §10.21**
> Each inference's `audit.model_handover.model_artifact_sha256` and `audit.model_handover.model_card_sha256` match the values in the referenced model-handover entry. Recomputing the SHA-256 of the production model artifact at any time matches the sealed value. The same is true for the model card and the fairness-audit report. The model-handover entry itself references the corresponding Lumière chain entry by `audit.model_handover.provider_chain_entry_id` — the cross-vendor anchor §10.21 normates. Bidirectional verification by hash-equality between deployer and provider chain entries is exactly the §10.21 cross-anchor mechanism. The handover entry also carries `audit.model_handover.audit_report_languages = ["fr", "en"]` (plural array per §10.21 — the spec elevates the plural form because multilingual reports are common when models cross jurisdictions, and BMW's vendor-management read is German-English-French) and `audit.model_handover.training_data_retention_floor_days = 720` (the 24-month commitment we'll come to in the lunch debrief — §10.20 retention floor).

"And if a new model comes in from Lumière," Mike said, "the inference entries from that day forward reference the new model-handover entry."

"Yes."

"And the old model — if it stops being used — its inference history is still verifiable."

"Yes. Each entry references the model-handover that was in force at the time. Past entries do not change when a new model arrives." (Per §10.3 append-only enforcement plus §4.4 genesis-block uniqueness — past chain entries are not retroactively mutable; the chain captures what was in force at capture time.)

Mike wrote that down. *Hash-bound provenance from inference back to model-handover, with cross-vendor anchor to Lumière. The §10.21 schema plus the §4.4 chain envelope plus the §4.4.2 deployment-intent attribute form a single integrity-bound record per inference.*

He had one more question. "Show me the canary."

Maximilian pulled up a parallel pipeline. Eberhardt operates a canary deployment in addition to production for any new model version — three-percent of the BMW iX fleet receives the candidate model for two weeks before promotion to full production. The canary entries carry `audit.deployment.intent = "canary"` per §4.4.2 with the conditional `audit.deployment.canary_traffic_pct` (currently `3.0` for the v4.3 candidate that is in canary now), `audit.deployment.experiment_id`, and `audit.deployment.policy_version` per the §4.4.2 conditional schema. The MRM disposition per §4.4.2 is "bounded production-validation activity" — Eberhardt's MRM committee reviews the canary's traffic-percentage trajectory, the canary's decision-equivalence record against production, and the rollout/rollback decisions.

> **✓ Confirmation #3a — Canary deployment is `audit.deployment.intent = "canary"` with required fields per §4.4.2**
> The canary pipeline emits `audit.deployment.intent = "canary"`, `audit.deployment.canary_traffic_pct = 3.0`, `audit.deployment.experiment_id = "exp-bh-v4.3-canary-2026-q2"`, and `audit.deployment.policy_version`. The §4.4.2 MRM disposition for `canary` is operating per the spec — Eberhardt's MRM committee reviews the canary trajectory monthly. A canary entry without `canary_traffic_pct` would be a control-completeness gap per §4.4.2; Mike spot-checked twenty canary entries from the prior two weeks and every entry carried the field.

"And the production pipeline — what intent does it carry?"

"`production`. The §4.4.2 emission requirement applies because we operate a canary, so we set `intent` on every entry. If we were single-version single-region with no canary or A/B, the §4.4.2 attribute set would be optional and we could omit it. We are not, so we emit it on every inference."

Mike noted the discipline. *Per §4.4.2, institutions running A/B / canary / multi-region drift / vendor-reroute / regulatory-sandbox / DI-test-run MUST emit `audit.deployment.intent`. Eberhardt operates a canary, so `intent` is REQUIRED on every inference, not optional. The chain carries the policy version that classified the invocation.*

He moved to one final check — the gen_ai response-model identifier. The §4.4 SDK-side enforcement rule is normative: SDKs MUST refuse to emit a chain entry whose attribute set includes any `gen_ai.*` namespace prefix attribute AND lacks either `gen_ai.request.model` or `gen_ai.response.model` non-empty. Maximilian demonstrated by attempting to write a malformed entry through the SDK harness. The SDK refused at write time with `GenAIModelIdentifierMissing` raised before MAC compute, exactly as §4.4 SDK-side enforcement specifies. The §7 step 12a check is defense-in-depth at the verifier; the SDK refusal closes the source.

> **✓ Confirmation #3b — SDK enforces `gen_ai.{request,response}.model` per §4.4 SDK-side enforcement rule**
> The SDK refuses at write time when `gen_ai.*` attributes are present without a non-empty `gen_ai.request.model` or `gen_ai.response.model` — `GenAIModelIdentifierMissing` raised before MAC compute. The §7 step 12a verifier check remains as defense-in-depth.

---

### 🧠 10:00 AM CET — Raj on the Lumière Chain (Paris)

Aurélien Marchand had set up a workstation for Raj in a small glass-walled office near the Lumière engineering bullpen. The bullpen was loud with quiet typing and one whiteboard that had a fairness-bias-bound proof half-erased on it. The view out the window was a courtyard with a single chestnut tree.

"You want to see a sealed model build," Aurélien said. He had switched to English without waiting for Raj to do so. "Pick a model. Pick a build date."

Raj picked the battery-health model — the same family that Stuttgart was looking at — and asked for the build that had produced `v4.2.1`.

Aurélien typed for ten seconds. "Battery-health v4.2.1. Build date April 1, 2026. Pre-handover to Eberhardt on the same day. The build pipeline emitted seventeen chain entries from training-start through artifact-finalization through fairness-audit-finalization through handover-prep. Pick one."

"Pick the artifact-finalization entry — the one that anchors the model file Eberhardt receives."

He pulled it up.

```
{
  "ffiec.chain.chain_kind": "audit",
  "ffiec.chain.run_id": "lum-mb-2026-04-01-bh",
  "ffiec.chain.seq": 1872391,
  "service.name": "model-build",
  "ffiec.chain.tenant_id": "lumiere-model-dev",
  "ffiec.chain.captured_at": "2026-04-01T08:42:11.502Z",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "model_card_sha256": "sha256:8e1b82...a4f1",
  "training_data_manifest_sha256": "sha256:b73c92...11f4",
  "training_run_id": "tr-2026-03-28-bh-001",
  "fairness_audit_report_sha256": "sha256:9d3c72...e5f2",
  "fairness_audit_entry_ref": "lum-fa-2026-03-31-bh-00007",
  "handover_target": "eberhardt-werkstoffe",
  "approver": "aurelien.marchand@lumiere.ai",
  "seal_ref": "lumiere-model-dev-2026-04-01-eod"
}
```

The model-artifact SHA-256 matched the value Stuttgart had pulled up an hour earlier — bit-for-bit. The model-card SHA-256 matched. The fairness-audit-report SHA-256 matched. (The hashes match by construction: §5 RFC 8785 JCS canonicalization plus SHA-256 yield byte-identical values for byte-identical inputs; the equality is the load-bearing property the §10.21 cross-anchor depends on.)

"Run the verifier."

Aurélien typed:

```
herald-verify --tenant=lumiere-model-dev --service=model-build \
  --date=2026-04-01 --entry-id=lum-mb-2026-04-01-bh-00031
```

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key lumiere-model-dev-2026-q2
```

Raj nodded.

> **✓ Confirmation #4 — Lumière chain integrity holds at four months on the model-build service**
> The model-development pipeline has been emitting sealed entries for 124 days. The verifier resolves the build entry that anchors the model Eberhardt is currently running in 4 seconds. Twelve verification steps per §7. The same shape as Eberhardt's chain, on a separate IKM (per-tenant HKDF binding per §4.1) under a separate HSM in a separate jurisdiction. Independent integrity statements that compose at the model-handover boundary per §10.21 — the §1.4 compositional security argument extended across an organizational boundary.

"What is sealing your seal?"

"OVHcloud Paris-region HSM. We chose OVH for two reasons — French sovereignty over the key material and GDPR data-residency for the training data. The HSM configuration follows ANSSI recommendations — that is the French national cybersecurity agency, Agence nationale de la sécurité des systèmes d'information. We have an ANSSI-aligned configuration document and key-ceremony minutes from December 2025."

"Show me."

Aurélien pulled up the configuration document. 41 pages. Hélène's signature on the last page. The key ceremony had been four people — the two founders, Hélène, and the SecOps lead — with a two-of-four threshold for key-recovery operations. The §10.17 `chain.partition_ceremony_attended` event for that December ceremony carried `signatories` for all four with `entity_affiliation = "lumiere-ai"`, plus a colocation engineer as `witness`, plus the `attendance_pdf_sha256` for the scanned ceremony document. The §10.17 cross-language CC8.1 discoverability rule applied here — the runbook is in French, so Lumière's CC8.1 (in English for cross-jurisdiction readers) cross-references the French runbook by title and named ceremony-procedure sections.

> **✓ Confirmation #5 — Lumière HSM configuration is ANSSI-aligned, documented, and chain-coupled per §10.17**
> OVHcloud Paris-region HSM with an ANSSI-recommended configuration. 41-page configuration document signed by Hélène Lefebvre. Key-ceremony minutes from December 2025 with two-of-four approval. French sovereignty and GDPR data-residency are the operational drivers. The §10.17 `chain.partition_ceremony_attended` events for partition creation, IKM rotation, and the founders' ceremony all carry the §10.17 schema with `entity_affiliation` per Round-17 M&A-P1. The cross-language CC8.1 discoverability rule (§10.17 closing paragraph plus §10.18 runbook cross-referencing) is satisfied — Lumière's English CC8.1 cross-references the French runbook by title and named sections so the German BSI auditor and BMW's German-English vendor-management readers can locate the procedures.

Raj asked the next question.

"Walk me through the fairness-audit entry."

Aurélien pulled up `lum-fa-2026-03-31-bh-00007`.

```
{
  "ffiec.chain.chain_kind": "audit",
  "ffiec.chain.run_id": "lum-fa-2026-03-31-bh",
  "ffiec.chain.seq": 287413,
  "service.name": "fairness-audit",
  "ffiec.chain.captured_at": "2026-03-31T16:12:48.119Z",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "audit_protocol_version": "lum-fa-protocol-v3.1",
  "audit_report_sha256": "sha256:9d3c72...e5f2",
  "audit_report_languages": ["fr"],
  "auditor_subject": "claire.dubois@lumiere.ai",
  "approver": "helene.lefebvre@lumiere.ai",
  "fairness_metrics_summary": {
    "disparate_impact_ratio_age_cohort": 0.97,
    "calibration_error_max": 0.013
  }
}
```

Raj read the entry twice. The `audit_report_languages` field caught his eye.

"Plural array. Currently `["fr"]`. The audit report is in French only."

"Yes. The auditor wrote it in French. Lumière is a French company. The model card is bilingual — French and English — but the full fairness-audit report is in French only at v4.2.1."

"And Eberhardt's deployment documentation is in German."

"Yes. The model-handover documentation that Eberhardt produces is in German — and bilingual where we cooperate, but the regulatory file at Eberhardt's side is German-language."

Raj wrote that down. *BMW examiner reads what — German, French, English? The chain captures the language array. The reader needs to know which translations are available.* Spec §10.21 makes the `audit.model_handover.audit_report_languages` field plural-array specifically so a downstream reader without the primary language can discover available translations through the chain itself; the plural-array discipline closes what would otherwise be a singular-field gap. (Round-17 M&A-N2 surfaced this point and the spec answered with the plural form.)

> **✓ Confirmation #6 — Fairness-audit report hash is anchored in the chain with language-array metadata captured per the §10.21 plural-array discipline**
> Every fairness-audit chain entry records the SHA-256 of the audit report, the auditor's identity, the audit-protocol version, the approver, and the report languages array. The chain captures `audit_report_languages` as a plural array — the same shape §10.21 normates for `audit.model_handover.audit_report_languages` (Round-17 M&A-N2). At v4.2.1 the array is single-element `["fr"]`; the post-engagement commitment Hélène will make at lunch is to populate it with English and German translations, and the array shape lets the chain carry the additional translations without a schema change. The procedural question — whether examiners reading the joint deliverable have access to the report in the language they read — is what the array discovery enables; the chain itself surfaces the translation gap to the reviewer.

Raj closed the fairness-audit window and opened the training-data manifest.

"Show me the training-data binding for v4.2.1."

Aurélien pulled it up. The manifest hash referenced a manifest document that listed every training-data shard, its source, its de-identification status, its consent basis under GDPR, and the SHA-256 of the shard. The manifest itself was the canonicalization of all of that. The manifest hash was sealed in the build entry.

"How long do you retain the training data shards themselves?"

Aurélien paused for a beat. "Ninety days post-model-delivery."

"Ninety days."

"Yes. After that the shards are deleted. The manifest hash and the manifest document are retained — those are part of the chain — but the underlying shards rotate out at ninety days post-delivery."

Raj wrote that down very carefully. *Ninety days post-model-delivery. The model is deployed at Eberhardt for as long as Eberhardt runs it. If a regression appears six months later and someone asks to retrace which training shards produced it — the shards are gone.*

"Hélène was the one who set the policy?"

"Yes. The policy is GDPR data-minimization-aligned. The argument was that we cannot retain training data indefinitely under GDPR Article 5(1)(c) data-minimization principle, and ninety days post-delivery was the operational compromise."

"How long does Eberhardt typically run a model in production?"

Aurélien thought. "Battery-health v3.9.0 was in production for fourteen months before v4.2.1 replaced it. v4.2.1 has been live for thirty-eight days. The pattern is — depending on the OEM rollout cadence — between nine and eighteen months per major model version."

"So the deployment window is typically four to six times longer than the training-data retention window."

"Yes."

This is the moment that drove the spec. At the time of the field engagement the spec did not yet name a retention floor; this finding plus the lunch conversation that followed are what the spec's §10.20 amendment cites as its worked example. Reading the finding back through the post-amendment spec:

> **⚠️ Finding #1 — Lumière's 90-day training-data retention violates §10.20 (training-data retention vs deployment-window discipline) — Partial conformance, remediation in flight**
> Spec §10.20 (normative; landed in v1.0a as the Wave-6 Eberhardt-Lumière fourth errata) requires that training-data shard retention be at least as long as the longest active deployment window of any model trained on the data, plus an investigation buffer (typically 60-90 days). Lumière retains training-data shards for 90 days post-model-delivery; Eberhardt typically runs a deployed model in production for 9-18 months. The asymmetry is the §10.20 worked example verbatim — the chain detects regression, but the regression's training-data root cause cannot be retrieved beyond the retention boundary. Per §10.20: GDPR Article 5(1)(c) data-minimization tension is resolved through the institution's Article 6(1)(f) legitimate-interests determination tied to EU AI Act Article 12 logging obligations; the retention is purpose-limited and the data is held under appropriate access controls. Per §10.20 bidirectional cross-vendor anchor implication: the provider-side retention floor governs the deployer-side forensic depth — Lumière's 90-day floor bounds Eberhardt's chain forensic reach into the training-data root-cause path. Severity: **Partial conformance**, not Nit. Hélène commits at lunch to the §10.20-compliant 24-month retention floor (max deployment window 22 months + 60-day investigation buffer rounded up to 24); the §10.21 `audit.model_handover.training_data_retention_floor_days` field on the next handover entry will carry `720` so the deployer's CC8.1 evidences the floor commitment without dependence on the underlying contract document at audit time.

Raj had not raised his voice. He had not even changed his expression. He simply wrote it down and went on.

"Show me the IAM for the model-build service."

---

### 🔐 11:00 AM CET — Diana on Eberhardt IAM (Stuttgart) and the Identity Walk in Paris

Diana started the IAM session in Stuttgart by sitting down with Eberhardt's identity-platform lead — a person named Andreas Vogt — and asking the question she always asked first.

"Show me a credential rotation for the predictive-maintenance service."

Andreas pulled up the chain entry for the most recent rotation — May 1, 2026, at 03:00 UTC. The service account `svc-pred-maint-inference` had its credentials rotated automatically. The rotation event emitted a sealed chain entry with the old credential fingerprint, the new credential fingerprint, the rotation reason (scheduled), the approver (Klaus Eberhardt, with a two-of-three approval from the platform engineering lead and the SecOps lead), the approval ticket reference, and the timestamp. The entry was a `chain_kind = "operational"` event per §3 / §10.2.

> **✓ Confirmation #7 — Eberhardt service-account IAM is chain-coupled end to end**
> Every credential lifecycle event — issuance, rotation, revocation, scope change — for the predictive-maintenance service produces a sealed chain entry under `chain_kind = "operational"` (§3 enumeration; §10.2 operational events list). Policy-level changes (who can rotate, who can approve, what the rotation interval is) require two-of-three approval and are themselves chain entries. Per §10.3 append-only enforcement, there is no path to change a service credential without producing a chain record. The §3 character class governs the tenant-id binding; Eberhardt's `tenant_id = "eberhardt-battery-health"` conforms natively without §3.1 legacy-aliasing migration.

"What does the engineer's identity look like? When Maximilian ran the verifier this morning — what was his SSO subject?"

"On-prem Active Directory federated to Azure AD. Maximilian's SSO subject is `max.brenner@eberhardt.de` — that is the Azure AD UPN. The federation is documented; the BSI IT-Grundschutz baseline says we have to be able to walk the federation chain back to the institutional credential, which is the on-prem AD account at Eberhardt. We can. There are no consumer identity providers in our trust chain."

"Twelve sites, twenty institutional IdPs?" Diana asked, half-joking, comparing notes against Helmstad's clinical-trial federation pattern.

Andreas smiled briefly. "Two. Eberhardt Stuttgart on-prem AD and Eberhardt Sindelfingen on-prem AD — same forest, two domain controllers. We do not federate to OEM customer IdPs because we do not give OEM customer staff access to the inference service directly. They consume the alerts on their side. Their identity is on their side."

Diana wrote that down. *Two-IdP shop. Easier than Helmstad's twelve. Easier than Olmstead's seventeen-school federation. Easier than Atrio's per-tenant SSO matrix. The Mittelstand culture shows in IAM — small footprint, well-controlled.*

> **✓ Confirmation #8 — Eberhardt's IAM trust chain is two on-prem AD domains, fully under chain governance**
> Eberhardt federates on-prem AD Stuttgart and on-prem AD Sindelfingen to Azure AD. Both are institutional. There are no consumer IdPs in the trust chain. Engineer SSO subjects in the chain entries resolve to the on-prem AD account in either domain. The BSI IT-Grundschutz baseline's identity-provenance requirement is satisfied without intermediate hops. The institution's CC8.1 names the federation per §10.18 cross-referencing rule (the Eberhardt runbook section "Identity Federation" cites §3 tenant-id and §10.1 IKM registration explicitly).

Meanwhile, in Paris, Raj had moved into the Lumière IAM walkthrough with Aurélien because Diana was on the Eberhardt side. Lumière's identity setup was Google Workspace as the primary IdP with a custom SSO layer in front of the model-build platform. Aurélien himself authenticated via Google Workspace. Hélène was the second factor on the model-build platform itself for any approval-bearing action — separation of duties between author and approver. Raj noted the §3.1 legacy-tenant-identifier handling — Lumière's `tenant_id = "lumiere-model-dev"` conformed natively to the §3 character class without aliasing migration. Aurélien confirmed the tenant-id was registered through Lumière's IAM provisioning at deployment time and SDK-side enforcement (per §3 — SDKs MUST reject a `tenant_id` that does not match the class at construct time) was active.

Each model-build chain entry recorded the author SSO subject and the approver SSO subject as separate fields. Raj cross-checked four entries; the author and approver were never the same person. The chain enforced that. A model-build entry with `author == approver` would fail the schema validator before it was sealed.

> **✓ Confirmation #9 — Lumière model-build IAM enforces author-approver separation in the chain schema**
> Every model-build chain entry has separate `author` and `approver` SSO-subject fields. The chain schema rejects an entry where the two are equal. Author and approver are typically a research-engineer (Aurélien or one of the team) and Hélène. Separation of duties is enforced by the chain, not by procedure — the SDK refuses the entry at write time per the §4.4 SDK-side enforcement pattern (the same shape as the §4.4 SDK-side `gen_ai.{request,response}.model` refusal that closes the source before MAC compute).

Diana, on the video bridge from Stuttgart, asked the cross-cutting question: "Klaus and Hélène — who can rotate the model-handover-entry approver list at Eberhardt? Specifically, who can change the policy that says 'a model-handover entry from Lumière requires Klaus's approval'?"

Klaus answered. "Two-of-three approval. Myself, the platform-engineering lead, the SecOps lead. The change itself is a chain entry. The chain trumps the org chart on this." (The policy-as-chain-entry pattern is the shape §10.18 endorses — the institution's CC8.1 names the policy version, the chain anchors policy changes, and the runbook cross-references the spec section the policy supports.)

Hélène, on the screen: "And on the Lumière side, the policy that says 'a model-build entry handed-over to Eberhardt requires Aurélien-author and Hélène-approver' is itself a chain entry under our IKM. Same shape — two-of-three approval to change the policy, and the policy change is sealed."

Diana wrote that down and underlined it.

She spent the next twenty minutes on credential-rotation cadence, scope-expansion procedures, and the §10.10 rotation-crossing-the-seal-boundary discipline. Eberhardt rotates the predictive-maintenance service-account credential quarterly. Each rotation crosses a seal boundary cleanly per §10.10.1 (hourly cadence does not apply here — Eberhardt's seal cadence is daily per §4.2.1). The `key_version` increments per §3 and the §4.4 attribute table; the §10.10.2 within-day algorithm-rotation discipline is documented but unused (Eberhardt has not rotated algorithms within a day).

She also pulled up the §10.1 key-fingerprint reconciliation evidence — Eberhardt maintains a tenant key registry that maps `tenant_id` to the IKM fingerprint and the KMS handle. The §4.4 `ffiec.chain.key_fingerprint` attribute on every chain entry asserts the looked-up IKM produces the recorded fingerprint BEFORE the verifier computes any MAC (per §7 step 8). The reconciliation is part of the §10.13 evidentiary-artifacts retention list — the institution's IT witness lays foundation from these artifacts under FRE 901(b)(9) without re-engineering the system.

> **✓ Confirmation #9a — Eberhardt's §10.1 key-fingerprint reconciliation and §10.10 rotation-crossing-the-seal-boundary discipline are documented and operating**
> Tenant key registry maps `tenant_id` to IKM fingerprint and KMS handle. Quarterly rotation of the predictive-maintenance service credential crosses the daily seal boundary cleanly per §10.10. The `key_version` increments and the §4.4 `ffiec.chain.key_fingerprint` attribute on every chain entry binds to the reconciliation registry. §10.10.2 within-day algorithm rotation is documented but not exercised. §10.13 evidentiary-artifacts retention covers the rotation history — SDK version manifest, HSM configuration, daily seal-job logs, change-management records, verifier output for the period.

---

### 🧪 12:00 PM CET — Joint Working Lunch by Video Bridge

The Stuttgart conference room had a small spread on a side table — Maultaschen, schnitzel, salad, a bowl of fruit. Klaus had insisted on Maultaschen, Swabian-style. The Paris side had a tartine with brie and a side of soup that Hélène had ordered up from the building's café. The video bridge stayed up. Both sides ate while talking.

Dawn put down her plate.

"Let's talk about what we've seen so far. Mike — the inference path."

Mike: "Eberhardt's chain holds at eight months. The inference entries are bound to model-handover entries by reference. The model-handover entries reference Lumière's chain by `audit.model_handover.provider_chain_entry_id` per §10.21. Klaus's two-of-three approval is on every handover. The §10.17 partition-ceremony events are present and well-formed."

Raj, on the screen from Paris: "Lumière's chain holds at four months. Build entries reference fairness-audit entries by ID. Author-approver separation is enforced by the schema. Aurélien-author and Hélène-approver on every build that handed over to Eberhardt."

Diana: "IAM on both sides is institutional and chain-governed. Two on-prem AD domains at Eberhardt. Google Workspace plus custom SSO at Lumière. Both sides have the policy-as-chain-entry pattern. Two-of-three approval to change either side's IAM policy."

Chen, on the screen from Paris: "I want to demonstrate the cross-anchor right now. Mike — share an Eberhardt model-handover entry. I have the corresponding Lumière build entry up."

Mike screen-shared his terminal. The model-handover entry `eb-mh-2026-04-01-lum-00012` from this morning was visible. The `audit.model_handover.model_artifact_sha256` field was clearly displayed: `sha256:7f2a91...c4d3`. The `audit.model_handover.provider_chain_entry_id` was `lum-mb-2026-04-01-bh-00031`.

Chen on his side opened the Lumière build entry. The `model_artifact_sha256` on Lumière's chain showed `sha256:7f2a91...c4d3`. Same value.

Chen ran a small Python helper on the Paris side that did a byte-equal compare of the three SHA-256 values pulled from the two terminals. The output:

```
Eberhardt audit.model_handover.model_artifact_sha256: 7f2a91...c4d3
Lumière   model_artifact_sha256:                       7f2a91...c4d3
match: True

Eberhardt audit.model_handover.model_card_sha256: 8e1b82...a4f1
Lumière   model_card_sha256:                       8e1b82...a4f1
match: True

Eberhardt audit.model_handover.fairness_audit_report_sha256: 9d3c72...e5f2
Lumière   fairness_audit_report_sha256:                       9d3c72...e5f2
match: True
```

Klaus, watching from the Stuttgart side: "That is the seam holding."

Hélène, on the screen: "That is what BMW asked for."

> **✓ Confirmation #10 — Cross-vendor anchor matches live across the partnership boundary per §10.21 cross-anchor verification**
> The `audit.model_handover.model_artifact_sha256`, `audit.model_handover.model_card_sha256`, and `audit.model_handover.fairness_audit_report_sha256` fields in Eberhardt's handover entry match the corresponding fields in Lumière's model-build entry bit-for-bit. The cross-vendor anchor composes per §10.21's bidirectional verification mechanism — independent verifier output on each side (§7), plus byte-equal hash join. RFC 8785 JCS canonicalization (§5) plus SHA-256 yields the byte-identical values for byte-identical inputs that the join depends on. This is the mechanism BMW's vendor-management team asked the joint engagement to demonstrate.

Dawn let the moment sit for a beat. Then she brought it back.

"Klaus — Hélène — for the deliverable, the language we'll use is that the cross-vendor anchor composes end-to-end without trust in either party's claim alone. Each side's chain is independently verifiable per §7. The match is what makes the joint claim. That language survives an examiner read on both sides."

Klaus nodded. Hélène nodded.

Tom: "Stefan Kuhn at BMW asked us specifically whether the cross-vendor anchor composes without trust. The answer is yes, by independent verification plus byte-equal hash match. That is the right answer."

Dawn: "Good. Save the live demo as evidence. We'll attach the screen-share frames to the deliverable per §10.13 evidentiary-artifacts retention discipline."

Klaus took the last bite of his Maultaschen. "Dawn — there is a question I want to put to you now, before the afternoon. We have been talking about the seam working. Tell me what does *not* work yet. Lumière and I both know there is something. Hélène and I would rather hear it now than at five-thirty."

Dawn looked at Raj on the screen. Raj had been the one who had found the 90-day retention issue this morning.

"Raj," Dawn said.

Raj walked them through it. The Lumière training-data shards rotated out at 90 days post-delivery. The Eberhardt deployment window was typically 9 to 18 months. If a regression appeared six months in and someone wanted to retrace which shards produced it, the shards were gone. The manifest hash was retained, the manifest document was retained, the chain entry was retained — but the shards themselves were not. This was the §10.20 partial.

Hélène was quiet for a beat.

"That is fair. The 90-day policy was set under GDPR data-minimization. I will say honestly — I had not modeled it against the Eberhardt deployment window. I should have. I will take a proposal to the founders this week to extend retention to the maximum deployment window with explicit GDPR justification, probably under Article 6(1)(f) legitimate interest balanced against Article 5(1)(c) data-minimization. Klaus — I want your engineering team's max-deployment-window number to use as the retention parameter."

Klaus nodded. "We can give you that. The longest we have run a model is twenty-two months. I would size it to twenty-four months." (Twenty-two months max-deployment-window plus a 60-day investigation buffer rounded up to twenty-four months — the §10.20 retention-floor calculation directly.)

"Twenty-four months it is."

Dawn: "Document the retention extension as a CAPA against this audit. The deliverable will name it as Finding #1 — Partial conformance against §10.20 — with the in-flight remediation. Tom — log it. The next handover entry should carry `audit.model_handover.training_data_retention_floor_days = 720` per §10.21 so the deployer's CC8.1 evidences the floor commitment from the chain."

Tom: "Logged."

Hélène was quiet for a moment more, then spoke again. "Dawn — let me say something for the record. The 90-day policy was not arbitrary. It was the answer we gave ourselves when we read GDPR Article 5(1)(c) the first time, four years ago, and decided we wanted to be conservative about training-data retention. It was a defensible position at the time. What I did not do was model the deployment-window asymmetry. The §10.20 retention floor mechanism — the longest active deployment window plus an investigation buffer, justified under Article 6(1)(f) legitimate interest tied to EU AI Act Article 12 — that is the right legal-basis architecture. We will adopt it. And the post-engagement spec section that names our case as the worked example is the durable answer to the question we should have asked four years ago. I would rather our policy be cited in spec text than rediscovered in a separate engagement."

Klaus, after a beat: "The chain that proves Lumière kept what they said they would keep is more useful than the contract that says they will."

Dawn: "That is the §10.21 `training_data_retention_floor_days` discipline exactly. The chain anchors the integer commitment on the handover entry; the deployer reads it from the chain, not from the contract binder; the post-close acquirer or post-incident regulator reads it from the chain alone. Round-17 M&A-G2 closed the contract-binding gap with the contract triple — `contract_id`, `contract_version`, `contract_hash_sha256` — and the §10.21 retention-floor field is the operational integer alongside it. The four attributes together are what make the §10.20 retention floor cryptographically auditable rather than procedurally trusted."

The bridge stayed up while everyone finished lunch. The mood was something like the last twenty minutes of a long meeting where the work was done and the people had settled into their seats. The §10.20 conversation had gone exactly as the spec section anticipates — the institution surfaces the asymmetry, the legal basis adjusts, the chain anchors the commitment, the next handover carries the integer. The mechanism the spec normates is the mechanism the field engagement adopts.

Dawn made a note for the deliverable's executive summary: *The §10.20 retention floor mechanism is not a remediation imposed by an outside auditor. It is the institution's own legal-basis architecture, made cryptographically auditable through the §10.21 model-handover schema. The audit surfaced the asymmetry; the spec named the resolution; the field engagement adopted the resolution within ninety minutes.*

---

### 🔄 1:00 PM CET — Mike on the API Layer (Stuttgart) and Elena on the CRM (Paris)

Mike's afternoon was the API surface. Eberhardt runs the inference service behind an internal API gateway with mutual-TLS authentication and a request-signing layer between the OEM-facing edge and the inference core. Every call to the inference service — from the OEM-customer-facing edge layer, from the internal validation harness, from the model-handover acceptance test — goes through this gateway. The gateway honored the §5.1 transport-encryption floor (TLS 1.3, with TLS 1.2 sunsetting 2028-01-01). OTLP transport identification per §4.4.3 is operating: the OTLP `Resource` carries `ffiec.chain.spec`, `service.name`, `service.version`, `ffiec.chain.posture`, and `ffiec.chain.format_version` so the receiver dispatches once per OTLP request before per-entry decode. Per §4.4.4, chain-of-custody traffic is severity-stamped at the receiver inside the `9..20` range so collector severity-filters do not silently drop chain entries.

Every call emits a chain entry. Request, headers (filtered for telemetry that could re-identify an individual vehicle owner), authorizer decision, downstream service, response code, response hash, latency. Sealed. The header-filtering happens at the SDK boundary per §10.22 redaction discipline — the canonical bytes that the per-event MAC covers are the redacted content; an attacker or insider reading the chain after capture sees only the redacted form. Mike checked: every inference entry that had carried any GDPR-special-category candidate field (vehicle owner names, billing addresses propagated by the OEM-edge layer in pre-redaction debug builds) emits the §10.22 `audit.redaction.*` attribute set with `disposition = "redacted_at_sdk"` (the conformant pre-MAC posture), `redacted_field_paths` enumerating the JSONPath identifiers, `redaction_method = ["sha256_hash"]` for the VIN and module-serial, and `policy_version` matching Eberhardt's published redaction-policy registry. The §10.22 binary posture statement is satisfied.

"Show me a call that returned a 5xx," Mike said.

Maximilian pulled up an entry from April 14. 503. The HSM at Sindelfingen had been unreachable for 90 seconds during a planned maintenance window. The chain entry recorded the failure. The downstream inference service had served a fallback `model-unavailable` response that explicitly told the OEM-side edge layer to fall back to its own conservative-default behavior — in this case, to suppress new alerts and continue surfacing the most recent confirmed alert if any. The fallback was itself sealed as a separate inference-shaped entry. The §4.3.1 HSM-unavailability notification SHOULD (72-hour notification window) was satisfied by Eberhardt's IR posture; the seal-job retry behavior aligned with the §10.16-style RTO discipline applied internally.

> **✓ Confirmation — API gateway and fallback paths are sealed end to end**
> Every inference API call — including failed authorizations, planned-maintenance HSM outages, and fallback responses — produces a sealed chain entry. The April 14 maintenance window shows up cleanly: 47 inference calls during the 90-second window, all routed to the fallback responder, all sealed. No silent gaps. The OEM-side edge layer received explicit fallback signals and behaved conservatively. Per §4.4 collector pass-through rule, no chain attributes were rewritten by the OTLP collector path during the outage.

Mike turned to the mutual-TLS authenticator. The mTLS certificate that Eberhardt presents to BMW's vehicle-side ingestion endpoint is rotated quarterly. The rotation produces a chain entry. The certificate fingerprint is bound to the rotation entry. BMW has the public-key portion of Eberhardt's quarterly rotation calendar in their vendor-management system.

> **✓ Confirmation — mTLS certificate rotation is chain-coupled with downstream OEM**
> Every quarterly rotation of the Eberhardt-to-BMW mTLS certificate produces a sealed chain entry with the new fingerprint and the rotation approval. BMW receives the public-key calendar through their vendor-management portal. The rotation is auditable from both ends. Eberhardt's rotation crosses the seal boundary cleanly per §10.10 (rotation crossing the seal boundary normative discipline).

Meanwhile, in Paris, Elena had spent thirty minutes with Lumière's revenue-operations lead — a person named Sophie Lacombe — going through the customer-engagement CRM. Lumière uses Pipedrive for sales and customer-engagement management. The Pipedrive instance has Eberhardt as the largest active account, with notes from quarterly business reviews, contract-renewal calendar entries, and contact-history.

The Pipedrive is not chain-coupled. It is a commercial CRM, not a technical-evidence system. The fairness-audit reports, model-handover records, deployment-tracking — none of those flow into Pipedrive. Pipedrive holds the relationship history, not the engineering history. Per §10.19 chain-coverage map discipline, Pipedrive lands on Lumière's map under the "third-party systems out of contractual chain reach" category, with a documented evidentiary substitute (Pipedrive's own audit trail) and an explicit rationale that the chain does not anchor commercial-relationship history because that history is not technical evidence.

Elena flagged this and moved on. *Out of scope for the AI-chain audit. Same shape as Helmstad's CTMS-as-KOL-CRM pattern, but cleaner — Pipedrive at Lumière holds nothing the chain needs to anchor. The CRM is the CRM and the chain is the chain. The §10.19 chain-coverage map names the boundary explicitly so the examiner does not discover it per-finding.*

> **✓ Confirmation — Lumière's CRM (Pipedrive) is correctly scoped as out of chain coverage per §10.19 chain-coverage map**
> Pipedrive holds commercial relationship history — contracts, renewal calendars, contact notes. It does not hold technical-evidence artifacts. The chain does not anchor Pipedrive entries because the chain does not need to. The boundary is documented and intentional. Per §10.19 the chain-coverage map enumerates Pipedrive under "third-party systems" with the documented evidentiary substitute. The map is version-stamped per §10.19 with a `chain.coverage_map_published` operational event in §10.2 carrying `coverage_map_version`, `effective_utc`, and `coverage_map_sha256` (Round-17 M&A-P3 — the chain anchor lets a lookback auditor determine which map version was in force on a given date). Same architectural decision as the Helmstad-pattern out-of-scope CRM.

Elena closed her notebook on the CRM section by 1:30 and moved over to help Chen on the data-pipeline walk that was about to start.

---

### 🧬 2:00 PM CET — Chen on the Lumière Training-Data Pipeline (Paris)

Chen had been waiting for this part. The chain at the model-build boundary is one thing. The chain at the training-data-input boundary is the other. They are not the same thing.

Aurélien walked Chen through the training-data ingestion pipeline. The raw signals — battery cell voltages, temperature curves, charge cycle counts, motor-controller telemetry — came from Eberhardt's field deployment, anonymized and aggregated by Eberhardt's data team, transferred to Lumière over a GDPR Article 28 processor-agreement-governed encrypted channel. The transfer arrived as a daily tarball, signed by Eberhardt's transfer key. Lumière's ingestion service verified the signature, computed a SHA-256 of the tarball, and recorded the hash as a referenced artifact in a chain entry of type `training-data-ingestion`. The records were then unpacked, parsed, validated against the schema, and written into the training-data warehouse. The warehouse row carried a foreign key to the ingestion chain entry.

"That is the same shape as the Eberhardt SFTP-ingestion pattern in reverse," Chen said. "Eberhardt sends, you receive. The chain entry on your side records what you received. The chain entry on Eberhardt's side records what they sent. The hashes match."

"Yes."

Chen pulled up a recent ingestion entry and the corresponding Eberhardt outbound entry. The tarball SHA-256 matched. The transfer key fingerprint matched. The tarball file count matched. The byte count matched.

> **✓ Confirmation — Training-data transfer is chain-coupled at both ends across the partnership boundary**
> Eberhardt's outbound transfer entry and Lumière's ingestion entry record matching SHA-256, file count, byte count, and signing-key fingerprint. The transfer is GDPR-Article-28-governed. The chain is the integrity record on both sides. This is the second cross-vendor anchor — alongside the §10.21 model-handover anchor — and it composes the same way (RFC 8785 JCS canonicalization per §5 plus SHA-256 yields the byte-equal hash join).

Chen circled the transfer-key-fingerprint match in his notebook.

"GDPR cross-border transfer," Chen said. "Stuttgart to Paris. Both within the EU. The Article 28 processor agreement governs it. Spec §10.21 names the cross-border-transfer composition explicitly — when the model handover crosses jurisdictions, the institution SHOULD also emit `audit.cross_border_transfer.*` on the same chain entry. Within-EU transfers do not need it under GDPR — same regime end-to-end — but the spec explicitly cites the Eberhardt-Lumière case as the exemplar where the explicit attribution is examiner-friendlier for non-EU vendor-management readers and for the BMW joint-supplier audit."

Aurélien paused. "The current entry carries the source country code, the destination country code, the processor-agreement reference, and the Article 28 reference. It does not carry the §4.4.1 `audit.cross_border_transfer.*` attribute set explicitly."

"Per §10.21 the spec says SHOULD for within-EU transfers — examiner-friendlier rather than required. That's the exact wording the Eberhardt-Lumière case is named under in the spec. Want to read the §4.4.1 schema?"

Aurélien did. Chen walked him through it: `audit.cross_border_transfer.contract_id` (the Article 28 processor-agreement identifier), `audit.cross_border_transfer.contract_version` (versioning lets auditors detect a contract amendment between two chain entries), `audit.cross_border_transfer.contract_hash_sha256` (binds the chain entry to the document), `audit.cross_border_transfer.source_jurisdiction` (`DE`), `audit.cross_border_transfer.destination_jurisdiction` (`FR`), `audit.cross_border_transfer.lawful_basis_type` (one of the §4.4.1 enumeration — `intra_group_agreement` is the right value here since Eberhardt and Lumière are in a long-running processor relationship, not full intra-group, so the right value is actually `standard_contractual_clauses` despite the within-EU framing being optional). The Round-17 NAIC-P4 elevation made the family REQUIRED when the chain entry is subject to a regulator-named privacy regime in the institution's CC8.1; for within-EU SHOULD posture remains.

Aurélien thought. "Probably yes — we should emit it. The CNIL would not need it. The German LfDI would not need it. But the BMW vendor-management read would want it. And if we ever add a non-EU customer who wanted to read the chain on the Eberhardt deployer side — they would find it useful. The §10.21 cross-border composition note flags exactly our case. Within-EU SHOULD, not REQUIRED — but the SHOULD is for examiner-readability, and the BMW examiner is the audience that benefits."

Chen: "And the institution-side trigger condition matters here. Per §4.4.1, institutions whose CC8.1 names a privacy-regime trigger emit the attribute set whenever the trigger condition holds. Within-EU GDPR Article 28 transfer is the trigger condition you would name in CC8.1. The Round-17 NAIC-P4 elevation made the family REQUIRED when the chain entry is subject to a regulator-named privacy regime in the institution's CC8.1 — for state-insurance-privacy this is binding, for GDPR within-EU this is SHOULD because the regime is the same end-to-end. The spec is careful about that distinction."

Aurélien: "We will name the trigger condition in CC8.1 and emit the attribute set on every training-data-transfer entry going forward. The schema is already normative in §4.4.1; the work is to add the institution-side emission to our ingestion pipeline. Two-week implementation, four-week rollout. I will document the change in our CC8.1 control description and cross-reference §4.4.1 per §10.18 runbook discipline."

Chen wrote that down. *The §10.18 cross-referencing rule applies — the runbook section that describes Lumière's Eberhardt-bound training-data ingestion needs to name §4.4.1 explicitly so the reviewer's path is runbook → spec → design → audit procedure → SOC engagement → examiner workpaper. Without the §10.18 cross-reference the reviewer's path is "read the runbook, then read the entire spec, then attempt to map," which is a CC8.1 discoverability deficiency.*

Chen wrote that down. *Same shape as Sun-Won's PIPA Section 28 cross-border-transfer attribute finding from a prior engagement — that finding drove the §4.4.1 `audit.cross_border_transfer.*` family in the first place. The Eberhardt-Lumière case is the one §10.21 cites by name as the within-EU exemplar where the attribute is examiner-friendlier rather than required.*

Reading the field finding back through the post-amendment spec:

> **⚠️ Finding #2 — `audit.cross_border_transfer.*` attribute family is implicit, not explicit per §10.21 cross-border composition recommendation — Nit (within-EU SHOULD)**
> Spec §4.4.1 (`audit.cross_border_transfer.*` attribute family, originally landed via the Sun-Won engagement and elevated under Round-17 NAIC-P4) plus §10.21's cross-border-composition recommendation handle this case. Within-EU transfers do not need the attribute set under GDPR (the spec is explicit: "the regime is the same end-to-end") — Eberhardt-Lumière transfers are Stuttgart-Paris under Article 28 processor agreement, GDPR end-to-end. Per §10.21 cross-border-composition note: the spec recommends the explicit attribute set for non-EU regulators and vendor-management auditors who are unfamiliar with the GDPR-internal model — and the spec **names the Eberhardt-Lumière case as the exemplar**. Severity: **Nit** (within-EU SHOULD, not REQUIRED). Hélène commits to emitting the explicit `audit.cross_border_transfer.*` attribute set on training-data transfer entries within four weeks for examiner-friendliness — primarily for the BMW German-English-French reader matrix and for non-EU customers Lumière may add in the future. Schema is already normative; the institution-side emission is the work item.

Chen moved on to the training-data manifest binding. Aurélien showed him the manifest document for the v4.2.1 training run. Every shard listed with source, ingestion date, de-identification scheme, consent basis, and SHA-256. The manifest itself was sealed — its hash was the `training_data_manifest_sha256` in the model-build entry. The optional `audit.model_handover.training_shard_manifest_sha256` field on the next handover entry (per Round-17 M&A-P2) will make the manifest hash directly bound on the chain at handover, not just transitively via the build entry; that closes the deal-window-lookback gap where the chain proved provider delivery but did not bind the enumerated shard list.

Chen also asked Aurélien to walk him through the routing layer. Lumière operates a routing classifier in front of the model-build feature-extraction pipeline — when a new training-data shard arrives, a content-type classifier decides whether the shard is in-scope for battery-health, motor-controller, or general-vehicle-telemetry feature extraction. Per §4.4.1, when routing is driven by a classifier, the institution MUST emit an `audit.routing.classifier_output` chain entry BEFORE the `audit.routing.attempt` event the classifier informs. The chain entry records `audit.routing.classifier_name`, `audit.routing.classifier_version`, `audit.routing.classifier_input_hash`, `audit.routing.classifier_scores`, `audit.routing.classifier_decision`, and `audit.routing.classifier_confidence`. Aurélien pulled up a recent classifier_output entry — the shard hash matched the manifest entry, the classifier scored battery-health at 0.94, motor-controller at 0.05, general at 0.01, and the decision was battery-health. The chained `audit.routing.attempt` followed, linked by `parent_run_id` / `parent_seq` per §4.4 with the classifier_output as parent.

> **✓ Confirmation — Pre-routing classifier capture is in place per §4.4.1 Round-17 NAIC-P4 / Sun-Won errata**
> Lumière's routing classifier emits `audit.routing.classifier_output` chain entries BEFORE the `audit.routing.attempt` events they inform, per the §4.4.1 normative rule for classifier-driven routing. The chain entry carries the classifier identity, version, input hash, per-class scores, decision, and confidence — the rationale-recoverable-from-chain shape the spec normates. Eberhardt-side has no analogous classifier (the inference pipeline is single-model per OEM-fleet binding); the §4.4.1 emission requirement does not apply to Eberhardt's inference path. The shape is correctly conditional.

"And the consent basis on each shard," Chen asked. "What is it?"

"GDPR Article 6(1)(f) legitimate interest, with the data-protection-impact assessment for the predictive-maintenance use case as the documented justification. The DPIA is signed by Hélène and the customer-side data-protection officer at Eberhardt. The DPIA hash is in the manifest. The §10.20 retention-floor extension to 24 months we agreed at lunch will be re-justified under the same Article 6(1)(f) framework, balanced against Article 5(1)(c) data-minimization, with Article 35 DPIA naming the longest-deployment-window justification — exactly the §10.20 GDPR data-minimization-tension resolution mechanism."

> **✓ Confirmation — Training-data manifest binds shards to consent basis with DPIA hash; §10.20 retention floor extension uses the same Article 6(1)(f) mechanism**
> Every shard in the training-data manifest is bound to its consent basis (GDPR Article 6(1)(f) legitimate interest), its DPIA reference, and its de-identification scheme. The DPIA itself is hashed into the manifest. The manifest hash is sealed in the model-build entry. End-to-end consent provenance is auditable from the model artifact back to the legal basis. The §10.20 retention floor extension reuses the same Article 6(1)(f) determination tied to EU AI Act Article 12 logging obligations — the spec's data-minimization tension resolution mechanism applies directly to the retention extension, no novel legal-basis work required.

Chen closed his notebook on the training-data section. *Two cross-vendor anchors. §10.21 model-handover and the training-data-transfer pattern. Both compose. The chain reaches the legal-basis layer.*

---

### 🧬 2:00 PM CET — Luis on Logs and Pipelines (Stuttgart)

Luis was on the Stuttgart side, at the same hour. He had spent the morning quietly going through the Eberhardt log-pipeline and the daily-seal job. By 2:00 PM he had a clean read on it.

The seal job ran at 23:55 local time at Sindelfingen. It pulled the day's chain entries, computed the daily Merkle root per §4.2 in `(run_id, seq)` ascending order (the canonical ordering — implementations MUST NOT use `received_at` or `captured_at` for Merkle ordering, per the §4.2 normative ordering rule), signed it with the on-prem HSM per §4.3 (Ed25519 over the v1.0a 10-line `sign_payload`, soon to migrate to the v1.0b 12-line form per the §12 change log), and published the seal record to an internal seal-archive bucket. The seal-archive bucket had versioning enabled, was on an immutable-storage tier with a 7-year retention lock under TISAX, and was replicated to a second data center in Munich.

The CloudWatch-equivalent — Eberhardt used Splunk Enterprise on-prem — had log-deletion-prevention configured at the index level. Splunk admin operations themselves emitted chain entries.

Luis asked the question he always asked.

"Who can disable the daily seal?"

The Eberhardt SecOps lead — a person named Bettina Hofer — answered. "No one, unilaterally. The seal job is a Kubernetes CronJob in a namespace where the deployment manifest is itself a chain entry under `chain_kind = "operational"`. To disable the seal you would have to issue a chain-bound deployment change with two-of-three approval. The two-of-three are Klaus, the platform-engineering lead, and me. We have not disabled it. We do not plan to."

"What if the HSM is unreachable at seal time?"

"The seal job retries for 90 minutes. If it cannot complete, it emits a `seal-failure` chain entry — the chain itself records the failure — and pages the SecOps lead. There is a one-time-only manual seal-completion procedure documented for the case where the HSM has a hardware failure; that procedure requires three-of-three approval and the manual seal is itself a chain entry of a special type. Per §4.3.1 the HSM-unavailability notification SHOULD applies — we file the notification within 72 hours of the unavailability if the seal cannot complete on retry."

> **✓ Confirmation — Seal job is chain-bound and chain-failure-aware on the Eberhardt side per §4.2 / §4.3 / §4.3.1**
> The daily seal job at Sindelfingen runs as a chain-governed Kubernetes CronJob. Disabling the job requires two-of-three approval and produces a chain entry per §10.2 operational events. Seal failures are themselves chain entries. Manual seal completion exists as a documented procedure with three-of-three approval and a special-type chain entry. The chain is reflexive about its own operation. §4.2.1 cadence (daily) is documented; §4.2.2 day-boundary semantics applied (late-arriving events use `ffiec.chain.late_binding = true` per §4.4 attribute table). The empty-day discipline (§4.2 — every tenant-day MUST receive a seal record, even tenant-days with zero events) is honored: Eberhardt's seal-archive bucket has a seal record for every day since deployment, including the holiday closures over Christmas 2025 when zero inferences ran.

Luis wrote that down. He moved to the log-replication path between Sindelfingen and Munich. The replication job was a separate CronJob, also chain-governed, and emitted the §10.2 `master.cross_region_replication_completed` event at the cadence the institution's CC8.1 named (per §10.15 Pattern A invariant 5). He spot-checked the cross-region replication-completion entry from the prior night. PASS, per-region count summed correctly to the seal-region count, no replication-loss anomaly. Eberhardt's run-resume discipline per §10.25 was in place — the SDK acquired the chain tail before emitting the next entry through the three-place lookup (in-memory → local persistence sidecar → ledger query rejoin), single-writer-per-run was enforced through file locks, and the ledger ingestion cross-check on `(prev_hash, seq)` monotonicity was active at every batch ingest. The §4.4 genesis-block uniqueness rule plus the §10.25 ingestion cross-check close the silent-restart attack class at both the SDK and ledger layers.

He asked Bettina the same question Lumière would be asked in Paris by Raj an hour later: "Where does the seal-archive bucket live, jurisdictionally?"

"On-prem in Sindelfingen, replicated to on-prem in Munich. Both are German jurisdiction. We do not replicate to cloud or to another country. BSI IT-Grundschutz baseline plus our own preference. The §10.15 Pattern A discipline applies — single seal region per tenant per `seal_date`, the seal region is Sindelfingen, Munich is the replication region. The `master.cross_region_replication_completed` event reads the per-region count synchronously from the replication pipeline state at emission time per the Round-17 errata that closed the Atrio cache-staleness Partial."

"Mittelstand," Luis said.

Bettina half-smiled. "Mittelstand."

> **✓ Confirmation — Eberhardt's multi-region resilience operates §10.15 Pattern A with synchronous-read replication-completion events per Round-17 third errata**
> Single seal region (Sindelfingen) per tenant per `seal_date`. Munich is the replication region. The `master.cross_region_replication_completed` operational event records per-region count, completion timestamp, and target seal region — synchronous read against the replication pipeline state at emission time, not poll-cached (per Round-17 third errata closing the Atrio Partial). Per-region event-count reconciliation is operating; the seal region's count equals the sum of regional counts. The institution's CC8.1 names Pattern A and the seal-region failover procedure per §10.15.

Luis pulled the seal record itself for April 8 — the day of the BMW iX inference Mike had been walking earlier. The seal carried the §4.3 `sign_payload` form: `format_version`, `hkdf_inputs_digest`, `seal_date`, `cadence`, `dev_mode`, `merkle_root`, the day's `key_versions` ascending list, the day's `kms_handle_uris` canonical sorted-distinct form. Eberhardt was on the v1.0a 10-line form at the time; the migration to the v1.0b 12-line form (binding `key_versions_canon` and `hex(kms_handle_uris_digest)` under the HSM signature, per Round-17 NIST-G1 / NIST-G2) is scheduled for the next quarterly maintenance window. Per §12 change log, v1.0a chains remain verifiable under any v1.0b verifier without re-sealing — the dispatch is monotonic, so existing v1.0a artifacts continue to verify under their original form indefinitely.

> **✓ Confirmation — Seal record carries §4.3 `sign_payload` v1.0a 10-line form correctly; v1.0b migration scheduled per §12 monotonic dispatch**
> Eberhardt's daily seal at Sindelfingen carries the v1.0a 10-line `sign_payload` form per §4.3 (with `sign_payload_version = "v1.0a"`). The lowercase-hex normative form is followed; LF-only line termination is followed; trailing-newline absence is followed; ISO 8601 date-only is followed; hex 64-character zero-padding on the merkle root is followed. The migration to v1.0b 12-line form (binding `key_versions_canon` and the kms_handle_uris digest under the HSM signature) is scheduled for the next quarterly maintenance window. Per §12 monotonic dispatch, v1.0a artifacts remain verifiable under any v1.0b verifier indefinitely.

---

### 📊 3:00 PM CET — Joint Reconciliation Test by Video Bridge

The video bridge was back up at 3:00 PM. Dawn had asked both sides to pick a single deployment and trace it end-to-end through both chains. She had asked Klaus and Hélène to be on the bridge for it. Tom had set up screen-recording per §10.13 evidentiary-artifacts retention discipline so the live trace would land in the deliverable as an integrity-bound artifact.

Dawn framed it for the bridge. "Two chains. One BMW vehicle. Seven legs to traverse. We're testing whether the §10.21 cross-vendor anchor composes end-to-end without trust in either party's claim — independent verifier output on each side, plus byte-equal hash join at every boundary. If any leg fails, we surface it. If every leg passes, we have the joint claim."

Maximilian on the Stuttgart side picked a recent one: an `urgent-service-required` predictive-maintenance alert from April 8 for a BMW iX. Anonymized as `Customer-X` for the deliverable. The vehicle VIN was hashed in the chain entry (the chain captures the hash; the VIN itself is not in the chain — §1.2 epistemic-scope discipline plus §10.22 redaction discipline composed at the SDK boundary).

Maximilian read out the trace. Dawn watched it on the screen. Tom took notes.

**Step 1.** The vehicle VIN-hash on the BMW iX. (BMW would have the unhashed VIN on their side. The hash matches.)

**Step 2.** The Eberhardt inference chain entry `eb-bh-2026-04-08-bmwix-00874`. Sealed. §7 verifier PASS in 4 seconds, twelve steps, exit code 0 per §10.12. The entry's `model_handover_entry_ref` points to `eb-mh-2026-04-01-lum-00012`.

**Step 3.** The Eberhardt model-handover entry `eb-mh-2026-04-01-lum-00012`. Sealed. §7 verifier PASS in 4 seconds. The entry's `audit.model_handover.provider_chain_entry_id` points to `lum-mb-2026-04-01-bh-00031` per §10.21.

**Step 4.** Hélène — by video bridge — granted Mike read access to the Lumière chain for the matching entry. Mike ran the Lumière verifier against `lum-mb-2026-04-01-bh-00031` from his Stuttgart terminal. Verifier PASS in 4 seconds. The entry's `model_artifact_sha256` matches Eberhardt's `audit.model_handover.model_artifact_sha256` by byte-equal compare (RFC 8785 JCS plus SHA-256 yields the byte-identical values per §5). The entry's `fairness_audit_entry_ref` points to `lum-fa-2026-03-31-bh-00007`.

**Step 5.** The Lumière fairness-audit entry `lum-fa-2026-03-31-bh-00007`. Sealed. §7 verifier PASS. The fairness-audit-report SHA-256 matches the value sealed in Eberhardt's model-handover entry (`audit.model_handover.fairness_audit_report_sha256`).

**Step 6.** The Lumière model-build entry's `training_data_manifest_sha256` resolves to a manifest document. The manifest document hashes to the value sealed. The manifest lists the shards used for v4.2.1 training. (Within 90 days; the shards themselves are still present for this April-1 build. After the §10.20 retention extension lands, the shards will remain available for 24 months from delivery — closing the §10.20 forensic gap for any future regression detected in the deployment-window tail.)

**Step 7.** Each shard in the manifest is hash-bound to its source. The Eberhardt outbound transfer entries on the Eberhardt side and the Lumière ingestion entries on the Lumière side reconcile by tarball SHA-256. End-to-end consent provenance traces back to the DPIA sealed in the manifest. The §10.20 GDPR Article 6(1)(f) legitimate-interests determination tied to EU AI Act Article 12 is the legal-basis layer the chain reaches.

The bridge sat in silence for a moment after Maximilian finished reading.

Dawn spoke first.

"Seven legs. Two chains. One BMW vehicle. End-to-end traversable with byte-equal hash matches at every cross-vendor anchor. That is the joint claim."

Klaus nodded once.

Hélène, on the screen: "That is what BMW asked for. That is what the AI Act Article 11 / 12 / 16 conformity assessment asks for. The §1.2 epistemic-scope language is what makes this defensible — the chain proves what was deployed and that the record was not tampered with; it does not prove the model is correct or unbiased. Those are separate evidence regimes. The conformity-assessment file references the chain for the integrity claim and the fairness-audit report for the bias claim, not the chain for everything."

> **✓ Confirmation — Joint end-to-end reconciliation traverses BMW vehicle → Eberhardt inference → Eberhardt §10.21 model-handover → Lumière model-build → Lumière fairness-audit → Lumière training-data manifest with byte-equal hash matches at every cross-vendor anchor**
> One BMW iX vehicle's predictive-maintenance alert traced through seven legs across two independently-IKM-governed chains (per-tenant HKDF binding per §4.1). §7 verifier PASS at every leg in 4 seconds per check, twelve steps, output format normative per §7. Byte-equal hash matches at every cross-vendor anchor (§5 RFC 8785 JCS canonicalization plus SHA-256). End-to-end provenance from the OEM customer's vehicle back to the GDPR consent basis on the training data, with the DPIA hashed into the manifest. §1.2 epistemic-scope discipline distinguishes what the chain proves (deployment integrity) from what it does not (model correctness, bias-freeness, or policy compliance) — the latter live in the fairness-audit report which the chain hash-anchors but does not evaluate. This is the joint-supplier audit BMW's vendor-management asked for.

Mike, on the Stuttgart side, looked up from his terminal.

"That trace took eleven minutes from kickoff to closing PASS. Including the Hélène-grants-credentials step."

Dawn: "Document the elapsed time. The deliverable will note that an end-to-end joint trace can be completed in under fifteen minutes once both sides cooperate. That is a number BMW will want to see."

Tom: "Logged."

Dawn turned the question back to the room. "What does the trace not cover, that BMW might ask about?"

Raj answered first. "BMW's vehicle-side data. Step 1 was the VIN-hash on the BMW iX, but the unhashed VIN, the raw sensor stream, the BMW-side ingestion logs — those are BMW's chain or BMW's equivalent of one. They are off our scope. The trace assumes BMW's side reconciles. If BMW's input does not match Eberhardt's `input_telemetry_hash`, the difference is in the BMW-to-Eberhardt transport. That is BMW's integration question, not Eberhardt's or Lumière's."

Chen: "And the trace does not cover the BMW-OEM customer's own decision-making once the alert reaches their vehicle. Whether the customer takes the car in for service, whether the BMW dealership confirms the alert against their own diagnostic — those are BMW-customer-side artifacts. The chain stops at the inference output sent to BMW. Beyond that boundary is BMW's audit trail."

Mike: "And the trace covers v4.2.1 specifically. If BMW asked about v3.9.0 or any prior version, we would walk a different model-handover entry — but the same seven-leg shape. The chain captures every model-handover historically; we can walk any version's trace at any time. The §10.20 retention extension matters here — once the 24-month retention floor is in place, we can walk the training-data layer for any model that has been deployed in the last 24 months. Before the extension, Lumière's training-data layer is reachable only within 90 days of delivery."

Klaus, after a beat: "So the trace today is fully traversable for v4.2.1 because we are within Lumière's 90-day retention. After June 30 — ninety days post-April-1 delivery — the training-data shards rotate out under the current policy. The §10.20 extension to 24 months keeps the trace fully traversable through the deployment window."

Hélène: "Yes. The retention extension is what keeps the seven-leg trace alive across the deployment window. The §10.21 chain anchors prove what was deployed; the §10.20 retention floor keeps the training-data root-cause path retrievable. The two amendments compose."

Dawn: "And that is the language the deliverable will use. The chain at Eberhardt proves what was deployed for any inference at any time. The chain at Lumière proves what was built for any model handover at any time. The §10.20 retention floor keeps the training-data root-cause path retrievable through the deployment window. The §10.21 cross-anchor composes the two chains without trust. The trace is fully traversable while the retention floor holds."

Tom: "Logged for the executive summary."

---

### 🧪 3:30 PM CET — The Audit-Report-Languages Question

Diana was the one who had been turning the language attribute over in her head since 10:00 AM. She brought it back up by video bridge.

"The fairness-audit report is in French. The Eberhardt deployment documentation is in German. The model card is bilingual French-English. BMW's vendor-management lead — Stefan Kuhn — what does he read?"

Klaus: "Stefan reads German first, English fluently, French only well enough to follow a technical paper."

Hélène: "The fairness-audit report has an executive summary in English. The full body is in French. The translation of the full body to English would take an external translation effort. We have not done it for v4.2.1."

Diana: "So a German examiner reading the joint deliverable would have the model card in two languages and the fairness-audit body in only one."

Klaus thought for a beat. "The BSI examiner reads English fluently. The LfDI Baden-Württemberg reads German. The CNIL reads French. The BMW vendor-management read is German-English. So the reader matrix is German, English, French — three languages. The model card is bilingual French-English. The fairness-audit body is French only with an English executive summary. The deployment-side documentation at Eberhardt is German with English appendices."

Dawn: "So an examiner reading across the seam — say BSI plus CNIL — has access to the artifacts in their preferred languages, but the *individual* fairness-audit body is French only."

Hélène: "Yes. We can produce an English translation of the fairness-audit body for v4.2.1 within four weeks. For new models we can include English-French bilingual full-body audit reports as a standing practice. And as we add German-language translations for the BMW reader matrix, we'll populate `audit.model_handover.audit_report_languages = ["fr", "en", "de"]` per §10.21."

Diana: "And the chain entry's `audit.model_handover.audit_report_languages` field is what flags this. It is captured. The reader knows from the chain what languages the report is in. The §10.21 plural-array discipline lets the chain entry surface available translations without a schema change. As more translations are produced, the array carries them."

Dawn: "Right. That is a procedural finding, not a chain-integrity finding. The chain captures the language array correctly per §10.21. The schema is already plural — Round-17 M&A-N2 closure. The procedural recommendation is to produce the English translation within four weeks and adopt bilingual-or-trilingual audit bodies as standing practice."

Reading the field finding back through the post-amendment spec:

> **⚠️ Finding #3 — `audit_report_language` is singular in the Lumière fairness-audit entry; §10.21 normates `audit.model_handover.audit_report_languages` plural-array (closed by Round-17 M&A-N2) — Nit**
> Spec §10.21 (closed by Round-17 M&A-N2 plural-array discipline) handles this case. The plural-array form (`audit.model_handover.audit_report_languages`) is normative — the spec elevates it because multilingual reports are common when models cross jurisdictions and downstream customers (e.g., automotive OEMs) include multilingual examiner audiences. The Eberhardt-Lumière BMW joint-supplier audit is the spec's worked example for the plural form; §10.21 names it directly in the "Audit-report languages worked example" paragraph. At capture time the Lumière fairness-audit entry uses singular `audit_report_language` (a Lumière-internal attribute name predating the §10.21 normative shape); the Eberhardt §10.21 model-handover entry already uses the plural-array `audit.model_handover.audit_report_languages` form correctly. Severity: **Nit** (Lumière-side schema-attribute alignment to the §10.21 plural form on the fairness-audit entry; the §10.21 normative attribute on the model-handover entry is already correct). Hélène commits to producing the v4.2.1 English translation within four weeks and aligning the Lumière fairness-audit entry to the plural-array form on the next chain-schema revision; the §10.21 model-handover plural form is already correctly emitted, so the BMW reader-matrix discoverability is already flowing through the chain.

Hélène: "Agreed. We will produce the v4.2.1 English translation in the next four weeks and adopt bilingual-or-trilingual audit bodies as standing practice for new models. The Lumière-internal attribute alignment to the §10.21 plural form is a one-line schema-change."

Dawn: "Document it both ways. Procedural CAPA at Lumière for the standing practice. Schema alignment for the next normative revision on Lumière's side. The §10.21 plural form on Eberhardt's side is already correct."

Tom: "Logged."

Diana, on the bridge: "There is one more thing about the language array I want to flag for the deliverable. The §10.21 worked-example paragraph on plural-array discipline names the case where a vendor-management auditor without German finds the English translation through the chain itself. That is the BMW examiner case verbatim. The chain alone serves all three audiences — German for the LfDI, French for the CNIL, English for BMW vendor-management — without separate per-language attribute emissions. The plural-array shape carries the discoverability."

Klaus: "And as the model evolves and the array grows — when v4.3 lands with German added to the bilingual French-English pair — the chain entry for that handover carries `["fr", "en", "de"]` and the BMW German examiner reads German from the chain without re-engineering the schema. The shape carries forward."

Hélène: "And the model card already does this — the model card is bilingual French-English. We will harmonize the fairness-audit-body language coverage with the model card going forward, with German added when the BMW reader matrix calls for it. The §10.21 plural-array form lets the chain reflect that growth without schema-change friction. That is the discipline I want for the team going forward — translations are a chain-anchor decision, not a procedural courtesy. The chain entry is the integrity-bound record that the translation exists."

Dawn: "And the per-language CC8.1 cross-referencing per §10.18 means that the institution's CC8.1 control description for the fairness-audit family — typically maintained in the institution's primary working language — cross-references the per-language audit-report bodies by file name, available languages, and §10.21 attribute mapping so a reviewer in any of the three languages can locate the right body without parsing the institution's free-form schema."

Tom: "I will draft the CC8.1 cross-reference as part of the deliverable. The §10.18 runbook discipline applies."

---

### 😬 3:45 PM CET — Friction: The 2024 Model-Drift Incident

Klaus had been quiet through the language discussion. He came back to the bridge as it ended.

"There is one more thing I want to put on the table. The 2024 model-drift incident. Hélène — how do you want to talk about it?"

Hélène took a beat before answering. The video screen caught the small motion as she sat back in her chair.

"Honestly. The 2024 incident was on us. We delivered a battery-health model that had silently regressed after a training-data update. The validation set had not caught the regression because the new training data shifted the distribution of the test set as well. The model passed our internal validation but was actually worse on the field distribution. We delivered it to Eberhardt. Eberhardt deployed it to the OEM fleet. About two weeks of degraded predictions before someone at BMW noticed and asked Eberhardt to investigate. Eberhardt asked us. We rolled back. Total degraded-prediction window was something like fifteen days. No safety incidents — the regressions were over-triggering of `service-recommended` alerts, not under-triggering of `urgent-service-required` — but the OEM customer-experience cost was real. We did not have the chain at the time. We could not retrace which deployments had received the regressed model versus the prior good model. We had to do it by reconstruction from delivery logs. It took three days. Three days that would have been three minutes with the chain."

Klaus, very evenly: "Dawn — does the chain prevent another 2024?"

Dawn took a beat.

"It does not prevent. It detects. The chain proves which model produced which inference per the §10.21 cross-anchor. It does not prove the model is good — that is exactly the §1.2 epistemic-scope split: the chain proves what was deployed, not whether it was correct. The detection happens earlier because the chain gives BMW or Eberhardt an immediate way to ask 'which model produced this anomaly?' without three days of reconstruction. The chain shrinks the silent-regression window from two weeks to whatever the field-monitoring cadence is. If field monitoring is daily, the window is one day. If it is weekly, the window is one week. The chain does not change the field-monitoring cadence."

Klaus: "So the answer is faster detection, not prevention."

Dawn: "Yes. Faster detection. The chain is the audit trail. It is not the validator. §1.2 makes that distinction normative."

Hélène: "And the §10.20 retention extension we agreed to this morning matters here. If a regression appears six months after deployment — which is the 2024 shape — the chain detects fast, but the *retraining* root-cause analysis needs the training shards. The 90-day window is too short. Twenty-four months — the maximum-deployment-window number plus investigation buffer — is the §10.20 retention floor. The forensic depth on the deployer side is governed by the provider-side retention floor per §10.20 — that's why the asymmetry mattered."

Dawn: "Right. The chain is the detection layer per §1.4 compositional security. The §10.20 retention extension is the root-cause-analysis layer. Both matter. Both are findings of this engagement."

Klaus, after a pause: "Detection without prevention is still better than two-week silent regression."

Dawn: "Yes. Materially better. But the deliverable has to be honest about the §1.2 distinction."

Tom: "I will draft the language for the executive summary. The chain detects faster; the chain does not prevent regression; the field-monitoring cadence and validator quality are separate concerns. The §10.20 retention extension supports retroactive root-cause analysis for the long-tail regression case."

Klaus and Hélène both nodded.

The bridge stayed quiet for another beat. Then Klaus said:

"Thank you for putting it that way. I wanted it on the record that the chain is not magic. The 2024 incident is something Hélène and I have lived with. I do not want anyone reading the deliverable to think this audit closes the 2024 question. It does not. It changes the cost of the *next* one."

Dawn: "That is the right framing. We will use it."

---

### 🔍 4:30 PM CET — The BMW Question

Klaus brought up the question that had been hanging over the engagement since kickoff. By video bridge.

"Stefan Kuhn at BMW asked us last month — and I am going to ask Dawn the same question. If a BMW customer's vehicle has an AI-driven false-positive predictive-maintenance alert that costs the customer time and BMW money, can BMW determine whether the false positive came from a regression in Lumière's model or from Eberhardt's integration or from BMW's own vehicle-side data?"

Dawn had been ready for the question. She had been thinking about it since lunch.

"The §10.21 cross-vendor anchor does that. Three chains, one root-cause path."

She walked it through.

"The BMW vehicle-side data — what the sensor reported, what the BMW-side ingestion logged — that is BMW's chain or BMW's equivalent of a chain. They have it on their side. It is the input.

"The Eberhardt inference chain entry shows the input telemetry hash that Eberhardt received from the BMW vehicle. If BMW's input does not match the input Eberhardt processed, the difference is in the BMW-to-Eberhardt transport — that is BMW's integration question. If they match, the input is reconciled.

"The Eberhardt inference entry shows the model artifact hash. The §10.21 model-handover entry shows the handover record. The Lumière model-build entry shows the model build provenance. If the model artifact at inference matches the model-build entry, the model is reconciled per the §10.21 cross-anchor (independent verifier output on each side plus byte-equal hash join). If the model artifact differs from any prior good build, the regression is in Lumière's build. That is Lumière's responsibility.

"The Eberhardt inference entry shows the inference output. If the output looks anomalous given the input — and the input and the model are both reconciled — the regression is in the model's behavior on this input class. That is again Lumière's responsibility, but in the modeling-quality sense, not the build-integrity sense (the §1.2 epistemic-scope split — the chain proves what was deployed; the modeling-quality is a separate evidence regime).

"The Eberhardt-to-BMW transport — what Eberhardt sent to BMW after inference — is the next leg. Eberhardt's outbound logging shows what was sent. BMW's inbound logging shows what was received. If they match, the transport is reconciled. If they don't, the integration question is on the boundary between Eberhardt and BMW.

"Three chains. One root-cause path. Each chain answers its own leg. The §10.21 cross-vendor anchor is what makes the join possible without trust in a single party's claim. That is exactly the joint-supplier audit BMW asked for."

Klaus had been listening with his fingers on his chin.

"That is the answer."

Hélène, on the screen: "That is the answer. And the auto-industry context matters here — the chain composes onto the M&A scenarios §10.24 entity succession addresses if BMW or Stellantis were ever to acquire Eberhardt, or if Lumière were acquired by a US AI lab. The §10.24 `chain.entity_succession` event with the dual-signature shape (per §10.17 signatory schema with `entity_affiliation`) means the chain stays continuous across a legal-entity change. That is a posture our boards understand."

> **✓ Confirmation — Three-chain root-cause path is end-to-end traversable for false-positive analysis per §10.21 + §1.2 + §1.4**
> A BMW false-positive predictive-maintenance alert can be root-caused across three chains — BMW's vehicle-side data, Eberhardt's inference chain, and Lumière's model-build chain — by reconciling input-telemetry-hash, model-artifact-hash, and inference-output-hash at the §10.21 cross-vendor anchors. Each chain answers its own leg without requiring trust in another party's claim — the §1.4 compositional-security argument extended across an organizational boundary. §1.2 epistemic-scope discipline distinguishes integrity claims (the chain) from modeling-quality claims (the fairness-audit + validator + field-monitoring). This is the BMW joint-supplier audit deliverable.

Tom: "I want to capture the answer verbatim. Stefan Kuhn will read it. He will know we addressed his question explicitly."

Dawn: "Capture it."

Klaus turned the question one more time. "Dawn — what does the chain say about the modeling-quality side? When BMW reads our deliverable they will ask not only 'where is the regression' but 'how does the chain support the validator?' I want to be honest with them about what the chain does and does not do."

Dawn took a beat. "Per §1.2 epistemic scope, the chain proves what was deployed and that the record was not tampered with. It does not prove the model is correct, that the model complies with policy, or that the model is bias-free. Those are three separate evidence regimes. The fairness-audit report — anchored by `audit.model_handover.fairness_audit_report_sha256` per §10.21 — is the bias claim. The validator harness on Lumière's side, and the field-monitoring cadence on Eberhardt's side, are the correctness claims. The chain proves which model was deployed when and which fairness-audit accompanied it. It does not adjudicate whether the model was good. Stefan will appreciate the distinction. EU AI Act conformity-assessment files specifically need this language — the chain is the integrity foundation; the bias-and-correctness claims live in separate documents that the chain hash-anchors but does not evaluate."

Klaus: "That is the right framing. The chain is the foundation. The validator and the fairness audit are separate, anchored to the foundation by hash."

Hélène: "And per §10.21 the fairness-audit anchor is bidirectional. Eberhardt reads the fairness-audit hash from the §10.21 model-handover entry. Lumière holds the report under retention. If the report changes, the hash changes; the chain entry does not. The post-hoc edit is detectable. That is the §10.21 cross-anchor verification mechanism applied to the fairness-audit boundary, not just the model-artifact boundary."

Tom: "Logged. The §1.2 distinction is going into the executive summary verbatim. The chain does not validate; it integrity-binds. The validator and the fairness audit are separate."

---

### 🌆 5:30 PM CET — Joint Debrief by Video Bridge

The full team was on the bridge at 5:30 PM. Stuttgart had moved everyone into the conference room. Paris had pulled chairs up around Hélène's desk. Klaus and Hélène were both on. Tom and Dawn led.

Dawn ran through the per-company summary first.

"Eberhardt-side individual findings: zero Gaps, zero Partials. The chain holds at eight months. IAM is two-AD-domain on-prem federated to Azure AD, fully chain-coupled per §3 and §10.2 operational events. Seal job is chain-governed with HSM at Sindelfingen per §4.2 / §4.3 / §10.5. §10.17 partition-ceremony events are present and well-formed. mTLS rotation to BMW is chain-coupled with downstream OEM per §10.10 rotation-crossing-the-seal-boundary. §10.15 Pattern A multi-region resilience operates correctly with synchronous-read `master.cross_region_replication_completed` events per Round-17 third errata. BSI IT-Grundschutz alignment, ISO 27001, ISO/IEC 27017, TISAX assessment artifacts all feed from the chain.

"Lumière-side individual findings: zero Gaps, zero Partials. The chain holds at four months. IAM is Google Workspace plus custom SSO with author-approver schema enforcement. Seal job is chain-governed with HSM at OVHcloud Paris. §10.17 partition-ceremony events are present and well-formed with the §10.17 cross-language CC8.1 discoverability discipline applied (English CC8.1 cross-references the French runbook by title and named sections per §10.18). ANSSI-aligned configuration. Fairness-audit hash anchoring works.

"Joint findings — the load-bearing part of this engagement, and the engagement that drove §10.20 + §10.21 into the spec body. Zero Gaps. One Partial. Two Nits.

"**Finding #1 (Partial, Joint): Lumière's 90-day training-data retention violates §10.20 (training-data retention vs deployment-window discipline).** Hélène has accepted the finding and committed to a 24-month retention extension under GDPR Article 6(1)(f) legitimate interest with documented justification per the §10.20 GDPR data-minimization tension resolution mechanism. The next §10.21 model-handover entry will carry `audit.model_handover.training_data_retention_floor_days = 720`. The retention extension is in flight as a CAPA against this audit. **Spec §10.20 (Wave-6 Eberhardt-Lumière fourth errata) is the close-out section; this engagement drove its addition to the spec body.**

"**Finding #2 (Nit, Joint): `audit.cross_border_transfer.*` attribute family is implicit in chain entries** (Stuttgart-Paris training-data transfers under Article 28 processor agreement). Within-EU transfers do not need the attribute set under GDPR — same regime end-to-end. Per §10.21's cross-border-composition recommendation (which **names the Eberhardt-Lumière case as the within-EU exemplar**), the explicit attribute set is examiner-friendlier for non-EU regulators and BMW's vendor-management read. Hélène commits to emitting the explicit set within four weeks. Schema is already normative; institution-side emission is the work.

"**Finding #3 (Nit, Joint): Lumière fairness-audit entry uses singular `audit_report_language`; §10.21 normates plural-array `audit.model_handover.audit_report_languages` (Round-17 M&A-N2 close-out).** The §10.21 model-handover entry on the Eberhardt side already uses the plural-array form correctly. The Lumière-side fairness-audit attribute is a Lumière-internal naming that predates the §10.21 normative shape. Hélène commits to producing the v4.2.1 English translation in the next four weeks and aligning the Lumière fairness-audit entry to the plural-array form on the next schema revision. The BMW reader-matrix discoverability is already flowing through the §10.21 model-handover plural form on Eberhardt's side. **§10.21 (Wave-6 Eberhardt-Lumière fourth errata) named this case in its plural-array worked example; this engagement drove the spec elevation.**

"Joint posture: the §10.21 cross-vendor anchor mechanism is the engagement's load-bearing finding. Two chains compose at the model-handover boundary. End-to-end reconciliation traverses seven legs from BMW vehicle through Eberhardt inference to Lumière training-data manifest, with byte-equal hash matches at every cross-vendor anchor (§5 RFC 8785 JCS canonicalization plus SHA-256), in under fifteen minutes elapsed."

Klaus: "Read out the regulator-by-regulator posture."

Dawn: "Five readers.

"**EU AI Act.** Article 11 deployer logging: satisfied by Eberhardt's inference chain. Article 12 conformity assessment: satisfied by Eberhardt's chain plus Lumière's fairness-audit chain plus the §10.21 cross-vendor anchor. Article 16 provider obligations: satisfied by Lumière's fairness-audit and model-build chains. The §10.20 retention floor extension closes the long-tail forensic gap that EU AI Act Article 12 logging contemplates for post-market surveillance. The joint chain is a defensible Article 11/12/16 evidence pack — and the §1.2 epistemic-scope discipline gives the conformity-assessment file the right vocabulary for what the chain proves vs. what other regimes prove.

"**BSI (Bundesamt für Sicherheit in der Informationstechnik).** IT-Grundschutz baseline: satisfied. The seal infrastructure (§4.2 / §4.3), HSM configuration (§10.5 / §10.17), log retention (§10.13), identity-provenance walk (§3 / §10.1) are all aligned. The ISO 27001 and ISO/IEC 27017 certifications provide the formal layer.

"**TISAX.** Auto-supply community assessment: the chain artifacts feed the TISAX assessment. Coordinated supplier security assessment is consistent with TISAX practice. The §10.21 cross-vendor anchor pattern extends the assessment across the Eberhardt-Lumière boundary cleanly.

"**LfDI Baden-Württemberg and CNIL (joint GDPR read).** Article 28 processor agreement governs Eberhardt-to-Lumière training-data transfer. Within-EU transfer documented per §4.4.1. Cross-border-transfer attribute is implicit (Finding #2 Nit). Article 6(1)(f) legitimate interest as the consent basis with DPIA hashed into the training-data manifest. Article 5(1)(c) data-minimization is the driver for the original 90-day retention; the §10.20 retention extension to 24 months will be re-justified under the §10.20 GDPR data-minimization-tension-resolution mechanism — Article 6(1)(f) legitimate-interests determination tied to EU AI Act Article 12 logging obligations, with Article 35 DPIA naming the longest-deployment-window justification.

"**BMW vendor-management (Stefan Kuhn).** §10.21 cross-vendor anchor composes end-to-end without trust in either party's claim — §10.21 cross-anchor verification mechanism (independent verifier output on each side plus byte-equal hash join). Three-chain root-cause path is traversable for false-positive analysis. The 2024 model-drift incident class would now resolve in minutes instead of days for detection, and within the §10.20 24-month retention window for retroactive root-cause analysis. Joint-supplier audit deliverable is satisfied. **§10.21 names the BMW joint-supplier audit case as its exemplar in the plural-array worked example and the cross-border-composition note.**

"**ISO 26262 functional-safety auditor.** The chain artifacts support the audit-trail requirements for AI-derived safety-related outputs. The chain is independent of ISO 26262 but compatible. ISO/SAE 21434 (automotive cybersecurity) and UNECE WP.29 R155/R156 (cybersecurity / software updates) compose alongside.

"That is the regulator matrix."

Klaus took a long breath in.

"Hélène. We are good?"

Hélène, on the screen: "We are good. The §10.20 retention extension is the only structural change. The two Nits are §10.21 institution-side emission alignment and translation work. The chain holds across the seam."

Klaus turned to the camera. "Dawn — Tom — the Stuttgart and Paris teams. Thank you. This is what I had hoped for. The deliverable for BMW will reference this audit explicitly."

Dawn: "Thank you. The formal deliverable will be in your inbox by end of next week. The joint posture, the per-regulator matrix, the three findings classified to the post-amendment spec, the live-trace evidence, and the language for Stefan Kuhn's three questions — all of it will be documented per §10.13 evidentiary-artifacts retention discipline. The deliverable will reference Appendix A consolidated chain envelope schema for the attribute-family lookup the BMW vendor-management read will want, and §13 stakeholder navigation will guide the BSI / LfDI / CNIL / BMW / ISO 26262 reader paths through the deliverable. The §11 references section pins the verifier version both sides use per §10.26."

Hélène: "And from Lumière — thank you. The 90-day retention conversation is one I should have had with myself months ago. Having it forced by the audit is the right reason to have it now. And — for the record — the spec section that names our retention asymmetry as its worked example is a useful place for our policy to land. The §10.20 retention floor will be a permanent line in our model-supply DPA going forward."

Dawn: "That is what audits are for. And — the spec is what audits land in. The §10.20 amendment is the durable answer to the question your audit raised."

Dawn caught Raj's eye on the screen as the meeting was winding down. "Same coffee debt as Northbridge?"

Raj, deadpan: "Ten engagements ago. Statute of limitations expired."

"Convenient."

The bridge held for another minute while everyone exchanged the small post-meeting pleasantries — handshakes by camera, the small wave from Klaus to the Paris side, a thumbs-up from Aurélien to Mike. Then the bridge dropped.

The Stuttgart team packed up. Maximilian stayed for a moment to trade contact details with Mike. Andreas stayed to walk Diana through one final question about the AD-to-Azure-AD federation that he wanted to clarify off the record. Bettina caught Luis in the corridor and gave him a recommendation about a Splunk admin-operation chain-coupling pattern that they had been considering for the next quarter. Tom shook Klaus's hand a second time at the door.

In Paris, Sophie walked Elena to the lobby. Aurélien stayed to chat with Chen and Raj about the §4.4.1 `audit.cross_border_transfer.*` emission alignment — Chen was already drafting it on his laptop. Hélène walked them all to the door at the bottom of the Haussmann staircase and said goodbye in three languages depending on whom she was speaking to.

---

## Final ✅ Comparison Block

### Per-Company Posture

| Dimension | Eberhardt (Stuttgart) | Lumière (Paris) | Joint (§10.21 Cross-Vendor Anchor) |
|---|---|---|---|
| Chain duration | 8 months | 4 months | 4 months at handover |
| HSM | Thales Luna on-prem, Sindelfingen | OVHcloud Paris-region, ANSSI-aligned | Independent IKMs per §4.1 |
| Verifier | `herald-verify`, 4 sec, 12 steps per §7, exit code 0 per §10.12 | `herald-verify`, 4 sec, 12 steps per §7 | byte-equal hash match at anchor per §10.21 + §5 |
| Service-account IAM | chain-coupled per §10.2 | chain-coupled per §10.2 | n/a |
| Author-approver separation | n/a | schema-enforced (SDK refusal pattern per §4.4) | n/a |
| §10.15 multi-region resilience | Pattern A (Sindelfingen seal, Munich replication) | Single-region | n/a |
| §10.17 partition-ceremony coupling | present, with §10.17 + §10.18 cross-referencing | present, with §10.17 cross-language CC8.1 discoverability | n/a |
| §10.19 chain-coverage map | published, version-stamped per Round-17 M&A-P3 | published, Pipedrive scoped out | n/a |
| Compliance baselines | BSI IT-Grundschutz, ISO 27001, ISO/IEC 27017, TISAX | ANSSI-aligned, GDPR | Article 28 processor agreement |
| Functional-safety alignment | ISO 26262 trail-supporting; ISO/SAE 21434 + WP.29 R155/R156 compatible | n/a | n/a |
| Individual gaps | 0 | 0 | n/a |
| Individual partials | 0 | 0 | n/a |
| Joint partials | n/a | n/a | 1 (§10.20 retention floor) |
| Joint nits | n/a | n/a | 2 (§4.4.1 cross-border emission, §10.21 plural-array alignment) |

### Per-Regulator Posture

| Reader | Posture | Rationale |
|---|---|---|
| EU AI Act (Article 11 / 12 / 16) | Satisfied | Joint chain is a defensible deployer-and-provider evidence pack; §10.20 retention floor closes Article 12 post-market surveillance forensic gap |
| BSI IT-Grundschutz | Satisfied | Seal infrastructure (§4.2 / §4.3), HSM (§10.5 / §10.17), identity-provenance (§3 / §10.1) aligned |
| ISO 27001 / ISO/IEC 27017 | Satisfied | Formal certification layer over the chain |
| TISAX | Satisfied | Chain artifacts feed the auto-supply community assessment |
| LfDI Baden-Württemberg | Satisfied with Nit | §4.4.1 cross-border-transfer attribute implicit, not blocking |
| CNIL | Satisfied with Nit | Same as LfDI; within-EU regime is end-to-end consistent |
| BMW vendor-management (Stefan Kuhn) | Satisfied | §10.21 cross-vendor anchor + three-chain root-cause path |
| ISO 26262 functional safety | Trail-supporting | Independent of TesseraSeal but compatible; ISO/SAE 21434, WP.29 R155/R156 compose alongside |

### The Three Findings (post-amendment classification)

| # | Severity | Spec Section | Finding | Status |
|---|---|---|---|---|
| 1 | Partial (Joint) | §10.20 (training-data retention vs deployment-window discipline) — **this engagement drove its addition to the spec body** | Lumière's 90-day training-data retention is shorter than Eberhardt's typical 9-18-month deployment window. §10.20 worked-example match. | Hélène committed to 24-month retention floor under GDPR Article 6(1)(f) legitimate-interests determination tied to EU AI Act Article 12 (the §10.20 GDPR-tension resolution mechanism). Next §10.21 model-handover entry to carry `audit.model_handover.training_data_retention_floor_days = 720`. CAPA in flight. |
| 2 | Nit (Joint) | §4.4.1 `audit.cross_border_transfer.*` family + §10.21 cross-border-composition recommendation — **§10.21 names this case as its within-EU exemplar** | `audit.cross_border_transfer.*` attribute set is implicit in Stuttgart-Paris training-data transfer entries. Within-EU SHOULD per §10.21. | Hélène commits to emitting the explicit §4.4.1 attribute set within four weeks for examiner-friendliness (BMW German-English-French reader matrix). Schema is already normative; institution-side emission is the work. |
| 3 | Nit (Joint) | §10.21 `audit.model_handover.audit_report_languages` plural-array (Round-17 M&A-N2 close-out) — **§10.21 names this case in its plural-array worked example** | Lumière fairness-audit entry uses singular `audit_report_language`; §10.21 normates plural-array. The §10.21 model-handover entry on Eberhardt's side already uses the plural form correctly. | Hélène committed to v4.2.1 English translation in 4 weeks; bilingual-or-trilingual audit bodies as standing practice for new models; Lumière fairness-audit entry alignment to plural-array form on next schema revision. BMW reader-matrix discoverability already flowing through §10.21 model-handover plural form on Eberhardt's side. |

---

### 🧾 Final Assessment Theme

The TGV from the Stuttgart side back across the border was tomorrow's plan. Tonight the team would all stay put — Dawn at the Stuttgart hotel, Raj at his Paris hotel near Bastille — and write up the deliverable on shared docs. The drive back to the hotel from Eberhardt was twenty minutes through the Stuttgart hills as the sun was going down. Dawn had her coffee, refilled, in the cup holder.

She thought about the day. About Maximilian's terminal at 9:15 in the morning showing a clean §7 PASS on the inference entry. About Aurélien's terminal at 10:00 in Paris showing a clean PASS on the model-build entry. About Chen at noon on the bridge running the byte-equal compare and the §10.21 cross-anchor hashes matching across both chains. About Klaus's question at 3:45 — *does the chain prevent another 2024?* — and her honest §1.2-grounded answer that it does not prevent, it detects faster, and faster detection is materially better than two-week silent regression. About Hélène's quiet acknowledgment of the §10.20 retention asymmetry and her unhesitating commitment to extend. About Klaus's question at 4:30 — *can BMW root-cause a false positive across the seam?* — and the §10.21 three-chain answer that took her ninety seconds to walk and that Stefan Kuhn would read in the deliverable.

She thought about the spec sections that the engagement had ended up reading itself through. §10.20 — the retention floor, with the GDPR Article 5(1)(c) data-minimization tension resolved through Article 6(1)(f) legitimate-interests determination tied to EU AI Act Article 12 logging obligations. §10.21 — the cross-vendor model-handover schema, with the contract triple closing Round-17 M&A-G2 and the plural-array `audit_report_languages` closing Round-17 M&A-N2 and the cross-border-composition note naming the Eberhardt-Lumière case as the within-EU exemplar. §4.4.1 — the cross-border-transfer attribute family that Sun-Won had driven into the spec last quarter. §10.16 — the SaaS-edge connector lag bounds that Northbridge had driven before that. §10.17 — the HSM partition ceremony attestation that NetiVa had driven. §10.15 Pattern A — the multi-region replication freshness clarification that Atrio had driven. §10.19 — the chain-coverage map plus `audit.external_artifact.*` that Salt Pond had driven. §10.18 — the runbook cross-referencing rule that landed alongside Atrio's freshness clarification.

Each engagement had landed one or two normative additions. The pattern was the same — the field finds a gap, the spec section is added, the next engagement of the same shape reads the spec and remediates against the named section rather than rediscovering the gap. The body of the spec was the accumulated answer to every audit-day finding the team had walked. The §12 change log read like a roster of the engagements.

She thought about Northbridge — the gold standard, fully sealed, one Nit reclassified to non-conformance under §10.16's "imprecise lag wording is non-conformance, not cosmetic" rule. About Mercator — the bifurcation, AI on the chain, EHR off, Patricia Okonkwo's funding-roadmap framing. About Stelvio — the three-zone, AI sealed, OT mutable, IT business legacy, Maria Costanza's triage. About Atrio — the multi-tenant, forty-seven tenants under twelve sponsor-bank IKMs, fourteen hundred verifier runs, the §10.15 Pattern A invariant 5 cache-staleness Partial that drove the §10.15 freshness clarification. About Helmstad — the biopharma CRO boundary, Quintessa's PGP-signed SFTP, source-side history outside the chain, FDA Part 11 alignment. About Pacific Crescent — the utility, AI gas-pipeline leak detection, Brentwood real-leak, NERC CIP composition. About Olmstead — the university, two override-down decisions where the rationale was gone, the FERPA-and-GLBA boundary. About NetiVa — the Tel Aviv vendor-management evaluation that drove §10.17 partition-ceremony attestation. About Sun-Won — the Korean engagement, PIPA Section 28 cross-border-transfer attribute that drove the §4.4.1 family. About Salt Pond — the toy supply chain, §10.19 chain-coverage map and `audit.external_artifact.*`.

Today was different. Today was the first joint engagement. Two companies, two countries, one chain that crossed the boundary. The §10.21 cross-vendor anchor was the load-bearing thing — a more rigorous version of what Helmstad's CRO had achieved with PGP-signed PDFs, scaled to mutual chain-on-both-sides verification with byte-equal hash matches at the seam. Not a SOC 2 attestation. Not a contractual representation. A cryptographic compose — the §1.4 compositional-security argument extended across an organizational boundary.

The auto-industry reader matrix made the engagement specific in a way the prior engagements had not been. ISO/SAE 21434 (automotive cybersecurity) and UNECE WP.29 R155 / R156 (cybersecurity / software updates) were composing alongside the chain — the chain artifacts feed both regimes without contradicting them. ASPICE (Automotive SPICE process assessment) is consistent with the chain's evidence shape. ISO 26262 functional-safety auditors accept the chain as part of the audit-trail for AI-derived safety-related outputs. TISAX is the auto-supply community's coordinated security assessment, and the chain artifacts feed the TISAX assessment cleanly. TÜV / DEKRA / SGS audit reports — the third-party assessment ecosystem German automotive lives in — read the chain as evidence of integrity-bound deployment posture. The chain is not an automotive-specific regime; it is the integrity foundation that the automotive regimes compose onto without re-engineering.

That had been Klaus's framing from the start. *Build it once, build it properly, trust it for thirty years.* The chain at Eberhardt is a Mittelstand artifact in the way Klaus thought of it — a thing that does not change capriciously, that the next generation of engineers can read and verify, that the next generation of OEM customers can audit without re-engineering. The §12 change log read like a record of careful changes — Wave 1 through Wave 6, Round-17 close-outs, the v1.0a wire form locked, the v1.0b 12-line `sign_payload` extension, the §10.20 + §10.21 errata that this engagement drove. The pattern is the spec growing through engagement findings, not through speculation. That is the discipline Klaus respects, and the discipline the chain delivers.

The chain at Eberhardt held. The chain at Lumière held. The seam held.

And the seam — *the seam was the engagement*. Each side's chain alone was easy to confirm. The interesting work was always going to be the join. The 90-day retention asymmetry that drove §10.20. The audit-report-languages matrix that drove the §10.21 plural-array shape. The implicit cross-border-transfer attribute that §10.21's cross-border-composition note now names as the within-EU exemplar. None of those were chain-integrity findings. They were findings about how two chains, written by two teams in two countries under two compliance baselines, talk to each other through a normative anchor mechanism that has to be intelligible to a third reader — BMW's Stefan Kuhn — who has to root-cause a false positive without trusting either party's claim alone. The spec amendments — §10.20 and §10.21 — are now the durable answer to those findings. The next institution that walks an Eberhardt-Lumière-shaped engagement reads the spec and remediates against the named sections rather than discovering the gap in the field.

The 2024 model-drift incident was the lesson. Pre-chain, three days of reconstruction. Post-chain, three minutes of detection. The §10.20 24-month retention extension closes the long-tail root-cause window. The chain does not prevent the next 2024. It changes the cost of finding out.

Dawn picked up her phone at a red light on the way down out of the hills and dictated a one-line note for the report's executive summary.

> *"The chain at Eberhardt proves what was deployed. The chain at Lumière proves what was built. The §10.21 cross-vendor anchor proves they are the same model — and that is what the joint-supplier audit asked for."*

The light turned green. She put the phone down and drove the rest of the way to the hotel as the Stuttgart sky went dark over the Schlossplatz.

---
