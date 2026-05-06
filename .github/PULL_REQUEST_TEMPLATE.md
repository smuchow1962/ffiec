## Summary

<!-- One paragraph: what changed and why. Cite the regulatory or operational driver if applicable. -->

## Type of change

- [ ] Bug fix (non-spec)
- [ ] Implementation behavior change (spec-conforming)
- [ ] Specification change (requires unanimous spec-editor approval, 14-day public comment)
- [ ] Documentation only
- [ ] Build / CI / governance
- [ ] Security fix (private review preferred — see SECURITY.md)

## The auditor's-lens summary

<!--
One paragraph explaining how a reviewer can verify this change without re-running it themselves.
Where applicable, cite:
  - The handbook section or regulatory framework satisfied
  - The cryptographic primitive used (and its FIPS/NIST status)
  - The deterministic output guarantees
  - The new or modified test vectors that demonstrate the property
-->

## Test vectors

- [ ] No new test vectors required (no behavior change)
- [ ] Test vectors added in `spec/test-vectors/` and pass `go test ./spec/test-vectors/...`
- [ ] Test vectors updated to reflect spec change

## Spec / docs updates

- [ ] No spec change
- [ ] `spec/chain-of-custody-vN.md` updated with normative changes
- [ ] `docs/design/` updated to reflect implementation behavior
- [ ] README / examples updated

## Checklist

- [ ] Commits are signed (`git commit -S`)
- [ ] `go test -race ./...` passes locally
- [ ] `gofmt` clean
- [ ] `golangci-lint run` clean
- [ ] No new dependencies on closed-source or non-permissive code
- [ ] No telemetry / phone-home behavior introduced
