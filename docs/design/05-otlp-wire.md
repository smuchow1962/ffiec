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
│  │  │  │      ffiec.chain.session_key_id: ".." │  │  │  │
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
| `ffiec.chain.session_key_id` | string | yes | Opaque ID for the session key (not the key itself) |
| `ffiec.chain.tenant_id` | string | yes | Tenant identifier; matches the HSM&rsquo;s key label |
| `ffiec.chain.captured_at` | int64 | yes | Unix nanoseconds, UTC |

### 4.2 The `gen_ai.*` namespace (referenced from OTel)

We use OTel&rsquo;s GenAI semconv unchanged. The chain extension complements rather than overlaps. Examples we expect on AI events:

- `gen_ai.system` &mdash; the LLM provider
- `gen_ai.request.model` &mdash; the model identifier
- `gen_ai.usage.prompt_tokens`, `gen_ai.usage.completion_tokens`
- `gen_ai.response.finish_reason`

These are not part of the chain hash unless the SDK explicitly includes them in the canonical payload. Implementations document which fields are included.

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

The `ffiec.chain.*` namespace is currently informal. The intent is to file an OTel semconv proposal for inclusion in the registry once the spec is at v1.0-final.

Open questions for the registration:

- Should the namespace be `ffiec.chain` or a more general `audit_chain.*`? The latter generalizes beyond banking.
- Should `payload_hash` and `prev_hash` be string-encoded (hex) or raw bytes? OTel semconv prefers strings; we have used bytes in v1.0-draft. Trade-off is wire size vs convention conformance.
- Should the spec version be a separate attribute or encoded in a tracestate header?

These are tracked in the spec change log.

## 9. Auditor's-lens review

| Question | Answer |
|---|---|
| Why depend on OTel? | OTel is the lingua franca. Banks already use it. Vendors compete on quality of OTel implementation. The spec rides on the standard; the standard is not ours to reinvent. |
| What if OTel changes the wire format? | The chain extension fields are attribute key/value pairs. Wire format changes (protobuf v2 to v3) preserve the attribute model. The spec stays compatible; the reference implementation tracks OTel releases. |
| Is the wire format part of what the auditor verifies? | No. The auditor verifies the chain (HMAC, Merkle, HSM signature). The wire is how the bytes arrived; verification does not depend on it. The verifier could in principle accept a CSV dump of the chain and verify it identically. |
| What if a vendor uses a custom wire format? | Non-conforming. The spec is OTLP. Vendors can transport OTLP over any pipe (gRPC, HTTP, Kafka), but the format on that pipe is OTLP protobuf with the chain attributes. |
| Open issue | Registration of `ffiec.chain.*` in the OTel semconv registry. Tracked for v1.0-final; sponsor in the OTel community to be identified. |
