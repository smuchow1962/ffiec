# N003-merkle-root-altered

## Expected verifier outcome

```
Status: FAIL
Step:   10
Reason: merkle root mismatch — ledger contents do not produce sealed root
```

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs.

Modify the seal record's merkle_root to arbitrary 32 bytes. Steps 1-9 pass; step 10 recomputes the Merkle root from the day's payload_hashes and detects mismatch.

## Conformance test

Verifier MUST execute spec §7 steps 1..(10-1) cleanly and fail at step 10 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 10 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 10 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 10 (the institution-side response)
