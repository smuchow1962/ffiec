# 🧾 Diary of an Audit Day — Helmstad BioSciences

**Engagement:** FDA BIMO Pre-Inspection Readiness — Phase II NSCLC Program
**Client:** Helmstad BioSciences — mid-size oncology biopharma, Cambridge, MA, ~$1.2B revenue, ~380 employees, two Phase III trials in progress, four Phase II, plus a discovery pipeline
**Status:** TesseraSeal live on the AI clinical-trial-eligibility screening tool for 4 months. Everything else is legacy.
**Date:** Tuesday, the week after Atrio
**Audit team lead:** Karen
**Client liaison:** Dr. Matti Østergaard, VP of Quality and Regulatory Affairs
**Trial in scope:** Phase II second-line therapy in metastatic NSCLC, ~280 enrolled patients across 12 sites in the US and EU

---

### Context

Helmstad BioSciences turned on TesseraSeal four months ago. Not on the QMS. Not on the lab systems. Not on the safety database. They put it on one thing: the AI clinical-trial-eligibility screening tool that decides whether a candidate identified in CRO-supplied data lakes is eligible for the Phase II NSCLC trial. The chain captures every model call, every classification, every confidence score, every reviewer accept-or-reject. The model card hash, the protocol document hash, the de-identified patient hash, the reviewer SSO subject. All of it, sealed.

Everything else at Helmstad runs on the same plumbing every other mid-size oncology sponsor runs on. Veeva Vault for QMS. A Medidata Rave EDC operated by the CRO. Argus for safety. Salesforce-based CTMS for site management. Email and SharePoint for everyday documents. Two CROs supply the data feeds — a global CRO for the Phase III programs and a specialty oncology CRO for the Phase II work. Where the chain doesn't extend, the regulatory record is built on tooling Helmstad does not own end-to-end.

Dr. Østergaard knows this. He has prepped for FDA inspections at three prior companies and has read the 483 letters from a half-dozen others. He picked the AI side first because the FDA's draft guidance on AI/ML in drug development told him to, the eligibility tool was the smallest blast radius that mattered, and four months was a realistic timeline. He wants Karen's team to run the BIMO readiness review like the inspector will run it in six weeks. Find the gaps. Name them. He will take the report into his pre-inspection prep meeting and decide which gaps get closed before the inspector arrives and which get explained.

This is the fourth audit Karen's team has done in four weeks. Last week was a multi-tenant BaaS platform — forty-seven tenants, fourteen hundred verifier runs, zero failures. The week before that was a specialty steel mill in northwest Indiana with three zones and three different posture grades. The week before that was a top-twenty health system that had sealed exactly one inference path and asked the audit to fund the rest. And the week before that was the gold standard.

Today is the third bifurcated audit in four weeks. Different industry. Different inspector. Same shape.

---

### Audit Team

| Name | Role |
|---|---|
| Karen | Lead Auditor — governance and narrative |
| Raj | Database specialist |
| Elena | CRM systems |
| Mike | Application / API layer |
| Diana | IAM and access control |
| Luis | DevOps / logs / pipelines |
| Chen | Data engineering / ETL |
| Tom | Internal-audit liaison specialist (visiting team; partners with the client CAE) |

---

### 🌅 8:30 AM — Kickoff

The drive in was twenty-five minutes from the hotel through Cambridge morning traffic. Karen had her coffee in the cup holder and the engagement brief on her tablet.

*Four bifurcated audits in five weeks*, she thought. *Different industries. Same architecture.*

Northbridge had been the gold standard. TesseraSeal everywhere — every credit decision, every wire, every IAM change, every ETL job. The team had spent four days trying to find a gap and found a stale comment in a YAML file. Karen had driven home that day with the report half-written in her head and the rest of it dictated into her phone.

Mercator had been the bifurcation introduced. AI imaging on the chain. Claims, billing, and the EHR off the chain. Patricia Okonkwo had asked the audit to draw the line cleanly so she could fund the next phase. The CAE had asked Karen to commit to language that would survive the board read. Karen had given him the language and watched him write it down word for word.

Stelvio had been the three-zone version. AI vision QC on the chain. OT on the floor. IT business systems in the back office. Maria Costanza had wanted the report to triage the gaps so the CFO could pick which zone to fund first. The audit had made the case for the OT zone. Maria had taken it to the CFO Friday.

Atrio had been last week. A mid-size BaaS platform with forty-seven fintech-tenant programs running on the same Herald deployment. Multi-tenant chain with per-tenant Merkle roots and a cross-tenant exposure index that the audit had stressed-tested for fourteen hundred verifier runs across two-hundred-and-eighty days of randomly selected dates and not produced a single FAIL. Karen had driven home from Atrio with a quiet confidence she wasn't sure she had felt before. The platform-vendor pattern wasn't theoretical anymore. It worked.

Today was Helmstad. Mid-size oncology biopharma. AI eligibility tool sealed. Everything else legacy. FDA BIMO inspection in six weeks.

*It never is*, Karen thought. *But sometimes the part that's chained is the part that decides, and that's enough — until the FDA inspector asks how the inputs got to the chain.*

She elaborated to herself as she pulled into the visitor lot off Kendall Square. The Mercator pattern was the closest precedent. Different regulator, different stakes, but the same architectural shape: AI side sealed, source side legacy, boundary at the vendor handoff. The difference was that the FDA inspector who would walk Helmstad in six weeks was a Bioresearch Monitoring inspector with a checklist that included 21 CFR Part 11, ALCOA+, ICH E6(R3), and the new AI/ML draft guidance. The HIPAA auditor had been the regulator at Mercator. The FDA inspector at Helmstad would ask sharper questions about source-of-truth. *Plaintiff bar versus FDA inspector*, she thought. *Different motivations, similar reads.*

The lobby was small. Glass-walled conference rooms on the right. A receptionist who recognized Tom's name on the calendar and waved them through.

The conference room was on the fourth floor with a view of the Charles River. Dr. Østergaard was already there with two of his direct reports and a single sheet of paper at each chair. No deck. The sheet was a system map: a green box on the left labeled "TesseraSeal scope: nsclc-phase2 eligibility classifier," a red box on the right labeled "Legacy scope: Veeva QMS, Argus, CTMS, EDC (CRO), CRO data feeds, lab/LIMS, email/SharePoint." A dotted line down the middle. Two arrows crossing the line — one labeled "CRO ingestion," one labeled "EDC extract."

