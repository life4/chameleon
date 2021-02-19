package chameleon

import "fmt"

type Pages struct {
	Repository Repository
}

func (ps Pages) Page(urlPath string) (*Page, error) {
	if urlPath == "" || urlPath == "/" {
		page := Page{
			Repository: &ps.Repository,
			Article: &Article{
				Repository: ps.Repository,
				FileName:   string(ps.Repository.Path().Join(ReadMe)),
			},
		}
		return &page, nil
	}

	p := ps.Repository.Storage.Join(urlPath)

	isdir, err := p.IsDir()
	if err != nil {
		return nil, err
	}
	if isdir {
		isfile, err := p.Join(ReadMe).IsFile()
		if err != nil {
			return nil, err
		}
		if !isfile {
			return nil, fmt.Errorf("README.md not found")
		}
		page := Page{
			Category: &Category{
				Repository: ps.Repository,
				DirName:    string(p.Relative(ps.Repository.Path())),
			},
			Article: &Article{
				Repository: ps.Repository,
				FileName:   string(p.Join(ReadMe).Relative(ps.Repository.Path())),
			},
			Repository: &ps.Repository,
		}
		return &page, nil
	}

	isfile, err := p.IsFile()
	if err != nil {
		return nil, err
	}
	if isfile {
		page := Page{
			Category: &Category{
				Repository: ps.Repository,
				DirName:    string(p.Parent()),
			},
			Article: &Article{
				Repository: ps.Repository,
				FileName:   string(p.Relative(ps.Repository.Path())),
			},
			Repository: &ps.Repository,
		}
		return &page, nil
	}

	return nil, fmt.Errorf("file not found")
}
