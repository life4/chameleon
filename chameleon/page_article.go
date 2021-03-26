package chameleon

import (
	"io"
	"net/http"
	"text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"go.uber.org/zap"
)

type PageArticle struct {
	Article  Article
	Parent   *Category
	Category *Category
	Views    *Views
	Template *template.Template
	Logger   *zap.Logger
}

func (p PageArticle) Render(w io.Writer) error {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	r, interw := io.Pipe()
	var errt error
	go func() {
		errt = p.Template.Execute(interw, &p)
		interw.Close()
	}()
	errm := m.Minify("text/html", w, r)
	if errt != nil {
		return errt
	}
	return errm
}

func (p PageArticle) Inc() {
	if p.Views == nil {
		return
	}
	go func() {
		err := p.Views.Inc()
		if err != nil {
			p.Logger.Error("cannot increment views", zap.Error(err))
		}
	}()
}

func (p PageArticle) Status() int {
	return http.StatusOK
}
