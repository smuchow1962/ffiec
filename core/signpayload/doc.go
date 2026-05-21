// Package signpayload reconstructs the byte form the HSM signed for
// a daily Merkle seal, dispatched on the seal record's
// sign_payload_version field per §4.3.
//
// Spec authority: §4.3 defines four byte forms:
//
//   - Pre-amendment 6-line form (sign_payload_version absent).
//   - v1.0a 10-line form (sign_payload_version = "v1.0a").
//   - v1.0b 12-line form (sign_payload_version = "v1.0b"; binds
//     key_versions_canon + hex(kms_handle_uris_digest)).
//   - v1.0c 13-line form (sign_payload_version = "v1.0c"; binds
//     hex(operational_events_log_root) per §10.79).
//
// Every line is separated by a single 0x0A byte; the terminal field
// carries NO trailing 0x0A. Encoding is exact: hex is lowercase,
// dev_mode is a single ASCII byte "0" or "1", key_versions_canon is
// ASCII comma-separated ascending decimal. A wrong byte anywhere
// produces a different signature input and the verifier reports
// `signature verification failed` at §7 step 11 rather than a more
// specific reason — which is why this package's byte-equivalence
// against the .NET and Python reference implementations is the most
// load-bearing thing in the cryptographic surface.
//
// Status: stub. The dispatcher lands in Commit 5 of the verifier
// upgrade, alongside §7 step 11 signature verification.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/SignPayload.cs is the byte-form
// authority. The Python reference folds reconstruction into
// _resealing.py + the verifier. Three-way symmetry on every
// dispatch case is the conformance bar.
//
// Commit-5 implementation shape (per Richard's cross-consult
// 2026-05-21):
//
//   - One `JoinLines` helper mirroring .NET + Python's DRY shape:
//     magic line carries its own embedded `\n`; the rest is a slice
//     `\n`-joined and concatenated. Don't write three separate
//     builders — dispatch on `sign_payload_version` and call the
//     same helper with different field counts.
//   - Lowercase hex via `hex.EncodeToString`. URI-sort via
//     `sort.Strings` (byte-ordinal — same gotcha as JCS; do NOT
//     reach for text/collate).
//   - Empty-day kms_handle_uris_digest collapses to SHA-256("") =
//     e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
//     — pin this hex constant in code so a reader can trace it
//     directly to the spec text.
//   - Typed `SignPayloadVersionError struct { Version string }` with
//     `Unwrap() error` so callers wrapping `errors.As` catch it
//     transparently.
//   - Vector gating: 018 + 019 for v1.0a/b (both materialized);
//     v1.0c defers to a follow-up commit if vector 027 is still
//     stub-only when Commit 5 lands — do not invent the test pin.
//
// Forward Go-specific risk Richard flagged: today's vector 018 has
// all-ASCII URIs only, so the byte-ordinal sort matches across
// .NET / Python / Go. A non-ASCII URI test extension would surface
// any future collation-sort drift; Heather may land an 018b
// extension if it merits.
package signpayload
