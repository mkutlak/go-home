package main

import "fmt"

func main() {

	a := make(map[string]string)
	a["test"] = "value"

	testPointer(a)

	fmt.Printf("A: %v\n", a)
}

func testPointer(a map[string]string) {
	a["test2"] = "another string"
}
