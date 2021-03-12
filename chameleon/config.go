package chameleon

import (
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	Address  string
	Pull     time.Duration
	RepoPath string
	RepoURL  string
	DBPath   string
	PProf    bool
}

func NewConfig() Config {
	return Config{
		Address:  "127.0.0.1:1337",
		Pull:     5 * time.Minute,
		RepoPath: ".repo",
		RepoURL:  "https://github.com/orsinium/notes.git",
		DBPath:   ".database.bin",
	}
}

func (c Config) Parse() Config {
	pflag.StringVar(&c.RepoPath, "path", c.RepoPath, "path to the repository")
	pflag.StringVar(&c.RepoURL, "url", c.RepoURL, "clone URL for repo if not exist")
	pflag.DurationVar(&c.Pull, "pull", c.Pull, "how often to pull the repo")
	pflag.StringVar(&c.Address, "addr", c.Address, "address to serve")
	pflag.StringVar(&c.DBPath, "db", c.DBPath, "path to database file")
	pflag.BoolVar(&c.PProf, "pprof", c.PProf, "serve pprof endpoints")
	pflag.Parse()
	return c
}
