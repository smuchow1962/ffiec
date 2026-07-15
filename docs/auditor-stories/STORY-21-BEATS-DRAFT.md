---
status: plot-beats-draft (not finished narrative)
audience: Steve (review and edit before full draft)
purpose: Story 21 plot skeleton — Texas regional bank merger, marketing-thread focus, first vendor-side Dawn appearance
last updated: 2026-05-11
---

# Story 21 — Brazos Federal × Mission Plaza Bank (working title)

## What this story is

A confirmation audit at a Texas regional bank ~6 months into a merger integration. The acquirer has TesseraSeal in production for 15 months; the acquired bank had a partial deployment going to seal at the merger boundary. The signature element is the **marketing chain** — Salesforce + Marketo + Adobe Experience Cloud campaign data spanning the brand transition, with a cross-vendor anchor pattern proving the legacy Marketo dump is byte-equal to what post-close TesseraSeal contains.

The structural pivot: **Dawn arrives at lunch as the MMPWorks TesseraSeal liaison**, invited by both sides. First story where she sits on the vendor side of the table. Raj has stepped into the Lead Auditor seat. The team handles the moment with the same dignity that defined her time as their lead — Tom's catch-up at lunch is brief, the work is still the work.

## Setting and posture

- **Brazos Federal Bancshares** — Houston-HQ Texas regional bank, ~$45B consolidated assets, Texas + Oklahoma footprint, OCC-supervised. **TesseraSeal in production for 15 months across all customer-facing surfaces.** Mirrors real-world Prosperity scale and footprint.
- **Mission Plaza Bank** — San Antonio-HQ community bank, ~$3.2B assets at close, acquired by Brazos Federal in a $269M cash-and-stock deal that closed Feb 1, 2026. **Mission Plaza had partial TesseraSeal — 6 months of deployment, AI-decisioning side only.** Merger-integration period runs until November 2026; brand transition in progress.
- **Date of engagement:** mid-September 2026, ~7 months post-close, ~6 weeks before full operational integration to Brazos brand.
- **Posture:** integration audit. Two chains have been running in parallel since close; the cross-vendor anchor was placed at close; the audit confirms the merged chain holds at every boundary.
- **Regulators in scope:** OCC (lead, both legacy and acquired charters), CFPB (UDAAP on rebrand-period customer communications), Texas Department of Banking (state-chartered Mission Plaza legacy).

## Cast — what's different from prior stories

**Auditor team:**
- **Raj** is the new Lead Auditor. He took the seat 8 months ago when Dawn left to join MMPWorks. The drive-in monologue pattern continues with him now leading it — first pairing him with **Diana** for this engagement. His "It never is" calibration is still being broken in; he says it slightly differently than Dawn did, and the team has noticed.
- **Dawn** is now MMPWorks's lead TesseraSeal liaison. She is not on the audit team's roster for this engagement. She arrives at noon on Day 1 at the joint invitation of Brazos's CMO and the audit team — pre-arranged, on the schedule, fully disclosed in the engagement letter.
- **Tom** continues as internal-audit liaison specialist; partners with both Brazos's CAE and (cross-engagement) Mission Plaza's legacy CAE. His posture this engagement: **mediating the merger-integration culture gap between acquirer and acquired internal-audit teams**. He treats Dawn the same way he always has when she's in the room.
- **Other team members:** Elena (CRM), Mike (App/API), Luis (DevOps/logs), Chen (Data engineering/ETL) — all in their usual seats. The marketing thread puts Elena front-and-center on Day 1 morning; Chen carries the legacy-Marketo-anchor verification.

