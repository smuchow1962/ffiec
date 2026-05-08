# Round-12 review — FFIEC Cybersecurity Specialist Examiner

## Persona

**Name.** Mateusz Kowalski.
**Role.** FFIEC Cybersecurity Specialist Examiner with the OCC. Thirteen years on the cyber-specialist track; recent secondment to ENISA on the EU-US cybersecurity cooperation program. Active member of the FFIEC working group on AI cybersecurity examination.
**Reading angle.** Threat model, IR playbook, NIST CSF 2.0 mapping, supply-chain SLSA L3 trust path, 36-hour computer-security incident notification triage, dual-algorithm transitional-period verifier behavior. I read the way a cyber-specialist reads a control description: I look for the named control, the operational event the control produces, the sample I can pull at examination time, and the path from a verifier failure to a regulator-notifiable determination. If any link in that chain is fuzzy, the control is not examinable.

## What I read

- `spec/chain-of-custody-v1.md` — close-read on §10 (operational requirements) and §7 step 11 (dual-algorithm verifier dispatch cases (a) through (e)), including the load-bearing co-signed-seal-failure case (e).
- `docs/design/00-overview.md` — system shape, four-attack table at §2, trust-boundary sketch.
- `docs/design/09-threat-model.md` — Adversaries A through I (with the new Adversary I institution-side reception procedure for regulator-held fingerprint rotation), R1–R13 residual-risk register, dual-algorithm transitional period at §2.8.
- `docs/regulator-pack/CSF-2.0.md` — the four-attack-to-subcategory table, the operational-events binding table (now including `mirror.reconciliation_completed` and the SLSA-verifier validation log), the GOVERN-function alignment one-pager.
- `docs/incident-response-playbook.md` — Scenarios 1–10, the 36-hour triage matrix, the named edge cases (concurrent rollup, vendor-hosted determination split, federal-regulator-routing-per-charter, CIRCIA-only triggers without FFIEC trigger), and Scenario 4+3 mapping for co-signed seal failure.
- `docs/cloud-hsm-guide.md` — conformance bar, per-provider provisioning, validation and DR sections.
- `docs/supply-chain.md` — SLSA L3 with institutional consumption MUST-tier and SLSA-aware tool-choice normative; cold-DR-key lifecycle (rotation cadence, key-holder lifecycle, annual dry-run, health check); GPG and cosign recovery; cosign recovery during active examination; mirror WORM + reconciliation control.
- `docs/regulator-pack/ai-policy-alignment.md` — Treasury RMF, EU AI Act Articles 12 / 14 / 26, DORA Articles 5 / 6 / 8 / 17 / 28, OCC supervisory posture.
- `docs/edge-and-federated-ai.md` — edge / federated / on-device / multi-agent patterns.

I did not read any prior-round feedback files. This is a fresh-eyes review against the shipped corpus.

## Headline impression

This is the most carefully evidenced AI-decision audit-trail spec I have read. Three constructions stand out for the cybersecurity examination workflow specifically.

First, the per-event `key_fingerprint` check at verifier step 8 — fired BEFORE any MAC compute — gives me a precise named failure mode for cross-tenant configuration drift (`key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint`) rather than a MAC-mismatch storm I would have to triage backwards. Combined with the §10.1 weekly fingerprint reconciliation control and the IR playbook Scenario 7 four-case triage tree, the alert / control / IR triage pipeline is something I can sample-test in one sitting.

Second, the four-attack-to-CSF-subcategory table at `CSF-2.0.md` lines 37–43 is exactly the working artifact a cyber-specialist examiner needs. I can sit with an institution's evidence repository, pick an attack approach, follow it to a CSF subcategory, and pull the named operational event the institution should be able to produce. The operational-events binding table now includes `mirror.reconciliation_completed` and the SLSA-verifier validation log entry, so the §2.4 examination-time-tooling-subversion attack resolves cleanly to GV.SC-04 and ID.RA-09 with named evidence in the same document — I do not have to reach across the corpus to find the binding.

