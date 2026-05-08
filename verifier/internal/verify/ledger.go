// Package verify implements the offline chain-of-custody verification
// procedure: chain linkage, RFC 6962 Merkle proof reconstruction, Ed25519
// seal signature, and (with --master-key) per-event HMAC re-derivation.
//
// The package has no I/O of its own beyond a one-shot ledger reader; it
// operates on parsed structures so the same primitives are testable against
// fixtures and reusable from contexts other than the verifier CLI.
//
// Bootstrap ledger format (NDJSON, one record per line):
//
//   - ChainEntry — every captured event, with chain linkage and Merkle leaf.
//   - SealRecord — exactly one per tenant-day, carrying the Merkle root and
//     the HSM-signed Ed25519 signature. Distinguished by the "type": "seal"
//     marker.
//
// The format is sufficient for the bootstrap; production ledgers may use a
// different on-wire encoding (binary, length-prefixed, etc.). Replacing the
// reader is the only change downstream consumers should need.
package verify

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ChainEntry is one event in the chain.
type ChainEntry struct {
	EntryID               string `json:"entry_id"`
	PrevHash              string `json:"prev_hash"`               // base64 SHA-256, empty on genesis
	EntryHash             string `json:"entry_hash"`              // base64 SHA-256
	HMACSHA256            string `json:"hmac_sha256"`             // base64, 32 bytes
	TenantBindingKDFLabel string `json:"tenant_binding_kdf_label"`
	EventPayloadJCS       string `json:"event_payload_jcs"` // base64 of the canonical event bytes
	MerkleLeafIndex       int    `json:"merkle_leaf_index"`
	SealDate              string `json:"seal_date"` // YYYY-MM-DD
}

// SealRecord is the daily Merkle seal, signed by the institution's
// HSM-resident Ed25519 signing key.
type SealRecord struct {
	Type                  string `json:"type"` // always "seal"
	SealDate              string `json:"seal_date"`
	MerkleRoot            string `json:"merkle_root"` // base64 SHA-256
	LeafCount             int    `json:"leaf_count"`
	SignatureEd25519      string `json:"signature_ed25519"`       // base64, 64 bytes
	SigningKeyFingerprint string `json:"signing_key_fingerprint"` // "sha256:<hex>"
}

// Ledger is the parsed contents of one tenant-day file.
type Ledger struct {
	Entries []ChainEntry
	Seal    SealRecord
}

// LoadLedger reads an NDJSON ledger file from path. Each line is one record.
// Returns the parsed entries and the (one) seal record. If the file has zero
// entries or no seal record, LoadLedger returns an error so the verifier
// fails closed.
func LoadLedger(path string) (*Ledger, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open ledger: %w", err)
	}
	defer f.Close()
	return readLedger(f)
}

// readLedger is the io.Reader-driven core of LoadLedger; split out so tests
// can drive it from in-memory bytes without touching the disk.
func readLedger(r io.Reader) (*Ledger, error) {
	led := &Ledger{}
	sealSeen := false
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}
		// A small two-step parse: peek at "type" to fork between ChainEntry
		// and SealRecord without two full unmarshals on every line.
		var probe struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal([]byte(raw), &probe); err != nil {
			return nil, fmt.Errorf("line %d: invalid JSON: %w", lineNo, err)
		}
		if probe.Type == "seal" {
			if sealSeen {
				return nil, fmt.Errorf("line %d: more than one seal record in ledger", lineNo)
			}
			var s SealRecord
			if err := json.Unmarshal([]byte(raw), &s); err != nil {
				return nil, fmt.Errorf("line %d: parse seal: %w", lineNo, err)
			}
			led.Seal = s
			sealSeen = true
			continue
		}
		var e ChainEntry
		if err := json.Unmarshal([]byte(raw), &e); err != nil {
			return nil, fmt.Errorf("line %d: parse entry: %w", lineNo, err)
		}
		led.Entries = append(led.Entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan ledger: %w", err)
	}

	if len(led.Entries) == 0 {
		return nil, errors.New("ledger has no chain entries")
	}
	if !sealSeen {
		return nil, errors.New("ledger has no seal record")
	}
	return led, nil
}
