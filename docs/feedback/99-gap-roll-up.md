# 99 — Gap roll-up (round 9 — first review of v1.0-rework corpus)

> **What this doc is.** Cumulative finding inventory across the six round-9 reviewer personas. Stopping criterion: 0 gaps and 0 partials per role document. Round 9 is the first pass against the v1.0-rework spec; the cycle continues until convergence.

## Per-role counts

| Role | Persona | Answered | Partial | Gap |
|---|---|---|---|---|
| Big Four — Cryptographic | Dr. Yuki Tanaka | 1 | 3 | 1 |
| Big Four — Cyber Risk | Anneliese Vandermeer | 6 | 0 | 0 |
| Big Four — SOC Audit | Marcus O'Brien | 2 | 2 | 1 |
| Big Four — Model Risk | Dr. Pradeep Iyer | 0 | 4 | 1 |
| FFIEC — IT Examiner | Robert Chen | 0 | 2 | 3 |
| FFIEC — Cybersecurity Specialist | Sofía Reyes (18 findings, mapped from severity) | 0 | 15 | 3 |
| **Totals** | | **9** | **26** | **9** |

Sofía Reyes used a severity-graded format (High / Medium / Low) rather than Status. Mapped to the standard Status: High → Gap, Medium → Partial, Low / Low-Medium → Partial. Eighteen findings total.

**Convergence not reached.** Nine gaps and twenty-six partials require closure before round 10.

## Findings inventory by target document

### Spec (`spec/chain-of-custody-v1.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| S-1 | Yuki Q1 | Gap | `hkdf_inputs_digest` has two definitions: per-tenant in §3/§4.4/§7 (and the test vectors); tenant-free in §4.2 and design 03 §3.7. Two implementations following different sections produce non-interoperable seals. **Pick per-tenant; fix §4.2 and design 03.** |
| S-2 | Yuki Q2 | Partial | HKDF info separator `b"\|"` is not prefix-free; tenant_id with `\|` is admitted by spec. Add tenant_id grammar restriction in §3 (e.g. `[A-Za-z0-9._-]+`) or length-prefix the tenant_id in info. |
| S-3 | Yuki Q3 | Partial | `algorithm` field not bound into `sign_payload`. Add `{algorithm}\n` to §4.3 sign_payload now (before v1.1 PQC algorithm ships); recompute chain_vectors.json sign_payload bytes. |
| S-4 | Yuki Q4 minor | Partial | §4.4 should note that per-entry `algorithm` is forensic-only, parallel to existing wording for `mac_computed_at_utc` and `kms_handle_uri`. |
| S-5 | Yuki Q5 + Sofía F17 | Partial | §10 needs normative IKM-registry retention rule: IKM MUST be retained as long as any chain entry under that key_version is retained; institution's IKM-retirement procedure MUST be documented. |
| S-6 | Yuki Q5 | Partial | §4.1.1 Model B (HSM-mediated HKDF) assumes HSM capability that most PKCS#11 devices don't natively expose. Document the fallback (PRK transits SDK process briefly) and what Model B actually requires. |
| S-7 | Yuki Q5 | Partial | §4.2.2 / §10 needs normative note covering "rotation crossing the seal boundary" (late-binding events under old IKM appearing in day-after seal). |
| S-8 | Pradeep Q2 | Partial | §5 canonical-form exclusion list excludes linkage fields (`parent_run_id`, `parent_seq`, `dag_parents`). Linkage fields should be IN the canonical bytes (integrity-bound under the MAC). One-line spec change. |
| S-9 | Sofía F3 | Partial | §10.2 add `audit_file.truncation_detected` operational event. |
| S-10 | Pradeep Q1 | Partial | §4.4 `gen_ai_parameters` expanded content guidance: name decoding parameters, sampler ID, model version (request + response), system prompt content/hash, retrieval context, intra-run data dependencies. Make `gen_ai.request.model` + `gen_ai.response.model` required (not optional) for chain entries representing model calls. |

### Design 02 chain construction (`docs/design/02-chain-construction.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| D02-1 | Pradeep Q2 | Partial | §8.1.2 (shared run_id handoff) lacks normative handoff event schema; verifier behavior on missing/malformed handoff event undefined. |

