# Contributing to configo

Contributions are welcome — the roadmap in the [README](./README.md#status) is a
good place to look for something to pick up.

## Getting set up

```shell
git clone https://github.com/Marfien/configo.git
cd configo
go test ./...
```

Go 1.26 or newer is required (the code uses `range` over integers and
`strings.SplitSeq`). The only dependency is
[testify](https://github.com/stretchr/testify), used for assertions in tests.

## Working on a change

- Open an issue first for anything that changes the public API or the tag
  format, so the design can be agreed on before you write code.
- Keep commits focused and write a short, imperative subject line.
- Run `go test ./...`, `go vet ./...` and `gofmt -l .` before pushing; `gofmt`
  should print nothing. CI runs the same three commands plus `go mod tidy`, so a
  clean local run means a green pull request.
- Note user-visible changes in [CHANGELOG.md](./CHANGELOG.md) under
  `## [Unreleased]`.

## Tests

The suite is table driven throughout — follow the existing shape rather than
introducing a new style:

- Config structs used by parser cases live in `parser_fixtures_test.go`, not
  inline in the test function.
- Each case is a struct in a `tests` slice with a `name`, run via
  `t.Run(tt.name, ...)`.
- Assertions use `github.com/stretchr/testify/assert`.

New behaviour needs a test. A new source needs coverage of the "not found" path
(`ErrNotFoundInSource`) as well as the happy path, since `Parse` relies on that
sentinel to fall through to the next source.

## Documentation

The package documentation on [pkg.go.dev](https://pkg.go.dev/github.com/Marfien/configo)
is generated from the comments in `doc.go` and from the doc comments on exported
identifiers. A change to the tag format, the source list or the supported field
types needs an update in `doc.go` *and* in the README, which cover the same
ground for different audiences.

## Releasing

For maintainers. A release *is* a git tag — there is no artifact to build and
nothing to upload. Pushing a `v*` tag runs
[`.github/workflows/release.yml`](./.github/workflows/release.yml), which
validates the tag, re-runs the full suite against it, creates the GitHub release
from the matching `CHANGELOG.md` section, and primes the module proxy so
pkg.go.dev picks up the docs within minutes rather than on its next poll.

```shell
# 1. Turn "## [Unreleased]" into "## [0.1.0] - YYYY-MM-DD" and commit
git commit -am "chore: release v0.1.0"
git push origin main

# 2. Tag and push; the workflow does the rest
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

Two things to know before you push a tag:

- **A published version is immutable.** `proxy.golang.org` caches a tag's
  content permanently on first fetch, so moving or deleting a tag does not take
  effect — anyone who already fetched keeps the old bytes, and anyone who has not
  gets a checksum mismatch. A bad release is fixed by publishing the next patch
  version, never by retagging.
- **`v2` and above need a module path suffix.** `go.mod` must declare
  `module github.com/Marfien/configo/v2`. The workflow refuses a mismatched tag
  before anything fetches it, which is the only point where the mistake is still
  cheap. Staying on `0.x` while the tag format and merge semantics are unsettled
  avoids the question entirely.

Tags of the form `v1.2.3-rc.1` are marked as prereleases automatically.

## Layout

| File               | Contents                                                         |
| ------------------ | ---------------------------------------------------------------- |
| `doc.go`           | Package documentation only                                       |
| `configo.go`       | Package level `Parse`, the `Environment` and `Default` sources    |
| `types.go`         | The `Source` interface and `ErrNotFoundInSource`                  |
| `parser.go`        | `Parser`, options, type dispatch and reflection-based assignment  |
| `explorer.go`      | Struct traversal: tag parsing, nesting, parent key paths          |
| `string_source.go` | Shared string parsing for string-based sources                    |
| `cli_source.go`    | `NewCliSource` and argument parsing                               |
| `map_source.go`    | `NewMapSource`                                                   |
