# Round 13 — Outside reviewer (AI platform reliability engineer)

**Reviewer.** Maya Patel. Principal AI Platform Engineer at a tier-1 bank holding company. Fourteen years on production observability and resilience for ML/AI systems. Lead architect for her bank's multi-LLM routing layer (Resilience4j + custom failover, ~40M LLM calls/day across eight providers, in production for nineteen months). Brought into the round-13 review cycle as an outside voice on a question that surfaced from a vendor's reference-implementation engineering decision and that the regular six personas did not raise.

**Reading angle.** Production multi-LLM operations. Specifically: does the chain spec adequately address the *routing decision* — the institution's choice of which LLM provider to call when, why, and what circuit-breaker state at the moment of decision — as a first-class captured event, or does the spec leave it to institution discipline?

**Scope read.** Spec `chain-of-custody-v1.md` §4 (capture and chain construction); design 02 (chain construction, hot path); design 05 (OTLP wire); design 09 §2.6 (Adversary F — application process compromise); audit-procedures (any P-N for capture-completeness); MRM-COMMITTEE-BRIEF (model-substitution oversight); regulator-pack/ai-policy-alignment (where routing-related policy might live); customer-dispute-procedures (whether routing-decisions enter the dispute evidence path).

**Stopping criterion declared up front.** One Partial. Escalates to Gap if round-13 does not take a position. The question is settle-able by either spec language or explicit silence with reasoning; the current state is silence without reasoning, which is the worst of the three options.

---

## Headline

The spec captures the LLM call. It does not address the *routing decision behind* the LLM call. In a single-provider deployment that distinction is zero. In a multi-provider deployment with failover, circuit breakers, and cost-routing, the routing decision is itself a model-risk-relevant and audit-relevant event the institution makes. Two different institutions can capture the same LLM call event and produce wildly different chains because Institution A captured the routing decision and Institution B did not. The examiner has no normative reference for what to expect.

This isn't a chain-construction defect — the chain captures whatever events the institution emits with byte-correct integrity. It is a **scope-of-capture** question the spec has not addressed and a comparability question the next round of bank examinations will surface.

---

## The question

Should AI-routing decisions be captured into the chain by SDK construction (the SDK refuses to call an LLM without first capturing the routing decision), or by institution discipline (the application emits `audit.routing.*` events through the existing capture decorators), or is the question out of scope for the spec entirely?

The three positions, with reasoning:

### Position 1: capture by SDK construction (the router lives in the SDK)

The SDK exposes `route(providers, request)` as the only path to an LLM call. Every call goes through it. Every routing decision is structurally captured before any provider is reached.

**Pros.** Completeness-of-capture is a structural property; an application that routes silently is impossible. The SOC examiner asks "show me your routing log" and gets a complete answer from the chain alone. MRM concerns about silent model substitution are addressed by construction.

**Cons.** Mission creep — the SDK absorbs resilience, which is a different domain from capture. Forces banks to adopt the SDK's router instead of their existing tooling (Polly in .NET shops, Resilience4j in JVM shops, custom-built in many production environments). Larger SDK surface = larger supply-chain risk (R5 in §6). Conflates the actor and auditor of routing.

### Position 2: capture by institution discipline (router lives in the application)

The SDK provides decorators (`@herald.audit("audit.routing.*")` or equivalent). The application's router emits routing events. The chain captures whatever the institution emits; completeness is institution discipline, verified by SOC team's CC8.1 procedure.

**Pros.** SDK stays focused on capture (CUPID Unix-philosophy). Banks reuse existing resilience tooling. Smaller SDK surface; cleaner trust boundary (the SDK is auditor of routing, not router-and-auditor). Composable with any future routing primitive.

**Cons.** Capture-of-routing-decisions becomes institution discipline; an application that forgets to decorate its router routes silently. SOC examiner verifies completeness via separate procedure (e.g., compare router-invocation count to chain-captured routing event count); meaningful audit step but non-load-bearing. Institution must document the discipline as a CC8.1 control.

### Position 3: silent (institution chooses)

The spec normates chain construction (capture, MAC, seal, signature). What the institution chooses to capture is institution policy. Routing decisions are events like any other; the chain captures them when emitted.

**Pros.** Smallest spec scope, least likely to constrain unforeseen patterns. Institutions with existing resilience tooling are unaffected.

**Cons.** Examiners examining institution-vs-institution comparability lack a normative answer. Institution A captures `audit.routing.failover`; Institution B captures `audit.failover`; Institution C doesn't capture routing at all. Sample-testing varies. MRM committees lack a single normative reference for what routing visibility they should expect.

---

## My recommendation, with reasoning

**Position 2, plus a normative attribute schema for routing events when emitted.**

The spec should:

1. **Normate the schema** of routing events when emitted: `audit.routing.attempt`, `audit.routing.success`, `audit.routing.failover`, `audit.routing.circuit_state_change`. Define canonical attribute names (`audit.routing.providers_attempted`, `audit.routing.provider_chosen`, `audit.routing.failover_reason`, `audit.routing.circuit_state.{provider}`) so two implementations agree byte-for-byte.

2. **Normate the institution's responsibility** to capture (not the SDK's). Add a §X.Y section: "The institution captures routing decisions through the chain's standard capture path. The institution's CC8.1 control description names the routing events it captures and the discipline it operates."

3. **Add an audit procedure** (analogous to P-5, P-6) for routing-completeness: SOC team samples X% of LLM-call events and confirms the corresponding routing event is in the chain at the same `(run_id, near_seq)`. Naming the procedure makes the discipline testable.

