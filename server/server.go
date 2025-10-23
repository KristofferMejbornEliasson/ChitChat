package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %v\n", lis.Addr())
	// var opts []grpc.ServerOption
	// grpcServer := grpc.NewServer(opts...)
}