Third, the spec §7 step 11 dual-algorithm dispatch cases (a) through (e) cover the post-quantum migration's verifier-side behavior with examiner-grade precision. Case (e) — both signatures present, one valid + one invalid — is the load-bearing case that distinguishes broken-algorithm-with-the-other-still-good (the un-broken algorithm provides integrity) from compromised-algorithm-specific-key (the un-broken algorithm confirms the un-compromised half) without asking the verifier to interpret which case applies. The verifier records both validation results and the institution's IR program does the interpretation. That separation of mechanism from policy is correct.

The concerns I raise below are about completeness on three specific edges. The foundation is sound.

## Five questions before sign-off

### Question 1 — Does the Adversary I institution-side reception procedure for a regulator-held fingerprint rotation contain enough to operate in real conditions?

The Adversary I residual at `09-threat-model.md` §2.9 (lines 189–211) names the regulator-held fingerprint as the trust anchor that makes the institution-published public key trustable. The institution-side reception-and-validation procedure at lines 199–207 enumerates five steps: receive the rotation notice, validate authenticity against the regulator's GPG signature, update the institution's verifier configuration, re-validate historical verifier reports, document the rotation in the control-evidence repository.

I tested this against my actual examination working pattern. A few specifics deserve to be tightened before this becomes a control I can sample-test cleanly.

**Step 2 — validation of the rotation notice's authenticity.** The text says the procedure "SHOULD be co-signed by the regulator's standing identity" and cross-checked against the regulator's authenticated-domain channel. The conformance keyword "SHOULD" is correct here because it is a regulator-side procedure the institution receives, not an institution-side control. But the institution's IR playbook does not name the operational event the institution emits when it RECEIVES a rotation notice and runs through this procedure. Without an operational event, I cannot sample-test that the procedure was operated. The institution-defined `mirror.reconciliation_completed` pattern is the precedent — an institution-defined operational event named, for example, `regulator_fingerprint.rotation_received` carrying `(received_at_utc, validated_at_utc, validation_method, regulator_signature_id, old_fingerprint, new_fingerprint, archived_old_fingerprint_at)` would let me pull a sample at examination time and confirm the reception procedure was operated.

**Step 4 — re-validation of historical verifier reports.** The text says the institution re-runs historical verifier reports "to confirm the institution's verifier output remains stable across the rotation." The action is right; the operational evidence is missing. A `regulator_fingerprint.rotation_revalidation_completed` event carrying `(revalidation_period, reports_revalidated_count, reports_unstable_count, action_taken)` would be the artifact. The "reports_unstable_count" field is the load-bearing one — if any historical verification flips PASS-to-FAIL across the rotation, that is a finding, and the institution needs to surface it in time for the institution to investigate (rotation operationally botched on the regulator side, an old fingerprint cached against a new key, etc.).

**Step 5 — documentation in the control-evidence repository.** This is the closing step but the playbook does not bind it to an IR scenario. A regulator-held fingerprint rotation that fails the institution-side reception procedure (forged notice, invalid signature, archive-failure) is itself a candidate for IR Scenario 3 (signature verification failed) under a sub-variant ("trust-anchor-rotation failure"). Right now Scenario 3 covers "the seal record's signature does not validate against the published key" but not "the regulator-held trust anchor changed and the institution failed to update or validated against a forged rotation notice." The two are different failure modes with different responses, and the IR playbook should explicitly name the second.

**Status: Partial.** The reception procedure is documented; the operational events for sample-testing it and the IR sub-variant for failure are missing.

### Question 2 — Does the cold-disaster-recovery-key lifecycle in `supply-chain.md` produce the institution-side evidence I need to verify the dual-compromise fallback is operationally exercised?

