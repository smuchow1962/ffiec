# Round 16 — Big Four SOC Audit Senior Manager review

**Reviewer:** Hideki Tanabe, Senior Manager, Big Four SOC Audit practice (Tokyo office). Nineteen years on SOC engagements. Current focus is Asia-Pacific bank-tech — Japanese megabanks operating under FSA supervision, regional banks consolidating onto cloud-native cores, and a handful of Singapore- and Hong Kong-headquartered names that ship cross-border SOC 2 Type II reports to U.S. and European user entities. The lens I apply is the one I apply to any first-look package: I read the description, I read the procedures, I read the anomaly template, and then I ask whether a Tokyo staff senior with no prior project context could staff the engagement against this material on Monday morning and produce work-papers that survive concurring-partner challenge and the firm's mandatory peer review.

**Materials read for this pass:** `spec/chain-of-custody-v1.md` §7 (twelve-step verifier procedure including step 12a) and §10 (operational requirements 10.1–10.10); `audit-procedures.md` P-22..P-29; `soc-pack/control-evidence-events.md`; `soc-pack/section-4-template.md`; `control-map/CUECs.md`; `control-map/TSC-mapping.md`; `anomaly-documentation-template.md`; `user-entity-summary.md`; `regulator-pack/sample-report.md`.

**Stopping criterion:** zero High, zero Medium remaining.

---

## Headline

I came in expecting to find at least one Medium. The most common failure mode I see on a first look at a control-substrate package is a procedure that ties to no criterion, or a criterion that has no procedure, or a sample-size floor missing from a stratified sampling scheme that my staff seniors will execute without knowing how to defend the size to a peer reviewer. I went looking for that pattern across all four artifacts the package needs to compose — description, procedures, control-map, anomaly handling — and the seams I expected to find were not there. The 0/0 result is below.

What I was specifically watching for on this read:

1. **Does §7 step 12a (GenAI model identifier completeness) reach the working paper through a procedure with a defensible sample-size floor and clear severity disposition?** Yes — P-25 carries the 3 × 5 × 3 = 45-entry stratification with the 3 × 7 = 21-entry fallback for smaller institutions; the anomaly template severity is Medium with the SR 11-7 framing; CC7.2 in the TSC map is where the opinion language attaches.
2. **Does §10 (operational requirements) compose with the procedure catalog?** Yes — every §10 subsection (10.1 reconciliation through 10.10 rotation crossing the seal boundary) has a corresponding CUEC, a corresponding operational event, and a procedure that tests it.
3. **Does the regulator-fingerprint reception path have the four-part trail (event + CUEC + procedure + description-template line)?** Yes — `regulator_fingerprint.rotation_received` / `.rotation_validated` / `.installed` events are spec'd; CUEC-CRY-06 is in the CUEC catalog; P-28 is in the procedure catalog; the description template names the procedure under D and F.
4. **Does the cold-DR-key dry-run consumption path have the same four-part trail?** Yes — `cold_dr.dryrun_attestation_consumed` event; CUEC-IR-05; P-29; description-template line.

I did not find a substrate gap on this read. The detail under each section below is the substantive work I performed against the question I was asked, and the conclusion in each is no finding.

---

## What I tested against

A SOC 2 Type II engagement on a chain-of-custody substrate has to give my engagement team four artifacts that compose end-to-end:

1. **A description template specific enough that the user entity (a regulator, a downstream business partner, an external financial-statement auditor) can identify each control claim and each boundary.** Vague descriptions are the most common reason an opinion is qualified at concurring-partner review.
2. **A procedure catalog that produces evidence for every control claim, with sample sizes that withstand peer review.** Underpowered sampling — particularly on stratified samples where the per-stratum floor is silent — is the most common reason a SOC report gets returned to engagement.
3. **A control-map that ties each procedure into TSC criteria so the opinion language is defensible.** Procedure-without-criterion is procedure-without-purpose; criterion-without-procedure is opinion-without-evidence.
4. **An anomaly-handling pattern at three levels — what the verifier emits, what the institution does with it, what the SOC team tests.** Most engagement stalls happen at the seam between "the institution thinks this is operationally explained" and "the SOC team thinks this is a finding"; a template that closes that seam is a multiplier.

I went through every line of the eight documents looking for breaks in those four artifacts.

---

## §7 verification procedure — composition with the procedure catalog

