package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	dateFormat = "20060102"
)

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
func WriteLog(info SearchInfo) {
	writer := getWriter()

	writer.Write(info.ToArray())
	writer.Flush()
}
