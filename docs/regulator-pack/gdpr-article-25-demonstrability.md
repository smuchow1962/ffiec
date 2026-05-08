# GDPR Article 25 — privacy-by-design demonstrability

> **What this doc is.** The five artifact types the institution produces to demonstrate privacy by design under GDPR Article 25. Strong technical design alone is not enough; Article 25 also requires the institution to be able to *prove* that privacy by design was implemented. The proof is a pack of artifacts an examiner, an auditor, or a Data Protection Authority can review without having to take the institution's word for any of it. This doc names the five artifact types, what each carries, and how they assemble into a demonstrability pack.

## Why demonstrability is a separate obligation

Article 25(1) requires the controller to "implement appropriate technical and organisational measures" for privacy by design. Article 25(2) requires the controller to ensure, by default, that only personal data necessary for the specific purpose is processed. Article 5(2) — the accountability principle — requires the controller to be "responsible for, and be able to demonstrate compliance with" the GDPR principles.

Implementation alone is not the standard. The institution must be able to *show* the implementation. A privacy-by-design design that cannot be evidenced is, for accountability purposes, indistinguishable from no design at all. The institution that cannot produce evidence loses the case before the merits are reached.

The five artifacts below collectively answer the question: "Show me what you did and how you know it worked." They are the institution's demonstrability pack.

## The five artifacts

```
1. Configuration   — version-controlled SDK + privacy configuration
2. Build           — build logs proving the configured artifact was deployed
3. Tests           — regression test results proving redaction/tokenization works
4. Operations      — operation logs proving privacy-store access is controlled
5. DPIA            — risk assessment proving the institution thought it through
```

The pack is reviewed in a single sitting. An examiner walks through it artifact-by-artifact and reaches a determination on whether privacy by design was implemented. Each artifact stands alone; together they answer the Article 25 + 5(2) accountability requirement.

## Artifact 1 — Configuration

The institution's SDK configuration determines what the chain captures, what it tokenizes, what it redacts, and what it omits entirely. The configuration is the central design choice that makes privacy by design real.

**What it carries:**

| Element | Content |
|---|---|
| Field-level capture rules | Which audit-namespace attributes are captured, which are tokenized, which are redacted, which are omitted. The privacy-by-design.md SDK-configuration section is the schema. |
| Tokenization scope | The list of fields that go through the privacy-store mapping (customer ID, account number, etc.). |
| Regex-redaction patterns | The patterns that scrub free-text fields (model prompts, reasoning chains) for inadvertent PII. |
| Special-category handling | Stricter rules for Article 9 fields (health data, biometric, behavioral inferences). |
| Retention overrides | Per-field retention rules where they differ from the chain default (rare, but documented when present). |

**Demonstrability requirement.** The configuration is stored in version control. Each version is tagged with its effective date. The institution can produce the configuration in effect on any specific date in the past seven years.

**Change management.** Configuration changes go through the institution's change-management process. The change record names: who proposed the change, who approved it (typically the DPO and the chain-operations lead), the date, the rationale, and the rollout date. Material changes — new tokenization scope, new redaction patterns, expanded capture — trigger DPO consultation per `gdpr-dpo-consultation.md`.

**Examiner question this answers.** "What did your privacy configuration look like in Q3 2026 when this customer's decisions were captured?" Answer: pull the version-controlled configuration tagged for that period.

## Artifact 2 — Build

The configuration in version control proves the design intent. The build artifacts prove the design was actually deployed.

**What it carries:**

| Element | Content |
|---|---|
| Build manifest | The SDK package version, the configuration version, and the configuration hash that the build was sealed against. |
| Deployment record | When the build was promoted to production, by whom, against which environment. |
| SLSA / supply-chain attestation | Per `supply-chain.md`: signed build provenance proving the deployed binary matches the source tree and configuration. |
| Software-key adapter status | Per spec §10.7: confirmation that the production build does NOT include the software-key adapter. The build manifest explicitly records the absence. |

**Demonstrability requirement.** For any chain entry in production, the institution can identify which build produced it (via the spec version stamp on the entry) and prove that build was the version intended (via the build manifest and provenance attestation).

**Examiner question this answers.** "Prove the SDK that captured this entry was running the configuration you say it was running." Answer: pull the chain-entry stamp, the build manifest, and the provenance attestation. The configuration hash on the build manifest must match the version-controlled configuration's hash for that period.

## Artifact 3 — Tests

The build artifacts prove what was deployed; the test results prove what it actually did. Privacy by design is asserted in design; it is verified in tests.

**What it carries:**

| Element | Content |
|---|---|
| Tokenization regression tests | Tests that confirm tokenized fields produce tokens (not original values) when captured. Run on every build. |
| Redaction regression tests | Tests that confirm regex-redaction patterns redact (not pass through) sensitive substrings in free-text fields. Includes both positive (sensitive patterns are redacted) and negative (non-sensitive patterns are not falsely redacted) cases. |
| Privacy-store access tests | Tests that confirm privacy-store access requires authentication and produces an access-log entry. |
| Conformance corpus results | Results against the project conformance corpus per `spec/test-vectors/`, confirming the institution's deployment matches the spec's privacy-relevant invariants. |
| Test retention | Test results retained per the build's retention policy (typically the same as the chain's seven-year retention). |

**Demonstrability requirement.** The institution can produce, for any deployed build, the test results that demonstrate the build's privacy controls work. Test results are not regenerated for the demonstration; the original results from the original test run are produced.

