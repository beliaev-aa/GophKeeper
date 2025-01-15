package main

import (
	"beliaev-aa/GophKeeper/internal/server/config"
	"beliaev-aa/GophKeeper/internal/server/grpc"
	"beliaev-aa/GophKeeper/internal/server/storage"
	"beliaev-aa/GophKeeper/pkg/utils"
	"go.uber.org/zap"
)

func main() {
	logger := utils.NewLogger()
	logger = utils.AddLoggerFields(logger, "GophKeeper Server")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	store := &storage.Storage{}
	store, err = storage.NewStorage(cfg.PostgresDSN)
	if err != nil {
		logger.Fatal("Database error", zap.Error(err))
	}

	server := grpc.NewServer(cfg, store, logger)
	err = server.Start()
	if err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
