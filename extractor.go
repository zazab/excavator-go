package extractor

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type convertError struct {
	from reflect.Type
	to   reflect.Type
}

func newConvertError(receiver, data reflect.Value) error {
	return &convertError{
		from: unrollType(data.Type()),
		to:   unrollType(receiver.Type()),
	}

}

func (err *convertError) Error() string {
	return fmt.Sprintf(
		"can't convert '%s' to '%s'",
		err.from, err.to,
	)
}

func Extract(receiver, data interface{}) error {
	receiverType := reflect.TypeOf(receiver)
	if receiverType.Kind() != reflect.Ptr {
		log.Printf("receiver: %#+v", receiver)
		log.Printf("data: %#+v", data)
		log.Printf("receiver kind: %#+v", reflect.TypeOf(receiver).Kind())
		log.Printf("data kind: %#+v", reflect.TypeOf(data).Kind())
		return errors.New("non ptr receiver")
	}

	var (
		receiverValue = unrollValue(reflect.ValueOf(receiver))
		dataValue     = unrollValue(reflect.ValueOf(data))
	)

	return extract(receiverValue, unrollValue(dataValue))
}

func extract(receiver, data reflect.Value) error {
	var (
		receiverKind = receiver.Type().Kind()
		dataKind     = data.Type().Kind()
	)

	switch receiverKind {
	case reflect.Struct:
		log.Println("exporting struct")
		return extractStruct(receiver, data)
	case reflect.Slice:
		log.Println("exporting slice")
		return extractSlice(receiver, data)
	case reflect.Map:
		log.Println("exporting map")
		return extractMap(receiver, data)
	default:
		log.Printf("exporting plain type %s", receiver.Type())
		if receiverKind != dataKind {
			return newConvertError(receiver, data)
		}

		receiver.Set(data)
		return nil
	}
}
