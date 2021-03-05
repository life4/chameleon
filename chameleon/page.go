package chameleon

import (
	"bytes"
	"io/fs"
	"text/template"
)

type Page struct {
	Article   Article
	Traceback []Category
	Category  *Category
	Templates fs.FS
}

func (p Page) Render() (string, error) {
	t, err := template.ParseFS(
		p.Templates,
		"templates/base.html",
		"templates/category.html",
	)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, p)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
