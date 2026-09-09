package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"personal-http-server/client"
	"sync"
)

type Server struct {
	Port           string
	Host           string
	Addr           string
	TCPAddr        *net.TCPAddr
	Listener       net.Listener
	ListenerConfig net.ListenConfig
	ctx            context.Context
	srvWG          *sync.WaitGroup
	SrvCancel      context.CancelFunc
	mu             sync.Mutex
	done           chan bool
}

func NewServer(serverName, host, port string, lc net.ListenConfig) *Server {
	addr := net.JoinHostPort(host, port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		Port:           port,
		Host:           host,
		Addr:           addr,
		TCPAddr:        tcpAddr,
		ctx:            ctx,
		SrvCancel:      cancel,
		ListenerConfig: lc,
	}
	return server
}
func (s *Server) StartServer() error {
	log.Printf("Will be listening on %s -- when we run our listeners", s.Addr)

	s.AcceptListener()
	return nil
}

func (s *Server) AcceptListener() {

	ctx, cancel := context.WithCancel(s.ctx)

	defer cancel()

	l, err := s.ListenerConfig.Listen(ctx, s.TCPAddr.Network(), s.TCPAddr.String())
	if err != nil {
		return // child that holds this context will be shutdown
	}

	s.Listener = l

	log.Printf("Listening on %s\n", l.Addr())

	//Accept the client connection and parse
	go s.acceptConnection(l, "Client", func(conn net.Conn) {
		//We can process the connection here as a client struct
		clientCreated := client.NewClient(conn)
		clientCreated.ReadRequest()
		clientCreated.WriteResponse()
		defer conn.Close()
	}, func(err error) bool {
		log.Println("Error handling a connection ")
		s.SrvCancel() //cancels any routine that uses this context
		return true

	})
	fmt.Println("Go routine for accepting connection has been fired")

}

func (s *Server) acceptConnection(l net.Listener, acceptName string, createFunc func(conn net.Conn), errFunc func(err error) bool) {

	for {
		conn, err := l.Accept()
		if err != nil {
			if errFunc != nil && errFunc(err) {
				log.Println("custom error called")
				break
			}
			break
		}

		log.Println("We got a connection")
		go createFunc(conn)
	}

	//Notifies that the server is done accepting connections
	s.done <- true
}

func (s *Server) IsDone() {
	<-s.done
}
func (s *Server) WaitForShutdown() {

	log.Println("waiting for shutdown")
	select {
	case <-s.done:
		log.Println("Server shutdown through done")
	case <-s.ctx.Done():
		fmt.Println("Server shutdown through context done")
	}

}
