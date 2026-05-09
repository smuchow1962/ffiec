# -*- coding: utf-8 -*-
"""Compute the streaming-verifier state-machine trace for FFIEC v1
test-vector case 022-streaming-verifier-incremental.

Spec §10.29 normates the streaming-mode verifier state machine
covering verdict states `4` (streaming-all-pass-so-far), `5`
(streaming-anomaly-detected), and `6` (streaming-key-rotation-
pending-confirmation), plus end-of-stream collapse to terminal
verdicts `0` (PASS) or `3` (chain-anomaly).

This case pins a sequence of inputs and the expected per-step verdict
transitions through the state machine, plus the finalize collapse,
so a clean-room implementation proves it walks the §10.29 table
identically to other implementations. The case is deterministic and
implementation-agnostic — the inputs are abstract `kind`+`outcome`
pairs (the per-input MAC / signature recomputation is the §7
verifier's job; §10.29 governs only the state-machine layer).

Run with: python _compute.py
"""

from __future__ import annotations

import json
import os

HERE = os.path.dirname(os.path.abspath(__file__))


# State-machine transitions per spec §10.29 (normative table).
# Keys: (current_verdict, kind, per_input_result). Result is one of
# "pass", "fail", or "n/a" (rotation events have no per-input PASS/FAIL).
# Values: next_verdict.
TRANSITIONS: dict[tuple[int, str, str], int] = {
    # From state 4 (streaming-all-pass-so-far)
    (4, "chain_entry", "pass"): 4,
    (4, "chain_entry", "fail"): 5,
    (4, "seal", "pass"): 4,
    (4, "seal", "fail"): 5,
    (4, "rotation", "n/a"): 6,
    # From state 6 (rotation pending)
    (6, "seal", "pass"): 4,    # confirms the rotation
    (6, "seal", "fail"): 5,    # rotation seal fails
    (6, "chain_entry", "pass"): 6,
    (6, "chain_entry", "fail"): 5,
    (6, "rotation", "n/a"): 5, # back-to-back rotation = anomaly
    # State 5 is sticky for any input that passes the kind check.
    (5, "chain_entry", "pass"): 5,
    (5, "chain_entry", "fail"): 5,
    (5, "seal", "pass"): 5,
    (5, "seal", "fail"): 5,
    (5, "rotation", "n/a"): 5,
}

# Finalize collapse: end-of-stream streaming verdict → terminal verdict.
# Per spec §10.29 finalize-collapse table.
FINALIZE_COLLAPSE: dict[int, int] = {
    4: 0,  # all-pass-so-far → PASS
    5: 3,  # anomaly-detected → chain-anomaly
    6: 3,  # rotation-pending-confirmation never confirmed → chain-anomaly
}


def step(current: int, kind: str, result: str) -> int:
    key = (current, kind, result)
    if key not in TRANSITIONS:
        raise ValueError(f"undefined transition: state={current} kind={kind!r} result={result!r}")
    return TRANSITIONS[key]


def trace(inputs: list[dict]) -> list[dict]:
    """Walk the state machine through `inputs` and emit a step-by-step trace.

    Each input is a dict {kind, result}. The trace records the verdict
    BEFORE and AFTER each step so a reader can spot the transition.
    """
    out: list[dict] = []
    state = 4  # streaming verifier starts at 4 per §10.29
    for i, inp in enumerate(inputs):
        before = state
        after = step(state, inp["kind"], inp["result"])
        out.append({
            "step": i + 1,
            "input": inp,
            "verdict_before": before,
            "verdict_after": after,
        })
        state = after
    return out


