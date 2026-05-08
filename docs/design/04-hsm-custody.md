# 04 — HSM custody

> **What this doc is.** How the daily Merkle seal is signed, how keys are managed, and how the institution and the regulator share the verification half. The HSM is the integrity anchor for the entire chain &mdash; everything else is plumbing around what the HSM proves.

## 1. Why HSM at all

Every other primitive in this system can be re-computed by anyone with sufficient access to the data and the keys. The HSM provides what software cannot: a key that exists *only* inside a tamper-resistant device, signing only what the institution&rsquo;s authenticated operators authorize.

Without HSM custody:

- An institution&rsquo;s admin with database access can rewrite history and re-sign the new history, indistinguishably from the original.
- A vendor with operational access to the institution&rsquo;s ledger can do the same.
- The integrity claim collapses to "trust us."

With HSM custody:

- The signing key is in a FIPS 140-2 Level 3 (or higher) device.
- Signing requires an authenticated operator action.
- Key extraction requires physical compromise of the device, which is detectable.
- The verifier confirms the signature against a public key the regulator and the institution both hold.

## 2. Algorithm choice

### 2.1 Ed25519, FIPS 186-5

- **Deterministic signatures.** Same payload &rarr; same signature. This makes retry semantics simple (see seal-job step 5 in [`03-merkle-seal.md`](03-merkle-seal.md)).
- **Small (64-byte signatures, 32-byte public keys).** Fits comfortably in a daily seal record.
- **Modern.** Ed25519 has been in stdlib since Go 1.13; FIPS 186-5 added it in 2023.
- **Side-channel resistant.** Ed25519&rsquo;s constant-time scalar multiplication and avoidance of nonce reuse make it more side-channel resistant than ECDSA.
- **Widely supported in HSMs.** AWS CloudHSM, Azure Key Vault Premium, Google Cloud HSM, and major on-prem HSMs (Thales Luna, Entrust nShield, Utimaco) all support Ed25519 by 2026.

### 2.2 Why not ECDSA P-256

ECDSA was the FIPS-approved alternative before FIPS 186-5. Two operational concerns:

- **Nonce randomization.** Each signature requires a per-signature nonce; failures of the RNG leak the private key (cf. PlayStation 3 keys). HSM hardware mitigates this in practice but does not eliminate the historical concern.
- **Variable-length signatures.** ECDSA signatures are DER-encoded and vary in length. Storage and transport are slightly less clean.

Ed25519 sidesteps both. The choice is forward-looking, FIPS-approved, and operationally simpler. Auditors familiar with both will accept either; we standardize on Ed25519 for the cleaner properties.

### 2.3 Why not RSA

RSA-2048 or RSA-3072 are still FIPS-approved and HSM-supported. Two reasons against:

- **Signature size.** 256 or 384 bytes versus Ed25519&rsquo;s 64 bytes.
- **Trajectory.** Post-quantum considerations push away from RSA. Ed25519 is also not post-quantum, but the spec contemplates a future spec version that adds a post-quantum signature alongside Ed25519 (Dilithium or SLH-DSA, both FIPS 204/205-approved as of 2024).

## 3. Key custody

### 3.1 IKM (Input Key Material — the HMAC chain key)

