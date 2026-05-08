🧾 A Day in the Life: 8 Auditors Verifying Customer Interaction Integrity at Northbridge Federal Savings

**Context:**
Northbridge Federal Savings — a regional US bank, ~$45B consolidated assets, OCC-supervised national bank, FDIC-insured. Engagement type: FFIEC IT Handbook supplementary review. The bank closed an MRA on customer-data integrity two quarters ago. This is the verification revisit.

**Posture going in:** the audit team has not seen this engagement before. Northbridge has been running something called "TesseraSeal" across customer-data capture for 18 months — independent of any prior conversation with this audit team. The team will encounter the product for the first time at the kickoff meeting. The bank's claim, going into the room, is that every customer-facing surface (CRM mirror, voice/recordings, branch tablets, API edges, IAM events, the AI advisor) lands in a sealed chain-of-custody ledger called Herald Core, that the regulator-facing surface is called Herald.Compliance, that daily seals run via a CloudHSM-managed signing key, that the verifier CLI is `herald-verify`, and that the whole stack conforms to a public spec called FFIEC chain-of-custody v1.0a. The team has heard pitches before that promised this much. Skeptical-but-listening is the working posture.

---

## 👥 The Audit Team

- **Karen** — Lead Auditor (governance + narrative)
- **Raj** — Database specialist
- **Elena** — CRM systems
- **Mike** — Application / API layer
- **Diana** — IAM & access control
- **Luis** — DevOps / logs / pipelines
- **Chen** — Data engineering / ETL
- **Tom** — Internal-audit liaison specialist (visiting team; partners with the client CAE)

**Client liaison:** Marcus Tan, Chief Audit Executive, Northbridge Federal Savings. Calm. Has done this before with FDIC examiners.

---

## 🌅 8:30 AM — Kickoff Meeting

The team rolled in to the Northbridge engagement room with the look of people who had spent last week somewhere unpleasant.

Karen poured coffee. She stared at the slide on the projector — a clean architecture diagram, every box labeled, every arrow ending at something called "Herald Core ledger." Above the diagram, in modest type: **TesseraSeal — chain-of-custody for customer-data capture.**

She had not heard the name before.

She turned to Raj.

"Last week was a graveyard," she said. "This week, I want to find at least one thing."

Raj nodded without looking up from his laptop. "Bet you a coffee you find a Gap by lunch."

"Bet you two coffees I find a Gap by 10 AM."

*It never is*, Karen thought. *Diagrams are clean until you ask the third question.*

Marcus Tan walked in. Mid-fifties. Pressed shirt. Coffee in his left hand, a thin folder under his right arm.

"Karen. Tom. Welcome to Northbridge."

Tom shook his hand. "Marcus."

"You'll see a name on the deck you may not have run into before," Marcus said, gesturing at the slide. "TesseraSeal. We've been on it for eighteen months. It's the chain-of-custody layer behind every customer-data capture path the bank operates. I'd rather walk you through what it is than make you guess from the diagram."

"Go ahead," Karen said.

Marcus stayed standing. He didn't sit. He didn't pull out a deck-of-decks.

"At the bottom is Herald Core. That's the append-only ledger. Above it is Herald.Py — that's the SDK we instrumented every customer-facing surface with: CRM mirror, voice transcription, branch tablets, the core-banking API edges, IAM, the AI wealth advisor. Every event lands in the ledger as a sealed entry. Daily, the system computes a Merkle root over the day's entries and signs it with a CloudHSM-resident Ed25519 key. The signed seal is published on a regulator-facing surface called Herald.Compliance. The verifier CLI is `herald-verify`. The whole stack conforms to a public spec — FFIEC chain-of-custody v1.0a — and the verifier is open-source, so you can run it on your own laptops without any of our credentials at any layer."

He paused.

"That's the elevator. I know it's a lot to take in cold."

Karen wrote on her notepad: **TesseraSeal — verify claims.**

Underneath: *spec public; verifier OSS; key on Compliance page; ledger append-only.* She underlined "verify."

Raj was already typing the spec name into his laptop search bar.

"Can I see the spec?" he said. "URL."

Marcus read it off. Raj wrote it down.

Mike said, "Eighteen months. So this isn't new to you, but it's new to us."

"It's new to you. It's not new to the FDIC — they've examined under it twice. The closing report from the prior-year MRA references the verifier outputs by entry-ID. I have copies in the workpaper pack."

Tom said, "Same drill as the FDIC visit in February, then?"

"Same drill," Marcus said. "I'll route you through the surfaces. SRE on-call is Greg today. Greg has done this before. Verifier credentials are already provisioned for your laptops — read-only, scoped to the Compliance surface."

Karen blinked. "You provisioned us before we asked."

"The verifier's design is that you don't need our credentials at all. The Ed25519 public key is published on the Herald.Compliance page. You can pull a seal record and verify it on a coffee shop wifi if you want. The credentials are just to save you the trouble of typing the tenant ID."

*Hm.*

Karen added a second line under her **TesseraSeal — verify claims** note: *seal verification is unprivileged.*

She paused. She added a third line: *check this claim before lunch.*

"Let's start with the architecture overview," she said. "I want to know what you think you have. Then we'll go look."

Marcus didn't bristle. He clicked to the next slide.

*Most CAEs bristle when I say 'then we'll go look,'* Karen thought. *He didn't. Either he is very tired, or he has nothing to defend. We'll find out which.*

*And we'll find out today whether 'TesseraSeal' is a slogan or a system.*

---

## 🧩 9:15 AM — First Crack in the Story

Marcus walked through the diagram. CRM mirror on the left, core-banking-API edges in the middle, the AI advisor box (a Llama-based wealth-recommendation model) tucked into a corner, the contact-center voice transcription path, the loan-decisioning workflow, IAM events, all draining into Herald Core.

Mike raised a hand at the AI advisor box.

"You're capturing model inputs and outputs?"

"Inputs, outputs, model version, system prompt fingerprint, retrieval-augmented context. Every recommendation that touches a customer file lands as a sealed entry."

"System prompt fingerprint as in a hash?"

"Hash. The full prompt is in the ledger too, but the fingerprint is what cross-references the model-governance registry."

Mike wrote that down. He didn't say anything.

Elena, who had been quietly reading the Salesforce architecture page, looked up.

"You're not running Salesforce-native logs."

"We are running Salesforce-native logs," Marcus said. "We also mirror every customer-touching field change into Herald.Py via a connector. The Salesforce-native log is the operational log. The Herald mirror is the chain-of-custody log."

"Two logs," Elena said.

"Two logs. The connector lag is something we can talk about later if you want."

Karen made a note: *connector lag — come back to this.*

She glanced at Raj. He was scrolling through a list of what looked like seal records. He didn't appear to be enjoying himself.

"Raj?" she said.

"I want to see the schema," Raj said. "Of the chain table. Right now."

"Greg," Marcus said into his phone, "Raj wants the chain table schema. Can you pull it up on the second screen?"

The second screen lit up.

> ### ✓ Confirmation #1 — Chain table schema is the spec
>
> The schema matched FFIEC chain-of-custody v1.0a line-for-line. `entry_id`, `prev_hash`, `entry_hash`, `hmac_sha256`, `tenant_binding_kdf_label`, `event_payload_jcs`, `merkle_leaf_index`, `seal_date`, `signature_ed25519`, `signing_key_fingerprint`. Indexed on `entry_id` and `seal_date`. No `updated_at` column — entries are append-only by schema, not by convention.

Raj said, "Where's the update audit trigger?"

"There isn't one," Marcus said.

"Why?"

"Because the table doesn't accept updates. The role that writes to it has INSERT only. The role that reads from it has SELECT only. There is no role with UPDATE or DELETE on the chain table in any environment, including production-DBA."

Raj leaned back.

"What about the role that *creates* roles?"

"Bootstrapped at deployment time. The role-creation role itself was retired after the system went live. It does not exist in IAM today. We can show you the IAM history — every role grant and revocation is, as you'd expect by now, a sealed chain entry."

Raj said, "I will want to see that."

"After the schema review, sure."

*Two coffees*, Karen thought. *I owe Raj two coffees.*

Elena, who had been listening, leaned forward. "You said the Salesforce mirror lands in the same chain. Same schema?"

"Same schema. Different `event_class` tag. The CRM mirror entries carry a `source=salesforce` annotation and the connector's run-id, but they go through the same hash, same MAC, same Merkle root, same daily seal."

"And the connector itself — its run-ids — are those chained?"

"The connector's lifecycle events are chained. Start, completion, failure, retry. Every batch is a chain entry that references the customer-data entries it produced."

