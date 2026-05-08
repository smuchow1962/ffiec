# Round-12 review — Big Four SOC Audit Engagement Partner

**Reviewer.** Patricia Zhang. Engagement Partner, Big Four firm. 24 years; APAC bank-tech SOC practice lead; recently rotated back to US engagements.
**Reading angle.** Auditability, testability, evidence-trail completeness, user-entity audience consumption, mechanical SOC procedure operability.
**Read at this round.** spec/chain-of-custody-v1.md (focus §7, §10); docs/audit-procedures.md (focus P-22 through P-25); docs/soc-pack/control-evidence-events.md; docs/soc-pack/section-4-template.md; docs/control-map/CUECs.md; docs/control-map/TSC-mapping.md; docs/anomaly-documentation-template.md; docs/user-entity-summary.md; docs/regulator-pack/sample-report.md.
**Stopping criterion.** 0 gaps / 0 partials.

---

## Headline

The SOC pack is operable as shipped. A SOC 2 Type II engagement team can pick up Section 4 starter, the audit-procedures menu, the operational-events schema, the anomaly-evaluation template, the CUEC catalogue, the TSC mapping, and the sample report; allocate procedures to criteria; pull evidence from the institution's operational-event stream; document anomalies on a structured template; and produce a defensible opinion. Every cross-reference I followed pointed at a real, named artifact. The vocabulary that previously drifted between CUEC text and spec text is now consistent at every cross-reference I tested. P-22 through P-25 — the four newer evidence-completeness procedures — are operable: each names the specific operational events to pull, the specific fields the SOC team reads, the disposition fork, and the escalation path. The sample verifier report's TSC-criterion appendix now covers all twelve §7 failure-record rows plus the mid-write-truncation row, with PASS and FAIL day examples adjacent inside the same examination period. I have one observation that does not rise to a finding (a small corpus-wide vocabulary residue in non-SOC-pack documents that the SOC pack itself no longer relies on).

**Verdict.** Pass for the SOC pack as a deliverable bundle. 0 gaps, 0 partials inside the SOC-pack scope. One non-SOC-pack hygiene observation noted below for the working group's reference, recorded as Observation rather than Finding because it does not affect a SOC engagement's mechanical operability.

---

## Findings table

| ID | Status | Location | What I tested | Result |
|---|---|---|---|---|
| F-1 | **PASS** | `docs/control-map/CUECs.md` CUEC-CRY-04 | Vocabulary check: does the CUEC name "key-fingerprint reconciliation" matching spec §10.1 and audit-procedures P-6? | CUEC-CRY-04 reads "**key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1)" with the correct field name `(tenant_id, key_version, key_fingerprint)` triple, the correct operational event `master.reconciliation_completed`, and the correct cross-reference to P-6. **Matches.** |
| F-2 | **PASS** | `docs/regulator-pack/sample-report.md` SOC team appendix | Coverage check: does the verifier-line → TSC-criterion mapping cover all twelve §7 steps plus the mid-write-truncation row? | All 12 spec §7 step rows present (steps 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12-cadence and 12-dev-mode as separate rows). Mid-write truncation row present. PASS day (2026-04-22) and FAIL day (2026-04-23) shown adjacent within the April 2026 examination period. Rotation day at 2026-04-15. The `--master-key` flag is shown in the verifier invocation with the institution-side disclosure cross-references. **Complete.** |
| F-3 | **PASS** | `docs/audit-procedures.md` P-22 through P-25 | Operability check: can a SOC team execute each procedure mechanically against the institution's operational-event stream? | Each procedure names: (a) the operational events to pull; (b) the specific fields read; (c) the disposition fork (operationally explained vs escalated); (d) the IR scenario for escalation; (e) the four-piece evidence package required for `key_fingerprint mismatch`. P-25 correctly frames `gen_ai_parameters` schema gaps as MRM-program findings rather than chain-integrity findings. **Operable.** |
| F-4 | **PASS** | `docs/soc-pack/section-4-template.md` | Completeness check: does the template's Procedures list cover the operational events that a SOC engagement attests to operating? | The Procedures bullet list covers: daily seal job; HSM PIN rotation; IKM rotation; **IKM retirement per §10.9 with `master_key.retired` event**; **key-fingerprint reconciliation per §10.1 with `master.reconciliation_completed`**; **mid-write-truncation recovery per IR Scenario 9 with `audit_file.truncation_detected`**; verifier validation and reproducible-build with SLSA provenance attestation; `--strict` for substantive SOC-evidence runs; IR per the playbook. The starter is now in lock-step with the operational-event schema. **Complete.** |
| F-5 | **PASS** | `docs/soc-pack/control-evidence-events.md` | Schema completeness for the operational events that P-22..P-25 consume | Schema covers `audit_file.truncation_detected` with `recovery_outcome ∈ {complete, partial, unrecoverable}`, `master_key.retired` with `chain_entries_referencing_remaining`, `master.reconciliation_completed` with `fingerprint_unmatched_count`, and `chain.verification_failure` with `step` field referencing spec §7 (1-12) and `detail` carrying both expected and recorded fingerprint for step-8 failures. The schema gives the SOC team and IR team mechanical mappings. **Complete.** |
| F-6 | **PASS** | `docs/anomaly-documentation-template.md` | Operability check: can the institution and the SOC team fill in a defensible record per anomaly type? | Template carries the anomaly-type enumeration matching the spec §7 step numbers. Higher-bar evidence requirements for `key_fingerprint mismatch` (four pieces), `hkdf_inputs_digest mismatch` (two acceptable explanations), and `audit_file_truncation_detected` (recovery outcome) are explicit. Worked example for the `key_fingerprint mismatch` four-piece package is template-quality and reduces first-time friction. **Operable.** |
| F-7 | **PASS** | `docs/control-map/TSC-mapping.md` | Coverage check: do the TSC subcriteria map cleanly to the chain primitives and to the new operational events? | CC6.1 carries the `master_key.retired` mapping. CC6.8 carries the `--strict` refusal of `dev_mode=true` and `kms_handle_uri = "plaintext-*"`. CC7.2 carries `audit_file.truncation_detected` and `master.reconciliation_completed`. CC8.1 carries the IKM-lifecycle events. PI1.1 / PI1.2 carry the chain-integrity headline. The headline-mapping paragraph correctly identifies the chain as a PI1.1 / PI1.2 / CC7.2 control with strong CC8.1 properties. **Complete.** |
| F-8 | **PASS** | `docs/user-entity-summary.md` | Audience check: does the user-entity audience receive a defensible summary that names the new operational events without requiring spec literacy? | Summary calls out `audit_file.truncation_detected` and `master_key.retired` as the two low-frequency, high-impact operational events the SOC report attests to, with plain-language descriptions of what each is and why it matters. The downstream-business-partner audience and the financial-statement-auditor audience both receive readable orientation. **Complete.** |

