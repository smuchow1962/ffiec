# Round 11 — FFIEC IT Examiner review

**Reviewer.** Devraj Patel, Senior FFIEC IT Examiner, OCC. Eighteen years in bank IT examination; current rotation through the FFIEC's joint AI Examination Guidance working group.

**Lens.** First-look reviewer. I run the verifier on examiner-laptop in real engagements; I read the quickstart, training, sample report, finding language, deployment package, and the response workflow as one integrated package. I have no insider history of prior iterations.

**Verdict.** Close to ready. The package is more usable than what I have seen from any prior chain-of-custody artifact in my career — the five-step response workflow is the kind of document that turns a panicked examiner into an effective one in 60 seconds. But there are real gaps in the finding-language paragraph coverage, an off-by-counting issue in the "12-row" failure modes table, and the spec §7 step 11 dual-algorithm text leaves an examiner without a quickstart row, a finding paragraph, or a workflow path when the new failure modes fire. Stopping criterion (0 gaps + 0 partials) is not met.

---

## Strengths I want on the record

The shape of the package is right.

- **The examination-response-workflow doc is the single best artifact in the package.** Five steps, ~60 seconds of reading, JSON record in / documented response out. The second worked example for step:10 was a smart addition — it proves the workflow generalizes across structurally different failure types (identity mismatch, content tampering, ledger-content mismatch, signature compromise) instead of being a step:8 special case. The "mismatch direction tells the auditor whether the seal record was altered or the ledger contents were altered or recovered from a backup" sentence in the step:10 example is the kind of forensic guidance examiners actually need at the moment of decision.
- **The PASS day at 2026-05-07 adjacent to the FAIL day at 2026-05-08** lets a new examiner see the structural difference (which fields populate, which read N/A, which carry the failure record) without flipping between two reports. The "Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)" line is small but load-bearing — it tells the reader exactly where the verifier stopped.
- **The cover-page summary changes accordingly block** (lines 263-280 of sample-report.md) is the right teaching device. New examiners need to see how a single failed day flows from per-day detail to cover-page summary to FAILED DAYS table to "investigate institution's IKM roster row" pointer to the response-workflow doc. That trail is intact.
- **The "stop and call the bank" framing** in the quickstart now correctly covers BOTH step 8 and step 9 with the right distinction (identity mismatch investigated through the IKM roster vs. content tampering investigated through the chain content). Before this language, an examiner could plausibly conflate the two and waste investigation cycles on the wrong evidence path.
- **Spec §7 step 11 carrying explicit `algorithm/key-type mismatch` reporting** is a courteous touch — examiners do not have to puzzle out from a generic "signature verification failed" whether the issue is key-curve mismatch or actual signature failure.
- **`audit_file.truncation_detected` consistently framed as operational, not tampering** across quickstart, finding-language, response-workflow, and IR scenario list. This is the kind of consistency that prevents an examiner from writing a Severe MRA when the actual root cause is an SDK process crash.
- **Handbook-mapping II.C.13 cryptographic-controls breakout** correctly enumerates the four examiner-confirmable items: IKM length ≥ 32, valid `key_fingerprint` bytes verified at step 8, no `plaintext-dev` URI in production, weekly reconciliation per §10.1. That converts an abstract "cryptographic controls" objective into a four-item examiner checklist.

---

## Gaps (must address before stopping)

### G1. Finding-language paragraphs do not cover all 12 steps

The prompt says "finding paragraphs all 12 steps." The severity table at lines 11-31 of finding-language.md is keyed by step number. But the "Sample finding paragraphs" section (lines 33-82) includes paragraphs for only some steps:

| Step | Severity row exists | Finding paragraph exists | Notes |
|---|---|---|---|
| 1 | yes | yes ("Format version not supported") | OK |
| 2 | yes | **no paragraph** | Severity row says "Severe MRA — investigate whether SDK or verifier constants drifted" but no candidate paragraph for the report |
| 3 | yes | **no paragraph** | Severity row says "Severe MRA; format-construction defect" but no paragraph |
| 4 | yes | yes ("Cross-chain lift detected") | OK |
| 5 | yes | **no paragraph** | Severity row says "MRA; mid-file format drift" but no paragraph |
| 6 | yes | **no titled paragraph by step** | The "Chain hash mismatch" paragraph at line 36 is generic — it does not name the step or distinguish the structural-walk failure from the MAC failure |
| 7 | yes | yes ("Unknown key_version") | OK |
| 8 | yes | yes ("Key fingerprint mismatch") | OK |
| 9 | yes | **no titled paragraph by step** | The "Chain hash mismatch" paragraph is the closest match but is not titled by step number, conflates structural and MAC failures, and does not match the step-keyed convention the rest of the section uses |
| 10 | yes | yes ("Merkle root mismatch (single day)") | OK |
| 11 | yes | yes ("Signature verification failed") | OK |
| 12-cadence | yes | yes ("Cadence mismatch") | OK |
| 12-dev-mode | yes | yes ("Software-key fallback used in production") | OK |
| pre-flight truncation | yes | yes ("Audit file ends mid-line") | OK |

