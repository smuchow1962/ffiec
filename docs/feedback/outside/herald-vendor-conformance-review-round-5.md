# Herald vendor conformance review — round 5

**Reviewer.** Same persona as rounds 1–4. Continuing the Herald-vendor lens against the v1.0-final-amendment spec as it stands today (post-2026-05-07 amendment-wave close-out).

**Why this round.** Two motivations. First, round 4's seven open findings (Q35–Q41). Second, the amendment wave introduced new normative material I haven't yet assessed: §1.2 forward-secrecy non-claim and effective security level, §10.6.1 IKM generation requirements, the §4.3 dual-algorithm AND-security paragraph, and the `sign_payload_version` discriminator lock. Round 5 confirms what closed, surfaces what didn't, and looks at the new sections with fresh eyes.

**Reading angle.** Unchanged from prior rounds: a vendor implementer who must build, ship, and operate a conformant SDK + verifier + KMS adapter, and who reads the spec as the contract. Findings are scoped to where the spec text is ambiguous, internally drifting, or absent in a way that makes a vendor write defensively or inconsistently. Cryptographic-strength concerns continue to defer to Reuven's lens (cryptographic security analyst, drop #04); privacy-law concerns defer to Elena (chief privacy officer, drop #05); evidentiary law concerns defer to Diego (digital forensics specialist, drop #03). My territory is the implementation contract.

**Distinct from chief privacy officer round (drop #05).** Elena's round examined privacy-law and data-protection obligations; this round examines implementation-contract precision in the new spec sections, plus the residual Daubert framing item Diego left for a later round.

---

## Round-4 closure confirmation

Six of round 4's seven findings closed:

- **Q36 (region definition).** §3 picked up `Region` as a defined term. Closed.
- **Q37, Q38 (Pattern A operational specificity).** §10.15 Pattern A invariants 4 and 6 absorbed the per-process region binding and the verifier's working-paper reference to per-region replication evidence. Closed.
- **Q39 (trusted-time placement).** §10.14 carries the v1.x forward-commitment line. Closed as forward-scope; the v1.x extension will name the placement when it lands.
- **Q40 (`ffiec.chain.late_binding` trust posture).** §4.2.2 added the explicit "Trust posture for `ffiec.chain.late_binding`" paragraph that names the attribute as ledger-stamped after ingest, NOT in the SDK-produced canonical bytes that the per-event MAC sealed. The §4.4 attribute-table row also names the trust posture inline. The implementer-misunderstanding scenario I flagged is closed at the spec layer.

**Q41 remains open** — the §1.1 Daubert four-factor framing still names three custody layers (IKM, ledger, HSM) without acknowledging the SDK-process compromise vector. I carry it into round 5 as Q48 below with the same closing language.

The cross-rounds tally now reads: 42 distinct items raised across rounds 1–4, of which 31 closed via spec text, 2 confirmed, 9 carried forward at round-4 close. Round 5 raises an additional 7 items below.

---

## Round-5 questions

### §10.6.1 IKM generation — verifier-side signal for RNG-source attestation

**Q42. The `master_key.generated` operational event records the RNG type, but the verifier never sees operational events.** §10.6.1 requires the institution to record the RNG source on the `master_key.generated` operational event (e.g. `"hsm.cloudhsm-classic"`, `"os.linux-urandom"`). The institution's CC8.1 control description names the RNG source the institution claims. The verifier (§7) processes the chain file and the seal records, not the operational-event stream. An institution claiming `"hsm.cloudhsm-classic"` in CC8.1 while their actual deployment used `"os.linux-urandom"` produces an indistinguishable chain — every per-event MAC is correct under whichever IKM was used, every Merkle root is correct, every seal signature is correct. The discrepancy surfaces only at SOC engagement when the auditor inspects the operational-event log and reconciles against CC8.1. There is no verifier-side working-paper line that the institution's MRM committee or examiner can read against the institution's declared posture.

**Status:** Partial.

**Closing language proposal.** §10.6.1 could add a paragraph naming the institution's RNG-source declaration as a tenant-public-key registry attribute parallel to the existing public-key fields, and require the verifier to print the institution's declared `rng_source` in its output (no FAIL/PASS impact, just a working-paper line). This makes the RNG-source decision verifier-visible without elevating it to a normative integrity check (which it cannot be — the chain holds under any unpredictable IKM regardless of RNG source).

### §10.6.1 IKM provisioning — wrapped-import is implied but not named

**Q43. The three conformant RNG patterns (§10.6.1) omit the wrapped-import workflow that AWS / Azure / Google Cloud KMS all support.** AWS KMS, Azure Key Vault Managed HSM, and Google Cloud KMS expose customer-supplied-key import workflows where the IKM is generated outside the HSM and flows in via a wrapped envelope (RSA-OAEP-256 or AES-KWP) to the destination HSM. This pattern is operationally common — an institution may want IKM generation under their own air-gapped FIPS-validated workstation rather than the cloud provider's RNG. The three patterns named (HSM internal RNG / OS-level CSPRNG / dedicated CSPRNG hardware) cover the generation side but the wrapped-import composition is not addressed. A vendor implementing the IKM-provisioning UX has to decide whether the import case is conformant by composition of "OS-level CSPRNG (generation)" + "HSM custody after import" — implied by §10.5's HSM custody requirement and §10.6.1's RNG patterns read together — but the spec doesn't name the composition explicitly.

**Status:** Partial.

**Closing language proposal.** §10.6.1 could add a fourth bullet: "**Wrapped-import.** The IKM is generated on a FIPS-validated CSPRNG outside the destination HSM (HSM, OS-level CSPRNG, or dedicated CSPRNG hardware per the patterns above) and imported into the destination HSM via the cloud provider's customer-supplied-key import workflow under RSA-OAEP-256 or AES-KWP wrapping. Conformant when (a) the source-side generation satisfies one of the three patterns above, (b) the destination-side custody satisfies §10.5 (FIPS 140-2 Level 3 or higher), and (c) the institution's CC8.1 names BOTH endpoints with their respective FIPS validation evidence."

### §1.2 forward-secrecy non-claim — session-key-only compromise scenario missing

**Q44. The forward-secrecy non-claim covers IKM compromise but skips the narrower session-key-only compromise.** §1.2's forward-secrecy paragraph addresses the IKM-compromise scenario: an attacker holding the IKM can derive any past session key for that tenant. The narrower scenario — an attacker holding a single `session_key` (e.g., extracted from the SDK process's RAM during a long-running run) but NOT the IKM — is not named. Under that compromise, the attacker can forge events for the specific (tenant_id, run_id) pair the session_key was derived for, but cannot forge cross-tenant or cross-run events because each (tenant, run) pair derives an isolated session key from HKDF over the tenant-bound `info` parameter and the spec's per-run state binding. The scope and the layered defenses against this scenario (run lifecycle, run-id isolation per §4.1, daily Merkle bound) are sound but not narrated. An implementer reading §1.2 might assume "no forward secrecy" implies "any compromise breaks any past event for any tenant," when the actual scope is narrower under session-key-only compromise.

