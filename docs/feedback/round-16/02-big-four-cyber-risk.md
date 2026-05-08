# Round 16 — Big Four Cyber Risk review

**Reviewer.** Marta Kowalczyk, Senior Manager, Cyber Risk Advisory (Warsaw office). 18 years across European bank cybersecurity advisory; partner-track on third-party-risk and cyber-supply-chain engagements for tier-1 and tier-2 EU institutions. CISA, CISM, CRISC.

**Engagement context.** First-look review of the v1.0 chain-of-custody artifacts. I have read only the documents the engagement letter named: spec §10, design 00, design 04 §3.2, design 09 §2.9, `byoc-deployment.md`, `supply-chain.md`, the full IR taxonomy 1–12 in `incident-response-playbook.md`, and `cloud-hsm-guide.md`. I have not read prior-round feedback or design history. My findings are framed against what an EU institution's third-party-risk function would put in front of a Cyber Risk partner before the institution signs an examination-grade adoption.

**Stopping criterion.** Engagement letter says 0/0 — I report all findings I see, no deferrals. I am not banking findings for a follow-up round.

---

## 1. Headline disposition

The artifacts, taken as a set, present a chain-of-custody design that is unusually well-shaped against supply-chain and trust-anchor risk. The points my partner will care about are positive ones: the design knows the line between what an institution can defend and what is project-side or regulator-side residual, and it is honest about that line in writing. That is not common in this category of work — most vendor specs we review smear institution-side and platform-side controls together until the SOC team cannot tell where to point a sample-test.

The discipline that earns the headline:

- Three independent trust anchors (cosign, GPG, regulator-held public-key fingerprint) for the verifier binary, with the dual-compromise edge case named and a cold-DR key lifecycle that the institution actually consumes via an annual dry-run attestation. The institution-side consumption procedure is what closes the loop — without it, the cold-DR fallback would be documented-but-unexercised, which is the exact gap a GV.SC-04 examiner will challenge. Supply-chain doc names that explicitly.
- BYOC IAM matrix that says the vendor MUST NOT reach the HSM, master key, or event-payload data, and gives the institution a unilateral revocation path. The vendor's support telemetry runs through a bank-controlled OTel Collector that enforces redaction. That is how this control should be drawn.
- IR playbook's clock-start matrix per scenario. The matrix separates the alert-point from the determination-point and keeps the 36-hour FFIEC clock from being started by every operational mismatch. CIRCIA path is tracked in parallel, not collapsed into the FFIEC path. Counsel reads this and does not have to reverse-engineer the chain's signal taxonomy.
- The `09-threat-model.md` §2.9 reception procedure for a regulator-held fingerprint rotation, with three named operational events providing the audit-evidence path. The reception-failure sub-variant (forged notice) is given its own IR scenario (Scenario 11 sub-variant), which is the right place to put it — a forged trust-anchor-rotation notice is a 36-hour-clock event the moment forgery is determined, and the playbook says exactly that.

Below I list the findings my engagement would put in writing. Some are observations rather than gaps — I include them because a partner-review will ask whether I considered them.

---

## 2. Findings — what I would write in the engagement memo

### 2.1 Cyber-supply-chain (CSF GV.SC, ID.RA-09) — what to commend

**The SLSA L3 institutional consumption clause is load-bearing and correctly positioned.** Most vendors ship a SLSA provenance attestation and stop. This spec writes the "publish is necessary but not sufficient — institution MUST consume" clause directly into the supply-chain doc, names `slsa-verifier` as a co-equal trust artifact alongside cosign and GPG, requires the institution to cache the `slsa-verifier` public key and binary SHA-256, and requires the institution's control description to name which SLSA-aware tool it operates with what version and what trust anchor. A GV.SC-04 examiner asks four questions: which tool, what version, what trust anchor, archived where. The institution's documentation has those four answers in one place. That is the right shape.

**The re-signing-mirror bridge invariant has retention and integrity attached.** The supply-chain doc requires the mirror's audit log to be on WORM storage with 7-year retention, plus a daily reconciliation between pulled-image set and audit-log-validation set. The bridge invariant ("the mirror's re-signature is conditional on the project's cosign signature having validated at mirror-pull time") rides on the audit log; if the log is rotation-deletable, the invariant cannot be examined retroactively. The retention-and-integrity clause closes that. The reconciliation cadence is the same control pattern as spec §10.1 fingerprint reconciliation — the institution's SOC team sees one pattern, applied twice, instead of two ad-hoc designs.

**The cold-DR key dry-run cadence is consumed institution-side.** The annual `KEY-DR-DRYRUN-{year}.asc` attestation has an explicit institution-side consumption procedure (validate, archive, retain per chain-event retention, fire IR Scenario 11 branch (c) on missed window). The CSF GV.SC-04 binding is named. This is the part most vendor designs omit — they document the project's dry-run and stop. The institution reads it and goes "what am I supposed to do with this?" Here, the institution knows.

