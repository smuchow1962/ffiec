# `ledgerctl` — institution-side admin CLI for examiner credential provisioning

Mint, revoke, and audit short-lived read-only credentials for arriving examiners. Every issuance is itself a sealed chain entry, so the credentials don't sit beside the audit trail — they live inside it.

## What this tool does, in plain English

Think of an old-fashioned notary's bound ledger.

The notary writes every consequential decision into the book by hand. Every page carries a small stamp linking it to the page before. Tear a page out and the link visibly breaks on the page that follows. At close of business, the notary presses one big wax seal across the whole day's pages with a signet ring kept in a safe. Nobody else has the ring. Nobody can fake the seal.

When an inspector arrives, the institution's senior officer fills out a **visitor pass**: who the inspector is, which rooms they may enter, what they may read, and a deadline. The officer signs the pass and stamps it. Then the officer glues the pass into the same bound ledger as one more stamped page — so the act of issuing it is part of the same tamper-evident record as everything else.

`ledgerctl` is the officer's pen and seal kit:

- `ledgerctl examiner issue` writes the visitor pass.
- `ledgerctl examiner revoke` writes the matching exit page when the inspector leaves early.
- `ledgerctl examiner audit` reads back the ledger for who got passes, when, scoped to what, signed by whom.

The inspector doesn't have to trust the institution to read the visitor pass later. The notary's signet pattern is published on the wall outside. Anyone with a photograph of a page can compare its seal to the public pattern and confirm the page is genuine. That's the design — the credentials are a courtesy to save the examiner some typing, not a trust boundary the verification depends on.

## Vocabulary, in plain English

### Sealed chain events

Every consequential thing the institution writes — every AI decision, every IAM change, every issued examiner credential — lands as one row in an append-only ledger. Each row carries the SHA-256 hash of the row before it. That's the **chain**. Tamper with any row and every row that follows shows the break.

Once a day, all of that day's rows are aggregated into a single fingerprint (a Merkle root), and that fingerprint is signed by a key held inside an HSM. That's the **seal**. The signature commits to the whole day's events at once, so flipping a single byte anywhere in the day requires forging the day's signature.

A "sealed chain event" is a row that's chained to the rows around it AND covered by the day's HSM signature. Both layers fail closed.

### Issuance key

An Ed25519 keypair the institution uses to sign credential bundles `ledgerctl` mints. When an examiner's verifier imports a bundle, it checks the signature against the institution's published issuance public key — so a tampered or fabricated bundle is caught at the laptop before the engagement starts.

The issuance key is separate from the daily-seal signing key. Different role, different rotation cadence, different blast radius. Compromising one does not compromise the other.

### Signing key (daily-seal key)

The Ed25519 keypair that signs each day's Merkle root. The private half lives inside the HSM under FIPS 140-2 Level 3 custody and never leaves. The public half is published on the institution's Compliance page. Anyone — examiner, regulator, an auditor on coffee-shop wifi — can pull a sealed day, recompute the Merkle root from the event rows, and verify the signature against the published public key.

The point of HSM custody: even the institution's own administrators can't extract this key. So even if an attacker compromises every server the institution operates, they still can't forge a backdated daily seal.

### Ed25519

A modern public-key signature scheme. Standard (NIST FIPS 186-5, RFC 8032). 32-byte private keys, 64-byte signatures, fast on every modern CPU, side-channel-resistant by design.

`ledgerctl` uses Ed25519 in two places:

- The **daily-seal signing key** (HSM-resident) signs the Merkle root once per tenant-day.
- The **issuance key** signs the credential bundles `ledgerctl` mints.

The reason for picking Ed25519 over alternatives: it's a fixed-parameter scheme. There's nothing for an operator to misconfigure. No curve choice, no padding choice, no random-nonce subtlety to get wrong. You either have a valid Ed25519 signature or you don't.

## Seal verification is unprivileged

A property of the system worth calling out, because it shapes everything else: an examiner can verify a sealed day with only the institution's published public key and the ledger bytes. No login. No tenant ID. No bearer token. No network call to the institution.

"Privileged" means you need the institution to give you something secret before you can do the thing. "Unprivileged" means you can do it with public information and the artifact in front of you. The institution cannot gate it, slow-walk it, or revoke it later.

Verifying a daily seal is two operations: recompute the Merkle root from the event rows, and check the Ed25519 signature against the published public key. Both are NIST-standard math. Both run on any laptop with no network.

The opposite design is the world where verification is privileged — the examiner logs into the institution's portal, the portal returns PASS or FAIL, and the examiner trusts the answer. That world makes the institution the oracle for its own integrity claims. The whole point of this design is to remove the institution from the trust chain.

