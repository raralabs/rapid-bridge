package setup

import (
	"log"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/adapter/cache"
	"rapid-bridge/internal/adapter/config"
	"rapid-bridge/internal/adapter/logger"
)

type Application struct {
	Config   port.ServerConfig
	Logger   port.Logger
	KeyCache port.KeyCache
}

type CLIApplication struct {
	Config   port.CLIConfig
	Logger   port.Logger
	KeyCache port.KeyCache
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

	keyCache := cache.NewKeyCacheAdapter()

	return &Application{
		Config:   cfg,
		Logger:   logger,
		KeyCache: keyCache,
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

	keyCache := cache.NewKeyCacheAdapter()

	return &CLIApplication{
		Config:   cfg,
		Logger:   logger,
		KeyCache: keyCache,
	}
}