4. **Stay silent on where the router lives.** Application or SDK is institution choice. The spec normates the captured event schema and the institution's responsibility; the implementation surface is open.

### Why not position 1

The "completeness by construction" argument is appealing but overreaches. The spec's integrity claim is about events that *are* captured; it does not claim completeness. Forcing the router into the SDK conflates two domains and forces a stack many banks don't want. CC8.1 already exists as the control area for completeness-of-capture; routing is one more event class that procedure tests. Building completeness-by-construction into the SDK adds maintenance burden, larger supply-chain risk, and a mission-creep narrative the FFIEC working group will rightly question.

### Why not position 3

Silence is operationally tempting but creates examination-comparability problems. Three institutions producing three different routing-event taxonomies is the failure mode. Schema standardization is cheap (a paragraph in the spec, a schema fragment in the OTel attribute namespace) and high-value (examiners sample-test consistently across institutions; MRM committees have a normative reference).

---

## What changes if accepted

Round-13 close-out adds:

- Spec §X.Y: routing events — schema, institution responsibility, capture-completeness audit procedure
- Audit procedures P-N: routing-completeness sample testing
- CSF 2.0 binding entry: routing events → DE.AE-03 + GV.AC-05
- IR scenario stub: routing event missing during incident reconstruction
- Negative test vector: routing event tampered with (extends N-vector approach)
- Design 02 §X: routing-event capture pattern
- Design 05 §X: `audit.routing.*` attribute namespace registration intent

The spec gets explicit; institutions get clarity; the SDK boundary stays clean. The Vidimus reference implementation (Herald.Py) keeps router code in `pipeline/router-demo/` rather than in the SDK proper, with the chain capturing routing events through existing decorators — exactly what Position 2 prescribes.

---

## Severity proposal

**Partial.** Round-13 takes a position (any of the three above) and the question is settled. If round-13 stays silent, the partial becomes a Gap because cross-institution comparability suffers and examiner sample-testing varies. The cost of taking a position is small (one §X.Y section + one audit procedure entry); the cost of staying silent compounds with each adopter.

---

## Acknowledgments

The question surfaced from the engineering team building Vidimus, the reference admin/SDK package the Herald.Py vendor is shipping. The team faced "router in SDK or application" as a concrete decision and asked the spec for guidance. The question generalizes from one institution's design choice to a spec-level concern affecting every multi-provider AI deployment. The vendor's reasoned position is Position 2 with the normative-schema addition; this review documents the analysis so round-13 reviewers can confirm or push back.

The vendor explicitly does not advocate for Position 1 (router-in-SDK) — they recognize the mission-creep and trust-boundary problems. They want the spec to either bless Position 2 with a normative schema, or take a different position with reasoning, so their engineering decision has spec backing rather than guesswork.

---

## Supporting questions surfaced during review

These are observations that came up while reading the corpus from the production-multi-LLM lens. Some may already be answered in spec text I haven't found; surfacing them lets round-13 reviewers either confirm "addressed in §X.Y" or take a position. Each is one paragraph, severity unclassified, intent is broaden the round-13 review surface.

### Q-1. Cross-process handoff IPC normate?

Design 02 §9.1 introduces the shared-`run_id` handoff pattern: the receiving process needs the sending process's last `payload_hash` to chain its first event. The spec normates the handoff event's schema (Q1 schema is solid) but is silent on the *IPC mechanism* that transfers `from_payload_hash` between processes. In a distributed deployment, processes don't share memory. Implementations may use Redis, a message queue, shared file storage, or pass-by-OTLP-attribute. Does the spec want to normate one mechanism so cross-vendor SDKs interoperate, or leave it institution-defined? If the latter, an institution-side documented procedure should be required (analogous to the master-key reception procedure in §2.9).

### Q-2. Empty-day seal cost during long-inactive tenants?

Spec §3 mandates a seal even on empty days (`merkle_root = SHA-256("")`). For a tenant that is truly inactive for weeks or months (paused pilot program, decommissioned business line awaiting compliance retention expiry), the HSM signing cost compounds. Is there a documented "dormant tenant" mode that relaxes the daily cadence to weekly or monthly with examiner approval? `regulator-pack/examiner-approval-template.md` covers cadence relaxation generically; specific dormant-tenant guidance would help institutions plan retirement of pilot programs without a multi-year empty-seal tail.

### Q-3. HSM unavailability SLA — normative threshold?

`docs/operator-guide.md` mentions seal-age thresholds: "< 25 hours under normal operation; > 72 hours triggers regulator notification." Are these normative or guidance? An institution operating with a stricter examiner relationship may face shorter de-facto thresholds. The working group should consider whether to normate the maximum acceptable seal-delay before mandatory regulator notification, or document that thresholds are institution-specific and named in the control description. Either is defensible; silence creates ambiguity.

### Q-4. OTLP-collector transformations vs canonical bytes?

OTel collectors routinely apply transformations (redaction processors, sampling processors, filter processors) between the SDK and the receiver. Any transformation that mutates the chain attributes — `ffiec.chain.payload_hash`, the `audit.*` payload, the `gen_ai.*` attributes captured under the MAC — breaks the canonical-bytes property the verifier depends on. Does the spec require these attributes to pass through transformations unchanged, and does it normate operator guidance for collector configuration? Two failure modes I have seen in production: (a) a redaction processor silently rewriting `gen_ai.request.messages` content, breaking the MAC at ingest re-verification; (b) a sampling processor dropping events, creating gaps the verifier reports as deletions. The spec should normate "the chain attributes are pass-through; collectors that mutate them are non-conformant for chain pipelines" and the operator guide should provide a worked configuration example.

