package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"html/template"
	"labix.org/v2/mgo"
	"math/rand"
	"net/http"
	"path/filepath"
	"shelman/sourgrapes/model"
	"shelman/sourgrapes/util"
	"strings"
)

var (
	frontEndRoot = "/Users/sam/code/skunk/src/shelman/sourgrapes/frontend"
)

func Start() {

	m := martini.Classic()

	m.Get("/", indexHandler)
	m.Get("/choose/:word", chooseHandler)
	m.Get("/movie/:title", movieHandler)
	m.Get("/keyword/:word", keywordHandler)
	m.Get("/search/movie/", movieSearchHandler)
	m.Get("/search_results/movie/:search", searchResultsHandler)

	m.Use(martini.Static(filepath.Join(frontEndRoot, "js")))
	m.Use(martini.Static(filepath.Join(frontEndRoot, "css")))

	m.Run()
}

func indexHandler(res http.ResponseWriter, req *http.Request) {

	randomKeywords, err := model.GetRandomKeywords(15)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
	}

	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "choose_keyword.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, &chooseInfo{Keywords: randomKeywords, Previous: []string{},
		Header: getHeader(true)})
}

type chooseInfo struct {
	Keywords []model.Keyword `json:"keywords"`
	Previous []string        `json:"previous"`
	Header   string          `json:"header"`
}

func chooseHandler(params martini.Params, res http.ResponseWriter, req *http.Request) {
	word := params["word"]
	keyword, err := model.FindKeyword(word)
	if err != nil {
		if err == mgo.ErrNotFound {
			res.Write([]byte(fmt.Sprintf("stop entering keywords into the url manually")))
			return
		}
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	// the previous words
	prev := req.FormValue("previous")
	allPrev := strings.Split(prev, ",")
	if prev == "" {
		allPrev = []string{}
	}

	// find all movies matching the keyword
	movies, err := model.FindMovies(keyword.Movies)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	// filter out to movies also matching the previous keywords
	allMatches := append(allPrev, word)
	moviesMatching := []model.Movie{}
	for _, movie := range movies {
		matches := true
		for _, w := range allMatches {
			if !util.SliceHasString(movie.Keywords, w) {
				matches = false
				break
			}
		}
		if matches {
			moviesMatching = append(moviesMatching, movie)
		}
	}

	// we narrowed it down!
	if len(moviesMatching) == 1 {
		tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "final_choice.html"))
		if err != nil {
			res.Write([]byte(fmt.Sprintf("error: %v", err)))
			return
		}
		tmpl.Execute(res, moviesMatching[0])
		return
	}

	// we have to narrow it down farther.  get one keyword from each movie matching
	// all of them so far
	newKeywords := []model.Keyword{}
	for _, movie := range moviesMatching {
		n := rand.Int() % len(movie.Keywords)
		kw, err := model.FindKeyword(movie.Keywords[n])
		if err != nil {
			res.Write([]byte(fmt.Sprintf("error: %v", err)))
			return
		}
		kw.Movies = []string{}
		// don't repeat a word.  this could possibly, in a tiny circumstance,
		// cause newKeywords to be empty
		if !util.SliceHasString(allMatches, kw.Word) {
			newKeywords = append(newKeywords, *kw)
		}
	}

	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "choose_keyword.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, &chooseInfo{Keywords: newKeywords, Previous: allMatches,
		Header: getHeader(false)})
}

func movieHandler(params martini.Params, res http.ResponseWriter, req *http.Request) {
	title := strings.ToLower(params["title"])
	movie, err := model.FindMovie(title)
	if err != nil {
		if err == mgo.ErrNotFound {
			res.Write([]byte(fmt.Sprintf("stop entering movie titles into the url manually")))
			return
		}
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "movie.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, movie)
}

func keywordHandler(params martini.Params, res http.ResponseWriter, req *http.Request) {
	word := params["word"]
	keyword, err := model.FindKeyword(word)
	if err != nil {
		if err == mgo.ErrNotFound {
			res.Write([]byte(fmt.Sprintf("stop entering keywords into the url manually")))
			return
		}
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "keyword.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, keyword)
}

func movieSearchHandler(res http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "movie_search.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, nil)
}

type searchResults struct {
	Results []string `json:"results"`
}

func searchResultsHandler(params martini.Params, res http.ResponseWriter, req *http.Request) {
	searchVal := params["search"]

	moviesWithPrefix, err := model.FindMoviesWithPrefix(searchVal)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	results := &searchResults{}
	for _, movie := range moviesWithPrefix {
		results.Results = append(results.Results, movie.Title)
	}

	marshalled, err := json.Marshal(results)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	res.Write(marshalled)
}