The cold-DR-key lifecycle at `supply-chain.md` lines 104–111 names four properties: 60-month rotation cadence aligned with NIST SP 800-57; key-holder lifecycle separate from key-rotation cadence with emergency rotation within 30 days of a key-holder departure; annual dry-run cadence with a published `KEY-DR-DRYRUN-{year}.asc` attestation; health check on the same annual cadence with deviation-triggered out-of-band notification to the FFIEC working group.

This is project-side governance, not institution-side. But the institution consumes the dry-run attestation as evidence the dual-compromise fallback is operationally exercised. I tested whether the institution-side consumption is documented well enough for me to sample-test it.

**The institution-side trust-anchor table.** The `supply-chain.md` "Trust anchors" table at lines 144–149 names cosign public key, GPG public key, and toolchain version pinning. It does NOT name the cold-DR key's public key fingerprint as an institution-side cached trust anchor. If the dual-compromise event fires and the project publishes a recovery notice signed by the cold-DR key, the institution must validate that notice against a cached cold-DR public key fingerprint. The table should include the cold-DR public key fingerprint as a fourth row, with the same caching procedure as the cosign and GPG keys.

**The institution-side dry-run-attestation consumption procedure.** The dry-run attestation `KEY-DR-DRYRUN-{year}.asc` is published annually. The institution consumes it as control evidence — but the operational pattern at lines 220–229 ("For each new release") does not include "validate the annual dry-run attestation and archive in the control-evidence repository" as a step. Without this, the dry-run attestation is published but the institution-side consumption is implicit. The cybersecurity examiner asking "does the institution validate the dry-run attestation annually" will not have a sample to pull.

**The institution-side health-check escalation path.** The text says deviations in the cold-DR key's hardware health "trigger out-of-band notification to the FFIEC working group." This is project-side. The institution receives the notification through the FFIEC working group's notification list. But the IR playbook does not name an institution-side scenario for "received a health-check-deviation notification on the cold-DR key from the project." This is a low-likelihood but high-consequence event — a degraded cold-DR key combined with a future cosign + GPG dual compromise removes the institution's last trust anchor for the verifier binary. The IR playbook should add a Scenario 11 (project-side trust-anchor degradation) with the institution's response posture: pause new verifier deployments, escalate to the chain operations lead, document in the control-evidence repository.

**Status: Partial.** The cold-DR key lifecycle is well-documented project-side; the institution-side trust-anchor caching, dry-run-attestation consumption, and health-check escalation IR scenario are missing.

### Question 3 — Does the spec §7 step 11 case (e) co-signed-seal-failure mapping to IR Scenarios 4+3 produce a clean clock-start determination?

Case (e) at `chain-of-custody-v1.md` lines 322–328 is the load-bearing dual-algorithm case: both signatures present, one valid + one invalid. Under `--strict` the verifier reports FAIL with `co-signed seal failure: algorithm X validated, algorithm Y did not`. The spec correctly says the verifier does NOT attempt to interpret which case applies (broken-algorithm-with-the-other-still-good versus compromised-algorithm-specific-key); the institution's IR program does the interpretation.

I traced the mapping. The IR playbook does not explicitly name a Scenario for case (e). The closest mappings are Scenario 3 (signature verification failed) and Scenario 4 (master key compromise suspected). Both apply: case (e) is a signature-verification failure on one algorithm AND a candidate for an algorithm-specific signing-key compromise.

The 36-hour clock-start matrix at `incident-response-playbook.md` lines 286–298 includes Scenario 3 (clock starts at alert receipt) and Scenario 4 (clock starts when the institution determines the suspect signal is credible). The matrix does not include a row for case (e). This matters because case (e) has a different determination posture from either Scenario 3 or Scenario 4 alone.

**The case (e) determination tree the institution actually needs.**

