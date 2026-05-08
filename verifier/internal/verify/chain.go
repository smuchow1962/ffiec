package verify

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// CheckChain walks entries in order and confirms two invariants:
//
//  1. Every entry_hash equals SHA-256 of the entry's canonical bytes (with
//     entry_hash zeroed). Tamper anywhere in the entry breaks this.
//  2. Every entry's prev_hash equals the previous entry's entry_hash, and
//     the genesis entry's prev_hash is empty. Tamper or reorder breaks this.
//
// Returns the index of the first bad entry on failure. The error message
// names which invariant tripped so an examiner can localise the break.
func CheckChain(entries []ChainEntry) (badIndex int, err error) {
	prev := ""
	for i := range entries {
		e := &entries[i]
		if err := checkEntryHash(e); err != nil {
			return i, fmt.Errorf("entry %d (%s): %w", i, e.EntryID, err)
		}
		if e.PrevHash != prev {
			return i, fmt.Errorf("entry %d (%s): prev_hash does not match prior entry_hash (got %q, want %q)",
				i, e.EntryID, e.PrevHash, prev)
		}
		prev = e.EntryHash
	}
	return -1, nil
}

func checkEntryHash(e *ChainEntry) error {
	canon, err := canonicalForEntryHash(e)
	if err != nil {
		return err
	}
	want, err := decodeBase64(e.EntryHash, "entry_hash")
	if err != nil {
		return err
	}
	got := sha256.Sum256(canon)
	if !bytesEqual(got[:], want) {
		return fmt.Errorf("entry_hash does not match canonical bytes (recomputed=%s, claimed=%s)",
			base64.StdEncoding.EncodeToString(got[:]), e.EntryHash)
	}
	return nil
}

// decodeBase64 wraps base64.StdEncoding.DecodeString with a field-name
// error context so the verifier's failure messages are localising.
func decodeBase64(s, field string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", field, err)
	}
	return b, nil
}

// bytesEqual is a tiny constant-time-friendly comparison. crypto/subtle is
// not strictly required here (we are not protecting a secret), but we use
// the constant-shape pattern anyway so the code reviews the same.
func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
