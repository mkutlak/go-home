package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan string)

	fmt.Println("t1")
	go testFunction(c)
	fmt.Println("t3")

	areWeFinished := <-c
	fmt.Println(areWeFinished)
}

func testFunction(c chan string) {
	for i := 0; i < 5; i++ {
		fmt.Println("checking...")
		time.Sleep(1 * time.Second)
	}

	c <- "we are fnshd"
}
