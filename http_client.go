package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiBaseURL = "https://api.ravelry.com"

func MakeGETRequest(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	login := os.Getenv("RAVELRY_LOGIN")
	password := os.Getenv("RAVELRY_PASSWORD")
	req.SetBasicAuth(login, password)

	req.Header.Set("User-Agent", "curl/7.68.0")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}
