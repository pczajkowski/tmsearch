package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	Name, AccessToken, Sid, BaseURL, AuthString string
	Languages                                   map[string]string
	Delay                                       time.Duration
}

func JsonDecoder(data io.ReadCloser, target interface{}) {
	decoder := json.NewDecoder(data)

	err := decoder.Decode(target)
	if err != nil {
		log.Printf("error reading json: %v", err)
	}
}

func (app *Application) LoadLanguages() {
	data, err := os.Open("./html/languages.json")
	defer data.Close()
	if err != nil {
		log.Printf("error reading languages: %v", err)
		return
	}

	app.Languages = make(map[string]string)
	JsonDecoder(data, &app.Languages)
}

func (app Application) CheckLanguage(language string) bool {
	_, ok := app.Languages[language]
	if !ok {
		return false
	}

	return true
}

func (app *Application) Login() {
	credentials, err := ioutil.ReadFile("./secrets.json")
	if err != nil {
		log.Printf("Error reading credentials: %v", err)
	}

	loginURL := app.BaseURL + "auth/login"

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(credentials))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error logging in: %v", err)
	}
	defer resp.Body.Close()

	JsonDecoder(resp.Body, &app)

	app.AuthString = "?authToken=" + app.AccessToken
	log.Println(app.AuthString, resp.Status)
}
