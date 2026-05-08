# 🧾 Diary of an Audit Day — Olmstead University

**Engagement:** Multi-framework audit-readiness assessment (FERPA review, NIH research-integrity audit, biennial GLBA audit, upcoming HHS OCR audit at the affiliated medical center)
**Client:** Olmstead University — private R1 research university, upper Midwest. ~28,000 students. ~$1.4B endowment. ~$680M annual research expenditure. Affiliated teaching hospital (Olmstead Medical Center) under shared board, separate compliance office.
**Posture:** TesseraSeal deployed 11 months ago on a single use case — undergraduate admissions AI screening — under a consent-to-resolve framework with a civil-rights firm. Everything else legacy: research computing, Olmstead Medical Center IT, financial aid (GLBA), advancement CRM, IRB system.
**Date:** Wednesday, the day after Pacific Crescent wrapped
**Auditor:** the same eight-person team that walked the diary baseline, Mercator the week after, Northbridge the week after that, then Stelvio, Atrio, Helmstad, and Pacific Crescent

---

## Context

Olmstead is a private R1 with fourteen schools and colleges. Top-25 medical school. Top-50 law school. Top-50 business school. Substantial engineering school. About fourteen thousand undergraduates and fourteen thousand graduate and professional students. The endowment runs about $1.4B. The annual research take is around $680M — NIH, NSF, DoD, foundation, and industry money in roughly that order.

Olmstead Medical Center sits on the same campus, shares the board, and is a separate legal entity with a separate compliance office. The two organizations share data on the medical-school side — clinical and research data both — and that data-sharing seam is its own audit problem.

Eleven months ago, an applicant-class disparate-impact threat letter from a civil-rights firm landed on the General Counsel's desk. The letter named the undergraduate admissions AI screening system as the source of the alleged disparity. Olmstead's response — written into the consent-to-resolve framework with the firm — was to put the AI screening system under TesseraSeal: every model score, every reviewer override, every retraining event, every fairness-audit report linked by hash. The chain is the university's principal defense if the firm files suit.

Everything else at Olmstead is legacy.

The research-computing side is faculty-led federalism. The HPC cluster is centrally managed; the labs are not. Each PI runs her own lab. Central IT can advise. Central IT cannot enforce. NIH and DoD have started asking questions that central IT does not have the leverage to answer.

The medical center is a separate audit problem. Epic EHR. MyChart. The medical school's research databases overlap with the hospital's clinical systems through documented data-sharing agreements. HIPAA-covered. No chain. The same Epic-side findings the team wrote up at Mercator three weeks ago.

Financial aid is GLBA-covered under the FTC Safeguards Rule. Banner. Mutable. The annual administrative-access review is paper-based.

Advancement uses Salesforce for donor relationship management. Major-gift cultivation notes overwriteable. Same shape as the diary baseline.

The IRB system is a homegrown SQL-backed app on the medical-school side. Protocol amendments are versioned. The audit trail for IRB approvals is in a database that admins can `UPDATE` directly.

Karen's team was engaged by Dr. Ines Achterberg, Vice Provost for Research Integrity and Compliance. PhD in epidemiology. Sixteen years in higher-ed compliance after a stint at NIH's Office of Research Integrity. The deliverable will be read by the General Counsel, the medical-center compliance office, and the Faculty Senate's research-integrity committee.

Four regulators are watching at once: the Department of Education on the FERPA side (a current OCR complaint from an admitted student about how her data was used to score her), NIH on the research-integrity side, the FTC on the GLBA side, and HHS OCR on the medical-center HIPAA side. Each wants something different. The chain on the admissions side fits one of them.

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
- **Tom** — Internal-audit liaison specialist (visiting team — partners with the client CAE)

Client-side liaison: **Dr. Ines Achterberg**, Vice Provost for Research Integrity and Compliance, Olmstead University. Direct. Politically aware. Knows where the bodies are buried in faculty-led labs without saying so.

---

## 🌅 8:30 AM — Kickoff and the Drive In

Karen rode in with Tom from the hotel. Twenty minutes north along the lakeshore. The campus came into view through a screen of bare oaks — Romanesque limestone, modern glass on the engineering-quad side, the medical-center towers a half-mile off behind the law school.

Tom had a thermos. Karen had her usual — black coffee in a travel cup, half gone before they pulled out of the hotel lot.

"Recap me," Tom said. "Just the headline of each one."

"Northbridge — banking, gold standard. Chain across the whole institution. We ran out of things to find by 3 PM."

"Mercator."

"Healthcare. AI imaging chained, claims side mutable. Bifurcated. We wrote the seam down the middle of the report and the CMO understood it instantly."

"Stelvio."

"Manufacturing. AI side chained, OT side legacy, IT business side legacy. Three zones. Maria took the four-second verifier video to her CFO."

"Atrio."

"BaaS multi-tenant. The chain shape was different but the seam was in the same place — between the regulated function and everything around it."

"Helmstad."

"Biopharma. AI eligibility decisioning chained, the rest of the clinical-trial stack mutable. Same shape again, different industry."

"Pacific Crescent."

"Utility. Yesterday. Public-safety stakes. AI on outage prediction chained, OT on the substations and the SCADA legacy. The seam mattered because the regulator there cares about what the AI told the operator and what the operator did with it."

Tom waited. Karen looked at the limestone tower coming up on the left.

"And today?"

"Today is the higher-education version. AI on a contested decision. Everything else loosely governed by faculty federalism. Different industry, same shape we keep finding."

"What's the wrinkle?"

"Two wrinkles. One — the AI side wasn't deployed because of best-practice. It was deployed because a civil-rights firm sent a threat letter eleven months ago and the chain is the consent-to-resolve framework. That changes the political calculus on every other system in the institution. The chain works because legal demanded it. Nothing else has had legal demand it."

"And two?"

"The medical center is on the same campus, shares the board, and has its own compliance office. We have no jurisdiction there. Ines wants a hallway tour we can write a recommendation off. She is honest about the split."

"What's the recurring line you keep saying?"

Karen drained her cup. "It never is."

"That one."

"It never is. Sometimes the chain is on the part that's being sued, and that's the only part that needs to be."

They pulled into the visitor lot at 8:24. The morning bell tolled across the quad. Students were walking to 8:30 classes in down jackets.

Ines met them at the badge desk in the administration building. Sweater, slacks, ID lanyard, the kind of handshake that came from sixteen years of compliance briefings she had to deliver before noon.

"Welcome. Coffee is in the conference room. We have a tight day. I would like to start with the admissions AI dashboard, do the database work, IAM, lunch, then a research-computing tour with two lab visits, then reconciliation. The medical-center hallway tour is at four o'clock — I will walk you over personally for thirty minutes."

Karen smiled. "That's a clean agenda."

"It needs to be. The OCR complaint review is in three weeks. The General Counsel reads my Friday memo. The Faculty Senate research-integrity committee meets next Tuesday. I need to know what is defensible and where the exposure is."

She paused.

"I am going to be straight with you. The admissions chain works. I built that program. The rest of the institution is not under the same regime and will not be for years. I want you to write what is true. I will read it Friday."

Tom nodded. "Same kickoff Maria gave us at Stelvio."

