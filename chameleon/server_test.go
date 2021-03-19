package chameleon

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootRedirect(t *testing.T) {
	is := require.New(t)
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	config := NewConfig()
	config.Pull = 0
	config.DBPath = ""
	config.RepoPath = "../.repo"
	s, err := NewServer(config, nil)
	is.Nil(err)
	s.ServeHTTP(response, request)

	is.Equal(response.Code, 307)
	is.Equal(response.Header()["Location"], []string{"/p/"})
}
