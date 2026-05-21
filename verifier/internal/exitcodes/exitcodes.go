// Package exitcodes is the §10.12 + §10.29 closed-enumeration source
// of truth for verifier process exit codes.
//
// Spec authority: §10.12 (terminal codes 0-3) + §10.29 (streaming-
// mode non-terminal codes 4-6). Implementations MAY define codes
// >= 7 for vendor-specific diagnostics; examiner harnesses MUST
// treat any code >= 7 as opaque and not branch on it. The codes 0-3
// are the conformance contract.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/VerifierVerdict.cs uses the same
// integer values; the Python reference at
// Herald.Py/src/herald/_verdict.py is the third witness. All three
// MUST agree byte-for-byte on the integer values, the
// terminal/streaming partition, and the closed-enum shape.
package exitcodes

import "fmt"

// Code is the verifier process exit code per §10.12 + §10.29. The
// underlying type is `int` to match os.Exit's signature and the
// integer values examiner harnesses branch on; the named type lets
// callers express intent without leaking magic numbers.
type Code int

// Terminal codes per §10.12. Returned by os.Exit at end-of-run.
const (
	// Pass — the verifier completed and the chain verified (or the
	// structural subset verified under witness mode). Bonus
	// verifications travel in the verdict object's
	// additional_verifications array, never via a new exit code.
	Pass Code = 0

	// Fail — the chain failed integrity verification at one of the
	// §7 steps. The reason and step number appear on stdout per the
	// §7 normative output format.
	Fail Code = 1

	// StructuralInputError — the verifier could not BEGIN the §7
	// procedure: file unreadable, JSON malformed, mandatory header
	// field missing entirely (distinct from an invalid value), file
	// empty. When the verifier reached §7 — even if step 1 rejects
	// on an unsupported format_version value — the code is Fail.
	StructuralInputError Code = 2

	// ConfigurationError — the verifier was invoked without a
	// required argument (e.g., --master-key under --strict), the
	// algorithm is unknown, or the posture flag does not match the
	// chain's posture. Also emitted by the §7 pre-flight JCS
	// self-test on a JCS implementation mismatch.
	ConfigurationError Code = 3
)

// Streaming-mode non-terminal codes per §10.29. Only emitted by
// streaming-mode verifiers (per_second / per_minute / per_hour
// cadences per §10.27); a daily-cadence verifier never returns
// these. The codes are non-terminal — a streaming-mode verifier
// may transition between them as the chain stream progresses,
// terminating in one of Pass / Fail / StructuralInputError /
// ConfigurationError when the streaming run ends.
const (
	// StreamingAllPassSoFar — every chain entry consumed so far has
	// verified; no anomaly or rotation pending.
	StreamingAllPassSoFar Code = 4

	// StreamingAnomalyDetected — an anomaly surfaced during the
	// streaming walk; the verifier continues consuming and may
	// transition back to StreamingAllPassSoFar after a clean
	// recovery, or terminate at Fail if the anomaly is
	// integrity-bearing.
	StreamingAnomalyDetected Code = 5

	// StreamingRotationPending — the streaming walk observed a
	// key-rotation event whose post-rotation seal has not yet
	// arrived. The verifier transitions back to
	// StreamingAllPassSoFar once the post-rotation seal arrives and
	// the cross-rotation chain links validate.
	StreamingRotationPending Code = 6
)

// IsTerminal reports whether c is one of the terminal exit codes
// (Pass, Fail, StructuralInputError, ConfigurationError) per
// §10.12. A streaming-mode verifier ends its run by returning a
// terminal code; non-terminal codes are emitted mid-stream only.
func (c Code) IsTerminal() bool {
	switch c {
	case Pass, Fail, StructuralInputError, ConfigurationError:
		return true
	}
	return false
}

// IsStreaming reports whether c is one of the §10.29 streaming-mode
// non-terminal codes (StreamingAllPassSoFar, StreamingAnomalyDetected,
// StreamingRotationPending).
func (c Code) IsStreaming() bool {
	switch c {
	case StreamingAllPassSoFar, StreamingAnomalyDetected, StreamingRotationPending:
		return true
	}
	return false
}

// String returns the spec-named identifier for c. Unrecognized
// values (vendor-specific codes >= 7) return "vendor-specific(N)"
// so a log line never silently hides a non-spec code.
func (c Code) String() string {
	switch c {
	case Pass:
		return "PASS"
	case Fail:
		return "FAIL"
	case StructuralInputError:
		return "STRUCTURAL_INPUT_ERROR"
	case ConfigurationError:
		return "CONFIGURATION_ERROR"
	case StreamingAllPassSoFar:
		return "STREAMING_ALL_PASS_SO_FAR"
	case StreamingAnomalyDetected:
		return "STREAMING_ANOMALY_DETECTED"
	case StreamingRotationPending:
		return "STREAMING_ROTATION_PENDING"
	}
	return fmt.Sprintf("vendor-specific(%d)", int(c))
}
