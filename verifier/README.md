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

# Walk an on-disk NDJSON audit file, as a chain sink writes it during capture.
# With no key this is witness mode: steps 1-6, and no MAC check.
verifier verify --audit-file ./custody.ndjson

# The same walk with the key: steps 1-9. --master-key-hex reads hex text; --master-key reads raw bytes.
verifier verify --audit-file ./custody.ndjson --master-key-hex ./custody.ikm

# The same walk with the day's seal: steps 1-11. Without the key: steps 1-6, 10 and 11.
# Take --root-key from a source other than whoever handed you the file.
verifier verify   --audit-file      ./custody.ndjson   --seal            ./seal-2026-11-04.json   --root-key        ./institution.seal.pub.pem   --master-key-hex  ./custody.ikm   --posture         ed25519

# The same run, where a posture problem fails. It needs the key, the seal and --posture.
verifier verify   --audit-file      ./custody.ndjson   --seal            ./seal-2026-11-04.json   --root-key        ./institution.seal.pub.pem   --master-key-hex  ./custody.ikm   --posture         ed25519   --strict

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

Try the seal on the published sample. Its keys come from published test seeds, so it proves nothing about any institution, and the output says so on its `Anomaly:` lines.

```bash
cd cmd/verifier/testdata
verifier verify \
  --audit-file      simulator-custody.golden \
  --seal            simulator-custody-seal.json \
  --root-key        simulator-custody-seal-pub.pem \
  --master-key-hex  simulator-custody-key.hex
```

The seal file is one JSON object with the field names of the `audit_file.seal` block in the conformance vectors. Step 11 here means the signature, the `key_versions` cross-check, and for `v1.0b` and `v1.0c` seals the KMS handle digest. A `v1.0c` seal's root leaves out `chain_kind = "operational"` entries, and step 11 checks those entries against the operational-events root on every run, so cutting them all still fails. An `Anomaly:` line names a check that did not run: a seal with no key version list skips that cross-check. `--root-key` is the only source of the public key. No field in the audit file or the seal can supply it, so a reader who takes the key from the file's holder has checked nothing.

`verify --audit-file` prints the §7 normative output and uses the §10.12 exit codes: `0` for a pass, `1` for a failed step, `2` for a file the walk cannot begin, `3` for a key it cannot read, a flag rule, or a seal it cannot judge. A pass lists under `Anomaly:` what the walk did not establish. With no `--seal`, steps 10 and 11 do not run, and the first anomaly says that a cut at the end of the file cannot be detected. The §7 `Verdict-Object` line is not printed yet.

The structural-only path completes whether or not `--master-key` is supplied. The MAC step records itself as skipped when no IKM is available, so the report is always complete.

## Examiner's-lens guarantees

- **Offline.** No DNS, no HTTPS, no telemetry. Verified by build-time policy: the binary imports nothing under `net/http` for external hosts.
- **Single binary, stdlib only.** No dynamic linking. No external configuration files except those passed on the command line.
- **Deterministic output.** A given ledger plus a given root-key produces a byte-for-byte identical report. Two examiners run the verifier independently and compare reports by hash.
- **Failure surfaces are loud.** Every failure is named, localised to the entry that tripped, and accompanied by the recomputed value next to the claimed value. The verifier never silently passes a partial verification.
- **Independent of the institution.** The verifier reads only the ledger bytes and the published public key. It does not call the institution's vendor, the institution's HSM, or any third party. The institution's only role is to hand over the ledger.
- **Fingerprint pinning.** The verifier confirms the supplied public key matches the seal's `signing_key_fingerprint` before accepting the signature. A verifier pointed at the wrong public key fails closed instead of silently validating a different signer.

## Several files and seals

`--audit-file` and `--seal` repeat. A day with 2 runs takes 2 files, because a run does not span files. A run that crosses midnight takes both days' seals.

```bash
verifier verify   --audit-file run-a.ndjson   --audit-file run-b.ndjson   --seal seal-2026-11-04.json   --seal seal-2026-11-05.json   --root-key institution.seal.pub.pem   --master-key-hex custody.ikm   --posture ed25519
```

An entry's seal is the UTC date of its `received_at`, plus one day when `late_binding` is true. In the file, both fields sit in the `audit_chain` block, outside the MAC. Each seal covers the entries of its date, and every entry must fall under a supplied seal.

| Input | Result |
|---|---|
| A file for the day is missing, or an entry was cut | Exit 1 at step 10: the root does not match |
| An entry's date matches no supplied seal | Exit 1 at step 10: `entry outside every supplied seal` |
| 2 different files claim one run | Exit 1 at step 10: the §10.25 fork reason. Steps 10 and 11 do not run |
| 2 byte-identical files, 2 seals for one date, or `received_at` on some entries only | Exit 2, before step 1 |
| Several seals, and no entry has `received_at` | Exit 2, before step 1. With one seal, all entries count as its day, and the root decides |
| A seal that is not daily | Its signature is checked. Exit 1 if it is bad, otherwise exit 3 |

