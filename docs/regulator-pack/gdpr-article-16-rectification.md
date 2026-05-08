# GDPR Article 16 — right to rectification

> **What this doc is.** The institution's procedure for honoring an Article 16 rectification request when the data subject says the AI decision relied on inaccurate facts. The chain is immutable; a direct edit would break MAC and Merkle integrity. The solution is the same shape as the rest of the chain's design: append a correction-record entry parent-linked to the original decision, update the privacy-store with the corrected fact, and notify the customer. The chain entry stays exactly as it was — that is the historical record of what the AI saw at the time. The correction-record records what the institution now knows. Article 16 is satisfied; chain integrity is preserved.

## What Article 16 grants

GDPR Article 16 grants the data subject the right to obtain "without undue delay the rectification of inaccurate personal data concerning him or her." The data subject also has the right to have incomplete personal data completed, including by means of providing a supplementary statement.

"Without undue delay, and in any event within one month" is the Article 12(3) timeline. The institution's commitment refines the milestones: investigation completes within 14 days; correction-record creation within 7 days of investigation conclusion; decision-impact review within 30 days. Total: 30 days from receipt to fully closed-out request, fitting within Article 12(3).

Rectification is not erasure. The customer is not asking the institution to forget; the customer is asking the institution to correct. The chain's append-only structure accommodates correction natively — by appending the correction, not by editing the original.

## The procedure

```
1. Receive          → rectification request; identity verified
2. Investigate      → fact-verification against authoritative sources (14 days)
3. Outcome A        → facts confirmed accurate; respond and close
   Outcome B        → facts found inaccurate; proceed to correction
4. Correction       → append correction-record entry (7 days)
5. Privacy-store    → update privacy-store with corrected fact
6. Decision review  → evaluate impact on the original decision (30 days)
7. Notify           → customer notification of outcome
8. Privacy log      → record fulfillment
```

## Step 1 — Receive and verify

The customer submits the rectification request through the institution's standard channel. The privacy team:

- Logs the request with date, requester identity, and the specific factual claim ("the AI's decision said my income was X; my income is actually Y").
- Verifies the requester's identity per the institution's standard authentication.
- Acknowledges receipt within five business days.
- Triages the request: which (tenant_id, run_id) pair(s) are touched? Which fact is challenged? Is the challenged fact in the privacy-store, the chain, or both?

