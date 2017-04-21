package main

import (
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

// GetLogger returns new logger
func GetLogger() *log.Logger {
	logFile := filepath.Join("log", (time.Now().Format("200612") + ".log"))
	logOutput, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
	}

	logger := log.New(logOutput, "", log.Ldate|log.Ltime)
	return logger
}

// Logger main function, saves event to the log.
func Logger(r *http.Request, resultsServed int) {
	host, searchPhrase, language := GetInfoFromRequest(r)
	logger := GetLogger()

	if searchPhrase != "" {
		logger.Printf(",%v,\"%v\",\"%v\",%v\n", host, searchPhrase, language, resultsServed)
	} else {
		logger.Printf(",%v,TMS,\"%v\",%v\n", host, language, resultsServed)
	}
}
