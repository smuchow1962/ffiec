# 03 — Big Four SOC Audit Engagement Partner (round 9)

> **Role.** SOC engagement and CUEC/TSC mapping lens. First-look review of the v1.0-rework corpus.

## Persona — Marcus O'Brien

**Background.** Engagement Partner at a Big Four firm; 25 years issuing SOC 1 Type II and SOC 2 Type II opinions for fintech and bank-tech vendors. Most of my recent docket sits in cryptographic-control attest engagements where the substance has to be testable mechanically against the AICPA Trust Services Criteria (CC, A, PI, C, P) and where the user-entity audience usually pulls the report into an ICFR or vendor-management workstream. I have spent more time than I want arguing with vendors about whether their "tamper-evident" claim is actually attestable, and I am reading this corpus with that scar tissue.

**Reading angle.** Auditability. Three concurrent tests on every control claim:
1. Is the claim worded narrowly enough that a SOC procedure can produce yes/no evidence?
2. Does the operational evidence trail tie one-to-one to the claim, in a schema my staff can ingest without bespoke parsing?
3. Does the corpus support the user-entity side as cleanly as the vendor side, so the report is consumable without ambiguity?

I am explicitly not reading for cryptographic novelty. I am reading for whether the next ten engagement teams that pick this up can run it the same way.

### Questions

**Q1. Is the updated audit-procedure P-6 (key-fingerprint reconciliation) testable end-to-end against the operational evidence the implementation actually emits, or does it leak unstated dependencies onto the audit team?**

**Status:** Answered

P-6 in `docs/audit-procedures.md` is one of the cleanest reconciliation procedures I have seen in this control class. It reads:

- Pull `master.reconciliation_completed` events
- Confirm cadence
- Confirm `fingerprint_unmatched_count` is zero or within baseline
- For each `(tenant_id, key_version)` pair observed, recompute `SHA-256(utf8(tenant_id) || ikm)[:16]` against the IKM-roster's expected IKM bytes
- For non-zero unmatched, walk the institution's investigation

Every input is named. The reconciliation event in `docs/soc-pack/control-evidence-events.md` carries `key_versions_observed`, `key_fingerprints_observed`, `fingerprint_unmatched_count`, and `period` — exactly the columns my staff needs to tie the procedure to evidence. Spec §10.1 anchors the cadence ("no more than weekly") and explicitly hands the evidence trail to P-6. Spec §3 defines `key_fingerprint` so the recomputation step is unambiguous: 16 raw bytes, `SHA-256(utf8(tenant_id) || ikm)[:16]`. There is no formula ambiguity for staff to argue about.

The one operational dependency P-6 inherits — that the institution actually maintains an "IKM roster" with expected fingerprints per `(tenant_id, key_version)` — is covered by CUEC-CRY-03 and CUEC-CRY-05 in `docs/control-map/CUECs.md`. CRY-05 specifically commits the institution to retaining historical public keys with validity windows; the equivalent expectation for IKM-roster maintenance is implicit in CRY-03 and CRY-04 but I would prefer it stated explicitly. That is a one-line tightening for round 10, not a gap.

The chain that closes P-6 is also satisfying for an auditor: spec §7 step 8 catches the same condition at the verifier (no MAC compute on fingerprint mismatch), the operational event records the failure with `step=8` and a parseable `detail` (per the example in `control-evidence-events.md`), and P-6 walks the operational evidence. Three independent observation points for the same condition. That is what defense-in-depth looks like in an attest engagement.

**Q2. Does P-3 (session-key handshake authentication) produce evidence the SOC team can sample-test, or does the per-tenant determinism property silently require the institution to operate test infrastructure my staff cannot reach?**

**Status:** Partial

P-3 is mostly answered. The first three bullets are textbook attest procedure (pull custodian auth logs, confirm authenticated workload identity, confirm anonymous handshakes were rejected). Spec §4.1.1 articulates the four properties cleanly and the recommended mechanisms (SPIFFE/SPIRE, mTLS, HSM-issued tokens) are conventional. CUEC-IAM-05 ties the procedure to a user-entity control my staff can confirm against the institution's IAM evidence.

The fourth bullet — "Confirm the per-tenant determinism property (spec §4.1.1 property 4): a sample handshake replayed for the same tenant produces a byte-identical `key_fingerprint`" — is where my pencil hovers. Two questions the procedure does not answer:

1. **Who runs the replay?** A SOC engagement does not normally execute the IKM custodian's HKDF derivation against production IKM. Two acceptable shapes exist: the institution operates a documented self-test that produces a re-derivable fingerprint and the SOC team inspects the self-test record; or the institution provides a non-production IKM the SOC team can re-derive against. Neither shape is named.
2. **What is "byte-identical" being compared to?** Implicitly, the previously-observed fingerprint stamped on a chain entry. Stating this — "compare the replay's fingerprint against any `key_fingerprint` previously stamped on a chain entry for the same `(tenant_id, key_version)`" — would close the loop for staff.

This is a minor procedural tightening. The spec's property is correct and the operational events carry the fingerprint, so the evidence is there. The procedure just needs one more sentence so a staff auditor knows what to actually do. I would not modify the SOC opinion for this; I would write a procedure-clarification note in the engagement working papers and proceed.

**Q3. Can the SOC team mechanically consume `chain.verification_failure` events with the new `key_fingerprint`, `key_version`, `format_version`, and `step` fields without bespoke parsing per engagement?**

**Status:** Answered

Yes, and this is the strongest piece of the corpus from an attest-engagement standpoint. The operational event schema in `docs/soc-pack/control-evidence-events.md` is now structured so a single mechanical query — filter by `event = "chain.verification_failure"`, group by `step` — produces the evidence sub-table for every spec §7 verification step. The `step` field references the spec §7 procedure step (1–12) explicitly. That tie is the kind of thing that lets a Big Four firm publish an engagement template instead of rebuilding consumption logic per client.

The two worked examples in the doc cover the two failure modes that matter most for attest evidence: `step=9` payload_hash MAC mismatch (the "the chain bytes were tampered" condition) and `step=8` key_fingerprint mismatch (the "the IKM is not what the institution thinks it is" condition). The `detail` field on the `step=8` example carries actionable text — `expected=…; recorded=…; investigate IKM roster row (tenant=…, key_version=1)` — which lands in my staff's evidence dump in a form they can paste straight into a working paper.

The schema discipline is what makes this work. Every event carries `tenant_id`, `correlation_id`, `timestamp`, and a `fields` block whose shape is event-specific but documented. The retention requirement ("at least as long as the chain events they relate to") is in §10.2 of the spec and restated in `control-evidence-events.md`. CUEC-OPS-02 is the user-entity-side commitment to the same retention. I would have asked for exactly this if it weren't already there.

The `master.reconciliation_completed` event is the same shape and consumes the same way: filter by event name, group by tenant, eyeball `fingerprint_unmatched_count` per period. Cross-walked against P-6, the procedure and the evidence are one-to-one.

One small observation, not a gap: `docs/soc-pack/control-evidence-events.md` shows `seal.job_completed` carrying a `master_version` field in the example, but the spec §3 / §4.1 vocabulary uses `key_version` (integer ≥ 1) on chain entries and the seal record carries `key_versions` (list). I read the example's `master_version` as the operational/custodian-side version label (e.g., `"v3"`), distinct from the per-entry `key_version` integer, and I would expect the doc to say so in one sentence. As-is, my staff would ask the question and the answer is benign, but the doc would be tighter if it pre-empted it.

**Q4. Does `docs/regulator-pack/sample-report.md` give my engagement team a defensible reference shape for the verifier output we will sample as evidence under P-13 (internal-audit verifier-run cadence) and as the substantive evidence behind PI1.1/PI1.2 in a SOC 2 opinion?**

**Status:** Partial

The sample report is good for the FFIEC-examiner audience it was written for. The cover page format, the per-day pass/fail detail, the anomaly section, and the methodology section are all there. The fictional 30-day result with three explained anomalies is realistic and gives my staff a template for what "pass with anomalies" looks like in production.

Where it sits at "Partial" for a SOC engagement specifically:

1. **No traceability between verifier output and the specific TSC the SOC opinion attests to.** The methodology section names HMAC-SHA-256, RFC 6962 Merkle, and Ed25519. A SOC team needs the next layer: which verifier output line is the substantive evidence for PI1.1, which line is for PI1.2, which line is for CC7.2. The TSC mapping in `docs/control-map/TSC-mapping.md` says PI1.1/PI1.2 are the headline, but the sample report doesn't surface that mapping in a column the SOC team would copy into a working paper. A two-column "verifier line → TSC criterion" appendix would bridge this for the SOC audience without complicating the examiner audience.

