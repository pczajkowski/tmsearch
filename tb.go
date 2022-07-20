package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

// TB stores information about TM.
type TB struct {
	NumEntries, AccessLevel                                                                         int
	Client, Domain, FriendlyName, Project, Subject, TBGuid, TBOwner string
	Languages []string
	LastModified time.Time
}

func (app *Application) getTBs(language string) []TB {
	tbURL := app.BaseURL + "tbs?"

	params := url.Values{}
	params.Add("authToken", app.AccessToken)

	if language != "" {
		params.Add("lang", language)
	}

	tbURL += params.Encode()

	var results []TB
	resp, err := getQuery(tbURL)
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

		return app.getTBs(language)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Problem getting TBs (%s)!", resp.Status)
		return results
	}

	err = jsonDecoder(resp.Body, &results)
	if err != nil {
		log.Printf("Error decoding TB results: %s", err)
	}

	return results
}
