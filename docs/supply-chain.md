# Supply chain

> **What this doc is.** Documentation of the verifier and ledger-server build pipeline, signing artifacts, and trust paths. Cybersecurity examiners and SOC teams use this to evaluate the institution's exposure to supply-chain risk.

## What ships per release

Every release of the verifier and ledger server publishes:

| Artifact | Format | Purpose |
|---|---|---|
| Binary (per platform) | Static ELF / Mach-O / PE | The executable |
| Cosign signature | `.sig`, `.cert` | Primary trust anchor |
| GPG-signed hash manifest | `.sha256.asc` | Fallback trust anchor |
| CycloneDX SBOM | `.cdx.json` | Vulnerability management input |
| Vulnerability scan report | `.scan.json` | Trivy + Grype output |
| Source tarball | `.src.tar.gz` | Reproducible-build input |
| Validate wrapper | `verifier-validate.sh` / `.ps1` | Runs all validation steps |
| **Spec PDF** | `chain-of-custody-v{version}.pdf` | The normative spec text, content-addressable; cosign-signed and GPG-signed alongside the binary |
| **Spec content-hash manifest** | `spec-manifest.sha256.asc` | SHA-256 of the spec PDF + each design doc + each test-vector file; GPG-signed; institutions rebuild the corpus from the spec text and confirm against this manifest |
| **Conformance corpus tarball** | `test-vectors-v{version}.tar.gz` | Full `spec/test-vectors/` content as a versioned bundle; cosign-signed and GPG-signed (corpus subversion defense) |
| **SLSA provenance attestation** | `verifier.intoto.jsonl` | SLSA Level 3 provenance per the in-toto attestation format; institution can verify build was produced by the declared pipeline from the declared source |

For containers:

| Artifact | Notes |
|---|---|
| Image | Distroless base; minimal attack surface |
| Image signature | Cosign signature on the image |
| Image SBOM | CycloneDX 1.5 |
| Image scan report | Trivy + Grype |

## Build pipeline

The pipeline runs in a hermetic environment. Properties:

- **Hermetic.** No external dependencies at build time. Source, toolchain, and dependencies are pinned.
- **Reproducible.** Independent rebuilds from the same source produce byte-for-byte identical artifacts.
- **Audited.** Build logs are retained; the pipeline's operations are subject to the project's standard audit cadence.

### Reproducibility

`CGO_ENABLED=0 go build -ldflags="-s -w -buildid=" -trimpath -o verifier ./cmd/verifier`

Flags:

- `CGO_ENABLED=0` — pure Go; no dynamic linking
- `-s -w` — strip debug symbols and DWARF info
- `-buildid=` — clears the build ID for reproducibility
- `-trimpath` — removes filesystem paths from the binary

Independent rebuilds with the same Go version and the same source produce byte-for-byte identical output.

### Toolchain pinning

Each release pins a specific Go version. The pinned version is part of the SBOM. Institutions performing reproducible-build verification use the same pinned version.

## Signing roles

The project's release signing is split across two roles:

### Cosign signing

A hardware-backed cosign key, held by the project's release-management role. The key signs every binary and every container image at release time. Sigstore is the trust root.

### GPG signing (fallback)

An offline-held GPG key, held by a *separate* small group within the project's release-management role from the cosign-key holders (the separation prevents a single compromise from collapsing both trust paths). The key signs the hash manifest at release time. The institution caches the GPG public key out-of-band as a fallback if Sigstore is compromised.

**GPG key rotation cadence.** The GPG key is offline-held but still requires a rotation plan; offline keys decay differently from online keys but are not exempt from rotation. Project-side cadence:

- **Default rotation: every 36 months** (3 years), aligned with NIST SP 800-57 guidance for offline signing keys.
- **Emergency rotation:** within 7 days of any credible compromise indicator, or within 30 days of a project-side governance event (key holder departure, hardware-token end-of-life).
- **Subkey rotation:** every 12 months for the signing subkey under the master GPG key (the master key remains on the same 36-month cadence).

Each rotation: project publishes the new GPG public key in the spec text (the next spec version), publishes the GPG public-key fingerprint in the release notes, and signs a transition statement (`KEY-ROLLOVER-{old_keyid}-to-{new_keyid}.asc`) under both the old and new keys for the institution's transition-validation.

### Cosign-key recovery procedure

If the cosign key is compromised (key material exfiltrated, holder credentials breached, hardware token lost):