"Good morning," Dr. Østergaard said. He had a soft Danish accent and a careful, deliberate way of speaking. "Before we start. The FDA inspectors arrive in six weeks. I have prepped for inspections at three prior companies. I have seen what happens when the AI side is good and the source side is not. Last year a mid-cap sponsor I will not name received a 483 because their AI eligibility tool was excellent and their CRO data feed was a black box. Two months ago another sponsor — a company called Vellisar Therapeutics, hypothetically — received a 483 because their EDC audit trail was on the CRO's SOC 2 and the inspector wanted Helmstad's word for it, not the CRO's. And six months ago there was a 483 over an Argus audit-trail-table being mutable by DBAs. I have read all three."

He tapped the system map.

"This is what I want from you. The AI side I am reasonably confident about — it has been live for four months and the engineers have been disciplined. The legacy side I want you to look at the way the FDA inspector will look at it. Find what they will find. Tell me what to fix in the next six weeks and what to explain. I would rather hear it from you in May than from the inspector in June."

Karen put her coffee down. "Thank you for being direct. It saves us a day."

Tom — the visiting team's internal-audit liaison — had been on a call with Helmstad's Chief Quality Officer the day before. He nodded. "We agreed yesterday on the bifurcation framing. The CQO is supportive. He wants the report to come out as one assessment with the AI-side and legacy-side findings drawn in separate sections so the inspector can see the boundary."

"Yes," Dr. Østergaard said. "That is what I want. Two sections. One report."

Karen looked around the table at her team. "Okay. Morning is the AI side. Mike and Chen on point — Helmstad's engineers will walk you through the eligibility classifier. Diana, you will do the IAM split — both sides of the line. The line is what we are mapping. Afternoon is legacy. Raj on Argus and the lab database. Elena on the Salesforce CTMS. Luis on the CRO ingestion pipeline. Chen on the EDC extract handoff. We reconvene at three for the reconciliation test. Five-thirty debrief."

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

"This is the entry for decision `nsclc-2026-03-22-mgh-00041`," Dr. Reisch said. "De-identified patient — the hash is HIPAA Safe-Harbor — was assessed against protocol version 4.2 of the NSCLC eligibility criteria. Model returned `eligible` at 0.93 confidence. The CRC reviewed and accepted at 14:18 local. Here." She rotated the screen.

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

Mike leaned back.

"Now pick one from ninety days ago," he said.

Dr. Reisch picked an entry from December 22, 2025 — same site, same protocol version. Ran the verifier. Same four seconds. Same PASS. Same twelve steps.

> **✓ Confirmation #1 — Chain integrity holds at four months and ninety days**
> The eligibility classifier has been emitting sealed entries for 124 days. The verifier resolves a recent entry in 4 seconds and a 90-day-old entry in 4 seconds. Twelve verification steps including HMAC recomputation, Merkle path resolution against the daily-seal Merkle root, and Ed25519 signature verification against the published quarterly public key. The chain endures across the seal boundary.

Mike asked: "What is sealing the seal?"

"Daily Ed25519 signature on AWS CloudHSM in `us-east-1`. We chose CloudHSM after our SOC 2 engagement said an on-prem HSM was overkill for our footprint. The IQ/OQ documentation for the CloudHSM configuration is in Veeva — I can show you if you want. The FDA accepts cloud HSM for Part 11 provided the sponsor's IQ/OQ is documented. Ours is."

"Show me."

She pulled up the IQ/OQ document in Veeva. It was 47 pages. Configuration captures, key-attribute definitions, partition isolation, cross-region failover behavior, key-ceremony minutes from January 2026. Dr. Østergaard's signature on page 47.

> **✓ Confirmation #2 — IQ/OQ for CloudHSM signing infrastructure is documented to Part 11 standard**
> CloudHSM configuration for the daily Ed25519 seal is documented in a 47-page IQ/OQ in the Veeva QMS. Key ceremony minutes are recorded with two-of-three approval (Dr. Østergaard, the platform engineering lead, and the SecOps lead). The document is signed and dated. The configuration is reproducible from the IQ/OQ.

Chen had been quiet, watching the chain entries scroll. He spoke up.

"Show me the model-card binding."

Dr. Reisch pulled up the model registry. Each model version had a SHA-256 of the weights file, a SHA-256 of the model card document, a SHA-256 of the validation report, and the date the version was promoted to production. The chain entry's `model_version` field was a hash of that bundle.

"Every chain entry's `model_version` references a specific tuple of weights, model card, and validation report. If any of the three changes, the hash changes. We can prove which model card was in force the day a decision was made."

> **✓ Confirmation #3 — Model card and validation report are hash-bound to every chain entry**
> The `model_version` field in each chain entry is a SHA-256 over a tuple of (model weights, model card document, validation report). Any of the three changes, the hash changes. The FDA's draft AI/ML guidance asks sponsors to prove which model card governed a specific decision. Helmstad can prove that for any classification by recomputing the hash from the registry.

Mike asked: "And the protocol — the eligibility criteria document?"

"Same pattern. The `criteria_doc_hash` references the protocol version and amendment in force on the decision date. Protocol 4.2 was promoted on January 14, 2026. Every entry from January 14 forward references the 4.2 hash. If we promote 4.3, the hash changes the moment the new criteria are in force."

Mike wrote that down. *That is the ALCOA+ Original attribute on the inputs side.*

> **✓ Confirmation #4 — Protocol document hash binds the eligibility criteria to the decision**
> Every chain entry references the SHA-256 of the protocol document version and amendment in force on the decision date. The classifier cannot decide against a version of the criteria that is not in the registry. The criteria document itself is stored in Veeva with version control. The hash is the bridge.

"Reviewer decisions," Chen said. "Walk me through the human-in-the-loop capture."

Dr. Reisch pulled up a record where the classifier had returned `human-review-required` at 0.62 confidence. The reviewer — a clinical research coordinator at the Memorial Sloan Kettering site — had reviewed the candidate's de-identified profile, made a manual determination of `ineligible`, and entered a reason code from a controlled vocabulary: `prior-systemic-therapy-exclusion`. The reviewer's SSO subject was captured. The decision timestamp was captured. The reason code was captured. The free-text justification — also captured, with a hash, into the chain.