Missing paragraph coverage: **steps 2, 3, 5, 6, 9** (five steps). The "Chain hash mismatch" paragraph at the top of the section reads as a leftover from before the step-keyed convention was adopted; it needs to be either retitled, split into per-step paragraphs, or removed in favor of step-numbered paragraphs.

The promise of the doc is: "examiner picks the paragraph that fits, edits for institution-specific facts, and inserts into the examination report." That promise breaks for an examiner who sees a step:6 `chain link broken at seq N` failure record — they have to write the finding paragraph from scratch or stretch the generic "Chain hash mismatch" language. Stretching is what produces inconsistency across the regulator's portfolio, which is exactly what the doc says it exists to prevent (line 4: "consistency across the regulator's portfolio matters more than the specific words").

**Required fix.** Add five named, step-keyed finding paragraphs:

- "Header HKDF inputs mismatch (spec §7 step 2)"
- "Header genesis hash mismatch (spec §7 step 3)"
- "Format version mismatch mid-file (spec §7 step 5)"
- "Chain link broken / seq out of order (spec §7 step 6)"
- "Payload hash MAC mismatch (spec §7 step 9)"

The existing generic "Chain hash mismatch" paragraph should be removed or explicitly noted as deprecated in favor of the step-keyed paragraphs. Examiners working from this doc need one row in the severity table to map cleanly to one paragraph in the sample-paragraph section.

### G2. Dual-algorithm transitional period (spec §7 step 11) is invisible to the examiner workflow

Spec §7 step 11 introduces three new verifier outcomes for the dual-algorithm transitional period:

- "partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture" (PASS-WITH-ANOMALY)
- "algorithm not on institution's declared posture list at seal_date {D}" (FAIL under --strict; PASS-WITH-ANOMALY otherwise)
- Co-signed seal carrying both algorithms (PASS)

None of these reason strings appear in:

- the quickstart 12-row failure modes table (the only step-11 row is "signature verification failed" / "algorithm/key-type mismatch")
- the finding-language severity table (same — only the existing two strings)
- the finding-language sample paragraphs (no paragraph for partial-coverage or undeclared-algorithm)
- the examination-response-workflow step-mapping table (the step:11 row maps to Scenario 3, but Scenario 3 is "signature verification failed" — not "the institution's posture commitment is incomplete on this seal-day")
- the handbook-mapping doc (II.E change management mentions `format_version` change management but not algorithm-posture change management; II.C.13 mentions "Ed25519 in HSM custody" but not the dual-algorithm transitional posture)

For an examiner running the verifier in 2027 against a bank that has begun a Dilithium pilot, the verifier will emit one of the three new reason strings, and the examiner will have no documented mapping from the JSON failure record to a severity guidance, an IR scenario, a finding paragraph, or a handbook-mapping section. The spec text exists; the examiner package does not yet reflect it.

**Required fix.** Three coordinated additions:

1. **Quickstart 12-row table.** Add three rows under step 11:
   - `partial-coverage seal: single-algorithm signature during dual-algorithm posture` → severity "Observation (control-completeness, NOT chain-integrity)"; IR scenario "(institution's posture commitment, not chain failure — direct to control description)"
   - `algorithm not on institution's declared posture list at seal_date {D}` (under --strict) → severity "Severe MRA"; IR scenario "Scenario 3 variant + algorithm-posture investigation"
   - `algorithm not on institution's declared posture list at seal_date {D}` (without --strict) → severity "Observation"; IR scenario same as above
2. **Finding-language paragraphs.** Add two named paragraphs:
   - "Partial-coverage seal during dual-algorithm posture (spec §7 step 11)" — explains the institution committed to dual-algorithm posture in its control description but produced a single-algorithm seal on this day; this is a control-completeness finding, NOT a chain-integrity finding; the institution's investigation path is its algorithm-rotation runbook, not the chain content.
   - "Algorithm not on declared posture list (spec §7 step 11)" — explains the seal carries an algorithm the institution did not declare in its posture commitment; under --strict, this fails; otherwise it surfaces as anomaly. The institution either updates its declared posture or removes the unauthorized algorithm.
