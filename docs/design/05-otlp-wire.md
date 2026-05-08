# 05 — OTLP wire format

> **What this doc is.** How chain-of-custody fields ride on the OpenTelemetry Protocol. The wire is OTLP; the chain is an attribute extension; the canonical-hash form is JCS (separate from the wire). This doc clarifies the layering.

## 1. Why OTLP

- **Industry standard.** OTel is CNCF-graduated, multi-language, vendor-neutral. The Splunk OpenTelemetry Distribution, the Datadog Agent, the AWS Distro for OpenTelemetry, and Honeycomb all consume OTLP.
- **Already in the bank&rsquo;s stack.** Most banks have deployed the OTel Collector or its vendor-distributed equivalents. Adding chain attributes does not require a new ingestion path.
- **Bidirectional compatibility.** A bank can run a SIEM that consumes OTLP plus a ledger server that consumes the same OTLP stream &mdash; tee, not fan-in.
- **GenAI semantic conventions.** OTel published `gen_ai.*` semconv in 2024. We carry chain fields in `ffiec.chain.*` alongside.

## 2. The wire layering

```
┌─────────────────────────────────────────────────────────┐
│  OTLP envelope (protobuf, gRPC or HTTP)                 │
│  ┌───────────────────────────────────────────────────┐  │
│  │  ResourceSpans                                    │  │
│  │  ┌─────────────────────────────────────────────┐  │  │
│  │  │  ScopeSpans                                 │  │  │
│  │  │  ┌───────────────────────────────────────┐  │  │  │
│  │  │  │  Span                                 │  │  │  │
│  │  │  │  - name, kind, start/end, status      │  │  │  │
│  │  │  │  - attributes (key/value pairs):      │  │  │  │
│  │  │  │      gen_ai.system: "anthropic"       │  │  │  │
│  │  │  │      gen_ai.request.model: "..."      │  │  │  │
│  │  │  │      ffiec.chain.spec: "v1.0"         │  │  │  │
│  │  │  │      ffiec.chain.run_id: "..."        │  │  │  │
│  │  │  │      ffiec.chain.seq: 42              │  │  │  │
│  │  │  │      ffiec.chain.prev_hash: <bytes>   │  │  │  │
│  │  │  │      ffiec.chain.payload_hash: <bytes>│  │  │  │
│  │  │  │      ffiec.chain.key_version: 3      │  │  │  │
│  │  │  │      ffiec.chain.key_fingerprint: ".."│  │  │  │
│  │  │  │      ffiec.chain.tenant_id: "..."     │  │  │  │
│  │  │  │      ffiec.chain.captured_at: ts      │  │  │  │
│  │  │  └───────────────────────────────────────┘  │  │  │
│  │  └─────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

The chain extension fields are **regular OTel span attributes**. Any OTel-aware backend treats them as ordinary key/value pairs. The chain logic (verification, sealing) reads them as structured data; non-chain backends pass them through unchanged.

## 3. JCS canonical form vs OTLP wire form

The two forms are independent and serve different purposes:

| Form | Purpose | When it&rsquo;s used |
|---|---|---|
| **JCS canonical JSON** | Deterministic payload-hash input | At capture, in the SDK, before the chain HMAC is computed |
| **OTLP protobuf** | Wire transport | After the chain HMAC is computed, when the event is exported |

The chain HMAC is computed over the JCS form. The wire transports the OTLP protobuf form *which carries the already-computed HMAC value*. Two implementations with different OTLP encoders produce identical chain HMACs because they compute over identical JCS canonical bytes &mdash; the protobuf encoding is irrelevant to the hash.

This separation is **load-bearing**: it lets the wire format evolve (OTLP v2, OTLP v3) without breaking chain agreement across implementations.

## 4. Attribute namespace

### 4.1 The `ffiec.chain.*` namespace

| Attribute | Type (OTel) | Required | Description |
|---|---|---|---|
| `ffiec.chain.spec` | string | yes | `"v1.0"` &mdash; the spec version this event conforms to |
| `ffiec.chain.run_id` | string | yes | Bounded-length opaque identifier for the run |
| `ffiec.chain.seq` | int64 | yes | Monotonically increasing within a run, starting at 1 |
| `ffiec.chain.prev_hash` | bytes | yes | 32 bytes; zeros for `seq==1` |
| `ffiec.chain.payload_hash` | bytes | yes | 32 bytes; HMAC-SHA-256 over JCS payload |
| `ffiec.chain.key_version` | int64 | yes | Integer (≥ 1) identifying the IKM generation (post-rework v1.0; replaces `session_key_id`) |
| `ffiec.chain.key_fingerprint` | bytes | yes | 16 raw bytes; `SHA-256(utf8(tenant_id) \|\| ikm)[:16]`; verifier asserts looked-up IKM produces this fingerprint BEFORE computing any MAC (post-rework v1.0) |
| `ffiec.chain.format_version` | string | yes | `"v1"` for this spec; verifier refuses unrecognized values most-specific-first (post-rework v1.0) |
| `ffiec.chain.chain_kind` | string | yes | Closed enumeration per spec §3: one of `"audit"` (default — application audit event) \| `"model_call"` (LLM invocation) \| `"tool_call"` (tool invocation) \| `"routing"` (routing-decision entry per spec §4.4.1) \| `"translation"` (ECOA translation entry per spec §10.11) \| `"operational"` (control-evidence operational events). The verifier MUST reject any value not in the enumerated set with `chain_kind out of v1 enumeration at seq N`. The value is integrity-bound under the OTel envelope per spec §5 (the `chain_kind` field appears in the §5 inclusion list alongside `trace_id`, `span_id`, `name`, `kind`, etc.) |
| `ffiec.chain.tenant_id` | string | yes | Tenant identifier; matches the HSM&rsquo;s key label |
| `ffiec.chain.captured_at` | int64 | yes | Unix nanoseconds, UTC |
| `ffiec.chain.canonical_encoding` | string | yes | Canonical-encoding identifier the SDK used to produce `payload_hash`. Default `"rfc8785-jcs"` at `format_version = "v1"`. The verifier asserts the encoding is one it knows how to reproduce; an unrecognized value is rejected the same way an unrecognized `format_version` is rejected. The field is forward-compatible: a future v2 or v3 spec MAY introduce a different canonical encoder (different deterministic JSON form, CBOR, or another stable byte representation) and stamp a different identifier here without retroactively breaking v1 chains. The verifier branches on this value before computing `payload_hash`; v1 chains continue to verify under the v1 encoder regardless of what later spec versions adopt |

### 4.1.1 No-key-material guarantee

The chain attributes are integrity-bearing identifiers (`payload_hash`, `prev_hash`, `key_fingerprint`, `key_version`, `format_version`) and routing metadata (`tenant_id`, `run_id`, `seq`, `captured_at`). **None of them are key material.** The `key_fingerprint` is a public per-tenant identity binding (`SHA-256(utf8(tenant_id) || ikm)[:16]`); the per-entry stamp is safe by design (the IKM-length minimum per spec §10.6 closes offline-grinding). SDKs MUST NOT route session keys, IKM bytes, or HSM tokens to any OTLP-bound stream.

### 4.1.2 Optional v1.0-final attributes (locked)

Working-group decision (2026-05-15): four optional attributes are added to v1.0-final. Implementations MAY emit them; verifiers MUST handle them when present and ignore them when absent. Implementations that emit `master_version` or `algorithm` SHOULD emit them on every event for consistency.

| Attribute | Type | Required | Purpose |
|---|---|---|---|
| `ffiec.chain.algorithm` | string | optional | The HMAC algorithm used (default `"HMAC-SHA-256"`). Avoids a spec-version bump if HMAC-SHA-256 weakens. |
| `ffiec.chain.master_version` | string | optional (deprecated under post-rework v1.0; use `ffiec.chain.key_version` integer field instead) | Legacy field name for the IKM generation. Post-rework v1.0 uses `ffiec.chain.key_version` (int64) per spec §4.4; implementations migrating from earlier drafts SHOULD emit `key_version` going forward. The seal record carries `key_versions` (list of integers) for the day's chain entries. |
| `ffiec.chain.parent_run_id` | string | optional | Links a child run to its parent in multi-process agent flows. |
| `ffiec.chain.parent_seq` | int64 | optional | The parent's seq at handoff. Required when `parent_run_id` is present. |
| `ffiec.chain.dag_parents` | string | optional | Comma-separated list of `(run_id, seq)` pairs for DAG-shaped multi-process flows. Mutually exclusive with `parent_run_id`. |
| `ffiec.chain.gen_ai_parameters` | string | optional | JCS-canonical JSON of model sampling parameters (temperature, top_p, seed, etc.) for SR 11-7 reproducibility. |
| `gen_ai.provider_attestation` | bytes or string | optional | Provider-issued cryptographic attestation that the named model produced the response (shipped by Anthropic, OpenAI, and Google for some endpoints). Captured byte-for-byte under the OTel envelope. Validated institution-side, not by the verifier. See §4.1.3 for the validation procedure shape. |

The latter two were added based on round-3 review: `dag_parents` for DAG-shaped multi-process runs, `gen_ai_parameters` for model-state reproducibility under SR 11-7 effective challenge. The `gen_ai.provider_attestation` attribute was added based on round-13 review (Maya Patel, AI platform reliability engineer, Q-6) for institutions whose evidence claim benefits from provider-side attestation alongside chain-side recording integrity.

### 4.1.3 Provider-attestation validation procedure (institution-side)

Some LLM providers ship a cryptographic attestation alongside their responses — Anthropic, OpenAI, and Google publish attestation public keys and produce per-response signatures binding "this model produced this content" to the provider's identity. An institution that captures the attestation into the chain gains a stronger, two-part evidence claim:

- **Chain-side claim** (verifier-validated): the chain proves the institution recorded what the provider sent. This is the `payload_hash` HMAC over the canonical bytes of the response, sealed under the daily Merkle root and HSM-signed.
- **Provider-side claim** (institution-validated): the attestation proves the provider sent it. This is the provider's signature over the response content, validated against the provider's published attestation public key.

The two claims compose without overlap. Neither replaces the other. The chain by itself proves recording integrity; the attestation by itself proves provider integrity; together they close both halves of "did this AI output really come from this provider, and did the institution really record exactly what the provider sent."

The validation is **institution-side responsibility, not verifier responsibility**. The spec keeps the verifier focused on chain-integrity decisions; the attestation is opaque to the chain by design. The institution operates the validation as a separate workflow.

**Attestation-validation worker (recommended shape).** The institution operates a worker that reads chain entries containing the `gen_ai.provider_attestation` attribute, validates each attestation against the relevant provider's published key, and emits validation-result chain entries back into the chain so the SOC team and the examiner can sample them. The worker shape:

```
1. Read chain entries with gen_ai.provider_attestation present.
2. Parse the attestation envelope (provider-specific format; the worker
   implements one parser per supported provider).
