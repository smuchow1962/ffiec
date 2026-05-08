# 02 — Big Four Cyber Risk Senior Manager (round 9)

> **Role.** Cyber-risk and key-management lens. First-look review of the v1.0-rework corpus.

## Persona — Anneliese Vandermeer

**Background.** Twenty years inside Big Four cyber risk advisory, the last twelve focused on top-25 US banks. Lead advisor on three large HSM-migration programs (one full Thales-to-CloudHSM cutover; two greenfield Azure Managed HSM rollouts), and engagement-quality reviewer on five SOC 1 / SOC 2 reports for AI-adjacent fintechs. Comfortable in PKCS#11 logs, IAM trust policies, NIST CSF 2.0 mapping work, and the conversation with FFIEC examiners about whether a control description holds up under examination. Spent the 2024–2025 cycle building BYOC vendor-management playbooks for two G-SIBs. Carries the scars of two real master-key compromise incidents (one a misconfigured backup, one a contractor with persistent KMS Decrypt rights).

**Reading angle.** I want to know whether the operational property the spec claims is the operational property the examiner will see. Cryptographic correctness is table stakes; what I look at is whether the controls are enforced at the right layer (compile time, not config), whether the rotation story survives contact with a real bank's change-management cadence, whether the BYOC IAM boundary actually prevents the vendor from reaching the keys, and whether IR has a named playbook for every chain-detected event the verifier can produce. I read the §10 operational requirements carefully because §10 is where the rubber meets the road in an examination — clean §4 cryptography behind sloppy §10 operations is a finding waiting to happen.

### Questions

**Q1. Does the IKM rotation story survive contact with a real bank's operational reality — multiple long-lived application processes, a bounded rotation window, a session-key cache, and a `key_versions` list on the rotation-day seal — without producing a window where a stale process is silently producing chain entries the verifier cannot resolve?**

**Status:** Answered

The rotation model is internally consistent and operationally honest. Spec §3 defines `key_version` as a non-negative integer (≥ 1) stamped on every entry, and §4.1 inviolate property #3 requires the verifier to assert `key_fingerprint` against the looked-up IKM **before** any MAC compute. `04-hsm-custody.md` §3.1 states the spec does not require periodic rotation — rotation is a tenant policy decision — and that "past entries (with the old `key_version` and old `key_fingerprint`) remain verifiable as long as the old IKM is retained in the tenant key registry." That last clause is the operational floor I needed to see: the institution carries an IKM-retention obligation aligned with its event retention.

The rotation-window mechanics are spelled out in `04-hsm-custody.md` §3.2.2: during the window, both versions are valid; events under the old `master_version` validate against the old IKM, events under the new `master_version` validate against the new IKM; the seal record covering the rotation day records `key_versions` as a list (e.g., `[3, 4]`). `06-ledger-server-design.md` §3.4 carries this through to the schema (`key_versions BIGINT[] NOT NULL`, with a `CHECK (array_length(key_versions, 1) >= 1)`). The verifier dispatches per-entry on `entry.key_version` (`07-verifier-design.md` §4.1 step 7), so a stale-process entry under the old version still resolves correctly as long as the old IKM remains registered.

The catch I was looking for — what happens when a stale process holds an IKM the registry has already retired — is correctly handled by spec §7 step 7: a null `ikm_lookup` produces `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` and **no MAC compute happens on lookup miss**. That is the right failure mode: the verifier names the operational fault (premature retirement of an IKM still being used by a long-lived process) rather than burying it in a MAC-mismatch storm. The IR playbook covers the analogous compromise scenario in Scenario 4.

What's missing from a strictly cryptographic reading but present in the operational reality: the IKM retention obligation should propagate into the `master_key.rotated` operational event. `06-ledger-server-design.md` §7.3.1 lists the event with fields `tenant_id, old_version, new_version, custodian` — it captures the rotation but doesn't flag the retention floor. That's a documentation refinement, not a gap in the spec.