Ines did not smile. "Maria and I were on a panel together in October. We compared notes."

> **🔍 Karen's note (internal):**
> *Two clients. Same disposition. They ran into each other on a compliance-conference panel and have been comparing notes for six months. The shape we keep finding has names attached to it now. That changes the politics of the report.*

The team filed into the conference room. Ines flipped on the wall display.

---

## 🧩 9:15 AM — The Admissions AI Dashboard

The dashboard lived inside the admissions office's intranet. Ines logged in with her university SSO and re-authed with a hardware key. The dashboard showed four panes — a count of applications scored to date for the current cycle, a fairness-monitoring panel with disparity ratios across protected classes, a reviewer-override panel showing the top reviewers by override volume, and a model-version tracker.

Mike pulled up his laptop. "Pick a decision from this cycle for me."

Ines selected an application from January 9, 2026. The application had been scored at 71. The reviewer had overridden up to 78 with the structured reason code `STRENGTH_OF_RECOMMENDATIONS`. The committee had admitted.

Mike copied the entry ID into his terminal:

```
herald-verify --tenant=olmstead \
              --service=undergrad-admissions-screen \
              --date=2026-01-09 \
              --entry-id=2026-01-09-UA-71418
```

Five seconds. The terminal returned:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key admissions-2026-q1
```

Mike turned the laptop. Karen read the output. Ines watched without comment.

"Good." Karen looked at the dashboard again. "Pick one from before the chain went live."

Ines paused. "Eleven months ago — March 2025."

"Pick one from February 2025."

She picked one. The application had been scored, reviewed, and the applicant had been waitlisted.

Mike copied the ID. Ran the verifier.

```
Status: REJECTED
Step: 1
Reason: entry-id not found in chain — chain effective
        date for tenant=olmstead service=undergrad-admissions-screen
        is 2025-05-12
```

Mike nodded. "Pre-deployment. Expected. Documented in the institution's CC8.1 as the chain's effective start date."

Ines: "May 12, 2025. The day the consent-to-resolve framework was finalized. The chain went live the same afternoon."

Karen wrote: *Chain effective date 2025-05-12. Pre-chain decisions are out of scope for chain-based defense. The threat letter referenced the prior cycle's decisions — those defenses live in the consent-to-resolve framework directly, not in the chain.*

> **✓ Confirmation #1**
> The undergraduate admissions screening chain is live and producing verifiable entries within the current cycle. Mike re-verified a January decision in five seconds. The chain effective date is documented as the day the consent-to-resolve framework was finalized. Pre-chain decisions are explicitly out of scope and the verifier rejects them by design rather than fabricating a result.

Ines walked them through the chain payload. Applicant hash — FERPA-de-identified at the screening boundary, the cleartext applicant ID lives in Slate. Model ID. Model version. Score zero to a hundred. Confidence. Top five SHAP feature attributions. Reviewer ID. Override decision — accept-AI, override-up, override-down, refer-committee. Override reason code from a controlled vocabulary.

Mike ran his finger across the screen. "Where is the free-text override rationale?"

Ines's face shifted half a degree. "Slate."

"Slate the admissions CRM."

"Slate. The reviewer types the rationale into a free-text field in the applicant record. The field has a 30-day edit history in Slate. Past 30 days, prior versions are not retained."

Mike's pen paused over his notebook. "And that field is overwriteable indefinitely."

"Yes."

"The chain has the override decision and the reason code. The chain does not have the rationale free-text."

"Correct."

Karen wrote: *Chain captures the structured override decision. Free-text rationale lives in Slate. 30-day field history. Indefinitely overwriteable. Watch this — this will be the legally interesting gap when the team gets to reconciliation at three.*

> **⚠️ Surprise #1**
> The override-rationale free-text — the reviewer's actual reasoning for changing the AI score — lives in Slate, not in the chain. Slate retains 30 days of field history; past 30 days, prior versions are not preserved. The chain captures the structured override decision and the reason code. It does not capture what the reviewer wrote into the free-text field.

Ines did not push back. "I have flagged this internally. The Slate integration was scoped to webhook the structured fields only. The free-text was deferred."

"Deferred to when?"

"Phase 2."

"Funded?"

"Not yet."

Karen wrote: *Same structure as Stelvio. Phase 2 unfunded. The thing that needs the chain is the thing that does not have it.*

The team split.

---

## 🧠 10:00 AM — Database Deep Dive

Raj had the corner of the conference table and three screens. The admissions ledger on one. Banner financial-aid on another. The IRB SQL backend on the third. The advancement CRM Salesforce backend was on a fourth tab in a browser window.

Four mutability profiles. Four audiences. He worked them in order.

### The admissions ledger

Append-only by design. The Herald.Py SDK signs each entry. HMAC-SHA-256 chain links entry N to entry N-1 with HKDF-per-tenant key binding. Daily Ed25519 seals on AWS CloudHSM in `us-east-2`.

Raj picked a random entry from three weeks ago. Verified.

```
Status: PASS
Step: 12
```

He picked the very first entry from May 12, 2025. Verified. PASS.

He attempted a direct UPDATE on the chain table. The database accepted it because nothing at the database layer prevents it. He ran the verifier on the next entry.

```
Status: FAIL
Step: 4
Reason: HMAC mismatch — entry payload does not produce
        the chained HMAC recorded in entry N+1
```

He attempted a multi-entry rewrite. The verifier failed at the daily seal:

```
Status: FAIL
Step: 9
Reason: Merkle root mismatch — recomputed root does not
        match sealed root for date 2025-11-18
