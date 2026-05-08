# 99 — Gap roll-up (round 12)

> **What this doc is.** Cumulative finding inventory across the six round-12 reviewer personas. Round 11 closed 4 gaps + 17 partials; round 12 assesses whether the closures held and whether new finds surfaced.

## Per-role counts

| Role | Persona | Status |
|---|---|---|
| Big Four — Cryptographic | Dr. Klaus Reinhardt | 3 Gap, 2 Partial |
| Big Four — Cyber Risk | Yuki Yamamoto-Wallace | **0/0 ✓** |
| Big Four — SOC Audit | Patricia Zhang | **0/0 within scope ✓** (3 non-finding observations) |
| Big Four — Model Risk | Dr. Sandeep Kumar Reddy | **0/0 ✓** (3 minors, 4 observations) |
| FFIEC — IT Examiner | Janet Buckley | 2 examiner-severity + ~14 MINOR/NOTE |
| FFIEC — Cybersecurity Specialist | Mateusz Kowalski | 0 Gap, 5 Partial |
| **Subtotal (strict)** | | **3 Gap, ~28 Partial** |

**Three roles converged.** Up from two in round 11.

**The gaps cluster in two areas:** Klaus's dual-algorithm extension (the round-11 dual-algo work was incomplete — spec language added but seal record schema, test vectors, design 03/04/07 not fully extended). Mateusz's trust-anchor lifecycle items extend the round-11 cold-DR-key + Adversary I work into institution-side operational evidence.

## Findings inventory

### Cryptographic engineering (Klaus Reinhardt)

| ID | Severity | Item |
|---|---|---|
| K-1 | Gap | Spec §4.2 `signatures` list + §7 step 11 cases (a)-(e) don't specify whether each entry's signature covers its own algorithm-bound `sign_payload` (Variant B, correct) or all entries cover one primary-algorithm payload (Variant A, breaks algorithm-confusion defense). Case (e) reasoning only holds under Variant B. |
| K-2 | Gap | No corpus fixture or negative tests for dual-algorithm dispatch. Need `015-dual-algorithm-cosigned-seal/` + N017 (partial-coverage seal), N018 (algorithm-not-in-posture), N019 (one-valid-one-invalid). |
| K-3 | Gap | Design 03 §3.7 schema table missing `signatures` row; design 04 §3.2 still says "one signing keypair per tenant" (stale under dual-algorithm posture). |
| K-4 | Partial | Design 02 §10 contains sibling `### 9.1 Closed findings from the Herald upgrade audit` (should be `### 10.1`) — round-11 §8/§9 renumbering missed §10's subsection. |
| K-5 | Partial | Design 07 §4.3 sign_payload single-algorithm pseudocode correct; dual-algorithm extension missing (gated on K-1 disposition). |

### FFIEC cybersecurity specialist (Mateusz Kowalski)

