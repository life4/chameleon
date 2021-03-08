package chameleon

import (
	"bytes"
	"text/template"
)

type Page struct {
	Article  Article
	Parent   *Category
	Category *Category
	Views    *Views
}

func (p Page) Render(t *template.Template) (string, error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, &p)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
