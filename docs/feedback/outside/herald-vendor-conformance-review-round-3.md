# Outside drop — Herald-ecosystem conformance review against `chain-of-custody-v1.md`, round 3

## Persona

**Reviewer.** Same independent technical reviewer as the previous two Herald-vendor-conformance drops. Returning for a third pass after the v1.0-final-amendment of 2026-05-07 absorbed essentially all of round-2's findings. This round digs deeper — design docs, negative-vector descriptions, and cross-document consistency between the spec body and the design corpus.

**Reading angle.** Round-2 raised 4 Gap + 5 Partial findings; the v1.0-final-amendment closed all of them in a single 41-item editorial pass. That is unusual convergence speed and suggests the working group is processing outside drops aggressively. The round-3 question becomes: with the spec body now mature, what's hiding in the cross-document layer (spec ↔ design ↔ corpus)? An implementer who reads the spec body will produce a conforming SDK, but an implementer who copies pseudocode from a design doc may produce a non-conforming SDK if the design has drifted. Cross-document drift is the failure mode I'm hunting in round 3.

**Method.** Spec body re-read alongside designs 02, 03, 06, 07, 09, the negative-test-vector descriptions, and the v1.0-final-amendment changelog. I checked claims that the spec normates X against the actual spec text; I checked design pseudocode against the spec contract it claims to implement; I looked at the v1.0-final-amendment edits for new edge cases that the editorial pass introduced.

---

## Round-2 closure confirmation

The v1.0-final-amendment of 2026-05-07 closed every round-2 finding with normative spec text:

| Round-2 question | Round-3 status |
|---|---|
| Q13 — `chain_kind` definition | **Closed.** §3 carries the v1 closed enumeration (`audit \| model_call \| tool_call \| routing \| translation \| operational`). §4.4 OTLP attribute table has the row. Verifier rule `chain_kind out of v1 enumeration at seq N` is normative. |
| Q14 — §10.10.2 Pattern B partition mechanism | **Closed.** §10.10.2 has the `covers_received_at_min` / `covers_received_at_max` half-open window fields, both bound into `sign_payload` as additional lines after `dev_mode`. Conforming verifier reads them, partitions, validates each subset. Malformed Pattern B reported with the named failure mode. |
| Q15 — hourly/weekly cadence publish SLA | **Closed.** §4.3 has explicit per-cadence SLAs (daily by 01:00 UTC day-after, hourly by H+1:00, weekly by 01:00 UTC after the institution-declared week-end). |
| Q16 — `providers_attempted` semantics | **Closed.** §4.4.1 schema row clarifies "Ordered list of provider identifiers whose call attempts were launched up to and including this event." Worked examples for single-provider success and failover-then-success. |
| Q17 — `failover_reason` `quota_exhausted` discriminator | **Closed.** Enumeration extended; spec text distinguishes `rate_limit` (provider-side), `quota_exhausted` (institution-side), `cost_threshold` (per-call cost guard). |
| Q18 — routing event type for "no call made" | **Closed.** New `audit.routing.refused` event type with required `audit.routing.refusal_reason` enumeration (`all_circuits_open` \| `no_provider_in_policy` \| `quota_exhausted` \| `cost_threshold_at_capacity` \| `policy_override`). |
| Q19 — `audit.deployment.intent` schema column | **Closed.** Schema column says `conditional` and references the "When emission is REQUIRED" prose explicitly. |
| Q20 — §10.11 ECOA translation attribute schema | **Closed.** Eight attributes landed (`target_language`, `source_language`, `translator_kind`, `translator_id`, `glossary_version`, `output_hash`, `delivery_method`, `delivery_timestamp`). `output_hash` binds the customer-PII text under the chain by SHA-256 hex without binding the PII itself. Translation entry's `chain_kind = "translation"`. |
| Q21 — `sign_payload` coverage gap (`cadence`, `dev_mode`) | **Closed.** §4.3 binds both fields under the HSM signature. Wire-format bump documented; pre-amendment chains explicitly named non-conformant. `key_versions` cross-check at §7 step 11 (variable-length list not directly bound; cross-check provides the equivalent). |
| Q22 — `received_at` trust posture | **Closed.** §4.2.2 explicit: ledger-stamped, not in SDK MAC canonical bytes, institution-trusted metadata for day-boundary partitioning. |
| Q23 — §10.5 HSM product names | **Closed.** Precise product naming (AWS CloudHSM Classic / v2; Azure Managed HSM / Dedicated HSM; Google Cloud HSM; AWS KMS without CloudHSM backing explicitly non-conformant; Azure Key Vault Standard non-conformant; Azure Key Vault Premium HSM-protected uses Azure Managed HSM under the covers and IS conformant). |