(Terminology note: "IKM" replaces the older "master key" term throughout the spec rework. The two names refer to the same thing — the long-lived secret a tenant's session keys are derived from. "Master key" remains acceptable in operational shorthand; the spec text uses "IKM" for precision.)

- Held in tenant-controlled storage, typically the same HSM that signs the daily seals (in `derive`/`mac` mode rather than `sign` mode), or in cloud KMS with envelope encryption.
- MAY reach the application host (Model A in `02-chain-construction.md` §4.0) OR may stay in-HSM with HKDF performed inside the device (Model B in §4.0a). Spec §4.1.1 specifies both delivery models.
- Per-tenant: one IKM per tenant. The IKM bytes plus the tenant_id together produce a per-tenant `key_fingerprint = SHA-256(utf8(tenant_id) || ikm)[:16]` that is stamped on every chain entry; the verifier asserts the looked-up IKM produces the recorded fingerprint before any MAC compute.
- Length: minimum 32 bytes per spec §10.6 (RFC 4868).
- Rotation: the spec does not require periodic rotation; rotation is a tenant policy decision. When rotation occurs, the new IKM gets a new `key_version` integer; chain entries captured after rotation stamp the new `key_version` and the new `key_fingerprint`. Past entries (with the old `key_version` and old `key_fingerprint`) remain verifiable as long as the old IKM is retained in the tenant key registry. Spec §4.2 seal record's `key_versions` field records the set of versions present on each tenant-day.

### 3.2 Daily seal signing key(s)

- Held in HSM, `sign`-mode only. The private key is non-extractable.
- **Per-tenant: one signing keypair per algorithm per tenant.** For single-algorithm posture (default v1.0): one Ed25519 keypair per tenant. For dual-algorithm transitional posture (v1.x post-quantum): one Ed25519 keypair AND one post-quantum keypair (Dilithium / SLH-DSA), both per tenant. The institution's HSM holds both keypairs; the seal job signs the day's Merkle root under both keys per spec §4.2 `signatures` list (Variant B per spec §4.2: each algorithm's signature covers its own algorithm-bound `sign_payload`).
- Public keys: published to the tenant key registry. Both the institution and the regulator hold copies. Under dual-algorithm posture, both algorithms' public keys are published; the registry entry carries the algorithm identifier per key.
- Rotation: rare. Triggered by HSM lifecycle events (HSM end-of-life, FIPS validation lapse) or by tenant policy. Rotation produces a new keypair; the institution publishes the new public key. Past seals remain verifiable against the previous public key (the registry retains historical keys with their validity windows).
- Algorithm-posture transitions: introducing the second algorithm (single → dual) requires the institution to provision the new-algorithm keypair, publish the public key, and update its declared algorithm posture (the configuration the verifier consults per spec §7 step 11). Retiring the first algorithm (dual → single under the new algorithm) is the second half of the transition, gated on the broken-algorithm timeline and regulator coordination.

### 3.2.1 Cloud HSM provisioning by provider

The FIPS 140-2 L3 requirement maps differently to each major cloud. Spec §10.5 names the precise conformant tiers; the table below is the design-doc summary:

| Cloud HSM | Default FIPS level | Notes |
|---|---|---|
| AWS CloudHSM Classic | FIPS 140-2 L3 | L3 by default; per-cluster pricing |
| AWS CloudHSM v2 | FIPS 140-2 L3 | L3 by default; the successor product line |
| Azure Managed HSM | FIPS 140-2 L3 | The conformant Azure tier; backs the Key Vault Premium "HSM-protected key" feature |
| Azure Dedicated HSM | FIPS 140-2 L3 | Thales Luna under the hood |
| Google Cloud HSM | FIPS 140-2 L3 | L3 by default for HSM-protected keys |
| AWS KMS (without CloudHSM backing) | FIPS 140-2 L2 | **Not conformant** for daily seal signing; use CloudHSM Classic / v2 instead |
| Azure Key Vault Standard | FIPS 140-2 L1 | **Not conformant** for daily seal signing |
| Azure Key Vault Premium (software-protected keys) | FIPS 140-2 L2 | **Not conformant** for daily seal signing; the Premium "HSM-protected key" feature is conformant because it uses Azure Managed HSM under the covers |

Bank IT teams routinely confuse Azure Key Vault Standard with the Managed HSM and Dedicated HSM tiers, and confuse Premium "software-protected" with Premium "HSM-protected." The spec's L3-or-higher requirement disqualifies Standard and Premium-software-protected. Implementations that depend on cloud HSM SHOULD document the specific tier and the FIPS validation certificate number in their control description.

**AWS KMS Custom Key Stores backed by CloudHSM** are conformant: the cryptographic operations are performed by the underlying CloudHSM, and the FIPS validation flows through. The KMS layer adds key-policy management without weakening the cryptographic property.

**Cross-region key replication.** AWS CloudHSM uses backup-and-restore for cross-region replication; per-tenant key labels are preserved through the replication. Azure Managed HSM supports geo-redundant configuration with key labels maintained. Google Cloud HSM supports cross-region replication with the same property. In all cases, the per-tenant signing key remains isolated across regions; cross-tenant signing remains structurally prevented. Institutions document the cross-region replication procedure in their DR plan.

### 3.2.2 IKM rotation window

When the institution rotates the IKM, some application processes hold session keys derived from the old IKM while others hold session keys derived from the new IKM until they restart. The rotation window is bounded by the longest-lived application process (typically minutes to hours).

During the window, both `key_version` values are valid. Events captured under the old `key_version` are validated against the old IKM; events under the new `key_version` are validated against the new IKM. The verifier resolves which IKM to apply per event via the entry's stamped `key_version` (spec §7 step 7) and asserts the looked-up IKM produces the entry's recorded `key_fingerprint` (spec §7 step 8) before any MAC compute.

The seal job covering the rotation day records `key_versions` as a list (e.g., `[3, 4]`) when both versions appear within the day's events (spec §4.2 seal record schema). The verifier handles a multi-`key_version` seal day correctly by applying the appropriate IKM per event via the per-entry lookup.

After the rotation window closes (all application processes have restarted under the new IKM), seals revert to a single-element `key_versions` list. Spec §10.10 ("Rotation crossing the seal boundary") names the operational reality that late-arriving events under the old IKM may appear in the day-after seal as late-binding entries; this is normal-operations PASS-with-anomaly behavior, not a failure.

### 3.3 The tenant key registry

A directory of `(tenant_id, key_version, public_key, valid_from, valid_until)` records.

In v1.0, the registry is implementation-flexible. Examples:

- A flat file the institution publishes alongside its seal exports
- A JWKS (RFC 7517) endpoint the institution operates
- A regulator-hosted registry (the FFIEC could in principle host one)

The verifier reads the registry as untrusted input and validates against a path the examiner trusts (typically the institution&rsquo;s public-key fingerprint as recorded with the regulator at registration time).

v1.x roadmap candidate: standardize the registry as a JWKS endpoint with FFIEC-defined fields, once enough institutions have shipped to inform the field set.

## 4. The signing payload

The text-format choice from spec §4.3, restated for design rationale. Under v1.0-final-amendment the canonical form is the v1.0a 10-line `sign_payload`:

```
ffiec.chain-of-custody.v1\n
{sign_payload_version}\n            // "v1.0a" under v1.0-final-amendment
{algorithm}\n                       // "ed25519" for v1.0
{format_version}\n                  // "v1" for this spec
{tenant_id}\n
{ISO 8601 date (YYYY-MM-DD UTC)}\n
{hex-encoded Merkle root, 64 chars lowercase}\n
{hex-encoded HKDF inputs digest, 64 chars lowercase}\n
{cadence}\n                         // "hourly" | "daily" | "weekly"
{dev_mode}                          // single byte: "1" or "0"; no trailing \n
```

Each inter-field separator is a single `\n` (0x0A) byte; the terminal `dev_mode` byte has NO trailing newline. Total LF separators: nine. CRLF is non-conformant. Hex fields use lowercase, no separators, zero-padded to 64 characters. The `{algorithm}` line carries the signature algorithm identifier (`ed25519` for v1.0; future post-quantum algorithms dispatch here per spec §4.3.2).

Pre-amendment chains (produced before 2026-05-07) omit the `sign_payload_version` line and use the pre-amendment 6-line form (magic line + `algorithm` + `format_version` + `tenant_id` + `iso8601_date` + `hex(merkle_root)` + `hex(hkdf_inputs_digest)`). The verifier dispatches on the seal record's `sign_payload_version` field per spec §4.3 and reconstructs the matching layout, so pre-amendment chains remain verifiable under amendment-aware verifiers without re-sealing.

Under dual-algorithm posture (spec §4.3.2, Variant B), each algorithm in the seal's `signatures` list is computed over its OWN algorithm-bound `sign_payload`: line 3 (`algorithm`) carries that algorithm's identifier; line 2 (`sign_payload_version`) remains `"v1.0a"` for every algorithm. A single shared `sign_payload` covering all algorithms is non-conformant because it leaks the algorithm-confusion defense.

### 4.1 Why text instead of binary

- **Auditor inspection.** An auditor can copy the payload from a hex dump and re-sign it (with their own key, for testing) to confirm the verifier's signature check.
- **Cross-platform stability.** Text is unambiguous; binary encodings can drift between implementations.
- **Slightly larger payload.** Costs ~140 bytes per day per tenant after the rework (was ~60 before). Negligible.

### 4.2 Why include the spec-prefix line

`ffiec.chain-of-custody.v1` in the payload prevents cross-version replay. A v1 signature cannot be presented as a v2 signature even if the rest of the payload happens to match. The naming aligns with the HKDF constants in spec §4.1 (`HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"`).

### 4.3 Why include the algorithm

The signature algorithm identifier closes the algorithm-confusion attack class (cf. JWT `alg=none`, SAML algorithm-substitution). v1.0 already admits the dual-algorithm transitional posture per spec §4.3.2 (the seal record's `signatures` list, Variant B, with each entry computed over its own algorithm-bound `sign_payload`); the binding is in force today and not deferred to a future spec version. Single-algorithm v1.0 chains ship Ed25519. Once an institution adds a v1.x post-quantum algorithm alongside Ed25519 under Variant B, the attacker class "present an Ed25519 signature as a Dilithium signature for a key that happens to match" is closed cryptographically. The cost is one ~10-byte line in the signed payload.

### 4.4 Why include the format_version

The chain-stamp `format_version` (currently `"v1"`) is independent of the spec version. A future spec v1.1 might keep `format_version=v1`; a future spec v2 might introduce `format_version=v2` with a different chain construction. Including `format_version` in the signed payload lets the verifier dispatch on it without trusting the seal record's metadata.

### 4.5 Why include the tenant_id

Prevents cross-tenant signature replay. A signature from tenant A cannot be presented as covering tenant B's data, even if the Merkle roots happen to coincide.

### 4.6 Why include the date as a string

ISO 8601 in UTC is unambiguous, sortable, and human-readable. The verifier parses it; the auditor reads it.

### 4.7 Why include the HKDF inputs digest

The `hkdf_inputs_digest` binds the seal to the specific per-tenant HKDF inputs (salt, info_for_tenant, length) in force on the day, where `info_for_tenant = HKDF_INFO_BASE || "|" || utf8(tenant_id)`. The per-tenant variant matches the file-header digest defined in spec §3 and used at the verifier's pre-flight (spec §7 step 2), so the seal-time and file-pre-flight call sites compute the same bytes. A future-version verifier reading a v1 seal walks this field first; mismatch refuses with the precise reason ("HKDF inputs do not match running v1 inputs"). Without this field, a misconfigured v2 verifier could accept v1 seals as v2 seals if the tenant + date + root happened to match.

### 4.8 Why include cadence and dev_mode

Both fields are bound under the HSM signature so a coordinated forgery cannot rewrite the institution's claimed posture without forging the signature itself. `cadence` (line 9) catches a seal-record rewrite that flips daily to weekly to claim a relaxed posture; before the v1.0a form, the only check was cross-document reconciliation against the institution's regulator-approved CC8.1 control description, which is operationally heavy. `dev_mode` (line 10) catches the more serious case: an attacker with seal-record write access could otherwise present a chain produced by the §10.7 development software-key adapter as a production chain, defeating the regulator-visible-line guarantee that §10.7 establishes. Binding `dev_mode` cryptographically means the flip now requires forging the HSM signature, which the FIPS 140-2 Level 3 custody posture rules out.

### 4.9 Why include sign_payload_version

`sign_payload_version` (line 2) is the form-generation discriminator. It binds itself into the signed bytes so a tampered value is detected at signature verification — the verifier reconstructs `sign_payload` using the field as written, and a mismatch between the field and the byte form actually used by the signer produces a signature failure. The discriminator also lets future v1.x amendments add lines to `sign_payload` (with values like `"v1.0b"`, `"v1.0c"`) without breaking pre-amendment verifiability: the verifier dispatches on the field and reconstructs the matching layout. A pre-amendment chain that omits the field is still verifiable under the pre-amendment 6-line form; an amendment chain that sets the field to `"v1.0a"` is verifiable under the 10-line form. The cost is one ASCII line; the benefit is forward-compatible byte-form evolution.

## 5. HSM operations

### 5.1 Signing flow

1. The seal job opens an HSM session (PKCS#11 token login or cloud-HSM API call).
2. The job authenticates with the HSM operator credentials (HSM-specific; not in scope here).
3. The job submits the `sign_payload` for signing under the tenant&rsquo;s key label.
4. The HSM returns a 64-byte Ed25519 signature.
5. The job closes the session.

The session-level authentication is a per-deployment policy. Production deployments use separation-of-duties: the seal job has authentication credentials that allow `sign` operations only, not `extract`, `delete`, or `import`.

### 5.2 HSM unavailability

If the HSM is unavailable:

- Captured events continue to be ingested and chained (the per-event HMAC is independent).
- The daily seal job fails for the affected day. The job retries with exponential backoff up to a configurable cap.
- After the cap, the job emits an alert. The seal is delayed.
- When the HSM is restored, the job runs and produces the seal. The signing time is recorded in the seal record.
- The verifier reports the delay explicitly. A multi-day delay is unusual and flagged in the examiner report.

**Notification thresholds.** Institutions SHOULD notify their primary regulator if the daily seal is delayed beyond 72 hours. The 72-hour window covers the common operational outages (HSM cluster failover, scheduled maintenance, cloud-provider regional incidents) without being so tight that every transient failure escalates. Persistent delays beyond 72 hours indicate a real operational problem and trigger the institution's cyber-incident notification framework.

### 5.2.1 HSM PIN and credential rotation

The PKCS#11 PIN authenticates the seal job to the HSM. Three operational practices apply:

- **Never embed the PIN in the ledger config file.** The reference config (`06-ledger-server-design.md` §6) reads the PIN from an environment variable; the environment variable is populated by the institution's secret-management system at process start.
- **Rotate the PIN on a documented cadence.** Quarterly is the typical bank operational cadence; institutions document the cadence in their control description. After rotation, the previous PIN MUST be invalidated at the HSM.
- **Use separation of duties.** The PIN holder (typically the seal-job operator) is separated from the HSM administrator (who can extract or import keys) and from the database administrator (who can alter ledger contents). Two-of-three collusion remains possible; the institution's audit program is the mitigation.

Rotation procedure:

1. The HSM administrator generates a new PIN at the HSM.
2. The seal-job operator updates the secret-management system with the new PIN.
3. The seal job's running processes re-read the secret on next reload (config-reload signal or container restart).
4. The HSM administrator invalidates the previous PIN.
5. The institution records the rotation in its control-evidence log.

### 5.3 HSM failover

Banks operating mission-critical workloads typically run an HSM cluster (active-active or active-passive). The seal job uses any healthy HSM in the cluster. The signing key is replicated across the cluster according to vendor-specific procedures (hardware-backed key replication for Thales/Entrust; key replication via secure backup for AWS CloudHSM).

### 5.4 HSM compromise

If the HSM is compromised (extraction of the signing key):

- The institution rotates the key immediately.
- All seals signed with the old key are still verifiable, but the integrity claim is now contingent: the verifier cannot prove that any seal signed in the compromise window was produced by the institution rather than by the attacker.
- The institution&rsquo;s incident response includes notification to the regulator, root-cause analysis, and re-attestation of seals from independent evidence (logs of the seal job's invocations, network traces, etc.).

The compromise of an HSM is a known-rare, high-effort attack. The threat model in [`09-threat-model.md`](09-threat-model.md) accepts it as a residual risk.

### 5.5 Verifier binary trust path

The verifier is signed with cosign at release time. Cosign signatures chain to a Sigstore root. Institutions and regulators SHOULD pin a known-good cosign public key from a trusted distribution channel (the project's GitHub release page hash, signed by a published GPG key the institution validates out-of-band). Two-stage verification:

1. **Cosign verification.** Confirms the binary was signed by the project's release pipeline.
2. **Reproducible-build verification.** Independent rebuild from source produces a byte-for-byte identical binary. Institutions that depend on the verifier in a regulatory submission SHOULD perform an independent rebuild and compare hashes at least once per release.

If Sigstore is compromised, the institution falls back to the GPG-signed manifest the project publishes alongside each release. The manifest lists the binary's SHA-256 hash; the institution validates against the GPG signature using a key it has cached out-of-band.

The cosign-trust path is governance, not cryptography. Documented in `GOVERNANCE.md`.

## 6. Software-key fallback (development only)

The reference implementation supports a software-backed Ed25519 keypair for development and CI:

- The fallback is enabled only when `--allow-software-key` is passed and `FFIEC_ALLOW_SOFTWARE_KEY=1` is set in the environment.
- Seals signed with a software key carry a clearly-marked `dev-mode: true` field in the seal record.
- The verifier refuses to validate a `dev-mode` seal as a production seal. A dev-mode seal can only be verified in dev mode.
- The reference implementation logs a warning every time the software fallback is used.

The double-flag-plus-marking design is to make the development path obvious and impossible to accidentally roll into production.

## 7. Auditor's-lens review

| Question | Answer |
|---|---|
| Who can sign with the tenant key? | Only the operators authenticated to the HSM with the seal-job role. The role grants `sign` only; not `extract`, `delete`, or `import`. |
| What if the seal-job operator credentials are compromised? | The compromise allows the attacker to sign seals over Merkle roots they control, as long as they can also write to the ledger. This is a privileged-insider scenario; the threat model treats it as a high-effort attack and recommends separation of duties (the seal-job operator should be separated from the database admin role). |
| Is the public key trustworthy? | The institution publishes the public key with its regulator at the time of tenant registration. The verifier validates the registry-published key against the regulator-held key fingerprint. The chain of trust is the regulator&rsquo;s. |
| What if FIPS 186-5 deprecates Ed25519 in the future? | Unlikely on the time scale of this spec, but the spec admits algorithm rotation today. The seal record carries the algorithm identifier explicitly and §4.3.2 admits Variant B dual-algorithm seals (per-algorithm `sign_payload`), so an institution can co-sign with Dilithium (FIPS 204) or SLH-DSA (FIPS 205) alongside Ed25519 once HSM products ship FIPS-validated post-quantum support. The verifier dispatches per-algorithm via the seal record's `signatures` list. |
| What if the HSM vendor&rsquo;s FIPS validation lapses? | The institution rotates to a different HSM vendor whose validation is current. Past seals remain verifiable against the past public key. |
| Tenant key registry format | v1.0 leaves it implementation-flexible — every institution already has an internal directory shape its IT and compliance teams know. A standardized JWKS-style endpoint with FFIEC-specified extensions is a v1.x roadmap commitment that lifts the lowest-friction shape into spec text once enough institutions have shipped to inform the field set. |
