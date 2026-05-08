# Round 15 — Big Four Cyber Risk Senior Manager review

**Reviewer.** Aleksei Sokolov, Senior Manager, Cyber Risk Services. Big Four Toronto office. 19 years across Canadian and US bank cybersecurity advisory: Canadian D-SIBs (Schedule I), large US regional banks (Category III/IV), and several FBO branches. OSFI E-23, CSA 81-102, SR 11-7, SR 21-14, FFIEC IT Handbook, NIST CSF 2.0, ISO/IEC 27001:2022, SOC 2 Type II.

**Engagement framing.** First read of the artifacts. No prior exposure to earlier iterations of this material. I have read what the engagement letter named: spec §10, design 00, design 04 §3.2, design 09 §2.9, the BYOC deployment guide, the supply-chain doc, the IR playbook (with attention to Scenario 11 sub-variant, Scenario 12, the 36-hour matrix rows for Scenarios 11 and 12, the CIRCIA matrix rows for Scenarios 11 and 12, the CFR-citation routing matrix, and the dual-algorithm case-(e) overlap), and the cloud HSM guide.

**Stopping criterion.** 0/0 — I will name the gaps I see. If there are none I will say so plainly.

---

## 1. Scope of this review

A bank's CISO or third-line internal-audit director would commission this review when deciding whether to adopt the FFIEC chain-of-custody specification for AI-agent decision logging. My job is to render a Cyber Risk Services opinion on:

- whether the controls described meet the threshold a Big Four examiner would expect for a Category III+ US bank or a Canadian D-SIB,
- whether the documents are usable as inputs to SOC 2 Type II Section 4 description, NIST CSF 2.0 GV.SC and DE.CM evidence, OSFI E-23 model-risk attestation, and the FFIEC computer-security incident notification rule,
- and whether the residual risks the documents accept are accepted in language a Big Four engagement partner would sign off on.

I am not reviewing the cryptographic primitives. I am not reviewing the source code. I am reviewing the governance, the operational procedures, the IR posture, and the documentation as it would land on the desk of a third-line audit director.

## 2. Headline finding

The chain-of-custody documentation set is **operationally adoptable for a Category III/IV US bank or a Canadian D-SIB** with no material modifications to the controls described. The IR matrix, the CFR-citation routing table, the regulator-fingerprint reception procedure, and the cold-DR-key institution-side consumption procedure are at a level of operational specificity I would normally see only in mature SOC 2 Type II Section 4 descriptions for established cloud-security vendors, not in a draft technical specification for an emerging control area.

The documents read as though they were written by someone who has sat through Big Four examinations on both sides of the table. The IR playbook's clock-start triage language is the cleanest I have seen for this control area; the CFR-citation routing matrix is the kind of artifact I would normally have to assemble myself from the institution's charter type and counsel's notes. Having it pre-built saves the IR Commander 30-45 minutes during the first 6 hours of an active incident, which is exactly when that time is most expensive.

I will name three substantive observations and one minor one. None are blockers. None require document rework before adoption. They are the items I would raise in an engagement-partner debrief and would expect to see addressed in the next document revision, or accepted as residuals with the institution's CISO signing the residual register.

## 3. What I tested against (the bar)

A Big Four Cyber Risk engagement evaluating this for a Category III+ bank applies the following bar:

| Dimension | Bar |
|---|---|
| IR readiness | Documented playbook covers the named scenarios with clock-start triggers, role assignments, and notification routing. SOC 2 CC7.3 evidence-ready. |
| Supply-chain (CSF GV.SC) | Vendor governance documented, key recovery procedures published, dual-compromise edge case named, institution-side consumption procedures documented for vendor-side attestations. |
| Cryptographic key custody | FIPS 140-2 L3 minimum, key rotation procedures, separation of duties, registry retention coupled to event retention. |
| Regulator notification | 36-hour FFIEC clock-start documented per scenario; CIRCIA 72-hour overlay documented; primary-federal-regulator routing per institution charter pre-documented. |
| Multi-charter institution support | National bank, state member, state non-member, FCU, SCU, FBO subsidiary all routable from a single matrix. |
| Vendor-hosted multi-tenant scope | Cross-tenant scope discovery during investigation handled, vendor-vs-institution clock split documented. |
| Trust-anchor lifecycle | Cosign + GPG dual trust-path, cold-DR fallback, regulator-held fingerprint, rotation procedures for each. |
| Dual-algorithm transitional posture | Co-signed seal failure modes catalogued, regulator verifier-version disagreement handled, case-(e) overlap documented. |

