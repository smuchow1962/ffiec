# Round 14 — FFIEC Cybersecurity Specialist Examiner review

**Reviewer:** Akira Nakamura, FFIEC Cybersecurity Specialist Examiner — Federal Reserve Bank, Second District (NY).
13 years on the cyber-specialist track. Background: information assurance,
incident response, supply-chain assessment under GV.SC, post-quantum readiness
working groups, joint-agency coordination on AI risk supervision.

**Stopping criterion:** 0 hard finding, 0 soft finding.

**Scope read for this pass.**

- `spec/chain-of-custody-v1.md` §10 (operational requirements, normative) and
  §7 verification procedure with focus on step 11 dual-algorithm dispatch and
  step 12a GenAI completeness check.
- `docs/design/00-overview.md` — system architecture, trust boundaries,
  deployment topologies, four named attack approaches.
- `docs/design/09-threat-model.md` — adversaries A through I, residual risk
  register R1–R13, Adversary I reception procedure.
- `docs/regulator-pack/CSF-2.0.md` — NIST CSF 2.0 mapping, headline alignment,
  attack-to-subcategory binding, operational-event-to-subcategory binding,
  GOVERN-function alignment.
- `docs/incident-response-playbook.md` — Scenarios 1 through 12, the 36-hour
  triage matrix, the CIRCIA matrix, federal-regulator routing per charter,
  dual-algorithm verifier-version timing-overlap, Scenario 11 sub-variant
  for forged rotation notice, Scenario 12 case-(e) branches.
- `docs/cloud-hsm-guide.md` — per-provider conformance matrix, provisioning,
  cross-region, cost picture.
- `docs/supply-chain.md` — release artifacts, signing roles, key-recovery
  procedures, mirror-registry patterns, SLSA L3 institutional consumption,
  spec-and-corpus integrity, cold-DR-key lifecycle.
- `docs/regulator-pack/ai-policy-alignment.md` — NIST AI RMF, EU AI Act,
  DORA, NIS2, EBA, MAS, JFSA, BCB, OCC supervisory posture.
- `docs/edge-and-federated-ai.md` — edge deployment patterns, federated
  learning scope boundary, on-device inference, multi-agent autonomous
  coordination.

I deliberately did not read prior-round feedback in this directory. The
review reflects a first-look reading of the shipped artifacts.

---

## 1. The integrity claim, evaluated against my examiner workflow

I examine institutions against the FFIEC IT Examination Handbook plus the
NIST CSF 2.0 frame. My working pattern on a cyber-specialist engagement is:
identify the institution's claimed CSF subcategory coverage; sample-test
the operational evidence per subcategory against the institution's log
store; trace any failed sample to the root cause and decide whether the
finding is a control gap (institution-side remediation) or a residual the
threat model accepts (institution documents and moves on).

The chain's headline integrity claim is **PR.DS-06: Integrity-checking
mechanisms verify data integrity.** The CSF mapping document binds the
chain's primitives to PR.DS-06 cleanly and resists the temptation to
overclaim against subcategories the chain composes with rather than
satisfies (PR.DS-01, PR.AA at principal level, DE.CM-01..08). That
restraint is the right discipline. An examiner reading the document
knows precisely which subcategories the chain alone supports and which
require the institution's broader posture to come along.

The four-attack approach in `00-overview.md` §2 maps cleanly onto the four
load-bearing CSF subcategories in `CSF-2.0.md`'s attack table. The
operational-event-to-subcategory binding gives me a working evidence
matrix on day one of an examination: I sample
`master.reconciliation_completed` events for PR.AA-05 + DE.CM-09;
`chain.verification_failure step=8` events for PR.DS-06 + ID.AM-08 +
RS.MA-04 (Scenario 7); `mirror.reconciliation_completed` events for
GV.SC-04. Each event has a primary subcategory, a secondary subcategory,
and a named IR scenario. The matrix is not invented at examination time;
it is shipped with the documentation.

