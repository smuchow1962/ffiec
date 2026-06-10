package vectors

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

// writeTempFile writes body to dir/name and returns the path, failing the
// test on error.
func writeTempFile(t *testing.T, dir, name string, body []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatalf("write temp file %s: %v", path, err)
	}
	return path
}

// CRLF-class regression: a CRLF-mangled corpus read MUST be caught, not
// silently tolerated.
//
// The .gitattributes guard PREVENTS the mangling on a fresh clone (LF-pin
// on text, binary on crypto material). This test is the Go-side regression
// for the same bug class: even if the guard is missing or bypassed, the
// byte-exact corpus-read path catches a CRLF-corrupted golden file.
//
// The load-bearing check is checkCanonicalSelfConsistency: it hashes the
// RAW bytes of expected_canonical.txt and compares to the pinned sha256.
// A Windows clone without .gitattributes rewrites LF→CRLF in that file, so
// the raw bytes — and thus the hash — change, and the check fails. That
// loud failure is the safety net.
//
// The JCS-idempotence check is deliberately CRLF-tolerant (splitNDJSON
// strips a trailing \r so a one-document-per-line file re-canonicalizes
// regardless of EOL). That tolerance is fine for idempotence but would
// SILENTLY pass a mangled file on its own — which is exactly why the
// self-consistency check (raw-byte hash) is the one that must catch the
// mangling. This test pins that division of labor.

// crlfMangle rewrites every LF to CRLF, simulating a Git autocrlf checkout
// on Windows of a file that was committed with LF and not protected by
// .gitattributes.
func crlfMangle(b []byte) []byte {
	return bytes.ReplaceAll(b, []byte("\n"), []byte("\r\n"))
}

// TestCanonicalSelfConsistency_CatchesCRLFMangling proves the raw-byte
// hash check rejects a CRLF-mangled canonical blob and accepts the clean
// LF blob. Self-contained: no corpus dependency, so it runs in CI.
func TestCanonicalSelfConsistency_CatchesCRLFMangling(t *testing.T) {
	// A two-line NDJSON canonical blob, the shape a multi-record vector
	// carries. LF-framed, as the corpus commits it.
	clean := []byte("{\"a\":1}\n{\"b\":2}\n")
	sum := sha256.Sum256(clean)
	pin := hex.EncodeToString(sum[:])

	// Clean (LF) bytes: self-consistency PASSES.
	cleanVector := CanonicalVector{Slot: "crlf-regression", Canonical: clean, PrimarySHA256: pin}
	if c := checkCanonicalSelfConsistency(cleanVector); !c.Pass {
		t.Fatalf("clean LF blob failed self-consistency: %s", c.Detail)
	}

	// CRLF-mangled bytes against the LF-derived pin: self-consistency FAILS.
	// This is the regression — a silently-mangled corpus read is caught.
	mangled := crlfMangle(clean)
	if bytes.Equal(mangled, clean) {
		t.Fatal("crlfMangle produced identical bytes — test is not exercising the mangle path")
	}
	mangledVector := CanonicalVector{Slot: "crlf-regression", Canonical: mangled, PrimarySHA256: pin}
	if c := checkCanonicalSelfConsistency(mangledVector); c.Pass {
		t.Fatal("CRLF-mangled blob PASSED self-consistency — a mangled corpus read would go undetected")
	}
}

// TestReadPrimarySHA_CRLFTolerant confirms the sha256-pin reader tolerates
// a CRLF-mangled expected_canonical_sha256.txt: the hex value itself is
// CRLF-agnostic (strings.Fields strips the trailing \r), so the PIN read
// survives mangling even though the canonical-BLOB read does not. That
// asymmetry is correct: the pin is a hex token whose bytes do not include
// an EOL, so its value is unchanged; the blob is the hashed content, whose
// bytes DO change. Pinning both behaviors documents where CRLF matters and
// where it does not.
func TestReadPrimarySHA_CRLFTolerant(t *testing.T) {
	dir := t.TempDir()
	pin := hex.EncodeToString(sha256.New().Sum(nil)) // a well-formed 64-char hex
	lf := []byte("full_blob " + pin + "\n")
	crlf := crlfMangle(lf)

	for _, tc := range []struct {
		name string
		body []byte
	}{
		{"lf", lf},
		{"crlf", crlf},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempFile(t, dir, "sha-"+tc.name+".txt", tc.body)
			got, err := readPrimarySHA(path)
			if err != nil {
				t.Fatalf("readPrimarySHA(%s): %v", tc.name, err)
			}
			if got != pin {
				t.Errorf("readPrimarySHA(%s) = %q, want %q (EOL must not corrupt the hex token)", tc.name, got, pin)
			}
		})
	}
}
