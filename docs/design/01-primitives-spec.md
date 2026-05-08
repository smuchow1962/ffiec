# 01 — The four primitives, design rationale

> **What this doc is.** For each of the four primitives, the alternatives we considered, the choice we made, and why an auditor would accept it.

## 1. Why four, not three or five

The set is sized to satisfy four distinct adversaries and four distinct verification questions. Each primitive defends against exactly one adversary class. Removing any primitive leaves a class undefended; adding more does not strengthen any existing defense.

| Adversary | Primitive that defends |
|---|---|
| External attacker forging events into the wire | HMAC chain at capture |
| Insider with database access altering past events | Daily Merkle seal |
| Insider with operational access forging the seal | HSM-rooted root signature |
| Vendor lock-in restricting independent verification | OpenTelemetry-native wire |

## 2. Primitive 1 — HMAC chain at capture

### 2.1 What it does

Binds every event to a sequence-of-events for a given run. Produces evidence that an event in a captured run is part of a deliberate, ordered set, and that an external attacker without the session key cannot insert plausible events.

### 2.2 Alternatives considered

| Alternative | Why rejected |
|---|---|
| Plain SHA-256 hash chain (no key) | An attacker who can write to the wire can construct a forged chain. SHA-256 alone proves nothing about authorship. |
| Per-event signature with asymmetric keys | Adds 64 bytes per event for an Ed25519 signature, plus a key-management problem at every host. The chain solves the same problem at lower cost. |
| MAC with AES-CMAC | FIPS-approved, but less common in Go stdlib. HMAC-SHA-256 is the same security level and stdlib-native. |
| Merkle leaf per event without a chain | Loses the per-run ordering property. The chain captures order in the prev-hash; a bare Merkle leaf doesn't. |

### 2.3 The choice — HMAC-SHA-256 with HKDF-derived session keys

- **HMAC-SHA-256** &mdash; FIPS-approved (FIPS 198-1), stdlib in every modern language, audited extensively. The auditor recognizes it without explanation.
- **HKDF-SHA-256 (RFC 5869)** &mdash; key derivation for the per-process session key. The master key never appears on the application host. A compromised process loses its session key but not the master.
- **prev_hash binding** &mdash; ties events into a tamper-evident sequence. Removing or inserting one event is detectable without re-running the keys.

### 2.4 The session-key lifecycle

- A process receives its session key at startup via a tenant-trusted handshake (out of scope here; documented per implementation).
- The key is held in process memory only.
- The key is destroyed on process exit.
- A new process (after a crash, restart, or scale event) gets a new session key via the same handshake.
- The session key is bound to `(tenant_id, process_uuid)` so two processes within the same tenant cannot accidentally chain their events together.

### 2.5 The auditor's-lens review

| Question | Answer |
|---|---|
| Why per-process keys instead of per-tenant? | Compromise of one process should not compromise other processes. Per-process scoping is the standard separation. |
| What stops the master key from leaking? | The master never reaches the application host. The handshake delivers a session key derived under HKDF; the master remains in tenant-controlled storage (HSM-backed where available). |
| What if HKDF is broken in the future? | The chain is still useful retrospectively because each event also includes its `payload_hash` under HMAC. A key compromise affects new events; sealed past events are still verifiable against the daily Merkle root. |
| Outstanding work | The handshake protocol that delivers session keys is out of scope for v1.0 — institutions document their handshake in CC8.1 and the chain inherits the institution's existing key-delivery posture. v1.x roadmap candidates: SPIFFE/SPIRE-based delivery, or vendored HSM-token-based delivery. Either is layered onto v1.0 without breaking the wire form. |

## 3. Primitive 2 — Daily Merkle seal

### 3.1 What it does

Catches retroactive tampering. An insider with database access who alters one event must also alter the Merkle root that depends on that event, plus the HSM-signed seal over that root. The Merkle structure is the witness that a single byte changed.

### 3.2 Alternatives considered

