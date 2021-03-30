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
		log.Fatalf("Log file error: %s", err)
	}

	writer := csv.NewWriter(logOutput)
	return writer, logOutput
}

func writeLog(info SearchInfo) {
	writer, file := getWriter()

	writer.Write(info.ToArray())
	writer.Flush()
	file.Close()
}
