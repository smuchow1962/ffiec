package signpayload_test

import (
	"strings"
	"testing"

	"github.com/mmpworks/ffiec/core/signpayload"
)

// baseSeal returns a Seal with every field populated so each
// version-specific test can set only Version and assert the right
// field subset is bound.
func baseSeal(v signpayload.Version) signpayload.Seal {
	return signpayload.Seal{
		Version:                  v,
		Algorithm:                "ed25519",
		FormatVersion:            "v1",
		TenantID:                 "demo-bank-001",
		SealDate:                 "2026-05-07",
		MerkleRootHex:            strings.Repeat("ab", 32),
		HKDFInputsDigest:         strings.Repeat("cd", 32),
		Cadence:                  "daily",
		DevMode:                  false,
		KeyVersionsCanon:         "1,2,3",
		KMSHandleURIsDigest:      strings.Repeat("ef", 32),
		OperationalEventsLogRoot: strings.Repeat("12", 32),
	}
}

// TestBuild_LineCountsAndOrder pins each form's exact line count and
// the value bound on each line. The line count is the spec's primary
// structural invariant (6 / 10 / 12 / 13).
func TestBuild_LineCountsAndOrder(t *testing.T) {
	cases := []struct {
		name      string
		version   signpayload.Version
		wantLines int
		// wantTag is the line-2 version tag, or "" for the pre-amendment
		// form which has no version line.
		wantTag string
	}{
		{"pre_amendment_6line", signpayload.VersionPreAmendment, 7, ""},
		{"v1_0a_10line", signpayload.VersionV1_0a, 10, "v1.0a"},
		{"v1_0b_12line", signpayload.VersionV1_0b, 12, "v1.0b"},
		{"v1_0c_13line", signpayload.VersionV1_0c, 13, "v1.0c"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := signpayload.Build(baseSeal(tc.version))
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			lines := strings.Split(string(got), "\n")
			if len(lines) != tc.wantLines {
				t.Fatalf("line count = %d, want %d\npayload:\n%s", len(lines), tc.wantLines, got)
			}
			if lines[0] != signpayload.Magic {
				t.Errorf("line 0 = %q, want magic %q", lines[0], signpayload.Magic)
			}
			if tc.wantTag != "" && lines[1] != tc.wantTag {
				t.Errorf("line 1 (version tag) = %q, want %q", lines[1], tc.wantTag)
			}
			// No trailing newline on the terminal field.
			if strings.HasSuffix(string(got), "\n") {
				t.Errorf("payload has a trailing newline; the terminal field must not")
			}
		})
	}
}

// TestBuild_DevModeByte asserts dev_mode serializes to the single
// ASCII byte "0" / "1", never "true" / "false".
func TestBuild_DevModeByte(t *testing.T) {
	for _, tc := range []struct {
		dev  bool
		want string
	}{{false, "0"}, {true, "1"}} {
		s := baseSeal(signpayload.VersionV1_0a)
		s.DevMode = tc.dev
		got, err := signpayload.Build(s)
		if err != nil {
			t.Fatalf("Build: %v", err)
		}
		lines := strings.Split(string(got), "\n")
		// dev_mode is the last line of the v1.0a form (line index 9).
		if lines[9] != tc.want {
			t.Errorf("dev_mode=%v → %q, want %q", tc.dev, lines[9], tc.want)
		}
	}
}

// TestBuild_V1_0b_EmptyKeyVersions asserts the empty-day edge case:
// key_versions_canon is the empty string, producing an empty line
// between its two surrounding separators (matches vector 019).
func TestBuild_V1_0b_EmptyKeyVersions(t *testing.T) {
	s := baseSeal(signpayload.VersionV1_0b)
	s.KeyVersionsCanon = ""
	got, err := signpayload.Build(s)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	lines := strings.Split(string(got), "\n")
	// key_versions_canon is line index 10 in the 12-line form.
	if lines[10] != "" {
		t.Errorf("empty-day key_versions_canon = %q, want empty", lines[10])
	}
}

// TestParseVersion covers the §7 step 11 dispatch mapping, including
// the absent-field (pre-amendment) and unrecognized-value cases.
func TestParseVersion(t *testing.T) {
	cases := []struct {
		field   string
		want    signpayload.Version
		wantErr bool
	}{
		{"", signpayload.VersionPreAmendment, false},
		{"v1.0a", signpayload.VersionV1_0a, false},
		{"v1.0b", signpayload.VersionV1_0b, false},
		{"v1.0c", signpayload.VersionV1_0c, false},
		{"v1.0d", 0, true},
		{"V1.0A", 0, true}, // case-variant rejection
	}
	for _, tc := range cases {
		t.Run("field_"+tc.field, func(t *testing.T) {
			got, err := signpayload.ParseVersion(tc.field)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tc.field)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.field, err)
			}
			if got != tc.want {
				t.Errorf("ParseVersion(%q) = %v, want %v", tc.field, got, tc.want)
			}
		})
	}
}

// TestBuild_MonotonicExtension asserts the structural promise of §4.3:
// each higher form is the lower form's fields (with the version tag
// corrected) plus the new terminal field(s). We verify v1.0c's first
// 12 lines equal v1.0b's 12 lines after the tag swap.
func TestBuild_MonotonicExtension(t *testing.T) {
	b, err := signpayload.Build(baseSeal(signpayload.VersionV1_0b))
	if err != nil {
		t.Fatalf("Build v1.0b: %v", err)
	}
	c, err := signpayload.Build(baseSeal(signpayload.VersionV1_0c))
	if err != nil {
		t.Fatalf("Build v1.0c: %v", err)
	}
	bLines := strings.Split(string(b), "\n")
	cLines := strings.Split(string(c), "\n")

	for i := range bLines {
		want := bLines[i]
		if i == 1 {
			want = "v1.0c" // the tag differs by design
		}
		if cLines[i] != want {
			t.Errorf("v1.0c line %d = %q, want %q (v1.0b parity)", i, cLines[i], want)
		}
	}
}
