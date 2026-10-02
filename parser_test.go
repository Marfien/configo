package configo

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    Parser
	}{
		{
			name: "defaults",
			want: Parser{
				tag:     defaultTag,
				sources: make([]Source, 0),
			},
		},
		{
			name:    "custom tag",
			options: []Option{WithTag("custom")},
			want: Parser{
				tag:     "custom",
				sources: make([]Source, 0),
			},
		},
		{
			name:    "with sources",
			options: []Option{WithSources(MapSource{})},
			want: Parser{
				tag:     defaultTag,
				sources: []Source{MapSource{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewParser(tt.options...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getterForType(t *testing.T) {
	var source MapSource

	tests := []struct {
		name    string
		typ     reflect.Type
		want    reflect.Type
		wantErr bool
	}{
		{
			name: "bool",
			typ:  reflect.TypeFor[bool](),
			want: reflect.TypeFor[bool](),
		},
		{
			name: "int",
			typ:  reflect.TypeFor[int](),
			want: reflect.TypeFor[int](),
		},
		{
			name: "int8",
			typ:  reflect.TypeFor[int8](),
			want: reflect.TypeFor[int](),
		},
		{
			name: "int32",
			typ:  reflect.TypeFor[int32](),
			want: reflect.TypeFor[int](),
		},
		{
			name: "int64",
			typ:  reflect.TypeFor[int64](),
			want: reflect.TypeFor[int](),
		},
		{
			name: "uint",
			typ:  reflect.TypeFor[uint](),
			want: reflect.TypeFor[uint](),
		},
		{
			name: "uint8",
			typ:  reflect.TypeFor[uint8](),
			want: reflect.TypeFor[uint](),
		},
		{
			name: "uint32",
			typ:  reflect.TypeFor[uint32](),
			want: reflect.TypeFor[uint](),
		},
		{
			name: "uint64",
			typ:  reflect.TypeFor[uint64](),
			want: reflect.TypeFor[uint](),
		},
		{
			name: "float32",
			typ:  reflect.TypeFor[float32](),
			want: reflect.TypeFor[float64](),
		},
		{
			name: "float64",
			typ:  reflect.TypeFor[float64](),
			want: reflect.TypeFor[float64](),
		},
		{
			name: "string",
			typ:  reflect.TypeFor[string](),
			want: reflect.TypeFor[string](),
		},
		{
			name: "pointer is dereferenced",
			typ:  reflect.TypeFor[*int](),
			want: reflect.TypeFor[int](),
		},
		{
			name: "bool slice",
			typ:  reflect.TypeFor[[]bool](),
			want: reflect.TypeFor[[]bool](),
		},
		{
			name: "int slice",
			typ:  reflect.TypeFor[[]int32](),
			want: reflect.TypeFor[[]int](),
		},
		{
			name: "uint slice",
			typ:  reflect.TypeFor[[]uint](),
			want: reflect.TypeFor[[]uint](),
		},
		{
			name: "float slice",
			typ:  reflect.TypeFor[[]float32](),
			want: reflect.TypeFor[[]float64](),
		},
		{
			name: "string slice",
			typ:  reflect.TypeFor[[]string](),
			want: reflect.TypeFor[[]string](),
		},
		{
			name: "bool map",
			typ:  reflect.TypeFor[map[string]bool](),
			want: reflect.TypeFor[map[string]bool](),
		},
		{
			name: "int map",
			typ:  reflect.TypeFor[map[string]int8](),
			want: reflect.TypeFor[map[string]int](),
		},
		{
			name: "uint map",
			typ:  reflect.TypeFor[map[string]uint](),
			want: reflect.TypeFor[map[string]uint](),
		},
		{
			name: "float map",
			typ:  reflect.TypeFor[map[string]float32](),
			want: reflect.TypeFor[map[string]float64](),
		},
		{
			name: "string map",
			typ:  reflect.TypeFor[map[string]string](),
			want: reflect.TypeFor[map[string]string](),
		},
		{
			name:    "error: non-string map key",
			typ:     reflect.TypeFor[map[int]string](),
			wantErr: true,
		},

		{
			name:    "error: unsupported map value",
			typ:     reflect.TypeFor[map[string]struct{}](),
			wantErr: true,
		},

		{
			name:    "error: unsupported slice element",
			typ:     reflect.TypeFor[[]struct{}](),
			wantErr: true,
		},

		{
			name:    "error: unsupported type",
			typ:     reflect.TypeFor[chan int](),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get := getterForType(tt.typ, source)
			val, err := get("dummy", nil)

			if tt.wantErr {
				assert.Error(t, err)
				assert.NotErrorIs(t, err, ErrNotFoundInSource)
				return
			}

			if err != nil {
				assert.ErrorIs(t, err, ErrNotFoundInSource)
			} else {
				assert.Equal(t, tt.want, reflect.TypeOf(val))
			}
		})
	}
}

func TestParser_Parse(t *testing.T) {
	parentsSource := stringSource{
		tag: "parents",
		lookupString: func(key string, parents []string) (string, bool) {
			return strings.Join(append(parents, key), "."), true
		},
	}

	tests := []struct {
		name        string
		options     []Option
		cfg         any
		want        any
		wantErr     bool
		errContains []string
	}{
		{
			name:    "error: not a struct pointer",
			cfg:     emptyConfig{},
			wantErr: true,
		},
		{
			name:    "error: maleformed tag",
			cfg:     &malformedTagConfig{},
			wantErr: true,
		},
		{
			name: "sets scalar fields",
			options: []Option{WithSources(NewMapSource("map", map[string]any{
				"str":   "value",
				"int":   -42,
				"uint":  uint(42),
				"float": 4.2,
				"bool":  true,
			}))},
			cfg: &scalarsConfig{},
			want: &scalarsConfig{
				Str:   "value",
				Int:   -42,
				Uint:  42,
				Float: 4.2,
				Bool:  true,
				Ptr:   new(-42),
			},
		},
		{
			name:    "nested struct uses parent keys",
			options: []Option{WithSources(parentsSource)},
			cfg:     &nestedConfig{},
			want:    &nestedConfig{Inner: nestedInnerConfig{Test: "parent.test"}},
		},
		{
			name: "later source overrides earlier one",
			options: []Option{WithSources(
				NewMapSource("first", map[string]any{"test": "first"}),
				NewMapSource("second", map[string]any{"test": "second"}),
			)},
			cfg:  &twoSourcesConfig{},
			want: &twoSourcesConfig{Test: "second"},
		},
		{
			name: "falls back when source misses the value",
			options: []Option{WithSources(
				NewMapSource("first", map[string]any{"test": "first"}),
				NewMapSource("second", map[string]any{}),
			)},
			cfg:  &twoSourcesConfig{},
			want: &twoSourcesConfig{Test: "first"},
		},
		{
			name: "sources without a key for the field are skipped",
			options: []Option{WithSources(
				NewMapSource("first", map[string]any{"test": "first"}),
				NewMapSource("second", map[string]any{"test": "second"}),
			)},
			cfg:  &secondSourceOnlyConfig{},
			want: &secondSourceOnlyConfig{Test: "second"},
		},
		{
			name:        "error: no source provides a value",
			options:     []Option{WithSources(NewMapSource("map", map[string]any{}))},
			cfg:         &twoFieldsConfig{},
			wantErr:     true,
			errContains: []string{"Test", "Other"},
		},
		{
			name:        "error: source fails to parse",
			options:     []Option{WithSources(NewMapSource("map", map[string]any{"test": "not a number"}))},
			cfg:         &intConfig{},
			wantErr:     true,
			errContains: []string{"Test"},
		},
		{
			name:    "error: unsupported field type",
			options: []Option{WithSources(NewMapSource("map", map[string]any{"test": 1}))},
			cfg:     &chanConfig{},
			wantErr: true,
		},
		{
			// GetIntSlice returns []int, but the source holds a string
			name:    "error: value not assignable to field",
			options: []Option{WithSources(NewMapSource("map", map[string]any{"test": "1,2"}))},
			cfg:     &intSliceConfig{},
			wantErr: true,
		},
		{
			name: "custom tag",
			options: []Option{
				WithTag("custom"),
				WithSources(NewMapSource("map", map[string]any{"test": "value"})),
			},
			cfg:  &customTagConfig{},
			want: &customTagConfig{Test: "value"},
		},
		{
			name:    "untagged and unexported fields are ignored",
			options: []Option{WithSources(NewMapSource("map", map[string]any{"test": "value"}))},
			cfg:     &untaggedConfig{},
			want:    &untaggedConfig{Test: "value"},
		},
		{
			name: "no fields to fill",
			cfg:  &emptyConfig{},
			want: &emptyConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewParser(tt.options...).Parse(tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
				for _, msg := range tt.errContains {
					assert.ErrorContains(t, err, msg)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, tt.cfg)
		})
	}
}

func Test_setValue(t *testing.T) {
	tests := []struct {
		name    string
		destVal any
		val     any
		want    any
		wantErr bool
	}{
		{
			name:    "nil default",
			destVal: new(int),
			val:     nil,
			want:    0,
		},
		{
			name:    "error: incompatible",
			destVal: new(string),
			val:     struct{}{},
			wantErr: true,
		},
		{
			name:    "assignable",
			destVal: new(int),
			val:     67,
			want:    67,
		},
		{
			name:    "converable",
			destVal: new(int),
			val:     int32(67),
			want:    67,
		},
		{
			name:    "pointer assignable",
			destVal: new(*int),
			val:     67,
			want:    new(67),
		},
		{
			name:    "pointer converable",
			destVal: new(*int),
			val:     int32(67),
			want:    new(67),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest := reflect.ValueOf(tt.destVal).Elem()
			err := setValue(dest, tt.val)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.EqualValues(t, tt.want, dest.Interface())
		})
	}
}
