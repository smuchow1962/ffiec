---
title: SEC Broker-Dealer / Investment-Adviser Articulation Overlay
status: informative
aligned-with:
  - 17 CFR §240.17a-4 (books-and-records retention; §240.17a-4(f) WORM path and 2022 audit-trail alternative)
  - 17 CFR §240.17a-3 (records to be made by broker-dealers)
  - 17 CFR §240.15l-1 (Regulation Best Interest)
  - 17 CFR Part 248 (Reg S-P) / 17 CFR §248.201 (Reg S-ID)
  - 17 CFR §275.206(4)-1 (Investment Advisers Marketing Rule)
date: 2026-07-02
version: 1.0.0
---

# SEC Broker-Dealer / Investment-Adviser Articulation Overlay

> **What this doc is.** A single articulation overlay mapping the FFIEC chain-of-custody specification (PRD-3 / 0.3.0) onto the SEC's broker-dealer and registered-investment-adviser framework. Written so an SEC examiner, a registered firm's compliance officer, and an enforcement attorney can read it alongside the spec and confirm which artifact discharges which obligation under the §240.17a-4 recordkeeping regime, Reg BI, Reg S-P / S-ID, and the Marketing Rule. It is the foundation the FINRA overlay (`finra-overlay.md`) composes on top of, because FINRA members are SEC-registered broker-dealers and FINRA Rule 4511 incorporates the §240.17a-3 / 17a-4 regime by reference.

> **What this doc is NOT.** Not a normative extension to the specification. Not a substitute for the firm's compliance program. The integrity primitives are framework-neutral cryptographic constructs; the SEC-supervisory translation is an overlay a firm layers on top. Spec §1.2 (epistemic scope) is the discipline: the chain proves what was captured and that the record was not tampered with after capture; it does not prove substantive correctness, model accuracy, or freedom from bias.

---

## 1. Scope and reading order

| Reader | Reading order |
|---|---|
| First-look SEC examiner | §2 (17a-4(f) recordkeeping) → §3 (Reg BI recommendation lineage) → §6 (bottom line) |
| Enforcement attorney | §2 (recordkeeping) → §5 (subpoena / privilege) → §4 (Reg S-P / S-ID) |
| Registered-firm compliance officer | §2 → §3 → §4 → §6 |

---

## 2. §240.17a-4(f) recordkeeping — the settled-rule anchor

The load-bearing mapping is spec **§5.2.3 (SEA 17a-4(f) recordkeeping-rule mapping)**. §240.17a-4(f) is the settled rule text for a tamper-evident electronic-recordkeeping system, and it offers **two co-equal compliance paths**, not one:

- **WORM path** — records preserved on write-once-read-many media (the pre-2022 requirement).
- **Audit-trail alternative** — added by SEC Release No. 34-96034 (effective Jan 3 2023, compliance date May 3 2023): a system maintaining "a complete time-stamped audit trail … to the extent necessary to permit the recreation of the original record if it is altered, over-written, or erased."

The chain is **not** WORM storage — it does not render the medium physically non-rewritable. It is a tamper-*evident* append-only audit trail, which is exactly the shape the audit-trail alternative describes. The audit-trail alternative is therefore the doctrinally-correct and stronger mapping; earlier "WORM-equivalent" framing understated the fit.

| 17a-4(f) audit-trail-alternative element | Chain primitive | Verifier demonstration |
|---|---|---|
| Complete | §4.1 monotonic `seq` + §4.2 daily Merkle seal | §7 chain-linkage rejects gap/reorder; Merkle leaf count binds the day's full set |
| Time-stamped | §4.1 `mac_computed_at_utc` + `captured_at_utc` (MAC-covered) | §7 per-event-MAC step confirms bind-at-capture |
| Audit trail | §7 chain of inferences: per-event MAC → Merkle root → HSM seal signature | Each step rejects a distinct tamper class |
| Recreation if altered or deleted | §10.3 append-only + Merkle seal (deletion) + per-event MAC (alteration) | Tamper surfaces as a named §7 `FAIL`; un-tampered bytes remain the recoverable original |

A verifier `Status: PASS` **is** the audit-trail-alternative conformance demonstration for the records the ledger covers. The firm's CC8.1 names the retention duration (17a-4 sets durations; the chain is duration-agnostic).

