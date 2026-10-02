# configo — Configuration with Go

[![CI](https://github.com/Marfien/configo/actions/workflows/ci.yml/badge.svg)](https://github.com/Marfien/configo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Marfien/configo.svg)](https://pkg.go.dev/github.com/Marfien/configo)
[![Go Report Card](https://goreportcard.com/badge/github.com/Marfien/configo)](https://goreportcard.com/report/github.com/Marfien/configo)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

configo supplies you with a tag driven configuration framework.

## Abstract

configo fills a config struct from any number of _sources_ — environment variables,
CLI arguments, in-memory maps, or anything you implement yourself. Every field
declares, per source, the key it should be read from:

```go
type Config struct {
    Host string `configo:"env=APP_HOST,cli=--host,default=localhost"`
}
```

There is no central key schema and no reflection-driven naming convention: the
struct tag _is_ the mapping. Sources are consulted in the order you register them
and the last one that supplies a value wins, which gives you layered
configuration (defaults → file → environment → flags) by simply ordering a slice.
Fields that no source can supply are collected and reported together, so a
misconfigured deployment fails once with a complete list instead of one key at a
time.

### Status

Early stage. The pieces below work and are covered by tests; the roadmap items
are not implemented yet.

- [x] Value sources with priorities
- [x] Struct tag driven key mapping, per source
- [x] Nested structs (parent key paths handed to the source)
- [x] Scalars, pointers, slices and string-keyed maps
- [ ] Merging of sources (currently last-writer-wins, not a deep merge)
- [ ] Validators

## Quickstart

```shell
go get github.com/Marfien/configo
```

```go
package main

import (
	"fmt"
	"log"

	"github.com/Marfien/configo"
)

type Config struct {
	Host    string `configo:"env=APP_HOST,default=localhost"`
	Port    int    `configo:"env=APP_PORT,default=8080"`
	Verbose bool   `configo:"env=APP_VERBOSE,default=false"`
}

func main() {
	var cfg Config
	if err := configo.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

```shell
$ go run .
{Host:localhost Port:8080 Verbose:false}

$ APP_HOST=example.com APP_PORT=9090 go run .
{Host:example.com Port:9090 Verbose:false}
```

The package level `configo.Parse` uses `Default` followed by `Environment`, so
environment variables override defaults. For anything else, build your own
parser (see below).

## Usage

### The `configo` tag

A tag is a comma separated list of `sourceTag=key` pairs:

```go
type Config struct {
    Host string `configo:"env=APP_HOST,cli=--host,default=localhost"`
}
```

- `sourceTag` selects the source — it must match that source's `Tag()`.
- `key` is the lookup key for _that_ source, so each source can use its own
  naming convention (`APP_HOST`, `--host`, `app.host`).
- A source that the field does not mention is skipped for that field.
- Fields without a `configo` tag and unexported fields are ignored.
- A pair without `=` is a malformed tag and `Parse` returns an error.

### Priorities

`Parse` walks the sources in registration order and assigns every value it
finds, so **the last source that resolves a field wins**. Order the slice from
lowest to highest priority:

```go
parser := configo.NewParser(
    configo.WithSources(
        configo.Default,                                  // lowest
        configo.NewMapSource("file", valuesFromYAML),
        configo.Environment,
        configo.NewCliSource("cli", os.Args[1:]),         // highest
    ),
)

var cfg Config
if err := parser.Parse(&cfg); err != nil {
    log.Fatal(err)
}
```

If no source supplies a value for a field, `Parse` continues with the remaining
fields and finally returns a `missing values: Host, Port` error naming every
unsatisfied field. A source that _finds_ a value but fails to convert it (for
example `APP_PORT=abc` for an `int`) aborts immediately with that error.

### Built-in sources

| Source                                      | Tag       | Description                                                                  |
| ------------------------------------------- | --------- | ---------------------------------------------------------------------------- |
| `configo.Environment`                       | `env`     | Reads `os.LookupEnv(key)`.                                                   |
| `configo.Default`                           | `default` | Echoes the key back as the value, so `default=8080` _is_ the literal `8080`. |
| `configo.NewCliSource(tag, args)`           | yours     | Parses command line arguments.                                               |
| `configo.NewMapSource(tag, map[string]any)` | yours     | Serves pre-typed values from a map.                                          |

`NewCliSource` accepts `--key=value`, `--key value` and bare `--flag` (which
yields `"true"`). A repeated key is joined with commas, which makes
`--tag a --tag b` usable as a slice. Note that the key in your tag includes the
dashes: `configo:"cli=--host"`. Anything starting with `-` is treated as the next
key, so negative numbers must be written as `--retries=-1`.

`MapSource` returns its values as-is and does not convert between types: an
`int` field needs an `int` in the map, a `float64` field a `float64`.

### Supported field types

`bool`, all sized and unsized `int`/`uint` variants, `float32`, `float64`,
`string`, pointers to those, slices of those, and maps with `string` keys and
those element types. Pointer fields are allocated on demand. Any other type
(channels, functions, nested pointers to structs, non-string map keys) is
rejected with an error.

String-based sources (environment, CLI, and anything built on them) split
collections on commas; maps additionally split each entry on the first `=`:

```go
type Config struct {
    Hosts   []string          `configo:"env=APP_HOSTS"`   // APP_HOSTS=a,b,c
    Weights map[string]int    `configo:"env=APP_WEIGHTS"` // APP_WEIGHTS=a=1,b=2
}
```

### Nested structs

A tagged struct field contributes its key to a _parent path_ that is handed to
the source alongside the leaf key, letting hierarchical sources address nested
values:

```go
type Database struct {
    Host string `configo:"yaml=host"`
}

type Config struct {
    Database Database `configo:"yaml=database"`
}
```

Here the source is asked for key `host` with parents `["database"]`. It is up to
the source to decide what that means — joining it into `database.host`, walking a
tree, and so on. If a nested struct is missing a source tag for a deeper nested
struct, a blank string is used instead. `Environment`, `Default` and `CliSource`
ignore the parent path and look up the leaf key directly.

### Custom sources

Implement the `Source` interface (`Tag()` plus a getter per supported type) to
plug in your own backend. If your source is fundamentally string based — files,
remote key-value stores, templates — the simplest route is to let configo do the
parsing for you and only supply the lookup. That is exactly how `Environment`
and `CliSource` are built. For this you may utilize `configo.StringSource` with
a provided lookup function. For an example take a look at the defintion of the `Environment` [source](./configo.go).

### Custom tag name

```go
parser := configo.NewParser(
    configo.WithTag("cfg"),
    configo.WithSources(configo.Environment),
)
```

## How to contribute

Contributions are welcome — the roadmap above is a good place to look for
something to pick up. [CONTRIBUTING.md](./CONTRIBUTING.md) covers the setup, the
expectations for a change, the shape of the test suite and the file layout.

Notable changes are recorded in [CHANGELOG.md](./CHANGELOG.md).

## License

[MIT](./LICENSE) © Marvin Haase
