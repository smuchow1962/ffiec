# Round 10 — FFIEC Cybersecurity Specialist Examiner

> **Reviewer.** Liang Zhao, FFIEC Cybersecurity Specialist Examiner with the OCC. Fourteen years on the cyber-specialist track. Recent FRB cross-rotation. NIST CSF 2.0, FFIEC computer-security incident notification rule (36-hour), CIRCIA, SLSA framework as my standing reference set.
>
> **Reading angle.** Threat model, IR playbook, cyber-specialist lens. CSF 2.0 mapping fidelity, supply-chain trust path, 36-hour notification triage. Honest assessment from a fresh read of the v1.0-rework corpus — I have no insider view of pre-rework iterations.

---

## What I read (in order)

1. `spec/chain-of-custody-v1.md` — full pass with focus on §10 operational requirements
2. `docs/design/00-overview.md` — system overview and the §2 four-attack catalog
3. `docs/design/09-threat-model.md` — adversary set A–I, residual-risk register R1–R13, dual-algorithm transitional period
4. `docs/regulator-pack/CSF-2.0.md` — operational-events binding table and GOVERN-function one-pager
5. `docs/incident-response-playbook.md` — full scenario set 1–10 and the 36-hour triage matrix
6. `docs/cloud-hsm-guide.md`
7. `docs/supply-chain.md` — SLSA L3 claim, cosign-key recovery, GPG rotation, mirror-registry, spec/corpus integrity binding
8. `docs/regulator-pack/ai-policy-alignment.md` — Treasury RMF, EU AI Act Art. 26, DORA Arts. 6 + 8, OCC supervisory posture
9. `docs/edge-and-federated-ai.md`

---

## Headline impression

The corpus reads like a standard that has already absorbed examiner feedback once. The four-attack catalog in `00-overview.md` §2 is the clearest piece of design framing I have seen on an AI-decision audit-trail standard — it names the attacker, names the goal, names the primitive that defends, and refuses to wave at the composition. The IR playbook's expansion into Scenarios 7–10 closes the gap I would have flagged immediately on a CSF-aligned cyber exam: "what does the institution do when the verifier emits the new step-7 / step-8 failure modes that the rework introduced." The 36-hour triage matrix per scenario is the artifact I want every institution adopting this spec to have on the wall in the SOC.

The supply-chain doc does the work I expect of a SLSA L3 claimant: it names which L3 requirements are met and how, separates cosign and GPG roles to non-overlapping holders, articulates GPG rotation cadence with NIST SP 800-57 grounding, documents both mirror-registry patterns with a non-conformant variant called out, and binds the spec text and conformance corpus into the same trust path as the binary. The corpus-rebuild defense against a subverted-corpus attack closes the residual concern I had about "what if the project's source-of-truth is compromised."

The threat model's R9–R13 additions and the named Adversary I are correctly placed against the rework primitives. R13 (algorithm-confusion at the seal layer) is bound at the right place — inside `sign_payload` per spec §4.3 — and the dual-algorithm transitional-period prose in §2.8 correctly states the integrity claim under both algorithm assumptions. The CSF 2.0 doc's operational-events binding table is the single artifact a cyber-specialist examiner needs to construct an evidence-to-subcategory matrix at scale; the GOVERN-function one-pager directly maps to the CSF 2.0 elevation of GOVERN to a top-level function.

There are five places I want to push on, articulated as questions below. None is a structural defect; all are convergence-tightening on a corpus that is already substantially close.

---

## Questions

### Question 1 — CSF 2.0 mapping completeness against the four attacks

The CSF mapping table addresses primitives and operational events but does not explicitly walk the **four attack approaches** from `00-overview.md` §2 against the function/category set. A cyber-specialist examiner constructing an evidence matrix wants to ask: for each of the four attacks (host-side tampering, cross-tenant key confusion, server-side history rewrite, examination-time tooling subversion), which CSF subcategories are the load-bearing detection-and-response evidence?

