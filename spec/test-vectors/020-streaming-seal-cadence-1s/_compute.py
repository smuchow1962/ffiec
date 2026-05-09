# -*- coding: utf-8 -*-
"""Compute the sign_payload byte form for FFIEC v1 test-vector case
020-streaming-seal-cadence-1s.

Spec §10.27 normates `cadence = "per_second"` as a §4.2.1 enumerated
value. The cadence value is bound under the §4.3 sign_payload form on
its dedicated line; a tampered cadence value in the seal record is
detected at signature verification because the verifier reconstructs
sign_payload using the field as written.

This case pins the v1.0b 12-line sign_payload form for a streaming-
mode 1-second-cadence seal so a clean-room implementation proves it
emits the §10.27 streaming cadence value byte-identical to other
implementations (the per_second / hourly / daily byte-distinct rule
from §10.27).

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


# Pinned inputs for this case — chosen to be visually distinguishable
# from the v1.0b mixed-day inputs in case 018 (different tenant_id +
# seal_date + merkle_root / hkdf_inputs_digest placeholders) so a
# careless copy-paste between fixtures does not silently produce the
# same bytes for two different cases.
ALGORITHM = "ed25519"
FORMAT_VERSION = "v1"
TENANT_ID = "demo-bank-streaming-001"
# A 1-second cadence interval landing on a clean second boundary.
# seal_period_start_utc is institution-trusted ledger-side metadata
# (not bound into sign_payload per §10.27 + §4.2 schema notes); it is
# carried in input.json for completeness.
SEAL_DATE = "2026-05-07"
SEAL_PERIOD_START_UTC = "2026-05-07T12:34:56.000000Z"
# Placeholder digest values — the v1.0b form is independent of the
# specific 32-byte values; the case pins the cadence byte form.
MERKLE_ROOT_HEX = (
    "0001020304050607080910111213141516171819202122232425262728293031"
)
HKDF_INPUTS_DIGEST_HEX = (
    "ffeeddccbbaa99887766554433221100ffeeddccbbaa99887766554433221100"
)
CADENCE = "per_second"  # §10.27 streaming-mode 1-second cadence
DEV_MODE = False  # serialized as the single ASCII byte "0"
KEY_VERSIONS = [1]  # single key version on the 1-second interval
KMS_HANDLE_URIS = [
    "aws-kms:arn:aws:kms:us-east-1:123456789012:key/streaming-mock-1"
]


def _canonicalize_key_versions(versions: list[int]) -> str:
    """v1.0b key_versions_canon line — sorted-distinct ascending decimal."""
    distinct_sorted = sorted({int(v) for v in versions})
    if distinct_sorted and distinct_sorted[0] < 0:
        raise ValueError(f"key_version must be >= 0, got {distinct_sorted[0]}")
    return ",".join(str(v) for v in distinct_sorted)


def _compute_kms_handle_uris_digest(uris: list[str]) -> bytes:
    """v1.0b kms_handle_uris_digest — SHA-256 over canonical sorted form."""
    canonical = "\n".join(sorted(set(uris))).encode("utf-8")
    return hashlib.sha256(canonical).digest()


def build_sign_payload_v1_0b(
    *,
    algorithm: str,
    format_version: str,
    tenant_id: str,
    seal_date: str,
    merkle_root_hex: str,
    hkdf_inputs_digest_hex: str,
    cadence: str,
    dev_mode: bool,
    key_versions: list[int],
    kms_handle_uris: list[str],
) -> bytes:
    """Construct the v1.0b 12-line sign_payload byte form per §4.3.

    Lines (each terminated with \\n except the terminal field):
      1. ffiec.chain-of-custody.v1
      2. v1.0b
      3. algorithm
      4. format_version
      5. tenant_id
      6. seal_date (YYYY-MM-DD)
      7. hex(merkle_root) — 64 lowercase hex chars
      8. hex(hkdf_inputs_digest) — 64 lowercase hex chars
      9. cadence (§10.27 enumerated value)
     10. dev_mode (single ASCII byte: "1" if true, "0" if false)
     11. key_versions_canon (comma-separated sorted-distinct decimals)
     12. hex(kms_handle_uris_digest) — 64 lowercase hex chars (terminal,
         no trailing newline)
    """
    dev_mode_byte = "1" if dev_mode else "0"
    key_versions_canon = _canonicalize_key_versions(key_versions)
    kms_uris_digest_hex = _compute_kms_handle_uris_digest(kms_handle_uris).hex()

    text = (
        "ffiec.chain-of-custody.v1\n"
        + "v1.0b\n"
        + algorithm + "\n"
        + format_version + "\n"
        + tenant_id + "\n"
        + seal_date + "\n"
        + merkle_root_hex + "\n"
        + hkdf_inputs_digest_hex + "\n"
        + cadence + "\n"
        + dev_mode_byte + "\n"
        + key_versions_canon + "\n"
        + kms_uris_digest_hex
    )
    return text.encode("utf-8")


def main() -> None:
    sign_payload = build_sign_payload_v1_0b(
        algorithm=ALGORITHM,
        format_version=FORMAT_VERSION,
        tenant_id=TENANT_ID,
        seal_date=SEAL_DATE,
        merkle_root_hex=MERKLE_ROOT_HEX,
        hkdf_inputs_digest_hex=HKDF_INPUTS_DIGEST_HEX,
        cadence=CADENCE,
        dev_mode=DEV_MODE,
        key_versions=KEY_VERSIONS,
        kms_handle_uris=KMS_HANDLE_URIS,
    )

    sign_payload_sha256_hex = hashlib.sha256(sign_payload).hexdigest()
    kms_uris_digest_hex = _compute_kms_handle_uris_digest(KMS_HANDLE_URIS).hex()

    # input.json — full input record for the case.
    input_record = {
        "_about": (
            "Input fixture for case 020-streaming-seal-cadence-1s — v1.0b "
            "sign_payload exercise of the §10.27 streaming-mode 1-second "
            "cadence value. The non-bound seal_period_start_utc field is "
            "included for completeness but is institution-trusted ledger-"
            "side metadata per §4.2 — it is NOT part of the sign_payload "
            "byte form."
        ),
        "seal": {
            "sign_payload_version": "v1.0b",
            "algorithm": ALGORITHM,
            "format_version": FORMAT_VERSION,
            "tenant_id": TENANT_ID,
            "seal_date": SEAL_DATE,
            "seal_period_start_utc": SEAL_PERIOD_START_UTC,
            "merkle_root_hex": MERKLE_ROOT_HEX,
            "hkdf_inputs_digest_hex": HKDF_INPUTS_DIGEST_HEX,
            "cadence": CADENCE,
            "dev_mode": DEV_MODE,
            "key_versions": KEY_VERSIONS,
        },
        "per_event_distribution": {
            "_comment": (
                "Distinct key_version and kms_handle_uri values present on "
                "the 1-second interval's chain entries. Case 020 verifies "
                "the sign_payload byte form, not the per-event walk."
            ),
            "key_versions_distinct": KEY_VERSIONS,
            "kms_handle_uris_distinct": KMS_HANDLE_URIS,
        },
        "computed_canonical_fields": {
            "key_versions_canon": _canonicalize_key_versions(KEY_VERSIONS),
            "kms_handle_uris_digest_hex": kms_uris_digest_hex,
        },
    }
    with open(os.path.join(HERE, "input.json"), "w", encoding="utf-8") as f:
        json.dump(input_record, f, indent=2, ensure_ascii=False)
        f.write("\n")

    # expected_sign_payload.txt — the canonical bytes verbatim, terminated
    # with NO trailing newline (the terminal field is the kms_handle_uris
    # digest). Open in binary mode so the platform's line-ending
    # conversion (Windows CRLF) cannot corrupt the embedded \n bytes.
    with open(os.path.join(HERE, "expected_sign_payload.txt"), "wb") as f:
        f.write(sign_payload)

    # expected_sign_payload_sha256.txt — quick compare without diffing
    # 200+ bytes by eye.
    with open(
        os.path.join(HERE, "expected_sign_payload_sha256.txt"),
        "w",
        encoding="utf-8",
        newline="\n",
    ) as f:
        f.write(sign_payload_sha256_hex + "\n")

    # Post-condition prints.
    print("[020] sign_payload byte length:", len(sign_payload))
    print("[020] sign_payload SHA-256:    ", sign_payload_sha256_hex)
    print("[020] cadence:                 ", CADENCE)
    print("[020] kms_handle_uris_digest:  ", kms_uris_digest_hex)


if __name__ == "__main__":
    main()
