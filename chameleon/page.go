package chameleon

import (
	"bytes"
	"io/fs"
	"text/template"
)

type Page struct {
	Article   Article
	Parent    *Category
	Category  *Category
	Templates fs.FS
	Views     *Views
}

func (p Page) Render(tname TemplateName) (string, error) {
	t, err := template.ParseFS(
		p.Templates,
		"templates/base.html.j2",
		string(tname),
	)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, &p)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
