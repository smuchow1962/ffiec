# 🧾 Auditor Story 09 — Sun-Won Cosmetics Group (Korea + Taiwan)

> **Engagement.** Sun-Won Holdings (KOSPI: 003410) and Sun-Won Taiwan Co. Ltd. (Taipei Exchange: 5891) — a coordinated annual review covering PIPA Section 28 cross-border transfers, FSS supervisory review of the BNPL consumer-finance arm, Taiwan FSC review of the listed subsidiary, PDPA Article 8 explicit-consent oversight, and a CPRA / GDPR sweep against the e-commerce platform.
>
> **Subject system.** TesseraSeal v1.0b, in production for sixteen months across four AI use cases — customer recommendation, inventory forecasting, multilingual chatbot, BNPL credit-scoring. Eight tenants, plus a ninth cross-jurisdiction tenant for inventory.
>
> **Audit team.** Eight people, split: four in Seoul (Sangam-dong HQ), four in Taipei (Xinyi District subsidiary office). One video bridge. One chain. Three regulators. Two HSMs.
>
> **Date.** April 9, 2026 — anchored to Seoul time (KST, UTC+9). Taipei time runs one hour behind (CST, UTC+8) and is noted explicitly when the scene crosses the strait.

---

## Context

Sun-Won is a Korean cosmetics group with about $1.8B in annual revenue. Three hundred-plus stores in Korea, eighty-plus in Taiwan, and one e-commerce platform that runs both fronts on a single codebase with localized storefronts, payment rails, and inventory pools. There is a small consumer-finance arm offering buy-now-pay-later, supervised by the Korean FSS (금융감독원). The company is a K-beauty house. The marketing model is celebrity-driven, which is relevant to today's review.

Twelve months before this engagement, a celebrity-endorsement controversy alleged that Sun-Won's AI personalization had used the celebrity's biometric features without explicit consent — facial-shape descriptors lifted from public modeling photos. The story was front-page in Korea for two weeks. It cost a Vice President her job. It moved the rollout of the audit chain forward by nine months. The chain was deployed partly to demonstrate, with cryptographic evidence, exactly which features the recommendation models used and which were excluded.

That history sits behind today's engagement. The team is here for an annual review. They are also here because two regulators in two countries — and a third regulator on the financial side — want to see the same chain entries through three different lenses.

The Seoul team works with **Park Hye-jin**, Sun-Won's Chief Compliance Officer. Park spent eleven years at the FSS before joining Sun-Won; she knows the supervisory templates by memory. She is fluent in Korean, English, and conversational Japanese. She has been preparing for this audit for six weeks.

The Taipei team works with **Lin Jia-hua**, Director of Legal & Compliance at Sun-Won Taiwan. Lin is a KPMG Taiwan alumna, native Mandarin, fluent English, legal-trained rather than compliance-trained — she asks clarifying questions about what FSC will read into a finding before agreeing to its phrasing.

Dawn has done eight of these now. She knows the cadence. The cross-border boundary is the part she has not seen tested under three regulators at once. The good news is the spec moved underneath her between the last engagement and this one — the v1.0b amendment locked the routing-classifier event type and the cross-border-transfer attribute family that Sun-Won's posture asked for almost word-for-word. Two of the three Partials she would have written eight months ago are now framed as "this is the work item the spec already named — here is the deadline."

---

## Audit Team

### Seoul (Sun-Won HQ — Sangam-dong)

- **Dawn** — Lead Auditor. Anchors the engagement from Seoul. Final sign-off on all four regulator-partitioned findings.
- **Raj** — Database specialist. BNPL credit-scoring chain walk and the per-applicant entry structure.
- **Diana** — IAM and access control. PASS-IT integration, per-tenant scoping, key custody on the Seoul HSM.
- **Tom** — Internal-audit liaison specialist. Bridge to Park-CCO's binder of prior FSS supervisory letters.

### Taipei (Sun-Won Taiwan — Xinyi District)

- **Elena** — CRM systems. Customer-record provenance, the recommendation-engine input side.
- **Mike** — Application / API layer. Recommendation engine and chatbot wire-side review.
- **Luis** — DevOps, logs, and pipelines. Daily seal landing, log retention, the language-detector microservice.
- **Chen** — Data engineering / ETL. Inventory-forecasting cross-jurisdiction data flow.

The two halves work in parallel across the morning. They join on a video bridge for the noon lunch and again at the 4:30 PM debrief. Timestamps in this diary are anchored to **Seoul time**. Taipei-local time is given in parentheses where it matters.

---

## 🌅 8:30 AM — Seoul, Sun-Won HQ Lobby

The Sangam-dong tower has a glass atrium that catches the morning light off the Han River. Dawn, Raj, Diana, and Tom badge in at the security desk. The receptionist switches to English the moment Dawn says her name and hands them four visitor lanyards in Sun-Won's house pink.

Park Hye-jin meets them at the elevator bank. She is in a charcoal suit, no jewelry, carrying a leather portfolio that Dawn recognizes as the FSS examiner-issue from about 2011.

"Welcome back," Park says. "Lin-Director's team should already be with your Taipei four. We have the bridge open in the executive conference room on twenty-eight."

"Thank you for the early start," Dawn says. "How is the binder?"

"Six weeks of preparation. The FSS supervisory letters from the last cycle are tabbed. The PIPC quarterly attestations are tabbed. The Taiwan FSC corresponding letters Lin's team holds — those are duplicated on her side. I have one set here as well."

*She is being thorough on purpose,* Dawn thinks. *She knows the difference between a routine year and a coordinated three-regulator review. She is not going to make us ask twice.*

They ride up. Park hands Dawn a printed agenda in Korean and English, side-by-side columns. The eight tenants are listed with their tenant_ids and their primary regulator mapping. The ninth — `sunwon-cross-inventory` — is at the bottom, in italics. Dawn notes that each tenant_id matches the §3 character class `^[A-Za-z0-9_.\-]{1,255}$` cleanly — Sun-Won's IAM provisioning enforces the class at registration so the SDK-side and verifier-side tenant_id checks (§3 enforcement and §7 step 3a) never have anything to reject. Sun-Won had three legacy CRM identifiers in the original tenant set; Park-CCO's team used §3.1 Pattern 1 (opaque hash-of-legacy with `tnt_` prefix) for two of them and §3.1 Pattern 2 (controlled aliasing) for the third. The legacy-mapping registry is institution-internal per §3.1's operational requirements; SOC 2 testing confirms append-only enforcement on the registry.

"That's the one we should talk about," Park says. "I marked it because you will see it before lunch."

"Thank you," Dawn says. "We will."

Park turns to Raj. "We pre-staged the BNPL chain on a read-only mirror. Your verifier credentials are in your packet — PASS-IT-bound, scoped to read-only, expire at six tonight. Diana, your IAM packet is the same shape, scoped to the IdP audit role. Tom, the compliance binder is on the second cart in the room. Photocopying is fine. Photography is not."

"Understood," Tom says.

"And the room has tea, coffee, and a coffee machine that none of us know how to operate," Park says. "We will figure it out together."

---

In the elevator, Dawn does her opening monologue for the Seoul four — quietly, the way she always does, while they ride up.

> "Two countries, four use cases, three regulators, eight tenants, one chain. Park-CCO has been preparing for this for six weeks. We have done eight of these now — Northbridge, Mercator, Stelvio, Atrio, Helmstad, Pacific Crescent, Olmstead, and we just came from a 23-tenant Israeli AI shop where the test was nation-state segregation. Today's test is the cross-border data-flow basis. The chain is the same chain. The questions are not."

Northbridge sat eight engagements back. One §10.16 non-conformance, the chain otherwise byte-for-byte clean — Dawn had not seen another engagement come within reach of it since. Sun-Won was the first East-Asian engagement of the cycle, PIPA + PDPA cross-border was new, the chain primitive was familiar; the question on the drive in was whether the gap between Northbridge and the rest was the product or the institution.

By the time the team rolled into Sun-Won, TesseraSeal had been audited at multiple US institutions plus a multi-tenant SaaS vendor in Tel Aviv. Korea + Taiwan was the first East-Asian engagement — PIPA + PDPA cross-border composition was new — but the chain primitive was familiar.

Raj nods. Diana is already pulling up the PASS-IT documentation on her tablet. Tom is making sure his recorder is on.

"The chain holds within a jurisdiction," Dawn says. "The cross-border evidence has to hold to two different regulators reading the same chain. Spec §1.4 calls that compositional security — three independent layers, per-event MAC plus daily Merkle seal plus HSM-rooted signature. None of those layers cares about jurisdiction. The institution's posture has to do the jurisdictional work, and §10.18 says that posture has to be testable."

"It never is," Tom says, half a second before Dawn can.

She lets him have it.

Diana, half to herself: "And §1.3 — the security definitions. Per-event MAC under HMAC-SHA-256 is EUF-CMA secure under FIPS 198-1; daily Merkle seal under RFC 6962 provides second-preimage resistance over SHA-256 (FIPS 180-4); HSM signature under Ed25519 (FIPS 186-5) is EUF-CMA. Three independent layers, each FIPS-current, each with a named formal security property. That is what the §1.4 compositional argument composes over."

---

## 🌅 8:30 AM — Taipei (7:30 AM CST), Sun-Won Taiwan, Xinyi District

Elena, Mike, Luis, and Chen are already in the subsidiary's compliance suite. Lin Jia-hua is pouring tea — proper Alishan oolong, not the office-machine kind — into four small cups arranged on a tray.

"Welcome to Taipei," Lin says in English. "I hope you slept. The hotel is two blocks. Please let me know if anything is wrong."

"Everything is good," Mike says. "We are ready."

Lin gestures to the wall monitor, which is showing the video bridge status. The Seoul side has not yet joined.

"We will go in parallel until twelve," Lin says. "Park-CCO and I have aligned the agenda. I will be on the bridge with you the whole time. If you want to ask a question that crosses to Seoul, I can pass it. If a question is FSC-shaped — Taiwan-side — I will answer first."

"Understood," Elena says. *KPMG-trained,* she thinks. *She is going to want to see the language before it goes in any letter.*

Lin sits down. "FSC will read this differently than FSS does. FSC reads cross-border data flow as a listing-disclosure question. PDPC reads it as an Article 8 consent question. The same chain entry. Two letters."

"That's why we are here," Mike says.

Lin pours the second round of tea. "One more piece of context. Sun-Won Taiwan is not a passive subsidiary in this engagement. The Taipei Exchange listing puts independent disclosure obligations on us. FSC will read the cross-border findings as listing material — not because the inventory data itself is material, but because the *evidence framework* for handling it is material to a reasonable investor. That is the FSC reading. PDPC reads the same finding under Article 8. We have to be precise about which sentence belongs in which letter."

Elena writes that down. *Lin is going to want to see every sentence before it goes anywhere.* "We will phrase findings neutrally. You and Dawn take the partitioning."

"Thank you," Lin says.

Mike has his laptop up showing the Taiwan tenant header. "Sun-Won runs FFIEC-conformance posture per §4.1.2 — `ffiec.chain.posture = ffiec` on every emitted Resource. Vendor-flag mode is not in play here. That matters for the FSC reading because FSC asked in the 2025 letter whether the chain's HKDF constants are jurisdictional, and the answer is no — the HKDF salt and info-base are FFIEC-conformance constants per §4.1, byte-identical across tenants. The jurisdictional binding is the per-tenant `info` parameter under §4.1, which encodes `tenant_id` into the HKDF input. Two tenants in the same HSM derive different session keys; two regulators reading the chain see the same construction.

"And the OTLP transport identification per §4.4.3 is set on every emitted Resource — `ffiec.chain.spec`, `service.name`, `service.version`, `ffiec.chain.posture`, `ffiec.chain.format_version`. The receiver dispatches once per OTLP request rather than per-entry. The HTTP transport carries `X-FFIEC-Chain-Spec: v1.0` and `X-FFIEC-Chain-Posture: ffiec` per §4.4.3's recommended header section, mirroring the Resource attributes. Severity per §4.4.4 is collector-passed-through unfiltered — the receiver stamps a `SeverityNumber` in the §4.4.4 range `9..20` so downstream filters do not silently drop chain traffic."

