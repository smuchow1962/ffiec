package jcs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// MaxDepth bounds the nesting the canonicalizer descends before
// rejecting the input. It matches the .NET reference's
// CanonicalJson.MaxDepth (256) so the two implementations reject the
// same deeply-nested-object attack class at the same boundary.
const MaxDepth = 256

// Canonicalize serializes value to RFC 8785 (JSON Canonicalization
// Scheme) UTF-8 bytes. The output is deterministic: identical inputs
// produce byte-identical outputs across runs and across the C# /
// Python / Go implementations.
//
// value is the Go decoding of a JSON document. Accepted dynamic types
// mirror what encoding/json produces from json.Unmarshal into `any`,
// plus the integer types a typed caller may pass:
//
//   - nil                      → null
//   - bool                     → true / false
//   - string                   → JSON string (minimal escape set)
//   - float64                  → ECMA-262 NumberToString
//   - int, int64               → decimal integer (no float round-trip)
//   - json.Number              → preserved digit string when integral
//   - map[string]any           → object, keys byte-ordinal sorted
//   - []any                    → array, element order preserved
//
// A non-finite float (NaN, ±Inf), a lone UTF-16 surrogate, an
// unsupported value type, or nesting past MaxDepth is a hard error —
// the same conditions the .NET reference rejects with ArgumentException.
//
// Why not encoding/json: the stdlib emits UPPERCASE \uXXXX escapes,
// HTML-escapes <>&, and routes every number through float64. All three
// are silent byte-divergences against the FFIEC corpus. We hand-roll
// the bytes (per D-1, locked 2026-05-21) to keep full control.
func Canonicalize(value any) ([]byte, error) {
	var b strings.Builder
	if err := writeValue(&b, value, 0); err != nil {
		return nil, err
	}
	return []byte(b.String()), nil
}

// writeValue dispatches one value to the right writer. The switch
// order is load-bearing only insofar as the concrete types are
// disjoint; map and slice are matched explicitly so an unsupported
// type falls through to the error arm rather than being mis-encoded.
func writeValue(b *strings.Builder, value any, depth int) error {
	switch v := value.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if v {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case string:
		return writeString(b, v)
	case float64:
		return writeNumber(b, v)
	case float32:
		return writeNumber(b, float64(v))
	case int:
		b.WriteString(strconv.FormatInt(int64(v), 10))
	case int64:
		b.WriteString(strconv.FormatInt(v, 10))
	case map[string]any:
		return writeObject(b, v, depth)
	case []any:
		return writeArray(b, v, depth)
	default:
		// json.Number arrives as a named string type; handle it by
		// reflection-free type name so callers decoding with
		// UseNumber() get integral preservation without a hard
		// dependency on encoding/json in this package.
		if n, ok := asJSONNumber(value); ok {
			return writeJSONNumber(b, n)
		}
		return fmt.Errorf("jcs: unsupported value type %T", value)
	}
	return nil
}

// writeObject emits a JSON object with keys sorted by byte-ordinal
// (UTF-8) order. Go's sort.Strings compares byte-by-byte, which is the
// same ordering as .NET's StringComparer.Ordinal over BMP keys and
// Python's sorted() — the three references agree because audit-chain
// keys are ASCII.
func writeObject(b *strings.Builder, m map[string]any, depth int) error {
	if depth >= MaxDepth {
		return fmt.Errorf("jcs: nesting exceeds MaxDepth (%d)", MaxDepth)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		if err := writeString(b, k); err != nil {
			return err
		}
		b.WriteByte(':')
		if err := writeValue(b, m[k], depth+1); err != nil {
			return err
		}
	}
	b.WriteByte('}')
	return nil
}

// writeArray emits a JSON array preserving element order (arrays are
// ordered in JSON; JCS does not reorder them).
func writeArray(b *strings.Builder, list []any, depth int) error {
	if depth >= MaxDepth {
		return fmt.Errorf("jcs: nesting exceeds MaxDepth (%d)", MaxDepth)
	}
	b.WriteByte('[')
	for i, item := range list {
		if i > 0 {
			b.WriteByte(',')
		}
		if err := writeValue(b, item, depth+1); err != nil {
			return err
		}
	}
	b.WriteByte(']')
	return nil
}

// writeString writes s as a JSON string with the minimal RFC 8785
// escape set: only " \ and the named short escapes \b \t \n \f \r,
// plus \u00xx (LOWERCASE hex) for the remaining C0 controls
// (U+0000..U+001F). Every other rune — forward slash, U+007F (DEL),
// and all non-ASCII — passes through as raw UTF-8 bytes.
//
// Go strings are UTF-8, so a well-formed string needs no surrogate
// handling: range yields runes directly. A lone surrogate would have
// decoded to utf8.RuneError (U+FFFD) on range, which is a lossy
// pass-through; we reject it explicitly so the Go output cannot
// silently diverge from the C# reference, which rejects lone
// surrogates loud at write time.
func writeString(b *strings.Builder, s string) error {
	b.WriteByte('"')
	for i, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\b':
			b.WriteString(`\b`)
		case '\t':
			b.WriteString(`\t`)
		case '\n':
			b.WriteString(`\n`)
		case '\f':
			b.WriteString(`\f`)
		case '\r':
			b.WriteString(`\r`)
		default:
			if r == 0xFFFD && !isValidUTF8Replacement(s, i) {
				return fmt.Errorf("jcs: lone surrogate / invalid UTF-8 at string index %d", i)
			}
			if r < 0x20 {
				// Remaining C0 control: \u00xx with lowercase hex.
				b.WriteString(`\u00`)
				b.WriteByte(hexLower(byte(r>>4) & 0xF))
				b.WriteByte(hexLower(byte(r) & 0xF))
			} else {
				// ASCII or any higher rune: emit the rune's UTF-8
				// bytes verbatim. WriteRune re-encodes the decoded
				// rune to the same bytes range read it from.
				b.WriteRune(r)
			}
		}
	}
	b.WriteByte('"')
	return nil
}

// isValidUTF8Replacement reports whether the U+FFFD at byte index i is
// a genuine replacement character in the source (3 bytes EF BF BD)
// rather than the decoder's substitution for invalid bytes. A genuine
// U+FFFD is conformant input; a substituted one means the source held
// invalid UTF-8 (e.g. a lone surrogate), which we reject.
func isValidUTF8Replacement(s string, i int) bool {
	return i+3 <= len(s) && s[i] == 0xEF && s[i+1] == 0xBF && s[i+2] == 0xBD
}

func hexLower(nybble byte) byte {
	if nybble < 10 {
		return '0' + nybble
	}
	return 'a' + (nybble - 10)
}
