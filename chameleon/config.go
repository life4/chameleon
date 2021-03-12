package chameleon

import (
	"time"

	"github.com/spf13/pflag"
)

type Config struct {
	Address  string
	RepoPath string
	RepoURL  string
	Pull     time.Duration
	Cache    int
	DBPath   string
	PProf    bool
}

func NewConfig() Config {
	return Config{
		Address:  "127.0.0.1:1337",
		RepoPath: ".repo",
		RepoURL:  "https://github.com/orsinium/notes.git",
		Pull:     5 * time.Minute,
		Cache:    1000,
		DBPath:   ".database.bin",
	}
}

func (c Config) Parse() Config {
	pflag.StringVar(&c.Address, "addr", c.Address, "address to serve")
	pflag.StringVar(&c.RepoPath, "path", c.RepoPath, "path to repository")
	pflag.StringVar(&c.RepoURL, "url", c.RepoURL, "clone URL for repo if not exist")
	pflag.DurationVar(&c.Pull, "pull", c.Pull, "how often pull repository, 0 to disable")
	pflag.IntVar(&c.Cache, "cache", c.Cache, "how many records to cache, 0 to disable")
	pflag.StringVar(&c.DBPath, "db", c.DBPath, "path to database file")
	pflag.BoolVar(&c.PProf, "pprof", c.PProf, "serve pprof endpoints")
	pflag.Parse()
	return c
}
