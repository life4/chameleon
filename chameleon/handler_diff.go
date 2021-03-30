package chameleon

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var rexHash = regexp.MustCompile(`^[a-f0-9]{40}$`)

type HandlerDiff struct {
	Server *Server
}

func (h HandlerDiff) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := h.Render(w, ps.ByName("hash"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h HandlerDiff) Render(w http.ResponseWriter, hash string) error {
	if !rexHash.MatchString(hash) {
		return errors.New("invalid commit hash")
	}
	cmd := h.Server.Repository.Command("show", "--pretty=%H|%cI|%an|%ae|%s", hash)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}

	lines := strings.SplitN(strings.TrimSpace(string(out)), "\n", 2)
	commit, err := ParseCommit(lines[0])
	if err != nil {
		return err
	}
	commit.Diff = lines[1]
	return TemplateDiff.Execute(w, &commit)
}