---

## Observation (non-finding)

**O-1.** Three documents outside the SOC pack still carry the pre-rework `session-key-id` vocabulary in user-facing prose: `docs/edge-and-federated-ai.md` line 27 ("Session-key-id reconciliation on edge devices..."), `docs/design/05-otlp-wire.md` line 32 (ASCII-art OTLP envelope showing `ffiec.chain.session_key_id: ".."`) and line 168 (SIEM correlation example listing `ffiec.chain.session_key_id` alongside `tenant_id` / `run_id`), and `web/content.js` line 1102 (master-key-exfiltration threat-model card text mentioning "Weekly session-key-id reconciliation"). Spec §12 change log line and `docs/design/10-glossary.md` `### session_key_id` entry are acceptable historical references (the change log explicitly explains what was dropped; the glossary entry preserves the term for readers encountering it in older artifacts). Design 05 line 67 is also acceptable because it explicitly says "post-rework v1.0; replaces `session_key_id`" — that is a deliberate cross-reference, not a drift.

**Why this is an Observation and not a Finding.** A SOC engagement does not consume `edge-and-federated-ai.md`, `design/05-otlp-wire.md`, or `web/content.js` as evidence. The SOC pack documents the engagement uses (CUECs.md, audit-procedures.md, control-evidence-events.md, section-4-template.md, TSC-mapping.md, anomaly-documentation-template.md, sample-report.md) all use the current `key_fingerprint` / `key-fingerprint reconciliation` vocabulary consistently. A SOC team would not encounter the residual drift during evidence collection, criterion mapping, or report drafting. The drift is a corpus-hygiene matter for the working group rather than an audit-defensibility matter.

**Suggested disposition.** The working group may schedule a corpus-wide vocabulary refresh in a future round; the four locations above are the residue. The fix is mechanical (search-and-replace from `session-key-id` / `session_key_id` to `key-fingerprint` / `key_fingerprint` with the surrounding sentence rewritten where needed). Not a SOC pack finding.

---

## What I tested mechanically (audit working-paper detail)

This is the section the engagement working paper would carry verbatim. I am recording it here so a future reviewer can follow exactly what I checked.

### CUEC-CRY-04 vocabulary verification

I read `docs/control-map/CUECs.md` lines 28-29. The text reads:

> CUEC-CRY-04 — The institution operates **key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1). The reconciliation matches every `(tenant_id, key_version, key_fingerprint)` triple observed on captured events against the institution's IKM-roster's expected fingerprint per `(tenant_id, key_version)` pair; mismatches are high-priority alerts (botched rotation, restored backup pointed at wrong tenant, cross-tenant key swap). The corresponding operational event is `master.reconciliation_completed` (audit-procedures P-6).