**Bank cast:**
- **Brazos Federal CAE** — early-50s woman, was CAE at Brazos for 6 years pre-merger, integration veteran (Brazos's third acquisition under her tenure). Calm. Knows the chain.
- **Mission Plaza legacy CAE** — late-40s man, two years in role at Mission Plaza, now reports to Brazos's CAE under the integration. Less familiar with TesseraSeal; was the one who stood up the partial deployment 6 months before close.
- **Brazos's CMO** — the one who invited Dawn. Was nervous about the brand-transition campaign window; called MMPWorks 60 days before close, asked for an integration plan, got one.
- **Mission Plaza's legacy CMO** — Marketo-stack person, stayed through integration as Director of Brand Transition, scheduled to roll off in November.

## The day's arc

**7:50 AM — drive-in (Raj + Diana, in from Hobby Airport).** Raj names eleven prior contexts and slots this engagement: *"Post-merger integration with marketing-data as load-bearing. Closest precedent is Atrio's coordinated-examiner-room — but here the chain crosses a vendor-merger boundary inside the acquiring bank itself."* Diana asks how he feels about Dawn being on the other side at lunch. *"It's the right answer. We trained her; she trained the next layer; the work compounds. Today is just the work."*

**8:30 AM — kickoff at Brazos's Houston operations center.** Both CAEs in the room. Brazos's CAE opens with the integration timeline. Tom asks the standard four questions; the legacy-side CAE answers crisply on three and pauses on the fourth (cross-charter retention policy). Note logged. Raj's first "It never is" — a little quieter than Dawn's was.

**9:30 AM — Elena opens the marketing-stack walkthrough.** Brazos uses Salesforce Financial Services Cloud + Salesforce Marketing Cloud + Datadog (observability). Mission Plaza used Adobe Marketo Engage + HubSpot Marketing Hub + Splunk. The cross-vendor anchor was placed at 23:59:59 UTC Jan 31, 2026 — a hash of the entire Marketo campaign-history export, sealed into the Brazos TesseraSeal chain on Feb 1 at 00:00:01 UTC. Elena pulls the seal, runs `herald-verify` — exit 0.

**11:15 AM — the first reconciliation pivot.** Five records traced end-to-end. The most interesting: customer #4, who received a pre-merger Mission Plaza "you're pre-qualified for a $50K HELOC at 6.875%" email on Jan 17, then a post-merger Brazos rebrand email on Feb 4 ("Same great service, new name"), then applied for a HELOC on May 22, got AI-screened, was approved at 7.250% on a different product structure. The audit question: can the bank prove what the customer was told, when, by whom, and that the post-merger product decision honored or properly superseded the pre-merger marketing promise? **It can.** All three marketing events chain through; the AI underwriting decision references the prior marketing-event hash in its `audit.marketing.prior_offer.*` field family.

**12:00 PM — Dawn arrives.** She walks into the Brazos conference room in a navy suit, carrying a laptop and two coffees. Greets the team by first name. Tom takes one coffee and offers her the second seat at the working lunch. The conversation for 90 seconds is personal — Tom asks about Steve, asks about the new house, then everyone moves on. *Brazos's CAE has already met Dawn in her vendor role twice in the last quarter; this isn't anyone's first meeting.* Dawn opens her laptop and pulls up MMPWorks's view of the cross-vendor anchor architecture. Elena moves a chair so she can see.

**1:30 PM — the Marketo legacy-anchor deep-dive.** Dawn walks Mission Plaza's legacy CMO and Chen through how the Marketo dump was hash-anchored. The dump was a 2.7-TB tar.gz of campaign history, customer lists, A/B variant data, send-time logs, click-stream telemetry. Hash committed at close. Mission Plaza's legacy CMO has been carrying anxiety about this for seven months — was the hash *actually* what the Marketo backup contained, or was there drift in the 47 minutes between Marketo's last write and the close timestamp? Chen pulls the receipt: Marketo's audit log shows last write at 23:42:18 UTC Jan 31; the hash was computed at 23:59:30; the seal sealed at 00:00:01 Feb 1. **No drift window.** Mission Plaza's legacy CMO exhales for the first time in seven months.

**2:30 PM — the second reconciliation pivot: a UDAAP-adjacent rebrand email.** On Feb 4, the rebrand campaign fired to ~340,000 Mission Plaza legacy customers. One variant accidentally referenced Mission Plaza's old fee schedule ("free overdraft protection") which Brazos does not offer at zero cost. The campaign-stop button was hit at 14 minutes; 47,000 emails had already sent. The bank's posture: settled the regulatory question with the CFPB in March via a customer-restitution package. **The audit question today: can the bank produce the exact list of 47,000 customers who received the misleading variant, the exact text they received, the timestamp, and the restitution receipt?** Mike pulls all four in 6 minutes. Verifier passes. Brazos's CAE: *"This is why we did the deployment."*

**4:30 PM — the client question, and Dawn's nagging question answered.** Brazos's CAE asks Dawn (in her vendor-liaison role, not Raj): *"If a Mission Plaza legacy customer sues Brazos three years from now over a pre-merger marketing promise we didn't honor, can we prove what they were actually told?"*

Dawn's answer (sober, concrete, the pattern preserved across her seat change): *"Yes. The Marketo legacy data is hash-anchored at close; the post-close marketing events are in your TesseraSeal proper; the AI underwriting decision references both. The chain spans every relevant event from the customer's first interaction in 2024 forward. Three-year window is well inside retention. The cross-vendor anchor is the load-bearing piece — Chen demonstrated its byte-equality this morning; the test will reproduce in court."*

She pauses for half a beat before adding the part that's been carrying her for two years:

*"I had this question once, when I was sitting on your side of the table. The engagement was Hill Country Federal Credit Union — a marketing-AI vendor swap, AWS-resident, NCUA exam coming. The byte-equality demonstration worked cleanly there, but I wrote a note in the engagement file: 'what happens when the substrate moves?' I didn't have a way to test it. When I came to MMPWorks, that was the first thing I put on Steve's desk. He had the §10.40 extension already roughed out. We built it. Today is the first production engagement that exercises it. The answer is: the chain doesn't care which cloud the artifacts live on. The substrate-trust boundary is just another anchor."*

Tom gives her a microscopic nod. Raj's expression doesn't change. The room understands that an auditor's two-year nagging question just landed as a vendor's deliverable, and that the seat change Dawn made earned its keep in this single 90-second answer.

**5:15 PM — Steve joins by video bridge for 20 minutes.** A specific technical question about the cross-vendor-anchor primitive came up during the Marketo deep-dive (a question about how the anchor handles a vendor's *own* internal hash format changing across a major Marketo platform upgrade in 2027). Steve walks the room through §10.21 + §10.40 + the migration path. Confirms. Signs off. Logs off.

**6:00 PM — debrief.** One Nit (the legacy-side CAE's cross-charter retention policy gap). Zero Gaps. Zero Partials. Engagement closes 18% under budget — not Northbridge's 30%, because the marketing-chain verification took more cycles than expected, but still well under.

## The marketing thread — what's distinctive

This is the first story where **marketing data is load-bearing audit content**, not background. The audit treats:

- **Salesforce Marketing Cloud + Marketo Engage** as production AI-decisioning surfaces (next-best-action, send-time optimization, audience segmentation are all AI-driven)
- **Marketing-to-credit pivot points** as in-scope for the chain (a marketing-driven offer that informs a credit decision is part of the credit-decision evidence)
- **Brand-transition campaign windows** as UDAAP risk periods (when banks make the most marketing mistakes, and when chain-of-custody pays its premium)
- **Cross-vendor anchors at vendor-platform boundaries** as a reusable pattern — Marketo's data leaving Marketo, hash-anchored on the way out, becomes byte-equal evidence inside the receiving chain

The spec sections exercised: §10.21 (cross-vendor model-handover), §10.40 (cross-vendor chain-merge anchor — this story's signature primitive), §10.11.1 (`audit.ecoa.adverse_action.*` family — for the HELOC reconciliation), plus a hypothetical §10.x `audit.marketing.prior_offer.*` family (TBD — may need to add to spec if not already there).

## Signature moments

1. **Dawn at the lunch table.** First time she sits across from her former team in the vendor seat. Tom's 90-second catch-up. Elena's chair-shifting. The team's confirmation that the personal-arc evolution is fully metabolized.
2. **Mission Plaza CMO's exhale.** Seven months of anxiety about a 47-minute drift window resolved in 6 minutes by Chen pulling a Marketo audit-log receipt. The kind of moment that justifies a chain to a non-technical executive.
3. **The 47,000-email UDAAP scene.** The bank can produce the exact misleading text, the exact list, the exact timestamps, the restitution receipt — in 6 minutes. The CAE's line: *"This is why we did the deployment."*
4. **Dawn's two-year-old nagging question, answered.** The character beat that earns her seat change. She carried *"what happens when the substrate moves?"* across the table from auditor to vendor; she got Steve's attention on day one at MMPWorks; the §10.40 extension shipped; today she gets to deliver the answer to a CAE. The room reads it as the work compounding, not as a personal triumph. Tom's microscopic nod is the only acknowledgment. (Callback structure: explicitly names "a Story-12 era engagement" as the original context, so Story 12 and Story 21 sit in dialogue.)
5. **Steve's 20 minutes by video.** Vendor-principal walking a room through a primitive his wife on the vendor team and his former colleagues on the audit team both already understand. Confirmation, not exposition.

## Held-line principles (per Stories 14-17 doctrine)

- **Raj's authority is load-bearing.** He runs this engagement; the story doesn't undercut him by framing it as "Dawn is missed." She's not missed; she's present, in her new seat, exactly where she should be.
- **Dawn's professional integrity remains load-bearing.** Her vendor-side answer at 4:30 PM is the same shape as her audit-side answers across Stories 01-11 — sober, concrete, specific. The seat changed; the discipline didn't.
- **Tom's posture is preserved.** He treats Dawn the same way he always has. The catch-up is short. The work is the work.
- **Steve appears with dignity, not drama.** 20 minutes, technical, signs off, logs off. No spotlight.
- **The marketing thread is operational, not editorial.** No "marketing is the new frontier" speeches. The chain extends across the marketing-data boundary because the spec says it does, and the bank built it because the regulators will ask.

## Implementation prerequisites

- Confirm the `audit.marketing.prior_offer.*` field family exists in the current spec, or add it. **(Open task.)**
- Confirm §10.40 cross-vendor chain-merge anchor language supports vendor-platform-upgrade migration semantics (the 2027 Marketo-platform-upgrade question Steve addresses by video). **(Open task.)**
- The story assumes Story 20 lore evolved from "Northbridge acquires" to "MMPWorks stays independent." **Confirm with Steve before writing the full draft.**

## Stack realism — mirroring TPB / Prosperity for narrative truth

Per the real-world TPB stack discovery on 2026-05-11 (`memory/tpb_stack_intel.md`), the fictional bank stacks in this story can be drawn from real-bank realism for additional weight. Recommended assignments:

**Mission Plaza Bank (fictional TPB analog) — pre-merger stack:**
- Jack Henry **Banno** (digital banking)
- **JHA Payment Solutions** (payments)
- **Yello** (banking-specific marketing CRM + account-opening journey)
- **Insider** (cross-channel marketing automation, AI personalization for offers + retention)
- **AWS-resident** — S3 backing storage, AWS KMS / CloudHSM for root-signature custody
- ASP.NET app stack
- Splunk for SIEM/observability (typical community-bank pairing)
- **Herald.Core.Aws** runs the chain — seal job on S3-resident artifacts

**Brazos Federal (fictional Prosperity analog) — acquirer stack:**
- Fiserv DNA (core banking — common at $40B+ tier)
- Q2 Holdings (digital banking)
- Salesforce Financial Services Cloud (CRM)
- Salesforce Marketing Cloud + Marketo Engage (marketing automation)
- Datadog (observability)
- **Azure-resident** — Blob Storage, Azure Key Vault Managed HSM, multi-region resilience per §10.15
- **Herald.Core.Azure** runs the chain — Azure-side runtime, larger HSM partition for the consolidated key-management footprint

**The cross-cloud anchor thread (this story's signature technical beat).** Mission Plaza is on AWS; Brazos Federal is on Azure. The merger integration is therefore a **cross-cloud chain consolidation** — Mission Plaza's chain artifacts in AWS S3 / Herald.Core.Aws get hash-anchored at close and the anchors land inside Brazos's Azure Blob Storage / Herald.Core.Azure chain. Two cloud substrates, two Herald runtimes, one continuous evidentiary chain.

**The narrative discipline: cross-cloud reads as routine.** Per Steve's 2026-05-11 design statement — *"making it almost routine is the goal"* — the cross-cloud byte-equality demonstration must not read as a heroic technical climax. Chen pulls the receipts in 6 minutes. The verifier returns PASS. The team treats cross-cloud as just another anchor type. Raj's drive-in monologue (or its equivalent moment) names this explicitly: *"Cross-cloud is the new same-cloud. The chain doesn't care where the artifacts live."* The audit team's calm familiarity with the operation is itself a form of evidence — chain-of-custody discipline scales because the routine has been engineered into the spec, not because the audit team happened to be exceptional.

**Spec dependency.** The full story is technically accurate only if PRD-1's §10.40 (cross-vendor chain-merge anchor) or a sibling section makes cross-cloud explicit. This is currently the parallel spec-hardening session's lane (`docs/review-2026-05-11/*`). Surface to Steve before full-draft to confirm.

The integration period creates the story's chain-of-custody pressure: **two ML-driven marketing stacks (Insider + Marketo/Salesforce) running in parallel, two digital-banking surfaces (Banno + Q2) presenting to the same legacy-customer population, two account-opening journeys (Yello + Salesforce FSC), and AI-driven decisioning at every customer-touchpoint hand-off.**

This is exactly the audit problem the chain-of-custody spec was designed to address. The story's marketing thread becomes structurally inevitable rather than editorial — the bank's reality is that marketing data is now the noisiest part of the merged chain, and the audit lives or dies on whether the cross-vendor anchors at the Insider→Salesforce, Marketo→Salesforce, and Banno→Q2 boundaries are byte-equal across the chain seal.

## Open questions for Steve

1. **Slot: confirmed.** This is Story 21 (post-wedding, post-Dawn-transition). Story 12 is preserved as an open foresight-cluster slot. Sequence: Stories 14-17 (Dawn-as-Lead, recusal protocol, foresight pattern) → Story 20 (wedding + whatever else) → Story 21 (post-wedding new normal).
2. **Story 12 callback.** The 4:30 PM moment explicitly names "a Story-12 era engagement" as the context where Dawn's nagging question was first written. Story 12's actual plot needs to support that callback (single-cloud, single-vendor marketing handover at a financial institution). Brainstorm in flight; Steve to confirm before full draft.
3. **Lore reconciliation: confirmed.** Per Steve's 2026-05-11 messages, Story 20 is being rebuilt with only "Steve and Dawn marry" as the fixed point. MMPWorks stays independent. Story 21 is consistent with that.
4. **Names.** "Brazos Federal Bancshares" + "Mission Plaza Bank" — keep, or pick different Texas-flavored names?
5. **New Lead Auditor identity.** Raj (default — natural promotion from Dawn's drive-in monologue partner) or someone else?
6. **Brazos's CMO and the CAEs.** Should they be named in advance, or invented in the prose during the draft?
7. **Real-vendor naming.** The stack-realism section names real products (Banno, Yello, Insider, Salesforce, Marketo, Q2, Fiserv DNA, Datadog, Splunk). Existing stories name real product *categories* and sometimes specific vendors. Comfort level with this many real-vendor names?
8. **The Steve-by-video moment.** Keep as a 20-minute appearance, or skip it (post-marriage, no recusal frame, the bar for Steve-on-page is lower — easy to include)?
9. **Length target for the full story.** Stories 01 and 04 are 40k-50k tokens. Match, or aim shorter for a "tighter" entry?
