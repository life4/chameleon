package main

import (
	"fmt"
	"os"

	"github.com/orsinium/chameleon/chameleon"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func run(logger *zap.Logger) error {
	var repoPath string
	pflag.StringVar(&repoPath, "--path", ".repo", "path to the repository")
	pflag.Parse()

	repo := chameleon.Repository{Path: chameleon.Path(repoPath)}
	s := chameleon.Server{
		Repository: repo,
	}

	logger.Info("initializing repos")
	err := s.Init()
	if err != nil {
		return fmt.Errorf("cannot init repos: %v", err)
	}
	defer s.Close()

	logger.Info("listening")
	return s.Serve()
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()

	err = run(logger)
	if err != nil {
		logger.Error("fatal error", zap.Error(err))
		os.Exit(1)
	}
}