3. Look up the provider's published attestation public key by key
   identifier carried in the envelope (institution-side key registry,
   refreshed on the provider's published rotation cadence).
4. Validate the signature over the response content the chain recorded.
5. Emit an audit.attestation_verification.* chain entry recording the
   validation result, the provider, the key identifier validated against,
   and the validation timestamp. The entry's parent linkage points back
   to the original (run_id, seq) of the chain entry whose attestation
   was validated.
```

**`audit.attestation_verification.*` chain-entry attributes (recommended schema).** The validation-result entry carries enough context for the SOC team to reproduce the validation independently:

| Attribute | Required | Notes |
|---|---|---|
| `audit.attestation_verification.target_run_id` | yes | The `run_id` of the chain entry whose attestation was validated |
| `audit.attestation_verification.target_seq` | yes | The `seq` of that entry |
| `audit.attestation_verification.provider` | yes | The provider name (`"anthropic"`, `"openai"`, `"google"`, etc.) |
| `audit.attestation_verification.key_identifier` | yes | The provider-key identifier from the attestation envelope |
| `audit.attestation_verification.result` | yes | `"valid"`, `"invalid_signature"`, `"expired_key"`, `"unknown_key"`, `"malformed_envelope"`, etc. |
| `audit.attestation_verification.validated_at_utc` | yes | RFC 3339 UTC timestamp |
| `audit.attestation_verification.validator_version` | optional | Forensic; the worker's software version |

**Failure handling — PASS-WITH-ANOMALY, not chain failure.** A failed attestation validation does NOT break chain integrity. The chain proves the institution recorded what the provider sent; the chain is correct regardless of whether the provider's attestation later validates. Validation failures are PASS-WITH-ANOMALY findings — the chain entry is sealed, the validation entry is sealed, and the institution investigates the anomaly through the appropriate path:

- `invalid_signature` — the recorded response and the recorded attestation disagree. Likely causes: a transport-layer mutation between provider and SDK that the SDK did not detect (an OTel-collector processor that touched `gen_ai.*` content despite the pass-through rule of spec §4.4), an attestation-format change the worker has not learned, or, in the threat-model worst case, a man-in-the-middle who altered the response without altering the attestation. The institution treats this as an investigation trigger.
- `expired_key` — the attestation key was retired before validation ran. Likely causes: validation backlog longer than the provider's key rotation window, institution's key-registry refresh stale. The institution adjusts the worker's cadence or refreshes the registry.
- `unknown_key` — the attestation envelope names a key the institution does not have. Likely causes: provider published a new key the institution has not picked up, or, less likely, an attestation purportedly from a provider the institution has not configured. The institution refreshes the registry.
- `malformed_envelope` — the envelope did not parse. Likely causes: provider format change the worker has not learned, or upstream encoding damage. The institution updates the worker.

None of these break the chain. The verifier's report on a chain containing attestation entries reads as PASS for the chain itself; the institution's attestation-validation report reads alongside it as a separate evidence stream the examiner can sample.

**Why opaque to the chain.** The verifier ignores `gen_ai.provider_attestation` because including it in chain integrity would entangle two different trust roots: the institution's HSM (which signs the chain) and the provider's attestation key (which signs the response). Keeping them separate keeps each evidence claim independently meaningful. The chain stays valid even if a provider's attestation key is later disclosed or rotated; the attestation stays valid evidence even if the chain's HSM is rotated. The two compose without one's compromise undermining the other's claim.

### 4.2 The `gen_ai.*` namespace (referenced from OTel)

We use OTel&rsquo;s GenAI semconv unchanged. The chain extension complements rather than overlaps. Examples we expect on AI events:

- `gen_ai.system` &mdash; the LLM provider
- `gen_ai.request.model` &mdash; the model identifier
- `gen_ai.usage.prompt_tokens`, `gen_ai.usage.completion_tokens`
- `gen_ai.response.finish_reason`

These are not part of the chain hash unless the SDK explicitly includes them in the canonical payload. Implementations document which fields are included.

### 4.3 The `audit.routing.*` namespace (institution-emitted, spec §4.4.1)

Spec §4.4.1 defines a normative attribute schema for AI-routing-decision chain entries: `audit.routing.event_type`, `audit.routing.providers_attempted`, `audit.routing.provider_chosen`, `audit.routing.failover_reason`, `audit.routing.refusal_reason`, `audit.routing.circuit_state.<provider>`, `audit.routing.decision_at`, `audit.routing.policy_version`, `audit.routing.cost_factor`, and `audit.routing.bypass_reason`. The namespace is registered here as part of the chain's documented OTel attribute namespace alongside `ffiec.chain.*`, `gen_ai.*`, and the broader `audit.*` namespace the institution operates.

The namespace is **institution-emitted**, NOT provided by an SDK semconv module. The institution's chain decorator subscribes to its existing router's routing-event hook (per `02-chain-construction.md` §11) and produces the routing chain entry with these attributes populated from the router's hook payload. The OTel project does not (as of v1.0-final) publish a `gen_ai.routing.*` semconv module — the chain's `audit.routing.*` namespace fills the gap until OTel formalises one, at which point the institution and the spec working group evaluate alignment.

The wire encoding follows the same `audit.*` rules as other audit attributes: each attribute appears as a regular OTel span attribute in the OTLP envelope, the chain HMAC is computed over the JCS canonical form of the attribute payload, and the institution's standard OTLP transport carries the bytes to the ledger. The OTel collector's pass-through rule (spec §4.4) applies — collector processors MUST NOT modify `audit.routing.*` attributes in transit; modifying them invalidates the chain MAC at the verifier.

Type conventions follow the spec §4.4.1 table: `event_type` and string-valued attributes are OTel `string`; `providers_attempted` is OTel `string[]`; `cost_factor` is OTel `double`; `decision_at` is the institution's standard timestamp encoding (RFC 3339 string per spec §4.4.1, OR Unix nanoseconds `int64` if the institution's broader audit-timestamp convention uses nanoseconds — institutions document the encoding choice in their CC8.1 control description so verifier-side reconstruction is unambiguous). Per-provider circuit-state attributes use the sub-key form `audit.routing.circuit_state.<provider>` where `<provider>` is the institution's provider identifier from the routing pool (e.g. `audit.routing.circuit_state.openai-gpt-4o`).

#### 4.3.1 Event types — five values closed at v1.0

Spec §4.4.1 names five routing event types, carried as the span `name` OR as the `audit.routing.event_type` attribute:

| Event type | Emitted when |
|---|---|
| `audit.routing.attempt` | The router selects a provider and is about to invoke |
| `audit.routing.success` | The selected provider returned successfully |
| `audit.routing.failover` | The previously-attempted provider failed; the router moved to the next provider in the policy |
| `audit.routing.circuit_state_change` | The router's circuit-breaker state for some provider transitioned (e.g., closed → open after a threshold) |
| `audit.routing.refused` | The router evaluates policy and determines no call can be made (all circuits open, no provider in policy, institution quota exhausted, cost threshold at capacity, or policy override) |

The `refused` event type closes the gap where the router evaluates policy and concludes no call can be made — the `circuit_state_change` and `failover` event types do not fit the no-call case mechanically (no state transitioned at the moment of decision; nothing was being failed over from). The chained `refused` entry IS the routing-decision evidence; the absence of a child LLM-call entry confirms no call followed.

#### 4.3.2 Refusal-reason attribute — closed enumeration

Spec §4.4.1 names `audit.routing.refusal_reason` as REQUIRED on `refused` events. The enumeration is closed at v1: `all_circuits_open` | `no_provider_in_policy` | `quota_exhausted` | `cost_threshold_at_capacity` | `policy_override`. The institution's chain decorator emits the value verbatim from the router's hook payload; SOC sample-comparison procedures (per `audit-procedures.md` P-33) confirm every `refused` entry's `refusal_reason` is one of the five enumerated values. The `audit.routing.provider_chosen` and `audit.routing.failover_reason` attributes are NOT present on `refused` entries (no provider was selected; no prior provider was being failed over from).

#### 4.3.3 Failover-reason discriminator semantics

Spec §4.4.1 names `audit.routing.failover_reason` with the closed enumeration `timeout` | `transport_error` | `provider_error` | `circuit_open` | `rate_limit` | `quota_exhausted` | `cost_threshold` | `policy_override` | `manual_override`. Three discriminators are distinct because their MRM committee dispositions, customer-dispute responses, and SOC sample-comparison procedures branch on the distinction:

- `rate_limit` is provider-side (the provider's API returned 429 or equivalent).
- `quota_exhausted` is institution-side (the institution's per-tenant or per-customer-tier quota was exhausted before the call left the institution's perimeter).
- `cost_threshold` is per-call cost guard (not aggregate quota).

The wire form preserves the value verbatim; consumers branch on the literal string. The `audit-procedures.md` P-33 procedure tests the discriminator's correctness against the institution's router behavior.

### 4.4 Transport identification (per spec §4.4.3)

Spec §4.4.3 makes a small set of OTLP Resource attributes normative on every chain-of-custody export, plus a recommended-but-not-required set of transport-layer headers (HTTP) and metadata (gRPC) so collectors can identify chain traffic at the transport layer without inspecting span content. The Resource attributes are integrity-bound under the OTel envelope; the headers and metadata are routing hints — collectors that drop or rewrite them MUST NOT alter the underlying Resource attributes.

#### 4.4.1 Required Resource attributes

Every OTLP `ResourceSpans` carrying chain entries MUST include the following Resource attributes:

| Attribute | Type | Notes |
|---|---|---|
| `ffiec.chain.spec` | string | The spec version, identical to the per-span attribute (e.g. `"v1.0"`). Stamped at the Resource level so receivers can identify chain traffic without descending into ScopeSpans |
| `service.name` | string | Standard OTel service-name attribute. The institution's SDK-host service identifier |
| `service.version` | string | Standard OTel service-version attribute. The SDK build version |
| `ffiec.chain.posture` | string | One of `"ffiec"`, `"hipaa"`, `"pci"`, `"soc2"`, or another posture identifier the institution operates under. Lets receivers route per-posture |
| `ffiec.chain.format_version` | string | The chain format version (e.g. `"v1"`). Identical to the per-span attribute; the Resource-level stamp lets receivers branch on format without descending into spans |

The Resource-level stamps are integrity-bound by spec §5 alongside the per-span chain attributes. A collector or receiver that rewrites Resource attributes at the OTLP boundary breaks the verifier's per-span MAC because the canonical payload includes the Resource-level identification on the spec-defined inclusion list.

#### 4.4.2 Recommended HTTP headers

OTLP/HTTP exporters SHOULD include the following request headers when shipping chain traffic. Headers are routing hints — receivers that filter, route, or rate-limit chain traffic at the HTTP-proxy layer use them; receivers that read the OTLP body anyway can ignore them. A collector that strips these headers MUST NOT also alter Resource attributes (the chain integrity-binds the Resource attributes; the headers are convenience).

| Header | Value | Notes |
|---|---|---|
| `X-FFIEC-Chain-Spec` | `"v1.0"` | The spec version, mirroring the Resource attribute. Lets HTTP-layer routers identify chain traffic without parsing protobuf |
| `X-FFIEC-Chain-Posture` | e.g. `"ffiec"` | The posture, mirroring the Resource attribute |

#### 4.4.3 Recommended gRPC metadata

OTLP/gRPC exporters SHOULD include the following metadata entries. gRPC metadata uses lowercase keys per the gRPC convention.

| Key | Value | Notes |
|---|---|---|
| `ffiec-chain-spec` | `"v1.0"` | The spec version |
| `ffiec-chain-posture` | e.g. `"ffiec"` | The posture |

#### 4.4.4 Worked example — OTLP/HTTP request shape

A representative OTLP/HTTP request shipping one chain entry. The headers and Resource section together let collectors and receivers identify the traffic before reading per-span content:

```http
POST /v1/traces HTTP/1.1
Host: ledger.example-bank.internal
Content-Type: application/x-protobuf
X-FFIEC-Chain-Spec: v1.0
X-FFIEC-Chain-Posture: ffiec
Content-Length: 1834

<protobuf body, decoded for clarity>

ResourceSpans {
  Resource {
    attributes {
      service.name: "loan-advisor-agent"
      service.version: "3.1.0"
      ffiec.chain.spec: "v1.0"
      ffiec.chain.format_version: "v1"
      ffiec.chain.posture: "ffiec"
    }
  }
  ScopeSpans {
    Scope { name: "ffiec.chain", version: "1.0" }
    Span {
      trace_id: ...
      span_id: ...
      name: "audit.advisory.decision"
      attributes {
        gen_ai.system: "anthropic"
        gen_ai.request.model: "claude-opus-4-7"
        ffiec.chain.spec: "v1.0"
        ffiec.chain.run_id: "r_a3f29b71c"
        ffiec.chain.seq: 42
        ffiec.chain.prev_hash: <32 bytes>
        ffiec.chain.payload_hash: <32 bytes>
        ffiec.chain.canonical_encoding: "rfc8785-jcs"
        ffiec.chain.format_version: "v1"
        ffiec.chain.key_version: 3
        ffiec.chain.key_fingerprint: <16 bytes>
        ffiec.chain.tenant_id: "tenant_acme_prod"
        ffiec.chain.captured_at: 1717564800000000000
        ffiec.chain.chain_kind: "audit"
      }
    }
  }
}
```

The Resource attributes appear once per `ResourceSpans` envelope; per-span attributes carry the per-event chain fields. A receiver that wants to cheaply count chain traffic per service version groups by `(service.name, service.version, ffiec.chain.posture)` at the Resource level without descending into per-span content. A receiver that needs to verify integrity descends into per-span content and runs the verifier per spec §7.

### 4.5 Severity treatment for chain-of-custody traffic (per spec §4.4.4)

Spec §4.4.4 makes severity-handling for chain traffic normative because operators have repeatedly tripped over the OTel default. By default, OTel SDKs emit log records at `SeverityNumber = 9` (INFO) and the standard collector configurations apply severity filters that drop everything below INFO or below WARN depending on environment. A chain-of-custody record dropped by a severity filter is a silent integrity gap — the chain verifies what arrived at the ledger, not what the SDK tried to send.

The spec closes this with two normative rules:

#### 4.5.1 Collector pass-through rule (normative)

Collectors MUST NOT filter chain-of-custody traffic by severity. Collector configurations that apply severity-based filters (the OTel `filter` processor with severity predicates, the `tail_sampling` processor with severity-based policies, vendor-specific severity gates) MUST exempt records carrying the `ffiec.chain.spec` Resource attribute. Identification at the collector layer is by Resource attribute — the same attribute the receiver later integrity-binds — so a collector that misidentifies chain traffic produces a non-conformant pipeline that the SOC team's P-38 procedure will surface.

#### 4.5.2 Receiver-stamping convention (normative)

Receivers MUST stamp chain-of-custody records with a non-default `SeverityNumber` so downstream filters that have not been updated for the pass-through rule still admit chain traffic by default. The reference implementation (Herald.Compliance receiver) stamps:

| Field | Value | Rationale |
|---|---|---|
| `SeverityNumber` | dynamic in range `9..20` (positioned by Herald's `QuickLogBuilder`) | Spec §4.4.4 mandates the receiver position the level via a resolver in the range `INFO ≤ N ≤ FATAL − 1` (i.e., `9..20`). Higher-range positioning resists routine `< INFO`, `< WARN`, and `< ERROR` filters more aggressively; staying below `FATAL = 21` keeps the record from being mis-treated as a system-fatal alert by routine alerting infrastructure. The Herald reference receiver uses its `QuickLogBuilder` component to pick the value dynamically per ingest based on institution policy. The spec deliberately does not pin a single numeric value because the right level is institution-tuned per the institution's resilience program and routine alerting baseline |
| `SeverityText` | `"OTLP"` | Unique enough to grep for in operator logs, dashboards, and SIEM searches. Operators looking for chain traffic in their log analysis tools filter on `SeverityText == "OTLP"` to pull every chain record without parsing OTel attribute payloads. Does not collide with the standard OTel severity texts (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`) |

