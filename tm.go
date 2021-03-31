package main

import (
	"log"
	"net/http"
	"time"
)

// TM stores information about TM.
type TM struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
}

func getQuery(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting query: %s", err)
	}

	return resp
}

func (app *Application) getTMs(language string) []TM {
	tmURL := app.BaseURL + "tms/"
	queryURL := tmURL + app.AuthString
	if language != "" {
		queryURL += "&targetLang=" + language
	}

	resp := getQuery(queryURL)
	defer resp.Body.Close()

	var results []TM
	if resp.StatusCode == http.StatusBadRequest {
		time.Sleep(app.Delay)

		status, err := app.login()
		if !status || err != nil {
			log.Printf("Couldn't log in: %s", err)
			return results
		}

		return app.getTMs(language)
	}

	err := jsonDecoder(resp.Body, &results)
	if err != nil {
		log.Printf("Error decoding results: %s", err)
	}

	return results
}
