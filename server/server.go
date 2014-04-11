package server

import (
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
	m.Get("/keyword/:word", keywordHandler)
	m.Get("/movie/:title", movieHandler)

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
	tmpl.Execute(res, &keywordsInfo{Keywords: randomKeywords, Previous: []string{},
		Header: getHeader(true)})
}

type keywordsInfo struct {
	Keywords []model.Keyword `json:"keywords"`
	Previous []string        `json:"previous"`
	Header   string          `json:"header"`
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
		res.Write([]byte(fmt.Sprintf("looks like you're watching %v", moviesMatching[0].Title)))
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
		newKeywords = append(newKeywords, *kw)
	}

	tmpl, err := template.ParseFiles(filepath.Join(frontEndRoot, "choose_keyword.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, &keywordsInfo{Keywords: newKeywords, Previous: allMatches,
		Header: getHeader(false)})
}

func movieHandler(params martini.Params) string {
	title := params["title"]
	movie, err := model.FindMovie(title)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return movie.Title
}
