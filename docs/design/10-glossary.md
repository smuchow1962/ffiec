# 10 — Glossary

> **What this doc is.** Terms that appear in the design docs and are not common everyday language. Each entry is short and defined in the context of this chain-of-custody system. Entries cross-reference the design doc that defines the concept in full.

## A

### Advisory lock
A coordination primitive in PostgreSQL (and similar databases) that lets a process declare "I'm working on X; nothing else should." Used by the seal job to prevent two workers from sealing the same `(tenant_id, day)` pair concurrently. See [`06-ledger-server-design.md`](06-ledger-server-design.md).

### Algorithm rotation
Replacing a cryptographic algorithm in the spec because the old one weakened. The seal record carries an explicit algorithm identifier so the verifier can dispatch on it and the spec can introduce, for example, a post-quantum signature alongside Ed25519. Different from key rotation, which keeps the algorithm and changes the key.

### Append-only
Storage where rows are inserted but never updated or deleted. Two layers: at the application level (no UPDATE or DELETE statements appear in the codebase) and at the storage level when supported. The integrity property the verifier depends on.

### Auditor's-lens
The convention that every design doc closes with a section answering questions a federal examiner or external auditor would ask. The lens names the worst questions, the design's answer, and any open issues acknowledged but not yet resolved.

### Authentic capture
The integrity property that an event was produced by the institution's legitimate AI agent process at the time and run-id it claims. One of three properties the chain provides. See [`09-threat-model.md`](09-threat-model.md).

## B

### Backpressure
What happens when a downstream component can't keep up. The slow stage fills its buffer, which slows the stage feeding it, which surfaces to the SDK as a retryable OTLP error. The chain never silently drops events; it pushes back instead.

### Big Four
The four major external audit firms (Deloitte, EY, KPMG, PwC) that financial institutions hire to attest to controls. One of the audiences the chain is designed for.

### BYOC (Bring-your-own-cloud)
A deployment topology where the bank operates the ledger server in its own cloud account. The vendor provides the image and operational support but does not hold the data or the HSM keys. See [`00-overview.md`](00-overview.md).

## C

### Canonical form
A deterministic byte representation of a JSON payload. Two implementations that agree on the canonical form produce identical bytes for the same logical event. The chain hashes the canonical form, not the wire bytes.

### Capture
The act of recording an AI event into the chain. Performed by the SDK at the moment the agent makes a decision, calls a tool, or produces a response.

### Chain of custody
The full record of who held what evidence, when, and what they did with it. In this system, an end-to-end audit trail of AI decisions whose integrity can be verified independently by a regulator.

### Conformance corpus
The directory of test vectors that distinguishes conforming from non-conforming implementations. Run the corpus; it passes or fails. See [`08-test-vectors.md`](08-test-vectors.md).

### Cosign
An open-source tool for signing software artifacts with short-lived keys backed by Sigstore. Used to sign verifier release binaries; examiners verify the signature before running the tool.

### Constant-time
Code that takes the same number of CPU cycles regardless of input values. Required for cryptographic comparisons (`hmac.Equal`) so that an attacker cannot infer secret data from timing differences.

### CUPID
A set of design properties (Composable, Unix philosophy, Predictable, Idiomatic, Domain-based) that the codebase prefers over rigid enforcement of SOLID. Referenced in `CODING_INSTRUCTIONS.md`.

### Cryptographic break
A discovery that an algorithm previously believed secure can be subverted in practice. The spec admits algorithm rotation as the response; past records remain verifiable against their original algorithm with the caveat noted in the residual-risk register.

## D

### Daily Merkle seal
The 32-byte root of an RFC 6962 Merkle tree built over every event for one `(tenant_id, UTC_date)` pair, signed by the HSM. The artifact that catches retroactive tampering. See [`03-merkle-seal.md`](03-merkle-seal.md).

### Daily root
Shorthand for the `merkle_root` that anchors one tenant-day's events. The verifier recomputes this from the ledger and compares it to the recorded value.

### DER encoding
A binary encoding format used for some cryptographic objects (e.g., ECDSA signatures). Variable-length, which is one reason the spec prefers Ed25519's fixed 64-byte signatures.

### Determinism
The property that the same inputs produce byte-for-byte identical outputs. Required for the chain hash, the Merkle root, the verifier report, and (with `--deterministic`) the PDF output. The conformance corpus exists to enforce it.

