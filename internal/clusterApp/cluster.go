package clusterapp

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/internal/handlers/quit"
	"github.com/teatou/distributor/pkg/mylogger"
)

type Cluster struct {
	N       int
	Targets []*Target
}

type Target struct {
	srv      *http.Server
	ReqCount int
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

	// questionable
	c := &Cluster{
		N: cfg.Cluster.N,
	}

	r := chi.NewRouter()

	r.Post("/quit", quit.New(c, logger))

	for i := 0; i < c.N; i++ {
		srv := &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", cfg.Cluster.Ports[i]),
			Handler: r,
		}

		t := &Target{
			srv: srv,
		}

		c.Targets = append(c.Targets, t)

		// if i == c.N-1 {
		// 	logger.Infof("last server starts on port, %d", cfg.Cluster.Ports[i])
		// 	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Cluster.Ports[i]), r)
		// 	if err != nil {
		// 		logger.Fatalf("last server is down: %v", err)
		// 	}
		// } else {

		go func() {
			logger.Infof("server starts on port, %d", cfg.Cluster.Ports[i])
			err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Cluster.Ports[i]), r)
			if err != nil {
				logger.Fatalf("server is down: %v", err)
			}
		}()
	}

	return nil
}
