package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/errata-ai/vale/check"
	"github.com/errata-ai/vale/core"
	"github.com/errata-ai/vale/lint"
	"github.com/recoilme/pudge"
	"github.com/recoilme/slowpoke"
	"github.com/tidwall/gjson"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// Author is a struct with contributor info
type Author struct {
	Login, Avatar string
}

// Article is a struct with article title and content
type Article struct {
	Category *Category
	Config   *Config

	File  string
	Title string
	Raw   string
	HTML  string
	Slug  string

	Authors   []Author
	CreatedAt time.Time
	UpdatedAt time.Time
	Alerts    []core.Alert

	Views uint32
}

func (article *Article) updateRaw() error {
	link, err := article.Category.makeLink(rawLinkT)
	if err != nil {
		return err
	}
	link += "/" + article.File
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	article.Raw = string(content)
	return nil
}

func (article *Article) updateMetaInfo() error {
	link, err := article.Category.makeLink(commitsLinkT)
	if err != nil {
		return err
	}
	link += "/" + article.File
	res, err := http.Get(link)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("invalid status code: %d (%s)", res.StatusCode, link)
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
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
	article.Authors = authors

	times := gjson.Get(string(content), "#.commit.author.date").Array()
	t, err := time.Parse(time.RFC3339, times[0].String())
	if err != nil {
		return err
	}
	article.UpdatedAt = t

	t, err = time.Parse(time.RFC3339, times[len(times)-1].String())
	if err != nil {
		return err
	}
	article.CreatedAt = t
	return nil
}

func (article *Article) getTitle() string {
	title := strings.Split(article.Raw, "\n")[0]
	if strings.Index(title, "# ") != 0 {
		return article.File
	}
	article.Raw = strings.TrimPrefix(article.Raw, title)
	return title[2:]
}

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")

func (article *Article) getHTML() (html string) {
	// convert markdown to html
	html = string(blackfriday.Run([]byte(article.Raw)))
	// fix relative paths
	html = strings.Replace(html, "src=\"./", "src=\"../", -1)
	html = strings.Replace(html, "href=\"./", "href=\"../", -1)
	// wrap images into link
	html = rexImg.ReplaceAllString(html, "<a href=\"$2\" target=\"_blank\">$1</a>")
	return
}

func (article *Article) updateAlerts() error {
	config := core.NewConfig()
	config.GBaseStyles = []string{"proselint", "write-good", "Joblint", "Spelling"}
	config.MinAlertLevel = 1
	config.InExt = ".html"
	linter := lint.Linter{Config: config, CheckManager: check.NewManager(config)}
	files, _ := linter.LintString(article.HTML)
	article.Alerts = files[0].SortedAlerts()
	return nil
}

func (article *Article) getFilename() string {
	return path.Join(article.Config.Project, ".storage", article.Category.Slug+".db")
}

func (article *Article) updateViews() error {
	if article.Views != 0 {
		return nil
	}
	bs, err := slowpoke.Get(article.getFilename(), []byte(article.Slug))
	if err == pudge.ErrKeyNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	article.Views = binary.LittleEndian.Uint32(bs)
	return nil
}

func (article *Article) incrementViews() error {
	err := article.updateViews()
	if err != pudge.ErrKeyNotFound && err != nil {
		return err
	}
	article.Views++
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, article.Views)
	return slowpoke.Set(article.getFilename(), []byte(article.Slug), bs)
}
