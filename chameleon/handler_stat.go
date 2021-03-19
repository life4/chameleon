package chameleon

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerStat struct {
	Server *Server
}

func (h HandlerStat) Handle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := h.Render(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h HandlerStat) Render(w http.ResponseWriter) error {
	stat, err := h.Server.Database.Views("").All()
	if err != nil {
		return err
	}
	stat.Sort()
	return TemplateStat.Execute(w, stat)
}
