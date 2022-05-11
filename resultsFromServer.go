package main

import (
	"time"
)

// ResultsFromServer stores results as received from server.
type ResultsFromServer struct {
	ConcResult []struct {
		ConcordanceTextRanges []struct {
			Length int `json:"Length"`
			Start  int `json:"Start"`
		} `json:"ConcordanceTextRanges"`
		ConcordanceTranslationRanges []struct {
			Length int     `json:"Length"`
			Score  float64 `json:"Score"`
			Start  int     `json:"Start"`
		} `json:"ConcordanceTranslationRanges"`
		Length   int `json:"Length"`
		StartPos int `json:"StartPos"`
		TMEntry  struct {
			SourceSegment string    `json:"SourceSegment"`
			TargetSegment string    `json:"TargetSegment"`
			EntryID       int       `json:"EntryId"`
			LastModified  time.Time `json:"LastModified"`
			Modifier      string    `json:"Modifier"`
			DocumentName  string    `json:"DocumentName"`
		} `json:"TMEntry"`
	} `json:"ConcResult"`
	ConcTransResult []struct {
		Expression string  `json:"Expression"`
		Score      float64 `json:"Score"`
	} `json:"ConcTransResult"`
	Errors []struct {
		ErrorType int    `json:"ErrorType"`
		QueryPart string `json:"QueryPart"`
	} `json:"Errors"`
	TotalConcResult int `json:"TotalConcResult"`
}