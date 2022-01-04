package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

func TestCreateArticle(t *testing.T) {
	app := &application{
		errorLog: log.New(io.Discard, "", 0), // Discard any logging info
		infoLog:  log.New(io.Discard, "", 0),
		articles: make(map[int]article),
	}

	testArticle := &article{
		ID:    "1",
		Title: "test title",
		Date:  "2022-01-01",
		Body:  "test body",
		Tags:  []string{"test"},
	}

	ta, err := json.Marshal(testArticle)
	if err != nil {
		log.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(ta))
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Content-Type", "application/json")

	app.createArticle(rr, r)
	if _, ok := app.articles[1]; !ok {
		t.Errorf("test article not found")
	}
	if app.articles[1].ID != testArticle.ID {
		t.Errorf("want ID: %s; got %s", testArticle.ID, app.articles[1].ID)
	}
	if app.articles[1].Title != testArticle.Title {
		t.Errorf("want Title: %s; got %s", testArticle.Title, app.articles[1].Title)
	}
	if app.articles[1].Date != testArticle.Date {
		t.Errorf("want Date: %s; got %s", testArticle.Date, app.articles[1].Date)
	}
	if app.articles[1].Body != testArticle.Body {
		t.Errorf("want Body: %s; got %s", testArticle.Body, app.articles[1].Body)
	}
	if !reflect.DeepEqual(app.articles[1].Tags, testArticle.Tags) {
		t.Errorf("want Tags: %s; got %s", testArticle.Tags, app.articles[1].Tags)
	}
}

func TestGetArticle(t *testing.T) {
	var responseArticle article

	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		articles: make(map[int]article),
	}

	testArticle := &article{
		ID:    "55",
		Title: "Bananas are superior to apples",
		Date:  "1970-07-29",
		Body:  "Surely not wrong",
		Tags:  []string{"Banana", "Apple", "Fruit"},
	}

	app.articles[55] = *testArticle

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/articles/55", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.getArticle(rr, r)
	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want status %d; got %d", http.StatusOK, rs.StatusCode)
	}

	err = json.NewDecoder(rs.Body).Decode(&responseArticle)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	// A quick way to check equality but you can't pinpoint errors
	/*
		if !reflect.DeepEqual(responseArticle, *testArticle) {
			t.Error("Not equal")
		}
	*/
	if responseArticle.ID != testArticle.ID {
		t.Errorf("want ID: %s; got %s", testArticle.ID, responseArticle.ID)
	}
	if responseArticle.Title != testArticle.Title {
		t.Errorf("want Title: %s; got %s", testArticle.Title, responseArticle.Title)
	}
	if responseArticle.Date != testArticle.Date {
		t.Errorf("want Date: %s; got %s", testArticle.Date, responseArticle.Date)
	}
	if responseArticle.Body != testArticle.Body {
		t.Errorf("want Body: %s; got %s", testArticle.Body, responseArticle.Body)
	}
	if !reflect.DeepEqual(app.articles[55].Tags, testArticle.Tags) {
		t.Errorf("want Tags: %s; got %s", testArticle.Tags, app.articles[55].Tags)
	}
}

func TestGetArticleByTagAndDate(t *testing.T) {
	var responseResult result
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		articles: make(map[int]article),
	}

	ta1 := &article{
		ID:    "3",
		Title: "The Three Bananas",
		Date:  "1994-11-24",
		Body:  "All for one and one for all, united we stand peeled we fall",
		Tags:  []string{"Banana", "Literature", "Fruit"},
	}
	ta2 := &article{
		ID:    "77",
		Title: "Carrots: the good, the bad and the orange",
		Date:  "1994-11-24",
		Body:  "You see in this world there's two kinds of carrots, my friend - those with Loaded Guns, and those who dig. You dig.",
		Tags:  []string{"Vegetable", "Carrot", "Western"},
	}
	ta3 := &article{
		ID:    "12",
		Title: "Pulp Fiction",
		Date:  "1994-11-24",
		Body:  "Is fruit pulp real?",
		Tags:  []string{"Defintely-not-a-movie", "Fruit", "Juice"},
	}
	ta4 := &article{
		ID:    "5",
		Title: "The Apple Strikes Back",
		Date:  "1980-05-06",
		Body:  "Luke, you can destroy the Apple. He has foreseen this. It is your destiny. Join me, and together we can rule the galaxy as father and son",
		Tags:  []string{"Apple", "Fruit", "Space"},
	}
	app.articles[3] = *ta1
	app.articles[77] = *ta2
	app.articles[12] = *ta3
	app.articles[5] = *ta4

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/tags/fruit/19941124/", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.getArticleByTagAndDate(rr, r)
	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want status %d; got %d", http.StatusOK, rs.StatusCode)
	}

	err = json.NewDecoder(rs.Body).Decode(&responseResult)
	if err != nil {
		t.Fatal(err)
	}

	if responseResult.Tag != "fruit" {
		t.Errorf("want Tag: %s, got %s", "fruit", responseResult.Tag)
	}
	if responseResult.Count != 2 {
		t.Errorf("want Count: %d, got %d", 2, responseResult.Count)
	}

	sort.Strings(responseResult.Articles)
	sort.Strings(responseResult.RelatedTags)

	if !reflect.DeepEqual(responseResult.Articles, []string{"12", "3"}) {
		t.Errorf("want Articles: %s, got %s", []string{"12", "3"}, responseResult.Articles)
	}
	if !reflect.DeepEqual(responseResult.RelatedTags, []string{"Banana", "Defintely-not-a-movie", "Juice", "Literature"}) {
		t.Errorf("want Articles: %s, got %s", []string{"Banana", "Defintely-not-a-movie", "Juice", "Literature"}, responseResult.RelatedTags)
	}
}
