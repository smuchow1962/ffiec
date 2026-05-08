# Computation method — case 008-jcs-edge-cases

This file records the provenance of `fixture.json` and `expected.json`. It exists so a future maintainer can regenerate the bytes deterministically, and so a reviewer can confirm the bytes were produced by a real, identifiable JCS implementation rather than transcribed by hand.

## Date computed

2026-05-07

## Primary implementation (produced `expected.json`)

- **Package**: `jcs` (Python)
- **Version**: 0.2.1
- **Source**: <https://pypi.org/project/jcs/>
- **Author**: Anders Rundgren — the editor of RFC 8785
- **Install command**: `pip install jcs==0.2.1`
- **Entry point used**: `jcs.canonicalize(value: Any) -> bytes`

Choosing the RFC author's own implementation as the primary is deliberate. Disagreement between this implementation and another one is more likely to indicate a bug in the other implementation than in this one.

### Library file hashes (SHA-256, lowercase hex)

```
__init__.py  85 bytes   6c1417adb5aa3a2fd64fac4aa29a75bf92abf56bbc47e0e183aa88d95109b155
_jcs.py      18228 bytes  f34cd9f5889dfa8ec045f5c2ba3f94c36e5fa11333086075e55803041a00887e
ntoj.py      3954 bytes   0adf12a7834986eb331c003e3f609340cd1843ebf0d0bba5cbd36c37cbd5b77d
```

A future regeneration whose library hashes match the values above produces byte-identical fixtures. If any hash differs, the regeneration MUST be cross-validated against the secondary implementation before the fixtures are accepted.

## Secondary implementation (cross-validation)

- **File**: `_validate_inline.py` in this directory
- **Approach**: minimal RFC 8785 reference implemented from the RFC text directly. Not a production canonicaliser — its only purpose is to catch a bug in the primary library.
- **Coverage**: every case in the fixture, including the rejection cases. The validator decodes the `__special__` markers in `fixture.json` to native non-finite floats before invoking its canonicaliser, the same way a conforming implementation does.
- **File hash (SHA-256)**: see below.

### Cross-validation result

```
Matched: 58
Rejection match: 3
Disagreements: 0
```

All 61 cases agree across the two independent implementations. There were no disagreements to report to the working group.

The validation script self-checks at every regeneration: a non-zero exit from `_validate_inline.py` blocks the fixture from being committed.

## Generator and validator file hashes (SHA-256, lowercase hex)

```
_compute.py          16121 bytes  c3fd5b2473bfee3e91cdaee719fb2d49f2d85ec38f0d9e4d2f5147858f75148b
_validate_inline.py  8043 bytes   b1bb555bd4618d31de31eabedfde2b5ba60791dc1a67c3c817ef385edc9ea6ae
```

## Output file hashes (SHA-256, lowercase hex)

```
fixture.json   108532 bytes  9fbea20632dc6a26b07f5129d43cf97804511598429c499f6eab4660360da1dc
expected.json  228735 bytes  de10ee0ef9b2a1f32976bb21b2ee4e1a1606f22b613e3cdb9e1fbd4966bc4b65
```

## Reproducing the fixtures

1. Install Python 3.10 or later (any 3.x with stable repr-shortest-round-trip floats).
2. `pip install jcs==0.2.1`.
3. Confirm the library file hashes above by running `python -c "import hashlib, jcs, os; d=os.path.dirname(jcs.__file__); print({f: hashlib.sha256(open(os.path.join(d,f),'rb').read()).hexdigest() for f in ['__init__.py','_jcs.py','ntoj.py']})"`.
4. Run `python _compute.py`. The script writes `fixture.json` and `expected.json`, then runs a round-trip self-check (re-reads the fixture, re-canonicalises, compares against expected). A non-zero exit blocks the regeneration.
5. Run `python _validate_inline.py`. The script independently canonicalises each fixture input and compares against `expected.json`. A non-zero exit blocks the regeneration.

If both scripts exit zero, the fixtures are byte-identical to the committed versions.

## Limitations and disclosures

- The Python `jcs` library treats all JSON numbers as IEEE-754 doubles (RFC 8785 §3.2.2.3 mandate). Integers above 2^53 lose precision when canonicalised — see the `above_safe_integer` case for the documented behaviour.
- The library raises `ValueError` for NaN and Infinity, and `TypeError` for unsupported Python types. The error class name is recorded in `expected.json`. A future patch release that changes the error class would require regenerating the fixture; the round-trip self-check would catch the drift.
- The fixture deliberately omits 1MB and 10MB string cases for repository-size reasons. An implementation that wants to claim conformance for those sizes can compute the SHA-256 of `{"v": "a" * 1000000}` and similar and compare against an out-of-band published hash.
- The 1000-level deeply-nested case is omitted because Python's default recursion limit is 1000 and most production JSON parsers reject far below that. Implementations supporting deeper nesting can publish their own extension fixtures.

## Independent third-party validation

The fixture was validated against two implementations (the published `jcs` 0.2.1 package and the inline reference). A future iteration may add validation against:

- Anders Rundgren's reference `webpki.org` JCS Java implementation (canonical source for the spec)
- Cyberphone's JavaScript JCS implementation (`canonicaljson.js`)
- The Go `gowebpki/jcs` package

All of these implementations existed at the time `expected.json` was computed; the fixture is amenable to validation against any of them. Disagreements found in future cross-validation runs MUST be filed as a finding for the FFIEC working group; the fixture is NOT silently updated to match a disagreeing implementation without working-group review.
