# Round-12 review — FFIEC IT Examiner

**Reviewer:** Janet Buckley, Senior FFIEC IT Examiner, FDIC
**Background:** 15 years on the IT examination beat; before that, four years as IT director at a community bank in Ohio. I run between two and four full-scope IT examinations a quarter and I read examiner-facing docs from the bench, not from a desk.
**Reading angle:** examiner workflow. Quickstart, training, sample report, finding language, IR-scenario mapping, and the response workflow. Can a working examiner pick this up cold and produce a defensible report?
**Date of review:** 2026-05-06.

---

## Headline

This material is at the level I'd accept on a national IT examination program. The five docs that matter to the bench examiner — quickstart, training, sample report, finding language, response workflow — read as a single coherent loop. A reasonably competent IT examiner can land the verifier on a laptop, run it against a bank's snapshot, identify a failure, look up the severity, lift the finding language, and document the response, without leaving the regulator-pack directory.

I want to call out three structural choices that materially help the bench:

- **The `step` field is the spine of the whole loop.** The verifier emits a numbered spec §7 step on every failure; the quickstart's table, the finding-language severity table, the training Module 5 table, and the response workflow all key off the same number. One number, one row, one finding paragraph, one IR scenario. That is examiner ergonomics done right.
- **The IR-scenario column is now in three places consistently.** Quickstart table, finding-language severity table, training Module 5. An examiner who sees `step: 8` in a JSON record is one column away from "IR Scenario 7" in any of the three docs. That redundancy is correct — examiners read the doc that's open, not the doc the author thinks is canonical.
- **The 2026-04-22 PASS / 2026-04-23 FAIL pairing in the sample report is the right pedagogy.** New examiners need to see the FAIL output sitting next to a PASS output from the same examination period, in the same report shape. Anchoring both days inside the April examination range makes "this could happen on any examination day" land in a way a hypothetical box at the end of the document never would.

The rest of the review walks through specific docs, calls out what works, and lists narrow follow-ups I'd want to see addressed. Severity grades are EXAMINER, MINOR, NOTE.

---

## Examiner-quickstart.md

The five-minute orientation does the job. I can hand this to a junior examiner the morning of an examination and they will be able to run the verifier by lunch.

### What works

- **The "what you do" block is a complete copy-pasteable invocation.** Including `--master-key` in the standard invocation (not as an optional add-on) is the right default — for FFIEC examiner work, the institution provides the IKM under a documented disclosure shape, and the verifier's structural-only fallback is for SOC engagements and customer-side disputes, not for examiner work. Showing the full invocation with `--master-key` means the examiner doesn't have to read the fail-closed-semantics paragraph in the spec to know the right shape for their workflow.
- **The "stop and call the bank" framing is clear and bounded.** Two specific failure modes (`step: 8` `key_fingerprint mismatch` and `step: 9` `payload_hash MAC mismatch`) get the framing. Both are Severe, both warrant immediate institution contact, and the doc spells out that the investigation paths differ — step 8 goes against the IKM roster, step 9 goes against the chain content. That distinction matters because an examiner who treats step 8 as a content-tampering finding will go down the wrong investigation path and the institution will be confused about what evidence to produce. This is the single most important pedagogical line in the whole quickstart.
- **The `audit file ends mid-line` row is correctly bracketed as operational, not tampering.** A new examiner's instinct on "the file ends in the middle of an entry" is "tampering." The doc heads that off explicitly. Good defensive writing.
- **The 14-row table covers every spec §7 failure mode plus the file pre-flight.** The IR scenario column on the right closes the loop to the institution-side response. An examiner reading the table can answer "what does this mean and what does the bank do about it" without leaving the row.

### Narrow follow-ups

