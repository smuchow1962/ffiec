# Round 12 — Big Four MRM review

**Reviewer.** Dr. Sandeep Kumar Reddy. Big Four MRM Senior Director. Twenty-two years on SR 11-7 model-validation engagements. Lead author of the firm's AI-MRM playbook for South Asian bank holding companies (also frequently consulted on US BHC engagements where SR 11-7 governs).

**Reading angle.** Model-risk substance. SR 11-7 reproducibility. MRM committee oversight workability. ECOA / customer-dispute chaining. AI-policy alignment beyond NIST AI RMF.

**Scope read.** Spec §4.4 (gen_ai_parameters) and §5 (linkage fields included in canonical bytes); MRM-COMMITTEE-BRIEF; customer-dispute-procedures; regulator-pack/ai-policy-alignment; audit-procedures P-22..P-25; design 02 §8/§9 (chain-construction); design 07 §4.3.1 (verifier DAG topology step 13); anomaly-documentation-template (worked example); management-summary; audit-committee-summary.

**Stopping criterion declared up front.** 0/0 — zero blockers, zero majors required to ship from an MRM-substance lens. Findings below either confirm 0/0 or document anything that would push above it.

---

## Headline

The corpus is the first AI logging-integrity package I have read where the MRM substance is genuinely close to what an SR 11-7-shaped validation engagement actually consumes. Most products in this category give me a tamper-evident log and assert it "supports SR 11-7." That assertion is usually correct in the narrow sense and useless in the substantive sense — the validator has nothing to reproduce against and the MRM committee has nothing to oversee beyond the existence of the log.

This package is structured differently. The chain integrity-binds the inputs an SR 11-7 reproduction needs (`gen_ai_parameters`, request-side AND response-side model identifiers, `audit.*` payload including system prompt and retrieval context). The MRM committee brief tells the chair what to ask. The customer-dispute procedure addresses the reproducibility question (would the AI have said the same thing on a different day) honestly rather than waving at it. The audit procedures (P-22..P-25) test the evidence completeness for the failure modes the chain construction actually produces, not generic IT-control wording.

The result is a logging-integrity control I would accept as one of the integrity-bearing components of an institution's broader AI-MRM framework, with the caveat MRM teams always carry: integrity is necessary, not sufficient. The chain-of-custody framework is correct in calling that out repeatedly (`MRM-COMMITTEE-BRIEF.md`: "The chain produces integrity, not validity"; "The MRM Committee's qualitative judgment is unchanged. The evidence supporting that judgment is now stronger"). That framing is the right framing; many vendors over-claim and the over-claim is what gets the institution into trouble downstream.

**Headline verdict.** 0/0. Zero blockers, zero majors. The findings below are either MINOR (worth doing, would not block) or OBSERVATIONS (worth recording, no action required).

---

## What's load-bearing for SR 11-7 substance, and whether the corpus delivers

SR 11-7 §V (model validation) requires three things from the records the validator consumes:

1. **Reproducibility.** The validator can reconstruct the model's input-output pair under known conditions.
2. **Authenticity.** The validator has confidence the input-output pair is what the model actually produced, not what someone said the model produced.
3. **Lineage.** The validator can trace the input-output pair to the upstream data, parameters, and configuration that produced it.

The chain delivers (2) directly via the per-event MAC + daily Merkle seal + HSM-rooted root signature. That is the un-controversial part and I will not dwell on it.

The chain delivers (1) and (3) **substantively** via three design choices that I want to flag explicitly because they are unusual in the AI logging space:

### Choice 1 — `gen_ai_parameters` is integrity-bound but schema-flexible

Spec §4.4 makes `gen_ai_parameters` an OPTIONAL JCS-canonical JSON field whose schema the institution determines. The chain integrity-binds whatever the institution puts there.