"The reviewer authenticates with SSO," Dr. Reisch said. "The SSO subject becomes part of the chain entry. The reason code is from a controlled vocabulary. The free-text justification is captured. All of it is sealed. We can prove who decided what at what time and why."

> **✓ Confirmation #5 — Human reviewer decisions are sealed end to end**
> Clinical research coordinators authenticate via SSO. Their accept/reject decision, the reason code from the controlled vocabulary, the free-text justification, and the timestamp are all captured into the chain entry. The reviewer's SSO subject is the attribution. The chain shows who decided what, when, and why.

Mike looked at Chen. Chen looked at Mike.

"That is Part 11 audit-trail discipline," Mike said. "On the AI side."

"On the AI side," Chen agreed.

---

### 🧠 10:00 AM — The Pipeline That Feeds the Classifier

Chen had been waiting for this part. The chain at the model boundary is one thing. The chain at the input boundary is the other thing. They are not the same thing.

Helmstad's clinical-informatics platform engineer was a person named Devansh Ramaswamy. Eight years at Helmstad, six of them on data-pipeline engineering. He had a Jupyter notebook open with the pipeline lineage rendered as a DAG.

"Walk me through the input data for the entry Dr. Reisch just pulled," Chen said.

Devansh pulled up the lineage for `nsclc-2026-03-22-mgh-00041`. The candidate's de-identified profile had been ingested at 06:14 UTC that morning from a SFTP delivery from Quintessa Research — Helmstad's global CRO. The SFTP delivery was a tarball of de-identified candidate records. Each tarball was PGP-signed by Quintessa's signing key. The Helmstad ingestion service verified the PGP signature, computed a SHA-256 of the tarball, and recorded the hash as a referenced artifact in a chain entry of type `cro-ingestion`. The candidate records were then unpacked, parsed, and written into the inference-input warehouse. The inference-input warehouse row carried a foreign key to the ingestion chain entry.

"So when the classifier ran on this candidate at 18:14 UTC," Chen said, "the input had been in the warehouse for twelve hours, and the warehouse row was bound by foreign key to a chain entry that references the SHA-256 of the tarball Quintessa sent us this morning."

"Yes."

"Show me the ingestion chain entry."

Devansh did. The entry recorded the SFTP source, the PGP signature verification result, the tarball SHA-256, the file count, the byte count, and the ingestion-service identity. All sealed.

"Good," Chen said. "Now — what about Quintessa's side?"

Devansh's expression changed slightly. Not unhappy. Realistic.

"Quintessa's source-side history is outside our chain. We can prove that the tarball we received was signed by Quintessa's PGP key and that the SHA-256 we recorded matches the tarball. We cannot prove that the contents of the tarball reflect the source EHR records as they existed in the originating site's clinical system. That is Quintessa's responsibility."

Chen wrote that down.

"Their SOC 2 covers the extraction process?"

"Yes. Their SOC 2 Type II report covers the extraction-and-de-identification pipeline. We have the report. It is in Veeva. The SOC 2 covers the period through December 31, 2025. The renewal report is expected in July."

"Have you reviewed the SOC 2 yourself?"

"Quality and Reg Affairs reviewed it. I've read the executive summary."

> **⚠️ Surprise #1 — CRO source-side history is outside the chain of custody**
> The chain captures what Helmstad ingested. It does not capture what was in the source EHR before Quintessa extracted it, or what changes Quintessa applied during extraction and de-identification. The provenance line stops at the SFTP boundary. The Quintessa SOC 2 covers their extraction process; Helmstad relies on the SOC 2 attestation as the integrity statement for upstream of the boundary. The FDA inspector will ask about this.

Chen circled "FDA will ask" twice in his notebook.

"What about the EDC?" he asked. "Medidata Rave. The trial-visit data."

"That is a one-way data feed from the CRO to us. Quintessa operates the EDC for this trial. We get extracts. The EDC's audit trail is the CRO's; we get extracts but the audit trail itself is not extractable in a way that we can hash and bind to the chain."

"So the EDC is a vendor-managed Part 11 audit trail that lives in the CRO's environment."

"Yes."

> **⚠️ Surprise #2 — EDC audit trail is in the CRO's SOC 2 scope, not Helmstad's**
> Medidata Rave is operated by Quintessa for the NSCLC trial. The Part 11 audit trail of trial-visit data lives in Quintessa's environment. Helmstad receives extracts but does not operate the EDC and does not hold the audit trail. Vendor-management dependency. The FDA will ask Helmstad to demonstrate that the CRO's controls meet Part 11. Helmstad can point to the SOC 2 and the contract.

Chen closed his notebook on the pipeline section. *The chain is honest about where it stops. That is the right answer. The question is whether the FDA inspector will accept it.*

---

### 🔐 11:00 AM — Diana on IAM, Both Sides of the Line

Diana started the IAM session by sitting down with Helmstad's identity-platform lead — a person named Rohan Patel — and asking the question she always asked first.

"Show me a credential rotation for the AI eligibility service."

Rohan pulled up the chain entry for the most recent rotation — May 1, 2026, at 03:00 UTC. The eligibility-classifier service account `svc-nsclc-warehouse-reader` had its credentials rotated automatically. The rotation event emitted a sealed chain entry with the old credential fingerprint, the new credential fingerprint, the rotation reason (scheduled), the approver (Dr. Østergaard, with a two-of-three approval), the approval ticket, and the timestamp.

> **✓ Confirmation #6 — AI service-account IAM is chain-coupled end to end**
> Every credential lifecycle event — issuance, rotation, revocation, scope change — for the eligibility-classifier service produces a sealed chain entry. Policy-level changes (who can rotate, who can approve, what the rotation interval is) require two-of-three approval and are themselves chain entries. There is no path to change a service credential without producing a chain record.

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

> **⚠️ Surprise #3 — Argus DBAs have UPDATE permission on the Part 11 audit-trail table**
> Twelve DBAs hold admin access to the Argus database. The Part 11 audit trail is a SQL Server table with no out-of-band integrity binding. The DBAs have UPDATE permission. There is no chain. There is no reconciliation. The application's audit trail is exactly as trustworthy as the discipline of the twelve people who can modify it. The FDA has issued 483s on this exact configuration in the past year.

"And the lab DBAs? The LIMS?"

