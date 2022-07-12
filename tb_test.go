package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveTBs() *httptest.Server {
	tbs, err := ioutil.ReadFile("./testFiles/tbs.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(tbs))
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestGetTBs(t *testing.T) {
	server := serveTBs()
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tbs := app.getTBs("")
	if len(tbs) != 2 {
		t.Fatalf("Not all TBs read! (%d)", len(tbs))
	}

	if tbs[0].FriendlyName != "Test TB 1" || tbs[1].FriendlyName != "Test TB 2" {
		t.Fatalf("Something went wrong while reading TBs!\n%v", tbs)
	}
}

func TestGetTVsWrongStatus(t *testing.T) {
	server := fakeServer(http.StatusBadRequest, "")
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tbs := app.getTBs("")
	if len(tbs) != 0 {
		t.Fatal("There should be no TBs!")
	}
}
