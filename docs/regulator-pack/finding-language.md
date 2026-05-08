# Examination-report finding language

> **What this doc is.** Standard examination-report language for findings the chain verifier produces. Examiners adopt or adapt the language; consistency across the regulator's portfolio matters more than the specific words.

## How to use

The verifier produces specific failure modes. Each failure mode below maps to one or more candidate examination-finding paragraphs. The examiner picks the paragraph that fits, edits for institution-specific facts, and inserts into the examination report.

Severity guidance keyed by spec §7 procedure step. The verifier emits a `step` field on every `chain.verification_failure` operational event and on every per-day failure record in the JSON report; map the step number to this table. The IR scenario column points to the institution-side response in `docs/incident-response-playbook.md`.

| `step` | Verifier reason string | Severity guidance | IR scenario |
|---|---|---|---|
| 1 | `format_version not supported by this verifier (running v1)` | Not an institutional finding — examiner verifier-version skew. Obtain the matching verifier from the project supply chain. | (none — examiner-side) |
| 2 | `header HKDF inputs do not match running v1 inputs` | Severe MRA; the file's HKDF inputs claim v1 but compute differently. Investigate whether the institution's SDK or the verifier's constants drifted. | Scenario 1 |
| 3 | `header genesis_hash does not match v1 constant` | Severe MRA; format-construction defect. Investigate the SDK's chain construction. | Scenario 1 |
| 4 | `cross-chain lift detected at seq N` | MRA; an event in the file claims a different `tenant_id` or `run_id` than the file header. Either the institution mis-bundled snapshots (operational) or an attacker attempted to lift an event from one chain into another (integrity). Require evidence of the event's original chain provenance. | Scenario 1 + evidence-handling investigation |
| 5 | `format_version mismatch at seq N` | MRA; mid-file format drift. Investigate the SDK's version handling. | Scenario 1 |
| 6 | `chain link broken at seq N` (or `seq out of order`) | MRA; possible tampering or insertion/deletion. Immediate response required. | Scenario 1 |
| 7 | `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` | MRA; institution's IKM roster is incomplete relative to its actual operations. Could be premature IKM retirement, missing provisioning, or stale verifier-side registry. | Scenario 8 |
| 8 | `key_fingerprint mismatch at seq N` | **Severe MRA — the load-bearing new failure mode.** The looked-up IKM does not match the entry's recorded fingerprint. Investigated against the IKM roster and §10.1 reconciliation log, NOT against chain content (this is an identity mismatch, not a content tampering). | Scenario 7 |
| 9 | `payload_hash MAC mismatch at seq N` | Severe MRA; content tampering or a defect in the chain construction. Immediate response required. | Scenario 1 |
| 10 | `merkle root mismatch — ledger contents do not produce sealed root` | Severe MRA; the day's events do not produce the sealed root. Possible server-side history rewrite. | Scenario 2 |
| 11 | `signature verification failed` (or `algorithm/key-type mismatch`) | Severe MRA; potential HSM key compromise OR misconfigured public-key registry. Immediate response required. | Scenario 3 |
| 11 (dual-algo) | `partial-coverage seal: single-algorithm signature during institution's declared dual-algorithm posture` | Observation under non-strict; PASS-WITH-ANOMALY (control-completeness, NOT chain-integrity); under `--strict` escalates to MRA. | Scenario 5 variant (control-completeness) |
| 11 (dual-algo) | `algorithm not on institution's declared posture list at seal_date {D}` | Under `--strict` MRA; under non-strict Observation. The institution's declared algorithm posture and the seal's actual algorithm disagree; investigate posture documentation and SDK configuration. | Scenario 5 variant |
| 11 (dual-algo) | `co-signed seal failure: algorithm X validated, algorithm Y did not` | Severe MRA. One algorithm validated; the other did not. Either an algorithm has been broken (un-broken algorithm's signature still provides integrity assurance) or a per-algorithm signing key has been compromised. The institution coordinates with the regulator on the migration timeline OR activates Scenario 4 (master-key compromise) for the failed algorithm. | Scenario 4 + Scenario 3 |
| 12 | `cadence mismatch` | Observation; the seal record's recorded cadence does not match the institution's claimed cadence in its control description. Control-description-accuracy finding; escalate to MRA on repeat. | Scenario 5 variant |
| 12 | `dev-mode seal in production verification — refused` | Severe MRA; production seals signed with non-HSM key. | Scenario 6 |
| (file pre-flight) | `audit file ends mid-line — possible mid-write crash` | Operational (sealing-delay-equivalent severity); writer-side crash, not tampering. Direct the institution to its incident log and crash-recovery procedure. | Scenario 9 |
| (anomaly) | `late-binding event count exceeds threshold` | Observation; operational concern, not integrity. | (anomaly only — no IR) |
| (anomaly) | `sealing delay > 24h` | Observation; investigate operational cause. | Scenario 5 |
| (anomaly) | `sealing delay > 72h without notification` | MRA; institution failed to notify per spec §4.3.1. | Scenario 5 |
| (anomaly) | `clock-skew anomalies` | Observation; operational time-sync issue. | (anomaly only — no IR) |
| (anomaly) | `master_key_rotation_observed` | Normal-operations PASS-with-anomaly when documented in the institution's incident log. Investigate when undocumented. | (PASS-with-anomaly when documented; Scenario 7 if undocumented) |
| 12a | `gen_ai_model_identifier_missing at seq N` | Medium MRA under non-strict; FAIL under `--strict`. Control-completeness finding for SR 11-7 reproducibility (the affected chain entries cannot be re-run to evaluate AI nondeterminism without a recorded model identifier). NOT a chain-integrity finding. | (audit-procedures P-25 stratification; MRM-program finding) |
| (verifier output) | Verifier output format mismatch | Not an institutional finding — verifier-implementation defect. Obtain a conforming verifier from the project supply chain and re-run. The institution is unaffected. | (none — examiner-side; vendor toolchain finding) |
| 3 | `chain_kind out of v1 enumeration at seq N` | MRA; the SDK emitted a `chain_kind` value not in the v1 closed enumeration. Investigate the SDK build and remediate the affected entries. NOT a chain-integrity finding (the chain still cryptographically binds whatever was captured); a control-completeness finding for the chain's documented attribute discipline. | Scenario 1 |
| 11 (Pattern B) | `Pattern B partition fields missing on multi-seal-day` | Severe MRA; the institution operated within-day algorithm rotation under Pattern B (spec §10.10.2) but at least one of the day's two seal records lacks `covers_received_at_min` / `covers_received_at_max`. Investigate the seal job's Pattern B implementation. | Scenario 5 variant (algorithm-rotation procedural gap) |
| (audit-procedures P-35) | ECOA translation-attribute completeness gap | Medium MRA when the gap rate is below the institution's documented threshold; High MRA when the gap rate exceeds the threshold OR when the gap pattern correlates with specific languages, translator kinds, or time windows. **Control-completeness finding, NOT chain-integrity finding** — the chain still cryptographically binds whatever was captured; the gap is in what was captured. | (audit-procedures P-35; MRM-program finding) |
| 9 / 8 | Constant-time discipline violation discovered at audit time | Severe MRA when the institution's verifier or SDK uses a non-constant-time equality primitive for the `key_fingerprint` check or the `payload_hash` MAC check (spec §10.8 MUST). Investigate the implementation and remediate to a constant-time primitive (`hmac.compare_digest`, `CryptographicOperations.FixedTimeEquals`, `subtle.ConstantTimeCompare`). | Scenario 1 + vendor-toolchain investigation |
| (audit-procedures P-37) | Multi-region replication-completeness gap (Pattern A discrepancy non-zero without documented tolerance) | Medium control-completeness when the gap is bounded and documented at examination time; High control-completeness when the gap is unbounded, undocumented, or the institution's CC8.1 multi-region procedure does not explain the disposition. **Control-completeness finding, NOT chain-integrity finding** — the chain still cryptographically binds whatever was captured in the seal region's ledger; the gap is in source-region events that did not replicate to the seal region by seal-time. | (audit-procedures P-37; CC8.1 multi-region replication procedure) |

The table has 20 informational rows covering the 12 spec §7 steps PLUS the file-pre-flight truncation refusal PLUS the anomaly-only rows PLUS the v1.0-final close-out rows (verifier output format, `chain_kind` enumeration, Pattern B partition fields, ECOA translation completeness, constant-time discipline, multi-region replication-completeness). The "12-step procedure" terminology refers to the spec §7 step count; the table grows because some steps (11 dual-algorithm sub-cases, 12 cadence vs dev-mode) emit multiple distinct reason strings, the anomaly rows are not §7 steps, and the v1.0-final close-out rows cover failure modes the normative material introduces. The `step` column maps each row to the spec §7 step that produced it (or `(file pre-flight)` / `(anomaly)` / `(verifier output)` / `(audit-procedures P-XX)` for non-step rows).

## Sample finding paragraphs

### Header HKDF inputs do not match running v1 inputs (spec §7 step 2)

> The chain verifier reported `header HKDF inputs do not match running v1 inputs` on [N] audit files for tenant [T] during the period. The verifier recomputed `expected_hkdf_inputs_digest = SHA-256(HKDF_SALT || info_for_tenant || length_LE32)` for the file's tenant_id and constant-time-compared against the file's `header.hkdf_inputs_digest`; the comparison failed. This is a Severe finding consistent with one of: (a) the institution's SDK shipped with HKDF constants different from the v1 spec constants (format-construction defect); (b) the institution's verifier is running the wrong constants (verifier-build defect — the institution should obtain a verifier built against the same v1 constants); (c) the file's header was tampered with after writing. The institution must identify which root cause applies and remediate accordingly.

### Header genesis_hash does not match v1 constant (spec §7 step 3)

> The chain verifier reported `header genesis_hash does not match v1 constant` on [N] audit files for tenant [T] during the period. The v1 spec mandates 32 zero bytes as the file's genesis hash; the file's recorded genesis differs. This is a Severe format-construction defect indicating the institution's SDK is producing chain entries whose first event's `prev_hash` is not the v1 zero-bytes constant. The institution must investigate the SDK build and re-issue the affected period from a v1-conformant SDK.

### Format version mismatch at entry (spec §7 step 5)

> The chain verifier reported `format_version mismatch at seq N` on [N] events for tenant [T] during the period. The file's header declared `format_version: v1` but one or more chain entries within the file declared a different `format_version`. This indicates mid-file format drift — typically an SDK that produces entries under multiple format versions in error, or chain entries that were copied from a different file with different format. The institution must investigate the SDK's format_version-handling code and remediate the affected entries.

### Chain link broken at seq N (spec §7 step 6)

> The chain verifier reported `chain link broken at seq N` on [N] events across [M] runs for tenant [T] during the period. The structural walk detected `entry.prev_hash != expected_prev_hash` (the previous entry's `payload_hash` did not match the current entry's `prev_hash`) OR `entry.seq != expected_seq` (the sequence integers were out of order). This is a High-severity finding consistent with insertion or deletion of chain entries between writing and verification — the most likely tampering signal that does not require key access. The institution must complete root-cause analysis and report findings to its primary regulator within [X] days.

### Chain hash mismatch — payload_hash MAC mismatch (spec §7 step 9)

> The chain verifier reported `payload_hash MAC mismatch at seq N` on [N] events across [M] runs for tenant [T] during the period. Mismatches indicate that the recomputed HMAC-SHA-256 over `(expected_prev_hash || canonical_bytes)` under the looked-up IKM's session key does not match the stored `payload_hash`. This is a Severe content-tampering finding consistent with either an in-flight modification of the captured event or a defect in the chain construction. Bank management's investigation [is/is not] underway. The institution should complete root-cause analysis and report findings to its primary regulator within [X] days.

### Merkle root mismatch (single day)

> The chain verifier reported a Merkle root mismatch on [date]. The recomputed Merkle root from the day's events does not match the HSM-signed root recorded for the day. This indicates either alteration of one or more events between the time of capture and the verifier's read, or a defect in the verifier's reconstruction. The institution must produce evidence of the original events and re-run the verifier; if the mismatch persists, the institution must treat the day as unverifiable and follow the incident-response playbook for chain-detected events.

### Sealing delay without notification

> The chain verifier reported [N] days where the daily Merkle seal was recorded more than 72 hours after UTC midnight of the seal day. Per spec design (`04-hsm-custody.md` §5.2), institutions are expected to notify their primary regulator when the daily seal is delayed beyond 72 hours. The institution did not notify. The institution should review its operational-monitoring posture and document a notification procedure for future seal-delay events.

### Signature verification failed

> The chain verifier reported that the Ed25519 signature on the daily seal for [date] failed verification against the institution's published public key. This is a severe finding consistent with one of the following: (a) corruption of the seal record, (b) a public-key mismatch between the registry and the verifying key, (c) substitution of the seal by an unauthorized party, or (d) compromise of the institution's HSM signing key. The institution must immediately investigate, rotate the signing key if compromise is confirmed, and notify both the regulator and any third parties whose verification depends on the affected key.

### Software-key fallback used in production

> The chain verifier reported that [N] daily seals during the examination period carry the dev-mode marker, which indicates the seal was signed with a software-backed Ed25519 key rather than an HSM-backed key. Software-key fallback is acceptable for development and CI use only; production seals must be signed in HSM custody. The institution must immediately investigate how dev-mode was enabled in the production environment, restore HSM-backed signing, and re-issue any affected seals from an HSM-backed key.

### Late-binding event rate exceeds threshold

> The chain verifier reported a late-binding event rate of [X.Y]% during the examination period, exceeding the threshold of [Z]%. Late-binding indicates events that arrived at the ledger after their day's seal had already been computed; included in the next day's seal. A high rate is an operational signal — typically network latency or SDK persistence delay — not a tampering signal. The institution should investigate the operational cause and document remediation.

### Key fingerprint mismatch (spec §7 step 8)

> The chain verifier reported `key_fingerprint mismatch` on [N] events on [date(s)]. The fingerprint check is the verifier's pre-flight identity binding: the looked-up IKM for `(tenant_id, key_version)` is hashed against the recorded `tenant_id` and the truncated SHA-256 is compared against the entry's stamped `key_fingerprint`. A mismatch indicates the institution's IKM lookup returned bytes that do not produce the per-entry fingerprint — botched key rotation that re-used `key_version=N` for a different IKM, restored backup pointed at the wrong tenant row, swap of tenant rows in the registry, or active cross-tenant configuration drift. **This is an identity mismatch, not a content tampering**; it is investigated against the institution's IKM roster and the most recent `master.reconciliation_completed` event (spec §10.1), not against the chain content itself. The institution must reconstruct the correct IKM-roster state, reconcile against the affected `(tenant_id, key_version)` triples, and re-run the verifier under the corrected roster. The mismatch is detected before any MAC compute (spec §7 step 8); the affected events' MAC integrity is independently confirmable once the IKM is restored.

### Unknown key_version (spec §7 step 7)

> The chain verifier reported `unknown key_version: no IKM for (tenant=T, key_version=V)` on [N] events on [date(s)]. The institution's IKM registry returned no IKM for the `(tenant_id, key_version)` pair stamped on the entry. Plausible root causes: (a) the institution rotated its IKM and did not register the new generation in the verifier's IKM registry; (b) the institution withdrew an old IKM from the registry while events stamped with that `key_version` still exist (premature retirement, prohibited by spec §10.9); (c) the verifier's IKM registry is stale relative to the institution's actual operations; (d) the entry's `key_version` field has been tampered with. The institution must retain its IKM-roster complete for the retention window of the events. The institution must investigate which root cause applies, restore the IKM if it was prematurely retired, and document the reconciliation.

### Format version not supported (spec §7 step 1)

> The chain verifier reported `format_version not supported by this verifier (running v1)` on the institution's audit-file header. **This is not an institutional finding.** The institution upgraded its SDK to a future spec version and the examiner is running an older verifier. The examiner resolves the finding by obtaining the newer verifier from the project supply chain (per `regulator-pack/deployment-package.md`) and re-running. No examination action against the institution is required.

### Audit file ends mid-line (mid-write truncation, spec §4.1)

> The chain verifier refused to verify the institution's audit file for [date] because the file's last byte is not a newline (`\n`). Per spec §4.1 mid-write truncation refusal, the verifier rejects truncated files rather than silently passing chains that lost their last entry. **This is an operational finding (writer-side crash recovery), NOT a tampering finding.** The institution's SDK process appears to have died mid-append; the chain content is intact through the last complete entry but the truncated tail is non-conformant. The institution should consult its incident log for the affected period, recover the dropped event from the SDK's local SQLite buffer or an upstream OTLP backend if available, and document the gap if recovery is not possible. Severity is sealing-delay-equivalent (Observation; MRA on persistent recurrence indicating systemic crash-loop).

### Cross-chain lift detected (spec §7 step 4)

> The chain verifier reported `cross-chain lift detected at seq N` on the institution's audit file for [date]. An event in the file claims a different `tenant_id` or `run_id` than the file header — the file declares chain `(tenant=T, run_id=R)` but the event declares `(tenant=T', run_id=R')`. Two plausible root causes: (a) the institution mis-bundled audit snapshots when producing the examination evidence (operational, low-severity finding for evidence-handling); (b) an attacker attempted to lift a valid event from another chain into this one to forge cross-chain provenance (integrity, severe finding). The institution must produce evidence of the event's original chain provenance and, in case (b), follow the institution's incident-response playbook for chain-detected events.

### Cadence mismatch (spec §7 step 12)

> The chain verifier reported `cadence mismatch` for the institution's seal record(s) on [date(s)]. The seal record's `cadence` field (`hourly` | `daily` | `weekly`) does not match the cadence the institution claims in its control description. This is a control-description-accuracy finding: the institution's documentation does not match its operations. The institution must reconcile its control description with the operating cadence — either by updating the control description to match the operating cadence (if the operating cadence is correct and the documentation drifted) or by adjusting the operating cadence to match the documented commitment (if the documentation is correct and the operations drifted; for cadence relaxation, the institution must obtain examiner approval per spec §4.2.1 before adjusting). **Severity is Observation on first occurrence within a regulator standard examination cycle; MRA on repeat across consecutive examination cycles** (the regulator's standard cycle is typically 12-18 months for IT examinations; the repeat threshold is two consecutive cycles with the same cadence-mismatch root cause).