- If the institution's threat intelligence confirms the broken-algorithm-with-the-other-still-good case (e.g., NIST published a deprecation notice for the failing algorithm), the institution's clock-start posture is "no immediate 36-hour clock; coordinate with the regulator on the broken-algorithm migration timeline." This is a control-completeness finding under spec §4.3.2, not a chain-integrity finding.
- If the institution's threat intelligence suggests the compromised-algorithm-specific-key case (e.g., HSM tamper detection on the algorithm-specific signing-key partition), the 36-hour clock starts under Scenario 3 immediately because the un-broken algorithm's signature confirms the un-compromised half but the broken half is suspected unauthorized signing.
- If neither case is yet confirmable, the institution's posture is "treat as Scenario 4 with reduced clock-start threshold; the un-broken algorithm provides interim integrity assurance while investigation proceeds."

The IR playbook should add a Scenario 11 (or extend Scenario 3 with a sub-variant) covering case (e) with this three-branch determination tree, and the 36-hour matrix should add a row.

**Status: Partial.** The verifier mechanism for case (e) is documented; the IR-side clock-start determination is implicit and needs a named scenario row.

### Question 4 — Does the operational-events binding table in `CSF-2.0.md` cover the new institution-defined events the cold-DR-key consumption and Adversary I reception procedure require?

The operational-events binding table at `CSF-2.0.md` lines 68–86 catalogs the spec §10.2 events plus the institution-defined `mirror.reconciliation_completed` and the SLSA-verifier validation log entry. Both additions are correct.

The table does NOT include the institution-defined events Question 1 and Question 2 above identify as missing:

- `regulator_fingerprint.rotation_received` (institution-defined; bound to GV.SC + RC.CO-03 — recovery activities communicated)
- `regulator_fingerprint.rotation_revalidation_completed` (institution-defined; bound to PR.DS-06 + DE.AE-02)
- `cold_dr_key.dryrun_attestation_validated` (institution-defined; bound to GV.SC-04 + ID.RA-09)
- `cold_dr_key.health_check_deviation_received` (institution-defined; bound to DE.CM-09 + GV.SC-04)

These are pre-conditions for the partials in Question 1 and Question 2 to close. Once those events exist (institution-defined per the precedent of `mirror.reconciliation_completed`), the binding table needs corresponding rows so a CSF examiner can resolve the institution-side trust-anchor controls to subcategories without reaching outside `CSF-2.0.md`.

**Status: Partial, dependent on Question 1 and Question 2 closure.** This is a roll-up partial; if Question 1 and Question 2 land their operational events, this partial closes by adding the four rows to the binding table.

### Question 5 — Does the 36-hour notification triage matrix handle institution-vs-regulator timing-overlap edge cases that arise during the dual-algorithm transitional period?

The 36-hour matrix at `incident-response-playbook.md` lines 286–298 and the named edge cases at lines 301–332 cover concurrent rollup, vendor-hosted determination split, federal-regulator-routing-per-charter, and CIRCIA-only-without-FFIEC. The federal-regulator-routing table at lines 311–319 is exactly the artifact a cyber-specialist examiner needs at scale — I can route a credit-union examination to NCUA without needing to reach for FFIEC organizational charts.

I tested one further edge case the matrix does not name: **dual-algorithm transitional-period institution-vs-regulator timing overlap.** During the multi-year post-quantum transitional period:

- The institution operates dual-algorithm posture (Ed25519 + post-quantum) with co-signed seals.
- The regulator may be operating a different posture transition timeline — some regulators may have transitioned their own verifier deployment to dispatch on both algorithms; others may still be operating Ed25519-only verifiers during their internal transition.
- A case (e) co-signed seal failure at the institution may produce different determinations depending on whether the regulator's verifier (which the institution does not directly control) validates only the un-broken algorithm versus dispatches on both algorithms.

The institution's IR Commander making the 36-hour clock-start decision under case (e) needs to know what the regulator's verifier will report when the institution's evidence packet reaches the regulator. If the regulator's verifier validates only the un-broken algorithm and reports PASS, the institution's "we have a co-signed seal failure" finding may not correspond to anything the regulator's verifier surfaces, creating a divergent narrative between the institution and the regulator at a load-bearing moment.