I then read spec §10.1 and audit-procedures P-6 to confirm the cross-references resolve. Both use the same vocabulary. The institution's procedure documentation, the spec, the CUEC, and the audit procedure are now in agreement. A SOC team querying the institution's procedures for "key-fingerprint reconciliation" will find the procedure and the `master.reconciliation_completed` event without translation. **Verified.**

### Sample-report TSC appendix coverage

I walked the table at `docs/regulator-pack/sample-report.md` lines 308-336 and matched each row against spec §7. Coverage map:

| Spec §7 step | Failure mode | Row present in appendix |
|---|---|---|
| 1 | format_version not supported | yes (lines 323) |
| 2 | header HKDF inputs do not match | yes (line 324) |
| 3 | header genesis_hash does not match v1 constant | yes (line 325) |
| 4 | cross-chain lift detected at seq N | yes (line 326) |
| 5 | format_version mismatch at seq N | yes (line 327) |
| 6 | chain link broken at seq N | yes (line 328) |
| 7 | unknown key_version | yes (line 329) |
| 8 | key_fingerprint mismatch | yes (line 319) |
| 9 | payload_hash MAC mismatch | yes (line 320) |
| 10 | merkle root mismatch | yes (line 321) |
| 11 | signature verification failed | yes (line 322) |
| 12 | cadence mismatch | yes (line 330) |
| 12 | dev-mode under --strict — refused | yes (line 331) |
| §4.1 | audit file ends mid-line | yes (line 332) |

PASS day at 2026-04-22 shown adjacent to FAIL day at 2026-04-23 (lines 220-261). Rotation day at 2026-04-15 in the per-day detail (lines 121-143). All within the April 2026 examination period (2026-04-01 to 2026-04-30). The `--master-key` flag is in the standard examiner invocation (lines 22-30) with the institution-side disclosure cross-references at lines 32 and the `--strict` SOC-engagement variant at lines 39-47. **Verified.**

### P-22 through P-25 operability

I walked each procedure as a SOC senior would on a Monday-morning evidence-pull.

**P-22 — `key_fingerprint mismatch` evidence completeness.** SOC team queries the institution's operational-event store for `chain.verification_failure` events with `step=8` for the period. For each, pulls the institution's anomaly-evaluation record (per the template) and confirms all four pieces: IKM-roster row identification, change-management approval for the correction, re-verification PASS on the corrected roster, reconciliation cross-check showing `fingerprint_unmatched_count=0`. Anomaly missing one or more pieces is documented as an unmitigated control gap. Mechanical, operable. The worked example in the anomaly template (lines 147-221) confirms a real anomaly produces a complete record.

**P-23 — `hkdf_inputs_digest mismatch` evidence completeness.** SOC team queries `chain.verification_failure` events with `step=2`. For each, confirms the root cause is one of two acceptable explanations (known SDK / verifier constants change documented in change management, or multi-region deployment-drift case with re-deployment to consistent constants). Mismatches without one of the two are escalated as format-construction defects under IR Scenario 1. Mechanical, operable.

**P-24 — `audit_file.truncation_detected` evidence completeness.** SOC team queries `audit_file.truncation_detected` operational events. For each, confirms the writer-process crash is identified in host-level monitoring AND `recovery_outcome` is documented (`complete`, `partial`, or `unrecoverable`). For `partial`, the recovered subset is restored and the un-recovered subset is documented as an unrecoverable gap with IR Scenario 9 disposition. For `unrecoverable`, the institution treats the gap as an integrity-control failure and starts the 36-hour clock per IR Scenario 9. Mechanical, operable.

**P-25 — `gen_ai_parameters` schema completeness for SR 11-7 reproducibility.** SOC team pulls representative chain entries with `gen_ai_parameters` populated; for chain entries representing model calls, confirms `gen_ai.request.model` and `gen_ai.response.model` are present (per spec §4.4 RECOMMENDED guidance — required for SR 11-7 reproducibility, since vendors silently re-route between model versions during outages). Confirms the institution's `gen_ai_parameters` schema covers the documented reproducibility surface (decoding parameters, sampler implementation identifier, system-prompt content/hash with prompt-version registry, retrieval context with RAG document IDs and content-hashes, intra-run data dependencies via `(run_id, seq)` of earlier chain entries). Crucially: the procedure correctly frames schema gaps as **MRM-program findings, not chain-integrity findings** (the chain integrity-binds whatever the institution puts in `gen_ai_parameters`; the chain does not specify what the institution must put there). This is the right separation of concerns and matches the user-entity-summary's "the chain proves the decision was made; the institution's MRM program validates correctness" framing. Mechanical, operable.

