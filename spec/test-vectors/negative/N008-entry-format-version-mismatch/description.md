# N008-entry-format-version-mismatch

## Expected verifier outcome

```
Status: FAIL
Step:   5
Reason: format_version mismatch at seq 1
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify one entry's format_version from v1 to v2. Header still says v1; the per-entry vs header mismatch fires at step 5.

## Conformance test

Verifier MUST execute spec §7 steps 1..(5-1) cleanly and fail at step 5 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 5 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 5 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 5 (the institution-side response)
