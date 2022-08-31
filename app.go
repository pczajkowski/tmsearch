package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Application stores main information needed to run the app
type Application struct {
	Name, AccessToken, Sid, BaseURL, AuthString string
	Languages                                   map[string]string
	Delay                                       time.Duration
}

func (app *Application) setBaseURL(baseURL string) {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	app.BaseURL = baseURL
}

func jsonDecoder(data io.ReadCloser, target interface{}) error {
	decoder := json.NewDecoder(data)
	return decoder.Decode(target)
}

func (app *Application) loadLanguages() bool {
	data, err := os.Open("./html/languages.json")
	if err != nil {
		log.Printf("Error reading languages: %s", err)
		return false
	}

	defer func() {
		if err := data.Close(); err != nil {
			log.Printf("Error closing file: %s", err)
		}
	}()

	app.Languages = make(map[string]string)
	err = jsonDecoder(data, &app.Languages)
	if err != nil {
		log.Printf("Error decoding languages: %s", err)
		return false
	}

	return true
}

func (app *Application) checkLanguage(language *string) bool {
	_, ok := app.Languages[*language]
	return ok
}

func (app *Application) login() (bool, error) {
	credentials, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		return false, fmt.Errorf("Error reading credentials: %s", err)
	}

	loginURL := app.BaseURL + "auth/login"

	resp, err := postQuery(loginURL, credentials)
	if err != nil {
		return false, fmt.Errorf("Error logging in: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Error logging in: %s", resp.Status)
	}

	err = jsonDecoder(resp.Body, &app)
	if err != nil {
		return false, fmt.Errorf("Error decoding login details: %s", err)
	}

	app.AuthString = "?authToken=" + app.AccessToken
	log.Println(app.AuthString, resp.Status)
	return true, nil
}
