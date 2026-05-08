package verify

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// makeFixture is the workhorse for every test below: produce a 5-entry
// ledger that is fully self-consistent and return it alongside the keys.
func makeFixture(t *testing.T) (led *Ledger, sealingPub ed25519.PublicKey, ikm []byte) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ikm = bytes.Repeat([]byte{0xAB}, 32)
	fb := &FixtureBuilder{
		TenantBindingKDFLabel: "tenant=acme;env=prod",
		SealDate:              "2026-06-15",
		MasterIKM:             ikm,
		SealingPriv:           priv,
	}
	var entries []ChainEntry
	prev := ""
	for i := 0; i < 5; i++ {
		payload := []byte(fmt.Sprintf(`{"event":"decision","seq":%d}`, i))
		e, err := fb.AppendEvent(prev, i, fmt.Sprintf("e_%04d", i), payload)
		if err != nil {
			t.Fatal(err)
		}
		entries = append(entries, e)
		prev = e.EntryHash
	}
	seal, err := fb.SealEntries(entries)
	if err != nil {
		t.Fatal(err)
	}
	return &Ledger{Entries: entries, Seal: seal}, pub, ikm
}

// ----- ledger reader -----------------------------------------------------

func TestReadLedger_RoundTripNDJSON(t *testing.T) {
	led, _, _ := makeFixture(t)
	raw, err := MarshalNDJSON(led.Entries, led.Seal)
	if err != nil {
		t.Fatal(err)
	}
	got, err := readLedger(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("readLedger: %v", err)
	}
	if len(got.Entries) != len(led.Entries) {
		t.Errorf("entry count mismatch: got %d, want %d", len(got.Entries), len(led.Entries))
	}
	if got.Seal.MerkleRoot != led.Seal.MerkleRoot {
		t.Errorf("seal MerkleRoot lost in round trip")
	}
}

func TestReadLedger_RejectsEmpty(t *testing.T) {
	if _, err := readLedger(bytes.NewReader(nil)); err == nil {
		t.Fatal("expected error on empty ledger")
	}
}

func TestReadLedger_RejectsNoSeal(t *testing.T) {
	led, _, _ := makeFixture(t)
	raw, _ := MarshalNDJSON(led.Entries, SealRecord{})
	// Strip the empty seal line.
	idx := bytes.LastIndexByte(raw[:len(raw)-1], '\n')
	raw = raw[:idx+1]
	if _, err := readLedger(bytes.NewReader(raw)); err == nil {
		t.Fatal("expected error when seal record missing")
	}
}

func TestReadLedger_RejectsTwoSeals(t *testing.T) {
	led, _, _ := makeFixture(t)
	raw, _ := MarshalNDJSON(led.Entries, led.Seal)
	sealLine, _ := json.Marshal(&led.Seal)
	raw = append(raw, sealLine...)
	raw = append(raw, '\n')
	if _, err := readLedger(bytes.NewReader(raw)); err == nil {
		t.Fatal("expected error on duplicate seal records")
	}
}

func TestReadLedger_RejectsMalformedJSON(t *testing.T) {
	if _, err := readLedger(strings.NewReader("{not json\n")); err == nil {
		t.Fatal("expected error on malformed line")
	}
}

// ----- chain linkage ----------------------------------------------------

func TestCheckChain_PassesOnFixture(t *testing.T) {
	led, _, _ := makeFixture(t)
	if i, err := CheckChain(led.Entries); err != nil {
		t.Fatalf("CheckChain failed at entry %d: %v", i, err)
	}
}

func TestCheckChain_FailsOnTamperedPayload(t *testing.T) {
	led, _, _ := makeFixture(t)
	// Flip the payload of entry 2 — entry_hash will now mismatch.
	led.Entries[2].EventPayloadJCS = base64.StdEncoding.EncodeToString([]byte("tampered"))
	bad, err := CheckChain(led.Entries)
	if err == nil {
		t.Fatal("CheckChain accepted a tampered payload")
	}
	if bad != 2 {
		t.Errorf("bad index = %d, want 2", bad)
	}
}

