package chameleon

import (
	"io"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
)

type Handlers struct {
	Repository Repository
	Templates  fs.FS
}

func (h Handlers) Register(prefix string, router chi.Router) {
	router.HandleFunc(prefix, h.handleRoot)
}

func (h Handlers) handleRoot(w http.ResponseWriter, r *http.Request) {
	err := h.renderRoot(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handlers) renderRoot(w io.Writer) error {
	paths, err := h.Repository.Path().SubPaths()
	if err != nil {
		return err
	}
	cats := make([]Category, 0)
	for _, p := range paths {
		println(p.String())
		isdir, err := p.IsDir()
		if err != nil {
			return err
		}
		if !isdir {
			continue
		}
		cat := Category{
			Repository: h.Repository,
			DirName:    p.Name()}
		cats = append(cats, cat)
	}

	t, err := template.ParseFS(
		h.Templates,
		"templates/base.html",
		"templates/categories.html",
	)
	if err != nil {
		return err
	}
	return t.Execute(w, cats)
}
