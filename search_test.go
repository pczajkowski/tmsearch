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

	segment := Segment{sourceSegment, targetSegment}
	segment.Clean()

	t.Log("Testing if <seg> tags will be removed from segments.")
	if segment.Source == cleanedSourceSegment && segment.Target == cleanedTargetSegment {
		t.Log("Segments have been cleaned!")
	} else {
		t.Fatalf("Segments still have tags!\nSource: %v\nTarget: %v", segment.Source, segment.Target)
	}
}

func ServeSearchResults() *httptest.Server {
	searchResults, err := ioutil.ReadFile("./testFiles/searchResults.json")
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(searchResults))
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestSearch(t *testing.T) {
	server := ServeSearchResults()
	defer server.Close()

	var app Application
	app.setBaseURL(server.URL)

	tmsJSON, err := os.Open("./testFiles/tms.json")
	if err != nil {
		t.Fatalf("error reading tms: %v", err)
		return
	}
	defer tmsJSON.Close()

	var tms []TM
	JSONDecoder(tmsJSON, &tms)

	testSourceSegment1 := "<bpt i='1' type='bold'>{}</bpt>Something Test/ Whatever<ept i='1'>{}</ept>"
	testSourceSegment2 := "<bpt i='1' type='bold'>{}</bpt>Another Test/ Anything<ept i='1'>{}</ept>"

	t.Log("Testing search method.")
	searchResults := app.Search(tms, "something")
	if searchResults.TotalResults != 4 {
		t.Fatalf("Not all results returned! (%v)", searchResults.TotalResults)
	}

	segment := searchResults.Results[0].Segments[0]
	if segment.Source == testSourceSegment1 || segment.Source == testSourceSegment2 {
		t.Log("Search results fine!")
	} else {
		t.Fatalf("Something is wrong with returned segment!\n%v", segment)
	}
}
