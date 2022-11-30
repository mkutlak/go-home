package main

import (
	"fmt"
	"reflect"
)

func main() {

	var t1 string = "this is a string!"
	discoverType(t1)

	var t2 *string = &t1
	discoverType(*t2)

	var t3 int = 123
	discoverType(t3)

	discoverType(nil)

}

func discoverType(t any) {
	switch v := t.(type) {
	case string:
		fmt.Printf("String found: %s\n", v)

	case *string:
		fmt.Printf("Pointer string found: %s\n", *v)

	case int:
		fmt.Printf("Integer found: %d\n", v)

	default:
		fmt.Printf("Type not found! (%v)\n", reflect.TypeOf(t))
	}
}