The stamping is institution-customizable when documented in CC8.1 — an institution operating its own receiver MAY use a different resolver and a different `SeverityText` provided the produced `SeverityNumber` lands in the `9..20` range mandated by spec §4.4.4. The institution's CC8.1 control description names the resolver mechanism, the produced range under its policy, and the `SeverityText`. The SOC team's P-38 procedure samples the receiver configuration against the institution's documented stamping.

#### 4.5.3 Worked example — OpenTelemetry Collector configuration

A representative collector configuration that exempts chain traffic from severity filters by routing it through a separate pipeline branch identified by the `ffiec.chain.spec` Resource attribute. The pattern is: receive once, route by Resource attribute, apply standard processing (filtering, sampling, batching) to the regular telemetry pipeline only, send chain traffic through a no-filter, no-sampling pipeline to the chain-of-custody exporter.

```yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  # Routing connector splits traffic by Resource attribute. Records with
  # ffiec.chain.spec present are routed to the chain pipeline; everything
  # else goes through the regular telemetry pipeline.
  routing:
    default_pipelines: [traces/regular]
    table:
      - statement: 'route() where resource.attributes["ffiec.chain.spec"] != nil'
        pipelines: [traces/chain]

  # Severity filter for the regular telemetry pipeline only.
  # MUST NOT appear in the chain pipeline per spec §4.4.4.
  filter/severity:
    error_mode: ignore
    logs:
      log_record:
        - 'severity_number < SEVERITY_NUMBER_INFO'

  # Tail-sampling for the regular telemetry pipeline only.
  # MUST NOT appear in the chain pipeline.
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
    # Chain-of-custody pipeline. NO severity filter. NO sampling. Batching
    # only. Resource attribute ffiec.chain.spec identifies this branch.
    traces/chain:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp/chain]

    # Regular telemetry pipeline. Standard severity filter + tail sampling.
    traces/regular:
      receivers: [otlp]
      processors: [filter/severity, tail_sampling, batch]
      exporters: [otlphttp/regular]
```

