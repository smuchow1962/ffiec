# Herald vendor-conformance round 5 — response

> **What this doc is.** Point-by-point response to the Herald vendor-conformance round 5 review (Q42-Q48), tied to the v1.0a amendment wave. The review carried 6 Partials and 1 Nit against the post-amendment spec, all clustered around implementer-precision, informative-completeness, and a residual court-facing item carried from round 4. None affect the cryptographic-integrity claim. This response cross-walks each finding to where it closes — either spec-text language already landed, a companion document where the closing language belongs, or a forward-scope deferral with rationale.

> **Audience.** FFIEC working group, Herald vendor implementation team, and any reviewer auditing the closure trail for the post-2026-05-07 amendment wave.

> **Round 4 closure baseline.** Round 4 had seven open findings (Q35-Q41). Six closed cleanly in the amendment wave (Q36 region definition, Q37/Q38 Pattern A operational specificity, Q39 trusted-time forward-scope acknowledgement, Q40 `ffiec.chain.late_binding` trust posture). Q41 carries forward as Q48 in round 5 with the same closing language proposed.

---

## Cross-walk summary

| Finding | Status | Severity | Close-out path | Where it lands |
|---|---|---|---|---|
| Q42 — RNG-source verifier signal | Partial | Implementer-precision | Companion-doc note in `audit-procedures-cluster-d-addenda.md` (P-54); spec §10.6.1 wording stays informative | Audit-procedure addenda |
| Q43 — Wrapped-import IKM workflow | Partial | Implementer-precision | Spec §10.6.1 informative addition (text in §3 below); compatible with existing §10.5 HSM custody | Spec amendment candidate (informative) |
| Q44 — Session-key-only compromise non-claim | Partial | Informative-completeness | Spec §1.2 informative paragraph (text in §4 below) | Spec amendment candidate (informative) |
| Q45 — Retention-vs-deprecation interaction | Partial | Informative-completeness | Spec §1.2 informative paragraph (text in §5 below) | Spec amendment candidate (informative) |
| Q46 — Linux LTS kernel wording | Nit | Editorial | Spec §10.6.1 wording fix (text in §6 below) | Spec amendment candidate (editorial) |
| Q47 — Verifier amendment-awareness CC8.1 | Partial | Operational | Spec §10.12 or §10.5 addition (text in §7 below); audit procedure P-55 | Spec amendment candidate + audit-procedure addenda |
| Q48 — §1.1 Daubert SDK-process compromise | Partial (carried) | Court-facing | Spec §1.1 informative extension (text in §8 below) | Spec amendment candidate (highest priority) |

All seven findings are addressable in editorial pass. None require new normative cryptographic surface. Q48 is highest priority for the next working-group cycle because it is defensive of the institution's existing courtroom posture under FRE 702 and Daubert and remained open from round 4.

---

## §1. Q42 — RNG-source verifier signal

### What the reviewer asked

The `master_key.generated` operational event records the RNG type. The verifier never sees operational events. An institution claiming `"hsm.cloudhsm-classic"` in CC8.1 while their actual deployment used `"os.linux-urandom"` produces an indistinguishable chain because every per-event MAC is correct under whichever IKM was used. The discrepancy surfaces only at SOC engagement when the auditor inspects the operational-event log against CC8.1. There is no verifier-side working-paper line that the institution's MRM committee or examiner can read against the institution's declared posture.

### Response

The chain holds under any unpredictable IKM regardless of RNG source. Promoting RNG-source declaration to a verifier-output integrity check would be incorrect — the chain remains valid. The closing path is to make the RNG-source decision verifier-visible without elevating it to a PASS/FAIL signal.

The institution's tenant-public-key registry (the existing public-key publication surface per spec §10.5) carries an `rng_source` attribute parallel to the existing public-key fields. The verifier prints the registry's declared `rng_source` in its working-paper output. The output line is informational; it does not affect PASS/FAIL but it gives the SOC auditor and the FFIEC examiner a visible cross-check against the institution's CC8.1 claim.

This closes the audit-visibility gap without disturbing the chain's integrity claim. The audit-procedure addition (P-54) names the cross-check as a sample procedure: audit pulls the registry's declared `rng_source` for the audit period, pulls the institution's CC8.1 control description for the same period, and confirms they match. Mismatches escalate to the AI Governance Committee for investigation.

### Closing language for spec §10.6.1 (informative)

> The institution's tenant-public-key registry SHOULD declare the `rng_source` attribute alongside the public-key fields. Conformant verifiers print the registry's declared `rng_source` in their working-paper output as an informational line. The verifier output is not affected by the declared `rng_source` value; the chain holds under any RNG source meeting the institution's CC8.1 commitment. Cross-check between the registry-declared `rng_source` and the institution's CC8.1 declaration is an audit procedure (P-54); inconsistency is an audit finding, not a verifier failure.

