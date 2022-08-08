package main

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

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
