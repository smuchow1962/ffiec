# 🧾 Diary of an Audit Day — Pacific Crescent Power & Gas

**Engagement:** NERC CIP audit readiness assessment with PHMSA pipeline integrity-management overlay
**Client:** Pacific Crescent Power & Gas — investor-owned utility, Pacific Northwest, ~3.1M customer accounts (electric + gas), multi-state (WA, OR, parts of northern CA and ID), HQ in Portland
**Posture:** Partial TesseraSeal deployment — chained on a single AI use case (gas-pipeline leak detection); not chained on SCADA, the PI historian, the AMI head-end, the OMS, the CIS, or the customer CRM
**Date:** Tuesday, the week after Helmstad Bio
**Auditor:** the same eight-person team — diary baseline, Mercator, Northbridge, Stelvio, Atrio, Helmstad

---

## Context

Pacific Crescent runs hydro from the Columbia, a fleet of natural-gas combined-cycle plants down the I-5 corridor, a recently added wind portfolio in eastern Washington, and a gas-distribution system that snakes through Portland, Vancouver, Salem, and parts of Eugene. Three regulators in three states. NERC CIP-002 through CIP-014 on the bulk-electric-system side. PHMSA pipeline integrity on the gas side. WUTC, OPUC, and a small slice of CA-CPUC on the customer side. DOE/CESER cybersecurity expectations on top of all of it. IEC 62443 on the industrial-control side because two of their large industrial customers have started asking for it in their procurement language.

Fourteen months ago, a near-miss gas leak in a Portland neighborhood turned political. The utility had detected residual methane on a routine patrol, but the detection happened roughly three hours after a homeowner had already smelled it and called the utility's customer line. The PUC opened an inquiry. The mayor's office called twice. The utility's CEO accelerated a procurement that had been quietly in progress: an AI-driven gas-pipeline leak-detection system that fuses methane sensor data, weather, soil-saturation models, and historic incident data to predict pipeline-failure probability in real time.

Nine months ago, that system went live. Pacific Crescent stood up TesseraSeal at the same time, chained on the leak-detection service. Three `service.name` values inside one tenant — `pipeline-leak-detection` (in production), `pipeline-integrity-trending` (in production, lower-priority analytics), and `outage-prediction-pilot` (pre-production, not yet in scope today).

Everything else at Pacific Crescent runs the way utilities have always run.

The OT side — GE iFIX SCADA on Windows Server 2016, OSIsoft PI historian, Itron OpenWay AMI head-end on version 5.2 — is not internet-connected. NERC CIP requires that. But "air-gapped from the internet" is not the same as "integrity-controlled," and that's a distinction the regulators are starting to write into their guidance. The customer-side stack — a customized Oracle CIS, an outage-management system, Salesforce CRM — runs on the IT business network. Different mutability shape. Same mutability problem.

Esme Yamashita, Pacific Crescent's Chief Compliance Officer, came over from BPA two years ago. She spent two years on a NERC standards drafting team for CIP-013. She knows what a NERC auditor will look for, and she knows what they will quietly skip. She doesn't oversell. On the prep call she told Karen: "We need a map. NERC CIP-007 is the area where the regulators are going to push hardest in the next two years. They're starting to ask AI questions. Tell me what we already have and where the gaps will land."

This is the diary of that day.

---

## Audit Team

- **Karen** — Lead Auditor (governance and narrative)
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

Karen rode in with Tom from the airport hotel. Twenty minutes north on I-5, then off at the river. The Pacific Crescent control center sat on a low rise overlooking the Willamette, glass and concrete, fenced with the kind of fence that meant something. NERC-CIP-classified facility. They saw the badge readers from the parking lot.

Tom drove. Karen drank coffee.

"Five engagements in five weeks." Tom merged across two lanes. "What's the through-line at this point?"

Karen watched the river. "Northbridge was the gold standard. Banking. TesseraSeal in everything. We ran out of things to find by 3 PM."

"Mercator."

"Half the river sealed, half not. Healthcare. Imaging side chained, claims side mutable. We wrote the seam down the middle of the report and the CMO understood it instantly."

"Stelvio."

"Manufacturing. Three zones. AI sealed. OT mill-floor mutable. IT business mutable. Maria knew where every gap was and wanted the report so she could take it to the CFO Friday. We wrote it that way."

"Atrio."

"Banking-as-a-service. Forty-seven tenants, all clean. The cleanest multi-tenant verifier work we've ever seen. We were on the plane home by 4."

"Helmstad."

"Biopharma. ALCOA+ on the AI side, CRO data legacy. The seam was in a different place again. The CRO data showed up in the chain by reference but the source-of-truth lived in someone else's system entirely."

Tom checked the GPS. Five minutes out.

"And today?"

Karen put her cup down. "Pacific Crescent has 3.1 million customers and a fleet of pipelines. The legacy gaps here are not just compliance findings."

"What are they."

"Stelvio's worst case was a yield-loss claim on a heat of steel. Mercator's worst case was a misadjudicated claim. Atrio's worst case was a sponsor-bank reconciliation finding. Helmstad's worst case was a CRO data-trace gap on a Phase 2 trial."

Tom waited.

"Pacific Crescent's worst case is a wrong reading and a dismissed alarm and someone's house blows up. Stelvio with public-safety consequences."

Tom didn't say anything for a while.

"It never is."

"It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach. We're going to find both today."

They pulled into the visitor lot at 8:24.

Esme met them at the badge desk. Mid-forties, dark blazer, the kind of badge clipped to her lapel that had three different access levels printed on it. She shook Karen's hand once, firmly, and Tom's the same way.

"You'll go through three checkpoints to get to the operations floor. The control center itself is CIP-categorized. Phones in the locker. Laptops registered at the second checkpoint. We have a clean room beside the control floor where you can work — same network segment as corporate but with read-only feeds from the OT side. The HSM and the AI ledger are reachable from the clean room. SCADA and the historian are not — you'll see those over a screen-share from one of my engineers."

"Understood."

