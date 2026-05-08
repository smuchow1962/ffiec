# 99 — Gap roll-up (round 11 — second close-out review)

> **What this doc is.** Cumulative finding inventory across the six round-11 reviewer personas. Round 10 closed 5 gaps + 21 partials; round 11 assesses whether the closures held and whether new finds surfaced.

## Per-role counts

| Role | Persona | Answered | Partial | Gap |
|---|---|---|---|---|
| Big Four — Cryptographic | Dr. Adaeze Okonkwo | 6 | 2 | 0 |
| Big Four — Cyber Risk | Eitan Schreiber | 6 | 0 | 0 ✓ |
| Big Four — SOC Audit | Camille Lefèvre | (strengths) | 2 | 1 |
| Big Four — Model Risk | Dr. Felipe Caro | 5 | 0 | 0 ✓ |
| FFIEC — IT Examiner | Devraj Patel | (strengths) | 6 | 3 |
| FFIEC — Cybersecurity Specialist | Astrid Holm | 9 | 7 | 0 |
| **Totals** | | — | **17** | **4** |

**Round 11 trajectory.** Round 10: 5 gaps + 21 partials. Round 11: 4 gaps + 17 partials. Net reduction: 26 → 21 (19% reduction).

**Two roles converged to 0/0** (Eitan, Felipe). Four roles still have items.

**Convergence not reached.** Round 12 will close the items below.

## Findings inventory

### Cryptographic engineering (Adaeze Okonkwo)

| ID | Severity | Item |
|---|---|---|
| AD-1 | Partial | Dual-algorithm seal encoding underspecified — spec §7 step 11 contemplates "both signatures present" but §4.2 seal record only admits one `algorithm` + `signature`. Specify either two seal records per `(tenant_id, seal_date)` or a `signatures` array. |
| AD-2 | Partial | design 02 §8/§8.1 numbering collision — top-level "## 8.1 Multi-process run semantics" collides with subsection "### 8.1 Verifier-side failures" inside §8. Round 10 fix addressed §4 only. |
| AD-3 (note) | (info) | design 08 §5.4.3 case-014 prose still uses pre-rework `master_version` vocabulary; chain_vectors.json bytes are correct. |

### Big Four cyber risk (Eitan Schreiber)

| ID | Status |
|---|---|
| All six focus areas | **Answered. Stopping criterion met.** Four documentation refinements flagged as non-gaps: default mirror-reconciliation cadence, cold-DR-key annual health check, "delayed determination" / re-classification rule in 36-hour matrix, `slsa-verifier` version-pinning posture. |

### Big Four SOC audit (Camille Lefèvre)

| ID | Severity | Item |
|---|---|---|
| F-1 | Gap (✅ CLOSED in this turn) | CUEC-CRY-04 vocabulary fix from round 10 didn't actually land in the shipped artifact. Verified and re-applied. |
| F-2 | Partial | Sample-report TSC appendix missing FAILURE-RECORD rows for §7 steps 1, 3, 4, 5, 6, 7, 12, and audit-file-truncation FAIL. SOC teams have to invent missing rows ad-hoc. |
| F-3 | Partial (✅ CLOSED in this turn) | "session-key-id" persists in audit-committee-summary, at-scale-operations, legal-disclosure, design 06, design 05, design 08, operator-guide, IR playbook. Corpus-wide vocabulary refresh applied. |

### Big Four model risk (Felipe Caro)

| ID | Status |
|---|---|
| All five focus areas | **Answered. Stopping criterion met.** No partials, no gaps. |

### FFIEC IT examiner (Devraj Patel)