### Design 03 merkle seal (`docs/design/03-merkle-seal.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| D03-1 | Yuki Q1 | Gap | §3.7 `hkdf_inputs_digest` formula needs to change to per-tenant variant (paired with S-1). |

### Design 04 HSM custody (`docs/design/04-hsm-custody.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| D04-1 | Yuki Q3 | Partial | §4 sign_payload definition needs `algorithm` line added (paired with S-3). |
| D04-2 | Pradeep Q3 | Partial | §3.1 IKM rotation discussion uses old `master_version` term in places; align vocabulary with `key_version`. |

### Design 07 verifier design (`docs/design/07-verifier-design.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| D07-1 | Pradeep Q2 | Partial | §7 verifier procedure needs DAG topology resolution step (call it step 13): every `(run_id, seq)` named in a `dag_parents` attribute must resolve to an existing chain entry. Optional under default; reportable as anomaly. |
| D07-2 | Robert Q2 | Partial | §5.1 PDF report structure does not show new fields (`key_versions`, `kms_handle_uri`, `dev_mode`, `master_key_rotation_observed` anomaly kind). |

### Test vectors (`spec/test-vectors/chain_vectors.json`)

| ID | Source | Severity | Item |
|---|---|---|---|
| TV-1 | Yuki Q3 | Partial | sign_payload bytes need recompute after `algorithm` is added to the format (paired with S-3). |

### Audit procedures (`docs/audit-procedures.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| AP-1 | Marcus Q2 | Partial | P-3 fourth bullet (per-tenant determinism replay) needs procedural detail: who runs the replay, what is the comparison target. |

### Operational events (`docs/soc-pack/control-evidence-events.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| OE-1 | Marcus Q3 minor | Partial | `seal.job_completed` example uses `master_version` field; spec uses `key_version` (per entry) and `key_versions` (seal record list). One sentence to clarify the operational/custodian-side label vs the per-entry integer. |
| OE-2 | Sofía F3 | Partial | Add `audit_file.truncation_detected` event (paired with S-9). |

### Anomaly-documentation template (`docs/anomaly-documentation-template.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| ADT-1 | Marcus Q5 | Gap | Add four new failure modes to enum + severity table: `key_fingerprint mismatch` (Critical), `unknown_key_version` (High), `format_version mismatch` (Medium), `hkdf_inputs_digest mismatch` (Critical). Extend "operationally explained" section to set a different bar for fingerprint mismatches (documented IKM-roster correction with change-management approval, re-run of affected period). |

### Examiner quickstart (`docs/examiner-quickstart.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| EQ-1 | Robert Q1 | Gap | "Common failure modes" 6-row table replaced with 12-row table using verifier's actual reason strings, one per spec §7 step. Add explicit callouts: `key_fingerprint mismatch` is the new "stop and call the bank" finding; `audit file ends mid-line` is operational, not tampering. |

### Sample report (`docs/regulator-pack/sample-report.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| SR-1 | Robert Q2 | Partial | Add a rotation day to the 30-day excerpt (`key_versions: [3, 4]`, `master_key_rotation_observed` anomaly with explanation). |
| SR-2 | Robert Q2 | Partial | Add per-day `key_fingerprint(s) observed` count, `kms_handle_uri` line per day. |
| SR-3 | Robert Q2 | Partial | Add a paired short FAIL-day sample report showing per-day block shape, failure-record structure (with `step`, `reason`, `seq`, `expected_fingerprint`, `recorded_fingerprint`), and how cover-page summary changes. |
| SR-4 | Robert Q2 | Partial | Methodology section page N+2 update to reference spec §7 twelve-step procedure by step number. |
| SR-5 | Marcus Q4 | Partial | Add `--strict` invocation example; SOC team's verifier-line → TSC criterion appendix. |

### Finding language (`docs/regulator-pack/finding-language.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| FL-1 | Robert Q3 | Gap | Add finding paragraphs for: `key_fingerprint mismatch`, `unknown_key_version`, `format_version not supported`, `audit file ends mid-line`, `cross-chain lift detected`, `cadence mismatch`. |
| FL-2 | Robert Q3 | Gap | Replace severity table with twelve-row table keyed by spec §7 step number, so examiner can map verifier output (`{"step": N, "reason": "..."}`) directly to a row. |

