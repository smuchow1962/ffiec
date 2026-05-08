package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// canonicalForEntryHash returns the bytes hashed to produce ChainEntry.EntryHash.
// The entry_hash field is cleared on the copy so the document signs/hashes
// itself without containing its own hash.
//
// Same struct-order assumption as the ledgerctl bundle canonicalization:
// stable for a fixed Go schema; production cross-language consumers must
// upgrade to RFC 8785 (JCS) before forking implementations.
func canonicalForEntryHash(e *ChainEntry) ([]byte, error) {
	if e == nil {
		return nil, fmt.Errorf("nil entry")
	}
	cp := *e
	cp.EntryHash = ""
	return marshalNoNewline(&cp)
}

// canonicalForSealSignature returns the bytes signed to produce
// SealRecord.SignatureEd25519. The signature_ed25519 field is cleared.
func canonicalForSealSignature(s *SealRecord) ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("nil seal")
	}
	cp := *s
	cp.SignatureEd25519 = ""
	return marshalNoNewline(&cp)
}

// marshalNoNewline encodes v with json.Encoder (which avoids HTML-escaping
// when SetEscapeHTML(false)) and strips the trailing newline so the bytes
// are exactly the document.
func marshalNoNewline(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("canonical encode: %w", err)
	}
	out := buf.Bytes()
	if n := len(out); n > 0 && out[n-1] == '\n' {
		out = out[:n-1]
	}
	return out, nil
}
