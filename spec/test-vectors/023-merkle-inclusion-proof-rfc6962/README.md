# Case 023 — Per-leaf Merkle inclusion proof (RFC 6962 / §10.31)

## What this case verifies

Spec §10.31 normates the per-cohort Merkle subtree disclosure surface for institutions serving multiple regulator-jurisdictions or running multi-tenant SaaS. Each disclosed leaf carries an `(leaf_index, leaf_hash, audit_path)` triple where `audit_path` is a sequence of `(sibling_hash, is_left)` pairs proceeding from leaf-adjacent outward to root-adjacent. The verifier folds the leaf hash through the path and asserts the recomputed root equals the seal's signed root.

This case pins the audit paths for **every leaf** of a 5-leaf RFC 6962 tree with real SHA-256 bytes. The 5-leaf shape exercises the §10.31 split-at-largest-power-of-2-strictly-less-than-n recursion (`n=5`, `k=4`) — leaf indices 0..3 sit in the 4-leaf left subtree, leaf index 4 sits in the 1-leaf right subtree.

A clean-room verifier MUST:

- Reproduce each `audit_path` byte-for-byte from the same input leaves.
- Round-trip every leaf (recompute the root by folding the leaf hash through the audit path) and obtain the same `merkle_root_hex`.

## Relationship to cases 016 / 017

- **Case 016** (`016-non-power-of-2-merkle/`) pins the Merkle **root** for non-power-of-2 trees (3, 5, 7 leaves) using a different right-promote shape (older spec).
- **Case 017** (`017-merkle-inclusion-partial-disclosure/`) pins a single index-3 audit path with placeholder hashes (structural fixture).
- **Case 023** pins **every leaf**'s audit path with real SHA-256 bytes under the §10.31 split-at-k recursion. A clean-room verifier passing 023 has byte-real coverage for inclusion proofs across all leaf positions.

The §10.31 split-at-largest-power-of-2 discipline differs from case 016's right-promote balancing — the spec moved to the split discipline for §10.31 because audit-path encoding is cleaner under split-at-k. Case 023 reflects the current §10.31 normative form.

## Inputs

Five leaves, each a deterministic `SHA-256(utf8(label))` so an auditor reading the fixture can identify each leaf at a glance:

| Index | Label | Leaf preimage (hex prefix) |
|---|---|---|
| 0 | `leaf-0-tenant-day-2026-05-07-event-1` | `5fa4...` |
| 1 | `leaf-1-tenant-day-2026-05-07-event-2` | (computed) |
| 2 | `leaf-2-tenant-day-2026-05-07-event-3` | (computed) |
| 3 | `leaf-3-tenant-day-2026-05-07-event-4` | (computed) |
| 4 | `leaf-4-tenant-day-2026-05-07-event-5` | (computed) |

Full preimage and leaf-hash hex values are in `expected.json` under `tree.leaves_preimages_hex` / `tree.leaves_hashed_hex`.

## Expected outputs

| Property | Value |
|---|---|
| n_leaves | `5` |
| merkle_root_hex | `12182c226015d4b8e19ac69184bf8563c21b5413fd3b2553c37b5b5a7a04cbe0` |
| audit_path length, leaves 0..3 | `3` (2 inner steps within the 4-leaf left subtree + 1 outer step combining with leaf 4) |
| audit_path length, leaf 4 | `1` (just the 4-leaf left-subtree root as a left-sibling) |

`expected.json` carries the per-leaf disclosure list. Each entry has the leaf index, label, preimage, leaf hash, the audit path as `(sibling_hex, is_left)` pairs, and a `verifies_to_root` boolean confirming the round-trip.

## RFC 6962 / §10.31 primitives

```
LEAF_HASH(leaf)  = SHA-256(0x00 || leaf)
INTERNAL_HASH(L, R) = SHA-256(0x01 || L || R)
```

`is_left` semantics:
- `is_left = true` — the sibling is on the **LEFT** of the current node. Combine via `INTERNAL_HASH(sibling, current)`.
- `is_left = false` — the sibling is on the **RIGHT**. Combine via `INTERNAL_HASH(current, sibling)`.

## Conformance behavior

A conforming implementation:

1. Computes the Merkle root over the 5 leaves per RFC 6962 §2.1 / §10.31 split-at-largest-power-of-2 and reproduces `merkle_root_hex` byte-for-byte.
2. For each leaf index `m`, produces an `audit_path` matching the published sequence of `(sibling_hex, is_left)` pairs.
3. Folds each published audit path back to the root (per `verify_audit_path` in `_compute.py`) and asserts byte-equality with `merkle_root_hex`.
4. Refuses an audit-path lookup against an empty tree per §10.31 (the empty-day Merkle root is `SHA-256(b"")` = `e3b0c44…b855`; producers MUST refuse audit paths against it).

## Negative cases this fixture supports

- An implementation that confuses `is_left` semantics (treats `is_left=true` as "sibling on the right") fails the round-trip on at least one leaf.
- An implementation that recurses into the right subtree first (rather than splitting at `k = largest power of 2 strictly less than n`) produces different audit paths for non-power-of-2 trees.
- An implementation that wraps subtree roots with an extra `LEAF_HASH` layer when combining at internal nodes (a common bug when porting from a hierarchical-tree implementation per §10.37) fails the recompute.

## Cross-references

- Spec §10.31 per-cohort Merkle subtree disclosure (the section this case pins)
- Spec §4.2 daily seal Merkle root (the flat-tree default §10.31 discloses subtrees of)
- RFC 6962 §2.1 / §2.1.1 (the leaf-hash, internal-hash, and audit-path constructions)
- Case 016 (`016-non-power-of-2-merkle/`) — non-power-of-2 root pin (older right-promote shape)
- Case 017 (`017-merkle-inclusion-partial-disclosure/`) — single-leaf audit path with placeholder hashes
- Case 026 (`026-hierarchical-merkle-aggregation/`) — §10.37 two-level tree using §10.31's audit-path encoding for the inner level

## Reproduction

The `_compute.py` script in this directory regenerates `expected.json` from the pinned 5 leaf labels. It uses only the Python standard library and includes a self-check that every audit path round-trips to the root before writing the file. Re-run after any spec change that affects §10.31 audit-path encoding:

```
python _compute.py
```

Cross-validation: the 5-leaf root reproduced by Herald.Py's `herald._merkle_disclosure.compute_merkle_root` from the same leaf bytes is byte-identical to the value above (`12182c…cbe0`).