### Dilithium
A post-quantum digital signature algorithm standardized in FIPS 204. v1.0 already admits the dual-algorithm transitional posture (spec §4.3.2 Variant B); a v1.x amendment will normate Dilithium support once HSM products ship FIPS-validated implementations at scale.

## E

### Ed25519
A modern elliptic-curve digital signature algorithm. FIPS 186-5 approved, deterministic (same payload yields same signature), 64-byte signatures, 32-byte public keys. The signing algorithm for daily seals. See [`04-hsm-custody.md`](04-hsm-custody.md).

### Empty day
A tenant-day with zero captured events. Still receives a seal, with the Merkle root defined as `SHA-256("")`. Continuity matters: a missing seal is a gap an examiner notices.

### ECDSA
Elliptic Curve Digital Signature Algorithm. The previous-generation signature standard. The spec rejects it for daily seals because of nonce-randomization brittleness and variable-length signatures.

## F

### FFIEC
Federal Financial Institutions Examination Council. The U.S. interagency body whose member agencies (OCC, Federal Reserve, FDIC, NCUA, CFPB) examine banks. The primary regulator audience for this spec.

### FinOps
Cloud financial operations. Tools like CloudZero and Vantage that ingest OTel traces with `gen_ai.usage.*` attributes for cost roll-up. The chain's OTLP wire passes through to FinOps tools unchanged.

### FIPS 140-2 / FIPS 140-3
U.S. federal cryptographic-module standards. Level 3 adds physical tamper resistance and identity-based authentication; the spec requires Level 3 or higher for HSMs holding signing keys. FIPS 140-3 supersedes 140-2 for new validations and is accepted equivalently.

### FIPS 186-5
The 2023 update to the U.S. digital signature standard, which added Ed25519 to the approved list.

### FIPS 198-1
The U.S. standard that approves HMAC. Cited so an auditor can recognize HMAC-SHA-256 as approved without further explanation.

### FIPS 204 / FIPS 205
Post-quantum signature standards (Dilithium and SLH-DSA, respectively), approved in 2024. v1.0 admits the dual-algorithm transitional posture today via the seal record's `signatures` list (Variant B per spec §4.3.2); a v1.x amendment will normate Dilithium / SLH-DSA support once HSM products ship FIPS-validated implementations at scale.

### Forward-only forgery
The bounded attack window that remains after an application process is compromised. The attacker can produce events going forward but cannot retroactively alter past events, which have already been chained, exported, and (eventually) sealed.

### Fsync
The OS-level operation that forces buffered data to physical disk. The SDK's local SQLite uses `synchronous=FULL` for audit events so a crash cannot lose a disclosed event.

### Fuzzing
Generating randomized inputs to surface defects that hand-written tests would miss. Two independent implementations are required to produce identical bytes on a randomly-generated event stream.

## G

### GenAI semantic conventions (gen_ai.*)
The OpenTelemetry attribute namespace for AI-related telemetry: `gen_ai.system`, `gen_ai.request.model`, `gen_ai.usage.prompt_tokens`, etc. The chain uses these unchanged; the chain extension lives in a separate `ffiec.chain.*` namespace.

### gRPC
A high-performance RPC framework using HTTP/2 and protobuf. The default OTLP transport.

## H

### HKDF
HMAC-based Key Derivation Function (RFC 5869). Derives session keys from a tenant master without ever exposing the master to the application host. Inputs: master (`IKM`), salt, info string; output: a fixed-length key.

### HMAC
Keyed Hash Message Authentication Code. A hash construction that produces a tag binding both data and a secret key. An attacker without the key cannot produce a matching tag. The chain uses HMAC-SHA-256.

### HMAC chain
The per-event integrity primitive: each `payload_hash` is `HMAC(session_key, prev_hash || canonical_payload)`. Catches in-flight wire tampering. See [`02-chain-construction.md`](02-chain-construction.md).

### Hot path
Code that runs once per captured event in the application process. Performance, determinism, and crash-safety all matter here. The hot path's wall-clock budget is under 1 ms per event.

### Hot store
The online query tier of the ledger server, typically Postgres + JSONB. Derived from the WAL; rebuildable byte-for-byte if it is corrupted.

### HSM (Hardware Security Module)
A tamper-resistant device that holds private keys and performs signing operations without the key ever leaving the device. The integrity anchor of the chain: even an attacker with full database admin and full process compromise cannot forge a daily root because they cannot access the HSM.

