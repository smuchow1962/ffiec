# 02 — Chain construction (the hot path)

> **What this doc is.** The detailed design of the HMAC chain construction step. This is the hot path &mdash; runs once per captured event, in the application process, before the event leaves the host. Performance, determinism, and crash-safety are all in scope.

## 1. The contract

For each event, the SDK produces a `payload_hash` defined by:

```
canonical    = JCS(canonical_event_payload)              // RFC 8785
prev_hash    = 0x00...00 (32 bytes)  if seq == 1
             = previous event's payload_hash  otherwise
payload_hash = HMAC-SHA-256(session_key, prev_hash || canonical)
```

`||` is byte concatenation. The output is exactly 32 bytes.

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

The hashed payload includes:

- All semantic event fields (`gen_ai.*`, `tool.*`, `audit.*`)
- The event timestamp
- The `tenant_id`, `run_id`, `seq`
- Excludes: `prev_hash`, `payload_hash` (those are around the payload, not in it)

Implementations include exactly the documented fields. Adding fields silently breaks chain agreement across vendors.

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

At process startup, before any event is captured:

```
1. Process reads the tenant_id from configuration.
2. Process generates a random process_uuid (16 bytes from crypto/rand).
3. Process performs the tenant key handshake (vendor-specific):
   - Out-of-band delivery from a tenant-controlled key broker, OR
   - HSM-token-mediated delivery (a token granting one-shot key derivation), OR
   - SPIFFE/SPIRE-based identity attestation
4. The handshake produces:
   session_key = HKDF-SHA256(
     IKM   = tenant_master,
     salt  = process_uuid || tenant_id,
     info  = "ffiec-ai-chain-v1",
     L     = 32
   )
5. The handshake also produces session_key_id = SHA-256(session_key) (16 leading bytes).
6. The session_key is held in process memory (locked pages where the OS supports it).
7. The session_key_id is exported on every event for verification.
```

The session-key handshake is **not normative in spec v1.0** because vendors must support multiple deployment shapes (cloud, on-prem, vendor-hosted). v1.0 specifies the inputs and outputs of HKDF; the delivery mechanism is implementation-flexible.

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

| Failure | Surface |
|---|---|
| Session key not yet delivered (process started, handshake pending) | SDK returns `ErrSessionKeyNotReady`; events are queued in a startup buffer up to a configurable cap. |
| SQLite full or write fails | SDK returns `ErrPersistenceFailed`; event is lost only if the caller does not retry. The audit-event API always retries until success or a hard timeout. |
| Cache and SQLite disagree | SDK aborts the chain construction and emits an `ErrChainInconsistent` alert. This should never happen in practice; if it does, the host is suspect. |
| `prev_hash` does not match expected | SDK rejects the event. Indicates a bug or a tampering attempt. |
| `seq` skips or reverses | SDK rejects the event. Same indication. |

All failures are loud, named, and actionable. Silent failure is a defect.

## 9. Auditor's-lens review

| Question | Answer |
|---|---|
| Can two implementations of v1.0 produce different `payload_hash` for the same event? | No. JCS canonicalization is deterministic. HMAC-SHA-256 is deterministic. The chain construction is byte-for-byte reproducible. The conformance corpus enforces this. |
| What if an attacker subverts the SDK's persistence ordering? | Events that have been disclosed (exported, used as prev_hash, read by the caller) are guaranteed persisted. An attacker who corrupts the SQLite cannot retroactively forge a persisted event without also corrupting the daily Merkle seal &mdash; which they cannot do without HSM access. |
| What about timing attacks against HMAC? | The session key never leaves process memory. An attacker with sufficient access to mount a timing attack already has access to the key directly. The HMAC subtle-time-equal property is preserved by Go's stdlib (`hmac.Equal`). |
| What about side-channel attacks against the SDK? | A side-channel attacker on the application host can observe events as they are constructed, but cannot forge past events (those are sealed under a daily Merkle root they cannot influence). The integrity model bounds the damage to the time between compromise and seal. |
| Open issue | Recovery after a host restart with corrupted SQLite. v1.0 leaves the recovery procedure to the implementation; v1.1 may standardize a deterministic re-export protocol. |