"I want to be straight with you about scope before we start." Esme walked them toward the first checkpoint as she spoke. "The leak-detection AI is in scope under NERC CIP-007 because it touches a BES Cyber Asset adjacent to gas-electric tie-points, and it's in scope under PHMSA because it's a pipeline-safety system. The PI historian is in scope under CIP. The AMI head-end is in scope. The OMS is partly in scope — work orders that touch BES assets are in. The CIS and the Salesforce CRM are out of scope under NERC but very much in scope for the PUC. We're going to look at all of it because that's how I think about the program. The report can split scope however you need it to."

Karen nodded. "That's how we'd write it anyway."

"Good. Let's start in the operations control room."

The team kitted up. Phones in lockers. Laptops registered. Three badge taps later, they were on the operations floor.

> **🔍 Karen's note (internal):**
> *It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach.*
>
> *Calibrate. The AI side will pass. The OT side won't. The customer-billing side won't. But the consequence axis is different from Stelvio's. Stelvio's worst case was money. Pacific Crescent's worst case is people. Document accordingly.*

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

Karen turned to Esme. "Action means dispatch?"

"Action means decision. Marcus gets four buttons — dispatch crew, monitor, dismiss, escalate-to-supervisor. Whatever he picks, that decision lands in the chain along with the model's prediction."

Mike opened his laptop, balanced it on the edge of Marcus's station, and pulled up Herald.Compliance. He filtered to `service.name = pipeline-leak-detection` and the last five minutes.

The Powell-44 entry was sitting at the top of the list. Confidence 0.31. No dispatcher action yet. The chain entry referenced a methane sensor ID, a weather snapshot hash, a soil-saturation model output hash, and a prediction.

Marcus glanced at the screen. The reading dropped to 0.9 ppm above background. Then 0.6. Then back to baseline.

"There. It's washing out."

The dashboard moved Powell-44 to the archive. Marcus tapped *dismiss* on the queue, picked a reason code from a dropdown — `wind shift, baseline restored` — and submitted.

A new chain entry hit Mike's screen.

"There it is." Mike tilted the laptop. "Decision entry. References the prediction entry. Dispatcher ID, reason code, timestamp."

Karen leaned in. "Run the verifier on it."

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

"Twelve steps, four seconds, on the operations clean-room network." Mike snapped the laptop shut. "That's the thing working. The dismissal is sealed. The model's inputs are sealed. The dispatcher's reason code is sealed. If anyone ever asks why that low-confidence flag was dismissed, you have an answer that's the same answer six months from now and six years from now."

Esme nodded once. "That's what it was procured to do."

> **✓ Confirmation #1**
> The leak-detection chain is live on the operations floor and producing verifiable entries within ~200 ms of each prediction or dispatcher decision. Mike re-verified a real-time dismissal in four seconds standing at a working dispatcher station. Reason codes are captured. Inputs are hashed. The infrastructure is real and observable on the production line.

Karen wrote in her notebook: *AI dashboard live. Dispatcher decisions sealed with reason codes. Inputs hashed. Verifier runs against the operations clean-room network in four seconds.*

She walked over to the wall screen — the gas-distribution map. Twelve hundred miles of distribution pipe. Several thousand methane sensors. The dashboard rendered every active prediction as a small dot. Three dots glowed amber. Most of the map was green.

"How many predictions a day?"

Esme answered. "On a normal day, four to six thousand inferences. About a hundred and fifty cross the 0.4 threshold and require dispatcher action. Maybe three to eight result in actual dispatch."

"And the model has been live for nine months?"

"Live for nine months. About 1.6 million inferences in the chain so far. Two hundred and forty actual dispatches. Eighty-something true positives — small leaks, mostly. Two of those would have escalated if we hadn't dispatched."

"And what was the false-positive rate before the chain?"

"There was no chain before. We started chained from day one. The vendor required it for the contract."

Karen wrote: *Chained from day one. 1.6M inferences. ~240 dispatches. ~80 true positives. Two would have escalated. Public-safety value already demonstrable.*

Esme led them off the operations floor and back through two of the three checkpoints to the clean room. Coffee. Whiteboard. Network jacks.

Karen pulled up a chair.

"Let's split. Raj — historian and AI ledger. Diana — IAM, all three sides. Mike and Chen — pipelines and the AI service. Elena — CIS and Salesforce. Luis — logs and ops. Tom — sit with Esme, work the NERC CIP and PHMSA mappings. Reconvene at noon."

They split.

---

## 🧠 10:00 AM — Database Deep Dive (AI Ledger and PI Historian)

Raj had a corner of the clean room and two screens. One showed the AI-side ledger — backed by a Postgres instance reachable from the clean-room network. The other showed a screen-share from one of Esme's engineers — the OSIsoft PI historian's PI Server admin console. Raj worked them in parallel.

### The AI ledger

Raj started with the chain. Append-only by design. The Herald.Py SDK signs each entry. HMAC-SHA-256 chain links entry N to entry N-1 with HKDF-per-tenant key binding. Daily Ed25519 seals close out each day's chain on the on-prem Thales Luna PCIe HSM in the operations control center — the same NERC-CIP-classified facility, on a CIP-categorized network segment, air-gapped from corporate IT.

He picked a random entry from three weeks ago. Copied its ID. Ran the verifier.

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key leak-prod-2026-q1
```

He picked a random entry from yesterday. Same result. He picked the very first entry from nine months ago — the day the leak-detection model went live. Same result.

He tried to mutate one. Issued an UPDATE on the chain table directly through psql. The database accepted it because nothing at the database layer prevents it. Raj ran the verifier on the next entry in the chain.

```
Status: FAIL
Step: 4
Reason: HMAC mismatch — entry payload does not produce
        the chained HMAC recorded in entry N+1
```

He tried to also rewrite entry N+1 to match. The verifier failed at the daily seal:

```
Status: FAIL
Step: 9
Reason: Merkle root mismatch — recomputed root does not
        match sealed root for date 2026-04-14
