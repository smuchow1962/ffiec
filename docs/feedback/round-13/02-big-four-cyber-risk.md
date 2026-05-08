# Round 13 — Big Four Cyber Risk review

> **Reviewer.** Lukas Holzer, Senior Manager, Cyber Risk Advisory (Frankfurt office).
> **Background.** Twenty-one years across European and US bank cybersecurity advisory; current focus on supply-chain SLSA L3 and post-quantum migration governance.
> **Lens.** HSM/KMS operational practice, BYOC vendor management, supply-chain trust path (cosign / GPG / SLSA / cold-DR), IR posture for chain-detected events, regulator-fingerprint lifecycle, dual-algorithm transition.
> **Materials reviewed.**
> - `spec/chain-of-custody-v1.md` §10 (operational events including the new `regulator_fingerprint.*` triple)
> - `docs/design/00-overview.md` §2 four-attack catalog
> - `docs/design/04-hsm-custody.md` §3.2 dual-algorithm keypair
> - `docs/design/09-threat-model.md` §2.9 Adversary I + reception procedure
> - `docs/byoc-deployment.md`
> - `docs/supply-chain.md` (cold-DR consumption procedure, slsa-verifier trust anchor, mirror WORM, GPG recovery, dual-compromise)
> - `docs/incident-response-playbook.md` (Scenarios 11 and 12, dual-algorithm verifier-version overlap, federal-regulator routing)
> - `docs/cloud-hsm-guide.md`
> **Stopping criterion.** 0 gaps, 0 partials.
> **Round-13 disposition.** PASS — 0 gaps, 0 partials.

---

## 1. Executive summary

The corpus crosses my Big Four cyber-risk threshold. Five things drove the disposition:

1. The supply-chain document closes the loop on the cold-disaster-recovery key. It does not just describe the project-side lifecycle — it specifies what the institution does each year with the dry-run attestation, what control-evidence artifact the institution archives, what the IR escalation looks like when the institution fails to obtain the attestation in time, and which CSF subcategory binds the consumption log. That is the difference between documented-but-unexercised and operationally exercised. GV.SC-04 examiners do not award credit for the former.
2. `slsa-verifier` is named as a co-equal trust artifact alongside cosign and the GPG-signed manifest, with its own public key and binary SHA-256 in the trust-anchor table. The institution's control description must name the SLSA-aware tool, version, trust anchor, and archive location in one place. That four-question test is the one I have been waiting to see in writing — without it, an L3 provenance claim is a marketing line, not a control.
3. The mirror-registry re-signing pattern carries a hard 7-year minimum on the audit log AND a WORM storage requirement AND a continuous-monitoring control (the `mirror.reconciliation_completed` event) that fires when the mirror silently stops validating project signatures. The bridge invariant is no longer a snapshot belief; it is a sampled-and-monitored discipline retroactively examinable.
4. The Adversary I reception procedure for regulator-held fingerprint rotation gives the institution three institution-defined operational events and a reception-failure sub-variant that activates IR Scenario 11. This is the right shape — the institution operates the procedure, the operational events are the audit trail, the SOC team and examiner sample-test against them, and forged-notice handling is named.
5. The IR playbook handles the dual-algorithm transition seriously: Scenario 12 names case-(e) co-signed seal failure with a three-branch triage tree (published break / per-algorithm key compromise / under-investigation), each branch with its own clock-start disposition. The verifier-version-overlap edge case names the institution-vs-regulator posture-disagreement matrix and the pre-coordination duty. The federal-regulator routing matrix is pre-documented per charter type, including FBO subsidiaries. These are the operationally specific pieces that distinguish a playbook from a checklist.

I have nothing left to escalate. The remaining items I would normally pursue are either (a) explicitly accepted residuals the threat model documents (HSM physical compromise, cryptographic primitive break, monoculture, application-host compromise) or (b) deferred by the spec to v1.1 with documented compensating controls (multi-region replication semantics; secure-enclave-attestation defeat at edge nodes). Both categories are appropriate for the v1.0-final disposition.

