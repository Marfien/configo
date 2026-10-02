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