**Observation, not a finding.** The supply-chain doc treats the GPG fallback's rotation cadence (36 months, 12-month subkey rotation) as load-bearing. That is the correct posture for offline keys per NIST SP 800-57; many vendor designs treat offline keys as rotation-exempt. The 36-month cadence aligned with NIST SP 800-57 is what I would advise an EU institution to require as a contractual term against any AI-decision-evidence vendor. The spec saves us that conversation.

### 2.2 Cyber-supply-chain — what I would still ask for

**Finding 2.2.1 — Reproducible-build evidence threshold across institutions.** The supply-chain doc's "operational pattern" step 4 says "at least one institution per release; ideally many" performs reproducible-build verification. As an institution-side reviewer evaluating my client's third-party-risk posture, I cannot verify this as a property of *the spec's adoption ecosystem* without a published count or registry of which institutions performed the rebuild for any given release. If five releases have shipped and only one institution rebuilds release 3 and zero rebuild releases 1, 2, 4, 5 — the "ideally many" property is unmet and the institution my engagement is advising is not seeing a published signal of how exercised the rebuild discipline is.

**Recommendation.** The project's release process publishes an aggregated count (anonymized if individual institutions don't want to be named) of independently-verified reproducible-build attestations per release, on the same release page as the binary. The institution's release-validation procedure consumes this count as a defense-depth signal: a release with zero independent rebuilds 30 days post-publication is a different risk posture than one with ten. Publishing the count without naming institutions preserves anonymity; the absolute number is the signal. This is roughly the same shape as the cold-DR-key dry-run attestation pattern — the project publishes operational evidence the institution consumes for its own posture.

**Severity.** Medium. The institution can rebuild itself (and per the supply-chain doc it is supposed to), so the institution's own posture is not impaired. The finding is about the spec's adoption-ecosystem signal — the institution wants to know whether independent reproducibility is being exercised across the ecosystem, not just by itself. A GV.SC-04 examiner asks whether the institution is alone in its rebuild discipline; the count answers that.

**Finding 2.2.2 — Conformance corpus rebuild — independent re-implementation provenance.** The threat model §5 names independent re-implementation from the spec text as a defense, and the supply-chain doc names "rebuilding the corpus from the spec text and confirming against `spec-manifest.sha256.asc`" as the recommended deeper verification. As written, the institution's only check is whether *its own* corpus rebuild matches the manifest. There is no mechanism for the institution to know whether *another implementation*'s rebuild also matches.

**Recommendation.** The project maintains a public registry of independent re-implementations that have rebuilt the corpus from the spec text and verified against the manifest. Each entry: implementation name, spec version, manifest hash matched, attesting party (independent of the project), date. The institution's third-party-risk-management consumes this registry as a signal of the spec text's interpretation-stability across implementations. If only the reference implementation has ever produced a corpus that matches the manifest, the spec text is not in fact independently re-implementable — it is reference-implementation-dependent in disguise.

**Severity.** Medium. Same shape as 2.2.1 — institution can defend its own posture via its own rebuild, but the ecosystem-signal is not visible.

### 2.3 BYOC topology — what to commend

**The vendor-cannot-reach-HSM control is tested at IAM, network, and process layers simultaneously.** The IAM matrix says no on the HSM signing operations row for the vendor's support role. The network-controls section places the HSM in a private subnet with private endpoints, no internet egress, security-group inbound only from the runtime workload. The deployment steps require the vendor's image to assume a bank-controlled IAM identity that the bank can revoke unilaterally. Three control layers, each enforcing the same property. A SOC 2 CC6.1 evaluator gets to sample-test all three independently and they all converge on the same answer.

**Image-pull integrity is bank-controlled.** The mirror-registry pattern with cosign signature validation at pull time is the right shape. The bank does not pull from the vendor's registry directly. If the vendor's registry credentials are compromised, the bank's pulls still validate signatures, and the bank's deployment pipeline pulls only from the bank's mirror. The compromise window collapses to a signature-substitution attack the cosign trust path is designed for.

### 2.4 BYOC topology — what I would still ask for

**Finding 2.4.1 — Vendor-support telemetry's redaction policy is referenced but not normative.** The BYOC doc's deployment step 5 says "the collector applies the bank's documented redaction policy (no event payload data; control-plane events only)." The redaction policy is bank-side and bank-defined. There is no published baseline redaction policy for the vendor's support telemetry — a tier-1 bank with a mature data-protection function will produce a strong policy, but a tier-2 community-bank-tier institution adopting a vendor-hosted shared-cloud deployment may underspec the policy and leak event payload data the spec assumes is bank-confidential.

**Recommendation.** The vendor (per its conformant deployment guidance) publishes a baseline redaction policy as part of the deployment artifacts. Institutions adopt the baseline as a starting point and tighten where their data-protection program requires. The baseline names the specific OTel attributes the collector MUST drop (e.g., `payload`, `prompt_text`, `response_text`, `gen_ai.*` content fields), and the specific ones it MAY pass through (e.g., `seq`, `run_id`, `tenant_id`, operational events from spec §10.2). This is a baseline-tightening posture, not a baseline-loosening posture, so the institution's autonomy on its own redaction policy is preserved.

