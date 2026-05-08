# 01 — Principal Cryptographic Engineer (round 9)

> **Role.** Cryptographic and key-management lens. First-look review of the v1.0-rework corpus.

## Persona — Dr. Yuki Tanaka

**Background.** Principal Cryptographic Engineer with 18 years operating HMAC and HKDF deployments at Tier-1 banks; three FIPS 140-2 L3 HSM integrations; active reviewer of IETF crypto drafts. I read this corpus the way I read a draft RFC — looking for the one underspecified byte that lets two implementations disagree, the one ordering choice that lets an attack class slip through, and the one custody assumption that cannot survive contact with a real bank operator.

**Reading angle.** The substance: HMAC + HKDF construction soundness, the per-tenant binding, the fingerprint-before-MAC ordering, the canonical-form exclusion rule, the cross-implementation reproducibility property, the IKM custody and rotation story, and the algorithm-rotation commitment. I am not looking at the operator narrative or the regulator framing; that is for the other personas.

### Questions

**Q1. Is the `hkdf_inputs_digest` definition consistent across the spec, design docs, seal record, file header, and test vectors? A single field with two definitions is a cross-implementation footgun.**

**Status:** Gap

The corpus has two different definitions of `hkdf_inputs_digest` and uses the same field name for both, which is the classic shape of a divergence trap.

