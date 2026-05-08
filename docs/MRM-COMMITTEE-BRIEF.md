# MRM Committee brief

> **What this doc is.** Template for a 2–4-page brief the bank's Model Risk Management committee chair reads instead of the full design docs. The bank fills it in and presents.

## Purpose

The bank's Model Risk Management Committee has authority over the model-risk framework under SR 11-7 and OCC Bulletin 2011-12. This brief informs the committee about the chain-of-custody implementation, which is a control supporting the bank's MRM program for AI agent decisions.

## What the chain is, in two paragraphs

The chain captures every AI agent decision the bank makes, with a cryptographic record that any later party — internal audit, external auditor, regulator — can verify the decision was the actual decision the agent made, in the order it made it, on the date it claims to have made it.

The chain produces integrity, not validity. It proves the record is authentic. Whether the AI decision the record describes was a *good* decision is the MRM Committee's question, not the chain's. The chain is the evidence layer; MRM operates on top of the evidence.

## What this changes for MRM

| Before the chain | After the chain |
|---|---|
| Vendor logs are the bank's evidence of AI decisions | Cryptographic chain proves the record is authentic regardless of vendor cooperation |
| Reproducing a specific decision requires vendor support | The bank can re-construct the run from chain data |
| Drift detection runs on uncorroborated logs | Drift detection runs on integrity-bearing records |
| Customer disputes about AI decisions are vendor-dependent | The bank produces the integrity-bearing record directly |
| Regulator requests for AI evidence go through vendor | The bank produces verifier-ready snapshots directly |

The MRM Committee's qualitative judgment is unchanged. The evidence supporting that judgment is now stronger.

## The four primitives (short version)

1. **HMAC chain at capture.** Each event is cryptographically tied to the event before it AND to the tenant identity. An outside party cannot insert a fake event into the chain. The MAC IS the chain entry — it is persisted byte-for-byte, not computed and discarded; the verifier compares against it.
2. **Daily Merkle seal.** Each day's collection of events is summarized by a single 32-byte hash. Modifying any event changes the hash; the hash is recorded under a hardware-protected signature. **This primitive lives at the ledger server (the bank's perimeter), not the SDK** — the seal exists to detect tampering by the server itself, so the seal must be unforgeable by the server.
3. **HSM-rooted root signature.** The daily hash is signed by a hardware key that even a fully compromised internal team cannot extract. This is the integrity anchor.
4. **OpenTelemetry-native wire.** The chain ships over the same wire as the bank's existing observability stack. No rip-and-replace.

The four primitives together produce the integrity property. Any one of them missing weakens the property.

### Three additional defensive properties from the v1.0-rework

Beyond the four primitives, the v1.0-rework adds three defensive properties the committee should understand:

**Per-tenant cryptographic binding.** The chain's key derivation binds the tenant identifier into the cryptographic step itself, so two tenants' chain entries are not interchangeable even if an operator accidentally points one tenant at another's key material. Cross-tenant chain-of-custody confusion is *structurally prevented* by the construction, not just operationally avoided. A botched rotation that re-uses a key version number for a different key produces a precise, named "wrong key" failure at the verifier — not an ambiguous failure that looks like data tampering. The corresponding operational evidence is the institution's weekly key-fingerprint reconciliation report (audit-procedures P-6).

**Compile-time prevention of dev-key material in production.** The development-only software-key adapter is excluded from production builds at compile time, not at run time. A misconfigured deployment that flips an environment flag cannot bring the development adapter online in production. The verifier additionally refuses any seal carrying the dev-mode marker under strict mode. Two independent enforcement layers; either alone would be a single point of failure.

**Constant-time comparison discipline.** The verifier uses constant-time comparison primitives for both the integrity check and the identity check. Side-channel attacks against the verifier process are closed at the cryptographic-discipline layer; no exposure during verification.

