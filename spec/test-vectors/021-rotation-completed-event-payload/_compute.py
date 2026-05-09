# -*- coding: utf-8 -*-
"""Compute the JCS-canonical bytes for FFIEC v1 test-vector case
021-rotation-completed-event-payload.

Spec §10.28 normates the `master.rotation.completed` operational event
payload that streaming-mode verifiers (§10.29) consume. The payload's
five fields are byte-locked across implementations; the
`rotation_at_utc` timestamp's RFC 3339 form with exactly 6-digit
microsecond precision and the trailing `Z` is load-bearing for
sub-second cadence.

This case pins the JCS-canonical bytes for one rotation event so a
clean-room implementation proves it produces the same bytes for the
same inputs. JCS canonicalization (RFC 8785) is what the Herald.Py
SDK uses for its event canonical form (verified against the reference
JCS library `jcs`).

Run with: python _compute.py
"""

from __future__ import annotations

import hashlib
import json
import os

import jcs

HERE = os.path.dirname(os.path.abspath(__file__))


# Pinned inputs for this case — a sub-second-cadence institution
# rotating IKM at 12:34:56.789012 UTC. The microsecond value is
# non-zero on purpose so any implementation that silently truncates
# to whole-second precision (a common bug under naive datetime
# formatting) produces a byte-different rotation_at_utc value and
# fails the SHA-256 pin loudly.
EVENT_TYPE = "master.rotation.completed"
CADENCE = "per_second"  # streaming-mode (§10.27)
PRIOR_KEY_VERSION = 1
NEW_KEY_VERSION = 2
ROTATION_AT_UTC = "2026-05-07T12:34:56.789012Z"  # 6-digit microsec + Z


def build_rotation_event_payload() -> dict:
    """Return the §10.28 normative event payload as a dict."""
    return {
        "event": EVENT_TYPE,
        "cadence": CADENCE,
        "prior_key_version": PRIOR_KEY_VERSION,
        "new_key_version": NEW_KEY_VERSION,
        "rotation_at_utc": ROTATION_AT_UTC,
    }


def main() -> None:
    payload = build_rotation_event_payload()

    # Use the RFC 8785 reference library so the corpus's canonical
    # bytes are the JCS reference implementation's output, not a
    # hand-rolled approximation.
    canonical_bytes = jcs.canonicalize(payload)
    canonical_text = canonical_bytes.decode("utf-8")
    canonical_sha256 = hashlib.sha256(canonical_bytes).hexdigest()

    input_record = {
        "_about": (
            "Input fixture for case 021-rotation-completed-event-payload — "
            "the §10.28 master.rotation.completed event payload byte-form "
            "pin. The `rotation_at_utc` value uses 6-digit microsecond "
            "precision and a trailing 'Z' per §10.28 normative; whole-"
            "second-resolution timestamps would erase the crossing-"
            "interval signal at sub-second cadence."
        ),
        "event_payload": payload,
    }
    with open(os.path.join(HERE, "input.json"), "w", encoding="utf-8") as f:
        json.dump(input_record, f, indent=2, ensure_ascii=False)
        f.write("\n")

    # expected_canonical.txt — JCS-canonical bytes verbatim. JCS
    # produces ASCII-only output for this payload (no non-ASCII string
    # values), but write in binary mode so platform line-ending
    # conversion cannot corrupt the bytes.
    with open(os.path.join(HERE, "expected_canonical.txt"), "wb") as f:
        f.write(canonical_bytes)

    with open(
        os.path.join(HERE, "expected_canonical_sha256.txt"),
        "w",
        encoding="utf-8",
        newline="\n",
    ) as f:
        f.write(canonical_sha256 + "\n")

    print("[021] event:                  ", payload["event"])
    print("[021] rotation_at_utc:        ", ROTATION_AT_UTC)
    print("[021] canonical byte length:  ", len(canonical_bytes))
    print("[021] canonical SHA-256:      ", canonical_sha256)
    print("[021] canonical text:         ", canonical_text)


if __name__ == "__main__":
    main()
