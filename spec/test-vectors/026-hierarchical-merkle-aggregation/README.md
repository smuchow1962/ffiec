# Case 026 — Hierarchical Merkle aggregation (§10.37)

## What this case verifies

Spec §10.37 normates a two-level Merkle tree for institutions whose flat-tree audit-path bandwidth dominates per-leaf overhead (Saraswati Microfinance: 15,000 tablets × ~1,000 events/day):

- **Level 1 — per-subtree tree.** Each subtree (typically per device under §10.32, but the boundary is institution-defined) is its own §10.31 RFC 6962 tree with `LEAF_HASH`-wrapped leaves.
- **Level 2 — top tree.** The top-level tree combines subtree roots **directly** via `INTERNAL_HASH` — no `LEAF_HASH` wrap at this level because subtree roots are already SHA-256 outputs.

A hierarchical inclusion proof is the **concatenation** of (a) the §10.31 inner path within the target subtree and (b) the §10.37 outer path between subtree roots. Both paths use the same `(sibling_hash, is_left)` step shape and combine via `INTERNAL_HASH`, so a §10.31 verifier verifies the concatenated path unchanged.

This case pins the top root + per-leaf audit paths for a 4-subtree top tree where each subtree carries different leaf counts (`2, 3, 1, 4`) so the case exercises:

- Multiple inner-tree shapes (perfect 2-leaf, non-power-of-2 3-leaf, degenerate 1-leaf, perfect 4-leaf)
- The outer top-tree split-at-largest-power-of-2 recursion (4 subtrees → `k=2`, balanced split)
- Leaves in different subtrees having different audit-path lengths
- The §10.31-compatible verify path (the Python implementation reproduces it directly to confirm the concatenation rule)

## Inputs

Four subtrees with deterministic leaf labels:

| Subtree | Leaf labels | Leaf count |
|---|---|---|
| 0 | `subtree-0-leaf-0-event-A`, `subtree-0-leaf-1-event-B` | 2 |
| 1 | `subtree-1-leaf-0-event-C`, `subtree-1-leaf-1-event-D`, `subtree-1-leaf-2-event-E` | 3 |
| 2 | `subtree-2-leaf-0-event-F` | 1 |
| 3 | `subtree-3-leaf-0-event-G`, `subtree-3-leaf-1-event-H`, `subtree-3-leaf-2-event-I`, `subtree-3-leaf-3-event-J` | 4 |

Each leaf preimage is `SHA-256(utf8(label))` so an auditor reading the fixture can identify each leaf at a glance and recompute leaf bytes from first principles.

## Tree shape

```
                      top_root
                          │
              INTERNAL_HASH(L_outer, R_outer)
              ┌───────────┴────────────┐
       L_outer (n=2 left)         R_outer (n=2 right)
       INTERNAL_HASH(s0, s1)      INTERNAL_HASH(s2, s3)
       ┌──────┴──────┐            ┌──────┴──────┐
       s0 (2 leaves) s1 (3 leaves) s2 (1 leaf)  s3 (4 leaves)
       §10.31 root   §10.31 root    leaf_hash    §10.31 root
       perfect-2     split-at-2     IS the root  perfect-4
```

The `s2` (1-leaf subtree) is degenerate — its root is the leaf hash with no further wrapping. The `s1` (3-leaf subtree) exercises the non-power-of-2 split discipline (`k=2`, left=2 leaves, right=1 leaf).

## Expected outputs

| Property | Value |
|---|---|
| number of subtrees | `4` |
| total leaves | `10` |
| subtree 0 root (hex prefix) | `e6042e6e409d23d3…` |
| subtree 1 root (hex prefix) | `fa0e648386c9fdb7…` |
| subtree 2 root (hex prefix) | `2c680f50f0e1ea06…` (degenerate — equals `LEAF_HASH("subtree-2-leaf-0-event-F")`) |
| subtree 3 root (hex prefix) | `83a56494859c3b8c…` |
| **top_root_hex** | `53451e5eff8083bc58fe0975d3582af3e64c4c9a3a445a181dee122e8c262156` |
| disclosures | 10 (one per leaf, each verifying to `top_root_hex`) |

`expected.json` carries the full per-leaf disclosure list. Each disclosure has the leaf identification, the leaf hash, the inner audit path within the subtree, the outer audit path between subtree roots, the concatenated path, the path length, and a `verifies_to_top_root` boolean confirming the round-trip.

## §10.37 conformance behavior

A conforming implementation supporting §10.37:

1. Computes each subtree's root per §10.31 RFC 6962 over `LEAF_HASH`-wrapped leaves.
2. Computes the top root by combining subtree roots **directly** via `INTERNAL_HASH` — NO `LEAF_HASH` wrap on the subtree roots. Reproduces `top_root_hex` byte-for-byte.
3. For each leaf, produces an audit path that is the concatenation of `inner_path` (within the leaf's subtree) and `outer_path` (between subtree roots), in that order.
4. Folds each disclosure's `concatenated_path` from the leaf hash via `INTERNAL_HASH(sibling, h)` if `is_left=true` else `INTERNAL_HASH(h, sibling)` and asserts byte-equality with `top_root_hex`.
5. Uses the §10.31 `verify_audit_path` algorithm unchanged — the §10.37 outer path encoding is structurally identical to the §10.31 audit-path encoding.

## Bandwidth observation (informative)

For Saraswati's 15,000-device × 1,000-events workload (15,000,000 leaves):

- Flat path: ~24 hashes per leaf × 1,000,000 leaves = ~24,000,000 hashes per device disclosure
- Hierarchical path: 1,000 × 10 (inner) + 14 (shared outer) = ~10,015 hashes per device disclosure
- ~58% bandwidth savings

The per-device disclosure shares the outer 14-deep path across all 1,000 leaves of that device — a saving the flat tree cannot realise.

## Negative cases this fixture supports

- An implementation that wraps subtree roots with `LEAF_HASH` at the top level (a common bug when porting from a flat-tree implementation) produces a different top root and fails the byte pin.
- An implementation that emits the audit path as `outer_path + inner_path` (wrong concatenation order) fails the round-trip on every leaf.
- An implementation that swaps `is_left` semantics on the outer path (treats the outer level differently from the inner level) fails leaves whose outer path includes both directions.

## Cross-references

- Spec §10.37 hierarchical Merkle aggregation (the section this case pins)
- Spec §10.31 per-cohort Merkle subtree disclosure (the inner-tree audit-path encoding §10.37 reuses)
- Spec §4.2 daily seal Merkle root (the flat-tree default §10.37 is an alternative to)
- RFC 6962 §2.1 / §2.1.1 (leaf-hash, internal-hash, audit-path constructions)
- Case 023 (`023-merkle-inclusion-proof-rfc6962/`) — §10.31 audit-path round-trip for a 5-leaf flat tree (the inner-tree encoding case 026 reuses)
- Case 024 (`024-per-device-derivation/`) — §10.32 the natural subtree boundary for §10.37

## Reproduction

The `_compute.py` script in this directory regenerates `expected.json` from the pinned subtree definitions. It uses only the Python standard library and includes a self-check that every leaf's concatenated audit path round-trips to the top root before writing the file.

```
python _compute.py
```

Cross-validation: the §10.31 inner-tree primitive is identical to case 023's `_compute.py` and Herald.Py's `herald._merkle_disclosure.compute_merkle_root` — both byte-validated above.