The committee chair should understand: where the per-event integrity primitive lives (SDK on the application host, on the bank's perimeter), where the daily integrity primitive lives (server-side at the ledger, with HSM signing in the bank's perimeter or per-tenant-segregated vendor equivalent), and where the trust anchors live (HSM hardware key for the daily signature; regulator-held public key for the verifier's signature check). The custody placement of each primitive is what makes the integrity claim load-bearing in vendor-hosted and BYOC topologies.

## How this fits SR 11-7

SR 11-7 requires effective challenge of model behavior across the model lifecycle. The chain supports this in specific ways:

- **Reproducibility.** Validators can pull a specific decision from the chain, including the prompt, the model response, and the parameters (when captured), and reproduce or evaluate it.
- **Ongoing monitoring.** Drift analysis, accuracy tracking, and bias evaluation all run on chain-bearing data with confidence the data is what the agent actually produced.
- **Independence.** The verifier produces output that does not depend on the AI vendor's continued cooperation. The MRM team can validate decisions even after a vendor change.
- **Documentation.** Chain output (PDF + JSON reports) supplements the MRM committee's documentation of effective challenge.

The chain does NOT satisfy SR 11-7 by itself. SR 11-7 is a model-validation framework; the chain is a logging-integrity control. The two compose; neither replaces the other.

## What the MRM Committee should ask

Standard questions for the chain owner:

1. Which AI agent platforms are in scope? Which are not, and why?
2. What is the seal cadence? Has it been relaxed? If so, with what regulator approval?
3. What model-state fields are captured in `gen_ai_parameters` (model id, decoding parameters, system prompt, retrieval context)? What are not, and what is the validator's compensating control? *(For SR 11-7 reproducibility, the chain integrity-binds whatever the institution puts in `audit.*` and `gen_ai_parameters`; the institution's `audit.*` schema is a committee-relevant decision and should be reviewed periodically against the SR 11-7 reproducibility surface.)* **Composition with model inventory:** the chain-captured per-model decision count (derived from `gen_ai.response.model` aggregation) is one input to the institution's MRM model-inventory completeness check — a model that produces decisions but does not appear in the inventory is a coverage gap; a model in the inventory that produces no chain-captured decisions is potentially decommissioned-but-still-listed. The committee asks the chain-ops team for the per-model decision-count distribution per quarter and cross-checks against the institution's MRM model inventory.
4. How are master-key (IKM) rotation and key-version transitions managed? **Specifically: walk us through the most recent key-fingerprint reconciliation report (the `master.reconciliation_completed` operational event), including any fingerprint mismatches during the period and how they were resolved (audit-procedures P-6).**
5. Has the verifier been run for the last quarter? What were the results? Were any seals stamped with `kms_handle_uri = "plaintext-dev"` (production dev-mode failure)?
6. Are there any open chain-detected events under investigation? What is the timeline? **For any `key_fingerprint mismatch` finding: was the IKM-roster row identified, the change-management approval recorded, and the verifier re-run against the corrected roster?**
7. What is the institution's incident response plan for chain failures? **Specifically: does the playbook cover the new failure modes — `key_fingerprint mismatch` (Scenario 7), `unknown_key_version` (Scenario 8), `audit_file.truncation_detected` (Scenario 9)?**
8. What is the institution's posture on concurrent model-version deployment? Specifically: when multiple `gen_ai.response.model` versions are observed in a single examination period for the same decision-class, is the institution operating an A/B test, a canary deployment, a vendor-side silent re-routing, or a multi-region deployment with version drift? The intent matters for SR 11-7: an A/B test is a deliberate model-validation activity that the MRM committee oversees; vendor-side silent re-routing is a control-completeness gap the institution must surface to the vendor; multi-region drift is an operational housekeeping issue. The chain captures the per-event response-model identifier (per spec §4.4 MUST requirement); the institution's deployment-intent classification is captured per spec §4.4.2 (`audit.deployment.intent`, with the conditionally-required `audit.deployment.policy_version` and intent-specific fields like `audit.deployment.experiment_id`, `audit.deployment.region`, and `audit.deployment.canary_traffic_pct`). The MRM committee asks for the institution's deployment-intent record per `gen_ai.response.model` value observed in the period and reviews it against the §4.4.2 normative schema.

   **Worked example — how the committee's analysis differs across the four intent types.** The committee pulls the per-model decision-count distribution stratified by `audit.deployment.intent` for the period. Suppose the chain shows two `gen_ai.response.model` values for a single decision-class (e.g., loan-application advisory). The four intent values produce four different committee responses:

   - `intent = ab_test` with a populated `audit.deployment.experiment_id`. The committee treats this as deliberate model-validation activity. It reviews the experiment's design document, the per-cohort decision distribution, the statistical-power evidence, and the experiment's relationship to the institution's SR 11-7 effective-challenge program. The two response-model values are expected — they are the experiment's arms.
   - `intent = canary` with a populated `audit.deployment.canary_traffic_pct`. The committee treats this as bounded production-validation. It reviews the canary's traffic-percentage trajectory over the period, the canary's decision-equivalence record against the production version, and the rollout/rollback decisions the institution made. The two response-model values are the production version + the canary; the committee documents the promotion-or-retraction outcome.
   - `intent = multi_region_drift` with a populated `audit.deployment.region`. The committee treats this as operational housekeeping. It reviews the regional-version drift alongside the institution's regional-config audit (audit-procedures.md P-26 extended) and confirms the drift is intentional (e.g., a phased regional rollout) rather than incidental.
   - `intent = vendor_reroute_observed`. The committee treats this as a control-completeness gap. It reviews the institution's vendor-management posture: did the contract specify model-version-change notification? Did the institution's detection logic surface the reroute through chain evidence rather than through vendor notification? An elevated `vendor_reroute_observed` count in the period triggers a vendor-management escalation.

   **Committee decision points (when each finding triggers escalation).** The committee's escalation rules attach to the §4.4.2 intent values:

   - **`audit.deployment.intent = vendor_reroute_observed` triggers a vendor-management escalation when:** (a) the period's `vendor_reroute_observed` count exceeds the institution's documented expected baseline (typically zero for vendors with contractual model-version-change notification, or a low single-digit count for vendors without); OR (b) the reroute pattern correlates with a specific time window (suggesting a vendor-side incident the vendor did not disclose); OR (c) the rerouted-to model has different SR 11-7 reproducibility characteristics than the contracted model (e.g., a different decoding-parameter envelope). The escalation goes to the institution's vendor-management committee with the working-paper from audit-procedures.md P-26 attached.
   - **`audit.deployment.intent = multi_region_drift` warrants a regional-config audit when:** (a) the regional drift was not part of a documented phased-rollout plan; OR (b) the drift persists across more than one period (transient drift during a rollout window is expected; persistent drift is a regional-config-management gap); OR (c) the drift correlates with elevated decision-disagreement rates between regions for the same decision-class (the regional model variants are producing materially different decisions). The regional-config audit goes to the institution's deployment-engineering team with the regional-version inventory and the per-region decision-distribution evidence from the chain.

   **Reference.** The full §4.4.2 schema, including the four-intent disposition table the committee uses, lives in `spec/chain-of-custody-v1.md` §4.4.2. The audit-procedure shape that produces the working-paper the committee reviews lives in `docs/audit-procedures.md` P-26 (extended for deployment-intent stratification).

