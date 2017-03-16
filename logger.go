package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// GetInfoFromRequest reads information from given request and returns them as tuple.
func GetInfoFromRequest(r *http.Request) (string, string, string) {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	searchPhrase := r.URL.Query().Get("phrase")
	language := r.URL.Query().Get("lang")
	if language == "" {
		language = "All languages"
	} else {
		language = app.Languages[language]
	}

	return host, searchPhrase, language
}

// WriteLog writes log entry to the file.
func WriteLog(logString string) error {
	logFile := filepath.Join("log", (time.Now().Format("200612") + ".log"))
	logOutput, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer logOutput.Close()

	_, err = logOutput.WriteString(logString)
	return err
}

// Logger main function, saves event to the log.
func Logger(r *http.Request, resultsServed int) {
	host, searchPhrase, language := GetInfoFromRequest(r)
	timeFormat := "2006-01-02 15:04:05"

	var logString string
	if searchPhrase != "" {
		logString = fmt.Sprintf("%v,%v,\"%v\",\"%v\",%v\n", time.Now().Format(timeFormat), host, searchPhrase, language, resultsServed)
	} else {
		logString = fmt.Sprintf("%v,%v,TMS,\"%v\",%v\n", time.Now().Format(timeFormat), host, language, resultsServed)
	}

	err := WriteLog(logString)
	if err != nil {
		log.Fatalf("Error writing log: %v", err)
	}
}
