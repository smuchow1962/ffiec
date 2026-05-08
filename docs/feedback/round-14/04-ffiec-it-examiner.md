# Round 14 — FFIEC IT Examiner review

**Reviewer.** Brendan Walsh, Senior IT Examiner, Office of the Comptroller of the Currency. 16 years on the IT examination side, mostly large-bank dedicated teams plus a stretch on the FFIEC Technology Subcommittee. I run examinations under the FFIEC IT Examination Handbook — Information Security, AIO, and Audit booklets are my daily work.

**What I read.** spec/chain-of-custody-v1.md §7 step 11 (dual-algorithm cases) and §7 step 12a; examiner-quickstart.md; regulator-pack/sample-report.md; regulator-pack/finding-language.md; regulator-pack/handbook-mapping.md; regulator-pack/deployment-package.md; regulator-pack/examiner-training.md; regulator-pack/examination-response-workflow.md; portfolio-comparison-procedures.md.

**What I am evaluating.** Whether the package is examinable on a normal IT examination — meaning a junior examiner can run it, a senior examiner can defend the finding language to bank counsel, and the EIC can roll it up across a supervisory portfolio. I am not evaluating the cryptography. I am evaluating the workflow.

**Stopping criterion.** 0 net-new gaps / 0 net-new defects.

---

## 1. The headline read

This is a regulator pack written by people who have actually sat in an examination room. I do not say that lightly. Most "examiner-facing" documents I receive from technology vendors are marketing material with the word "auditor" pasted over the word "customer." This pack is not that. The five-step workflow doc is the giveaway — somebody on the writing team knows that the moment the verifier produces a JSON failure record, the examiner needs one document that runs JSON-record → severity → IR scenario → evidence request → finding paragraph in under a minute, and that document needs to exist or the workflow falls apart at the desk.

The package would survive a Handbook-aligned examination. The mapping at handbook-mapping.md hits IS II.C.10 (Logging) and II.C.13 (Cryptographic controls) where the chain belongs, and the AIO II.E (Change management) row is the one I would have flagged as missing if it were not already there — `format_version` and the dual-algorithm posture are change-management concerns, not just cryptographic concerns, and the doc treats them that way.

The dual-algorithm transitional period at §7 step 11 cases (a)-(e) is the part of the spec I expected to be hand-wavy. It is not. The five cases are enumerated, each has a strict-mode and non-strict disposition, case (e) carries an explicit severity-stays-Severe-regardless-of-bracket note that I would otherwise have to argue with bank counsel about, and the working-paper convention for recording both algorithms' validation results is written down in the spec text where it belongs rather than buried in an examiner doc.

Step 12a (GenAI model identifier completeness) is a small addition with a load-bearing purpose: it makes SR 11-7 reproducibility a chain-detected property rather than a documentary attestation. That matters because SR 11-7 model risk findings are some of the hardest to write up — you usually have a vendor-provided spec sheet, a bank-side validation report, and no third document that says the actual production calls used the model the validation report covers. The chain now produces that third document.

Bottom line: I would take this to an examination tomorrow.

The remainder of this review is the things I would still want to see, ordered by whether they would block an examination today (none would) or would improve the next iteration (a handful).

## 2. Things I checked specifically and found in good shape

### 2.1 The five-step workflow doc

`examination-response-workflow.md` is the document I would most miss if it were not in the pack. The shape — JSON record → step lookup → IR scenario → evidence request → finding paragraph — matches how examiners actually work. The two worked examples (step 8 key_fingerprint mismatch, step 10 Merkle root mismatch) are structurally different failure types, which is the right pair to choose; if both worked examples had been step 8 variants the doc would have under-demonstrated its generality.

The "When this workflow does NOT apply" section at the bottom is small and load-bearing. PASS-with-anomaly results are the case where new examiners over-reach — they treat operational anomalies as findings and the institution rightly pushes back. Having one paragraph that says "PASS-with-anomaly is an operational signal, not a finding; it lands in the report under Anomalies noted with a sentence each" prevents that.

