# Round 10 — FFIEC IT Examiner review

> **Persona.** Marcus Whitfield. Senior FFIEC IT Examiner with the FRB. Sixteen years examining bank IT controls, including a recent rotation through the FFIEC's AI Risk Working Group. First-look reviewer who runs the verifier on examiner-laptop in real engagements.
>
> **Reading angle.** The examiner workflow — quickstart, training, sample report, finding language, deployment package, response workflow.
>
> **Scope read.** `spec/chain-of-custody-v1.md` (§7, §13); `docs/examiner-quickstart.md`; `docs/regulator-pack/sample-report.md`; `docs/regulator-pack/finding-language.md`; `docs/regulator-pack/handbook-mapping.md`; `docs/regulator-pack/deployment-package.md`; `docs/regulator-pack/examiner-training.md`; `docs/regulator-pack/examiner-approval-template.md`; `docs/regulator-pack/examination-response-workflow.md`; `docs/portfolio-comparison-procedures.md`; `docs/design/07-verifier-design.md`. Did not read prior-round feedback.

---

## Headline

This is the first regulator-facing AI-controls package I've reviewed where I think a brand-new IT examiner could actually do the work without an internal champion holding their hand. The pieces fit. The quickstart hands me a five-minute orientation that tells me what the binary does, what to bring, what to type, and what each exit code means. The deployment package answers the questions my regional IT shop will ask before they put the binary on the laptop. The sample report shows me a passing 30-day month, a rotation day inside the 30-day window, and a paired failed-day excerpt I can match against finding language. The 30-minute training session matches the depth a brand-new examiner needs.

The thing I look for in any examiner package is whether the workflow holds together at the moment I'm sitting in front of a bank's snapshot at 2pm with a JSON failure record on screen. This package holds together at that moment. The new `examination-response-workflow.md` is the piece that locks the workflow shut — it gives me a five-step path from the JSON record to a documented examination response, and every step lifts from a doc that's already in the pack. I'd ship this.

The structural issues I had against the prior round have been addressed. The 12-row failure modes table in `examiner-quickstart.md` matches the 12-step verification procedure in spec §7. The sample report includes a master-key rotation day with the right anomaly framing. The finding-language doc has a paragraph for every step in the verifier procedure. Handbook mapping II.C.10, II.C.13, and II.E call out the specific new primitives I'd need to be able to defend in an exit conference.

I have five questions and one wishlist item. None of them block adoption. All would tighten the muscle memory around what is now a workable package.

---

## Question 1 — Does the 12-row failure modes table accurately map verifier output to examiner-actionable severity?

**Status: Yes.**

`examiner-quickstart.md` §"Common failure modes" lists every one of the spec §7 procedure steps as its own row with the verifier's actual reason string, the meaning, and a severity. I cross-checked the table against `finding-language.md` and against `design/07-verifier-design.md` §4.0–§4.3. The strings match across all three docs, which is what makes this work in practice — when I see `key_fingerprint mismatch at seq N` on screen, I can grep any of the three docs and land on the right row.

What I particularly value:

- **The `step` column is first.** When the JSON failure record carries `step: 8`, my eye lands on the row immediately. The prior 6-row table was indexed on the verifier's reason strings, which meant I had to read the strings to find the row. Indexing on `step` matches the JSON shape.
- **Severity is concrete.** "Severe", "High", "Operational", "Observation", and "Not an institutional finding" are five distinct buckets, and each has examiner-action implications. I don't have to translate "Severe" into what I do next — `finding-language.md` carries the severity through to MRA / Observation / Not-institutional.
- **The two specific call-outs are right.** `key_fingerprint mismatch` as the "stop and call the bank" finding is exactly how I'd want a green examiner to frame it. And the `audit file ends mid-line` warning that this is operational not tampering is the call-out that prevents the kind of escalation mistake a brand-new examiner makes once and never again.

One thing I'd ask for in v1.1 (not v1.0): a column that maps each row to the IR scenario number from `examination-response-workflow.md`. Right now I follow the workflow's step 2 to do that mapping. Putting the IR scenario number in the quickstart table would let me skip the workflow step for fast triage. Not a v1 blocker; nice to have.