func TestCheckChain_FailsOnReorderedEntries(t *testing.T) {
	led, _, _ := makeFixture(t)
	led.Entries[1], led.Entries[2] = led.Entries[2], led.Entries[1]
	if _, err := CheckChain(led.Entries); err == nil {
		t.Fatal("CheckChain accepted a reordered ledger")
	}
}

func TestCheckChain_FailsOnTamperedPrevHash(t *testing.T) {
	led, _, _ := makeFixture(t)
	led.Entries[3].PrevHash = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0xFF}, 32))
	bad, err := CheckChain(led.Entries)
	if err == nil {
		t.Fatal("CheckChain accepted a wrong prev_hash")
	}
	if bad != 3 {
		t.Errorf("bad index = %d, want 3", bad)
	}
}

// ----- Merkle ------------------------------------------------------------

func TestCheckMerkleRoot_PassesOnFixture(t *testing.T) {
	led, _, _ := makeFixture(t)
	if err := CheckMerkleRoot(led.Entries, led.Seal.MerkleRoot); err != nil {
		t.Fatalf("CheckMerkleRoot: %v", err)
	}
}

func TestCheckMerkleRoot_FailsOnPayloadFlip(t *testing.T) {
	led, _, _ := makeFixture(t)
	led.Entries[0].EventPayloadJCS = base64.StdEncoding.EncodeToString([]byte("different"))
	if err := CheckMerkleRoot(led.Entries, led.Seal.MerkleRoot); err == nil {
		t.Fatal("CheckMerkleRoot accepted a flipped leaf payload")
	}
}

func TestMerkleTreeHash_RFC6962KnownVectors(t *testing.T) {
	// RFC 6962 §2.1: MTH({}) = SHA256(empty). We capture the property here so
	// any future tweak to the construction trips this test.
	if got := MerkleTreeHash(nil); !bytesEqual(got, sha256Empty()) {
		t.Errorf("MTH({}) does not equal SHA256() — RFC 6962 §2.1")
	}
}

func TestMerkleTreeHash_OneLeafEqualsLeaf(t *testing.T) {
	leaf := MerkleLeafHash([]byte("hello"))
	if got := MerkleTreeHash([][]byte{leaf}); !bytesEqual(got, leaf) {
		t.Errorf("MTH({L}) != L")
	}
}

func TestLargestPowerOfTwoLessThan(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{2, 1}, {3, 2}, {4, 2}, {5, 4}, {7, 4}, {8, 4}, {9, 8}, {16, 8}, {17, 16},
	}
	for _, tc := range cases {
		if got := largestPowerOfTwoLessThan(tc.in); got != tc.want {
			t.Errorf("largestPowerOfTwoLessThan(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

// ----- signature ---------------------------------------------------------

func TestCheckSealSignature_PassesOnFixture(t *testing.T) {
	led, pub, _ := makeFixture(t)
	if err := CheckSealSignature(&led.Seal, pub); err != nil {
		t.Fatalf("CheckSealSignature: %v", err)
	}
}

func TestCheckSealSignature_RejectsWrongKey(t *testing.T) {
	led, _, _ := makeFixture(t)
	other, _, _ := ed25519.GenerateKey(rand.Reader)
	if err := CheckSealSignature(&led.Seal, other); err == nil {
		t.Fatal("CheckSealSignature accepted the wrong public key")
	}
}

func TestCheckSealSignature_RejectsTamperedRoot(t *testing.T) {
	led, pub, _ := makeFixture(t)
	led.Seal.MerkleRoot = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x00}, 32))
	if err := CheckSealSignature(&led.Seal, pub); err == nil {
		t.Fatal("CheckSealSignature accepted a tampered merkle_root")
	}
}

func TestCheckSealSignature_RejectsFingerprintMismatch(t *testing.T) {
	led, pub, _ := makeFixture(t)
	// A correct fingerprint matches pub; flip it to cause mismatch.
	led.Seal.SigningKeyFingerprint = "sha256:" + strings.Repeat("00", 32)
	if err := CheckSealSignature(&led.Seal, pub); err == nil {
		t.Fatal("CheckSealSignature did not catch fingerprint mismatch")
	}
}