The GOVERN function alignment section is the right depth. CSF 2.0 made
GOVERN a top-level function and examiners now expect institutions to
articulate how integrity controls flow from the risk-tolerance statement
down to operational evidence. The document maps GV.OC, GV.RM, GV.RR,
GV.PO, GV.OV, and GV.SC to chain-specific artifacts an institution can
copy into its own GOVERN policy. That keeps the institution from
reinventing the GOVERN narrative under examination pressure.

## 2. The threat model, evaluated against the cybersecurity-specialist lens

`09-threat-model.md` enumerates Adversaries A through I with capability,
goal, defense, and residual risk for each. The residual risk register
(R1–R13) is the auditor's-lens artifact I need: each risk has a
likelihood, an impact, a mitigation status, and an owner.

The construction is honest about what it does not claim. The §4 "Properties
NOT claimed" section names AI decision correctness, real-time prevention
of bad behavior, side-channel attacks on the application host, and
operational monoculture as out of scope. That keeps the integrity claim
defensible against a cybersecurity-specialist sampling the institution's
controls; I know where the chain's story ends and where the institution's
broader posture begins.

The Adversary I institution-side reception procedure for a regulator-held
fingerprint rotation is the most operationally demanding piece of the
threat model. The procedure spells out three institution-defined
operational events (`regulator_fingerprint.rotation_received`,
`regulator_fingerprint.rotation_validated`,
`regulator_fingerprint.installed`), the two-channel validation discipline,
the reception-failure sub-variant (forged-notice suspected), and the
re-validation-of-historical-reports requirement. This is the kind of
trust-anchor lifecycle most institutions get wrong because they treat the
trust anchor as static. The chain's document treats it as a managed
artifact with its own change-management discipline. That maps to GV.SC-04
(third-party assessments) in the operational-event binding.

The §3.1.2 identity attribution boundary is the right scope. The chain
captures tenant-level identity via `key_fingerprint`; principal-level
identity is the institution's IAM layer. The composition narrative
(IAM session ID resolves to a principal via the bank's IAM logs;
chain event references the IAM session ID; verifier confirms chain
integrity; IAM logs answer "which principal authored this") is what an
SR 11-7 examiner cross-checking accountability needs.

R12 (edge-device physical compromise) being deferred to v1.1 with explicit
documentation of the residual is the honest disposition. The
`edge-and-federated-ai.md` document operates Pattern A (per-device IKM in
TPM/secure-enclave) as the v1.0 compensating control with reconciliation
cadence baseline tuned per fleet. An institution running edge AI knows
the residual; the threat model does not pretend the v1.0 spec closes it.

## 3. Spec §10 operational requirements — examination posture

Section 10 is normative. I read it the way I read the FFIEC IT Handbook:
each subsection is a control I sample-test.

**§10.1 Key-fingerprint reconciliation.** Weekly cadence as a SHOULD,
matched against the institution's IKM roster, evidenced as
`master.reconciliation_completed` with `unmatched_count`. The "no more
than weekly" framing bounds the master-compromise detection window; the
institution documents the cadence and the response procedure for any
non-zero unmatched count. P-6 in the audit-procedures sample is the
SOC team's cross-check. From my sampling angle: pull the most recent
`master.reconciliation_completed` event, compute the cadence delta to the
prior reconciliation, confirm the delta does not exceed the institution's
declared cadence; pull any `unmatched_count > 0` event, confirm a
corresponding Scenario 7 IR record exists.

