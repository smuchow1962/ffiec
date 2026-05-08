# Round-16 review — FFIEC Cybersecurity Specialist Examiner

**Reviewer.** Yejide Adeleke, Cybersecurity Specialist Examiner, FDIC. 11 years on the cyber-specialist track; current rotations on community-bank IT examinations and the inter-agency cyber working group that consumes outputs from the FFIEC computer-security incident notification rule. I work CSF 2.0 assessments at scale — meaning the artifacts I read have to map cleanly to subcategories in a way another examiner can pick up cold.

**Materials read.** spec/chain-of-custody-v1.md §10 + §7; design/00-overview.md; design/09-threat-model.md; regulator-pack/CSF-2.0.md; regulator-pack/ai-policy-alignment.md; incident-response-playbook.md; cloud-hsm-guide.md; supply-chain.md; edge-and-federated-ai.md.

**Persona discipline.** I have not seen prior rounds of this material. What's below is a first-look reaction to the artifacts as shipped. Where I describe a strength I'm describing what I observe in the document; where I describe a gap I'm describing what I cannot find as a cybersecurity examiner working a CSF 2.0 examination.

**Stopping criterion.** 0 / 0 — when I have zero new gaps and zero new questions on a re-read of the materials, I stop. This is my first read; that bar is not yet met. See §5 for the gap and question count.

---

## 1. What lands well — the artifacts I would carry into an examination

### 1.1 The four-attack indexing in CSF-2.0.md

The table in `regulator-pack/CSF-2.0.md` ("Four attack approaches mapped to CSF 2.0 subcategories") is the single most useful artifact in the pack from my seat. As a cybersecurity examiner I am building an evidence-to-subcategory matrix every time I open an institution. The four-attack indexing gives me four named operational events tied to four CSF subcategories tied to four IR scenarios — that is exactly the working pattern I use to sample-test PR.DS-06 evidence against an institution's log store.

The headline wins:

- **Per-attack named operational artifact.** `chain.verification_failure step=9` for §2.1, `step=8` for §2.2, `step=10` and `step=11` for §2.3, the SLSA + cosign + GPG bundle for §2.4. I can write a working paper that says "sample 5 events of type X from the institution's WORM log and confirm each fired the expected IR scenario disposition" without inventing the evidence trail. That is not normal; usually I am the one inventing it.

- **PR.DS-06 boundaries are explicit.** The "Subcategories the chain does NOT satisfy" table is unusually honest. PR.DS-01 (data at rest) and PR.AA principal-level identity are correctly disowned. An institution that over-claims PR.DS-06 coverage is a finding I can write; the pack giving me the boundary up front means I am writing a finding against an institution that ignored the pack, not against the pack itself.

- **Operational events bound to subcategories.** The table at lines 68-91 of `CSF-2.0.md` resolves each operational event to a primary CSF subcategory and a secondary set. The evidence binding is granular enough that I can sample-test by event type rather than by control as a whole. That cuts the time-on-engagement materially.

### 1.2 The 36-hour clock-start triage matrix

The matrix in `incident-response-playbook.md` lines 340-353 is the artifact I have been wanting to see in a vendor pack since the FFIEC computer-security rule landed. The rule is clear about the 36 hours; the rule is not clear about the determination point. Most vendor packs hand-wave this. This one does not.

What works:

- **Per-scenario clock-start trigger.** Each scenario has a named determination event that starts the clock. Scenario 1 (chain hash mismatch) starts the clock when investigation rules out SDK defect. Scenario 7 (key_fingerprint mismatch) starts the clock at triage case 4 only. Scenario 9 (truncation) almost never starts the clock. That is the right shape — it matches how counsel actually decides.

- **Concurrent multi-scenario rollup rule** (lines 357-359 of the playbook). The constellation of Scenario 1 + Scenario 7 + Scenario 8 within a 60-minute window being treated as one incident is exactly the operational reality. Treating each clock independently double-counts incidents in regulator notifications, which is a finding I have actually written against a community bank. The rollup rule preempts that finding.

- **Federal-regulator routing per institution charter** (lines 365-375). The CFR citations per charter type are correct for the four charter types I see most often (national bank → 12 CFR Part 53; state member → 12 CFR Part 225 App. F; state non-member → 12 CFR Part 304 Subpart C; FCU → 12 CFR Part 748 App. B). The state-side cadence note ("varies by state") is the right level of caveat — a vendor pack that pretended to enumerate every state's timing would be wrong within a quarter.

