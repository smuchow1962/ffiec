# Round-11 review — FFIEC Cybersecurity Specialist Examiner

## Persona

**Name.** Astrid Holm.
**Role.** FFIEC Cybersecurity Specialist Examiner with the FDIC. Eleven years on the cyber-specialist track; two years seconded to ENISA on the EU-US cybersecurity cooperation program. Active in the FFIEC working group on AI cybersecurity examination.
**Reading angle.** Threat model, IR playbook, NIST CSF 2.0 mapping fidelity, supply-chain SLSA L3 trust path, 36-hour cyber-incident notification triage with edge cases. I read documents the way an examiner reads control descriptions: I look for the control, the evidence the control produces, the operational event that lets me sample the evidence, and the path that lets me reach a finding when the control fails.

## What I read

- `spec/chain-of-custody-v1.md` — the v1.0-final spec text, with §10 operational requirements and §7 step 11 dual-algorithm transitional-period verifier behavior as my close-read targets.
- `docs/design/00-overview.md` — system overview and the §2 four-attack table.
- `docs/design/09-threat-model.md` — R1 through R13 residual-risk register; Adversaries A through I including the new Adversary I framing; the §2.8 dual-algorithm transitional period.
- `docs/regulator-pack/CSF-2.0.md` — the four-attack-to-subcategory table, the operational-events binding table, and the GOVERN-function alignment one-pager.
- `docs/incident-response-playbook.md` — Scenarios 1 through 10, the 36-hour triage matrix, and the named edge cases (concurrent rollup, vendor-hosted determination split, CIRCIA-only triggers without FFIEC trigger).
- `docs/cloud-hsm-guide.md` — the conformance bar and per-provider provisioning.
- `docs/supply-chain.md` — SLSA L3 with institutional consumption MUST-tier, mirror WORM + reconciliation, GPG recovery with dual-compromise, cosign recovery during examination.
- `docs/regulator-pack/ai-policy-alignment.md` — Treasury RMF, EU AI Act Article 14 + 26, DORA Articles 6 + 8, OCC supervisory posture.
- `docs/edge-and-federated-ai.md` — the deployment patterns I am increasingly being asked about.

I did not read any prior-round feedback files. This is a fresh-eyes review against the shipped corpus.

## Headline impression

This is the most rigorously evidenced AI-decision audit-trail spec I have read in the last twelve months. The construction I find most convincing is the per-event `key_fingerprint` check at verifier step 8, fired BEFORE any MAC compute. That single design choice means that at examination time I get a precise, named failure mode for cross-tenant configuration drift — `key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint` — instead of a MAC-mismatch storm I would have to investigate my way out of. Combined with the spec §10.1 weekly fingerprint reconciliation control and the IR playbook Scenario 7 triage tree (rotation-in-flight / restored-backup / tenant-row-restored / unauthorized substitution), the cybersecurity examination workflow has the precise alert / precise control / precise IR triage pipeline I want to see.

The four-attack-to-subcategory table in `CSF-2.0.md` is also exactly the artifact a cyber-specialist examiner needs. I can sit down with an institution's evidence repository, pick an attack approach, follow it to a CSF subcategory, and pull the specific operational event the institution should be able to produce. That is a working examination tool; most spec authors do not produce it.

The concerns I raise below are about completeness and edge cases, not about the design's foundation. The foundation is sound.

## Five questions before sign-off

### Question 1 — Is the four-attack-to-subcategory table complete enough that a cybersecurity examiner can sample-test all four attacks without reaching outside the document?

The CSF four-attack table at `CSF-2.0.md` lines 37–43 names a load-bearing CSF subcategory, a load-bearing operational evidence artifact, and the IR scenario for each of the four attacks. I tested this against my actual examination working pattern by walking each row.

**§2.1 Host-side chain tampering.** PR.DS-06 → `chain.verification_failure step=9` → Scenario 1. I can trace this end-to-end in the corpus. PASS.

**§2.2 Cross-tenant or cross-version key confusion.** PR.DS-06 + ID.AM-08 → `chain.verification_failure step=8` and `master.reconciliation_completed.fingerprint_unmatched_count` → Scenario 7. The dual operational evidence (verifier-emitted plus reconciliation-emitted) is exactly what I want to sample-test against the institution's log store. PASS.

