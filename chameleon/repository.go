package chameleon

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	giturls "github.com/whilp/git-urls"
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

func (r Repository) Remote() (*url.URL, error) {
	c := r.Command("remote", "get-url", "origin")
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}
	rawURL := strings.TrimSpace(string(out))
	return giturls.Parse(rawURL)
}

func (r Repository) Branch() (string, error) {
	c := r.Command("rev-parse", "--abbrev-ref", "HEAD")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(string(out))
	return branch, nil
}

func (r Repository) Clone(url string) error {
	isdir, err := r.Path.IsDir()
	if err != nil {
		return fmt.Errorf("cannot access local repo: %v", err)
	}
	if isdir {
		return nil
	}
	if url == "" {
		return fmt.Errorf("repo not found and no clone URL specified")
	}

	c := exec.Command("git", "clone", url, r.Path.String())
	out, err := c.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot clone repo: %v: %s", err, out)
	}
	return nil
}
