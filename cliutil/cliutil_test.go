package cliutil

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ----- errors.go ----------------------------------------------------------

func TestUsagef_TagsAsUsageError(t *testing.T) {
	err := Usagef("missing %s", "thing")
	if !IsUsage(err) {
		t.Fatal("Usagef result is not detected as UsageError")
	}
	if !strings.Contains(err.Error(), "missing thing") {
		t.Errorf("error message wrong: %q", err.Error())
	}
}

func TestIsUsage_FalseForOtherErrors(t *testing.T) {
	if IsUsage(nil) {
		t.Error("IsUsage(nil) should be false")
	}
	if IsUsage(errors.New("plain error")) {
		t.Error("IsUsage(plain) should be false")
	}
}

func TestExitCodeFor(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"plain", errors.New("boom"), 1},
		{"usage", Usagef("nope"), 2},
		{"wrapped usage", &wrapErr{Usagef("nope")}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExitCodeFor(tc.err); got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestReportAndExit_PrintsAndReturnsCode(t *testing.T) {
	cases := []struct {
		name      string
		err       error
		wantCode  int
		wantInErr string
	}{
		{"nil → 0", nil, 0, ""},
		{"plain → 1", errors.New("boom"), 1, "ledgerctl: boom"},
		{"usage → 2", Usagef("missing flag"), 2, "ledgerctl: missing flag"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var se bytes.Buffer
			got := ReportAndExit("ledgerctl", &se, tc.err)
			if got != tc.wantCode {
				t.Errorf("exit code = %d, want %d", got, tc.wantCode)
			}
			if tc.wantInErr != "" && !strings.Contains(se.String(), tc.wantInErr) {
				t.Errorf("stderr = %q, want it to contain %q", se.String(), tc.wantInErr)
			}
			if tc.wantInErr == "" && se.Len() != 0 {
				t.Errorf("stderr should be empty, got %q", se.String())
			}
		})
	}
}

// wrapErr is a trivial wrapper that exercises errors.As-through-Unwrap.
type wrapErr struct{ inner error }

func (w *wrapErr) Error() string { return "wrapped: " + w.inner.Error() }
func (w *wrapErr) Unwrap() error { return w.inner }

// ----- flags.go -----------------------------------------------------------

func TestStringSlice_AppendsAndJoins(t *testing.T) {
	var s StringSlice
	if err := s.Set("a"); err != nil {
		t.Fatal(err)
	}
	if err := s.Set("b"); err != nil {
		t.Fatal(err)
	}
	if got := s.String(); got != "a,b" {
		t.Errorf("String() = %q, want %q", got, "a,b")
	}
	if len(s) != 2 || s[0] != "a" || s[1] != "b" {
		t.Errorf("slice contents wrong: %v", s)
	}
}

func TestParseFlags_SuccessReturnsContinue(t *testing.T) {
	fs := flag.NewFlagSet("sub", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	val := fs.String("foo", "", "")

	stop, err := ParseFlags(fs, []string{"--foo", "bar"}, io.Discard, "tool")
	if err != nil {
		t.Fatalf("ParseFlags: %v", err)
	}
	if stop {
		t.Error("expected stop=false on successful parse")
	}
	if *val != "bar" {
		t.Errorf("flag did not take value: %q", *val)
	}
}

func TestParseFlags_HelpStopsCleanly(t *testing.T) {
	fs := flag.NewFlagSet("sub", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.String("foo", "", "the foo flag")

	var so bytes.Buffer
	stop, err := ParseFlags(fs, []string{"--help"}, &so, "tool")
	if err != nil {
		t.Fatalf("ParseFlags returned err on --help: %v", err)
	}
	if !stop {
		t.Error("expected stop=true on --help")
	}
	if !strings.Contains(so.String(), "Usage of tool sub:") {
		t.Errorf("help text missing usage banner: %q", so.String())
	}
	if !strings.Contains(so.String(), "the foo flag") {
		t.Errorf("help text missing flag description: %q", so.String())
	}
}

func TestParseFlags_BadFlagBecomesUsageError(t *testing.T) {
	fs := flag.NewFlagSet("sub", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	stop, err := ParseFlags(fs, []string{"--no-such-flag"}, io.Discard, "tool")
	if err == nil {
		t.Fatal("expected error on unknown flag")
	}
	if stop {
		t.Error("expected stop=false on parse error")
	}
	if !IsUsage(err) {
		t.Errorf("expected UsageError, got %T: %v", err, err)
	}
}

func TestRequireFlags_MissingProducesSortedUsageError(t *testing.T) {
	err := RequireFlags(map[string]string{
		"alpha":  "",
		"bravo":  "",
		"charlie": "set",
	})
	if !IsUsage(err) {
		t.Fatalf("expected UsageError, got %T: %v", err, err)
	}
	msg := err.Error()
	idxA := strings.Index(msg, "--alpha")
	idxB := strings.Index(msg, "--bravo")
	if idxA < 0 || idxB < 0 {
		t.Fatalf("missing flags not listed: %q", msg)
	}
	if idxA >= idxB {
		t.Errorf("missing flags not sorted: %q", msg)
	}
	if strings.Contains(msg, "--charlie") {
		t.Errorf("set flag listed as missing: %q", msg)
	}
}

func TestRequireFlags_AllSetReturnsNil(t *testing.T) {
	err := RequireFlags(map[string]string{"x": "1", "y": "2"})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

// ----- keys.go ------------------------------------------------------------

func TestGenerateAndLoadPKCS8RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "k.pem")

	pubA, err := GenerateEd25519PKCS8(path)
	if err != nil {
		t.Fatalf("GenerateEd25519PKCS8: %v", err)
	}
	priv, pubB, err := LoadEd25519PKCS8(path)
	if err != nil {
		t.Fatalf("LoadEd25519PKCS8: %v", err)
	}
	if !bytesEq(pubA, pubB) {
		t.Errorf("public keys differ across generate/load")
	}
	// Sign-verify proves the private half loaded matches the public half.
	msg := []byte("smoke test")
	sig := ed25519.Sign(priv, msg)
	if !ed25519.Verify(pubA, msg, sig) {
		t.Error("loaded private key does not match generated public key")
	}
}

func TestLoadPKCS8_RejectsMissingFile(t *testing.T) {
	if _, _, err := LoadEd25519PKCS8(filepath.Join(t.TempDir(), "no.pem")); err == nil {
		t.Fatal("expected error on missing file")
	}
}

func TestLoadPKCS8_RejectsWrongPEMType(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wrong.pem")
	junk := []byte("-----BEGIN CERTIFICATE-----\nQUFB\n-----END CERTIFICATE-----\n")
	if err := os.WriteFile(path, junk, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := LoadEd25519PKCS8(path); err == nil {
		t.Fatal("expected error on wrong PEM type")
	}
}

func TestPublicKeyPEMRoundTrip(t *testing.T) {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "pub.pem")
	if err := WriteEd25519PublicPEM(path, pub); err != nil {
		t.Fatalf("WriteEd25519PublicPEM: %v", err)
	}
	got, err := LoadEd25519PublicPEM(path)
	if err != nil {
		t.Fatalf("LoadEd25519PublicPEM: %v", err)
	}
	if !bytesEq(pub, got) {
		t.Error("public key did not round-trip")
	}
}

// ----- helpers ------------------------------------------------------------

func bytesEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