3. **Examination-response-workflow step-mapping table.** Split the step:11 row into three sub-rows so the examiner sees that "signature verification failed", "algorithm/key-type mismatch", and the new dual-algorithm reason strings have different IR scenarios and different evidence paths.
4. **Handbook-mapping II.E.** Add explicit text: "Algorithm rotation across the dual-algorithm transitional period (spec §4.3.2 / §7 step 11) is a change-management event; institutions document their declared algorithm-posture list, the timeline of any addition/removal, and the change-management approvals. Examiners should confirm the institution's `algorithm` field in seal records aligns with the declared posture list across the examination period."

Without these, an examiner running the verifier in the dual-algorithm window will see a reason string they cannot map to a documented response. That is the exact scenario the response-workflow doc was built to prevent.

### G3. The "12-row failure modes table" is not 12 rows

The quickstart says "12-row failure modes table" (line 55: "The verifier follows the spec §7 twelve-step procedure. Each step's failure produces a specific named reason."). The actual table has 14 rows: steps 1-11 (11 rows), step 12 cadence (1 row), step 12 dev-mode (1 row), file pre-flight (1 row) = 14 rows.

This is not a deal-breaker, but it is the kind of inconsistency a SOC reviewer will flag. Either the framing should change ("the failure modes table covers each spec §7 step plus the file pre-flight refusal — 14 named failure modes across 12 procedure steps"), or step 12 should be presented as one row with the two sub-cases compacted, or the pre-flight row should be moved to a separate table. The same off-by-row issue propagates into finding-language.md severity table (which is also presented as the "12-step" table but has 15 rows once you include both anomaly entries).

I would not block the package on this if the other gaps were resolved, but per the stopping criterion, every partial counts.

---

## Partials (close before stopping)

### P1. Sample-report PASS/FAIL pairing has internal inconsistency the disclaimer does not fully cover

The "Sample failed-day output (page N+3 — paired example)" block at lines 213-280 of sample-report.md inserts hypothetical 2026-05-07 PASS and 2026-05-08 FAIL days into a report whose verification range is 2026-04-01 to 2026-04-30 (per cover page line 64 and the verifier invocation at lines 21-29). The disclaimer at line 214 ("Hypothetical; not part of the 30-day pass example above") covers the dating mismatch, but the "Cover page summary changes accordingly" block at lines 263-280 then shows `Days verified: 30, Days passed: 29, Days failed: 1` — that is a 30-day total but with a 2026-05-08 failure that is outside the 2026-04-01..2026-04-30 range.

A new examiner reading this in sequence is confused for ~30 seconds before they parse the disclaimer. Two cleaner options:

- **Option A.** Make the paired example a fully self-contained report with its own date range (e.g. 2026-05-01 to 2026-05-31, 31 days, 30 passed, 1 failed). The example then stands alone as "what a single-failure report looks like."
- **Option B.** Move the paired example to a separate section ("Example: a single-day failure in a different examination period") and present it explicitly as a different report, not as a "summary changes accordingly" view of the original report.

Option A is cleaner pedagogically; the report should be internally consistent on dates.

Also: the FAIL example uses `key_version: 1` while the original 30-day setup is on `key_versions [3]` and `[3, 4]`. The disclaimer covers this, but Option A would let the example use a `key_version` value consistent with its own setup, removing one more piece of cognitive load.

### P2. Examination-response-workflow doc does not cover step:11 dual-algorithm

Per G2 above. The five-step path is exactly right for the existing failure modes; it needs the dual-algorithm sub-cases to remain complete once spec §7 step 11's transitional text is in force. Listed as partial because the workflow doc's structure works — what is missing is the dual-algorithm rows in the step→IR-scenario table at lines 38-52, and a third worked example (or a sub-bullet in step:3 evidence collection) for "what evidence to ask the institution for when the seal carries a partial-coverage or undeclared-algorithm signature."

The evidence path is: the institution's algorithm-posture declaration (control description + any change-management records since), the institution's `master_key.rotated` events for any signing-key changes, and the institution's dual-algorithm signing-job logs (which signing path executed and why only one signature was produced). Two questions of the form the workflow doc uses elsewhere would close this.

### P3. Quickstart "What you need" + verifier invocation does not mention `--master-key`

