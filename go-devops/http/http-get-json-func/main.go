package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("Words: %s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	out := []string{}

	for word, occurrence := range o.Words {
		out = append(out, fmt.Sprintf("%s (%d)", word, occurrence))
	}

	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: http-get <url>\n")
		os.Exit(1)
	}

	resp, err := doRequest(args[1])
	if err != nil {
		if requestErr, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", requestErr.Err, requestErr.HTTPCode, requestErr.Body)
			os.Exit(1)
		}
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	if resp == nil {
		fmt.Printf("No response!\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", resp.GetResponse())
}

func doRequest(requestUrl string) (Response, error) {

	if _, err := url.ParseRequestURI(requestUrl); err != nil {
		return nil, fmt.Errorf("URL is not valid URL: %s", err)
	}

	resp, err := http.Get(requestUrl)

	if err != nil {
		return nil, fmt.Errorf("Error: Http Get error: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("Error: ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid HTTP Status Code: %d\n%s\n", resp.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("No valid JSON returned!"),
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)

	if err != nil {
		return nil, RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("Page unmarshal error: %s", err),
		}

	}

	switch page.Name {

	case "words":
		var words Words

		err := json.Unmarshal(body, &words)

		if err != nil {
			return nil, RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Words unmarshal error: %s", err),
			}
		}

		return words, nil

	case "occurrence":
		var occurrence Occurrence

		err := json.Unmarshal(body, &occurrence)

		if err != nil {
			return nil, RequestError{
				HTTPCode: resp.StatusCode,
				Body:     string(body),
				Err:      fmt.Sprintf("Occurence unmarshal error: %s", err),
			}
		}

		return occurrence, nil
	}

	return nil, nil

}
