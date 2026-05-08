# Round 11 — Big Four SOC Audit Senior Manager review

**Reviewer.** Camille Lefèvre, Senior Manager, Big Four SOC Audit Practice (US bank-tech, prior EU fintech). Engagement-quality reviewer on five SOC reports for cryptographic-control implementations.

**Reading angle.** Auditability, testability, evidence-trail completeness, user-entity audience consumption, mechanical SOC procedure operability.

**Materials read.**

1. `spec/chain-of-custody-v1.md` — focus §7 verification, §10 operational requirements
2. `docs/audit-procedures.md` — focus updated P-3 + new P-22..P-25
3. `docs/soc-pack/control-evidence-events.md` — focus the two new events plus vocabulary clarification
4. `docs/soc-pack/section-4-template.md` — focus updated Procedures section
5. `docs/control-map/CUECs.md` — focus CUEC-CRY-04 (round-10 fix claim)
6. `docs/control-map/TSC-mapping.md` — focus CC6.1, CC6.7, CC6.8, CC7.2, CC8.1 updates
7. `docs/anomaly-documentation-template.md` — focus higher-bar "operationally explained" + worked example
8. `docs/user-entity-summary.md` — focus added paragraph on the two new events
9. `docs/regulator-pack/sample-report.md` — focus rotation/FAIL day, methodology, --strict, TSC appendix

I have not read prior-round feedback.

---

## Five questions I would ask the engagement team

1. The four-piece evidence package for `key_fingerprint mismatch` is the strongest part of the entire corpus from a SOC-evidence standpoint. The template's worked example threads exhibit identifiers (`CM-2026-04-12-039`, `CM-2026-04-15-002`), reconciliation week numbers, and explicit fingerprint values. I have signed off on Section 4 starter exhibits with weaker grounding than this template. **Is the worked example the institution's responsibility to adapt verbatim, or does the project intend it to read as normative-quality?** The text says "template-quality, not normative-quality" — fine, but the caveat is buried. A SOC team picking up this template should know the structure is the contract, the specific values are illustrative.

2. The `audit_file.truncation_detected` `recovery_outcome` enum (`complete` | `partial` | `unrecoverable`) is now consistent across spec §4.1, the operational-event JSON, the anomaly template, P-24, the Section-4 procedures bullet, and the user-entity summary. **Is the enum value-set frozen?** A future spec patch that adds a fourth value (e.g., `recovered_to_known_good_subset`) would need to ripple through six documents. A short note in the spec change log about the enum's stability commitment would close the maintenance question.

3. P-23 names two acceptable root causes for `hkdf_inputs_digest mismatch`: documented constants change OR multi-region deployment-drift. The multi-region case assumes the institution's working papers carry per-region SDK version inventory. **Does the SOC engagement scope require that inventory as standing evidence**, or is the SOC team meant to request it on an ad-hoc basis when the anomaly fires? My experience says ad-hoc requests in the middle of fieldwork eat days; a standing artifact (e.g., monthly SDK-version-by-region report) would let the SOC team confirm the explanation in minutes.

4. The Section 4 Procedures bullet on IKM retirement reads "IKM retirement on `[institution-defined]` cadence per spec §10.9; retirement requires `chain_entries_referencing_remaining = 0`." That bullet is correct on the math. **Is there a Section 4 review checklist item that confirms the institution's IKM-retirement procedure was actually executed during the period** (vs documented but never run)? The `master_key.retired` event provides the evidence; I would want a Section G note for any retirements during the period and a P-procedure that pulls those events. Right now retirement is named in Procedures (Section C) but doesn't appear in the audit-procedures inventory by name. P-5 covers IKM custody and rotation but does not call out retirement specifically.

5. CC6.7 and CC8.1 in TSC-mapping cite `master_key.retired` as evidence for the IKM-lifecycle. CC6.1 cites the same event for RBAC. **Is the same operational event being credited against three TSC criteria intentional**, and if so, is the SOC team expected to triangulate the same evidence to three findings? Triangulation is fine when the criterion shapes are distinct (logical-access-administer is a different shape from change-management); it is double-counting when the criterion shapes overlap. A short note in TSC-mapping clarifying the distinct evidence each criterion expects would help.

---

## Findings

### Status

| ID | Severity | Status | Area | Summary |
|---|---|---|---|---|
| F-1 | Gap | Open | CUECs.md | CUEC-CRY-04 still reads `session-key-id reconciliation`; the brief asserts round-10 fixed this. The drift is unfixed in the shipped artifact. |
| F-2 | Partial | Open | sample-report.md TSC appendix | Late-binding and `plaintext-` prefix rows are present; FAILURE-RECORD rows for steps 1–7 (other than step 8) and step 12 (cadence + dev-mode) are absent; truncation-refusal failure is also absent. The appendix is incomplete relative to the methodology section's own §7 step inventory. |
| F-3 | Partial | Open | CUECs.md cross-references | Even after F-1 is fixed, three other documents (`audit-committee-summary.md`, `at-scale-operations.md`, `legal-disclosure.md`) and one design doc (`design/06-ledger-server-design.md`) still use `session-key-id reconciliation`. The vocabulary is inconsistent across the corpus, not just CUECs.md. |

