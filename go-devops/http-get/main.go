package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Todos struct {
	//Id        int    `json:"id"`
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {

	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: http-get <url>\n")
		os.Exit(1)
	}

	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("Not a URL: %v\n", args[1])
		os.Exit(1)
	}

	resp, err := http.Get(args[1])
	if err != nil {
		log.Fatalf("GET request has failed: %v\n", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Invalid HTTP Status Code: %d\n%s\n", resp.StatusCode, body)
	}

	var todos []Todos

	err = json.Unmarshal(body, &todos)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range todos {
		//fmt.Printf("JSON Parsed:\nUserId: %d\nId: %d\nTitle: %s\nCompleted: %t\n", t.UserId, t.Id, t.Title, t.Completed)
		fmt.Printf("JSON Parsed:\nUserId: %d\nTitle: %s\nCompleted: %t\n", t.UserId, t.Title, t.Completed)
	}

}