This is the right separation. The chain doesn't dictate what reproducibility surface the institution captures; it provides the integrity property over whatever the institution chooses to capture. The RECOMMENDED contents (decoding parameters, sampler implementation identifier, request-side AND response-side model IDs, system-prompt content or content-hash, retrieval context, intra-run data dependencies) are exactly the surface SR 11-7 reproducibility actually needs.

The most under-appreciated item in that RECOMMENDED list is the request-side / response-side model-ID distinction. Vendors silently re-route between model versions during outages and capacity events. An MRM team reproducing a decision under what it thinks is the original model version, when the response actually came from a fallback model, produces a reproduction that lies. The corpus catches this explicitly:

> "model identifier and version (request-side AND response-side, per the OTel `gen_ai.request.model` / `gen_ai.response.model` semconv — vendors silently re-route between model versions during outages)"

That sentence is the kind of thing that costs an MRM team a year to learn in production. Naming it in spec §4.4 is institution-friendly.

P-25 in audit-procedures completes the loop: the SOC team confirms the institution's `gen_ai_parameters` schema covers the reproducibility surface, with schema gaps explicitly classified as MRM-program findings rather than chain-integrity findings ("the chain integrity-binds whatever the institution puts in `gen_ai_parameters`; they are MRM-program findings that the institution's effective-challenge capability is weaker than the spec's reproducibility surface allows"). That classification is correct and important — without it, the auditor's instinct is to log the schema gap as a chain-integrity finding, which would make the chain look worse than it is and would conflate two separate control areas.

### Choice 2 — Cross-run linkage fields are inside the canonical bytes

Spec §5: `parent_run_id`, `parent_seq`, `dag_parents` are INCLUDED in the canonical bytes the MAC covers. They are not chain-stamp fields; they are substantive evidence about the agent's decision graph.

This is the right call and an unusual one. Most logging-integrity systems treat parent linkage as metadata and exclude it from the integrity surface. The threat that closes is: an attacker who can write to the ledger after capture but before seal cannot rewrite the parent linkage without breaking the MAC. The lineage property (3) above is therefore integrity-bearing.

For a DAG-shaped agent flow (multiple parents joining at a final-decision step), this is what makes the post-hoc reconstruction trustworthy. The MRM team reconstructing the agent's decision graph from the chain has cryptographic assurance the parent edges weren't rewritten between capture and reconstruction.

Design 07 §4.3.1 step 13 (DAG topology resolution) closes the resolution side: the verifier, after the per-day Merkle and signature checks, walks each event's `dag_parents` and asserts each parent resolves to a chain entry that exists in the same tenant's ledger. Anomalies (typo in `dag_parents`, missing parent event, parent in a different tenant) are reported in the verifier's day-level summary with a `dag_topology_check` section rather than discovered downstream during a SR 11-7 reconstruction.

The DAG step is correctly placed at the day-level, not per-event during the chain walk, because parent events may be in earlier days within the verification range. The reasoning is right and the placement is right.

### Choice 3 — ECOA adverse-action notice is itself a chain entry

`customer-dispute-procedures.md` line 42:

> "The translation [of the AI's response into ECOA-compliant adverse-action language] is itself a chain entry, with the AI's original response as the parent (per spec §4.4 `parent_run_id` / `parent_seq` linkage). This is the institution's choice but it is the recommended posture: a translation logged outside the chain is repudiable in the same way a vendor's logs are repudiable, which is exactly the gap the chain exists to close. The MRM committee opines on whether the institution chains the translation step or logs it outside the chain; the recommended decision is to chain."

This is the single most institution-protective sentence in the corpus from an MRM-defensibility perspective. The ECOA adverse-action notice translation is the moment a customer's claim ("the bank's notice does not match what the AI actually said") becomes investigable. If the translation is logged outside the chain, the customer's counsel has the same repudiation argument against the translation that the chain exists to close against the AI decision itself. By chaining the translation step with the AI's response as the parent, the institution makes the notice's reasons cryptographically traceable to the AI's actual response.

The MRM committee's role in that decision is exactly the right placement of authority: the committee opines on the institution's posture, the chain's primitives support either choice, the recommended posture is named.