**Status:** Partial.

**Closing language proposal.** §1.2 could add a paragraph after the existing forward-secrecy paragraph: "**Session-key-only compromise (informative).** An attacker who extracts a single `session_key` from a running SDK process (e.g., process-memory dump during the run) but does NOT compromise the tenant's IKM can forge events for the specific (tenant_id, run_id) pair the session_key was derived for, until the run terminates and the SDK destroys its session_key state. The attacker cannot forge events for other tenants or other runs because each (tenant, run) pair derives an isolated session_key from HKDF over the tenant-bound `info` parameter. The layered defenses against this scenario are: (a) run lifecycle (§4.1 — runs are bound by `run_id`, and the writer destroys session_key when the run closes), (b) per-run state isolation (§4.1 inviolate property 7 — MAC input is canonical JSON of application content, not raw session_key fragments, so memory-disclosure of one event does not leak the session_key), and (c) the daily Merkle bound (events the attacker tries to retroactively insert into a prior day are caught at §7 step 12 because that day's seal already covered the legitimate event set). This non-claim is documented for institutions whose threat model includes hostile process-memory access; institutions requiring stronger session-key isolation operate compensating controls outside the chain (process-isolation, runtime attestation, key-derivation-on-demand from HSM rather than process-resident session_key)."

### §1.2 effective security level — retention vs algorithm deprecation interaction

**Q45. The 128-bit effective-security-level statement does not address retention-vs-deprecation timing.** §1.2 says the chain's effective security level (128 bits) is "appropriate for FFIEC banking-regulation horizons (typically 7-year retention)." NIST SP 800-57 Part 1 Rev 5 names 2030 as the year ECDSA-P256-class signatures (~128-bit security) are no longer recommended for new deployments; SHA-2 family at 128-bit second-preimage resistance is recommended through 2030+ with caveats around long-term protection horizons. A chain produced today and retained 7 years lands its verifier in 2033 — past the deprecation boundary for new-deployment recommendations. The spec's §4.3.2 quantum-readiness commitment (30-day spec-patch SLA) covers v1.x evolution, but a chain ALREADY SEALED under v1.0 algorithms cannot retroactively benefit from v1.x's stronger algorithms. The 2033 verifier verifying a 2026-sealed chain is verifying under what will then be deprecated-for-new-deployments algorithms. Whether that is acceptable for the chain's evidentiary purpose is a regulatory question, not a cryptographic one — but the spec doesn't name the question.

