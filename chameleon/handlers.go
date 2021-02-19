package chameleon

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
)

type TemplateContext struct {
	Article    *Article
	Categories []Category
	Articles   []Article
}

type Handlers struct {
	Repository Repository
	Templates  fs.FS
}

func (h Handlers) Register(prefix string, router chi.Router) {
	router.HandleFunc(prefix, h.handleRoot)
	router.HandleFunc(prefix+"*", h.handleSubPath)
}

func (h Handlers) handleRoot(w http.ResponseWriter, r *http.Request) {
	executor, err := h.render("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = executor(w)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func (h Handlers) handleSubPath(w http.ResponseWriter, r *http.Request) {
	suffix := chi.URLParam(r, "*")
	executor, err := h.render(suffix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = executor(w)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func (h Handlers) render(suffix string) (func(io.Writer) error, error) {
	currentPath := h.Repository.Path().Join(suffix)
	isdir, err := currentPath.IsDir()
	if err != nil {
		return nil, err
	}

	cats := make([]Category, 0)
	arts := make([]Article, 0)
	if isdir {
		paths, err := currentPath.SubPaths()
		if err != nil {
			return nil, err
		}
		for _, p := range paths {
			isdir, err := p.IsDir()
			if err != nil {
				return nil, err
			}
			if isdir {
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
				continue
			}

			isfile, err := p.IsFile()
			if err != nil {
				return nil, err
			}
			if isfile {
				art := Article{
					Repository: h.Repository,
					FileName:   string(p.Relative(h.Repository.Path())),
				}
				if !art.IsMarkdown() {
					continue
				}
				if strings.HasSuffix(p.Name(), ReadMe) {
					continue
				}
				arts = append(arts, art)
				continue
			}
		}
	}

	articlePath := currentPath
	isdir, err = articlePath.IsDir()
	if err != nil {
		return nil, err
	}
	if isdir {
		articlePath = articlePath.Join(ReadMe)
	}
	a := Article{
		Repository: h.Repository,
		FileName:   string(articlePath.Relative(h.Repository.Path())),
	}
	ctx := TemplateContext{
		Article:    &a,
		Categories: cats,
		Articles:   arts,
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
