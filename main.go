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

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	"github.com/spf13/pflag"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

const (
	commitsLinkT = "https://api.github.com/repos/{{.Repo}}/commits?sha={{.Branch}}&path={{.Dir}}"
	rawLinkT     = "https://raw.githubusercontent.com/{{.Repo}}/{{.Branch}}/{{.Dir}}"
	dirAPILinkT  = "https://api.github.com/repos/{{.Repo}}/contents/{{.Dir}}?ref={{.Branch}}"
)

// Config is the TOML config with attached handlers
type Config struct {
	Listen    string
	Root      string
	Templates string

	Cache        bool
	Contributors bool
	Lint         bool
	Views        bool

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
	categorySlug := chi.URLParam(r, "category")
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
	categorySlug := chi.URLParam(r, "category")
	articleSlug := strings.TrimPrefix(chi.URLParam(r, "*"), "/"+categorySlug+"/")

	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// if there is no extension then it is probably article without trailing slash (/lol -> /lol/)
	if !strings.ContainsRune(articleSlug, '.') {
		http.Redirect(w, r, articleSlug+"/", http.StatusTemporaryRedirect)
		return
	}

	// if it has extension of article then it is definetly article (/lol.md -> /lol/)
	if strings.HasSuffix(articleSlug, category.Ext) {
		http.Redirect(w, r, strings.TrimSuffix(articleSlug, category.Ext)+"/", http.StatusTemporaryRedirect)
	}

	// if it has different extension then it is file (/lol.jpg -> githubusercontent.com/.../lol.jpg)
	link, err := category.makeLink(rawLinkT)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	link += "/" + articleSlug
	http.Redirect(w, r, link, http.StatusTemporaryRedirect)

}

func (config *Config) handleArticle(w http.ResponseWriter, r *http.Request) {
	categorySlug := chi.URLParam(r, "category")
	articleSlug := chi.URLParam(r, "article")

	category, ok := config.Categories[categorySlug]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if strings.HasSuffix(articleSlug, category.Ext) {
		newURL := fmt.Sprintf("/%s/%s/", categorySlug, strings.TrimSuffix(articleSlug, category.Ext))
		http.Redirect(w, r, newURL, http.StatusPermanentRedirect)
		return
	}

	category.Slug = categorySlug
	file := articleSlug + category.Ext
	article := Article{
		Category: category,
		File:     file,
		Slug:     articleSlug,
	}
	err := article.updateRaw()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	article.Title = article.getTitle()
	article.HTML = article.getHTML()

	if config.Contributors {
		err := article.updateMetaInfo()
		if err != nil {
			fmt.Println(err)
		}
	}

	if config.Lint {
		err := article.updateAlerts()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if config.Views {
		_, err := r.Cookie("viewed")
		if err != http.ErrNoCookie {
			// article already viewed
			err := article.updateViews()
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// article hasn't been viewed yet
			err := article.incrementViews()
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cookie := http.Cookie{
				Name:   "viewed",
				Value:  "1",
				Domain: r.Host,
				Path:   fmt.Sprintf("/%s/%s/", category.Slug, article.Slug),
				MaxAge: 3600 * 24,
			}
			http.SetCookie(w, &cookie)
		}
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

	r := chi.NewRouter()
	r.Get("/", conf.handleCategories)
	r.Get("/{category}/", conf.handleCategory)
	r.Get("/{category}/{article}/", conf.handleArticle)
	r.Get("/{category}/*", conf.handleArticleRedirect)

	if conf.Cache {
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

		http.Handle(conf.Root, cacheClient.Middleware(r))
	} else {
		http.Handle(conf.Root, r)
	}

	fmt.Printf("Ready to get connections on %s%s\n", conf.Listen, conf.Root)
	log.Fatal(http.ListenAndServe(conf.Listen, nil))
}
