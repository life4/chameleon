package chameleon

import (
	"io"
	"net/http"
)

type Page404 struct {
}

func (p Page404) Render(w io.Writer) error {
	return Template404.Execute(w, &p)
}

func (p Page404) Inc() error {
	return nil
}

func (p Page404) Status() int {
	return http.StatusNotFound
}
