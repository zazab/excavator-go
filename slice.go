package excavator

import (
	"fmt"
	"log"
	"reflect"
)

func excavateSlice(receiverValue, dataValue reflect.Value, tag string) error {
	log.Println("processing slice")
	dataKind := unrollType(dataValue.Type()).Kind()
	if dataKind != reflect.Slice {
		return newConvertError(receiverValue, dataValue)
	}
	log.Println("kind", dataKind)

	sliceType := receiverValue.Type()
	newSlice := reflect.MakeSlice(
		sliceType, dataValue.Len(), dataValue.Cap(),
	)

	if dataValue.CanInterface() && dataValue.IsNil() {
		receiverValue.Set(newSlice)
		return nil
	}

	for index := 0; index < dataValue.Len(); index++ {
		log.Println("processing item", index)
		sliceElem := unrollValue(dataValue.Index(index))

		newSliceElem := newSlice.Index(index)
		err := excavate(newSliceElem, sliceElem, tag)
		if err != nil {
			return fmt.Errorf("can't set element #%d: %s", index, err)
		}
	}

	receiverValue.Set(newSlice)
	return nil
}
