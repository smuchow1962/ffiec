# Round 14 — Big Four SOC Audit review

**Reviewer.** Marisol Hernández, Senior Manager, SOC 1 / SOC 2 practice (Mexico City). Seventeen years in SOC engagements across Latin American and US bank-tech market — most of it on service organizations whose customers are US regional banks, processors, and cross-border fintechs. I read SOC reports the way the issuer's external auditor reads them: Section 4 description first, CUECs second, evidence shape third, and then I work backward to the test of operating effectiveness.

**Scope of this read.** Spec §7 (especially the new step 12a `gen_ai_model_identifier_missing` check) and §10 (the new `regulator_fingerprint.*` and `cold_dr.dryrun_attestation_consumed` operational events). `audit-procedures.md` P-22 through P-27 — the new P-26 model-inventory composition and P-27 customer-dispute reproduction procedures. `soc-pack/control-evidence-events.md` (new schema rows for the regulator-fingerprint lifecycle and the cold-DR-key dry-run consumption). `soc-pack/section-4-template.md`. `control-map/CUECs.md` (CUEC-CRY-04 vocabulary). `control-map/TSC-mapping.md`. `anomaly-documentation-template.md`. `user-entity-summary.md`. `regulator-pack/sample-report.md`.

**Stopping criterion.** 0/0. I read everything in scope; I report what I find. No prior-round context.

---

## 1. Headline

The package is testable. That is the bar a SOC engagement actually needs and most service-organization descriptions miss. Every control claim in this stack ties to (a) a specific operational event with a defined schema, (b) a verifier output line with a named step number, and (c) a sample-test procedure with stated evidence. The cross-walk in `sample-report.md`'s SOC team appendix — verifier line → TSC criterion — is the cleanest example of this I have seen in a SOC-supporting artifact in this asset class. An engagement team can plan the testing program from this stack alone.

The four pieces that make the package land for SOC purposes:

- **The verifier output is mechanically consumable.** Every FAIL row carries `step`, `reason`, `run_id`, `seq`, `key_version`, `expected_fingerprint`, `recorded_fingerprint`. A SOC team querying the institution's log store can pull `chain.verification_failure` events filtered by `fields.step` and produce a population for sampling within the engagement's normal evidence-collection cadence. No bespoke parsing.
- **The CUEC list is enumerated and tiered.** CUEC-IAM-01 through CUEC-CFG-03 with Tier 1 / Tier 2 / Tier 3 ramp-up, plus a one-line mapping from each chain primitive to the relevant CUECs. Section 4F is a fill-in-the-blanks problem for the institution rather than an architecture problem.
- **The anomaly evaluation has a written floor.** The "Higher bar for `key_fingerprint mismatch`" in the anomaly-documentation-template — four required pieces of evidence (IKM-roster row, change-management approval, re-verification, reconciliation cross-check) — is what a SOC engagement actually needs to dispose of an anomaly without a finding. Same for the two acceptable explanations on `hkdf_inputs_digest mismatch` and the three-value `recovery_outcome` enum on `audit_file_truncation_detected`. The template tells the institution what evidence to assemble and tells the SOC team what to test against. That alignment is the difference between a clean engagement and a multi-week back-and-forth.
- **The vocabulary is consistent across the stack.** `key_version` (per-entry, integer) and `key_versions` (seal-record, list) is reconciled in CUEC-CRY-04, in P-6, in the seal-job event schema, and in the verifier sample report. The `master_version`-deprecation footnote in `control-evidence-events.md` is exactly the migration note SOC practitioners need when comparing a Type II report's prior period to its current period.

The package reads like the authors have been on the receiving end of SOC engagement audit-program reviews. Most service-organization descriptions in AI-adjacent tooling fail their first SOC engagement because the description names controls the system does not actually emit evidence for. This stack does not have that failure mode.

---

## 2. Specific testability findings

These are the items where the SOC team's testing program lands cleanly without requiring institution-side custom evidence builds.

### 2.1 Step 12a `gen_ai_model_identifier_missing` is properly framed as control-completeness, not chain-integrity

Spec §7 step 12a fires only on entries that carry any `gen_ai.*` attribute, leaves audit-only and tool-call entries alone, and reports as PASS-WITH-ANOMALY under non-strict / FAIL under `--strict`. That distinction matters for the SOC opinion. A chain-integrity finding cascades into PI1.1 / PI1.2 modification of the SOC opinion. A control-completeness finding under SR 11-7 reproducibility is a separate matter — it is the institution's MRM program that has not populated `gen_ai.request.model` and `gen_ai.response.model` consistently, NOT the chain that has been tampered with. The spec text in step 12a names the distinction in plain language ("control-completeness for SR 11-7 reproducibility, NOT chain-integrity"), and audit-procedures P-25 echoes it ("Schema gaps are not chain-integrity findings... they are MRM-program findings"). The SOC engagement can dispose of a step 12a anomaly without modifying the PI1.1 / PI1.2 opinion language; the finding flows to the institution's MRM-program control narrative instead. That is the right shape.

