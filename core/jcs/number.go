package jcs

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// asJSONNumber recognizes encoding/json's json.Number — the type a
// caller decoding with json.Decoder.UseNumber() produces for every
// numeric value. Matching the concrete type (not a broad Stringer
// interface) keeps the dispatch closed: only a real json.Number takes
// the integral-preservation path; any other type falls through to the
// unsupported-type error in writeValue.
func asJSONNumber(value any) (string, bool) {
	if n, ok := value.(json.Number); ok {
		return string(n), true
	}
	return "", false
}

// writeJSONNumber emits a json.Number. Integral values are written as
// their exact digit string (no float64 round-trip, which would lose
// precision past 2^53). Non-integral values route through the float64
// ECMA-262 path so they match the corpus's pinned float forms.
func writeJSONNumber(b *strings.Builder, num string) error {
	if !strings.ContainsAny(num, ".eE") {
		// Pure integer literal: emit verbatim after validating it
		// parses, so a malformed number is a loud error not a silent
		// pass-through.
		if _, err := strconv.ParseInt(num, 10, 64); err != nil {
			// Could be a big integer beyond int64; preserve the digits
			// if they are all decimal (RFC 8785 keeps integer digits).
			if isAllDigits(num) {
				b.WriteString(num)
				return nil
			}
			return fmt.Errorf("jcs: malformed integer number %q: %w", num, err)
		}
		b.WriteString(num)
		return nil
	}
	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return fmt.Errorf("jcs: malformed number %q: %w", num, err)
	}
	return writeNumber(b, f)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' {
		start = 1
	}
	if start == len(s) {
		return false
	}
	for i := start; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// writeNumber emits value per RFC 8785 §3.2.2.3, which defers to
// ECMA-262 §6.1.6.1.13 NumberToString. NaN and ±Inf are rejected:
// they have no canonical JSON form.
func writeNumber(b *strings.Builder, value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("jcs: NaN/Infinity has no canonical form (got %v)", value)
	}
	if value == 0 {
		// RFC 8785: both +0 and -0 serialize as "0".
		b.WriteByte('0')
		return nil
	}
	b.WriteString(formatEcma262(value))
	return nil
}

// formatEcma262 implements ECMA-262 §6.1.6.1.13 NumberToString for a
// finite, non-zero double. The algorithm finds integers s, k, n with
// 10^(k-1) ≤ s < 10^k and s × 10^(n-k) = value and k minimal, then
// dispatches on n's relationship to k to pick fixed-form,
// leading-zero form, or exponent form.
//
// This is a line-by-line port of the .NET reference's
// CanonicalJson.FormatEcma262Number so the two produce byte-identical
// output. Go's strconv.FormatFloat(value, 'g', -1, 64) is the
// shortest-round-trip decimal — the equivalent of .NET's
// double.ToString("R"). We parse it to extract the significant digits
// and the decimal-point position, then apply the case dispatch.
//
// Cognitive-complexity note: the function is long because the ECMA-262
// spec is a five-case dispatch. Each case is a single return with a
// comment naming the spec case (d/e/f/g/h); the surrounding code only
// normalizes the parsed digits. Splitting the cases into helpers would
// scatter the dispatch and obscure that they are mutually exclusive
// branches of one spec rule, so they stay inline.
func formatEcma262(value float64) string {
	negative := value < 0
	abs := math.Abs(value)

	// Shortest round-trip decimal. Examples:
	//   1.0  → "1"        1.5   → "1.5"
	//   1e21 → "1e+21"    1e-7  → "1e-07"
	//   0.1  → "0.1"
	raw := strconv.FormatFloat(abs, 'g', -1, 64)

	// Parse: <int>[.<frac>][e[+-]?<exp>]
	mantissa := raw
	explicitExp := 0
	if eIdx := strings.IndexAny(raw, "eE"); eIdx >= 0 {
		mantissa = raw[:eIdx]
		exp, err := strconv.Atoi(raw[eIdx+1:])
		if err == nil {
			explicitExp = exp
		}
	}

	// Extract significant digits and the integer-part length (the
	// position of the implicit decimal point relative to those digits).
	var digits string
	var integerLen int
	if dotIdx := strings.IndexByte(mantissa, '.'); dotIdx < 0 {
		digits = mantissa
		integerLen = len(digits)
	} else {
		digits = mantissa[:dotIdx] + mantissa[dotIdx+1:]
		integerLen = dotIdx
	}

	// Strip leading zeros (defensive — "R"/'g' shouldn't emit them for
	// a non-zero abs value, but normalize so k counts true significant
	// digits).
	leadingZeros := 0
	for leadingZeros < len(digits)-1 && digits[leadingZeros] == '0' {
		leadingZeros++
	}
	digits = digits[leadingZeros:]
	integerLen -= leadingZeros

	// Strip trailing zeros (positional, not significant). n is
	// unaffected: dropping low-significance digits shifts s and k by
	// the same amount.
	trailingZeros := 0
	for trailingZeros < len(digits)-1 && digits[len(digits)-1-trailingZeros] == '0' {
		trailingZeros++
	}
	if trailingZeros > 0 {
		digits = digits[:len(digits)-trailingZeros]
	}

	k := len(digits)
	n := integerLen + explicitExp

	sign := ""
	if negative {
		sign = "-"
	}

	// ECMA-262 §6.1.6.1.13 case dispatch:
	switch {
	case k <= n && n <= 21:
		// case d: digits followed by (n-k) zeros.
		return sign + digits + strings.Repeat("0", n-k)
	case 0 < n && n <= 21:
		// case e: first n digits, decimal point, remaining digits.
		return sign + digits[:n] + "." + digits[n:]
	case -6 < n && n <= 0:
		// case f: "0." then (-n) zeros then digits.
		return sign + "0." + strings.Repeat("0", -n) + digits
	default:
		// cases g/h: exponent form. Exponent is n-1, explicit sign,
		// no leading zeros, lowercase 'e'.
		exp := n - 1
		expStr := strconv.Itoa(exp)
		if exp >= 0 {
			expStr = "+" + expStr // strconv omits the + for non-negatives
		}
		if k == 1 {
			return sign + digits + "e" + expStr
		}
		return sign + digits[:1] + "." + digits[1:] + "e" + expStr
	}
}
