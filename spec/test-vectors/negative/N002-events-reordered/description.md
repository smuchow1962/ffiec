# N002-events-reordered

## Expected verifier outcome

```
Status: FAIL
Step:   6
Reason: chain link broken at seq 2
```

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs.

Swap events at seq=2 and seq=4 in the file (without recomputing prev_hash chain). The structural walk at step 6 detects entry.prev_hash does not match expected_prev_hash.

## Conformance test

Verifier MUST execute spec §7 steps 1..(6-1) cleanly and fail at step 6 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 6 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 6 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 6 (the institution-side response)
