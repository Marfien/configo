package configo

import (
	"fmt"
	"reflect"
	"strings"
)

const defaultTag = "configo"

var defaultExplorer = explorer{
	tag: defaultTag,
}

type exploredField struct {
	name             string
	value            reflect.Value
	sourceKeys       map[string]string
	parentSourceKeys map[string][]string
}

type explorer struct {
	tag string
}

func (e explorer) explore(s any) ([]exploredField, error) {
	return e.exploreWithParents(s, make(map[string][]string))
}

func (e explorer) exploreWithParents(s any, parents map[string][]string) ([]exploredField, error) {
	t := reflect.TypeOf(s)

	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s must be of kind *struct", t)
	}

	if parents == nil {
		parents = make(map[string][]string)
	}

	elemT := t.Elem()
	elemV := reflect.ValueOf(s).Elem()

	numFields := elemT.NumField()
	fields := make([]exploredField, 0, numFields)

	for i := range numFields {
		fieldT := elemT.Field(i)
		if !fieldT.IsExported() {
			continue
		}

		tagLine, ok := fieldT.Tag.Lookup(e.tag)
		if !ok {
			continue
		}

		sourceKeys, err := buildSourceKeys(tagLine)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", fieldT.Name, err)
		}

		fieldV := elemV.Field(i)
		if fieldT.Type.Kind() != reflect.Struct {
			fields = append(fields, exploredField{
				name:             fieldT.Name,
				value:            fieldV,
				sourceKeys:       sourceKeys,
				parentSourceKeys: parents,
			})
			continue
		}

		nestedParents := make(map[string][]string)
		for tag, v := range sourceKeys {
			// TODO: handle missing entries in parents
			if parentsKeys, ok := parents[tag]; ok {
				nestedParents[tag] = append(parentsKeys, v)
			} else {
				nestedParents[tag] = []string{v}
			}
		}

		// add empty entries for those sources that are not present at this field, but may be deeper
		for tag, v := range parents {
			if _, ok := sourceKeys[tag]; !ok {
				nestedParents[tag] = append(v, "")
			}
		}

		nestedFields, err := e.exploreWithParents(fieldV.Addr().Interface(), nestedParents)
		if err != nil {
			return fields, fmt.Errorf("field %s: %w", fieldT.Name, err)
		}

		fields = append(fields, nestedFields...)
	}

	return fields, nil
}

func buildSourceKeys(tagLine string) (map[string]string, error) {
	sourceKeys := make(map[string]string)
	if len(tagLine) == 0 {
		return sourceKeys, nil
	}

	for arg := range strings.SplitSeq(tagLine, ",") {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("maleformed source-key: %s", arg)
		}

		sourceKeys[kv[0]] = kv[1]
	}

	return sourceKeys, nil
}
