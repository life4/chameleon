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
	return s + "/"
}

func (urls URLs) Article() string {
	return ArticlePrefix + urls.suffix()
}

func (urls URLs) Linter() string {
	return LinterPrefix + urls.suffix()
}

func (urls URLs) Commits() string {
	return CommitsPrefix + urls.suffix()
}
