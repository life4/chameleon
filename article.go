package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Article is a struct with article title and content
type Article struct {
	Category  Category
	File      string
	Title     string
	Raw       string
	HTML      string
	Slug string
}

func (article *Article) getRaw() (string, error) {
	link, err := article.Category.makeLink(rawLinkT)
	if err != nil {
		return "", err
	}
	link += "/" + article.File
	res, err := http.Get(link)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (article *Article) getTitle() string {
	title := strings.Split(article.Raw, "\n")[0]
	if strings.Index(title, "# ") != 0 {
		return article.File
	}
	article.Raw = strings.TrimPrefix(article.Raw, title)
	return title[2:]
}

func (article *Article) init() error {
	raw, err := article.getRaw()
	if err != nil {
		return err
	}
	article.Raw = raw
	article.Title = article.getTitle()
	return nil
}