```

Raj rolled back the mutations. The chain returned to PASS. He noted the test in his workbook.

> **✓ Confirmation #2**
> The admissions-side ledger is append-only in practice. Direct database mutation is technically possible — the Postgres backend is mutable like any Postgres backend — but the verifier catches single-entry tamper at the HMAC layer and multi-entry tamper at the Merkle/seal layer. The Ed25519 daily seal is on AWS CloudHSM in `us-east-2`. The database engineers do not have HSM access. The split is real.

### The Banner financial-aid backend

Raj opened Banner. SQL Server. He had read-only access through Ines's audit role.

He pulled the schema for `SFAREGS` and `RPRAWRD`. Award records. He ran a query against the audit-log configuration.

Audit logging was enabled. Retention was set to 90 days.

He picked a random award record from the current academic year. Looked at the change history.

The history showed the award amount edited twice in the last quarter. Two different counselor IDs. The reason codes were both `ADMINISTRATIVE_ADJUSTMENT`. No free-text justification.

"And the original award amount?"

Ines pulled up the record. "The chain of edits is visible. The original is the row at the top of the change-history list."

"For records edited within the 90-day window."

"Yes."

"Records edited before the 90-day window — the prior values are gone."

"Correct."

Raj wrote: *Banner audit-log retention 90 days. Original values older than 90 days not preserved. GLBA Safeguards Rule 2023 update expects 'audit and oversight' of administrative access — the audit configuration is on, the retention is short.*

> **⚠️ Surprise #2**
> The Banner financial-aid audit log is configured but retention is 90 days. Award amounts can be edited; the change history is preserved within the 90-day window and overwritten after. For an academic year that runs longer than 90 days, edits to a fall award visible in early spring will have lost their original values by the time the GLBA biennial audit runs.

### The IRB SQL backend

Raj opened the IRB system. Homegrown app on the medical-school side. SQL Server backend.

The schema had `irb_protocol`, `irb_amendment`, `irb_approval`, and `irb_audit_log`. Protocol amendments were versioned — each amendment was a row, with a reference to the prior amendment. The approval table was different. Each approval was a single row with a status field — `pending`, `approved`, `expired` — and a `last_modified_by` and `last_modified_at`.

Raj asked Ines, "Who can edit the approval table directly?"

She paused. "The medical-center compliance office maintains this system. I would have to ask them."

"Best guess?"

"The system administrator. There is one named system administrator. There is also a service account that the application uses. I do not know if there are other accounts with write access."

"Have any IRB-approval records been edited in the past year?"

"I do not know."

Raj waited.

Ines looked at the screen. "I will ask. I will not have the answer today."

He wrote: *IRB approval table — single-row, last-modified-by/at columns, mutable. Compliance office on medical-center side maintains. Cross-organization. Out of Ines's direct reach. Document and recommend the question be asked.*

> **⚠️ Surprise #3**
> The IRB approval audit trail is in a homegrown SQL database on the medical-school side. The approval table is a single row per protocol with `last_modified_by` and `last_modified_at` columns — mutable by anyone with database write access. The medical-center compliance office maintains the system. Ines does not know whether IRB-approval records have been edited in the past year and cannot answer today.

### The advancement CRM (Salesforce)

Raj opened the Salesforce backend. Opportunity records for major-gift cultivation. He pulled a record at random.

Field history was enabled on the financial fields — committed amount, payment schedule, fund designation. Field history was not enabled on the cultivation-notes field — the long free-text where development officers documented their conversations with prospects.

Raj asked Elena to pick this one up. She nodded.

> **⚠️ Surprise #4**
> The advancement CRM in Salesforce has field history enabled on financial fields but not on the cultivation-notes free-text field. Development-officer conversation notes are overwriteable indefinitely with no version history. Same pattern as the diary baseline two months ago.

Raj closed his three screens. He stacked the workbook page.

Four backends. Four mutability profiles. Four audiences:

- Admissions ledger — immutable in practice. FERPA's audience.
- Banner — auditable for 90 days, mutable past that. GLBA's audience.
- IRB approval — mutable, cross-organization governance. NIH's audience and the IRB itself.
- Advancement CRM — mutable, no version history on the field that matters. Donor-relations governance, not regulatory.

He wrote at the bottom of the page: *One chain. Three legacy backends. Each backend speaks to a different regulator. The chain answers one question. The other three do not.*

---

## 🔐 11:00 AM — IAM Review (Four Columns)

Diana had a workbook with four columns this time, not three. Admissions AI. Banner. Slate. Salesforce advancement. She filled them in order.

### Admissions AI side

Every credential the admissions screening service used — database creds for the chain backing store, S3 creds for the model artifacts and fairness-audit reports, the AWS CloudHSM PIN for the daily seal — every rotation was a chain entry. `event.type = credential.rotated`, with rotator identity, rotation reason, and the new key fingerprint.

Diana picked the last six rotations. Verified each. All PASS.

She ran:

```
herald-verify --tenant=olmstead \
              --service=undergrad-admissions-screen \
              --event-type=credential.rotated \
              --date-range=2025-05-12:2026-04-09
```

Eleven entries returned across the eleven months. Quarterly rotations on database creds. Two ad-hoc rotations on S3 keys after personnel changes. All PASS.

> **✓ Confirmation #3**
> Credential rotation on the admissions screening service is captured in the chain with rotator identity, fingerprint, and reason. Eleven rotations across eleven months. All verifier PASS. Multi-factor authentication on every service account that accesses the chain. The IAM under the AI service is chain-coupled.

### Banner financial-aid IAM

Diana asked Ines to log in to the Banner counselor workstation that financial-aid staff used during peak FAFSA season.

The login was through the university SSO. Diana checked the user roster. There were 23 financial-aid counselors with individual accounts. There was also one shared account: `aid_admin`.

"Who uses `aid_admin`?"

Ines: "During peak season — January through April — counselors use it for batch operations. Loading FAFSA imports, generating award letters, running aid-package recalculations."

"How many people know the password?"

"Six? Seven? I would have to ask the financial-aid office."

"MFA?"

"On the SSO, yes. On `aid_admin`, no — it bypasses SSO because it is a service account for batch jobs. Counselors keystroke-paste the password from a shared password manager."

"Last password rotation?"

"I do not know offhand."

Diana wrote: *Shared account during peak season. SSO bypass. Six or seven counselors keystroke-pasting from a shared password manager. Same shape as the Stelvio `Plant_Engineer` shared account. Different industry, same compounding risk.*

> **⚠️ Surprise #5**
> The Banner financial-aid system has a shared `aid_admin` account used during peak FAFSA season for batch operations. The account bypasses university SSO. Six or seven counselors share the password through a shared password manager. The last rotation date is not known offhand. GLBA Safeguards Rule 2023 update requires "audit and oversight" of administrative access — the audit log records actions taken under `aid_admin`, but the account binding to a specific human is operational discipline, not technical enforcement.

### Slate (admissions CRM)

Slate was per-counselor SSO with MFA. Diana checked.

Then she found the bypass. The admissions office had documented an SSO-bypass for "admissions readers traveling overseas during the holiday application-reading period" — counselors stationed in Europe and Asia who reviewed early applications. The bypass was a static-token authentication that did not require MFA.

"How many counselors used the bypass last year?"

Ines pulled up the access log. Eleven.

"Is the static token rotated?"

"Annually."

"And the bypass is documented?"

"Yes. There's a memo. Approved by the Director of Admissions and the CISO."

Diana wrote: *Bypass is documented. Bypass exists. Token rotation is annual. The bypass weakens MFA enforcement on an audience that includes the admissions readers who are doing the override decisions in the chain. Document.*

> **⚠️ Surprise #6**
> Slate has a documented SSO-bypass for admissions readers traveling overseas. Static-token authentication. No MFA on the bypass path. Token rotated annually. Eleven counselors used the bypass last year. The bypass is documented and approved, but it weakens the IAM posture on the audience whose override decisions are the chained artifact.

### Salesforce advancement IAM

Per-user accounts with SSO and MFA. Field-level security enforced. The development officers' access matched their portfolio assignments.

The IAM on the advancement side was the cleanest of the four. The data being protected was the weakness — the cultivation notes field had no field history, so the IAM controls protected access to a record whose contents could be silently overwritten.

Diana wrote: *Salesforce IAM clean. The mutability of the protected data is the gap, not the access control.*

She stacked the four columns.

Admissions AI — chain-coupled rotation history.
Banner — shared account during peak season, SSO bypass on the shared account.
Slate — per-user with documented bypass.
Salesforce — clean IAM, weak data integrity.

Four columns. Four shapes. The AI column had the strongest IAM, because the chain forced the cleanup. The others had whatever they had when nobody forced them to look.

She wrote at the bottom of the page: *Same finding as Stelvio. Where the chain is wired in, IAM behaves. Where it is not, IAM is whatever the operations team can sustain without enforcement.*

---

## 🧪 12:00 PM — Lunch (Reporting Frame)

The catering was in the small dining room across from the conference room. Soup and sandwiches. Coffee that was actually warm.

Karen and Tom took a corner. The rest of the team ate at the long table.

Karen unwrapped a turkey-and-swiss. "Tom. The reporting frame."

Tom set his soup spoon down. "Each of these gaps is in a different regulator's house."

"Which means?"

"Which means we map findings per-regulator. FERPA gets the admissions-side report. GLBA gets the financial-aid report. NIH gets the research-integrity report. HHS OCR gets the medical-center hallway-tour informal advisory. Civil-rights litigation defense gets a separate addendum because it cuts across FERPA and across the override-rationale gap."

Karen ate. Chewed. "Five reports out of one engagement."

"Five sections of one report. One severity scale. Five regulator audiences. Each section ends with a per-regulator summary. The General Counsel reads the whole thing. Each compliance lead reads their section."

"Phase prioritization?"

"Phase 2 closes the override-rationale gap. That's the litigation exposure. Everything else is on a longer timeline."

"Five years for the rest?"

"Three to five. NIH and DoD will force the research-computing piece in the next two cycles. HIPAA the medical center is on its own schedule. GLBA the financial aid is biennial — the next audit is the forcing function. Slate is twelve months."

Karen took a bite. Looked out the dining-room window at the quad. Two students were arguing about something on the steps of the library, gesticulating, laughing.

"Tom."

"Yeah."

"Ines told me on the prep call that the civil-rights firm's threat letter pointed specifically at override-down decisions for applicants whose AI score would have admitted them. That's the suspect class."

"You think the reconciliation test today is going to find one."

"I think the reconciliation test is going to find more than one. And I think the rationale fields on those records are going to be the legally interesting part."

"Twenty bucks."

"You're on."

They finished lunch in twelve minutes.

---

## 🔄 1:00 PM — API Layer Inspection

Mike pulled up the admissions screening service's API logs. He had a tail on the chain stream and a tail on the Slate webhook stream.

```
herald-tail --tenant=olmstead \
            --service=undergrad-admissions-screen \
            --follow --since=now