- **CIRCIA composability** (lines 388-403). The CIRCIA-without-FFIEC scenarios are correctly enumerated. Scenario 5 (sealing delay associated with broader cyber incident) is the canonical case where CIRCIA fires and FFIEC may not — that is how the inter-agency working group has been thinking about it.

### 1.3 Verifier supply-chain discipline

`supply-chain.md` is the strongest artifact in the pack for GV.SC and ID.RA-09 evidence. The chain-of-defenses (cosign + GPG fallback + cold-DR-key + reproducible builds + SLSA L3 + WORM-retained mirror audit log) is the right composition. The specific things I would carry into an examination:

- **The cold-DR-key lifecycle** (lines 104-119). Three-key separation (cosign holders + standard-GPG holders + cold-DR-key holders), 60-month rotation, annual dry-run with `KEY-DR-DRYRUN-{year}.asc` signed under standard cosign or GPG — that is governance I can sample. The institution-side consumption procedure (download, validate, archive, IR Scenario 11 trigger if missed) gives me the institution-side artifact to sample-test against. Without that consumption procedure I would have written a "documented-but-unexercised control" finding under GV.SC-04. The procedure preempts that.

- **Mirror audit-log retention and integrity** (lines 197-202). 7-year retention on WORM storage for the re-signing-mirror pattern's bridge invariant is the right call. CloudWatch Logs default retention is shorter than 7 years; institutions get this wrong routinely and my examination teams have written it as a finding before. Naming "AWS CloudWatch Logs with locked retention policy + S3 Object Lock" as an acceptable pattern, alongside Azure and GCP equivalents and on-prem WORM, gives me exactly what I need to confirm the institution's choice maps to an acceptable backend.

- **SLSA L3 institutional consumption MUST-tier** (lines 316-325). Publishing the L3 attestation is necessary but not sufficient — the institution must consume it with `slsa-verifier`, the consumption log is retained for the binary's deployment lifetime, and `slsa-verifier` itself must be signed and trust-anchored. That last point is the gap I would otherwise write up as a "supply-chain-trust-path-with-an-unsigned-verification-tool" finding. The pack closes it.

### 1.4 The Adversary I reception procedure (regulator-fingerprint rotation)