- **NOTE — IKM file-mode reference.** The `--master-key` paragraph says the IKM file is "32 raw bytes, file mode 0600." On a Windows examiner laptop that descriptor doesn't translate. I'd add one parenthetical: "(file mode 0600 on POSIX; equivalent ACL restricting read access to the examiner's user on Windows)." Not a finding; a clarity edit.
- **NOTE — `--master-key absent` row.** The quickstart's invocation shows `--master-key`, but the table doesn't have a row for "examiner ran without `--master-key` and got `structurally consistent, key-bound verification skipped`." That's a real outcome an examiner will hit if they forget the flag, and the result text is technically a non-finding, but a row in the table would short-circuit the "did I do something wrong?" anxiety. Cross-references the spec §7 fail-closed paragraph.

---

## Sample report (regulator-pack/sample-report.md)

The PASS-with-anomalies report and the paired PASS/FAIL example are the right two artifacts for examiner training. The TSC appendix is the right artifact for the SOC team consuming the same report.

### What works

- **The 2026-04-22 PASS day next to the 2026-04-23 FAIL day, both inside the April examination range.** This is the load-bearing pedagogical decision. New examiners need to see what a real failure looks like inside an otherwise-healthy month. Anchoring the FAIL on April 23 (next to a documented PASS on April 22) makes the failure feel like "this could happen on any day in your examination" rather than "this is a contrived edge case." The visual continuity in the per-day detail block — same field shape, same step list, same anomaly column — is exactly what the bench examiner needs.
- **`--master-key` is in the standard invocation.** The sample report's invocation matches the quickstart's: `--master-key ./acme-ikm.bin` is shown as standard. The cross-reference paragraph (institution provides under a documented disclosure shape) is right. The SOC variant correctly drops `--master-key` and adds `--strict`, which matches the customer-side-dispute and legal-disclosure docs the cross-references point to.
- **The FAIL day shows `Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)`.** Showing how far through the procedure the verifier got is excellent. The examiner doesn't have to wonder "did the verifier check the Merkle? did it check the signature?" — the field tells them. And `Merkle match: N/A (chain walk failed before reaching step 10)` makes the early-stop semantics explicit.
- **The expanded TSC appendix.** The cross-walk table now covers every spec §7 step and every operational outcome. For the SOC engagement team this is gold — they can allocate verifier evidence to TSC criteria without inventing the mapping themselves. The PI1.1 / PI1.2 / CC6.7 / CC6.8 / CC8.1 mapping is consistent with how the Big Four allocate substantive evidence behind processing-integrity criteria. Particular call-out: the `step: 4 cross-chain lift` row maps to CC6.7 + CC9.2 (vendor management — evidence-handling), which captures the "mis-bundled snapshot" branch of the failure correctly without prejudging it as integrity tampering.
- **The methodology section matches the spec §7 procedure exactly.** Steps 1-12 in the order the verifier executes them; pass/fail rules block at the bottom. An auditor or examiner can reconstruct what the verifier did from this section without reading the spec.

### Narrow follow-ups

- **MINOR — Cover page re-statement on the FAIL example.** The "Cover page summary changes accordingly" block at the end of the FAIL example is correct, but it sits below the per-day FAIL detail rather than above it. An examiner reading the doc top-to-bottom sees the FAIL detail before they see the cover-page rollup. A bench examiner reads cover page first, then per-day. Re-ordering the FAIL example to (a) cover page summary first, then (b) PASS day for context, then (c) FAIL day with failure record, would mirror the actual reading order. Optional reorganization, not a content gap.
- **NOTE — JSON snippet pairing.** The sample report's FAILURE RECORD shows the field names (`step`, `reason`, `run_id`, `seq`, `tenant_id`, `key_version`, `expected_fingerprint`, `recorded_fingerprint`) but not as JSON. The examination-response-workflow.md doc shows the same record AS JSON. An examiner moving between the two docs has to mentally translate "the report's pretty-print" into "the JSON the workflow doc describes." Optional addition: one line in the sample-report indicating that `--json-report` produces the same record in JSON form (and pointing at examination-response-workflow.md for the JSON shape). Cross-link, not a content change.

---

## Finding-language.md

