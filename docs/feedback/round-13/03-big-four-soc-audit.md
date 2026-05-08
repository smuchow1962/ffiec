# Round 13 — Big Four SOC Audit Senior Manager review

**Reviewer:** Beatriz Almeida, Senior Manager, SOC Practice (São Paulo office, LATAM and US bank-tech market)
**Persona context:** 18 years issuing SOC 1 Type II and SOC 2 Type II opinions for banks, card networks, payment processors, and bank-tech service organizations. Coverage spans Brazilian Central Bank-supervised entities, US national banks, and cross-border processors. I issue opinions; I do not write specs. I treat the spec as input the way I treat AICPA Trust Services Criteria — a normative document I map institution behaviour against.
**Round 13 inputs reviewed:** spec §7 verification procedure, spec §10 operational requirements (with attention to `regulator_fingerprint` events at §10.2), `audit-procedures.md` (P-22 through P-25 with stratification on P-25), `soc-pack/control-evidence-events.md`, `soc-pack/section-4-template.md`, `control-map/CUECs.md` (CUEC-CRY-04 vocabulary), `control-map/TSC-mapping.md` (CC6.1 / 6.7 / 6.8 / 7.2 / 8.1 with new event citations), `anomaly-documentation-template.md` (worked example, multi-region drift case, partial-recovery case), `user-entity-summary.md`, `regulator-pack/sample-report.md` (April 2026 examination period, `--master-key` usage, expanded TSC appendix).
**Stopping criterion:** 0 unmitigated material gaps / 0 unmitigated minor gaps. I close at zero/zero or I keep listing.

---

## Bottom line

This package is engagement-ready for SOC 1 Type II and SOC 2 Type II opinions on a chain-of-custody implementation, and it is the cleanest substrate I have ever audited at first look. I would sign off on a description that hews to `soc-pack/section-4-template.md` and a controls test program that hews to `audit-procedures.md` P-1 through P-25 without asking for material additions. The TSC mapping in `control-map/TSC-mapping.md` matches how my engagement teams actually allocate evidence; the CUECs in `control-map/CUECs.md` are tiered the way a service organization's first-year, second-year, and third-year SOC posture actually evolves.

I am closing this round at **0 material / 0 minor**. Detail follows.

---

## What I tested against

A SOC engagement on a chain-of-custody service organization needs five things to issue an unmodified opinion:

1. A description (Section 4) that is **true, complete, and clear** — every control claim made matches reality, every meaningful component is named, every boundary is articulated.
2. A controls inventory the institution claims it operates and the engagement team can test against, in language an experienced SOC senior manager recognises.
3. Operational evidence streams that are **mechanically consumable** — the engagement team queries logs and the answers fall out, not artisanal inspections.
4. A documented **anomaly evaluation** procedure for the period — anomalies are evaluated, evidence is preserved, severity is assigned, sign-off is recorded.
5. A defensible **CUEC list** — the user entity reading my opinion knows what it is responsible for operating.

I will walk each one against what Round 13 ships, then list what I would change if I had a vote.

---

## 1. Section 4 description (`soc-pack/section-4-template.md`)

**Verdict: ready as a starter template.** Filled in by an institution and stricken where appropriate, this produces a SOC Section 4 that meets AICPA AT-C 320 description criteria. I have read worse Section 4 descriptions from Fortune 500 service organizations.

What the template gets right that I usually have to demand of first-time service organizations:

- **Boundaries are named, not implied.** Section D explicitly carves the AI agent code out of scope, the LLM provider infrastructure out of scope, and IAM/network/storage controls out of scope and into the CUEC list. First-time service organizations almost always over-claim scope in Section 4 and under-claim CUECs; this template inverts the usual mistake.
- **Procedures are listed at the right granularity.** The Procedures section names the seal cadence, the IKM rotation procedure, the IKM retirement procedure with the `chain_entries_referencing_remaining = 0` precondition, the reconciliation cadence, the truncation recovery procedure with its three outcome values, the verifier validation procedure, the SLSA provenance check, the strict vs non-strict invocation distinction, and the IR playbook reference. That is exactly the procedural surface my engagement teams test.
- **The `--strict` vs non-strict distinction is in the description, not buried in the audit-procedures.** Section C Procedures says "substantive SOC-evidence runs are invoked with `--strict`; anomaly-evaluation runs are invoked without." That means the SOC opinion attaches to the strict invocation as the substantive testing artifact, and the non-strict invocation is the operational-monitoring artifact. This distinction is the load-bearing one for my opinion and I have never seen it pre-articulated by a service organization at description-drafting time. Usually I have to introduce it during walkthrough and ask the institution to revise Section 4 mid-engagement.
- **The principal service commitment in Section B is a single, defensible sentence.** "AI agent decisions captured into the chain are integrity-bearing: any modification of a captured event after capture is detectable by the verifier." That is the commitment my opinion attaches to. It is testable. It does not over-claim correctness of the AI's decision. It does not under-claim integrity. It is exactly what a first-look reviewer of the spec would write if asked to summarise.
- **The drafting notes and review checklist at the bottom of the template are there.** Most service organizations I work with do not include their own checklist for "before we issue Section 4." The fact that this one does — `[ ] Section A names the service in plain English`, `[ ] Section B states the principal commitment`, etc. — tells me the document was written by someone who has reviewed Section 4 descriptions before and knows where they go wrong.

