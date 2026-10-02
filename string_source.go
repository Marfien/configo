package configo

import (
	"fmt"
	"strconv"
	"strings"
)

// LookupFunc resolves a key to its raw string representation. The second return
// value reports whether the key was present at all; a found-but-empty value is
// ("", true). StringSource turns a false into ErrNotFoundInSource so the Parser
// falls through to the next source.
type LookupFunc func(key string, parents []string) (string, bool)

// StringSource adapts a string based backend to the Source interface. It handles
// all conversions itself: scalars via the strconv package, slices by splitting on
// commas, and maps by additionally splitting each entry on the first "=".
type StringSource struct {
	tag          string
	lookupString LookupFunc
}

// NewStringSource builds a Source that answers every lookup through lookup.
func NewStringSource(tag string, lookup LookupFunc) StringSource {
	return StringSource{
		tag:          tag,
		lookupString: lookup,
	}
}

func (s StringSource) Tag() string {
	return s.tag
}

func (s StringSource) GetString(key string, parents []string) (string, error) {
	val, ok := s.lookupString(key, parents)
	if ok {
		return val, nil
	} else {
		return val, ErrNotFoundInSource
	}
}

func (s StringSource) GetBool(key string, parents []string) (bool, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return false, err
	}
	val, err := strconv.ParseBool(sVal)
	return val, err
}

func (s StringSource) GetInt(key string, parents []string) (int, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(sVal)
	return val, err
}

func (s StringSource) GetUint(key string, parents []string) (uint, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(sVal, 10, 0)
	return uint(v), err
}

func (s StringSource) GetFloat(key string, parents []string) (float64, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(sVal, 64)
	return v, err
}

func (s StringSource) GetStringSlice(key string, parents []string) ([]string, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return nil, err
	}

	return strings.Split(sVal, ","), nil
}

func (s StringSource) GetBoolSlice(key string, parents []string) ([]bool, error) {
	raw, err := s.GetStringSlice(key, parents)
	if err != nil {
		return nil, err
	}

	res := make([]bool, len(raw))
	for i, v := range raw {
		if p, err := strconv.ParseBool(v); err != nil {
			return nil, err
		} else {
			res[i] = p
		}
	}
	return res, nil
}

func (s StringSource) GetIntSlice(key string, parents []string) ([]int, error) {
	raw, err := s.GetStringSlice(key, parents)
	if err != nil {
		return nil, err
	}

	res := make([]int, len(raw))
	for i, v := range raw {
		if p, err := strconv.Atoi(v); err != nil {
			return nil, err
		} else {
			res[i] = p
		}
	}
	return res, nil
}

func (s StringSource) GetUintSlice(key string, parents []string) ([]uint, error) {
	raw, err := s.GetStringSlice(key, parents)
	if err != nil {
		return nil, err
	}

	res := make([]uint, len(raw))
	for i, v := range raw {
		if p, err := strconv.ParseUint(v, 10, 0); err != nil {
			return nil, err
		} else {
			res[i] = uint(p)
		}
	}
	return res, nil
}

func (s StringSource) GetFloatSlice(key string, parents []string) ([]float64, error) {
	raw, err := s.GetStringSlice(key, parents)
	if err != nil {
		return nil, err
	}

	res := make([]float64, len(raw))
	for i, v := range raw {
		if p, err := strconv.ParseFloat(v, 64); err != nil {
			return nil, err
		} else {
			res[i] = p
		}
	}
	return res, nil
}

func (s StringSource) GetStringMap(key string, parents []string) (map[string]string, error) {
	raw, err := s.GetStringSlice(key, parents)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(raw))
	for _, v := range raw {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) != 2 {
			return res, fmt.Errorf("entry '%s' cannot be interpreted as key-value pair (missing =)", v)
		}
		res[kv[0]] = kv[1]
	}

	return res, nil
}

func (s StringSource) GetBoolMap(key string, parents []string) (map[string]bool, error) {
	sMap, err := s.GetStringMap(key, parents)
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	for k, v := range sMap {
		if p, err := strconv.ParseBool(v); err != nil {
			return res, err
		} else {
			res[k] = p
		}
	}

	return res, nil
}

func (s StringSource) GetIntMap(key string, parents []string) (map[string]int, error) {
	sMap, err := s.GetStringMap(key, parents)
	if err != nil {
		return nil, err
	}

	res := make(map[string]int)
	for k, v := range sMap {
		if p, err := strconv.Atoi(v); err != nil {
			return res, err
		} else {
			res[k] = p
		}
	}

	return res, nil
}

func (s StringSource) GetUintMap(key string, parents []string) (map[string]uint, error) {
	sMap, err := s.GetStringMap(key, parents)
	if err != nil {
		return nil, err
	}

	res := make(map[string]uint)
	for k, v := range sMap {
		if p, err := strconv.ParseUint(v, 10, 0); err != nil {
			return res, err
		} else {
			res[k] = uint(p)
		}
	}

	return res, nil
}

func (s StringSource) GetFloatMap(key string, parents []string) (map[string]float64, error) {
	sMap, err := s.GetStringMap(key, parents)
	if err != nil {
		return nil, err
	}

	res := make(map[string]float64)
	for k, v := range sMap {
		if p, err := strconv.ParseFloat(v, 64); err != nil {
			return res, err
		} else {
			res[k] = p
		}
	}

	return res, nil
}
