package main

import (
	"net"
	"personal-http-server/balancer"
	"personal-http-server/server"
)

func main() {
	//spin up the load balancer
	go balancer.Run()
	lc := net.ListenConfig{}
	server := server.NewServer("Test Server", "localhost", "5173", lc)
	defer server.SrvCancel()

	server.StartServer()

	//blocks the main routine until the server is finished
	server.WaitForShutdown()

}
