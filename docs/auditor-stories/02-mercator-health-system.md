## 🧾 Diary of an Audit — Mercator Health System

**Engagement:** HITRUST CSF v11 + HIPAA Security Rule + FDA SaMD Post-Market Combined Assessment
**Client:** Mercator Health System (top-20 US integrated health system — academic medical center, regional health-insurance carrier, multi-state physician group)
**Status:** AI tracing live for 90 days. Everything else is legacy.
**Audit team lead:** Dawn
**Client liaison:** Dr. Patricia Okonkwo, system-wide CISO

**Posture going in:** same product the team saw at Northbridge Bank two weeks ago, but applied to a healthcare setting. The crypto substrate is familiar; the question now is whether the operational controls translate to HIPAA + FDA-regulated workflows. Mercator turned on TesseraSeal 90 days ago, independent of any prior conversation with this audit team. The institution's claim, going into the room, is that every prediction the FDA-cleared sepsis CDS model makes lands in a sealed chain-of-custody ledger built on Herald Core, that the daily seals are rooted in AWS CloudHSM, that the verifier CLI is `herald-verify`, and that the whole stack conforms to a public spec called FFIEC chain-of-custody v1.0b. The team is confident in the product; the open question is whether Mercator's deployment lives up to it.

---

### Context

Mercator turned on TesseraSeal 90 days ago. Not across the enterprise — just on the inference path of one model. Their FDA-cleared sepsis-prediction clinical decision support system, running in three hospitals' ICUs, is sealed end-to-end. Every prediction, every clinician override, every lab-pull tool call is in the chain. The chain's evidentiary scope is exactly what spec §1.2 names: it proves what the model said at a specific time and that the record was not tampered with after capture. It does NOT prove the model's clinical judgment was right. Patricia knows the §1.2 line; she had it printed on the wall above her desk before the audit team arrived.

Everything else — the EHR (Epic), the billing platform, the claims-processing engine for the insurance arm, the Salesforce-based CRM for member services, the lab results pipeline, the imaging archive — runs on the same plumbing every other healthcare system in the country runs on. Logs in CloudWatch. Notes overwritten in place. Backups on tape. "We trust the engineers." None of those systems has a chain-coverage map of the kind §10.19 requires; the AI-side has one and Patricia wants the audit to confirm the boundary the map draws.

Patricia knows the difference. She picked the AI side first because the FDA SaMD post-market surveillance rule does not give her a choice, the plaintiff bar is paying attention to model-driven clinical errors, and the blast radius was small enough to land in 90 days. She wants Dawn's team to write a real assessment so she can take it to the board and ask for the budget to extend the chain.

This is the second audit Dawn's team has done in three weeks. Last week was a graveyard. The week before that was a different kind of graveyard.

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

The drive in was forty minutes through hospital-district traffic. Dawn had her coffee in the cup holder and the engagement brief on her tablet.

*Last week was a graveyard*, she thought. *The week before that was a different kind of graveyard.*

The week before last had been Northbridge Bank. TesseraSeal everywhere. Every credit decision, every wire transfer, every IAM change, every ETL job — sealed. Dawn's team had spent four days trying to find a gap and the worst thing they had found was a Salesforce-mirror runbook that described connector lag with the phrase "near real-time" instead of the four quantified bounds spec §10.16 requires. That had been a non-conformance the spec explicitly forbids downgrading to a Nit, and Dawn had written it up as Finding-001 against an otherwise spotless engagement. Northbridge had been the cleanest engagement Dawn had run in nine years. The CAE there had asked her, on the last day, whether she had any "wishes" for what they should do next, and the only wish she could think of was that they should write up their internal playbook so other banks could follow it. He had said they were considering exactly that.

Last week had been Continental Mutual — a mid-size financial services firm in the suburbs of a city Dawn tried not to think about on her own time. No chain anywhere. CRM notes overwritten, database backups deletable, CloudWatch logs purgeable by the same engineers who wrote the code. Twelve people with temporary admin that had no expiration date. Dawn had walked out of that one with a forty-page report and a feeling she had not been able to shake on the drive home. The CFO at Continental had asked her, on the last day, whether the report was really going to say what the draft said, and she had told him yes, and he had asked whether they could "soften the language" and she had said no.

Northbridge was still sitting on her like an unanswered question. One §10.16 non-conformance and a clean chain — the cleanest engagement she had run in nine years. The week before Continental, twelve Gaps and four Material Findings on a different bank's report Dawn had reviewed for a colleague. The open question on the drive in was whether Northbridge had been the product or the institution. Mercator was the test.

Today was something else. Today was half-and-half.

*It never is*, she thought. *Except sometimes, on the parts that are.*

She elaborated to herself as she pulled into the visitor lot. The bifurcation was the whole story. The thing that made Mercator interesting was not that they had chained their AI side — plenty of vendors are starting to do that. It was that they had chained *only* their AI side, and they had been honest about it, and they wanted the audit to draw the line clearly so they could fund the rest. That was unusual. Most clients want the auditor to find no gaps. Patricia wanted them found and named so she could go ask for money.

The auditor's instinct, Dawn thought, is to grade. Pass, fail, partial. The Mercator engagement was going to need a different shape. Two grades, side by side. One AI side. One legacy side. Two entirely different posture assessments stapled together. She would need to be careful that the report did not let the AI-side grade dilute the legacy-side grade, or vice versa. Patricia had said "two reports stapled together" and Dawn had agreed — but the stapling itself was going to take some thought. The spec's §10.19 chain-coverage map was the conceptual spine of the bifurcation: the chain reaches one set of systems and not another, and the map names which is which.

She wondered, briefly, what the budget number actually looked like. Patricia had not said. Dawn had a guess based on the surface area — extending the chain to billing, EHR, and the lab pipeline at a system Mercator's size was not going to be cheap. She had seen vendors quote eight figures for engagements like this. Whether the board would approve it was Patricia's problem, not hers. Her problem was making sure the report would carry weight in front of a board that would be looking at the number and asking whether the spend was really necessary.

The conference room was on the fourth floor of the administration tower. Patricia was already there with two of her direct reports and a printed deck. Coffee on the table. A single sheet of paper at each chair: a one-page system map with two colored zones. Green on the left labeled "TesseraSeal scope (sepsis-cds)." Red on the right labeled "Legacy scope (everything else)." No marketing language. Just a map. Dawn recognized the shape immediately — the green zone was a §10.19 chain-coverage map; the red zone named the institutional systems not yet chain-instrumented and the evidentiary substitute (CloudWatch + tape backups + "we trust the engineers") that today substitutes for chain coverage.

"Good morning," Patricia said. "Before we start. I want you to know what you're walking into. Ninety days ago we turned on TesseraSeal for our sepsis CDS model. Three ICUs. One model. One inference path. That is the only thing in this building that is sealed. Everything else looks like every other hospital you have ever audited. I picked the AI side first because the FDA gave me a deadline. I want your assessment to back the budget request to extend the rest. So please find the gaps. That is what I am paying you to do."

Dawn put her coffee down. "Thank you for saying that out loud. It saves us a day."

She had a follow-up.

"You said FFIEC chain-of-custody v1.0b. Walk me through the §1.2 epistemic scope. If the plaintiff bar pulls a sepsis decision into court, what does your expert witness lay foundation on?"

Patricia did not pause. "§1.2 is on the wall above my desk. The chain proves what the model said at a specific time and that the record was not tampered after capture. Two things. The chain does NOT prove the model's clinical statement was accurate, that it complied with our internal AI Governance Committee policy, or that it was free of bias. The §1.1 Daubert grounding lives in the spec itself — testability via the §7 procedure and the test-vector corpus, peer review under the working-group process, known error rate per §1.3 security definitions, general acceptance of HMAC-SHA-256, RFC 6962 Merkle, and Ed25519. A false negative requires simultaneous compromise of the three independent custody layers §1.4 names — IKM, ledger storage, and HSM — plus the §1.2 SDK-process scenario for forward-only forgery. We name the line clearly so the plaintiff's expert cannot drag our witness onto the truth foundation under cross-examination."

Dawn wrote: *§1.1 / §1.2 / §1.3 / §1.4 — Daubert framing is in the spec text, not in vendor marketing.* Same line she had written at Northbridge. Same answer, almost word for word. Different building.

Tom — the visiting team's internal-audit liaison — had already been on a call with Mercator's Chief Audit Executive the day before. He nodded. "We agreed yesterday on the bifurcation framing. The CAE is supportive. He wants the report to read as two reports stapled together."

"Two reports stapled together," Patricia repeated. "Yes. That is exactly right."

Dawn looked around the table at her team. "Okay. Morning is the AI side. Mike and Chen, you're on point — Patricia's team will walk you through the inference chain. Diana, you'll do the IAM split — both sides of the line, because the line is exactly what we want to map. Afternoon is legacy. Raj on databases. Elena on the CRM. Luis on the pipelines. We reconvene at three for the reconciliation test. Five-thirty debrief."

Patricia nodded along. "The sepsis team is expecting Mike and Chen at nine. Dr. Wei — the lead clinical informaticist — will walk you through the inference chain. She built most of it herself, with the platform team. She is very direct."

"Good," Mike said. "I prefer direct."

"You'll like her."

Elena had been listening quietly. "Patricia. Quick orienting question. The Salesforce CRM — is that for the insurance arm or for the physician group?"

"Both. Member services for the insurance arm runs on Salesforce. The physician group's patient outreach team uses the same Salesforce instance as a separate org. Two business lines, one Salesforce footprint. None of it is in the chain. It is exactly the same Salesforce setup we had before any of the AI work started. I want you to look at it because I want it documented in the report alongside everything else."

"That's clear."

Luis, half to himself: "Two business lines on one Salesforce. Member-services CSAs and patient-outreach navigators in the same tenant. That's going to be interesting. No connector under §10.16; no `audit.connector_source.*` family per §4.4.6. It's just Salesforce."

Patricia smiled, just a little. "It is interesting. It is also typical. Hospitals do this everywhere."

She stood up and gestured at the door. "Mike, Chen — let's go to the engineering floor. The sepsis team is expecting you."

---

### 🧩 9:15 AM — First Crack in the Story (or Rather, Not)

Mike had expected the AI walkthrough to follow the usual pattern. Engineering team puts up slides. Talks about "comprehensive logging." Shows a Splunk dashboard. Hand-waves around the parts that cannot be reconstructed.

It did not go that way.

The sepsis-CDS team lead was a clinical informaticist named Dr. Wei, a former ICU attending who had moved into informatics six years ago. She had a laptop open and one terminal window. No deck.