---

## 2. What I tested

### 2.1 HSM custody and operational practice

| Test | Reference | Result |
|---|---|---|
| FIPS 140-2 L3 conformance bar named, with disqualified products explicitly enumerated | `cloud-hsm-guide.md` matrix; `04-hsm-custody.md` §3.2.1 | PASS — AWS KMS default tier and Azure Key Vault Standard are flagged non-conformant by name; the cross-tier-mixing pitfall is called out |
| Per-tenant signing keypair, sign-only operator role, separation from HSM admin | `04-hsm-custody.md` §3.2, §5.2.1; `09-threat-model.md` §2.3 | PASS — separation between seal-job operator, HSM administrator, DBA documented; small-institution dual-control compensating control named |
| HSM PIN rotation cadence and procedure | `04-hsm-custody.md` §5.2.1 | PASS — quarterly default, never-in-config rule, five-step rotation, prior-PIN invalidation step included |
| HSM unavailability with 72-hour SHOULD-notify | spec §4.3.1; `04-hsm-custody.md` §5.2 | PASS — distinct from the 36-hour cyber-incident path; both clocks documented |
| Cross-region key replication preserves per-tenant labels | `04-hsm-custody.md` §3.2.1; `cloud-hsm-guide.md` per-cloud sections | PASS — cross-tenant signing structurally prevented across regions |
| Dual-algorithm keypair for post-quantum coexistence | `04-hsm-custody.md` §3.2 (per-tenant per-algorithm); spec §4.2 `signatures` list (Variant B) | PASS — algorithm-bound `sign_payload` rules out cross-algorithm replay (R13); both keys per tenant, both publics in the registry |

### 2.2 BYOC vendor management

| Test | Reference | Result |
|---|---|---|
| Vendor cannot reach HSM, master key, or event-payload data | `byoc-deployment.md` IAM matrix and network controls table | PASS — runtime IAM identity is bank-controlled and bank-revocable; HSM is in private subnet with no internet egress |
| Mirror-registry pull discipline (cosign-validated, scanned) | `byoc-deployment.md` Step 2 | PASS — pull-rules with cosign verification at replication; vendor never pushes to bank compute |
| Vendor support telemetry redaction | `byoc-deployment.md` Step 5 | PASS — bank-controlled OTel Collector enforces redaction; bank can revoke egress unilaterally |
| Vendor-relationship termination procedure | `byoc-deployment.md` "When the vendor relationship ends" | PASS — past events remain verifiable under the original public key after vendor switch; the open-standard plus bank-held keys make this tractable |

### 2.3 Supply-chain trust path

