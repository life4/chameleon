package chameleon

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type HandlerMain struct {
	Template *template.Template
	Server   *Server
}

func (h HandlerMain) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := ps.ByName("filepath")
	path = strings.TrimRight(path, "/")

	// get page
	page, err := h.Page(path)
	if err != nil {
		h.Server.Logger.Error("cannot get page", zap.Error(err), zap.String("path", path))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// increment views count
	_, err = r.Cookie("viewed")
	if err == http.ErrNoCookie {
		cookie := http.Cookie{
			Name:   "viewed",
			Value:  "1",
			Path:   r.URL.Path,
			MaxAge: 3600 * 24,
		}
		http.SetCookie(w, &cookie)
		page.Inc()
	} else if err != nil {
		h.Server.Logger.Error("cannot get cookie", zap.Error(err), zap.String("path", path))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render the page
	w.WriteHeader(page.Status())
	err = page.Render(w)
	if err != nil {
		h.Server.Logger.Error("cannot render page", zap.Error(err), zap.String("path", path))
		_, err = fmt.Fprint(w, err.Error())
		if err != nil {
			h.Server.Logger.Warn("cannot write response", zap.Error(err), zap.String("path", path))
		}
		return
	}
}

func (h HandlerMain) Page(urlPath string) (Page, error) {
	if strings.Contains(urlPath, "/.") {
		return Page403{}, nil
	}

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
			Views:    h.Server.Database.Views(p),
			Template: h.Template,
			Logger:   h.Server.Logger,
		}
		if urlPath != "" && urlPath != "/" {
			page.Parent = &Category{
				Repository: h.Server.Repository,
				Path:       p.Parent(),
			}
		}
		return page, nil
	}

	// raw file
	isfile, err := p.IsFile()
	if err != nil {
		return nil, err
	}
	if isfile {
		page := PageAsset{Path: p}
		return page, nil
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
			Views:    h.Server.Database.Views(p),
			Template: h.Template,
			Logger:   h.Server.Logger,
		}
		return page, nil
	}

	return Page404{}, nil
}