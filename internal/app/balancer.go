package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/pkg/mylogger"
)

func New(cfg *config.Config, logger mylogger.Logger) error {
	var servers []*Server
	for _, port := range cfg.Cluster.Ports {
		s := NewServer(port)
		servers = append(servers, s)
		s.Start()
	}

	// for _, server := range servers {
	// 	go func(s *Server) {
	// 		for range time.Tick(time.Second * 10) {
	// 			res, err := http.Get(s.URL.String())
	// 			if err != nil || res.StatusCode >= 500 {
	// 				s.Healthy = false
	// 			} else {
	// 				s.Healthy = true
	// 			}
	// 		}
	// 	}(server)
	// }

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

	go healthCheck(servers)
	go func() {
		time.Sleep(12 * time.Second)
		servers[0].SetUnhealthy()
	}()
	go func() {
		time.Sleep(12 * time.Second)
		servers[0].SetHealthy()
	}()

	log.Println("Starting server on port", cfg.Balancer.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Balancer.Port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}

	return nil
}

func healthCheck(servers []*Server) {
	t := time.NewTicker(time.Second * 5)
	for {
		<-t.C
		log.Println("Starting health check...")
		for _, s := range servers {
			log.Printf("server %d: %d conns, healthy: %v", s.Port, s.ActiveConnections, s.Healthy)
		}
		log.Println("Health check completed")
	}

}