### Q-5. Adversarial inputs to JCS canonicalization?

JCS (RFC 8785) canonicalization is rigorous, but the conformance corpus may not test adversarial inputs comprehensively. Examples of inputs that have caused divergence in canonical-form implementations elsewhere: very deeply nested JSON (stack-recursion limits), Unicode normalization edge cases (NFC vs NFD conflicts), control characters in strings (escape-sequence variants), surrogate pair handling at the encoder boundary, very long strings (size-limit edge cases at the JSON parser). Does the conformance corpus include an adversarial-input fixture set? If not, round-13 should add one. The threat model classifies this as part of R4 (insider collusion) rather than a separate adversary; in practice, an attacker who can craft an LLM prompt that produces a divergent canonical form on the SDK vs the verifier has an exploit. Worth specifying.

### Q-6. Provider-attestation chaining?

The chain captures what the SDK observed from the LLM provider's response. Some providers (notably Anthropic, OpenAI, Google for their newer endpoints) ship cryptographic attestation alongside responses ("we attest this model produced this content"). Should the spec normate a `gen_ai.provider_attestation` attribute that captures the provider's signature alongside the SDK's chain? The institution's evidence claim becomes stronger ("the chain proves we recorded what the provider sent; the provider's attestation proves the provider sent it"); the audit-procedure load is small (the institution captures the attestation as data; the verifier or an independent tool validates it against the provider's published key). The spec doesn't have to require it, but normating the attribute name and the validation expectation would let institutions for whom this matters operate consistently.

### Q-7. Verifier resource limits at tier-1 volumes?

A tier-1 bank captures on the order of 1M events/sec; an annualized retention of seven years times that volume is ~30 trillion events. The verifier as a single static binary is the design today. Can it complete a one-year-coverage examination in the typical 8-hour window? Streaming Merkle is in the glossary; is there a documented streaming verifier mode for tier-1 deployments, with bounded memory and resumable progress? What is the working group's tier-1 deployment expectation — does the verifier scale, or do tier-1 institutions run partitioned verifications (per business line, per region)? The cost-model documentation hints at this but doesn't fully resolve it; round-13 should commit to a position.

### Q-8. Tenant-onboarding workflow — normate the steps?

Adding a new tenant to the chain requires: (a) institution generates master HMAC key in tenant-controlled custody; (b) institution generates Ed25519 signing keypair in HSM; (c) institution publishes public key to tenant key registry; (d) institution registers public-key fingerprint with regulator; (e) institution provisions tenant routing in ledger config; (f) SDK handshakes for `key_version=1` IKM. Is this institution-defined or does the spec want to normate the process so cross-institution practice is consistent? `04-hsm-custody.md` covers the HSM steps; what's missing is the end-to-end onboarding choreography (and the order of operations: do you register the public key with the regulator before or after the first chain event under it?). A normative onboarding sequence would help new institutions stand up the chain confidently.

### Q-9. Within-day algorithm rotation?

Round-12 K-1 addresses cross-day dual-algorithm cosigned seals: a tenant-day fully under Ed25519 vs a tenant-day fully under Dilithium. What about within-day? An institution that rotates Ed25519→Dilithium at noon UTC on Tuesday: does the same day's seal cosign with both algorithms (so events captured before noon are covered by Ed25519 and events after noon by Dilithium), or does the day get split into two seals? If split, what is the day-boundary semantic for the sealed root (does each half-day produce its own root, and do they chain)? This is a pure spec question round-13 should address explicitly so K-1's resolution doesn't leave the within-day case ambiguous.

### Q-10. Edge / federated-AI connectivity gaps?

`docs/edge-and-federated-ai.md` addresses some of the long-disconnected pattern. For AI agents on intermittently-connected devices (banker laptops in branches with poor connectivity, remote-office tooling, federated learning workflows), captured events accumulate locally and upload when connectivity returns. The chain integrity holds (the SDK chains locally; the ledger re-verifies on ingest). The examination question is: at examination time, the institution's chain may show gaps between captured-at and received-at timestamps spanning days or weeks. What is the examiner's expected disposition? Are there normative thresholds beyond which delayed-upload becomes a finding, or is this institution-discretion documented in the control description? Round-13 could add specific guidance for federated/edge cases that complement the existing connectivity-loss IR scenario.

### Q-11. Declared-algorithm-posture publication mechanism?

Spec §7 step 11 case (c) says the verifier resolves the institution's declared algorithm posture from the institution's published configuration ("out of scope for the spec; the institution's control description names the file or registry entry"). This is undefined enough that two verifier implementations could read the posture from two different sources and disagree on case (c) disposition for the same seal. Is the posture (a) a flat file in the ledger snapshot, (b) a separate registry handed to the examiner, (c) a field in the seal record, or (d) institution-discretion-documented? The first three are conformance-relevant; the fourth creates a verifier-portability gap. Round-13 should pick one or normate that the institution's control description names the source so cross-vendor verifiers don't drift.

### Q-12. Tenant-id character constraint and legacy migration?

Spec §3 mandates `^[A-Za-z0-9_.\-]{1,255}$` to make the HKDF info-parameter unambiguously parseable (the `|` byte cannot appear inside `tenant_id`). Many institutions have legacy tenant identifiers using slashes (`acme/prod`), colons (`tenant:acme:prod`), Unicode names (CJK characters in some Asia-Pacific deployments), or longer-than-255-byte identifiers from upstream IAM systems. Is there a documented migration path for institutions whose existing tenant identifiers don't conform — opaque hash-of-legacy-id, controlled aliasing, simple rejection of non-conforming new tenants? The constraint is correct cryptographically; the operational migration shape is not in the corpus today. Adding a short §3.1 ("Legacy tenant identifier handling") would close the question.

### Q-13. Seal-timing audit procedure — 60-min SLA testability gap?

**Partially addressed — recording for completeness.** P-8 (Seal-age monitoring) tests the seal-age metric over the period. P-9 (Notification on 72-hour seal delay) confirms notification evidence for the 72-hour threshold per §4.3.1. Together they cover the SHOULD-notification posture. What's not crisply named is a procedure that specifically tests the §4.3 60-minute SLA ("signed root MUST be appended within 60 minutes of UTC midnight"). The 60-minute SLA is normative ("MUST"), the 72-hour threshold is regulatory-notification ("SHOULD"); the two are testable from the same `signed_at` data but a separate P-N for the 60-minute SLA would let the SOC team test directly against the spec's normative requirement rather than infer from the metric. Lightweight close: extend P-8 with a 60-min sub-test, or add P-N "60-minute seal-publication SLA: sample N seal records, confirm `signed_at - day_boundary < 60 min` for each."

### Q-14. Step 12a observability cost at tier-1 volumes?

Spec §7 step 12a (GenAI model identifier completeness check) runs inline during the per-event walk and is normative. For tier-1 institutions emitting on the order of 1M events/sec at peak, the additional namespace-prefix scan and field-presence check adds per-event verifier cost. Has the working group estimated the latency impact, or is this assumed to fit inside the existing per-event budget? The deferral-prohibition language is correct ("MUST NOT defer it to a second pass"), which means the cost lands on the hot verification path. A back-of-envelope budget table for tier-1 volumes would let institutions plan verification window capacity. Not a Gap; an operational-readiness question.

### Q-15. Verifier `--master-key` UX during examination?

Spec §7 fail-closed semantics: "an absent IKM (verifier has no master key) is a `--strict` FAIL when the verifier was invoked without `--master-key`; otherwise the verifier reports `structurally consistent, key-bound verification skipped`." This is a useful split — structural verification without master key, full verification with it — but the institution's procedure for handing the master key to the examiner during examination is not normated. Janet Buckley's round-12 J-1 finding (Training Module 3 invocation missing `--master-key`) suggests this is already a known confusion source for new examiners. Round-13 should add a master-key-handover procedure to either the operator guide or the regulator pack — what handover channel, what receipt/disposal evidence, what control-description language. Without it, examiners default to structural-only verification (the easier path) and miss the load-bearing per-event MAC check the rework was designed around.

### Q-16. Spec §12 change-log chronology presentation?

Spec §12 lists three change-log entries: v1.0-draft (2026-05-06), v1.0-final (2026-05-15), v1.0-rework (2026-05-06). Draft and rework carry the same date; final is nine days later. The narrative implied by the entries is "rework integrated into draft on the same day, then v1.0-final issued nine days later as the polished form" — but the chronological ordering as presented invites reader confusion ("rework AFTER final?"). For the FFIEC submission package the change log is read by examiners forming an opinion on spec maturity; ambiguous chronology is a presentation issue, not a normative one, but it costs the project credibility. Round-13 should consider re-ordering the entries to v1.0-rework → v1.0-draft → v1.0-final, OR re-dating the rework entry to align with v1.0-final issuance. Cosmetic, but visible.

### Q-17. Edge AI Pattern B conformance under v1.0?

Threat model R12 marks edge-device physical compromise as "Mitigation deferred to v1.1" with the note that v1.0 deployments "operate compensating controls (Pattern A per `edge-and-federated-ai.md`: per-device IKM in TPM/secure-enclave; reconciliation cadence baseline tuned per fleet)." But edge-and-federated-ai.md describes Pattern B (bulk session-key issuance) as the operational answer for "fleets without secure-enclave capability" — and Scenario 14 explicitly handles Pattern B in-service device compromise. This creates a conformance ambiguity: is Pattern B conformant for v1.0 (with documented residual + Scenario 14 IR coverage), or non-conformant for v1.0 (institutions must use Pattern A or compensating controls)? Institutions running ATM-class deployments with mass-market hardware face this question at procurement time. Round-13 should take an explicit position: Pattern B is conformant under v1.0 with [list of compensating controls], OR Pattern B is non-conformant for production v1.0. The current implicit "OK with residual" is operationally correct but conformance-ambiguous.

### Q-18. Federated-learning training scope boundary as the AI/MRM examiner pool grows?

edge-and-federated-ai.md states "The chain captures decisions, not training. Federated learning is the model-training phase; the chain operates on the inference / decision-making phase." Treasury's Financial Services AI Risk Management Framework (Feb 2026, cited in spec §11) and emerging CFPB/OCC guidance on training-data integrity may push examiner expectations toward training-phase coverage. The spec correctly scopes itself to inference; the question is whether the project commits to revisit if examiner consensus shifts. Round-13 could add a one-sentence forward-commitment to spec §1 or §11: "the chain's scope is decision-time integrity; training-phase integrity is out of scope for v1.x and is a candidate for a future spec version if examiner consensus on training-data integrity matures." Naming the scope boundary explicitly is institution-friendly; silently inheriting examiner discomfort is not.

### Q-19. Within-day rotation seal semantics under non-default cadence?

Spec §10.10 (Rotation crossing the seal boundary) handles the rotation-window edge case for daily cadence: late-arriving events under the old IKM are included in the next day's seal as late-binding entries; the day-after seal records `key_versions = [old, new]`. The reasoning is sound for daily cadence. For institutions operating hourly cadence (per §4.2.1, configurable per tenant), the rotation can cross MULTIPLE seal boundaries within a single day. Does §10.10's reasoning generalize directly (each crossed seal records `key_versions = [old, new]` until rotation completes), or are there hourly-cadence-specific edge cases? The verifier per-entry `key_version` lookup handles the case mechanically; the operational guidance for the institution may differ. Round-13 should confirm §10.10 applies at any cadence or document the hourly-specific variant.

### Q-20. Scenario 9 truncation-rate baseline and actionability threshold?

IR playbook Scenario 9 (audit file truncation) says the 36-hour clock starts only when "deliberate corruption" is suspected. What threshold separates baseline-noise truncations (normal SDK crash recovery) from actionable signal? An institution with frequent SDK crashes (resource-constrained edge fleet, dev environment promotion gaps) has a chronic Scenario 9 baseline that would, under strict reading, generate weekly 36-hour clocks. The institution's IR program should establish a baseline + alerting threshold (e.g., "truncation rate > 2× the institution's 90-day rolling baseline triggers Scenario 9 escalation review"). Round-13 should add P-N: "Scenario 9 truncation-rate baseline establishment" — sample N truncation events per period, confirm the institution has baseline metrics, confirm alerting thresholds are documented and reviewed quarterly. Without it, the "deliberate corruption" determination is ad-hoc.

### Q-21. Privacy-key rotation cadence vs chain-master rotation cadence relationship?

privacy-by-design.md prescribes `privacy_key_rotation_cadence: quarterly` for the privacy-store key. Chain-master IKM rotation is "tenant policy decision" (no spec-mandated cadence). The two keys serve different threat models but interact compositely: a privacy-key compromise + chain-master compromise enables an attacker to map tokens back to PII and forge new chain entries with valid tokens. Is round-13 going to take a position on the relationship — e.g., "privacy-key rotation MUST be at least as frequent as chain-master rotation" or "privacy-key and chain-master rotations MUST NOT both occur in the same operational window"? Without coordination guidance, an institution could legitimately rotate both simultaneously, double-extending the at-risk window for a single rotation operational error.

### Q-22. Customer-dispute "did the AI hallucinate fact X" question scope?

customer-dispute-procedures.md handles reproducibility ("would the AI have said the same thing on a different day"). It does not address the harder question: "the AI cited fact X about my account in its denial reasoning; was X actually in the bank's records, or did the AI hallucinate it?" The chain captures what the AI said with cryptographic integrity; whether what the AI said was true is a different question requiring cross-check against the bank's authoritative source-of-truth records (account systems, document repositories). The spec correctly scopes itself to integrity-of-AI-output; the customer-facing question of "did the AI invent facts" needs explicit handling. Round-13 could add a customer-dispute-procedures section: "When the customer disputes an AI-cited fact, the institution cross-checks the cited fact against [authoritative source]; the chain records the cross-check itself as an `audit.fact_verification.*` entry with the original decision's `(run_id, seq)` as parent." This makes the hallucination-detection procedure first-class without expanding the chain's spec scope.

### Q-23. Vendor-hosted topology master-key custody-handoff exit procedure?

vendor-hosted-controls.md describes the institution-vs-vendor control distribution. For vendor-held master-key custody, the spec says "document per-tenant isolation." It does NOT document the procedure for an institution exiting vendor-hosted (vendor migrates the master from vendor's HSM to the bank's HSM, OR the bank rotates to a new IKM and re-derives the chain forward). Without this procedure, the institution faces vendor-lock-in: switching vendors or going BYOC requires either (a) abandoning the historical chain (loses 7 years of audit history continuity) or (b) trusting the vendor to faithfully extract and transfer the master. Round-13 should add a vendor-hosted master-key handoff procedure to either vendor-hosted-controls.md or m-and-a-handoff.md — at minimum, document the conformant approaches (e.g., HSM-to-HSM key wrap with documented vendor-side authorization, or fork-and-rotate with explicit chain-discontinuity disclosure).

### Q-24. ECOA adverse-action notice translation — RECOMMENDED vs MUST?

customer-dispute-procedures.md §"Adverse-action notices (ECOA)" says "the recommended decision is to chain the translation step" and notes "a translation logged outside the chain is repudiable in the same way a vendor's logs are repudiable, which is exactly the gap the chain exists to close." The reasoning supports MUST, not RECOMMENDED. An institution that doesn't chain the translation has the customer-side repudiation gap the chain exists to close — that's the textbook case for a normative requirement. Round-13 should consider promoting "chain the ECOA translation" to MUST (or at minimum SHOULD with the gap explicitly named in the spec text, not just the customer-dispute-procedures doc). This is an asymmetric question: the institution that already chains the translation gains nothing from the change; the institution that doesn't gains everything from the change being normative.

### Q-25. Community-bank affordability under cadence-relaxation savings — quantify?

cost-model.md gives community bank ($1B–$10B) annual cost as $25k–$60k, dominated by HSM cost. Cadence relaxation (daily → weekly with examiner approval) is mentioned in spec §4.2.1 but the cost-model doesn't quantify the savings. A community bank's CFO making the adoption decision needs the math: what does the $25k–$60k drop to under weekly cadence? Under monthly cadence (with examiner approval)? The HSM cost may not drop much (per-cluster pricing dominates, signing-volume is a small fraction); the savings may be in operational complexity rather than direct cost. Round-13 should add a cost-model variant table: cadence vs annual cost for community / mid-size / tier-1, so the affordability-gating institution sees the explicit savings curve. This is a presentation-completeness question that affects spec adoption breadth.

### Q-26. At-scale emergency-rotation 60-second termination normativity?

at-scale-operations.md §"Master-key rotation at scale" / "Emergency rotation" describes a "Forced-handshake mechanism" where "Processes that don't respond within 60 seconds are terminated and replaced." For globally-systemic banks running thousands of processes across multiple regions and business lines, terminating non-responsive processes within 60 seconds requires deployment-orchestration sophistication some institutions don't have (especially institutions with batch-processing or scheduled-job AI usage where processes legitimately don't respond on the rotation timescale). Is the 60-second termination normative or operational guidance? Round-13 should clarify: institutions whose deployment posture cannot guarantee 60-second forced-handshake operate compensating controls (extended rotation window with documented bounded compromise + monitoring) rather than declaring the spec non-conformant. The current at-scale doc reads as guidance; the spec text should confirm.

