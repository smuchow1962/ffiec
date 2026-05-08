# Vendor-conformance attestation

> **What this doc is.** The procedure by which a vendor attests that its chain-of-custody implementation passes the FFIEC conformance corpus, the project-side registry of authorized-implementer attestations, and the institution-side consumption procedure that turns the attestation into examination-grade evidence. This procedure is separate from, and complementary to, the vendor's annual SOC report.

## Purpose

A vendor hosting or shipping a chain-of-custody implementation consumes the published FFIEC conformance corpus (`spec/test-vectors/`) to validate its implementation. The vendor's annual SOC report attests to vendor-side **operational** controls — HSM operations, ledger operations, release-pipeline discipline, change-management. The SOC report does NOT attest that the vendor's implementation **passes the FFIEC conformance corpus**. That is a separate question with a separate evidence path.

An institution that accepts a vendor's word on "we are FFIEC-chain-conformant" without an independent attestation has a control gap. The institution's CC8.1 control description names the attestation evidence used in vendor selection; without the attestation, the institution depends on the vendor's marketing claim with no project-side verification path.

The vendor-conformance attestation closes that gap. The vendor produces a signed attestation against a specific corpus version. The project-side working group operates a public registry of attestations. The institution validates the attestation signature, confirms the corpus-version match against the corpus the institution itself runs against, and archives the attestation as part of vendor-management evidence.

This procedure is the project-side mechanism for establishing trust in vendor-hosted and vendor-shipped chain-of-custody implementations. The procedure is normative for vendors claiming FFIEC chain-of-custody conformance and for institutions consuming such claims.

## Procedure shape

### Cadence

The vendor produces an attestation **annually**, aligned with the vendor's product release cycle. A new attestation is also produced when:

- The vendor ships a new product version that changes chain-of-custody behavior (a SHOULD — vendors that ship multiple product versions per year produce one attestation per such release; vendors that ship one product version per year produce one annual attestation).
- The FFIEC conformance corpus publishes a new version. The vendor MUST produce a re-attestation against the new corpus within **90 days** of the corpus version's publication date. An attestation against a superseded corpus version remains valid for the institution's archived-evidence purposes but is NOT current for new vendor-selection or ongoing-review use.
- An exception report on the previous attestation's exception list is closed (a previously-skipped test now passes). The vendor SHOULD produce a re-attestation reflecting the closed exception within the next attestation cycle.

### Vendor responsibilities

The vendor:

1. Runs the FFIEC conformance corpus against the vendor's product version under attestation.
2. Records the test-execution date, environment, and complete results.
3. Computes a deterministic test-results hash (SHA-256 over the canonical-form-encoded results table — see "Test-results hash construction" below).
4. Produces the attestation document per the schema in "Attestation document content."
5. Signs the attestation document with the vendor's published public-key signing key (see "Signature mechanism").
6. Submits the attestation to the project-side working group for registry inclusion (see "Working-group governance").
7. Retains the test-execution artifacts (raw corpus output, environment manifest, results-table source) for the longer of (a) the institution's chain-event retention period applicable to the vendor's customer base, or (b) seven years from the attestation date.

### Project-side working-group responsibilities

The working-group sub-committee (see "Working-group governance"):

1. Reviews the submitted attestation for schema conformance, signature validity, and corpus-version alignment.
2. Cross-checks the test-results hash against the working group's own canonical run output for the corpus version when feasible (the working group does not re-run the vendor's product, but it confirms the test-results hash is well-formed and the corpus-version reference is current).
3. Publishes the attestation to the public registry within the published submission-review timeline.
4. Maintains the registry's currency — superseded attestations move to the registry's "historical" view; revoked attestations move to the registry's "revoked" view with a public revocation log entry.

### Institution-side responsibilities

The institution:

1. Consults the registry at vendor selection and on its ongoing-review cadence.
2. Validates the attestation signature against the vendor's published public key.
3. Confirms the corpus-version match against the corpus version the institution itself runs against (or commits to run against in the deployment timeframe).
4. Archives the attestation document, the validation log, and the registry-consultation record in the institution's control-evidence repository per CUEC-VND-06.
5. Re-validates annually as part of the institution's vendor-management cycle.

## Attestation document content (normative)