This is the kind of detail a Big Four MRM engagement looks for as a marker that the institution understood the customer-facing integrity question. Most institutions don't. This corpus does.

---

## SR 11-7 reproducibility — does the customer-dispute flow actually work?

`customer-dispute-procedures.md` §"Reproducibility and the customer's 'would the AI have said the same thing on a different day' question" handles the reproducibility question with a level of honesty I rarely see in vendor documentation. The institution's response template:

> "The decision's full input set is recorded in our integrity-bearing audit trail, including the AI model's version, the decoding parameters, the system prompt, and any retrieved context. We can re-run the decision under the same inputs to evaluate whether the AI would produce the same output. Significant variation across re-runs (a phenomenon called model nondeterminism) is itself information our model-risk-management committee evaluates; consistent variation in your favor would prompt us to revisit the original decision."

Three things to call out:

1. **The acknowledgment of model nondeterminism.** Most institutions do not acknowledge this in customer-facing language. The acknowledgment puts the institution in a defensible posture — "yes, the model is nondeterministic, here is how we handle that" — rather than the brittle posture of asserting the model is deterministic when it isn't.

2. **"Consistent variation in your favor would prompt us to revisit the original decision."** This is the institution committing to act on reproduction findings. From an MRM perspective, this is the right commitment to make and it is the MRM committee's call whether to make it. Naming it in the customer-facing template is institution-protective: it tells the customer the institution will actually do something with the reproduction, and it tells the regulator the institution has thought about the reproduction-result-handling protocol.

3. **The MRM committee establishes the institution's reproduction posture.** Correct placement of authority. The chain provides the substrate; the MRM committee owns the policy.

The IKM-access section is where the corpus does the most institution-protective work I have seen on this topic:

> "1. Protective-order disclosure (litigation context). Per `docs/legal-disclosure.md` §'Court-ordered master-key disclosure,' the IKM is disclosed under a court-issued protective order...
> 2. HSM-mediated verification (operational context). The institution provides the customer's expert with an HSM-mediated verification path (the expert's verifier dispatches HMAC operations through the institution's HSM API for the affected runs without ever holding the IKM bytes)..."

Naming the two acceptable shapes — and specifying the threshold at which the institution shifts from operational to protective-order — is exactly the MRM-program-integrity work that most institutions skip. The risk it closes is real:

> "Without naming the shape, the institution risks the customer's expert encountering `unknown_key_version` at step 7 mid-verification, which leaves the verification incomplete and the dispute unresolved."

That is precisely the failure mode I have seen in litigation engagements where the institution had not pre-decided its IKM-disclosure posture. The expert hits step 7, the verification stalls, the dispute drags, and the institution's regulatory posture deteriorates because the dispute remains open.

The MRM committee's annual review of the chosen shape (line 108) is the right cadence and the right oversight placement.

---

## MRM committee brief — workability assessment

`MRM-COMMITTEE-BRIEF.md` is structured for a 2-4-page brief the committee chair reads instead of the full design docs. I scored this against what a real MRM committee chair (typically a senior risk officer or a former regulator) actually needs:

| Need | Coverage |
|---|---|
| Plain-English description of what the chain is | Yes — first two paragraphs nail it; "integrity, not validity" is the right opening framing |
| What changes for MRM after adoption | Yes — the before/after table is the right shape |
| Why the four primitives compose | Yes — including the "any one missing weakens the property" callout |
| Where each primitive lives in the institution's perimeter | Yes — explicit on per-event SDK vs. server-side seal vs. HSM-rooted signature; this is the load-bearing custody-placement reasoning |
| What questions to ask the chain owner | Yes — seven specific questions, each tied to a verifiable artifact |
| How chain composes with SR 11-7 | Yes — four explicit composition points |
| Cost summary | Yes — $80k–$200k/year all-in for mid-size, with reference to the cost model |
| What committee decisions to take | Yes — four explicit decisions named |