**Status:** Partial.

**Closing language proposal.** §1.2 could add a paragraph after the effective-security-level paragraph: "**Retention vs deprecation interaction (informative).** A chain sealed under v1.0 algorithms is verifiable under those algorithms throughout its retention horizon. A 2033 verifier verifying a 2026 v1.0 chain operates under v1.0 algorithm dispatch even though NIST guidance for new deployments may have moved on by 2033. Re-sealing of historical chains under post-quantum algorithms is OPTIONAL and is the institution's CC8.1 procedure decision; the institution may operate v1.x dual-algorithm posture for new chains while continuing to verify legacy chains under v1.0 dispatch. The chain's integrity claim under v1.0 algorithms during the retention horizon does not depend on those algorithms remaining NIST-recommended for new deployments — the integrity claim depends on the algorithms not being broken in practice (i.e., a publicly-known practical attack producing verifying tampered seals). Institutions with retention horizons exceeding 10 years and concerns about algorithm deprecation operate under v1.x dual-algorithm posture (Ed25519 + Dilithium per §4.3.2) for new chains AND maintain a re-sealing procedure to migrate historical chains to post-quantum algorithms before the deprecation horizon."

### §10.6.1 OS-CSPRNG list — Linux 5.18+ wording is restrictive for LTS kernels

**Q46. The `/dev/urandom on Linux (which is CSPRNG-quality on modern kernels per Linux 5.18+ documentation)` parenthetical is read literally by implementers and excludes LTS kernels in current production support windows.** RHEL 7 / 8 ship 4.18 / 5.4; Ubuntu 18.04/20.04/22.04 LTS ship 4.15 / 5.4 / 5.15; AWS Linux 2 ships 4.14/5.10/5.15. All within their LTS support windows but pre-5.18. The semantics on those kernels are good — `/dev/urandom` post-boot-entropy-completion produces CSPRNG-quality output going back well before 5.18 — but the boot-time entropy gathering changed in 5.18 (the `RANDOM_TRUST_CPU` / `random: getrandom` semantics). A vendor reading the spec literally would worry about LTS kernel acceptability and would either over-reject conformant deployments or push customers to upgrade kernels they cannot upgrade in their support window. The spec's intent is "modern CSPRNG quality"; the wording reads as "kernel ≥ 5.18 only."

**Status:** Nit.

**Closing language proposal.** Rephrase the bullet to: "/dev/urandom on Linux meeting the boot-time entropy condition (Linux 5.18+ unconditionally; earlier kernels acceptable when the institution's CC8.1 evidences boot-entropy gathering completed before IKM generation per the kernel's `getrandom` semantics — typically by initialising IKM after the system has reached multi-user state and the kernel's entropy pool is fully seeded)."

### §4.3 sign_payload_version — verifier amendment-awareness should be a CC8.1 attribute

**Q47. A verifier compiled before the v1.0a amendment landed is silently forward-incompatible.** The §7 step 11 dispatch handles the bidirectional case correctly: an amendment-aware verifier reads `sign_payload_version` and dispatches to the pre-amendment 6-line form (when absent) or the amendment 10-line form (when `"v1.0a"`). A verifier compiled BEFORE the v1.0a amendment landed has no `sign_payload_version` reading code at all — it reconstructs the 6-line form for every seal. Such a verifier presented with a v1.0a seal will compute a signature over the 6-line form, the signer wrote the signature over the 10-line form, the signature verification fails, the verifier reports `signature verification failed` (the spec's normative reason from §7 step 11). The reason string is correct but uninformative — the actual root cause is verifier amendment-staleness, not signature tampering. An institution running stale verifier code against post-amendment seals will see chains-fail-verification at the moment of upgrade, not at the moment the verifier was deployed. The spec's CC8.1 obligations don't currently name verifier amendment-awareness as a tracked attribute.

**Status:** Partial.

**Closing language proposal.** §10.12 (verifier CLI exit-code contract) or §10.5 (custody) could add: "Institutions track verifier amendment-awareness as a CC8.1 attribute. The institution's verifier deployment names the spec amendment level the verifier supports (e.g., `v1.0` pre-amendment, `v1.0-final-amendment` with `v1.0a` `sign_payload_version` support). A verifier rolled out before a spec amendment landed is forward-incompatible with seals produced under the new amendment; the institution's CC8.1 procedure tests verifier amendment-awareness against the seals the institution actually produces and refreshes verifier deployment when the institution's posture progresses to a new amendment level."

### §1.1 Daubert four-factor framing — SDK-process compromise still missing

**Q48. Round-4 Q41 carried forward unchanged.** The §1.1 Daubert four-factor "known error rate" paragraph still names three custody layers (IKM, ledger, HSM) without naming the SDK-process compromise vector. A compromised SDK process holding a legitimate session_key produces verifying chain entries the institution did not authorize. The chain's §7 dispatch reports `Status: PASS`. Under cross-examination, an expert witness laying foundation under FRE 702 / Daubert is asked "what about a compromised SDK process?" — and the §1.1 informative paragraph does not have a stock answer. The institution's actual defense exists (SDK process integrity controls per §10.5, runtime attestation, signed-binary attestation) but is not named in the §1.1 informative grounding. This was Diego's and my round-4 finding; it remains open.

**Status:** Partial (carried from round 4).

**Closing language proposal.** §1.1's "known error rate" paragraph could be extended: "A successful false-negative — a tampered chain that verifies as PASS — requires either (a) the simultaneous compromise of three independent custody layers (the tenant's IKM held in HSM/KMS, the institution's ledger storage held under append-only operator-side controls, and the HSM signing key held under FIPS 140-2 Level 3 or higher) — these three layers are operated by different roles under separation-of-duties controls and the compromise of any one alone does not produce a verifying tamper; OR (b) the compromise of an SDK process during the window the SDK held an active `session_key`, in which case the false-negative is bounded to events the compromised process could produce during the run lifecycle (forge-events-as-this-run-during-this-process). Defenses against (b) are SDK process integrity (FIPS 140-2 host validation, runtime attestation under confidential-computing posture, signed-binary attestation per §10.5 and §10.7), run lifecycle hygiene (each run terminates and destroys its session_key, bounding the forge window), and the daily Merkle bound (events forged during the compromise window cannot be retroactively inserted into prior days because the prior days' seals already cover the legitimate event sets); these are institutional concerns parallel to the chain's cryptographic primitives, and the institution's CC8.1 procedure documents the SDK-process integrity controls the institution operates."

