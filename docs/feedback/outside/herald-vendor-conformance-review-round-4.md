# Outside drop — Herald-ecosystem conformance review against `chain-of-custody-v1.md`, round 4

## Persona

**Reviewer.** Same independent technical reviewer as the previous three Herald-vendor-conformance drops. Returning for a fourth pass after the v1.0-final-amendment Wave-3 close-out absorbed essentially all of round-3's findings AND landed substantial user-directed scope additions (multi-region resilience, late-binding events, empty-day seal continuity, Merkle ordering normative).

**Reading angle.** Round-3 raised 5 Gap + 5 Partial findings; the Wave-3 close-out (alongside Diego Hernández's evidentiary work) closed all of them in one 36-item editorial pass. Round-4 looks at the new sections introduced in Wave-3 + the user-directed additions (§1.1 Daubert grounding, §1.2 epistemic scope, §5.2 best-evidence posture, §10.13 evidentiary artifacts, §10.14 trusted-time integration, §10.15 multi-region resilience) for new edge cases the editorial pass introduced.

**Method.** Re-read spec body cover-to-cover. Spot-checked the new `sign_payload_version` field's verifier dispatch text (§7 step 11). Walked the §10.15 Pattern A invariants against the threat model (§9 / `09-threat-model.md`). Checked the late-binding attribute's integrity posture against the §5 inclusion list. Did NOT re-derive byte values this round — Wave-3 regenerated the corpus and round-13 cryptographic-engineer's drop already pinned that work.

---

## Round-3 closure confirmation

The Wave-3 close-out absorbed every round-3 finding with normative spec text:

| Round-3 question | Round-4 status |
|---|---|
| Q24 — `docs/design/03-merkle-seal.md` §1 `captured_at` vs spec `received_at` | **Closed.** Design 03 §1 fix landed. |
| Q26 — pre-amendment vs amendment chain disambiguation | **Closed.** New `sign_payload_version` field on the seal record (`"v1.0a"` for amendment seals; absent for pre-amendment); §7 step 11 dispatches on the field with three explicit cases (absent → 6-line form, `"v1.0a"` → 10-line form, unrecognized → named refusal). The field is bound into the 10-line form's second line so a tampered value is detected at signature verification. This is exactly the (a) variant I recommended in round 3. Clean close. |
| Q27 — §10.10.2 prose "captured" vs partition `received_at` | **Closed.** Pattern B prose terminology fix. |
| Q28 — Pattern B `key_versions` cross-check semantics | **Closed.** Per-subset reading is normative. |
| Q29 — §7 step 12a "after step 12" contradiction | **Closed.** Per-event walk runs inline after step 9. |
| Q30 — `key_versions` cross-check ordering re: dual-algorithm | **Closed.** Cross-check runs after signature dispatch in ALL cases (a)–(e). |
| Q31 — §10.12 exit code 1 vs 2 boundary | **Closed.** "Could the §7 procedure begin?" is the discriminator. |
| Q32 — empty-file pre-flight ordering | **Closed.** Zero-byte file rejected before byte-level seek check. |
| Q33 — `attempt` REQUIRED-on-success coupling | **Closed.** §4.4.1 names required event types per call shape. |
| Q34 — N021 byte-level fixture | **Closed.** Fixture regenerated for the v1.0a wire form. |
| Q35 — format-version negative-test boundary coverage | **Closed.** New N023 format-version-case-variant case. |

**Net round-3 outcome.** 11 closures via normative spec text. Across rounds 1+2+3, the spec has absorbed 30 of 35 findings as direct edits. The other 5 were either confirmations or vendor-side gaps. Convergence pace is striking — three reviewer waves on the same calendar day (2026-05-07) produced a 130+-item editorial close-out.

Bonus credit: Diego Hernández's wave-3 evidentiary additions (§1.1 Daubert grounding, §1.2 epistemic scope, §5.2 best-evidence posture, §10.13 evidentiary artifacts, §10.14 trusted-time, full `litigation-support.md`) are exactly the kind of evidence-stance work an examiner-laptop / IT-witness path needs and the prior rounds did not surface. Substantive lift.

---

## Round-4 questions

These cluster into three groups: (a) the new §10.15 multi-region section's enforcement boundaries; (b) the new §10.14 trusted-time integration's integrity posture; (c) the §1.1 Daubert four-factor framing's residual-risk completeness.

### Multi-region resilience (§10.15)

**Q36. "Region" is referenced normatively but not defined.** §10.15 Pattern A invariants 1-6 use the term "region" repeatedly: "region-agnostic", "single seal region per tenant per `seal_date`", "the promoted region MUST have all events", "per-region event-count reconciliation". §3 (Definitions) does not define region. An institution's CC8.1 control description "names the chosen pattern and the operational discipline" but the spec body has no normative shape for what counts as a region.

**Status:** Gap

Three implementer-facing readings of "region":

- **Cloud provider region.** AWS `us-east-1` vs `eu-west-1`. The unit AWS / Azure / GCP exposes for failover topology. Most multi-region deployments mean this.
- **Datacenter / availability zone.** Sub-region facility identifier (AWS AZ, Azure availability zone). Tighter than provider region.
- **Institution-defined logical region.** A grouping the institution defines for its own resilience model — could be cross-cloud, could be on-prem-plus-cloud, could be regulator-jurisdiction-bound.

The three readings produce different operational postures and different §10.15 invariant 5 reconciliation populations. A spec definition would close this. Recommend §3 entry:

```
| **Region** | An institution-defined unit of operational and
  cryptographic locality for multi-region deployments. The
  institution's CC8.1 control description names what each region
  encompasses (e.g., "AWS region us-east-1 + the institution's
  on-prem East datacenter") and the per-region operational
  identity (the regional ledger endpoint, the regional HSM cluster
  membership, the regional public-key registry entry). The
  institution's risk-tolerance statement governs the granularity
  choice. |
```

This makes §10.15 Pattern A invariants enforceable against a documented institution-side definition rather than ambient industry vocabulary.

**Q37. Pattern A invariant 2 (run-locality) is enforceable at the SDK process boundary but not from the chain alone.** §10.15.2: "The SDK MUST NOT compute a chain entry whose `prev_hash` references an event captured in a different region." This is an SDK-internal invariant — the SDK process knows which region it's running in and refuses to compute chain entries whose prev_hash links to a different-region event. But the chain entries themselves do NOT carry a `region` attribute (no `ffiec.chain.region` in §4.4's table); the verifier walks the seal region's ledger and cannot tell, from the chain bytes alone, whether each event was captured in the seal region or replicated from a different region.

