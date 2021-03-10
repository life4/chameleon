package chameleon

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Template *template.Template
	Server   *Server
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := ps.ByName("filepath")
	if !strings.HasSuffix(path, "/") {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusTemporaryRedirect)
		return
	}
	path = strings.TrimRight(path, "/")
	if strings.HasSuffix(path, Extension) {
		url := "/p/" + strings.TrimSuffix(path, Extension)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	page, err := h.Page(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = page.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = page.Inc()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) Page(urlPath string) (Page, error) {
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
		page := PageArticle{
			Category: &Category{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Article: Article{
				Repository: h.Server.Repository,
				Path:       p.Join(ReadMe),
			},
			Views: h.Server.Database.Views(p),
			URLs: URLs{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Template: h.Template,
		}
		if urlPath != "" && urlPath != "/" {
			page.Parent = &Category{
				Repository: h.Server.Repository,
				Path:       p.Parent(),
			}
		}
		return &page, nil
	}

	// raw file
	isfile, err := p.IsFile()
	if err != nil {
		return nil, err
	}
	if isfile {
		page := PageAsset{Path: p}
		return &page, nil
	}

	// article page
	p = h.Server.Repository.Path.Join(urlPath + Extension)
	isfile, err = p.IsFile()
	if err != nil {
		return nil, err
	}
	if isfile {
		page := PageArticle{
			Parent: &Category{
				Repository: h.Server.Repository,
				Path:       p.Parent(),
			},
			Article: Article{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Views: h.Server.Database.Views(p),
			URLs: URLs{
				Repository: h.Server.Repository,
				Path:       p,
			},
			Template: h.Template,
		}
		return &page, nil
	}

	page := Page404{Views: h.Server.Database.Views(p)}
	return page, nil
}
