package chameleon

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

var rexSafeChars = regexp.MustCompile(`[^a-zA-Z\-\.\_ ]`)

type HandlerSearch struct {
	Server *Server
}

func (h HandlerSearch) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := r.URL.Query()["q"]
	if len(query) == 0 {
		http.Error(w, "empty query", http.StatusBadRequest)
		return
	}
	err := h.Render(w, query[0])
	if err != nil {
		h.Server.Logger.Error("cannot handle search query", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h HandlerSearch) Render(w http.ResponseWriter, query string) error {
	query = rexSafeChars.ReplaceAllString(query, "")
	cmd := h.Server.Repository.Command(
		"grep",
		"--full-name",
		"--ignore-case",
		"--recursive",
		"--word-regexp",
		"--name-only",
		"--fixed-strings",
		query,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}

	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	arts := make([]*Article, 0)
	for _, path := range paths {
		if path == "" {
			continue
		}
		art := &Article{
			Repository: h.Server.Repository,
			Path:       h.Server.Repository.Join(path),
		}
		valid, err := art.Valid()
		if err != nil {
			return err
		}
		if !valid {
			continue
		}
		arts = append(arts, art)
	}
	return TemplateSearch.Execute(w, &arts)
}