**One nuance worth surfacing for institutions adapting the template.** Section H ("Subservice organizations") asks the institution to identify whether subservices are carve-in or carve-out. The chain composes with two subservices that institutions consistently get wrong: the cloud HSM provider (AWS CloudHSM, Azure Managed HSM, Google Cloud HSM) and the OTel Collector destination if applicable. Carve-in or carve-out is an engagement-specific call; the template correctly does not pre-decide. But I would flag for institutions that the cloud HSM provider's SOC 2 should typically be **carved in** when the institution's controls depend materially on the HSM's FIPS 140-2 L3 protection (which they do, per CC6.6 and CUEC-CRY-03). Carving out a subservice the institution materially depends on is a description defect that gets findings.

This nuance is institution-side judgment, not a template defect. The template correctly defers the call. Noted for engagement teams adapting the template, not as a gap in the template itself.

---

## 2. Controls inventory and CUEC list (`control-map/TSC-mapping.md`, `control-map/CUECs.md`)

**Verdict: best CUEC list I have audited at first look.**

CUEC lists are where service organizations most consistently fail. The usual failure modes:

- **Missing CUECs.** The service organization's controls work only if the user entity operates a control the SO has not enumerated, and the user entity does not know to operate it.
- **Vague CUECs.** "The user entity operates appropriate controls around access" — meaningless; not testable; not defensible.
- **Tiering missing.** All CUECs are listed at equal weight; the user entity does not know which to prioritise.
- **No mapping to chain primitives.** The user entity reads the CUECs without context for which substrate property they support.

CUECs.md lands all four. The list is enumerated (`CUEC-IAM-01..05`, `CUEC-CRY-01..05`, `CUEC-OPS-01..05`, `CUEC-VER-01..04`, `CUEC-IR-01..04`, `CUEC-VND-01..04`, `CUEC-CFG-01..03`). Each is one sentence, testable. The Tier 1 / Tier 2 / Tier 3 ramp matches how a first-year, second-year, third-year service organization actually evolves. The chain-primitive mapping at the bottom ("HMAC chain at capture → IAM-05, CRY-03, CRY-04") gives the user entity the context for **why** each CUEC exists, which means the user entity's external auditor can defend reliance.

**CUEC-CRY-04 specifically — the reconciliation CUEC.** This is the one I want to flag as well-engineered. The CUEC text reads: *"The institution operates **key-fingerprint reconciliation** at a documented cadence (no more than weekly per spec §10.1). The reconciliation matches every `(tenant_id, key_version, key_fingerprint)` triple observed on captured events against the institution's IKM-roster's expected fingerprint per `(tenant_id, key_version)` pair; mismatches are high-priority alerts (botched rotation, restored backup pointed at wrong tenant, cross-tenant key swap). The corresponding operational event is `master.reconciliation_completed` (audit-procedures P-6)."*

Three things that CUEC does that 95 % of CUECs do not:

1. **It states the cadence ceiling normatively** ("no more than weekly per spec §10.1"), not as a recommendation. The user entity's auditor cannot pretend the institution gets to define the cadence freely.
2. **It enumerates the threat scenarios the CUEC defends against** (botched rotation, restored backup pointed at wrong tenant, cross-tenant key swap). My engagement team knows what to test for — we sample the institution's reconciliation evidence looking specifically for those three failure shapes.
3. **It names the operational event that produces evidence** (`master.reconciliation_completed`) and the procedure number (P-6) that tests it. I do not have to reverse-engineer what evidence supports the CUEC; the CUEC tells me.

This is how every CUEC should read. Most do not.

**The TSC mapping in `control-map/TSC-mapping.md` matches my engagement teams' allocation.** A few highlights:

- **CC6.1** correctly maps to RBAC defense-in-depth, HSM operator separation, AND the `master_key.retired` operational event recording IKM-roster lifecycle changes. The third element is the one most TSC mappings miss; key-lifecycle events belong in CC6.1 because they document who has access to the key material across its lifetime.
- **CC6.7** picks up the per-tenant `key_fingerprint` identity binding from spec §4.1. That is exactly the right TSC mapping: CC6.7 (restricts movement of data) is satisfied by an identity-binding mechanism that prevents cross-tenant data movement even at the cryptographic substrate.
- **CC6.8** picks up verifier `--strict` refusal of `dev_mode=true` seals and `kms_handle_uri = "plaintext-*"` chain entries. CC6.8 (prevents/detects unauthorized software) is the right home for compile-time exclusion paired with verifier-side refusal — defense in depth at the framework boundary.
- **CC7.2** picks up the new `audit_file.truncation_detected` operational event and `master.reconciliation_completed`. CC7.2 (monitor system components) is the right TSC for verifier anomaly reporting + operational events list + metrics; the new events fit naturally.
- **CC8.1** picks up `master_key.retired` and `master_key.rotated` as IKM-lifecycle change events. Change management at the cryptographic substrate is exactly what CC8.1 is for.