### Q-27. Tier-1 reconciliation procedure variant — exhaustive vs sampled?

Spec §10.1 mandates weekly key-fingerprint reconciliation: match every observed `(tenant_id, key_version, key_fingerprint)` triple against the IKM roster. At billions of events/day per tenant, the weekly reconciliation has tens of billions of triples to compare. Is exhaustive comparison normative, or can institutions sample? The audit-procedures P-6 procedure references reconciliation but doesn't address the at-scale variant. Round-13 should add to spec §10.1 or to at-scale-operations.md: at tier-1 volumes, exhaustive reconciliation may be replaced by stratified sampling with documented coverage (e.g., 10K triples per `(tenant_id, key_version)` pair, with the strata covering all observed pairs in the period). The institution's MRM committee approves the sampling design; the SOC team confirms via a tier-1 sample-of-samples procedure.

### Q-28. Vendor-conformance attestation separate from vendor SOC report?

vendor-hosted-controls.md describes vendor SOC report consumption — the institution receives the vendor's SOC report annually and identifies CUECs. But the vendor's SOC report covers vendor-side OPERATIONAL controls of the chain implementation (HSM operations, ledger operations, etc.); it does NOT necessarily attest that the vendor's implementation passes the FFIEC conformance corpus. An institution accepting a vendor's word on "we're FFIEC-chain-conformant" without an independent attestation has a gap. Round-13 should consider establishing a vendor-conformance attestation procedure: the vendor produces an annual attestation that their implementation passes the FFIEC conformance corpus (test-vectors/), signed by the vendor's CTO or equivalent, with the corpus-version + test-results-hash recorded. This is separate from the SOC report; it's an FFIEC-spec-specific attestation. Alternatively, the spec could leverage existing patterns (e.g., a "FFIEC chain-of-custody conformant" certification mark with an authorized-implementer registry).