The current table maps the chain-as-control to PR.DS-06 as headline and threads everything else as composition. That is true and correct, but it leaves the cyber-specialist to do the cross-walk from "attack approach" to "subcategory evidence" themselves. A four-row table — one row per attack — naming the load-bearing CSF subcategory and the load-bearing operational event would let the cyber-specialist examiner sample-test against the institution's evidence store directly, without re-deriving the mapping per exam.

Worked sketch of what I have in mind:

| Attack approach (00-overview.md §2) | Load-bearing CSF subcategory | Load-bearing operational evidence | IR scenario |
|---|---|---|---|
| §2.1 Host-side tampering | PR.DS-06 (information integrity) | `chain.verification_failure step=9` | Scenario 1 |
| §2.2 Cross-tenant / cross-version key confusion | PR.DS-06 + ID.AM-08 | `chain.verification_failure step=8` and `master.reconciliation_completed.fingerprint_unmatched_count` | Scenario 7 |
| §2.3 Server-side history rewrite | PR.DS-06 + DE.CM-09 | `chain.verification_failure` at Merkle/signature steps; verifier `merkle root mismatch` and `signature verification failed` | Scenarios 2, 3 |
| §2.4 Examination-time tooling subversion | GV.SC + ID.RA-09 | Reproducible-build log + cosign/GPG validation; `bundleverify` output | Scenario 3 (binary-signature variant per supply-chain.md) |

The information is already in the corpus — it just isn't indexed by attack approach. Adding it makes the CSF mapping doc the single reference a cyber-specialist examiner uses, not a cross-walking exercise.

### Question 2 — 36-hour clock-start triggers: edge cases the matrix doesn't yet name

The triage matrix is the right artifact at the right granularity, and the per-scenario clock-start triggers are correctly framed around *determination* (not alert receipt) per the rule. Three edge cases I want to surface for explicit naming:

**2a. Concurrent multi-scenario alerts.** If an institution receives a Scenario 1 (chain hash mismatch) alert AND a Scenario 7 (key_fingerprint mismatch) alert AND a Scenario 8 (unknown key_version) alert in the same 60-minute window, the institution may be looking at a single root-cause event that triggered three detection paths. The matrix as written treats each scenario's clock independently. In practice, the institution's IR program should treat the constellation as one incident and start the clock at the earliest qualifying determination. Adding a one-paragraph "concurrent-alert rollup rule" to the matrix preface would prevent both (a) double-counting incidents in the regulator notification and (b) under-counting where the constellation is the signal even though no individual scenario crossed its threshold.

**2b. Cross-tenant scope discovery during investigation.** Scenarios 4 and 7 may discover during investigation that the suspected compromise spans more than one tenant in a multi-tenant deployment (vendor-hosted topology, shared-cloud-HSM at community-bank tier per `00-overview.md` §6.5). The 36-hour clock applies per institution; in a vendor-hosted deployment one investigation may produce N institution-side determinations on different timelines. The playbook should name the institution-vs-vendor responsibility split for the determination event in vendor-hosted topology, since the institution's clock starts when the institution determines, not when the vendor first alerted.

**2c. CIRCIA-only triggers without FFIEC trigger.** The "CIRCIA and additional notification frameworks" section names CIRCIA's 72-hour path. Some chain-detected events (e.g., Scenario 5 sealing-delay associated with a broader cyber incident at the institution; Scenario 9 truncation pattern across hosts) may rise to CIRCIA's "substantial cyber incident" threshold without rising to the FFIEC rule's "computer-security incident affecting safety/soundness" threshold. The matrix as written triggers on the FFIEC rule; I would want the matrix to also have a CIRCIA-trigger column so the institution's IR program does not miss the 72-hour CISA path because the FFIEC clock did not start.

### Question 3 — Dual-algorithm transitional period: verifier-side strict-mode behavior during transition