- Spec §3 (the definition table): "`SHA-256(salt \|\| info_for_tenant \|\| length_LE32)`" — per-tenant, includes `info_for_tenant = HKDF_INFO_BASE \|\| "\|" \|\| utf8(tenant_id)`.
- Spec §4.4 (file-header attribute table): "`SHA-256(HKDF_SALT \|\| info_for_tenant \|\| length_LE32)` for the file's tenant" — per-tenant, matches §3.
- Spec §7 step 2 (verifier procedure): "`SHA-256(HKDF_SALT \|\| info_for_header_tenant \|\| length_LE32)`" — per-tenant.
- Verifier-design 07 §4.0 step 2: same per-tenant variant.
- Test-vector README and `chain_vectors.json` `hkdf_inputs_digest_hex = 6f8a5005cabb...` — computed as the per-tenant variant (verified via the README's recipe `SHA-256(HKDF_SALT \|\| info_for("tenant-ffiec-test-1") \|\| (32).to_bytes(4, "little"))`).

But:

- Spec §4.2 (seal-record schema): "`SHA-256(HKDF_SALT \|\| HKDF_INFO_BASE \|\| length_LE32)` for the chain construction in force on this day" — **tenant-free, no `|` separator, no tenant_id.**
- Design 03 §3.7: "`SHA-256(HKDF_SALT \|\| HKDF_INFO_BASE \|\| length_LE32)` for the day's chain construction. Same value on every seal for v1." — same tenant-free variant.
- Design 04 §4.6: discusses the field generically without naming the bytes hashed.

These cannot both be right. The seal record's signed `sign_payload` (§4.3) embeds `hex(hkdf_inputs_digest)`; if implementation A computes the per-tenant variant and implementation B computes the tenant-free variant, the verifier rejects every seal it did not produce locally. Two conforming implementations of v1.0 will fail to interoperate at the seal-verification step. The test vectors only cover the per-tenant variant, so an implementation that follows §4.2 verbatim passes its own seals but fails the corpus's expected seal payload (`sign_payload_single_hex` decodes to a payload whose last hex line is `6f8a5005...`, which is the per-tenant value).

The spec needs to pick one definition and propagate it everywhere. The corpus's choice (per-tenant) is the better one because it binds the seal more tightly — a seal accidentally built with the wrong tenant's HKDF inputs fails the digest check before signature compute. The spec §4.2 and design 03 §3.7 text needs to be corrected to match.

**Q2. Is the binary structure of the HKDF `info` parameter unambiguous given an unrestricted tenant_id character set?**

**Status:** Partial

Spec §4.1 defines `info_for_tenant = HKDF_INFO_BASE \|\| b"\|" \|\| utf8(tenant_id)` and §4.1.1 inviolate property #1 makes per-tenant binding via `info` normative. RFC 5869 §3.2 supports binding context into `info`, so the construction is sound at the primitive level.

The gap is at the encoding boundary. The separator is a single `|` byte (0x7C). Spec §3 defines `tenant_id` as a UTF-8 string and uses it freely throughout (`tenant_acme_prod`, `tenant-ffiec-test-1`). The spec does not constrain the tenant_id character set. A tenant_id that itself contains `|` (e.g., `tenant|prod` versus `tenant` with a literal pipe handled differently) produces a HKDF info value that is structurally indistinguishable from a different (tenant_id_prefix, suffix) split. This is the classic length-extension / boundary-confusion shape.

Concrete attack: tenant `acme` with `info = "ffiec.chain-of-custody.v1.info|acme"` and tenant `acme|x` with `info = "ffiec.chain-of-custody.v1.info|acme|x"` derive different keys (good), but a registry that allows both tenants creates an unstated equivalence class: `info_base|tenant_id` is not a prefix-free encoding. Two verifiers that disagree on whether `tenant_id` may contain `|` could disagree on whether two tenants are truly distinct.

Two paths close this:

1. Constrain `tenant_id` to a character set that excludes the separator byte (e.g., RFC 4648 base32-alphabet, or `[A-Za-z0-9_-]` per the Herald.Py reference). Spec §3 should add a normative grammar.
2. Length-prefix the `tenant_id` instead of using a separator: `info = HKDF_INFO_BASE \|\| utf8(tenant_id).length.to_bytes(4, "big") \|\| utf8(tenant_id)`. This is unambiguous regardless of tenant_id contents.

The spec does neither today. Implementations that follow option (1) but disagree on the allowed character set, or that follow option (2) instead of the literal `b"\|"`, produce incompatible HKDF outputs and the chain breaks at vendor boundaries. The conformance corpus uses tenant_id `tenant-ffiec-test-1` which dodges the question; an institution with a tenant_id containing `|` would discover the underspecification at deployment time.

**Q3. Is the algorithm identifier in the seal record bound into the signed payload, or can an attacker present an Ed25519-signed seal with a flipped algorithm identifier?**

**Status:** Partial

Spec §4.3.2 ("Algorithm rotation and quantum-readiness") makes the seal record carry an explicit `algorithm` field and requires the verifier to dispatch on it. The future-version commitment is reasonable: Dilithium (FIPS 204) and SLH-DSA (FIPS 205) are named.

The gap: the `algorithm` field is **not** included in `sign_payload` (spec §4.3). The signed payload is:

```
ffiec.chain-of-custody.v1\n
{format_version}\n
{tenant_id}\n
{ISO 8601 date}\n
{hex(merkle_root)}\n
{hex(hkdf_inputs_digest)}
```

The verifier dispatches to the algorithm named in the seal record's `algorithm` field, then verifies the signature against the public key resolved from `public_key_id`. In v1.0 with only Ed25519, a flipped `algorithm` field causes the verifier to look for a non-Ed25519 public key for `public_key_id` and fail with "signature verification failed" or a key-type mismatch. The chain holds. But the failure mode is generic; a more specific reason ("algorithm field does not match the public key's algorithm") would be preferable.

Once a second algorithm ships (say v1.1 adds Dilithium), the surface widens. An attacker with a valid Ed25519 signature over `sign_payload` cannot forge a Dilithium signature, so cross-algorithm forgery is not immediately possible. But algorithm-identifier confusion attacks have a long history (JWT `alg=none`, the various "alg=HS256 with public key as HMAC key" attacks, and SAML algorithm-substitution). Best practice since 2018 has been to bind the algorithm identifier into the signed payload.

The clean fix is one extra line in `sign_payload`:

```
ffiec.chain-of-custody.v1\n
{algorithm}\n                 ← added line
{format_version}\n
{tenant_id}\n
{ISO 8601 date}\n
{hex(merkle_root)}\n
{hex(hkdf_inputs_digest)}
```

This costs ~10 bytes per seal and removes the entire class of algorithm-confusion attacks before v1.1 ships a second algorithm. The spec's current text is defensible for v1.0 (only one algorithm exists), but the v1.1 algorithm-rotation commitment becomes harder to honour cleanly without breaking v1.0 seal-payload compatibility. Lift it now.

**Q4. Does the canonical-form exclusion rule produce byte-identical output across two independent v1.0 implementations for any logical event the spec admits?**

**Status:** Answered

Spec §5 names the exclusion rule plainly: "the canonical bytes input to the MAC MUST EXCLUDE every chain-stamp field" and lists the exact fields (`prev_hash`, `payload_hash`, `key_version`, `key_fingerprint`, `format_version`, `mac_computed_at_utc`, `kms_handle_uri`, `algorithm`, `seq`, plus the linkage fields). Spec §5 commits to RFC 8785 (JCS) for the canonical form. Design 02 §2.3 enumerates the included field set explicitly (OTel envelope, GenAI semconv attributes, application audit payload, explicit metadata).

The `chain_vectors.json` fixture provides byte-level evidence: each event's `event_canonical_hex` decodes to a JCS-canonical JSON object whose keys are sorted Unicode-codepoint-ascending and whose chain-stamp fields are absent. Spot-checking event 1's canonical bytes (decoded from the hex on line 25): the keys are `attributes, chain_kind, duration_ns, event_id, kind, name, parent_span_id, resource, run_id, severity, span_id, tenant_id, timestamp_ns, trace_id` — sorted, no `seq`, no `prev_hash`, no `payload_hash`, no `key_*`, no `format_version`, no `mac_computed_at_utc`, no `kms_handle_uri`. The MAC computed over `prev_hash || canonical` reproduces the published `payload_hash_hex`. The corpus enforces what the spec requires.

The reasoning in design 02 §2.3 ("Why the chain-stamp fields are excluded") is correct: the MAC is one of the chain-stamp fields, so including them in the canonical bytes would create a circular dependency. The other chain-stamp fields' tamper coverage is defended elsewhere (key_version → wrong IKM lookup → MAC fail; key_fingerprint → fingerprint check at step 8; format_version → header mismatch at step 5).

One minor observation, not a gap: the spec excludes `algorithm` from the canonical bytes (§5 last paragraph). This is correct per the design rationale, but it means an attacker can flip `algorithm` on a per-entry basis without breaking the chain entry's MAC. The verifier never reads per-entry `algorithm` for security purposes — it dispatches on the seal record's `algorithm` field at step 11. Worth a sentence in §4.4 making explicit that per-entry `algorithm` is forensic-only, parallel to the existing wording for `mac_computed_at_utc` and `kms_handle_uri`.

**Q5. Is the IKM custody and rotation story complete enough that two implementations following the spec arrive at the same operational behaviour during a rotation window?**

**Status:** Partial

The spec covers the cryptographic mechanics of rotation cleanly:

- Per-entry `key_version` and `key_fingerprint` resolve to the right IKM at verifier time (spec §4.1, §7 step 7-8).
- The seal record's `key_versions` is a list to capture mid-day rotation (spec §4.2, design 03 §3.7).
- Test vector 010 provides the byte-level rotation fixture (`rotation_chain` in `chain_vectors.json`); the verifier walks the rotation seamlessly.
- Negative test N014 (botched rotation: `key_version=1` re-used for a different IKM) is caught at the fingerprint step before any MAC compute.

What is incomplete is the operational picture during the rotation window:

1. **No retirement procedure for the old IKM.** Design 04 §3.1 says "Past entries (with the old `key_version` and old `key_fingerprint`) remain verifiable as long as the old IKM is retained in the tenant key registry." The spec does not say *how long* the old IKM must be retained, who decides, or what happens when retention ends. An institution that retires an IKM after the regulatory minimum (typically 7 years for chain-of-custody data) loses the ability to re-verify any chain entry signed under that IKM. Spec §10 should add a normative retention rule for the IKM registry: at minimum, the IKM must be retained as long as any chain entry under that `key_version` is retained, and the institution's IKM-retirement procedure must be documented.

2. **HSM-mediated HKDF (Model B) assumes capability not all HSMs expose.** Design 02 §4.0a ("Model B — Session-key-delivered") says "HSM uses the IKM bound to its key label, HSM computes info_for_tenant and runs HKDF-SHA-256 internally, HSM returns the 32-byte session_key." Most PKCS#11 HSMs do not expose HKDF as a native mechanism. They expose `CKM_SHA256_HMAC` and the application chains the extract+expand calls itself. If the HSM only exposes HMAC, then either (a) the SDK does the HKDF chaining outside the HSM, which contradicts "HSM never lets the IKM out" (the extract step's PRK is a 32-byte key derived from the IKM and is itself sensitive), or (b) the HSM exposes a proprietary HKDF wrapper. Spec §4.1.1 should call out that Model B requires an HSM with native HKDF or a vendor wrapper, and document the fallback (the PRK transits the SDK process briefly, with the same memory-protection posture as Model A's IKM).

3. **Rotation window race during seal computation.** Design 04 §3.2.2 ("Master-key rotation window") says the seal job covering the rotation day records the multi-version state. But the seal job runs at UTC-midnight + 60 minutes (§4.3 / design 03 §4.1). If the rotation completes at 23:55 UTC and one application process holds an old session key for 90 more minutes (process-restart pending), the seal job runs against a partial picture: the late-arriving events under the old IKM appear in the *next* day's seal as late-binding, with `key_version` from the old generation. The verifier sees `key_versions=[old]` on the rotation day's seal and `key_versions=[old, new]` on the day-after seal, with old-IKM late-binding entries on the day after. This is correct behaviour but the spec does not tell the verifier to expect this shape, so a strict-mode verifier might flag the day-after seal's old-IKM late-binding entries as suspicious. A normative note in §4.2.2 / §10 covering "rotation crossing the seal boundary" would close the operational gap.

The cryptographic primitives are right. The operational rotation story has these three loose threads.

---

## Per-role roll-up

| Status | Count |
|---|---|
| Answered | 1 |
| Partial | 3 |
| Gap | 1 |
