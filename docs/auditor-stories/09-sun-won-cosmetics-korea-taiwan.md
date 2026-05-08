# 🧾 Auditor Story 09 — Sun-Won Cosmetics Group (Korea + Taiwan)

> **Engagement.** Sun-Won Holdings (KOSPI: 003410) and Sun-Won Taiwan Co. Ltd. (Taipei Exchange: 5891) — a coordinated annual review covering PIPA Section 28 cross-border transfers, FSS supervisory review of the BNPL consumer-finance arm, Taiwan FSC review of the listed subsidiary, PDPA Article 8 explicit-consent oversight, and a CPRA / GDPR sweep against the e-commerce platform.
>
> **Subject system.** TesseraSeal v1.0a, in production for sixteen months across four AI use cases — customer recommendation, inventory forecasting, multilingual chatbot, BNPL credit-scoring. Eight tenants, plus a ninth cross-jurisdiction tenant for inventory.
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

Karen has done eight of these now. She knows the cadence. The cross-border boundary is the part she has not seen tested under three regulators at once.

---

## Audit Team

### Seoul (Sun-Won HQ — Sangam-dong)

- **Karen** — Lead Auditor. Anchors the engagement from Seoul. Final sign-off on all four regulator-partitioned findings.
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

The Sangam-dong tower has a glass atrium that catches the morning light off the Han River. Karen, Raj, Diana, and Tom badge in at the security desk. The receptionist switches to English the moment Karen says her name and hands them four visitor lanyards in Sun-Won's house pink.

Park Hye-jin meets them at the elevator bank. She is in a charcoal suit, no jewelry, carrying a leather portfolio that Karen recognizes as the FSS examiner-issue from about 2011.

"Welcome back," Park says. "Lin-Director's team should already be with your Taipei four. We have the bridge open in the executive conference room on twenty-eight."

"Thank you for the early start," Karen says. "How is the binder?"

"Six weeks of preparation. The FSS supervisory letters from the last cycle are tabbed. The PIPC quarterly attestations are tabbed. The Taiwan FSC corresponding letters Lin's team holds — those are duplicated on her side. I have one set here as well."

*She is being thorough on purpose,* Karen thinks. *She knows the difference between a routine year and a coordinated three-regulator review. She is not going to make us ask twice.*

They ride up. Park hands Karen a printed agenda in Korean and English, side-by-side columns. The eight tenants are listed with their tenant_ids and their primary regulator mapping. The ninth — `sunwon-cross-inventory` — is at the bottom, in italics.

"That's the one we should talk about," Park says. "I marked it because you will see it before lunch."

"Thank you," Karen says. "We will."

Park turns to Raj. "We pre-staged the BNPL chain on a read-only mirror. Your verifier credentials are in your packet — PASS-IT-bound, scoped to read-only, expire at six tonight. Diana, your IAM packet is the same shape, scoped to the IdP audit role. Tom, the compliance binder is on the second cart in the room. Photocopying is fine. Photography is not."

"Understood," Tom says.

"And the room has tea, coffee, and a coffee machine that none of us know how to operate," Park says. "We will figure it out together."

---

In the elevator, Karen does her opening monologue for the Seoul four — quietly, the way she always does, while they ride up.

> "Two countries, four use cases, three regulators, eight tenants, one chain. Park-CCO has been preparing for this for six weeks. We have done eight of these now — Northbridge, Mercator, Stelvio, Atrio, Helmstad, Pacific Crescent, Olmstead, and we just came from a 23-tenant Israeli AI shop where the test was nation-state segregation. Today's test is the cross-border data-flow basis. The chain is the same chain. The questions are not."

Raj nods. Diana is already pulling up the PASS-IT documentation on her tablet. Tom is making sure his recorder is on.

"The chain holds within a jurisdiction," Karen says. "The cross-border evidence has to hold to two different regulators reading the same chain."

"It never is," Tom says, half a second before Karen can.

She lets him have it.

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

Elena writes that down. *Lin is going to want to see every sentence before it goes anywhere.* "We will phrase findings neutrally. You and Karen take the partitioning."

"Thank you," Lin says.

---

## 🧩 9:15 AM — Seoul, Conference Room 28-A

Raj opens the BNPL credit-scoring chain. He has been thinking about it on the plane — the BNPL tenant is the highest-stakes one in the portfolio because the FSS examiner standard for AI-driven credit decisions is the strictest standard at the table today.

He pulls up the tenant config first.

