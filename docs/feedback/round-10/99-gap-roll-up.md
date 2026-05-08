# 99 — Gap roll-up (round 10 — first review of post-round-9-closure corpus)

> **What this doc is.** Cumulative finding inventory across the six round-10 reviewer personas. Stopping criterion: 0 gaps and 0 partials per role document. Round 10 is the second pass against the v1.0-rework spec; round 9 closed 9 gaps + 26 partials, round 10 assesses whether the closures held and whether new finds surfaced.

## Per-role counts

| Role | Persona | Answered | Partial | Gap |
|---|---|---|---|---|
| Big Four — Cryptographic | Dr. Tomás Bergström | 4 | 0 | 4 |
| Big Four — Cyber Risk | Priya Vasishta | 8 | 5 | 0 |
| Big Four — SOC Audit | Diane Mehrotra | 14 strengths | 8 | 0 |
| Big Four — Model Risk | Dr. Hiroshi Bergmann | 5 | 0 | 0 |
| FFIEC — IT Examiner | Marcus Whitfield | 5 | 4 (below-bar asks) | 0 |
| FFIEC — Cybersecurity Specialist | Liang Zhao | 14 strengths | 4 | 1 |
| **Totals** | | — | **21** | **5** |

**Round 10 trajectory.** Round 9 closed 9 gaps + 26 partials. Round 10 lands at 5 gaps + 21 partials. Net reduction in items: 35 → 26 (26% reduction). The remaining items are concentrated in cryptographic plumbing, surface-doc consistency, and cyber-risk operational refinements; no structural redesign.

**Convergence not yet reached.** Round 11 will close the items below.

## Findings inventory

### Cryptographic engineering (Tomás Bergström)

| ID | Severity | Item |
|---|---|---|
| T-1 | Gap | `docs/design/07-verifier-design.md` §4.3 sign_payload pseudocode is missing the `algorithm` line. Six lines instead of seven; contradicts spec §4.3 and the regenerated test vectors. |
| T-2 | Gap | `docs/design/04-hsm-custody.md` §3.2.2 still uses pre-rework `master_version`/`session_key_id` terminology, contradicting the rework's `key_version`/`key_fingerprint` model. |
| T-3 | Gap | `docs/design/02-chain-construction.md` has DUPLICATE §4.1 "Handshake security floor" sections (one rewritten with four properties; one pre-rework leftover with three). Likely affects §4.2 similarly. Pre-rework leftover must be removed. |
| T-4 | Gap | `spec/test-vectors/negative/` ships only README.md — none of the 16 documented negative cases (especially N006/N013/N014, the load-bearing rework cases) have per-case fixture files on disk. |

### Big Four cyber risk (Priya Vasishta)

