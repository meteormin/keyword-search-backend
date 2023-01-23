package utils

import "reflect"

func HasField(data interface{}, field string) bool {
	metaValue := reflect.ValueOf(data).Elem()

	fieldValue := metaValue.FieldByName(field)
	if fieldValue == (reflect.Value{}) {
		return false
	}

	return true
}
