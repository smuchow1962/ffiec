package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mmpworks/ffiec/cliutil"
	"github.com/mmpworks/ffiec/verifier/internal/verify"
)

func runDiff(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("diff", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var (
		beforePath = fs.String("before", "", "Path to the earlier ledger (required)")
		afterPath  = fs.String("after", "", "Path to the later ledger (required)")
	)
	if stop, err := cliutil.ParseFlags(fs, args, stdout, progName); err != nil {
		return err
	} else if stop {
		return nil
	}
	if err := cliutil.RequireFlags(map[string]string{
		"before": *beforePath,
		"after":  *afterPath,
	}); err != nil {
		return err
	}

	a, err := verify.LoadLedger(*beforePath)
	if err != nil {
		return fmt.Errorf("--before: %w", err)
	}
	b, err := verify.LoadLedger(*afterPath)
	if err != nil {
		return fmt.Errorf("--after: %w", err)
	}

	sameBytes, err := sameFileBytes(*beforePath, *afterPath)
	if err != nil {
		return err
	}
	mismatches := compareLedgers(a, b)
	if !sameBytes && len(mismatches) == 0 {
		mismatches = []string{"file bytes differ, outside the seal and entry fields this command names"}
	}
	if len(mismatches) == 0 {
		printLine(stdout, "ledgers identical (every byte)")
		return nil
	}
	printLine(stdout, "%d differences found:", len(mismatches))
	for _, m := range mismatches {
		printLine(stdout, "  %s", m)
	}
	return fmt.Errorf("ledgers differ")
}

// sameFileBytes is what "identical" means here. compareLedgers names the fields an examiner reads
// first, and a ledger has more fields than that, some of them outside the seal's signature.
func sameFileBytes(before, after string) (bool, error) {
	a, err := os.ReadFile(before)
	if err != nil {
		return false, fmt.Errorf("--before: %w", err)
	}
	b, err := os.ReadFile(after)
	if err != nil {
		return false, fmt.Errorf("--after: %w", err)
	}
	return bytes.Equal(a, b), nil
}

// compareLedgers produces a deterministic list of human-readable
// mismatch descriptions. The function is split out for unit testability.
func compareLedgers(a, b *verify.Ledger) []string {
	var out []string
	if a.Seal.MerkleRoot != b.Seal.MerkleRoot {
		out = append(out, fmt.Sprintf("seal merkle_root differs: before=%s after=%s",
			short(a.Seal.MerkleRoot), short(b.Seal.MerkleRoot)))
	}
	if a.Seal.SignatureEd25519 != b.Seal.SignatureEd25519 {
		out = append(out, "seal signature_ed25519 differs")
	}
	if a.Seal.LeafCount != b.Seal.LeafCount {
		out = append(out, fmt.Sprintf("seal leaf_count differs: before=%d after=%d",
			a.Seal.LeafCount, b.Seal.LeafCount))
	}

	n := len(a.Entries)
	if n2 := len(b.Entries); n2 != n {
		out = append(out, fmt.Sprintf("entry count differs: before=%d after=%d", n, n2))
		if n2 < n {
			n = n2
		}
	}
	for i := 0; i < n; i++ {
		ae, be := &a.Entries[i], &b.Entries[i]
		if ae.EntryID != be.EntryID {
			out = append(out, fmt.Sprintf("entry %d entry_id differs: before=%s after=%s",
				i, ae.EntryID, be.EntryID))
		}
		if ae.EntryHash != be.EntryHash {
			out = append(out, fmt.Sprintf("entry %d entry_hash differs (id=%s)", i, ae.EntryID))
		}
	}
	return out
}
