# Outside drop — Herald-ecosystem conformance review against `chain-of-custody-v1.md`, round 2

## Persona

**Reviewer.** Same independent technical reviewer as the first Herald-vendor-conformance drop (`herald-vendor-conformance-review.md`). Returning to walk the spec as it stands today against the same vendor candidate. Read v1.0-final in its current form; consulted the published Herald SDK surface; deliberately did not re-read my round-1 drop before rereading the spec, so the second reading is independent.

**Reading angle.** First-pass round-1 found 7 Gap and 4 Partial findings against the previous version of the spec. Returning to confirm the spec absorbed what it agreed to absorb, and to find what the new sections (the v1.0-final additions to §3.1, §4.1.2, §4.4.1, §4.4.2, §10.7, §10.10.1, §10.10.2, §10.11) introduce. The new sections are substantial — about a third of the spec body — and adding that much normative material in one cycle inevitably produces small seams.

**Method.** Re-read spec §1–§13 cover-to-cover; spot-checked design 02 (chain construction), design 05 (OTLP wire), and the v1.0-final test-vector corpus including `008-jcs-edge-cases/`. I did not re-derive byte values this round — round-16's cryptographic-engineer drop and round-1's confirmation already pinned that work; my contribution is on the operational and schema-coherence dimension.

---

## Round-1 closure confirmation

The v1.0-final spec absorbed substantially everything round-1 raised. Recording the confirmations here so the convergence tracker can close the prior drop:

| Round-1 question | Round-2 status |
|---|---|
| Q1 — vendor-namespace conformance pathway | **Closed.** §4.1.2 lifts the binary FFIEC-vs-vendor posture, the on-disk `hkdf_inputs_digest` witness, the verifier `--posture=ffiec` invocation, and the institution-side CC8.1 documentation requirement. The spec text is exactly the (b) variant I recommended ("vendor-flag mode permitted; chain entry is FFIEC-conformant only when produced under the §4.1 byte values"). Clean close. |
| Q2 — verifier refusal on namespace drift | **Closed (was already correct).** §7 step 2 unchanged; cross-posture detection works. |
| Q3 — RFC 8785 JCS conformance bar | **Closed.** §5 now says "Implementations passing only the basic fixtures are not v1.0-conformant — the JCS edge cases are part of the conformance bar." `spec/test-vectors/008-jcs-edge-cases/` shipped with `fixture.json`, `expected.json`, `_compute.py`, `_validate_inline.py`, and a `computation_method.md`. The bar is now testable. |
| Q4 — mid-write truncation refusal mechanism | **Closed.** §7 has a normative "Implementation note" naming the byte-level seek discriminator, calling out the stdlib-line-reader pitfall by language, and tying the failure to the `audit_file.truncation_detected` operational event. |
| Q5 — Ed25519 only at v1.0 | **Confirmed.** §4.3.2 dispatch is unchanged; v1.0 is single-algorithm Ed25519. Vendors shipping non-Ed25519 signers are non-conformant on Primitive 3 unless the signer ships an Ed25519 path alongside. |
| Q6 — software-key adapter exclusion strictness | **Closed.** §10.7 lifted to "unreachable in production through any combination of build-flag, packaging, and configuration that a normal misconfigured deployment cannot bypass at run time." Compile-time exclusion, packaging exclusion, and equally-strict alternatives are all conformant; run-time-only env-var gating is explicitly non-conformant. The text is more flexible than my round-1 phrasing while keeping the regulator-visible line. |
| Q7 — tenant_id character-class enforcement boundaries | **Closed.** §3 names enforcement at SDK construct + verifier file-header pre-flight; §7 step 3a does the verifier check; §3.1 documents three legacy-migration patterns with named CC8.1 documentation requirements. |
| Q8 — gen_ai.* SDK-side enforcement | **Closed.** §4.4 has a dedicated "SDK-side enforcement of `gen_ai.{request,response}.model` (normative)" paragraph requiring SDK refusal at write time before the MAC compute. Verifier check (§7 step 12a) remains as defense-in-depth. |
| Q9 — multi-vendor topology (SDK + ledger from different vendors) | **Confirmed.** §4.2 + spec README §"Conformance" make this explicit. |
| Q10 — in-process attribute names vs OTLP wire-form | **Closed.** §4.4 has both an "In-process attribute names vs wire form" paragraph and a complementary "OTLP-collector transformation pass-through" paragraph. The mapping requirement is named in CC8.1 testable terms. |
| Q11 — tenant_id and run_id integrity-binding | **Confirmed.** §5 inclusion list is unchanged. |
| Q12 — per-file envelope vs §4.2 daily seal disambiguation | **Closed.** §4.2 schema has a "Supplemental per-file envelopes (informative)" paragraph naming the three reasons per-file envelopes do NOT satisfy §4.2 (cadence, coverage, sign_payload shape), and explicitly stating "naming a per-file envelope as 'the seal' in audit documentation is non-conformant." |

