package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func ServeTMs() *httptest.Server {
	tms, err := ioutil.ReadFile("./testFiles/tms.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(tms))
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestGetTMs(t *testing.T) {
	server := ServeTMs()
	defer server.Close()

	var app Application
	app.SetBaseURL(server.URL)

	t.Log("Testing if TMs are properly read from the server.")
	tms := app.GetTMs("")
	if len(tms) != 2 {
		t.Fatalf("Not all TMs read! (%v)", len(tms))
	} else if tms[0].FriendlyName == "Test TM 1" || tms[0].FriendlyName == "Test TM 2" {
		t.Log("TMs properly read!")
	} else {
		t.Fatalf("Something went wrong while reading TMs!\n%v", tms)
	}
}
