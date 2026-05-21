package exitcodes

import "testing"

// TestCodeIntegerValues pins the integer values per §10.12 + §10.29.
// Two implementations disagreeing on these integers would silently
// mis-route examiner-harness branches; the byte-pin is the contract.
func TestCodeIntegerValues(t *testing.T) {
	cases := []struct {
		code Code
		want int
	}{
		{Pass, 0},
		{Fail, 1},
		{StructuralInputError, 2},
		{ConfigurationError, 3},
		{StreamingAllPassSoFar, 4},
		{StreamingAnomalyDetected, 5},
		{StreamingRotationPending, 6},
	}
	for _, tc := range cases {
		if got := int(tc.code); got != tc.want {
			t.Errorf("%s = %d, want %d (spec §10.12 / §10.29)", tc.code, got, tc.want)
		}
	}
}

func TestIsTerminal(t *testing.T) {
	terminal := []Code{Pass, Fail, StructuralInputError, ConfigurationError}
	for _, c := range terminal {
		if !c.IsTerminal() {
			t.Errorf("%s.IsTerminal() = false, want true", c)
		}
	}
	streaming := []Code{StreamingAllPassSoFar, StreamingAnomalyDetected, StreamingRotationPending}
	for _, c := range streaming {
		if c.IsTerminal() {
			t.Errorf("%s.IsTerminal() = true, want false (streaming codes are non-terminal)", c)
		}
	}
}

func TestIsStreaming(t *testing.T) {
	streaming := []Code{StreamingAllPassSoFar, StreamingAnomalyDetected, StreamingRotationPending}
	for _, c := range streaming {
		if !c.IsStreaming() {
			t.Errorf("%s.IsStreaming() = false, want true", c)
		}
	}
	terminal := []Code{Pass, Fail, StructuralInputError, ConfigurationError}
	for _, c := range terminal {
		if c.IsStreaming() {
			t.Errorf("%s.IsStreaming() = true, want false", c)
		}
	}
}

func TestStringNamesAreStable(t *testing.T) {
	cases := []struct {
		code Code
		want string
	}{
		{Pass, "PASS"},
		{Fail, "FAIL"},
		{StructuralInputError, "STRUCTURAL_INPUT_ERROR"},
		{ConfigurationError, "CONFIGURATION_ERROR"},
		{StreamingAllPassSoFar, "STREAMING_ALL_PASS_SO_FAR"},
		{StreamingAnomalyDetected, "STREAMING_ANOMALY_DETECTED"},
		{StreamingRotationPending, "STREAMING_ROTATION_PENDING"},
	}
	for _, tc := range cases {
		if got := tc.code.String(); got != tc.want {
			t.Errorf("Code(%d).String() = %q, want %q", int(tc.code), got, tc.want)
		}
	}
}

func TestStringVendorSpecific(t *testing.T) {
	if got := Code(42).String(); got != "vendor-specific(42)" {
		t.Errorf("Code(42).String() = %q, want %q", got, "vendor-specific(42)")
	}
}

func TestVendorCodeIsNotTerminalNorStreaming(t *testing.T) {
	c := Code(42)
	if c.IsTerminal() {
		t.Error("vendor code (>= 7) reported terminal; only 0-3 are terminal")
	}
	if c.IsStreaming() {
		t.Error("vendor code (>= 7) reported streaming; only 4-6 are streaming")
	}
}
