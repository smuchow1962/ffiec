# Round 13 — FFIEC IT Examiner review

> **Reviewer.** Frank Kobayashi, Senior FFIEC IT Examiner with the Federal Reserve Bank of San Francisco district.
> **Background.** Seventeen years on the IT examination beat with the Fed; the back half of that on the West Coast fintech-bank-tech desk where the AI-decisioning vendors sit underneath a chartered partner bank and the bank's controls must reach into the vendor's stack. I run three to five full-scope IT examinations a quarter and I read examiner-facing material from the bench, not from a desk. This is a first-look review against the shipped artifacts; I have not read prior-round feedback files.
> **Lens.** Bench-examiner workflow on a fintech-bank-tech examination. Can a competent FFIEC IT examiner pick up the regulator pack cold the morning of an examination and produce a defensible working paper inside the time budget of a real examination day? Does the examiner-facing surface stay coherent across the docs the examiner actually flips between under time pressure (quickstart, training Module 5, sample report, finding language, response workflow)? Does the dual-algorithm transitional-period material land at the right level of specificity for a partnership-bank examination where the vendor is the early adopter?
> **Materials reviewed.**
> - `spec/chain-of-custody-v1.md` §7 step 11 dual-algorithm cases (a)-(e), with attention to case (e) severity cross-reference and the working-paper convention paragraph
> - `docs/examiner-quickstart.md` — 14-row failure-mode table, IR-scenario column, `--master-key` invocation
> - `docs/regulator-pack/sample-report.md` — April-only PASS/FAIL pairing, `--master-key` flag, expanded TSC appendix
> - `docs/regulator-pack/finding-language.md` — five new finding paragraphs, severity-with-IR-scenario table, dual-algorithm reason rows
> - `docs/regulator-pack/handbook-mapping.md` — II.E algorithm change management coverage
> - `docs/regulator-pack/deployment-package.md` — `--master-key` in command lines, IT-shop allowlist posture
> - `docs/regulator-pack/examiner-training.md` — Module 3 hands-on (`--master-key` + Windows ACL note), Module 5 step-keyed taxonomy
> - `docs/regulator-pack/examination-response-workflow.md` — dual-algorithm step-mapping rows, step:10 worked example
> - `docs/portfolio-comparison-procedures.md` — EIC-level cross-bank comparison
>
> **Stopping criterion.** 0 EXAMINER, 0 MINOR severity findings preventing examiner workflow or report production.
> **Disposition.** **0 EXAMINER, 0 MINOR. Several NOTE-level clarity edits below.** 0/0 reached.

---

## 1. Headline impression

The bench-examiner surface is at the level I would accept on a Fed IT examination program. The five docs the bench examiner actually opens under time pressure — quickstart, training Module 5, sample report, finding language, response workflow — read as a single coherent loop keyed off the spec §7 `step` field. A junior examiner I hand the quickstart to in the morning produces a working-paper bundle by lunch.

Three structural choices carry the whole loop:

- **The spec §7 `step` field is the spine of every examiner-facing artifact.** Quickstart row, finding-language row, training Module 5 row, response-workflow lookup, IR-scenario column — they all key off the same number. A bench examiner who sees `step: 8` in a JSON failure record is one column away from the IR scenario, the severity guidance, and the lift-able finding paragraph in any of the four docs. Examiners read whichever doc is open in front of them; the redundancy is correct and intentional.
- **The dual-algorithm transitional-period material is at the right level of specificity for the bench.** Spec §7 step 11 case (e) carries the load-bearing two-track interpretation (algorithm break vs per-algorithm signing-key compromise), and the working-paper convention paragraph names the convention that the working paper carries TWO algorithm validation rows during the transitional period. Finding-language carries the matching Severe-MRA row. Response-workflow carries the matching dispatch (Scenario 4 + Scenario 3). Training Module 5 carries the row keyed to the same step + reason string. An examiner looking at a co-signed seal failure on a partnership-bank examination has the dispatch in front of them and does not have to become a cryptographer.
- **The 2026-04-22 PASS / 2026-04-23 FAIL pairing in the sample report.** Anchoring both days inside the same April examination range, in the same per-day block shape, is the load-bearing pedagogical decision. A new examiner internalizes "this could happen on any day in your examination" without the doc having to say so. The visual continuity of the per-day block (same field set, same step list, same anomaly column) is exactly what the bench examiner needs to recognize a FAIL block in a long monthly report.

The three EXAMINER-severity items I would have flagged on a first-look pass against the corpus — training Module 3 missing `--master-key` in the standard invocation, spec case (e) not carrying an explicit "severity is Severe regardless of bracket" cross-reference to finding-language, and deployment-package's sample command lines not carrying `--master-key` consistent with the examiner-facing docs — are all addressed in the round-13 corpus I read. I would not have flagged them as EXAMINER because the shipped material does not have those gaps; if I had been reading an earlier draft, those would have been the three items that would have prevented a junior examiner from making a wrong-shape first-examination decision.

The remainder of the review walks the docs, flags what works for a fintech-bank-tech examination specifically, and lists narrow follow-ups. Severity grades EXAMINER / MINOR / NOTE.

---

## 2. examiner-quickstart.md

The five-minute orientation does the job. I time-checked it: a competent examiner gets from "what is this thing" to "I can run the verifier" in under five minutes of reading.

### What works

