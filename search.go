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
	Source, Target, DocumentName string
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
	Results      []*CleanedResults
	TotalResults int
}

func getCleanedResults(tempResults *ResultsFromServer, TMFriendlyName string) CleanedResults {
	var tmResults CleanedResults
	var numberOfSegments = len(tempResults.ConcResult)

	tmResults.Segments = make([]Segment, numberOfSegments)
	tmResults.TMName = TMFriendlyName

	for index := 0; index < numberOfSegments; index++ {
		result := tempResults.ConcResult[index]
		segment := Segment{result.TMEntry.SourceSegment, result.TMEntry.TargetSegment, result.TMEntry.DocumentName}
		segment.clean()
		tmResults.Segments[index] = segment
	}

	return tmResults
}

func getSearchJSON(info *SearchInfo) []byte {
	query := searchQuery{}
	query.SearchExpression = append(query.SearchExpression, info.Phrase)
	query.Options.CaseSensitive = false
	query.Options.ReverseLookup = info.Reverse
	query.Options.ResultsLimit = info.SearchLimit

	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Printf("Error marshalling query: %s", err)
		return []byte{}
	}

	return queryJSON
}

func (app *Application) getResultsFromTM(tmURL string, tm *TM, searchJSON []byte) (retry bool, result ResultsFromServer) {
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("Problem getting results (%s)!", resp.Status)
		return false, tempResults
	}

	err = jsonDecoder(resp.Body, &tempResults)
	if err != nil {
		log.Printf("Error decoding results: %s", err)
	}

	return false, tempResults
}

func (app *Application) search(tms []TM, info *SearchInfo) SearchResults {
	var finalResults SearchResults
	finalResults.SearchPhrase = info.Phrase

	searchJSON := getSearchJSON(info)
	if len(searchJSON) == 0 {
		return finalResults
	}

	tmURL := app.BaseURL + "tms/"
	max := len(tms)
	for i := 0; i < max; i++ {
		retry, tempResults := app.getResultsFromTM(tmURL, &tms[i], searchJSON)
		if retry {
			_, tempResults = app.getResultsFromTM(tmURL, &tms[i], searchJSON)
		}

		if tempResults.TotalConcResult <= 0 {
			continue
		}

		tmResults := getCleanedResults(&tempResults, tms[i].FriendlyName)
		finalResults.Results = append(finalResults.Results, &tmResults)
		finalResults.TotalResults += len(tmResults.Segments)
	}

	return finalResults
}