This is the document I'd lift from on the bench. It now carries five new finding paragraphs (steps 2, 3, 5, 6, 9) plus the existing paragraphs for steps 7, 8, 10, 11, 12, the file-pre-flight truncation, the cross-chain lift, and the dual-algorithm cases. The severity table at the top with the IR scenario column on the right is the right layout.

### What works

- **The five new finding paragraphs cover the previously-unaddressed steps.** Step 2 (HKDF inputs mismatch), step 3 (genesis_hash mismatch), step 5 (format_version mismatch at entry), step 6 (chain link broken), and step 9 (payload_hash MAC mismatch) all now have lift-and-edit paragraphs. The bench examiner no longer has to invent finding language for these failure modes — they're written, edited for clarity, and in the same voice as the existing paragraphs. The step 2 paragraph in particular does the right thing by enumerating the three plausible root causes (institution SDK constants, verifier-build defect, header tampering) so the examiner doesn't pre-judge the institution.
- **The severity table's IR scenario column.** Now in finding-language.md too, matching quickstart and training. An examiner working from finding-language.md can pick the row, lift the severity guidance, lift the finding paragraph, and direct the institution to the matching IR scenario without leaving the doc.
- **The three dual-algorithm reason rows.** The table now carries three rows for spec §7 step 11 in dual-algorithm posture: partial-coverage seal, algorithm-not-on-posture-list, and co-signed seal failure. The three rows correctly distinguish severity by case: partial-coverage and algorithm-not-on-posture-list are control-completeness/control-description findings (Observation under non-strict, MRA under --strict), while co-signed seal failure is a Severe MRA in either mode. Bracketing the first two as control findings rather than chain-integrity findings is the right call — a partial-coverage seal is integrity-bearing under the present algorithm; the institution's posture commitment is what's incomplete.
- **The "identity mismatch, not content tampering" emphasis on step 8.** The key_fingerprint mismatch paragraph spells out, in bold, that this is investigated against the IKM roster and the §10.1 reconciliation log, NOT against chain content. This is the most likely-to-be-confused finding mode in the whole spec, and the paragraph defends against the confusion explicitly. The §7-step-8-detected-before-MAC-compute reassurance ("the affected events' MAC integrity is independently confirmable once the IKM is restored") is the right closing — the institution and the examiner both leave the conversation knowing the chain content is not under suspicion.
- **The 14-row table count clarification.** The paragraph explaining why the table has 14 rows for a 12-step spec procedure (some steps emit multiple reason strings; anomaly rows aren't §7 steps) closes the obvious question a careful reader will ask. Defensive writing in the right place.

### Narrow follow-ups

- **MINOR — Step 9 paragraph and "[is/is not] underway" template field.** The new step 9 finding paragraph contains "Bank management's investigation [is/is not] underway." That bracket-template style is unusual in this doc — the other finding paragraphs use `[N]`, `[date]`, `[X]` for substituted values, not "[is/is not]" for a binary state. Suggest harmonizing to: "Bank management's investigation is [underway / has not been initiated]." Same semantics, consistent template style.
- **NOTE — Step 6 paragraph and "the most likely tampering signal that does not require key access" framing.** The step 6 paragraph (chain link broken) describes the failure as "the most likely tampering signal that does not require key access." That's accurate — an attacker who can write to the ledger but doesn't have the IKM can produce a step 6 failure but not a step 9 failure. But "the most likely" is a conclusion the examiner shouldn't write into a finding without thinking through the institution's specific facts. Suggest softer framing: "consistent with insertion or deletion of chain entries by an actor who does not have access to the IKM, including some classes of operational defect at the storage layer." Keeps the substance, removes the implicit attribution.
- **NOTE — Repeat-finding language and the new failure modes.** The repeat-finding language section at the bottom is generic, which is correct. But the new step 2/3/5/6/9 paragraphs don't yet have explicit examples of repeat-finding language. An examiner working a second-cycle examination on a step 2 finding would adapt the generic language correctly, so this isn't a gap, but a one-line example per new paragraph ("If the institution had a related finding on [prior date]: ...") would close the loop. Optional.

---

## Examination-response-workflow.md

This is the document that turns a JSON failure record into a documented examination response. ~60 seconds of reading is exactly right for the bench examiner who is sitting in front of a failure and needs to act.

### What works

- **The five-step path is the correct shape.** Map step → severity (+ finding paragraph), map step → IR scenario, pull the institution's reconciliation/seal-job/HSM evidence, lift the finding paragraph, confirm the institution's response. Five steps, each one cross-references a specific source doc, none of them require the examiner to assemble the workflow from scratch.
- **The dual-algorithm sub-cases are now in the step-mapping table.** Three sub-rows under step 11: partial-coverage seal, algorithm-not-on-declared-posture-list, and co-signed seal failure with one valid + one invalid. Each maps to a specific IR scenario (Scenario 5 variant for the first two, Scenario 4 + Scenario 3 for the third). The bench examiner doesn't have to read the spec §7 step 11 cases (a)-(e) to know what to do — the table has the dispatch.
- **The second worked example (step 10 Merkle root mismatch).** Walking the same five-step path through a structurally different failure type confirms the workflow generalizes. The "computed_root vs recorded_root mismatch direction tells the auditor whether the seal record was altered or the ledger contents were altered" paragraph is the right level of forensic specificity for the bench — concrete enough to produce a question the institution can answer, not so deep that the examiner has to be a cryptographer.
- **The "when this workflow does NOT apply" section.** Three exclusions (step 1 verifier-version skew, PASS-with-anomaly results, PASS results) keep the examiner from over-applying the workflow. The PASS-with-anomaly bracket ("they land in the examination report under 'Anomalies noted' with a sentence each") is exactly how examiners actually write up operational signals.

### Narrow follow-ups

- **NOTE — Step-mapping table alignment with spec §7 cases (a)-(e).** The spec §7 step 11 dual-algorithm dispatch defines five cases (a) both valid, (b) single-algorithm-during-dual-posture (partial-coverage), (c) algorithm-not-on-list, (d) single-algorithm-posture default reduction, (e) co-signed failure with one valid + one invalid. The workflow doc's step-mapping table covers (b), (c), and (e) with explicit rows but omits (a) and (d) because they're PASS / PASS reductions. That's correct — the workflow handles failures, not passes — but a one-line note above the table ("Cases (a) and (d) of spec §7 step 11 are PASS results and do not require this workflow; see the 'when this workflow does NOT apply' section") would close the cross-reference loop. Optional.
- **NOTE — Case (e) and the regulator-coordination phrase.** The case (e) row in the step-mapping table says "regulator coordination on migration timeline if attributable to a published algorithm break." That's correct. But the IR scenario column says "Scenario 4 (master-key compromise on the failed algorithm) + Scenario 3 (signature failure path)." A bench examiner reading the row may not immediately register that case (e) is a TWO-scenario response (one scenario per algorithm), not one. Suggest a parenthetical in the IR-scenario column: "(both scenarios apply because both algorithms must be evaluated; the un-broken algorithm's signature still provides integrity assurance under Scenario 3, while the broken algorithm requires Scenario 4 evaluation)." Adds a sentence; clarifies the two-track response.