**Scope boundary — [OPEN QUESTION].** 17a-4(f) governs *how* records are preserved once they are records. Whether a given AI prompt or output is itself a "record … relating to its business as such" under 17a-4(b)(4) is a separate, currently-unsettled question. It is **not** resolved by the SEC or FINRA as of mid-2026 (see FINRA RN 25-07, April 2025, which requested comment on exactly this). This overlay does not answer it; the chain binds whatever the firm captures, and the firm makes the risk-based recordkeeping judgment.

---

## 3. Reg BI — recommendation lineage

An AI-generated recommendation to a retail customer is still a recommendation under Reg BI (17 CFR §240.15l-1); the care and compliance obligations sit with the firm. The chain expresses a recommendation-with-lineage through composition — no Reg-BI-specific attribute family is required:

| Reg BI element | Chain-of-custody evidence |
|---|---|
| The recommendation was made and captured | §10.47 generation four-tuple (system/user prompt, retrieval root, output — all hash-bound) |
| The reasoning basis | §14.7 `audit.reasoning.substrate_kind` (the reasoning architecture class) + §10.49 retrieval-source integrity |
| Human review / override | §10.50 output-grounding review (signed reviewer, role, outcome) |
| Downstream action taken on the recommendation | §14.8 `audit.downstream_action.*` (action kind, system-of-record, applied-at) |
| Prior-event lineage | §4.4 `parent_run_id` / `parent_seq` cross-binding (the same lineage shape §10.11.1 uses for the ECOA prior-offer link) |

The chain proves the recommendation was captured with cryptographic integrity and links to its reasoning, review, and downstream action; it does not prove the recommendation was in the customer's best interest — that is the firm's substantive Reg BI compliance evidence.

---

## 4. Reg S-P / Reg S-ID — privacy and identity-theft

The chain never binds raw PII: customer identifiers enter as hashes (e.g., §10.47 `user_prompt_sha256`, §14.6 `authenticated_user_id_hash`), and the token-vault discipline (`token-vault-architecture.md`) holds the token→PII mapping outside the chain. This posture supports Reg S-P's safeguarding expectation and Reg S-ID's red-flags recordkeeping without placing consumer financial information into the integrity ledger. Deletion requests are bound as separate chain entries while the vault destroys the mapping (see §10.74 crypto-erasure composition).

---

## 5. Subpoena and privilege interaction

SEC subpoena and privilege interaction runs through spec **§10.70 (privileged-investigation overlay)**: a verifier invoked with a cleared-role claim returns full content; without it, redacted-with-existence-attestation. Selective production for a subpoena uses `selective-production-and-sampling.md`. The verifier output is independently reproducible by the SEC without trusting the firm's logging infrastructure — the asymmetric-evidence posture an enforcement attorney looks for.

---

## 6. Bottom line for the SEC examiner

1. **17a-4(f) recordkeeping** (spec §5.2.3) — the chain satisfies the settled audit-trail alternative, element-by-element, demonstrable by running the verifier.
2. **Reg BI recommendation lineage** (spec §10.47 / §14.7 / §10.50 / §14.8) — recommendation, reasoning, review, and downstream action are integrity-bound and linked.
3. **Reg S-P / S-ID** — PII never enters the chain raw; hashes and the token vault carry identity.
4. **Subpoena / privilege** (spec §10.70) — role-gated production with independent verifier reproduction.

What the chain does NOT deliver: substantive best-interest compliance, model accuracy, or an answer to whether a given AI output is a "business as such" record. Those are firm-side compliance and unsettled-law questions the SEC evaluates through normal channels.

---

## 7. Cross-references

- Spec §5.2.3 — SEA 17a-4(f) recordkeeping-rule mapping
- Spec §10.47 / §10.48 — generation four-tuple + model provenance
- Spec §14.7 — reasoning substrate_kind; §10.50 — output-grounding review; §14.8 — downstream action
- Spec §10.70 — privileged-investigation overlay; §10.22 — redaction discipline
- Spec §13 — SEC examiner and FINRA-examiner (SRO) stakeholder entries
- `finra-overlay.md` — the SRO overlay that composes on top of this one
- `token-vault-architecture.md`; `selective-production-and-sampling.md`
