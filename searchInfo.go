package main

import (
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04"
)

// SearchInfo represents concise information about search query
type SearchInfo struct {
	Date time.Time
	Host, Phrase, Language, LanguageCode string
	ResultsServed int
}

func (s *SearchInfo) ToArray() []string {
	return []string{s.Date.Format(timeFormat), s.Host, s.Phrase, s.Language, strconv.Itoa(s.ResultsServed)}
}

func (s *SearchInfo) GetInfoFromRequest(r *http.Request) {
	s.Date = time.Now()
	s.Host, _, _ = net.SplitHostPort(r.RemoteAddr)

	s.Phrase = r.URL.Query().Get("phrase")
	if s.Phrase == "" {
		s.Phrase = "TMS"
	}

	s.LanguageCode = r.URL.Query().Get("lang")
	s.Language = app.Languages[s.LanguageCode]
}