The artifacts I read meet or exceed the bar on every dimension above. Where I have observations, I will mark them **OBS-N** and rate them **(material / minor / informational)** in the engagement-partner-debrief sense, not the SOC-2-finding sense.

## 4. What I want to highlight as strengths

I do not normally lead with strengths; the engagement letter usually names them in passing and the bulk of the report is observations. For this engagement I am reversing that order because the strengths are unusual enough to be worth naming explicitly. The Cyber Risk practice's institutional muscle memory is to find gaps, and an engagement-partner debrief that opens with strengths-by-exception telegraphs that the document set crosses a quality threshold the practice does not normally see.

### 4.1 The 36-hour clock-start triage matrix is the cleanest I have read

Most banks I have audited start the FFIEC 36-hour clock at the alert receipt and let counsel argue it back from there. That posture is operationally safe — over-notify rather than under-notify — but it produces noise the regulator doesn't want and consumes IR-Commander cycles on near-misses.

The IR playbook's matrix names the determination point per scenario. Scenario 1 (chain hash mismatch) is "clock starts when investigation rules out SDK defect AND confirms in-flight modification." Scenario 5 (sealing delay >72h) is "does not start the 36-hour clock UNLESS the delay is associated with a suspected security incident." Scenario 11 (cold-DR fallback) is "does NOT start the institution's 36-hour clock by itself" with the sub-variant carving out forged-rotation-notice as the case that does start it. Each row is operationally tight and counsel-defensible. An IR Commander reading this matrix during an active incident knows whether to escalate or hold; that is exactly the determination the FFIEC rule asks the institution to make.

The matrix also recognizes the difference between the alert and the determination — many institutions blur these because their IR programs were built before the FFIEC rule introduced the determination concept. The chain spec's matrix is built around the determination from the start, which means an institution adopting the spec inherits the right operational posture rather than having to retrofit it.

### 4.2 The CFR-citation routing matrix is a working artifact, not a survey

The "Federal-regulator routing per institution charter" subsection in the IR playbook is the single most useful artifact in this corpus from an operational standpoint. The matrix names six charter types, the primary federal regulator for each, the CFR citation, and the state-side cadence where applicable. National bank to OCC under 12 CFR Part 53; state member bank to FRB under 12 CFR Part 225 App. F; state non-member to FDIC under 12 CFR Part 304 Subpart C; FCU to NCUA under 12 CFR Part 748 App. B; FBO subsidiary routed by chartering structure.

I have spent multi-hour engagement segments helping incident-response teams figure this out during a live incident. Having it pre-rendered at the IR-playbook level — not buried in counsel's compliance-archive — is operationally significant. The phrase "the IR Commander confirms the routing per institution charter at the time of the incident; the institution's standing IR documentation pre-records the primary regulator and the secondary state-side counterpart where applicable" tells me the document author understood that this artifact has to be pre-rendered in the institution's standing documentation and re-confirmed at incident time, not assembled from scratch.

For a multi-charter holding company (the typical Category III/IV US bank with a national-bank subsidiary plus a state non-member subsidiary plus possibly a federal credit-union affiliate), the matrix gives the IR Commander the routing without requiring counsel intervention in the first hour. Counsel still owns the final determination; the matrix gets the operational posture started.

### 4.3 The cold-DR-key institution-side consumption procedure closes a gap most vendors leave open

In every supply-chain review I have run on an open-source dependency claiming SLSA Level 3, the project publishes a key-recovery procedure but does not document what the institution does with that procedure on a recurring basis. The institution's third-line auditor then asks the question I am trained to ask: "what is the institution-side evidence that the project's published cold-DR fallback is operationally exercised, not just documented?"

The supply-chain document's "Institution-side consumption of cold-DR-key dry-run attestation" subsection answers that question directly. The institution downloads the year's `KEY-DR-DRYRUN-{year}.asc` attestation, validates it against the cached cosign or GPG public key, archives it in the control-evidence repository for the year, and treats failure to obtain the attestation within 30 days of the project's annual dry-run window as IR Scenario 11 branch (c). The institution's `KEY-DR-DRYRUN-{year}.asc consumption log` is the operational evidence for CSF GV.SC-04.

