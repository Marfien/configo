package configo

import "os"

var Environment Source = stringSource{
	tag: "env",
	lookupString: func(key string, _ []string) (string, bool) {
		return os.LookupEnv(key)
	},
}

var Default Source = stringSource{
	tag: "default",
	lookupString: func(key string, _ []string) (string, bool) {
		return key, true
	},
}

var defaultParser = NewParser(
	WithSources(
		Default,
		Environment,
	),
)

func Parse(cfg any) error {
	return defaultParser.Parse(cfg)
}
