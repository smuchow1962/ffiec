# Round 14 — Big Four Model Risk Director review

> **Reviewer.** Dr. Elena Vasilescu, MRM Director, Bucharest office. 18 years on
> SR 11-7 model-validation engagements; current focus on European AI Act
> implementation for cross-border financial institutions.
>
> **Reviewer posture.** First-look. No prior exposure to the chain spec, the
> design corpus, or earlier review rounds. I read what an MRM director on a
> live engagement would read in the order they would read it — committee
> brief, customer-dispute procedures, AI policy alignment, and then the spec
> itself for the model-call fields and verifier procedure. I cross-checked
> the audit-procedures and anomaly-template only after I had a thesis.
>
> **Stopping criterion.** 0 BLOCKER / 0 MAJOR. I record minor observations
> the working group may close at its discretion; none of them gate adoption
> on my engagements.

## Engagement context

Three of my current engagements include AI-driven decisioning in scope:

- A pan-EU retail bank with high-risk AI under EU AI Act Article 6 and a
  consumer-credit decisioning model that triggers Article 14 effective-human-
  oversight requirements.
- A US bank holding company with an OCC primary supervisor and an SR 11-7
  model inventory north of 200 models, ~25 of which are AI-agent-class.
- A Singapore-headquartered bank with operations in three jurisdictions
  (MAS, BaFin, FCA) and a Veritas-aligned validation framework.

The chain spec lands on these engagements as a logging-integrity control. My
test is the same on all three: does the chain do enough to make SR 11-7
effective challenge possible against vendor AI, and does it carry the
shipping evidence European regulators are starting to ask for under the AI
Act and DORA?

The short answer, after one pass through the corpus: yes, on both axes.

## What the chain does for SR 11-7 effective challenge

SR 11-7 §III §IV require validators to challenge model behavior with
records they trust. The validator's nightmare is the vendor's logging
black box: the model API returned an answer, the vendor's log says it
returned that answer, and the validator has no way to detect a vendor that
re-routed the request to a different model version, dropped a tool call, or
silently re-ordered events. The validator's signature on the validation
working paper is signing off on records the vendor's interests
shaped.

The chain closes this gap at the layer the validator actually consumes.

The two model-identifier fields the spec made REQUIRED on chain entries
representing model calls — `gen_ai.request.model` and
`gen_ai.response.model` (§4.4) — are the load-bearing pair for SR 11-7
reproducibility. The request-side identifier is what the institution asked
for; the response-side identifier is what the vendor actually answered
with. **The spec records both, integrity-bound under the per-event MAC.**
The verifier reports `gen_ai_model_identifier_missing at seq N` (§7 step
12a) when either is absent. Under `--strict` this is a FAIL; under
non-strict it is PASS-WITH-ANOMALY characterised as control-completeness,
not chain-integrity. That distinction is the right one — a missing model
identifier does not break chain integrity, but it breaks the institution's
ability to challenge the decision.

The framing matters because OpenAI and Anthropic both silently re-route
between model versions during outages and capacity events. A validator who
trusts only the request-side identifier records the model the institution
*requested*; a validator who reads both records the model the institution
*got*. SR 11-7 §V (ongoing monitoring) is the bucket where this matters —
drift analysis on records that name the wrong model is drift analysis on
mislabeled data, and the validator's findings are wrong.

The integrity-binding of the OTel envelope (§5: `gen_ai.*`, `tool.*`,
`audit.*`, plus `parent_run_id` / `parent_seq` / `dag_parents`) is the
second load-bearing property. The reproducibility surface the spec
recommends inside `gen_ai_parameters` (§4.4) — decoding parameters,
sampler implementation identifier, system prompt content or hash with
prompt-version registry reference, RAG document IDs and content-hashes,
intra-run data dependencies — is exactly the input set a validator needs
to re-run the decision and challenge the model under the same conditions.
The institution sets the schema; the chain integrity-binds whatever the
institution puts in the field. That is the right division of labour: the
chain does not opine on what the institution must capture; it guarantees
the capture is unrepudiable.

## What the chain does for AI Act Article 14 effective human oversight

Article 14 of the EU AI Act is on my mind in every European engagement.
The Article 14 phrase that does the work is "*effective* human
oversight": the human reviewer has to have the input set the AI actually
used, in a form they can actually read, in time to override the AI's
decision before it affects the customer.