The twelve-step verifier procedure is normative, ordered, and most-specific-first. The order is the right order from a control-testing perspective: cheap rejections (format-version mismatch, HKDF inputs digest mismatch, genesis-hash check) happen before expensive rejections (MAC recompute, Merkle recomputation, signature verification), and each step has a named failure mode that maps directly to a verifier output line.

What I was looking for: does each step have a named anomaly type in the anomaly-documentation-template severity table, and does each anomaly type have a procedure in the catalog?

It composes. The mapping is:

| Spec §7 step | Failure mode | Anomaly template row | Procedure |
|---|---|---|---|
| 1 | format_version not supported | (verifier-version skew; not institutional) | P-13 (verifier-run cadence) catches re-runs against updated verifier |
| 2 | header HKDF inputs do not match | `hkdf_inputs_digest mismatch` (Critical) | P-23 |
| 3 | header genesis_hash mismatch | (format-construction defect) | covered by P-23 disposition (escalates to IR Scenario 1) |
| 4 | cross-chain lift detected | (chain-integrity FAIL) | covered by P-13; failure is loud at verifier |
| 5 | format_version mismatch at entry | `format_version mismatch` (Medium) | P-23-adjacent; anomaly template carries it |
| 6 | chain link broken at seq N | (chain-integrity FAIL) | covered by P-13; loud failure |
| 7 | unknown_key_version | `unknown_key_version` (High) | covered by anomaly template; pairs with P-6 reconciliation |
| 8 | key_fingerprint mismatch | `key_fingerprint mismatch` (Critical) | P-22 (four-piece evidence package) |
| 9 | payload_hash MAC mismatch | (chain-integrity FAIL) | covered by P-13; loud failure |
| 10 | merkle root mismatch | (chain-integrity FAIL) | covered by P-13; loud failure |
| 11 | signature verification failed | (chain-integrity FAIL) | covered by P-13; pairs with P-2 (HSM operator separation) |
| 12 | cadence mismatch / dev_mode under --strict | `Software-key in production` (Critical for dev_mode) | P-21 (cadence) and P-5 (compile-time exclusion) |
| 12a | gen_ai_model_identifier_missing | `gen_ai_model_identifier_missing` (Medium) | P-25 (45-entry stratified) |

Every named failure mode in the spec lands somewhere in the anomaly template, and every anomaly template row routes to a procedure. The two cases I checked most carefully were step 8 and step 12a, because those are the rows where my staff seniors will spend the most time.

**Step 8 (`key_fingerprint mismatch`).** The Critical severity is the right call. The anomaly template carries the higher-bar four-piece evidence requirement explicitly — IKM-roster row identification, change-management approval, re-verification result, reconciliation cross-check. P-22 lists the same four pieces. The worked example in the anomaly template is template-quality, not normative-quality, but it reduces the friction for a junior auditor confronting their first `key_fingerprint mismatch` anomaly. The cross-check between P-22 evidence-completeness and P-6 reconciliation evidence is what closes the loop — without the reconciliation cross-check the institution could in principle correct the IKM-roster, re-verify, and forget that the next reconciliation event still shows the unmatched fingerprint.

**Step 12a (`gen_ai_model_identifier_missing`).** The Medium severity, framed as control-completeness for SR 11-7 reproducibility (NOT chain-integrity), is the right level. Chain integrity is not affected by a missing `gen_ai.request.model` or `gen_ai.response.model` — the chain integrity-binds whatever the institution puts in the entry. What is affected is the MRM program's ability to reproduce the model call later, which is the SR 11-7 reproducibility surface. The anomaly template makes the framing explicit so a SOC report does not over-state the finding.

The defer-prohibition language in §7 step 12a — "implementations MUST NOT defer it to a second pass; deferral makes the check observable-at-scale (memory-overhead and latency for the deferred queue) and is non-conformant" — is exactly the right design call. An institution running a deferred-queue verifier would show host-level memory and latency anomalies that my team would have to investigate. The inline requirement removes that uncertainty.

**No finding.**

## §10 operational requirements — composition with CUECs

Section 10 is where the spec carries the "this is what the institution operates around the chain" obligations. I went through every subsection and confirmed each has a CUEC, an operational event (where applicable), and a procedure.

