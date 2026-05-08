# GDPR Article 17 procedures — erasure response posture

> **What this doc is.** The institution's operational response posture for GDPR Article 17 ("right to erasure") requests against AI decisions captured in the chain. The chain is append-only and integrity-bound. A traditional "delete this row" would break the MAC and the Merkle root. This doc names the only safe path: erasure happens in the privacy-store, not the chain. Tokens stay; the mapping that makes them reversible goes away. That is GDPR-compliant de-identification, not deletion.

## The tension Article 17 creates

GDPR Article 17(1) gives the data subject the right to obtain erasure of personal data without undue delay. Article 17(3) carves out exceptions, including 17(3)(b): processing necessary for compliance with a legal obligation. The chain is retained for seven years to satisfy FFIEC banking regulations (12 CFR 1002 ECOA records, OCC/Fed/FDIC examination records, CFPB supervisory records — see `gdpr-lawful-basis.md` for the full Article 6 + 17(3)(b) basis). During that retention window, the institution may refuse erasure of the chain entries themselves; the legal-obligation exception applies.

But the institution does not get to refuse forever. Outside the retention mandate, or upon a documented institution-side decision to release a record early, an Article 17 request must be honored. The chain's integrity-bearing design must accommodate erasure without breaking. The privacy-by-design.md tokenization pattern is what makes that possible.

## The four-step response

```
1. Receive request   → privacy team logs the inbound DSAR/erasure request
2. Retention check   → still within 7-year FFIEC mandate? document refusal per 17(3)(b) and 17(4)
3. Out-of-mandate    → delete the privacy-store mapping; chain tokens become anonymous
4. Privacy log       → record date, customer ID token, record count, retention rationale
```

The chain itself is never modified. Erasure in this design means severing the token-to-PII mapping in the privacy-store. Once the mapping is gone, the chain entries — which carry only the tokens — are anonymized data under GDPR Article 4(1) and Recital 26 (see `gdpr-pseudonymization-vs-anonymization.md` for the analysis). The chain stays append-only; the integrity proofs continue to verify; the customer's personal data is no longer present in the chain in any form that can be traced back to them.

## Step 1 — Receive the request

The customer (or their authorized representative) submits an erasure request through the institution's standard DSAR intake channel. The privacy team logs the request with:

- date and time of receipt
- requester identity and authorization (if a representative)
- customer identifier the institution can correlate to (tenant_id, run_id) pairs via the customer-correlation index (CUI; see `gdpr-dsar-fulfillment.md`)
- scope of the request (full erasure, partial erasure, specific decision events)

GDPR Article 12(3) starts the clock at receipt: the institution has 30 days to respond, extensible by 60 days for complex requests with notification to the data subject within the first 30. An erasure request that touches multiple AI decisions across multiple runs is "complex" in this sense; the institution may use the extension if the CUI lookup and the retention assessment require it.

## Step 2 — Retention check

The privacy team confirms the retention status of every (tenant_id, run_id) pair the request touches. The check uses the institution's retention register, not the chain. For each pair the answer is one of:

| Status | Article 17(3)(b) result | Action |
|---|---|---|
| Within 7-year FFIEC mandate, no litigation hold | Refuse erasure under 17(3)(b) (legal obligation) | Document refusal per 17(4) |
| Within 7-year mandate, litigation hold attached | Refuse erasure under 17(3)(b) + 17(3)(e) (legal claim) | Document refusal; cite the specific litigation matter |
| Outside the 7-year mandate, no litigation hold | Erasure required | Proceed to step 3 |
| Outside the mandate, but business continuity case (e.g., active customer relationship) | Refuse under 17(1)(a) only if processing is still "necessary" for the original purpose | Legal team review required; document rationale |

Most requests in the first seven years of operation will fall into the first row. The institution refuses, but the refusal is documented per Article 17(4): the data subject is notified of the refusal, the legal basis (17(3)(b)), and the retention period after which erasure becomes possible.

The retention-check output is itself a privacy-log entry. It carries the request identifier, the (tenant_id, run_id) pairs evaluated, the disposition for each, and the responsible privacy officer's identity.

## Step 3 — Out-of-mandate erasure

When the retention check determines erasure is required, the privacy team executes the privacy-store deletion:

1. Look up the customer's privacy-store records via the CUI.
2. For each (tenant_id, run_id) pair within scope, identify the privacy-store mapping rows.
3. Execute a hard delete on the mapping rows (not a logical-delete flag — actual row removal, with the deletion logged in the privacy-store's own audit log).
4. Confirm the deletion by attempting to reverse one of the affected tokens; the lookup must return "no mapping" rather than the original PII.
5. Record the operation count: how many mapping rows were deleted, across how many runs.

The chain itself is untouched. The chain's integrity-proof artifacts (HMAC-SHA-256 per entry, daily Merkle seal, HSM Ed25519 signature) continue to verify because the chain content has not changed. The tokens in those entries no longer reference any natural person's identity; the data has been irreversibly de-identified.

This is what the rest of GDPR considers anonymization. Recital 26 says anonymization renders data outside GDPR scope. The chain entries that remain are records of AI decisions, but they are not personal data once the mapping is gone.

## Step 4 — Privacy log documentation

Every erasure operation is recorded in the institution's privacy log. The log entry carries:

| Field | Content |
|---|---|
| Request ID | Institution-internal request identifier |
| Date received | When the data subject's request arrived |
| Date completed | When the erasure was executed |
| Customer token | The customer's ID token (NOT the original ID — the privacy log itself must not reintroduce the PII) |
| Tenant scope | The tenant_id values touched |
| Records erased | Count of privacy-store mapping rows deleted |
| Disposition | "Erased" or "Refused — within retention mandate" |
| Refusal basis | If refused: 17(3)(b) + the specific FFIEC retention reference |
| Officer | Identity of the privacy officer who executed the operation |

The privacy log itself is retained per the institution's records-retention policy. It supports examiner inquiries about the institution's Article 17 response posture and provides evidence that erasure decisions are deliberate and documented.

## Worked example — request received during retention period

A customer in 2026-05-07 requests erasure of all AI decisions about them. The institution's retention register shows the affected decisions are from 2024-08 through 2026-04 — all within the seven-year FFIEC mandate.

The privacy team's response:

> Your erasure request has been received and reviewed. The AI decision records you have referenced are subject to retention requirements under FFIEC banking regulations (12 CFR 1002.12 ECOA records). Under GDPR Article 17(3)(b), we are unable to erase these records during the retention period, which expires for the earliest record on 2031-08-15 and the latest on 2033-04-22. We have logged your request and will execute erasure when each record reaches the end of its retention period unless a litigation hold is then in effect. You may submit a renewed request at that time or contact our Data Protection Officer if you have additional questions.

The privacy log records the request, the retention-based refusal, and the calendar dates at which erasure becomes possible. No chain entries are modified. No privacy-store mappings are deleted.

## Worked example — request received outside retention mandate

A customer in 2034-01 requests erasure of AI decisions from 2026. The retention register shows all affected decisions are past the seven-year mandate; no litigation holds are attached.

The privacy team executes the deletion (step 3), confirms it (token lookup returns no mapping), and logs the operation (step 4). The customer is notified within the 30-day Article 12(3) window:

> Your erasure request has been completed. The mapping records in our privacy-store that linked AI decision artifacts to your identity have been deleted. The decision artifacts themselves remain in our integrity-bearing audit trail, but they no longer carry any data that can be associated with you. This is documented in our privacy log under request ID `[id]` dated `[date]`.

## Why the chain is never modified

The chain's integrity guarantees rest on append-only, MAC-bound, Merkle-sealed, HSM-signed entries. Modifying an entry — even to remove personal data — would invalidate the MAC for that entry, invalidate the Merkle root for the seal day that contained it, and invalidate the HSM signature on that seal. The verifier (spec §7) would surface the modification as tampering. The institution would lose the tamper-evidence property that makes the chain useful for ECOA defense, MRM oversight, and examination.

Tokenization plus privacy-store erasure preserves both properties: GDPR-compliant erasure (the personal data is gone) and chain integrity (the chain is unchanged). This is what GDPR Article 25 calls privacy by design — a design choice that makes the privacy outcome possible without compromising the technical control. The demonstrability artifacts are listed in `gdpr-article-25-demonstrability.md`.

## Cross-references

- `gdpr-lawful-basis.md` — Article 6 + 17(3)(b) legal-obligation basis for the retention period
- `gdpr-dsar-fulfillment.md` — customer-correlation index (CUI) and the broader DSAR workflow
- `gdpr-article-16-rectification.md` — companion procedure when the customer disputes accuracy rather than presence
- `gdpr-pseudonymization-vs-anonymization.md` — why post-erasure tokens are anonymized data outside GDPR scope
- `gdpr-dpo-consultation.md` — when the DPO must be consulted on erasure-policy edge cases
- `privacy-by-design.md` — the tokenization + privacy-store design pattern this procedure depends on
- spec §10.9 (IKM retention) and `incident-response-playbook.md` litigation-hold sections — the seven-year mandate basis
