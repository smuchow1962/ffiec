# -*- coding: utf-8 -*-
"""
Compute byte-level fixtures for FFIEC v1 test-vector case 016-non-power-of-2-merkle.

Three sub-fixtures over distinguishable `payload_hash` leaf values:
  - 3 leaves
  - 5 leaves
  - 7 leaves

RFC 6962 §2.1 right-promote balancing: when a level has an odd number of
elements, the unpaired rightmost element is PROMOTED to the next level
unchanged (NOT duplicated to pair with itself, NOT padded with a zero hash).

Tree shapes:

  3-leaf:                             root
                                      /  \\
                              h(L0,L1)    L2     <- L2 right-promoted
                              /     \\
                             L0     L1

  5-leaf:                                  root
                                          /    \\
                                  h(A,B)        L4   <- L4 right-promoted
                                 /      \\
                          h(L0,L1)     h(L2,L3)
                            /  \\         /  \\
                           L0  L1       L2   L3

  7-leaf:                                 root
                                         /     \\
                                  h(A,B)        h(C,L6)   <- L6 right-promoted
                                  /     \\          /  \\
                          h(L0,L1)  h(L2,L3)  h(L4,L5) L6
                            /  \\      /  \\     /  \\
                           L0  L1    L2  L3   L4  L5

Each leaf is hashed as `H(0x00 || leaf)` per RFC 6962; each internal node is
`H(0x01 || left || right)`. The compute script implements the streaming
balance directly; rerunning produces byte-identical roots.

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import json
import os


HERE = os.path.dirname(os.path.abspath(__file__))


def leaf_hash(leaf: bytes) -> bytes:
    return hashlib.sha256(b"\x00" + leaf).digest()


def internal_hash(left: bytes, right: bytes) -> bytes:
    return hashlib.sha256(b"\x01" + left + right).digest()


def merkle_root_rfc6962(payload_hashes: list[bytes]) -> bytes:
    """RFC 6962 §2.1 streaming Merkle with right-promote balancing.

    Cognitive-complexity note: the inner loop walks the level pairwise; an
    unpaired rightmost element promotes to the next level unchanged. No
    duplication, no zero-padding — that is the load-bearing property the case
    exercises.
    """
    if not payload_hashes:
        return hashlib.sha256(b"").digest()

    level = [leaf_hash(h) for h in payload_hashes]
    while len(level) > 1:
        next_level: list[bytes] = []
        i = 0
        while i + 1 < len(level):
            next_level.append(internal_hash(level[i], level[i + 1]))
            i += 2
        if i < len(level):
            next_level.append(level[i])
        level = next_level
    return level[0]


def _make_distinguishable_leaves(count: int) -> list[bytes]:
    """Pin distinguishable byte-level leaves so a parser bug surfaces clearly.

    Leaf k = SHA-256(b"ffiec-016-leaf-" || str(k).encode("ascii")).
    Choosing SHA-256 outputs (32 bytes each) keeps the leaves the right size
    for `payload_hash` semantics and makes the fixture trivial to recompute.
    """
    leaves: list[bytes] = []
    for k in range(count):
        seed = b"ffiec-016-leaf-" + str(k).encode("ascii")
        leaves.append(hashlib.sha256(seed).digest())
    return leaves


def main() -> None:
    sub_fixtures: dict[str, dict] = {}

    for size in (3, 5, 7):
        leaves = _make_distinguishable_leaves(size)
        root = merkle_root_rfc6962(leaves)
        sub_fixtures[f"{size}_leaves"] = {
            "leaf_count": size,
            "leaves_hex": [b.hex() for b in leaves],
            "merkle_root_hex": root.hex(),
        }

    fixture = {
        "_about": (
            "FFIEC v1 test-vector case 016-non-power-of-2-merkle. RFC 6962 "
            "§2.1 right-promote balancing for odd-leaf trees. Three sub-"
            "fixtures: 3, 5, and 7 leaves. Leaves are SHA-256 outputs of "
            "stable seeds (see _compute.py) so reproducing them does not "
            "require this script — any conforming implementation can derive "
            "the leaves from the seeds and compute the root."
        ),
        "leaf_seed_format": "leaf k = SHA-256(b\"ffiec-016-leaf-\" || ascii(k)) for k in [0, leaf_count)",
        "rfc6962_balancing": (
            "When a level has an odd count, the rightmost unpaired node is "
            "promoted unchanged to the next level (NOT duplicated, NOT zero-"
            "padded). Each leaf is hashed as H(0x00 || leaf); each internal "
            "node is H(0x01 || left || right)."
        ),
        "sub_fixtures": sub_fixtures,
    }

    with open(os.path.join(HERE, "fixture.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    expected = {
        "_about": (
            "Expected Merkle roots for case 016. A conforming v1 verifier "
            "implementing RFC 6962 §2.1 right-promote balancing MUST "
            "reproduce these byte-for-byte from the leaves in fixture.json."
        ),
        "merkle_roots": {
            f"{size}_leaves": sub_fixtures[f"{size}_leaves"]["merkle_root_hex"]
            for size in (3, 5, 7)
        },
    }

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(expected, f, indent=2, ensure_ascii=False)
        f.write("\n")

    for size in (3, 5, 7):
        print(f"[016] {size}-leaf merkle_root: {sub_fixtures[f'{size}_leaves']['merkle_root_hex']}")


if __name__ == "__main__":
    main()