```

A reader was working through applications in real time. The screening service was returning scores. Mike watched a sequence.

```
[2026-04-09T13:14:22.317Z] entry_id=2026-04-09-UA-19044
  service=undergrad-admissions-screen tenant=olmstead
  event=applicant.scored
  applicant_hash=a4f1...c92e
  model_id=admissions-readiness-v2.1
  model_version=2026-Q1
  score=68 confidence=0.82
  shap_top5=[gpa_weighted:0.31, course_rigor:0.18,
             essay_voice:0.11, recommendation_strength:0.09,
             extracurricular_depth:0.07]
```

Two minutes later:

```
[2026-04-09T13:16:48.802Z] entry_id=2026-04-09-UA-19045
  event=reviewer.override
  applicant_hash=a4f1...c92e
  reviewer_id=R-48117
  override_decision=override-up
  override_reason_code=STRENGTH_OF_RECOMMENDATIONS
```

Mike pointed at the chain entry. "Override-up. Reason code in the chain. Reason free-text — let's see."

He flipped to the Slate webhook log. The webhook fired for the structured override fields. There was no webhook for the rationale field. The rationale lived in Slate's own database, accessible through Slate's own API, not chained.

He pulled up Slate's audit configuration through Ines's audit role. The rationale field had a `lastModifiedAt` and `lastModifiedBy`. Slate kept 30 days of field-history snapshots for fields that had field-history enabled. The rationale field had field-history enabled.

"So for 30 days the prior versions of this rationale are recoverable from Slate."

"Yes."

"And after 30 days, only the current version remains."

"Yes."

Mike wrote: *Slate webhook covers structured fields. Free-text rationale is in Slate's own database. Slate field-history is 30 days. The chain is upstream of the rationale.*

> **⚠️ Surprise #7**
> The Slate webhook into the chain covers the structured override fields — decision, reason code, reviewer ID. The free-text rationale field is not webhooked. Slate retains 30 days of field-history snapshots for that field; past 30 days, prior versions are gone. The chain has no view of the rationale at any point — current or historical.

Mike kept going. He checked the Slate API for the rationale field on five applications from the past week. The field was readable. The field-history endpoint was readable.

He made a note: *Phase 2 line-item. Webhook the rationale field on edit. Even without chaining the field, capturing every edit at the chain boundary closes the historical-rationale gap. This is engineering effort, not architecture.*

He wrote in his workbook the recommendation that the team would put in the report.

---

## 🧬 2:00 PM — Pipeline Reality (Training Pipeline + Research Computing)

The team split. Chen went to the screening-model training pipeline. Luis and Diana went on the research-computing tour with Ines.

### The training pipeline

Chen pulled up the model-training pipeline. The screening model was retrained quarterly on a fresh slice of the prior cycle's data, validated against a fairness audit, and signed off by the admissions office plus an external fairness-audit vendor before deployment.

He drew the flow on the conference-room whiteboard.

```mermaid
sequenceDiagram
    participant src as application data store
    participant tr as training pipeline
    participant fa as fairness-audit vendor
    participant art as model artifact
    participant ai as admissions chain
    src->>tr: extract training set with seed
    tr->>tr: train candidate model
    tr->>fa: submit candidate for fairness audit
    fa-->>tr: audit report with disparity ratios
    tr->>art: store model weights, hash artifact
    tr->>ai: chain entry with training seed, audit report hash, model hash
    ai-->>tr: entry id, hmac
```

Chen pointed at the third arrow. "The fairness-audit report is hashed and the hash is in the chain entry. The audit report itself sits in S3 referenced by the hash. If the audit report is altered, the hash mismatches and the verifier fails on the linked entry."

He pulled up the most recent retraining event — January 14, 2026. Verified.

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key admissions-2026-q1
```

He rehashed the audit report PDF in S3. The hash matched the chain entry.

He pulled up the model artifact. Hashed it. Matched.

He pulled up the prior retraining — October 2025. Verified. Audit report hash matched. Model hash matched.

He went back further — July 2025, May 2025. All PASS. All hashes matched.

Four retraining events. Four fairness audits. Four chain entries. Each linked the training seed, the model artifact hash, and the fairness-audit report hash. Each verified.

> **✓ Confirmation #4**
> The training pipeline integrity for the admissions screening model is the cleanest the team has seen at a university. Quarterly retraining. Fairness audit by external vendor for each retraining. Audit report hash and model artifact hash both in the chain entry. Four retraining events sampled across eleven months. All PASS. All hashes match. The disparate-impact threat letter forced this discipline up front, and it shows in the artifacts.

Chen wrote in his workbook: *Best training-pipeline integrity I have seen at a university. It is because they had to defend it in advance.*

### The research-computing tour

Luis and Diana went with Ines across the quad to the engineering and biological-sciences buildings. Two lab visits. Forty-five minutes total.

**Lab A — microbiology, biological sciences building.**

The PI was a senior microbiologist. NIH-funded. Two postdocs and three graduate students. The lab kept code in a university-managed GitLab. Experimental data was snapshotted nightly to a campus institutional-data bucket. Each experimental run was tagged by date and protocol version. The PI walked Luis through the snapshot policy. Luis nodded.