The attestation document MUST contain the following fields. The document is JSON, UTF-8 encoded, JCS-canonicalized (RFC 8785) before signing.

```json
{
  "attestation_version": "1.0",
  "attestation_id": "<UUID v4 unique to this attestation>",
  "issued_at_utc": "<RFC 3339 timestamp>",
  "validity_end_utc": "<RFC 3339 timestamp; typically issued_at + 1 year>",
  "vendor": {
    "name": "<vendor legal name>",
    "url": "<vendor product URL>",
    "contact_email": "<conformance-program contact>"
  },
  "product": {
    "name": "<product name>",
    "version": "<product version string>",
    "git_tag_or_commit": "<source-of-truth identifier for the product build>"
  },
  "corpus": {
    "version": "<corpus version, e.g., v1.0>",
    "git_tag": "<spec/test-vectors/v1.0 git-tag>",
    "content_hash_sha256": "<SHA-256 of the canonical-corpus tarball spec-manifest.sha256.asc references>"
  },
  "test_execution": {
    "executed_at_utc": "<RFC 3339 timestamp>",
    "environment_summary": "<OS, runtime version, HSM model if applicable>",
    "results_hash_sha256": "<deterministic hash of the canonical results table; see Test-results hash construction>",
    "results_summary": {
      "total_cases": "<integer>",
      "passes": "<integer>",
      "fails": "<integer>",
      "skips": "<integer>"
    },
    "exceptions": [
      {
        "case_id": "<corpus case identifier, e.g., 010-tenant-ikm-rotation-mid-day>",
        "disposition": "skip | fail",
        "rationale": "<why the case was skipped or failed; what compensating control covers the gap if any>"
      }
    ]
  },
  "signatory": {
    "name": "<full legal name>",
    "role": "<CTO, VP Engineering, or equivalent senior engineering role>",
    "title_at_vendor": "<as reported on vendor's corporate registry>",
    "commitment_statement": "I attest that the vendor product named above, at the version named above, was tested against the FFIEC conformance corpus version named above, with the results recorded above. The test execution was performed honestly, the results are accurate, and any exceptions are disclosed. I commit, on behalf of the vendor, to publish a re-attestation when material changes affect any of the named scope items, and to notify the project-side working group of any subsequent finding that would invalidate this attestation."
  }
}
```

### Test-results hash construction

The vendor's test-execution output is reduced to a canonical results table — one row per corpus case, containing case ID, disposition (pass/fail/skip), and any case-specific output digest (e.g., for cases that compare a verifier output to an expected output, the SHA-256 of the canonical-form verifier output). The table is sorted lexicographically by case ID. The canonical form is JCS-encoded JSON. The SHA-256 of the canonical-form table is the test-results hash.

The hash is deterministic — two independent vendors running the same product version against the same corpus version produce the same hash. A hash mismatch between two vendors claiming attestation against the same corpus version is a working-group escalation point (one or both vendors are running modified corpus content or modified product behavior).

### Exception list

The exception list is the load-bearing transparency mechanism. A vendor that runs the corpus and skips three cases must list those cases with rationale. An empty exception list is the strongest attestation posture; a non-empty list is acceptable when the rationale is transparent and the institution can evaluate whether the skipped cases are material to the institution's deployment.

Examples of acceptable exception rationale:

- "Case `015-dual-algorithm-cosigned-seal` skipped because the vendor product does not support post-quantum algorithms in this version. The vendor commits to coverage in product version v1.3 (Q3 2026)."
- "Case `010-tenant-ikm-rotation-mid-day` failed in environment-specific manner traced to the vendor's HSM-vendor's library. The vendor disclosed the issue to the HSM vendor; remediation expected in HSM library version X.Y."

Unacceptable rationale (working group rejects the attestation):

- "Case skipped, no rationale."
- "Case is irrelevant to our customer base." (Working group decides corpus relevance, not the vendor.)

## Signature mechanism (normative)

### Signing key requirements

The vendor's signatory signs the attestation document with the vendor's **published public-key signing key**. Acceptable signing schemes:

- **Ed25519** (RFC 8032) — the project's primary recommendation for new vendor attestations.
- **RSA-PSS-SHA256** (RFC 8017 with SHA-256 as both hash and MGF) — accepted for institutions whose existing trust infrastructure is RSA-only and for vendors aligned with that posture. Vendors using RSA MUST use a key length of at least 3072 bits.

