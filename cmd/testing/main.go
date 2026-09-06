package main

import (
	"net"
	"personal-http-server/balancer"
	"personal-http-server/server"
	"strconv"
	"sync"
)

const NUMBACKENDS int = 10
const STARTPORT int = 5173

func main() {
	//spin up the load balancer
	serverPool := balancer.ServerPool{
		Backends: []*balancer.Backend{},
		Current:  0,
	}
	var wg sync.WaitGroup

	for i := range NUMBACKENDS {

		lc := net.ListenConfig{}
		wg.Add(1)

		wg.Go(func() {
			defer wg.Done()
			port := STARTPORT + i
			server := server.NewServer("Test Server", "localhost", strconv.Itoa(port), lc)
			defer server.SrvCancel()

			go server.StartServer()
			server.WaitForShutdown()
		})

		serverPool.AddBackend(balancer.NewBackend("http://localhost:5173"))
	}
	go balancer.Run(&serverPool)
	wg.Wait()

}