### Verifier output format mismatch (spec §7 normative output format)

> The chain verifier the institution presented at the examination produced output that does not conform to the spec §7 normative output format. The expected format for a failed run is three lines — `Status: FAIL`, `Step: N`, `Reason: <text>` — with exact field labels (capitalization, colon, single-space separator) and `0x0A` line terminators. **This is not an institutional finding.** The institution presented a verifier whose output format the examiner could not parse mechanically; the resolution is to obtain a conforming verifier from the project supply chain (per `regulator-pack/deployment-package.md`) and re-run. The institution is unaffected; the finding is against the verifier's vendor toolchain. The examiner re-runs verification with the conforming verifier and any chain-integrity findings produced by the conforming verifier are documented separately under their respective spec §7 step entries.

### `chain_kind` enumeration violation (spec §3 closed enumeration)

> The chain verifier reported `chain_kind out of v1 enumeration at seq N` on [N] events for tenant [T] during the period. Spec §3 normates a closed enumeration for the `chain_kind` field — `audit` | `model_call` | `tool_call` | `routing` | `translation` | `operational`. The institution's SDK emitted a value outside the enumerated set. This is a control-completeness finding for the chain's documented attribute discipline: the chain still cryptographically binds whatever was captured (the per-entry MAC, the daily Merkle seal, and the HSM signature are unaffected), but the SDK is producing chain entries whose `chain_kind` discriminator does not conform to v1. The institution must investigate the SDK build, identify where the non-conforming `chain_kind` value originated (typically a vendor SDK extension or a custom chain decorator that did not enforce the closed enumeration), and remediate the affected SDK to emit only enumerated values. The affected entries remain integrity-bearing under the chain's cryptographic property; the discriminator gap is a maintainability concern that prevents downstream consumers from mechanically separating event classes per the documented enumeration.