---

## Examiner-training.md (Module 5 step-keyed taxonomy refresh)

The 30-minute training session structure is the right shape for new-examiner onboarding. Module 5 is the one I focused on.

### What works

- **Module 5 is now keyed off the JSON `step` field.** The table's left column is "Pattern (JSON `step` + `reason`)" — that matches the JSON record the verifier produces and matches the keying in finding-language.md and the response workflow. A new examiner finishing Module 5 has the same mental model the response workflow expects, so there's no translation step at the moment of failure evaluation.
- **The dual-algorithm row.** The `step: 11 (dual-algo)` `co-signed seal failure` row is now in Module 5, with the correct severity ("Severe; coordinate with regulator on migration timeline") and the correct IR scenario reference ("IR Scenario 4 + Scenario 3"). New examiners learn the dual-algorithm case as part of basic training, not as advanced material. That's the right call given the multi-year transitional period the spec contemplates.
- **The "stop and call the bank" framing in the step 8 row.** Training Module 5 explicitly carries the framing, and the step 8 row's "investigate against IKM roster (NOT chain content)" parenthetical reinforces the identity-vs-content distinction. New examiners absorb the distinction in training rather than having to learn it the hard way on a real examination.
- **The "Steps 1-6 (header pre-flight + cross-chain lift + structural walk) are the format-construction and structural defenses" closing paragraph.** This bracket tells the new examiner that the table doesn't enumerate every step — only the most-likely failure modes — and points them at the full taxonomy in quickstart and finding-language. Defensive writing for "what about steps 1, 2, 3, 4, 5, 6?"

