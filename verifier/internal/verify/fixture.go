package verify

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// FixtureBuilder constructs a self-consistent ledger for tests, demos, and
// the verifier's smoke tests. Real ledgers come out of the ledger server;
// FixtureBuilder mirrors the same chain-construction rules so verifier
// tests can produce inputs without bringing the server in.
type FixtureBuilder struct {
	TenantBindingKDFLabel string
	SealDate              string

	// MasterIKM is required so per-event MACs are computable.
	MasterIKM []byte
	// SealingPriv signs the daily Merkle root. Its public half is what the
	// examiner's verifier loads via --root-key.
	SealingPriv ed25519.PrivateKey
}

// AppendEvent computes one chain entry over payload (the canonical event
// bytes) and returns it with all per-entry fields populated. The caller
// holds the slice of entries; FixtureBuilder is stateless beyond config.
func (fb *FixtureBuilder) AppendEvent(prevEntryHash string, index int, entryID string, payload []byte) (ChainEntry, error) {
	if len(fb.MasterIKM) == 0 {
		return ChainEntry{}, fmt.Errorf("FixtureBuilder.MasterIKM is required")
	}
	mac := hmacSHA256(hkdf32(fb.MasterIKM, []byte(fb.TenantBindingKDFLabel), []byte(entryID)), payload)
	e := ChainEntry{
		EntryID:               entryID,
		PrevHash:              prevEntryHash,
		HMACSHA256:            base64.StdEncoding.EncodeToString(mac),
		TenantBindingKDFLabel: fb.TenantBindingKDFLabel,
		EventPayloadJCS:       base64.StdEncoding.EncodeToString(payload),
		MerkleLeafIndex:       index,
		SealDate:              fb.SealDate,
	}
	canon, err := canonicalForEntryHash(&e)
	if err != nil {
		return ChainEntry{}, err
	}
	sum := sha256.Sum256(canon)
	e.EntryHash = base64.StdEncoding.EncodeToString(sum[:])
	return e, nil
}

// SealEntries computes the Merkle root over entries' payloads and signs the
// resulting SealRecord with fb.SealingPriv.
func (fb *FixtureBuilder) SealEntries(entries []ChainEntry) (SealRecord, error) {
	if len(fb.SealingPriv) != ed25519.PrivateKeySize {
		return SealRecord{}, fmt.Errorf("FixtureBuilder.SealingPriv has wrong size: %d", len(fb.SealingPriv))
	}
	leaves := make([][]byte, 0, len(entries))
	for i := range entries {
		payload, err := decodeBase64(entries[i].EventPayloadJCS, "event_payload_jcs")
		if err != nil {
			return SealRecord{}, err
		}
		leaves = append(leaves, MerkleLeafHash(payload))
	}
	root := MerkleTreeHash(leaves)

	pub := fb.SealingPriv.Public().(ed25519.PublicKey)
	pubSum := sha256.Sum256(pub)
	s := SealRecord{
		Type:                  "seal",
		SealDate:              fb.SealDate,
		MerkleRoot:            base64.StdEncoding.EncodeToString(root),
		LeafCount:             len(entries),
		SigningKeyFingerprint: "sha256:" + hex.EncodeToString(pubSum[:]),
	}
	msg, err := canonicalForSealSignature(&s)
	if err != nil {
		return SealRecord{}, err
	}
	s.SignatureEd25519 = base64.StdEncoding.EncodeToString(ed25519.Sign(fb.SealingPriv, msg))
	return s, nil
}

// MarshalNDJSON serialises entries followed by the seal as one
// newline-delimited JSON document — the bootstrap on-disk ledger format.
func MarshalNDJSON(entries []ChainEntry, seal SealRecord) ([]byte, error) {
	var sb strings.Builder
	for i := range entries {
		raw, err := json.Marshal(&entries[i])
		if err != nil {
			return nil, err
		}
		sb.Write(raw)
		sb.WriteByte('\n')
	}
	raw, err := json.Marshal(&seal)
	if err != nil {
		return nil, err
	}
	sb.Write(raw)
	sb.WriteByte('\n')
	return []byte(sb.String()), nil
}