The pattern is symmetric for OTLP `logs` and `metrics` pipelines — replace `traces/*` with `logs/*` or `metrics/*` and use the analogous routing connector and processors. Institutions running multiple postures (FFIEC plus HIPAA, for example) extend the routing table with additional `ffiec.chain.posture` predicates and add per-posture pipelines.

The collector config above is the load-bearing artifact the SOC team's P-38 procedure samples. Operators document the config in source control alongside the institution's other infrastructure-as-code artifacts; changes to the chain pipeline branch route through the institution's standard change-management procedure.

### 4.6 Receiver-policy discovery (informative)

The chain-of-custody-v1 spec does NOT mandate a receiver-policy discovery endpoint. The endpoint described here is implementation-specific to the Herald reference topology, where the SDK and the chain receiver run as distinct processes (the TesseraSeal / Herald.Compliance split). Other implementations may run the receiver as a library inside the SDK process and skip the network hop entirely; the new spec §4 implementation-topology framing makes this distinction explicit.

This section documents the Herald reference shape for implementers building a similar topology and for auditors who want to understand what the SDK fetches from the receiver before traffic starts flowing.

**Purpose.** The SDK queries the receiver to obtain per-tenant data the receiver will apply on ingest:

- The actual `SeverityNumber` the receiver-side resolver (Herald reference: `QuickLogBuilder`) will position for this tenant. Sits within the spec §4.4.4 range `[9, 20]`. The SDK uses this to surface the actual position alongside the spec range in operator views.
- The set of receiver-side filter exemptions configured in the receiver's internal pipeline.
- The accepted `sign_payload_version` forms — the SDK validates its emitted version matches at least one accepted form before emitting traffic.
- The accepted `format_version` and posture values.

