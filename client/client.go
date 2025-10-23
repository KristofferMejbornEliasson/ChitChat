package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	. "chitchat/m/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	txt := "This is a message from a client! :)"
	msg := &Message{}
	msg.Text = &txt
	err = stream.Send(msg)
	if err != nil {
		log.Fatalf("fail to call Send: %v", err)
	}
	wait := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}
			if err != nil {
				log.Fatalf("Connection was closed.")
			}
			log.Println(in.GetText())
		}
	}()
	reader := bufio.NewScanner(os.Stdin)
	for {
		reader.Scan()
		if reader.Err() != nil {
			log.Fatalf("fail to call Read: %v", err)
		}
		text := reader.Text()
		if text == "exit" {
			break
		}
		err = stream.Send(&Message{Text: &text})
		if err != nil {
			log.Fatalf("fail to call Send: %v", err)
		}
	}
	_ = stream.CloseSend()
	<-wait
}
