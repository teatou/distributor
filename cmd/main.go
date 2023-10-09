package main

import (
	"os"

	clusterapp "github.com/teatou/distributor/internal/clusterApp"
	"github.com/teatou/distributor/internal/config"
	distapp "github.com/teatou/distributor/internal/distApp"
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

	// for _, port := range cfg.Cluster.Ports {
	// 	serverUrl, err := url.Parse(fmt.Sprintf("localhost:%d", port))
	// 	if err != nil {
	// 		panic("cannot parse config ports")
	// 	}
	// 	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	// 	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
	// 		log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
	// 		retries := GetRetryFromContext(request)
	// 		if retries < 3 {
	// 			select {
	// 			case <-time.After(10 * time.Millisecond):
	// 				ctx := context.WithValue(request.Context(), Retry, retries+1)
	// 				proxy.ServeHTTP(writer, request.WithContext(ctx))
	// 			}
	// 			return
	// 		}

	// 		// after 3 retries, mark this backend as down
	// 		serverPool.MarkBackendStatus(serverUrl, false)

	// 		// if the same request routing for few attempts with different backends, increase the count
	// 		attempts := GetAttemptsFromContext(request)
	// 		log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
	// 		ctx := context.WithValue(request.Context(), Attempts, attempts+1)
	// 		lb(writer, request.WithContext(ctx))
	// 	}

	// 	serverPool.AddBackend(&Backend{
	// 		URL:          serverUrl,
	// 		Alive:        true,
	// 		ReverseProxy: proxy,
	// 	})
	// 	log.Printf("Configured server: %s\n", serverUrl)
	// }

	if err := distapp.New(val); err != nil {
		panic(err)
	}

	err = clusterapp.New(val)
	if err != nil {
		panic(err)
	}
}
