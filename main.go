package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/tidwall/gjson"
)

const (
	viewLinkT     = "https://github.com/{{.Repo}}/blob/{{.Branch}}/{{.Path}}"
	rawLinkT      = "https://raw.githubusercontent.com/{{.Repo}}/{{.Branch}}/{{.Path}}"
	editLinkT     = "https://github.com/{{.Repo}}/edit/{{.Branch}}/{{.Path}}"
	feedbackLinkT = "https://github.com/{{.Repo}}/issues/new"
	dirAPILinkT   = "https://api.github.com/repos/{{.Repo}}/contents/{{.Dir}}?ref={{.Branch}}"
)

// Category is struct with all information about content in given category
type Category struct {
	Repo, Branch, Dir, Name, Ext string
}

func (category *Category) makeLink(linkTemplate string) (content string, err error) {
	t, err := template.New("linkTemplate").Parse(linkTemplate)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, *category)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

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

	link, err := category.makeLink(dirAPILinkT)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	namesQuery := gjson.Get(string(content), "#.name")
	var names []string
	for _, name := range namesQuery.Array() {
		names = append(names, name.String())
	}

	t, err := template.ParseFiles("templates/category.html")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, names)
}

func (config *Config) handleArticle(w http.ResponseWriter, r *http.Request) {
	// ...
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
