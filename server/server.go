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
	connections []Connection
}

type Connection struct {
	conn    grpc.BidiStreamingServer[Message, Message]
	channel chan string
}

func (s *ChitChatService) RouteChat(server grpc.BidiStreamingServer[Message, Message]) error {
	conn := Connection{server, make(chan string)}
	s.connections = append(s.connections, conn)
	go conn.Listen()
	for {
	}
	return nil
}

func (c *Connection) Listen() {
	for {
		recv, err := c.conn.Recv()
		if err != nil {
			return
		}
		log.Print("Msg: " + recv.GetText())
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Listening on %v\n", lis.Addr())
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	service := ChitChatService{}
	RegisterChitChatServer(grpcServer, &service)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