### 2.2 Severity table parity across documents

The severity guidance in finding-language.md, the failure-mode table in examiner-quickstart.md, and the 30-minute training table at examiner-training.md Module 5 all key off the same `step` field and produce the same severity column. This is the single biggest source of cross-document drift in a regulator pack and the doc set has it consistent. I checked step 8, step 9, step 10, step 11 (single-algo), step 11 (dual-algo case e), step 12 dev-mode, and the file pre-flight truncation row across all three docs and got the same severity in every case.

The one nuance worth keeping consistent in future iterations: the file pre-flight truncation row says "Operational; sealing-delay-equivalent severity" in two docs and "Operational, NOT tampering — Medium severity" in the quickstart. Both convey the same thing but the wording is slightly different. Pick one phrasing and use it everywhere on the next pass.

### 2.3 Dual-algorithm case (e) treatment

§7 step 11 case (e) — both signatures present, one valid + one invalid — is the case where I would have expected the spec to soften the severity language because the un-broken algorithm provides integrity assurance. The spec does not soften it. The text at §7 step 11 case (e) says the severity is Severe regardless of bracket, the PASS-WITH-ANOMALY disposition under non-strict reflects integrity assurance for downstream consumers and is NOT a low-severity claim, and finding-language.md row "11 (dual-algo) co-signed seal failure" carries Severe MRA explicitly.

That language survives bank counsel. If the verifier dispositioned case (e) as PASS-WITH-ANOMALY and the finding-language doc treated it as Observation, counsel would argue the finding away on the grounds that the verifier itself disagrees with the severity. Because the disposition and the severity are separated explicitly — one for downstream-consumer integrity, the other for examination response — the finding stands.

### 2.4 Step 8 vs step 9 distinction

Quickstart §"Three specific points to note" carries the distinction that step 8 (key_fingerprint mismatch) is investigated through the IKM roster and step 9 (payload_hash MAC mismatch) is investigated through the chain content. Both are Severe. Both are "stop and call the bank" findings. The investigation paths are different.

This distinction matters because the natural examiner instinct on any cryptographic-controls finding is to ask the institution for the chain content. For step 8 that is the wrong evidence — the chain content is fine, the IKM-roster row is wrong, and asking for chain content wastes examination time and can introduce sampling errors. The quickstart catches this. The finding-language paragraph for step 8 reinforces it. The five-step workflow's evidence-request column points to the institution's `master.reconciliation_completed` event for step 8 and to the `chain.verification_failure` event for step 9. The doc set is consistent on the distinction across all three layers.

### 2.5 Step 12a (GenAI model identifier) treatment

The completeness check fires only on chain entries representing model calls (entries with any `gen_ai.*` attribute), which means tool calls and audit-only events are not affected. That scoping decision is correct — over-broad scoping would make every chain entry a candidate for the check and would produce noise on entries where the GenAI semconv attributes do not apply.

The strict-mode disposition is FAIL; non-strict is PASS-WITH-ANOMALY treated as control-completeness for SR 11-7 reproducibility. That mapping is the right one. SR 11-7 examiners work from model documentation; the chain produces a per-call record of the actual model the vendor's API answered with (the response-side identifier, which the spec rightly notes is the load-bearing record because vendors silently re-route between model versions during outages). I have written up "vendor silently swapped models during validation period" findings on three examinations in the last four years; each took weeks of evidence-gathering to substantiate. The chain produces the evidence as a chain-detected property.

