package chameleon

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//go:embed assets/*
var assets embed.FS

type Server struct {
	Repository Repository
	Database   *Database
	router     *httprouter.Router
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

	s.router = httprouter.New()
	s.router.Handler(
		http.MethodGet,
		"/",
		http.RedirectHandler("/p/", http.StatusTemporaryRedirect),
	)
	s.router.Handler(
		http.MethodGet,
		"/assets/*filepath",
		http.FileServer(http.FS(assets)),
	)
	s.router.GET(
		"/p/*filepath",
		Handler{Server: s, Template: TemplateArticle}.Handle,
	)
	s.router.GET(
		"/linter/*filepath",
		Handler{Server: s, Template: TemplateLinter}.Handle,
	)
	s.router.GET(
		"/commits/*filepath",
		Handler{Server: s, Template: TemplateCommits}.Handle,
	)

	return nil
}

func (s *Server) Close() error {
	return s.Database.Close()
}

func (s *Server) Serve() error {
	return http.ListenAndServe("127.0.0.1:1337", s.router)
}
