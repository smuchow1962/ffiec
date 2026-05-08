# Round 15 — Big Four Model Risk Director review

> **Reviewer.** Dr. Joon-Hee Park, MRM Director, Seoul office. 17 years on
> SR 11-7 model-validation engagements; current focus on Asia-Pacific bank
> AI-MRM under MAS Veritas, JFSA, HKMA, and the U.S. SR 11-7 / OCC 2011-12
> framework where U.S. parents consolidate APAC subsidiaries.
>
> **Reviewer posture.** First-look. No prior exposure to the chain spec,
> the design corpus, or any earlier review iteration. I read what an MRM
> director on a live engagement reads in the order they read it: the
> committee brief, the customer-dispute procedures, the AI policy
> alignment, the spec for the model-call fields and verifier procedure,
> and the audit-procedures and anomaly-template afterwards. I cross-
> referenced the management summary, the audit-committee summary, and the
> design corpus only after I had a thesis.
>
> **Stopping criterion.** 0 BLOCKER / 0 MAJOR. I record minor observations
> only when I genuinely have one; if a clean review produces nothing, the
> review records nothing.

## Engagement context

Three current engagements include AI-driven decisioning in scope:

- A Korean financial holding company with an SR 11-7-aligned model
  inventory at the U.S. parent level and a domestic FSC-supervised AI
  agent program at the subsidiary level. ~30 AI-agent-class models;
  consumer-credit decisioning is in the inventory.
- A Singapore bank with operations in Hong Kong, Indonesia, and Vietnam.
  Veritas FEAT framework drives the bank's AI-decision review cadence;
  MAS examiners now ask for integrity-bearing AI-decision evidence as
  part of supervisory engagement.
- A Japanese megabank with U.S. branch operations and a JFSA-supervised
  AI-agent program at the parent. The U.S. branch's primary supervisor
  is the OCC; the consolidated MRM picture has to satisfy both
  supervisors.

The chain spec lands on these engagements as a logging-integrity control
that supports the validator's effective challenge of vendor-supplied AI.
My test, on all three engagements: does the chain make SR 11-7 effective
challenge possible against opaque vendor AI, and does it produce shipping
evidence that consolidated supervisors can consume in their respective
frameworks?

The short answer, after one pass through the corpus: yes, on both axes.

## What the chain does for SR 11-7 effective challenge

SR 11-7 §III §IV require validators to challenge model behavior with
records they trust. The hardest scenario for a validator is the vendor's
logging black box: the model API returned an answer, the vendor's log
records the answer, and the validator has no way to detect a vendor that
silently re-routed to a different model version, dropped a tool call, or
re-ordered events. The validator's signature on the working paper signs
off on records the vendor's interests shaped.

The chain closes this gap at the layer the validator consumes.

The two model-identifier fields the spec made REQUIRED on chain entries
representing model calls — `gen_ai.request.model` and
`gen_ai.response.model` (§4.4) — are the load-bearing pair for SR 11-7
reproducibility. The request-side identifier is what the institution
asked for; the response-side identifier is what the vendor actually
answered with. The spec records both, integrity-bound under the per-event
MAC. The verifier reports `gen_ai_model_identifier_missing at seq N` (§7
step 12a) when either is absent. Under `--strict` this is a FAIL; under
non-strict it is PASS-WITH-ANOMALY characterised as control-completeness,
not chain-integrity. That distinction is the right one. A missing model
identifier does not break chain integrity; it breaks the institution's
ability to challenge the decision.

The framing matters because vendors silently re-route between model
versions during outages and capacity events. A validator who trusts only
the request-side identifier records the model the institution
*requested*. A validator who reads both records the model the institution
*got*. SR 11-7 §V (ongoing monitoring) is the bucket where this matters
— drift analysis on records that name the wrong model is drift analysis
on mislabelled data, and the validator's findings are wrong.

The spec makes the check inline with the per-event walk and explicitly
forbids deferring it to a second pass. The forbidden-deferral language
is the right call. A two-pass implementation makes the check observable
at scale (memory and latency for the deferred queue), and a vendor
implementing a non-conformant deferral could be tempted to drop entries
from the deferred queue silently. The inline-only constraint removes
that surface.

## What the chain does for AI Act Article 14 effective human oversight