---

## §2. Q43 — Wrapped-import IKM workflow

### What the reviewer asked

AWS KMS, Azure Key Vault Managed HSM, and Google Cloud KMS expose customer-supplied-key import workflows where the IKM is generated outside the HSM and flows in via a wrapped envelope (RSA-OAEP-256 or AES-KWP) to the destination HSM. This pattern is operationally common — institutions may want IKM generation under their own air-gapped FIPS-validated workstation rather than the cloud provider's RNG. The three patterns named in §10.6.1 (HSM internal RNG, OS-level CSPRNG, dedicated CSPRNG hardware) cover the generation side but the wrapped-import composition is implied by §10.5 + §10.6.1 read together — not named explicitly.

### Response

The wrapped-import case is conformant by composition. The closing language makes the composition explicit so vendors implementing IKM-provisioning UX do not have to read §10.5 and §10.6.1 together to confirm conformance.

### Closing language for spec §10.6.1 (informative)

> **Wrapped-import composition.** The IKM is generated on a FIPS-validated CSPRNG outside the destination HSM (under one of the three patterns above — HSM internal RNG at a separate institution-controlled HSM, OS-level CSPRNG, or dedicated CSPRNG hardware) and imported into the destination HSM via the cloud provider's customer-supplied-key import workflow under RSA-OAEP-256 or AES-KWP wrapping. Conformant when (a) the source-side generation satisfies one of the three patterns above, (b) the destination-side custody satisfies §10.5 (FIPS 140-2 Level 3 or higher), and (c) the institution's CC8.1 names BOTH endpoints with their respective FIPS validation evidence and the wrapping algorithm in use.

This is added as a fourth bullet under the existing three RNG patterns. The chain's integrity claim is unchanged. The institution's CC8.1 covers both endpoints; the institution's evidence under HITRUST 06.d (Cryptographic Key Management, where applicable) is parallel.

---

## §3. Q44 — Session-key-only compromise non-claim

### What the reviewer asked

§1.2's forward-secrecy paragraph addresses the IKM-compromise scenario. The narrower scenario — an attacker holding a single `session_key` (extracted from the SDK process's RAM during a long-running run) but NOT the IKM — is not named. An implementer reading §1.2 might assume "no forward secrecy" implies "any compromise breaks any past event for any tenant," when the actual scope is narrower under session-key-only compromise.

### Response

This is informative-completeness, not a normative gap. The chain's integrity properties under session-key-only compromise are well-defined; the spec just doesn't narrate them. Adding the paragraph makes the implementer's threat-model analysis cleaner and reduces over-claiming or under-claiming the chain's properties in vendor security documentation.

### Closing language for spec §1.2 (informative)

> **Session-key-only compromise (informative).** An attacker who extracts a single `session_key` from a running SDK process (process-memory dump during the run) but does NOT compromise the tenant's IKM can forge events for the specific (`tenant_id`, `run_id`) pair the session_key was derived for, until the run terminates and the SDK destroys its session_key state. The attacker cannot forge events for other tenants or other runs because each (tenant, run) pair derives an isolated session_key from HKDF over the tenant-bound `info` parameter. Layered defenses against this scenario: (a) run lifecycle (§4.1 — runs are bound by `run_id`, and the writer destroys session_key when the run closes); (b) per-run state isolation (§4.1 inviolate property 7 — MAC input is canonical JSON of application content, not raw session_key fragments, so memory-disclosure of one event does not leak the session_key); and (c) the daily Merkle bound (events the attacker tries to retroactively insert into a prior day are caught at §7 step 12 because that day's seal already covered the legitimate event set). This non-claim is documented for institutions whose threat model includes hostile process-memory access. Institutions requiring stronger session-key isolation operate compensating controls outside the chain (process-isolation, runtime attestation, key-derivation-on-demand from HSM rather than process-resident session_key).

---

## §4. Q45 — Retention-vs-deprecation interaction

### What the reviewer asked

§1.2 says the chain's effective security level (128 bits) is "appropriate for FFIEC banking-regulation horizons (typically 7-year retention)." NIST SP 800-57 Part 1 Rev 5 names 2030 as the year ECDSA-P256-class signatures (~128-bit security) are no longer recommended for new deployments. A chain produced today and retained 7 years lands its verifier in 2033 — past the deprecation boundary for new-deployment recommendations. A chain ALREADY SEALED under v1.0 algorithms cannot retroactively benefit from v1.x's stronger algorithms. Whether 2033 verification under deprecated-for-new-deployments algorithms is acceptable for the chain's evidentiary purpose is a regulatory question.