**Endpoint shape (Herald reference).**

```
GET https://<receiver-base-url>/api/v1/tenants/<tenant_id>/receiver-policy

Authorization: Bearer <token>          (or mTLS client cert)
Accept: application/json
If-None-Match: "<etag>"                (refresh path only)
```

Other implementations are free to use a different URL path or a different discovery mechanism (configuration file, control-plane API, environment variable). The spec does not normate the URL; it normates only the transport-security floor and the data the SDK needs.

**Authentication.** HTTPS-only; the discovery endpoint inherits the OTLP transport-security floor mandated in §5.1 — TLS 1.3 minimum, server-authenticated TLS, Bearer token or mTLS client auth. Implementations SHOULD refuse plaintext URLs at SDK construction time so an operator misconfiguration cannot leak a Bearer token over plain HTTP.

**Response shape (example).** The body is a JSON object with the fields the SDK consumes. Implementations are free to add additional fields; SDKs ignore fields they don't recognise.

```json
{
  "policy_version": "v1",
  "receiver_stamp": {
    "severity_number_range": [9, 20],
    "severity_number_resolver": "herald.compliance.QuickLogBuilder",
    "severity_text": "OTLP",
    "actual_value": 17
  },
  "filter_exemptions": ["filter/severity", "tail_sampling"],
  "sign_payload": {
    "accepted_forms": ["v1"],
    "rejected_forms": []
  },
  "format_version": "v1",
  "posture_accepted": ["ffiec"]
}
```