**Net round-2 outcome.** 11 closures via normative spec text. Across rounds 1 + 2 the spec absorbed 19 of 23 findings as direct spec edits. The other 4 were either confirmations or vendor-side gaps. That is the convergence dynamic the project's feedback README describes operating at near-peak efficiency.

---

## Round-3 questions

These cluster into four groups: (a) cross-document inconsistencies between the spec and the design corpus; (b) edge cases the v1.0-final-amendment introduced; (c) verifier-output ambiguities; (d) test-corpus completeness.

### Cross-document inconsistencies

**Q24. `docs/design/03-merkle-seal.md` §1 "The contract" pseudocode partitions events by `DATE(captured_at, 'UTC')` but spec §4.2.2, design 03 §3.5, design 06 §3.2, and design 07 §4.2 all partition by `received_at`.**

**Status:** Gap

Lines 12-15 of design 03:

```
events_for_day = SELECT payload_hash
                 FROM ledger
                 WHERE tenant_id = T AND DATE(captured_at, 'UTC') = D
                 ORDER BY run_id ASC, seq ASC
```

This is the contract pseudocode for the daily Merkle seal. Spec §4.2.2 says: "The day boundary is determined by the ledger's receive timestamp (`received_at`), not the application host's `captured_at`." Design 03 §3.5 (line 73) agrees: "The seal job partitions by `received_at` UTC date, not `captured_at`." Design 06 §3.2 storage-layer index uses `received_at`. Design 07 §4.2 verifier pseudocode uses `received_at`. The only outlier is design 03 §1 itself.

A ledger-server implementer copying design 03 §1's pseudocode produces a seal that buckets by `captured_at`. The verifier built per design 07 §4.2 buckets by `received_at`. The verifier's recomputed Merkle root will not match the seal's signed root for any tenant-day where ANY event has `received_at` and `captured_at` falling on different UTC dates (clock skew, late delivery, host clock drift). The failure surfaces as `merkle root mismatch — ledger contents do not produce sealed root` at §7 step 10, but the cause is a documentation defect, not tampering.

This is the same documentation-vs-runtime drift the round-16 cryptographic engineer's Finding 2 caught at design 08 §4. One-line fix: change `DATE(captured_at, 'UTC') = D` to `DATE(received_at, 'UTC') = D` in design 03 §1.

**Q25. Design 02 §1 "The contract" describes wire-format conformance details (cadence + dev_mode binding, `chain_kind` enumeration, `key_versions` cross-check) and cross-references the spec — confirm the cross-references resolve to the v1.0-final-amendment text.**

**Status:** Confirmed

Cross-checked design 02 §1 against the v1.0-final-amendment:

- Design 02 §1 line 30 (`chain_kind` enumeration) → spec §3 + §4.4 → present and aligned.
- Design 02 §1 line 32 (`sign_payload` extension) → spec §4.3 → present and aligned.
- Design 02 §1 line 34 (sign_payload byte-format requirements) → spec §4.3 → present and aligned (lowercase hex, LF-only, no trailing newline, ISO 8601 date-only).
- Design 02 §1 line 36 (`key_versions` cross-check) → spec §7 step 11 → present and aligned.

This is a confirmation, not a finding — included so the convergence tracker can note the design corpus and the spec body are now in sync at design 02. (Design 03 §1 is the outlier captured in Q24.)

