package app

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	URL               *url.URL   // URL of the backend server.
	ActiveConnections int        // Count of active connections
	Mutex             sync.Mutex // A mutex for safe concurrency
	Healthy           bool
	Srv               *http.ServeMux
	Port              int
}

func NewServer(port int) *Server {
	u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))

	srv := http.NewServeMux()
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		log.Printf("completed %d", port)
	})

	return &Server{
		URL:     u,
		Srv:     srv,
		Port:    port,
		Healthy: true,
	}
}

func (s *Server) Start() {
	go func() {
		log.Printf("Server started on port: %v", s.Port)
		http.ListenAndServe(fmt.Sprintf(":%v", s.Port), s.Srv)
	}()
}

func (s *Server) Proxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
}

func (s *Server) SetUnhealthy() {
	s.Healthy = false
}

func (s *Server) SetHealthy() {
	s.Healthy = true
}

func nextServerLeastActive(servers []*Server) *Server {
	leastActiveConnections := -1
	leastActiveServer := servers[0]
	for _, server := range servers {
		server.Mutex.Lock()
		if (server.ActiveConnections < leastActiveConnections || leastActiveConnections == -1) && server.Healthy {
			leastActiveConnections = server.ActiveConnections
			leastActiveServer = server
		}
		server.Mutex.Unlock()
	}

	return leastActiveServer
}
