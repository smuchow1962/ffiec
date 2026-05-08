# Round 13 — FFIEC Cybersecurity Specialist Examiner review

> **Reviewer.** Aditi Krishnan, FFIEC Cybersecurity Specialist Examiner with the FDIC.
> **Background.** Twelve years on the cyber-specialist track within the FDIC's Division of Risk Management Supervision. The last rotation was through the CIRCIA implementation working group bridging CISA, the federal banking agencies, and the FFIEC member-agency notification rules. My day-to-day examination work is FDIC-supervised state non-member banks and savings associations, so the reading angle leans into 12 CFR Part 304 Subpart C ("Computer-Security Incident Notification"), the parallel state-banking-supervisor coordination, and the CIRCIA 72-hour CISA path that runs alongside the 36-hour banking-agency path.
> **Lens.** NIST CSF 2.0 mapping; 36-hour computer-security incident notification cadence and clock-start determination; CIRCIA 72-hour parallel; threat-model Adversary I (regulator-fingerprint reception); cold-DR-key consumption as institution-side control evidence; supply-chain trust path (cosign / GPG / cold-DR / `slsa-verifier`); IR playbook scenario coverage; dual-algorithm transitional-period verifier-version overlap.
> **Materials reviewed.**
> - `spec/chain-of-custody-v1.md` §10 (operational requirements, including the new `regulator_fingerprint.*` triple) and §7 step 11 (dual-algorithm verifier dispatch with cases (a) through (e))
> - `docs/design/00-overview.md` — system shape and four-attack catalog (§2)
> - `docs/design/09-threat-model.md` — Adversaries A through I including the §2.9 institution-side reception procedure for the regulator-held fingerprint rotation, plus the explicit reception-failure sub-variant; R1–R13 register
> - `docs/regulator-pack/CSF-2.0.md` — operational-events binding table including the four new rows for the regulator-fingerprint lifecycle and the cold-DR consumption log
> - `docs/incident-response-playbook.md` — Scenarios 1 through 12, federal-regulator routing per charter, dual-algorithm verifier-version overlap edge case
> - `docs/cloud-hsm-guide.md` — conformance bar and per-provider provisioning
> - `docs/supply-chain.md` — cold-DR-key lifecycle and institution-side consumption; `slsa-verifier` trust-anchor row; mirror WORM and continuous-monitoring control
> - `docs/regulator-pack/ai-policy-alignment.md` — Treasury RMF, EU AI Act Art. 12 / 14 / 26, DORA Art. 6 / 8 / 17 / 28
> - `docs/edge-and-federated-ai.md` — edge / federated / on-device patterns
>
> Per the engagement brief, I have not consulted prior-round feedback. This is a first-look review against the shipped artifacts on the morning of 2026-05-06.
>
> **Stopping criterion.** 0 Gap, 0 Partial.
> **Disposition this round.** **4 Gap, 2 Partial.** Detailed in §3.

---

## 1. Headline impression

The corpus is in good shape against my cyber-specialist lens. Three constructions are above the bar I usually see at this level of supervision:

- The four-attack-to-CSF-subcategory table at `CSF-2.0.md` lines 37–43 is the artifact a CSF examiner builds an evidence matrix from. Each attack approach resolves to a primary subcategory, a load-bearing operational event, and an IR scenario. The table is dense enough that I can pull a sample at the table row level rather than spelunking through the prose.
- The Adversary I catalog entry in `09-threat-model.md` §2.9 is the right shape. The reception procedure carries three institution-defined operational events (`regulator_fingerprint.rotation_received`, `.rotation_validated`, `.installed`), each sample-testable, plus a reception-failure sub-variant for the forged-notice path. Forging a regulator-side rotation notice is exactly the threat I would press on at the FDIC; the catalog entry names it explicitly.
- The Scenario 12 three-branch triage tree for spec §7 step 11 case (e) co-signed seal failure is operationally honest about ambiguity. The verifier produces a mechanical signal (one algorithm validated, one did not). The IR program produces the interpretation (published break / per-algorithm key compromise / under investigation), each branch with its own clock-start posture and a bounded 48–72h window on the third branch so the clock-start cannot be deferred indefinitely. The bounded window is the part most playbooks omit.

My findings sit on a tighter set of edges. New artifacts (Scenarios 11 and 12, the reception-failure sub-variant, the cold-DR consumption log, the `regulator_fingerprint.*` triple) were added in this iteration. The artifacts themselves are correctly specified. What is missing is the connective tissue — the points where consolidated reference tables and downstream consumption sites should reflect the new artifacts but do not. Three of those are Severe in operation; one is Severe at the SOC-evidence layer; two are Partial extensions to existing matrices.