### Q-29. Operational events JSON schema — versioning?

`docs/operational-events.schema.json` is referenced as the canonical schema for operational events. The spec mandates implementations conform to it. As operational events evolve across spec versions (round-12 added `master_key.retired`, `cold_dr.dryrun_attestation_consumed`; future versions may add more), the schema versions independently of the spec. Is there a documented schema-version-vs-spec-version compatibility matrix? Institutions running v1.0 SDK + v1.0 ledger + v1.x verifier may have operational-events that the verifier-side consumer doesn't recognize. Round-13 should add to spec §10.2 or to the schema doc: schema versioning convention, backward-compatibility commitments, and a matrix or commitment statement.

### Q-30. SOC-report vs FFIEC-examination evidence overlap — sequencing guidance?

audit-procedures.md and the regulator pack describe both SOC engagements and FFIEC examinations. The SOC engagement happens during the period (typically annual); the FFIEC examination happens on the regulator's cycle (typically 18-month for community banks, more frequent for tier-1). Both consume the same evidence base (operational events, verifier output, IKM roster, etc.). For institutions undergoing both within the same window: is there sequencing guidance — does the SOC report's findings inform the FFIEC examination, or vice versa? Does the institution share the SOC report's working papers with the FFIEC examiner? Round-13 should consider adding short guidance on SOC-FFIEC evidence-handoff: how the institution provides one body of evidence to two consumers, what cross-references are appropriate, what timing dependencies exist. Without it, institutions duplicate evidence-production work.

