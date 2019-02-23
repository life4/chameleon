package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/errata-ai/vale/check"
	"github.com/errata-ai/vale/core"
	"github.com/errata-ai/vale/lint"
	"github.com/tidwall/gjson"
)

// Author is a struct with contributor info
type Author struct {
	Login, Avatar string
}

// Article is a struct with article title and content
type Article struct {
	Category Category
	File     string
	Title    string
	Raw      string
	HTML     string
	Slug     string
	Authors  []Author
	Alerts   []core.Alert
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

func (article *Article) getAuthors() ([]Author, error) {
	link, err := article.Category.makeLink(commitsLinkT)
	if err != nil {
		return nil, err
	}
	link += "/" + article.File
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d (%s)", res.StatusCode, link)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	authorsMap := make(map[string]Author)
	var login string
	for _, subtree := range gjson.Get(string(content), "#.author").Array() {
		login = subtree.Get("login").String()
		if login != "" {
			authorsMap[login] = Author{Login: login, Avatar: subtree.Get("avatar_url").String()}
		}
	}
	var authors []Author
	for _, author := range authorsMap {
		authors = append(authors, author)
	}

	return authors, nil
}

func (article *Article) getTitle() string {
	title := strings.Split(article.Raw, "\n")[0]
	if strings.Index(title, "# ") != 0 {
		return article.File
	}
	article.Raw = strings.TrimPrefix(article.Raw, title)
	return title[2:]
}

func (article *Article) lintHTML() ([]core.Alert, error) {
	config := core.NewConfig()
	config.GBaseStyles = []string{"proselint", "write-good", "Joblint", "Spelling"}
	config.MinAlertLevel = 1
	config.InExt = ".html"
	linter := lint.Linter{Config: config, CheckManager: check.NewManager(config)}
	files, _ := linter.LintString(article.HTML)
	return files[0].SortedAlerts(), nil
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