"Same architecture. Lab DBAs have admin access. The LIMS audit trail is on the same database. No chain. Different team, different group, similar exposure."

> **⚠️ Surprise #4 — LIMS audit trail has the same DBA-mutable shape as Argus**
> The lab information management system audit trail is a database table mutable by the LIMS DBA team. No out-of-band integrity binding. No chain. The lab data feeds into the trial through the EDC and through case report forms. The integrity of the lab audit trail is operator discipline.

Diana made a note in her book that the AI-side IAM finding was clean within its boundary, the legacy-side IAM finding was as bad as Helmstad's peers, and the relationship between the two was the actually-interesting finding — exactly the same shape she had written down at Mercator three weeks ago.

---

### 🧪 12:00 PM — Lunch and the Argument About ALCOA+ "Original"

The team gathered in a side conference room with sandwiches from the building's cafeteria. Tom was on the phone with the CQO. Dr. Østergaard had ducked out to a regulatory affairs huddle.

Karen put her sandwich down before she'd taken a bite.

"Let's talk about the morning."

Mike: "AI side is real. The verifier resolves a 90-day-old entry in four seconds. Model card and protocol document are hash-bound. Reviewer decisions are sealed under SSO. This is the Part 11 audit trail the FDA's draft AI/ML guidance asks for."

Chen: "The pipeline-to-chain reconciliation works at the SFTP boundary. We can prove that the tarball Quintessa sent us is the tarball we processed. We cannot prove what was in Quintessa's source systems before extraction. That is the right boundary to draw, but the FDA inspector is going to push on it."

Diana: "AI service IAM is sealed. CRC IAM is SSO-coupled — the chain trusts the institutional SAML subject. Argus DBAs and LIMS DBAs are not chained. Same legacy posture as the diary baseline."

Raj: "I have the full Argus database review after lunch. Going in expecting what Diana described — DBA-mutable audit table, no integrity binding, no reconciliation."

Karen: "Okay. Question. ALCOA+ has nine attributes. Walk me through the AI side. Attributable — yes, SSO subject. Legible — yes. Contemporaneous — yes, the chain entry is within 200 ms of the decision. Original — that's where I want to argue. Original on the AI side means the chain entry is the original record of what the model said. But the FDA cares about Original on the inputs. The candidate's lab values, the tumor staging, the prior-therapy history. Where does Original live on the inputs side?"

The room got quiet for a beat.

Mike: "The chain is the original record of what the AI saw. It is not the original record of what the EHR said before Quintessa touched it."

Karen: "Exactly. Two different definitions of Original. Original-as-input-to-the-model and Original-as-source-of-truth. The chain handles the first cleanly. The CRO handles the second through their SOC 2 and their PGP signature on the SFTP delivery. The FDA inspector is going to ask which Original we are claiming."

Tom, off the phone: "The CQO is asking me what I think the inspector will weight. I told him my read — the inspector will accept the model-side Original because the chain is convincing. They will push on the source-side Original, especially if there is a discrepancy in the trace."

Chen: "There is going to be a discrepancy. There always is. The question is whether we find it before the inspector does."

Karen tapped her pen on the table. "Tom — tell the CQO that the report will frame ALCOA+ as a per-attribute walk-through with the boundary called out for each attribute. Original gets two paragraphs. One for AI-output Original, which the chain handles. One for source-data Original, which Quintessa's SOC 2 handles. We are honest about the boundary."

Tom relayed it. He came off the call after a minute. "He is good with that framing."

Diana: "Three more attributes I want to flag. Complete — the AI side is complete on what the chain captures, but the chain does not capture pre-screen exclusions that happen in the CRO's pipeline before the candidate ever reaches our classifier. So Complete on the AI side is complete-for-classified-candidates, not complete-for-all-screened-candidates. Consistent — yes for the chain. Enduring — yes, the seal cadence and the daily Ed25519 signature on CloudHSM mean the seal endures across rotations and key generations. Available — yes, Herald.Compliance retrieval works in four seconds for a 90-day-old entry. We confirmed this morning."

Karen: "Good. Write up the ALCOA+ walk for both sides. Boundary called out per attribute."

Luis had been quiet. He looked up from his laptop.

"I was reading their CRO ingestion runbook while you were all talking. The SFTP delivery from Quintessa lands in an S3 bucket. The bucket has versioning enabled and a 7-year retention lock. CloudTrail is enabled. The CloudTrail logs are written to a separate AWS account and the cross-account permissions are configured with an IAM role that the Helmstad SecOps team controls. The cross-account write path looks clean. The S3 bucket itself — the one that holds the Quintessa tarballs — is owned by the Helmstad data-engineering team. They have `s3:PutBucketLogging` permission. Three engineers."

Karen put her sandwich down again.

"Three engineers can disable CloudTrail logging on the ingestion bucket."

"Three."

"And the chain entry references a SHA-256 of the tarball at ingestion time. So even if someone disabled CloudTrail and replaced the tarball, the chain entry would still verify against the recorded SHA-256."

"Yes. The chain catches the tampering at the recorded-hash layer. But the FDA inspector will ask why three engineers can disable CloudTrail at all. That is a separate finding."

Karen wrote it down. "Good catch."

Tom, off the phone again: "The CQO asked one more question. He wants to know — between us — whether this assessment is going to recommend that Helmstad delay the FDA inspection."

Karen looked at him.

"That is his question, not mine," Tom said.

"That is a question for Dr. Østergaard," Karen said. "Not for us. Our job is the assessment."

The team finished lunch. Elena had already wandered off at 12:25 with the CTMS admin's calendar invite on her laptop. The rest of them rinsed coffee cups and walked back out to the engineering floor.

---

### 🔄 1:00 PM — Mike on the API Layer

Mike's afternoon was the API surface. Helmstad runs the eligibility classifier behind an internal API gateway with a custom Lambda authorizer and a request-signing layer. Every call to the classifier — from the CRC review UI, from the batch-screening daily job, from the CRO data-feed ingestion side — goes through this gateway.

Every call emits a chain entry. Request, headers (filtered for PHI — there is no PHI in the request anyway, since the inputs are all de-identified), authorizer decision, downstream service, response code, response hash, latency. Sealed.

"Show me a call that returned a 5xx," Mike said.

