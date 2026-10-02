package configo

import (
	"fmt"
	"strings"
)

var (
// CliSource      = NewCliSource("cli", os.Args[1:])
// CliShortSource = NewCliSource("scli", os.Args[1:])
)

func NewCliSource(tag string, args []string) Source {
	namedArgs := parseArgs(args)
	return NewStringSource(tag, func(key string, parents []string) (string, bool) {
		val, ok := namedArgs[key]
		return val, ok
	})
}

func parseArgs(args []string) map[string]string {
	namedArgs := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		kv := strings.SplitN(arg, "=", 2)
		name := kv[0]

		// key=value
		if len(kv) == 2 {
			setOrAppendToString(namedArgs, name, kv[1])
			continue
		}

		// --flag
		next := i + 1
		if next >= len(args) || strings.HasPrefix(args[next], "-") {
			setOrAppendToString(namedArgs, name, "true")
			continue
		}

		// --key value
		setOrAppendToString(namedArgs, name, args[next])
		i++
	}

	return namedArgs
}

func setOrAppendToString(m map[string]string, key, value string) {
	if val, ok := m[key]; ok {
		m[key] = fmt.Sprintf("%s,%s", val, value)
	} else {
		m[key] = value
	}
}
