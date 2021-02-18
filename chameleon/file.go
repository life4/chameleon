package chameleon

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")

type File struct {
	string
}

func (f File) Markdown() ([]byte, error) {
	return os.ReadFile(f.string)
}

func (f File) HTML() (string, error) {
	md, err := f.Markdown()
	if err != nil {
		return "", err
	}
	html := string(blackfriday.Run(md))
	// fix relative paths
	html = strings.Replace(html, "src=\"./", "src=\"../", -1)
	html = strings.Replace(html, "href=\"./", "href=\"../", -1)
	// wrap images into link
	html = rexImg.ReplaceAllString(html, "<a href=\"$2\" target=\"_blank\">$1</a>")
	return html, nil
}
