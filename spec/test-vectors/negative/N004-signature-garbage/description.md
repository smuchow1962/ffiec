# N004-signature-garbage

## Expected verifier outcome

```
Status: FAIL
Step:   11
Reason: signature verification failed
```

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs.

Replace the seal record's signature with 64 random bytes. Steps 1-10 pass; step 11 Ed25519.Verify fails.

## Conformance test

Verifier MUST execute spec §7 steps 1..(11-1) cleanly and fail at step 11 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 11 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 11 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 11 (the institution-side response)
