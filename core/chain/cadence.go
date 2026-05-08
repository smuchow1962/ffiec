// Package chain implements the per-event MAC chain primitive and supporting
// types for the FFIEC chain-of-custody specification.
//
// This file defines Cadence per spec §10.27 (Configurable seal cadence,
// which extends the §4.2.1 enumeration with sub-daily streaming-mode values).
// The default cadence is daily; per_second / per_minute / per_hour are the
// streaming-mode values per §10.27.
package chain

import "time"

// Cadence is the seal aggregation cadence per §10.27 (extending §4.2.1).
//
// The string value is the wire-format byte value recorded in the seal
// record's "cadence" field and bound under the §4.3 sign_payload form.
// Verifiers reject any unrecognized value at §7 step 12 with the reason
// `cadence "X" is not in the §10.27 enumeration`.
type Cadence string

// Normative cadence values per §10.27.
//
// Sub-daily values (PerSecond, PerMinute, PerHour) are streaming-mode per
// §10.27. CadencePerHour and CadenceHourly are byte-distinct values that
// both produce a one-hour interval; per §10.27 they are NOT interchangeable
// within a single chain.
const (
	CadencePerSecond Cadence = "per_second"
	CadencePerMinute Cadence = "per_minute"
	CadencePerHour   Cadence = "per_hour"
	CadenceHourly    Cadence = "hourly"
	CadenceDaily     Cadence = "daily"
	CadenceWeekly    Cadence = "weekly"
)

// cadenceProperties is the single source of truth for the validity, the
// streaming-mode classification, and the interval of every normative
// cadence. Adding a new cadence value is a single-line table edit; the
// methods below dispatch through this map.
//
// This table-driven shape keeps the three exported predicates (IsValid,
// IsStreaming, Interval) consistent: a future amendment cannot accidentally
// add a value to one method's switch and forget the others.
var cadenceProperties = map[Cadence]struct {
	streaming bool
	interval  time.Duration
}{
	CadencePerSecond: {streaming: true, interval: time.Second},
	CadencePerMinute: {streaming: true, interval: time.Minute},
	CadencePerHour:   {streaming: true, interval: time.Hour},
	CadenceHourly:    {streaming: false, interval: time.Hour},
	CadenceDaily:     {streaming: false, interval: 24 * time.Hour},
	CadenceWeekly:    {streaming: false, interval: 7 * 24 * time.Hour},
}

// IsValid reports whether the cadence is one of the spec-normative values
// per §10.27. Invalid values include the empty string, case variants
// (e.g. "Daily"), and any string outside the enumeration.
func (c Cadence) IsValid() bool {
	_, ok := cadenceProperties[c]
	return ok
}

// IsStreaming reports whether the cadence is sub-daily streaming-mode per
// §10.27 (per_second, per_minute, per_hour). Returns false for valid
// non-streaming cadences (hourly, daily, weekly) AND for invalid cadences;
// callers that need to distinguish "non-streaming" from "invalid" MUST
// also check IsValid.
func (c Cadence) IsStreaming() bool {
	return cadenceProperties[c].streaming
}

// Interval returns the time-duration interval the cadence implies, or 0
// if the cadence is not in the §10.27 enumeration. A 0 return is
// unambiguous because no valid cadence has a zero-duration interval.
//
// Note that CadencePerHour and CadenceHourly produce the same interval
// (1 hour); the spec-byte-form distinction is preserved by the wire-format
// value, not by the interval.
func (c Cadence) Interval() time.Duration {
	return cadenceProperties[c].interval
}
