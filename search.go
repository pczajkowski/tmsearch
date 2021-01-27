package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Segment stores source and translated texts.
type Segment struct {
	Source, Target string
}

// Clean cleans <seg> tags from source and translated texts.
func (s *Segment) Clean() {
	re := regexp.MustCompile("</?seg>")
	s.Source = re.ReplaceAllString(s.Source, "")
	s.Target = re.ReplaceAllString(s.Target, "")
}

// CleanedResults stores processed results from given TM.
type CleanedResults struct {
	TMName   string
	Segments []Segment
}

// SearchResults stores processed results from all TMs.
type SearchResults struct {
	SearchPhrase string
	Results      []CleanedResults
	TotalResults int
}

// ResultsFromServer stores results as received from server.
type ResultsFromServer struct {
	ConcResult []struct {
		ConcordanceTextRanges []struct {
			Length, Start int
		}
		ConcordanceTranslationRanges []string
		Length, StartPos             int
		TMEntry                      struct {
			SourceSegment, TargetSegment string
		}
	}
	ConcTransResult, Errors []string
	TotalConcResult         int
}

func postQuery(requestURL string, searchJSON []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(searchJSON))
	if err != nil {
		return nil, fmt.Errorf("Error creating post request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error posting query: %v", err)
	}

	return resp, nil
}

func getCleanedResults(tempResults ResultsFromServer, TMFriendlyName string) CleanedResults {
	var tmResults CleanedResults
	//64 is maximum returned by server
	var numberOfSegments int
	if tempResults.TotalConcResult > 64 {
		numberOfSegments = 64
	} else {
		numberOfSegments = tempResults.TotalConcResult
	}
	//Allocating Segments array beforehand
	tmResults.Segments = make([]Segment, 0, numberOfSegments)
	tmResults.TMName = TMFriendlyName

	for _, result := range tempResults.ConcResult {
		segment := Segment{result.TMEntry.SourceSegment, result.TMEntry.TargetSegment}
		segment.Clean()
		tmResults.Segments = append(tmResults.Segments, segment)
	}
	return tmResults
}

type searchQuery struct {
	SearchExpression []string
}

func getSearchJSON(text string) []byte {
	query := searchQuery{}
	query.SearchExpression = append(query.SearchExpression, text)

	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Printf("Error marshalling query: %v", err)
		return []byte{}
	}

	return queryJSON
}

// Search for given phrase in given TMs.
func (app *Application) Search(TMs []TM, text string) SearchResults {
	searchJSON := getSearchJSON(text)

	tmURL := app.BaseURL + "tms/"

	var finalResults SearchResults
	finalResults.SearchPhrase = text

	for _, tm := range TMs {
		getTM := tmURL + tm.TMGuid
		concordanceURL := getTM + "/concordance"
		requestURL := concordanceURL + app.AuthString

		resp, err := postQuery(requestURL, searchJSON)
		if err != nil {
			log.Println(err)
			return finalResults
		}
		defer resp.Body.Close()

		if resp.StatusCode == 401 {
			time.Sleep(app.Delay)
			app.Login()
			return app.Search(TMs, text)
		}

		var tempResults ResultsFromServer
		JSONDecoder(resp.Body, &tempResults)

		if tempResults.TotalConcResult > 0 {
			tmResults := getCleanedResults(tempResults, tm.FriendlyName)
			finalResults.Results = append(finalResults.Results, tmResults)
			finalResults.TotalResults += len(tmResults.Segments)
		}
	}

	return finalResults
}
