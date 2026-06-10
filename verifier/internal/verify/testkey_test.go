package verify

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TesseraSeal TEST signing-key loader for the verify package's signature
// tests.
//
// The verifier itself NEVER reads the private seed — it consumes only
// public-key material (CheckSealSignatureV1 takes an ed25519.PublicKey).
// These helpers exist for TESTS that must SIGN a fixture to drive the
// positive step-11 path. They read the seed from TESSERASEAL_TEST_KEY_DIR
// (default E:\dev\testing\private-keys\tesseraseal) and SKIP with a clear
// message when the seed is absent, so CI without the key still runs the
// rest of the suite green.

const (
	// testKeyDirEnvVar overrides the default test-key directory.
	testKeyDirEnvVar = "TESSERASEAL_TEST_KEY_DIR"
	// defaultTestKeyDir is the local-only key directory the corpus
	// generators and these tests read by convention.
	defaultTestKeyDir = `E:\dev\testing\private-keys\tesseraseal`
	// seedHexFile holds the 32-byte private seed as hex (LOCAL ONLY).
	seedHexFile = "test-signing-key.seed.hex"
	// pubHexFile holds the 32-byte public key as hex (publishable).
	pubHexFile = "test-signing-key.pub.hex"
	// publishedPubHex is the corpus's published public key. A loaded seed
	// whose derived public key differs from this is a key-rotation error,
	// not a silent test pass.
	publishedPubHex = "0985603b6c0e099bac783bcd7801664ed376a7888eafbf191859854bf8ff7f35"
)

// testKeyDir resolves the test-key directory from the env var with the
// local-only default.
func testKeyDir() string {
	if v := os.Getenv(testKeyDirEnvVar); v != "" {
		return v
	}
	return defaultTestKeyDir
}

// loadTestSigningKey reads the private seed and returns the signing key.
// It SKIPS the calling test (never fails it) when the seed file is absent,
// so a CI machine without the local-only key still passes the rest of the
// gate. When the seed IS present, a derived-public-key mismatch against
// the published pub.hex is a hard failure — the key was rotated without
// regenerating the corpus.
func loadTestSigningKey(t *testing.T) ed25519.PrivateKey {
	t.Helper()
	dir := testKeyDir()
	seedPath := filepath.Join(dir, seedHexFile)
	raw, err := os.ReadFile(seedPath)
	if err != nil {
		t.Skipf("test signing key not available at %s (set %s to enable signature tests): %v",
			seedPath, testKeyDirEnvVar, err)
	}
	seed, err := hex.DecodeString(strings.TrimSpace(string(raw)))
	if err != nil {
		t.Fatalf("decode seed hex %s: %v", seedPath, err)
	}
	if len(seed) != ed25519.SeedSize {
		t.Fatalf("seed %s has %d bytes, want %d", seedPath, len(seed), ed25519.SeedSize)
	}
	priv := ed25519.NewKeyFromSeed(seed)
	got := hex.EncodeToString(priv.Public().(ed25519.PublicKey))
	if got != publishedPubHex {
		t.Fatalf("loaded seed derives public key %s, but the corpus publishes %s — the test key was rotated without regenerating the corpus",
			got, publishedPubHex)
	}
	return priv
}

// publishedTestPublicKey returns the corpus-published public key, parsed
// through the verifier's own ParsePublicKeyHex entry point. It needs no
// private material, so it never skips.
func publishedTestPublicKey(t *testing.T) ed25519.PublicKey {
	t.Helper()
	pub, err := ParsePublicKeyHex(publishedPubHex)
	if err != nil {
		t.Fatalf("parse published public key: %v", err)
	}
	return pub
}
