# N021 — Routing event tampered with (`audit.routing.provider_chosen` altered after capture)

## Purpose

Verify that routing chain entries (per spec §4.4.1) are covered by the chain MAC the same way ordinary `audit.*` chain entries are. An attacker who alters a routing attribute after capture — for example, changing `audit.routing.provider_chosen` from `openai-gpt-4o` to `anthropic-claude-sonnet` to make the institution's audit trail suggest a different provider was selected at the moment of decision — MUST be detected at spec §7 step 9 (`payload_hash MAC mismatch`). The routing attributes are part of the canonical bytes the chain MAC covers; tampering with them invalidates the MAC the same way tampering with any other event attribute does.

A verifier that passes a chain whose routing attributes were altered post-capture is non-conforming. The case proves the chain's integrity property extends to the routing-decision evidence the spec normates in §4.4.1, not just to the LLM-call evidence and the broader `audit.*` namespace.

## Tampering recipe

Start from a passing-case input that includes a routing chain entry — typically a small fixture demonstrating the spec §4.4.1 shape: one `audit.routing.attempt` entry naming `openai-gpt-4o` as `provider_chosen`, followed by an `audit.routing.success` entry, followed by an LLM-call entry whose `gen_ai.response.model` resolves to the OpenAI provider, all linked via `parent_run_id` / `parent_seq` per spec §4.4.

Tamper with the routing entry's `audit.routing.provider_chosen` attribute: change the value from `"openai-gpt-4o"` to `"anthropic-claude-sonnet"` while leaving the entry's persisted `payload_hash` field untouched. Do NOT recompute the MAC — the attacker does not have the IKM, and even if they did, the verifier's `expected_prev_hash` machinery (spec §4.1 inviolate property #8) closes the recompute path against substituted-prev-hash attacks. The tampering simulation is deliberately the simplest form: alter the attribute byte payload, leave the persisted MAC alone, present the result to the verifier.

Variants the same case covers (any one is sufficient to exercise the failure-mode):

- Change `audit.routing.event_type` from `"audit.routing.attempt"` to `"audit.routing.success"` (collapses two distinct event types into one, masking that the routing decision was the initial attempt versus the terminating success).
- Change `audit.routing.failover_reason` from `"timeout"` to `"cost_threshold"` on a `failover` entry (rewrites the institution's narrative about why it moved away from a provider — useful for an adversary trying to make a vendor-relationship issue look like a cost-policy decision).
- Change `audit.routing.policy_version` from `"v3.2"` to `"v3.5"` (rewrites the policy-version attribution so the routing decision appears to have been made under a different policy than the one actually in force at the time).
- Change `audit.routing.providers_attempted` from `["openai-gpt-4o"]` to `["openai-gpt-4o", "anthropic-claude-sonnet"]` (fabricates a failover that did not occur).

All variants land at the same verifier outcome: the canonical bytes covering the routing attributes have changed, the recomputed HMAC over the tampered canonical bytes does not match the persisted `payload_hash`, and the verifier reports the MAC mismatch at step 9.

## Expected verifier outcome

```
Status: FAIL
Step:   9
Reason: payload_hash MAC mismatch at seq N
```

Where `seq N` is the seq of the tampered routing entry within the run.

The verifier MUST report the mismatch at step 9 (the per-event MAC check). It MUST NOT report any routing-specific reason — the verifier does not interpret the `audit.routing.*` attribute schema; the verifier's integrity claim covers the bytes of the routing decision the same way it covers any other event.

A verifier that passes the tampered chain (returns `ok=True`) is non-conforming. A verifier that reports a different step or a different reason (e.g. `chain link broken at seq N+1` because the verifier accidentally treated the routing entry's `payload_hash` as the source of the next entry's `prev_hash` and the linkage walked correctly even though the routing entry's MAC was wrong) is also non-conforming — the routing entry's MAC failure surfaces at step 9 BEFORE the next entry's structural walk, per spec §7's ordered procedure.

## What this case proves

The case proves three load-bearing properties of spec §4.4.1's chain integration:

1. **Routing attributes are inside the canonical bytes.** The chain MAC covers the `audit.routing.*` attributes the same way it covers any other event attribute. The §4.4.1 attribute schema is institution-emitted content but the integrity coverage is the chain's standard MAC coverage; nothing about routing entries is structurally different from ordinary `audit.*` entries.

2. **The verifier does not need routing-specific logic.** The verifier walks routing entries per spec §7 without interpreting the routing-attribute schema. Tampering with routing attributes surfaces as the standard MAC mismatch at step 9 — the same failure mode the verifier produces for tampering with any other attribute.

3. **The §4.4.1 integrity claim composes cleanly with the existing §7 verification procedure.** A verifier built against v1.0-pre-§4.4.1 (no awareness of routing entries as a distinct concept) walks routing entries correctly because they look like any other chain entry. A verifier built against v1.0 with §4.4.1 awareness has nothing extra to do for integrity — the routing-attribute schema is consumed by the audit-procedure layer (P-33) and the IR program (Scenario 15), not by the verifier.

The auditor reading the verifier output for a chain containing tampered routing entries sees `payload_hash MAC mismatch at seq N` and routes the investigation to spec §7 step 9's standard failure mode. The auditor reads the affected entry's content and observes that the entry is a routing entry; the investigation discipline is the same as for any other tampered entry — confirm whether the cause is SDK defect, network corruption, or tampering, and respond per the institution's IR program.

## Evidence the tampering was detected

The verifier's output for the tampered chain contains:

- The named failure step (9) and reason string (`payload_hash MAC mismatch at seq N`).
- The affected `(tenant_id, run_id, seq)` tuple identifying the tampered routing entry.
- The entry's recorded `payload_hash` (the persisted MAC the attacker did not recompute).
- The verifier's recomputed MAC over the tampered canonical bytes (which differs from the recorded `payload_hash`).
- The constant-time comparison result confirming the two MACs disagree.

The institution's SOC team and the FFIEC examiner read the verifier output and identify the tampering as a chain-integrity finding at step 9. The investigation routes through IR Scenario 1 (`docs/incident-response-playbook.md` Scenario 1 — chain hash mismatch) for the integrity-side response. Because the affected entry happens to be a routing entry, the institution's MRM committee is notified through the standard Scenario 1 post-incident review path; the routing-attribute tampering is documented in the institution's MRM evidence record alongside the chain-integrity disposition.

## Negative case provenance

This case fills the gap noted in AI platform engineer review (Maya Patel, F-1 cascade closure): spec §4.4.1 normates the routing-attribute schema and names the chain MAC as the integrity property covering it, but no negative case in the corpus exercised tampering with a routing attribute. The case demonstrates that the integrity claim spec §4.4.1 makes about routing entries is exactly the same integrity claim the verifier already enforces for every other chain entry — no new verifier logic is needed, and the existing step-9 failure mode is the load-bearing detection.

## Implementation note — fixture bytes deferred

The byte-level fixture for this case is NOT generated at this step. The `description.md` is sufficient to scope the case for the conformance corpus; the byte-level inputs and the recorded expected verifier output land when the conformance corpus is regenerated alongside the next batch of negative cases. The reference Go runner (per `docs/design/08-test-vectors.md` §5.5) produces the tampered inputs from a passing-case fixture that includes routing chain entries; the conformance-corpus regeneration adds that passing-case fixture and this case's tampered variant in the same regeneration pass.

Until the byte-level fixture lands, conforming verifier implementations may exercise the case in the implementation's own test suite by constructing a passing chain with a routing entry, applying any of the tampering variants above, and asserting the verifier's output matches the expected outcome.