The signature covers the JCS-canonicalized attestation document bytes. The signature is appended to the document as a detached signature artifact (e.g., `attestation-v1.0-product-v1.2.json.sig`).

### Public-key publication

The vendor publishes the public key on a documented vendor-controlled URL (e.g., `https://vendor.example.com/.well-known/ffiec-attestation-key`). The URL MUST:

- Be served over TLS 1.3 (TLS 1.2 acceptable until the project's TLS 1.2 sunset date — currently 2028-01-01 per spec §5.1).
- Return the public key in a documented format (e.g., PEM-encoded for the public key, with a JSON metadata file naming the algorithm, the key fingerprint, the issuance date, and the key-rotation-notification URL).
- Be reachable from public networks without authentication (institutions and the working group MUST be able to fetch the key without vendor coordination).

### Key rotation

The vendor commits to a **30-day rotation-notification cadence**. When the vendor rotates the attestation signing key, the vendor publishes the new key at the same URL, retains the old key alongside the new key for a transition window of at least 30 days, and notifies the project-side working group within 14 days of the rotation. Institutions that consume attestations under the old key during the transition window remain on a valid trust path; institutions that begin consumption under the new key validate against the new key.

A key rotation does NOT invalidate prior attestations signed under the old key — the old key's archived public key remains a valid validation anchor for the prior attestations' validity windows. The vendor's archived-key documentation MUST preserve the old keys with their issuance and rotation dates so historical attestations remain validatable.

### Institution-side validation

The institution's validation procedure:

1. Fetch the vendor's published public key from the documented URL. Cache it out-of-band.
2. Fetch the attestation document and signature artifact from the vendor's publication URL or from the registry.
3. Validate the signature against the cached public key.
4. Confirm the signatory identity in the attestation document matches a published vendor-side commitment (the vendor's signatory must be a named individual on the vendor's website or corporate registry; institutions cross-check).
5. Confirm the corpus version named in the attestation matches the corpus version the institution itself runs against (or commits to run against in its deployment timeframe).
6. Archive the attestation document, signature, public key, and validation log per CUEC-VND-06.

A failed signature validation triggers a vendor-management escalation. A corpus-version mismatch triggers an institution-side review (the institution either upgrades its corpus consumption to match the attestation, downgrades the attestation evidence to a "for context only" disposition, or requests a re-attestation against the institution's current corpus version).

## Registry shape (normative)

The project-side working group operates a public registry of authorized-implementer attestations. The registry is reachable at a stable URL pattern under the project's domain.

### URL pattern

```
https://ffiec-chain-of-custody.org/registry/                         (registry index)
https://ffiec-chain-of-custody.org/registry/{vendor-slug}/           (vendor-specific page)
https://ffiec-chain-of-custody.org/registry/{vendor-slug}/{product-slug}/{version}/  (per-attestation entry)
https://ffiec-chain-of-custody.org/registry/revoked/                 (revocation log)
https://ffiec-chain-of-custody.org/registry/historical/              (superseded attestations, browsable)
```

The actual web infrastructure is operational (the working group operates the website per the project's standard public-website discipline). The URL pattern and the content shape per URL are normated here so cross-vendor and cross-institution consumption is consistent.

### Per-entry content shape (normative)

Each registry entry contains:

| Field | Description |
|---|---|
| Vendor name | Vendor legal name |
| Vendor URL | Vendor product URL |
| Product name and version | The attested product and version |
| Corpus version | The corpus version the attestation is against |
| Corpus content-hash | The SHA-256 the attestation references |
| Attestation document URL | The vendor-controlled URL where the institution downloads the signed attestation |
| Attestation public-key URL | The vendor-controlled URL where the institution fetches the validation public key |
| Attestation signature | The signature value (or a content-hash of the signature artifact for compact registry display) |
| Registration date | When the working group accepted the attestation |
| Validity end date | When the attestation expires (typically registration date + 1 year) |
| Status | Active / Superseded / Revoked |
| Notes | Working-group annotations (e.g., "Re-attestation pending after corpus v1.1 publication; see vendor's published timeline") |

### Public access

The registry is public. Institutions consult it without authentication. The registry's content is signed by the working group's release-management role using the project's standard signing keys (cosign primary, GPG fallback) — institutions consuming registry data may validate the registry-document signatures against the project's published trust anchors.

The registry also publishes a **machine-readable feed** (JSON Lines or equivalent) for institutions integrating registry consumption into their vendor-management automation. The feed's schema versioning follows the project's standard schema-versioning convention.

## Working-group governance (normative)

### Sub-committee composition

A **vendor-conformance sub-committee** within the project-side working group operates the registry and the attestation review process. The sub-committee:

- Comprises at minimum three project maintainers (separated from the spec-editor role to avoid concentration of decision-making).
- Operates under the project's standard public-decision-making discipline (decisions logged in public issue tracking; rationale published).
- Rotates members on a documented cadence (per the project's standard maintainer-rotation discipline).

### Vendor submission process

A vendor submits an attestation to the registry via the project's standard submission channel (currently a GitHub issue or pull request against the registry's public repository). The submission package contains:

1. The signed attestation document (JSON + signature artifact).
2. A pointer to the vendor's published public-key URL.
3. An evidence package: the vendor's environment manifest, a description of the test-execution procedure, and any per-case output artifacts the vendor includes for transparency.
4. A working-group reviewer assignment request.

The sub-committee acknowledges the submission within **5 business days** and completes review within **30 calendar days** of acknowledgement. Outcomes:

- **Accepted.** The attestation is published to the registry; the vendor is notified.
- **Returned for revision.** The sub-committee names specific items requiring vendor action (schema gap, signature validation issue, exception-list rationale insufficient). The vendor revises and resubmits; the 30-day clock restarts on resubmission.
- **Rejected.** The sub-committee documents the rejection rationale publicly. The vendor may appeal per the project's standard appeals process.

### Revocation procedure

The sub-committee revokes an attestation when:

- The corpus version is superseded and the vendor fails to re-attest within the 90-day grace period.
- A subsequent finding establishes the attestation misrepresents the test results (vendor-side disclosure or external discovery).
- The vendor fails a subsequent corpus version (publishes a new attestation that changes prior pass dispositions to fail without intervening product changes that explain the regression).
- The vendor fails to maintain the attestation public-key URL (e.g., the URL becomes unreachable for more than the published-rotation transition window).

Revocation is published to the registry's revocation log within **14 calendar days** of the sub-committee's revocation decision. The revocation log entry contains the attestation ID, the vendor name, the product and version, the revocation rationale, and the revocation date.

Institutions monitor the revocation log as part of vendor-management ongoing-review. A revoked attestation triggers institution-side vendor-management review per the institution's standard procedure.

### Cadence of registry review

The sub-committee reviews the registry **quarterly**. The review confirms:

- Each registry entry's vendor-published-key URL is reachable.
- Each registry entry's attestation document URL is reachable.
- Pending re-attestations (e.g., 90-day-clock for corpus version updates) are tracked.
- The revocation log is current.

The quarterly review's output is published as a registry-review summary signed by the working group's release-management role.

### Corpus-version update cadence

The FFIEC conformance corpus updates per the project's standard spec-versioning discipline (see GOVERNANCE.md "Specification changes"). When a new corpus version publishes:

- The working group announces the corpus version with a **90-day vendor re-attestation grace period**.
- During the grace period, prior-corpus attestations remain "Active" in the registry with a "re-attestation pending" annotation.
- After the grace period expires, the sub-committee revokes the prior-corpus attestation per the revocation procedure if no re-attestation has been submitted.

The 90-day grace period is calibrated against typical vendor product-release cadence (one or two major releases per year) and the test-execution time required to run the corpus against a new version (typically days to weeks at vendor scale).

## Institution-side consumption (normative)

### CC8.1 control description language

The institution's CC8.1 control description names:

1. The attestation evidence consumed in vendor selection. The institution names the attestation document, signature artifact, vendor public key, and the date of validation.
2. The cadence of attestation re-validation. The institution names "annually as part of the vendor-management cycle, and whenever the registry indicates a revocation or supersession of a current vendor's attestation."
3. The trigger conditions for ad-hoc re-validation. The institution names: corpus-version change, vendor product-version change that affects chain-of-custody behavior, registry revocation, vendor-side notification of an exception that affects the institution's deployment.
4. The validation procedure operated. The institution names: signature validation against the vendor's published key, corpus-version match confirmation, exception-list review against the institution's deployment posture.

### Validation evidence retained

The institution retains:

- The signed attestation document.
- The signature artifact.
- The vendor's public key (cached out-of-band).
- The validation log: a record per validation noting the validating individual, the date, the validation outcome, and the cross-check against registry-current state.
- The registry-consultation record: a record of registry queries (date, vendor entries reviewed, status of those entries on the consultation date).

Retention is the longer of (a) the institution's chain-event retention period, or (b) seven years from validation.

### SOC test of the attestation-validation control

The institution's SOC team tests CUEC-VND-06 (see `docs/control-map/CUECs.md`):

1. Sample a population of vendor relationships in scope.
2. Confirm the institution's vendor-management evidence repository contains the attestation, signature, public key, and validation log per the scope vendor.
3. Confirm the validation log shows the validation was performed within the institution's documented cadence.
4. Confirm the signature validation, corpus-version match, and exception-list review are recorded in the validation log.
5. Confirm any registry revocation or supersession events during the period triggered the institution's documented response.

Findings on the control test follow the institution's standard SOC-finding disposition.

### Annual re-validation

The institution's vendor-management cycle includes annual re-validation:

1. Re-fetch the vendor's published public key (a key rotation may have occurred; the institution updates its cached key).
2. Re-fetch the current attestation from the registry (the registration may have been superseded by a newer attestation).
3. Re-validate the current attestation's signature.
4. Confirm the corpus version named in the current attestation matches the corpus version the institution runs against in the new period.
5. Update the validation log.

## Failure scenarios

### Vendor refuses to attest

If a vendor declines to produce a conformance attestation, the institution treats this as a control gap and escalates per its vendor-management procedure. The institution's CC8.1 control description names the gap; the institution either (a) requires the vendor to attest within a defined remediation window, (b) operates compensating controls (e.g., the institution itself runs the conformance corpus against the vendor's product on each release and retains the institution-run results as the vendor-conformance evidence), or (c) replaces the vendor.

The compensating-control posture (option b) is acceptable for a transition window but is operationally costly — the institution shoulders the test-execution responsibility the vendor would otherwise own. Most institutions move to option (a) or option (c) on a defined timeline.

### Attestation found to misrepresent

If a subsequent finding establishes the attestation misrepresents the test results — vendor self-disclosure, external security researcher disclosure, working-group cross-check identifying inconsistency — the working group revokes the attestation per the revocation procedure. The institution's response on receiving notification:

1. Pull the revoked attestation from active vendor-management evidence.
2. Re-evaluate the vendor relationship per the institution's standard policy. The disposition depends on the misrepresentation's materiality (a single skipped case rationalized incorrectly is different from a fabricated test-results hash).
3. Document the re-evaluation outcome and any consequent control changes.

The institution's IR playbook MAY include a Scenario for vendor-attestation revocation; the project's IR Scenario catalog references this case.

### Corpus-version change

When the FFIEC conformance corpus publishes a new version:

1. The working group announces the new corpus version with the 90-day vendor re-attestation grace period.
2. The vendor produces a re-attestation against the new corpus version within the grace period.
3. The institution updates its corpus consumption on its standard cadence (synchronized with the institution's verifier and ledger update cadence).
4. The institution's vendor-management evidence references both the prior-corpus and current-corpus attestations during the institution's transition window. The institution's control description names which corpus version is currently authoritative for the institution's deployment.

If a vendor fails to re-attest within the 90-day grace period, the institution treats the vendor's prior attestation as superseded-but-not-current. The institution's vendor-management cycle escalates per its standard procedure.

## Examples

### Example A — Standard attestation flow

Vendor X produces product Y v1.2. Vendor X runs FFIEC conformance corpus v1.0 against product Y v1.2 in the vendor's release-validation environment. All cases pass; no exceptions. Vendor X's CTO signs the attestation; the vendor publishes the attestation at the documented URL and the public key at the vendor's `.well-known` URL. Vendor X submits the attestation to the working group; the sub-committee reviews and accepts within 18 days. The attestation appears in the public registry.

Institution Z is evaluating Vendor X for chain-of-custody hosting. Institution Z's vendor-management team consults the registry, finds Vendor X's product Y v1.2 attestation, validates the signature against the vendor's published key, confirms corpus v1.0 matches the corpus the institution runs, archives the attestation document, signature, public key, and validation log. Institution Z proceeds with vendor selection.

One year later, Institution Z's vendor-management cycle re-validates: re-fetches the vendor's key (no rotation occurred), re-fetches the current attestation (Vendor X has issued an attestation for product Y v1.3 against corpus v1.0; institution Z confirms the current attestation), re-validates, updates the validation log.

### Example B — Corpus-version transition

The working group publishes corpus v1.1, announcing a 90-day vendor re-attestation grace period. Vendor X has 90 days to re-attest product Y v1.3 against corpus v1.1.

Institution Z is in the middle of its annual vendor-management review when the corpus update lands. Institution Z's evidence base currently references Vendor X's product Y v1.3 / corpus v1.0 attestation. Institution Z's vendor-management note records: "Pending re-attestation under corpus v1.1; expected by [date 90 days hence]; vendor-watchlist status: active."

Vendor X completes re-attestation against corpus v1.1 within 60 days of publication. The registry adds the new attestation; the prior corpus v1.0 attestation moves to "superseded." Institution Z updates its evidence base to reference the new attestation, re-validates per its standard procedure, removes the vendor-watchlist annotation.

### Example C — Vendor refuses to attest

Vendor W ships a chain-of-custody product but declines to produce a vendor-conformance attestation, claiming SOC reporting is sufficient. Institution Z's vendor-management cycle identifies the gap: the SOC report covers operational controls but does NOT attest the implementation passes the FFIEC conformance corpus.

Institution Z's CC8.1 control description records the gap. The institution operates the compensating control of running the corpus itself against Vendor W's product on each release, archiving the institution-run results as the vendor-conformance evidence. The institution's vendor-management cycle includes a quarterly review of whether Vendor W is willing to attest; the institution's vendor-management committee evaluates whether to require attestation as a contractual condition of renewal.

## Optional: certification mark

The working group MAY operate a **"FFIEC chain-of-custody conformant"** certification mark. Vendors with a current attestation in the registry may use the mark in marketing material under the working group's mark-usage rules.

### Mark-usage rules (normative if the mark is operated)

- The mark MAY be used only by vendors with a current (non-superseded, non-revoked) attestation in the registry.
- The mark MUST be accompanied by the corpus version the attestation is against (e.g., "FFIEC chain-of-custody conformant — corpus v1.0").
- The mark MUST link to the vendor's registry entry (a URL the working group provides at attestation acceptance).
- The mark MUST NOT be used on product versions different from the version under attestation (a vendor with a v1.2 attestation does NOT use the mark for v1.3 until v1.3 is independently attested).
- The mark MUST NOT be used during a corpus-version-transition window after the vendor's prior attestation supersedes — the vendor either re-attests within the grace period and continues mark use, or stops using the mark when the attestation supersedes.

The working group monitors mark usage as part of the quarterly registry review. Mark-usage violations are addressed per the working group's standard appeals and remediation discipline.

### Why the mark is useful

The mark serves vendor marketing and institution-side trust signaling. An institution's procurement team that recognizes the mark on vendor marketing material has a fast-path identifier for "this vendor is at minimum in the registry"; the institution's vendor-management team performs the full validation procedure as the load-bearing trust path.

The mark is intentionally NOT a substitute for the full validation procedure. The mark is a recognition signal; the institution's evidence base is the validation log, not the mark.

## Summary

The vendor-conformance attestation procedure establishes a project-side mechanism for institutions to confirm vendor-implementation conformance against the FFIEC conformance corpus. The procedure is separate from vendor SOC reporting (which covers operational controls). The procedure operates through:

- An annual vendor-produced signed attestation against a specific corpus version.
- A project-side working-group sub-committee that reviews and publishes attestations to a public registry.
- Institution-side validation, archival, and ongoing-review per CUEC-VND-06.
- A revocation procedure for attestations found to misrepresent or for vendors that fail subsequent corpus versions.
- Optional certification mark usage with documented rules.

The procedure closes the trust gap between "vendor's word" and "institution's examination-grade evidence." Institutions consume the attestation as part of CC8.1 control discipline; the working group operates the registry as part of project-side governance commitments named in `GOVERNANCE.md`.
