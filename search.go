package main

import (
	"bytes"
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

// PostQuery sends POST query to server and returns response.
func PostQuery(requestURL string, searchJSON []byte) *http.Response {
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(searchJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error posting query: %v", err)
	}

	return resp
}

// Search searches for given phrase in given TMs.
func (app *Application) Search(TMs []TM, text string) SearchResults {
	searchString := "{ \"SearchExpression\": [ \"" + text + "\" ]}"
	searchJSON := []byte(searchString)

	tmURL := app.BaseURL + "tms/"

	var finalResults SearchResults
	finalResults.SearchPhrase = text

	var results []CleanedResults
	for _, tm := range TMs {
		getTM := tmURL + tm.TMGuid
		concordanceURL := getTM + "/concordance"
		requestURL := concordanceURL + app.AuthString

		resp := PostQuery(requestURL, searchJSON)
		defer resp.Body.Close()
		if resp.StatusCode == 401 {
			time.Sleep(app.Delay)
			app.Login()
			return app.Search(TMs, text)
		}

		var tempResults ResultsFromServer
		JSONDecoder(resp.Body, &tempResults)

		if tempResults.TotalConcResult > 0 {
			var tmResults CleanedResults
			//Allocating Segments array beforehand
			tmResults.Segments = make([]Segment, 0, tempResults.TotalConcResult)
			tmResults.TMName = tm.FriendlyName

			for _, result := range tempResults.ConcResult {
				segment := Segment{result.TMEntry.SourceSegment, result.TMEntry.TargetSegment}
				segment.Clean()
				tmResults.Segments = append(tmResults.Segments, segment)
			}
			results = append(results, tmResults)
			finalResults.TotalResults += len(tmResults.Segments)
		}
	}
	finalResults.Results = results
	return finalResults
}