| ID | Severity | Item |
|---|---|---|
| D-1 | Gap | Finding-language paragraphs missing for steps 2, 3, 5, 6, 9. Generic "Chain hash mismatch" paragraph is a leftover from before the step-keyed convention. |
| D-2 | Gap | Dual-algorithm transitional period verifier reason strings ("partial-coverage seal", "algorithm not on declared posture list", co-signed seal) appear in spec §7 but in zero downstream examiner artifacts (quickstart, finding-language, response-workflow, handbook-mapping II.E). |
| D-3 | Gap | "12-row" table is actually 14 rows — counting/framing inconsistency in both quickstart and finding-language. |
| D-4 | Partial | Sample-report PASS/FAIL paired example uses dates outside the 2026-04 examination range (2026-05-07 and 2026-05-08); cover-page summary inconsistency only partially covered by disclaimer. |
| D-5 | Partial | Examination-response-workflow doc lacks step-mapping rows and worked-example coverage for dual-algorithm sub-cases. |
| D-6 | Partial | `--master-key` flag never mentioned in quickstart or sample-report invocations despite being load-bearing under `--strict`. |
| D-7 | Partial | Handbook-mapping II.E covers `format_version` change management thoroughly but not algorithm change management for the dual-algorithm transitional period. |
| D-8 | Partial | Finding-language severity table has no IR-scenario column; examiner has to cross-reference back to quickstart. |
| D-9 | Partial | Examiner-training Module 5 "Common patterns" table predates the step-keyed taxonomy. |

### FFIEC cybersecurity specialist (Astrid Holm)

| ID | Severity | Item |
|---|---|---|
| A-1 | Partial | CSF operational-events binding table lacks rows for `mirror.reconciliation_completed` and SLSA-verifier validation log → GV.SC-04. |
| A-2 | Partial | 36-hour triage edge cases don't name federal-regulator routing per institution charter (credit unions to NCUA, state member banks to Fed, etc.). |
| A-3 | Partial | Spec §7 step 11 missing case (e): both signatures present, one valid + one invalid. |
| A-4 | Partial | Dual-algorithm examiner working-paper convention (recording both algorithm validation results) unstated. |
| A-5 | Partial | SLSA-aware tool-choice documentation requirement implicit, not normative — examiner needs to know which tool the institution chose and its trust anchor. |
| A-6 | Partial | Cold-disaster-recovery key's own rotation cadence, key-holder lifecycle, and dry-run cadence undocumented. |
| A-7 | Partial | Regulator-held fingerprint rotation pointer: (a) dependence on unpublished `regulator-pack/regulator-procedures.md` should be named as a bounding residual, (b) institution-side reception-and-validation procedure for the rotation event is unstated. |

## Round-12 close-out plan

The 19 remaining items (after F-1 and F-3 closed in this turn) cluster into five focused work items:

1. **Cryptographic refinements** (AD-1 dual-algorithm seal encoding, AD-2 §8 numbering collision, AD-3 design 08 case-014 vocab)
2. **Examiner-facing dual-algorithm cascade** (D-2, D-5, D-7, A-3, A-4 — propagate dual-algorithm from spec §7 to quickstart, finding-language, workflow, handbook-mapping II.E)
3. **Examiner-facing polish** (D-1 missing finding paragraphs, D-3 row-count framing, D-4 dates-in-range, D-6 --master-key flag, D-8 IR column, D-9 training Module 5)
4. **SOC + CSF additions** (F-2 TSC appendix rows, A-1 CSF operational events, A-5 SLSA tool-choice normative)
5. **Supply chain + threat model refinements** (A-2 federal-regulator routing, A-6 cold-DR-key cadence, A-7 regulator-held fingerprint reception)

Round 12 follows after closures with six fresh personas (different from rounds 9, 10, 11) and the convergence reassessment.

## Trajectory

| Round | Total Gaps | Total Partials | Notes |
|---|---|---|---|
| Round 9 | 9 | 26 | First pass against v1.0-rework |
| Round 10 | 5 | 21 | After round-9 closure (26% reduction) |
| Round 11 | 4 | 17 | After round-10 closure (19% reduction); 2 roles converged |
| Round 12 (target) | 0 | 0 | Convergence target |

If round 12 reaches 0/0, round 13 confirms convergence with six more fresh personas.
