package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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

func (app *Application) setBaseURL(url string) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	app.BaseURL = url
}

func jsonDecoder(data io.ReadCloser, target interface{}) {
	decoder := json.NewDecoder(data)

	err := decoder.Decode(target)
	if err != nil {
		log.Printf("Error reading json: %s", err)
	}
}

func (app *Application) loadLanguages() {
	data, err := os.Open("./html/languages.json")
	if err != nil {
		log.Fatalf("Error reading languages: %s", err)
	}
	defer data.Close()

	app.Languages = make(map[string]string)
	jsonDecoder(data, &app.Languages)
}

func (app Application) checkLanguage(language string) bool {
	_, ok := app.Languages[language]
	return ok
}

func (app *Application) login() {
	credentials, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Fatalf("Error reading credentials: %s", err)
	}

	loginURL := app.BaseURL + "auth/login"

	resp, err := postQuery(loginURL, credentials)
	if err != nil {
		log.Fatalf("Error logging in: %s", err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Error logging in: %s", resp.Status)
	}
	defer resp.Body.Close()

	jsonDecoder(resp.Body, &app)

	app.AuthString = "?authToken=" + app.AccessToken
	log.Println(app.AuthString, resp.Status)
}
