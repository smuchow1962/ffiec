# -*- coding: utf-8 -*-
"""
Compute byte-level fixture for FFIEC v1 test-vector case 003-multi-run-same-day.

Two runs in the same tenant-day. Run A has 3 events under key_version=1; run B
has 3 events under key_version=1. Run B's first event resets `prev_hash` to 32
zero bytes (cross-run isolation — chain links never cross run boundaries).

The daily Merkle tree aggregates all 6 `payload_hash` values in `(run_id, seq)`
ordering per spec §4.2 (lexicographic by run_id ASC, then seq ASC).

Cognitive-complexity note: each helper is a pure function; the script reads
top-to-bottom as input -> derive -> chain -> Merkle -> sign_payload.

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import hmac
import json
import os
from typing import Any

HERE = os.path.dirname(os.path.abspath(__file__))

# Constants from spec §4.1 (FFIEC-conformance posture)
HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"
HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"
HKDF_LENGTH = 32
LENGTH_LE32 = (32).to_bytes(4, "little")  # 0x20 0x00 0x00 0x00

# Test inputs (mirrors central corpus)
TENANT_ID = "tenant-ffiec-test-1"
RUN_A = "RUN-FFIEC-VECTOR-A"
RUN_B = "RUN-FFIEC-VECTOR-B"
SEAL_DATE = "2026-05-06"
ALGORITHM = "ed25519"
FORMAT_VERSION = "v1"
CADENCE = "daily"
DEV_MODE_BYTE = "0"  # dev_mode=false -> single ASCII byte "0"

# Same ikm_v1 as central corpus
IKM_V1_HEX = (
    "66666965632d63726f73732d6c616e67756167652d746573742d766563"
    "746f722d696b6d2d76312d3332627974657321"
)
IKM_V1 = bytes.fromhex(IKM_V1_HEX)


def _hkdf_extract(salt: bytes, ikm: bytes) -> bytes:
    return hmac.new(salt, ikm, hashlib.sha256).digest()


def _hkdf_expand(prk: bytes, info: bytes, length: int) -> bytes:
    """RFC 5869 expand. length=32 fits in the first round; T(1) IS the output."""
    okm = b""
    t = b""
    counter = 1
    while len(okm) < length:
        t = hmac.new(prk, t + info + bytes([counter]), hashlib.sha256).digest()
        okm += t
        counter += 1
    return okm[:length]


def hkdf_sha256(ikm: bytes, salt: bytes, info: bytes, length: int) -> bytes:
    return _hkdf_expand(_hkdf_extract(salt, ikm), info, length)


def info_for_tenant(tenant_id: str) -> bytes:
    return HKDF_INFO_BASE + b"|" + tenant_id.encode("utf-8")


def session_key(ikm: bytes, tenant_id: str) -> bytes:
    return hkdf_sha256(ikm, HKDF_SALT, info_for_tenant(tenant_id), HKDF_LENGTH)


def key_fingerprint(tenant_id: str, ikm: bytes) -> bytes:
    return hashlib.sha256(tenant_id.encode("utf-8") + ikm).digest()[:16]


def hkdf_inputs_digest(tenant_id: str) -> bytes:
    return hashlib.sha256(
        HKDF_SALT + info_for_tenant(tenant_id) + LENGTH_LE32
    ).digest()


def jcs_canonicalize(value: Any) -> bytes:
    """Minimal JCS canonicalisation for the simple shapes this fixture uses.

    The corpus's per-event canonical bytes for cases 001/002/010 use sorted
    object keys, no whitespace, UTF-8 string output, and no float values
    (timestamps are ints, span_id/trace_id are base64 strings). We replicate
    that here directly with `json.dumps(..., sort_keys=True, separators=(",", ":"))`.
    """
    return json.dumps(
        value, sort_keys=True, separators=(",", ":"), ensure_ascii=False
    ).encode("utf-8")


def event_canonical(
    *,
    tenant_id: str,
    run_id: str,
    seq: int,
    name: str,
    timestamp_ns: int,
    span_id_b64: str,
    trace_id_b64: str,
    event_id: str,
    data: str,
    step: int,
) -> bytes:
    obj = {
        "attributes": {"data": data, "step": step},
        "chain_kind": "audit",
        "duration_ns": None,
        "event_id": event_id,
        "kind": "audit",
        "name": name,
        "parent_span_id": None,
        "resource": {},
        "run_id": run_id,
        "severity": "AUDIT",
        "span_id": span_id_b64,
        "tenant_id": tenant_id,
        "timestamp_ns": timestamp_ns,
        "trace_id": trace_id_b64,
    }
    return jcs_canonicalize(obj)


def payload_hash(session_key_bytes: bytes, prev_hash: bytes, canonical: bytes) -> bytes:
    return hmac.new(session_key_bytes, prev_hash + canonical, hashlib.sha256).digest()


# RFC 6962 Merkle (right-promote balancing)
def _rfc6962_leaf_hash(leaf: bytes) -> bytes:
    return hashlib.sha256(b"\x00" + leaf).digest()


def _rfc6962_internal_hash(left: bytes, right: bytes) -> bytes:
    return hashlib.sha256(b"\x01" + left + right).digest()


def merkle_root_rfc6962(payload_hashes: list[bytes]) -> bytes:
    """RFC 6962 §2.1 streaming Merkle. Right-promote unpaired leaf at each level.

    Note: spec §4.2 specifies the Merkle leaves are the per-event `payload_hash`
    values directly. We hash each as a leaf per RFC 6962 (`H(0x00 || leaf)`).
    """
    if not payload_hashes:
        return hashlib.sha256(b"").digest()

    level = [_rfc6962_leaf_hash(h) for h in payload_hashes]
    while len(level) > 1:
        next_level: list[bytes] = []
        i = 0
        while i + 1 < len(level):
            next_level.append(_rfc6962_internal_hash(level[i], level[i + 1]))
            i += 2
        if i < len(level):
            # Right-promote the unpaired rightmost leaf without re-hashing.
            next_level.append(level[i])
        level = next_level
    return level[0]


def build_sign_payload(
    *,
    algorithm: str,
    format_version: str,
    tenant_id: str,
    seal_date: str,
    merkle_root_hex: str,
    hkdf_inputs_digest_hex: str,
    cadence: str,
    dev_mode_byte: str,
) -> bytes:
    text = (
        "ffiec.chain-of-custody.v1\n"
        + algorithm + "\n"
        + format_version + "\n"
        + tenant_id + "\n"
        + seal_date + "\n"
        + merkle_root_hex + "\n"
        + hkdf_inputs_digest_hex + "\n"
        + cadence + "\n"
        + dev_mode_byte
    )
    return text.encode("utf-8")


def _build_run_events(run_id: str, base_ns: int, base_step: int) -> list[dict]:
    """Build three events for one run, each with distinguishable canonical bytes.

    The span_id / trace_id / event_id schema mirrors the central corpus shape so
    a verifier can compare all the byte-level fields the same way.
    """
    suffix = "A" if run_id.endswith("A") else "B"
    events: list[dict] = []
    for offset in range(3):
        seq = offset + 1
        step = base_step + offset
        # Distinguishable byte payloads per (run, seq).
        # span_id is 8 bytes; trace_id is 16 bytes. Patterns are run-stamped.
        run_byte = ord(suffix)              # 'A' = 0x41, 'B' = 0x42
        seq_byte = seq                       # 1, 2, 3
        span_id_bytes = bytes([run_byte, seq_byte] * 4)
        trace_id_bytes = bytes([run_byte, seq_byte] * 8)
        import base64
        span_id_b64 = base64.b64encode(span_id_bytes).decode("ascii")
        trace_id_b64 = base64.b64encode(trace_id_bytes).decode("ascii")
        events.append({
            "seq": seq,
            "name": f"audit_event_{suffix.lower()}_{seq}",
            "timestamp_ns": base_ns + offset * 1_000_000_000,
            "span_id_b64": span_id_b64,
            "trace_id_b64": trace_id_b64,
            "event_id": f"01HFFIEC{suffix}000000000000000{seq:02X}",
            "data": f"payload-{suffix.lower()}-{seq}",
            "step": step,
        })
    return events


def main() -> None:
    sk_v1 = session_key(IKM_V1, TENANT_ID)
    kfp_v1 = key_fingerprint(TENANT_ID, IKM_V1)
    digest = hkdf_inputs_digest(TENANT_ID)

    # Run A: events 1..3, base ns chosen to mirror central corpus style.
    run_a_events = _build_run_events(RUN_A, base_ns=1735689601000000000, base_step=1)
    # Run B: independent base ns, distinguishable payloads.
    run_b_events = _build_run_events(RUN_B, base_ns=1735693201000000000, base_step=10)

    def _chain_run(run_id: str, raw_events: list[dict]) -> list[dict]:
        out: list[dict] = []
        prev = b"\x00" * 32  # Run starts with genesis prev_hash regardless of other runs.
        for ev in raw_events:
            canonical = event_canonical(
                tenant_id=TENANT_ID,
                run_id=run_id,
                seq=ev["seq"],
                name=ev["name"],
                timestamp_ns=ev["timestamp_ns"],
                span_id_b64=ev["span_id_b64"],
                trace_id_b64=ev["trace_id_b64"],
                event_id=ev["event_id"],
                data=ev["data"],
                step=ev["step"],
            )
            ph = payload_hash(sk_v1, prev, canonical)
            out.append({
                "seq": ev["seq"],
                "event_canonical_hex": canonical.hex(),
                "prev_hash_hex": prev.hex(),
                "payload_hash_hex": ph.hex(),
                "key_version": 1,
                "key_fingerprint_hex": kfp_v1.hex(),
            })
            prev = ph
        return out

    chain_a = _chain_run(RUN_A, run_a_events)
    chain_b = _chain_run(RUN_B, run_b_events)

    # Daily Merkle aggregates all six payload_hash values in (run_id, seq) order.
    # RUN_A < RUN_B lexicographically, so chain_a precedes chain_b.
    all_hashes = [bytes.fromhex(e["payload_hash_hex"]) for e in chain_a] + [
        bytes.fromhex(e["payload_hash_hex"]) for e in chain_b
    ]
    merkle_root = merkle_root_rfc6962(all_hashes)

    sign_payload = build_sign_payload(
        algorithm=ALGORITHM,
        format_version=FORMAT_VERSION,
        tenant_id=TENANT_ID,
        seal_date=SEAL_DATE,
        merkle_root_hex=merkle_root.hex(),
        hkdf_inputs_digest_hex=digest.hex(),
        cadence=CADENCE,
        dev_mode_byte=DEV_MODE_BYTE,
    )

    fixture = {
        "_about": (
            "FFIEC v1 test-vector case 003-multi-run-same-day. Two runs in one "
            "tenant-day; run B's first event resets prev_hash to 32 zero bytes "
            "(cross-run isolation). Daily Merkle aggregates all 6 payload_hash "
            "values in (run_id, seq) order per spec §4.2."
        ),
        "inputs": {
            "HKDF_SALT": HKDF_SALT.decode("utf-8"),
            "HKDF_INFO_BASE": HKDF_INFO_BASE.decode("utf-8"),
            "tenant_id": TENANT_ID,
            "run_id_a": RUN_A,
            "run_id_b": RUN_B,
            "ikm_v1_hex": IKM_V1_HEX,
            "seal_date": SEAL_DATE,
            "algorithm": ALGORITHM,
            "format_version": FORMAT_VERSION,
            "cadence": CADENCE,
            "dev_mode": False,
        },
        "expected": {
            "session_key_v1_hex": sk_v1.hex(),
            "key_fingerprint_v1_hex": kfp_v1.hex(),
            "hkdf_inputs_digest_hex": digest.hex(),
            "merkle_root_daily_hex": merkle_root.hex(),
            "sign_payload_hex": sign_payload.hex(),
            "sign_payload_text": sign_payload.decode("utf-8"),
        },
        "run_a_chain": chain_a,
        "run_b_chain": chain_b,
        "merkle_leaves_in_order": [
            {"run_id": RUN_A, "seq": e["seq"], "payload_hash_hex": e["payload_hash_hex"]}
            for e in chain_a
        ] + [
            {"run_id": RUN_B, "seq": e["seq"], "payload_hash_hex": e["payload_hash_hex"]}
            for e in chain_b
        ],
    }

    with open(os.path.join(HERE, "fixture.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    expected = {
        "_about": (
            "Expected outputs for fixture.json. A conforming v1 verifier MUST "
            "reproduce these byte-for-byte. The Merkle root aggregates all 6 "
            "payload_hash values in (run_id, seq) ordering. The sign_payload "
            "uses the post-2026-05-07 spec §4.3 form (cadence + dev_mode bound)."
        ),
        "session_key_v1_hex": sk_v1.hex(),
        "key_fingerprint_v1_hex": kfp_v1.hex(),
        "hkdf_inputs_digest_hex": digest.hex(),
        "run_a_payload_hashes_hex": [e["payload_hash_hex"] for e in chain_a],
        "run_b_payload_hashes_hex": [e["payload_hash_hex"] for e in chain_b],
        "run_b_first_prev_hash_hex": chain_b[0]["prev_hash_hex"],
        "merkle_root_daily_hex": merkle_root.hex(),
        "sign_payload_hex": sign_payload.hex(),
        "sign_payload_text": sign_payload.decode("utf-8"),
    }

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(expected, f, indent=2, ensure_ascii=False)
        f.write("\n")

    # Post-condition prints for human inspection.
    print("[003] hkdf_inputs_digest:", digest.hex())
    print("[003] session_key_v1:    ", sk_v1.hex())
    print("[003] key_fingerprint_v1:", kfp_v1.hex())
    print("[003] run A payload_hashes:")
    for e in chain_a:
        print("       seq", e["seq"], e["payload_hash_hex"])
    print("[003] run B payload_hashes:")
    for e in chain_b:
        print("       seq", e["seq"], e["payload_hash_hex"])
    print("[003] run B[0] prev_hash:", chain_b[0]["prev_hash_hex"])
    print("[003] daily merkle_root: ", merkle_root.hex())
    print("[003] sign_payload SHA-256:", hashlib.sha256(sign_payload).hexdigest())


if __name__ == "__main__":
    main()