**Severity.** Medium for tier-2 institutions, low for tier-1. A tier-1 bank's data-protection officer will write the policy from scratch. A community-bank-tier institution adopting BYOC at the cost-floor will use whatever the vendor recommends, and "consult your data-protection officer" is not adequate guidance for that tier.

**Finding 2.4.2 — Vendor-image upgrade gating across IAM-permission drift.** The BYOC doc's "When the vendor relationship ends" section names the migration steps. It does not name the corresponding upgrade-gating control: when the vendor publishes a new image version, the institution validates whether the new image's IAM permission requirements are unchanged versus the deployed version. A vendor that adds new functionality may require expanded IAM permissions on `chain-ledger-runtime` (e.g., new HSM operation, new S3 prefix, new secret-manager path). The institution's IT change-management may not catch this if the vendor's release notes do not flag the IAM delta explicitly.

**Recommendation.** The vendor's release pipeline publishes a per-version IAM-requirements manifest (the specific permissions the image requires, in JSON or HCL form) alongside the binary. The institution's deployment pipeline diffs the new manifest against the deployed manifest before applying. A non-empty diff is a change-management gating event; the institution's change-management ticket must approve the new permission set before the new image deploys. This is the same pattern as the SLSA L3 attestation deployment-gating control, applied to IAM.

**Severity.** Medium. A vendor expanding its image's IAM footprint silently is the kind of supply-chain drift a GV.SC examiner will challenge. Closing it with a published manifest the deployment pipeline diffs is low-cost and high-value.

### 2.5 IR playbook — what to commend

**Clock-start triage tree for Scenario 7 (key_fingerprint mismatch) is the playbook's most useful single page.** Three of the four triage cases do not start the 36-hour clock (rotation-in-flight, restored-backup, tenant-row-restored). Only the fourth — suspected unauthorized substitution — does. This is the correct posture. A naive playbook would treat every fingerprint mismatch as a 36-hour-clock event, which would cause the institution to over-notify and burn regulator goodwill on operational mismatches. Counsel reads the triage tree and routes per case. Same shape on Scenario 8 (unknown key_version).

**The constellation-rollup rule for concurrent multi-scenario alerts.** A Scenario 1 + Scenario 7 + Scenario 8 alert in a 60-minute window is one incident, not three, and the IR Commander documents the rollup decision. This is the rule that prevents the institution from accidentally over-notifying because three detection paths fired on one root-cause event. EU AI Act and DORA reporting framings will increasingly impose this kind of incident-correlation discipline; the spec is ahead of where most vendor playbooks are.

**Federal-regulator routing per institution charter.** The playbook pre-documents which regulator is the primary path per charter type with CFR citations. State-side cadences are also tracked. A multi-charter holding company knows it notifies all applicable primary regulators. This is the kind of pre-documentation that takes most institutions a separate engagement to produce; the spec ships it.

**The dual-algorithm verifier-version timing-overlap matrix (institution-vs-regulator) is correct and unusually mature.** Three posture-disagreement scenarios, each with its disposition. The case-(e) interaction with a regulator's X-only verifier is named as a load-bearing regulator-side blind spot during the dual-algorithm transitional period. The institution's pre-coordination obligation with the regulator is named. This is post-quantum migration governance the FFIEC has not yet published guidance on; the spec is ahead.

### 2.6 IR playbook — what I would still ask for

