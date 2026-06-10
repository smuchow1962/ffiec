package vectors

import (
	"encoding/hex"
	"fmt"

	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

// disclosure is the shared shape of one audit-path disclosure across the
// 023 (flat §10.31) and 026 (hierarchical §10.37) vectors. 023 carries
// `audit_path` + `merkle_root_hex` at the tree level; 026 carries
// `concatenated_path` + a top-level `top_root_hex`. Both fold through
// the same verify.FoldAuditPath, so the decode normalizes to this shape.
type disclosure struct {
	LeafLabel   string
	LeafHashHex string
	Path        []verify.AuditPathStep
}

// pathStepJSON is one audit-path step as the corpus pins it.
type pathStepJSON struct {
	SiblingHex string `json:"sibling_hex"`
	IsLeft     bool   `json:"is_left"`
}

func decodePath(steps []pathStepJSON) ([]verify.AuditPathStep, error) {
	out := make([]verify.AuditPathStep, 0, len(steps))
	for i, s := range steps {
		sib, err := hex.DecodeString(s.SiblingHex)
		if err != nil {
			return nil, fmt.Errorf("decode sibling at step %d: %w", i, err)
		}
		out = append(out, verify.AuditPathStep{Sibling: sib, IsLeft: s.IsLeft})
	}
	return out, nil
}

// ---- 023: RFC 6962 inclusion proof (§10.31) ----

func runInclusionProof(v RichVector) *Report {
	r := &Report{}

	var exp struct {
		Tree struct {
			MerkleRootHex string `json:"merkle_root_hex"`
		} `json:"tree"`
		Disclosures []struct {
			LeafIndex   int            `json:"leaf_index"`
			LeafHashHex string         `json:"leaf_hash_hex"`
			AuditPath   []pathStepJSON `json:"audit_path"`
		} `json:"disclosures"`
	}
	if err := readVectorJSON(v.Dir, "expected.json", &exp); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	for _, d := range exp.Disclosures {
		name := fmt.Sprintf("%s/leaf-%d-folds-to-root", v.Slot, d.LeafIndex)
		path, err := decodePath(d.AuditPath)
		if err != nil {
			r.Checks = append(r.Checks, failCheck(name, "%v", err))
			continue
		}
		if !verify.VerifyAuditPathHex(d.LeafHashHex, path, exp.Tree.MerkleRootHex) {
			r.Checks = append(r.Checks, failCheck(name, "leaf %s did not fold to root %s", d.LeafHashHex, exp.Tree.MerkleRootHex))
			continue
		}
		r.Checks = append(r.Checks, passCheck(name))
	}
	return r
}

// ---- 026: hierarchical Merkle aggregation (§10.37) ----
//
// Two proofs per disclosure: (a) the subtree roots build the top root via
// MerkleNodeHash DIRECTLY (no leaf-hash wrap on subtree roots — the
// spec's top_level_no_leaf_hash_note), and (b) each leaf's concatenated
// path (inner subtree path ++ outer top-tree path) folds to the top root
// through the same FoldAuditPath the flat §10.31 case uses.

func runHierarchicalMerkle(v RichVector) *Report {
	r := &Report{}

	var exp struct {
		Subtrees []struct {
			SubtreeRootHex string `json:"subtree_root_hex"`
		} `json:"subtrees"`
		TopRootHex  string `json:"top_root_hex"`
		Disclosures []struct {
			SubtreeIndex     int            `json:"subtree_index"`
			LeafIndex        int            `json:"leaf_index"`
			LeafHashHex      string         `json:"leaf_hash_hex"`
			ConcatenatedPath []pathStepJSON `json:"concatenated_path"`
		} `json:"disclosures"`
	}
	if err := readVectorJSON(v.Dir, "expected.json", &exp); err != nil {
		r.Checks = append(r.Checks, failCheck(v.Slot+"/read", "%v", err))
		return r
	}

	r.Checks = append(r.Checks, checkTopRootFromSubtreeRoots(v.Slot, exp.Subtrees, exp.TopRootHex))

	for _, d := range exp.Disclosures {
		name := fmt.Sprintf("%s/subtree-%d-leaf-%d-folds-to-top-root", v.Slot, d.SubtreeIndex, d.LeafIndex)
		path, err := decodePath(d.ConcatenatedPath)
		if err != nil {
			r.Checks = append(r.Checks, failCheck(name, "%v", err))
			continue
		}
		if !verify.VerifyAuditPathHex(d.LeafHashHex, path, exp.TopRootHex) {
			r.Checks = append(r.Checks, failCheck(name, "leaf %s did not fold to top root %s", d.LeafHashHex, exp.TopRootHex))
			continue
		}
		r.Checks = append(r.Checks, passCheck(name))
	}
	return r
}

// checkTopRootFromSubtreeRoots builds the §10.37 outer tree over the
// subtree roots, combining them DIRECTLY via MerkleNodeHash (no
// leaf-hash wrap — subtree roots are already SHA-256 outputs). Asserts
// the result equals top_root_hex.
func checkTopRootFromSubtreeRoots(slot string, subtrees []struct {
	SubtreeRootHex string `json:"subtree_root_hex"`
}, wantTopRootHex string) Check {
	name := slot + "/top-root-from-subtree-roots"
	roots := make([][]byte, 0, len(subtrees))
	for i, s := range subtrees {
		b, err := hex.DecodeString(s.SubtreeRootHex)
		if err != nil {
			return failCheck(name, "decode subtree root %d: %v", i, err)
		}
		roots = append(roots, b)
	}
	// MerkleTreeHash treats a 1-element slice as the element itself and
	// combines via MerkleNodeHash (0x01-prefixed) — exactly the outer
	// tree's no-leaf-hash discipline, since the inputs are already roots.
	got := hex.EncodeToString(verify.MerkleTreeHash(roots))
	if got != wantTopRootHex {
		return failCheck(name, "top root mismatch: got=%s want=%s", got, wantTopRootHex)
	}
	return passCheck(name)
}
