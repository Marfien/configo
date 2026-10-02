package configo

import (
	"fmt"
	"strconv"
	"strings"
)

type stringSource struct {
	tag          string
	lookupString func(key string, parents []string) (string, bool)
}

func (s stringSource) Tag() string {
	return s.tag
}

func (s stringSource) GetString(key string, parents []string) (string, error) {
	val, ok := s.lookupString(key, parents)
	if ok {
		return val, nil
	} else {
		return val, ErrNotFoundInSource
	}
}

func (s stringSource) GetBool(key string, parents []string) (bool, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return false, err
	}
	val, err := strconv.ParseBool(sVal)
	return val, err
}

func (s stringSource) GetInt(key string, parents []string) (int, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(sVal)
	return val, err
}

func (s stringSource) GetUint(key string, parents []string) (uint, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(sVal, 10, 0)
	return uint(v), err
}

func (s stringSource) GetFloat(key string, parents []string) (float64, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseFloat(sVal, 64)
	return v, err
}

func (s stringSource) GetStringSlice(key string, parents []string) ([]string, error) {
	sVal, err := s.GetString(key, parents)
	if err != nil {
		return nil, err
	}

	return strings.Split(sVal, ","), nil
}

func (s stringSource) GetBoolSlice(key string, parents []string) ([]bool, error) {
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

func (s stringSource) GetIntSlice(key string, parents []string) ([]int, error) {
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

func (s stringSource) GetUintSlice(key string, parents []string) ([]uint, error) {
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

func (s stringSource) GetFloatSlice(key string, parents []string) ([]float64, error) {
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

func (s stringSource) GetStringMap(key string, parents []string) (map[string]string, error) {
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

func (s stringSource) GetBoolMap(key string, parents []string) (map[string]bool, error) {
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

func (s stringSource) GetIntMap(key string, parents []string) (map[string]int, error) {
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

func (s stringSource) GetUintMap(key string, parents []string) (map[string]uint, error) {
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

func (s stringSource) GetFloatMap(key string, parents []string) (map[string]float64, error) {
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