**Finding 2.6.1 — Concurrent multi-tenant scope discovery in vendor-hosted topology — institution-side awareness.** The "Cross-tenant scope discovery during investigation" edge case names the situation where Scenario 4 or 7 spans more than one tenant in a multi-tenant deployment. The vendor's IR team produces the alert; each institution's IR team receives the alert through the vendor's notification channel. The vendor's contractual notification SLA (4–12 hours from vendor's determination) is a separate clock from the institution's 36-hour FFIEC clock.

The gap: the institution does not know, from the alert it receives, which other tenants are in scope. It cannot conduct its own cross-correlation with cyber-threat intelligence feeds because it does not know whether the same compromise pattern affected its peer institutions on the same vendor platform. The vendor cannot disclose other tenants' identities for confidentiality reasons. This produces a governance gap: the institution's IR program operates blind on whether the suspected compromise is institution-specific or vendor-platform-wide, even though the answer materially affects the institution's incident-response posture (a vendor-platform-wide compromise indicates a vendor-side root cause and a different remediation than an institution-specific compromise).

**Recommendation.** The vendor's contractual obligations to its institution clients include an anonymized cross-tenant-scope indicator: when the vendor's IR team determines a Scenario 4 or 7 incident spans multiple tenants, the vendor's notification to each affected institution includes the count (anonymized) of other affected tenants and the broad pattern (e.g., "5 tenants affected, all using shared-cloud-HSM tier per `00-overview.md` §6.5"). The institution then knows whether to treat as institution-specific or vendor-platform-wide without learning peer institutions' identities. The vendor's incident-disclosure policy can be reviewed by each institution's third-party-risk-management at vendor onboarding.

**Severity.** High for vendor-hosted topology adopters; not applicable for self-hosted or BYOC. The vendor-hosted topology is named in `00-overview.md` §6.3 and §6.5 as the cost-appropriate deployment for community-bank-tier institutions; precisely those institutions have the least mature IR programs and benefit most from the cross-tenant-scope signal.

**Finding 2.6.2 — Long-dwell adversary detection composes with chain reconciliation but the composition is not a control.** The IR playbook's "Long-dwell adversary considerations" section names UEBA, anomalous-access detection on the master-key custodian, and the spec §10.1 reconciliation cadence (weekly) as composing controls. The reconciliation cadence bounds the visibility window to one week plus remediation time. This is correct as far as it goes, but the spec §10.1 reconciliation evaluates the *fingerprint* — it does not evaluate access-pattern anomalies on the master-key custodian or the key-bytes themselves.

The gap: a long-dwell adversary with credentials but no key-bytes (as in `00-overview.md` §2.5 residual "application-host compromise that exfiltrates the IKM") will not be detected by fingerprint reconciliation. The reconciliation detects key-fingerprint drift (a wrong-key configuration), not unauthorized-but-legitimate-key access. Detection of the access pattern is bank-side UEBA, which the IR playbook references but does not bind to a chain-side operational event.

**Recommendation.** The institution-defined operational event catalog includes a `master_key.access_anomaly_detected` event the bank-side UEBA emits when it detects an anomalous access pattern on the master-key custodian. The event is correlated with the chain's `master.reconciliation_completed` events at SOC review time. A `master_key.access_anomaly_detected` followed by no chain-side `chain.verification_failure` is the long-dwell-observation pattern (attacker reading without writing); a `master_key.access_anomaly_detected` followed by chain-side fingerprint mismatch or hash mismatch is the long-dwell-action pattern (attacker writing). Naming the correlation pattern in the IR playbook gives the SOC team the playbook entry; today the playbook says "treat as IR Scenario 4 to preempt the action phase," which is correct guidance but does not say what evidence the SOC team gathers to support the determination.

**Severity.** Medium. The institution will run UEBA regardless; the gap is that the chain doesn't help the SOC team correlate UEBA signals with chain signals because the operational-event catalog does not include the access-anomaly event. A documented correlation pattern in the IR playbook closes it.

**Finding 2.6.3 — Scenario 11 sub-variant (forged rotation notice) — out-of-band cross-check channel diversity.** The remediation step says: "contact the regulator's primary IT examination liaison via a separate authenticated channel (phone call, in-person visit, regulator's encrypted-email system if separate from the notice channel)." This is correct. The gap: the institution-side procedure does not require the institution to *pre-establish* the separate authenticated channel as a documented contact path. In an actual forgery scenario, the institution's incident commander will scramble to find a phone number for "the regulator's primary IT examination liaison" at midnight on a Saturday. The institution that has not pre-established the contact path will spend hours on hold-and-callback discipline; the institution that has pre-established it will resolve the cross-check inside the notification window.

**Recommendation.** The institution's IR documentation pre-records the separate-authenticated-channel contact path: name, role, phone number, in-person address, alternate contact for off-hours, and the channel's authentication method (e.g., regulator's published GPG key for encrypted email, regulator's published phone number cross-checked against a separate regulator publication). The institution validates the contact path on the same annual cadence as the cold-DR-key dry-run consumption (per §2.1 of supply-chain doc). The institution emits a `regulator_contact_path.validated` operational event annually; the SOC team and FFIEC examiner sample-test the institution's annual validation.

**Severity.** Medium. The forged-rotation-notice scenario is rare but is a 36-hour-clock event the moment forgery is determined; the determination depends on the cross-check completing in time. Pre-establishing the channel is low-cost and high-value.

**Finding 2.6.4 — Scenario 12 branch (i) "published algorithm break" — credibility-determination discipline is missing.** Branch (i) says the disposition is no clock-start (cryptographic-deprecation event regulator-aware). The investigation triage assumes a NIST-published or FIPS-published practical attack on algorithm Y was announced. The gap: a "published algorithm break" can be ambiguous in real time. A research-paper announcement is not the same as a NIST formal deprecation; a NIST informal advisory is not the same as a formal SP 800-56 update; a vendor security advisory referencing a published attack may overstate or understate the practical-attack threshold. The institution's IR Commander needs guidance on what counts as "published" for the purposes of the no-clock disposition.

**Recommendation.** The IR playbook names the threshold for branch (i) explicitly: branch (i) applies when (a) NIST has issued a formal deprecation in the SP 800-56 or SP 800-57 series, OR (b) FIPS has issued a formal validation-impacting advisory under the FIPS 140-3 program, OR (c) the project's own security advisory (signed by the project's GPG-key role) declares the algorithm broken with reference to (a) or (b). Anything below this threshold (a research paper, a vendor advisory, threat-intel speculation) routes to branch (iii) (under investigation, with the bounded 48-72h window for the credibility determination). This protects the institution against a false-positive branch-(i) disposition during a research-paper-driven panic that turns out not to be a practical attack.

**Severity.** Medium. The post-quantum transition is on the calendar; the first published algorithm break (real or perceived) will produce a coordination challenge the spec should pre-document the threshold for.

### 2.7 Spec §10 (operational requirements) — what to commend