"You version your data."

"We version our data. The institutional snapshot policy is documented. Every PhD in this lab knows the rule. We retain three years of nightly snapshots."

"And the snapshots are immutable?"

"S3 with versioning enabled. Object lock for the past 18 months."

Luis wrote: *Lab A is the well-run case. Code in GitLab. Data snapshots nightly. Object lock on 18 months of snapshots. Tagging convention by date and protocol version. This is what compliance looks like when the PI cares.*

> **✓ Confirmation #5**
> Lab A — microbiology — runs a research-computing operation that meets NIH research-integrity expectations: code in version control, data snapshots nightly to an immutable bucket, 18 months of object-locked retention, tagging by date and protocol version. This is achievable at the lab level when the PI prioritizes it.

**Lab B — quantitative finance, business school.**

Different building. Different culture. The PI was a finance professor with two co-PIs and four PhDs. The lab ran on a single Linux server in a closet behind a grad-student bullpen. Diana logged in over SSH using credentials Ines had been given for the audit.

The login prompt did not ask for a username. It asked for a password. The username was hardcoded — `qfin_lab`.

"Shared account."

The PI nodded. "Six of us use it. Same password."

Diana asked, "How long has this password been the same?"

"I don't know."

Luis pulled `last` on the server. The login history showed two thousand eight hundred sessions in the past year, all from `qfin_lab`. Different IPs. Different times of day. No way to tell which human was logged in at any given moment.

Luis: "Where is the working data?"

The PI pointed at the server. "Some of it is on `/data`. The current working set is on a USB drive in my office."

"USB drive."

"It's backed up."

"To?"

"My laptop."

Luis wrote: *Lab B is the median case. Shared password to a Linux server. Two PIs, four PhDs, all logged in as the same user. Working data on a USB drive in the PI's office, backed up to the PI's laptop. No version control. NIH and DoD expectations are not met.*

> **⚠️ Surprise #8**
> Lab B — quantitative finance — runs a single shared `qfin_lab` Linux account that two PIs and four PhDs all log in as. Working data lives on a USB drive in the PI's office, backed up to the PI's laptop. No version control. No snapshots. No tagging. Twenty-eight hundred login sessions in the past year, all under the same shared identity. This is the median case for faculty-led labs at Olmstead, not the exception.

Ines, walking out of the building with Luis and Diana, did not say much. She said enough.

"Lab B is what we have to figure out. The Faculty Senate will resist any chain mandate. Central IT cannot enforce. NIH and DoD are increasingly asking. This is a five-year roadmap, not a twelve-month one."

Diana wrote: *Five-year roadmap. Faculty federalism. Phased policy. Grant-funded labs that publish are the priority. Document.*

They walked back across the quad to the conference room. The wind off the lake had picked up.

---

## 📊 3:00 PM — Reconciliation Test

Tom set the test. Five admissions decisions from the past 60 days. Trace each one end to end. Ines picked them — five entry IDs sent to Karen by email at 2:55. The team did not know which five until 3:00.

```
2026-02-14-UA-08221
2026-02-22-UA-09817
2026-03-04-UA-12044
2026-03-18-UA-15903
2026-04-02-UA-18441
```

Mike, Chen, Raj, Diana, and Elena each took one. Twenty-five minutes.

### 2026-02-14-UA-08221

Mike: "Verifier — PASS. Twelve steps. Public key admissions-2026-q1. Score 73, confidence 0.86. Override decision: accept-AI. Override reason code: not applicable — no override. Reviewer ID R-48104."

"And the rationale?"

He pulled Slate. "Rationale field present. Last modified February 14 at 3:48 PM. No subsequent edits. Field-history shows the original entry — three sentences explaining the reviewer agreed with the AI score and the application was complete. Original is the only version."

"Clean trace."

Karen wrote: *Reconciliation 1 — full trace. AI clean. Reviewer rationale traceable. No override.*

### 2026-02-22-UA-09817

Chen: "Verifier — PASS. Score 64, confidence 0.71. Override decision: override-up. Override reason code: STRENGTH_OF_RECOMMENDATIONS. Reviewer R-48117."

"Rationale?"

"Slate field present. Last modified February 22. No subsequent edits. Field-history shows original. Two paragraphs explaining the recommendation letters from the applicant's research mentor weighted heavily. Original is the only version."

Karen wrote: *Reconciliation 2 — full trace. AI clean. Override-up captured. Rationale original preserved.*

### 2026-03-04-UA-12044

Raj: "Verifier — PASS. Score 81, confidence 0.91. Override decision: override-down. Override reason code: ESSAY_INCONSISTENCY. Reviewer R-48104."

The room shifted half a degree.

"Rationale?"

Raj pulled Slate. "Rationale field present. Last modified — March 4 at 2:11 PM. And then — let me check field-history."

He scrolled. The field-history showed two entries. The second entry — the current value — was added April 6.

"Edited."

Tom: "What was the original?"

"Field-history shows the March 4 entry. Two sentences flagging an inconsistency between the applicant's essay and the supplemental questions. The current April 6 entry — three sentences, expanded reasoning, slightly different language."

Mike: "Within the 30-day window. Original is recoverable."

Diana: "Just barely. April 6 plus 30 days is May 6. Today is April 9. We're inside the window. If this same record had been examined in May, the original would be gone."

Karen wrote: *Reconciliation 3 — AI clean. Override-down captured. Rationale edited April 6, original recoverable in field-history. Within window.*

### 2026-03-18-UA-15903

Diana: "Verifier — PASS. Score 77, confidence 0.84. Override decision: override-down. Override reason code: ACADEMIC_FIT_CONCERN. Reviewer R-48117."

The room shifted again.

"Rationale?"

She pulled Slate. The field-history showed the current value as the only entry. Last modified February 19.

"February 19? That's before the AI score date."

She looked again. Her face changed.

"That's the wrong applicant. Let me re-pull. Sorry."

She re-pulled with the applicant_hash from the chain entry. The Slate record came back. The rationale field's last-modified was March 18 at 4:02 PM. Field-history showed one prior version, edited away on March 19.

"Edited the next day."

"Within field-history window. Recoverable. The March 18 original was four sentences. The March 19 current is two sentences."

"Same gist?"

She read both. "Same gist. Trimmed. The original mentioned a specific course the applicant had not taken. The current is more general."

Karen wrote: *Reconciliation 4 — AI clean. Override-down captured. Rationale edited March 19, original recoverable. Trimmed, not contradicted.*

### 2026-04-02-UA-18441

Elena: "Verifier — PASS. Score 79, confidence 0.88. Override decision: override-down. Override reason code: HOLISTIC_REVIEW_PRIORITY. Reviewer R-48104."

The room got quiet.

Karen: "Rationale?"

Elena pulled Slate. The rationale field-history showed only the current value. Last modified April 2 at 10:14 AM. No field-history entries before that.

"Field-history shows no prior versions."

"How can that be? The field is supposed to retain 30 days."

Elena dug. "The field-history retention is on the field schema. If the field was originally null and the first edit was the current value, there is no prior version to record."

