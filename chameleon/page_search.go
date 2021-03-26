package chameleon

import (
	"io"
	"net/http"
)

type PageSearch struct {
	Query   string
	Results []*Article
}

func (p PageSearch) Render(w io.Writer) error {
	return TemplateSearch.Execute(w, &p)
}

func (p PageSearch) Inc() {}

func (p PageSearch) Status() int {
	return http.StatusOK
}
