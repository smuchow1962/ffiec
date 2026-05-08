# Case 016 — Non-power-of-2 Merkle (RFC 6962 §2.1 right-promote balancing)

## Purpose

Pin the byte-level expected Merkle roots for trees whose leaf count is not a power of 2. The case forces a verifier to implement the RFC 6962 §2.1 right-promote rule correctly: an unpaired rightmost element at any level is promoted to the next level UNCHANGED, never duplicated to pair with itself, never padded with a zero hash.

The case answers an implementation-level question that has historically been a silent-disagreement source between Merkle implementations: when a tenant-day has 3, 5, 7, 9, ... events, what does a level whose count is odd look like?

## Why this matters

Three common-but-wrong balancing schemes produce three different roots from the same leaves:

1. **Right-duplicate (Bitcoin-style):** the unpaired rightmost element pairs with a copy of itself. `H(L2 || L2)` instead of promoting `L2`. Produces a different root for every odd-count level.
2. **Zero-pad:** the unpaired element pairs with 32 zero bytes. `H(L2 || 0x00..00)` instead of promoting. Produces a different root for every odd-count level.
3. **RFC 6962 right-promote (correct):** the unpaired element passes through to the next level unchanged. The root over 3 leaves combines `h(L0, L1)` (an internal hash) with `L2` (a promoted leaf hash) in the final compression.

A v1 verifier implementing scheme 1 or 2 fails this case. A v1 verifier implementing scheme 3 reproduces the pinned roots byte-for-byte.

## Sub-fixtures

Three trees, increasing in size, each chosen so the right-promote rule fires at multiple levels:

| Size | Where right-promote fires |
|---|---|
| 3 | Top level: `[h(L0, L1), L2]` -> root `h(h(L0,L1), L2)` |
| 5 | Mid level: `[h(L0,L1), h(L2,L3), L4]` -> top: `[h(...,...), L4]` -> root |
| 7 | Mid level: `[h(L0,L1), h(L2,L3), h(L4,L5), L6]` -> top: `[h(...), h(h(L4,L5), L6)]` -> root |

The 7-leaf case is the most informative — `L6` is right-promoted at the second level into a regular pair with `h(L4, L5)`, exercising the "promoted node is then paired" path that the 3-leaf and 5-leaf cases do not reach.

## Inputs

Leaves are SHA-256 outputs of stable seeds, chosen so any conforming implementation can reproduce them without this fixture file:

```
leaf_k = SHA-256(b"ffiec-016-leaf-" || ascii(k))   for k in [0, leaf_count)
```

The hex bytes for each leaf are also pinned in `fixture.json` under `sub_fixtures.{size}_leaves.leaves_hex` for direct copy.

## Expected outputs

| Size | `merkle_root_hex` |
|---|---|
| 3 | `dd6351d5116bd0e6e52ed6731a4dc22de6b1a7811a79fec21a52b1dfcac10810` |
| 5 | `ba019ff1e3c94d606603b6fef0dfd33e1a0ab7202babb9e5e6029b2cdbdaa339` |
| 7 | `ac2e348512498a8d06b1bf432a303d1025ced42a25a7dbe68a8454bd22126779` |

These roots are also pinned in `expected.json`.

## Hashing rules (RFC 6962)

```
leaf_node     = SHA-256(0x00 || leaf_bytes)
internal_node = SHA-256(0x01 || left || right)
empty_root    = SHA-256("")  // = e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

The `0x00` and `0x01` domain-separator bytes prevent leaf-vs-internal collision. Implementations that omit the domain separators produce different roots and are non-conformant.

## Conformance test

A v1 verifier implementing the §4.2 daily Merkle MUST reproduce the three pinned roots byte-for-byte from the published leaves. A verifier that succeeds on `001`, `002`, and `010` (all of which use 1 or 5 leaves with even-paired sub-levels in the 5-leaf case once the right-promote fires) but fails this case has the right-promote rule wrong and would silently disagree with a clean-room implementation on any tenant-day whose leaf count is odd.

## Failure modes

If the 3-leaf root mismatches and reads `e3...` (the empty-tree root), the implementation tried to pad with the empty hash. If it mismatches and reads as `H(0x01 || L0L1 || L2L2)`, the implementation right-duplicated rather than right-promoted. The 5-leaf and 7-leaf cases narrow the failure further by exercising right-promote at non-top levels.

## Provenance

`fixture.json` and `expected.json` are produced by `_compute.py`. Run with `python _compute.py`. The script implements the RFC 6962 §2.1 streaming balance directly; cross-verification by an independent hand-walk produces the same roots (recorded in this case's git history).