**Status:** Partial

The verifier's role here is documented as silent: "Cross-region replication-loss is detected by the institution's reconciliation, not by the verifier directly — the verifier's output is silent on replication completeness" (§10.15 Pattern A — verifier behavior). Same applies to run-locality: the verifier cannot check that all events in a run were captured in the same region.

That posture is internally consistent — the verifier walks chain integrity, the institution's CC8.1 procedures handle multi-region invariants — but two follow-up questions:

(a) Should the chain entries carry an OPTIONAL `ffiec.chain.region` attribute (institution-emitted, MAC-bound, like `audit.routing.*`) so a Pattern A failover-incident reconstruction can identify which events came from which region? Today the institution's regional-ledger storage carries that attribution; making it MAC-bound on the entry would make the attribution evidence rather than just operational metadata.

(b) §10.15 Pattern A invariant 2 says the SDK MUST refuse to chain across regions. Under what mechanism does an SDK detect that an event came from a different region? If a single SDK process serves multiple tenants from multiple regional sources, the SDK has to consult some institution-supplied region tag at event capture. If it's per-process (one SDK process = one region), the invariant is trivially enforced. The spec is silent on which model is conformant. Worth naming: "SDKs MUST be configured per-region; serving events from multiple regions in a single SDK process is non-conformant."

**Q38. Pattern A invariant 6 (seal-region failover) — the promoted region's "all events for the tenant-day" assertion is not verifier-detectable.** §10.15 invariant 6 says the institution promotes a replication region to seal-region status when the seal region is unavailable, and "The promoted region MUST have all events for the tenant-day before producing the seal." If the failover happens with incomplete replication, the produced seal seals an incomplete view of the day. The verifier's output is `Status: PASS` — the seal's Merkle root equals the recomputed root over the events the (promoted) seal-region has. The events that didn't replicate before failover are silently absent from the chain.

