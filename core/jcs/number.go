package jcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// writeInt64 emits a signed integer in decimal form. JCS §3.2.2.3
// requires IEEE-754-double serialization, but integers in the safe
// range [-(2^53-1), 2^53-1] have an exact double representation and
// serialize as their decimal-integer form (no decimal point, no
// exponent).
//
// Integers ABOVE the safe range are accepted (Go's int64 can hold
// them); they're emitted in decimal form per the
// `above_safe_integer` case in vector 008 numeric_edge_cases
// (input 9007199254740992 → output "9007199254740992"). The vector
// is silent on integer inputs >= 10^21 because Go's int64 maxes at
// ~9.2e18 — far short of the 1e21 threshold where ECMAScript
// number-to-string switches to scientific notation. v1.0 emits all
// in-range int64 values as decimal.
func writeInt64(buf *bytes.Buffer, v int64) error {
	buf.WriteString(strconv.FormatInt(v, 10))
	return nil
}

// writeUint64 emits an unsigned integer in decimal form. Same
// rationale as writeInt64; the safe-range cap is observed but
// uint64 values up to math.MaxUint64 = 2^64-1 fall short of the
// 1e21 scientific-notation threshold.
func writeUint64(buf *bytes.Buffer, v uint64) error {
	buf.WriteString(strconv.FormatUint(v, 10))
	return nil
}

// writeJSONNumber emits a json.Number per the callers that decoded
// JSON with Decoder.UseNumber. The branch logic: if the value
// parses cleanly as int64 we serialize as a decimal integer
// (preserving precision the float64 round-trip would lose); else
// we fall through to the float64 path.
func writeJSONNumber(buf *bytes.Buffer, n json.Number) error {
	if i, err := n.Int64(); err == nil {
		return writeInt64(buf, i)
	}
	f, err := n.Float64()
	if err != nil {
		return fmt.Errorf("jcs: parse json.Number %q: %w", n, err)
	}
	return writeFloat64(buf, f)
}

