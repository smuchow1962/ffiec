// Package wire holds the spec §4.4 + §5 + §6 wire-format types: file
// header, chain entry, seal record. These mirror the OTLP attribute
// model the SDK emits and the on-disk form the verifier reads.
//
// Spec authority: §4.4 (OpenTelemetry-native wire), §5 (wire format
// inclusion list), §6 (storage form), plus the per-entry stamp table
// in §4.1 and the seal-record schema in §4.2.
//
// Status: stub. The types land in Commit 3 of the verifier upgrade,
// alongside the §7 header preflight (steps 1-3a). The current
// bootstrap NDJSON shape in verifier/internal/verify/ledger.go is
// being hard-cutover (D-3, locked 2026-05-21) to this spec-conformant
// wire form; the existing simplified types are removed in Commit 3.
//
// Cross-implementation reference: the .NET reference at
// Herald.Compliance/Audit/Chain/{AuditChainEntry,AuditFileHeader,
// WireFormat}.cs is the structural model. Field names are spec-named
// snake_case on the wire, idiomatic Go CamelCase in the types.
package wire