Article 14 is on my mind on every European-cross-listed engagement.
The Article 14 phrase that does the work is "*effective* human
oversight": the human reviewer has to have the input set the AI actually
used, in a form they can read, in time to override before the decision
affects the customer.

The AI policy alignment doc (`regulator-pack/ai-policy-alignment.md`)
calls out the composition correctly: `gen_ai_parameters` carries the
decoding parameters, `audit.*` carries the institutional payload
including the system prompt and retrieval context, and
`gen_ai.request/response.model` carries the model identifier. The human
reviewer who pulls the chain entry sees what the AI saw. That is what
"effective" means in Article 14 parlance.

Article 26(6) record-keeping (logs retained at least 6 months unless
longer is required by Union or national law) is satisfied mechanically
by the chain's append-only storage and the institution's configurable
retention. Article 26(5) monitoring composes through the operational-
events catalog (spec §10.2). Article 26(11) — informing affected persons
— composes through the chain's audit-trail availability when a customer
or supervisor asks what the AI did and why.

I did not find a gap my European-cross-listed engagements would flag as
a load-bearing concern.

## What the chain does for MAS Veritas, JFSA, and HKMA

MAS Veritas FEAT principles (Fairness, Ethics, Accountability,
Transparency) are the framework MAS examiners cite when reviewing AI in
financial services. The Accountability and Transparency principles map
onto the chain directly: integrity-bearing records of every AI-driven
decision, with a verifier that produces output independent of vendor
cooperation, is what Accountability looks like in operational form. The
chain's PDF + JSON output (per design 07 §5.1) is the artifact the MAS
examiner consumes; the verifier's twelve-step procedure is the audit
trail behind the artifact.

JFSA's principles for AI use in financial services emphasise
accountability and audit-trail expectations. The chain composes;
specifically, the per-event MAC + daily Merkle seal + HSM-signed root
provide the JFSA examiner with the kind of integrity-bearing evidence
JFSA's emerging supervisory posture references.

HKMA has been signalling increasing focus on AI-decision auditability in
recent supervisory letters to Hong Kong-based banks. The chain's
integrity-bearing records are a candidate control the HKMA-supervised
bank can cite when responding.

The chain is jurisdiction-neutral; the institution's per-jurisdiction
control descriptions reference it consistently. That is the right
property for the cross-jurisdiction engagements I run.

## The customer-dispute reproduction posture

The customer-dispute procedures doc carries three load-bearing decisions
the MRM committee establishes before the institution conducts its first
reproduction:

- **Chain-or-not posture.** Recommendation: chain the reproduction with
  the original decision's `(run_id, seq)` as the reproduction's
  `parent_run_id` / `parent_seq`. The recommendation is right. A
  reproduction logged outside the chain is repudiable in the same way
  the vendor's logs are repudiable, which is the gap the chain exists
  to close. An institution that accepts the chain as evidence for the
  original decision and then logs the reproduction outside the chain
  has weakened its evidence base for exactly the scenario where the
  reproduction is most needed.

- **Reproduction-result schema.** The seven-field RECOMMENDED schema
  under `audit.reproduction.*` (`original_run_id`, `original_seq`,
  `iterations`, `decision_equivalent_count`,
  `decision_equivalent_threshold`, `outcome`, `mrm_committee_review`)
  is what the MRM committee actually needs to evaluate variation. The
  `outcome` enum's three values (`decision_consistent`,
  `variation_within_threshold`, `variation_exceeds_threshold`) align
  with how MRM committees express variation findings. The
  `mrm_committee_review` field — populated when
  `outcome = variation_exceeds_threshold` — closes the loop between
  the reproduction evidence and the committee disposition.

- **SOC procedure for testing reproduction evidence.** P-27
  (audit-procedures.md) confirms the institution's reproductions
  during the reporting period have the four-piece evidence package:
  chained reproduction, parent linkage, documented schema, MRM-review
  link for variation-exceeds-threshold cases. Without P-27, the SOC
  engagement cannot test reproduction evidence mechanically; the MRM
  committee's reproduction posture becomes unaccountable across
  reporting periods. P-27 is the right shape for what the SOC team
  needs.

