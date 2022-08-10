package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
	"encoding/csv"
	"fmt"
	"io"
)

var host = flag.String("h", "localhost", "host")
var port = flag.String("p", "80", "port")
var baseURL = flag.String("b", "", "API URL")
var app Application
var errorPage = template.Must(template.ParseFiles("./html/error.html"))

func serveIndex(w http.ResponseWriter, _ *http.Request) {
	t := template.Must(template.ParseFiles("./html/index.html"))
	t.Execute(w, app.Languages)
}

func displaySearchResults(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.ParseRequest(r)

	if info.Phrase == "" {
		errorPage.Execute(w, "You need to enter search phrase!")
		return
	}

	if info.LanguageCode != "" && !app.checkLanguage(&info.LanguageCode) {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	tms := app.getTMs(info.LanguageCode)
	if len(tms) == 0 {
		errorPage.Execute(w, "Couldn't get TMs!")
		return
	}

	searchResults := app.search(tms, &info)
	info.ResultsServed = searchResults.TotalResults
	writeLog(info)

	if len(searchResults.Results) == 0 {
		errorPage.Execute(w, "Nothing found!")
		return
	}

	t := template.Must(template.New("results.html").ParseFiles("./html/results.html"))
	t.Execute(w, searchResults)
}

func displayTMs(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.ParseRequest(r)

	if info.LanguageCode != "" && !app.checkLanguage(&info.LanguageCode) {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	TMList := app.getTMs(info.LanguageCode)
	info.ResultsServed = len(TMList)
	writeLog(info)

	if info.ResultsServed == 0 {
		errorPage.Execute(w, "No TMs to display!")
		return
	}

	t := template.Must(template.New("tms.html").ParseFiles("./html/tms.html"))
	t.Execute(w, TMList)
}

func prepareResponseForCSV(w http.ResponseWriter, name string) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", name))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("charset", "utf-8")

	io.WriteString(w, "\xEF\xBB\xBF")
}

func serveTMsAsCSV(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.ParseRequest(r)

	if info.LanguageCode != "" && !app.checkLanguage(&info.LanguageCode) {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	tmList := app.getTMs(info.LanguageCode)
	info.ResultsServed = len(tmList)
	writeLog(info)

	if info.ResultsServed == 0 {
		errorPage.Execute(w, "No TMs to display!")
		return
	}

	prepareResponseForCSV(w, "tms")
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write(tmList[0].Header()); err != nil {
		errorPage.Execute(w, fmt.Sprintf("error writing header to csv: %s", err))
	}

	for i := range tmList {
		if err := csvWriter.Write(tmList[i].ToArray()); err != nil {
			errorPage.Execute(w, fmt.Sprintf("error writing record to csv: %s", err))
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		errorPage.Execute(w, err.Error)
	}
}

func displayTBs(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.ParseRequest(r)

	if info.LanguageCode != "" && !app.checkLanguage(&info.LanguageCode) {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	TBList := app.getTBs(info.LanguageCode)
	info.ResultsServed = len(TBList)
	writeLog(info)

	if info.ResultsServed == 0 {
		errorPage.Execute(w, "No TBs to display!")
		return
	}

	t := template.Must(template.New("tbs.html").ParseFiles("./html/tbs.html"))
	t.Execute(w, TBList)
}

func serveTBsAsCSV(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.ParseRequest(r)

	if info.LanguageCode != "" && !app.checkLanguage(&info.LanguageCode) {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	tbList := app.getTBs(info.LanguageCode)
	info.ResultsServed = len(tbList)
	writeLog(info)

	if info.ResultsServed == 0 {
		errorPage.Execute(w, "No TBs to display!")
		return
	}

	prepareResponseForCSV(w, "tbs")
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write(tbList[0].Header()); err != nil {
		errorPage.Execute(w, fmt.Sprintf("error writing header to csv: %s", err))
	}

	for i := range tbList {
		if err := csvWriter.Write(tbList[i].ToArray()); err != nil {
			errorPage.Execute(w, fmt.Sprintf("error writing record to csv: %s", err))
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		errorPage.Execute(w, err.Error)
	}
}

func main() {
	flag.Parse()
	if *baseURL == "" {
		log.Fatalln("Can't do anything without URL to API")
	}

	app.setBaseURL(*baseURL)

	status, err := app.login()
	if !status || err != nil {
		log.Fatalf("Couldn't log in: %s", err)
	}

	if !app.loadLanguages() {
		log.Fatal("Couldn't load languages!")
	}
	app.Delay = time.Duration(20 * time.Second)

	hostname := *host + ":" + *port
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/q", displaySearchResults)
	http.HandleFunc("/tms", displayTMs)
	http.HandleFunc("/tmscsv", serveTMsAsCSV)
	http.HandleFunc("/tbs", displayTBs)
	http.HandleFunc("/tbscsv", serveTBsAsCSV)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(hostname, nil))
}