## I

### Idempotent
A property of an operation: running it twice yields the same result as running it once. The seal job is idempotent so a partial failure can be retried safely.

### IKM (Input Keying Material)
HKDF's term for the secret input from which session keys are derived. The tenant master is the IKM in the chain's session-key handshake.

### Independent verifiability
The integrity property that a regulator with no access to the institution beyond a public key can verify the chain. The verifier is a single static binary with no network calls; it produces a defensible report from untrusted input.

### Iceberg
Apache Iceberg, a table format on top of Parquet that supports snapshot management and time-travel queries. Used for the cold-store tier. See [`06-ledger-server-design.md`](06-ledger-server-design.md).

### ISO 8601
The international standard for date/time formatting (e.g., `2026-04-01`). Used in the seal payload for unambiguous, sortable, human-readable date representation.

## J

### JCS (JSON Canonicalization Scheme, RFC 8785)
A deterministic JSON encoding: object keys sorted by Unicode codepoint, whitespace normalized, numbers in a defined form, strings escaped consistently. The chain hashes the JCS form, not the wire form. Two implementations that agree on JCS produce byte-for-byte identical chain hashes.

### JSONB
PostgreSQL's binary JSON type. Indexed via GIN for path queries. The hot store stores canonical event payloads as JSONB.

### JWKS (JSON Web Key Set, RFC 7517)
A standard format for publishing public keys. A v1.x roadmap candidate for the tenant key registry; v1.0 leaves the registry implementation-flexible and institutions document their retrieval path in CC8.1.

## K

### Key rotation
Replacing a key with a new one of the same algorithm. Keeps past records verifiable against the previous key and starts new records under the new key. Different from algorithm rotation.

## L

### Late-binding event
An event that arrives at the ledger after its day's seal has already been computed. Marked with a `late_binding` flag and included in the next day's Merkle leaves; the original day's seal is not altered. A high rate of late-binding signals operational issues, not tampering.

### Leaf hash / leaf-prefix
RFC 6962 hashes leaves with a `0x00` prefix and internal nodes with a `0x01` prefix. The domain separation prevents an attacker from presenting an internal-node hash as a leaf or vice versa.

### Ledger
The append-only event store. In this system specifically: the durable record of every captured AI event, indexed by `(tenant_id, run_id, seq)`, that the verifier reads.

### Ledger server
The reference ingest server (`ledger/`). Receives OTLP, re-verifies the HMAC chain, writes to the WAL, computes the daily Merkle seal, and signs in HSM custody. See [`06-ledger-server-design.md`](06-ledger-server-design.md).

## M

### Master key
The per-tenant HMAC key from which session keys are derived via HKDF. Held in tenant-controlled storage (typically the same HSM that signs daily seals, in derive-mode rather than sign-mode). Never reaches the application host.

### Merkle root
The 32-byte hash at the top of an RFC 6962 binary Merkle tree built over a tenant-day's events. The artifact the HSM signs.

### Merkle tree
A tree of hashes where each internal node is the hash of its two children. Modifying any leaf changes every hash on the path from leaf to root, making tampering detectable. The spec uses the RFC 6962 binary construction.

### mTLS (mutual TLS)
A TLS variant where both client and server present certificates. Used as one of the authentication options between the SDK and the ledger server.

## N

### Negative tests
Test cases where the verifier MUST report failure. A verifier that produces "pass" on any negative test is broken. Guards against false-positive verifier behavior.

### Nonce
A number used once in a cryptographic protocol. ECDSA requires a per-signature nonce; reuse leaks the private key. Ed25519 sidesteps the issue by deriving the nonce deterministically from the message and key.

## O

### OpenTelemetry (OTel)
The CNCF-graduated, multi-language, vendor-neutral observability framework. Banks already deploy OTel collectors; the chain rides on the OTel wire to avoid forcing a separate ingestion path.

### OTLP (OpenTelemetry Protocol)
The wire format for OTel data: protobuf over gRPC or HTTP. The chain transports its extension fields as ordinary OTLP span attributes.

### OTel Collector
A daemon that receives OTel data, processes it (filter, sample, redact), and exports it to backends. A common deployment puts a collector at the bank's edge, forwarding to both the ledger and a SIEM.

## P

### Parquet
A columnar storage format optimized for analytics. The cold-store tier uses Parquet for long-retention archives.

