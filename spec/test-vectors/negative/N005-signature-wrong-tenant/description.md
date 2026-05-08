# N005-signature-wrong-tenant

## Expected verifier outcome

```
Status: FAIL
Step:   11
Reason: signature verification failed
```

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs.

Use a signature produced for a different tenant_id over the same merkle_root. The sign_payload includes tenant_id; the signature does not validate against the file's tenant_id.

## Conformance test

Verifier MUST execute spec §7 steps 1..(11-1) cleanly and fail at step 11 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 11 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 11 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 11 (the institution-side response)