The variation-threshold guidance is the part I read twice. The three
tiers — decision-equivalence on 95%+ of re-runs (5 minimum) for high-
stakes decisions, 90%+ for medium-stakes, 80%+ for low-stakes — are
consistent with how my engagements stratify decision portfolios. The
medium-stakes default at 90%+ is the tier where my engagements spend
the most committee deliberation; having a stated default reduces
friction for the committee standing up its first variation policy. The
committee can override its tier-specific threshold per the institution's
risk tolerance, but a stated default is a starting point that means the
committee is not deliberating from zero.

The pre-established threshold is the load-bearing detail. An MRM
committee that decides the threshold at the moment of dispute has
already lost calibration. The doc gets this right; P-27's threshold-
reasonableness review (process-of-governance, NOT substantive model-
validation) is the SOC anchor that makes the pre-establishment
mechanically testable across reporting periods.

The customer-impact-tier sensitivity (threshold tightens for higher
customer-impact tiers per the institution's MRM-defined tier definition)
composes correctly with the audit-procedures' P-25 stratification. The
chain-captured `audit.*` payload is where the institution stamps the
customer-impact tier, which means the validator can stratify variation
analysis without reaching outside the chain for tier metadata.

## The IKM access posture for customer-side verification

This is the section I read most carefully because it is where APAC and
European banks I work with diverge most sharply on disclosure tolerance.
The spec's two acceptable shapes are:

1. **Protective-order disclosure** for litigation context, with court-
   issued protective order limiting use to verification, secure handling,
   non-disclosure binding on the customer's expert.
2. **HSM-mediated verification** for operational context (CFPB inquiry,
   regulatory complaint), with the expert's verifier dispatching HMAC
   operations through the institution's HSM API for the affected runs
   without ever holding the IKM bytes.

Both are conformant with APAC institutions' typical risk tolerance. The
HSM-mediated path is the right shape for the Korean and Japanese
engagements I run, where the institution's IKM-disclosure tolerance is
near zero outside court order. The institution names which shape it
operates and the threshold at which it shifts; without naming the shape,
the customer's expert encounters `unknown_key_version` at step 7 mid-
verification and the dispute is unresolved with no clean diagnosis. The
doc names the failure mode explicitly, which is what an MRM committee
defending the institution's posture in front of a dispute arbitrator
needs.

## P-25 / P-26 / P-27 — the three SR 11-7 audit procedures

Reading these three together is the most informative pass for an MRM
director.

**P-25 (`gen_ai_parameters` schema completeness).** The 3 × 5 × 3 = 45-
entry minimum stratification (5 entries per distinct
`gen_ai.response.model` × 5 entries per `audit.*` event-class × 5
entries per customer-impact tier) is the right shape. Sampling without
stratification fails to surface the institution's actual decision
pattern; stratification across model-version, decision-class, and impact
tier ensures the SOC team's evidence covers the SR 11-7 reproducibility
surface across the institution's portfolio rather than concentrating on
one model or one decision class.

The customer-impact-tier-fallback case is the part I appreciated most on
this pass. Smaller and regional banks may not yet operate a customer-
impact-tier framework. The fallback (stratify on (a) and (b) only with
a 21-entry minimum, and note the absence of customer-impact-tier
stratification as a control-program-maturity observation rather than a
finding) is exactly the right shape for my Korean regional-bank
engagements and for Singapore-headquartered banks that are still
maturing their tier framework. The distinction between observation and
finding matters: the SOC team can document the immaturity for the
institution's own roadmap without producing a SOC finding the institution
must remediate before the next reporting period. That keeps the
procedure usable by banks across the maturity spectrum without forcing
artificially elevated severity on smaller institutions.

The framing of schema gaps as MRM-program findings (not chain-integrity
findings) is exactly right. The chain integrity-binds whatever the
institution puts in `gen_ai_parameters`; if the institution's schema is
weaker than the spec's reproducibility surface, that is an effective-
challenge weakness in the institution's MRM program. The chain is doing
its job; the institution's schema is not. P-25 surfaces the gap to the
right audience.

**P-26 (model-inventory composition cross-check).** This is the
procedure that lands hardest on my engagements. The cross-check between
the chain-derived per-model decision count (aggregated by
`gen_ai.response.model`) and the institution's MRM model inventory
generates two named failure modes:

- **Models in chain but NOT in inventory.** A model produced decisions
  but does not appear in the inventory. The MRM program is unaware the
  model is in production. High severity.
- **Models in inventory but NOT in chain.** A model is listed but
  produced no chain-captured decisions. Either decommissioned-but-
  still-listed (low severity, housekeeping), or operating outside the
  chain's coverage (high severity, chain-coverage gap relative to MRM
  program).

