# 02 — Chain construction (the hot path)

> **What this doc is.** The detailed design of the HMAC chain construction step. This is the hot path &mdash; runs once per captured event, in the application process, before the event leaves the host. Performance, determinism, and crash-safety are all in scope.
>
> **Reference implementation.** Herald.Py v1.0 (`src/herald/_crypto/chain.py` and `_crypto/keys.py`). Three Auditor rounds against the implementation; every BLOCKER and MAJOR finding closed. The cross-language byte-level test vectors at `tests/fixtures/chain_vectors.json` are the byte-level contract any conforming implementation reproduces verbatim. The FFIEC `spec/test-vectors/` corpus is seeded from that fixture's shape.

## 1. The contract

For each event, the SDK produces a chain entry whose `payload_hash` is defined by spec §4.1:

```
ikm                  = the tenant's IKM (raw bytes; minimum 32 per spec §10.6)
tenant_id            = the event's ffiec.chain.tenant_id (UTF-8 string)
HKDF_SALT            = b"ffiec.chain-of-custody.v1.salt"             // FFIEC-conformance constant; vendor-flag mode permitted per spec §4.1.2
HKDF_INFO_BASE       = b"ffiec.chain-of-custody.v1.info"             // FFIEC-conformance constant; vendor-flag mode permitted per spec §4.1.2
info_for_tenant      = HKDF_INFO_BASE || b"|" || utf8(tenant_id)
session_key          = HKDF-SHA-256(IKM=ikm, salt=HKDF_SALT, info=info_for_tenant, length=32)
key_fingerprint      = SHA-256(utf8(tenant_id) || ikm)[:16]          // 16 raw bytes
canonical            = canonical_json(event_payload minus chain-stamp fields)   // RFC 8785 JCS
prev_hash            = 32 zero bytes for seq=1; previous payload_hash for seq>1 (raw 32 bytes)
payload_hash         = HMAC-SHA-256(session_key, prev_hash || canonical)        // raw 32 bytes
```

`||` is byte concatenation. `payload_hash`, `prev_hash` are exactly 32 raw bytes; `key_fingerprint` is exactly 16 raw bytes.

**Per-entry stamp.** Each chain entry persists `(seq, prev_hash, payload_hash, key_version, key_fingerprint, format_version, mac_computed_at_utc, kms_handle_uri)`. The MAC IS the `payload_hash` field, persisted byte-for-byte. The MAC input excludes the chain-stamp fields; the field carrying the MAC does not self-reference.

**Vendor-flag mode (spec §4.1.2).** A vendor SDK MAY accept `HKDF_SALT` and `HKDF_INFO_BASE` as construct-time configuration so a single codebase serves both FFIEC-conformance posture and a vendor-namespaced posture (internal product telemetry, non-banking deployments, regulator-aligned-but-not-FFIEC jurisdictions). The byte values shown above are the FFIEC-conformance values. Chains produced under any other constants are non-FFIEC chains by definition; the audit-file header's `hkdf_inputs_digest` records which constants were in force, and a verifier walking under FFIEC posture refuses non-FFIEC chains at the file-header pre-flight (spec §7 step 2). Posture is binary at the chain-file level — every entry under a given file header inherits the posture; the SDK does not toggle posture mid-file. An SDK exposing both postures documents the configuration mechanism (constructor argument, factory function, environment variable) so the institution's CC8.1 control description names the posture in force for FFIEC-supervised activity.

**`chain_kind` enumeration (spec §3, closed at v1).** The `chain_kind` field on every chain entry classifies the entry's event class. The v1 enumeration is closed: `"audit"` (default — application audit event), `"model_call"` (LLM invocation), `"tool_call"` (tool invocation), `"routing"` (routing-decision entry per spec §4.4.1), `"translation"` (ECOA translation entry per spec §10.11), `"operational"` (control-evidence operational events). The closed enumeration is what lets two SDKs producing the same logical event under the same canonicalization rules emit byte-identical canonical bytes — an open or implicit `chain_kind` would let one SDK emit `"audit"` and another `"audit.event"` for the same logical event and break the cross-vendor byte-identity claim. The discriminator is distinct from the OTel `kind` field (which carries the OpenTelemetry `SpanKind` — `server`, `client`, `internal`, `producer`, `consumer`); both fields are integrity-bound under the per-event MAC and both appear on every entry. The verifier MUST reject any value not in the enumerated set with `chain_kind out of v1 enumeration at seq N`. Cross-reference spec §3 for the full normative text.

**`sign_payload` v1.0a canonical form (spec §4.3).** The §4.3 `sign_payload` carries every byte under the HSM signature. Under v1.0-final-amendment the canonical form is the 10-line v1.0a layout: line 1 is the magic line `"ffiec.chain-of-custody.v1"`, line 2 is `sign_payload_version = "v1.0a"` (the discriminator that lets future amendments add lines without breaking pre-amendment verifiability), line 3 is `algorithm`, line 4 is `format_version`, line 5 is `tenant_id`, line 6 is `iso8601_date(tenant_day)`, line 7 is `hex(merkle_root)`, line 8 is `hex(hkdf_inputs_digest)`, line 9 is `cadence` (one of `"hourly"` | `"daily"` | `"weekly"`), and line 10 is `dev_mode` (a single ASCII byte: `"1"` when true, `"0"` when false or absent). Each inter-field separator is a single `0x0A` byte; the terminal `dev_mode` byte has NO trailing `\n`. Total LF separators: nine.

The form binds `sign_payload_version`, `cadence`, and `dev_mode` under the HSM signature. Each binding closes a path that an earlier form left open. `sign_payload_version` is the form-generation discriminator: a tampered value at line 2 is detected at signature verification because the verifier reconstructs the form using the field as written, so a mismatch between the field and the byte form actually used by the signer produces a signature failure. `cadence` rewrite (e.g., flipping `daily` to `weekly` to claim a relaxed posture) is now caught cryptographically rather than via cross-document reconciliation against the institution's CC8.1 control description. `dev_mode = true → false` is the more serious case: an attacker with seal-record write access could otherwise present a chain produced by the §10.7 development software-key adapter as a production chain, defeating the regulator-visible-line guarantee that §10.7 establishes. Binding `dev_mode` under the HSM signature means a `dev_mode` flip now requires forging the HSM signature, which the FIPS 140-2 Level 3 custody posture rules out.

