# Round 15 — FFIEC IT Examiner review

**Reviewer:** Elizabeth O'Donnell, Senior IT Examiner, OCC Boston Field Office (FFIEC IT cadre, 19 years)
**Scope of this read:** spec §7 step 11 + step 12a; `examiner-quickstart.md` 15-row failure-mode table (now including step 12a + bundle sufficiency); `regulator-pack/sample-report.md`; `regulator-pack/finding-language.md` (12a row, counsel-rebuttal framing, CAT crosswalk note, cadence-mismatch repeat threshold); `regulator-pack/handbook-mapping.md`; `regulator-pack/deployment-package.md`; `regulator-pack/examiner-training.md` (Module 5 step 12a row); `regulator-pack/examination-response-workflow.md`; `portfolio-comparison-procedures.md`.
**Stopping criterion:** 0 blockers / 0 should-fix.

---

## Headline

The package is in production-usable shape from an FFIEC IT-examiner standpoint. I can take this set of documents to a field examination of a mid-sized BYOC institution tomorrow, run the verifier on the bank's snapshot, produce a report, write up findings using the prepared finding-language paragraphs, and defend the report at exit conference. Every step of that workflow is covered by exactly one doc that names itself as the place to look — quickstart for orientation, deployment-package for laptop allowlisting, training for the 30-minute onboarding, sample-report for the artifact shape, finding-language for the report-paragraph wording, examination-response-workflow for the 60-second JSON-to-finding path, handbook-mapping for the Handbook crosswalk, portfolio-comparison-procedures for the EIC's cross-bank read.

I have **no blockers and no should-fix items**. The two narrow observations below are file-and-forget; they do not gate adoption.

---

## What works for an FFIEC IT examiner

### The five-step path is the backbone, and it holds

`examination-response-workflow.md` is the single doc I will print and keep in my examination kit. The five-step path — `step` → severity, `step` → IR scenario, pull the institution's reconciliation/seal-job/HSM evidence, lift the finding paragraph, confirm the institution's four-item response — is exactly the path I'd describe if you asked me to teach a junior examiner how to handle a chain-detected finding. The fact that the workflow doc explicitly says "this doc exists so the examiner does not have to assemble the workflow from four other docs at the moment of finding evaluation" tells me the author understands what an examination day actually feels like. We do not have time, in front of a bank's chain-operations team at 2pm, to cross-reference four documents.

The two worked examples (step 8 `key_fingerprint mismatch` and step 10 `merkle root mismatch`) are well chosen — they're structurally different failure types (identity mismatch vs. ledger-content mismatch), and the workflow walks through both without needing custom handling. That generalizes the path credibly for the other steps I haven't seen worked.

### Step 11 (signature verification) is examined cleanly

Step 11 is the failure mode my management would push hardest on if it ever fired in my portfolio. Possible HSM key compromise reads as the highest-impact line item in the entire examination response. The spec §7 step 11 text handles this carefully:

- The single-algorithm path is unambiguous: signature verifies against the tenant public key, dispatching on the seal record's `algorithm`. Mismatch → Severe. The verifier emits `algorithm/key-type mismatch` rather than the generic message when the seal's claimed algorithm and the resolved public key disagree, which is the right diagnostic — it tells me which way the mismatch went without my having to reason about it.
- The dual-algorithm transitional period text (cases a–e) is a substantial piece of the spec, and it lands. Case (e) — both signatures present, one valid + one invalid — is the case I had to read twice, and the spec anticipates that. The bracket explicitly says "the severity of case (e) is Severe regardless of bracket" and tells me to cite the row in finding-language.md regardless of which bracket the verifier reported. That's important. Without that bracket-bridging language, an examiner reading a non-strict PASS-WITH-ANOMALY for case (e) might under-write the finding; the spec closes that gap by name.
- The examiner working-paper convention is also right: both algorithms' validation results recorded in the working paper, both algorithm validations cited in the examination report. That gives me both rows on the day, and gives the next examiner (in a multi-year transitional period) the full picture without needing to reconstruct it.