### Pattern B partition fields missing on multi-seal-day (spec §10.10.2)

> The chain verifier reported `Pattern B partition fields missing on multi-seal-day` for the institution's seal records on [date]. The institution operated within-day algorithm rotation under Pattern B (per spec §10.10.2) — two seal records share the seal-date and split the day's events by capture time — but at least one of the day's two seal records lacks the `covers_received_at_min` / `covers_received_at_max` partition fields the verifier requires to reconstruct each seal's event-subset boundary. This is a Severe finding: without the partition fields, the verifier cannot determine which events each seal covered, and the day's chain cannot be mechanically validated under the Pattern B convention. The institution must investigate its Pattern B seal-job implementation, confirm both seal records carry the `covers_received_at_min` / `covers_received_at_max` fields, confirm the two windows are contiguous and non-overlapping (together covering the full UTC day from `00:00:00.000000Z` inclusive to the next day's `00:00:00.000000Z` exclusive), and re-issue affected seals from a Pattern B-conformant seal job. Pattern A (cosigned same-day seals) is the RECOMMENDED alternative for institutions whose tooling supports dual-algorithm cosign; the institution may evaluate switching to Pattern A as a remediation.

### ECOA translation-attribute completeness gap (audit-procedures P-35)

> The institution's SOC engagement (or the examiner's direct sample) found [N] ECOA translation chain entries during the period missing one or more REQUIRED attributes per spec §10.11. The REQUIRED attributes are `audit.ecoa.translation.target_language`, `audit.ecoa.translation.translator_kind`, and `audit.ecoa.translation.output_hash`; `audit.ecoa.translation.translator_id` is REQUIRED conditionally when `translator_kind = "llm"`. **This is a control-completeness finding, NOT a chain-integrity finding.** The chain still cryptographically binds whatever was captured (the per-entry MAC, the daily Merkle seal, and the HSM signature are unaffected); the gap is in what was captured. Severity is Medium when the gap rate is below the institution's documented threshold for routine variance; High when the gap rate exceeds the threshold OR when the gap pattern correlates with specific languages, translator kinds, or time windows in a way that suggests systematic rather than incidental loss. The institution must investigate its ECOA translation pipeline, identify where the REQUIRED attributes were dropped (typically a chain decorator that did not capture all attributes from the translation step's output, or a translation-pipeline configuration that elided certain attributes), and remediate the pipeline. The institution's CFPB / ECOA examiner cross-checks the chain entries against the customer-side delivery record per `customer-dispute-procedures.md` §"Adverse-action notices (ECOA)"; gaps that prevent the examiner from answering "did this customer receive the adverse-action notice in their preferred language within the regulatory window?" from the chain alone are escalated per the institution's ECOA-supervision posture.

