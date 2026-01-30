package utils

import (
	"reflect"
	"strings"
)

func Map[T any, R any](arr []T, fn func(T) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

func StructToMap(obj any, tag ...string) map[string]any {
	result := make(map[string]any)

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return result
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()

	useTag := ""
	if len(tag) > 0 {
		useTag = tag[0]
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if !value.CanInterface() {
			continue
		}

		// Default = struct field name
		key := field.Name

		// If tag provided, use it
		if useTag != "" {
			tagValue := field.Tag.Get(useTag)
			if tagValue == "-" {
				continue
			}
			if tagValue != "" {
				key = strings.Split(tagValue, ",")[0]
			}
		}

		// Skip zero values (critical for UpdateItem)
		if value.IsZero() {
			continue
		}

		result[key] = value.Interface()
	}

	return result
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
