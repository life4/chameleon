package chameleon

import (
	"path"
	"strings"
)

const (
	Extension = ".md"
	ReadMe    = "README.md"
)

type Category struct {
	Repository Repository
	DirName    string
}

func (c Category) Path() Path {
	return c.Repository.Path().Join(c.DirName)
}

func (c Category) Title() (string, error) {
	p := c.Path().Join(ReadMe)
	isfile, err := p.IsFile()
	if err != nil {
		return "", err
	}
	if !isfile {
		return c.DirName, nil
	}
	a := Article{
		Repository: c.Repository,
		FileName:   path.Join(c.DirName, ReadMe),
	}
	t, err := a.Title()
	if err != nil {
		return "", err
	}
	t = strings.TrimSuffix(t, "/"+ReadMe)
	return t, nil
}
