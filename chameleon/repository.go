package chameleon

import (
	"fmt"
	"os/exec"
)

const CloneURL = "https://github.com/%s.git"

type Repository struct {
	Slug    string
	Storage Storage
}

func (r Repository) Path() string {
	return r.Storage.Join(r.Slug)
}

func (r Repository) URL() string {
	return fmt.Sprintf(CloneURL, r.Slug)
}

func (r Repository) Command(args ...string) *exec.Cmd {
	c := exec.Command("git", args...)
	c.Dir = r.Path()
	return c
}

func (r Repository) Clone() error {
	err := r.Storage.Ensure()
	if err != nil {
		return fmt.Errorf("cannot ensure storage dir: %v", err)
	}

	c := exec.Command("git", "clone", r.URL())
	c.Dir = r.Storage.Path()
	err = c.Run()
	if err != nil {
		return fmt.Errorf("cannot clone repo: %v", err)
	}
	return nil
}
