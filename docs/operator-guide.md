# Operator guide

> **What this doc is.** Bank-operator-facing guide for running the ledger and SDK in production. Consolidates content scattered across the design docs into a single reference for the bank's IT-ops team.

## Audience

- The bank's chain-operations team
- The bank's HSM administration team
- The bank's IT-ops team that supports the AI agent platform

If you're an examiner, see `examiner-quickstart.md`. If you're on the SOC team, see `soc-pack/`. If you're on the MRM committee, see `MRM-COMMITTEE-BRIEF.md`.

## System overview

The implementation has three runtime components plus an HSM:

1. **The SDK** — runs in the bank's AI agent processes. Captures events, computes the chain hash, persists locally, exports over OTLP.
2. **The ledger server** — runs as a standalone service. Receives OTLP, re-verifies the chain, writes to the WAL, runs the daily seal job.
3. **The HSM** — holds the signing key. Receives signing requests from the seal job; returns signatures.
4. **The verifier** — the offline binary auditors run; not part of runtime.

## Deployment

### Topology

Pick one based on the institution's posture:

- **Self-hosted.** Bank runs everything in its own data center or cloud account.
- **BYOC (Bring-Your-Own-Cloud).** Vendor's image runs in the bank's cloud account. See `byoc-deployment.md`.
- **Vendor-hosted.** Vendor operates the entire stack on the bank's behalf, with per-tenant key isolation.

### Reference configuration

The ledger server config (`docs/design/06-ledger-server-design.md` §6) is YAML. Fill in:

- `listen` — gRPC and HTTP ports
- `tls` — server cert/key, client CA for mTLS
- `storage` — WAL DSN, hot DSN, cold S3 prefix
- `hsm` — PKCS#11 module path, slot, PIN env-var name
- `tenants` — per-tenant key labels and seal cadence

No defaults for credentials. No defaults for HSM PIN. Sourced from the institution's secret manager at process start.

### HSM provisioning

Pick HSM provider (`docs/design/04-hsm-custody.md` §3.2.1):

| Provider | FIPS 140-2 L3 tier |
|---|---|
| AWS CloudHSM Classic | Default |
| Azure Managed HSM | Default; Premium tier of Key Vault |
| Google Cloud HSM | Default |
| AWS KMS Custom Key Stores | Conformant when backed by CloudHSM |
| Azure Key Vault Standard | Not conformant — does not meet L3 |
| AWS KMS default tier | Not conformant — L2 only |

Per-tenant key labels per the institution's tenant naming scheme. Document the FIPS validation certificate number in the institution's control description.

### SDK deployment

The SDK runs in the AI agent process. Configuration:

- Tenant ID
- Master-key custodian endpoint (the handshake target)
- Workload identity credential (SPIFFE, mTLS cert, or HSM token)
- Local SQLite path with adequate disk
- OTLP target (the ledger server)
- TLS configuration (1.3 minimum; mTLS preferred)

Memory protection: configure `mlock`-or-equivalent per platform (`docs/design/02-chain-construction.md` §4.2). On capability-restricted environments, use HSM-backed key handles instead of in-memory keys.

## Tenant onboarding

Adding a new tenant to the chain is a multi-step choreography that crosses the institution, the HSM, the tenant key registry, and the regulator. The HSM-side mechanics are covered in `docs/design/04-hsm-custody.md`; the M&A and decommissioning transitions are covered in `docs/m-and-a-handoff.md`. This section is the end-to-end onboarding sequence — what the institution does, in what order, and what operational events are emitted at each step. New institutions standing up the chain follow it the first time; institutions adding tenants to an existing deployment follow the same sequence per tenant.

The sequencing rule is load-bearing. Regulator registration MUST complete BEFORE the institution captures the first chain event under the new key. The reasoning is mechanical: when the regulator examines the first event under that key, the regulator's fingerprint reference must already exist on the regulator's side, otherwise the fingerprint check the regulator runs at examination time has nothing to compare against. Institutions that reverse this order (capture first, register later) introduce a verification window during which the regulator cannot independently confirm the public key the institution claims, and that window is documented as an integrity-control finding regardless of whether the institution registers the key promptly afterward.

### Onboarding sequence

1. **Generate the IKM in tenant-controlled custody.** The institution generates the master HMAC key (IKM) in HSM-resident storage under the tenant's custody domain, 32 bytes minimum per spec §10.6. The IKM is non-extractable from the HSM in Model B deployments and bounded-extractable under documented memory-protection posture in Model A deployments (per spec §4.1.1). The institution captures `master_key.generated` as an operational event with the tenant identifier, the HSM provider, the FIPS validation certificate number, and the model (A or B).

2. **Generate the Ed25519 signing keypair in the HSM with non-extractable private key.** Per spec §10.5, the signing-key private half MUST NOT leave the HSM. The institution generates the keypair, records the public key, and binds the keypair to the tenant identifier in the HSM's key labelling scheme. The institution captures `signing_key.generated` with the tenant identifier, the key-label, the algorithm (`Ed25519`), and the HSM key-handle reference.

3. **Publish the public key to the tenant key registry.** The tenant key registry is the institution's authoritative store of `(tenant_id, key_version, public_key, fingerprint, valid_from, valid_until)` tuples. The institution publishes the new tenant's public key to the registry with `key_version=1`, the SHA-256 fingerprint of the public key, and `valid_from` set to the planned activation timestamp. The institution captures `tenant_key.registered` with the tenant identifier, key version, fingerprint, and registry reference.

4. **Register the public-key fingerprint with the regulator.** The institution submits the public-key fingerprint to the regulator using the regulator-fingerprint reception procedure documented in `docs/design/09-threat-model.md` §2.9. The regulator acknowledges receipt; the institution captures the acknowledgement reference. The institution captures `regulator_fingerprint.registration_acknowledged` with the tenant identifier, the fingerprint, the regulator's identity, and the regulator's acknowledgement reference. This event is the institution's evidence that step 4 completed before any chain event under the new key was captured. Retain the acknowledgement document for the institution's standard audit-document retention period.

5. **Provision tenant routing in the ledger config.** Update the ledger server config (per `docs/design/06-ledger-server-design.md` §6) with the new tenant's key labels, seal cadence, and routing rules. Restart or hot-reload the ledger server per the institution's change-management procedure. The institution captures `tenant.provisioned` with the tenant identifier, the seal cadence, the ledger version that received the configuration, and the change-management record reference.

6. **SDK handshakes for `key_version=1` IKM.** The institution starts AI agent processes for the new tenant. The SDK performs the handshake against the master-key custodian per spec §4.1.1 (Model A or Model B as the institution chose in step 1). The custodian returns either the IKM bytes (Model A) or the per-tenant session key derived inside the HSM (Model B). The SDK reports `tenant.handshake_completed` with the tenant identifier, the handshake model, and the key version. This event is institution-defined; the schema follows the existing operational-events vocabulary conventions.

7. **Capture the first chain event and confirm end-to-end verification.** The institution generates a synthetic test event (or the first production event) for the new tenant and confirms the ledger receives and re-verifies it without error. The institution runs the verifier against the day's events ending in the seal that includes the test event. A successful verifier run with PASS disposition is the onboarding-complete signal. The institution captures `tenant.first_event_verified` with the tenant identifier, the event run-id and sequence number, the seal that covered the event, and the verifier output reference.

The seven steps form one logical operation; the institution's runbook captures the sequence as a checklist with each operational event as a gating evidence artifact. Skipping or reordering steps creates audit gaps the SOC team will surface during the next CC8.1 review.

### Onboarding rollback