This is a CSF GV.SC-04 control I would otherwise have to construct myself in a compliance gap-analysis report. The institution adopting the chain spec inherits the consumption procedure; the SOC engagement team samples it; the FFIEC cybersecurity examiner validates the consumption log against the archived attestations. The control is closed-loop without the institution's third-line audit needing to assemble it.

### 4.4 The regulator-fingerprint reception procedure is institutionally honest about responsibility

Adversary I in the threat model could have stopped at "verifier supply-chain discipline" and left the trust-anchor reception side as institution-defined-but-unspecified. Most threat models do exactly that, because the trust-anchor reception is institution-side responsibility and the spec author has limited authority over what the institution does.

The threat-model authors went further. They documented a five-step reception procedure with three operational events (`regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`), they named the reception-failure sub-variant that triggers IR Scenario 11 sub-variant, and they were explicit that "without a documented institution-side reception procedure, the institution risks (a) accepting a forged rotation notice, (b) failing to update the verifier and verifying against the old fingerprint when the new fingerprint is in force, or (c) under-instrumenting the trust-anchor lifecycle so the SOC team and examiner cannot sample-test the institution's reception discipline."

That is the kind of language a Cyber Risk practice writes after the institution has had an examiner challenge a missing trust-anchor reception procedure. Having it baked into the threat model means the next institution adopting the spec does not pay the lesson the first institution paid.

### 4.5 The dual-algorithm case-(e) regulator-overlap matrix is forward-looking

The IR playbook's "Dual-algorithm verifier-version timing-overlap (institution-vs-regulator)" subsection enumerates the three posture-disagreement scenarios when the institution adopts a post-quantum algorithm before the regulator's verifier supports it. The third row — institution reports FAIL case (e) (X-algorithm signature fails, Y-algorithm signature validates), regulator's X-only verifier reports PASS — is the genuinely concerning one, because the regulator's verifier is structurally blind to the case-(e) signal.

The matrix names this as "the load-bearing reason the institution coordinates verifier-version compatibility BEFORE the transition begins" and prescribes the institution proactively notifying the regulator and supplying the institution's case-(e) verifier output as supplementary evidence, with the regulator's IT examination program SHOULD-clause to upgrade within 6 months. That coordination cadence will not be perfect in practice, but having the matrix pre-render the disagreement scenarios means the institution's IR program knows what to communicate to the regulator before the disagreement surfaces in an examination.

The CIRCIA-only-triggers-without-FFIEC-trigger matrix lower in the same document is a similar artifact for a different boundary. Both are evidence the document author thought about the seams between control regimes, not just the regimes individually.

## 5. Observations

### OBS-1 — Cross-tenant scope discovery during investigation: vendor-side communication SLA is referenced but not normative (informational)