### Q-31. CSF 2.0 mapping vs FFIEC IT Handbook mapping — parity analysis?

The chain has two parallel mappings: FFIEC IT Handbook (primary: II.C.10 Logging in the IS booklet) and NIST CSF 2.0 (primary: PR.DS-06 Information Integrity). Both mappings exist as separate documents (`docs/regulator-pack/handbook-mapping.md` and `docs/regulator-pack/CSF-2.0.md`) with overlapping subcategories. For institutions subject to both frameworks, do the mappings agree on which controls are headline vs supporting? Are there subtle disagreements (e.g., a control listed as primary in one and supporting in the other)? Round-13 should add a short parity analysis confirming the two mappings are coherent — institutions whose examiners work both frameworks can be sure the same control evidence satisfies both equally. If divergences exist, name them explicitly so the institution's compliance team can prepare both responses.

### Q-32. CUEC tiering rationale — gating criterion?

CUECs are tiered into Critical (Tier 1, day-1) / Important (Tier 2, first 6 months) / Operational hygiene (Tier 3, first year). The doc says "without these, the chain's claims are weaker than the spec implies" for Tier 1. But the tiering isn't explicitly tied to the spec's normative MUST/SHOULD/MAY lattice — some Tier 1 CUECs are spec-mandated (CUEC-IAM-01..05 RBAC defense-in-depth maps to spec §10.3 SHOULD; CUEC-CRY-03 IKM custody maps to spec §10.6 MUST), and some are spec-recommended (CUEC-CRY-04 weekly fingerprint reconciliation maps to spec §10.1 SHOULD). Round-13 should consider re-tiering CUECs against the spec's MUST/SHOULD/MAY lattice so institutions know which CUECs are mandatory for v1.0 conformance vs which are recommended-but-deferrable. The current tiering is "production-readiness staging," which is useful but distinct from "conformance-required."

