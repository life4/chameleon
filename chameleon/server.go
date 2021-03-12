package chameleon

import (
	"embed"
	"net/http"
	"net/http/pprof"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

//go:embed assets/*
var assets embed.FS

type Server struct {
	Repository Repository
	Database   *Database
	Logger     *zap.Logger
	router     *httprouter.Router
}

func (s *Server) Init(debug bool) {
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
		MainPrefix+"*filepath",
		HandlerMain{Server: s, Template: TemplateArticle}.Handle,
	)
	s.router.GET(
		LinterPrefix+"*filepath",
		HandlerMain{Server: s, Template: TemplateLinter}.Handle,
	)
	s.router.GET(
		CommitsPrefix+"*filepath",
		HandlerMain{Server: s, Template: TemplateCommits}.Handle,
	)
	s.router.GET(
		DiffPrefix+":hash",
		HandlerDiff{Server: s}.Handle,
	)

	if debug {
		s.Logger.Debug("debugging enabled", zap.String("endpoint", "/debug/pprof/"))
		s.router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
		s.router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
		s.router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
		s.router.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)
	}
}

func (s *Server) Close() error {
	return s.Database.Close()
}

func (s *Server) Serve(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
