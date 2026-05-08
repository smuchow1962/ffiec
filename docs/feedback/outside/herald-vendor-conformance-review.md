# Outside drop — Herald-ecosystem conformance review against `chain-of-custody-v1.md`

## Persona

**Reviewer.** Independent technical reviewer assessing whether a candidate vendor implementation — the Herald ecosystem (Herald.Core, Herald.Compliance, Herald.Server, Herald.Py) — can satisfy the v1.0-final spec as published. Background: financial-services compliance engineering and applied cryptography on tamper-evident logging primitives. No prior commercial relationship with the Herald team; reading the spec text first, the published Herald SDK surface second, and comparing.

**Reading angle.** The spec's v1.0-rework changelog cites the Herald HMAC-SHA-256 + HKDF audit-chain construction as the resolved-through-Auditor-rounds prior art the rework anchored on. That makes Herald the natural first vendor to test conformance against. Treating the spec as the authoritative bar and asking: where does Herald already satisfy the bar, where is it close, and where is the gap a Herald customer would have to close before claiming v1.0 conformance?

**Method.** I read the spec, then consulted the publicly-published Herald artifacts:

- `Modules/Herald.Compliance/src/Audit/Chain/` — observable SDK shapes (`HmacChainCore`, `KeyDerivation`, `ChainVerifier`, `AuditChainEntry`, `AuditFileHeader`, `Constants`, `CanonicalJson`).
- `Modules/Herald.Compliance/src/Audit/` — `ImmutableAuditFileSink`, `AuditFileSigner`, `AuditFileManifest`, `AuditFileVerifier`, plus the new pipeline-level decorator at `Modules/Herald.Compliance/src/Audit/Pipeline/HmacChainPipelineLogger`.
- `Modules/Herald.Compliance/src/Audit/Kms/IKmsClient` — the verifier's IKM-source abstraction.
- `Herald.Py` published cross-language fixture at `Herald.Py/tests/fixtures/chain_vectors.json`.

I deliberately did **not** read the Herald team's prior round notes or change history — the assessment treats Herald as a candidate vendor I would evaluate from its shipped surface alone.

---

## Headline summary

Herald implements **Primitive 1** (HMAC-SHA-256 + HKDF chain-at-capture) at the construction-architecture level, with the right per-entry stamp, the right verifier order (fingerprint-before-MAC, `expected_prev_hash` into recompute, constant-time compares), and the right append-only sink contract. The gap to FFIEC v1.0 conformance is principally **byte-form**: Herald's HKDF salt + info constants live in the `herald.*` namespace; the spec demands `ffiec.*`. A vanilla Herald deployment produces byte-divergent `session_key`, `key_fingerprint`, `hkdf_inputs_digest`, and `payload_hash` from `spec/test-vectors/chain_vectors.json`. Same primitive, different bytes — non-conforming as published unless the constants are made configurable.

**Primitives 2, 3, and 4 are not present** in the Herald ecosystem today as I read the public surface:

- No daily Merkle seal across the tenant-day boundary.
- No HSM-rooted **Ed25519** root signature (Herald's `AuditFileSigner` ships RSA-PSS-SHA-256 and HMAC-SHA-256, both per-file, both with a different `sign_payload` shape than spec §4.3).
- No OTLP attribute encoder under the `ffiec.chain.*` namespace; the chain seal lands on Herald's in-memory `LogEvent.Context` as a nested `audit_chain` block. The wire-out is whatever sink the inner pipeline ends in.

That is the true gap-shape an honest conformance assessment would surface: Herald is approximately one-of-four primitives in, with the per-event construction at near-conformance and the day/seal/signature/wire layers absent. A Herald customer claiming v1.0 conformance today would be claiming partial coverage.

The findings below are structured as questions to the spec working group, in the format suggested by `docs/feedback/README.md`. Several of them are also implicit asks of the Herald team; I am surfacing them through the spec because the spec's working group is the natural arbiter of what conformance requires.

---

## Questions

**Q1. Will the spec keep `HKDF_SALT = b"ffiec.chain-of-custody.v1.salt"` and `HKDF_INFO_BASE = b"ffiec.chain-of-custody.v1.info"` byte-locked at v1.0, or is there a conformance pathway for vendors who chose a vendor-namespaced constant before the rework (e.g., `herald.audit-chain.v1.*`)?**

**Status:** Gap

The spec test-vector corpus (`spec/test-vectors/chain_vectors.json` and `spec/test-vectors/README.md` §"Constants") locks the salt and info-base byte values. The `hkdf_inputs_digest = SHA-256(salt || info || length_LE32)` propagates these constants into both the file header and the `sign_payload` (§4.3), so a vendor running `herald.*` constants produces a different `hkdf_inputs_digest` byte sequence and a different signed `sign_payload`. The chain primitives compose identically, but the verifier walking a Herald file under spec §7 step 2 fails the `header HKDF inputs do not match running v1 inputs` check.

This is the load-bearing decision for "can a vendor that shipped before v1.0-final claim conformance?" Three possible spec answers:

- **(a) Hard byte-lock.** Vendors flip their constants to `ffiec.*` in a new format_version, the Herald pre-flip chain becomes a `herald.audit-chain.v1` chain (vendor-namespaced) and the FFIEC chain becomes `ffiec.chain-of-custody.v1`. Two contracts coexist; vendor migration is a flag-day.
- **(b) Vendor-flag mode.** A conforming vendor SDK ships a configuration knob that selects FFIEC constants when the vendor is asked to produce conforming output. Pre-flip data remains under the vendor's constants and is documented as non-FFIEC-conformant; post-flip data is FFIEC-conformant. Spec adds one normative line: "implementations MAY parameterize the salt and info-base constants at SDK-construct time, but a chain entry is FFIEC-conformant only when produced under the constants in spec §4.1." The `hkdf_inputs_digest` field already detects which constants were in force, so the on-disk witness is unambiguous.
- **(c) Per-tenant constants.** The spec adds a tenant-namespace prefix (e.g., `ffiec.chain-of-custody.v1.salt:bank-of-acme`) that lets the same SDK serve multiple regulatory regimes simultaneously. Substantially more work; closes the dual-jurisdiction case.

I would lift (b) into the spec body. Herald's published `chain_vectors.json` shows the construction is identical at the function level — only the byte values of two constants differ. Forcing every vendor with prior-art under namespace-X to fork is friction without security payoff. Spec §3 already constrains the `tenant_id` character set so the `info` boundary is unambiguous; an analogous one-line constraint on the constants ("MUST be the v1 byte values during FFIEC-conformant operation; MAY be alternate bytes during vendor-internal operation, in which case the chain is not FFIEC-conformant") closes the question.

**Q2. Does the spec require a verifier to refuse a chain whose constants do not match v1, even if `format_version == "v1"` and the construction is otherwise correct?**

**Status:** Answered

Spec §7 step 2 (HKDF inputs digest) does this implicitly: the verifier recomputes `expected_hkdf_inputs_digest` from the v1 constants and constant-time compares against `header.hkdf_inputs_digest`. A Herald-produced file under `herald.*` constants fails this check at the file-header pre-flight, with the right error message ("header HKDF inputs do not match running v1 inputs"). The mechanism is correct; the practical effect is what Q1 addresses. The verifier behavior is unambiguous.

**Q3. Is RFC 8785 JCS canonical-JSON conformance a normative bar at v1.0, or is "produces byte-identical canonical bytes for the field set in §5" the actual bar?**

**Status:** Gap

Spec §5 says: "The canonical-JSON form used inside `payload_hash` MUST follow [RFC 8785]" and "Two implementations of v1 MUST produce byte-identical canonical bytes for the same logical event." Those are not the same bar. RFC 8785 has rules about float canonicalization (no exponent if avoidable, no trailing zeros, NaN/Infinity rejected), Unicode (UTF-8, escape only the minimal control set per RFC 8259 §7), and number formatting that go beyond "sort keys + no whitespace."

The Herald public surface for canonical JSON (`MMP.Herald.Compliance.Audit.Chain.CanonicalJson`) carries an XML-doc note describing it as a "simple-cases scaffold" with deferred full orjson byte-equivalence. The deferred edge cases include float-format and Unicode-escape parity. A vendor reading "Phase 2 lifts the bar to full orjson byte-equivalence" alongside the spec's MUST statement could reasonably conclude that v1.0 conformance requires the deferred work to land, not the scaffold.

The spec test vectors at `spec/test-vectors/chain_vectors.json` exercise integer timestamps, base64-encoded byte strings, and ASCII-only attribute values — they do **not** exercise the JCS edge cases (no float values, no non-ASCII Unicode, no NaN/Infinity rejection paths). A vendor passing the published test vectors byte-for-byte can still fail RFC 8785 conformance on the first event whose `attributes` map includes a float or a non-ASCII string.

I would either tighten the test-vector corpus to include a JCS edge-case fixture (case `008-jcs-edge-cases/` is on the future-cases list per `spec/test-vectors/README.md`) or relax §5 to "produces byte-identical canonical bytes for the field set, where 'canonical' means the rules in §5.x — RFC 8785 is informative but not the conformance bar." Without one of those, two vendors agreeing on the published fixture can disagree on the first non-fixture event, and the chain-of-custody guarantee is silently broken.

**Q4. Does the spec require mid-write truncation refusal at the verifier (last-byte-must-be-`\n`)? If so, what is the exact discriminator a verifier must implement, given that some line-oriented file formats use `\r\n` on Windows or omit the terminator on the final line?**

**Status:** Partial

Spec §4.1 says: "The verifier MUST refuse to verify a file whose last byte is not `\n`." Spec §6 echoes the same rule. A vendor SDK whose verifier consumes the file via a stdlib line reader (Python's `for line in file:`, .NET's `StreamReader.ReadLine()`, Go's `bufio.Scanner.Scan()`) gets line-stripped output and silently tolerates a missing trailing `\n` — none of those readers raises on a final-line-without-newline. A conforming verifier therefore needs an out-of-band check that opens the file at byte level and reads the last byte before letting the line reader run.

The Herald public surface for verification (`MMP.Herald.Compliance.Audit.Chain.ChainVerifier.VerifyChain` and the `ImmutableAuditFileSink.VerifyChain` overload) does not advertise a mid-write-truncation refusal in its XML doc. Whether the implementation does the byte-level check is not visible from the SDK shape alone. A vendor reading §4.1 and reaching for `StreamReader.ReadLine()` would build a non-conforming verifier without realizing it.

I would add a one-paragraph implementation note to §7 that names the failure mode ("a stdlib line-reader will silently tolerate the missing terminator; verifiers MUST do an out-of-band byte-level check, e.g., `seek(end-1, SEEK_END); read(1)` and assert the byte equals `\n`"). Same point in `docs/design/02-chain-construction.md` if it covers writer-side implementation. The `audit_file.truncation_detected` operational event in §10.2 is the right thing to emit; the spec just needs to say *how* the verifier detects the condition.

**Q5. The spec normatively requires Ed25519 in HSM custody (§4.3) for daily root signing. Is RSA-PSS-SHA-256 a conformant alternative for v1.0, or is Ed25519 the only conforming algorithm?**

**Status:** Answered

§4.3.2 names Ed25519 as the v1.0 algorithm and reserves dispatch on a per-seal `algorithm` field for future post-quantum candidates. RSA-PSS is not in the algorithm list. A vendor whose signer ships only RSA-PSS-SHA-256 (the Herald `AuditFileSigner.SignWithRsaManifest` path; observable from the public `AuditFileSignerAlgorithms.RsaPssSha256` constant) is non-conforming on Primitive 3 — not because RSA-PSS is cryptographically weaker, but because the spec locked the `sign_payload` byte form including the `algorithm` line, and a verifier dispatching on `seal.algorithm == "ed25519"` will not recompute against an `RSA-PSS-SHA256` signature.

The spec is unambiguous; the gap is on the vendor side, not in the spec. I am flagging it as Answered to confirm I read §4.3 + §4.3.2 correctly. A Herald customer pursuing FFIEC conformance needs the Herald team (or a third-party plugin) to add Ed25519 alongside the existing RSA-PSS path.

**Q6. The spec requires `dev_mode = true` seals to be refused under `--strict` (§4.2 schema, §10.7) and requires the dev adapter to be compile-time excluded from production builds (§10.7). Is run-time gating-only sufficient if the vendor stamps `kms_handle_uri = "plaintext-dev"` on every entry?**

**Status:** Gap

§10.7 is explicit: "Run-time environment-variable gating is NOT sufficient — a misconfigured deployment that flips the flag must NOT be able to bring the software adapter online in production." Herald's published surface includes `AuditChainConstants.KmsHandlePlaintextDev = "plaintext-dev"` and a constructor parameter `kmsHandleUri` on `ImmutableAuditFileSink` that defaults to that marker. The marker stamps onto every entry as required. The compile-time exclusion is **not** observable on the public surface — the `ImmutableAuditFileSink` and `HmacChainPipelineLogger` constructors accept arbitrary `ikm` bytes at run time without any compile-time switch separating "production" from "dev" Herald assemblies.

A vendor following spec §10.7 strictly would need to ship two Herald assemblies (dev and production) with the in-process IKM acceptance path stripped from production. A Herald customer claiming conformance today can do this externally (CI pipeline excludes the adapter; production deploy uses a KMS-backed `IKmsClient` only) but the SDK does not enforce it.

I would lift the spec language one notch: rather than "MUST be excluded from production release builds at compile time," say "MUST be unreachable in production deployments through any combination of build-flag, packaging, and configuration that a normal misconfigured deployment cannot bypass at run time." That accommodates compile-out (the strictest pattern) and equally-strict alternatives like "production builds ship without the adapter assembly on disk; configuration cannot resurrect it." A vendor that ships a single assembly with a run-time flag is non-conformant regardless of which intent the operator declared.

**Q7. Does the spec require the `tenant_id` character-class restriction (§3, `^[A-Za-z0-9_.\-]{1,255}$`) to be enforced at the SDK boundary (refuse construction with a tenant_id containing `|`), at the verifier boundary (refuse a chain whose header `tenant_id` violates the class), or both?**

**Status:** Gap

§3 normatively defines the tenant_id character class for HKDF info-boundary safety. The spec body does not name where enforcement lives. A reasonable reading: writer-side enforcement (refuse to construct an `HmacChainWriter`/`ImmutableAuditFileSink` with a non-conforming `tenant_id`) closes the attack at the source; verifier-side enforcement (refuse a chain whose `header.tenant_id` violates the class) catches misuse from older or non-conforming writers.

Herald's public surface (`AuditFileHeader.Create`, `ImmutableAuditFileSink` constructor, `HmacChainPipelineLogger` constructor) validates `tenantId.Length > 0` only — it does not enforce the spec's character class. A `tenant_id` containing `|` would land in `KeyDerivation.BuildHkdfInfo` and the resulting `info_for_tenant` bytes would be ambiguous (`info_base || "|" || A|B` parses two ways). The chain still produces deterministic bytes, but the spec's claim that "two distinct tenant identifiers cannot produce the same `info` byte sequence" is broken at that input.

The Herald-published cross-language fixture uses `tenant-vector-1` which conforms to the class. Production deployments that accept arbitrary-character tenant IDs from external systems (multi-tenant SaaS where tenant_id derives from a CRM record) would silently allow class violations.

Spec language to add to §3 or §4.1: "Implementations MUST reject a `tenant_id` that does not match the §3 character class at SDK-construct time. Verifiers MUST reject a chain whose `header.tenant_id` does not match the class with `tenant_id violates §3 character class`." Same enforcement point, named both ends.

**Q8. Spec §4.4 names `gen_ai.request.model` and `gen_ai.response.model` as REQUIRED on any chain entry that represents a model call (§7 step 12a control-completeness check). Is this a property the SDK enforces (refuses to log a chain entry with `gen_ai.*` attributes that lacks both `gen_ai.request.model` and `gen_ai.response.model`), or is the missing-attribute case detected only at verification?**

**Status:** Partial

§7 step 12a names a verifier check ("PASS-WITH-ANOMALY (control-completeness for SR 11-7 reproducibility, NOT chain-integrity)"). The spec is silent on whether the SDK SHOULD/MUST also refuse the missing-attribute case at write time. A "verifier-only" enforcement model is acceptable cryptographically — the chain integrity is unaffected by the omission — but it puts the burden on the verifier to detect a control-completeness gap that the SDK could have caught at write time.

The Herald public surface (`HmacChainPipelineLogger`, `ImmutableAuditFileSink`) is a general-purpose log decorator + sink. It does not introspect `gen_ai.*` attributes; the chain seals whatever fields are on the event without enforcing semantic-convention completeness. A Herald customer producing FFIEC-conformant chains would need a separate validating processor in their pipeline before the chain decorator (Herald has the enricher / processor surface for this; the validating processor is not shipped today as far as I observe).

I would lift §4.4 with a one-paragraph normative requirement: "SDKs MUST refuse to write a chain entry whose attribute set includes any `gen_ai.*` namespace attribute and lacks either `gen_ai.request.model` or `gen_ai.response.model`. Verifiers MUST emit `gen_ai_model_identifier_missing` per §7 step 12a; the SDK-side refusal is a defense-in-depth against a misconfigured pipeline that would otherwise produce chains failing at audit time." Documenting both ends keeps the chain integrity-bearing AND the model-identifier completeness enforced where it matters most.

**Q9. The spec requires daily Merkle aggregation (§4.2) and HSM-rooted root signing (§4.3) at the **ledger server**, not the SDK. Is a vendor whose published surface is SDK-only (per-event chain + per-file/per-pipeline sign) conformant if the ledger-server work is supplied by a different vendor or by the institution itself?**

**Status:** Answered

§4.2 names the ledger server as the location of Primitive 2 ("on the bank's perimeter (or vendor-hosted with per-tenant key segregation)"). §4.3 names the HSM as the location of Primitive 3. The spec does not require a single vendor to ship all four primitives — and the architecture document explicitly contemplates multi-vendor topology ("an institution running an SDK from one vendor and a ledger from another vendor and a verifier from this repository should see byte-for-byte agreement on every chain artifact" — `spec/README.md`).

Herald's public surface stops at Primitive 1 (per-event chain) plus a per-file integrity envelope (`AuditFileManifest`) that is not what spec §4.3 calls a daily seal. A Herald-as-SDK + institution-supplied-ledger topology can be FFIEC-conformant: the institution's ledger server consumes Herald-emitted events, computes the daily Merkle seal, signs in the HSM, exposes the seal record in the §4.2 schema, and Herald is on the hook only for Primitive 1 byte-conformance plus the §4.4 OTLP wire.

This question is Answered, but it implies a follow-up: who ships the ledger? The spec project's `ledger/` reference implementation is in design phase per `README.md` ("Implementation of `core/`, `ledger/`, and `verifier/` follows once the spec is locked"). A Herald customer pursuing v1.0 conformance today has Primitive 1 from Herald, no ledger, and no verifier from any vendor. That's a vendor-ecosystem gap the spec working group should be tracking — Herald is positioned to supply Primitive 1, not Primitives 2-3-4.

**Q10. The §4.4 `ffiec.chain.*` OTLP attribute namespace is defined normatively. Does the spec consider an SDK that emits the chain seal onto a non-OTLP transport (e.g., NDJSON file, syslog, custom JSON over HTTPS) conformant, provided the field set is preserved and a separate OTLP encoder can re-encode the events for the wire?**

**Status:** Partial

§4.4 says: "Events MUST ship over OpenTelemetry Protocol (OTLP) using the OTel GenAI Semantic Conventions for AI-specific attributes." A strict reading mandates OTLP at the wire. A practical reading allows non-OTLP transports for delivery as long as the OTLP envelope is "encoded within their payload" (§4.4 last paragraph addresses indirect transports — Kafka, Kinesis — that wrap an OTLP envelope).

Herald's chain decorator stamps the seal onto `LogEvent.Context` as a nested `audit_chain` block (observable from `HmacChainPipelineLogger.ContextKey = "audit_chain"`). Whether that lands as `ffiec.chain.*` OTLP attributes or as some other on-the-wire encoding depends on the inner sink: a Herald customer feeding events to an OTLP-OTel sink would need a translator that converts the `audit_chain` block into the eight `ffiec.chain.*` attribute names spec §4.4 requires.

A spec clarification would help: either (a) "the `ffiec.chain.*` attribute names are the wire-form contract; intermediate in-process representations MAY use vendor-specific naming, but the OTLP wire emission MUST translate to `ffiec.chain.*`" or (b) "implementations MAY use vendor-specific in-process attribute names; the §4.4 attribute table is the conformance contract the verifier reads from the captured events." Either is workable; the current text is silent on the in-process-versus-wire distinction.

**Q11. Spec §5 requires the canonical-form to EXCLUDE every chain-stamp field listed in §4.1 (the verifier-relevant fields) AND requires the canonical-form to INCLUDE the cross-run linkage fields (`parent_run_id`, `parent_seq`, `dag_parents`). Are the `tenant_id` and `run_id` fields in or out of the canonical bytes?**

**Status:** Answered

Spec §5 enumerates the included fields explicitly: "The OTel envelope: `trace_id`, `span_id`, `parent_span_id`, `name`, `timestamp_ns`, `duration_ns`, `attributes`, `resource`, `severity`, `kind`, `chain_kind`. ... The explicit chain metadata: `tenant_id`, `run_id`, `captured_at`. ... The cross-run linkage fields: `parent_run_id`, `parent_seq`, `dag_parents`." So `tenant_id` and `run_id` ARE in the canonical bytes; the chain-stamp fields (which include `seq`, `prev_hash`, `payload_hash`, `key_version`, `key_fingerprint`, `format_version`, `mac_computed_at_utc`, `kms_handle_uri`, `algorithm`) are out.

Herald's `CanonicalEventBuilder` (observable from the published Compliance project structure) builds a map containing `tenant_id` and `run_id` and excludes the `audit_chain` block at canonicalization time — same shape. The published `chain_vectors.json` event_canonical_hex confirms `tenant_id` and `run_id` are in the canonical bytes. Answered.

**Q12. Does the spec consider a SOC 2 Section-4 description that names Herald.Compliance as the SDK vendor, an institution-supplied ledger server, and the FFIEC reference verifier as the verifier, conformant — given that Herald.Compliance ships a per-file `AuditFileManifest` that overlaps with but does NOT match the spec's §4.2 daily-seal schema?**

**Status:** Gap

The institution would have to choose: either (a) discontinue the Herald `AuditFileManifest` sidecar (it is not a §4.2 seal; it is a per-file RSA-PSS or HMAC envelope, with a different schema and a different signature algorithm), or (b) keep it as supplemental institution-internal evidence and ship a separate FFIEC §4.2 seal alongside. Option (b) is the practical answer; option (a) loses the per-file evidence Herald customers may already depend on.

The spec is silent on supplemental envelopes. A spec note clarifying that "a vendor MAY ship per-file integrity envelopes for institution-internal use; FFIEC conformance is determined by the §4.2 daily seal alone, and per-file envelopes do not satisfy §4.2" would prevent confusion at audit time. The §4.2 seal record schema and Herald's `AuditFileManifest` schema differ in algorithm (`ed25519` vs `RSA-PSS-SHA256`/`HMAC-SHA-256`), in cadence (daily vs per-file-roll), in coverage (Merkle root over a tenant-day vs file-digest over one chain file), and in `sign_payload` shape. Naming both as "the seal" in audit documentation would mislead an examiner.

---

## Per-role roll-up

| Aspect | Status |
|---|---|
| Primitive 1 — HMAC chain construction (HKDF + HMAC + tenant binding + per-entry stamp) | Answered (architecturally) |
| Primitive 1 — byte-level conformance to `chain_vectors.json` | Gap (Q1: vendor uses `herald.*` constants, spec demands `ffiec.*`) |
| Primitive 1 — `expected_prev_hash` into MAC recompute | Answered |
| Primitive 1 — fingerprint-before-MAC + constant-time compares | Answered |
| Primitive 1 — IKM minimum 32 bytes | Answered |
| Primitive 1 — per-entry `format_version` + `key_version` + `key_fingerprint` + `mac_computed_at_utc` + `kms_handle_uri` | Answered |
| Primitive 1 — file-header `hkdf_inputs_digest` | Answered |
| Primitive 1 — RFC 8785 JCS conformance for the canonical bytes | Partial (Q3: scaffold today; full edge-case parity deferred) |
| Primitive 1 — mid-write truncation refusal | Partial (Q4: stdlib line readers silently tolerate missing trailing `\n`) |
| Primitive 1 — `tenant_id` character-class enforcement at SDK + verifier boundary | Gap (Q7) |
| Primitive 2 — Daily Merkle seal across the tenant-day boundary | Gap (Q9: not in Herald; spec says ledger-server, OK if institution supplies) |
| Primitive 3 — HSM-rooted Ed25519 root signing with §4.3 `sign_payload` shape | Gap (Q5: Herald ships RSA-PSS / HMAC, not Ed25519; per-file, not daily) |
| Primitive 3 — algorithm rotation + dual-algorithm transitional period | Out-of-scope at v1.0 (spec defers to v1.x post-quantum) |
| Primitive 4 — OTLP wire encoding under `ffiec.chain.*` namespace | Gap (Q10: in-process `audit_chain` block today; OTLP translator absent) |
| Primitive 4 — `gen_ai.request.model` / `gen_ai.response.model` SDK enforcement | Partial (Q8: §7 step 12a verifier-only today; SDK-side refusal would close at source) |
| Operational — software-key adapter compile-time exclusion | Gap (Q6: run-time `kms_handle_uri` marker only; no compile-out) |
| Operational — append-only enforcement at app + DB layer | Answered (§10.3 satisfiable; OS-level append shown in `ImmutableAuditFileSink`) |
| Operational — fail-closed verifier semantics | Answered |
| Operational — vendor-flag conformance pathway | Gap (Q1: spec is silent on vendor-prior-art namespace) |
| Multi-vendor topology — SDK from one vendor + ledger from another + verifier from FFIEC repo | Answered (Q9: spec contemplates this) |
| Per-file integrity envelope vs §4.2 daily seal naming clarity | Gap (Q12: Herald `AuditFileManifest` overlaps; spec silent on supplemental envelopes) |

| Status | Count |
|---|---|
| Answered | 8 |
| Partial | 4 |
| Gap | 7 |

---

## Where I would prioritize

If the spec working group treats my findings as a closure backlog:

1. **Q1 (vendor namespace pathway).** This is the single decision that changes whether Herald — or any vendor with prior-art under a different namespace — can claim conformance without a flag-day migration. One paragraph in §4.1 + §3 closes it. Highest leverage.
2. **Q5 + Q9.** Combined: name explicitly which primitives a SDK-only vendor is responsible for, and which the institution's ledger + HSM are responsible for. Spec already implies the answer; making it explicit prevents the "Herald is FFIEC-conformant" overclaim at procurement.
3. **Q3 (JCS bar).** Either tighten the test vectors to include JCS edge cases (case `008-jcs-edge-cases/` already on the future-cases list — promote it to v1.0) or relax §5 to a field-set-and-canonical-rules bar that is independently testable. The current state has two MUSTs that don't agree.
4. **Q4, Q6, Q7, Q8.** These are operational-discipline questions where the spec's normative requirement is in the right place but the enforcement boundary is not named. Each is a one-paragraph spec patch. Together they would close most of the "vendor implementation looks conformant but isn't" risk.
5. **Q12.** Doc-only clarification. Prevents an examiner from treating a per-file envelope as a §4.2 seal at audit time.

Q2 and Q11 are confirmations — no action needed.

---

## Stopping criterion

Per `docs/feedback/README.md`: this drop carries 7 Gap and 4 Partial findings. Until those close (or the spec working group records reasons not to close them), this drop should remain open in the convergence tracker.

I would expect at least Q1, Q5, and Q9 to close before any vendor — Herald or otherwise — issues a "we are FFIEC v1.0-conformant" public claim. The other items are spec-quality issues a working-group cycle can resolve in normal cadence.