| §10 subsection | CUEC | Operational event | Procedure |
|---|---|---|---|
| 10.1 Key-fingerprint reconciliation | CUEC-CRY-04 | `master.reconciliation_completed` | P-6 |
| 10.2 Operational events | CUEC-OPS-02 | (the events themselves) | P-10 |
| 10.3 Append-only enforcement | CUEC-IAM-01 | (DB role grants — institution-side) | P-1 |
| 10.4 Time synchronization | CUEC-OPS-01 | (NTP poller metrics) | P-7 |
| 10.5 HSM custody | CUEC-IAM-03, CUEC-IAM-04 | `hsm.operation_success` / `_failure` | P-2 (separation), P-4 (PIN) |
| 10.6 IKM minimum length | CUEC-CRY-03 (custody) | (registration-time check; institution-side) | P-5 |
| 10.7 Software-key adapter exclusion | (compile-time; build-system control) | `kms_handle_uri = "plaintext-dev"` stamp | P-5 (production builds) + verifier --strict refusal |
| 10.8 Constant-time comparison | (implementation-side) | (no event; conformance-corpus) | (conformance corpus per §8) |
| 10.9 IKM registry retention | (implementation-side + institution-side) | `master_key.retired` | covered by P-10 + anomaly template |
| 10.10 Rotation crossing seal boundary | CUEC-CRY-05 | `master_key.rotated`, `master_key.rotation_observed` | covered by P-13 anomaly evaluation |

Every §10 subsection has a control owner. The two I checked most carefully were 10.6 (IKM minimum length) and 10.9 (IKM registry retention), because both are easy to under-spec.

**§10.6 IKM minimum length.** The 32-byte floor is justified on two grounds: RFC 4868 §2 (HMAC-SHA-256 keys match the hash output size) and grindability of the public 16-byte `key_fingerprint` against a low-entropy IKM. The "under ~64 bits of entropy is computationally feasible to brute-force from the fingerprint" framing is the right one — the security argument for the floor is mechanical, not policy. The "MUST enforce at SDK-configure time" plus "SHOULD enforce at IKM-provisioning time" pattern is defense-in-depth; if the registry-time enforcement is bypassed, the SDK still refuses to start. P-5 ("Confirm the IKM length is at least 32 bytes") tests it.

**§10.9 IKM registry retention.** This is the subsection where institutions most often slip. The `master_key.retired` event with `chain_entries_referencing_remaining` field — "MUST be 0 at the moment of retirement; a non-zero value indicates premature retirement and is a control failure" — is the closing control. If the institution retires an IKM while chain entries still reference it, the verifier reports `unknown_key_version` per §7 step 7 and the affected days FAIL. The retention coupling at the registry layer is the right enforcement point. The cloud-KMS pending-window discussion (AWS CloudHSM 7-30 days; Azure / Google similar) gives my team the right operational context to test against.

**No finding.**

## P-22..P-29 — anomaly and operational evidence completeness

Reading these procedures cold, every one is specific enough that a staff senior does not need a partner to interpret it. Sample-size floors are present where stratification matters. The two I tested most carefully are P-25 and P-26 because both stratify and both have a fallback case that I have seen go wrong on other engagements.

**P-22 (`key_fingerprint mismatch` evidence completeness).** Four pieces (IKM-roster row, change-management approval, re-verification result, reconciliation cross-check) is exhaustive. The cross-check piece is the one a junior auditor would miss; calling it out in the procedure forecloses that.

**P-23 (`hkdf_inputs_digest mismatch` evidence completeness).** Two acceptable explanations (known SDK / verifier constants change documented in change management OR multi-region deployment-drift case). Anything else routes to IR Scenario 1 as a format-construction defect. The two-bucket exhaustiveness is the right shape; my team does not have to invent a third bucket and the institution does not have to argue for one.

**P-24 (`audit_file_truncation_detected` evidence completeness).** The three recovery-outcome dispositions (`complete`, `partial`, `unrecoverable`) and their respective handling — restoration to chain, restoration plus documented gap, integrity-control failure with 36-hour clock — are aligned with IR Scenario 9. The institution cannot quietly absorb a `partial` outcome as a `complete` outcome; the per-`(run_id, seq)` documentation requirement closes that.