The four counsel-rebuttal patterns in `finding-language.md` for step 11 ("the registry has the wrong key, not the HSM has been compromised") are the kind of preparation I rarely see in regulator-pack documentation. Counsel will rebut that way; I've heard it. Having the prepared response — registry mismatch is one of four plausible IR Scenario 3 root causes; in the absence of HSM-side audit-log evidence supporting the registry-mismatch hypothesis, the finding remains Severe pending HSM-side investigation — saves me from making the response up at the moment of rebuttal.

### Step 12a (`gen_ai_model_identifier_missing`) lands well

This is the new failure mode in the package I read, and it's positioned exactly right.

The 15-row failure-mode table in `examiner-quickstart.md` includes it as a Medium MRA under non-strict / FAIL under `--strict`, with the right framing: control-completeness for SR 11-7 reproducibility, NOT chain-integrity. The IR-scenario column says "(anomaly only — no IR; MRM-program finding under audit-procedures P-25)" — and that classification is defensible. A missing `gen_ai.request.model` or `gen_ai.response.model` is not a tampering signal, and the audit-procedures P-25 stratification is the right institutional response path. It belongs in the MRM program's audit cycle, not in the chain-integrity IR playbook.

The 12a row in `finding-language.md` reads cleanly:

> Medium MRA under non-strict; FAIL under `--strict`. Control-completeness finding for SR 11-7 reproducibility (the affected chain entries cannot be re-run to evaluate AI nondeterminism without a recorded model identifier). NOT a chain-integrity finding.

