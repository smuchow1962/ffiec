# `ledger` — institution-side reference ingest server

The server that receives OTLP-shaped chain entries from instrumented applications, verifies the HMAC chain on ingest, writes to append-only storage, computes the daily Merkle seal, and signs that seal under HSM custody.

This is a **bootstrap**. The server loop, OTLP receivers, storage backends, and HSM clients are scaffolded — the dispatch shape, config schema, and validation surface are real and exercised by the test suite. `serve` produces a readiness summary today; the actual receive loop lands in follow-up work.

## What this tool does, in plain English

Picking up the notary's bound ledger from `ledgerctl`'s README:

`ledger` is the **notary and the clerks**. The server stands at the inbound side. When an instrumented application emits a chain entry — a sealed AI decision, an IAM change, a customer interaction — the entry arrives over the wire (OTLP) at the server's receiver. The server verifies the per-event MAC, walks the proposed prev-hash linkage, and only then writes the entry into the append-only ledger.

At end of day, the server stops accepting new entries for the current day, builds the day's Merkle tree, and signs the resulting root with the institution's signet held inside an HSM. That signed seal record is the artifact `verifier` later checks.

The bound book itself is the **append-only storage** — files on disk under a directory the operator configures, or a Postgres database with INSERT-only role grants. Whichever. The contract is the same: writes never modify or delete prior entries. Schema-level role grants enforce that even an attacker with database administrator access cannot edit a past row.

The signet ring lives in the **HSM** — a hardware module that signs but does not export the private key. Even the institution's own administrators cannot extract it. The bootstrap server config includes a `software` HSM driver for development, which is explicitly forbidden in production mode.

This CLI is the operator's interface to that server: validate config, write a default config, dry-run startup, run a health probe.

## Vocabulary, in plain English

### OTLP (OpenTelemetry Line Protocol)

The wire format chain entries arrive on. OTLP is the standard OpenTelemetry transport: gRPC by default, HTTP/JSON as an alternative, both ports configurable. The institution's instrumented apps already speak it for traces and metrics; the chain extension uses the same encoding so the institution's network plumbing carries audit events alongside everything else without a separate channel.

### Append-only writer

The storage layer that physically receives writes. The contract: only INSERT, ever. Never UPDATE, never DELETE. The bootstrap supports a directory driver (NDJSON files per tenant-day) and a postgres driver (DSN to a database where the writer role's grants are INSERT-only on the chain table). Both enforce append-only at a layer the application code can't bypass.

If a chain entry arrives whose `prev_hash` does not match the current tail of the chain, the writer rejects it and records the rejection in a side channel. Rejections do not corrupt the ledger because they were never committed.

### Daily Merkle seal

At end of day for a given tenant, the server collects every entry written during that day, hashes each entry's payload into a Merkle leaf, builds the tree, and produces a single 32-byte root that commits to all of them. That root, plus a few metadata fields (date, leaf count, signing-key fingerprint), becomes the **seal record** — one row per tenant-day.

The seal is what the verifier checks. Recompute the tree from the entries, compare the root to the seal record's claim, and verify the signature on the canonical seal bytes. Three independent integrity layers (chain, Merkle, signature) all hang together.

### HSM-rooted custody

The seal record is signed by an Ed25519 keypair whose private half lives inside a hardware security module (HSM) under FIPS 140-2 Level 3 custody. The HSM signs but does not export — the private key bytes never leave the device. The public key is published on the institution's Compliance page.

The bootstrap config supports three HSM drivers:

- **pkcs11** — production, talks to an on-prem HSM via the PKCS#11 module path.
- **cloud-kms** — production, talks to a managed key (AWS CloudHSM, Azure Managed HSM, GCP Cloud HSM).
- **software** — development only. The server refuses to start with `production_mode: true` and `driver: software` so a misconfigured deployment fails closed at boot.

## Build

```bash
CGO_ENABLED=0 go build -o /tmp/ledger ./cmd/ledger
```

Single static binary. No CGO. Stdlib-only.

## Commands

```bash
# Write a default config to disk (refuses to overwrite without --overwrite).
ledger config init --out ./ledger.json

# Validate a config file. Prints a one-line summary on success.
ledger config validate --config ./ledger.json

# Run a readiness probe against config (no network, no server start).
ledger health --config ./ledger.json

# Print readiness summary and exit; useful as a deploy smoke.
ledger serve --config ./ledger.json --dry-run

# (Bootstrap) Start the server. Currently exits with "not yet implemented".
ledger serve --config ./ledger.json
```

`config validate` prints one summary line — the version, storage driver, HSM driver, production-mode flag, log level. `serve --dry-run` prints a multi-line readiness block with the OTLP listen addresses, storage settings, and HSM driver. `health` is the silent pass/fail check intended for orchestration probes.

The server's current `serve` (without `--dry-run`) is intentionally noisy: it prints readiness, then exits 1 with a "not yet implemented" message. That keeps the shape honest — an operator can see exactly where the bootstrap stops and what config it would have used.

## Config shape

```json
{
  "version": "1.0",
  "otlp": {
    "grpc_addr": "0.0.0.0:4317",
    "http_addr": "0.0.0.0:4318"
  },
  "storage": {
    "driver": "directory",
    "path": "/var/lib/ffiec-ledger"
  },
  "hsm": {
    "driver": "software",
    "production_mode": false
  },
  "logging": {
    "level": "info"
  }
}
```

`storage.driver` is `directory` (NDJSON files) or `postgres` (DSN-driven). `hsm.driver` is `pkcs11`, `cloud-kms`, or `software`. Every field is validated at load time; misconfigured deployments fail closed at boot rather than half-starting.

The DSN is redacted in `serve --dry-run` output (printed as `<set>` rather than the full string) so operator screen-shares don't leak credentials.

## Operator's-lens guarantees

- **Fail-closed config validation.** Every field has a known set of values. Unknown driver names, missing required fields, and the production-mode + software-HSM combination all reject at load time.
- **No half-start.** `serve` validates the entire config before opening any port or touching storage. There is no path where the receiver is up but the storage backend is misconfigured.
- **DSN redaction.** `serve --dry-run` prints `dsn="<set>"` rather than the literal DSN, so a captured terminal log never contains storage credentials.
- **Health is silent on success.** `ledger health --config <path>` prints `ok` and exits 0 when the config validates. Suitable as a Kubernetes startup probe or a deploy smoke step.
- **Single binary, stdlib only.** Same posture as the `verifier` — no CGO, no third-party dependencies, deterministic output.

## Cross-references

- `ledger/README.md` — module-level architecture overview.
- `verifier/README.md` — the offline tool that consumes the seal records this server signs.
- `ledger/cmd/ledgerctl/README.md` — the institution-side admin CLI for examiner credentials.
- `spec/chain-of-custody-v1.md` — the wire format and chain construction rules the receiver enforces.

## License

Apache 2.0. See `../../../LICENSE`.