**§2.3 Server-side or privileged-insider history rewrite.** PR.DS-06 + DE.CM-09 → `chain.verification_failure step=10/11` plus `seal.job_completed` for cross-correlation → Scenarios 2, 3. PASS.

**§2.4 Examination-time tooling subversion.** GV.SC + ID.RA-09 → reproducible-build log + cosign + GPG validation logs + `bundleverify` output + SLSA-verifier validation log → Scenario 3 (binary-signature variant). PASS, with a small request: the operational-events binding table at `CSF-2.0.md` lines 68–86 catalogs the spec §10.2 events but does NOT catalog `mirror.reconciliation_completed` or the institution-side SLSA-verifier validation log entry as binding to GV.SC-04. I can trace these from `supply-chain.md` (lines 178–183 for mirror reconciliation; lines 297–304 for SLSA institutional consumption), but the operational-events binding table is the artifact a cyber-specialist examiner reaches for first. **Status: Partial.** The four-attack table is complete; the operational-events binding table needs two additional rows to fully resolve §2.4 evidence to subcategory.

### Question 2 — Do the 36-hour triage edge cases cover concurrent rollup, vendor-hosted determination split, and CIRCIA-only triggers without FFIEC trigger?

I read the triage matrix at `incident-response-playbook.md` lines 286–298 and the edge-cases section at lines 301–319 carefully. The three edge cases the prompt asks about are all present and well-handled.

**Concurrent multi-scenario alerts (rollup rule).** Lines 303 explicitly addresses this: when Scenarios 1, 7, and 8 fire in the same 60-minute window, the IR Commander treats the constellation as one incident and starts the 36-hour clock at the earliest qualifying determination across the constellation. The playbook explicitly warns against both double-counting and under-counting. This is the right answer; it matches the FFIEC working group's emerging consensus on AI-incident rollup. PASS.

**Cross-tenant scope discovery during investigation (vendor-hosted).** Lines 305–307 split the responsibility correctly: vendor's IR team produces the alert and conducts the cross-tenant investigation; each institution's IR program receives the alert through the vendor's notification channel; the institution's clock starts when the institution determines (not when the vendor first alerted). The vendor's contractual SLA (typically 4–12 hours from vendor's determination) is a SEPARATE clock. This is exactly the framing I would use in a finding. PASS.

**CIRCIA-only triggers without FFIEC trigger.** Lines 309–319 with the per-scenario CIRCIA-vs-FFIEC table is the artifact I want. The Scenario 5 / Scenario 9 "Maybe / Likely" entries are honest — most CIRCIA-only situations arise from broader institutional context, not the chain-detected event in isolation, and the matrix communicates that. The recommendation that "the institution's IR program SHOULD have CIRCIA notification on its standard checklist for chain-detected events, even when the FFIEC clock has not started" is the right operational guidance. PASS.

**One concern across all three edge cases.** None of the three edge cases name the **NCUA notification path** for credit unions or the **state-banking-department notification path** for state-chartered institutions whose primary federal regulator is the FDIC or Fed but whose state regulator may also have a notification expectation. The 36-hour FFIEC computer-security incident notification rule applies federally; state notification frameworks may add additional clocks. For a credit union with material AI deployment, NCUA is the primary federal regulator and the 36-hour clock applies through NCUA, not OCC/FDIC/Fed. The IR playbook's recipient column at lines 351–357 says "Primary federal regulator" generically, which is correct, but the edge-cases section should explicitly name "Different primary federal regulator per institution charter (OCC for national banks, FDIC for state non-member banks, Fed for state member banks, NCUA for credit unions)" so a credit-union examiner reading the playbook does not mistakenly route the notification through the wrong agency. **Status: Partial.** The federal-regulator-routing edge case is missing.

### Question 3 — Does the spec §7 step 11 dual-algorithm transitional-period verifier behavior cover the failure modes a cyber-specialist examiner will encounter during the post-quantum migration?