**Q2. Is the plaintext-KMS adapter exclusion enforceable in a way that survives a misconfigured deployment — i.e., does the exclusion live at compile time, not at run-time configuration, and is the failure path on the verifier side regulator-visible?**

**Status:** Answered

This is the question I was most prepared to find a partial answer to, and the spec answers it cleanly. §10.7 has three normative requirements stacked:

1. **Compile-time exclusion.** "Be excluded from production release builds at compile time. Run-time environment-variable gating is NOT sufficient — a misconfigured deployment that flips the flag must NOT be able to bring the software adapter online in production." That is the right enforcement layer; it removes the entire class of "operator forgets to set FFIEC_DISABLE_PLAINTEXT=1" failures.
2. **Per-entry stamp.** Every chain entry the dev adapter produces is stamped with `kms_handle_uri = "plaintext-dev"` (or another `"plaintext-"` prefix). That stamp is recorded on the verifier-visible chain entry per §3 and §4.4.
3. **Verifier refusal under `--strict`.** The verifier rejects any seal whose `dev_mode` is `true` or whose `kms_handle_uri` begins with `"plaintext-"`. Spec §7 step 12 codifies this; `07-verifier-design.md` §4.3 shows the implementation.

The seal record carries `dev_mode` independently (§4.2 schema and `06-ledger-server-design.md` §3.4), and the IR playbook Scenario 6 names the failure mode ("Software-key fallback in production") with a clear remediation path: re-issue the affected seals from the HSM-backed key.

What I assessed for completeness: the double-protection (compile-time exclusion **and** verifier-side refusal) is the regulator-visible line. Either alone would be a single point of failure — a build pipeline that mis-tagged a release, or a verifier with a missing flag — and the spec correctly insists on both. The `kms_handle_uri = "plaintext-dev"` stamp on every entry is the third leg: even if a seal somehow shipped without `dev_mode=true`, the per-entry stamp catches it at the entry level. That is defense-in-depth at the right granularity.

The only operational refinement worth flagging in a real bank context: §10.7 says "Implementations MAY ship a software-key adapter for development and test." The MAY leaves vendor distributions free to ship without the adapter at all — which is the safest posture for a regulated institution and which I'd recommend the BYOC `vendor-management` controls call out explicitly. That is an institution-side procurement specification, not a spec gap.

**Q3. Does the BYOC IAM boundary prevent the vendor from reaching the bank's HSM, master-key custodian, and event-payload data through any path the documentation actually exercises — including the image-pull path, the support-telemetry path, and the secret-manager read path?**

**Status:** Answered

`byoc-deployment.md` §"IAM permission matrix" is the load-bearing artifact. Three rows answer the question:

- **HSM signing operations.** Bank's chain-ops role: yes (via seal-job role only). Vendor's support role: no. Vendor's image: yes (via seal-job role only; runs as the image's IAM identity). The vendor's image runs under a bank-controlled IAM identity that the bank can revoke unilaterally. The vendor's account cannot assume the role from outside the bank's account.
- **HSM administration (extract, import, delete).** Bank: yes (via separate HSM admin role). Vendor: no. Vendor's image: no. That's the right separation: the workload signs, but cannot extract.
- **Bank's event payload data.** Bank: yes. Vendor's support role: **No** (bolded in the source). Vendor's image: yes (the image processes it). The vendor's runtime sees the data because it has to; the vendor's people don't.

The two cross-boundary flows that worried me — image pull and support telemetry — are explicitly mediated through bank-controlled hops in `06-ledger-server-design.md` §6.1.1 and `byoc-deployment.md` §"Step 2" + §"Step 5":

- **Image pull.** The bank operates a mirror registry (Harbor, ECR, ACR, GAR). The mirror pulls from the vendor's registry through an explicit, audit-logged egress path; the mirror verifies cosign signatures before storing the image. The ledger workload pulls from the bank's mirror, never from the vendor's registry directly. That removes in-flight image substitution and gives the bank a chokepoint for scan-and-approve.
- **Support telemetry.** Bank-controlled OTel Collector receives operational events and metrics; applies a documented redaction policy (no event payload data; control-plane events only); forwards the redacted stream to the vendor over a bank-managed egress. The bank can revoke the egress unilaterally without coordinating with the vendor. The privacy-impact assessment for the support relationship documents the data flow.

