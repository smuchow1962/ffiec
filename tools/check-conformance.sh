#!/usr/bin/env bash
# Runs the conformance vectors and fails unless the run read what the pinned tree holds.
#
# The vector tests skip when a vector's files are absent, and some fall back to a
# contract-only check, so a green run alone cannot show that the vectors ran. This
# script fails on any skip, and it compares each corpus gate's line with
# tools/conformance-expected.txt. That file goes with the pin in .github/workflows/ci.yml:
# move one, and update the other.
#
# Usage: FFIEC_PUBLIC_VECTORS_DIR=<ffiec-public>/spec/test-vectors tools/check-conformance.sh
set -u

here="$(cd "$(dirname "$0")" && pwd)"
expected="$here/conformance-expected.txt"
jcs_expected=61

if [ -z "${FFIEC_PUBLIC_VECTORS_DIR:-}" ]; then
  echo "conformance: FFIEC_PUBLIC_VECTORS_DIR is not set" >&2
  exit 2
fi
# The tests fall back to a sibling checkout when this folder has no fixture, which
# would run a tree other than the one named here.
if [ ! -f "$FFIEC_PUBLIC_VECTORS_DIR/chain_vectors.json" ]; then
  echo "conformance: no chain_vectors.json in: $FFIEC_PUBLIC_VECTORS_DIR" >&2
  exit 2
fi

log="$(mktemp)"
go test -count=1 -v ./verifier/internal/vectors/ ./core/jcs/ >"$log" 2>&1
status=$?

passed=$(grep -c -E '^ *--- PASS' "$log")
failed=$(grep -c -E '^ *--- FAIL' "$log")
skipped=$(grep -c -E '^ *--- SKIP' "$log")
jcs=$(grep -c -E '^ *--- PASS: TestJCS_008Corpus/' "$log")

grep -E 'corpus gate:|conformance gate:|live-walk:' "$log"
grep -E '^ *--- (FAIL|SKIP)' "$log"
echo "conformance: $passed passed, $failed failed, $skipped skipped, $jcs JCS corpus cases, against $FFIEC_PUBLIC_VECTORS_DIR"

if [ "$status" -ne 0 ] || [ "$failed" -ne 0 ]; then
  echo "conformance: a test failed" >&2
  tail -n 40 "$log" >&2
  exit 1
fi
if [ "$skipped" -ne 0 ]; then
  echo "conformance: a skip means a vector did not run" >&2
  exit 1
fi

missing=0
while IFS= read -r line || [ -n "$line" ]; do
  line="${line%$'\r'}"
  [ -z "$line" ] && continue
  if ! grep -q -F -- "$line" "$log"; then
    echo "conformance: expected line is absent: $line" >&2
    missing=1
  fi
done <"$expected"
if [ "$missing" -ne 0 ]; then
  exit 1
fi
if [ "$jcs" -ne "$jcs_expected" ]; then
  echo "conformance: $jcs JCS corpus cases ran, expected $jcs_expected" >&2
  exit 1
fi
