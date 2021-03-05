package chameleon

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/enescakir/emoji"
	"gopkg.in/russross/blackfriday.v2"
)

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")

type Article struct {
	Repository Repository
	Path       Path

	// cache
	raw   []byte
	title string
}

func (a Article) Valid() (bool, error) {
	if !strings.HasSuffix(a.Path.String(), Extension) {
		return false, nil
	}
	return a.Path.IsFile()
}

func (a Article) IsReadme() bool {
	return a.Path.Name() == ReadMe
}

func (a *Article) Linter() Linter {
	return Linter{Article: a}
}

func (a *Article) Raw() ([]byte, error) {
	if a.raw != nil {
		return a.raw, nil
	}
	raw, err := os.ReadFile(a.Path.String())
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}
	a.raw = raw
	a.trimTitle()
	return raw, nil
}

func (a *Article) Content() (string, error) {
	raw, err := a.Raw()
	return string(raw), err
}

func (a *Article) HTML() (string, error) {
	raw, err := a.Raw()
	if err != nil {
		return "", err
	}
	html := string(blackfriday.Run(raw))
	html = emoji.Parse(html)
	// fix relative paths
	html = strings.Replace(html, "src=\"./", "src=\"../", -1)
	html = strings.Replace(html, "href=\"./", "href=\"../", -1)
	// wrap images into link
	html = rexImg.ReplaceAllString(html, "<a href=\"$2\" target=\"_blank\">$1</a>")
	return html, nil
}

// trimTitle extracts title from raw content
func (a *Article) trimTitle() {
	title := bytes.SplitN(a.raw, []byte{'\n'}, 2)[0]
	if bytes.Index(title, []byte{'#', ' '}) != 0 {
		if a.IsReadme() {
			a.title = a.Path.Parent().Name()
			return
		}
		a.title = a.Path.Name()
		return
	}
	a.raw = bytes.TrimPrefix(a.raw, title)
	a.title = strings.TrimSuffix(string(title[2:]), "\n")
}

func (a *Article) Title() (string, error) {
	if a.title == "" {
		_, err := a.Raw()
		if err != nil {
			return "", err
		}
	}
	return a.title, nil
}