```

He rolled back his mutations. The chain returned to PASS. He noted the test in his workbook.

> **✓ Confirmation #2**
> The AI-side ledger is append-only in practice. Direct database mutation is technically possible, but the verifier catches it at the HMAC layer (single-entry tamper) or the Merkle/seal layer (multi-entry tamper). The seal HSM lives on a CIP-categorized network segment that the database engineers do not have access to. To forge undetectably, an attacker would need both database write access and HSM signing access, and Pacific Crescent has those split across two different network zones.

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

> **⚠️ Surprise #1**
> The OSIsoft PI historian's audit trail is enabled but with 60-day retention. After 60 days, sensor-value changes are unrecoverable from the historian. Three engineers hold `PIWorld\db_admin` with override authority. A methane sensor reading from 90 days ago could have been edited at any time after the original write, and there is no way to determine whether that occurred. The AI service references sensor data by ID and timestamp, not by content hash, so the chain captures what the AI saw — not what the sensor actually measured.

Raj wrote in his workbook: *Same shape as Stelvio's Proficy historian. Different vendor, same mutability surface, longer retention than Stelvio's row-level (which had none) but still finite. The AI trusts whatever PI gives it. Document the boundary.*

He moved on to the GE iFIX SCADA backing store.

### The iFIX SCADA backing store

Aaron pulled up the SQL Server instance behind GE iFIX. The HMI itself runs on the operator workstations, but operator notes — the free-text notes an operator types when responding to an alarm — land in a SQL Server table called `iFIXAlarmAck`.

Raj asked for the schema. The table had `alarm_id`, `ack_user`, `ack_timestamp`, `ack_note`, and `last_modified`.

"Last modified. Does that mean what I think it means?"

Aaron clicked into the table definition. "There's a trigger that updates `last_modified` on UPDATE. There is no trigger that captures what the prior value was. So if I edit my ack_note on an alarm from last Tuesday, `last_modified` will tick to today, but the original note is gone."

Raj wrote: *iFIX. Operator ack notes mutable. Trigger captures last_modified, not prior value. UPDATE silently overwrites.*

"And there's an audit log somewhere?"

"GE iFIX writes an INSERT-only audit log to a separate table for the initial ack. UPDATE is not captured anywhere."

> **⚠️ Surprise #2**
> GE iFIX SCADA stores operator alarm-acknowledgment notes in a SQL Server backing table. The audit log captures the original INSERT. Subsequent UPDATEs to the ack note are not captured anywhere. An operator can revise what they said about an alarm response after the fact, and there is no record that the revision occurred.

Raj closed the screen-share. He had enough for the morning.

---

## 🔐 11:00 AM — IAM Review (AI Side, OT Side, Customer Side)

Diana had a workbook that walked through identity, access, and credential rotation — once for the AI side, once for the OT side, once for the customer side. Three columns. She filled them in the same order.

### AI side IAM

Every credential the leak-detection service uses — Postgres creds for the chain backing store, methane-sensor data feed creds, weather-API keys, the Luna HSM PIN for the daily seal, the dispatcher application's service account — every rotation was a chain entry. `event.type = credential.rotated`, with rotator identity, rotation reason, and the new key fingerprint.

Diana picked the last six rotations from the chain. Verified each. All PASS.

She asked Esme for the most recent rotation. Esme pulled it up.

```
herald-verify --tenant=pacific-crescent --service=pipeline-leak-detection \
              --event-type=credential.rotated \
              --date-range=2026-04-15:2026-05-05