---

## 2. What I tested (PASS rows omitted for brevity; failures detailed in §3)

### 2.1 NIST CSF 2.0 mapping

| Test | Reference | Result |
|---|---|---|
| PR.DS-06 headline mapping with chain as integrity-checking control | `CSF-2.0.md` §"Headline mapping" | PASS |
| Four-attack-to-subcategory table with operational evidence and IR scenario per row | `CSF-2.0.md` lines 37–43 | PASS |
| Operational-events binding table covers the `regulator_fingerprint.*` triple | `CSF-2.0.md` lines 88–90 | PASS |
| Cold-DR `KEY-DR-DRYRUN-{year}.asc` consumption log bound to GV.SC-04 + ID.RA-09 with IR Scenario 11 pairing | `CSF-2.0.md` line 91 | PASS |
| GOVERN-function alignment one-pager | `CSF-2.0.md` §"GOVERN function alignment" | PASS |

### 2.2 Threat model and Adversary I (regulator-fingerprint reception)

| Test | Reference | Result |
|---|---|---|
| Adversary I named and first-class in catalog | `09-threat-model.md` §2.9 | PASS |
| Three institution-defined operational events for the reception procedure | `09-threat-model.md` §2.9 steps 1–3; spec §10.2 | PASS |
| Reception-failure sub-variant cross-reference resolves in the IR playbook | `09-threat-model.md` §2.9 step 2 → `incident-response-playbook.md` Scenario 11 | **FAIL — see Gap A-1.** Cross-reference target absent. |
| Examiner-laptop hygiene scoped OUT of institution's CC6 | `09-threat-model.md` §2.9 closing paragraph | PASS |
| R13 algorithm-confusion at seal layer closed by spec §4.3 sign_payload extension | `09-threat-model.md` §6 R13 | PASS |

### 2.3 Spec §7 step 11 dual-algorithm cases

| Test | Reference | Result |
|---|---|---|
| Cases (a) through (e) named with --strict and non-strict dispositions | spec §7 step 11 | PASS |
| Case (e) severity-is-Severe-regardless-of-bracket convention | spec §7 step 11 case (e) | PASS |
| Examiner working-paper convention records BOTH algorithms' validation results | spec §7 step 11 closing paragraph | PASS |
| IR scenario named for case (e) | `incident-response-playbook.md` Scenario 12 | PASS |

### 2.4 IR playbook — scenario coverage and clock-start determination

| Test | Reference | Result |
|---|---|---|
| Scenarios 1 through 12 cover the chain-detected event surface | `incident-response-playbook.md` §"Common scenarios" | PASS |
| Scenario 11 trigger variants (a) failed dry-run, (b) dual-compromise, (c) missed window | `incident-response-playbook.md` Scenario 11 | PASS for the three named triggers |
| Scenario 11 sub-variant for trust-anchor reception failure (the threat-model cross-reference target) | `incident-response-playbook.md` Scenario 11 | **FAIL — see Gap A-1.** |
| Scenario 12 three-branch triage tree | `incident-response-playbook.md` Scenario 12 | PASS |
| 36-hour cyber-incident triage matrix covers all twelve scenarios | `incident-response-playbook.md` lines 324–336 | **FAIL — see Gap A-2.** Matrix terminates at Scenario 10. |
| CIRCIA 72-hour vs FFIEC 36-hour matrix covers Scenarios 11 and 12 | `incident-response-playbook.md` lines 371–377 | **FAIL — see Gap A-3.** CIRCIA matrix carries five rows (Scenarios 4, 5, 6, 7, 9). |
| Federal-regulator routing per institution charter | `incident-response-playbook.md` lines 347–358 | PASS for federal routing; **PARTIAL — see Partial A-5** for per-agency CFR citation and state-side cadence. |
| Dual-algorithm verifier-version overlap matrix | `incident-response-playbook.md` lines 360–367 | **PARTIAL — see Partial A-6.** Two rows present; case-(e) interaction missing. |

### 2.5 Supply-chain cold-DR-key consumption (institution-side)