### Response

The integrity claim under v1.0 algorithms during the retention horizon does not depend on those algorithms remaining NIST-recommended for new deployments. The integrity claim depends on the algorithms not being broken in practice. Adding the paragraph names the distinction explicitly so institutions evaluating long-retention deployments under v1.0 understand the posture.

### Closing language for spec §1.2 (informative)

> **Retention-vs-deprecation interaction (informative).** A chain sealed under v1.0 algorithms is verifiable under those algorithms throughout its retention horizon. A 2033 verifier verifying a 2026 v1.0 chain operates under v1.0 algorithm dispatch even though NIST guidance for new deployments may have moved on by 2033. Re-sealing of historical chains under post-quantum algorithms is OPTIONAL and is the institution's CC8.1 procedure decision; the institution may operate v1.x dual-algorithm posture for new chains while continuing to verify legacy chains under v1.0 dispatch. The chain's integrity claim under v1.0 algorithms during the retention horizon does not depend on those algorithms remaining NIST-recommended for new deployments — the integrity claim depends on the algorithms not being broken in practice (a publicly-known practical attack producing verifying tampered seals). Institutions with retention horizons exceeding 10 years and concerns about algorithm deprecation operate under v1.x dual-algorithm posture (Ed25519 + Dilithium per §4.3.2) for new chains AND maintain a re-sealing procedure to migrate historical chains to post-quantum algorithms before the deprecation horizon.

For healthcare deployments under multi-decade retention (per `docs/regulator-pack/healthcare-overlay.md` §H4 covering pediatric records under state law), the re-sealing procedure is operationally significant. The healthcare-overlay names re-keying as a forward-scope option.

---

## §5. Q46 — Linux LTS kernel wording

### What the reviewer asked

The `/dev/urandom on Linux (which is CSPRNG-quality on modern kernels per Linux 5.18+ documentation)` parenthetical reads literally as "kernel ≥ 5.18 only." LTS kernels in current production support windows (RHEL 7/8 ship 4.18/5.4; Ubuntu 18.04/20.04/22.04 LTS ship 4.15/5.4/5.15; AWS Linux 2 ships 4.14/5.10/5.15) are pre-5.18 but produce CSPRNG-quality output post-boot-entropy-completion. The spec's intent is "modern CSPRNG quality"; the wording reads as "kernel ≥ 5.18 only."

### Response

Editorial fix. The substantive intent is preserved; the wording change avoids the over-rejection of conformant LTS deployments.

### Closing language for spec §10.6.1