### F-1 — CUEC-CRY-04 vocabulary still drifts from the spec

**Severity.** Gap (the brief asserted this was fixed in round 10; it was not).

The brief introducing this round states "confirm CUEC-CRY-04 vocabulary now uses key-fingerprint reconciliation (round 10 fixed drift)." I confirmed the opposite. CUECs.md line 28 reads:

> The institution operates session-key-id reconciliation at a documented cadence (no more than weekly per spec recommendation).

Spec §10.1 and audit-procedures P-6 use "key-fingerprint reconciliation." The `session_key_id` field was removed in the v1.0-rework per the spec change log; the surviving identity-binding field is `key_fingerprint`. The CUEC is referring to a field that no longer exists in the normative spec.

**Why this matters mechanically.** A SOC team reading the CUEC pulls the institution's procedure documentation looking for "session-key-id reconciliation." The institution's procedure (which tracks the spec) calls the procedure "key-fingerprint reconciliation." The SOC team has to reconcile vocabulary mid-engagement before they can test the control. That is exactly the kind of friction the SOC pack exists to eliminate.

**Fix.** Replace CUECs.md line 28 with:

> The institution operates **key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1). Reconciliation evidence is the `master.reconciliation_completed` operational event with `fingerprint_unmatched_count` field.

Cross-reference P-6 by name. Update the "Mapping to chain primitives" footer (line 121) — `CRY-04` is correctly cited, but the table reader will hit the stale text.

### F-2 — Sample-report TSC appendix is incomplete relative to the §7 step inventory

**Severity.** Partial.

The TSC appendix at sample-report.md lines 305–321 carries the new rows the brief called out — Late-binding count (line 320) and `KMS handle URI` plaintext-prefix (line 321) — both correct. However, the appendix's FAILURE-RECORD coverage is uneven:

**Present in appendix.**
- step 8: key_fingerprint mismatch
- step 9: payload_hash MAC mismatch
- step 10: merkle root mismatch
- step 11: signature verification failed

**Absent from appendix despite being in the methodology section (lines 180–198).**
- step 1: format_version not supported (PI1.1 / CC8.1 — format-dispatch failure)
- step 2: hkdf_inputs_digest mismatch (CC7.2 — already cited as a non-failure row at line 312, but not as a failure row)
- step 3: genesis_hash mismatch (PI1.1 / CC8.1)
- step 4: cross-chain lift detected (CC6.1 — tenant/run binding)
- step 5: format_version mismatch at entry (CC8.1)
- step 6: chain link broken / seq out of order (PI1.1 — link-walk failure)
- step 7: unknown key_version (CC6.1 + CC7.2 — registry-retention failure; this is the IR Scenario 8 failure mode and warrants its own row)
- step 12: cadence mismatch under non-strict (anomaly) and dev_mode under --strict (FAIL — already cited at line 321 as a non-failure-record row, but not as a FAIL-day row)

**Also absent.** `audit file ends mid-line — possible mid-write crash` (file-pre-flight FAIL per spec §4.1). The methodology section names it (line 209); the appendix does not credit it to PI1.2 / CC7.5 / CC8.1.

**Why this matters mechanically.** When my engagement team consumes the sample report and a real institution's verifier emits a step-7 failure, the SOC team has no row in the cross-walk to allocate the evidence. They invent one mid-engagement, which means the next engagement team invents a different one, which means the firm's working-paper inventory is inconsistent across institutions. The cross-walk is the artifact that prevents that. Half-coverage forces the team to do the missing half ad-hoc.

**Fix.** Add rows for steps 1, 3, 4, 5, 6, 7, 12, and the truncation-refusal FAIL. The mapping for step 7 specifically (unknown_key_version) should include the IKM-retention CUEC reference (CUEC-CRY-05 + spec §10.9) since this failure mode is the spec's evidence path for the retention rule.

### F-3 — Vocabulary drift extends beyond CUECs.md

**Severity.** Partial.

Even after F-1 is fixed, the term `session-key-id reconciliation` appears in:

- `audit-committee-summary.md` line 36 (Audit Committee evidence list)
- `at-scale-operations.md` lines 102, 116 (operational scaling guidance)
- `legal-disclosure.md` line 48 (court-ordered-disclosure record list)
- `design/06-ledger-server-design.md` line 330 (operational events emission table; refers to "session_key_id_count, unmatched_count")

