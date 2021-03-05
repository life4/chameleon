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
