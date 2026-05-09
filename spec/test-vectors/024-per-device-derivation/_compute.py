# -*- coding: utf-8 -*-
"""Compute per-device HKDF info bytes + session keys for FFIEC v1
test-vector case 024-per-device-derivation.

Spec §10.32 normates the per-device session-key derivation extension
to §4.1's tenant-bound HKDF. The extended derivation adds a third
name segment binding a device identity:

  info = HKDF_INFO_BASE || "|" || utf8(tenant_id) || "|" || utf8(device_id)
  session_key = HKDF-SHA256(ikm, salt = HKDF_SALT, info = info, length = 32)

The `|` separator (byte `0x7c`) is the same fixed separator §4.1 uses.
A §10.32-bound chain is byte-distinct from a §4.1-only chain — a §4.1
verifier walking a §10.32 chain fails its MAC checks because the
device-bound info produces a different session key.

This case pins the exact info bytes, derived session keys, and key
fingerprints for two devices under the same tenant + IKM, plus the
§4.1-only baseline derivation for byte-distinctness comparison.

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import hmac
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


# Constants from spec §4.1 (FFIEC-conformance posture)
HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"
HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"
HKDF_LENGTH = 32

# Pinned inputs — matching the central-corpus tenant_id and a 32-byte
# IKM derived deterministically from a label so any reader can
# recompute the IKM bytes from first principles.
TENANT_ID = "tenant-ffiec-test-1"
IKM_LABEL = "ffiec-cross-language-test-vector-ikm-v1-32bytes!"
IKM = IKM_LABEL.encode("utf-8")  # exactly 48 bytes
# Two device_id values exercising distinct identity-authority shapes
DEVICE_A = "tablet-A1B2-C3D4"          # field-tablet inventory tag
DEVICE_B = "tpm-7eb3-c000-0001-aabb"   # TPM 2.0 endorsement-key fingerprint


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
    """§4.1 base derivation info (no device segment)."""
    return HKDF_INFO_BASE + b"|" + tenant_id.encode("utf-8")


def info_for_device(tenant_id: str, device_id: str) -> bytes:
    """§10.32 extended derivation info — adds a third `|`-separated segment."""
    return (
        HKDF_INFO_BASE
        + b"|" + tenant_id.encode("utf-8")
        + b"|" + device_id.encode("utf-8")
    )


def session_key(ikm: bytes, info: bytes) -> bytes:
    return hkdf_sha256(ikm, HKDF_SALT, info, HKDF_LENGTH)


def key_fingerprint(tenant_id: str, ikm: bytes) -> bytes:
    """Doc §3.5 fingerprint — identical recipe under §4.1 and §10.32.

    The fingerprint binds tenant + IKM only; device_id is NOT part of
    the fingerprint because the fingerprint exists to detect IKM-swap
    misconfigurations at lookup time, not device-binding drift.
    """
    return hashlib.sha256(tenant_id.encode("utf-8") + ikm).digest()[:16]


def main() -> None:
    info_base = info_for_tenant(TENANT_ID)
    info_a = info_for_device(TENANT_ID, DEVICE_A)
    info_b = info_for_device(TENANT_ID, DEVICE_B)

    sk_base = session_key(IKM, info_base)
    sk_a = session_key(IKM, info_a)
    sk_b = session_key(IKM, info_b)

    fp = key_fingerprint(TENANT_ID, IKM)

    # Sanity: the three session keys must be pairwise distinct. If
    # any two collide, the generator broke and the fixture is invalid.
    keys = {sk_base, sk_a, sk_b}
    assert len(keys) == 3, "session keys collided — generator bug"

    fixture = {
        "_about": (
            "Per-device HKDF derivation byte pin per spec §10.32. The "
            "extended `info` adds a third `|`-separated segment binding "
            "a device identity. A §4.1-only verifier walking a §10.32 "
            "chain fails its MAC checks because the device-bound info "
            "produces a different session key. This fixture pins the "
            "info bytes + session keys for two devices under the same "
            "tenant + IKM, plus the §4.1-only baseline for byte-"
            "distinctness comparison."
        ),
        "constants": {
            "hkdf_salt_utf8": HKDF_SALT.decode("utf-8"),
            "hkdf_salt_hex": HKDF_SALT.hex(),
            "hkdf_info_base_utf8": HKDF_INFO_BASE.decode("utf-8"),
            "hkdf_info_base_hex": HKDF_INFO_BASE.hex(),
            "hkdf_length": HKDF_LENGTH,
            "separator_byte": "0x7c (the ASCII pipe character)",
        },
        "inputs": {
            "tenant_id": TENANT_ID,
            "ikm_label": IKM_LABEL,
            "ikm_hex": IKM.hex(),
            "ikm_byte_length": len(IKM),
            "device_a": DEVICE_A,
            "device_b": DEVICE_B,
        },
        "computed_info_bytes": {
            "section_4_1_baseline": {
                "info_recipe": "HKDF_INFO_BASE || '|' || utf8(tenant_id)",
                "info_utf8": info_base.decode("utf-8"),
                "info_hex": info_base.hex(),
                "info_byte_length": len(info_base),
            },
            "section_10_32_device_a": {
                "info_recipe": (
                    "HKDF_INFO_BASE || '|' || utf8(tenant_id) || '|' || "
                    "utf8(device_a)"
                ),
                "info_utf8": info_a.decode("utf-8"),
                "info_hex": info_a.hex(),
                "info_byte_length": len(info_a),
            },
            "section_10_32_device_b": {
                "info_recipe": (
                    "HKDF_INFO_BASE || '|' || utf8(tenant_id) || '|' || "
                    "utf8(device_b)"
                ),
                "info_utf8": info_b.decode("utf-8"),
                "info_hex": info_b.hex(),
                "info_byte_length": len(info_b),
            },
        },
        "expected": {
            "section_4_1_baseline_session_key_hex": sk_base.hex(),
            "section_10_32_device_a_session_key_hex": sk_a.hex(),
            "section_10_32_device_b_session_key_hex": sk_b.hex(),
            "key_fingerprint_hex": fp.hex(),
            "_byte_distinctness_note": (
                "All three session keys are pairwise distinct. The "
                "device_id segment in the §10.32 info changes the HKDF "
                "expand input, so the same (ikm, tenant_id) pair "
                "produces a different 32-byte session key for each "
                "device_id value AND for the §4.1 baseline (no device "
                "segment)."
            ),
        },
    }

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    print(f"[024] info_4_1   ({len(info_base):3d} bytes): {info_base!r}")
    print(f"[024] info_10_32_a ({len(info_a):3d} bytes): {info_a!r}")
    print(f"[024] info_10_32_b ({len(info_b):3d} bytes): {info_b!r}")
    print(f"[024] sk_4_1_baseline:    {sk_base.hex()}")
    print(f"[024] sk_10_32_device_a:  {sk_a.hex()}")
    print(f"[024] sk_10_32_device_b:  {sk_b.hex()}")
    print(f"[024] key_fingerprint:    {fp.hex()}")


if __name__ == "__main__":
    main()