**Status:** Partial

The institution-side per-region reconciliation (invariant 5) catches this asynchronously: each region reports its event count, and the seal-region's count is compared against the sum. A failover with incomplete replication produces a count mismatch surfaced in the `master.cross_region_replication_completed` operational event. So the gap is detected — just not by the verifier, and not synchronously with the seal walk.

For high-stakes audit trails (litigation support, regulator inquiry), the gap is a documentation question: the verifier's working-paper output should link to the institution's reconciliation evidence so an examiner reviewing the verifier output knows where to find the replication-completeness evidence. Recommend a one-line addition to §10.15 Pattern A — verifier behavior:

```
The verifier's output for a Pattern A tenant references the
institution's `master.cross_region_replication_completed` operational
events for the seal-date as supporting evidence; the examiner
cross-checks the verifier's `Status: PASS` against the per-region
event counts the institution reported.
```

This makes the institution-side reconciliation a named link in the audit-evidence chain rather than ambient process knowledge.

### Trusted-time integration (§10.14)

**Q39. Trusted-time integration token placement is unspecified.** §10.14: "A future v1.x extension MAY define a normative `audit.timestamp.rfc3161_token` attribute (Base64-encoded RFC 3161 TimeStampToken for the entry's `captured_at` value) to provide an independent time-authority attestation alongside the institution's NTP-synchronized clock."

**Status:** Gap (forward-scope)

Two implementation models for the token:

- **Pre-MAC binding.** SDK requests the token from the TSA (RFC 3161) BEFORE computing the per-event MAC. The token is in the canonical bytes the MAC seals; tampering with the token is detected at §7 step 9. Cost: one TSA round-trip per chain entry on the hot path (network IO, TSA latency). For high-volume institutions, this may be unacceptable — a TSA round-trip turns a microsecond-scale MAC compute into a hundred-millisecond-scale attestation.

- **Post-MAC binding.** SDK computes the MAC, writes the chain entry, then asynchronously requests the TSA token and stores it adjacent to the chain entry (e.g., in a sidecar). The token is NOT in the canonical bytes; tampering with the token is detected only by independent TSA validation, not by the chain MAC. This decouples hot-path latency from TSA availability.

Both models compose with the chain, but they answer different evidentiary questions. Pre-MAC binding asserts "the SDK had a TSA-attested timestamp at the moment the MAC sealed the event"; post-MAC binding asserts "an institution-internal worker obtained a TSA-attested timestamp for this captured_at after capture." For litigation support (the §5.2 best-evidence framing), the distinction matters — pre-MAC binding gives the SDK's clock the strongest possible trust anchor, post-MAC binding lets the institution decouple SDK availability from TSA availability.

§10.14 should pick one model (or normate both with explicit naming) when the v1.x extension lands. Currently the section says "MAY define a normative attribute" but doesn't anticipate the placement question. Worth adding to the v1.x candidate scope:

```
A v1.x extension that adds RFC 3161 timestamps MUST specify whether
the attribute is bound into the per-event canonical bytes (pre-MAC,
hot-path TSA round-trip) or recorded post-MAC (decoupled TSA path,
institution-trusted attestation). The two models compose differently
with the §5 integrity binding; the spec extension MUST name the
posture and the SDK MUST implement the chosen model exclusively per
chain.
```

This pre-empts the implementer-side ambiguity that today's "RECOMMENDED" wording leaves open.

### Late-binding attribute (§4.2.2 + §4.4)

**Q40. `ffiec.chain.late_binding` integrity binding posture.** §4.4 attribute table line 357: `ffiec.chain.late_binding | bool | optional`. §4.2.2 says: "Events that arrive at the ledger after the daily seal for their `received_at` UTC date is sealed are recorded with the per-entry attribute `ffiec.chain.late_binding = true`." The phrasing "are recorded with" is ambiguous about whether the attribute is added by the SDK at capture time (would be MAC-bound per §5) or by the ledger after ingest determines the late-binding condition (would NOT be MAC-bound; institution-trusted ledger metadata).

