package chameleon

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthNoRedirect(t *testing.T) {
	is := require.New(t)
	request, err := http.NewRequest(http.MethodGet, "/p/", nil)
	is.Nil(err)
	response := httptest.NewRecorder()

	s, err := NewServer(newTestConfig(), nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	is.Equal(response.Code, 200)
}

func TestAuthRedirect(t *testing.T) {
	is := require.New(t)
	request, err := http.NewRequest(http.MethodGet, "/p/", nil)
	is.Nil(err)
	response := httptest.NewRecorder()

	c := newTestConfig()
	c.Password = "123"
	s, err := NewServer(c, nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	is.Equal(response.Code, 307)
	is.Equal(response.Header()["Location"], []string{"/auth/"})
}

func TestAuthForm(t *testing.T) {
	is := require.New(t)
	request, err := http.NewRequest(http.MethodGet, "/auth/", nil)
	is.Nil(err)
	response := httptest.NewRecorder()

	c := newTestConfig()
	c.Password = "123"
	s, err := NewServer(c, nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	b := response.Body.String()
	is.Equal(response.Code, 403, b)
	is.Contains(b, "403: Password Required")
}

func TestAuthWrongPass(t *testing.T) {
	is := require.New(t)
	data := url.Values{}
	data.Set("password", "1")
	request, err := http.NewRequest(
		http.MethodPost, "/auth/", strings.NewReader(data.Encode()),
	)
	is.Nil(err)
	response := httptest.NewRecorder()

	c := newTestConfig()
	c.Password = "123"
	s, err := NewServer(c, nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	b := response.Body.String()
	is.Equal(response.Code, 403, b)
	is.Contains(b, "403: Password Required")
}

func TestAuthCorrectPass(t *testing.T) {
	is := require.New(t)
	data := url.Values{}
	data.Set("password", "123")
	request, err := http.NewRequest(
		http.MethodPost, "/auth/", strings.NewReader(data.Encode()),
	)
	is.Nil(err)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	response := httptest.NewRecorder()

	c := newTestConfig()
	c.Password = "123"
	s, err := NewServer(c, nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	is.Equal(response.Code, 303)
	is.Equal(response.Header()["Location"], []string{"/p/"})
}
