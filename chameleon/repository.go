package chameleon

import (
	"fmt"
	"os/exec"
	"strings"
)

const CloneURL = "https://github.com/%s.git"

type Repository struct {
	URL     string
	Storage Path
}

func (r Repository) Slug() string {
	slug := r.URL
	slug = strings.TrimPrefix(slug, "https://github.com/")
	slug = strings.TrimSuffix(slug, ".git")
	return slug
}

func (r Repository) Path() Path {
	return r.Storage.Join(r.Slug())
}

func (r Repository) Command(args ...string) *exec.Cmd {
	c := exec.Command("git", args...)
	c.Dir = r.Path().String()
	return c
}

func (r Repository) Clone() error {
	isdir, err := r.Path().IsDir()
	if err != nil {
		return fmt.Errorf("cannot access local repo: %v", err)
	}
	if isdir {
		return nil
	}

	c := exec.Command("git", "clone", r.URL, r.Slug())
	c.Dir = r.Storage.String()
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot clone repo: %v: %s", err, out)
	}
	return nil
}

func (r Repository) Pull() error {
	c := r.Command("pull")
	err := c.Run()
	if err != nil {
		return fmt.Errorf("cannot pull repo: %v", err)
	}
	return nil
}
