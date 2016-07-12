package excavator

import (
	"reflect"

	"github.com/seletskiy/hierr"
)

func excavateMap(receiverValue, dataValue reflect.Value) error {
	dataKind := unrollType(dataValue.Type()).Kind()
	if dataKind != reflect.Map {
		return newConvertError(receiverValue, dataValue)
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
