package chameleon

import (
	"path"
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

func (c Category) HasReadme() (bool, error) {
	p := c.Path().Join(ReadMe)
	return p.IsFile()
}

func (c Category) Title() (string, error) {
	a := Article{
		Repository: c.Repository,
		FileName:   path.Join(c.DirName, ReadMe),
	}
	return a.Title()
}
