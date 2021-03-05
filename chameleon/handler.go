package chameleon

import (
	"fmt"
	"net/http"
)

type TemplateName string

const (
	TemplateArticle = TemplateName("templates/category.html.j2")
	TemplateLinter  = TemplateName("templates/linter.html.j2")
)

type Handler struct {
	Template TemplateName
	Server   *Server
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page, err := h.Page(r.URL.Path[3:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := page.Render(h.Template)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = page.Views.Inc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Page(urlPath string) (*Page, error) {
	p := h.Server.Repository.Path.Join(urlPath)

	// category page
	isdir, err := p.IsDir()
	if err != nil {
		return nil, err
	}
	if isdir {
		isfile, err := p.Join(ReadMe).IsFile()
		if err != nil {
			return nil, err
		}
		if !isfile {
			return nil, fmt.Errorf("README.md not found")
		}
		page := Page{
			Category: &Category{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Article: Article{
				Repository: h.Server.Repository,
				Path:       p.Join(ReadMe),
			},
			Templates: h.Server.Templates,
			Views:     h.Server.Database.Views(p),
		}
		return &page, nil
	}

	// article page
	isfile, err := p.IsFile()
	if err != nil {
		return nil, err
	}
	if isfile {
		page := Page{
			Parent: &Category{
				Repository: h.Server.Repository,
				Path:       p.Parent(),
			},
			Article: Article{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Templates: h.Server.Templates,
			Views:     h.Server.Database.Views(p),
		}
		return &page, nil
	}

	return nil, fmt.Errorf("file not found")
}
