package chameleon

import "fmt"

type Category struct {
	Repository Repository
	Path       Path
}

func (c Category) Valid() (bool, error) {
	isDir, err := c.Path.IsDir()
	if err != nil {
		return false, err
	}
	if !isDir {
		return false, nil
	}
	return c.Path.Join(ReadMe).IsFile()
}

func (c Category) Article() *Article {
	return &Article{
		Repository: c.Repository,
		Path:       c.Path.Join(ReadMe),
	}
}

func (c Category) Categories() ([]Category, error) {
	cats := make([]Category, 0)
	paths, err := c.Path.SubPaths()
	if err != nil {
		return nil, fmt.Errorf("get subpaths for %s: %v", c.Path.Name(), err)
	}
	for _, p := range paths {
		cat := Category{
			Repository: c.Repository,
			Path:       p,
		}
		valid, err := cat.Valid()
		if err != nil {
			return nil, fmt.Errorf("validate category %s: %v", p.Name(), err)
		}
		if !valid {
			continue
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (c Category) Articles() (Articles, error) {
	arts := make(Articles, 0)
	paths, err := c.Path.SubPaths()
	if err != nil {
		return nil, fmt.Errorf("get subpaths for %s: %v", c.Path.Name(), err)
	}
	for _, p := range paths {
		art := &Article{
			Repository: c.Repository,
			Path:       p,
		}
		valid, err := art.Valid()
		if err != nil {
			return nil, fmt.Errorf("validate article %s: %v", p.Name(), err)
		}
		if !valid {
			continue
		}
		if art.IsReadme() {
			continue
		}
		arts = append(arts, art)
	}
	arts.Sort()
	return arts, nil
}

func (c Category) URLs() URLs {
	return URLs(c)
}