Elena wrote: *connector lifecycle is auditable.*

---

## 🧠 10:00 AM — Database Deep Dive

Before Raj opened the database session, he wanted something else.

"Marcus," he said. "Print me the chain envelope schema. The full attribute table. Whatever the spec calls Required and Optional. I want to read it before I touch the database."

"Section 4.4 of the spec," Marcus said. He nodded at the second screen. "Greg, pull the v1.0a attribute table and the canonical-bytes definition. Both."

Two blocks appeared on the screen. The first one was the attribute table. Raj read it slowly, top to bottom.

```
# FFIEC chain-of-custody v1.0b — chain envelope (per-event)
# Per spec §4.4 attribute table

ffiec.chain.spec               = "v1.0"        # Required
ffiec.chain.format_version     = "v1"          # Required
ffiec.chain.posture            = "ffiec"       # Required (resource-level)
ffiec.chain.chain_kind         = <enum>        # Required: audit | model_call | tool_call | routing | translation | operational
ffiec.chain.run_id             = <string>      # Required: run identifier (= chain_id)
ffiec.chain.tenant_id          = <string>      # Optional: 1-255 chars [A-Za-z0-9_.-] per §3
ffiec.chain.captured_at        = <RFC3339-ns>  # Required: UTC nanosecond precision
ffiec.chain.seq                = <int>         # Required: 1-indexed position
ffiec.chain.payload_hash       = <64-hex>      # Required: HMAC-SHA-256 (lowercase hex)
ffiec.chain.prev_hash          = <64-hex>      # Required: previous payload_hash; 32 zero bytes for seq=1
ffiec.chain.key_version        = <int>         # Required: 1-indexed IKM generation
ffiec.chain.key_fingerprint    = <32-hex>      # Required: SHA-256(utf8(tenant_id) || ikm)[:16]
ffiec.chain.mac_computed_at_utc = <RFC3339>    # Required: writer's wallclock (forensic, not security)
ffiec.chain.kms_handle_uri     = <string>      # Required: KMS handle URI; "plaintext-dev" marks dev adapter
ffiec.chain.canonical_encoding = "rfc8785-jcs" # Optional, default at format_version=v1
ffiec.chain.late_binding       = <bool>        # Optional: true for events arriving after their day was sealed
ffiec.chain.region             = <string>      # Optional: SDK's binding region for multi-region resilience

# Herald-internal canonical-bytes-load-bearing fields (per OtlpAttributeKeys.cs)
herald.event_id                = <UUID>        # Required for byte-reconstruction
herald.duration_ns             = <int>         # Optional; absence encodes None at the receiver
herald.kind                    = <SpanKind>    # OTel SpanKind shape
herald.severity                = <string>      # AUDIT | INFO | WARN | ERROR | FATAL

# Standard OTel
service.name                   = <string>      # Required at Resource
service.version                = <string>      # Required at Resource
```

Raj didn't say anything for a full minute. He read the block twice. Then he started in.

"What goes into `payload_hash` exactly?"

"HMAC-SHA-256 over the canonical bytes of the event," Marcus said. "Per §4.4.1. The canonical bytes are RFC 8785 JCS over the entire event, with a documented exclusion set — the chain stamp itself is excluded, since the chain stamp is computed from those bytes. The HMAC key is the per-tenant HKDF-derived key from the IKM. Lowercase hex, 64 characters."

"And `kms_handle_uri` is Required, but it's not in the canonical bytes?"

"It's Required because every entry needs to identify which KMS-managed IKM generated the per-tenant HKDF key. It's not in the per-event MAC's canonical bytes because the MAC is over the *event*, not over the operational metadata that names the key. The handle URI is operational metadata — what KMS, which key alias, which version. It binds the entry to a key without circular-defining what the MAC covers. Forensic, not security."

"Why the distinction?"

"Because if you put the URI inside the MAC, then rotating the URI string — say, renaming the key alias for an unrelated reason — would invalidate every prior MAC. We pin the *fingerprint* under the MAC via `key_fingerprint`. The URI is a label that points at the same key the fingerprint identifies. The fingerprint is what the verifier checks."

Raj nodded slowly. He underlined `key_fingerprint` on his notepad.

"What's the difference between `mac_computed_at_utc` and `captured_at`?"

"`captured_at` is the wallclock at the moment of the application event — the customer interaction, the API call, the model inference. That's the timestamp the business cares about. `mac_computed_at_utc` is the wallclock at the moment the SDK actually sealed the entry. They're usually within a few milliseconds, but they can diverge if the SDK is under buffer pressure or recovering from a sidecar fault. Per §4.4, `mac_computed_at_utc` is forensic — it tells you what the writer's clock said at MAC time, even if the writer's clock was wrong. We don't trust it for security. We trust it for forensic reconstruction."

"If the writer's clock is wrong, does the verifier care?"

"The verifier doesn't trust either timestamp for cryptographic purposes. It cares about chain order, prev_hash linkage, MAC validity, and Merkle inclusion. Timestamps are for humans investigating an incident. The spec is explicit about that — §4.4 calls out 'forensic, not security' on `mac_computed_at_utc`."

Raj wrote that down. Then he wrote a follow-up underneath, with a star next to it: *come back to clock-skew handling under §10.16 lag bounds.*

"`tenant_id` is Optional. What do the canonical bytes look like for a single-tenant deployment?"

"If the field is absent at the source, it's absent in the canonical bytes — JCS doesn't emit a key for a field the producer didn't set. The HKDF derivation falls back to a documented single-tenant label. The fingerprint computation uses the empty string for the `utf8(tenant_id)` portion of `SHA-256(utf8(tenant_id) || ikm)[:16]`. Section 3 covers it. Northbridge sets `tenant_id="northbridge-bank-prod"` everywhere — we're explicitly not single-tenant in the SDK's sense, even though we're a single bank."

"Why?"

"Because the bank has subsidiary entities. Each subsidiary is a separate `tenant_id` under the same Herald Core deployment. The prod tenant is `northbridge-bank-prod`. The wealth-management subsidiary is `northbridge-wealth-prod`. The merchant-services subsidiary is `northbridge-merchant-prod`. Three tenants, three HKDF derivations, three independent chains."

Raj nodded.

"One more. What stops someone from re-emitting an entry with a different `key_version` to claim it was sealed under a rotated key?"

Marcus paused. Not because he didn't have the answer — because this was the question he was waiting for.

"Three things. First, the `key_fingerprint` is in the canonical bytes — re-emitting with a different `key_version` doesn't change the fingerprint, so the MAC won't recompute under a different key. Second, the per-tenant HKDF derivation is keyed on IKM generation; a wrong `key_version` produces a different HKDF output and the MAC fails. Third, the daily Merkle seal is computed over the day's entries as written; you can't slip a re-emitted entry into a sealed day. The combination is what makes key-rotation transparent to the verifier — the verifier reads the entry, picks the correct IKM generation by version, recomputes, and either matches or doesn't."

Raj said, "I want to see the canonical bytes the MAC actually covers."

Marcus typed. The second block appeared on the screen.

```
canonical bytes (v1.0a, per spec §5) = RFC 8785 JCS over the event with:
  EXCLUDED: ffiec.chain.payload_hash, ffiec.chain.prev_hash,
            ffiec.chain.key_version, ffiec.chain.key_fingerprint,
            ffiec.chain.mac_computed_at_utc, ffiec.chain.kms_handle_uri
            (the chain stamp itself is excluded — it's computed from these bytes)
  INCLUDED: everything else, including ffiec.chain.{spec, format_version,
            chain_kind, run_id, tenant_id, captured_at, seq}
            plus the entire payload (event-class-specific attributes + body)
```

Raj read it twice.

"That's a clean exclusion set," he said. "The chain stamp is the only thing not under its own MAC. Everything else is bound."

"Section 5 is the smallest section in the spec," Marcus said. "It needs to be unambiguous more than it needs to be long."

Raj wrote *§5 is short on purpose* on his notepad and underlined it.

"OK," he said. "Now I want to look at the database."

Raj opened the database session. He had a list of forty queries he runs against any chain-of-custody system. He started with the soft ones.

```
SELECT COUNT(*) FROM chain_entries WHERE prev_hash IS NULL;
```

One row. The genesis block. As specified.

```
SELECT entry_id, prev_hash, entry_hash FROM chain_entries
ORDER BY entry_id LIMIT 10;
```

Ten rows. Each `prev_hash` matched the previous row's `entry_hash`. Raj didn't say anything. He kept going.

