package main

import (
	"strings"
	"testing"
)

func TestSetBaseURL(t *testing.T) {
	toTest := []string{"http://test.com:880/test/", "http://test.com:880/test"}

	var app Application

	for _, testCase := range toTest {
		app.setBaseURL(testCase)
		if strings.HasSuffix(app.BaseURL, "//") || !strings.HasSuffix(app.BaseURL, "/") {
			t.Errorf("URL has been malformed: %v", app.BaseURL)
		}
	}
}

func TestLoadLanguages(t *testing.T) {
	var app Application
	app.loadLanguages()

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
