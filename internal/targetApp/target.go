package targetapp

import (
	"fmt"

	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/pkg/mylogger"
)

func New(configFileName string) error {
	cfg, err := config.LoadConfig(configFileName)
	if err != nil {
		return fmt.Errorf("uploading config error: %w", err)
	}

	logger, err := mylogger.NewZapLogger(cfg.Logger.Level)
	if err != nil {
		return fmt.Errorf("making mylogger error: %w", err)
	}
	defer logger.Sync()

	return nil
}
