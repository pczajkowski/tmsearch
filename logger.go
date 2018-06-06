package main

import (
	"encoding/csv"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"strconv"
)

const (
	timeFormat = "2006-01-02 15:04"
	dateFormat = "20060102"
)

type searchInfo struct {
	Date time.Time
	Host, Phrase, Language string
	ResultsServed int
}

func (s *searchInfo) ToArray() []string {
	return []string{s.Date.Format(timeFormat), s.Host, s.Phrase, s.Language, strconv.Itoa(s.ResultsServed)}
}

func getInfoFromRequest(r *http.Request) searchInfo {
	info := searchInfo{Date: time.Now()}
	info.Host, _, _ = net.SplitHostPort(r.RemoteAddr)
	info.Phrase = r.URL.Query().Get("phrase")

	language := r.URL.Query().Get("lang")
	if language == "" {
		info.Language = "All languages"
	} else {
		info.Language = app.Languages[language]
	}

	return info
}

func getWriter() *csv.Writer {
	logFile := filepath.Join("log", (time.Now().Format(dateFormat) + ".log"))
	logOutput, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}

	writer := csv.NewWriter(logOutput)
	return writer
}

// WriteLog main function, saves event to the log.
func WriteLog(r *http.Request, resultsServed int) {
	info := getInfoFromRequest(r)
	info.ResultsServed = resultsServed

	writer := getWriter()

	writer.Write(info.ToArray())
	writer.Flush()
}