**§10.2 Operational events.** The catalog is comprehensive and the
retention rule ("retained at least as long as the chain events they
relate to") aligns with the FFIEC retention expectations. The
master-key lifecycle events (`master_key.rotated`,
`master_key.rotation_observed`, `master_key.retired`) and the
regulator-fingerprint lifecycle events (`regulator_fingerprint.*`) are
the load-bearing trust-anchor evidence I sample-test.

**§10.3 Append-only enforcement at two layers.** Application level (no
UPDATE/DELETE on events or daily_seals tables) plus database-role level
(INSERT and SELECT only; UPDATE/DELETE/TRUNCATE revoked). The defense in
depth is the right shape; the Merkle seal catches deletion regardless of
which layer it occurs at, and the role-level enforcement is a
defense-in-depth control I sample-test the institution's database-role
configuration against.

**§10.5 HSM custody.** FIPS 140-2 Level 3 minimum, non-extractable private
key, sign-only authority for the seal-job operator role, separation of
duties from the HSM administrator. Small-institution dual-control as a
documented compensating control is the right proportionality. The
`cloud-hsm-guide.md` matrix tells me which provider-and-tier combinations
are conformant and which are not; the disqualification of AWS KMS default
tier and Azure Key Vault Standard is correct and useful.

**§10.6 IKM minimum length.** 32 bytes per RFC 4868. The rationale
(HMAC-SHA-256 keying recommendation; the public 16-byte fingerprint must
not be offline-grindable) closes R8 (session key leakage to logs) for
the fingerprint side-channel. Enforcement at IKM-provisioning time
(refuse to register an under-length IKM) plus enforcement at SDK-configure
time (refuse to start the chain writer with a short IKM) is the right
two-point control.

**§10.7 Software-key adapter compile-time exclusion.** The compile-time
gate plus the runtime stamp (`kms_handle_uri = "plaintext-dev"`) plus the
verifier's --strict refusal of `dev_mode = true` is triple-defense
against software-key fallback in production. The "run-time
environment-variable gating is NOT sufficient" sentence is the load-bearing
discipline. A misconfigured deployment that flips the flag must NOT be
able to bring the software adapter online in production. Scenario 6 in
the IR playbook handles the rare case where the exclusion is bypassed.

**§10.8 Constant-time comparison.** Required for fingerprint check
(step 8) and MAC check (step 9). Stdlib helpers named per language
(`hmac.compare_digest` Python; `subtle.ConstantTimeCompare` Go;
`CryptographicOperations.FixedTimeEquals` .NET). The discipline applies
to both checks even though the fingerprint is publicly stamped on every
entry — the rationale ("a future maintainer extending the verifier does
not reach for `==` on the MAC compare") is the right defensive posture.

**§10.9 IKM registry retention.** "Every IKM generation MUST be retained
for at least as long as any chain entry stamped with that key_version is
retained." The retirement coupling at the registry layer (explicit
override required to retire an IKM whose key_version is still referenced)
plus the `master_key.retired` operational event closes R11 (IKM-registry
premature retirement). The KMS provider retention windows (7-30 day
pending-window) are documented; the conservative posture (longer of the
chain entry retention and the institution's regulatory minimum, typically
7 years) is the institution's documented baseline.

**§10.10 Rotation crossing the seal boundary.** The day-after seal's
`key_versions` field listing both old and new generations during the
transient window is the right operational discipline; the verifier handles
this case by per-entry `key_version` lookup with no special-case logic.
The `master_key.rotation_observed` operational event provides the
audit-evidence for the rotation event itself.

The §10 catalog is the right normative shape: each subsection names a
control, the rationale, the operational evidence, and the verifier
behavior. An examination-grade institution can adopt the §10 controls
verbatim and sample-test against them.

## 4. Spec §7 verification — the dispatch logic that matters at examination

Step 11 is where the dual-algorithm transitional period dispatch lives.
The case enumeration (a) through (e) is the load-bearing logic an
institution operating dual-algorithm posture commits to. I read it
carefully:

- **Case (a):** both signatures present and both valid → PASS. Working
  paper records both validations.
- **Case (b):** single algorithm signature during dual-algorithm posture
  → PASS-WITH-ANOMALY (control-completeness, not chain-integrity). The
  institution investigates as a control-completeness finding.
- **Case (c):** seal carrying an algorithm not on the institution's
  declared posture list → --strict FAIL; non-strict PASS-WITH-ANOMALY.
  The institution's declared posture is consulted via institution
  configuration (out of scope for the spec).
- **Case (d):** single-algorithm posture (default for v1.0) → reduces to
  single-algorithm verification.
- **Case (e):** both signatures present, one valid, one invalid →
  --strict FAIL; non-strict PASS-WITH-ANOMALY. **Severity is Severe
  regardless of bracket.** The PASS-WITH-ANOMALY disposition under
  non-strict reflects the un-broken algorithm's signature still
  providing integrity assurance, NOT that the failure is itself
  low-severity.

The case-(e) framing is precisely correct from my examination angle. The
verifier output dispenses bracket-level dispositions, but the institution's
finding-language repository (`finding-language.md` row "11 (dual-algo)
co-signed seal failure") cites Severe MRA regardless of bracket. That
keeps the examiner from misreading PASS-WITH-ANOMALY as low-severity. The
working-paper convention (record BOTH the valid-algorithm validation and
the invalid-algorithm failure) gives the regulator the full picture even
when the verifier's headline disposition is bracket-dependent.

The Scenario 12 IR triage in the playbook fits the spec dispatch:
branch (i) attributable to a published algorithm break (no FFIEC clock,
regulator-coordinated), branch (ii) per-algorithm signing-key compromise
(STARTS the clock at credibility determination per Scenario 4),
branch (iii) under investigation (STARTS the clock at investigation
conclusion within a bounded 48–72-hour window). The branch enumeration
is honest about the determination point and the bounded window. The
under-investigation branch's bounded window prevents the clock-start
from being indefinitely deferred.

**Step 12a GenAI model identifier completeness check** is the
SR 11-7 reproducibility evidence. The check fires only on chain entries
representing model calls (entries with any `gen_ai.*` attribute);
entries without `gen_ai.*` (tool calls, audit-only events) are
unaffected. The disposition is --strict FAIL or non-strict
PASS-WITH-ANOMALY framed as control-completeness for SR 11-7
reproducibility, NOT chain-integrity. That separation of concerns is
correct: a missing model identifier is a model-risk control gap (the
MRM team's work), not an integrity gap (the chain construction's work).

## 5. Incident response playbook — examination posture

The playbook covers Scenarios 1 through 12 with severity classification,
investigation triage trees, containment, remediation, and notification.
I focus on the cybersecurity-specialist sections:

**The 36-hour clock-start trigger matrix per scenario.** The matrix
distinguishes between "alert receipt" and "determination." The
determination point starts the clock. Scenario 1 (chain hash mismatch)
clock starts when investigation rules out SDK defect AND confirms
in-flight modification. Scenario 7 (key_fingerprint mismatch) clock
starts only when triage rules out documented operations (rotation in
flight, restored backup, restored tenant row) AND treats as suspected
unauthorized substitution. The triage tree pre-sorts the operational
cases (no clock) from the security-incident case (clock starts), which
matches the FFIEC computer-security incident notification rule's
"determining" framing precisely.

**Scenario 11 sub-variant — forged rotation notice.** When the
institution's regulator-fingerprint reception validation fails (notice's
GPG signature does not validate, OR cross-channel parallel notification
disagrees), the institution does NOT install the new fingerprint; the
verifier configuration retains the previous fingerprint as the active
trust anchor. Three plausible root causes are enumerated: forged notice
(severe; full IR activation per Adversary I framing; STARTS the
36-hour clock at the determination of forgery); operational error at
the regulator (re-issue corrected notice; no clock); institution-side
validation tooling defect (update tooling, re-run reception; no clock).
The CIRCIA matrix flags the forgery case as both FFIEC-trigger (yes) and
CIRCIA-trigger (yes — substantial cyber incident, attacker attempted
trust-anchor substitution).

**Scenario 12 case-(e) interaction with verifier-version timing-overlap.**
The most operationally subtle scenario in the playbook. The institution
is operating a Y-algorithm-aware verifier (post-quantum); the regulator's
verifier supports only X-algorithm. A case-(e) failure (institution
sees Y validates, X does not) may be invisible to the regulator's
X-only verifier (the regulator sees the X-algorithm signature as the
only signature and validates it; the Y-algorithm half is structurally
invisible). The institution MUST proactively notify the regulator of the
case-(e) failure and supply the institution's case-(e) verifier output as
supplementary evidence. This is a known regulator-side blind spot during
the dual-algorithm transitional period; the regulator's IT examination
program SHOULD upgrade to a Y-algorithm-aware verifier within 6 months
of the institution's adoption.

The pre-coordination guidance ("the institution's IR program SHOULD
pre-coordinate verifier-version compatibility with the regulator BEFORE
the multi-year transitional period begins") matches my district's
working-group expectation. We don't want institutions surprising us at
examination time with a co-signed seal our verifier cannot evaluate; the
institution surfaces the verifier-version state in advance and we
provision compatibility on our side.

**CIRCIA-only triggers without FFIEC trigger.** Some chain-detected
events rise to CIRCIA's "substantial cyber incident" threshold (72-hour
CISA path) without rising to the FFIEC computer-security rule's
"safety/soundness" threshold. The matrix names which scenarios are
CIRCIA-likely without FFIEC-likely (Scenario 5 sealing delay associated
with broader incident; Scenario 9 truncation pattern across hosts;
Scenario 11 dual-compromise impacting institution operations at scale).
The "counsel works the matrix per incident" framing recognises CIRCIA's
broader scope. The institution's IR program SHOULD have CIRCIA
notification on its standard checklist for chain-detected events, even
when the FFIEC clock has not started.

**Federal-regulator routing per institution charter with CFR citations.**
The matrix maps charter to primary federal regulator with the relevant
12 CFR citation: National bank / federal savings → OCC under Part 53;
state member bank / BHC / state savings & loan HC → FRB under Part 225
App. F; state non-member bank / state savings → FDIC under Part 304
Subpart C plus state regulator; federal credit union → NCUA under Part
748 App. B; state credit union → NCUA plus state regulator. The
multi-charter holding-company posture (notify all applicable primary
regulators at the affected entity's chartering structure) is the right
disposition. State-side cadence column tracks the state-specific timing
where applicable. The IR Commander confirms routing per institution
charter at incident time; the standing IR documentation pre-records
the primary regulator and secondary state-side counterpart. That's the
right operational posture.

**Concurrent multi-scenario alerts (rollup rule).** Scenario 1 + Scenario
7 + Scenario 8 alerts in the same 60-minute window may be one root-cause
event triggering three detection paths. The rollup rule (treat as one
incident; start the 36-hour clock at the earliest qualifying determination
across the constellation, NOT at the union of per-scenario clocks) is
correct. Treating each scenario's clock independently risks both
double-counting and under-counting depending on whether the constellation
is the signal.

**Cross-tenant scope discovery during investigation (vendor-hosted
topology).** The 36-hour clock applies per institution; in vendor-hosted
deployment one investigation may produce N institution-side determinations
on different timelines. The vendor's contractual notification SLA
(typically 4-12 hours from vendor's determination) is a SEPARATE clock
from the institution's 36-hour FFIEC clock — both apply, both are tracked.
The vendor-vs-institution responsibility split is the right disposition:
vendor's IR team produces the alert and conducts cross-tenant
investigation; each institution's IR program receives the alert through
the vendor's notification channel and produces its own determination.

## 6. Supply-chain controls — examination posture

The supply-chain document covers what most spec-of-this-class documents
gloss over. The release artifacts table names eleven artifacts including
spec PDF + spec content-hash manifest + conformance corpus tarball + SLSA
provenance attestation. The signing-role separation (cosign held by
release-management role; GPG held by a separate small group; cold-DR
key held by a third separated small group) is the right
defense-in-depth.

**Cosign-key recovery procedure.** ≤ 24 hours project response (revoke
certificate, publish revocation notice signed by GPG); ≤ 48 hours
institution notification through normal channels plus FFIEC working group;
≤ 7 days new key authentication chain (provision new cosign key, publish
fingerprint in spec text, sign KEY-INTRODUCTION-{new_keyid}.asc under
GPG). Institutions remove revoked key, install new key validated against
GPG-signed introduction notice, re-validate binaries signed under the new
key. Binaries signed under the revoked key remain verifiable for historical
audit purposes; the revoked-but-archived key stays on the institution's
trust-anchor list for historical seal verification.

**GPG-key recovery procedure.** Symmetric to cosign-key recovery, with
the dual-compromise edge case explicitly named. If both keys are
confirmed compromised in the same window, the project escalates to the
cold-disaster-recovery key. The institution's emergency response: stop
deploying new releases of verifier and ledger server until the project
completes the new-key authentication chain; continue using the historical
verifier binary against historical seals; coordinate directly with the
primary regulator on examination posture during the multi-week
dual-compromise recovery window.

**Cold-DR key lifecycle.** 60-month rotation cadence aligned with NIST
SP 800-57 guidance for cold-storage keys. Annual dry-run cadence. The
dry-run output published as `KEY-DR-DRYRUN-{year}.asc` signed under the
standard cosign or GPG key. Key-holder rotation per departure /
governance event, not per key rotation; emergency rotation within 30
days of any key-holder departure. Health check on annual cadence
(without exposing key bytes — HSM tamper detection state, hardware-token
battery indicator).

**Institution-side consumption of cold-DR-key dry-run attestation.** The
institution's release-validation procedure includes annual consumption
of `KEY-DR-DRYRUN-{year}.asc`. The institution validates the
attestation's signature, archives in control-evidence repository for the
year, treats failure to obtain within 30 days of project's annual
dry-run window as Scenario 11 branch (c) (missed dry-run window). The
institution's `KEY-DR-DRYRUN-{year}.asc consumption log` is the
operational evidence per CSF-2.0 GV.SC-04 binding; SOC and examination
teams sample-test the consumption log against the institution's archived
attestations. Without this consumption procedure, the institution treats
the cold-DR fallback as documented-but-unexercised, which a GV.SC-04
examiner would challenge as un-tested control evidence. That's exactly
the challenge I would raise; the document anticipates and answers it.

**SLSA Level 3 institutional consumption.** The institution MUST validate
the published `verifier.intoto.jsonl` against `slsa-verifier` BEFORE
deploying a new release version. The validation is a deployment-gating
control; a binary whose provenance attestation does not validate MUST NOT
be deployed. The SLSA-verifier binary itself MUST be signed and
trust-anchored alongside the verifier. The trust-anchor table includes
`slsa-verifier` public key + binary SHA-256 cached out-of-band. The
SLSA-aware tool-choice documentation (which tool, what version, what
trust anchor, archived where) is normative for examination-grade adopters.

**Mirror-registry signature handling.** Two patterns documented:
signature-preserving mirror (preferred; bytes-faithful; project signature
preserved) and re-signing mirror (institution validates project
signature at pull, re-signs under institution's internal cosign key,
documents the trust-path bridge in the institution's IR playbook). The
re-signing pattern requires:

- Bridge invariant documented ("the mirror's re-signature is conditional
  on the project's cosign signature having validated at mirror-pull
  time, recorded in the mirror's audit log").
- Audit log on WORM (write-once-read-many) storage, NOT rotation-deletable.
- Audit log retention 7 years minimum or longer if institution's
  chain-event retention exceeds 7 years.
- Mirror continuous-monitoring control (`mirror.reconciliation_completed`)
  to detect mirror silently stopping signature validation.

The mirror reconciliation control is the supply-chain-side analog of
spec §10.1 fingerprint reconciliation. An institution operating the
re-signing pattern without the WORM audit log + 7-year retention + the
reconciliation discipline has a GV.SC-04 gap I would write up.

## 7. Cloud HSM conformance — examination posture

The cloud HSM matrix is the answer to the question I most often have at
the start of an examination: "is the institution's HSM tier conformant?"
The matrix disqualifies AWS KMS default tier (FIPS 140-2 L2 only) and
Azure Key Vault Standard (FIPS 140-2 L1) and Google Cloud KMS software
keys (software-backed; not L3). Azure Key Vault Premium legacy SKU is
"Partial — FIPS 140-2 L2 default; L3 available; consult Azure docs." The
disqualification list catches the most common institution-side
misconfiguration: an institution that thinks "we're using AWS KMS" when
they should be using AWS CloudHSM Custom Key Stores backed by CloudHSM.

The cost picture per provider is documented (CloudHSM ≈ $13k/year per
member; Azure Managed HSM ≈ $15k/year per partition; Google Cloud HSM
$5k–$15k/year depending on volume; on-prem $30k–$80k/year all-in). The
proportionality picture (community bank shared cloud HSM; mid-size bank
dedicated cloud HSM single region; tier-1 multi-region HSM cluster) lets
the institution defend its tier choice to me on cost grounds.

The pitfalls section (mixing HSM tiers; sharing keys across tenants;
backup keys without replication; PIN management drift) names what
institutions get wrong. PIN management drift is the most common gap I
catch: institutions provision the HSM and never rotate the PIN; a leaked
PIN that has not been rotated is a vulnerability.

## 8. AI policy alignment — examination posture

The AI policy alignment document is the document my agency sends back to
an institution when they ask "how does the chain align with AI RMF?" or
"how does the chain support EU AI Act Article 12?" The document maps
the chain's primitives to:

- **NIST AI RMF** Govern + Map + Measure + Manage functions with specific
  function-to-primitive bindings.
- **Treasury AI RMF (Feb 2026)** model-inventory, model-lifecycle
  controls, third-party model governance, ongoing monitoring — four
  composition points the MRM director cites in front of an OCC examiner.
- **OCC supervisory posture (no formal new AI rule, operationalised under
  SR 11-7 / OCC Bulletin 2011-12).** The framing avoids over-claiming
  the OCC has explicitly required the chain (it has not) while
  recognising the chain's relevance to the OCC's evolving supervisory
  posture under existing model-risk-management framework.
- **EU AI Act Article 12 (record-keeping)** as the headline alignment;
  Articles 14 (human oversight) and 26 (deployer obligations) as the
  load-bearing secondary alignments where the chain is a substrate for
  the deployer's compliance posture, not just a record-keeping artifact.
- **DORA Articles 5, 6, 8, 17, 28** with the chain composing within the
  institution's broader ICT framework, not as a standalone control.
- **NIS2 Articles 21, 23.**
- **EBA Guidelines on AI in Financial Services Section 5.**
- **UK FCA / PRA, Singapore MAS Veritas / FEAT, JFSA, BCB.**

The document is honest about cross-jurisdictional adaptation being the
institution's compliance program's responsibility, with the chain
providing the consistent substrate. The chain's primitive set is
policy-stable: integrity, tamper evidence, independent verifiability are
durable concepts across policy evolution. The chain may need to add new
attributes or extensions to align with new policy specifics; the spec
evolves with the policy landscape.

## 9. Edge and federated AI — examination posture

`edge-and-federated-ai.md` accommodates edge / federated / on-device /
multi-agent autonomous-coordination patterns without spec modification.
The patterns:

- **Edge AI:** SDK runs in agent process on edge device; HMAC chain
  computed at capture; local SQLite persistence; OTLP export when
  connectivity available. Pattern A (per-device IKM in TPM/secure-enclave)
  and Pattern B (bulk session-key issuance at commissioning) are the two
  master-key custody options.
- **Federated learning:** chain captures decisions (inference phase),
  not training. Training is MLOps controls outside chain scope; the
  institution's MLOps program operates training; chain composes with the
  resulting model deployment.
- **Multi-agent autonomous coordination:** DAG semantics handle the
  topology. Each agent is a chain participant with own session key; agent-
  to-agent communication captured as cross-run references (`parent_run_id`
  linkage). Verifier walks each run independently; cross-run topology is
  the institution's analytical responsibility.
- **On-device inference (mobile, IoT):** same pattern as edge AI;
  privacy-by-design tokenization applies; privacy-store handshake
  decision based on device-trust posture.
- **Cross-institution federated patterns:** consortium-style deployments
  may operate shared chain (consortium-level master key custody) or
  separate per-member chains (per-member master keys). Choice affects
  tenant_id structure but not substance.

R12 (edge-device physical compromise; secure-enclave attestation
defeated) is the v1.0 residual the institution documents and operates
compensating controls against. The threat model's R12 row names the
mitigation status: deferred to v1.1; v1.0 deployments document the
residual and operate Pattern A with reconciliation cadence baseline tuned
per fleet. From my examination angle: an institution running edge AI
documents R12 in the institution's risk-management framework as an
accepted residual with the compensating control (Pattern A + reconciliation
cadence) named and operationally exercised.

## 10. Stopping criterion

I read the documents listed at the top of this review as a first-look
cybersecurity-specialist examiner. I evaluated the integrity claim
against the FFIEC IT Examination Handbook plus the NIST CSF 2.0 frame.
I sample-tested the operational-evidence catalog against the
subcategory-binding table. I traced the four named attack approaches
through the threat-model adversaries, the residual-risk register, the
incident-response scenarios, the 36-hour clock-start triage matrix, the
CIRCIA matrix, the federal-regulator routing per institution charter,
the dual-algorithm verifier-version timing-overlap, and the Scenario 11
sub-variant for forged rotation notice. I read the supply-chain controls
including the cold-DR-key lifecycle and the institution-side consumption
procedure for the dry-run attestation. I read the cloud HSM conformance
matrix and the AI policy alignment for the policy frameworks I expect to
see at examination time.

The spec §10 normative requirements give me the controls I sample-test.
The §7 verification procedure gives me the dispatch logic I trace at
examination time, including step 11's dual-algorithm case enumeration
with case (e)'s Severe-regardless-of-bracket disposition and step 12a's
GenAI completeness check framed as control-completeness for SR 11-7
reproducibility (not chain-integrity). The CSF 2.0 mapping gives me the
attack-to-subcategory binding plus the operational-event-to-subcategory
binding I work into an evidence matrix on day one of the examination.
The threat model's residual risk register R1–R13 documents the accepted
residuals. The IR playbook covers the chain-detected scenarios with the
36-hour clock-start determination point precisely framed and the CIRCIA
broader scope explicitly named. The federal-regulator routing matrix
tells me which CFR citation I am examining the institution against. The
supply-chain controls cover the cold-DR fallback with institution-side
consumption discipline an examiner can sample-test. The cloud HSM
matrix disqualifies the common misconfigurations. The AI policy
alignment maps the chain to the frameworks my agency receives questions
about.

I have no hard findings.

I have no soft findings.

The documentation set is at the depth and rigor I expect from an
examination-grade integrity-control specification. I would not write up
a finding against an institution presenting the chain in this form. The
institution's adoption of the v1.0 spec, with the operational-event
catalog wired to its log store, the IR playbook adapted to its IR
framework, the cloud HSM provisioned to the matrix, the supply-chain
controls operationally exercised including the SLSA L3 consumption and
the cold-DR-key dry-run attestation, plus the AI policy alignment cited
in the institution's control description, gives my examination team the
evidence base we need to confirm PR.DS-06 coverage at scale.

**Stopping criterion met: 0 hard finding, 0 soft finding.**

---

*Akira Nakamura*
*FFIEC Cybersecurity Specialist Examiner*
*Federal Reserve Bank of New York, Second District*
*Round 14, first-look review — 2026-05-06*
