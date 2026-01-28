package setup

import (
	"log"
	"rapid-bridge/domain/port"
	"rapid-bridge/internal/adapter/cache"
	"rapid-bridge/internal/adapter/config"
	keymanagementfs "rapid-bridge/internal/adapter/keymanagement_fs"
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

	keyLoader := keymanagementfs.NewFSKeyLoader()

	keyCacheAdapter := cache.NewKeyCacheAdapter(keyLoader, logger)

	return &Application{
		Config:   cfg,
		Logger:   logger,
		KeyCache: keyCacheAdapter,
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

	keyLoader := keymanagementfs.NewFSKeyLoader()

	keyCache := cache.NewKeyCacheAdapter(keyLoader, logger)

	return &CLIApplication{
		Config:   cfg,
		Logger:   logger,
		KeyCache: keyCache,
	}
}
