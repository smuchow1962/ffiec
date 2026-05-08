---
title: CFPB / Bureau Articulation Overlay
status: informative
aligned-with:
  - 12 CFR Part 1002 (Regulation B / ECOA)
  - 12 CFR Part 1022 (Regulation V / FCRA)
  - 12 USC §5531, §5536 (Bureau §1031/1036 UDAAP authority)
  - 15 USC §1691 et seq. (Equal Credit Opportunity Act)
  - 15 USC §1681 et seq. (Fair Credit Reporting Act)
  - CFPB CIDs (Civil Investigative Demands)
date: 2026-05-07
version: 1.0.0
---

# CFPB / Bureau Articulation Overlay

> **What this doc is.** A single articulation overlay that maps the FFIEC chain-of-custody v1.0b specification onto the Consumer Financial Protection Bureau's supervisory and enforcement framework. Written so a Bureau Office of Supervision Examinations examiner, a covered entity's compliance officer, a respondent preparing a CID response, and the consumer's representative in a dispute can read this document alongside the spec and confirm what the chain delivers, what it does not, and which artifact discharges which obligation under Reg B (ECOA), Reg V (FCRA), and the Bureau's UDAAP authority.

> **What this doc is NOT.** Not a normative extension to the v1.0a/v1.0b specification. Not a substitute for the institution's compliance program documents. The integrity primitives are framework-neutral cryptographic constructs; the Bureau-supervisory translation is an articulation overlay an institution layers on top, not a change to the underlying specification. Spec §1.2 (epistemic scope) is the discipline this document operates under: the chain proves what was said and that the record was not tampered with after capture; it does not prove substantive correctness, the model's accuracy, or freedom from bias.

> **Round-17 close-out.** This overlay was added per Round-17 CFPB-N2 close-out. The previously-deferred CFPB section in `customer-dispute-procedures.md` §259 ("could add an explicit CFPB section in v1.1") collapses into this regulator-pack document under the close-out directive that pulled v1.1 candidates into v1.0a/v1.0b companion docs.

---

## 1. Scope and reading order

This overlay is consumed in three reading orders.

| Reader | Reading order |
|---|---|
| First-look CFPB examiner | §2 (Reg B headline mapping) → §3 (FCRA dispute-trail framing) → §6 (CID response posture) → §8 (bottom line) |
| Bureau enforcement attorney preparing a CID | §6 (CID response posture) → §7 (consumer-keyed retrieval) → §4 (Bureau-mediated verification) |
| Consumer's representative in an active dispute | §4 (Bureau-mediated verification path) → §5 (witness-mode posture) → §3 (FCRA reinvestigation timing) |

Cross-references are exact. If this overlay sends the reader to `customer-dispute-procedures.md` §"Bureau-mediated verification" for the consumer-side IKM-access shape, that document is the load-bearing source.

---

## 2. Reg B / ECOA — adverse-action notice framing

The headline mapping is spec §10.11 (Adverse-action notice translation, ECOA and state-insurance analog) plus the new §10.11.1 (ECOA adverse-action reasons schema) added at v1.0b. Together they cover Reg B §1002.9(a)(2)(i)'s "principal reason(s)" requirement.

| Reg B obligation | Chain-of-custody evidence |
|---|---|
| §1002.9(a)(1) — 30-day notification window | Spec §10.11 `audit.ecoa.translation.delivery_timestamp` (REQUIRED when `delivery_method` is also recorded per Round-17 CFPB-N1); the within-window check anchors against `entry 1's captured_at` and `entry 2's delivery_timestamp` per `customer-dispute-procedures.md` §109 |
| §1002.9(a)(2)(i) — disclosure of "principal reasons" | Spec §10.11.1 `audit.ecoa.adverse_action.reasons` (REQUIRED on the underlying entry); the listed reasons are integrity-bound under the per-event MAC |
| Match between listed reasons and model output | Spec §10.11.1 `audit.ecoa.adverse_action.feature_attributions` and `audit.ecoa.adverse_action.model_explanation_method`; lets a Bureau examiner answer "do the listed reasons match the model's actual weights?" mechanically rather than circumstantially |
| Translation discipline for limited-English-proficiency consumers | Spec §10.11 translation entry (chained to the underlying entry via `parent_run_id` / `parent_seq`); §10.11.1's reasons are the integrity-bound source the translation derives from |

The four-piece evidence package referenced in `customer-dispute-procedures.md` (chain integration, parent linkage, documented schema, MRM-review link) is the operational shape a Bureau examiner reviewing an ECOA case expects to see.

---

## 3. Reg V / FCRA — dispute-trail framing

