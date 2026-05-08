# Privacy by design

> **What this doc is.** Privacy patterns for the chain. Closes round-3 GDPR alignment partial and fresh-batch privacy-criterion partials. Addresses how the chain's integrity property composes with privacy obligations.

## The fundamental tension

The chain is append-only, integrity-bearing. Privacy regulations (GDPR, CCPA, state-level laws) frequently require erasure or modification of records. These appear in conflict.

Resolution: **the chain captures hashes or tokens of personal information, never the personal information itself.** Personal information lives in a separate store the institution can govern under privacy regulations. The chain's integrity property is preserved; the personal information is governable.

## Design pattern

### What the chain captures

The chain captures the canonical event payload (per `02-chain-construction.md`). The institution configures the SDK to:

1. Identify fields containing personal information
2. Replace those fields with deterministic tokens before canonicalization
3. Store the original-field-to-token mapping in a separate, privacy-governed store

The chain payload then contains tokens, not personal information. The token-to-original mapping is the institution's separate privacy-store that responds to erasure requests.

### Token shape

Tokens are deterministic but irreversible without the mapping. Two patterns:

**HMAC tokenization.** `token = HMAC-SHA-256(privacy_key, original)`. The privacy_key lives in a separate HSM-protected store; without it, tokens are not reversible. The privacy_key is rotatable independently of the chain's master key.

**Reversible-encryption tokenization.** The original is encrypted with a privacy-store key; the encrypted value is the token. Erasure deletes the privacy-store record; the token in the chain remains but is no longer reversible without the deleted key.

### What's in the chain

| Field type | What's captured |
|---|---|
| Customer identifier (account number, SSN) | Token (not original) |
| Email address, phone number | Token (not original) |
| Free-text content potentially containing PII | Tokenized or redacted at SDK |
| Decision-relevant facts (loan amount, decision) | Original (not personal-information-bearing in most jurisdictions) |
| Timestamps, technical metadata | Original |
| Model-state (model id, parameters) | Original |

### Erasure flow

1. Customer requests erasure under GDPR / CCPA / similar
2. Institution identifies all customer-related tokens via its customer-correlation index
3. Institution erases the customer's records in the privacy-store (the original-to-token mapping)
4. The chain's records remain (tokens preserved, integrity preserved)
5. Post-erasure, the tokens in the chain no longer resolve to identifying information
6. The institution documents the erasure in its privacy-log

## GDPR alignment

### Recital 26 — pseudonymization

"The principles of data protection should apply to any information concerning an identified or identifiable natural person. Personal data which have undergone pseudonymisation, which could be attributed to a natural person by the use of additional information should be considered to be information on an identifiable natural person."

Tokens with a separate mapping store are pseudonymized data. With the mapping erased, the data is anonymized (no additional information available). The chain's tokens transition from pseudonymized to anonymized at erasure.

### Article 5(1)(c) — data minimization

The chain captures what is necessary for integrity-bearing audit. Free-text content potentially containing PII is tokenized at the SDK before canonicalization. The chain's payload-completeness assertion (per `legal-evidence-procedures.md`) describes what was captured; the institution's privacy-impact assessment evaluates necessity.

### Article 5(1)(c) data-minimisation balance

The chain's payload-completeness principle (spec §4.4) captures the decision-relevant attributes plus the audit-context the institution needs to reproduce the AI decision under SR 11-7 model-risk-management discipline (spec §3). A stricter reading of Article 5(1)(c) might argue only the decision and timestamp are strictly necessary; everything else (model inputs, parameters, reasoning chain) is excessive. The balancing test the institution conducts in its DPIA arrives at the opposite conclusion for the chain's defined purpose.

The reasoning runs as follows. The chain's purpose is tamper-evident, comprehensive evidence of what the AI decided AND on what basis. A reduced capture (decision plus timestamp only) would not support reproduction of the decision under the same inputs, would not support fairness audits stratified by decision-relevant attributes, and would not support the customer-dispute hallucination cross-check (`customer-dispute-procedures.md` §"Hallucination cross-check") because there would be nothing in the chain to cross-check against the institution's authoritative records. The reduced capture defeats the audit purpose the chain exists to serve. Comprehensive capture is necessary for that purpose, not over-inclusive relative to it.

The minimisation work happens at the field layer, not the event layer. The institution's privacy configuration (`privacy-by-design.md` §SDK configuration) names which fields are tokenized, which are redacted, which are retained as originals. PII in the customer identifier becomes a token before canonicalization; PII inadvertently captured in free-text reasoning is redacted at the SDK or replaced with a redaction-with-hash placeholder so the chain still binds the redaction without binding the original text. The chain captures the full decision context with personal data minimised to tokens and hashes — the comprehensive scope serves the audit purpose and the field-level controls keep the personal-data exposure proportionate.

