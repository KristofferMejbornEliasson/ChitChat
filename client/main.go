package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "chitchat/m/client/client"
	. "chitchat/m/grpc"
)

func main() {
	// Create a gRPC channel
	// If credentials are needed, these are set in the options.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("error closing connection:\n%v", err)
		}
	}(conn)

	client := NewChitChatClient(conn)
	stream, err := client.RouteChat(context.Background())
	if err != nil {
		log.Fatalf("fail to call RouteChat: %v", err)
	}

	txt := "This is a message from a new client! :)"
	msg := &Message{}
	msg.Text = &txt
	err = stream.Send(msg)
	if err != nil {
		log.Fatalf("fail to call Send: %v", err)
	}

	c := NewClient(stream)
	c.Run()
}