§611 reinvestigation timing under FCRA was the load-bearing missing piece in the v1.0a posture (Round-17 CFPB Gap). The Bureau examines whether the institution's dispute-trail discipline meets the 30/45-day clock under §611(a)(1) and §611(a)(3). For a covered entity using AI in a credit-reporting capacity, the chain entries supporting the dispute-trail are emitted under the institution's `audit.fcra.*` event family (institution-side schema; not currently normated in v1.0b — see `customer-dispute-procedures.md` for the institution's documented anchors). The chain's role is integrity-binding the institution's claimed timing anchors; the substantive timing compliance is institution-side.

A Bureau examiner reading the chain for a FCRA reinvestigation case looks for:

- The institution's named dispute-receipt event (typically with attribute `audit.fcra.dispute_received_at`)
- The additional-info event (extends the clock to 45 days under §611(a)(3))
- The furnisher-notification event
- The reinvestigation-completed event
- The consumer-notified event

Each event's `captured_at` is the chain-of-custody anchor for the timing the Bureau examines. The chain proves the events occurred at the recorded times; the substantive question of whether the institution's actions discharged the §611 obligations is institution-side compliance evidence the Bureau evaluates separately.

---

## 4. Bureau-mediated verification (consumer-side full §7 verification)

Per `customer-dispute-procedures.md` §"Bureau-mediated verification" (Round-17 CFPB-P3 close-out), the Bureau (or a court-appointed master) holds the IKM under a standing protective order and runs the full §7 verification on the consumer's behalf when witness-mode `Status: PASS-STRUCTURALLY, key-bound verification skipped` is insufficient for the dispute.

The Bureau-mediated path applies when:
1. The consumer has filed a complaint with the Bureau or a federal court has referred the matter
2. The consumer's expert has run witness mode and reported a structural PASS
3. The consumer's representative believes the per-event MAC verification is material to the dispute resolution
4. The protective-order shape is unavailable or impractical

The Bureau's role is delegated cryptographic verification — the consumer's expert prepares the chain artifacts; the Bureau executes the IKM-bound steps (§7 steps 7, 8, 9) the witness-mode procedure skips. The Bureau reports the per-event MAC result to both parties using the spec §7 normative output format (`Status: PASS` or `Status: FAIL` with named §7 step).

---

## 5. Witness-mode posture for consumer disputes

Spec §7 witness mode produces `Status: PASS-STRUCTURALLY, key-bound verification skipped` when a consumer's expert runs the verifier without IKM access. This is real integrity property — the structural PASS confirms chain-link continuity, Merkle tree consistency, and HSM signature validity over the seal records. What it does NOT confirm is the per-event MAC.

For most consumer disputes, structural PASS is sufficient evidence — it confirms the institution did not silently rewrite the chain after capture. For disputes where the per-event MAC verification is material (e.g., a consumer alleges the chain entries were fabricated rather than rewritten), the Bureau-mediated path under §4 above closes the gap.

---

## 6. CID response posture

A Bureau CID typically reads "produce all adverse-action decisions for consumers in [ZIP, demographic class, product line] during [period]." The institution's response uses the spec's selective-production primitive (`docs/selective-production-and-sampling.md` §"18 USC §2703(d) selective production") with the appropriate cover letter for the CID context. The institution produces:

- The chain artifacts in witness-runnable form (NDJSON + seal records + public-key resolution)
- The institution's customer-correlation-index (CUEC) entries mapping the named consumers to runs
- Verifier output for the produced runs (institution-side, full §7 verification)

The Bureau's verifier (running independently, witness-mode by default) confirms the institution's verification result without trusting the institution's logging infrastructure — that is the asymmetric-evidence-balancing posture the Bureau looks for.

---

## 7. Consumer-keyed retrieval and CUEC integrity

The chain is keyed by `(tenant_id, run_id, seq)`, not by consumer identity. Consumer-keyed retrieval (the typical CID shape) depends on the institution's customer-correlation-index (CUEC). The CUEC's integrity is institution-side; the spec does not normate its shape. A Bureau examiner reading a CID response confirms:

- The CUEC entries themselves are chain entries (or are anchored to the chain via cross-anchor per §10.21-style sidecar discipline)
- The CUEC's coverage matches the chain's coverage map (spec §10.19)
- A CID-class production excluding consumers the institution would prefer not to surface is detectable through the chain-coverage map (§10.19 pre Round-17 update) and the `chain.coverage_map_published` operational event (§10.2 post Round-17 M&A-P3 update)

If the CUEC's integrity is not anchored, a Bureau examiner files a clarification request asking the institution to name the CUEC's integrity mechanism — this is the load-bearing question for any CID-class production.

---

## 8. Bottom line for the Bureau examiner

The chain delivers, under v1.0b, the following integrity-bound evidence the Bureau examines under Reg B (ECOA), Reg V (FCRA), and the Bureau's UDAAP authority:

1. **Adverse-action reasons** (spec §10.11.1) — what the model said, integrity-bound
2. **Translation discipline** (spec §10.11) — what the consumer received in their language, integrity-bound
3. **30-day window evidence** (spec §10.11 with REQUIRED `delivery_timestamp` per CFPB-N1)
4. **Witness-mode verifier path** (spec §7) — consumer-side independent verification without IKM access
5. **Bureau-mediated verification path** (`customer-dispute-procedures.md` §"Bureau-mediated verification") — consumer-side full §7 verification with Bureau-held IKM under protective order
6. **Selective production for CID response** (`selective-production-and-sampling.md`) — partial-disclosure under 18 USC §2703(d)-class cover letters
7. **Redaction discipline** (spec §10.22, Round-17 CFPB-P2) — pre-MAC SDK redaction with normative `audit.redaction.*` family

What the chain does NOT deliver: substantive correctness of the AI's decisions, freedom from bias, or compliance with substantive Reg B / Reg V obligations. Those are institution-side compliance evidence the Bureau evaluates through normal supervisory channels.

---

## 9. Cross-references

- Spec §10.11 — Adverse-action notice translation (ECOA and state-insurance analog)
- Spec §10.11.1 — ECOA adverse-action reasons schema (Round-17 CFPB-P1)
- Spec §10.22 — Redaction discipline (Round-17 CFPB-P2)
- Spec §7 witness-verifier mode
- `customer-dispute-procedures.md` §"Bureau-mediated verification" (Round-17 CFPB-P3)
- `selective-production-and-sampling.md` §"18 USC §2703(d) selective production"
- `regulator-pack/nydfs-part500-overlay.md` (parallel state-supervisory overlay structure)
