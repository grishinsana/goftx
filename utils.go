package goftx

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

func PrepareQueryParams(params interface{}) (map[string]string, error) {
	result := make(map[string]string)

	val := reflect.ValueOf(params).Elem()
	if val.Kind() != reflect.Struct {
		return result, nil
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("json")

		switch valueField.Kind() {
		case reflect.Ptr:
			if valueField.IsNil() {
				continue
			}
			result[tag] = fmt.Sprintf("%v", valueField.Elem().Interface())
		default:
			if valueField.IsZero() {
				return result, errors.Errorf("required field: %v", tag)
			}
			result[tag] = fmt.Sprintf("%v", valueField.Interface())
		}
	}

	return result, nil
}