He ran the chain-walk verifier against an arbitrary 50,000-row window. The verifier finished in 11 seconds. Exit code 0.

He ran it again with `--strict`. Same result. 11 seconds. Exit code 0.

He picked a random row, took its `entry_hash`, and ran a manual SHA-256 over the canonicalized payload (RFC 8785 JCS) plus the `prev_hash`. The hash matched the stored `entry_hash`.

He did the same thing for the HMAC, recomputing it with the per-tenant HKDF-derived key. The HMAC matched. Constant-time comparison was visible in the verifier source code.

> ### ✓ Confirmation #2 — Per-event MAC and chain hash both recompute
>
> Raj independently recomputed both the SHA-256 entry hash and the HMAC-SHA-256 MAC for a sampled entry, using the documented v1.0a 10-line `sign_payload` form. Both matched. The HKDF tenant-binding label resolved to the documented `tenant=northbridge` derivation.

Raj sat back. He took a long drink of coffee.

"Sample size?" Karen asked him quietly.

"Fifty thousand entries on the chain walk. One entry recomputed by hand. I'll do another twenty by hand before lunch."

"Take your time."

Raj didn't take his time. He ran twenty more by hand inside fifteen minutes. They all matched. He picked entries from the start of the 18-month window, the middle, and the most recent week. He picked entries across two CloudHSM key-rotation boundaries.

Each one matched.

*This is what consistent hashing across an 18-month window is supposed to look like*, he thought. *I have been doing this job for nine years. I have never actually seen it.*

"What's the daily seal cadence?"

"Daily," Marcus said. "Merkle root over the day's entries. Signed Ed25519 by a CloudHSM-resident key. Key fingerprints rotate quarterly. The current fingerprint is on the Compliance page."

"Show me a daily seal record."

Marcus pulled up the seal for 2026-04-15. Merkle root, signature, public-key fingerprint `7f3a9...`, leaf count, the date range, and a JCS hash of the metadata block.

Raj copied the public-key fingerprint and pasted it into a comparison against the published Compliance page. Match.

"Run the verifier on this seal."

```
herald-verify --tenant=northbridge --date=2026-04-15 --strict
```

Output:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key 7f3a9...
Elapsed: 4.1s
```

Raj read the output twice.

"`--explain`," he said.

```
herald-verify --tenant=northbridge --date=2026-04-15 --strict --explain
```

A 47-line trace scrolled past. Every step from genesis traversal through Merkle tree resolution to signature verification was named and timed.

Raj closed the laptop halfway. Not all the way. Halfway.

"I want to look at IAM next."

---

## 🔐 11:00 AM — IAM Review

Diana took over. She had a specific scenario she wanted to test: the "temporary admin" pattern that breaks every chain-of-custody system she has ever audited.

"Walk me through how a DBA gets emergency write access to the chain table."

"They don't," Marcus said.

"That's not an answer. That's a slogan."

Marcus smiled. "Fair. Let me re-answer. Emergency access to the chain table is not a feature of the system. There is no break-glass account. The deployment runbook for chain-corruption recovery is to roll forward from the last sealed Merkle root and reconstruct downstream views — never to mutate the chain in place. We tested this in disaster-recovery drills in Q1 and Q3."

"What about temporary admin elevation for *other* surfaces? Salesforce admin, AI advisor model deployment, that kind of thing."

"Temporary admin works the way it works in any decent IAM system. Elevation request, approval, time-boxed role grant, auto-revocation at 24 hours."

"And the auto-revocation — is that a cron job that someone could turn off?"

"It's a chain-driven workflow. The grant itself is a chain entry. The expiration is a chain entry. The revocation is enforced by a worker that reads the chain and applies the role removal. If the worker is down, a separate health check fires. If both the worker and the health check are down, IAM fails closed — the role lookup defaults to the unprivileged baseline."

> ### ✓ Confirmation #3 — IAM events are themselves chain-of-custody captured
>
> Every IAM grant, revocation, and elevation request lands as a sealed chain entry in the same Herald Core ledger as customer-data events. The auto-revocation worker is chain-driven, not cron-driven. Diana sampled three temporary-admin grants from the past 90 days. Each had a matching revocation entry, each landed within 30 seconds of the 24-hour mark, each was sealed in the daily Merkle root.

Karen's pen paused over the notepad.

*It never is*, she thought. *Except apparently this time.*

She crossed out the *It never is* she had written at 8:30. She didn't write anything in its place.

Diana asked, "Does the verifier work on IAM entries the same way it works on customer-data entries?"

"Same verifier. Same exit codes. Same chain. The IAM entries are tagged with `event_class=iam` for filtering, but the chain-walk and seal verification don't distinguish."

Diana ran:

```
herald-verify --tenant=northbridge --date=2026-04-15 --event-class=iam --strict
```

Status: PASS. Step: 12. 3.7 seconds.

She ran it without the filter.

Status: PASS. Step: 12. 4.0 seconds.

She closed her laptop fully.

"Lunch?"

She added one more line to her notepad before standing up: *the IAM-as-chain pattern is the part I want to write down for other engagements. Don't bury it.*

---

## 🧪 12:00 PM — Lunch (But Not Really)

The team ordered sandwiches into the engagement room. Nobody left the building.

Karen walked over to where Tom and Marcus were standing by the window, mid-conversation about audit-procedure cross-references.

"Tom, what are we at on findings?"

"Zero Gaps. Zero Partials. One thing Elena flagged that I want to come back to after lunch — the Salesforce mirror lag wording."

"Is it a Gap?"

"It's a Nit. The mirror works. The seal works. The documentation says 'near real-time' without quantifying it."

"That's a Nit."

"That's a Nit."

Karen looked at Marcus. "Your CAE liaison just told my internal-audit liaison we have zero Gaps at noon."

Marcus said, "It's noon. There's still time."

Karen laughed. She actually laughed. She hadn't laughed during a workpaper-week since 2024.

She took a sandwich. She sat down. She looked at Raj across the table.

Raj said, "Two coffees."

"Two coffees," Karen agreed.

---

## 🔄 1:00 PM — API Layer Inspection

Mike had been waiting. He liked the API layer because the API layer is where systems lie.

He pulled up the core-banking API gateway logs. He picked a single request — a wire-transfer authorization, timestamped 2026-04-12T14:23:11.847Z, customer ID redacted, transaction ID `tx_8a7f...`.

He found the corresponding chain entry. He extracted the `entry_id`. He ran:

```
herald-verify --entry-id=ce_4f29d8a3b1... --strict
```

Output:

```
Status: PASS
Step: 12
Reason: chain integrity verified, HMAC recomputed,
        Merkle path resolved, signature verified
        against public key 7f3a9...
