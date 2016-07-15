package excavator

import (
	"reflect"
	"strings"
)

const (
	// this tag is used, if nothing passed to Excavate function
	DefaultStructTag = "excavator"
)

func getFieldKey(
	field reflect.StructField, tag string,
) reflect.Value {
	tagName := field.Tag.Get(tag)
	if tagName != "" {
		return reflect.ValueOf(tagName)
	}

	return reflect.ValueOf(strings.ToLower(field.Name))
}

func getFieldFromMap(
	field reflect.StructField,
	dataValue reflect.Value,
	tag string,
) reflect.Value {
	tagName := field.Tag.Get(tag)
	if tagName != "" {
		return dataValue.MapIndex(reflect.ValueOf(tagName))
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
