package chameleon

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	Config    Config
	Router    chi.Router
	Templates fs.FS

	repos []Repository
}

func (s *Server) InitRepos() error {
	storage := s.Config.CachePath
	err := storage.EnsureDir()
	if err != nil {
		return fmt.Errorf("cannot ensure dir: %v", err)
	}
	s.repos = make([]Repository, len(s.Config.Repos))
	for _, crepo := range s.Config.Repos {
		repo := Repository{URL: crepo.URL, Storage: storage}
		s.repos = append(s.repos, repo)

		err = repo.Clone()
		if err != nil {
			return fmt.Errorf("cannot clone repo: %v", err)
		}
		err = repo.Pull()
		if err != nil {
			return fmt.Errorf("cannot pull repo: %v", err)
		}

		h := Handlers{Repository: repo, Templates: s.Templates}
		h.Register(crepo.URLPath, s.Router)
	}
	return nil
}

func (s *Server) Serve() error {
	http.Handle("/", s.Router)
	return http.ListenAndServe("127.0.0.1:1337", nil)
}