### Constant-time discipline violation (spec §10.8 MUST)

> The institution's [SOC engagement / vendor-conformance attestation / examiner verifier review] found that the institution's verifier (or SDK) uses a non-constant-time equality primitive for the `key_fingerprint` check (spec §7 step 8) or the `payload_hash` MAC check (spec §7 step 9). Spec §10.8 makes constant-time comparison a normative MUST (changed from RECOMMENDED in v1.0-final). The conformant primitives are the platform's stdlib helpers — `hmac.compare_digest` in Python, `CryptographicOperations.FixedTimeEquals` in .NET, `subtle.ConstantTimeCompare` in Go — or an equivalent that has been measured under a CI harness asserting compare-time independence from the position of the first byte mismatch. Naive byte-equality (`bytes.Equal`, `==` on Python `bytes`, `Span<byte>.SequenceEqual` without the cryptographic-equality variant) leaks early-mismatch position via timing and lets an attacker forge MACs against a running verifier. This is a Severe finding: the institution's chain-integrity verification is operating without the timing-side-channel discipline the spec requires. The institution must investigate the implementation, replace the non-conforming comparison with a constant-time primitive, and validate the fix under the project's conformance corpus. The verifier-toolchain's vendor is responsible for shipping the fix; the institution is responsible for ensuring its deployed verifier carries the conformant primitive.

