package verify

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// CheckEntryMAC re-derives the per-event MAC key from the master IKM and
// the entry's tenant_binding_kdf_label + entry_id, then checks the entry's
// hmac_sha256 over the event_payload_jcs bytes.
//
// This is the heavier verification path the operator-guide calls "full
// per-event integrity" — it requires the institution's master IKM. Without
// it, the verifier is structural-only (chain linkage + Merkle + signature).
func CheckEntryMAC(e *ChainEntry, ikm []byte) error {
	payload, err := decodeBase64(e.EventPayloadJCS, "event_payload_jcs")
	if err != nil {
		return err
	}
	wantMAC, err := decodeBase64(e.HMACSHA256, "hmac_sha256")
	if err != nil {
		return err
	}
	key := hkdf32(ikm, []byte(e.TenantBindingKDFLabel), []byte(e.EntryID))
	gotMAC := hmacSHA256(key, payload)
	if !hmac.Equal(gotMAC, wantMAC) {
		return fmt.Errorf("entry %s: per-event MAC does not match", e.EntryID)
	}
	return nil
}

// CheckChainMACs runs CheckEntryMAC across every entry, returning the
// index of the first failing entry on error.
func CheckChainMACs(entries []ChainEntry, ikm []byte) (badIndex int, err error) {
	if len(ikm) == 0 {
		return -1, fmt.Errorf("master IKM is empty")
	}
	for i := range entries {
		if err := CheckEntryMAC(&entries[i], ikm); err != nil {
			return i, err
		}
	}
	return -1, nil
}

// hmacSHA256 returns HMAC-SHA-256(key, msg).
func hmacSHA256(key, msg []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return h.Sum(nil)
}

// hkdf32 is HKDF-SHA-256 truncated to 32 bytes (RFC 5869). The salt is
// the tenant_binding_kdf_label and the info is the entry_id, so each
// per-event key is bound to both the tenant context and the specific event.
//
// We implement HKDF inline because the verifier ships as a single static
// stdlib-only binary; the construction is small and exercise-tested below.
func hkdf32(ikm, salt, info []byte) []byte {
	prk := hkdfExtract(salt, ikm)
	return hkdfExpand(prk, info, 32)
}

func hkdfExtract(salt, ikm []byte) []byte {
	if len(salt) == 0 {
		// Per RFC 5869 §2.2, an empty salt is replaced by HashLen zeros.
		salt = make([]byte, sha256.Size)
	}
	return hmacSHA256(salt, ikm)
}

func hkdfExpand(prk, info []byte, length int) []byte {
	if length <= 0 {
		return nil
	}
	out := make([]byte, 0, length)
	var t []byte
	for i := byte(1); len(out) < length; i++ {
		h := hmac.New(sha256.New, prk)
		h.Write(t)
		h.Write(info)
		h.Write([]byte{i})
		t = h.Sum(nil)
		out = append(out, t...)
	}
	return out[:length]
}