### Q-33. Mirror-registry pattern — explicit project recommendation?

`docs/supply-chain.md` §"Container supply chain" describes two mirror-registry patterns: signature-preserving (preferred) and re-signing (with documented bridging). The re-signing pattern requires substantial controls — audit log on WORM storage with 7-year retention, continuous monitoring (`mirror.reconciliation_completed`), bridge invariant documentation, IR Scenario 13. For institutions deciding which pattern to adopt at procurement time, the cost-of-control delta is material. Signature-preserving is described as "preferred"; round-13 could promote this to an explicit recommendation: "Signature-preserving is the spec's recommended pattern for v1.0 deployments; institutions adopting re-signing operate the documented bridging controls and accept the additional CC8.1 control burden." Without a clear recommendation, institutions default to whichever is operationally simpler at the moment, which is sometimes re-signing without the bridging controls — a non-conformant posture.

### Q-34. MRM model-deployment-intent record schema?

`MRM-COMMITTEE-BRIEF.md` asks the committee chair to ask "the institution's deployment-intent record per `gen_ai.response.model` value observed in the period" — distinguishing between A/B test, canary deployment, vendor silent re-routing, and multi-region deployment with version drift. The four scenarios have different MRM dispositions (A/B is deliberate model-validation activity; vendor re-routing is a control-completeness gap; multi-region drift is operational housekeeping; canary is bounded production-validation). But there's no documented schema for the deployment-intent record. Round-13 could add a schema (or extend `audit.deployment.*` namespace under spec §4.4) so the institution's record-keeping is consistent and the SOC team can mechanically test it via P-26 (model-inventory composition cross-check) extended with deployment-intent. This makes the per-model decision-count distribution analysis substantive rather than circumstantial.

### Q-35. SLSA L3 institution-side validation — CUEC formalization?

Supply-chain.md's "What examiners verify" list includes "the institution validates the SLSA provenance attestation against `slsa-verifier` at deployment time (deployment-gating control); the SLSA-verifier output log is archived for the binary's deployment lifetime (typically 7 years)." This is a substantial new institution-side control burden, equivalent in maturity to cosign + GPG validation discipline but requiring a third trust anchor (`slsa-verifier` public key). Is round-13 going to formalize this as a new CUEC (CUEC-VER-05 or CUEC-VND-05) or is it absorbed under existing CUECs? Without explicit naming, examiners may not test it consistently and institutions may operate it inconsistently. The deployment-gating posture is the institution's primary defense against R5 (vendor supply chain compromise); explicit naming raises its visibility.

### Q-36. `mirror.reconciliation_completed` recommended schema?

Supply-chain.md introduces `mirror.reconciliation_completed` as an institution-defined operational event paralleling spec §10.1 fingerprint reconciliation but for the mirror's signature-validation discipline. "Institution-defined" means schema discretion is the institution's. For mechanical SOC consumption (P-N audit-procedure for the bridge invariant), a normative or RECOMMENDED schema would let the SOC team test consistently across institutions running the re-signing pattern. Round-13 could add the schema to either `docs/soc-pack/control-evidence-events.md` or as an annex to supply-chain.md. Suggested fields: `event`, `timestamp`, `tenant_id`, `correlation_id`, `fields: {pulled_image_count, validated_image_count, unvalidated_image_count, audit_log_uri, reconciliation_window}`.

### Q-37. GOVERN function — institution risk-tolerance statement template?

CSF 2.0's GOVERN function (top-level in CSF 2.0) expects integrity controls to flow from risk-tolerance statements down to operational evidence. CSF-2.0.md §"GOVERN function alignment (one-page institution adoption)" describes what the institution's risk-tolerance statement should name (per-day Merkle seal as load-bearing tampering-detection, per-event HMAC chain as load-bearing real-time-tampering-detection, regulator-held public key as integrity trust anchor). Round-13 could add a template or example risk-tolerance statement that institutions can adopt with minimal modification. This lowers the adoption bar for institutions that haven't written this statement yet — currently, the institution is on its own to translate the chain's primitives into risk-tolerance language consumable by the institution's policy framework. A template both accelerates adoption AND surfaces consistency across institutions for cross-institution examiner comparison.

