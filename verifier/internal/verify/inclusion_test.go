package verify

import (
	"encoding/hex"
	"testing"
)

// TestFoldAuditPath_RFC6962_5Leaf round-trips leaf 0 and leaf 4 of the
// 023 5-leaf tree to the pinned root. Leaf 0 sits in the 4-leaf left
// subtree (3-step path); leaf 4 is the lone right-subtree leaf (1-step
// path). Both must fold to the same root, exercising the n=5 → k=4 split.
func TestFoldAuditPath_RFC6962_5Leaf(t *testing.T) {
	const root = "12182c226015d4b8e19ac69184bf8563c21b5413fd3b2553c37b5b5a7a04cbe0"

	tests := []struct {
		name     string
		leafHash string
		path     []hexStep
	}{
		{
			name:     "leaf-0",
			leafHash: "51e216412f613ba7855183d09da1c153a8b1da2a64afa963989bbe1f89e2563e",
			path: []hexStep{
				{"f173c6da9148d721b0adbecbe5753e082ad856a5c386153a8b85ed4d20f3f0e4", false},
				{"051db98292681c93c011eebe90be9f15fe1071a6b7d6b77640ee69b51eeebbb3", false},
				{"999c377d3b9e9ad37a989e351e0269676eac2bb76a7ee0e26196500094c7658a", false},
			},
		},
		{
			name:     "leaf-4-lone-right",
			leafHash: "999c377d3b9e9ad37a989e351e0269676eac2bb76a7ee0e26196500094c7658a",
			path: []hexStep{
				{"539fcb6ef3d7df93496e6fde4d54b106ecbf6366650575714bd33b3fc3e0cb9f", true},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !VerifyAuditPathHex(tc.leafHash, decodeSteps(t, tc.path), root) {
				t.Errorf("leaf %s did not fold to root %s", tc.leafHash, root)
			}
		})
	}
}

// TestFoldAuditPath_Tampered confirms a flipped is_left does NOT verify —
// a proof that folds to a different root must fail, which is the same
// disposition a tampered proof gets.
func TestFoldAuditPath_Tampered(t *testing.T) {
	const root = "12182c226015d4b8e19ac69184bf8563c21b5413fd3b2553c37b5b5a7a04cbe0"
	leaf := "51e216412f613ba7855183d09da1c153a8b1da2a64afa963989bbe1f89e2563e"
	// Same path as leaf-0 but with the first step's side flipped.
	path := []hexStep{
		{"f173c6da9148d721b0adbecbe5753e082ad856a5c386153a8b85ed4d20f3f0e4", true},
		{"051db98292681c93c011eebe90be9f15fe1071a6b7d6b77640ee69b51eeebbb3", false},
		{"999c377d3b9e9ad37a989e351e0269676eac2bb76a7ee0e26196500094c7658a", false},
	}
	if VerifyAuditPathHex(leaf, decodeSteps(t, path), root) {
		t.Error("tampered path (flipped is_left) unexpectedly verified")
	}
}

type hexStep struct {
	siblingHex string
	isLeft     bool
}

func decodeSteps(t *testing.T, steps []hexStep) []AuditPathStep {
	t.Helper()
	out := make([]AuditPathStep, 0, len(steps))
	for _, s := range steps {
		sib, err := hex.DecodeString(s.siblingHex)
		if err != nil {
			t.Fatalf("decode sibling %s: %v", s.siblingHex, err)
		}
		out = append(out, AuditPathStep{Sibling: sib, IsLeft: s.isLeft})
	}
	return out
}
