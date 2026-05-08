package verify

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// RFC 6962 leaf and node prefixes. Domain-separating leaves from internal
// nodes prevents a forged tree from substituting a node hash for a leaf
// hash; the prefixes are checked into the spec, not parameters.
const (
	merkleLeafPrefix byte = 0x00
	merkleNodePrefix byte = 0x01
)

// MerkleLeafHash returns SHA-256(0x00 || payload) per RFC 6962 §2.1.
func MerkleLeafHash(payload []byte) []byte {
	h := sha256.New()
	h.Write([]byte{merkleLeafPrefix})
	h.Write(payload)
	return h.Sum(nil)
}

// MerkleNodeHash returns SHA-256(0x01 || left || right) per RFC 6962 §2.1.
func MerkleNodeHash(left, right []byte) []byte {
	h := sha256.New()
	h.Write([]byte{merkleNodePrefix})
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

// MerkleTreeHash returns the root of a Merkle tree built over leaves per
// RFC 6962 §2.1. Empty input returns SHA-256 of the empty string. Odd
// counts are handled by the natural largest-power-of-two split — leaves
// are not duplicated.
func MerkleTreeHash(leaves [][]byte) []byte {
	switch len(leaves) {
	case 0:
		// RFC 6962: MTH({}) = SHA256().
		h := sha256.Sum256(nil)
		return h[:]
	case 1:
		return leaves[0]
	}
	k := largestPowerOfTwoLessThan(len(leaves))
	left := MerkleTreeHash(leaves[:k])
	right := MerkleTreeHash(leaves[k:])
	return MerkleNodeHash(left, right)
}

// largestPowerOfTwoLessThan returns the largest power of two strictly less
// than n, for n >= 2. Used by RFC 6962 to split a non-leaf range.
func largestPowerOfTwoLessThan(n int) int {
	k := 1
	for k*2 < n {
		k *= 2
	}
	return k
}

// CheckMerkleRoot recomputes the Merkle root from entries' event-payload
// JCS bytes (in the order entries appear) and compares to claimedRoot.
//
// The function only needs the payload bytes — the chain linkage is a
// separate invariant validated by CheckChain. This keeps the two failure
// modes legible at the verifier's report level.
func CheckMerkleRoot(entries []ChainEntry, claimedRoot string) error {
	leaves := make([][]byte, 0, len(entries))
	for i := range entries {
		payload, err := decodeBase64(entries[i].EventPayloadJCS, "event_payload_jcs")
		if err != nil {
			return fmt.Errorf("entry %d (%s): %w", i, entries[i].EntryID, err)
		}
		leaves = append(leaves, MerkleLeafHash(payload))
	}
	want, err := decodeBase64(claimedRoot, "merkle_root")
	if err != nil {
		return err
	}
	got := MerkleTreeHash(leaves)
	if !bytesEqual(got, want) {
		return fmt.Errorf("merkle root does not match (recomputed=%s, claimed=%s)",
			base64.StdEncoding.EncodeToString(got), claimedRoot)
	}
	return nil
}