**§10.1 reconciliation cadence (weekly default, daily/hourly/continuous for elevated risk) bounds the master-compromise detection window to at most one week plus remediation time.** This is the load-bearing operational property — the institution can compose its broader threat-detection program with a known upper bound on the chain-side detection window. The reconciliation evidence is emitted as `master.reconciliation_completed` per §10.2, sampled by SOC and examiner against the institution's IKM roster (audit-procedures P-6 reference). Operational, evidenced, sample-testable.

**§10.6 IKM minimum length 32 bytes with the offline-grindable-fingerprint attack named explicitly.** The spec text doesn't just say "use a 32-byte key" — it explains *why* (RFC 4868 §2 HMAC-SHA-256 keying recommendation AND the public 16-byte fingerprint's offline-grindability against a low-entropy IKM). The dual enforcement (provisioning-time SHOULD + SDK-configure-time MUST) means a misconfigured deployment that attempts to register a short IKM is refused at the registry layer; if the registry is bypassed, the SDK refuses at startup. Defense in depth without requiring institutional discipline to remember.

**§10.7 software-key adapter compile-time exclusion is the right shape.** The clause says compile-time exclusion is the rule, runtime environment-variable gating is NOT sufficient. The verifier refuses any seal whose `dev_mode` is `true` or whose `kms_handle_uri` begins with `"plaintext-"` under `--strict` mode. This is the second-layer enforcement: even if the compile-time exclusion is bypassed, the verifier catches it. Three independent failures (build pipeline, deployment configuration, verifier `--strict` flag) would have to align for plaintext-key material to silently ride through to a passing examination.

**§10.8 constant-time comparison applied to BOTH fingerprint and MAC.** The spec text explicitly notes the discipline applies to both even though the fingerprint is publicly stamped on every entry — "the discipline carries, and a future maintainer extending the verifier does not reach for `==` on the MAC compare." This is the kind of implementation-discipline guidance most specs leave implicit. A future-maintainer-resilient property.

**§10.9 IKM registry retention coupled to chain-event retention.** The retention rule is stated as a coupling: retain IKMs at least as long as any chain entry stamped with that `key_version` is retained. The retire-with-active-references case requires explicit override and emits a `master_key.retired` operational event. Cloud-KMS provider default behaviors (AWS CloudHSM 7-30-day pending-window, Azure/Google similar) are named so the institution knows to configure the retention window per IKM generation. Institutional autonomy preserved; the load-bearing property (don't lose the ability to verify retained chain entries) is defended at the registry layer.

### 2.8 Spec §10 — what I would still ask for

**Finding 2.8.1 — §10.1 reconciliation evidence's chain-side integrity coupling.** The `master.reconciliation_completed` operational event is itself a chain-side event subject to spec §10.2 retention. The reconciliation procedure cross-checks the institution's IKM roster against `(tenant_id, key_version, key_fingerprint)` triples in the captured chain. The gap: if the reconciliation evidence is itself emitted into the same chain it is reconciling, an attacker who can write to the chain can also forge reconciliation events showing zero unmatched fingerprints when the actual reconciliation found mismatches.

The spec's `00-overview.md` §2.3 is honest about this — "an adversary with full server access AND full IKM access still cannot forge the day's signed root" — but the per-event chain construction is HMAC-keyed by the IKM, which the §2.3 adversary has. So the per-event integrity property does not hold against the §2.3 adversary; only the daily Merkle seal does. The reconciliation event would be sealed in the next day's Merkle root, which provides the integrity property — but the institution's incident response on a fingerprint-mismatch alert depends on reading the most recent reconciliation event, which may not yet be sealed.

**Recommendation.** The reconciliation procedure's evidence path includes an out-of-band attestation: the SOC team's manual sample-test (audit-procedures P-6) against the institution's IKM roster produces a separate operational record (signed by the SOC team's identity, not chain-side) that the FFIEC examiner consults independently of the chain's `master.reconciliation_completed` event. The chain's event is the operational signal; the SOC team's separate attestation is the integrity-anchored confirmation. This is the same pattern as the dual-trust-anchor cosign+GPG defense, applied to reconciliation evidence.

**Severity.** Medium. The §2.3 adversary is the strongest in the threat model; the reconciliation evidence's integrity-against-this-adversary depends on the daily Merkle seal having sealed the reconciliation event, which means the most recent reconciliation event is not integrity-anchored until the next day's seal. Closing the gap with a SOC-team-signed separate attestation is procedural, not cryptographic, and the spec already requires P-6 — it just doesn't name the attestation as a separate evidence artifact.

**Finding 2.8.2 — §10.7 verifier `--strict` flag is the second-layer enforcement, but the institution's examination procedure does not mandate strict mode by default.** The spec says "the verifier refuses the produced chain under `--strict` mode." This implies strict mode is opt-in. An institution's examination procedure that runs the verifier without `--strict` will silently accept seals carrying `dev_mode: true` and `kms_handle_uri="plaintext-dev"` — precisely the misconfiguration the §10.7 second-layer enforcement is supposed to catch.

