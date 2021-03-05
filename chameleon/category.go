package chameleon

const (
	Extension = ".md"
	ReadMe    = "README.md"
)

type Category struct {
	Repository Repository
	Path       Path
}

func (c Category) HasReadme() (bool, error) {
	p := c.Path.Join(ReadMe)
	return p.IsFile()
}

func (c Category) Title() (string, error) {
	a := Article{
		Repository: c.Repository,
		Path:       c.Path.Join(ReadMe),
	}
	return a.Title()
}