## Specific scenarios

### When the chain reports a clean pass

MRM treats this as supporting evidence that the AI agent program operated under integrity-bearing controls during the period. The MRM Committee continues its substantive review of model behavior.

### When the chain reports an integrity failure

MRM coordinates with the chain-operations team and the institution's IR program. The MRM Committee may require additional model-validation work for the period covered by the failure (the AI decisions during that period are repudiable until the failure is investigated).

### When the chain detects an operational anomaly

MRM treats this as a supplementary signal. Operational anomalies are not directly model-risk events, but persistent anomalies may indicate underlying issues with the AI agent program's infrastructure that warrant attention.

## Cost summary

The chain has an operating cost dominated by HSM operations. The mid-size institution spends roughly $80k–$200k per year all-in. The detailed cost picture is in `docs/cost-model.md`.

## Decisions the MRM Committee may take

- Endorse the chain as part of the bank's model-risk framework (typically yes, given the integrity property)
- Set the seal cadence appropriate for the AI agent program's risk profile (the institution's choice within the spec defaults)
- Set the documentation expectations for chain-detected events (incorporate into the existing model-risk documentation)
- Set the cadence for MRM-Committee review of chain-output reports (typically aligned with the existing model-risk review cadence)

## Reference materials available to the committee

- The full chain spec at `spec/chain-of-custody-v1.md`
- The design docs at `docs/design/`
- The threat model at `docs/design/09-threat-model.md`
- The cost model at `docs/cost-model.md`
- The most recent verifier report (provided by the chain-operations team)
- Sample SOC report Section 4 description at `docs/soc-pack/section-4-template.md`

## Glossary (one paragraph)

The chain uses standard cryptographic terminology: HMAC (a keyed hash), Merkle tree (a hash-of-hashes structure that lets a single root summarize many leaves), HSM (a hardware security module that protects the signing key), Ed25519 (the signing algorithm). The full glossary is at `docs/design/10-glossary.md`. The MRM Committee does not need to memorize the cryptographic detail to evaluate the control; the chain-operations team and the SOC team handle that layer.