---

## Per-role roll-up

**For the spec working group.** Five Partials and one Nit against the new §1.2 / §10.6.1 / §4.3 sections, plus the round-4 Q41 carry-forward. None affect the cryptographic-integrity claim of the chain. All affect either implementer-precision (Q42, Q43, Q46, Q47), informative-completeness (Q44, Q45), or court-facing residual-risk framing (Q48).

**For institutions evaluating Herald vendor candidates.** No vendor-side gaps net-new from this round. Q5 (Ed25519 absence on the candidate vendor's signer) remains the only vendor-side open from round 1 and is unchanged.

**For Herald-the-vendor.** No new vendor-conformance gaps; the spec amendments closed substantively all of round-4's items. The vendor's Compliance module already implements §10.6.1's OS-level CSPRNG via .NET's `RandomNumberGenerator`; the wrapped-import case (Q43) is a deferred implementation question for the vendor's IKM-provisioning UX, not a current shipping gap.

---

## Where I would prioritize

1. **Q48 (§1.1 Daubert framing — SDK-process compromise).** Same as round 4: court-facing posture, expert witness foundation under FRE 702, cross-examination-defensive. One-paragraph spec patch. Highest priority because it's defensive of the institution's actual courtroom posture.

2. **Q44 (§1.2 session-key-only compromise non-claim).** Closes a subtle scope misunderstanding; one paragraph in §1.2.

3. **Q45 (§1.2 retention vs deprecation).** Forward-scope clarity; one paragraph in §1.2.

4. **Q42, Q43, Q46 (§10.6.1 details).** Three editorial passes; the Linux-LTS wording (Q46) is the most operationally pressing because it changes how vendors test against LTS distributions.

5. **Q47 (verifier amendment-awareness CC8.1 attribute).** Operational; informs the institution's verifier-upgrade procedure.

---

## Stopping criterion

This drop carries 6 Partial + 1 Nit findings against the v1.0-final-amendment spec as it stands today (post-2026-05-07 amendment-wave close-out). The findings cluster into informative-completeness in §1.2 (Q44, Q45), implementer-precision in §10.6.1 (Q42, Q43, Q46), forward-incompatibility in verifier deployment (Q47), and the carried-forward §1.1 Daubert framing residual-risk gap (Q48). None of them break the four primitives' cryptographic integrity claims; they affect implementer-side defensiveness, court-facing evidence framing, and the institution's CC8.1 procedure completeness. A working-group cycle can close them in editorial pass; Q48 is the one I'd stage first because it's defensive of the institution's existing courtroom posture and remains open from round 4.

The cross-rounds tally now: 49 distinct items raised across rounds 1–5 (12 + 11 + 12 + 7 + 7). Of those, 31 closed via spec text, 2 confirmed, 16 carried forward at round-5 close (1 vendor-side Q5, 15 spec-quality items pending working-group cycle). The spec is substantially mature against the implementation contract; remaining items are precision and informative-completeness, not integrity-bearing.
