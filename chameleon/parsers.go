package chameleon

import (
	"regexp"
	"strings"
)

var rexImg = regexp.MustCompile("(<img src=\"(.*?)\".*?/>)")

type Parser interface {
	HTML(raw []byte) (string, error)
	ExtractTitle(raw []byte) (string, []byte)
}

func GetParser(path Path) Parser {
	if strings.HasSuffix(path.String(), ExtensionMarkdown) {
		return MarkdownParser{}
	}
	if strings.HasSuffix(path.String(), ExtensionJupyter) {
		return JupyterParser{}
	}
	return nil
}
