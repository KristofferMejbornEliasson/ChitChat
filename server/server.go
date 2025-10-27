package main

import (
	"fmt"
	"log"
	"net"

	. "chitchat/m/grpc"
	. "chitchat/m/server/connection"
	. "chitchat/m/server/service"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %v\n", lis.Addr())
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	service := ChitChatService{
		Connections: []Connection{},
		Channel:     make(chan *Message),
	}
	RegisterChitChatServer(grpcServer, &service)
	go service.ManageChannels()
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
