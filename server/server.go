package server

import (
	"fmt"
	"github.com/go-martini/martini"
	"html/template"
	"net/http"
	"path/filepath"
	"shelman/sourgrapes/model"
)

var (
	tmplRoot = "/Users/sam/code/skunk/src/shelman/sourgrapes/frontend"
)

func Start() {

	m := martini.Classic()

	m.Get("/", indexHandler)
	m.Get("/keyword/:word", keywordHandler)
	m.Get("/movie/:title", movieHandler)

	m.Run()
}

func indexHandler(res http.ResponseWriter, req *http.Request) {

	randomKeywords, err := model.GetRandomKeywords(15)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
	}

	tmpl, err := template.ParseFiles(filepath.Join(tmplRoot, "main.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, randomKeywords)
}

func keywordHandler(params martini.Params, res http.ResponseWriter, req *http.Request) {
	word := params["word"]
	keyword, err := model.FindKeyword(word)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join(tmplRoot, "keyword.html"))
	if err != nil {
		res.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
	tmpl.Execute(res, keyword)
}

func movieHandler(params martini.Params) string {
	title := params["title"]
	movie, err := model.FindMovie(title)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return movie.Title
}
