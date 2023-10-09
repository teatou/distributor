package main

import (
	"os"

	serverpool "github.com/teatou/distributor/internal/app/serverPool"
	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/pkg/mylogger"
)

const configEnv = "CONFIG"

func main() {
	val, ok := os.LookupEnv(configEnv)
	if !ok {
		panic("no config env")
	}

	cfg, err := config.LoadConfig(val)
	if err != nil {
		panic("uploading config error")
	}

	logger, err := mylogger.NewZapLogger(cfg.Logger.Level)
	if err != nil {
		panic("making mylogger error")
	}
	defer logger.Sync()

	if err = serverpool.New(cfg, logger); err != nil {
		panic(err)
	}
}