The spec §7 step 11 dual-algorithm dispatch at lines 321–326 names four cases: (a) both signatures present and both valid (PASS), (b) single-algorithm signature on a seal during dual-algorithm posture (PASS-WITH-ANOMALY with `partial-coverage seal` reason), (c) seal carrying an algorithm not on the institution's declared algorithm-posture list (FAIL under `--strict`, PASS-WITH-ANOMALY otherwise), (d) single-algorithm posture default for v1.0 (reduces to single-algorithm verification).

The four-case dispatch is correct on the chain-integrity dimension. The behavior I would want to see explicitly added: **case (e) — both signatures present, one valid and one invalid**. This is the case where the institution operates dual-algorithm posture, both algorithms produced a signature, and the verifier finds that one signature validates but the other does not. The spec text does not name this case. The natural reading — that the seal is integrity-bearing under the valid-signature algorithm and the invalid signature is a control-completeness finding — is plausible but not normative. During the multi-year transitional period, the institution's HSM may produce a signature under the post-quantum algorithm that fails post-quantum-key validation due to a vendor regression (a known class of issue with new HSM firmware shipping post-quantum support); the verifier needs a defined behavior. My recommendation is to add case (e) explicitly: "Both signatures present, one valid and one invalid. Under `--strict`: FAIL with `algorithm signature mismatch on dual-signed seal: one of two signatures failed validation`. Under non-strict: PASS-WITH-ANOMALY with the same reason. The institution investigates as a control-completeness finding; the integrity-bearing claim is reduced to the validating algorithm's strength." **Status: Partial.** Case (e) is missing.

I would also like the spec to name the **examiner's working-paper recommendation** for dual-algorithm posture during a live transition — specifically, that an examiner sampling a tenant-day's seal during the transition records BOTH algorithm validation results in the working papers, NOT just the PASS/FAIL roll-up. The reason: a year after the transition completes, when the post-quantum algorithm is the only one in use and the institution retires the Ed25519 algorithm from its declared posture list, an auditor reviewing the historical examination working paper needs to see that BOTH algorithms validated at the time the seal was created, not just that the seal "passed." This is a working-paper hygiene point that a fresh examiner reading the spec for the first time will not derive on their own. The dual-algorithm spec text should reference an examiner-working-paper subsection in `examiner-quickstart.md` or `regulator-pack/sample-report.md` that names the dual-algorithm working-paper convention. **Status: Partial.** The working-paper convention is unstated.

### Question 4 — Does SLSA L3 institutional consumption with `slsa-verifier` close the supply-chain trust path I would challenge as a cyber-specialist examiner?

The SLSA L3 institutional consumption MUST-tier at `supply-chain.md` lines 297–304 is the right framing. The four properties — (1) institution validates against `slsa-verifier` BEFORE deploying, (2) SLSA-verification log retained for binary's deployment lifetime (typically 7 years), (3) `slsa-verifier` itself signed and trust-anchored alongside the verifier, (4) `slsa-verifier` is co-equal trust artifact alongside cosign and GPG manifest — directly close the gap a GV.SC-04 examiner would challenge. The framing that "an institution that publishes the L3 claim without operating `slsa-verifier` consumption has a gap a GV.SC-04 examiner will challenge" is exactly the language I would use in a finding.

The mirror WORM retention + reconciliation control at `supply-chain.md` lines 178–183 is similarly sound. The two retention properties (7-year minimum, WORM storage) and the continuous-monitoring control (`mirror.reconciliation_completed`) close the bridge-invariant gap. The framing that "a re-signing mirror that silently stops validating project signatures … would produce mirror-side re-signed images with no project-side validation evidence — and the institution's deployment pipeline, validating only the institution's signature, would not detect this" is the precise finding language I would use.

**Status: Pass on the construction.** Two clarification requests:

First, the SLSA L3 institutional consumption section at lines 297–304 names `slsa-verifier` as the canonical tool. SLSA-aware tooling has multiple implementations (the SLSA project's `slsa-verifier`, in-toto's `attestation-verifier`, `cosign verify-attestation`, vendor-specific tools from cloud providers). The MUST-tier framing currently reads "validate against `slsa-verifier` (or equivalent SLSA-aware tooling)" which is fine, but a cyber-specialist examiner sampling the institution's SLSA-verification log needs to know which tool the institution chose AND whether that tool's signature trust path is itself documented. For an institution running `cosign verify-attestation` instead of `slsa-verifier`, the trust anchor is `cosign`'s public key, which the institution already manages. For an institution running a cloud-vendor-specific tool, the trust anchor is the cloud vendor's tool-signing identity, which is a NEW trust anchor the institution needs to manage. The supply-chain doc should explicitly require the institution to document which SLSA-aware tool is in use AND the trust anchor for that tool. **Status: Partial.** The tool-choice documentation requirement is implicit, not normative.

Second, the dual-compromise edge case at `supply-chain.md` lines 97–104 names the "cold-disaster-recovery key" maintained by yet another separated small group as the emergency response. This is a sound design. The cyber-specialist examiner question I would ask: what is the **cold-disaster-recovery key's own rotation cadence and key-holder lifecycle**? An offline-held cold-disaster-recovery key that has not been rotated in 10 years, held by a group whose composition has churned over that period, is not a credible emergency response. The supply-chain doc should name (a) the cold-disaster-recovery key's rotation cadence (offline keys still need rotation; NIST SP 800-57 guidance applies), (b) the key-holder-lifecycle procedure (who is added, who is removed, how that procedure stays auditable), and (c) the institution-side validation that the cold-disaster-recovery key remains operational (a "cold-key recovery dry-run" cadence — perhaps every 2 years — that exercises the recovery path without an actual disaster). **Status: Partial.** The cold-disaster-recovery key's own lifecycle is undocumented.

### Question 5 — Does the threat model's R9–R13 register, the dual-algorithm framing, the regulator-held fingerprint rotation pointer, and the new Adversary I framing close the residual-risk surface I would want closed?

I walked R9 through R13 against the spec text and the IR playbook. Each is closed by a specific spec section or operational control:

- **R9 (cross-tenant chain confusion via HKDF input collision).** Closed by spec §4.1 inviolate property #1 (tenant_id binding into HKDF info) plus per-entry `key_fingerprint` check at verifier step 8. I can verify this by inspection: two tenants whose IKMs are swapped derive verifiably different session keys, AND the fingerprint check rejects the wrong-key configuration before any MAC compute. PASS.

- **R10 (future-maintainer relaxation of structural prev_hash check).** Closed by spec §4.1 inviolate property #8 (verifier feeds `expected_prev_hash` into MAC recompute, NOT `entry.prev_hash`). The "latent footgun" framing is exactly right; this is the kind of subtle property that survives spec-text changes only if it is named as inviolate. PASS.

- **R11 (IKM-registry premature retirement).** Closed by spec §10.9 (IKM retention coupling) + audit-procedures P-5 + `master_key.retired` operational event + IR Scenario 8 (unknown_key_version). The retention-coupling normative requirement that "a request to retire an IKM whose `key_version` is still referenced by retained chain entries MUST require an explicit override and MUST be logged" is the precise control I would sample-test. PASS.

- **R12 (edge-device physical compromise — secure-enclave attestation defeated).** Mitigation deferred to v1.1; v1.0 documents the residual and operates compensating controls (Pattern A per `edge-and-federated-ai.md`: per-device IKM in TPM/secure-enclave; reconciliation cadence baseline tuned per fleet). The deferral framing is honest. The compensating controls are operationally credible for current edge-AI deployments. PASS as a residual; the v1.1 candidate work is named.

- **R13 (algorithm-confusion attack at the seal layer — post-quantum coexistence).** Closed by spec §4.3 sign_payload extension (algorithm bound into signed payload). The closure is correct: cross-algorithm replay fails at signature verification regardless of public-key resolution. PASS.

**Adversary I framing.** The new examiner-laptop-hygiene-as-regulator-side framing at `09-threat-model.md` lines 199 is exactly right — a cyber-specialist examiner reads this and immediately understands that examiner-laptop hygiene is the FFIEC's own examiner-IT program, not an institution-side control to be evaluated under SOC 2 CC6. The explicit guidance that "institution-side reviewers (SOC engagement partners) should NOT misread this residual as an institution-owned control gap" is the kind of framing that prevents a SOC engagement partner from generating a false finding. PASS.