If onboarding fails partway — the regulator does not acknowledge the fingerprint registration in step 4, or the SDK handshake fails in step 6, or the first-event verification fails in step 7 — the institution does NOT proceed. The institution captures `tenant.onboarding_aborted` with the failure step and the reason, retires any partially-published artifacts (deactivate the tenant in the key registry, remove the tenant's routing from the ledger config), and re-runs onboarding from step 1 after the failure cause is corrected. The institution MUST NOT capture production chain events for a tenant whose onboarding has not completed through step 7.

## Daily operations

### Health monitoring

Standard probes on the ledger admin port (default `:4319`):

- `/healthz` — process is up
- `/readyz` — all subsystems initialized, HSM session established, storage healthy
- `/metrics` — Prometheus-format metrics

Key metrics to watch:

| Metric | Healthy range |
|---|---|
| `ledger_otlp_received_total{status="ok"}` | Tracks expected event volume |
| `ledger_chain_verifications_total{result="pass"}` | All ingest passes; a fail is an alert |
| `ledger_seal_runs_total{status="success"}` | Daily count of successful seals |
| `ledger_seal_age_seconds` | < 25 hours under normal operation; > 72 hours triggers regulator notification |
| `ledger_hsm_operations_total{status="error"}` | Should be near zero; spikes indicate HSM issues |

### Daily seal job

Default trigger: UTC 00:00 + 60 minutes per tenant. The institution monitors the seal-completion event (`seal.job_completed`) and alerts on `seal.job_failed`.

For incident response, the seal job can be triggered on demand:

```
ledger-cli seal trigger --tenant tenant_acme_prod --date 2026-04-01
```

(Or HTTP POST to the admin port with appropriate authentication.)

### Cadence-aware seal-publication SLA (per spec §4.3)

Spec §4.3 makes seal publication a normative MUST: the signed root MUST be appended within 60 minutes of the END of the tenant-day's seal window. The 60-minute number is the same across cadences; what differs is the reference moment (the "end of the seal window") that the 60 minutes runs from. The operations team's monitoring posture varies by cadence.

**Daily cadence (default).** The seal job runs at UTC 00:00 + 60 minutes for the previous UTC day. Operations teams alert if `signed_at` for any tenant-day extends beyond `01:00 UTC` of the day after the seal day. The monitoring posture is the steady-state `ledger_seal_age_seconds` threshold under 25 hours; alerts fire on the H+1:00 UTC SLA breach.

**Hourly cadence (per-tenant configurable per §4.2.1).** A seal covering hour H must be signed and appended by `H+1:00 UTC` — a seal covering 13:00–14:00 UTC must be signed by 15:00 UTC. The operations team monitors per-hour seal-completion timestamps against the H+1:00 UTC SLA. The monitoring threshold is tighter than daily — `ledger_seal_age_seconds` should not exceed roughly 75 minutes (60-minute SLA plus 15-minute alerting margin). Hourly cadence usually accompanies a high-throughput tenant where the per-hour seal volume is part of the operational rhythm; the institution's runbook documents the per-hour expectation and the alert thresholds.

**Weekly cadence (institution-approved relaxation per §4.2.1).** The signed root is due by `01:00 UTC` on the day after the seal-week's end. The default seal-week ends Monday (the seal covers Tuesday-through-Monday); institutions MAY declare a different week-end day in their CC8.1 control description, in which case the SLA shifts to 01:00 UTC on the day following the declared week-end. Operations teams align their monitoring to the institution's declared week-end day. For example, an institution declaring Friday as the seal-week-end has each weekly seal due by 01:00 UTC Saturday; alerts fire on the 01:00 UTC Saturday SLA breach. The CC8.1 declaration is the load-bearing document — the verifier confirms the seal record's `cadence` field is `"weekly"` and the SLA monitoring confirms the institution's actual publication rhythm matches the declared week-end day.

The cadence value appears on every seal record per spec §4.2.1; the verifier's cadence check (spec §7 step 12) confirms the seal record's cadence matches the institution's claim. Spec §4.3's `sign_payload` extension binds `cadence` under the HSM signature — a cadence rewrite (e.g., flipping `daily` to `weekly` to claim a relaxed posture) is a forge attempt that requires forging the HSM signature.

### Multi-region operational guidance (per spec §10.15)

Spec §10.15 normates two conformant patterns for multi-region resilience. The institution selects per tenant and documents the choice in CC8.1. Operations teams supporting a multi-region deployment configure cross-region replication and per-region reconciliation per the chosen pattern.

**Pattern A — Active-active with seal-region pinning (RECOMMENDED).** A single canonical `tenant_id` operates in multiple regions; one region is the **seal region** for the tenant; the others are **replication regions** that ship their events to the seal region before seal-time.

- **Seal-region designation per tenant.** The institution's CC8.1 control description names the seal region per tenant. One region per `tenant_id`. The designation is the load-bearing document — the seal region is where the day's Merkle root is computed and where the HSM signature is produced.
- **Cross-region replication mechanism.** Institution chooses the mechanism that fits its data-platform posture: Postgres streaming replication, Kafka cross-region, S3 cross-region replication, or application-level event-streaming. The mechanism is opaque to the chain spec; what matters is that source-region events land at the seal region's ledger before seal-time.
- **Replication-completion SLA.** Events captured in a source region MUST arrive at the seal region before seal-time. The SLA aligns with the spec §4.3 publish window (60 minutes after seal-window end per spec §4.3 / §10.15 invariant 5). A replication mechanism whose worst-case lag exceeds the publish window is non-conformant for Pattern A — the institution either tightens the lag, accepts the seal-record will exclude late-arriving events (which become next-day's events under the seal region's ingest clock per spec §4.2.2), or operates Pattern B instead.
- **Per-region event-count reconciliation.** Each region emits one `master.cross_region_replication_completed` operational event per `(tenant_id, seal_date)` per `soc-pack/control-evidence-events.md`. The seal region's actual event count for the tenant-day MUST equal the sum of `events_replicated_to_seal_region` across all source regions. Operations teams monitor the equality alongside the standard seal-completion monitoring; a non-matching aggregation routes to the institution's CC8.1 multi-region replication procedure.
- **Failover.** If the seal region becomes unavailable before seal-time, the institution promotes a replication region per the documented CC8.1 failover procedure. The promoted region MUST have all events for the tenant-day before producing the seal — operations teams confirm the per-region replication evidence completed for every source region before the promoted region runs its seal job.
- **Run-locality (normative for v1.0).** Workloads MUST be region-pinned in v1.0: a run starts and ends in one region. Cross-region run continuation (an agent migrating mid-run between regions) is deferred to v1.1. Workloads requiring cross-region run continuation either route to a single region or operate under Pattern B.

**Pattern B — Per-region `tenant_id` (CONFORMANT alternative).** Each region operates an independent chain under its own per-region `tenant_id`.

- **Per-region tenants.** Each regional tenant has its own IKM, its own seals, and its own verifier runs. Cross-region correlation is institution-side only.
- **Cross-region correlation registry.** The institution maintains an institution-side registry naming the regional tenants comprising one logical deployment (typical: `tenant_acme_prod_us_east_1`, `tenant_acme_prod_us_west_2`, `tenant_acme_prod_eu_west_1` all roll up to the institution's "Acme prod" deployment). Operations teams reference the registry during human disambiguation — examiner inquiries, customer disputes, or cross-region operational reviews.
- **Verifier runs.** The verifier runs once per regional tenant per audit period. The institution's audit-period burden is O(regions) rather than O(1); operations teams plan verifier-run capacity accordingly.

**Pattern selection guidance.** Pattern A is the lower-cost option for institutions whose risk posture admits cross-region replication trust — the verifier runs once per tenant per audit period, and the seal-region's signed root aggregates events from every region. Pattern B is the right choice when regional regulatory regimes mandate in-region key custody (some EU banking jurisdictions; APAC data-sovereignty regimes), or when the institution's risk posture treats cross-region replication as an unacceptable trust boundary (the seal region's compromise would corrupt the seal even though events were captured securely elsewhere). The patterns are mutually exclusive per tenant — an institution operating Pattern A for one tenant MAY operate Pattern B for another.

Cross-reference: spec §10.15; `docs/design/00-overview.md` §6.4; `docs/dr-and-resilience.md` "Multi-region resilience"; `docs/soc-pack/control-evidence-events.md` "Cross-region replication reconciliation"; `docs/audit-procedures.md` P-37.

### OTLP collector configuration for chain-of-custody traffic (per spec §4.4.3 and §4.4.4)

OTel collectors that ingest chain-of-custody traffic alongside regular telemetry need a configuration that exempts chain traffic from severity-based filters and severity-based sampling. The need is mechanical: OTel collectors ship with a default posture that drops anything below INFO or below WARN depending on the deployment. A chain record dropped by a severity filter is a silent integrity gap — the verifier reports PASS on what arrived at the ledger, not on what the SDK tried to send. Spec §4.4.4 closes this with normative rules; this section translates those rules into operator-facing configuration guidance.

**Required.** The collector configuration MUST exempt chain-of-custody traffic from severity filters per spec §4.4.4. Chain traffic is identified at the collector layer by the Resource attribute `ffiec.chain.spec` (per spec §4.4.3). A collector that applies severity filters or severity-based sampling to records carrying this Resource attribute is non-conformant and produces a control-completeness gap routed through the institution's CC8.1 procedure when the SOC team's P-38 procedure samples it.

**Recommended.** Route chain-of-custody traffic through a dedicated pipeline branch with NO severity filter and NO sampling processor. The pattern is: one receiver accepts all OTLP traffic, a routing connector splits chain traffic from regular telemetry by inspecting the `ffiec.chain.spec` Resource attribute, and the two pipelines apply different processing — the regular pipeline retains the institution's standard severity filtering and tail-sampling posture, the chain pipeline applies only batching. The split lets operators continue applying normal cost-control filters to non-chain telemetry without affecting chain integrity.

**Receiver-side stamping.** The Herald.Compliance receiver positions the `SeverityNumber` per chain record using its `QuickLogBuilder` resolver — a Herald-side component that picks a value within the spec's normative range `9 ≤ N ≤ 20` (INFO floor through ERROR4 ceiling, just below `FATAL = 21`). The `SeverityText` is `"OTLP"`. The resolver tunes the level per ingest based on the institution's policy; higher values (closer to 20) resist routine `< WARN` and `< ERROR` filters more aggressively, while staying below `FATAL = 21` keeps routine alerting infrastructure from treating the record as a system-fatal alert. Operators configuring downstream filters or aggregations within Herald.Compliance MUST exempt records with `SeverityText = "OTLP"` from filtering and from any aggregation that drops based on severity. The text is unique enough to grep for in log analysis tools — operators searching their SIEM for chain records use `SeverityText == "OTLP"` to pull every chain entry without parsing OTel attribute payloads.

Institutions operating their own receiver MAY use a different resolver and a different `SeverityText` provided the produced `SeverityNumber` lands in the `9..20` range and the resolver mechanism plus the chosen text are documented in CC8.1. The Herald reference values are the default; institution-side deviations are governed by the institution's CC8.1.

**Worked example — OTel Collector YAML.** The pattern below routes traffic by the `ffiec.chain.spec` Resource attribute. The chain pipeline applies only batching; the regular pipeline applies the institution's standard severity filter and tail sampling. Institutions running additional postures (HIPAA, PCI) extend the routing table with `ffiec.chain.posture` predicates and add per-posture pipelines.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  # Routing connector — splits chain traffic from regular telemetry
  # by inspecting the ffiec.chain.spec Resource attribute.
  routing:
    default_pipelines: [traces/regular]
    table:
      - statement: 'route() where resource.attributes["ffiec.chain.spec"] != nil'
        pipelines: [traces/chain]

  # Severity filter — applied to the regular pipeline only.
  # MUST NOT appear in the chain pipeline per spec §4.4.4.
  filter/severity:
    error_mode: ignore
    logs:
      log_record:
        - 'severity_number < SEVERITY_NUMBER_INFO'

  # Tail sampling — applied to the regular pipeline only.
  tail_sampling:
    policies:
      - name: sample_errors
        type: status_code
        status_code: { status_codes: [ERROR] }
      - name: sample_slow
        type: latency
        latency: { threshold_ms: 1000 }

  batch:
    timeout: 5s

exporters:
  # Chain-of-custody exporter — ships to the institution's ledger.
  otlphttp/chain:
    endpoint: https://ledger.example-bank.internal/v1/traces
    tls:
      insecure: false

  # Regular telemetry exporter — ships to the institution's APM backend.
  otlphttp/regular:
    endpoint: https://apm.example-bank.internal/v1/traces

service:
  pipelines:
    # Chain-of-custody pipeline. NO severity filter. NO sampling. Batching only.
    traces/chain:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp/chain]

    # Regular telemetry pipeline. Standard severity filter and tail sampling.
    traces/regular:
      receivers: [otlp]
      processors: [filter/severity, tail_sampling, batch]
      exporters: [otlphttp/regular]
```

The collector config is the load-bearing artifact the SOC team's P-38 procedure samples. Operators store the config in source control alongside the institution's other infrastructure-as-code artifacts; changes to the chain pipeline branch route through the institution's standard change-management procedure.

#### FAQ — Why does Herald stamp `OTLP` as the severity text?

The `OTLP` severity text serves two operational purposes that a default `INFO` text does not.

First, it identifies chain-of-custody traffic in operator-visible logs and dashboards without requiring operators to know the OTel attribute schema. An operator looking at a SIEM dashboard sees a column labelled `severity_text` and can filter on `OTLP` to pull every chain record. The same filter does not need to know about `ffiec.chain.*` attributes; it works with whatever log-analysis tool the institution operates.

Second, it resists routine severity-based filters. Most OTel collector deployments ship with a default that drops `< INFO` (severity number < 9) records, and institutions running cost-control filters often raise the floor to `< WARN` (severity number < 13) or `< ERROR` (severity number < 17). Herald's `QuickLogBuilder` positions chain records in the `9..20` range mandated by spec §4.4.4 with the specific value tuned per institution policy. Higher positioning resists more filters; staying below `FATAL = 21` keeps routine alerting infrastructure from misinterpreting the record. The stamp does not replace the §4.4.4 collector pass-through rule — the rule is the load-bearing requirement — but it provides defense-in-depth for downstream consumers that have not yet been updated for the rule.

The text "OTLP" is short, unique enough to grep for, and operator-meaningful. It does not collide with the standard OTel severity texts (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`) and is not a number, so it stands out in log analysis tools that group by severity text.

Cross-reference: spec §4.4.3, §4.4.4; `docs/design/05-otlp-wire.md` §"Transport identification" and §"Severity treatment for chain-of-custody traffic"; `docs/audit-procedures.md` P-38.

### Receiver-policy discovery (Herald-specific)

The chain-of-custody-v1 spec does NOT mandate a receiver-policy discovery endpoint. The configuration described here applies to the Herald reference topology, where the SDK and the chain receiver run as distinct processes. Institutions running the receiver as a library inside the SDK process can skip this section; institutions running the Herald.Compliance receiver alongside the SDK use this section to wire the SDK's discovery client.

The endpoint exists so the SDK can ask the receiver, before emitting traffic for a tenant, what the receiver will actually do with that traffic — which `SeverityNumber` value `QuickLogBuilder` positions for this tenant, which filter exemptions are configured, which `sign_payload_version` forms are accepted. Operators see the spec range alongside the actual position the receiver picks; SDKs in strict mode refuse to emit traffic when the receiver disagrees with what the SDK is configured to send. The endpoint inherits the OTLP transport-security floor in spec §5.1: TLS 1.3 minimum, server-authenticated TLS, Bearer token or mTLS client auth.

**SDK configuration.** Operators wire the SDK with the receiver-policy endpoint URL plus an auth credential. The credential comes from the same source that holds the institution's other backend credentials — the credential vault, the secrets manager, or the operator's configuration management system. Two auth modes:

- **Bearer token.** A long-lived or rotating token issued by the institution's auth service. The SDK sends it as `Authorization: Bearer <token>` on every fetch. The token is never logged; the SDK refuses to construct against a plaintext `http://` URL so a misconfigured scheme cannot leak the token over plain HTTP.
- **mTLS.** A client certificate plus key issued by the institution's PKI. The SDK loads the cert chain into its SSL context and presents it on the TLS handshake. No `Authorization` header is set; auth rides on the certificate. Suitable for high-assurance deployments where the institution's posture mandates mTLS for control-plane traffic.

The SDK caches the fetched policy with a short TTL (15 minutes by default) and uses `If-None-Match: <etag>` for conditional refresh. Operators changing receiver-side policy should expect the change to propagate to the SDK fleet within one TTL window; institutions running large fleets may want to lower the TTL or trigger an explicit refresh from the admin tooling.

**Charles / Postman observability.** During development, operators may want to inspect the receiver-policy fetch flow. The SDK's stdlib `urllib` transport honours the standard `HTTP_PROXY` / `HTTPS_PROXY` environment variables, so pointing the SDK at a Charles or Postman proxy is straightforward:

```bash
export HTTPS_PROXY=http://127.0.0.1:8888  # Charles default
# Install the Charles root CA on the developer's machine and either
# set ca_bundle_path to the Charles CA OR temporarily disable cert
# verification in the SDK config (developer machines only).
```

**Charles MITM warning.** Charles intercepts HTTPS traffic by terminating the SDK's TLS connection and presenting its own self-signed cert. The developer must explicitly install the Charles root CA on their machine for the SDK to trust the proxy. Production traffic is HTTPS-only with the institution's PKI; the Charles cert is never installed in production. Operators MUST NOT ship a configuration with `ca_bundle_path` pointing at the Charles CA into a production deployment — this is the same hygiene rule that applies to any developer-tooling cert. The SDK refuses plaintext `http://` URLs at construction time, but it cannot detect a misconfigured CA bundle; that is on the operator's change-management procedure.

Postman behaves the same way; operators using Postman's proxy feature install the Postman CA on the developer machine and follow the same production-hygiene rule.

**Failure modes.** When the receiver-policy fetch fails (network error, 5xx, malformed JSON), the SDK's behaviour depends on whether a cached policy exists:

- **Cached policy available.** The SDK logs a WARN (`herald: receiver-policy fetch failed for <tenant>: <reason> — using cached entry`) and returns the cached policy. Chain emit continues uninterrupted. The cache freshness clock is NOT bumped, so the next emit re-tries the fetch.
- **No cached policy.** The SDK raises `ReceiverPolicyUnavailable`. The fail-open / fail-closed choice is made by the institution's SDK configuration:
  - **Fail-open.** The SDK catches the unavailable exception, logs WARN, and emits traffic without a policy check. Suitable for institutions where chain emit availability matters more than strict policy alignment.
  - **Fail-closed.** The SDK propagates the exception and refuses to emit traffic to the affected tenant. Suitable for high-assurance deployments where a misconfigured receiver shouldn't see traffic. Set `ReceiverPolicyClient(strict=True, ...)` in the SDK configuration.

The fail-mode choice is documented in CC8.1 alongside the rest of the institution's chain-of-custody control description.

**Operator inspection.** Operators diagnosing a receiver-policy issue can fetch the policy manually with `curl`. The SDK's auth credential rides on the same headers the operator passes to `curl`:

```bash
# Bearer token
curl -H "Authorization: Bearer ${HERALD_RECEIVER_POLICY_TOKEN}" \
     -H "Accept: application/json" \
     --cacert /etc/herald/receiver-ca.pem \
     https://compliance.bank.internal/api/v1/tenants/acme/receiver-policy

# mTLS
curl --cert /etc/herald/client.crt \
     --key  /etc/herald/client.key \
     --cacert /etc/herald/receiver-ca.pem \
     -H "Accept: application/json" \
     https://compliance.bank.internal/api/v1/tenants/acme/receiver-policy
```

The response body matches the SDK's parsed shape one-to-one — operators comparing the curl output to the SDK's logged policy can confirm the SDK is fetching what the receiver is serving. The `actual_value` field on the response is the value `QuickLogBuilder` reports for this tenant; operators see both the spec range `[9, 20]` and the actual position the receiver picked.

Cross-reference: spec §4 (implementation topology); spec §4.4.4 (severity range); spec §5.1 (transport security); `docs/design/05-otlp-wire.md` §4.6 "Receiver-policy discovery" for the design rationale and endpoint shape.

### `chain_kind` operational note (per spec §3)

Spec §3 normates a closed enumeration for the `chain_kind` field on every chain entry. Operations teams will see six values in production traffic:

| `chain_kind` | When operations sees it |
|---|---|
| `audit` | Default — application audit event the institution emits |
| `model_call` | Chain entry representing an LLM invocation |
| `tool_call` | Chain entry representing a tool invocation |
| `routing` | Chain entry from the routing-event surface (per spec §4.4.1) |
| `translation` | Chain entry from the ECOA translation step (per spec §10.11) |
| `operational` | Chain entry for control-evidence operational events |

The `operational` value primarily appears on control-evidence events (operational events recording reconciliation, key-rotation observation, mirror-reconciliation, and similar control-evidence material — see `soc-pack/control-evidence-events.md`). These are operational events under spec §10.2 and not chain-of-custody integrity events; the `chain_kind` distinction lets dashboards, alerting rules, and SOC sample-comparison procedures separate operational evidence from integrity-bearing chain content. The verifier MUST reject any `chain_kind` value not in the enumerated set with `chain_kind out of v1 enumeration at seq N` per spec §3 — operations teams that see verifier-reported `chain_kind out of v1 enumeration` failures route the SDK or chain decorator under investigation; the value is closed and a non-conforming value indicates the producer is non-conformant.

### Operational events

The implementation emits operational events per `soc-pack/control-evidence-events.md`. The institution routes these to:

- The bank's SIEM for security monitoring
- The bank's observability backend (Splunk, Datadog) for operations
- A long-term archive for SOC and examination evidence (retention matches the events table)

### Backup and recovery

The WAL is the source of truth (`docs/design/06-ledger-server-design.md` §3.1). Backups are continuous via:

- PostgreSQL streaming replication (preferred; sub-second RPO)
- S3 versioned bucket replication (for cold-store archives)

Recovery procedure:

1. Restore the WAL from the most recent backup
2. Re-derive the hot store and indexes from the WAL
3. If the backup is older than the most recent seal, replay the gap from any available source (SDK local SQLite buffers, downstream OTLP backends)
4. Re-run the verifier to confirm integrity

If the gap cannot be filled, the affected days are unverifiable; treat as an integrity-control failure per `incident-response-playbook.md`.

## Periodic operations

### Quarterly

- Rotate HSM PIN (`docs/design/04-hsm-custody.md` §5.2.1)
- Review operational-event log retention
- Review seal-age metric for any near-threshold events

### Annually

- Re-evaluate seal cadence appropriateness for the institution's AI agent risk profile
- Update vendor SOC report on file
- Disaster recovery exercise (table-top or actual failover)
- Review and update the IR playbook

### As needed

- Master-key rotation (institution-defined cadence)
- Spec version migration (when a new spec version ships)
- Verifier release adoption (within the project's 30-day pre-announcement window)

## Master-key rotation procedure

1. Notify the chain-operations team and the institution's MRM committee
2. Generate a new master at the master-key custodian
3. Mark the new master as the primary
4. Force application processes to re-handshake (rolling restart)
5. Verify new events carry the new `master_version`
6. After all processes have re-handshaked, retire the old master from the active key set
7. Retain the old master in the registry with valid_from/valid_until for verification of past events
8. Document the rotation in the institution's control-evidence repository

The seal record automatically captures `master_version`; the verifier resolves which master to apply per event.

## Examination-time master-key handover procedure

Spec §7 splits verifier behavior into two paths. Without `--master-key`, the verifier confirms structural consistency only — Merkle roots match, signatures validate, the chain is well-formed. With `--master-key`, the verifier additionally re-derives every per-event MAC from the IKM and confirms each entry's `payload_hash` MAC matches the chain's claimed value. The second path is the load-bearing per-event integrity check the spec exists to provide. Without it, the verifier reports `structurally consistent, key-bound verification skipped` and the examiner has not actually tested per-event integrity — only the seal-and-signature envelope.

Round-12 J-1 (Training Module 3 invocation missing `--master-key`) signaled that this distinction confuses new examiners. The institution's responsibility is to make the master-key handover procedure explicit, repeatable, and auditable so the examiner gets the full verification path by default rather than the structural-only path by accident.

This section covers the handover for an in-person or remote-supervised examination. Court-ordered disclosure and law-enforcement disclosure are separate concerns covered in `docs/legal-disclosure.md`; this is the routine examination case.

### Conformant handover channels

The institution chooses one of three approaches and documents the choice in the institution's control description. All three preserve the IKM's confidentiality and produce auditable handover and disposal evidence.

**Approach A — Physical secure-courier of the IKM file on hardware-encrypted media.** The institution exports the IKM (under appropriate HSM authorization with a documented key-export ceremony for institutions whose IKM is HSM-resident) to a hardware-encrypted USB device or equivalent FIPS-validated medium. The medium is sealed in tamper-evident packaging and delivered to the examiner's secure-storage location via the institution's secure-courier service. The examiner returns the medium at the end of the examination. The institution destroys the medium per its standard sensitive-media disposal procedure and captures the destruction evidence. This approach is appropriate for Model A institutions whose IKM is exportable under documented procedure; Model B institutions whose IKM is non-extractable cannot use this approach.

**Approach B — HSM-side in-place key wrap with examiner-bounded access.** The institution does not export the IKM. Instead, the institution provisions a per-examination HSM operator role with `decrypt`-via-HSM permission for the affected tenant keys, scoped to the examiner's identity for the duration of the examination. The examiner runs the verifier with `--master-key-via-hsm` (or equivalent) which routes MAC re-derivation through the HSM's HMAC primitive rather than computing the MAC in verifier process memory. This approach is appropriate for Model B institutions and for Model A institutions whose risk posture forbids IKM bytes leaving the HSM. The verifier runtime is higher (each MAC re-derivation incurs an HSM round-trip), which is a planning consideration for tier-1 examination windows.

**Approach C — Per-examination bounded operator role.** The institution provisions a per-examination service account with `master-key.read` permission, scoped to the examiner's identity, the affected tenant identifiers, and the examination period. The examiner runs the verifier with the service account's credentials. The bounded role expires at a documented end-of-examination timestamp and the institution revokes it explicitly. This approach is appropriate for institutions whose master-key custodian is a software service rather than an HSM and whose risk posture supports time-bounded credential issuance to external parties under contractual examiner identity.

The institution MAY support multiple approaches and choose per examination. The institution MUST document the chosen approach in the examination's evidence package so the examiner can confirm the handover proceeded under named procedure rather than ad-hoc arrangement.

### Receipt and disposal evidence

The institution captures two operational events bracketing every examination master-key handover. Both events are institution-defined under the operational-events vocabulary in `soc-pack/control-evidence-events.md`; institutions adding the events to their event stream conform to the existing schema (event name, timestamp, tenant_id, correlation_id, fields).

**Receipt event.** Captured when the examiner receives access. Schema:

```json
{
  "event": "master_key.examiner_handover_received",
  "timestamp": "2026-06-15T09:00:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "exam-2026-06-15-occ",
  "fields": {
    "examiner_identity": "OCC-IT-EX-Cert-2026",
    "examination_period_start": "2026-01-01",
    "examination_period_end": "2026-06-14",
    "handover_channel": "approach_b_hsm_in_place_wrap",
    "handover_authorization_record_id": "CM-2026-06-15-009",
    "expires_at": "2026-06-22T17:00:00Z"
  }
}
```

**Disposal event.** Captured when the examiner returns access (returns the medium under Approach A, the institution revokes the role under Approaches B and C). Schema:

```json
{
  "event": "master_key.examiner_handover_returned",
  "timestamp": "2026-06-22T15:30:00Z",
  "tenant_id": "tenant_acme_prod_us_east_1",
  "correlation_id": "exam-2026-06-15-occ",
  "fields": {
    "returned_at": "2026-06-22T15:30:00Z",
    "disposal_method": "hsm_role_revoked",
    "disposal_evidence_record_id": "CM-2026-06-22-014",
    "examination_outcome_reference": "exam-report-2026-06-occ"
  }
}
```

The two events bracket every handover. The SOC team confirms one-to-one pairing during CC8.1 review; an unpaired receipt event (no matching return event) is investigated as a potentially-unrevoked examiner credential.

### Control-description language

The institution adapts the following text for its control description. The text is illustrative; institution-specific identifiers, approach choice, and procedural references replace the bracketed values.

> **Control.** Examination-time master-key handover.
>
> **Description.** When a regulatory examiner requires `--master-key` access to perform full per-event MAC verification under spec §7, the institution provides access through Approach `[A | B | C]` per the institution's chosen handover channel. Receipt evidence is captured as `master_key.examiner_handover_received` with the examiner identity, examination period, handover channel, and expiration timestamp. Disposal evidence is captured as `master_key.examiner_handover_returned` when access is revoked at examination conclusion. Both events are routed to the institution's SOC evidence archive and retained for the standard audit-document retention period.
>
> **Frequency.** Per examination. Typical cadence is 12–18 months for community banks and quarterly-to-continuous for tier-1 banks under heightened supervision.
>
> **Owner.** The chain-operations team in coordination with the institution's regulatory-relations function.

### IR procedure for mishandled handover

If the institution detects that an examination master-key handover has been mishandled — the courier medium is lost in transit (Approach A), the examiner-bounded credential is leaked or used outside the examination scope (Approach C), the HSM role remains active beyond the documented expiration (Approach B), or the examiner retains access beyond the agreed window — the institution treats the event as an IR incident.

The handover-mishandling IR procedure:

1. Activate IR playbook scenario 4 (master-key compromise response) treating the handover-channel exposure as a potential master-key compromise. The procedural rigor is the same regardless of whether the institution believes the compromise actually occurred.
2. Capture `master_key.examiner_handover_compromise_suspected` as an institution-defined operational event with the suspected compromise vector, the affected tenant identifiers, and the time window of suspected exposure.
3. Notify the regulator within 36 hours per the cyber-incident notification rule. The notification names the handover-channel exposure as the trigger; the regulator MAY classify the event as a notifiable incident or MAY accept it as a procedural anomaly with no chain-integrity impact, depending on the facts.
4. Rotate the affected master keys per the master-key rotation procedure. The rotation is precautionary; the historical chain remains verifiable under the previous IKM, the new IKM covers events captured after rotation.
5. Capture `master_key.examiner_handover_remediation_completed` with the remediation evidence and the regulator's classification.
6. Conduct a post-incident review and update the handover procedure if the incident reveals a procedural gap. The review's findings flow into the institution's CC8.1 control description.

The institution's IR playbook MAY add a handover-specific scenario referencing this procedure. The SOC team includes handover-mishandling readiness in CC8.1 review by sampling examination handover records and confirming the receipt and disposal events are paired and the disposal evidence is retained.

## Forensic handoff procedure

The examination handover procedure above covers the routine regulatory-examination case. A different handoff occurs when the institution gives chain evidence to a forensic examiner — law enforcement under a warrant, defense counsel preparing for litigation, a neutral third-party expert appointed by the court. The forensic posture is more disciplined than the examination posture because the receiving party is not the institution's regulator; the chain-of-custody documentation becomes the foundation for the examiner's later courtroom testimony that the chain was received in integrity state.

The forensic handoff is a one-time transfer for a specific case. The institution documents which handoff method was used and why so the receiving examiner can testify cleanly under cross-examination.

**Required handoff documentation.** Each forensic handoff produces:

1. A chain-of-custody form signed by both parties. The form names the IKM (or the HSM endpoint, if HSM-mediated) being transferred, the SHA-256 hash of the key material, the date, the time, and the signatures of the institution's transferring custodian and the receiving examiner. The hash lets the receiving examiner confirm later that the IKM bytes have not changed since transfer.
2. A two-person witness to the transfer. The witness signs the chain-of-custody form alongside the transferring and receiving parties. The two-person rule prevents a later challenge that the transfer was unwitnessed.
3. A receipt confirmation. The receiving examiner signs that the sealed evidence was received in integrity state. If the seal was broken or the medium was damaged, the receipt names the discrepancy and the institution treats the handoff as compromised under the IR playbook's mishandled-handover procedure.

**Method 1 — Physical IKM transfer in tamper-evident envelope.** The institution exports the IKM to hardware-encrypted media per the examination handover Approach A procedure, places the media in a tamper-evident envelope, and transfers physical custody to the receiving examiner under the chain-of-custody form. This method is appropriate when the receiving party cannot use an HSM-mediated path — for example, defense-side examiners working from a laboratory without HSM connectivity to the institution.

**Method 2 — HSM-mediated verification (preferred for defense-side examiner handoff).** The institution provisions an HSM API endpoint with documented credentials scoped to the examiner's identity for the duration of the engagement. The examiner's verification tool queries the HSM without ever holding the IKM bytes. This method is preferred for evidence handoff to a defense-side examiner because the examiner can still perform full per-event MAC verification per spec §7 without taking custody of key material; the institution retains operational control of the IKM throughout. The chain-of-custody form names the HSM endpoint, the credential identifier, the credential expiration, and the receiving examiner's identity rather than a physical artifact.

Either method is acceptable. The institution documents which method was used and why in the case file. The handoff documentation is the foundation for the examiner's later testimony that the chain was received in integrity state and remained unaltered throughout the examination. Without this documentation, opposing counsel argues the chain's integrity was lost during transfer; with it, the examiner testifies cleanly from the chain-of-custody form.

**Cross-reference.** `docs/litigation-support.md` (forensic handoff in the litigation context, examiner-side qualification), `docs/incident-response-playbook.md` (mishandled handover IR procedure), spec §7 (verifier paths — structural-only versus full-MAC verification).

## Litigation-hold and subpoena-response timing

When the institution receives notice of litigation, anticipates a subpoena, or learns a customer dispute is likely to escalate to litigation, the IT team triggers a litigation hold on the chain data. The hold is a discrete operational shift — the institution's standard retention schedule is paused for affected data, IKM rotation is frozen or carefully documented, and access to the affected period's chain entries is logged at every touch. The institution's discovery posture is the standard FRCP 34 deadline (14 days after service of the request, or as extended by court order); the litigation-hold operations begin earlier — at the moment of anticipation — so the institution has the evidence ready when the formal request arrives.

**Required hold operations.** When the institution triggers the hold, the chain-operations team:

1. Extends chain-data retention beyond the standard schedule for affected tenants and date ranges. The extension covers chain entries, seal records, IKM history, verifier output, and operational events.
2. Freezes IKM rotation, OR documents the rotation with explicit before-and-after IKM custody records. If a rotation is operationally required during the hold (a scheduled rotation that cannot be deferred without breaking the institution's standard rotation cadence), the institution captures the IKM custody chain for both the pre-rotation and post-rotation IKM so the chain remains verifiable across the rotation boundary.
3. Prevents any deletion or modification of chain entries or seal records for the affected period. The institution's storage controls (append-only enforcement per spec §10.3) already prevent modification; the hold makes deletion explicit by documenting that any retention-job action against the affected data requires legal-team approval.
4. Preserves verifier output for the affected period. The "before" verifier output documents the chain's integrity state at the time of the dispute; preservation prevents an opposing party from claiming the institution generated post-hoc verifier output that does not match the chain's actual state.

**Subpoena-response package.** When the formal request arrives, the institution produces:

- Chain entries for the affected customer ID and date range. Production is in native NDJSON form per the eDiscovery posture (one JSON object per line per spec §6).
- Full verifier output covering the period, showing PASS or documenting any failures.
- Metadata about the chain deployment during the period — SDK version, ledger configuration, HSM configuration. The metadata lets the opposing party re-create the verification environment if they choose to.
- Change-management records showing no integrity-affecting changes during the period. Records covering routine deployments (SDK upgrades that do not break verification, monitoring changes, non-chain configuration) are produced if requested but are not load-bearing for integrity.

**Discovery timeline.** The standard FRCP 34 deadline is 14 days after service of the request. The institution plans for expedited verification during the discovery window — running the verifier against the affected period, generating the discovery package, and confirming the package against the affected customer ID and date range. Tier-1 institutions with frequent litigation may operate a continuous-verification posture that produces the verifier output as a routine artifact; smaller institutions run the verifier on demand when a hold is triggered.

**When to escalate.** If the chain-operations team detects an integrity violation during litigation-hold preparation (verifier reports a failure for the affected period), the IR playbook's Critical-severity scenarios apply — preserve forensic evidence first, then proceed with remediation. The institution's litigation posture under FRCP 37(e) is materially different if the integrity violation is detected during a litigation hold versus during routine operations; the forensic preservation step is required either way, but the spoliation-defense narrative needs the preservation evidence in hand before remediation begins.

**Cross-reference.** `docs/litigation-support.md` (FRCP 26/34 discovery scope, FRCP 37(e) spoliation defense), `docs/incident-response-playbook.md` §"Forensic preservation for Critical-severity scenarios" and §"FRCP 37(e) spoliation defense", spec §10.9 (IKM retention rules), spec §7 (verifier procedure for the affected period).

## Verifier CLI exit codes (per spec §10.12)

Spec §10.12 normates the verifier CLI's exit-code contract. Operations teams scripting batch verification — examiner harnesses, SOC sample-comparison scripts, internal-audit verifier-run automation — branch on these exit codes. The codes are stable across implementations; vendor-specific codes are reserved for codes ≥ 4 and MUST NOT be branched on by examiner harnesses or SOC scripts.

| Exit code | Meaning |
|---|---|
| `0` | PASS or PASS-STRUCTURALLY (witness mode). The verifier completed and the chain (or the structural subset under witness mode) verified |
| `1` | FAIL. The chain failed integrity verification at one of the §7 steps. The reason and step number appear on stdout per the normative output format |
| `2` | Structural / input error. The verifier could not parse the file, required headers were missing, or the file format is unsupported by the verifier version |
| `3` | Configuration error. The verifier was invoked without a required argument (e.g., `--master-key` under `--strict`), the algorithm is unknown, or the posture flag does not match the chain's posture |
| `≥ 4` | Vendor-specific diagnostics. Examiner harnesses and SOC sample-comparison scripts MUST treat exit codes ≥ 4 as opaque diagnostic output and MUST NOT branch on them; the normative reason string on stdout is the load-bearing signal |

The verifier's stdout format is also normative per spec §7. For a failed run the output is three lines: `Status: FAIL`, `Step: N`, `Reason: <text>`. For PASS the output is `Status: PASS` (one line). For witness mode the output is `Status: PASS-STRUCTURALLY, key-bound verification skipped` (one line). Field labels are exact (capitalization, colon, single-space separator); examiner harnesses parse the first three lines and ignore any additional diagnostic lines the implementation appends.

### Worked example — examiner harness scripting batch verification

An examiner running batch verification across an examination period typically scripts a loop over the days in the period, capturing stdout and branching on exit code. A representative shape (POSIX shell):

```bash
#!/usr/bin/env bash
# Examiner harness — batch verification across the period
PERIOD_START="2026-01-01"
PERIOD_END="2026-06-30"
TENANT="tenant_acme_prod_us_east_1"
LEDGER_DIR="/examiner/ledger-snapshot"
MASTER_KEY="/examiner/secure/ikm.bin"

current="$PERIOD_START"
while [[ "$current" <= "$PERIOD_END" ]]; do
    output=$(verifier walk \
        --tenant "$TENANT" \
        --date "$current" \
        --ledger "$LEDGER_DIR" \
        --master-key "$MASTER_KEY" \
        --strict \
        --posture=ffiec 2>&1)
    rc=$?

    case "$rc" in
        0)  echo "[$current] PASS"
            ;;
        1)  # Chain integrity failure. Parse Step and Reason from stdout.
            step=$(echo "$output" | awk -F': ' '/^Step:/ {print $2}')
            reason=$(echo "$output" | awk -F': ' '/^Reason:/ {print $2}')
            echo "[$current] FAIL step=$step reason=$reason"
            # Examiner working-paper entry; route to finding-language.md
            ;;
        2)  echo "[$current] STRUCTURAL ERROR — file unreadable or unsupported"
            # Investigate the ledger snapshot itself
            ;;
        3)  echo "[$current] CONFIG ERROR — verifier invocation problem"
            # Re-check the harness arguments and posture flag
            exit 1
            ;;
        *)  echo "[$current] vendor-specific exit $rc — ignored, stdout: $output"
            # Vendor-specific codes are diagnostic only; harness does not branch
            ;;
    esac

    current=$(date -I -d "$current + 1 day")
done
```

The harness branches on exit codes 0-3 and treats codes ≥ 4 as opaque diagnostic output. The normative reason string on stdout is the load-bearing signal for the FAIL case; the harness routes the step number and reason to the examiner's working-paper entry per `regulator-pack/finding-language.md`.

### Worked example — SOC sample-comparison script

The SOC team's sample-comparison procedure (per `audit-procedures.md` P-13 internal-audit verifier-run cadence) typically scripts verifier runs against the sampled days for a tenant. A representative shape:

```bash
#!/usr/bin/env bash
# SOC sample-comparison — verifier runs against sampled days
SAMPLE_FILE="/soc/working-papers/sample-days.txt"
TENANT="tenant_acme_prod_us_east_1"
RESULTS_DIR="/soc/working-papers/verifier-runs"

while IFS= read -r seal_date; do
    output_file="$RESULTS_DIR/$seal_date.txt"
    verifier walk \
        --tenant "$TENANT" \
        --date "$seal_date" \
        --master-key /soc/secure/ikm.bin \
        --strict > "$output_file" 2>&1
    rc=$?

    status=$(head -n 1 "$output_file" | awk -F': ' '{print $2}')

    if [[ $rc -eq 0 ]]; then
        echo "$seal_date,PASS" >> "$RESULTS_DIR/summary.csv"
    elif [[ $rc -eq 1 ]]; then
        # FAIL — parse step and reason for the working paper
        step=$(grep '^Step:' "$output_file" | awk -F': ' '{print $2}')
        reason=$(grep '^Reason:' "$output_file" | awk -F': ' '{print $2}')
        echo "$seal_date,FAIL,$step,\"$reason\"" >> "$RESULTS_DIR/summary.csv"
        # The SOC team escalates per audit-procedures.md anomaly evidence
        # completeness procedures (P-22 through P-25)
    elif [[ $rc -eq 2 ]] || [[ $rc -eq 3 ]]; then
        echo "$seal_date,ERROR,$rc,\"$status\"" >> "$RESULTS_DIR/summary.csv"
        # Investigate harness configuration before re-running
    else
        # Vendor-specific code — record but do not branch
        echo "$seal_date,VENDOR_DIAG,$rc" >> "$RESULTS_DIR/summary.csv"
    fi
done < "$SAMPLE_FILE"
```

The script produces a CSV summary the SOC team consumes during evidence review. Exit codes 0, 1, 2, and 3 drive the disposition; vendor-specific codes are recorded for transparency but do not change the working-paper outcome.

## Spec version migration

When a new spec version ships:

1. Review the spec's change log; identify any breaking changes
2. Update the implementation to the new spec version
3. The `ffiec.chain.spec` attribute on every event records the spec version of the implementation; the verifier confirms continuity
4. If the spec change is breaking, plan a transition: events under the old spec verify under the old spec rules; events under the new spec verify under the new

## Common operational patterns

### Pattern: Cadence relaxation

The institution wants to relax cadence (daily → weekly). Procedure:

1. Submit examiner-approval request per `regulator-pack/examiner-approval-template.md`
2. Wait for regulator response
3. Update the institution's control description with the approval reference
4. Update the ledger config with the new cadence
5. Restart the seal job (will run weekly going forward)
6. The seal record carries the new cadence; the verifier confirms

### Pattern: Master-key compromise response

Triggered by key-fingerprint reconciliation finding `fingerprint_unmatched_count > 0` on the `master.reconciliation_completed` event, or by external intelligence:

1. Activate IR playbook scenario 4
2. Rotate master immediately
3. Mark the compromise window
4. Notify the regulator within 36 hours

### Pattern: Vendor SDK update

The vendor publishes a new SDK version:

1. Validate the new SDK against the conformance corpus
2. Validate the new SDK against the institution's smoke-test environment
3. Roll out per the institution's change-management procedure
4. Confirm new events appear in the chain with expected `ffiec.chain.spec` value

## SDK local-buffer saturation contract

> Closes G-1. Operationally bounds what the SDK does when its local persistence ring fills, so a regional OTLP outage does not turn into a silent integrity gap or an unbounded customer-call latency spike.

The SDK persists every chain entry to a local ring buffer (SQLite WAL by default per `docs/design/02-chain-construction.md` §4) before returning to the caller. The buffer is bounded. Under sustained OTLP backpressure or a receiver outage, the buffer reaches its high-water mark. What the SDK does at that mark is operator policy, not a default the SDK picks for the institution.

The contract has three components: the policy choice, the operational events emitted, and the partial-write rule.

### Policy choice (configurable, no default)

The SDK exposes `local_buffer.saturation_policy` with two conformant values:

- **`fail_closed`** — when the buffer is full, the SDK drops the new entry, returns control to the caller without raising, and emits `seal.local_buffer_overflow` with the dropped entry's `(run_id, seq)` and the buffer's bytes-used / bytes-capacity at drop time. The caller proceeds; the LLM call completes normally on the customer-facing path. The chain has a gap. The verifier reports the gap at examination time as `chain link broken at seq N` per spec §7. This policy prefers customer-facing availability over chain completeness.
- **`fail_open`** — when the buffer is full, the SDK blocks the calling thread until space becomes available (a successful OTLP export drains the buffer, or the ring overwrites an already-exported entry). The blocking is bounded by `local_buffer.max_block_seconds` (default 30 seconds); if the bound expires, the SDK raises `LocalBufferSaturated` to the caller. This policy prefers chain completeness over customer-facing availability — the LLM call latency spikes, but no gap appears in the chain.

The policy choice is documented in CC8.1. Institutions running customer-facing AI agents typically pick `fail_closed` and accept that an extended OTLP outage produces a documented integrity gap; institutions running back-office workflows where a 30-second pause is acceptable typically pick `fail_open`.

The institution MUST NOT silently overwrite un-exported entries. Overwriting an entry the verifier has not yet seen produces a `chain link broken` finding without an operational-event paper trail, and that is the worst of both options.

### Operational events emitted at the high-water mark

The SDK emits two events bracketing every saturation episode:

| Event | When | Fields |
|---|---|---|
| `audit.buffer.saturation_warned` | Buffer crosses the warning threshold (default 80% of capacity) | `tenant_id`, `bytes_used`, `bytes_capacity`, `oldest_unexported_age_seconds` |
| `audit.buffer.saturation_blocked` (`fail_open`) or `seal.local_buffer_overflow` (`fail_closed`) | Buffer reaches 100% and policy fires | `tenant_id`, `bytes_used`, `bytes_capacity`, `dropped_run_id`, `dropped_seq`, `policy` |

The warning event is the operator's leading indicator. A buffer that crosses 80% during a regional incident gives the on-call ~15 minutes of headroom (typical capture-rate dependent) to shed load, scale the receiver, or accept the upcoming saturation policy outcome. The blocking / overflow event is the actual saturation moment; alerts on the second event are pages, not warnings.

The metric the operator alerts on is `sdk_local_buffer_bytes_used / sdk_local_buffer_bytes_capacity`. The threshold is institution-policy; conservative deployments page at 70%, typical deployments page at 80%, aggressive deployments accept the policy firing without a page.

### Partial-write rule (per spec §4 wire-bound observation)

A half-written record is not observable. Spec §4 anchors this: only on-disk-or-wire bytes are integrity-citable; in-process state is not. The SDK MUST NOT export a record whose persistence layer reports anything other than a complete, fsynced append. If the persistence write fails partway (disk full between `write` and `fsync`, the SDK process dies between the two syscalls, the SQLite WAL checkpoint stalls and aborts), the partially-written record is treated as if it had never been captured. The verifier reads only the last byte-complete `\n`-terminated record and stops the chain at the previous valid `seq`. The SDK on restart re-reads the WAL, finds the last byte-complete entry, and resumes from `seq + 1`.

This rule means a saturation event during a partial-write window cannot produce a malformed record. The buffer is full only of byte-complete records. Partial writes are invisible to the chain — a failure that did not produce wire-observable bytes is, per the spec, not a chain event.

### Backpressure semantics — fail-open vs fail-closed at a glance

| Posture | Buffer-full behavior | Customer-call latency | Chain completeness | Recommended for |
|---|---|---|---|---|
| `fail_closed` | drop new entry, emit `seal.local_buffer_overflow`, return to caller | unaffected | gap appears in chain; verifier reports `chain link broken` | customer-facing AI agents where availability is contractually required |
| `fail_open` | block caller up to `max_block_seconds`, drain via OTLP export, then resume | spikes to `max_block_seconds` worst-case | preserved | back-office, internal-only AI workflows |

Cross-reference: spec §4 (wire-bound observation), §10.3 (append-only enforcement); `docs/design/02-chain-construction.md` §4 (persistence-before-disclosure); `docs/dr-and-resilience.md` "Single ledger instance failure".

## Cardinality budget for OTLP attributes

> Closes G-2 and G-9. Operators wiring chain telemetry into Prometheus, Datadog, or any TSDB-shaped backend need a per-attribute cardinality ceiling so they do not melt the metrics ingest. This section gives the budget.

Chain-of-custody attributes ride on every chain log record per spec §4.4. Naively labelling a Prometheus counter with every chain attribute pushes label cardinality into the millions and tips over a typical Prometheus ingestor at 1-2M active series. Datadog bills per custom-metric-tag, so the same mistake shows up as a budget surprise. The mitigation is a cardinality table the operator consults before adding any chain attribute as a metric label.

### Cardinality table — what to label and what to drop

The table format below matches what Datadog and Prometheus operators recognize: per-attribute cardinality estimate, per-tenant-day distinct values, and operator guidance for each attribute. The estimates are typical-FI numbers; institutions with unusual postures adjust per their tenant count and AI agent surface area.

| Attribute | Typical cardinality | Per tenant-day | Guidance |
|---|---|---|---|
| `ffiec.chain.tenant_id` | ~10^3 to ~10^4 in a tier-1 FI | 1 (constant per tenant) | Safe to label. Bounded by institution tenant count |
| `ffiec.chain.spec` | ~10^0 (one per spec version) | 1 | Safe to label. Bounded by deployed spec versions |
| `ffiec.chain.canonical_encoding` | ~10^0 (one per encoder version) | 1 | Safe to label. Bounded by deployed SDK versions |
| `gen_ai.request.model` | ~10^2 (model + version variants) | ~10^1 typical | Safe to label. Bounded by institution's approved-model list |
| `gen_ai.system` | ~10^1 (vendors per institution) | ~10^1 | Safe to label |
| `audit.routing.destination` | ~10^1 (routing destinations) | ~10^1 | Safe to label |
| `audit.deployment.intent` | ~10^0 (canary, baseline, rollback) | small enum | Safe to label |
| `audit.deployment.canary_traffic_pct` | ~10^2 (0..100 integer) | varies | Aggregate before labelling — bucket into 10% bands |
| `service.name` | ~10^1 to ~10^2 (services per institution) | ~10^1 | Safe to label |
| `service.version` | ~10^2 to ~10^3 (deploy history) | ~10^0 | Drop label after 24h; aggregate by service.name only |
| `ffiec.chain.run_id` | UNBOUNDED — one per AI agent run | thousands to millions | DO NOT label. Use as a trace identifier, not a metric dimension |
| `ffiec.chain.kms_handle_uri` | bounded but rarely useful | 1 | Drop. Per-tenant constant; redundant with tenant_id |
| `ffiec.chain.key_fingerprint` | unbounded over rotation history | 1 within a key version | DO NOT label. Per-rotation distinct value, accumulates over years |
| `audit.routing.circuit_states` (string array form) | 1 attribute, N states inside | bounded array | Safe to label as a single attribute. Do not explode the array into per-state labels |

### The 10^5 rule

Any attribute that exceeds 10^5 distinct values per day per tenant is flagged for redesign. The threshold is the upper bound at which a typical Prometheus ingestor remains performant; above it, query latency degrades and disk usage grows out-of-budget. The operator's monthly cardinality review pulls the top-N attributes by `prometheus_tsdb_head_series` per chain metric and confirms no attribute crosses the threshold.

### Recording-rule pattern — aggregate unbounded dimensions before metric leaves the SDK

The SDK's exported metrics carry the rich attribute set on every record; the operator's recording rule aggregates away the unbounded dimensions before the metric lands in the long-term TSDB. A representative Prometheus recording rule:

```yaml
groups:
  - name: chain_metric_aggregation
    interval: 30s
    rules:
      # Aggregate away run_id and key_fingerprint before storage.
      - record: chain:events_per_tenant_per_minute
        expr: |
          sum by (tenant_id, ffiec_chain_spec, gen_ai_request_model) (
            rate(chain_events_total[1m])
          )
      # Bucket canary percentage into 10% bands for label safety.
      - record: chain:canary_events_per_band
        expr: |
          sum by (tenant_id, canary_band) (
            label_replace(
              rate(chain_events_total{audit_deployment_intent="canary"}[1m]),
              "canary_band", "${1}0pct", "audit_deployment_canary_traffic_pct", "([0-9])[0-9]?"
            )
          )
```

The recording rule keeps the chain metric's high-cardinality fidelity in the raw event stream (where it is needed for forensic replay) while keeping the long-term TSDB at sustainable cardinality (where it powers dashboards and alerts).

### Cardinality budget — audit-event volume noise

A related operational risk is operational-event volume drowning the chain-of-custody signal. Operational events (per `soc-pack/control-evidence-events.md`) carry `chain_kind="operational"` and serve a different purpose from integrity-bearing chain content. An institution's monthly volume budget per event class is the operational discipline:

| Event class | Typical volume budget per tenant per day | Action when exceeded |
|---|---|---|
| Integrity-bearing chain entries (`audit`, `model_call`, `tool_call`, `routing`, `translation`) | per tenant's AI traffic profile | rate-limit the AI agent, not the chain |
| Operational events — control evidence (reconciliation, key rotation, replication completion) | < 100 per tenant-day | investigate emit pattern; one event per cadence per region is the design |
| Operational events — saturation / outage signals | < 10 per tenant-day under steady state | spikes are intentional alerting; do not throttle |
| `seal.*` events (one per cadence per tenant) | 1 per cadence per tenant (24/day for hourly, 1/day for daily, ~0.14/day for weekly) | spikes indicate seal-job retries; investigate HSM availability |

The volume-budget is a guardrail against operational-event storm conditions where a misconfigured SDK emits hundreds of `audit.buffer.saturation_warned` per second during a sustained backpressure event. The SDK MUST rate-limit the warning event to one per minute per tenant per buffer state transition; the saturation-blocked / overflow event is not rate-limited because each one corresponds to a distinct dropped entry.

### Late-binding rate alarm (closes P-3)

`ffiec.chain.late_binding=true` events are normal at low rates and pathological at high rates. The operator alerts on the per-tenant-day ratio:

| Threshold | Alert level | Likely cause |
|---|---|---|
| `late_binding_ratio > 1%` | Warning | replication-mechanism slowdown; investigate cross-region lag |
| `late_binding_ratio > 10%` | Critical | replication is failing; the chain's day-boundary semantics are degraded |

The threshold is institution-default; institutions with cadence-appropriate tolerances declare their own values in CC8.1.

Cross-reference: spec §4.4 attribute table; `docs/design/05-otlp-wire.md` §4.5 attribute schema; `soc-pack/control-evidence-events.md` "Operational event volume budgets".

## Per-event byte budget for capacity planning

> Closes G-6. Operators sizing disk, network, and storage need bytes-per-event at every layer. This section gives the reference numbers and a worked example for both regional and national bank profiles.

Capacity planning starts from one number — bytes per chain event on the wire — and propagates through every layer: SDK SQLite, OTLP wire, ledger WAL, hot store, cold store, replication egress. Without those numbers the operator cannot size disk, network, or backup retention.

### Reference per-event byte budget

The numbers below are typical-FI references. Institutions with verbose attribute payloads (long `gen_ai.completion.text`, large `audit.routing.context_payload`) scale upward; institutions with terse payloads scale downward. The dominant variable is the model-completion text size, which the chain captures as a hash by default but as full text when the institution's CC8.1 declares full-text retention.

| Layer | Bytes per event | Notes |
|---|---|---|
| `gen_ai.completion.completed` log record on the wire | 1.5–3 KB typical | Larger when full-text retention is declared (5–20 KB) |
| Daily seal record on the wire | ~600 bytes plus signature (Ed25519 = 64 bytes) | Constant per tenant per cadence |
| Audit-path entry within the seal | ~32 bytes per Merkle level (~10 levels for 10^4 events) | Tree-depth dependent |
| SDK SQLite WAL (uncompressed, with indices) | ~1.3× the wire bytes | WAL overhead plus index bytes |
| Ledger Postgres WAL (uncompressed) | ~1.5× the wire bytes | WAL plus indexes plus row-versioning |
| Ledger hot store (compressed, indexed) | ~0.7× the wire bytes | After zstd compression |
| Cold store S3 (compressed, parquet/jsonl) | ~0.5× the wire bytes | After columnar compression |
| Cross-region replication egress | 1.0× the wire bytes | Per-region per replicated event |

### Worked example — 10M events/day per tenant

An institution running 10M chain events per day per tenant at 2 KB mean event size on the wire:

```
Daily wire bytes per tenant         = 10M × 2 KB             = 20 GB/day
Daily SDK SQLite WAL                = 20 GB × 1.3            = 26 GB/day
Daily ledger Postgres WAL           = 20 GB × 1.5            = 30 GB/day
Daily hot-store landed (compressed) = 20 GB × 0.7            = 14 GB/day
Daily cold-store landed             = 20 GB × 0.5            = 10 GB/day
Cross-region egress (Pattern A)     = 20 GB × N replicas     = 20 GB × N/day
```

Multi-year retention multiplies cold-store bytes by retention days. A 7-year retention window:

```
7-year cold store per tenant        = 10 GB/day × 365 × 7    = ~25 TB
7-year cold store across 100 tenants                         = ~2.5 PB
```

S3 Standard at typical pricing puts the institution at ~$60K/month for the cold store alone; S3 Glacier Deep Archive drops it to ~$2K/month with the trade-off that recovery for an examination takes hours. Most institutions split: hot tier for the most recent 90 days, warm tier for the next 12 months, deep-archive for the remainder.

### Worked example — regional bank profile

A regional bank running 50K chain events per day per tenant at 2 KB mean event size, 7-year retention, single tenant:

```
Daily wire bytes                    = 50K × 2 KB              = 100 MB/day
Daily SDK SQLite WAL                = 100 MB × 1.3            = 130 MB/day
Daily ledger Postgres WAL           = 100 MB × 1.5            = 150 MB/day
Daily cold-store landed             = 100 MB × 0.5            = 50 MB/day
7-year cold store                   = 50 MB × 365 × 7         = ~125 GB
```

A regional bank's chain footprint is small. A single S3 bucket on Standard tier costs single-digit dollars per month; the cost discussion is dominated by the HSM and the verifier compute, not the storage.

### Worked example — national bank profile

A national bank running 10M chain events per day across 100 tenants at 1.5 KB mean event size, 7-year retention:

```
Per-tenant daily wire bytes         = 10M × 1.5 KB            = 15 GB/day
Total daily wire bytes (100 tenants)                          = 1.5 TB/day
Total daily ledger Postgres WAL                               = 2.25 TB/day
Total daily cold-store landed                                 = 750 GB/day
7-year cold store (compressed)                                = ~1.9 PB
Cross-region egress (Pattern A, 2 replicas)                   = 3.0 TB/day
```

The national-bank profile benefits from columnar compression on the cold store (parquet with zstd typically reaches 0.4× the wire bytes for chain workloads, dropping the 7-year footprint to ~1.5 PB). The cross-region egress dominates the data-platform cost; institutions running Pattern B (per-region tenants) trade verifier-run multiplicity for replication-egress reduction.

### Capacity-planning worksheet

The operator fills in the table below with the institution's actual profile and computes the layer footprints. Reference numbers above provide the multipliers.

| Input | Value | Source |
|---|---|---|
| Events/day per tenant | _____ | AI-agent traffic projection |
| Mean event size (KB) | _____ | sample 10K production events; compute mean |
| Retention window (days) | _____ | regulatory requirement (typical 7 years for audit data) |
| Number of tenants | _____ | institution's tenant inventory |
| Number of replication regions | _____ | Pattern A: 1 seal region + N replicas |

**Computed outputs:**

| Output | Formula |
|---|---|
| Daily wire bytes per tenant | events/day × mean event size |
| Daily wire bytes total | (above) × number of tenants |
| Daily SDK SQLite WAL total | daily wire bytes × 1.3 |
| Daily ledger Postgres WAL | daily wire bytes × 1.5 |
| Daily hot store landed | daily wire bytes × 0.7 |
| Daily cold store landed | daily wire bytes × 0.5 |
| Cold store total (retention) | daily cold store × retention days |
| Cross-region egress per region | daily wire bytes × (N replicas) |
| Verifier compute per audit period | events/day × audit-period days × ~5 µs MAC re-derivation per entry |

The verifier compute estimate is dominated by the per-event MAC re-derivation (HMAC-SHA-256 on 32-byte canonical input). At ~5 µs per entry on commodity hardware, a 7-year audit period for a 10M-events/day tenant runs ~25.5 billion entries × 5 µs = ~35 hours. Examiner harnesses parallelize per day; the wall-clock examination time is hours, not days, when the harness is configured for parallel-by-day.

Cross-reference: spec §4.2 (Merkle tree depth and seal-record size); `docs/at-scale-operations.md` "Hot-path budget"; `docs/cost-model.md`.

## Synthetic-canary control specification

> Closes G-3 and partially G-5. The verifier runs once per audit period; without a continuous-integrity probe, an SDK bug introduced on Friday surfaces Monday morning. The synthetic canary is the load-bearing detection control between deploys and the next verifier run.

The canary is a tagged tenant emitting known-shape events on a fixed cadence; an out-of-band watchdog tails the canary's events and pages on first divergence. The canary is institution-operated, not SDK-vendor-operated — running a vendor canary inside an SDK process tests only the SDK; the institution's canary tests the entire pipeline from SDK through OTLP through the receiver through the seal job through the verifier.

### Canary tenant configuration

The canary tenant is tagged at the resource level so dashboards and alerts can isolate it from production tenants:

```yaml
tenant_id: tenant_canary_us_east_1
tenant_role: canary
seal_cadence: daily
emit_interval_seconds: 60
emit_payload_template: deterministic
```

The canary emits one synthetic chain event per `emit_interval_seconds` (default 60). The emit payload is deterministic per emit-time — `(tenant_id, run_id, seq, captured_at_minute_floor)` deterministically derives the event's `gen_ai.request.prompt`, `gen_ai.completion.text`, and routing payload. Every chain entry's `payload_hash` is therefore pre-computable by the watchdog from the emit-time alone.

### Watchdog process

The watchdog runs out-of-band — distinct host, distinct credentials, distinct deployment — and tails the canary tenant's events from the receiver. For each event the watchdog:

1. Reads `(tenant_id, run_id, seq, captured_at, payload_hash)` from the wire.
2. Re-derives the expected `payload_hash` from the deterministic emit template.
3. Compares the re-derived hash to the on-wire value.

A mismatch is paged within seconds. The watchdog also checks for sequence gaps (`seq` non-monotonic), reordering (`captured_at` regressing), and seal-mismatch at end-of-day (the canary's seal-record's `merkle_root` must match a watchdog-side recomputed root from the day's events).

### Canary invariants

The canary's invariants are the spec's invariants applied to a workload whose ground truth is known:

- **Sequence monotonicity.** `seq` is strictly increasing per `run_id`. Verifier rule per spec §7 step 3.
- **Time monotonicity.** `captured_at` is non-decreasing per `run_id`. Verifier rule per spec §7 step 5.
- **Hash determinism.** `payload_hash` matches the watchdog's pre-computed value byte-for-byte. This is the canary's load-bearing signal — a divergence here means the SDK's canonical-form encoder is broken, the OTLP transport corrupted bytes, or the receiver mutated the record.
- **Seal coverage.** Every emitted event appears in the day's seal record. Watchdog confirms by recomputing the Merkle root.
- **Cross-region replication evidence.** If the canary tenant runs Pattern A, every source region's `master.cross_region_replication_completed` arrives with the expected event count.

### SLO for paging and recovery

The institution's CC8.1 names the canary's SLO. Typical defaults:

| SLO | Target |
|---|---|
| Mean detection time (anomaly to page) | ≤ 5 minutes |
| Mean recovery time (page to root cause identified) | ≤ 15 minutes |
| Watchdog uptime | ≥ 99.9% (the watchdog is the load-bearing control between verifier runs) |
| Canary emit success rate | ≥ 99.9% (failures here are themselves a signal) |

The detection SLO is tight because the canary's purpose is short-circuiting the audit-period detection latency. A weekly verifier run finds an SDK bug introduced on day 1 only at the end of the week; the canary finds it within minutes.

### Operational events emitted by the canary

The canary and watchdog emit institution-defined operational events under the existing schema:

| Event | When |
|---|---|
| `audit.canary.emit_completed` | Every successful canary emit |
| `audit.canary.divergence_detected` | Watchdog finds a mismatch (any invariant) |
| `audit.canary.recovery_completed` | Operator marks the divergence resolved |

The watchdog routes `audit.canary.divergence_detected` to the institution's pager. The on-call's first action is to read the event's `divergence_class` field — `payload_hash_mismatch`, `sequence_gap`, `time_regression`, `seal_mismatch`, `replication_count_mismatch` — and branch into the receiver-outage runbook or the SDK-version-coexistence runbook accordingly.

### Clock-skew detection by the canary (partially closes G-5)

The canary's `captured_at` values come from the canary's host clock; the watchdog's expected `captured_at` comes from the watchdog's clock. A drift between the two surfaces clock skew at canary-emit cadence rather than at verifier-run cadence. The watchdog alerts if the per-emit drift exceeds the SDK's clock-skew threshold (default 30 seconds per the next section); the alert is the same skew surface the SDK's emit refusal would catch on a production tenant, but the canary surfaces it before the production fleet is affected.

Cross-reference: spec §4 (wire-bound observation — the canary tests on-wire bytes, not in-process state); spec §7 (verifier procedure — the canary applies the same checks at higher cadence); `docs/dr-and-resilience.md` "DR exercise procedure".

## Dual-running deploys (rolling SDK upgrades)

> Closes G-7. Standard rolling deploys mix SDK versions for ~30 minutes. The wire-format dispatch token is the contract that lets receivers handle both at once safely; the operational rule is to never roll across a wire-format bump without a drain.

A rolling deploy is the default for any institution at scale. For ~30 minutes during the deploy window, two SDK versions emit traffic to the same receiver and contribute to the same tenant-day's seal. What stays safe and what must be drained is the contract this section names.

### The dispatch token — `sign_payload_version`

Spec §4.3 names `sign_payload_version` as the seal-record field that dispatches verifier reconstruction. The receiver MUST handle whatever `sign_payload_version` a participating SDK emits. The current values in the field's enumeration:

- absent — pre-amendment 6-line `sign_payload` form
- `"v1.0a"` — current 10-line amendment form

A future `"v1.0b"` would be a wire-format bump that the receiver must explicitly support before any SDK in the fleet emits it. The dispatch is byte-level — the verifier reconstructs `sign_payload` from the field as written and verifies against the seal-time signature. A mismatch between the field and the byte form actually used by the signer produces a signature failure at verifier time.

### What's safe in a same-wire-version deploy

If both SDK versions emit `sign_payload_version = "v1.0a"` (no wire-format change between versions), a rolling deploy is safe with these guardrails:

1. **Both versions agree on `canonical_encoding`.** The `canonical_encoding` attribute on every chain entry names the encoder version (per spec §4.4). Both SDK versions in the deploy window MUST emit the same `canonical_encoding` value; if they differ, the verifier reports the cross-vendor reproduction divergence at the next audit-period verifier run. The institution's pre-deploy checklist confirms `canonical_encoding` parity before traffic crosses.
2. **Test-vector parity.** Both versions pass the same conformance test vectors at deploy time. Vector divergence on the same logical input means the SDKs disagree on canonical bytes, which means the seal verifier will pass per-event MAC (each entry signed under its own SDK's canonical form) but cross-vendor reproduction fails. The pre-deploy checklist runs the conformance corpus against both versions.
3. **Routing-event coupling intact.** Both versions emit the same `audit.routing.*` schema (per spec §4.4.1). A coupling-rule divergence (success without paired attempt, or vice versa) shows up in the verifier's P-33 check.

### What's unsafe — wire-format bumps

A wire-format bump means at least one SDK version in the deploy emits a `sign_payload_version` value the other does not understand, or emits attributes under a different schema, or changes the canonical-form encoding in a way that produces different bytes for the same logical input.

The operational rule for wire-format bumps:

1. **Drain the old version before any new-version traffic crosses.** "Drain" means: the old SDK fleet stops accepting new requests, finishes in-flight requests, exports its local SQLite buffer to the receiver, and goes quiet. The receiver waits for the seal job's completion of any tenant-day that includes old-version traffic before accepting new-version traffic.
2. **Never roll forward and back across the bump in the same deploy window.** A canary that promotes v1.5 (new wire format) and then rolls back to v1.4 (old wire format) within the same tenant-day produces a chain interleaved between two `sign_payload_version` values that the verifier cannot stitch into one seal. The rollback path is operationally a new deploy, not a rollback within the same window.
3. **Deploy the receiver-side change first.** The receiver MUST handle the new `sign_payload_version` before any SDK in the fleet emits it. The pre-deploy checklist confirms receiver readiness via the receiver-policy discovery endpoint — the `accepted_sign_payload_versions` list returned by the receiver names every wire form the receiver can process, and the SDK refuses to start under strict mode if its emit version is not in the list.

### Pre-deploy checklist

Operators copy this checklist into the deploy ticket. Each item is a gate; failing any item blocks the deploy.

```
[ ] Both SDK versions present in the deploy emit the same canonical_encoding
[ ] Both SDK versions pass the conformance corpus byte-for-byte at deploy time
[ ] If wire-format bump: receiver fleet handles the new sign_payload_version
[ ] If wire-format bump: drain the old version before new-version traffic
[ ] Receiver-policy discovery returns the accepted_sign_payload_versions list
    that includes the new SDK's emit version
[ ] Synthetic canary (per the canary section) is healthy in the target region
[ ] The institution's change-management record references this checklist
```

### Detection — verifier reports cross-version coexistence

The verifier flags any tenant-day that crossed a `canonical_encoding` boundary or a `sign_payload_version` boundary. The verifier's report names the affected `seq` ranges and the SDK-version transition observed. Operations teams reading the verifier output see the deploy window in the integrity report and can correlate with the change-management record.

### `audit.canary_traffic_pct` source documentation (closes P-5)

For deployments using `audit.deployment.intent="canary"`, the SDK emits `audit.deployment.canary_traffic_pct` per spec §4.4.2. The operator wires the SDK's source for this value:

| Source | When appropriate | Failure mode |
|---|---|---|
| Environment variable, set by orchestrator | static canary deployments | stale if orchestrator does not refresh |
| Feature-flag SDK (LaunchDarkly, etc.) | dynamic canary deployments | flag-service outage produces stale value |
| File watch on a config file the orchestrator updates | high-trust orchestrators | file-watcher latency under fast rollouts |
| HTTP poll against the orchestrator's API | cloud-native deployments | poll-cadence-bounded staleness |

The SDK emits `audit.deployment.canary_source_unavailable` when the source becomes unreadable; the institution's CC8.1 names the source, the staleness tolerance, and the fallback (typically: `audit.deployment.canary_traffic_pct` is omitted from chain entries while the source is unavailable).

Cross-reference: spec §4.3 (sign_payload_version dispatch), §4.4.2 (deployment-intent schema); `docs/design/02-chain-construction.md` §4 (canonical-form encoder); `docs/design/05-otlp-wire.md` §4.6 (receiver-policy `accepted_sign_payload_versions`).

## Receiver-outage runbook

> Closes G-4. The receiver-policy discovery endpoint, OTLP receiver, or seal store goes down. The on-call's runbook is here. The local-buffer saturation contract (above) is the SDK's load-bearing protection during the outage; this section is the operator's procedure.

The receiver-outage runbook covers three failure surfaces that the on-call sees from one paging signal: the receiver-policy discovery endpoint returning 5xx, the OTLP receiver refusing connections, and the seal store accepting writes but failing the daily seal job. The first two are SDK-side observable; the third is ledger-side observable. The runbook below assumes the on-call has been paged; it does not cover the alarm-tuning that produces the page.

### Detection signals

| Signal | What it means | Where it surfaces |
|---|---|---|
| Synthetic canary `audit.canary.divergence_detected` (class `seal_mismatch` or `sequence_gap`) | The receiver is dropping or reordering events | Watchdog page; per the canary section above |
| SDK emits `seal.receiver_unreachable` | The OTLP receiver is unreachable from the SDK | SDK fleet's operational event stream |
| SDK emits `audit.receiver_policy.cache_expired` | The receiver-policy endpoint has been unavailable past the cache TTL | SDK fleet's operational event stream |
| SDK `sdk_local_buffer_bytes_used` rising | Buffers are filling because the receiver is not draining them | Prometheus alert |
| Ledger `ledger_seal_runs_total{status="success"}` not advancing | The daily seal job is not completing | Prometheus alert |
| Verifier reports `chain link broken` for a recent day | Buffers saturated under `fail_closed`; entries dropped | Verifier output |

### Contain — first 15 minutes

1. **Confirm scope.** Is the outage SDK-fleet-wide (every SDK reports `seal.receiver_unreachable`) or partial (one region, one SDK pool)? The synthetic canary's region tag isolates per-region scope.
2. **Verify SDK fallback is firing.** The SDK's `local_buffer.saturation_policy` (per the saturation contract above) is the load-bearing fallback. Confirm via the SDK's emitted operational events that the configured policy is actually firing — `fail_closed` deployments emit `seal.local_buffer_overflow` as buffers fill; `fail_open` deployments emit `audit.buffer.saturation_blocked`.
3. **Acknowledge the alarm trail.** The on-call updates the IR ticket with the alarm trail (canary divergence event, SDK saturation event, ledger seal-job failure). The trail is the timeline opposing counsel reads later if the outage produces a litigation exposure.
4. **Confirm receiver-policy fail-mode behavior.** SDKs in `fail_open` mode emit traffic without a fresh policy fetch; SDKs in `fail_closed` mode refuse to emit and contribute to the buffer saturation. The CC8.1 control description names the institution's choice; the on-call confirms the runtime behavior matches the documented choice.

### Recover — once the receiver is back

The recovery procedure drains buffered events and re-establishes continuity:

1. **Confirm receiver health.** OTLP grpc/http endpoints accept connections, `/readyz` returns 200, the HSM session is established, the seal store accepts writes.
2. **Allow buffered drains.** SDK fleets with non-empty local buffers begin exporting on the next emit attempt. Watch `sdk_local_buffer_bytes_used` decline; the drain rate is bounded by the OTLP receiver's accept rate, not the SDK's emit rate. Monitor for thundering-herd patterns where every SDK in the fleet drains simultaneously and overwhelms the receiver.
3. **Re-verify continuity.** Run the verifier against the affected tenant-days. The verifier's `seq` gap detection flags any dropped entries (chains under `fail_closed` policy will show gaps; chains under `fail_open` policy should be intact). The gaps are documented as integrity-control findings per the IR playbook.
4. **Document the saturation events.** Each `seal.local_buffer_overflow` event captured during the outage is part of the institution's chain-integrity record. The institution's CC8.1 review references the saturation events as documented gaps; the verifier's gap report is reconciled against the saturation events to confirm 1:1 mapping.
5. **Re-fetch receiver policy.** After SDK buffers drain, the SDK's next policy fetch refreshes the cache. Operators wanting an immediate refresh trigger the SDK's `refresh_receiver_policy` admin operation per the SDK's CLI; otherwise the cache refreshes on the next TTL boundary (default 15 minutes).
6. **Re-run the synthetic canary.** Confirm the canary is emitting and the watchdog is verifying. The canary's recovery confirms the end-to-end pipeline rather than just the receiver-side health.

### Receiver-policy cache-expiry escalation (closes G-4 partial)

The receiver-policy fetch fails open by default — the SDK uses the cached policy when fetch fails. If the receiver-policy endpoint is unavailable past the cached-policy max-age (default 4 hours), the SDK MUST escalate to fail-closed regardless of the configured fail-mode. The escalation is the safety net against a long control-plane outage where the SDK runs unbounded on stale policy. The metric `sdk_receiver_policy_age_seconds` exposes the cache age; alerts fire at 1 hour (warning), 3 hours (critical), 4 hours (escalation to fail-closed).

The 4-hour escalation produces `audit.receiver_policy.cache_expired` and the SDK begins refusing to emit traffic. The on-call's runbook step at this point is to either restore the receiver-policy endpoint (returns the SDK to normal fail-open behavior) or to bypass policy enforcement by issuing an institution-signed override (an emergency CC8.1 procedure documented separately; not a routine path).

### HSM unavailability backoff (closes P-1)

A related receiver-side failure is HSM unavailability during the seal job. The seal job retries with exponential backoff per spec §4.3.1; the operational policy:

| Parameter | Value |
|---|---|
| Initial retry delay | 1 second |
| Backoff multiplier | 2 |
| Backoff cap | 5 minutes |
| Jitter | full (random delay between 0 and current backoff) |
| Immediate retry trigger | success-by-other-cluster-member event |
| Regulator notification | 72 hours per spec §4.3.1 |

The metric `hsm_backoff_current_seconds` exposes the current backoff; alerts fire when the backoff caps at 5 minutes (indicating sustained unavailability). The "HSM is back" declaration is automatic — the next successful signing operation resets the backoff to the initial delay.

Cross-reference: spec §4.3.1 (HSM unavailability), §10.3 (append-only enforcement), §10.15 (multi-region replication-completion SLA — closes P-2 with regional emission within 15 minutes of seal-window end and a 30-minute fail-closed reconciliation window per the at-scale operational guidance).

## Disaster recovery — ransomware encrypts the seal store

> Closes G-8. The seal store is the chain's load-bearing artifact. Ransomware encrypting the seal store is a recoverable-but-painful event with explicit restore procedures, documented data-loss bounds, and a chain-of-custody-as-evidence narrative.

The scenario is concrete: ransomware encrypts the institution's Postgres ledger and the most recent S3 backups in the affected region. Three days of seal records are gone. The institution's recovery procedure restores from the last good cross-region replica, verifies continuity, and documents the ransomware event itself as a chain-of-custody event so the chain becomes evidence in the incident-response investigation.

### Recovery steps

1. **Activate the IR playbook.** Ransomware against the seal store is an IR event; the institution activates the IR playbook's Critical-severity scenario for control-system compromise. The IR team's first action is forensic preservation of the encrypted artifacts (the ransomware-encrypted Postgres files, the encrypted S3 objects, any ransom-note artifacts left by the attacker). Preservation is the foundation for the post-incident criminal investigation and any insurance claim.

2. **Restore from the last good cross-region replica.** Pattern A institutions (per spec §10.15) have replication regions that ship events to the seal region before seal-time. The replication regions retain their own copies of the daily seals — the seal-region's seal-record was replicated outward via the institution's standard replication mechanism, so the replication regions hold seal copies up to the moment of replication-cessation. The institution restores from the last replication region whose seal store is intact and whose `master.cross_region_replication_completed` events confirm the day's events were fully replicated before the ransomware event.

3. **Verify continuity via per-day seal records on the surviving replica.** The verifier runs against the restored seal store. The verifier's per-day output names every day that has a seal and every day that does not. The gap between the last-good-seal day and the ransomware-event day is the data-loss window; the institution documents this window in the IR ticket, the regulator notification, and the institution's CC8.1 control failure record.

4. **Document the ransomware event itself as a chain-of-custody event.** The institution captures `master_key.ransomware_compromise_detected` as an institution-defined operational event with the suspected compromise vector, the affected tenant identifiers, the time window of suspected exposure, and the recovery method. This event becomes part of the chain — the chain-of-custody artifact is, at this point, evidence of the institution's response to the ransomware event. The IR investigation uses the chain as evidence rather than as a target.

5. **Re-key the affected tenants.** The institution performs a master-key rotation per the master-key rotation procedure for every tenant whose IKM was potentially exposed. The historical chain remains verifiable under the previous IKM (the cross-region replica preserves it); the new IKM covers events captured after rotation. The institution captures `master_key.rotation_post_ransomware` with the rotation evidence.

6. **Notify the regulator within 36 hours.** The cyber-incident notification rule applies (per the institution's regulator-pack documentation). The notification names the ransomware event as the trigger, the data-loss window, the recovery method, and the post-incident remediation plan. The regulator MAY require additional examination cadence following the event; the institution accommodates per the regulator's findings.

### RPO / RTO under ransomware

The RPO under ransomware is bounded by the cross-region replication lag (typically 15 minutes per the spec §10.15 invariant 5 SLA). The RTO depends on the institution's restore mechanism:

| Restore mechanism | Typical RTO | Trade-off |
|---|---|---|
| Active-active replication restoration | < 1 hour | Highest cost; second region runs continuously |
| Cross-region snapshot restore | 1–4 hours | Snapshot cadence determines RPO |
| S3 Object Lock immutable backup restore | 4–24 hours | Lowest cost; explicit immutable retention |
| Tape or air-gapped offline restore | 24–72 hours | Air-gap protection against networked ransomware |

The institution's CC8.1 names the chosen mechanism, the RTO, and the RPO. Tier-1 institutions typically run active-active replication; smaller institutions choose one of the snapshot-or-immutable-backup options.

### Air-gapped immutable backups

S3 Object Lock with `compliance` retention mode is the recommended pattern for the seal store specifically. The retention mode prevents deletion or modification by any IAM principal — the institution's own root account cannot delete the locked objects until the retention window expires. This is the strongest protection against ransomware that obtains AWS credentials; the attacker can encrypt the live ledger but cannot delete the locked seal records.

The retention window for seal-store immutable backups is institution-policy; tier-1 institutions typically set 7 years (matching the audit-document retention period); regional banks set 1–3 years and accept the cost trade-off for the shorter window.

### Pattern B institutions

Pattern B institutions (per-region tenant_id) recover per region independently. The recovery procedure above applies per regional tenant. The institution's verifier-run plan executes per regional tenant; cross-region correlation is institution-side per the spec §10.15 Pattern B documentation.

### Recovery evidence package

The institution produces a recovery evidence package combining:

- The IR playbook activation timestamp and incident classification
- The forensic preservation evidence of the encrypted artifacts
- The cross-region replica restoration evidence (replication-region seal records, `master.cross_region_replication_completed` events)
- The verifier output covering pre-ransomware and post-recovery periods
- The master-key rotation evidence (`master_key.rotation_post_ransomware` events)
- The regulator notification record
- The data-loss window documentation
- Any insurance-claim documentation referencing the recovery

The recovery evidence package is retained for the institution's standard audit-document retention period; the package is the institution's evidence in any post-incident regulator examination, criminal investigation, or insurance proceeding.

Cross-reference: spec §10.15 (multi-region invariants and `master.cross_region_replication_completed` event); `docs/dr-and-resilience.md` "Backup tampering" and "Multi-region resilience"; `docs/incident-response-playbook.md` Critical-severity scenarios; `docs/legal-disclosure.md` (criminal investigation handoff).

## Clock skew, NTP failure, leap-second handling, cross-region drift

> Closes G-5. The chain's day-boundary semantics depend on accurate time. Clock skew at SDK-emit time produces detectable anomalies at verifier time; cross-region drift produces silent reconciliation failures. This section names the SDK-side and operator-side rules.

The chain's `captured_at` is the SDK's wall-clock reading at emit time; the chain's `received_at` is the receiver's wall-clock reading at ingest. Day-boundary semantics partition by `received_at` per spec §4.2.2. A drift between SDK and receiver clocks therefore produces `late_binding` events — entries arriving at the receiver after their `captured_at` would suggest. The threshold and the operational response are below.

### SDK-side normative behavior on clock skew

The SDK MUST refuse to emit a chain entry whose `captured_at` is earlier than the previous entry's `captured_at` in the same `run_id`. This is the monotonicity rule; the verifier enforces it at audit time per spec §7 step 5, and the SDK enforces it at emit time so the violation never reaches the wire.

Beyond monotonicity, the SDK refuses to emit if its local-clock skew vs NTP exceeds a configurable threshold:

| Configuration | Default | Range |
|---|---|---|
| `clock.skew_threshold_seconds` | 30 | 1–300 |
| `clock.ntp_check_interval_seconds` | 60 | 30–600 |
| `clock.action_on_threshold_exceeded` | `refuse_emit` | `refuse_emit`, `warn_only` |

The default threshold is 30 seconds because spec §4.2.2's day-boundary semantics already absorb up to a 5-minute drift without producing late_binding noise; 30 seconds gives a generous safety margin while keeping NTP-step recovery fast. The institution's CC8.1 names the chosen threshold.

When the threshold is exceeded under `refuse_emit`, the SDK:

1. Emits `audit.clock.skew_exceeded` with the measured skew, the NTP source, and the captured timestamp at detection.
2. Refuses subsequent emit attempts. The local-buffer saturation contract applies — the SDK's `fail_closed` or `fail_open` policy fires once the buffer fills under the refusal.
3. Re-checks NTP at `clock.ntp_check_interval_seconds` cadence. When skew falls below threshold, the SDK emits `audit.clock.skew_recovered` and resumes normal emit.

### NTP-step recovery procedure

When NTP corrects a fleet's clocks (an `ntpd -g` step, a chrony recovery, an AWS time-sync service correction), some hosts may step backwards across a day boundary. The procedure:

1. **Detect.** Synthetic canary alerts on `captured_at` regression (a host's emit at 14:32 followed by the same host's emit at 14:30 after the step). The watchdog's `time_regression` divergence class fires.
2. **Bound the affected window.** Operations team reads the SDK fleet's `audit.clock.skew_exceeded` events to identify which hosts stepped and when. The affected window is the time between the first skew_exceeded event and the last skew_recovered event.
3. **Hold emits during the step.** Hosts that detected the skew already refused to emit (per the SDK's normative behavior above). Hosts that did not detect the skew (skew below threshold) emitted normally; their `captured_at` values are within the day-boundary tolerance and the chain is intact.
4. **Resume.** When the synthetic canary's `time_regression` alarm clears, the operations team confirms the emit refusal lifted across the fleet and watches `late_binding_ratio` for a return to baseline.

### Leap-second handling

Leap seconds are absorbed by NTP. The SDK does not implement leap-second-specific logic; the host's NTP implementation (chrony, ntpd, AWS time-sync) handles leap-second smearing or stepping per the host's configuration. Institutions running smearing-mode NTP (Google's leap-smear, AWS's amzn-time-sync) experience a ~1ppm clock-rate adjustment over 24 hours around the leap-second event; the SDK sees this as a slow drift well below the 30-second threshold. Institutions running stepping-mode NTP experience a 1-second step at the leap-second moment; the SDK sees this as a 1-second skew, well below threshold.

The operator's runbook for leap-second days: confirm the institution's NTP posture is consistent across the SDK fleet (no mix of smearing and stepping hosts within one tenant), confirm the synthetic canary is healthy across the leap-second moment, and document the leap-second event in the change-management log.

### Cross-region drift

Cross-region drift is bounded by the spec §10.15 region-binding invariants. The seal region's `received_at` partitions the seal day; replication regions' `captured_at` events arrive at the seal region with cross-region network latency (typical: < 100ms within a continent, < 300ms intercontinental). The drift between region clocks affects `late_binding` thresholds, not chain integrity — a replication region with a 2-second clock skew vs the seal region produces `captured_at` values 2 seconds in the past or future, which the seal region absorbs into the day-boundary semantics.

The operator's monitoring posture for cross-region drift:

| Metric | Threshold | Action |
|---|---|---|
| Cross-region clock skew (seal region vs replication region, NTP source comparison) | > 5 seconds | Investigate NTP source for affected region |
| `late_binding_ratio` per region | > 1% sustained | Investigate replication mechanism lag |
| `master.cross_region_replication_completed` arrival latency vs seal-time | > 15 minutes | Closes P-2; the regional emission MUST land within 15 minutes; reconciliation fails-closed at 30 minutes with `replication_evidence_missing` |

### Escalation

A clock-skew incident that exceeds 5 minutes in any one host or 1 minute across the fleet is escalated to the network-operations team. The escalation is not a chain-integrity event by default — the SDK's refusal-to-emit behavior protects the chain — but it is a fleet-health event that warrants investigation. The IR playbook does not activate for clock-skew incidents unless the synthetic canary detects a divergence that the SDK's emit refusal failed to catch (which would itself be a separate IR event for the SDK's clock-detection logic).

Cross-reference: spec §4.2.2 (day-boundary semantics), §10.15 (region-binding invariants); `docs/incident-response-playbook.md` "Clock-skew incident response"; `docs/at-scale-operations.md` "Hot-path budget" (for the timing budget that absorbs sub-second skew).

## OpenTelemetry collector compatibility checklist

> Closes the operational consequence of P-6 and provides the deployment-readiness checklist operators run before chain traffic crosses the OTel collector. The chain spec is OTLP-native per spec §4.4.3; this section confirms what the operator validates.

The chain spec is built on standard OTLP. Operators deploying chain-of-custody traffic on top of an existing OTel collector posture validate that the standard collector binary, standard processors, and popular backends handle chain traffic without semantic violations. This is a checklist, not a configuration guide — the configuration guide is the OTel collector configuration section earlier in this document.

### Compatibility checklist

The operator runs each item before declaring the collector deployment ready for chain traffic:

| Item | What to validate | Expected result |
|---|---|---|
| Standard OTLP/HTTP receiver | OTLP/HTTP `:4318` endpoint accepts a sample chain log record | 200 OK; record appears in receiver-side logs |
| Standard OTLP/gRPC receiver | OTLP/gRPC `:4317` endpoint accepts a sample chain log record | gRPC `OK`; record appears in receiver-side logs |
| Batch processor | Standard `batch` processor accumulates chain records and flushes per timeout | Records flushed at timeout; no truncation |
| Attributes processor | Standard `attributes` processor reads / writes / preserves the `ffiec.chain.*` attribute namespace | Attributes preserved byte-for-byte |
| Resource processor | Standard `resource` processor preserves `ffiec.chain.spec`, `service.name`, `service.version` | Resource attributes preserved |
| Routing connector | Routing on `ffiec.chain.spec` Resource attribute splits chain from regular | Chain pipeline isolated; regular pipeline unaffected |
| Severity-filter exemption | Severity filter NOT applied to chain pipeline (per spec §4.4.4) | All chain records pass through, regardless of `SeverityNumber` |
| Honeycomb export | Chain traffic exports to Honeycomb without attribute loss | Chain records visible in Honeycomb with full attribute set |
| Lightstep export | Chain traffic exports to Lightstep without attribute loss | Chain records visible in Lightstep with full attribute set |
| Datadog export | Chain traffic exports to Datadog without attribute loss | Chain records visible in Datadog with full attribute set |
| Prometheus export | Chain metrics export to Prometheus with cardinality budget intact | No high-cardinality label explosions |
| Loki export | Chain log records export to Loki with `SeverityText="OTLP"` preserved | Log records queryable by `SeverityText` |

The exporter validations are commodity: the chain spec is OTLP-native, so any OTLP-compliant backend handles chain traffic. Failures in this row indicate either the backend has a non-OTLP-compliant attribute renaming (rare in modern backends), or the operator's pipeline introduces a custom processor that mutates chain attributes. The fix is the pipeline, not the backend.

### Vendor-collector compatibility

Operators running a vendor-customized collector (Datadog Agent, Splunk OpenTelemetry Collector, Honeycomb Refinery) validate that the vendor's customization preserves the chain pipeline's processing posture. Specifically:

- The vendor's collector MUST NOT apply severity filters to chain traffic. Per spec §4.4.4 the chain pipeline is exempt; a vendor-defaults filter that drops `< INFO` is a control-completeness gap. The institution's CC8.1 procedure samples for this via P-38.
- The vendor's collector MUST preserve the `ffiec.chain.spec` Resource attribute. A collector that strips or renames Resource attributes makes the routing connector unable to identify chain traffic, which collapses the pipeline's chain-vs-regular split.
- The vendor's collector MUST preserve `SeverityText="OTLP"`. Vendors that rewrite `SeverityText` to a normalized value (e.g., always `INFO`) break the SIEM-side filter pattern documented in the SIEM exemption section below.

### SIEM alert exemption pattern (closes P-6)

Most SIEMs (Splunk, Datadog, ELK) ship with default alerting rules that fire on `severity == ERROR` or `severity == any non-INFO value`. Chain records carrying `SeverityNumber` in the `9..20` range will trigger these alerts unless the operator carves an exemption. The exemption is per-SIEM; concrete examples:

**Splunk SPL.**

```
index=otlp_chain
| where SeverityText != "OTLP"
| <existing alert rule>
```

The SPL filter excludes chain records before the alert rule evaluates. Operators add the filter to every existing severity-based alert that ingests from the chain receiver's index.

**Datadog log monitor.**

```
status:error -severity_text:OTLP
```

The Datadog query syntax excludes chain records via the negative filter. Operators update existing log monitors to include `-severity_text:OTLP`.

**ELK / Kibana KQL.**

```
status: "error" AND NOT severity_text: "OTLP"
```

The KQL filter excludes chain records from the dashboard panels and alert rules. Operators update saved searches and visualizations.

The `SeverityText="OTLP"` value is unique enough across standard OTel severity texts (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`) that the filter does not collide with non-chain traffic. Institutions running their own receiver with a different severity text apply the same pattern with their chosen text.

### Severity treatment in operations

Spec §4.4.4 normates the severity range and the `SeverityText` value; the operational consequence is that chain traffic cannot be filtered out by routine severity-based filters. Operators adapt their dashboards, alert rules, and aggregation queries to either include chain traffic intentionally (for chain-of-custody dashboards) or exclude it explicitly (for the SIEM exemption pattern above). The default — applying a severity filter that drops chain records — is non-conformant and produces a control-completeness gap.

The operator's pre-deploy verification: the institution's existing SIEM dashboards and alert rules are reviewed for severity-based filters; each filter is updated to include the chain exemption per the SIEM-specific pattern above; the change-management record references this document.

### Ed25519 strict canonicalization audit (closes P-4)

A related operational discipline is verifier-side audit of Ed25519 strict canonicalization. Spec §4.3 requires the verifier to reject non-canonical signatures. The audit-procedure to confirm production verifier behavior:

1. Sample one signature per tenant per quarter from production seal records.
2. Re-encode the signature into canonical form (using a reference encoder distinct from the production verifier's encoder).
3. Compare the re-encoded canonical bytes to the persisted bytes byte-for-byte.
4. If they differ, the production verifier accepted a non-canonical signature; this is a control failure routed through the institution's CC8.1 procedure.

The audit-procedure is run quarterly per the institution's SOC cadence. The reference encoder is institution-chosen — typical choices are the OpenSSL Ed25519 implementation, the Go `crypto/ed25519` library, or the BoringSSL implementation. Cross-vendor reference encoder usage prevents a single-vendor library defect from masking the issue.

Cross-reference: spec §4.4.3 (OTLP transport identification), §4.4.4 (severity treatment), §4.3 (Ed25519 strict canonicalization); `docs/design/05-otlp-wire.md` §4.5 (severity treatment) and §4.6 (receiver-policy informative); `docs/audit-procedures.md` P-38 (severity exemption sampling).

## Capacity-planning worksheet

> Closes G-6 and G-10. The operator fills in the worksheet with institution-specific numbers; the tables compute the layer footprints. Reference numbers and worked examples appear in the per-event byte-budget section above.

This worksheet is the consolidated planning artifact. Operators copy the tables into the institution's infrastructure-planning document, fill in the inputs, and compute the outputs. The worksheet supports two profiles: a regional bank profile and a national bank profile. Institutions outside these ranges adapt by interpolating from the worked examples.

### Inputs

| Input | Regional bank example | National bank example | Your value |
|---|---|---|---|
| Events/day per tenant | 50K | 10M | _____ |
| Mean event size on the wire (KB) | 2.0 | 1.5 | _____ |
| Retention window (days) | 2,555 (7 years) | 2,555 (7 years) | _____ |
| Number of tenants | 1 | 100 | _____ |
| Number of replication regions (Pattern A) | 1 | 2 | _____ |
| Daily seal-job duration target (minutes) | 5 | 30 | _____ |
| Verifier audit-period duration (months) | 12 | 12 | _____ |

### Computed outputs — capacity

| Output | Formula | Regional example | National example |
|---|---|---|---|
| Daily wire bytes per tenant | events/day × mean event size | 100 MB | 15 GB |
| Daily wire bytes total | (above) × number of tenants | 100 MB | 1.5 TB |
| Daily SDK SQLite WAL total | daily wire bytes × 1.3 | 130 MB | 1.95 TB |
| Daily ledger Postgres WAL | daily wire bytes × 1.5 | 150 MB | 2.25 TB |
| Daily hot store landed | daily wire bytes × 0.7 | 70 MB | 1.05 TB |
| Daily cold store landed | daily wire bytes × 0.5 | 50 MB | 750 GB |
| 7-year cold store total | daily cold store × retention days | ~125 GB | ~1.9 PB |
| Cross-region egress per region per day | daily wire bytes × replication regions | 0 | 3.0 TB |
| Verifier compute per audit period | events/day × 365 × ~5 µs | ~25 minutes | ~3.5 hours |

The verifier compute estimate parallelizes per-day; the wall-clock time scales down by the number of parallel-by-day workers the harness uses. A tier-1 institution running 10 parallel workers completes a 12-month verifier audit in ~21 minutes wall-clock.

### Computed outputs — network

| Output | Formula | Regional example | National example |
|---|---|---|---|
| Peak wire throughput (events/sec) | events/day / (24 × 3600) × peak-to-mean ratio (3) | ~1.7 events/sec | ~350 events/sec |
| Peak wire bandwidth (Mbps) | peak throughput × mean event size × 8 / 1024 | ~0.03 Mbps | ~4 Mbps |
| OTLP gRPC connections from SDK fleet | SDK pool size × 1 per pool | 1–5 | 100–1000 |
| Cross-region replication bandwidth per region (Mbps) | daily egress × 8 / 86400 / 1024 | 0 | ~280 Mbps |

The peak-to-mean ratio of 3 is typical for AI-agent workloads; institutions with bursty traffic profiles (regulatory filings, batch inference jobs) see ratios up to 10. Operators sizing network capacity use the institution's measured peak ratio rather than the default.

### Computed outputs — tail-latency planning (closes G-10)

The SDK's per-event capture latency has a tail. Operators planning for the tail size SDK host capacity to absorb periodic stalls without customer-facing impact.

| Percentile | Typical latency | Likely cause of stalls at this percentile |
|---|---|---|
| p50 | ~200 µs | normal capture path |
| p99 | ~1.8 ms | normal capture path with brief lock contention |
| p99.9 | ~10 ms | SQLite WAL checkpoint pause |
| p99.99 | ~850 ms | SQLite WAL checkpoint stall on large WAL; OTLP gRPC backpressure |
| max | ~3 s | HSM round-trip outliers in Model B HMAC-via-HSM dispatch; OTLP receiver brief unavailability |

For a 10K-events/sec workload, the p99.99 percentile produces ~1 stall per second. AI agents calling the SDK from request-handling code see periodic latency spikes. Operators mitigate with:

- **Async emit.** The SDK's `emit_async` path returns to the caller before the WAL append completes; the WAL append happens on a background thread. Trades the tail-latency-on-the-caller-path for a slight increase in the customer-facing happy-path latency (the `emit_async` queue insert is ~5 µs).
- **WAL checkpoint tuning.** Increase the SQLite `wal_autocheckpoint` to spread checkpoints over more events; reduces stall frequency at the cost of slower hot-store recovery on SDK restart.
- **OTLP gRPC connection pooling.** Multiple gRPC connections from the SDK to the receiver smooth backpressure; the SDK's `otlp.grpc.pool_size` defaults to 4 per pool.
- **Model A vs Model B trade-off.** Model B (HMAC-via-HSM) introduces an HSM round-trip per event, which dominates the tail. Institutions running Model B for risk-posture reasons accept the tail; institutions running Model A do not see the HSM-round-trip tail at all.

### SDK metric for tail-latency monitoring

`sdk_capture_duration_seconds_bucket` exposes the per-event capture latency as a Prometheus histogram. The operator's recommended bucket boundaries:

```
1ms, 2ms, 5ms, 10ms, 50ms, 100ms, 500ms, 1s, 5s, +Inf
```

Buckets cover p50 through max. Operators alert on sustained p99.99 above 1 second (indicating the tail has moved into pathological territory) and on max above 5 seconds (indicating an HSM or OTLP receiver hard fault).

### Per-cadence seal-job sizing

The seal job's compute budget per cadence:

| Cadence | Events per seal | Merkle-tree depth | HSM signing operations | Wall-clock budget |
|---|---|---|---|---|
| Hourly | events/24 | log2(events/24) | 1 | ≤ 60 min per spec §4.3 |
| Daily | events/day | log2(events/day) | 1 | ≤ 60 min per spec §4.3 |
| Weekly | events × 7 | log2(events × 7) | 1 | ≤ 25h per spec §4.3 (weekly cadence's relaxed window) |

Hourly cadence on a high-throughput tenant (1B events/hour) needs careful budget planning — the 60-minute window covers reading 1B events from the WAL, building a 30-level Merkle tree, dispatching one HSM signature, and appending the seal record. Reference per-phase allotment for the 1B-events/hour case:

| Phase | Allotment |
|---|---|
| Read events from WAL | 25 minutes |
| Build Merkle tree | 25 minutes |
| HSM signing round-trip | 1 minute |
| Append seal record | 1 minute |
| Slack for retries / backpressure | 8 minutes |

Tier-1 institutions running hourly cadence at this scale typically distribute the Merkle build across multiple ledger workers to keep the wall-clock under budget; the institution's CC8.1 names the seal-job's compute topology.

### Worksheet completion checklist

After filling in the worksheet, the operator confirms:

```
[ ] Daily wire bytes total fits the planned network capacity at peak
[ ] 7-year cold store fits the planned long-term storage budget
[ ] Cross-region egress fits the planned inter-region network budget
[ ] Verifier compute per audit period fits the planned audit-window budget
[ ] Per-cadence seal-job sizing fits the spec §4.3 publish window
[ ] Tail-latency percentiles fit the AI agent's customer-facing latency budget
[ ] Cardinality budget per attribute fits the metrics ingest's cardinality ceiling
```

The checklist is the institution's pre-deployment artifact. Failure to complete any item is a deployment blocker; the institution either resizes the budget or revises the input assumptions.

Cross-reference: spec §4.2 (Merkle tree depth and seal-record size), §4.3 (publish window per cadence); `docs/at-scale-operations.md` "Hot-path budget"; `docs/cost-model.md`; `docs/design/02-chain-construction.md` §4 (SDK capture path).

## Where to go for more

- **Design rationale.** `docs/design/`
- **Threat model.** `docs/design/09-threat-model.md`
- **Cost picture.** `docs/cost-model.md`
- **Glossary.** `docs/design/10-glossary.md`
- **Incident response.** `docs/incident-response-playbook.md`
- **BYOC specifics.** `docs/byoc-deployment.md`
- **Cloud HSM specifics.** `docs/cloud-hsm-guide.md`
- **DR specifics.** `docs/dr-and-resilience.md`