Elapsed: 0.8s
```

> ### ✓ Confirmation #4 — Single-entry verification resolves end-to-end
>
> A single API-call entry, picked from operational logs by transaction ID, verified through the full chain-of-custody pipeline: per-event hash, per-tenant HMAC, Merkle inclusion proof, daily seal Ed25519 signature against the published public key. 0.8 seconds, no Northbridge credentials beyond read scope on the Compliance surface.

Mike rotated. He picked a different request. A failed authorization. A retry. A reversal.

All three: PASS, PASS, PASS.

In the verifier output for the reversal, Mike noticed an attribute he hadn't seen before. He scrolled back through his earlier samples. It wasn't on most of them.

```
ffiec.chain.late_binding = true
```

"Marcus, what's `late_binding=true`?"

Marcus said, "Section 4.2.2 of the spec. An event whose `received_at` UTC date is after the seal day was already sealed gets stamped with `late_binding=true` and lands in the next day's seal. The original seal isn't altered. The verifier reports `late-binding entries: N` as an anomaly line under PASS. Normal operational case — a connector backlog or a clock-skew event — but it's visible. You can't sneak a late event in without it being labeled as late."

"So this reversal arrived after its day was sealed."

"Probably a pre-dawn reversal that was stuck behind a connector retry. The chain entry's `captured_at` is the original transaction time. The seal that includes it is the next day's seal. The `late_binding=true` flag tells the verifier to expect a seal-date offset. Without the flag, that offset would look like a tampered timestamp."

Mike wrote: *late_binding=true is a positive declaration, not a hidden state.*

"And the verifier counts them?"

"Counts them, prints the total in the anomaly section under PASS. If a day suddenly has thousands of late-binding entries, that's a connector-health signal worth chasing. Most days, it's single digits."

Mike nodded. He moved on.

He picked an AI-advisor recommendation — the kind of event that he expected to be the thinnest seam. The model output, the system prompt fingerprint, the retrieval context, the customer ID, the recommendation text.

PASS.

"How does the model recommendation get from the model into the chain?"

"Herald.Py wraps the inference call. The wrapper captures inputs, outputs, model version, prompt fingerprint, retrieval context. Synchronous capture. The chain entry lands before the recommendation is rendered to the customer."

"Synchronous? Latency cost?"

"Single-digit milliseconds. The chain write is local-buffer with a write-ahead log; the seal is asynchronous later. The recommendation isn't rendered until the buffer write returns."

Mike asked another thing.

"What about read paths? Customer data flowing *out* of the bank — is that captured?"

"Read paths are captured at the same edge. Every API response that returns customer data lands as a sealed chain entry tagged `event_class=read`. The payload includes the request fingerprint, the responding system, the customer ID, and the field set returned."

"Including bulk-export jobs?"

"Especially bulk-export jobs. The export job itself is a chain entry. Each row produced by the export is a chain entry. The relationship between the parent job and the child rows is encoded in the chain. We can audit a regulator-style export end-to-end — who ran it, what they pulled, what they got, when it was sealed."

Mike said, "That's better than what most banks have for write paths."

Marcus didn't smile. He just nodded.

Mike asked one more thing.

"What happens if the buffer write fails?"

"The recommendation isn't rendered. The customer sees a soft error. The retry logic is in the Herald.Py wrapper. There's a circuit breaker; if it trips, the AI advisor fails closed and customers get a 'temporarily unavailable' message until the path recovers. The bank prefers a degraded-experience customer to an un-audited recommendation."

> ### ✓ Confirmation #5 — AI advisor fails closed when capture fails
>
> The customer-facing AI recommendation surface is gated on successful chain capture. A failed capture results in a degraded customer experience, not an un-audited recommendation. This is enforced in the Herald.Py wrapper, not as an operational policy.

Mike wrote that down. He underlined it.

"That's the part I usually have to argue people into," he said. "It's already done."

---

## 🧬 2:00 PM — Data Pipeline Reality

Chen and Luis tag-teamed the next hour.

Luis went first. He wanted to know what the Herald Core retention story looked like, and specifically whether anyone could delete log groups.

"Append-only," Marcus said. "The chain table itself is append-only by role. The seal records are append-only by role. The retention policy is enforced by the storage tier — object lock, immutability window matching the FFIEC retention requirement, no role with delete permission inside the window."

"Even an account root?"

"Even an account root. The signing key is in CloudHSM, and the storage account has a separate trust boundary. Account root in the application AWS account cannot reach into the storage account's bucket."

"What about the storage account's root?"

"Object lock with a compliance-mode retention period. Account root in the storage account cannot bypass it either. The retention period exceeds the FFIEC requirement by a margin."

> ### ✓ Confirmation #6 — Append-only at the storage tier, not at the convention tier
>
> Log retention is enforced by S3 object lock in compliance mode, not by IAM convention. The storage account has a separate trust boundary from the application account. No principal — including either account root — can delete or mutate sealed chain entries within the retention window.

Luis closed his laptop. He looked at Karen.

"That's the thing the last bank could not show me."

Karen nodded.

Luis went on. "I want to see one more thing. The CloudWatch retention story for non-Herald operational logs. Application logs, infrastructure logs, the stuff that *isn't* the chain but lives next to it."

Marcus said, "Standard CloudWatch with retention policies set per-log-group. Engineers can stop a log stream but cannot retroactively delete entries within retention. The policy itself is in the chain — every retention-policy change lands as an `event_class=ops` chain entry."

"Even retention-policy changes are chained."

"Especially those. The one thing we never want is for someone to be able to silently shorten retention. Retention shortening is itself a sealed event with the role that requested it, the prior policy, the new policy, and the time-to-effect."

Luis wrote that down.

*This is the part where most banks tell me 'we'll get to it next quarter,'* he thought. *They've gotten to it.*

Luis paused.

"One more thing before I hand off to Chen," he said. "I want to see what's actually getting captured. Source side and chain side, end to end. The Salesforce mirror connector specifically — that's the surface Elena flagged this morning. Show me the raw CDC event Salesforce produces, and show me what the chain wrote for the same event. Side by side."

Marcus said, "I can do that. Greg, split-pane terminal. Pull a single CDC event from the test-replay archive and the matching chain entry. The case-assigned event from yesterday afternoon would be a good one — it crossed the connector lag window cleanly."

The second screen lit with two panes.

**Left pane — Salesforce CDC stream (raw source side):**

```jsonc
{
  "schema": "OYZGY...PUL",
  "payload": {
    "ChangeEventHeader": {
      "entityName": "Case",
      "recordIds": ["500Hp00001abcXYZ"],
      "changeType": "UPDATE",
      "changedFields": ["Status", "OwnerId", "Description"],
      "transactionKey": "001a2b3c-4d5e-6789-abcd-ef0123456789",
      "sequenceNumber": 1,
      "commitTimestamp": 1746719234123,
      "commitNumber": 12345678,
      "commitUser": "0051Hp00001AGENT001"
    },
    "Status": "Working",
    "OwnerId": "0051Hp00001AGENT001",
    "Description": "Customer requesting refund per email dated 2026-05-08. Damaged item ORD-7741922.",
    "LastModifiedDate": "2026-05-08T13:43:02.123Z"
  },
  "event": {
    "replayId": 9874321
  }
}
```

**Right pane — the chain-of-custody log entry written by the Salesforce mirror connector (a single line of NDJSON in the chain file):**

```jsonc
{
  "ffiec.chain.spec": "v1.0",
  "ffiec.chain.format_version": "v1",
  "ffiec.chain.chain_kind": "audit",
  "ffiec.chain.run_id": "500Hp00001abcXYZ",
  "ffiec.chain.tenant_id": "northbridge-bank-prod",
  "ffiec.chain.captured_at": "2026-05-08T13:43:02.481000000Z",
  "ffiec.chain.seq": 3,
  "ffiec.chain.payload_hash": "9b3c7e1a4f8d62b0c5a1e9f4d2b8c0e3a7d1f5b9c4e8d2a6b0f3e7d1a4b8c2e5",
  "ffiec.chain.prev_hash":    "ad7e2c9f1b4d8e3a0c5b9d2e6f1a4c7e0b3d9f2a5c8e1b4d7f0a3c6e9b2d5f8a",
  "ffiec.chain.key_version": 4,
  "ffiec.chain.key_fingerprint": "7f3a9c1e5b8d2f4a6c0e9b3d7f1a4c8e",
  "ffiec.chain.mac_computed_at_utc": "2026-05-08T13:43:03.122Z",
  "ffiec.chain.kms_handle_uri": "aws-cloudhsm:cluster/cluster-7y8a/key/k-northbridge-prod-2026q2",
  "ffiec.chain.canonical_encoding": "rfc8785-jcs",
  "ffiec.chain.region": "us-east-1",
  "herald.event_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
  "herald.kind": "INTERNAL",
  "herald.severity": "AUDIT",
  "service.name": "northbridge-sf-mirror",
  "service.version": "2.4.1",
  "event": "case.assigned",
  "case_id": "500Hp00001abcXYZ",
  "assignee_id": "0051Hp00001AGENT001",
  "queue": "refunds-tier-1",
  "assigned_at_utc": "2026-05-08T13:43:02.481Z",
  "source_replay_id": 9874321,
  "source_commit_timestamp": 1746719234123,
  "source_lag_observed_ms": 1358,
  "redaction_policy_id": "northbridge-pii-v7",
  "redaction_policy_version": "2026-04-15"
}
```

Luis leaned forward.

"Walk me through the lineage," he said. "How does the right side prove it captured the left side."

Marcus traced it with a stylus.

"Salesforce's `replayId=9874321` and `commitTimestamp=1746719234123` are both captured into the chain entry as `source_replay_id` and `source_commit_timestamp`. Anyone can verify the Salesforce side independently — Salesforce keeps replay IDs for 72 hours, longer with EventLogFile. Pull the raw CDC stream by replay ID, line up against the chain entry, confirm the source-side metadata matches. The chain entry's MAC binds those two fields into the canonical bytes, so the connector can't claim a different replay-ID after the fact."

He moved the stylus.

"The mirror connector recorded `source_lag_observed_ms=1358` — meaning 1.358 seconds elapsed between Salesforce's commit and the chain entry's MAC computation. That's well inside the §10.16 lag bound the bank published in CC8.1."

Karen wrote down: *§10.16 — the lag bound is published.*

Marcus moved the stylus again.

"`redaction_policy_id` and `redaction_policy_version` are bound under the per-event MAC. If anyone retroactively edits the policy, the chain entry's signed reference to the old version doesn't change — the mismatch surfaces at audit. The policy itself is in a separate signed registry; the chain entry pins which policy version was active at capture time."

"And `tenant_id` plus `service.name` together —"

"— makes the source unambiguous. `ffiec.chain.tenant_id="northbridge-bank-prod"` plus `service.name="northbridge-sf-mirror"` says: this entry was written by the Salesforce mirror connector running under the production-bank tenant, not by some other Northbridge process. The Wealth subsidiary's connector would write under `northbridge-wealth-prod`. They don't overlap. They can't overlap — the IKM generations are different, so the MACs wouldn't validate cross-tenant."

Karen's pen had not stopped moving. She had a different question.

"What if Salesforce silently drops a CDC event? The chain wouldn't see it."

Marcus said, "The connector reconciles against Salesforce's source-side commitNumber sequence on a 5-minute cadence. Gaps trigger `connector.outage` events that get chained. If Salesforce drops one, the gap shows up in the chain itself, not just in the connector's logs. Section 10.16 of the spec covers SaaS-edge connectors specifically — the freshness invariant requires the connector to either reconcile a missing source-side sequence within the published bound, or emit a chained outage entry naming the gap. We can't claim a gap doesn't exist by staying silent — silence inside the lag window is itself a violation."

"Show me one of the outage entries."

Marcus typed. A third pane appeared.

```jsonc
{
  "ffiec.chain.run_id": "northbridge-sf-mirror",
  "ffiec.chain.chain_kind": "operational",
  "ffiec.chain.seq": 8847,
  "ffiec.chain.captured_at": "2026-03-17T14:08:23.000000000Z",
  "service.name": "northbridge-sf-mirror",
  "event": "connector.outage",
  "outage_kind": "source_sequence_gap",
  "expected_commitNumber_low": 12340012,
  "expected_commitNumber_high": 12340034,
  "observed_commitNumber_at_resume": 12340035,
  "missing_count": 23,
  "remediation_action": "salesforce_eventlogfile_replay_initiated",
  "remediation_run_id": "sf-replay-2026-03-17-008"
}
```

"That's the connector saying 'I lost 23 events between commitNumbers 12340012 and 12340034, and I'm replaying them from Salesforce's EventLogFile.' The remediation run is its own chain — `sf-replay-2026-03-17-008` — and each replayed event lands as a chained entry tagged `late_binding=true` because they arrive after the day they were originally committed. The verifier reports them as anomaly-line late-binding under PASS. Visible, named, traceable."

Luis wrote: *outages are positive events. Silence is a violation.*

> ### ✓ Confirmation #7 — SaaS-edge connector lineage is end-to-end visible
>
> Luis side-by-sided the Salesforce CDC source event and the chain entry written by the mirror connector. The Salesforce-side metadata (`replayId`, `commitTimestamp`) is bound under the per-event MAC as `source_replay_id` and `source_commit_timestamp`. The connector's observed lag is recorded against the §10.16-published bound. Source-side sequence gaps surface as chained `connector.outage` events; replayed events land with `late_binding=true`. The chain itself reports gaps — silence is not an option.

Luis didn't move on.

"One more thing," he said. "Show me the actual connector. The Python code that turns Salesforce CDC events into chain entries. I want to see what your engineer wrote."

Marcus said, "Greg, pull `subscriber.py` from the sf_mirror repo. Read-only, current main."

The connector source filled the second screen.

```python
# E:/northbridge/connectors/sf_mirror/subscriber.py
"""
Salesforce CDC mirror connector — Northbridge production.

Subscribes to Salesforce Change Data Capture (CDC) channels for the
Case, EmailMessage, and Account entities. Translates each CDC event
to a chain-of-custody entry under tenant_id="northbridge-bank-prod"
via the Herald.Py SDK. The chain handles the integrity bind; this
connector is responsible for fidelity (no source events lost) and
freshness (lag stays inside §10.16 bounds published in CC8.1).

Per FFIEC chain-of-custody v1.0b §10.16 + §10.18 + §10.19.
"""

