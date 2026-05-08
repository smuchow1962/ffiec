# N012-cross-chain-tenant-mismatch

## Expected verifier outcome

```
Status: FAIL
Step:   4
Reason: cross-chain lift detected at seq 1 (event.tenant_id mismatch)
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify the event's tenant_id to a value different from header.tenant_id (e.g., 'tenant-other'). The cross-chain-lift defense at step 4 fires.

## Conformance test

Verifier MUST execute spec §7 steps 1..(4-1) cleanly and fail at step 4 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 4 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 4 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 4 (the institution-side response)