**Net round-1 outcome.** 8 closures + 3 confirmations + 1 vendor-side gap (Q5, which is on the Herald team to close, not the spec). Spec contributed real text in response to 8 of my 12 questions. That is exactly the convergence dynamic the project's feedback README describes.

---

## Round-2 questions

These are new questions surfaced by the v1.0-final additions. They cluster around schema-coherence inside the new sections (§4.4.1 routing, §4.4.2 deployment-intent, §10.10.2 algorithm-rotation Pattern B) and a small set of definitional gaps (`chain_kind`, `received_at`, hourly/weekly cadence SLA).

**Q13. `chain_kind` is referenced in §5's MAC-covered field list and in `docs/design/02-chain-construction.md` §3.1 but is not defined anywhere in the spec body. What values does `chain_kind` take, what is the type, and what is the semantic discriminator from the OTel `kind` field that sits next to it?**

**Status:** Gap

§5 line 401 enumerates the integrity-bound OTel-envelope fields:

> The OTel envelope: `trace_id`, `span_id`, `parent_span_id`, `name`, `timestamp_ns`, `duration_ns`, `attributes`, `resource`, `severity`, `kind`, `chain_kind`.

`kind` is the OTel SpanKind (`server`, `client`, `internal`, `producer`, `consumer`) per the OpenTelemetry spec. `chain_kind` is FFIEC-specific — it appears in the published Herald cross-language fixture (`event_canonical_hex` decodes to a JCS object whose `chain_kind` field is the literal string `"audit"`) but the spec body never enumerates its values. Two implementations cannot agree on `chain_kind` without a spec definition: one might emit `"audit"`, another `"audit.event"`, another `"AI"`. The MAC then diverges on a field that's part of the canonical bytes.

Recommend adding to §3 (Definitions) and to the §4.4 attribute table:

```
chain_kind  string   yes   Chain-of-custody event class. One of:
                            "audit"      — application audit event (the default)
                            "model_call" — chain entry representing an LLM invocation
                            "tool_call"  — chain entry representing a tool invocation
                            "routing"    — chain entry from §4.4.1
                            "translation" — chain entry from §10.11
                            "operational" — chain entry for control-evidence operational events
                           The verifier MUST reject any value not in the enumerated set
                           with `chain_kind out of v1 enumeration at seq N`.
```

The values are illustrative; what matters is that the enumeration is locked and the verifier rejects non-conforming values. Without this, two SDKs producing the same logical event under different `chain_kind` strings produce different canonical bytes and the cross-vendor "byte-identical" claim breaks.

The published Herald fixture uses `chain_kind = "audit"` for a generic audit event, suggesting the working assumption is a small enumerated set. Lifting that assumption into the spec is the close.

**Q14. §10.10.2 Pattern B (split-day with two seal records) does not name the field that lets a verifier partition events between the two seals. The seal record schema (§4.2) has no "covers events captured before time T" or "covers seq range" field. How does the verifier mechanically know which events go with which Merkle root?**

**Status:** Gap

§10.10.2 Pattern B says:

> The institution's seal job produces TWO seal records for the rotation day: one covering events captured before the algorithm-rotation boundary (signed under the old algorithm), one covering events captured after (signed under the new algorithm). Each seal records its own `merkle_root` over its event subset. ... The institution's chain-operations runbook documents the split-time so the verifier knows to expect two seal records for the day.

