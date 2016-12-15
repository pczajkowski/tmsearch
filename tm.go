package main

import (
	"log"
	"net/http"
	"time"
)

type TMList struct {
	TMs []struct {
		NumEntries, AccessLevel                                                                         int
		Client, Domain, FriendlyName, Project, SourceLangCode, Subject, TMGuid, TMOwner, TargetLangCode string
	}
}

func GetQuery(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error getting query: %v", err)
	}

	return resp
}

func GetTMs(language string) TMList {
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
		return GetTMs(language)
	}

	var results TMList
	JsonDecoder(resp.Body, &results.TMs)

	return results
}