The verifier reason string is `gen_ai_model_identifier_missing at seq N: {field_name} required for chain entries representing model calls`. The `{field_name}` substitution distinguishes "request.model missing" from "response.model missing" in the failure record, which is the distinction the institution needs to investigate the SDK build (request.model missing typically means SDK shipped a chain entry without populating the OTel GenAI semconv attribute) versus the vendor integration (response.model missing typically means the SDK is not reading the response-side identifier from the vendor's API response).

Recommend the failure-mode tables in examiner-quickstart.md and finding-language.md add a row for step 12a so it gets the same severity-and-IR-scenario treatment as the other rows. I checked: step 12a is in the spec, in the workflow doc's IR-scenario table (under "12 (cadence)" and "12 (dev-mode)" the table does not have a 12a row yet), and in the training doc Module 5 — but it is not in the failure-mode tables in the quickstart and finding-language docs. Add it. See Gap 5.1.

## 3. Things I would have flagged as missing if they were missing

### 3.1 Verifier behaviour when the IKM is not provided

Examination practice varies on whether the institution provides the IKM at examination time. Some institutions provide it routinely under a documented disclosure shape. Some require a court order or protective-order disclosure first. Some refuse and force structural-only verification.

The pack handles all three cases. The verifier degrades to structural-only verification when `--master-key` is absent (chain links, Merkle, signature, but no per-event HMAC equality), reports as PASS-WITH-ANOMALY with `structurally consistent, key-bound verification skipped`, and elevates to FAIL under `--strict`. The deployment package, the examiner training, the quickstart, and the sample report all carry this language consistently.

The legal-disclosure cross-reference (`legal-disclosure.md` §"Court-ordered master-key disclosure") and the customer-dispute cross-reference (`customer-dispute-procedures.md` §"IKM access for customer-side verification") give the examiner a documented path back to the institution-side disclosure shape. I would not normally check those cross-references during a routine examination, but if the institution refused the IKM and we needed to escalate, having the path documented prevents the examination from getting stuck on a procedural argument.

### 3.2 Verifier-version skew

The pack treats "the institution's SDK is newer than my verifier" as `step: 1 format_version not supported by this verifier (running v1)` and explicitly says this is not an institutional finding. The examiner obtains the matching verifier from the project supply chain.

That handling is correct. Without it, an examiner with an older verifier and a bank with a newer SDK would produce a FAIL report against the bank for what is actually a tooling issue on the regulator's side. The most-specific-first refusal at §7 step 1 catches the version skew before any other check fails, and the failure-mode tables across all three docs label the row as "Not an institutional finding."

This is the kind of detail that sounds obvious in the abstract and is routinely gotten wrong in regulator packs that try to be comprehensive. The pack gets it right.

### 3.3 The bundle as a working-paper artifact

`--bundle` produces a tarball that becomes the working-paper artifact. The deployment-package doc says the verifier is read-only on the ledger snapshot and the bundle's SHA-256 is filed in the working-paper system. The training doc says the same.

That separation matters. The working-paper artifact has to be byte-identical regardless of which examiner ran the verifier, and the bundle has to capture enough state that a second examiner can re-run the verification without going back to the bank. The `--deterministic` flag at deployment-package.md §"Operational behavior" gives byte-identical output for the same inputs. The bundle includes the ledger snapshot, the public key, the verifier output, and the methodology section.

I did not see explicit text that says "the bundle is sufficient for a second examiner to re-run the verification without re-engaging the bank," but the contents listed are sufficient for that purpose. Recommend a line in deployment-package.md or sample-report.md that says this explicitly. See Gap 5.2.

### 3.4 The portfolio comparison procedures

`portfolio-comparison-procedures.md` is the document I most expected to be missing. EIC-level cross-bank comparison is the part of FFIEC IT examination that gets the least tooling investment from technology vendors because individual banks do not need it. Regulators do.

The doc covers cross-bank, multi-region within one institution, and cross-period within one institution. The metrics table is the right metrics — pass rate, anomaly rate, sealing-delay frequency, late-binding rate, master-rotation events, verifier-validation cadence. The "Common patterns and their implications" section recognises that cross-bank patterns can indicate a shared vendor issue (HSM provider degradation, cloud region issue) and that the supervisory response is to investigate the shared root cause rather than write up each institution individually. That distinction is exactly the one the FFIEC Cybersecurity Assessment Tool was supposed to surface but did not.

The acknowledgement that manual comparison works up to ~20-50 institutions and the v1.1 `verifier portfolio` subcommand is the planned tooling — that is the right framing. Regulators with very large portfolios will build internal aggregation tools that consume the JSON output; the chain produces deterministic JSON, so aggregation is straightforward. The doc does not over-promise the v1.1 tooling.

### 3.5 Repeat-finding language

`finding-language.md` §"Repeat-finding language" carries the language for when a finding from a prior examination remains open or recurs. Examination report drafting routinely includes repeat findings, and the language for "the institution previously remediated a related finding at the [date] examination; the current examination finds the issue has recurred" is the language we use. Having it pre-drafted reduces the time to write up the recurrence and ensures the severity escalation (next-tier severity / enforcement track) is consistent across examiners.

The "Public-disclosure language" and "Enforcement-action lifecycle" sections are correctly scoped to formal enforcement actions (consent orders, formal agreements, civil money penalty orders). Most chain findings will not reach that level. Having the language for the case where they do means the regulator's enforcement function can adopt the language without re-drafting from scratch.

## 4. Things I want to push back on (small)

### 4.1 The spec §7 step count terminology

`finding-language.md` notes that the table has 14 informational rows covering 12 steps plus file pre-flight plus anomaly rows, and the "12-step procedure" terminology refers to the spec §7 step count. That is fine for someone who has read the spec. For an examiner whose first encounter with the procedure is the failure-mode table, the row count and the step count not matching can be momentarily confusing — they think they are looking at a 14-step procedure.

The note at the bottom of the table explains it. The note works. But a one-line preamble to the table that says "The 12 spec §7 steps map to the rows below; some steps emit multiple distinct reason strings (step 11 dual-algorithm sub-cases, step 12 cadence vs dev-mode), and rows for file-pre-flight and anomaly checks are added alongside" would land the framing before the reader counts rows. Not blocking; reader-experience nit.

### 4.2 The cadence-mismatch escalation path

`finding-language.md` row for step 12 cadence mismatch says "Observation; control-description-accuracy finding; escalate to MRA on repeat." That is the right severity for a first occurrence. The escalation-on-repeat is the right disposition for recurrence.

What is missing is the threshold for "repeat." Does one prior examination with a cadence-mismatch finding constitute repeat? Two? Within how many examination cycles? If the institution's cadence-mismatch finding from three years ago lapsed because the institution updated its control description and then drifted again, is that a repeat or a fresh first-occurrence?

This is a regulator-pack question, not a spec question. The spec is silent because it should be. But the pack could carry one paragraph in finding-language.md or examination-response-workflow.md that says "repeat is defined per the regulator's standard examination-cycle conventions; for OCC, repeat is the same finding within two consecutive examination cycles" or similar. If the answer is "this is up to each regulator," say so explicitly. Right now the reader fills in the answer from their own examination practice and may end up inconsistent across examiners.

### 4.3 The "Bank counsel will ask" voice

The pack has the right severities and the right finding paragraphs. What it does not have is the explicit acknowledgement of how bank counsel will push back on each finding type and what the examiner's response should be.

For example, on step 8 key_fingerprint mismatch, bank counsel will argue that the finding is a documentation issue (the institution rotated the IKM and forgot to update the verifier-side registry, but the chain content is intact and the actual control operated correctly). The examiner's response is that the spec §10.1 weekly reconciliation is a control in its own right, the institution's reconciliation should have caught the rotation discrepancy, and the finding is against the reconciliation control rather than the chain content. That response is substantively in the finding paragraph but the rebuttal-of-counsel framing is not explicit.

Similarly on step 11 case (e) co-signed seal failure under non-strict (PASS-WITH-ANOMALY), bank counsel will argue that the verifier itself disposed the seal as passing and an MRA against a passing seal is unsupported. The examiner's response is that PASS-WITH-ANOMALY reflects downstream-consumer integrity assurance from the un-broken algorithm, not a determination that the failure is low-severity, and the spec text at §7 step 11 case (e) says explicitly that severity is Severe regardless of bracket. That response is in the spec text but is not lifted into the finding-language doc as an explicit "if counsel argues X, the response is Y" paragraph.

I would not block an examination on this. The substantive answer is in the documents. But for examiners new to chain examinations, an explicit "common pushback patterns" section in finding-language.md would shorten the learning curve. See Gap 5.3.

### 4.4 The training doc is a 30-minute session, not a substitute for examination experience

`examiner-training.md` is described as 30 minutes from "I have a ledger snapshot and a public key" to "I have a defensible report." The structure is sound — five modules, hands-on in module 3, sample exercises in `docs/regulator-pack/exercises/` separately distributed. That is the right shape for onboarding a junior examiner who already has IT examination experience.

What it is not is sufficient training for a cybersecurity-track examiner whose normal work is FFIEC CAT-style assessments rather than Handbook-aligned IT examinations. Those examiners will need additional context on the IS booklet II.C.10 / II.C.13 framing, the AIO II.E change-management framing, and how chain output integrates with their normal evidence collection.

This is not a defect in the training doc — the training doc is correctly scoped to "examiner who has read the quickstart and has the verifier binary." It is a gap in the broader pack: there is no document that maps chain output to FFIEC CAT controls, only to Handbook controls. CSF-2.0.md exists, which covers the NIST CSF mapping. CAT does not have its own crosswalk document.

For the OCC, this is fine — we do not use CAT as our primary control framework. For state regulators and some Federal Reserve teams that lean more heavily on CAT, the absence may be more noticeable. Recommend a CAT crosswalk document in the regulator pack if the project intends to support CAT-aligned examinations directly. See Gap 5.4.

## 5. Net-new gaps for this round

### 5.1 Step 12a missing from failure-mode tables

The new step 12a (GenAI model identifier completeness) is in the spec text, in the dual-algorithm-aware finding-language row layout, and in the training doc Module 5 by reference. It is NOT in the failure-mode taxonomy tables at:

- `examiner-quickstart.md` §"Common failure modes (quick reference)"
- `finding-language.md` severity-by-step table at the top
- `examiner-training.md` Module 5 patterns table

Add a row for step 12a to all three tables with:
- Verifier reason string: `gen_ai_model_identifier_missing at seq N: {field_name} required for chain entries representing model calls`
- Severity: PASS-WITH-ANOMALY under non-strict (Observation), FAIL under `--strict` (control-completeness MRA)
- IR scenario: an SR 11-7 model-risk variant (the institution's SR 11-7 program documents how the institution responds to model identifier completeness gaps; the chain check produces the evidence)

Severity rationale: control-completeness for SR 11-7 reproducibility is not chain-integrity; the chain content is sound, the SR 11-7 reproducibility evidence is incomplete on the affected entries.

### 5.2 Bundle sufficiency statement

Add one line to either `regulator-pack/deployment-package.md` or `regulator-pack/sample-report.md` that says:

> The bundle is sufficient for a second examiner to re-verify the institution's chain without re-engaging the institution. The bundle contains the ledger snapshot, the institution's public key, the verifier output, and the methodology section; with the bundle and a verified copy of the verifier binary, a second examiner produces byte-identical output to the original examination.

Substantively the contents already imply this. Saying it explicitly closes the loop on regulator-internal review processes (peer review of an examiner's working paper) and on FOIA-adjacent disclosure considerations where the institution's chain output may be reviewed by the regulator's enforcement function without re-engagement of the bank.

### 5.3 Counsel-rebuttal framing in finding-language.md

Add a section to `regulator-pack/finding-language.md` titled "Common pushback patterns and the examiner response." For each of the high-severity finding types (step 8, step 9, step 10, step 11 single-algo, step 11 case (e) co-signed seal failure, step 12 dev-mode), add 2-3 sentences naming the common counsel argument and the examiner's response.

Format suggestion (one example):

> **Step 8 key_fingerprint mismatch — common pushback.** Counsel may argue that the finding is a documentation issue rather than a control failure, on the grounds that the chain content is intact and the actual chain-of-custody control operated correctly. The examiner response: spec §10.1 weekly key-fingerprint reconciliation is itself a control in the institution's control description; the institution's reconciliation should have surfaced the discrepancy between the IKM-roster row and the per-entry fingerprint before the examination. The finding is against the reconciliation control's effectiveness, not the chain content. The institution's most recent `master.reconciliation_completed` event is the relevant evidence; if the event does not record the affected `(tenant_id, key_version, key_fingerprint)` triple as unmatched, the reconciliation either did not run on the affected period or produced a false negative.

Three to six such paragraphs, one per high-severity finding type. The substance is already in the documents; the rebuttal-of-counsel framing makes it usable at the moment of pushback rather than reconstructed from four other docs.

### 5.4 CAT crosswalk (only if the project intends to support CAT-aligned examinations)

If the project intends to support FFIEC Cybersecurity Assessment Tool aligned examinations (some state regulators and some Fed teams), add a `regulator-pack/cat-mapping.md` document with the same shape as `handbook-mapping.md` but keyed to CAT domain assessment factors and component statements. If the project's intent is to support Handbook-aligned and CSF-aligned examinations only, document that scope explicitly in the regulator-pack README so reviewers do not infer a CAT crosswalk that is not coming.

This is the lowest-priority gap of the four. Many examiners (myself included) work primarily from the Handbook and CSF. CAT is a self-assessment tool more than an examination tool in current OCC practice. The gap is real but small.

### 5.5 Cadence-mismatch repeat-threshold definition

Add one paragraph to `regulator-pack/finding-language.md` §"Cadence mismatch (spec §7 step 12)" or to the closing-language section that defines what "repeat" means for the escalation-on-repeat disposition. If the answer is "per regulator standard examination-cycle conventions," say so explicitly so examiners do not infer a project-specific definition.

Suggested wording:

> "Repeat" for the escalation-on-repeat disposition is defined per the regulator's standard examination-cycle conventions. For most regulators this is the same finding within two consecutive examination cycles; the project does not impose a definition, and examiners apply their regulator's convention.

Smallest of the five gaps; reader-experience improvement, not a defect.

## 6. Stopping criterion check

The doc set is examinable today. The five gaps above are net-new for round 14 and are all improvements rather than blockers. The headline workflow (JSON record → severity → IR scenario → evidence → finding paragraph) is intact, the dual-algorithm transitional treatment is sound, the step 12a addition is correctly scoped, the portfolio comparison procedures are written and acknowledged as manual until v1.1 ships, and the failure-mode taxonomy is consistent across the quickstart, training, finding-language, and workflow documents (modulo the step 12a row I want added to three tables).

Net-new gaps: 5 (all improvements; none block examination).
Net-new defects: 0.

The stopping criterion 0/0 is not yet met because of the five net-new gaps. None are blocking; the most significant is gap 5.1 (step 12a row in failure-mode tables) because it affects the core failure-mode taxonomy. The other four are reader-experience improvements (5.2 bundle sufficiency, 5.5 repeat threshold), counsel-rebuttal framing (5.3), or scope-clarification (5.4 CAT crosswalk).

I would expect round 15 to close gap 5.1 and 5.2 trivially, gap 5.3 with 2-3 hours of writing, gap 5.5 with one paragraph, and gap 5.4 to be addressed by a scope statement rather than a new document. After those five, I would not have further substantive gaps; the package would be ready to support OCC IT examinations without further regulator-pack iteration.

— Brendan Walsh, OCC IT Examination