The runbook is institution-side documentation, not a chain-integrity-bound artifact. A verifier walking the ledger sees two seal records both with `seal_date = "2026-05-12"` (say), each with its own `merkle_root`. Without a field on the seal record that names the events the seal covers, the verifier has three choices:

- Recompute Merkle over ALL events for the day, fail to match either seal's root, report `merkle root mismatch` — incorrectly.
- Recompute over the first half of events (by some heuristic — `received_at` < some unknown threshold), fail to find the threshold, give up.
- Read the institution's runbook out-of-band, learn the split-time, partition events accordingly, recompute. This breaks the "verifier walks the ledger without trusting the institution's runbook" property.

Recommend extending the seal record schema with two optional fields under Pattern B:

```
covers_received_at_min  RFC 3339 UTC  required for Pattern B  Inclusive lower bound of received_at
covers_received_at_max  RFC 3339 UTC  required for Pattern B  Exclusive upper bound of received_at
```

A verifier reading two seals for the same `seal_date` and seeing `covers_received_at_*` on both partitions events deterministically. The fields are bound into `sign_payload` (extending §4.3 by two more lines) so the partition is itself integrity-protected: a malicious actor rewriting the boundary would need a fresh HSM signature.

A simpler alternative if Pattern B is rare in practice: deprecate Pattern B in v1.x and require Pattern A (same-day cosign) for all algorithm rotations. §10.10.2 already says Pattern A is preferred; promoting it to MUST closes the partition-ambiguity question.

**Q15. §4.2.1 names hourly, daily, and weekly cadences. §4.3 names a 60-minute publish SLA "of UTC midnight on the day after the tenant-day it covers." That phrasing is daily-cadence-shaped. What is the publish SLA for hourly cadence (each seal covers a 1-hour window — must the seal be published within 60 minutes of the hour boundary, or shorter)? What is the publish SLA for weekly cadence?**

**Status:** Gap

A vendor implementing hourly cadence has to pick a publish SLA without spec guidance. Three candidates:

- **(a) Same 60-minute SLA.** A seal covering 13:00–14:00 must be published by 15:00. Cleanest.
- **(b) Proportional SLA.** A seal covering N hours must be published within max(60min, N hours / K) for some K. Lets hourly cadence have a tighter SLA than daily.
- **(c) Cadence-specific override.** Hourly: 15 min. Daily: 60 min. Weekly: 240 min. Documented per cadence.

§4.3 implies (a) by default ("within 60 minutes of UTC midnight on the day after"), but the language only handles daily. Recommend extending §4.3:

```
The signed root MUST be appended to the ledger within 60 minutes of the END
of the tenant-day's seal window:
- daily cadence:  by 01:00 UTC on the day after
- hourly cadence: by H+1:00 UTC where H is the hour the seal covers
- weekly cadence: by 01:00 UTC on the Monday after the seal-week (or
                  whichever day-of-week the institution declared as week-end)
```

The 60-minute number transfers naturally; the spec just needs to make the "end of the seal window" notion explicit so hourly and weekly are unambiguous.

**Q16. §4.4.1 `audit.routing.providers_attempted` semantics on the first `attempt` event.** The schema says: "Ordered list of provider identifiers attempted before this event ... The list MAY be a single element on `attempt`." Does the first `attempt` of a request carry an empty list (no prior attempts) or a single-element list (just the current attempt)? On a `success` after one failover, does the list contain `[failed_provider_X]` or `[failed_provider_X, current_provider_Y]`?**

**Status:** Gap

The "before this event" wording reads as "prior attempts only," but the "MAY be a single element on `attempt`" phrase reads as if the current attempt counts. Two reasonable readings exist; institutions implementing routing capture under different readings produce non-comparable evidence.

Recommend a worked example in §4.4.1:

```
Worked example — single-provider success:
  attempt:    providers_attempted = ["openai-gpt-4o"], provider_chosen = "openai-gpt-4o"
  success:    providers_attempted = ["openai-gpt-4o"], provider_chosen = "openai-gpt-4o"

Worked example — failover-then-success:
  attempt:    providers_attempted = ["openai-gpt-4o"], provider_chosen = "openai-gpt-4o"
  failover:   providers_attempted = ["openai-gpt-4o"], failover_reason = "timeout"
  attempt:    providers_attempted = ["openai-gpt-4o", "anthropic-claude-sonnet"],
              provider_chosen = "anthropic-claude-sonnet"
  success:    providers_attempted = ["openai-gpt-4o", "anthropic-claude-sonnet"],
              provider_chosen = "anthropic-claude-sonnet"
```

