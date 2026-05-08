# N022-format-version-v1-1

## Expected verifier outcome

```
Status: FAIL
Step:   1
Reason: format_version v1.1 not supported by this verifier (running v1)
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify the chain file's header so `format_version = "v1.1"`. Everything else (`hkdf_inputs_digest`, `genesis_hash`, the entry's `format_version` field, every per-event MAC, the seal record) is left valid. The verifier (running v1) refuses at step 1 most-specific-first.

The sibling case `N009-header-format-version-v2/` covers the unrecognized major-version case (`"v2"`); this case covers the unrecognized minor-version case within the v1 family. Both fail at step 1 — the verifier MUST refuse anything other than the literal `format_version` strings it knows. A v1 verifier knows `"v1"` and nothing else; `"v1.1"` is rejected even though it shares the major-version prefix.

This is the load-bearing closure for the version-negotiation policy: a v1 verifier MUST NOT silently accept a `"v1.1"` chain on the assumption that minor-version bumps stay backward-compatible. The chain construction or seal binding may have changed in v1.1 — only a v1.1 verifier knows. A v1 verifier that accepts `"v1.1"` and then validates with v1 rules can produce a false-pass on a chain whose v1.1 construction differs in some byte that v1 rules don't check.

## Conformance test

Verifier MUST execute spec §7 steps 1..(1-1) cleanly and fail at step 1 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

A verifier that emits the reason string from N009 (`format_version v2 not supported by this verifier (running v1)`) on a `"v1.1"` header is non-conforming — the version string in the reason MUST reflect what the chain claimed, not the canonical "next major" placeholder.

## Fixture shape

`fixture.json` describes the modification recipe rather than carrying a full pre-tampered chain. A reference runner produces the tampered input from `../../001-single-event-empty-prev/` by overwriting the header's `format_version` field with the literal string `"v1.1"`.

## See also

- spec §7 step 1 (the failure step)
- `N009-header-format-version-v2/` (sibling case for unrecognized major version)
- `docs/regulator-pack/finding-language.md` row for step 1 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 1 (the institution-side response)
