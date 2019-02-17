package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	viewLinkT     = "https://github.com/{{.Repo}}/blob/{{.Branch}}/{{.Path}}"
	rawLinkT      = "https://raw.githubusercontent.com/{{.Repo}}/{{.Branch}}/{{.Dir}}"
	editLinkT     = "https://github.com/{{.Repo}}/edit/{{.Branch}}/{{.Path}}"
	feedbackLinkT = "https://github.com/{{.Repo}}/issues/new"
	dirAPILinkT   = "https://api.github.com/repos/{{.Repo}}/contents/{{.Dir}}?ref={{.Branch}}"
)

// Config is the TOML config with attached handlers
type Config struct {
	Categories map[string]Category
}

func (config *Config) handleCategories(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(config.Categories))
	for k := range config.Categories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// ...
	w.WriteHeader(http.StatusOK)
}

func (config *Config) handleCategory(w http.ResponseWriter, r *http.Request) {
	categorySlug := mux.Vars(r)["category"]
	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	articles, err := category.getArticles()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("templates/base.html", "templates/category.html")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, struct {
		Articles []Article
		Category Category
	}{
		Articles: articles,
		Category: category,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (config *Config) handleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categorySlug := vars["category"]
	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	article := Article{
		Category: category,
		File:     vars["article"],
	}
	raw, err := article.getRaw()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	article.Raw = raw
	article.Title = article.getTitle()
	article.HTML = string(blackfriday.Run([]byte(raw)))

	t, err := template.ParseFiles("templates/base.html", "templates/article.html")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, article)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	var conf Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", conf.handleCategories)
	r.HandleFunc("/{category}/", conf.handleCategory)
	r.HandleFunc("/{category}/{article}", conf.handleArticle)

	http.Handle("/", r)
	fmt.Println("Ready")
	log.Fatal(http.ListenAndServe("127.0.0.1:1337", nil))
}
