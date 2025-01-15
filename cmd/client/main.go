package main

import (
	"beliaev-aa/GophKeeper/internal/client/config"
	"beliaev-aa/GophKeeper/internal/client/tui/app"
	"beliaev-aa/GophKeeper/pkg/utils"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logger := utils.NewLogger()
	logger = utils.AddLoggerFields(logger, "GophKeeper Client")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	cfg.BuildVersion = buildVersion
	cfg.BuildDate = buildDate
	cfg.BuildCommit = buildCommit

	if err != nil {
		logger.Fatal("Error initializing model", zap.Error(err))
	}
	tuiApp := app.NewTuiApplication(cfg, logger)
	tuiApp.Start()
}
