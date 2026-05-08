package main

import (
	"flag"
	"fmt"
	"io"

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

	mismatches := compareLedgers(a, b)
	if len(mismatches) == 0 {
		fmt.Fprintln(stdout, "ledgers identical (chain entries and seal)")
		return nil
	}
	fmt.Fprintf(stdout, "%d differences found:\n", len(mismatches))
	for _, m := range mismatches {
		fmt.Fprintf(stdout, "  %s\n", m)
	}
	return fmt.Errorf("ledgers differ")
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
