package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/spf13/pflag"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	rawLinkT    = "https://raw.githubusercontent.com/{{.Repo}}/{{.Branch}}/{{.Dir}}"
	dirAPILinkT = "https://api.github.com/repos/{{.Repo}}/contents/{{.Dir}}?ref={{.Branch}}"
)

// Config is the TOML config with attached handlers
type Config struct {
	Listen     string
	Root       string
	Categories map[string]Category
}

func (config *Config) handleCategories(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(config.Categories))
	for k := range config.Categories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var categories []Category
	var category Category
	for _, key := range keys {
		category = config.Categories[key]
		category.Slug = key
		categories = append(categories, category)
	}

	t, err := template.ParseFiles("templates/base.html", "templates/categories.html")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, categories)
	if err != nil {
		log.Fatal(err)
		return
	}
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
	category.Slug = categorySlug
	file := vars["article"] + category.Ext
	article := Article{
		Category: category,
		File:     file,
		Slug:     vars["article"],
	}
	raw, err := article.getRaw()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	article.Raw = raw
	article.Title = article.getTitle()
	article.HTML = string(blackfriday.Run([]byte(article.Raw)))

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
	path := pflag.String("config", "config.toml", "path to config file")
	listen := pflag.String("listen", "", "server and port to listen (value from config by default)")
	pflag.Parse()

	var conf Config
	if _, err := toml.DecodeFile(*path, &conf); err != nil {
		log.Fatal(err)
		return
	}
	if *listen != "" {
		conf.Listen = *listen
	}

	r := mux.NewRouter()
	r.HandleFunc("/", conf.handleCategories)
	r.HandleFunc("/{category}/", conf.handleCategory)
	r.HandleFunc("/{category}/{article}/", conf.handleArticle)

	http.Handle(conf.Root, r)
	fmt.Println("Ready")
	log.Fatal(http.ListenAndServe(conf.Listen, nil))
}