The first three are user-facing documents that committee chairs, ops engineers, and legal counsel read. A search-and-replace across the corpus closes this. The design doc (line 330) is the most concerning — it's the schema reference that drove the CUECs.md text, and it still names `session_key_id_count` as a field on the `master.reconciliation_completed` event. The shipped event (per `control-evidence-events.md` lines 280–294) carries `key_versions_observed`, `key_fingerprints_observed`, and `fingerprint_unmatched_count` — no `session_key_id_count` anywhere.

**Fix.** Search the corpus for `session-key-id` and `session_key_id` (both spellings) and replace with the current vocabulary (`key-fingerprint reconciliation` / `key_fingerprint`). The design doc at `06-ledger-server-design.md` §7.3.1 needs the schema fields refreshed to match the shipped event.

---

## What is solid

I want to be specific about what I would NOT raise in a SOC engagement, because the artifacts shipped this round are unusually good.

- **The four-piece evidence package for `key_fingerprint mismatch`** (anomaly-template lines 121–130; P-22 lines 155–161) is the highest-quality SOC-evidence specification I have read for a cryptographic control. The four pieces correspond to four independent failure modes (roster row identification → change approval → re-verification → reconciliation cross-check) and the absence of any one piece tells the SOC team something specific about which control surface is weak. That is exactly the shape an evidence package should take.

- **The worked example anomaly record** (anomaly-template lines 147–221) eliminates the friction that usually surrounds the first time an institution writes one of these. The cited identifiers (`CM-2026-04-12-039`, `ANM-2026-04-15-001`, week-15 reconciliation) read like a real institution's records, which means a real institution can use the example as a copy-and-adapt template. I have asked for exactly this kind of friction-reducer in past engagements; it is rare to see it shipped.

- **The methodology section in the sample report** (lines 167–210) walks the §7 twelve-step procedure with explicit pass/fail rules per step. The "NO MAC compute on miss" and "NO MAC compute on mismatch" annotations on steps 7 and 8 (lines 189, 191) call out the load-bearing defense properties at exactly the spot where a SOC reviewer needs to confirm them.

- **The `master_key.retired` event** (control-evidence-events.md lines 258–274) carries `chain_entries_referencing_remaining` as a control assertion. That is the right field shape — a non-zero value at retirement time IS the control failure, and the SOC team can test it mechanically by querying the event log for retirement events with non-zero values. No interpretation required.

- **P-22 through P-25** are operationally testable as written. P-22's four-piece checklist matches the anomaly template's four-piece evidence requirement byte-for-byte; P-24's `complete` / `partial` / `unrecoverable` enum matches the operational event's enum byte-for-byte; P-25 distinguishes chain-integrity findings from MRM-program findings, which is the right boundary for a SOC team to draw. I would issue these procedures to my team without modification.

- **TSC-mapping CC7.2 and CC8.1 updates** correctly credit the new events to monitoring (CC7.2) and change management (CC8.1). The CC6.1 credit for `master_key.retired` is reasonable as RBAC-administer evidence (the retirement is an authorized administrative action). See question 5 above on the triple-credit.

- **The user-entity summary's audience-appropriate paragraph** on the two new events (lines 28–34) is calibrated correctly: it names what the event is, names the response procedure, names where the recovery outcome is recorded. A downstream business partner's auditor can read that paragraph and know what reliance is reasonable. No cryptographic detail leaks; no operational detail is omitted.

- **Section 4 starter additions** (procedures bullets at lines 56, 58, 59) cover IKM retirement, mid-write-truncation recovery, --strict cadence, and slsa-verifier deployment-gating. The IKM-retirement bullet correctly conditions retirement on `chain_entries_referencing_remaining = 0`. The slsa-verifier bullet is the right level of detail for Section 4 (deployment-gating control with archived output log).

---

## Roll-up

**Gaps:** 1 (F-1 — CUEC-CRY-04 vocabulary drift, asserted-fixed but not actually fixed)
**Partials:** 2 (F-2 — sample-report TSC appendix incomplete; F-3 — vocabulary drift extends beyond CUECs.md to four other documents)

**Stopping criterion (0 gaps + 0 partials).** Not met. Three open items, all small, all mechanical.

**Path to clean.** F-1 is a one-line edit. F-3 is a search-and-replace across four documents plus a schema refresh on `06-ledger-server-design.md` §7.3.1. F-2 is adding eight rows to the sample-report TSC appendix and one row for the truncation-refusal FAIL. None of the three is structural.

**Engagement-grade verdict.** With F-1 and F-3 fixed (one edit + a search-and-replace) and F-2 closed (eight rows added to one appendix), the SOC pack is at the level I expect for issuance. The procedures match how my team would actually test; the four-piece evidence package and worked example are above the floor; the new operational events are correctly credited to TSC. The remaining work is mechanical, not structural.
