package verify

import "encoding/hex"

// AuditPathStep is one step of an RFC 6962 inclusion proof: the sibling
// hash to combine with the running node, and which side the sibling sits
// on. IsLeft=true means the sibling is on the LEFT (combine as
// H(0x01 || sibling || current)); IsLeft=false means the sibling is on
// the RIGHT (combine as H(0x01 || current || sibling)). This matches the
// §10.31 is_left_semantics the corpus pins.
type AuditPathStep struct {
	Sibling []byte
	IsLeft  bool
}

// FoldAuditPath folds a leaf hash up an RFC 6962 audit path and returns
// the resulting root. The fold is the §10.31 inclusion-proof verifier:
// at each step it combines the running node with the step's sibling on
// the correct side using MerkleNodeHash (the 0x01-prefixed internal-node
// hash), so a clean-room verifier reproduces the same root from the same
// path.
//
// The §10.37 hierarchical case (case 026) folds a CONCATENATED path —
// the inner subtree path followed by the outer top-tree path — through
// this same fold unchanged, because both levels use the same
// (sibling, is_left) step shape. That is why this is one function, not
// two: the spec's hierarchical note guarantees a §10.31 verifier
// verifies the concatenated path without special-casing.
func FoldAuditPath(leafHash []byte, path []AuditPathStep) []byte {
	current := leafHash
	for _, step := range path {
		if step.IsLeft {
			current = MerkleNodeHash(step.Sibling, current)
		} else {
			current = MerkleNodeHash(current, step.Sibling)
		}
	}
	return current
}

// VerifyAuditPath folds leafHash up path and reports whether the result
// equals root. All three arguments are raw bytes (not hex); the caller
// decodes the corpus's hex pins before calling. Returns true on a
// byte-identical match.
func VerifyAuditPath(leafHash []byte, path []AuditPathStep, root []byte) bool {
	got := FoldAuditPath(leafHash, path)
	return bytesEqual(got, root)
}

// VerifyAuditPathHex is the hex-string convenience wrapper for callers
// holding the corpus's hex pins directly. Returns false (not an error)
// on any decode failure — a malformed proof does not verify, which is
// the same disposition a tampered proof gets.
func VerifyAuditPathHex(leafHashHex string, path []AuditPathStep, rootHex string) bool {
	leaf, err := hex.DecodeString(leafHashHex)
	if err != nil {
		return false
	}
	root, err := hex.DecodeString(rootHex)
	if err != nil {
		return false
	}
	return VerifyAuditPath(leaf, path, root)
}
