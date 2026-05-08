package cliutil

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// LoadEd25519PKCS8 reads a PEM-encoded PKCS#8 Ed25519 private key from path
// and returns the (private, public) pair. The expected PEM type is
// "PRIVATE KEY" — the format produced by `openssl genpkey -algorithm ed25519`
// and by GenerateEd25519PKCS8.
func LoadEd25519PKCS8(path string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, nil, errors.New("no PEM block found")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("unsupported PEM type %q (expected PRIVATE KEY / PKCS#8)", block.Type)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parse PKCS#8: %w", err)
	}
	priv, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("PKCS#8 key is not Ed25519")
	}
	pub, ok := priv.Public().(ed25519.PublicKey)
	if !ok {
		return nil, nil, errors.New("could not derive Ed25519 public key")
	}
	return priv, pub, nil
}

// LoadEd25519PublicPEM reads a PEM-encoded Ed25519 public key (PKIX format,
// "PUBLIC KEY" PEM type) and returns it. Used by the verifier to load the
// institution's published seal-signing public key without ever touching a
// private key.
func LoadEd25519PublicPEM(path string) (ed25519.PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("unsupported PEM type %q (expected PUBLIC KEY / PKIX)", block.Type)
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX: %w", err)
	}
	pub, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("PKIX key is not Ed25519")
	}
	return pub, nil
}

// GenerateEd25519PKCS8 mints a new Ed25519 keypair and writes the private
// half to path as PEM-encoded PKCS#8 ("PRIVATE KEY"). Returns the public
// half so the caller can print a fingerprint.
func GenerateEd25519PKCS8(path string) (ed25519.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("marshal pkcs8: %w", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		return nil, fmt.Errorf("write %s: %w", path, err)
	}
	return pub, nil
}

// WriteEd25519PublicPEM writes pub to path in PKIX PEM form. Used by the
// ledger server (when it publishes the seal-signing public key) and by
// tests that need a paired key on disk.
func WriteEd25519PublicPEM(path string, pub ed25519.PublicKey) error {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return fmt.Errorf("marshal pkix: %w", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