> /dev/urandom on Linux meeting the boot-time entropy condition (Linux 5.18+ unconditionally; earlier kernels acceptable when the institution's CC8.1 evidences boot-entropy gathering completed before IKM generation per the kernel's `getrandom` semantics — typically by initialising IKM after the system has reached multi-user state and the kernel's entropy pool is fully seeded).

This replaces the existing parenthetical. Vendors testing against LTS distributions can confirm conformance via the boot-entropy condition rather than over-rejecting based on kernel-version literal reading.

---

## §6. Q47 — Verifier amendment-awareness as a CC8.1 attribute

### What the reviewer asked

A verifier compiled before the v1.0a amendment landed has no `sign_payload_version` reading code. Such a verifier presented with a v1.0a seal computes a signature over the 6-line form; the signer wrote the signature over the 10-line form; signature verification fails; the verifier reports `signature verification failed` (the spec's normative reason from §7 step 11). The reason string is correct but uninformative — the actual root cause is verifier amendment-staleness. The spec's CC8.1 obligations don't currently name verifier amendment-awareness as a tracked attribute.

### Response

Operational gap. The institution's CC8.1 names the verifier amendment-level the institution operates against; the SOC engagement tests for staleness against the seals the institution actually produces. The closing path adds the CC8.1 attribute and pairs it with an audit procedure (P-55).

### Closing language for spec §10.12 (or §10.5)

> Institutions track verifier amendment-awareness as a CC8.1 attribute. The institution's verifier deployment names the spec amendment level the verifier supports (e.g., `v1.0` pre-amendment, `v1.0-final-amendment` with `v1.0a` `sign_payload_version` support). A verifier rolled out before a spec amendment landed is forward-incompatible with seals produced under the new amendment; the institution's CC8.1 procedure tests verifier amendment-awareness against the seals the institution actually produces and refreshes verifier deployment when the institution's posture progresses to a new amendment level. The audit procedure (P-55) samples the institution's verifier output against current-amendment seals and confirms the verifier's amendment-awareness.

The audit procedure addition (P-55) is in `docs/audit-procedures-cluster-d-addenda.md` — the institution's auditor confirms the verifier deployment in production handles the seal payload-version dispatch correctly for the seals the institution actually produces.

---

## §7. Q48 — §1.1 Daubert SDK-process compromise (carried from round 4)

### What the reviewer asked

The §1.1 Daubert four-factor "known error rate" paragraph names three custody layers (IKM, ledger, HSM) without naming the SDK-process compromise vector. A compromised SDK process holding a legitimate session_key produces verifying chain entries the institution did not authorize. The chain's §7 dispatch reports `Status: PASS`. Under cross-examination, an expert witness laying foundation under FRE 702 / Daubert is asked "what about a compromised SDK process?" — and the §1.1 informative paragraph does not have a stock answer. The institution's actual defense exists (SDK process integrity controls per §10.5, runtime attestation, signed-binary attestation) but is not named in the §1.1 informative grounding.

### Response

This is the highest-priority finding from round 5 because it is defensive of the institution's existing courtroom posture. The closing language extends the §1.1 paragraph to acknowledge the SDK-process compromise vector and name the institution's actual defenses. The chain's integrity properties do not change; the framing of the residual risk does.

### Closing language for spec §1.1

> A successful false-negative — a tampered chain that verifies as PASS — requires either (a) the simultaneous compromise of three independent custody layers (the tenant's IKM held in HSM/KMS, the institution's ledger storage held under append-only operator-side controls, and the HSM signing key held under FIPS 140-2 Level 3 or higher) — these three layers are operated by different roles under separation-of-duties controls and the compromise of any one alone does not produce a verifying tamper; OR (b) the compromise of an SDK process during the window the SDK held an active `session_key`, in which case the false-negative is bounded to events the compromised process could produce during the run lifecycle (forge-events-as-this-run-during-this-process). Defenses against (b) are SDK process integrity (FIPS 140-2 host validation, runtime attestation under confidential-computing posture, signed-binary attestation per §10.5 and §10.7), run lifecycle hygiene (each run terminates and destroys its session_key, bounding the forge window), and the daily Merkle bound (events forged during the compromise window cannot be retroactively inserted into prior days because the prior days' seals already cover the legitimate event sets); these are institutional concerns parallel to the chain's cryptographic primitives, and the institution's CC8.1 procedure documents the SDK-process integrity controls the institution operates.

The expert witness now has a stock answer when asked about SDK-process compromise on cross-examination: the false-negative is bounded; the bounding mechanism is run lifecycle and the daily Merkle bound; the institution's compensating controls are named in CC8.1.

---

## §8. Vendor-side posture

The Herald vendor's Compliance module already implements §10.6.1's OS-level CSPRNG via .NET's `RandomNumberGenerator`. The wrapped-import case (Q43) is a deferred implementation question for the vendor's IKM-provisioning UX — addressing it does not require code change to the existing chain implementation; it requires UX work to expose the wrapped-import workflow to the institution's IKM-provisioning operator.

The round-1 Q5 (Ed25519 absence on the candidate vendor's signer) remains the only vendor-side open from rounds 1-5 and is unchanged by round 5. Round-5 raises no new vendor-side gaps.

---

## §9. Cross-rounds tally

The cross-rounds tally now reads: 49 distinct items raised across rounds 1-5 (12 + 11 + 12 + 7 + 7). Of those:

- 31 closed via spec text (post-amendment-wave).
- 2 confirmed (no action needed; spec already addresses).
- 16 carried forward at round-5 close: 1 vendor-side (Q5 — Ed25519 absence on the candidate vendor's signer); 15 spec-quality items including the 7 in this response and 8 from earlier rounds pending working-group cycle.

The spec is substantially mature against the implementation contract. Remaining items are precision and informative-completeness, not integrity-bearing. A working-group cycle can close all 7 round-5 items in editorial pass; Q48 stages first because it is court-facing and was open from round 4.

---

## §10. Forward-scope acknowledgements

- **§10.6.1 wrapped-import composition.** Conformant by composition with §10.5 today; spec text addition is informative clarification.
- **§1.2 session-key-only compromise non-claim.** Informative; no normative change.
- **§1.2 retention-vs-deprecation.** Informative; institutions choose v1.0 retention or v1.x dual-algorithm posture for new chains based on retention horizon.
- **§1.1 Daubert framing.** Informative extension; no normative change to integrity properties.
- **§10.12 verifier amendment-awareness CC8.1 attribute.** Informative addition; existing verifier deployment posture is unchanged.

None of the closures require new normative cryptographic surface or breaking changes to wire format. The closure trail is editorial-pass-friendly.