The IR playbook should add an edge-case section under "36-hour cyber-incident notification triage matrix" naming this:

- During the dual-algorithm transitional period, the institution's IR program SHOULD include in the regulator notification (a) the institution's verifier's per-algorithm validation results, (b) the verifier-binary version and SLSA-attestation digest the institution operated, and (c) a request that the regulator confirm which verifier-binary version the regulator is operating for the receiving examination.
- This pre-empts the divergent-narrative problem and lets both institution and regulator align on which seal validations are in scope before the determination becomes a finding.

**Status: Partial.** The standing edge cases are well-handled; the dual-algorithm institution-vs-regulator timing-overlap edge case is missing.

## Strengths I want to call out

The following are not findings; they are constructions I want to praise so the working group does not lose them in subsequent revision rounds.

- **The `key_fingerprint` check before MAC compute (spec §7 step 8).** This is the single best examination-workflow choice in the corpus. Cross-tenant configuration drift produces a precise named alert instead of a MAC-mismatch storm. The IR playbook's Scenario 7 four-case triage tree (rotation-in-flight / restored-backup / tenant-row-restored / unauthorized substitution) is the right examination companion.
- **The four-attack-to-CSF-subcategory table.** A working examination tool with named operational evidence per attack approach. Most spec authors do not produce this artifact.
- **The federal-regulator-routing table per institution charter.** Saves me from having to reach for organizational charts during the IR triage. Multi-charter holding companies are correctly handled.
- **The mirror-registry signature handling section.** The signature-preserving versus re-signing distinction with documented bridging is the right framing. The WORM-storage and 7-year retention requirements on the re-signing pattern are exactly the controls a cyber-specialist examiner reaches for.
- **The cosign-key recovery during active examination procedure.** The mid-engagement procedure is a real operational concern; the GPG-fallback validation as interim posture with post-recovery re-validation against the new cosign key is the right answer.
- **The dual-algorithm transitional-period verifier behavior at spec §7 step 11.** Cases (a) through (e) cover the post-quantum migration's verifier-side behavior with examiner-grade precision. Case (e) separating mechanism (verifier records both results) from policy (institution's IR program interprets) is correct.
- **The Adversary I institution-side reception procedure.** The five-step procedure is the right shape; my partial in Question 1 is about adding operational events for sample-testing, not redesigning the procedure.

## Roll-up

| # | Question | Finding |
|---|---|---|
| 1 | Adversary I institution-side reception procedure operational events | Partial |
| 2 | Cold-DR-key institution-side trust-anchor caching, dry-run-attestation consumption, health-check escalation IR scenario | Partial |
| 3 | Spec §7 step 11 case (e) co-signed-seal-failure mapping to a named IR scenario with clock-start determination tree | Partial |
| 4 | Operational-events binding table additions for the new institution-defined events from Q1 and Q2 | Partial (rolls up from Q1 and Q2) |
| 5 | 36-hour matrix dual-algorithm institution-vs-regulator timing-overlap edge case | Partial |

**Total: 0 gaps, 5 partials.**

The five partials are completeness items in the institution-side operational evidence and IR scenario coverage for trust-anchor-related events (regulator-held fingerprint rotation, cold-DR-key lifecycle consumption, dual-algorithm case (e) determination). They do not require changes to the chain construction, the verifier procedure, or the threat model. They require additions to the operational-events catalog (institution-defined per the `mirror.reconciliation_completed` precedent), additions to the IR playbook (one or two new scenarios plus one matrix row plus one edge-case paragraph), and additions to the supply-chain trust-anchor table (one row for the cold-DR public key fingerprint).

If the working group lands these five items, the cybersecurity examination workflow has end-to-end coverage from chain construction through verifier dispatch through IR determination through regulator notification, with named operational events at every load-bearing point. That is the bar I want to sign off at.