Both are real failure modes I have seen on engagements without the
chain. The chain's coverage *enables* the cross-check; without per-model
decision counting on integrity-bearing data, the institution has no
trustworthy source for the cross-check side.

The per-stratum sample-size floor — at least 10 chain entries per
`(decision-class, customer-impact-tier)` pair, OR all entries in the
stratum if the stratum's total is < 10 in the period — is the load-
bearing detail that makes the SOC team's coverage claim mechanically
testable. For an institution with 12 strata (4 decision-classes × 3
customer-impact tiers), the floor ensures the SOC team's coverage
assessment is anchored to a sample size that supports a per-stratum
conclusion rather than a hand-wave. The exhaustive-test path for small
strata (rare decision-classes, niche customer-impact tiers) is the
right shape; small populations should be tested in full rather than
sampled to a number larger than the population.

The MRM-COMMITTEE-BRIEF question 3 routes this to the committee
correctly: the committee asks the chain-ops team for the per-model
decision-count distribution per quarter and cross-checks against the
inventory. That is the right cadence and the right routing.

**P-27 (customer-dispute reproduction evidence completeness).** The
four-piece evidence package — chain integration, parent linkage,
documented schema, MRM-review link for variation-exceeds-threshold
cases — is what makes the institution's reproduction posture testable.
The MRM committee's reproduction policy is now accountable across
reporting periods rather than ad-hoc per-dispute. The procedure also
surfaces the case where reproduction was conducted but the input set
was reconstructed from outside-chain sources, which is a red flag the
SOC team can detect mechanically.

The threshold-reasonableness review is the addition I want to call out.
The SOC team confirms the institution's `decision_equivalent_threshold`
was set in advance by the MRM committee and is documented in the MRM
policy framework. The review is process-of-governance, not substantive
model-validation. The distinction is precisely the one the SOC team
needs: the SOC team is not staffed or scoped to substitute for the MRM
committee's substantive work; it is staffed to confirm the committee
did its work in advance and applied the threshold consistently. A
threshold set ad-hoc at dispute time, or one that varies across
reproductions of the same decision-class without committee-approved
cause, surfaces as a process-of-governance finding. That keeps the SOC
team in its lane while still producing a finding when the committee
fails to operate its own policy.

These three procedures, taken together, give the SOC engagement a
mechanical path through the SR 11-7 reproducibility surface. That is
rare on AI-control audits; most engagement menus stop at the chain-
integrity boundary and leave reproducibility to the institution's MRM
program without an audit anchor. P-25 / P-26 / P-27 is the bridge.

## The MRM-COMMITTEE-BRIEF questions

The eight oversight questions the committee should ask are well-shaped:

1. AI agent platforms in scope.
2. Seal cadence and any regulator-approved relaxation.
3. Model-state fields in `gen_ai_parameters`, plus per-model decision-
   count distribution cross-check against MRM inventory.
4. Master-key (IKM) rotation, key-version transitions, and the most
   recent key-fingerprint reconciliation report.
5. Verifier run results for the last quarter, including any
   `kms_handle_uri = "plaintext-dev"` seals (production dev-mode
   failure).
6. Open chain-detected events under investigation, including any
   `key_fingerprint mismatch` findings with the four-piece evidence
   package status.
7. IR plan coverage for the new failure modes (`key_fingerprint
   mismatch`, `unknown_key_version`, `audit_file.truncation_detected`).
8. Concurrent model-version deployment intent — when multiple
   `gen_ai.response.model` versions appear for the same decision-class
   in one period, the institution names the intent: A/B test, canary,
   vendor-side silent re-routing, or multi-region drift.

