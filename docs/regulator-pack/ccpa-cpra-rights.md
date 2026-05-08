# CCPA/CPRA Consumer Rights — Chain Interaction

> **What this doc is.** Procedures for honoring California Consumer Privacy Act (CCPA) and California Privacy Rights Act (CPRA) consumer rights when the chain is part of the institution's processing landscape. Written so the institution's privacy team can respond to a verified consumer request within the 45-day statutory deadline without confusion about which artifact is the response and which is internal evidence.

> **Why this matters.** California's privacy regime grants consumers four operational rights that interact with the chain in different ways. A delete request is not a tampering event. A correction request is a chain-append, not a chain-edit. A know request returns privacy-store records, not chain entries. Without a written procedure, the institution will conflate these and either over-disclose chain internals or under-respond to legitimate requests.

---

## 1. The four primary CCPA/CPRA consumer rights

| Right | Statute | What the consumer can demand |
|---|---|---|
| Right to Know | Cal. Civ. Code §1798.110 | Disclosure of categories and specific pieces of personal information collected, sources, business/commercial purposes, third parties with whom shared |
| Right to Delete | §1798.105 | Deletion of personal information the business has collected from the consumer |
| Right to Correct | §1798.106 (added by CPRA) | Correction of inaccurate personal information the business maintains |
| Right to Opt-Out of Automated Decision-Making | §1798.185(a)(16) and forthcoming CPRA regulations | Opt-out of profiling and use of automated decision-making for decisions producing legal or similarly significant effects |

CPRA also adds the Right to Limit Use of Sensitive Personal Information (§1798.121) and the Right to Data Portability (overlap with Right to Know). These follow the same procedural pattern as the four primary rights and are covered as variants below.

---

## 2. The chain in CCPA/CPRA terms

The chain is **business-internal evidence**. It is a record of automated decisions, retained for legal-defense and regulatory-examination purposes, with a 7-year retention period grounded in `retention-justification.md`.

| Question | Answer |
|---|---|
| Is the chain "personal information" under §1798.140(v)? | The chain entries contain pseudonymized references (tokens) that resolve to personal information via the privacy-store. While the privacy-store mapping is retained, the chain entries are personal information |
| Is the chain disclosed in a Right to Know response? | No, by default. The Right to Know response is built from the privacy-store records that name the consumer's personal information directly. The chain is supporting evidence, not the response artifact |
| Is the chain modified for a Right to Delete? | No. The privacy-store mapping is deleted; chain tokens become irreversible. The chain ledger itself is never modified |
| Is the chain modified for a Right to Correct? | No. A correction-record entry is appended to the chain (parent-linked to the original decision). The original entry remains as-is |
| Does opt-out of automated decision-making apply prospectively or retrospectively? | Prospectively. The chain captures historical decisions; those entries are legal-defense evidence and remain |

---

## 3. Right to Know — §1798.110

### 3.1 What the consumer can request

Per §1798.110(a), the consumer may request:

1. The categories of personal information collected.
2. The categories of sources from which personal information was collected.
3. The business or commercial purpose for collecting the personal information.
4. The categories of third parties with whom the business shares personal information.
5. The specific pieces of personal information collected (the "data portability" element).

### 3.2 Fulfillment procedure

| Step | Action | Owner |
|---|---|---|
| 1 | Receive request via consumer-facing portal, designated email, or toll-free number per §1798.130(a)(1) | Privacy team intake |
| 2 | Verify consumer identity per §1798.130(a)(7) regulations (matching information, government-issued ID, or signed declaration depending on data sensitivity) | Privacy team verification |
| 3 | Pull privacy-store records for the consumer (the canonical personal-information source) | Privacy data protection team |
| 4 | Identify categories: contact info, financial info, automated-decision results, etc. | Privacy team |
| 5 | Identify sources: consumer directly, credit bureaus, application forms, etc. | Privacy team |
| 6 | Identify business purposes: credit decisions, fraud prevention, customer service, etc. | Privacy team |
| 7 | Identify third-party recipients: regulators, processors, auditors, courts on subpoena | Privacy team |
| 8 | Compile response package per §1798.130(a)(2): readily usable format, transmittable to another entity | Privacy team |
| 9 | Send response within 45 days of receipt; one 45-day extension allowed if the consumer is notified per §1798.130(a)(2)(B) | Privacy team |
| 10 | Log fulfillment in privacy-log (request date, verification method, response date, scope) | Privacy team |

