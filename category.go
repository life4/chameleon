package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/tidwall/gjson"
)

// Category is a struct with all information about content in given category
type Category struct {
	Repo, Branch, Dir, Name, Ext, Slug string
}

func (category *Category) makeLink(linkTemplate string) (string, error) {
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

func (category *Category) getArticles() (articles []Article, err error) {
	link, err := category.makeLink(dirAPILinkT)
	if err != nil {
		return
	}

	// get filenames
	res, err := http.Get(link)
	if err != nil {
		return
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return
	}
	names := gjson.Get(string(content), "#.name")

	// get content
	var wg sync.WaitGroup
	var article Article
	articlesChan := make(chan Article, len(names.Array()))
	for _, name := range names.Array() {
		if strings.HasSuffix(name.String(), category.Ext) {
			wg.Add(1)
			article = Article{
				Category: category,
				File:     name.String(),
				Slug:     strings.TrimSuffix(name.String(), category.Ext),
			}
			go func(article Article) {
				defer wg.Done()
				err := article.updateRaw()
				if err != nil {
					log.Fatal(err)
					return
				}
				article.Title = article.getTitle()
				articlesChan <- article
			}(article)
		}
	}
	wg.Wait()

	close(articlesChan)
	for article := range articlesChan {
		articles = append(articles, article)
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Title < articles[j].Title
	})

	return articles, nil
}
