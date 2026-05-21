package jcs

import (
	"bytes"
	"fmt"
)

// writeString encodes a Go string as a JSON string per RFC 8785
// §3.2.2.2.
//
// The escape rules (from RFC 8785 §3.2.2.2 + the deliberate
// references to RFC 8259 §7):
//
//   - U+0022 (double quote) → \"
//   - U+005C (backslash) → \\
//   - U+0008, U+0009, U+000A, U+000C, U+000D → short escapes
//     (\b, \t, \n, \f, \r) in that exact order
//   - Every other character in the ASCII control range
//     (U+0000..U+001F) → \u00xx with LOWERCASE hex (Richard's
//     cross-consult flagged this; the stdlib's encoding/json emits
//     UPPERCASE \uXXXX and is therefore non-canonical for JCS)
//   - Every other character (including U+007F DEL — JCS does NOT
//     escape it; encoding/json's safe-mode does — and including all
//     non-ASCII Unicode) → raw UTF-8 bytes
//   - U+002F (forward slash) → emitted as-is, NEVER escaped (RFC
//     8259 permits the escape; JCS forbids unnecessary escapes per
//     vector 008 control_characters category)
//
// Why we don't delegate to json.Marshal: it emits uppercase \uXXXX
// escapes (non-canonical), it escapes non-ASCII characters by
// default (also non-canonical for JCS), and its safe-mode escapes
// U+007F (non-canonical). Three deviations is more than enough to
// justify the ~40-line in-place encoder below.
func writeString(buf *bytes.Buffer, s string) error {
	buf.WriteByte('"')
	// We walk the input string's BYTES (not runes), because raw UTF-8
	// passes through unchanged — only the ASCII range needs inspection
	// for escape rules. A multi-byte UTF-8 sequence's continuation
	// bytes are all 0x80..0xBF, none of which match any escape rule,
	// so they're written through as-is.
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\t':
			buf.WriteString(`\t`)
		case '\n':
			buf.WriteString(`\n`)
		case '\f':
			buf.WriteString(`\f`)
		case '\r':
			buf.WriteString(`\r`)
		default:
			if c < 0x20 {
				// Control character without a short escape — emit as
				// \u00xx with LOWERCASE hex. The format spec "%04x" is
				// lowercase per Go's fmt convention; the explicit width
				// (4) handles the leading zeros the spec mandates.
				fmt.Fprintf(buf, `\u%04x`, c)
			} else {
				// Includes U+007F DEL (the canonical-form-preserves-it
				// case the vector 008 control_characters category pins),
				// every byte >= 0x80 (UTF-8 leading + continuation bytes
				// of any multi-byte sequence), and ordinary ASCII.
				buf.WriteByte(c)
			}
		}
	}
	buf.WriteByte('"')
	return nil
}
