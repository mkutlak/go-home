package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client ClientIface, requestURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)

	if err != nil {
		return "", fmt.Errorf("Marshal Error: %s", err)
	}

	resp, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("Error: Http Post error: %s", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("Error: ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid HTTP Status Code: %d\n%s\n", resp.StatusCode, string(respBody))
	}

	if !json.Valid(respBody) {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(respBody),
			Err:      fmt.Sprintf("No valid JSON returned!"),
		}
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(respBody, &loginResponse)

	if err != nil {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(respBody),
			Err:      fmt.Sprintf("LoginResponse unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(respBody),
			Err:      "Empty token reply!",
		}
	}

	return loginResponse.Token, nil
}
