package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type convertError struct {
	from reflect.Type
	to   reflect.Type
}

func newConvertError(receiver, data interface{}) error {
	return &convertError{
		from: unrollType(reflect.TypeOf(data)),
		to:   unrollType(reflect.TypeOf(receiver)),
	}

}

func (err *convertError) Error() string {
	return fmt.Sprintf(
		"can't convert '%s' to '%s'",
		err.from, err.to,
	)
}

func Export(receiver, data interface{}) error {
	receiverType := reflect.TypeOf(receiver)
	if receiverType.Kind() != reflect.Ptr {
		log.Printf("receiver: %#+v", receiver)
		log.Printf("data: %#+v", data)
		log.Printf("receiver kind: %#+v", reflect.TypeOf(receiver).Kind())
		log.Printf("data kind: %#+v", reflect.TypeOf(data).Kind())
		return errors.New("non ptr receiver")
	}

	receiverType = unrollType(receiverType)

	switch receiverType.Kind() {
	case reflect.Slice:
		return exportSlice(receiver, data)
	case reflect.Map:
		return exportMap(receiver, data)
	case reflect.Struct:
		return exportStruct(receiver, data)
	default:
		dataKind := unrollType(reflect.TypeOf(data)).Kind()
		if receiverType.Kind() != dataKind {
			return newConvertError(receiver, data)
		}

		receiveTo := unrollValue(reflect.ValueOf(receiver))
		receiveTo.Set(unrollValue(reflect.ValueOf(data)))
		return nil
	}
}

func export(receiver, data reflect.Value) error {
	receiverType = unrollType(receiver.Type())

	switch receiverType.Kind() {
	case reflect.Slice:
		return exportSlice(receiver, data)
	case reflect.Map:
		return exportMap(receiver, data)
	case reflect.Struct:
		return exportStruct(receiver, data)
	default:
		dataKind := unrollType(reflect.TypeOf(data)).Kind()
		if receiverType.Kind() != dataKind {
			return newConvertError(receiver, data)
		}

		receiveTo := unrollValue(reflect.ValueOf(receiver))
		receiveTo.Set(unrollValue(reflect.ValueOf(data)))
		return nil
	}
}

func exportSlice(receiver, data reflect.Value) error {
	if err := check(receiver, data, reflect.Slice); err != nil {
		return err
	}

	dataKind := unrollType(reflect.TypeOf(data)).Kind()
	if dataKind != reflect.Slice {
		return newConvertError(receiver, data)
	}

	dataValue := reflect.ValueOf(data)
	sliceType := reflect.TypeOf(receiver).Elem()
	newSlice := reflect.MakeSlice(
		sliceType, dataValue.Len(), dataValue.Cap(),
	)

	for index := 0; index < dataValue.Len(); index++ {
		sliceElem := unrollValue(dataValue.Index(index))
		if sliceElem.Kind() != sliceType.Elem().Kind() {
			return fmt.Errorf(
				"element #%d of data slice is %s, not %s",
				index, sliceElem.Kind(), sliceType.Elem().Kind(),
			)
		}

		newSliceElem := newSlice.Index(index)
		//newElemPointer := newSliceElem.Addr().Pointer()
		//dataElemPointer := dataValue.Index(index).Addr().Pointer()
		//		err := export(newElemPointer, dataElemPointer)
		//		if err != nil {
		//			return fmt.Errorf("can't set element #%d: %s", index, err)
		//		}
		newSliceElem.Set(sliceElem)
	}

	receiverValue := reflect.ValueOf(receiver).Elem()
	receiverValue.Set(newSlice)

	return nil
}