### Handbook mapping (`docs/regulator-pack/handbook-mapping.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| HM-1 | Robert Q4 | Partial | II.C.13 cryptographic controls: add coverage of per-entry `key_fingerprint` identity provenance (cryptographic control beyond "HMAC + Ed25519 in HSM custody"). |
| HM-2 | Robert Q4 | Partial | II.C.10 logging: add coverage of weekly key-fingerprint reconciliation per spec §10.1. |
| HM-3 | Robert Q4 | Partial | II.C.13: add IKM minimum length (§10.6, RFC 4868) and software-adapter compile-time exclusion (§10.7). |
| HM-4 | Robert Q4 | Partial | II.E change management: add `format_version` change-management coverage. |

### IR playbook (`docs/incident-response-playbook.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| IR-1 | Sofía F1 + Robert Q5 | Gap | Add Scenario 7: `key_fingerprint mismatch detected at verification`. Triage tree: rotation in flight? backup restored? tenant row restored from different period? Most often P2 ops misconfig, not P1 compromise. Tie-in to §10.1 weekly reconciliation. |
| IR-2 | Sofía F2 | Gap | Add Scenario 8: `unknown_key_version`. Three plausible root causes (decommissioned IKM, provisioning gap, tampered key_version). Different IR triage per cause. |
| IR-3 | Sofía F3 | Partial | Add Scenario 9: `audit_file.truncation_detected`. Mostly availability event; recover dropped event from SDK SQLite buffer. |
| IR-4 | Sofía F8 | Gap | Add 36-hour-applicability triage matrix per chain-detected event type; clock-start trigger lines per scenario. |
| IR-5 | Sofía F16 | Partial | Add Scenario 10: backup-integrity verifier failure (RC.RP-04 + RC.CO-03 evidence callouts). |

