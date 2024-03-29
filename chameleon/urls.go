package chameleon

import (
	"fmt"
	"strings"
)

const (
	githubEditURL = "https://github.com/%s/edit/%s/%s"
)

type URLs struct {
	Repository Repository
	Path       Path
}

func (urls URLs) suffix() string {
	s := urls.Path.Relative(urls.Repository.Path).String()
	s = strings.TrimSuffix(s, ExtensionMarkdown)
	s = strings.TrimSuffix(s, ExtensionJupyter)
	if s == "" {
		return s
	}
	return s + "/"
}

func (urls URLs) Main() string {
	return MainPrefix + urls.suffix()
}

func (urls URLs) Linter() string {
	return LinterPrefix + urls.suffix()
}

func (urls URLs) Commits() string {
	return CommitsPrefix + urls.suffix()
}

func (urls URLs) Raw() string {
	s := urls.Path.Relative(urls.Repository.Path).String()
	return MainPrefix + s
}

func (urls URLs) Edit() (string, error) {
	remote, err := urls.Repository.Remote()
	if err != nil {
		return "", nil
	}

	if remote.Hostname() == "github.com" {
		branch, err := urls.Repository.Branch()
		if err != nil {
			return "", err
		}
		repo := strings.TrimSuffix(remote.Path, ".git")
		url := fmt.Sprintf(githubEditURL, repo, branch, urls.suffix())
		return url, nil
	}

	return "", nil
}
