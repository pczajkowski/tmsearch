package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func postQuery(requestURL string, jsonBytes []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("Error creating post request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error posting query: %s", err)
	}

	return resp, nil
}

func getQuery(destination string) (*http.Response, error) {
	resp, err := http.Get(destination)
	if err != nil {
		return nil, fmt.Errorf("Error getting query: %s", err)
	}

	return resp, nil
}
