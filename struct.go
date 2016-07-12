package excavator

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/seletskiy/hierr"
)

func excavateStruct(receiverValue, dataValue reflect.Value) error {
	switch dataValue.Type().Kind() {
	case reflect.Map:
		return excavateMapToStruct(receiverValue, dataValue)
	case reflect.Struct:
		// @TODO convert struct to struct
		return errors.New("not implemented")
	default:
		return newConvertError(receiverValue, dataValue)
	}

}

func excavateMapToStruct(receiverValue, dataValue reflect.Value) error {
	var (
		structType   = receiverValue.Type()
		fieldsNumber = structType.NumField()

		zero reflect.Value
	)

	//mapElemType := dataType.Elem()

	if dataValue.Type().Key().Kind() != reflect.String {
		return fmt.Errorf(
			"%s not supported, can work only with string keys",
			dataValue.Type(),
		)
	}

	for index := 0; index < fieldsNumber; index++ {
		field := structType.Field(index)
		log.Println("processing", field.Name)

		if field.Anonymous {
			continue // skip anonimous fields
		}

		fieldValue := getFieldFromMap(field, dataValue)
		if fieldValue == zero {
			log.Printf("'%s' not found", field.Name)
			continue
		}

		newFieldValue := reflect.New(field.Type).Elem()
		err := excavate(newFieldValue, unrollValue(fieldValue))
		if err != nil {
			return hierr.Errorf(err, "can't excavate field '%s'", field.Name)
		}

		receiverValue.Field(index).Set(newFieldValue)
	}

	return nil
}
