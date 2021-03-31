package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func serveTMs() *httptest.Server {
	tms, err := ioutil.ReadFile("./testFiles/tms.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(tms))
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestGetTMs(t *testing.T) {
	server := serveTMs()
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tms := app.getTMs("")
	if len(tms) != 2 {
		t.Fatalf("Not all TMs read! (%d)", len(tms))
	}

	if tms[0].FriendlyName != "Test TM 1" || tms[1].FriendlyName != "Test TM 2" {
		t.Fatalf("Something went wrong while reading TMs!\n%v", tms)
	}
}
