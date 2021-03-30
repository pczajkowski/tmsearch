package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fakeServer(statusCode int, data string) *httptest.Server {
	function := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprint(w, data)
	}

	return httptest.NewServer(http.HandlerFunc(function))
}

func TestSetBaseURL(t *testing.T) {
	toTest := []string{"http://test.com:880/test/", "http://test.com:880/test"}

	var app Application

	for _, testCase := range toTest {
		app.setBaseURL(testCase)
		if strings.HasSuffix(app.BaseURL, "//") || !strings.HasSuffix(app.BaseURL, "/") {
			t.Errorf("URL has been malformed: %s", app.BaseURL)
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
		t.Fatalf("There's no key '%s'!", testLanguageCode)
	}

	if fromMap != testLanguage {
		t.Fatalf("Value of key '%s' isn't '%s'", testLanguageCode, testLanguage)
	}
}

func TestLogin(t *testing.T) {
	loginResponse := `{
    "Name": "admin",
    "Sid": "00000000-0000-0000-0001-000000000001",
    "AccessToken": "fde0f7ed-d585-48ec-a0a9-397aea195ccd"
}`

	server := fakeServer(http.StatusOK, loginResponse)
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)
	status, err := app.login()
	if !status || err != nil {
		t.Fatalf("Status: %v, error: %s", status, err)
	}
}

func TestLoginBadURL(t *testing.T) {
	var app Application
	app.setBaseURL("badURL")

	status, err := app.login()
	if status || err == nil {
		t.Fatalf("Status: %v, error: %s", status, err)
	}
}

func TestLoginWrongStatus(t *testing.T) {
	server := fakeServer(http.StatusBadRequest, "")
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)
	status, err := app.login()
	if status || err == nil {
		t.Fatalf("Status: %v, error: %s", status, err)
	}
}