2. **The report doesn't surface the new spec §7 verification steps in a way that lets the SOC team confirm the verifier exercised them.** A reader of this sample report knows the chain walk passed; they cannot immediately see that the verifier executed step 1 (format-version check), step 2 (HKDF inputs digest), step 4 (per-entry binding / cross-chain-lift defence), step 7 (IKM lookup), and step 8 (fingerprint check). For SOC purposes this matters because the absence of a `chain.verification_failure` event for a given step is meaningful only if the verifier exercised the step. A line in the methodology section enumerating "steps executed: 1–12" with a tick per day would close this.

3. **No worked example of a `--strict` invocation.** P-5 commits the institution to no `kms_handle_uri = "plaintext-dev"` in production; spec §10.7 commits the verifier to refusing dev-mode seals under `--strict`. The sample report shows an invocation without `--strict`. Production-engagement reality is that `--strict` is what the SOC team would run for the substantive test. A second invocation block, or a sentence noting that `--strict` is the SOC team's mode, would remove ambiguity.

These are sample-report polish items, not corpus gaps. The spec, the operational events, and the audit procedures are sufficient for me to walk a SOC engagement without them. The sample report would just be a more useful starting point for my engagement template.

**Q5. Does the anomaly-documentation template (`docs/anomaly-documentation-template.md`) accommodate the new failure modes the rework introduced — particularly `key_fingerprint` mismatch and `unknown_key_version` — without forcing the institution to invent a category?**

**Status:** Gap

This is the one place the rework hasn't fully landed in the corpus.

The template's "Anomaly type" enum lists: `Sealing delay | Late-binding rate elevated | Clock-skew | Missing seal continuity | Software-key in production | Other`. The severity mapping table at the bottom of the doc covers the same set plus `Sealing delay associated with HSM cluster outage` and `Software-key (dev-mode) in production`.

The new failure modes the rework introduced are not in either list:

1. **`key_fingerprint` mismatch (spec §7 step 8)** — this is the loud detection of botched rotation, restored-from-wrong-backup, or cross-tenant key swap. It is not a routine anomaly; it is the headline tampering-or-misconfiguration signal the new construction is designed to surface. It needs its own row in the severity table (default: Critical, with the IR playbook reference) and its own enum value in the template.
2. **`unknown_key_version` (spec §7 step 7)** — verifier was handed a chain entry whose `key_version` has no corresponding IKM in the registry. Default severity: High at minimum (potential rotation procedure failure, potential unauthorised IKM in use). Needs the same treatment.
3. **`format_version` mismatch (spec §7 step 5)** — should be impossible inside a single file because §7 step 1 refuses unrecognized header versions, but the per-entry check at step 5 catches a mid-file format drift. Probably Medium default. Needs an enum entry.
4. **`hkdf_inputs_digest` mismatch (spec §7 step 2)** — header-pre-flight failure. Critical. Needs an enum entry.

Right now if a `key_fingerprint` mismatch surfaces during a reporting period, the institution would file the record under "Other" with a free-text description. A SOC team reviewing that record cannot mechanically distinguish "Other = network glitch the institution chose not to enumerate" from "Other = key_fingerprint mismatch which is a Critical-severity event with an IR playbook activation." The whole point of the template is to make severity assessment mechanical and the SOC team's evaluation reproducible.

This is a documentation update, not a spec or implementation gap — the spec and the operational events handle the cases correctly. But until the template is updated, my engagement team would have to issue a procedural recommendation that the institution adopt a local extension to the template covering the four modes above. That recommendation would land as a written observation in the SOC report, not a finding, but it is avoidable with one round of editing on the template.

I would also extend the template's "What 'operationally explained' means" section. The current four bullets are fine for sealing delays and clock-skew. For a `key_fingerprint` mismatch, "operationally explained" carries a different bar — at minimum a documented IKM-roster correction with change-management approval, and a re-run of the affected period under the corrected roster. The template should call this out so the institution and the SOC team agree on what the bar is before they're standing in front of an exception.

---

## Per-role roll-up

| Status | Count |
|---|---|
| Answered | 2 |
| Partial | 2 |
| Gap | 1 |
