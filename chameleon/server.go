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
	s.Database = &Database{}
	err := s.Database.Open()
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
		ArticlePrefix+"*filepath",
		Handler{Server: s, Template: TemplateArticle}.Handle,
	)
	s.router.GET(
		LinterPrefix+"*filepath",
		Handler{Server: s, Template: TemplateLinter}.Handle,
	)
	s.router.GET(
		CommitsPrefix+"*filepath",
		Handler{Server: s, Template: TemplateCommits}.Handle,
	)

	return nil
}

func (s *Server) Close() error {
	return s.Database.Close()
}

func (s *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