import asyncio
import logging
from datetime import datetime, timezone
from typing import Any

import herald
from aiosfstream import SalesforceStreamingClient

logger = logging.getLogger(__name__)

CDC_CHANNELS = (
    "/data/CaseChangeEvent",
    "/data/EmailMessageChangeEvent",
    "/data/AccountChangeEvent",
)


async def subscribe(
    sf_client: SalesforceStreamingClient,
    tenant_id: str = "northbridge-bank-prod",
) -> None:
    """
    Subscribe to all CDC channels for the configured tenant. For each
    event, write a chain-of-custody entry with the Salesforce-side
    metadata bound under the per-event MAC.
    """
    async for sf_event in sf_client.subscribe(CDC_CHANNELS):
        await _emit_chain_entry(sf_event, tenant_id)


async def _emit_chain_entry(sf_event: dict[str, Any], tenant_id: str) -> None:
    payload = sf_event["payload"]
    header = payload["ChangeEventHeader"]
    entity = header["entityName"]
    record_id = header["recordIds"][0]  # CDC envelopes one record at a time
    change_type = header["changeType"]
    commit_user = header["commitUser"]
    commit_ts_ms = header["commitTimestamp"]

    # Run is keyed on the Salesforce record ID. A given Case stays in one
    # run for its entire lifecycle, so seq advances monotonically across
    # all CDC events for that record.
    run_id = record_id

    # Source-side observed lag — CDC commit to here. Used for §10.16
    # connector.lag_observation events emitted on a separate cadence
    # (see lag_observer.py).
    now_utc = datetime.now(timezone.utc)
    commit_dt = datetime.fromtimestamp(commit_ts_ms / 1000, tz=timezone.utc)
    lag_ms = int((now_utc - commit_dt).total_seconds() * 1000)

    event_name, attrs = _translate(entity, change_type, payload)

    with herald.run(
        run_id=run_id,
        actor_id=commit_user,
        tenant_id=tenant_id,
    ):
        herald.audit(
            event=event_name,
            **attrs,
            source_replay_id=sf_event["event"]["replayId"],
            source_commit_timestamp=commit_ts_ms,
            source_lag_observed_ms=lag_ms,
            source_change_type=change_type,
        )


def _translate(
    entity: str, change_type: str, payload: dict[str, Any]
) -> tuple[str, dict[str, Any]]:
    """
    Translate a Salesforce CDC entity + change_type into a chain event
    name + attribute dict. The mapping is documented in
    docs/sf-cdc-to-chain-event-mapping.md and reviewed quarterly.
    """
    if entity == "Case":
        if change_type == "CREATE":
            return "case.opened", {
                "case_id": payload["ChangeEventHeader"]["recordIds"][0],
                "subject": payload.get("Subject", ""),
                "origin": payload.get("Origin", ""),
                "priority": payload.get("Priority", ""),
            }
        if change_type == "UPDATE":
            changed = payload["ChangeEventHeader"].get("changedFields", [])
            if "OwnerId" in changed:
                return "case.assigned", {
                    "case_id": payload["ChangeEventHeader"]["recordIds"][0],
                    "assignee_id": payload["OwnerId"],
                    "queue": _resolve_queue(payload["OwnerId"]),
                    "assigned_at_utc": payload["LastModifiedDate"],
                }
            if "Status" in changed and payload.get("Status") == "Closed":
                return "case.closed", {
                    "case_id": payload["ChangeEventHeader"]["recordIds"][0],
                    "closure_reason": payload.get("Reason", ""),
                    "resolution": payload.get("Resolution__c", ""),
                    "closed_at_utc": payload["LastModifiedDate"],
                }
        # ... other Case change types
    if entity == "EmailMessage":
        return _translate_email(change_type, payload)
    if entity == "Account":
        return _translate_account(change_type, payload)
    raise ValueError(f"unmapped: entity={entity} change_type={change_type}")


def _resolve_queue(owner_id: str) -> str:
    # Resolves a Salesforce Group ID to its DeveloperName, cached.
    # Implementation in queue_resolver.py.
    ...