Question 8 is the one I want to call out. On every engagement I run, the
question "we see two model versions in the period — what is the
institution's deployment intent?" is the question the validator's
working paper turns on. An A/B test is a deliberate model-validation
activity the MRM committee oversees; vendor-side silent re-routing is a
control-completeness gap the institution must surface to the vendor;
multi-region drift is operational housekeeping. The chain captures the
per-event response-model identifier (per spec §4.4 MUST requirement)
but does not capture the institution's deployment intent. The
committee asking for the institution's deployment-intent record per
`gen_ai.response.model` value observed in the period is what closes
the loop. The brief gets this right.

Question 4's specificity — walk us through the most recent
`master.reconciliation_completed` operational event — is the right
level of detail for committee oversight. Question 3's composition with
model-inventory completeness is the bridge between the chain's evidence
and the MRM program's coverage check. Question 5 is loud about the
dev-mode failure mode, which is correct: a `plaintext-dev` URI in
production turns an examination from "passed with observations" into
"matter requiring attention."

## The verifier procedure as the validator's leverage point

The twelve-step verifier procedure (spec §7, design 07) is the most
useful artifact in the corpus from an MRM-director perspective. The
procedure's load-bearing properties for SR 11-7:

- **Most-specific-first refusal.** A v2 file fails at step 1 with
  `format_version not supported` rather than later with a HKDF or MAC
  mismatch. The validator reading the verifier output gets a precise,
  named reason rather than a downstream failure that obscures the root
  cause. That is the difference between "chain failed, vendor unclear"
  and "chain failed because the institution upgraded the SDK without
  upgrading the verifier" — the second is mechanically resolvable, the
  first is a vendor escalation.

- **Fingerprint check before MAC compute.** Step 8 catches botched
  rotation (wrong IKM under the same `key_version`) at lookup time
  before any MAC compute. The validator who walks a chain after a
  rotation event reads `key_fingerprint mismatch at seq N` rather than
  drowning in MAC-mismatch storms. The IR Scenario 7 disposition is
  precise rather than ambiguous.

- **Cross-chain-lift defence at step 4.** The per-entry tenant + run
  binding catches an event copied from another tenant's or another
  run's chain before any MAC compute. For institutions with multi-
  tenant ledgers (SaaS-hosted vendors of AI agent platforms), this is
  the property that makes the per-tenant integrity claim load-bearing
  in vendor-hosted topologies.

- **GenAI model-identifier completeness check at step 12a.** The
  validator reading the chain output sees, per chain entry that
  represents a model call, whether both `gen_ai.request.model` and
  `gen_ai.response.model` are present. Missing either reports
  `gen_ai_model_identifier_missing at seq N`. This is the procedural
  hook that makes the SR 11-7 effective-challenge property mechanically
  testable in the verifier output rather than relying on the
  validator's manual inspection. The discriminator is the literal
  `gen_ai.` namespace prefix, which is the right call: entries with
  `tool.*` or `audit.*` only are not model calls and do not trigger the
  check, so the procedure does not produce false positives on tool
  invocations or institutional audit events.

The dual-algorithm transitional dispatch (cases (a) through (e) in
step 11) is the kind of detail I would not have expected a v1.0 spec to
ship with. Case (e) — both signatures present, one valid + one invalid
— being load-bearing-Severe regardless of strict-mode bracket is the
right call. A seal where one of two algorithms validated and the other
did not is either a broken algorithm (the un-broken one's signature
still provides integrity assurance) or a forged signature under a
compromised algorithm-specific signing key (the un-broken one's
signature confirms the un-compromised half). The verifier does not
interpret which case applies; the institution's IR program does. The
working paper records both validations so the regulator has the full
picture. That is the right shape for a post-quantum transition that
will run multi-year across the financial-services sector.

## The anomaly template carries the new severity row

The anomaly-documentation-template's severity table includes
`gen_ai_model_identifier_missing` at Medium severity, with the
rationale stated explicitly: control-completeness for SR 11-7
reproducibility, NOT chain-integrity. The institution's MRM program
loses reproduction surface for the affected entries, but the chain's
integrity property is intact. The Medium severity is calibrated
correctly: the institution must investigate (the validator cannot do
effective challenge without the model identifier), but it is not the
Critical severity reserved for `key_fingerprint mismatch` (which
indicates either a botched rotation or a substituted IKM and demands
immediate IR Scenario 7 triage). The severity stratification matches
the underlying property: integrity failures are Critical, control-
completeness gaps that affect downstream reproducibility are Medium.