| Alternative | Why rejected |
|---|---|
| No periodic seal &mdash; rely on HMAC alone | HMAC chains can be re-computed by anyone who recovers the session key. Without a separate sealing step, an insider with the key can rewrite history. |
| Per-event signature instead of a daily aggregate | Costs 64 bytes per event &times; billions of events. The aggregate is the same security at materially lower cost. |
| Merkle seal per hour or per minute | Higher signing frequency increases HSM-operational cost and audit-burden without a corresponding security gain. The daily granularity matches the typical examination cadence. |
| Blockchain-style proof-of-work | Energy and complexity not justified for a private-tenant ledger. The HSM signature provides equivalent integrity. |

### 3.3 The choice &mdash; RFC 6962 binary Merkle, daily, leaf-prefix

- **RFC 6962** is the Certificate Transparency Merkle scheme. Battle-tested. Auditor precedent.
- **Leaf-prefix and node-prefix** (`0x00` and `0x01`) prevent second-preimage attacks where an internal-node hash could be presented as a leaf or vice versa.
- **Daily** matches retention period definitions and examination cadence.

### 3.4 Edge cases

- **Empty days.** A tenant-day with zero events still gets a seal. The Merkle root for an empty leaf set is `SHA-256("")`. The seal is published to maintain continuity of the seal chain.
- **Late-arriving events.** Events that arrive after the daily seal is sealed are recorded with a `late_binding` flag and included in the next day's seal. The original day's seal is not altered. The verifier reports late-binding events explicitly.
- **Out-of-order events.** Events are ordered for Merkle construction by `(run_id, seq)` ascending. Implementations must NOT use receive timestamp for Merkle ordering &mdash; that would make the seal non-deterministic across implementations.

### 3.5 The auditor's-lens review

