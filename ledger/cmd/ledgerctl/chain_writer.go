package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/mmpworks/ffiec/ledger/internal/examiner"
)

// WriterChainWriter is the bootstrap ChainWriter. It writes each event as
// pretty-printed JSON to W. ledgerctl wires W = stdout in production and
// W = a buffer in tests. Replace with the real append-only ledger writer
// once that path lands in ledger/internal.
type WriterChainWriter struct {
	W io.Writer
}

// Write satisfies examiner.ChainWriter.
func (w WriterChainWriter) Write(_ context.Context, evt examiner.OperationalEvent) error {
	enc := json.NewEncoder(w.W)
	enc.SetIndent("", "  ")
	if err := enc.Encode(evt); err != nil {
		return fmt.Errorf("writer chain writer: %w", err)
	}
	return nil
}
