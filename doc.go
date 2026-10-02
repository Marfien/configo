// Package configo fills a configuration struct from any number of value
// sources — environment variables, CLI arguments, in-memory maps, or anything
// implementing [Source].
//
// Every field declares, per source, the key it should be read from. There is no
// central key schema and no naming convention: the struct tag is the mapping.
//
//	type Config struct {
//	    Host string `configo:"env=APP_HOST,cli=--host,default=localhost"`
//	}
//
// # The configo tag
//
// A tag is a comma separated list of sourceTag=key pairs. The sourceTag selects
// the source and must match that source's [Source.Tag]; the key is the lookup
// key for that source alone, so each source keeps its own naming convention. A
// source a field does not mention is skipped for that field. Fields without a
// configo tag and unexported fields are ignored. A pair without "=" is a
// malformed tag and [Parse] returns an error.
//
// # Priorities
//
// [Parser.Parse] walks the sources in registration order and assigns every
// value it finds, so the last source that resolves a field wins. Order the
// sources from lowest to highest priority:
//
//	parser := configo.NewParser(
//	    configo.WithSources(
//	        configo.Default,                          // lowest
//	        configo.NewMapSource("file", fromYAML),
//	        configo.Environment,
//	        configo.NewCliSource("cli", os.Args[1:]), // highest
//	    ),
//	)
//
//	var cfg Config
//	if err := parser.Parse(&cfg); err != nil {
//	    log.Fatal(err)
//	}
//
// Fields that no source can supply are collected and reported together, so a
// misconfigured deployment fails once with a complete list instead of one key at
// a time. A source that finds a value but fails to convert it (for example
// APP_PORT=abc for an int field) aborts immediately with that error.
//
// The package level [Parse] uses [Default] followed by [Environment], so
// environment variables override defaults.
//
// # Built-in sources
//
//   - [Environment] (tag "env") reads os.LookupEnv.
//   - [Default] (tag "default") echoes the key back as the value, so
//     default=8080 is the literal 8080.
//   - [NewCliSource] parses command line arguments.
//   - [NewMapSource] serves pre-typed values from a map.
//
// # Supported field types
//
// bool, all sized and unsized int/uint variants, float32, float64, string,
// pointers to those, slices of those, and maps with string keys and those
// element types. Pointer fields are allocated on demand. Any other type
// (channels, functions, pointers to structs, non-string map keys) is rejected
// with an error.
//
// String based sources split collections on commas; maps additionally split
// each entry on the first "=":
//
//	type Config struct {
//	    Hosts   []string       `configo:"env=APP_HOSTS"`   // APP_HOSTS=a,b,c
//	    Weights map[string]int `configo:"env=APP_WEIGHTS"` // APP_WEIGHTS=a=1,b=2
//	}
//
// # Nested structs
//
// A tagged struct field contributes its key to a parent path that is handed to
// the source alongside the leaf key, letting hierarchical sources address nested
// values:
//
//	type Database struct {
//	    Host string `configo:"yaml=host"`
//	}
//
//	type Config struct {
//	    Database Database `configo:"yaml=database"`
//	}
//
// Here the source is asked for key "host" with parents ["database"]. It is up to
// the source to decide what that means — joining it into "database.host",
// walking a tree, and so on. [Environment], [Default] and the source returned by
// [NewCliSource] ignore the parent path and look up the leaf key directly.
//
// # Custom sources
//
// Implement [Source] to plug in a new backend. A source reports a missing key by
// returning [ErrNotFoundInSource], which is what makes [Parser.Parse] fall
// through to the next source. For a fundamentally string based backend — files,
// remote key-value stores, templates — [NewStringSource] does all the parsing
// and conversion and only needs a [LookupFunc]; that is how [Environment] and
// [NewCliSource] are built.
package configo