`09-threat-model.md` §2.9 lines 199-212 — the institution-side reception procedure for a regulator-held fingerprint rotation, with three named operational events (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`).

This addresses a real problem. Trust-anchor rotation is the kind of event where an institution operates without a procedure for years until the rotation happens, then improvises, and the improvisation is the finding. The three-event sequence with the reception-failure sub-variant (forged notice → IR Scenario 11 sub-variant → 36-hour clock starts at forgery determination) gives me a clean evidence trail to sample.

The "Examiner-laptop hygiene is a regulator-side control, NOT an institution-side control" note (lines 209-211) is correct and clears up a confusion I see in vendor packs routinely. Examiner-laptop posture is the FFIEC's own examiner-IT program; institution-side reviewers should not write a CC6 finding against an institution for a residual the institution does not operate. The pack getting this right means I am not arguing it with the SOC engagement partner downstream.

### 1.5 Spec §10 normative operational requirements

Section 10 of `spec/chain-of-custody-v1.md` reads like an examination checklist:

- **§10.1 weekly key-fingerprint reconciliation** at no-more-than-weekly cadence with `master.reconciliation_completed.unmatched_count` as the operational evidence. Bounds the master-compromise detection window to one week + remediation. That is a number I can put in a working paper.
- **§10.2 operational events catalog.** I get a full list with retention coupled to chain events.
- **§10.3 application-level + database-role-level append-only enforcement.** Defense in depth I can verify with two checks (codebase grep + database role grant query).
- **§10.5 HSM custody — `sign`-only seal-job role; separate authorization for `extract` / `delete` / `import`; separation of duties between seal-job operator and HSM administrator REQUIRED at institutions where role separation is operationally feasible.** The "where operationally feasible" caveat with documented dual-control as a compensating control accommodates community banks without weakening the spec.
- **§10.6 IKM minimum 32 bytes** with the dual-rationale (RFC 4868 §2 HMAC keying + 16-byte fingerprint offline-grind defense). Cited in R8 of the residual-risk register.
- **§10.7 software-key adapter compile-time exclusion + `kms_handle_uri = "plaintext-dev"` stamp + verifier refusal under `--strict`.** Three-layer defense against the dev-key-in-production misconfiguration I have seen in actual examinations. Run-time gating not being sufficient is exactly correct.
- **§10.9 IKM registry retention coupled to chain-event retention.** The premature-retirement failure mode (R11) is bounded by a normative requirement.
- **§10.10 rotation crossing the seal boundary** with `key_versions = [old, new]` on the day-after seal — the operational reality of rotation timing is captured normatively, not waved off.

Section 10 is the section I would copy into an examination request for evidence as the "what we expect to see" list.

### 1.6 §7 verifier procedure ordering

Section 7 of the spec puts cheap rejections before expensive ones in a way that closes specific failure modes:

- **Step 1 most-specific-first format-version check.** v2 file fails at step 1, not buried at step 9.
- **Step 4 cross-chain-lift check before MAC compute.** Cross-tenant or cross-run lift fails at step 4 with a specific message, not at step 9 with a generic MAC mismatch.
- **Steps 7 and 8 IKM lookup + fingerprint check both before MAC compute.** Botched-rotation failure mode surfaces at the fingerprint check, not in a MAC-mismatch storm. The constant-time-comparison requirement at step 8 (and §10.8) is the right discipline.
- **Step 9 verifier feeds `expected_prev_hash` (structurally walked) into the MAC recompute, NOT `entry.prev_hash`.** This is the R10 defense — even if a future maintainer relaxes the structural check, the MAC compute still uses the value derived from the previous entry's payload_hash, so the relaxation does not become a footgun. That is unusually careful design.
- **Step 11 algorithm/key-type mismatch reported as a specific message** rather than a generic signature-verification-failed. That helps examiners working a dual-algorithm transitional period (R13) distinguish "wrong-algorithm-public-key" from "key-actually-compromised."
- **Step 11 dual-algorithm dispatch with the five cases (a)-(e) and the Severe-regardless-of-bracket call on case (e).** Case (e) — both signatures present, one valid + one invalid — being Severe regardless of the verifier's PASS-WITH-ANOMALY vs FAIL bracket is the right call. That is the case where one of the algorithms has been broken OR one of the seals is forged under a compromised algorithm-specific key, and the institution's IR program does the interpretation per Scenario 12. The verifier does not pretend to know which.

## 2. The novel concepts that change how I'd run an examination

### 2.1 Per-entry `key_fingerprint` checked before any MAC compute

The §4.1 inviolate property #3 — fingerprint stamped on every entry, verifier asserts looked-up IKM produces the fingerprint BEFORE computing any MAC — is a control I have not seen articulated this cleanly elsewhere. The operational consequence is that a botched rotation that re-uses `key_version=1` for a different IKM is detected at the fingerprint check with a specific message ("looked-up IKM does not match the entry's recorded fingerprint"), not buried in a MAC-mismatch storm.

For an examiner, this is the difference between an institution being able to triage a fingerprint mismatch in minutes (Scenario 7 triage tree: rotation in flight? backup restored? tenant row restored? unauthorized substitution?) versus spending hours sorting through a MAC-mismatch storm without a clear root cause.

The IKM-minimum-32-bytes requirement (§10.6) protects the fingerprint from offline-grinding — the public 16-byte fingerprint is not offline-grindable against an IKM with the 32-byte minimum. R8 of the residual-risk register cites this.

### 2.2 Algorithm binding into `sign_payload`

The §4.3 sign_payload includes `algorithm` as a normative line. R13 of the residual-risk register names the defense: an attacker with a valid Ed25519 signature cannot present it as a Dilithium signature even if `public_key_id` happened to match. This is the JWT `alg=none` / SAML algorithm-substitution lesson applied at the seal layer.

The dual-algorithm coexistence rules in §7 step 11 case (b)-(e) are the operational manifestation. The institution's IR program runs Scenario 12 when case (e) fires; the verifier's working-paper convention records both algorithm validations regardless of bracket; the examiner cites `regulator-pack/finding-language.md` row "11 (dual-algo) co-signed seal failure" (Severe MRA) regardless of bracket.

This is forward-looking design. Cybersecurity examiners working post-quantum migration in five years will be reading this section and finding the answers already there.

### 2.3 Variant B per-algorithm `sign_payload`

The §4.2 schema note on the `signatures` list:

> Each signature in the list MUST cover its own algorithm-bound `sign_payload` (Variant B): each algorithm's `sign_payload` is constructed per §4.3 with that algorithm's identifier in the second line. A single shared `sign_payload` covering all algorithms (Variant A) is non-conformant — it leaks the algorithm-confusion defense by letting an attacker present an algorithm-X signature on a payload that names algorithm-Y.

That paragraph closes a subtle attack class I have not seen articulated in vendor packs at all. Variant A would have looked superficially fine — same payload, two signatures — but it would have re-opened the algorithm-confusion attack the §4.3 sign_payload bind closed. Calling Variant B normative and Variant A non-conformant is the right call.

### 2.4 Dual-algorithm verifier-version timing-overlap

Lines 378-386 of the playbook — the institution-vs-regulator verifier-version timing-overlap matrix during the dual-algorithm transitional period.

The load-bearing case is the third row: institution sees case (e) FAIL, regulator sees PASS under its X-only verifier because the regulator's verifier doesn't process the `signatures` list. The regulator is structurally blind to the Y-algorithm validation. The institution must proactively notify the regulator and supply the case-(e) verifier output as supplementary evidence.

This is the kind of operational reality that gets discovered at examination time and creates an emergency. Pre-coordinating verifier-version compatibility BEFORE the multi-year transitional period begins is the right operational posture, and the SHOULD on the regulator's IT examination program upgrading to a Y-algorithm-aware verifier within ~6 months gives the project a number to track.

I have not seen a vendor pack articulate this concern at all before. It is the kind of thing that surfaces only when someone has thought through the post-quantum transition end-to-end.

### 2.5 The R8 fingerprint-not-offline-grindable framing

R8 in the residual-risk register:

> Session key never logged; key handles only. The public per-entry `key_fingerprint` is NOT offline-grindable due to spec §10.6 IKM minimum (32 bytes per RFC 4868) — an attacker observing chain entries cannot brute-force the IKM from the fingerprint at acceptable cost.

Most threat models would have left R8 as "session key leakage to logs — operational concern" without addressing the publicly-observable-fingerprint concern at all. The fingerprint is on every entry; an attacker observing chain entries can see it; an examiner asking "could that fingerprint be inverted?" is asking the right question and the spec has the answer pre-loaded.

## 3. Where I have gaps or questions on first read

This is what I would walk into a working-paper meeting with on a first-read engagement. Each item is something I cannot resolve from the materials as shipped.

### 3.1 Reconciliation cadence justification at small institutions

Spec §10.1 says weekly cadence is the floor; "institutions with elevated risk profile MAY operate reconciliation at higher frequency (daily, hourly, or continuous)." The corresponding R8 / R9 / Scenario 7 framing carries the weekly assumption.

Question: for a community bank with low AI-decision volume, is weekly the right floor, or is the floor actually "the institution's articulated risk-tolerance, with weekly as the default"? The §10.1 text says "no more than weekly cadence" which I read as a ceiling on the gap between reconciliations — i.e., the institution may not let a week pass without reconciliation. That is a strong floor for the master-compromise detection window. But the weekly cadence at a community bank with 50 AI decisions per week may produce more operational noise than signal.

I am not asking the spec to weaken the cadence. I am asking whether `regulator-pack/finding-language.md` or `audit-procedures.md` has a line about how to write the finding when an institution operates the reconciliation at the spec floor but the operational signal is dominated by noise (false-positive unmatched_count from documented-rotation-in-flight, etc.). The Scenario 7 triage tree handles the per-event response; I'm asking about the per-week posture.

**Gap count for this item: 1 question.**

### 3.2 Reconciliation evidence from §10.1 — what the actual evidence looks like

`master.reconciliation_completed` is named as the operational event with `unmatched_count`. That is fine for an examiner walking the institution's log store. What I cannot find is what the institution archives as the working-paper evidence: is it the single event line? The full per-(tenant_id, key_version) match table the reconciliation operated over? An attestation from the SOC team that the reconciliation ran with the expected match set? The audit-procedures P-6 reference suggests P-6 is the procedure but I do not have the P-6 text in the materials I read.

For an examiner sample-testing the reconciliation, I want to walk in and ask: "Show me the reconciliation evidence for the week of {date}." The institution should produce a specific artifact. The pack tells me the event fires; it does not tell me the artifact's shape.

**Gap count for this item: 1 question.**

### 3.3 Pattern A edge-AI fingerprint reconciliation cadence

`edge-and-federated-ai.md` describes Pattern A (per-device IKM in TPM/secure-enclave) and notes "session-key-id reconciliation on edge devices may show 'unmatched' entries during disconnection windows; the institution baselines the unmatched count for the edge fleet separately from cloud."

R12 in the residual-risk register cites Pattern A and "reconciliation cadence baseline tuned per fleet" as the v1.0 compensating control until v1.1 ships explicit physical-attacker handling.

What is missing is the institution-side baseline procedure. If an edge fleet of 10,000 devices produces N unmatched entries per week during disconnection windows, what N is the institution allowed to baseline as normal-operations? At what threshold does the unmatched_count cross from baseline-noise to genuine-signal? `at-scale-operations.md` may have this — I did not read it. The materials I read do not have a number, an institution-articulation procedure, or a sample baseline to reference.

For a cybersecurity examiner reviewing an institution's edge deployment, I would want to see the institution's baseline articulation in the control description. The pack should give institutions a starting shape for that articulation.

**Gap count for this item: 1 gap (missing institution-articulation procedure for edge fleet reconciliation baselines).**

### 3.4 The `gen_ai_model_identifier_missing` finding — disposition under PASS-WITH-ANOMALY

Spec §7 step 12a: under `--strict` FAIL, under non-strict PASS-WITH-ANOMALY. The reason is "control-completeness for SR 11-7 reproducibility, NOT chain-integrity."

Question for me as a cybersecurity examiner (not a model-risk examiner): if an institution's verifier output shows PASS-WITH-ANOMALY for `gen_ai_model_identifier_missing` on N% of model-call entries during the period, do I write a finding? My instinct is no — this is SR 11-7 reproducibility, which is the model-risk examiner's lane. But if the institution's control description lists `gen_ai.request.model` and `gen_ai.response.model` as integrity-bound fields and the verifier output shows missing identifiers, that is a control-description-versus-actual-behavior gap that is a CSF GV.PO-01 finding.

The ai-policy-alignment doc's Article 14 (EU AI Act human oversight) reference treats the model identifiers as load-bearing for effective oversight. If an EU-deployer institution operates the chain and the model identifiers are missing, the Article 14 evidence claim is partially broken. That feels like it should be written as a finding even when the institution operates only under FFIEC supervision domestically — because the institution's control description names the field and the field is missing.

What I want is a line in `regulator-pack/finding-language.md` that addresses the cybersecurity-examiner's disposition on `gen_ai_model_identifier_missing` PASS-WITH-ANOMALY at scale. Is it a control-completeness finding I write, or do I refer it to the model-risk examiner? My instinct is the former; the materials I read don't tell me.

**Gap count for this item: 1 question.**

### 3.5 Cross-tenant scope discovery in a vendor-hosted topology — the institution's notification procedure

Lines 360-363 of the playbook:

> **Cross-tenant scope discovery during investigation (vendor-hosted topology).** Scenarios 4 and 7 may discover during investigation that the suspected compromise spans more than one tenant in a multi-tenant deployment (vendor-hosted topology, shared-cloud-HSM at community-bank tier per `00-overview.md` §6.5). The 36-hour clock applies per institution; in a vendor-hosted deployment one investigation may produce N institution-side determinations on different timelines.

The vendor-vs-institution responsibility split is articulated. What I cannot find is the institution-side procedure for "the vendor told me on Tuesday at 14:00 that they have a credible Scenario 4 signal that may affect my tenant; the vendor's investigation will conclude on Thursday or Friday; my 36-hour clock starts at my determination, not at the vendor's notification." The institution's IR program needs a documented procedure for "what evidence does the institution gather between Tuesday's vendor notification and the institution's own determination?" so the institution's clock-start is defensible to the FFIEC.

The current text says "the institution's clock starts when the institution determines, not when the vendor first alerted." That is the right rule. What is missing is the institution-side decision procedure — e.g., "the institution gathers (a) the vendor's preliminary investigation findings, (b) the institution's own log-store sample of the suspect (tenant_id, key_version) range, (c) the institution's IR Commander's documented determination — and the determination point is the IR Commander's signed entry in the incident-management system."

Without that procedure, the institution risks under-determining (waiting until the vendor's investigation completes, missing the FFIEC clock) or over-determining (starting the clock at vendor notification, double-counting the incident with the vendor's own determination).

**Gap count for this item: 1 gap (missing institution-side decision procedure for cross-tenant scope discovery in vendor-hosted topology).**

### 3.6 Scenario 11 sub-variant — when the rotation notice is forged, the institution coordinates with the regulator

Lines 283-297 of the playbook describe the trust-anchor reception failure sub-variant. Containment: do not install the new fingerprint; retain the previous fingerprint as the active trust anchor until validation succeeds. Remediation: out-of-band cross-check via phone call, in-person visit, or regulator's encrypted-email system if separate from the notice channel.

Question: in case 1 (forged notice), the institution starts the 36-hour clock at the forgery determination. The institution then coordinates with the regulator on whether the forgery attempt indicates a broader attack. Two operational questions a cybersecurity examiner working the post-incident review would ask:

1. Does the institution's notification to the FFIEC under the 36-hour rule include the forgery attempt's evidence (the forged notice's bytes, the failed-validation log entry, the cross-channel disagreement evidence)? My instinct is yes — the regulator needs that evidence to investigate whether the attempt indicates broader compromise. The materials don't say.

2. If the institution's verifier was running against the regulator's published key-rotation channel directly (not via the institution's caching), would the institution still be able to complete a verifier run during the reception-failure window? My read of supply-chain.md is that the institution caches the trust anchors out-of-band, so a regulator-channel disruption does not immediately break the verifier. But that is the right thing to confirm in the playbook explicitly — the institution's verifier uses cached trust anchors, so the institution's verifier is not blocked by the trust-anchor reception failure.

Both are questions a post-incident review would ask. The materials should pre-load the answers.

**Gap count for this item: 1 question (specifically: confirmation that the institution's verifier continues to function against cached trust anchors during a reception-failure window).**

### 3.7 SLSA L4 candidacy timeline

`supply-chain.md` line 314: "SLSA Level 4 (which adds two-party review and reproducibility verification at the build service) is a v1.1 candidate; the current pipeline meets the reproducibility property but does not yet operate the two-party review at the build-service layer."

Question: does the spec working group have a target window for SLSA L4? A v1.1 candidate is good; a date is better. For an institution adopting v1.0 today, the institution's GV.SC posture has to articulate the institution's tolerance for v1.0 SLSA L3 vs v1.1 SLSA L4. If v1.1 is 18 months out the institution's posture is one thing; if v1.1 is 4 years out the posture is another.

The 30-day emergency-spec-patch SLA in §4.3.2 is a project-side commitment with a number. The SLSA L4 v1.1 candidacy does not have a number.

**Gap count for this item: 1 question.**

### 3.8 The `mirror.reconciliation_completed` event — institution-defined but bound to CSF subcategories

`CSF-2.0.md` lists `mirror.reconciliation_completed` as institution-defined and binds it to DE.CM-09 + GV.SC-04. The supply-chain doc names the reconciliation as the institution's continuous-monitoring control on the mirror's signature-validation discipline.

Question: what does the institution's IR Scenario disposition look like when `mirror.reconciliation_completed` shows a mirror-side signature-validation gap? The supply-chain doc says "fires an operational alert"; the IR playbook does not have a Scenario for it (Scenario 11 is project-side trust-anchor degradation, not mirror-side). I would expect the playbook to have a sub-scenario or a reference to a mirror-side scenario.

The realistic flow is: `mirror.reconciliation_completed` fires with an unmatched count > 0 → the institution's IR Commander triages: was this a mirror infrastructure outage (operational), a configuration drift at the mirror (operational), or a vendor-side compromise (security event)? The triage tree is symmetric to Scenario 7's. The materials do not give me that triage tree.

**Gap count for this item: 1 gap (missing IR sub-scenario or full scenario for mirror-side signature-validation gap).**

### 3.9 The CIRCIA matrix — what about Scenario 10 (backup integrity failure)?

Lines 389-401 of the playbook enumerate Scenario 5, 9, 7, 4, 6, 11, 11-sub, 12-i, 12-ii, 12-iii. Scenario 10 (backup integrity failure causing verifier failure) is not on the CIRCIA matrix.

Scenario 10 is described separately at lines 315-334. The disposition: "operational unless gap is unfillable; clock starts when the gap is unfillable AND the affected days remain unverifiable."

Question: under what condition does Scenario 10 cross the CIRCIA "substantial cyber incident" threshold? My read is: an unfillable gap that includes a material business decision affecting customers may cross the CIRCIA threshold even if the FFIEC trigger is also engaged. The matrix should cover Scenario 10 explicitly.

**Gap count for this item: 1 gap (Scenario 10 missing from CIRCIA matrix).**

### 3.10 Edge AI master-key custody — Pattern A vs Pattern B and the IR posture

`edge-and-federated-ai.md` describes Pattern A (per-device IKM in TPM/secure-enclave) and Pattern B (bulk session-key issuance at commissioning). R12 names Pattern A as the v1.0 compensating control until v1.1 ships explicit physical-attacker handling.

Question: under Pattern B, what does the IR playbook scenario look like when an edge device is suspected of physical compromise during the device's operational lifetime? Pattern A's failure mode is "the device's hardware secure element is defeated"; Pattern B's failure mode is "the batch of session keys issued at commissioning is exfiltrated and used by an attacker after device retirement, before the institution's reconciliation cadence catches the unmatched fingerprints."

The materials cover Pattern A's residual under R12. They do not articulate Pattern B's residual or the IR scenario for "physical compromise of an in-service edge device." The institution choosing Pattern B because the device fleet doesn't have secure-enclave capability needs a documented IR posture.

**Gap count for this item: 1 gap (missing IR scenario for physical compromise of in-service edge device under Pattern B).**

## 4. Use cases I would carry into an examination

Five examination shapes I would walk into using these materials:

### 4.1 Community-bank IT examination, low-volume AI-decision use case

The institution operates a small fraud-detection AI agent at one branch. Volume: ~200 decisions/day. Deployment: shared-cloud HSM (per `00-overview.md` §6.5 community-bank tier); daily seal cadence relaxed to weekly with examiner approval per `regulator-pack/examiner-approval-template.md`.

The CSF-2.0 mapping gives me the four-attack matrix; I sample-test:

- PR.DS-06 evidence via 5 random `chain.verification_failure` events from the institution's log store (any step). Confirm each fired the expected IR scenario disposition per the playbook.
- ID.AM-08 evidence via 1 sampled `master.reconciliation_completed` event for each of 4 weeks in the look-back period, confirming `unmatched_count` and the institution's IKM-roster cross-reference.
- GV.SC evidence via the institution's cosign + GPG + SLSA-verifier trust-anchor cache, with a sample of 1 binary deployment in the look-back period validated end-to-end.
- The institution's annual `KEY-DR-DRYRUN-{year}.asc` consumption log — one record per year for the look-back.

The pack gives me the artifact list per CSF subcategory. The institution gives me the artifacts. The working paper writes itself.

### 4.2 Mid-size bank, BYOC topology, dual-region deployment

The institution operates AI agents in two regions (us-east-1, us-west-2) with per-region tenant_id (`tenant_acme_prod_us_east_1`, `tenant_acme_prod_us_west_2`) per `00-overview.md` §6.4 multi-region resilience workaround.

The institution's CSF-2.0 evidence is per-region. I sample-test:

- PR.DS-06 evidence per region (one set of `chain.verification_failure` samples per tenant_id).
- DE.CM-09 evidence on the per-region reconciliation cadence — both regions show `master.reconciliation_completed` events at the spec-floor cadence.
- The institution's correlation procedure for cross-region events (described in the institution's control description per §6.4).

The materials give me the topology; the institution gives me the evidence; the working paper covers both regions.

### 4.3 Vendor-hosted multi-tenant deployment with shared cloud HSM

Per `00-overview.md` §6.3 vendor-hosted topology with per-tenant key separation. The institution has a SOC report on the vendor; the institution's own controls are limited to the `audit.*` payload schema and the IKM custodian's reconciliation procedure.

I sample-test:

- The vendor's SOC report covering the chain primitives (PR.DS-06 evidence comes from the vendor; the institution's evidence is the SOC report consumption).
- The institution's IR program for cross-tenant scope discovery during a Scenario 4 or 7 (per playbook lines 360-363) — confirm the institution has a documented procedure for the vendor-vs-institution responsibility split AND a documented decision procedure for the institution's clock-start (the gap I named in §3.5 above).
- The institution's regulator-fingerprint reception procedure per Adversary I (`09-threat-model.md` §2.9) — confirm the three operational events fire and the institution's verifier configuration update is change-managed.

The pack gives me the cross-tenant rollup rule and the responsibility split. The gap in §3.5 is the one I would write up as a finding if the institution's IR procedure is silent on the institution-side decision shape.

### 4.4 Post-incident review of a Scenario 7 triage case 4 (unauthorized substitution suspected)

The institution's verifier reported `key_fingerprint mismatch at seq N`. Triage cases 1-3 (rotation-in-flight, restored-backup, tenant-row-restored) all eliminated. Case 4 confirmed: suspected unauthorized key substitution. Clock started; FFIEC notification within 36 hours; CIRCIA notification within 72 hours.

I review the institution's incident-management system for:

- The triage decision evidence (cases 1-3 elimination evidence).
- The clock-start determination (IR Commander's signed entry).
- The 36-hour FFIEC notification text matching the playbook's communication template.
- The CIRCIA 72-hour notification.
- The post-incident review per playbook lines 419-429.

The pack gives me the triage tree and the notification triggers; the institution gives me the evidence trail. Where a rollup with concurrent Scenario 1 alerts applied (lines 357-359), I confirm the IR Commander documented the rollup decision.

### 4.5 GV.SC examination on the institution's verifier supply-chain posture

The institution adopts the v1.0 verifier. I sample-test:

- The institution's cached cosign + GPG + cold-DR + SLSA-verifier public keys, cross-referenced against the spec text and the project's published fingerprints.
- The institution's reproducible-build log for the deployed binary.
- The institution's `slsa-verifier` consumption log showing the L3 attestation validated at deployment time.
- The institution's `KEY-DR-DRYRUN-{year}.asc` archive for the previous year.
- For re-signing-mirror pattern: the mirror's audit log on WORM storage with 7-year retention, the `mirror.reconciliation_completed` events, and the bridge documentation.

The supply-chain doc gives me each artifact's expected shape. The institution gives me the artifact. The pack lets me write the GV.SC working paper without inventing the evidence trail.

## 5. Stopping criterion check and gap count

**0 / 0 stopping criterion.** I stop when zero new gaps and zero new questions emerge on a re-read.

**This pass.** First read.

| Item | Type |
|---|---|
| 3.1 Reconciliation cadence floor at small institutions | Question |
| 3.2 Reconciliation evidence artifact shape | Question |
| 3.3 Edge fleet reconciliation baseline procedure | Gap |
| 3.4 `gen_ai_model_identifier_missing` PASS-WITH-ANOMALY disposition for cybersecurity examiner | Question |
| 3.5 Institution-side decision procedure for cross-tenant scope discovery | Gap |
| 3.6 Verifier behavior during Scenario 11 sub-variant reception-failure window | Question |
| 3.7 SLSA L4 v1.1 candidacy target window | Question |
| 3.8 IR sub-scenario for mirror-side signature-validation gap | Gap |
| 3.9 Scenario 10 missing from CIRCIA matrix | Gap |
| 3.10 IR scenario for physical compromise of in-service edge device under Pattern B | Gap |

**Count.** 5 gaps + 5 questions. Stopping criterion not met.

**Re-read trigger.** Address the items above and I re-read. The artifacts I would expect to see next round to drive these to 0 / 0:

- A line in `regulator-pack/finding-language.md` or `audit-procedures.md` covering the cybersecurity-examiner's disposition on PASS-WITH-ANOMALY at scale (closes 3.4).
- A reconciliation-evidence-artifact-shape paragraph in `audit-procedures.md` P-6 or in `regulator-pack/CSF-2.0.md`'s operational events binding (closes 3.2).
- An edge-fleet baseline-articulation procedure in `at-scale-operations.md` or `edge-and-federated-ai.md` (closes 3.3).
- An institution-side cross-tenant scope-discovery decision procedure paragraph in the playbook's vendor-hosted-topology edge case (closes 3.5).
- A confirmation paragraph in the playbook Scenario 11 sub-variant on cached-trust-anchor verifier continuity during a reception-failure window (closes 3.6).
- A target window for SLSA L4 in `supply-chain.md` (closes 3.7).
- An IR scenario or sub-scenario for mirror-side signature-validation gap (closes 3.8).
- Scenario 10 added to the CIRCIA matrix (closes 3.9).
- An IR scenario for physical compromise of in-service edge device under Pattern B in `edge-and-federated-ai.md` or as a new IR scenario (closes 3.10).
- A clarification paragraph in `regulator-pack/CSF-2.0.md` or `audit-procedures.md` on the §10.1 weekly cadence floor at small institutions and the noise-vs-signal posture (closes 3.1).

If next round closes those items, the 0/0 stopping criterion is met from my seat.

## 6. Bottom line

This is the strongest cybersecurity-examination-ready vendor pack I have read on first look. The four-attack indexing, the 36-hour triage matrix with charter-routing and CIRCIA composability, the supply-chain discipline with cold-DR-key lifecycle and SLSA L3 institutional consumption, the Adversary I reception procedure, and the §10 normative operational requirements are five artifacts I would not have to invent — they would map straight into my working papers.

The gaps and questions in §3 are real but they are operational refinements, not structural issues. The chain's design holds together; the pack's CSF mapping is honest about the boundaries; the IR playbook's clock-start triage is the right shape for the FFIEC computer-security rule.

I would write a clean PR.DS-06 evaluation against this pack. The gaps in §3 are the items I'd attach as supplemental questions to the institution; none of them are findings against the spec or the pack.

Yejide Adeleke
Cybersecurity Specialist Examiner, FDIC
Round-16, first read