Lin makes a note. "FSC will accept that framing. It is the framing they wanted in the previous letter."

---

## 🧩 9:15 AM — Seoul, Conference Room 28-A

Raj opens the BNPL credit-scoring chain. He has been thinking about it on the plane — the BNPL tenant is the highest-stakes one in the portfolio because the FSS examiner standard for AI-driven credit decisions is the strictest standard at the table today.

He pulls up the tenant config first.

```
tenant_id:        sunwon-kr-bnpl
hsm:              seoul-sangam-hsm-01 (KISA-certified)
ikm:              kr-bnpl-ikm-2026-q2
seal_cadence:     daily, 23:59:59 KST
spec_version:     1.0b
sign_payload_ver: v1.0b
attribute_set:    [applicant_id_hash, model_id, model_version,
                   features_hash, score, decision, reviewer_override,
                   audit.underwriting.features.*, audit.deployment.*]
```

Raj points at the `sign_payload_version` line. "That's the v1.0b 12-line form. Round-17 closed two NIST partials — the day's distinct `key_version` set and the day's distinct `kms_handle_uri` set both bind under the seal signature now (§4.3 v1.0b amendment). v1.0a chains in storage continue to verify under the original 10-line form indefinitely. New seals from the SDK upgrade in February are 12-line."

Park watches over Raj's shoulder. "The reviewer_override field — that's the one the FSS examiner asked about in the last letter."

"I see," Raj says. "What did the letter ask?"

"It asked whether the override is captured deterministically — meaning, can the chain show that a human reviewer overrode the model, distinct from the model's own conditional output."

Raj nods. He pulls a sample entry from April 4 — five days ago, well within the daily-seal window per §4.2.1 cadence.

```
entry_id:               kr-bnpl-2026-04-04-00041823
applicant_id_hash:      sha256:e7c2...9f1a
model_id:               sunwon-bnpl-v3
model_version:          3.4.2
features_hash:          sha256:b81f...0c7e
score:                  712
decision:               conditional
reviewer_override:      none
gen_ai.request.model:   sunwon-bnpl-v3
gen_ai.response.model:  sunwon-bnpl-v3
audit.deployment.intent:        production
audit.deployment.policy_version: kr-bnpl-deploy-2026-q2
audit.underwriting.features.feature_vector_hash:    sha256:b81f...0c7e
audit.underwriting.features.feature_store_version:  fs-kr-bnpl-2026-q2
audit.underwriting.features.feature_categories:     [income_band, debt_ratio_band, employment_band, ...]
sealed_at:              2026-04-04T23:59:59+09:00
seal_signature:         ed25519:7a2f...4b8c
```

"Conditional with no override," Raj says. "So the model itself returned conditional. Both `gen_ai.request.model` and `gen_ai.response.model` are present per the §4.4 MUST requirement — without either, the SDK would have refused to emit per the §4.4 SDK-side enforcement rule. And the Round-17 NAIC-P1 close-out lands the underwriting features family — `audit.underwriting.features.*` is REQUIRED on any chain entry whose `chain_kind = 'audit'` and whose event corresponds to a model-driven underwriting/triage/pricing decision per §4.4.5. BNPL credit scoring qualifies."

"Correct," Park says.

Raj scrolls to the next sample, picked by his fuzzer earlier this morning.

```
entry_id:               kr-bnpl-2026-04-04-00041901
applicant_id_hash:      sha256:1d4a...c0f7
model_id:               sunwon-bnpl-v3
model_version:          3.4.2
features_hash:          sha256:9c2e...44b1
score:                  689
decision:               approved
reviewer_override:      override_to_approved by reviewer_id:kr-rv-014
audit.deployment.intent:        production
audit.deployment.policy_version: kr-bnpl-deploy-2026-q2
sealed_at:              2026-04-04T23:59:59+09:00
seal_signature:         ed25519:7a2f...4b8c
```

"There it is," Raj says. "Score 689 would have been declined under the auto-cutoff at 700. Reviewer 014 overrode to approved. The chain shows the model's score and the reviewer's decision separately. FSS-grade. And this is the discipline §10.11.1 wants on adverse-action shape — `audit.ecoa.adverse_action.reasons` plus the reviewer's structured-reasons code, integrity-bound under the per-event MAC. Sun-Won's BNPL emits the analog for FSS even though §10.11.1 was authored against the US ECOA model."

Park exhales just slightly. *That was the question she was holding.*

"Run the verifier," Dawn says from across the table.

Raj types.

```
$ herald-verify --tenant=sunwon-kr-bnpl \
                --service=credit-score \
                --date=2026-04-09 \
                --strict

[verify] tenant=sunwon-kr-bnpl
[verify] hsm=seoul-sangam-hsm-01
[verify] entries scanned: 18,442
[verify] seal coverage: 365/365 days
[verify] strict mode: ON
[verify] sign_payload_version: v1.0b (12-line form, §4.3 + §7 step 11)
[verify] reviewer_override field: present in 1.7% of entries (expected range)
[verify] JCS self-test: PASS (per §7 pre-flight, NIST-P2)
Status: PASS
```

> ✅ **Confirmation #1 — BNPL credit-scoring (FSS-grade) under §4.4 + §4.4.5.**
> Per-applicant entries capture applicant_id_hash, model_id, model_version, features_hash, score, decision, reviewer_override, plus `audit.underwriting.features.*` per §4.4.5 (Round-17 NAIC-P1) and `audit.deployment.intent`/`policy_version` per §4.4.2. The override is captured as a distinct field with reviewer_id, separable from the model's own conditional output per the §10.11.1 adverse-action reasons discipline applied by analogy. Verifier strict-mode PASS over 365 days, 18,442 entries; output format matches §7 normative `Status: PASS`. JCS self-test (§7 pre-flight, Round-17 NIST-P2) executed before chain walk. Per-event MAC verification under §7 step 9 with `expected_prev_hash` (not `entry.prev_hash`) per the named footgun-avoidance rule.

"That's the one the FSS letter asked about," Park says. "We can put that in front of an examiner without commentary."

"Yes," Dawn says. "That entry stands on its own."

Raj has one more probe. He runs the verifier in `--diff-features` mode against an entry from January, when the model_version was 3.3.x — checking that an upgrade between versions does not silently drop a feature.

```
$ herald-verify --tenant=sunwon-kr-bnpl \
                --service=credit-score \
                --diff-features \
                --from=2026-01-15 --to=2026-04-04

[verify] feature-set delta:
[verify]   2026-01-15 model_version=3.3.7  features=[..., income_proxy_v2, ...]
[verify]   2026-04-04 model_version=3.4.2  features=[..., income_proxy_v3, ...]
[verify] delta documented in feature-change manifest: kr-bnpl-fcm-2026-q1
[verify] PASS
```

"Feature-change manifest is wired," Raj says. "Every model_version bump that changes the feature set lands a manifest entry. Verifier reconciles it. The features_hash on each chain entry covers the canonical-bytes of the feature vector under §4.1's per-event MAC, so a quietly-removed feature shows up as a features_hash divergence from the manifest — not a chain-integrity finding, but a control-completeness signal that triggers MRM review per §4.4.2's MRM dispositions."

Park: "That was the 2026-Q1 letter. Already answered. I just like seeing it answered twice."

---

## 🧩 9:15 AM — Taipei (8:15 AM CST), Sun-Won Taiwan Compliance Suite

Mike opens the recommendation-engine chain for the Taiwan tenant.

```
tenant_id:        sunwon-tw-rec
hsm:              taipei-chunghwa-hsm-02 (CNS 27001 aligned)
ikm:              tw-rec-ikm-2026-q2
seal_cadence:     daily, 23:59:59 CST
spec_version:     1.0b
sign_payload_ver: v1.0b
attribute_set:    [user_id_hash, model_id, model_version,
                   features_hash, recommended_skus, served_at]
```

"Different HSM," Lin says. "Different IKM. Different cadence anchor — Taipei time, not Seoul time."

"That's the right shape for what §10.15 calls Pattern B," Mike says. "Per-jurisdiction tenant means per-jurisdiction key material. Cross-region correlation is institution-side; cross-region run continuation is not attempted. Each regional tenant has its own IKM, its own seals, its own verifier runs. PIPA Section 28 and PDPA Article 8 both want that — the data crossing the strait should not be sealing under the same IKM as data that stayed home, and Pattern B preserves per-region cryptographic isolation for institutions whose regional regulatory regimes mandate in-region key custody."

Lin nods slowly. "FSC will look for that explicitly. They have asked, in past letters, about co-mingling of key material between the parent and the subsidiary. The §10.15 Pattern B framing is the one I want in the FSC letter — it is the spec naming the posture, not us."

Mike pulls a sample.

```
entry_id:               tw-rec-2026-04-09-00982341
user_id_hash:           sha256:4a8c...11de
model_id:               sunwon-rec-tw-v7
model_version:          7.2.0
features_hash:          sha256:c2f1...9087
recommended_skus:       [SK-77241, SK-91123, SK-44290]
gen_ai.request.model:   sunwon-rec-tw-v7
gen_ai.response.model:  sunwon-rec-tw-v7
audit.deployment.intent: production
served_at:              2026-04-09T08:14:22+08:00
seal_signature:         ed25519:b9d4...7e23
```

"Run the verifier," Lin says.

```
$ herald-verify --tenant=sunwon-tw-rec \
                --service=customer-rec \
                --date=2026-04-09

[verify] tenant=sunwon-tw-rec
[verify] hsm=taipei-chunghwa-hsm-02
[verify] entries scanned: 142,887
[verify] seal coverage: 365/365 days
[verify] sign_payload_version: v1.0b
Status: PASS
```

"Clean," Mike says.

> ✅ **Confirmation #2 — Per-jurisdiction tenant isolation under §10.15 Pattern B.**
> The KR and TW recommendation tenants use distinct HSMs, distinct IKMs, and distinct seal-cadence anchors per §4.2.1 cadence configuration. Per §4.1 HKDF binding, each regional tenant derives its session key from the per-tenant `info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)`, so the two tenants' session keys are byte-disjoint even on the same HSM. No co-mingling of key material across the strait. PIPA Section 28 and PDPA Article 8 both satisfied at the chain layer. §10.15 Pattern B (per-region `tenant_id`, single-tenant chain integrity per regional tenant) honored.

Lin makes a note in her own portfolio. "FSC will accept that. That is a clean answer."

---

## 🧠 10:00 AM — Seoul, Conference Room 28-A

Diana takes over for the IAM walk. Sun-Won uses **PASS-IT** — a Korean government-administered authentication provider that issues mobile-PKI tokens for KISA-certified services. PASS-IT is integrated into Sun-Won's SSO via SAML 2.0 with a Korean-resident IdP cluster.

"PASS-IT is the right choice for this jurisdiction," Diana says. "The Korean residency requirement in K-ISMS makes it the obvious pick. Is it the only auth path into the chain operator console? §10.5 HSM custody plus separation-of-duties controls require named operator identities; the institution's CC8.1 control description names them per §10.18."

"It is the only auth path into the production console," Park says. "Break-glass exists. Break-glass requires PASS-IT plus a Yubikey plus a recorded business justification. And the IKM behind the session-key derivation per §4.1 is held under §10.6 — 32-byte minimum length, generated under §10.6.1 (RNG of cryptographic strength, FIPS 140-3 attestation available on request). Sun-Won's IKM registry retention follows §10.9 — IKMs older than the chain's retention horizon stay in the registry indefinitely so any historical chain entry remains verifiable."

Diana asks to see a break-glass record. Park opens the access-log archive.

```
break_glass_event_id:   bg-2026-03-14-00007
operator:               kr-ops-009 (PASS-IT verified)
yubikey_serial:         5A:71:F2:9E
business_justification: "Investigating malformed seal entry on
                         sunwon-kr-rec, 2026-03-14. CCO approved."
duration_minutes:       42
actions_taken:          [read-only chain inspection]
post-event_review:      attested by Park Hye-jin, 2026-03-15
```