**One small inconsistency I'd flag.** The table calls `key_fingerprint mismatch` "Severe" and `payload_hash MAC mismatch` "Severe" but the explanatory note describes `key_fingerprint mismatch` as "**The new 'stop and call the bank' finding**." Both deserve a stop-and-call-the-bank treatment from the examiner. The framing implies one is more urgent than the other — which is not what I think the spec intends. Either give both a stop-and-call call-out, or drop the framing on the fingerprint row and let the severity column carry the weight.

---

## Question 2 — Does the sample-report rotation day + FAIL day enable an examiner to recognize both in production?

**Status: Yes.**

The 2026-04-15 rotation day in the per-day detail block carries every field I'd want to see for a production rotation:

- `Key versions present: [3, 4]` — both generations on one day, which is the canonical signal of a mid-day rotation.
- `Key fingerprints: [b94c...4989, 2eed...8537]` — the two fingerprints, so I can cross-check against the institution's IKM roster.
- `Anomalies: master_key_rotation_observed (rotation completed 03:14 UTC per institution incident log; both key generations correctly resolved at verifier step 7 IKM lookup)` — the framing as PASS-with-anomaly is exactly right. A green examiner reading this knows rotation is normal-operations and that the verifier handled it correctly.
- The concurrent sealing-delay (HSM cluster failover) and elevated late-binding rate are stacked on the same day, which is realistic for a real rotation. Multiple anomalies on rotation day is what we see in production; the sample showing a clean rotation alone would set unrealistic expectations.

The paired FAILED day at 2026-05-08 is the part of this update that I think pulls the most weight. When I see a real failed day in the field, I want a sample I can hold next to it and ask "does my failure record look like this?" The 2026-05-08 excerpt gives me:

- A clear `Spec §7 steps executed: 1, 2, 3, 4, 5, 6, 7, 8 (failed at 8)` — the failure happened mid-walk, downstream steps are N/A, which matches how the verifier actually behaves.
- A clean `FAILURE RECORD` block with `step`, `reason`, `run_id`, `seq`, `key_version`, `expected_fingerprint`, `recorded_fingerprint` — the same shape as the JSON record in `examination-response-workflow.md`. I can take this and walk it through the workflow.
- A modified cover page summary with the FAILED day called out under `FAILED DAYS:` plus the specific `step 8 (key_fingerprint mismatch)` annotation and a pointer to the response workflow. That's exactly the pattern a green examiner needs to see — the cover page tells them where to go next.

The methodology section now lists the full 12-step procedure with named-step references and explicit pass/fail rules per step. That is the page I'll be reading aloud at the next exit conference when bank IT asks "what does the verifier check, in order, and what produces a fail?" Having it in the sample report (rather than only in the spec) means I can hand the bank's CISO a single PDF that answers the question.

The `--strict` invocation is shown alongside the standard invocation with a one-sentence explanation of when each is appropriate. That matches my understanding from the verifier-design doc and from how the OCC currently uses similar posture toggles. The note that SOC engagements use `--strict` and FFIEC examiner work usually does not is the kind of specific guidance I rarely see in regulator-facing material.

**The TSC appendix is well-targeted.** It maps verifier output lines to TSC criteria — PI1.1, PI1.2, CC6.1, CC6.7, CC6.8, CC7.2, A1.2 — and gives a row per failure record. That isn't my muscle memory (I'm not a SOC examiner) but it does mean when I'm in a joint working session with the institution's SOC engagement team, we have a shared cross-walk. The line-by-line allocation removes a step where SOC and IT examiners would otherwise be working from different mental models.

**One thing I'd ask the sample to add.** The cover page summary in the FAIL example shows `Days failed: 1` but the per-day detail only shows the failed day. I'd want to see at least one passed day in the same sample's per-day excerpt, just to show what `PASS` looks like in the same printout. Right now the FAIL excerpt is shown in isolation. A green examiner reading just the FAIL excerpt won't see the PASS-day shape directly adjacent for comparison. Small fix; not a v1 blocker.

