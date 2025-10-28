package main

import (
	"log"
	"net"
	"os"

	. "chitchat/m/grpc"
	. "chitchat/m/server/connection"
	. "chitchat/m/server/service"

	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "[]: ", 0)
	defer logger.Printf("Server shut down.")
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	logger.Printf("Listening on %v\n", lis.Addr())
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	service := ChitChatService{
		Connections: []*Connection{},
		Channel:     make(chan *Message),
		Logger:      logger,
	}
	RegisterChitChatServer(grpcServer, &service)
	go service.ManageChannels()
	err = grpcServer.Serve(lis)
	if err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
