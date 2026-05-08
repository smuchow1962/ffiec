# Verifier deployment package

> **What this doc is.** The information the regulator's IT shop needs to allowlist and deploy the verifier on examiner laptops.

## Verifier binary specifications

| Field | Value |
|---|---|
| Binary name (Linux/macOS) | `verifier` |
| Binary name (Windows) | `verifier.exe` |
| File type | Static ELF / Mach-O / PE binary, no dynamic linking |
| Approximate size | ~12 MB stripped |
| Network behavior | None — verifier makes no outbound connections |
| File system reads | `--ledger`, `--root-key`, `--bundle <path>` (read-only) |
| File system writes | `--report`, `--json-report`, `--bundle <path>` (write only to specified paths) |
| Privileged operations | None — runs as ordinary user |
| Crash behavior | Exits with non-zero status; no telemetry leaves the machine |
| Memory footprint | Streaming construction; O(log N) hashes plus event-row working set |

## Trust artifacts per release

Every release publishes:

1. **Binaries** — one per platform target (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
2. **Cosign signatures** — `.sig` and `.cert` per binary
3. **GPG-signed hash manifest** — `verifier-v<version>.sha256.asc`, signed by the project's GPG key
4. **CycloneDX SBOM** — `verifier-v<version>.cdx.json`
5. **Vulnerability scan report** — `verifier-v<version>.scan.json` (trivy + grype output)
6. **Source tarball** — `verifier-v<version>.src.tar.gz` for reproducible-build verification
7. **Validate wrapper scripts** — `verifier-validate.sh` (POSIX) and `verifier-validate.ps1` (Windows)

## Trust anchors

Cached out-of-band by the regulator's IT shop, validated against fingerprints published in the spec text:

| Anchor | Purpose | Fingerprint location |
|---|---|---|
| Project cosign public key | Validates binary signatures | `spec/trust-anchors.md` (canonical) |
| Project GPG public key | Validates the hash manifest (cosign fallback) | Same |
| Spec-version-pinned key | Per spec version, the trust anchors valid for that version | Same |

## Allowlisting

Two acceptable postures:

**Conservative (SHA-256 allowlist).** The regulator's IT shop allowlists the binary by SHA-256 hash. Each new release requires a manual allowlist update. Lower automation, higher control.

**Signature allowlist (cosign-based).** The regulator's IT shop allowlists by cosign signature against the pinned project key. New releases auto-allowlist on signature validation. Higher automation, requires institution-side cosign verification capability.

Both are conformant with FFIEC IT shop standards. Many regulators run both — signature allowlist for routine examinations, hash allowlist for high-sensitivity engagements.

## Operational behavior

The verifier on examiner laptops:

- **No telemetry.** The verifier emits no logs to external services.
- **No network access.** The verifier opens no sockets. (Verifiable: `strace` or equivalent shows no `socket()` syscalls.)
- **Read-only ledger access.** The verifier opens the ledger snapshot in read-only mode. (Verifiable: examiner can mount the snapshot from a read-only filesystem.)
- **Deterministic output (with `--deterministic`).** Two examiners produce identical PDFs from identical inputs.

## Sample command lines

The `--master-key` flag points at the institution's IKM file (32 raw bytes; file mode 0600 on POSIX, ACL restricted to the examiner account on Windows). Without `--master-key`, the verifier performs structural verification only and skips per-event HMAC equality (under `--strict` the absent IKM elevates to FAIL); the institution provides the IKM at examination time per a documented disclosure shape (protective-order disclosure or HSM-mediated derivation per `legal-disclosure.md`).

```
# Standard verification (with key-bound verification — recommended)
verifier verify --ledger ./ledger-snapshot --root-key ./tenant.pub \
  --master-key ./tenant-ikm.bin \
  --report ./acme-2026-april.pdf --json-report ./acme-2026-april.json

# With strict mode (conservative — substantive SOC engagements use this)
verifier verify --ledger ./ledger-snapshot --root-key ./tenant.pub \
  --master-key ./tenant-ikm.bin \
  --report ./acme-2026-april.pdf --strict

# Bundle output for working papers
verifier verify --ledger ./ledger-snapshot --root-key ./tenant.pub \
  --master-key ./tenant-ikm.bin \
  --bundle ./acme-2026-april-bundle.tar.gz

# Structural-only verification (when IKM is not yet provided — initial review)
verifier verify --ledger ./ledger-snapshot --root-key ./tenant.pub \
  --report ./acme-2026-april-structural.pdf
# Reports as PASS-WITH-ANOMALY: structural verification only; key-bound verification skipped

# Walk a single run for incident investigation
verifier walk --ledger ./ledger-snapshot --run-id r_a3f29b71c \
  --master-key ./tenant-ikm.bin

# Compare two snapshots
verifier diff --before ./ledger-2026-04-01 --after ./ledger-2026-04-02
```

## Validate-before-run script

Examiners run `verifier-validate.sh` before running the verifier. The script:

1. Verifies the cosign signature against the cached project public key
2. Verifies the GPG-signed hash manifest matches the binary
3. Optionally re-runs reproducible-build verification if the institution's build pipeline is available
4. Exits 0 only if all checks pass

```
$ verifier-validate.sh ./verifier
[1/3] Verifying cosign signature...   PASS
[2/3] Verifying GPG hash manifest...  PASS
[3/3] Reproducible-build check...     SKIPPED (no build pipeline configured)
verifier-validate: OK to run

$ ./verifier verify --ledger ./snapshot --root-key ./tenant.pub --report ./out.pdf
```

If any step fails, the wrapper exits non-zero and the examiner does not run the verifier.

## Deployment cadence

The project follows a 30-day pre-announcement window for new versions. The regulator's IT shop updates the allowlist between announcement and release. Existing examinations using older versions continue without disruption; a passing report from an older verifier remains a valid working-paper artifact.