The "NOT a chain-integrity finding" emphasis matters. SR 11-7 is the model-risk-management framework an FFIEC IT examiner cross-references but does not directly own (model risk lives with the bank's MRM function and the MRM examiner, not the IT examiner). The clarification — this is control-completeness, not integrity — is what tells me whether to lift the finding into the IT examination report or hand it off to the MRM examiner. The spec §7 step 12a normative text confirms it: the check fires inline during the per-event walk, the verifier reports `gen_ai_model_identifier_missing at seq N`, and the institution's response goes into the MRM program's audit cycle.

The Module 5 row in `examiner-training.md` is well calibrated to the rest of the table:

> Medium under non-strict; FAIL under `--strict`. Investigate the institution's SDK-instrumentation for AI model calls; MRM-program finding (audit-procedures P-25).

That's the right action — the IT examiner notes the finding, points the institution at its SDK-instrumentation and at audit-procedures P-25, and hands the substantive MRM follow-up to the MRM-program function.

The bundle-sufficiency paragraph at the end of the failure-mode table closes a question I had on first read of step 12a:

> The bundle is the working-paper artifact; downstream re-verification is a property of the bundle alone, not of the institution's continued cooperation.

That is the property I needed to confirm for working-paper retention. A receiving examiner runs `bundleverify` on their own laptop, validates the verifier-binary signature, re-runs against the bundled ledger snapshot, and reproduces the chain verification independently. The bundle contents per `07-verifier-design.md` §5.3.4 (report.pdf, report.json, verifier.sha256, ledger.sha256, public_key.pem, metadata.json) are the right minimal set.

### The cadence-mismatch repeat threshold is calibrated to the regulatory cycle

This one I want to call out specifically because it's the kind of detail that's usually wrong when I read a regulator-pack draft.

`finding-language.md` step 12 (cadence) row says:

> Severity is Observation on first occurrence within a regulator standard examination cycle; MRA on repeat across consecutive examination cycles (the regulator's standard cycle is typically 12-18 months for IT examinations; the repeat threshold is two consecutive cycles with the same cadence-mismatch root cause).

That's the right frame. The OCC's IT examination cycle for community banks is 12-18 months (closer to 18 for healthy small banks; closer to 12 for higher-risk institutions); for mid-size regional banks it's typically annual. Calibrating "repeat" to two consecutive cycles with the same root cause means an institution that drifts on cadence between examinations gets a chance to reconcile its control description before the finding escalates, but a persistent cadence-vs-documentation mismatch escalates to MRA without the examiner having to invent the threshold ad-hoc.

The "with the same root cause" qualifier is the load-bearing piece. An institution might have cadence drift for two cycles for two different reasons — first cycle the documentation drifted, remediated to update the documentation; second cycle the operations drifted on a different control. Treating those as one repeat would be unfair. The threshold language correctly says "same cadence-mismatch root cause," which keeps the repeat threshold tied to the underlying control gap rather than to the surface symptom.

### The CAT crosswalk scope note is right

`finding-language.md` includes a scope note that this document is FFIEC IT Handbook–scoped, not CAT-scoped, and that a formal chain → CAT crosswalk is a v1.x candidate. That's the right disposition. The CAT (Cybersecurity Assessment Tool) is a separate FFIEC product with its own Inherent Risk Profile + Cybersecurity Maturity assessment domains and its own finding-language conventions. Trying to crosswalk chain output to CAT in a v1.0 examiner pack would either over-promise (claim CAT coverage the chain doesn't have) or under-deliver (a partial crosswalk that examiners would have to mentally complete). The honest disposition is what's there: "examiners working CAT-track examinations follow the CAT's own finding-language conventions for CAT-specific assessment domains," with a pointer back to handbook-mapping.md for the IT Handbook primary mapping and a note that a formal crosswalk is a future deliverable.

This framing also matches how the examiner cadre is organized inside my agency. IT examiners and CAT examiners are distinct skills, often distinct people, and they read distinct finding-language conventions. A v1.0 chain pack that explicitly says "we cover IT examinations, here's the IT-side language; CAT examiners adapt as needed pending a formal crosswalk" is the right disposition for both cadres.

### The handbook mapping is the right shape

`handbook-mapping.md` correctly identifies II.C.10 (Logging) as the headline. That's where I'd put the chain in any IT examination of mine; it's the integrity-of-logging requirement the IS booklet mandates and the chain extends to the property "logs are not just present, they are integrity-bearing — verifiable by an independent examiner without trusting the institution's logging infrastructure."

The II.C.13 (Cryptographic controls) deepening — per-entry `key_fingerprint`, IKM minimum 32 bytes, software-key adapter compile-time exclusion — is the substantive examination evidence beyond the generic "HMAC + Ed25519 in HSM custody" line that most cryptographic-control reviews stop at. The four examiner-confirmation items at the end of the II.C.13 commentary are concrete:

1. The institution's IKM length is at least 32 bytes (spec §10.6).
2. The institution's chain entries carry valid `key_fingerprint` bytes that the verifier confirms at step 8.
3. The institution's production build does NOT include the software-key adapter (no chain entries with `kms_handle_uri = "plaintext-dev"` in production).
4. The institution operates weekly key-fingerprint reconciliation per spec §10.1.

Those are the four checks I'd run on a II.C.13 review of an institution operating this chain. The fact that they're written down and tied to specific spec sections means a junior examiner can run them without my supervision.

The II.E (Change management) section's algorithm-change-management text covers the dual-algorithm transitional period — change management on the `algorithm` field, the `signatures` list, the institution's declared algorithm-posture configuration, IR program response procedures for spec §7 step 11 dual-algorithm failure modes, regulator coordination on multi-year migration timelines. That's the substantive change-management examination text for a multi-year algorithm transition, and it gives me concrete things to check.

### Sample report is the artifact I'd want examiners to model on

`sample-report.md` is well constructed:

- The cover page summary tells me everything in the first 30 seconds: 30 days verified, 30 passed, 0 failed, 47 late-binding events (0.001%), 3 days with sealing delays > 1h. OVERALL RESULT: PASS WITH ANOMALIES. That's the headline. The examiner-signature line is right where it should be.
- The per-day detail with `Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12` per day is concrete evidence that the verifier exercised every defense-in-depth step. That row is what I'd cite if asked at exit conference whether the verifier actually performed the integrity checks I'm relying on.
- The 2026-04-15 master-key-rotation-day example with both `[3, 4]` key versions present and both fingerprints recorded is exactly the rotation-day pattern an examiner needs to recognize. The `Anomalies: master_key_rotation_observed (rotation completed 03:14 UTC per institution incident log; both key generations correctly resolved at verifier step 7 IKM lookup)` text is the working-paper-quality finding I'd lift into my report verbatim.
- The Sample failed-day output (2026-04-23, step 8 `key_fingerprint mismatch`) is the paired example I needed to see. It demonstrates the verifier behavior on a day that fails — the partial step list (`Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)`), the FAILURE RECORD section with `expected_fingerprint` and `recorded_fingerprint`, and the cover page summary update with FAILED DAYS callout pointing to `regulator-pack/examination-response-workflow.md`. The visual structure makes failed days unmissable on cover-page review.
- The SOC team appendix (verifier output line → TSC criterion mapping) is a thoughtful inclusion. I am not the SOC reviewer, but the cross-pack consistency — the same verifier output feeding both the FFIEC IT examination and the SOC engagement — is exactly the property that lets us share evidence with the bank's SOC auditor without re-running anything. The TSC mapping table also tells me, indirectly, what control objectives the SOC engagement is evaluating, which helps me coordinate my IT examination with the bank's SOC reporting cycle.

### Deployment package addresses what regulator IT actually cares about

`deployment-package.md` answers the right questions for a regulator's IT shop:

- Static binary, no dynamic linking, ~12 MB, no network behavior, no telemetry, no privileged operations, runs as ordinary user. That's the security-architecture-review checkboxes for an examiner-laptop allowlist.
- Both allowlisting postures (SHA-256 hash allowlist for high-sensitivity engagements; cosign signature allowlist for routine examinations) are accurate descriptions of how regulator IT shops actually operate. "Many regulators run both" matches what I see in our shop.
- The trust artifacts per release (binaries, cosign signatures, GPG-signed hash manifest, CycloneDX SBOM, vulnerability scan report, source tarball, validate wrapper scripts) are what a regulator's software-assurance team needs.
- The 30-day pre-announcement window for new versions and the "existing examinations using older versions continue without disruption; a passing report from an older verifier remains a valid working-paper artifact" property are the operational guarantees we need to commit to a multi-cycle examination program. We do not want to discover mid-examination that a verifier we ran last week is no longer valid.

### Portfolio comparison procedures fit my EIC workflow

`portfolio-comparison-procedures.md` is right-sized for where the project is today. The honest acknowledgment that v1.1 `verifier portfolio` is the planned tooling, but the manual procedure suffices for portfolios up to ~20–50 institutions, matches the reality of how supervisory teams aggregate findings during the v1.0 deployment window. The sample comparison workspace (Bank A through Bank F) shows the table shape an EIC would build in a portfolio-management spreadsheet, with columns for pass days, anomaly days, sealing >24h, late-binding rate, and notes. That's the table I'd expect an EIC to maintain in their portfolio-management system.

The cross-portfolio patterns with their typical implications are directly useful:
- "All institutions in a sub-portfolio show similar anomaly pattern" → shared vendor / HSM / cloud-region issue → supervisory function investigates the shared root cause.
- "One institution differs materially from the portfolio average" → misalignment between claimed and observed cadence, operational immaturity, or unique configuration → deeper examination at next cycle.
- "Portfolio-wide degradation over time" → emerging operational issue, regulatory framework change, or vendor-side degradation → cross-agency coordination if appropriate.

These are the analytical patterns an EIC actually thinks about. Having them written down means the supervisory team's portfolio-comparison work is consistent across EICs and across cycles.

The cross-agency coordination paragraph correctly notes that for institutions supervised by multiple agencies (Fed + OCC for BHC + bank-subsidiary structure), the verifier output is the same across agencies and can be shared through standard supervisory channels. That's important; it means the chain produces a portable evidence artifact even where the supervisory framework is fragmented.

---

## Two narrow observations (neither blocks adoption)

These are file-and-forget. I'd accept the package as-is.

### Observation 1 — Working-paper hash recording

`examiner-quickstart.md` step 2 says `sha256sum ./ledger-snapshot.dump` and notes "record for working paper." `examiner-training.md` Module 3 step 2 and 3 also says to record snapshot and public-key hashes for the working paper. `sample-report.md` shows the snapshot SHA-256 on the cover page (line 61: `SHA-256: 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`).

What's not explicit anywhere I read: where the examiner records *the public-key fingerprint* in the working paper as a separate artifact from the snapshot hash. The cover page shows `Fingerprint: SHA256:HQ1cLP3OwY/...UZNiMjA=` (line 64), but the quickstart step 2 only records the snapshot. An examiner following the quickstart literally might not record the public-key fingerprint as a separate working-paper line item, which means the working paper records what was verified but not which key it was verified against.

Practical impact: the bundle (`bundle.tar.gz`) contains `public_key.pem` per `07-verifier-design.md` §5.3.4, so the receiving examiner can re-verify against the same key. The working-paper-only audit trail would not be incomplete in practice, because the bundle preserves the key. The observation is a documentation tightening, not a functional gap.

### Observation 2 — `verifier-validate.sh` cadence in the working paper

`deployment-package.md` describes the validate-before-run script and `examiner-training.md` Module 3 step 1 calls `./verifier-validate.sh ./verifier`. `portfolio-comparison-procedures.md` lists "Verifier-validation cadence" as a comparison metric ("How often the institution validated the verifier binary"). 

Two threads converge here that I want to read together: (a) the examiner runs `verifier-validate.sh` "every time, before you run it" per the quickstart; (b) the EIC reviews verifier-validation cadence as a portfolio-comparison metric.

What's not explicit: the working-paper convention for *recording* the validate-script output. Each examination presumably should preserve the `verifier-validate.sh` PASS line ("verifier-validate: OK to run") as a working-paper artifact, so the EIC's portfolio-comparison cadence metric has something to read. The quickstart and training docs run the script but don't say to capture its output.

Practical impact: an EIC computing the verifier-validation-cadence metric across the portfolio either (a) infers it from the examination metadata (examination date implies a recent validate-script run) or (b) asks each examiner to attest the run. The second is plausible and doesn't require the documentation change; the first is an inference. Either way, the metric is computable; the documentation could tighten by saying "capture the validate-script output as a working-paper line."

Both observations are tightening opportunities, not gaps. I'd ship the package as-is and address these on the next documentation pass.

---

## Final disposition

**Stopping criterion: 0 blockers / 0 should-fix.**

- Step 11 (signature verification) is examined cleanly under both single-algorithm and dual-algorithm postures; the case (e) bracket-bridging language is the load-bearing piece and it's there.
- Step 12a (`gen_ai_model_identifier_missing`) is positioned correctly as control-completeness for SR 11-7 reproducibility, not chain-integrity, and the IR/MRM-program routing (audit-procedures P-25) is defensible.
- The five-step path in `examination-response-workflow.md` is the right backbone, and the two worked examples (step 8, step 10) generalize it credibly.
- Bundle sufficiency for downstream re-verification is the property I needed for working-paper retention; it's stated explicitly and the bundle contents support it.
- The cadence-mismatch repeat threshold (Observation on first occurrence; MRA on two consecutive cycles with the same root cause) is calibrated to the actual regulatory cycle.
- The CAT crosswalk is correctly scoped out of v1.0 with a v1.x candidate flag and a pointer back to handbook-mapping.md.
- The handbook mapping puts the chain at II.C.10 (Logging) headline + II.C.13 (Cryptographic controls) deepening + II.E (Change management) for algorithm-posture transitions, which is where it belongs.
- The deployment package addresses the regulator IT shop's allowlisting and trust-artifact requirements correctly.
- The portfolio-comparison procedures fit the EIC's actual cross-bank workflow.

I would take this package to a field examination tomorrow without modification.

— Elizabeth O'Donnell, Senior IT Examiner, OCC Boston
