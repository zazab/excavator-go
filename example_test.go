package excavator_test

import (
	"fmt"
	"log"

	"github.com/zazab/excavator-go"
)

type S struct {
	Field1 int      `excavator:"a"`
	Field2 []string `excavator:"b"`
}

type S2 struct {
	Field1 S `excavator:"nom"`
	Field2 S `excavator:"bon"`
}

func ExampleTest() {
	var (
		b    = []string{"foo", "bar", "baz"}
		data = map[string]interface{}{
			"nom": map[string]interface{}{
				"a": 10,
				"b": &b,
			},
			"bon": map[string]interface{}{
				"a": 13,
				"b": nil,
			},
		}

		receiver S2
	)

	err := excavator.Excavate(&receiver, data)
	if err != nil {
		log.Fatalf("can't export: %s", err)
	}

	fmt.Printf("result: %v", receiver)
	// Output:
	// result: {{10, ["foo", "bar", "baz"]}, {13}}
}
