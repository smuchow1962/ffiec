package chain

import (
	"testing"
	"time"
)

// invalidCadences is the shared set of strings that MUST be rejected as
// non-conformant per §10.27. The set covers empty, case-variants, whitespace
// variants, and out-of-enumeration values. IsValid, IsStreaming, and
// Interval all reject these consistently — that's an invariant the table-
// driven implementation in cadence.go upholds.
var invalidCadences = []Cadence{
	Cadence(""),
	Cadence("monthly"),
	Cadence("yearly"),
	Cadence("Daily"),
	Cadence("DAILY"),
	Cadence(" daily"),
	Cadence("daily "),
	Cadence("PER_SECOND"),
	Cadence("daily\n"), // trailing-newline (line-terminator-leakage rejection)
	Cadence("per-second"), // hyphen-instead-of-underscore mistake
	Cadence("perSecond"), // camelCase mistake
}

// validCadences is the complete §10.27 enumeration in the order the spec
// lists them.
var validCadences = []Cadence{
	CadencePerSecond,
	CadencePerMinute,
	CadencePerHour,
	CadenceHourly,
	CadenceDaily,
	CadenceWeekly,
}

func TestCadenceIsStreaming(t *testing.T) {
	streamingValid := map[Cadence]bool{
		CadencePerSecond: true,
		CadencePerMinute: true,
		CadencePerHour:   true,
		CadenceHourly:    false,
		CadenceDaily:     false,
		CadenceWeekly:    false,
	}
	for c, want := range streamingValid {
		t.Run(string(c), func(t *testing.T) {
			if got := c.IsStreaming(); got != want {
				t.Errorf("Cadence(%q).IsStreaming() = %v, want %v", c, got, want)
			}
		})
	}
	// Per the doc-comment contract, invalid cadences return false.
	for _, c := range invalidCadences {
		t.Run("invalid:"+string(c), func(t *testing.T) {
			if c.IsStreaming() {
				t.Errorf("Cadence(%q).IsStreaming() = true, want false (invalid input)", c)
			}
		})
	}
}

func TestCadenceIsValid(t *testing.T) {
	for _, c := range validCadences {
		t.Run(string(c), func(t *testing.T) {
			if !c.IsValid() {
				t.Errorf("Cadence(%q).IsValid() = false, want true", c)
			}
		})
	}
	for _, c := range invalidCadences {
		t.Run("invalid:"+string(c), func(t *testing.T) {
			if c.IsValid() {
				t.Errorf("Cadence(%q).IsValid() = true, want false", c)
			}
		})
	}
}

func TestCadenceInterval(t *testing.T) {
	expectedIntervals := map[Cadence]time.Duration{
		CadencePerSecond: time.Second,
		CadencePerMinute: time.Minute,
		CadencePerHour:   time.Hour,
		CadenceHourly:    time.Hour,
		CadenceDaily:     24 * time.Hour,
		CadenceWeekly:    7 * 24 * time.Hour,
	}
	for c, want := range expectedIntervals {
		t.Run(string(c), func(t *testing.T) {
			if got := c.Interval(); got != want {
				t.Errorf("Cadence(%q).Interval() = %v, want %v", c, got, want)
			}
		})
	}
	// Per the doc-comment contract, invalid cadences return 0.
	for _, c := range invalidCadences {
		t.Run("invalid:"+string(c), func(t *testing.T) {
			if got := c.Interval(); got != 0 {
				t.Errorf("Cadence(%q).Interval() = %v, want 0 (invalid input)", c, got)
			}
		})
	}
}

// TestCadencePerHourAndHourlyAreEquivalentInterval pins the §10.27 normative
// invariant: CadencePerHour and CadenceHourly are byte-distinct values
// that produce the same interval. The streaming-mode flag differs; the
// interval does not.
func TestCadencePerHourAndHourlyAreEquivalentInterval(t *testing.T) {
	if CadencePerHour.Interval() != CadenceHourly.Interval() {
		t.Errorf("per_hour and hourly intervals differ: %v vs %v",
			CadencePerHour.Interval(), CadenceHourly.Interval())
	}
	if !CadencePerHour.IsStreaming() {
		t.Errorf("per_hour MUST be streaming-mode per §10.27")
	}
	if CadenceHourly.IsStreaming() {
		t.Errorf("hourly MUST NOT be streaming-mode per §10.27")
	}
}

// TestCadenceStreamingPartitionInvariant pins the partition-soundness
// invariant: every IsStreaming()=true value has Interval() < 24h, and
// every IsStreaming()=false valid value has Interval() >= 1h. This
// matches the spec §10.27 streaming-vs-non-streaming partition.
func TestCadenceStreamingPartitionInvariant(t *testing.T) {
	for _, c := range validCadences {
		if c.IsStreaming() {
			if c.Interval() >= 24*time.Hour {
				t.Errorf("streaming cadence %q has interval >= 24h: %v", c, c.Interval())
			}
		} else {
			if c.Interval() < time.Hour {
				t.Errorf("non-streaming cadence %q has interval < 1h: %v", c, c.Interval())
			}
		}
	}
}

// TestCadenceConstantsAreDistinct pins the wire-format invariant that
// the six exported normative constants have pairwise distinct string
// values. A copy-paste regression that aliased two constants to the
// same byte form would produce ambiguous wire-format output and a chain
// that survives the type system but fails examiner-side reproduction.
func TestCadenceConstantsAreDistinct(t *testing.T) {
	seen := make(map[Cadence]bool, len(validCadences))
	for _, c := range validCadences {
		if seen[c] {
			t.Errorf("duplicate cadence constant value: %q", c)
		}
		seen[c] = true
	}
	if len(seen) != len(validCadences) {
		t.Errorf("expected %d distinct cadences, got %d", len(validCadences), len(seen))
	}
}
