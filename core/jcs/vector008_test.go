package jcs

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
)

// vector008CorpusDirEnvVar lets a contributor override the corpus
// location for the vector-008 test. Default falls back to the same
// ancestor-walk logic the conformance-gate runner uses.
const vector008CorpusDirEnvVar = "FFIEC_PUBLIC_VECTORS_DIR"

// TestCanonicalize_Vector008Conformance is the load-bearing
// integration test for this package. It loads the materialized
// vector 008 fixture + expected pins from ffiec-public/spec/test-
// vectors/008-jcs-edge-cases/ and runs every case (61 across 9
// categories) against Canonicalize.
//
// Per the vector's description.md: "An implementation is v1.0-
// conformant only if it reproduces every byte sequence in
// expected.json for the corresponding input in fixture.json, AND
// rejects every input in the nan_infinity_rejection category with
// an error."
//
// Skips with a diagnostic when the corpus isn't available locally
// (CI machines or contributors without ffiec-public side-by-side);
// fails with a per-case breakdown when any case diverges.
func TestCanonicalize_Vector008Conformance(t *testing.T) {
	dir, ok := findVector008Dir(t)
	if !ok {
		return
	}

	fx, exp, err := loadVector008(dir)
	if err != nil {
		t.Fatalf("load vector 008: %v", err)
	}

	var (
		pass, fail, skipped int
		failures            []string
	)

	for category, fixtureCat := range fx.Categories {
		expectedCat, ok := exp.Categories[category]
		if !ok {
			t.Errorf("category %q present in fixture but absent in expected", category)
			continue
		}
		for caseName, fixtureCase := range fixtureCat.Cases {
			expectedCase, ok := expectedCat.Cases[caseName]
			if !ok {
				t.Errorf("case %s/%s present in fixture but absent in expected", category, caseName)
				continue
			}

			input := decodeFixtureInput(fixtureCase.Input)

			switch {
			case expectedCase.Error != "":
				// Rejection case: Canonicalize MUST return an error.
				if _, err := Canonicalize(input); err == nil {
					failures = append(failures, fmt.Sprintf(
						"%s/%s: expected error %s but Canonicalize returned no error",
						category, caseName, expectedCase.Error))
					fail++
				} else {
					pass++
				}
			case expectedCase.CanonicalHex != "":
				wantBytes, hexErr := hex.DecodeString(expectedCase.CanonicalHex)
				if hexErr != nil {
					t.Errorf("%s/%s: malformed canonical_hex in fixture: %v", category, caseName, hexErr)
					continue
				}
				gotBytes, err := Canonicalize(input)
				if err != nil {
					failures = append(failures, fmt.Sprintf(
						"%s/%s: Canonicalize errored — %v (expected bytes %s)",
						category, caseName, err, expectedCase.CanonicalHex))
					fail++
					continue
				}
				if !bytes.Equal(gotBytes, wantBytes) {
					failures = append(failures, fmt.Sprintf(
						"%s/%s:\n  got:  %s (%q)\n  want: %s (%q)",
						category, caseName,
						hex.EncodeToString(gotBytes), string(gotBytes),
						expectedCase.CanonicalHex, expectedCase.CanonicalUTF8))
					fail++
					continue
				}
				pass++
			default:
				skipped++
			}
		}
	}

	t.Logf("vector 008 conformance: %d passed, %d failed, %d skipped (no expected bytes pinned)",
		pass, fail, skipped)
	for _, f := range failures {
		t.Errorf("vector 008: %s", f)
	}
	if fail > 0 {
		t.Fatalf("vector 008 conformance: %d of %d cases failed", fail, pass+fail)
	}
}

