package chameleon

import "path"

const Extension = ".md"

type Category struct {
	Repository Repository
	Slug       string
	DirName    string
	Branch     string
	Name       string
}

func (c Category) Path() string {
	return path.Join(c.Repository.Path(), c.DirName)
}
