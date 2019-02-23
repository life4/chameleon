package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/pflag"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

const (
	commitsLinkT = "https://api.github.com/repos/{{.Repo}}/commits?sha={{.Branch}}&path={{.Dir}}"
	rawLinkT     = "https://raw.githubusercontent.com/{{.Repo}}/{{.Branch}}/{{.Dir}}"
	dirAPILinkT  = "https://api.github.com/repos/{{.Repo}}/contents/{{.Dir}}?ref={{.Branch}}"
)

// Config is the TOML config with attached handlers
type Config struct {
	Listen     string
	Root       string
	Templates  string
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

	t, err := template.ParseFiles(
		path.Join(config.Templates, "base.html"),
		path.Join(config.Templates, "categories.html"),
	)
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

	t, err := template.ParseFiles(
		path.Join(config.Templates, "base.html"),
		path.Join(config.Templates, "category.html"),
	)
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

func (config *Config) handleArticleRedirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categorySlug := vars["category"]
	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, strings.TrimSuffix(vars["article"], category.Ext)+"/", http.StatusPermanentRedirect)
}

func (config *Config) handleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categorySlug := vars["category"]
	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if strings.HasSuffix(vars["article"], category.Ext) {
		newURL := fmt.Sprintf("/%s/%s/", categorySlug, strings.TrimSuffix(vars["article"], category.Ext))
		http.Redirect(w, r, newURL, http.StatusPermanentRedirect)
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

	authors, err := article.getAuthors()
	if err != nil {
		fmt.Println(err)
	} else {
		article.Authors = authors
	}

	t, err := template.ParseFiles(
		path.Join(config.Templates, "base.html"),
		path.Join(config.Templates, "article.html"),
	)
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
	configPath := pflag.StringP("config", "c", "config.toml", "path to config file")
	templatesPath := pflag.StringP("templates", "t", "templates/", "path to templates directory")
	listen := pflag.StringP("listen", "l", "", "server and port to listen (value from config by default)")
	pflag.Parse()

	var conf Config
	if _, err := toml.DecodeFile(*configPath, &conf); err != nil {
		log.Fatal(err)
		return
	}
	if *listen != "" {
		conf.Listen = *listen
	}
	conf.Templates = *templatesPath

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(20),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/", conf.handleCategories)
	r.HandleFunc("/{category}/", conf.handleCategory)
	r.HandleFunc("/{category}/{article}", conf.handleArticleRedirect)
	r.HandleFunc("/{category}/{article}/", conf.handleArticle)

	http.Handle(conf.Root, cacheClient.Middleware(r))
	fmt.Printf("Ready to get connections on %s%s\n", conf.Listen, conf.Root)
	log.Fatal(http.ListenAndServe(conf.Listen, nil))
}