# Pinned input scenarios — one fixture exercises each major state-
# machine path so a clean-room implementation can verify all the rows
# of the §10.29 table.
SCENARIOS: list[dict] = [
    {
        "scenario_id": "happy_path_no_rotation",
        "description": (
            "Three chain entries (all PASS) followed by a seal (PASS). "
            "Verdict stays at 4 throughout; finalize collapses to 0 (PASS)."
        ),
        "inputs": [
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "seal", "result": "pass"},
        ],
    },
    {
        "scenario_id": "rotation_then_confirming_seal",
        "description": (
            "Two chain entries PASS → rotation event → one chain entry PASS "
            "(stays at 6) → confirming seal PASS (drops back to 4) → one more "
            "entry PASS. Finalize collapses to 0 (PASS)."
        ),
        "inputs": [
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "rotation", "result": "n/a"},
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "seal", "result": "pass"},
            {"kind": "chain_entry", "result": "pass"},
        ],
    },
    {
        "scenario_id": "rotation_with_failing_seal",
        "description": (
            "Rotation event → seal under the new key FAILS → 5 (sticky). "
            "Subsequent entries stay at 5. Finalize collapses to 3."
        ),
        "inputs": [
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "rotation", "result": "n/a"},
            {"kind": "seal", "result": "fail"},
            {"kind": "chain_entry", "result": "pass"},
        ],
    },
    {
        "scenario_id": "back_to_back_rotation",
        "description": (
            "Rotation event → another rotation before any confirming seal → 5. "
            "§10.28 integrity-claim violation: two rotations without an "
            "intervening seal under the new key."
        ),
        "inputs": [
            {"kind": "rotation", "result": "n/a"},
            {"kind": "rotation", "result": "n/a"},
        ],
    },
    {
        "scenario_id": "unconfirmed_rotation_at_eof",
        "description": (
            "Rotation event → chain entries PASS but stream ends at state 6 "
            "without a confirming seal. Finalize collapses to 3 (control-"
            "completeness failure)."
        ),
        "inputs": [
            {"kind": "chain_entry", "result": "pass"},
            {"kind": "rotation", "result": "n/a"},
            {"kind": "chain_entry", "result": "pass"},
        ],
    },
    {
        "scenario_id": "anomaly_stickiness",
        "description": (
            "Chain entry FAILS at step 1 → 5. Subsequent passes do not "
            "clear the verdict per §10.29 stickiness rule."
        ),
        "inputs": [
            {"kind": "chain_entry", "result": "fail"},
            {"kind": "seal", "result": "pass"},
            {"kind": "chain_entry", "result": "pass"},
        ],
    },
]


def main() -> None:
    fixture = {
        "_about": (
            "Streaming-verifier state-machine pin per spec §10.29. Each "
            "scenario walks a sequence of inputs through the §10.29 "
            "transition table and records the verdict before and after "
            "each step, plus the end-of-stream finalize collapse to a "
            "terminal verdict (0 PASS or 3 chain-anomaly)."
        ),
        "transitions_table": [
            {"from": k[0], "kind": k[1], "result": k[2], "to": v}
            for k, v in TRANSITIONS.items()
        ],
        "finalize_collapse_table": [
            {"streaming_verdict": k, "terminal_verdict": v}
            for k, v in FINALIZE_COLLAPSE.items()
        ],
        "scenarios": [],
    }

    for scenario in SCENARIOS:
        steps = trace(scenario["inputs"])
        final_streaming_verdict = steps[-1]["verdict_after"]
        terminal_verdict = FINALIZE_COLLAPSE[final_streaming_verdict]
        fixture["scenarios"].append({
            "scenario_id": scenario["scenario_id"],
            "description": scenario["description"],
            "inputs": scenario["inputs"],
            "steps": steps,
            "final_streaming_verdict": final_streaming_verdict,
            "terminal_verdict_after_finalize": terminal_verdict,
        })

    with open(os.path.join(HERE, "expected.json"), "w", encoding="utf-8") as f:
        json.dump(fixture, f, indent=2, ensure_ascii=False)
        f.write("\n")

    for s in fixture["scenarios"]:
        print(
            f"[022] {s['scenario_id']:30s} "
            f"final={s['final_streaming_verdict']} "
            f"terminal={s['terminal_verdict_after_finalize']}"
        )


if __name__ == "__main__":
    main()
