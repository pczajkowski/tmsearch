package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
	"strings"
)

// TM stores information about TM.
type TM struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
	LastModified string
}

func (t *TM) LastModifiedDate() *time.Time {
	if !strings.HasSuffix(t.LastModified, "Z") {
		t.LastModified += ".000Z"
	}

	modified, err := time.Parse(time.RFC3339, t.LastModified)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &modified
}

func (app *Application) getTMs(language string) []TM {
	tmURL := app.BaseURL + "tms?"

	params := url.Values{}
	params.Add("authToken", app.AccessToken)
	params.Add("targetLang", language)

	tmURL += params.Encode()

	var results []TM
	resp, err := getQuery(tmURL)
	if err != nil {
		log.Println(err)
		return results
	}

	defer resp.Body.Close()

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

	err = jsonDecoder(resp.Body, &results)
	if err != nil {
		log.Printf("Error decoding TM results: %s", err)
	}

	return results
}
