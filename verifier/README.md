# `verifier` — offline standalone CLI for examiners

Walks a tenant-day ledger, recomputes the chain hashes, recomputes the daily Merkle root, and verifies the institution's HSM-signed seal against a published Ed25519 public key. Optional `--master-key` adds a per-event MAC re-derivation pass that catches tampering individual events without touching the seal.

Single static binary. No network calls. Stdlib-only — no third-party dependencies. Deterministic reports.

## What this tool does, in plain English

Picking up the notary's bound ledger from the `ledgerctl` README:

The institution keeps every consequential decision in a bound book. Every page carries a stamp linking it to the page before. At close of business, the notary presses one big wax seal across the day's pages with a signet held in a safe. The signet's pattern is published on the wall outside the building.

`verifier` is the examiner with a magnifying glass walking the book.

- **Chain linkage** — for every page, the verifier confirms that the stamp linking it to the prior page is the right shape and that no page was lifted, swapped, or inserted.
- **Merkle root** — the verifier recomputes the day's wax-seal pattern from the page contents and confirms it matches the impression the notary actually pressed.
- **Seal signature** — the verifier compares the signet imprint on the seal to the pattern published on the wall, using only public information. The institution does not need to be in the room.
- **Per-event MAC** *(optional, with `--master-key`)* — the heavier check. The verifier re-derives the tamper-evident mark on each individual entry using the institution's master ribbon (the master IKM). If any single entry's mark doesn't match, that entry is flagged.

The institution can refuse to be present, lose its records of the visit, or dispute the conclusion later. None of that changes the math. The verifier reaches the same answer on a coffee shop's wifi as it does in the institution's data centre.

## Vocabulary, in plain English

### Chain linkage

Every entry on the ledger has a `prev_hash` field — the SHA-256 hash of the entry before it — and an `entry_hash` field — the SHA-256 of its own canonical bytes. The verifier walks the entries in order and confirms two things on each one: the entry's recomputed hash matches its stated `entry_hash`, and the next entry's `prev_hash` matches this entry's `entry_hash`.

Tamper with any entry and that entry's hash changes. The change cascades — the next entry's `prev_hash` no longer matches, and every entry after that point fails. Insertion or deletion has the same effect. The chain breaks at the spot the tamper happened, and the verifier names which entry tripped.

### Merkle root

A Merkle tree is a way to commit to a list of items with one fixed-size hash. Each item (each event payload) becomes a leaf. Pairs of leaves are hashed together into internal nodes. Pairs of internal nodes are hashed together up the tree. The single hash at the top is the **Merkle root**.

The verifier reconstructs the tree from the entries' payload bytes and confirms the root it computes matches the root the seal record claims. Flipping a single byte in any payload changes its leaf hash, which changes every internal node on the path to the root, which changes the root.

The construction follows RFC 6962 with the standard 0x00/0x01 leaf/node domain-separation prefixes. Empty tree, single leaf, and odd-count cases all match the spec.

### Seal signature

Once a day, the institution presses an Ed25519 signature over the canonical bytes of the seal record (with the signature field cleared). The signing private key lives inside an HSM under FIPS 140-2 Level 3 custody — even the institution's administrators can't extract it.

The verifier needs only the matching **public** key (32 bytes) to confirm the signature. The institution publishes that public key on its Compliance page. The verifier also confirms that the public key it was given matches the seal's `signing_key_fingerprint` — a SHA-256 of the public key bytes — so a verifier handed the wrong key fails closed instead of silently passing on a different signer.

### Per-event MAC and the master IKM

Each chain entry carries its own per-event HMAC over the event payload. The MAC key is derived per-entry from the institution's master input keying material (the **IKM**) using HKDF-SHA-256, salted with the tenant binding label and personalised with the entry id.

Without the IKM, the verifier cannot recompute these per-event MACs and reports `per-event-mac: skipped (structural-only verification)`. Structural-only is still a meaningful check — the chain, Merkle, and signature are all verified — but per-event MAC adds a finer-grained tamper-evidence layer.