"You want to see a sealed inference," she said. "Pick a date. Pick a hospital."

Mike picked April 12. Memorial campus.

Dr. Wei typed for ten seconds. "Okay. April 12, Memorial. We had 247 sepsis predictions that day across the ICU. Pick one."

"Pick the one with the highest confidence score that fired before 10 AM."

She typed again. The terminal showed an entry — JSON, structured, with a string of fields Mike recognized: `model_id`, `model_version`, `prompt`, `response`, `tool_calls`, `clinician_override`. Each one had a hash next to it. A tenant binding. A sequence number. A signature reference. Mike recognized the shape from the spec — §4.4 attribute table for the chain envelope, plus the OpenTelemetry GenAI attributes that §4.4 makes REQUIRED on any model-call entry.

"This is the entry for prediction `sep-2026-04-12-mem-00194`," Dr. Wei said. "Patient was a 67-year-old admitted overnight with pneumonia. The model was called at 09:42:11. The prompt — meaning the inputs the model saw — included six lab values and four vital signs. Here." She rotated the screen.

```
{
  "entry_id": "sep-2026-04-12-mem-00194",
  "tenant": "mercator-memorial",
  "service": "sepsis-cds",
  "seq": 8472913,
  "ts": "2026-04-12T09:42:11.483Z",
  "model_id": "mercator/sepsis-pred-v3.2",
  "model_version": "v3.2.1-fda-cleared-2025-11",
  "gen_ai.request.model": "mercator/sepsis-pred-v3.2",
  "gen_ai.response.model": "mercator/sepsis-pred-v3.2",
  "prompt": {
    "lab_values": [...],
    "vitals": [...],
    "patient_context_hash": "sha256:..."
  },
  "response": {
    "sepsis_probability": 0.87,
    "confidence": 0.92,
    "feature_attributions": [...]
  },
  "tool_calls": [
    {"tool": "lab_pull", "query": "...", "result_hash": "sha256:..."}
  ],
  "clinician_override": null,
  "audit.deployment.intent": "production",
  "audit.deployment.policy_version": "mercator-mrm-2026q1",
  "hmac": "...",
  "merkle_path": [...],
  "daily_seal_ref": "2026-04-12-d-seal-mercator-memorial"
}
```

Mike looked at it for a long moment. "The clinician override is null. And `gen_ai.request.model` and `gen_ai.response.model` are both populated."

"Right. Spec §4.4 makes both REQUIRED on any chain entry that represents a model call — the request side is what we asked for, the response side is what the model actually answered with. The two are equal here because we run the model in-house, but the spec normates both fields specifically to catch silent vendor-side rerouting where they would diverge. The attending agreed with the prediction and started the sepsis bundle. If she had overridden it, that field would have a separate sealed entry referenced here, with her override reason and her signature."

"And `audit.deployment.intent` is `production`."

"Per spec §4.4.2 deployment-intent capture. Steady-state production traffic. When MRM runs an A/B test or a canary we flip the value to `ab_test` or `canary` for the cohort and the chain captures the per-decision intent classification. Today's run is straight production under MRM policy version `mercator-mrm-2026q1`. The `audit.deployment.policy_version` attribute is conditionally REQUIRED whenever any `audit.deployment.*` attribute is present per §4.4.2. We emit it on every model call so the MRM committee can correlate any change-point with the policy that was in force."

"Show me one with an override."

Dr. Wei filtered. "Plenty to choose from. The model fires more aggressively than our ICU team would like. Override rate is around 30%."

She picked one. The override entry was its own sealed record: clinician ID, override timestamp, override reason (free-text from a structured pick-list), and a hash linking back to the original prediction via `parent_run_id` / `parent_seq` per spec §4.4.

> **✓ Confirmation #1 — Inference entries are content-complete and sealed**
> Every sepsis-CDS inference is captured as a structured entry with the model identifier, the model version (tied to the FDA clearance vintage), the full prompt the model saw, the full response, all tool calls made by the model, and any clinician override. Each entry is HMAC-SHA-256 chained at capture per spec §4.1 with HKDF-derived per-tenant binding. The per-tenant determinism property of §4.1 holds — same IKM and same `tenant_id` produce a byte-identical session key across processes. The session-key handshake conforms to §4.1.1 (Mercator runs Model B — HSM-resident PRK with SDK-side Expand, so the IKM never leaves CloudHSM), the IKM is 32 bytes minimum per §10.6 with generation requirements per §10.6.1 (HSM hardware RNG, never exposed to application memory), and the §10.7 software-key adapter exclusion holds in production (no `plaintext-dev` `kms_handle_uri` ever appears on a production entry). Both `gen_ai.request.model` and `gen_ai.response.model` are populated per the §4.4 SDK-side enforcement rule. Daily Ed25519 seals are signed by AWS CloudHSM in `us-east-1` under FIPS 140-2 Level 3 custody per §10.5, on a daily cadence per §4.2.1, with day-boundary semantics per §4.2.2 anchored to the ledger's `received_at`. The entry shown — `sep-2026-04-12-mem-00194` — recomputed cleanly under `herald-verify` per the §7 verification procedure.

Mike asked Dr. Wei to run the verifier in front of him.

```
$ herald-verify --tenant=mercator-memorial \
                --service=sepsis-cds \
                --date=2026-04-12 \
                --entry-id=sep-2026-04-12-mem-00194 \
                --strict

Status: PASS
Step:   12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key sepsis-prod-2026-q2
```

Twelve verification steps. Pass. Exit code 0 per the §10.12 verifier CLI exit-code contract.

"Re-run it on a randomly chosen entry from last week," Mike said.

Dr. Wei re-ran. Pass.

"Pick one from the day the model was retrained. November 2025."

She paused. "We rotated the signing key the morning of the retraining. Different daily seal key. Should still resolve under the §4.3 algorithm rotation rules and the §10.10 rotation-crossing-the-seal-boundary procedure."

Pass.

> **✓ Confirmation #2 — Verifier resolves across signing-key rotations**
> The chain survives daily seal-key rotation. Each daily seal is bound to its own public key reference (`sepsis-prod-2025-q4`, `sepsis-prod-2026-q1`, etc.) per the §4.2 seal record schema's `public_key_id` field, and the verifier walks rotation boundaries without manual reconciliation. The IKM rotation procedure conforms to §10.10 hourly cadence rules, and the per-entry `key_version` field per §4.4 lets the verifier resolve which IKM generation produced each entry's session key. Constant-time comparison per §10.8 was visible in the verifier source.

Chen had been quiet. She spoke up. "What feeds the prompt? The lab values. Where do those come from in this entry?"

Dr. Wei nodded slowly. "That is the right question. The lab values come from the lab pipeline. The pipeline writes to the warehouse. The model reads from the warehouse at inference time. The `patient_context_hash` field you see in the prompt is a hash of the inputs *as the model saw them*. That hash is sealed. So we can prove what the model received. We cannot prove that what it received matches what was in the warehouse a minute earlier or what was in Epic an hour earlier. The chain stops at our service boundary. The §10.19 chain-coverage map names exactly that boundary."

Chen looked at Mike. Mike looked at Patricia, who had walked in halfway through.

"That is the question," Patricia said. "And that is why I want the budget."

Elena, who had drifted in to observe before her own afternoon session, was watching the screen. "The prompt has a `patient_context_hash`. That hash is sealed. What's in the hash?"

"Lab values, vitals, demographics, and a stable patient identifier — pseudonymized at the boundary. The hash is computed at the model service. Inputs the model used. Sealed. The hash itself is in the canonical bytes the per-event MAC covers per §5; the underlying PHI lives in a separate access-controlled store under HIPAA Security Rule §164.312 controls. The chain entry binds the hash without binding the PHI — that's the §10.22 redaction discipline applied at the SDK boundary, pre-MAC, which is the only conformant posture per §10.22."

"So you can prove the model saw these inputs at this time."

"Yes."

"You cannot prove these inputs are what the lab actually measured."

"Correct."

Elena wrote that down. The phrasing was load-bearing — §1.2 epistemic scope precisely.

Mike kept going for another forty-five minutes. He asked Dr. Wei to show him a model retraining event. Sealed. He asked her to show him a model rollback. Sealed — there had been one in February, when the team had pushed v3.2.0 and rolled back to v3.1.7 within four hours after a feature attribution drift was caught in canary. Both transitions sealed. He asked her to show him a clinician account that had been deactivated. Sealed.

> **✓ Confirmation #3 — Model lifecycle events are in the chain**
> Model deployments, rollbacks, and retraining events are sealed entries with the same structure as inference entries. The February v3.2.0 → v3.1.7 rollback shows as two transition entries (deploy, rollback) with a four-hour separation, signed by the on-call MLOps engineer's hardware key. The canary that surfaced the drift carried `audit.deployment.intent = "canary"` per §4.4.2 with `audit.deployment.canary_traffic_pct` populated, satisfying the §4.4.2 conditional-required rule for canary entries. The chain answered the FDA SaMD post-market question — what version was running when, and what triggered the rollback — from the chain alone, not from email reconstruction.

> **✓ Confirmation #4 — Tool calls within model inference are first-class chain entries**
> When the sepsis model issues a tool call — typically a structured lab-pull query against the warehouse — the tool call, its parameters, and the result hash are recorded inside the parent inference entry. The result hash is bound to the warehouse query response at the time the model received it. The `chain_kind = "tool_call"` discriminator per §3 makes tool-call entries machine-distinguishable from `model_call` and `audit` entries; the §7 verification procedure walks them identically.

Mike asked one more thing before he closed his notebook. "What about the model's own confidence calibration? You said override rate is around 30%. Is that captured anywhere as a chain-coupled metric, or is it derived offline?"

Dr. Wei rotated her screen back. "Both. Each inference entry includes the model's stated confidence. Each override entry includes the structured override reason. We compute calibration weekly — predicted confidence versus override rate, sliced by hospital, ICU unit, and patient demographic — and the calibration report is itself a sealed artifact. Here." She pulled up the most recent one. April 28. The report was a structured document with its own sealed entry, signed by the MLOps lead, and referenced in the MLOps team's monthly drift-monitoring submission to the FDA. It was hash-anchored using the §10.19 `audit.external_artifact.*` family — `audit.external_artifact.kind = "fda_quarterly_calibration_report"`, `audit.external_artifact.sha256 = ...`, `audit.external_artifact.received_at_utc = ...`, `audit.external_artifact.evidentiary_role = "regulatory_compliance"`.

