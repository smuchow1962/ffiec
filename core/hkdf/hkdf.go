// Package hkdf provides HKDF-SHA-256 (RFC 5869) extracted from the
// verifier's inline implementation so the test-vector consumer and
// future per-event MAC code can share one source of truth.
//
// Spec authority: §4.1 normates HKDF-SHA-256 as the per-tenant
// session-key derivation function with a fixed 32-byte output
// length. The construction is RFC 5869 §2.2 + §2.3 — extract over
// (salt, IKM) then expand over (PRK, info) with a length parameter.
//
// Why a small in-repo package instead of x/crypto/hkdf:
//
//   - The verifier ships as a stdlib-only static binary per the
//     existing repo posture (see verifier/internal/verify/mac.go's
//     comment: "we implement HKDF inline because the verifier ships
//     as a single static stdlib-only binary").
//   - HKDF is ~30 lines; the dependency-surface cost of pulling in
//     x/crypto would dwarf the implementation cost.
//   - Two callers now (vectors runner + verify pipeline) need the
//     same function; extracting prevents duplication.
package hkdf

import (
	"crypto/hmac"
	"crypto/sha256"
)

// Derive returns HKDF-SHA-256(IKM=ikm, salt=salt, info=info, L=length).
// Per RFC 5869 §2.2, an empty salt is replaced with HashLen (32)
// zero bytes — the helper applies that rule so callers don't need
// to special-case it.
//
// A zero or negative length returns nil.
func Derive(ikm, salt, info []byte, length int) []byte {
	if length <= 0 {
		return nil
	}
	prk := extract(salt, ikm)
	return expand(prk, info, length)
}

func extract(salt, ikm []byte) []byte {
	if len(salt) == 0 {
		salt = make([]byte, sha256.Size)
	}
	return hmacSHA256(salt, ikm)
}

// expand implements HKDF-Expand per RFC 5869 §2.3. The output is the
// concatenation T(1) || T(2) || ... truncated to length bytes,
// where T(i) = HMAC(PRK, T(i-1) || info || i) and T(0) is empty.
func expand(prk, info []byte, length int) []byte {
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

func hmacSHA256(key, msg []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return h.Sum(nil)
}
