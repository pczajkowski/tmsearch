package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"
)

var host = flag.String("h", "localhost", "host")
var port = flag.String("p", "80", "port")
var url = flag.String("b", "", "API URL")
var app Application
var errorPage = template.Must(template.ParseFiles("./html/error.html"))

// ServeIndex serves index page.
func serveIndex(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./html/index.html"))
	t.Execute(w, app.Languages)
}

//Addition for counter
func add(x, y int) int {
	return x + y
}

// DisplaySearchResults displays search results as HTML page.
func displaySearchResults(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.GetInfoFromRequest(r)

	if info.Phrase != "" {
		var searchResults SearchResults
		if info.LanguageCode == "" || app.CheckLanguage(info.LanguageCode) {
			searchResults = app.Search(app.GetTMs(info.LanguageCode), info.Phrase)
			info.ResultsServed = searchResults.TotalResults
			WriteLog(info)
		} else {
			errorPage.Execute(w, "Language not valid!")
			return
		}

		if len(searchResults.Results) > 0 {
			funcs := template.FuncMap{"add": add}
			t := template.Must(template.New("results.html").Funcs(funcs).ParseFiles("./html/results.html"))
			t.Execute(w, searchResults)
		} else {
			errorPage.Execute(w, "Nothing found!")
		}
	} else {
		errorPage.Execute(w, "You need to enter search phrase!")
	}
}

// DisplayTMs displays TMs as HTML page.
func displayTMs(w http.ResponseWriter, r *http.Request) {
	var info SearchInfo
	info.GetInfoFromRequest(r)

	var TMList []TM
	if info.LanguageCode == "" || app.CheckLanguage(info.LanguageCode) {
		TMList = app.GetTMs(info.LanguageCode)
		info.ResultsServed = len(TMList)
		WriteLog(info)
	} else {
		errorPage.Execute(w, "Language not valid!")
		return
	}

	if len(TMList) > 0 {
		t := template.Must(template.New("tms.html").ParseFiles("./html/tms.html"))
		t.Execute(w, TMList)
	} else {
		errorPage.Execute(w, "No TMs to display!")
	}
}

func main() {
	flag.Parse()
	app.SetBaseURL(*url)
	if app.BaseURL == "" {
		log.Panicln("Can't do anything without URL to API")
	}

	app.Login()
	app.LoadLanguages()
	app.Delay = time.Duration(20 * time.Second)

	hostname := *host + ":" + *port
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/q", displaySearchResults)
	http.HandleFunc("/tms", displayTMs)
	log.Fatal(http.ListenAndServe(hostname, nil))
}
