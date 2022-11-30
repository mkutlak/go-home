package main

import (
	"fmt"
	"reflect"
)

func main() {
	var t1 int = 123
	fmt.Printf("plusOne: %d (type: %v)\n", plusOne(t1), reflect.TypeOf(t1))

	var t2 float64 = 123.2
	fmt.Printf("plusOne: %f (type: %v)\n", plusOne(t2), reflect.TypeOf(t2))

	fmt.Printf("plusOne: %v (type: %v)\n", sum(t1, t1), reflect.TypeOf(sum(t1, t1)))
	fmt.Printf("plusOne: %v (type: %v)\n", sum(t2, t2), reflect.TypeOf(sum(t2, t2)))
}

func plusOne[V int | float64 | int64 | float32 | int32](t V) V {
	return t + 1
}

func sum[V int | float64 | int64 | float32 | int32](t1, t2 V) V {
	return t1 + t2
}
