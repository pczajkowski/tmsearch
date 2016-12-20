package main

import (
	"log"
	"net/http"
	"time"
)

type TM struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
}

func GetQuery(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error getting query: %v", err)
	}

	return resp
}

func (app *Application) GetTMs(language string) []TM {
	tmURL := app.BaseURL + "tms/"
	var queryURL string
	if language == "" {
		queryURL = tmURL + app.AuthString
	} else {
		queryURL = tmURL + app.AuthString + "&targetLang=" + language
	}

	resp := GetQuery(queryURL)
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		time.Sleep(app.Delay)
		app.Login()
		return app.GetTMs(language)
	}

	var results []TM
	JsonDecoder(resp.Body, &results)

	return results
}