### Narrow follow-ups

- **NOTE — Module 5 row for `step: 11 (dual-algo)` partial-coverage seal.** Module 5 has the row for `co-signed seal failure` (the most severe dual-algorithm case) but doesn't have a row for the partial-coverage or algorithm-not-on-posture-list cases. Those are Observation/MRA-under-strict findings — less severe but more common during the transitional period. A new examiner trained only on the most-severe case will not recognize the others. Optional addition: one row for partial-coverage ("`step: 11 (dual-algo)` `partial-coverage seal` — Observation/MRA-under-strict; control-completeness, NOT chain-integrity"). Keeps Module 5's "common patterns" framing while extending coverage to the realistic transitional-period cases.
- **NOTE — Module 5 row count vs quickstart row count.** Quickstart has 14 rows; finding-language has 14 rows; Module 5 currently has 11 rows (the most-likely patterns). The discrepancy is intentional ("common patterns" is a subset) but a one-line note above the Module 5 table ("This table covers the most common patterns; the full 14-row taxonomy is in `examiner-quickstart.md` and `finding-language.md`") would head off the confusion. Optional.

---

## Handbook-mapping.md (II.E algorithm change management)

The Handbook mapping is the doc I take into a section-by-section examination. The IS booklet headline mapping is II.C.10 (Logging) and that hasn't changed; the AIO booklet mapping now carries an expanded II.E (Change management) cell that addresses algorithm change management explicitly.

### What works

- **The II.E cell covers both `format_version` change management AND algorithm change management.** That's the right scope. An institution adopting a post-quantum algorithm alongside Ed25519 is doing change management on the seal record's algorithm field, the signatures list, AND the institution's declared algorithm-posture configuration. Naming all three in the same cell tells the examiner what to look for in the institution's change-management documentation.
- **The cell explicitly cross-references spec §7 step 11 cases (a)-(e).** An examiner working II.E on a dual-algorithm-adopting institution can read the cell, look up the spec §7 step 11 cases, and ask the institution "show me how your change-management process handles each of cases (a)-(e)." That's the right examination question.
- **The cell calls out IR program coverage.** "The institution's IR program has documented response procedures for spec §7 step 11 dual-algorithm failure modes" is what an examiner should be confirming. Pointing at the IR program (not just the change-management process) recognizes that algorithm posture is a multi-team operation.
- **The "multi-year migration timelines" cross-reference to `09-threat-model.md` §2.8.** An examiner who needs the threat-model context for the migration timeline knows where to find it. The cross-reference is exactly the right specificity — section number, not just doc name.

### Narrow follow-ups

- **NOTE — II.E and the "examiners confirm" list.** The II.E cell ends with three things examiners should confirm. Suggest adding a fourth: "the institution's verifier validation procedure exercises the dual-algorithm dispatch (i.e., the institution's internal-audit or chain-operations team has run the verifier against a co-signed seal during testing, confirming the verifier reports BOTH algorithms' validation results per spec §7 step 11 working-paper convention)." This closes the loop on whether the institution has tested the dispatch path before relying on it in production. Cross-references the spec §7 working-paper convention paragraph.

---

## Deployment-package.md

The deployment-package doc is for the regulator's IT shop, not the bench examiner. I read it from the angle of "what does the IT shop need to allowlist this binary?" That work has been done.

