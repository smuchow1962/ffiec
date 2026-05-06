# Chain-of-Custody Specification, Version 1.0 (draft)

> **Status:** v1.0-draft. Design phase. Subject to change before v1.0-final per the process in [GOVERNANCE.md](../GOVERNANCE.md).
> **Audience:** implementers of conforming SDKs, ledger servers, and verifiers; auditors and examiners reviewing the standard; regulators evaluating whether to adopt it.
> **Conformance keywords.** The keywords MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are to be interpreted as described in [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) and [RFC 8174](https://www.rfc-editor.org/rfc/rfc8174).

---

## 1. Scope

This specification defines a chain-of-custody primitive set for capturing AI-driven decisions in regulated systems. The primitives are designed to satisfy the integrity-of-logging requirements in the FFIEC IT Examination Handbook (Information Security and AIO booklets) and related model-risk-management guidance (SR 11-7, OCC Bulletin 2011-12).

The primitives are language-neutral, transport-agnostic in their core construction (though a wire-format binding to OTLP is normative), and produce output that can be independently verified by an examiner with no access to the institution's vendor or operations team.

## 2. Out of scope

This specification does not define:

- The runtime environment in which AI agents execute
- The semantic meaning of captured events (those follow the OpenTelemetry GenAI semantic conventions, referenced normatively below)
- The retention duration of ledger artifacts (regulatory frameworks set this; implementations support configurable retention)
- The specific HSM model or vendor used for root signing (any FIPS 140-2 Level 3 or higher device is conformant)

## 3. Definitions

| Term | Definition |
|---|---|
| **Event** | A discrete record of one observable AI activity (a model call, a tool call, a decision, a retry). |
| **Run** | A bounded sequence of events sharing one `run_id`, representing one logical agent invocation. |
| **Tenant** | A regulated institution subscribed to a chain-of-custody implementation. Each tenant has independent keys and a separate ledger. |
| **Tenant-day** | The set of events for one tenant within one UTC calendar day, used as the Merkle-seal aggregation unit. |
| **Master key** | A long-lived HMAC key bound to a tenant. Held in escrow at tenant-controlled storage; never on application hosts. |
| **Session key** | A short-lived HMAC key derived from the master via HKDF, bound to a single application process. |
| **HSM** | A Hardware Security Module conforming to FIPS 140-2 Level 3 or higher (or Common Criteria EAL4+). Used for daily root signing. |

## 4. The four primitives

### 4.1 Primitive 1 — HMAC chain at capture (normative)

Every captured event MUST carry the SHA-256 of the previous event in the same `run_id`, hashed under HMAC-SHA-256 with the per-process session key. The sequence forms a hash chain per run.

For the first event of a run, `prev_hash` MUST be 32 zero bytes.

```
session_key  = HKDF-SHA256(master_key, salt = process_uuid || tenant_id, info = "ffiec-ai-chain-v1", L = 32)
prev_hash    = 0x00...00 (32 bytes) for seq=1
             = payload_hash of prior event for seq>1
canonical    = canonical_json(event_payload)   // RFC 8785 JSON Canonicalization Scheme
payload_hash = HMAC-SHA-256(session_key, prev_hash || canonical)
```

The chain construction MUST occur in the application process at the moment the event is created, before the event leaves the host.

### 4.2 Primitive 2 — Daily Merkle seal (normative)

For each tenant-day, the ledger server MUST construct a binary Merkle tree over the `payload_hash` of every event captured during that UTC day, ordered by `(run_id, seq)` ascending.

The Merkle tree MUST use SHA-256 with the leaf-prefix and node-prefix scheme of [RFC 6962](https://www.rfc-editor.org/rfc/rfc6962):

```
leaf_hash     = SHA-256(0x00 || payload_hash)
internal_hash = SHA-256(0x01 || left || right)
```

Odd-leaf trees MUST be balanced by promoting the unpaired leaf to the next level (the same scheme used by RFC 6962). The Merkle root is the apex hash.

### 4.3 Primitive 3 — HSM-rooted root signature (normative)

The daily Merkle root MUST be signed using Ed25519 in HSM custody. The signature MUST cover:

```
sign_payload = "ffiec-ai-chain-v1\n" ||
               tenant_id || "\n" ||
               iso8601_date(tenant_day) || "\n" ||
               hex(merkle_root)
```

The HSM MUST hold the Ed25519 private key under FIPS 140-2 Level 3 or higher protection. The corresponding public key is published via the institution's tenant key registry (out of scope of this specification, but implementations MUST provide a documented retrieval path).

The signed root MUST be appended to the ledger within 60 minutes of UTC midnight on the day after the tenant-day it covers. Implementations MAY publish provisional roots earlier; the final root is the one the examiner verifies.

### 4.4 Primitive 4 — OpenTelemetry-native wire (normative)

Events MUST ship over OpenTelemetry Protocol (OTLP) using the OTel GenAI Semantic Conventions for AI-specific attributes. Chain-of-custody fields MUST be encoded as OTLP span attributes under the `ffiec.chain.*` namespace:

| Attribute | Type | Required | Description |
|---|---|---|---|
| `ffiec.chain.spec` | string | yes | Spec version, e.g. `"v1.0"` |
| `ffiec.chain.run_id` | string | yes | The run identifier |
| `ffiec.chain.seq` | int64 | yes | Sequence number within the run, starting at 1 |
| `ffiec.chain.prev_hash` | bytes | yes | 32-byte SHA-256 of previous event |
| `ffiec.chain.payload_hash` | bytes | yes | 32-byte HMAC-SHA-256 of canonical payload |
| `ffiec.chain.session_key_id` | string | yes | Identifier for the per-process session key (not the key itself) |
| `ffiec.chain.tenant_id` | string | yes | The tenant identifier |
| `ffiec.chain.captured_at` | timestamp | yes | UTC timestamp of capture, with nanosecond precision |

OTLP gRPC and HTTP transports are both conformant. Other transports (Kafka, Kinesis, etc.) MAY be used for delivery but MUST encode the OTLP envelope within their payload.

## 5. Wire format

The on-the-wire encoding for chain extension fields MUST follow OTLP protobuf encoding. The reference implementation in [`core/otlp/`](../core/otlp/) provides encode and decode functions.

The canonical-JSON form used inside `payload_hash` MUST follow [RFC 8785 (JSON Canonicalization Scheme, JCS)](https://www.rfc-editor.org/rfc/rfc8785). Implementations MUST NOT use the OTLP protobuf encoding itself for hashing &mdash; protobuf is non-deterministic for some field types and is unsuitable as a canonical form.

## 6. Storage

Implementations MUST persist captured events in append-only form. UPDATE and DELETE operations on stored events are non-conformant. Retention period is set by regulatory framework and tenant configuration.

The Merkle tree's intermediate nodes MAY be re-computed at verification time rather than stored. Storage of the per-day Merkle root and its HSM signature is REQUIRED.

## 7. Verification

A conforming verifier:

1. Re-computes every event's `payload_hash` using the published `prev_hash` and canonical JSON, and compares to the stored `payload_hash`.
2. Re-computes the daily Merkle root from the ordered set of `payload_hash` values.
3. Verifies the Ed25519 signature over `sign_payload` using the institution's published public key.
4. Reports per-day pass/fail with the specific failure mode if any step fails.

A verifier MUST fail-closed: if any step cannot be evaluated unambiguously, the day is reported as failed.

## 8. Conformance test vectors

[`test-vectors/`](test-vectors/) contains the canonical conformance corpus. A conforming implementation produces output identical to the test vectors for the cases they cover. Implementations SHOULD extend the corpus when they add features.

## 9. Security considerations

See [`docs/design/09-threat-model.md`](../docs/design/09-threat-model.md) for the threat model and adversary capabilities considered. The auditor's lens applies: any change that weakens an integrity property is a normative break, not an implementation choice.

## 10. References

Normative references:

- RFC 2119 / RFC 8174 — Conformance keywords
- RFC 6962 — Certificate Transparency Merkle tree construction
- RFC 8785 — JSON Canonicalization Scheme (JCS)
- RFC 5869 — HKDF
- RFC 2104 — HMAC
- FIPS 140-2 / FIPS 140-3 — HSM protection levels
- FIPS 180-4 — SHA-2 family
- FIPS 186-5 — Ed25519
- OpenTelemetry Specification, OTLP
- OpenTelemetry Semantic Conventions for Generative AI (gen_ai.*)

Informative references:

- FFIEC IT Examination Handbook — Information Security booklet (Sept 2016)
- FFIEC IT Examination Handbook — Architecture, Infrastructure, and Operations booklet (June 2021)
- Federal Reserve SR 11-7 — Model Risk Management
- OCC Bulletin 2011-12
- U.S. Treasury Financial Services AI Risk Management Framework (Feb 2026)

## 11. Change log

| Version | Date | Change |
|---|---|---|
| v1.0-draft | 2026-05-06 | Initial draft. Four primitives defined. Wire format normative. Test vectors planned. |