| Test | Reference | Result |
|---|---|---|
| Cosign as primary, GPG-signed hash manifest as fallback, role-separation between signers | `supply-chain.md` "Signing roles" | PASS — single-compromise across both anchors does not collapse the trust path |
| GPG key rotation cadence (36 months default; 7-day emergency; 12-month subkey) | `supply-chain.md` "GPG key rotation cadence" | PASS — aligned with NIST SP 800-57 offline-key guidance; transition statement co-signed under old and new key |
| Cosign-key recovery procedure with timing (24h project response / 48h institution notification / 7 days new-key chain) | `supply-chain.md` "Cosign-key recovery procedure" | PASS — mid-examination variant covers the GPG-only interim posture and post-recovery re-validation, with `bundleverify --gpg-only` named |
| GPG-key recovery procedure (symmetric; cosign signs the GPG-revocation notice) | `supply-chain.md` "GPG-key recovery procedure" | PASS — the cross-signing handles the obvious edge case |
| Dual-compromise edge case (cosign AND GPG simultaneously compromised) | `supply-chain.md` "Dual-compromise edge case" | PASS — cold-DR key is the third anchor; institution falls back to the historical pre-compromise verifier; coordinates directly with primary regulator on examination posture |
| Cold-DR-key lifecycle: 60-month rotation, third separated key-holder group, annual dry-run, hardware health-check | `supply-chain.md` "Cold-disaster-recovery key lifecycle" | PASS — the dry-run is published as `KEY-DR-DRYRUN-{year}.asc` so the institution can consume it as control evidence |
| Institution-side consumption of the cold-DR dry-run attestation | `supply-chain.md` "Institution-side consumption of cold-DR-key dry-run attestation" | PASS — five-step procedure, 30-day grace window, IR Scenario 11 branch (c) for missed window, control-evidence binding to CSF GV.SC-04 |
| SLSA L3 provenance attestation institutional consumption (deployment-gating) | `supply-chain.md` "Institutional consumption of the SLSA L3 attestation" | PASS — `slsa-verifier` is a co-equal trust anchor with its own public key and binary SHA-256; the four-question control-description test (which tool, what version, what anchor, archived where) is named |
| Mirror-registry re-signing pattern documented bridging (when used) | `supply-chain.md` "Mirror-registry signature handling" | PASS — bridge invariant asserted; missing mirror-pull-validation log entry is a Scenario 3 finding |
| Mirror audit-log retention (7 years) and WORM integrity | `supply-chain.md` "Mirror audit-log retention and integrity" | PASS — acceptable backends per cloud named (S3 Object Lock, Azure immutable storage, GCP locked retention bucket, on-prem WORM) |
| Mirror continuous-monitoring control (`mirror.reconciliation_completed`) | `supply-chain.md` "Mirror continuous-monitoring control" | PASS — same control pattern as spec §10.1 fingerprint reconciliation, applied to the mirror's signature-validation discipline; closes the silent-mirror-misconfig gap (DE.CM-09) |
| Spec-text and conformance-corpus signed-artifact trust path | `supply-chain.md` "Spec text and conformance corpus integrity" | PASS — corpus tarball is cosign-and-GPG-signed; institutions that adopt for examination-grade use rebuild from spec text |

### 2.4 Regulator-fingerprint lifecycle (Adversary I reception)