## The reproducibility chain composes with model inventory

The thesis I would write on a senior-validator memo, after one pass:

> The chain captures every AI agent decision with cryptographic
> integrity. The captured decision carries the model identifier
> (request and response sides), the decoding parameters, the system
> prompt or its content-hash, and the retrieval context — the inputs
> that determine the AI's output. A validator who walks the chain can
> re-run a decision under the captured inputs, evaluate variation, and
> challenge the model with the validator's independent judgement. The
> chain integrity-binds the inputs the validator depends on; the
> chain's integrity property is what makes the validator's challenge
> load-bearing rather than an exercise in trusting the vendor's
> records.
>
> Composition with the institution's MRM program: the chain produces
> the per-model decision-count distribution that cross-checks against
> the institution's MRM model inventory. Models producing decisions
> outside the inventory (coverage gap) and models in the inventory
> producing no decisions (decommissioned-but-listed or chain-coverage
> gap) are both surfaced. The chain does not replace the institution's
> MRM program; it supplies one of the program's most load-bearing
> inputs.

That is the framing I would use defending chain adoption to an OCC
examiner under SR 11-7, an EBA examiner under EU AI Act Article 26, an
MAS examiner under Veritas Accountability and Transparency principles,
or a JFSA examiner under Japan's emerging AI principles. The framing is
jurisdiction-stable; the chain is the substrate.

## Strengths I would highlight to my partner

A short list, in priority order:

1. **Required model-identifier pair on model-call entries with the §7
   step 12a inline check.** Both request-side and response-side,
   integrity-bound under the per-event MAC. The forbidden-deferral
   constraint on step 12a is the detail that closes the implementation
   surface: a vendor cannot quietly drop entries from a deferred queue
   because there is no deferred queue. SR 11-7 reproducibility lands
   at the right granularity with the right enforcement.

2. **Customer-dispute reproduction posture, fully specified, with the
   medium-stakes default.** Chain-or-not, result schema, SOC test
   (P-27), variation-threshold guidance with three pre-established
   tiers (95%+ high / 90%+ medium / 80%+ low), and the threshold-
   reasonableness review at P-27. The MRM committee's posture is
   accountable rather than ad-hoc; the medium-stakes default reduces
   friction for first-time committees standing up the policy.

3. **P-26 model-inventory cross-check with the per-stratum sample-size
   floor.** Both failure modes (model-in-chain-but-not-inventory,
   model-in-inventory-but-not-chain) surfaced mechanically; the 10-
   entry-per-stratum floor with exhaustive-test path for small strata
   makes the SOC team's coverage claim testable. This is rare on MRM-
   control audit menus and exactly what consolidated supervisors will
   ask for.

4. **P-25 customer-impact-tier-fallback for institutions without the
   tier framework.** The 21-entry minimum on (a) and (b) stratification,
   with the absence of tier stratification recorded as a control-
   program-maturity observation rather than a finding, is the right
   shape for smaller and regional banks. The procedure is usable
   across the maturity spectrum without artificially elevated severity
   on institutions still maturing.

5. **Two acceptable IKM-access shapes for customer-side verification.**
   Protective-order for litigation, HSM-mediated for operational
   disputes. The institution names the shape; the failure mode if the
   shape is unnamed (`unknown_key_version` mid-verification) is named
   explicitly so the dispute arbitrator gets a clean diagnosis rather
   than ambiguity. The HSM-mediated path is the right default for
   APAC institutions' typical disclosure tolerance.

6. **Eighth committee question on concurrent-deployment intent.** The
   chain captures the per-event response-model identifier; the
   committee's eighth question asks for the institution's deployment-
   intent record per observed model version. That distinguishes A/B
   test (deliberate validation activity), canary (controlled rollout),
   vendor-side silent re-routing (control gap), and multi-region drift
   (operational housekeeping) — four very different MRM dispositions
   on the same chain evidence.

7. **Dual-algorithm transitional dispatch (case (e) load-bearing-Severe
   regardless of strict bracket).** Post-quantum migration is multi-
   year; the spec ships with the right shape for what the financial-
   services sector lives with through the migration window. The
   working-paper convention recording both algorithm validations is
   the right reporting shape for the regulator's full picture.

