# Round 16 review — FFIEC IT Examiner

**Reviewer.** Marcus Reynolds, Senior FFIEC IT Examiner, Federal Reserve Bank — Dallas district. 17 years on the IT examination beat across SIFIs, mid-tier regional banks, and a handful of cross-agency holding-company engagements. I have run plenty of independent log-integrity reviews over the years where the institution's own logging stack was the only thing the examiner had to lean on; that is the lens I am bringing here.

**Documents read.** spec/chain-of-custody-v1.md §7; examiner-quickstart.md; regulator-pack/sample-report.md; regulator-pack/finding-language.md; regulator-pack/handbook-mapping.md; regulator-pack/deployment-package.md; regulator-pack/examiner-training.md; regulator-pack/examination-response-workflow.md; portfolio-comparison-procedures.md.

**Stopping criterion.** 0 gaps and 0 partials.

---

## Headline

I can take the verifier into an examination tomorrow. The five-minute quickstart, the thirty-minute training, the response workflow, and the finding-language paragraphs together get an examiner from "I have never seen this thing before" to "I have a defensible report and a working-paper bundle" inside one engagement-prep day.

The pieces that matter to me as an IT examiner are all here:

- A precise twelve-step procedure in the spec (§7) with named failure modes per step.
- A response workflow that keys off the verifier's `step` field — one lookup, one severity row, one IR scenario, one finding paragraph.
- A finding-language doc keyed by the same `step` field with paragraph-ready language for each failure mode.
- A handbook-mapping doc that ties the chain to specific IT Examination Handbook control objectives — II.C.10, II.C.13, II.E — with enough depth that I can defend the mapping to my supervisor or to OCC counterparts on a joint engagement.
- A sample report I can hand to a junior examiner as the reference for what they should expect to see.
- A portfolio-comparison procedure that works manually today and articulates what the v1.1 tooling will automate.

I tried to find the gap. I could not.

## What I checked, in examination-prep order

### 1. The verifier's failure taxonomy is precise and addressable

A common pattern in vendor-supplied verification tools is a generic "verification failed" output that pushes the diagnostic burden back onto the examiner. The spec §7 procedure is the opposite shape. Every failure carries a `step` integer, a precise reason string, and the `(run_id, seq)` or `seal_date` that scopes the failure. That precision is what makes the examination-response-workflow's five-step path possible at all — without the `step` field the examiner would be reconstructing the workflow from prose at the moment of decision, which is what the doc explicitly avoids.

The taxonomy I checked against my own examination experience:

- **Identity vs content separation (step 8 vs step 9).** This is the distinction that examiners-in-the-field most often get wrong on first encounter. The doc-set handles it three times — the quickstart's "stop and call the bank" callout, the finding-language paragraph for step 8 emphasizing it is investigated against the IKM roster (not the chain content), and the workflow's step 3 evidence path pointing at `master.reconciliation_completed`. The repetition is intentional and correct; an examiner reading any one of the three docs at the moment of a step-8 failure lands on the right investigation path.
- **Operational vs integrity separation (file pre-flight vs steps 6/9/10).** The mid-write truncation refusal explicitly carries "operational, NOT tampering" wording in three places (quickstart, finding-language, training). On a busy examination day this matters — the examiner who sees `audit file ends mid-line` on a Friday afternoon needs to reach the right severity bucket without reading the spec.
- **Verifier-version skew vs institutional finding (step 1).** Calling out "this is not an institutional finding" in three places (quickstart, finding-language, training, workflow's "When this workflow does NOT apply" section) is correct. I have seen examiners write up vendor-side version skew as institutional findings because the tool said "FAIL" and the examiner did not have the vocabulary to push back.
- **Dual-algorithm transitional period (step 11 sub-cases (a) through (e)).** The sub-case-(e) load-bearing call-out — "Severe regardless of bracket" — is exactly the language a Severity-determination defender needs. Counsel will rebut that the un-broken algorithm's signature provides integrity assurance and therefore the finding should be Observation; the spec's pre-framing in §7 step 11 case (e) plus the finding-language counsel-rebuttal section give the examiner the response without scrambling.

### 2. The workflow doc is the doc I would actually use at the moment of finding evaluation

The five-step path in `examination-response-workflow.md` is keyed by the verifier's `step` field — not by the institution's name, not by the date, not by the failure description. That is the right key. An examiner looking at a JSON failure record has the `step` integer in front of them; everything else flows from that.

The two worked examples (step 8 key_fingerprint mismatch and step 10 Merkle root mismatch) cover structurally different failure shapes — identity mismatch vs ledger-content mismatch — which demonstrates that the workflow generalizes. A junior examiner reading those two examples can apply the same five-step path to step 9 (content tampering) or step 11 (signature compromise) without further training.

The "When this workflow does NOT apply" section is the part that examiners reading procedures most often skip and most often need. PASS-with-anomaly results are not findings; calling that out explicitly prevents over-writeup at the report-drafting stage.

### 3. The handbook mapping defends the chain against the right control objectives

The Information Security booklet mapping is the one I cared about most. II.C.10 (Logging) is the headline; the mapping calls it out explicitly and adds the weekly key-fingerprint reconciliation per spec §10.1 as the audit-evidence artifact. That is correct — institutions claiming II.C.10 compliance via this chain need to produce the `master.reconciliation_completed` operational event as the periodic evidence, not just the chain itself.

II.C.13 (Cryptographic controls) gets the depth treatment with three rework primitives — per-entry `key_fingerprint` as a public identity binding, IKM minimum 32 bytes per RFC 4868 keying floor, and software-key adapter compile-time exclusion. That is the right level for a cryptographic-controls examiner working from the booklet; the language gives me the four specific items to confirm at examination time (IKM length, valid `key_fingerprint` bytes verified at step 8, no production entries with `kms_handle_uri = "plaintext-dev"`, weekly reconciliation cadence operating).

The AIO booklet's II.E (Change management) mapping correctly treats `format_version` as the load-bearing change-management primitive, plus the dual-algorithm transitional-period algorithm-posture change-management addition. That is forward-looking; institutions adopting post-quantum will need exactly this examination posture, and having it documented now means the change-management examination step is ready when it is needed rather than being constructed during the first PQC migration.

The Audit booklet and OTS booklet mappings are short and correct — the verifier is the audit artifact, the chain is topology-neutral so vendor-hosted implementations get examined through the same verifier.

### 4. The sample report and the deployment package are usable artifacts

The sample report's per-day detail block is exactly the shape a working-paper-friendly report takes. The "Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12" line is the diagnostic that lets an examiner confirm every defense-in-depth step actually ran on that day's data. The paired PASS-day-then-FAIL-day excerpt at the end is the right teaching device — it shows the reader what the same report shape looks like under both outcomes without making them imagine.

The SOC team appendix mapping verifier-output lines to TSC criteria is helpful for the joint-engagement scenario where the FFIEC examination happens alongside a SOC 2 audit at the same institution. I do not work SOC engagements directly, but I have shared workspaces with SOC engagement teams enough times to value having the cross-walk in hand rather than having to construct it during the joint review.

The deployment package is the IT-shop-ready document. The verifier's operational behavior — no telemetry, no network access, read-only ledger access, deterministic output — is exactly what my IT shop needs to see before allowlisting a binary onto examiner laptops. The two acceptable allowlisting postures (SHA-256 hash vs cosign signature) match the two postures my IT shop currently runs for other examination tools. The 30-day pre-announcement window is reasonable; the "passing report from older verifier remains valid" guarantee is the operational continuity statement my IT shop will look for first.

The validate-before-run script is the right control point. An examiner running an unsigned or tampered verifier binary against an institution's snapshot is a finding-the-finding scenario; the wrapper makes the validation step a reflex rather than a remembered procedure.

### 5. The training doc is well-paced for the 30-minute target

The five modules in `examiner-training.md` (concepts → inputs → run → read → patterns) match how I would actually train a new IT examiner on a verification tool. The Module 1 three-layer concept (HMAC chain → daily Merkle seal → Ed25519 in HSM custody) is the right depth for a new examiner — enough to defend the report's conclusions without diving into the cryptographic primitives.

The Module 5 patterns table covers the failure modes a new examiner is most likely to encounter. The "stop and call the bank" callout for step 8 is the same callout the quickstart has, which is correct — across documents, a Severe finding gets the same handling instruction.

The "Beyond the 30 minutes" topics list (`verifier walk`, `verifier diff`, `verifier consolidate`, NIST CSF mapping, FFIEC Handbook mapping) maps to the next training session a returning examiner would take after their first engagement. That is the right scope for a 30-minute onboarding.

The sample exercises set (clean-pass, failure-mode, anomaly-only, multi-region) is the right shape for sandbox practice. The doc notes the exercises are separately distributed; on first read I would have wanted them inline in the doc, but distributing them separately makes operational sense — the sample exercises will rev faster than the training doc, and bundling them into the training doc would create version-skew between the two.

### 6. Portfolio comparison covers the EIC role appropriately

The portfolio-comparison-procedures doc covers cross-bank, multi-region, and cross-period comparison patterns. The metrics table (pass rate, anomaly rate, sealing-delay frequency, late-binding rate, master-rotation events, verifier-validation cadence) is the right summary set for an EIC working a portfolio.

The "common patterns and their implications" section calls out the three patterns that materially drive supervisory action — sub-portfolio shared anomaly (vendor-side investigation), institution-differs-materially (deeper examination at next cycle), portfolio-wide degradation (industry-level escalation). Those are the three patterns I have actually seen play out in supervisory practice; having them documented as procedures rather than instinct is the right move for a younger EIC.

The "manual today, v1.1 tooling planned" framing for the portfolio subcommand is honest. The procedure as written produces equivalent comparison data manually for portfolios up to 20-50 institutions; that covers the Dallas district's typical SIFI + mid-tier portfolio scope. A regulator with a much larger portfolio gets the JSON output and builds internal aggregation; the chain produces deterministic JSON, so the aggregation is straightforward.

The cross-agency coordination section correctly identifies that the verifier output is the same across agencies and is shared through standard supervisory MOUs. For holding-company examinations under Fed + bank-subsidiary examinations under OCC, that is the right operational shape — each agency conducts its own examination using the same verifier output as the substrate.

## Specific items I deliberately checked for and found

- **Severity guidance per failure mode.** Every `step` in spec §7 has a severity row in `finding-language.md`. PASS-with-anomaly cases are explicitly distinguished from MRA-track findings. The dual-algorithm sub-cases get their own severity rows.
- **Counsel-rebuttal pre-framing.** The four most-rebutted findings (step 8 key_fingerprint, step 9 payload_hash MAC, step 11 signature, IR Scenario 11 cold-DR attestation) have prepared examiner responses. Pre-framing is preparation, not script — that framing is correct.
- **Repeat-finding language.** The "had a related finding at the [date] examination" template is exactly the language I would use; the escalation tier (next-tier severity / enforcement track) is the right hook for repeat findings.
- **Public-disclosure language.** The consent-order / formal-agreement / civil-money-penalty paragraphs are concise enough to flow into public-disclosure documents without further editing. The technical detail correctly stays in the underlying examination report.
- **Enforcement-action lifecycle.** The initial-action / monitoring / removal-of-action triple is the right shape for an enforcement-track finding involving the chain. The "verifier output for [N consecutive months] shows pass results" criterion for removal is the right substantive criterion.
- **Bundle sufficiency.** The quickstart's "bundle is sufficient for a second examiner to re-verify the institution's chain output independently, without re-engaging the institution" guarantee is the criterion that makes the bundle a defensible working-paper artifact for cross-examiner re-verification — including for examination-quality-review purposes within the regulator.
- **Strict mode posture.** Substantive SOC engagements use `--strict`; FFIEC examiner work uses non-strict and evaluates anomalies in operational context. That is the right posture split — substantive testing wants the conservative refusal; institutional examination wants the operational context to land in the report.
- **CAT crosswalk scope note.** The finding-language doc explicitly scopes itself to FFIEC IT Examination Handbook examinations and notes that CAT-track examiners use the doc as a reference for chain-detected events but follow CAT's own conventions otherwise. That is the right scope discipline; pretending a single document covers both tracks would invite over-claiming.
- **SR 11-7 reproducibility (step 12a).** The new GenAI model identifier completeness check correctly carries Medium severity under non-strict and FAIL under `--strict`, with the explicit "control-completeness for SR 11-7 reproducibility, NOT chain-integrity" framing. That keeps the model-risk-program finding distinct from the chain-integrity finding pile.

## Items where I considered raising concern and decided not to

- **Step 11 dual-algorithm sub-case (e) "Severe regardless of bracket."** I considered whether the PASS-WITH-ANOMALY disposition under non-strict would let an examiner under-write the finding by treating the disposition label as the severity. The finding-language row explicitly cites "row 11 (dual-algo) co-signed seal failure (Severe MRA) regardless of which bracket the verifier reported." The pre-framing in spec §7 step 11 case (e) reinforces that. An examiner who reads either doc lands on the right severity. Concern resolved.
- **The "12-step procedure" terminology when the finding-language table has 14 informational rows.** I noted this and the doc itself addresses it: "The table has 14 informational rows covering the 12 spec §7 steps PLUS the file-pre-flight truncation refusal PLUS the anomaly-only rows... some steps emit multiple distinct reason strings, and the anomaly rows are not §7 steps." That note is the right place for the explanation. Concern resolved.
- **Verifier-version skew (step 1) being labelled "not an institutional finding" might be read as "not a finding at all."** I considered whether an examiner could close out a step-1 verifier output as no-action when the appropriate response is "obtain newer verifier and re-run." The finding-language paragraph for step 1 explicitly says "the examiner resolves the finding by obtaining the newer verifier from the project supply chain (per `regulator-pack/deployment-package.md`) and re-running. No examination action against the institution is required." That language directs the action correctly. Concern resolved.
- **The portfolio-comparison doc's "manual today, v1.1 automated" framing.** I considered whether the manual-today posture would produce comparison data of materially worse quality than the v1.1 tooling will. The metrics defined manually are the same metrics the v1.1 subcommand will compute; the JSON output is deterministic; aggregation is straightforward. The manual procedure produces equivalent data for the portfolio scopes that matter today. Concern resolved.
- **Cross-agency coordination on holding companies.** I considered whether the doc-set adequately addresses the Fed + OCC joint-supervision case that comes up routinely in the Dallas district. The portfolio-comparison doc's cross-agency section calls out exactly this case and routes through standard supervisory MOUs. The verifier output being identical across agencies is the load-bearing property. Concern resolved.
- **The IKM disclosure shape.** I considered whether the `--master-key` provisioning at examination time creates an awkward operational dependency on the institution. Three docs (quickstart, sample-report, deployment-package, training) consistently route the IKM disclosure through `legal-disclosure.md` §"Court-ordered master-key disclosure" and `customer-dispute-procedures.md` §"IKM access for customer-side verification." That is the right escalation shape — the protective-order or HSM-mediated derivation path is the documented control. The structural-only fallback under `--master-key absent` produces PASS-WITH-ANOMALY rather than silent partial verification, and `--strict` elevates to FAIL. The fail-closed posture is correct. Concern resolved.

## Final disposition

**0 gaps. 0 partials.**

For an IT-examination program adopting this chain as a control input under FFIEC IT Examination Handbook authority, the regulator-pack documents above are sufficient to:

1. Allowlist and deploy the verifier on examiner laptops (`deployment-package.md`).
2. Onboard new examiners to the verification workflow in 30 minutes (`examiner-training.md`).
3. Produce defensible examination reports keyed off verifier output (`sample-report.md`, `finding-language.md`).
4. Walk from a JSON failure record to a documented examination response in 60 seconds (`examination-response-workflow.md`).
5. Defend the chain's mapping to specific IT Examination Handbook control objectives (`handbook-mapping.md`).
6. Conduct portfolio-level comparison across institutions and across multi-region deployments at one institution (`portfolio-comparison-procedures.md`).

The five-minute quickstart serves as the engagement-prep entry point and is consistent with the deeper docs.

I would take this into an examination tomorrow without further preparation.

— Marcus Reynolds, Senior FFIEC IT Examiner, FRB Dallas