§2.8 of the threat model articulates the dual-algorithm transitional period correctly: the verifier dispatches on the per-seal `algorithm` field, both signatures coexist during the transition, the integrity claim is "sound under both algorithm assumptions until one is broken." Spec §4.3 binds `algorithm` into `sign_payload`, closing the algorithm-confusion attack class (R13).

The piece I want explicit handling for: **what is the verifier's `--strict` mode behavior during the transition?** Concretely, if an institution operates under v1.x with both Ed25519 and a post-quantum algorithm available, and a tenant-day's seal arrives signed *only* under Ed25519 (the institution had not yet enabled the PQ co-signing), is `--strict` PASS or PASS-WITH-ANOMALY? I would expect `--strict` to PASS the seal if the algorithm is still on the conformant list and the signature verifies, with PASS-WITH-ANOMALY reserved for "the institution declared dual-algorithm posture but the seal is single-algorithm." This is a per-institution control-description question more than a spec question, but the spec's §7 step 11 verifier procedure should name dual-algorithm dispatch explicitly so an examiner reading the spec knows what the verifier does. Currently §7 step 11 handles single-algorithm dispatch crisply; the transitional case is handled in `09-threat-model.md` §2.8 prose but not in the verifier procedure.

A one-paragraph addition to §7 step 11 along the lines of: "When the institution operates dual-algorithm posture, the seal record's `algorithm` is the dispatch key and the verifier resolves the public key for that algorithm. A seal carrying an algorithm not on the institution's declared algorithm-posture list is reported as `algorithm not on institution's declared posture list` and is `--strict` FAIL. A seal carrying a single algorithm during dual-algorithm posture is `--strict` PASS-WITH-ANOMALY (the institution's declared posture is incomplete on this seal-day)." That gives the cyber-specialist examiner a deterministic verifier outcome to test against during the multi-year transition.

### Question 4 — Supply-chain trust path: SLSA L3 to L4 transition and mirror-registry telemetry

The SLSA L3 claim is well-supported. The "Level 4 is a v1.1 candidate" framing is the right shape for now — L4's two-party review at the build service is a project-side governance evolution, not a chain-construction defect. Two follow-ups for the cyber-specialist examiner reading this:

**4a. SLSA L3 evidence the institution presents.** The doc says the cybersecurity examiner uses the SLSA L3 claim as evidence for GV.SC-04 and ID.RA-09. The substrate is the published `verifier.intoto.jsonl`. What I want explicitly named: the institution-side workflow for consuming the provenance attestation. Specifically — does the institution run `slsa-verifier` (or equivalent) at the same point in the deployment pipeline as cosign signature validation, and is the `slsa-verifier` exit code a gating control on deployment? The "What examiners verify" list at the bottom of `supply-chain.md` lists cosign and GPG validation but does not list provenance-attestation validation as a checked item. Adding a line — "the institution validates the SLSA provenance attestation against `slsa-verifier` (or equivalent) at deployment time and archives the validation log" — would close the cyber-specialist's evidence ask.

**4b. Mirror-registry telemetry.** The re-signing mirror pattern documents the trust-path bridge correctly: the mirror's audit log records the project-side cosign validation at pull time, and a missing log entry routes to Scenario 3. What is not yet named: the institution's continuous-monitoring control that asserts the audit log is being written. A re-signing mirror that silently stops validating project signatures (operator misconfig at the mirror, registry-side bug, vendor regression) would produce mirror-side re-signed images with no project-side validation evidence — and the institution's deployment pipeline, validating only the institution's signature, would not detect this. The DE.CM-09 evidence for the mirror should be a periodic (daily or per-pull) reconciliation between the mirror's pulled-image set and the audit-log-validation set, with mismatch firing an operational alert. This is the same control pattern as the spec §10.1 fingerprint reconciliation, applied to the mirror's signature-validation discipline.

### Question 5 — Adversary I and the regulator-held trust anchor

