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

const ISO8601 = "2006-01-02T15:04:05-07:00"

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")
var rexLang = regexp.MustCompile("<code class=\"language-([a-zA-Z]+)\">")

type Article struct {
	Repository Repository
	Path       Path

	// cache
	raw     []byte
	title   string
	commits Commits
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
	var err error
	if a.raw != nil {
		return a.raw, nil
	}
	a.raw, err = os.ReadFile(a.Path.String())
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}
	a.trimTitle()
	return a.raw, nil
}

func (a *Article) HTML() (string, error) {
	raw, err := a.Raw()
	if err != nil {
		return "", err
	}
	html := string(blackfriday.Run(raw))
	html = emoji.Parse(html)
	// fix relative paths
	html = strings.ReplaceAll(html, "src=\"./", "src=\"../")
	html = strings.ReplaceAll(html, "href=\"./", "href=\"../")
	// wrap images into link
	html = rexImg.ReplaceAllString(html, "<a href=\"$2\" target=\"_blank\">$1</a>")
	return html, nil
}

func (a *Article) Languages() ([]string, error) {
	html, err := a.HTML()
	if err != nil {
		return nil, err
	}
	matches := rexLang.FindAllStringSubmatch(html, -1)
	result := make([]string, 0)
	set := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		lang := m[1]
		_, ok := set[lang]
		if !ok {
			set[lang] = struct{}{}
			result = append(result, lang)
		}
	}
	return result, nil
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

func (a *Article) Commits() (Commits, error) {
	if a.commits != nil {
		return a.commits, nil
	}

	p := a.Path.Relative(a.Repository.Path)
	cmd := a.Repository.Command("log", "--pretty=%H|%cI|%an|%ae|%s", "--follow", p.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	a.commits = make(Commits, len(lines))
	for i, line := range lines {
		a.commits[i], err = ParseCommit(line)
		if err != nil {
			return nil, err
		}
	}
	return a.commits, nil
}

func (a Article) URLs() URLs {
	return URLs{
		Repository: a.Repository,
		Path:       a.Path,
	}
}