The `actual_value` field is the value Herald.Compliance reports for this tenant on this fetch — the value `QuickLogBuilder` will position when the next chain entry arrives. It sits inside the `severity_number_range` window. The SDK surfaces both the range and the actual value in operator-facing views so operators see the spec window plus the receiver's actual position.

**Caching.** Short TTL plus ETag-based revalidation. The Herald reference SDK uses a 15-minute default TTL; refresh calls send `If-None-Match: <etag>`. The receiver returns 304 when the ETag matches and 200 with a fresh body when the policy has changed. The SDK refreshes opportunistically — operators don't see the receiver-policy fetch on the chain emit path; the cache amortises it.

**Failure modes.** Implementations choose between fail-open and fail-closed when the discovery endpoint is unreachable:

- **Fail-open.** Use the cached policy when one exists; emit traffic without a fresh policy check when the cache is empty. Suitable for institutions where availability of the chain emit path matters more than strict policy alignment.
- **Fail-closed.** Refuse to emit traffic to the affected tenant until the policy fetch succeeds. Suitable for high-assurance deployments where a misconfigured receiver shouldn't see traffic.

The Herald reference SDK exposes both modes via a `strict` flag on its `ReceiverPolicyClient`.

**Out-of-spec scope.** The endpoint is implementation-specific. Spec §5.1 mandates only the transport-security floor that any such endpoint inherits. Implementations that skip the network hop entirely (the receiver-as-library topology) satisfy the spec without implementing this endpoint. Implementations that run a receiver in a different process MAY implement a receiver-policy endpoint; if they do, the transport-security floor in §5.1 applies.