### Q-38. EU AI Act Article 26(6) retention vs spec defaults — multi-jurisdiction coordination?

EU AI Act Article 26(6) requires logs retained at least 6 months unless longer required by Union or national law. Spec defaults and operator-guide assume 7-year retention (typical US FFIEC banking norm). For institutions deploying in EU, US-style 7-year retention is over-compliant with the EU minimum but consistent with US standards. For institutions deploying in jurisdictions with neither the US 7-year norm nor explicit EU long-retention law, what's the recommended posture? Round-13 could add to ai-policy-alignment.md a multi-jurisdiction retention coordination short section: "retention is the longer of (a) US FFIEC 7-year if applicable, (b) Article 26(6) 6-month, (c) institution's commercial retention policy, (d) any other applicable jurisdictional minimum; institutions document the retention rationale in their control description." This is a presentation question affecting institutions making the adoption decision in multi-region operations.

### Q-39. DORA Article 17 24-hour vs FFIEC 36-hour incident reporting — IR playbook clock-start jurisdiction handling?

DORA Article 17 requires 24-hour reporting for major ICT incidents. FFIEC computer-security incident notification rule is 36 hours. The IR playbook references "36 hours per the cyber-incident notification rule" but doesn't explicitly handle the 24-hour DORA window for EU-supervised institutions. For multinational institutions subject to both rules, the tighter window dominates per incident; the institution's IR playbook should articulate which clock applies under which jurisdictional incident-class. Round-13 could update incident-response-playbook.md to add: "For institutions with multi-jurisdiction supervision, the IR playbook's clock-start triggers apply per the tightest applicable rule. DORA Art. 17 24-hour clock applies to institutions with EU-supervised entities; FFIEC 36-hour clock applies to institutions with US-federal-supervised entities; institutions with both operate the 24-hour clock as the binding floor." Without this clarification, multinational institutions face an institutional-discretion question with substantial regulatory consequences.

### Q-40. Institution-side AI policy framework — template artifact?

ai-policy-alignment.md describes which controls the chain satisfies under each framework (NIST AI RMF, EO 14110, Treasury AI RMF, CFPB, OCC, EU AI Act, EBA, DORA, NIS2). It does NOT provide a template institution-side AI policy framework. Adoption institutions face the question: "what does our internal AI policy framework look like that references the chain as one technical control among many?" Round-13 could add a short template institution-side AI policy framework that institutions adapt — at minimum, a one-page outline pointing to AI-RMF / EO 14110 / Treasury / SR 11-7 / EU AI Act / DORA / NIS2 sections and naming the chain's role within each. This complements the MRM-COMMITTEE-BRIEF (operational) and the management-summary (executive) with a policy-level artifact. Specific institutions write better policy than templates produce, but a starting template accelerates the conversation between IT-Risk, Legal, and Compliance during chain adoption — especially institutions where AI policy doesn't yet exist or doesn't yet reference logging-integrity controls.

---

## Stopping criterion verdict

**1 Partial (the routing question, F-1) plus 40 supporting questions Q-1 through Q-40.**

The routing question is settled when round-13 takes a position. The verdict assumes round-13 will respond; if it doesn't, F-1 becomes a Gap in the round-13 close-out.

The 40 supporting questions are unclassified intentionally — round-13 reviewers triage each as Already Addressed, Adopt as Partial, Adopt as Gap, or Out of Scope. The intent is to broaden the round-13 review surface so the spec at FFIEC submission has answered the questions a production-multi-LLM operator, a tier-1 institution, a cross-functional bank operations team, a CSF/SOC compliance team, and a multinational deployer subject to both U.S. and EU supervision would ask, not just the questions the regular six personas raised.

**Question groupings.** Q-1 through Q-10 surfaced from the production-multi-LLM-operator lens (routing, mock LLM, edge cases of multi-process flows). Q-11 through Q-15 surfaced from a deeper read of the v1.0-final spec text and the Go reference implementation; they are spec-text and audit-procedure questions. Q-16 through Q-30 surfaced from the cross-functional read pass (privacy-by-design, legal-disclosure, customer-dispute, BYOC, vendor-hosted, cost-model, at-scale-operations); they are integration questions across the corpus and audience tracks. Q-31 through Q-37 surfaced from the controls-mapping and supply-chain-detail read pass (CSF 2.0, CUECs, mirror-registry handling, MRM brief, SLSA L3 institution-side validation); they are controls-framework and adoption-readiness questions. Q-38 through Q-40 surfaced from the AI-policy-alignment read pass (EU AI Act, DORA, multi-jurisdiction); they are international-deployment and policy-framework questions.

**Disposition expectation.** Most of the questions are likely Partials closing with short text additions or Already-Addressed with cross-reference clarifications. A handful may rise to Gaps if round-13 takes positions opposite to the corpus's current tone — F-1 (routing) is the most likely Gap; Q-17 (Pattern B conformance) and Q-23 (vendor-hosted exit procedure) are the next most likely. The remainder cluster as presentation-quality, integration-quality, and forward-commitment items.

**Cumulative submission-quality posture.** Closing F-1 plus Q-1 through Q-30 by round-13 close-out positions the spec for FFIEC submission in October as a corpus that has been substantively cross-examined from outside the regular six-persona review cycle, with each integration question taking an explicit position (or explicit out-of-scope with reasoning). The submission package's strength is the explicit positioning, not the absence of questions — every spec ships with open questions; mature specs ship with documented positions on each.