| Test | Reference | Result |
|---|---|---|
| Three institution-defined operational events for the reception procedure | `09-threat-model.md` §2.9; spec §10.2 (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`) | PASS — events are first-class citizens in the §10.2 catalog; field schemas specified per event |
| Reception-failure sub-variant (forged-notice suspected) | `09-threat-model.md` §2.9 step 2 | PASS — institution does NOT install; activates IR Scenario 11 sub-variant for trust-anchor reception failure; cross-checks via out-of-band regulator contact |
| Historical-verifier-report re-validation across the rotation | `09-threat-model.md` §2.9 step 4 | PASS — confirms verifier output stability across the trust-anchor change |
| Examiner-laptop hygiene scoped OUT of the institution's CC6 | `09-threat-model.md` §2.9 | PASS — the document explicitly names this as a regulator-side control and tells SOC engagement partners not to misread it as an institution-owned gap |

### 2.5 IR posture for chain-detected events

| Test | Reference | Result |
|---|---|---|
| Twelve scenarios cover the chain-detected event surface | `incident-response-playbook.md` §"Common scenarios" | PASS — Scenarios 1-12 cover hash mismatch, Merkle mismatch, signature failure, master compromise, sealing delay, software-key fallback, fingerprint mismatch, unknown key_version, file truncation, backup integrity, cold-DR fallback, co-signed seal failure |
| 36-hour cyber-incident clock-start triage matrix per scenario | `incident-response-playbook.md` §"36-hour cyber-incident notification triage matrix" | PASS — clock-start is the determination point, NOT the alert point; per-scenario triggers documented |
| Scenario 11 (cold-DR fallback) — three trigger variants | `incident-response-playbook.md` Scenario 11 | PASS — failed dry-run / dual-compromise activation / missed dry-run window each have their own containment; clock does NOT start by itself (project-side governance, not institution-side cyber incident) |
| Scenario 12 (co-signed seal failure case (e)) — three-branch triage | `incident-response-playbook.md` Scenario 12 | PASS — published break (no clock) / per-algorithm key compromise (clock starts per Scenario 4) / under-investigation (clock starts at investigation-conclusion, bounded 48-72h to prevent indefinite deferral); spec §7 step 11 case (e) examiner working-paper convention referenced |
| Concurrent multi-scenario rollup rule | `incident-response-playbook.md` §"Edge cases" | PASS — earliest qualifying determination across the constellation, NOT union of per-scenario clocks; prevents double-counting and under-counting |
| Cross-tenant scope discovery (vendor-hosted, shared-cloud-HSM) | `incident-response-playbook.md` §"Edge cases" | PASS — vendor SLA and FFIEC 36-hour clock are separate, both tracked; institution's clock starts at institution determination, not at vendor alert |
| Federal-regulator routing per institution charter | `incident-response-playbook.md` §"Edge cases" | PASS — OCC / FRB / FDIC / NCUA mapping per charter type, including FBO subsidiaries; pre-documented in IR runbook |
| Dual-algorithm verifier-version timing-overlap (institution-vs-regulator posture) | `incident-response-playbook.md` §"Edge cases" | PASS — two posture-disagreement scenarios with named dispositions; pre-coordination duty before the multi-year transitional period begins |
| CIRCIA 72-hour vs FFIEC 36-hour scenario-by-scenario matrix | `incident-response-playbook.md` §"Edge cases" + §"CIRCIA and additional notification frameworks" | PASS — CIRCIA "substantial cyber incident" can fire without FFIEC "safety/soundness" firing; counsel works the matrix per incident |

### 2.6 Threat-model boundary correctness

| Test | Reference | Result |
|---|---|---|
| Adversaries A through I named, each with defense and residual risk | `09-threat-model.md` §2.1-§2.9 | PASS — Adversary I (examiner-side tooling subversion) explicit and first-class in the catalog; mirrors §2.4 of `00-overview.md` |
| Backup-integrity boundary (chain composes with WORM/immutable backups) | `09-threat-model.md` §3.1.1 | PASS — backup tampering is denial-of-service against the verifier report, not silent rewrite |
| Identity-attribution boundary (chain captures tenant-level; IAM layer captures principal) | `09-threat-model.md` §3.1.2 | PASS — composition with IAM produces end-to-end attribution; chain alone produces only the integrity-bearing record |
| Forgery-allegation evidence path | `09-threat-model.md` §3.1.3 | PASS — process audit logs, network logs, IDS/EDR, HSM operations log, key-rotation records named |
| Residual-risk register R1-R13 | `09-threat-model.md` §6 | PASS — each risk owner-tagged with mitigation status; R13 (algorithm-confusion at seal layer) closed by the v1.0-rework `sign_payload` extension |

---

## 3. What works particularly well

**The cold-DR key lifecycle is what I want every project-side governance document to look like.** Annual dry-run, published attestation, institution-side consumption procedure with a 30-day window, control-evidence binding to a specific CSF subcategory, IR playbook branch when consumption fails. The procedure makes the cold-DR fallback exercised rather than aspirational; the GV.SC-04 examiner has something to sample-test.

**The mirror-registry re-signing pattern's three controls compose cleanly.** Bridge documentation requirement closes the trust-anchor-substitution gap. WORM-backed 7-year audit-log retention closes the retroactive-examination gap. The continuous-monitoring control closes the silent-misconfig gap. Each gap is closed by a specific control, each control fires its own operational event, each event is sample-tested. That is the layering I look for.

**Scenario 12's three-branch triage is operationally honest.** The playbook does not pretend that a co-signed seal failure has a single root cause. Branch (i) names the published-break case where the regulator already knows. Branch (ii) names the per-algorithm key compromise where Scenario 4 takes over. Branch (iii) names the under-investigation case with a bounded 48-72h window so the clock-start cannot be deferred indefinitely. The bounded window is the part most playbooks omit and is the part that matters when counsel is asking "when does the 36-hour clock start running."

**The federal-regulator routing matrix per charter type is a piece I rarely see pre-documented.** OCC for national banks, FRB for state member banks and BHCs, FDIC for state non-member banks, NCUA for credit unions, and the FBO-subsidiary case named explicitly. The IR Commander does not stand up the routing decision under incident pressure; the standing IR documentation has it pre-recorded.

**Adversary I's reception procedure for the regulator-held fingerprint rotation is well-shaped.** Three operational events form the audit trail. The reception-failure sub-variant covers forged-notice handling. The institution does NOT install on validation failure; it cross-checks out-of-band and activates IR Scenario 11. That is the difference between a procedure that works under attack and a procedure that only works on the happy path.

**The dual-algorithm transitional period gets precise spec-and-IR coverage.** Variant B `sign_payload` per algorithm closes the algorithm-confusion attack class (R13) before the transition begins. The verifier dispatches per-algorithm and reports both validation outcomes in the working paper. Scenario 12 names the case-(e) severity-is-Severe-regardless-of-bracket convention so an examiner cannot misread non-strict PASS-WITH-ANOMALY as a low-severity disposition. The verifier-version-overlap matrix names what happens when the institution's verifier supports algorithm Y but the regulator's verifier still supports only algorithm X. Pre-coordination duty is named.

---

## 4. Residual concerns I am not escalating

I want to be explicit about residuals I considered and chose not to escalate. They are documented residuals, not gaps.

- **HSM physical compromise.** Accepted as residual R3 in `09-threat-model.md` §6. FIPS 140-2 L3 tamper detection is the technical defense; tamper detection plus immediate revocation is the operational response. A nation-state-level adversary may compromise the HSM; the institution cannot eliminate this and the threat model does not pretend otherwise.
- **Cryptographic primitive break.** Accepted as residual R6. The 30-day spec-patch SLA per spec §4.3.2 plus the algorithm-rotation provision plus the dual-algorithm coexistence design are the operational response. The migration lag is operational, not technical.
- **Operational monoculture.** Accepted as a not-claimed item in `09-threat-model.md` §4. The defense is independent re-implementation per §5; institutions operating two implementations in parallel get the diversity benefit. The chain spec cannot mandate diversity at the institution layer; it makes diversity achievable.
- **Application-host compromise (Adversary F).** Accepted as residual R1 with the mitigations named (host hardening, rotation on detection). The compromise window between attacker access and institution detection is bounded by the institution's broader detection program, not by the chain construction.
- **Multi-region replication semantics.** Deferred to v1.1 per `00-overview.md` §6.4. The v1.0 workaround (one ledger per region with per-region tenant_id) is heavier than the eventual design but defensible; institutions whose resilience program requires multi-region operate it today.
- **Edge-device physical compromise (R12).** Deferred to v1.1. The compensating control (per-device IKM in TPM/secure-enclave; reconciliation cadence baseline tuned per fleet) is documented in `edge-and-federated-ai.md`; institutions that operate edge today document the residual and operate the compensating controls.

Each of these is a documented residual at the spec-version level. None reads as "we forgot to think about this." All read as "we know about this and the v1.0-final scope ends here."

---

## 5. Disposition

**0 gaps. 0 partials.** The corpus passes the Big Four cyber-risk threshold for v1.0-final.

The HSM custody story, the BYOC vendor-management story, the supply-chain trust path, the regulator-fingerprint reception procedure, and the IR playbook compose into a defensible cybersecurity-control story. The pieces I would normally challenge in this lens — silent mirror misconfig, undocumented SLSA tool choice, un-exercised cold-DR fallback, undefined post-quantum transition controls, missing federal-regulator routing under incident pressure, ambiguous case-(e) co-signed seal handling, regulator-fingerprint reception under forged-notice attack — are each closed with a specific control, an operational event, an audit-evidence binding, and a sample-test path.

I have nothing to add to the round-13 inventory.

---

## 6. Sign-off

**Reviewer.** Lukas Holzer, Senior Manager, Cyber Risk Advisory.
**Date.** 2026-05-06.
**Round.** 13.
**Disposition.** PASS — 0 gaps, 0 partials.