// writeFloat64 emits a double-precision float per RFC 8785 §3.2.2.3
// (ECMAScript Number.prototype.toString).
//
// Algorithm:
//
//  1. NaN or Infinity → return ErrNonFinite per §3.2.2.3.
//  2. Negative zero → emit "0" (not "-0") per vector 008
//     float_canonicalization negative_zero case.
//  3. Integer-valued float in the safe range → emit as integer.
//  4. For everything else, use Go's strconv with 'e' format to get
//     the shortest-round-trip digits + the exponent, then format
//     per the ECMAScript thresholds:
//     - 1e-6 <= |n| < 1e21 → fixed-decimal form
//     - otherwise → exponential form with lowercase 'e', explicit
//       sign, no leading zeros in the exponent
//
// Implementation notes:
//
//   - The integer-detection happens BEFORE the strconv call so values
//     like 1.0 and 10000000000.0 emit as "1" and "10000000000" rather
//     than going through the float-format path.
//   - The fixed-decimal expansion walks the mantissa digits + the
//     decimal-point position implied by the exponent; it does NOT
//     call strconv.FormatFloat with 'f' verb because that path can
//     emit redundant trailing zeros for very-small-magnitude inputs.
//   - The exponential form's exponent value strips Go's two-digit
//     zero padding ("1e-07" → "1e-7", "1e+00" → not emitted because
//     0.0 hits the negative-zero branch above).
func writeFloat64(buf *bytes.Buffer, f float64) error {
	switch {
	case math.IsNaN(f):
		return fmt.Errorf("jcs: %w (NaN)", ErrNonFinite)
	case math.IsInf(f, 1):
		return fmt.Errorf("jcs: %w (+Inf)", ErrNonFinite)
	case math.IsInf(f, -1):
		return fmt.Errorf("jcs: %w (-Inf)", ErrNonFinite)
	}

	// Negative-zero normalization: -0.0 and 0.0 both serialize as "0".
	// Without this branch the strconv path produces "-0" for -0.0,
	// which violates the vector 008 negative_zero pin.
	if f == 0 {
		buf.WriteByte('0')
		return nil
	}

	// Integer-valued float in the int64 safe range — emit as integer.
	// math.Trunc(f) == f catches integer-valued doubles; the bound check
	// avoids overflow on math.MaxFloat64 (which is integer-valued in the
	// math.Trunc sense but vastly exceeds int64 range).
	if !math.IsInf(f, 0) && f == math.Trunc(f) && f >= math.MinInt64 && f <= math.MaxInt64 {
		// Additional precision guard: only treat as integer when the
		// magnitude is below 2^53 (the IEEE-754-double safe-integer
		// boundary). Above 2^53 the float64 representation has gaps and
		// the integer form might mislead — but vector 008's
		// `above_safe_integer` case explicitly expects the integer form,
		// so we accept the precision risk for consistency with the spec.
		buf.WriteString(strconv.FormatInt(int64(f), 10))
		return nil
	}

	// Get the shortest-round-trip mantissa + exponent via Go's 'e' verb.
	// The output shape is "[-]d.ddddde[+|-]dd" — predictable + parseable.
	es := strconv.FormatFloat(f, 'e', -1, 64)
	mantissa, expSign, expAbs, err := splitExp(es)
	if err != nil {
		// Defensive: strconv's 'e' output is well-formed for finite floats,
		// so this branch is unreachable in practice. Surfacing as a typed
		// error keeps the caller's `errors.Is` story consistent.
		return fmt.Errorf("jcs: parse strconv output %q: %w", es, err)
	}

	// Compute the effective decimal exponent (the power of ten the
	// canonical integer-mantissa carries). This is what determines the
	// fixed-vs-exponent dispatch below.
	exp10 := expAbs
	if expSign == '-' {
		exp10 = -exp10
	}

	// ECMAScript Number-to-String dispatch: fixed-decimal form for
	// -6 <= effective_exp < 21; exponential form otherwise. The exact
	// boundary values come from the float_canonicalization vector cases.
	mantissaSign, mantissaDigits, fracOffset := decomposeMantissa(mantissa)
	totalDigits := len(mantissaDigits)
	// The decimal exponent of the LEADING digit, with the mantissa
	// treated as a whole integer. E.g., mantissa "1.5" with exp +1
	// represents 15, leading-digit exponent = exp + (totalDigits - 1 -
	// fracOffset). We want the exponent of the MSB digit; that's the
	// fixed-vs-scientific discriminator.
	msbExp := exp10 + (totalDigits - 1 - fracOffset)

	if msbExp >= -6 && msbExp < 21 {
		writeFixedDecimal(buf, mantissaSign, mantissaDigits, msbExp)
		return nil
	}
	writeScientific(buf, mantissaSign, mantissaDigits, msbExp)
	return nil
}

// splitExp parses strconv.FormatFloat(_, 'e', -1, 64) output into its
// (mantissa, exponent-sign, exponent-absolute-value) components.
// Input shape: "[-]d.ddddde[+|-]dd". The mantissa is returned with
// its decimal point intact (or absent for whole-mantissa cases like
// "1e+00").
func splitExp(s string) (mantissa string, expSign byte, expAbs int, err error) {
	i := strings.IndexByte(s, 'e')
	if i < 0 {
		return "", 0, 0, fmt.Errorf("no exponent marker in %q", s)
	}
	mantissa = s[:i]
	rest := s[i+1:]
	if len(rest) < 2 || (rest[0] != '+' && rest[0] != '-') {
		return "", 0, 0, fmt.Errorf("malformed exponent in %q", s)
	}
	expSign = rest[0]
	expAbs, err = strconv.Atoi(rest[1:])
	if err != nil {
		return "", 0, 0, fmt.Errorf("parse exponent in %q: %w", s, err)
	}
	return mantissa, expSign, expAbs, nil
}

