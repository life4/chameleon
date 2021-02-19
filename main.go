package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	"github.com/orsinium/chameleon/chameleon"
	"go.uber.org/zap"
)

//go:embed config.toml
var config string

//go:embed templates/*.html
var templates embed.FS

func run(logger *zap.Logger) error {
	c := chameleon.Config{}
	_, err := toml.Decode(config, &c)
	if err != nil {
		return fmt.Errorf("cannot parse config: %v", err)
	}

	s := chameleon.Server{
		Config:    c,
		Templates: templates,
		Router:    chi.NewRouter(),
	}

	logger.Info("initializing repos")
	err = s.InitRepos()
	if err != nil {
		return fmt.Errorf("cannot init repos: %v", err)
	}

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
