package excavator

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

// Excavate tries to convert data to receiver. If your struct already have some
// tags (for example for marshaling form some format), then you can pass this
// tag as structTag. If structTag not given, DefaultStructTag is used.
func Excavate(receiver, data interface{}, structTag ...string) error {
	tag := DefaultStructTag
	if len(structTag) > 1 {
		return errors.New("more than one tag field given")
	}
	if len(structTag) == 1 {
		tag = structTag[0]
	}

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

	return excavate(receiverValue, unrollValue(dataValue), tag)
}

func excavate(receiver, data reflect.Value, tag string) error {
	var (
		receiverKind = receiver.Type().Kind()
		dataKind     = data.Type().Kind()
	)

	switch receiverKind {
	case reflect.Struct:
		log.Println("exporting struct")
		return excavateStruct(receiver, data, tag)
	case reflect.Slice:
		log.Println("exporting slice")
		return excavateSlice(receiver, data, tag)
	case reflect.Map:
		log.Println("exporting map")
		return excavateMap(receiver, data, tag)
	default:
		log.Printf("exporting plain type %s", receiver.Type())
		if receiverKind != dataKind {
			return newConvertError(receiver, data)
		}

		receiver.Set(data)
		return nil
	}
}