The IKM handover is the procedure described in `docs/operator-guide.md` §"Examination-time master-key handover procedure". The bootstrap accepts it as a 32-byte raw file via `--master-key <path>`.

### Structural-only vs. full (key-bound) verification

Spec §7 calls these out by name. Structural verification confirms that the chain, the Merkle tree, and the seal signature all hang together — three independent integrity layers. Full verification adds the per-event MAC pass that proves each individual entry was generated by the legitimate session-key chain and not synthesised after the fact.

The verifier's report names which path ran. An examiner who needs the full path and got back `PASS (structural)` knows they didn't actually run the load-bearing per-event integrity check — and can ask the institution for the IKM handover.

## Build

```bash
CGO_ENABLED=0 go build -o /tmp/verifier ./cmd/verifier

# Cross-compile for examiner laptops
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/verifier.exe ./cmd/verifier
GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -o /tmp/verifier-mac ./cmd/verifier
GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o /tmp/verifier-linux ./cmd/verifier
```

Single static binary. No CGO. No third-party dependencies — stdlib only. The binary plus its arguments is the whole trusted-computing-base.

## Commands

```bash
# Structural verification (chain + merkle + signature). No master key needed.
verifier verify \
  --ledger    ./tenant-day.ledger \
  --root-key  ./institution.seal.pub.pem

# Full verification including per-event MACs. Requires the master IKM.
verifier verify \
  --ledger      ./tenant-day.ledger \
  --root-key    ./institution.seal.pub.pem \
  --master-key  ./ikm.bin

# Print every chain entry in order, plus the seal record.
verifier walk --ledger ./tenant-day.ledger

# Diff two ledgers; surfaces any chain or seal divergence.
verifier diff \
  --before ./yesterday.ledger \
  --after  ./today.ledger

# Generate a self-consistent demo ledger + seal-key pair (for smoke tests).
verifier gen-fixture \
  --out-ledger ./demo.ledger \
  --out-pub    ./demo.pub.pem \
  --out-ikm    ./demo.ikm
```

`verify` produces a deterministic report on stdout: one line per named step (`chain-linkage`, `merkle-root`, `seal-signature`, `per-event-mac`), each marked `[PASS]` or `[FAIL]`, with an overall verdict.

The structural-only path completes whether or not `--master-key` is supplied. The MAC step records itself as skipped when no IKM is available, so the report is always complete.

## Examiner's-lens guarantees

- **Offline.** No DNS, no HTTPS, no telemetry. Verified by build-time policy: the binary imports nothing under `net/http` for external hosts.
- **Single binary, stdlib only.** No dynamic linking. No external configuration files except those passed on the command line.
- **Deterministic output.** A given ledger plus a given root-key produces a byte-for-byte identical report. Two examiners run the verifier independently and compare reports by hash.
- **Failure surfaces are loud.** Every failure is named, localised to the entry that tripped, and accompanied by the recomputed value next to the claimed value. The verifier never silently passes a partial verification.
- **Independent of the institution.** The verifier reads only the ledger bytes and the published public key. It does not call the institution's vendor, the institution's HSM, or any third party. The institution's only role is to hand over the ledger.
- **Fingerprint pinning.** The verifier confirms the supplied public key matches the seal's `signing_key_fingerprint` before accepting the signature. A verifier pointed at the wrong public key fails closed instead of silently validating a different signer.

## Cross-references

- `spec/chain-of-custody-v1.md` §7 — the verification procedure this tool implements.
- `docs/operator-guide.md` §"Examination-time master-key handover procedure" — the IKM-handover ceremony that produces the `--master-key` input.
- `ledger/cmd/ledgerctl/README.md` — the institution-side admin CLI that mints examiner credentials.
- `testkit/clitest` — the harness this tool's CLI tests run against.

## License

Apache 2.0. See `../LICENSE`.