The API engineer pulled up an entry from April 8. 503. The authorizer had failed because CloudHSM was briefly unreachable during a planned maintenance window. The chain entry recorded the failure. The downstream classifier had served a fallback `model-unavailable` response. The fallback was itself sealed as a separate inference-shaped entry with `model_id: fallback-noop`, `classification: unavailable`, and a reference to the upstream API failure.

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

Karthik did. The table had 387,000 rows. Each row had a timestamp, an actor, an action, an entity-id, and a before/after JSON. The actor field was a username from the Argus authentication layer. The timestamps were monotonic.

"Has anyone modified rows in this table in the past twelve months?"

Karthik paused. "The DBA team has UPDATE permission. I do not believe anyone has modified rows. I would have to query the SQL Server transaction log to be sure, and the transaction log retention is 30 days."

"So you can confirm no modifications in the past 30 days but not before that."

"Correct."

"And the SQL Server transaction log retention — who controls it?"

"The DBA team. Same group of twelve."

Raj wrote that down. "Who has UPDATE permission on the audit-trail table specifically?"

Karthik queried it. "Twelve DBAs and the application service account. The application writes new rows. The DBAs are technically able to modify."

"Has anyone reviewed who has UPDATE permission on this table in the past year?"

"Not that I know of. It would be in the access-review records that QA maintains."

"Can you pull those?"

Karthik pulled them up. Quarterly access reviews. The most recent was Q1 2026 — March. The review confirmed that twelve DBAs were members of the Argus-DBA group. It did not specifically call out UPDATE-on-audit-trail-table as a privilege, because the group grants admin and admin includes UPDATE on the table by default.

"So the access review confirms that twelve people are admins. It does not specifically attest that those twelve have not modified the audit-trail table."

"Correct."

Raj wrote that down twice.

> **⚠️ Surprise #8 — Argus access review attests group membership, not table-level modification**
> The quarterly access review confirms that twelve DBAs hold the Argus-DBA group, which grants admin. The review does not specifically attest that the audit-trail table has not been modified. The review is shape-checking the IAM, not integrity-checking the audit trail. The FDA inspector who looks at this finding will ask the same question.

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

> **⚠️ Surprise #9 — CTMS Site Monitoring Visit Summary field history is disabled**
> The Long Text field that holds CRA visit notes does not have Salesforce field history enabled. Edits after submission are not recorded. The field is the canonical record of what the CRA observed at the site. The decision to disable field history was made for storage-cost reasons. The HITRUST-style attestation pattern from peer engagements — where this is renewed as a Partial finding cycle after cycle without the remediation being funded — is the same pattern here.

"What about KOL outreach? Same Salesforce instance, separate org?"

"Same instance. Same Salesforce org-within-the-tenant. The commercial team uses a Contact-and-Activity model for KOL relationships. Field history on the Activity comment field is also off, for the same reason."

"PHI in the KOL comments? Or only commercial context?"

"Only commercial context. The commercial team is not supposed to log PHI. We have a DLP rule that scans the field for known PHI patterns and flags it for review. The DLP rule has caught fewer than ten instances in three years."

"Has the DLP rule caught anything in the past quarter?"

"Two cases. Both reviewed by Compliance. Both turned out to be false positives — names of physician contacts that the rule confused with patient names. No PHI in the field."

> **⚠️ Surprise #10 — CTMS-as-KOL-CRM has DLP-only protection on PHI in free-text fields**
> The same Salesforce instance that holds Site Monitoring Visit Reports also holds KOL outreach activities for the commercial team. PHI is policy-prohibited but technically possible. The DLP rule has caught two cases in the past quarter, both false positives. The control is operator discipline plus DLP, not chain-grade. The CTMS is the same instance the diary baseline saw at peer companies.

Elena closed her notebook on the CTMS. "Jordan. Thank you. I appreciate the directness."

"I have been waiting for someone to ask the right questions about Visit Summary. The platform team has been trying to fund the field-history uplift for years. The remediation plan has been the same for four cycles."

"Of course."

She walked back to the conference room. *Same answer. Different industry. Fourth time in five weeks.*

---

### 📊 3:00 PM — The Reconciliation Test

Karen called the team together at 3 PM in the main conference room. Dr. Østergaard was back from his regulatory-affairs huddle. The CQO was on Zoom from the Boston office.

"We are going to pick five eligibility decisions at random and trace them end to end," Karen said. "Backwards from the chain entry to the CRO source. Forwards from the chain entry to enrollment and to the trial-visit data. The AI side is the chained side. Both ends are the legacy side. We are testing the boundary."

Dr. Østergaard nodded. "Pick the five. I will not interfere."

Karen turned to Mike. Mike pulled up the chain database in the conference room's projector and used a deterministic pseudo-random sampler — based on the date, seeded so the choice was reproducible — to select five decisions from the past 60 days. The sampler returned five entry IDs.

```
nsclc-2026-04-02-mgh-00018
nsclc-2026-04-15-mskcc-00074
nsclc-2026-04-21-uchicago-00049
nsclc-2026-04-28-msd-00031
nsclc-2026-05-03-stanford-00027
```

"Run the verifier on all five," Karen said.

Mike ran the verifier on all five. Four seconds each. Five PASS results. Five times twelve verification steps each. Sixty steps. All clean.

> **The AI side: 5 of 5 PASS.**

"Now backwards," Karen said. "For each decision, find the CRO ingestion entry that brought the source data into our warehouse. Then ask Quintessa whether the source records are still on their side."

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

Karen: "Yes. That is a clinical-quality issue, not an audit finding. The audit finding is that the screening pipeline does not have a re-screening loop for source-data corrections in the window between screening and enrollment."

> **⚠️ Surprise #11 — Source-data correction in the screening-to-enrollment window is not re-evaluated**
> The classifier sees the candidate's data at the moment of screening. If the source EHR is corrected between screening and enrollment, the correction is not re-evaluated by the classifier. The CRC's accept decision is bound to the data the classifier saw, not to the corrected data. The patient in question may have been enrolled into a trial they would have been excluded from under the corrected staging. This is an audit finding *and* a clinical-quality follow-up.

Dr. Østergaard wrote it down. "We will pursue this with the medical monitor today. The audit finding is what we asked you to find. The clinical follow-up is mine."

"Forward trace," Karen said. "Of the five decisions, how many trace forward to enrollment?"