### payload_hash
The 32-byte HMAC-SHA-256 output for one event: `HMAC(session_key, prev_hash || canonical_payload)`. The chain links each event to the previous one in its run.

### PEM
A text encoding for cryptographic objects (keys, certificates) using base64 with `-----BEGIN ...-----` framing. Public keys in the tenant key registry are PEM-encoded Ed25519.

### PKCS#11
A standardized C API for talking to HSMs and other cryptographic tokens. The ledger configures a PKCS#11 module path to access the HSM.

### Postgres / PostgreSQL
The relational database used for the WAL and hot store in the reference ledger server. JSONB plus GIN indexes handle the canonical event payload; range partitioning handles per-day retention.

### prev_hash
The 32-byte payload_hash of the previous event in the same run. For `seq == 1`, all-zero. The link that ties events into a tamper-evident sequence.

### Process UUID
A 16-byte random identifier the SDK generates at startup and uses as part of the HKDF salt when deriving its session key. Ensures two processes within the same tenant get independent session keys.

### Protobuf
Google's binary serialization format. The OTLP wire is protobuf. Non-deterministic for several constructs, which is why the chain hashes JCS-canonical JSON instead of protobuf bytes.

## R

### Reproducible builds
Builds where independent rebuilds produce byte-for-byte identical artifacts. Combined with cosign signatures, they let third parties verify that a published binary matches the published source.

### Repudiable
Of an event: not provably the institution's. Events captured during a master-key compromise window are repudiable because the institution cannot prove they were produced by legitimate processes versus by an attacker holding the leaked key.

### Residual risk
A risk the threat model accepts as acknowledged but not fully mitigated. Each residual risk has an owner and a mitigation status; the register documents what the chain catches and what it does not.

### RFC 6962
The Certificate Transparency Merkle tree specification. Adopted because it is the most widely-trusted Merkle construction in regulator-adjacent systems and because second-preimage attacks are explicitly mitigated.

### RFC 8785
The JSON Canonicalization Scheme. See JCS.

### Ring buffer
A bounded circular buffer. The SDK's local SQLite acts as a ring buffer with a configurable cap; once full, the oldest entries are evicted (and the eviction is recorded).

### RSA
A long-established public-key algorithm. Still FIPS-approved but rejected for daily seals because of larger signature size and worse post-quantum trajectory.

### Run
A logical AI agent execution. Identified by `run_id`. Events within a run share the same chain (each event's `prev_hash` is the previous event's `payload_hash`); inter-run parallelism is unbounded.

## S

### Seal job
The periodic task in the ledger server that builds the daily Merkle tree, submits the root to the HSM for signing, and records the resulting seal. Runs at UTC midnight + a configurable delay (default 60 minutes) per tenant.

### seq
A monotonically increasing integer within a run, starting at 1. The SDK rejects out-of-order or skipping `seq` values.

### Session key
A per-process HMAC key derived from the tenant master via HKDF. Held in process memory only; destroyed on process exit. A new process gets a new session key via the same handshake.

### session_key_id
An opaque 16-byte identifier (the leading 16 bytes of `SHA-256(session_key)`) exported on every event. Lets the ledger and verifier identify which session key signed an event without exposing the key itself.

### SHA-256
The 256-bit member of the SHA-2 family. The hash used inside HMAC and inside the Merkle tree construction. FIPS-approved.

### Side-channel attack
An attack that infers secret data from indirect observations (timing, power consumption, electromagnetic emissions). The SDK uses constant-time HMAC comparison; the HSM's FIPS 140-2 L3 protections cover the signing key.

### SIEM (Security Information and Event Management)
A system that ingests events from many sources for detection and correlation. Splunk, Sentinel, QRadar, CrowdStrike Falcon. The OTLP stream tees to the SIEM alongside the ledger.

### Signed root
A daily Merkle root together with its Ed25519 signature. The artifact the verifier validates against the tenant's public key.

