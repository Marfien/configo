package configo

import "fmt"

type MapSource struct {
	tag      string
	delegate map[string]any
}

func NewMapSource(tag string, delegate map[string]any) MapSource {
	return MapSource{
		tag:      tag,
		delegate: delegate,
	}
}

func (m MapSource) Tag() string {
	return m.tag
}

func get[T any](m MapSource, key string) (T, error) {
	val, found := m.delegate[key]
	if !found {
		var def T
		return def, ErrNotFoundInSource
	}

	if parsed, ok := val.(T); ok {
		return parsed, nil
	} else {
		var def T
		return def, fmt.Errorf("key %s: %s cannot be casted to %T", key, val, def)
	}
}

func (m MapSource) GetBool(key string, parents []string) (bool, error) {
	return get[bool](m, key)
}

func (m MapSource) GetString(key string, parents []string) (string, error) {
	return get[string](m, key)
}

func (m MapSource) GetInt(key string, parents []string) (int, error) {
	return get[int](m, key)
}

func (m MapSource) GetUint(key string, parents []string) (uint, error) {
	return get[uint](m, key)
}

func (m MapSource) GetFloat(key string, parents []string) (float64, error) {
	return get[float64](m, key)
}

func (m MapSource) GetBoolSlice(key string, parents []string) ([]bool, error) {
	return get[[]bool](m, key)
}

func (m MapSource) GetStringSlice(key string, parents []string) ([]string, error) {
	return get[[]string](m, key)
}

func (m MapSource) GetIntSlice(key string, parents []string) ([]int, error) {
	return get[[]int](m, key)
}

func (m MapSource) GetUintSlice(key string, parents []string) ([]uint, error) {
	return get[[]uint](m, key)
}

func (m MapSource) GetFloatSlice(key string, parents []string) ([]float64, error) {
	return get[[]float64](m, key)
}

func (m MapSource) GetBoolMap(key string, parents []string) (map[string]bool, error) {
	return get[map[string]bool](m, key)
}

func (m MapSource) GetStringMap(key string, parents []string) (map[string]string, error) {
	return get[map[string]string](m, key)
}

func (m MapSource) GetIntMap(key string, parents []string) (map[string]int, error) {
	return get[map[string]int](m, key)
}

func (m MapSource) GetUintMap(key string, parents []string) (map[string]uint, error) {
	return get[map[string]uint](m, key)
}

func (m MapSource) GetFloatMap(key string, parents []string) (map[string]float64, error) {
	return get[map[string]float64](m, key)
}
