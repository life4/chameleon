package chameleon

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Name() string {
	return path.Base(p.String())
}

func (p Path) Relative(to Path) Path {
	return Path(strings.TrimPrefix(p.String(), to.String()+"/"))
}

func (p Path) Parent() Path {
	return Path(path.Dir(p.String()))
}

func (p Path) Join(fname string) Path {
	if fname == "" {
		return p
	}
	return Path(path.Join(p.String(), fname))
}

func (p Path) IsDir() (bool, error) {
	stat, err := os.Stat(p.String())
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}

func (p Path) IsFile() (bool, error) {
	stat, err := os.Stat(p.String())
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !stat.IsDir(), nil
}

func (p Path) EnsureDir() error {
	exists, err := p.IsDir()
	if err != nil {
		return fmt.Errorf("cannot get dir stat: %v", err)
	}
	if exists {
		return nil
	}
	err = os.MkdirAll(p.String(), 0777)
	if err != nil {
		return fmt.Errorf("cannot create dir: %v", err)
	}
	return nil
}

func (p Path) SubPaths() ([]Path, error) {
	infos, err := ioutil.ReadDir(p.String())
	if err != nil {
		return nil, err
	}

	result := make([]Path, len(infos))
	for i, info := range infos {
		result[i] = p.Join(info.Name())
	}
	return result, nil
}

func (p Path) Open() (*os.File, error) {
	return os.Open(p.String())
}