She pulled the Slate audit log directly. The audit log showed the rationale field had been edited twice — once on April 2 at 10:14 AM (the current value) and once on April 2 at 9:42 AM, 32 minutes earlier. The 9:42 AM version had been overwritten 32 minutes later.

"Field-history fired on the second edit but not the first."

"Why?"

"Slate's field-history captures snapshots on edit. The first edit replaces null with content. The snapshot of the prior state — null — is not retained as a useful version. The 9:42 AM content was overwritten by the 10:14 AM content within the 30-day window, but the audit log shows it was overwritten, not what it was."

Mike: "So the original 9:42 AM rationale is gone."

"The original is gone."

The room held still for a second.

Tom: "Override-down with no recoverable original rationale."

Elena: "Override-down with no recoverable original rationale. The April 2 record is what we have."

Karen wrote: *Reconciliation 5 — AI clean. Override-down captured. Rationale edited within 32 minutes of being entered. Original gone. Audit log records the overwrite event but not the prior content.*

### Tally

Karen put it on the board.

| Decision | AI side | Reviewer decision | Rationale traceable |
|---|---|---|---|
| 2026-02-14-UA-08221 | PASS | accept-AI | YES, original preserved |
| 2026-02-22-UA-09817 | PASS | override-up | YES, original preserved |
| 2026-03-04-UA-12044 | PASS | override-down | YES, edited April 6, original recoverable |
| 2026-03-18-UA-15903 | PASS | override-down | YES, edited March 19, original recoverable |
| 2026-04-02-UA-18441 | PASS | override-down | NO, original overwritten same day, gone |

Five-of-five AI side PASS.
Five-of-five reviewer decision captured in chain.
Three-of-five reviewer rationale recoverable.
Two-of-five reviewer rationale gone.

Both of the gone-rationale cases were override-down decisions for applicants whose AI score would have admitted them.

The room was quiet. Ines was sitting with her hands folded on the table.

Karen broke the silence. "Ines, the threat letter."

Ines did not look up. "The threat letter named override-down decisions for applicants whose AI score would have admitted them as the suspect class. Both of these would be in that class if the firm files."

"And you cannot prove what the rationale was at the time of the override-down on one of these two."

"No."

> **✓ Confirmation #6**
> Five admissions decisions traced end to end. Five-of-five AI-side PASS. Five-of-five reviewer decisions captured in the chain with reason codes. Three-of-five reviewer rationales fully traceable to the original entry. The chain itself behaves exactly as designed — the chain is not where the gap is.

> **⚠️ Surprise #9**
> Two-of-five reviewer rationales are gone. Both are override-down decisions on applicants whose AI score would have admitted them. One has the original recoverable through Slate field-history because the edit was within the 30-day window. The other has no recoverable original because the rationale field was edited a second time within 32 minutes of being entered, and Slate's field-history did not capture the first content. This is the litigation exposure.

Karen looked at Tom. He looked back. Neither of them said anything about the bet.

---

## 😬 3:45 PM — The Friction Builds (Lab B's PI)

Ines had asked Lab B's PI to come over for a sit-down. He showed up two minutes late, in a sport coat over a Patagonia vest, looking like he had three other things to do.

He sat down. He did not put his phone away.

Karen introduced herself. "We visited your lab earlier this morning."

"Right. Diana and Luis."

"We have some questions about the data-integrity practices."

He waited.

"The shared `qfin_lab` account. The USB drive in your office. The lack of version control."

He set his phone face-down on the table. His eyes narrowed.

"Look. We have other priorities. This is research. We publish. We get grants. We are not running a financial-services audit operation. The data we work with is publicly available — equity returns, options chains, macro indicators. There is no PII. There is no PHI. There is nothing that would survive a HIPAA audit because there is nothing here that HIPAA would care about. I appreciate that there are integrity questions, but the framing is mismatched."

Karen heard him out. She did not interrupt.

"Professor, I am not here to argue your priorities. I am here to document what is and what is not. NIH and DoD have started asking questions about research-computing integrity at universities. That is the context. The framing is not 'we are auditing your lab.' The framing is 'when the funding agency asks, what do you say.'"

He looked at Ines.

Ines: "We are not going to require chain on every lab. We are going to require it on grant-funded labs that publish, and we are going to phase it. Your lab is grant-funded. Your lab publishes. You are in the first cohort, but the cohort is twelve months from now and the implementation will be light-touch. I am telling you the answer before you ask the question."

The PI's posture shifted. He picked his phone back up but did not unlock it.

"Twelve months."

"Twelve to eighteen. We will have a working group. You will be on it. We will not impose a system from above. We will design it with you."

"Working group I can do."

"Working group it is."

He stood up. He shook Karen's hand. He left.

Karen wrote in her notebook: *Lab B PI defused. Ines handled it. Faculty federalism is a political problem before it is a technical problem. The chain is not the answer to faculty federalism. The chain is the answer to a specific use case where the institution has decided enforcement is required. The decision precedes the chain. Without the decision, there is no chain.*

> **⚠️ Surprise #10**
> Faculty governance is the political third rail. Central IT cannot enforce. The Faculty Senate will resist any chain mandate. NIH and DoD are increasingly asking. The remediation is a working-group-led twelve-to-eighteen-month design effort, not a top-down system rollout. The cost of getting this wrong is loss of faculty trust, which costs more than the audit finding.

Ines closed her notebook. "Now the medical center."

---

## 🔍 4:30 PM — The Medical-Center Hallway Tour and the Litigation Question

The team walked across the quad. The medical-center towers were three minutes away on foot. Ines badged them through the connector building. The medical-center compliance office was on the fourth floor. She had cleared a thirty-minute slot with the medical-center CISO for a hallway conversation.

The CISO met them in the corridor. Black turtleneck, badge on a lanyard, the look of someone who had been at the hospital for twelve hours already.

"Thirty minutes. What do you need to see."

Karen: "Epic clinical-notes mutability. MyChart audit posture. The medical school's research-database overlap with clinical systems. The lab's specimen-tracking pipeline."

The CISO walked them through. Epic clinical notes — same as Mercator. Notes were mutable through addendum until co-signed; after co-sign they were locked but the original-vs-addendum diff was retained. MyChart audit logs were on, retention was 90 days for free-text fields and one year for structured fields. The medical school's research databases shared a data-warehouse layer with the hospital's clinical data warehouse — research could query clinical, with consent and IRB approval, and the queries were logged but the query-result snapshots were not always retained. The specimen-tracking pipeline through the pathology lab was instrumented at the LIMS layer; retention there was three years.

Karen: "What about anything chain-coupled?"

The CISO: "Nothing. We have a tenant ID reserved for medical-center deployment. We have not deployed."

"Why?"

"Time and money. We have other priorities. The board has not made a decision. We are watching what the admissions side does."

Mike wrote: *Medical-center tenant reserved but not deployed. Same Epic-side pattern as Mercator. The medical-center compliance office is staffed and operational; the audit problem is not absence of governance but absence of cryptographic enforcement.*

> **⚠️ Surprise #11**
> The medical-center side mirrors the Mercator findings: Epic clinical-notes mutable through co-signed addendum, MyChart audit retention 90 days for free-text, research-clinical data-warehouse overlap with logged queries but inconsistently retained query-result snapshots. The medical-center compliance office is competent and aware. The board has not decided to fund chain deployment. A tenant is reserved.

