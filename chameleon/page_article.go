package chameleon

import (
	"io"
	"net/http"
	"text/template"

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
	return p.Template.Execute(w, &p)
}

func (p PageArticle) Inc() {
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
