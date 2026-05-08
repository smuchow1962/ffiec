# Round 11 — Cryptographic Engineer Review

## Reviewer

**Dr. Adaeze Okonkwo**
Principal Cryptographic Engineer, 17 years implementing FIPS-validated HMAC/HKDF/Ed25519 systems for African Tier-1 banks and the IMF. Reviewer for ETSI TC CYBER. Author of the ITU-T X.509 chain integrity advisory.

**Reading angle.** Cryptographic substance — HMAC + HKDF construction, per-tenant binding, fingerprint-before-MAC ordering, canonical-form exclusion, cross-implementation reproducibility, IKM custody and registry retention, algorithm rotation, dual-algorithm transitional-period verifier semantics, negative-test-vector completeness.

**Read scope.**

- `spec/chain-of-custody-v1.md` (§3, §4.1, §4.1.1, §4.2, §4.3, §4.4, §5, §6, §7, §10)
- `docs/design/02-chain-construction.md`
- `docs/design/03-merkle-seal.md` (§3.7)
- `docs/design/04-hsm-custody.md` (§3.2.2)
- `docs/design/07-verifier-design.md` (§4.3)
- `docs/design/08-test-vectors.md`
- `spec/test-vectors/README.md`, `chain_vectors.json`
- `spec/test-vectors/negative/README.md` and the 16 per-case `description.md` files
- `spec/test-vectors/{001,002,010}-*/description.md`

I have no prior exposure to this codebase or to any earlier reviewer round. My findings are grounded in the shipped artifacts as I read them.

---

## Q1. Is the dual-algorithm transitional-period verifier semantics in spec §7 step 11 cryptographically sound?

**What I read.** Spec §7 step 11 (lines 319–326) defines four verifier behaviors when an institution operates Ed25519 + a post-quantum algorithm coexisting per §4.3.2:

1. **Both signatures present and both valid** → PASS, "co-signed integrity under both algorithm assumptions."
2. **Single algorithm signature on a seal during dual-algorithm posture** → PASS-WITH-ANOMALY (`partial-coverage seal`), explicitly framed as a *control-completeness* finding, NOT a chain-integrity finding.
3. **Algorithm not on the institution's declared posture list** → `--strict`: FAIL; non-strict: PASS-WITH-ANOMALY. Posture list resolution is delegated to institution-published configuration, with the institution's control description naming the file or registry entry.
4. **Single-algorithm posture (default for v1.0)** → reduces to single-algorithm verification.

The seal record schema in §4.2 carries `algorithm` as a single string. The signed payload in §4.3 binds `algorithm` inside the signature. Spec §7 step 11 also requires a precise `algorithm/key-type mismatch at signature verification` reason when the seal's claimed algorithm and the resolved public key's algorithm disagree.

**Cryptographic assessment.**

