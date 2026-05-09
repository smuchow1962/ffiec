# 🧾 Diary of an Audit Day — Pacific Crescent Power & Gas

**Engagement:** NERC CIP audit readiness assessment with PHMSA pipeline integrity-management overlay
**Client:** Pacific Crescent Power & Gas — investor-owned utility, Pacific Northwest, ~3.1M customer accounts (electric + gas), multi-state (WA, OR, parts of northern CA and ID), HQ in Portland
**Posture:** Partial TesseraSeal deployment — chained on a single AI use case (gas-pipeline leak detection); not chained on SCADA, the PI historian, the AMI head-end, the OMS, the CIS, or the customer CRM
**Date:** Tuesday, the week after Helmstad Bio
**Auditor:** the same eight-person team — diary baseline, Mercator, Northbridge, Stelvio, Atrio, Helmstad

---

## Context

Pacific Crescent runs hydro from the Columbia, a fleet of natural-gas combined-cycle plants down the I-5 corridor, a recently added wind portfolio in eastern Washington, and a gas-distribution system that snakes through Portland, Vancouver, Salem, and parts of Eugene. Three regulators in three states. NERC CIP-002 through CIP-014 on the bulk-electric-system side. PHMSA pipeline integrity on the gas side. WUTC, OPUC, and a small slice of CA-CPUC on the customer side. DOE/CESER cybersecurity expectations on top of all of it. IEC 62443 on the industrial-control side because two of their large industrial customers have started asking for it in their procurement language. FERC Order 887 cybersecurity-audit expectations layered on top of NERC. WECC regional-entity auditors do the actual NERC walkthroughs.

Fourteen months ago, a near-miss gas leak in a Portland neighborhood turned political. The utility had detected residual methane on a routine patrol, but the detection happened roughly three hours after a homeowner had already smelled it and called the utility's customer line. The PUC opened an inquiry. The mayor's office called twice. The utility's CEO accelerated a procurement that had been quietly in progress: an AI-driven gas-pipeline leak-detection system that fuses methane sensor data, weather, soil-saturation models, and historic incident data to predict pipeline-failure probability in real time.

Nine months ago, that system went live. Pacific Crescent stood up TesseraSeal at the same time, chained on the leak-detection service. Three `service.name` values inside one tenant — `pipeline-leak-detection` (in production), `pipeline-integrity-trending` (in production, lower-priority analytics), and `outage-prediction-pilot` (pre-production, not yet in scope today). The `tenant_id` is `pacific-crescent`, conformant under spec §3 character class (alphanumerics plus hyphen), so HKDF tenant-binding under §4.1's `info_for_tenant` parameter is unambiguous.

Everything else at Pacific Crescent runs the way utilities have always run.

The OT side — GE iFIX SCADA on Windows Server 2016, OSIsoft PI historian, Itron OpenWay AMI head-end on version 5.2 — is not internet-connected. NERC CIP-005 requires that. But "air-gapped from the internet" is not the same as "integrity-controlled," and that's a distinction the regulators are starting to write into their guidance. The customer-side stack — a customized Oracle CIS, an outage-management system, Salesforce CRM — runs on the IT business network. Different mutability shape. Same mutability problem.

Esme Yamashita, Pacific Crescent's Chief Compliance Officer, came over from BPA two years ago. She spent two years on a NERC standards drafting team for CIP-013, which makes her one of the few people in the room who has read the FFIEC chain-of-custody v1.0b spec end-to-end alongside the NERC supply-chain standards. She knows what a NERC auditor will look for, and she knows what they will quietly skip. She doesn't oversell. On the prep call she told Dawn: "We need a map. NERC CIP-007 is the area where the regulators are going to push hardest in the next two years. They're starting to ask AI questions. Tell me what we already have and where the gaps will land. And give me the chain-coverage map per spec §10.19 so the WECC auditor and the PHMSA inspector see the same picture."

By the time the team arrived at Pacific Crescent, they had seen TesseraSeal hold up under banking, healthcare, BaaS multi-tenancy, industrial three-tier, and biopharma. Northbridge had been the cleanest engagement of the cycle — one §10.16 non-conformance, the chain otherwise held byte-for-byte across forty-three citations. Five engagements later Dawn had stopped expecting it to repeat. Utility was a new regulatory composition — NERC CIP, FERC Order 887 — but the chain primitive was familiar. The new question wasn't whether the chain held; it was whether NERC CIP-002 through CIP-014 composed with the spec. Particularly the BES Cyber System scope boundary under CIP-002 categorization and the OT/IT integration shape — SCADA, EMS, and OMS mirrors — under CIP-005 and CIP-007. That was the composition the team was here to test.

This is the diary of that day.

---

## Audit Team

- **Dawn** — Lead Auditor (governance and narrative)
- **Raj** — Database specialist
- **Elena** — CRM systems
- **Mike** — Application and API layer
- **Diana** — IAM and access control
- **Luis** — DevOps, logs, pipelines
- **Chen** — Data engineering and ETL
- **Tom** — Internal-audit liaison specialist (visiting team — partners with the client CCO)

Client-side liaison: **Esme Yamashita**, Chief Compliance Officer, Pacific Crescent Power & Gas. Ex-BPA. NERC CIP-013 drafting veteran. Direct. Technical. Doesn't oversell.

---

## 🌅 8:30 AM — Kickoff and the Drive In

Dawn rode in with Tom from the airport hotel. Twenty minutes north on I-5, then off at the river. The Pacific Crescent control center sat on a low rise overlooking the Willamette, glass and concrete, fenced with the kind of fence that meant something. NERC-CIP-classified facility. They saw the badge readers from the parking lot.

Tom drove. Dawn drank coffee.

"Five engagements in five weeks." Tom merged across two lanes. "What's the through-line at this point?"

Dawn watched the river. "Northbridge was the gold standard. Banking. TesseraSeal in everything. We ran out of things to find by 3 PM. Forty-three citations across the report, one §10.16 non-conformance on a Salesforce mirror, and a §10.17 partition-ceremony posture that was already running clean."

"Mercator."

"Half the river sealed, half not. Healthcare. Imaging side chained, claims side mutable. We wrote the seam down the middle of the report and the CMO understood it instantly. §10.19 chain-coverage map made the seam testable rather than rhetorical."

"Stelvio."

"Manufacturing. Three zones. AI sealed. OT mill-floor mutable. IT business mutable. Maria knew where every gap was and wanted the report so she could take it to the CFO Friday. We wrote it that way. Stelvio's Proficy historian is the same boundary problem we'll see today on the PI side — chain begins at the AI service per §1.2 epistemic scope, upstream is trust-by-policy."

"Atrio."

"Banking-as-a-service. Forty-seven tenants, all clean. The cleanest multi-tenant verifier work we've ever seen. We were on the plane home by 4. Per-tenant key fingerprint reconciliation under §10.1 was running daily across all forty-seven."

"Helmstad."

"Biopharma. ALCOA+ on the AI side, CRO data legacy. The seam was in a different place again. The CRO data showed up in the chain by reference but the source-of-truth lived in someone else's system entirely. §10.19 third-party-without-contractual-inspection category. We hash-anchored what we could under `audit.external_artifact.*` and named the rest."

Tom checked the GPS. Five minutes out.

"And today?"

Dawn put her cup down. "Pacific Crescent has 3.1 million customers and a fleet of pipelines. The legacy gaps here are not just compliance findings."

"What are they."

"Stelvio's worst case was a yield-loss claim on a heat of steel. Mercator's worst case was a misadjudicated claim. Atrio's worst case was a sponsor-bank reconciliation finding. Helmstad's worst case was a CRO data-trace gap on a Phase 2 trial."

Tom waited.

"Pacific Crescent's worst case is a wrong reading and a dismissed alarm and someone's house blows up. Stelvio with public-safety consequences. The §1.2 epistemic-scope distinction — chain proves what the AI said at time T, chain does not prove the upstream sensor was honest — is the whole report in one sentence."

Tom didn't say anything for a while.

"It never is."

"It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach. We're going to find both today."

They pulled into the visitor lot at 8:24.

Esme met them at the badge desk. Mid-forties, dark blazer, the kind of badge clipped to her lapel that had three different access levels printed on it. She shook Dawn's hand once, firmly, and Tom's the same way.

"You'll go through three checkpoints to get to the operations floor. The control center itself is CIP-categorized. Phones in the locker. Laptops registered at the second checkpoint. We have a clean room beside the control floor where you can work — same network segment as corporate but with read-only feeds from the OT side. The HSM and the AI ledger are reachable from the clean room. SCADA and the historian are not — you'll see those over a screen-share from one of my engineers."

"Understood."

"I want to be straight with you about scope before we start." Esme walked them toward the first checkpoint as she spoke. "The leak-detection AI is in scope under NERC CIP-007 because it touches a BES Cyber Asset adjacent to gas-electric tie-points, and it's in scope under PHMSA because it's a pipeline-safety system. The PI historian is in scope under CIP. The AMI head-end is in scope. The OMS is partly in scope — work orders that touch BES assets are in. The CIS and the Salesforce CRM are out of scope under NERC but very much in scope for the PUC. We're going to look at all of it because that's how I think about the program. The report can split scope however you need it to."

Dawn nodded. "That's how we'd write it anyway. And that lines up with spec §10.19 — the chain-coverage map enumerates chain-instrumented institutional systems, not-yet-instrumented institutional systems, third-party-with-inspection, and third-party-without-inspection. Four categories. Your scope split maps cleanly."

"That's the document I want at the end of today."

"You'll have it."

The team kitted up. Phones in lockers. Laptops registered. Three badge taps later, they were on the operations floor.

> **🔍 Dawn's note (internal):**
> *It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach.*
>
> *Calibrate. The AI side will pass. The OT side won't. The customer-billing side won't. But the consequence axis is different from Stelvio's. Stelvio's worst case was money. Pacific Crescent's worst case is people. Document accordingly. The §1.2 epistemic-scope language is the controlled vocabulary the report needs.*

---

## 🧩 9:15 AM — The Control-Room Walkthrough

The operations floor was quiet in the way only a working control room is quiet. Eight workstations across two rows. Three large screens on the wall — one showing the natural-gas distribution map, one showing the electric-transmission system, one showing the AI leak-detection dashboard. Methane sensor readings updated every fifteen seconds. Soil-saturation tiles colored a dozen pipeline segments in shades from green through amber.

Esme walked the team to a station along the back wall. A dispatcher sat at it — name tag said *Marcus*. Headset on his neck. Coffee at his elbow.

"Marcus is on the morning gas-distribution shift. The AI dashboard you see is `pipeline-leak-detection`. Every prediction the model produces is on this screen for thirty seconds before it auto-archives. Any prediction above 0.4 confidence sits on the queue until someone actions it. Marcus, talk them through what you see right now."

Marcus turned in his chair. "We've got a low-confidence flag on segment SE-Powell-44. Methane uptick of 1.2 ppm above background, but the wind shifted right when it triggered, so the model knocked the confidence down. I'm watching it."

Mike looked at the screen. "What's the confidence on it?"

"0.31. Below dispatch threshold."

"What happens next?"

"If it climbs above 0.4 in the next ten minutes, the queue holds it for me to action. If it drops back to baseline for three readings in a row, the model writes it off and the prediction archives."

Dawn turned to Esme. "Action means dispatch?"

"Action means decision. Marcus gets four buttons — dispatch crew, monitor, dismiss, escalate-to-supervisor. Whatever he picks, that decision lands in the chain along with the model's prediction."

Mike opened his laptop, balanced it on the edge of Marcus's station, and pulled up TesseraSeal. He filtered to `service.name = pipeline-leak-detection` and the last five minutes.

The Powell-44 entry was sitting at the top of the list. Confidence 0.31. No dispatcher action yet. The chain entry referenced a methane sensor ID, a weather snapshot hash, a soil-saturation model output hash, and a prediction. Mike tilted the laptop so Dawn could see the per-entry stamp — `chain_kind = "audit"` per §3 enumeration, `format_version = "v1"`, `key_version = 4`, `key_fingerprint` 16 raw bytes per §4.1, `kms_handle_uri = "thales-luna:partition=pacific-crescent-leak-prod"`, `gen_ai.request.model` and `gen_ai.response.model` both populated per §4.4 normative requirement and §7 step 12a.

Marcus glanced at the screen. The reading dropped to 0.9 ppm above background. Then 0.6. Then back to baseline.

"There. It's washing out."

The dashboard moved Powell-44 to the archive. Marcus tapped *dismiss* on the queue, picked a reason code from a dropdown — `wind shift, baseline restored` — and submitted.

A new chain entry hit Mike's screen.

"There it is." Mike tilted the laptop. "Decision entry. References the prediction entry via `parent_run_id` and `parent_seq` per §4.4. Dispatcher ID, reason code, timestamp, model_id, model_version. The dispatcher decision is the child; the model prediction is the parent. Same parent-linkage discipline §10.11 uses for translation entries."

Dawn leaned in. "Run the verifier on it."

Mike copied the entry ID into his terminal:

```
herald-verify --tenant=pacific-crescent --service=pipeline-leak-detection \
              --date=2026-05-05 --entry-id=2026-05-05-PowellSE-44-dismiss
```

