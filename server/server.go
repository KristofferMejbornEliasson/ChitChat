package main

import (
	"fmt"
	"log"
	"net"

	. "chitchat/m/grpc"

	"google.golang.org/grpc"
)

type ChitChatService struct {
	UnimplementedChitChatServer
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %v\n", lis.Addr())
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterChitChatServer(grpcServer, &ChitChatService{})
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