Mike pulled up the enrollment records. Four of the five had been enrolled. The fifth — the April 21 University of Chicago decision — had been classified as eligible but the patient had declined to participate. So the forward trace to enrollment was 4 of 5.

"Of the four enrolled, how many trace forward to clinical-visit data in the EDC?"

Mike worked with Marisol to confirm. Two of the four had completed Cycle 1 visits with EDC data captured. One had been enrolled but had not yet had their first visit. One had been enrolled, had Cycle 1 visit data captured, but the EDC extract for that site had been delayed because of a CRO-side processing backlog.

> **The forward trace: 4 of 5 to enrollment. 2 of 5 to clinical-visit data in the EDC.**

Karen wrote the summary on the whiteboard:

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

Dr. Østergaard, after a beat: "That is the picture. The AI side is solid. The boundary at the CRO source is the inspection risk. The forward trace into the EDC is fine when it works but it is not under our chain. Karen — that is the framing for the report."

"Yes. That is the framing."

The CQO, on Zoom: "I want it on record that the AI-side reconciliation is 100% and that the backward and forward traces are mixed. The mixed-ness is not a TesseraSeal failure. It is the boundary."

Tom was writing the language down on his iPad.

---

### 😬 3:45 PM — Friction Between the AI Team and Clinical Operations

Dr. Østergaard had asked the clinical operations lead — a person named Sandra Mendelsson — to join the team for the friction conversation he knew was coming. Sandra had run clinical operations at Helmstad for five years and at two other sponsors before that. She had been nominally supportive of the eligibility-classifier project but had pushed back on extending the chain into the CRO data feed. Her position was that Quintessa would not accept it and that pushing for it would damage the relationship.

She had been listening to the reconciliation results from her office on Zoom and had walked over for the live conversation.

"I want to be honest about something," Sandra said. "The AI team has been disciplined about the chain on the model side. I think the four-month run is a real accomplishment. But the ask to extend the chain into Quintessa is unrealistic. Quintessa is not going to instrument their pipeline for our chain. That is not in our contract. That is not what they sell. We use Quintessa because they are the global Phase III CRO and they handle 14 of our trials. Asking them to instrument is not a conversation that ends in yes."

The AI team lead — Dr. Reisch — pushed back.

"I am not asking Quintessa to instrument. I am asking that the contract require Quintessa to provide hash-bound source records, with their own integrity attestation, that we can store as referenced artifacts in our chain. That is a contract clause, not an engineering integration. The Northbridge Bank pattern from last quarter — the bank required their KYC vendor to deliver hash-bound records. The vendor agreed because the bank made it a contract requirement. We can do the same."

Sandra: "The bank is a customer of the vendor's retail product. We are a customer of the CRO's full-service trial operations. Different relationship. Different leverage."

Dr. Reisch: "I understand the leverage difference. I am still asking that we make the contract change for the next renewal."

Karen looked at Dr. Østergaard.

Dr. Østergaard, calmly: "The right answer is a contract update, not an engineering integration. Sandra is right that Quintessa will not instrument. Hannah is right that a hash-bound delivery clause is a contract change, not an engineering ask. Both can be true. The next contract renewal is in November. We will draft the clause this quarter and put it on the November agenda."

Sandra, quieter: "If we draft it carefully — with their language — they may agree. Hash-bound delivery is not a heavy ask. The integrity attestation language is."

Dr. Østergaard: "We will draft it carefully. Sandra, you and Hannah will lead the drafting. I want a clause we can ship to legal by end of June."

Sandra nodded. Hannah nodded.

Karen wrote it down. *Friction between AI discipline and clinical-operations pragmatism. Mediated by the VP. Resolved as a contract action with named owners and a date. That is what mediation looks like when the principals are senior enough to make the decision.*

The friction did not return for the rest of the day.

---

### 🔍 4:30 PM — The Inspector's Question

Dr. Østergaard had asked Karen the question on the kickoff call three weeks ago. He asked it again now, in the conference room, with Tom and the CQO on Zoom and the team gathered around.

"If the FDA inspector picks a random patient and asks me to demonstrate that the AI eligibility decision was correct — what can I show?"

Karen had been writing the answer in her head all day. She gave it now.

"You can show seven things, in this order.

"One — the chain entry for the eligibility decision. Twelve fields, sealed, with a verifier output that resolves in four seconds and produces twelve verification steps. The inspector can run the verifier on a laptop you provide them.

"Two — the verifier output itself, as a printable artifact. PASS. Twelve steps. Reason string.

"Three — the model-card hash, which resolves to a tuple of (weights, model card document, validation report). The inspector can ask to see the model card document; you produce it from the registry; you recompute the SHA-256 in front of them; the hash matches the chain entry's `model_version` field.

"Four — the protocol document hash. Same drill. The inspector asks to see protocol version 4.2. You produce it from Veeva. You recompute the SHA-256. It matches the `criteria_doc_hash` field.

"Five — the reviewer's accept decision, with the SSO subject, the reason code, the timestamp, and the free-text justification. All sealed. The inspector can ask who the reviewer is — the SSO subject resolves through Okta to the CRC's institutional identity.

"Six — the daily Ed25519 seal from CloudHSM, with the public key and the IQ/OQ document. The inspector can verify the seal independently.

"Seven — the retention attestation. Helmstad's policy is 25-year retention for trial records under Part 11. The chain entries are stored in a managed retention envelope.

"That is a Part 11 defensible evidence pack. Seven artifacts. The inspector can independently verify each one in under twenty minutes.

"If the inspector asks where the source EHR data came from — you point to the SFTP delivery, the PGP signature, the tarball SHA-256, and the Quintessa SOC 2 Type II report. The provenance line stops at the SFTP boundary. Quintessa is the responsible party upstream of that boundary.

"If the inspector asks whether the source EHR was correct in the originating site's clinical system — you say that is outside Helmstad's chain of custody. The originating site is the responsible party for source EHR integrity. Quintessa is the responsible party for extraction integrity. Helmstad is the responsible party from the SFTP boundary forward.

"If the inspector pushes — and they will — on the April 15 patient with the corrected staging that we found this afternoon, you say: we identified that case during a pre-inspection audit. We have notified the medical monitor. The corrective action is in flight. The clinical-quality follow-up is documented. The audit finding has been raised: the screening pipeline does not have a re-screening loop for source corrections in the screening-to-enrollment window. The remediation plan is in the CAPA system.