Cross-references: spec §4 (implementation topology framing); spec §4.4.4 (severity range); spec §5.1 (transport security); `docs/operator-guide.md` §"Receiver-policy discovery" for operator-side configuration.

## 5. The canonical payload

The JCS-hashed payload is the SDK&rsquo;s explicit canonical representation, *not* the OTLP wire bytes. The included fields are:

```json
{
  "captured_at_ns": 1717564800000000000,
  "kind": "audit",
  "run_id": "r_a3f29b71c",
  "seq": 42,
  "tenant_id": "tenant_acme_prod",
  "event": {
    "decision": "approve",
    "amount": 1240.18,
    "rationale": "...",
    "agent_version": "3.1.0"
  },
  "gen_ai": {
    "system": "anthropic",
    "request": {"model": "claude-opus-4-7"},
    "response": {"finish_reason": "end_turn"}
  }
}
```

Implementations document the *complete* schema of the canonical payload. Adding fields silently breaks chain agreement; the conformance corpus enforces this.

## 5.1 Transport encryption (normative for v1.0-final)

The chain's HMAC catches injection and tampering on the wire; it does not provide confidentiality. AI agent payloads frequently contain customer-sensitive content (PII, financial data, private business context). Implementations MUST encrypt OTLP transport between the SDK and the ledger:

- **TLS 1.3 minimum.** TLS 1.2 is acceptable only for legacy environments where the institution documents the exception in its control description. **The TLS 1.2 exception sunsets on 2028-01-01.** After that date, all conforming deployments MUST use TLS 1.3 or higher; institutions with legacy hardware that cannot upgrade by then must deploy a TLS-terminating proxy that upgrades to 1.3 toward the ledger and retire the legacy hardware on its normal refresh cycle.
- **Server-authenticated TLS** is the floor; **mutual TLS** is required where the institution's posture mandates it (typical for production banking deployments).
- **Cipher suites** follow the institution's standard policy; the spec does not enumerate.

For indirect transports (Kafka, Kinesis), the wire encryption is the broker's encryption-at-rest plus the producer-to-broker and broker-to-consumer encryption-in-transit settings; the institution documents the path in its control description.

## 6. Transport options

### 6.1 OTLP/gRPC (recommended)

- Default OTLP transport
- Streaming, multiplexed, binary
- Supported by every OTel SDK and every modern backend
- The reference implementation supports both gRPC and HTTP receivers

### 6.2 OTLP/HTTP (also supported)

- Acceptable when gRPC is blocked by network policy
- Same protobuf encoding, sent as `application/x-protobuf` over HTTPS
- Minor performance penalty, no functional difference

### 6.3 Indirect transport (Kafka, Kinesis)

- Acceptable for delivery
- The wire encoding is OTLP protobuf inside the message body
- The chain extension fields ride along unchanged
- Conforming implementations document the indirect transport configuration

The verifier does not care how the events arrived. It reads the ledger; the ledger&rsquo;s internal representation is implementation-flexible as long as the chain extension fields are recoverable.

## 7. Backend interop

### 7.1 SIEM ingestion

Splunk, Sentinel, QRadar, CrowdStrike Falcon all ingest OTLP. The chain extension fields appear as named attribute key/value pairs in the SIEM&rsquo;s event format. Detection rules can correlate on `ffiec.chain.tenant_id`, `ffiec.chain.run_id`, `ffiec.chain.session_key_id`.

### 7.2 FinOps ingestion

CloudZero, Vantage, and similar tools ingest OTel traces with `gen_ai.usage.*` attributes for cost roll-up. The chain extension attributes are passed through; FinOps tools generally ignore them but can join on `ffiec.chain.tenant_id` and `ffiec.chain.run_id` to attribute cost per agent run.

### 7.3 Existing observability backends

Honeycomb, Datadog, Lightstep, Grafana Tempo all ingest OTLP. The chain extension attributes appear unchanged. The institution&rsquo;s existing observability stream remains useful for engineering debugging without a separate path.

The substrate framing is the architecture: one captured stream serves multiple downstream consumers without translation.

## 8. The OTel Semantic Conventions registration

The `ffiec.chain.*` and `audit.*` namespaces are currently informal. Registration with OpenTelemetry semantic conventions under FFIEC stewardship is a v1.0-final-amendment post-submission milestone — the spec is locked at v1.0-final-amendment, the conformance corpus is built, and the registration request goes to the OTel semconv community once an FFIEC-side sponsor is identified. The chain attributes are stable in v1.0; registration is paperwork, not a wire-format change.

Open questions for the registration:

- Should the namespace be `ffiec.chain` or a more general `audit_chain.*`? The latter generalizes beyond banking.
- Should `payload_hash` and `prev_hash` be string-encoded (hex) or raw bytes? OTel semconv prefers strings; v1.0 uses bytes for byte-form determinism (the canonical bytes go straight into the per-event MAC and the daily Merkle leaf without an intermediate hex round-trip). The OTel semconv tradeoff was settled in favor of bytes; the registration submission documents the rationale.
- Should the spec version be a separate attribute or encoded in a tracestate header?

These are tracked in the spec change log.

## 9. Auditor's-lens review

| Question | Answer |
|---|---|
| Why depend on OTel? | OTel is the lingua franca. Banks already use it. Vendors compete on quality of OTel implementation. The spec rides on the standard; the standard is not ours to reinvent. |
| What if OTel changes the wire format? | The chain extension fields are attribute key/value pairs. Wire format changes (protobuf v2 to v3) preserve the attribute model. The spec stays compatible; the reference implementation tracks OTel releases. |
| Is the wire format part of what the auditor verifies? | No. The auditor verifies the chain (HMAC, Merkle, HSM signature). The wire is how the bytes arrived; verification does not depend on it. The verifier could in principle accept a CSV dump of the chain and verify it identically. |
| What if a vendor uses a custom wire format? | Non-conforming. The spec is OTLP. Vendors can transport OTLP over any pipe (gRPC, HTTP, Kafka), but the format on that pipe is OTLP protobuf with the chain attributes. |
| OTel semconv registration | The `ffiec.chain.*` and `audit.*` namespaces are stable in v1.0; registration with the OTel semconv registry under FFIEC stewardship is a v1.0-final-amendment post-submission milestone, gated only on identifying an FFIEC-side sponsor in the OTel community. The spec is locked at v1.0-final-amendment and the conformance corpus is built; registration is paperwork, not a wire-format change. |