### Multi-region replication-completeness gap (audit-procedures P-37; spec §10.15 Pattern A invariant 5)

> **Multi-region replication-completeness gap.** The institution operates spec §10.15 Pattern A (active-active multi-region) and reports `discrepancy != 0` on the `master.cross_region_replication_completed` operational event without an institution-documented tolerance explaining the gap. SOC procedure P-37 confirms cross-region event counts do not reconcile against the seal region. Finding language: "Multi-region replication-completeness — the institution's per-region event-count reconciliation discloses a discrepancy between events emitted in source regions and events replicated to the seal region. The chain integrity holds (the seal accurately seals what's in the seal region's ledger), but the institution's CC8.1 multi-region procedure does not document the disposition of the discrepancy. Control finding; remediation: institution documents the replication-loss disposition and either tightens replication SLAs or expands the institution's tolerance statement with a documented rationale."

## Counsel-rebuttal framing for high-severity findings

For the most-severe findings the institution's counsel is most likely to rebut, the examiner SHOULD pre-frame the response. Common rebuttal patterns and the examiner's prepared response:

**For `key_fingerprint mismatch` (spec §7 step 8) — counsel argues "this is just a key-management hiccup, not a security finding":**

> The verifier's check at step 8 fires before any MAC compute. The match-or-mismatch determination is purely arithmetic: the looked-up IKM either produces the recorded fingerprint or it does not. A non-matching fingerprint means the institution's IKM-roster entry for `(tenant_id, key_version)` is wrong; whether that wrong entry resulted from a misconfiguration (low severity) or unauthorized substitution (high severity) is the institution's investigation question, NOT the verifier's. The IR Scenario 7 triage tree is the institution's tool for classifying root cause; until the institution completes the triage, the finding remains Severe pending classification.

