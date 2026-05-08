# N001-payload-hash-bit-flip

## Expected verifier outcome

```
Status: FAIL
Step:   9
Reason: payload_hash MAC mismatch at seq 1
```

## Tampering recipe

Start from `../../001-single-event-empty-prev/` inputs.

Flip one byte of payload_hash on the chain entry. Verifier walks steps 1-8, reaches step 9 MAC compute, fails because the recomputed MAC over (expected_prev_hash || canonical_bytes) does not match the flipped payload_hash.

## Conformance test

Verifier MUST execute spec §7 steps 1..(9-1) cleanly and fail at step 9 with the exact reason string above. A verifier that produces `Status: PASS` is broken; a verifier that fails at a different step (or with a different reason string) is non-conforming for the rework's named-failure-mode taxonomy.

## See also

- spec §7 step 9 (the failure step)
- `docs/regulator-pack/finding-language.md` row for step 9 (the examiner-finding paragraph)
- `docs/incident-response-playbook.md` IR scenario mapped to step 9 (the institution-side response)
