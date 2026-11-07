package vectors

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// The posture decision against the bytes the corpus pins. N017 and N018 stay contract-only in
// the walk, because their signatures are placeholders. The decision reads no signature, so
// these 2 vectors can judge it directly: input.json gives the seal and the outside posture,
// and expected_output.txt gives the result.
func TestPostureDecisionMatchesThePinnedVectors(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skip("corpus not found")
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		slot   string
		strict bool
	}{
		{"N017-dual-algo-partial-coverage", false},
		{"N018-dual-algo-not-in-posture", true},
	} {
		t.Run(tc.slot, func(t *testing.T) {
			vector := filepath.Join(dir, "negative", tc.slot)
			want, err := parseExpectedOutput(vector)
			if err != nil {
				t.Fatal(err)
			}
			seal, posture := postureInput(t, vector)
			got := verify.JudgeSealPosture(seal, verify.SealPolicy{Posture: posture, Strict: tc.strict})
			if got.Status != want.Status || got.Step != want.Step || got.Reason != want.Reason || got.ExitCode != want.ExitCode {
				t.Errorf("got %+v, pinned %+v", got, want)
			}
			if !strings.Contains(strings.Join(got.Anomalies, "\n"), want.Anomaly) {
				t.Errorf("anomalies %q lack the pinned %q", got.Anomalies, want.Anomaly)
			}
		})
	}
}

func postureInput(t *testing.T, vector string) (verify.AuditSeal, []string) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(vector, "input.json"))
	if err != nil {
		t.Fatal(err)
	}
	var input struct {
		Tamper struct {
			Declared []string `json:"declared"`
		} `json:"tamper"`
		AuditFile struct {
			Seal json.RawMessage `json:"seal"`
		} `json:"audit_file"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		t.Fatal(err)
	}
	seal, loaded := verify.LoadAuditSeal(input.AuditFile.Seal)
	if loaded.Status == "FAIL" {
		t.Fatal(loaded.Reason)
	}
	return seal, input.Tamper.Declared
}