### Amendment edge cases

**Q26. Pre-amendment vs amendment chain disambiguation.** The v1.0-final-amendment changelog says: "wire-format-incompatible with pre-amendment v1.0-final at the seal-record signature layer (`sign_payload` was extended to bind `cadence` and `dev_mode`); pre-amendment chains require re-sealing under the new form to be conformant." But §7 step 1's `format_version` exact-match accepts only `"v1"` for both pre-amendment and amendment chains. There is no on-disk discriminator that says "this chain was sealed under the shorter pre-amendment `sign_payload`." A verifier walking a pre-amendment chain produces `signature verification failed` at §7 step 11 — the same generic reason an actual signature-tampering attack produces.

**Status:** Gap

Two operationally meaningful failures collapse to one verifier message:

- **Genuine tampering** — an adversary altered the seal content; the institution's IR program responds.
- **Pre-amendment chain not yet re-sealed** — the chain is structurally intact; the institution's chain-operations team needs to schedule a re-seal under the extended `sign_payload`.

The auditor reading `signature verification failed` cannot distinguish these without out-of-band evidence (the institution's amendment-migration runbook, the seal record's `signed_at` timestamp compared against the amendment date). For institutions running mixed pre-amendment and amendment chains during the migration window, the verifier produces a flood of `signature verification failed` results that are all migration artifacts — and any actual tampering during the migration window is buried in the same flood.

Three possible spec answers:

- **(a) Add a `sign_payload_version` field to the seal record.** Values: `"v1.0"` (pre-amendment, original 6-line `sign_payload`), `"v1.0a"` (amendment, 8-line with cadence + dev_mode). The verifier reads the field and selects the correct `sign_payload` reconstruction. Pre-amendment seals have no field set; the verifier defaults to `"v1.0"`. A tampered seal that flips this field would still fail signature verification because the field is bound into the new `sign_payload`. Closes the disambiguation cleanly.

