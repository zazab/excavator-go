package extractor

import (
	"fmt"
	"log"
	"reflect"
)

func extractSlice(receiverValue, dataValue reflect.Value) error {
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
		err := extract(newSliceElem, sliceElem)
		if err != nil {
			return fmt.Errorf("can't set element #%d: %s", index, err)
		}
	}

	receiverValue.Set(newSlice)
	return nil
}
