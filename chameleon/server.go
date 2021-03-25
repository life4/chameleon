package chameleon

import (
	"embed"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

//go:embed assets/*
var assets embed.FS

type Server struct {
	Repository Repository
	Database   *Database
	Logger     *zap.Logger
	Config     Config

	cache  *Cache
	auth   Auth
	router *httprouter.Router
}

func NewServer(config Config, logger *zap.Logger) (*Server, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	repo := Repository{Path: Path(config.RepoPath)}

	cache, err := NewCache(config.Cache)
	if err != nil {
		return nil, fmt.Errorf("cannot create cache: %v", err)
	}

	server := Server{
		Repository: repo,
		Database:   &Database{},
		Logger:     logger,
		cache:      cache,
		auth: Auth{
			Password: config.Password,
			Logger:   logger,
		},
		Config: config,
	}

	if config.DBPath != "" {
		err = server.Database.Open(config.DBPath)
		if err != nil {
			return nil, fmt.Errorf("cannot open database: %v", err)
		}
	}

	err = repo.Clone(config.RepoURL)
	if err != nil {
		return nil, fmt.Errorf("cannot clone repo: %v", err)
	}

	server.Init()

	if config.Pull != 0 {
		sch := gocron.NewScheduler(time.UTC)
		job, err := sch.Every(config.Pull).Do(func() {
			logger.Debug("pulling the repo")
			err := repo.Pull()
			if err != nil {
				logger.Error("cannot pull the repo", zap.Error(err))
			}
		})
		if err != nil {
			return nil, fmt.Errorf("cannot schedule the task: %v", err)
		}
		job.SingletonMode()
		sch.StartAsync()
	}

	return &server, nil
}

func (s *Server) Init() {
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
		s.auth.Wrap(
			s.cache.Wrap(
				HandlerMain{Server: s, Template: TemplateArticle}.Handle,
			),
		),
	)
	s.router.GET(
		LinterPrefix+"*filepath",
		s.auth.Wrap(
			s.cache.Wrap(
				HandlerMain{Server: s, Template: TemplateLinter}.Handle,
			),
		),
	)
	s.router.GET(
		CommitsPrefix+"*filepath",
		s.auth.Wrap(
			HandlerMain{Server: s, Template: TemplateCommits}.Handle,
		),
	)
	s.router.GET(
		DiffPrefix+":hash",
		s.auth.Wrap(
			HandlerDiff{Server: s}.Handle,
		),
	)
	s.router.GET(
		StatPrefix,
		s.auth.Wrap(
			HandlerStat{Server: s}.Handle,
		),
	)
	s.router.GET(
		SearchPrefix,
		s.auth.Wrap(
			HandlerSearch{Server: s}.Handle,
		),
	)

	s.router.GET(AuthPrefix, s.auth.HandleGET)
	s.router.POST(AuthPrefix, s.auth.HandlePOST)

	if s.Config.PProf {
		s.Logger.Debug("debugging enabled", zap.String("endpoint", "/debug/pprof/"))
		s.router.HandlerFunc("GET", "/debug/pprof/", pprof.Index)
		s.router.HandlerFunc("GET", "/debug/pprof/cmdline", pprof.Cmdline)
		s.router.HandlerFunc("GET", "/debug/pprof/profile", pprof.Profile)
		s.router.HandlerFunc("GET", "/debug/pprof/symbol", pprof.Symbol)
		s.router.HandlerFunc("GET", "/debug/pprof/trace", pprof.Trace)
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
