package verify

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/mmpworks/ffiec/core/constants"
)

// KeyFingerprint computes the §4.1 key fingerprint:
// SHA-256(utf8(tenant_id) || ikm) truncated to the first
// constants.KeyFingerprintLen (16) bytes. The fingerprint binds a
// session key to its tenant + IKM generation so a botched key rotation
// surfaces at §7 step 8 before any MAC compute.
func KeyFingerprint(tenantID string, ikm []byte) []byte {
	h := sha256.New()
	h.Write([]byte(tenantID))
	h.Write(ikm)
	full := h.Sum(nil)
	return full[:constants.KeyFingerprintLen]
}

// KeyFingerprintHex is the lowercase-hex form of KeyFingerprint, the
// shape the corpus pins.
func KeyFingerprintHex(tenantID string, ikm []byte) string {
	return hex.EncodeToString(KeyFingerprint(tenantID, ikm))
}
