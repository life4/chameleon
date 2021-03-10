package chameleon

import (
	"crypto/md5"
	"embed"
	"fmt"
	"reflect"
	"text/template"
	"time"
)

//go:embed templates/*.html.j2
var templates embed.FS

var (
	TemplateArticle = parseTemplate("templates/category.html.j2")
	TemplateLinter  = parseTemplate("templates/linter.html.j2")
	TemplateCommits = parseTemplate("templates/commits.html.j2")
	Template404     = parseTemplate("templates/404.html.j2")
)

func parseTemplate(tname string) *template.Template {
	return template.Must(template.New("base.html.j2").Funcs(funcs).ParseFS(
		templates,
		"templates/base.html.j2",
		string(tname),
	))
}

var funcs = template.FuncMap{
	"first": func(item reflect.Value) (reflect.Value, error) {
		item, isNil := indirect(item)
		if isNil {
			return reflect.Value{}, fmt.Errorf("index of nil pointer")
		}
		item = item.Index(0)
		return item, nil
	},
	"last": func(item reflect.Value) (reflect.Value, error) {
		item, isNil := indirect(item)
		if isNil {
			return reflect.Value{}, fmt.Errorf("index of nil pointer")
		}
		item = item.Index(item.Len() - 1)
		return item, nil
	},
	"date": func(item time.Time) string {
		return item.Format("2006-01-02")
	},
	"gravatar": func(mail string) string {
		hash := md5.Sum([]byte(mail))
		return fmt.Sprintf("https://www.gravatar.com/avatar/%x", hash)
	},
}

func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}
