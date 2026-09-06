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
	Backends []*Backend
	Current  int
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
func (s *ServerPool) AddBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func (s *ServerPool) NextIdx() int {
	return ((s.Current + (1)) % (len(s.Backends)))
}

func (s *ServerPool) GetNextPeer() *Backend {
	nextIdx := s.NextIdx()

	for {
		fmt.Printf(" nextIdx %d\n", nextIdx)
		nextIdx %= len(s.Backends)

		if !s.Backends[nextIdx].IsAlive() {
			nextIdx += 1
			continue
		}

		s.Current += nextIdx
		return s.Backends[nextIdx]
	}
}

func Run(serverPool *ServerPool) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		fmt.Println("Load balancer signal interrupt")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
	}()
}