That shapes how to read the credential bundles `ledgerctl` mints. The bundle does not grant verification capability — verification is already free. The bundle grants *fetch* capability against the live Compliance read surface so the examiner doesn't have to ask an SRE to email a ledger file every time they want a different day. Convenience plumbing. The examiner can throw the bundle away, get the ledger bytes through any other channel, and verify every sealed day on a coffee-shop wifi.

In one sentence: *the institution cannot withhold verification by withholding access.*

## Build

```bash
CGO_ENABLED=0 go build -o /tmp/ledgerctl ./cmd/ledgerctl
```

Single static binary. No CGO. Same posture as `verifier`: install on an operator workstation; no need to bring the ledger server.

## Commands

```bash
# Mint a credential bundle for an arriving examination team
ledgerctl examiner issue \
  --identity            "OCC-IT-EX-Cert-2026-Karen" \
  --tenant              tenant_acme_prod_us_east_1 \
  --period-start        2026-01-01 \
  --period-end          2026-06-14 \
  --expires-at          2026-06-22T17:00:00Z \
  --correlation-id      exam-2026-06-15-occ \
  --scope               compliance.read \
  --laptop-fingerprint  sha256:9f3e...c01a \
  --out                 ./karen.examiner.bundle

# Revoke before expiry
ledgerctl examiner revoke \
  --bundle-id  bnd_3f29b71c \
  --reason     examination_concluded

# Roster of currently-active examination credentials
ledgerctl examiner list --active

# Reconstruct the issuance/revocation history from the chain
ledgerctl examiner audit --correlation-id exam-2026-06-15-occ
```

`issue` does three things atomically:

1. Mints the credential bundle (signed JSON; see shape below).
2. Provisions the bounded role on the Compliance read surface, scoped to the named tenants and the time window.
3. Writes a `master_key.examiner_handover_received` chain entry per the receipt event in `docs/operator-guide.md`. The issuance is itself sealed evidence the examiner can verify on arrival.

`revoke` symmetrically tears down the role and writes `master_key.examiner_handover_returned`. Expiry triggers an automatic `_returned` event without an operator step.

`audit` reads the chain, not a side database. The chain is the source of truth for who got what, when.

## Credential bundle

A signed JSON file the institution hands to the examiner — by email, USB, secure share, whatever the engagement allows. The examiner imports it with `verifier identity import`.

```json
{
  "bundle_id": "bnd_3f29b71c",
  "version": "1.0",
  "issued_at": "2026-06-15T08:30:00Z",
  "expires_at": "2026-06-22T17:00:00Z",
  "issuer": {
    "institution_id": "northbridge-federal",
    "issuance_key_fingerprint": "sha256:..."
  },
  "examiner": {
    "identity": "OCC-IT-EX-Cert-2026-Karen",
    "laptop_fingerprint": "sha256:9f3e...c01a"
  },
  "scope": {
    "surface": "compliance.read",
    "tenants": ["tenant_acme_prod_us_east_1"],
    "examination_period_start": "2026-01-01",
    "examination_period_end": "2026-06-14",
    "correlation_id": "exam-2026-06-15-occ"
  },
  "endpoint": {
    "compliance_url": "https://compliance.northbridge.example",
    "tls_cert_fingerprint": "sha256:..."
  },
  "credential": {
    "type": "mtls_client_cert",
    "cert_pem": "...",
    "encrypted_key_pem": "..."
  },
  "signature_ed25519": "..."
}
```

The whole document is signed by the institution's issuance key. The examiner's verifier checks the signature on import against the institution's published issuance public key.

## Operator's-lens guarantees

- **Issuance is auditable.** Every `issue` and every `revoke` writes a chain entry under the same seal as every other consequential event. There is no separate audit log for credential operations to drift away from the chain.
- **Read-only is enforced, not advisory.** The `--scope compliance.read` flag lands as a bounded role on the read surface, not as an honor system. The role has no path to write, no path to operational systems, no path to master keys.
- **Verification is unprivileged.** Losing the bundle does not strand the examiner. `verifier verify --root-key <pub> --ledger <file>` works with no identity loaded; the bundle is convenience plumbing, not a trust boundary.
- **Expiry is automatic.** The bounded role expires on the bundle's `expires_at` timestamp without an operator action. The `_returned` event fires automatically.
- **Two-key separation.** The issuance key (signs bundles) is separate from the daily-seal signing key (signs the Merkle root). Different rotation cadence, different blast radius, different custody.

## Cross-references

- `docs/operator-guide.md` §"Examination-time master-key handover procedure" — the heavier `master-key.read` scope follows the same shape with a different role grant.
- `docs/auditor-stories/01-northbridge-federal-savings.md` — the read-only-Compliance-surface case in narrative form.
- `spec/chain-of-custody-v1.md` §7 — the verification procedure the examiner runs against pulled ledgers.
- `verifier/README.md` — the offline examiner CLI that imports the bundles `ledgerctl` mints.

## License

Apache 2.0. See `../../../LICENSE`.
