# N011-header-genesis-nonzero

## Expected verifier outcome

```
Status: FAIL
Step:   3
Reason: header genesis_hash does not match v1 constant
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify header.genesis_hash to non-zero bytes (e.g., 0x01 * 32). Verifier asserts the v1 constant of 32 zero bytes.

## Conformance test

Verifier MUST execute spec §7 steps 1..(3-1) cleanly and fail at step 3 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 3 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 3 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 3 (the institution-side response)
