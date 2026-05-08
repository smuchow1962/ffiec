# Case 017 — Merkle inclusion proof for partial disclosure

## Status

**Structural fixture.** Hash placeholders pinned with stable seeds; byte-level hash recomputation deferred to the verifier-toolchain landing the partial-disclosure mode (design 07 §10). The structural recipe and expected verifier behavior are normative as of this case.

## Purpose

Verify the partial-disclosure verifier mode per design 07 §10 against a 7-entry seal-day. The case exercises:

- The RFC 6962 §2.1.1 audit-path shape with directional-bit encoding.
- The right-promote balancing rule from case 016 in an audit-path context — the 7-leaf tree right-promotes the rightmost leaf at the second level, which an audit path for any leaf in the upper half of the tree must encode correctly.
- The per-entry verification procedure in design 07 §10.6 (steps P1-P5).
- The output shape design 07 §10.3 specifies, including the explicit `completeness_assertion: NOT-ASSERTED` field.

This case builds on case 016. A v1 verifier that fails 016 will fail 017 because the audit-path walk depends on the same right-promote rule; a verifier that passes 016 should be able to extend its Merkle code to handle audit paths without algorithmic change.

## Why this matters

A 2703(d) federal selective production, a SOC 2 sampling-driven audit, and a DORA Article 11-14 incident-reporting filing each produce one or more entries from a sealed daily set without producing the remainder. The partial-disclosure verifier authenticates each produced entry against the seal — an inclusion claim, not a completeness claim. The verifier's correctness depends on:

1. Recomputing the leaf hash from the entry's canonical bytes (step P1).
2. Walking the audit path to the sealed root, applying the right-promote rule when the proof carries a `"promote"` step (step P2).
3. Comparing the recomputed root to the seal's `merkle_root` (step P3).
4. Validating the seal's signature against the registry-anchored public key (step P4 — once per batch).
5. Confirming each entry's `tenant_id` matches the seal's `tenant_id` (step P5).

This case fixes the byte-level shape of each step's input and output for the index-3 entry of a 7-entry tenant-day.

## Sub-fixtures and tree shape

The 7-leaf tree is the case-016 7-leaf shape. Recall the right-promote pattern:

```
Level 0 (leaves):  L0   L1   L2   L3   L4   L5   L6
Level 1:           h(L0,L1)  h(L2,L3)  h(L4,L5)   L6'  ← L6 promoted unchanged
Level 2:           h(h(L0,L1), h(L2,L3))   h(h(L4,L5), L6')
Level 3 (root):    h(level 2 left, level 2 right)
```

The audit path for `L3` (index 3, the 4th entry from left) walks:

| Level | Sibling at this level | Position relative to running hash |
|---|---|---|
| 0 | `L2` | left (sibling on the left, running on the right) |
| 1 | `h(L0, L1)` | left |
| 2 | `h(h(L4,L5), L6')` | right (sibling on the right) |

Three audit-path steps. No `"promote"` step in this path — the right-promote happens in the right subtree, while `L3` lives in the left subtree. Other entries in the tree would produce paths with `"promote"` steps; `L3` produces a clean three-step path.

A separate path for `L6` (the right-promoted leaf) would carry a `"promote"` step at level 1 — the rightmost leaf has no sibling at level 1, so the running hash is carried forward unchanged. The case fixture pins the `L3` proof; an extended fixture covering `L6` is a v1.x extension that would exercise the `"promote"` step end-to-end.

## Inputs

```
tenant_id  = "tenant-ffiec-test-1"
tenant_day = "2026-05-08"
spec_version = "v1.0"
format_version = "v1"
sign_payload_version = "v1.0a"
algorithm  = "ed25519"
cadence    = "daily"
dev_mode   = false
```

The 7 entries' canonical bytes are stable seeds chosen so any conforming implementation can reproduce the leaf hashes without this fixture file:

```
entry_k_bytes = b"ffiec-017-entry-" || ascii(k)   for k in [0, 7)
leaf_k        = SHA-256(0x00 || entry_k_bytes)    // RFC 6962 leaf-domain prefix
```

The seal's `merkle_root` is the case-016 7-leaf root produced by these leaves under the same right-promote rule.

## Files

