# 🧾 Diary of an Audit Day — Eberhardt Werkstoffe × Lumière AI

**Engagement:** Joint pre-audit ahead of EU AI Act enforcement and the BMW joint-supplier audit
**Clients (joint):**
- **Eberhardt Werkstoffe GmbH** — German automotive-electronics Mittelstand, ~€1.4B revenue, ~3,800 employees, Stuttgart HQ, Sindelfingen primary plant, third-generation family-owned. Sensors, ECUs, AI-driven predictive-maintenance modules for OEM customers (BMW, VW, Mercedes-Benz, Stellantis, Renault). Recent push into EV-specific modules (battery health, motor-controller predictive maintenance).
- **Lumière AI Sàrl** — French AI consultancy, Paris HQ in the 1er arrondissement, ~80 employees, ex-INRIA founders (Institut national de recherche en informatique et en automatique), strong fairness-and-explainability practice. Builds custom ML models for industrial customers including Eberhardt.
**Posture:** TesseraSeal in production at both companies. **The chain extends across the partnership boundary.**
- Eberhardt: 8 months on the OEM-facing inference path (every predictive-maintenance alert is sealed in Eberhardt's chain).
- Lumière: 4 months on the model-development pipeline (training data hashes, model artifact hashes, fairness-audit reports — sealed in Lumière's chain).
- The handoff: when Lumière delivers a new model, Lumière's model-artifact-hash + model-card-hash + fairness-audit-report-hash are recorded into an Eberhardt chain entry of `service.name = "model-handover"` as cross-vendor anchors. Two chains compose at the handover event.
**Date:** Tuesday, the week after the engagement-before-this-one
**Audit team lead:** Karen
**Stuttgart liaison:** Klaus Eberhardt, third-generation owner-CEO. 64. Engineering background — Maschinenbau from RWTH Aachen.
**Paris liaison:** Hélène Lefebvre, COO. 47. Ex-INRIA researcher, PhD in computational learning theory.
**Cities:** Stuttgart and Paris, both CET. About one hour apart by TGV-Eurostar. The team works in parallel and joins by video bridge midday and at debrief.

---

## Context

Eberhardt Werkstoffe builds the boxes the OEMs hide behind the dashboard. Sensors, ECUs (electronic control units), and the AI-driven predictive-maintenance modules that look at battery cells and motor controllers and decide whether the car needs to send the driver a warning or schedule a service appointment. Family business, third generation. Klaus Eberhardt's grandfather started it after the war making relays for Daimler-Benz. Klaus runs it now from a glass-fronted office in a Stuttgart industrial park that overlooks the Sindelfingen plant in the distance. The family still owns the company outright. The Mittelstand model — the German mid-size family-business sector — is a culture as much as an economic category. Klaus thinks of the chain the way his father thought of the relay-test bench: a thing you build once, build properly, and trust for thirty years.

Lumière AI is a different shape entirely. Eighty people in a renovated Haussmann building near the Louvre, founded six years ago by two ex-INRIA researchers who decided industrial ML had a fairness problem nobody was solving in production. They publish papers. They speak at NeurIPS. They also ship models that go into cars. Their customer list is short and exclusive — Eberhardt is the largest. Hélène Lefebvre runs the operations side; the founders run research. The two cultures meet at the model-handover boundary, which is where the chain meets the chain.

Eight months ago Eberhardt deployed TesseraSeal across the OEM-facing inference path. Every predictive-maintenance alert that goes from an Eberhardt module out to a BMW or VW vehicle is sealed in Eberhardt's chain. The hash of the model artifact that produced the alert is in the chain entry. The hash of the model card is in the chain entry. The vehicle ID is hashed into the chain entry. The OEM customer's deployment-side data is what BMW or VW logs on their side. Three chains, conceptually — Eberhardt's chain, Lumière's chain, the OEM's vehicle-side data — and the cross-vendor anchor is the bridge between the first two.

Four months ago Lumière deployed TesseraSeal on the model-development pipeline. Training data hashes, model artifact hashes, fairness-audit reports — all sealed in Lumière's chain under their own IKM. The deployment was driven by a 2024 incident. A battery-health prediction model regressed silently after a training-data update. Lumière delivered the regressed model to Eberhardt. Eberhardt deployed it to BMW. Two weeks of degraded predictions before someone noticed. Lumière did not have the chain at the time. They could not retrace which deployments had received the regressed model. Hélène had been the one who put the case together for the founders that the chain was no longer optional.

The handoff is the load-bearing thing. When Lumière delivers a new model to Eberhardt, an Eberhardt chain entry of `service.name = "model-handover"` records a JCS-canonical attribute set:

```
{
  "vendor": "lumiere-ai",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "abc123...",
  "model_card_sha256": "def456...",
  "fairness_audit_report_sha256": "ghi789...",
  "lumiere_chain_entry_id": "lum-2026-04-08-..."
}
```

The hash of Lumière's chain entry is the cross-vendor anchor. Verifying Eberhardt's chain confirms what Eberhardt received. Verifying Lumière's chain (with credentials Lumière grants on a request basis) confirms what Lumière handed over. Both must match for end-to-end verification. That mechanism is the centre of gravity of this engagement.

Klaus and Hélène commissioned the joint audit because three things had converged. The EU AI Act enforcement deadlines for high-risk AI systems were approaching — predictive-maintenance for safety-critical automotive systems sits inside Annex III's high-risk category, and Lumière as the model provider had Article 16 obligations while Eberhardt as the deployer had Article 11 logging and Article 12 conformity-assessment obligations. BMW's vendor-management team had asked Eberhardt and Lumière for a joint-supplier audit that demonstrated end-to-end chain coverage at the OEM-supplier boundary. And the 2024 model-drift incident was in the rear view but not far enough that BMW had stopped asking about it.

The deliverable will be read by both companies' boards. By BMW's vendor-management team. By the German BSI (Bundesamt für Sicherheit in der Informationstechnik) for cybersecurity sign-off. And, for the cross-border GDPR aspect, by both the German LfDI Baden-Württemberg and the French CNIL. Five readers. Two countries. One report.

---

## Audit Team

The same eight-person team that walked Northbridge a quarter ago and Olmstead and the others. This week they split.

### Stuttgart (Eberhardt HQ)

- **Karen** — Lead Auditor (governance + narrative)
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

Karen had landed at Stuttgart Airport the night before and taken the S-Bahn into Mitte. The hotel was a small one near the Schlossplatz. The drive out to the Eberhardt office in the morning was twenty-five minutes through the Stuttgart hills. She had her coffee in the cup holder and the engagement brief on her tablet.

*Two countries*, she thought as she came up the autobahn ramp. *Two companies. One chain that crosses the boundary. EU AI Act for both. BMW watching. 2024 model-drift in the rear view. We test whether the cross-vendor anchor actually composes — what one chain says, the other confirms.*

She took stock of the prior engagements. Northbridge had been the gold standard — full single-tenant deployment, eighteen months mature, the team had spent four days trying to find a gap and found a stale comment in a YAML file. Mercator had been the bifurcation — sepsis CDS on the chain, the EHR off the chain, Patricia Okonkwo's funding-roadmap framing. Stelvio had been the three-zone version — AI side sealed, OT mutable, IT business legacy, Maria Costanza's triage. Atrio had been the multi-tenant test — forty-seven tenants under twelve sponsor-bank IKMs, fourteen hundred verifier runs, zero failures, Naomi Reisinger's coordinated examiner room.

Then Helmstad. The biopharma. The CRO data feed where Quintessa had PGP-signed the SFTP delivery and Helmstad had recorded the SHA-256 at the boundary, and the source side beyond the boundary lived on Quintessa's SOC 2. Pacific Crescent — the utility, AI gas-pipeline leak detection on the chain, OT historian off the chain, the Brentwood alert that turned out to be a real small leak. Olmstead — the university, AI admissions screening on the chain, Slate free-text rationale-fields off the chain, two override-down decisions where the rationale was gone.

After Olmstead there had been two more. The Korean engagement at Sun-Won where the PIPA Section 28 cross-border-transfer attribute had been the central finding — the chain had not carried an explicit transfer-basis flag and the PIPC examiner had pointed that out specifically. And the engagement after that. Nine prior engagements before today. Today is the tenth-and-eleventh — joint, parallel, two countries.

*It never is*, Karen thought. *Today the question is whether the joint chain holds at the handoff. Each side's chain looks fine alone. The question is the seam.*

The Eberhardt building was in a Stuttgart industrial park, glass front, the Eberhardt name in plain stainless steel on the wall. Klaus Eberhardt was waiting in the lobby in a navy jumper and grey trousers. He was 64, lean, with engineer's hands and a directness Karen recognized in the first thirty seconds.

"Karen. Welcome to Stuttgart." His English was almost unaccented. "You came in last night?"

"I did. Hotel near the Schlossplatz."

"Good. The airport hotel is a tragedy." He shook hands with each of the Stuttgart team in turn — Mike, Diana, Luis, Tom — as they came through the revolving door from the parking lot.

The conference room was on the second floor with a view of a lawn that ran down to a pond. There was a long table with a Polycom video unit at the head and a screen on the far wall. Hélène Lefebvre was already on the screen from Paris, in a small office with a window behind her that showed grey Parisian morning. The Paris team — Raj, Elena, Chen — were visible in a corner of the screen, just arrived themselves.

Klaus did not waste time.

"Good morning. Let me say what I want from this. Eberhardt has had the chain for eight months on the OEM-facing inference path. I am reasonably confident in it. Lumière deployed four months ago on the model-development side. Hélène is reasonably confident in hers. Neither of us has tested the seam in front of an audit team. The EU AI Act enforcement deadline is in seven months. BMW's vendor-management team has asked us for a joint-supplier audit that demonstrates the seam works. We have a 2024 incident in our recent history where Lumière delivered a regressed battery-health model and we both saw it late. I want today's audit to test the seam. I want to know — if a BMW customer's car has a false-positive predictive-maintenance alert tomorrow — whether we can together identify the root cause: was it Lumière's model, or our integration, or BMW's vehicle-side data."

Hélène, on the screen: "I will say the same thing from Paris. Lumière has obligations as a provider under Article 16 of the EU AI Act — fairness audit, transparency to deployers. The chain entries support those obligations on our side. Klaus's chain on the deployer side covers Article 11 logging and Article 12 conformity assessment. The cross-vendor anchor is what makes the two halves a whole. I want it tested today."

Karen put her coffee down on the table.

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

He typed again. The terminal showed a JSON entry — structured, with fields Mike recognized.

```
{
  "entry_id": "eb-bh-2026-04-08-bmwix-00874",
  "tenant": "eberhardt-battery-health",
  "service": "pred-maint-inference",
  "seq": 4827193,
  "ts": "2026-04-08T11:42:17.083Z",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "model_card_sha256": "sha256:8e1b82...a4f1",
  "vehicle_vin_hash": "sha256:a91f8b...ee27",
  "module_serial_hash": "sha256:c12e4f...8a91",
  "input_telemetry_hash": "sha256:d83a01...b3ef",
  "classification": "urgent-service-required",
  "confidence": 0.94,
  "downstream_oem": "bmw-ag",
  "model_handover_entry_ref": "eb-mh-2026-04-01-lum-00012",
  "prev_hmac": "sha256:...",
  "this_hmac": "sha256:...",
  "merkle_path": [...],
  "seal_ref": "eberhardt-battery-health-2026-04-08-eod"
}
```

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

Mike leaned back.

"Pick one from a hundred and eighty days ago."

Maximilian picked an entry from October 12, 2025. Same fleet, same model family — though running on `battery-health-v3.9.0`, the predecessor. Ran the verifier. Same four seconds. Same PASS. Same twelve steps.

> **✓ Confirmation #1 — Eberhardt chain integrity holds at eight months and 180 days**
> The predictive-maintenance service has been emitting sealed entries for 244 days. The verifier resolves a recent entry in 4 seconds and a 180-day-old entry in 4 seconds. Twelve verification steps including HMAC recomputation, Merkle path resolution against the daily-seal Merkle root, and Ed25519 signature verification against the published quarterly public key. The chain endures across the seal boundary and across model-version transitions.

Mike asked: "What is sealing the seal?"

"Daily Ed25519 signature on a Thales Luna network HSM at the Sindelfingen plant. On-prem. We are Mittelstand — for things this load-bearing we prefer on-prem. The IQ/OQ documentation for the HSM configuration is in our QMS. BSI IT-Grundschutz baseline — that is the German Federal Office for Information Security's compliance baseline. We are ISO 27001 certified, ISO/IEC 27017 certified for the cloud-touching parts, and TISAX certified for the auto-supply community — that is the Trusted Information Security Assessment Exchange the German automotive industry uses for coordinated supplier security assessments. The chain artifacts feed the TISAX assessment."

"Show me the IQ/OQ."

Maximilian pulled it up. 62 pages. Configuration captures, key-attribute definitions, partition isolation, key-ceremony minutes from August 2025 when the daily seal had been first stood up. Klaus's signature on page 62.

> **✓ Confirmation #2 — IQ/OQ for the on-prem HSM signing infrastructure is documented to BSI IT-Grundschutz and ISO 27001 standards**
> The Thales Luna network HSM at Sindelfingen has a 62-page IQ/OQ in the Eberhardt QMS. Key ceremony minutes are recorded with two-of-three approval (Klaus Eberhardt, the platform engineering lead, and the SecOps lead). The configuration is reproducible from the IQ/OQ. The artifacts feed into Eberhardt's TISAX assessment. ISO 26262 functional-safety auditors accept the chain artifacts as part of the audit trail for AI-derived safety-related outputs.

Mike asked: "Show me the model-handover binding."

Maximilian pulled up the entry referenced by `model_handover_entry_ref` in the inference entry — `eb-mh-2026-04-01-lum-00012`.

```
{
  "entry_id": "eb-mh-2026-04-01-lum-00012",
  "tenant": "eberhardt-battery-health",
  "service": "model-handover",
  "seq": 4801237,
  "ts": "2026-04-01T09:15:42.001Z",
  "vendor": "lumiere-ai",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "model_card_sha256": "sha256:8e1b82...a4f1",
  "fairness_audit_report_sha256": "sha256:9d3c72...e5f2",
  "lumiere_chain_entry_id": "lum-mb-2026-04-01-bh-00031",
  "received_by": "max.brenner@eberhardt.de",
  "approver": "klaus.eberhardt@eberhardt.de",
  "prev_hmac": "sha256:...",
  "this_hmac": "sha256:...",
  "seal_ref": "eberhardt-battery-health-2026-04-01-eod"
}
```

"Each model-handover entry references the Lumière chain entry ID," Maximilian said. "When we receive a model from Lumière, the engineering lead approves the artifact bundle, the chain entry is written, and the inference service starts using the new model only after the chain entry is sealed. The inference entry's `model_handover_entry_ref` field points back to that handover entry. So every inference is bound to the model-handover record that admitted the model."

> **✓ Confirmation #3 — Model artifact, model card, and fairness-audit-report hashes are bound to every Eberhardt inference via the model-handover entry**
> Each inference's `model_artifact_sha256` and `model_card_sha256` match the values in the referenced model-handover entry. Recomputing the SHA-256 of the production model artifact at any time matches the sealed value. The same is true for the model card and the fairness-audit report. The model-handover entry itself references the corresponding Lumière chain entry by ID — the cross-vendor anchor.

"And if a new model comes in from Lumière," Mike said, "the inference entries from that day forward reference the new model-handover entry."

"Yes."

"And the old model — if it stops being used — its inference history is still verifiable."

"Yes. Each entry references the model-handover that was in force at the time. Past entries do not change when a new model arrives."

Mike wrote that down. *Hash-bound provenance from inference back to model-handover, with cross-vendor anchor to Lumière.*

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
  "entry_id": "lum-mb-2026-04-01-bh-00031",
  "tenant": "lumiere-model-dev",
  "service": "model-build",
  "seq": 1872391,
  "ts": "2026-04-01T08:42:11.502Z",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "model_card_sha256": "sha256:8e1b82...a4f1",
  "training_data_manifest_sha256": "sha256:b73c92...11f4",
  "training_run_id": "tr-2026-03-28-bh-001",
  "fairness_audit_report_sha256": "sha256:9d3c72...e5f2",
  "fairness_audit_entry_ref": "lum-fa-2026-03-31-bh-00007",
  "handover_target": "eberhardt-werkstoffe",
  "approver": "aurelien.marchand@lumiere.ai",
  "prev_hmac": "sha256:...",
  "this_hmac": "sha256:...",
  "seal_ref": "lumiere-model-dev-2026-04-01-eod"
}
```

The model-artifact SHA-256 matched the value Stuttgart had pulled up an hour earlier — bit-for-bit. The model-card SHA-256 matched. The fairness-audit-report SHA-256 matched.

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
> The model-development pipeline has been emitting sealed entries for 124 days. The verifier resolves the build entry that anchors the model Eberhardt is currently running in 4 seconds. Twelve verification steps. The same shape as Eberhardt's chain, on a separate IKM under a separate HSM in a separate jurisdiction. Independent integrity statements that compose at the model-handover boundary.

"What is sealing your seal?"

"OVHcloud Paris-region HSM. We chose OVH for two reasons — French sovereignty over the key material and GDPR data-residency for the training data. The HSM configuration follows ANSSI recommendations — that is the French national cybersecurity agency, Agence nationale de la sécurité des systèmes d'information. We have an ANSSI-aligned configuration document and key-ceremony minutes from December 2025."

"Show me."

Aurélien pulled up the configuration document. 41 pages. Hélène's signature on the last page. The key ceremony had been four people — the two founders, Hélène, and the SecOps lead — with a two-of-four threshold for key-recovery operations.

> **✓ Confirmation #5 — Lumière HSM configuration is ANSSI-aligned and documented**
> OVHcloud Paris-region HSM with an ANSSI-recommended configuration. 41-page configuration document signed by Hélène Lefebvre. Key-ceremony minutes from December 2025 with two-of-four approval. French sovereignty and GDPR data-residency are the operational drivers. The ANSSI-aligned configuration is a discriminator for French regulators reading the deliverable.

Raj asked the next question.

"Walk me through the fairness-audit entry."

Aurélien pulled up `lum-fa-2026-03-31-bh-00007`.

```
{
  "entry_id": "lum-fa-2026-03-31-bh-00007",
  "tenant": "lumiere-fairness-audit",
  "service": "fairness-audit",
  "seq": 287413,
  "ts": "2026-03-31T16:12:48.119Z",
  "model_id": "battery-health-v4.2.1",
  "model_artifact_sha256": "sha256:7f2a91...c4d3",
  "audit_protocol_version": "lum-fa-protocol-v3.1",
  "audit_report_sha256": "sha256:9d3c72...e5f2",
  "audit_report_language": "fr-FR",
  "auditor_subject": "claire.dubois@lumiere.ai",
  "approver": "helene.lefebvre@lumiere.ai",
  "fairness_metrics_summary": {
    "disparate_impact_ratio_age_cohort": 0.97,
    "calibration_error_max": 0.013
  },
  "prev_hmac": "sha256:...",
  "this_hmac": "sha256:..."
}
```

Raj read the entry twice. The `audit_report_language` field caught his eye.

"`fr-FR`. The audit report is in French."

"Yes. The auditor wrote it in French. Lumière is a French company. The model card is bilingual — French and English — but the full fairness-audit report is in French only."

"And Eberhardt's deployment documentation is in German."

"Yes. The model-handover documentation that Eberhardt produces is in German — and bilingual where we cooperate, but the regulatory file at Eberhardt's side is German-language."

Raj wrote that down. *BMW examiner reads what — German, French, English? The chain captures the language attribute. The reader needs to know.*

> **✓ Confirmation #6 — Fairness-audit report hash is anchored in the chain with language metadata captured**
> Every fairness-audit chain entry records the SHA-256 of the audit report, the auditor's identity, the audit-protocol version, the approver, and the report language. The chain captures the language attribute as a normative field; the question is whether examiners reading the joint deliverable have access to the report in the language they read. That is a procedural question, not a chain-integrity question — but it is worth flagging.

Raj closed the fairness-audit window and opened the training-data manifest.

"Show me the training-data binding for v4.2.1."

Aurélien pulled it up. The manifest hash referenced a manifest document that listed every training-data shard, its source, its de-identification status, its consent basis under GDPR, and the SHA-256 of the shard. The manifest itself was the canonicalization of all of that. The manifest hash was sealed in the build entry.

"How long do you retain the training data shards themselves?"

Aurélien paused for a beat. "Ninety days post-model-delivery."

"Ninety days."

"Yes. After that the shards are deleted. The manifest hash and the manifest document are retained — those are part of the chain — but the underlying shards rotate out at ninety days post-delivery."

Raj wrote that down very carefully. *Ninety days post-model-delivery. The model is deployed at Eberhardt for as long as Eberhardt runs it. If a regression appears six months later and someone asks to retrace which training shards produced it — the shards are gone.*

"Hélène was the one who set the policy?"

"Yes. The policy is GDPR data-minimization-aligned. The argument was that we cannot retain training data indefinitely under GDPR Article 5(1)(e) storage-limitation principle, and ninety days post-delivery was the operational compromise."

"How long does Eberhardt typically run a model in production?"

Aurélien thought. "Battery-health v3.9.0 was in production for fourteen months before v4.2.1 replaced it. v4.2.1 has been live for thirty-eight days. The pattern is — depending on the OEM rollout cadence — between nine and eighteen months per major model version."

"So the deployment window is typically four to six times longer than the training-data retention window."

"Yes."

> **⚠️ Partial #1 — Lumière's 90-day training-data retention is shorter than Eberhardt's typical deployment window**
> Lumière retains training-data shards for 90 days post-model-delivery under a GDPR data-minimization policy. Eberhardt typically runs a deployed model in production for 9-18 months. If a model regression appears six months after deployment and someone wants to retrace which training shards produced the regression, the shards are gone. The manifest hash and the manifest document are retained — the chain proves *what manifest was used* — but the underlying shards have rotated out. This is the 2024-incident-class forensic gap: the chain detects the regression, but the regression's training-data root cause cannot be retrieved beyond 90 days. Hélène should consider extending retention to match the maximum deployment window, with appropriate GDPR justification.

Raj had not raised his voice. He had not even changed his expression. He simply wrote it down and went on.

"Show me the IAM for the model-build service."

---

### 🔐 11:00 AM CET — Diana on Eberhardt IAM (Stuttgart) and the Identity Walk in Paris

Diana started the IAM session in Stuttgart by sitting down with Eberhardt's identity-platform lead — a person named Andreas Vogt — and asking the question she always asked first.

"Show me a credential rotation for the predictive-maintenance service."

Andreas pulled up the chain entry for the most recent rotation — May 1, 2026, at 03:00 UTC. The service account `svc-pred-maint-inference` had its credentials rotated automatically. The rotation event emitted a sealed chain entry with the old credential fingerprint, the new credential fingerprint, the rotation reason (scheduled), the approver (Klaus Eberhardt, with a two-of-three approval from the platform engineering lead and the SecOps lead), the approval ticket reference, and the timestamp.

> **✓ Confirmation #7 — Eberhardt service-account IAM is chain-coupled end to end**
> Every credential lifecycle event — issuance, rotation, revocation, scope change — for the predictive-maintenance service produces a sealed chain entry. Policy-level changes (who can rotate, who can approve, what the rotation interval is) require two-of-three approval and are themselves chain entries. There is no path to change a service credential without producing a chain record.

"What does the engineer's identity look like? When Maximilian ran the verifier this morning — what was his SSO subject?"

"On-prem Active Directory federated to Azure AD. Maximilian's SSO subject is `max.brenner@eberhardt.de` — that is the Azure AD UPN. The federation is documented; the BSI IT-Grundschutz baseline says we have to be able to walk the federation chain back to the institutional credential, which is the on-prem AD account at Eberhardt. We can. There are no consumer identity providers in our trust chain."

"Twelve sites, twenty institutional IdPs?" Diana asked, half-joking, comparing notes against Helmstad's clinical-trial federation pattern.

Andreas smiled briefly. "Two. Eberhardt Stuttgart on-prem AD and Eberhardt Sindelfingen on-prem AD — same forest, two domain controllers. We do not federate to OEM customer IdPs because we do not give OEM customer staff access to the inference service directly. They consume the alerts on their side. Their identity is on their side."

Diana wrote that down. *Two-IdP shop. Easier than Helmstad's twelve. Easier than Olmstead's seventeen-school federation. Easier than Atrio's per-tenant SSO matrix. The Mittelstand culture shows in IAM — small footprint, well-controlled.*

> **✓ Confirmation #8 — Eberhardt's IAM trust chain is two on-prem AD domains, fully under chain governance**
> Eberhardt federates on-prem AD Stuttgart and on-prem AD Sindelfingen to Azure AD. Both are institutional. There are no consumer IdPs in the trust chain. Engineer SSO subjects in the chain entries resolve to the on-prem AD account in either domain. The BSI IT-Grundschutz baseline's identity-provenance requirement is satisfied without intermediate hops.

Meanwhile, in Paris, Raj had moved into the Lumière IAM walkthrough with Aurélien because Diana was on the Eberhardt side. Lumière's identity setup was Google Workspace as the primary IdP with a custom SSO layer in front of the model-build platform. Aurélien himself authenticated via Google Workspace. Hélène was the second factor on the model-build platform itself for any approval-bearing action — separation of duties between author and approver.

Each model-build chain entry recorded the author SSO subject and the approver SSO subject as separate fields. Raj cross-checked four entries; the author and approver were never the same person. The chain enforced that. A model-build entry with `author == approver` would fail the schema validator before it was sealed.

> **✓ Confirmation #9 — Lumière model-build IAM enforces author-approver separation in the chain schema**
> Every model-build chain entry has separate `author` and `approver` SSO-subject fields. The chain schema rejects an entry where the two are equal. Author and approver are typically a research-engineer (Aurélien or one of the team) and Hélène. Separation of duties is enforced by the chain, not by procedure.

Diana, on the video bridge from Stuttgart, asked the cross-cutting question: "Klaus and Hélène — who can rotate the model-handover-entry approver list at Eberhardt? Specifically, who can change the policy that says 'a model-handover entry from Lumière requires Klaus's approval'?"

Klaus answered. "Two-of-three approval. Myself, the platform-engineering lead, the SecOps lead. The change itself is a chain entry. The chain trumps the org chart on this."

Hélène, on the screen: "And on the Lumière side, the policy that says 'a model-build entry handed-over to Eberhardt requires Aurélien-author and Hélène-approver' is itself a chain entry under our IKM. Same shape — two-of-three approval to change the policy, and the policy change is sealed."

Diana wrote that down and underlined it.

---

### 🧪 12:00 PM CET — Joint Working Lunch by Video Bridge

The Stuttgart conference room had a small spread on a side table — Maultaschen, schnitzel, salad, a bowl of fruit. Klaus had insisted on Maultaschen, Swabian-style. The Paris side had a tartine with brie and a side of soup that Hélène had ordered up from the building's café. The video bridge stayed up. Both sides ate while talking.

Karen put down her plate.

"Let's talk about what we've seen so far. Mike — the inference path."

Mike: "Eberhardt's chain holds at eight months. The inference entries are bound to model-handover entries by reference. The model-handover entries reference Lumière's chain by ID. Klaus's two-of-three approval is on every handover."

Raj, on the screen from Paris: "Lumière's chain holds at four months. Build entries reference fairness-audit entries by ID. Author-approver separation is enforced by the schema. Aurélien-author and Hélène-approver on every build that handed over to Eberhardt."

Diana: "IAM on both sides is institutional and chain-governed. Two on-prem AD domains at Eberhardt. Google Workspace plus custom SSO at Lumière. Both sides have the policy-as-chain-entry pattern. Two-of-three approval to change either side's IAM policy."

Chen, on the screen from Paris: "I want to demonstrate the cross-anchor right now. Mike — share an Eberhardt model-handover entry. I have the corresponding Lumière build entry up."

Mike screen-shared his terminal. The model-handover entry `eb-mh-2026-04-01-lum-00012` from this morning was visible. The `model_artifact_sha256` field was clearly displayed: `sha256:7f2a91...c4d3`. The `lumiere_chain_entry_id` was `lum-mb-2026-04-01-bh-00031`.

Chen on his side opened the Lumière build entry. The `model_artifact_sha256` on Lumière's chain showed `sha256:7f2a91...c4d3`. Same value.

Chen ran a small Python helper on the Paris side that did a byte-equal compare of the two SHA-256 values pulled from the two terminals. The output:

```
Eberhardt model_artifact_sha256: 7f2a91...c4d3
Lumière  model_artifact_sha256: 7f2a91...c4d3
match: True

Eberhardt model_card_sha256: 8e1b82...a4f1
Lumière  model_card_sha256: 8e1b82...a4f1
match: True

Eberhardt fairness_audit_report_sha256: 9d3c72...e5f2
Lumière  fairness_audit_report_sha256: 9d3c72...e5f2
match: True
```

Klaus, watching from the Stuttgart side: "That is the seam holding."

Hélène, on the screen: "That is what BMW asked for."

> **✓ Confirmation #10 — Cross-vendor anchor matches live across the partnership boundary**
> The model_artifact_sha256, model_card_sha256, and fairness_audit_report_sha256 fields in Eberhardt's model-handover entry match the corresponding fields in Lumière's model-build entry bit-for-bit. The cross-vendor anchor composes. Verification is independent on each side; the match is the join. This is the mechanism BMW's vendor-management team asked the joint engagement to demonstrate.

Karen let the moment sit for a beat. Then she brought it back.

"Klaus — Hélène — for the deliverable, the language we'll use is that the cross-vendor anchor composes end-to-end without trust in either party's claim alone. Each side's chain is independently verifiable. The match is what makes the joint claim. That language survives an examiner read on both sides."

Klaus nodded. Hélène nodded.

Tom: "Stefan Kuhn at BMW asked us specifically whether the cross-vendor anchor composes without trust. The answer is yes, by independent verification plus byte-equal hash match. That is the right answer."

Karen: "Good. Save the live demo as evidence. We'll attach the screen-share frames to the deliverable."

Klaus took the last bite of his Maultaschen. "Karen — there is a question I want to put to you now, before the afternoon. We have been talking about the seam working. Tell me what does *not* work yet. Lumière and I both know there is something. Hélène and I would rather hear it now than at five-thirty."

Karen looked at Raj on the screen. Raj had been the one who had found the 90-day retention issue this morning.

"Raj," Karen said.

Raj walked them through it. The Lumière training-data shards rotated out at 90 days post-delivery. The Eberhardt deployment window was typically 9 to 18 months. If a regression appeared six months in and someone wanted to retrace which shards produced it, the shards were gone. The manifest hash was retained, the manifest document was retained, the chain entry was retained — but the shards themselves were not. This was the 2024-incident-class forensic gap.

Hélène was quiet for a beat.

"That is fair. The 90-day policy was set under GDPR data-minimization. I will say honestly — I had not modeled it against the Eberhardt deployment window. I should have. I will take a proposal to the founders this week to extend retention to the maximum deployment window with explicit GDPR justification, probably under Article 6(1)(f) legitimate interest balanced against the storage-limitation principle. Klaus — I want your engineering team's max-deployment-window number to use as the retention parameter."

Klaus nodded. "We can give you that. The longest we have run a model is twenty-two months. I would size it to twenty-four months."

"Twenty-four months it is."

Karen: "Document the retention extension as a CAPA against this audit. The deliverable will name it as Partial #1 with a noted in-flight remediation. Tom — log it."

Tom: "Logged."

The bridge stayed up while everyone finished lunch. The mood was something like the last twenty minutes of a long meeting where the work was done and the people had settled into their seats.

---

### 🔄 1:00 PM CET — Mike on the API Layer (Stuttgart) and Elena on the CRM (Paris)

Mike's afternoon was the API surface. Eberhardt runs the inference service behind an internal API gateway with mutual-TLS authentication and a request-signing layer between the OEM-facing edge and the inference core. Every call to the inference service — from the OEM-customer-facing edge layer, from the internal validation harness, from the model-handover acceptance test — goes through this gateway.

Every call emits a chain entry. Request, headers (filtered for telemetry that could re-identify an individual vehicle owner), authorizer decision, downstream service, response code, response hash, latency. Sealed.

"Show me a call that returned a 5xx," Mike said.

Maximilian pulled up an entry from April 14. 503. The HSM at Sindelfingen had been unreachable for 90 seconds during a planned maintenance window. The chain entry recorded the failure. The downstream inference service had served a fallback `model-unavailable` response that explicitly told the OEM-side edge layer to fall back to its own conservative-default behavior — in this case, to suppress new alerts and continue surfacing the most recent confirmed alert if any. The fallback was itself sealed as a separate inference-shaped entry.

> **✓ Confirmation — API gateway and fallback paths are sealed end to end**
> Every inference API call — including failed authorizations, planned-maintenance HSM outages, and fallback responses — produces a sealed chain entry. The April 14 maintenance window shows up cleanly: 47 inference calls during the 90-second window, all routed to the fallback responder, all sealed. No silent gaps. The OEM-side edge layer received explicit fallback signals and behaved conservatively.

Mike turned to the mutual-TLS authenticator. The mTLS certificate that Eberhardt presents to BMW's vehicle-side ingestion endpoint is rotated quarterly. The rotation produces a chain entry. The certificate fingerprint is bound to the rotation entry. BMW has the public-key portion of Eberhardt's quarterly rotation calendar in their vendor-management system.

> **✓ Confirmation — mTLS certificate rotation is chain-coupled with downstream OEM**
> Every quarterly rotation of the Eberhardt-to-BMW mTLS certificate produces a sealed chain entry with the new fingerprint and the rotation approval. BMW receives the public-key calendar through their vendor-management portal. The rotation is auditable from both ends.

Meanwhile, in Paris, Elena had spent thirty minutes with Lumière's revenue-operations lead — a person named Sophie Lacombe — going through the customer-engagement CRM. Lumière uses Pipedrive for sales and customer-engagement management. The Pipedrive instance has Eberhardt as the largest active account, with notes from quarterly business reviews, contract-renewal calendar entries, and contact-history.

The Pipedrive is not chain-coupled. It is a commercial CRM, not a technical-evidence system. The fairness-audit reports, model-handover records, deployment-tracking — none of those flow into Pipedrive. Pipedrive holds the relationship history, not the engineering history.

Elena flagged this and moved on. *Out of scope for the AI-chain audit. Same shape as Helmstad's CTMS-as-KOL-CRM pattern, but cleaner — Pipedrive at Lumière holds nothing the chain needs to anchor. The CRM is the CRM and the chain is the chain.*

> **✓ Confirmation — Lumière's CRM (Pipedrive) is correctly scoped as out of chain coverage**
> Pipedrive holds commercial relationship history — contracts, renewal calendars, contact notes. It does not hold technical-evidence artifacts. The chain does not anchor Pipedrive entries because the chain does not need to. The boundary is documented and intentional. Same architectural decision as the Helmstad-pattern out-of-scope CRM.

Elena closed her notebook on the CRM section by 1:30 and moved over to help Chen on the data-pipeline walk that was about to start.

---

### 🧬 2:00 PM CET — Chen on the Lumière Training-Data Pipeline (Paris)

Chen had been waiting for this part. The chain at the model-build boundary is one thing. The chain at the training-data-input boundary is the other. They are not the same thing.

Aurélien walked Chen through the training-data ingestion pipeline. The raw signals — battery cell voltages, temperature curves, charge cycle counts, motor-controller telemetry — came from Eberhardt's field deployment, anonymized and aggregated by Eberhardt's data team, transferred to Lumière over a GDPR Article 28 processor-agreement-governed encrypted channel. The transfer arrived as a daily tarball, signed by Eberhardt's transfer key. Lumière's ingestion service verified the signature, computed a SHA-256 of the tarball, and recorded the hash as a referenced artifact in a chain entry of type `training-data-ingestion`. The records were then unpacked, parsed, validated against the schema, and written into the training-data warehouse. The warehouse row carried a foreign key to the ingestion chain entry.

"That is the same shape as the Eberhardt SFTP-ingestion pattern in reverse," Chen said. "Eberhardt sends, you receive. The chain entry on your side records what you received. The chain entry on Eberhardt's side records what they sent. The hashes match."

"Yes."

Chen pulled up a recent ingestion entry and the corresponding Eberhardt outbound entry. The tarball SHA-256 matched. The transfer key fingerprint matched. The tarball file count matched. The byte count matched.

> **✓ Confirmation — Training-data transfer is chain-coupled at both ends across the partnership boundary**
> Eberhardt's outbound transfer entry and Lumière's ingestion entry record matching SHA-256, file count, byte count, and signing-key fingerprint. The transfer is GDPR-Article-28-governed. The chain is the integrity record on both sides. This is the second cross-vendor anchor — alongside the model-handover anchor — and it composes the same way.

Chen circled the transfer-key-fingerprint match in his notebook.

"GDPR cross-border transfer," Chen said. "Stuttgart to Paris. Both within the EU. The Article 28 processor agreement governs it. But the chain entry — the transfer entry on either side — does it carry an explicit cross-border-transfer attribute?"

Aurélien paused. "It carries the source country code, the destination country code, the processor-agreement reference, and the transfer-mechanism reference. It does not carry an explicit `cross_border_transfer = true` boolean."

"Within the EU you do not need the explicit boolean for legal purposes — the GDPR is the GDPR end-to-end. But for examiner-readability — and for the cross-jurisdictional registry that a non-EU regulator might be reading — would it be useful?"

Aurélien thought. "Probably yes. The CNIL would not need it. The German LfDI would not need it. But the BMW vendor-management read might want it. And if we ever add a non-EU customer who wanted to read the chain on the Eberhardt deployer side — they would find it useful."

Chen wrote that down. *Same shape as Sun-Won's PIPA Section 28 cross-border-transfer attribute finding — but within the EU and thus less pressing. Worth flagging as a Nit.*

> **⚠️ Nit #1 — Cross-border-transfer attribute is implicit in the chain, not explicit**
> Eberhardt-to-Lumière training-data transfers cross the German-French border within the EU. The chain entries on both sides record the source country code, destination country code, and processor-agreement reference, but do not carry an explicit `cross_border_transfer` boolean attribute. Within-EU transfers do not need this for GDPR purposes — the regime is the same end-to-end. But for examiner-readability, especially for non-EU regulators or for BMW's vendor-management read, an explicit attribute would be useful. Same shape as the Sun-Won PIPA Section 28 finding from a prior engagement, scaled down to within-EU pressing-ness. Recommended as a normative attribute addition for the next chain-schema revision.

Chen moved on to the training-data manifest binding. Aurélien showed him the manifest document for the v4.2.1 training run. Every shard listed with source, ingestion date, de-identification scheme, consent basis, and SHA-256. The manifest itself was sealed — its hash was the `training_data_manifest_sha256` in the model-build entry.

"And the consent basis on each shard," Chen asked. "What is it?"

"GDPR Article 6(1)(f) legitimate interest, with the data-protection-impact assessment for the predictive-maintenance use case as the documented justification. The DPIA is signed by Hélène and the customer-side data-protection officer at Eberhardt. The DPIA hash is in the manifest."

> **✓ Confirmation — Training-data manifest binds shards to consent basis with DPIA hash**
> Every shard in the training-data manifest is bound to its consent basis (GDPR Article 6(1)(f) legitimate interest), its DPIA reference, and its de-identification scheme. The DPIA itself is hashed into the manifest. The manifest hash is sealed in the model-build entry. End-to-end consent provenance is auditable from the model artifact back to the legal basis.

Chen closed his notebook on the training-data section. *Two cross-vendor anchors. Model-handover and training-data-transfer. Both compose. The chain reaches the legal-basis layer.*

---

### 🧬 2:00 PM CET — Luis on Logs and Pipelines (Stuttgart)

Luis was on the Stuttgart side, at the same hour. He had spent the morning quietly going through the Eberhardt log-pipeline and the daily-seal job. By 2:00 PM he had a clean read on it.

The seal job ran at 23:55 local time at Sindelfingen. It pulled the day's chain entries, computed the daily Merkle root, signed it with the on-prem HSM, and published the seal record to an internal seal-archive bucket. The seal-archive bucket had versioning enabled, was on an immutable-storage tier with a 7-year retention lock under TISAX, and was replicated to a second data center in Munich.

The CloudWatch-equivalent — Eberhardt used Splunk Enterprise on-prem — had log-deletion-prevention configured at the index level. Splunk admin operations themselves emitted chain entries.

Luis asked the question he always asked.

"Who can disable the daily seal?"

The Eberhardt SecOps lead — a person named Bettina Hofer — answered. "No one, unilaterally. The seal job is a Kubernetes CronJob in a namespace where the deployment manifest is itself a chain entry. To disable the seal you would have to issue a chain-bound deployment change with two-of-three approval. The two-of-three are Klaus, the platform-engineering lead, and me. We have not disabled it. We do not plan to."

"What if the HSM is unreachable at seal time?"

"The seal job retries for 90 minutes. If it cannot complete, it emits a `seal-failure` chain entry — the chain itself records the failure — and pages the SecOps lead. There is a one-time-only manual seal-completion procedure documented for the case where the HSM has a hardware failure; that procedure requires three-of-three approval and the manual seal is itself a chain entry of a special type."

> **✓ Confirmation — Seal job is chain-bound and chain-failure-aware on the Eberhardt side**
> The daily seal job at Sindelfingen runs as a chain-governed Kubernetes CronJob. Disabling the job requires two-of-three approval and produces a chain entry. Seal failures are themselves chain entries. Manual seal completion exists as a documented procedure with three-of-three approval and a special-type chain entry. The chain is reflexive about its own operation.

Luis wrote that down. He moved to the log-replication path between Sindelfingen and Munich. The replication job was a separate CronJob, also chain-governed. He spot-checked the cross-region replication-completion entry from the prior night. PASS.

He asked Bettina the same question Lumière would be asked in Paris by Raj an hour later: "Where does the seal-archive bucket live, jurisdictionally?"

"On-prem in Sindelfingen, replicated to on-prem in Munich. Both are German jurisdiction. We do not replicate to cloud or to another country. BSI IT-Grundschutz baseline plus our own preference."

"Mittelstand," Luis said.

Bettina half-smiled. "Mittelstand."

---

### 📊 3:00 PM CET — Joint Reconciliation Test by Video Bridge

The video bridge was back up at 3:00 PM. Karen had asked both sides to pick a single deployment and trace it end-to-end through both chains. She had asked Klaus and Hélène to be on the bridge for it.

Maximilian on the Stuttgart side picked a recent one: an `urgent-service-required` predictive-maintenance alert from April 8 for a BMW iX. Anonymized as `Customer-X` for the deliverable. The vehicle VIN was hashed in the chain entry.

Maximilian read out the trace. Karen watched it on the screen. Tom took notes.

**Step 1.** The vehicle VIN-hash on the BMW iX. (BMW would have the unhashed VIN on their side. The hash matches.)

**Step 2.** The Eberhardt inference chain entry `eb-bh-2026-04-08-bmwix-00874`. Sealed. Verifier PASS in 4 seconds. The entry's `model_handover_entry_ref` points to `eb-mh-2026-04-01-lum-00012`.

**Step 3.** The Eberhardt model-handover entry `eb-mh-2026-04-01-lum-00012`. Sealed. Verifier PASS in 4 seconds. The entry's `lumiere_chain_entry_id` points to `lum-mb-2026-04-01-bh-00031`.

**Step 4.** Hélène — by video bridge — granted Mike read access to the Lumière chain for the matching entry. Mike ran the Lumière verifier against `lum-mb-2026-04-01-bh-00031` from his Stuttgart terminal. Verifier PASS in 4 seconds. The entry's `model_artifact_sha256` matches Eberhardt's by byte-equal compare. The entry's `fairness_audit_entry_ref` points to `lum-fa-2026-03-31-bh-00007`.

**Step 5.** The Lumière fairness-audit entry `lum-fa-2026-03-31-bh-00007`. Sealed. Verifier PASS. The fairness-audit-report SHA-256 matches the value sealed in Eberhardt's model-handover entry.

**Step 6.** The Lumière model-build entry's `training_data_manifest_sha256` resolves to a manifest document. The manifest document hashes to the value sealed. The manifest lists the shards used for v4.2.1 training. (Within 90 days; the shards themselves are still present for this April-1 build.)

**Step 7.** Each shard in the manifest is hash-bound to its source. The Eberhardt outbound transfer entries on the Eberhardt side and the Lumière ingestion entries on the Lumière side reconcile by tarball SHA-256. End-to-end consent provenance traces back to the DPIA sealed in the manifest.

The bridge sat in silence for a moment after Maximilian finished reading.

Karen spoke first.

"Seven legs. Two chains. One BMW vehicle. End-to-end traversable with byte-equal hash matches at every cross-vendor anchor. That is the joint claim."

Klaus nodded once.

Hélène, on the screen: "That is what BMW asked for. That is what the AI Act conformity assessment asks for."

> **✓ Confirmation — Joint end-to-end reconciliation traverses BMW vehicle → Eberhardt inference → Eberhardt model-handover → Lumière model-build → Lumière fairness-audit → Lumière training-data manifest with byte-equal hash matches at every cross-vendor anchor**
> One BMW iX vehicle's predictive-maintenance alert traced through seven legs across two independently-IKM-governed chains. Verifier PASS at every leg in 4 seconds per check. Byte-equal hash matches at every cross-vendor anchor. End-to-end provenance from the OEM customer's vehicle back to the GDPR consent basis on the training data, with the DPIA hashed into the manifest. This is the joint-supplier audit BMW's vendor-management asked for.

Mike, on the Stuttgart side, looked up from his terminal.

"That trace took eleven minutes from kickoff to closing PASS. Including the Hélène-grants-credentials step."

Karen: "Document the elapsed time. The deliverable will note that an end-to-end joint trace can be completed in under fifteen minutes once both sides cooperate. That is a number BMW will want to see."

Tom: "Logged."

---

### 🧪 3:30 PM CET — The Audit-Report-Language Question

Diana was the one who had been turning the language attribute over in her head since 10:00 AM. She brought it back up by video bridge.

"The fairness-audit report is in French. The Eberhardt deployment documentation is in German. The model card is bilingual French-English. BMW's vendor-management lead — Stefan Kuhn — what does he read?"

Klaus: "Stefan reads German first, English fluently, French only well enough to follow a technical paper."

Hélène: "The fairness-audit report has an executive summary in English. The full body is in French. The translation of the full body to English would take an external translation effort. We have not done it for v4.2.1."

Diana: "So a German examiner reading the joint deliverable would have the model card in two languages and the fairness-audit body in only one."

Klaus thought for a beat. "The BSI examiner reads English fluently. The LfDI Baden-Württemberg reads German. The CNIL reads French. The BMW vendor-management read is German-English. So the reader matrix is German, English, French — three languages. The model card is bilingual French-English. The fairness-audit body is French only with an English executive summary. The deployment-side documentation at Eberhardt is German with English appendices."

Karen: "So an examiner reading across the seam — say BSI plus CNIL — has access to the artifacts in their preferred languages, but the *individual* fairness-audit body is French only."

Hélène: "Yes. We can produce an English translation of the fairness-audit body for v4.2.1 within four weeks. For new models we can include English-French bilingual full-body audit reports as a standing practice."

Diana: "And the chain entry's `audit_report_language` field is what flags this. It is captured. The reader knows from the chain what language the report is in. The question is whether the reader has access to the report in their preferred language."

Karen: "Right. That is a procedural finding, not a chain-integrity finding. The chain captures the language attribute correctly. The procedural recommendation is to provide English translations of all fairness-audit bodies as standing practice and to add a normative `audit_report_languages` array — plural — that enumerates available translations."

> **⚠️ Nit #2 — `audit_report_language` is singular; multilingual audit reports would benefit from an array**
> The chain entry's `audit_report_language` field captures the primary language of the audit report (`fr-FR` for v4.2.1). When the report is translated, the field does not naturally accommodate the multiple languages now available. Recommend extending the schema to `audit_report_languages` (plural, array) so a translated report can be discovered through the chain. As a procedural matter, providing English translations of all fairness-audit bodies as standing practice would address the cross-jurisdictional reader-matrix question for the German-French-English BMW examiner audience.

Hélène: "Agreed. We will produce the v4.2.1 English translation in the next four weeks and adopt bilingual-or-trilingual audit bodies as standing practice for new models. The schema-attribute change is something I would coordinate through the chain-schema revision channel."

Karen: "Document it both ways. Procedural CAPA at Lumière for the standing practice. Schema recommendation for the next normative revision."

Tom: "Logged."

---

### 😬 3:45 PM CET — Friction: The 2024 Model-Drift Incident

Klaus had been quiet through the language discussion. He came back to the bridge as it ended.

"There is one more thing I want to put on the table. The 2024 model-drift incident. Hélène — how do you want to talk about it?"

Hélène took a beat before answering. The video screen caught the small motion as she sat back in her chair.

"Honestly. The 2024 incident was on us. We delivered a battery-health model that had silently regressed after a training-data update. The validation set had not caught the regression because the new training data shifted the distribution of the test set as well. The model passed our internal validation but was actually worse on the field distribution. We delivered it to Eberhardt. Eberhardt deployed it to the OEM fleet. About two weeks of degraded predictions before someone at BMW noticed and asked Eberhardt to investigate. Eberhardt asked us. We rolled back. Total degraded-prediction window was something like fifteen days. No safety incidents — the regressions were over-triggering of `service-recommended` alerts, not under-triggering of `urgent-service-required` — but the OEM customer-experience cost was real. We did not have the chain at the time. We could not retrace which deployments had received the regressed model versus the prior good model. We had to do it by reconstruction from delivery logs. It took three days. Three days that would have been three minutes with the chain."

Klaus, very evenly: "Karen — does the chain prevent another 2024?"

Karen took a beat.

"It does not prevent. It detects. The chain proves which model produced which inference. It does not prove the model is good. The detection happens earlier because the chain gives BMW or Eberhardt an immediate way to ask 'which model produced this anomaly?' without three days of reconstruction. The chain shrinks the silent-regression window from two weeks to whatever the field-monitoring cadence is. If field monitoring is daily, the window is one day. If it is weekly, the window is one week. The chain does not change the field-monitoring cadence."

Klaus: "So the answer is faster detection, not prevention."

Karen: "Yes. Faster detection. The chain is the audit trail. It is not the validator."

Hélène: "And the 90-day retention extension we agreed to this morning matters here. If a regression appears six months after deployment — which is the 2024 shape — the chain detects fast, but the *retraining* root-cause analysis needs the training shards. The 90-day window is too short. Twenty-four months — the maximum-deployment-window number — is the right size."

Karen: "Right. The chain is the detection layer. The retention extension is the root-cause-analysis layer. Both matter. Both are findings of this engagement."

Klaus, after a pause: "Detection without prevention is still better than two-week silent regression."

Karen: "Yes. Materially better. But the deliverable has to be honest about the distinction."

Tom: "I will draft the language for the executive summary. The chain detects faster; the chain does not prevent regression; the field-monitoring cadence and validator quality are separate concerns. The 90-day retention extension supports retroactive root-cause analysis for the long-tail regression case."

Klaus and Hélène both nodded.

The bridge stayed quiet for another beat. Then Klaus said:

"Thank you for putting it that way. I wanted it on the record that the chain is not magic. The 2024 incident is something Hélène and I have lived with. I do not want anyone reading the deliverable to think this audit closes the 2024 question. It does not. It changes the cost of the *next* one."

Karen: "That is the right framing. We will use it."

---

### 🔍 4:30 PM CET — The BMW Question

Klaus brought up the question that had been hanging over the engagement since kickoff. By video bridge.

"Stefan Kuhn at BMW asked us last month — and I am going to ask Karen the same question. If a BMW customer's vehicle has an AI-driven false-positive predictive-maintenance alert that costs the customer time and BMW money, can BMW determine whether the false positive came from a regression in Lumière's model or from Eberhardt's integration or from BMW's own vehicle-side data?"

Karen had been ready for the question. She had been thinking about it since lunch.

"The cross-vendor anchor does that. Three chains, one root-cause path."

She walked it through.

"The BMW vehicle-side data — what the sensor reported, what the BMW-side ingestion logged — that is BMW's chain or BMW's equivalent of a chain. They have it on their side. It is the input.

"The Eberhardt inference chain entry shows the input telemetry hash that Eberhardt received from the BMW vehicle. If BMW's input does not match the input Eberhardt processed, the difference is in the BMW-to-Eberhardt transport — that is BMW's integration question. If they match, the input is reconciled.

"The Eberhardt inference entry shows the model artifact hash. The model-handover entry shows the model-handover record. The Lumière model-build entry shows the model build provenance. If the model artifact at inference matches the model-build entry, the model is reconciled. If the model artifact differs from any prior good build, the regression is in Lumière's build. That is Lumière's responsibility.

"The Eberhardt inference entry shows the inference output. If the output looks anomalous given the input — and the input and the model are both reconciled — the regression is in the model's behavior on this input class. That is again Lumière's responsibility, but in the modeling-quality sense, not the build-integrity sense.

"The Eberhardt-to-BMW transport — what Eberhardt sent to BMW after inference — is the next leg. Eberhardt's outbound logging shows what was sent. BMW's inbound logging shows what was received. If they match, the transport is reconciled. If they don't, the integration question is on the boundary between Eberhardt and BMW.

"Three chains. One root-cause path. Each chain answers its own leg. The cross-vendor anchor is what makes the join possible without trust in a single party's claim. That is exactly the joint-supplier audit BMW asked for."

Klaus had been listening with his fingers on his chin.

"That is the answer."

Hélène, on the screen: "That is the answer."

> **✓ Confirmation — Three-chain root-cause path is end-to-end traversable for false-positive analysis**
> A BMW false-positive predictive-maintenance alert can be root-caused across three chains — BMW's vehicle-side data, Eberhardt's inference chain, and Lumière's model-build chain — by reconciling input-telemetry-hash, model-artifact-hash, and inference-output-hash at the cross-vendor anchors. Each chain answers its own leg without requiring trust in another party's claim. This is the BMW joint-supplier audit deliverable.

Tom: "I want to capture the answer verbatim. Stefan Kuhn will read it. He will know we addressed his question explicitly."

Karen: "Capture it."

---

### 🌆 5:30 PM CET — Joint Debrief by Video Bridge

The full team was on the bridge at 5:30 PM. Stuttgart had moved everyone into the conference room. Paris had pulled chairs up around Hélène's desk. Klaus and Hélène were both on. Tom and Karen led.

Karen ran through the per-company summary first.

"Eberhardt-side individual findings: zero Gaps, zero Partials. The chain holds at eight months. IAM is two-AD-domain on-prem federated to Azure AD, fully chain-coupled. Seal job is chain-governed with HSM at Sindelfingen. mTLS rotation to BMW is chain-coupled with downstream OEM. BSI IT-Grundschutz alignment, ISO 27001, ISO/IEC 27017, TISAX assessment artifacts all feed from the chain.

"Lumière-side individual findings: zero Gaps, zero Partials. The chain holds at four months. IAM is Google Workspace plus custom SSO with author-approver schema enforcement. Seal job is chain-governed with HSM at OVHcloud Paris. ANSSI-aligned configuration. Fairness-audit hash anchoring works.

"Joint findings — the load-bearing part of this engagement: zero Gaps, one Partial, two Nits.

"Partial #1: Lumière's 90-day training-data retention is shorter than Eberhardt's typical deployment window. Hélène has accepted the finding and committed to a 24-month retention extension under GDPR Article 6(1)(f) legitimate interest with documented justification. The retention extension is in flight as a CAPA against this audit.

"Nit #1: Cross-border-transfer attribute is implicit in the chain. Recommend an explicit `cross_border_transfer` boolean for examiner-readability, especially for non-EU readers and for BMW's vendor-management read. Within-EU transfers do not need this for GDPR purposes — same regime end-to-end — but the explicit attribute is procedurally cleaner. Same shape as the Sun-Won PIPA Section 28 finding from a prior engagement, scaled down to within-EU pressing-ness.

"Nit #2: `audit_report_language` is singular; multilingual audit reports would benefit from an array. Recommend the schema extend to `audit_report_languages` plural. Procedural recommendation: provide English translations of all fairness-audit bodies as standing practice for the German-French-English BMW examiner audience. Hélène has committed to producing the v4.2.1 English translation in the next four weeks.

"Joint posture: the cross-vendor anchor mechanism is the engagement's load-bearing finding. Two chains compose at the model-handover boundary. End-to-end reconciliation traverses seven legs from BMW vehicle through Eberhardt inference to Lumière training-data manifest, with byte-equal hash matches at every cross-vendor anchor, in under fifteen minutes elapsed."

Klaus: "Read out the regulator-by-regulator posture."

Karen: "Five readers.

"**EU AI Act.** Article 11 deployer logging: satisfied by Eberhardt's inference chain. Article 12 conformity assessment: satisfied by Eberhardt's chain plus Lumière's fairness-audit chain plus the cross-vendor anchor. Article 16 provider obligations: satisfied by Lumière's fairness-audit and model-build chains. The joint chain is a defensible Article 11/12/16 evidence pack.

"**BSI (Bundesamt für Sicherheit in der Informationstechnik).** IT-Grundschutz baseline: satisfied. The seal infrastructure, HSM configuration, log retention, identity-provenance walk are all aligned. The ISO 27001 and ISO/IEC 27017 certifications provide the formal layer.

"**TISAX.** Auto-supply community assessment: the chain artifacts feed the TISAX assessment. Coordinated supplier security assessment is consistent with TISAX practice.

"**LfDI Baden-Württemberg and CNIL (joint GDPR read).** Article 28 processor agreement governs Eberhardt-to-Lumière training-data transfer. Within-EU transfer documented. Cross-border-transfer attribute is implicit (Nit #1). Article 6(1)(f) legitimate interest as the consent basis with DPIA hashed into the training-data manifest. Article 5(1)(e) storage limitation is the driver for the 90-day retention; the retention extension to 24 months will be re-justified under legitimate-interest balancing.

"**BMW vendor-management (Stefan Kuhn).** Cross-vendor anchor composes end-to-end without trust in either party's claim. Three-chain root-cause path is traversable for false-positive analysis. The 2024 model-drift incident class would now resolve in minutes instead of days for detection, and within the 24-month retention window for retroactive root-cause analysis. Joint-supplier audit deliverable is satisfied.

"**ISO 26262 functional-safety auditor.** The chain artifacts support the audit-trail requirements for AI-derived safety-related outputs. The chain is independent of ISO 26262 but compatible.

"That is the regulator matrix."

Klaus took a long breath in.

"Hélène. We are good?"

Hélène, on the screen: "We are good. The 90-day retention extension is the only structural change. The two Nits are schema and translation work. The chain holds across the seam."

Klaus turned to the camera. "Karen — Tom — the Stuttgart and Paris teams. Thank you. This is what I had hoped for. The deliverable for BMW will reference this audit explicitly."

Karen: "Thank you. The formal deliverable will be in your inbox by end of next week. The joint posture, the per-regulator matrix, the three findings, the live-trace evidence, and the language for Stefan Kuhn's three questions — all of it will be documented."

Hélène: "And from Lumière — thank you. The 90-day retention conversation is one I should have had with myself months ago. Having it forced by the audit is the right reason to have it now."

Karen: "That is what audits are for."

The bridge held for another minute while everyone exchanged the small post-meeting pleasantries — handshakes by camera, the small wave from Klaus to the Paris side, a thumbs-up from Aurélien to Mike. Then the bridge dropped.

The Stuttgart team packed up. Maximilian stayed for a moment to trade contact details with Mike. Andreas stayed to walk Diana through one final question about the AD-to-Azure-AD federation that he wanted to clarify off the record. Bettina caught Luis in the corridor and gave him a recommendation about a Splunk admin-operation chain-coupling pattern that they had been considering for the next quarter. Tom shook Klaus's hand a second time at the door.

In Paris, Sophie walked Elena to the lobby. Aurélien stayed to chat with Chen and Raj about the cross-border-transfer attribute schema proposal — Chen was already drafting it on his laptop. Hélène walked them all to the door at the bottom of the Haussmann staircase and said goodbye in three languages depending on whom she was speaking to.

---

## Final ✅ Comparison Block

### Per-Company Posture

| Dimension | Eberhardt (Stuttgart) | Lumière (Paris) | Joint (Cross-Vendor Anchor) |
|---|---|---|---|
| Chain duration | 8 months | 4 months | 4 months at handover |
| HSM | Thales Luna on-prem, Sindelfingen | OVHcloud Paris-region, ANSSI-aligned | Independent IKMs |
| Verifier | `herald-verify`, 4 sec, 12 steps | `herald-verify`, 4 sec, 12 steps | byte-equal hash match at anchor |
| Service-account IAM | chain-coupled | chain-coupled | n/a |
| Author-approver separation | n/a | schema-enforced | n/a |
| Compliance baselines | BSI IT-Grundschutz, ISO 27001, ISO/IEC 27017, TISAX | ANSSI-aligned, GDPR | Article 28 processor agreement |
| Functional-safety alignment | ISO 26262 trail-supporting | n/a | n/a |
| Individual gaps | 0 | 0 | n/a |
| Individual partials | 0 | 0 | n/a |
| Joint partials | n/a | n/a | 1 (90-day retention) |
| Joint nits | n/a | n/a | 2 (cross-border attribute, language array) |

### Per-Regulator Posture

| Reader | Posture | Rationale |
|---|---|---|
| EU AI Act (Article 11 / 12 / 16) | Satisfied | Joint chain is a defensible deployer-and-provider evidence pack |
| BSI IT-Grundschutz | Satisfied | Seal infrastructure, HSM, identity-provenance aligned |
| ISO 27001 / ISO/IEC 27017 | Satisfied | Formal certification layer over the chain |
| TISAX | Satisfied | Chain artifacts feed the auto-supply community assessment |
| LfDI Baden-Württemberg | Satisfied with Nit | Cross-border-transfer attribute implicit, not blocking |
| CNIL | Satisfied with Nit | Same as LfDI; within-EU regime is end-to-end consistent |
| BMW vendor-management (Stefan Kuhn) | Satisfied | Cross-vendor anchor + three-chain root-cause path |
| ISO 26262 functional safety | Trail-supporting | Independent of TesseraSeal but compatible |

### The Three Findings

| # | Severity | Finding | Status |
|---|---|---|---|
| 1 | Partial (Joint) | Lumière's 90-day training-data retention is shorter than Eberhardt's typical 9-18-month deployment window. Forensic gap for retroactive regression root-cause analysis. | Hélène committed to 24-month extension under GDPR Article 6(1)(f) with documented justification. CAPA in flight. |
| 2 | Nit (Joint) | Cross-border-transfer attribute is implicit in chain entries. Within-EU transfers do not need it for GDPR purposes, but explicit attribute helps non-EU readers and BMW's vendor-management read. | Schema-revision recommendation for next normative chain-schema update. |
| 3 | Nit (Joint) | `audit_report_language` is singular; multilingual reports would benefit from an array (`audit_report_languages`). Procedural: provide English translations of fairness-audit bodies as standing practice for the BMW German-French-English reader matrix. | Hélène committed to v4.2.1 English translation in 4 weeks; bilingual-or-trilingual audit bodies as standing practice; schema-extension recommendation. |

---

### 🧾 Final Assessment Theme

The TGV from the Stuttgart side back across the border was tomorrow's plan. Tonight the team would all stay put — Karen at the Stuttgart hotel, Raj at his Paris hotel near Bastille — and write up the deliverable on shared docs. The drive back to the hotel from Eberhardt was twenty minutes through the Stuttgart hills as the sun was going down. Karen had her coffee, refilled, in the cup holder.

She thought about the day. About Maximilian's terminal at 9:15 in the morning showing a clean PASS on the inference entry. About Aurélien's terminal at 10:00 in Paris showing a clean PASS on the model-build entry. About Chen at noon on the bridge running the byte-equal compare and the hashes matching across both chains. About Klaus's question at 3:45 — *does the chain prevent another 2024?* — and her honest answer that it does not prevent, it detects faster, and faster detection is materially better than two-week silent regression. About Hélène's quiet acknowledgment of the 90-day retention asymmetry and her unhesitating commitment to extend. About Klaus's question at 4:30 — *can BMW root-cause a false positive across the seam?* — and the three-chain answer that took her ninety seconds to walk and that Stefan Kuhn would read in the deliverable.

She thought about Northbridge — the gold standard, fully sealed, one Nit. About Mercator — the bifurcation, AI on the chain, EHR off. About Stelvio — the three-zone, AI sealed, OT mutable, IT business legacy. About Atrio — the multi-tenant, forty-seven tenants under twelve sponsor-bank IKMs, fourteen hundred verifier runs, zero failures. About Helmstad — the biopharma CRO boundary, Quintessa's PGP-signed SFTP, source-side history outside the chain. About Pacific Crescent — the utility, AI gas-pipeline leak detection, Brentwood real-leak. About Olmstead — the university, two override-down decisions where the rationale was gone. About Sun-Won — the Korean engagement, PIPA Section 28 cross-border-transfer attribute. About the engagement after Sun-Won.

Today was different. Today was the first joint engagement. Two companies, two countries, one chain that crossed the boundary. The cross-vendor anchor was the load-bearing thing — a more rigorous version of what Helmstad's CRO had achieved with PGP-signed PDFs, scaled to mutual chain-on-both-sides verification with byte-equal hash matches at the seam. Not a SOC 2 attestation. Not a contractual representation. A cryptographic compose.

The chain at Eberhardt held. The chain at Lumière held. The seam held.

And the seam — *the seam was the engagement*. Each side's chain alone was easy to confirm. The interesting work was always going to be the join. The ninety-day retention asymmetry. The audit-report-language matrix. The implicit cross-border-transfer attribute. None of those were chain-integrity findings. They were findings about how two chains, written by two teams in two countries under two compliance baselines, talk to each other through a normative anchor mechanism that has to be intelligible to a third reader — BMW's Stefan Kuhn — who has to root-cause a false positive without trusting either party's claim alone.

The 2024 model-drift incident was the lesson. Pre-chain, three days of reconstruction. Post-chain, three minutes of detection. The 24-month retention extension closes the long-tail root-cause window. The chain does not prevent the next 2024. It changes the cost of finding out.

Karen picked up her phone at a red light on the way down out of the hills and dictated a one-line note for the report's executive summary.

> *"The chain at Eberhardt proves what was deployed. The chain at Lumière proves what was built. The cross-vendor anchor proves they are the same model — and that is what the joint-supplier audit asked for."*

The light turned green. She put the phone down and drove the rest of the way to the hotel as the Stuttgart sky went dark over the Schlossplatz.

---
