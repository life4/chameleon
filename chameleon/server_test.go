package chameleon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestConfig() Config {
	config := NewConfig()
	config.Pull = 0
	config.DBPath = ""
	config.RepoPath = "../.repo"
	return config
}

func TestRootRedirect(t *testing.T) {
	is := require.New(t)
	request, err := http.NewRequest(http.MethodGet, "/", nil)
	is.Nil(err)
	response := httptest.NewRecorder()

	s, err := NewServer(newTestConfig(), nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	is.Equal(response.Code, 307)
	is.Equal(response.Header()["Location"], []string{"/p/"})
}
