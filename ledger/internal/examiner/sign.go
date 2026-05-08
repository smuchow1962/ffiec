package examiner

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

// Sign returns the signed JSON bytes of b. priv is the institution's
// issuance private key. b is mutated to carry the resulting signature so the
// caller can inspect it.
func Sign(b *Bundle, priv ed25519.PrivateKey) ([]byte, error) {
	if l := len(priv); l != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("issuance private key: have %d bytes, want %d", l, ed25519.PrivateKeySize)
	}
	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("bundle invariants: %w", err)
	}
	msg, err := canonicalForSigning(b)
	if err != nil {
		return nil, err
	}
	sig := ed25519.Sign(priv, msg)
	b.SignatureEd25519 = base64.StdEncoding.EncodeToString(sig)
	out, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal signed bundle: %w", err)
	}
	return out, nil
}

// Verify parses bundleJSON and confirms the signature under pub. On success
// it returns the parsed Bundle. The caller is still responsible for
// enforcing time-window invariants (expires_at, examination period) — Verify
// confirms cryptographic and structural integrity only.
func Verify(bundleJSON []byte, pub ed25519.PublicKey) (*Bundle, error) {
	if l := len(pub); l != ed25519.PublicKeySize {
		return nil, fmt.Errorf("issuance public key: have %d bytes, want %d", l, ed25519.PublicKeySize)
	}
	var b Bundle
	if err := json.Unmarshal(bundleJSON, &b); err != nil {
		return nil, fmt.Errorf("parse bundle: %w", err)
	}
	if b.SignatureEd25519 == "" {
		return nil, errors.New("bundle has no signature")
	}
	sig, err := base64.StdEncoding.DecodeString(b.SignatureEd25519)
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}
	msg, err := canonicalForSigning(&b)
	if err != nil {
		return nil, err
	}
	if !ed25519.Verify(pub, msg, sig) {
		return nil, errors.New("signature does not verify under the supplied public key")
	}
	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("bundle invariants: %w", err)
	}
	return &b, nil
}
