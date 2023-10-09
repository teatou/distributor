package serverpool

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/teatou/distributor/internal/app/target"
	"github.com/teatou/distributor/internal/config"
	"github.com/teatou/distributor/internal/helpers"
	httpres "github.com/teatou/distributor/internal/http-res"
	"github.com/teatou/distributor/pkg/mylogger"
)

type ServerPool struct {
	targets []*target.Target
	current uint64
}

func New(cfg *config.Config, logger mylogger.Logger) error {
	var serverPool ServerPool

	for _, port := range cfg.Cluster.Ports {
		serverUrl, err := url.Parse(fmt.Sprintf("localhost:%d", port))
		if err != nil {
			return fmt.Errorf("cannot parse config ports")
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
			retries := httpres.GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), httpres.Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// after 3 retries, mark this backend as down
			serverPool.MarkTargetStatus(serverUrl, false)

			// if the same request routing for few attempts with different backends, increase the count
			attempts := httpres.GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), httpres.Attempts, attempts+1)
			serverPool.Lb(writer, request.WithContext(ctx))
		}

		serverPool.AddTarget(&target.Target{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Dist.Port),
		Handler: http.HandlerFunc(serverPool.Lb),
	}

	// start health checking
	go serverPool.HealthCheck2Mins()

	log.Printf("Load Balancer started at :%d\n", cfg.Dist.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("server stopped")
	}

	return nil
}

// AddTarget to the server pool
func (s *ServerPool) AddTarget(target *target.Target) {
	s.targets = append(s.targets, target)
}

// NextIndex atomically increase the counter and return an index
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.targets)))
}

// MarkTargetStatus changes a status of a target
func (s *ServerPool) MarkTargetStatus(targetUrl *url.URL, alive bool) {
	for _, t := range s.targets {
		if t.URL.String() == targetUrl.String() {
			t.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer returns next active peer to take a connection
func (s *ServerPool) GetNextPeer() *target.Target {
	// loop entire targets to find out an Alive target
	next := s.NextIndex()
	l := len(s.targets) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(s.targets)     // take an index by modding
		if s.targets[idx].IsAlive() { // if we have an alive target, use it and store if its not the original one
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.targets[idx]
		}
	}
	return nil
}

// HealthCheck pings the targets and update the status
func (s *ServerPool) HealthCheck() {
	for _, t := range s.targets {
		status := "up"
		alive := helpers.IsTargetAlive(t.URL)
		t.SetAlive(alive)
		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", t.URL, status)
	}
}

func (s *ServerPool) HealthCheck2Mins() {
	t := time.NewTicker(time.Minute * 2)
	for {
		<-t.C
		log.Println("Starting health check...")
		s.HealthCheck()
		log.Println("Health check completed")
	}
}

// lb load balances the incoming request
func (s *ServerPool) Lb(w http.ResponseWriter, r *http.Request) {
	attempts := httpres.GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := s.GetNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
