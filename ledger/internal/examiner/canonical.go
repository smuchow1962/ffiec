package examiner

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// canonicalForSigning returns the deterministic JSON encoding of the bundle
// used as input to the Ed25519 signature. The signature_ed25519 field is
// excluded from the input — a signature cannot cover its own bytes.
//
// This canonicalization is sufficient for the v1.0 bundle schema:
//
//   1. Encoding is via encoding/json on a fixed struct, so field order is
//      stable across Go versions and matches the source declaration order.
//   2. All fields are string, time.Time, or string-slice — no maps, no
//      floats. The JSON encoder is byte-stable for this shape.
//   3. signature_ed25519 carries `omitempty`; we clear it before encoding
//      so the signing input is identical on issuer and verifier.
//
// Upgrade note (cross-language conformance): once non-Go consumers verify
// bundles directly (rather than going through verifier/), this function
// must be replaced with a full RFC 8785 (JCS) implementation. JCS is the
// normative canonicalization used elsewhere in the chain spec; the bundle
// has stayed Go-only during bootstrap, so struct-order encoding suffices.
func canonicalForSigning(b *Bundle) ([]byte, error) {
	if b == nil {
		return nil, fmt.Errorf("nil bundle")
	}
	cp := *b
	cp.SignatureEd25519 = ""

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&cp); err != nil {
		return nil, fmt.Errorf("canonical encode: %w", err)
	}
	// json.Encoder.Encode appends a trailing newline; drop it so the bytes
	// are exactly the document.
	out := buf.Bytes()
	if n := len(out); n > 0 && out[n-1] == '\n' {
		out = out[:n-1]
	}
	return out, nil
}