"One break-glass in the last quarter," Diana says. "Read-only. Forty-two minutes. Justified. Reviewed. That maps to a §10.2 operational event — the kind the verifier walks the same way as any other operational event under §4.4 with `chain_kind = 'operational'`. The break-glass record itself is chained, so a post-hoc edit of the justification surfaces as a per-event MAC mismatch at §7 step 9."

"That is the pattern," Park says.

Diana checks the per-tenant scoping. PASS-IT identities are bound to tenant scopes via the SSO claim set. Operator `kr-ops-009` has scope `sunwon-kr-rec` and `sunwon-kr-bnpl` only — not `sunwon-tw-rec`, not `sunwon-cross-inventory`, not the chatbot tenants.

"Per-tenant scoping is enforced at the IdP," Diana says. "The chain operator console will not let `kr-ops-009` even see the TW tenants. That is correct. Cross-tenant lift attempts would also be caught at §7 step 4 (per-entry binding — `event.tenant_id == header.tenant_id`) before any MAC compute, so even an operator who somehow got bytes lifted from the TW chain into a KR chain file would see `cross-chain lift detected at seq N`. Defense-in-depth from IdP scoping plus structural verifier rejection."

> ✅ **Confirmation #3 — IAM via PASS-IT, per-tenant scoping enforced at the IdP, with §7 step 4 backstop.**
> Korean SSO uses PASS-IT mobile-PKI as the only path into the production console per §10.5 HSM custody and operator-identity discipline. Break-glass requires PASS-IT + Yubikey + recorded justification + post-event CCO review; the break-glass event itself is chained as a §10.2 operational event under `chain_kind = "operational"` per §3 enumeration. Per-tenant scoping is enforced at the IdP claim set, not at application code. K-ISMS access-control requirements satisfied. Cross-chain lift is structurally rejected at §7 step 4 before any MAC compute, providing the cryptographic backstop to the IdP-side scoping.

"FSS will read this and stop reading," Park says. "That is exactly what they asked for in the 2024 letter."

Tom adds, half to his recorder: "And §10.18 says the runbook section for break-glass procedure has to cross-reference the spec section the requirement derives from. Park, your runbook?"

"Section 4.2 of the operations runbook names §10.5 HSM custody and §10.2 operational events. Inline, in the heading, per §10.18."

"Good."

---

## 🧠 10:00 AM — Taipei (9:00 AM CST), Compliance Suite

Elena walks the CRM side of the recommendation engine — the input side, where customer records become feature vectors.

The Taiwan-side CRM is called **Sun-Won-Connect-TW**, a Salesforce-derived deployment with localized fields for Taiwanese ID-card prefixes and the National Health Insurance card identifier (which is *not* used as a feature, but is in the CRM for shipping). Sun-Won runs the CRM as a §10.16 SaaS-edge mirror — a connector subscribes to Salesforce CDC and replicates each captured record into the institution's chain-instrumented store. The connector emits the chain entries from there. Per §10.16, the institution's CC8.1 names the four quantified lag bounds: **median lag 18 seconds, 95th-percentile SLO 90 seconds over a 30-day rolling window, alerting threshold 150 seconds (between 1× and 2× the SLO per §10.16), connector-outage RTO 5 minutes**. None of those numbers is "near real-time" or "low-latency mirror" wording — Sun-Won's runbook learned that lesson from the Northbridge engagement findings published last quarter. Per §4.4.6 SaaS-edge connector source attribution, every connector-emitted chain entry carries `audit.connector_source.system = "salesforce-cdc"`, `audit.connector_source.replay_id` (the Salesforce CDC ReplayId), `audit.connector_source.commit_timestamp` (the Salesforce-side commit time), `audit.connector_source.commit_user` (the Salesforce User ID), `audit.connector_source.lag_observed_ms`, and `audit.connector_source.change_kind`. The `run_id` derives from the Salesforce Account ID per §4.4.6's stable-`run_id` discipline — every CDC event for an account chains within the same run, and connector restarts pick up the existing tail per §10.25 run-resume rules rather than starting a new genesis.

"The features that actually feed the recommendation model," Elena says, "are which?"

Lin pulls up the documented feature list.

```
recommendation_features_tw:
  - skin_type_self_reported          # tier-1 PII, customer-entered
  - age_band                         # 18-24, 25-34, 35-44, 45-54, 55+
  - prior_purchase_categories        # category, not SKU
  - season                           # ambient
  - location_region                  # county-level, not address
  - language_preference              # zh_TW, en

excluded_features:
  - facial_features_from_photo       # explicitly excluded post-controversy
  - voice_features                   # explicitly excluded
  - precise_location                 # never collected
  - national_id                      # never used as feature
  - nhi_card_id                      # never used as feature
```

"The features_hash in the chain entry is the hash of the feature vector," Elena says. "Not the source data. And per §5 canonical-form rules, the feature vector is part of the canonical bytes the per-event MAC covers — the OTel envelope plus the `audit.*` namespace plus the `gen_ai.*` semconv attributes. The `payload_hash` is HMAC-SHA-256 over `prev_hash || canonical_bytes` per §4.1, so the absence of a feature is not a policy claim — it is a cryptographic claim about what bytes went into the MAC."

"Correct."

"And the excluded features — facial features from photo — are excluded at ingest, not at model time."

"Correct," Lin says. "After the controversy, we moved the exclusion to ingest. The model layer cannot see those fields because they never enter the feature pipeline. The pipeline produces the canonical bytes; the SDK MACs over them; the `payload_hash` is the cryptographic record that those bytes — and only those bytes — were what the model saw."

Elena writes this down. *That is the thing that matters about the post-controversy redesign,* she thinks. *Exclusion at the model layer is a promise. Exclusion at ingest is structural. And the chain promotes "structural" from a process word to a cryptographic word.*

Mike chimes in from his laptop. "And the redaction posture is §10.22 conformant — pre-MAC at the SDK boundary, not post-MAC sidecar. The captured JSON is the redacted form, and the per-event MAC covers the redacted form per §10.22's binary posture statement. An `audit.redaction.*` attribute set is emitted on entries where the SDK redacted PII before MAC compute — `redaction.policy_id`, `redaction.policy_version`, `redaction.redacted_field_paths`, `redaction.redaction_method`, `redaction.disposition = 'redacted_at_sdk'`. PIPA-compliant. PDPA-compliant. CFPB-compliant if a US examiner ever reads it."

> ✅ **Confirmation #4 — Excluded features removed at ingest, not at model layer; redaction discipline §10.22 conformant.**
> Facial features from photos and voice features are excluded at the ingest stage, before the feature pipeline. The model layer cannot observe them because they never enter the pipeline. The features_hash in each chain entry is over the actual feature vector under §4.1's canonical-bytes inclusion rule (§5) — making the absence of excluded features structurally provable, not policy-provable. Redaction occurs pre-MAC at the SDK boundary per §10.22 (Round-17 CFPB-P2 close-out): the `audit.redaction.*` attribute set names the policy_id, policy_version, redacted_field_paths, redaction_method, and disposition, all bound under the per-event MAC. Post-MAC sidecar redaction is non-conformant unless cross-anchored per §10.21-style anchor; Sun-Won does not run a sidecar. Composition with §5.2 best-evidence: the captured JSON is the redacted form; the canonical bytes integrity-bind the redacted form; both originals under FRE 1001(d).

---

## 🔐 11:00 AM — Seoul, Conference Room 28-A

Tom is reviewing Park's binder of prior FSS supervisory letters. There are nine of them across the past three years. He matches each one against the current chain configuration to see whether the chain has answered each prior question structurally.