**Recommendation.** The examiner-quickstart and operator-guide documentation requires strict mode by default. The verifier's default behavior under no-flag invocation is strict; an explicit `--allow-dev` flag is required to relax. This is the same pattern as fail-closed-by-default discipline applied to TLS verification: don't rely on the operator to remember to set the strict flag; make strict the default and require explicit relaxation.

**Severity.** Medium. The compile-time exclusion is the primary defense, so the strict-mode verifier flag is defense-in-depth on a deployment-misconfiguration path. Making strict the default closes the path without changing the verifier's behavior for institutions running examinations.

### 2.9 Design 04 §3.2 (HSM custody — daily seal signing keys) — what to commend

**Per-tenant: one signing keypair per algorithm per tenant.** The spec is precise on this: not a shared signing key with per-tenant labels, not a keypair that signs across tenants with per-tenant key-derivation, but one keypair per algorithm per tenant, held in HSM, `sign`-mode only, non-extractable. Cross-tenant signing is structurally prevented at the HSM layer. This is the strongest possible posture and the spec adopts it as the default.

**Public key registry retains historical keys with their validity windows.** Past seals remain verifiable against the previous public key after a rotation. The institution does not lose the ability to verify historical examinations after rotating the signing key. This is the right shape — most vendor designs we review treat rotation as a clean break, which produces a "do we still have the old key?" scramble at the next examination cycle.

**Algorithm-posture transitions (single → dual → single under new algorithm) are operationally documented.** The institution's transition path through the post-quantum migration is named: provision the new-algorithm keypair, publish the public key, update declared algorithm posture, then later retire the first algorithm. This is the kind of pre-documentation that takes most institutions a separate engagement to produce.

**Cross-region key replication preserves per-tenant labels.** AWS CloudHSM backup-and-restore, Azure Managed HSM geo-redundant, Google Cloud HSM cross-region — all preserve the per-tenant key-isolation property. Cross-tenant signing remains structurally prevented across regions. The institution's DR plan documents the cross-region replication procedure; the spec does not normate it because the procedure is HSM-vendor-specific.

### 2.10 Design 04 §3.2 — what I would still ask for

**Finding 2.10.1 — Per-tenant signing keypair scaling at vendor-hosted topology.** The single-tenant case is clear: one keypair per algorithm. For a vendor-hosted multi-tenant deployment serving N institutions, the HSM holds N keypairs per algorithm. At the cloud-HSM cost structures named in `cloud-hsm-guide.md` (AWS CloudHSM ~$13k/year per cluster member; Azure Managed HSM ~$15k/year per partition; Google Cloud HSM pay-per-use with $1/month per key version), the per-tenant keypair model imposes a cost floor that scales linearly with tenant count.

The gap: `00-overview.md` §6.5 names community-bank-tier institutions as candidates for shared-cloud-HSM deployment, but the per-tenant keypair model means even a "shared" cloud HSM holds one keypair per tenant — sharing applies to the HSM hardware, not to the keypair allocation. The cost model documentation (`docs/cost-model.md`, not in scope for this review) presumably addresses this, but the institution's third-party-risk-management for a vendor-hosted deployment will want to see the per-tenant-keypair cost passthrough explicitly named.

**Recommendation.** The vendor-hosted topology documentation explicitly states the per-tenant keypair cost structure and how it is passed through to the institution's contract. This is not a security finding — the per-tenant keypair model is correct — it is a transparency finding. The institution adopting vendor-hosted should not be surprised at year-2 renewal by a per-tenant cost-passthrough line item it did not see at onboarding.

**Severity.** Low (transparency, not security). Including for completeness because it is a partner-review question.

**Finding 2.10.2 — Cross-region replication's per-tenant key-isolation property is asserted but not normated.** The spec text says cross-region replication preserves the per-tenant signing key's isolation across regions ("cross-tenant signing remains structurally prevented"). This is the load-bearing property. The institution's cross-region DR plan tests the replication, but the institution does not have a normative test the spec defines for the cross-tenant-signing-prevention property post-replication.

**Recommendation.** The conformance corpus (`spec/test-vectors/`) includes a cross-region replication test vector: institution provisions tenant A's keypair in region 1, replicates to region 2, attempts to sign tenant B's seal under tenant A's key in region 2. The test passes only if the attempt is structurally refused at the HSM layer. The institution runs this test as part of its annual DR plan exercise; the test produces a `cross_region_isolation_test.completed` operational event that the SOC team and FFIEC examiner sample-test.

**Severity.** Medium. The cross-region property is load-bearing; an institution operating multi-region production without testing the property is defending it on documentation alone.

### 2.11 Design 09 §2.9 (Adversary I — examiner-side tooling subversion) — what to commend

**The reception-and-validation procedure for a regulator-held fingerprint rotation has three named operational events providing audit-evidence.** `regulator_fingerprint.rotation_received`, `regulator_fingerprint.rotation_validated`, `regulator_fingerprint.installed`. Each event has explicit fields. The reception-failure sub-variant (forged notice suspected) is named with its own clock-start trigger (Scenario 11 sub-variant). This is exactly the documentation depth a third-party-risk function looks for when evaluating a vendor's trust-anchor-rotation discipline.

