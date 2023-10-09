package target

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Target struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

func (t *Target) SetAlive(alive bool) {
	t.mux.Lock()
	t.Alive = alive
	t.mux.Unlock()
}

// IsAlive returns true when target is alive
func (t *Target) IsAlive() (alive bool) {
	t.mux.RLock()
	alive = t.Alive
	t.mux.RUnlock()
	return
}