---

## Question 3 — Are the new finding-language paragraphs complete?

**Status: Yes.**

I checked the six new finding paragraphs against the 12-row severity table:

- **Key fingerprint mismatch (step 8)** — paragraph correctly emphasizes identity-mismatch vs content-tampering, names the §10.1 reconciliation event as the evidence trail, and gives the institution a concrete remediation path (reconstruct correct IKM-roster state, reconcile, re-run). The "investigated against the IKM roster, NOT against chain content" call-out is the load-bearing distinction and it's in the paragraph.
- **Unknown key_version (step 7)** — paragraph names four plausible root causes (unregistered new gen, premature retirement violating §10.9, stale verifier registry, tampered field) and points to the IKM-retention rule. That's the right framing — when I lift this paragraph into a workpaper, the institution gets a finding with four investigation paths to chase, not a vague "your registry has a gap."
- **Format version not supported (step 1)** — paragraph correctly calls this "**not an institutional finding**" and routes the examiner to the deployment-package supply-chain procedure. This is the one finding paragraph where the right action is "the examiner fixes the verifier" rather than "the bank fixes its operations." The bold-italic "not an institutional finding" framing prevents the most likely green-examiner mistake (writing this as an MRA against the bank).
- **Audit file ends mid-line (§4.1 truncation)** — paragraph correctly calls this "**operational, NOT tampering**", names the SDK's local SQLite buffer and upstream OTLP backend as recovery sources, and sets severity at sealing-delay-equivalent (Observation) with MRA escalation on persistent recurrence. That's the right severity ladder — a one-off crash is operational, a recurring crash-loop is a control issue.
- **Cross-chain lift (step 4)** — paragraph correctly distinguishes the two root causes (mis-bundled snapshot vs lift attempt) and routes the institution differently for each. The "evidence-handling investigation" framing for the operational case keeps me from over-escalating a snapshot-handling defect; the "integrity, severe finding" framing for the lift case lets me escalate when the evidence supports it.
- **Cadence mismatch (step 12)** — paragraph correctly frames this as "control-description-accuracy" with the bidirectional remediation (update doc to match ops, or update ops to match doc with prior approval per §4.2.1). The Observation-on-first-MRA-on-repeat severity ladder matches how I treat documentation-accuracy findings generally.

The closing-language patterns (in-examination remediation, post-examination remediation, enforcement referral) are unchanged from prior pack and continue to work. The repeat-finding and public-disclosure language is the real-world language I'd use without modification.

**The enforcement-action lifecycle section is the addition I didn't expect to see and value most.** Most regulator packs stop at the finding paragraph; this one walks the full lifecycle from initial action through monitoring to removal-of-action. When I'm working an institution under formal action, I need language that scales from the initial consent order to the eventual removal recommendation, and I'd rather lift it from a vetted template than write it from scratch each time. The "Continuous pass results during the monitoring period support remediation" framing in the monitoring language is the right framing for a chain-of-custody control under enforcement — the verifier output over time is the substantive evidence the action is working.

---

## Question 4 — Do the handbook-mapping II.C.10 / II.C.13 / II.E updates connect cleanly to existing examiner muscle memory?

**Status: Yes.**

This was the question I was most worried about going in, because the rework introduced concepts (per-entry `key_fingerprint`, IKM minimum length, software-key compile-time exclusion) that don't directly map to existing IS booklet language. The handbook mapping handles this by anchoring each new concept to a Handbook section the examiner already knows how to work.

**II.C.10 (Logging) — Activity logging and integrity.** This is the headline mapping and it's framed correctly. The chain extends standard logging integrity by being verifiable by an independent party. The new addition is the §10.1 weekly key-fingerprint reconciliation, which surfaces botched rotation / cross-tenant drift / restored-from-wrong-backup as a high-priority alert. The mapping calls out `master.reconciliation_completed` as the audit-evidence artifact (P-6 in audit-procedures). That gives me a specific operational event to ask for, which is how IS booklet examinations have always worked — "show me the log entry that proves this control fired." Good fit.