The quickstart at lines 9-41 shows the standard invocation with `--ledger`, `--root-key`, `--tenant-id`, `--from`, `--to`, `--report`, `--bundle`. It does not mention `--master-key`. But the verifier-design doc §4.4 makes clear that without `--master-key`, the verifier degrades to structural verification only (PASS-WITH-ANOMALY: "structural verification only; key-bound verification skipped"); under `--strict`, missing `--master-key` is a FAIL.

For an examiner running their first verification on a bank that has not yet provided IKM access (a real first-engagement scenario), the report will come back PASS-WITH-ANOMALY with a confusing message and the examiner will not know whether to escalate. A two-sentence note in "What you need" — "If the bank has provided IKM access, also pass `--master-key <path>`. Without it, the verifier degrades to structural verification only and reports PASS-WITH-ANOMALY; consult the design doc §4.4 for handling" — would close this.

The sample report's standard invocation at lines 21-29 has the same omission. The SOC `--strict` invocation at lines 36-44 also does not show `--master-key`, which would actively FAIL under `--strict` in master-key-less mode.

### P4. Handbook-mapping II.E change-management coverage incomplete for algorithm rotation

Per G2 fix #4. II.E currently covers `format_version` change management thoroughly but not `algorithm` change management. The spec's §4.3.2 algorithm-rotation text and §7 step 11's dual-algorithm dispatch make algorithm change management a first-class concern; the handbook-mapping needs a sentence or two acknowledging that and pointing examiners to where in the institution's documentation to look (the algorithm-posture declaration, the change-management records for any addition/removal of an algorithm).

### P5. Quickstart and finding-language IR scenario column / severity-table column do not match the response-workflow's step→scenario table on every row

Quickstart row for step 6 says "Scenario 1." Response-workflow row for step 6 says "Scenario 1 (chain hash mismatch)." Quickstart row for step 4 says "Scenario 1 + evidence-handling investigation." Response-workflow row for step 4 says "Scenario 1 + evidence-handling investigation (mis-bundled snapshot OR lift attempt)." These match.

Quickstart row for step 12 (cadence) says "Scenario 5 variant." Response-workflow row says "Scenario 5 (sealing delay) variant." OK.

But the finding-language severity table (lines 11-31) has NO IR-scenario column at all. An examiner working from finding-language.md alone has to mentally cross-reference back to the quickstart or the response-workflow to know the IR scenario. Adding an IR-scenario column to the finding-language severity table would make finding-language self-contained for the field examiner who already knows the failure mode and is going straight to the paragraph-lift step.

### P6. Examiner-training Module 5 "Common patterns" table predates the step-keyed taxonomy

The Module 5 table at lines 87-95 has rows like "One or more days fail with chain hash mismatch" — which is the pre-rework framing. Today's verifier produces step-keyed reasons (`chain link broken at seq N`, `payload_hash MAC mismatch at seq N`, etc.) and the examiner is expected to map the step number, not a free-text "chain hash mismatch" string. The training table should be updated to use the step-keyed taxonomy so a new examiner trains on the same vocabulary they will see in the JSON report.

A small but meaningful gap: a new examiner trained on "chain hash mismatch" will look at a `chain.verification_failure` JSON record carrying `step: 6, reason: "chain link broken at seq 47"` and not immediately know whether that maps to "chain hash mismatch" in their training. Two minutes' worth of reframing closes this.

---

## Items I checked and found OK

I want to be explicit about what I did NOT find issue with, since the stopping criterion is 0 gaps + 0 partials and the absence of a mention should not be ambiguous.

