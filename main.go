package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/orsinium/chameleon/chameleon"
	"go.uber.org/zap"
)

func run(logger *zap.Logger) error {
	config := chameleon.NewConfig().Parse()
	server, err := chameleon.NewServer(config, logger)
	if err != nil {
		return err
	}

	defer func() {
		err := server.Database.Close()
		if err != nil {
			logger.Error("cannot close connection", zap.Error(err))
		}
	}()

	logger.Info("listening", zap.String("addr", config.Address))
	return http.ListenAndServe(config.Address, server)
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync() // nolint

	logger.Info("starting...")
	err = run(logger)
	if err != nil {
		logger.Error("fatal error", zap.Error(err))
		os.Exit(1)
	}
}