| ID | Severity | Item |
|---|---|---|
| P-1 | Partial | Mirror-registry bridge audit-log retention duration unspecified; recommend WORM posture for 7 years to match chain retention. |
| P-2 | Partial | SLSA L3 institutional consumption — `slsa-verifier` should be MUST-tier with co-equal trust anchor; add to "What examiners verify" checklist. |
| P-3 | Partial | GPG-key recovery procedure analogous to cosign-key recovery; dual-compromise institution-side fallback path. |
| P-4 | Partial | Adversary I residual scope clarification — examiner-laptop hygiene named as regulator-side control (out of scope for institution's CC6 evaluation). |
| P-5 | Partial | Cosign-key recovery during active examination — bank-continuity posture during the 24/48/7-day recovery window. |

### Big Four SOC audit (Diane Mehrotra)

| ID | Severity | Item |
|---|---|---|
| D-1 | Partial | P-3 self-test cadence unstated; recommend lesser-of-IKM-rotation-or-P-6-cadence floor. |
| D-2 | Partial | P-3 sample-comparison non-production IKM provisioning, transport, retention rules unstated. |
| D-3 | Partial | CUECs.md CUEC-CRY-04 vocabulary drift: still says "session-key-id reconciliation"; spec/P-6 use "key-fingerprint reconciliation." |
| D-4 | Partial | Section 4 template missing three procedure bullets: IKM retirement (§10.9), mid-write-truncation recovery (IR Scenario 9), `--strict` vs operational verifier-run cadence. |
| D-5 | Partial | TSC mapping does not explicitly home `audit_file.truncation_detected` (CC7.2) or `master_key.retired` (CC6.1 + CC8.1). |
| D-6 | Partial | Sample-report appendix missing two rows: `Late-binding count` (PI1.2), `KMS handle URI` with `plaintext-` prefix under `--strict` (CC6.8). |
| D-7 | Partial | No P-22/P-23/P-24 reciprocal SOC tests for the four-piece evidence package on `key_fingerprint mismatch`, `hkdf_inputs_digest mismatch`, `audit_file_truncation_detected`. |
| D-8 | Partial | User-entity summary should mention the two new events (`audit_file.truncation_detected`, `master_key.retired`) one-line each. |
| D-9 | Partial | Anomaly template should ship a worked example record for `key_fingerprint mismatch` (template-quality, friction-reducer). |
| D-10 | Partial | Anomaly template `hkdf_inputs_digest mismatch` wording should address the multi-region deployment-drift case (two production SDKs at the same `format_version` producing different digests). |
| D-11 | Partial | Anomaly template `audit_file_truncation_detected` wording should specify `partial` recovery disposition (recovered events restored, unrecovered subset → IR Scenario 9 unrecoverable-gap). |

### FFIEC IT examiner (Marcus Whitfield)

| ID | Severity | Item |
|---|---|---|
| M-1 | Partial | examiner-quickstart.md "stop and call the bank" framing inconsistency — `key_fingerprint mismatch` and `payload_hash MAC mismatch` both deserve the call-out, OR neither. |
| M-2 | Partial | sample-report.md FAIL day excerpt should show at least one PASS day directly adjacent for visual comparison. |
| M-3 | Partial | examiner-quickstart.md 12-row table should add an IR scenario column (lifts from `examination-response-workflow.md` step 2). |
| M-4 | Partial | examination-response-workflow.md should add a second worked JSON example for a structurally different failure type (e.g., `step: 10` Merkle root mismatch). |
| M-5 (v1.1) | Wishlist | `verifier explain <step>` CLI subcommand — NOT a v1 blocker; tracked for v1.1. |

### FFIEC cybersecurity specialist (Liang Zhao)

| ID | Severity | Item |
|---|---|---|
| L-1 | Gap | CSF mapping doc not indexed by `00-overview.md` §2 four-attack catalog; add four-row table mapping each attack to load-bearing CSF subcategory + operational event + IR scenario. |
| L-2 | Partial | 36-hour triage matrix missing three edge cases: concurrent multi-scenario alert rollup, vendor-hosted institution-vs-vendor determination split, CIRCIA-only triggers without FFIEC trigger. |
| L-3 | Partial | spec §7 step 11 verifier procedure does not name `--strict` mode behavior under dual-algorithm transitional period (single-algorithm seal during dual-algorithm posture; algorithm not on declared posture list). |
| L-4 | Partial | SLSA provenance validation not in supply-chain "What examiners verify" list; mirror-registry section missing continuous-monitoring control on signature-validation discipline (re-signing pattern). |
| L-5 | Partial | Regulator-held public-key fingerprint rotation procedure referenced as "documented separately" but no pointer or stub exists. |

### Big Four model risk (Hiroshi Bergmann)

| ID | Severity | Item |
|---|---|---|
| H-1 | Partial (noted) | Audit procedure for confirming the institution's `gen_ai_parameters` schema covers the SR 11-7 reproducibility surface (decoding, sampler ID, model version, system prompt, retrieval context, intra-run dependencies). Recommended P-22-equivalent. |

## Round-11 close-out plan

The 26 items cluster into six focused work items:

1. **Tomás's mechanical gaps** (T-1..T-4) — copy-paste leftovers and missing test-vector files. Most critical for cross-implementation interoperability.
2. **Diane's SOC audit cascade** (D-1..D-11) — cadence, vocabulary, Section 4 additions, TSC homing, sample-report rows, P-22/P-23/P-24, user-entity, anomaly-template refinements.
3. **Priya's supply-chain operational** (P-1..P-5) — mirror retention, slsa-verifier MUST-tier, GPG recovery, Adversary I framing, cosign during examination.
4. **Liang's CSF + 36-hour + dual-algorithm** (L-1..L-5) — CSF four-attack table, triage edge cases, verifier dual-algorithm strict, SLSA in checklist, fingerprint rotation pointer.
5. **Marcus's examiner polish** (M-1..M-4) — quickstart consistency, sample-report adjacent PASS, IR scenario column, second worked example.
6. **Hiroshi's noted P-25** (H-1) — gen_ai_parameters schema check audit procedure.

Round 11 follows after closures with six fresh personas (different names from rounds 9 and 10) and the convergence reassessment.

## Trajectory

| Round | Total Gaps | Total Partials | Notes |
|---|---|---|---|
| Round 9 | 9 | 26 | First pass against v1.0-rework |
| Round 10 | 5 | 21 | After round-9 closure (26% net reduction) |
| Round 11 (target) | 0 | 0 | Convergence target |

If round 11 reaches 0/0, round 12 spawns six more fresh personas to confirm convergence holds across rounds.