1. **Project response (≤ 24 hours of confirmed compromise).** The release-management role revokes the cosign certificate at Sigstore (using Fulcio's revocation procedure if applicable) and publishes a `KEY-REVOCATION-{keyid}.asc` notice signed by the project's GPG key.
2. **Institution notification (≤ 48 hours).** The project publishes the revocation notice on the spec website, the GitHub release page, the project's announcements channel, and to the FFIEC working group's notification list. Institutions are notified through the same channels their normal release-notification arrives on.
3. **New key authentication chain (≤ 7 days).** The project provisions a new cosign key, publishes its public-key fingerprint in the spec text and as a `KEY-INTRODUCTION-{new_keyid}.asc` notice signed by the project's GPG key. The new key's first signed binary is the next release.
4. **Institution-side action.** Institutions remove the revoked cosign key from their trust anchors, install the new cosign key (validating against the GPG-signed introduction notice), and re-validate any binaries signed under the new key. Binaries signed under the revoked key remain verifiable for historical audit purposes (old signatures still verify against the revoked-but-archived key) but new deployments require the new key.

The procedure is documented in `GOVERNANCE.md` at the project level. The institution's IR playbook references it in Scenario 3 (signature verification failed) when the failure is at the binary-signature layer rather than the seal-signature layer.

### GPG-key recovery procedure

Symmetric to the cosign-key recovery procedure, with the dual-compromise edge case explicitly named. If the GPG signing key is compromised:

1. **Project response (≤ 24 hours of confirmed compromise).** The release-management role publishes a `KEY-REVOCATION-{gpg_keyid}.asc` notice signed by the project's cosign-signed-via-Sigstore release-management identity. Note: the GPG revocation notice is signed by cosign in this case, because the GPG key itself is the compromised artifact.
2. **Institution notification (≤ 48 hours).** Same notification channels as cosign-key recovery (spec website, GitHub release page, project announcements, FFIEC working group's notification list).
3. **New GPG key authentication chain (≤ 7 days).** The project provisions a new GPG key (typically held by the same separated small group, possibly with new key holders if compromise involved a key-holder breach). The new key's public-key fingerprint is published in the spec text and as a `KEY-INTRODUCTION-{new_gpg_keyid}.asc` notice signed by the project's cosign identity. The first GPG-signed manifest under the new key is the next release.
4. **Institution-side action.** Institutions remove the revoked GPG key from their trust anchors, install the new GPG key (validating against the cosign-signed introduction notice), and re-validate any historical manifests as needed for ongoing examinations.

**Dual-compromise edge case (cosign AND GPG simultaneously compromised).** The two-trust-anchor independence claim relies on cosign-key holders and GPG-key holders not being subverted simultaneously. If both keys are confirmed compromised in the same window:

- The project's only remaining trust anchor is the regulator-held public-key fingerprint (which is the cosign public key as recorded with the regulator at tenant onboarding). Per Adversary I, three independent compromises must align for forged data to silently pass.
- The project escalates: published recovery notice signed by an out-of-band emergency GPG key the release-management role pre-provisioned for this scenario (the "cold-disaster-recovery key" maintained by yet another separated small group, never used until a dual-compromise event).
- The institution's emergency response: stop deploying new releases of the verifier and ledger server until the project completes the new-key authentication chain. Continue using the historical (pre-compromise) verifier binary against historical seals; the historical binary's signatures remain validatable against the pre-compromise cached cosign public key.
- The institution coordinates directly with its primary regulator on examination posture during the multi-week dual-compromise recovery window.

**Cold-disaster-recovery key lifecycle (project-side governance).** The cold-DR key has its own lifecycle independent of the cosign and GPG keys it backs:

- **Rotation cadence: every 60 months** (5 years), aligned with NIST SP 800-57 guidance for cold-storage signing keys. Longer than the 36-month standard-GPG cadence because the cold-DR key is air-gapped and used only for emergencies.
- **Key-holder lifecycle:** the cold-DR key is held by a third separated small group within the project's release-management role (distinct from cosign-key holders and standard-GPG-key holders). Key-holder rotation per departure / governance event, not per key rotation; emergency rotation within 30 days of any key-holder departure.
- **Annual dry-run cadence:** the project performs a documented annual dry-run exercise to confirm the cold-DR key's hardware is functional, the air-gap procedure works, the key holders can be reached within target timelines (24 hours), and the published recovery procedure produces a successfully-validated transition statement. The dry-run output (without the actual key material) is published as `KEY-DR-DRYRUN-{year}.asc` signed under the standard cosign or GPG key.
- **Health check:** the project records (without exposing the key bytes) the cold-DR key's hardware health (HSM tamper detection state, hardware-token battery indicator if applicable) on the same annual cadence; deviations trigger out-of-band notification to the FFIEC working group.

**Institution-side consumption of cold-DR-key dry-run attestation.** The institution's release-validation procedure includes annual consumption of the `KEY-DR-DRYRUN-{year}.asc` attestation:

1. The institution downloads the year's dry-run attestation from the project's release page.
2. The institution validates the attestation's signature against the cached standard cosign or GPG public key (whichever signed the attestation per the project's procedure).
3. The institution archives the attestation in its control-evidence repository for the year (retention per the institution's chain-event retention period).
4. Failure to obtain the attestation within 30 days of the project's annual dry-run window triggers IR Scenario 11 (project-side trust-anchor degradation) with branch (c) "missed dry-run window" disposition.
5. The institution's `KEY-DR-DRYRUN-{year}.asc consumption log` is the operational evidence (per CSF-2.0 GV.SC-04 binding); SOC and examination teams sample-test the consumption log against the institution's archived attestations.

This is the institution-side consumption procedure for the project-side cold-DR-key lifecycle. Without it, the institution treats the cold-DR fallback as documented-but-unexercised, which a GV.SC-04 examiner would challenge as un-tested control evidence.

The cold-DR key's procedures are documented in `GOVERNANCE.md` §"Crypto Key Disaster Recovery." Institution-side reviewers consume the DR-DRYRUN attestations as evidence the dual-compromise fallback is operationally exercised, not just documented. Dual-compromise is a project-side governance event; the institution-side response is documented in IR Scenario 3 (signature verification failed) under the "binary-signature variant" with a pointer to this section.

### Cosign-key recovery during an active examination

If the cosign-key revocation notice arrives mid-examination, when the FFIEC EIC is on-site validating verifier signatures:

1. **The bank's verifier remains usable against historical seals.** The historical binary's signature was valid at deployment time, recorded in the bank's reproducible-build log. Past verifications using the historical verifier remain authoritative.
2. **The bank pauses new-version verifier deployment** until the new cosign key is available and validated. The bank's deployment-gating control (per the SLSA L3 institutional consumption section above) refuses to deploy the next release until the new cosign key is on the institution's trust-anchor list.
3. **If the examination is mid-engagement,** the bank communicates the supply-chain incident to the EIC immediately, offers the GPG-fallback validation as the interim posture, and continues the examination using the GPG-validated binary. The bank's `bundleverify` script (per `07-verifier-design.md` §8.5) supports the GPG-only validation mode for this case; the EIC re-runs the bank's verifier outputs through `bundleverify --gpg-only` to confirm the GPG fallback validates.
4. **Post-recovery,** when the new cosign key is published and the bank installs it, the bank re-runs the affected verifier outputs through the new-cosign-key-validated binary to confirm the verification results are stable across the trust-anchor change. The examination's working papers record both the GPG-fallback and post-cosign-recovery validation results for evidence completeness.

This is the institution's operational continuity posture during the cosign recovery window. It is documented in IR Scenario 3 (signature verification failed) and referenced from this section of the supply-chain doc.

## Trust path summary

**Primary path.**
- Institution receives binary
- Institution validates cosign signature against pinned project public key
- Sigstore validates the signing certificate

**Fallback path (if Sigstore is compromised).**
- Institution validates GPG-signed hash manifest against cached project GPG public key
- Institution computes binary SHA-256 and compares to the manifest

**Defense-in-depth (additional).**
- Institution rebuilds from source using the same toolchain version
- Institution compares the rebuilt binary's SHA-256 to the published binary
- Institution archives the rebuild log as control evidence

## Trust anchors

The institution caches:

| Anchor | Purpose | Where to obtain |
|---|---|---|
| Cosign public key (per spec version) | Primary signature validation | Spec text, regulator, project website |
| GPG public key (per spec version) | Fallback signature validation | Same |
| Cold-DR public key (per spec version) | Dual-compromise fallback (only used during emergency dual-compromise recovery; never used in normal operation) | Spec text + project's GOVERNANCE.md; cached out-of-band by the institution's release-validation procedure |
| Toolchain version pinning | Reproducible-build verification | Per-release release notes |
| `slsa-verifier` public key + binary SHA-256 | SLSA L3 provenance attestation validation (deployment-gating per "Institutional consumption of the SLSA L3 attestation") | Co-published with the verifier release; cached out-of-band |

The project publishes the cosign and GPG public-key fingerprints in the spec text itself, so the trust anchors are validated against the spec rather than against a separate distribution channel.

## Threat model

Threats addressed by the supply chain:

- **Compromised binary.** Defense: cosign signature validates the binary was signed by the project. Sigstore compromise is mitigated by the GPG fallback.
- **In-flight modification of the binary.** Defense: the institution validates the signature on the binary it received; modifications produce a signature mismatch.
- **Substitution of the project's signing key.** Defense: institution caches the public key out-of-band; substitution requires compromising both the project's key infrastructure and the institution's cached key.
- **Compromise of the project's release pipeline.** Defense: reproducible builds let independent parties (institutions, auditors, security researchers) rebuild from source and confirm the published binary matches.

Threats not addressed:

- **Compromise of an institution's local trust anchors.** The institution's IT shop is responsible for protecting its cached cosign and GPG public keys. The chain spec assumes the institution operates standard IT-shop controls.
- **Side-channel attacks against the verifier on the examiner's laptop.** Examiner-laptop hygiene is the regulator's responsibility.

## Container supply chain

The verifier Docker image has additional considerations:

- Distroless base image (currently Google distroless or Chainguard wolfi)
- No package manager in the image; dependencies are statically linked
- SBOM published as CycloneDX 1.5
- Trivy and Grype scan results published per release
- Image signed with cosign

The institution's container registry mirror pulls the image, validates the signature, scans for vulnerabilities, and publishes to the institution's deployment pipeline. The institution does not pull from the project's registry directly.

**Project recommendation.** Signature-preserving mirroring is the v1.0 recommended pattern for institution mirror registries. The institution's deployment pipeline validates the project's original cosign signature against the project's published cosign public key; trust path is project → mirror (byte-faithful) → deployment, with no additional trust anchor and no bridge invariant to test. Institutions adopting re-signing operate the documented bridging controls (WORM 7-year audit log per "Mirror audit-log retention and integrity"; continuous monitoring via the `mirror.reconciliation_completed` operational event; bridge-invariant documentation; IR Scenario 13 readiness). Re-signing without these controls is a non-conformant posture under v1.0; the institution's CC8.1 control description names which pattern is in force and the corresponding evidence the SOC team and FFIEC examiner consume. The recommendation is project-level and applies to every v1.0 deployment regardless of size, regulator, or vendor topology — institutions that already operate a mature internal cosign trust path are not exempt from the bridging-control burden if they choose re-signing; they document the bridge controls and the SOC team tests them. This recommendation is intended to be load-bearing at examination time: the examiner's first question is which pattern the institution operates; the examiner's second question is which controls back the chosen pattern.

**Mirror-registry signature handling.** Two patterns the institution operates, both conformant when their controls are in place:

1. **Signature-preserving mirror (preferred).** The mirror registry pulls the image bytes-faithfully and preserves the original cosign signature. The deployment pipeline pulls from the mirror and validates the original signature against the project's published cosign public key. The mirror is byte-faithful infrastructure; trust path: project → mirror (preserves bytes) → deployment (validates project signature). No additional trust anchor needed.
2. **Re-signing mirror (with documented bridging).** The mirror registry pulls the image, validates the project's signature, and **re-signs** the image under the institution's own internal cosign key. The deployment pipeline validates the institution's signature only. This is acceptable when the institution operates a mature internal cosign trust path, but the institution MUST document the trust-path bridge: "the mirror's re-signature is conditional on the project's cosign signature having validated at mirror-pull time, recorded in the mirror's audit log; the institution's IR playbook treats a missing mirror-pull validation log entry as a Scenario 3 finding."

Re-signing without documenting the bridge is non-conformant — the institution loses the project-side cosign as a trust anchor and depends entirely on its own internal cosign key. Examiners verify which pattern the institution operates and confirm the bridge documentation if the re-signing pattern is in use.

**Mirror audit-log retention and integrity (re-signing pattern only).** The bridge invariant — "the mirror's re-signature is conditional on the project's cosign signature having validated at mirror-pull time" — is asserted on the mirror's audit log. The audit log MUST satisfy two retention properties:

1. **Retention duration: 7 years minimum** (or longer if the institution's chain-event retention exceeds 7 years). The bridge invariant is examined long after the binary's deployment; an FFIEC examiner sampling the bridge on a binary deployed two years ago needs the original validation log entry available. Standard cloud-logging targets (CloudWatch Logs, Azure Monitor, GCP Cloud Logging) default to shorter retention; the institution explicitly extends.
2. **Integrity posture: WORM (write-once-read-many) storage.** The audit log MUST be on append-only / WORM storage, NOT standard log rotation that the operator can delete. Acceptable backends: AWS CloudWatch Logs with locked retention policy + S3 Object Lock; Azure Monitor with immutable storage; GCP Logging with locked retention bucket; on-prem WORM appliance. The bridge invariant rides on the audit log; if the log is rotation-deletable, the invariant cannot be examined retroactively.

**Mirror continuous-monitoring control (DE.CM-09).** A re-signing mirror that silently stops validating project signatures (operator misconfig at the mirror, registry-side bug, vendor regression) would produce mirror-side re-signed images with no project-side validation evidence — and the institution's deployment pipeline, validating only the institution's signature, would not detect this. The institution operates a periodic (daily or per-pull) reconciliation between the mirror's pulled-image set and the audit-log-validation set; mismatch fires an operational alert. This is the same control pattern as spec §10.1 fingerprint reconciliation, applied to the mirror's signature-validation discipline. The reconciliation event is the institution's `mirror.reconciliation_completed` operational event (institution-defined; not part of the chain's spec §10.2 catalog).

## SBOM consumption

The CycloneDX SBOM names every dependency in the binary. Institutions feed it to their vulnerability management:

- Snyk, Mend, JFrog Xray, etc. consume CycloneDX directly
- Comparison to the published vulnerability scan reports lets the institution confirm no new vulnerabilities have appeared since release

If a vulnerability is discovered post-release in a dependency, the project issues a patch release; institutions update on their standard cadence.

## Reproducible-build evidence

When the institution performs an independent rebuild, it emits a `reproducible-build-log` JSON record:

```json
{
  "rebuild_timestamp": "2026-05-15T14:30:00Z",
  "source_commit": "a3f29b71c8...",
  "go_version": "go1.22.3",
  "build_flags": "CGO_ENABLED=0 -ldflags='-s -w -buildid=' -trimpath",
  "binary_sha256": "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
  "matches_published": true,
  "verifier": "internal-build-pipeline-v2.1"
}
```

The log is archived in the institution's control-evidence repository. Sample-tested by SOC and examination teams.

## Operational pattern

For each new release:

1. Project publishes binary + signatures + SBOM + scan results + source tarball
2. Project announces the release with a 30-day adoption window
3. Institution updates its allowlist (SHA-256 hash or cosign signature)
4. Institution performs reproducible-build verification (at least one institution per release; ideally many)
5. Institution archives the rebuild log
6. Institution rolls out the new version per its standard change-management procedure
7. SOC and examination teams sample-test the institution's adoption procedure on the next examination cycle

## Project-side governance

The chain spec depends on the project (Apache 2.0 open-source community) operating its own controls. Institution adopters trust the project; the project earns the trust through public governance.

### Maintainer controls

The project operates the following controls:

- **Multi-maintainer review.** Code changes to the reference implementation require approval from 2-of-N maintainers; security-sensitive changes require approval from a security-cleared subset.
- **Public issue tracking.** All discussions, disagreements, and decisions occur in public issue tracking with full history.
- **Public commit history.** All commits are pushed to a public git repository; force-push is restricted; signed commits required.
- **External security review.** The project commissions external security reviews on a defined cadence; results published.
- **Conformance corpus integrity.** The corpus is part of the open-source repository; changes require multi-maintainer review; corpus content is reviewable independently of the implementation.

The project's governance is documented at `GOVERNANCE.md`. Institutions adopting the spec verify the project operates these controls as part of standard third-party-risk-management for open-source dependencies.

### Conformance corpus integrity

A subverted corpus that produces wrong-but-deterministic output is theoretically possible if an attacker compromises the project's source-of-truth. Defenses:

- Multi-maintainer review of corpus changes (procedural)
- Public git history of corpus modifications (auditable)
- Independent re-implementations from the spec text rather than the corpus (reduces single-source-of-truth)
- External security review of the corpus content periodically

Institutions adopting the spec for examination-grade use rebuild the corpus from the spec text rather than relying solely on the published corpus.

### Release-pipeline operator governance

The release pipeline is operated by the project's release-management role. The role has:

- Hardware-backed cosign key (held by a small group of release-managers)
- Offline GPG key (held by a separate small group, primarily for emergency-response scenarios)
- Documented release procedure (every release follows the same procedure)
- Audit logging of release operations

The project publishes its release-management governance in `GOVERNANCE.md`. The project's standard practice is to publish a SOC 2 report on the release pipeline's controls; institution adopters consume this as part of their third-party-risk-management.

## What examiners verify

In an examination:

- The institution has cached the project's cosign public key AND the project's GPG public key (for fallback) AND the project's SLSA-verifier public key
- The institution validates signatures before deploying new versions
- **The institution validates the SLSA provenance attestation against `slsa-verifier` at deployment time** (deployment-gating control); the SLSA-verifier output log is archived for the binary's deployment lifetime (typically 7 years)
- The institution's reproducible-build log shows verification of the deployed version
- The institution's SBOM consumption is integrated with its vulnerability management
- The container deployment uses a bank-controlled mirror registry with image-signature validation
- For institutions using the re-signing mirror pattern: the mirror's audit log shows the project-side cosign signature validated at pull time, AND the audit log is on WORM storage with 7-year retention, AND the mirror's continuous-monitoring control fires on signature-validation discipline gaps (the `mirror.reconciliation_completed` operational event)
- For institutions performing independent corpus rebuilds: the rebuild log shows the corpus reproduces from the spec text and matches `spec-manifest.sha256.asc`
- During an active cosign-key recovery window: the institution's `bundleverify --gpg-only` validations are recorded as the interim posture, with post-recovery re-validations against the new cosign key documenting evidence stability across the trust-anchor change

## SLSA level claim

The build pipeline meets the **SLSA Level 3** properties per the SLSA v1.0 specification:

| SLSA L3 requirement | How the pipeline satisfies it |
|---|---|
| Source — version controlled | Public git repository with signed commits and audit logging of force-push restriction |
| Source — verified history | All commits are signed; multi-maintainer review on the security-sensitive subset |
| Build — scripted build | Pipeline definition is in-repo, version-controlled |
| Build — build service | Hosted CI runs the pipeline (project-controlled hosted runner; not contributor laptops) |
| Build — ephemeral environment | Hermetic build environment per release; no persistent build infrastructure |
| Build — isolated | Build inputs are pinned (toolchain version, dependencies); no network access during build |
| Provenance — available | `verifier.intoto.jsonl` published per release per the in-toto attestation format |
| Provenance — authenticated | Provenance is signed by the build system's identity (Sigstore Fulcio cert or equivalent) |
| Provenance — service generated | Provenance is generated by the hosted build service, not by the developer |
| Provenance — non-falsifiable | Provenance signing happens inside the build service's trusted boundary |

**The SLSA L3 claim is auditable.** Institutions consuming the verifier may verify the published provenance attestation against their own SLSA-aware tooling (e.g., `slsa-verifier`). The cybersecurity examiner working a CSF assessment uses the SLSA L3 claim as evidence for GV.SC-04 (cybersecurity supply chain — third-party assessments) and ID.RA-09 (authenticity of hardware and software).

SLSA Level 4 (which adds two-party review and reproducibility verification at the build service) is a v1.1 candidate; the current pipeline meets the reproducibility property but does not yet operate the two-party review at the build-service layer.

**Institutional consumption of the SLSA L3 attestation (MUST-tier).** Publishing the SLSA provenance attestation is necessary but not sufficient — the institution MUST consume it. Specifically:

1. The institution MUST validate the published `verifier.intoto.jsonl` against `slsa-verifier` (or equivalent SLSA-aware tooling) BEFORE deploying a new release version of the verifier or ledger server. The validation is a deployment-gating control; a binary whose provenance attestation does not validate MUST NOT be deployed.
2. The institution's SLSA-verification log is a control-evidence artifact retained for the same period as the binary's deployment lifetime (typically the chain-event retention period, 7 years).
3. The institution's `slsa-verifier` binary itself MUST be signed and trust-anchored alongside the verifier — if the institution validates the chain's verifier with cosign but runs an unsigned `slsa-verifier`, the supply-chain trust path has a gap that defeats the L3 evidence claim.
4. `slsa-verifier` is a co-equal trust artifact alongside cosign and the GPG-signed manifest; the trust-anchor table at the top of "Trust anchors" should include `slsa-verifier` public key (and binary SHA-256) cached out-of-band.

**SLSA-aware tool-choice documentation (normative for examination-grade adopters).** The institution's control description MUST name which SLSA-aware tool the institution operates (`slsa-verifier` is the canonical reference; equivalent tools must be enumerated by name with the institution's evaluation rationale). The institution MUST document the tool's trust anchor (public key, binary SHA-256, signature-validation chain) in its trust-anchor table. The cybersecurity examiner asks: which tool, what version, what trust anchor, archived where; the institution's documentation has those four answers in one place. Without this documentation, the examiner cannot verify the L3 provenance-attestation consumption is operationally exercised against a known-good tool.

The cybersecurity examiner working a CSF assessment uses the institution's SLSA-verification log as evidence for GV.SC-04 and ID.RA-09. An institution that publishes the L3 claim without operating `slsa-verifier` consumption has a gap a GV.SC-04 examiner will challenge.

## Spec text and conformance corpus integrity

The spec PDF and the conformance corpus are published as signed artifacts (per the "What ships per release" table at the top of this doc). The institution's trust path for the spec and corpus mirrors the trust path for the binary:

- **Cosign-signed spec PDF.** The institution validates the cosign signature before adopting the spec; the spec's content-hash is in `spec-manifest.sha256.asc`.
- **GPG-signed spec content-hash manifest.** The fallback path; the institution validates the manifest against the cached GPG public key.
- **Cosign-signed conformance corpus tarball.** The institution validates the corpus's cosign signature; rebuilding the corpus from the spec text and confirming against `spec-manifest.sha256.asc` is the recommended deeper verification (per `09-threat-model.md` §5 independent re-implementation as defense).

Without these signed artifacts, an attacker who compromises the spec or corpus distribution channel could ship a spec/corpus that produces wrong-but-deterministic output and let it propagate through the institution's adoption pipeline. The signed-artifact trust path closes this with the same defense-in-depth as the binary trust path.

## Vendor-conformance attestation

The signed-artifact trust path establishes that the spec and corpus the institution consumes are authentic. A separate question is whether a **vendor's implementation** passes the corpus the institution consumes. The vendor's annual SOC report covers vendor-side operational controls — HSM operations, ledger operations, release-pipeline discipline. The SOC report does NOT necessarily attest that the vendor's implementation passes the FFIEC conformance corpus. An institution accepting a vendor's marketing claim of "FFIEC chain-of-custody conformant" without an independent attestation has a control gap.

The vendor-conformance attestation procedure closes that gap. The vendor produces an annual signed attestation against a specific corpus version; the project-side working group operates a public registry of attestations; the institution validates the attestation signature, confirms the corpus-version match against the corpus the institution itself runs against, and archives the attestation as vendor-management evidence.

The full procedure is documented in [`docs/vendor-conformance-attestation.md`](vendor-conformance-attestation.md). The high-level institution-side consumption flow:

1. The institution consults the registry at vendor selection and on the institution's ongoing-review cadence (typically annual, with ad-hoc re-validation on corpus-version change, vendor product-version change affecting chain behavior, or registry revocation).
2. The institution fetches the vendor's published attestation public key from the vendor's `.well-known` URL and caches it out-of-band.
3. The institution validates the attestation signature against the cached public key.
4. The institution confirms the corpus version named in the attestation matches the corpus version the institution itself runs against.
5. The institution archives the attestation document, signature, public key, and validation log per CUEC-VND-06.
6. The institution's CC8.1 control description names the attestation evidence used in vendor selection and the cadence of re-validation.
7. The institution's SOC team tests CUEC-VND-06 as part of the vendor-management control test.

A vendor that refuses to attest is treated as a control gap; the institution either requires attestation as a contractual condition, operates the compensating control of running the corpus itself against the vendor's product, or replaces the vendor. A vendor whose attestation is revoked by the working group triggers institution-side vendor-management re-evaluation per the institution's standard procedure.

The procedure is separate from, and complementary to, vendor SOC reporting. Institutions consume both: the SOC report covers operational controls; the vendor-conformance attestation covers conformance against the published corpus. Together they satisfy the institution's CC8.1 vendor-management evidence base for chain-of-custody implementations operated or shipped by vendors.

---

## SLSA, Sigstore, in-toto, SBOM, and SSDF composition

The supply-chain primitives this document established cover cosign + GPG dual signing, hermetic and reproducible builds, SLSA Level 3 provenance per SLSA v1.0, CycloneDX SBOM publication, signed-binary attestation per spec §10.7, the conformance corpus as a versioned signed artifact, the cold-DR key with annual dry-run attestation, and explicit revocation procedures. This section composes those primitives across the SLSA, Sigstore, in-toto, SBOM, and SSDF frameworks so the institution's deployment-gating control runs end-to-end across SDK, ledger, and verifier components and so the chain's posture maps cleanly to OpenSSF Scorecard, GitHub artifact attestations, and OMB M-22-18 federal-procurement requirements.

### Per-component SLSA Build Level alignment

The chain ecosystem has three load-bearing software components: the SDK (links into the institution's application and produces chain entries), the ledger server (ingests entries, produces seal records), and the verifier binary (produces the PASS/FAIL evidence the examiner relies on). SLSA v1.0 attests per artifact; each component carries its own provenance.

| Component | Minimum SLSA Build Level | Per-release attestation | Build platform identity | Source-control authentication |
|---|---|---|---|---|
| SDK (per language) | Level 3 | `sdk-{language}.intoto.jsonl` | GitHub-hosted runner with hardened isolation | Signed commits + branch protection |
| Ledger server | Level 3 | `ledger.intoto.jsonl` | GitHub-hosted runner with hardened isolation | Signed commits + branch protection |
| Verifier binary | Level 3 | `verifier.intoto.jsonl` | GitHub-hosted runner with hardened isolation | Signed commits + branch protection |

The provenance-signing identity is the project's release-management identity (Fulcio short-lived certificate's GitHub OIDC claim binding the workflow to the project's release branch). The institution's deployment-gating control extends from validating one attestation to validating every per-component attestation for the components the institution deploys. Institutions deploying only the verifier (e.g., examiner-side verification of a vendor-hosted chain) validate only the verifier attestation; institutions deploying the full stack validate all three.

The institution's CC8.1 control description names the SLSA-verifier procedure for each component and the per-release attestation file consumed.

### Sigstore Rekor transparency-log inclusion

Cosign signing without Rekor inclusion provides authentication (the signature verifies against the project's pinned cosign public key) but not transparency (an external observer cannot independently confirm that a signature was issued). For supply-chain scenarios where a vendor-side compromise could include silent re-signing under the legitimate cosign key, Rekor's transparency log is the canonical detection mechanism: every legitimate signature appears in the public Rekor log; an unexpected signature absent from Rekor is anomalous.

Every cosign signature produced by the project's release pipeline is appended to the Sigstore Rekor transparency log at signing time. The institution's deployment-gating control validates Rekor inclusion alongside cosign signature validation. The validator (`cosign verify --rekor-url ...` or equivalent tooling) confirms:

- the signature's Rekor log entry is present,
- the entry's signing identity matches the project's release-management identity (Fulcio short-lived certificate's GitHub OIDC claim binding the workflow to the project's release branch),
- the entry's inclusion proof verifies against Rekor's signed tree head.

A signature that validates against the project's pinned cosign key but is absent from Rekor is non-conformant. The institution's IR procedure treats Rekor-inclusion-missing as a Scenario 3 escalation (signature verification failed) with branch (d) `Rekor inclusion missing` distinguishing it from primary signature failure.

### in-toto Layout Manifest for the build pipeline

SLSA L3 attestation provides one provenance statement per output artifact. The in-toto Layout Manifest specification 2024 attests every step in a multi-step pipeline (source checkout, dependency fetch, compile, test, sign, release-publish) with per-step attestation binding the declared functionary, inputs, and outputs. SLSA names provenance availability; in-toto names step-shape-conforming-pipeline. The two compose: SLSA gives "this binary came from this platform"; in-toto gives "this platform's pipeline followed the declared steps with the declared functionary identities."

Without the Layout, an attacker who compromises the release-management role's pipeline definition (alters which steps run or which functionary signs) can produce an SLSA L3 provenance attestation that is technically valid but pipeline-divergent. The Layout closes this gap by binding the pipeline structure as an attestation.

The project publishes an in-toto Layout Manifest per the in-toto Layout Specification 2024 declaring:

- the build-pipeline steps (source-checkout, dependency-fetch, compile, test, sign, release-publish),
- each step's expected functionary identity (Fulcio short-lived cert, GPG fingerprint, or hardware-token serial),
- each step's expected input and output material,
- the chain of step attestations binding step N's outputs to step N+1's inputs.

The Layout Manifest is signed under the project's cosign key with Rekor inclusion. Institutions operating in-toto-aware tooling (the Tekton chains controller, the GitHub artifact-attestations validator, the Sigstore policy-controller) opt into pipeline-shape validation in addition to per-artifact provenance validation, providing defense-in-depth against pipeline-structure subversion.

v1.0 ships the Layout Manifest informatively. v1.1 makes pipeline-shape validation a deployment-gating control alongside SLSA-verifier.

### in-toto Layout Manifest for the conformance corpus

The conformance corpus at `spec/test-vectors/` is a load-bearing artifact: every vendor's implementation is conformance-tested against it, and corpus subversion that produces wrong-but-deterministic output across implementations is the textbook supply-chain attack against the chain ecosystem. The corpus signature today attests "this tarball was signed by the project's cosign key." An in-toto Layout for the corpus attests "this tarball is the output of the declared spec-authoring → test-generation → working-group-approval pipeline, with each step's functionary identity captured."

The spec working group publishes an in-toto Layout Manifest for the conformance corpus declaring:

- spec authoring (with the spec-text commit hash as the input),
- test-vector generation (from the spec text with the generator binary's SHA-256 as the functionary tooling),
- working-group review (with the reviewer functionaries' signatures),
- working-group approval (with the WG-chair signature),
- release publication (with the release-manager signature).

The Layout Manifest is signed under the project's cosign key with Rekor inclusion. The institution's corpus-integrity validation extends from "validate the cosign signature on the tarball" to "validate the in-toto Layout Manifest's pipeline-shape conformance." This is the highest-leverage in-toto application in the spec because the corpus is the conformance-bar floor across the ecosystem. A SolarWinds-class adversary compromising the release-publication step but unable to retroactively forge the working-group-approval functionary signatures is detected by Layout validation.

### SBOM dual-format publication: CycloneDX 1.6 and SPDX 3.0

The chain SBOM is published in CycloneDX 1.6 and SPDX 3.0 formats in parallel. CycloneDX 1.6 (late 2024) adds explicit support for ML model components via `component.type=machine-learning-model`, formula-based composition for derived components, signed SBOMs per JSON Schema Forms, and OpenSSF Vulnerability Disclosure Report integration. SPDX 3.0 adds AI/ML profile coverage through SPDX-AI 3.0 and improved transitive-dependency modeling.

NTIA SBOM minimum elements accept either format. Federal procurement under OMB M-22-18 tracks toward both formats in parallel for downstream-consumer flexibility, and CISA's SBOM-a-rama findings recommend dual-format publication for federal-agency vendors. The dual-format commitment removes the institution's tooling-choice gating dependency on the project: institutions feeding SBOMs to vulnerability management (Snyk, Mend, JFrog Xray, Anchore) consume the format their tooling supports.

#### AI/ML model-component SBOM coverage

The SBOM uses CycloneDX 1.6's `component.type=machine-learning-model` for any model component the SDK / ledger / verifier ships with. The SPDX SBOM uses the SPDX-AI 3.0 profile in parallel. Examples of model components: a tokenizer model, a classification model used in the verifier's anomaly-detection path, model-fingerprinting weights for vendor-implementation testing.

### NIST SP 800-218 SSDF mapping

OMB Memorandum M-22-18 (September 2022) and follow-up M-23-16 (June 2023) require federal agencies' software providers to self-attest compliance with NIST SP 800-218 SSDF v1.1. CISA's Secure Software Development Attestation Form (June 2023) is the canonical document. National banks regulated by OCC, federal-chartered credit unions regulated by NCUA, and federal-agency-adjacent banking infrastructure (FedNow participants, FRB-supervised entities) are inheriting SSDF self-attestation as a supplier requirement.

The chain's published SSDF mapping (informative; banks supplement with their own institutional controls):

| SSDF v1.1 practice | Chain control satisfying it |
|---|---|
| PO.1.1 (Define Roles) | `GOVERNANCE.md` working-group structure |
| PO.5.1 (Implement Secure Environments) | hermetic build environment per §"Build pipeline" |
| PO.5.2 (Use Secure Environments) | hermetic build environment + GitHub-hosted hardened runners |
| PS.1.1 (Protect Code from Unauthorized Access) | multi-maintainer review + signed commits + branch protection |
| PS.2.1 (Verify Software Release Integrity) | cosign + GPG dual-signing + Rekor inclusion |
| PS.3.1 (Archive and Protect Each Software Release) | Sigstore Rekor transparency log + signed-artifact archive |
| PW.1.1 (Design Software with Security Requirements) | spec body §§1-10 normative requirements |
| PW.4.1 (Reuse Well-Secured Software) | published CycloneDX 1.6 + SPDX 3.0 SBOM with vulnerability scan |
| PW.4.4 (Verify Reuse Components) | per-release dependency vulnerability scan published |
| PW.7.1 (Configure Compilation) | reproducible-build configuration documented |
| PW.7.2 (Review Code) | mandatory review on protected branches |
| PW.8.1 (Test Executable) | per-release test-suite with conformance-corpus run |
| PW.9.1 (Configure Software for Default Secure Settings) | spec defaults + §10.7 software-key adapter exclusion |
| RV.1.1 (Identify Vulnerabilities) | published scan reports + the §4.3.2 emergency-patch SLA |
| RV.1.3 (Have a Vulnerability Disclosure Program) | `SECURITY.md` + OSV.dev / GHSA publication (see §"Vulnerability response") |
| RV.2.1 (Analyze Vulnerabilities) | working-group response with 30-day cryptographic-attack SLA |
| RV.2.2 (Develop and Implement Mitigations) | spec-patch + new release per §4.3.2 |
| RV.3.1 (Analyze Vulnerabilities for Root Cause) | post-incident review per the IR playbook |

The mapping is published alongside each release as a derivative document. Institutions submitting downstream SSDF self-attestation to federal agencies cite the project's published mapping. Federal-agency software providers consuming the chain inherit the chain's SSDF posture for the chain components and supplement with their own institutional controls for surrounding integration.

### Reproducible-build extension: independent-rebuilder attestation

The spec's existing reproducible-build commitment covers byte-for-byte build reproducibility under the documented Go toolchain flags. The verifier-as-evidence-tool case has a stronger requirement: the FFIEC examiner using the verifier as forensic evidence in a matter where the institution and an adverse party (CFPB-disputing customer, plaintiff class) need the binary's provenance to flow back to source independent of the project's signing key. The Reproducible Builds project's standard goes beyond byte-identical output: it expects a published `buildinfo` file recording every build input and a published rebuild attestation from at least one independent rebuilder.

The project commits to:

- publishing a per-release `buildinfo` file in the Reproducible Builds project's standard format, listing every input the build consumed (Go toolchain version, dependency module versions, build flags, source-tree commit hash, hermetic-environment image hash);
- operating an independent-rebuilder attestation pipeline where at least one non-project-controlled rebuilder (a separate cloud account, a partner foundation, a CMU-SEI test rebuilder, an FFIEC working-group-coordinated rebuilder) rebuilds from source on each release and publishes a signed rebuild-attestation;
- publishing the rebuild-attestations on the project's release page alongside `verifier.intoto.jsonl`.

The institution's reproducible-build evidence (the existing `reproducible-build-log` JSON) becomes one input among at least three: project-built binary, project's SLSA L3 attestation, independent-rebuilder attestation. Each is independently validatable. The trust-floor for the verifier-as-evidence-tool case rises from "trust the project's published binary" to "trust the binary if at least three independent parties rebuild from source and produce identical bytes."

### AI/ML model-supply-chain composition

The chain captures `gen_ai.request.model` and `gen_ai.response.model` for every chain entry representing a model call. The chain does NOT directly capture the model's supply-chain provenance — training-data manifest, training-code commit, model-weights hash, model-card URL, model-fingerprint, or any cryptographic attestation by the provider that the model was produced by the declared training pipeline. Hugging Face's model-attestation work (Sigstore-signed model artifacts since late 2024), the OpenSSF Model Signing project, and the proposed in-toto AI Layout extension address this gap on the model-producer side. The chain composes naturally with model-provenance attestation when a provider ships one.

The v1.x candidate optional attribute `gen_ai.model_provenance` (string or bytes) captures the provider's model-supply-chain attestation when one is shipped. Schema: a JCS-canonical JSON envelope containing the model identifier, the producer identity (provider's signing key fingerprint or model-hub URL), the training-pipeline identifier (commit hash of training code, training-data manifest hash, hyperparameter set hash), and the provenance signature. The chain captures the attestation byte-for-byte under the OTel envelope per spec §5; the verifier ignores the attribute for chain-integrity decisions (parallel to the existing `gen_ai.provider_attestation` posture).

The three-layer evidence model:

1. **Chain integrity.** The chain proves recording integrity of the AI invocation.
2. **Model identity.** The chain attribute names the model the institution requested and the model the provider delivered.
3. **Model provenance.** The captured attestation, validated separately, proves the named model was produced by the declared training pipeline.

This positions the chain for emerging AI-supply-chain frameworks (NIST AI RMF Generative AI Profile, EU AI Act Article 12 logging requirement composing with Article 13 transparency requirement on model provenance, the OpenSSF Model Signing v1 work) without requiring spec changes that depend on those frameworks finalizing.

### Verifier output supply-chain manifest

The verifier produces a PASS/FAIL output that travels to downstream consumers (examiner working papers, plaintiff expert deposition exhibits, regulator enforcement file). The supply-chain doc covers the verifier binary's provenance — institution validates SLSA L3 + cosign + Rekor + reproducible-build before deploying. The verifier's output is itself a downstream supply-chain artifact and benefits from accompanying supply-chain context.

The verifier emits a `verifier-output-manifest.json` alongside every PASS/FAIL run capturing:

- the verifier binary's SHA-256,
- the verifier's release identifier,
- a content-addressed reference to the verifier's `verifier.intoto.jsonl` SLSA attestation (via Rekor URL or cosign signature handle),
- the chain file's SHA-256,
- the seal record's SHA-256,
- the run timestamp,
- the PASS/FAIL disposition with reason string per spec §7.

The manifest is plain JSON. Cryptographic signing of the manifest is a v1.1 candidate (the chain's existing capture-and-archival posture for verifier output covers v1.0). The downstream consumer reading the manifest re-fetches the verifier's SLSA attestation, validates the verifier binary's supply-chain posture, and confirms the manifest's binary SHA-256 matches the attestation's output hash. This closes the verifier-output-as-supply-chain-input gap without expanding the chain's normative cryptographic surface.

### Subservice supply-chain composition

The ledger storage layer for many deployments rides on managed cloud storage (AWS S3, Azure Blob Storage, Google Cloud Storage). These services have their own supply-chain attestation: the cloud provider's SLSA / SOC 2 / FedRAMP / ISO 27001 reports cover the storage subservice. The institution's defense-in-depth posture under spec §10.5 (HSM custody) and §10.13 (evidentiary artifacts) implicitly inherits the cloud provider's posture.

The institution's CC8.1 control description names which storage subservice is in use and which supply-chain attestation (the cloud provider's SOC 2 / FedRAMP / ISO 27001 / SLSA where available) the institution consumes annually for the subservice. The SOC-pack subservice-organization scoping (carve-out method per the existing SOC-pack guidance) names the subservice's supply-chain posture as a complementary user-entity control.

### Vendor-conformance attestation as in-toto step attestation

Spec §10.12 (verifier CLI exit-code contract) means multiple vendors ship verifiers conforming to the same exit-code semantics. The vendor-conformance attestation procedure today (`docs/vendor-conformance-attestation.md`) is shaped as: vendor signs an attestation document; institution validates the signature; institution archives. An in-toto step attestation gives this procedure a structured machine-readable shape.

Each "vendor-X verifier passes corpus version Y" attestation is an in-toto step attestation per the in-toto Specification 2024:

- declared functionary identity: vendor's signing key (cosign or HSM-backed),
- declared materials: corpus version Y's SHA-256 (matching `test-vectors-v{version}.tar.gz` published by the project),
- declared products: vendor verifier binary SHA-256 + test-execution log SHA-256 + PASS/FAIL summary.

The attestation is signed under the vendor's cosign key with Rekor inclusion. The institution's validation is mechanical: validate the in-toto step attestation signature, confirm Rekor inclusion, confirm corpus-version match against the corpus the institution itself runs against, confirm the test-execution log's PASS summary, archive all four artifacts. The project-side working-group-operated registry publishes attestations as in-toto step attestations rather than free-form PDFs. This raises conformance-attestation rigor across vendors and makes the registry mechanically queryable.

### Vulnerability response and coordinated disclosure

Spec §4.3.2 commits to a 30-day spec-patch publication SLA on credible cryptographic-attack demonstration. This composes with industry-standard advisory channels. The project commits to:

- the OpenSSF Vulnerability Disclosure Working Group's coordinated disclosure timeline (90 days standard; 30 days for actively-exploited; 7 days emergency for cryptographic-substrate breaks per §4.3.2);
- publishing every CVE under both OSV.dev (for SBOM-tool consumption) and GitHub Security Advisories (for repository-level visibility);
- coordinating CNA assignment through GitHub's CNA or another working-group-named CNA;
- maintaining `SECURITY.md` at the repository root naming the disclosure procedure, the project's PGP-encrypted disclosure email, and the response SLA.

The institution's vulnerability-management tooling (Snyk, Mend, JFrog Xray) ingests OSV.dev / GHSA feeds and matches against the SBOM. A chain SDK CVE shipped only as a project-side announcement without OSV.dev / GHSA publication is invisible to the institution's standard vulnerability-management pipeline. The dual-channel publication closes that gap. The 30-day cryptographic-attack SLA composes with OSV.dev / GHSA publication windows (typically 7 days from credible report) so the institution's tooling sees the vulnerability through standard channels.

The §"SBOM consumption" section above ties to this section: institutions consuming the SBOM for vulnerability management rely on the OSV.dev / GHSA publication for vulnerability identification. The two sections compose: SBOM for inventory; OSV.dev / GHSA for vulnerability stream.

---

## Per-role roll-up for the supply-chain composition

**For the spec working group.** The composition layer extends the existing primitives (SLSA L3, cosign, GPG, hermetic builds, reproducible builds, CycloneDX SBOM, signed-binary attestation, vendor-conformance attestation, cold-DR key) across SLSA + Sigstore + in-toto + SBOM dual-format + SSDF + OSV. None of the additions change the spec's normative cryptographic surface. All position the chain for OpenSSF Scorecard upper-quartile scoring, federal-procurement bidder readiness under OMB M-22-18, and natural composition with the AI/ML model-supply-chain dimension.

**For institutions targeting OpenSSF Scorecard and federal-procurement readiness.** The chain's posture closes the per-component SLSA gap, the Rekor-inclusion gap, the in-toto Layout gap, and the SSDF mapping gap that affect Scorecard scoring. Institutions inheriting the chain's posture for the SDK, ledger, and verifier extend their deployment-gating control to validate the institution's own deployment artifacts under the same framework.

**For vendors shipping chain implementations.** The in-toto step attestation framing for vendor-conformance attestation (per `docs/vendor-conformance-attestation.md`) makes vendor-conformance evidence consumable through standard supply-chain tooling rather than vendor-specific attestation review. Vendors shipping under SLSA L3 with Rekor inclusion plus in-toto step attestation against the published corpus produce machine-readable conformance evidence at one CI-integration cost.

**For downstream consumers.** The verifier-output supply-chain manifest and the corpus in-toto Layout give consumers the supply-chain context to validate the producing-verifier's posture and the conformance bar the chain artifacts were validated against, without relying on institution-discretionary archival. A consumer reading verifier output years later reconstructs the supply-chain trust path back to the spec working group's published Layout Manifest — defense-in-depth surviving institution turnover, vendor exit, and project leadership change.

The project-side working-group governance commitments — sub-committee composition, registry review cadence (quarterly), revocation publication timeline (14 days), corpus-version re-attestation grace period (90 days) — are in `GOVERNANCE.md` §"Vendor-conformance attestation registry."
