package extractor

import (
	"reflect"
	"strings"
)

const (
	structTag = "zhash"
)

func getFieldFromMap(
	field reflect.StructField,
	dataValue reflect.Value,
) reflect.Value {
	tag := field.Tag.Get(structTag)
	if tag != "" {
		return dataValue.MapIndex(reflect.ValueOf(tag))
	}

	nameKey := reflect.ValueOf(field.Name)
	value := dataValue.MapIndex(nameKey)
	if value != (reflect.Value{}) {
		return value
	}

	lowerNameKey := reflect.ValueOf(strings.ToLower(field.Name))
	return dataValue.MapIndex(lowerNameKey)
}

func unrollType(value reflect.Type) reflect.Type {
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}

func unrollValue(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return value
}
