package main

import "fmt"

type person struct {
	firstName string
	lastName  string
}

func main() {

	alex := person{
		firstName: "Alex",
		lastName:  "Anderson",
	}

	fmt.Printf("%s %s\n", alex.firstName, alex.lastName)

	// Declare variable jim of type person
	var jim person

	fmt.Printf("%v %v\n", jim.firstName, jim.lastName)

	jim.firstName = "Jim"
	jim.lastName = "Jones"

	fmt.Printf("%v %v\n", jim.firstName, jim.lastName)

}
