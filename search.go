package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	maxReturnedBySever := 64

	numberOfSegments := tempResults.TotalConcResult
	if numberOfSegments > maxReturnedBySever {
		numberOfSegments = maxReturnedBySever
	}

	tmResults.Segments = make([]Segment, numberOfSegments)
	tmResults.TMName = TMFriendlyName

	for index := 0; index < numberOfSegments; index++ {
		result := tempResults.ConcResult[index]
		segment := Segment{result.TMEntry.SourceSegment, result.TMEntry.TargetSegment}
		segment.clean()
		tmResults.Segments[index] = segment
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
		log.Printf("Error marshalling query: %s", err)
		return []byte{}
	}

	return queryJSON
}

func (app Application) getResultsFromTM(tmURL string, tm TM, searchJSON []byte) (retry bool, result ResultsFromServer) {
	getTM := tmURL + tm.TMGuid
	concordanceURL := getTM + "/concordance"
	requestURL := concordanceURL + app.AuthString

	var tempResults ResultsFromServer
	resp, err := postQuery(requestURL, searchJSON)
	if err != nil {
		log.Println(err)
		return false, tempResults
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		time.Sleep(app.Delay)

		status, err := app.login()
		if !status || err != nil {
			log.Printf("Couldn't log in: %s", err)
			return false, tempResults
		}

		return true, tempResults
	}

	err = jsonDecoder(resp.Body, &tempResults)
	if err != nil {
		log.Printf("Error decoding results: %s", err)
	}

	return false, tempResults
}

func (app Application) search(TMs []TM, text string) SearchResults {
	var finalResults SearchResults
	finalResults.SearchPhrase = text

	searchJSON := getSearchJSON(text)
	if len(searchJSON) == 0 {
		return finalResults
	}

	tmURL := app.BaseURL + "tms/"
	for _, tm := range TMs {
		retry, tempResults := app.getResultsFromTM(tmURL, tm, searchJSON)
		if retry {
			_, tempResults = app.getResultsFromTM(tmURL, tm, searchJSON)
		}

		if tempResults.TotalConcResult <= 0 {
			continue
		}

		tmResults := getCleanedResults(tempResults, tm.FriendlyName)
		finalResults.Results = append(finalResults.Results, tmResults)
		finalResults.TotalResults += len(tmResults.Segments)
	}

	return finalResults
}
