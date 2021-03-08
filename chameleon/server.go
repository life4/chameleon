package chameleon

import (
	"fmt"
	"net/http"
)

type Server struct {
	Repository Repository
	Database   *Database
}

func (s *Server) Init() error {
	err := s.Repository.Pull()
	if err != nil {
		return fmt.Errorf("cannot pull repo: %v", err)
	}
	s.Database = &Database{}
	err = s.Database.Open()
	if err != nil {
		return fmt.Errorf("cannot open database: %v", err)
	}
	return nil
}

func (s *Server) Close() error {
	return s.Database.Close()
}

func (s *Server) Serve() error {
	http.HandleFunc("/", s.redirect)
	http.Handle("/p/", Handler{Server: s, Template: TemplateArticle})
	http.Handle("/l/", Handler{Server: s, Template: TemplateLinter})
	return http.ListenAndServe("127.0.0.1:1337", nil)
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/p/", 301)
}
