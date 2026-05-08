# N016-prev-hash-and-payload-recomputed-by-attacker

## Expected verifier outcome

```
Status: FAIL
Step:   6
Reason: chain link broken at seq 2
```

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs.

Attacker without IKM substitutes entry-2's prev_hash AND attempts to recompute payload_hash. Without the session_key, the recomputed payload_hash is wrong. Step 6 structural walk catches the prev_hash substitution; even if relaxed, step 9 catches the wrong MAC because the verifier feeds expected_prev_hash (not entry.prev_hash) into the MAC compute, and the MAC over (expected_prev_hash || canonical_bytes) under the real session_key does not match the attacker's MAC.

## Conformance test

Verifier MUST execute spec §7 steps 1..(6-1) cleanly and fail at step 6 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 6 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 6 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 6 (the institution-side response)