func exportMap(receiver, data interface{}) error {
	if err := check(receiver, data, reflect.Map); err != nil {
		return err
	}

	dataKind := unrollType(reflect.TypeOf(data)).Kind()
	if dataKind != reflect.Map {
		return newConvertError(receiver, data)
	}

	mapType := reflect.TypeOf(receiver).Elem()
	mapElemType := mapType.Elem().Kind()
	newMap := reflect.MakeMap(mapType)

	dataValue := reflect.ValueOf(data)
	keys := dataValue.MapKeys()
	for _, key := range keys {
		mapElem := unrollValue(dataValue.MapIndex(key))
		if mapElem.Type().Kind() != mapElemType {
			return fmt.Errorf(
				"element '%v' of data map is %s, not %s",
				key, mapElem.Type().Kind(), mapElemType,
			)
		}

		newMap.SetMapIndex(key, mapElem)
	}

	receiverValue := reflect.ValueOf(receiver).Elem()
	receiverValue.Set(newMap)

	return nil
}

func exportStruct(receiver, data interface{}) error {
	if err := check(receiver, data, reflect.Struct); err != nil {
		return err
	}

	dataValue := unrollValue(reflect.ValueOf(data))
	structValue := unrollValue(reflect.ValueOf(receiver))
	structType := unrollType(reflect.TypeOf(receiver))
	dataType := unrollType(reflect.TypeOf(data))
	switch dataType.Kind() {
	case reflect.Struct:
		// @TODO convert struct to struct
		return errors.New("not implemented")
	case reflect.Map:
		//mapElemType := dataType.Elem()
		zero := reflect.Value{}
		fieldsNumber := structValue.NumField()
		if dataType.Key().Kind() != reflect.String {
			return fmt.Errorf(
				"can't extract from '%s', can't operate with non string keys",
				dataType,
			)
		}

		for index := 0; index < fieldsNumber; index++ {
			field := structType.Field(index)
			if field.Anonymous {
				continue // skip anonimous fields
			}

			var (
				tagValue       reflect.Value
				nameValue      reflect.Value
				lowerNameValue reflect.Value
				fieldValue     reflect.Value
			)

			nameValue = dataValue.MapIndex(reflect.ValueOf(field.Name))
			lowerNameValue = dataValue.MapIndex(
				reflect.ValueOf(strings.ToLower(field.Name)),
			)
			tag := field.Tag.Get("zhash")
			if tag != "" {
				// if tag not empty, erasing named values
				nameValue = reflect.Value{}
				lowerNameValue = reflect.Value{}

				tagValue = dataValue.MapIndex(reflect.ValueOf(tag))
			}

			log.Printf("tag: %v", tagValue)
			log.Printf("Name: %v", nameValue)
			log.Printf("name: %v", lowerNameValue)
			switch {
			case tagValue != zero:
				fieldValue = tagValue
			case nameValue != zero:
				fieldValue = nameValue
			case lowerNameValue != zero:
				fieldValue = lowerNameValue
			}

			log.Printf("field value: %v", fieldValue)

			if fieldValue == zero {
				log.Printf("'%s' not found", field.Name)
				continue
			}

			fieldValue = unrollValue(fieldValue)

			if field.Type != fieldValue.Type() {
				return fmt.Errorf(
					"can't assign %s to field '%s' of type %s",
					fieldValue.Type(), field.Name, field.Type,
				)
			}

			structValue.Field(index).Set(fieldValue)
		}
	default:
		return newConvertError(receiver, data)
	}

	return nil
}

func check(receiver, data interface{}, kind reflect.Kind) error {
	if reflect.TypeOf(receiver).Kind() != reflect.Ptr {
		return errors.New("receiver is not a pointer")
	}

	receiverKind := reflect.TypeOf(receiver).Elem().Kind()
	if receiverKind != kind {
		return fmt.Errorf("receiver is not a pointer to %s", kind)
	}

	return nil
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

func main() {
	var receiver = []string{}

	b := "10"

	var data = []interface{}{
		"A", &b,
	}

	err := export(&receiver, data)
	if err != nil {
		log.Fatalf("can't export slice: %s", err)
	}

	log.Printf("result %#+v", receiver)

	/*s := struct {
		Field1 int `zhash:"fld"`
		Field2 string
	}{}

	err := exportStruct(&s, map[string]interface{}{
		"Field1": 12,
		"Field2": "lol",
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("result %#+v", s)*/

}