- **Quickstart "stop and call the bank" framing covering both step 8 and step 9.** Correct. The investigation-path distinction (IKM roster vs. chain content) is exactly the right framing.
- **Sample-report rotation day at 2026-04-15.** Correct. `Key versions present: [3, 4]`, both fingerprints listed, HMAC chain walk explicitly notes "events under v3 verified against ikm_v3; events under v4 verified against ikm_v4," anomaly section explains the rotation completed 03:14 UTC per institution incident log. This is the cleanest rotation-day example I have read in any chain-of-custody document.
- **Sample-report methodology refresh.** The order-of-operations block (lines 180-198) now matches spec §7's twelve-step procedure exactly. The PASS/FAIL rules (lines 200-209) correctly distinguish "step 12 cadence mismatch → PASS with anomaly" from "step 12 dev_mode under --strict → FAIL," which is examiner-essential.
- **Sample-report --strict invocation.** Lines 36-44 correctly show `--strict` with the explanatory note "Use --strict for substantive SOC testing; without it for FFIEC examiner work where anomalies are evaluated in operational context." This is the right framing — the FFIEC examiner does NOT default to --strict; the SOC engagement does.
- **TSC criterion appendix.** The verifier-output-line → TSC-criterion mapping at lines 305-321 is the kind of cross-walk the SOC team needs. The mappings are sensible (PI1.1 for chain-walk, PI1.2 for Merkle, CC6.7 for HSM, CC6.8 for unauthorized-modification defenses).
- **Examination-response-workflow second worked example (step:10 Merkle root mismatch).** Coverage is complete — JSON record, severity mapping, IR scenario mapping, evidence collection (seal-job log questions), finding paragraph, four-item institution response. The "mismatch direction" explanation (computed correct + recorded wrong vs. recorded correct + computed wrong) is the kind of forensic interpretation guidance examiners will use directly.
- **Handbook-mapping II.C.10 logging coverage.** The "weekly key-fingerprint reconciliation per spec §10.1" addition correctly elevates the logging mapping beyond "the chain provides logging integrity" into "the chain plus its operational reconciliation provide ongoing key-identity assurance with bounded detection window." Examiners working II.C.10 now have a concrete operational evidence artifact to ask for (`master.reconciliation_completed` events) and a documented cross-reference into audit-procedures P-6.
- **Handbook-mapping II.C.13 cryptographic controls four-item checklist.** Four examiner-confirmable items, each with a spec-section pointer. Converts an abstract control objective into operational evidence questions.
- **Spec §7 step 11 single-algorithm posture default for v1.0.** The text correctly states "the dual-algorithm dispatch reduces to single-algorithm verification; PASS / FAIL on the single signature" — meaning v1.0 deployments today see no behavior change. This is the right transitional design from a regulatory-adoption perspective; institutions don't have to do anything new in v1.0 to be conformant with the §7 step 11 dual-algorithm text.
- **Deployment package allowlisting postures.** Both SHA-256 hash allowlist and cosign signature allowlist are presented as conformant with FFIEC IT shop standards, which matches the heterogeneity I see across the regulator's IT shops. The "many regulators run both" note is realistic.

---

## What "ready" looks like from this lens

Three concrete additions:

1. **Five new finding-language paragraphs** for steps 2, 3, 5, 6, 9. The generic "Chain hash mismatch" paragraph either retired or split. (G1)
2. **Dual-algorithm transitional-period coverage** in quickstart, finding-language, response-workflow, and handbook-mapping. Three new failure-mode rows, two new finding paragraphs, one third worked example (or a substantive sub-bullet under the second), one handbook-mapping addendum. (G2)
3. **Counting consistency** on the "12-row" table framing — either reframe as "14 named failure modes across 12 spec §7 steps" or compact step 12 into one row with two sub-cases. (G3)

Plus the partials, which are individually small but cumulatively change the package from "examiner-usable" to "examiner-trustworthy under any failure mode the verifier can produce":

4. Sample-report paired example as self-contained report with its own date range. (P1)
5. Workflow-doc third worked example or sub-bullets for dual-algorithm. (P2)
6. `--master-key` mention in quickstart and sample-report invocations. (P3)
7. Handbook-mapping II.E algorithm-rotation change-management addendum. (P4 — covered by G2 fix #4)
8. IR-scenario column added to finding-language severity table. (P5)
9. Examiner-training Module 5 retitled to step-keyed taxonomy. (P6)

When all nine are in, an examiner running the verifier on day one of an engagement can take any JSON failure record the verifier produces, walk the five-step workflow, lift the right finding paragraph, document the IR scenario, and have a defensible report at the end. That is the bar. We are not at the bar yet.

---

## Files I read

- `E:/dev/ffiec/spec/chain-of-custody-v1.md` (full read; §7 step 11 close read)
- `E:/dev/ffiec/docs/examiner-quickstart.md`
- `E:/dev/ffiec/docs/regulator-pack/sample-report.md`
- `E:/dev/ffiec/docs/regulator-pack/finding-language.md`
- `E:/dev/ffiec/docs/regulator-pack/handbook-mapping.md`
- `E:/dev/ffiec/docs/regulator-pack/deployment-package.md`
- `E:/dev/ffiec/docs/regulator-pack/examiner-training.md`
- `E:/dev/ffiec/docs/regulator-pack/examination-response-workflow.md`
- `E:/dev/ffiec/docs/portfolio-comparison-procedures.md`
- `E:/dev/ffiec/docs/design/07-verifier-design.md`

3 gaps. 6 partials. Stopping criterion not met.