**The regulator-held fingerprint rotation pointer.** At `09-threat-model.md` line 197 the threat model now names the regulator-held fingerprint as the trust anchor that makes the institution-published public key trustable, AND notes that "rotating it is a regulator-side procedure (referenced in `regulator-pack/regulator-procedures.md` — v1.1 candidate stub; until the procedure is published, the trust-anchor rotation event is recognized as an unbounded residual the institution cannot control institution-side)." The honest framing that this is unbounded and out of institution control is correct.

**Status: Partial.** Two concerns on the regulator-held fingerprint rotation pointer:

First, the v1.1 candidate stub framing means the procedure does not exist yet. From the FFIEC working-group seat, I will tell you that this procedure needs to exist before institutions can credibly defend their dual-trust-path (cosign + regulator-held-fingerprint) claim to an examiner. The current state is "the institution's defense rests on a regulator-side procedure that is not yet published." A cyber-specialist examiner reading this will note that the dual-trust-path claim is conditional on a future regulator publication; that is a real residual the threat model should name as such. **Status: Partial — the dependence on an unpublished regulator-side procedure should be named as a residual that bounds the dual-trust-path strength claim until the procedure is published.**

Second, even when the procedure is published, the regulator-side fingerprint rotation event has institution-side implications: the institution must update its trust anchors to the new regulator-held fingerprint within some bounded window, AND the institution must validate that the rotation notice is itself authentic (not a phishing-style attempt). The threat model should name the institution-side reception-and-validation procedure for the regulator-held-fingerprint rotation event, even if the rotation itself is regulator-side. **Status: Partial — the institution-side reception procedure for regulator-side fingerprint rotation is unstated.**

## Roll-up status

| Area | Status |
|---|---|
| CSF four-attack-to-subcategory table | **Partial** — operational-events binding table missing rows for `mirror.reconciliation_completed` and SLSA-verifier validation log binding to GV.SC-04 (Question 1) |
| 36-hour triage edge cases (concurrent rollup) | **Pass** (Question 2) |
| 36-hour triage edge cases (vendor-hosted determination split) | **Pass** (Question 2) |
| 36-hour triage edge cases (CIRCIA-only) | **Pass** (Question 2) |
| 36-hour triage edge cases (federal-regulator routing per institution charter) | **Partial** — credit-union NCUA path and state-banking-department paths should be named (Question 2) |
| Spec §7 step 11 dual-algorithm dispatch (cases a–d) | **Pass** (Question 3) |
| Spec §7 step 11 dual-algorithm dispatch (case e — one valid, one invalid) | **Partial** — case missing (Question 3) |
| Dual-algorithm examiner working-paper convention | **Partial** — convention unstated (Question 3) |
| SLSA L3 institutional consumption (`slsa-verifier`) construction | **Pass** (Question 4) |
| SLSA-aware tool choice documentation | **Partial** — tool-choice documentation requirement is implicit (Question 4) |
| Mirror WORM retention + reconciliation control | **Pass** (Question 4) |
| Cold-disaster-recovery key rotation and lifecycle | **Partial** — cold-key's own lifecycle is undocumented (Question 4) |
| Adversary I framing (regulator-side examiner-laptop hygiene) | **Pass** (Question 5) |
| R9, R10, R11, R12, R13 residual-risk closures | **Pass** (Question 5) |
| Regulator-held fingerprint rotation pointer (procedure dependency) | **Partial** — dependence on unpublished regulator-side procedure should be named as bounding residual (Question 5) |
| Regulator-held fingerprint rotation pointer (institution-side reception procedure) | **Partial** — institution-side reception-and-validation procedure unstated (Question 5) |

**Overall.** 9 Pass, 7 Partial, 0 Fail. The corpus is at signing-distance for a cybersecurity examination. The seven partials are not foundational defects; they are completeness items that a cyber-specialist examiner will reach for in the field. Closing them moves the spec from "I can examine against this" to "I can examine against this without reaching for clarification."

Stopping criterion (0 gaps + 0 partials) not met. Recommend round-12 closes the seven partials enumerated above; if closed, this reviewer's sign-off follows.
