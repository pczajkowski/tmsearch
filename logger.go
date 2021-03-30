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

func getWriter() (*csv.Writer, *os.File) {
	logFile := filepath.Join("log", (time.Now().Format(dateFormat) + ".log"))
	logOutput, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Printf("Log file error: %s", err)
		return nil, nil
	}

	writer := csv.NewWriter(logOutput)
	return writer, logOutput
}

func writeLog(info SearchInfo) {
	writer, file := getWriter()
	if writer == nil || file == nil {
		return
	}

	writer.Write(info.ToArray())
	if err := writer.Error(); err != nil {
		log.Printf("Error writing csv: %s", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Printf("Error flushing csv: %s", err)
	}

	file.Close()
	if err := file.Close(); err != nil {
		log.Printf("Error closing csv: %s", err)
	}
}