**The institution-side responsibility is named explicitly: examiner-laptop hygiene is regulator-side.** The spec text says: "Examiner-laptop hygiene is a regulator-side control, NOT an institution-side control." This is the line the institution's CC6 evaluation does not need to defend. A SOC engagement partner reading this knows the residual is out of scope for the institution's evaluation; the FFIEC's own examiner-IT program defends it.

**Three-independent-compromises framing for the silent-pass attack.** Cosign + GPG-signed manifest + regulator-held fingerprint must align for the verifier to silently pass forged data. The framing is correct and conservative — even cosign-AND-GPG dual-compromise (the supply-chain doc's edge case) does not collapse the chain's integrity claim because the regulator-held fingerprint is an independent third anchor.

### 2.12 Design 09 §2.9 — what I would still ask for

**Finding 2.12.1 — Step 4 historical-report re-validation does not specify the SOC artifact.** The §2.9 procedure step 4 says: "Re-validate any historical verifier reports that depended on the old fingerprint, against the new fingerprint, to confirm the institution's verifier output remains stable across the rotation." The clause is operationally correct but does not name the artifact the institution produces from the re-validation. The next FFIEC examination will ask: "show me the re-validation evidence for the fingerprint rotation 18 months ago." Without a named artifact, the institution's response varies by IR Commander — some institutions will produce a memo, others will produce a verifier-output diff, others will produce nothing because the verifier outputs are kept in a staff-shared drive and no formal artifact was emitted.

**Recommendation.** The §2.9 procedure step 4 names the artifact: a `regulator_fingerprint.historical_revalidation_completed` operational event with fields `revalidation_window` (date range), `reports_revalidated` (list of historical verifier reports re-checked), `pre_rotation_outcome` (PASS/FAIL/anomaly summary), `post_rotation_outcome` (PASS/FAIL/anomaly summary), `outcome_stable` (boolean). The SOC team and FFIEC examiner sample-test this event against the institution's archived verifier reports.

**Severity.** Medium. The institution's CSF GV.SC-04 evidence depends on the re-validation being operationally exercised; without a named artifact, the exercise is documented but not auditable retroactively.

**Finding 2.12.2 — Reception channel diversity is named but not normated.** The §2.9 procedure step 1 says the regulator publishes the new fingerprint via "the institution's established regulator-communication channel (typically the regulator's encrypted messaging system or a paper notice with hand-delivered envelope to the institution's compliance officer)." The clause is "typically" — meaning the institution's procedure may use a single channel. A single-channel reception is vulnerable to a channel-substitution attack: if the attacker controls the institution's primary regulator-communication channel, the attacker controls the rotation-notice path.

The §2.9 procedure step 2 names the cross-check: "cross-checks against the regulator's parallel notification on the regulator's authenticated-domain channel." But the cross-check is in step 2 (validation), and the channels are not pre-bound to the institution's IR documentation. The institution may operate the cross-check ad hoc.

**Recommendation.** The institution's IR documentation pre-records two independent regulator-communication channels (e.g., regulator's encrypted-email system AND regulator's authenticated-domain web portal AND regulator's published phone number, of which the institution selects two as the primary cross-check pair). The two channels are validated annually (per §2.6.3 above) and emit a `regulator_communication_channel.validated` operational event. A rotation notice arriving on only one of the two channels does NOT validate; the institution requires both-channels-converging before installing the new fingerprint.

**Severity.** Medium. The forged-rotation-notice scenario is rare but high-impact; the dual-channel reception requirement is a concrete defense the institution operationally exercises rather than relying on the IR Commander to invent at incident time.

---

## 3. Cross-document gaps — partner-level observations

### 3.1 Operational-event catalog completeness (cross-cutting across §10.2, §2.9, IR playbook)

The spec §10.2 catalog names the operational events the implementation MUST emit. The IR playbook references additional events the institution-side procedures emit (`regulator_fingerprint.*`, `master.reconciliation_completed`, `audit_file.truncation_detected`, `mirror.reconciliation_completed`). My findings 2.6.2, 2.6.3, 2.10.2, 2.12.1 above propose four additional events: `master_key.access_anomaly_detected`, `regulator_contact_path.validated`, `cross_region_isolation_test.completed`, `regulator_fingerprint.historical_revalidation_completed`.

The cross-cutting observation: the operational-event catalog is split between spec §10.2 (chain-side, normative) and various institution-defined event categories (institution-side, non-normative). The institution's SOC team and FFIEC examiner consume both. The split is correct (institution-defined events are institution-policy-driven, not chain-construction-driven) but the spec does not provide a single referenced catalog of expected institution-defined events for the FFIEC examiner to sample-test against.

**Recommendation.** A separate document — `docs/soc-pack/institution-defined-events.md` (referenced from spec §10.2 and from the SOC/examiner pack documentation) — enumerates the institution-defined operational events the institution's procedures emit, with field definitions, retention requirements, and the FFIEC-examiner sample-test guidance per event type. This is the SOC-pack-side counterpart to spec §10.2's chain-side catalog. The institution adopting the spec instantiates the document by selecting which events it emits (with rationale where it omits any) and which fields its procedures populate.

