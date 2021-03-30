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

// GetQuery sends GET query and returns response.
func GetQuery(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting query: %v", err)
	}

	return resp
}

// GetTMs returns list of TMs for given target language.
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
	jsonDecoder(resp.Body, &results)

	return results
}