The seven oversight questions are workable. Question 4 in particular is the kind of question a regulator-trained committee chair asks:

> "How are master-key (IKM) rotation and key-version transitions managed? Specifically: walk us through the most recent key-fingerprint reconciliation report (the `master.reconciliation_completed` operational event), including any fingerprint mismatches during the period and how they were resolved (audit-procedures P-6)."

That question is directly answerable from the operational events catalog and forces the chain-ops team to either produce the artifact or admit they don't have it. Either answer is useful to the committee.

Question 7 catches the new failure modes from the v1.0-rework:

> "What is the institution's incident response plan for chain failures? Specifically: does the playbook cover the new failure modes — `key_fingerprint mismatch` (Scenario 7), `unknown_key_version` (Scenario 8), `audit_file.truncation_detected` (Scenario 9)?"

Naming the scenarios by number lets the committee chair cross-check against the IR playbook directly. This is the kind of specificity that distinguishes a workable committee brief from a brochure.

The "three additional defensive properties from the v1.0-rework" section (per-tenant cryptographic binding; compile-time prevention of dev-key material; constant-time comparison discipline) is the right level of detail for a committee chair. Specific enough to be checkable; not so detailed that the chair is doing cryptographic review.

---

## Audit procedures P-22..P-25 — evidence-completeness substance

The four procedures (P-22 `key_fingerprint mismatch`, P-23 `hkdf_inputs_digest mismatch`, P-24 `audit_file_truncation_detected`, P-25 `gen_ai_parameters` schema completeness) are the procedures that test whether the institution actually operated the controls the spec mandates.

**P-22** is the load-bearing one. The four required evidence pieces:

1. IKM-roster row identification — which `(tenant_id, key_version)` pair had an IKM that did not produce the recorded `key_fingerprint`
2. Change-management approval for the IKM-roster correction
3. Re-verification result on the corrected roster (PASS for the affected period, OR documented gap)
4. Reconciliation cross-check — the institution's most recent `master.reconciliation_completed` event shows the corrected triple matching the IKM roster

That is a complete evidence set. An institution that produces all four for each `key_fingerprint mismatch` finding has demonstrated the control operated end-to-end, not just at the verifier-detection moment. An institution missing one of the four has a gap that warrants escalation.

The worked example in `anomaly-documentation-template.md` lines 147-221 is the friction-reducer that makes P-22 actually operable. The first time an institution stands up the four-piece evidence package, having a worked example to copy from drops the time-to-first-record from days to hours. The example is institution-realistic (botched key rotation; HSM operator misread the procedure and overwrote `key_version=1` instead of provisioning `key_version=2`; remediation includes restoration of the original IKM bytes from KMS history, change-management approval for the restoration, verifier re-run, reconciliation cross-check). The remediation step that updates the rotation procedure to prevent recurrence is the right closing posture — operational anomaly explained, root cause addressed, control updated.

**P-23** correctly names the two acceptable root causes (known SDK / verifier constants change documented in change management; multi-region deployment-drift case). The multi-region case is the one most institutions have not thought about — an SDK in us-east-1 patched to a new `format_version` while us-west-2 is still on the old one will produce different `hkdf_inputs_digest` for the same tenant. Naming it as an acceptable explanation, with the institution's response named (re-deployment to consistent constants and re-verification), is institution-protective.

**P-24** correctly distinguishes `complete` / `partial` / `unrecoverable` recovery outcomes for `audit_file.truncation_detected`. The partial case is the one that requires the most specific evidence: the recovered events are restored, the un-recovered subset is documented as an unrecoverable gap with IR Scenario 9 disposition for those specific events. That `(run_id, seq)`-level granularity is what the SOC team needs to evaluate the gap.

**P-25** is the bridge from chain-integrity audit to MRM-program audit. The procedure samples the institution's chain entries, confirms `gen_ai.request.model` and `gen_ai.response.model` are present on chain entries representing model calls, and confirms the institution's `gen_ai_parameters` schema covers the documented reproducibility surface. The classification of schema gaps as MRM-program findings (not chain-integrity findings) is the load-bearing distinction I called out above — without it, the audit produces the wrong category of finding and the institution remediates the wrong control.