| Test | Reference | Result |
|---|---|---|
| Cold-DR-key lifecycle (60-month rotation, separated key holders, annual dry-run, hardware health check) | `supply-chain.md` §"Cold-disaster-recovery key lifecycle" | PASS |
| Institution-side consumption procedure with five steps, 30-day grace, IR Scenario 11 branch (c) for missed window, control-evidence binding | `supply-chain.md` §"Institution-side consumption of cold-DR-key dry-run attestation" | PASS |
| Cold-DR public-key fingerprint as institution-side cached trust anchor in trust-anchor table | `supply-chain.md` §"Trust anchors" line 158 | PASS |
| `slsa-verifier` co-equal trust artifact alongside cosign and GPG with deployment-gating documentation | `supply-chain.md` §"Institutional consumption of the SLSA L3 attestation" | PASS |

### 2.6 SOC-evidence schema (control-evidence operational events)

| Test | Reference | Result |
|---|---|---|
| Schema row for `regulator_fingerprint.rotation_received` with field vocabulary matching threat model §2.9 step 1 | `docs/soc-pack/control-evidence-events.md` | **FAIL — see Gap A-4.** Event named in spec and CSF-2.0 binding; schema entry absent. |
| Schema row for `regulator_fingerprint.rotation_validated` with field vocabulary matching threat model §2.9 step 2 | `docs/soc-pack/control-evidence-events.md` | **FAIL — see Gap A-4.** Same pattern. |
| Schema row for `regulator_fingerprint.installed` with field vocabulary matching threat model §2.9 step 3 | `docs/soc-pack/control-evidence-events.md` | **FAIL — see Gap A-4.** Same pattern. |
| Schema or field vocabulary for the `KEY-DR-DRYRUN-{year}.asc consumption log` | `docs/soc-pack/control-evidence-events.md` | **FAIL — see Gap A-4.** Bound in CSF-2.0; schema absent. |

### 2.7 Boundary correctness for FDIC-supervised institutions

| Test | Reference | Result |
|---|---|---|
| Per-agency CFR citations for the joint computer-security incident notification rule (12 CFR Part 53 OCC; 12 CFR Part 225 Appendix F FRB; 12 CFR Part 304 Subpart C FDIC) | `incident-response-playbook.md` lines 322 and 413 | **PARTIAL — see Partial A-5.** Doc uses "FFIEC computer-security rule" generically; per-agency citation absent. |

---

## 3. Findings

### Gap A-1 (Severe). The IR playbook is missing the "trust-anchor reception failure" sub-variant that the threat model explicitly cross-references.

`09-threat-model.md` §2.9 step 2 (line 202) names the failure path for a forged regulator-fingerprint rotation notice: *"if validation FAILS (forged notice suspected), the institution does NOT install the fingerprint; activates IR Scenario 11 sub-variant for 'trust-anchor reception failure' (project- or regulator-side governance event); cross-checks via out-of-band regulator contact."*

I traced the cross-reference. Scenario 11 (`incident-response-playbook.md` lines 261–281) carries three triggers — (a) failed dry-run, (b) dual-compromise activation, (c) missed dry-run window — and each is properly handled. None of the three is the trust-anchor reception failure sub-variant. The sub-variant does not appear anywhere else in the playbook either.

**Why Severe rather than Partial.** Forged-notice handling is the load-bearing path for the Adversary I attack class. The institution receives a rotation notice that LOOKS authentic, fails the validation check at threat-model §2.9 step 2, and then needs the IR scenario to tell it what to do next. The threat model points to the IR playbook for that procedure. The IR playbook does not have it. An institution following the cross-reference will arrive at Scenario 11 and find no procedure for the case the threat model said the procedure exists for. The downstream effect is the institution's general IR framework swallowing the forged-notice case without chain-specific containment guidance — the verifier's trust anchor sits in a state where the institution cannot be sure whether to install the new fingerprint or freeze, and the runbook silence is exactly the wrong artifact to encounter under incident pressure.

The corresponding examiner-side question — "show me the runbook entry for a forged regulator-fingerprint rotation notice" — has no artifact to point to. I would write this as a Severe finding in working papers because the forged-notice path is the most consequential of the three Adversary I scenarios and the documentation chain breaks at the point where the institution most needs it.

**Disposition.** Add a sub-variant to Scenario 11 — call it Scenario 11(d) or a named "Trust-anchor reception failure" sub-section — covering:

- **Trigger.** Validation step at threat-model §2.9 step 2 returns FAIL (forged notice suspected); the corresponding `regulator_fingerprint.rotation_validated` event has `overall_validation = FAIL`.
- **Investigation.** Cross-check out-of-band against the regulator's authenticated channel; pull the regulator's standing identity GPG signature; confirm whether the regulator has actually issued a rotation; preserve the forged notice and its `notice_artifact_sha256` as forensic evidence.
- **Containment.** Do NOT install the fingerprint; freeze the verifier's trust-anchor configuration at the previous fingerprint; pause new-version verifier deployments until the determination resolves.
- **Remediation.** Await regulator confirmation of true rotation status; if confirmed forgery, treat the determination as a candidate to escalate to IR Scenario 4 disposition because the forging party may have broader access to regulator-side material.
- **Notification.** Clock-start trigger: STARTS the 36-hour clock when the institution determines the notice is forged AND the determination implies an external party with regulator-side knowledge is operating. That meets the "computer-security incident affecting safety/soundness" determination point under 12 CFR Part 304 Subpart C / 12 CFR Part 53 / 12 CFR Part 225 Appendix F.

### Gap A-2 (Severe). The 36-hour cyber-incident notification triage matrix is missing rows for Scenarios 11 and 12.

The matrix at `incident-response-playbook.md` lines 324–336 carries one row per scenario for Scenarios 1 through 10. Scenarios 11 and 12 were added to the playbook this iteration and are NOT in the matrix. The matrix's stated purpose, in the same section's lead paragraph (line 322), is to name "what triggers the determination and therefore the clock-start" per chain-detected event type. Two scenarios are now outside that scope.

**Why this matters in practice.** The FFIEC rule's "determination point" is the load-bearing concept. The matrix is where examiners and IR Commanders look it up at incident time. If a Scenario 12 alert fires at 02:00 and the IR Commander needs to know the clock-start posture, the answer needs to be in the matrix at the bottom of the playbook, not buried in the §"Scenario 12" prose at the top. A consolidated reference table that omits two scenarios forces the IR Commander to read two scenario sections to find what should be in two rows of one table. Under incident pressure that produces hesitation; hesitation produces clock-start drift.

**Why Severe rather than Partial.** "Show me the institution's consolidated clock-start matrix per scenario" is a question I ask routinely at examination. An incomplete matrix is a finding because the matrix is the artifact the examiner samples. The playbook explicitly states (line 322) that the matrix lists the determination triggers; the omission of Scenarios 11 and 12 contradicts that scope.

**Disposition.** Add two rows to the matrix:

- **Scenario 11 (project-side trust-anchor degradation, including the trust-anchor reception failure sub-variant per Gap A-1).** Default disposition: project-side governance, no FFIEC clock by itself. Clock-start trigger: does not start the 36-hour clock unless the trust-anchor reception failure sub-variant per Gap A-1 fires (forged-notice handling), in which case the clock starts at the institution's determination of forgery.
- **Scenario 12 (co-signed seal failure case (e)).** Default disposition: branch-dependent per the three-branch triage tree. Clock-start trigger: branch (i) NO clock (published-break case; coordinate with regulator on migration timeline); branch (ii) STARTS at credibility determination per Scenario 4; branch (iii) STARTS at investigation-conclusion determination, bounded 48–72 hours per the scenario text.

### Gap A-3 (Severe). The CIRCIA 72-hour matrix is missing Scenarios 11 and 12.

The CIRCIA matrix at `incident-response-playbook.md` lines 371–377 carries rows for Scenarios 4, 5, 6, 7, and 9. Scenarios 11 and 12 are absent. (Scenarios 1, 2, 3, 8, and 10 are also absent, but those are partial-coverage residuals from earlier rounds; the round-13-introduced Scenarios 11 and 12 are the new omissions.)

**Why the omissions matter, specifically for Scenarios 11 and 12 in the CIRCIA dimension.** CIRCIA reporting is a separate notification path on a different clock to a different recipient (CISA, not the primary banking regulator). Counsel works from the matrix when deciding which path applies. An incomplete CIRCIA matrix produces under-reporting on the CIRCIA path when a CIRCIA-qualifying event fires and counsel does not see it in the matrix. It also produces over-reporting when a non-qualifying event fires and counsel cannot see "No" in the matrix and defaults to filing on conservative grounds. Both directions are findings the cyber-specialist examiner challenges.

**The substantive call on the missing rows.**

