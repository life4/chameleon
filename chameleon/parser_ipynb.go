package chameleon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	Type   string   `json:"cell_type"`
	Source []string `json:"source"`
}

type JupyterParser struct{}

func (JupyterParser) HTML(raw []byte) (string, error) {
	// parse JSON
	var nb Notebook
	err := json.Unmarshal(raw, &nb)
	if err != nil {
		return "", fmt.Errorf("parse JSON: %v", err)
	}

	// prepare Markdown parser fo "markdown" cells
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

	// convert each cell into HTML
	result := make([]string, 0)
	for _, cell := range nb.Cells {
		if cell.Type == "markdown" {
			src := strings.Join(cell.Source, "\n")
			err := mdparser.Convert([]byte(src), &buf)
			if err != nil {
				return "", err
			}
			html := buf.String()
			buf.Reset()
			result = append(result, html)
		}
	}
	return strings.Join(result, "\n"), nil
}

// ExtractTitle extracts title from raw content
func (JupyterParser) ExtractTitle(raw []byte) (string, []byte) {
	return "", raw
}
