package excavator

import (
	"fmt"
	"log"
	"reflect"
)

func excavateSlice(receiverValue, dataValue reflect.Value, tag string) error {
	dataKind := unrollType(dataValue.Type()).Kind()
	if dataKind != reflect.Slice {
		log.Println("e1")
		return newConvertError(receiverValue, dataValue)
	}

	sliceType := receiverValue.Type()
	newSlice := reflect.MakeSlice(
		sliceType, dataValue.Len(), dataValue.Cap(),
	)

	for index := 0; index < dataValue.Len(); index++ {
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
