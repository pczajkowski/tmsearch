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

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./html/index.html"))
	t.Execute(w, app.Languages)
}

//Addition for counter
func add(x, y int) int {
	return x + y
}

func DisplaySearchResults(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("lang")
	searchPhrase := r.URL.Query().Get("phrase")

	if searchPhrase != "" {
		var searchResults SearchResults
		if language == "" || app.CheckLanguage(language) {
			searchResults = app.Search(app.GetTMs(language), searchPhrase)
			Logger(r, searchResults.TotalResults)
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

func DisplayTMs(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("lang")

	var TMList []TM
	if language == "" || app.CheckLanguage(language) {
		TMList = app.GetTMs(language)
		Logger(r, len(TMList))
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
	http.HandleFunc("/", ServeIndex)
	http.HandleFunc("/q", DisplaySearchResults)
	http.HandleFunc("/tms", DisplayTMs)
	log.Fatal(http.ListenAndServe(hostname, nil))
}
