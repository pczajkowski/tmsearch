package main

import (
	"encoding/json"
	"log"
	"regexp"
	"time"
)

// Segment stores source and translated texts.
type Segment struct {
	Source, Target string
}

func (s *Segment) clean() {
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
		segment.clean()
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

			status, err := app.login()
			if !status || err != nil {
				log.Printf("Couldn't log in: %s", err)
				return finalResults
			}

			return app.Search(TMs, text)
		}

		var tempResults ResultsFromServer
		jsonDecoder(resp.Body, &tempResults)

		if tempResults.TotalConcResult > 0 {
			tmResults := getCleanedResults(tempResults, tm.FriendlyName)
			finalResults.Results = append(finalResults.Results, tmResults)
			finalResults.TotalResults += len(tmResults.Segments)
		}
	}

	return finalResults
}