The institution documents the minimisation balance in its DPIA (Article 35) and reviews it annually. Triggers for re-evaluation include: a change to the institution's privacy configuration that adds or removes tokenized fields, regulatory guidance that narrows the minimisation boundary (an updated EDPB opinion, a national DPA enforcement action, an evolved interpretation of "necessary" under Article 5(1)(c)), or a fairness-audit finding that the captured fields are insufficient for the audit's stratification surface. The DPIA records the balancing-test inputs and conclusions so the institution can defend the comprehensive-capture posture in a regulatory inquiry.

### Article 9 special-category data handling

GDPR Article 9(1) prohibits processing special-category personal data except under enumerated Article 9(2) exceptions. The chain's `audit.*` payload may contain inferred special categories the institution does not explicitly request from the customer but which arise from decision-relevant features the AI consumes. Examples: marital status inferred from joint-account or dependent fields; age proxies inferred from credit-history depth or product-eligibility flags; disability inferred from payment patterns (irregular cadence aligned with disability-benefit deposit cycles); health data inferred from healthcare-financed loan applications (the lender knows the loan funds a medical procedure); education data, which is a special category in some EU member states for student-finance products.

The institution's privacy configuration (per §SDK configuration above) identifies which fields are or may contain special-category data and applies tokenization to them before canonicalization. Health data on healthcare-financed loans is tokenized; behavioral-pattern features that proxy for disability are tokenized; education data on student-finance products is tokenized when the institution operates in member states where it is a special category. Tokenization keeps the chain's integrity property over the decision while the special-category content lives in the privacy-store under tightened access control.

The lawful-basis analysis for the chain (Article 6) is paired with the Article 9 exception analysis. The institution confirms a valid Article 9(2) exception covers the special-category processing — typically Article 9(2)(b) (legal obligation in employment, social security, or social protection law) for credit decisions tied to a regulated activity, plus Article 9(2)(a) (explicit consent) for healthcare-financed lending under HIPAA where explicit consent is the controlling rule. The institution documents the exception per processing activity in its RoPA and per data category in its DPIA.

The DPIA explicitly addresses special-category risk. The risk register names: which fields may carry inferred special categories, the tokenization control applied to each, the residual risk if the privacy-store is compromised (the tokens become reversible to the special-category content), and the institution's incident-response posture for a privacy-store breach involving special-category data (heightened notification thresholds under GDPR Article 33-34, plus state-law healthcare-data-specific notification rules). The institution's DPO consultation per Article 38 covers the Article 9 analysis; the DPO's opinion on the special-category exception is documented in the DPIA.

### Article 17 — right to erasure

Erasure of the privacy-store records de-identifies the chain's tokens. The chain's data, once tokens are de-mapped, is no longer personal data within the meaning of GDPR. The institution's erasure log demonstrates compliance.

### Article 25 — privacy by design

The pattern is privacy-by-design: tokenization happens at SDK time, before personal data enters the chain. The chain itself never holds personal data.

## CCPA / state-law alignment

### Right to deletion (CCPA §1798.105)

Same pattern as GDPR Article 17. The privacy-store records are deletable; the chain's tokens persist as anonymized data.

### Right to know (CCPA §1798.110)