// findVector008Dir locates the 008-jcs-edge-cases directory using
// the same env-var-then-ancestor-walk strategy the conformance-gate
// runner uses. Returns false (with a skip on the test) when the
// corpus is not available locally.
func findVector008Dir(t *testing.T) (string, bool) {
	t.Helper()
	const subdir = "008-jcs-edge-cases"

	candidates := []string{}
	if v := os.Getenv(vector008CorpusDirEnvVar); v != "" {
		candidates = append(candidates, filepath.Join(v, subdir))
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Skipf("getwd failed: %v", err)
		return "", false
	}
	for dir := cwd; ; {
		candidates = append(candidates, filepath.Join(dir, "..", "ffiec-public", "spec", "test-vectors", subdir))
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	tried := []string{}
	for _, c := range candidates {
		abs, err := filepath.Abs(c)
		if err != nil {
			tried = append(tried, fmt.Sprintf("%s (abs failed: %v)", c, err))
			continue
		}
		if _, err := os.Stat(filepath.Join(abs, "fixture.json")); err == nil {
			return abs, true
		}
		tried = append(tried, abs)
	}
	t.Skipf("vector 008 corpus not found; set %s to override (tried %d candidates)",
		vector008CorpusDirEnvVar, len(tried))
	return "", false
}

// loadVector008 reads fixture.json + expected.json from the vector's
// directory and decodes them into typed shapes.
func loadVector008(dir string) (*vector008Fixture, *vector008Expected, error) {
	fxRaw, err := os.ReadFile(filepath.Join(dir, "fixture.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("read fixture.json: %w", err)
	}
	expRaw, err := os.ReadFile(filepath.Join(dir, "expected.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("read expected.json: %w", err)
	}

	// Decode with UseNumber so int64 precision survives the round-trip
	// for the numeric_edge_cases above-safe-integer case.
	var fx vector008Fixture
	dec := json.NewDecoder(bytes.NewReader(fxRaw))
	dec.UseNumber()
	if err := dec.Decode(&fx); err != nil {
		return nil, nil, fmt.Errorf("decode fixture: %w", err)
	}
	var exp vector008Expected
	if err := json.Unmarshal(expRaw, &exp); err != nil {
		return nil, nil, fmt.Errorf("decode expected: %w", err)
	}
	return &fx, &exp, nil
}

// vector008Fixture mirrors the fixture.json top-level shape. Each
// category's Cases map carries one entry per named case, whose
// Input field is the JSON value to canonicalize.
type vector008Fixture struct {
	About      string                              `json:"_about"`
	Categories map[string]vector008FixtureCategory `json:"categories"`
}

type vector008FixtureCategory struct {
	Description string                          `json:"description"`
	Cases       map[string]vector008FixtureCase `json:"cases"`
}

type vector008FixtureCase struct {
	Input       interface{} `json:"input"`
	Description string      `json:"description,omitempty"`
}

// vector008Expected mirrors expected.json. Each case carries either
// CanonicalHex + CanonicalUTF8 + CanonicalSHA256 (success case) or
// Error + ErrorClass (rejection case).
type vector008Expected struct {
	Categories map[string]vector008ExpectedCategory `json:"categories"`
}

type vector008ExpectedCategory struct {
	Description string                           `json:"description"`
	Cases       map[string]vector008ExpectedCase `json:"cases"`
}

type vector008ExpectedCase struct {
	CanonicalHex    string `json:"canonical_hex,omitempty"`
	CanonicalUTF8   string `json:"canonical_utf8,omitempty"`
	CanonicalSHA256 string `json:"canonical_sha256,omitempty"`
	CanonicalLength int    `json:"canonical_length,omitempty"`
	Error           string `json:"error,omitempty"`
	ErrorClass      string `json:"error_class,omitempty"`
}

// decodeFixtureInput converts a fixture's Input field into the form
// Canonicalize expects. The fixture encodes NaN / Infinity as
// {"__special__": "NaN"} markers because JSON can't carry the raw
// float values; we detect the marker and substitute the native
// non-finite float so the rejection path is exercised the same way
// a production-decoded input would exercise it.
//
// The recursion walks every nested map and array so a marker placed
// arbitrarily deep is still detected.
func decodeFixtureInput(v interface{}) interface{} {
	switch x := v.(type) {
	case map[string]interface{}:
		// Special-marker detection.
		if len(x) == 1 {
			if marker, ok := x["__special__"]; ok {
				if s, isString := marker.(string); isString {
					switch s {
					case "NaN":
						return math.NaN()
					case "Infinity":
						return math.Inf(1)
					case "-Infinity":
						return math.Inf(-1)
					}
				}
			}
		}
		out := make(map[string]interface{}, len(x))
		for k, val := range x {
			out[k] = decodeFixtureInput(val)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(x))
		for i, elem := range x {
			out[i] = decodeFixtureInput(elem)
		}
		return out
	case json.Number:
		// Mirror what CanonicalizeJSON's UseNumber path does: emit
		// int64 when it fits, float64 otherwise. This keeps the test
		// honest about integer precision the production decode path
		// preserves.
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x
	}
	return v
}

// (helper removed — using bytes.Equal directly from the stdlib)