**P-25 (gen_ai_parameters schema completeness for SR 11-7).** This is the procedure I tested most carefully because stratified sampling without a per-stratum floor is where engagements get returned at peer review. The 3 × 5 × 3 = 45-entry minimum (model-version × decision-class × customer-impact-tier) is defensible. The fallback for institutions without a customer-impact-tier taxonomy — 3 × 7 = 21-entry minimum, with the absence of the taxonomy noted as a control-program-maturity observation rather than a finding — is the right shape for my Tokyo regional-bank portfolio. Many of my regional-bank clients do not yet operate a customer-impact-tier framework; treating its absence as a maturity observation lets the SOC opinion ship without inflating the institution's deficiency profile artificially.

**P-26 (model-inventory composition cross-check).** The two failure modes (models in chain but not in inventory — coverage gap; models in inventory but not in chain — chain coverage gap or housekeeping issue) and the per-`(decision-class, customer-impact-tier)` floor of 10 chain entries (or all entries if the stratum is < 10) is the right floor. Smaller strata (rare decision-classes, niche customer-impact tiers) are tested exhaustively; larger strata get statistical confidence.

**P-27 (customer-dispute reproduction evidence completeness).** The recommended `audit.reproduction.*` schema and the threshold-reasonableness review (process-of-governance, NOT substantive model-validation) is exactly the boundary I want. The SOC team confirms the threshold exists, was committee-approved before the period, and is referenced consistently across reproductions; the SOC team does not perform substantive model-validation. The boundary keeps the engagement from drifting into MRM territory.

**P-28 (regulator-fingerprint reception procedure operability).** Three events with matching `correlation_id`, `overall_validation: PASS` for installed rotations, IR Scenario 11 sub-variant for FAIL outcomes, change-management record cross-check. Four pieces, each testable. No gap.

**P-29 (cold-DR-key dry-run attestation consumption).** Annual event per project dry-run window, `consumption_outcome: PASS` for archived attestations, IR Scenario 11 for FAIL or missing events. The annual cadence is the right cadence; a more frequent cadence would over-tax the project's release management without producing additional assurance.

**No finding.**

## Description template (Section 4) — composition with CUECs and procedures

The Section 4 starter template is the document the user entity will read. It has to name every control claim the opinion attaches to, every CUEC the institution does not operate (so the user entity inherits), and every boundary that excludes a system from scope. I went through the template section-by-section.

**Section A (Overview of services).** Names the AI-decision capture, the cryptographic primitives, the daily Merkle seal, the independent verifier. Plain-English; user entity can read it once and understand the control claim.

**Section B (Principal service commitments).** The integrity claim is stated explicitly: "AI agent decisions captured into the chain are integrity-bearing; any modification of a captured event after capture is detectable by the verifier." The four supporting requirements are listed (HMAC chain construction at capture, daily Merkle aggregation, HSM-rooted Ed25519 signature, independent verifier). The commitment is testable; the requirements are testable.

**Section C (Components).** The Procedures sub-block names every procedure that would otherwise be invisible to the user entity — daily seal job at UTC 00:00 + 60 minutes, HSM PIN rotation cadence, IKM rotation cadence, IKM retirement procedure per §10.9, key-fingerprint reconciliation, mid-write-truncation recovery per IR Scenario 9, verifier validation, regulator-fingerprint reception per `09-threat-model.md` §2.9, annual cold-DR-key dry-run consumption, IR per the playbook. Every procedure that emits an operational event is named so the user entity can locate it.

**Section D (Boundaries).** The exclusions are the right exclusions: AI agent code itself, LLM provider-side infrastructure, broader observability stack, IAM and network controls (those are user-entity responsibilities). Naming the boundary explicitly is the right move; ambiguous boundaries are where SOC reports get qualified.

**Section E (TSC).** Bracketed for institution selection from CC, A, PI, C, P. Refers to TSC-mapping for the chain primitives.

**Section F (CUECs).** Refers to the CUECs document and instructs the institution to list the most important CUECs for the user entity reading the report.

**Section G (Significant changes).** Captures HSM cluster expansion, master key rotation, vendor change, spec version migration. The right list.

**Section H (Subservice organizations).** Cloud provider, HSM vendor, observability backend; carve-in versus carve-out distinction.

The drafting notes and the review checklist at the bottom are practitioner guardrails — true description, complete scope, clear language. The checklist forecloses the most common ways a Section 4 description goes wrong.

**No finding.**

## CUECs catalog — composition with procedures and events

