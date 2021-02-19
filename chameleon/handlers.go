package chameleon

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
)

type TemplateContext struct {
	Article    *Article
	Categories []Category
}

type Handlers struct {
	Repository Repository
	Templates  fs.FS
}

func (h Handlers) Register(prefix string, router chi.Router) {
	router.HandleFunc(prefix, h.handleRoot)
}

func (h Handlers) handleRoot(w http.ResponseWriter, r *http.Request) {
	executor, err := h.renderRoot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = executor(w)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func (h Handlers) renderRoot() (func(io.Writer) error, error) {
	paths, err := h.Repository.Path().SubPaths()
	if err != nil {
		return nil, err
	}
	cats := make([]Category, 0)
	for _, p := range paths {
		isdir, err := p.IsDir()
		if err != nil {
			return nil, err
		}
		if !isdir {
			continue
		}
		cat := Category{
			Repository: h.Repository,
			DirName:    p.Name(),
		}
		hasReadme, err := cat.HasReadme()
		if err != nil {
			return nil, err
		}
		if !hasReadme {
			continue
		}
		cats = append(cats, cat)
	}

	a := h.Repository.ReadMe()
	ctx := TemplateContext{
		Article:    &a,
		Categories: cats,
	}

	t, err := template.ParseFS(
		h.Templates,
		"templates/base.html",
		"templates/category.html",
	)
	if err != nil {
		return nil, err
	}
	executor := func(w io.Writer) error {
		return t.Execute(w, ctx)
	}
	return executor, nil
}
