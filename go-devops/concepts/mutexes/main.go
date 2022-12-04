package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type myType struct {
	counter int
	mu      sync.Mutex
}

func (m *myType) IncreaseCounter() {
	m.mu.Lock()
	m.counter++
	m.mu.Unlock()
}

func main() {

	myTypeInstance := &myType{}

	finished := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(myTypeInstance *myType) {
			myTypeInstance.mu.Lock()
			fmt.Printf("IN Counter: %d\n", myTypeInstance.counter)

			myTypeInstance.counter++
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))

			fmt.Printf("PAST Counter: %d\n", myTypeInstance.counter)

			if myTypeInstance.counter == 5 {
				fmt.Printf("Found Counter: %d == 5\n", myTypeInstance.counter)
			}

			finished <- true
			myTypeInstance.mu.Unlock()
		}(myTypeInstance)
	}
	for i := 0; i < 10; i++ {
		<-finished
	}

	fmt.Printf("OUT Counter: %d\n", myTypeInstance.counter)
}