The 2024-Q2 letter asked about reviewer-override capture — answered (Confirmation #1 above).

The 2024-Q4 letter asked about model versioning — answered, model_version is in every entry under the §4.4 GenAI semconv requirement that `gen_ai.response.model` be present on any entry carrying any `gen_ai.*` attribute.

The 2025-Q1 letter asked about cross-tenant operator access — answered (Confirmation #3 above).

The 2025-Q3 letter asked something Tom has to read twice.

> "The supervised entity should demonstrate that the BNPL model's training data does not include retail-side customer behavior absent explicit consent at the training-data collection point."

Tom looks up. "Park, this letter is about training-data lineage, not inference. Does the chain cover training?"

"The chain covers inference," Park says. "Training-data lineage is in a separate manifest. We attest to it quarterly. Per §1.2 epistemic scope, the chain proves what the AI said and that the record was not tampered with — it does not prove training-data lineage on its own. The manifest is the institution-side evidence regime that composes alongside the chain."

"Can I see the manifest?"

Park pulls it. The training manifest lists the data sources, their consent basis, and the cutoff timestamps. The BNPL model's training set draws from BNPL applications only — not from retail purchase history, unless the customer signed the cross-use consent at BNPL application time.

"Cross-use consent rate?" Tom asks.

"Forty-one percent. The model is trained on the consenting applicants only. The non-consenting applicants are not in the training set. That is enforced by the training pipeline, not by the model."

Tom nods. *That is the right answer for the 2025-Q3 letter. The chain does not need to cover training. The manifest covers training. They are separable concerns.*

He notes: "Training-data lineage out of chain scope per §1.2 (the chain proves what the AI said, not what the training set contained); manifest covers it; attested quarterly; enforced at pipeline. FSS 2025-Q3 letter answered. And the manifest itself — Park, is that hash-anchored to the chain via §10.19's `audit.external_artifact.*` family?"

Park: "Yes. Each quarterly manifest is published with a SHA-256, and the SHA-256 is anchored via `audit.external_artifact.kind = 'training_manifest_quarterly'`, `audit.external_artifact.identifier = 'kr-bnpl-train-manifest-2026-q1'`, `audit.external_artifact.sha256 = ...`, `audit.external_artifact.received_at_utc = ...`, `audit.external_artifact.source_party = 'sunwon_internal_ml'`, `audit.external_artifact.evidentiary_role = 'regulatory_compliance'`. The manifest is retained per §10.13 evidentiary artifacts retention plus §10.20 training-data retention floor."

Tom: "And §10.20 training-data retention vs deployment-window discipline — the BNPL model's deployment window is what?"

"Eighteen months on the longest path. Training shards retain for 24 months — the deployment window plus a six-month investigation buffer. That is above the §10.20 60-90-day floor. Documented in CC8.1 with §10.18 cross-reference."

> ✅ **Confirmation #4a — Training-data lineage hash-anchored under §10.19 + §10.20.**
> Quarterly training manifest is anchored on the chain via the §10.19 `audit.external_artifact.*` family (advisory; institution-controlled). Training-data shard retention is 24 months — 18-month deployment window + 6-month investigation buffer — above the §10.20 floor (deployment-window-plus-60-to-90-days). GDPR Article 5(1)(c) data-minimization tension resolved through the §10.20-named legitimate-interest determination. Composition with §1.2 epistemic scope: the chain proves what was deployed; the manifest proves what trained the deployed model; the §10.20 retention floor governs forensic depth.

---

## 🔐 11:00 AM — Taipei (10:00 AM CST), Compliance Suite

Luis pulls up the daily-seal landings on the Taipei HSM.

```
$ herald-verify --hsm=taipei-chunghwa-hsm-02 --seal-coverage --year=2025

[verify] HSM: taipei-chunghwa-hsm-02
[verify] daily seals expected: 365
[verify] daily seals landed:   365
[verify] gaps:                 0
[verify] late seals (>2 min):  0
[verify] empty-day seals:      6 (per §4.2 every tenant-day MUST receive a seal record)
[verify] PASS
```

"Six empty-day seals," Luis notes. "Days when the Taiwan chatbot tenant had zero customer interactions — typically Lunar New Year holidays. Per §4.2 (Round-17 close-out clarification), every tenant-day MUST receive a seal record, including tenant-days with zero events. Empty-day Merkle root pinned at `SHA-256(b"")` = `e3b0c44...b855` per §4.2. A missing empty-day seal would be reported as `missing seal for tenant-day {D}` — control-completeness failure, not chain-integrity failure. We get a clean PASS instead."

He reruns for the year-to-date 2026.

```
$ herald-verify --hsm=taipei-chunghwa-hsm-02 --seal-coverage --year=2026

[verify] daily seals expected: 99
[verify] daily seals landed:   99
[verify] gaps:                 0
[verify] PASS
```

"Three hundred sixty-five plus ninety-nine," Luis says. "Four hundred sixty-four consecutive daily seals on the Taipei HSM. No gaps. No late landings."

He runs the same on the Seoul HSM via the bridge — Diana, in Seoul, reads the result back.

```
$ herald-verify --hsm=seoul-sangam-hsm-01 --seal-coverage --year=2025

[verify] daily seals expected: 365
[verify] daily seals landed:   365
[verify] PASS
```

> ✅ **Confirmation #5 — Daily seals on both jurisdiction HSMs, no gaps, empty-day continuity preserved.**
> Seoul HSM (Sangam-dong, KISA-certified): 365 + 99 consecutive daily seals, zero gaps, zero late landings. Taipei HSM (Chunghwa Telecom data center, CNS 27001 aligned): 365 + 99 consecutive daily seals, zero gaps, zero late landings. Daily seal cadence operationally sound across both jurisdictions per §4.2.1. Empty-day seals (six on the Taiwan side during Lunar New Year holidays) recorded with the canonical `SHA-256(b"")` Merkle root per the §4.2 every-tenant-day-MUST-receive-a-seal rule. NTP discipline per §10.4 confirmed by audit-procedures P-7.

Luis adds, half to himself: "Two HSMs, two operations teams, two on-call rotations. Zero gaps. That is not luck. That is §10.5 HSM custody discipline plus §10.18 runbook cross-referencing executed without slippage for sixteen months."

Lin asks: "And the SaaS-edge connector — the Salesforce CDC mirror — has it had any outages?"

Luis pulls the `connector.outage` event log per §10.16. "Two outages in the past year. One was 2025-08-14, lasted 3 minutes 22 seconds — the Salesforce CDC subscription dropped after a network blip. The connector reconnected automatically; the back-log was 12 records, all caught up within the next 30 seconds. RTO of 5 minutes was satisfied. The second was 2026-01-09 during a planned Salesforce maintenance window — Sun-Won had it on the schedule, the connector emitted `connector.outage` at start, the back-log was 47 records, replay was clean. Both outages well under the §10.16 RTO."

He pulls the `connector.lag_observation` events. "Median lag actual: 14 seconds rolling 30-day. 95th-percentile actual: 64 seconds. Both well under the SLO numbers Sun-Won committed to in CC8.1. Alerting threshold of 150 seconds has fired zero times in the past year — the steady-state numbers are well clear of it."

> ✅ **Confirmation #5b — SaaS-edge mirror connector under §10.16 + §4.4.6.**
> Sun-Won's Salesforce CDC connector for the Taiwan CRM operates within the four quantified bounds named in CC8.1 per §10.16 (median 18s SLO, p95 90s SLO, alert 150s, RTO 5min) — actual rolling 30-day numbers (median 14s, p95 64s) clear the SLOs comfortably, and zero alert-threshold breaches in twelve months. Two outages in the past year, both within RTO, both replay-recovered cleanly, both chained as `connector.outage` operational events per §10.2 + §10.16. Per-entry `audit.connector_source.*` attribution per §4.4.6 lets an examiner cross-reference any chain entry to its Salesforce CDC ReplayId and source-side commit timestamp.

Lin smiles. "It is not. Park-CCO and our Taipei ops director have a running bet about who is going to break the streak first. Neither has."

Mike adds: "And the HSM partition ceremony attestation — §10.17. Each HSM partition creation, IKM rotation, and partition-PIN reset emits a `chain.partition_ceremony_attended` event with the signatories array, the witness signature, the SHA-256 hash of the attendance-log PDF, and per Round-17 M&A-P1 the `entity_affiliation` field on each signatory. Park-CCO is named on the Seoul ceremonies; Lin's controlling-person counterpart is on the Taipei ceremonies."

Luis: "I confirmed three ceremonies in the last year for each HSM — initial partition creation in 2024-Q4, IKM rotation in 2025-Q2, and the partition-PIN reset in 2026-Q1. All three carry the `chain.partition_ceremony_attended` event. The PDF retention per §10.17 is in the institution's compliance vault; the chain entry is the integrity-bound attestation; a discrepancy between PDF and chain event would be a P-6 control failure."

> ✅ **Confirmation #5a — HSM partition ceremony attestation under §10.17 (Wave-6 second errata close).**
> Both Seoul and Taipei HSMs operate under §10.17 dual-control attestation. Partition creation, IKM rotation, and partition-PIN reset ceremonies are emitted as `chain.partition_ceremony_attended` events with `ceremony_type`, `partition_handle`, `ceremony_started_at_utc`, `ceremony_completed_at_utc`, signatories array (with `role`, `name`, and `entity_affiliation` per Round-17 M&A-P1), `witness` object, and `attendance_pdf_sha256`. Composition with §10.5 HSM custody preserved. Cross-language CC8.1 discoverability per §10.17's last paragraph: the Korean and Mandarin operational runbooks are cross-referenced from the English CC8.1 by title, table-of-contents structure, and named ceremony-procedure section, so a non-Korean-reading or non-Mandarin-reading customer-bank auditor can identify and request translations.

---

## 🔐 11:30 AM — Seoul + Taipei, Bridge Open

The two halves of the team join the video bridge for the first time. Park and Lin are visible side-by-side on each other's screens, in their own conference rooms.

Dawn opens. "Before lunch I want to surface one item, so we can talk about it over food. The inventory-forecasting tenant. Park, Lin — you both know the shape. Can someone walk us through it?"

Park nods to Lin. Lin nods back.

"It is the only tenant that crosses the strait," Lin says. "Single tenant: `sunwon-cross-inventory`. Inventory data from Korean stores and Taiwanese stores both feed in. The model forecasts SKU-level demand at the regional warehouse layer. The model has to see both jurisdictions because the inventory rebalances between them — particularly at quarter-end and around lunar holidays. We operate it under §10.15 Pattern A — single seal region (Seoul), per-region event-count reconciliation, run-locality enforced via SDK per-process region binding per §4.4. Storage is append-only per §10.3 — UPDATE and DELETE operations on stored events are non-conformant under §6, and the institution's WORM-compatible storage layer enforces this at the storage tier. Per §10.1, daily key-fingerprint reconciliation runs against the IKM registry — a fingerprint mismatch surfaces at §7 step 8 before any MAC compute."

"Cross-border transfer basis?" Dawn asks.

"Contract," Park says. "There is a documented intra-group data transfer agreement between Sun-Won Holdings and Sun-Won Taiwan, registered with PIPC and acknowledged by PDPC. The contract identifies inventory data as a category, identifies the model as a recipient, and identifies the lawful basis as legitimate business interest with appropriate safeguards."

"And the chain entries?" Dawn asks.

Pause.

"The chain entries do not currently carry the cross-border transfer basis as an attribute," Park says. "The contract carries it. The compliance binder carries it. The chain entry shows that an inventory data point from store KR-Seoul-014 was used in a forecast — and that the forecast was generated on the cross-tenant. The cross-border basis is not stamped into the entry."

Dawn pauses. *That is the thing — and the spec moved on this. §4.4 added the `audit.cross_border_transfer.*` attribute family in the Wave-6 fourth erratum specifically because Sun-Won's posture surfaced the gap.*

"That is what I thought," Dawn says. "Lin, Park — that finding drove a spec amendment. The Round-17 Wave-6 fourth erratum (§12 change log) lifted the cross-border-transfer attribute family to the spec body. §4.4 normates it now. The attribute set is `audit.cross_border_transfer.contract_id`, `contract_version`, `contract_hash_sha256`, `source_jurisdiction`, `destination_jurisdiction`, `lawful_basis_type`. It is REQUIRED on entries subject to a regulator-named privacy regime that the institution's CC8.1 names — and PIPA Section 28 plus PDPA Article 8 are precisely those regimes. The advisory posture from the v1.0a draft is gone. The institution's CC8.1 names the trigger; the attribute set is REQUIRED whenever the trigger holds."

Park looks up. "REQUIRED, not advisory?"

"REQUIRED, when the institution's CC8.1 names a privacy regime that triggers it. Sun-Won's CC8.1 names PIPA §28 and PDPA Art 8 as triggers. The chain entries on `sunwon-cross-inventory` should be carrying the attribute set today and are not. That is no longer an advisory recommendation — it is a non-conformance against the v1.0b spec text under §4.4."

"Lunch?" Lin asks, drily.

"Lunch."

---

## 🧪 12:00 PM — Seoul + Taipei, Working Lunch via Video Bridge

Korean lunch in Seoul: bibimbap, mandu, pickled radish. Taiwanese lunch in Taipei: lu rou fan and pickled mustard greens. Both teams eat with their cameras on. Park and Lin are on the call together for the first time today.

Dawn does not let the lunch slide into chitchat. She takes the inventory tenant head-on.

"Park, Lin. The inventory tenant is the one that needs the conversation. The chain works. The contract works. The two pieces of evidence sit in different binders. If FSS, PIPC, FSC, and PDPC all asked the same question on the same day — show me the cross-border basis for this one inventory data point — could you give all four the same answer?"

Park considers. "Today, the answer is: the chain entry plus the contract reference. The chain proves what the inventory model saw. The contract proves the consent basis. They have to be read together."

"Lin?"

"FSC and PDPC will accept that," Lin says, "but they will note in their letter that the link between the two is procedural rather than cryptographic. PDPC has been asking for cryptographic linkage in cross-border-flow attestations since 2024."

Dawn takes a breath. "Right — and that is now the spec's normative answer. §4.4 binds it. The attribute set lives on the chain entry. The contract is institution-published; the contract's `contract_hash_sha256` anchors it to the entry. A post-hoc edit of the contract is detectable. The link is no longer procedural. It is cryptographic. Per §12 (Wave-6 fourth erratum), this story drove that amendment — Sun-Won's exact posture, surfaced as a Partial in the v1.0a draft, was the worked example the spec used to lock the family. So the recommendation is not 'consider adding a future capability' — it is 'remediate to the v1.0b spec text by emitting the attribute set the spec already names.' The CC8.1 update declares which privacy regimes are triggers; the SDK update emits the attribute set on cross-jurisdiction entries; the contract repository publishes the contract as a versioned, hash-anchored document so the `contract_version` and `contract_hash_sha256` resolve."

Mike, on the Taipei side, leans into his camera. "Mechanical work. Six attributes per cross-border entry. The inventory tenant's attribute set is six fields today; adding the cross-border-transfer six is a config change at the SDK boundary, not a code change. The contract-versioning side is the work — Sun-Won's compliance team has to publish the contract as a versioned, hash-anchored document the SDK can resolve at MAC time. The attribute set is part of the canonical bytes per §5, so a tampered `contract_hash_sha256` would surface as a MAC mismatch at §7 step 9."

"Six weeks of work on the legal side," Lin says. "Two weeks on the chain side. Maybe."

"This is no longer 'advisory recommendation,'" Dawn says. "It is 'remediation to the v1.0b spec.' The current setup answered the regulators in 2025 under the v1.0a posture; under the v1.0b posture (active since the February SDK upgrade Sun-Won is on), the attribute set is REQUIRED whenever the institution's CC8.1 names a privacy-regime trigger. Sun-Won's CC8.1 names the triggers. The remediation closes the gap to v1.0b conformance."

> ⚠️ **Finding-001 — Cross-border transfer basis not stamped into inventory chain entries (non-conformance against v1.0b §4.4).**
> The `sunwon-cross-inventory` tenant aggregates inventory data from both Korea (source) and Taiwan (destination, and vice versa for KR-bound forecasts). The lawful basis for cross-border transfer is documented in an intra-group data transfer agreement registered with PIPC and acknowledged by PDPC. The contract is sound. However, the chain entries do not carry the `audit.cross_border_transfer.*` attribute set per §4.4. **Severity: non-conformance against §4.4.** Sun-Won's CC8.1 names PIPA §28 and PDPA Art 8 as privacy-regime triggers; the v1.0b §4.4 elevation makes the attribute set REQUIRED on chain entries subject to those regimes, not advisory. **This story drove the Wave-6 fourth erratum (§12 change log)** that lifted the attribute family to the spec body. The institution's posture is the worked example the spec used to lock the family. Remediation: publish the intra-group data transfer agreement as a versioned, hash-anchored document; emit `contract_id`, `contract_version`, `contract_hash_sha256`, `source_jurisdiction`, `destination_jurisdiction`, `lawful_basis_type` (= `intra_group_agreement`) on every `sunwon-cross-inventory` entry. Six attributes; bound under the per-event MAC per §5; cryptographic linkage between chain and contract closes the procedural-vs-cryptographic gap the v1.0a posture left. Target: Q3 2026.

Park writes the recommendation into her binder. "Q3. The legal versioning is the long pole, but it is now spec-mandated, not advisory. That changes the internal conversation."

"Agreed," Lin says. "And FSC will read 'remediation to spec' as a stronger commitment than 'advisory upgrade.' The framing helps us."

---

## 🔄 1:00 PM — Taipei (12:00 PM CST), Compliance Suite

After lunch, Mike turns to the chatbot.

The multilingual customer-service chatbot has three tenant configurations:

```
tenant_id:        sunwon-chatbot-ko    (Korean model — kakao-style)
tenant_id:        sunwon-chatbot-zh    (Mandarin model — Taiwan dialect)
tenant_id:        sunwon-chatbot-en    (English fallback)
```

The chatbot's first step, on every customer interaction, is **language detection**. The detector is a separate microservice — a small fastText classifier running in the Taipei region. It looks at the first 50 characters of the customer's message and emits a language label: `ko`, `zh`, or `en`.

The router then dispatches to the appropriate tenant. The chain entry, once the model serves a response, is sealed under the tenant that actually served — `sunwon-chatbot-ko` if the Korean model served, `sunwon-chatbot-zh` if the Mandarin model served, and so on.

Mike pulls a sample chain entry.

```
entry_id:               tw-chatbot-zh-2026-04-09-00128941
user_id_hash:           sha256:7a82...e441
model_id:               sunwon-chatbot-zh-v4
model_version:          4.1.7
gen_ai.request.model:   sunwon-chatbot-zh-v4
gen_ai.response.model:  sunwon-chatbot-zh-v4
prompt_hash:            sha256:b193...7c0f
response_hash:          sha256:c4d2...8e9b
audit.deployment.intent: production
served_at:              2026-04-09T11:42:18+08:00
seal_signature:         ed25519:b9d4...7e23
```

"The chain shows which model served," Mike says. "It does not show how the routing decision was made. The detector logs are separate."

"Where do the detector logs live?" Lin asks.

"Separate microservice, separate log system. Ninety-day retention, by default. Beyond ninety days, the detector input, the per-class scores, and the chosen output are not recoverable from anywhere — chain or detector logs."

Lin presses. "And the chatbot models themselves — the Korean kakao-style model, the Mandarin Taiwan-dialect model, the English fallback — those are vendor-supplied?"

Mike nods. "Two of the three. The Korean model came from a Seoul-based AI consultancy under a model-supply contract. The Mandarin and English models are in-house. The §10.21 cross-vendor model-handover discipline applies to the Korean model — `audit.model_handover.provider`, `audit.model_handover.model_id`, `audit.model_handover.model_version`, `audit.model_handover.model_artifact_sha256`, `audit.model_handover.model_card_sha256`. Sun-Won emits the handover entry on each model upgrade. And per Round-17 M&A-G2, the `audit.model_handover.contract_id`, `contract_version`, and `contract_hash_sha256` are present too — the model-supply contract is hash-anchored, so a post-close auditor (if Sun-Won is ever acquired) could answer 'which contract version governed this delivery' from the chain alone."

Lin: "And the `audit_report_languages` array — plural per the M&A-N2 close?"

Mike: "Plural array. The Korean consultancy's fairness audit report is in Korean and English; the chain entry records `audit.model_handover.audit_report_languages = ["ko", "en"]` per §10.21's plural-array discipline. A non-Korean-reading vendor-management auditor finds the English translation through the chain itself."

> ✅ **Confirmation #6a — Vendor model handover under §10.21 (Round-17 M&A-G2 contract binding).**
> The Korean chatbot model is supplied by a Seoul AI consultancy. Each handover emits `audit.model_handover.*` per §10.21 with provider, model_id, model_version, model_artifact_sha256, model_card_sha256, fairness_audit_report_sha256, audit_report_languages (plural), provider_chain_entry_id, training_data_retention_floor_days. The Round-17 M&A-G2 contract-binding extension applies — `contract_id`, `contract_version`, `contract_hash_sha256` cryptographically link the chain entry to the supply contract. Per §10.20 training-data retention vs deployment-window discipline, the consultancy's training-data shard retention is committed contractually at 24 months — above the §10.20 60-90-day floor and matching Sun-Won's deployment-window-plus-investigation-buffer requirement.

Mike runs the verifier on the Mandarin tenant.

```
$ herald-verify --tenant=sunwon-chatbot-zh \
                --service=chatbot \
                --date=2026-04-09

[verify] tenant=sunwon-chatbot-zh
[verify] hsm=taipei-chunghwa-hsm-02
[verify] entries scanned: 9,217
[verify] seal coverage: 365/365 days
Status: PASS
```

"Clean," Mike says. "The chain shows the Mandarin model served 9,217 interactions yesterday. It does not show why those particular interactions were routed to the Mandarin model rather than the Korean or English ones."

> ✅ **Confirmation #6 — Chatbot per-language tenant separation.**
> Three chatbot tenants — Korean, Mandarin, English — each with its own model, model_version, IKM, and seal stream. Per-language separation is structural at the tenant layer. PIPA and PDPA both want the model populations separated; both are satisfied. §4.1 HKDF binding plus §10.15 Pattern B per-region tenant isolation jointly honored. `gen_ai.request.model` and `gen_ai.response.model` present on every entry per the §4.4 SDK-side enforcement rule (the SDK refuses to emit a chain entry whose attribute set includes any `gen_ai.*` namespace prefix attribute AND lacks either model identifier).

Dawn, on the Seoul bridge, leans in. "Mike — the routing decision. That's the §4.4.1 question."

Mike nods. "Right. §4.4.1 is the routing-classifier event family. And per the Wave-6 fourth erratum (§12), this story drove the spec's sixth event type — `audit.routing.classifier_output`. The v1.0a §4.4.1 had five event types: `attempt`, `success`, `failover`, `circuit_state_change`, `refused`. None of those covers a pre-routing classifier. The v1.0b amendment added `classifier_output` precisely for this case."

Lin looks up. "Read me the schema."

Mike reads from the spec on his second screen. "§4.4.1 sixth event type, emitted BEFORE the `audit.routing.attempt` it informs, linked via `parent_run_id` / `parent_seq` per §4.4 — classifier_output is the parent of the attempt. Six new attributes, all REQUIRED on the classifier_output event: `audit.routing.classifier_name` — the classifier service or model identifier; `audit.routing.classifier_version` — version identifier for the classifier; `audit.routing.classifier_input_hash` — SHA-256 lowercase hex of the canonicalized classifier input; `audit.routing.classifier_scores` — JCS-canonical object mapping class identifier to score; `audit.routing.classifier_decision` — the class identifier the classifier selected, MUST be a key in `classifier_scores`; `audit.routing.classifier_confidence` — confidence in [0.0, 1.0] for the chosen class. Per the §4.4.1 normative text: 'Without the pre-routing entry, reconstructing why a user was routed to a specific provider depends on the classifier service's logs, which typically retain shorter than the chain itself; pre-chaining the classifier output makes the rationale recoverable from the chain alone for the chain's full retention period.' That is verbatim what the chatbot is missing. Institutions whose routing policy is purely rule-based MAY omit; institutions with classifier-driven routing — like Sun-Won's chatbot — MUST emit it."

Dawn: "Park, Lin — same posture as the cross-border attribute. This is no longer 'advisory.' The spec body normates it now under §4.4.1. Sun-Won's chatbot operates classifier-driven routing; the spec MUST applies. Sun-Won is non-conformant against v1.0b until the classifier_output event is emitted."

Lin reads through her screen. "And — the spec text says the entry is BEFORE the attempt event it informs. The classifier_output is the parent. So we have to chain the detector before the model selection, not as a side annotation."

"Right. Linked via `parent_run_id` / `parent_seq` per §4.4. The two entries are a parent-child pair. The classifier_output sits in the chain at `(run_id, classifier_seq)`; the `attempt` sits at `(run_id, attempt_seq)` with `parent_run_id = run_id, parent_seq = classifier_seq` — even though they are within the same run, the parent-linkage discipline binds them. Per §5, those linkage fields are part of the canonical bytes the per-event MAC covers, so the linkage itself is integrity-bound."

> ⚠️ **Finding-002 — Language-detection routing decision not chained ahead of model selection (non-conformance against v1.0b §4.4.1).**
> The chatbot's first step is a language-detection microservice that picks which model serves the user. The chain entry records which model served — not which classifier output drove the routing. Reconstructing "why was this user served by the Mandarin model rather than the Korean model" requires the language-detector microservice's logs, which live in a separate log system on a 90-day retention. Beyond 90 days, the routing rationale is not recoverable from any source. **Severity: non-conformance against §4.4.1.** Sun-Won's chatbot operates classifier-driven routing; the spec text "Institutions with classifier-driven routing MUST emit it" applies. **This story drove the Wave-6 fourth erratum (§12 change log)** that added the sixth event type `audit.routing.classifier_output` and the six new classifier attributes to §4.4.1. Remediation: emit `audit.routing.classifier_output` chain entries BEFORE the `audit.routing.attempt` they inform, linked via `parent_run_id` / `parent_seq` per §4.4. Required attributes: `classifier_name` (= `language-detector-fasttext-v3`), `classifier_version`, `classifier_input_hash`, `classifier_scores` (per-class scores: ko / zh / en), `classifier_decision` (the chosen class), `classifier_confidence`. The classifier's input itself MAY be retained separately under detector-log retention; the input hash on the chain entry is the load-bearing reference per §4.4.1's normative text. Bundle remediation with Finding-001's Q3 work: same SDK config-change shape; same six-week order of magnitude.

Lin looks at the recommendation. "I prefer the chained approach over a log-retention extension. Chaining makes the answer self-contained; log retention extensions get rolled back when finance reviews them."

Mike nods. "Same. And the spec made the choice for us — §4.4.1 normates the chained event type; log-retention extension is not a §4.4.1-conformant alternative. It is the spec answering the question."

Dawn: "Required-pairing rule per §4.4.1 — for a classifier-driven routing decision, the chain MUST carry classifier_output → attempt → success/failover. Audit-procedures.md P-33 samples for the pairing. Sun-Won's audit cycle has to add the P-33 procedure once the remediation lands."

Lin asks one more thing. "What about the BNPL declination path? When the BNPL model declines an applicant, does the §10.11 adverse-action translation apply?"

Park nods. "FSS does not enforce ECOA — that is US-jurisdictional. But §10.11's last paragraph applies the adverse-action translation discipline by analogy to state-insurance-law adverse-action notices, and the Round-17 NAIC-N2 close renamed §10.11 to 'Adverse-action notice translation (ECOA and state-insurance analog)' specifically so a non-US regulator citing the section reaches the analogous discipline. We apply the §10.11 attribute schema to BNPL declinations sent to Korean-language consumers — `audit.ecoa.translation.target_language = "ko-KR"`, `translator_kind = "human"` (Sun-Won uses human translators for legal-impact letters), `output_hash` over the customer-facing translated declination text, `delivery_method = "mail"`, `delivery_timestamp` per the Round-17 CFPB-N1 close requiring `delivery_timestamp` whenever `delivery_method` is recorded. The translation entry's `chain_kind = "translation"` per §3 enumeration. Parent-linkage to the underlying decline entry via `parent_run_id` / `parent_seq` per §4.4."

Dawn: "And the underlying decline entry carries the §10.11.1 adverse-action reasons schema — `audit.ecoa.adverse_action.reasons` (the structured-reasons code), `audit.ecoa.adverse_action.feature_attributions` when the model exposes attribution at decision time, `audit.ecoa.adverse_action.model_explanation_method`. The chain proves what reasons the model produced and what text the consumer received. Two integrity-bound records, one chain. FSS-grade."

> ✅ **Confirmation #6b — BNPL adverse-action discipline under §10.11 + §10.11.1 (analog application).**
> §10.11 normates ECOA adverse-action notice translation; per the section's last paragraph (Round-17 NAIC-N2), the discipline applies by analogy to state-insurance-law adverse-action notices and equivalent regimes — including FSS-supervised BNPL declination notices. Sun-Won emits the §10.11 `audit.ecoa.translation.*` attribute set on each declination translation entry, and the §10.11.1 `audit.ecoa.adverse_action.*` family on the underlying decline entry per Round-17 CFPB-P1. `delivery_timestamp` is REQUIRED whenever `delivery_method` is recorded, per Round-17 CFPB-N1. Composition with §1.2 epistemic scope: the chain proves what reasons the model produced and what text the consumer received; it does not prove either reason was the *actual* policy-compliant reason — that proof lives in Sun-Won's policy-as-code system per §1.2's separate-evidence-regime framing.

---

## 🧬 2:00 PM — Taipei (1:00 PM CST), Compliance Suite

Chen takes the inventory-forecasting cross-jurisdiction data flow.

The architecture:

```
Korean stores (300+) ──┐
                       ├──> Cross-jurisdiction data lake (Seoul region — seal region)
Taiwanese stores (80+) ┘                │
                                        ▼
                            Inventory forecasting model
                          (sunwon-cross-inventory tenant; §10.15 Pattern A)
                                        │
                                        ▼
                          Forecast output → warehouse rebalancing
```

The data lake holds inventory data only — SKU-level stock levels, sales velocity, returns, regional demand signals. No customer-side data. No PII. The PIPA Section 28 question is whether the inventory data itself counts as personal data; PIPC has held in prior letters that aggregated SKU-level inventory data does not, but the cross-border transfer agreement covers it anyway as a precaution.

Chen hash-anchors the input data feeds.

```
$ herald-verify --tenant=sunwon-cross-inventory \
                --service=inventory-forecast \
                --date=2026-04-09 \
                --check-input-anchors

[verify] tenant=sunwon-cross-inventory
[verify] §10.15 pattern: A (single seal region — Seoul)
[verify] input feeds: 2 (kr-stores, tw-stores)
[verify] kr-stores feed anchor: sha256:8e2c...41bf  (matches manifest)
[verify] tw-stores feed anchor: sha256:3a91...0d22  (matches manifest)
[verify] entries scanned: 2,194
[verify] seal coverage: 365/365 days
[verify] master.cross_region_replication_completed: 365 (per §10.15 invariant 5)
Status: PASS
```

"Both feeds anchored," Chen says. "The chain shows that a forecast generated on April 9 used the Korean feed at hash 8e2c... and the Taiwan feed at hash 3a91... — and the manifest confirms those hashes correspond to the inventory snapshots taken at midnight KST and midnight CST respectively."

"And the §10.15 Pattern A reconciliation?" Lin asks.

Chen pulls it up. "Per §10.15 invariant 5, every Pattern A tenant emits the `master.cross_region_replication_completed` operational event on each seal day. The event records the per-region event count and the replication-completion timestamp the seal region observed. The Round-17 Wave-6 third erratum tightened this — the per-region count and timestamp MUST reflect the replication pipeline's actual state at emission time, not a poll-cached representation. Sun-Won's implementation reads synchronously against the replication pipeline's state at emission per the §10.15 invariant 5 normative text."

He runs the reconciliation:

```
$ herald-verify --tenant=sunwon-cross-inventory \
                --check-replication-completeness \
                --date=2026-04-09

[verify] seal region:        seoul-sangam (per §10.15 invariant 3)
[verify] kr-stores region count:  1,832
[verify] tw-stores region count:    362
[verify] sum:                    2,194
[verify] seal-region count:      2,194
[verify] match:                  yes
[verify] master.cross_region_replication_completed event: present
[verify] freshness: synchronous-read at emission (per §10.15 invariant 5)
[verify] PASS
```

"So if PDPC asks 'show me which Taiwanese inventory snapshot fed the cross-tenant on April 9'," Lin says, "the chain answers."

"Yes. Per §4.4 chain envelope and the per-feed anchor in the canonical bytes."

"And if PIPC asks 'show me that the Korean inventory snapshot was the only Korean data that fed the cross-tenant'," Chen says, "the chain answers — kr-stores feed anchor, no other Korean source. The cross-jurisdiction data-flow traceability is structurally provable at the data-lineage level."

"Good."

The cross-border-transfer-basis attribute remains the gap — already noted at lunch under Finding-001. The hash-anchor side of the cross-jurisdiction flow is clean and the §10.15 Pattern A reconciliation is clean.

> ✅ **Confirmation #7 — Cross-jurisdiction inventory feeds hash-anchored, both directions traceable; §10.15 Pattern A reconciliation clean.**
> The `sunwon-cross-inventory` tenant's input feeds — Korean stores and Taiwanese stores — are both hash-anchored at the daily snapshot point. The chain entries reference the input anchors, allowing per-forecast attribution back to the originating jurisdiction's inventory snapshot. §10.15 Pattern A operated cleanly: single seal region (Seoul), per-region event-count reconciliation with synchronous-read freshness per Round-17 Wave-6 third erratum's §10.15 invariant 5 tightening, `master.cross_region_replication_completed` operational event emitted per §10.2 every seal day, no replication gaps. Cross-border data-flow traceability is structurally provable at the data-lineage level. The remaining gap is the cross-border-transfer-basis attribute (Finding-001 above).

---

## 📊 3:00 PM — Seoul, Conference Room 28-A

Reconciliation test, Seoul side. Dawn has Raj pull five BNPL credit decisions from yesterday and trace each one end-to-end.

```
Sample 1: applicant_id_hash:e7c2...9f1a, decision:conditional → traced
Sample 2: applicant_id_hash:1d4a...c0f7, decision:approved (override) → traced
Sample 3: applicant_id_hash:6b81...3e92, decision:declined → traced
Sample 4: applicant_id_hash:f0a4...2c1d, decision:approved → traced
Sample 5: applicant_id_hash:9c0e...77a8, decision:conditional → traced
```

Each entry resolves to:

- The applicant's submitted features (hash-matched to the manifest under §4.4.5 `audit.underwriting.features.feature_vector_hash`).
- The model version that scored the application (`gen_ai.response.model` per §4.4 MUST requirement).
- The decision and any reviewer override.
- The seal under which it landed (per §4.2 + §4.3 + §7 step 11 v1.0b 12-line `sign_payload` form).
- The PASS-IT identity that recorded the override on Sample 2 (per §10.5 + §10.2 operational event linkage).

5 of 5 PASS. Tom records each one.

> ✅ **Confirmation #8 — BNPL reconciliation, 5 of 5 traced end-to-end.**
> Five BNPL credit decisions from 2026-04-08 traced from applicant submission to sealed chain entry to reviewer override (where applicable) to operator identity. All five resolve cleanly. FSS examiner-grade reconciliation discipline. §4.4 chain envelope + §4.4.5 underwriting features family + §4.4.2 deployment-intent + §10.5 HSM custody + §10.11.1 ECOA-analog adverse-action reasons schema (applied to FSS-jurisdiction by §10.11's analogy clause for state-insurance and equivalent regimes) jointly honored. §7 verifier procedure executed end-to-end including JCS self-test pre-flight, structural walk, MAC recompute under `expected_prev_hash`, Merkle recomputation, and signature verification dispatch on `sign_payload_version = "v1.0b"`.

---

## 📊 3:00 PM — Taipei (2:00 PM CST), Compliance Suite

Reconciliation test, Taipei side. Mike pulls five chatbot interactions from yesterday — three Mandarin, one Korean (a Taiwan-resident Korean speaker), one English — and traces each.

```
Sample 1: zh, served_at 2026-04-08T10:14, model:sunwon-chatbot-zh-v4 → traced
Sample 2: zh, served_at 2026-04-08T11:22, model:sunwon-chatbot-zh-v4 → traced
Sample 3: ko, served_at 2026-04-08T13:45, model:sunwon-chatbot-ko-v4 → traced (in TW region)
Sample 4: en, served_at 2026-04-08T15:08, model:sunwon-chatbot-en-v4 → traced
Sample 5: zh, served_at 2026-04-08T19:31, model:sunwon-chatbot-zh-v4 → traced
```

Each entry resolves to:

- The user_id_hash (per institution's §10.22 redaction discipline — the user-side identifier hashed pre-MAC).
- The prompt_hash and response_hash.
- The model_id and model_version (per §4.4 `gen_ai.{request,response}.model` MUST requirement).
- The seal under which it landed.

But the routing rationale — *why was this interaction sent to this model and not another* — requires the language-detector microservice's logs. Mike pulls the detector logs for the same five samples.

```
Sample 1: detector input "我想找適合敏感肌的精華液", classified zh (confidence 0.987) → log hit
Sample 2: detector input "請問口紅有什麼顏色", classified zh (confidence 0.992) → log hit
Sample 3: detector input "보습 크림 추천해 주세요", classified ko (confidence 0.961) → log hit
Sample 4: detector input "do you ship to Singapore", classified en (confidence 0.978) → log MISS (88-day-old)
Sample 5: detector input "口紅含色素嗎", classified zh (confidence 0.989) → log MISS (89-day-old)
```

Mike pauses on the misses. "Samples 4 and 5 fall outside the 90-day detector log retention. The chain entries are intact. The routing rationale is not recoverable."

Lin: "Acceptable for these two? Yes — they are old. The point is the structural exposure: anything older than ninety days, the routing rationale is gone."

"Yes," Mike says. "3 of 5 routing rationale recoverable from the detector logs. Chain entries 5 of 5 PASS at §7. The exposure is the routing side, not the chain side. And per §4.4.1 — once Finding-002 is remediated and the `audit.routing.classifier_output` events land on the chain — the routing rationale is recoverable from the chain alone for the chain's full retention period, which is bounded by §6 storage retention plus the institution's regulatory retention setting (typically 7 years for FSS-supervised entities), not by the detector's 90-day log retention. The 90-day exposure becomes a 7-year window."

> ⚠️ **Finding-003 (closure path) — Chatbot reconciliation: 5/5 chain PASS, 3/5 routing rationale recoverable beyond detector log retention.**
> Five chatbot interactions traced. Chain entries: 5 of 5 PASS at §7. Routing rationale: 3 of 5 recoverable from the language-detector microservice logs; the other 2 fell outside the 90-day detector log retention. **This finding is the worked example for Finding-002's spec amendment** — the Wave-6 fourth erratum (§12) lifted `audit.routing.classifier_output` to spec body precisely because Sun-Won's chatbot reconciliation showed that without pre-chained classifier output, the routing rationale's recoverability is bounded by the classifier service's log retention rather than by the chain's retention period. **This story drove the §4.4.1 sixth-event-type addition.** Once Finding-002 is remediated (Q3 2026), the routing rationale is recoverable from the chain alone for the chain's full retention period (typically 7 years for FSS supervision). Severity: tracking finding rather than independent non-conformance — the underlying remediation is Finding-002.

---

## 😬 3:45 PM — Seoul + Taipei, Bridge Open

Park and Lin both join the bridge. Dawn opens.

"We have two Findings against v1.0b spec — the cross-border attribute (Finding-001) and the language-detection-routing chained event (Finding-002). Both drove the Wave-6 fourth erratum that landed the §4.4 cross-border-transfer family and the §4.4.1 classifier_output event. The chatbot reconciliation 3-of-5 result (Finding-003) is the worked example for Finding-002 rather than an independent non-conformance — same remediation closes both. The pre-chain era retention question is a separate discussion — let's discuss it last."

Park nods.

"On the cross-border attribute," Dawn says, "Park and Lin both heard my recommendation at lunch. The spec already lifted the family to §4.4 normative text under the Wave-6 fourth erratum. Hash-anchor the contract version. Six weeks of legal work to publish the contract as a versioned hash-anchored document, two weeks of chain config work to add the attribute set. Park, Lin — disagreement on the framing as 'remediation to spec' rather than 'advisory upgrade'?"

Park: "No disagreement. The spec moved. The remediation is to v1.0b conformance."

Lin: "Agreed. I want to add — FSC will appreciate this in the next supervisory letter. They have asked for cryptographic linkage in cross-border evidence since 2024; the spec amendment delivers it; Sun-Won remediating to spec is the cleanest framing for that letter."

"Good. On the language-detection-routing," Dawn says, "the Wave-6 fourth erratum normates the chained classifier_output event under §4.4.1. Mike — your reading?"

"Spec answers the question," Mike says. "§4.4.1 sixth event type, six new attributes, parent-linkage to the attempt. The previous v1.0a posture left log-retention extension as a possible alternative; v1.0b's §4.4.1 normative text closed that — institutions with classifier-driven routing MUST emit the chained event. Sun-Won's chatbot operates classifier-driven routing. The spec applies. Remediation is to v1.0b conformance."

"Agreed," Park says. "Make it part of the same Q3 work as the cross-border attribute."

Lin: "Agreed. Both are essentially the same shape of work — one new attribute set on an existing tenant, one new entry-type for a pre-routing event. Bundle them."

"Done."

Dawn pauses. *This is where the question Park has been waiting to ask is going to come.*

It does.

"Dawn," Park says. "If FSS, PIPC, and FSC ask the same question — 'demonstrate consent for cross-border transfer of an inventory data point' — can we give all three the same answer today?"

The room goes a little quieter.

Dawn takes a beat.

"Today, the answer is the chain plus the contract. The chain proves what the inventory model saw under §1.2's epistemic scope — the captured event is integrity-bound under §1.4's compositional security (per-event MAC + daily Merkle seal + HSM-rooted signature). The contract proves the consent basis, currently as institution-side parallel evidence. You have to hand over both — and both have to be read together to construct the full answer. That is procedurally sound but evidentially compound.

After the Q3 remediation lands — the `audit.cross_border_transfer.*` attribute set on inventory entries per the v1.0b §4.4 amendment — the chain plus the verifier output answers all three regulators on its own. The contract is referenced inside the chain entry by `contract_id`, `contract_version`, and `contract_hash_sha256`. Verifier dumps the entry, the entry shows the contract version, the contract repository serves the contract by hash, and the regulator gets a single self-consistent evidence package. Per §5.2 best-evidence posture, the captured JSON is the content-bearing form and the canonical bytes are the integrity-bearing form — both originals under FRE 1001(d) but the regulator does not have to assemble them from separate binders.

That is the upgrade you are paying for. Today: two binders, both required, link is procedural. After Q3: one verifier output, contract referenced by hash, link is cryptographic per §4.4 + §5, single answer to three regulators."

Park writes that down word for word. *That is the language she will use in her summary memo.*

Lin: "Two binders today, link procedural; one verifier output after Q3, link cryptographic per spec §4.4. That's the right framing for FSC."

---

## 🔍 4:30 PM — Seoul + Taipei, Bridge Open

The pre-chain era — the celebrity controversy lookback.

Dawn takes it directly. "The chain was deployed sixteen months ago. The celebrity controversy was eighteen months ago. There are about four months of recommendation-engine activity from before the chain that fall inside the lookback window. We cannot verify those four months through the chain because the chain did not exist for them. Per §1.2 epistemic scope, the chain proves what the AI said and that the record was not tampered with — for the pre-chain window, neither claim is available because no chain entries exist."

Park nods. "We acknowledge that. The legacy recommendation-engine logs cover those four months. They are append-only on a write-once-read-many storage tier — that was already best practice before the chain was deployed. They are admissible. They are not chain-grade."

"Sun-Won is being honest about that," Dawn says. "I appreciate it. The audit deliverable will document the chain's effective-start date and note that pre-chain activity is verifiable through the legacy logs only. Per the Wave-6 fourth erratum (§12 change log), the spec explicitly names this as an institution-side legacy-log dependency rather than a spec concern — §1.2 epistemic scope plus the institution's CC8.1 chain-coverage map per §10.19 jointly carry the framing."

She pulls up §10.19. "And the §10.19 chain-coverage map is the right place for this. The map enumerates 'institutional systems not yet chain-instrumented' — the pre-chain era recommendation engine fits that category exactly. Sun-Won's CC8.1 names the rollout posture as 'deferred (effective end-of-period 2024-Q4)' and the evidentiary substitute as 'legacy WORM-storage logs'. That is the §10.19 map's exact shape. Per Round-17 M&A-P3, the chain-coverage map itself is version-stamped and chain-anchored — every publication or update emits the `chain.coverage_map_published` operational event under §10.2 carrying `coverage_map_version`, `effective_utc`, and `coverage_map_sha256`. Sun-Won's current map version is `v3.2.0`, effective 2026-01-15 after the chatbot-tenant addition."

Dawn turns to Tom. "Tom — for the litigation framing, what does §1.1 give us if Sun-Won ever ends up in court over the celebrity-controversy lookback?"

Tom: "Daubert four-factor grounding. Per §1.1's four-factor analysis: testability is the §7 verifier procedure with normative reason strings; peer review is the OpenTelemetry ecosystem plus the FFIEC reference verifier; known error rate is bounded by the three-layer compromise model in §1.1 plus the §1.2 fourth class (SDK-process compromise). General acceptance is RFC 8785 + RFC 6962 + FIPS 198-1 + FIPS 186-5 — all published Internet Standards or FIPS. An expert witness laying foundation under FRE 702 has the full residual-risk picture. For the pre-chain era, §1.1 doesn't help — there's no chain to ground in Daubert — but the legacy WORM-storage logs admissibility under FRE 803(6) (records of regularly conducted activity) plus §10.13 evidentiary artifacts retention plus §10.14 trusted-time integration (RFC 3161 RECOMMENDED for v1.0; Sun-Won has not adopted yet but is on the v1.x roadmap) gives Sun-Won an admissible-but-not-chain-grade record for that window."

"Good. That goes in the deliverable framing."

> ✅ **Confirmation-by-spec — Pre-chain era retention gap (the celebrity-controversy lookback) is institution-side legacy-log dependency, not a spec concern.**
> The chain was deployed in 2024-Q4. The celebrity-controversy lookback window extends back four months prior to chain deployment (mid-2024). Pre-chain activity in those four months is verifiable through legacy recommendation-engine logs (append-only WORM storage), not through the chain. Per §1.2 epistemic scope, the chain does not retroactively cover events that predate its deployment — the chain's claims (what the AI said + non-tampering) attach to events captured under the chain's primitives. Per §10.19 chain-coverage map, the institution's CC8.1 documents the pre-chain rollout posture (deferred / completed) and the evidentiary substitute (legacy WORM logs). **The Wave-6 fourth erratum (§12 change log) explicitly names the pre-chain era retention gap as an institution-side legacy-log dependency, NOT a spec concern.** Not a finding of negligence — the chain was deployed as quickly as practical after the controversy. Documented for completeness; framed by §1.2 + §10.19 + §12.

Lin: "PDPC will accept that. They have asked equivalent questions in past letters and accepted equivalent answers. The §10.19 + §1.2 + §12 framing gives the answer a spec citation rather than an institution-side claim."

Park: "FSS will accept that too. They know when we deployed."

Dawn: "Good."

She moves to the close.

---

## 🔍 4:30 PM — Seoul + Taipei, Bridge Open (continued)

Two more confirmations to land before the debrief.

The seventh Confirmation — Dawn pulls it from her notes.

> ✅ **Confirmation #9 — K-ISMS and CNS 27001 alignment across the two HSMs under §10.5.**
> The Seoul HSM is hosted in a KISA-certified data center in Sangam-dong, in compliance with Korea's K-ISMS (Korea Information Security Management System) certification. The Taipei HSM is hosted in a Chunghwa Telecom data center in compliance with Taiwan's CNS 27001 (Taiwanese localization of ISO 27001). Both certifications were re-validated in the past twelve months. Both certifications were inspected by their respective regulators in the past twenty-four months. Per §10.5 HSM custody, both HSMs operate under FIPS 140-2 Level 3 or higher with separation-of-duties controls; per §10.17, partition ceremonies are chain-coupled with dual-control attestation; per §10.18, both jurisdictions' operational runbooks cross-reference §10.5 + §10.17 in their procedural sections.

Dawn also calls out the verifier discipline. "Per §10.26, Sun-Won runs the reference verifier under the pinned `v1.0b-verifier` release. CC8.1 names the implementation, the version, and the verification key per §10.26's three-name citation discipline. Reproducible-build evidence is in the binder; Cosign signature is in the binder; SBOM is in the binder. Examiners across both jurisdictions can rebuild the binary from source and confirm against the published artifact. That is the §10.26 conformance bar."

Dawn adds one more for the record. "And the entity-succession discipline per §10.24 — Sun-Won Holdings has no acquisition or divestiture activity in scope this period, but the CC8.1 control description names §10.24's procedure for any future legal-entity transition. The `chain.entity_succession` operational event is wired in the SDK; the dual-signature requirement per §10.17 schema with `entity_affiliation` per Round-17 M&A-P1 is documented; the from-tenant-id / to-tenant-id discipline is named for the rename case. Acquirer-side IT due-diligence per §13's stakeholder-navigation entry has a binding spec section to cite if Sun-Won is ever acquired or itself acquires another entity. And §10.25 run-resume discipline — the SDK uses local-persistence sidecar (SQLite) with file-locked writer-per-run as the steady-state path; the ledger chain-tail endpoint is the rejoin path for DR. The single-writer-per-run rule under §10.25 prevents any two SDK processes from racing for the same `(tenant_id, run_id)`."

> ✅ **Confirmation #10 — Reference verifier distribution and citation per §10.26.**
> Sun-Won's CC8.1 cites the reference verifier with three names per §10.26: the implementation (the spec's reference verifier), the version (`v1.0b-verifier`, the spec-pinned version per §11), and the verification key (the institution's accepted Cosign key fingerprint). Reproducible-build evidence, Cosign signature, SHA-256/SHA-512 manifest, and CycloneDX SBOM are retained in the institution's binder per §10.13 evidentiary artifacts list.

The team is in a good place. Dawn looks at her notes and counts: ten Confirmations plus the pre-chain era confirmation-by-spec, two Findings against v1.0b spec (cross-border attribute and chained classifier_output) and one tracking finding (Finding-003 — chatbot reconciliation, closure path tied to Finding-002). The chain holds within both jurisdictions. The cross-border boundary holds procedurally today and will hold cryptographically after Q3 once the spec-mandated remediations land.

She closes her notebook.

---

## 🌆 5:30 PM — Seoul + Taipei, Joint Debrief on Video Bridge

Full team on. Park-CCO and Lin-Director both present. Dawn runs the per-regulator finding table.

### PIPA (Korea Personal Information Protection Commission)

| Item | Spec § | Finding |
|---|---|---|
| Per-tenant key isolation (KR side) | §4.1, §10.15 Pattern B | ✅ PASS — distinct HSM, IKM, seal cadence; per-tenant HKDF binding |
| Redaction discipline at SDK boundary | §10.22 | ✅ PASS — pre-MAC redaction; `audit.redaction.*` family emitted per Round-17 CFPB-P2 close |
| BNPL training-data lineage + retention floor | §1.2, §10.20, §10.19 | ✅ PASS — out of chain scope per §1.2; manifest hash-anchored via §10.19; 24-month retention above §10.20 floor |
| Cross-border transfer basis (PIPA Section 28) | §4.4 (Wave-6 fourth erratum) | ⚠️ Finding-001 — non-conformance against v1.0b §4.4; remediation Q3; this story drove the §4.4 amendment |
| Pre-chain era lookback (celebrity controversy) | §1.2, §10.19, §12 | ✅ Institution-side legacy-log dependency per §12 erratum; not a spec concern |
| K-ISMS alignment | §10.5, §10.17 | ✅ PASS — Sangam-dong KISA-certified data center; partition ceremonies chain-coupled |

### PDPA (Taiwan Personal Data Protection Commission)

| Item | Spec § | Finding |
|---|---|---|
| Per-tenant key isolation (TW side) | §4.1, §10.15 Pattern B | ✅ PASS — distinct HSM, IKM, seal cadence |
| Excluded features at ingest | §4.1, §5, §10.22 | ✅ PASS — facial features and voice excluded structurally; integrity-bound by per-event MAC over canonical bytes |
| Article 8 explicit-consent for chatbot tenants | §4.1, §10.15 Pattern B | ✅ PASS — per-language tenant separation, consent captured at chatbot opt-in |
| Cross-border transfer basis (Article 8) | §4.4 (Wave-6 fourth erratum) | ⚠️ Finding-001 — same finding as PIPA above; same v1.0b §4.4 remediation |
| Language-detection routing rationale | §4.4.1 (Wave-6 fourth erratum) | ⚠️ Finding-002 — non-conformance against v1.0b §4.4.1; this story drove the §4.4.1 sixth-event-type addition |
| Chatbot reconciliation 3/5 routing rationale | §4.4.1 | ⚠️ Finding-003 — tracking; closure path is Finding-002 remediation |
| CNS 27001 alignment | §10.5, §10.17 | ✅ PASS — Chunghwa Telecom data center; partition ceremonies chain-coupled |

### FSS (Korea Financial Supervisory Service — BNPL)

| Item | Spec § | Finding |
|---|---|---|
| BNPL credit-scoring chain | §4.1, §4.4, §4.4.5, §7 | ✅ PASS — strict-mode verifier, 18,442 entries, 365 days, JCS self-test PASS, sign_payload v1.0b 12-line form |
| Reviewer override capture (adverse-action shape) | §10.11.1 (analog) | ✅ PASS — distinct field with reviewer_id; structured-reasons code |
| Operator identity binding (PASS-IT) | §10.5, §10.2 | ✅ PASS — break-glass discipline enforced; chained as §10.2 operational event |
| Training-data lineage manifest + retention | §1.2, §10.19, §10.20 | ✅ PASS — quarterly attestation; pipeline-enforced consent; 24-month retention floor |
| Reconciliation (5 of 5) | §7, §4.4, §4.4.5 | ✅ PASS — all five decisions traced end-to-end |
| HSM partition ceremony attestation | §10.17 | ✅ PASS — `chain.partition_ceremony_attended` events with `entity_affiliation` per Round-17 M&A-P1 |

### FSC (Taiwan Financial Supervisory Commission — listed-subsidiary disclosure)

| Item | Spec § | Finding |
|---|---|---|
| Per-jurisdiction tenant isolation | §4.1, §10.15 Pattern B | ✅ PASS — no co-mingling of key material between parent and subsidiary |
| Daily seal cadence (Taipei HSM) | §4.2.1, §4.3, §10.4 | ✅ PASS — 464 consecutive daily seals, no gaps, NTP discipline |
| Cross-jurisdiction inventory traceability | §4.4, §10.15 invariant 5 | ✅ PASS — input feeds hash-anchored both directions; replication-completion event freshness per Round-17 Wave-6 third erratum |
| Cross-border transfer basis disclosure | §4.4 (Wave-6 fourth erratum) | ⚠️ Finding-001 — same finding; same v1.0b §4.4 remediation; FSC will appreciate the upgrade |
| CNS 27001 alignment | §10.5, §10.17 | ✅ PASS |
| Reference verifier citation | §10.26 | ✅ PASS — three-name citation in CC8.1 |

---

Dawn closes.

"Three regulators, one chain. Today, the chain answers each regulator inside its jurisdictional lens cleanly within the v1.0a posture you operated through 2025. Two findings against the v1.0b spec text — the cross-border attribute family per §4.4 (Finding-001) and the chained classifier_output event per §4.4.1 (Finding-002) — both of which were lifted to the spec body in the Wave-6 fourth erratum precisely because Sun-Won's posture surfaced them. That last point matters for the regulator letters: the spec moved because of work like this. Sun-Won's remediation is to a normative posture the institution helped shape. Park, Lin — Sun-Won has done this well. The discipline shows. The HSM streams have not skipped a day in sixteen months. The PASS-IT integration is clean. The post-controversy redesign — moving feature exclusion to ingest — is structurally provable per §4.1 + §5 + §10.22, not just policy-provable.

The §10.18 runbook cross-referencing is in place across both jurisdictions; the Korean and Mandarin operational runbooks are cross-referenced from the English CC8.1 per §10.17's last paragraph; the §10.26 reference-verifier citation is a clean three-name citation. The §1.4 compositional-security argument composes in both jurisdictions because the chain's three layers — per-event HMAC under §4.1, daily Merkle seal under §4.2, HSM-rooted signature under §4.3 — are each independently strong and the §10.5 + §10.17 custody discipline closes the obvious online attack surface.

We will partition the deliverable by regulator audience. PIPA letter, PDPA letter, FSS letter, FSC letter. One chain. Four readers. Same answers, framed for each."

Park: "Thank you. Six weeks of preparation. Worth it."

Lin: "Thank you. We will see you in Q3 for the follow-up — we will bring the v1.0b §4.4 + §4.4.1 remediation evidence."

Dawn looks at her team — four faces in Seoul, four faces in Taipei, eight people who have spent a long day inside two binders and one chain. "We will draft the letters this week. Park, Lin — you will see the partitioned drafts before they go to any regulator. Standard process. Comments back inside ten business days."

Park: "Standard. Thank you."

Lin: "Standard. Thank you."

Tom adds — quietly, half to the room and half to himself — "Eight engagements in. The shape of this one is going to stay with me. Three regulators, one chain, two HSMs, one strait, two §4.4 amendments that this story drove. The cross-border boundary held procedurally today and will hold cryptographically after Q3 once the §4.4 + §4.4.1 remediation lands. That is a clean answer to a hard question — and the spec section numbers anchor the answer."

Dawn nods. *That is the right summary line. Tom heard it the way I heard it.*

The bridge stays open another few minutes for handshakes — Korean bows on Park's side, slight Taiwanese inclines on Lin's, easy nods from the team. The recordings stop. The cameras stay on for one more minute. Then off.

The Seoul team gathers their notebooks. Dawn looks out the conference-room window — the Han River is dark blue under the late-afternoon sun, and the Sangam-dong towers are catching gold off their west faces. She thinks about the Israeli engagement last week, where the test was nation-state segregation. She thinks about Atrio, where the test was multi-tenant key isolation across a banking platform. Today the test was the cross-border data-flow basis under three regulators.

*Different test,* she thinks. *Same chain. The chain held. It always does, when the operator has done the work and the spec answers the question. Sun-Won did the work. The spec answered the question — twice, in the Wave-6 fourth erratum, with the §4.4 cross-border-transfer family and the §4.4.1 classifier_output event. That is what a maturing spec looks like.*

---

---

## 📝 Spec-Citation Index (lead-auditor working paper)

Dawn's working paper for this engagement carries a per-section spec-citation index — the spec sections each finding or confirmation lands against. The index is part of the deliverable so every regulator receiving a partition of the report can walk from a finding to the spec text by section number without inferring the mapping.

| Spec § | Use in Sun-Won Engagement |
|---|---|
| §1.1 | Daubert four-factor grounding for litigation posture (Tom's note on the celebrity-controversy lookback) |
| §1.2 | Epistemic scope — the chain proves what the AI said and that the record was not tampered with; training-data lineage and pre-chain era are out of scope by spec |
| §1.3 | Security definitions — EUF-CMA / second-preimage / EUF-CMA composition |
| §1.4 | Compositional security — three independent layers |
| §3 | Definitions; tenant_id character class (`^[A-Za-z0-9_.\-]{1,255}$`); chain_kind enumeration |
| §3.1 | Legacy tenant identifier handling — Sun-Won's three legacy CRM identifiers under Pattern 1 / Pattern 2 |
| §4.1 | Per-event MAC and HKDF binding |
| §4.1.2 | FFIEC-conformance posture (`ffiec.chain.posture = ffiec`) |
| §4.2 | Daily Merkle seal; empty-day seal continuity |
| §4.2.1 | Cadence (daily for Sun-Won) |
| §4.3 | HSM-rooted root signature; v1.0b 12-line `sign_payload` form |
| §4.4 | Chain envelope, `gen_ai.{request,response}.model` MUST, `audit.cross_border_transfer.*` family (Wave-6 fourth erratum), parent_run_id / parent_seq |
| §4.4.1 | Routing schema; `audit.routing.classifier_output` sixth event type (Wave-6 fourth erratum); required-pairing rule |
| §4.4.2 | Deployment-intent capture |
| §4.4.3 | OTLP transport identification (Resource attributes + HTTP headers) |
| §4.4.4 | Severity for chain-of-custody traffic (collector pass-through; receiver stamping `9..20`) |
| §4.4.5 | Underwriting features family (Round-17 NAIC-P1) |
| §4.4.6 | SaaS-edge connector source attribution (`audit.connector_source.*`); stable-`run_id` discipline |
| §5 | Wire format; canonical-form exclusion rule |
| §5.2 | Best-evidence posture under FRE 1001-1004 |
| §6 | Storage; chain-stamp preservation |
| §7 | Verifier procedure; JCS self-test pre-flight; normative reason strings |
| §10.1 | Key-fingerprint reconciliation |
| §10.2 | Operational events |
| §10.3 | Append-only enforcement |
| §10.4 | Time synchronization (NTP) |
| §10.5 | HSM custody |
| §10.6 + §10.6.1 | IKM minimum length (32 bytes); IKM generation requirements |
| §10.9 | IKM registry retention |
| §10.11 + §10.11.1 | Adverse-action notice translation (analog applied to FSS BNPL declination); reasons schema |
| §10.13 | Evidentiary artifacts retention |
| §10.15 | Multi-region resilience — Pattern A (cross-inventory) + Pattern B (per-region tenants) |
| §10.16 | SaaS-edge capture connectors (four quantified bounds); Salesforce CDC mirror for TW CRM |
| §10.17 | HSM partition ceremony attestation; `entity_affiliation` per Round-17 M&A-P1 |
| §10.18 | CC8.1 and runbook cross-referencing |
| §10.19 | Chain-coverage boundary documentation; `audit.external_artifact.*`; pre-chain era is institution-side legacy-log dependency |
| §10.20 | Training-data retention vs deployment-window discipline (Sun-Won 24 months) |
| §10.21 | Cross-vendor model-handover schema (Korean chatbot model from Seoul consultancy); Round-17 M&A-G2 contract binding |
| §10.22 | Redaction discipline (pre-MAC at SDK boundary) |
| §10.24 | Entity succession (no current activity; CC8.1 names the procedure) |
| §10.25 | Run resume and chain-tail acquisition; SQLite sidecar + ledger rejoin |
| §10.26 | Reference verifier distribution; three-name CC8.1 citation |
| §11 | References — pinned `v1.0b-verifier` |
| §12 | Change log — Wave-6 fourth erratum lifted §4.4 cross-border-transfer family and §4.4.1 sixth-event-type to spec body |
| §13 | Stakeholder navigation — acquirer-side IT due-diligence entry |

**Spec amendments this engagement drove (§12 Wave-6 fourth erratum):** §4.4 `audit.cross_border_transfer.*` attribute family lifted to spec body; §4.4.1 sixth event type `audit.routing.classifier_output` added with six new attributes. Sun-Won's posture surfaced both gaps under v1.0a; the spec body now normates both. Sun-Won's Q3 remediation is to v1.0b conformance, not advisory upgrade.

---

## 🧾 Final Assessment Theme

> **"Three regulators reading the same chain through three different lenses. Today, the chain plus the contract answers each lens within v1.0a posture. After Q3, the verifier output answers all three on its own under v1.0b §4.4 + §4.4.1 — the contract is referenced inside the chain entry by hash, the routing rationale is integrity-bound through the chained classifier_output event, and the spec's normative text matches the work item. That is the upgrade Sun-Won is paying for — and it is the cleanest cross-border story we have audited yet, because the spec moved to meet the posture rather than the posture having to bend to a fixed spec."**
>
> — Dawn, lead auditor, Sun-Won Cosmetics Group, April 9, 2026

---

*End of Story 09.*
