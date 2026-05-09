# -*- coding: utf-8 -*-
"""Compute the per-leaf audit paths for FFIEC v1 test-vector case
023-merkle-inclusion-proof-rfc6962.

Spec §10.31 normates the per-cohort Merkle subtree disclosure surface:
each disclosed leaf carries an `(leaf_index, leaf_hash, audit_path)`
triple, where `audit_path` is a sequence of `(sibling_hash, is_left)`
pairs proceeding leaf-adjacent → root-adjacent. The verifier folds
the leaf hash through the path via `INTERNAL_HASH(sibling, h)` if
`is_left == true` (sibling on the LEFT) else `INTERNAL_HASH(h, sibling)`,
and asserts the recomputed root equals the seal's signed root.

This case pins the audit paths for every leaf of a 5-leaf tree —
chosen for the non-power-of-2 split discipline (§10.31 splits at
`k = largest power of 2 strictly less than n`). The 5-leaf split
produces a 4-leaf left subtree and a 1-leaf right subtree, exercising:

  - The single-leaf-tree base case (`n=1`, audit path is empty)
  - The recursive split direction for a leaf in the left subtree
  - The recursive split direction for a leaf in the right subtree
  - Mixed-direction audit paths (entries whose path includes both
    `is_left=true` and `is_left=false` steps)

Case 016 pins the **root** for the same 7-leaf right-promote shape;
case 017 pins a single index-3 audit path with placeholder hashes;
case 023 pins **every leaf**'s audit path with real SHA-256 bytes so
a clean-room verifier round-trips inclusion of any leaf without
relying on case 017's structural placeholders.

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


# RFC 6962 §2.1 leaf and internal hash primitives.
def leaf_hash(leaf: bytes) -> bytes:
    return hashlib.sha256(b"\x00" + leaf).digest()


def internal_hash(left: bytes, right: bytes) -> bytes:
    return hashlib.sha256(b"\x01" + left + right).digest()


def merkle_root(leaves_unhashed: list[bytes]) -> bytes:
    """RFC 6962 §2.1 root over `leaves_unhashed` — applies LEAF_HASH."""
    if not leaves_unhashed:
        return hashlib.sha256(b"").digest()
    hashed = [leaf_hash(b) for b in leaves_unhashed]
    return _root_from_hashed(hashed)


def _root_from_hashed(hashed: list[bytes]) -> bytes:
    """Compute root from already-leaf-hashed values using §10.31's
    split-at-largest-power-of-2-strictly-less-than-n recursion.
    """
    n = len(hashed)
    if n == 1:
        return hashed[0]
    k = _largest_pow2_less_than(n)
    left = _root_from_hashed(hashed[:k])
    right = _root_from_hashed(hashed[k:])
    return internal_hash(left, right)


def _largest_pow2_less_than(n: int) -> int:
    """Largest k such that k = 2^j and k < n (n >= 2)."""
    if n < 2:
        raise ValueError("largest_pow2_less_than requires n >= 2")
    k = 1
    while k * 2 < n:
        k *= 2
    return k


def audit_path_from_hashed(
    hashed: list[bytes], m: int
) -> list[dict]:
    """Return the audit path for leaf at index `m` per §10.31.

    Each step is a dict {sibling_hex, is_left}. is_left=true means the
    sibling appears on the LEFT of the current node when combining
    (i.e., the current node is the right child).
    """
    n = len(hashed)
    if not (0 <= m < n):
        raise IndexError(f"leaf index {m} out of range for n={n}")
    if n == 1:
        return []  # base case: leaf IS the root
    k = _largest_pow2_less_than(n)
    if m < k:
        # Recurse into left half; right subtree root is sibling on right.
        inner = audit_path_from_hashed(hashed[:k], m)
        sibling = _root_from_hashed(hashed[k:])
        return inner + [{"sibling_hex": sibling.hex(), "is_left": False}]
    else:
        # Recurse into right half; left subtree root is sibling on left.
        inner = audit_path_from_hashed(hashed[k:], m - k)
        sibling = _root_from_hashed(hashed[:k])
        return inner + [{"sibling_hex": sibling.hex(), "is_left": True}]


def verify_audit_path(
    leaf_h: bytes, path: list[dict], expected_root: bytes
) -> bool:
    """Fold `leaf_h` through `path` and check it equals `expected_root`."""
    h = leaf_h
    for step in path:
        sibling = bytes.fromhex(step["sibling_hex"])
        if step["is_left"]:
            h = internal_hash(sibling, h)
        else:
            h = internal_hash(h, sibling)
    return h == expected_root


# Pinned 5-leaf inputs. Each leaf is a deterministic SHA-256 of a
# distinguishable label so an auditor reading the fixture can tell at
# a glance which leaf is which without computing hashes.
LEAF_LABELS = [
    "leaf-0-tenant-day-2026-05-07-event-1",
    "leaf-1-tenant-day-2026-05-07-event-2",
    "leaf-2-tenant-day-2026-05-07-event-3",
    "leaf-3-tenant-day-2026-05-07-event-4",
    "leaf-4-tenant-day-2026-05-07-event-5",
]
LEAVES = [hashlib.sha256(label.encode("utf-8")).digest() for label in LEAF_LABELS]


def main() -> None:
    n = len(LEAVES)
    hashed = [leaf_hash(b) for b in LEAVES]
    root = _root_from_hashed(hashed)

    # Build a per-leaf disclosure list. Each entry has the leaf index,
    # the leaf's pre-image bytes (so a verifier can recompute the
    # LEAF_HASH from scratch), the leaf hash, the audit path, and a
    # round-trip verification result (which MUST be true).
    disclosures: list[dict] = []
    for m in range(n):
        path = audit_path_from_hashed(hashed, m)
        ok = verify_audit_path(hashed[m], path, root)
        disclosures.append({
            "leaf_index": m,
            "leaf_label": LEAF_LABELS[m],
            "leaf_preimage_hex": LEAVES[m].hex(),
            "leaf_hash_hex": hashed[m].hex(),
            "audit_path": path,
            "audit_path_length": len(path),
            "verifies_to_root": ok,
        })

    fixture = {
        "_about": (
            "Per-leaf audit-path pin for a 5-leaf RFC 6962 tree per "
            "§10.31. The 5-leaf shape exercises the split-at-largest-"
            "power-of-2-strictly-less-than-n recursion (n=5, k=4) — leaf "
            "indices 0..3 sit in the 4-leaf left subtree, leaf index 4 "
            "sits in the 1-leaf right subtree. A clean-room verifier "
            "MUST reproduce each `audit_path` byte-for-byte and round-"
            "trip every leaf to the same root. The Python `verify_audit_"
            "path` in `_compute.py` is the reference fold."
        ),
        "tree": {
            "n_leaves": n,
            "leaves_preimages_hex": [b.hex() for b in LEAVES],
            "leaves_hashed_hex": [b.hex() for b in hashed],
            "merkle_root_hex": root.hex(),
            "split_discipline": (
                "k = largest_power_of_2_less_than(n). For n=5: k=4. Left "
                "subtree = leaves 0..3 (4-leaf, perfect binary). Right "
                "subtree = leaves 4..4 (1-leaf, root IS leaf hash)."
            ),
        },
        "constants": {
            "leaf_hash_recipe": "SHA-256(0x00 || leaf_preimage)",
            "internal_hash_recipe": "SHA-256(0x01 || left_hash || right_hash)",
            "is_left_semantics": (
                "is_left=true means the sibling is on the LEFT of the "
                "current node; combine via INTERNAL_HASH(sibling, current). "
                "is_left=false means the sibling is on the RIGHT; combine "
                "via INTERNAL_HASH(current, sibling)."
            ),
        },
        "disclosures": disclosures,
    }

    # Sanity: every disclosure verifies. If this fails, the generator
    # itself is broken — fail loud rather than emit a bad fixture.
    for d in disclosures:
        assert d["verifies_to_root"], (
            f"audit-path verification failed at leaf_index={d['leaf_index']}"
        )

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    print(f"[023] n_leaves:               {n}")
    print(f"[023] merkle_root:            {root.hex()}")
    for d in disclosures:
        print(
            f"[023] leaf {d['leaf_index']} path_len={d['audit_path_length']} "
            f"verifies={d['verifies_to_root']}"
        )


if __name__ == "__main__":
    main()
