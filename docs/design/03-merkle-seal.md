# 03 — Daily Merkle seal

> **What this doc is.** The detailed design of the daily Merkle seal step. Runs in the ledger server, once per tenant per UTC day. Produces the artifact an HSM signs.
>
> **Where this primitive lives.** Server-side, on the bank's perimeter (or vendor-hosted with per-tenant key segregation). The SDK does NOT produce the Merkle seal. The seal is the load-bearing defence against the §2.3 attack class (server-side or privileged-insider history rewrite); see `00-overview.md` §2.3 for the placement reasoning.

## 1. The contract

For each `(tenant_id, UTC_date)` pair, the ledger server produces a single 32-byte `merkle_root` defined by:

```
events_for_day = SELECT payload_hash
                 FROM ledger
                 WHERE tenant_id = T AND DATE(received_at, 'UTC') = D
                 ORDER BY run_id ASC, seq ASC
merkle_root    = MerkleRoot(events_for_day)   // RFC 6962 binary Merkle
```

Any two implementations that produce the same `events_for_day` ordered set MUST produce the same `merkle_root`.

## 2. RFC 6962 Merkle construction

The Certificate Transparency Merkle scheme. Adopted because it is the most widely-trusted Merkle construction in regulator-adjacent systems and because second-preimage attacks are explicitly mitigated.

```go
// Leaf hash — domain-separated by a 0x00 prefix
leafHash := sha256.Sum256(append([]byte{0x00}, payload_hash...))

// Internal hash — domain-separated by a 0x01 prefix
internalHash := sha256.Sum256(append(append([]byte{0x01}, left...), right...))
```

For a tree of N leaves where N is not a power of 2, the right-most subtree contains the unpaired leaves. The construction matches RFC 6962 Section 2.1 exactly.

## 3. Edge cases

### 3.1 Empty days

A tenant-day with zero captured events still gets a seal. The Merkle root for an empty leaf set is defined as `SHA-256("")` per RFC 6962 §2.1:

```
empty_root = SHA-256("")
           = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

The seal is published to maintain continuity. A missing seal is a gap an examiner notices.

**Streaming-Merkle empty-input contract.** The streaming Merkle implementation MUST return the empty-root constant (`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`) when no leaves have been added. A streaming implementation that returns nil / empty bytes / a sentinel value on empty input is non-conformant — it produces a recomputed root that does not match the seal's recorded `merkle_root` for empty days, causing the verifier to FAIL legitimate empty-day seals. The pseudocode in §4.3 below correctly returns the empty-root constant on empty input.

The verifier's per-day Merkle recomputation (per design 07 §4.2) compares `seal.merkle_root` to `computed_root` literally; for empty days both values MUST equal the constant above. Implementations that build the streaming Merkle from a generic library that returns nil on empty input must wrap it with an empty-input check that substitutes the SHA-256("") constant.

### 3.2 Single-event days

A tenant-day with one event has `merkle_root = leaf_hash(event_payload_hash)`. The leaf-prefix domain separation prevents the root from being indistinguishable from a single-event tree-of-trees.

### 3.3 Late-arriving events

Events that arrive after the daily seal has been computed are recorded with a `late_binding` flag and included in the *next* day's Merkle leaves. The original day's seal is not altered.

The verifier reports late-binding events explicitly. A high rate of late-binding events is a signal of operational issues (network, SDK persistence delay) rather than tampering.

### 3.4 Out-of-order events

Events are ordered for Merkle construction by `(run_id, seq)` ascending. **Implementations MUST NOT use receive timestamp for Merkle ordering.** That would make the seal non-deterministic across implementations: two ledger servers receiving the same events in different network orders would produce different roots.

The deterministic ordering rule is the single most important invariant of the seal computation.

### 3.5 Time-stamp authority and the day boundary

The chain hash is independent of the timestamp; `payload_hash` is correct regardless of clock state. But the daily seal aggregates by UTC date, and an event whose `captured_at` falls on the wrong day lands in the wrong seal. Two defenses:

- **Application-host clock SHOULD be NTP-synchronized.** Standard bank IT practice; the chain inherits this control. An institution that does not run NTP cannot rely on the day-boundary semantics.
- **The ledger server's receive clock is authoritative for the day boundary.** When the ledger ingests an event, it records both `captured_at` (from the host) and `received_at` (server-side). The seal job partitions by `received_at` UTC date, not `captured_at`. The `captured_at` field is retained as advisory and as a clock-skew telemetry signal.

The verifier reports clock-skew anomalies (events where `received_at - captured_at` exceeds a configurable threshold, default 5 minutes). High skew is an operational signal, not a tampering signal — the chain catches tampering through HMAC and Merkle, not through timestamp comparison.

### 3.6 Cadence proportionality

Daily is the spec default; the seal cadence is configurable per tenant. Three defensible operating points:

| Cadence | When appropriate | Trade-off |
|---|---|---|
| Hourly | Tier-1 banks with high-volume agent traffic and tight retroactive-tamper-detection windows | Higher HSM operational cost; more seal records to verify |
| Daily | Default; aligns with examination cadence and retention period definitions | Spec default; well-understood |
| Weekly | Community banks with low-volume agent traffic and 18-month examination cycles | Larger time-to-detect for retroactive tampering; document in control description |

Cadence is recorded in the seal record so the verifier can confirm the institution's claimed cadence matches the recorded cadence. Changing cadence is a tenant-policy change documented in the institution's control description.

**Cadence-change examiner-approval workflow.** Relaxing cadence (daily → weekly, weekly → monthly) requires examiner approval. The institution submits a cadence-relaxation request to its primary regulator with:

1. The proposed new cadence and effective date
2. The rationale (typically: AI agent traffic volume, examination cycle, cost considerations)
3. The compensating controls (typically: more aggressive monitoring of integrity-alert events, faster IR response to chain-detected events)
4. The duration of the relaxation (typically until the next examination cycle)

The regulator approves, denies, or requests modifications. The approval is referenced in the institution's control description with the regulator's response document ID. Tightening cadence (weekly → daily) requires only notification, not approval.

A template for the request is in `docs/regulator-pack/examiner-approval-template.md`.

## 3.7 Seal record schema (normative; mirrors spec §4.2)

The signed seal record persisted to the `daily_seals` table carries:

| Field | Type | Source / notes |
|---|---|---|
| `tenant_id` | string | The tenant whose events the seal covers |
| `seal_date` | date (UTC, YYYY-MM-DD) | The tenant-day the seal covers |
| `spec_version` | string | `"v1.0"` for this spec |
| `format_version` | string | `"v1"` for this spec's chain-stamp format. A single seal covers one `format_version`. |
| `merkle_root` | bytes[32] | Apex hash from §2 construction |
| `algorithm` | string | Signature algorithm (`"ed25519"` for v1.0); see §4.3.2 |
| `public_key_id` | string | Resolves to the tenant public key entry that verifies the signature |
| `key_versions` | list of int (≥ 1) | The `key_version` values present in the day's chain entries. Single-element list `[v]` on non-rotation days; multi-element list `[v3, v4]` on a rotation day where some events were captured under the old IKM and others under the new. |
| `hkdf_inputs_digest` | bytes[32] | `SHA-256(HKDF_SALT \|\| info_for_tenant \|\| length_LE32)` for the day's chain construction, where `info_for_tenant = HKDF_INFO_BASE \|\| "\|" \|\| utf8(tenant_id)`. Per-tenant; differs across tenants. Recorded on the seal so a future-version verifier can detect format-drift independent of the chain entries' per-entry stamps. |
| `signature` | bytes[64] | Ed25519 signature over `sign_payload` (§4 below). |
| `signatures` | list of `{algorithm, signature}` | Optional dual-algorithm list per spec §4.2. Each entry covers its own algorithm-bound `sign_payload` (Variant B). When present, the verifier dispatches per-algorithm and produces per-algorithm validation results. Single-algorithm seals OMIT this field. |
| `signed_at` | RFC 3339 UTC | Wall-clock when the HSM signed the root |
| `cadence` | enum | `"hourly"` \| `"daily"` \| `"weekly"`; matches the institution's claimed cadence |
| `late_binding_count` | int64 | Count of events included in the next day's seal due to late arrival |
| `hsm_cluster_member` | string | Optional; advisory pointer to which HSM signed |
| `dev_mode` | bool | Optional; `true` only when the seal was signed by a development software-key adapter (spec §10.7). Verifier under `--strict` refuses `dev_mode=true`. |

**Why `key_versions` is a list.** A tenant-day that crosses a master-key rotation contains events captured under both the old and new IKM. The seal covers the day's full event set (the Merkle root is over all `payload_hash` values regardless of which IKM signed each one); the `key_versions` field records both versions so the verifier resolves the right IKM per entry. The conformance corpus includes a rotation-mid-day case (test vector 010).

**Why `hkdf_inputs_digest` is on the seal.** Chain entries already carry `format_version`, but the seal record needs an independent attestation of the HKDF inputs in force on the day. The digest is computed per-tenant (`info_for_tenant` includes `tenant_id`) so the same `hkdf_inputs_digest` value applies across days for one tenant but differs across tenants. A future-version verifier reading a v1 seal walks this digest first; mismatch refuses the file with a precise reason ("HKDF inputs do not match v1") rather than a downstream MAC-mismatch storm. The digest matches the per-tenant variant defined in spec §3 and used at the file-header pre-flight (spec §7 step 2); both call sites compute the same bytes.

## 4. The seal job

The seal job runs as a periodic task in the ledger server.

### 4.1 Trigger

The job runs:

- At UTC-midnight + a configured delay (default 60 minutes, allowing late-arriving events to settle)
- On manual operator request (for incident response or replay testing)
- On crash recovery, for any unsealed days

### 4.2 Steps

1. Acquire a per-tenant-per-day advisory lock (prevents concurrent seal jobs for the same day).
2. Stream the day's events from append-only ledger storage in `(run_id, seq)` order.
3. Compute the Merkle root incrementally (no need to materialize the whole tree in memory).
4. Build the `sign_payload` per the spec.
5. Submit the signing request to the HSM. Block until the signature returns.
6. Append the `(merkle_root, signature, sign_metadata)` tuple to the seals table.
7. Release the advisory lock.

The job is idempotent: re-running it for an already-sealed day is a no-op (the existing seal is verified to match the recomputed root; if they disagree, an alert fires and operator intervention is required).

### 4.3 Scaling

For tenants with billions of events per day, the Merkle tree cannot be held in memory. The implementation streams events from storage and uses an iterative Merkle construction:

```go
// Pseudocode for streaming Merkle root
type StreamingMerkle struct {
    levels [][]byte // hash at each level, indexed by tree depth
}

