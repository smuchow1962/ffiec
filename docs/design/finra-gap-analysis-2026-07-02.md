# FINRA gap analysis — spec + verifier (2026-07-02)

**Author:** Jared (Go/Rust systems, TesseraSeal verifier)
**Spec under analysis:** `E:\dev\ffiec-public\spec\chain-of-custody-DRAFT-0.3.0.md` (PRD-3, local; public repo pins 0.2.0)
**Research input:** `E:\dev\Herald.Documentation\private\research\finra-ai-requirements-2026-07-02.md` (verified 2026-07-02)
**Question (Steve, via team-lead):** (1) Does the product need anything IN THE SPEC for FINRA? (2) What code change does TesseraSeal need?

## Accuracy discipline applied throughout

Two statuses are kept distinct and never conflated:

- **[SETTLED RULE]** — text in force today: SEA 17a-3/17a-4 (incl. the 2022 17a-4(f) audit-trail alternative), FINRA 4511, Rule 3110, Rule 2210, Reg BI.
- **[EXAM SIGNAL]** — the 2026 FINRA Annual Regulatory Oversight Report's GenAI expectations (prompt/output logs, model-version provenance, "track and log AI agent actions and decisions", HITL evidence). Stated exam expectations / observed practices, **not rule text.**

There is **no FINRA AI rule** and none is formally proposed (mid-2026). The RN 25-07 question — whether AI outputs are 17a-4(b)(4) "business as such" records — is **open**; the spec must not encode an answer.

---

## 1. Verdict table — requirement → spec verdict → action

| # | FINRA requirement | Status | Spec today | Verdict | Action |
|---|---|---|---|---|---|
| 1 | **SEA 17a-4(f) audit-trail alternative** — "complete time-stamped audit trail permitting recreation of an original record if altered or deleted" | SETTLED | §240.17a-4(f) named only as **"WORM-equivalence"** (§11 refs; §5.239/§5.242 SEC stakeholder; §10.3 anchor) | **Covered-but-mis-mapped** — the chain is tamper-*evident*, not WORM storage; the audit-trail alternative is the doctrinally-correct AND stronger hook, and it is settled rule text | **Add explicit mapping clause** mapping the four rule elements to primitives; **correct** the WORM-only framing to name both compliance paths. Highest-value, lowest-risk delta. |
| 2 | **FINRA Rule 4511** — books/records per 17a-3/17a-4 | SETTLED | Not named anywhere; §11 names the SEC rule (17a-4) but not the FINRA rule that incorporates it | **Covered-but-unnamed** | Name 4511 as an anchor beside 17a-4 in §11 + SEC stakeholder section. One-line addition. |
| 3 | **Rule 3110 supervision** — supervise AI activity; firm owns outputs "human or AI" | SETTLED | Heavy supervisory vocab: SR 11-7, CC8.1, §10.50 HITL, §10.80 three-lines-of-defense, §10.81 anomaly taxonomy | **Covered** (unnamed) | Regulator-pack cross-walk note only. No spec-mechanism change. |
| 4 | **Rule 2210 principal pre-approval** — registered principal approves *before* retail comms are used | SETTLED | §10.50 review (`pending_review→reviewed`, role enum, signed `signed_at_utc`) + §14.8 `downstream_action` (`action_kind`, `applied_at_utc`) + parent-linkage | **Narrow gap** — the *substrate* to express "approval precedes send" exists (two timestamped, parent-linked events); what is absent is (a) a named **ordering discipline** and (b) two semantic tags: a `registered_principal` role value and a retail-vs-institutional audience classifier (the >25-retail-in-30-days trigger) | **Optional** §10.x soft-enforcement discipline + two enum extensions. The one place new normative text (and, if greenlit, bounded code) has real value. Do **not** build speculatively. |
| 5 | **Reg BI** — an AI recommendation to a retail customer is still a recommendation; care obligation on the firm | SETTLED | §10.47 generation four-tuple + §14.7 `substrate_kind` + §10.49/§10.50 grounding/review; §10.11.1 carries a *lineage* pattern (`prior_offer_run_id/seq` parent-linkage) but ECOA-shaped | **Covered-by-pattern, unnamed for securities** | Regulator-pack mapping note showing the generation+substrate+review+downstream_action families express a Reg-BI recommendation-with-lineage. No code. |
| 6 | **Model-version provenance** — "which model version was used and when" | EXAM SIGNAL | §10.47 `model_id` (REQUIRED) + `inference_at_utc` (REQUIRED); §10.48 `model_version` + `model_weight_hash`; §10.33 model-update; §10.21 handover | **Covered** | Regulator-pack pointer at the existing fields. No change. |
| 7 | **Prompt / output logs** — "storing prompt and output logs" | EXAM SIGNAL | §10.47 binds `system_prompt_sha256` / `user_prompt_sha256` / `output_sha256`; content retained under §10.69 (PII never enters the chain raw) | **Covered** | Regulator-pack note making the integrity-vs-content split explicit (chain binds hashes; the 17a-4-governed store holds the logs). No change. |
| 8 | **"Track and log AI agent actions and decisions"** | EXAM SIGNAL | §14.6 `audit.actor.*` (+ delegation chain for on-behalf-of/agent), §14.8 `audit.downstream_action.*`, §14.7 `substrate_kind` — all SHIPPED C#/Py/Go + vectors 050-084 | **Covered** | None. |
| 9 | **Human-in-the-loop evidence** | EXAM SIGNAL | §10.50 signed-review (role enum, signed_at_utc, signature) + §14.6 actor | **Covered** | None. |
| 10 | **SRO-as-examiner** (FINRA is an SRO, not a government agency) | FRAMING | Scope §1 line 143 already accommodates FINRA "under its own supervisory mandate without spec-amendment"; verifier trust model (§10.76, design §07) is "principal holds the public key" — agnostic to gov-vs-SRO | **No structural assumption blocks FINRA.** But: no FINRA stakeholder-nav entry, and **`docs/regulator-pack/sec-overlay.md` (referenced by §5.239) does not exist** | Author the missing SEC overlay + a FINRA overlay; add a §13 stakeholder entry. Doc work, not spec-mechanism. |
| 11 | **AI outputs as "business as such" records** | OPEN (RN 25-07) | Spec is technology-neutral; does not opine | **Correctly silent** | Note only. Do not encode an answer; the SEC/FINRA have not resolved it. |