def _translate_email(change_type: str, payload: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    ...


def _translate_account(change_type: str, payload: dict[str, Any]) -> tuple[str, dict[str, Any]]:
    ...
```

Luis read it twice. He stopped at the `with herald.run(...)` block.

"OK," he said. "This is the part I want to understand. `herald.run(run_id=record_id)` — when this connector restarts mid-stream, how does the SDK know what `prev_hash` to use for the next event? You can't just pick up at seq=N if you don't remember N. And if the connector is wrong about N, the chain breaks."

Marcus nodded. He had been waiting for this question, too.

"Layered answer. The `with herald.run(...)` block, on entry, asks the SDK for the chain tail for `(tenant_id, run_id)`. The SDK checks three places in order."

He held up one finger.

"One. In-memory state. If the run is already open in this process, the tail is in memory — `(latest_seq, latest_payload_hash, key_version)`. The block reuses it. Cheap."

Two fingers.

"Two. Local persistence sidecar. A tiny SQLite file at `<state-dir>/<tenant_id>/<run_id>.state` carrying `(latest_seq, latest_payload_hash, key_version, last_commit_utc)`. Updated after every successful commit. File-locked for single-writer. If the process crashed mid-stream and restarted, the SDK reads the sidecar and resumes."

Three fingers.

"Three. Ledger query — the rejoin path. When local persistence is missing or corrupted — fresh container, lost disk, full DR scenario — the SDK calls Herald Core's ingestion API. `GET /chains/{tenant_id}/{run_id}/tail` returns the latest seq, payload_hash, key_version. Network access is required for this path; it's the fallback."

He paused.

"If none of the three find a tail, the run is genuinely new. Genesis: `seq=1`, `prev_hash=32 zero bytes`. The spec calls this the genesis-block anti-spoof — only a fresh run gets zero bytes."

Luis was writing.

Marcus continued.

"On flush — when the SDK ships a batch to the ledger — the ledger does a cross-check. The SDK declares its claimed `prev_hash` for the batch's first entry. The ledger compares against its own last-known `payload_hash` for `(tenant_id, run_id)`. Mismatch, ledger refuses with a named reason. Per §7 step 9 — the spec calls it `expected_prev_hash` discipline. The SDK can't unilaterally claim continuity; the ledger has to agree."

Karen looked up. "What if two connector processes try to write to the same run at once?"

"The SDK's file lock on the persistence sidecar fails the second one immediately. The application layer is responsible for one-process-per-run. The SDK enforces it locally so the second process can't even open the run, and the ledger enforces it remotely because the prev_hash disagreement at flush would surface."

Marcus pulled up the runtime configuration on the screen.

```python
herald.configure(
    tenant_id="northbridge-bank-prod",
    ikm_provider=herald.kms.AwsCloudHsmProvider(
        cluster="cluster-7y8a",
        key_id="k-northbridge-prod-2026q2",
    ),
    state_dir="/var/lib/herald/state",   # local persistence
    ledger_url="https://herald-core.northbridge-internal.com",
    rejoin_on_cold_start=True,           # fall back to ledger query if state_dir is empty
)
```

"That's the production config. `state_dir` is the sidecar location. `rejoin_on_cold_start=True` enables the third path — without it, a cold start with empty sidecar would refuse to write. We want rejoin in production because we run on ECS and tasks rotate; we leave it off in some test environments to make sure local persistence is exercised."

Luis asked the next question.

"What stops someone from silently restarting the chain at seq=1 to hide entries? If I'm a rogue actor and I want to disappear yesterday's bad assignment, can I just `rm -rf` the sidecar and let the SDK think it's a fresh run?"

"The genesis-block anti-spoof. A genuine new run has `prev_hash=zeros`. If you wipe the sidecar and try to silently restart an existing `run_id` from seq=1, the SDK on the rejoin path queries the ledger — and the ledger sees an existing tail for that `run_id`. The SDK's claimed `prev_hash=zeros` doesn't match the ledger's last-known `payload_hash`. Ledger refuses. Named reason."

"What if the rogue actor disables the rejoin path and just writes locally?"

"Then they have a chain on local disk that says seq=1, prev_hash=zeros. They flush to the ledger. The ledger's cross-check fires — the ledger has a real tail for that `run_id`, the SDK is claiming zero — refuse. The 'restart' shows up in the ledger's refusal log as a fork attempt. The verifier reports it as anomaly. The forensic trail is in two places: the ledger's refusal log and the gap in the legitimate chain where the rogue's entries should have continued."

Luis was nodding. He wrote: *genesis-block anti-spoof — fork shows up in two places.*

Karen asked: "Is the SDK source open?"

Marcus said, "Herald.Py is Apache 2.0 open source. The state-management code is in `herald/_buffer.py` and `herald/_runtime.py`. Anyone — your team included — can audit the resume logic. The repo is on GitHub."

Karen wrote that down. *SDK is OSS. Resume logic is auditable.*

She stared at the line for a moment.

*It never is*, she thought. *Except today, apparently, the SDK is also open source and the resume logic is in two named files.*

Chen took over.

"Multi-region setup?"

"Pattern A from spec §10.15 — multi-region active-active. Both regions write to local Herald Core. ETL reconciliation runs on a schedule, publishes a sealed `master.cross_region_replication_completed` event each batch. The reconciliation entry itself is in the chain."

"So the cross-region reconciliation is auditable as a chain entry."

"Yes."

Chen pulled up the most recent reconciliation event. Sealed. Verified. The reconciliation report metadata showed a delta of zero between regions for the previous 24 hours.

"Has there ever been a non-zero delta?"

"Twice. Once in February, once in March. Both were resolved within the reconciliation window. Both resolution events are sealed chain entries. I can pull them up if you want."

"Pull up the February one."

The February delta showed up: 3 events, all from a region failover during a maintenance window, all eventually replicated, all reconciled. The reconciliation event chained to the original entries.

> ### ✓ Confirmation #8 — Cross-region reconciliation is itself a sealed event
>
> The bank runs Pattern A multi-region active-active. ETL reconciliation events are sealed. Historical reconciliation deltas (including non-zero deltas) are themselves chained and reviewable. Chen verified the February event end-to-end.

Chen paused.

"Have you had any actual data-integrity events this year that weren't reconciliation deltas?"

"One," Marcus said. "March 17. A connector retry storm produced duplicate Salesforce mirror events. The deduplication ran inside the chain — every duplicate was captured, every dedup decision was a sealed event. The audit trail of the dedup is in the chain. No data was lost. No data was silently dropped. I have the incident report."

He slid a folder across the table. Karen flipped it open. The incident report referenced 14 sealed chain entries by `entry_id`. She picked one at random and ran the verifier.

PASS. Step: 12. 1.2 seconds.

She closed the folder.

Chen had one last question.

"Backup integrity. The chain is in the database. The seal records are in object storage. What about backups of the database — are *those* themselves auditable?"

"Backups are written to a separate object storage tier with the same compliance-mode lock. Each backup completion is a sealed chain entry. The chain entry includes a hash of the backup artifact. Restoring from a backup that doesn't match its hash fails — the restore tool refuses to load a backup whose chain entry doesn't validate."

"So a tampered backup is detectable on restore."

"It's detectable before restore. The hash check happens before the restore tool will read past the artifact header."

Chen nodded slowly. "I want that pattern in my notes."

---

## 📊 3:00 PM — Reconciliation Test

Karen wanted to do the reconciliation test herself. This was her usual test, the one she always ran on chain-of-custody systems, and it was where most systems folded.

She picked a sample window — 1,000 customer interactions across a single business day from the prior quarter. She asked Marcus for two things:

1. The operational-system view of those interactions (Salesforce, core-banking, voice transcription, AI advisor outputs).
2. The chain-of-custody view of those same interactions.

Marcus pulled both. The team spent forty minutes diffing them.

The diff returned zero.

Not "zero meaningful." Zero. Every event in the operational view had a sealed chain entry. Every sealed chain entry had a corresponding operational event. Timestamps matched within the documented capture latency. Payloads matched byte-for-byte after JCS canonicalization.

Karen ran the diff again with a different sample. 5,000 events this time, randomly selected from across the prior twelve months.

Zero.

She ran it a third time with a sample from a known-noisy day (March 17, the connector retry storm).

Zero — once the dedup events were factored in. The dedup events themselves were in both views.

She ran a fourth diff. This one she didn't tell anyone about. She picked a sample from the day immediately before the prior-year MRA closed — a day she knew, from the closing report, had been operationally tense.

Zero.

*The system worked under the eyes of FDIC examiners closing an MRA*, she thought. *That is a non-trivial test environment.*

She looked up at Marcus.

"What was your false-positive rate during the MRA close?"

"On the chain side, zero. On the operational side, we had two near-misses where Salesforce reporting and the chain disagreed momentarily during the connector lag window. Both reconciled within minutes. Both reconciliations are in the chain."

"Did you tell the FDIC?"

"I told the FDIC. I showed them the chain entries for the disagreement and the reconciliation. They closed the MRA on time."

Karen wrote: *FDIC saw the lag window during MRA close. Closed anyway. The Nit I'm about to write is not new information to the regulator.*

> ### ✓ Confirmation #9 — Operational and chain views reconcile to zero
>
> Three independent samples (1,000 events, 5,000 events, and a known-noisy-day sample) reconciled byte-for-byte between the operational system view and the Herald Core chain. Latency offsets were within the documented capture window. Dedup events were visible and traceable.

Karen put her pen down.

"Marcus, when did this go in?"

"Eighteen months ago."

"Was the prior-year MRA the trigger?"

"The MRA was the trigger. The procurement was already in flight. The MRA accelerated it by about a quarter."

Karen made a note. *NB-prior-MRA closed cleanly. This is the verification revisit. Verification holds.*

---

## 😬 3:45 PM — The Friction Builds (In a New Direction)

Karen wanted to push harder. She had a half-formed sense that something was off — not because she had found anything, but because she hadn't found anything, and her professional instinct was that this was the time things broke.

"Marcus, can you pull in the SRE on-call? I want to watch a seal happen live."

"Greg's already in the building. Hang on."

Greg came in. He was wearing a fleece pullover with a coffee stain on it. He nodded at the team and sat down at the second screen without much ceremony.

"What do you want to see?"

"Live seal. Today's batch. From the moment the seal job kicks off to the moment the verifier returns PASS."

Greg shrugged. "It runs at 4:15 PM Eastern. We can wait twenty minutes or I can trigger a manual seal of the partial day. Manual seal lands the same, just on a partial-day Merkle root."

"Manual seal."

Greg hit two keys. A job kicked off. The team watched the log stream.

```
[15:46:02] seal job started: tenant=northbridge, partial=true
[15:46:02] gathering chain entries: 1,847,392 leaves
[15:46:04] computing Merkle tree: depth 22
[15:46:05] requesting Ed25519 signature from CloudHSM
[15:46:05] signature received, fingerprint=7f3a9...
[15:46:05] writing seal record to Herald Core
[15:46:05] seal record written, entry_id=ce_8b1c...
[15:46:05] seal job complete: duration=3.1s
```

Greg said, "Verifier?"

Karen ran:

```
herald-verify --tenant=northbridge --seal-id=ce_8b1c... --strict
```

PASS. Step: 12. 0.6 seconds.

Greg stood up.

"Anything else?"

"No," Karen said.

"Cool."

He walked out.

> ### ✓ Confirmation #10 — Live seal demonstrated end-to-end in under 4 seconds
>
> Manual seal job kicked off, completed, and verified during the engagement window. The SRE on-call demonstrated the workflow without ceremony. CloudHSM signature acquired, Merkle root sealed, verifier returned PASS in under one second after seal completion.

Tom looked at Marcus.

"Greg has done this before."

"Greg has done this for the FDIC examiners three times this year."

Tom nodded. He wasn't sweating. He was, by his own internal measurement, pleased.

---

## 🔍 4:30 PM — Final Stress Test

Karen wanted to break it. Not because she thought she could, but because she wanted to know what it felt like to try.

"Pick me ten random entries. Across eighteen months. Different event classes. Different regions. Different customers."

Marcus typed. A list appeared.

```
ce_3a8f1d... (2024-11-04, voice transcription, region=us-east-1)
ce_7c2b91... (2025-02-19, AI advisor, region=us-west-2)
ce_91e8a4... (2025-06-30, IAM grant, region=us-east-1)
ce_4f29d8... (2026-04-12, API call, region=us-west-2)
ce_2d7b6c... (2025-09-11, loan decisioning, region=us-east-1)
ce_8a3f72... (2025-12-22, branch tablet, region=us-west-2)
ce_5e1a08... (2025-04-08, CRM mirror, region=us-east-1)
ce_b27c9f... (2024-08-15, voice transcription, region=us-west-2)
ce_6d9a31... (2026-01-03, AI advisor, region=us-east-1)
ce_f04e8b... (2025-11-27, API call, region=us-west-2)
```

Karen ran the verifier on each one.

```
herald-verify --entry-id=ce_3a8f1d... --strict
```
PASS. 4.0s.

```
herald-verify --entry-id=ce_7c2b91... --strict
```
PASS. 3.8s.

```
herald-verify --entry-id=ce_91e8a4... --strict
```
PASS. 4.2s.

She kept going. She ran all ten.

Ten passes. Average 4 seconds. The longest one was 4.4 seconds (the 2024 entry, which had to walk further back in the chain).

Karen raised an eyebrow. Genuine, not theatrical.

She picked an eleventh, off-script. A random entry from the day of a known incident — March 17, 2026, the connector retry storm.

PASS. 4.1s.

She picked a twelfth. A signing-key-rotation boundary. Q1 to Q2 of last year, one entry on each side.

PASS. PASS. The verifier handled the key rotation transparently — both entries verified against their respective signing-key fingerprints, with the rotation event itself being a sealed chain entry that linked the two key periods.

> ### ✓ Confirmation #11 — Verifier handles signing-key rotations transparently
>
> Quarterly key rotation events are themselves sealed chain entries. The verifier resolves the correct signing-key fingerprint per entry based on seal-date metadata. Cross-rotation verification works without any operator intervention. Karen sampled both sides of a Q1→Q2 rotation boundary; both passed.

She picked a thirteenth — a deliberately torturous one. An entry from a tenant-binding label that she couldn't find in the public registry.

```
herald-verify --entry-id=ce_x... --strict
```

```
Status: FAIL
Step: 4
Exit: 1
Reason: tenant binding label not resolvable.
        Procedure could not begin.
```

Karen looked up.

"Why exit 1 and not exit 3?"

"Exit 1 is procedure-could-not-begin," Marcus said. "Exit 3 is chain-anomaly. The verifier distinguishes between 'I cannot start because something upstream is wrong' and 'I started and found a chain inconsistency.' This entry was from a deprecated test tenant from 2024. The label was retired."

"Show me a real exit 3."

"I'd have to corrupt an entry. I'd rather not corrupt an entry."

"Fair. Show me the test fixture."

Marcus pulled up the spec test vectors. The exit-3 fixture was there. Karen ran the verifier against the fixture.

```
Status: FAIL
Step: 7
Exit: 3
Reason: chain anomaly detected. prev_hash mismatch
        at entry_id=ce_test_corrupt_002.
```

> ### ✓ Confirmation #12 — Verifier exit codes are meaningfully distinct
>
> Exit 0 (PASS), exit 1 (procedure-could-not-begin), exit 2 (procedure-began-and-failed), exit 3 (chain-anomaly) are all reachable and meaningfully distinct. Karen exercised exit 0, exit 1 against a deprecated-tenant entry, and exit 3 against the spec test vector.

She closed the laptop.

She opened it again.

"One more," she said. "I want to verify a seal record on a laptop with no Northbridge credentials at all. Not even the read-scope ones."

She switched to her personal laptop. She pulled up the Herald.Compliance public page. She copied the published Ed25519 public-key fingerprint. She pulled down a seal record from the same page — the bank's documentation said this surface was unprivileged-readable for any seal older than 24 hours, and she picked one from the prior week.

She ran the standalone verifier locally.

```
herald-verify-standalone --seal-file=northbridge-2026-04-30.seal \
                         --pubkey=7f3a9...
```

```
Status: PASS
Step: 12
Reason: Merkle root matches signed seal,
        signature verified against provided public key
        (fingerprint 7f3a9...)
Elapsed: 2.4s
```

She closed her personal laptop.

"That's the property I needed to see. The chain verifies without us trusting Northbridge at all. We trust the public key on the Compliance page, and we trust the open-source verifier we ran. Everything else is mathematics."

> ### ✓ Confirmation #13 — Seal verification works with zero Northbridge-side trust
>
> Karen ran the standalone verifier on her personal laptop using only the published Ed25519 public-key fingerprint and a seal record pulled from the public Herald.Compliance surface. Verification passed in 2.4 seconds. No Northbridge credentials were used at any layer of the verification path. This is the assurance property that makes the system useful to a regulator who has not personally inspected the bank's infrastructure.

Tom had been watching the standalone-verifier run from the next chair. He had a question.

"How is `herald-verify` distributed?" he said. "Where did Karen pull that binary from? If a future examiner is going to download it cold, what does that path look like?"

Marcus said, "Separate repo from the spec. `github.com/<vendor>/herald-verify`. Go binary, Apache 2.0, reproducible builds. Each release ships Linux, macOS, and Windows binaries plus a manifest of SHA-256 and SHA-512 hashes, Cosign signatures tied to a published public key — sigstore.dev — a CycloneDX-format SBOM, and a source tarball. An examiner downloads from GitHub Releases, verifies the Cosign signature against the published key, runs the binary. No connection to Northbridge required at any stage."

"Where's the spec repo, then?"

"`github.com/ffiec-chain-spec/spec`. That repo references `herald-verify` as the reference implementation in the §11 References section. It pins a specific verifier version per spec version — v1.0b pins to herald-verify v2.x. Anyone else can write a competing verifier in any language; the spec is what they conform to. The test-vector corpus in the spec repo is the conformance harness."

Tom was writing. He underlined something.

*Verifier is OSS, separate repo, Cosign-signed releases.*

Karen asked the question Tom had been about to ask.

"Why isn't the verifier in the spec repo? Wouldn't that be simpler — one place, one download?"

Marcus had four reasons. He gave them in order.

"One. Spec is locked. Verifiers need bug fixes, performance improvements, new platform builds. Putting the verifier in the spec repo means every verifier patch becomes a spec-repo commit, and that pollutes the change history of a regulatory artifact."

"Two. Spec is a regulatory artifact — it's licensed CC-BY or public domain depending on the jurisdiction. The verifier is engineering software, Apache 2.0. Different licenses, different review processes, different communities. Mixing them confuses the contribution model."

"Three. When the spec eventually transfers to FFIEC's actual repo — and that transfer is the explicit goal — the federal regulator should not inherit a Go-binary maintenance burden. They get a spec. The reference implementation continues to live with the engineering community."

"Four. An examiner downloading from a GitHub Releases page with Cosign signatures tied to a published key is a stronger trust statement than downloading from a spec repo. The release-signing path is purpose-built for binary trust. The spec repo's signing path isn't. We picked the right channel for each artifact."

Karen nodded slowly. *That's the right separation.*

She thought about it for another beat.

"You're saying the verifier and the spec are co-developed by the same community but distributed through different trust channels because the trust models for 'specification' and 'binary you'll run on your audit laptop' are not the same."

"Correct."

"That's the right answer."

Tom finished writing.

> ### ✓ Confirmation #14 — Verifier is OSS, distributed separately, signed at release
>
> `herald-verify` lives in its own GitHub repo, Apache 2.0 licensed, reproducible builds, Cosign-signed releases tied to a published public key. The examiner's trust path is download → cosign-verify → run. No Northbridge-side credential, no spec-repo dependency. The spec at `github.com/ffiec-chain-spec/spec` references the verifier as the reference implementation but does not bundle it. Karen and Tom independently agreed on the rationale for the separation — different licenses, different change cadences, different trust channels, and a clean handover path when the spec transfers to FFIEC's repo.

"Tom," she said. "Are we done?"

"We're done."

---

## 🌆 5:30 PM — Auditor Debrief

The team gathered in the engagement room. Marcus had stepped out to give them privacy. Tom closed the door.

Karen wrote on the whiteboard.

```
Gaps:     0
Partials: 0
Nits:     1
```

Under Nits, she wrote:

> **Nit-001: Salesforce mirror lag wording.** The Salesforce-to-Herald.Py mirror is captured via a connector. The mirror lag is occasionally up to 90 seconds during peak Salesforce load. The runbook documentation says "near real-time" without quantifying it. Recommend updating documentation to specify a 95th-percentile lag bound and an alerting threshold.

> ### ⚠️ Nit-001 — Documentation precision
>
> The Salesforce SaaS edge is captured via a Herald.Py mirror connector. The connector is reliable; the seal coverage is complete; the reconciliation diff was zero. The wording "near real-time" in the runbook is imprecise. Replace with a quantified bound (e.g., "95th-percentile lag under 90 seconds; alerting fires above 120 seconds").

She turned around.

"Anyone want to add anything?"

Raj said, "I want to go on record that I have never finished a chain-of-custody audit 30% under budget time."

Diana said, "I want to go on record that the IAM events being chain-captured solves a class of problems I usually have to escalate."

Mike said, "The AI advisor failing closed when capture fails is the part I'm taking back to other engagements as a benchmark."

Luis said, "Object-lock at the storage tier with a separate trust boundary. That's the answer when anyone asks me what 'real append-only' looks like."

Chen said, "Cross-region reconciliation as a sealed event. I'm stealing that pattern."

Elena said, "The mirror lag Nit is the only thing I have. The Salesforce side itself is fine. The mirror itself is fine. The wording is fine *to a Northbridge engineer who knows what 'near real-time' means in their ops context.* It is not fine *to an examiner who has never seen the system before.*"

Tom said, "I told Marcus we'd have a draft report to him by end of day tomorrow. He said no rush. The CAE function here is staffed for this. I appreciated that."

Tom added, "He also asked me whether the report was something he could share with his board's risk committee verbatim or whether he'd need to summarize it. I told him verbatim is fine. The verifier outputs speak for themselves. He was visibly relieved — apparently last quarter's vendor risk review of a different system left them with a 40-page document the committee couldn't follow."

Karen looked at the whiteboard.

"Last week I wrote a report with twelve Gaps and four Material Findings. This week I'm writing a report with zero and a Nit."

She paused.

"I want to be careful in the report not to sound like a brochure. State the facts. Show the verifier output. Note the Nit. The bank knows what it has. Our job is to confirm it, not to celebrate it."

Tom nodded.

"One more thing," Karen said. "When the FFIEC examiners come back next year, this report should still be useful to them. I want the workpapers to include the verifier outputs we collected. Marcus already pulled a SOC 2 evidence pack for us — let's reference it as Appendix A. Spec version v1.0a. Public-key fingerprint as of engagement date. Sample entry IDs. The Nit. That's the report."

Karen capped her marker.

"One last thing for the workpapers. Tom, I want a paragraph in the cover memo about what 'verification revisit' meant in this engagement. Specifically: the prior-year MRA closed cleanly two quarters ago. This engagement was scoped to confirm the close held. It held. The control environment we examined today is materially the same as the one the MRA close report described. That is what we wrote down. We did not find new control degradation. We did not find drift."

Tom wrote.

"And Karen?"

"Mm."

"Do we want to flag for the FFIEC reviewer next year that we ran our verifier independently — that the public-key fingerprint and seal records reconciled without Northbridge-side privileges?"

"Yes. Specifically yes. That is the part of the assurance posture I want a future examiner to understand without having to ask. The bank could be hostile or compromised at the operational layer and the chain would still verify. The reviewer should know we tested that property by exercising it."

Tom nodded.

The team packed up.

Raj stopped at the door.

"You owe me two coffees."

"I'll buy you four," Karen said. "You earned the headache."

Raj said, "I didn't have a headache today."

"That's the headache."

He almost smiled. He turned and walked out.

---

## ❌ What They Expected vs ✅ What They Found

**❌ What They Expected (based on last week, and on the prior-year MRA history):**

- Operational logs and chain-of-custody logs would diverge under sampling.
- The chain table would have a hidden update path, somewhere.
- IAM elevation would have a break-glass account that bypassed audit.
- "Append-only" would mean "by convention," not "by storage policy."
- Cross-region replication would have un-reconciled deltas.
- The AI advisor surface would be the weak seam — capture would be best-effort.
- Verifier credentials would require Northbridge-side privileges.
- Engineers would be pulled in to make excuses about edge cases.
- The day would end with at least three Partials and one Material Finding.

**✅ What They Found:**

- Operational and chain views reconciled to zero across three independent samples.
- The chain table has no UPDATE or DELETE role in any environment, including production-DBA.
- IAM elevation is itself chain-captured. Auto-revocation is chain-driven, not cron-driven. Failures fail closed to the unprivileged baseline.
- Append-only is enforced by S3 object lock in compliance mode, in a separate trust boundary from the application account. Neither account root can bypass it within the retention window.
- Cross-region reconciliation publishes a sealed event each batch. Historical non-zero deltas (Feb, Mar) are themselves chained and reviewable.
- The AI advisor wrapper captures synchronously. The customer-facing surface fails closed when capture fails. Customers see a soft error rather than an un-audited recommendation.
- The verifier's design is unprivileged. The Ed25519 public key is published. Verification works on a coffee-shop wifi.
- The SRE on-call demonstrated a live seal in 3.1 seconds without ceremony. He had done it before.
- The day ended with zero Gaps, zero Partials, and one Nit about documentation wording.

---

## 🧾 Final Assessment Theme

> "The organization can demonstrate, byte-for-byte, that customer interaction data is complete, accurate, and unaltered."
