package main

import "fmt"

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contact   contactInfo
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

	jim.contact = contactInfo{
		email:   "jim@jones.com",
		zipCode: 90011,
	}

	fmt.Println(jim)

	jim.print()

	// point to address of the jim (Reference)
	jimPointer := &jim
	jimPointer.updateName("Jimmy")
	jim.print()

}

func (p person) print() {
	fmt.Printf("%+v\n", p)
}

// point to value of the person pointer (Dereference)
func (pointerToPerson *person) updateName(newFirstName string) {
	(*pointerToPerson).firstName = newFirstName
}