---

## 2. Code-change recommendation

**Headline: zero required verifier code for the settled-rule and exam-signal requirements. One optional, bounded predicate — gated on a spec-delta decision — for Rule 2210 ordering.**

Why the settled requirements need no code:

- **Model-version, prompt/output-log hashing, agent-action logging, HITL** all land on **already-shipped attribute families** (§10.47/§10.48 generation, §14.6/§14.7/§14.8 actor/substrate/downstream_action) with C#+Python+Go reference implementations and conformance vectors 050-084. FINRA *reuses* these; it adds no field. The Go verifier already dispatches these as additive `additional_verifications` checks (`ValidateActorAttributes`, `ValidateReasoningAttributes`, `ValidateDownstreamActionAttributes` in `verifier/internal/verify/predicates.go`).
- **17a-4(f) audit-trail alternative** is already *proved mechanically* by the existing §7 walk: **complete** = `seq` monotonic, no gaps; **time-stamped** = `mac_computed_at_utc` / `captured_at_utc`; **audit trail** = the §7 chain of inferences (per-event MAC → Merkle root → HSM signature); **recreation if altered or deleted** = the Merkle seal catches deletion, the per-event MAC catches alteration. A verifier `Status: PASS` **is** the audit-trail-alternative conformance demonstration. No new code; the value is in *naming* this in the doc layer (and optionally in a regulator-pack), not in the binary.

