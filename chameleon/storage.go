package chameleon

import (
	"fmt"
	"os"
	"path"
)

type Storage string

func (s Storage) Path() string {
	return string(s)
}

func (s Storage) Join(fname string) string {
	return path.Join(s.Path(), fname)
}

func (s Storage) Exists() (bool, error) {
	_, err := os.Stat(s.Path())
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s Storage) Ensure() error {
	exists, err := s.Exists()
	if err != nil {
		return fmt.Errorf("cannot get dir stat: %v", err)
	}
	if exists {
		return nil
	}
	err = os.MkdirAll(s.Path(), 0x770)
	if err != nil {
		return fmt.Errorf("cannot create dir: %v", err)
	}
	return nil
}