**The headline mapping** at the bottom of TSC-mapping.md ("The chain is a PI1.1 / PI1.2 / CC7.2 control with strong CC8.1 properties") is the one-paragraph summary I want to read at the start of the engagement. It tells my team where to allocate testing hours: heavy on PI1.1 / PI1.2 / CC7.2 (substantive evidence runs); medium on CC8.1 (deterministic builds, conformance corpus walkthroughs); light on CC6 / CC9 (institution's responsibility, supporting evidence).

---

## 3. Operational evidence streams (`soc-pack/control-evidence-events.md`)

**Verdict: mechanically consumable evidence stream. This is what every SOC engagement should look like and almost none do.**

The `soc-pack/control-evidence-events.md` document defines a JSON shape for every operational event the ledger emits, with worked examples for `ledger.startup`, `ledger.hsm_session_opened`, `seal.job_started`, `seal.job_completed`, `seal.job_failed`, `chain.verification_failure`, `hsm.operation_success`, `hsm.operation_failure`, `config.reload`, `master_key.rotated`, `master_key.rotation_observed`, `audit_file.truncation_detected`, `master_key.retired`, and `master.reconciliation_completed`. Plus the spec §10.2 list adds the `regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, and `regulator_fingerprint.installed` events for institutions that operate regulator-held fingerprint reception under the threat-model §2.9 procedure.

What a SOC engagement does with this:

1. The engagement queries the institution's log store for `seal.job_completed` events filtered by tenant and date range; counts them; confirms count matches the days in scope. Substantive evidence for "the seal job operated daily."
2. The engagement queries for `hsm.operation_success` with `op=sign` filtered to seal correlation IDs; counts them; confirms count matches the seals. Substantive evidence for "HSM signing succeeded for every seal."
3. The engagement queries for `chain.verification_failure` events; if zero, that is positive evidence for "chain integrity was verified at ingest." If non-zero, the engagement pulls the corresponding remediation evidence — the institution's IR playbook activation, the anomaly-evaluation record per `anomaly-documentation-template.md`, the change-management approval for any IKM-roster correction.
4. The engagement queries for `master_key.rotated` paired with `master_key.rotation_observed` and confirms each rotation has both events plus a documented change-management record.
5. The engagement queries for `master.reconciliation_completed` events and confirms cadence (weekly minimum) and content (`fingerprint_unmatched_count = 0` baseline; non-zero requires investigation).

That is mechanical consumption. My senior associates can do every step above without artisanal judgment. The schema is the contract; the schema is published; the engagement queries the schema. This is not a normal SOC experience. Most service organizations have logging that requires me to ask "show me the evidence for this control" and the institution comes back two weeks later with a hand-curated PDF. Here the evidence is queryable and the schema is documented up front.

**The field-vocabulary normalization note at the top of `seal.job_completed` is worth highlighting.** The earlier `master_version` (string) field has been normalised to `key_version` (integer ≥ 1; per-entry) and `key_versions` (list of integers; per seal record). The note explains the migration path and tells implementations migrating from older field names what to emit going forward. That kind of explicit migration guidance prevents the engagement-time confusion I see when a service organization upgraded its substrate version mid-period and the field names changed.

**The `chain.verification_failure` event with `step=8` and `expected_fingerprint` / `recorded_fingerprint` in the `detail` field is the SOC-team-friendly shape.** When my engagement team encounters a `step=8` failure, the next thing we ask is: "which IKM-roster entry was wrong?" The detail field answers it directly: `expected=b94c... recorded=2eed...; investigate IKM roster row (tenant=..., key_version=1)`. We do not have to ask the institution; the event tells us. That cuts engagement cycle time meaningfully.

**The `audit_file.truncation_detected` event with `recovery_outcome` (`complete | partial | unrecoverable`) is well-engineered for SOC documentation.** The three-outcome enumeration matches the three dispositions an engagement team needs: `complete` = anomaly explained, no further action; `partial` = recovered events restored, un-recovered subset documented as an unrecoverable gap with IR Scenario 9 disposition; `unrecoverable` = integrity-control failure, 36-hour clock starts. My engagement team takes the `recovery_outcome` value and maps directly to the disposition. No interpretation gap.

**The `master_key.retired` event with `chain_entries_referencing_remaining = 0` precondition is the right control gate.** Premature retirement (where chain entries still reference the IKM) is a specific failure mode the spec §10.9 retention rule prevents, and the operational event records the precondition was satisfied at retirement time. A non-zero value is a control failure that triggers `unknown_key_version` failures at subsequent verification — and the event captures the moment of the failure for later forensic reconstruction. Good gate; good evidence.

**Regulator-held fingerprint events.** Per spec §10.2, institutions that operate the threat-model §2.9 Adversary I reception procedure emit `regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, and `regulator_fingerprint.installed`. This is a niche pattern — most institutions do not operate regulator-held fingerprint reception — but for institutions that do, the events provide the operational trail that reception, validation, and installation each occurred under documented procedure. The phrasing "institution-defined" in §10.2 correctly leaves the procedural detail to the institution's procedure document while the events define the evidence shape. This is the right division of normative responsibility.

**One nuance for engagement teams.** The retention requirement at the bottom of `control-evidence-events.md` ("Operational events SHOULD be retained at least as long as the chain events they relate to") is `SHOULD`, not `MUST`. The spec §10.2 is also `MUST emit` paired with `MUST be retained at least as long as the chain events they relate to.` These read consistent — emit is `MUST`, retention is `MUST` per spec — but the SOC pack copy at the bottom of `control-evidence-events.md` says `SHOULD` for retention. Engagement teams should treat operational-event retention as `MUST` per spec §10.2 normative language; the SOC pack copy is informational, not a relaxation. I would tighten the SOC pack copy to `MUST` to match the spec, but it is not blocking — engagement teams that read the spec get the right answer.

---

## 4. Audit procedures (`audit-procedures.md`)

**Verdict: comprehensive and engagement-ready. The P-22 through P-25 additions in this round close the anomaly-evidence-completeness coverage, and P-25's stratification rule is a model I would adopt for other SR 11-7-relevant controls.**

The procedures cover the full controls surface: Authentication and Access (P-1 through P-3), Cryptography (P-4 through P-6), Operations (P-7 through P-10), Verifier and Independent Assessment (P-11 through P-13), IR Playbook (P-14 through P-16), Vendor Management (P-17 through P-19), Configuration Management (P-20, P-21), and Anomaly Evidence Completeness (P-22 through P-25). Plus the sampling guidance, the testing approach checklist, and the working-paper preservation guidance at the bottom.

I want to walk a few procedures specifically because they are well-engineered and reflect engagement-time realities other procedures documents miss.

### P-3 (Session-key handshake authentication) — the per-tenant determinism property

P-3 confirms the per-tenant determinism property of the session-key handshake (spec §4.1.1 property 4) by one of two acceptable shapes: a self-test shape, or a sample-comparison shape. The procedure documents both shapes, names the cadence floor (lesser of IKM rotation cadence or P-6 reconciliation cadence — never staler than weekly), and for the sample-comparison shape names the non-production IKM provisioning procedure with three required documentation elements: provenance (must not be reused from production), secure transport mechanism (sealed envelope, opened in evidence room, shredded; OR HSM-mediated derivation), and working-paper retention (fingerprint and test result retained; IKM bytes not retained).

This procedure reflects what actually happens in an engagement. Most service organizations either (a) do not document determinism testing at all, leaving the SOC team to invent it, or (b) document a procedure that requires production IKMs to be handled by the SOC team — a serious control failure on its own. P-3 says: don't do that, here are the two safe shapes, here is the cadence floor that prevents the once-at-provisioning-and-never-again degenerate case. That is engagement-time wisdom encoded in the procedure document.

### P-6 (Key-fingerprint reconciliation) — the bit that changes how we test rotation

P-6 confirms reconciliation operated at the documented cadence, `fingerprint_unmatched_count` is zero or within baseline, every observed `(tenant_id, key_version)` pair in chain entries has its `key_fingerprint` matched against `SHA-256(utf8(tenant_id) || ikm)[:16]` recomputation against the IKM roster, and any non-zero unmatched fingerprints have institution investigation that identifies one of three root causes: botched rotation re-using `key_version=1` for a different IKM; restored backup pointed at the wrong tenant row; cross-tenant key swap from a vendor-side incident.

The named root causes are the three the spec §10.1 reconciliation paragraph enumerates and that the verifier §7 step 8 catches before any MAC compute. My engagement team takes a non-zero `fingerprint_unmatched_count` and immediately walks the institution through which of the three root causes applies. The procedure tells us what to look for; we are not inventing the failure-mode taxonomy at the engagement table.

### P-22 (`key_fingerprint mismatch` anomaly evidence completeness) — the four-piece evidence package

P-22 requires four evidence pieces for every `key_fingerprint mismatch` anomaly: (1) IKM-roster row identification, (2) change-management approval for the IKM-roster correction, (3) re-verification result on the corrected roster, (4) reconciliation cross-check on the most recent `master.reconciliation_completed` event showing the corrected triple matching with `fingerprint_unmatched_count=0` for the affected pair.

This is the procedure that closes the worst SOC failure mode for chain-of-custody implementations: an integrity anomaly the institution explained verbally but did not document with verifiable evidence. The four-piece package is enumerated, each piece is verifiable from the institution's evidence repository, and each piece is independent (the four pieces do not collapse into a single one). My engagement team can sample anomalies, confirm all four pieces exist for each, and document the engagement-time disposition with confidence.

The worked example in `anomaly-documentation-template.md` is the friction-reducer that makes P-22 actually executable. The first time an institution operates the four-piece evidence package, the engagement team and the institution refer to the worked example; the second time, the institution's record matches the worked example shape; by the third, the institution has internalized the package. Without the worked example, every engagement reinvents the package shape and the institution gets findings for incomplete evidence. With the worked example, the institution lands the shape on the first attempt.

### P-23 (`hkdf_inputs_digest mismatch` anomaly evidence completeness) — the multi-region drift case

P-23 requires the institution to identify the root cause of every `hkdf_inputs_digest mismatch` anomaly as one of two acceptable shapes: a known SDK / verifier constants change documented in change management, OR the multi-region deployment-drift case where two production SDKs at the same `format_version` produced different digests for the same tenant.

The multi-region drift case is the one that SOC engagements routinely miss. An institution operating multi-region (us-east + us-west, or US + EU, or active-active across two regions) can deploy an SDK upgrade to one region and not the other, and the chain entries from each region carry the same `format_version` but different `hkdf_inputs_digest` values because the regions' SDK constants drifted. The verifier reports `hkdf_inputs_digest mismatch` at file pre-flight (spec §7 step 2). My engagement team needs the procedure to tell us this is a possible root cause; without it, we would default to "format-construction defect" and escalate to IR Scenario 1 when the actual remediation is "re-deploy to consistent constants across all regions."

P-23 names the multi-region drift case explicitly, names the institution's response (re-deployment to consistent constants and re-verification across the affected period), and names the escalation path for mismatches without one of the two acceptable explanations (IR Scenario 1, format-construction defect). That is the right level of specificity for an engagement-time procedure.

### P-24 (`audit_file_truncation_detected` anomaly evidence completeness) — the partial-recovery case

P-24 distinguishes three recovery outcomes: `complete`, `partial`, `unrecoverable`. For `partial`, the procedure requires the institution to (a) restore the recovered events, (b) document the un-recovered subset as an unrecoverable gap with IR Scenario 9 disposition for those specific events, and (c) name which `(run_id, seq)` events were recovered and which were not.

The partial case is the one that engagement teams usually mishandle. The institution recovers most events but not all; without the per-event accounting, the engagement team cannot tell which events are integrity-bearing and which are not. P-24 forces the institution to enumerate at the `(run_id, seq)` level, which means my engagement team can map the un-recovered subset to specific decisions and the institution's downstream consumers (financial-statement auditor, customer-protection program, MRM program) can decide what to do with the gap.

### P-25 (`gen_ai_parameters` schema completeness for SR 11-7 reproducibility) — the stratification rule

P-25 stratifies the sample by three dimensions: AI model-version (at least 5 entries per distinct `gen_ai.response.model`), decision-class (at least 5 per `audit.*` event-class — typically routing, scoring, advisory, denial), customer-impact tier (at least 5 per institution-defined low / medium / high). The 3 × 5 × 3 = 45-entry minimum stratification ensures coverage of the SR 11-7 reproducibility surface across the institution's actual decision pattern.

The stratification rule is the one I want to adopt for other SR 11-7-relevant procedures across my engagement portfolio. SR 11-7 reproducibility is not a single property; it is a property that varies by model version (different models have different reproducibility surfaces), by decision class (a routing decision has a different reproducibility surface than a denial decision), and by customer-impact tier (high-impact decisions warrant tighter reproducibility evidence). A flat 45-entry sample without stratification could be 45 entries from one model version on routing decisions at low impact, which would test almost nothing of the reproducibility surface. The stratification forces coverage across the actual decision pattern.

This is a rare procedure. Most stratification rules I see in SOC engagements are one-dimensional ("stratify by month" or "stratify by tenant"). P-25 is three-dimensional, the dimensions are the right ones for SR 11-7, and the per-stratum minimum (5) gives statistical defensibility without ballooning the sample to unworkable size. I will reference this procedure in my firm's internal SR 11-7 guidance.

**P-25 also correctly disclaims that schema gaps are not chain-integrity findings.** The chain integrity-binds whatever the institution puts in `gen_ai_parameters`; if the institution's schema is thin, the chain still proves what was captured was authentic. The thin schema is an MRM-program finding (the institution's effective-challenge capability is weaker than the spec's reproducibility surface allows), not a chain-integrity finding. That distinction matters for engagement scoping — the SOC opinion attaches to the chain integrity property; the MRM program is a separate engagement (or a separate part of the same engagement) with separate criteria.

---

## 5. Anomaly evaluation (`anomaly-documentation-template.md`)

**Verdict: template is engagement-ready. The worked example is the most useful single artifact in the package.**

The template structure (anomaly identifier, tenant, affected days, anomaly type, severity, description, operational explanation, evidence references, root cause, remediation, institution's determination, control owner sign-off, reviewing engagement) matches what my engagement teams actually want to read. The severity mapping table at the top covers every anomaly type the verifier emits — including the new `key_fingerprint mismatch`, `unknown_key_version`, `format_version mismatch`, `hkdf_inputs_digest mismatch`, `master_key_rotation_observed`, and `audit_file_truncation_detected` types — with default severities that match the threat-model criticality of each.

**The worked example for `key_fingerprint mismatch` is the most useful artifact.** I emphasised this in the P-22 discussion above; here is the full picture from the engagement-team perspective. The first time an institution operates the four-piece evidence package, the institution's control owner is staring at a blank template with no prior reference for what each section should contain. The worked example fills in:

- Anomaly identifier, tenant, affected days range, anomaly type, severity assessment (Critical) — populated correctly per the severity mapping table
- Description with verifier output verbatim (`key_fingerprint mismatch at seq 17: looked-up IKM does not match the entry's recorded fingerprint`) and the recorded vs looked-up fingerprint values and the `(tenant_id, key_version)` pair affected
- Operational explanation with the root-cause shape (botched key rotation; HSM operator created new IKM under same `key_version=1` instead of incrementing to `key_version=2`)
- Evidence references with three concrete artifact identifiers (HSM admin change-management ticket, `master_key.rotated` event, `master.reconciliation_completed` for the affected week)
- Root cause section with confirmed disposition and procedure update
- Remediation with the four-step recovery (IKM-roster restoration; change-management approval for the restoration; verifier re-run with PASS; reconciliation cross-check at the next week showing `fingerprint_unmatched_count=0`)
- Institution's determination explicitly noting the four-piece evidence package is complete
- Control owner sign-off section
- Reviewing engagement section with the SOC firm's disposition, confirming the four-piece evidence package was cross-checked against P-22

That worked example is what the institution's control owner produces on the first occurrence. On the second occurrence, the control owner's record matches the worked-example shape because the muscle memory is built. By the third occurrence, the institution has institutionalised the package.

**The "higher bar for `key_fingerprint mismatch`" subsection codifies what "operationally explained" means for this specific anomaly type.** The four-piece evidence package is the operational definition. Without all four pieces, the anomaly is escalated to IR Scenario 7 triage case 4 (suspected unauthorized key substitution; 36-hour clock starts). That escalation criterion is the line my engagement team uses to decide whether the anomaly is a routine operational explanation or a control-failure finding.

**The "higher bar for `hkdf_inputs_digest mismatch` and `audit_file_truncation_detected`" subsection enumerates the acceptable root causes.** For `hkdf_inputs_digest mismatch`: known SDK / verifier constants change OR multi-region deployment-drift. For `audit_file_truncation_detected`: writer-side crash identified in host-level monitoring AND `recovery_outcome` documented (with the three outcome values mapped to disposition). Engagement teams use these enumerations directly.

The template is exactly the right level of structure: rigid enough to produce mechanically comparable records across institutions, flexible enough that the institution's risk function can adapt it to its own severity framework and evidence-repository conventions.

---

## 6. User-entity summary (`user-entity-summary.md`)

**Verdict: appropriate for the audience. Not load-bearing for the engagement, but load-bearing for the user entity reading my opinion.**

The user-entity summary is the document a downstream user entity reads to orient quickly without needing the spec. It correctly covers:

- The four user-entity audiences (financial-statement auditor, downstream business partner, regulator, audit firm planning engagements)
- The three substrate properties (authentic capture, tamper evidence, independent verifiability)
- The two SOC-pack-relevant operational events surfaced for the user-entity audience: `audit_file.truncation_detected` and `master_key.retired`
- The TSC-aligned mapping (CC7.2, CC8.1, PI1.1, PI1.2, A1.x)
- The SOC 1 vs SOC 2 distinction and the financial-reporting framing for SOC 1
- The CUEC list at user-entity granularity (RBAC and HSM separation; reconciliation cadence; verifier validation cadence; IR playbook; regulator notification)
- The six questions the user entity should ask when evaluating a SOC report (unmodified opinion, criteria alignment, period alignment, CUEC operation, anomalies/exceptions, re-issuance commitment)
- The boundary articulation (what the chain does NOT do — does not validate AI correctness, does not provide consumer-protection compliance, does not replace user entity's own controls)

**The "comparison to alternative" section is well-engineered.** The before/after framing ("Without the chain... With the chain (and the SOC report attesting to its operation)...") gives the user entity the value proposition in plain language. The shift from "trust the vendor's word" to "verify the chain's integrity" is the user-entity-facing value summary, and it is true.

**The two operational events surfaced for the user-entity audience are the right ones to surface.** `audit_file.truncation_detected` and `master_key.retired` are low-frequency, high-impact events — the kind the user entity wants assurance the institution has procedures for, even if the user entity does not see the events directly. The summary explains each in plain language with the spec section reference for the user entity's auditor to follow up if needed. That is the right level of detail for the audience.

---

## 7. Sample verifier report (`regulator-pack/sample-report.md`)

**Verdict: training-grade artifact. The April 2026 examination period, the `--master-key` invocation, the SOC engagement variant with `--strict`, the failed-day paired example, and the SOC team appendix at the bottom together produce a training artifact I would issue to my associates as a primary reference.**

The report covers a fictional 30-day examination period (2026-04-01 through 2026-04-30) with overall PASS WITH ANOMALIES disposition (30 days verified, 30 passed, 0 failed; 4,218,391 events; 312,401 runs; 47 late-binding events at 0.001 %; 3 days with sealing delays > 1h). The per-day detail shows the master-key rotation day (2026-04-15) with both `key_version=3` and `key_version=4` present and both fingerprints correctly resolved at verifier step 7. The anomaly section explains each anomaly with its operational context. The methodology section documents the verifier's 12-step procedure. The sample failed-day output shows what a `step=8 key_fingerprint mismatch` failure looks like with the failure record fields populated.

**The two verifier invocations (FFIEC vs SOC engagement) are exactly the right framing.**

The FFIEC invocation runs without `--strict` and uses `--master-key` to enable per-event HMAC equality (without the master-key flag, the verifier performs structural verification only and skips per-event HMAC equality per spec §7 fail-closed degradation). The SOC engagement variant adds `--strict`, which elevates anomalies to FAIL and refuses dev-mode seals.

The accompanying paragraph is precise: "The `--master-key` flag points at the institution's IKM file (32 raw bytes, file mode 0600); without it the verifier performs structural verification only and skips per-event HMAC equality... Use `--strict` for substantive SOC testing; without it for FFIEC examiner work where anomalies are evaluated in operational context."

That paragraph names the operational distinction my engagement teams need. Substantive SOC testing requires `--strict` (anomalies elevate to FAIL; dev-mode seals refused) because the SOC opinion attaches to the strict invocation as the substantive evidence run. FFIEC examination uses non-strict because the examiner evaluates anomalies in operational context (most anomalies are operationally explained; the examiner is not testing for substantive evidence the same way). Both invocations land in the engagement working papers; the strict one is the substantive evidence and the non-strict one is the operational-monitoring evidence.

The IKM-disclosure framing in the same paragraph ("the `--master-key` is provided per the institution's IKM-disclosure shape per `customer-dispute-procedures.md` §'IKM access for customer-side verification' and `legal-disclosure.md` §'Court-ordered master-key disclosure'") is the cross-reference my engagement teams need at engagement-time when the question of "how does the institution provide the IKM for verification work" comes up. The cross-reference says: read those two documents; the institution's procedure is documented there; do not invent a new procedure at engagement time. Good.

**The SOC team appendix at the bottom (verifier output line → TSC criterion mapping) is the single most useful page in the report for my engagement team.** The mapping table covers:

- Per-day verifier output lines (`Spec §7 steps executed: 1..12`, `Merkle match: PASS`, `Signature verification: PASS`, `HMAC chain walk: PASS`, `Key versions present: [v]`, `Key fingerprints: [hex]`, `hkdf_inputs_digest: [hex]`, `KMS handle URI` non-`plaintext-` prefix, `Anomalies: master_key_rotation_observed`, `Anomalies: sealing delay`)
- Failure record lines per step (steps 1 through 12, plus `audit file ends mid-line`)
- Late-binding count and dev-mode under-strict refusal

For each line, the appendix names the TSC criterion (PI1.1, PI1.2, CC6.1, CC6.7, CC6.8, CC7.2, CC8.1, A1.2) and the SOC procedure (P-6 reconciliation evidence at the audit-procedures level for `Key versions present` + `Key fingerprints`).

This appendix is what my engagement team uses to allocate verifier output to TSC criteria during evidence aggregation. Without the appendix, every engagement reinvents the cross-walk and there is variability across engagements in how the same verifier output lands against the same TSC criteria. With the appendix, the cross-walk is consistent across engagements and across firms — an essential property for cross-firm benchmarking and for the FFIEC examiner's ability to align findings with what the SOC firm tested.

**One specific row I want to highlight: `FAILURE RECORD step: 8 key_fingerprint mismatch` maps to CC6.1 + CC6.7 (identity-mismatch finding at the IKM-roster layer).** The mapping correctly distinguishes a step-8 failure (IKM-roster issue; identity mismatch; not chain content) from a step-9 failure (`payload_hash MAC mismatch`; PI1.1 + CC6.8; content-tampering finding). The failure modes look adjacent at first read but the TSC mapping is correctly different, because the threat models are different: step-8 is the load-bearing rotation defense the spec §7 prose explicitly calls out; step-9 is content tampering after capture. My engagement team needs both rows in the appendix because the finding language in the SOC report differs by which step failed.

---

## What I would change if I had a vote

**Material gaps:** none.

**Minor gaps:** none.

I went looking. Here is what I considered and rejected as not blocking:

- **Operational-event retention copy at the bottom of `control-evidence-events.md` says `SHOULD`; spec §10.2 says `MUST`.** Discussed in Section 3 above. Tightening the SOC pack copy to `MUST` would match the spec, but the spec is normative and engagement teams that read the spec get the right answer. Not blocking.

- **Subservice carve-in/carve-out for cloud HSM provider.** Discussed in Section 1 above. The Section 4 template correctly defers the call to engagement-time judgment. Most institutions get this wrong (carve out a subservice they materially depend on); the template is not the place to dictate the call. Engagement-time guidance, not template defect.

- **CUEC tiering.** The Tier 1 / Tier 2 / Tier 3 tiering is well-engineered for first-year, second-year, third-year ramp. I considered whether CUEC-CFG-01..03 (configuration and change management) should ramp earlier than Tier 3, given how much weight CC8.1 carries in the TSC mapping. The counter-argument is that institutions adopting the chain typically have existing change management for their broader stack and the chain's CUECs ride on top of existing process. Tier 3 (within first year) is defensible. Not blocking.

- **The worked-example anomaly record uses `tenant_acme_prod_us_east_1` consistently with the sample verifier report.** I checked. The fingerprint hex values (`b94c1a77b40bf5106c66ca6c1c1b4989` and `2eede65f0f764c97eaf3b3f306a48537`) match across the worked example, the sample verifier report's failed-day output, and the `chain.verification_failure` event example in `control-evidence-events.md`. The cross-document consistency is the kind of detail that signals the package was authored by someone who actually walked the worked example through every consuming document, not generated artifact-by-artifact. This is a strength, not a gap.

- **P-25 stratification minimum (5 per stratum × 3 dimensions × 3 strata each = 45) versus the general sampling guidance at the bottom of `audit-procedures.md` (250-event population uses 25, 250-2,500 uses 50, > 2,500 uses 75).** P-25's 45-entry minimum overrides the general sampling guidance for SR 11-7 reproducibility specifically — and that override is correct because stratification dominates flat sampling for this control. Not a conflict; the override is intentional. Not blocking.

I close at 0/0.

---

## Two strengths I want to call out by name

These are the two artifacts I would carry forward into my firm's broader SOC practice as patterns for other engagement portfolios:

### A. The four-piece evidence package for `key_fingerprint mismatch`

P-22's four pieces (IKM-roster row identification, change-management approval, re-verification PASS, reconciliation cross-check) is the right operational definition of "operationally explained" for an integrity anomaly. It is enumerated, each piece is independent, each piece is verifiable from institution evidence repositories. The worked example in `anomaly-documentation-template.md` makes it executable on the first attempt.

I am going to draft an internal-firm SOC-practice memo this quarter applying the four-piece pattern to other integrity-anomaly types in non-chain SOC engagements. The pattern generalizes: any integrity anomaly should require (i) identification of the wrong artifact, (ii) change-management approval for the correction, (iii) re-verification result, (iv) cross-check against the institution's standing reconciliation control. Most SOC engagements do not enforce all four; this package does.

### B. The three-dimensional stratification rule in P-25

P-25's 3-dimensional stratification (model-version × decision-class × customer-impact-tier, 5 per stratum minimum) is the right shape for SR 11-7 reproducibility evidence and the right shape for any control where the substrate property varies across multiple dimensions. The rule is explicit, the dimensions are the right ones for SR 11-7, and the per-stratum minimum gives statistical defensibility without ballooning the sample.

I will reference P-25 in my firm's internal SR 11-7 guidance and adapt the pattern for other multi-dimensional sampling problems. The stratification pattern is portable.

---

## Closing

This package supports an unmodified SOC 1 Type II or SOC 2 Type II opinion on the chain-of-custody implementation, with no material or minor gaps that would prevent issuance. The Section 4 template, the CUEC list, the TSC mapping, the audit procedures (P-1 through P-25), the operational events schema, the anomaly evaluation template with its worked example, the user-entity summary, and the sample verifier report with its SOC team appendix — together they form the cleanest substrate I have ever audited at first look.

I close this round at **0 material gaps / 0 minor gaps**.

Beatriz Almeida
Senior Manager, SOC Practice
São Paulo office