One refinement worth surfacing: the anomaly-documentation-template's severity table does not list `gen_ai_model_identifier_missing` as a row. The template's enumerated anomaly types in the "Anomaly type:" line should add it explicitly so an institution's first-time control owner does not have to invent a label. Suggested addition under the severity table:

> | **`gen_ai_model_identifier_missing`** (spec §7 step 12a) | **Medium** when isolated to a small subset of model-call entries; **High** when the gap is systemic across an `audit.*` event-class — control-completeness finding for SR 11-7 reproducibility, NOT chain-integrity |

The Medium/High split mirrors how a SOC engagement would approach the disposition: a handful of entries missing `gen_ai.response.model` is a low-friction housekeeping issue; a systemic gap across (say) all routing-decision events is an MRM-program-coverage finding the institution's MRM committee owns. The label tracking through the template, the severity table, and the verifier output keeps the engagement working paper consistent.

### 2.2 The `regulator_fingerprint.*` lifecycle events are testable but the SOC procedure is missing

The three new events (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`) are well-shaped — `notice_artifact_sha256` on the receipt event, `validation_paths_attempted` + `validation_results` on the validation event, `change_management_record_id` on the install event. A SOC team can sample-test the trio: pull all three events for a rotation, confirm the trio matches a single `correlation_id`, confirm `overall_validation: PASS` on the validation event, confirm the change-management record exists in the institution's CM system. That is a clean test design.

The gap: there is no audit-procedure entry in `audit-procedures.md` covering it. P-22 through P-27 cover anomaly evidence completeness and customer-dispute reproduction, but none of them tests the regulator-fingerprint lifecycle. The CUECs list does not name it either — there is no CUEC-CRY entry for "the institution operates a documented procedure for receiving and installing rotated regulator-held public-key fingerprints." Without those, the SOC engagement has the evidence shape but no audit step that pulls it.

Recommended additions:

- **CUEC-CRY-06.** "When the institution rotates its tenant signing key (per HSM-custody guidance) or when the regulator's fingerprint storage rotates, the institution operates a documented reception-and-validation procedure that produces the `regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, and `regulator_fingerprint.installed` operational events. The procedure validates the rotation notice through at least two independent paths (e.g., GPG signature verification AND cross-channel cross-check) before installation. A rotation notice that fails validation triggers IR Scenario 11."
- **P-28 (audit-procedures).** "Regulator-fingerprint rotation lifecycle. For each regulator-fingerprint rotation during the period: pull the matched `regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed` triple by `correlation_id`. Confirm the trio is complete, `overall_validation: PASS`, and the `change_management_record_id` exists in the institution's CM system. For any rotation where `overall_validation: FAIL` was recorded: confirm IR Scenario 11 was activated and disposition is documented. For periods with no rotation: confirm by absence and document the absence in the working paper."
- **TSC-mapping addition.** Under CC8.1 (change management): "regulator-fingerprint rotation lifecycle (`regulator_fingerprint.rotation_received` + `.rotation_validated` + `.installed`)". Under CC6.7 (cryptographic controls): the same triple, framed as "external-trust-anchor reception integrity."
- **Sample-report SOC team appendix addition.** A new row under the "Anomalies" section of the verifier output:

> | `regulator_fingerprint.rotation_*` triple complete + valid | CC6.7 + CC8.1 (external-trust-anchor reception); P-28 evidence |

These four additions close the gap between "the events exist and have a schema" and "the SOC engagement has a written test for the events." Without P-28 and CUEC-CRY-06, an engagement team has to invent the test on the fly, which is exactly the friction the rest of this stack avoids.

### 2.3 The `cold_dr.dryrun_attestation_consumed` event is the right shape; same gap as above

The annual cold-DR-key dry-run consumption event is a clever evidence pattern. The institution consumes a project-published attestation, validates the signing identity (`signing_identity_validated`), records the consumption outcome (`PASS` / `FAIL`), and archives the attestation. This is a repeatable annual control with a deterministic evidence shape — the SOC engagement can test it in a single pull: confirm the event exists for the prior year, confirm `consumption_outcome: PASS`, confirm `attestation_artifact_sha256` matches the project's published attestation hash.

Same gap as 2.2: no CUEC, no audit-procedure entry, no TSC-mapping line. Recommended additions:

- **CUEC-IR-05.** "The institution operates an annual procedure to consume the project-published cold-DR-key dry-run attestation (`KEY-DR-DRYRUN-{year}.asc`) within 30 days of the project's annual dry-run window. The procedure validates the signing identity, records the consumption outcome, and archives the attestation as evidence of the cold-DR fallback path. A `consumption_outcome: FAIL` triggers IR Scenario 11."
- **P-29 (audit-procedures).** "Cold-DR-key dry-run attestation consumption. For each reporting period that includes the project's annual dry-run window: pull the `cold_dr.dryrun_attestation_consumed` event. Confirm `consumption_outcome: PASS`, the `attestation_artifact_sha256` matches the project's published artifact hash, and the attestation is archived in the institution's evidence repository. For periods that do not include the annual window: confirm by absence and document the absence."
- **TSC-mapping addition.** Under CC7.5 (recover from security incidents): "cold-DR-key dry-run attestation (`cold_dr.dryrun_attestation_consumed`)" — the institution exercises the cold fallback path annually as evidence of recovery readiness.

The cold-DR pattern is one of the more defensible items in the package for an engagement covering Availability or Processing Integrity criteria, but only if the audit-program lists it. Right now it does not.

### 2.4 P-26 model-inventory composition is the right cross-check at the right layer

P-26 testing the institution's MRM model inventory against the per-model decision-count distribution from chain entries is a clean composition: the chain provides the ground-truth list of models that actually produced decisions, the institution's MRM program provides the list of models the MRM program is aware of, and the cross-check surfaces gaps either direction. The two failure modes are correctly labelled:

- "Models in chain but NOT in inventory" (High; model-governance gap — production model the MRM program does not know about)
- "Models in inventory but NOT in chain" with the two-root-cause split (decommissioned-but-still-listed = housekeeping; chain-coverage-gap = High)

The procedure is institution-side reconciliation, not chain-integrity verification, which is the right framing — the chain's integrity property is unaffected by inventory composition. The SOC opinion on PI1.1 / PI1.2 is independent of whether the MRM inventory is current.

One refinement: the stratification language ("decision-class (routing, scoring, advisory, denial)" and "customer-impact tier") relies on the institution operating those taxonomies. Most institutions in the regional-bank market I work with do not yet have a formal customer-impact tier on their AI inventory. The procedure should add a fallback: "If the institution does not operate a customer-impact tier taxonomy, the SOC team confirms the absence is documented in the institution's MRM-program scope statement and stratifies the cross-check across the institution's available risk taxonomy (typically materiality-of-decision or regulatory-exposure tier)." Without that fallback, the procedure reads like it requires a taxonomy the institution may not have, which puts the engagement in the awkward position of finding the institution non-conformant on a procedure rather than on a control.

### 2.5 P-27 customer-dispute reproduction is the most ambitious procedure in the stack and the one most likely to drift in practice

P-27 names a specific schema for reproduction-result evidence under `audit.reproduction.*`: `original_run_id`, `original_seq`, `iterations`, `decision_equivalent_count`, `decision_equivalent_threshold`, `outcome` ∈ {`decision_consistent`, `variation_within_threshold`, `variation_exceeds_threshold`}. The schema is well-shaped. The procedure correctly notes that the reproduction's input set must be reconstructed from chain-captured `gen_ai_parameters` + `audit.*` payload (not from outside-chain sources) and that the reproduction's output is itself a chain entry with the original decision's `(run_id, seq)` as parent.

The drift risk: the threshold is "institution's MRM-defined tier-specific threshold" — an institution-defined value with no reference floor. In practice, an institution will either (a) set the threshold so loose that `variation_within_threshold` always fires, defeating the cross-check, or (b) set it so tight that legitimate model nondeterminism produces `variation_exceeds_threshold` constantly, drowning the MRM committee in non-actionable escalations. The SOC engagement testing this procedure will see the threshold and have to evaluate reasonableness, which is a judgement call the engagement team is not always equipped to make on a model-by-model basis.

Recommended addition to P-27:

> **Threshold reasonableness review.** For each `decision_equivalent_threshold` value the institution operates: confirm the threshold is documented in the institution's MRM-program tier definitions, confirm the threshold was approved by the MRM committee within the period (or a documented period of the institution's choosing), and confirm the threshold's last-review date is within the institution's MRM-program review cadence. Threshold values without committee approval are documented as MRM-program control gaps. The SOC team does NOT independently judge the numeric value; the SOC opinion attests that the institution's documented governance over the threshold is operating, not that the threshold itself is correct.

