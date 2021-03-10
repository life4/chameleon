package chameleon

import (
	"io"
	"net/http"
	"text/template"
)

type PageArticle struct {
	Article  Article
	Parent   *Category
	Category *Category
	Views    *Views
	Template *template.Template
}

func (p PageArticle) Render(w io.Writer) error {
	return p.Template.Execute(w, &p)
}

func (p PageArticle) Inc() error {
	return p.Views.Inc()
}

func (p PageArticle) Status() int {
	return http.StatusOK
}