The one optional piece — **Rule 2210 approval-precedes-send** (requirement #4):

- **Only if** Steve wants principal pre-approval as a first-class conformance check. It requires a spec-delta **first** (Steve-gated), because it adds a closed-enum `additional_verifications` marker (a doc-version revision per §10.12's closed-enumeration rule) and/or a new PASS-anomaly line.
- Shape (mirrors §14.8 `downstream_action` closely): a soft-enforcement anomaly `principal approval did not precede communication send at seq N` under `Status: PASS`, plus an optional marker `communication_principal_preapproval_verified`; a `registered_principal` value added to the §10.50 `audit.review.role` enum; an audience classifier (`retail` | `institutional` | `correspondence`) so the >25-retail-in-30-days pre-approval trigger is expressible.
- **Effort if greenlit:** ~1–1.5 engineer-days across the three reference implementations (Go predicate + C#/Py emitters) + a positive and a negative conformance vector + the byte-equivalence check. It is *soft-enforcement* (completeness, not chain-integrity), so it does **not** gate `Status: PASS/FAIL` and does not risk regressing existing vectors.
- **My recommendation:** propose it in the delta list below; **do not build it until greenlit.** The existing §10.50 + §14.8 + parent-linkage already let an institution *emit* the approval-before-send sequence; what's missing is only the verifier *asserting* the ordering. That is worth adding for FINRA members specifically, but it is a discrete decision, not a prerequisite for the settled-rule story.

An "attestation output naming 17a-4(f) conformance" (a verdict-object field) was considered and **rejected** for the verifier: the verifier's job is chain-integrity, not regulatory legal-conclusion labeling. The 17a-4(f) element-to-primitive mapping belongs in the spec's legal-mapping section and the regulator pack, where a compliance officer and an examiner read it — not baked into the binary's output where it would read as the tool asserting a legal conclusion.

---

## 3. Proposed spec-delta list (for a spec-editing pass — Steve reviews)

Analysis + proposal only; no spec text edited here. Section numbers are against 0.3.0.

**Settled-rule deltas (recommended — all additive/corrective, no wire-format or code impact):**

1. **New §5.2.3 "SEA 17a-4(f) recordkeeping-rule mapping (normative when applicable)."** Parallel to the §5.2.1/§5.2.2 FRE mappings. Maps the four 17a-4(f) audit-trail-alternative elements (complete / time-stamped / audit trail / permits recreation if altered or deleted) to spec primitives (§4.1 per-event MAC, §4.2 Merkle seal, §4.3 HSM signature, §7 walk, §10.3 append-only). States both compliance paths (WORM **and** audit-trail alternative) and that the chain naturally satisfies the audit-trail alternative.
2. **Correction at §11 references + §5.239/§5.242 (SEC stakeholder).** Change "§240.17a-4(f) WORM-equivalent retention discipline" to name **both** the WORM path and the **2022 audit-trail alternative** (effective Jan 2023 / compliance May 2023); cross-reference the new §5.2.3.
3. **§11 references — add FINRA Rule 4511** as an anchor beside 17 CFR §240.17a-4 (4511 requires FINRA-member books/records per 17a-3/17a-4).
4. **§13 stakeholder navigation — new "FINRA examiner (SRO)" entry.** Mirrors the SEC-examiner entry; notes FINRA members are SEC-registered broker-dealers, 4511 incorporates 17a-4, and the verifier trust model is SRO-agnostic (principal holds the public key). Reinforce the §1 line-143 "own supervisory mandate" framing.

**Optional delta (Rule 2210 — build only if greenlit):**

5. **New §10.x "Communication principal-preapproval ordering (normative when applicable)."** Reuses §10.50 review + §14.8 downstream_action + parent-linkage. Adds: `registered_principal` to the §10.50 `audit.review.role` enum; an audience classifier (`retail`/`institutional`/`correspondence`); a soft-enforcement anomaly `principal approval did not precede communication send at seq N` under `Status: PASS`; and (if the marker route is chosen) a §10.12 closed-enum addition `communication_principal_preapproval_verified` (doc-version revision).

**Doc-layer (regulator pack — not spec text, but prerequisite for the FINRA story):**

6. **Author `docs/regulator-pack/sec-overlay.md`** — referenced by §5.239 but **missing today.** FINRA composes on top of it (members are SEC-registered broker-dealers).
7. **Author `docs/regulator-pack/finra-overlay.md`** — the SRO cross-walk: 4511/17a-4(f), 3110, 2210, Reg BI → chain families; and an explicit, clearly-labeled **[EXAM SIGNAL]** section for the 2026-report GenAI expectations pointing at the existing §10.47/§10.48/§14.6/§14.7/§14.8 fields. Keep the RN 25-07 records question flagged as open; do not resolve it.

---

## Bottom line

- **Spec:** four small, recommended deltas (one is a genuine accuracy *correction* — 17a-4(f) is settled and the current WORM-only framing understates the fit) + one optional Rule-2210 discipline. FINRA needs **naming and mapping**, not new mechanism.
- **Code:** nothing required. The FINRA requirements land on already-shipped, already-vectored families. The only code on the table is the optional Rule-2210 ordering predicate (~1–1.5 days, three implementations + vectors, soft-enforcement so no regression risk) — and only after a spec-delta decision.
- **Prereq surfaced:** the SEC overlay the spec already points to does not exist; author it before/with the FINRA overlay.
