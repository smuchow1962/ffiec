# N013 — Audit file ends mid-line (mid-write truncation, load-bearing rework case)

## Purpose

Verify the spec §4.1 mid-write truncation refusal. A verifier that silently passes a file whose last byte is not `\n` is non-conforming because it would silently accept a chain that lost its last entry to a writer-side crash. The verifier MUST refuse the file with a precise named reason at the file pre-flight, BEFORE walking any events.

## Tampering recipe

Start from `../../002-multi-event-same-run/` inputs (the five-event chain). Take the audit file's bytes and remove the trailing `\n` from the last line, OR remove the last 1-N bytes from the file (simulating a writer-process crash mid-append).

Two specific test files SHOULD ship in this case directory:

- `inputs/audit-truncated-no-newline.ndjson` — file ends with the JSON closing brace `}` (no trailing `\n`).
- `inputs/audit-truncated-mid-event.ndjson` — file ends mid-event in the middle of the last event's JSON serialization.

Both MUST be refused.

## Expected verifier outcome

```
Status: FAIL (file pre-flight)
Step:   pre-flight (before steps 1-12)
Reason: audit file ends mid-line — possible mid-write crash; the writer's last append did not complete
```

Specifically:

- The verifier MUST reject the file at the pre-flight check, before parsing the header.
- The verifier MUST NOT silently truncate to the last `\n` and proceed (a permissive verifier that does this would silently lose events).
- The verifier MUST report the failure as operational (writer-side crash recovery), distinct from chain integrity failures (steps 1-12).

## What this case proves

A verifier that passes this case — by refusing the file with `audit file ends mid-line — possible mid-write crash` — defends against the silent-event-loss failure mode the spec §4.1 mid-write truncation refusal exists to prevent. The IR Scenario 9 (audit_file.truncation_detected) is the operational disposition path; the verifier's job is to surface the truncation as a refusable input, not to silently work around it.

## Negative case provenance

This case exists to prevent the regression where a verifier implementer adds "permissive" mode that splits on `\n` and processes whatever lines parse — silently losing the truncated tail. The truncation refusal is a strict-mode check by design; the operational disposition is in the IR playbook, not in the verifier.
