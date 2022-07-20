package chameleon

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/enescakir/emoji"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")

type Parser interface {
	HTML(raw []byte) (string, error)
	ExtractTitle(raw []byte) (string, []byte)
}

func GetParser(path Path) Parser {
	return MarkdownParser{}
}

type MarkdownParser struct{}

func (MarkdownParser) HTML(raw []byte) (string, error) {
	mdparser := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	err := mdparser.Convert(raw, &buf)
	if err != nil {
		return "", err
	}
	html := buf.String()
	html = emoji.Parse(html)
	// fix relative paths
	html = strings.ReplaceAll(html, "src=\"./", "src=\"../")
	html = strings.ReplaceAll(html, "href=\"./", "href=\"../")
	// wrap images into link
	html = rexImg.ReplaceAllString(html, "<a href=\"$2\" target=\"_blank\">$1</a>")
	return html, nil
}

// ExtractTitle extracts title from raw content
func (MarkdownParser) ExtractTitle(raw []byte) (string, []byte) {
	title := bytes.SplitN(raw, []byte{'\n'}, 2)[0]
	if bytes.Index(title, []byte{'#', ' '}) != 0 {
		return "", raw
	}
	raw = bytes.TrimPrefix(raw, title)
	titleS := strings.ReplaceAll(strings.TrimSuffix(string(title[2:]), "\n"), "`", "")
	return titleS, raw
}
