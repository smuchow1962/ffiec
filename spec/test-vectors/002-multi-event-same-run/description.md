# Case 002 — Multi-event, same run

## Purpose

Verifies the chain links across multiple events: every `prev_hash` (after `seq=1`) equals the previous entry's `payload_hash`, and the MAC compute uses the structural chain value not a stamped one.

## Inputs

From `../chain_vectors.json` `single_chain` array — all five events under `key_version=1` with `ikm_v1`. Same `tenant_id`, same `run_id`, sequential `seq` 1..5.

## Expected outputs

| `seq` | `prev_hash_hex` | `payload_hash_hex` |
|---|---|---|
| 1 | `0000...0000` | `e84168071562a9e866f28977a204b43abfdc959b7f84aad83c16fe790f63b049` |
| 2 | `e84168071562a9e866f28977a204b43abfdc959b7f84aad83c16fe790f63b049` | `d9acb553eb1fc59911ebf04e78f9b97fd032c68b8b50bb9073ea4bff1a5de89a` |
| 3 | `d9acb553eb1fc59911ebf04e78f9b97fd032c68b8b50bb9073ea4bff1a5de89a` | `c6f45090813d8d81a8425e07dcc33318133539ddea9e53b0437060ad4ed732c6` |
| 4 | `c6f45090813d8d81a8425e07dcc33318133539ddea9e53b0437060ad4ed732c6` | `57f11c41a5b2194cb0dad5d5dbd56b09d9bc71f7ea09522c2b6f0815bbf6aa9c` |
| 5 | `57f11c41a5b2194cb0dad5d5dbd56b09d9bc71f7ea09522c2b6f0815bbf6aa9c` | `bb0c0deee87b7ed56d48a3bd306728ff45c7e44e767aef0045ef38c15ecec257` |

Merkle root over the five `payload_hash` values: `927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9`.

## What this case proves

- `prev_hash` of entry N+1 equals `payload_hash` of entry N (the chain link).
- The session key persists across the run; same key derives the same MAC.
- RFC 6962 streaming Merkle produces the expected root from the ordered `payload_hash` list.

## Failure mode

If the chain breaks at some `seq=K`, suspect:

1. `prev_hash` for `seq=K` was computed against the wrong previous entry (typically the canonical bytes vs the MAC).
2. The session key changed mid-run (it should not — the cache is per-tenant within a process).
3. The canonical-form excluded different fields between entries (must be deterministic per-entry).
