package site

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// MarshalerMap implements a custom [MarshalToMap].
type MarshalerMap interface {
	MarshalMap() (map[string]any, error)
}

type options struct {
	key string
	fn  func(any) (any, bool)
}

func newOpt(key string, fn func(any) (any, bool)) options {
	return options{key, fn}
}

var opts = []options{
	newOpt("omitempty", func(v any) (any, bool) {
		if v == nil {
			return nil, false
		}
		refVal := reflect.ValueOf(v)
		if reflect.DeepEqual(v, reflect.Zero(refVal.Type()).Interface()) {
			return v, false
		}
		return v, true
	}),
	newOpt("string", func(v any) (any, bool) {
		if cv, ok := v.(fmt.Stringer); ok {
			return cv.String(), true
		}
		return fmt.Sprintf("%v", v), true
	}),
}

func getElem(v reflect.Value) (any, error) {
	if !v.CanInterface() {
		return nil, nil
	}
	val := v.Interface()
	if r, ok := val.(Record); ok {
		val = AsJSON(r)
	}
	if conv, ok := val.(MarshalerMap); ok {
		return conv.MarshalMap()
	}
	switch v.Kind() {
	case reflect.Struct:
		return MarshalToMap(val)
	case reflect.Pointer:
		if v.IsNil() {
			return nil, nil
		}
		return getElem(v.Elem())
	default:
		return val, nil
	}
}

// MarshalToMap transforms a struct into a map.
//
// If v is not a map, it returns nil.
//
// Implements [MarshalerMap] to have a custom behavior.
func MarshalToMap(v any) (map[string]any, error) {
	ref := reflect.ValueOf(v)
	switch ref.Kind() {
	case reflect.Struct:
	case reflect.Pointer:
		if ref.IsNil() {
			return nil, nil
		}
		return MarshalToMap(ref.Elem().Interface())
	default:
		return nil, nil
	}
	if conv, ok := v.(MarshalerMap); ok {
		return conv.MarshalMap()
	}
	refType := ref.Type()
	fields := ref.NumField()
	mp := make(map[string]any, fields)
	for i := range fields {
		field := ref.Field(i)
		fieldType := refType.Field(i)
		val, err := getElem(field)
		if err != nil {
			return nil, err
		}
		name := fieldType.Name
		data := strings.Split(refType.Field(i).Tag.Get("json"), ",")
		if len(data) > 0 {
			if len(data[0]) > 0 {
				name = data[0]
			}
			if name == "-" {
				continue
			}
			if len(data) > 1 {
				tagOpts := data[1:]
				ok := true
				i := 0
				for i < len(opts) && ok {
					opt := opts[i]
					if slices.Contains(tagOpts, opt.key) {
						val, ok = opt.fn(val)
					}
					i++
				}
				if !ok {
					continue
				}
			}
		}
		mp[name] = val
	}
	return mp, nil
}