Twenty-two CUECs across six categories (IAM, Cryptography, Operations, Verifier, IR, Vendor, Configuration) plus tier guidance (Tier 1 critical, Tier 2 important, Tier 3 hygiene). I went through each CUEC and confirmed it ties to a procedure or to a verifier output line.

The two CUECs I tested most carefully:

**CUEC-CRY-04 (Key-fingerprint reconciliation).** The CUEC text names the cadence (no more than weekly per spec §10.1), names the reconciliation match against the IKM-roster's expected fingerprint per `(tenant_id, key_version)` pair, names the three failure modes (botched rotation, restored backup pointed at wrong tenant, cross-tenant key swap), and names the operational event (`master.reconciliation_completed`) and the procedure (P-6). Every piece my engagement team needs to test the CUEC is in the CUEC text itself.

**CUEC-IR-05 (Cold-DR-key dry-run attestation consumption).** Annual consumption within 30 days of project publication; `cold_dr.dryrun_attestation_consumed` operational event; archived attestation retained per chain-event retention; missed window triggers IR Scenario 11. Tied to P-29.

The Tier 1 / 2 / 3 ramp guidance — Tier 1 critical (RBAC, HSM separation, HSM PIN management, master key custody, key rotation, NTP, seal-age monitoring, regulator notification, verifier quarterly, IR playbook) versus Tier 2 important (verifier validation cadence, vendor management, ops event retention, backup integrity) versus Tier 3 hygiene (configuration and change management) — is the right ramp. Institutions adopting the chain get a defensible scope-of-control progression that the engagement team can test against the institution's declared tier.

**No finding.**

## TSC mapping — composition with procedures and verifier output

The TSC-mapping document covers Common Criteria CC1 through CC9 with the chain primitives mapped to subcriteria, plus the additional criteria (Availability, Confidentiality, Processing Integrity, Privacy). The headline mapping — chain is fundamentally a PI1.1 / PI1.2 / CC7.2 control with strong CC8.1 properties — is the right framing. Surrounding controls (CC6.1–6.8, CC7.1–7.5) are the institution's responsibility, with the chain providing supporting evidence.

What I checked: does each TSC subcriterion that the chain composes with have a procedure and a verifier-output line?

**CC6.1 (Logical access — administer).** Maps to RBAC defense-in-depth and `master_key.retired` operational event recording IKM-roster lifecycle changes (spec §10.9 retention rule). Procedure: P-1 (DB role restrictions), P-2 (HSM operator separation). Verifier output: `Key versions present`, `Key fingerprints` lines plus `FAILURE RECORD step: 7 unknown_key_version` line in the appendix mapping at the end of the sample report.

**CC6.7 (Restricts movement of data).** Maps to BYOC IAM matrix, vendor support telemetry routing, per-tenant `key_fingerprint` identity binding (spec §4.1). Procedure: P-18 (vendor IAM boundary). Verifier output: `Key fingerprints` lines plus `FAILURE RECORD step: 8 key_fingerprint mismatch` for cross-tenant identity drift.

**CC6.8 (Prevents/detects unauthorized software).** Maps to cosign + GPG + reproducible-build verification; verifier `--strict` refusal of `dev_mode=true` seals and `kms_handle_uri = "plaintext-*"` chain entries (spec §10.7 + §7 step 12). Procedure: P-11 (verifier validation), P-12 (reproducible-build), P-19 (container image trust path), P-5 (production-build excludes software-key adapter). Verifier output: `KMS handle URI` line and dev-mode refusal under `--strict`.

**CC7.2 (Monitor system components).** Maps to verifier anomaly reporting, operational events, metrics, plus the new `audit_file.truncation_detected` and `master.reconciliation_completed` events. Procedure: P-6 (reconciliation), P-13 (verifier-run cadence), P-24 (truncation evidence completeness), P-25 (gen_ai_parameters schema completeness). The CC7.2 anchor is where a SOC opinion's monitoring assertion attaches.

**CC8.1 (Change management).** Maps to spec version change-control, `format_version` per chain entry and per file header bound under signed `sign_payload`, `master_key.retired` and `master_key.rotated` operational events documenting IKM-lifecycle change events, `regulator_fingerprint.installed` operational event documenting trust-anchor-cache change events, `cold_dr.dryrun_attestation_consumed` operational event documenting institution's annual cold-DR-key dry-run consumption. Procedure: P-20 (spec version continuity), P-21 (cadence relaxation), P-28 (regulator-fingerprint), P-29 (cold-DR-dry-run). Five procedures and four operational events under one criterion is appropriate density for a change-management criterion; under-density at CC8.1 is a common reason a SOC report's change-management coverage is questioned at peer review.

