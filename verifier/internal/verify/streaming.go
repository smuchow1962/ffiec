package verify

// §10.29 streaming-verifier state machine.
//
// A streaming-mode verifier (per_second / per_minute / per_hour cadence
// per §10.27) holds one of three non-terminal verdicts as it consumes a
// chain stream, then collapses to a terminal verdict at end-of-stream:
//
//	4 — all-pass-so-far     → terminal 0 (PASS)
//	5 — anomaly-detected    → terminal 3 (chain anomaly; sticky)
//	6 — rotation-pending    → terminal 3 (unconfirmed rotation)
//
// The transition table and finalize collapse below are the §10.29
// normative dispatch (pinned byte-for-byte by test vector 022). The
// machine is verifier capability — a streaming CLI mode drives it — not
// a test-only construct, so it lives in the verify package with an
// exported API.

// Streaming verdict states per §10.29 (= exitcodes 4/5/6).
const (
	StreamStateActive          = 4 // all-pass-so-far
	StreamStateAnomaly         = 5 // anomaly detected (sticky)
	StreamStateRotationPending = 6 // rotation observed, confirming seal pending
)

// Stream input kinds per §10.29.
const (
	StreamKindChainEntry = "chain_entry"
	StreamKindSeal       = "seal"
	StreamKindRotation   = "rotation"
)

// Stream input results per §10.29.
const (
	StreamResultPass = "pass"
	StreamResultFail = "fail"
)

// StreamState is a §10.29 streaming-verifier verdict, advanced one input
// at a time via Step and collapsed to a terminal exit code via Finalize.
// The zero value is NOT valid; construct with NewStreamState so the
// initial verdict is the active (all-pass-so-far) state per §10.29.
type StreamState struct {
	verdict int
}

// NewStreamState returns a streaming verdict initialized to the active
// (all-pass-so-far) state, the §10.29 starting verdict before any input.
func NewStreamState() *StreamState {
	return &StreamState{verdict: StreamStateActive}
}

// Verdict returns the current streaming verdict (4, 5, or 6).
func (s *StreamState) Verdict() int { return s.verdict }

// Step advances the verdict given one stream input (kind, result),
// implementing the §10.29 transition table:
//
//   - Anomaly (5) is STICKY: once entered, every input keeps it at 5.
//   - A rotation moves active→rotation_pending (4→6); a SECOND rotation
//     before a confirming seal is an integrity-claim violation (6→5).
//   - A confirming seal under rotation_pending drops back to active
//     (6→4 on pass, 6→5 on fail).
//   - Any failing chain_entry or seal moves to anomaly (→5).
//
// Unknown kinds/results leave the verdict unchanged — a streaming
// verifier must not crash or guess on an input shape it does not
// recognize (§7 forward-compatibility fault-tolerance).
func (s *StreamState) Step(kind, result string) {
	s.verdict = nextStreamVerdict(s.verdict, kind, result)
}

func nextStreamVerdict(state int, kind, result string) int {
	if state == StreamStateAnomaly {
		return StreamStateAnomaly // sticky anomaly
	}

	switch kind {
	case StreamKindRotation:
		// active → rotation_pending; rotation_pending → anomaly
		// (back-to-back rotation without a confirming seal).
		if state == StreamStateRotationPending {
			return StreamStateAnomaly
		}
		return StreamStateRotationPending

	case StreamKindChainEntry:
		if result == StreamResultFail {
			return StreamStateAnomaly
		}
		return state // pass keeps the current state (4 or 6)

	case StreamKindSeal:
		if result == StreamResultFail {
			return StreamStateAnomaly
		}
		// A passing seal confirms a pending rotation, dropping to active.
		if state == StreamStateRotationPending {
			return StreamStateActive
		}
		return state
	}

	return state
}

// Finalize collapses the current streaming verdict to its terminal exit
// code at end-of-stream per the §10.29 finalize table: active→0,
// anomaly→3, rotation_pending→3 (an unconfirmed rotation at EOF is a
// control-completeness failure).
func (s *StreamState) Finalize() int {
	switch s.verdict {
	case StreamStateActive:
		return 0
	case StreamStateAnomaly, StreamStateRotationPending:
		return 3
	}
	return 3 // unknown verdict collapses conservatively to a failure
}