---

## AI policy alignment — the regulator-facing surface

`docs/regulator-pack/ai-policy-alignment.md` covers a wider surface than I expected. The Treasury RMF (Feb 2026) composition table is the one I would lean on most heavily in an MRM engagement defending chain adoption to OCC examiners under the existing SR 11-7 / 2011-12 framework. Four explicit compositions:

| Treasury RMF expectation | Chain composition |
|---|---|
| Model inventory | Chain captures every decision; per-tenant per-day record set supports model-inventory completeness checks via `gen_ai.request.model` |
| Model lifecycle controls | Chain output supports validation, ongoing monitoring, retirement evidence |
| Third-party model governance | Chain output is independent of vendor cooperation |
| Ongoing monitoring | Drift analysis runs on chain-bearing data with confidence the data is what the model actually produced |

The model-inventory composition is the one I had not seen articulated this way in any other vendor's documentation. The reasoning is simple and correct: every model that produced a decision is enumerated by the `gen_ai.request.model` field on the integrity-bound canonical bytes. The institution's model-inventory completeness check can therefore consume the chain's data and have cryptographic confidence the inventory is complete (within the scope of decisions the chain captured).

The OCC supervisory-letter framing (line 53-59) is honest about what the OCC has and has not done:

> "As of v1.0-rework publication, the OCC has not issued formal new AI-specific guidance; instead, it operationalises AI model-risk supervision under the existing SR 11-7 / OCC Bulletin 2011-12 model-risk-management framework."

That is exactly the right framing. Over-claiming "the OCC requires the chain" would be wrong and would damage the institution's posture if challenged. Under-claiming "the OCC has no expectations" would miss the supervisory-letter signaling. The middle posture — the chain is the technical control implementing the integrity-of-AI-decision-evidence requirement under the existing SR 11-7 / 2011-12 framework, adopted to position the institution for current and emerging OCC expectations including expectations the OCC has signalled through supervisory letters to peer institutions — is the defensible posture for an MRM director.

The EU AI Act Article 14 + 26 composition (line 75-77) is the secondary alignment that matters most outside the US. Article 14 (effective human oversight) requires the human reviewer to have sufficient information to second-guess the AI; the chain's `gen_ai_parameters`, `audit.*`, and `gen_ai.request/response.model` are the substrate that makes Article 14 oversight effective. Article 26 (deployer obligations) imposes monitoring, logging, and informing-affected-persons obligations; the chain composes into all three.

The Singapore MAS Veritas / FEAT, JFSA, and BCB callouts are correctly scoped — chain supports the accountability and audit-trail expectations; the institution's broader compliance program handles the rest. For an institution operating across APAC, the chain is jurisdiction-neutral and the institution's per-jurisdiction control descriptions adapt the framing.

---

## What I checked specifically against my SR 11-7 / AI-MRM playbook

Going through the playbook's "things that distinguish a substantive AI logging-integrity package from a brochure" checklist:

| Playbook item | Result |
|---|---|
| Does the integrity-bearing record include the model's parameters at decision time? | Yes — `gen_ai_parameters` integrity-bound; RECOMMENDED contents enumerate the reproducibility surface |
| Does it distinguish request-side from response-side model identifiers? | Yes — explicit in spec §4.4 with the rationale (vendors silently re-route during outages) |
| Does it integrity-bind the cross-run lineage? | Yes — `parent_run_id`, `parent_seq`, `dag_parents` inside canonical bytes per spec §5 |
| Does it address the customer's reproducibility question honestly? | Yes — customer-dispute-procedures §"Reproducibility" acknowledges nondeterminism and commits to act on reproduction findings |
| Does it specify the institution's IKM-disclosure posture for litigation vs. operational disputes? | Yes — two acceptable shapes named with the threshold at which the institution shifts |
| Does it chain the ECOA adverse-action notice translation? | Yes — recommended posture explicitly named; MRM committee opines on the choice |
| Does it test the `key_fingerprint mismatch` evidence completeness? | Yes — P-22 with four-piece evidence requirement |
| Does it provide a worked example of the four-piece evidence package? | Yes — anomaly-documentation-template lines 147-221 |
| Does the verifier resolve DAG topology? | Yes — design 07 §4.3.1 step 13 with day-level placement and anomaly reporting |
| Does the MRM committee brief tell the chair what to ask? | Yes — seven specific questions, each tied to a verifiable artifact |
| Does the policy-alignment doc enumerate explicit composition points beyond NIST AI RMF? | Yes — Treasury RMF (4 compositions); EU AI Act Articles 12, 14, 26; DORA Articles 5, 6, 8, 17, 28; UK FCA/PRA; MAS Veritas/FEAT; JFSA; BCB |

Twelve of twelve. I have not previously found an AI logging-integrity corpus that checked more than seven of these.

---

## Findings

### MINOR — none rising to MAJOR

**M-1. P-25 sampling guidance for `gen_ai_parameters` schema completeness could be tightened.**

P-25 says "pull a representative sample of the institution's chain entries with `gen_ai_parameters` populated" and confirms the schema covers the reproducibility surface. The general sampling table at the bottom of audit-procedures (< 25, 25–250, 250–2,500, > 2,500) applies, but P-25 is unusual in that the population the SOC team cares about is "model-call-representing chain entries," not all chain entries. In a high-volume institution, the population of model-call entries may be 1% of the chain-entry total; sampling against the chain-entry total may miss the model-call distribution.

Recommend P-25 say something like: "Stratify the sample by chain-entry kind so that model-call-representing entries are over-sampled relative to non-model-call entries; the model-call distribution is the population P-25 tests against."

This is a wording tightening, not a substantive change. The procedure works as written; the stratification guidance reduces the chance an inexperienced SOC team mis-samples.

**M-2. MRM-COMMITTEE-BRIEF question 3 could explicitly call out the model-inventory composition.**

Question 3 covers the `gen_ai_parameters` and `audit.*` schema decisions. It does not explicitly mention the model-inventory completeness composition (Treasury RMF first row in `ai-policy-alignment.md`). For a committee chair familiar with the institution's broader Treasury RMF posture, the connection is natural; for a chair newer to AI-MRM, naming the composition explicitly would help.

Suggested edit (additive, near the end of question 3):

> "Specifically: confirm that `gen_ai.request.model` is captured on every chain entry representing a model call so the institution's model-inventory completeness check can consume chain data."

Again, this is wording. The substance is in `ai-policy-alignment.md`; the committee brief defers to it implicitly.

**M-3. Customer-dispute-procedures §"Reproducibility" could name the MRM committee's reproduction-result-handling threshold.**

The institution's response template commits to "consistent variation in your favor would prompt us to revisit the original decision." That commitment is correct; the institution-protective addition would be naming the MRM committee's threshold for what counts as "consistent variation." Today the threshold is implicit in "the MRM committee establishes the institution's reproduction posture"; an institution standing up its first dispute-response procedure would benefit from a worked threshold (e.g., "the MRM committee defines significant variation for each model class; for credit-decision models, the typical threshold is X% of re-runs producing a different decision over Y re-runs").

This is a deliberately soft suggestion — the institution's MRM committee is the right authority to set the threshold, not the spec. The suggestion is to add a line saying "the institution documents its variation threshold in its MRM-committee charter or model-risk policy" so a first-time institution knows where the threshold lives.

### OBSERVATIONS — no action required

**O-1. The "MRM Committee's qualitative judgment is unchanged" framing is exactly right.**

