package utils

import (
	"fmt"
	"reflect"
)

func Map[T any, R any](arr []T, fn func(T) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v)
	}
	return result
}

func StructToMap(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)

	// If the object is a pointer, get the element it points to
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	// Ensure it's a struct
	if v.Kind() != reflect.Struct {
		fmt.Printf("StructToMap only accepts structs; got %T\n", obj)
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		// Optional: use json tags if present
		// jsonTag := t.Field(i).Tag.Get("json")
		// if jsonTag != "" {
		//     fieldName = jsonTag
		// }
		fieldValue := v.Field(i).Interface()
		result[fieldName] = fieldValue
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
