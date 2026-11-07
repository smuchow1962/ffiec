package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/testkit/clitest"
)

// Richard's probe13e inputs. A ledger's fields have no signature over these lines, so a line break
// could print a false "overall: PASS" above the true result, and an ESC byte could reach a terminal.
const forgedOverall = "x\x1b[2J\n  overall:       PASS (full, key-bound)"

type ledgerFixture struct {
	path, pub string
	rows      []map[string]any
}

func genLedger(t *testing.T) ledgerFixture {
	t.Helper()
	dir := t.TempDir()
	path, pub := filepath.Join(dir, "ledger.ndjson"), filepath.Join(dir, "pub.pem")
	clitest.Run(Main, "gen-fixture", "--out-ledger", path, "--out-pub", pub, "--entries", "4").MustSucceed(t)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	f := ledgerFixture{path: path, pub: pub}
	for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte("\n")) {
		var row map[string]any
		if err := json.Unmarshal(line, &row); err != nil {
			t.Fatal(err)
		}
		f.rows = append(f.rows, row)
	}
	return f
}

// edited writes a copy with one field changed on the seal row, or on the first entry row.
func (f ledgerFixture) edited(t *testing.T, onSeal bool, field string, value any) string {
	t.Helper()
	var out bytes.Buffer
	done := false
	for _, row := range f.rows {
		copied := map[string]any{}
		for k, v := range row {
			copied[k] = v
		}
		if !done && (row["type"] == "seal") == onSeal {
			copied[field] = value
			done = true
		}
		line, err := json.Marshal(copied)
		if err != nil {
			t.Fatal(err)
		}
		out.Write(line)
		out.WriteByte('\n')
	}
	path := filepath.Join(t.TempDir(), "edited.ndjson")
	if err := os.WriteFile(path, out.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func assertNoForgedOutput(t *testing.T, stdout, stderr string) {
	t.Helper()
	overall := 0
	for _, line := range strings.Split(strings.ReplaceAll(stdout, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "  overall:") {
			overall++
		}
	}
	if overall > 1 {
		t.Errorf("%d overall lines:\n%s", overall, stdout)
	}
	if strings.ContainsRune(stdout, 0x1b) || strings.ContainsRune(stderr, 0x1b) {
		t.Errorf("an ESC byte was printed: stdout %q stderr %q", stdout, stderr)
	}
}

func TestLedgerOutput_AFieldCannotForgeALine(t *testing.T) {
	f := genLedger(t)
	for _, tc := range []struct {
		name   string
		onSeal bool
		field  string
	}{
		{"seal fingerprint", true, "signing_key_fingerprint"},
		{"seal date", true, "seal_date"},
		{"seal merkle root", true, "merkle_root"},
		{"entry id", false, "entry_id"},
		{"entry seal date", false, "seal_date"},
		{"entry kdf label", false, "tenant_binding_kdf_label"},
		{"entry prev hash", false, "prev_hash"},
		{"entry hash", false, "entry_hash"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			file := f.edited(t, tc.onSeal, tc.field, forgedOverall)
			for name, args := range map[string][]string{
				"verify": {"verify", "--root-key", f.pub, "--ledger", file},
				"walk":   {"walk", "--ledger", file},
				"diff":   {"diff", "--before", f.path, "--after", file},
			} {
				t.Run(name, func(t *testing.T) {
					result := clitest.Run(Main, args...)
					assertNoForgedOutput(t, result.Stdout, result.Stderr)
				})
			}
		})
	}
}

// The --ledger walk applies no posture and no floor, so it refuses these flags. An empty value must
// not slip past that refusal and read as a run that applied them.
func TestLedger_RefusesPolicyFlagsEvenWhenEmpty(t *testing.T) {
	f := genLedger(t)
	for _, flag := range []string{"--posture=", "--min-sign-payload-version=", "--posture=ed25519", "--strict"} {
		t.Run(flag, func(t *testing.T) {
			clitest.Run(Main, "verify", "--root-key", f.pub, "--ledger", f.path, flag).MustExit(t, 2)
		})
	}
}

// Richard's probe. diff said "ledgers identical" at exit 0 over edits it did not compare.
func TestDiff_IdenticalMeansEveryByte(t *testing.T) {
	f := genLedger(t)
	clitest.Run(Main, "diff", "--before", f.path, "--after", f.path).MustSucceed(t).StdoutContains(t, "ledgers identical")
	for _, tc := range []struct {
		onSeal bool
		field  string
	}{
		{true, "signing_key_fingerprint"},
		{false, "tenant_binding_kdf_label"},
	} {
		t.Run(tc.field, func(t *testing.T) {
			result := clitest.Run(Main, "diff", "--before", f.path, "--after", f.edited(t, tc.onSeal, tc.field, "edited"))
			if result.ExitCode == 0 || strings.Contains(result.Stdout, "identical") {
				t.Fatalf("exit %d: %s", result.ExitCode, result.Stdout)
			}
		})
	}
}
