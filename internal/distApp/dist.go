package distapp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	clusterapp "github.com/teatou/distributor/internal/clusterApp"
	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/internal/handlers/wait"
	"github.com/teatou/distributor/pkg/mylogger"
)

type Distributor struct {
	srv     *http.Server
	Cluster clusterapp.Cluster
}

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

	r := chi.NewRouter()

	// bs

	d := &Distributor{}

	r.Get("/wait", wait.New(d, logger))
	//

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", cfg.Dist.Port),
		Handler: r,
	}

	logger.Infof("starting server on port: %d", cfg.Dist.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Errorf("failed to start server")
	}

	logger.Errorf("server stopped")

	return nil
}