Adversary I in `09-threat-model.md` §2.9 is named correctly and the residual risk is articulated honestly: a simultaneous compromise of cosign + GPG + the regulator-held public-key fingerprint would silently pass forged data. The defense composition (three independent compromises must align) is the right shape.

The piece I want explicit: **what is the regulator-side procedure for rotating the regulator-held fingerprint?** The doc says "rotating it is a regulator-side procedure documented separately." I want that procedure to either be in the corpus (as a stub or pointer) or to be named as a v1.1 candidate with the gap acknowledged in the residual-risk register. Without it, the trust-anchor rotation is "documented separately, somewhere"; for the cyber-specialist examiner working a federal rotation event (institution merger, regulator re-organization, fingerprint key-decay event), the absence of a referenced procedure is a real ambiguity.

A one-line addition to R3 or R5 in the residual-risk register: "the regulator-held public-key fingerprint rotation procedure is regulator-side and is referenced in `regulator-pack/regulator-procedures.md` (v1.1 candidate)" — even if the referenced doc does not yet exist — would tell the examiner the gap is recognized and bounded, not silently inherited.

---

## Findings

### Strengths

- **Four-attack catalog is the cleanest design framing I have read on this class of system.** Naming the attacker, the goal, the primitive that defends, and refusing to wave at composition is the right posture for FFIEC consumption.
- **IR Scenarios 7–10 close the rework's detection-surface expansion.** Scenario 7's triage tree (rotation-in-flight, restored-backup, tenant-row-restored, suspected substitution) correctly avoids the "treat every fingerprint mismatch as a 36-hour event" failure mode that would burn out IR teams. Scenario 8's three-cause triage is similarly well-shaped.
- **36-hour triage matrix is the most operationally usable artifact in the corpus.** Per-scenario clock-start triggers framed around determination (not alert receipt) match the rule's actual text. The matrix is the document I would want every adopting institution's SOC to have laminated.
- **R9–R13 residual-risk additions are correctly bound to spec primitives.** R13 in particular — algorithm-confusion at the seal layer — is closed at the right place (`sign_payload` extension) and the dual-algorithm transitional prose is the right shape for the multi-year migration window.
- **Adversary I named for first-class catalog completeness.** Mirroring the §2.4 attack approach in the adversary register removes the asymmetry between the four attacks and the eight prior adversaries.
- **CSF 2.0 operational-events binding table.** This is the artifact a cyber-specialist examiner uses to construct an evidence-to-subcategory matrix at scale. Headline mapping correctly identifies PR.DS-06 as the chain's primary subcategory; surrounding subcategories correctly framed as composition.
- **GOVERN-function one-pager.** Walks GV.OC, GV.RM, GV.RR, GV.PO, GV.OV, GV.SC against the chain. The "institution copies this section into its own GV-alignment policy document" framing is the right posture — adopt, don't re-invent.
- **SLSA L3 claim is properly substantiated.** Each L3 requirement has a "how the pipeline satisfies it" cell; the L4 gap is honestly named. The provenance attestation publishing is concrete (`verifier.intoto.jsonl`), not aspirational.
- **Cosign-key recovery procedure is operationally specific.** 24/48/7-day timeline; new-key authentication chain via GPG-signed introduction notice; old-binary historical verifiability preserved.
- **GPG rotation cadence is grounded.** 36-month default per NIST SP 800-57; 12-month subkey rotation; 7-day emergency rotation; transition statement signed under both old and new keys.
- **Mirror-registry handling.** Both signature-preserving and re-signing patterns named with the bridge-documentation requirement on the re-signing variant. Non-conformant variant explicitly called out.
- **Spec / corpus integrity binding.** Cosign + GPG over the spec PDF and the conformance corpus tarball, with the corpus-rebuild defense recommended for examination-grade adopters. Closes the subverted-corpus attack class against an attacker compromising the project's source-of-truth.
- **EU AI Act Article 26 + DORA Articles 6 + 8 mapping.** Specific article-level composition for deployer obligations, ICT risk-management framework, and ICT-supported business function identification. The MRM director defending chain adoption to a European regulator has the language ready.
- **OCC supervisory posture framing.** "Adopted to position the institution for current and emerging OCC expectations the OCC has signalled through supervisory letters to peer institutions" is the correct posture — chain-adoption is forward-looking but not over-claiming an explicit OCC mandate that does not yet exist.

