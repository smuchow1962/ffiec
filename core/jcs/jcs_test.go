package jcs_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
)

// TestCanonicalize_Numbers pins the ECMA-262 NumberToString dispatch
// without a corpus dependency, so the float-formatting logic is gated
// on every machine. The expected forms are the RFC 8785 canonical
// representations the 008 corpus also pins.
func TestCanonicalize_Numbers(t *testing.T) {
	cases := []struct {
		name string
		in   float64
		want string
	}{
		{"integer_one", 1.0, "1"},
		{"small_decimal", 0.1, "0.1"},
		{"one_third", 0.3333333333333333, "0.3333333333333333"},
		{"ten_billion", 10000000000.0, "10000000000"},
		{"exp_high_threshold", 1e21, "1e+21"},
		{"below_exp_threshold", 1e20, "100000000000000000000"},
		{"exp_low_threshold", 1e-6, "0.000001"},
		{"below_low_threshold", 1e-7, "1e-7"},
		{"negative", -42.5, "-42.5"},
		{"negative_zero", -0.0, "0"},
		{"positive_zero", 0.0, "0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jcs.Canonicalize(map[string]any{"v": tc.in})
			if err != nil {
				t.Fatalf("canonicalize: %v", err)
			}
			want := `{"v":` + tc.want + "}"
			if string(got) != want {
				t.Errorf("got %q, want %q", string(got), want)
			}
		})
	}
}

// TestCanonicalize_KeyOrdering asserts byte-ordinal key sorting — the
// property that makes the output deterministic across implementations.
func TestCanonicalize_KeyOrdering(t *testing.T) {
	in := map[string]any{"b": 2.0, "a": 1.0, "c": 3.0, "A": 0.0}
	got, err := jcs.Canonicalize(in)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	// Uppercase 'A' (0x41) sorts before lowercase letters (0x61+).
	want := `{"A":0,"a":1,"b":2,"c":3}`
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}

// TestCanonicalize_StringEscapes pins the minimal RFC 8785 escape set:
// only the named short escapes + lowercase \u00xx for other C0
// controls; forward slash, DEL, and non-ASCII pass through raw.
func TestCanonicalize_StringEscapes(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"quote", "a\"b", "\"a\\\"b\""},
		{"backslash", "a\\b", "\"a\\\\b\""},
		{"newline", "a\nb", "\"a\\nb\""},
		{"tab", "a\tb", "\"a\\tb\""},
		{"unit_separator", "a\x1fb", "\"a\\u001fb\""}, // C0 control → lowercase \u00xx
		{"forward_slash", "a/b", "\"a/b\""},           // NOT escaped
		{"del", "a\x7fb", "\"a\x7fb\""},               // U+007F raw
		{"unicode_passthrough", "café", "\"café\""},   // raw UTF-8
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := jcs.Canonicalize(tc.in)
			if err != nil {
				t.Fatalf("canonicalize: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %q, want %q", string(got), tc.want)
			}
		})
	}
}

// TestCanonicalize_RejectsNonFinite asserts NaN and ±Inf are rejected
// — they have no canonical JSON form per RFC 8785 §3.2.2.3.
func TestCanonicalize_RejectsNonFinite(t *testing.T) {
	for _, bad := range []float64{
		mustInf(1), mustInf(-1), mustNaN(),
	} {
		if _, err := jcs.Canonicalize(bad); err == nil {
			t.Errorf("expected rejection for %v", bad)
		}
	}
}

// TestCanonicalize_IntegralPreservation asserts a json.Number integer
// past 2^53 keeps its exact digits (no float64 round-trip loss). This
// is the contract a typed caller using UseNumber() relies on.
func TestCanonicalize_IntegralPreservation(t *testing.T) {
	const big = "9007199254740993" // 2^53 + 1, not representable in float64
	var v any
	dec := json.NewDecoder(strings.NewReader(`{"n":` + big + `}`))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	got, err := jcs.Canonicalize(v)
	if err != nil {
		t.Fatalf("canonicalize: %v", err)
	}
	want := `{"n":` + big + `}`
	if string(got) != want {
		t.Errorf("got %q, want %q", string(got), want)
	}
}

func mustInf(sign int) float64 {
	if sign < 0 {
		return -1.0 / zero()
	}
	return 1.0 / zero()
}

func mustNaN() float64 { return zero() / zero() }

// zero defeats the compiler's constant-folding so 1/0 and 0/0 produce
// runtime Inf/NaN rather than a compile error.
func zero() float64 {
	var z float64
	return z
}
