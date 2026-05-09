# -*- coding: utf-8 -*-
"""Compute the attestation-envelope JCS-canonical bytes for FFIEC v1
test-vector case 025-attestation-android-keystore.

Spec §10.35 normates the `ffiec.chain.attestation` attribute on chain
entries plus a delegate-dispatch validator pattern that platform-
specific attestation parsers plug into at chain-walk time. The
attestation envelope schema is byte-locked across implementations.

This case pins the JCS-canonical bytes for an Android Keystore
attestation envelope so a clean-room implementation proves it
constructs the §10.35 envelope identically to other implementations.
The attestation-document body is a synthetic placeholder (case 025
exercises the **envelope** byte form, not Android Keystore CBOR/X.509
parsing — that's the platform validator's job).

Run with: python _compute.py
"""

from __future__ import annotations

import base64
import hashlib
import json
import os

import jcs

HERE = os.path.dirname(os.path.abspath(__file__))


# Pinned inputs — synthetic Android Keystore attestation envelope.
PLATFORM = "android_keystore"
# A 256-byte synthetic attestation document. In production this would
# be a CBOR-or-X.509 chain rooted at Google's Hardware-Backed Keys CA;
# case 025 pins the envelope byte form, so the document content is a
# deterministic placeholder (SHA-256 chain of a label, repeated to
# fill 256 bytes) so a reader can recompute it from first principles.
_DOCUMENT_LABEL = "025-android-keystore-synthetic-document-v1"


def _build_synthetic_document(seed_label: str, n_bytes: int) -> bytes:
    """Deterministic byte sequence for the placeholder attestation doc.

    Recipe: seed = SHA-256(label); body = seed || SHA-256(seed) ||
    SHA-256(SHA-256(seed)) || ... truncated to n_bytes. Anyone can
    recompute the exact byte sequence from the label.
    """
    out = b""
    h = hashlib.sha256(seed_label.encode("utf-8")).digest()
    while len(out) < n_bytes:
        out += h
        h = hashlib.sha256(h).digest()
    return out[:n_bytes]


_DOCUMENT_BYTES = _build_synthetic_document(_DOCUMENT_LABEL, 256)
ATTESTATION_DOCUMENT_B64 = base64.b64encode(_DOCUMENT_BYTES).decode("ascii")
ROOT_CERTIFICATE_ID = "google.hardware-backed-keys.root.2023"
# Optional public-key fingerprint — SHA-256 of a synthetic 32-byte
# device public key.
_DEVICE_PUBKEY = hashlib.sha256(
    b"025-android-keystore-synthetic-device-pubkey"
).digest()
PUBLIC_KEY_FINGERPRINT = hashlib.sha256(_DEVICE_PUBKEY).hexdigest()


def build_attestation_envelope() -> dict:
    """§10.35 attestation envelope — four fields, byte-locked schema."""
    return {
        "platform": PLATFORM,
        "attestation_document_b64": ATTESTATION_DOCUMENT_B64,
        "root_certificate_id": ROOT_CERTIFICATE_ID,
        "public_key_fingerprint": PUBLIC_KEY_FINGERPRINT,
    }


def main() -> None:
    envelope = build_attestation_envelope()
    canonical_bytes = jcs.canonicalize(envelope)
    canonical_sha256 = hashlib.sha256(canonical_bytes).hexdigest()

    input_record = {
        "_about": (
            "Input fixture for case 025-attestation-android-keystore — "
            "the §10.35 attestation envelope byte-form pin for the "
            "Android Keystore platform. The attestation document body "
            "is a synthetic placeholder; case 025 exercises the "
            "envelope byte form, not Android Keystore CBOR/X.509 "
            "parsing (that's the platform validator's job)."
        ),
        "envelope": envelope,
        "synthetic_document": {
            "label": _DOCUMENT_LABEL,
            "byte_length": len(_DOCUMENT_BYTES),
            "recipe": (
                "seed = SHA-256(utf8(label)); body = seed || SHA-256(seed) "
                "|| SHA-256(SHA-256(seed)) || ... truncated to byte_length"
            ),
            "first_32_bytes_hex": _DOCUMENT_BYTES[:32].hex(),
        },
        "synthetic_device_pubkey_sha256_hex": PUBLIC_KEY_FINGERPRINT,
    }
    with open(os.path.join(HERE, "input.json"), "w", encoding="utf-8") as f:
        json.dump(input_record, f, indent=2, ensure_ascii=False)
        f.write("\n")

    with open(os.path.join(HERE, "expected_canonical.txt"), "wb") as f:
        f.write(canonical_bytes)

    with open(
        os.path.join(HERE, "expected_canonical_sha256.txt"),
        "w",
        encoding="utf-8",
        newline="\n",
    ) as f:
        f.write(canonical_sha256 + "\n")

    print(f"[025] platform:                   {PLATFORM}")
    print(f"[025] doc base64 len:             {len(ATTESTATION_DOCUMENT_B64)}")
    print(f"[025] envelope canonical len:     {len(canonical_bytes)}")
    print(f"[025] envelope canonical SHA-256: {canonical_sha256}")


if __name__ == "__main__":
    main()
