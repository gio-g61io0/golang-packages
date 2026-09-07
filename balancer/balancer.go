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
	"personal-http-server/server"
	"sync"
	"syscall"
	"time"
)

const (
	Attempts int = iota
	Retry
)

type Backend struct {
	URL    *url.URL
	Alive  bool
	mux    sync.RWMutex
	rp     *httputil.ReverseProxy
	server *server.Server
}
type ServerPool struct {
	Backends         []*Backend
	Current          int
	mux              sync.RWMutex
	LoadDistribution []int
}

func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

func GetAttemptFromContext(r *http.Request) int {
	if attempt, ok := r.Context().Value(Attempts).(int); ok {
		return attempt
	}
	return 0
}

func (b *Backend) IsAlive() bool {
	timeout := time.Second * 2
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)

	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}

	defer conn.Close()

	return true

}

func (s *ServerPool) HealthCheck() {
	for _, backend := range s.Backends {
		status := "up"
		alive := backend.IsAlive()
		backend.SetAlive(alive)

		if !alive {
			status = "down"
		}
		log.Printf("%s [%s]\n", backend.URL, status)
	}

}

func HealthCheck(s *ServerPool) {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-t.C:
			log.Println("Starting Health Check")
			s.HealthCheck()
			log.Println("Health Check Done")
		}
	}

}

func (s *ServerPool) MarkBackendStatus(url *url.URL, status bool) bool {
	for backend := range len(s.Backends) {
		if s.Backends[backend].URL == url {
			s.Backends[backend].SetAlive(status)
			return true
		}
	}
	return false //Nothing is set
}

func (s *ServerPool) NewBackend(refUrl string, server *server.Server) *Backend {
	url, _ := url.Parse(refUrl)
	rp := httputil.NewSingleHostReverseProxy(url)

	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, e error) {

		log.Printf("[%s] %s\n", url.Host, e.Error())
		retries := GetRetryFromContext(req)

		if retries < 3 {
			select {
			case <-time.After(time.Millisecond * 10):
				ctx := context.WithValue(req.Context(), Retry, retries+1)
				rp.ServeHTTP(rw, req.WithContext(ctx))
			}
			return
		}

		s.MarkBackendStatus(url, false)

		attempts := GetAttemptFromContext(req)
		log.Printf("%s(%s) Attempting retry %d\n", req.RemoteAddr, req.URL.Path, attempts)
		ctx := context.WithValue(req.Context(), Attempts, attempts+1)
		lb := s.ProvideLb()
		lb(rw, req.WithContext(ctx))
	}

	return &Backend{
		URL:    url,
		Alive:  true,
		mux:    sync.RWMutex{},
		rp:     rp,
		server: server,
	}
}
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()

}
func (s *ServerPool) AddBackend(backend *Backend) {
	s.Backends = append(s.Backends, backend)
}

func (s *ServerPool) NextIdx() int {
	s.mux.Lock()
	defer s.mux.Unlock()
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
		s.Current = nextIdx
		s.LoadDistribution[nextIdx] += 1 //tracks the distribution of load
		return s.Backends[nextIdx]
	}
}

func (s *ServerPool) ProvideLb() func(w http.ResponseWriter, rq *http.Request) {
	return func(w http.ResponseWriter, rq *http.Request) {
		peer := s.GetNextPeer()
		if peer != nil {
			peer.rp.ServeHTTP(w, rq)
			return
		}
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}
}
func Run(serverPool *ServerPool) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lb := serverPool.ProvideLb()

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

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Load balancer signal interrupt")

		//shutdown backend servers managed by the load balancer
		for _, b := range serverPool.Backends {
			b.server.SrvCancel()
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		server.Shutdown(shutdownCtx)
	}()

	server.Serve(listener)

}
