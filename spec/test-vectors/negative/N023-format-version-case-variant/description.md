# N023-format-version-case-variant

## Expected verifier outcome

```
Status: FAIL
Step:   1
Reason: format_version "V1" not supported by this verifier (running v1)
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify the chain file's header so `format_version = "V1"` (uppercase V; everything else valid). The `hkdf_inputs_digest`, `genesis_hash`, the entry's `format_version`, every per-event MAC, and the seal record are left valid. The verifier (running v1) refuses at step 1 because the literal string `"V1"` is NOT in the set of values a v1 verifier recognizes.

The sibling cases position case-variant matching as part of the version-string contract:

- `N009-header-format-version-v2/` covers an unrecognized major-version string (`"v2"`).
- `N022-format-version-v1-1/` covers an unrecognized minor-version string within the v1 family (`"v1.1"`).
- `N023-format-version-case-variant/` (this case) covers a case-variant of the recognized string (`"V1"`).

All three fail at step 1 — the verifier MUST refuse anything other than the literal `format_version` strings it knows. The byte-level contract is: a v1 verifier knows the lowercase ASCII string `"v1"` and nothing else. Uppercase `"V1"`, mixed-case `"V1"`, whitespace-padded `" v1 "`, and any other variant whose bytes differ from `0x76 0x31` is rejected.

This is the load-bearing closure for the case-sensitivity policy. The `format_version` field is bound into `sign_payload` per spec §4.3 and the byte form is normative. A verifier that lowercases the field before comparison would silently accept `"V1"` as `"v1"`, but the seal's signature is computed over the bytes actually written — so a `"V1"` chain whose seal was signed with `"V1"` on `sign_payload` line 4 would have its signature break under any verifier that lowercased the field before reconstructing `sign_payload`. The cleaner guarantee is the byte-level one: refuse anything that is not byte-identical to `"v1"`.

## Why case sensitivity matters

Three independent reasons:

1. **`sign_payload` is byte-bound.** The `format_version` field lives on line 4 of the v1.0a `sign_payload` per spec §4.3. The signature covers those bytes. A verifier that lowercases (or otherwise normalizes) the field before comparison can disagree with a verifier that does not, and the signature check is the only mechanical defense against that disagreement. The simplest and most portable rule is: bytes-on-the-wire must equal bytes-the-verifier-knows.

2. **Cross-implementation portability.** Two implementations that disagree about case-folding produce different verification results on the same chain. The conformance corpus exists so that two independent implementations cannot disagree on what v1 means; this case pins that the disagreement on case-folding is closed at the verifier-refuses-non-canonical-bytes layer.

3. **Future-version dispatch.** When v1.1 ships, its `format_version` string is `"v1.1"` (lowercase). A v1 verifier that case-folds would also accept `"V1.1"`, `"V1.1"`, etc. — and the v1.1 spec would be silently accepting four byte-distinct chain populations as the same version. Pinning case sensitivity at v1.0 means v1.1 inherits the same byte-level discipline.

## Conformance test

Verifier MUST execute spec §7 steps 1..(1-1) cleanly and fail at step 1 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

A verifier that emits the reason string from N009 or N022 on a `"V1"` header is non-conforming — the version string in the reason MUST reflect what the chain claimed (the literal `"V1"` bytes), not the canonical `"v2"` placeholder or the canonical `"v1.1"` placeholder.

## Fixture shape

`fixture.json` describes the modification recipe rather than carrying a full pre-tampered chain. A reference runner produces the tampered input from `../../001-single-event-empty-prev/` by overwriting the header's `format_version` field with the literal byte sequence `0x56 0x31` (`"V1"`).

## See also

- spec §7 step 1 (the failure step)
- `N009-header-format-version-v2/` (sibling case for unrecognized major version)
- `N022-format-version-v1-1/` (sibling case for unrecognized minor version)
- `docs/regulator-pack/finding-language.md` row for step 1 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 1 (the institution-side response)