### What works

- **The binary specifications table.** Static binary, no dynamic linking, no network behavior, no privileged operations, runs as ordinary user. Every line is something an IT shop's allowlist process needs to know. The "verifier makes no outbound connections" line is the one that closes the most IT-shop questions in advance.
- **Two acceptable allowlist postures (SHA-256 vs cosign signature).** Recognizing that some regulators allowlist by hash and some by signature, and naming both as conformant, mirrors how IT shops actually work. The "many regulators run both" closing is realistic.
- **The validate-before-run script.** `verifier-validate.sh` runs cosign and GPG checks before the verifier itself, with explicit exit semantics. An examiner who is told "always run the validator first" can reliably do that without understanding the cryptography behind cosign.

### Narrow follow-ups

- **NOTE — Sample command lines and `--master-key` consistency.** The deployment-package's sample command lines (Standard verification / Strict mode / Bundle output) do NOT show `--master-key`. The quickstart and sample-report DO show `--master-key` in the standard examiner invocation. The deployment doc is for IT-shop allowlist purposes (where the command-line shape is illustrative, not normative for examiner work) but the inconsistency is the kind of thing a bench examiner notices. Suggest one of: (a) add `--master-key ./tenant-ikm.bin` to the sample command lines for consistency with the examiner-facing docs, or (b) add a one-line note above the samples ("Examiner-side invocations include `--master-key`; see `examiner-quickstart.md` for the standard examiner invocation. The samples below are abbreviated for IT-shop allowlist illustration."). Either resolves the apparent contradiction.

---

## Portfolio-comparison-procedures.md

For the EIC and supervisory team comparing across banks. This is not bench examiner reading; I read it from the angle of "is this the right doc for the EIC who needs to write a portfolio analysis?"

### What works

- **The realistic acknowledgement that comparison is manual until v1.1 ships.** "Until the v1.1 `verifier portfolio` subcommand ships, comparison is manual; this doc articulates the manual procedure." That's the right framing — the EIC does the work today with a spreadsheet; the doc tells them how. Promising tooling that doesn't exist yet would be unhelpful.
- **The metrics table.** Pass rate, anomaly rate, sealing-delay frequency, late-binding rate, master-rotation events, verifier-validation cadence. Six metrics is the right size for a portfolio comparison spreadsheet.
- **The sample 6-bank workspace.** Concrete example with realistic numbers (Bank C with HSM cluster issue Q1 remediated; Bank E with anomaly cluster under investigation). An EIC can copy this shape into their own portfolio review without inventing the columns.
- **The "common patterns and their implications" section.** Three patterns: portfolio-wide degradation, single-institution outlier, sub-portfolio commonality. Each has a "this often indicates" list and an "action" line. That's the right shape for an EIC's pattern-recognition reference.

### Narrow follow-ups

- **NOTE — Cross-period comparison and the new failure modes.** The cross-period comparison section calls out trend metrics (pass rate, anomaly rate, sealing-delay frequency, operational-events volume, configuration changes) but doesn't call out trends in specific spec §7 failure step categories. An EIC doing longitudinal analysis on an institution that has had two step 8 (`key_fingerprint mismatch`) findings in two consecutive examinations should be looking at that as a control-degradation pattern, not just an anomaly-rate pattern. Optional addition: one bullet under cross-period metrics ("Trends in spec §7 failure step categories — recurring step 7/8 findings indicate IKM-management control degradation; recurring step 6/9 findings indicate chain-integrity control degradation; recurring step 11 findings indicate signing-key-management control degradation"). Keys the EIC's longitudinal lens to the same step-numbered taxonomy the rest of the regulator pack uses.

---

## Spec §7 step 11 dual-algorithm cases (a)-(e)

I read the spec section because the new case (e) and the working-paper convention paragraph are called out as material additions. The spec is the implementer's doc, not the examiner's, but I want the bench examiner to be able to confirm the spec text matches the finding-language and response-workflow docs.

