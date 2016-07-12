package excavator

import (
	"reflect"

	"github.com/seletskiy/hierr"
)

func excavateMap(receiverValue, dataValue reflect.Value) error {
	dataKind := unrollType(dataValue.Type()).Kind()

	switch dataKind {
	case reflect.Map:
		return excavateMapFromMap(receiverValue, dataValue)
	case reflect.Struct:
		return excavateMapFromStruct(receiverValue, dataValue)
	default:
		return newConvertError(receiverValue, dataValue)
	}
}

func excavateMapFromStruct(receiverValue, dataValue reflect.Value) error {
	var (
		mapType     = receiverValue.Type()
		mapElemType = mapType.Elem()
		newMap      = reflect.MakeMap(mapType)

		dataType     = dataValue.Type()
		fieldsNumber = dataType.NumField()
	)

	for index := 0; index < fieldsNumber; index++ {
		field := dataType.Field(index)
		key := getFieldKey(field)

		dataFieldValue := unrollValue(dataValue.Field(index))
		mapElem := reflect.New(mapElemType).Elem()

		err := excavate(mapElem, dataFieldValue)
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

func excavateMapFromMap(receiverValue, dataValue reflect.Value) error {
	mapType := receiverValue.Type()
	mapElemType := mapType.Elem()

	newMap := reflect.MakeMap(mapType)
	for _, key := range dataValue.MapKeys() {
		dataMapElem := unrollValue(dataValue.MapIndex(key))
		mapElem := reflect.New(mapElemType).Elem()

		err := excavate(mapElem, dataMapElem)
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