| ID | Severity | Item |
|---|---|---|
| M-1 | Partial | Adversary I rotation-reception procedure needs institution-defined operational events (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`) for sample-testing, plus an IR sub-variant for reception-failure (notice arrives but validation fails). |
| M-2 | Partial | Cold-DR-key consumption needs (a) institution-side trust-anchor cache row in supply-chain `Trust anchors` table (cold-DR public-key fingerprint cached); (b) annual dry-run-attestation consumption step in institution's release-validation procedure (parse `KEY-DR-DRYRUN-{year}.asc`); (c) Scenario 11 in IR playbook for project-side trust-anchor degradation (cold-DR key cannot be reached during dual-compromise). |
| M-3 | Partial | Spec §7 step 11 case (e) co-signed-seal-failure needs a named IR scenario with a three-branch clock-start determination tree: (i) attributable to published algorithm break (regulator coordination, no clock); (ii) attributable to per-algorithm signing-key compromise (Scenario 4 starts clock); (iii) under investigation (clock starts at investigation-conclusion determination). |
| M-4 | Partial | CSF-2.0 operational-events binding table needs four additional rows: `regulator_fingerprint.rotation_received` → DE.CM-09 + GV.SC-04; `regulator_fingerprint.rotation_validated` → DE.AE-03; `regulator_fingerprint.installed` → CC8.1 + GV.SC-04; cold-DR `KEY-DR-DRYRUN-{year}.asc` consumption → GV.SC-04. |
| M-5 | Partial | 36-hour matrix needs a dual-algorithm institution-vs-regulator verifier-version timing-overlap edge case: when the institution upgrades to a verifier supporting algorithm Y but the regulator's verifier still only supports algorithm X, and a co-signed seal fails on the regulator's verifier (which sees only the X-algorithm half), the institution and regulator may produce different verification results from the same ledger snapshot. The matrix should name the disposition: institution coordinates with regulator on verifier-version compatibility before submitting evidence. |

### FFIEC IT examiner (Janet Buckley)

| ID | Severity | Item |
|---|---|---|
| J-1 | Examiner-severity | Training Module 3 invocation missing `--master-key`. New examiner trained on Module 3 will produce a structural-only verification on first real examination. |
| J-2 | Examiner-severity | Spec §7 case (e) needs a cross-reference to finding-language severity. The spec says non-strict is PASS-WITH-ANOMALY but finding-language says case (e) is Severe regardless of bracket. Junior examiner could write up case (e) as Observation. |
| J-3..J-16 | MINOR/NOTE | ~14 minor clarity edits: Windows ACL note for IKM file mode, FAIL-example block ordering, JSON snippet pairing, repeat-finding examples for new step 2/3/5/6/9 paragraphs, Module 5 row coverage for the two non-Severe dual-algorithm sub-cases, deployment-package sample command lines missing --master-key. |

### Big Four SOC audit (Patricia Zhang) — non-finding observations

| ID | Severity | Item |
|---|---|---|
| P-OBS-1 | (non-finding) | `session-key-id` / `session_key_id` vocabulary residue in three non-SOC-pack docs: `docs/edge-and-federated-ai.md` line 27, `docs/design/05-otlp-wire.md` lines 32 and 168 (line 67 acceptable — explicit cross-reference), `web/content.js` line 1102. Not a SOC-pack finding (these docs aren't SOC-consumed); recorded for working-group corpus-hygiene pass. |

### Big Four model risk (Sandeep Kumar Reddy) — minors

| ID | Severity | Item |
|---|---|---|
| S-1 | Minor | P-25 sampling stratification — name the stratification dimension (model-version, decision-class, customer-impact tier) the SOC team uses when sampling chain entries for `gen_ai_parameters` schema completeness testing. |
| S-2 | Minor | MRM-COMMITTEE-BRIEF question 3 should call out model-inventory composition (the chain-captured per-model decision count is one input to the institution's MRM model-inventory completeness check). |
| S-3 | Minor | Customer-dispute procedure should add MRM committee variation-threshold guidance — what variation across re-runs constitutes "significant" for the institution's reproduction posture. |

## Round-13 close-out plan

The 28 remaining items cluster into four focused work items:

1. **Klaus's dual-algorithm extension** (K-1..K-5) — Specify Variant B per-algorithm sign_payload; extend chain_vectors.json + add 015 case + N017/N018/N019; update design 03 §3.7 schema + design 04 §3.2 dual-keypair text + design 07 §4.3 dual-algo pseudocode; fix §10/§9.1 renumbering.
2. **Mateusz's trust-anchor lifecycle** (M-1..M-5) — Three new operational events; cold-DR consumption procedure; IR Scenario 11; case (e) IR scenario with three-branch clock-start; dual-algorithm verifier-version-overlap edge case.
3. **Janet's examiner polish** (J-1..J-16) — Training Module 3 + deployment-package --master-key; spec §7 case (e) severity cross-reference; ~14 minor clarity edits.
4. **Patricia's vocab residue + Sandeep's minors** (P-OBS-1, S-1..S-3) — Corpus-wide vocab refresh on edge-and-federated-ai, design 05, web/content.js; sampling stratification; MRM brief addition; variation-threshold.

Round 13 follows after closures with six fresh personas (different from rounds 9-12) and the convergence reassessment.

## Trajectory

| Round | Gaps | Partials | Roles converged |
|---|---|---|---|
| Round 9 | 9 | 26 | 0 of 6 |
| Round 10 | 5 | 21 | 0 of 6 |
| Round 11 | 4 | 17 | 2 of 6 |
| Round 12 | 3 | ~25 | **3 of 6** |
| Round 13 (target) | 0 | 0 | 6 of 6 |

If round 13 reaches 6 of 6, round 14 confirms convergence with six more fresh personas.
