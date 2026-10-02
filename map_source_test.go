package configo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// mapGetter wraps a typed getter so all of them can be driven by one table.
type mapGetter struct {
	name string
	get  func(m MapSource, key string) (any, error)
}

var mapGetters = []mapGetter{
	{"GetBool", func(m MapSource, k string) (any, error) { return m.GetBool(k, nil) }},
	{"GetString", func(m MapSource, k string) (any, error) { return m.GetString(k, nil) }},
	{"GetInt", func(m MapSource, k string) (any, error) { return m.GetInt(k, nil) }},
	{"GetUint", func(m MapSource, k string) (any, error) { return m.GetUint(k, nil) }},
	{"GetFloat", func(m MapSource, k string) (any, error) { return m.GetFloat(k, nil) }},

	{"GetBoolSlice", func(m MapSource, k string) (any, error) { return m.GetBoolSlice(k, nil) }},
	{"GetStringSlice", func(m MapSource, k string) (any, error) { return m.GetStringSlice(k, nil) }},
	{"GetIntSlice", func(m MapSource, k string) (any, error) { return m.GetIntSlice(k, nil) }},
	{"GetUintSlice", func(m MapSource, k string) (any, error) { return m.GetUintSlice(k, nil) }},
	{"GetFloatSlice", func(m MapSource, k string) (any, error) { return m.GetFloatSlice(k, nil) }},

	{"GetBoolMap", func(m MapSource, k string) (any, error) { return m.GetBoolMap(k, nil) }},
	{"GetStringMap", func(m MapSource, k string) (any, error) { return m.GetStringMap(k, nil) }},
	{"GetIntMap", func(m MapSource, k string) (any, error) { return m.GetIntMap(k, nil) }},
	{"GetUintMap", func(m MapSource, k string) (any, error) { return m.GetUintMap(k, nil) }},
	{"GetFloatMap", func(m MapSource, k string) (any, error) { return m.GetFloatMap(k, nil) }},
}

func TestMapSource_Tag(t *testing.T) {
	assert.Equal(t, "map", NewMapSource("map", nil).Tag())
}

func TestMapSource_get(t *testing.T) {
	m := NewMapSource("map", map[string]any{
		"bool":   true,
		"string": "value",
		"int":    42,
		"uint":   uint(42),
		"float":  4.2,

		"boolSlice":   []bool{true, false},
		"stringSlice": []string{"a", "b"},
		"intSlice":    []int{1, -2},
		"uintSlice":   []uint{1, 2},
		"floatSlice":  []float64{1.5, -2.5},

		"boolMap":   map[string]bool{"a": true},
		"stringMap": map[string]string{"a": "b"},
		"intMap":    map[string]int{"a": -1},
		"uintMap":   map[string]uint{"a": 1},
		"floatMap":  map[string]float64{"a": 1.5},
	})

	tests := []struct {
		name   string
		getter mapGetter
		key    string
		want   any
	}{
		{name: "bool", getter: mapGetters[0], key: "bool", want: true},
		{name: "string", getter: mapGetters[1], key: "string", want: "value"},
		{name: "int", getter: mapGetters[2], key: "int", want: 42},
		{name: "uint", getter: mapGetters[3], key: "uint", want: uint(42)},
		{name: "float", getter: mapGetters[4], key: "float", want: 4.2},

		{name: "bool slice", getter: mapGetters[5], key: "boolSlice", want: []bool{true, false}},
		{name: "string slice", getter: mapGetters[6], key: "stringSlice", want: []string{"a", "b"}},
		{name: "int slice", getter: mapGetters[7], key: "intSlice", want: []int{1, -2}},
		{name: "uint slice", getter: mapGetters[8], key: "uintSlice", want: []uint{1, 2}},
		{name: "float slice", getter: mapGetters[9], key: "floatSlice", want: []float64{1.5, -2.5}},

		{name: "bool map", getter: mapGetters[10], key: "boolMap", want: map[string]bool{"a": true}},
		{name: "string map", getter: mapGetters[11], key: "stringMap", want: map[string]string{"a": "b"}},
		{name: "int map", getter: mapGetters[12], key: "intMap", want: map[string]int{"a": -1}},
		{name: "uint map", getter: mapGetters[13], key: "uintMap", want: map[string]uint{"a": 1}},
		{name: "float map", getter: mapGetters[14], key: "floatMap", want: map[string]float64{"a": 1.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := tt.getter.get(m, tt.key)

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapSource_get_notFound(t *testing.T) {
	m := NewMapSource("map", map[string]any{"other": "value"})

	for _, getter := range mapGetters {
		t.Run(getter.name, func(t *testing.T) {
			got, gotErr := getter.get(m, "unknown")

			assert.ErrorIs(t, gotErr, ErrNotFoundInSource)
			assert.Zero(t, got)
		})
	}
}

// Values are returned as stored, so a mismatching type is an error instead of a conversion.
func TestMapSource_get_wrongType(t *testing.T) {
	m := NewMapSource("map", map[string]any{"key": struct{}{}})

	for _, getter := range mapGetters {
		t.Run(getter.name, func(t *testing.T) {
			got, gotErr := getter.get(m, "key")

			assert.Error(t, gotErr)
			assert.NotErrorIs(t, gotErr, ErrNotFoundInSource)
			assert.ErrorContains(t, gotErr, "key")
			assert.Zero(t, got)
		})
	}
}

// The map is flat: keys of nested structs are looked up without their parents.
func TestMapSource_get_ignoresParents(t *testing.T) {
	m := NewMapSource("map", map[string]any{"key": "value"})

	got, gotErr := m.GetString("key", []string{"parent1", "parent2"})

	assert.NoError(t, gotErr)
	assert.Equal(t, "value", got)
}

func TestMapSource_get_nilDelegate(t *testing.T) {
	m := NewMapSource("map", nil)

	got, gotErr := m.GetString("key", nil)

	assert.ErrorIs(t, gotErr, ErrNotFoundInSource)
	assert.Empty(t, got)
}
