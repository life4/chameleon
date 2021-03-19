package chameleon

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Auth struct {
	Password string
	Logger   *zap.Logger
}

func (a Auth) generate(date string) string {
	h := sha512.New().Sum([]byte(date + "|" + a.Password))
	return base64.StdEncoding.EncodeToString(h)
}

func (a Auth) valid(r *http.Request) (bool, error) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		return false, nil
	}
	parts := strings.SplitN(cookie.Value, "|", 2)
	if len(parts) != 2 {
		return false, errors.New("auth date not found")
	}

	given := []byte(parts[0])
	expected := []byte(a.generate(parts[1]))
	eq := subtle.ConstantTimeCompare(given, expected) == 1
	return eq, nil
}

func (a Auth) make(r *http.Request) *http.Cookie {
	date := time.Now().Format(time.RFC3339)
	token := a.generate(date) + "|" + date
	return &http.Cookie{
		Name:   "auth",
		Value:  token,
		Path:   "/",
		MaxAge: 3600 * 24 * 7,
	}
}

func (a Auth) render(w http.ResponseWriter, ok bool) {
	w.WriteHeader(http.StatusForbidden)
	_ = TemplateAuth.Execute(w, ok)
}

func (a Auth) Wrap(h httprouter.Handle) httprouter.Handle {
	if a.Password == "" {
		return h
	}
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		valid, err := a.valid(r)
		if err != nil {
			a.Logger.Debug("auth error", zap.Error(err), zap.String("ip", r.RemoteAddr))
			http.Redirect(w, r, AuthPrefix, http.StatusTemporaryRedirect)
			return
		}
		if !valid {
			http.Redirect(w, r, AuthPrefix, http.StatusTemporaryRedirect)
			return
		}
		h(w, r, ps)
	}
}

func (a Auth) HandleGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	valid, _ := a.valid(r)
	if valid {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	a.render(w, true)
}

func (a Auth) HandlePOST(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		a.Logger.Debug("cannot parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	given := []byte(r.PostForm.Get("password"))
	expected := []byte(a.Password)
	ok := subtle.ConstantTimeCompare(given, expected) == 1
	if !ok {
		a.render(w, false)
		return
	}
	http.SetCookie(w, a.make(r))
	http.Redirect(w, r, MainPrefix, http.StatusSeeOther)
}
