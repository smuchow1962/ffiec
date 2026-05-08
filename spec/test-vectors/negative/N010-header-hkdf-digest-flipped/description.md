# N010-header-hkdf-digest-flipped

## Expected verifier outcome

```
Status: FAIL
Step:   2
Reason: header HKDF inputs do not match running v1 inputs
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Flip one byte of header.hkdf_inputs_digest. The recomputed digest from the running v1 constants does not match.

## Conformance test

Verifier MUST execute spec §7 steps 1..(2-1) cleanly and fail at step 2 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 2 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 2 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 2 (the institution-side response)
