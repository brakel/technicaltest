package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// This blog post was a big help
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
// Only one json object can be received at a time
func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	// Return 405 if not a POST request
	if r.Method != http.MethodPost {
		// Set the 'Allow' header to let the user know what methods are accepted
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Return 415 if the header is not correctly set
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not set to application/json"
		app.errorLog.Println(msg)
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var a article
	// Initialise the decoder and return an error if there are any unexpected fields
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	// Decode the request body into the article struct
	err := dec.Decode(&a)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Check the ID is a valid number
	id, err := strconv.Atoi(a.ID)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Invalid ID - Not a valid number", http.StatusBadRequest)
		return
	}

	// Check if local storage already contains an articles with the current ID
	if _, ok := app.articles[id]; ok {
		app.errorLog.Println("Article with that ID already exists")
		http.Error(w, "Article with that ID already exists", http.StatusBadRequest)
		return
	}

	// Store the article in a map with the id as the key
	app.articles[id] = a

	app.infoLog.Printf("%s request to %s", r.Method, r.URL)

}

// Using query strings seems more intuitive but I will stick to the requirements - /articles/{id}
// You could also use a third party library (gorilla, httprouter) to handle URL variables
func (app *application) getArticle(w http.ResponseWriter, r *http.Request) {
	// Return 405 if not a GET request
	if r.Method != http.MethodGet {
		// Set the 'Allow' header to let the user know what methods are accepted
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Trim the URL prefix
	trimmedURL := strings.TrimPrefix(r.URL.Path, "/articles/")
	// Check if id a valid number
	id, err := strconv.Atoi(trimmedURL)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Invalid ID - Not a valid number", http.StatusBadRequest)
		return
	}
	// Does the article exist?
	article, ok := app.articles[id]
	if !ok {
		app.infoLog.Println("Article not found")
		w.Write([]byte("Article not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Initialise the encoder to stream to the responsewriter
	enc := json.NewEncoder(w)
	// Pretty print
	enc.SetIndent("", " ")
	enc.Encode(article)

	app.infoLog.Printf("%s request to %s", r.Method, r.URL)
}

func (app *application) getArticleByTagAndDate(w http.ResponseWriter, r *http.Request) {
	// Return 405 if not a GET request
	if r.Method != http.MethodGet {
		// Set the 'Allow' header to let the user know what methods are accepted
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Trim the URL prefix and any trailing slashes
	trimmedURL := strings.TrimPrefix(r.URL.Path, "/tags/")
	trimmedURL = strings.TrimSuffix(trimmedURL, "/")
	// Split the URL params and return if we don't have exactly two
	params := strings.Split(trimmedURL, "/")
	if len(params) != 2 {
		app.errorLog.Println("Invalid URL parameters")
		http.Error(w, "Invalid URL parameters", http.StatusBadRequest)
		return
	}

	// Check if the date is of valid length (8)
	if len(params[1]) != 8 {
		app.errorLog.Println("Invalid date format - use yyyymmdd")
		http.Error(w, "Invalid date format - use yyyymmdd", http.StatusBadRequest)
		return
	}
	// Insert dashes into date to match the format in article
	date := params[1][:4] + "-" + params[1][4:6] + "-" + params[1][6:]
	tag := strings.ToLower(params[0])

	result := result{Tag: tag, Count: 0, Articles: make([]string, 0), RelatedTags: make([]string, 0)}
	// Map to act like a set for related tags
	relatedTags := make(map[string]bool)

	// Loop through the articles stored
	for id, article := range app.articles {
		// Check if the dates match
		if article.Date == date {
			// Loop through the matched articles tags
			for _, t := range article.Tags {
				// If there is a tag that matches
				if strings.ToLower(t) == tag {
					// Add the tag then loop again to add all other tags
					// We only want the last 10 articles of the day
					if len(result.Articles) < 10 {
						result.Articles = append(result.Articles, strconv.Itoa(id))
					}
					for _, rt := range article.Tags {
						if strings.ToLower(rt) != tag {
							relatedTags[rt] = true
						}
					}
					result.Count += 1
				}
			}
		}
	}
	// Finally add the related tags stored in the map
	for rt := range relatedTags {
		result.RelatedTags = append(result.RelatedTags, rt)
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(result)

	app.infoLog.Printf("%s request to %s", r.Method, r.URL)
}
