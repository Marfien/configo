package configo

// Config structs used by the table driven cases of TestParser_Parse.

type scalarsConfig struct {
	Str   string  `configo:"map=str"`
	Int   int     `configo:"map=int"`
	Uint  uint    `configo:"map=uint"`
	Float float64 `configo:"map=float"`
	Bool  bool    `configo:"map=bool"`
	Ptr   *int    `configo:"map=int"`
}

type nestedInnerConfig struct {
	Test string `configo:"parents=test"`
}

type nestedConfig struct {
	Inner nestedInnerConfig `configo:"parents=parent"`
}

type malformedTagConfig struct {
	Test string `configo:"invalid"`
}

type twoSourcesConfig struct {
	Test string `configo:"first=test,second=test"`
}

type secondSourceOnlyConfig struct {
	Test string `configo:"second=test"`
}

type twoFieldsConfig struct {
	Test  string `configo:"map=test"`
	Other string `configo:"map=other"`
}

type intConfig struct {
	Test int `configo:"map=test"`
}

type intSliceConfig struct {
	Test []int `configo:"map=test"`
}

type chanConfig struct {
	Test chan int `configo:"map=test"`
}

type customTagConfig struct {
	Test string `custom:"map=test"`
}

type untaggedConfig struct {
	Untagged string
	Test     string `configo:"map=test"`
}

type emptyConfig struct{}