The AI policy alignment doc (`regulator-pack/ai-policy-alignment.md`)
calls out the composition correctly: `gen_ai_parameters` carries the
decoding parameters, `audit.*` carries the institutional payload
including the system prompt and retrieval context, and
`gen_ai.request/response.model` carries the model identifier. The human
reviewer who pulls the chain entry sees what the AI saw. That is what
"effective" means in Article 14 parlance.

The Article 26(6) record-keeping obligation (logs retained at least 6
months unless longer is required by Union or national law) is satisfied
mechanically by the chain's append-only storage and the institution's
configurable retention. The Article 26(5) monitoring obligation is
satisfied through the operational-events catalog (spec §10.2). Article
26(11) — informing affected persons — composes through the chain's
audit-trail availability when a customer or supervisor asks what the AI
did and why.

The spec's posture against the AI Act is competent. I did not find a gap
my European engagements would flag as a load-bearing concern for v1.0.

## The customer-dispute reproduction posture

The customer-dispute procedures doc (`customer-dispute-procedures.md`)
carries the three load-bearing decisions the MRM committee establishes
before the institution conducts its first reproduction:

- **Chain-or-not posture.** Recommendation: chain the reproduction with
  the original decision's `(run_id, seq)` as the reproduction's
  `parent_run_id` / `parent_seq`. The recommendation is right. A
  reproduction logged outside the chain is repudiable in the same way the
  vendor's logs are repudiable, which is the gap the chain exists to
  close. An institution that accepts the chain as its evidence for the
  original decision and then logs the reproduction outside the chain has
  weakened its evidence base for exactly the scenario where the
  reproduction is most needed.

- **Reproduction-result schema.** The seven-field RECOMMENDED schema
  under `audit.reproduction.*` (`original_run_id`, `original_seq`,
  `iterations`, `decision_equivalent_count`,
  `decision_equivalent_threshold`, `outcome`, `mrm_committee_review`) is
  what the MRM committee actually needs to evaluate variation. The
  `outcome` enum's three values (`decision_consistent`,
  `variation_within_threshold`, `variation_exceeds_threshold`) align with
  how MRM committees express variation findings in practice. The
  `mrm_committee_review` field — populated when
  `outcome = variation_exceeds_threshold` — closes the loop between the
  reproduction evidence and the committee disposition.

- **SOC procedure for testing reproduction evidence.** P-27
  (audit-procedures.md) confirms the institution's reproductions during
  the reporting period have the four-piece evidence package: chained
  reproduction, parent linkage, documented schema, MRM-review link for
  variation-exceeds-threshold cases. Without P-27, the SOC engagement
  cannot test reproduction evidence mechanically; the MRM committee's
  reproduction posture becomes unaccountable across reporting periods.
  P-27 is the right shape for what the SOC team needs.

The variation-threshold guidance (decision-equivalence on 95%+ of re-runs
for high-stakes; 80%+ for low-stakes; tightening for higher
customer-impact tiers) is consistent with how my engagements stratify
their decision portfolios. The pre-established threshold is the
load-bearing detail: an MRM committee that decides the threshold at the
moment of dispute has already lost the calibration. The doc gets this
right.

## The IKM access posture for customer-side verification

This is the section I read most carefully because it is where US and
European banks I work with diverge most sharply on disclosure tolerance.
The spec's two acceptable shapes —

1. **Protective-order disclosure** for litigation context, with court-
   issued protective order limiting use to verification, secure handling,
   non-disclosure binding on the customer's expert.
2. **HSM-mediated verification** for operational context (CFPB inquiry,
   regulatory complaint), with the expert's verifier dispatching HMAC
   operations through the institution's HSM API for the affected runs
   without ever holding the IKM bytes.

— are both conformant with European institutions' typical risk tolerance.
The HSM-mediated path in particular is the right shape for the European
engagements I run, where the institution's IKM-disclosure tolerance is
near zero outside court order. The institution names which shape it
operates and the threshold at which it shifts; without naming the shape,
the customer's expert will encounter `unknown_key_version` at step 7
mid-verification and the dispute is unresolved with no clean diagnosis.
The doc names the failure mode explicitly, which is precisely what an MRM
committee defending the institution's posture in front of a dispute
arbitrator needs.

