package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSegmentCleanup(t *testing.T) {
	sourceSegment := "<seg>This is test</seg>"
	targetSegment := "<seg>This is test for target</seg>"

	cleanedSourceSegment := "This is test"
	cleanedTargetSegment := "This is test for target"

	segment := Segment{sourceSegment, targetSegment, "test.txt"}
	segment.clean()

	if segment.Source != cleanedSourceSegment || segment.Target != cleanedTargetSegment {
		t.Fatalf("Segments still have tags!\nSource: %s\nTarget: %s", segment.Source, segment.Target)
	}
}

func serveSearchResults() *httptest.Server {
	searchResults, err := ioutil.ReadFile("./testFiles/searchResults.json")
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(searchResults))
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestSearch(t *testing.T) {
	server := serveSearchResults()
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tmsJSON, err := os.Open("./testFiles/tms.json")
	if err != nil {
		t.Fatalf("Error reading tms: %s", err)
	}
	defer tmsJSON.Close()

	var tms []TM
	err = jsonDecoder(tmsJSON, &tms)
	if err != nil {
		t.Fatalf("Error decoding tms: %s", err)
	}

	info := SearchInfo{Phrase: "something"}
	searchResults := app.search(tms, info)
	if searchResults.TotalResults != 4 {
		t.Fatalf("Not all results returned! (%d)", searchResults.TotalResults)
	}

	testSourceSegments := []string{"<bpt i='1' type='bold'>{}</bpt>Something Test/ Whatever<ept i='1'>{}</ept>",
		"<bpt i='1' type='bold'>{}</bpt>Another Test/ Anything<ept i='1'>{}</ept>"}

	for index, segment := range searchResults.Results[0].Segments {
		if segment.Source != testSourceSegments[index] {
			t.Fatalf("Something is wrong with returned segment!\nShould be:\n%s\nBut is:\n%s",
				testSourceSegments[index], segment)
		}
	}
}

func TestSearchWrongStatus(t *testing.T) {
	server := fakeServer(http.StatusBadRequest, "")
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tmsJSON, err := os.Open("./testFiles/tms.json")
	if err != nil {
		t.Fatalf("Error reading tms: %s", err)
	}
	defer tmsJSON.Close()

	var tms []TM
	err = jsonDecoder(tmsJSON, &tms)
	if err != nil {
		t.Fatalf("Error decoding tms: %s", err)
	}

	info := SearchInfo{Phrase: "something"}
	searchResults := app.search(tms, info)
	if searchResults.TotalResults != 0 {
		t.Fatal("There should be no results!")
	}
}
