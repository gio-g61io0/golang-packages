package main

import (
	"fmt"
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
		Backends:         []*balancer.Backend{},
		Current:          0,
		LoadDistribution: make([]int, NUMBACKENDS),
	}
	var wg sync.WaitGroup

	for i := range NUMBACKENDS {

		lc := net.ListenConfig{}
		wg.Add(1)
		port := STARTPORT + i

		server := server.NewServer("Test Server", "localhost", strconv.Itoa(port), lc)

		wg.Go(func() {
			defer wg.Done()
			defer server.SrvCancel()

			go server.StartServer()
			server.WaitForShutdown()
		})

		serverPool.AddBackend(serverPool.NewBackend(fmt.Sprintf("http://localhost:%d", port), server))
	}
	go balancer.Run(&serverPool)
	go balancer.HealthCheck(&serverPool)
	wg.Wait()

}