`MRM-COMMITTEE-BRIEF.md` lines 13 and 25 carry this framing: the chain is the evidence layer; MRM operates on top of the evidence. I want to reinforce this because vendors in this category routinely over-claim. An MRM director defending chain adoption uses this framing verbatim and the framing holds up under regulator challenge. The corpus carries the framing consistently across MRM-COMMITTEE-BRIEF, audit-committee-summary, and management-summary; no drift between the documents.

**O-2. The "integrity is necessary, not sufficient" recurring framing is the right inoculation against scope creep.**

`management-summary.md` §"What the chain does NOT do" is the plainest statement of this and it is the right placement. A CEO reading the summary understands the chain is one control alongside many others. An MRM committee chair reading the brief understands the chain doesn't replace the committee's substantive review. An audit committee chair reading the audit-committee-summary understands the chain is one input to broader audit oversight.

**O-3. The dual-algorithm post-quantum dispatch (spec §4.3.2 + §7 step 11 cases a-e) is forward-thinking and correctly scoped.**

For an MRM director planning multi-year posture, knowing the chain has a defined dual-algorithm transitional path (Ed25519 + post-quantum coexistence) is institution-friendly. The cases (both-valid PASS; single-algorithm during dual-posture PASS-WITH-ANOMALY; one-valid-one-invalid FAIL under strict) are correctly scoped — the institution's IR program does the interpretation of which case applies; the verifier reports both algorithm validations for the working-paper trail. This is the right separation of mechanism and policy.

**O-4. The reference-implementation provenance (Herald.Py v1.0; three Auditor rounds; cross-language byte-level test vectors) is mentioned at the right level.**

Design 02 line 5 carries the reference-implementation pointer without over-claiming. The FFIEC test-vectors corpus is seeded from the Herald.Py fixture's shape. For an institution evaluating chain adoption, knowing there is a reference implementation that has been audit-reviewed and produces byte-level conformance vectors lowers the implementation-risk perception. The corpus does not promise the reference implementation is the only conforming implementation, which is the right posture.

---

## Stopping criterion

**Declared at the top:** 0/0 — zero blockers, zero majors required to ship from an MRM-substance lens.

**Actual:** 0 blockers, 0 majors, 3 minors, 4 observations. The minors are wording tightenings that an editor pass closes; none would block an MRM committee from endorsing chain adoption, and none would block an SR 11-7 validation engagement from accepting chain output as integrity-bearing evidence.

**0/0 met.**

---

## What I would tell the institution's MRM director

If an institution's MRM director asked me whether to adopt this chain as part of the institution's AI-MRM framework, my answer would be: yes, with the standard caveats.

The standard caveats:

1. The chain provides integrity, not validity. Your MRM committee's qualitative judgment of model behavior is unchanged. The chain makes that judgment defensible; it does not replace it.
2. Operate the controls. The chain's integrity property depends on the institution operating the controls (key-fingerprint reconciliation, IKM rotation, HSM custody, append-only enforcement). Read the audit procedures (P-1 through P-25) and confirm your operations team can produce the evidence each procedure tests.
3. Pre-decide your IKM-disclosure posture. Customer disputes will arrive. Operational shape (HSM-mediated verification) for CFPB inquiries and regulator complaints; protective-order shape only when the dispute escalates to litigation. Documenting the threshold in your customer-dispute procedure now saves the institution weeks during a live dispute.
4. Configure `gen_ai_parameters` to the full SR 11-7 reproducibility surface. The chain integrity-binds whatever you put there. Putting only decoding parameters and not capturing the system prompt content-hash or the retrieval context means your reproductions can lie. The RECOMMENDED contents in spec §4.4 are the right floor.
5. Chain the ECOA adverse-action notice translation. The recommended posture is the right posture from a customer-dispute-defensibility perspective.

With those five operational commitments, the chain is the integrity-bearing component of an AI-MRM framework that I would defend in front of an OCC examiner, an EU AI Act conformity assessor, or a CFPB enforcement attorney. None of the five is unusual; all five are standard MRM-program operating practice once the chain is in place.

**Sign-off.** Dr. Sandeep Kumar Reddy.
