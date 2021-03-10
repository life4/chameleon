package chameleon

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/enescakir/emoji"
	"gopkg.in/russross/blackfriday.v2"
)

const ISO8601 = "2006-01-02T15:04:05-07:00"

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

func (a Article) Slug() string {
	return strings.TrimSuffix(a.Path.Name(), Extension)
}

func (a Article) Commits() ([]Commit, error) {
	p := a.Path.Relative(a.Repository.Path)
	cmd := a.Repository.Command("log", "--pretty=%H|%cI|%an|%ae", p.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	commits := make([]Commit, len(lines))
	for i, line := range lines {
		line := strings.TrimSpace(line)
		parts := strings.Split(line, "|")
		t, err := time.Parse(ISO8601, parts[1])
		if err != nil {
			return nil, err
		}
		commits[i] = Commit{
			Hash: parts[0],
			Time: t,
			Name: parts[2],
			Mail: parts[3],
		}
	}
	return commits, nil
}

func (a Article) URLs() URLs {
	return URLs{
		Repository: a.Repository,
		Path:       a.Path,
	}
}
