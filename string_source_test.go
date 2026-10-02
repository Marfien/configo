package configo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestStringSource serves every key from a flat map, ignoring parents.
func newTestStringSource(values map[string]string) StringSource {
	return NewStringSource("string", func(key string, parents []string) (string, bool) {
		val, ok := values[key]
		return val, ok
	})
}

func TestStringSource_Tag(t *testing.T) {
	assert.Equal(t, "string", newTestStringSource(nil).Tag())
}

func TestStringSource_lookup(t *testing.T) {
	var gotKey string
	var gotParents []string

	s := NewStringSource("string", func(key string, parents []string) (string, bool) {
		gotKey = key
		gotParents = parents
		return "value", true
	})

	got, gotErr := s.GetString("key", []string{"parent1", "parent2"})

	assert.NoError(t, gotErr)
	assert.Equal(t, "value", got)
	assert.Equal(t, "key", gotKey)
	assert.Equal(t, []string{"parent1", "parent2"}, gotParents)
}

func TestStringSource_GetString(t *testing.T) {
	tests := []struct {
		name    string
		values  map[string]string
		key     string
		want    string
		wantErr error
	}{
		{
			name:   "found",
			values: map[string]string{"key": "value"},
			key:    "key",
			want:   "value",
		},
		{
			name:   "empty value is found",
			values: map[string]string{"key": ""},
			key:    "key",
			want:   "",
		},
		{
			name:    "missing",
			values:  map[string]string{"other": "value"},
			key:     "key",
			wantErr: ErrNotFoundInSource,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(tt.values).GetString(tt.key, nil)

			if tt.wantErr != nil {
				assert.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetBool(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    bool
		wantErr bool
	}{
		{name: "true", value: "true", want: true},
		{name: "false", value: "false", want: false},
		{name: "1", value: "1", want: true},
		{name: "0", value: "0", want: false},
		{name: "TRUE", value: "TRUE", want: true},
		{name: "invalid", value: "yes", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetBool("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetInt(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    int
		wantErr bool
	}{
		{name: "positive", value: "42", want: 42},
		{name: "negative", value: "-42", want: -42},
		{name: "zero", value: "0", want: 0},
		{name: "float", value: "4.2", wantErr: true},
		{name: "not a number", value: "abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetInt("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetUint(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    uint
		wantErr bool
	}{
		{name: "positive", value: "42", want: 42},
		{name: "zero", value: "0", want: 0},
		{name: "negative", value: "-1", wantErr: true},
		{name: "not a number", value: "abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetUint("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetFloat(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    float64
		wantErr bool
	}{
		{name: "decimal", value: "4.2", want: 4.2},
		{name: "integer", value: "42", want: 42},
		{name: "negative", value: "-4.2", want: -4.2},
		{name: "exponent", value: "4.2e2", want: 420},
		{name: "not a number", value: "abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetFloat("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{name: "single", value: "a", want: []string{"a"}},
		{name: "many", value: "a,b,c", want: []string{"a", "b", "c"}},
		{name: "empty is a single empty entry", value: "", want: []string{""}},
		{name: "trailing separator", value: "a,", want: []string{"a", ""}},
		{name: "whitespace is kept", value: "a, b", want: []string{"a", " b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetStringSlice("key", nil)

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetBoolSlice(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []bool
		wantErr bool
	}{
		{name: "single", value: "true", want: []bool{true}},
		{name: "many", value: "true,false,1,0", want: []bool{true, false, true, false}},
		{name: "one invalid entry", value: "true,yes", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetBoolSlice("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetIntSlice(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []int
		wantErr bool
	}{
		{name: "single", value: "42", want: []int{42}},
		{name: "many", value: "1,-2,3", want: []int{1, -2, 3}},
		{name: "one invalid entry", value: "1,abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetIntSlice("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetUintSlice(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []uint
		wantErr bool
	}{
		{name: "single", value: "42", want: []uint{42}},
		{name: "many", value: "1,2,3", want: []uint{1, 2, 3}},
		{name: "negative entry", value: "1,-2", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetUintSlice("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetFloatSlice(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []float64
		wantErr bool
	}{
		{name: "single", value: "4.2", want: []float64{4.2}},
		{name: "many", value: "1.5,-2.5,3", want: []float64{1.5, -2.5, 3}},
		{name: "one invalid entry", value: "1.5,abc", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetFloatSlice("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetStringMap(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]string
		wantErr bool
	}{
		{name: "single", value: "a=1", want: map[string]string{"a": "1"}},
		{name: "many", value: "a=1,b=2", want: map[string]string{"a": "1", "b": "2"}},
		{name: "value contains =", value: "a=k=v", want: map[string]string{"a": "k=v"}},
		{name: "empty value", value: "a=", want: map[string]string{"a": ""}},
		{name: "repeated key wins last", value: "a=1,a=2", want: map[string]string{"a": "2"}},
		{name: "missing =", value: "a", wantErr: true},
		{name: "one entry missing =", value: "a=1,b", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetStringMap("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetBoolMap(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]bool
		wantErr bool
	}{
		{name: "single", value: "a=true", want: map[string]bool{"a": true}},
		{name: "many", value: "a=true,b=0", want: map[string]bool{"a": true, "b": false}},
		{name: "one invalid entry", value: "a=true,b=yes", wantErr: true},
		{name: "missing =", value: "a", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetBoolMap("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetIntMap(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]int
		wantErr bool
	}{
		{name: "single", value: "a=42", want: map[string]int{"a": 42}},
		{name: "many", value: "a=1,b=-2", want: map[string]int{"a": 1, "b": -2}},
		{name: "one invalid entry", value: "a=1,b=abc", wantErr: true},
		{name: "missing =", value: "a", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetIntMap("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetUintMap(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]uint
		wantErr bool
	}{
		{name: "single", value: "a=42", want: map[string]uint{"a": 42}},
		{name: "many", value: "a=1,b=2", want: map[string]uint{"a": 1, "b": 2}},
		{name: "negative entry", value: "a=-1", wantErr: true},
		{name: "missing =", value: "a", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetUintMap("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringSource_GetFloatMap(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    map[string]float64
		wantErr bool
	}{
		{name: "single", value: "a=4.2", want: map[string]float64{"a": 4.2}},
		{name: "many", value: "a=1.5,b=-2.5", want: map[string]float64{"a": 1.5, "b": -2.5}},
		{name: "one invalid entry", value: "a=1.5,b=abc", wantErr: true},
		{name: "missing =", value: "a", wantErr: true},
		{name: "empty", value: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTestStringSource(map[string]string{"key": tt.value}).GetFloatMap("key", nil)

			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}

			assert.NoError(t, gotErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Every derived getter must surface a missing key as ErrNotFoundInSource.
func TestStringSource_missingKey(t *testing.T) {
	s := newTestStringSource(map[string]string{"other": "value"})

	getters := map[string]func() error{
		"GetString": func() error { _, err := s.GetString("key", nil); return err },
		"GetBool":   func() error { _, err := s.GetBool("key", nil); return err },
		"GetInt":    func() error { _, err := s.GetInt("key", nil); return err },
		"GetUint":   func() error { _, err := s.GetUint("key", nil); return err },
		"GetFloat":  func() error { _, err := s.GetFloat("key", nil); return err },

		"GetBoolSlice":   func() error { _, err := s.GetBoolSlice("key", nil); return err },
		"GetStringSlice": func() error { _, err := s.GetStringSlice("key", nil); return err },
		"GetIntSlice":    func() error { _, err := s.GetIntSlice("key", nil); return err },
		"GetUintSlice":   func() error { _, err := s.GetUintSlice("key", nil); return err },
		"GetFloatSlice":  func() error { _, err := s.GetFloatSlice("key", nil); return err },

		"GetBoolMap":   func() error { _, err := s.GetBoolMap("key", nil); return err },
		"GetStringMap": func() error { _, err := s.GetStringMap("key", nil); return err },
		"GetIntMap":    func() error { _, err := s.GetIntMap("key", nil); return err },
		"GetUintMap":   func() error { _, err := s.GetUintMap("key", nil); return err },
		"GetFloatMap":  func() error { _, err := s.GetFloatMap("key", nil); return err },
	}

	for name, get := range getters {
		t.Run(name, func(t *testing.T) {
			assert.ErrorIs(t, get(), ErrNotFoundInSource)
		})
	}
}