The CISO checked his watch. Twenty-eight minutes had passed.

Karen wrapped up. "Thank you. We will write this up as informal advisory, not a finding. Ines will route it to your office formally."

The CISO nodded. "Appreciated."

They walked back to the connector. Ines spoke for the first time since they had crossed.

"Karen. The litigation question."

"Go."

"If the civil-rights firm files suit and asks for the rationale behind the override-down decisions, what can we produce?"

Karen had been holding the answer since the reconciliation test.

"For decisions in the chain era, you produce the override_reason code, which is in the chain and sealed. You produce the rationale field's current state. If the rationale has been edited within the past 30 days, you produce the field-history snapshot from Slate showing the prior version. If the rationale has been edited and the original is older than the 30-day window, you cannot prove what the original was. The audit log records that an edit happened; it does not record what was overwritten. For the two override-down records we found this afternoon, one is recoverable through field-history because it is within the window — the other is gone because the second edit happened the same day the first edit happened, and field-history did not snapshot the first content."

"What does that look like in court?"

Tom answered. "It looks like an expert witness for the firm asking, 'Was the rationale you have today the rationale you had at the moment of the override?' For two of the five we sampled, the answer is 'we cannot demonstrate that.' If the same sampling rate holds across the full population of override-down decisions in the chain era, the firm will find cases where the answer is 'no, the rationale was edited and the original is gone, and the only thing in the chain is the structured reason code.' The structured reason code is defensible — it came from a controlled vocabulary, it was selected by the reviewer at the time of the decision, and it is sealed. The free-text rationale is a different artifact. The argument the firm will make is that the free-text rationale was the actual decision-rationale and the structured reason code was a checkbox. They may or may not prevail on that argument, but it is the argument they will make."

Ines was nodding slowly.

"What is the remediation?"

Mike: "Two parts. One — webhook the rationale field on edit, so every prior version of the rationale lands in the chain at the boundary. That closes the historical-rationale gap from this point forward. Two — extend Slate field-history retention beyond 30 days. Slate supports it; it is a configuration change. That gives you a longer lookback window for any rationale that was edited."

"Phase 2."

"Phase 2."

"Funded?"

Ines did not answer immediately. Then: "After today, yes. I will take this to the General Counsel Friday."

Karen wrote: *Phase 2 funding catalyzed by today's findings. The 30-second elevator pitch to GC: 'we cannot prove the original rationale on two-of-five override-down decisions in our sample, both for applicants whose AI score would have admitted them, both in the suspect class. Webhook the rationale field. Close the gap.'*

> **✓ Confirmation #7**
> The chain on the AI side supports the disparate-impact litigation defense for the structured artifacts — score, model version, override decision, override reason code, fairness audit by retraining. The remediation for the rationale gap is engineering effort at the Slate webhook boundary plus a Slate configuration change, not architectural rework. Phase 2 is now funded.

They walked back to the conference room.

---

## 🌆 5:30 PM — Auditor Debrief

The team reconvened. Coffee was cold. The afternoon light had turned that flat upper-Midwest gray. The library across the quad had its lamps on.

Karen stood at the whiteboard. Five rows.

| Regulator | Status |
|---|---|
| FERPA (admissions) | 0 Gaps on AI side, 1 Gap (override-rationale free-text not chained), 1 Partial (Slate field-history retention) |
| GLBA (financial aid) | 3 Gaps, 2 Partials |
| HIPAA (medical center) | Out-of-scope, informal advisory list mirrors Mercator |
| NIH (research integrity) | 1 Gap, 2 Partials |
| Civil-rights litigation defense | Chain supports AI-side defense; override-rationale gap is the exposure to track |

"That's the shape."

Ines stood at the side, arms folded, listening.

"FERPA admissions side. The chain is mature. Confirmed: chain integrity, append-only ledger behavior, credential rotation under chain, training pipeline integrity with quarterly fairness audits, reconciliation five-of-five PASS on the AI side, customer-litigation defense fully serviceable for the structured artifacts. Gap: the override-rationale free-text lives in Slate, not in the chain, and Slate's 30-day field-history retention is the only protection against silent overwrite. Partial: Slate field-history is configurable but currently set to 30 days. Phase 2 in twelve months — webhook the rationale at edit, extend Slate retention. Funded as of today."

She moved to GLBA.

"GLBA financial-aid side. Three gaps. Banner audit-log retention at 90 days. Shared `aid_admin` account during peak FAFSA season with SSO bypass. Paper-based annual administrative-access review. Two partials. Audit configuration coverage is not 100% across all Banner modules. Encryption-at-rest documentation is incomplete. The biennial GLBA audit is the forcing function. Phase 3 territory. Eighteen to twenty-four months."

She moved to HIPAA.

"HIPAA medical-center side. Out of scope for this engagement. Informal advisory list — Epic clinical-notes addendum mutability, MyChart 90-day retention on free-text fields, research-clinical data-warehouse query-result snapshots inconsistently retained. Same shape as Mercator three weeks ago. Recommendation: route through Ines's office to the medical-center compliance office formally. The medical-center board has not decided to fund chain deployment. A tenant is reserved."

She moved to NIH.

"NIH research-integrity side. One gap. No enforcement at the lab level — central IT cannot mandate, faculty federalism is the political constraint. Two partials. Institutional snapshot policy is uneven across labs — Lab A is the well-run case, Lab B is the median. IRB approval audit trail is in a homegrown SQL database with mutable approval rows. Phased five-year roadmap. Working-group-led design with grant-funded publishing labs as the first cohort. Twelve to eighteen months for the first cohort."

She moved to civil-rights litigation.

"Civil-rights litigation defense. The chain on the AI side supports the disparate-impact defense for the structured artifacts — score, model version, override decision, override reason code, fairness audit. That's a strong defense and it is the reason the chain was deployed. The exposure is the override-rationale gap. We sampled five decisions today; two-of-five had unrecoverable original rationales, both override-down decisions on applicants whose AI score would have admitted them. Both in the suspect class the threat letter named. Phase 2 closes this gap going forward. The historical exposure — for decisions made between May 12, 2025, and the Phase 2 deployment — is what it is. Document the boundary in the litigation-defense memo to the General Counsel."

She put the pen down.

Ines spoke. "What do I take to the GC Friday?"

Tom answered. "Three things. The five-row summary. The reconciliation test as a real artifact — five decisions, three traceable, two not. The Phase 2 scope and cost — Slate webhook plus retention extension — with the framing that this closes the litigation exposure going forward."

"And the rest?"

"The rest is a longer conversation. GLBA is biennial, the next audit is your forcing function, and the GLBA report writes itself off this material. NIH is a working group and a five-year roadmap. The medical center is a hallway tour we route formally to their compliance office. The Faculty Senate item we put in the appendix because the political handling matters more than the audit framing."

Ines nodded. "That tracks."

Karen closed her notebook. "We will have the report Thursday. Friday morning before your GC review."

Ines's shoulders dropped that quarter-inch. "Thank you."

