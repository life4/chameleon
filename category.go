package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
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
	var article Article
	for _, name := range names.Array() {
		article = Article{
			Category: *category,
			File:     name.String(),
		}
		article.Raw, err = article.getRaw()
		if err != nil {
			return
		}
		article.Title = article.getTitle()
		articles = append(articles, article)
	}
	return articles, nil
}