**For `payload_hash MAC mismatch` (spec §7 step 9) — counsel argues "this is an SDK defect, not tampering":**

> SDK defect is one of three plausible root causes; the verifier does not distinguish among them at step 9. The institution must complete root-cause analysis (per IR Scenario 1) and produce evidence supporting its determination. An SDK-defect determination requires evidence of the defect's reach, the affected period, and the remediation; in the absence of that evidence, the finding remains Severe pending the institution's analysis.

**For `signature verification failed` (spec §7 step 11) — counsel argues "the registry has the wrong key, not the HSM has been compromised":**

> A registry-mismatch root cause is one of four plausible IR Scenario 3 root causes (registry mismatch, seal record corruption, unauthorized substitution, HSM signing-key compromise). The institution must produce evidence ruling out the latter three; in the absence of HSM-side audit-log evidence supporting the registry-mismatch hypothesis, the finding remains Severe pending HSM-side investigation.

**For `cold_dr.dryrun_attestation_consumed = FAIL` or missed (IR Scenario 11) — counsel argues "this is a project-side governance issue, not the institution's responsibility":**

> The institution's CUEC-IR-05 commits the institution to consume the project's annual dry-run attestation. The institution's failure to consume (whether due to project-side delay or institution-side oversight) leaves the institution's cold-DR fallback path unexercised, which is exactly the residual the institution's CUEC commits the institution to close. The finding is institution-side; the institution's response includes coordination with the project to obtain the missing attestation.