- The decomposition of "missing co-signature" as a *control-completeness anomaly* rather than a *chain-integrity failure* is correct. A seal with only an Ed25519 signature during declared dual posture is still cryptographically valid under Ed25519 — the integrity claim under that algorithm assumption holds. Treating it as a control-completeness issue (the institution's posture commitment was not honored on this seal-day) is the right separation of concerns. It cleanly avoids the JWT `alg=none`-style trap where verifiers conflate "the integrity property doesn't hold" with "the institution didn't follow its own policy."
- The seal record's `algorithm` field being a single string while a dual-algorithm seal carries two co-signatures is the only structural ambiguity I see. The spec text reads as if dual posture produces *two seal records* (each with its own `algorithm` and `signature`) for the same `seal_date`/`tenant_id` — which is the only encoding that makes sense given the §4.2 schema. But the schema does not say so explicitly. A v1.1 implementer could reasonably read "dual-algorithm posture" as "extend the schema with a `signatures: list of {algorithm, signature, public_key_id}`" and produce a different on-the-wire shape than another implementer who reads it as "two parallel seal rows."
- The `algorithm` binding inside `sign_payload` (per §4.3) closes algorithm-confusion the way it should — an Ed25519 signature cannot be presented as a Dilithium signature for a coincidentally-matching key, because the signed bytes include the algorithm identifier. This is best practice since 2018 and the placement is correct.
- The "algorithm not on declared posture list" check is delegated to institution-side configuration. The spec correctly notes this is out-of-scope for the spec proper but requires the institution's control description to name the file or registry entry. The verifier's implementation cost is low (parse a YAML/JSON list) and the audit trail is clean (`--strict` failure with seal-date attribution).

**Status: Partial.**

The dual-algorithm dispatch *semantics* are sound. The *seal encoding* during dual-algorithm posture is underspecified: spec §4.2 schema admits exactly one `algorithm` and one `signature` per seal record, while spec §7 step 11 contemplates "both signatures present" without saying how they coexist on the wire. Two conforming v1.1 implementers could produce different shapes. Concrete guidance — either "dual posture produces two seal records per (`tenant_id`, `seal_date`), one per algorithm" or "extend §4.2 with a `signatures` array" — closes the ambiguity before the first PQ-capable implementation ships.

This is forward-looking, not blocking for v1.0 (where §7 step 11 paragraph 4 explicitly says single-algorithm posture is the default). But it is a normative-text gap that v1.0 should fix or explicitly defer to v1.1.

---

## Q2. Is the negative-test-vector corpus complete and load-bearing for the rework defenses?

**What I read.** `spec/test-vectors/negative/` contains 16 cases (N001–N016), each with a `description.md` documenting tampering recipe, expected step number, expected reason string, and provenance. The README ranks N006 (flipped fingerprint), N013 (mid-write truncation), and N014 (botched rotation) as "the rework's signature contributions."

**Per-case verification (cryptographic substance).**

I walked each description and verified the tampering recipe maps cleanly to the spec §7 step it claims to exercise:

| Case | Step | Tampering | Reason ordering | Verdict |
|---|---|---|---|---|
| N001 | 9 | bit-flip in `payload_hash` | MAC mismatch after fingerprint passes | OK |
| N002 | 6 | swap two events without re-chaining | structural `prev_hash` walk catches | OK |
| N003 | 10 | seal.merkle_root garbage | streaming Merkle recompute mismatch | OK |
| N004 | 11 | signature replaced with garbage | Ed25519 verify fails | OK |
| N005 | 11 | signature for wrong tenant | `tenant_id` is in `sign_payload` so signature fails | OK and correct (this is the cross-tenant signature replay defence) |
| N006 | 8 | `key_fingerprint` flipped to arbitrary 16 bytes | recompute fingerprint from looked-up IKM, constant-time compare fails BEFORE any MAC compute | OK and load-bearing |
| N007 | 7 | `key_version` set to a generation absent from registry | IKM lookup miss, NO MAC compute | OK |
| N008 | 5 | entry's `format_version = "v2"` | per-entry format_version check after step 4 | OK |
| N009 | 1 | header's `format_version = "v2"` | most-specific-first refusal at file pre-flight | OK |
| N010 | 2 | header's `hkdf_inputs_digest` flipped | recompute and constant-time compare | OK |
| N011 | 3 | header's `genesis_hash` set non-zero | constant compare against 32 zero bytes | OK |
| N012 | 4 | event's `tenant_id` differs from header's | cross-chain-lift defence catches before MAC | OK |
| N013 | pre-flight | audit file's last byte not `\n` | file-level refusal before steps 1–12 | OK |
| N014 | 8 | botched rotation: `key_version=1` re-used for different IKM (same tenant) | fingerprint check catches it; the load-bearing operational case | OK and load-bearing |
| N015 | 6 | `prev_hash` substituted, `payload_hash` left | structural walk catches | OK |
| N016 | 6 | `prev_hash` substituted AND `payload_hash` recomputed by an attacker without the IKM | structural walk catches the substitution; the spec §4.1 inviolate property #8 (verifier feeds `expected_prev_hash`, not `entry.prev_hash`, into MAC recompute) closes the latent footgun even if the structural check were relaxed | OK and forward-looking |

**N006 description quality (load-bearing case).** The description correctly forbids the verifier from running step 9 after step 8 mismatch and explicitly states that a verifier reporting `payload_hash MAC mismatch` instead of `key_fingerprint mismatch` is non-conforming for the rework. The provenance section names the future-maintainer regression it exists to prevent ("a future maintainer reorders steps 8 and 9 for performance reasons") with the right cost analysis (~100ns vs ~1µs). This is exemplary defense-in-depth documentation.

**N013 description quality.** Correctly mandates the verifier MUST NOT silently truncate to the last `\n` and proceed. Two test files are required (`audit-truncated-no-newline.ndjson` and `audit-truncated-mid-event.ndjson`) which exercise both "no trailing newline" and "mid-event truncation" failure modes. Maps to IR Scenario 9 (`audit_file.truncation_detected`) for operational disposition. The separation of "verifier refuses" from "operational recovery" is clean.

**N014 description quality.** Correctly identifies the operational scenario (KMS/HSM `key_version` re-use after deletion). The description names the auditor's investigation path (IKM roster, not chain content) — the load-bearing reason for the precise reason-string distinction. The simulation walks through `ikm_lookup → fingerprint recompute → constant-time compare → FAIL` step-by-step, which is the right level of pseudocode for a normative test vector.

**Gaps in the corpus I would add for v1.1 (not blocking for v1.0).**

- **N017 (algorithm-confusion across PQ posture).** Once v1.1 ships dual-algorithm, a negative case where a Dilithium signature is presented in a seal claiming `algorithm="ed25519"` (or vice versa) — verifier MUST report `algorithm/key-type mismatch at signature verification`. The spec §7 step 11 already mandates this reason; the corpus has no case to lock the implementation.
- **N018 (sign_payload cross-version replay).** A v1 signature presented in a seal record marked `format_version="v2"` (or `spec_version="v1.1"`). The signed payload's `ffiec.chain-of-custody.v1\n` prefix and `format_version` line make this attack class structurally impossible, but a negative test that asserts "verifier refuses with sign_payload mismatch when seal claims a different format_version than the signature was made for" would lock the property.
- **N019 (canonical-form drift via included field).** A canonical-form encoder that accidentally includes one of the chain-stamp fields (e.g. `seq` or `kms_handle_uri`) in the JCS bytes would silently produce wrong `payload_hash` values that no verifier could catch by inspection. A test vector with two chain entries that differ only by a chain-stamp field but produce the same `payload_hash` would lock the canonical-form exclusion rule (spec §5).

**Status: Answered.**

The 16 cases are complete for the rework's signature defenses. N006/N013/N014 are correctly identified as load-bearing and their descriptions are exemplary. The three suggestions above are v1.1 hardening, not v1.0 gaps.

---

## Q3. Is the design 02 chain-construction document free of duplicate sections, and is its section numbering consistent?

**What I read.** I extracted every `##` and `###` heading from `docs/design/02-chain-construction.md`. Section list:

```
##  1. The contract
##  2. The canonical-form choice
        ### 2.1 Why JCS over JSON
        ### 2.2 Why not protobuf
        ### 2.3 What goes into the canonical payload
##  3. Persistence ordering — crash safety
##  4. The session-key handshake
        ### 4.0 Model A — IKM-delivered (most common)
        ### 4.0a Model B — Session-key-delivered
        ### 4.1 Handshake security floor
        ### 4.2 Memory protection per platform
        ### 4.3 Memory zeroisation honesty note
##  5. Performance budget
##  6. The hash-cache contract
##  7. Concurrency
##  8. What can go wrong, and how it surfaces
        ### 8.0 Writer-side (SDK) failures
        ### 8.1 Verifier-side failures (per spec §7 ordered procedure)
        ### 8.2 Mid-write truncation
##  8.1 Multi-process run semantics                    <-- top-level §8.1 collides with subsection ### 8.1
        ### 8.1.1 Per-process runs with parent linkage
        ### 8.1.2 Shared run_id across processes
        ### 8.1.3 DAG-shaped multi-process flows
##  9. Auditor's-lens review
        ### 9.1 Closed findings from the Herald upgrade audit
```

**Cryptographic substance.** §4 (the cryptographic body) contains no duplicate sections. §4.0, §4.0a, §4.1, §4.2, §4.3 are unique. The Model A / Model B exposition is clean. The HKDF inputs and per-tenant `key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` are stated identically in §1, §4.0, §4.0a, and the per-platform memory protection table — same byte sequence, no drift.

**Numbering anomaly.** The `## 8.1 Multi-process run semantics` heading at the top level collides with the `### 8.1 Verifier-side failures` subsection inside `## 8`. Two sections are addressable as "§8.1" within the same document. This is not a cryptographic substance issue, but it is an internal-reference hazard: a future cross-reference to "design 02 §8.1" is ambiguous. The fix is to renumber the multi-process section to `## 8.2` (and the existing `### 8.2 Mid-write truncation` to something like `### 8.0.2` or move it under §3 with the rest of the persistence material) OR promote `## 8.1 Multi-process run semantics` to `## 9` and renumber subsequent sections.

This was flagged for me as something round 10 found and fixed; the fix landed for §4 (which is clean) but the §8/§8.1 collision is still present.

**Status: Partial.**

Cryptographic content has no duplicates. The §4 cleanup is complete. The §8/§8.1 numbering collision remains and creates an internal-reference hazard. Renumber the top-level multi-process section to break the collision.

---

## Q4. Is design 04 §3.2.2 using the post-rework `key_version` / `key_fingerprint` vocabulary throughout?

**What I read.** `docs/design/04-hsm-custody.md` §3.2.2 ("IKM rotation window," lines 86–94).

The section uses the post-rework vocabulary consistently:

- "the institution rotates the IKM" (not "master key")
- "session keys derived from the old IKM"
- `key_version` (not the pre-rework `master_version`)
- `key_fingerprint` (not pre-rework `master_fingerprint`)
- spec cross-references are correct: "the entry's stamped `key_version` (spec §7 step 7)" and "the entry's recorded `key_fingerprint` (spec §7 step 8) before any MAC compute"
- The seal-day schema reference is correct: `key_versions` (plural list, e.g. `[3, 4]`)
- The cross-reference to spec §10.10 ("Rotation crossing the seal boundary") is correct

§3.2.2's terminology aligns with spec §3, spec §4.1, spec §7, and design 02. I cross-checked §3.1 (the IKM definition) for vocabulary consistency — the section opens with an explicit terminology bridge: "'IKM' replaces the older 'master key' term throughout the spec rework. The two names refer to the same thing... 'Master key' remains acceptable in operational shorthand; the spec text uses 'IKM' for precision." This is the right place to put the bridge and the right wording for an audience that may have seen pre-rework material.

**One observation (informational, not a finding).** Test vector case 014 in `docs/design/08-test-vectors.md` §5.4.3 still says "Master-rotation window" and uses `master_version = "v3,v4"` and `session_key_id` in its description — pre-rework vocabulary. The case `010-tenant-ikm-rotation-mid-day` in `spec/test-vectors/README.md` and `chain_vectors.json` correctly uses `key_version` and the rotation chain itself uses post-rework byte values. The case-014 prose in design 08 was not updated when the spec terminology was reworked. This is a documentation drift, not a cryptographic substance issue (the bytes in `chain_vectors.json` are right), but it would mislead a first-time implementer.

**Status: Answered.**

Design 04 §3.2.2 is post-rework-clean. The pre-rework terminology in design 08 §5.4.3 (case 014 description) is a documentation drift to clean up but does not affect cryptographic correctness — the byte-level fixture in `chain_vectors.json` uses the correct post-rework vocabulary.

---

## Q5. Does the design 07 verifier-design §4.3 sign_payload pseudocode include the `algorithm` line, and does the verification pseudocode dispatch on algorithm?

**What I read.** `docs/design/07-verifier-design.md` §4.3 (lines 215–236).

The pseudocode:

```
sign_payload = "ffiec.chain-of-custody.v1\n" ||
               seal.algorithm                || "\n" ||
               seal.format_version           || "\n" ||
               T                             || "\n" ||
               iso8601(D)                    || "\n" ||
               hex(seal.merkle_root)         || "\n" ||
               hex(seal.hkdf_inputs_digest)
# Algorithm dispatch: resolves the public key's algorithm from public_key_id.
# Mismatch between seal.algorithm and the public key's algorithm reports
# "algorithm/key-type mismatch at signature verification" (spec §7 step 11).
IF NOT verify_for_algorithm(seal.algorithm, public_key, sign_payload, seal.signature):
  FAIL: signature verification failed
```

Cross-checking against spec §4.3 (lines 174–192):

```
sign_payload = "ffiec.chain-of-custody.v1\n" ||
               algorithm                  || "\n" ||
               format_version             || "\n" ||
               tenant_id                  || "\n" ||
               iso8601_date(tenant_day)   || "\n" ||
               hex(merkle_root)           || "\n" ||
               hex(hkdf_inputs_digest)
```

Byte-for-byte, the verifier-design and spec sign_payload constructions match. The verifier-design comment also names the precise spec-§7-step-11 reason string for the algorithm/key-type mismatch case. The `verify_for_algorithm(seal.algorithm, ...)` dispatch matches the spec §7 step 11 mandate ("dispatching on `algorithm`").

I cross-checked against `chain_vectors.json` (`sign_payload_single_text`):

```
ffiec.chain-of-custody.v1
ed25519
v1
tenant-ffiec-test-1
2026-05-06
927adc88e5d843c5847674f7245d1e3c762bacc3f0b225f3322300cddf4eefe9
6f8a5005cabb2eab8b347254d0c94c2d585a7e5a5a82398c5aeaea074c727d65
```

This is exactly the byte sequence the verifier-design pseudocode produces from `algorithm="ed25519"`, `format_version="v1"`, `tenant_id="tenant-ffiec-test-1"`, `seal_date="2026-05-06"`, the published Merkle root, and the published HKDF inputs digest. Spec, design, and test vector are byte-for-byte aligned.

The only minor delta is that the spec §4.3 normative description names the inputs as `algorithm`, `format_version`, `tenant_id`, `iso8601_date(tenant_day)` — describing the abstract bytes — while the verifier-design pseudocode resolves them as `seal.algorithm`, `seal.format_version`, `T`, `iso8601(D)` from the seal record. That is correct verifier perspective: the verifier reads from the seal record, not from abstract inputs. The bytes that result are identical.

**Status: Answered.**

§4.3 includes the `algorithm` line. The verifier dispatch on `seal.algorithm` is documented. The construction matches spec §4.3 byte-for-byte and reproduces the test-vector `sign_payload_single_text` exactly.

---

## Q6. Is the algorithm field forensic-only note in spec §4.4 accurate, and does the security-decision boundary make sense?

**What I read.** Spec §4.4 (line 236) carries the per-entry `ffiec.chain.algorithm` attribute with this language:

> The HMAC algorithm used; default `"HMAC-SHA-256"`. Verifiers MUST handle when present and assume default when absent. **Forensic only** — the verifier does NOT trust this per-entry value for security decisions; the algorithm dispatch happens on the seal record's `algorithm` field at §7 step 11, and the `payload_hash` MAC compute uses the algorithm constant for v1 (HMAC-SHA-256). Parallel to `mac_computed_at_utc` and `kms_handle_uri` in this respect.

**Cryptographic assessment.** This is the correct security-decision boundary, and the wording is precise.

The per-entry `algorithm` attribute could otherwise be a JWT-`alg=none`-style attack surface: an attacker who can flip per-entry attributes on the wire could claim a chain entry was MAC'd under a weaker algorithm, and a credulous verifier might dispatch on the entry's own claimed algorithm. The spec correctly:

1. **Pins the per-entry MAC compute to the v1 constant (HMAC-SHA-256), not the entry's claim.** The verifier at spec §7 step 9 always recomputes with HMAC-SHA-256 — there is no algorithm dispatch at the per-entry level. An entry that claims `algorithm="HMAC-SHA-1"` and presents a SHA-1 MAC will fail the SHA-256 recompute regardless of what the attribute says.
2. **Routes algorithm dispatch to the seal record's `algorithm` field** (which is *signed* by the HSM as part of `sign_payload`). The seal-level algorithm cannot be flipped without the signature failing.
3. **Explicitly groups `ffiec.chain.algorithm` with `mac_computed_at_utc` and `kms_handle_uri` as forensic-only fields.** This is the right peer group — fields the verifier records but does not trust for security.

The canonical-form exclusion rule in spec §5 (line 265) correctly lists `algorithm` among the chain-stamp fields excluded from the JCS bytes. So the field is doubly contained: not in the MAC input, not trusted at MAC dispatch.

The "default `HMAC-SHA-256`" semantics ("verifiers MUST handle when present and assume default when absent") is appropriate for a v1.0 spec where only one HMAC algorithm is defined. A v1.1 that adds a second HMAC algorithm would need to elevate this attribute to security-relevant — but that elevation would also require binding the per-entry algorithm into the canonical bytes (or into a new per-entry signed envelope), which is a v1.1 design problem the spec correctly defers.

**Status: Answered.**

The forensic-only framing is cryptographically correct. The peer grouping with `mac_computed_at_utc` and `kms_handle_uri` is the right shape. The v1 constant pinning at MAC compute closes the algorithm-confusion attack at the per-entry level cleanly.

---

## Q7. Is the `tenant_id` grammar in spec §3 sufficient to prevent HKDF info-parameter boundary confusion?

**What I read.** Spec §3 (line 30):

> The tenant identifier (`tenant_id`) MUST match the regular expression `^[A-Za-z0-9_.\-]{1,255}$` — alphanumerics, underscore, hyphen, and dot only, with a maximum length of 255 bytes. The character set is constrained so that the HKDF `info` parameter `info = HKDF_INFO_BASE || "|" || utf8(tenant_id)` is unambiguously parseable: the `|` byte (0x7C) cannot appear inside `tenant_id`, eliminating boundary-confusion attacks where two distinct tenant identifiers could produce the same `info` byte sequence.

**Cryptographic assessment.**

The grammar is sufficient and the rationale is correct. The character class `[A-Za-z0-9_.\-]` excludes `|` (0x7C), which is the boundary byte between `HKDF_INFO_BASE` and `utf8(tenant_id)`. Without the exclusion, two tenants such as `acme` and `acme|extra` could in principle construct different `tenant_id` values that, when the boundary byte is appended, produce the same byte sequence as the constructor `HKDF_INFO_BASE || "|" || "acme|extra"` and `HKDF_INFO_BASE || "|" || "acme" + something` — the HKDF Expand step would derive the same session key for two distinct tenant_id values. The constraint correctly closes this.

The 1–255 byte length constraint is also right: HKDF info parameter is unbounded in RFC 5869, but the 255-byte cap is small enough to keep the per-tenant `info_for_tenant` under 286 bytes (`30 + 1 + 255`), which fits comfortably in the HKDF-SHA-256 single-block info path and avoids any encoder-specific length-prefix surprises.

ASCII-restricted character class also avoids UTF-8 normalization drift: an attacker cannot produce two distinct Unicode codepoint sequences (NFC vs NFD) that both UTF-8-encode to bytes the verifier might accept as "the same" tenant. Restricting to ASCII alphanumerics + `._-` makes the byte sequence unambiguous.

**One observation (informational).** The spec does not say where the grammar is enforced (SDK at configure-time? IKM-registration at provisioning-time? Both?). Section 10.6 (IKM minimum length) explicitly says "Implementations SHOULD enforce the minimum at IKM-provisioning time... and MUST enforce it at SDK-configure time." A parallel statement for the `tenant_id` grammar — enforcement at IKM-registration AND at SDK-configure-time — would close the door on a misconfigured deployment that registers a malformed tenant_id in the registry but the SDK then refuses to serve, leaving the registry in an inconsistent state.

**Status: Answered.**

The grammar is cryptographically sufficient and the boundary-confusion rationale is correctly stated. The enforcement-point ambiguity is an operational hardening note, not a cryptographic gap.

---

## Q8. Is the IKM-registry retention coupling in spec §10.9 sound, and does it close the cryptographic deletion-risk failure mode?

**What I read.** Spec §10.9 (lines 409–415):

> The tenant key registry MUST retain every IKM generation for at least as long as any chain entry stamped with that `key_version` is retained. An institution that retires an IKM out of recoverability while chain entries that reference it still exist loses the ability to verify those entries; the verifier reports `unknown key_version: no IKM for (tenant=T, key_version=V)` per §7 step 7 and the affected days FAIL key-bound verification.
>
> Implementations SHOULD enforce the retention coupling at the registry layer: a request to retire an IKM whose `key_version` is still referenced by retained chain entries MUST require an explicit override and MUST be logged as a `master_key.retired` operational event...
>
> Cloud KMS providers' default retention behavior varies (AWS CloudHSM keys can be marked for deletion with a 7-30-day pending-window; Azure Managed HSM and Google Cloud HSM offer similar windows). Institutions document the configured retention window per IKM generation; the conservative posture is to retain IKMs for the longer of (a) the retention period of any chain entry referencing them, and (b) the institution's regulatory minimum (typically 7 years for FFIEC chain-of-custody data).

**Cryptographic assessment.**

The retention coupling is *necessary and sufficient* to keep the chain verifiable. The spec correctly:

1. **Names the failure mode precisely** (`unknown key_version: no IKM for (tenant=T, key_version=V)` per §7 step 7) so a verifier output points the auditor at the registry, not the chain.
2. **Requires a registry-layer override** for retiring an IKM whose `key_version` is still referenced — this is the correct enforcement point because the registry is where the binding `(tenant_id, key_version) → IKM` lives.
3. **Mandates a `master_key.retired` operational event** so the audit trail records every IKM retirement explicitly. (The event name predates the IKM/master-key vocabulary rework but is grandfathered for backward compatibility — this is the right call.)
4. **Enumerates cloud-KMS pending-window behavior per provider** (AWS CloudHSM 7–30 days; Azure Managed HSM and GCP similar) so institutions know to configure the pending window to be at least as long as their chain-entry retention.
5. **Sets the conservative posture** as `max(chain_entry_retention, regulatory_minimum)`, with the typical 7-year FFIEC chain-of-custody floor.

This composes cleanly with §3.2.1 of design 04 (cloud HSM provisioning) and §3.2.2 of design 04 (IKM rotation window). The verifier's behavior at §7 step 7 (unknown key_version → FAIL) is documented and exercised by negative test N007.

The coupling correctly handles the late-binding case from §10.10: a tenant-day spanning a rotation has both the old and new `key_version` in the seal's `key_versions` list, and as long as both IKMs are retained per §10.9, the verifier walks both halves. The verifier needs both IKMs because the `payload_hash` for entries under the old `key_version` was MAC'd with the old session key, derived from the old IKM.

**Cryptographic forward-look.** The retention coupling also bounds the *post-quantum migration* problem cleanly. When v1.1 ships a PQ algorithm alongside Ed25519, the Merkle seal's `algorithm` field rotates but the per-entry HMAC-SHA-256 keeps using the same IKM-derived session keys. So the IKM retention requirement does not change with algorithm rotation — it remains keyed to chain-entry retention. This is the right separation of concerns; a less-careful spec might have entangled IKM retention with signing-key retention.

**Status: Answered.**

The retention coupling is cryptographically sound, operationally enforceable at the registry layer, and composes correctly with rotation, late-binding, and (forward-looking) algorithm rotation.

---

## Per-role roll-up

| # | Question | Status |
|---|---|---|
| Q1 | Dual-algorithm transitional-period verifier semantics in spec §7 step 11 | **Partial** — semantics sound, but seal-record encoding for the dual-signature case is underspecified |
| Q2 | Negative-test-vector completeness (16 cases, especially N006/N013/N014) | **Answered** |
| Q3 | Design 02 free of duplicate sections | **Partial** — §4 is clean (round 10 fix held); §8/§8.1 numbering collision remains |
| Q4 | Design 04 §3.2.2 post-rework vocabulary | **Answered** (with informational note: design 08 §5.4.3 case-014 prose still uses pre-rework vocabulary; bytes in `chain_vectors.json` are correct) |
| Q5 | Design 07 §4.3 sign_payload pseudocode includes `algorithm` line | **Answered** |
| Q6 | §4.4 algorithm field forensic-only framing | **Answered** |
| Q7 | tenant_id grammar prevents HKDF info-parameter boundary confusion | **Answered** |
| Q8 | IKM-registry retention coupling in §10.9 | **Answered** |

**Tally:** 6 Answered, 2 Partial, 0 Gap.

## Gaps and partials — what would close them

**Q1 (Partial → Answered).** Add one paragraph to spec §4.2 (seal record schema) or §7 step 11 explicitly specifying the on-the-wire shape of a dual-algorithm seal. Two clean choices:

- *Option A (preferred for backwards compatibility):* "When the institution operates dual-algorithm posture, the seal job produces two seal records for the same `(tenant_id, seal_date)` — one per algorithm. Each record carries its own `algorithm`, `signature`, and `public_key_id`. The verifier resolves both records from storage and applies the §7 step 11 dispatch."
- *Option B (cleaner schema, breaks v1.0 → v1.1 record shape):* "Extend §4.2 schema with `signatures: list of {algorithm, signature, public_key_id}`. Single-algorithm posture has a one-element list; dual-algorithm has two. The `algorithm` and `signature` top-level fields are deprecated in v1.1 in favor of the list form."

Option A is non-breaking and clearly the right choice for v1.0 forward-compatibility.

**Q3 (Partial → Answered).** In `docs/design/02-chain-construction.md`, renumber `## 8.1 Multi-process run semantics` to break the collision with `### 8.1 Verifier-side failures`. The simplest fix: promote the multi-process section to `## 9. Multi-process run semantics` and shift `## 9. Auditor's-lens review` to `## 10`. Internal cross-references inside the file would need a one-pass update; nothing outside the file references those numbers (I grep-checked the spec and other design docs).

---

## Cryptographic substance — overall assessment

The HMAC + HKDF + Ed25519 construction is cryptographically clean across spec, design, and test vectors. The load-bearing rework defenses (per-tenant HKDF binding via `info` parameter; per-entry `key_fingerprint` checked before any MAC compute; canonical-form exclusion of chain-stamp fields; verifier feeds `expected_prev_hash` into MAC recompute, not `entry.prev_hash`; algorithm binding inside `sign_payload`) are stated normatively in the spec, exercised by N006/N014/N015/N016 in the negative corpus, and reproduced byte-for-byte in `chain_vectors.json`.

The forward-look on post-quantum is well-positioned: the algorithm dispatch lives at the seal level (where it is signed), per-entry HMAC stays pinned to the v1 constant (closes algorithm-confusion at the per-entry level), and IKM retention is keyed to chain-entry retention rather than signing-key retention (so PQ algorithm rotation does not perturb the IKM lifecycle).

The two open partials are documentation-shape issues — the dual-algorithm seal encoding ambiguity is forward-looking (v1.0 default is single-algorithm so the ambiguity is dormant until v1.1), and the §8/§8.1 numbering collision is an internal-reference hazard rather than a cryptographic substance issue. Both close with localized text edits.

This corpus, as it stands, is implementable by an independent SDK or verifier and the byte-level outputs will match the published vectors. That is the test that matters.