The chain's role: **none, by default.** The chain is internal evidence; the privacy-store is the response source. If the consumer specifically asks for "automated-decision history" and the institution's privacy notice has named the chain as a source category, the institution may include high-level chain-derived data (e.g., "automated decision: approved on 2026-03-15") drawn from the privacy-store's correlation index. The chain's raw entries are not disclosed.

### 3.3 Edge case — consumer demands chain entries

If the consumer or the consumer's counsel specifically demands chain entries by reference (e.g., during pre-litigation negotiation or via subpoena), the institution treats this as a litigation matter, not a CCPA Right to Know matter. `litigation-support.md` governs the chain's role as evidence; the chain entries are produced under FRCP discovery rules, not CCPA disclosure rules. The consumer's CCPA right is satisfied separately via the privacy-store records.

---

## 4. Right to Delete — §1798.105

### 4.1 What the consumer can request

Deletion of personal information the business has collected from the consumer, subject to the exceptions in §1798.105(d). The exceptions most relevant to the chain:

- §1798.105(d)(1) — necessary to complete the transaction or perform a contract.
- §1798.105(d)(4) — to detect security incidents, protect against malicious or fraudulent activity, prosecute those responsible.
- §1798.105(d)(7) — to comply with a legal obligation.
- §1798.105(d)(8) — to enable solely internal uses reasonably aligned with the consumer's expectations.

The 7-year FFIEC retention obligation (per `retention-justification.md`) is a §1798.105(d)(7) legal-obligation exception during the retention period. The institution's deletion response distinguishes between data outside the retention obligation (deleted) and data within (retained, with reason cited).

### 4.2 Fulfillment procedure

| Step | Action | Owner |
|---|---|---|
| 1 | Receive verified deletion request | Privacy team |
| 2 | Identify the consumer's personal information across all systems (privacy-store, transactional records, marketing systems, etc.) | Privacy team + each data steward |
| 3 | Determine which records are within a §1798.105(d) exception | Legal team |
| 4 | For records outside exceptions: delete from the originating system; if the privacy-store is the source, delete the privacy-store mapping (which renders chain tokens anonymized — see `pseudonymization-vs-anonymization.md`) | Privacy data protection team |
| 5 | For records within an exception: do not delete; document the exception cited | Privacy team |
| 6 | Notify the consumer of the deletion outcome within 45 days, including any exceptions invoked | Privacy team |
| 7 | Log fulfillment | Privacy team |

The chain's role: **the privacy-store deletion is the operative action.** Once the privacy-store mapping is deleted, the chain tokens for that consumer become irreversible and the chain entries no longer constitute personal information. The chain ledger itself is not touched.

### 4.3 Interaction with the 7-year retention period

For chain entries within the 7-year FFIEC retention window, the institution invokes the §1798.105(d)(7) legal-obligation exception. The institution responds: "Your request is partially fulfilled. Personal information outside our regulatory-retention obligations has been deleted. Chain-of-custody records of automated credit decisions are retained for 7 years pursuant to 12 CFR 1002 (ECOA) and FFIEC supervisory record-retention requirements. At the end of the retention period, the linked personal-information mapping will be automatically deleted."

The institution's privacy notice (Article 14 / CCPA §1798.130(a)(5)) discloses this retention practice in advance, so the partial-fulfillment response is not a surprise.

---

## 5. Right to Correct — §1798.106 (added by CPRA)

### 5.1 What the consumer can request

Correction of inaccurate personal information the business maintains. The consumer must provide the asserted-correct information; the institution evaluates against authoritative sources.

### 5.2 Fulfillment procedure

| Step | Action | Owner |
|---|---|---|
| 1 | Receive verified correction request | Privacy team |
| 2 | Pull the affected records (privacy-store + chain entries that referenced the inaccurate data) | Privacy team + chain operations |
| 3 | Investigate: cross-check the asserted facts against authoritative sources (credit bureau, official records, source documents) | Privacy team + relevant business units |
| 4 | If the asserted correction is verified: update the privacy-store record (the source-of-truth) and **append a correction-record chain entry** parent-linked to the original decision | Privacy team + chain operations |
| 5 | If the asserted correction is not verified: respond to the consumer with the basis for the determination per §1798.106(c) | Privacy team |
| 6 | Notify the consumer of the outcome within 45 days | Privacy team |
| 7 | Log fulfillment | Privacy team |

### 5.3 The correction-record chain entry

The correction is **never an in-place edit** to the chain. The chain is integrity-bound; in-place edits break the HMAC and Merkle root. Instead, the institution appends a new chain entry of type `audit.correction` carrying:

| Field | Content |
|---|---|
| `audit.correction.original_decision_run_id` | Run-ID of the original decision being corrected |
| `audit.correction.original_decision_seq` | Sequence number of the original decision entry |
| `audit.correction.inaccurate_fact` | Description of the fact found to be inaccurate (tokenized if PII) |
| `audit.correction.accurate_fact` | Description of the corrected fact (tokenized if PII) |
| `audit.correction.source` | Authoritative source consulted |
| `audit.correction.impact_on_decision` | One of: `decision_unaffected`, `decision_affected_pending_review`, `decision_should_be_reversed` |
| `audit.correction.reviewer` | Identifier of the privacy-team or business-unit reviewer |

The chain now contains both the original decision (untouched) and the correction record (parent-linked). A reader of the chain can reconstruct the decision history including the correction; the integrity property is preserved.

If the correction's `impact_on_decision` is `decision_should_be_reversed`, the institution opens a separate dispute-review process per `customer-dispute-procedures.md`. The correction record is the evidentiary anchor for that review.

---

## 6. Right to Opt-Out of Automated Decision-Making

### 6.1 What the consumer can request

CPRA §1798.185(a)(16) authorizes regulations governing automated decision-making, including a consumer's right to opt-out and a right to meaningful information about the logic involved. Final regulations from the California Privacy Protection Agency (CPPA) define the operational scope; the institution operates conservatively pending final regulations.

The opt-out applies prospectively: the consumer requests that the business not use automated decision-making for future decisions affecting the consumer.

### 6.2 Fulfillment procedure

| Step | Action | Owner |
|---|---|---|
| 1 | Receive verified opt-out request | Privacy team |
| 2 | Flag the consumer in the institution's automated-decision systems | Privacy team + AI product team |
| 3 | Configure those systems to route the consumer's future decisions through a human-review pathway | AI product team |
| 4 | Confirm that the chain continues to log decisions made by the human-review pathway (the chain is decision-system-agnostic; human review is itself a logged decision) | Chain operations |
| 5 | Acknowledge the opt-out to the consumer within 45 days | Privacy team |
| 6 | Log the opt-out as a privacy-store record (and optionally as an `audit.opt_out` chain entry for auditable reference) | Privacy team |

The chain's role: **historical entries are not modified.** Decisions made before the opt-out date remain in the chain as legal evidence. The opt-out changes the institution's prospective behavior, not its historical record.

### 6.3 Right to Meaningful Information

The CPPA regulations are expected to require disclosure of meaningful information about the logic of automated decisions. The chain's role: the chain entry captures the model version, parameters, inputs, and outputs. The institution's disclosure to the consumer is built from this data, abstracted to a non-trade-secret level appropriate for consumer disclosure. The raw chain entry is not disclosed; the institution prepares a derived summary.

---

## 7. The 45-day response timeline

Per §1798.130(a)(2), the institution must respond to a verified consumer request within 45 days of receipt. One 45-day extension is permitted if the consumer is notified within the original 45 days and given the reason for the extension.

The institution's internal cadence:

| Day | Milestone |
|---|---|
| 0 | Request received |
| 1–3 | Identity verification |
| 3–10 | Data assembly across systems |
| 10–25 | Investigation (especially for correction requests) |
| 25–40 | Response drafting and legal review |
| 40–45 | Response sent |

Complex requests (correction with disputed facts, multi-system deletion) trigger the extension and are tracked separately.

---

## 8. Cross-references

- `customer-dispute-procedures.md` — handles the dispute process when a consumer challenges an automated-decision outcome (overlaps with Right to Correct).
- `pseudonymization-vs-anonymization.md` — explains why chain tokens are personal information until the privacy-store mapping is deleted.
- `retention-justification.md` — names the §1798.105(d)(7) legal-obligation exception that limits deletion during the 7-year window.
- `breach-notification-matrix.md` — California's data-breach notification statute (§1798.82) is one of the state-law clocks coordinated through the privacy office.
- `token-vault-architecture.md` — privacy-store custody by a separate team supports privacy-team-only deletion authority.
- `dpia-template.md` §4 — Rights & Safeguards section names this document.
- §§1798.100–1798.199 — operationalized here.

---

## 9. Review cadence

This document is reviewed annually by the privacy office and legal counsel. Triggers for early review:

- The CPPA issues final regulations on automated decision-making.
- The California Attorney General publishes new enforcement guidance.
- A peer institution's enforcement action reveals a procedural gap.
- A new amendment extends or narrows consumer rights.
- The institution adopts a new automated-decision system requiring re-evaluation of prospective opt-out scope.