The counsel-rebuttal framing is preparation, not a script. Counsel will adapt; the examiner adapts in turn. The framing exists so the examiner is not making the response up at the moment of rebuttal.

## CAT (FFIEC Cybersecurity Assessment Tool) crosswalk scope note

This finding-language document is scoped to chain-of-custody examinations under the FFIEC IT Examination Handbook (per `regulator-pack/handbook-mapping.md`). Examiners working CAT-track examinations (FFIEC Cybersecurity Assessment Tool, Inherent Risk Profile + Cybersecurity Maturity assessment) use this document as a reference for chain-detected events but follow the CAT's own finding-language conventions for CAT-specific assessment domains. A formal chain → CAT crosswalk is a v1.x candidate; until then, examiners working both tracks consult the FFIEC IT Handbook crosswalk in `regulator-pack/handbook-mapping.md` for the chain's primary mapping and adapt CAT-specific language as needed.

## Closing-language patterns

For findings that the institution remediates during the examination:

> The institution remediated this finding during the examination by [specific action]. The verifier's subsequent run on [date] confirmed the issue was resolved. No further action is required.

For findings that require post-examination remediation:

> The institution must complete remediation by [date] and provide documentation including [list]. The verifier output for the remediated period will be reviewed in the next examination cycle.

For findings escalated to enforcement consideration:

> Given the severity of [finding type] and the institution's [contributing factors], this finding is referred to the [regulator's enforcement function] for further action.

## Repeat-finding language

When a finding from a prior examination remains open or recurs at a subsequent examination:

> The institution had a related finding at the [date] examination ([prior finding identifier]). The current examination finds [same/similar pattern]. The institution's prior remediation [was insufficient/has not yet been completed]. Given the recurrence, the finding is escalated to [next-tier severity / enforcement track].

For findings closed at the prior examination but recurring at this one:

> The institution previously remediated a related finding at the [date] examination. The current examination finds the issue has recurred. The institution's [process / control] that supported the prior remediation has [drifted / been disabled / failed]. The institution must investigate the cause of recurrence and report to its primary regulator within [X] days.

## Public-disclosure language

For findings that flow into public-disclosure documents (consent orders, formal agreements, civil money penalty orders):

> The institution failed to maintain effective controls over the integrity of artificial-intelligence-driven decisions captured by the institution. Specifically, [the institution did not maintain effective integrity verification of its AI agent activity logs / the institution's daily seal-signing operations were materially delayed without notification to the institution's primary regulator / similar specific finding]. The institution shall [specific remediation requirements].

Public-disclosure language is concise. The technical detail is in the institution's response document and the underlying examination report; the public-disclosure document names the failure and the required action.

## Enforcement-action lifecycle

For institutions under formal enforcement action (formal agreement, consent order, etc.) where a chain finding is part of the action:

**Initial action language:**

> The institution shall maintain effective controls over the integrity of AI-driven decisions, including operating the [chain control] in accordance with the institution's documented control description. The institution shall provide the [primary regulator] with the verifier output for the institution's chain operation on a [monthly / quarterly] basis.

**Monitoring language:**

> The institution's chain operation is monitored monthly via verifier output. Continuous pass results during the monitoring period support remediation; persistent anomalies indicate ongoing concerns warranting continued monitoring.

**Removal-of-action language:**

> The institution has consistently demonstrated effective chain control operation over the [period]. Verifier output for [N consecutive months] shows pass results with no material anomalies. The institution's CUEC operation is documented and tested; SOC reporting (if applicable) confirms the controls operate effectively. The [enforcement action] is recommended for removal.
