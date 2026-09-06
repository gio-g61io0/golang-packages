package balancer

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	mux   sync.RWMutex
	rp    *httputil.ReverseProxy
}
type ServerPool struct {
	backends []*Backend
	current  int
}

func NewBackend(refUrl string) *Backend {
	url, _ := url.Parse(refUrl)
	rp := httputil.NewSingleHostReverseProxy(url)
	return &Backend{
		URL:   url,
		Alive: true,
		mux:   sync.RWMutex{},
		rp:    rp,
	}
}
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()

}
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive

}

func (s *ServerPool) NextIdx() int {
	return ((s.current + (1)) % (len(s.backends)))
}

func (s *ServerPool) GetNextPeer() *Backend {
	nextIdx := s.NextIdx()

	for {
		fmt.Printf(" nextIdx %d\n", nextIdx)
		nextIdx %= len(s.backends)

		if !s.backends[nextIdx].IsAlive() {
			nextIdx += 1
			continue
		}

		s.current += nextIdx
		return s.backends[nextIdx]
	}
}

func Run() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverPool := ServerPool{
		backends: make([]*Backend, 1),
		current:  0,
	}
	serverPool.backends[0] = NewBackend("http://localhost:5173")

	lb := func(w http.ResponseWriter, rq *http.Request) {
		peer := serverPool.GetNextPeer()
		if peer != nil {
			peer.rp.ServeHTTP(w, rq)
			return
		}
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}

	server := http.Server{
		Handler: http.HandlerFunc(lb),
	}

	listenerConfig := net.ListenConfig{}
	addr := net.JoinHostPort("localhost", "6106")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	listener, err := listenerConfig.Listen(ctx, tcpAddr.Network(), tcpAddr.String())

	if err != nil {
		log.Fatal("Cannot initialize listener")
	}

	server.Serve(listener)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
	}()
}