**Status:** Partial

The semantic argument: late-binding is determined by comparing the entry's `received_at` to the seal status of that UTC day. The SDK at capture time cannot know this (it doesn't know when the ledger will seal the day, or whether ingest delays will land the event after seal). Only the ledger, post-ingest, knows. So the attribute must be added by the ledger AFTER the chain MAC was computed; it cannot be in the canonical bytes the MAC sealed.

That makes `ffiec.chain.late_binding` parallel to `received_at` (§4.2.2): institution-trusted ledger metadata, not chain-MAC-bound. Spec §4.2.2 explicitly names received_at's trust posture; §4.4 + §4.2.2 should do the same for `late_binding`. Recommend a one-paragraph addition to §4.2.2:

```
Trust posture for `ffiec.chain.late_binding`. The attribute is
determined and stamped by the ledger after ingest (the SDK at
capture time cannot know whether its event will arrive late
relative to the seal job). The attribute is NOT part of the SDK-
produced canonical bytes that the per-event MAC sealed; the
verifier reads it as institution-trusted ledger-side metadata
parallel to `received_at`. A ledger that mutates the attribute
after the entry is appended is a control failure independent of
chain integrity; SOC engagements test this control via the
audit-procedures.md storage-integrity sample.
```

Without this, an implementer reading §4.4 + §5 together might believe the attribute is MAC-bound and a tampered value would be detected at §7 step 9 — which is false; tampering with `ffiec.chain.late_binding` is detected only by the institution's storage-integrity controls, not by the chain.

### Daubert four-factor framing (§1.1)

**Q41. §1.1 "Known error rate" omits SDK-process compromise (Adversary F per `09-threat-model.md` §2.6).** The spec text:

> "A successful false-negative — a tampered chain that verifies as PASS — requires the simultaneous compromise of three independent custody layers: the tenant's IKM (held in HSM/KMS), the institution's ledger storage (append-only with operator-side controls), and the HSM signing key (FIPS 140-2 Level 3 or higher)."

**Status:** Partial

The `09-threat-model.md` §2.6 names a fourth class:

> "**Adversary F — Application process compromise.** The AI agent's host process is compromised (RCE, container escape, malicious dependency). Goal: forge events going forward, claim they reflect legitimate decisions. Defense: bounded forward-only attack window. The compromised process holds the session key and can produce valid-looking events."

A compromised SDK process, holding the legitimate session key derived from the legitimate IKM, can produce chain entries that pass §7 step 8 (fingerprint match — same IKM) AND step 9 (MAC match — same session key). The institution's IKM is NOT compromised; the ledger storage is NOT compromised; the HSM is NOT compromised. The chain entries the compromised process produced will appear in the next daily Merkle root and the seal will sign them. The verifier will report `Status: PASS`. The chain has produced a verifying false-negative: events that the legitimate AI agent did not generate, presented as if it had.

§1.1 omits this scenario. For a court witness laying foundation under FRE 702 / *Daubert*, the omission is meaningful: the witness asked "what could produce a false-negative in your system?" answers from §1.1 by naming three custody layers, and the cross-examination produces an SDK-process compromise scenario the witness did not anticipate. The four-factor framing's "Known error rate" should be honest about this.

Recommend extending §1.1 with a fourth bullet:

```
- **Application-process compromise as a fourth class.** A fourth
  scenario — compromise of the SDK process holding the session
  key — produces chain entries that verify as PASS for as long as
  the compromise persists. This is a forward-only attack window
  (past chain entries cannot be retroactively altered) bounded by
  the institution's host-hardening, intrusion-detection, and
  master-key-rotation controls. The chain composes alongside
  these controls; the chain alone does not defend against an
  attacker who has root on the SDK's host. See
  `09-threat-model.md` §2.6 (Adversary F) for the residual-risk
  posture and the compensating controls (anomaly detection on
  the captured stream, out-of-band agent-behavior monitoring,
  third-party intrusion detection).
```

