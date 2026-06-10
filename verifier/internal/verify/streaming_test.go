package verify

import "testing"

// TestStreamState_TransitionTable exercises the §10.29 transition table
// directly (independent of the corpus) so the state machine's contract
// is pinned even when the 022 vector is absent. Each row is one
// (from, kind, result) → to transition from the spec table.
func TestStreamState_TransitionTable(t *testing.T) {
	tests := []struct {
		name   string
		from   int
		kind   string
		result string
		want   int
	}{
		{"active+entry-pass→active", StreamStateActive, StreamKindChainEntry, StreamResultPass, StreamStateActive},
		{"active+entry-fail→anomaly", StreamStateActive, StreamKindChainEntry, StreamResultFail, StreamStateAnomaly},
		{"active+seal-pass→active", StreamStateActive, StreamKindSeal, StreamResultPass, StreamStateActive},
		{"active+seal-fail→anomaly", StreamStateActive, StreamKindSeal, StreamResultFail, StreamStateAnomaly},
		{"active+rotation→rotation-pending", StreamStateActive, StreamKindRotation, "n/a", StreamStateRotationPending},
		{"rotation-pending+seal-pass→active", StreamStateRotationPending, StreamKindSeal, StreamResultPass, StreamStateActive},
		{"rotation-pending+seal-fail→anomaly", StreamStateRotationPending, StreamKindSeal, StreamResultFail, StreamStateAnomaly},
		{"rotation-pending+entry-pass→rotation-pending", StreamStateRotationPending, StreamKindChainEntry, StreamResultPass, StreamStateRotationPending},
		{"rotation-pending+entry-fail→anomaly", StreamStateRotationPending, StreamKindChainEntry, StreamResultFail, StreamStateAnomaly},
		{"rotation-pending+rotation→anomaly", StreamStateRotationPending, StreamKindRotation, "n/a", StreamStateAnomaly},
		{"anomaly+entry-pass→anomaly(sticky)", StreamStateAnomaly, StreamKindChainEntry, StreamResultPass, StreamStateAnomaly},
		{"anomaly+seal-pass→anomaly(sticky)", StreamStateAnomaly, StreamKindSeal, StreamResultPass, StreamStateAnomaly},
		{"anomaly+rotation→anomaly(sticky)", StreamStateAnomaly, StreamKindRotation, "n/a", StreamStateAnomaly},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &StreamState{verdict: tc.from}
			s.Step(tc.kind, tc.result)
			if s.Verdict() != tc.want {
				t.Errorf("from=%d %s/%s: got=%d want=%d", tc.from, tc.kind, tc.result, s.Verdict(), tc.want)
			}
		})
	}
}

// TestStreamState_Finalize pins the §10.29 finalize-collapse table.
func TestStreamState_Finalize(t *testing.T) {
	tests := []struct {
		state int
		want  int
	}{
		{StreamStateActive, 0},
		{StreamStateAnomaly, 3},
		{StreamStateRotationPending, 3},
	}
	for _, tc := range tests {
		s := &StreamState{verdict: tc.state}
		if got := s.Finalize(); got != tc.want {
			t.Errorf("finalize(%d): got=%d want=%d", tc.state, got, tc.want)
		}
	}
}

// TestStreamState_StartsActive confirms NewStreamState begins at the
// active (all-pass-so-far) verdict per §10.29.
func TestStreamState_StartsActive(t *testing.T) {
	if v := NewStreamState().Verdict(); v != StreamStateActive {
		t.Errorf("new stream state verdict = %d, want %d", v, StreamStateActive)
	}
}
