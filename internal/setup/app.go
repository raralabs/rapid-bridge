package setup

import (
	"log"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/adapter/config"
	"rapid-bridge/internal/adapter/logger"
)

type Application struct {
	Config port.ServerConfig
	Logger port.Logger
}

type CLIApplication struct {
	Config port.CLIConfig
	Logger port.Logger
}

func NewApplication() *Application {
	logger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &Application{
		Config: cfg,
		Logger: logger,
	}
}

func NewCLIApplication() *CLIApplication {
	logger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg, err := config.LoadCLIConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &CLIApplication{
		Config: cfg,
		Logger: logger,
	}
}
