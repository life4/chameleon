package main

import (
	"bytes"
	"text/template"
)

// Category is a struct with all information about content in given category
type Category struct {
	Repo, Branch, Dir, Name, Ext string
}

func (category *Category) makeLink(linkTemplate string) (string, error) {
	t, err := template.New("linkTemplate").Parse(linkTemplate)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	err = t.Execute(buffer, *category)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
