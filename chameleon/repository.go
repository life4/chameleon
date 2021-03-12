package chameleon

import (
	"fmt"
	"os/exec"
)

const CloneURL = "https://github.com/%s.git"

type Repository struct {
	Path
}

func (r Repository) Command(args ...string) *exec.Cmd {
	c := exec.Command("git", args...)
	c.Dir = r.Path.String()
	return c
}

func (r Repository) Pull() error {
	c := r.Command("pull")
	err := c.Run()
	if err != nil {
		return fmt.Errorf("cannot pull repo: %v", err)
	}
	return nil
}

func (r Repository) Clone(url string) error {
	isdir, err := r.Path.IsDir()
	if err != nil {
		return fmt.Errorf("cannot access local repo: %v", err)
	}
	if isdir {
		return nil
	}

	c := exec.Command("git", "clone", url, r.Path.String())
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot clone repo: %v: %s", err, out)
	}
	return nil
}
