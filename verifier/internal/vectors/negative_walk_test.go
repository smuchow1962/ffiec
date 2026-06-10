package vectors

import (
	"errors"
	"strings"
	"testing"
)

// expectedClass pins the wave-3 classification (the plan-doc table) so a
// regression that silently reclassifies a live vector as contract-only —
// the way a verifier bug could quietly de-arm the gate — is caught. The
// 20 live-walk slots assert the verifier's own behavior; the rest are
// contract-only with a recorded spec-section reason.
var expectedClass = map[string]negativeClass{
	// Base §7 1-10 walk (18).
	"N001": classBaseWalk, "N002": classBaseWalk, "N003": classBaseWalk,
	"N006": classBaseWalk, "N007": classBaseWalk, "N008": classBaseWalk,
	"N009": classBaseWalk, "N010": classBaseWalk, "N011": classBaseWalk,
	"N012": classBaseWalk, "N013": classBaseWalk, "N014": classBaseWalk,
	"N015": classBaseWalk, "N016": classBaseWalk, "N022": classBaseWalk,
	"N023": classBaseWalk, "N030": classBaseWalk, "N033": classBaseWalk,
	// Structural §7-step-11 algorithm/key-type compare (1).
	"N020": classAlgKeyType,
	// §10.42 backfill root recompute (1).
	"N025": classBackfillRoot,
	// §7 step-11 live Ed25519 signature walk (2). The DRIVER decides per-
	// fixture whether to run the live walk or fall back to contract-only
	// (live-readiness gate in negative_signature.go); the CLASS is always
	// signature.
	"N004": classSealSignature, "N005": classSealSignature,
	// Contract-only (16).
	"N017": classContractOnly, "N018": classContractOnly, "N019": classContractOnly,
	"N021": classContractOnly, "N024": classContractOnly, "N026": classContractOnly,
	"N027": classContractOnly, "N028": classContractOnly, "N029": classContractOnly,
	"N031": classContractOnly, "N032": classContractOnly, "N034": classContractOnly,
	"N035": classContractOnly, "N036": classContractOnly, "N037": classContractOnly,
	"N038": classContractOnly,
}

// TestNegativeClassification_PinsTheTable guards the live-vs-contract-only
// map against silent drift. A verifier change that reclassified a live
// vector would de-arm a live assertion without failing any other test;
// this catches it.
func TestNegativeClassification_PinsTheTable(t *testing.T) {
	for slot, want := range expectedClass {
		if got := classifyNegative(slot); got != want {
			t.Errorf("classifyNegative(%s) = %d, want %d", slot, got, want)
		}
	}
}

// TestNegativeLiveVectors_RunLivePath proves the 20 live-walk vectors are
// asserted by the live driver (the §7 walk / §10.42 recompute / §7-step-11
// structural compare) and NOT silently via the contract-only fallback.
// The proof is the Check name: a live assertion carries the live-path
// suffix; the contract-only path carries "/expected-output". This is the
// honesty guard — a green gate alone does not prove the vectors are live.
func TestNegativeLiveVectors_RunLivePath(t *testing.T) {
	dir, err := CorpusDir()
	if errors.Is(err, ErrCorpusNotFound) {
		t.Skipf("negative corpus not available: %v", err)
	}
	if err != nil {
		t.Fatalf("CorpusDir: %v", err)
	}

	vectors, err := DiscoverNegativeVectors(dir)
	if err != nil {
		t.Fatalf("DiscoverNegativeVectors: %v", err)
	}

	// Live-path suffixes split into the always-live classes (base walk /
	// structural step-11 / backfill recompute) and the signature class,
	// whose per-fixture liveness depends on Heather's materialization.
	alwaysLiveSuffixes := []string{"/§7-walk", "/§7-step-11-structural", "/§10.42-root-recompute"}
	const signatureSuffix = "/§7-step-11-signature"

	alwaysLiveCount, signatureLiveCount, contractCount := 0, 0, 0

	for _, v := range vectors {
		class := classifyNegative(v.Slot)
		if class == classContractOnly {
			contractCount++
			continue
		}
		result := RunNegativeVector(v)
		if result.Skipped {
			t.Errorf("%s classified live but SKIPPED (not materialized?)", v.Slot)
			continue
		}
		if !result.Report.Pass() {
			t.Errorf("%s live assertion failed", v.Slot)
		}

		// Signature vectors run live OR fall back to contract-only based on
		// the live-readiness gate; both outcomes pass. The other live
		// classes MUST carry an always-live suffix.
		if class == classSealSignature {
			if anyCheckNameHasSuffix(result.Report.Checks, []string{signatureSuffix}) {
				signatureLiveCount++
			} else {
				contractCount++ // fell back to /expected-output
			}
			continue
		}
		if !anyCheckNameHasSuffix(result.Report.Checks, alwaysLiveSuffixes) {
			t.Errorf("%s classified live but its checks carry no live-path name: %v",
				v.Slot, checkNames(result.Report.Checks))
		}
		alwaysLiveCount++
	}

	// The base/structural/backfill live count is fixed at 20 (18 base + 1
	// alg-key-type + 1 backfill). The signature count is whatever is
	// materialized — at minimum N004 (garbage sig, decodable), and N005
	// once its real wrong-tenant signature lands.
	const wantAlwaysLive = 20
	if alwaysLiveCount != wantAlwaysLive {
		t.Errorf("ran %d always-live vectors, expected %d", alwaysLiveCount, wantAlwaysLive)
	}
	t.Logf("live-walk: %d always-live (§7/§10.42/step-11-structural) + %d step-11-signature live; contract-only: %d",
		alwaysLiveCount, signatureLiveCount, contractCount)
}

func anyCheckNameHasSuffix(checks []Check, suffixes []string) bool {
	for _, c := range checks {
		for _, s := range suffixes {
			if strings.HasSuffix(c.Name, s) {
				return true
			}
		}
	}
	return false
}

func checkNames(checks []Check) []string {
	names := make([]string, len(checks))
	for i, c := range checks {
		names[i] = c.Name
	}
	return names
}