| Question | Answer |
|---|---|
| Why daily? Could an attacker do damage in less than a day? | The HMAC chain catches in-flight tampering instantly; the daily seal catches retroactive tampering at the day boundary. The window is bounded by the seal interval. Higher frequency is available as an option (configurable per tenant). |
| Can the cadence be relaxed for smaller institutions? | Yes. Daily is the default; weekly or monthly is defensible for community banks where examination cadence is 18-month and HSM cost is meaningful. The trade-off is the time-to-detect for retroactive tampering — institutions document the cadence in their control description and the examiner accepts or pushes back. |
| What if the seal job fails? | Events continue to be captured and HMAC-chained. The daily seal is computed and signed when the HSM becomes available. The seal&rsquo;s timestamp records the actual signing time; the day boundary is unambiguous. The verifier reports the delay. |
| Why RFC 6962 specifically? | It is the standard Merkle construction in regulator-trusted systems (Certificate Transparency, CONIKS, Trillian). Auditors recognize it. Re-using a standard is a feature. |
| What does it mean to "seal" no events? | Resolved in spec §4.2 (normative). Every tenant-day MUST receive a seal record, including tenant-days with zero events; the empty-day Merkle root is pinned at `SHA-256(b"")` = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`. A missing empty-day seal is reported as `missing seal for tenant-day {D}` — control-completeness failure, not chain-integrity failure. The seal is a continuous attestation that the institution is operating the ledger. |

## 4. Primitive 3 — HSM-rooted root signature

### 4.1 What it does

Anchors integrity in hardware. Even an attacker with full database admin and full process compromise cannot forge a daily root because they cannot access the HSM signing key.

### 4.2 Alternatives considered

| Alternative | Why rejected |
|---|---|
| Software-held signing key | A compromised host can extract the key. The integrity model collapses. |
| TPM-backed key | TPM is fine for individual workstation authentication; banks need shared, auditable, FIPS 140-2 L3-or-higher devices. |
| Multi-party signing (threshold Ed25519) | Operationally complex. The simpler HSM-with-redundancy model satisfies the same regulator concern. |
| Bitcoin-style timestamping | External dependency. Banks resist external dependencies in their integrity path. |

### 4.3 The choice &mdash; Ed25519 in a FIPS 140-2 L3 (or higher) HSM

- **Ed25519** &mdash; FIPS 186-5 approved (since 2023), modern, fast, deterministic signatures (no nonce issues), 64-byte signatures, 32-byte keys. Stdlib in Go.
- **FIPS 140-2 Level 3 minimum** &mdash; physical tamper resistance plus identity-based authentication. Cloud HSMs at this level: AWS CloudHSM Classic (FIPS 140-2 L3), Azure Key Vault Premium (FIPS 140-2 L2 default; L3 available), Google Cloud HSM (FIPS 140-2 L3). On-prem options: Thales, Entrust, Utimaco, Yubico YubiHSM 2.
- **Per-tenant keys** &mdash; one signing key per tenant. Cross-tenant signing is structurally prevented.

### 4.4 The signing payload

The signed payload is constructed as text rather than binary, to maximize auditor readability. Under v1.0-final-amendment the canonical form is the v1.0a 10-line `sign_payload` defined normatively in spec §4.3:

```
ffiec.chain-of-custody.v1\n
{sign_payload_version}\n            // "v1.0a" under v1.0-final-amendment
{algorithm}\n                       // "ed25519" for v1.0
{format_version}\n                  // "v1" for this spec
{tenant_id}\n
{ISO 8601 date (YYYY-MM-DD UTC)}\n
{hex-encoded Merkle root, 64 chars lowercase}\n
{hex-encoded HKDF inputs digest, 64 chars lowercase}\n
{cadence}\n                         // "hourly" | "daily" | "weekly"
{dev_mode}                          // single byte: "1" or "0"; no trailing \n
```

Each inter-field separator is a single `0x0A` (LF) byte; the terminal `dev_mode` byte has no trailing newline. The payload is exactly nine `0x0A` separators. CRLF is non-conformant. An auditor can read the payload, copy it into a verifier, and confirm the signature by hand if necessary. Binary encoding would be marginally more efficient and substantially less inspectable.

Pre-amendment chains (produced before 2026-05-07) omit the `sign_payload_version` line and use the pre-amendment 6-line form (magic line + `algorithm` + `format_version` + `tenant_id` + `iso8601_date` + `hex(merkle_root)` + `hex(hkdf_inputs_digest)`). The verifier dispatches on the seal record's `sign_payload_version` field and reconstructs the matching form, so pre-amendment chains remain verifiable under amendment-aware verifiers without re-sealing. Spec §4.3 carries the byte-level normative text.

### 4.5 The auditor's-lens review

| Question | Answer |
|---|---|
| Who holds the signing key? | The tenant's HSM. The vendor never holds the key. In vendor-hosted topologies, the HSM is per-tenant or supports per-tenant key labels enforced by HSM ACLs. |
| Who holds the public key? | The tenant publishes the public key to a tenant-controlled registry. The institution shares it with examiners and external auditors. The verifier reads it from a flat file or from the registry. |
| What if the HSM key is lost? | A tenant rotation event. The new HSM key signs going forward; previous days' seals are still verifiable against the previous public key. The institution must document key rotation in the tenant key registry. |
| What if FIPS 140-2 is replaced by FIPS 140-3? | FIPS 140-3 supersedes 140-2 for new validations. The spec accepts &ldquo;FIPS 140-2 Level 3 or higher&rdquo;, which includes FIPS 140-3 Level 3. |
| Tenant key registry format | v1.0 leaves the registry implementation-flexible — &ldquo;institution provides a documented retrieval path&rdquo; — because every institution already has an internal directory shape its IT and compliance teams know. A standardized JWKS-style registry endpoint with FFIEC-specified extensions is a v1.x roadmap commitment that lifts the lowest-friction shape into spec text once enough institutions have shipped to inform the field set. |

## 5. Primitive 4 — OpenTelemetry-native wire

### 5.1 What it does

Lets banks integrate the chain-of-custody primitives into their existing observability infrastructure without a rip-and-replace. Lets vendors compete on quality of implementation without forcing the institution to commit to one stack.

### 5.2 Alternatives considered

| Alternative | Why rejected |
|---|---|
| Custom binary protocol | Forces every vendor and every backend to implement decode logic. OTel solved this. |
| JSON-over-HTTPS, hand-rolled schema | Same problem as custom binary. |
| Apache Avro / Apache Thrift | Less common in observability. OTel is the de facto standard. |
| Use OTel logs instead of OTel traces | Logs lack the run/span hierarchy needed for agent execution. Traces are the right shape. |

### 5.3 The choice &mdash; OTLP with `ffiec.chain.*` attribute extension

- **OTLP** is the wire format the rest of observability is converging on. Splunk ingests it. Sentinel ingests it. Datadog ingests it. Banks already deploy OTel collectors.
- **GenAI semantic conventions** carry the AI-specific fields (`gen_ai.system`, `gen_ai.request.model`, etc.) without our spec defining them.
- **`ffiec.chain.*` namespace** carries the integrity-extension fields. The namespace is requested for inclusion in the OTel semantic conventions registry.

### 5.4 The escape hatch

The wire format is OTLP. Implementations may export to *any OTLP-compatible backend*. The chain extension fields pass through as ordinary span attributes. A bank running Honeycomb today doesn&rsquo;t need to migrate &mdash; the chain attributes ride along with their existing observability stream.

### 5.5 The auditor's-lens review

| Question | Answer |
|---|---|
| Why depend on OTel? | OTel is a CNCF-graduated, industry-standard, multi-language, vendor-neutral observability framework. Depending on it is depending on the lingua franca. |
| What if OTel changes? | The chain extension fields are encoded as plain attribute key/value pairs. Future OTel versions either preserve them or deprecate the attribute mechanism, in which case our spec migrates with them. |
| Are we registered with the OTel semantic conventions? | Not yet. Registration of the `ffiec.chain.*` and `audit.*` namespaces under FFIEC stewardship is a v1.0-final-amendment post-submission milestone — the spec is locked at v1.0-final-amendment, the conformance corpus is built, and the registration request goes to the OTel semconv community once an FFIEC-side sponsor is identified. The chain attributes are stable in v1.0; registration is paperwork, not a wire-format change. |
| Outstanding work | The registration step. Track in [`spec/`](../../spec/) once a sponsor in the OTel community is identified. |

## 5.6 Algorithm identifiers and migration

The seal record carries an explicit algorithm identifier so a future spec version can introduce a post-quantum signature alongside Ed25519 without invalidating past seals. Spec §4.3.2 already admits the dual-algorithm transitional posture via the seal record's `signatures` list (Variant B, per-algorithm `sign_payload`); the verifier dispatches on the algorithm identifier and the algorithm-confusion attack class is closed before the v1.x post-quantum transition begins. The HMAC chain has no per-event algorithm identifier in v1.0 — HMAC-SHA-256 is implicit. The optional attribute `ffiec.chain.algorithm` (defaulting to HMAC-SHA-256 for v1.0 events) is part of the spec §4.4 attribute table; it costs nothing to emit today and makes a future migration structurally tractable. Without it, an HMAC-SHA-256 weakening would require a spec version bump rather than an attribute addition.

## 6. Cross-primitive review

The four primitives compose. The composition matters:

- HMAC chain alone catches wire-tampering but not retroactive insider abuse.
- Merkle seal alone catches retroactive abuse but not wire-tampering.
- HSM signature alone is integrity for the *seal*, but the seal must cover something tamper-resistant or the signature is over fiction.
- OTLP wire alone is just a pipe; without the integrity primitives, it&rsquo;s just bytes.

The composition is what an auditor accepts. Each defends one adversary class; together they cover the threat model in [`09-threat-model.md`](09-threat-model.md).
