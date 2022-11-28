package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type MockClient struct {
	ResponseOutput *http.Response
	PostResponse   *http.Response
}

func (mc *MockClient) Get(url string) (resp *http.Response, err error) {
	return mc.ResponseOutput, nil
}

func (m MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	return m.PostResponse, nil
}

func TestDoRequest(t *testing.T) {
	words := WordsPage{
		Page: Page{"words"},
		Words: Words{
			Input: "abc",
			Words: []string{"a", "b"},
		},
	}

	wordsBytes, err := json.Marshal(words)

	if err != nil {
		t.Errorf("Marshal Error: %s", err)
	}

	apiInstance := api{
		Options: Options{},
		Client: &MockClient{
			ResponseOutput: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
			},
		},
	}

	response, err := apiInstance.DoGetRequest("http://localost:8080/words")

	if err != nil {
		t.Errorf("DoGetRequest Error: %s", err)
	}

	if response == nil {
		t.Errorf("Response is empty")
	}

	if response.GetResponse() != `Words: a, b` {
		t.Errorf("Unexpected Response: %s", response.GetResponse())
	}

}