Pre-amendment chains (produced before 2026-05-07) omit the `sign_payload_version` line and use the pre-amendment 6-line form; the verifier dispatches on the seal record's `sign_payload_version` field per §4.3 and reconstructs the matching layout. Pre-amendment chains remain verifiable under amendment-aware verifiers without re-sealing.

The `sign_payload` byte-format requirements are byte-for-byte normative: each inter-field separator is a single `\n` (0x0A) byte; the terminal field has NO trailing `\n`; all hex-encoded values use lowercase ASCII digits and lowercase a-f (uppercase hex is non-conformant) and are zero-padded to 64 characters; all line terminations are LF (0x0A) only, never CRLF; the `iso8601_date(tenant_day)` is `YYYY-MM-DD` with no time component, no timezone, no fractional seconds. Implementations on every platform MUST emit `0x0A` when computing `sign_payload` — platform-default newline conversion in serialization libraries (e.g., Go's `text/template`, Windows-host text-mode file writers) MUST be disabled on this hot path. Per-algorithm `sign_payload` under dual-algorithm posture (spec §4.3.2, Variant B) means each algorithm's signature is computed over its OWN algorithm-bound `sign_payload` — line 3 (`algorithm`) carries that algorithm's identifier; line 2 (`sign_payload_version`) remains `"v1.0a"` for both algorithms. A single shared `sign_payload` covering all algorithms is non-conformant because it leaks the algorithm-confusion defense. Cross-reference spec §4.3 for the byte layout and §4.3.2 for the dual-algorithm rule.

**`key_versions` verifier cross-check (spec §7 step 11).** After signature validation the verifier MUST cross-check the seal record's `key_versions` list against the actual `key_version` distribution observed in the day's chain entries: `seal.key_versions == sorted(set(entry.key_version for entry in day_events))`. Mismatch reports `seal.key_versions does not match per-event key_version distribution`. The cross-check exists because `key_versions` is a variable-length list whose direct binding under `sign_payload` would complicate signature reproducibility (a sorted-list serialization is straightforward but adds wire-format surface that other consumers cross-check independently); the cross-check provides the equivalent integrity property by independent assertion after signature validation. A coordinated forgery rewriting both the seal's `key_versions` and a per-entry `key_version` would still be caught by step 8 (fingerprint mismatch) — the adversary would need to produce an IKM-fingerprint pair under the rewritten `key_version` that matches the entry's stamped fingerprint, which the IKM lookup at step 7 either resolves to the correct IKM (and step 8 fails) or fails to resolve at all (step 7 fails with `unknown_key_version`). The verifier's defense is composite: signature validation establishes the seal's authenticity, the `key_versions` cross-check establishes the seal's `key_versions` list is consistent with the per-event distribution, and the per-entry fingerprint check establishes the IKM lookup is consistent with the entry's claimed identity. Cross-reference spec §7 step 11 for the normative text.

**Why each per-entry field exists.**

| Field | Why it lands on every entry |
|---|---|
| `seq` | Monotonic per-run ordering. Verifier asserts `seq + 1` walk; gaps or reversals are rejected. |
| `prev_hash` | The chain link. Verifier asserts `entry.prev_hash == previous payload_hash`. Fixed 32 bytes for length-prefix unambiguity in the MAC input. |
| `payload_hash` | The MAC, persisted byte-for-byte. **This IS the seal.** A writer that computes the MAC and discards it provides zero tamper-evidence. |
| `key_version` | Names the IKM generation. Lets the verifier resolve `(tenant_id, key_version) -> IKM` across rotations. |
| `key_fingerprint` | The per-entry tenant↔IKM identity binding. Verifier asserts the looked-up IKM produces this fingerprint BEFORE computing any MAC. Catches botched rotation (wrong IKM under same `key_version`) at lookup time, not buried in MAC-mismatch storms. |
| `format_version` | Lets the verifier refuse files of the wrong format most-specific-first. A v2 file fails with "format_version not supported" rather than "HKDF inputs do not match v1". |
| `mac_computed_at_utc` | Forensic-only. The verifier does NOT trust this for security decisions. Useful for after-the-fact incident review. |
| `kms_handle_uri` | Provenance pointer (e.g. `"aws-kms:arn:..."`, `"plaintext-dev"`). Auditor cross-references against IKM custody documentation. |

## 2. The canonical-form choice

### 2.1 Why JCS over JSON

[RFC 8785 (JSON Canonicalization Scheme)](https://www.rfc-editor.org/rfc/rfc8785) defines a deterministic JSON encoding:

- Object keys are sorted by Unicode codepoint
- Whitespace is normalized
- Numbers are encoded in a defined form
- String escaping is normalized

Two implementations of JCS produce byte-for-byte identical output for the same logical event. That is the determinism property the chain depends on.

### 2.2 Why not protobuf

The OTLP wire format is protobuf. Protobuf is non-deterministic for several constructs:

- Map field ordering is not specified
- Unknown fields may be re-ordered
- Field encoding for some types is not unique (varint vs zigzag, for instance)

A chain that hashes protobuf encoding produces different `payload_hash` values on different implementations. Two conforming implementations must agree byte-for-byte; protobuf cannot guarantee that. JCS can.

**The chain hashes JCS-canonical JSON; the wire transports OTLP protobuf carrying the chain fields as opaque bytes.** The two are independent.

### 2.3 What goes into the canonical payload

**Included** (the MAC covers these):

- The OTel envelope: `trace_id`, `span_id`, `parent_span_id`, `name`, `timestamp_ns`, `duration_ns`, `attributes`, `resource`, `severity`, `kind`, `chain_kind`
- The GenAI semconv attributes: `gen_ai.*`, `tool.*`
- The application audit-event payload: `audit.*`
- The explicit metadata: `tenant_id`, `run_id`, `captured_at`
- The cross-run linkage fields: `parent_run_id`, `parent_seq`, `dag_parents` (when present)

**Excluded** (the MAC input excludes these to avoid self-reference and to allow chain-stamp metadata to be added without changing the MAC):

- `seq` — derived from per-run ordering; not in the canonical bytes
- `prev_hash` — concatenated outside the canonical bytes (`prev_hash || canonical`)
- `payload_hash` — the MAC itself
- `key_version`, `key_fingerprint`, `format_version`, `mac_computed_at_utc`, `kms_handle_uri` — chain-stamp metadata
- `algorithm` — the MAC algorithm identifier itself

**Note: cross-run linkage fields (`parent_run_id`, `parent_seq`, `dag_parents`) are INCLUDED in the canonical bytes** — they are substantive evidence about the agent's decision graph, not chain-stamp metadata. Including them under the MAC binds the cross-run linkage with the same integrity property the per-event payload enjoys.

Implementations include exactly the documented fields. Adding fields silently breaks chain agreement across vendors. The conformance test vectors (spec §8 + `spec/test-vectors/`) lock the exact byte sequence the canonical bytes must produce.

**Why the chain-stamp fields are excluded.** The MAC IS one of the chain-stamp fields. Including them in the canonical bytes would create a circular dependency. The chain-stamp fields' tamper coverage comes from elsewhere:

- `key_version` flip → verifier looks up the wrong IKM → wrong session key → MAC fails. Detected.
- `key_fingerprint` flip → verifier recomputes fingerprint from looked-up IKM and constant-time compares; mismatch detected before any MAC compute (closes the botched-rotation failure mode). Detected.
- `format_version` flip → verifier asserts `header.format_version == entry.format_version`; mismatch refused. Detected.
- `mac_computed_at_utc` / `kms_handle_uri` flip → forensic fields, not security-relevant to the seal. The verifier does not trust them for security decisions; they exist for after-the-fact incident review.

## 3. Persistence ordering &mdash; crash safety

The hot path must guarantee one property: **for any event whose `payload_hash` has been disclosed, the event is persisted to disk**.

Disclosure happens in three ways:
- The event is exported over OTLP
- The event's `payload_hash` is used as the next event's `prev_hash`
- A user reads the event from the SDK's API

The SDK enforces ordering by writing-then-disclosing:

```
1. Compute payload_hash
2. Insert into local SQLite (with synchronous=FULL for audit events,
   synchronous=NORMAL for telemetry)
3. Wait for fsync to complete
4. Update the in-memory "last hash for this run" cache
5. Return to caller
6. (Async) Export over OTLP
```

A crash between step 2 and step 4 is recoverable: the SDK reads the persisted event on restart and rebuilds the in-memory cache. A crash before step 2 means the event never existed, which is correct.

## 4. The session-key handshake

At process startup, before any event is captured. Two delivery models are conformant per spec §4.1.1; both produce a per-tenant session key the SDK can use for HMAC compute.

### 4.0 Model A — IKM-delivered (most common)

```
1. SDK reads the tenant_id from configuration (or from per-event tenant scope).
2. SDK authenticates to the IKM custodian (HSM/KMS) via SPIFFE/SPIRE, mTLS, or HSM-issued bearer token.
3. Custodian validates the workload identity, returns the IKM bytes over an authenticated TLS 1.3 channel.
4. SDK derives, for each tenant the SDK serves:
   info_for_tenant = HKDF_INFO_BASE || b"|" || utf8(tenant_id)
   session_key     = HKDF-SHA-256(IKM=ikm, salt=HKDF_SALT, info=info_for_tenant, length=32)
   key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]
5. SDK caches (session_key, key_fingerprint) per-tenant in a bounded LRU.
6. The IKM bytes remain in process memory for the process lifetime (best-effort zeroisation
   on shutdown; honesty note in §4.3 on memory-zeroisation in GC runtimes).
```

The IKM-delivered model is what Herald.Py and most reference implementations use today. The SDK can serve any tenant whose IKM it has fetched; cache misses on a new tenant trigger a fresh derivation (microseconds).

### 4.0a Model B — Session-key-delivered (HSM-mediated, no IKM in process memory)

For environments where holding IKM bytes in process memory is unacceptable (capability-restricted PaaS, strict zero-trust postures):

```
1. SDK reads the tenant_id from configuration.
2. SDK authenticates to the HSM (PKCS#11, KMIP, or vendor SDK).
3. The HSM performs HKDF inside the device with the SDK-supplied tenant_id:
   - HSM uses the IKM bound to its key label
   - HSM computes info_for_tenant and runs HKDF-SHA-256 internally
   - HSM returns the 32-byte session_key + 16-byte key_fingerprint
4. SDK never sees the IKM bytes. HMAC operations may run locally in the SDK process
   (the session_key is in SDK memory) or dispatch through the HSM API for fully-isolated
   key custody.
```

This satisfies the same per-tenant determinism property: the same `(IKM, tenant_id)` pair always produces the same `(session_key, key_fingerprint)`, so the verifier can recompute the MAC from the IKM alone — the verifier need not know which delivery model the SDK used.

### 4.1 Handshake security floor (normative for v1.0-final)

A weak handshake undoes everything: an attacker who can intercept session-key delivery (Model B) or IKM delivery (Model A) produces events that pass HMAC verification. Implementations MUST satisfy four properties regardless of model:

- **Authentication.** The application host MUST authenticate to the IKM custodian using a cryptographic credential bound to the host's identity (workload-identity attestation, mTLS client certificate, HSM-issued token, or equivalent).
- **Confidentiality.** The IKM (Model A) or session key (Model B) MUST be transported over an authenticated, confidential channel. TLS 1.3 with appropriate cipher suites is the minimum; HSM-mediated key wrap is preferred where available.
- **No-cache at the custodian.** The custodian MUST NOT cache derived session keys across handshakes. (The SDK MAY cache per-tenant within a single process; bounded LRU is RECOMMENDED.)
- **Per-tenant determinism.** Given the same IKM and the same `tenant_id`, the derivation MUST produce a byte-identical `session_key` and `key_fingerprint` across processes, hosts, and time. The verifier depends on this property.

**SPIFFE/SPIRE is the recommended mechanism** for zero-trust deployments. Other mechanisms (HSM-issued bearer tokens, mTLS, hardware-attested workload identity) are conformant if they satisfy the four properties above.

**Bare-metal handshake example (no orchestrator).** For institutions running on bare-metal where SPIFFE/SPIRE is not available:

1. The host holds a long-lived mTLS client certificate, issued at host commissioning by the institution's PKI and bound to the host's hardware identity (TPM-backed where available).
2. At process start, the SDK performs an mTLS handshake with the IKM custodian, presenting the host certificate.
3. The custodian validates the certificate against the institution's PKI, confirms the host is on the active roster, and (Model A) returns the IKM, OR (Model B) issues a short-lived (10-minute) one-shot HSM token bound to a fresh derivation request.
4. The SDK proceeds per Model A or Model B above.
5. Tokens (Model B) are invalidated at first use; replaying them fails.

This satisfies all four properties: authenticated (mTLS + HSM token), confidential (TLS 1.3), no-cache (single-use token or per-call derivation), per-tenant deterministic (the HKDF inputs are fixed). Audit evidence is the institution's PKI logs (certificate validations) plus the HSM access log.

### 4.2 Memory protection per platform

Session keys (and, in Model A, IKM bytes) are held in process memory only. Platform-specific protection:

| Platform | Mechanism | Implementation note |
|---|---|---|
| Linux | `mlock(2)` on the page containing the key | Set `RLIMIT_MEMLOCK` adequately; root or `CAP_IPC_LOCK` may be required outside containers |
| Windows | `VirtualLock` | Process working-set limits apply |
| macOS | `mlock(2)` | Same as Linux |
| Containers | Inherit host policy | `--cap-add IPC_LOCK` on Docker; `securityContext.capabilities.add` on Kubernetes |
| Hardware-backed key storage (Model B) | HSM-mediated key handle, not raw key in memory | Preferred where available; the SDK never sees the IKM |

Implementations MUST document their memory-protection posture for each supported platform. The verifier's report includes a section confirming the institution's memory-protection claims, which the examiner cross-checks against the implementation's documentation.

**Compensating controls when IPC_LOCK is forbidden.** Many bank container platforms (PaaS, managed Kubernetes with strict pod security) forbid privileged capabilities. In those environments the application CANNOT mlock pages. The conformant alternative is Model B: the SDK uses an HSM-mediated session-key handle (or HMAC-via-HSM) rather than holding raw key bytes in memory.

Implementations targeting capability-restricted environments MUST support Model B. Software-only Model A implementations on capability-restricted environments are non-conformant for production. Test/development environments may use software keys with the dev-mode marker (`kms_handle_uri = "plaintext-dev"`; spec §10.7).

### 4.3 Memory zeroisation honesty note

Calling `mlock` keeps the key out of swap; it does not erase it from memory. On process exit, the memory is freed but the bytes may persist briefly until reused. Implementations SHOULD attempt best-effort zeroisation on shutdown (`OPENSSL_cleanse` in C; `CryptographicOperations.ZeroMemory` in .NET) but the threat model assumes the IKM may persist in process memory for the process lifetime. Zeroisation is hardening, not a guarantee.

What actually bounds the damage:

- **Rotation cadence.** The IKM rotates on a schedule (typical: 90 days for high-value tenants, 365 days for standard).
- **Process restart.** A short-lived service worker derives a new session key per startup; an attacker who reads memory at hour 8 loses the key when the worker recycles at hour 12.
- **Session-key max age.** Long-lived processes SHOULD enforce a per-process re-handshake cadence (default: 24 hours) so memory residence of any one IKM is bounded.

CPython's `bytes` is immutable, so a Python-side IKM cannot be zeroed at all without C-extension tricks that are not portable. The Herald.Py reference implementation documents this residual risk explicitly.

## 5. Performance budget

The hot path runs once per captured event. AI agents emit on the order of 10&ndash;100 events per agent invocation. At a target throughput of 10K agent invocations per second (an upper bound for the largest banks), the hot path runs ~1M times per second.

Per-event cost budget on a single core:

| Step | Budget |
|---|---|
| JCS encoding | <50 &micro;s for typical event |
| HMAC-SHA-256 of `prev_hash` + canonical | <10 &micro;s for events under 4 KB |
| SQLite insert with `synchronous=FULL` | <500 &micro;s on SSD |
| Total wall-clock per event | <1 ms |

Note: the SQLite insert dominates. Implementations that need higher throughput batch SQLite writes (with the WAL fsync invariant preserved per batch). For audit events (`synchronous=FULL`), batching is bounded; for telemetry (`synchronous=NORMAL`), batching is liberal.

## 6. The hash-cache contract

The SDK maintains an in-memory cache of `(run_id) -> last_payload_hash`. The cache is the source of truth for `prev_hash` lookups during normal operation. The SQLite store is the source of truth across restarts.

Cache eviction policy:

- Evict on run completion (`@agent` decorator's exit point clears the cache for that run)
- Evict on inactivity timeout (configurable, default 1 hour)
- Hard cap on cache size (configurable, default 100K entries)

Cache miss behavior: read the most recent `payload_hash` for the run from SQLite, populate the cache, proceed. A cache miss costs a single SELECT plus the normal hot-path operations.

## 7. Concurrency

Multiple goroutines (or threads) within one process may capture events for different runs concurrently. The chain construction is **per-run sequential** but inter-run parallel.

Implementation guidance:

- A `sync.Map[run_id]*runState` holds per-run state.
- `runState` contains `last_payload_hash` and a per-run `sync.Mutex`.
- Acquiring the per-run mutex serializes events within one run.
- Inter-run parallelism is unbounded (limited only by SQLite write throughput).

Events within one run are by definition sequential because they share `run_id` and `seq` is monotonically increasing. The SDK rejects out-of-order `seq` values.

## 8. What can go wrong, and how it surfaces

### 8.0 Writer-side (SDK) failures

| Failure | Surface |
|---|---|
| Session key not yet delivered (process started, handshake pending) | SDK returns `ErrSessionKeyNotReady`; events are queued in a startup buffer up to a configurable cap. |
| IKM under 32 bytes (spec §10.6 violation) | SDK refuses to start the chain writer; raises `ErrIkmTooShort` at configure time. Never reaches the hot path. |
| Software-key adapter selected in production build (spec §10.7 violation) | Compile-time exclusion: production binaries do not link the software adapter. A configuration that names the software adapter at runtime fails closed with `ErrSoftwareAdapterNotInBuild`. |
| SQLite full or write fails | SDK returns `ErrPersistenceFailed`; event is lost only if the caller does not retry. The audit-event API always retries until success or a hard timeout. |
| Cache and SQLite disagree | SDK aborts the chain construction and emits an `ErrChainInconsistent` alert. This should never happen in practice; if it does, the host is suspect. |
| `prev_hash` does not match expected | SDK rejects the event. Indicates a bug or a tampering attempt. |
| `seq` skips or reverses | SDK rejects the event. Same indication. |

### 8.1 Verifier-side failures (per spec §7 ordered procedure)

The verifier walks the procedure in order; the first failure produces a specific named reason. The order matters: cheap rejections happen before expensive ones, and structural errors are reported before cryptographic ones.

| Step | Failure | Reason string | Closes which attack |
|---|---|---|---|
| 1 | `header.format_version != "v1"` | `format_version not supported by this verifier` | Future-version misread (most-specific-first) |
| 2 | `header.hkdf_inputs_digest` mismatch | `header HKDF inputs do not match running v1 inputs` | Format-drift / wrong constants |
| 3 | `header.genesis_hash` mismatch | `header genesis_hash does not match v1 constant` | Format-drift |
| 4 | `event.tenant_id != header.tenant_id` | `cross-chain lift detected at seq N` | Cross-chain lift (event copied from another chain) |
| 4 | `event.run_id != header.chain_id` | `cross-chain lift detected at seq N` | Same |
| 5 | `entry.format_version != header.format_version` | `format_version mismatch at seq N` | Format-drift mid-file |
| 6 | `entry.seq != expected_seq` | `seq out of order at seq N` | Reorder / gap / duplicate |
| 6 | `entry.prev_hash != expected_prev_hash` | `chain link broken at seq N` | Insertion / deletion |
| 7 | `ikm_lookup` returns null | `unknown key_version: no IKM for (tenant=T, key_version=V) at seq N` | Decommissioned key generation |
| 8 | `key_fingerprint` mismatch | `key_fingerprint mismatch at seq N: looked-up IKM does not match the entry's recorded fingerprint` | **Botched rotation / cross-tenant key swap** (caught at lookup, no MAC compute) |
| 9 | `payload_hash` MAC mismatch | `payload_hash MAC mismatch at seq N` | Payload tampering |
| 10 | Merkle root mismatch | `merkle root mismatch — ledger contents do not produce sealed root` | Server-side history rewrite (composed with §11) |
| 11 | Ed25519 signature invalid | `signature verification failed` | Server-side / HSM compromise (composed with §10) |
| 11 | `seal.key_versions != sorted(set(entry.key_version))` | `seal.key_versions does not match per-event key_version distribution` | Silent rewrite of the seal's `key_versions` list (caught by independent cross-check after signature validation per spec §7 step 11) |
| 12 | Cadence/dev-mode under `--strict` | `cadence mismatch` or `dev-mode seal in production verification — refused` | Misconfiguration |
| 12a | `gen_ai.{request,response}.model` missing on a model-call entry | `gen_ai_model_identifier_missing at seq N: {field_name} required for chain entries representing model calls` | Control-completeness for SR 11-7 reproducibility (NOT chain-integrity) |
| 3a | `header.tenant_id` violates §3 character class | `tenant_id violates §3 character class` | Non-conforming SDK or tampered header |

All failures are loud, named, and actionable. Silent failure is a defect.

### 8.2 Mid-write truncation

Per spec §4.1: every persisted chain-entry line MUST end with a single `\n` byte. The verifier refuses any file whose last byte is not `\n` with `audit file ends mid-line — possible mid-write crash`. A permissive verifier would silently pass a chain that lost its last entry; the strict refusal is what surfaces the crash.

**Implementation footgun.** Stdlib line readers in common languages strip the terminator from each returned line and silently tolerate a final-line-without-newline — Python's `for line in file:` and `readline()`, .NET's `StreamReader.ReadLine()`, Go's `bufio.Scanner.Scan()`, Java's `BufferedReader.readLine()`. A verifier that delegates straight to one of these readers will MISS the truncation case and incorrectly PASS a truncated chain. The verifier MUST therefore perform an explicit byte-level check before walking the file: open the file at byte level, seek to the last byte (POSIX `seek(-1, SEEK_END); read(1)`; .NET `stream.Seek(-1, SeekOrigin.End)`; Go `fs.Seek(-1, io.SeekEnd)`), and assert the byte equals `0x0A`. The check costs one syscall regardless of file size and runs once at pre-flight time. Spec §7 calls this implementation requirement out as normative.

The reference implementation in Herald.Py (`_crypto/chain.py:read_audit_file`) does the explicit byte-level check before line iteration and raises `AuditFileFormatError` rather than returning `ok=True`. The .NET implementation follows the same pattern.

## 9. Multi-process run semantics

Production agent platforms increasingly fan out one logical agent invocation across multiple processes (orchestrator + worker pattern; pre-processing + LLM call + post-processing; tool execution in a sandboxed process). The chain handles this with two patterns:

### 9.1 Per-process runs with parent linkage (preferred)

Each process generates its own `run_id` for events it captures. Processes that spawn other processes record a parent-child relationship as an event attribute:

```
ffiec.chain.parent_run_id = "r_orchestrator_a3f29b71c"
ffiec.chain.parent_seq    = 17
```

The chain within each `run_id` is independent; the parent-child relationship is a graph the verifier reconstructs at examination time. This is the preferred pattern because it keeps the per-process chain construction simple and lets each process verify locally.

### 9.2 Shared run_id across processes (acceptable for tightly-coupled handoffs)

When a process explicitly hands off control to another process and the chain must continue without a break, both processes share `run_id`. The handoff event carries the previous process's last `payload_hash` (as `prev_hash` of the first event in the receiving process) and the receiving process derives its own session key under its own tenant_id. The tenant_id and key_version stay constant across the handoff (same tenant); the per-tenant determinism property (spec §4.1.1 property 4) means both processes derive byte-identical session keys.

**Handoff event normative schema.** The handoff event MUST be emitted by the receiving process as the first event in its half of the run. It MUST carry:

| Field | Required | Notes |
|---|---|---|
| `audit.handoff.from_run_id` | yes | The sending process's `run_id` (same as receiving when they share `run_id`) |
| `audit.handoff.from_seq` | yes | The sending process's last `seq` value |
| `audit.handoff.from_payload_hash` | yes | The sending process's last `payload_hash` (32 raw bytes; equals the handoff event's `prev_hash`) |
| `audit.handoff.handoff_at` | yes | RFC 3339 UTC timestamp of the handoff |
| `audit.handoff.from_host` | optional | Forensic; the sending process's hostname or workload identity |

The handoff event itself participates in the chain like any other event: its `payload_hash = HMAC-SHA-256(session_key, prev_hash || canonical_bytes)` where `prev_hash` is `audit.handoff.from_payload_hash`. The verifier validates the handoff event by the same per-event procedure (spec §7 steps 4-9); the only requirement is that the receiving process's first `prev_hash` equals the sending process's last `payload_hash`.

**Verifier behavior on missing or malformed handoff event.** When the verifier walks a run that crosses a process boundary (detected by a discontinuity in `process_uuid` or `host` if recorded as observability metadata), it expects the receiving side's first event to be a handoff event. A first event without `audit.handoff.*` fields is structurally indistinguishable from a regular event AS LONG AS its `prev_hash` matches the previous event's `payload_hash` — the per-event MAC will fail otherwise. A handoff event whose `audit.handoff.from_payload_hash` does not equal the entry's `prev_hash` is rejected at spec §7 step 6 (`chain link broken at seq N`).

**Handoff IPC transport (normative for the institution's control description, not for the spec).** The schema above normates the byte content of the handoff event. It does not normate the IPC mechanism that carries `audit.handoff.from_payload_hash` between the sending and receiving processes. In a distributed deployment processes do not share memory; an implementation may use a Redis stream, a Kafka or other message-queue topic, a shared filesystem entry, an OTLP attribute on a coordination event, or any other transport the institution operates. The spec deliberately stays out of the wire choice — institutions arrive with existing IPC stacks, and the chain's integrity claim does not depend on which one carries the bytes.

What the spec DOES require is the institution's CC8.1 control description name the transport and the integrity property it provides. Examples of conformant control-description language:

- "Redis stream `chain.handoff.<tenant_id>`, mutual-TLS to the cluster, append-only entries with consumer-group acknowledgment."
- "Kafka topic `chain-handoff-prod`, mTLS to brokers, topic-level retention 7 days, version-pinned client library."
- "Shared filesystem `/var/lib/chain/handoff/`, per-process subdirectories, filesystem-level access control, append-only flag on the file."
- "OTLP attribute `audit.handoff.from_payload_hash` on a coordination span the orchestrator emits, in-transit TLS 1.3 from orchestrator to worker, attribute pass-through verified per the OTLP-collector pass-through rule (spec §4.4)."

The transport MUST satisfy two normative properties regardless of mechanism:

1. **Origin authentication within the trust boundary.** The receiving process MUST authenticate the handoff event's origin to a process within the same trust boundary as the sending process. Cross-tenant pickup, cross-region pickup, or pickup by a process outside the institution's documented orchestration topology is non-conformant. The authentication mechanism is the transport's standard one (mTLS peer identity for Redis or Kafka, filesystem ACL plus host identity for shared-filesystem, the OTel resource attributes for OTLP-attribute carry).

2. **Byte-exact preservation of `from_payload_hash`.** The transport MUST preserve the `from_payload_hash` value end-to-end as 32 raw bytes. No encoding-layer truncation, no character-set substitution (an UTF-8 encoder treating raw bytes as a string and re-normalizing them is non-conformant), no compression that loses bytes, no length-coercion that drops a leading or trailing byte. If the transport requires a string-typed field, the institution encodes the 32 bytes as hex or base64 and decodes byte-for-byte on the receiving side; the institution's control description names the encoding and the test that confirms round-trip equality.

Once the receiving process has the `from_payload_hash` value via whatever transport, the chain treatment is identical regardless of transport: the receiving process emits its first event with `prev_hash = from_payload_hash`, computes `payload_hash = HMAC-SHA-256(session_key, prev_hash || canonical_bytes)`, and the chain is byte-continuous across the process boundary. The verifier does not see the transport — it sees the chain entries and asserts `entry.prev_hash == previous payload_hash` per spec §7 step 6, which holds whether the bytes traveled by Redis stream, Kafka topic, shared filesystem, or OTLP attribute. The spec normates the bytes; the institution operates the wire.

The shared-run_id pattern is heavier on the implementation; the per-process pattern is simpler and is the default. The conformance corpus includes test vectors for both.

### 9.3 DAG-shaped multi-process flows (third pattern)

Production agent platforms increasingly use DAG-shaped flows where one event has multiple parents — two parallel reasoning paths join at a final-decision step, or an aggregation step combines results from several worker processes. The pattern:

```
ffiec.chain.dag_parents = "r_left:42,r_right:38"
```

The attribute carries a comma-separated list of `(run_id, seq)` pairs identifying the contributing parents. The chain integrity within each contributing run is validated independently; the DAG topology is a graph the verifier reconstructs at examination time.

DAG and parent-child are mutually exclusive on the same event. An event has either `parent_run_id`/`parent_seq` (single parent) or `dag_parents` (multiple parents), not both. The conformance corpus includes a DAG test vector demonstrating the pattern.

## 10. Auditor's-lens review

| Question | Answer |
|---|---|
| Can two implementations of v1.0 produce different `payload_hash` for the same event? | No. JCS canonicalization is deterministic. HMAC-SHA-256 is deterministic. The HKDF inputs are fixed constants plus the tenant_id. Given the same `(IKM, tenant_id, canonical_bytes, prev_hash)` tuple, every conforming implementation produces the same `payload_hash` byte-for-byte. The conformance corpus at `spec/test-vectors/` (seeded from Herald.Py's `tests/fixtures/chain_vectors.json`) enforces this. |
| What if an attacker subverts the SDK's persistence ordering? | Events that have been disclosed (exported, used as prev_hash, read by the caller) are guaranteed persisted. An attacker who corrupts the SQLite cannot retroactively forge a persisted event without also corrupting the daily Merkle seal — which they cannot do without HSM access. |
| What if an operator points a tenant at the wrong IKM (botched rotation, restored backup pointed at the wrong row)? | Caught at the verifier's `key_fingerprint` check (spec §7 step 8) before any MAC compute. The recomputed `SHA-256(utf8(tenant_id) || ikm)[:16]` does not match the entry's recorded `key_fingerprint`; the verifier reports `key_fingerprint mismatch` with a message that names the rotation/restore as the likely cause. No MAC-mismatch storm; the failure is precise. |
| What if an attacker discards the MAC and persists a SHA-256 fingerprint instead (the Herald-discovered headline gap)? | Non-conformant per spec §4.1 inviolate property #6. The verifier's recomputed MAC does not match the persisted `payload_hash` (which would be the SHA-256 fingerprint, not the HMAC), and the day fails. The conformance corpus explicitly tests this case. |
| What if a maintainer relaxes the structural `prev_hash` check in a later version? | The verifier feeds `expected_prev_hash` (the structurally walked value), NOT `entry.prev_hash`, into the MAC recompute (spec §4.1 inviolate property #8). Even if the structural check is relaxed, the MAC compute still uses the value derived from the previous entry's `payload_hash`; an attacker who substitutes `entry.prev_hash` does not also substitute the previous payload_hash, and the MAC mismatch still surfaces. This closes a future-maintainer footgun. |
| What about timing attacks against HMAC? | The session key never leaves process memory. The verifier uses constant-time comparison (`hmac.compare_digest` / `FixedTimeEquals` / `ConstantTimeCompare`) for both the fingerprint check and the MAC check (spec §10.8). |
| What about side-channel attacks against the SDK? | A side-channel attacker on the application host can observe events as they are constructed, but cannot forge past events — those are sealed under a daily Merkle root the attacker cannot influence (server-side seal, HSM-signed). The integrity model bounds the damage to the time between compromise and seal. |
| What about an attacker who flips a chain-stamp field (`key_version`, `key_fingerprint`, `format_version`)? | All three are detected. `key_version` flip → wrong IKM looked up → wrong session key → MAC mismatch. `key_fingerprint` flip → recomputed fingerprint from the looked-up IKM does not match the flipped value → fingerprint mismatch (no MAC compute). `format_version` flip → header/entry mismatch refused. Each is independently caught, with a specific named reason. |
| What about an attacker who lifts a valid event from chain A into chain B's verifier batch? | Detected. The verifier asserts `event.tenant_id == header.tenant_id` AND `event.run_id == header.chain_id` per entry (spec §7 step 4); a lifted event fails one or both before any MAC compute. |
| Recovery after a host restart with corrupted SQLite | v1.0 leaves the recovery procedure to the implementation — institutions document their recovery posture in CC8.1, and `07-verifier-design.md` §8.6 names the verifier-side path for replaying a gap from any source available. A standardized deterministic re-export protocol is a v1.x roadmap candidate; until then, the procedure is institution-defined and the verifier's existing gap-handling carries the audit posture. |

### 10.1 Closed findings from the Herald upgrade audit

The Herald `HMAC-SHA256-HKDF-upgrade.md` design doc was reviewed by an AUDITOR persona (Principal Security Engineer with the Herald codebase loaded fresh). Three BLOCKER, five MAJOR, and six MINOR findings were closed in the design doc and propagated into the Herald.Py reference implementation. The FFIEC v1.0 spec lifts the resolutions verbatim where applicable:

| Herald finding | Resolution adopted by FFIEC v1.0 spec |
|---|---|
| **B1.** Python writer discarded the MAC; .NET design persisted it | Spec §4.1 inviolate property #6 makes "MAC IS the chain entry, persisted byte-for-byte" normative. Non-conformant to compute and discard. |
| **B2.** Length-prefix confusion in `prev_hash || canonical` | Spec §4.1 inviolate property #2 mandates fixed 32 raw bytes for `prev_hash`. |
| **B3.** Constant salt + constant info created multi-tenant cross-contamination risk | Spec §4.1 mandates `info = HKDF_INFO_BASE || "|" || utf8(tenant_id)`; spec §4.1 inviolate property #1. |
| **M1.** `key_version` as a bare integer was insufficient for botched rotation | Spec §4.1 mandates per-entry `key_fingerprint`; spec §7 step 8 mandates fingerprint check before MAC compute. |
| **M2.** Missing audit fields (10-year audit defensibility) | Spec §4.1 per-entry stamp adds `mac_computed_at_utc` and `kms_handle_uri`; spec §3 defines them. |
| **M3.** Memory-zeroisation overpromise | Design 02 §4.3 (above) carries the honest best-effort wording. |
| **M4.** Plaintext-KMS adapter must be compile-time-excluded | Spec §10.7 mandates compile-time exclusion. |
| **M5.** Server-side compromise + KMS Decrypt grant | Spec §4.2/§4.3 (Merkle + HSM rollup) is the composed defence; design 00 §2.3 states the placement reasoning. |

All BLOCKER and MAJOR findings closed. The FFIEC v1.0 spec and this design doc are the formal record.

## 11. Routing-event capture pattern

Spec §4.4.1 makes routing-decision capture normative for institutions operating multi-provider LLM deployments with failover, circuit breakers, or cost-routing. The spec normates the schema of the routing event when emitted (the four event types, the attribute set, the chain-entry shape) but deliberately stays out of the router itself — institutions arrive with existing resilience tooling (Polly in .NET, Resilience4j in JVM, custom-built routers in Python, vendor SDKs that include routing internally) and the chain composes with that tooling rather than replacing it. This section documents the design-time pattern by which the institution wires its existing router to the chain decorator so the four routing event types fire reliably.

### 11.1 The integration shape

The pattern has three components that the institution operates as a single composed unit:

1. **The router emits a routing-event hook.** Every conformant router exposes a hook the institution subscribes to. Polly's `OnRetry` / `OnBreak` / `OnHalfOpen` policies expose direct callbacks; Resilience4j's `CircuitBreaker.EventPublisher` exposes typed event streams; custom-built routers typically expose a delegate, an event-bus publication, or a callback list at the policy-evaluation point. The hook fires with enough information to populate the spec §4.4.1 attribute schema: the providers attempted, the provider chosen, the failover reason (when applicable), the per-provider circuit-breaker states observed, the routing-policy version in force, and the decision timestamp.
2. **The chain decorator subscribes to the hook.** At process startup, alongside the LLM-call decorator and the audit decorator, the institution wires a routing decorator that subscribes to the router's hook. The subscription is a single configuration point in the institution's process-startup wiring; one decorator instance per router instance. The decorator is the bridge: it receives the router's hook payload, translates it into the spec §4.4.1 attribute schema, and produces a chain entry with `audit.routing.event_type` set to the appropriate value.
3. **The decorator emits the chain entry through the same chain pipeline used for ordinary `audit.*` events.** The routing chain entry is structurally a regular chain entry — the canonical bytes cover the routing attributes the same way they cover any other `audit.*` namespace, the per-event `payload_hash` MAC binds the routing decision under the chain's normal integrity property, and the entry persists through the same SQLite + OTLP path as any other event. Routing entries link to the downstream LLM-call entry (when one occurs) via `parent_run_id` / `parent_seq` per spec §4.4. The verifier walks routing entries the same way it walks any other entry; the routing-attribute schema is institution-emitted content and outside the cryptographic verification path.

The pattern keeps the router unchanged. The institution's existing Polly / Resilience4j / custom-built router continues to implement routing logic; the chain decorator is a new subscriber to an existing event surface. Institutions that already emit observability events from the router (typical: a metrics counter for circuit-breaker transitions, a log line per failover) are wiring the chain decorator as one more subscriber alongside those existing emitters.

### 11.2 Implementation constraint — hook fires BEFORE LLM call attempt

The routing-event hook MUST fire BEFORE any LLM call is attempted. This is the load-bearing implementation constraint that makes spec §4.4.1's "absence of a child LLM-call entry is itself evidence" property work. Concretely:

- **`audit.routing.attempt`** fires at the moment the router has selected a provider and is about to invoke. The hook fires BEFORE the HTTP client is given the request. An institution that fires the hook AFTER the call returns (a common naive pattern) produces routing entries only for successful calls — failed calls produce no routing entry, and the chain cannot distinguish "the router never selected this provider" from "the router selected the provider and the call subsequently failed." The before-attempt placement closes that ambiguity.
- **`audit.routing.failover`** fires at the moment the router decides to move away from the previously-attempted provider. The hook fires AFTER the previous attempt has failed but BEFORE the next provider's call is initiated. Same reasoning as above: the chain entry for the failover precedes the next attempt's chain entry.
- **`audit.routing.circuit_state_change`** fires at the moment the router's circuit-breaker transitions for some provider. The hook fires inside the policy's state-transition logic, NOT inside the call path — circuit-state changes are observed regardless of whether a call is in flight, and the routing entry records the transition independently of any specific call.
- **`audit.routing.success`** fires AFTER the LLM call returns successfully. This is the only routing event type that fires after the call; it's the terminating event for a successful routing path and links to the LLM-call entry as its parent.

The before-attempt placement is what makes routing decisions chained even when the call subsequently fails or is precluded. A circuit-breaker that opens before any call is attempted (the policy's threshold has been crossed by previous failures across the institution's fleet) emits a `circuit_state_change` chain entry without any subsequent LLM-call entry — the chain records that the router decided not to attempt a call, and the absence of the LLM-call child is the substantive evidence about the institution's behavior at that moment. A rate-limit kill that prevents a retry emits a `failover` followed by no further chain entry — the chain records that the router gave up, and the absence of further entries is the substantive evidence the institution's MRM committee reviews.

Institutions wiring the pattern document the before-attempt placement in their CC8.1 control description so the SOC team and the FFIEC examiner have a documented anchor for the audit shape P-33 tests against.

### 11.3 Worked-example wiring per common router

The pattern is router-agnostic but the wiring point varies by tooling:

- **Polly (.NET).** The institution registers an `IAsyncPolicy` with `OnRetryAsync`, `OnBreakAsync`, `OnHalfOpenAsync`, and (for fallback policies) `OnFallbackAsync` callbacks. The chain decorator's subscription wires these callbacks to a `RoutingChainEmitter` class that produces the chain entry. The emitter runs synchronously inside the callback so the routing entry's `payload_hash` is computed and persisted before the policy resumes the call path.
- **Resilience4j (JVM).** The institution attaches an `EventPublisher` listener to the relevant `CircuitBreaker`, `Retry`, and `TimeLimiter` instances. The listener implements `onEvent` for each typed event the router emits and dispatches to a `RoutingChainEmitter` bean. Spring-managed wiring places the listener registration in a `@PostConstruct` block on the router-configuration class.
- **Custom-built (Python).** The institution's router exposes a hook list (typical: a `routing_hooks: list[Callable[[RoutingEvent], None]]` field on the router instance). The chain decorator's subscription appends a `routing_chain_hook` callable to the list at process startup. The callable receives a `RoutingEvent` dataclass and produces the chain entry through the SDK's standard audit-emission path.
- **Vendor SDK with internal routing.** Some vendor SDKs (notably some BYOC LLM gateways) implement routing internally without exposing a hook surface. The institution wraps the vendor SDK in a thin proxy that observes the SDK's outbound HTTP calls (or the SDK's own observability events when published) and reconstructs routing decisions from the observed pattern. This wrapping is heavier than the hook pattern — the institution's CC8.1 control description names the wrapping shape, the observability surface the wrapper consumes, and the residual gap if any (typical: circuit-state transitions internal to the SDK that the wrapper cannot observe). For institutions on this path, P-33's failover sub-sample disposition often surfaces gaps the wrapper could not bridge.

The reference implementation in Herald.Py demonstrates the custom-built Python pattern with an explicit `routing_hooks` registry (`src/herald/_runtime.py`); the .NET reference implementation demonstrates the Polly-callback pattern. Other languages and frameworks follow the same shape.

### 11.4 What the verifier sees

From the verifier's perspective, routing chain entries are ordinary chain entries. The verifier walks them per spec §7 — header pre-flight, per-entry MAC, daily Merkle seal, HSM signature — without any routing-specific logic. The `audit.routing.*` attributes are institution-emitted content the verifier does not interpret; the verifier's integrity claim covers the bytes of the routing decision the same way it covers any other event.

The institution's SOC team and the FFIEC examiner consume routing entries through P-33 (`docs/audit-procedures.md` P-33), which is the audit-procedure layer that interprets the routing-attribute schema. The verifier produces the integrity evidence; P-33 produces the control-completeness evidence; the two compose into the institution's full coverage claim per spec §4.4.1.