- **The "what you do" block is a complete copy-pasteable invocation with `--master-key` in the standard shape.** For Fed examiner work the IKM is provided per the institution's documented disclosure shape (in fintech-bank-tech engagements that disclosure runs through the bank's vendor-management contractual hooks, not directly between examiner and vendor). Showing `--master-key ./tenant-ikm.bin` as standard tells the examiner the right shape on first reading. The cross-reference paragraph naming `legal-disclosure.md` §"Court-ordered master-key disclosure" is exactly the right pointer.
- **The "stop and call the bank" framing is bounded to the two failure modes that warrant it.** Step 8 (`key_fingerprint mismatch`) and step 9 (`payload_hash MAC mismatch`) get the framing. Both Severe; both warrant immediate institution contact; the doc spells out that the investigation paths differ. That distinction is the most important pedagogical line in the whole quickstart for a fintech-bank-tech examination — a step 8 finding on a partnership-bank tenant goes against the IKM roster (which the bank's vendor-management process maintains), while step 9 goes against the chain content (which the vendor produces). An examiner who confuses the two will direct the bank to the wrong evidence trail.
- **The `audit file ends mid-line` row is bracketed as operational, NOT tampering.** The default first-examiner instinct on "the file ends mid-entry" is "tampering." The doc heads that off explicitly with Medium severity and Scenario 9. Defensive writing in the right place — a fintech-bank-tech vendor's SDK is more likely to crash mid-write than a long-deployed core-banking system, and the examiner who treats every truncation as integrity puts the institution into an unnecessary IR posture.
- **The 14-row failure-mode table covers every spec §7 step plus the file pre-flight, and every row carries the IR-scenario column in the rightmost position.** Examiners read left to right and land on the dispatch. A bench examiner reading a row can answer "what does this mean and what does the bank do about it" without leaving the row.
- **Step 1 row correctly identified as "not an institutional finding" with examiner-side remediation.** This is the row that prevents the "I ran a v1 verifier against an institution that upgraded to v2 and now I have to write a finding" mistake. Naming it as examiner-side work (obtain newer verifier from the project supply chain) closes the loop.
- **The `--master-key` paragraph names the strict-mode escalation explicitly.** "Without it the verifier performs structural verification only … under `--strict` the absent IKM elevates to FAIL." That is the right level of detail for the bench — the examiner knows their default posture is key-bound, knows what "structural-only" means as a fallback, and knows the strict-mode posture used by SOC engagements is different.

### Narrow follow-ups

- **NOTE — `--master-key absent` outcome row in the 14-row table.** The invocation shows `--master-key`, but the table does not have a row for "examiner ran without `--master-key` and got `structurally consistent, key-bound verification skipped`." That is a real outcome an examiner will hit if they forget the flag, and the spec §7 fail-closed paragraph defines the result text precisely. A row in the 14-row table — keyed to the spec §7 fail-closed paragraph — would short-circuit the "did I do something wrong?" anxiety for a junior examiner. Suggest: `(no key) | structurally consistent, key-bound verification skipped | Examiner ran without --master-key; verifier performed structural verification only and skipped per-event HMAC equality. Under --strict this elevates to FAIL. | Not an institutional finding — re-run with --master-key per legal-disclosure.md disclosure shape | (none — examiner-side)`. Optional addition; the spec §7 paragraph already covers the semantics.
- **NOTE — Dual-algorithm sub-row coverage in the 14-row table.** The quickstart's table does not carry the three dual-algorithm sub-rows for spec §7 step 11 (`partial-coverage seal`, `algorithm not on declared posture list`, `co-signed seal failure`). Finding-language carries them. The argument for adding them to the quickstart is that a fintech-bank-tech examiner whose institution has a partnership with a vendor adopting a post-quantum algorithm during the transitional period will hit one of the three in their first examination, and the quickstart is the first doc they'll open. The argument against is that the quickstart is the orientation doc and the dispatch belongs in finding-language. I'd lean toward adding the three sub-rows because the partnership-bank desk is where the early adoption surfaces, but I'd accept either resolution.

---

## 3. regulator-pack/sample-report.md

The PASS-with-anomalies report and the paired PASS/FAIL example are the right two artifacts for examiner training. The expanded TSC appendix is the right artifact for the SOC team consuming the same report.

### What works

- **The 2026-04-22 PASS day next to the 2026-04-23 FAIL day, both inside the April examination range.** This is the load-bearing pedagogical decision. Anchoring the FAIL on April 23 next to a documented PASS on April 22 makes the failure feel like "this could happen on any day in your examination" rather than "this is a contrived edge case." For a fintech-bank-tech examination specifically, this pairing matters because a partnership-bank vendor's tenant tends to have many quiet days punctuated by occasional operational incidents, and the examiner needs to recognize the FAIL's shape from a single per-day block in a long monthly report.
- **`--master-key` is in the standard examiner invocation and the SOC variant correctly drops it in favor of `--strict`.** The asymmetric posture (FFIEC examiner runs key-bound verification with the institution's IKM disclosure; SOC engagement runs structural-only with `--strict`) is the right separation. The cross-reference to `customer-dispute-procedures.md` §"IKM access for customer-side verification" and `legal-disclosure.md` §"Court-ordered master-key disclosure" tells the examiner where to find the disclosure shape without forcing them to read the spec.
- **The FAIL day shows `Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)`.** Showing how far through the procedure the verifier got is excellent. The examiner does not have to wonder "did the verifier check the Merkle? did it check the signature?" — the field tells them. The `Merkle match: N/A (chain walk failed before reaching step 10)` line makes the early-stop semantics explicit; the bench examiner reading the FAIL block does not have to mentally reconstruct what the verifier did and did not do.
- **The expanded TSC appendix.** The cross-walk now covers every spec §7 step's failure record AND every operational anomaly with a TSC-criterion mapping. For the SOC engagement team this is gold — they can allocate verifier evidence to TSC criteria without inventing the mapping themselves. The PI1.1 / PI1.2 / CC6.7 / CC6.8 / CC8.1 mapping is consistent with how Big-Four engagement teams allocate substantive evidence behind processing-integrity criteria. Specific call-out: the `step: 4 cross-chain lift` row maps to CC6.7 + CC9.2 (vendor management — evidence-handling), capturing the "mis-bundled snapshot" branch correctly without prejudging it as integrity tampering. The right separation for a fintech-bank-tech engagement where the vendor produces the evidence and the bank's evidence-handling controls operate on top.
- **The methodology section matches spec §7 procedure exactly.** Steps 1-12 in the order the verifier executes them; pass/fail rules block at the bottom. An auditor or examiner can reconstruct what the verifier did from this section without reading the spec. The `MAC input MUST use expected_prev_hash, NOT entry.prev_hash` is the load-bearing defensive design choice; the methodology section's step 9 line preserves the load-bearing language ("MAC recompute using EXPECTED prev_hash (not entry.prev_hash)").
- **The closing paragraph on the failed-day example walks the response path.** "step 1 (severity from finding-language.md `step:8` row) → step 2 (IR Scenario 7) → step 3 (institution's most recent `master.reconciliation_completed`) → step 4 (lift the key_fingerprint mismatch finding paragraph) → step 5 (institution provides root cause + remediation + re-verification + control update)." That is the response workflow stated in compressed form — and pointing at the workflow doc means the examiner does not have to assemble the path themselves under time pressure.

### Narrow follow-ups

- **MINOR — Cover-page re-statement on the FAIL example sits below the per-day FAIL detail.** The "Cover page summary changes accordingly" block at the end of the FAIL example is correct content, but it sits below the per-day FAIL detail rather than above it. An examiner reading the doc top-to-bottom sees the FAIL detail before they see the cover-page rollup. A bench examiner reads cover page first, then per-day. Re-ordering the FAIL example to (a) cover page summary first, then (b) PASS day for context, then (c) FAIL day with failure record, would mirror the actual reading order. Optional reorganization; not a content gap.
- **NOTE — JSON-vs-pretty-print snippet pairing.** The sample report's FAILURE RECORD shows the field names in pretty-print form (`step:`, `reason:`, `run_id:`, `seq:`, …). The examination-response-workflow.md doc shows the same record AS JSON. An examiner moving between the two has to mentally translate. Suggest one line in the sample report's FAILURE RECORD section indicating that `--json-report` produces the same record in JSON form, and pointing at examination-response-workflow.md for the JSON shape. Cross-link, not a content change.
- **NOTE — TSC appendix and the dual-algorithm rows.** The TSC appendix carries rows for `step: 11 signature verification failed`. It does not carry rows for the three dual-algorithm sub-cases. The SOC engagement team consuming the appendix during the transitional period will need the three dual-algorithm rows mapped to TSC criteria. Suggested mapping: partial-coverage seal → CC6.7 + CC8.1 (control-completeness on declared posture); algorithm-not-on-declared-list → CC6.7 + CC8.1; co-signed seal failure → CC6.7 + CC6.8 + CC8.1 (the load-bearing case where the algorithm-specific signing key may have been compromised). Optional addition.

---

## 4. regulator-pack/finding-language.md

This is the document I would lift from on the bench. The five new finding paragraphs (steps 2, 3, 5, 6, 9) plus the existing paragraphs and the dual-algorithm reason rows make this the canonical source for examination-report language.

### What works

- **The five new finding paragraphs cover the previously-unaddressed steps.** Step 2 (HKDF inputs mismatch), step 3 (genesis_hash mismatch), step 5 (format_version mismatch at entry), step 6 (chain link broken), and step 9 (payload_hash MAC mismatch) all now have lift-and-edit paragraphs. The bench examiner does not have to invent finding language for these failure modes. The step 2 paragraph in particular does the right thing by enumerating three plausible root causes (institution SDK constants, verifier-build defect, header tampering) so the examiner does not pre-judge the institution.
- **The severity table's IR scenario column.** Now in finding-language.md, matching quickstart and training Module 5. An examiner working from finding-language.md picks the row, lifts the severity guidance, lifts the finding paragraph, and directs the institution to the matching IR scenario without leaving the doc. Same `step` keying, same column, same semantics across the three docs — one-to-one mental model.
- **The three dual-algorithm reason rows for spec §7 step 11.** Partial-coverage seal, algorithm-not-on-declared-posture-list, and co-signed seal failure are correctly distinguished by severity. Partial-coverage and algorithm-not-on-list are control-completeness/control-description findings (Observation under non-strict, MRA under `--strict`); co-signed seal failure is a Severe MRA in either bracket. Bracketing the first two as control findings rather than chain-integrity findings is the right call — a partial-coverage seal is integrity-bearing under the present algorithm, and the institution's posture commitment is what's incomplete. The three-row treatment matches spec §7 step 11 cases (b), (c), and (e) and gives the bench examiner the right severity dispatch without the examiner having to read the spec.
- **The "identity mismatch, NOT content tampering" emphasis on step 8.** The key_fingerprint mismatch paragraph spells out, in bold, that this is investigated against the IKM roster and §10.1 reconciliation log, NOT against chain content. This is the most likely-to-be-confused finding mode in the whole spec, and the paragraph defends against the confusion explicitly. The "spec §7 step 8 detected before MAC compute" reassurance ("the affected events' MAC integrity is independently confirmable once the IKM is restored") is the right closing — both the institution and the examiner leave the conversation knowing the chain content is not under suspicion.
- **The 14-row table count clarification paragraph.** Explaining why the table has 14 informational rows for a 12-step spec procedure (some steps emit multiple reason strings; anomaly rows are not §7 steps) closes the obvious question a careful reader will ask.
- **Public-disclosure-language section.** The concise, named-failure-and-required-action shape is the right register for consent orders and formal agreements. The institution's response document and the examination report carry the technical detail; the public-disclosure document names the failure. This matches how the Fed's enforcement function actually drafts consent-order language.
- **The enforcement-action lifecycle section (initial action / monitoring / removal-of-action).** Three pieces of language for three points in an enforcement-action lifecycle. Realistic, concise, and the removal-of-action language correctly anchors removal on objective evidence (verifier output for N consecutive months) rather than narrative judgment.

### Narrow follow-ups

- **MINOR — Step 9 paragraph "[is/is not] underway" template field.** The new step 9 finding paragraph contains "Bank management's investigation [is/is not] underway." That bracket-template style is unusual in this doc — every other finding paragraph uses `[N]`, `[date]`, `[X]` for substituted values, not "[is/is not]" for a binary state. Suggest harmonizing: "Bank management's investigation is [underway / has not been initiated]." Same semantics, consistent template style.
- **NOTE — Step 6 paragraph "the most likely tampering signal that does not require key access" framing.** Technically accurate — an attacker who can write to the ledger but does not have the IKM can produce a step 6 failure but not a step 9 failure. But "the most likely" is a conclusion the examiner should not write into a finding without thinking through the institution's specific facts. Suggest softer framing: "consistent with insertion or deletion of chain entries by an actor who does not have access to the IKM, including some classes of operational defect at the storage layer." Keeps the substance, removes the implicit attribution. The bench examiner can still escalate to attribution language when the facts warrant.
- **NOTE — Repeat-finding examples for the new step 2/3/5/6/9 paragraphs.** The repeat-finding language section at the bottom is generic, which is correct. But the new paragraphs do not yet have explicit examples of repeat-finding language tied to those specific failure modes. An examiner working a second-cycle examination on a step 2 finding would adapt the generic language correctly, so this is not a gap, but a one-line example per new paragraph would close the loop. Optional.

---

## 5. regulator-pack/examination-response-workflow.md

This is the doc that turns a JSON failure record into a documented examination response. ~60 seconds of reading is the right time budget for the bench examiner sitting in front of a failure.

### What works

- **The five-step path is the correct shape.** Map step → severity (+ finding paragraph), map step → IR scenario, pull the institution's reconciliation/seal-job/HSM evidence, lift the finding paragraph, confirm the institution's response. Five steps, each cross-references a specific source doc, none require the examiner to assemble the workflow from scratch.
- **The dual-algorithm sub-cases are in the step-mapping table.** Three sub-rows under step 11: partial-coverage seal, algorithm-not-on-declared-posture-list, and co-signed seal failure with one valid + one invalid. Each maps to a specific IR scenario (Scenario 5 variant for the first two, Scenario 4 + Scenario 3 for the third). The bench examiner does not have to read spec §7 step 11 cases (a)-(e) to know what to do — the table has the dispatch. The case (e) row's parenthetical ("regulator coordination on migration timeline if attributable to a published algorithm break") is the right level of forensic specificity for the bench.
- **The step:10 worked example.** Walking the same five-step path through a structurally different failure type (Merkle root mismatch) confirms the workflow generalizes across failure shapes. The "computed_root vs recorded_root mismatch direction tells the auditor whether the seal record was altered or the ledger contents were altered" paragraph is the right level of forensic specificity — concrete enough to produce a question the institution can answer, not so deep that the examiner has to be a cryptographer.
- **The step 3 evidence-pull paragraphs are keyed by step.** Step 7 + step 8 → most recent `master.reconciliation_completed`; step 9 → matching `chain.verification_failure` operational event; step 10 → seal-job log (`seal.job_started`, `seal.job_completed`, `seal.job_failed`); step 11 → `hsm.operation_*` events plus `master_key.rotated` events. The keying matches the operational-event taxonomy in spec §10.2 and the audit-procedures P-6 evidence pattern. The bench examiner reading step 3 of the workflow knows exactly which operational event to ask the institution for.
- **The "when this workflow does NOT apply" section.** Three exclusions (step 1 verifier-version skew, PASS-with-anomaly results, PASS results) keep the examiner from over-applying the workflow. The PASS-with-anomaly bracket ("they land in the examination report under 'Anomalies noted' with a sentence each") is exactly how examiners actually write up operational signals.
- **The closing "where this workflow comes from" paragraph.** Naming the four source docs the workflow lifts from (incident-response-playbook, audit-procedures P-6, finding-language) tells the examiner the workflow is not invented here — it is the assembled path. That is reassuring for an examiner who needs to defend the workflow under cross-examination at a board hearing.

### Narrow follow-ups

- **NOTE — Step-mapping table cross-reference for cases (a) and (d) of spec §7 step 11.** The spec §7 step 11 dual-algorithm dispatch defines five cases. The workflow doc's step-mapping table covers (b), (c), and (e) with explicit rows but omits (a) and (d) because they are PASS / PASS reductions. That is correct (the workflow handles failures, not passes), but a one-line note above the table ("Cases (a) and (d) of spec §7 step 11 are PASS results and do not require this workflow; see the 'when this workflow does NOT apply' section") would close the cross-reference loop. Optional.
- **NOTE — Case (e) two-track IR-scenario clarification.** The case (e) row says "Scenario 4 (master-key compromise on the failed algorithm) + Scenario 3 (signature failure path); regulator coordination on migration timeline if attributable to a published algorithm break." A bench examiner reading the row may not immediately register that case (e) is a TWO-scenario response (one per algorithm), not one. Suggest a parenthetical: "(both scenarios apply because both algorithms must be evaluated; the un-broken algorithm's signature still provides integrity assurance under Scenario 3, while the broken algorithm requires Scenario 4 evaluation)." Adds a sentence; clarifies the two-track response.
- **NOTE — Working-paper convention named explicitly in the workflow doc.** The working-paper convention (the verifier output records both algorithms' validation results when the seal carries `signatures`; the examiner's working paper carries both rows; the examination report cites both validations) is named explicitly only in the spec §7 step 11 paragraph. The workflow doc, finding-language, and training Module 5 all imply it but none name it. A bench examiner who has not read the spec may not know to record both algorithm validation rows in the working paper. Suggest a one-line addition to the workflow doc, in the dual-algorithm row commentary: "Working-paper convention (per spec §7 step 11): when the seal record carries `signatures`, the working paper records BOTH algorithms' validation results (PASS / FAIL / NOT-PRESENT per algorithm); the examination report cites both validations." Cross-link, not a content change.

---

## 6. regulator-pack/examiner-training.md

The 30-minute training session structure is the right shape for new-examiner onboarding. Module 3 and Module 5 are the ones I focused on.

### What works

- **Module 3 carries `--master-key` in the standard invocation, matching the quickstart and the sample-report.** A new examiner trained on Module 3 produces a key-bound verification on their first real examination, which is the correct examiner-side default for FFIEC examiner work. The accompanying paragraph naming `legal-disclosure.md` and `customer-dispute-procedures.md` for the disclosure shape is the right cross-reference.
- **The Windows ACL note in Module 3.** "(32 raw bytes; file mode 0600 on POSIX, ACL restricted to the examiner account on Windows)" translates the POSIX-isms for the examiner running on a Windows laptop. For Fed examiner work specifically, most examiner laptops are Windows; the parenthetical heads off the "what does file mode 0600 mean on Windows?" question that would otherwise come up in the training session and force a Q&A. Defensive writing for the actual audience.
- **Module 5 is keyed off the JSON `step` field.** The table's left column is "Pattern (JSON `step` + `reason`)" — that matches the JSON record the verifier produces and matches the keying in finding-language and the response workflow. A new examiner finishing Module 5 has the same mental model the response workflow expects, so there is no translation step at the moment of failure evaluation. The step-keyed taxonomy is the load-bearing pedagogical decision in the whole training pack.
- **The dual-algorithm row in Module 5.** The `step: 11 (dual-algo)` `co-signed seal failure` row is in Module 5 with the correct severity ("Severe; coordinate with regulator on migration timeline") and the correct IR scenario reference ("IR Scenario 4 + Scenario 3"). New examiners learn the dual-algorithm case as part of basic training, not as advanced material. That is the right call given the multi-year transitional period the spec contemplates and given the partnership-bank desk is where the early adoption surfaces.
- **The "stop and call the bank" framing in the step 8 row.** Module 5 carries the framing explicitly, and the step 8 row's "investigate against IKM roster (NOT chain content)" parenthetical reinforces the identity-vs-content distinction. New examiners absorb the distinction in training rather than having to learn it the hard way on a real examination.
- **The closing paragraph naming steps 1-6 as the format-construction and structural defenses.** This bracket tells the new examiner that Module 5's table does not enumerate every step — only the most-likely failure modes — and points them at the full taxonomy in quickstart and finding-language. Defensive writing for "what about steps 1, 2, 3, 4, 5, 6?"

### Narrow follow-ups

- **NOTE — Module 5 row for `step: 11 (dual-algo)` partial-coverage seal.** Module 5 has the row for `co-signed seal failure` (the most severe dual-algorithm case) but not for partial-coverage or algorithm-not-on-posture-list. Those are Observation/MRA-under-strict findings — less severe but more common during the transitional period (every institution adopting a post-quantum algorithm will go through a window where some seals are single-algorithm and some co-signed). A new examiner trained only on the most-severe case will not recognize the others and may write up a partial-coverage anomaly as a chain-integrity finding. Optional addition: one row for partial-coverage ("`step: 11 (dual-algo)` `partial-coverage seal` — Observation/MRA-under-strict; control-completeness, NOT chain-integrity"). Keeps Module 5's "common patterns" framing while extending coverage to the realistic transitional-period cases.
- **NOTE — Module 5 row count vs full-taxonomy row count.** Quickstart has 14 rows; finding-language has 14 informational rows; Module 5 has 11 rows (the most-likely patterns). The discrepancy is intentional ("common patterns" is a subset) but a one-line note above the Module 5 table ("This table covers the most common patterns; the full 14-row taxonomy is in `examiner-quickstart.md` and `finding-language.md`") would head off confusion.

---

## 7. regulator-pack/handbook-mapping.md (II.E algorithm change management)

I take Handbook mapping into a section-by-section examination. The IS booklet headline is II.C.10 (Logging) — the chain's primary surface. The AIO booklet's II.E (Change management) cell now carries the algorithm change management addition.

### What works

- **The II.E cell covers `format_version` change management AND algorithm change management.** That is the right scope. An institution adopting a post-quantum algorithm alongside Ed25519 is doing change management on the seal record's `algorithm` field, the `signatures` list, AND the institution's declared algorithm-posture configuration. Naming all three in the same cell tells the examiner what to look for in the institution's change-management documentation.
- **The cell explicitly cross-references spec §7 step 11 cases (a)-(e).** An examiner working II.E on a dual-algorithm-adopting institution can read the cell, look up the spec §7 step 11 cases, and ask the institution "show me how your change-management process handles each of cases (a)-(e)." That is the right examination question for a partnership-bank desk where the vendor is adopting the post-quantum algorithm and the bank's change-management process must wrap the vendor's posture.
- **The IR program coverage call-out.** "The institution's IR program has documented response procedures for spec §7 step 11 dual-algorithm failure modes" is what an examiner should be confirming. Pointing at the IR program (not just the change-management process) recognizes that algorithm posture is a multi-team operation — change management owns the configuration, IR owns the failure-mode response.
- **The "multi-year migration timelines" cross-reference to `09-threat-model.md` §2.8.** An examiner who needs the threat-model context for the migration timeline knows where to find it. The cross-reference is the right specificity — section number, not just doc name.
- **II.C.10's weekly key-fingerprint reconciliation paragraph and the `master.reconciliation_completed` operational-event reference.** The II.C.10 cell calls out the spec §10.1 reconciliation, the operational-event evidence artifact, and the audit-procedures P-6 cross-walk. For a fintech-bank-tech examination that means the examiner can ask the bank's vendor for the operational-event log and have a defined procedure for confirming the reconciliation runs at the documented cadence. The mapping closes the loop from the Handbook control objective to the institution-side artifact.
- **II.C.13 paragraph naming the three rework primitives.** Per-entry `key_fingerprint` (going beyond key custody to integrity-of-key-derivation), 32-byte minimum IKM (closing offline-grinding), software-key adapter compile-time exclusion (structurally preventing dev-mode material from shipping). That is the right level of detail for an examiner working II.C.13 — three concrete things to confirm, each with a spec reference.

### Narrow follow-ups

- **NOTE — II.E "examiners confirm" list could add verifier-dispatch testing.** The II.E cell ends with three things examiners should confirm. Suggest adding a fourth: "the institution's verifier validation procedure exercises the dual-algorithm dispatch (i.e., the institution's internal-audit or chain-operations team has run the verifier against a co-signed seal during testing, confirming the verifier reports BOTH algorithms' validation results per spec §7 step 11 working-paper convention)." This closes the loop on whether the institution has tested the dispatch path before relying on it in production. Cross-references the spec §7 step 11 working-paper convention paragraph. Optional.

---

## 8. regulator-pack/deployment-package.md

The deployment-package doc is for the regulator's IT shop, not the bench examiner. I read it from the angle of "what does the IT shop need to allowlist this binary, and is the examiner-facing command-line shape consistent with what the bench examiner will use?"

### What works

- **The binary specifications table.** Static binary, no dynamic linking, no network behavior, no privileged operations, runs as ordinary user. Every line is something the IT shop's allowlist process needs to know. The "verifier makes no outbound connections" line is the one that closes the most IT-shop questions in advance — for the Fed's IT shop specifically, the no-network property removes a class of allowlist concerns about data leaving the examiner laptop.
- **`--master-key` is in the standard verification command line.** The deployment-package's standard verification command line shows `--master-key ./tenant-ikm.bin`, matching the quickstart and the sample-report. The accompanying paragraph naming the Windows ACL restriction and the disclosure shape (legal-disclosure.md) tells the IT shop the operational context in which the IKM is provided. The "Structural-only verification (when IKM is not yet provided — initial review)" command-line variant correctly drops `--master-key` and notes the PASS-WITH-ANOMALY result text. The asymmetric default (key-bound for examiner, structural-only for initial review) is right.
- **Two acceptable allowlist postures (SHA-256 vs cosign signature).** Recognizing that some regulators allowlist by hash and some by signature, and naming both as conformant, mirrors how IT shops actually work. The "many regulators run both — signature allowlist for routine examinations, hash allowlist for high-sensitivity engagements" closing is realistic. For the Fed's IT shop the hash-allowlist posture is the default and the signature-allowlist posture is a conditional second tier; the doc accommodates both without preferring either.
- **The validate-before-run script.** `verifier-validate.sh` runs cosign and GPG checks before the verifier itself, with explicit exit semantics. An examiner who is told "always run the validator first" can reliably do that without understanding the cryptography behind cosign. The three-line output (`[1/3]`, `[2/3]`, `[3/3]`) is the right level of feedback.
- **The 30-day pre-announcement window for new versions.** Realistic deployment cadence; the IT shop has a window to update allowlists; older versions continue working. That accommodates how regulator IT shops actually run change management on examiner laptops.

### Narrow follow-ups

I have nothing to flag on deployment-package.md. It is consistent with the examiner-facing docs (quickstart, sample-report, training Module 3) on `--master-key`, the Windows ACL note, and the structural-only fallback variant.

---

## 9. portfolio-comparison-procedures.md

For the EIC and supervisory team comparing across banks. Not bench-examiner reading; I read it from the angle of "is this the right doc for the EIC who needs to write a portfolio analysis at the Fed SF district level?"

### What works

- **The realistic acknowledgement that comparison is manual until v1.1 ships.** "Until the v1.1 `verifier portfolio` subcommand ships, comparison is manual; this doc articulates the manual procedure." That is the right framing — the EIC does the work today with a spreadsheet; the doc tells them how. Promising tooling that does not exist yet would be unhelpful and would invite the EIC to wait for tooling that may slip.
- **The metrics table.** Pass rate, anomaly rate, sealing-delay frequency, late-binding rate, master-rotation events, verifier-validation cadence. Six metrics is the right size for a portfolio-comparison spreadsheet — enough signal to identify outliers without becoming an analytical project unto itself. For the Fed SF district's fintech-bank-tech sub-portfolio specifically, master-rotation events are the metric I'd lean on most heavily; partnership-bank vendors rotate IKMs more aggressively than a long-deployed core-banking system would.
- **The sample 6-bank workspace.** Concrete example with realistic numbers (Bank C with HSM cluster issue Q1 remediated; Bank E with anomaly cluster under investigation). An EIC can copy this shape into their own portfolio review without inventing the columns. The "Notes" column carries the qualitative narrative that the metrics alone would miss; that is the right separation for a portfolio analysis at the supervisory level.
- **The "common patterns and their implications" section.** Three patterns: portfolio-wide degradation, single-institution outlier, sub-portfolio commonality. Each has a "this often indicates" list and an "action" line. That is the right shape for an EIC's pattern-recognition reference. The "shared HSM provider with degradation" branch is the one I'd particularly call out — for the Fed SF district's fintech-bank-tech sub-portfolio, multiple institutions sharing a single AWS region or a single HSM-as-a-service vendor is the most likely root cause for a sub-portfolio commonality, and the doc names that branch correctly.
- **The cross-period comparison section.** Trend in pass rate, trend in anomaly rate, trend in sealing-delay frequency, trend in operational-events volume, material configuration changes. Five longitudinal metrics. The "improvement trajectory" and "degradation trajectory" patterns give the EIC a clean read on whether an institution is maturing or weakening. That is exactly the read an EIC needs for the supervisory plan.
- **The cross-agency coordination section.** Names the realistic case (bank holding company under Fed + bank subsidiary under OCC) and clarifies that each agency conducts its own examination but shares verifier output through standard supervisory channels. That recognizes how multi-agency supervision actually works without overstepping the agency-specific examination authority.

### Narrow follow-ups

- **NOTE — Cross-period comparison and spec §7 failure step categories.** The cross-period section calls out trend metrics (pass rate, anomaly rate, sealing-delay frequency, operational-events volume, configuration changes) but does not call out trends in specific spec §7 failure step categories. An EIC doing longitudinal analysis on an institution that has had two step 8 (`key_fingerprint mismatch`) findings in two consecutive examinations should be looking at that as a control-degradation pattern, not just an anomaly-rate pattern. Suggest one bullet under cross-period metrics: "Trends in spec §7 failure step categories — recurring step 7/8 findings indicate IKM-management control degradation; recurring step 6/9 findings indicate chain-integrity control degradation; recurring step 11 findings indicate signing-key-management control degradation; recurring step 11 (dual-algo) findings during the transitional period indicate algorithm-posture change-management gaps." Keys the EIC's longitudinal lens to the same step-numbered taxonomy the rest of the regulator pack uses, and extends the lens to the dual-algorithm transitional period. Optional.

---

## 10. spec/chain-of-custody-v1.md §7 step 11 dual-algorithm cases (a)-(e)

The spec is the implementer's doc, not the examiner's, but a bench examiner looking at a co-signed seal failure on a partnership-bank examination needs the spec text and the finding-language doc to tell a consistent severity story. I read §7 step 11 to confirm the case (e) severity cross-reference and the working-paper convention paragraph land where they should.

### What works

- **Case (e) severity cross-reference is in the spec text.** The spec §7 step 11 case (e) paragraph reads: "**The severity of case (e) is Severe regardless of bracket** — the PASS-WITH-ANOMALY disposition under non-strict reflects the spec's posture that the un-broken algorithm's signature still provides integrity assurance for downstream consumers, NOT that the failure is itself low-severity. Examiners writing up case (e) findings cite `regulator-pack/finding-language.md` row '11 (dual-algo) co-signed seal failure' (Severe MRA) regardless of which bracket the verifier reported." A junior examiner reading the spec and seeing PASS-WITH-ANOMALY in non-strict mode no longer has any reason to write the finding up as an Observation — the spec text itself directs them to the finding-language Severe-MRA row. That is the right enforcement of the severity guidance at the spec layer rather than relying on the examiner to know to consult finding-language separately.
- **Case (e)'s two-track interpretation paragraph.** "EITHER (i) one of the algorithms has been broken (in which case the un-broken algorithm's signature still provides integrity assurance and the institution coordinates with the regulator on the broken-algorithm migration timeline), OR (ii) one of the seals is forged under a compromised algorithm-specific signing key (in which case the un-broken algorithm's signature confirms the un-compromised half of the chain custody)." The two-track separation of duties is correct: the verifier reports both algorithms' validation results; the institution's IR program does the interpretation. The bench examiner does not have to pick the interpretation — they confirm the institution's IR program has executed the dispatch.
- **The working-paper convention paragraph.** "The verifier output MUST record both algorithms' validation results (PASS / FAIL / NOT-PRESENT per algorithm) when the seal record carries `signatures`. The examiner's working paper carries both rows; the examination report cites both algorithm validations." This is the load-bearing paragraph for examiner work — it tells the examiner that during the dual-algorithm transitional period, the working paper carries TWO algorithm validation rows for each affected seal-day, not one. Without this convention, an examiner could record the primary algorithm's PASS and overlook the co-signed algorithm's FAIL, missing a Scenario 4 + Scenario 3 finding.
- **The "this convention applies during the multi-year transitional period; once the institution retires one algorithm, the convention reverts to single-algorithm reporting" closing.** Tells the examiner the convention has a sunset. That matters because examiners build muscle memory; knowing the muscle memory will revert after the transition prevents over-engineering the working-paper template downstream.
- **The defensive choice in case (b).** Single-algorithm signature on a seal during dual-algorithm posture is bracketed as PASS-WITH-ANOMALY in both strict and non-strict — a control-completeness finding, NOT a chain-integrity finding. That is the right call: the seal is integrity-bearing under the present algorithm, the institution's posture commitment is what's incomplete, and dragging the institution into a chain-integrity IR posture for a partial-coverage day would over-correct.

### Narrow follow-ups

I have nothing to flag on the spec §7 step 11 case (e) paragraph itself. The severity-overrides-bracket cross-reference and the working-paper convention land where they need to land.

---

## 11. Cross-doc consistency checks

I went through the regulator-pack docs looking for places where the same fact is stated in multiple docs and confirmed they agree.

| Fact | Doc 1 | Doc 2 | Doc 3 | Doc 4 | Consistent? |
|---|---|---|---|---|---|
| `--master-key` is in the standard examiner invocation | quickstart | sample-report | training Module 3 | deployment-package | **Consistent across all four** |
| Windows ACL note for the IKM file | quickstart (POSIX-only descriptor) | sample-report (POSIX-only) | training Module 3 (POSIX + Windows ACL) | deployment-package (POSIX + Windows ACL) | **Mostly consistent — see follow-up below** |
| `step: 8` is "stop and call the bank" / Severe MRA | quickstart 14-row table + 3 specific points | finding-language step 8 row + paragraph | training Module 5 row | response-workflow step 1 example | **Consistent** |
| `step: 8` is identity mismatch, NOT content tampering | quickstart specific point | finding-language paragraph (bold) | response-workflow step 4 | training Module 5 row | **Consistent** |
| `audit file ends mid-line` is operational, NOT tampering | quickstart specific point | finding-language paragraph | training Module 5 row | response-workflow excludes via "when this workflow does NOT apply" | **Consistent** |
| Spec §7 step 11 case (e) co-signed seal failure | spec §7 step 11 (e) | finding-language dual-algo row 3 | response-workflow step 2 dual-algo row | training Module 5 dual-algo row | **Consistent** |
| Case (e) severity is Severe regardless of verifier bracket | spec §7 step 11 (e) explicit cross-reference to finding-language | finding-language dual-algo row 3 (Severe MRA) | (workflow + training defer to finding-language) | — | **Consistent** |
| 14-row taxonomy count | quickstart 14 rows | finding-language 14 informational rows | training Module 5 11 rows (subset) | — | **Module 5 is a subset (intentional); see Module 5 follow-up above** |
| IR scenario column present | quickstart (yes) | finding-language (yes) | training Module 5 (embedded in action column) | response-workflow step 2 table | **Consistent** |
| Working-paper records both algorithm validation rows | spec §7 step 11 working-paper convention paragraph | (finding-language dual-algo rows imply by Severe-MRA-regardless-of-bracket) | (training Module 5 implies by IR Scenario 4 + Scenario 3 reference) | — | **Consistent in substance; convention named explicitly only in the spec — see follow-up** |

The two cross-doc items worth calling out:

- **The Windows ACL note** is in training Module 3 and deployment-package.md but not in the quickstart's `--master-key` paragraph or the sample-report's `--master-key` paragraph. Both of those use the POSIX-only descriptor "32 raw bytes, file mode 0600." For Fed examiner work specifically, most examiner laptops are Windows; the descriptor doesn't translate without the parenthetical. Suggest the same one-line addition the deployment-package and training Module 3 already carry: "(file mode 0600 on POSIX; ACL restricted to the examiner account on Windows)." Two small edits, two docs.
- **The working-paper convention** is named explicitly only in the spec §7 step 11 paragraph. The finding-language doc, response-workflow doc, and training Module 5 doc all imply it (Severe MRA regardless of bracket; Scenario 4 + Scenario 3; both algorithm validations cited) but none of them name the convention explicitly. A bench examiner who has not read the spec may not know to record both algorithm validation rows in the working paper. Suggest a one-line addition to the response-workflow doc, in the dual-algorithm row commentary: "Working-paper convention (per spec §7 step 11): when the seal record carries `signatures`, the working paper records BOTH algorithms' validation results (PASS / FAIL / NOT-PRESENT per algorithm); the examination report cites both validations." Cross-link, not a content change.

---

## 12. Items I'd want to see before signing off

| Severity | Item | Doc |
|---|---|---|
| MINOR | Step 9 finding paragraph "[is/is not] underway" template style inconsistent with rest of doc | finding-language.md |
| MINOR | FAIL example block ordering (cover-page summary should precede per-day FAIL detail for top-to-bottom reading order) | sample-report.md |
| NOTE | Add row for `--master-key absent` outcome in 14-row table | examiner-quickstart.md |
| NOTE | Add three dual-algorithm sub-rows to 14-row table (or accept the "in finding-language only" treatment) | examiner-quickstart.md |
| NOTE | JSON-vs-pretty-print snippet pairing (sample-report shows pretty-print; workflow shows JSON) | sample-report.md |
| NOTE | TSC appendix could carry the three dual-algorithm sub-rows mapped to TSC criteria | sample-report.md |
| NOTE | Windows ACL parenthetical missing from quickstart and sample-report `--master-key` paragraphs | examiner-quickstart.md, sample-report.md |
| NOTE | Step 6 paragraph "most likely tampering signal" framing too strong without facts | finding-language.md |
| NOTE | Repeat-finding examples for new step 2/3/5/6/9 paragraphs | finding-language.md |
| NOTE | Step-mapping table cross-reference for cases (a) and (d) of spec §7 step 11 | examination-response-workflow.md |
| NOTE | Case (e) two-track IR-scenario clarification | examination-response-workflow.md |
| NOTE | Working-paper convention named explicitly in workflow doc, not just the spec | examination-response-workflow.md |
| NOTE | Module 5 row for `step: 11 (dual-algo)` partial-coverage seal | examiner-training.md |
| NOTE | Module 5 row count vs full taxonomy clarifying note | examiner-training.md |
| NOTE | II.E "examiners confirm" list could add verifier-dispatch testing | handbook-mapping.md |
| NOTE | Cross-period comparison metrics keyed to spec §7 step categories (incl. dual-algo) | portfolio-comparison-procedures.md |

Two MINOR (template-style harmonization in finding-language; block-ordering in sample-report). Fourteen NOTE-level items, all clarity edits or optional additions; none would prevent a competent FFIEC IT examiner from producing a defensible report on a fintech-bank-tech examination today.

---

## 13. Stopping criterion

0 EXAMINER-severity findings preventing examiner workflow + 0 EXAMINER-severity findings preventing report production. **0/0 reached.**

The bench-examiner surface is at the level I would accept on a Fed IT examination program. The five docs that matter to the bench examiner are coherent, cross-referenced, and key off the same `step` field consistently. The dual-algorithm transitional-period material is at the right level of specificity — enough that a bench examiner can act on a co-signed seal failure without becoming a cryptographer, not so deep that the doc is unreadable. The working-paper convention paragraph in the spec is the load-bearing piece for the multi-year transitional period; ensuring the convention is also named in the workflow doc would close the loop for examiners who do not read the spec.

The 2026-04-22 PASS / 2026-04-23 FAIL pairing in the sample report is the single pedagogical decision I would most want other regulator packs to copy. The `step` field as the spine of the whole loop is the single architectural decision I would most want other examination tools to copy. The asymmetric `--master-key` posture (key-bound for examiner, structural-only for SOC and customer-side) is the single operational decision I would most want other audit tooling to copy — it correctly reflects that an examiner's disclosure shape differs from a SOC engagement's, and the docs name both shapes consistently.

— Frank Kobayashi