"You file these with FDA?"

"Quarterly. The post-market surveillance package includes the calibration reports plus aggregate adverse-event tracking. Everything in the package is a sealed chain artifact. The package itself has a Merkle root that ties the included artifacts together per the §4.2 RFC 6962 construction, and the root is what FDA receives. If they ask us to reproduce any single artifact, we re-pull from the chain and verify."

Mike: "And FDA understands the chain."

"FDA understands cryptographic integrity in principle. They are still learning what to do with it operationally. We submitted our first quarterly package in February. The reviewer's first question was whether they could verify the chain themselves. We sent them the verifier binary and the public-key references — Apache 2.0 verifier per the §10.26 reference verifier distribution discipline, with reproducible builds, Cosign signatures, and the per-platform binary the reviewer's laptop needed. They came back two weeks later with verifier output that matched ours, exit code 0 per §10.12."

"They ran the verifier."

"They ran the verifier. We were the second SaMD they had seen do this. The first was a different vendor's diagnostic imaging tool."

> **✓ Confirmation #X (extra) — FDA reviewer-side verification works in practice**
> Mercator's quarterly post-market surveillance submission to FDA includes a Merkle root tying together calibration reports, adverse-event records, and inference-entry samples. FDA reviewers ran the verifier independently in February. Verifier output matched. The chain is not a private artifact; it works at the regulator boundary. The §10.26 reference verifier distribution discipline — repository separation, reproducible builds, Cosign-signed binaries, SHA-256/SHA-512 manifests, SBOM — is what made the cross-organization handoff possible without trusting Mercator-side credentials.

Mike closed his notebook for the morning side. "I have not seen a chain this clean outside Northbridge."

Dr. Wei smiled. "We hired the engineer who built theirs."

Mike looked at her for a long beat. He filed the comment for later.

---

### 🧠 10:00 AM — Database Deep Dive (AI Side)

Raj joined Chen and Mike on the engineering floor. The morning database review was the warehouse — specifically the read path the sepsis model used at inference time. The afternoon would be everything else.

The warehouse engineer, a soft-spoken person named Devansh, walked Raj through the schema. The inference-relevant tables were a separate schema — `sepsis_inference_inputs` — with append-only semantics enforced at the database level via a row-versioning trigger. Every update created a new version row; the application never saw an in-place update. Raj recognized the discipline — it was the warehouse-side parallel of the §10.3 append-only enforcement the chain itself runs at the storage tier.

"That's not the chain," Raj said.

"No. The chain is one layer up. This is what feeds it. We made the warehouse read append-only on the inference path so we could checksum what the model saw against what we have stored. The chain entry has the hash of what the model saw. The warehouse has the row versions. We can join them after the fact."

Raj asked for a sample. Devansh pulled the warehouse rows that had fed the same `sep-2026-04-12-mem-00194` prediction. The lab values — six rows — were each in `sepsis_inference_inputs` with their version_id and a `served_at` timestamp. The hash of the served set, computed on the fly, matched the `patient_context_hash` in the chain entry.

> **✓ Confirmation #5 — Inference-input warehouse is reconcilable to the chain**
> The warehouse tables that feed the sepsis model are append-only with row versioning. The sealed chain entry's `patient_context_hash` recomputes from the warehouse rows the model read. For sepsis predictions in the last 90 days, this reconciliation works backward from any inference entry to the exact warehouse state the model saw. Reconciliation depth: the warehouse boundary. Beyond that — the lab pipeline that filled the warehouse — the chain ends. The §10.19 chain-coverage map names the warehouse as a chain-instrumented institutional system; the lab pipeline upstream of it lands in the §10.19 "institutional systems not yet chain-instrumented" category with the institution's CC8.1 naming the rollout posture as deferred.

Raj wrote the boundary observation in his notes. He underlined the second half.

He asked Devansh: "What about the 90+ day window? You said the inference pipeline went live ninety days ago. What's in the warehouse from before that?"

"Same data structure. Same tables. Just no chain entries pointing at it. We backfilled the warehouse from the EHR archive when we set up the pipeline, because we needed historical training-set parity. But there's no inference chain referencing pre-launch data. The chain starts at go-live."

"And the EHR archive itself?"

"Epic. That's the legacy side."

Raj wrote that down too.

He had one more line of inquiry before he moved on. "Devansh. If the inference-input warehouse is append-only and the pre-launch backfill came from Epic — what guards against the backfill itself having been corrupted on the way in?"

Devansh: "Honest answer, nothing structural. We did a one-time validation against Epic at backfill time. We sampled 5,000 patient-encounter records and reconciled. Match rate was 100%. The validation report is a chain entry — anchored under the §10.19 `audit.external_artifact.*` family with `kind = "ehr_backfill_validation_report"`. But the backfill itself wasn't streamed through the chain — that would have required the chain to exist before go-live, and it didn't."

"So the backfilled rows are append-only now, but their original loading wasn't sealed."

"Correct."

"Does the model use backfilled rows in production inference?"

Devansh thought. "It uses the warehouse as a feature store at inference time. Some features are time-series — say, lab values over the last 72 hours. If a patient's lab value from 70 hours ago happens to be a pre-launch backfilled row, then yes, the model is reading a row whose original loading wasn't sealed. After 72 hours from go-live that risk goes away. It's gone for new patients. It is theoretically still possible for a patient who was admitted before go-live and has a long ICU stay."

Raj: "How long was the longest ICU stay overlapping the go-live boundary?"

Devansh checked. "Eighteen days. So for the first eighteen days post-launch, the model could have read backfilled rows for patients with stays spanning the boundary. After that, no."

"That's a known gap."

"Yes. And documented. We added it to the FDA submission as a known limitation of the post-launch surveillance window. FDA accepted it."

Raj wrote it down. He marked it with a small annotation — "addressed at submission, not a finding for this engagement, note for completeness." He cited spec §10.19 in the margin: chain-coverage map enumerates exactly this kind of bounded-period evidentiary substitute.

> **✓ Confirmation #X (extra-2) — Pre-launch backfill is documented as an FDA-acknowledged limitation**
> The 90-day inference chain begins at go-live. Patients whose ICU stay spanned the go-live boundary may have been served features from pre-launch backfilled warehouse rows during the first 18 days post-launch. Mercator disclosed this in the FDA SaMD post-market submission as a known limitation, hash-anchored to the chain via the §10.19 `audit.external_artifact.*` family. The submission was accepted. This is not a finding for this engagement; it is noted for completeness. The §10.20 training-data retention floor discipline applies in spirit — the deployment window is bounded, the backfill provenance is bounded, and Mercator's CC8.1 names both bounds.

He closed his notebook for the AI-side database review. The afternoon was going to be a different kind of work.

---

### 🔐 11:00 AM — IAM Review

Diana had been doing this for fifteen years and she could spot a bifurcated IAM environment in about forty minutes. Mercator's was bifurcated to the millimeter.

The sepsis service had its own service-account namespace. Every credential issued to the sepsis pipeline — the model serving runtime, the warehouse reader, the CloudHSM signer — was registered through an internal IAM workflow that emitted a chain entry on every change. Credential issued: chain entry. Credential rotated: chain entry. Credential revoked: chain entry. Credential's permissions changed: chain entry. The shape was the IAM-as-chain pattern Diana had seen at Northbridge two weeks ago — every IAM lifecycle event lands as an `chain_kind = "operational"` entry per the §3 enumeration. Diana asked to see the entry for the most recent rotation. Patricia's IAM lead pulled it up.

```
{
  "entry_id": "iam-2026-05-01-sepsis-cds-rotate-001",
  "tenant": "mercator-platform",
  "service": "iam-sepsis",
  "actor": "iam-automation-service",
  "approver": "p.okonkwo@mercator.health",
  "action": "credential_rotate",
  "subject": "svc-sepsis-cds-warehouse-reader",
  "old_cred_fingerprint": "sha256:...",
  "new_cred_fingerprint": "sha256:...",
  "ts": "2026-05-01T03:00:00Z",
  "approval_ticket": "IAM-2026-3417",
  "ffiec.chain.chain_kind": "operational",
  ...
}
```

Diana asked: "Who can change the IAM rotation policy itself?"

"Three people. Patricia. Me. The platform engineering director. Any change to the policy emits a chain entry. The change requires two-of-three approval."

> **✓ Confirmation #6 — Sepsis service IAM is chain-coupled end to end**
> Every credential lifecycle event — issuance, rotation, revocation, scope change — for the sepsis service emits a sealed chain entry under `chain_kind = "operational"` per the §3 enumeration. Policy-level changes (who can rotate, who can approve, what the rotation interval is) require two-of-three approval and are themselves chain entries. There is no path to change a sepsis service credential without producing a chain record. The append-only enforcement at the storage tier per §10.3 means even a privileged operator cannot retroactively delete the rotation history.

Diana asked, almost out of habit: "And the clinician users? The doctors who get the alerts?"

The IAM lead's expression shifted. "That's legacy AD."

"Walk me through it."

He did. Mercator's clinician identity is in Active Directory, federated to the EHR via SAML, federated again to the sepsis-CDS UI via the same SAML pipeline. When a sepsis prediction fires and the attending overrides it, the override entry's `clinician_id` is the SAML subject — which is stable — but the underlying AD account is managed by the central IT identity team using the same processes Mercator has used for fifteen years.

"How many people have temporary admin in the AD environment?"

The IAM lead checked. "Right now? Twenty-three. Fourteen are over thirty days old. Six are over a year. The oldest is from 2019."

Diana wrote that down very carefully.

> **⚠️ Surprise #1 — Clinician identity (legacy AD) has no chain and has the same temporary-admin sprawl as every other AD environment**
> Twenty-three temporary-admin grants are currently active in the clinician AD. Fourteen are older than 30 days. Six older than a year. The oldest from 2019. There is no audit trail of what these admins did beyond the AD security log, which retains 90 days and can be cleared by Domain Admins. The sepsis chain captures *which clinician* signed off on an override (the SAML subject is stable). It does not capture *whether that clinician's account was acting under a privilege grant that should have expired four years ago*. The remediation pattern is the IAM-as-chain shape Mercator already runs on the sepsis service side — every grant a chain entry, every revocation a chain entry, time-boxed expiration enforced by a chain-driven worker rather than by a cron job an admin can disable. The sepsis side runs that pattern; the legacy AD does not. Spec §10.19 chain-coverage-map placement: legacy AD is "institutional systems not yet chain-instrumented"; the §10.19 evidentiary-substitute column reads "AD security log, 90-day retention, mutable by Domain Admins."

