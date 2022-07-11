package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

// TM stores information about TM.
type TM struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
}

func getQuery(destination string) *http.Response {
	resp, err := http.Get(destination)
	if err != nil {
		log.Printf("Error getting query: %s", err)
	}

	return resp
}

func (app *Application) getTMs(language string) []TM {
	tmURL := app.BaseURL + "tms?"

	params := url.Values{}
	params.Add("authToken", app.AccessToken)
	params.Add("targetLang", language)

	tmURL += params.Encode()

	resp := getQuery(tmURL)
	defer resp.Body.Close()

	var results []TM
	if resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(app.Delay)

		status, err := app.login()
		if !status || err != nil {
			log.Printf("Couldn't log in: %s", err)
			return results
		}

		return app.getTMs(language)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Problem getting TMs (%s)!", resp.Status)
		return results
	}

	err := jsonDecoder(resp.Body, &results)
	if err != nil {
		log.Printf("Error decoding TM results: %s", err)
	}

	return results
}
