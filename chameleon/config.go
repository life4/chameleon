package chameleon

import (
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	Address  string
	Pull     time.Duration
	RepoPath string
}

func NewConfig() Config {
	return Config{
		Address:  "127.0.0.1:1337",
		Pull:     5 * time.Minute,
		RepoPath: ".repo",
	}
}

func (c Config) Parse() Config {
	pflag.StringVar(&c.RepoPath, "--path", c.RepoPath, "path to the repository")
	pflag.DurationVar(&c.Pull, "--pull", c.Pull, "how often to pull the repo")
	pflag.StringVar(&c.Address, "--addr", c.Address, "address to serve")
	pflag.Parse()
	return c
}