### What works

- **Case (e) is the right addition.** Both signatures present, one valid + one invalid, is the case the previous five-case enumeration was missing. The spec text correctly identifies the two interpretations (algorithm broken with un-broken algorithm still providing integrity assurance, or per-algorithm signing key compromise) and correctly leaves the interpretation to the institution's IR program. The verifier reports both algorithms' validation results; the institution interprets. That's the right separation of duties.
- **The working-paper convention paragraph.** "The verifier output MUST record both algorithms' validation results (PASS / FAIL / NOT-PRESENT per algorithm) when the seal record carries `signatures`. The examiner's working paper carries both rows; the examination report cites both algorithm validations." This is the load-bearing paragraph for examiner work — it tells the examiner that during the dual-algorithm transitional period, the working paper carries TWO algorithm validation rows for each affected seal-day, not one. Without this convention, an examiner could record the primary algorithm's PASS and overlook the co-signed algorithm's FAIL, missing a Scenario 4 finding.
- **The "this convention applies during the multi-year transitional period; once the institution retires one algorithm, the convention reverts to single-algorithm reporting" closing.** Tells the examiner the convention has a sunset. That matters because examiners build muscle memory; knowing the muscle memory will revert after the transition prevents over-engineering downstream.

### Narrow follow-ups

- **EXAMINER — Spec text severity for case (e).** The spec says case (e) is `--strict` FAIL and non-strict PASS-WITH-ANOMALY. The finding-language doc says case (e) is Severe MRA in either mode. These are not contradictory (the spec describes verifier output; the finding language describes examiner severity guidance) but the bench examiner reading both will need to know that the verifier's PASS-WITH-ANOMALY in non-strict mode does NOT translate to "Observation" severity in the examination report. The finding-language guidance overrides — case (e) is Severe regardless of the verifier's pass/fail bracket. Suggest one sentence in the spec §7 case (e) paragraph: "Examiner severity treatment for case (e) is in `regulator-pack/finding-language.md` and is Severe regardless of the verifier's pass/fail bracket." This stops a junior examiner from seeing PASS-WITH-ANOMALY in non-strict and writing it up as an Observation.

---

## Cross-doc consistency checks

I went through the regulator-pack docs looking for places where the same fact is stated in multiple docs and confirmed they agree.

| Fact | Doc 1 | Doc 2 | Doc 3 | Consistent? |
|---|---|---|---|---|
| `--master-key` is in the standard examiner invocation | quickstart §"What you do" | sample-report §"Verifier invocation" | training Module 3 (no `--master-key`) | **Inconsistent — see follow-up below** |
| `step: 8` is "stop and call the bank" / Severe MRA | quickstart 14-row table + 3 specific points | finding-language step 8 row + paragraph | training Module 5 row | **Consistent** |
| `step: 8` is identity mismatch, NOT content tampering | quickstart specific point | finding-language paragraph (bold) | response-workflow step 4 | **Consistent** |
| Audit file ends mid-line is operational, NOT tampering | quickstart specific point | finding-language paragraph | training Module 5 row | **Consistent** |
| Spec §7 step 11 case (e) co-signed seal failure | spec §7 step 11 (e) | finding-language dual-algo row 3 | response-workflow step 2 dual-algo row | **Consistent** |
| 14-row taxonomy count | quickstart 14 rows | finding-language 14 rows | training Module 5 11 rows | **Module 5 is a subset (intentional); see narrow follow-up above** |
| IR scenario column present | quickstart (yes) | finding-language (yes) | training Module 5 (yes) | **Consistent** |

The one cross-doc inconsistency is the `--master-key` flag: examiner-quickstart and sample-report show it in the standard invocation; training Module 3 does not. The training Module 3 invocation should be updated to match — a new examiner trained on the Module 3 invocation will produce a structural-only verification on their first real examination and have to be told to re-run with `--master-key`.

