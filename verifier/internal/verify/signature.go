package verify

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// CheckSealSignature confirms the seal's Ed25519 signature under pub. The
// signed message is the canonical encoding of the seal record with the
// signature_ed25519 field cleared.
//
// Additionally checks that pub matches the seal's signing_key_fingerprint
// — a mismatch means the verifier was handed the wrong public key, which
// would silently pass signature verification on a different signer.
func CheckSealSignature(s *SealRecord, pub ed25519.PublicKey) error {
	if l := len(pub); l != ed25519.PublicKeySize {
		return fmt.Errorf("seal-signing public key: have %d bytes, want %d", l, ed25519.PublicKeySize)
	}
	if err := checkFingerprint(pub, s.SigningKeyFingerprint); err != nil {
		return err
	}
	msg, err := canonicalForSealSignature(s)
	if err != nil {
		return err
	}
	sig, err := decodeBase64(s.SignatureEd25519, "signature_ed25519")
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, msg, sig) {
		return errors.New("seal signature does not verify under the supplied public key")
	}
	return nil
}

// checkFingerprint compares SHA-256 of pub bytes to the seal's stated
// signing_key_fingerprint. A mismatch is the loudest signal that the
// verifier has been pointed at a different signer than the institution
// claims — fail closed.
func checkFingerprint(pub ed25519.PublicKey, claimed string) error {
	if claimed == "" {
		return errors.New("seal has empty signing_key_fingerprint")
	}
	want := strings.TrimPrefix(claimed, "sha256:")
	if want == claimed {
		return fmt.Errorf("signing_key_fingerprint missing 'sha256:' prefix: %q", claimed)
	}
	sum := sha256.Sum256(pub)
	got := hex.EncodeToString(sum[:])
	if got != want {
		return fmt.Errorf("signing_key_fingerprint does not match supplied public key (computed=%s, claimed=%s)", got, want)
	}
	return nil
}