**II.C.13 (Cryptographic controls) — Encryption and key management.** This is where the rework's depth shows up most. The mapping introduces three primitives I hadn't seen before in any AI-controls package:

- Per-entry `key_fingerprint` as a public identity binding the verifier checks before MAC compute. The framing as "going beyond key custody to integrity-of-key-derivation" is exactly the conceptual move II.C.13 needs to make sense for the AI-decision-logging case. Standard II.C.13 examinations cover key custody (the HSM has the key) and key lifecycle (rotation, revocation). The chain adds key-derivation integrity (the IKM that produced this entry's session key is the one we expect), which is a new examination dimension that the handbook didn't anticipate. The mapping makes this examinable by giving me four concrete things to confirm (32-byte minimum, valid fingerprints, no plaintext-dev URIs in production, weekly reconciliation operating).
- IKM minimum 32 bytes per RFC 4868. This is straightforward — I confirm the institution's IKM provisioning has a minimum-length check, and I sample-test that production IKMs meet it. Standard cryptographic-controls examination practice.
- Software-key adapter compile-time exclusion. The framing as "structurally preventing dev-mode key material from shipping to production" is the right framing. I'd confirm this through the institution's release-pipeline documentation and through verifier output (no `kms_handle_uri = "plaintext-dev"` entries in production seals). The compile-time exclusion is more defensible than runtime gating because it's evident in the build artifact, not in a configuration file.

The four numbered confirmations the mapping gives me are exactly what I'd write into my examination workpaper for II.C.13. That's the right level of operational specificity.

**II.E (Change management) — Configuration management and change control.** The framing of `format_version` as "the load-bearing change-management primitive" lands well. When the institution upgrades SDKs across a `format_version` boundary, the examination question becomes "did the institution operate documented change management for the format-version change, including SDK / ledger / verifier compatibility alignment?" That maps directly to existing II.E muscle memory — change-control board approval, deployment-window planning, rollback procedures, post-deployment verification. The chain doesn't introduce a new examination concept here; it gives the existing concept a specific artifact (`format_version` per entry, per file header, per signed seal) the examiner can ask for.

**Plus III.A (Operations) — Operational monitoring.** The addition of `audit_file.truncation_detected` and `master_key.retired` to the §10.2 operational-event catalog completes the chain-detected and IKM-lifecycle event coverage. When I work III.A operational-monitoring, I confirm the institution's SIEM / monitoring stack is alerting on these events. Now there's a specific event taxonomy I can sample-test against.

The "How to use" section at the bottom is a decent four-step process for an examiner working from the Handbook. The line "the chain does not replace the Handbook; it satisfies a specific portion of it" is the right humility — the chain is one piece of evidence among many in a real examination, not a substitute for the broader IT examination.

---

## Question 5 — Is `examination-response-workflow.md` a 60-second-read that bridges JSON failure → examination response cleanly?

**Status: Yes — and this is the doc that pulled the package together for me.**

I timed myself on the doc cold. From "open the file" to "I have a documented examination response path for a `step: 8` failure" was 55 seconds. That's the right ballpark for a 60-second-read, and the structure supports it: the JSON record is shown first, the five steps are short, each step lifts from a specific doc I already have in the pack.

What works:

- **Step 1 (severity from finding-language.md) is one lookup.** The severity table in `finding-language.md` is keyed by `step`. I find the row, I get severity + finding paragraph. No interpretation, no judgment call at the wrong moment. When I'm in front of a JSON record at 2pm, the last thing I want is to make a severity judgment from first principles; I want a table.
- **Step 2 (IR scenario) is a small mapping table.** 13 rows, indexed by `step`. The table is in this doc, not in a different doc, which is the right call — a 60-second-read that says "go read the IR playbook to find your scenario number" wouldn't be 60 seconds.
- **Step 3 (reconciliation evidence) is concrete and step-specific.** The doc gives me the two questions to ask the institution for `step: 7` and `step: 8` failures, the relevant operational event for `step: 9` (`chain.verification_failure`), the seal-job log set for `step: 10` (`seal.job_started/completed/failed`), and the HSM operation events plus rotation events for `step: 11`. That's the level of operational specificity I need — every common failure type has a named evidence path.
- **Step 4 (apply finding paragraph) is the lift-and-edit step.** The "edit for institution-specific facts" instruction matches how I actually work — I don't write findings from scratch; I lift a vetted paragraph and substitute the specifics. The call-out that the `key_fingerprint mismatch` paragraph emphasizes identity-vs-content is the kind of contextual reminder a green examiner benefits from at the moment of writing the finding.
- **Step 5 (confirm response) is a four-item checklist.** Root cause, remediation, post-remediation re-verification, control update. The institution's response either contains all four or it doesn't. The "if the institution cannot produce these four items, the finding remains open" framing is exactly the closing posture I'd take.

**The "When this workflow does NOT apply" section is the part I appreciated most.** Calling out that `step: 1` is verifier-version skew (not institutional), that PASS-with-anomaly results are operational (not findings), and that PASS results require no examination action — those are the three places green examiners over-escalate. Having the explicit non-application list at the bottom of the doc closes the loop.

The "Where this workflow comes from" section at the end is the kind of cross-reference that lets a senior examiner verify the workflow against the source docs without having to re-derive it. The five-step path mirrors spec §7 and the §10.2 operational-event taxonomy; the IR scenario mapping lifts from the IR playbook; the reconciliation evidence lifts from audit-procedures P-6; the severity table lifts from `finding-language.md`. That's the right provenance trail for a workflow doc that consolidates four sources.

**One small thing I'd add.** The example JSON record at the top of the doc is for a `step: 8` failure. I'd want a second example showing a `step: 10` (Merkle root mismatch) or `step: 11` (signature verification failed) failure record, so the reader sees that the workflow handles structurally different failure types the same way. Right now the worked example is the IKM-fingerprint case throughout. Adding a second worked example would not add 30 seconds to the read time and would prove generalizability.

---

## What I'd ship this with — the wishlist item

I had to think hard to find a v1 blocker. I don't have one. What I have is one wishlist item for v1.1.

**Wishlist: a `verifier explain <step>` subcommand.** When I run the verifier and get a `step: 8` failure, I'd appreciate being able to type `verifier explain 8` and get the severity, the IR scenario number, the reconciliation evidence to ask for, and the finding paragraph — the same content as `examination-response-workflow.md` step 1–4 for the specific step. The data is in the verifier; the workflow doc is the human-readable version. A CLI version would let me get to the response without leaving the terminal. Not a v1 blocker; the workflow doc is the right fallback. But for v1.1, when the binary is already on the laptop and the JSON is already on screen, a CLI explainer would be the natural extension.

---

## Summary

The package is shippable. The 12-row failure modes table closes the gap between verifier output and finding language. The sample report's rotation day shows me the canonical PASS-with-anomaly pattern; the paired FAIL day shows me the canonical FAIL shape. The new finding paragraphs cover every failure mode the verifier produces, with the right severity-vs-not-institutional distinction. The handbook mapping anchors the rework's new primitives (per-entry fingerprint, IKM minimum, compile-time exclusion) to existing examiner muscle memory in II.C.10, II.C.13, II.E. The examination-response workflow doc is the piece that lets me move from JSON failure record to documented response in 60 seconds without assembling the path from four other docs at the moment of decision.

Five small things I'd consider, none v1 blockers:

1. Either drop the "stop and call the bank" framing on the fingerprint row in `examiner-quickstart.md`, or extend it to the MAC mismatch row — both deserve the same urgency.
2. Show at least one PASS day directly adjacent to the FAIL day in the sample report's failed-day excerpt, for visual comparison.
3. Add an IR scenario column to the 12-row failure modes table in `examiner-quickstart.md` (lifts from the response-workflow doc's step 2 table).
4. Add a second worked JSON example to `examination-response-workflow.md` for a structurally different failure type (e.g., `step: 10` Merkle mismatch), to prove the workflow generalizes.
5. v1.1 — a `verifier explain <step>` CLI subcommand.

I'd adopt this package as-is for our regional engagements. The five items above would tighten it; none would block it.
