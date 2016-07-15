package excavator

import (
	"reflect"

	"github.com/seletskiy/hierr"
)

func excavateMap(
	receiverValue, dataValue reflect.Value, tag string,
) error {
	dataKind := unrollType(dataValue.Type()).Kind()

	switch dataKind {
	case reflect.Map:
		return excavateMapFromMap(receiverValue, dataValue, tag)
	case reflect.Struct:
		return excavateMapFromStruct(receiverValue, dataValue, tag)
	default:
		return newConvertError(receiverValue, dataValue)
	}
}

func excavateMapFromStruct(
	receiverValue, dataValue reflect.Value, tag string,
) error {
	var (
		mapType     = receiverValue.Type()
		mapElemType = mapType.Elem()
		newMap      = reflect.MakeMap(mapType)

		dataType     = dataValue.Type()
		fieldsNumber = dataType.NumField()
	)

	for index := 0; index < fieldsNumber; index++ {
		field := dataType.Field(index)
		key := getFieldKey(field, tag)

		dataFieldValue := unrollValue(dataValue.Field(index))
		mapElem := reflect.New(mapElemType).Elem()

		err := excavate(mapElem, dataFieldValue, tag)
		if err != nil {
			return hierr.Errorf(
				err, "can't excavate element '%v'", key,
			)
		}

		newMap.SetMapIndex(key, mapElem)
	}

	receiverValue.Set(newMap)
	return nil
}

func excavateMapFromMap(
	receiverValue, dataValue reflect.Value, tag string,
) error {
	mapType := receiverValue.Type()
	mapElemType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)
	for _, key := range dataValue.MapKeys() {
		dataMapElem := unrollValue(dataValue.MapIndex(key))
		mapElem := reflect.New(mapElemType).Elem()

		err := excavate(mapElem, dataMapElem, tag)
		if err != nil {
			return hierr.Errorf(
				err, "can't excavate element '%v'", key,
			)
		}

		newMap.SetMapIndex(key, mapElem)
	}

	receiverValue.Set(newMap)
	return nil
}