**Severity.** Medium. The spec's normative discipline is on chain-side events; institution-side events are operational policy. The gap is that the institution's policy-formation is unguided — the institution's SOC team writes the institution-defined event catalog from scratch, which produces inconsistency across institutions and additional examination overhead.

### 3.2 Trust-anchor rotation cross-document timing alignment

Three trust anchors operate on three rotation cadences:

- Cosign: per release (event-driven, with emergency rotation within 24h of compromise).
- GPG: 36 months default, 12-month subkey, emergency 7 days.
- Cold-DR: 60 months default, key-holder rotation per departure, annual dry-run cadence.
- Regulator-held public-key fingerprint: institution-driven (per `09-threat-model.md` §2.9, "rotated by institution's tenant signing key rotation OR by regulator's fingerprint-storage system rotation").

The cadences are independent. The institution's release-validation procedure consumes attestations from all four. The institution's IR playbook references all four (Scenario 3, Scenario 11, Scenario 11 sub-variant). There is no single document the institution's compliance officer consults to confirm all four cadences are on schedule, no missed-rotation overdue, and the institution's archived attestations are complete.

**Recommendation.** The institution's standing IR documentation includes a trust-anchor calendar: a dashboard or static document listing each trust anchor, its rotation cadence, the next expected rotation date, the institution's archived attestations for past rotations, and the current status (on-schedule / overdue / not-applicable). The dashboard is reviewed annually as part of the institution's IR-program update cadence. The FFIEC examiner consults the dashboard as a single artifact rather than reconstructing it from four separate document references.

**Severity.** Medium. Each cadence is correctly documented in isolation; the gap is the consolidated view. A tier-1 institution will produce this dashboard from its own controls inventory; a tier-2 institution will not, and the FFIEC examiner will spend hours reconstructing it.

---

## 4. Disposition by CSF function

For the engagement memo's CSF-function summary:

| CSF function | Disposition | Notes |
|---|---|---|
| GV.SC (cybersecurity supply chain) | **Strong** | SLSA L3 institutional consumption, GPG/cosign dual trust path, cold-DR fallback with annual dry-run consumption, mirror audit-log retention. Findings 2.2.1 and 2.2.2 are ecosystem-signal gaps, not institution-side gaps. |
| ID.RA (risk assessment) | **Strong** | Threat model §2.9 reception procedure, three-independent-compromises framing, residual-risk register. Finding 2.12.1 is an artifact-naming gap. |
| PR.AC (access control) | **Strong (BYOC)** | IAM permission matrix, vendor-cannot-reach-HSM control at three layers, unilateral revocation. Finding 2.4.2 is an upgrade-gating gap. |
| PR.DS (data security) | **Strong** | Per-tenant HKDF binding, key fingerprint stamping, constant-time comparison, IKM minimum length, registry retention. Finding 2.10.2 is a cross-region test-vector gap. |
| DE.CM (continuous monitoring) | **Adequate** | Reconciliation cadence (weekly) is the primary signal. Finding 2.6.2 (long-dwell composition) and Finding 2.8.1 (reconciliation evidence integrity coupling) are composition gaps. |
| RS.MA (incident management) | **Strong** | Clock-start triage, constellation rollup, federal-regulator routing per charter, dual-algorithm verifier-version timing-overlap matrix. Findings 2.6.1, 2.6.3, 2.6.4 are coordination gaps. |
| RC.RP (recovery planning) | **Adequate** | Backup integrity failure has its own scenario (10), unverifiable-gap clock-start trigger is named. The gap is upstream — see Finding 2.10.2 on cross-region replication test. |

No CSF function rated below adequate. Findings I would brief at partner-review are 2.4.1 (BYOC redaction policy), 2.6.1 (cross-tenant scope discovery), 2.10.2 (cross-region isolation test), and 2.12.2 (dual-channel reception requirement). Partner-level escalation is reserved for findings that change the engagement's recommend/not-recommend disposition; none of the findings above do.

---

## 5. Partner-review summary

The chain-of-custody v1.0 artifacts present a design unusually well-aligned with FFIEC examination posture, CSF GV.SC controls, and BYOC third-party-risk management. The findings I have written are refinements, not gates. An EU institution's third-party-risk function evaluating this for adoption would brief partner-review with a recommend disposition, conditional on the institution's own SOC team adopting the institution-defined operational event catalog implied by Findings 2.6.2 / 2.6.3 / 2.10.2 / 2.12.1, the trust-anchor calendar implied by Finding 3.2, and the BYOC redaction-policy baseline implied by Finding 2.4.1.

The strongest properties — three-independent-compromises trust-anchor framing, clock-start triage with named operational signals, and the institution-side reception procedure for regulator-held fingerprint rotation — would individually distinguish this design from the AI-decision-evidence vendor designs my engagement has reviewed in the EU market over the last 12 months. Partner will read the engagement memo and ask why the rest of the market is not at this bar yet. The honest answer is that the rest of the market has not yet absorbed the FFIEC computer-security incident notification rule's clock-start discipline into their IR taxonomies, and the spec authors have.

End of review.
