package utils

import "reflect"

func IsStructEmpty(s interface{}) bool {
	rv := reflect.ValueOf(s)
	if rv.Kind() != reflect.Struct {
		panic("isStructEmpty called with a non-struct type")
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.String() != "" {
			return false
		}
	}

	return true
}
