// Package jcs implements RFC 8785 (JSON Canonicalization Scheme) —
// the byte-exact canonical encoding the spec's per-event MAC, daily
// seal, and verdict-object compute over.
//
// Spec authority: RFC 8785; FFIEC chain-of-custody spec §5 (canonical
// bytes are the MAC input) + §4.1 (the MAC covers prev_hash ||
// canonical_bytes) + §4.2 (canonical bytes are the Merkle leaf
// preimage) + §10.12 (the verdict object is JCS-canonical).
//
// Conformance bar: test-vector 008-jcs-edge-cases. An implementation
// is v1.0-conformant only if it reproduces every byte sequence the
// vector pins AND rejects every NaN / Infinity input with an error
// (per §3.2.2.3 of RFC 8785).
//
// Implementation notes:
//
//   - Input is interface{} (the Go decoded form of any JSON value).
//     Callers parsing JSON should use encoding/json to decode into
//     interface{}, then call Canonicalize. The decode path inherits
//     encoding/json's int64-vs-float64 disposition (every JSON
//     number becomes float64 unless Decoder.UseNumber is set —
//     callers needing 64-bit-integer precision MUST set UseNumber
//     and the canonicalizer handles json.Number).
//   - Stdlib-only per decision D-1. No x/text/collate, no
//     gowebpki/jcs. ~250 lines across this file + number.go +
//     string.go.
//   - The §7 pre-flight self-test lives in selftest.go and is what
//     the verifier runs at startup to confirm the JCS implementation
//     wasn't replaced at build time with a non-conformant one.
package jcs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"unicode/utf16"
)

