#!/usr/bin/env bash
# Runs the verifier's conformance vectors and fails on any failure or any skip.
#
# The vector tests skip when a vector's files are absent, so a green run alone
# cannot show that the vectors ran. This script makes a skip a failure.
#
# Usage: FFIEC_PUBLIC_VECTORS_DIR=<ffiec-public>/spec/test-vectors tools/check-conformance.sh
set -u

if [ -z "${FFIEC_PUBLIC_VECTORS_DIR:-}" ]; then
  echo "conformance: FFIEC_PUBLIC_VECTORS_DIR is not set" >&2
  exit 2
fi
if [ ! -d "$FFIEC_PUBLIC_VECTORS_DIR" ]; then
  echo "conformance: no such folder: $FFIEC_PUBLIC_VECTORS_DIR" >&2
  exit 2
fi

log="$(mktemp)"
go test -count=1 -v ./verifier/internal/vectors/ >"$log" 2>&1
status=$?

passed=$(grep -c -E '^ *--- PASS' "$log")
failed=$(grep -c -E '^ *--- FAIL' "$log")
skipped=$(grep -c -E '^ *--- SKIP' "$log")

grep -E 'corpus gate:|conformance gate:' "$log"
grep -E '^ *--- (FAIL|SKIP)' "$log"
echo "conformance: $passed passed, $failed failed, $skipped skipped, against $FFIEC_PUBLIC_VECTORS_DIR"

if [ "$status" -ne 0 ] || [ "$failed" -ne 0 ]; then
  echo "conformance: a test failed" >&2
  tail -n 40 "$log" >&2
  exit 1
fi
if [ "$skipped" -ne 0 ]; then
  echo "conformance: a skip means a vector did not run" >&2
  exit 1
fi
if [ "$passed" -eq 0 ]; then
  echo "conformance: no test ran" >&2
  exit 1
fi
