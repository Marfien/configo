package configo

import (
	"errors"
)

var ErrNotFoundInSource = errors.New("not found in source")

type Source interface {
	Tag() string
	GetBool(key string, parents []string) (bool, error)
	GetString(key string, parents []string) (string, error)
	GetInt(key string, parents []string) (int, error)
	GetUint(key string, parents []string) (uint, error)
	GetFloat(key string, parents []string) (float64, error)

	GetBoolSlice(key string, parents []string) ([]bool, error)
	GetStringSlice(key string, parents []string) ([]string, error)
	GetIntSlice(key string, parents []string) ([]int, error)
	GetUintSlice(key string, parents []string) ([]uint, error)
	GetFloatSlice(key string, parents []string) ([]float64, error)

	GetBoolMap(key string, parents []string) (map[string]bool, error)
	GetStringMap(key string, parents []string) (map[string]string, error)
	GetIntMap(key string, parents []string) (map[string]int, error)
	GetUintMap(key string, parents []string) (map[string]uint, error)
	GetFloatMap(key string, parents []string) (map[string]float64, error)
}