A record with `format_version` and none of `seq`, `prev_hash`, `payload_hash`, `key_version` and `key_fingerprint` has no leaf, no MAC and no root over it. It is outside the verified evidence, and the anomaly line `unstamped records: N` counts it. A record with some of the 5 fails at step 6. An added unstamped record stays exit 0, with `--strict` too: read its count.

Late-binding entries print as `late-binding entries: N`, with each entry's `(run_id, seq, received_at, captured_at)`.

**What a pass covers.** A sealed run prints `seal dates covered:`. The result is about those days only. If a run crosses midnight, and its later entries are cut from the end along with that day's seal, every supplied entry is covered and no step can fail. Compare the dates with the days you expected.

**`kms_handle_uri` is provenance, outside the MAC.** On a `v1.0a` seal, no signature binds it, so an editor can change it. That includes clearing a `plaintext-` handle before `--strict` reads it. `v1.0b` and `v1.0c` bind the day's distinct values, so an edit that changes that set fails, and removing one entry's copy of a value others still hold does not. `dev_mode` is signed from `v1.0a` on.

## Algorithm posture

`--posture` lists the signature algorithms the institution declares. Take the list from the institution's CC8.1 control description, and never from the seal. A seal field can be edited, and no signature binds it.

| Input with a seal | Result |
|---|---|
| No `--posture` | `Status: PASS-WITH-ANOMALY`, exit 0: a stripped signature cannot be detected. A witness run keeps its own status and prints the anomaly line. |
| No `--posture`, with `--strict` | Exit 3 |
| The list names the seal's one algorithm | `Status: PASS`, in both modes |
| The list names a second algorithm, and the seal lacks it | `Status: PASS-WITH-ANOMALY`, exit 0. With `--strict`: exit 1 |
| The seal names an algorithm the list lacks | `Status: PASS-WITH-ANOMALY`, exit 0. With `--strict`: exit 1 |
| The seal carries a `signatures` list that the posture agrees with | Exit 3: this verifier cannot check a second signature |

Exit 3 here does not show that a second signature exists. The list is unsigned, so anyone who can edit the seal file can add an entry with no key, and that moves a stripped seal from exit 1 to exit 3. Before you rely on such a seal, obtain a verifier that checks the second signature.

Exit 0 alone does not show a clean result. Read the `Status:` and `Anomaly:` lines.

`--strict` also refuses, with exit 1 at step 12, a seal with `dev_mode` true and a seal over any entry whose `kms_handle_uri` begins with `plaintext-`. Without `--strict`, these print anomaly lines. `--strict` also refuses a seal with no `sign_payload_version`: that form signs no `dev_mode`, so nothing shows it is a production seal.

## Sign-payload floor

`--min-sign-payload-version v1.0b` fails a seal below that version with exit 1, before signature dispatch. Take the value from the institution's CC8.1. A seal's own `declared_sign_payload_floor` field is also honoured, but no signature binds it, so an editor can remove it. Use the flag. Only the flag protects against that edit. If the flag repeats, the last one wins, so a script must not append a weaker one. An empty value is exit 3.

## Limits against spec 0.4.0

The spec published first, and this verifier follows it in steps. It does not yet do the following:

1. **Hourly and weekly seals.** No signed field names their window, so a run with one exits 3 after its signature verifies. The spec tracks this as an open issue.
2. **Dual-algorithm signatures.** Cases (a) and (e) of §7 step 11 are not built, so such a seal is exit 3, and exit 3 proves nothing about its list. Go 1.26.2's standard library has no ML-DSA package (`go doc crypto/mldsa` finds none).
3. **Posture sources.** Only the command-line list is built, which is source (iii) of case (c).
4. **Unknown record kinds.** A stamped record of a later kind is walked and covered by the roots. The `unknown wire-format kind present` line and the `additional_verifications` marker are not printed.
5. **Keys.** One IKM per run, for the one `--key-version`. A set whose entries use several generations fails at step 7.
6. **`--ledger` walk.** It applies no posture and no floor, takes one seal inside its file, and refuses `--posture`, `--strict` and `--min-sign-payload-version`.
7. **`--strict` rules not built.** Step 12a (the GenAI model identifier check) does not run, and `seal.cadence` is not compared with a claimed cadence, because no input names one. A `--strict` run prints both as anomaly lines. The tenant key registry and flat-file floor sources are not built.

## Cross-references

- `spec/chain-of-custody-v1.md` §7 — the verification procedure this tool implements.
- `docs/operator-guide.md` §"Examination-time master-key handover procedure" — the IKM-handover ceremony that produces the `--master-key` input.
- `ledger/cmd/ledgerctl/README.md` — the institution-side admin CLI that mints examiner credentials.
- `testkit/clitest` — the harness this tool's CLI tests run against.

## License

Apache 2.0. See `../LICENSE`.