Four seconds. The terminal returned:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key leak-prod-2026-q1
```

Mike turned the laptop. Marcus looked at it. Esme looked at it.

"Twelve steps, four seconds, on the operations clean-room network." Mike snapped the laptop shut. "That's the thing working. Per the §7 verification procedure, exit code 0 per §10.12 means PASS — the JCS self-test fired at startup per the Round-17 NIST-P2 amendment, the format-version check passed at step 1, the HKDF-inputs digest matched at step 2, the structural walk passed at step 6, the IKM lookup resolved at step 7, the fingerprint matched at step 8, the MAC recompute matched at step 9, the Merkle root matched at step 10, the §4.3 v1.0b 12-line `sign_payload` reconstructed and the Ed25519 signature verified at step 11. The dismissal is sealed. The model's inputs are sealed. The dispatcher's reason code is sealed. If anyone ever asks why that low-confidence flag was dismissed, you have an answer that's the same answer six months from now and six years from now."

Esme nodded once. "That's what it was procured to do."

> **✓ Confirmation #1 — chain integrity at the AI boundary**
> The leak-detection chain is live on the operations floor and producing verifiable entries within ~200 ms of each prediction or dispatcher decision. Mike re-verified a real-time dismissal in four seconds standing at a working dispatcher station, the §7 12-step procedure returning PASS with §10.12 exit code 0. Reason codes are captured. Inputs are hashed under the per-event MAC per §4.1. The infrastructure is real and observable on the production line. Per spec §1.4 compositional security, the three independent authentication layers — per-event HMAC, daily Merkle seal, HSM-rooted Ed25519 signature — each add an independent integrity property; the demonstrated PASS exercises all three layers.

Dawn wrote in her notebook: *AI dashboard live. Dispatcher decisions sealed with reason codes. Inputs hashed. Verifier runs against the operations clean-room network in four seconds. §7 12-step path exercised. §10.12 exit-code contract honored. §4.4.2 deployment-intent attribute family in scope — confirm `intent` value during database review.*

She walked over to the wall screen — the gas-distribution map. Twelve hundred miles of distribution pipe. Several thousand methane sensors. The dashboard rendered every active prediction as a small dot. Three dots glowed amber. Most of the map was green.

"How many predictions a day?"

Esme answered. "On a normal day, four to six thousand inferences. About a hundred and fifty cross the 0.4 threshold and require dispatcher action. Maybe three to eight result in actual dispatch."

"And the model has been live for nine months?"

"Live for nine months. About 1.6 million inferences in the chain so far. Two hundred and forty actual dispatches. Eighty-something true positives — small leaks, mostly. Two of those would have escalated if we hadn't dispatched."

"And what was the false-positive rate before the chain?"

"There was no chain before. We started chained from day one. The vendor required it for the contract."

Dawn wrote: *Chained from day one. 1.6M inferences. ~240 dispatches. ~80 true positives. Two would have escalated. Public-safety value already demonstrable. §10.13 evidentiary-artifacts retention will care about the SDK-version manifest and HSM-config history across the nine-month window.*

Esme led them off the operations floor and back through two of the three checkpoints to the clean room. Coffee. Whiteboard. Network jacks.

Dawn pulled up a chair.

"Let's split. Raj — historian and AI ledger. Diana — IAM, all three sides. Mike and Chen — pipelines and the AI service. Elena — CIS and Salesforce, with the §10.16 SaaS-edge connector lens. Luis — logs and ops, with the §10.26 reference-verifier-distribution lens. Tom — sit with Esme, work the NERC CIP and PHMSA mappings into the §10.19 chain-coverage map. Reconvene at noon."

They split.

---

## 🧠 10:00 AM — Database Deep Dive (AI Ledger and PI Historian)

Raj had a corner of the clean room and two screens. One showed the AI-side ledger — backed by a Postgres instance reachable from the clean-room network. The other showed a screen-share from one of Esme's engineers — the OSIsoft PI historian's PI Server admin console. Raj worked them in parallel.

### The AI ledger

Raj started with the chain. Append-only by design per spec §10.3 — both at the application level (the codebase contains no UPDATE or DELETE statements on the events or daily_seals tables) and at the database role level (the ledger writer's Postgres role had INSERT and SELECT only; UPDATE, DELETE, and TRUNCATE were revoked). The Vidimus SDK signs each entry per §4.1. HMAC-SHA-256 chain links entry N to entry N-1 with HKDF-per-tenant key binding under §4.1's `info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)`. The IKM is 32 bytes per §10.6 minimum. The `key_fingerprint` is the 16-byte truncation `SHA-256(utf8("pacific-crescent") || ikm)[:16]` per §3 and §4.1.

Daily Ed25519 seals close out each day's chain on the on-prem Thales Luna PCIe HSM in the operations control center — FIPS 140-2 Level 3 conformant per §10.5, on a CIP-categorized network segment, air-gapped from corporate IT. The seal job operator role grants `sign` only; `extract`, `delete`, and `import` require separate authorization, which Pacific Crescent split between the SOC-team operator and the HSM administrator per §10.5 separation-of-duties.

Raj asked Esme which `sign_payload_version` the seals carried. "v1.0b for everything sealed since 2026-05-07. v1.0a for the seals before that. We did the cutover the day the spec amendment landed."

"Both verify under §7 step 11 dispatch — the verifier reads `sign_payload_version` from the seal record and reconstructs the appropriate byte form. v1.0a is the 10-line form; v1.0b is the 12-line form that additionally binds `key_versions_canon` and `hex(kms_handle_uris_digest)` per §4.3 v1.0b amendment. Pre-amendment seals from before v1.0a omit the field and use the 6-line form. Your verifier handles all three."

"The verifier handles all three. We tested all three on the migration day."

Raj picked a random entry from three weeks ago. Copied its ID. Ran the verifier.

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key leak-prod-2026-q1
```

He picked a random entry from yesterday. Same result. He picked the very first entry from nine months ago — the day the leak-detection model went live. Same result. The first-day entry's seal record was a v1.0 pre-amendment 6-line `sign_payload` form per §4.3 dispatch; the verifier handled the older form mechanically.

He tried to mutate one. Issued an UPDATE on the chain table directly through psql — the database accepted it because Raj had been temporarily elevated to a role with UPDATE permission for this exact red-team test, with Esme present. Production roles per §10.3 do not have this permission. Raj ran the verifier on the next entry in the chain.

```
Status: FAIL
Step: 4
Reason: HMAC mismatch — entry payload does not produce
        the chained HMAC recorded in entry N+1
```

The verifier output named §7 step 9 — `payload_hash MAC mismatch at seq N` — translated through Pacific Crescent's verifier into a slightly more conversational form, but the normative reason-string prefix per §7 was preserved verbatim. Exit code 1 per §10.12.

He tried to also rewrite entry N+1 to match. The verifier failed at the daily seal:

```
Status: FAIL
Step: 9
Reason: Merkle root mismatch — recomputed root does not
        match sealed root for date 2026-04-14
```

Per §7 step 10, the streaming RFC 6962 Merkle root recomputation over the day's `payload_hash` values in `(run_id, seq)` order produced a root different from the sealed root. The leaf-prefix `0x00` and node-prefix `0x01` per RFC 6962 §2.1 prevented confusing internal nodes with leaves. Per §1.3 security definitions, the Merkle seal's second-preimage resistance is what blocked Raj from finding a different leaf set that matched the sealed root.

He rolled back his mutations. The chain returned to PASS. He noted the test in his workbook.

> **✓ Confirmation #2 — append-only behavior under direct-mutation attempt**
> The AI-side ledger is append-only in practice per §10.3, enforced at both the application level (no UPDATE/DELETE in the codebase) and the database-role level (INSERT/SELECT only for production roles). Direct database mutation under a temporarily-elevated role is technically possible, but the verifier catches it at the HMAC layer (single-entry tamper, §7 step 9) or the Merkle/seal layer (multi-entry tamper, §7 step 10). The seal HSM lives on a CIP-categorized network segment that the database engineers do not have access to, satisfying §10.5 separation of duties (sign-only role; extract/delete/import require separate authorization). To forge undetectably, an attacker would need to simultaneously compromise three independent custody layers per §1.1 Daubert "Known error rate": the tenant's IKM (held in HSM custody under §10.5/§10.6.1 RNG provenance), the institution's ledger storage (append-only with operator-side controls), and the HSM signing key (FIPS 140-2 Level 3, non-extractable). Pacific Crescent has those split across three different network zones and three different operator role-sets.

Raj also confirmed the IKM provenance: per the `master_key.generated` operational event under §10.2 and §10.6.1 RNG-source naming, the IKM was generated inside the Luna HSM by the HSM's internal CSPRNG (RNG type `"hsm.thales-luna-internal"`), conformant under §10.6.1's first listed pattern (HSM internal RNG, the highest-assurance posture). Raj noted: *§10.6.1 conformant. Not the OS-CSPRNG pattern, not the dedicated-TRNG pattern. HSM-native.*

### The PI historian

Raj turned to the second screen. Esme's engineer — a man named Aaron — was sharing his SCADA-clean workstation. The PI Server admin console was open.

"Aaron, can you show me the audit-trail configuration?"

Aaron clicked through two menus. The audit trail was enabled. Retention was set to 60 days.

"Sixty days." Raj wrote it down. "Was that the default, or did you set it?"

Aaron looked at Esme. Esme nodded.

"That's what we set it to when the historian came in. Disk space was the constraint. We sized the audit-trail volume against the retention. Sixty days fit."

"And after sixty days?"

"After sixty days the audit trail rolls. Who-changed-what for any value older than sixty days is unrecoverable from the historian itself."

"And writes — who can write to PI?"

"PI Server uses its own identity model. PIWorld is the default group. We have three engineers with `PIWorld\db_admin` who can write or override values. The historian service account writes the normal stream. The three engineers have override authority on top of that."

Raj wrote: *PI historian. Audit trail 60 days. Three engineers with db_admin override. After 60 days, override invisible.*

He thought for a beat.

"Has any historian sensor value from 90 or more days ago been edited?"

Aaron looked at Esme. Esme looked back at him.

"I don't know. I would have to check, and after 60 days I can't check."

"Could it be?"

"Yes."

"Three people could edit a methane sensor reading from 90 days ago, and the historian itself would not retain a record of who edited it or when."

"Yes."

Raj leaned back. "And the AI service ingests methane readings from PI."

"Yes."

"In real time, as they come in?"

"Yes."

"And the AI's chain entry references the methane reading by sensor ID and timestamp, not by content hash?"

Aaron looked at Esme again. "That's correct."

> **⚠️ Surprise #1 — PI historian audit trail and §10.13 / §10.19 implications**
> The OSIsoft PI historian's audit trail is enabled but with 60-day retention. After 60 days, sensor-value changes are unrecoverable from the historian. Three engineers hold `PIWorld\db_admin` with override authority. A methane sensor reading from 90 days ago could have been edited at any time after the original write, and there is no way to determine whether that occurred. The AI service references sensor data by ID and timestamp, not by content hash, so the chain captures what the AI saw — not what the sensor actually measured. Per spec §10.13 evidentiary-artifacts retention, institutions whose chain entries may enter litigation must retain documentary evidence (SDK version manifest, HSM configuration history, daily seal-job logs, change-management records) for the chain-data retention period; PI's 60-day window is far shorter than the 7-year FFIEC chain retention horizon and far shorter than the 3-year NERC CIP-008-6 / CIP-009-6 evidentiary minimum, let alone the longer windows for grid-impacting events. Per spec §10.19 chain-coverage map, the PI historian belongs in the "institutional systems not yet chain-instrumented" category with the rollout posture named (deferred pending Phase 2 disk re-sizing) and the evidentiary substitute (the 60-day rolling audit trail) named alongside its retention boundary.

Raj wrote in his workbook: *Same shape as Stelvio's Proficy historian. Different vendor, same mutability surface, longer retention than Stelvio's row-level (which had none) but still finite. The AI trusts whatever PI gives it. §1.2 epistemic-scope language — the chain proves what the AI said, the chain does not prove what the upstream sensor actually measured. §10.19 chain-coverage map: PI in "institutional systems not yet chain-instrumented" with named substitute. Document the boundary.*

He moved on to the GE iFIX SCADA backing store.

### The iFIX SCADA backing store

Aaron pulled up the SQL Server instance behind GE iFIX. The HMI itself runs on the operator workstations, but operator notes — the free-text notes an operator types when responding to an alarm — land in a SQL Server table called `iFIXAlarmAck`.

Raj asked for the schema. The table had `alarm_id`, `ack_user`, `ack_timestamp`, `ack_note`, and `last_modified`.

"Last modified. Does that mean what I think it means?"

Aaron clicked into the table definition. "There's a trigger that updates `last_modified` on UPDATE. There is no trigger that captures what the prior value was. So if I edit my ack_note on an alarm from last Tuesday, `last_modified` will tick to today, but the original note is gone."

Raj wrote: *iFIX. Operator ack notes mutable. Trigger captures last_modified, not prior value. UPDATE silently overwrites. Same shape problem §10.3 names — UPDATE permission on an audit-bearing table is the violation, regardless of whether the application discipline says "we don't UPDATE." §10.3 application-level enforcement is "the codebase contains no UPDATE or DELETE statements on the events or daily_seals tables" and database-role level enforcement is "UPDATE, DELETE, and TRUNCATE permissions SHOULD be revoked." iFIX violates both at the alarm-ack table, since the product itself UPDATEs.*

"And there's an audit log somewhere?"

"GE iFIX writes an INSERT-only audit log to a separate table for the initial ack. UPDATE is not captured anywhere."

> **⚠️ Surprise #2 — GE iFIX UPDATE audit gap**
> GE iFIX SCADA stores operator alarm-acknowledgment notes in a SQL Server backing table. The audit log captures the original INSERT. Subsequent UPDATEs to the ack note are not captured anywhere. An operator can revise what they said about an alarm response after the fact, and there is no record that the revision occurred. This is a §10.3 append-only-enforcement gap by analogy: the iFIX product's own audit shape does not satisfy what §10.3 requires of a chain-of-custody system at the application or role level. Migrating iFIX behind a chain-of-custody-aware capture layer, or layering a third-party append-only audit add-on (the SCADA vendors offer one), would close the gap. Per §10.19 chain-coverage map, iFIX belongs in the "institutional systems not yet chain-instrumented" category; the evidentiary substitute (the INSERT-only iFIX log) is named alongside its known limitation (UPDATEs invisible).

Raj closed the screen-share. He had enough for the morning.

---

## 🔐 11:00 AM — IAM Review (AI Side, OT Side, Customer Side)

Diana had a workbook that walked through identity, access, and credential rotation — once for the AI side, once for the OT side, once for the customer side. Three columns. She filled them in the same order.

### AI side IAM

Every credential the leak-detection service uses — Postgres creds for the chain backing store, methane-sensor data feed creds, weather-API keys, the Luna HSM PIN for the daily seal, the dispatcher application's service account — every rotation was a chain entry. `event.type = credential.rotated`, `chain_kind = "audit"` per §3 enumeration, with rotator identity, rotation reason, and the new key fingerprint. The key-fingerprint discipline aligns with §10.1 reconciliation (fingerprint reconciliation runs daily at Pacific Crescent, more aggressive than the §10.1 RECOMMENDED weekly cadence).

Diana picked the last six rotations from the chain. Verified each per §7 12-step procedure. All PASS.

She asked Esme for the most recent rotation. Esme pulled it up.

```
herald-verify --tenant=pacific-crescent --service=pipeline-leak-detection \
              --event-type=credential.rotated \
              --date-range=2026-04-15:2026-05-05
```

Two entries returned. One Postgres cred and one weather-API key. Both PASS, exit code 0 per §10.12.

Diana also asked about the dispatcher application identity itself — the user_id that gets recorded when a dispatcher dismisses or actions a prediction.

