<!--
Thanks for contributing. Keep the description short — the checklist below is the
part that saves a review round trip.
-->

## What this changes

<!-- One or two sentences. What behaviour is different after this PR? -->

## Why

<!--
Link the issue this implements: "Closes #123".
Changes to the public API or the tag format should have an issue agreeing on the
design first — see CONTRIBUTING.md.
-->

## Checklist

- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] `gofmt -l .` prints nothing
- [ ] New behaviour has a test, following the table driven shape of the existing
      suite (fixtures in `parser_fixtures_test.go`, cases in a `tests` slice)
- [ ] A new source covers the `ErrNotFoundInSource` path as well as the happy
      path — *N/A if this PR adds no source*
- [ ] `doc.go` and the README are updated if the tag format, the source list or
      the supported field types changed
- [ ] User-visible changes are noted in `CHANGELOG.md` under `## [Unreleased]`

## Breaking change?

<!--
Yes / No. If yes, spell out what existing code stops compiling or silently
resolves differently, so it can be captured in the changelog.
-->
