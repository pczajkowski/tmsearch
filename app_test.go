package main

import (
	"strings"
	"testing"
)

func TestSetBaseURL(t *testing.T) {
	urlWithSlash := "http://test.com:880/test/"
	urlWithoutSlash := "http://test.com:880/test"

	t.Log("Testing if mandatory slash is added to the end of BaseURL.")
	var app Application
	app.SetBaseURL(urlWithSlash)
	if strings.HasSuffix(app.BaseURL, "//") {
		t.Errorf("URL has been malformed: %v", app.BaseURL)
	} else {
		t.Log("URL with slash was set correctly!")
	}

	app.SetBaseURL(urlWithoutSlash)
	if !strings.HasSuffix(app.BaseURL, "/") {
		t.Errorf("URL has been malformed: %v", app.BaseURL)
	} else {
		t.Log("URL without slash was set correctly!")
	}
}

func TestLoadLanguages(t *testing.T) {
	var app Application
	app.LoadLanguages()

	testLanguageCode := "dan"
	testLanguage := "Danish"

	t.Log("Testing if languages have been successfully loaded to app.Languages dictionary.")
	_, ok := app.Languages[testLanguageCode]
	if !ok {
		t.Fatalf("There's no key '%v'!", testLanguageCode)
	} else if app.Languages[testLanguageCode] == testLanguage {
		t.Log("Languages are in dictionary!")
	} else {
		t.Fatalf("Value of key '%v' isn't '%v'", testLanguageCode, testLanguage)
	}
}
