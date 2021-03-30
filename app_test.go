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

	fromMap, ok := app.Languages[testLanguageCode]
	if !ok {
		t.Fatalf("There's no key '%v'!", testLanguageCode)
	}

	if fromMap != testLanguage {
		t.Fatalf("Value of key '%v' isn't '%v'", testLanguageCode, testLanguage)
	}
}
