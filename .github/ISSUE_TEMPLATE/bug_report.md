---
name: Bug report
about: Report a defect in the chain-of-custody implementation
title: '[bug] '
labels: ['bug']
---

## Component

- [ ] `core` (primitive library)
- [ ] `ledger` (ingest server)
- [ ] `verifier` (offline verifier)
- [ ] `spec` (specification)
- [ ] Documentation
- [ ] Other

## Severity

- [ ] Critical (chain integrity defect — please follow [SECURITY.md](../../SECURITY.md) instead)
- [ ] High (verifier produces wrong result, but chain integrity intact)
- [ ] Medium (incorrect behavior in normal operation)
- [ ] Low (cosmetic, performance, or non-functional)

## Description

<!-- What is the defect? What did you expect to happen? -->

## Reproduction

```bash
# Minimal reproduction commands
```

## Test case

<!-- Attach a minimal Go test that exercises the defect. The test should fail before the fix and pass after. -->

```go
func TestReproduces(t *testing.T) {
  // ...
}
```

## Environment

- Version / commit hash:
- Go version (`go version`):
- OS and architecture:
- Platform (bare metal / VM / container / managed cloud):

## Auditor's-lens question

<!--
If a regulator examiner saw this defect, what would they say?
This question is the bar for whether the bug is critical or merely operational.
-->