The institution responds with a copy of the privacy-store records relating to the consumer. The chain itself is not the response artifact (it's institution-internal); the privacy-store is the response artifact.

### State variations

Different states (Virginia VCDPA, Colorado CPA, Connecticut CTDPA, Texas TDPSA) impose similar requirements with variations. The pattern (tokenization + separate privacy-store + erasure) accommodates all of them.

## Cross-border data transfers

For institutions transferring chain data across jurisdictions:

- The chain's tokens are not personal data once mapping is erased; cross-border transfer of tokenized data may not trigger transfer-protection requirements
- The privacy-store is residency-controlled; the institution operates the privacy-store in the appropriate jurisdiction
- Standard contractual clauses (SCCs) and adequacy decisions apply to the privacy-store, not the chain

The chain's data-residency-neutrality (`docs/dr-and-resilience.md`) supports per-jurisdiction privacy-store configuration.

## Privacy-key and chain-master rotation coordination

The privacy-store key and the chain's master IKM serve different threat models but interact compositely. The privacy-store key reverses tokens to PII; the chain-master IKM authenticates chain entries. An attacker holding both keys at the same time can map tokens back to PII AND forge new chain entries with valid tokens. The two-key compromise is more dangerous than either single-key compromise because it enables fabrication of evidence that looks authentic and identifies real customers.

Rotation is the institution's primary defense against undetected key compromise. A simultaneous rotation of both keys — even when each rotation succeeds individually — doubles the operational at-risk window: any operational error during the overlap (a botched rotation, a stale roster entry, a missed handshake) extends the window during which both keys are in an indeterminate state. The institution's detection latency for either compromise increases because the reconciliation signals from both keys arrive in the same operational noise.

The two rotations MUST NOT occur in the same operational window. The institution's chain-operations team and privacy-data-protection team coordinate via the change-management calendar so that:

1. Each rotation has its own scheduled window with its own pre-rotation roster snapshot, post-rotation reconciliation, and post-rotation evidence-of-success record (per spec §10.1 for the chain-master, per the institution's privacy program for the privacy-key).
2. The two scheduled windows do not overlap. Specifically, no privacy-key rotation begins until the chain-master rotation has completed reconciliation (the `master.reconciliation_completed` event records `unmatched_count = 0` after the new IKM generation is in steady state), and no chain-master rotation begins until the privacy-key rotation has completed its parallel reconciliation step.

**Operational window definition.** "Same operational window" means the same calendar week, OR — when the institution operates either rotation cadence faster than weekly — until both rotations have completed reconciliation under the institution's standard reconciliation procedure. The institution's change-management calendar records both rotation windows; the institution's IR program references the calendar when investigating any post-rotation anomaly.

**Failure-mode-of-coincidence.** A single rotation operational error during overlap doubles the at-risk window for both keys. The combined detection latency increases because the institution's reconciliation evidence for both keys arrives mixed in the same operational period — the chain-operations team and the privacy-data-protection team each see noise in their own reconciliation signal that is not separable from the other team's rotation activity. The institution loses the per-key root-cause attribution that non-overlapping rotations preserve. The non-overlap rule is the cheapest control that preserves attribution.

The institution documents the coordination procedure in its control description; the SOC team confirms via sample of the change-management calendar that no two rotations were scheduled in the same operational window during the period.

## SDK configuration

The institution's SDK configuration declares what gets tokenized:

```yaml
privacy:
  tokenization:
    fields:
      - path: "event.customer_id"
        method: hmac
        privacy_key_label: tenant-acme-privacy-key-v1
      - path: "event.amount_payee"
        method: hmac
        privacy_key_label: tenant-acme-privacy-key-v1
      - path: "gen_ai.request.messages[*].content"
        method: redaction-with-hash
        privacy_key_label: tenant-acme-privacy-key-v1
    privacy_key_rotation_cadence: quarterly
    privacy_store_endpoint: https://privacy-store.bank.internal/v1
```

The configuration is part of the institution's control description. Changes go through change management.

## TSC Privacy criterion (SOC 2)

For SOC 2 with Privacy criterion claimed:

| Privacy Criterion | Chain support |
|---|---|
| **P1 (Notice)** | Chain composes; institution's notice describes AI use, including chain operation |
| **P2 (Choice and consent)** | Institution governs; chain neutral |
| **P3 (Collection)** | Tokenization at collection; chain captures tokens not originals |
| **P4 (Use, retention, disposal)** | Privacy-store erasure de-identifies chain tokens |
| **P5 (Access)** | Privacy-store responds to access requests; chain provides integrity assurance |
| **P6 (Disclosure to third parties)** | Privacy-store responds; chain integrity protects against unauthorized disclosure of records |
| **P7 (Security for privacy)** | Standard security applies to privacy-store; chain provides integrity layer |
| **P8 (Quality)** | Chain provides integrity-bearing record |

The chain is not a Privacy control by itself; it composes with the institution's Privacy program.

## Privacy-impact assessment

When adopting the chain, the institution's privacy team conducts a privacy-impact assessment (PIA) covering:

1. What personal data is captured (in the privacy-store) versus what tokens are captured (in the chain)
2. How tokens map to personal data (privacy-store schema)
3. Privacy-store retention vs chain retention (typically privacy-store retention is shorter)
4. Erasure workflow (how erasure affects chain tokens — they become anonymous)
5. Transfer flows (chain to backup, chain to verifier, chain to SOC team)

The PIA is part of the institution's privacy program documentation. Updated when the SDK configuration changes.

## Operational pattern

| Operation | Privacy-store | Chain |
|---|---|---|
| Customer interaction | Records personal data | Records tokens of personal data |
| Customer requests erasure | Deletes record | Tokens become anonymous |
| SOC team review | Reviews privacy-store separately under privacy controls | Reviews chain for integrity |
| Examination | Examiner sees chain (with tokens); examiner does not see privacy-store unless required | |
| Discovery | Privacy-store responds to discovery; chain confirms record integrity | |

## Audit considerations

SOC and examination teams test:

- The institution operates the privacy-store as a separate control
- The privacy-store records are erasable; the chain's tokens remain anonymous post-erasure
- The privacy-store key is rotatable; rotation is documented
- The PIA is current
- Cross-border transfer flows are documented

## Limitations

This pattern works for personal-data-as-fields. It does not work cleanly for:

- Free-text content where personal information is interspersed with non-personal information (the SDK can apply NER-based redaction with hash, but the redacted content may not be reversibly de-identifiable)
- Personal data that is itself decision-relevant (e.g., customer's name being used in the AI's reasoning); the institution's privacy program addresses

For these cases, the institution applies institution-specific procedures within its privacy program. The chain's integrity property is preserved; privacy responsiveness depends on the institution's broader procedures.