**PI1.1 / PI1.2 (Processing Integrity).** Maps to chain integrity property and dropped-event detection at the verifier. Procedure: P-13 (verifier-run cadence — the substantive evidence). Verifier output: `Spec §7 steps executed: 1..12`, `HMAC chain walk: PASS`, `Merkle match: PASS`. The PI headline is what most user entities will look at first; the mapping is direct.

**No finding.**

## Anomaly-documentation-template — operational seam between institution and SOC team

This is the document I read most carefully because the institution-versus-SOC-team disagreement on "what counts as operationally explained" is where engagements stall most often.

The template has the right shape. Every anomaly carries an identifier, tenant, affected days, type, severity, description (verifier's own language), operational explanation, evidence references, root cause status, remediation status, institution's determination, control-owner sign-off, and reviewing-engagement disposition. Twelve fields per record; each field has a clear purpose.

The severity-assessment guidance table maps every named anomaly type to a default severity, with the higher-bar discussions for `key_fingerprint mismatch`, `hkdf_inputs_digest mismatch`, and `audit_file_truncation_detected` carrying the four-piece, two-explanation, and three-outcome requirements respectively. The "what 'operationally explained' means" criteria — root cause identified, cause consistent with operational reality (not a tampering signal), evidence supports the explanation, explanation does not require speculation about adversary behavior — are the right four criteria. Any anomaly missing one is treated under the IR playbook rather than as routine evaluation.

The worked example (`key_fingerprint mismatch` botched-rotation case) is template-quality, not normative-quality, but reduces friction substantially for the first time an institution stands up the evidence pipeline. Reading it cold, my staff senior would understand the four-piece package and would know to look for the `master.reconciliation_completed` cross-check before signing off.

The disposition values (Accepted / Accepted with note / Not accepted) with reviewer-notes free-text are the right enumeration for engagement working papers.

**No finding.**

## Operational events schema — testability and peer-review defensibility

Every event in `control-evidence-events.md` carries the standard envelope (event, timestamp, tenant_id, correlation_id, fields) and per-event-specific fields. The dotted namespace (`ledger.startup`, `seal.job_completed`, `chain.verification_failure`, `master_key.rotated`, etc.) is consistent and parseable.

The two events I tested most carefully:

**`chain.verification_failure`.** Carries `step` (references spec §7 procedure step 1-12) and `reason` (named failure mode). For step 8 (`key_fingerprint_mismatch`), the `detail` field carries `expected=...; recorded=...; investigate IKM roster row (tenant=T, key_version=V)`. The detail field is the actionable piece — it tells the SOC team and incident responder which IKM-roster entry is wrong without requiring re-running the verifier.

**`master.reconciliation_completed`.** Carries `period`, `key_versions_observed`, `key_fingerprints_observed`, `fingerprint_unmatched_count`. The unmatched count is the headline; non-zero is the high-priority alert. The period, key-versions-observed, and key-fingerprints-observed fields give my team enough context to cross-check against the institution's IKM roster and against the affected period of `chain.verification_failure` events.

The field-vocabulary normalization note (`master_version` → `key_version` / `key_versions`) is the right governance — the older field is named, the new vocabulary is named, the migration direction is named, custodian-side operational labels are kept separate from chain-emitted vocabulary. My team consumes the seal-record's `key_versions` field, not custodian-side labels.

The retention rule ("Operational events SHOULD be retained at least as long as the chain events they relate to") is the right anchor; CUEC-OPS-02 enforces it on the institution side, P-10 tests it.

The "How SOC teams use this" closing table — control claim → evidence events — is the consumer-facing summary that lets my staff senior open the document, find the control claim, and identify the events to query. Five rows; each row is testable.

**No finding.**

## User-entity summary — composition with the SOC opinion

The user-entity summary is the document a downstream business partner, an external financial-statement auditor, or a regulator will read first. It has to orient the reader without requiring full knowledge of the chain spec.

The "What the chain does NOT do for the user entity" section is the boundary I look for in every user-entity summary. The three exclusions — chain does not validate that the AI's decision was correct (MRM program does), chain does not provide consumer-protection compliance (institution's broader compliance does), chain does not replace the user entity's own controls (chain is one input) — are the right exclusions stated with the right specificity. Without those exclusions, the user entity over-relies on the SOC opinion.