This makes the §1.1 framing honest under cross-examination and aligns the spec body with the threat-model document.

### Cross-document drift

**Q42. §10.15 Pattern A names `master.cross_region_replication_completed` as a §10.2 operational event but the §10.2 list reads "Multi-region replication" (no period) — confirm the event-name spelling is locked.**

**Status:** Confirmed

§10.2 line 654: `Multi-region replication: master.cross_region_replication_completed (per-tenant per-day per-source-region replication evidence; per §10.15 Pattern A reconciliation)`. Matches the §10.15 text. Good.

This is a confirmation, not a finding. Included so the convergence tracker can confirm the cross-reference is intact.

---

## Per-role roll-up

| Aspect | Status |
|---|---|
| Round-3 closures (Q24–Q35) | 11 closed via Wave-3 close-out |
| §10.15 "region" undefined (Q36) | Gap |
| §10.15 run-locality enforcement model + region attribute (Q37) | Partial |
| §10.15 Pattern A failover replication-completeness verifier blindness (Q38) | Partial |
| §10.14 trusted-time token placement (pre-MAC vs post-MAC) (Q39) | Gap (forward-scope, v1.x candidate) |
| §4.2.2 + §4.4 `ffiec.chain.late_binding` integrity-binding posture (Q40) | Partial |
| §1.1 Daubert "Known error rate" omits SDK-process compromise (Q41) | Partial |
| §10.2 ↔ §10.15 cross-reference (Q42) | Confirmed |

| Status | Count |
|---|---|
| Answered (round-3 closures) | 11 |
| Confirmed | 1 |
| Gap | 2 (Q36, Q39) |
| Partial | 4 (Q37, Q38, Q40, Q41) |

---

## Where I would prioritize

1. **Q41 (§1.1 Daubert framing — SDK-process compromise omission).** Court-facing posture. An expert witness laying foundation under *Daubert* SHOULD have a complete residual-risk picture; today's text invites a cross-examination gotcha. One-paragraph spec patch. Highest priority because it's defensive of the institution's actual courtroom posture, not a future-vendor-conformance concern.

2. **Q40 (`ffiec.chain.late_binding` trust posture).** Same shape as the round-3 received_at fix that landed in Wave-3: explicit naming of the trust posture so an implementer doesn't believe the attribute is MAC-bound. One-paragraph addition to §4.2.2.

3. **Q36 (§10.15 "region" definition).** §3 is where definitions live; one-row addition closes it.

4. **Q37, Q38 (multi-region operational details).** Two paragraphs in §10.15 cover them: (a) per-process region binding for run-locality enforcement, (b) the verifier's working-paper reference to the per-region replication evidence.

5. **Q39 (trusted-time placement).** Forward-scope, v1.x candidate. The v1.0 text "MAY define a normative attribute" does not need to resolve the placement today; the v1.x extension MUST name the posture when it lands. One forward-commitment line in §10.14.

---

## Stopping criterion

This drop carries 2 Gap + 4 Partial findings against the v1.0-final-amendment spec as it stands today (post-Wave-3 close-out). The findings cluster into multi-region operational specificity (Q36-Q38), late-binding integrity posture (Q40), trusted-time forward-commitment (Q39), and one Daubert-framing residual-risk completeness gap (Q41). None of them break the four primitives' cryptographic integrity claims; they affect cross-vendor implementation comparability, court-facing evidence framing, and forward-scope clarity. A working-group cycle can close them in editorial pass; Q41 is the one I'd stage first because it's defensive of the institution's existing courtroom posture.

The Q5 vendor-side gap (Ed25519 absence on the candidate vendor's signer) is unchanged from rounds 1, 2, and 3.

The cross-rounds tally: 42 distinct questions raised across rounds 1-4 (12 + 11 + 12 + 7). Of those, 30 closed via spec text, 2 confirmed, 10 open at round-4 close (1 vendor-side Q5, 9 spec-quality items pending working-group cycle). The convergence dynamic remains striking — three reviewer waves on a single calendar day, plus a fourth round the same week, has produced a spec that's substantially more mature than v1.0-final issued nine days ago.