### Threat model (`docs/design/09-threat-model.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| TM-1 | Sofía F4 | Partial | §2.7 + §3.1.2: vocabulary refresh `session_key_id` → `key_fingerprint` consistent with spec §10.1. |
| TM-2 | Sofía F5 | Partial | Add residual-risk row(s) for per-tenant HKDF binding and `expected_prev_hash` defence (closed by spec §4.1 inviolate properties #1 and #8). |
| TM-3 | Sofía F14 | Partial | Add Adversary I: examiner-side tooling (mirror of `00-overview.md` §2.4 attack 4). |
| TM-4 | Sofía F15 | Partial | Add edge / federated residual risk (physical attacker on edge device with secure-enclave attestation). |
| TM-5 | Sofía F17 | Partial | IKM-registry retention obligation explicit in residual-risk register. |
| TM-6 | Sofía F18 | Partial | §2.8 Adversary H: add dual-algorithm transitional-period guidance (PQC coexistence, verifier preference). |
| TM-7 | Sofía F5 | Partial | R8 (session key leakage to logs): add bound that the public `key_fingerprint` is not offline-grindable due to 32-byte IKM minimum (RFC 4868). |

### CSF 2.0 mapping (`docs/regulator-pack/CSF-2.0.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| CSF-1 | Sofía F6 | Partial | Add GV.OC and GV.RM evidence binding: one-page "GV alignment" doc the institution can copy into its policy framework. |
| CSF-2 | Sofía F7 | Partial | Annotate operational events (spec §10.2) against CSF subcategories at finer granularity (DE.CM-09 vs DE.AE-02 vs DE.AE-03). |

### Supply chain (`docs/supply-chain.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| SC-1 | Sofía F9 | Partial | Cosign-key recovery procedure (revocation, institution notification, new-key authentication chain). |
| SC-2 | Sofía F10 | Partial | GPG-key rotation cadence. |
| SC-3 | Sofía F11 | Partial | Claim SLSA level (build pipeline meets SLSA L3 properties). |
| SC-4 | Sofía F12 | Partial | Spec text + conformance corpus integrity binding in per-release artifact table; sign both. |
| SC-5 | Sofía F13 | Partial | Container mirror-registry signature handling (institution mirror should preserve the cosign signature; if mirror re-signs, document the trust-path bridge). |

### MRM committee brief (`docs/MRM-COMMITTEE-BRIEF.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| MRM-1 | Pradeep Q3 | Gap | Extend §"The four primitives (short version)" with rework defensive primitives (per-tenant HKDF binding, key_fingerprint identity check, server-side Merkle+HSM placement, MAC-IS-payload_hash, compile-time dev-adapter exclusion, constant-time compare). Update §"What the MRM Committee should ask" item 4 to reference P-6 reconciliation report. |

### Customer dispute procedures (`docs/customer-dispute-procedures.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| CDP-1 | Pradeep Q4 | Partial | Add ECOA / adverse-action `gen_ai_parameters` reproducibility note; address whether the bank's compliance translation is itself a chain entry (parent of the AI's response). |
| CDP-2 | Pradeep Q4 | Partial | Address customer-side IKM access path (under protective order vs HSM-mediated verification). |

### AI policy alignment (`docs/regulator-pack/ai-policy-alignment.md`)

| ID | Source | Severity | Item |
|---|---|---|---|
| AI-1 | Pradeep Q5 | Partial | Treasury AI RMF section: enumerate Treasury RMF expectations (model inventory, lifecycle, third-party governance, ongoing monitoring) and where the chain composes. |
| AI-2 | Pradeep Q5 | Partial | EU AI Act: add Article 26 (deployer obligations); strengthen Article 14 (effective oversight). |
| AI-3 | Pradeep Q5 | Partial | DORA: add Article 6 (ICT risk-management framework) and Article 8 (identification of ICT-supported business functions). |
| AI-4 | Pradeep Q5 | Partial | Add OCC emerging supervisory posture on AI under existing SR 11-7 / 2011-12 framework. |

### Workflow doc (NEW)

| ID | Source | Severity | Item |
|---|---|---|---|
| WF-1 | Robert Q5 | Gap | Create new workflow doc: "From `chain.verification_failure` JSON record to documented examination response." Names the path: failure record → finding-language paragraph → IR playbook scenario → reconciliation event lookup → institution response template. ~60 seconds of doc to read. |

## Round-10 close-out plan

The findings cluster into seven focused work items. Round 10 follows after each is closed.

1. **Spec consistency fixes** (S-1, S-2, S-3, S-4, S-8 → spec §3, §4.2, §4.3, §4.4, §5; D03-1 → design 03; D04-1 → design 04). Plus regenerate `chain_vectors.json` with the new sign_payload format.
2. **Spec operational additions** (S-5, S-6, S-7, S-9, S-10 → spec §10 IKM retention; §4.1.1 Model B honesty; §4.2 rotation-crossing-seal note; §10.2 truncation event; §4.4 gen_ai_parameters guidance).
3. **Examiner-facing docs cascade** (EQ-1, FL-1, FL-2, SR-1..5, HM-1..4, WF-1 → examiner-quickstart, finding-language, sample-report, handbook-mapping, new workflow doc). Single editing pass replacing four-failure vocabulary with twelve-step vocabulary.
4. **IR playbook + threat model + CSF expansion** (IR-1..5, TM-1..7, CSF-1..2 → 4 new IR scenarios + 36-hour triage matrix; vocabulary refresh + 6 new residual-risk rows + 2 new adversaries; GV alignment + operational-events binding).
5. **Anomaly template + audit procedures + operational events** (ADT-1, AP-1, OE-1, OE-2, D02-1, D07-1, D07-2 → enum updates; P-3 procedural detail; vocabulary clarification; truncation event; handoff event schema; verifier step 13 DAG topology; PDF report new fields).
6. **Supply chain hardening** (SC-1..5 → cosign-key recovery; GPG rotation cadence; SLSA-L3 claim; corpus integrity binding; container mirror signature handling).
7. **MRM + dispute + AI policy expansion** (MRM-1, CDP-1..2, AI-1..4 → committee brief defensive primitives; ECOA chaining + IKM access; Treasury RMF + EU AI Act + DORA + OCC supervisory posture).

## Stopping criterion

**0 gaps and 0 partials per role document, across at least three consecutive convergence-attempt rounds with fresh personas each time.** Round 10 closes the items above; round 11 spawns six fresh personas to re-assess; round 12 spawns six more if any remain. Convergence at round 10 (single-pass) would be unusual; rounds 11-13 are the realistic target.
