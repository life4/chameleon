package chameleon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	htmlRenderer "github.com/yuin/goldmark/renderer/html"
)

type Notebook struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	Type    string   `json:"cell_type"`
	Source  []string `json:"source"`
	Outputs []Output `json:"outputs"`
}

type Output struct {
	Data Data
}

type Data struct {
	HTML  []string `json:"text/html"`
	Plain []string `json:"text/plain"`
	PNG   string   `json:"image/png"`
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
			extension.Typographer,
		),
		goldmark.WithRendererOptions(
			htmlRenderer.WithUnsafe(),
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
			result = append(result, buf.String())
			buf.Reset()
			continue
		}

		if cell.Type == "code" {
			src := strings.Join(cell.Source, "")
			line := fmt.Sprintf(`<pre><code class="language-python">%s</code></pre>`, src)
			result = append(result, line)

			hasOut := false
			for _, output := range cell.Outputs {
				if output.Data.PNG != "" {
					src := fmt.Sprintf(`<img src="data:image/png;base64,%s"/>`, output.Data.PNG)
					result = append(result, src)
					continue
				}

				src := strings.Join(output.Data.Plain, "")
				if src != "" {
					if !hasOut {
						hasOut = true
						result = append(result, `<div class="card text-dark bg-light">`)
						result = append(result, `<div class="card-body">`)
					}
					src = html.EscapeString(src)
					src := fmt.Sprintf("<pre>%s</pre>", src)
					result = append(result, src)
				}
			}
			if hasOut {
				result = append(result, `</div></div>`)
			}
			continue
		}

	}
	return strings.Join(result, "\n"), nil
}

// ExtractTitle extracts title from raw content
func (JupyterParser) ExtractTitle(raw []byte) (string, []byte) {
	return "", raw
}
