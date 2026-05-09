# -*- coding: utf-8 -*-
"""Compute the §10.37 hierarchical-Merkle top root + per-leaf
audit paths for FFIEC v1 test-vector case 026-hierarchical-merkle-
aggregation.

Spec §10.37 normates a two-level Merkle tree:

  - Level 1 — per-subtree tree. Each subtree is a §10.31 RFC 6962
    tree. Leaves are LEAF_HASH-wrapped.
  - Level 2 — top tree. Combines subtree roots directly via
    INTERNAL_HASH (NO additional LEAF_HASH wrapper at this level —
    subtree roots are already SHA-256 outputs).

A hierarchical inclusion proof is the concatenation of (a) the inner
path within the target subtree (per §10.31 audit-path encoding) and
(b) the outer path between subtree roots (§10.37 no-LEAF_HASH
encoding, otherwise identical to §10.31). Both paths use the same
(sibling_hash, is_left) step shape and combine via INTERNAL_HASH so
§10.31's verify_audit_path verifies the concatenated path unchanged.

This case pins the top root + per-leaf audit paths for a 4-subtree
top tree where each subtree carries different leaf counts (2, 3, 1,
4) so the case exercises:

  - Multiple inner-tree shapes (perfect 2-leaf, non-power-of-2 3-leaf,
    degenerate 1-leaf, perfect 4-leaf)
  - The outer top-tree split-at-largest-power-of-2 recursion (4
    subtrees → k=2, balanced split)
  - Leaves-in-different-subtrees having different audit-path lengths
  - The §10.31-compatible verify path (Python implementation
    reproduces it directly to confirm the concatenation rule)

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


def leaf_hash(leaf: bytes) -> bytes:
    """RFC 6962 §2.1 LEAF_HASH — applied at the inner (per-subtree) level only."""
    return hashlib.sha256(b"\x00" + leaf).digest()


def internal_hash(left: bytes, right: bytes) -> bytes:
    """RFC 6962 §2.1 INTERNAL_HASH — used at both levels."""
    return hashlib.sha256(b"\x01" + left + right).digest()


def _largest_pow2_less_than(n: int) -> int:
    if n < 2:
        raise ValueError("largest_pow2_less_than requires n >= 2")
    k = 1
    while k * 2 < n:
        k *= 2
    return k


def _root_from_hashed(hashed: list[bytes]) -> bytes:
    """Compute root from already-hashed values via split-at-k recursion.

    This function combines already-hashed values via INTERNAL_HASH
    only — no LEAF_HASH wrap. Used for both the inner per-subtree
    tree (whose `hashed` inputs are LEAF_HASH-wrapped leaves) and
    the outer top tree (whose `hashed` inputs are subtree roots,
    already SHA-256 outputs per §10.37).
    """
    n = len(hashed)
    if n == 0:
        return hashlib.sha256(b"").digest()
    if n == 1:
        return hashed[0]
    k = _largest_pow2_less_than(n)
    left = _root_from_hashed(hashed[:k])
    right = _root_from_hashed(hashed[k:])
    return internal_hash(left, right)


def subtree_root_with_leaf_hash(leaves: list[bytes]) -> bytes:
    """§10.31 inner-tree root over LEAF_HASH-wrapped leaves."""
    if not leaves:
        return hashlib.sha256(b"").digest()
    return _root_from_hashed([leaf_hash(b) for b in leaves])


def top_tree_root(subtree_roots: list[bytes]) -> bytes:
    """§10.37 top-tree root — INTERNAL_HASH only, NO LEAF_HASH wrap.

    Subtree roots are already SHA-256 outputs; treating them as raw
    leaves and applying LEAF_HASH would add an unnecessary domain-
    separation layer (per §10.37 normative).
    """
    return _root_from_hashed(subtree_roots)


def _audit_path_inner(hashed: list[bytes], m: int) -> list[dict]:
    """§10.31 audit path within the inner per-subtree tree."""
    n = len(hashed)
    if not (0 <= m < n):
        raise IndexError(f"leaf index {m} out of range for n={n}")
    if n == 1:
        return []
    k = _largest_pow2_less_than(n)
    if m < k:
        inner = _audit_path_inner(hashed[:k], m)
        sibling = _root_from_hashed(hashed[k:])
        return inner + [{"sibling_hex": sibling.hex(), "is_left": False}]
    else:
        inner = _audit_path_inner(hashed[k:], m - k)
        sibling = _root_from_hashed(hashed[:k])
        return inner + [{"sibling_hex": sibling.hex(), "is_left": True}]


def _audit_path_outer(subtree_roots: list[bytes], s: int) -> list[dict]:
    """§10.37 outer audit path — same recursion shape as §10.31's inner
    audit path, but the inputs are subtree roots (no LEAF_HASH) so the
    function is structurally _audit_path_inner over already-hashed values.
    """
    return _audit_path_inner(subtree_roots, s)


def _verify_concatenated_path(
    leaf_h: bytes, path: list[dict], expected_root: bytes
) -> bool:
    h = leaf_h
    for step in path:
        sibling = bytes.fromhex(step["sibling_hex"])
        if step["is_left"]:
            h = internal_hash(sibling, h)
        else:
            h = internal_hash(h, sibling)
    return h == expected_root


# Pinned 4-subtree top tree. Each subtree carries a different leaf
# count to exercise multiple inner-tree shapes; total leaves = 10.
SUBTREES_DEFINITION = [
    # Subtree 0 — 2 leaves (perfect 2-leaf binary tree).
    [
        "subtree-0-leaf-0-event-A",
        "subtree-0-leaf-1-event-B",
    ],
    # Subtree 1 — 3 leaves (non-power-of-2; exercises split-at-2).
    [
        "subtree-1-leaf-0-event-C",
        "subtree-1-leaf-1-event-D",
        "subtree-1-leaf-2-event-E",
    ],
    # Subtree 2 — 1 leaf (degenerate; subtree root IS the leaf hash).
    [
        "subtree-2-leaf-0-event-F",
    ],
    # Subtree 3 — 4 leaves (perfect 4-leaf binary tree).
    [
        "subtree-3-leaf-0-event-G",
        "subtree-3-leaf-1-event-H",
        "subtree-3-leaf-2-event-I",
        "subtree-3-leaf-3-event-J",
    ],
]


def main() -> None:
    # Compute leaf preimages and per-subtree leaf-hash arrays.
    subtrees: list[dict] = []
    subtree_roots: list[bytes] = []
    for s_index, labels in enumerate(SUBTREES_DEFINITION):
        leaves = [hashlib.sha256(label.encode("utf-8")).digest() for label in labels]
        hashed_leaves = [leaf_hash(b) for b in leaves]
        s_root = subtree_root_with_leaf_hash(leaves)
        subtree_roots.append(s_root)
        subtrees.append({
            "subtree_index": s_index,
            "n_leaves": len(leaves),
            "leaf_labels": labels,
            "leaves_preimages_hex": [b.hex() for b in leaves],
            "leaves_hashed_hex": [b.hex() for b in hashed_leaves],
            "subtree_root_hex": s_root.hex(),
        })

    top_root = top_tree_root(subtree_roots)

    # Build the per-leaf hierarchical disclosure list. Each entry is
    # (subtree_index, leaf_index, leaf_hash, audit_path) where the
    # audit path is the CONCATENATION of inner path + outer path.
    disclosures: list[dict] = []
    for s_index, sub in enumerate(subtrees):
        labels = sub["leaf_labels"]
        leaves = [hashlib.sha256(label.encode("utf-8")).digest() for label in labels]
        hashed_leaves = [leaf_hash(b) for b in leaves]
        for m, leaf_label in enumerate(labels):
            inner_path = _audit_path_inner(hashed_leaves, m)
            outer_path = _audit_path_outer(subtree_roots, s_index)
            full_path = inner_path + outer_path
            ok = _verify_concatenated_path(hashed_leaves[m], full_path, top_root)
            disclosures.append({
                "subtree_index": s_index,
                "leaf_index": m,
                "leaf_label": leaf_label,
                "leaf_hash_hex": hashed_leaves[m].hex(),
                "inner_path": inner_path,
                "outer_path": outer_path,
                "concatenated_path": full_path,
                "concatenated_path_length": len(full_path),
                "verifies_to_top_root": ok,
            })

    # Sanity: every disclosure round-trips to the top root.
    for d in disclosures:
        assert d["verifies_to_top_root"], (
            f"hierarchical audit-path verification failed at "
            f"subtree={d['subtree_index']} leaf={d['leaf_index']}"
        )

    fixture = {
        "_about": (
            "Hierarchical Merkle aggregation pin per spec §10.37. A "
            "4-subtree top tree with subtree leaf counts (2, 3, 1, 4) "
            "exercises multiple inner-tree shapes and the outer split-"
            "at-largest-power-of-2 recursion. Per-leaf hierarchical "
            "audit paths are the CONCATENATION of (a) the §10.31 inner "
            "path within the leaf's subtree and (b) the §10.37 outer "
            "path among subtree roots — both using the same "
            "(sibling_hex, is_left) step shape, so a §10.31 verifier "
            "verifies the concatenated path unchanged."
        ),
        "constants": {
            "leaf_hash_recipe": "SHA-256(0x00 || leaf_preimage)  (inner level only)",
            "internal_hash_recipe": "SHA-256(0x01 || left || right)  (both levels)",
            "top_level_no_leaf_hash_note": (
                "The outer (top-level) tree combines subtree roots "
                "DIRECTLY via INTERNAL_HASH — no LEAF_HASH wrap. "
                "Subtree roots are already SHA-256 outputs; treating "
                "them as raw leaves and applying LEAF_HASH would add "
                "an unnecessary domain-separation layer."
            ),
            "split_discipline": (
                "k = largest_power_of_2_less_than(n). For n=4 (top-"
                "level subtree count): k=2 (balanced split). For n=3 "
                "(subtree 1's leaf count): k=2."
            ),
        },
        "subtrees": subtrees,
        "top_root_hex": top_root.hex(),
        "disclosures": disclosures,
    }

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    print(f"[026] subtrees:               {len(subtrees)}")
    for s in subtrees:
        print(
            f"[026]   subtree {s['subtree_index']} "
            f"n_leaves={s['n_leaves']:1d} root={s['subtree_root_hex'][:16]}..."
        )
    print(f"[026] top_root:               {top_root.hex()}")
    print(f"[026] disclosures:            {len(disclosures)}")
    print(f"[026] all verify to top root: {all(d['verifies_to_top_root'] for d in disclosures)}")


if __name__ == "__main__":
    main()
