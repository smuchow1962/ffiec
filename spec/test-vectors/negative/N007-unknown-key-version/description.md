# N007-unknown-key-version

## Expected verifier outcome

```
Status: FAIL
Step:   7
Reason: unknown key_version: no IKM for (tenant=tenant-ffiec-test-1, key_version=99) at seq 1
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Modify the chain entry's key_version to 99 (a value not in the test IKM registry). Step 7 lookup returns null; verifier reports unknown_key_version. NO MAC compute happens.

## Conformance test

Verifier MUST execute spec §7 steps 1..(7-1) cleanly and fail at step 7 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 7 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 7 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 7 (the institution-side response)