## P-25 / P-26 / P-27 — the three SR 11-7 audit procedures

Reading these three together is the most informative pass for an MRM
director.

**P-25 (`gen_ai_parameters` schema completeness).** The 3 × 5 × 3 = 45-
entry minimum stratification (5 entries per distinct
`gen_ai.response.model` × 5 entries per `audit.*` event-class × 5 entries
per customer-impact tier) is the right shape. Sampling without
stratification fails to surface the institution's actual decision
pattern; stratification across model-version, decision-class, and impact
tier ensures the SOC team's evidence covers the SR 11-7 reproducibility
surface across the institution's portfolio rather than concentrating on
one model or one decision class.

The framing of schema gaps as MRM-program findings (not chain-integrity
findings) is exactly right. The chain integrity-binds whatever the
institution puts in `gen_ai_parameters`; if the institution's schema is
weaker than the spec's reproducibility surface, that is an effective-
challenge weakness in the institution's MRM program. The chain is doing
its job; the institution's schema is not. P-25 surfaces the gap to the
right audience.

**P-26 (model-inventory composition cross-check).** This is the
procedure that surprised me most positively. The cross-check between the
chain-derived per-model decision count (aggregated by
`gen_ai.response.model`) and the institution's MRM model inventory is a
load-bearing finding generator that I have not seen on prior MRM-control
audit menus. The two failure modes —

- **Models in chain but NOT in inventory.** A model produced decisions
  but does not appear in the inventory. The MRM program is unaware the
  model is in production. High severity.
- **Models in inventory but NOT in chain.** A model is listed but
  produced no chain-captured decisions. Either decommissioned-but-still-
  listed (low severity, housekeeping), or operating outside the chain's
  coverage (high severity, chain-coverage gap relative to MRM program).

— are both real failure modes I have seen on engagements without the
chain. The chain's coverage *enables* the cross-check; without the
chain's per-model decision counting, the institution has no integrity-
bearing source for the cross-check side.

The MRM-COMMITTEE-BRIEF (question 3) routes this to the committee
correctly: the committee asks the chain-ops team for the per-model
decision-count distribution per quarter and cross-checks against the
inventory. This is the right cadence and the right routing.

**P-27 (customer-dispute reproduction evidence completeness).** The
four-piece evidence package — chain integration, parent linkage,
documented schema, MRM-review link for variation-exceeds-threshold cases
— is what makes the institution's reproduction posture testable. The MRM
committee's reproduction policy is now accountable across reporting
periods rather than ad-hoc per-dispute. The procedure also surfaces the
case where reproduction was conducted but the input set was reconstructed
from outside-chain sources (a red flag the SOC team can detect
mechanically).

These three procedures, taken together, give the SOC engagement a
mechanical path through the SR 11-7 reproducibility surface. That is rare
on AI-control audits; most engagement menus stop at the chain-integrity
boundary and leave reproducibility to the institution's MRM program
without an audit anchor. P-25 / P-26 / P-27 is the bridge.

## The MRM-COMMITTEE-BRIEF questions

The seven oversight questions the committee should ask
(MRM-COMMITTEE-BRIEF.md) are well-shaped:

1. AI agent platforms in scope.
2. Seal cadence and any regulator-approved relaxation.
3. Model-state fields captured in `gen_ai_parameters`, plus per-model
   decision-count distribution cross-check against MRM inventory.
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

Question 4's specificity (walk us through the most recent
`master.reconciliation_completed` operational event) is the right level
of detail for committee oversight. Question 3's composition with model
inventory completeness is the bridge between the chain's evidence and
the MRM program's coverage check. Question 5 is loud about the dev-mode
failure mode, which is correct — a `plaintext-dev` URI in production is
the kind of finding that turns an examination from "passed with
observations" into "matter requiring attention."

## The verifier procedure as the validator's leverage point

The twelve-step verifier procedure (spec §7, design 07) is the most
useful artifact in the corpus from an MRM-director perspective.

The procedure's load-bearing properties for SR 11-7:

- **Most-specific-first refusal.** A v2 file fails at step 1 with
  `format_version not supported` rather than later with a HKDF or MAC
  mismatch. The validator reading the verifier output gets a precise,
  named reason rather than a downstream failure that obscures the root
  cause. This is the difference between "chain failed, vendor unclear"
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
  validator's manual inspection.

The dual-algorithm transitional dispatch (cases (a) through (e) in step
11) is the kind of detail I would not have expected v1.0 to ship with.
Case (e) — both signatures present, one valid + one invalid — being
load-bearing-Severe regardless of strict-mode bracket is the right call.
A seal where one of two algorithms validated and the other did not is
either a broken algorithm (the un-broken one's signature still provides
integrity assurance) or a forged signature under a compromised
algorithm-specific signing key (the un-broken algorithm's signature
confirms the un-compromised half). The verifier does not interpret which
case applies; the institution's IR program does. The working paper
records both validations so the regulator has the full picture. This is
the right shape for a post-quantum transition that may take 5+ years
across the financial-services sector.

## The reproducibility chain composes with model inventory

The thesis I would write on a senior-validator memo, after one pass:

> The chain captures every AI agent decision with cryptographic integrity.
> The captured decision carries the model identifier (request and response
> sides), the decoding parameters, the system prompt or its content-hash,
> and the retrieval context — the inputs that determine the AI's output.
> A validator who walks the chain can re-run a decision under the captured
> inputs, evaluate variation, and challenge the model with the validator's
> independent judgement. The chain integrity-binds the inputs the
> validator depends on; the chain's integrity property is what makes the
> validator's challenge load-bearing rather than an exercise in trusting
> the vendor's records.
>
> Composition with the institution's MRM program: the chain produces the
> per-model decision-count distribution that cross-checks against the
> institution's MRM model inventory. Models producing decisions outside
> the inventory (coverage gap) and models in the inventory producing no
> decisions (decommissioned-but-listed or chain-coverage gap) are both
> surfaced. The chain does not replace the institution's MRM program; it
> supplies one of the program's most load-bearing inputs.

That is the framing I would use defending chain adoption to an OCC
examiner under SR 11-7, an EBA examiner under EU AI Act Article 26, or a
MAS examiner under Veritas Accountability and Transparency principles.
The framing is jurisdiction-stable; the chain is the substrate.

## Strengths I would highlight to my partner

A short list, in priority order:

1. **The required model-identifier pair on model-call entries.** Both
   request-side and response-side, integrity-bound under the per-event
   MAC, with the verifier's step 12a procedural check. This is the SR
   11-7 reproducibility property landed at the right granularity.
2. **The customer-dispute reproduction posture, fully specified.**
   Chain-or-not, result schema, SOC test (P-27), and the variation-
   threshold guidance with pre-established tiers. The MRM committee's
   posture is now accountable rather than ad-hoc.
3. **P-26 model-inventory cross-check.** Surfaces both
   model-in-chain-but-not-inventory and model-in-inventory-but-not-chain
   failure modes mechanically. This is rare on MRM-control audit menus
   and is exactly what the OCC wants to see in a documented framework.
4. **Two acceptable IKM-access shapes for customer-side verification.**
   Protective-order for litigation, HSM-mediated for operational
   disputes. The institution names the shape; the failure mode if the
   shape is unnamed (`unknown_key_version` mid-verification) is named
   explicitly so the dispute arbitrator gets a clean diagnosis rather
   than ambiguity.
5. **The dual-algorithm transitional dispatch (case (e) load-bearing-
   Severe regardless of strict bracket).** Post-quantum migration is a
   multi-year transition; the spec is shipping with the right shape for
   what the financial-services sector will live with through the
   migration window.
6. **Per-tenant cryptographic binding via HKDF info.** Cross-tenant
   confusion is structurally prevented at the cryptographic layer, not
   operationally avoided. For SaaS-hosted vendor topologies this is the
   property that makes the per-tenant integrity claim hold under a
   compromised vendor.
7. **Compile-time exclusion of the dev-key adapter from production
   builds.** A misconfigured deployment cannot bring development key
   material online in production. Combined with the verifier's `--strict`
   refusal of `dev_mode = true` seals, this is two independent
   enforcement layers; either alone would be a single point of failure.

## Findings

**BLOCKER.** None.

**MAJOR.** None.