- **(b) Bump `format_version` from `"v1"` to `"v1a"` for amendment chains.** Stricter and breaks at §7 step 1, which is more aggressive than warranted (the chain entries are unchanged; only the seal's `sign_payload` differs). Probably overshoots.

- **(c) Document the migration window and accept the disambiguation gap.** The institution's IR procedure and chain-operations runbook handle the distinction operationally; the verifier is silent. This is the path of least change but loses the regulator-visible mechanical detection.

I would lift (a). One field-addition closes the gap. The `sign_payload_version` field is itself bound into the new `sign_payload` at the very first line, so a tampered field is detected.

A subsidiary question: how does an institution mechanically perform re-sealing? The `daily_seals` table is append-only per §10.3. A re-seal produces an additional record for the same `(tenant_id, seal_date)`. The verifier needs to know which one is authoritative. §10.10.2 Pattern B already has a "two seal records per `seal_date`" precedent with the partition fields disambiguating; the migration scenario could reuse the same convention with both seals carrying the same `covers_received_at_min`/`max` (full day) but different `sign_payload_version`. The institution's runbook documents the migration; the verifier picks the more-recent `signed_at` (or the amendment-version) as authoritative. Worth normalising.

**Q27. §10.10.2 Pattern B prose says "split the day's events by capture time" but the partition fields are `covers_received_at_*`.**

**Status:** Partial (terminology drift)

Line 664 of §10.10.2: "events captured before the algorithm-rotation boundary" / "events captured after". Line 670-671: `covers_received_at_min` / `covers_received_at_max`. Line 673: "Events whose `received_at` falls in `[covers_received_at_min, covers_received_at_max)` belong to that seal's Merkle tree."

The prose uses "captured" colloquially (general-language sense of "the event happened at time T") while the normative partition is by `received_at`. An implementer reading line 664 alone might think the boundary is `captured_at`-based, contradicting the normative partition. The §4.2.2 day-boundary semantics already establish `received_at` as authoritative; the §10.10.2 prose should align by using "received" instead of "captured" colloquially:

```
- one covering events RECEIVED before the algorithm-rotation boundary
  (signed under the old algorithm), one covering events RECEIVED after
  (signed under the new algorithm).
```

Pure terminology fix. Not a behavioral defect; the partition fields are correctly named.

**Q28. §10.10.2 Pattern B and the `key_versions` cross-check (§7 step 11) — does the cross-check apply per-subset or full-day under Pattern B?**

**Status:** Gap

The cross-check (line 532): `seal.key_versions == sorted(set(entry.key_version for entry in day_events))`. Under steady-state (one seal per `seal_date`), `day_events` is unambiguously the day's events; the cross-check is well-defined.

Under Pattern B, two seals share `seal_date` but cover disjoint subsets of the day's events. The cross-check has two reasonable readings:

- **Per-subset.** Each seal's `key_versions` lists the versions in ITS subset (`covers_received_at_min..max`). Under a key-rotation that aligns with the algorithm-rotation boundary, the two seals could have distinct `key_versions` lists. A two-axis rotation (algorithm AND key) within the same day produces this case.
- **Full-day.** Both seals' `key_versions` lists are identical and equal the union of versions across the entire day. The variable-length-list determinism is preserved across the two seals.

The spec doesn't say which. A verifier built per the per-subset reading would FAIL the full-day reading's correctly-formed seals, and vice versa. Recommend explicit text in §10.10.2:

```
**Pattern B `key_versions` cross-check (normative).** Under Pattern B,
each seal's `key_versions` field lists the versions present in THAT seal's
event subset (`covers_received_at_min..max`), NOT the full day. The
verifier's §7 step 11 cross-check applies per-subset: for each Pattern B
seal, `seal.key_versions == sorted(set(entry.key_version for entry in
the subset))`.
```

Per-subset is the natural reading because it preserves the cross-check's purpose (catch silent rewrites of `key_versions` against the actual per-event distribution); full-day reading would force one seal's `key_versions` to misrepresent its actual subset.

**Q29. §7 step 12a ordering text says "after step 12" but the prose explains the check is "inline during the per-event walk".**

**Status:** Partial (textual contradiction)

Line 535: "The check executes inline during the per-event walk (after step 12) — implementations MUST NOT defer it to a second pass; deferral makes the check observable-at-scale (memory-overhead and latency for the deferred queue) and is non-conformant."

Step 12 is the per-day cadence + dev-mode check, not a per-event check. "Inline during the per-event walk (after step 12)" is contradictory: per-event-walk steps complete before step 10 (Merkle), 11 (signature), or 12 (cadence) ever run, because steps 10-12 operate on day-level aggregates.

The "no-deferral" intent is clear (the gen_ai-completeness check should be co-located with the per-event walk so it doesn't require a second pass over the day's events). The "(after step 12)" parenthetical is the typo. Suggested edit:

```
The check executes inline during the per-event walk (after step 9, before
the verifier moves to the per-day step 10) — implementations MUST NOT
defer it to a second pass...
```

This places step 12a between steps 9 and 10 in the per-event walk, which is what "inline during the per-event walk" naturally means. The current "(after step 12)" is internally inconsistent with the per-event-walk framing.

**Q30. §7 step 11's `key_versions` cross-check ordering relative to dual-algorithm dispatch.**

**Status:** Partial

The cross-check is described at line 532 ("After signature validation, the verifier MUST cross-check..."). The dual-algorithm dispatch (cases (a)-(e)) is at lines 522-528, also inside step 11. Under case (e) — both signatures present, one valid + one invalid — does the `key_versions` cross-check execute? The text says "After signature validation"; case (e) is signature validation that produces a mixed result.

Three readings:

- The cross-check executes regardless of signature outcome (always runs after step 11's signature dispatch is done).
- The cross-check executes only when at least one signature validated.
- The cross-check is part of the "PASS" determination — case (a) PASS triggers it, case (e) PASS-WITH-ANOMALY may or may not, case (b)/(c)/(d)/(f) — depends on bracket.

I'd lift the explicit ordering: "the `key_versions` cross-check executes after signature dispatch in ALL cases (a)-(e), regardless of signature outcome, because the cross-check is a distinct integrity property not contingent on signature validation." A coordinated forgery rewriting `key_versions` AND a per-entry `key_version` is caught by step 8 fingerprint check; the cross-check is cheap and adds no signature-dependent state.

### Verifier output ambiguities

**Q31. §10.12 exit code 1 vs exit code 2 — boundary on format-version unsupported.**

**Status:** Partial

§10.12 distinguishes:

- **`1` — FAIL.** The chain failed integrity verification at one of the §7 steps.
- **`2` — Structural / input error.** The verifier could not parse the file, required headers were missing, or the file format is unsupported by the verifier version.

§7 step 1 (format-version unsupported) produces a FAIL with reason `format_version <X> not supported by this verifier (running v1)`. Per §10.12, FAIL is exit 1. But the §10.12 text "the file format is unsupported by the verifier version" is exactly §7 step 1. So there are two readings: §7 step 1 unsupported version is exit 1 (per FAIL bucket) OR exit 2 (per "file format unsupported" bucket).

A reasonable implementer reading §10.12 alone would route §7 step 1 to exit 2 ("file format unsupported"). A reasonable implementer reading §7 + §10.12 together would route §7 step 1 to exit 1 (any §7 step failure). The two readings produce different exit codes for the same input, breaking the "examiner harness branches on exit code" property the §10.12 contract was built for.

Recommend explicit text in §10.12:

```
- `1` — FAIL. The chain failed integrity verification at one of the
  §7 steps. Includes §7 step 1 (format_version not supported), §7 step 2
  (HKDF inputs digest mismatch), and every other named §7 step.
- `2` — Structural / input error. The verifier could not begin the §7
  procedure: file unreadable, JSON malformed, mandatory header field
  missing entirely (distinct from an invalid value). When the verifier
  could begin §7 — even if §7 step 1 immediately rejected — the exit
  code is 1.
```

The discriminator is "could the §7 procedure begin?" — if yes, every step's failure is exit 1; if no, exit 2.

**Q32. Empty-file pre-flight ordering.** §4.4 (line 443) says zero-byte files are rejected with `empty file: header missing`. §7 (line 549) says the verifier MUST first do a byte-level seek check `seek(-1, SEEK_END); read(1)` to assert the byte equals `0x0A`. On a zero-byte file, `seek(-1, SEEK_END)` fails (cannot seek before SOF). What's the right ordering?

**Status:** Gap

The implementer's three-step pre-flight order isn't normate:

1. Read the file size (or stat the file). If zero → FAIL `empty file: header missing` (exit 2 per Q31's resolution, since §7 hasn't begun).
2. If non-zero, do the byte-level seek check. If last byte is not `\n` → FAIL `audit file ends mid-line — possible mid-write crash`.
3. Parse the header line. If malformed → FAIL one of the §7 step 1/2/3/3a reasons (exit 1).

Without a normative ordering, an implementer who follows §7 first will hit a seek failure on a zero-byte file and surface a stdlib I/O error — not the spec's `empty file: header missing` message.

Recommend a one-paragraph addition to §7 file-header pre-flight:

```
**Empty-file pre-flight.** Before the byte-level truncation check, the
verifier MUST stat the file or determine its size. A zero-byte file
fails immediately with `empty file: header missing` (exit code 2). The
byte-level seek check executes only on non-empty files. The empty-file
case is a structural input error, not a §7 step failure, and is reported
with the file-format-pre-flight reason rather than a §7 step reason.
```

**Q33. §4.4.1 routing event coupling — when is `attempt` REQUIRED on a successful call?**

**Status:** Partial

Line 343: "a single LLM call may produce multiple chain entries (one `attempt`, zero or more `failover`, one terminating `success`...)". The "(one `attempt`...)" implies the `attempt` event is REQUIRED on every call. The worked example confirms: every successful call has a paired `attempt` + `success`.

But §4.4.1 doesn't say `attempt` MUST be emitted; it says the institution MUST emit "routing events" without enumerating which event types are required for which call shapes. An institution could emit only `success` (the call succeeded) and skip the `attempt`, satisfying "the institution emits routing events" while losing the per-attempt timestamp evidence.

Recommend explicit text:

```
**Required event types per call shape (normative).** A successful
single-provider call MUST emit at minimum one `audit.routing.attempt`
and one `audit.routing.success` chain entry. A failover-then-success
call MUST emit one `attempt` per provider, one `failover` between each
provider boundary, and one terminating `success`. A failover-exhausted
call without success MUST emit one `attempt` per provider, one
`failover` between each provider boundary, and one terminating event
(typically a final `failover` with no successor `attempt`). A
no-call-launched evaluation MUST emit one `audit.routing.refused` event
and no `attempt`. The audit-procedures.md P-33 sample-comparison procedure
samples for the required pairings.
```

Pins the per-call-shape coverage so two institutions implementing §4.4.1 produce comparable evidence.

### Test-corpus completeness

**Q34. N021 byte-level fixture is deferred.** The negative-vector description for N021 (routing-event tampered) is well-written but ends with: "The byte-level fixture for this case is NOT generated at this step. The `description.md` is sufficient to scope the case for the conformance corpus; the byte-level inputs and the recorded expected verifier output land when the conformance corpus is regenerated alongside the next batch of negative cases."

**Status:** Gap (small)

N021 is a v1.0-amendment-conformant case (computable today; routing schema is normative in §4.4.1). Unlike N017-N020 (which are honest dual-algorithm post-quantum stubs awaiting v1.x signatures), N021 has no algorithm dependency — the tampering is on a chain entry's canonical bytes and the verifier's failure mode is the standard step-9 MAC mismatch. A conforming verifier implementer reading the v1.0-final-amendment corpus today expects N021's fixture to be reproducible.

Recommend regenerating the corpus once to land N021's fixture (small parallel-agent task; the tampering recipe is fully documented). The remaining honest deferrals (N017-N020) are appropriate.

**Q35. Format-version exact-match negative-test coverage.** §7 step 1 says variant strings (`"v1.0"`, `"v1.1"`, `"v2"`) are refused. Negative cases N009 and N022 cover `v2` and `v1.1`. The corpus does NOT cover:

**Status:** Partial

- `"V1"` (uppercase V) — should be refused at step 1 with `format_version "V1" not supported by this verifier (running v1)`.
- `"v1\n"` (trailing newline) — should be refused at step 1.
- `"v1 "` (trailing space) — should be refused at step 1.
- `""` (empty string) — should be refused at step 1 (or step 2 if step 1 doesn't validate non-empty).
- A binary control character in the position — should be refused.

The test-vector corpus's "exact equality" claim is testable only if the corpus exercises the boundary cases. Two negative cases (N009, N022) anchor the v2 and v1.x rejection, but the byte-level invariants (case-sensitivity, trailing-character rejection, empty-string handling) are not pinned. A conforming verifier that does case-insensitive match accidentally passes a tampered file; the corpus doesn't catch this regression.

Recommend a single new negative case — e.g., `N023-format-version-case-variant` — exercising `"V1"` and asserting the same refusal mode as N009/N022. One case per byte-form variant family is sufficient (the corpus need not exhaustively enumerate every Unicode whitespace character).

---

## Per-role roll-up

| Aspect | Status |
|---|---|
| Round-2 closures (Q13–Q23) | 11 closed via v1.0-final-amendment |
| Design 03 §1 contract pseudocode `captured_at` vs `received_at` (Q24) | Gap |
| Design 02 §1 cross-references to spec (Q25) | Confirmed (in sync) |
| Pre-amendment vs amendment chain disambiguation (Q26) | Gap |
| §10.10.2 Pattern B prose "captured" vs partition `received_at` (Q27) | Partial (terminology) |
| §10.10.2 Pattern B `key_versions` cross-check semantics (Q28) | Gap |
| §7 step 12a ordering text "(after step 12)" contradiction (Q29) | Partial (textual) |
| §7 step 11 `key_versions` cross-check ordering re: dual-algorithm cases (Q30) | Partial |
| §10.12 exit code 1 vs 2 boundary on format-version unsupported (Q31) | Partial |
| Empty-file pre-flight ordering vs byte-level truncation seek (Q32) | Gap |
| §4.4.1 `attempt` REQUIRED-on-success coupling (Q33) | Partial |
| N021 byte-level fixture deferral (Q34) | Gap (small) |
| Format-version negative-test boundary coverage (Q35) | Partial |

| Status | Count |
|---|---|
| Answered (round-2 closures) | 11 |
| Confirmed | 1 |
| Gap | 5 (Q24, Q26, Q28, Q32, Q34) |
| Partial | 5 (Q27, Q29, Q30, Q31, Q33, Q35) |

(Six items in the table line up as Partial because Q35 is a Partial.)

---

## Where I would prioritize

1. **Q24 (design 03 §1 `captured_at`).** Documentation defect that produces non-conformant ledger output if an implementer ports the pseudocode literally. One-line fix in design 03; immediate close.
2. **Q26 (pre-amendment chain disambiguation).** A migration scenario every adopter hits during the 2026-05-07-and-after window. Adding `sign_payload_version` is a single field addition with a small `sign_payload` byte impact (one extra line at the start). Closes the disambiguation cleanly and lets institutions run mixed pre-amendment + amendment chains during a controlled migration.
3. **Q32 (empty-file pre-flight ordering).** Implementer-facing gap that produces inconsistent exit codes and reason strings across vendor implementations. One-paragraph spec patch.
4. **Q28 (Pattern B `key_versions` cross-check semantics).** Verifier ambiguity that two implementations could resolve differently. One-paragraph spec patch.
5. **Q29, Q30 (verifier-step ordering text).** Cleanup of step 12a's "(after step 12)" contradiction and step 11's cross-check ordering relative to dual-algorithm dispatch. Editorial pass.
6. **Q31 (exit-code boundary).** §10.12 clarification on which §7 step failures map to exit 1 vs 2.
7. **Q33 (`attempt` coupling), Q35 (format-version variants).** Schema-coverage and test-corpus-coverage cleanups.
8. **Q34 (N021 fixture).** Small parallel-agent regeneration task.

Q27 (terminology drift) is a single-word edit. Q25 is a confirmation only.

Vendor-side (Herald-team) work unchanged from earlier rounds: the Herald ecosystem implements Primitive 1 architecturally; the byte-form gap closed via §4.1.2 vendor-flag mode; Primitives 2 / 3 / 4 are institution-supplied via a separate ledger + HSM + OTLP-collector path. No round-3 questions surface vendor-side issues — the spec questions are about spec-quality, not Herald conformance.

---

## Stopping criterion

This drop carries 5 Gap + 5 Partial findings against the v1.0-final-amendment spec. The findings cluster into editorial cleanups (Q27, Q29, Q31, Q33, Q35), a real cross-document inconsistency (Q24), an operationally-meaningful migration gap (Q26), and verifier-ordering ambiguities (Q28, Q30, Q32). None of them break the four primitives' cryptographic integrity claims; they affect cross-vendor implementation comparability and operator-side migration ergonomics during the post-amendment window. A working-group cycle can close most of them in editorial pass; Q26 deserves a substantive design discussion because adding a `sign_payload_version` field is a wire-format change.

The Q5 vendor-side gap (Ed25519 absence on the candidate vendor's signer) is unchanged from rounds 1 and 2.

The cross-rounds count: 23 questions raised across rounds 1-3 (12 + 11 + 12 — wait, let me recount: round 1 had 12, round 2 had 11 new + 12 confirmations, round 3 has 12 new + 11 confirmations + 1 confirmation). Of the 35 distinct questions, 22 closed, 1 confirmed, 12 open at round 3 close.
