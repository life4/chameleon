package chameleon

import (
	"io"
)

type Page404 struct {
	Views *Views
}

func (p Page404) Render(w io.Writer) error {
	return Template404.Execute(w, &p)
}

func (p Page404) Inc() error {
	return p.Views.Inc()
}