"If the inspector asks what about the EDC — you say the EDC is operated by Quintessa under their SOC 2. You provide the SOC 2 report. You explain the vendor-management dependency.

"If the inspector asks what about Argus — you say the Argus audit trail is mutable by DBAs and that you have a CAPA in flight to add an integrity binding. You explain the timeline. You do not pretend the gap is not there.

"That is the inspector-facing posture."

Dr. Østergaard had been listening with both hands flat on the table. He nodded once when Karen finished.

"That is the posture," he said. "Tom — please make sure that articulation is in the report verbatim."

Tom was already typing it.

The CQO on Zoom: "Karen — that articulation. Use exactly that language. The seven-artifact evidence pack is the inspection-day playbook."

Karen: "Use it."

---

### 🌆 5:30 PM — Debrief

The conference room was quieter than it had been at kickoff. Dr. Østergaard had ordered coffee and pastries. The CQO was on Zoom from his car — he had a hard stop at 6:30. Sandra was on Zoom from her office. Hannah and Devansh were in the room.

Karen stood up at the whiteboard.

"Two findings sections," she said. "AI side. Legacy side."

She wrote on the whiteboard.

#### ✅ What They Confirmed on the AI Side

| # | Finding |
|---|---|
| 1 | Chain integrity holds at 4 months and at 90 days. Verifier resolves any entry in ~4 seconds in 12 steps. |
| 2 | CloudHSM signing infrastructure is documented to Part 11 standard. IQ/OQ in Veeva, signed and dated. |
| 3 | Model card and validation report are hash-bound to every entry. Provable per-decision model lineage. |
| 4 | Protocol document hash binds eligibility criteria to the decision. ALCOA+ Original on the criteria side. |
| 5 | Reviewer decisions are sealed end to end under SSO subject. Reason code, timestamp, justification all captured. |
| 6 | AI service-account IAM is chain-coupled. Two-of-three approval on policy changes. |
| 7 | CRC reviewer decisions are sealed under SSO subject attribution. |
| 8 | API gateway and fallback paths are sealed. CloudHSM maintenance windows are visible. No silent gaps. |

**AI side: 0 Gaps. 0 Partials. ALCOA+ defensible end to end on the surface area the chain covers.**

#### ❌ What They Found on the Legacy Side

| # | Finding |
|---|---|
| 1 | Argus DBAs have UPDATE permission on the Part 11 audit-trail table. No out-of-band integrity binding. |
| 2 | LIMS audit trail has the same DBA-mutable shape as Argus. |
| 3 | Veeva Vault QMS has a "Late Effective Date" workflow that allows up to 30 days of effective-date back-dating without showing in the standard audit-trail report unless the user knows the filter. |
| 4 | CTMS Site Monitoring Visit Summary field history is disabled. Salesforce CRM-pattern. |
| 5 | EDC extract pipeline has no chain. Trial-visit warehouse is loaded from CRO SFTP without HMAC binding. |
| 6 | Argus-to-warehouse ETL has no chain. Quarterly reconciliation is a count check, not a contents check. |
| 7 | CRO source-side history is outside the chain. Provenance stops at the SFTP boundary. |
| 8 | EDC audit trail is in the CRO's SOC 2 scope, not Helmstad's. Vendor-management dependency. |
| 9 | Source-data correction in the screening-to-enrollment window is not re-evaluated by the classifier. April 15 patient. Clinical follow-up in flight. |
| 10 | Three engineers can disable CloudTrail logging on the CRO ingestion S3 bucket. |
| 11 | CTMS-as-KOL-CRM has DLP-only protection on PHI-prohibited free-text fields. |
| 12 | Quarterly access review attests group membership, not table-level modification. |
| 13 | Quintessa's retention policy removed source records for 2 of 5 reconciliation samples (90-day rotation). |