This frames the SOC team's role correctly — process-of-governance, not substantive-judgement-of-numeric-threshold. Without that framing, an engagement team can find itself drifting into model-validation territory that belongs to the institution's MRM program.

### 2.6 The CUEC-CRY-04 vocabulary update reads cleanly

The reframing of CUEC-CRY-04 from session-key reconciliation to key-fingerprint reconciliation, with explicit reference to the `(tenant_id, key_version, key_fingerprint)` triple, the failure modes (botched rotation / restored backup / cross-tenant key swap), and the pairing with `master.reconciliation_completed` — reads as an institution-control description that a SOC engagement can test directly. The cadence floor ("no more than weekly per spec §10.1") gives the SOC team a baseline to test against. The audit-procedures P-6 cross-reference closes the loop.

The user-entity-summary's surfacing of `master_key.retired` and `audit_file.truncation_detected` for downstream user-entity audiences is a thoughtful touch — the user entity reading a SOC report on this asset class typically does not have visibility into IKM-retention semantics, and surfacing the events at the user-entity layer with the spec-rule rationale ("the IKM MUST be retained as long as any chain entry referencing it is retained; premature retirement causes affected events to fail verification at spec §7 step 7") gives the user entity the right level of detail without dumping the full normative spec on them. That is the right composition for a downstream-user audience.