```
tenant_id:        sunwon-kr-bnpl
hsm:              seoul-sangam-hsm-01 (KISA-certified)
ikm:              kr-bnpl-ikm-2026-q2
seal_cadence:     daily, 23:59:59 KST
spec_version:     1.0a
attribute_set:    [applicant_id_hash, model_id, model_version,
                   features_hash, score, decision, reviewer_override]
```

Park watches over Raj's shoulder. "The reviewer_override field — that's the one the FSS examiner asked about in the last letter."

"I see," Raj says. "What did the letter ask?"

"It asked whether the override is captured deterministically — meaning, can the chain show that a human reviewer overrode the model, distinct from the model's own conditional output."

Raj nods. He pulls a sample entry from April 4 — five days ago, well within the daily-seal window.

```
entry_id:               kr-bnpl-2026-04-04-00041823
applicant_id_hash:      sha256:e7c2...9f1a
model_id:               sunwon-bnpl-v3
model_version:          3.4.2
features_hash:          sha256:b81f...0c7e
score:                  712
decision:               conditional
reviewer_override:      none
sealed_at:              2026-04-04T23:59:59+09:00
seal_signature:         ed25519:7a2f...4b8c
```

"Conditional with no override," Raj says. "So the model itself returned conditional."

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
sealed_at:              2026-04-04T23:59:59+09:00
seal_signature:         ed25519:7a2f...4b8c
```

"There it is," Raj says. "Score 689 would have been declined under the auto-cutoff at 700. Reviewer 014 overrode to approved. The chain shows the model's score and the reviewer's decision separately. FSS-grade."

Park exhales just slightly. *That was the question she was holding.*

"Run the verifier," Karen says from across the table.

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
[verify] reviewer_override field: present in 1.7% of entries (expected range)
[verify] PASS
```

> ✅ **Confirmation #1 — BNPL credit-scoring (FSS-grade).**
> Per-applicant entries capture applicant_id_hash, model_id, model_version, features_hash, score, decision, and reviewer_override. The override is captured as a distinct field with reviewer_id, separable from the model's own conditional output. Verifier strict-mode PASS over 365 days, 18,442 entries. Spec section 4.6 (BNPL/consumer-finance attribute set) honored.

"That's the one the FSS letter asked about," Park says. "We can put that in front of an examiner without commentary."

"Yes," Karen says. "That entry stands on its own."

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

"Feature-change manifest is wired," Raj says. "Every model_version bump that changes the feature set lands a manifest entry. Verifier reconciles it. That answers the implicit question — did anything quietly leave the model? No."

Park: "That was the 2026-Q1 letter. Already answered. I just like seeing it answered twice."

---

## 🧩 9:15 AM — Taipei (8:15 AM CST), Sun-Won Taiwan Compliance Suite

Mike opens the recommendation-engine chain for the Taiwan tenant.

```
tenant_id:        sunwon-tw-rec
hsm:              taipei-chunghwa-hsm-02 (CNS 27001 aligned)
ikm:              tw-rec-ikm-2026-q2
seal_cadence:     daily, 23:59:59 CST
spec_version:     1.0a
attribute_set:    [user_id_hash, model_id, model_version,
                   features_hash, recommended_skus, served_at]
```

"Different HSM," Lin says. "Different IKM. Different cadence anchor — Taipei time, not Seoul time."

"That's the right shape," Mike says. "Per-jurisdiction tenant means per-jurisdiction key material. PIPA Section 28 and PDPA Article 8 both want that — the data crossing the strait should not be sealing under the same IKM as data that stayed home."

Lin nods slowly. "FSC will look for that explicitly. They have asked, in past letters, about co-mingling of key material between the parent and the subsidiary."

Mike pulls a sample.

```
entry_id:               tw-rec-2026-04-09-00982341
user_id_hash:           sha256:4a8c...11de
model_id:               sunwon-rec-tw-v7
model_version:          7.2.0
features_hash:          sha256:c2f1...9087
recommended_skus:       [SK-77241, SK-91123, SK-44290]
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
[verify] PASS
```

"Clean," Mike says.

> ✅ **Confirmation #2 — Per-jurisdiction tenant isolation.**
> The KR and TW recommendation tenants use distinct HSMs, distinct IKMs, and distinct seal-cadence anchors. No co-mingling of key material across the strait. PIPA Section 28 and PDPA Article 8 both satisfied at the chain layer. Spec section 3.2 (per-tenant key isolation) honored.

Lin makes a note in her own portfolio. "FSC will accept that. That is a clean answer."

---

## 🧠 10:00 AM — Seoul, Conference Room 28-A

