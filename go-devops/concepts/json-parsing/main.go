package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type MyJson struct {
	Test  any    `json:"test"`
	Test3 string `json:"test3"`
}

func main() {
	var jsonParsed MyJson

	err := json.Unmarshal([]byte(`{"test": {"test2": [1,2,3]}}`), &jsonParsed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", reflect.TypeOf(jsonParsed))

	switch v := jsonParsed.Test.(type) {

	case map[string]any:
		fmt.Printf("Map found: %v\n", v)
		field1, ok := v["test2"]
		if ok {
			switch v2 := field1.(type) {
			case []any:
				fmt.Printf("I Found Slice any?\n")
				for _, v2Element := range v2 {
					fmt.Printf("Type any: %v\n", reflect.TypeOf(v2Element))
				}
			default:
				fmt.Printf("Type not found %v\n", reflect.TypeOf(v2))
			}
		}

	default:
		fmt.Printf("Type not found %v\n", reflect.TypeOf(v))

	}
}