### Section 4 template completeness

I walked the Procedures bullet list at `docs/soc-pack/section-4-template.md` lines 53-61. The list covers:

- Daily seal job runs at UTC 00:00 + 60 minutes per tenant (default cadence).
- HSM PIN rotation on `[quarterly | semi-annual]` cadence.
- IKM rotation on `[institution-defined]` cadence with documented rotation procedure.
- **IKM retirement on `[institution-defined]` cadence per spec §10.9; retirement requires `chain_entries_referencing_remaining = 0` and is recorded as a `master_key.retired` operational event signed off by `[role]`.**
- **Key-fingerprint reconciliation on `[weekly]` cadence per spec §10.1; recorded as `master.reconciliation_completed` operational event with `fingerprint_unmatched_count` field.**
- **Mid-write-truncation recovery per IR Scenario 9; the SDK-local SQLite buffer and OTLP retention are the documented recovery sources; recovery outcome is recorded on the `audit_file.truncation_detected` operational event.**
- Verifier validation before each use; reproducible-build verification at least once per verifier release; SLSA provenance attestation validated against `slsa-verifier` at deployment time.
- Verifier runs on `[institution-defined]` cadence; substantive SOC-evidence runs are invoked with `--strict`; anomaly-evaluation runs are invoked without.
- Incident response per `docs/incident-response-playbook.md` (Scenarios 1-10 + 36-hour triage matrix).

The starter Section 4 description now names every operational event the SOC opinion attests to. The `[bracketed]` items are the institution-specific facts the institution fills in. The structure matches what an engagement team would expect to read in Section 4 of a SOC 2 Type II report. **Complete.**

---

## What I would have flagged in earlier rounds and no longer flag here

Round-12 readers may want to know what is no longer a concern in the SOC pack:

- **CUEC-CRY-04 vocabulary.** Now reads `key-fingerprint reconciliation` matching the spec, the audit procedure, and the operational event. The reconciliation cross-reference points at `master.reconciliation_completed`. A SOC team querying the institution's procedure documentation for the CUEC's name will find the institution's actual procedure without vocabulary translation.
- **Section 4 template Procedures list.** The IKM-retirement, key-fingerprint reconciliation, and mid-write-truncation procedures are now first-class bullet points alongside HSM PIN rotation and the daily seal job. The reproducible-build line carries SLSA provenance attestation.
- **Sample-report TSC appendix.** All twelve spec §7 steps and the truncation row are mapped to TSC subcriteria. PASS and FAIL days are shown adjacent in the same examination period for examiner training. The `--master-key` flag is in the invocation with the disclosure-procedure cross-references.
- **Anomaly template higher-bar evidence.** The four-piece package for `key_fingerprint mismatch` is explicit and the worked example reduces first-time friction. The two acceptable explanations for `hkdf_inputs_digest mismatch` and the three recovery outcomes for `audit_file_truncation_detected` are explicit.
- **TSC-mapping coverage of new operational events.** CC6.1 carries `master_key.retired`; CC6.8 carries `--strict` refusal of dev-mode and `plaintext-*` adapters; CC7.2 carries `audit_file.truncation_detected` and `master.reconciliation_completed`; CC8.1 carries the IKM-lifecycle events.
- **User-entity-summary surfacing of `audit_file.truncation_detected` and `master_key.retired`.** The downstream-audience document calls these out in plain language so a financial-statement auditor or downstream business partner reads them without spec literacy.

These were SOC-pack-affecting items; they are closed.

---

## Final stopping criterion check

- **Findings that block the SOC pack from operating mechanically: 0.**
- **Partials inside the SOC pack: 0.**
- **Observations recorded for working-group consideration: 1 (corpus-wide vocabulary residue in three non-SOC-pack documents; not affecting engagement operability).**

The SOC pack passes round-12 review at the 0/0 stopping criterion. The Observation is captured for the working group's next corpus-hygiene pass.

---

## Engagement-partner sign-off note

When I told my Asia-Pacific engagement teams what to look for in a SOC pack that gets out of their way, my list was: (1) the CUECs say what the institution actually does, in the institution's own vocabulary; (2) the audit procedures map to operational events the institution actually emits; (3) the Section 4 starter mentions every event we will see in the operational-event stream so the description is complete; (4) the anomaly template is fillable by the institution's control owner without requiring a cryptographer in the room; (5) the sample report shows what a passing day looks like AND what a failing day looks like, in the same period, so the examiner training is grounded; (6) the TSC mapping ties every verifier output line to a subcriterion so we can allocate evidence without arguing.

This pack does all six. The SOC engagement team picks it up and works.

— Patricia Zhang
