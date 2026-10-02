package configo

import (
	"fmt"
	"reflect"
	"strings"
)

type Parser struct {
	tag     string
	sources []Source
}

type Option func(p *Parser)

func WithTag(tag string) Option {
	return func(p *Parser) {
		p.tag = tag
	}
}

func WithSources(sources ...Source) Option {
	return func(p *Parser) {
		p.sources = sources
	}
}

func NewParser(opts ...Option) Parser {
	p := Parser{
		tag:     defaultTag,
		sources: make([]Source, 0),
	}

	for _, o := range opts {
		o(&p)
	}

	return p
}

func (p Parser) Parse(cfg any) error {
	explorer := explorer{
		tag: p.tag,
	}
	fields, err := explorer.explore(cfg)
	if err != nil {
		return err
	}

	missingValues := make([]string, 0)

	for _, f := range fields {
		var found bool
		for _, source := range p.sources {
			key, ok := f.sourceKeys[source.Tag()]
			if !ok {
				continue
			}

			parents := f.parentSourceKeys[source.Tag()]

			get := getterForType(f.value.Type(), source)
			val, err := get(key, parents)

			if err == ErrNotFoundInSource {
				continue
			}

			found = true

			if err != nil {
				return fmt.Errorf("field %s with source %s: %w", f.name, source.Tag(), err)
			}

			if err := setValue(f.value, val); err != nil {
				return fmt.Errorf("set value of field %s: %w", f.name, err)
			}
		}
		if !found {
			missingValues = append(missingValues, f.name)
		}
	}

	if len(missingValues) > 0 {
		return fmt.Errorf("missing values: %s", strings.Join(missingValues, ", "))
	}
	return nil
}

func getterForType(t reflect.Type, source Source) func(string, []string) (any, error) {
	switch t.Kind() {
	case reflect.Pointer:
		return getterForType(t.Elem(), source)
	case reflect.Bool:
		return wrapFunc(source.GetBool)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return wrapFunc(source.GetInt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return wrapFunc(source.GetUint)
	case reflect.Float32, reflect.Float64:
		return wrapFunc(source.GetFloat)
	case reflect.String:
		return wrapFunc(source.GetString)

	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			return func(string, []string) (any, error) {
				return nil, fmt.Errorf("map keys may only of type string")
			}
		}

		switch t.Elem().Kind() {
		case reflect.Bool:
			return wrapFunc(source.GetBoolMap)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return wrapFunc(source.GetIntMap)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return wrapFunc(source.GetUintMap)
		case reflect.Float32, reflect.Float64:
			return wrapFunc(source.GetFloatMap)
		case reflect.String:
			return wrapFunc(source.GetStringMap)
		}

	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Bool:
			return wrapFunc(source.GetBoolSlice)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return wrapFunc(source.GetIntSlice)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return wrapFunc(source.GetUintSlice)
		case reflect.Float32, reflect.Float64:
			return wrapFunc(source.GetFloatSlice)
		case reflect.String:
			return wrapFunc(source.GetStringSlice)
		}
	}

	return func(string, []string) (any, error) {
		return nil, fmt.Errorf("type %s is not supported, yet", t)
	}
}

func wrapFunc[T any](f func(string, []string) (T, error)) func(string, []string) (any, error) {
	return func(s1 string, s2 []string) (any, error) {
		return f(s1, s2)
	}
}

func setValue(dest reflect.Value, val any) error {
	if val == nil {
		return nil
	}
	if !dest.IsValid() || !dest.CanSet() {
		return fmt.Errorf("destination is not settable")
	}

	src := reflect.ValueOf(val)
	srcType := src.Type()
	destType := dest.Type()

	if src.Type().AssignableTo(destType) {
		dest.Set(src)
		return nil
	}

	if src.Type().ConvertibleTo(destType) {
		dest.Set(src.Convert(destType))
		return nil
	}

	if destType.Kind() == reflect.Pointer {
		elemType := destType.Elem()

		if srcType.AssignableTo(elemType) {
			ptr := reflect.New(elemType)
			ptr.Elem().Set(src)
			dest.Set(ptr)
			return nil
		}

		if srcType.ConvertibleTo(elemType) {
			ptr := reflect.New(elemType)
			ptr.Elem().Set(src.Convert(elemType))
			dest.Set(ptr)
			return nil
		}
	}

	return fmt.Errorf("value of type %s cannot be assinged to field of type %s", srcType, destType)
}