Most rectification requests touch a fact in the privacy-store; the chain entry contains the tokenized version. Some requests touch a fact captured in the chain directly (a non-tokenized decision factor like the AI's assessment of debt-to-income ratio). The triage determines which path the procedure takes.

## Step 2 — Investigate (14 days)

The privacy team works with the AI program owner and the customer-service team to verify the facts. The investigation:

| Investigation step | Source |
|---|---|
| Pull the chain entry for the decision in question | Chain ledger via verifier or query interface |
| Pull the privacy-store record for the affected customer | Privacy-store via standard extract tool |
| Identify the specific fact challenged | Chain entry's audit-namespace fields; privacy-store fields |
| Cross-check against authoritative sources | Customer's primary-source documentation (if provided), institution's account systems, third-party data sources (credit bureau, employment verification), customer's submitted application |
| Determine the accurate fact | Authoritative source resolution; if multiple sources disagree, the institution names which source it is treating as authoritative |
| Document the investigation | Investigation log with dates, sources consulted, conclusions |

The 14-day commitment is a service-level target; the substantive standard is Article 12(3) "without undue delay." Most investigations close in days; complex ones (multi-source disagreement, third-party verification needed) may use the full 14.

## Step 3 — Outcomes

### Outcome A — facts confirmed as accurate

The institution determines the facts as recorded match the authoritative source. The chain entry stands as captured; the AI relied on accurate inputs.

The institution responds to the customer:

> Your rectification request received on 2026-05-07 has been investigated. Our records show the fact you challenged ("[fact]") matches our authoritative source ("[source]") as of the date of the decision. The AI decision was based on accurate information. We have logged your inquiry and the investigation's findings. If you have additional supporting documentation that supports the change you are requesting, please submit it and we will reopen the investigation.

The chain is not modified. The privacy log records the request and the outcome (Article 16 inquiry, no rectification required).

### Outcome B — facts found inaccurate

The institution determines the facts as captured did not match the authoritative source, or the customer's documentation establishes a more accurate version. Rectification is required.

Proceed to step 4.

## Step 4 — Append correction-record entry (7 days)

The institution appends a correction-record chain entry. This is not a modification of the original entry; it is a new entry parent-linked to the original, recording the correction.

The correction-record uses the spec's translation-entry schema as a model (spec §10.11 documents the parent-linked translation pattern; correction-record extends it for rectification semantics):

| Field | Content |
|---|---|
| `audit.correction.original_decision_run_id` | The run_id of the decision being corrected |
| `audit.correction.original_decision_seq` | The seq of the original entry within that run |
| `audit.correction.original_decision_chain_hash` | The chain_hash of the original entry; binds the correction to the exact original record |
| `audit.correction.inaccurate_fact` | The specific fact as originally captured (tokenized if PII) |
| `audit.correction.accurate_fact` | The corrected version of the fact (tokenized if PII) |
| `audit.correction.fact_source` | The authoritative source the institution relied on for the correction |
| `audit.correction.investigation_id` | The institution-internal investigation log identifier |
| `audit.correction.impact_on_decision` | Enum: `decision_unaffected` / `decision_affected_pending_review` / `decision_should_be_reversed` |
| `audit.correction.requested_by` | "data_subject" (Article 16) — distinguishes from institution-initiated corrections |
| `audit.correction.requested_date` | When the request was received |
| `audit.correction.completed_date` | When the correction-record was appended |

The correction-record is itself an immutable chain entry. It carries the standard chain-entry properties: HMAC, Merkle inclusion at the next seal, HSM signature on the seal-day root. The integrity guarantees apply to the correction-record as they do to any other entry.

The customer can verify the correction-record exists by following the parent link from the original entry. Independent verifiers (regulators, courts, the customer's counsel) can confirm both entries are intact and their relationship is recorded.

## Step 5 — Update privacy-store

The chain holds tokens; the privacy-store holds the originals. When the corrected fact is in the privacy-store, the privacy-store is updated to reflect the corrected version. The update is logged in the privacy-store's own audit log:

| Privacy-store update field | Content |
|---|---|
| Previous value | The inaccurate value as previously stored |
| New value | The corrected value |
| Update timestamp | When the correction was applied |
| Source authority | The authoritative source name |
| Updater identity | The privacy officer who applied the correction |
| Cross-reference | The correction-record entry's chain location |

The privacy-store is now the authoritative source-of-truth for the customer's personal data. Subsequent uses (DSAR responses, future AI decisions, dispute defense for the original decision) read from the privacy-store and find the corrected value. The chain entry from the original decision is preserved as historical record of what the AI saw at the time, which matters for fairness audits and dispute defense.

## Step 6 — Decision-impact review (30 days)

A correction may or may not change the decision the AI made. The institution evaluates:

| Decision-impact category | Action |
|---|---|
| `decision_unaffected` | The corrected fact does not change the decision (the AI's logic was not load-bearing on this fact, or the correction is small enough to leave the decision in the same outcome class). The institution closes the review with a documented finding. |
| `decision_affected_pending_review` | The corrected fact may change the decision; a re-evaluation is warranted. The institution re-runs the decision (same model version, corrected inputs) or routes to manual review per the institution's adverse-action procedures. |
| `decision_should_be_reversed` | The corrected fact clearly indicates the original decision was wrong. The institution reverses the decision, applies the corrective action (e.g., approves the loan that was incorrectly denied, refunds the pricing differential), and notifies the customer. |

For `decision_should_be_reversed` specifically, the institution operates the same workflow as a customer-dispute outcome that resolves in the customer's favor (`customer-dispute-procedures.md`). The reversal is documented; the corrective action is taken; the customer is notified.

The review's outcome is itself a chain event in the institution's broader audit trail (a separate entry capturing the review and its outcome) but does not modify the original decision entry or the correction-record.

## Step 7 — Notify the customer

Within 30 days of receipt, the institution notifies the customer:

| Outcome | Notification content |
|---|---|
| Outcome A (facts accurate) | Investigation result; explanation of the authoritative source; invitation to provide additional documentation |
| Outcome B + decision_unaffected | Correction recorded; the AI decision stands; explanation that the correction does not change the outcome |
| Outcome B + decision_affected_pending_review | Correction recorded; re-evaluation in progress; estimated completion date |
| Outcome B + decision_should_be_reversed | Correction recorded; original decision reversed; corrective action being taken; explanation of next steps |

The notification includes the institution's reference to the correction-record (so the customer can ask for the chain extract if they want to verify the correction was made).

## Step 8 — Privacy log

The institution records the rectification in its privacy log:

| Field | Content |
|---|---|
| Request ID | Institution-internal request identifier |
| Date received | When the request arrived |
| Date investigation completed | When step 2 closed |
| Outcome | Outcome A / Outcome B |
| Decision impact (Outcome B only) | `decision_unaffected` / `decision_affected_pending_review` / `decision_should_be_reversed` |
| Correction-record chain location | (tenant_id, run_id, seq) of the correction-record entry |
| Privacy-store updated | Yes / No |
| Reversal completed | Yes / No / Pending |
| Customer notification date | When step 7 was completed |
| Officer | Identity of the privacy officer who completed the operation |

## Worked example — accurate facts (Outcome A)

A customer disputes the AI's stated income figure on a 2025 loan application. Investigation pulls the original application and the institution's authoritative source (the customer's W-2 submitted at application time). The W-2 figure matches the AI's input.

The privacy team responds:

> Your rectification request has been investigated. Our records show your income figure on the loan application matches the W-2 documentation you submitted at the time. The AI decision was based on the figure you provided. If you have updated W-2 or pay-stub documentation reflecting a different figure, please share it and we will reopen the investigation. We have logged your inquiry under request ID `[id]`.

The chain is not modified. The privacy log records "Article 16 inquiry, no rectification required."

## Worked example — inaccurate facts, decision unaffected (Outcome B)

A customer disputes the AI's stated employment-tenure figure. Investigation finds the institution's record was 18 months; the customer's employer-letter confirms 24 months at the time of the application.

The institution appends a correction-record:

```
audit.correction.original_decision_run_id = r_2025-08-22_loan_app_4471
audit.correction.original_decision_seq = 14
audit.correction.original_decision_chain_hash = [hash from original entry]
audit.correction.inaccurate_fact = "employment_tenure_months: 18"
audit.correction.accurate_fact = "employment_tenure_months: 24"
audit.correction.fact_source = "employer_verification_letter dated 2025-08-15"
audit.correction.investigation_id = inv_2026-05-09_3389
audit.correction.impact_on_decision = decision_unaffected
audit.correction.requested_by = data_subject
audit.correction.requested_date = 2026-05-07
audit.correction.completed_date = 2026-05-14
```

The privacy-store is updated. The decision-impact review confirms the outcome would not change (the AI's logic was satisfied at both 18 months and 24 months for this customer's risk profile). The customer is notified:

> Your correction has been recorded. Our authoritative records have been updated to reflect the 24-month employment tenure. We have reviewed the original decision; the corrected information does not change the outcome. The original decision stands as made.

## Worked example — inaccurate facts, decision reversed (Outcome B)

A customer disputes the AI's stated debt-to-income ratio. Investigation finds the institution misread a paid-off loan as still active, inflating the customer's reported debt. The corrected debt-to-income ratio falls below the threshold the AI used to deny the loan.

The institution appends the correction-record with `impact_on_decision = decision_should_be_reversed`. The privacy-store is updated. The decision is re-evaluated; the AI (or manual review) confirms approval at the corrected ratio. The customer is notified:

> Your correction has been recorded. The original decision was based on an inflated debt-to-income ratio that included a loan you had already paid off. With the corrected information, our re-evaluation approves the loan you originally applied for. We are processing the approval now and a representative will contact you with next steps within five business days.

The reversal is itself an institution-side event documented in the chain (a separate run capturing the re-evaluation). The original denial entry remains as historical record. The correction-record links the two. The customer's audit trail is now complete: original denial → correction request → corrected fact → reversal approval.

## Why the correction-record approach satisfies Article 16

Article 16 requires the inaccurate data to be rectified. The institution's authoritative source-of-truth (the privacy-store) is updated; that is rectification. The chain's historical record of what the AI saw at the time is preserved; that is integrity. The two are not in conflict because the chain is institution-internal evidence, not the source-of-truth.

The customer's interest is in the corrected outcome. The privacy-store update gives them that. The chain's integrity preservation protects against tampering claims that would weaken the institution's defense in fairness audits or examination — a protection that benefits the customer as much as the institution because it ensures the institution's records remain trustworthy.

## Cross-references

- `gdpr-article-17-procedures.md` — companion procedure for erasure
- `gdpr-dsar-fulfillment.md` — companion procedure for access; the DSAR response shows the privacy-store as it stands today (post-correction if Article 16 has been honored)
- `customer-dispute-procedures.md` — broader dispute procedure including non-Article-16 cases (institution-initiated corrections, customer-disputed decisions where no factual inaccuracy is named)
- `gdpr-ropa-template.md` — RoPA entry that names rectification in the data-subject-rights field
- `privacy-by-design.md` — tokenization design that places authoritative customer data in the privacy-store
- spec §10.11 — translation-entry schema that the correction-record extends