The HSM PIN read path is the third sensitive flow — `06-ledger-server-design.md` §6 reads the PIN from an environment variable populated by the bank's secret-management system at process start. The IAM matrix grants the vendor's image read-only on one secret only (the HSM PIN), and `04-hsm-custody.md` §5.2.1 spells out the rotation procedure: HSM admin generates a new PIN at the HSM; seal-job operator updates the secret-management system; the running processes re-read on next reload; HSM admin invalidates the previous PIN. Quarterly cadence is the typical bank operational rhythm and the spec admits it.

The CC-mapping in §"Audit considerations" maps cleanly to SOC 2 TSC: CC6.1 (logical access — IAM permission matrix), CC6.7 (restricts movement of data — bank-controlled OTel Collector, redaction policy), CC6.8 (prevents/detects unauthorized software — mirror-registry signature validation, cosign on every pull), CC9.2 (vendor management — bank can revoke vendor access unilaterally). Those are the four that examiners reach for in a vendor-hosted-or-BYOC engagement, and they're addressed at the right layer.

**Q4. Does the cosign-plus-GPG dual-signing trust path actually compose into a defense against the realistic compromise scenarios — Sigstore root compromise, project release-pipeline compromise, in-flight binary substitution at the institution — and is the institution-side validation procedure mechanizable rather than a checklist of human steps?**

**Status:** Answered

The dual-signing path is well-thought-through and the institution-side procedure is mechanizable, which is the part I most wanted to confirm.

`supply-chain.md` §"Trust path summary" lays out three composed paths:

- **Primary.** Institution receives binary; validates cosign signature against pinned project public key; Sigstore validates the signing certificate.
- **Fallback (if Sigstore is compromised).** Institution validates GPG-signed hash manifest against cached project GPG public key; computes binary SHA-256 and compares to the manifest.
- **Defense-in-depth.** Institution rebuilds from source using the same toolchain version; compares the rebuilt binary's SHA-256 to the published binary; archives the rebuild log as control evidence.

The cosign and GPG keys are held by separate small groups within the project's release-management role, with the GPG key offline-held primarily for emergency-response (`supply-chain.md` §"Release-pipeline operator governance"). That separation is what makes the fallback path actually a fallback rather than a co-compromise. Institution-side trust anchors (cosign public key, GPG public key, toolchain version pinning) are cached out-of-band and validated against the spec text itself rather than against a separate distribution channel — that's the right anchoring choice; the spec text is what the examiner has direct visibility into.

The mechanization is in `07-verifier-design.md` §8.5 — `verifier-validate.sh` (Linux/macOS) and `verifier-validate.ps1` (Windows) run the full validation chain (cosign + reproducible-build manifest + GPG fallback) and produce a single pass/fail. Examiners run the wrapper before running the verifier; the wrapper exits 0 only if all three checks pass. The script is itself signed and reproducible; its hash is published alongside the verifier binary.

The reproducible-build evidence schema (§"Reproducible-build evidence") records `rebuild_timestamp, source_commit, go_version, build_flags, binary_sha256, matches_published, verifier`. That's the right shape for a control-evidence repository: machine-readable, sample-testable by SOC and examination teams, archivable.

Three compromise scenarios I checked the spec against:

- **Sigstore root compromise.** Falls back to GPG-signed manifest. The institution holds the GPG public key out-of-band, so a Sigstore-only compromise does not collapse the trust path.
- **Project release-pipeline compromise** (someone pushes a malicious binary signed with the project's cosign key). Defense is reproducible builds — the institution rebuilds from source using the pinned toolchain and compares hashes. This is the residual-risk scenario the project specifically calls out: "A subverted build that produces incorrect output fails the corpus" + "External security review of the corpus content periodically." The conformance corpus is the second leg — a subverted binary that produces deterministic-but-wrong output gets caught by the corpus on the next run.
- **In-flight binary substitution at the institution** (someone MITMs the institution's pull from the project). Cosign signature mismatch catches it; GPG manifest catches it independently.

The compound-compromise threshold the spec explicitly names: "Three independent compromises (cosign, GPG manifest, regulator-held fingerprint) would have to align for the verifier to silently pass forged data" (`00-overview.md` §2.4). That's a defensible posture for the regulator audience.

What I assessed for completeness: the SBOM is published in CycloneDX 1.5, which is what most bank vulnerability-management tools consume directly (Snyk, Mend, JFrog Xray are explicitly named). Trivy and Grype scan results are published per release. The container path follows the same trust chain (distroless base image, SBOM, signed image, scan results). The institution's container registry mirror — not a direct pull — is the deployment pattern. That's the pattern I'd expect from a mature program, and it's documented at the right level for an SOC team to test against.

**Q5. Does IR have a named playbook for every chain-detected event the verifier can produce — including the load-bearing `key_fingerprint mismatch` and `unknown_key_version` failures from spec §7 steps 7 and 8 — with notification thresholds that align with the FFIEC computer-security incident notification rule and CIRCIA?**

**Status:** Answered

The IR playbook has named scenarios for the load-bearing chain-detected events, with severity classifications and notification timelines that align with the federal frameworks.

`incident-response-playbook.md` §"Scope" lists the events the playbook covers: chain hash mismatch, chain link broken, Merkle root mismatch, signature verification failed, sealing delay >72h, master key compromise, software-key fallback in production, backup integrity failure. The severity classification matrix is explicit: Critical (15 min ack / 1h contain) for signature verification failure, suspected master compromise, software-key in production; High (30 min ack / 4h contain) for chain hash mismatch and merkle root mismatch; Medium (1h ack / 24h remediate) for sealing delay >72h and backup-recovery integrity failure; Low (24h ack / next business day) for late-binding and clock-skew anomalies.

The two specific failure modes I came in checking for — `key_fingerprint mismatch` (spec §7 step 8) and `unknown_key_version` (spec §7 step 7) — are addressed but mapped onto Scenarios 1 and 4 rather than getting their own scenarios. The mapping is correct: a `key_fingerprint mismatch` is the operational signal the §10.1 weekly reconciliation produces, which feeds Scenario 4 (master-key compromise suspected) when the unmatched fingerprints don't trace to an authorized rotation; an `unknown_key_version` is the signal that a stale process is producing chain entries against a retired IKM, which traces back to either an incomplete rotation rollout (operational; Scenario 1 root-cause investigation) or a process-level compromise producing forged events (Scenario 4 again). Both signals are captured by the verifier's per-step failure taxonomy (`07-verifier-design.md` §6 failure-record shape includes `step` referencing the spec §7 step number) and both have a documented containment + remediation path.

The notification framework alignment is explicit and useful: §"CIRCIA and additional notification frameworks" tabulates the parallel paths — FFIEC computer-security rule (36 hours, primary federal regulator), CIRCIA (72 hours, CISA), state data-breach notification (state-specific, state AGs), CFPB UDAAP (per CFPB rules), sealing-delay (chain-specific, primary regulator at 72h). That's the right framing for a bank IR program: multiple paths apply, counsel decides which.

The HSM tamper-detection integration in §"HSM tamper-detection integration" is a touch I appreciated — the HSM emits a vendor-specific tamper event; the playbook routes it preemptively to Scenario 3 (signature verification potentially compromised) before the verifier sees a failure. That's the right risk posture: tamper detection is a leading indicator of seal compromise, not a coincident one.

The long-dwell adversary section names the chain's blind spot honestly: "long-dwell observation does not produce chain-detectable signals." The mitigation is composition with the institution's broader detection posture (UEBA on chain access patterns, anomalous-access detection on master-key custodian, weekly reconciliation cadence, threat-intel cross-correlation). That's the right level of honesty for a federal examination — the chain detects modification at action time, and the institution's broader controls cover the dwell period.

**Q6. Are HSM cluster failover procedures documented at the level a bank ops team can operate against — per-cloud, with explicit cross-region replication semantics and a tested DR cadence — and does the spec admit the operational reality that HSM unavailability is a "delay the seal" event rather than a "stop ingest" event?**

**Status:** Answered

The HSM unavailability semantics are correctly partitioned: per-event HMAC is independent of the HSM, so capture continues; the daily seal job retries with exponential backoff; institutions notify their primary regulator if the seal is delayed beyond 72 hours. Spec §4.3.1 codifies this and `04-hsm-custody.md` §5.2 spells out the operational behavior. The 72-hour threshold is the right framing for a bank ops team — it accommodates HSM cluster failover, scheduled maintenance, and regional cloud-provider incidents without requiring escalation for every transient blip, which is what would otherwise produce regulator-notification fatigue. The Cyber-incident notification under the FFIEC computer-security incident notification rule (36 hours) is correctly distinguished as a separate path that applies when the HSM unavailability is associated with a suspected security incident.

Per-cloud HSM cluster failover and replication is in `cloud-hsm-guide.md`:

- **AWS CloudHSM.** Cluster of at least 2 HSM members for HA; cross-region uses backup-and-restore; the bank's DR plan tests the procedure annually. Per-tenant key labels are preserved through replication; cross-tenant signing remains structurally prevented.
- **Azure Managed HSM.** Geo-redundant configuration with key labels maintained; the bank's DR plan tests failover.
- **Google Cloud HSM.** Cross-region replication with the same property.
- **On-prem (Thales Luna, Entrust nShield, Utimaco, Yubico YubiHSM 2).** Vendor-specific procedures; "all current models meet L3."

`04-hsm-custody.md` §3.2.1 reinforces the per-cloud cross-region story and adds the detail that "AWS KMS Custom Key Stores backed by CloudHSM are conformant: the cryptographic operations are performed by the underlying CloudHSM, and the FIPS validation flows through." That's a useful clarification — a lot of bank teams default to KMS for the IAM ergonomics and don't realize the FIPS validation flows through when CloudHSM backs it.

The cluster-member field in the seal record (`seal.hsm_cluster_member` per spec §4.2 schema, advisory) is the right level of granularity for forensic purposes — you can trace which cluster member signed a given seal without making it a security-bearing field.

`06-ledger-server-design.md` §7.5 ties the HSM failover into the broader DR / RPO / RTO discussion: WAL streaming replication for sub-second RPO, hot store derivable from WAL, daily seals stored in WAL-backed table. The async-vs-sync replication trade-off table is the correct framing for a bank's DR architect: same-region sync adds 1–3 ms (fits the <2 ms p99 hot path budget); cross-region sync adds tens of ms (defeats the budget). Most banks land on async same-region for the latency, with async cross-region for resilience without hot-path cost.

Pitfalls in `cloud-hsm-guide.md` §"Pitfalls" are the practical ones I'd ask an SOC team to test against: don't mix HSM tiers (a tenant with one L3 and one L2 key has L2 as the weakest link); don't share keys across tenants (per-tenant required by spec; cross-tenant sharing is non-conformant); don't keep backup keys without testing restore (a backup that is not tested for restore is not a backup); don't let PIN management drift (a leaked-but-not-rotated PIN is a vulnerability). Those map directly to the SOC 2 controls examination evidence.

The one operational refinement worth noting for a real bank context: the cost picture in the per-cloud sections ($26k/year for a 2-member CloudHSM cluster, ~$15k/year per Azure Managed HSM partition, $5–15k/year for Google Cloud HSM at typical volumes, $30–80k/year all-in for on-prem) is the right level of transparency for a CIO-level decision conversation. The implicit message — daily cadence is operationally affordable at every tier; weekly cadence with examiner approval is the lever for community-bank deployments — matches what I'd recommend in a procurement workshop.

---

## Per-role roll-up

| Status | Count |
|---|---|
| Answered | 6 |
| Partial | 0 |
| Gap | 0 |