Diana paused. "If the AD account that signed an override was actually being shadowed by a temporary admin who was impersonating, the chain would still show the SAML subject of the impersonated user."

"Yes. The chain trusts the upstream identity provider."

"And the upstream identity provider is the legacy AD."

"Yes."

She wrote that down twice, in two different places in her notebook.

> **⚠️ Surprise #2 — Service-account chain rests on a legacy-AD root of trust for clinician overrides**
> The sepsis chain sees what the SAML subject says it sees. The legacy AD is the root of the SAML subject for clinician overrides. A compromise of the AD identity layer compromises the clinician-attribution side of the sepsis chain even though the chain itself is intact. The chain proves *what the system recorded*; the legacy AD is what determines *whether the recorded actor is who they appear to be*. This is an §1.2 epistemic-scope boundary playing out at a layer the §1.4 compositional security argument does NOT cover — §1.4 names IKM, ledger storage, and HSM as the three custody layers; the upstream identity provider is a fourth trust assumption that lives outside the chain's cryptographic substrate. A taller wall on the AD foundation does not fix the foundation. The remediation is to firm up the AD identity root before extending the chain further; the §10.21 cross-vendor model-handover discipline doesn't apply directly here, but the spirit (bind cryptographic evidence at every handoff) does.

Diana made a note that the AI-side IAM finding was clean *within its boundary*, and the legacy-side IAM finding was as bad as Mercator's peer institutions, and the relationship between the two was the actually-interesting finding.

---

### 🧪 12:00 PM — Lunch (But Not Really)

The team gathered in a side conference room with sandwiches from the cafeteria. Tom was on the phone with the CAE. Patricia had ducked out to a board prep call.

Dawn put her sandwich down before she'd taken a bite. "Let's talk about the morning."