**Examiner question this answers.** "How do you know your tokenization actually works on the customer-ID field?" Answer: pull the regression test that runs the customer-ID through the SDK and asserts the output is a token, not the original value. Show the test result from the build that was deployed at the time of the question.

## Artifact 4 — Operations

The configuration, build, and tests prove the design was implemented and verified. The operations log proves it kept working in production.

**What it carries:**

| Element | Content |
|---|---|
| Privacy-store access log | Per Article 32 audit-control requirement: every access to the privacy-store is logged with timestamp, accessor identity, the customer record accessed, and the business reason (e.g., "DSAR fulfillment for request `[id]`," "Article 16 rectification for request `[id]`"). |
| Key-fingerprint reconciliation | The `master.reconciliation_completed` operational events per spec §10.1 P-6. These confirm chain identity remained intact; they are also relevant to privacy because a reconciliation failure could indicate privacy-store key drift. |
| Operational events for privacy operations | When erasure runs, rectification correction-records are appended, or DSAR responses are produced, the institution emits an operational event tying the action to its request identifier. |
| Access-review records | Quarterly review confirming the access list to the privacy-store is current; departures are removed promptly; no excessive privileges accumulated. |

**Demonstrability requirement.** The institution can produce a complete history of privacy-store access for any specified period. The history is searchable by accessor, by customer (token), by time. Anomalous access patterns (off-hours, bulk access, unfamiliar accessor) are surfaced by the institution's monitoring.

**Examiner question this answers.** "Who has accessed this customer's privacy-store record, and why?" Answer: pull the access log filtered by the customer's token. Each access has an accessor, a timestamp, and a documented business reason.

## Artifact 5 — DPIA

The Data Protection Impact Assessment is the document where the institution writes down what it considered, what it decided, and what residual risk it accepted. It is the cognitive trail of the privacy-by-design choice.

**What it carries:**

(See `gdpr-dpia-template.md` for the full template content.) Briefly:

| Section | Demonstrability content |
|---|---|
| Necessity & proportionality | Documents that the institution evaluated the chain's purpose against less invasive alternatives and chose the current design. |
| Data categories | Names the data the chain captures and the privacy-by-design treatment for each. |
| Risks | Names the privacy risks (privacy-store compromise, IKM compromise, retention beyond necessity, data-subject-access scope) and the mitigations. |
| Rights & safeguards | Maps each data-subject right to the institution's procedure for honoring it. |
| Residual risk | Documents what the institution accepts and on what assumptions. |
| Review cadence | Names the annual review and material-change triggers. |

**Demonstrability requirement.** The DPIA is current (reviewed at least annually, and after any material change). The DPIA references the configuration, build, test, and operations artifacts as evidence supporting its claims. The DPIA was reviewed by the DPO and signed off by the appropriate executive (typically the Chief Risk Officer or Chief Information Security Officer).

**Examiner question this answers.** "Did you think about the privacy implications before deploying?" Answer: produce the DPIA with the date it was first completed and the date of the most recent review. The DPIA's risk analysis is concrete (named risks, named mitigations, named assumptions), not boilerplate.

## How the artifacts assemble

The five artifacts are not independent boxes; they reference each other and together support a complete demonstration:

- The DPIA names the configuration and references the test results.
- The configuration is deployed via the build, which is provenance-attested.
- The tests run against the build and produce the artifacts the DPIA references.
- The operations log records the actual production behavior of the deployed build.

When an examiner asks the Article 25 question, the institution answers in the order the artifacts logically chain. "Here is the DPIA showing what we considered. Here is the configuration we chose. Here is the build that deployed it. Here are the test results showing it works. Here is the operations log showing it kept working." Five artifacts, one chain of evidence.

## Where the demonstrability pack lives

The pack is assembled on demand. The components live where the institution naturally stores them:

- Configuration in version control (Git, typically)
- Build artifacts in the build server / artifact registry / supply-chain attestation store
- Test results in the build server / test-result archive
- Operations logs in the institution's centralized logging / SIEM
- DPIA in the institution's GRC tool / privacy-program document store

The pack is not a single static document. It is a procedure: when asked, retrieve the components, confirm they are mutually consistent, and present them as a coherent set. The institution's privacy team owns the procedure; the chain-operations team supports the build and operations components.

## Annual demonstrability rehearsal

Once per year, the privacy team rehearses the pack assembly. The rehearsal:

1. Picks an arbitrary date in the past 12 months.
2. Asks the team to produce the five artifacts as if responding to an examiner question dated for that day.
3. Confirms the artifacts are coherent (the build references a configuration version that is still retrievable; the test results match the build; the operations log covers the period).
4. Documents the rehearsal in the institution's privacy-program records.

The rehearsal surfaces gaps before an examiner does. A gap discovered in rehearsal is a routine operational issue; a gap discovered during examination is a finding.

## Cross-references

- `privacy-by-design.md` — the SDK configuration whose schema artifact 1 captures
- `supply-chain.md` — the build provenance that artifact 2 relies on
- `spec/test-vectors/` — the conformance corpus that artifact 3 references
- `gdpr-dpia-template.md` — the DPIA template that artifact 5 instantiates
- `gdpr-dpo-consultation.md` — the DPO review that signs off the artifacts
- `gdpr-ropa-template.md` — the RoPA entry that names the security measures these artifacts demonstrate
- spec §10.1 — the operational events that artifact 4 records
