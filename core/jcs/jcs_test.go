package jcs

import (
	"errors"
	"math"
	"testing"
)

// TestCanonicalize_NullAndBoolean covers the three trivial literal
// cases. These never break in practice but a regression here would
// surface drift in the dispatch switch, so the test is fast +
// self-documenting.
func TestCanonicalize_NullAndBoolean(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{nil, "null"},
		{true, "true"},
		{false, "false"},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("Canonicalize(%v): unexpected error %v", tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("Canonicalize(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_Integers covers the integer-typed convenience
// path. JCS treats all numbers as IEEE-754-double per §3.2.2.3,
// so the byte output here matches the integer-decimal form vector
// 008 numeric_edge_cases pins.
func TestCanonicalize_Integers(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{int64(9007199254740991), "9007199254740991"},   // max_safe_integer
		{int64(-9007199254740991), "-9007199254740991"}, // min_safe_integer
		{int64(9007199254740992), "9007199254740992"},   // above_safe_integer
		{uint64(1), "1"},
		{int32(42), "42"},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("Canonicalize(%v): %v", tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("Canonicalize(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_Floats pins the ECMAScript Number.prototype.toString
// behaviour against the exact byte pins vector 008
// float_canonicalization will require. Each case here mirrors one of
// the named sub-cases in the vector.
//
// A failure here means writeFloat64 is non-conformant; running the
// vector 008 conformance test before fixing the unit tests would
// surface the same failure ten times slower.
func TestCanonicalize_Floats(t *testing.T) {
	cases := []struct {
		name string
		in   float64
		want string
	}{
		{"integer_value_one", 1.0, "1"},
		{"small_decimal", 0.1, "0.1"},
		{"decimal_one_third_approx", 0.3333333333333333, "0.3333333333333333"},
		{"ten_billion", 10000000000.0, "10000000000"},
		{"exponent_threshold_high", 1e21, "1e+21"},
		{"below_exponent_threshold", 1e20, "100000000000000000000"},
		{"exponent_threshold_low", 1e-6, "0.000001"},
		{"below_low_threshold", 1e-7, "1e-7"},
		{"very_small_finite", 5e-324, "5e-324"},
		{"very_large_finite", 1.7976931348623157e+308, "1.7976931348623157e+308"},
		{"positive_zero", 0.0, "0"},
		{"negative_zero", math.Copysign(0, -1), "0"},
		{"trailing_zero_strip", 1.5, "1.5"},
		{"long_mantissa_e_minus_100", 1.234567890123e-100, "1.234567890123e-100"},
		{"small_negative_decimal", -0.0001, "-0.0001"},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("%s: Canonicalize(%v) error %v", tc.name, tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("%s: Canonicalize(%v) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_NonFinite confirms NaN and Infinity reject per
// RFC 8785 §3.2.2.3 + vector 008 nan_infinity_rejection. errors.Is
// against ErrNonFinite is the documented branch surface for callers.
func TestCanonicalize_NonFinite(t *testing.T) {
	cases := []struct {
		name string
		in   float64
	}{
		{"NaN", math.NaN()},
		{"PositiveInfinity", math.Inf(1)},
		{"NegativeInfinity", math.Inf(-1)},
	}
	for _, tc := range cases {
		_, err := Canonicalize(tc.in)
		if err == nil {
			t.Errorf("%s: Canonicalize(%v) returned no error; want ErrNonFinite", tc.name, tc.in)
			continue
		}
		if !errors.Is(err, ErrNonFinite) {
			t.Errorf("%s: error %v is not ErrNonFinite", tc.name, err)
		}
	}
}

// TestCanonicalize_StringBasic exercises the easy string cases that
// don't require the escape table. Non-ASCII is verified separately
// in TestCanonicalize_StringUnicode below.
func TestCanonicalize_StringBasic(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "\"\""},
		{"a", "\"a\""},
		{"hello", "\"hello\""},
		{"a/b", "\"a/b\""}, // forward slash NOT escaped
		{"line1\nline2", "\"line1\\nline2\""},
		{"tab\there", "\"tab\\there\""},
		{"quoted \"value\"", "\"quoted \\\"value\\\"\""},
		{"back\\slash", "\"back\\\\slash\""},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("Canonicalize(%q): %v", tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("Canonicalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_StringControlCharacters covers every distinct
// behaviour of the ASCII control range (U+0000..U+001F): the five
// short-escape characters (\b \t \n \f \r), one hex-escape character
// (U+000B, which has no short form), the boundary at U+001F, the
// null byte (U+0000), and the U+007F DEL byte that JCS preserves raw.
// The control-character category in vector 008 pins these byte-for-
// byte.
//
// All `want` strings are written as INTERPRETED Go string literals
// (double-quoted with backslash escapes) so the source file contains
// no literal control bytes — Go's compiler hard-errors on a literal
// NUL byte in source code, so any want-value containing the escape
// sequence `\u00xx` (six visible characters) must be written using
// backslash-escapes in the Go source.
func TestCanonicalize_StringControlCharacters(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		// Want: literal quote, backslash, u, 0, 0, 0, 0, quote.
		{"null_byte", "\x00", "\"\\u0000\""},
		// Want: literal quote, backslash, u, 0, 0, 1, f, quote.
		{"boundary_1f", "\x1f", "\"\\u001f\""},
		// DEL is NOT escaped per JCS — passes through as raw byte.
		{"del_not_escaped", "\x7f", "\"\x7f\""},
		// Short escapes.
		{"backspace_short_escape", "\b", "\"\\b\""},
		{"tab_short_escape", "\t", "\"\\t\""},
		{"newline_short_escape", "\n", "\"\\n\""},
		{"formfeed_short_escape", "\f", "\"\\f\""},
		{"carriage_return_short_escape", "\r", "\"\\r\""},
		// Vertical tab U+000B has no short escape; hex form.
		{"vertical_tab_uses_hex", "\x0b", "\"\\u000b\""},
		// Confirm lowercase hex (RFC 8785 mandates lowercase; Go's
		// %04x is lowercase by convention).
		{"lowercase_hex_check", "\x10", "\"\\u0010\""},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("%s: Canonicalize(%q): %v", tc.name, tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("%s: Canonicalize(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_StringUnicode confirms non-ASCII bytes pass
// through as raw UTF-8 rather than being escaped as \uXXXX. This is
// the JCS rule that diverges from encoding/json's default behaviour
// (json.Marshal escapes <, >, & by default and can escape non-ASCII
// when SetEscapeHTML is on); the vector 008 non_ascii_unicode
// category pins these.
func TestCanonicalize_StringUnicode(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"latin_1_with_accent", "café", "\"café\""},
		{"cjk_ideographs", "日本語", "\"日本語\""},
		{"emoji_supplementary_plane", "🔒", "\"🔒\""},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("%s: Canonicalize(%q): %v", tc.name, tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("%s: Canonicalize(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestCanonicalize_ObjectKeyOrdering covers the ASCII mixed-case
// case from vector 008's object_key_ordering category. The case
// confirms digits sort before uppercase, uppercase before
// underscore, underscore before lowercase — the natural UTF-16
// code-unit ordering, which for ASCII is identical to byte order.
func TestCanonicalize_ObjectKeyOrdering(t *testing.T) {
	in := map[string]interface{}{
		"b": int64(1),
		"B": int64(2),
		"a": int64(3),
		"A": int64(4),
		"0": int64(5),
		"_": int64(6),
	}
	want := "{\"0\":5,\"A\":4,\"B\":2,\"_\":6,\"a\":3,\"b\":1}"
	got, err := Canonicalize(in)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if string(got) != want {
		t.Errorf("Canonicalize = %q, want %q", got, want)
	}
}

// TestCanonicalize_NumericStringKeysAreLexicographic confirms keys
// like "10" and "2" sort lexicographically (so "10" < "2") rather
// than numerically. Vector 008 object_key_ordering names this as a
// trap that catches numeric-aware sorts.
func TestCanonicalize_NumericStringKeysAreLexicographic(t *testing.T) {
	in := map[string]interface{}{
		"10": int64(1),
		"2":  int64(2),
		"1":  int64(3),
	}
	want := "{\"1\":3,\"10\":1,\"2\":2}"
	got, err := Canonicalize(in)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if string(got) != want {
		t.Errorf("Canonicalize = %q, want %q", got, want)
	}
}

// TestCanonicalize_NestedStructures exercises arrays nested inside
// objects and objects nested inside arrays. The byte output is
// pinned end-to-end so a future change to writeArray / writeObject
// dispatch shows up here.
func TestCanonicalize_NestedStructures(t *testing.T) {
	in := map[string]interface{}{
		"outer": []interface{}{
			map[string]interface{}{"b": int64(2), "a": int64(1)},
			int64(42),
			"string",
		},
	}
	want := "{\"outer\":[{\"a\":1,\"b\":2},42,\"string\"]}"
	got, err := Canonicalize(in)
	if err != nil {
		t.Fatalf("Canonicalize: %v", err)
	}
	if string(got) != want {
		t.Errorf("Canonicalize = %q, want %q", got, want)
	}
}

// TestCanonicalize_EmptyContainers covers the empty-object and
// empty-array cases — the smallest valid JCS output is "{}" or
// "[]", and a regression in writeArray / writeObject's prefix /
// suffix handling would surface here.
func TestCanonicalize_EmptyContainers(t *testing.T) {
	cases := []struct {
		in   interface{}
		want string
	}{
		{map[string]interface{}{}, "{}"},
		{[]interface{}{}, "[]"},
		{map[string]interface{}{"empty": []interface{}{}}, "{\"empty\":[]}"},
	}
	for _, tc := range cases {
		got, err := Canonicalize(tc.in)
		if err != nil {
			t.Errorf("Canonicalize(%v): %v", tc.in, err)
			continue
		}
		if string(got) != tc.want {
			t.Errorf("Canonicalize(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// TestCanonicalizeJSON_UseNumberPreservesInt64Precision confirms
// the convenience wrapper preserves int64 precision through the
// JSON decode. Without Decoder.UseNumber, encoding/json decodes
// every number as float64 and the above_safe_integer case
// (9007199254740992 → 9007199254740993 due to float64 precision
// gap) would silently drift.
func TestCanonicalizeJSON_UseNumberPreservesInt64Precision(t *testing.T) {
	in := []byte("{\"v\":9007199254740992}")
	got, err := CanonicalizeJSON(in)
	if err != nil {
		t.Fatalf("CanonicalizeJSON: %v", err)
	}
	want := "{\"v\":9007199254740992}"
	if string(got) != want {
		t.Errorf("CanonicalizeJSON(%q) = %q, want %q", in, got, want)
	}
}

// TestCanonicalize_UnsupportedTypeRejected confirms types we don't
// support surface a localising error rather than producing silently
// wrong output. The error format names the offending type so a
// caller can grep their codebase.
func TestCanonicalize_UnsupportedTypeRejected(t *testing.T) {
	type S struct{ X int }
	_, err := Canonicalize(S{X: 1})
	if err == nil {
		t.Fatal("Canonicalize accepted a struct; want unsupported-type error")
	}
}
