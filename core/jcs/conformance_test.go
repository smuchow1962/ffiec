package jcs_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmpworks/ffiec/core/jcs"
)

// vectorsDirEnvVar mirrors the verifier loader's env-var name so a
// contributor sets one path and both the core and verifier conformance
// gates find the corpus.
const vectorsDirEnvVar = "FFIEC_PUBLIC_VECTORS_DIR"

// defaultVectorsDir is the side-by-side ffiec + ffiec-public layout.
const defaultVectorsDir = "../../../ffiec-public/spec/test-vectors"

// jcsCorpusMarker is a file that exists only inside the 008 corpus
// directory; we use it to confirm a candidate path is the right one.
const jcsCorpusMarker = "008-jcs-edge-cases/fixture.json"

// resolveVectorsDir finds the test-vector corpus. Resolution order:
//  1. $FFIEC_PUBLIC_VECTORS_DIR.
//  2. defaultVectorsDir relative to the cwd and each ancestor.
//
// Returns ("", false) when no candidate resolves — the caller SHOULD
// t.Skip so a contributor without ffiec-public checked out still gets
// a green local build, with the gate re-arming where the corpus exists.
func resolveVectorsDir() (string, bool) {
	if v := os.Getenv(vectorsDirEnvVar); v != "" {
		if ok := hasMarker(v); ok {
			return v, true
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for dir := cwd; ; {
		candidate := filepath.Join(dir, defaultVectorsDir)
		if hasMarker(candidate) {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func hasMarker(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, jcsCorpusMarker))
	return err == nil
}

// jcsFixture / jcsExpected mirror the 008 corpus shape:
// categories.<cat>.cases.<name>.input (a JSON object) on the fixture
// side; .canonical_utf8 / .canonical_sha256 (or .error) on the
// expected side.
type jcsFixture struct {
	Categories map[string]struct {
		Cases map[string]struct {
			Input json.RawMessage `json:"input"`
		} `json:"cases"`
	} `json:"categories"`
}

type jcsExpected struct {
	Categories map[string]struct {
		Cases map[string]struct {
			CanonicalUTF8   string `json:"canonical_utf8"`
			CanonicalSHA256 string `json:"canonical_sha256"`
			Error           string `json:"error"`
		} `json:"cases"`
	} `json:"categories"`
}

// TestJCS_008Corpus exercises the canonicalizer against the
// 008-jcs-edge-cases corpus — the same answer key the .NET reference's
// JcsRfc8785ConformanceTests consumes. Per spec §5 (MUST in the
// v1.0-final-amendment), an implementation passing the basic fixtures
// but failing any 008 case is non-conformant. This is the load-bearing
// JCS conformance gate.
func TestJCS_008Corpus(t *testing.T) {
	dir, ok := resolveVectorsDir()
	if !ok {
		t.Skipf("test-vector corpus not found; set %s to enable the JCS conformance gate", vectorsDirEnvVar)
	}

	base := filepath.Join(dir, "008-jcs-edge-cases")
	fx := loadJCSFixture(t, filepath.Join(base, "fixture.json"))
	exp := loadJCSExpected(t, filepath.Join(base, "expected.json"))

	for cat, fc := range fx.Categories {
		ec, ok := exp.Categories[cat]
		if !ok {
			t.Errorf("category %q present in fixture but absent in expected", cat)
			continue
		}
		for name, fcase := range fc.Cases {
			ecase, ok := ec.Cases[name]
			if !ok {
				t.Errorf("case %s/%s present in fixture but absent in expected", cat, name)
				continue
			}
			t.Run(cat+"/"+name, func(t *testing.T) {
				runJCSCase(t, fcase.Input, ecase.CanonicalUTF8, ecase.CanonicalSHA256, ecase.Error)
			})
		}
	}
}

// runJCSCase canonicalizes one fixture input and asserts either the
// pinned canonical bytes (success case) or rejection (error case).
func runJCSCase(t *testing.T, input json.RawMessage, wantUTF8, wantSHA, wantErr string) {
	t.Helper()

	value, special := decodeJCSInput(t, input)

	// Rejection cases carry an `error` on the expected side and encode
	// the non-finite float via the {"__special__": "..."} marker.
	if wantErr != "" || special {
		_, err := jcs.Canonicalize(value)
		if err == nil {
			t.Fatalf("expected rejection (%s) but canonicalize succeeded", wantErr)
		}
		return
	}

	got, err := jcs.Canonicalize(value)
	if err != nil {
		t.Fatalf("canonicalize failed: %v", err)
	}

	if string(got) != wantUTF8 {
		t.Errorf("canonical bytes diverge\n got: %q\nwant: %q", string(got), wantUTF8)
	}

	if wantSHA != "" {
		sum := sha256.Sum256(got)
		gotSHA := hex.EncodeToString(sum[:])
		if gotSHA != wantSHA {
			t.Errorf("canonical sha256 mismatch\n got: %s\nwant: %s", gotSHA, wantSHA)
		}
	}
}

// decodeJCSInput decodes a fixture input into the `any` shape the
// canonicalizer accepts, using UseNumber so integral values keep their
// exact digits. It also detects the {"__special__": "NaN"|...} marker
// the corpus uses to encode non-finite floats (which JSON cannot
// carry) and substitutes the native float, returning special=true so
// the caller asserts rejection.
func decodeJCSInput(t *testing.T, raw json.RawMessage) (value any, special bool) {
	t.Helper()
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&value); err != nil {
		t.Fatalf("decode fixture input: %v", err)
	}
	return substituteSpecials(value)
}

// substituteSpecials walks the decoded value replacing any
// {"__special__": "NaN"|"Infinity"|"-Infinity"} object with the native
// float64 it encodes. Returns special=true when at least one marker was
// found, so the case is treated as a rejection case.
func substituteSpecials(value any) (any, bool) {
	m, ok := value.(map[string]any)
	if !ok {
		return value, false
	}
	if s, ok := specialMarker(m); ok {
		return s, true
	}
	found := false
	for k, v := range m {
		nv, sp := substituteSpecials(v)
		m[k] = nv
		found = found || sp
	}
	return m, found
}

func specialMarker(m map[string]any) (float64, bool) {
	if len(m) != 1 {
		return 0, false
	}
	v, ok := m["__special__"]
	if !ok {
		return 0, false
	}
	switch v {
	case "NaN":
		return math.NaN(), true
	case "Infinity":
		return math.Inf(1), true
	case "-Infinity":
		return math.Inf(-1), true
	}
	return 0, false
}

func loadJCSFixture(t *testing.T, path string) jcsFixture {
	t.Helper()
	var fx jcsFixture
	mustDecodeFile(t, path, &fx)
	return fx
}

func loadJCSExpected(t *testing.T, path string) jcsExpected {
	t.Helper()
	var exp jcsExpected
	mustDecodeFile(t, path, &exp)
	return exp
}

func mustDecodeFile(t *testing.T, path string, into any) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}