**Legacy side: 4 Gaps. 6 Partials.** (The Partials are items where there is *some* protection — Veeva's audit trail is real if you know the filter; Salesforce field history is on for some fields; the SFTP signature does verify upstream provenance to a point — but the protection is at the discretion of the system operator or the vendor and is therefore not chain-grade evidence. The remaining items are Gaps in the strict sense.)

Karen ticked off the Gaps on her fingers as she summarized.

"Argus DBA write access is a Gap. Veeva late-effective-date workflow is a Gap. EDC extract has no chain — Gap. CTMS overwrites — Gap.

"The Partials: CRO source-side history is a Partial — there is the SOC 2, there is the SFTP signature, there is the tarball hash, but the source-side itself is outside our perimeter. Lab values discrepancy investigation — Partial pending the medical monitor's follow-up. Quintessa rotation policy is a Partial because their 90-day retention is contractually permitted but limits backward traceability. IAM bifurcation for non-AI users is a Partial — the legacy AD is not chain-coupled but the access reviews exist. Retention policy variance between Helmstad's 25-year requirement and Quintessa's 90-day staging — Partial. Audit-trail filter discoverability in Veeva — Partial because the filter exists but is not the default."

#### 🔁 Side-by-side comparison

| Dimension | AI side (eligibility classifier) | Legacy side (Veeva, Argus, CTMS, EDC, CROs, LIMS) |
|---|---|---|
| Record integrity | Sealed via HMAC chain + daily Ed25519 seal on CloudHSM | DBA-mutable in Argus and LIMS; Veeva back-datable; CTMS overwrite-able |
| Identity coupling | Service IAM chain-coupled; CRC SSO chain-coupled | Argus and LIMS DBAs not chained; commercial CTMS users not chained |
| Reconciliation | Cross-checkpoint sealed events; verifier in 4 sec | Quarterly count checks in Argus; manual reconciliation elsewhere |
| Retention | Helmstad 25-year, sealed | Quintessa 90-day staging; SQL Server 30-day txlog |
| Verifier | `herald-verify` resolves in 12 steps | No verifier; audit by inspection |
| FDA Part 11 audit-trail | Defensible per draft AI/ML guidance | Defensible only on shape, not on integrity |
| ALCOA+ Attributable | Yes (chained SSO subject) | Yes (Argus authentication; CTMS Salesforce subject) |
| ALCOA+ Original (model output) | Yes (chain entry is the original record) | n/a |
| ALCOA+ Original (source data) | Boundary at CRO SFTP | Quintessa SOC 2 |
| ALCOA+ Contemporaneous | Yes (200 ms chain latency) | Mostly yes; CTMS allows post-hoc edits |
| BIMO inspection posture | Strong — 7-artifact evidence pack | Mixed — vendor-management dependencies |

Karen paused on the table for a beat longer than the rest.

"The chain proves what the model said. The chain proves what the criteria document was. The chain proves who the reviewer was. The chain proves the reviewer's decision and reason. That is the FDA's question, on the AI side.

"The chain does not prove what the source EHR said before Quintessa touched it. The chain does not prove what the EDC's audit trail says about a trial visit. The chain does not prove that an Argus row was not modified by a DBA with UPDATE permission.

"That is the inspection risk. Six weeks. Argus first. EDC contract clause for November renewal. Veeva late-effective-date filter as a CAPA. CTMS field-history funding for the next budget cycle. The April 15 patient is medical-monitor work and is in flight as of this afternoon."

The CQO on Zoom: "That is the line for the inspection-prep memo. Karen — Tom — thank you. The seven-artifact evidence pack is going to be how we run the first morning of the inspection. The legacy-side findings are going to be the second morning. We will know what is coming."

Dr. Østergaard: "Permission to share the bifurcation framing in the BIMO prep documents to the executive committee on Friday."

Karen: "Permission granted. We will send the formal report by end of week. The bifurcation framing is the framing. Use it."

The CQO on Zoom: "One more thing. I want it on record that this assessment found the AI-side controls to meet or exceed best-known practices observed in deployed AI clinical-decision-support systems audited in the last twelve months. And that the legacy-side controls are at-or-below the median for mid-size oncology sponsors. That is the comparative posture we are taking to the executive committee. Karen — you good with that characterization?"

Karen took a beat. "I am, with the same caveat I gave at Mercator. The population of sponsors running chain-grade AI controls in production is small. Helmstad is in the top quartile of a fairly small group. The 'meets or exceeds best-known practices' phrasing is more accurate. I would also add — for the executive committee read — that the legacy-side gaps include items where the gap is known and the remediation has been deferred for budget or contractual reasons. The CTMS field-history-disabled finding has been on the HITRUST-equivalent attestation for prior audits. The Argus DBA finding has been a known industry concern for years. Helmstad is not surprised by these findings. The committee should not be surprised either."

The CQO: "Better. Use that."

Sandra, on Zoom: "Agreed."

Tom was writing the language down verbatim.

Hannah: "One more for the record. The April 15 patient. I want it documented that the audit team found the discrepancy during a pre-inspection readiness audit, not during the inspection. The fact that we found it before the inspector matters for the BIMO read."

Karen: "Documented. The reconciliation methodology — five-decision random sample with backward and forward trace — is in the report. The April 15 finding is one of the five. It will read as a discrepancy detected by the audit, with corrective action in flight."

Dr. Østergaard stood up, shook Karen's hand, then Tom's, then went around the room and shook each team member's hand individually. It took ninety seconds. He thanked each of them by name.

"Drive safely. The Massachusetts Pike is bad in the rain."

Karen looked out the window. It had started raining at some point during the afternoon.

The team packed up. There was the usual quiet shuffle of laptops closing and notebooks going into bags. Diana paused on the way out and asked Rohan a private question about extending the AD access-review process to attest table-level privileges, not just group membership. Rohan answered it. Mike traded contact details with Hannah, who had stayed for the debrief tail. Chen and Devansh exchanged GitHub handles. Luis caught Karthik in the hallway and gave him an unsolicited recommendation about a SQL Server Always Encrypted pattern that would not have prevented DBA write access but would have made it visible faster. Karthik thanked him. Elena thanked Jordan for being honest about the field-history-disabled situation.

Sandra was the last one out. She stopped Karen by the door.

"I want to say — I was wrong about the Quintessa contract clause. The audit gave me the framing I needed. The November renewal is going to be a conversation, not a fight. Hannah and I will get the clause to legal in June."

Karen: "Good. The contract change is the right move. The chain into the CRO is not."

Sandra: "Understood."

Karen watched the room empty. She gathered her notes, clipped her pen back to the cover, and slung her bag over her shoulder.

Tom held the door for her on the way out. "Round four is in the rear-view."

"Round four is the second bifurcation in five weeks," Karen said. "Mercator was the first. Helmstad is the same shape with a different inspector. The pattern is real."

"Different inspector. Same architecture."

"Same architecture. That is the line."

---

### 🧾 Final Assessment Theme

The drive back to the hotel was twelve minutes through Cambridge rain. Karen had her coffee, refilled, in the cup holder. The wipers were on intermittent.

She thought about the day. About the bifurcation. About Hannah's terminal showing a clean verifier resolution at 9:45 in the morning. About Luis pointing at the S3 bucket policy at lunch and saying "three engineers can disable CloudTrail." About the April 15 patient and the staging correction that had moved between screening and enrollment. About Dr. Østergaard's question at 4:30 and the seven-artifact answer she had given.

She thought about Northbridge, four weeks ago, where everything had been sealed and the audit had been almost boring.

She thought about Mercator, three weeks ago, where the chain had been on the imaging path and the EHR had been mutable and Patricia had asked the audit to draw the line so she could fund the rest.

She thought about Stelvio, two weeks ago, where the chain had been on the QC vision and the OT had been on PLCs and Maria had asked the audit to triage three zones.

She thought about Atrio, last week, where forty-seven tenants had run on one chain and the verifier had not failed once across fourteen hundred runs.

Helmstad was the second bifurcated audit. AI side sealed. Legacy side legacy. Different industry. Different regulator. Same architectural pattern. Same line drawn between the part that decides and the part that supplies the inputs.

The thing that made Helmstad interesting was not that the chain held. The chain has held in every audit since Northbridge. The thing that made Helmstad interesting was that the FDA inspection in six weeks would test the *boundary* — and the boundary was where the next class of audit findings would live, in this and every other bifurcated deployment for the rest of Karen's career.

The chain proves what the model decided.
The contract has to prove what the model was given.

Karen picked up her phone at a red light and dictated a one-line note for the report's executive summary.

> *"On the screening path, Helmstad can demonstrate what the classifier said and that the record was not tampered with. On every other path, integrity is still vendor-managed."*

The light turned green. She put the phone down and drove the rest of the way to the hotel through the rain.

---
