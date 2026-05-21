// Package constants holds the FFIEC chain-of-custody v1 byte-pinned
// constants. Every value in this package is normative per spec §4.1
// (HKDF salt + info), §4.2 (Merkle tree constants), and §4.1.2 (the
// FFIEC-conformance posture).
//
// Changing a value here is a wire-format break — the spec's
// hkdf_inputs_digest is computed over the salt + info bytes, so any
// edit invalidates every existing FFIEC-posture chain. Vendor-flag
// mode (§4.1.2) lets implementations parameterize the salt + info at
// SDK-construct time; the FFIEC-conformance values live here and are
// the bar a verifier asserts when run under --posture=ffiec.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/Constants.cs uses the byte-identical
// "ffiec.chain-of-custody.v1.salt" / ".info" UTF-8 strings; the
// Python reference at Herald.Py/src/herald/_crypto/__init__.py
// exports the same byte sequences. Cross-language byte-equivalence
// depends on these three definitions matching.
package constants

// FFIEC-conformance HKDF parameters per §4.1. These byte strings are
// the salt and info-base for HKDF-SHA-256 derivation of the per-tenant
// session key. The info parameter the SDK passes to HKDF is
// HKDFInfoBase || "|" || utf8(tenant_id) — the separator byte is
// asserted to be 0x7C (`|`) by §3's tenant_id character class, which
// excludes that byte so the boundary is unambiguous.
const (
	HKDFSalt     = "ffiec.chain-of-custody.v1.salt"
	HKDFInfoBase = "ffiec.chain-of-custody.v1.info"

	// HKDFOutputLength is the HKDF-SHA-256 output length in bytes,
	// fixed at 32 (the session key length). The value is bound into
	// hkdf_inputs_digest as its 4-byte little-endian encoding per §3.
	HKDFOutputLength = 32

	// InfoTenantSeparator is the ASCII pipe `|` (0x7C) that separates
	// HKDFInfoBase from utf8(tenant_id) in the HKDF info parameter.
	// Per §3 the tenant_id character class excludes 0x7C so the
	// boundary is unambiguously parseable.
	InfoTenantSeparator = "|"
)

// FormatVersion is the v1 wire-format identifier per §0 Version
// policy. A v1 verifier rejects any other value at §7 step 1
// (format_version check) with the spec-mandated reason string.
const FormatVersion = "v1"

// Merkle prefix bytes per RFC 6962 §2.1, normative per §4.2. A
// leaf hash is SHA-256(0x00 || payload_hash); an internal node hash
// is SHA-256(0x01 || left || right). The domain separation prevents
// presenting an internal node as a leaf.
const (
	MerkleLeafPrefix byte = 0x00
	MerkleNodePrefix byte = 0x01
)

// GenesisPrevHashLen is the fixed 32-byte length of every chain
// entry's prev_hash field per §4.1 inviolate property 2. The genesis
// entry (seq=1) MUST carry 32 zero bytes; thereafter prev_hash is the
// previous entry's payload_hash.
const GenesisPrevHashLen = 32

// KeyFingerprintLen is the byte length of the truncated fingerprint
// per §4.1 — SHA-256(utf8(tenant_id) || ikm)[:16].
const KeyFingerprintLen = 16
