package excavator_test

import "log"

type S struct {
	Field1 int      `zhash:"a"`
	Field2 []string `zhash:"b"`
}

type S2 struct {
	Field1 S `zhash:"nom"`
	Field2 S `zhash:"bon"`
}

func main() {
	var (
		b    = []string{"ol", "gon", "don"}
		data = map[string]interface{}{
			"nom": map[string]interface{}{
				"a": 10,
				"b": &b,
			},
			"bon": map[string]interface{}{
				"a": 13,
				"b": &b,
			},
		}

		receiver S2
	)

	err := excavator.Excavate(&receiver, data)
	if err != nil {
		log.Fatalf("can't export: %s", err)
	}

	log.Printf("result %#+v", receiver)
}
