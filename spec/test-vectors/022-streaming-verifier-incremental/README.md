# Case 022 — streaming-verifier state machine (§10.29)

## What this case verifies

Spec §10.29 normates the streaming-mode verifier procedure with three streaming-state codes (`4` streaming-all-pass-so-far, `5` streaming-anomaly-detected, `6` streaming-key-rotation-pending-confirmation), a transition table over input kinds (`chain_entry` / `seal` / `rotation`), and an end-of-stream finalize collapse to terminal verdicts (`0` PASS / `3` chain-anomaly).

This case pins the state-machine table and a set of scenario walks so a clean-room implementation proves it walks the §10.29 contract identically. Unlike cases 020 / 021 (which pin **byte forms**), case 022 pins **state-machine behaviour**: a clean-room verifier reads the scenarios from `expected.json` and confirms its own state-transition logic produces the same per-step verdicts and the same terminal verdict after finalize.

## Inputs

The case is implementation-agnostic. Each scenario's input is a sequence of `(kind, result)` pairs:

- `kind` ∈ {`chain_entry`, `seal`, `rotation`} — the §10.29 input enumeration
- `result` ∈ {`pass`, `fail`, `n/a`} — the per-input PASS/FAIL outcome from the §7 inner verifier (rotation events have `n/a` because the state machine does not recompute a MAC or signature for them)

The per-input PASS/FAIL outcome is the §7 verifier's job (per-event MAC recompute, seal signature verify); §10.29 governs only the state-machine layer that consumes those outcomes. Case 022 abstracts the inner verifier and pins the outer state-machine behaviour directly.

## Six pinned scenarios

| Scenario | Description | Final streaming verdict | Terminal verdict |
|---|---|---|---|
| `happy_path_no_rotation` | Three chain entries (PASS) → one seal (PASS) | `4` | `0` PASS |
| `rotation_then_confirming_seal` | Two entries PASS → rotation → entry PASS at state `6` → confirming seal PASS → entry PASS | `4` | `0` PASS |
| `rotation_with_failing_seal` | Rotation → seal under new key FAILS → state `5` (sticky); subsequent entries stay at `5` | `5` | `3` chain-anomaly |
| `back_to_back_rotation` | Two rotation events with no intervening seal — §10.28 integrity-claim violation | `5` | `3` chain-anomaly |
| `unconfirmed_rotation_at_eof` | Rotation observed but stream ends at state `6` without a confirming seal — control-completeness failure | `6` | `3` chain-anomaly |
| `anomaly_stickiness` | Chain entry FAIL at step 1 → `5`; subsequent passes do not clear the verdict | `5` | `3` chain-anomaly |

The six scenarios collectively exercise every row of the §10.29 transition table and both directions of the finalize collapse.

## Expected outputs

`expected.json` carries:

1. The full §10.29 transition table (`transitions_table`) — one entry per `(from_state, kind, result) → to_state` row.
2. The finalize-collapse table (`finalize_collapse_table`) — `4 → 0`, `5 → 3`, `6 → 3`.
3. A `scenarios` array. Each scenario carries the input sequence, a step-by-step trace recording the verdict before and after each step, the final streaming verdict, and the terminal verdict after finalize.

A conforming implementation reads each scenario's input sequence, walks its own state machine, and asserts byte-equality with the published per-step trace and terminal verdict.

## Conformance behavior

A conforming implementation:

1. Starts every streaming session at verdict `4` per §10.29.
2. Transitions per the published table — including the **rotation confirmation** rule (state `6` returns to `4` only on a passing seal), the **back-to-back rotation** rule (a second rotation before a confirming seal transitions to `5`), and the **anomaly stickiness** rule (state `5` is a fixed point under any input that passes the kind check).
3. Rejects inputs whose `kind` is outside the three-element enumeration with `procedure-could-not-begin` (exit code `1`); the kind check runs **before** the verdict-state dispatch (per §10.29 normative).
4. Collapses the streaming verdict to a terminal verdict at end-of-stream per the finalize-collapse table.
5. Refuses to accept further inputs after a terminal verdict has been emitted (verifier-state immutability per §10.29).

## Negative cases this fixture supports

- An implementation that allows state `5` to clear (e.g., resets to `4` on a passing seal) fails the `anomaly_stickiness` scenario.
- An implementation that re-arms state `6` on a back-to-back rotation (rather than transitioning to `5`) fails the `back_to_back_rotation` scenario.
- An implementation that collapses state `6` at EOF to `0` PASS (rather than `3` chain-anomaly) fails the `unconfirmed_rotation_at_eof` scenario.
- An implementation that accepts a fourth input kind (e.g., a custom `metadata` kind) without rejecting per the kind check is non-conformant and fails the §10.29 conformance bar in general (case 022 does not exercise this directly — the kind check is structural, separate from the transitions this case pins).

## Cross-references

- Spec §10.29 streaming-mode verifier procedure (the section this case pins)
- Spec §10.27 configurable cadence (the precondition for streaming-mode operation)
- Spec §10.28 streaming-mode IKM rotation discipline (defines what `rotation` events mean)
- Spec §10.12 verifier exit-code contract (the `4`/`5`/`6` extension)
- Case 021 (`021-rotation-completed-event-payload/`) — the byte form of the rotation events case 022 abstracts as `(kind="rotation", result="n/a")`

## Reproduction

The `_compute.py` script in this directory regenerates `expected.json` from the pinned transition table and scenario list. It uses only the Python standard library; running it walks each scenario through the same state machine a verifier would and emits the per-step trace.

```
python _compute.py
```