func (s *StreamingMerkle) Add(leafHash []byte) {
    h := sha256.Sum256(append([]byte{0x00}, leafHash...))
    level := 0
    for level < len(s.levels) && s.levels[level] != nil {
        h = sha256.Sum256(append(append([]byte{0x01}, s.levels[level]...), h...))
        s.levels[level] = nil
        level++
    }
    if level == len(s.levels) {
        s.levels = append(s.levels, nil)
    }
    s.levels[level] = h[:]
}

func (s *StreamingMerkle) Root() []byte {
    // Empty-day contract per §3.1: when no leaves have been added, return
    // SHA-256("") rather than nil. A nil/empty return on empty input is
    // non-conformant — see §3.1 normative text.
    if len(s.levels) == 0 {
        h := sha256.Sum256([]byte{})
        return h[:]
    }
    var current []byte
    for _, h := range s.levels {
        if h == nil { continue }
        if current == nil {
            current = h
        } else {
            x := sha256.Sum256(append(append([]byte{0x01}, h...), current...))
            current = x[:]
        }
    }
    if current == nil {
        // All level slots were nil (still empty after Add() calls failed
        // to set any level — unreachable in normal operation, defensive only).
        h := sha256.Sum256([]byte{})
        return h[:]
    }
    return current
}
```

Memory footprint: O(log N) regardless of event count.

### 4.4 Failure recovery

If the seal job fails partway:

- After step 4 but before step 5: re-run from step 4 (Merkle root is deterministic; same root will be re-computed).
- After step 5 but before step 6: the HSM signed the root; the ledger does not yet record it. On retry, the HSM re-signs (Ed25519 is deterministic; same signature for same payload). Step 6 inserts the record.
- After step 6: the seal is committed. Idempotency check on retry confirms the recorded seal matches.

Ed25519 determinism (FIPS 186-5) is what makes step 5/6 retry safe. With ECDSA, the nonce randomization would produce a different signature on retry, complicating recovery.

## 5. What this design does not protect

- **Loss of events before they reach the ledger.** If the SDK fails to deliver an event, the ledger doesn&rsquo;t know it existed, so the seal won&rsquo;t cover it. The SDK&rsquo;s local-buffer-with-fsync property is the defense.
- **Modification of the daily seals table itself.** A sufficiently privileged attacker could delete the seal record. This is mitigated by the seal-table being append-only at the application level and by external archival of seals (the institution publishes daily seals to a tenant-controlled archive that the regulator also has read access to).
- **HSM compromise.** If the HSM signing key is extracted (e.g., physical attack on the HSM), the attacker can forge daily roots. FIPS 140-2 Level 3 is designed to make this physically detectable and operationally impractical, but it is not impossible. The threat model in [`09-threat-model.md`](09-threat-model.md) discusses HSM compromise as a high-effort, low-probability adversary.

## 6. Auditor's-lens review

| Question | Answer |
|---|---|
| Why RFC 6962 specifically? | It is the precedent. Certificate Transparency uses it. CONIKS uses it. Trillian implements it. Auditors recognize the construction. |
| Why daily and not hourly? | Daily matches the typical retention-period and examination cadence. Hourly is a configurable option for tenants that want tighter sealing. |
| What if a tenant has no events for several days? | Each empty day still produces a seal. Continuity is the audit property. |
| Is the Merkle root reproducible from the ledger? | Yes. Given the ordered set of `payload_hash` values and the RFC 6962 construction, the root is deterministic. The verifier re-computes it from the ledger and compares to the signed root. |
| What happens during a regulatory inquiry where the examiner needs an old day&rsquo;s seal? | The seal is in append-only storage with the institution&rsquo;s standard retention. The institution provides the ledger snapshot to the examiner; the verifier produces the report. |
| Cross-tenant aggregation | The RFC 6962 construction is per-tenant; a multi-tenant ledger composes per-tenant seals. v1.0 does not specify a meta-Merkle of seals — each tenant's chain stands alone, and cross-tenant comparisons are institution-correlated outside the chain. A meta-Merkle of seals is a v1.x roadmap candidate if multi-tenant verification scenarios warrant it; the v1.0 per-tenant shape is sufficient for single-tenant FFIEC submissions. |
