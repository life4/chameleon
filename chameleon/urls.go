package chameleon

import (
	"strings"
)

type URLs struct {
	Repository Repository
	Path       Path
}

func (urls URLs) suffix() string {
	s := urls.Path.Relative(urls.Repository.Path).String()
	s = strings.TrimSuffix(s, Extension)
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
	return MainPrefix + strings.TrimSuffix(urls.suffix(), "/") + Extension
}
