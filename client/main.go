package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	. "chitchat/m/client/client"
	. "chitchat/m/grpc"
)

func main() {
	file, err := os.CreateTemp("", "client-*.txt")
	logger := log.New(file, "", 0)
	var c *Client
	defer func(c *Client) {
		if c != nil {
			c.Log("Shutting down client.")
		} else {
			logger.Print("Shutting down client")
		}
		_ = file.Close()
	}(c)

	// Create a gRPC channel
	// If credentials are needed, these are set in the options.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		logger.Fatalf("Failed to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			logger.Fatalf("Error closing connection:\n%v", err)
		}
	}(conn)

	client := NewChitChatClient(conn)
	stream, err := client.RouteChat(context.Background())
	if err != nil {
		logger.Fatalf("Failed to create a stream via call to RouteChat: %v", err)
	}

	txt := "Establishing connection."
	msg := &Message{
		Text: &txt,
	}
	err = stream.Send(msg)
	if err != nil {
		logger.Fatalf("Failed to send establishing message: %v", err)
	}

	c = NewClient(stream, logger)
	c.Run()
}