- **Scenario 12 branch (ii) — per-algorithm signing-key compromise — is a CIRCIA candidate.** A confirmed compromise of an algorithm-specific signing key is HSM-side compromise of cryptographic material binding the institution's AI-decision audit trail. CIRCIA's "substantial cyber incident" definition reaches this: the institution's ability to produce trustworthy AI-decision evidence is degraded across the affected algorithm, and at scale this is reportable to CISA. The matrix should explicitly say "Yes" for branch (ii). For branch (i), the matrix should say "No" because the published break is a known cryptographic-deprecation event, not a substantial cyber incident at the institution. For branch (iii), the matrix posture follows Scenario 4's posture pending determination.
- **Scenario 11 (cold-DR fallback) is NOT a CIRCIA event at the institution layer** — but the matrix should explicitly say so to prevent over-reporting. The cold-DR fallback fires on project-side governance events; the institution is downstream. Listing the row with "No" is the artifact counsel reaches for. The exception case is the trust-anchor reception failure sub-variant per Gap A-1 — if that sub-variant determines forgery and the forging party is suspected to have broader access, escalate per Scenario 4 and the CIRCIA posture follows Scenario 4.

**Why Severe rather than Partial.** Same reasoning as Gap A-2 — the matrix is the consolidated reference under incident pressure. The omissions are not abstractly addressable through "counsel decides per incident"; the matrix exists precisely so counsel does not have to make a category determination from first principles when a clock has started.

**Disposition.** Add three rows (the Scenario 12 row breaks into the three branches because the branch postures differ):

- **Scenario 11 (project-side trust-anchor degradation, including the trust-anchor reception failure sub-variant per Gap A-1).** FFIEC 36-hour: per Gap A-2 disposition. CIRCIA 72-hour: No (project-side governance; not an institution-side substantial cyber incident). Exception: if the reception-failure sub-variant fires and the forging party is suspected to have broader access, escalate per Scenario 4 and the CIRCIA matrix follows Scenario 4.
- **Scenario 12 branch (i) (published algorithm break).** FFIEC 36-hour: No. CIRCIA 72-hour: No.
- **Scenario 12 branch (ii) (per-algorithm signing-key compromise).** FFIEC 36-hour: Yes. CIRCIA 72-hour: Yes (substantial cyber incident at HSM-side cryptographic-material layer; the affected algorithm's audit-trail integrity is degraded across the institution's deployment).
- **Scenario 12 branch (iii) (under investigation).** FFIEC 36-hour: per the bounded 48–72h determination window. CIRCIA 72-hour: posture follows the FFIEC determination — once the FFIEC clock starts, evaluate CIRCIA at the same determination point.

### Gap A-4 (Severe). The SOC-evidence schema document does not carry schema rows for the four new operational events the spec and CSF binding name.

`docs/soc-pack/control-evidence-events.md` is the document SOC engagements consume to standardise how operational events are read as control evidence. Its lead paragraph (lines 5–7) explicitly says the schema is the artifact making consumption mechanical. The document carries schema entries for `ledger.startup`, `ledger.hsm_session_opened`, `seal.job_started`, `seal.job_completed`, `seal.job_failed`, `chain.verification_failure`, the HSM operations, configuration, master-key rotation, audit-file truncation, master-key retired, and reconciliation. It does NOT carry schema entries for:

- `regulator_fingerprint.rotation_received` (named in spec §10.2 and bound in `CSF-2.0.md` line 88)
- `regulator_fingerprint.rotation_validated` (named in spec §10.2 and bound in `CSF-2.0.md` line 89)
- `regulator_fingerprint.installed` (named in spec §10.2 and bound in `CSF-2.0.md` line 90)
- The annual `KEY-DR-DRYRUN-{year}.asc` consumption log (bound in `CSF-2.0.md` line 91 and operationally named in `supply-chain.md` §"Institution-side consumption")

**Why this is a SOC-evidence problem and a cyber-specialist problem.** The CSF binding table tells the examiner which subcategory each event supports. The SOC schema tells the SOC engagement what the event's `fields` block is supposed to contain so the engagement can compare event instances mechanically across the reporting period. Without the schema rows, the SOC engagement cannot produce mechanical comparability for the four new events; the engagement reverts to ad-hoc field-extraction per institution. That is the failure mode the SOC schema document was written to prevent. The threat-model §2.9 reception procedure names the field vocabulary for each event (received_at, received_via, regulator_signing_identity, notice_artifact_sha256, etc.); that vocabulary needs to make the round-trip into the schema document so the SOC engagement consumes the same field names the threat model specifies.

**Why Severe rather than Partial.** The Adversary I reception procedure is the load-bearing institution-side response to the most consequential examiner-side adversary the threat model carries. The SOC engagement is the institution's substantive operational evidence that the procedure is operated. A schema gap in the SOC-evidence document is a control-evidence gap that propagates into every SOC engagement covering the round-13 spec version forward. The cold-DR consumption log row is similar — bound in CSF as the load-bearing GV.SC-04 evidence, but no schema for the SOC engagement to consume.

