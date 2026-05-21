package hkdf

import (
	"encoding/hex"
	"testing"
)

// TestDerive_RFC5869_TestCase1 pins the canonical RFC 5869 Appendix A
// Test Case 1 vector. SHA-256, basic inputs. If this drifts, the
// HKDF implementation is broken at the foundation and every
// downstream conformance check is meaningless.
//
// Source: https://www.rfc-editor.org/rfc/rfc5869#appendix-A.1
func TestDerive_RFC5869_TestCase1(t *testing.T) {
	ikm, _ := hex.DecodeString("0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b")
	salt, _ := hex.DecodeString("000102030405060708090a0b0c")
	info, _ := hex.DecodeString("f0f1f2f3f4f5f6f7f8f9")
	wantOKM, _ := hex.DecodeString("3cb25f25faacd57a90434f64d0362f2a2d2d0a90cf1a5a4c5db02d56ecc4c5bf34007208d5b887185865")

	got := Derive(ikm, salt, info, 42)
	if hex.EncodeToString(got) != hex.EncodeToString(wantOKM) {
		t.Errorf("HKDF RFC 5869 test case 1 mismatch:\n got: %x\nwant: %x", got, wantOKM)
	}
}

// TestDerive_RFC5869_TestCase3 covers the empty-salt path per RFC
// 5869 §2.2 — an empty salt MUST be replaced with HashLen zero
// bytes.
//
// Source: https://www.rfc-editor.org/rfc/rfc5869#appendix-A.3
func TestDerive_RFC5869_TestCase3(t *testing.T) {
	ikm, _ := hex.DecodeString("0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b")
	wantOKM, _ := hex.DecodeString("8da4e775a563c18f715f802a063c5a31b8a11f5c5ee1879ec3454e5f3c738d2d9d201395faa4b61a96c8")

	got := Derive(ikm, nil, nil, 42)
	if hex.EncodeToString(got) != hex.EncodeToString(wantOKM) {
		t.Errorf("HKDF RFC 5869 test case 3 (empty-salt path) mismatch:\n got: %x\nwant: %x", got, wantOKM)
	}
}

func TestDerive_LengthZeroReturnsNil(t *testing.T) {
	if got := Derive([]byte("ikm"), nil, nil, 0); got != nil {
		t.Errorf("Derive(..., length=0) = %x, want nil", got)
	}
}

func TestDerive_NegativeLengthReturnsNil(t *testing.T) {
	if got := Derive([]byte("ikm"), nil, nil, -1); got != nil {
		t.Errorf("Derive(..., length=-1) = %x, want nil", got)
	}
}

func TestDerive_Determinism(t *testing.T) {
	a := Derive([]byte("ikm"), []byte("salt"), []byte("info"), 32)
	b := Derive([]byte("ikm"), []byte("salt"), []byte("info"), 32)
	if hex.EncodeToString(a) != hex.EncodeToString(b) {
		t.Errorf("Derive not deterministic for identical inputs")
	}
}