Diana takes over for the IAM walk. Sun-Won uses **PASS-IT** — a Korean government-administered authentication provider that issues mobile-PKI tokens for KISA-certified services. PASS-IT is integrated into Sun-Won's SSO via SAML 2.0 with a Korean-resident IdP cluster.

"PASS-IT is the right choice for this jurisdiction," Diana says. "The Korean residency requirement in K-ISMS makes it the obvious pick. Is it the only auth path into the chain operator console?"

"It is the only auth path into the production console," Park says. "Break-glass exists. Break-glass requires PASS-IT plus a Yubikey plus a recorded business justification."

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

"One break-glass in the last quarter," Diana says. "Read-only. Forty-two minutes. Justified. Reviewed."

"That is the only one," Park confirms.

Diana checks the per-tenant scoping. PASS-IT identities are bound to tenant scopes via the SSO claim set. Operator `kr-ops-009` has scope `sunwon-kr-rec` and `sunwon-kr-bnpl` only — not `sunwon-tw-rec`, not `sunwon-cross-inventory`, not the chatbot tenants.

"Per-tenant scoping is enforced at the IdP," Diana says. "The chain operator console will not let `kr-ops-009` even see the TW tenants. That is correct."

> ✅ **Confirmation #3 — IAM via PASS-IT, per-tenant scoping enforced at the IdP.**
> Korean SSO uses PASS-IT mobile-PKI as the only path into the production console. Break-glass requires PASS-IT + Yubikey + recorded justification + post-event CCO review. Per-tenant scoping is enforced at the IdP claim set, not at application code. K-ISMS access-control requirements satisfied. Spec section 5.4 (operator identity binding) honored.

"FSS will read this and stop reading," Park says. "That is exactly what they asked for in the 2024 letter."

---

## 🧠 10:00 AM — Taipei (9:00 AM CST), Compliance Suite

Elena walks the CRM side of the recommendation engine — the input side, where customer records become feature vectors.

The Taiwan-side CRM is called **Sun-Won-Connect-TW**, a Salesforce-derived deployment with localized fields for Taiwanese ID-card prefixes and the National Health Insurance card identifier (which is *not* used as a feature, but is in the CRM for shipping).

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

"The features_hash in the chain entry is the hash of the feature vector," Elena says. "Not the source data."

"Correct."

"And the excluded features — facial features from photo — are excluded at ingest, not at model time."

"Correct," Lin says. "After the controversy, we moved the exclusion to ingest. The model layer cannot see those fields because they never enter the feature pipeline."

Elena writes this down. *That is the thing that matters about the post-controversy redesign,* she thinks. *Exclusion at the model layer is a promise. Exclusion at ingest is structural.*

> ✅ **Confirmation #4 — Excluded features removed at ingest, not at model layer.**
> Facial features from photos and voice features are excluded at the ingest stage, before the feature pipeline. The model layer cannot observe them because they never enter the pipeline. The features_hash in each chain entry is over the actual feature vector — making the absence of excluded features structurally provable, not policy-provable. Spec section 4.2 (feature-set attestation) honored.

---

## 🔐 11:00 AM — Seoul, Conference Room 28-A

Tom is reviewing Park's binder of prior FSS supervisory letters. There are nine of them across the past three years. He matches each one against the current chain configuration to see whether the chain has answered each prior question structurally.