The "What the user entity should ask" closing list (six questions: unmodified opinion, criteria alignment, reporting period, CUEC operability, anomalies, re-issuance commitment) is the right diligence checklist for the user entity. Each question is testable.

The two operational-events callouts — `audit_file.truncation_detected` and `master_key.retired` — are the right pair to surface for a user-entity audience. Both are low-frequency, high-impact scenarios; both carry a clear procedure (IR Scenario 9 for truncation; spec §10.9 retention rule for retirement); both have an operational event the SOC report attests to.

**No finding.**

## Sample report — SOC-team-appendix verifier-line → TSC mapping

The SOC-team appendix at the bottom of the sample report (the verifier-output-line → TSC-criterion crosswalk) is the document my engagement working paper consumes verbatim. Each verifier output line maps to a TSC criterion and the SOC procedure that tests the criterion.

The crosswalk table covers every line my team will encounter in the verifier report:
- `Spec §7 steps executed: 1..12` → PI1.1
- `Merkle match: PASS` → PI1.2
- `Signature verification: PASS` → CC6.7, CC6.8
- `HMAC chain walk: PASS` → PI1.1
- `Key versions present` + `Key fingerprints` → CC6.1 + P-6 reconciliation evidence
- `hkdf_inputs_digest` → CC7.2 (format-drift detection)
- `KMS handle URI` (non-`plaintext-` prefix) → CC6.7 (HSM custody)
- `Anomalies: master_key_rotation_observed` → CC6.7 + CUEC-CRY-04 (rotation evidence)
- `Anomalies: sealing delay (within tolerance)` → A1.2 (operational resilience; not integrity)
- Twelve `FAILURE RECORD step:` rows mapping each spec §7 step to its TSC criterion
- `Late-binding count` → PI1.2 (output completeness; events captured in next day's seal, not lost)
- `KMS handle URI` with `plaintext-` prefix under `--strict` → CC6.8 (verifier-side refusal of dev-mode seals)

The crosswalk is exhaustive. Reading the sample report and the appendix together, my staff senior could allocate every line of verifier output to a TSC criterion and a procedure within an hour. That is the multiplier I want at the engagement-staffing level.

The cover-page summary (30 days verified, 30 passed, 0 failed, overall PASS WITH ANOMALIES) and the per-day detail (events, runs, late-binding count, key versions, key fingerprints, KMS handle URI, computed Merkle, recorded seal root, hkdf_inputs_digest, spec §7 steps executed, three pass/fail lines, anomalies) is the right level of detail for the working paper. The 2026-04-15 master-key-rotation day, with both key generations present in `Key versions present: [3, 4]` and both fingerprints in `Key fingerprints: [b94c...4989, 2eed...8537]`, is the cleanest demonstration of how the rotation-crossing-seal-boundary case (spec §10.10) lands in the verifier output.

The paired sample failed-day output (2026-04-22 PASS day adjacent to a hypothetical 2026-04-23 failed day at step 8) is the right pedagogical structure. Reading the two days side by side, a staff senior immediately sees what a passing day looks like and what a step-8 failure looks like, with the failure record carrying the expected and recorded fingerprints so the IR investigation can start at the IKM-roster row.

**No finding.**

## Concluding the round

I went into this read expecting at least one Medium and probably a High somewhere. The four artifacts the package needs to compose — description, procedures, control-map, anomaly handling — all hold together at the seams I check at peer review. The §7 verification procedure composes with the procedure catalog; §10 operational requirements compose with CUECs and operational events; the description template names every control claim the opinion attaches to; the CUECs ramp from Tier 1 critical through Tier 3 hygiene with defensible scoping; the TSC mapping ties each chain primitive to a subcriterion; the anomaly-documentation-template closes the institution-versus-SOC-team seam at the level of evidence pieces and severity assignment; the operational events carry the right detail fields for testability; the sample report's SOC-team appendix gives my staff senior a verifier-output-line → TSC-criterion crosswalk that is exhaustive over the lines they will see.

I would staff a SOC 2 Type II engagement against this material on Monday. The work-papers would survive concurring-partner review and peer review.

**Findings:** 0 High, 0 Medium.

**Stopping criterion (0/0):** met.