- **EXAMINER — Training Module 3 invocation should include `--master-key`.** Training Module 3 currently shows the invocation without `--master-key`. Per the quickstart and sample-report, the standard examiner invocation includes it (the institution provides the IKM under a documented disclosure shape per legal-disclosure.md and customer-dispute-procedures.md). Update Module 3 to include `--master-key ./tenant-ikm.bin` in the invocation, and add one line in the surrounding text mirroring the quickstart's explanation ("The `--master-key` flag points at the institution's IKM file; without it the verifier performs structural verification only and skips per-event HMAC equality. For FFIEC examiner work the IKM is provided per the institution's disclosure shape — see `legal-disclosure.md`."). This is the only finding I'd write up as EXAMINER-severity because it produces wrong examiner behavior on the first real examination.

---

## Items I'd want to see before signing off (none EXAMINER-severity except the one above)

| Severity | Item | Doc |
|---|---|---|
| EXAMINER | Training Module 3 invocation missing `--master-key`; will produce wrong first-examination behavior | examiner-training.md |
| MINOR | Step 9 finding paragraph "[is/is not] underway" template style inconsistent with rest of doc | finding-language.md |
| MINOR | FAIL example block ordering (cover page summary should precede per-day FAIL detail for reading order) | sample-report.md |
| NOTE | IKM file-mode reference doesn't translate to Windows | examiner-quickstart.md |
| NOTE | Add row for `--master-key absent` outcome in 14-row table | examiner-quickstart.md |
| NOTE | JSON snippet pairing across sample-report and response-workflow | sample-report.md |
| NOTE | Step 6 paragraph "most likely tampering signal" framing is too strong without facts | finding-language.md |
| NOTE | Repeat-finding examples for new step 2/3/5/6/9 paragraphs | finding-language.md |
| NOTE | Step-mapping table cross-reference for cases (a) and (d) of spec §7 step 11 | examination-response-workflow.md |
| NOTE | Case (e) two-track IR-scenario clarification | examination-response-workflow.md |
| NOTE | Module 5 row for partial-coverage and algorithm-not-on-posture-list dual-algorithm cases | examiner-training.md |
| NOTE | Module 5 row count vs full taxonomy clarifying note | examiner-training.md |
| NOTE | II.E "examiners confirm" list could add verifier-dispatch testing | handbook-mapping.md |
| NOTE | Sample command lines `--master-key` consistency | deployment-package.md |
| NOTE | Cross-period comparison metrics keyed to spec §7 step categories | portfolio-comparison-procedures.md |
| EXAMINER | Spec case (e) paragraph cross-reference to finding-language severity overriding verifier pass/fail bracket | spec/chain-of-custody-v1.md |

Two EXAMINER-severity items, both narrow and both fixable with single-line edits. One is in training (the `--master-key` invocation); the other is in the spec (a one-sentence cross-reference to finding-language so the examiner doesn't misread the verifier's PASS-WITH-ANOMALY as an Observation).

Everything else is MINOR or NOTE — places where a clarity edit or an optional addition would tighten the loop, but nothing that would prevent a competent examiner from producing a defensible report today.

---

## Stopping criterion

0 EXAMINER-severity findings preventing examiner workflow + 0 EXAMINER-severity findings preventing report production. Two EXAMINER-severity findings remain, both narrow, both single-line edits to source text (one in training, one in the spec cross-reference). Neither prevents an examiner from producing a defensible report — they prevent a NEW examiner from making a first-examination mistake. With those two edits, the examiner-pack reaches 0/0.

The corpus is at the level I'd accept on a national IT examination program. The five docs that matter to the bench examiner are coherent, cross-referenced, and key off the same `step` field consistently. The dual-algorithm transitional-period material is the right level of specificity — enough that a bench examiner can act on a co-signed seal failure without becoming a cryptographer, not so deep that the doc is unreadable.

The 2026-04-22 PASS / 2026-04-23 FAIL pairing in the sample report is the single pedagogical decision I'd most want other regulator packs to copy. The `step` field as the spine of the whole loop is the single architectural decision I'd most want other examination tools to copy.

— Janet Buckley