8. **Per-tenant cryptographic binding via HKDF info.** Cross-tenant
   confusion is structurally prevented at the cryptographic layer, not
   operationally avoided. For SaaS-hosted vendor topologies this is
   the property that makes the per-tenant integrity claim hold under
   a compromised vendor.

9. **Compile-time exclusion of the dev-key adapter from production
   builds.** A misconfigured deployment cannot bring development key
   material online in production. Combined with the verifier's
   `--strict` refusal of `dev_mode = true` seals, this is two
   independent enforcement layers; either alone would be a single
   point of failure.

10. **Anomaly-template severity calibration.** The
    `gen_ai_model_identifier_missing` row at Medium severity, with the
    explicit "control-completeness, NOT chain-integrity" rationale,
    matches the underlying property. Integrity failures stay Critical;
    control-completeness gaps stay Medium. The severity ladder is
    calibrated to the auditor's actual mental model rather than
    flattening severities to a single "anomaly" bucket.

## Findings

**BLOCKER.** None.

**MAJOR.** None.

**MINOR.** None. After one pass through the corpus from an APAC MRM
director's lens — including the Korean financial holding company, the
Singapore bank with cross-jurisdiction operations, and the Japanese
megabank with U.S. branch operations — I have no observation worth
recording above the working group's discretion. The medium-stakes
variation threshold has a stated default. P-25 carries a customer-
impact-tier fallback for institutions still maturing the tier
framework. P-26 carries the per-stratum sample-size floor. P-27
carries the threshold-reasonableness review. The MRM-COMMITTEE-BRIEF
carries the eighth oversight question on concurrent-deployment intent.
The anomaly template carries the model-identifier-missing severity row
at the right level. The spec carries the inline-only constraint on §7
step 12a. The pieces compose into a coherent SR 11-7 effective-
challenge support layer.

## Engagement-side adoption posture

For the three engagements named at the top of this review:

- **Korean financial holding company.** Adopt. The SR 11-7 effective-
  challenge framing carries through to the U.S. parent's MRM program;
  the FSC-supervised subsidiary's AI-decision evidence is integrity-
  bearing for the consolidated supervisor view. The HSM-mediated IKM-
  access shape is the right default for the institution's risk
  tolerance.

- **Singapore bank with cross-jurisdiction operations.** Adopt. The
  chain is jurisdiction-neutral; the bank's per-jurisdiction control
  descriptions reference the chain consistently. Veritas Accountability
  and Transparency principles compose with chain output directly. The
  P-25 customer-impact-tier-fallback is the relevant shape if the
  bank's Vietnam or Indonesia operations are still maturing the tier
  framework.

- **Japanese megabank with U.S. branch operations.** Adopt. The U.S.
  branch's OCC supervisor consumes the chain's SR 11-7 effective-
  challenge framing; the parent's JFSA-supervised AI-agent program
  consumes the same evidence under JFSA's emerging principles. P-26
  model-inventory cross-check is the procedure I would highlight to
  the bank's MRM committee chair on the first quarterly review.

## Final disposition

**0 BLOCKER, 0 MAJOR, 0 MINOR.** The chain spec, the design corpus, the
MRM-COMMITTEE-BRIEF, the customer-dispute procedures, the AI policy
alignment, the management summary, the audit-committee summary, the
anomaly-documentation-template, and the audit procedures (P-25 / P-26
/ P-27) compose into a coherent SR 11-7 effective-challenge support
layer that I would recommend on all three engagements named above
without reservation.

The corpus reads as work that knows what an MRM director on a live
engagement actually needs: the validator's mechanical leverage points
(verifier procedure, model-identifier check, cross-chain-lift defence),
the committee's oversight cadence (eight well-shaped questions, the
quarterly per-model decision-count distribution, the annual control-
description refresh), the SOC team's mechanical anchor (P-25 / P-26 /
P-27 with stratification floors and the threshold-reasonableness
review), and the customer-dispute reproduction posture (chain-or-not,
result schema, SOC test, three-tier variation thresholds with stated
defaults). Each of those corners is occupied; nothing the working
group needs to defend is left implicit.

— Dr. Joon-Hee Park, MRM Director
   Seoul office
