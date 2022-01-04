package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Application-wide dependencies
// I am using a map as the "database" in this application for simplicity's sake
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	articles map[int]article
}

// https://mholt.github.io/json-to-go/ - handy
type article struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

// I found the description of 'count' in the design brief confusing, the example gives a value of 17 for two articles. For the purpose of this application it will equal the number of results found. The other alternative was the total number of tags of every article matched.
type result struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}

func main() {
	// Command line flag to set the network address - defaults to *:4000
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Custom loggers to standard streams
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	articles := make(map[int]article)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		articles: articles,
	}

	// Initialise a new server struct with our config
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	errorLog.Fatal(srv.ListenAndServe())
}
