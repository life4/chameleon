package chameleon

import (
	"fmt"
	"os/exec"
)

func CloneRepo(r Repository, url string) error {
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
