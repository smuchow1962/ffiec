package verify

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/core/signpayload"
)

// §7 step 11 — daily-seal Ed25519 signature verification (§4.3).
//
// Step 11 is the last cryptographic check in the §7 walk: it proves the
// Merkle root sealed at step 10 was signed by the institution's HSM key,
// not by anyone with ledger-write access. The verifier consumes ONLY
// public-key material — it never holds or reads the private seed. The
// public key is supplied out-of-band (corpus-published pub.hex, or the
// verifier config); the chain does not self-assert it.
//
// The check is two independent assertions, in this order:
//
//  1. **Reconstruct-and-compare.** Rebuild the §4.3 sign_payload from the
//     seal's structured fields (algorithm, tenant_id, merkle_root, …) via
//     core/signpayload.Build — the same builder vectors 018/019/020/035
//     gate byte-for-byte — and confirm it equals the seal's published
//     sign_payload_hex. A published payload that does not reconstruct from
//     the real fields is a step-11 failure: it means the signature covers
//     bytes that do not bind the day's actual root/tenant/version. Doing
//     this BEFORE the Verify call means a tampered tenant_id or a swapped
//     sign_payload_version is caught even when an attacker also forged a
//     matching signature over their substituted payload.
//
//  2. **Verify.** Ed25519-verify the decoded signature_b64 over the
//     reconstructed bytes under pub. A garbage signature (N004), a
//     signature for a different tenant's payload (N005), or any tamper
//     that changed the signed bytes fails here.
//
// Both failures render the §7 normative reason "signature verification
// failed" with exit code 1, matching the N004/N005 pins. A structural
// problem (unparseable version, undecodable signature, wrong key length)
// renders the same reason — the seal did not verify — at exit code 1,
// because a malformed signature field is a verification failure, not a
// pre-flight input error.

// CheckSealSignatureV1 runs §7 step 11 over the audit-file seal using the
// supplied institution public key. It returns the clean PASS outcome when
// the reconstructed sign_payload both matches the published bytes and
// verifies under pub, or a FAIL("11", "signature verification failed")
// otherwise.
//
// pub MUST be a 32-byte Ed25519 public key. A wrong-sized key is a
// configuration error surfaced as a step-11 failure (the verifier cannot
// verify against a key it cannot use).
func CheckSealSignatureV1(seal AuditSeal, pub ed25519.PublicKey) Outcome {
	const reason = "signature verification failed"

	if len(pub) != ed25519.PublicKeySize {
		return fail("11", reason)
	}

	reconstructed, err := reconstructSignPayload(seal)
	if err != nil {
		return fail("11", reason)
	}

	// Assertion 1: the published sign_payload_hex MUST equal the bytes the
	// verifier reconstructs from the seal's structured fields. This binds
	// every structured field (tenant_id, merkle_root, sign_payload_version,
	// …) to the signature even before the Verify call.
	published, err := hex.DecodeString(seal.SignPayloadHex)
	if err != nil {
		return fail("11", reason)
	}
	if !bytes.Equal(published, reconstructed) {
		return fail("11", reason)
	}

	// Assertion 2: the signature verifies over the reconstructed bytes.
	sig, err := base64.StdEncoding.DecodeString(seal.SignatureB64)
	if err != nil {
		return fail("11", reason)
	}
	if !ed25519.Verify(pub, reconstructed, sig) {
		return fail("11", reason)
	}
	return pass()
}

// reconstructSignPayload maps the audit-file seal into a signpayload.Seal
// and builds the §4.3 byte form for the seal's sign_payload_version. The
// builder is the same one the positive sign_payload vectors gate, so the
// reconstruction is the spec's normative byte form, not a re-derivation.
func reconstructSignPayload(seal AuditSeal) ([]byte, error) {
	version, err := signpayload.ParseVersion(seal.SignPayloadVersion)
	if err != nil {
		return nil, err
	}
	return signpayload.Build(signpayload.Seal{
		Version:                  version,
		Algorithm:                seal.Algorithm,
		FormatVersion:            seal.FormatVersion,
		TenantID:                 seal.TenantID,
		SealDate:                 seal.SealDate,
		MerkleRootHex:            seal.MerkleRootHex,
		HKDFInputsDigest:         seal.HKDFInputsDigestHex,
		Cadence:                  seal.Cadence,
		DevMode:                  seal.DevMode,
		KeyVersionsCanon:         seal.KeyVersionsCanon,
		KMSHandleURIsDigest:      seal.KMSHandleURIsDigestHex,
		OperationalEventsLogRoot: seal.OperationalEventsLogRoot,
	})
}

// SealCarriesLiveSignature reports whether the seal carries a real
// Ed25519 signature the live step-11 walk can drive: a base64 string that
// decodes to exactly 64 bytes. A placeholder (N005's textual
// "TEST-SIGNATURE-PLACEHOLDER…") or an absent signature decodes to the
// wrong length, so the driver keeps such a vector contract-only rather
// than asserting a live verify it cannot perform.
//
// Note: N004's 88-char base64 of 'Z' DECODES to 66 bytes, not 64 — so a
// caller wanting to drive N004's garbage-signature path live should not
// gate on this helper for the negative direction. This helper answers
// "could this seal possibly carry a valid signature", which is the gate
// for the POSITIVE live path; the negative driver (garbage/wrong-tenant)
// has its own gate. See the vectors-package live driver.
func SealCarriesLiveSignature(seal AuditSeal) bool {
	sig, err := base64.StdEncoding.DecodeString(seal.SignatureB64)
	if err != nil {
		return false
	}
	return len(sig) == ed25519.SignatureSize
}

// errInvalidPublicKeyHex is returned by ParsePublicKeyHex for a hex string
// that does not decode to a 32-byte Ed25519 public key.
type errInvalidPublicKeyHex struct{ msg string }

func (e errInvalidPublicKeyHex) Error() string { return e.msg }

// ParsePublicKeyHex decodes a 64-char lowercase-hex Ed25519 public key
// (the corpus-published pub.hex form). It is the verifier's single entry
// point for public-key material — the verifier holds the public key only,
// never the private seed.
func ParsePublicKeyHex(h string) (ed25519.PublicKey, error) {
	raw, err := hex.DecodeString(h)
	if err != nil {
		return nil, errInvalidPublicKeyHex{fmt.Sprintf("public key is not valid hex: %v", err)}
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, errInvalidPublicKeyHex{
			fmt.Sprintf("public key has %d bytes, want %d", len(raw), ed25519.PublicKeySize)}
	}
	return ed25519.PublicKey(raw), nil
}
