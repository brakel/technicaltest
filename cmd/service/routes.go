package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// Create a new servemux (router)
	mux := http.NewServeMux()

	// Register handler functions
	// 404 everything not specified
	mux.HandleFunc("/", http.NotFound)
	mux.HandleFunc("/articles", app.createArticle)
	mux.HandleFunc("/articles/", app.getArticle)
	mux.HandleFunc("/tags/", app.getArticleByTagAndDate)

	return mux
}