The team packed up. Raj and Luis loaded the boxes of evidence into the rental SUV. Diana and Elena said goodbye to Ines at the connector building. Mike and Chen took one last look across the quad at the library lamps.

Karen walked out last. She turned at the doorway and looked back at the conference-room window — at the empty whiteboard, the coffee cups, the five rows that would become Friday's memo.

> **🔍 Karen's note (internal):**
> *It never is. Sometimes the chain is on the part that's being sued, and that's the only part that needs to be.*
>
> *The chain works. The chain is not the gap. The gap is a free-text field one webhook away from being closed. Phase 2 is twelve months. The litigation question is going to be asked between now and Phase 2 deployment. That window is the report.*

---

## ✅ vs ❌ — The Five-Regulator Summary

### ✅ FERPA (Admissions AI)

| Item | Status |
|---|---|
| Chain integrity (HMAC + Merkle + daily Ed25519 seal on AWS CloudHSM `us-east-2`) | PASS |
| Append-only ledger behavior under direct DB mutation attempt | PASS — verifier catches at HMAC layer |
| Multi-entry tamper attempt | PASS — verifier catches at Merkle/seal layer |
| Credential rotation under chain | PASS — eleven rotations sampled across eleven months, all PASS |
| Training pipeline integrity (quarterly retraining, fairness audit linked by hash) | PASS — four retrainings sampled, all hashes match |
| Reconciliation test (5-of-5 AI side) | PASS — verifier PASS for all five, override decisions captured |
| Pre-chain decisions (before May 12, 2025) | REJECTED by verifier as designed; documented as out of chain scope |
| Override-rationale free-text in Slate | GAP — not chained, lives in Slate, 30-day field-history retention |
| Slate field-history retention | PARTIAL — 30 days, configurable, currently the only historical-rationale protection |

### ❌ GLBA (Financial Aid — Banner)

| Item | Status |
|---|---|
| Banner audit-log retention | GAP — 90 days, edits older than 90 days have no recoverable original |
| Shared `aid_admin` account | GAP — six or seven counselors during peak season, SSO bypass, last rotation date unknown |
| Annual administrative-access review | GAP — paper-based, well-intentioned, largely unverifiable |
| Audit configuration coverage | PARTIAL — enabled on most modules, gaps in some |
| Encryption-at-rest documentation | PARTIAL — present but incomplete |
| GLBA biennial audit timeline | Forcing function — Phase 3 territory, 18-24 months |

### ⚠️ HIPAA (Olmstead Medical Center) — Informal Advisory

| Item | Status |
|---|---|
| Epic clinical-notes mutability | Same as Mercator — addendum-based, original-vs-addendum diff retained post co-sign |
| MyChart audit retention | 90 days free-text, 1 year structured |
| Research-clinical data-warehouse overlap | Logged queries, inconsistently retained query-result snapshots |
| Pathology LIMS specimen-tracking | 3-year retention, instrumented |
| Chain deployment | Reserved tenant, not deployed, board has not funded |
| Engagement scope | Out-of-scope; route formally to medical-center compliance office |

### ❌ NIH (Research Integrity)

| Item | Status |
|---|---|
| Lab-level enforcement | GAP — central IT advisory only, faculty federalism, Faculty Senate political constraint |
| Institutional snapshot policy uniformity | PARTIAL — Lab A meets expectations, Lab B does not, median is closer to Lab B |
| IRB approval audit trail | PARTIAL — homegrown SQL database, mutable approval rows, cross-organizational governance |
| Lab A (microbiology) | Meets NIH expectations — git, nightly snapshots, object lock, tagging |
| Lab B (quantitative finance) | Shared account, USB drive, no version control — does not meet NIH expectations |
| Roadmap | Five-year, working-group-led, grant-funded publishing labs first cohort, 12-18 months |

### ⚖️ Civil-Rights Litigation Defense

| Item | Status |
|---|---|
| AI-side structured artifacts (score, model version, override decision, reason code) | DEFENSIBLE — chain-sealed, verifier PASS, hash-linked fairness audits |
| Training pipeline (quarterly retraining + external fairness audit) | DEFENSIBLE — four retrainings sampled, audit reports hash-linked in chain |
| Reviewer rationale free-text (current state) | DEFENSIBLE — current value retrievable from Slate |
| Reviewer rationale free-text (historical, within 30 days) | DEFENSIBLE via Slate field-history snapshot |
| Reviewer rationale free-text (historical, edited and overwritten same-day) | NOT DEFENSIBLE — original gone, audit log records overwrite event but not prior content |
| Sampling result | 5-of-5 AI side PASS, 5-of-5 reviewer decision captured, 3-of-5 rationale traceable, 2-of-5 rationale gone |
| Both gone-rationale cases | Override-down decisions on applicants whose AI score would have admitted them — the suspect class named in the threat letter |
| Phase 2 remediation | Webhook the rationale field at edit; extend Slate field-history retention. Funded as of today. 12 months. |

---

## 🧾 Final Assessment Theme

> *"Olmstead can defend the AI admissions decision under the disparate-impact threat. Olmstead cannot prove that the human review was made on the basis of the documented reasoning at the time of the decision. The chain is on the part that decides — the algorithm. The free-text rationale is the part that explains the human override of the algorithm, and that part is one webhook away from the chain. Twelve months. Phase 2. Funded today."*

Olmstead University demonstrates AI-decision integrity within scope. The undergraduate admissions screening chain is mature, verifiable, and tied to a consent-to-resolve framework with a civil-rights firm. The training pipeline is the cleanest the team has seen at a university — quarterly retraining, external fairness audit linked by hash for each retraining, four retrainings sampled across eleven months, all hashes match. The chain on the structured override decision and reason code is defensible. The credential-rotation history and the IAM around the AI service are chain-coupled.

The gap is the free-text rationale. The Slate webhook covers the structured override fields. It does not cover the rationale free-text. Slate's 30-day field-history retention is the only historical-rationale protection. In a five-decision reconciliation, three rationales were traceable to the original; two were not. Both of the not-traceable cases were override-down decisions on applicants whose AI score would have admitted them — the suspect class the civil-rights threat letter named. Phase 2 closes this gap going forward — webhook the rationale field on edit, extend Slate field-history retention — and is funded as of today.

Outside the AI scope, Olmstead is a higher-education hybrid of Mercator and Stelvio. The medical center is a separate audit problem under HIPAA, with a competent compliance office, a reserved-but-undeployed tenant, and the same Epic-side findings the team wrote at Mercator three weeks ago. The financial-aid GLBA stack is mutable Banner with 90-day audit retention and a shared peak-season account. The research-computing side is faculty-led federalism with one well-run lab and a median lab that does not meet NIH expectations; central IT cannot enforce; the Faculty Senate will resist any chain mandate; the remediation is a five-year working-group-led design effort with grant-funded publishing labs as the first cohort.

Five regulators. Five sections. One severity scale. Five remediation timelines. The seam between the chained AI service and the unchained free-text field is the single most important line in the litigation-defense memo. Olmstead knows where the seam is. The civil-rights firm is going to ask. The window between now and Phase 2 deployment is the exposure to track.

---

*End of diary. Filed Wednesday evening. Report drafted Thursday. Delivered Friday morning before the General Counsel review.*