### Gaps (1)

**G-CSE-1.** CSF mapping doc does not index CSF subcategories by the `00-overview.md` §2 four-attack catalog. The information is present in the corpus but requires the cyber-specialist examiner to cross-walk from attack approach to subcategory evidence rather than reading it directly. Per Question 1.

### Partials (4)

**P-CSE-1.** 36-hour triage matrix does not name three operational edge cases: concurrent multi-scenario alerts (rollup rule), cross-tenant scope discovery in vendor-hosted topology (institution-vs-vendor determination split), and CIRCIA-only triggers without FFIEC trigger (CISA 72-hour path). Per Question 2.

**P-CSE-2.** Spec §7 step 11 verifier procedure handles single-algorithm dispatch crisply but does not name `--strict` mode behavior under the dual-algorithm transitional period. Threat-model §2.8 has the prose but the verifier procedure is the operational artifact and should carry the deterministic outcome. Per Question 3.

**P-CSE-3.** Supply-chain "What examiners verify" list does not include SLSA provenance-attestation validation as a deployment-gating control item. The L3 claim is supported by published `verifier.intoto.jsonl` but the institution-side consumption is implicit. Mirror-registry section also lacks a continuous-monitoring control on the mirror's signature-validation discipline (re-signing pattern). Per Question 4.

**P-CSE-4.** Regulator-held public-key fingerprint rotation procedure is referenced as "documented separately" but no pointer or stub exists. Adversary I residual-risk acceptance is honest but the trust-anchor rotation event is unbounded. Per Question 5.

### Status

| Category | Count |
|---|---|
| Strengths | 14 |
| Gaps | 1 |
| Partials | 4 |

**Overall posture.** The corpus is substantially close to a CSF-2.0-aligned cyber-specialist examiner's expectations. The four observations above are convergence-tightening, not structural redesign. The single gap (G-CSE-1) is an indexing add to an existing doc; the four partials each resolve with one to three paragraphs of prose against an existing artifact. None of the five concerns blocks examination-grade use of the corpus today; resolving them would let the cyber-specialist examiner work the spec at production cadence without cross-walking documents during the exam.

---

## Roll-up

| ID | Type | Title | Resolves with |
|---|---|---|---|
| G-CSE-1 | Gap | CSF mapping not indexed by four-attack catalog | Add four-row table to `regulator-pack/CSF-2.0.md` mapping each `00-overview.md` §2 attack to load-bearing CSF subcategory + operational event + IR scenario |
| P-CSE-1 | Partial | 36-hour triage matrix missing concurrent-alert rollup, vendor-hosted determination split, CIRCIA-only column | Three paragraph additions to `incident-response-playbook.md` matrix preface; CIRCIA column added to matrix |
| P-CSE-2 | Partial | Verifier `--strict` behavior under dual-algorithm transitional period not in spec §7 step 11 | One paragraph addition to spec §7 step 11 naming dual-algorithm dispatch and the PASS / PASS-WITH-ANOMALY / FAIL outcomes |
| P-CSE-3 | Partial | SLSA provenance validation and mirror-registry telemetry not in examiner-verifies list / monitoring controls | Add `slsa-verifier` validation and archive to "What examiners verify"; add mirror audit-log reconciliation control to mirror-registry section |
| P-CSE-4 | Partial | Regulator-held fingerprint rotation procedure unreferenced | Add pointer (or v1.1-candidate stub) in residual-risk register against R3 / R5 / Adversary I |

Convergence target for round 11: zero gaps and zero partials, achieved by the five touch-ups above against the existing artifacts.
