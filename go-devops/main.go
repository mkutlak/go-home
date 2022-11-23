package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Not enough arguments!\n")
		os.Exit(1)
	}

	fmt.Printf("Hello World\nArguments: %v\n", strings.Join(args[1:], " "))
}