func TestCheckSealSignature_RejectsBadKeyLength(t *testing.T) {
	led, _, _ := makeFixture(t)
	if err := CheckSealSignature(&led.Seal, ed25519.PublicKey(make([]byte, 8))); err == nil {
		t.Fatal("CheckSealSignature accepted an 8-byte 'public key'")
	}
}

// ----- per-event MAC -----------------------------------------------------

func TestCheckChainMACs_PassesOnFixture(t *testing.T) {
	led, _, ikm := makeFixture(t)
	if i, err := CheckChainMACs(led.Entries, ikm); err != nil {
		t.Fatalf("CheckChainMACs failed at %d: %v", i, err)
	}
}

func TestCheckChainMACs_RejectsWrongIKM(t *testing.T) {
	led, _, _ := makeFixture(t)
	if _, err := CheckChainMACs(led.Entries, bytes.Repeat([]byte{0x01}, 32)); err == nil {
		t.Fatal("CheckChainMACs accepted the wrong IKM")
	}
}

func TestCheckChainMACs_RejectsEmptyIKM(t *testing.T) {
	led, _, _ := makeFixture(t)
	if _, err := CheckChainMACs(led.Entries, nil); err == nil {
		t.Fatal("CheckChainMACs accepted empty IKM")
	}
}

func TestHKDFExpand_OutputLengthAndDeterminism(t *testing.T) {
	a := hkdf32([]byte("ikm"), []byte("salt"), []byte("info"))
	b := hkdf32([]byte("ikm"), []byte("salt"), []byte("info"))
	if !bytesEqual(a, b) {
		t.Error("hkdf32 is non-deterministic for fixed inputs")
	}
	if len(a) != 32 {
		t.Errorf("hkdf32 length = %d, want 32", len(a))
	}
	if c := hkdf32([]byte("ikm"), []byte("salt"), []byte("DIFFERENT")); bytesEqual(a, c) {
		t.Error("hkdf32 produced equal outputs for distinct info")
	}
}

// ----- pipeline ---------------------------------------------------------

func TestVerify_StructuralOnlyWhenNoIKM(t *testing.T) {
	led, pub, _ := makeFixture(t)
	r, err := Verify(led, Plan{SealingKey: pub, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !r.StructuralPass {
		t.Error("StructuralPass false on a clean ledger")
	}
	if r.MACPass {
		t.Error("MACPass true even though no IKM was supplied")
	}
	last := r.Steps[len(r.Steps)-1]
	if last.Name != "per-event-mac" || last.OK {
		t.Errorf("expected per-event-mac step recorded as skipped, got %+v", last)
	}
	if !strings.Contains(last.Note, "structural-only") {
		t.Errorf("skip note missing 'structural-only': %q", last.Note)
	}
}

func TestVerify_FullPathPassesWithIKM(t *testing.T) {
	led, pub, ikm := makeFixture(t)
	r, err := Verify(led, Plan{SealingKey: pub, MasterIKM: ikm, StopOnFirstFailure: true})
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !r.StructuralPass || !r.MACPass {
		t.Errorf("expected both StructuralPass and MACPass, got %+v", r)
	}
}

func TestVerify_FailsLoudOnTamper(t *testing.T) {
	led, pub, _ := makeFixture(t)
	led.Entries[2].EventPayloadJCS = base64.StdEncoding.EncodeToString([]byte("tampered"))
	r, err := Verify(led, Plan{SealingKey: pub, StopOnFirstFailure: true})
	if err == nil {
		t.Fatal("Verify accepted a tampered ledger")
	}
	if r == nil || r.StructuralPass {
		t.Error("StructuralPass should be false on tamper")
	}
}

func TestVerify_RequiresSealingKey(t *testing.T) {
	led, _, _ := makeFixture(t)
	if _, err := Verify(led, Plan{}); err == nil {
		t.Fatal("Verify accepted plan with no SealingKey")
	}
}

// ----- helpers ----------------------------------------------------------

func sha256Empty() []byte {
	h := sha256SumEmpty
	return h[:]
}

// Pre-computed SHA-256 of the empty input. Inlined as a constant so the
// test's intent stays readable.
var sha256SumEmpty = [32]byte{
	0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14,
	0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24,
	0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c,
	0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55,
}
