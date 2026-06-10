package vectors

import "testing"

// TestReasonMatches_QuotingIsByteExact is the regression guard the N023
// tolerance-tightening leaves behind. The earlier reasonMatches stripped
// ASCII double-quotes from both sides so a fixture's quoted format_version
// reason compared equal to the verifier's then-unquoted rendering. That
// tolerance is gone: Heather re-rendered N009/N022/N023 uniformly to the
// quoted form and the verifier renders the quoted form too, so the compare
// is byte-exact (modulo the normative `: detail` suffix).
//
// These cases lock that in. A future verifier change that reverted to the
// unquoted rendering, or a future fixture re-rendered with the wrong
// quoting, would make the quoted/unquoted pair compare UNEQUAL here and
// fail the gate loudly — which is the entire point of removing the
// tolerance.
func TestReasonMatches_QuotingIsByteExact(t *testing.T) {
	const quoted = `format_version "V1" not supported by this verifier (running v1)`
	const unquoted = `format_version V1 not supported by this verifier (running v1)`

	tests := []struct {
		name string
		got  string
		pin  string
		want bool
	}{
		{
			name: "quoted verifier reason matches quoted pin (exact)",
			got:  quoted,
			pin:  quoted,
			want: true,
		},
		{
			name: "unquoted verifier reason does NOT match quoted pin (tolerance removed)",
			got:  unquoted,
			pin:  quoted,
			want: false,
		},
		{
			name: "quoted verifier reason does NOT match unquoted pin (tolerance removed)",
			got:  quoted,
			pin:  unquoted,
			want: false,
		},
		{
			name: "normative ': detail' suffix still tolerated",
			got:  `key_fingerprint mismatch at seq 4: looked-up IKM does not match`,
			pin:  `key_fingerprint mismatch at seq 4`,
			want: true,
		},
		{
			name: "a different ':' that is not the detail boundary is not tolerated",
			got:  `chain link broken at seq 4`,
			pin:  `chain link broken at seq 5`,
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := reasonMatches(tc.got, expectedOutput{Reason: tc.pin})
			if got != tc.want {
				t.Errorf("reasonMatches(%q, pin=%q) = %v, want %v", tc.got, tc.pin, got, tc.want)
			}
		})
	}
}
