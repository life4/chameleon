package chameleon

import (
	"fmt"
	"io/fs"
	"net/http"
)

type Server struct {
	Repository Repository
	Templates  fs.FS
}

func (s *Server) Init() error {
	err := s.Repository.Pull()
	if err != nil {
		return fmt.Errorf("cannot pull repo: %v", err)
	}
	return nil
}

func (s *Server) Serve() error {
	http.HandleFunc("/", s.handle)
	return http.ListenAndServe("127.0.0.1:1337", nil)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	page, err := s.Page(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := page.Render()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(content))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s Server) Page(urlPath string) (*Page, error) {
	// root page
	if urlPath == "" || urlPath == "/" {
		page := Page{
			Article: Article{
				Repository: s.Repository,
				FileName:   string(s.Repository.Path.Join(ReadMe)),
			},
			Templates: s.Templates,
		}
		return &page, nil
	}

	p := s.Repository.Path.Join(urlPath)

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
				Repository: s.Repository,
				DirName:    string(p.Relative(s.Repository.Path)),
			},
			Article: Article{
				Repository: s.Repository,
				FileName:   string(p.Join(ReadMe).Relative(s.Repository.Path)),
			},
			Templates: s.Templates,
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
			Traceback: []Category{{
				Repository: s.Repository,
				DirName:    string(p.Parent()),
			}},
			Article: Article{
				Repository: s.Repository,
				FileName:   string(p.Relative(s.Repository.Path)),
			},
			Templates: s.Templates,
		}
		return &page, nil
	}

	return nil, fmt.Errorf("file not found")
}