### Signing payload
The text-format input the HSM signs. Under v1.0-final-amendment the canonical form is the v1.0a 10-line `sign_payload` (spec §4.3):
```
ffiec.chain-of-custody.v1\n
v1.0a\n                            (sign_payload_version)
{algorithm}\n                      (e.g. "ed25519")
{format_version}\n                 ("v1")
{tenant_id}\n
{YYYY-MM-DD UTC}\n
{hex(merkle_root) — 64 chars lowercase}\n
{hex(hkdf_inputs_digest) — 64 chars lowercase}\n
{cadence}\n                        ("hourly" | "daily" | "weekly")
{dev_mode}                         (single byte: "1" or "0"; no trailing \n)
```
Nine `0x0A` separators total; no trailing newline; CRLF non-conformant. Pre-amendment chains use the 6-line form (no `sign_payload_version`, no `cadence`, no `dev_mode`); the verifier dispatches on the seal record's `sign_payload_version` field. Text rather than binary so an auditor can read it.

### SLH-DSA
A stateless hash-based digital signature algorithm standardized in FIPS 205. Like Dilithium, supported under v1.0's dual-algorithm transitional posture (spec §4.3.2 Variant B) and a candidate for normative v1.x amendment text once HSM products mature.

### SPIFFE / SPIRE
A standard and reference implementation for issuing identity credentials to workloads. A v1.x roadmap candidate for the tenant key handshake mechanism; v1.0 leaves the handshake protocol institution-defined and the chain inherits whatever delivery posture the institution operates.

### SR 11-7
The Federal Reserve's supervisory letter on model risk management. Sets expectations for model validation. Out of scope for the chain; the chain captures evidence of decisions, not whether the decisions are correct.

### Streaming Merkle
An iterative construction of the Merkle root that holds only `O(log N)` hashes in memory regardless of leaf count. Necessary for tenants with billions of events per day.

### Substrate
A captured stream that serves multiple downstream consumers without translation. The OTLP wire is the chain's substrate: one stream feeds the ledger, the SIEM, FinOps tools, and any other OTel-aware backend.

## T

### Tamper evidence
The integrity property that any modification to a captured event is detectable. Achieved by the HMAC chain (catches in-flight tampering) plus the Merkle seal under HSM signature (catches retroactive tampering).

### Tenant
A logical customer of the chain-of-custody system. Each tenant has its own master HMAC key, its own Ed25519 signing keypair, and its own ledger partition. Cross-tenant signature replay is prevented by including `tenant_id` in the signing payload.

### Tenant-day
A `(tenant_id, UTC_date)` pair. The unit at which the daily Merkle seal is produced.

### Tenant key registry
A directory of `(tenant_id, key_version, public_key, valid_from, valid_until)` records. The institution publishes its public keys here; both institution and regulator hold copies. Format is implementation-flexible in v1.0; a standardized JWKS-style endpoint with FFIEC-specified extensions is a v1.x roadmap candidate.

### Threat model
The structured enumeration of adversaries the chain defends against, the integrity properties it claims, and the residual risks it accepts. See [`09-threat-model.md`](09-threat-model.md).

### TPM (Trusted Platform Module)
A small cryptographic chip on individual workstations. Sufficient for workstation authentication; not sufficient for shared, auditable, FIPS 140-2 L3 signing in a bank-grade system.

### Trust boundary / trust zone
A line across which a different trust assumption applies. The chain has three: the application process (TZ1), the bank's infrastructure (TZ2), and the regulator's verifier (TZ3). The verifier in TZ3 does not trust TZ2.

### Test vectors
Self-contained test cases that any conforming implementation must reproduce byte-for-byte. The discriminator between conforming and non-conforming implementations.

### TLA+ / Tamarin
Formal-methods tools for modeling and verifying protocols. v1.x roadmap candidates for a formal model of the threat model and the chain construction; the v1.0 conformance posture relies on the corpus and cross-implementation fuzzing for assurance, with formal modeling as additive evidence.

## V

### Verifier
The standalone offline CLI in `verifier/` that an examiner runs against a ledger snapshot. No network calls, no dependencies, single static binary. Reads the ledger and the public key as untrusted inputs and produces a defensible report. See [`07-verifier-design.md`](07-verifier-design.md).

## W

### WAL (Write-ahead log)
The source-of-truth append-only store in the ledger server. Every accepted event lands in the WAL before any other storage operation. The hot store and cold store are derived; the WAL is what the verifier reads.

## Z

### Zero-knowledge mode
A v1.x roadmap candidate verification mode where the institution proves chain knowledge to the examiner without disclosing the master IKM. v1.0 verifiers degrade to structural verification (prev_hash linking, Merkle, signature) when the IKM is not shared; under `--strict` the absent IKM elevates to FAIL.