**MINOR.** Three observations, none of which gate adoption on my
engagements. The working group may close any or all at its discretion.

### MINOR-1. Variation-threshold default for medium-stakes decisions is unstated

The customer-dispute procedures doc names the variation-threshold default
for high-stakes decisions (95%+ decision-equivalence on 5+ re-runs) and
low-stakes decisions (80%+). The implicit middle is left to the MRM
committee. On my engagements with three-tier classifications (low /
medium / high) the unstated middle is where the committee spends its
deliberation.

A non-normative recommended default — say, 90%+ decision-equivalence on
5+ re-runs for medium-stakes — would reduce friction for first-time MRM
committees standing up the policy. The committee can override; the
default gives them a starting point.

### MINOR-2. P-26 cross-check stratification could specify a sample-size floor

P-26 requires stratification across decision-class and customer-impact
tier and confirms cross-check completeness across all strata. The
sampling guidance in the audit-procedures doc (Sampling table) gives
population-based sample sizes (< 25, 25-250, 250-2500, > 2500) but does
not specify per-stratum minimums for P-26.

For an institution with 200 models and 4 decision-classes × 3 customer-
impact tiers = 12 strata, the SOC team's professional judgement currently
drives whether each stratum gets at least 1 model represented. A floor of
"at least 1 model per stratum that produced any decision in the period"
would make the SOC team's coverage claim mechanically testable. This is a
SOC-engagement workflow detail rather than a chain-integrity concern.

### MINOR-3. Model-version concurrent-deployment recordkeeping is implicit

The chain captures `gen_ai.response.model` per decision. When an
institution operates two model versions concurrently (A/B testing,
canary rollout, blue-green deployment of a new model version), the
per-model decision count distinguishes the versions automatically. What
is not explicit in the corpus is whether the institution should record
the *intent* of the concurrent deployment as a chain entry — "we are
operating model X-v1.2 and model X-v1.3 concurrently for the rollout
period from D1 to D2."

For SR 11-7 model-validation work, the deployment intent is what the
validator validates against. The chain records the *observed* per-model
decision counts; the validator's working paper benefits from a chain
entry recording the *intended* deployment topology. This is an
institutional practice rather than a chain-spec concern, but the MRM-
COMMITTEE-BRIEF could surface it as an oversight question. Adding an
eighth committee question — "What is the institution's posture on
recording concurrent-deployment intent in the chain?" — would close the
loop.

## Engagement-side adoption posture

For the three engagements named at the top of this review:

- **Pan-EU retail bank.** Adopt. The Article 14 effective-oversight
  framing, Article 26 record-keeping mechanics, and Article 12 logging-
  integrity property all land. The customer-dispute reproduction posture
  composes with the bank's existing CFPB-equivalent dispute procedures
  for the bank's UK and EU consumer-credit operations. The HSM-mediated
  IKM-access shape is the right default for the bank's risk tolerance.

- **US bank holding company.** Adopt. The SR 11-7 effective-challenge
  framing is the headline; P-26 model-inventory cross-check is the
  finding-generator that gives the bank's MRM program a mechanical
  coverage check it does not currently have. The per-model decision-
  count distribution is what the OCC will ask for at the next examination
  cycle and the bank should be ready to produce.

- **Singapore-headquartered bank with cross-jurisdiction operations.**
  Adopt. The chain is jurisdiction-neutral; the bank's per-jurisdiction
  control descriptions reference the chain consistently. Veritas
  Accountability and Transparency principles compose with chain output
  directly. DORA Article 6 / 8 / 17 / 28 composition is documented;
  bank's ICT risk function consumes the chain's threat model and IR
  playbook as the chain-specific contributions to the bank's DORA
  framework.

## Final disposition

**0 BLOCKER, 0 MAJOR.** The chain spec, the design corpus, the MRM-
COMMITTEE-BRIEF, the customer-dispute procedures, the AI policy
alignment, and the audit procedures (P-25 / P-26 / P-27) compose into a
coherent SR 11-7 effective-challenge support layer that I would
recommend on all three engagements named above without reservation.

The minor observations are working-group-discretion items; none of them
gate the partner's sign-off on adoption.

— Dr. Elena Vasilescu, MRM Director
   Bucharest office