The 2024-Q2 letter asked about reviewer-override capture — answered (Confirmation #1 above).

The 2024-Q4 letter asked about model versioning — answered, model_version is in every entry.

The 2025-Q1 letter asked about cross-tenant operator access — answered (Confirmation #3 above).

The 2025-Q3 letter asked something Tom has to read twice.

> "The supervised entity should demonstrate that the BNPL model's training data does not include retail-side customer behavior absent explicit consent at the training-data collection point."

Tom looks up. "Park, this letter is about training-data lineage, not inference. Does the chain cover training?"

"The chain covers inference," Park says. "Training-data lineage is in a separate manifest. We attest to it quarterly."

"Can I see the manifest?"

Park pulls it. The training manifest lists the data sources, their consent basis, and the cutoff timestamps. The BNPL model's training set draws from BNPL applications only — not from retail purchase history, unless the customer signed the cross-use consent at BNPL application time.

"Cross-use consent rate?" Tom asks.

"Forty-one percent. The model is trained on the consenting applicants only. The non-consenting applicants are not in the training set. That is enforced by the training pipeline, not by the model."

Tom nods. *That is the right answer for the 2025-Q3 letter. The chain does not need to cover training. The manifest covers training. They are separable concerns.*

He notes: "Training-data lineage out of chain scope, in manifest, attested quarterly, enforced at pipeline. FSS 2025-Q3 letter answered."

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
[verify] PASS
```

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

> ✅ **Confirmation #5 — Daily seals on both jurisdiction HSMs, no gaps.**
> Seoul HSM (Sangam-dong, KISA-certified): 365 + 99 consecutive daily seals, zero gaps, zero late landings. Taipei HSM (Chunghwa Telecom data center, CNS 27001 aligned): 365 + 99 consecutive daily seals, zero gaps, zero late landings. Daily seal cadence operationally sound across both jurisdictions. Spec section 6.1 (daily seal cadence) honored.

Luis adds, half to himself: "Two HSMs, two operations teams, two on-call rotations. Zero gaps. That is not luck."

Lin smiles. "It is not. Park-CCO and our Taipei ops director have a running bet about who is going to break the streak first. Neither has."

---

## 🔐 11:30 AM — Seoul + Taipei, Bridge Open

The two halves of the team join the video bridge for the first time. Park and Lin are visible side-by-side on each other's screens, in their own conference rooms.

Karen opens. "Before lunch I want to surface one item, so we can talk about it over food. The inventory-forecasting tenant. Park, Lin — you both know the shape. Can someone walk us through it?"

Park nods to Lin. Lin nods back.

"It is the only tenant that crosses the strait," Lin says. "Single tenant: `sunwon-cross-inventory`. Inventory data from Korean stores and Taiwanese stores both feed in. The model forecasts SKU-level demand at the regional warehouse layer. The model has to see both jurisdictions because the inventory rebalances between them — particularly at quarter-end and around lunar holidays."

"Cross-border transfer basis?" Karen asks.

"Contract," Park says. "There is a documented intra-group data transfer agreement between Sun-Won Holdings and Sun-Won Taiwan, registered with PIPC and acknowledged by PDPC. The contract identifies inventory data as a category, identifies the model as a recipient, and identifies the lawful basis as legitimate business interest with appropriate safeguards."

"And the chain entries?" Karen asks.

Pause.

"The chain entries do not carry the cross-border transfer basis as an attribute," Park says. "The contract carries it. The compliance binder carries it. The chain entry shows that an inventory data point from store KR-Seoul-014 was used in a forecast — and that the forecast was generated on the cross-tenant. The cross-border basis is not stamped into the entry."

"That is what I thought," Karen says. *That is the thing.*

"Lunch?" Lin asks.

"Lunch."

---

## 🧪 12:00 PM — Seoul + Taipei, Working Lunch via Video Bridge

Korean lunch in Seoul: bibimbap, mandu, pickled radish. Taiwanese lunch in Taipei: lu rou fan and pickled mustard greens. Both teams eat with their cameras on. Park and Lin are on the call together for the first time today.

Karen does not let the lunch slide into chitchat. She takes the inventory tenant head-on.

"Park, Lin. The inventory tenant is the one that needs the conversation. The chain works. The contract works. The two pieces of evidence sit in different binders. If FSS, PIPC, FSC, and PDPC all asked the same question on the same day — show me the cross-border basis for this one inventory data point — could you give all four the same answer?"

Park considers. "Today, the answer is: the chain entry plus the contract reference. The chain proves what the inventory model saw. The contract proves the consent basis. They have to be read together."

"Lin?"

"FSC and PDPC will accept that," Lin says, "but they will note in their letter that the link between the two is procedural rather than cryptographic. PDPC has been asking for cryptographic linkage in cross-border-flow attestations since 2024. They have not made it a normative requirement yet. They will, eventually."

Karen nods. "So the recommendation is: add a `cross_border_transfer_basis` attribute to the inventory tenant's chain entries. The attribute is a hash anchor — it points to the version of the contract that authorized the transfer at the time the entry was sealed. The contract changes; the hash changes; the entry shows which contract version was in force."

Mike, on the Taipei side, leans into his camera. "That is a small change. The inventory tenant's attribute set is six fields. Adding a seventh is a config change, not a code change. The contract-versioning side is the work — Sun-Won's compliance team has to publish the contract as a versioned, hash-anchored document."

"Six weeks of work on the legal side," Lin says. "Two weeks on the chain side. Maybe."

"Advisory recommendation," Karen says. "Not a normative spec gap. The current setup answers the regulators today. The recommendation upgrades the answer from procedural to cryptographic, which positions Sun-Won for the PDPC requirement when it lands."

> ⚠️ **Surprise / Partial #1 — Cross-border transfer basis not stamped into inventory chain entries.**
> The `sunwon-cross-inventory` tenant aggregates inventory data from both Korea and Taiwan. The lawful basis for cross-border transfer is documented in an intra-group data transfer agreement registered with PIPC and PDPC. The contract is sound. However, the chain entries themselves do not carry a `cross_border_transfer_basis` attribute pointing to a versioned, hash-anchored contract reference. The two pieces of evidence — chain entry and contract — are linked procedurally rather than cryptographically. **Advisory recommendation:** add a `cross_border_transfer_basis` attribute to the inventory tenant's attribute set, populated with a hash anchor pointing to the contract version in force at seal time. Not a normative spec gap. Positions Sun-Won for the PDPC's anticipated cryptographic-linkage requirement.

Park writes the recommendation into her binder. "We will plan it for Q3. The legal versioning is the long pole."

"Agreed," Lin says.

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
prompt_hash:            sha256:b193...7c0f
response_hash:          sha256:c4d2...8e9b
served_at:              2026-04-09T11:42:18+08:00
seal_signature:         ed25519:b9d4...7e23
```

"The chain shows which model served," Mike says. "It does not show how the routing decision was made. The detector logs are separate."

"Where do the detector logs live?" Lin asks.

"Separate microservice, separate log system. Ninety-day retention, by default."

Mike runs the verifier on the Mandarin tenant.

```
$ herald-verify --tenant=sunwon-chatbot-zh \
                --service=chatbot \
                --date=2026-04-09

[verify] tenant=sunwon-chatbot-zh
[verify] hsm=taipei-chunghwa-hsm-02
[verify] entries scanned: 9,217
[verify] seal coverage: 365/365 days
[verify] PASS
```

"Clean," Mike says. "The chain shows the Mandarin model served 9,217 interactions yesterday. It does not show why those particular interactions were routed to the Mandarin model rather than the Korean or English ones."

> ✅ **Confirmation #6 — Chatbot per-language tenant separation.**
> Three chatbot tenants — Korean, Mandarin, English — each with its own model, model_version, IKM, and seal stream. Per-language separation is structural at the tenant layer. PIPA and PDPA both want the model populations separated; both are satisfied. Spec section 3.2 (per-tenant key isolation) honored.

> ⚠️ **Surprise / Partial #2 — Language-detection routing decision is not chained.**
> The chatbot's first step is a language-detection microservice that picks which model serves the user. The chain entry records which model served — not which classifier output drove the routing. Reconstructing "why was this user served by the Mandarin model rather than the Korean model" requires the language-detector microservice's logs, which live in a separate log system on a 90-day retention. Beyond 90 days, the routing rationale is not recoverable from the chain alone. **Advisory recommendation:** either (a) chain the detector's classifier output as a pre-routing entry, or (b) extend the detector's log retention to match the chain retention. Spec is silent on routing-rationale capture; this is a defensibility-of-evidence finding, not a normative gap.

Lin looks at the recommendation. "I prefer (a). Chaining the detector's output makes the answer self-contained. Log retention extensions get rolled back when finance reviews them."

Mike nods. "Same. Chaining it is also cheap — the detector emits one classification per interaction; the chain entry is small."

Lin notes it.

---

## 🧬 2:00 PM — Taipei (1:00 PM CST), Compliance Suite

Chen takes the inventory-forecasting cross-jurisdiction data flow.

The architecture:

```
Korean stores (300+) ──┐
                       ├──> Cross-jurisdiction data lake (Seoul region)
Taiwanese stores (80+) ┘                │
                                        ▼
                            Inventory forecasting model
                          (sunwon-cross-inventory tenant)
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
[verify] input feeds: 2 (kr-stores, tw-stores)
[verify] kr-stores feed anchor: sha256:8e2c...41bf  (matches manifest)
[verify] tw-stores feed anchor: sha256:3a91...0d22  (matches manifest)
[verify] entries scanned: 2,194
[verify] seal coverage: 365/365 days
[verify] PASS
```

"Both feeds anchored," Chen says. "The chain shows that a forecast generated on April 9 used the Korean feed at hash 8e2c... and the Taiwan feed at hash 3a91... — and the manifest confirms those hashes correspond to the inventory snapshots taken at midnight KST and midnight CST respectively."

"So if PDPC asks 'show me which Taiwanese inventory snapshot fed the cross-tenant on April 9'," Lin says, "the chain answers."

"Yes."

"And if PIPC asks 'show me that the Korean inventory snapshot was the only Korean data that fed the cross-tenant'," Chen says, "the chain answers — kr-stores feed anchor, no other Korean source."

"Good."

The cross-border transfer-basis attribute is the gap — already noted at lunch. The hash-anchor side of the cross-jurisdiction flow is clean.

> ✅ **Confirmation #7 — Cross-jurisdiction inventory feeds hash-anchored, both directions traceable.**
> The `sunwon-cross-inventory` tenant's input feeds — Korean stores and Taiwanese stores — are both hash-anchored at the daily snapshot point. The chain entries reference the input anchors, allowing per-forecast attribution back to the originating jurisdiction's inventory snapshot. Multi-region Pattern A (per spec section 7.3) is honored: one logical tenant, two regional input streams, hash-linked to the entries. Cross-border data-flow traceability is structurally provable at the data-lineage level. The remaining gap is the cross-border-transfer-basis attribute (Surprise #1 above).

---

## 📊 3:00 PM — Seoul, Conference Room 28-A

Reconciliation test, Seoul side. Karen has Raj pull five BNPL credit decisions from yesterday and trace each one end-to-end.

```
Sample 1: applicant_id_hash:e7c2...9f1a, decision:conditional → traced
Sample 2: applicant_id_hash:1d4a...c0f7, decision:approved (override) → traced
Sample 3: applicant_id_hash:6b81...3e92, decision:declined → traced
Sample 4: applicant_id_hash:f0a4...2c1d, decision:approved → traced
Sample 5: applicant_id_hash:9c0e...77a8, decision:conditional → traced
```

Each entry resolves to:

- The applicant's submitted features (hash-matched to the manifest).
- The model version that scored the application.
- The decision and any reviewer override.
- The seal under which it landed.
- The PASS-IT identity that recorded the override on Sample 2.

5 of 5 PASS. Tom records each one.

> ✅ **Confirmation #8 — BNPL reconciliation, 5 of 5 traced end-to-end.**
> Five BNPL credit decisions from 2026-04-08 traced from applicant submission to sealed chain entry to reviewer override (where applicable) to operator identity. All five resolve cleanly. FSS examiner-grade reconciliation discipline. Spec sections 4.6 (BNPL attribute set) + 5.4 (operator identity binding) jointly honored.

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

- The user_id_hash.
- The prompt_hash and response_hash.
- The model_id and model_version.
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

"Yes," Mike says. "3 of 5 routing rationale recoverable. Chain entries 5 of 5 PASS. The exposure is the routing side, not the chain side."

> ⚠️ **Surprise / Partial #3 — Chatbot reconciliation: 5/5 chain PASS, 3/5 routing rationale recoverable.**
> Five chatbot interactions traced. Chain entries: 5 of 5 PASS — model, model_version, prompt_hash, response_hash, seal all resolve. Routing rationale (why this user was routed to this model): 3 of 5 recoverable from the language-detector microservice logs. The other 2 fell outside the 90-day detector log retention. The exposure is the routing-rationale side, not the chain side. Reinforces the recommendation in Surprise #2 — chain the detector's classifier output as a pre-routing entry.

---

## 😬 3:45 PM — Seoul + Taipei, Bridge Open

Park and Lin both join the bridge. Karen opens.

"We have three Surprise items and eight Confirmations. Of the Surprises, only the cross-border-attribute and the language-detection-routing items have any near-term action. The pre-chain era retention gap is historical — let's discuss it last."

Park nods.

"On the cross-border attribute," Karen says, "Park and Lin both heard my recommendation at lunch. Add the attribute on the inventory tenant. Hash-anchor the contract version. Six weeks of legal work to publish the contract as a versioned hash-anchored document, two weeks of chain config work to add the attribute. Park, Lin — disagreement?"

Park: "Agreed."

Lin: "Agreed. I want to add — FSC will appreciate this in the next supervisory letter. They are watching for it."

"Good. On the language-detection-routing," Karen says, "the question is whether to chain the detector output or to extend detector log retention. Mike — you took the position that chaining is preferable."

"Chaining is preferable," Mike says. "It makes the routing rationale self-contained inside the chain. Log retention extensions get rolled back when budget reviews them. The chain is harder to roll back."

"Agreed," Park says. "Make it part of the same Q3 work as the cross-border attribute."

Lin: "Agreed. Both are essentially the same shape of work — one new attribute on an existing tenant, one new entry-type for a pre-routing event. Bundle them."

"Done."

Karen pauses. *This is where the question Park has been waiting to ask is going to come.*

It does.

"Karen," Park says. "If FSS, PIPC, and FSC ask the same question — 'demonstrate consent for cross-border transfer of an inventory data point' — can we give all three the same answer today?"

The room goes a little quieter.

Karen takes a beat.

"Today, the answer is the chain plus the contract. The chain proves what the inventory model saw. The contract proves the consent basis. You have to hand over both — and both have to be read together to construct the full answer. That is procedurally sound but evidentially compound.

After the recommendations land — the cross-border-attribute on the inventory tenant — the chain plus the verifier output answers all three regulators on its own. The contract is referenced inside the chain entry, by hash. Verifier dumps the entry, the entry shows the contract version, the contract repository serves the contract by hash, and the regulator gets a single self-consistent evidence package.

That is the upgrade you are paying for. Today: two binders, both required. After Q3: one verifier output, contract referenced by hash, single answer to three regulators."

Park writes that down word for word. *That is the language she will use in her summary memo.*

Lin: "Two binders today. One verifier output after Q3. That's the right framing for FSC."

---

## 🔍 4:30 PM — Seoul + Taipei, Bridge Open

The pre-chain era — the celebrity controversy lookback.

Karen takes it directly. "The chain was deployed sixteen months ago. The celebrity controversy was eighteen months ago. There are about four months of recommendation-engine activity from before the chain that fall inside the lookback window. We cannot verify those four months through the chain because the chain did not exist for them."

Park nods. "We acknowledge that. The legacy recommendation-engine logs cover those four months. They are append-only on a write-once-read-many storage tier — that was already best practice before the chain was deployed. They are admissible. They are not chain-grade."

"Sun-Won is being honest about that," Karen says. "I appreciate it. The audit deliverable will document the chain's effective-start date and note that pre-chain activity is verifiable through the legacy logs only."

> ⚠️ **Surprise / Partial #3a — Pre-chain era retention gap (the celebrity-controversy lookback).**
> The chain was deployed in 2024-Q4. The celebrity-controversy lookback window extends back four months prior to chain deployment (mid-2024). Pre-chain activity in those four months is verifiable through legacy recommendation-engine logs (append-only WORM storage), not through the chain. Sun-Won acknowledges this honestly. The audit deliverable documents the chain's effective-start date and notes that pre-chain activity is verifiable through legacy logs only. Not a finding of negligence — the chain was deployed as quickly as practical after the controversy. Documented for completeness.

Lin: "PDPC will accept that. They have asked equivalent questions in past letters and accepted equivalent answers."

Park: "FSS will accept that too. They know when we deployed."

Karen: "Good."

She moves to the close.

---

## 🔍 4:30 PM — Seoul + Taipei, Bridge Open (continued)

Two more confirmations to land before the debrief.

The seventh Confirmation — Karen pulls it from her notes.

> ✅ **Confirmation #9 — K-ISMS and CNS 27001 alignment across the two HSMs.**
> The Seoul HSM is hosted in a KISA-certified data center in Sangam-dong, in compliance with Korea's K-ISMS (Korea Information Security Management System) certification. The Taipei HSM is hosted in a Chunghwa Telecom data center in compliance with Taiwan's CNS 27001 (Taiwanese localization of ISO 27001). Both certifications were re-validated in the past twelve months. Both certifications were inspected by their respective regulators in the past twenty-four months. Spec section 6.4 (HSM hosting and certification) honored on both jurisdictions.

The eighth — Mike's reconciliation Mandarin tenant verifier output, already PASS, and the Korean tenant verifier run earlier on Raj's terminal. Both clean. Already covered above; bundled here as the closing confirmation.

The team is in a good place. Karen looks at her notes and counts: nine Confirmations, three Surprises (one historical, two with concrete Q3 work). The chain holds within both jurisdictions. The cross-border boundary holds procedurally today and will hold cryptographically after Q3.

She closes her notebook.

---

## 🌆 5:30 PM — Seoul + Taipei, Joint Debrief on Video Bridge

Full team on. Park-CCO and Lin-Director both present. Karen runs the per-regulator finding table.

### PIPA (Korea Personal Information Protection Commission)

| Item | Finding |
|---|---|
| Per-tenant key isolation (KR side) | ✅ PASS — distinct HSM, IKM, seal cadence |
| BNPL training-data lineage | ✅ PASS — out of chain scope, manifest covers, attested quarterly |
| Cross-border transfer basis (PIPA Section 28) | ⚠️ Partial — contract sound; chain attribute not stamped; advisory recommendation to add `cross_border_transfer_basis` |
| Pre-chain era lookback (celebrity controversy) | ⚠️ Acknowledged — verifiable through legacy WORM logs only |
| K-ISMS alignment | ✅ PASS — Sangam-dong KISA-certified data center |

### PDPA (Taiwan Personal Data Protection Commission)

| Item | Finding |
|---|---|
| Per-tenant key isolation (TW side) | ✅ PASS — distinct HSM, IKM, seal cadence |
| Excluded features at ingest | ✅ PASS — facial features and voice excluded structurally |
| Article 8 explicit-consent for chatbot tenants | ✅ PASS — per-language tenant separation, consent captured at chatbot opt-in |
| Cross-border transfer basis (Article 8) | ⚠️ Partial — same finding as PIPA above; same advisory recommendation |
| Language-detection routing rationale | ⚠️ Partial — 90-day retention exposure; advisory recommendation to chain detector output |
| CNS 27001 alignment | ✅ PASS — Chunghwa Telecom data center |

### FSS (Korea Financial Supervisory Service — BNPL)

| Item | Finding |
|---|---|
| BNPL credit-scoring chain | ✅ PASS — strict-mode verifier, 18,442 entries, 365 days |
| Reviewer override capture | ✅ PASS — distinct field with reviewer_id |
| Operator identity binding (PASS-IT) | ✅ PASS — break-glass discipline enforced |
| Training-data lineage manifest | ✅ PASS — quarterly attestation, pipeline-enforced consent |
| Reconciliation (5 of 5) | ✅ PASS — all five decisions traced end-to-end |

### FSC (Taiwan Financial Supervisory Commission — listed-subsidiary disclosure)

| Item | Finding |
|---|---|
| Per-jurisdiction tenant isolation | ✅ PASS — no co-mingling of key material between parent and subsidiary |
| Daily seal cadence (Taipei HSM) | ✅ PASS — 464 consecutive daily seals, no gaps |
| Cross-jurisdiction inventory traceability | ✅ PASS — input feeds hash-anchored both directions |
| Cross-border transfer basis disclosure | ⚠️ Partial — same finding; same advisory recommendation; FSC will appreciate the upgrade |
| CNS 27001 alignment | ✅ PASS |

---

Karen closes.

"Three regulators, one chain. Today, the chain answers each regulator inside its jurisdictional lens cleanly. The cross-border boundary is the work item — and it is concrete, scoped to Q3, and bundled with the chatbot routing-rationale upgrade. Park, Lin — Sun-Won has done this well. The discipline shows. The HSM streams have not skipped a day in sixteen months. The PASS-IT integration is clean. The post-controversy redesign — moving feature exclusion to ingest — is structurally provable, not just policy-provable.

We will partition the deliverable by regulator audience. PIPA letter, PDPA letter, FSS letter, FSC letter. One chain. Four readers. Same answers, framed for each."

Park: "Thank you. Six weeks of preparation. Worth it."

Lin: "Thank you. We will see you in Q3 for the follow-up."

Karen looks at her team — four faces in Seoul, four faces in Taipei, eight people who have spent a long day inside two binders and one chain. "We will draft the letters this week. Park, Lin — you will see the partitioned drafts before they go to any regulator. Standard process. Comments back inside ten business days."

Park: "Standard. Thank you."

Lin: "Standard. Thank you."

Tom adds — quietly, half to the room and half to himself — "Eight engagements in. The shape of this one is going to stay with me. Three regulators, one chain, two HSMs, one strait. The cross-border boundary held procedurally today and will hold cryptographically after Q3. That is a clean answer to a hard question."

Karen nods. *That is the right summary line. Tom heard it the way I heard it.*

The bridge stays open another few minutes for handshakes — Korean bows on Park's side, slight Taiwanese inclines on Lin's, easy nods from the team. The recordings stop. The cameras stay on for one more minute. Then off.

The Seoul team gathers their notebooks. Karen looks out the conference-room window — the Han River is dark blue under the late-afternoon sun, and the Sangam-dong towers are catching gold off their west faces. She thinks about the Israeli engagement last week, where the test was nation-state segregation. She thinks about Atrio, where the test was multi-tenant key isolation across a banking platform. Today the test was the cross-border data-flow basis under three regulators.

*Different test,* she thinks. *Same chain. The chain held. It always does, when the operator has done the work. Sun-Won did the work.*

---

## 🧾 Final Assessment Theme

> **"Three regulators reading the same chain through three different lenses. Today, the chain plus the contract answers each lens. After Q3, the verifier output answers all three on its own. That is the upgrade Sun-Won is paying for — and it is the cleanest cross-border story we have audited yet."**
>
> — Karen, lead auditor, Sun-Won Cosmetics Group, April 9, 2026

---

*End of Story 09.*
