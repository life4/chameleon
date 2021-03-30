package chameleon

import (
	"io"
	"net/http"
)

type Page403 struct {
}

func (p Page403) Render(w io.Writer) error {
	return Template403.Execute(w, &p)
}

func (p Page403) Inc() {}

func (p Page403) Status() int {
	return http.StatusForbidden
}