// decomposeMantissa splits "[-]d[.dddd]" into:
//   - sign: '-' or 0
//   - digits: the concatenated digit string with the decimal point
//     removed (e.g., "1.5" → "15", "3.333333" → "3333333")
//   - fracOffset: the number of digits AFTER the original decimal
//     point (e.g., "1.5" → 1, "3" → 0)
//
// The (digits, fracOffset) pair represents the mantissa as an
// integer with an implied decimal-point position; the writeFixed /
// writeScientific helpers below reuse this representation.
func decomposeMantissa(m string) (sign byte, digits string, fracOffset int) {
	if len(m) > 0 && m[0] == '-' {
		sign = '-'
		m = m[1:]
	}
	if dot := strings.IndexByte(m, '.'); dot >= 0 {
		fracOffset = len(m) - dot - 1
		digits = m[:dot] + m[dot+1:]
		return sign, digits, fracOffset
	}
	return sign, m, 0
}

// writeFixedDecimal emits the canonical fixed-decimal form for a
// value whose MSB exponent is in [-6, 21). The output rules:
//
//   - msbExp >= 0 and msbExp+1 >= len(digits) → integer-form
//     (left-pad with zeros if needed). E.g., digits="1", msbExp=10
//     → "10000000000".
//   - msbExp >= 0 and msbExp+1 < len(digits) → split the digit
//     string at position msbExp+1, emit with a decimal point.
//     E.g., digits="15", msbExp=0 → "1.5".
//   - msbExp < 0 → "0." + (-msbExp - 1) zeros + digits.
//     E.g., digits="1", msbExp=-6 → "0.000001".
//
// The function is small enough to read end-to-end; each branch
// implements one of the three structural cases.
func writeFixedDecimal(buf *bytes.Buffer, sign byte, digits string, msbExp int) {
	if sign == '-' {
		buf.WriteByte('-')
	}
	switch {
	case msbExp >= 0 && msbExp+1 >= len(digits):
		// Integer form: emit digits, then trailing zeros to pad out to
		// msbExp+1 positions.
		buf.WriteString(digits)
		for i := 0; i < msbExp+1-len(digits); i++ {
			buf.WriteByte('0')
		}
	case msbExp >= 0:
		// Mixed integer + fraction: split at msbExp+1.
		intPart := digits[:msbExp+1]
		fracPart := digits[msbExp+1:]
		buf.WriteString(intPart)
		buf.WriteByte('.')
		buf.WriteString(fracPart)
	default:
		// Pure fraction (msbExp < 0): "0." + leading-zero pad + digits.
		buf.WriteString("0.")
		for i := 0; i < -msbExp-1; i++ {
			buf.WriteByte('0')
		}
		buf.WriteString(digits)
	}
}

// writeScientific emits the canonical exponential form for a value
// whose MSB exponent is outside [-6, 21). The output shape:
// "[-]d[.dddd]e[+|-]N" where N is the MSB exponent with NO leading
// zeros (strip Go's two-digit padding).
//
// Special case: when len(digits) == 1 the mantissa has no decimal
// point ("1e+21", not "1.0e+21"). When len(digits) > 1 the decimal
// point splits between digits[0] and digits[1:] ("1.234e-100").
func writeScientific(buf *bytes.Buffer, sign byte, digits string, msbExp int) {
	if sign == '-' {
		buf.WriteByte('-')
	}
	buf.WriteByte(digits[0])
	if len(digits) > 1 {
		buf.WriteByte('.')
		buf.WriteString(digits[1:])
	}
	buf.WriteByte('e')
	if msbExp >= 0 {
		buf.WriteByte('+')
	} else {
		buf.WriteByte('-')
	}
	buf.WriteString(strconv.Itoa(absInt(msbExp)))
}

// absInt returns the absolute value of x as an int. Avoids importing
// math.Abs (which is float64-typed) for a one-line helper.
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
