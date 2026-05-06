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

### 3.1 Master key (HMAC chain key)

- Held in tenant-controlled storage, typically the same HSM that signs the daily seals (in `keep`/`derive` mode rather than `sign` mode).
- Never reaches the application host.
- Per-tenant: one master per tenant.
- Rotation: the spec does not require periodic rotation; rotation is a tenant policy decision. When rotation occurs, all session keys derived from the previous master remain valid for verification of past events; new sessions derive from the new master.

### 3.2 Daily seal signing key

- Held in HSM, `sign`-mode only. The private key is non-extractable.
- Per-tenant: one signing keypair per tenant.
- Public key: published to the tenant key registry. Both the institution and the regulator hold a copy.
- Rotation: rare. Triggered by HSM lifecycle events (HSM end-of-life, FIPS validation lapse) or by tenant policy. Rotation produces a new keypair; the institution publishes the new public key. Past seals remain verifiable against the previous public key (the registry retains historical keys with their validity windows).

### 3.3 The tenant key registry

A directory of `(tenant_id, key_version, public_key, valid_from, valid_until)` records.

In v1.0, the registry is implementation-flexible. Examples:

- A flat file the institution publishes alongside its seal exports
- A JWKS (RFC 7517) endpoint the institution operates
- A regulator-hosted registry (the FFIEC could in principle host one)

The verifier reads the registry as untrusted input and validates against a path the examiner trusts (typically the institution&rsquo;s public-key fingerprint as recorded with the regulator at registration time).

v1.1 candidate: standardize the registry as a JWKS endpoint with FFIEC-defined fields.

## 4. The signing payload

The text-format choice from the spec, restated for design rationale:

```
ffiec-ai-chain-v1\n
{tenant_id}\n
{ISO 8601 date in UTC}\n
{hex-encoded Merkle root}
```

### 4.1 Why text instead of binary

- **Auditor inspection.** An auditor can copy the payload from a hex dump and re-sign it (with their own key, for testing) to confirm the verifier&rsquo;s signature check.
- **Cross-platform stability.** Text is unambiguous; binary encodings can drift between implementations.
- **Slightly larger payload.** Costs 60ish bytes per day per tenant. Negligible.

### 4.2 Why include the spec version

`ffiec-ai-chain-v1` in the payload prevents cross-version replay. A v1 signature cannot be presented as a v2 signature even if the rest of the payload happens to match.

### 4.3 Why include the tenant_id

Prevents cross-tenant signature replay. A signature from tenant A cannot be presented as covering tenant B&rsquo;s data, even if the Merkle roots happen to coincide.

### 4.4 Why include the date as a string

ISO 8601 in UTC is unambiguous, sortable, and human-readable. The verifier parses it; the auditor reads it.

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

### 5.3 HSM failover

Banks operating mission-critical workloads typically run an HSM cluster (active-active or active-passive). The seal job uses any healthy HSM in the cluster. The signing key is replicated across the cluster according to vendor-specific procedures (hardware-backed key replication for Thales/Entrust; key replication via secure backup for AWS CloudHSM).

### 5.4 HSM compromise

If the HSM is compromised (extraction of the signing key):

- The institution rotates the key immediately.
- All seals signed with the old key are still verifiable, but the integrity claim is now contingent: the verifier cannot prove that any seal signed in the compromise window was produced by the institution rather than by the attacker.
- The institution&rsquo;s incident response includes notification to the regulator, root-cause analysis, and re-attestation of seals from independent evidence (logs of the seal job's invocations, network traces, etc.).

The compromise of an HSM is a known-rare, high-effort attack. The threat model in [`09-threat-model.md`](09-threat-model.md) accepts it as a residual risk.

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
| What if FIPS 186-5 deprecates Ed25519 in the future? | Unlikely on the time scale of this spec, but the spec admits algorithm rotation. v1.1 may add Dilithium (FIPS 204) or SLH-DSA (FIPS 205) alongside Ed25519. The seal record includes the algorithm identifier explicitly, so the verifier can dispatch on it. |
| What if the HSM vendor&rsquo;s FIPS validation lapses? | The institution rotates to a different HSM vendor whose validation is current. Past seals remain verifiable against the past public key. |
| Open issue | Defining the tenant key registry format. v1.0 leaves it implementation-flexible; v1.1 may standardize a JWKS-style endpoint with FFIEC-specified extensions. |