Mike: "AI side is real. I've audited four AI deployments this year that claimed sealed inference. Mercator's is the second one that actually does it. The §4.4 attribute table is populated correctly — both `gen_ai.request.model` and `gen_ai.response.model` on every model-call entry, `audit.deployment.intent` and `audit.deployment.policy_version` on every entry, `chain_kind` populated correctly per the §3 enumeration." (The first one Mike had seen was Northbridge's quant trading model the week before last.)

Chen: "The pipeline-to-chain reconciliation works. The warehouse is append-only on the inference path. The chain entry's input hash recomputes from the warehouse. That's a real reconciliation, not a check-the-box log. Backward to the chain's input boundary, the §10.19 map names the boundary correctly."

Diana: "Service IAM is sealed per §10.3 and the §3 operational `chain_kind`. Clinician IAM is legacy. Any clinician override is chain-recorded but only as deep as the SAML subject is trusted. The legacy AD is the root of trust for the chain on the clinician side — that's a §1.2 epistemic-scope boundary outside the §1.4 three-layer compositional security argument."

Raj: "I haven't done the full warehouse review yet — that's after lunch — but the inference-input schema is properly designed. Append-only with versioning."

Dawn: "Okay. Question. Is half-chained better than not chained at all?"

The room got quiet for a beat.

Tom, off the phone: "That's the question Patricia is going to ask at five-thirty."

Mike: "Yes. It's better. You can prove what the model saw and what it said. That's the FDA's biggest concern. Half a chain is a real chain on the half it covers. §1.2 epistemic scope is exactly that line — the chain proves what the model said at time T, no more and no less."

Diana: "Better, yes. Sufficient, no. Half a chain creates a false sense of completeness. People look at the AI side, see how clean it is, and assume the rest is at least decent."

Chen: "And the lab pipeline that feeds the model is on the wrong side of the line. If a lab value is wrong before it gets to the warehouse, the chain seals the wrongness. The chain proves what the model received; it doesn't prove what the lab actually measured. §1.2 again."

Dawn tapped her pen. "Better, yes. Sufficient, no. That's my position. Tom?"

"That's the position the CAE is bracing for. He told me so this morning."

Elena: "I want to flag — I haven't been to the CRM yet, but if the morning is any indication, the CRM is going to be the diary all over again. Salesforce. Member services. Notes. Backups. No `audit.connector_source.*` family per §4.4.6 because there's no connector."

Dawn: "Yes. We'll see at one."

Luis had been quiet. He looked up from his laptop. "I was reading their internal runbook for the lab pipeline while you were all talking. The pipeline writes to S3. CloudTrail logs the writes. CloudTrail logging can be disabled per-bucket by anyone with `s3:PutBucketLogging` permission. Three engineers have that permission. None of them are on a watchlist. Compare that to the §10.5 HSM custody and §10.3 append-only enforcement on the AI side — the chain's storage tier on the sepsis side is in object lock with compliance-mode retention; the lab S3 bucket is just a bucket."

Dawn put her sandwich down again, after one bite.

"Better, yes. Sufficient, no," she said again. To no one in particular.

Tom, off the phone again: "The CAE asked me to put a question to the team. He wants to know how we frame this in the report so the AI-side findings don't look like a marketing piece for TesseraSeal."

Dawn thought about that for a moment. "Frame it as boundary-marking. The AI side is sealed because Mercator chose to seal it; the legacy side is not because Mercator hasn't gotten there yet. The report is not endorsing TesseraSeal. The report is documenting that on the surface area where Mercator has applied chain-grade controls, the controls hold; on the surface area where they haven't, they don't. The vendor is incidental to the finding. Spec §10.19 is the framework — a chain-coverage map, with categories, and an evidentiary substitute named for each unchained system."

Mike: "I'd add — the AI side is sealed *to the controls Mercator built*. The vendor product is the substrate. The sealing posture is Mercator's. The same product deployed without two-of-three approval, without append-only warehouse semantics, without the reconciliation events Chen showed us, would not have produced this audit posture. The product is a tool. The posture is the team. Same Vidimus SDK, same Herald Core ledger, same `herald-verify` CLI — the difference between Mercator's posture and a hypothetical sloppy deployment is the institutional discipline at the SDK and operational layer."

"That's the framing," Dawn said. "Tom, tell the CAE we'll make sure the report reads that way."

Tom relayed the message. He came off the call after another minute. "He's on board. He also asked — between us — whether he should brace the board for a number with a B in it."

Dawn looked at him.

"That's his question, not mine," Tom said.

"That's a question for Patricia," Dawn said. "Not for us. Our job is the assessment."

The team finished lunch. Elena had already started on the CRM walkthrough — she'd wandered off at 12:25 with Jordan-the-CRM-admin's calendar invite on her laptop. The rest of them rinsed coffee cups and walked back out to the engineering floor.

---

### 🧪 12:30 PM — CRM Walkthrough (Elena's Detour)

Elena had a 45-minute slot before the formal afternoon kicked off. She used it on the Salesforce CRM, because she'd asked Patricia at kickoff and Patricia had said "go now if you want, you'll lose half the day to the rest of the legacy stack."

The CRM admin was a person named Jordan. Friendly. Senior. Twelve years at Mercator. They had seen four CRM platforms come and go. The current Salesforce instance had been live since 2019.

Elena: "Walk me through how a member services call gets recorded."

Jordan: "Member calls in. CSA picks up. Authenticates the member with the standard challenge questions. Pulls up their record. Logs the call as a Case object. Adds notes to the Case. If escalation is needed, the Case is reassigned. If resolved, the Case is closed. Standard Salesforce flow."

"Notes on the Case — are those Case Comments or are they on the Case description field?"

"Both, depending on the CSA. Some of them put everything in the description. Some of them use Comments. We don't enforce a structure."

"If a CSA edits a Case description after closing it — is that visible?"

Jordan paused. "Salesforce field history. If field history is enabled for the Description field. Let me check."

They checked. Field history was enabled for some fields, not for Description. The Description field was the one most CSAs put their notes in.

"So a CSA can edit the Description after the call and there's no record."

"There's the Salesforce setup audit log, which would catch the field-history-disabled config decision. But the edits themselves — no, not visible to the audit log if field history isn't on."

"Why isn't field history on for Description?"

Jordan thought. "Storage cost. Salesforce charges per field-history row at scale. We turned it off for high-volume free-text fields years ago to control the bill. I think that's why. I'd have to check with the platform team."

Elena had heard this answer before. Last week, in fact, in a different building, from a different person.

> **⚠️ Surprise #X (CRM-1) — Salesforce Description field history disabled for storage-cost reasons**
> The primary free-text field for member-services Case notes does not have Salesforce field history enabled. The decision was made years ago to control storage bill. Edits to Case Description after the call cannot be reconstructed from Salesforce metadata. There is no chain. There is no out-of-band log. The story is "we keep backups," and the backups are on a 30-day rolling window before they go to long-term storage that the CRM platform team can manage. Compare to §10.3 append-only enforcement at the chain's storage tier (compliance-mode retention, no role with delete permission inside the window) and to §4.4.6 SaaS-edge connector source attribution (which would let a chain-instrumented CRM mirror prove what Salesforce committed at what time, with a stable source-side `run_id`). Today the CRM has neither. The §10.19 chain-coverage map places this CRM in the "institutional systems not yet chain-instrumented" category with the evidentiary-substitute column reading "Salesforce field history (selectively enabled), 30-day backups."

"What about the patient-outreach navigators? Same Salesforce instance, separate org?"

"Same instance. Separate Salesforce org-within-the-tenant. Same field-history posture. Same retention. The navigators do clinical outreach for the physician group — care management for high-risk patients. They put a lot of clinical context in the Case notes. PHI by definition."

Elena set her pen down. "PHI in a free-text field with no edit audit and a 30-day backup window."

"Yes."

"Has that ever come up in a HIPAA review?"

"It's been flagged as a Partial finding in the HITRUST attestation for the last three cycles. The remediation plan has been the same for three cycles."

"What's the remediation plan?"

"Enable field history. Increase retention. Estimate from Salesforce was around six figures annually for the storage uplift. Each cycle the budget didn't land. Each cycle the Partial finding got renewed."

Elena wrote that down. Slowly.

> **⚠️ Surprise #X (CRM-2) — Patient-outreach navigators store PHI in unaudited free-text Salesforce field**
> The physician group's care-management navigators use Salesforce Cases for high-risk-patient outreach. Clinical context — including PHI — goes into the Description field. Field history is disabled. The HITRUST Partial finding has been renewed for three cycles. The remediation has not been funded. The CRM is functionally identical to the CRM at last week's financial-services audit, except that the data sitting in it is medical rather than financial. Under the §10.13 evidentiary-artifacts retention guidance applied with HIPAA 6-year-floor controls per 45 CFR 164.530(j)(2) and the state pediatric-records floors that often run to 7-10 years past age 18, the 30-day backup window is not even close to the operative retention floor. The §10.22 redaction discipline does not apply here because there is no chain to redact pre-MAC into; the PHI lives in the source field with no integrity binding at all. If Mercator extends the chain to the CRM via the §4.4.6 connector pattern, the institution's CC8.1 will need to name the lag bound under §10.16's four quantified numbers (median, 95th-percentile SLO, alerting threshold, RTO) — descriptive language like "near real-time" is non-conformant per §10.16's severity-classification clause and would land as a non-conformance Dawn could not downgrade.

She closed her notebook on the CRM. "Jordan. Thank you. I appreciate the directness."

"Of course. I've been waiting for someone to ask the right questions about Description. The platform team is going to be happy you flagged it. They've been trying to fund the remediation for years."

Elena walked back to the conference room. *Of course they have*, she thought. *And they'll be trying to fund it next year, too, unless this report changes something.*

---

### 🔄 1:00 PM — API Layer Inspection

Mike's afternoon was the API gateway. Mercator runs the sepsis-CDS API behind an internal AWS API Gateway with a custom Lambda authorizer and a request-signing layer. Every call to the sepsis API — from the EHR's CDS hook, from the bedside monitor integration, from the clinician override UI — goes through this gateway.

Every call emits a chain entry. Request, headers (filtered for PHI per §10.22 redaction discipline applied pre-MAC at the SDK boundary), authorizer decision, downstream service, response code, response hash, latency. Sealed. The OTLP transport carries the §4.4.3 Resource attributes — `ffiec.chain.spec`, `service.name`, `service.version`, `ffiec.chain.posture = "ffiec"`, `ffiec.chain.format_version` — so the receiver dispatches correctly before per-entry decode. The §4.4.4 severity discipline keeps chain traffic from being filtered by the OTLP collector pipeline.

"Show me a call that returned a 5xx," Mike said.

The API engineer pulled up an entry from April 18. 503. The authorizer had failed because CloudHSM was briefly unreachable during a planned maintenance window. The chain entry recorded the failure. The downstream sepsis service had served a fallback "model unavailable" response. The fallback was itself sealed as a separate inference-entry-shaped record with `model_id: fallback-noop`, `response: unavailable`, and a reference to the upstream API failure. The HSM unavailability event had been emitted as an `chain_kind = "operational"` event per §4.3.1 HSM unavailability and notification.

> **✓ Confirmation #7 — API gateway and fallback paths are sealed**
> Every sepsis API call — including failed authorizations, planned-maintenance outages, and fallback responses — produces a sealed chain entry. The April 18 CloudHSM maintenance window shows up cleanly: 247 sepsis API calls during the window, all routed to the fallback responder, all sealed. No silent gaps. The §4.3.1 HSM unavailability event captured the maintenance start and end; the fallback responder's chain entries link to the operational event via `parent_run_id` / `parent_seq` per §4.4. A regulator inquiring "what was your model doing during the CloudHSM outage" answers from the chain alone — 247 fallback responses, none of which advised the bundle, with the operational context preserved.

Mike turned to Mercator's lead API engineer. "What about the Epic-to-billing API? The one that sends the inpatient discharge data over to the billing platform when a patient leaves?"

The API engineer looked at Patricia, who had returned from her board prep call.

Patricia: "That is the legacy API gateway. Different team. Different vendor. No chain. It runs on the same Mulesoft stack we've had since 2016."

"How are calls audited?"

"Mulesoft has an audit log. It stores six weeks. After six weeks, calls older than six weeks are aggregated into a daily summary and the per-call records are deleted to control storage costs."

"Per-call records are deleted by the system or by a person?"

"By a scheduled job. The job runs as a service account that the Mulesoft team manages. The team has eight people. Any of them can change the retention from six weeks to anything else. There's no chain on that change either."

Mike wrote that down.

> **⚠️ Surprise #3 — Epic-to-billing handoff has no chain and a vendor-managed retention dial**
> The legacy Mulesoft gateway that brokers Epic ↔ billing calls keeps per-call records for six weeks. After six weeks, records are aggregated and originals deleted by a scheduled job. The retention dial is a service-account-controlled value that any of the eight Mulesoft engineers can change. There is no record of changes to the retention dial. There is no chain on any of this. Compared to the chain's §10.13 evidentiary-artifacts retention guidance and the HIPAA 6-year floor at 45 CFR 164.530(j)(2), the 6-week window is dramatically below what the regulatory floor demands. The §10.3 append-only enforcement that protects the sepsis chain's storage tier has no parallel here — eight engineers can shorten retention silently. Under spec §10.19 the chain-coverage map places this in "institutional systems not yet chain-instrumented" with "Mulesoft 6-week audit log, mutable retention dial" as the evidentiary substitute. CMS-0057-F prior-authorization-timeliness audits will not be answerable from this Mulesoft log once decisions roll outside the 6-week window — and CMS retention expectations run far longer than that.

He kept going. The claims-processing API for Mercator's insurance arm was on a third stack. AWS, but built on a 2018-era reference architecture with Kinesis and Lambda and a DynamoDB audit table that the SREs had unilateral access to. Not chain-coupled. Standard CloudWatch retention.

> **⚠️ Surprise #4 — Claims-processing audit table is mutable by SREs**
> The DynamoDB table that audits claims-processing decisions is writable and deletable by the SRE on-call. Five SREs have the relevant IAM. There is no out-of-band signing layer. Standard CloudWatch retention with a 90-day TTL set by a SRE-controlled CloudFormation template. The §10.5 HSM custody discipline (separation of duties, FIPS 140-2 Level 3, signing key non-extractable) and the §10.3 append-only storage discipline — both of which protect the sepsis chain — have no parallel here. Five SREs hold the keys to the audit-evidence kingdom for claims processing, and a careless or hostile actor could erase the trail. Where the insurance arm's claims processing makes adverse-determination decisions on health-coverage benefits, the absent chain trail also leaves the institution unable to reconstruct decision timing the way §10.11.2's FCRA reinvestigation timing schema reconstructs ECOA / FCRA dispute timing — the spec's `audit.fcra.reinvestigation.*` family does not directly bind to insurance adverse-determination workflows but the discipline (a chain entry per dispute-clock-starting event, integrity-bound timing) is exactly the pattern that would protect Mercator under state-insurance-privacy regimes (NY DFS Insurance Circular Letter 7, Cal Insurance Code §791.10) by analogy under §10.11's translation discipline.

Mike sat back in his chair. He'd heard this story before. About a week ago. In a different building, in a different industry, with different acronyms.

---

### 🧬 2:00 PM — Data Pipeline Reality

Chen was sitting next to Luis at a long desk, two laptops open. The morning had been the inference pipeline. The afternoon was every other pipeline.

The lab results pipeline was first. Lab analyzers in the three hospitals' labs send results over an HL7 v2 interface to an integration engine. The integration engine writes to two places: Epic (for clinical use) and an S3 bucket (for downstream analytics, including the warehouse the sepsis model reads).

"How do you know the value Epic sees is the same value the warehouse sees?" Chen asked.

The lab pipeline lead, a person named Marcus, scratched his head. "We don't, really. They come from the same HL7 message. The integration engine fans out. There's no checksum reconciliation between the two destinations."

"What if the integration engine writes successfully to Epic but fails on S3?"

"It retries. If retries exhaust, it writes to a dead-letter queue. The DLQ is monitored."

"By who?"

"The integration team. Three people. They have a Slack channel."

"Does the DLQ persist?"

"Three days, then it's purged."

> **⚠️ Surprise #5 — Lab pipeline has no cross-destination reconciliation and a 3-day DLQ**
> The HL7 integration engine fans the same lab message to Epic and to the analytics S3 bucket without a reconciliation checksum. If the S3 write fails and Epic succeeds, the lab value the sepsis model eventually sees may differ from what Epic shows the clinician. The DLQ retains three days. After three days, evidence of the divergence is gone. Compare to the sepsis ETL's cross-checkpoint sealed reconciliation events (Confirmation #8 below) and to the §4.4.6 SaaS-edge connector pattern's source-commit binding (`audit.connector_source.commit_timestamp`, `audit.connector_source.replay_id`, `audit.connector_source.lag_observed_ms`). A reconciliation-checksum chain entry pattern modeled on the §10.16 `connector.lag_observation` operational event would close the gap; today there is nothing. The §10.19 chain-coverage map places this as "institutional systems not yet chain-instrumented" with "DLQ 3-day retention, monitored by Slack" as the evidentiary substitute. The 3-day window is dramatically below the HIPAA 6-year retention floor at 45 CFR 164.530(j)(2) for clinical data the lab values eventually inform.

Chen kept going. Luis had pulled up the S3 bucket policy on his laptop in parallel.

"Luis. What's the bucket policy?"

"CloudTrail logging is enabled. The CloudTrail target bucket is in the same account. The same three engineers who manage the pipeline have `s3:PutBucketLogging` and `cloudtrail:StopLogging`. Logging can be disabled. Lambda triggers can be disabled. Nothing chain-couples the pipeline state."

> **⚠️ Surprise #6 — Lab pipeline storage layer is mutable and CloudTrail-disable-able**
> The S3 bucket the lab pipeline writes to has its CloudTrail logging configured by the same engineers who run the pipeline. Three engineers can disable logging. Three engineers can edit Lambda triggers. There is no out-of-band sealing. The lab values that ultimately become the sepsis model's prompt inputs sit in a fully mutable storage layer for the entire window between lab analyzer write and warehouse ingest. The §10.3 append-only enforcement that protects the chain's storage tier (object lock with compliance-mode retention, no role with delete permission inside the retention window) has no parallel on this S3 bucket. Three engineers run the pipeline; three engineers can turn off the pipeline's own audit trail. The §1.4 compositional-security argument — three independent custody layers required for a verifying tamper — collapses here because there is only one custody layer (the engineers) and they hold every relevant permission.

Chen looked at Luis. Luis looked back.

"This is the diary," Luis said.

"This is the diary," Chen agreed.

Patricia, who had stopped by to listen, said nothing. She just nodded.

The imaging archive — PACS — was next. PACS had its own audit log. Standard DICOM auditing. Centrally controlled by the radiology IT team. Six engineers. Same pattern.

The claims-processing engine for the insurance arm — already named in Mike's API review — turned out to also write to a separate analytics warehouse with its own ETL pipeline. That pipeline had checksums between source and destination, but the checksum log was in a relational table the data engineering team could update. Chen found three rows in that table that had been edited within the last sixty days, with no edit history.

> **⚠️ Surprise #7 — Claims ETL checksum log is editable**
> The reconciliation checksums for the claims-processing ETL are stored in a Postgres table the data engineering team can `UPDATE`. Three rows show edits in the last sixty days. The previous values cannot be reconstructed; the table has no audit trigger. The reconciliation evidence the team relies on for SOC 1 attestation is mutable by the team that owns the attestation. Compare to the sepsis ETL pattern (Confirmation #8 below) — every reconciliation event is itself a sealed chain entry under `chain_kind = "operational"` per §3, integrity-bound under §4.1's per-event MAC and the §4.2 daily Merkle seal. The claims ETL's mutable Postgres table is exactly the §10.23 consumer-correlation-index integrity gap the spec closes for CFPB-class production: an unchained, institution-controlled retrieval substrate is exactly the kind of asymmetric-evidence move §10.23 names; the same anti-pattern exists here for SOC 1 attestation. Either Shape 1 (chain-anchored, every checksum a chain entry) or Shape 2 (daily attestation event under §10.2) would close it.

Chen wrote that down with a circle around it.

She came back to the sepsis pipeline at the end of the hour, almost as relief.

> **✓ Confirmation #8 — Sepsis ETL has cross-checkpoint reconciliation as sealed events**
> Each stage of the sepsis inference pipeline emits a sealed reconciliation event with the input checksum, the output checksum, and the row counts. The reconciliation events are themselves chain entries under `chain_kind = "operational"` per the §3 enumeration. A divergence at any stage produces a sealed alert and a halt-the-line response. In the last 90 days, three reconciliation events have triggered halts. All three were investigated, resolved, and re-run with sealed before/after evidence. The append-only enforcement at §10.3 means the halt evidence cannot be erased even by an SRE with elevated permissions; the §10.5 HSM custody on the daily seal means the institution itself cannot retroactively rewrite the day's reconciliation history.

She closed her notebook on the pipeline review. "The sepsis side has an integrity story. The other pipelines have a customer story."

---

### 📊 3:00 PM — Reconciliation Test

The reconciliation test was Dawn's design. She had sketched it the night before. The premise was simple: pick five sepsis alerts. Trace each one end-to-end — backward to the lab values that fed it, forward to the clinical action that followed. Score what reconciles and what does not.

The team gathered in a small huddle room with a screen on the wall. Mike, Chen, Diana, Raj, Dawn. Patricia was there. Dr. Wei was there. The CAE was on a Zoom.

Mike picked the alerts at random — five from a single week, spread across the three hospitals.

For each alert, the same three-part trace.

**Alert 1.** `sep-2026-04-09-uni-00081`. University campus. 04:15 AM. Sepsis probability 0.91. No override.

- AI inference: chain entry resolves under `herald-verify` per the §7 verification procedure. PASS, twelve steps. Exit code 0 per §10.12.
- Backward to lab values: warehouse rows present, hash matches the chain entry's `patient_context_hash`. Lab pipeline upstream of the warehouse — six lab values came in via HL7 the previous evening at 22:14, 22:15, and 22:16. S3 raw landing data is present. CloudTrail records the writes. *No chain, but the data is there.* The §10.19 chain-coverage map names exactly this hand-off boundary; the data presence is the institution's evidentiary substitute and Mercator was honest about that.
- Forward to clinical action: EHR shows the sepsis bundle was started at 04:18, three minutes after the alert. The order set author, the medication administration, and the nursing documentation are all in Epic. The order set has not been edited since. The nursing note has not been edited since. PASS.

**Alert 2.** `sep-2026-04-10-mem-00102`. Memorial. 14:40 PM. Sepsis probability 0.84. Clinician override at 14:42 — attending overrode citing "known UTI with localized infection, low concern for sepsis."

- AI inference: PASS per §7.
- Backward to lab values: warehouse rows present, hash matches. *No chain on the lab pipeline upstream, but the data is there.*
- Forward to clinical action: override is sealed (a separate chain entry of the override-record shape, linked via §4.4 `parent_run_id` / `parent_seq` to the prediction). Attending's reasoning is in the override entry. Subsequent EHR documentation: the attending added a progress note at 16:15 that referenced the override decision. The progress note has been *edited* twice since — once at 18:30 the same day, once on April 14. The edits are visible in Epic's clinical-note-version history. The edits are not in any chain.

**Alert 3.** `sep-2026-04-11-uni-00045`. University. 02:33 AM. Sepsis probability 0.78. No override.

- AI inference: PASS per §7.
- Backward to lab values: warehouse rows present, hash matches. **One lab value is missing from the warehouse.** A potassium level the model used in its prompt is no longer in `sepsis_inference_inputs`. The chain entry shows the model received a value of 5.8 mEq/L. The warehouse has no row for it. Cross-checking with the lab pipeline S3 raw data: the original HL7 message is also gone — purged 60 days ago by a retention job that has since been disabled but ran on this date. The chain entry retains the *hash* of what the model saw. The plaintext value (5.8) is in the chain entry's prompt field, sealed. But the *underlying source record* is gone. The chain says "the model saw 5.8" with §1.2 epistemic certainty. The source no longer says anything. This is a §1.2 epistemic-scope-boundary case in operation: the chain can answer "what did the model see at time T" but cannot independently corroborate "5.8 was the actual lab value" without the source record.
- Forward to clinical action: bundle started, documented in Epic. Note unedited.

**Alert 4.** `sep-2026-04-12-mem-00194`. Memorial. 09:42 AM. (The same alert from this morning's walkthrough.)

- AI inference: PASS.
- Backward to lab values: PASS in warehouse. Lab pipeline S3 data present and unmodified.
- Forward to clinical action: bundle started, documented in Epic, note unedited. PASS.

**Alert 5.** `sep-2026-04-13-uni-00211`. University. 11:08 AM. Sepsis probability 0.94. Clinician override.

- AI inference: PASS.
- Backward to lab values: warehouse rows present, hash matches. Lab pipeline S3 data present and unmodified.
- Forward to clinical action: override is sealed. Attending's reasoning is in the override entry. **Subsequent EHR documentation has been substantially rewritten.** The progress note shows three edits between April 13 and April 22. The current version of the note describes a clinical reasoning that is plausibly different from what the override entry's structured pick-list reason indicated at the time. Epic shows the version history. The chain shows the override. Reconciling the two requires reading both, side by side, and deciding which one represents what the clinician actually believed at 11:08 on April 13.

Mike, who had run the trace for all five, sat back from his laptop.

Chen, doing the warehouse cross-check: "On Alert 3 — the missing potassium. The chain entry has the *value* the model received. We could reconstruct the underlying source by trusting the chain. If the question is 'what did the model see,' the answer is 5.8. If the question is 'was 5.8 the actual lab value,' we can no longer prove it independently. That's the §1.2 scope line played out in a real case — the chain answers the first question; the second is outside the chain's evidentiary reach."

Diana: "And on Alert 5 — the rewritten note. The override entry is sealed. The structured pick-list reason is sealed. The free-text reasoning is in a note that has been edited three times. If we are defending the model's behavior, the chain is sufficient. If we are defending the clinician's judgment, we are reading three versions of a note and arguing about which one is the 'real' one."

Patricia: "That's the reality of clinical documentation. Note editing is medically appropriate. It is also an evidentiary problem and we know it."

Dawn tallied on the whiteboard.

| Alert | AI inference | Backward (lab inputs) | Forward (clinical action) |
|---|---|---|---|
| 1 | PASS | PASS (no chain, data present) | PASS |
| 2 | PASS | PASS (no chain, data present) | PARTIAL — note edited twice |
| 3 | PASS | **FAIL — input source deleted** | PASS |
| 4 | PASS | PASS (no chain, data present) | PASS |
| 5 | PASS | PASS (no chain, data present) | **FAIL — note rewritten three times** |

She put down the marker.

"AI side reconciles five of five. Backward, four of five — one had its source record deleted by a retention job sixty days ago. Forward, three of five — two have edited Epic notes that we cannot prove represent the original clinical reasoning."

The room was quiet.

Patricia broke the silence. "That is the budget request."

Dawn: "That is exactly the budget request."

> **⚠️ Surprise #8 — The chain has a rooting failure at the EHR boundary**
> The sepsis chain proves what the model saw and what it said. Backward from the chain — to the lab values that fed the model — the warehouse remains reconcilable for the 90-day window the inference pipeline has been live, but the upstream lab pipeline is mutable and one of five tested alerts had its source lab record deleted by a retention job. Forward from the chain — to the clinical action that followed — Epic notes can be edited indefinitely. Two of five tested alerts have post-hoc note edits. The chain is sealed in the middle. The two ends are not. The §10.19 chain-coverage map needs to be re-read against this finding: the boundary is not just "what is chain-instrumented vs not" but "what does the institution rely on the chain for, and at what end is the rooting weakest." Per §10.13 evidentiary-artifacts retention, the documentation that substantiates FRE 901(b)(9) authentication of the process is in the institution's hands; for the sepsis chain that documentation is the SDK manifest, source-code hash, HSM configuration, daily seal-job logs, change-management records, and verifier output for each tenant-day. For the EHR notes, the comparable evidentiary-artifacts package does not exist — the note version history is in Epic's hands and Epic's retention controls; an Epic-side retention or correction policy invalidates the institution's reconstruction at the edges of the chart's life. The remediation is to extend the chain to cover the EHR write-back path under the §10.21-style cross-vendor model-handover discipline (the chain entries reference the Epic-side write events with hash anchors); the operational details land under §10.16's connector-lag pattern and §4.4.6 connector-source attribution.

Mike, quietly: "On the inference path you can prove what the model said. On every other path the integrity is still hope."

Dawn looked at him sharply. Then she wrote that sentence down in her notebook, word for word.

---

### 😬 3:45 PM — The Friction Builds

The clinicians arrived at 3:45 — three of them, all ICU attendings, brought in by Dr. Wei to talk through the override pattern.

Dr. Friedman went first. He was direct in a way only ICU attendings get to be.

"The model fires too often. We override about thirty percent of the time. Mostly because the model doesn't know context. It sees a patient with elevated lactate and a fever and calls sepsis. It doesn't know the patient is two hours post-op and the lactate is from the surgical stress, not from sepsis."

Dawn: "When you override, what do you record?"

"In the override UI we pick a reason from a list. There are about twelve options. Most of the time my reason is 'clinical context not captured by model.' Then I go write a real progress note in Epic explaining what's actually going on."

Diana: "The override entry is sealed. Your structured reason is sealed. Both are part of the canonical bytes the per-event MAC covers under §5; the override entry binds to the parent prediction via §4.4 `parent_run_id` / `parent_seq`."

"Yes. The note in Epic is not."

Dawn: "And the note in Epic is where your actual clinical reasoning lives."

"Yes."

"And you can edit that note later."

"Yes. We do, sometimes. Patient's clinical course evolves. We update the progress note to reflect new information. Standard medical practice. We are required to do it."

"Required to update the note. Not required to keep the original."

Dr. Friedman paused. "Epic has version history. We can see what the note used to say."

"For how long?"

"Forever, I think. As long as the patient's chart exists."

Patricia: "Epic retains note versions for the life of the chart, yes. But the version history is in a separate Epic table. We have a 90-day retention on certain audit views. The note versions themselves are retained longer but the audit view of who-edited-what-when is on a 90-day rolling window."

Diana: "So an edit from a hundred days ago shows up as a different version of the note, but the audit trail of *who* made the edit and *when* is gone."

"For administrative-action edits, yes. For clinical-content edits, the metadata is retained as part of the note version itself."

"And clinical-content edits include the kind Dr. Friedman is describing."

"Most of them. Yes."

Dr. Friedman: "Look. I'm an attending. I write notes. I update them as the patient evolves. That's medical practice. If you're going to chain my notes, that's a different conversation, and frankly I have concerns about it."

Dawn: "I am not going to chain your notes today. I am noting that the chain you have today does not extend to your notes, and that the clinical reasoning behind your overrides — the override decisions that the chain *does* capture — is in your notes, which are mutable. That is a finding. It is not a request."

He thought about that. "Okay. That's fair."

> **⚠️ Surprise #9 — Override decisions are sealed; the reasoning behind them is in mutable Epic notes**
> The chain captures every clinician override of the sepsis model, with the structured pick-list reason and the timestamp, sealed. The substantive clinical reasoning — written in Epic progress notes — is editable indefinitely. Epic retains note versions for the life of the chart, but the audit view of administrative-action edits is on a 90-day window. For an audit defending an override decision more than 90 days after the fact, the reasoning behind the override may be reconstructible from Epic version history but the integrity of that reconstruction depends on Epic's behavior, not on the chain. The §10.13 evidentiary-artifacts retention discipline names the documentary evidence the chain itself depends on; an Epic-side note version is outside that package and outside the chain's §10.5 HSM-rooted custody. Where the override decision later becomes load-bearing in litigation or in a state-medical-board review, the chain answers what was decided and when; Epic answers (or attempts to answer) why. The two evidentiary regimes compose without overlap — chain for decision-time integrity, Epic for clinical narrative — but the institution should know the line between them and not confuse one for the other.

Dr. Wei added, after the attendings had left: "We've talked about extending the chain to clinical notes. The clinicians don't want it. Not because they're hiding anything — because they are required to update notes as a matter of medical practice and they don't want every update to be a sealed event that a plaintiff lawyer can wave around in court."

Dawn: "That is a real concern. Worth discussing separately. It is also worth knowing what the gap is, even if we choose not to close it. The §10.22 redaction-discipline framework would let you redact pre-MAC at the SDK boundary so plaintiff's counsel cannot demand unredacted note content from the chain — but redaction does not solve the version-tracking problem. That's a deeper design question."

"Agreed."

She wrote that down too.

---

### 🔍 4:30 PM — The Boundary Question

Patricia came back into the small huddle room at 4:30. The CAE had dropped off the Zoom. It was just Dawn, Tom, and Patricia.

Patricia leaned forward. "I need to ask you something directly. Will FDA accept the AI-side chain even if the EHR-side isn't chained?"

Dawn had been waiting for this question all day. She'd thought about her answer on the drive in. She thought about it again now.

"FDA SaMD post-market surveillance: yes. The AI side is what they review. They want to know that the model's behavior in the field matches the model's behavior in the clearance submission. They want predictive performance monitoring. They want adverse-event tracking. They want a credible record of what the model was asked and what it answered. You have that. The chain you have today satisfies the FDA SaMD post-market evidentiary burden as I understand it. The §1.1 Daubert grounding gives the FDA reviewer the four-factor answers from the spec text alone, the §10.26 reference-verifier distribution lets the reviewer run the verifier independently, the §10.18 CC8.1-and-runbook cross-referencing rule means your runbook points back at the spec section that governs each control, and the §10.13 evidentiary-artifacts retention package is what you'd produce under cross-examination. The NTP discipline behind the timestamps is per §10.4 (and §10.14 names the optional RFC 3161 trusted-time integration that you have not adopted yet — that's an open candidate for FDA-litigation-anticipated workloads). I am not your FDA counsel and you should confirm with them. But based on what I have seen today, yes."

Patricia exhaled. "Okay."

"HIPAA Security Rule audit-control: partial. The Security Rule requires audit controls for systems that handle ePHI. The sepsis service has audit controls that are stronger than the rule requires — chain integrity per §4.1, append-only storage per §10.3, HSM custody per §10.5, all binding the access logs and the decision logs together. Epic has audit controls that are at the field median for hospital systems. The lab pipeline and the claims engine and the CRM are at or below the field median. If an OCR auditor walked in tomorrow, the sepsis service would impress them. The rest of the environment would not. You would not be cited, in my opinion. You would not be commended either. You would be average."

"And the plaintiff bar?"

"That depends on the case theory."

Tom looked up.

Dawn continued. "If a plaintiff sues you for a clinical error and the case theory is *the model's output was wrong*, the chain saves you. You can produce a sealed record of exactly what the model was asked, exactly what it said, exactly which clinician saw it, exactly what they did. Plaintiff's expert cannot rewrite that record. The chain is your defense. The §1.2 epistemic-scope claim — what the AI said at time T, and that the record was not tampered after capture — is exactly what defends you in that case theory.

If the case theory is *the clinical judgment that followed was wrong*, the chain doesn't help. The chain shows the clinician saw the alert. It doesn't show what the clinician was thinking. The clinical reasoning is in Epic notes. Plaintiff's expert will pull the note version history and argue about edits. Your defense in that case is the standard medical-malpractice defense — it doesn't get worse because of the chain, but it doesn't get better. §1.2 names this explicitly: the chain does not prove statement-level accuracy or compliance.

If the case theory is *the lab values that fed the model were wrong*, the chain partially helps. You can prove what the model received. You cannot prove what the lab actually measured. The lab pipeline is upstream of the chain. Plaintiff will subpoena the lab analyzer logs and the integration engine logs and the S3 raw data, and the integrity of those records depends on whether the engineers running those systems happened to have left logging on for the relevant period. That's the §1.2 epistemic boundary again, and it's also a §10.19 chain-coverage-map gap — the lab pipeline is on the wrong side of the line."

Patricia listened. She had her hands folded on the table.

"That's the assessment," she said.

"That's the assessment."

"Okay. So my budget pitch is — extend the chain to billing in Q3, EHR in Q4, lab pipeline in Q1 next year. Roughly in priority order of which gap is most likely to bite us in court."

"That ordering is wrong."

Patricia paused. "Tell me."

"Lab pipeline first. Then EHR. Then billing. The lab pipeline is the *input* side of the chain. Right now your sealed inference rests on a mutable input. That's the most fragile boundary. EHR is the *output* side. Less fragile because clinical judgment is ultimately a human decision, not a record-integrity question. Billing has its own evidentiary patterns and is the least likely to come up in a clinical-error suit — it'll come up in a CMS audit, but you have time. CMS-0057-F prior-authorization timeliness is your billing-side regulatory exposure and the deadline is January 1 2026 for prior-auth API conformance, but that's a different audit than a malpractice claim. When you do extend the chain to billing, the prior-auth adverse-determination notices follow the §10.11 translation discipline by analogy — state-insurance-law adverse-action notices already get cited by §10.11 as a conformant use-case for that schema, and the §10.11.1 reasons schema gives you the integrity-bound list of decision reasons. The §10.11.2 FCRA reinvestigation timing schema is the right shape for prior-authorization appeal timing — by analogy, not directly, since FCRA is a credit-reporting regime, but the integrity-bound clock-start, clock-extension, completion, and consumer-notification timestamps map cleanly onto CMS-0057-F's 7-day-standard / 72-hour-expedited windows."

Patricia thought for a long beat. "Lab. Then EHR. Then billing. Okay."

"And in parallel: clinician AD. The chain you have today rests on a legacy AD identity layer. If you extend the chain without firming up the identity root, you're building taller walls on the same foundation. The IAM-as-chain pattern your sepsis service runs on the service-account side is the model — every AD grant a chain entry, every revocation a chain entry, time-boxed expiration enforced by a chain-driven worker. That's not specific to any spec section; it's the §10.3 append-only enforcement plus the §3 `chain_kind = "operational"` discipline applied to identity events. And the §10.1 key-fingerprint reconciliation discipline — every `key_fingerprint` rotation auditable, no silent IKM swaps — is the kind of identity-root hygiene the AD modernization should adopt."

Dawn had one more.

"And §10.25. When the chain extends to those new systems, the SDK on each system has to honor the run-resume contract. If the SDK loses local state, it queries the ledger's chain-tail endpoint before emitting the next entry. Three-place tail acquisition — in-memory, sidecar, ledger — and the SDK refuses to emit if it cannot find a tail and the ledger is unreachable. That closes the silent-restart attack class. The §10.17 HSM partition ceremony attestation is also worth knowing about — every CloudHSM partition ceremony you run for the new tenants emits a chain-coupled attestation event. And §10.24 entity succession governs whether the tenant_id stays stable across any future organizational change. We won't run into those today, but they're going to come up the moment Mercator extends the chain to the insurance arm — that subsidiary may need its own tenant_id under §10.15 multi-region resilience Pattern B if the regulatory regime treats it as a separate entity."

"Diana already told me. I have an AD modernization project in flight. It was scheduled to land in Q4. I'll move it forward."

"Good."

Tom: "I'll write the cross-walk between the lab-first ordering and HITRUST CSF v11 control families. The CAE will want it for the board materials."

Patricia: "Thank you."

She stood up. "Five-thirty debrief?"

"Five-thirty."

---

### 🌆 5:30 PM — Auditor Debrief

The full team reconvened. Patricia, the CAE on Zoom, Patricia's two direct reports, Dawn's eight-person team. The conference room was the same one from the morning.

Dawn ran the debrief.

"This is a bifurcated assessment. We are going to give you two findings sets. One for the AI side. One for the legacy side. Both are real. Neither cancels the other. The bifurcation tracks the §10.19 chain-coverage map you printed for us at kickoff."

She pulled up a single slide.

#### ✅ What They Found on the AI Side

| # | Finding | Spec anchor |
|---|---|---|
| 1 | Inference entries are content-complete and sealed. Model ID, version, prompt, response, tool calls, clinician override — all in the chain. | §4.1, §4.4 |
| 2 | Verifier resolves cleanly across signing-key rotations and across the 90-day operating window. | §4.2, §10.10, §10.12 |
| 3 | Model lifecycle events (deploy, rollback, retrain) are sealed entries with the same structure as inference entries. | §4.4.2, §3 |
| 4 | Tool calls within model inference are first-class chain entries with parameter and result-hash sealing. | §3, §4.4 |
| 5 | Inference-input warehouse is reconcilable to the chain. The `patient_context_hash` recomputes from append-only warehouse rows. | §5, §10.3 |
| 6 | Sepsis service IAM is chain-coupled end-to-end. Every credential lifecycle event and every policy change is sealed. | §3 (operational), §10.3 |
| 7 | API gateway and fallback paths are sealed. Even the April 18 CloudHSM maintenance window shows up as 247 sealed fallback entries. | §4.3.1, §4.4.3, §4.4.4 |
| 8 | Sepsis ETL has cross-checkpoint reconciliation as sealed events. Halt-the-line responses to checksum divergence are themselves sealed. | §3 (operational), §4.1, §4.2 |

**AI side: 0 Gaps. 0 Partials. 0 non-conformances per §10.16. Audit-passes.**

#### ❌ What They Found on the Legacy Side

| # | Finding | Spec anchor (would-remediate) |
|---|---|---|
| 1 | Clinician identity (legacy AD) has no chain and 23 active temporary-admin grants, 6 older than a year. | §3 (operational), §10.3 |
| 2 | Service-account chain rests on legacy-AD root of trust for clinician overrides. SAML subject is only as trustworthy as the upstream IdP. | §1.2, §1.4 |
| 3 | Epic-to-billing handoff has no chain. Mulesoft retention dial is service-account-controlled by 8 engineers with no audit. | §10.13, §10.3 |
| 4 | Claims-processing audit table is mutable by 5 SREs. Standard CloudWatch retention with no out-of-band signing. | §10.3, §10.5, §10.23 |
| 5 | Lab pipeline has no cross-destination reconciliation between Epic and warehouse. 3-day DLQ. Storage layer fully mutable. | §4.4.6, §10.16 (pattern), §10.19 |
| 6 | Lab pipeline storage layer (S3) is CloudTrail-disable-able by the 3 engineers running the pipeline. | §10.3, §1.4 |
| 7 | Claims ETL checksum log is editable. 3 rows show edits in last 60 days with no audit trigger. | §10.23, §10.3 |
| 8 | The chain has a rooting failure at the EHR boundary. Of 5 tested sepsis alerts: 5/5 inference reconciled, 4/5 backward (1 source record deleted by retention job), 3/5 forward (2 Epic notes rewritten post-hoc). | §10.13, §10.19, §10.21 (pattern) |
| 9 | Override decisions are sealed; reasoning behind them is in mutable Epic notes with 90-day audit-view retention. | §10.13, §10.22 (pattern) |
| 10 | Salesforce CRM (member services + patient-outreach navigators) — Description field history disabled for storage-cost reasons. PHI sits in unaudited free-text. HITRUST Partial renewed three cycles, remediation unfunded. | §4.4.6, §10.16, §10.19, §10.22 |

**Legacy side: 5 Gaps. 7 Partials.** (The Partials are the items that have *some* protection — Epic version history, Mulesoft 6-week logs, CloudTrail when it's enabled — but where the protection is at the discretion of the system operator and is therefore not chain-grade evidence.)

---

#### 🔁 Side-by-side comparison

| Dimension | AI side (sepsis-cds) | Legacy side (Epic, lab, billing, claims, CRM) |
|---|---|---|
| Record integrity | Sealed via HMAC chain (§4.1) + daily Ed25519 seal (§4.2, §4.3) | In-place edits possible; backups on tape |
| Identity coupling | Service IAM chain-coupled (§3 operational, §10.3) | Legacy AD with no chain; 23 temp-admin grants active |
| Reconciliation | Cross-checkpoint sealed events (§4.1, §4.2) | Manual; checksum tables editable |
| Retention | Indefinite per §10.13, sealed by §10.3 storage | 90-day windows controlled by system operators |
| Verifier | `herald-verify` resolves in 12 steps (§7, §10.12, §10.26) | No verifier; audit by inspection |
| FDA SaMD post-market | Defensible (§1.1 Daubert, §10.13 evidentiary artifacts) | Not in scope |
| HIPAA audit-control | Above field median | At or below field median |
| Plaintiff defense — model output | Strong (§1.2 epistemic scope a) and b)) | n/a |
| Plaintiff defense — clinical judgment | Indirect (chain shows the alert was seen) — §1.2 (c) explicitly out of scope | Standard malpractice posture |
| Plaintiff defense — input data | Sealed at model boundary (§5 canonical-bytes); mutable upstream | Mutable upstream |
| Custody-layer count behind a verifying false negative | Three (§1.4: IKM, ledger storage, HSM) plus §1.2 SDK-process residual | One (the engineers who run each system) |

Dawn paused on the slide for a beat longer than the rest.

"Half the river is sealed. The half upstream is not. You can prove what your model said. You cannot prove what it was given."

The CAE, on Zoom: "That's the line for the board memo, right there."

Patricia: "That's the line."

She turned to the room. "Dawn. Tom. Thank you. This is the assessment I asked for. The sequencing recommendation — lab first, EHR second, billing third — is going into the budget request next week. The AD modernization is moving from Q4 to Q3. The board meets June 14. I'd like permission to share the bifurcation framing in the board materials."

Dawn: "Permission granted. We'll send the formal report by Friday. The bifurcation framing is the framing. Use it."

The CAE on Zoom: "One more thing. I want it on record that this assessment found the AI-side controls to be at the top quartile for healthcare systems we benchmark against, and the legacy-side controls to be at or below the median. That is the comparative posture we are going to take to the board. Dawn — you good with that characterization?"

Dawn took a beat. "I am, with one note. The 'top quartile' phrasing comes with a caveat — the population of healthcare systems running chain-grade AI controls in production is small. Mercator is in the top quartile of a fairly small group. I'd rather the report say 'meets or exceeds best-known practices observed in deployed healthcare AI systems audited in the last twelve months.' That phrasing is more defensible and it is more accurate."

The CAE nodded. "Better. Use that."

Patricia: "Agreed."

Tom was writing the language down verbatim on his iPad.

Elena raised a hand. "One more for the record. The CRM Partial findings — the field-history-disabled situation in Salesforce — has been on the HITRUST attestation for three cycles. The remediation has not been funded. I want the report to flag that the legacy-side gaps include items where the *gap is known* and the *remediation has been deferred for budget reasons*, not just items where the gap is newly discovered. That distinction matters for the board read. And one more thing — I want to flag a forward-looking risk per §10.16. If Mercator extends the chain to the CRM via a SaaS-edge mirror connector under §4.4.6, the CC8.1 control description for that connector MUST quantify the four lag bounds — median, 95th-percentile SLO, alerting threshold, RTO — before it ships. Imprecise lag wording in the runbook (`'near real-time'`, `'low-latency mirror'`, anything similar) is a non-conformance per §10.16's severity-classification clause and the audit team will not have discretion to downgrade it. We saw that exact pattern at last week's banking engagement; it cost them their otherwise-spotless audit."

Patricia: "Noted. The CC8.1 will name the four numbers before any connector ships."

Elena: "Good."

Patricia stood up, shook Dawn's hand, then Tom's, then went around the room and shook each team member's hand individually. It took ninety seconds. She thanked each of them by name.

"Drive safely."

The team packed up. There was the usual quiet shuffle of laptops closing and notebooks going into bags. Diana paused on the way out and asked Patricia a private question about the AD modernization timeline. Patricia answered it. Mike traded business cards with Dr. Wei, who had come up for the debrief tail. Chen and Devansh exchanged GitHub handles. Luis caught Marcus from the lab pipeline team in the hallway and gave him an unsolicited recommendation about a CloudTrail-immutability tool that would not have prevented the disable risk but would have made the disable visible faster — closer to the spirit of §10.3 append-only enforcement applied at the lab S3 bucket. Marcus thanked him.

Dawn watched the room empty. She gathered her notes, clipped her pen back to the cover, and slung her bag over her shoulder.

Tom held the door for her on the way out. "Round two is in the rear-view."

"Round two is the new comparison point," Dawn said. "Last week was the floor. Two weeks ago was the ceiling. Today is what real-world transition looks like."

"Good story."

"Good story."

---

### 🧾 Final Assessment Theme

The drive home was forty-five minutes. Dawn had her coffee, refilled, in the cup holder. The sun was setting over the hospital district behind her.

She thought about the day. About the bifurcation. About Dr. Wei's terminal showing a clean verifier resolution at 9:30 in the morning, exit code 0 against §7's twelve-step procedure. About Luis pointing at the S3 bucket policy at lunch and saying "This is the diary." About Patricia's question at 4:30 and the answer she had given. About Dr. Wei's offhand comment that they had hired the engineer who built Northbridge's chain — which explained why two sites separated by an industry boundary looked so similar at the chain-construction layer and so different everywhere else.

She thought about Northbridge, two weeks ago, where everything had been sealed and the audit had been almost boring — except for the §10.16 mirror-runbook non-conformance Elena had caught at lunch.

She thought about the financial services firm last week, where nothing had been sealed and the audit had been an indictment.

Mercator was neither. Mercator was a system mid-transition. Honest about the line — the §10.19 chain-coverage map drawn on one printed page on every chair at kickoff. Seeking the assessment to fund the work to extend the line. The AI side was as good as any she had seen. The legacy side was as bad as any she had seen. And the boundary between them was the actually-interesting story.

She picked up her phone at a red light and dictated a one-line note for the report's executive summary.

> *"On the inference path, Mercator can demonstrate what the model said and that the record was not tampered with — exactly what spec §1.2 names as the chain's evidentiary scope. On every other path, integrity is still hope."*

The light turned green. She put the phone down and drove home.

---

*— end of audit diary —*
