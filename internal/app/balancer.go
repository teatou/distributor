package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/pkg/mylogger"
)

func New(cfg *config.Config, logger mylogger.Logger) error {
	var servers []*Server
	for _, serverUrl := range cfg.Cluster.Ports {
		u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", serverUrl))
		servers = append(servers, &Server{URL: u})
	}

	for _, server := range servers {
		go func(s *Server) {
			for range time.Tick(time.Second * 10) {
				res, err := http.Get(s.URL.String())
				if err != nil || res.StatusCode >= 500 {
					s.Healthy = false
				} else {
					s.Healthy = true
				}
			}
		}(server)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := nextServerLeastActive(servers)
		server.Mutex.Lock()
		server.ActiveConnections++
		server.Mutex.Unlock()
		server.Proxy().ServeHTTP(w, r)
		server.Mutex.Lock()
		server.ActiveConnections--
		server.Mutex.Unlock()
	})

	log.Println("Starting server on port", cfg.Balancer.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Balancer.Port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

	return nil
}