// Canonicalize returns the RFC 8785 canonical UTF-8 bytes of v.
//
// Supported Go types (mirroring what encoding/json produces):
//
//   - nil → "null"
//   - bool → "true" / "false"
//   - string → RFC 8785 §3.2.2.2 string escape, UTF-8 emitted as raw
//     bytes, control chars and double-quote / backslash escaped
//   - float64 → ECMAScript Number.prototype.toString form
//   - json.Number → preserves int-form when the value parses as
//     int64; otherwise serializes as ECMAScript Number form via
//     float64 round-trip (per JCS §3.2.2.3 IEEE-754-double mandate)
//   - int / int8..int64 / uint / uint8..uint64 — convenience types
//     so callers passing typed integer values don't have to wrap
//     them; serialized as decimal integer
//   - []interface{} → "[" + canonical(elem) + "," + ... + "]"
//   - map[string]interface{} → "{" + key + ":" + canonical(value)
//     + "," + ... + "}" with keys sorted by UTF-16 code-unit order
//     per §3.2.3
//
// Returns an error for:
//
//   - NaN, +Inf, -Inf float64 (per §3.2.2.3 — these have no canonical
//     JSON form)
//   - Unsupported Go types (e.g., channels, functions, structs that
//     callers should marshal through encoding/json first)
//   - Map keys that are not strings (Go's interface{}-typed JSON
//     decode produces map[string]interface{}, so this is normally
//     unreachable; surfaces only if a caller hand-builds a
//     map[interface{}]interface{})
func Canonicalize(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := writeValue(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// CanonicalizeJSON is a convenience wrapper for the common case:
// raw JSON bytes in, canonical JSON bytes out. Decodes with
// json.Decoder + UseNumber so int64 precision is preserved through
// the round-trip (per JCS §3.2.2.3 numeric edge cases — values in
// the [-(2^53-1), 2^53-1] safe-integer range serialize exactly).
func CanonicalizeJSON(raw []byte) ([]byte, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("jcs: decode input: %w", err)
	}
	return Canonicalize(v)
}

// writeValue dispatches one Go value to its canonical encoding. The
// switch covers every type Canonicalize's contract names; the
// default branch rejects unsupported types with a localising error.
func writeValue(buf *bytes.Buffer, v interface{}) error {
	switch x := v.(type) {
	case nil:
		buf.WriteString("null")
		return nil
	case bool:
		if x {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
		return nil
	case string:
		return writeString(buf, x)
	case json.Number:
		return writeJSONNumber(buf, x)
	case float64:
		return writeFloat64(buf, x)
	case float32:
		return writeFloat64(buf, float64(x))
	case int:
		return writeInt64(buf, int64(x))
	case int8:
		return writeInt64(buf, int64(x))
	case int16:
		return writeInt64(buf, int64(x))
	case int32:
		return writeInt64(buf, int64(x))
	case int64:
		return writeInt64(buf, x)
	case uint:
		return writeUint64(buf, uint64(x))
	case uint8:
		return writeUint64(buf, uint64(x))
	case uint16:
		return writeUint64(buf, uint64(x))
	case uint32:
		return writeUint64(buf, uint64(x))
	case uint64:
		return writeUint64(buf, x)
	case []interface{}:
		return writeArray(buf, x)
	case map[string]interface{}:
		return writeObject(buf, x)
	}
	return fmt.Errorf("jcs: unsupported type %T (decode JSON to interface{} with json.Decoder.UseNumber and try again)", v)
}

// writeArray encodes a JSON array. Element order is preserved (per
// RFC 8785 §3.2.3 — array element order MUST NOT be changed); the
// canonicalizer only sorts object keys.
func writeArray(buf *bytes.Buffer, a []interface{}) error {
	buf.WriteByte('[')
	for i, elem := range a {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := writeValue(buf, elem); err != nil {
			return err
		}
	}
	buf.WriteByte(']')
	return nil
}

// writeObject encodes a JSON object with keys sorted per RFC 8785
// §3.2.3 — by UTF-16 code-unit lexicographic comparison. The
// implementation precomputes each key's UTF-16 representation once
// and reuses it both for the comparator and as the original-string
// reference so the output emits the raw UTF-8 (not a re-encoding)
// per §3.2.2.2.
//
// Cognitive complexity note: the sort comparator inspects two
// UTF-16 sequences as unsigned 16-bit code units. .NET's
// StringComparer.Ordinal walks chars (which ARE UTF-16 code units in
// C#); Python's sorted() compares Unicode code points; Go's
// sort.Strings compares UTF-8 bytes. None of those match RFC 8785
// for non-BMP strings — we implement the UTF-16 walk explicitly so a
// future supplementary-plane key (the case the vector 008
// `unicode_keys_with_supplementary` test pins) sorts correctly.
func writeObject(buf *bytes.Buffer, m map[string]interface{}) error {
	keys := sortedKeys(m)
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		if err := writeString(buf, k); err != nil {
			return err
		}
		buf.WriteByte(':')
		if err := writeValue(buf, m[k]); err != nil {
			return err
		}
	}
	buf.WriteByte('}')
	return nil
}

// sortedKeys returns m's keys ordered per RFC 8785 §3.2.3 (UTF-16
// code-unit lexicographic). The sort is stable across runs because
// the comparator is total over the (de-duplicated) key set.
//
// JSON parsing eliminates duplicate keys upstream (encoding/json
// keeps the last value), so the map already has unique keys by the
// time it reaches here.
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return lessUTF16(keys[i], keys[j])
	})
	return keys
}

// lessUTF16 reports whether a < b under RFC 8785 §3.2.3 UTF-16
// code-unit lexicographic ordering. Implemented as two parallel
// UTF-16 encodings + a code-unit-by-code-unit comparison.
//
// Why not a precomputed-UTF-16 cache: the sort comparator runs
// O(N log N) times; for typical chain entries N <= 20 keys and the
// cost of the repeated UTF-16 encode is dwarfed by the JSON I/O
// either side. If a future profile shows this as a hot spot, the
// sort can switch to a (key, []uint16) tuple with the UTF-16 pre-
// encoded once per key.
func lessUTF16(a, b string) bool {
	au := utf16.Encode([]rune(a))
	bu := utf16.Encode([]rune(b))
	n := len(au)
	if len(bu) < n {
		n = len(bu)
	}
	for i := 0; i < n; i++ {
		if au[i] != bu[i] {
			return au[i] < bu[i]
		}
	}
	// Equal prefix; shorter string sorts first.
	return len(au) < len(bu)
}

// ErrNonFinite is returned by Canonicalize when a NaN or Infinity
// float64 reaches the encoder. Per RFC 8785 §3.2.2.3 these have no
// canonical JSON form and conforming implementations MUST reject
// them rather than coerce to null or string. Errors returned from
// the canonicalizer wrap ErrNonFinite so callers can branch on
// errors.Is.
var ErrNonFinite = errors.New("jcs: NaN or Infinity is not representable in canonical JSON")
