package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	links := []string{
		"http://google.com",
		"http://stackoverflow.com",
		"http://facebook.com",
		"http://golang.org",
		"http://amazon.com",
		"http://google.com",
		"http://stackoverflow.com",
		"http://facebook.com",
		"http://golang.org",
		"http://amazon.com",
	}

	ch := make(chan string)

	for _, link := range links {
		// run checkLink in go routine (fork thread)
		go checkLink(link, ch)
	}

	// infinite loop
	for l := range ch {
		// function literal (lambda function)
		go func(link string) {
			time.Sleep(5 * time.Second)
			checkLink(link, ch)
		}(l)
	}

}

func checkLink(link string, c chan string) {
	_, err := http.Get(link)

	if err != nil {
		log.Fatalln(link, "is most likely DOWN.")
		c <- link
		return
	}

	fmt.Println(link, "is UP!")
	c <- link
}
