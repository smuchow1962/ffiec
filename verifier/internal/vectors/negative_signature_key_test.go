package vectors

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"path/filepath"
	"testing"
)

// TestCorpusPublicKey_VerifiesAGenuineSignature pins the key file itself. N004 and N005
// are negative, so their step-11 walks fail under any key, and a wrong key file would pass
// them. N005 holds a genuine signature over its published payload, so this key must verify
// it, and a changed key must not.
func TestCorpusPublicKey_VerifiesAGenuineSignature(t *testing.T) {
	corpus, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("corpus not available: %v", err)
	}
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(corpus, "negative", "N005-signature-wrong-tenant")
	af, _, err := loadAuditFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := resolveCorpusPublicKey(dir)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := hex.DecodeString(af.Seal.SignPayloadHex)
	if err != nil {
		t.Fatal(err)
	}
	signature, err := base64.StdEncoding.DecodeString(af.Seal.SignatureB64)
	if err != nil {
		t.Fatal(err)
	}
	if !ed25519.Verify(pub, payload, signature) {
		t.Fatal("the corpus public key does not verify N005's genuine signature")
	}
	wrong := append(ed25519.PublicKey(nil), pub...)
	wrong[0] ^= 1
	if ed25519.Verify(wrong, payload, signature) {
		t.Fatal("a changed key verified N005's signature")
	}
}