**Disposition.** Add four schema entries to `docs/soc-pack/control-evidence-events.md`:

- `regulator_fingerprint.rotation_received` with fields `received_at` (RFC 3339 UTC), `received_via` (channel name string), `regulator_signing_identity` (string), `notice_artifact_sha256` (hex string).
- `regulator_fingerprint.rotation_validated` with fields `validated_at` (RFC 3339 UTC), `validation_paths_attempted` (list of strings), `validation_results` (per-path PASS/FAIL list), `overall_validation` (PASS/FAIL).
- `regulator_fingerprint.installed` with fields `installed_at` (RFC 3339 UTC), `old_fingerprint_archived_at` (RFC 3339 UTC), `new_fingerprint_active_from` (RFC 3339 UTC), `change_management_record_id` (string).
- The `KEY-DR-DRYRUN-{year}.asc consumption log` shape — either as an operational event the institution emits when it consumes the year's dry-run attestation, or (if the institution-side consumption log is a different artifact than an operational event) as a documented log-shape with the same field-vocabulary discipline (`consumed_at`, `attestation_year`, `signature_validation_result`, `archived_at`, `archive_location`).

The fourth item is structural — the spec §10.2 enumeration treats the consumption log as institution-defined; the SOC schema document needs to call that out and either define the operational event shape if the project recommends one, or document the field-vocabulary discipline the institution should apply if it shapes the log institution-side.

### Partial A-5. Federal-regulator routing matrix is correct for federal routing but is silent on (a) per-agency CFR citation for the joint computer-security incident rule, and (b) the parallel state-banking-supervisor cadence for FDIC-supervised state non-member banks.

The routing table at `incident-response-playbook.md` lines 349–356 names the primary federal regulator per institution charter, including the FBO subsidiary case. The federal coverage is complete and correct. Two specific operational details are missing.

**Sub-issue (a): per-agency CFR citation.** The text uses "FFIEC computer-security rule" and "FFIEC 36-hour rule" throughout. The rule was adopted jointly by the three federal banking agencies as separate codifications: 12 CFR Part 53 (OCC), 12 CFR Part 225 Appendix F (Federal Reserve), and 12 CFR Part 304 Subpart C (FDIC). The FFIEC member-agency working group coordinated the adoption; the FFIEC itself did not adopt a rule, and an institution that cites "FFIEC computer-security rule" in its 36-hour notification will be politely asked by the supervising agency's examiner to use the agency's correct citation. The doc should carry the per-agency CFR citation alongside the routing table so the institution's runbook references the right rule per its charter.

**Sub-issue (b): state-banking-supervisor cadence for FDIC-supervised state non-member banks.** The routing row says "FDIC (federal); state banking regulator (state)" but does not name the state-side cadence. Several state banking departments operate their own computer-security incident notification rules with cadences ranging from 24 to 72 hours; some states have stricter cadence than the federal 36-hour rule. New York DFS Part 500.17, for example, carries a 72-hour cadence with a different determination point than the federal rule. The IR Commander needs to know which clock starts first when running an FDIC-supervised state non-member bank's incident — the state-side clock often does. The same posture applies for state credit unions where the state credit-union regulator may have its own cadence distinct from NCUA's parallel expectation.

**Why Partial rather than Gap.** The federal-routing coverage is correct; the omissions are operational specifics layered on top of a correct framework. An institution with a sophisticated counsel team will navigate this correctly even with the matrix as currently written. An institution running a smaller IR program — the cyber-specialist examiner's typical population at the FDIC-supervised tier — benefits from explicit citations and an explicit state-side-cadence column.

**Disposition.** Two additions:

1. Add a footnote or fourth column to the routing table naming the per-agency CFR citation: National banks → 12 CFR Part 53; FRB-supervised institutions → 12 CFR Part 225 Appendix F; FDIC-supervised state non-member banks → 12 CFR Part 304 Subpart C; NCUA-supervised credit unions → 12 CFR Part 748 Appendix B (the NCUA's parallel incident-notification expectation).
2. Add a paragraph after the routing table naming the state-side coordination requirement for state-chartered institutions: the institution's IR runbook MUST pre-record the state-side cadence and clock-start trigger alongside the federal cadence. The state-by-state catalogue is out of scope for this document; naming the requirement that the state-side cadence be pre-recorded is in scope.

### Partial A-6. Dual-algorithm verifier-version timing-overlap matrix is correct for the two PASS / PASS-WITH-ANOMALY case-(b) interactions but does not cover the case-(e) interaction.

The matrix at `incident-response-playbook.md` lines 362–365 covers two posture-disagreement scenarios both involving spec §7 step 11 case (b) "partial-coverage seal." That coverage is correct.

A third interaction the matrix should name: **case-(e) co-signed seal failure interacting with the verifier-version overlap.** When the institution's verifier supports algorithm Y and the regulator's verifier still supports only algorithm X, and a co-signed seal fires the case-(e) failure path on the institution's verifier (one algorithm validated, one did not — Scenario 12 fires institution-side), the regulator's verifier sees only the X-algorithm signature and reports PASS. The institution has a clock-starting determination on its end (per Scenario 12's three-branch triage); the regulator has no anomaly to investigate from its verifier output.

**Why this is the harder posture-disagreement case.** The PASS / PASS-WITH-ANOMALY case (b) interactions are control-completeness misalignments. The case-(e) interaction is qualitatively different: the institution sees a security-relevant FAIL on the Y-algorithm signature; the regulator sees a clean PASS on its X-algorithm-only verifier. Without a documented disposition, the institution may delay surfacing the case-(e) determination because the regulator-facing artifact (the regulator's verifier output) shows no problem. That delay is exactly the kind of clock-start deferral the bounded 48–72h window in Scenario 12 was meant to prevent. The matrix needs a row that says explicitly: the institution surfaces the case-(e) determination to the regulator on the institution's clock, not the regulator's clock.

**Why Partial rather than Gap.** The matrix exists; the structure is right; the case-(e) row slots into the existing matrix without restructuring. The two PASS / case-(b) rows are correct; the missing case-(e) row is a coverage extension, not a structural issue.

**Disposition.** Add a third row to the verifier-version timing-overlap matrix:

| Posture disagreement | Disposition |
|---|---|
| Institution: FAIL case (e) (one algorithm validated, one did not — Scenario 12 fires). Regulator: PASS under its X-only posture (regulator's verifier sees only the X-algorithm signature, which validated). | The institution's case-(e) determination governs the institution's clock-start per Scenario 12 three-branch triage. The institution MUST surface the case-(e) determination to the regulator within the institution's clock-start window even though the regulator's verifier shows no anomaly; the institution's IR program does NOT defer the surfacing on the basis of the regulator's verifier output. The institution's working paper records both verifier outputs (institution's case-(e) FAIL plus regulator's case-(d) PASS) so the regulator has the full picture; the institution proposes a working-paper procedure for the regulator either to obtain the Y-algorithm-aware verifier from the project supply chain or to accept the institution's Y-algorithm-aware verifier output as supplementary evidence. |

The added row makes the verifier-version overlap matrix complete across cases (a) through (e). The two existing rows handle case (b); the new row handles case (e); cases (a), (c), (d) do not produce posture disagreements that require named disposition (case (a) is PASS on both verifiers when both algorithms are visible; case (c) FAILs on both verifiers under --strict; case (d) is the single-algorithm posture default which the verifier-version overlap does not affect).

---

## 4. What works particularly well

Specific call-outs on what landed.

**The reception procedure for the regulator-held fingerprint rotation has the right shape.** Five steps — receive, validate, install, re-validate historical reports, document — each with operational evidence. The three operational events (`regulator_fingerprint.rotation_received`, `.rotation_validated`, `.installed`) form a sample-testable audit trail. The institution-side procedure for a regulator-side governance event is rare and well-handled. The defects (Gap A-1, Gap A-4) are downstream-consumption gaps; the procedure itself is correctly specified at the threat-model layer and at the spec §10.2 layer.

**The cold-DR-key consumption procedure binds project-side governance to institution-side examinable control evidence.** The institution downloads the annual `KEY-DR-DRYRUN-{year}.asc` attestation, validates the signature, archives in the control-evidence repository, treats a missed window as IR Scenario 11 branch (c). The CSF binding to GV.SC-04 + ID.RA-09 closes the loop — a CSF examiner asking "show me the institution exercises the cold-DR fallback path" gets the consumption log as the artifact. This is a model for how project-side governance becomes institution-side control evidence. Gap A-4's schema-row omission does not detract from the design; it is a follow-through gap at the SOC-evidence layer.

**Scenario 12's three-branch triage is operationally honest about ambiguity.** Branch (i) names the published-break case where the regulator already knows. Branch (ii) names the per-algorithm key compromise where Scenario 4 takes over. Branch (iii) names the under-investigation case with a bounded 48–72h window so the clock-start cannot be deferred indefinitely. The bounded window is what most playbooks omit; including it makes the case-(e) clock-start determination work under counsel pressure.

**The federal-regulator routing matrix is pre-documented per charter type, including the FBO subsidiary case.** Pre-documenting the routing matrix removes the routing-decision-under-pressure problem at incident time. The two operational-detail omissions noted in Partial A-5 do not detract from the value of the existing matrix; they extend it.

**The Adversary I framing — examiner-laptop hygiene scoped OUT of the institution's CC6 evaluation — is the correct boundary call.** A compromised examiner laptop can produce a false report regardless of the chain's integrity properties; the mitigation is the FFIEC's own examiner-IT program, not an institution-side control. The §2.9 closing paragraph names this explicitly so an institution-side reviewer cannot misread the residual as an institution-owned control gap. From my position at the FDIC, this is the right call: examiner-IT lifecycle is the agency's responsibility, and conflating it with institution CC6 obligations confuses the boundary.

---

## 5. Residual concerns I am NOT escalating

Documented residuals at the spec-version level, not gaps at the round-13 level.

- **HSM physical compromise (R3).** Accepted residual. FIPS 140-2 L3 tamper detection plus immediate revocation is the operational response. Nation-state-level adversary cannot be eliminated.
- **Cryptographic primitive break (R6).** Accepted residual. The 30-day spec-patch SLA, the algorithm-rotation provision, and the dual-algorithm coexistence design are the operational response. The migration lag is operational, not technical.
- **Master-key compromise (R2).** Accepted residual with documented mitigation: weekly fingerprint reconciliation per spec §10.1 bounds the detection window to ≤ 1 week plus remediation. Events captured during the compromise window are repudiable; forensic analysis (network logs, process audit logs) is the supplementary evidence path. The chain construction does not eliminate the compromise window; the threat model accepts it.
- **Application-host compromise (R1, Adversary F).** Accepted residual. Forward-only forgery, bounded by detection window. Standard host hardening plus rotation on detection is the response.
- **Edge-device physical compromise (R12).** Deferred to v1.1. Compensating controls documented in `edge-and-federated-ai.md` (Pattern A: per-device IKM in TPM/secure-enclave; reconciliation cadence baseline tuned per fleet).
- **Examiner-laptop hygiene.** Out of scope for institution's CC6 evaluation per `09-threat-model.md` §2.9 closing paragraph. Correctly framed as an FDIC-side / regulator-side IT control, not an institution-side gap.
- **AI decision correctness.** Out of scope by design; SR 11-7 model validation is the separate workflow.
- **State-by-state computer-security cadence catalogue.** Out of scope for the IR playbook itself; naming the requirement that the state-side cadence be pre-recorded (per Partial A-5 sub-issue (b)) is in scope.

---

## 6. Disposition

**4 Gap, 2 Partial.** Does not meet the 0/0 stopping criterion.

The four gaps are connective-tissue holes between artifacts added this iteration and their downstream consumption sites:

- Gap A-1: threat-model cross-reference target (the trust-anchor reception failure sub-variant) does not exist in the IR playbook.
- Gap A-2: 36-hour matrix terminates at Scenario 10; new Scenarios 11 and 12 are not in the consolidated reference table.
- Gap A-3: CIRCIA matrix omits the round-13-introduced scenarios.
- Gap A-4: SOC-evidence schema document carries no entries for the four new events the spec and CSF binding name.

Each gap is closed by adding the corresponding row, sub-variant, or schema entry; no design change is required.

The two partials are operational-detail extensions to existing artifacts (per-agency CFR citations and state-side cadence in the federal-regulator routing matrix; case-(e) interaction row in the verifier-version overlap matrix). Each partial is closed by additive edits.

The corpus is close to PASS for this lens. The six items above are the remaining work. The design constructions themselves are correct; the gaps are at the points where the new constructions connect to the consolidated reference tables and to the SOC-evidence consumption layer. That is a classic "new artifact added; downstream consumption sites not extended" pattern, addressable by additive edits within the existing structure.

---

## 7. Sign-off

**Reviewer.** Aditi Krishnan, FFIEC Cybersecurity Specialist Examiner with the FDIC.
**Date.** 2026-05-06.
**Round.** 13.
**Disposition.** 4 Gap, 2 Partial. Does not meet 0/0.
