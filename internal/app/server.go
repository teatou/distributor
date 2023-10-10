package app

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL               *url.URL   // URL of the backend server.
	ActiveConnections int        // Count of active connections
	Mutex             sync.Mutex // A mutex for safe concurrency
	Healthy           bool
}

func (s *Server) Proxy() *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(s.URL)
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