| File | Contents |
|---|---|
| `entries.json` | The 7 entries' canonical bytes plus their seq numbers. Hash fields marked `// PLACEHOLDER:hex(32)` for the verifier toolchain to populate. |
| `seal.json` | The seal record for `(tenant-ffiec-test-1, 2026-05-08)` over all 7 entries. Hash and signature fields marked as placeholders. |
| `partial-proof-entry-3.json` | The inclusion proof for entry index 3 (seq=4) — leaf hash, three-step audit path, expected root, signature. |
| `expected-result.json` | The verifier output shape for the partial-disclosure run, naming step P1-P5 outcomes per entry plus the aggregated summary. |

## Expected verifier outcome

```
Status: PASS (partial-disclosure)
Mode: partial-disclosure
tenant_id: tenant-ffiec-test-1
tenant_day: 2026-05-08
Entries disclosed: 1
Entries passed: 1
Entries failed: 0
Completeness assertion: NOT-ASSERTED
Per-entry results:
  - seq=4: PASS (steps P1, P2, P3, P5)
Per-batch result:
  - Seal signature: VALID (step P4)
Banner: This report does NOT attest chain completeness for tenant_day 2026-05-08.
```

The exit code is `0` per design 07 §10.9.

## Failure modes the case discriminates

A verifier that mishandles the partial-disclosure mode fails 017 in distinct ways:

- If the verifier walks the audit path with no domain-separator prefix on internal-node hashes (`SHA-256(left || right)` instead of `SHA-256(0x01 || left || right)`), the running hash at level 1 differs from the seal's recorded subtree hash and step P3 fails with `audit path does not lead to sealed root`.
- If the verifier mishandles the directional bit (treats `"left"` as `"right"`), the running hash at level 0 differs and step P3 fails for the same reason. The error is at the audit-path-walk layer; the leaf hash computed at step P1 is correct.
- If the verifier omits step P5 (the `tenant_id` cross-check), an entry sealed in tenant A presented as a member of tenant B's daily seal would falsely PASS. The case's expected output names step P5 explicitly; a verifier whose output omits step P5 from the per-entry breakdown is non-conformant.
- If the verifier emits a single combined PASS over the whole input rather than per-entry verdicts, the output shape diverges from `expected-result.json` and the case fails by output-shape comparison rather than by cryptographic mismatch.
- If the verifier omits the `completeness_assertion: NOT-ASSERTED` field, a downstream consumer could mistake the partial-disclosure result for a full-day pass. The output-shape comparison catches the omission.

## Conformance test

A v1 verifier supporting the partial-disclosure mode per design 07 §10 MUST:

1. Recompute `leaf_3 = SHA-256(0x00 || b"ffiec-017-entry-3")` and compare against the proof's `leaf_hash`.
2. Walk the three-step audit path with the directional bits and produce a running hash that matches `seal.merkle_root`.
3. Verify the seal signature against the registry-anchored public key (step P4 — case treats this as PASS by precondition since the registry snapshot is fixture-defined).
4. Confirm `entry.tenant_id == seal.tenant_id` for the disclosed entry.
5. Emit per-entry output in the shape `expected-result.json` pins, including the `completeness_assertion: NOT-ASSERTED` field and the partial-disclosure banner.

A verifier that produces all five outcomes from the inputs in `entries.json`, `seal.json`, and `partial-proof-entry-3.json` is conformant for the partial-disclosure mode.

## Cross-references

- Test-vector 016 — non-power-of-2 Merkle. The 7-leaf tree shape and the right-promote rule are pinned in 016; case 017 reuses them in an audit-path context.
- Design 07 §10 — partial-disclosure verifier mode. The §10.6 procedure (steps P1-P5) is what the verifier runs against this fixture.
- Design 07 §10.5 — audit-path shape with directional bits. The proof in `partial-proof-entry-3.json` follows this shape.
- `docs/selective-production-and-sampling.md` — the consumer-side procedural framing for 2703(d), SOC 2, and DORA contexts.
- Spec §5 — canonical bytes per entry; step P1 reads from this form.
- Spec §10.12 — verifier exit codes; design 07 §10.9 maps the partial mode onto them.

## Provenance

The fixture files are placeholder-shaped. Hash fields carry `// PLACEHOLDER:hex(32)` markers where a future verifier-toolchain run will compute the byte-level value from the leaf seeds and the right-promote rule. The structural shape of every file — field names, ordering, directional-bit encoding — is normative as of this case. Replacing the placeholders with computed hex from a conforming implementation produces the fixture in its final form.
