# Cybersecurity-specialist examiner review — chain-of-custody v1.0

**Reviewer.** Sofía Reyes, FFIEC Cybersecurity Specialist Examiner (FRB), 12 years on the cyber-specialist track. Lens: NIST CSF 2.0 mapping fidelity, IR playbook coverage of the chain's distinct failure modes, threat-model alignment with the shipped defensive primitives, supply-chain trust path, and the 36-hour computer-security incident notification rule.

**What I read.** The shipped artifacts: `spec/chain-of-custody-v1.md` (focus §10), `docs/design/00-overview.md` (focus §2), `docs/design/09-threat-model.md`, `docs/design/02-chain-construction.md` §4, `docs/regulator-pack/CSF-2.0.md`, `docs/incident-response-playbook.md`, `docs/cloud-hsm-guide.md`, `docs/supply-chain.md`, `docs/regulator-pack/ai-policy-alignment.md`, `docs/edge-and-federated-ai.md`. No prior round of this spec was loaded.

**Bottom line.** The spec is the most coherent attempt I have seen to take an AI-decision audit-trail problem and render it in primitives an FFIEC cyber-specialist can examine without a vendor in the room. The cryptographic floor is conservative and FIPS-anchored. The four-attack framing in `00-overview.md` §2 is clean and maps to CSF 2.0 categories with only one real gap. The supply-chain dual-signing posture (cosign + GPG) is defensible. The IR playbook covers the heavy hitters. The 36-hour articulation works for the headline cases and needs sharpening on three of the new chain-detected event types. Findings below, prioritised.

---

## 1. CSF 2.0 alignment of the four attack approaches

The four attacks in `00-overview.md` §2 sort cleanly into CSF 2.0 functions when I overlay them on the mapping in `docs/regulator-pack/CSF-2.0.md`:

| Attack (`00-overview.md` §2) | CSF 2.0 function | Subcategory landing |
|---|---|---|
| 2.1 Host-side chain tampering | PROTECT + DETECT | PR.DS-06 (integrity) primary; DE.CM-09 + DE.AE-02 (anomalous events surfaced by the verifier) |
| 2.2 Cross-tenant or cross-version key confusion | IDENTIFY + PROTECT | ID.AM-08 (asset/key lifecycle); PR.AA-05 (least privilege at IKM lookup); PR.DS-06 (the fingerprint check is integrity-of-derivation) |
| 2.3 Server-side or privileged-insider history rewrite | PROTECT + DETECT + RESPOND | PR.DS-06 (Merkle+HSM integrity); PR.PS-01 (HSM platform hardening); DE.AE-03 (correlation across operational events); RS.MA-01..05 (IR Scenarios 2 and 3) |
| 2.4 Examination-time tooling subversion | GOVERN + IDENTIFY | GV.SC-01..10 (cybersecurity supply chain); ID.RA-09 (authenticity of hardware and software) |

What I like, as a cyber-specialist:

- The headline mapping in `CSF-2.0.md` (PR.DS-06 as the chain's primary subcategory) is honest. The doc resists the temptation to claim every subcategory the chain "touches"; it only claims the ones the chain is the load-bearing control for. PR.DS-01, PR.AA at principal level, and DE.CM-01..08 are explicitly disclaimed. That is exactly the boundary an examiner needs the institution to draw.
- The four attack approaches are not just a security-engineering taxonomy — they are written in a way that lets me, in a finding letter, point a CISO at attack 2.3 and say "this is the PR.DS-06 evidence I am examining." The cross-walk holds.
- The IDENTIFY function is implicitly handled by attack 2.2 (cross-tenant key confusion is, in CSF terms, an asset-management failure of key generations). The new per-entry `key_fingerprint` is the right primitive for that subcategory; ID.AM-08's "throughout lifecycle" language is exactly what the per-entry fingerprint enforces.

What is missing or weak:

- **GV.OC (Organizational Context) and GV.RM (Risk Management Strategy).** The CSF mapping table only cites GV.RM-06 (risk responses) and GV.SC. CSF 2.0 made GOVERN a top-level function for a reason — examiners now expect to see how integrity controls flow from risk-tolerance statements. The mapping table does not point to a document that says "here is how the institution's risk tolerance for AI-decision integrity is articulated and how the chain implements it." The closest thing I found is the residual-risk register in `09-threat-model.md` §6, which is a project-side artifact, not an institution-side one. **Finding-grade gap:** the regulator pack should add a one-page "GV alignment" doc the institution copies into its policy framework, naming risk tolerance, ownership, and evidence cadence.
- **DE.CM-09 evidence binding.** The mapping cites DE.CM-09 ("Computing resources are monitored to find anomalous events") and points to operational events. The list of operational events in spec §10.2 is good, but the mapping doesn't name which specific operational events satisfy DE.CM-09 versus DE.AE-02 versus DE.AE-03. An examiner working a CSF assessment wants the evidence-to-subcategory binding at one cell of granularity finer than the table currently goes. **Recommendation:** annotate `CSF-2.0.md` so each operational event from spec §10.2 lands against the subcategory it most strongly supports.
- **RC.RP-04 ("Critical mission functions are restored").** The mapping points to `06-ledger-server-design.md` §7.5 for DR. I read the IR playbook §5 (sealing delay) and the cloud-HSM guide §"Failover and DR." Both are competent on the HSM-availability path. Neither walks through "the verifier returns FAIL on a recovered ledger because of a backup-integrity gap" — which `09-threat-model.md` §3.1.1 acknowledges as a real DoS vector. The IR playbook lists "Backup integrity failure" in scope but only Scenario 5 cleanly addresses it. **Recommendation:** add a Scenario 7 for backup-integrity verifier failure, with explicit RC.RP-04 + RC.CO-03 evidence callouts.

Verdict on CSF alignment: **the four attacks map cleanly to CSF 2.0; the mapping doc is honest about what it does not cover; the gaps are at the GOVERN function and at evidence-to-subcategory granularity.**

---

## 2. IR playbook coverage of the new failure modes

The rework introduces three failure modes that did not exist in earlier audit-chain shapes I have seen. I checked whether the IR playbook (`docs/incident-response-playbook.md`) addresses each one with operational specificity.

### 2.1 `key_fingerprint` mismatch

This is the failure mode that fires when the verifier looks up an IKM for `(tenant_id, key_version)` and `SHA-256(utf8(tenant_id) || ikm)[:16]` does not match the entry's stamped fingerprint. Per spec §7 step 8 and §4.1 inviolate property #3, this catches botched rotation, restored backup pointed at the wrong tenant row, and cross-tenant key swap **before any MAC compute happens**.

What the IR playbook does cover:

- Scenario 1 (chain hash mismatch at ingest) — but this is the wrong scenario. A `key_fingerprint` mismatch is a different failure with a different root-cause distribution. SDK defects produce hash mismatches; key-fingerprint mismatches almost always indicate operational misconfiguration at the IKM custodian, restored-backup confusion, or a botched rotation.
- Scenario 4 (master key compromise) — touches the right operational team but assumes confirmed compromise; a fingerprint mismatch is more often a misconfiguration than a compromise.

What is missing:

- **No dedicated scenario for `key_fingerprint` mismatch.** This is the failure mode the rework was designed to surface cleanly. The whole point of catching it before MAC compute is to give the SOC a precise alert that says "this is a key-management problem, not a tampering event." The IR playbook's scenario list reads as if the rework's defensive primitive doesn't exist. **Finding-grade gap.** Add Scenario 7 explicitly named "key_fingerprint mismatch detected at verification" with the right triage tree: was a rotation in flight? was a backup restored? has a tenant row been restored from a different period? Most of the time this is a P2 ops misconfig, not a P1 compromise — and the IR playbook should let the on-call differentiate.
- **No tie-in to the §10.1 reconciliation cadence.** Spec §10.1 makes weekly key-fingerprint reconciliation a SHOULD. The IR playbook should reference it directly: "if a fingerprint mismatch fires outside the weekly reconciliation, escalate per Scenario 7; if it fires during reconciliation, route to the reconciliation owner first."

### 2.2 `unknown_key_version`

Per spec §7 step 7, when the verifier resolves `(tenant_id, key_version) -> IKM` and the lookup returns null, the result is `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N`. The spec is clear: no MAC compute happens.

What the IR playbook does cover:

- Nothing directly. Scenario 3 (signature verification) covers seal-level key resolution failure. Scenario 4 (master compromise) addresses suspect IKMs. Neither addresses the case where the verifier looks for an IKM the registry never knew about.

What is missing:

- **No dedicated scenario.** This failure mode has three plausible root causes: (a) IKM was decommissioned and the registry forgot to retain it for the legally-required retention period; (b) IKM was provisioned at the SDK but never registered with the verifier-side registry; (c) tampered `key_version` field on the entry. The IR triage for each is different. The first is a retention-control failure (PR.DS-06 supporting evidence), the second is a provisioning-control failure (ID.AM-08), the third is tampering (RS.MA-05). **Add Scenario 8.**
- **The retention-period angle deserves a callout in `09-threat-model.md`.** R7 in the residual-risk register names "cross-tenant key reuse" but does not name "decommissioned IKM + chain entries that still reference it." For a 7-year audit retention, IKM-registry retention has to keep pace; if an institution decommissions an IKM at 90-day rotation and discards the IKM bytes, every entry stamped with that `key_version` becomes unverifiable. That is a known footgun the threat model should articulate.

### 2.3 Mid-write truncation

Per spec §4.1 (mid-write truncation refusal) and §7 (verifier MUST refuse a file whose last byte is not `\n`), this fires when a writer crashes mid-append.

What the IR playbook does cover:

- Nothing directly. The closest scenario is Scenario 1 (chain hash mismatch at ingest), which addresses bad data arriving over OTLP, not corruption of the persisted file.

What is missing:

- **No dedicated scenario.** Mid-write truncation is mostly an availability event (the writer crashed; the audit file is now in a state the verifier refuses); occasionally it indicates storage failure or, very rarely, deliberate truncation. The IR triage is straightforward but distinct: identify the writer process that crashed, recover the dropped event from the SDK's local SQLite buffer, replay if available, document the gap if not. **Add Scenario 9.**
- **Operational-event coverage.** The list in spec §10.2 includes `chain.verification_failure` but doesn't name a `audit_file.truncation_detected` event. If the verifier surfaces this failure mode, the corresponding operational event needs to fire so the SOC sees it before the next examination. **Recommendation:** add the event to spec §10.2.

### 2.4 What the IR playbook gets right

To be balanced about it: Scenarios 3 (signature verification failed) and 4 (master compromise suspected) are well-shaped, with the right escalation chain, the right notification targets, and concrete remediation steps including HSM rotation and re-signing. The 36-hour notification timing is referenced explicitly in both. The CIRCIA + 72-hour-CISA + state-AG decomposition in §"CIRCIA and additional notification frameworks" is the right multi-path articulation; counsel can work from this table.

The HSM tamper-detection integration section is also strong — it correctly treats a hardware tamper event as a preemptive Scenario 3 activation rather than waiting for a verifier failure.

Verdict on IR coverage: **the playbook covers the high-severity headline scenarios but the rework's three new chain-detected failure modes do not yet have dedicated scenarios. This is the largest single gap I found in the docs.**

---

## 3. Threat model fidelity to the rework's defensive primitives

I compared `09-threat-model.md` against the spec's actual primitives in §4.1 and the rework callout in §12 (change log).

What the threat model gets right:

- The eight adversary types (A–H) cover the realistic capability set. The composition argument in `00-overview.md` §2.3 — "HMAC defends against attackers without the key; Merkle+HSM defends against attackers with the key" — is articulated cleanly in adversaries B and C. The trust-zone diagram in §3 matches the deployment topology in `00-overview.md` §6.
- §3.1.1 (backup integrity), §3.1.2 (identity attribution), and §3.1.3 (forgery-allegation evidence path) draw the right boundaries between what the chain proves and what the institution's other controls have to prove. Forgery-allegation evidence path is a particularly good piece of writing — it is exactly what counsel needs when responding to a litigation hold.
- The R1–R8 residual-risk register names owners and mitigation status. R7 (cross-tenant key reuse) is correctly classified as low likelihood now that per-tenant HKDF binding is normative.

What is mis-aligned with the actual rework:

- **The rework's two headline defensive primitives are not first-class adversaries in the threat model.** Per-tenant HKDF binding (`info = info_base || "|" || utf8(tenant_id)`) and `expected_prev_hash` (not `entry.prev_hash`) in MAC recompute are the two changes that close the most subtle attack classes. Neither appears explicitly in §2's adversary catalog. Adversary B (insider with DB access) is the closest, but it does not articulate the `expected_prev_hash` defence — which is specifically a defence against a future-maintainer relaxation of the structural check. That is a real, named, codified design property; it deserves a residual-risk row of its own with the resolution status "closed by spec §4.1 inviolate property #8."
- **Adversary G (master key compromise) is still written in the pre-rework vocabulary.** It uses `session_key_id` reconciliation (the old field name) instead of `key_fingerprint` reconciliation (the new term). Spec §10.1 has been reframed around `key_fingerprint`; the threat model still says `session_key_id`. This is a doc-drift item but a real one — an examiner reading the threat model and the spec side by side will notice the term change and ask why. **Update §2.7 to use `key_fingerprint` terminology** consistent with spec §10.1.
- **§3.1.2 (identity attribution) still references `session_key_id binds to a process_uuid`.** Same drift. The shipped primitives bind via `key_fingerprint` and `key_version`; the doc still references the old shape.
- **R8 ("Session key leakage to logs") is medium likelihood, with mitigation "Session key never logged; key handles only."** The rework's per-entry `key_fingerprint` is publicly stamped on every entry. That is fine — the fingerprint is a public-by-design identity binding and the spec §10.6 IKM-length minimum closes the offline-grinding attack against a low-entropy IKM. But the threat model does not call out the offline-grinding bound explicitly; an examiner verifying the IKM-length control will want a residual-risk row that says "fingerprint is public; offline-grinding is closed by 32-byte IKM minimum (RFC 4868); residual risk is therefore acceptable." Add this.
- **Adversary H (cryptographic break) carries the 30-day spec-patch SLA and the 90/180-day migration windows.** Good. The PQC roadmap is reasonable for an FFIEC-aligned spec at this point in the NIST timeline. One thing missing: when the spec admits a post-quantum signature alongside Ed25519 (§4.3.2), the seal record's algorithm-identifier field is the dispatch key. The threat model does not articulate the transitional period during which both signatures are valid; an examiner thinking about a multi-year migration wants a residual-risk note for "dual-algorithm period: which algorithm does the verifier prefer when both are present, and what does that mean for the institution's evidence claims during the transition?"

What is genuinely missing from the threat model:

- **Edge / federated deployment threat additions.** `docs/edge-and-federated-ai.md` describes patterns that change the attacker capability surface materially — edge devices have weaker memory protection (called out in §"Operational considerations") and may be physically accessible. The threat model in §2 does not have a dedicated adversary for "physical attacker on an edge device with secure-enclave attestation." This is a v1.1 candidate at most, and the edge doc gestures at it correctly, but it should appear in the residual-risk register so the FFIEC examiner of an edge-deployed institution sees the residual risk explicitly.
- **Verifier-binary substitution at examination time.** Attack 2.4 in `00-overview.md` is articulated in the overview, but the threat model in `09-threat-model.md` §2 does not have a corresponding adversary I (examiner-side tooling). The supply-chain doc covers it, but a CSF-aligned threat model should have the adversary explicit in the central catalog. **Recommendation:** add Adversary I.

Verdict on threat model: **the model's bones are right and the trust-boundary discipline is exemplary. The doc has not fully absorbed the rework's vocabulary changes (`session_key_id` → `key_fingerprint`) and does not yet articulate the rework's two headline defensive primitives as first-class residual-risk rows. These are doc-hygiene items, not design gaps — but examiners will notice them.**

---

## 4. Supply-chain controls under the cyber-specialist lens

I read `docs/supply-chain.md` end-to-end and cross-checked against the spec's §4.4 (file-header attributes) and §10.7 (software-key adapter compile-time exclusion).

What I like:

- **Dual signing (cosign + GPG) with an explicitly-named fallback path.** This is the right shape. Sigstore is excellent but is not yet at the point where I would single-source the trust anchor for a regulated institution. The GPG-signed hash manifest as a fallback is exactly the kind of belt-and-suspenders an examiner can defend in an examination report. The trust-anchor table (cosign public key + GPG public key + toolchain version pinning) gives the institution a concrete checklist.
- **Reproducible builds with hermetic toolchain pinning.** `CGO_ENABLED=0`, `-trimpath`, `-buildid=` — these are the correct flags and the doc explains why each one matters for byte-for-byte reproducibility. The `reproducible-build-log` JSON record is a clean control-evidence artifact; SOC 2 examiners can sample-test it directly.
- **CycloneDX SBOM + Trivy/Grype scans published per release.** Standard at this point but worth confirming. Bank vulnerability-management programs consume CycloneDX directly.
- **The "what examiners verify" checklist at the end.** This is the most useful five-line section in the doc. It tells me, as an examiner, the exact five things to look for. The institution's evidence-collection cadence flows from it directly.
- **Software-key adapter compile-time exclusion (spec §10.7).** This is the right enforcement level. Run-time gating is a misconfiguration footgun; compile-time exclusion is what an FFIEC examiner can verify with a single binary inspection. The doc says it; the spec mandates it; the IR playbook Scenario 6 catches the runtime case if it ever escapes the build pipeline. That is defense in depth I can sign off on.

Where the supply-chain controls feel thin under the cyber-specialist lens:

- **Cosign key custody is named ("hardware-backed, held by the project's release-management role") but the recovery procedure is not.** What happens if the cosign key is compromised? The doc names the GPG fallback as the path institutions take if Sigstore is compromised, but doesn't articulate the project-side response: how is the key revoked, how are institutions notified, what is the new key's authentication chain? **Finding-grade gap for the project-side governance.** GOVERNANCE.md is referenced but I cannot evaluate it from this doc alone.
- **GPG key rotation cadence.** The doc says GPG is offline-held for emergency-response. It does not name a rotation cadence. Offline keys decay differently from online keys, but they still need a rotation plan; without one, an institution cannot articulate to its own auditor when the cached anchor needs to be refreshed.
- **No evidence the conformance corpus itself is signed.** Spec §8 references `spec/test-vectors/` and the doc names "Independent re-implementations from the spec text rather than the corpus" as a defense against corpus subversion. That's correct as a defense, but the corpus artifacts themselves are not listed in the per-release artifacts table at the top of the supply-chain doc. If an institution is going to rebuild the corpus from the spec text (the recommended pattern), the spec text needs an integrity binding too. **Recommendation:** add the spec PDF + a content-hash to the per-release artifact table; sign both.
- **Container supply chain — the institution mirror requirement.** The doc says "the institution does not pull from the project's registry directly." Good. But the doc does not name what the institution's mirror does in terms of signature re-verification. If the mirror just copies bytes, the cosign signature is preserved. If the mirror re-signs (some institutional registries do), the signature chain is broken and the institution has to bridge the trust path. This needs to be called out.
- **"Conformance corpus integrity" subsection.** The doc names multi-maintainer review and public git history. These are procedural; they are necessary but not sufficient. A supply-chain attacker who compromises a maintainer account can submit a covertly-malicious corpus entry. The defense the doc names — "external security review" — is good but unscheduled in the public artifact set. **Recommendation:** publish the external review cadence and the most recent review date in the spec text (or in `GOVERNANCE.md` referenced from the spec).
- **No SLSA level claim.** The build pipeline as described meets SLSA Level 3 properties (hermetic, reproducible, signed provenance). The doc does not claim a SLSA level explicitly. Cyber-specialist examiners increasingly look for SLSA-level claims because they are a single number that maps to a known set of properties. **Recommendation:** claim the SLSA level in the supply-chain doc and back it with the per-release attestation.

Verdict on supply chain: **the architecture is sound and the dual-signing + reproducible-build floor is at the level I expect for an examination-grade artifact. The gaps are mostly project-side governance items (key recovery, rotation cadence, SLSA-level claim) and one institutional-side item (mirror-registry signature handling) that should be lifted into the cyber-specialist's evidence checklist.**

---

## 5. The 36-hour cyber-incident notification rule articulation

The FFIEC computer-security incident notification rule requires a banking organization to notify its primary federal regulator within 36 hours of determining that a "computer-security incident" has occurred that has "materially disrupted or degraded, or is reasonably likely to materially disrupt or degrade" the bank's operations. The threshold is determination-of-incident, not detection-of-anomaly. This distinction matters for chain-detected events because the chain detects a wide range of anomalies, only some of which constitute notification-grade incidents.

What the docs get right:

- The IR playbook Scenario 3 (signature verification failed) explicitly cites the 36-hour rule and routes notification to the primary federal regulator. Good — signature verification failure is at or near the threshold for "reasonably likely to materially disrupt" the integrity of the bank's AI-decision evidence.
- IR Scenario 4 (master key compromise, suspected) cites the 36-hour rule. Master-key compromise is unambiguously notification-grade.
- The IR playbook §"CIRCIA and additional notification frameworks" decomposes the path matrix correctly: FFIEC 36-hour, CIRCIA 72-hour-CISA, state AG, CFPB UDAAP, sealing-delay-72h. Counsel can work from this table directly.
- Spec §4.3.1 (HSM unavailability) correctly distinguishes "operational delay → 72-hour SHOULD" from "associated with suspected security incident → 36-hour rule." This is the right separation; the 72-hour seal-delay notification is a different control surface from the cyber-incident notification rule.

Where the articulation needs to be more concrete for the new chain-detected event types:

- **`key_fingerprint` mismatch.** Is this a 36-hour-rule event? It depends on root cause. A botched rotation (operations error) is not. A confirmed cross-tenant key swap that the institution can attribute to malicious action is. The IR playbook does not give the on-call CISO a decision tree. **Recommendation:** add to the proposed Scenario 7 a triage matrix: "if the mismatch is attributable to a documented operations event (rotation in flight, restored backup) within X hours of the alert, route per change-management; if not, treat as suspected unauthorized key substitution and start the 36-hour clock."
- **`unknown_key_version`.** The triage is similar but the threshold question is different. If an IKM was decommissioned per retention policy and the chain entries that reference it are now unverifiable, that is a control failure — not a security incident. If an IKM appears that was never registered, that is suspicious and likely 36-hour-grade. The proposed Scenario 8 needs the same triage matrix.
- **Mid-write truncation.** Almost never 36-hour-grade. The default disposition is "writer crashed, recover dropped event from SDK buffer, document the gap." If the truncation pattern recurs across hosts in a way that suggests deliberate corruption, escalate. This belongs in Scenario 9 with explicit "do not start the 36-hour clock unless [specific indicators]" guidance.
- **The 36-hour clock-start question generally.** The rule says the clock starts at "determination" of the incident, not at "detection of the alert." The chain produces precise alerts; determination is a human judgment downstream. The IR playbook should articulate, for each chain-detected scenario, what evidence triggers the determination — so the clock-start moment is auditable. **Recommendation:** add a "clock-start trigger" line to each scenario in the IR playbook; counsel will use it when the bank explains its notification timing to the FRB.

What the spec itself gets right and could amplify:

- Spec §4.3.1 paragraph "Cyber-incident notification under the FFIEC computer-security incident notification rule (36 hours) applies separately when the HSM unavailability is associated with a suspected security incident. The two notification paths are distinct." This is the correct articulation. **Recommendation:** add a parallel paragraph to spec §10 (operational requirements) covering the chain-detected event types — `key_fingerprint` mismatch, `unknown_key_version`, mid-write truncation, signature verification failure, Merkle root mismatch — naming for each whether the default disposition is an operational SHOULD-notify or a 36-hour MUST.
- The `dev_mode` field on the seal record (spec §4.2 schema) and the `kms_handle_uri = "plaintext-dev"` marker are useful audit-trail evidence for the institution to demonstrate that the 36-hour clock did not start because the alert was a known-development environment leak. The IR playbook Scenario 6 (software-key fallback in production) cites the institution's "standard control-failure framework" for notification — that is the right framing because a dev-key in production is a control failure, not a cyber-security incident. The doc could be slightly more explicit that this is a deliberately-not-36-hour disposition.

Verdict on 36-hour articulation: **the headline cases (signature verification failure, master compromise) are correctly cited and routed. The new chain-detected event types do not yet have explicit 36-hour-applicability guidance, and the clock-start trigger is implicit. Both are addressable in the IR playbook update.**

---

## 6. Cross-cutting observations

A few notes that span the documents.

### 6.1 The trust-zone story holds at examination

The TZ1 → TZ2 → TZ3 trust topology in `09-threat-model.md` §3 is the cleanest way I have seen to explain to a cyber-specialist colleague why a vendor-hosted ledger is examinable to the same standard as a self-hosted one. The verifier reads TZ2 as untrusted input; the regulator-held public-key fingerprint binds the integrity claim. When I prepare an FFIEC examination of a vendor-hosted institution, this diagram is what I would use to scope the examination boundary. It does the work.

### 6.2 The cloud-HSM matrix is examiner-grade

The matrix in `cloud-hsm-guide.md` of "FIPS 140-2 L3 conformant?" by service is exactly what an institution needs to put in its control description and exactly what I need to verify. The explicit non-conformant entries (AWS KMS default, Google Cloud KMS software keys, Azure Key Vault Standard) are the most useful part — examiners frequently see institutions citing AWS KMS or Azure Key Vault Standard as their "HSM" and the matrix gives the bank's compliance team a clean reference to push back with.

The cost guidance at the end of each provider section is also useful for proportionality conversations. A community bank claiming HSM cost is prohibitive can be pointed at the Google Cloud HSM pay-per-use option and the §6.5 per-tier proportionality discussion in `00-overview.md`.

### 6.3 The AI-policy alignment doc is broad but uneven

`docs/regulator-pack/ai-policy-alignment.md` covers a lot of ground — NIST AI RMF, EO 14110, Treasury, CFPB, EU AI Act, EBA, DORA, NIS2, FCA/PRA, MAS Veritas, JFSA, BCB, U.S. state laws. The breadth is impressive and the EU AI Act Article 12 mapping is the strongest individual cell. The CFPB section is appropriately limited (the chain provides substrate; consumer-facing requirements are the bank's compliance program). DORA's 24-hour reporting is correctly noted.

What I would add: an explicit cross-reference to the IR playbook's notification matrix, so a cyber-specialist working a multi-jurisdictional institution can see the chain-detected event → notification-path mapping for non-U.S. frameworks. The DORA 24-hour clock and the NIS2 24-hour early warning are specifically tighter than the FFIEC 36-hour rule; institutions operating in EU jurisdictions need to know the chain's operational-event timing supports those tighter clocks.

### 6.4 Edge and federated AI — the doc is correct that v1.0 supports the patterns without modification

I evaluated `docs/edge-and-federated-ai.md` against the spec and the threat model. The claim that the chain handles edge / federated / on-device deployments without spec modification is defensible: the SDK's local-persistence + asynchronous-export pattern works in disconnected scenarios; the per-tenant HKDF binding and `key_fingerprint` work for per-device session keys; the DAG semantics in spec §4.4 cover multi-agent coordination.

The two patterns articulated for master-key custody at edge (Pattern A: per-device master in TPM/secure-enclave; Pattern B: bulk session-key issuance at commissioning) are pragmatic. Pattern A is preferred from a cyber-specialist lens; the doc says so. The reconciliation cadence on edge fleets (the "baseline the unmatched count for the edge fleet separately" guidance) is good operational pragmatism that examiners should expect to see in the institution's control description.

The threat model gap I named in §3 above (no dedicated adversary for "physical attacker on edge device with secure enclave") is this doc's main outstanding need.

### 6.5 What I would tell my colleagues

If a colleague on the cyber-specialist track asks me whether to use this spec as a reference point for an FFIEC examination of a bank's AI-decision audit trail, I would say yes, with three caveats:

1. The IR playbook update for the three new chain-detected event types is necessary before I would treat the doc set as complete for cybersecurity-examiner use. Without dedicated scenarios for `key_fingerprint` mismatch, `unknown_key_version`, and mid-write truncation, the SOC's response to those alerts is improvised, which is exactly what an examiner does not want to see.
2. The threat model needs the vocabulary refresh and the two new residual-risk rows (HKDF binding and `expected_prev_hash` defence). These are doc-hygiene items but they are visible to anyone reading the spec and the threat model side by side.
3. The supply-chain doc needs the cosign-key recovery procedure and the SLSA-level claim. Without those, the institution's third-party-risk-management has a gap.

None of these are design gaps. The design is sound. They are documentation completeness items that translate the design into examination-grade artifacts.

---

## 7. Findings summary

| # | Finding | Severity | Location | Recommended owner |
|---|---|---|---|---|
| F1 | No dedicated IR scenario for `key_fingerprint` mismatch (the rework's most distinctive failure mode) | High | `incident-response-playbook.md` §"Common scenarios" | IR playbook author |
| F2 | No dedicated IR scenario for `unknown_key_version` (decommissioned IKM and provisioning-gap cases) | High | `incident-response-playbook.md` §"Common scenarios" | IR playbook author |
| F3 | No dedicated IR scenario for mid-write truncation; no operational event named for it in spec §10.2 | Medium | `incident-response-playbook.md` + spec §10.2 | IR + spec authors |
| F4 | Threat model uses `session_key_id` vocabulary; spec §10.1 uses `key_fingerprint` (doc drift) | Medium | `09-threat-model.md` §2.7 and §3.1.2 | Threat-model author |
| F5 | Per-tenant HKDF binding and `expected_prev_hash` are not first-class residual-risk rows | Medium | `09-threat-model.md` §6 | Threat-model author |
| F6 | No GV.OC / GV.RM evidence binding for institution-side risk-tolerance articulation | Medium | `regulator-pack/CSF-2.0.md` | Regulator-pack author |
| F7 | Operational events from spec §10.2 not bound to specific CSF subcategories at one cell finer granularity | Low-Medium | `regulator-pack/CSF-2.0.md` | Regulator-pack author |
| F8 | No 36-hour-applicability triage matrix per chain-detected event type; clock-start trigger implicit | High | `incident-response-playbook.md` (per scenario) | IR playbook author + counsel |
| F9 | Cosign-key recovery procedure not articulated in supply-chain doc | Medium | `supply-chain.md` | Project governance + supply-chain author |
| F10 | GPG-key rotation cadence not named | Low-Medium | `supply-chain.md` | Project governance |
| F11 | No SLSA-level claim despite the build pipeline meeting SLSA L3 properties | Low | `supply-chain.md` | Supply-chain author |
| F12 | Spec text + conformance corpus integrity binding not in per-release artifact table | Medium | `supply-chain.md` | Supply-chain author |
| F13 | Container mirror-registry signature handling not articulated | Low-Medium | `supply-chain.md` | Supply-chain author |
| F14 | No threat-model adversary for examiner-side tooling (attack 2.4 in overview is not in §2 catalog) | Low | `09-threat-model.md` §2 | Threat-model author |
| F15 | No threat-model adversary for physical attacker on edge device with secure-enclave | Low | `09-threat-model.md` §6 + `edge-and-federated-ai.md` | Threat-model author |
| F16 | RC.RP-04 evidence does not include a backup-integrity verifier-failure scenario | Low-Medium | `incident-response-playbook.md` + `regulator-pack/CSF-2.0.md` | IR + regulator-pack authors |
| F17 | IKM-registry retention obligation (must outlive any chain entry that references it) is not explicit | Medium | `09-threat-model.md` and spec §10 | Spec + threat-model authors |
| F18 | Dual-algorithm transitional-period guidance (post-quantum coexistence) absent | Low | `09-threat-model.md` §2.8 | Threat-model author |

The High-severity findings (F1, F2, F8) are the ones that would be in my preliminary findings letter to the institution if I were examining a deployment of this spec today. The Medium-severity findings would land as MRA (Matter Requiring Attention). The Low-severity findings are doc-quality items.

---

## 8. Closing assessment

The cryptographic floor is right. The trust-zone topology is right. The verifier procedure is correctly defense-in-depth-ordered. The supply-chain dual-signing posture and the FIPS-140-2-L3 HSM requirement are at the level a cyber-specialist examination expects. The CSF 2.0 mapping is honest about what the chain is and is not the headline control for. The IR playbook covers the heavy hitters. The 36-hour notification rule is correctly cited where it most clearly applies.

The work that remains is mostly documentation completeness rather than design change. The rework introduced specific defensive primitives — per-tenant HKDF binding, per-entry `key_fingerprint` checked before MAC compute, `expected_prev_hash` in MAC recompute, mid-write truncation refusal — and the IR playbook and the threat model have not yet caught up to those primitives in their failure-mode vocabulary. Closing the F1–F8 findings would bring the operational-side documentation to the same standard as the spec.

I would be comfortable using this spec as an examination reference for an FFIEC cybersecurity assessment of a regulated institution's AI-decision audit trail, with the IR playbook updates as a precondition for examination-grade institutional adoption.

---

*Reviewed by Sofía Reyes, FFIEC Cybersecurity Specialist Examiner, FRB. No prior round of this spec was loaded for this review.*
