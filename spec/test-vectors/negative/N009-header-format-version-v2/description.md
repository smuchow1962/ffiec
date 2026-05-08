# N009-header-format-version-v2

## Expected verifier outcome

```
Status: FAIL
Step:   1
Reason: format_version v2 not supported by this verifier (running v1)
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify the header's format_version from v1 to v2. The verifier (running v1) refuses at step 1 most-specific-first.

## Conformance test

Verifier MUST execute spec §7 steps 1..(1-1) cleanly and fail at step 1 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 1 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 1 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 1 (the institution-side response)