### 2.7 Section 4 template covers the spec-defined surface but lags the new events

`section-4-template.md` Section C "Procedures" lists daily seal job, HSM PIN rotation, IKM rotation, IKM retirement, key-fingerprint reconciliation, mid-write-truncation recovery, and verifier validation. It does NOT list:

- The regulator-fingerprint rotation reception procedure (§2.2 above)
- The annual cold-DR-key dry-run attestation consumption (§2.3 above)

Both should be added to Section C as bullet items so the institution's Section 4 description names the procedure as in-scope and the SOC opinion attaches to it. Without the description naming the procedure, the SOC team cannot test the procedure as part of the engagement scope. Suggested additions:

- "Regulator-held fingerprint rotation reception per `[institution-defined procedure document]`; recorded as the `regulator_fingerprint.rotation_received` / `.rotation_validated` / `.installed` operational event triple with `change_management_record_id` per spec §10.2 and threat-model §2.9"
- "Annual cold-DR-key dry-run attestation consumption per `[institution-defined procedure document]`; recorded as the `cold_dr.dryrun_attestation_consumed` operational event with `signing_identity_validated` and `consumption_outcome` per spec §10.2 and supply-chain.md guidance"

### 2.8 The sample-report SOC team appendix is the strongest evidence cross-walk in the package

Lines like:

> `KMS handle URI` (non-`plaintext-` prefix) → CC6.7 (HSM custody); verification that production seals are HSM-backed
> `FAILURE RECORD step: 8 key_fingerprint mismatch` → CC6.1 + CC6.7; identity-mismatch finding at the IKM-roster layer
> `FAILURE RECORD step: 9 payload_hash MAC mismatch` → PI1.1 + CC6.8; content-tampering finding at the chain layer

This is the cross-walk a SOC engagement working paper carries verbatim. I have not seen another asset class where the verifier output line maps one-to-one to TSC criteria without engagement-team interpretation. It is the difference between a SOC engagement that can defend its evidence in a peer review and one that has to reconstruct the mapping from primary sources.

The two rows missing from the appendix that should be added (per §2.2 and §2.3):

> | `regulator_fingerprint.rotation_*` triple complete + valid | CC6.7 + CC8.1 (external-trust-anchor reception); P-28 evidence |
> | `cold_dr.dryrun_attestation_consumed` annual event PASS | CC7.5 (recovery-readiness); P-29 evidence |

And one row that warrants tightening — the `gen_ai_model_identifier_missing` (step 12a) anomaly is not in the appendix at all. Suggested:

> | `Anomalies: gen_ai_model_identifier_missing at seq N` | (Not a TSC chain-integrity finding); SR 11-7 reproducibility-completeness — flows to institution's MRM-program control narrative |

The "Not a TSC chain-integrity finding" framing is important because the appendix's other rows all map to TSC criteria, and a reader scanning the table for step 12a would otherwise infer it should map to PI1.1 or PI1.2 the same way step 9 does. Calling out the non-TSC framing explicitly avoids the inference error.

---

## 3. Disposition and engagement-readiness

For a Type II SOC 2 engagement covering Common Criteria + Processing Integrity + Availability, with the institution operating Tier 1 + Tier 2 CUECs:

- **Section 4 description.** Buildable from the template with the two procedure additions in §2.7. Estimated effort: institution control owner + SOC engagement Section 4 lead, two working sessions.
- **Test of operating effectiveness over the period.** Buildable from audit-procedures P-1 through P-25 + the recommended P-28 + P-29. Estimated effort: standard SOC 2 Type II engagement scope; no asset-class-specific custom build.
- **Anomaly disposition during the period.** Buildable from the anomaly-documentation-template + the higher-bar evidence floors. With the §2.1 addition (severity row for `gen_ai_model_identifier_missing`), the template covers the verifier's full anomaly surface.
- **CUEC list in Section 4F.** Buildable directly from `CUECs.md` with the two recommended additions in §2.2 (CUEC-CRY-06) and §2.3 (CUEC-IR-05).

The engagement is run-rate-ready. The four small additions I have flagged are the difference between an engagement team that has to invent test design on the fly and one that pulls from a written audit program. That is a meaningful difference in repeatability across institutions and across engagement years; the same engagement team running this stack twelve months later should not have to reinvent the test design.

---

## 4. Items I would surface to the engagement partner

If I were briefing the engagement partner on what to expect from this stack going into a first-year Type II engagement, three points:

1. **The chain's integrity property is independently verifiable.** That removes the SOC engagement's typical exposure on processing-integrity assertions for AI-driven decisions. The verifier produces the evidence; the engagement team confirms the evidence shape and the institution's CUEC operation. The opinion's PI1.1 / PI1.2 language is supportable on substantive testing rather than design-and-implementation testing alone.
2. **The `key_fingerprint mismatch` higher-bar evidence floor is engagement-defining.** Most SOC engagements I have run on AI-adjacent service organizations dispose of cryptographic-control anomalies on management's narrative because the documentation does not specify what the four-piece evidence package looks like. Here the four pieces are written down: IKM-roster row identification, change-management approval, re-verification on corrected roster, reconciliation cross-check. An engagement team accepting management's narrative without those four pieces is leaving the engagement exposed; with the pieces in writing, the engagement working paper is straightforward.
3. **The MRM-program adjacency is intentional and bounded.** P-26 and P-27 reach into MRM-program territory (model inventory, customer-dispute reproduction, decision-equivalent threshold), but the procedures correctly frame the SOC engagement's role as institution-side reconciliation rather than substantive model-validation. With the §2.4 fallback language and the §2.5 threshold-reasonableness framing, the engagement team can run those procedures without drifting into model-validation work that is not the SOC engagement's scope. Worth coordinating with the institution's MRM committee chair and the engagement's model-risk specialist before the engagement period opens.

---

## 5. Items I would NOT raise as findings

Three items I would NOT raise in a peer review of this package — flagged so the engagement team does not waste cycles on them:

1. **The single-step `gen_ai_model_identifier_missing` anomaly is correctly scoped to control-completeness, not chain-integrity.** A peer reviewer unfamiliar with the spec might read step 12a as a chain-integrity check and try to map it to PI1.1. The spec text and the procedure text both name the distinction; the engagement team holds the line.
2. **The `master_version` → `key_version` / `key_versions` vocabulary migration footnote in `control-evidence-events.md`.** Some peer reviewers will flag the migration as a gap if their reading of the prior period used the older vocabulary. The footnote is the correct disposition; comparing across periods uses the institution's documented migration date as the cutover.
3. **The `--master-key` provisioning shape for examiner-time key-bound verification.** The sample-report's invocation correctly provides the IKM at examination time per the institution's IKM-disclosure shape (court order, customer-dispute, examiner request). A peer reviewer might flag the IKM-disclosure as a confidentiality concern; the cross-references to `customer-dispute-procedures.md` and `legal-disclosure.md` are the correct disposition. The engagement team confirms the institution's IKM-disclosure shape is documented and operates per institution policy; the engagement does not re-evaluate the policy itself.

---

## 6. Closing

The package is ready for a first-year Type II engagement with the four small additions I have flagged in §2.2 (CUEC-CRY-06 + P-28 + TSC-mapping line + sample-report appendix row), §2.3 (CUEC-IR-05 + P-29 + TSC-mapping line + sample-report appendix row), §2.1 (anomaly severity row for step 12a), and §2.7 (Section 4 procedure list additions). None of these are blockers; all of them are paper-only additions to documents that already exist in the stack. With them, the stack covers the engagement's testing program end-to-end: every control claim ties to an event, every event has a schema, every schema has an audit step, every audit step has an evidence floor, and every evidence shape ties back to a TSC criterion. That is what a SOC engagement actually needs.

— Marisol Hernández
  Senior Manager, SOC 1 / SOC 2 practice
  Mexico City office
