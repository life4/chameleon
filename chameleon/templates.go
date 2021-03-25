package chameleon

import (
	"crypto/md5"
	"embed"
	"fmt"
	"text/template"
	"time"
)

//go:embed templates/*.html.j2
var templates embed.FS

var (
	TemplateArticle = parseTemplate("templates/category.html.j2")
	TemplateLinter  = parseTemplate("templates/linter.html.j2")
	TemplateCommits = parseTemplate("templates/commits.html.j2")
	TemplateAuth    = parseTemplate("templates/auth.html.j2")
	TemplateDiff    = parseTemplate("templates/diff.html.j2")
	TemplateStat    = parseTemplate("templates/stat.html.j2")
	TemplateSearch  = parseTemplate("templates/search.html.j2")

	Template403 = parseTemplate("templates/403.html.j2")
	Template404 = parseTemplate("templates/404.html.j2")
)

func parseTemplate(tname string) *template.Template {
	return template.Must(template.New("base.html.j2").Funcs(funcs).ParseFS(
		templates,
		"templates/base.html.j2",
		string(tname),
	))
}

var funcs = template.FuncMap{
	"date": func(item time.Time) string {
		return item.Format("2006-01-02")
	},
	"gravatar": func(mail string) string {
		hash := md5.Sum([]byte(mail))
		return fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash)
	},
	"percent": func(a, b uint32) uint32 {
		return a * 100 / b
	},
}
