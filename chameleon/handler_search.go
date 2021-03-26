package chameleon

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type HandlerSearch struct {
	Server *Server
}

func (h HandlerSearch) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	page, err := h.Page(r)
	if err != nil {
		h.Server.Logger.Error("cannot get page", zap.Error(err), zap.String("url", r.URL.RawPath))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(page.Status())
	err = page.Render(w)
	if err != nil {
		h.Server.Logger.Error("cannot render page", zap.Error(err), zap.String("url", r.URL.RawPath))
		_, err = fmt.Fprint(w, err.Error())
		if err != nil {
			h.Server.Logger.Warn("cannot write response", zap.Error(err), zap.String("url", r.URL.RawPath))
		}
		return
	}
}

func (h HandlerSearch) Page(r *http.Request) (Page, error) {
	queries, ok := r.URL.Query()["q"]
	if !ok || len(queries) == 0 {
		return PageSearch{}, nil
	}

	query := safeText(queries[0])
	if query == "" {
		return PageSearch{}, nil
	}

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
	if cmd.ProcessState.ExitCode() == 1 {
		return PageSearch{Query: query}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("%v: %s", err, out)
	}

	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	page := PageSearch{
		Query:   query,
		Results: make([]*Article, 0),
	}
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
			return nil, err
		}
		if !valid {
			continue
		}
		page.Results = append(page.Results, art)
	}
	return page, nil
}