Esme paused for a beat. "The dispatcher application authenticates the user against our Active Directory. The application records the user's display name as the dispatcher_id in the chain entry. So when the chain says `dispatcher: M.Reyes`, that's Marcus's display name from AD."

Diana wrote it down. "Display name, not federated SSO identity?"

"Display name."

"And display names can change."

"They can. AD doesn't enforce uniqueness on display names if a manager edits one. The underlying SID is unique. We chain the display name, not the SID."

Diana wrote: *Dispatcher chain captures display name, not federated SSO identity. SID is the durable identifier. The chain entry would survive a display-name change — the per-event MAC under §4.1 covers whatever bytes the SDK emitted at capture time, byte-stable forever — but the human-readable label would become ambiguous on a post-hoc display-name change. The chain integrity holds; the operational interpretation degrades.*

> **✓ Confirmation #3 — credential rotation under chain**
> Credential rotation on the leak-detection AI service is captured in the chain with rotator identity, fingerprint, and reason. Two rotations sampled in the most recent window, both PASS. Eight rotations sampled across the past nine months, all PASS. The AI service account itself is not shared, has MFA, and rotates every 90 days under chain. The IKM rotation discipline aligns with §10.10 (rotation crossing the seal boundary): when an IKM rotation completes near a seal-time boundary, late-arriving events under the old IKM are included in the next day's seal as late-binding entries per §4.2.2, and the day-after seal's `key_versions` field lists both generations per §10.10. Verifier handles via per-entry `key_version` lookup at §7 step 7 with no special-case logic. The IKM-registry retention is conformant under §10.9 — retired IKMs remain queryable as long as any chain entry stamped with that `key_version` is retained, which Pacific Crescent's CC8.1 names as "the longer of 7 years or the chain-data retention period."

> **⚠️ Surprise #3 — dispatcher_id binding to AD display name (operational, not integrity)**
> The dispatcher_id field captured in chain entries records the user's Active Directory display name, not their durable SID or federated SSO identity. A display-name change in AD — for example after a marriage or a department transfer — would render historical chain entries human-ambiguous, although the underlying chain integrity is unaffected (the per-event MAC under §4.1 covers the bytes as captured, including whatever string was the display name at that moment). Federated SSO with SID-binding would resolve it. The finding is a control-completeness gap, not a chain-integrity gap — the institution's CC8.1 control description per §10.18 should name the dispatcher-identity binding mechanism (display name vs SID) so a SOC-engagement reviewer reads the choice mechanically rather than inferring it from architecture. Phase 2 remediation: switch to SID-binding under the same `audit.*` namespace; the §4.4 attribute set lets institutions name the binding without a spec change.

### OT side IAM

Diana asked Esme for a screen-share from one of the HMI workstations on the operations floor. Aaron came back on the line.

The login screen showed `Operator_ControlRoom` as the account.

"Who is `Operator_ControlRoom`?"

Aaron's voice was slightly flat. "All three control-room shifts. Six operators on rotation. Same password. We rotate the password every 180 days."

"MFA?"

"No. Smart-card readers are on the procurement plan."

"Procurement plan timeline?"

"Phase 2. Twelve months out."

Diana wrote: *iFIX HMI. Shared `Operator_ControlRoom` account. Six operators, same password. No MFA. Smart cards 12 months out. The shared-account shape compounds the §10.3-style mutability problem at the iFIX backing store — when an UPDATE to `iFIXAlarmAck` is silent, six possible operators could have done it, and the auditor cannot narrow further. NERC CIP-007 systems-security-management asks for unique-user accountability; the shared account is non-conformant against CIP-007 standalone, and the chain cannot recover what NERC CIP-007 itself does not deliver upstream.*

She kept going. "AMI head-end. Itron OpenWay. Who has override authority?"

Aaron pulled up a different screen. "We have a `meter_data_engineer` role. Two people. They can override a meter reading if they decide it's malfunctioning."

"Override means edit a reading after it was recorded?"

"Override means submit a corrected reading. The original reading stays in the head-end audit log. The corrected reading replaces it in the meter data management system."

"And the head-end audit log captures the override?"

"Yes."

"With what fields?"

Aaron clicked through. "User ID, timestamp, meter ID, original value, override value, reason code."

"Reason code is mandatory?"

"It's a free-text field."

"Is the reason free-text mandatory or just present?"

"Mandatory but free-text."

Diana wrote: *AMI head-end. Two people with `meter_data_engineer` override authority. Override captured with original + corrected + reason. Reason is free-text — same shape as Pacific Crescent's iFIX ack notes, but at least the original value is preserved here, unlike iFIX. Override captured but not chain-coupled — the override is in the head-end's own log, not in the AI-side chain that uses the meter data downstream.*

She kept going. "What version of OpenWay are you on?"

Aaron checked. "5.2."

"OpenWay 5.4 ships an integrity-checking option. You're version-locked behind it."

"5.4 was released last fall. We have it on the upgrade plan for next year. There are some downstream system compatibility issues — our meter data management is on the older Itron interface."

> **⚠️ Surprise #4 — Itron OpenWay 5.2 / 5.4 integrity-checking option**
> The Itron OpenWay AMI head-end is on version 5.2. Version 5.4 introduced an integrity-checking option that would cryptographically attest meter readings end-to-end. Pacific Crescent has the upgrade on the roadmap but is constrained by downstream system compatibility (the meter data management is on the older Itron interface). Two `meter_data_engineer` accounts have override authority on individual meter readings; the override is captured with original value, corrected value, and a free-text reason, but the chain of custody from meter to head-end to MDM is not cryptographically attested. Per §1.2 epistemic-scope language, the chain proves what the AI saw at the AI service boundary; the chain does not prove the meter's emitted value matches the sensor's actual measurement. Extending §1.2's "what the chain proves" leftward — to the sensor itself — is exactly what OpenWay 5.4 enables. Phase 2 closes this gap. Per §10.19 chain-coverage map, the AMI head-end belongs in the "institutional systems not yet chain-instrumented" category with the rollout posture named (Phase 2 — AMI 5.4 upgrade pending MDM compatibility) and the evidentiary substitute (the head-end's own override audit log, free-text reason) named.

> **⚠️ Surprise #5 — shared HMI account on the operations floor**
> The GE iFIX HMI workstations on the operations floor use a shared `Operator_ControlRoom` account across all three control-room shifts. Six operators rotate through it. Password rotates every 180 days. No MFA. The same shared identity is what stamps operator alarm-acknowledgment notes in `iFIXAlarmAck`. Smart-card readers are on the 12-month roadmap. Per spec §10.18 CC8.1 cross-referencing, Pacific Crescent's CC8.1 control description should name the shared-account posture, the migration timeline, and the cross-reference to §10.3 (the shape the spec requires of an append-only audit-bearing system). The shared account is a NERC CIP-007 unique-user-accountability finding standalone; the chain-of-custody implication is that even if iFIX captured UPDATEs, the captured user_id would still be the shared role rather than the actual person. Two non-conformances stack here; the chain layer cannot recover what the upstream IAM layer does not deliver.

### Customer side IAM (CIS, OMS, Salesforce)

Diana shifted to the customer-side stack. The customized Oracle CIS, the OMS, the Salesforce CRM. She went through them quickly.

CIS had an internal user model — call-center reps had role-based access, with a small admin group that could adjust accounts. Audit trail was in an Oracle table that kept 18 months. Adjustments to billing records were logged with user ID and reason code. Adjustments to interaction notes were not.

OMS — the work-order system — had its own user model that mostly federated to AD via SAML. Work orders had version history. Field crew updates came in via tablets that authenticated with AD; some updates came in via paper forms transcribed by dispatchers, in which case the dispatcher's user_id was the one stamped on the work order.

Salesforce was Salesforce. Field-level audit on some fields. Not on others. Backups, not version history. Same shape as the diary baseline two weeks ago and same shape as Helmstad's Salesforce piece. Diana flagged: *Salesforce → §10.16 SaaS-edge connector territory if Pacific Crescent ever stands a mirror for the leak-detection AI to ingest customer-call data. Today they don't, so §10.16 is not yet load-bearing — but if Phase 4 stands the mirror, the §10.16 four-number lag posture (median, 95th-percentile SLO, alerting threshold, RTO) becomes the entry-fee for that connector's conformance. Imprecise lag wording will be a non-conformance per §10.16 severity-classification clause, not a Nit. Same posture Northbridge's Salesforce mirror landed two weeks ago.*

Diana wrote a half-page summary in her workbook. The customer side wasn't pretty, but it wasn't NERC's primary concern either. That was the PUC's concern.

---

## 🧪 12:00 PM — Lunch (Dawn and Tom in the cafeteria)

The cafeteria sat on the second floor with a view east toward Mount Hood, white today against a flat blue sky. Dawn and Tom found a corner table. The rest of the team was scattered — Diana and Mike on a sandwich run, Raj and Chen working through their morning notes.

Tom unwrapped a sandwich. "How are we framing this."

Dawn had a salad and a notebook open. "Three tiers. AI side passes. OT side mostly fails. Customer-billing side is a different audit. Per §10.19 chain-coverage map, the three tiers map to three categories — chain-instrumented institutional systems, institutional systems not yet chain-instrumented, and out-of-scope systems with named substitutes."

"NERC won't care about the customer side."

"NERC won't care about Salesforce. Right. NERC scopes to BES Cyber Systems per CIP-002 categorization. The pipeline leak detection is in scope under CIP-007. The AMI head-end is in scope under CIP-005 because of how it's networked into the OT segment. The PI historian is in scope under CIP-007. The OMS is partly in scope — work orders that touch BES assets are in, work orders for residential gas leak responses are in under the public-safety overlay even if they're not strictly BES."

"And the CIS and Salesforce —"

"PUC scope. PHMSA might care about parts of OMS for incident reconstruction. Salesforce mostly nobody cares about until a regulator subpoenas customer-service interactions in a litigation."

Tom nodded. "So the report has to split scope."

"Three columns. NERC + PHMSA in the first column. PUC in the second. Out-of-scope but operationally relevant in the third. That's how Esme's already thinking about it. And the §10.19 chain-coverage map gives all three columns a single artifact the WECC auditor and the PHMSA inspector and the WUTC regulator can read as the same picture."

Tom took a bite of his sandwich. "What's the public-safety angle."

Dawn looked at her notebook. "If a leak alert is dismissed and a week later a house explodes, what evidence do we have that the dismissal was reasonable at the time. The chain entry of the dismissal. The AI's confidence and inputs at that moment. The dispatcher's reason code. The dispatcher's identity. The model_id and version per §4.4 normative `gen_ai.request.model` and `gen_ai.response.model` requirement. All of that is in the chain, sealed daily under §4.2/§4.3 with the v1.0b 12-line `sign_payload` form. Per §1.2 the chain proves the AI said X at time T and the record was not tampered with after capture; that's what we get to hand the regulator."

"And if the inputs were tampered with upstream of the chain?"

"Then we're back to PI. The chain captures what the AI saw. The chain doesn't capture whether what the AI saw was true. Per §1.2 epistemic-scope, the chain does NOT prove (c) factual accuracy of the AI's statement, (d) policy compliance, or (e) absence of bias. Upstream-input authenticity sits closer to (c) — the chain is silent on it by design. The PI historian is the upstream weakness, and §10.19 names PI as 'institutional systems not yet chain-instrumented' with the audit-trail-rolls-at-60-days substitute."

Tom finished his sandwich. "That's the line."

"That's the line. Esme already knows it. We're going to write it down so her CEO knows it and her PUC commissioners know it and her insurance carrier knows it. And the §1.2 language is the controlled vocabulary — it lets every reader land on the same understanding without re-deriving it."

Dawn closed her notebook.

---

## 🔄 1:00 PM — A Real Alert in the Control Room

The team was back in the clean room and just settling in when Esme stuck her head through the door. "We have an alert. Different pipeline. SE-Brentwood. If you want to watch a live one through to dispatch, follow me."

The team grabbed laptops and followed.

On the operations floor, Marcus was no longer at the gas-distribution station. A different dispatcher — name tag *Yolanda* — was in his seat. The AI dashboard had a single bright amber dot pulsing on the SE Portland map. Confidence reading 0.78. Methane uptick 4.6 ppm above background. Soil saturation low. No recent precipitation. Wind moderate, steady direction.

Yolanda was already on the phone. "Brentwood crew, this is dispatch, copy?"

A voice came back. "Brentwood, copy dispatch."

"We have a 0.78 confidence leak prediction on segment SE-Brentwood-12. Methane 4.6 above background. Crew respond to coordinates I'm sending now."

"Copy. Three minutes."

Yolanda hung up. She tapped *dispatch crew* on her console. A reason code dropdown opened. She picked `model confidence above 0.75 threshold, no environmental confound`. Submit.

Mike was already at his laptop. The chain entry hit the ledger. He filtered to it.

```
herald-verify --tenant=pacific-crescent --service=pipeline-leak-detection \
              --date=2026-05-05 --entry-id=2026-05-05-Brentwood-12-dispatch
```