Pin the convention. Worked example is also useful to the audit-procedures.md P-33 sample-comparison procedure; without it, two SOC teams produce different sample-evaluation outcomes.

**Q17. §4.4.1 `audit.routing.failover_reason` enumeration is missing institution-side rate-limit / quota cases.** The list is `timeout | transport_error | provider_error | circuit_open | rate_limit | cost_threshold | policy_override | manual_override`. `rate_limit` is provider-side (the provider's API returned 429). The institution-side cases — the institution's per-tenant API quota, the institution's customer-tier-quota guard — are missing.

**Status:** Partial

A failover triggered by the institution's own rate-limiting infrastructure (e.g., the institution rejects an LLM call because the calling tenant exceeded its monthly quota) is a substantively different signal from a provider-side 429: the customer-side dispute response, the MRM committee's interpretation, and the SOC team's sample-comparison procedure all branch on the distinction. Today the institution would write `policy_override` for both; that loses the discriminator.

Recommend extending the enumeration:

```
quota_exhausted   The institution's per-tenant quota was exhausted; the
                  router moved to the next provider (or terminated with no call).
                  Distinct from rate_limit (provider-side) and cost_threshold
                  (per-call cost guard, not aggregate quota).
```

Or more comprehensively, split `rate_limit` into `rate_limit_provider` and `rate_limit_institution`, and add `quota_exhausted_institution` for aggregate-quota cases. The current single `rate_limit` value collapses three different signals.

**Q18. §4.4.1 routing-event coverage for "no call made" cases.** Spec §4.4.1 says: "Routing decisions that DO NOT result in a call (circuit-open precludes the attempt; rate-limit kills before retry) are still chained — the absence of a child LLM-call entry is itself evidence." But the four event types are `attempt | success | failover | circuit_state_change`. None of them name "the router evaluated policy and decided no provider was reachable, no call was made." Which event type does that case use?

**Status:** Partial

A pure "all providers circuit-open" case at the moment the request arrives looks like:

- The router observes circuits are all open. No `attempt` happens.
- The request returns to the caller with an error.
- The chain entry that records this... is what? `circuit_state_change` doesn't fit (no state change happened — circuits were already open). `failover` doesn't fit (nothing was being failed-over from). `attempt` with `provider_chosen = null`? `success` no, `success` is for the provider returning successfully.

Recommend adding a fifth event type:

```
audit.routing.refused
  Emitted when the router evaluates policy and determines no call can be made.
  Required attributes: providers_attempted (the list the router considered),
  refusal_reason (one of `all_circuits_open` | `no_provider_in_policy` |
  `quota_exhausted` | `cost_threshold_at_capacity` | `policy_override`).
  No provider_chosen (no provider was selected). No failover_reason (no
  prior provider was being failed over from).
```

Without this event type, the "no-call" case either gets squeezed into one of the other four events (mis-categorizing the evidence) or doesn't get recorded at all (the §4.4.1 stated requirement of "still chained" is not actually achievable). Given §4.4.1 explicitly says the absence-of-child-call case is part of the routing-decision evidence, the spec needs the event type to support the requirement.

**Q19. §4.4.2 `audit.deployment.intent` schema column says `optional`, but §4.4.2 "When emission is REQUIRED" makes it conditionally required for four named cases (A/B test, canary, multi-region drift, vendor-reroute observed). The table column and the prose disagree.**

**Status:** Partial

§4.4.2 schema row 344:

> `audit.deployment.intent` | string | optional | One of `production` | `ab_test` | ...

§4.4.2 "When emission is REQUIRED":

> Institutions operating any of the following MUST emit `audit.deployment.intent` ... A/B test ... canary ... multi-region deployment ... vendor relationship ...

The schema column should match the conditional language of `audit.deployment.canary_traffic_pct` ("required when `intent` is `canary`") and `audit.deployment.policy_version` ("required when any `audit.deployment.*` attribute is present"). Suggested edit:

```
audit.deployment.intent | string | conditional |
  REQUIRED when the institution operates any of the postures named in
  §4.4.2 "When emission is REQUIRED" (A/B test, canary, multi-region drift,
  vendor-reroute observed). Optional otherwise. One of `production` |
  `ab_test` | `canary` | `multi_region_drift` | `vendor_reroute_observed` |
  `unknown`.
```

Pure documentation defect; one-line fix.

**Q20. §10.11 ECOA adverse-action notice translation MUST chain the translation step but defines no attribute schema for the translation entry. Two institutions implementing §10.11 produce non-comparable evidence.**

**Status:** Gap

§10.11 says:

> Chain entries for ECOA adverse-action notices MUST chain the customer-language translation step ... The translation entry binds to the AI's original response via `parent_run_id` / `parent_seq` per §4.4 so the customer-side disclosure produces a coherent narrative.

That defines the binding but not the content. Required attributes a CFPB / ECOA examiner would expect on the translation entry:

| Attribute | Suggested |
|---|---|
| `audit.ecoa.translation.target_language` | The language code (BCP 47) the translation produced (e.g., `es-US`, `zh-CN`). REQUIRED. |
| `audit.ecoa.translation.source_language` | The language the original AI response was in. RECOMMENDED. |
| `audit.ecoa.translation.translator_kind` | One of `human` \| `llm` \| `translation_api` \| `glossary_lookup`. The institution's translation method. REQUIRED. |
| `audit.ecoa.translation.translator_id` | If `translator_kind` is `llm`, the model identifier (matches `gen_ai.response.model` shape). Required when `translator_kind` is `llm`. |
| `audit.ecoa.translation.glossary_version` | If the institution operates a regulated-language glossary the translation conforms to, the glossary version. Conditional. |
| `audit.ecoa.translation.output_hash` | SHA-256 hex of the customer-facing translated text. The text itself MAY be customer-PII; the hash binds the translation under the chain without binding the PII. REQUIRED. |
| `audit.ecoa.translation.delivery_method` | One of `mail` \| `secure_message` \| `email` \| `phone` \| `in_person`. The method by which the translated notice was delivered to the customer. RECOMMENDED. |
| `audit.ecoa.translation.delivery_timestamp` | RFC 3339 UTC of delivery. RECOMMENDED. |

A CFPB inquiry into "did this customer receive the adverse-action notice in their preferred language within the regulatory window?" answers from these attributes alone, without the institution surfacing supplementary evidence. The current §10.11 leaves each institution to invent its own schema, and the resulting cross-institution comparison breaks.

A subsidiary point: §10.11 should also say whether the actual delivery (the "this letter was mailed" event) is itself a chain entry, or whether the translation entry is the load-bearing record and delivery is institution-tracked separately. The customer-dispute-procedures.md document is referenced; if it answers this, a one-line cross-reference in §10.11 closes the question.

**Q21. The §4.3 `sign_payload` covers `algorithm`, `format_version`, `tenant_id`, `seal_date`, `merkle_root`, `hkdf_inputs_digest`. It does not cover `cadence`, `key_versions`, `late_binding_count`, `signed_at`, `public_key_id`, `dev_mode`. A malicious actor with write access to the seal record can rewrite any of those fields without invalidating the HSM signature.**

**Status:** Partial

Three of the unprotected fields have downstream verifier checks that catch silent rewrites:

- `cadence` — §7 step 12 cadence-mismatch check IF the institution's externally-asserted cadence matches the seal record. A coordinated forgery that rewrites both the seal-record `cadence` AND the institution's externally-claimed cadence would silently pass. The institution's cadence is normally recorded in regulator-approved control documents and the SOC engagement letter, so the coordinated forgery is harder than rewriting the seal alone — but it's not impossible.
- `dev_mode` — §7 step 12 refuses `dev_mode = true` under `--strict`. A flip from `true` to `false` would silently pass; a chain that the dev adapter produced could be presented as a production chain. This is a real attack against the §10.7 regulator-visible-line guarantee.
- `key_versions` — §7 step 7 lookup against per-event `key_version` distribution. A `key_versions` rewrite would surface as a per-entry mismatch only if the verifier cross-checks the seal's `key_versions` list against the events' actual `key_version` distribution. The spec does not make this cross-check normative.

The remaining unprotected fields (`late_binding_count`, `signed_at`, `public_key_id`, `hsm_cluster_member`) are advisory or forensic and rewriting them does not immediately compromise integrity.

Recommend extending the §4.3 `sign_payload` to include `cadence` and `dev_mode` at minimum:

```
sign_payload = "ffiec.chain-of-custody.v1\n" ||
               algorithm                  || "\n" ||
               format_version             || "\n" ||
               tenant_id                  || "\n" ||
               iso8601_date(tenant_day)   || "\n" ||
               hex(merkle_root)           || "\n" ||
               hex(hkdf_inputs_digest)    || "\n" ||
               cadence                    || "\n" ||
               (dev_mode ? "1" : "0")
```

Two more lines closes both rewrite paths. A v1.0-conformance test vector for this gets generated mechanically. The cost is one wire-format bump (chains using the longer `sign_payload` are not byte-identical to chains using the shorter form), so this is a v1.x change rather than a v1.0 patch — but the gap should be tracked.

`key_versions` is harder to bind into `sign_payload` because it's a list of variable length. The cleanest close is a verifier rule: §7 step 11 also asserts `seal.key_versions == sorted(unique(entry.key_version for entry in day))`. Mismatch → `seal.key_versions does not match per-event key_version distribution`. That doesn't bind `key_versions` into the signature but catches silent rewrites at verification time.

**Q22. §4.2.2 names `received_at` as the day-boundary discriminator, but `received_at` is not in the §5 MAC-covered field list and is not in the §4.4 OTLP attribute table. The verifier partitions events into seals by `received_at` UTC date — but the verifier reads `received_at` from the ledger's storage, not from the canonical bytes. The spec should make this explicit.**

**Status:** Partial

The current §4.2.2 phrasing implies `received_at` is institution-trusted ledger-side metadata, but the spec never says this in so many words. A first-time implementer reading §4.2.2 alongside §5's MAC-covered field list reasonably wonders: should `received_at` be added to the canonical bytes so the verifier doesn't have to trust the ledger's storage?

The answer is no: `received_at` is server-side and post-capture, so binding it under the SDK's MAC would require the SDK to know its own future `received_at` value, which is impossible. The right design is exactly what the spec has — `received_at` is a ledger-side column the verifier reads as institution-trusted metadata. But that design choice should be stated, with the trust posture named.

Recommend a one-paragraph addition to §4.2.2:

```
Trust posture for received_at. The received_at field is stamped by the
ledger server upon ingest and stored alongside the captured event; it is
NOT part of the SDK-produced canonical bytes that go into the per-event
MAC. The verifier consumes received_at as institution-trusted ledger-side
metadata for day-boundary partitioning. The institution's CC8.1 (or
equivalent) control description names the ledger's append-only storage
posture and the operational controls preventing post-ingest rewriting
of received_at. A ledger that mutates received_at after ingest is a
control failure independent of the chain's cryptographic integrity;
SOC engagements test this control via the audit-procedures.md storage-
integrity sample.
```

The trust posture is already implicit. Making it explicit prevents an implementer from over-engineering (binding received_at into the canonical bytes, which can't work) and prevents an examiner from under-evaluating (failing to test the ledger's storage-integrity controls because they assume the chain's MAC covers received_at).

**Q23. §10.5 lists "AWS KMS at default tier and Azure Key Vault Standard are not conformant." "Default tier" is not an AWS product term.**

**Status:** Nit

AWS distinguishes between "AWS KMS" (multi-tenant managed service, FIPS 140-2 Level 2) and "AWS CloudHSM" (single-tenant dedicated HSM, FIPS 140-2 Level 3). "Default tier" is ambiguous — implementers asking "does this AWS KMS feature qualify" cannot match "default tier" against the AWS console.

Suggested precise language:

```
Acceptable HSM products at v1.0 publication include AWS CloudHSM Classic
(FIPS 140-2 Level 3), AWS CloudHSM v2 (FIPS 140-2 Level 3), Azure Managed
HSM (FIPS 140-2 Level 3), Azure Dedicated HSM (FIPS 140-2 Level 3), and
Google Cloud HSM (FIPS 140-2 Level 3). AWS KMS without CloudHSM backing
(the multi-tenant managed-key tier; FIPS 140-2 Level 2) is not conformant.
Azure Key Vault Standard (FIPS 140-2 Level 1) and Azure Key Vault Premium
(FIPS 140-2 Level 2 with software-protected keys) are not conformant.
The Azure Key Vault Premium "HSM-protected key" feature uses Azure
Managed HSM under the covers and IS conformant.
```

Pure documentation precision. The cloud-hsm-guide is the authoritative product-mapping but the spec's one-line summary should be unambiguous to a reader who hasn't yet found the guide.

---

## Per-role roll-up

| Aspect | Status |
|---|---|
| Round-1 closures (Q1–Q12) | 8 closed / 3 confirmed / 1 vendor-side gap (Q5) |
| `chain_kind` definition | Gap (Q13) |
| §10.10.2 Pattern B partition mechanism | Gap (Q14) |
| §4.2.1 hourly/weekly cadence publish SLA | Gap (Q15) |
| §4.4.1 `providers_attempted` semantics | Gap (Q16) |
| §4.4.1 `failover_reason` enumeration | Partial (Q17) |
| §4.4.1 "no call made" event type | Partial (Q18) |
| §4.4.2 `intent` schema column "optional" vs prose REQUIRED | Partial (Q19) |
| §10.11 ECOA translation entry attribute schema | Gap (Q20) |
| §4.3 `sign_payload` coverage gaps (`cadence`, `dev_mode`, `key_versions`) | Partial (Q21) |
| §4.2.2 `received_at` trust posture | Partial (Q22) |
| §10.5 "AWS KMS at default tier" terminology | Nit (Q23) |

| Status | Count |
|---|---|
| Answered | 8 (round-1 closures) |
| Confirmed | 3 (round-1 confirmations) |
| Gap | 4 (Q13, Q14, Q15, Q20) |
| Partial | 5 (Q17, Q18, Q19, Q21, Q22) |
| Nit | 1 (Q23) |

---

## Where I would prioritize

1. **Q13 (`chain_kind` definition).** Cross-vendor byte-identity claim depends on a closed enumeration. One-paragraph spec patch.
2. **Q14 (Pattern B partition).** Without the partition field, Pattern B is operationally unverifiable. Either deprecate Pattern B in favor of mandatory Pattern A, or add the time-window fields.
3. **Q20 (ECOA translation schema).** CFPB / ECOA evidence is not just an integrity claim — it's a cross-institution comparison artifact. Without a normative attribute schema, the cross-comparison is anecdotal.
4. **Q21 (`sign_payload` coverage).** `dev_mode` rewrite is a real attack against the §10.7 regulator-visible line; closing it is a v1.x signature-format change but worth tracking now.
5. **Q15 (cadence SLA), Q16-Q19 (routing + deployment-intent schema cleanup).** All single-paragraph patches; group them into one editorial pass.
6. **Q22 (`received_at` trust posture).** Doc-only clarification.
7. **Q23 (HSM product terminology).** Doc-only.

Q5 is on Herald to close, not the spec. The Herald team's published surface is one-of-four primitives in (Primitive 1 architecturally complete, byte-form gap closed via §4.1.2 vendor-flag mode); Primitives 2 / 3 / 4 are vendor-side work the institution supplies via a separate ledger + HSM + OTLP-collector path.

---

## Stopping criterion

This drop carries 4 Gap + 5 Partial findings against the v1.0-final spec. Per `docs/feedback/README.md`, the convergence target is 0 Gap + 0 Partial. The findings cluster into spec-quality work (mostly under §4.4.1, §4.4.2, §10.10.2, §10.11) and one definitional gap (`chain_kind`). None of them block the four primitives' cryptographic integrity claims; they affect schema-coherence, multi-institution comparability, and rare-rotation operability. A working-group cycle can close most of them in editorial pass.

The Q5 vendor-side gap (Ed25519 absence on the candidate vendor) is unchanged from round 1 and is on the Herald team's roadmap, not the spec.