```

Two entries returned. One Postgres cred and one weather-API key. Both PASS.

Diana also asked about the dispatcher application identity itself — the user_id that gets recorded when a dispatcher dismisses or actions a prediction.

Esme paused for a beat. "The dispatcher application authenticates the user against our Active Directory. The application records the user's display name as the dispatcher_id in the chain entry. So when the chain says `dispatcher: M.Reyes`, that's Marcus's display name from AD."

Diana wrote it down. "Display name, not federated SSO identity?"

"Display name."

"And display names can change."

"They can. AD doesn't enforce uniqueness on display names if a manager edits one. The underlying SID is unique. We chain the display name, not the SID."

Diana wrote: *Dispatcher chain captures display name, not federated SSO identity. SID is the durable identifier. The chain entry would survive a display-name change but the human-readable label would become ambiguous.*

> **✓ Confirmation #3**
> Credential rotation on the leak-detection AI service is captured in the chain with rotator identity, fingerprint, and reason. Two rotations sampled in the most recent window, both PASS. Eight rotations sampled across the past nine months, all PASS. The AI service account itself is not shared, has MFA, and rotates every 90 days under chain.

> **⚠️ Surprise #3**
> The dispatcher_id field captured in chain entries records the user's Active Directory display name, not their durable SID or federated SSO identity. A display-name change in AD — for example after a marriage or a department transfer — would render historical chain entries human-ambiguous, although the underlying chain integrity is unaffected. Federated SSO with SID-binding would resolve it.

### OT side IAM

Diana asked Esme for a screen-share from one of the HMI workstations on the operations floor. Aaron came back on the line.

The login screen showed `Operator_ControlRoom` as the account.

"Who is `Operator_ControlRoom`?"

Aaron's voice was slightly flat. "All three control-room shifts. Six operators on rotation. Same password. We rotate the password every 180 days."

"MFA?"

"No. Smart-card readers are on the procurement plan."

"Procurement plan timeline?"

"Phase 2. Twelve months out."

Diana wrote: *iFIX HMI. Shared `Operator_ControlRoom` account. Six operators, same password. No MFA. Smart cards 12 months out.*

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

Diana wrote: *AMI head-end. Two people with `meter_data_engineer` override authority. Override captured with original + corrected + reason. Reason is free-text.*

She kept going. "What version of OpenWay are you on?"

Aaron checked. "5.2."

"OpenWay 5.4 ships an integrity-checking option. You're version-locked behind it."

"5.4 was released last fall. We have it on the upgrade plan for next year. There are some downstream system compatibility issues — our meter data management is on the older Itron interface."

> **⚠️ Surprise #4**
> The Itron OpenWay AMI head-end is on version 5.2. Version 5.4 introduced an integrity-checking option that would cryptographically attest meter readings end-to-end. Pacific Crescent has the upgrade on the roadmap but is constrained by downstream system compatibility. Two `meter_data_engineer` accounts have override authority on individual meter readings; the override is captured with original value, corrected value, and a free-text reason, but the chain of custody from meter to head-end to MDM is not cryptographically attested.

> **⚠️ Surprise #5**
> The GE iFIX HMI workstations on the operations floor use a shared `Operator_ControlRoom` account across all three control-room shifts. Six operators rotate through it. Password rotates every 180 days. No MFA. The same shared identity is what stamps operator alarm-acknowledgment notes in `iFIXAlarmAck`. Smart-card readers are on the 12-month roadmap.

### Customer side IAM (CIS, OMS, Salesforce)

Diana shifted to the customer-side stack. The customized Oracle CIS, the OMS, the Salesforce CRM. She went through them quickly.

CIS had an internal user model — call-center reps had role-based access, with a small admin group that could adjust accounts. Audit trail was in an Oracle table that kept 18 months. Adjustments to billing records were logged with user ID and reason code. Adjustments to interaction notes were not.

OMS — the work-order system — had its own user model that mostly federated to AD via SAML. Work orders had version history. Field crew updates came in via tablets that authenticated with AD; some updates came in via paper forms transcribed by dispatchers, in which case the dispatcher's user_id was the one stamped on the work order.

Salesforce was Salesforce. Field-level audit on some fields. Not on others. Backups, not version history. Same shape as the diary baseline two weeks ago and same shape as Helmstad's Salesforce piece.

Diana wrote a half-page summary in her workbook. The customer side wasn't pretty, but it wasn't NERC's primary concern either. That was the PUC's concern.

---

## 🧪 12:00 PM — Lunch (Karen and Tom in the cafeteria)

The cafeteria sat on the second floor with a view east toward Mount Hood, white today against a flat blue sky. Karen and Tom found a corner table. The rest of the team was scattered — Diana and Mike on a sandwich run, Raj and Chen working through their morning notes.

Tom unwrapped a sandwich. "How are we framing this."

Karen had a salad and a notebook open. "Three tiers. AI side passes. OT side mostly fails. Customer-billing side is a different audit."

"NERC won't care about the customer side."

"NERC won't care about Salesforce. Right. NERC scopes to BES Cyber Systems. The pipeline leak detection is in scope under CIP-007. The AMI head-end is in scope under CIP-005 because of how it's networked into the OT segment. The PI historian is in scope under CIP-007. The OMS is partly in scope — work orders that touch BES assets are in, work orders for residential gas leak responses are in under the public-safety overlay even if they're not strictly BES."

"And the CIS and Salesforce —"

"PUC scope. PHMSA might care about parts of OMS for incident reconstruction. Salesforce mostly nobody cares about until a regulator subpoenas customer-service interactions in a litigation."

Tom nodded. "So the report has to split scope."

"Three columns. NERC + PHMSA in the first column. PUC in the second. Out-of-scope but operationally relevant in the third. That's how Esme's already thinking about it."

Tom took a bite of his sandwich. "What's the public-safety angle."

Karen looked at her notebook. "If a leak alert is dismissed and a week later a house explodes, what evidence do we have that the dismissal was reasonable at the time. The chain entry of the dismissal. The AI's confidence and inputs at that moment. The dispatcher's reason code. The dispatcher's identity. The model_id and version. All of that is in the chain, sealed, daily."

"And if the inputs were tampered with upstream of the chain?"

"Then we're back to PI. The chain captures what the AI saw. The chain doesn't capture whether what the AI saw was true. The PI historian is the upstream weakness."

Tom finished his sandwich. "That's the line."

"That's the line. Esme already knows it. We're going to write it down so her CEO knows it and her PUC commissioners know it and her insurance carrier knows it."

Karen closed her notebook.

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

Mike turned the laptop toward Karen and Esme. "Live alert, dispatched, sealed, verified. Twelve steps, four seconds, on the clean-room network."

Esme nodded once. She didn't say anything but her shoulders dropped a half-inch.

The team watched the dashboard for the next twenty-two minutes. A separate screen showed the work-order ID that had auto-opened in the OMS. The crew arrived at the location. A second voice came over the dispatch line.

"Brentwood crew on site. We have positive methane reading on the residential service line, three feet east of the meter. Small leak. Initiating standard isolation."

Yolanda put her hand to her headset. "Copy positive small leak, isolation in progress. Updating work order."

Esme's phone rang. She stepped two paces away to take it. The team heard her say "small leak, isolation, no evac" and "yes, I'll send a note to the comms team." She hung up.

She came back. "PUC liaison. They get a courtesy notice on every dispatch that escalates to crew on site. Standard."

Karen wrote in her notebook: *Live alert from prediction to crew on site to confirmed real leak in 22 minutes. Chain captured the AI's call and the dispatcher's decision. Field crew finding will land in chain via follow-up entry once crew leader logs it.*

Mike asked Esme: "When does the field finding hit the chain?"

"When the crew leader closes the work order in the OMS. The OMS pushes a follow-up event to the leak-detection service, and the service writes a follow-up chain entry referencing the original prediction. Outcome — true positive, false positive, equipment fault, environmental confound. Sealed in the daily seal that night."

"How long until the crew leader closes the work order?"

"Today probably six to eight hours from now. Sometimes longer if it's a complex repair."

Karen wrote: *Follow-up entry mechanism exists. Crew leader closes OMS work order. OMS pushes outcome to AI service. AI service writes follow-up chain entry referencing original prediction. Sealed that night.*

> **✓ Confirmation #4**
> A live, real alert observed end-to-end on the operations floor: AI prediction at 0.78 confidence, dispatcher action with reason code, work order auto-opened in OMS, crew on site within 22 minutes, confirmed real leak. Chain captured the prediction and the dispatcher decision in real time, verifier returned PASS in four seconds for the dispatch entry. Follow-up chain entry mechanism exists for the field-crew outcome and will land in tonight's seal.

The team filed back to the clean room.

---

## 🧬 2:00 PM — Pipeline Reality (Chen on the Sensor-to-AI Path)

Chen had been working through the data path for an hour by the time the rest of the team came back from the control room. He had a diagram on the whiteboard.

```
methane sensor -> cellular AMI feed -> PI historian -> AI ingestion adapter -> leak-detection service -> chain entry
```

He pointed. "Here is where the chain begins." He tapped *leak-detection service*. "Everything to the left of that arrow is unauthenticated."

Karen read the diagram. "The methane sensor itself. It just emits a value over the cellular link?"

"It emits a value. The value goes to the AMI feed, which writes it into PI, which the AI ingestion adapter reads. There's no signature on the sensor's emission. There's no signature when PI receives it. The first cryptographic operation in the path is the chain entry the AI service writes after it has already trusted the PI value."

"So if PI is tampered with —"

"The chain captures whatever PI says. Faithfully. The chain entry will say `methane_value: 4.6 ppm` because that's what the AI ingested. Whether 4.6 was the real reading at the sensor or whether someone overrode it in PI is invisible to the chain."

"Same shape as Stelvio's historian boundary."

"Exactly the same shape. Different industry, different vendor, identical boundary problem. The AI trusts what its upstream gave it. The chain attests to the AI's behavior. The chain does not attest to the upstream's behavior."

Chen wrote in the corner of the whiteboard: *Trust boundary = first cryptographic operation. Pacific Crescent's first cryptographic operation is at the AI service. Everything upstream is trust-by-policy.*

Karen wrote in her notebook: *Sensor-to-AI path is unauthenticated for the first three hops. Documented as a Phase 2 remediation: hash the value at the sensor side, carry the hash forward through PI, verify at AI ingestion. Itron OpenWay 5.4 enables it natively for the AMI portion. PI Server has integrity-extension options that Pacific Crescent has not turned on.*

> **⚠️ Surprise #6**
> The data path from methane sensor to AI ingestion is unauthenticated for the first three hops. The chain begins at the leak-detection service, which means the chain captures what the AI saw, not what the sensor actually measured. If the PI historian were tampered with — by one of the three engineers with override authority — the AI's chain entry would faithfully record the tampered value. The chain would verify PASS. The forensic trail would not detect the tamper.

Chen pointed at the diagram one more time. "If you fix the AMI side and the PI side, the trust boundary moves left to the sensor itself. That's where it should be in a public-safety AI."

---

## 📊 3:00 PM — Reconciliation Test (Four Leak Predictions, End to End)

The team did the reconciliation test together. Esme picked four leak predictions from the prior 120 days, spanning recent and old. The team traced each one end to end — from the methane sensor reading that fed it, through the AI's prediction, through the dispatcher's decision, through the OMS work order, through the crew leader's confirmed finding.

### Prediction 1 — 2026-05-02 — SE-Powell-22 — high confidence — true positive

Three days old. AI-side chain: PASS. Verifier four seconds.

Forward trace: dispatcher decision in chain (`dispatch crew, model confidence 0.86, no confound`), OMS work order opened automatically, crew leader closed it that evening with finding `true positive — service line leak — repaired`. Follow-up chain entry referenced the original prediction. PASS.

Backward trace: methane sensor reading at the time of inference — Aaron pulled it from PI within the 60-day audit retention. Original value, no edits, captured by `historian_writer` service account. Confidence the value is unaltered: high.

Status: clean four-of-four end to end.

### Prediction 2 — 2026-04-21 — NW-Lovejoy-18 — moderate confidence — false positive

Two weeks old. AI-side chain: PASS. Verifier four seconds.

Forward trace: dispatcher decision (`dispatch crew, confidence 0.62, light precipitation environmental confound but above threshold`), OMS work order opened, crew leader closed it with finding `false positive — wet soil signature — no leak`. Follow-up chain entry referenced original. PASS.

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

Status: AI-side PASS, dispatcher decision PASS, field-crew outcome PASS, sensor-side unprovable.

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

Diana spoke up. "And the original paper form?"

"In a filing cabinet at the crew dispatch yard. We scan the photos and attach them to the OMS work order. The paper form goes into a 7-year retention box."

"The work order's text in the OMS is the dispatcher's transcription of the paper form?"

"Or the crew leader's tablet entry, depending on whether the crew leader had the tablet that day. About sixty percent of work orders are tablet-entered, forty percent paper-transcribed."

Mike wrote: *OMS work-order linkage to chain is by ID, not by content hash. Paper-to-OMS transcription introduces a gap. The chain captures that the dispatch happened and that a follow-up was written. The chain does not capture whether the field finding's text in the OMS today is what the field finding actually said.*

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

Karen looked at the whiteboard.

"4 out of 4 AI-side PASS. 2 out of 4 trace back cleanly to the source sensor data. 4 out of 4 trace forward to the dispatcher decision. 2 out of 4 trace forward to a confirmed field-crew finding that we'd be willing to put in front of an investigator without a caveat. The chain works for the recent past. The legacy systems erode it as time goes on."

Esme looked at the whiteboard for a long beat. "That's what I expected."

"The 60-day PI retention is the single biggest forensic limit. After 60 days, the upstream half of the trail is gone. That's a NERC CIP-007 finding waiting to be written. PHMSA will care too — the integrity-management standards expect 5+ year evidence retention on incident-relevant data."

"Yes."

"The OMS transcription gap is a public-safety-investigation finding. If a leak is dismissed and the dismissal turns out to have been wrong, the chain will tell you the dismissal was logged and what the AI's inputs were and what the dispatcher said. It won't tell you what the field crew found three weeks earlier on a different alarm if that finding was paper-transcribed."

"Yes."

> **✓ Confirmation #5**
> Reconciliation test on four predictions traced end-to-end. AI-side: 4/4 PASS — chain integrity, verifier returns under 5 seconds for any single entry, twelve-step output. Dispatcher decisions: 4/4 PASS — captured with user, reason code, model inputs, model version. Forward trace to OMS work order: 4/4 by ID linkage. The chain holds where it was designed to hold.

> **⚠️ Surprise #7**
> Reconciliation test reveals two erosion points outside the AI scope. Backward trace to sensor source — 2/4 PASS, 2/4 UNPROVABLE due to PI's 60-day audit retention. Forward trace to confirmed field-crew finding — 3/4 PASS (recent + tablet-entered), 1/4 UNPROVABLE due to paper-form transcription. The chain does what it was scoped to do; the legacy systems on either side of the chain limit how far the forensic trail can be carried.

---

## 😬 3:45 PM — Friction in the Room (OT Engineering Defensive)

Esme had pulled in three OT engineering staff at the start of the 3:00 PM session — Aaron from PI, a SCADA engineer named Hugh, and the senior engineer over the AMI head-end, a woman named Pavithra. They had been listening from the back. By 3:45 they were not listening passively anymore.

Hugh spoke first. "I want to push back on the framing. You're describing GE iFIX like it's broken. iFIX is a NERC CIP-compliant HMI deployed in compliance with our CIP-007 controls. The audit log captures the INSERT. That's the standard."

Karen put her pen down. "Hugh, I hear you. We're not writing iFIX up as broken. We're writing it up as a system whose audit semantics are INSERT-only, which means a subsequent UPDATE to an operator note doesn't generate a record. That's a finding about the audit shape, not about iFIX as a product."

"Auditors always make the same finding. It doesn't lead anywhere. The product doesn't support UPDATE auditing without a third-party add-on."

"That's right. And that's exactly what the finding will say. INSERT-only audit, third-party add-on available, evaluate against the cost of a tamper scenario in the public-safety context."

Hugh sat back. He didn't say anything for a beat.

Pavithra picked it up. "AMI 5.4. We know about the integrity-checking option. We can't deploy it because our MDM is on the older interface."

"Right. And the finding will say: 'AMI 5.2 is version-locked behind an integrity-checking option that will close a finding when MDM is upgraded.' That's a roadmap item, not a violation."

"It's going to read like a violation in front of the PUC."

Esme cut in. "It's not going to. The team is identifying gaps for our roadmap, not writing a NERC violation. We've already self-disclosed two of these to the PUC liaison in the last quarter."

Pavithra and Hugh both looked at Esme. The friction in the room shifted by a quarter-turn. Esme had clearly had this conversation with them before, but never in front of outside auditors.

Aaron was the quietest of the three. He spoke last. "The 60-day PI retention. I didn't set that. I inherited it from the previous PI admin. When I came on, I asked about extending it. The disk volume isn't sized for 18 months of audit trail at our sensor density."

Karen wrote: *PI 60-day retention disclosed by current admin as inherited. Disk sizing constraint. Phase 2 includes disk re-sizing on the historian server.*

"Aaron, that's a clean disclosure. We'll write it that way. The finding will be on the retention shape and the sizing constraint, not on you personally."

The friction in the room subsided. The three engineers didn't relax exactly, but they stopped pushing back. Esme had defused it with one sentence and Karen had received the disclosures cleanly.

Karen made a note in the margin: *NERC engineers always defensive. Esme manages it well. Disclosures land cleaner because she creates the space for them.*

---

## 🔍 4:30 PM — The Public-Safety Question

Esme had asked the engineers to step out for the last hour. The team was alone in the clean room with her. The afternoon sun was angling through the high windows. She had her elbows on the table and her chin in her hands.

She looked across at Karen.

"Let me ask you a question I've been trying to find a clean answer to for a year."

Karen waited.

"If a leak alert is dismissed and a house explodes a week later, what evidence do we have that the dismissal was reasonable at the time?"

The team went still.

Karen took a long breath. "Sober answer. The chain entry of the dismissal — that gives you the dispatcher's identity, the reason code, the timestamp. The AI's confidence and inputs at that moment — that's in the same chain entry. The model_id and model_version — also in the chain. The seal record for that day — ties everything to a public key with a known fingerprint."

"And that's enough?"

"That's enough to demonstrate, with cryptographic confidence, that the dismissal happened the way the chain says it happened. That's what your insurance carrier and the state regulator and the Class A NERC auditor will want to see in that scenario."

"And if the inputs to the AI were tampered with in the historian?"

Karen paused.

"Then we cannot prove the inputs were tampered with. And we cannot prove they weren't. The historian is your weakest evidence in that scenario. The chain entry will faithfully record the inputs the AI saw. If those inputs were already tampered with when the AI ingested them, the chain reflects the tampered version. The audit trail in PI rolls at 60 days, so going back further, even the question of who-changed-what becomes unanswerable."

"And the dispatcher in that scenario — they made what they thought was a reasonable call based on what the AI showed them."

"They made a defensible call based on the inputs they saw. The chain proves that. The question of whether the inputs were what the sensors actually measured is a different question that the chain cannot answer."

Esme sat back. Her chin came off her hands. "That's the answer I was afraid of."

"The remediation is the Phase 2 work. Hashing at the sensor — Itron OpenWay 5.4. Hashing at the historian boundary — PI integrity extension. Once the chain extends to the sensor, the question of input authenticity becomes answerable. As long as the chain starts at the AI service, the answer to that question depends on the historian's discipline, and the historian's discipline is operational, not cryptographic."

Esme nodded. She wrote something in her own notebook.

Tom leaned forward. "Esme, the value of doing this assessment now is precisely so that you have the answer in your hand before you ever need it. The PUC is going to ask. Insurance carriers are starting to ask. If a public-safety event ever does happen, having the chain in place for the AI side and a documented Phase 2 plan for the upstream is a substantially better position than not having either of those things."

"I know. I'm going to be straight with you — I've been trying to fund Phase 2 for six months. Today's report is the lever I needed."

Karen wrote in her notebook: *Public-safety question. Esme's been preparing for this question for a year. The report has to answer it directly. The line is: chain proves the AI's behavior, chain does not prove the upstream's behavior, Phase 2 closes the gap. Document.*

> **✓ Confirmation #6**
> The chain provides cryptographic evidence sufficient to demonstrate that an alarm dismissal was made on the inputs and reason recorded, by the dispatcher recorded, at the timestamp recorded. In the public-safety scenario where a regulator, insurance carrier, or litigation discovery process asks how a dismissed alarm was reasoned through, the chain produces a verifiable artifact in seconds.

> **⚠️ Surprise #8**
> The chain cannot extend its evidentiary value upstream of the AI service. If a leak prediction were dismissed based on tampered sensor inputs, the chain would faithfully record the tampered inputs and the dispatcher's reasonable decision based on them. Detecting upstream tamper requires extending cryptographic integrity to the AMI head-end (Itron OpenWay 5.4) and to the historian (PI integrity extension). Both are on the Phase 2 roadmap. Both are gating items for the public-safety evidentiary story Pacific Crescent will need.

---

## 🌆 5:30 PM — Auditor Debrief

The team reconvened in the clean room. Esme stayed. Coffee was cold. Mount Hood had moved into the late-afternoon haze.

Karen stood at the whiteboard. Three columns.

| Tier | Status |
|---|---|
| AI side (`pipeline-leak-detection`, `pipeline-integrity-trending`) | 0 Gaps, 0 Partials |
| OT side (PI historian, GE iFIX SCADA, Itron OpenWay AMI, OMS work-order linkage) | 5 Gaps, 6 Partials |
| Customer-billing side (Salesforce, CIS, OMS customer interactions) | 3 Gaps, 4 Partials |

"That's the shape. Three tiers. One passes. Two don't. Today is Stelvio with public-safety consequences."

Esme stood with her arms crossed, listening.

"AI side." Karen pointed. "TesseraSeal. Six confirmations — chain integrity verified across nine months and 1.6 million inferences, append-only ledger behavior under direct DB mutation attempt, credential rotation under chain (eight rotations sampled, all PASS), live alert observed end-to-end on the operations floor with verifier PASS in four seconds, four-of-four reconciliations PASS on the AI side, the public-safety evidence question answerable in two paragraphs. The AI side maps cleanly to NERC CIP-007 for the AI/ML control points and to PHMSA pipeline-integrity expectations for the leak-detection scope. IEC 62443-3-3 SR 6.1, SR 6.2, and SR 7.5 are explicitly satisfied for this service."

She moved to the OT column.

"OT side. Five gaps. PI historian retention is 60 days with three engineers holding override authority — the upstream half of the forensic trail goes invisible past 60 days. GE iFIX backing store captures INSERTs only on operator alarm-acknowledgment notes; subsequent UPDATEs are not recorded. Itron OpenWay AMI is on 5.2, version-locked behind the 5.4 integrity-checking option. The HMI workstations use a shared `Operator_ControlRoom` account across three shifts with no MFA. The dispatcher_id captured in chain entries is the AD display name, not federated SSO with SID-binding. Six partials around OMS work-order content-coupling, paper-to-OMS transcription, AMI override authority discipline, sensor-to-AI path authentication, IEC 62443-3-3 SR 7.5 coverage on legacy systems, and historian disk-volume sizing."

"Customer-billing side. Three gaps. Salesforce overwrite shape on customer-interaction notes. OMS work-order linkage from field finding to AI prediction is by ID only, not by content hash. CIS audit-trail retention on interaction notes is shorter than on billing records. Four partials around backup-vs-change-history, retention variance across three states' regulators, and customer-CRM IAM federation. Phase 4 territory. The CIO's budget conversation."

Karen put the pen down.

"Three observations to close."

"One. The chain works. Pacific Crescent has been running it on the leak-detection service for nine months. Today we observed a real alert dispatched to a real crew and confirmed as a real leak — chain captured every step that mattered, verifier returned in four seconds. The investment paid off. When the PUC asks, you have an answer."

"Two. The chain is on the part that decides whether to dispatch. The chain is not on the part that produces what the decision is made about. Stelvio's seam was the same shape — AI sealed, OT mutable. Pacific Crescent's seam is in the same place but the consequence is different. A wrong reading and a dismissed alarm at Stelvio is a yield-loss claim on a heat of steel. A wrong reading and a dismissed alarm at Pacific Crescent is a public-safety incident. The seam matters more here. Closing it through Phase 2 is not a compliance project. It is a public-safety project that happens to also close a compliance gap."

"Three. The roadmap you already have is the right roadmap. Phase 2 in 12 months — AMI head-end to 5.4, PI integrity extension, PI retention extended via disk re-sizing, MFA on the HMI workstations, federated SSO with SID-binding for the dispatcher application. Phase 3 in 18 months — OMS work-order content-coupling, paper-to-OMS gap closure, IEC 62443-3-3 SR 7.5 coverage on legacy. Phase 4 deferred — customer-CRM. The order is correct because Phase 2 closes the public-safety story first. Document it that way for the CFO."

Esme nodded. "What do I take to the CEO Friday?"

Tom answered. "The three-tier summary. The roadmap costed against the three-tier summary. The four-second verifier output Mike captured during the live Brentwood alert as a real artifact showing what the prior investment delivered. The reconciliation table from this afternoon — four predictions, what worked, what didn't, where the legacy erosion is. And the public-safety paragraph from the 4:30 conversation — the chain proves the AI's behavior, the chain does not prove the upstream's behavior, Phase 2 closes the gap, the gap is fundable today."

"Send me the Brentwood verifier capture."

Mike held up his phone. "Got it on the laptop. I'll send it tonight."

"And the report?"

"Thursday morning. Before your Friday CEO review."

Karen closed her notebook. "We'll send it."

Esme's shoulders dropped that half-inch again. She didn't smile but she nodded twice.

The team packed up. Raj and Luis loaded boxes of evidence into the rental SUV. Diana and Elena said goodbye to Esme at the badge desk. Mike and Chen took one last look at the operations floor on the way past — the gas-distribution map glowing green except for two amber dots, both being watched, both at confidence below dispatch threshold.

Karen walked out last. She paused at the badge desk and looked back through the glass at the operations floor. Yolanda was still on shift. The methane sensor readings were still updating every fifteen seconds. The Brentwood crew had finished isolation and were filing the work-order paperwork on a tablet.

> **🔍 Karen's note (internal):**
> *It never is. The closer the AI sits to a public-safety decision, the more the chain matters. The further the historian sits from the AI, the more the chain doesn't reach.*
>
> *Today the chain reached far enough to dispatch a real crew to a real leak in real time and prove it after the fact. Today the chain did not reach the sensor. Phase 2 closes the gap before the gap closes a neighborhood.*

---

## ✅ vs ❌ — The Three-Tier Summary

### ✅ AI Side (TesseraSeal — `pipeline-leak-detection`, `pipeline-integrity-trending`)

| Item | Status |
|---|---|
| Chain integrity (HMAC + Merkle + daily Ed25519 seal on Thales Luna PCIe HSM, CIP-categorized network) | PASS |
| Append-only ledger behavior under direct DB mutation attempt | PASS — verifier catches at HMAC layer |
| Multi-entry tamper attempt | PASS — verifier catches at Merkle/seal layer |
| Credential rotation under chain | PASS — eight rotations sampled across nine months, all PASS |
| Inference -> chain latency | ~200 ms observed live on operations floor |
| Verifier latency | ~4 seconds for any single entry, 12 steps |
| Live alert observed end-to-end (Brentwood, 2026-05-05) | PASS — prediction, dispatch, crew on site, confirmed real leak, chain captured each step |
| Reconciliation test (4 predictions, AI-side) | 4/4 PASS |
| Dispatcher decision capture (identity, reason code, model inputs) | PASS — every dispatcher action sealed in chain |
| Field-crew follow-up entry mechanism | PASS — OMS closes work order, leak-detection service writes follow-up referencing original prediction |
| IEC 62443-3-3 SR 6.1, SR 6.2, SR 7.5 mapping for AI scope | PASS — explicitly satisfied |
| NERC CIP-007 readiness for AI/ML control points | Demonstrable |
| PHMSA pipeline-integrity evidence for leak-detection scope | Demonstrable |
| Public-safety evidence question (dismissed-alarm scenario) | Answerable in two paragraphs with chain artifact |

### ❌ OT Side (PI historian, GE iFIX SCADA, Itron OpenWay AMI, OMS work-order linkage)

| Item | Status |
|---|---|
| PI historian audit-trail retention | 60 DAYS — past 60 days, who-changed-what is unrecoverable |
| PI override authority | THREE ENGINEERS with `PIWorld\db_admin`, no second control |
| Sensor-to-AI path authentication | NONE — first cryptographic op is at the AI service, three hops upstream are trust-by-policy |
| GE iFIX `iFIXAlarmAck` UPDATE audit | NONE — INSERT captured, subsequent UPDATEs silently overwrite |
| HMI workstation account | SHARED — `Operator_ControlRoom`, six users, no MFA, smart-card readers Phase 2 |
| Dispatcher_id binding in chain | DISPLAY NAME, not federated SSO with durable SID |
| Itron OpenWay version | 5.2 — version-locked behind 5.4 integrity-checking option, MDM compatibility constraint |
| AMI override authority (`meter_data_engineer`) | TWO accounts, free-text reason code, override captured but no chain extension |
| OMS work-order linkage to chain | ID-LINKED, not content-hashed — text changes after follow-up entry undetectable |
| Paper-to-OMS transcription | ~40% OF WORK ORDERS — dispatcher transcribes paper, chain references work order ID, content-coupling broken |
| Reconciliation test sensor-side trace | 2/4 PASS, 2/4 UNPROVABLE due to PI 60-day retention |
| IEC 62443-3-3 SR 7.5 on legacy systems | NOT MET |
| Phase 2 remediation timeline | 12 months — AMI 5.4 upgrade, PI integrity extension, PI disk re-sizing, MFA on HMI, federated SSO |
| Phase 3 remediation timeline | 18 months — OMS content-coupling, paper-to-OMS closure, IEC 62443-3-3 SR 7.5 on legacy |

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
| Phase 4 (customer-CRM) | Deferred pending CIO budget. Document explicitly. |

---

## 🧾 Final Assessment Theme

> *"The chain is on the part that decides whether to dispatch. The chain is not on the part that produces what the decision is made about. Pacific Crescent's customers don't know that line. Their regulators don't yet either. A dismissed alarm and an exploded house would teach the line in the worst possible way. Phase 2 teaches it the right way, in 12 months, on a fundable budget."*

Pacific Crescent Power & Gas demonstrates AI-decision integrity within scope. The investment in TesseraSeal on the leak-detection service returns verifiable provenance under NERC CIP-007 for the AI/ML control points, PHMSA pipeline-integrity expectations for the leak-detection scope, and IEC 62443-3-3 SR 6.1, 6.2, and 7.5. The verifier runs in four seconds on the operations clean-room network. A live alert observed end-to-end during the engagement — prediction, dispatch, crew on site, confirmed real leak — left a sealed forensic trail that will satisfy a NERC auditor, a PHMSA investigator, and an insurance carrier.

Outside the AI scope, integrity is operational discipline. The OT side carries the same mutability shape as the manufacturing OT engagement the team documented two weeks ago — different vendors, identical boundary problem. The PI historian's 60-day audit retention erodes the upstream half of the forensic trail past 60 days. The HMI's INSERT-only audit semantics on operator notes erode the operator-side accountability. The shared `Operator_ControlRoom` account compounds. The Itron OpenWay 5.2 version-lock keeps the cryptographic boundary at the AI service rather than at the sensor. The OMS work-order content-coupling gap leaves a paper-transcription seam that grows with time-since-incident. None of these surprised Esme. All of them are on the roadmap she has been building for a year.

The customer-billing side is a different audit — PUC scope rather than NERC scope, with the same overwrite-and-backup shape that has appeared in five of the team's last six engagements.

The framing for the CEO and the board is straightforward. Phase 2 is not a compliance project. Phase 2 is a public-safety project that happens to also close a compliance gap. The chain proves the AI's behavior. The chain does not prove the upstream's behavior. Closing the upstream gap before the gap closes a neighborhood is the operationally correct sequence and it is the morally correct sequence and it is fundable today.

Three tiers, one report, one severity scale, three remediation timelines, one public-safety paragraph that decides everything else. Pacific Crescent knows where the seam is. Their CEO is about to. Their PUC commissioners need to.

---

*End of diary. Filed Tuesday evening. Report drafted Wednesday and Thursday. Delivered Thursday morning before the Friday CEO review.*