Four seconds.

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key leak-prod-2026-q1
```

Mike turned the laptop toward Dawn and Esme. "Live alert, dispatched, sealed, verified. Twelve steps, four seconds, on the clean-room network. The §7 verifier procedure exercised end-to-end on a real-time entry — JCS pre-flight self-test (Round-17 NIST-P2 amendment), file-header pre-flight, per-event walk through structural / IKM-lookup / fingerprint / MAC-recompute, per-day Merkle and signature, GenAI completeness check at step 12a since the entry carries `gen_ai.*` attributes."

Esme nodded once. She didn't say anything but her shoulders dropped a half-inch.

The team watched the dashboard for the next twenty-two minutes. A separate screen showed the work-order ID that had auto-opened in the OMS. The crew arrived at the location. A second voice came over the dispatch line.

"Brentwood crew on site. We have positive methane reading on the residential service line, three feet east of the meter. Small leak. Initiating standard isolation."

Yolanda put her hand to her headset. "Copy positive small leak, isolation in progress. Updating work order."

Esme's phone rang. She stepped two paces away to take it. The team heard her say "small leak, isolation, no evac" and "yes, I'll send a note to the comms team." She hung up.

She came back. "PUC liaison. They get a courtesy notice on every dispatch that escalates to crew on site. Standard."

Dawn wrote in her notebook: *Live alert from prediction to crew on site to confirmed real leak in 22 minutes. Chain captured the AI's call and the dispatcher's decision. Field crew finding will land in chain via follow-up entry once crew leader logs it. The §4.4 `parent_run_id` / `parent_seq` discipline will link the follow-up to the original prediction the same way §10.11 binds an ECOA translation to its parent decision.*

Mike asked Esme: "When does the field finding hit the chain?"

"When the crew leader closes the work order in the OMS. The OMS pushes a follow-up event to the leak-detection service, and the service writes a follow-up chain entry referencing the original prediction. Outcome — true positive, false positive, equipment fault, environmental confound. Sealed in the daily seal that night."

"How long until the crew leader closes the work order?"

"Today probably six to eight hours from now. Sometimes longer if it's a complex repair."

Dawn wrote: *Follow-up entry mechanism exists. Crew leader closes OMS work order. OMS pushes outcome to AI service. AI service writes follow-up chain entry referencing original prediction via `parent_run_id` / `parent_seq`. Sealed that night under v1.0b 12-line `sign_payload` per §4.3. Late-binding flag per §4.4 / §4.2.2 if the follow-up arrives after the day's seal closes — handled mechanically by the verifier as PASS-with-anomaly per §7 late-binding-entry-reporting.*

> **✓ Confirmation #4 — live alert end-to-end with §7 PASS**
> A live, real alert observed end-to-end on the operations floor: AI prediction at 0.78 confidence, dispatcher action with reason code, work order auto-opened in OMS, crew on site within 22 minutes, confirmed real leak. Chain captured the prediction and the dispatcher decision in real time, verifier returned PASS in four seconds for the dispatch entry — §7 12-step procedure, exit code 0 per §10.12. Follow-up chain entry mechanism exists for the field-crew outcome and will land in tonight's seal under §4.3 v1.0b 12-line `sign_payload` form. Late-binding handling via §4.4.`ffiec.chain.late_binding` and §4.2.2 day-boundary semantics is mechanical — the late-binding entry's `received_at` exceeds the seal-time of the day it was captured under, so it lands in the next day's seal with `late_binding_count` incremented and the verifier reports it as PASS-with-anomaly rather than as an integrity violation.

The team filed back to the clean room.

---

## 🧬 2:00 PM — Pipeline Reality (Chen on the Sensor-to-AI Path)

Chen had been working through the data path for an hour by the time the rest of the team came back from the control room. He had a diagram on the whiteboard.

```
methane sensor -> cellular AMI feed -> PI historian -> AI ingestion adapter -> leak-detection service -> chain entry
```

He pointed. "Here is where the chain begins." He tapped *leak-detection service*. "Everything to the left of that arrow is unauthenticated."

Dawn read the diagram. "The methane sensor itself. It just emits a value over the cellular link?"

"It emits a value. The value goes to the AMI feed, which writes it into PI, which the AI ingestion adapter reads. There's no signature on the sensor's emission. There's no signature when PI receives it. The first cryptographic operation in the path is the chain entry the AI service writes after it has already trusted the PI value. Per §4.1 'where this primitive lives' — inside the application process running the SDK, on the bank's host. The MAC is computed at the moment of event capture, before the event leaves the host. Pacific Crescent's host is the leak-detection service; everything to the left of that host is outside §4.1 protection."

"So if PI is tampered with —"

"The chain captures whatever PI says. Faithfully. The chain entry will say `methane_value: 4.6 ppm` because that's what the AI ingested. Whether 4.6 was the real reading at the sensor or whether someone overrode it in PI is invisible to the chain. §1.2 epistemic-scope — the chain proves what the AI said, not whether what the AI said was true."

"Same shape as Stelvio's historian boundary."

"Exactly the same shape. Different industry, different vendor, identical boundary problem. The AI trusts what its upstream gave it. The chain attests to the AI's behavior. The chain does not attest to the upstream's behavior. Per §10.19 chain-coverage map, the upstream gets named explicitly — institutional systems not yet chain-instrumented (PI), the rollout posture (Phase 2), and the evidentiary substitute (PI's 60-day rolling audit log, which does not extend to incidents older than 60 days)."

Chen wrote in the corner of the whiteboard: *Trust boundary = first cryptographic operation. Pacific Crescent's first cryptographic operation is at the AI service. Everything upstream is trust-by-policy. §1.2 epistemic-scope language — the chain proves (a) what the AI said at time T and (b) that the record was not tampered with after capture. The chain does NOT prove (c) factual accuracy of the AI's input or output. Upstream-input authenticity falls under (c).*

Dawn wrote in her notebook: *Sensor-to-AI path is unauthenticated for the first three hops. Documented as a Phase 2 remediation: hash the value at the sensor side, carry the hash forward through PI, verify at AI ingestion. Itron OpenWay 5.4 enables it natively for the AMI portion. PI Server has integrity-extension options that Pacific Crescent has not turned on. Once the chain extends leftward, §1.2's "what the chain proves" reaches the sensor — closer to the public-safety question than the current AI-service boundary allows. The §10.19 map will redraw at that point — PI moves from "not yet chain-instrumented" to "chain-instrumented institutional system."*

> **⚠️ Surprise #6 — sensor-to-AI path unauthenticated**
> The data path from methane sensor to AI ingestion is unauthenticated for the first three hops. The chain begins at the leak-detection service, which means the chain captures what the AI saw, not what the sensor actually measured. If the PI historian were tampered with — by one of the three engineers with override authority — the AI's chain entry would faithfully record the tampered value. The chain would verify PASS (the §7 procedure has nothing to compare against upstream of the chain's first cryptographic operation). The forensic trail would not detect the tamper. Per spec §1.2 epistemic-scope, this is exactly what the chain does not prove (the chain proves what the AI said, not whether what the AI said was true). Per §10.19 chain-coverage map, the sensor / AMI / PI hops belong in the "institutional systems not yet chain-instrumented" category with the rollout posture (Phase 2 — AMI 5.4, PI integrity extension) and the evidentiary substitute (PI's rolling audit log + AMI's override log) named.

Chen pointed at the diagram one more time. "If you fix the AMI side and the PI side, the trust boundary moves left to the sensor itself. That's where it should be in a public-safety AI."

> **⚠️ Surprise #6a — `audit.connector_source.*` family by analogy for PI ingestion**
> Per §4.4.6 SaaS-edge connector source attribution, when a connector subscribes to a source-platform's change stream and emits chain entries derived from source-platform events, the chain entries must carry the `audit.connector_source.*` attribute family — `audit.connector_source.system`, `audit.connector_source.replay_id`, `audit.connector_source.commit_timestamp`, `audit.connector_source.commit_user`, `audit.connector_source.lag_observed_ms`, `audit.connector_source.change_kind`. The PI-to-AI ingestion adapter is the same shape as a §4.4.6 mirror connector, but Pacific Crescent's ingestion adapter does not emit the `audit.connector_source.*` attributes. The closest source-side identifier is the PI sensor ID + PI write timestamp; binding those into the chain entry under `audit.connector_source.replay_id` and `audit.connector_source.commit_timestamp` would let an examiner cross-reference the chain entry to PI's audit trail directly. Phase 2 work item — emit `audit.connector_source.*` from the PI ingestion adapter so the chain at the AI boundary already carries the upstream-correlation handles, even before the chain is extended further left to the sensor itself. Per §4.4.6 stable run_id discipline, the connector-emitted entries' `ffiec.chain.run_id` must derive from a stable source-side identifier (the methane sensor ID + the PI write epoch is one institution-side choice; the institution's CC8.1 names the canonicalization), not from a per-process UUID — that lets an examiner enumerate every chain entry tied to a given sensor by `run_id` alone, surviving connector restarts.

---

## 🧪 2:30 PM — The AI Vendor Handover and Wire-Form Discipline (Mike on the §10.21 / §10.20 lens)

Mike had been tracking a parallel thread for the past hour — the AI model itself, the vendor delivery, the training-data retention, and the wire-form discipline of the chain entries on the OTLP path between the SDK and the TesseraSeal ledger. Pacific Crescent did not build the leak-detection model in-house. Schneider Electric's Industrial Cybersecurity AI division delivered the model nine months ago under a model-supply contract, with the model card, a fairness-audit report (covered for the leak-detection AI under the EU AI Act high-risk shadow regime even though the deployment is US-only — Pacific Crescent's procurement standard required the model card and fairness audit regardless of jurisdiction), and a training-data manifest naming the methane-sensor + weather + soil-saturation training shards. Spec §10.21 (cross-vendor model-handover schema) is the lens.

Mike opened the chain entry that captured the original model handover — the deployer-side `audit.model_handover.*` attribute family per §10.21. The chain entry carried `audit.model_handover.provider = "schneider-industrial-cybersecurity-ai"`, `audit.model_handover.model_id`, `audit.model_handover.model_version`, `audit.model_handover.model_artifact_sha256` (the SHA-256 of the canonicalized model archive), `audit.model_handover.model_card_sha256`, `audit.model_handover.fairness_audit_report_sha256`. Per Round-17 M&A-G2, the Round-17 close-out also requires `audit.model_handover.contract_id`, `audit.model_handover.contract_version`, and `audit.model_handover.contract_hash_sha256` when the handover happens under a written supply contract — Pacific Crescent emits all three. The contract-hash binding is what advances the post-close evidence posture from "chain-plus-contract-binder (procedural)" to "chain-plus-bound-contract (cryptographic)" — a post-close auditor answers "which contract version governed this delivery?" from the chain alone.

Mike asked Esme about the training-data retention floor. Per §10.20 normative posture, the party retaining training-data shards must set the retention floor at the longest active deployment window across all models trained on the data, plus an investigation buffer.

"Schneider's contract names a five-year training-data retention floor," Esme said. "Our deployment window is open-ended — we keep the model running until the next major upgrade, expected in 24 months. Five years is well above the deployment window plus a 60-90 day investigation buffer, so we're conformant on §10.20."

"And the `audit.model_handover.training_data_retention_floor_days` field on the handover entry?"

"Populated. 1825. Cross-referenced in our CC8.1 to the contract clause."

Mike confirmed: per Round-17 M&A-P2, the `audit.model_handover.training_shard_manifest_sha256` field is also populated — Schneider published a canonical sorted list of training-shard hashes (newline-joined ASCII, no trailing newline) alongside the model card, and the SHA-256 over that manifest is bound on the handover chain entry. A post-close auditor can recompute the manifest from the surviving shards and confirm the provider delivered what was committed. Closes the deal-window-lookback gap.

Dawn wrote: *§10.21 cross-vendor model-handover conformant — full attribute family populated including Round-17 M&A-G2 contract-binding triple and Round-17 M&A-P2 training-shard manifest hash. §10.20 training-data retention floor (1825 days) exceeds the deployment-window-plus-buffer requirement.*

Mike turned to the wire-form discipline. The OTLP path between the leak-detection SDK and the TesseraSeal ledger ran over OTLP/gRPC with TLS 1.3 — per §5.1 transport encryption, "transport encryption MUST be applied between the SDK and the ledger" (the spec's normative wire-encryption posture). The Resource attributes on every OTLP request carried `ffiec.chain.spec = "v1.0"`, `service.name = "pipeline-leak-detection"`, `service.version = "0.13.0"` (the Vidimus SDK build), `ffiec.chain.posture = "ffiec"` (FFIEC-conformance posture per §4.1.2, not vendor-flag mode), `ffiec.chain.format_version = "v1"` per §4.4.3 OTLP transport identification. The gRPC metadata carried `ffiec-chain-spec: v1.0` and `ffiec-chain-posture: ffiec` per §4.4.3 RECOMMENDED out-of-band signals; receiver dispatch happens once per OTLP request before per-entry decode, and the metadata cross-checks against the Resource attributes on body decode.

Mike asked whether collectors mutate any chain attributes between SDK and ledger. "Per §4.4 OTLP-collector transformation pass-through, the chain attributes (`ffiec.chain.*`) and the integrity-bound payload attributes (`gen_ai.*`, `tool.*`, `audit.*`, the OTel envelope per §5) must pass through collector transformations unchanged. A collector that mutates any of these produces non-conformant output downstream."

Esme confirmed Pacific Crescent runs a single OTel collector between the SDK and the ledger, configured for transparent pass-through on the chain pipeline branch with no redaction, no sampling, no severity-filter on chain-of-custody traffic — per §4.4.4 severity-for-chain-of-custody-traffic, the institution's collector is exempt from severity filters on the chain branch, and the receiver stamps the TesseraSeal `QuickLogBuilder`-positioned `SeverityNumber` on receipt (per §4.4.4 receiver level-resolution, in the `9..20` INFO-to-FATAL-minus-one range), with `SeverityText = "OTLP"` per the TesseraSeal reference convention.

The redaction discipline was the next question. Per §10.22 normative posture, redaction must happen pre-MAC at the SDK boundary; post-MAC sidecar redaction is non-conformant unless the sidecar is itself a separate chain that points to a parent unredacted chain via §10.21-style cross-anchor. Pacific Crescent's leak-detection chain entries do not redact anything pre-MAC — the chain captures the full prediction context including methane sensor IDs, dispatcher reason codes, and model inputs. No `audit.redaction.*` attribute family is emitted because no redaction is occurring; the absence of the family is conformant under §10.22 (the family is REQUIRED only on chain entries whose content was redacted at the SDK boundary).

For Critical Energy Infrastructure Information (CEII), Pacific Crescent's redaction discipline differs by audience: the chain-of-custody record contains the full data; CEII redaction happens on the audience-facing reports the institution generates from the chain (the WUTC submission, the FERC notice under DOE Form OE-417 disturbance reports). The audience-facing reports are derived artifacts, not chain entries; the chain remains the integrity-bound canonical record. Per §5.2 best-evidence posture (informative), the captured JSON IS the content-bearing form; CEII redaction in audience-facing artifacts is a downstream operation that the chain does not interfere with.

Mike checked one more thing — the per-event canonical bytes. Per §5 wire format normative, the canonical-JSON form used inside `payload_hash` follows RFC 8785 (JCS); implementations MUST NOT use the OTLP protobuf encoding itself for hashing because protobuf is non-deterministic for some field types. Per §5 canonical-form exclusion rule, the canonical bytes input to the MAC EXCLUDE the chain-stamp fields (`prev_hash`, `payload_hash`, `key_version`, `key_fingerprint`, `format_version`, `mac_computed_at_utc`, `kms_handle_uri`, `algorithm`, `seq`). Pacific Crescent's SDK is the Vidimus reference implementation v0.13.0; the canonicalization passes the §7 pre-flight JCS self-test (Round-17 NIST-P2) on the verifier side at every startup.

> **✓ Confirmation #6b — wire-form, transport, redaction, and vendor-handover discipline**
> §10.21 cross-vendor model-handover schema is conformant on the leak-detection model — all required attributes populated including Round-17 M&A-G2 contract-binding triple (`contract_id`, `contract_version`, `contract_hash_sha256`) and Round-17 M&A-P2 training-shard manifest hash. §10.20 training-data retention floor (1825 days) exceeds the deployment-window-plus-investigation-buffer requirement and is bound on the handover entry. §5 wire format conformant via RFC 8785 JCS canonicalization with the §5 exclusion-rule list of chain-stamp fields. §5.1 transport encryption applied via TLS 1.3 between SDK and TesseraSeal ledger. §4.4.3 OTLP transport identification conformant — all five required Resource attributes set; recommended gRPC metadata mirrors. §4.4.4 severity-for-chain-of-custody-traffic conformant — collector exempts chain pipeline from severity filters; receiver stamps `SeverityNumber` in the `9..20` range with `SeverityText = "OTLP"`. §10.22 redaction discipline — no redaction occurring on the leak-detection chain; absence of `audit.redaction.*` is conformant; CEII redaction happens on derived audience-facing artifacts, not on chain entries, preserving the §5.2 best-evidence captured-JSON content-bearing form. §6 storage discipline conformant — chain entries persisted byte-for-byte to fsync'd Postgres before `payload_hash` is disclosed in any way (per §4.1 construction-location requirement).

Dawn wrote in her notebook: *§10.21 / §10.20 / §5 / §5.1 / §5.2 / §6 / §4.4.3 / §4.4.4 / §10.22 — wire-form and vendor-handover side of the chain is clean. Phase 2 / Phase 3 work is all upstream of the SDK or downstream of the OMS, not on the chain itself. The chain is doing what the spec says it should do.*

---

## 🏛️ 2:45 PM — Two Future-Scenario Questions (Esme on §10.24 succession and §10.15 multi-region)

Esme had two questions she'd been holding for the right moment.

"Two future scenarios. One: Pacific Crescent is in early conversation with a Northwest peer about a possible merger. Not signed, not announced, but the boards are talking. If that closes, our chain history under `tenant_id = pacific-crescent` carries forward — what's the spec say?"

Dawn answered. "Spec §10.24 entity succession. When the chain operating under `(tenant_id, run_id)` keying experiences a legal-entity change of operator — merger, acquisition, divestiture, rename, subsidiary transfer — the institution emits a `chain.entity_succession` operational event under §10.2 marking the legal-entity transition. The chain itself stays under the same `(tenant_id, run_id)` keying — the succession event is the integrity-bound record of the legal-entity change, not a re-keying. Chain entries from before the succession remain verifiable under the original entity's binding; chain entries after the succession bind under the successor entity's signature on the transfer-day seal. The event schema carries `from_entity_legal_name`, `to_entity_legal_name`, RECOMMENDED `from_entity_lei` and `to_entity_lei` (the 20-character LEI per RFC 9101 / ISO 17442), `effective_utc`, `kind`, optional `regulator_filing_id` (FERC merger filing ID for utility M&A), and `dual_signatures` per the §10.17 signatory schema with the `entity_affiliation` field discriminating from-entity and to-entity signers. The companion document `docs/m-and-a-handoff.md` is normative-supplement and provides the operational shape for merger / acquisition / divestiture / spin-off scenarios anchored against §10.24."

"And the `tenant_id` itself?"

"Optional rename. If the surviving entity operates as `tenant_id = pacific-crescent-northwest-merged` or whatever, both the old and new `tenant_id` values are recorded on the succession event itself, both bound under the seal of the transfer-day. The verifier walks both chains independently using the per-`tenant_id` HKDF binding per §4.1. If you preserve the original `tenant_id`, even simpler — chain integrity holds across the boundary with the succession event as the legal-entity-change record."

Esme wrote it down. "Question two. We have a hot/cold DR posture on the leak-detection AI service. Hot in the operations control center, cold-failover to a Bonneville Power facility 90 miles north. Today the cold side runs in receive-only mode — replication of the chain artifacts but not active SDK emission. If we ever flip to active-active across the two sites for resilience, what does the spec say?"

"Spec §10.15 multi-region resilience. A single `tenant_id` may operate across multiple regions; you operate one of two patterns. Pattern A — active-active with seal-region pinning. The seal region aggregates events from all regions for the tenant-day, computes the daily Merkle root in `(run_id, seq)` order across all regions, and produces the HSM-signed seal. Per-event MAC is region-agnostic (the HKDF binding is region-blind). Run-locality at v1.0 — a run starts and ends in the same region; cross-region run continuation is deferred to v1.1. SDK per-process region binding is the conformant enforcement model — one SDK process per region, each pinned to its region's IKM custody endpoint and ledger endpoint. Pattern A invariant 5 requires per-region event-count reconciliation per tenant-day via the `master.cross_region_replication_completed` operational event under §10.2; the seal region's count must equal the sum of regional counts. Pattern B — per-region `tenant_id`, with cross-region correlation institution-side. Pattern A reduces verifier-run count to one per audit period; Pattern B preserves per-region cryptographic isolation. Your risk-tolerance statement governs the choice."

"And the `ffiec.chain.region` attribute?"

"Optional under §4.4. RECOMMENDED for Pattern A so a Pattern A failover-incident reconstruction can identify which events came from which region under integrity binding. The attribute is in the canonical bytes per §5 — a tampered value surfaces as a MAC mismatch at §7 step 9. Single-region deployments (which is where you are today) and Pattern B deployments may omit the attribute; Pacific Crescent's CC8.1 names the single-region posture today."

Esme wrote it down. "Phase 5 conversation if we ever flip to active-active. Or if the merger lands."

Dawn wrote: *§10.24 entity succession framework documented for the prospective merger. §10.15 multi-region pattern selection deferred to a Phase 5 conversation. Both are normative spec text Pacific Crescent's CC8.1 author can cite on the day either scenario lands.*

> **✓ Confirmation #6c — succession and multi-region forward-readiness**
> §10.24 entity-succession framework documented. The companion `docs/m-and-a-handoff.md` is normative-supplement; the `chain.entity_succession` operational event under §10.2 carries the dual-signature attestation per §10.17 schema with `entity_affiliation` discriminating from-entity and to-entity signers. The chain history under `tenant_id = pacific-crescent` carries forward across a merger boundary with chain integrity preserved per §4.1 HKDF binding. §10.15 multi-region resilience patterns named for a future active-active Phase 5 conversation; Pattern A vs Pattern B selection is institution-side under the risk-tolerance statement, with §10.15 invariant 5 reconciliation via `master.cross_region_replication_completed` under §10.2 as the load-bearing evidence on the active-active path.

Esme had one more question. "DR rejoin. We have a 90-mile cold-failover site at the BPA facility. Today it's receive-only — replicates the chain artifacts but doesn't emit. If a regional event takes the operations control center offline and we have to spin up the cold side hot — what does the spec say?"

Dawn answered. "Spec §10.25 run resume and chain-tail acquisition. Once a run is opened with `(tenant_id, run_id)`, the SDK must acquire the run's chain tail before emitting the next entry, regardless of whether the run is fresh, in-process, or being resumed across a process boundary. The chain tail is the triple `(latest_seq, latest_payload_hash, key_version)` plus optional metadata. Three-place tail acquisition — in-memory state, local persistence sidecar, or ledger query (the rejoin path). When local persistence is permanently lost — disk corruption, container disposal, region failover without state replication — the SDK must query the ledger's chain-tail endpoint before emitting the next entry under the affected run. The rejoin mechanism must NOT degrade silently to genesis if the ledger is unreachable. If the ledger is unreachable AND local persistence is missing, the SDK must refuse to emit until the tail is acquired or until an operator explicitly authorizes a fresh genesis under a NEW `run_id`. A fresh genesis under the SAME `run_id` is non-conformant; the §4.4 genesis-block uniqueness rule plus the ledger's ingestion cross-check catch the resulting fork at the ledger layer. Single-writer-per-run rule — file-locked or row-locked sidecar prevents two SDK processes from racing for the same run identifier. Ledger ingestion cross-check rejects a batch whose `prev_hash` doesn't match the run's known tail with a normative reason string."

"And the cold-side ledger?"

"Per §10.25 disaster-recovery rejoin discipline, the cold-side ledger needs to be reachable at the documented chain-tail endpoint URL (named in your CC8.1 per §10.25's required CC8.1 cross-reference). The ledger replication keeps the cold-side ledger up to date; an SDK process that comes up cold queries the cold-side ledger's chain-tail endpoint, picks up the existing tail, and resumes from there. No silent re-genesis. If your operations team flips the cold side hot under documented procedure, the runbook names the procedure and the §10.25 rejoin path is the load-bearing mechanism."

Esme wrote it down. "Documented in our DR runbook today by component, but not cross-referenced to §10.25. That's a §10.18 cross-referencing Nit on the DR runbook the same shape as Luis's earlier finding on the dispatcher-application runbook."

Dawn wrote: *§10.25 run resume and chain-tail acquisition discipline documented; §10.18 cross-referencing Nit on the DR runbook (does not cite §10.25 / §4.4 genesis-block uniqueness / §4.4 ingestion cross-check). Combined with Luis's dispatcher-application runbook Nit — two §10.18 cross-referencing fixes total, both five-line additions.*

Then Esme asked the consumer-correlation question.

"§10.23 — consumer-correlation index. We don't operate one for the leak-detection AI today (the chain isn't keyed by customer; it's keyed by sensor segment and dispatch decision). But our customer-side regulators (WUTC, OPUC, the CA-CPUC sliver) increasingly ask consumer-keyed retrieval questions when they investigate complaints — 'produce all leak-detection-related dispatch decisions that affected the household at [address X] during [time window Y].'"

"Spec §10.23 closes that. If you ever stand a consumer-correlation index over chain-bound consumer-facing decisions, you operate one of two shapes. Shape 1 — chain-anchored index — each CUEC entry is itself a chain entry under `chain_kind = "operational"` per §3 enumeration. The append-only property of the chain (per §10.3) makes the index append-only by construction. Shape 2 — index attestation — the institution emits a daily `consumer_index.attestation` operational event under §10.2 carrying the index's snapshot hash and consumer count over the period. The CFPB's verifier (or the PUC's) independently recomputes the index hash from the chain and compares against the attestation. Mismatch is a control-completeness finding."

"For us today the trigger isn't ECOA-shaped. But the §10.23 mechanic applies by analogy."

"Yes. The §10.23 attribute schema (`consumer_index.consumer_id_hash`, `consumer_index.run_id`, `consumer_index.seq`, `consumer_index.relationship`) is institution-side as long as your CC8.1 names the canonicalization of the consumer identifier — typical: the lowercased premise address normalized to USPS canonicalization. Phase 4 work item if a PUC consumer-keyed retrieval question becomes recurrent."

Dawn wrote: *§10.23 consumer-correlation index by analogy — Phase 4 if PUC consumer-keyed retrieval becomes recurrent. Not load-bearing today.*

---

## 📊 3:00 PM — Reconciliation Test (Four Leak Predictions, End to End)

The team did the reconciliation test together. Esme picked four leak predictions from the prior 120 days, spanning recent and old. The team traced each one end to end — from the methane sensor reading that fed it, through the AI's prediction, through the dispatcher's decision, through the OMS work order, through the crew leader's confirmed finding.

### Prediction 1 — 2026-05-02 — SE-Powell-22 — high confidence — true positive

Three days old. AI-side chain: PASS. Verifier four seconds. §7 12-step procedure exercised; exit code 0 per §10.12.

Forward trace: dispatcher decision in chain (`dispatch crew, model confidence 0.86, no confound`), OMS work order opened automatically, crew leader closed it that evening with finding `true positive — service line leak — repaired`. Follow-up chain entry referenced the original prediction via `parent_run_id` / `parent_seq` per §4.4. PASS.

Backward trace: methane sensor reading at the time of inference — Aaron pulled it from PI within the 60-day audit retention. Original value, no edits, captured by `historian_writer` service account. Confidence the value is unaltered: high. Per §10.19 chain-coverage map, this is within the PI-substitute window the institution names as evidentiarily-usable; past 60 days it isn't.

Status: clean four-of-four end to end.

### Prediction 2 — 2026-04-21 — NW-Lovejoy-18 — moderate confidence — false positive

Two weeks old. AI-side chain: PASS. Verifier four seconds.

Forward trace: dispatcher decision (`dispatch crew, confidence 0.62, light precipitation environmental confound but above threshold`), OMS work order opened, crew leader closed it with finding `false positive — wet soil signature — no leak`. Follow-up chain entry referenced original via §4.4 parent linkage. PASS.

Backward trace: methane reading still within PI's 60-day audit retention. Original value, no edits, captured by `historian_writer`. Clean.

Status: clean four-of-four end to end.

### Prediction 3 — 2026-02-08 — SE-Holgate-7 — high confidence — true positive

Eighty-six days old. AI-side chain: PASS. Verifier four seconds.

Forward trace: dispatcher decision in chain (`dispatch crew, confidence 0.91, low soil saturation no confound`), OMS work order opened, crew leader closed it with finding `true positive — main line corrosion seep — repaired`. Follow-up chain entry referenced original. PASS.

Backward trace: methane reading from 2026-02-08 is past the PI 60-day audit retention as of today. Aaron pulled the value itself — it was there, 5.2 ppm, but the audit-trail entry that would tell us whether the value had been edited at any point was no longer recoverable.

Chen asked the question. "Could that 5.2 have been overridden between 60 days ago and now?"

Aaron paused. "I can't answer that from PI. The audit trail rolled."

"Could it have been?"

"Yes."

Status: AI-side PASS, dispatcher decision PASS, field-crew outcome PASS, sensor-side unprovable. Per §1.2 epistemic-scope — the chain still proves what the AI said at time T; what the chain does not prove (that the input the AI saw matched the sensor's actual reading) is exactly the (c) clause of §1.2.

### Prediction 4 — 2025-11-19 — NE-Killingsworth-3 — high confidence — true positive

One hundred sixty-nine days old. AI-side chain: PASS. Verifier four seconds.

Forward trace: dispatcher decision in chain (`dispatch crew, confidence 0.83, no confound`), OMS work order opened, crew leader closed it with finding... and here the trace got interesting.

Esme pulled up the OMS work-order detail. The "actual finding" field said `small leak — repaired — see attached photos`. The photos were attached as JPEGs. The work order also referenced a paper form completed by the crew leader on site.

Mike asked. "Where's the cryptographic linkage between the field crew's finding and the AI prediction?"

Esme looked at the work order. "The work order references the alert ID. The alert ID is in the original chain entry. The follow-up chain entry references the work order ID."

"So the linkage is by work-order ID."

"By ID, yes. Not by content hash."

"So if someone edited the work order's actual-finding field after the follow-up chain entry was written —"

"The chain entry would still reference the work order ID. The work order's content might have changed. The chain wouldn't catch it."

Mike paused. "Per §10.19 spec, there's an attribute family `audit.external_artifact.*` that hash-anchors external evidentiary artifacts on the chain. `audit.external_artifact.kind` (institution-named — for example `oms_work_order_finding`), `audit.external_artifact.identifier` (the OMS work-order ID), `audit.external_artifact.sha256` (SHA-256 over the canonicalized OMS work-order content at follow-up time), `audit.external_artifact.received_at_utc` (when the institution received the OMS finding), `audit.external_artifact.source_party` (`oms`), `audit.external_artifact.evidentiary_role` (`chain_of_custody_handoff`). With those six attributes on the follow-up chain entry, post-hoc edits to the OMS work-order content would be detectable from the chain alone. Without them, the linkage is by ID only — which is what we have today."

Esme wrote it down. "That's a Phase 3 remediation. The OMS exposes a stable canonicalization?"

"Institutional choice — your CC8.1 names the canonicalization. Even a JCS-canonical form over `(work_order_id, finding_text, photo_hashes_sorted, last_modified)` would do it. The hash is what binds; the content lives in the OMS under whatever retention OMS provides."

Diana spoke up. "And the original paper form?"

"In a filing cabinet at the crew dispatch yard. We scan the photos and attach them to the OMS work order. The paper form goes into a 7-year retention box."

"The work order's text in the OMS is the dispatcher's transcription of the paper form?"

"Or the crew leader's tablet entry, depending on whether the crew leader had the tablet that day. About sixty percent of work orders are tablet-entered, forty percent paper-transcribed."

Mike wrote: *OMS work-order linkage to chain is by ID, not by content hash. Paper-to-OMS transcription introduces a gap. The chain captures that the dispatch happened and that a follow-up was written. The chain does not capture whether the field finding's text in the OMS today is what the field finding actually said. §10.19 `audit.external_artifact.*` family closes it as a Phase 3 work item — six attributes, no spec change required, pure institution-side emission discipline.*

Backward trace: PI audit retention long since rolled. Sensor value present in PI but origin unprovable.

Status: AI-side PASS, dispatcher decision PASS, follow-up chain entry PASS by ID, field-crew finding text not chain-coupled, sensor-side unprovable.

### Reconciliation summary

The team wrote it on the whiteboard.

| Prediction | AI chain | Dispatcher | OMS work order | Field finding | Sensor source |
|---|---|---|---|---|---|
| 1 (3 days old) | PASS | PASS | PASS | PASS — within OMS, content-coupled by recency | PASS — within PI 60-day retention |
| 2 (14 days old) | PASS | PASS | PASS | PASS — within OMS, content-coupled by recency | PASS — within PI 60-day retention |
| 3 (86 days old) | PASS | PASS | PASS | PASS — within OMS | UNPROVABLE — past PI 60-day retention |
| 4 (169 days old) | PASS | PASS | PASS | UNPROVABLE — paper transcription gap | UNPROVABLE — past PI 60-day retention |

Dawn looked at the whiteboard.

"4 out of 4 AI-side PASS. 2 out of 4 trace back cleanly to the source sensor data. 4 out of 4 trace forward to the dispatcher decision. 2 out of 4 trace forward to a confirmed field-crew finding that we'd be willing to put in front of an investigator without a caveat. The chain works for the recent past. The legacy systems erode it as time goes on. Per §1.2 epistemic-scope, the chain's claims are bounded — what falls outside the bound is exactly what we're documenting. Not as integrity gaps, but as scope gaps. The §10.19 chain-coverage map will name them as such."

Esme looked at the whiteboard for a long beat. "That's what I expected."

"The 60-day PI retention is the single biggest forensic limit. After 60 days, the upstream half of the trail is gone. That's a NERC CIP-007 finding waiting to be written. PHMSA will care too — the integrity-management standards expect 5+ year evidence retention on incident-relevant data. Per spec §10.13 evidentiary-artifacts retention, institutions whose chain entries may enter litigation must retain documentary evidence for the chain-data retention period or longer if a litigation hold extends it; a 60-day rolling audit log on the sensor source is well below the §10.13 standard, the NERC 3-year CIP-008-6 minimum, and the longer windows for grid-impacting events."

"Yes."

"The OMS transcription gap is a public-safety-investigation finding. If a leak is dismissed and the dismissal turns out to have been wrong, the chain will tell you the dismissal was logged and what the AI's inputs were and what the dispatcher said. It won't tell you what the field crew found three weeks earlier on a different alarm if that finding was paper-transcribed and the OMS text was edited after the follow-up entry. §10.19's `audit.external_artifact.*` family closes it on the digital side by hash-anchoring the OMS finding at follow-up time."

"Yes."

> **✓ Confirmation #5 — reconciliation test against AI scope**
> Reconciliation test on four predictions traced end-to-end. AI-side: 4/4 PASS — chain integrity per §7 12-step procedure, verifier returns under 5 seconds for any single entry, twelve-step output. Dispatcher decisions: 4/4 PASS — captured with user, reason code, model inputs, model version (per §4.4 normative `gen_ai.request.model` and `gen_ai.response.model` — both populated, §7 step 12a passes). Forward trace to OMS work order: 4/4 by ID linkage. The chain holds where it was designed to hold — per §1.2 the chain proves what the AI said at time T and that the record was not tampered with after capture, and per §1.4 compositional security all three independent authentication layers (per-event HMAC, daily Merkle seal, HSM-rooted Ed25519 signature) are exercised on the AI side.

> **⚠️ Surprise #7 — erosion outside AI scope; §10.19 substitutes named**
> Reconciliation test reveals two erosion points outside the AI scope. Backward trace to sensor source — 2/4 PASS, 2/4 UNPROVABLE due to PI's 60-day audit retention. Forward trace to confirmed field-crew finding — 3/4 PASS (recent + tablet-entered), 1/4 UNPROVABLE due to paper-form transcription. The chain does what it was scoped to do; the legacy systems on either side of the chain limit how far the forensic trail can be carried. Per §10.19 chain-coverage map, both erosion points belong in the "institutional systems not yet chain-instrumented" or "third-party / external evidence with substitute" categories with the rollout posture and substitute named — PI's 60-day retention noted alongside the Phase 2 disk re-sizing remediation; OMS work-order content-coupling noted alongside the Phase 3 `audit.external_artifact.*` emission. The §10.19 family makes the gap explicit, testable, and accountable to a vendor-management auditor reviewing Pacific Crescent's chain-of-custody control rather than implicit and per-finding.

---

## 😬 3:45 PM — Friction in the Room (OT Engineering Defensive)

Esme had pulled in three OT engineering staff at the start of the 3:00 PM session — Aaron from PI, a SCADA engineer named Hugh, and the senior engineer over the AMI head-end, a woman named Pavithra. They had been listening from the back. By 3:45 they were not listening passively anymore.

Hugh spoke first. "I want to push back on the framing. You're describing GE iFIX like it's broken. iFIX is a NERC CIP-compliant HMI deployed in compliance with our CIP-007 controls. The audit log captures the INSERT. That's the standard."

Dawn put her pen down. "Hugh, I hear you. We're not writing iFIX up as broken. We're writing it up as a system whose audit semantics are INSERT-only, which means a subsequent UPDATE to an operator note doesn't generate a record. That's a finding about the audit shape, not about iFIX as a product. Per spec §10.3, append-only enforcement applies at both the application level (no UPDATE / DELETE statements on audit-bearing tables) and the database role level (UPDATE / DELETE / TRUNCATE permissions revoked). iFIX violates both because the product itself UPDATEs the alarm-ack note in place. The finding language is precise about that — we're naming the product behavior, not the operator's deployment."

"Auditors always make the same finding. It doesn't lead anywhere. The product doesn't support UPDATE auditing without a third-party add-on."

"That's right. And that's exactly what the finding will say. INSERT-only audit, third-party add-on available, evaluate against the cost of a tamper scenario in the public-safety context. Per §10.18 CC8.1 cross-referencing, your CC8.1 control description should name the iFIX append-only posture (currently INSERT-only), the add-on that would close the gap, and the spec section §10.3 that the gap pertains to. The cross-reference makes the finding discoverable to a vendor-management auditor without requiring that auditor to re-derive the spec from scratch."

Hugh sat back. He didn't say anything for a beat.

Pavithra picked it up. "AMI 5.4. We know about the integrity-checking option. We can't deploy it because our MDM is on the older interface."

"Right. And the finding will say: 'AMI 5.2 is version-locked behind an integrity-checking option that will close a finding when MDM is upgraded.' That's a roadmap item, not a violation."

"It's going to read like a violation in front of the PUC."

Esme cut in. "It's not going to. The team is identifying gaps for our roadmap, not writing a NERC violation. We've already self-disclosed two of these to the PUC liaison in the last quarter. And the §10.19 chain-coverage map has the AMI head-end in the 'not yet chain-instrumented' category with the rollout posture and the named substitute already in our CC8.1. The PUC reads the same artifact as the WECC auditor."

Pavithra and Hugh both looked at Esme. The friction in the room shifted by a quarter-turn. Esme had clearly had this conversation with them before, but never in front of outside auditors.

Aaron was the quietest of the three. He spoke last. "The 60-day PI retention. I didn't set that. I inherited it from the previous PI admin. When I came on, I asked about extending it. The disk volume isn't sized for 18 months of audit trail at our sensor density."

Dawn wrote: *PI 60-day retention disclosed by current admin as inherited. Disk sizing constraint. Phase 2 includes disk re-sizing on the historian server. Per §10.13 evidentiary-artifacts retention, the chain-data retention period plus litigation-hold extensions sets the floor; PI's current 60 days is well below the §10.13 standard and below the NERC CIP-008-6 / CIP-009-6 3-year minimum. The disk-sizing remediation is operational, not chain-architectural — extending PI to 18 months brings PI's substitute window inside §10.13's retention contract.*

"Aaron, that's a clean disclosure. We'll write it that way. The finding will be on the retention shape and the sizing constraint, not on you personally. And the §10.19 map names this as the substitute Pacific Crescent already operates while Phase 2 closes the gap."

The friction in the room subsided. The three engineers didn't relax exactly, but they stopped pushing back. Esme had defused it with one sentence and Dawn had received the disclosures cleanly.

Dawn made a note in the margin: *NERC engineers always defensive. Esme manages it well. Disclosures land cleaner because she creates the space for them. The §10.19 chain-coverage map is the artifact that makes a defensive engineer's contribution land as 'here's the substitute we operate' rather than 'here's the gap I'm responsible for.'*

---

## 🔍 4:30 PM — The Public-Safety Question

Esme had asked the engineers to step out for the last hour. The team was alone in the clean room with her. The afternoon sun was angling through the high windows. She had her elbows on the table and her chin in her hands.

She looked across at Dawn.

"Let me ask you a question I've been trying to find a clean answer to for a year."

Dawn waited.

"If a leak alert is dismissed and a house explodes a week later, what evidence do we have that the dismissal was reasonable at the time?"

The team went still.

Dawn took a long breath. "Sober answer. The chain entry of the dismissal — that gives you the dispatcher's identity, the reason code, the timestamp. The AI's confidence and inputs at that moment — that's in the same chain entry. The model_id and model_version per §4.4 normative `gen_ai.request.model` and `gen_ai.response.model` — also in the chain. The seal record for that day — ties everything to a public key with a known fingerprint via §4.3 v1.0b 12-line `sign_payload` form. Per §1.2 epistemic-scope, the chain proves (a) what the AI said at a specific time and (b) that the record was not tampered with after capture. Both are exactly what your forensic question needs."

"And that's enough?"

"That's enough to demonstrate, with cryptographic confidence, that the dismissal happened the way the chain says it happened. That's what your insurance carrier and the state regulator and the Class A NERC auditor will want to see in that scenario. Per §1.1 Daubert four-factor grounding, the chain has a concrete answer to each factor — testability under §7's ordered byte-exact verification procedure with the test-vector corpus pinning byte values; peer review under the FFIEC working-group process; known error rate under §1.1's three-layer compromise model (IKM + ledger + HSM, each compromised by different roles under separation of duties); general acceptance under the NIST-standardized primitives (FIPS 180-4 SHA-256, FIPS 186-5 Ed25519, FIPS 198-1 HMAC, RFC 5869 HKDF, RFC 8785 JCS, RFC 6962 Merkle). An expert witness laying foundation under FRE 702 has the four answers in a one-page table."

"And if the inputs to the AI were tampered with in the historian?"

Dawn paused.

"Then we cannot prove the inputs were tampered with. And we cannot prove they weren't. The historian is your weakest evidence in that scenario. The chain entry will faithfully record the inputs the AI saw. If those inputs were already tampered with when the AI ingested them, the chain reflects the tampered version. The audit trail in PI rolls at 60 days, so going back further, even the question of who-changed-what becomes unanswerable. Per §1.2 epistemic-scope, this falls in the (c) non-claim — the chain does not prove the AI's statement is factually accurate, and 'the input the AI saw matched what the sensor actually emitted' is a factual question about the AI's input. The chain is silent on it by design."

"And the dispatcher in that scenario — they made what they thought was a reasonable call based on what the AI showed them."

"They made a defensible call based on the inputs they saw. The chain proves that. The question of whether the inputs were what the sensors actually measured is a different question that the chain cannot answer."

Esme sat back. Her chin came off her hands. "That's the answer I was afraid of."

"The remediation is the Phase 2 work. Hashing at the sensor — Itron OpenWay 5.4. Hashing at the historian boundary — PI integrity extension. Once the chain extends to the sensor, the question of input authenticity becomes answerable from the chain alone. As long as the chain starts at the AI service, the answer to that question depends on the historian's discipline, and the historian's discipline is operational, not cryptographic. Per §10.19 chain-coverage map, the boundary moves left when those Phase 2 items land — PI moves from 'not yet chain-instrumented' to 'chain-instrumented institutional system,' and the AMI head-end follows once 5.4 is live."

Esme nodded. She wrote something in her own notebook.

Tom leaned forward. "Esme, the value of doing this assessment now is precisely so that you have the answer in your hand before you ever need it. The PUC is going to ask. Insurance carriers are starting to ask. If a public-safety event ever does happen, having the chain in place for the AI side and a documented Phase 2 plan for the upstream is a substantially better position than not having either of those things. Per §10.13 evidentiary-artifacts retention, the litigation-foundation evidence — SDK version manifest, SDK source-code hash with build reproducibility, HSM configuration, daily seal-job logs, change-management records, verifier output — is exactly the documentary evidence the institution's IT witness lays foundation from at deposition under FRE 901(b)(9). The chain produces the verifier output; §10.13's named artifacts produce the process. Both together are the foundation."

"I know. I'm going to be straight with you — I've been trying to fund Phase 2 for six months. Today's report is the lever I needed."

Dawn wrote in her notebook: *Public-safety question. Esme's been preparing for this question for a year. The report has to answer it directly. The line is: chain proves the AI's behavior per §1.2 (a) and (b), chain does not prove the upstream's behavior per §1.2 (c), Phase 2 closes the gap. Document. The §10.13 evidentiary-artifacts retention list is the litigation-foundation backbone — name each artifact and its retention.*

> **✓ Confirmation #6 — public-safety evidence question answerable from chain alone**
> The chain provides cryptographic evidence sufficient to demonstrate that an alarm dismissal was made on the inputs and reason recorded, by the dispatcher recorded, at the timestamp recorded. In the public-safety scenario where a regulator, insurance carrier, or litigation discovery process asks how a dismissed alarm was reasoned through, the chain produces a verifiable artifact in seconds. Per §1.2 epistemic-scope, the chain proves (a) what the AI said at a specific time and (b) that the record was not tampered with after capture; per §1.1 Daubert grounding, the institution's expert witness has a concrete answer to each of the four factors (testability under §7, peer review under the FFIEC working-group process, known error rate under §1.1's three-layer compromise model, general acceptance under the NIST-standardized primitives). Per §10.13, the chain output composes with the documentary artifacts (SDK version manifest, SDK source hash, HSM configuration, daily seal-job logs, change-management records, verifier output) that substantiate FRE 901(b)(9) authentication of the process; the institution's IT witness lays foundation from these at deposition without re-engineering the system.

> **⚠️ Surprise #8 — chain cannot extend evidentiary value upstream of AI**
> The chain cannot extend its evidentiary value upstream of the AI service. If a leak prediction were dismissed based on tampered sensor inputs, the chain would faithfully record the tampered inputs and the dispatcher's reasonable decision based on them. Detecting upstream tamper requires extending cryptographic integrity to the AMI head-end (Itron OpenWay 5.4) and to the historian (PI integrity extension). Both are on the Phase 2 roadmap. Both are gating items for the public-safety evidentiary story Pacific Crescent will need. Per §1.2 epistemic-scope, this is the (c) non-claim made explicit — the chain proves what the AI said, not whether what the AI said was true; upstream-input authenticity is a (c) question, and the chain is silent on it by design until the chain's first cryptographic operation moves left to the sensor itself.

---

## 🔧 5:00 PM — Verifier Distribution and CC8.1 Hygiene (Luis closes the loop)

Luis had spent the afternoon on logs, ops, and the verifier distribution posture — the §10.26 lens. He came to the whiteboard with three findings and one observation.

"§10.26 reference-verifier distribution. Pacific Crescent operates the reference verifier off USB media for in-person examiner sessions and from a CIP-categorized internal artifact registry for production verifications. The reference verifier ships in a separate repository from the spec under Apache 2.0, which is normative per §10.26. Per-release artifact discipline is conformant — reproducible builds, Cosign signatures, per-platform binaries (Linux x86_64 + Linux ARM64 are minimum, Pacific Crescent additionally pulls Windows x86_64 for the WECC examiner laptop and macOS ARM64 for the CCO's laptop), SHA-256 + SHA-512 manifests, CycloneDX SBOM. The verification key fingerprint is named in CC8.1 and the institution's posture matches §10.26 CC8.1 citation discipline — implementation named (the reference verifier), version named (pinned per §11), verification key named."

Dawn wrote: *§10.26 verifier distribution conformant. Examiners run binaries off USB. Pacific Crescent matches the §10.26 CC8.1 three-name citation rule.*

"Spec-version pinning per §10.26 is also conformant — the institution cites the verifier version pinned in spec §11 References for v1.0b. Later verifier releases that maintain back-compat are acceptable; the pinned version is the floor."

Luis's second finding was on §10.18 CC8.1 cross-referencing.

"The runbooks support the chain-of-custody program but the cross-referencing is uneven. The seal-job runbook names §4.2 / §4.3 by section number — clean. The HSM partition-ceremony runbook names §10.5 and §10.17 — clean. The IKM-rotation runbook names §10.10 — clean. The dispatcher-application runbook does NOT cite a spec section, which per §10.18 is a Nit — the runbook describes the chain integration without naming §4.4 / §4.4.1 / §4.4.2 / §4.4.6, so a SOC-engagement reviewer reading the runbook needs to map across by content. Five-line addition closes it."

Dawn wrote: *§10.18 CC8.1 cross-referencing — Nit on dispatcher-application runbook. Five-line fix.*

Luis's third finding was on §10.17.

"The HSM partition-ceremony attestation under §10.17 is being emitted — `chain.partition_ceremony_attended` operational events for partition creation and IKM rotation, with `ceremony_type`, `partition_handle`, `ceremony_started_at_utc`, `ceremony_completed_at_utc`, `signatories` array, `witness`, `attendance_pdf_sha256`. The signatory `entity_affiliation` field per Round-17 M&A-P1 is populated. The `attendance_pdf_holder` field is also present (Pacific Crescent retains the original attendance log in its compliance vault). The `hsm_attestation_token_b64` field is RECOMMENDED at v1.0b and Pacific Crescent emits it — Thales Luna exposes ceremony-bound attestation tokens through its attestation API, and the verification path is named in CC8.1. Conformant."

Dawn wrote: *§10.17 partition-ceremony attestation conformant — `chain.partition_ceremony_attended` events emitted with all required fields including the RECOMMENDED `hsm_attestation_token_b64`. Pacific Crescent ahead of the v1.0b posture.*

Luis's observation was on §10.4 trusted-time.

"§10.14 trusted-time integration is RECOMMENDED but NOT REQUIRED for v1.0 conformance. Pacific Crescent is interesting because they have PMU-grade time synchronization on the BES side — IRIG-B from a GPS-disciplined master clock — that they could use as a NTP discipline upgrade. Per §10.4, NTP-synchronized application hosts and ledger servers are SHOULD; Pacific Crescent's existing PMU-grade time source exceeds the SHOULD bar. Their CC8.1 already names the time source. The §10.14 informative paragraph names RFC 3161 trusted-timestamp integration as a candidate v1.x extension; until that lands, NTP discipline (or PMU-grade equivalent) is the v1.0 timestamp foundation."

Dawn wrote: *§10.4 NTP synchronization conformant via PMU-grade GPS-disciplined master clock — exceeds SHOULD bar. §10.14 trusted-time integration RECOMMENDED, candidate v1.x extension. Note for the CFO: the time-discipline investment pays off on the chain side without a separate integration.*

> **✓ Confirmation #7 — verifier distribution and CC8.1 hygiene**
> The §10.26 reference-verifier distribution posture is conformant — repository separation, per-release reproducible builds, Cosign signatures, per-platform binaries, SHA-256 / SHA-512 manifests, CycloneDX SBOM, spec-version pinning per §11, three-name CC8.1 citation discipline (implementation, version, verification key). The §10.17 HSM partition-ceremony attestation is conformant including the v1.0b RECOMMENDED `hsm_attestation_token_b64` field. The §10.4 NTP discipline is exceeded by PMU-grade time synchronization. One §10.18 CC8.1 cross-referencing Nit on the dispatcher-application runbook (does not cite §4.4 / §4.4.1 / §4.4.2 / §4.4.6 by section number); five-line fix.

---

## 🌆 5:30 PM — Auditor Debrief

The team reconvened in the clean room. Esme stayed. Coffee was cold. Mount Hood had moved into the late-afternoon haze.

Dawn stood at the whiteboard. Three columns.

| Tier | Status |
|---|---|
| AI side (`pipeline-leak-detection`, `pipeline-integrity-trending`) | 0 Gaps, 0 Partials, 1 Nit (§10.18 dispatcher-runbook cross-reference) |
| OT side (PI historian, GE iFIX SCADA, Itron OpenWay AMI, OMS work-order linkage) | 5 Gaps, 6 Partials |
| Customer-billing side (Salesforce, CIS, OMS customer interactions) | 3 Gaps, 4 Partials |

"That's the shape. Three tiers. One passes (with one Nit). Two don't. Today is Stelvio with public-safety consequences."

Esme stood with her arms crossed, listening.

"AI side." Dawn pointed. "TesseraSeal. Seven confirmations — chain integrity verified across nine months and 1.6 million inferences (§7 12-step procedure, §10.12 exit code contract), append-only ledger behavior under direct DB mutation attempt (§10.3), credential rotation under chain (eight rotations sampled, all PASS, §10.10 rotation discipline), live alert observed end-to-end on the operations floor with verifier PASS in four seconds (§4.4 / §4.4.2 / §10.11-style parent-linkage), four-of-four reconciliations PASS on the AI side, the public-safety evidence question answerable in two paragraphs (§1.2 / §1.1 / §10.13), and verifier-distribution + CC8.1 hygiene conformant (§10.26 / §10.18 / §10.17 / §10.4). The AI side maps cleanly to NERC CIP-007 for the AI/ML control points and to PHMSA pipeline-integrity expectations for the leak-detection scope. IEC 62443-3-3 SR 6.1, SR 6.2, and SR 7.5 are explicitly satisfied for this service. The single Nit on §10.18 cross-referencing is a five-line fix on the dispatcher-application runbook."

She moved to the OT column.

"OT side. Five gaps. PI historian retention is 60 days with three engineers holding override authority — the upstream half of the forensic trail goes invisible past 60 days; per §10.13 evidentiary-artifacts retention this is well below the chain-data retention floor. GE iFIX backing store captures INSERTs only on operator alarm-acknowledgment notes; subsequent UPDATEs are not recorded — same shape as the §10.3 application + role enforcement requirement, which iFIX's product behavior violates. Itron OpenWay AMI is on 5.2, version-locked behind the 5.4 integrity-checking option — the chain's first cryptographic op cannot move leftward until 5.4 is live. The HMI workstations use a shared `Operator_ControlRoom` account across three shifts with no MFA. The dispatcher_id captured in chain entries is the AD display name, not federated SSO with SID-binding — operational not integrity. Six partials around OMS work-order content-coupling (§10.19 `audit.external_artifact.*` family closes it as a Phase 3 work item — six attributes, no spec change), paper-to-OMS transcription (Phase 3), AMI override authority discipline (Phase 2), sensor-to-AI path authentication (Phase 2 — `audit.connector_source.*` per §4.4.6 by analogy from PI ingestion adapter), IEC 62443-3-3 SR 7.5 coverage on legacy systems, and historian disk-volume sizing (Phase 2)."

"Customer-billing side. Three gaps. Salesforce overwrite shape on customer-interaction notes. OMS work-order linkage from field finding to AI prediction is by ID only, not by content hash (Phase 3 closure via §10.19 `audit.external_artifact.*`). CIS audit-trail retention on interaction notes is shorter than on billing records. Four partials around backup-vs-change-history, retention variance across three states' regulators, customer-CRM IAM federation, and a future §10.16 SaaS-edge mirror connector posture should Pacific Crescent ever stand a Salesforce mirror for the leak-detection AI to ingest customer-call data — not load-bearing today, but pre-named in the chain-coverage map so the four-number lag posture (median, 95th-percentile SLO, alerting threshold, RTO) is the entry-fee for that connector if it lands. Phase 4 territory. The CIO's budget conversation."

Dawn put the pen down.

"Three observations to close."

"One. The chain works. Pacific Crescent has been running it on the leak-detection service for nine months. Today we observed a real alert dispatched to a real crew and confirmed as a real leak — chain captured every step that mattered, verifier returned in four seconds. The investment paid off. When the PUC asks, you have an answer. Per §1.1 Daubert grounding, your IT witness has a one-page answer to each of the four factors; per §1.2 the chain's claims are bounded honestly; per §10.13 the documentary artifacts that compose with the chain output substantiate FRE 901(b)(9) authentication of the process."

"Two. The chain is on the part that decides whether to dispatch. The chain is not on the part that produces what the decision is made about. Stelvio's seam was the same shape — AI sealed, OT mutable. Pacific Crescent's seam is in the same place but the consequence is different. A wrong reading and a dismissed alarm at Stelvio is a yield-loss claim on a heat of steel. A wrong reading and a dismissed alarm at Pacific Crescent is a public-safety incident. The seam matters more here. Closing it through Phase 2 is not a compliance project. It is a public-safety project that happens to also close a compliance gap. Per §10.19 chain-coverage map, the Phase 2 + Phase 3 closures redraw the map — PI moves from 'not yet chain-instrumented' to 'chain-instrumented' once integrity-extension is live; the AMI head-end follows; OMS work-order content-coupling lands via `audit.external_artifact.*`; the upstream half of the forensic trail extends from 60 days to the institution's full chain-data retention period."

"Three. The roadmap you already have is the right roadmap. Phase 2 in 12 months — AMI head-end to 5.4, PI integrity extension, PI retention extended via disk re-sizing, MFA on the HMI workstations, federated SSO with SID-binding for the dispatcher application, `audit.connector_source.*` emission from PI ingestion. Phase 3 in 18 months — OMS work-order content-coupling via §10.19 `audit.external_artifact.*`, paper-to-OMS gap closure, IEC 62443-3-3 SR 7.5 coverage on legacy. Phase 4 deferred — customer-CRM, including §10.16 SaaS-edge connector posture if a Salesforce mirror is stood. The order is correct because Phase 2 closes the public-safety story first. Document it that way for the CFO."

Esme nodded. "What do I take to the CEO Friday?"

Tom answered. "The three-tier summary. The §10.19 chain-coverage map costed against the three-tier summary. The four-second verifier output Mike captured during the live Brentwood alert as a real artifact showing what the prior investment delivered — name the §7 step exit code 0 per §10.12 in the artifact caption so the CEO sees the conformance citation. The reconciliation table from this afternoon — four predictions, what worked, what didn't, where the legacy erosion is. And the public-safety paragraph from the 4:30 conversation — the chain proves the AI's behavior per §1.2(a)(b), the chain does not prove the upstream's behavior per §1.2(c), Phase 2 closes the gap, the gap is fundable today."

"Send me the Brentwood verifier capture."

Mike held up his phone. "Got it on the laptop. I'll send it tonight."

"And the report?"

"Thursday morning. Before your Friday CEO review."

Dawn closed her notebook. "We'll send it."

Esme's shoulders dropped that half-inch again. She didn't smile but she nodded twice.

The team packed up. Raj and Luis loaded boxes of evidence into the rental SUV. Diana and Elena said goodbye to Esme at the badge desk. Mike and Chen took one last look at the operations floor on the way past — the gas-distribution map glowing green except for two amber dots, both being watched, both at confidence below dispatch threshold.

Dawn walked out last. She paused at the badge desk and looked back through the glass at the operations floor. Yolanda was still on shift. The methane sensor readings were still updating every fifteen seconds. The Brentwood crew had finished isolation and were filing the work-order paperwork on a tablet.

> **🔍 Dawn's note (internal):**
> *It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach.*
>
> *Today the chain reached far enough to dispatch a real crew to a real leak in real time and prove it after the fact. Today the chain did not reach the sensor. Phase 2 closes the gap before the gap closes a neighborhood. The §1.2 epistemic-scope language is the one paragraph that tells the CEO what the chain delivers and what it doesn't. The §10.19 chain-coverage map is the one artifact the WECC auditor and the PHMSA inspector and the WUTC regulator read as the same picture. Two pieces of writing carry the whole report.*

---

## ✅ vs ❌ — The Three-Tier Summary

### ✅ AI Side (TesseraSeal — `pipeline-leak-detection`, `pipeline-integrity-trending`)

| Item | Status |
|---|---|
| Chain integrity (HMAC + Merkle + daily Ed25519 seal on Thales Luna PCIe HSM, FIPS 140-2 Level 3, CIP-categorized network) per §4.1 / §4.2 / §4.3 / §10.5 | PASS |
| §7 verifier procedure 12-step pass; §10.12 exit code 0 | PASS |
| §1.4 compositional security (three independent authentication layers) | PASS |
| Append-only ledger behavior under direct DB mutation attempt per §10.3 (application + role layers) | PASS — verifier catches at HMAC layer (§7 step 9) |
| Multi-entry tamper attempt | PASS — verifier catches at Merkle/seal layer (§7 step 10) |
| §1.1 Daubert four-factor grounding (testability / peer review / known error rate / general acceptance) | PASS |
| Credential rotation under chain per §10.10 boundary discipline | PASS — eight rotations sampled across nine months, all PASS |
| IKM length and generation per §10.6 / §10.6.1 (Thales Luna internal CSPRNG, RNG type recorded in `master_key.generated`) | PASS |
| §10.1 key-fingerprint reconciliation (daily cadence, exceeds RECOMMENDED weekly) | PASS |
| §10.9 IKM-registry retention (longer of 7 years or chain-data retention) | PASS |
| Inference -> chain latency | ~200 ms observed live on operations floor |
| Verifier latency | ~4 seconds for any single entry, 12 steps, exit code 0 |
| Live alert observed end-to-end (Brentwood, 2026-05-05) | PASS — prediction, dispatch, crew on site, confirmed real leak, chain captured each step |
| Reconciliation test (4 predictions, AI-side) | 4/4 PASS |
| Dispatcher decision capture (identity, reason code, model inputs) per §4.4 attribute schema | PASS — every dispatcher action sealed in chain |
| `gen_ai.request.model` + `gen_ai.response.model` per §4.4 normative requirement and §7 step 12a check | PASS |
| Field-crew follow-up entry mechanism with §4.4 `parent_run_id` / `parent_seq` linkage | PASS — OMS closes work order, leak-detection service writes follow-up referencing original prediction |
| §10.17 HSM partition-ceremony attestation (`chain.partition_ceremony_attended`, full schema including `entity_affiliation`, `attendance_pdf_sha256`, `attendance_pdf_holder`, RECOMMENDED `hsm_attestation_token_b64`) | PASS |
| §10.26 reference-verifier distribution (repo separation, reproducible builds, Cosign sigs, per-platform binaries, SHA-256/512 manifests, CycloneDX SBOM, three-name CC8.1 citation) | PASS |
| §10.4 NTP discipline via PMU-grade GPS-disciplined master clock | EXCEEDS SHOULD |
| §4.3 v1.0b 12-line `sign_payload` form (binds `key_versions_canon` + `hex(kms_handle_uris_digest)` per Round-17 NIST-G1/G2) | PASS |
| §10.7 software-key-adapter exclusion in production (compile-time exclusion documented in CC8.1) | PASS |
| §10.8 constant-time comparison for fingerprint and MAC checks | PASS |
| §4.4.2 deployment-intent (`production` for the leak-detection service; CC8.1 names the single-region single-version production posture) | PASS |
| §10.18 CC8.1 cross-referencing for chain-of-custody runbooks | NIT — dispatcher-application runbook does not cite §4.4 / §4.4.1 / §4.4.2 / §4.4.6; five-line fix |
| IEC 62443-3-3 SR 6.1, SR 6.2, SR 7.5 mapping for AI scope | PASS — explicitly satisfied |
| NERC CIP-007 readiness for AI/ML control points | Demonstrable |
| PHMSA pipeline-integrity evidence for leak-detection scope | Demonstrable |
| Public-safety evidence question (dismissed-alarm scenario) per §1.2 / §1.1 / §10.13 | Answerable in two paragraphs with chain artifact |

### ❌ OT Side (PI historian, GE iFIX SCADA, Itron OpenWay AMI, OMS work-order linkage)

| Item | Status |
|---|---|
| PI historian audit-trail retention | 60 DAYS — past 60 days, who-changed-what is unrecoverable; below §10.13 evidentiary-artifacts retention floor and below NERC CIP-008-6 / CIP-009-6 3-year minimum |
| PI override authority | THREE ENGINEERS with `PIWorld\db_admin`, no second control |
| Sensor-to-AI path authentication | NONE — first cryptographic op (§4.1) is at the AI service; three hops upstream are trust-by-policy; per §1.2 epistemic-scope (c) non-claim |
| `audit.connector_source.*` family on PI-to-AI ingestion adapter (§4.4.6 by analogy) | NOT EMITTED — Phase 2 work item, six attributes |
| GE iFIX `iFIXAlarmAck` UPDATE audit | NONE — INSERT captured, subsequent UPDATEs silently overwrite; product violates §10.3 application + role enforcement |
| HMI workstation account | SHARED — `Operator_ControlRoom`, six users, no MFA, smart-card readers Phase 2; NERC CIP-007 unique-user accountability finding standalone |
| Dispatcher_id binding in chain | DISPLAY NAME, not federated SSO with durable SID; control-completeness (chain integrity unaffected per §4.1) |
| Itron OpenWay version | 5.2 — version-locked behind 5.4 integrity-checking option, MDM compatibility constraint; Phase 2 closure moves §1.2 trust boundary leftward |
| AMI override authority (`meter_data_engineer`) | TWO accounts, free-text reason code, override captured in head-end log but no chain-coupling to AI-side chain |
| OMS work-order linkage to chain | ID-LINKED, not content-hashed — §10.19 `audit.external_artifact.*` family closes it as a Phase 3 work item (six attributes, no spec change) |
| Paper-to-OMS transcription | ~40% OF WORK ORDERS — dispatcher transcribes paper, chain references work-order ID, content-coupling broken; same Phase 3 closure path |
| Reconciliation test sensor-side trace | 2/4 PASS, 2/4 UNPROVABLE due to PI 60-day retention |
| §10.19 chain-coverage map (PI / iFIX / AMI / OMS in "institutional systems not yet chain-instrumented" with rollout posture and substitute named) | DOCUMENTED |
| IEC 62443-3-3 SR 7.5 on legacy systems | NOT MET |
| Phase 2 remediation timeline | 12 months — AMI 5.4 upgrade, PI integrity extension, PI disk re-sizing, MFA on HMI, federated SSO, `audit.connector_source.*` emission |
| Phase 3 remediation timeline | 18 months — OMS content-coupling via §10.19 `audit.external_artifact.*`, paper-to-OMS closure, IEC 62443-3-3 SR 7.5 on legacy |

### ❌ Customer-Billing Side (Salesforce, CIS, OMS customer interactions)

| Item | Status |
|---|---|
| Salesforce field-level audit | PARTIAL — enabled on some fields, disabled on customer-interaction notes |
| Salesforce retention shape | "Backups, not version history" — same as diary baseline |
| CIS audit-trail retention | 18 MONTHS on billing records, shorter on interaction notes |
| OMS customer-interaction linkage | NONE between customer-service interactions and AI predictions even when the interaction triggered the alert |
| Multi-state PUC retention variance (WA, OR, CA) | THREE different retention shapes, not harmonized |
| Customer-CRM IAM federation | PARTIAL — Salesforce SSO is federated, CIS internal user model is not |
| ERP/CIS billing audit clearable by admin role | YES — same as diary baseline |
| Future §10.16 SaaS-edge mirror connector posture (if Pacific Crescent ever stands a Salesforce mirror for AI ingest) | NOT YET LOAD-BEARING — pre-named in §10.19 map; four-number lag posture (median, 95th-percentile SLO, alerting threshold, RTO) is the entry-fee per §10.16; imprecise lag wording will be a §10.16 non-conformance, not a Nit |
| Phase 4 (customer-CRM) | Deferred pending CIO budget. Document explicitly. |

---

## 🧾 Final Assessment Theme

> *"The chain is on the part that decides whether to dispatch. The chain is not on the part that produces what the decision is made about. Pacific Crescent's customers don't know that line. Their regulators don't yet either. A dismissed alarm and an exploded house would teach the line in the worst possible way. Phase 2 teaches it the right way, in 12 months, on a fundable budget. Per §1.2 epistemic-scope, the line is in the spec — the chain proves what the AI said and that the record was not tampered with after capture; the chain does not prove the AI's input was what the sensor actually measured. Phase 2 moves §1.2's first cryptographic operation leftward to the sensor."*

Pacific Crescent Power & Gas demonstrates AI-decision integrity within scope. The investment in TesseraSeal on the leak-detection service returns verifiable provenance under NERC CIP-007 for the AI/ML control points, PHMSA pipeline-integrity expectations for the leak-detection scope, and IEC 62443-3-3 SR 6.1, 6.2, and 7.5. The §7 verifier procedure runs in four seconds on the operations clean-room network with the §10.12 exit code 0 contract. A live alert observed end-to-end during the engagement — prediction, dispatch, crew on site, confirmed real leak — left a sealed forensic trail under §4.3 v1.0b 12-line `sign_payload` that will satisfy a NERC auditor, a PHMSA investigator, and an insurance carrier. The chain composes with the §10.13 evidentiary-artifacts retention list (SDK manifest, source-code hash with build reproducibility, HSM configuration, daily seal-job logs, change-management records, verifier output) so the institution's IT witness can lay foundation under FRE 901(b)(9) at deposition without re-engineering the system. Per §1.1 Daubert four-factor grounding, the institution has a one-page answer to each factor — testability under §7's ordered byte-exact verification procedure with the test-vector corpus pinning byte values; peer review under the FFIEC working-group process and the public Apache-2.0 reference implementation; known error rate under §1.1's three-layer compromise model (IKM + ledger + HSM, each compromised by different roles under separation-of-duties); general acceptance under the NIST-standardized primitives.

Outside the AI scope, integrity is operational discipline. The OT side carries the same mutability shape as the manufacturing OT engagement the team documented two weeks ago — different vendors, identical boundary problem. The PI historian's 60-day audit retention erodes the upstream half of the forensic trail past 60 days and falls below the §10.13 chain-data retention floor and the NERC CIP-008-6 / CIP-009-6 3-year minimum. The HMI's INSERT-only audit semantics on operator notes (which violate §10.3 application + role append-only enforcement at the product level) erode the operator-side accountability. The shared `Operator_ControlRoom` account compounds the §10.3-shape gap with a NERC CIP-007 unique-user-accountability finding standalone. The Itron OpenWay 5.2 version-lock keeps the cryptographic boundary at the AI service rather than at the sensor; per §1.2 epistemic-scope, the chain's claims about what the AI said are bounded honestly while the (c) non-claim about the AI's input authenticity is exactly what Phase 2 closes. The OMS work-order content-coupling gap leaves a paper-transcription seam that grows with time-since-incident; per §10.19 `audit.external_artifact.*` family, six attributes on the follow-up chain entry close it without a spec change. None of these surprised Esme. All of them are on the roadmap she has been building for a year. The §10.19 chain-coverage map names them as "institutional systems not yet chain-instrumented" with the rollout posture and substitute named — explicit, testable, accountable to a vendor-management auditor.

The customer-billing side is a different audit — PUC scope rather than NERC scope, with the same overwrite-and-backup shape that has appeared in five of the team's last six engagements. Phase 4 is the customer-CRM conversation. If Pacific Crescent ever stands a Salesforce mirror for the leak-detection AI to ingest customer-call data, §10.16 SaaS-edge connector posture lands as the entry-fee — four quantified numbers (median, 95th-percentile SLO, alerting threshold, RTO) in the runbook, `connector.lag_observation` and `connector.outage` operational events under §10.2, and the §10.16 severity-classification clause that imprecise lag wording is non-conformance not Nit. Pre-named in the §10.19 chain-coverage map so the future Pacific Crescent CC8.1 author does not re-derive it.

The framing for the CEO and the board is straightforward. Phase 2 is not a compliance project. Phase 2 is a public-safety project that happens to also close a compliance gap. The chain proves the AI's behavior. The chain does not prove the upstream's behavior. Closing the upstream gap before the gap closes a neighborhood is the operationally correct sequence and it is the morally correct sequence and it is fundable today.

Three tiers, one report, one severity scale, three remediation timelines, one public-safety paragraph that decides everything else. Pacific Crescent knows where the seam is. Their CEO is about to. Their PUC commissioners need to. The §1.2 epistemic-scope language is the controlled vocabulary that lets every reader land on the same understanding without re-deriving it. The §10.19 chain-coverage map is the artifact every regulator reads as the same picture.

The report's reading order follows §13 stakeholder navigation — the executive summary (CEO, audit committee, board) lands first; the chain-operations and DevOps detail (the SOC team, internal audit, the FFIEC IT examiner) lands second; the cryptographic-witness foundation (§1.1 / §1.2 / §1.3 / §1.4 / §10.13) lands third for the litigation-support team. The §11 References section pins the verifier version Pacific Crescent's CC8.1 cites; later patch releases that maintain back-compat are acceptable, but the pinned version is the conformance floor. Every section number cited in the report points at the v1.0b spec text; an examiner with the spec in one hand and the report in the other can walk every determination back to the normative spec language without inferring or re-deriving any of it.

---

*End of diary. Filed Tuesday evening. Report drafted Wednesday and Thursday. Delivered Thursday morning before the Friday CEO review.*