The IR playbook §"Cross-tenant scope discovery during investigation (vendor-hosted topology)" names that the vendor's IR team produces the alert and conducts the cross-tenant investigation, and that the vendor's contractual notification SLA (typically 4-12 hours from vendor's determination) is a separate clock from the institution's 36-hour FFIEC clock. The 4-12 hour range is the right ballpark for what I see in commercial vendor MSAs.

What is missing — and what I would expect to see in a SOC 2 Type II description backed by this playbook — is a normative statement on what the institution does when the vendor's notification arrives **at the boundary** of the institution's 36-hour clock. Two specific cases:

- The vendor notifies the institution at hour 30 of the vendor's own clock. The vendor's investigation has produced a cross-tenant compromise determination. The institution's clock starts at hour 0 of the institution's determination — but the institution's investigation must consume the vendor's findings, conduct any institution-specific evidence gathering, and reach the institution's own determination. If the institution's investigation takes 8 hours and the vendor's notification consumed all 30 of the vendor's contractual hours, the institution is at hour 38 from the vendor's original determination point — a regulator might ask why the chain of determinations took 38 hours when the FFIEC rule is 36.
- The vendor's IR team and the institution's IR team disagree on whether the vendor's findings rise to the institution's determination threshold. The vendor's findings might be "cross-tenant compromise observed" while the institution's evidence might show no impact on the institution's specific tenant. The institution needs a documented disposition for "vendor determined; institution did not."

A normative paragraph after the "Vendor-vs-institution responsibility split" subsection would address both. Suggested language for the document author to consider:

> When the vendor's notification arrives within 6 hours of the institution's 36-hour FFIEC clock window (i.e., the vendor's contractual SLA pushed close to the institution's regulatory window), the institution's IR Commander SHOULD treat the vendor's notification arrival as a presumptive determination start and run the institution's investigation in parallel with notification preparation, deferring final determination only if the parallel investigation produces evidence that the vendor's findings do not implicate the institution's specific tenant. The institution's working-paper records both the parallel-investigation timeline and the vendor's notification timestamp.

This is institutional, not technical. I rate it **informational** because the existing matrix gets the institution most of the way there; the proposed paragraph is for the institution that wants the boundary cases pre-rendered.

### OBS-2 — BYOC vendor-support telemetry redaction: testing cadence is named but not the failure-mode response (minor)

The BYOC deployment doc §"Common pitfalls" entry "HSM PIN leaks via vendor support telemetry" prescribes "strict redaction policy; the PIN is never in any log or trace; the redaction is tested." The redaction-is-tested clause is correct but understated. In a CC6.7 evaluation a Big Four engagement team will ask: tested how often, by whom, with what coverage, and what is the response if a redaction-policy regression is detected post-deployment.

The deployment doc could close this with a one-paragraph addition naming:

- redaction policy testing cadence (typical: per-release of the bank-controlled OTel Collector configuration plus quarterly sample-test against production telemetry stream),
- the response if a redaction regression is detected (typical: pause vendor-support telemetry egress, conduct retroactive sampling of the prior 30 days of telemetry against the corrected redaction policy, notify the vendor's security team if any sensitive data was exposed in the regression window),
- and the operational event the bank emits on regression detection (suggested name: `support_telemetry.redaction_regression_detected`, parallel to spec §10.2's other operational events).

The pattern is the same as spec §10.1 fingerprint reconciliation applied to a different control surface. The mechanism for adding it is documentation, not specification rework. I rate this **minor** because the existing language is defensible — a Big Four reviewer would accept "the redaction is tested" with the institution's standard testing cadence applied — but the paragraph above would close the question pre-emptively rather than at examination time.

### OBS-3 — Regulator-fingerprint rotation: cross-channel parallel notification authenticated-domain dependency is not enumerated (material)

The threat model §2.9 step 2 prescribes that the institution validates the rotation notice against the regulator's published GPG signature and "cross-checks against the regulator's parallel notification on the regulator's authenticated-domain channel." The reception procedure assumes the institution has standing access to the regulator's authenticated-domain channel for cross-channel verification.

For a national bank communicating with the OCC, the authenticated-domain channel is well-established (OCC examiner portal, BankNet, etc.). For a state credit union communicating with NCUA + a state credit-union regulator, the cross-channel verification posture is less uniform — the state regulator's authenticated-domain channel is per-state and varies in maturity. For an FBO subsidiary the cross-channel verification posture spans both the chartering federal regulator and the home-country supervisor.

What is missing is a statement on what the institution does when the cross-channel parallel notification is **not available** through an authenticated-domain channel of the regulator. Three specific dispositions the document could enumerate:

- The state-side regulator does not operate an authenticated-domain channel separate from the rotation-notice channel. Disposition: institution treats the cross-channel verification step as N/A for that regulator and documents the gap; the institution's reception-procedure SOC 2 evidence carries a CUEC (Complementary User Entity Control) on the state-side regulator's communication maturity.
- The cross-channel parallel notification is delayed (regulator's secondary channel has not yet published the parallel notification when the institution receives the primary). Disposition: institution does not install the fingerprint until the parallel notification arrives or until a documented bounded waiting window expires (suggested: 7 calendar days), at which point the institution falls back to out-of-band verification per Scenario 11 sub-variant.
- The parallel notifications **disagree** (the primary notice's fingerprint does not match the parallel notification's fingerprint). Disposition: this is the Scenario 11 sub-variant case 1 (forged notice suspected), but the threat-model document says the institution does NOT install the fingerprint and activates IR Scenario 11 sub-variant; it does not name the disposition for the disagreement specifically.

I rate this **material** because the cross-channel verification is the load-bearing step in the reception procedure — without it, the institution is validating the GPG signature against the regulator's published key, which is itself in the institution's cached anchor set, which closes the loop only if the cached anchor itself has not been substituted. Absent the cross-channel verification, the institution is one anchor-cache substitution away from accepting a forged notice. The threat model recognizes this implicitly (the cross-channel step is named) but does not document the dispositions for the cases above.

The fix is documentation, not architecture. I would expect the document author to add a paragraph after the §2.9 step 2 enumeration naming the three dispositions above, and an additional row in the IR playbook's Scenario 11 sub-variant for the disagreement-disposition case.

### OBS-4 — IR playbook Scenario 11 sub-variant: missing investigation-conclusion bounded window (minor)

The IR playbook's Scenario 11 sub-variant prescribes that case 1 (forged notice) STARTS the 36-hour clock at the institution's determination of forgery, and cases 2 (operational error at regulator) and 3 (institution-side validation tooling defect) do NOT start the clock. What is not named: the institution must commit to a determination across the three cases within a bounded window, parallel to Scenario 12 branch (iii)'s 48-72 hour bounded window for "under investigation."

Without a bounded window, an institution's IR program might let the determination drift indefinitely. The drift is operationally tempting because case 1 carries regulator notification and the institution's IR program may be motivated to determine "case 2 or 3" rather than commit to "case 1." A bounded window prevents the drift.

Suggested language:

> The institution MUST commit to a determination across the three cases within a bounded window (typically 48-72 hours from the failed validation event). The bounded window prevents the determination from drifting indefinitely; if the institution cannot commit within the window, the institution defaults to case 1 (forged notice suspected), starts the 36-hour clock at the bounded-window expiration, and continues investigation in parallel with notification preparation.

I rate this **minor** because Scenario 12 branch (iii) already documents the bounded-window pattern, so the institution's IR program has a precedent to follow even without explicit Scenario 11 sub-variant language. Adding it makes the playbook self-consistent.

## 6. Things I evaluated and found acceptable as-is

These are areas a Big Four engagement team would specifically test and where the documentation already meets the bar.

| Test | Documentation answer | Acceptable? |
|---|---|---|
| Multi-charter holding company FFIEC routing | CFR-citation matrix in IR playbook §36-hour matrix | Yes |
| FBO subsidiary routing | Last row of CFR-citation matrix | Yes |
| Vendor-hosted multi-tenant cross-tenant compromise discovery | "Cross-tenant scope discovery during investigation" subsection | Yes (subject to OBS-1) |
| Vendor MSA notification SLA vs FFIEC clock | Vendor-vs-institution responsibility split subsection | Yes (subject to OBS-1) |
| Cosign + GPG dual-trust-path independence | Supply-chain doc §"Dual-compromise edge case" | Yes |
| Cold-DR-key institution-side consumption evidence | Supply-chain doc §"Institutional consumption of cold-DR-key dry-run attestation" | Yes |
| FIPS 140-2 L3 disqualification list per cloud | Cloud HSM guide conformance matrix | Yes |
| Per-tenant key isolation enforcement | Cloud HSM guide §"Pitfalls" + spec §10.5 + design 04 §3.2 | Yes |
| HSM tamper-detection integration to IR | IR playbook §"HSM tamper-detection integration" | Yes |
| Long-dwell adversary residual treatment | IR playbook §"Long-dwell adversary considerations" | Yes |
| CIRCIA 72-hour overlay per scenario | IR playbook §"CIRCIA-only triggers without FFIEC trigger" | Yes |
| Dual-algorithm regulator-overlap matrix | IR playbook §"Dual-algorithm verifier-version timing-overlap" | Yes |
| Software-key fallback in production response | IR playbook Scenario 6 + spec §10.7 | Yes |
| IKM premature retirement detection and remediation | Spec §10.9 + IR playbook Scenario 8 | Yes |
| Audit-file truncation refusal | Spec §4.1 + IR playbook Scenario 9 | Yes |
| Constant-time comparison discipline | Spec §10.8 | Yes |
| 32-byte IKM minimum (RFC 4868) | Spec §10.6 | Yes |
| Append-only enforcement at two layers | Spec §10.3 | Yes |
| Operational-event catalog | Spec §10.2 | Yes |
| BYOC IAM permission matrix | BYOC doc §"IAM permission matrix" | Yes |
| BYOC mirror registry trust-path | BYOC doc + supply-chain doc §"Mirror-registry signature handling" | Yes |
| Vendor support telemetry egress control | BYOC doc step 5 + §"Common pitfalls" | Yes (subject to OBS-2) |
| SLSA Level 3 institution-side consumption | Supply-chain doc §"Institutional consumption of the SLSA L3 attestation" | Yes |
| SLSA-aware tool-choice documentation | Supply-chain doc §"SLSA-aware tool-choice documentation" | Yes |
| Mirror audit-log WORM retention | Supply-chain doc §"Mirror audit-log retention and integrity" | Yes |
| Mirror continuous-monitoring control (DE.CM-09) | Supply-chain doc §"Mirror continuous-monitoring control" | Yes |
| Co-signed seal failure (case (e)) triage | IR playbook Scenario 12 | Yes |

The breadth of this list is the headline finding. A typical Cyber Risk engagement on a draft control specification produces a 30-60 line gap list at this stage; this one produces three observations, two of which are minor or informational.

## 7. Cross-document consistency

I cross-checked the six documents the engagement letter named for inconsistency at boundaries where one document references another. The references resolve cleanly:

- IR playbook Scenario 11 sub-variant references threat-model §2.9 reception procedure → resolves; the reception procedure documents the validation step the sub-variant triggers on.
- IR playbook Scenario 6 references spec §10.7 software-key adapter compile-time exclusion → resolves; the spec normative requirement is the source-of-truth for what production builds MUST exclude.
- Supply-chain doc §"Cosign-key recovery during an active examination" references IR Scenario 3 → resolves; Scenario 3 is the binary-signature variant the supply-chain doc points to.
- Supply-chain doc §"Mirror continuous-monitoring control" references spec §10.1 fingerprint reconciliation as the pattern source → resolves; the mirror reconciliation is structurally parallel to the IKM reconciliation.
- Cloud HSM guide §"Mixing HSM tiers" pitfall is consistent with spec §10.5 HSM custody minimum.
- Threat-model §2.9 references `regulator-pack/regulator-procedures.md` as a v1.1 candidate stub. I did not read that file; the cross-reference is honest about its v1.1-stub status.

The cross-document consistency is at the level I would expect to see in a published SOC 2 Type II report's narrative section, where one description's claims are mirrored in the other description's evidence. This is unusual for draft specifications.

## 8. What I would tell the engagement partner

If I were debriefing the engagement partner ahead of the institution's CISO meeting, I would say this:

The chain-of-custody specification and its operational documents are at a maturity level where the institution can adopt them without restructuring the institution's existing IR program, SOC 2 control description, or vendor management framework. The institution inherits a clock-start triage matrix, a CFR-citation routing matrix, a regulator-fingerprint reception procedure, a cold-DR-key consumption procedure, and a dual-algorithm regulator-overlap matrix that the institution would otherwise have built itself over multiple examination cycles.

Three observations to raise with the institution: the vendor-clock-boundary case (OBS-1), the BYOC redaction-regression response (OBS-2), and the cross-channel verification dispositions for the regulator-fingerprint reception procedure (OBS-3). All three are documentation additions, not architectural changes. The institution can adopt the spec on the current document set and treat the three observations as items to raise with the project's working group for inclusion in the next document revision.

The fourth observation (OBS-4) is internal-consistency cleanup; the institution's IR Commander can apply the Scenario 12 branch (iii) bounded-window pattern to Scenario 11 sub-variant by analogy without waiting for the document revision.

The chain-of-custody specification is the kind of artifact a Cyber Risk practice recommends to its banking clients without reservation. I would be comfortable putting this opinion into the engagement-partner-signed cover letter.

---

## 9. Stopping criterion check

The engagement letter named 0/0 — I find the gaps I see. I have named four observations: one material (OBS-3), two minor (OBS-2, OBS-4), one informational (OBS-1). The material observation is documentation, not architecture. None block adoption.

I have nothing else to raise.

---

**Aleksei Sokolov**
Senior Manager, Cyber Risk Services
Toronto office
Round 15
