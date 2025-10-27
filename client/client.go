package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	. "chitchat/m/client/clock"
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
	txt := "This is a message from a new client! :)"
	msg := &Message{}
	msg.Text = &txt
	err = stream.Send(msg)
	if err != nil {
		log.Fatalf("fail to call Send: %v", err)
	}
	in, err := stream.Recv()
	if err != nil {
		log.Fatalf("No initial response received from server:\n%v", err)
	}
	c := newClient(in)
	c.Run(stream)
}

type Client struct {
	id          int32
	vectorClock Clock
}

func (c *Client) IncrementClock() {
	c.vectorClock.Increment(c.id)
}

func newClient(msg *Message) *Client {
	return &Client{
		id:          msg.GetId(),
		vectorClock: NewClock(msg.GetClock()),
	}
}

func (c *Client) Listen(wait chan struct{}, stream grpc.BidiStreamingClient[Message, Message]) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			close(wait)
			return
		}
		if err != nil {
			log.Fatalf("Connection was closed.")
		}
		c.vectorClock.Update(in.GetClock())
		log.Println(in.GetText())
	}
}

func (c *Client) AwaitUserInput(stream grpc.BidiStreamingClient[Message, Message]) error {
	reader := bufio.NewScanner(os.Stdin)
	for {
		reader.Scan()
		if reader.Err() != nil {
			log.Fatalf("fail to call Read: %v", reader.Err())
		}
		text := reader.Text()
		if text == "exit" {
			break
		}
		err := stream.Send(&Message{Text: &text})
		if err != nil {
			log.Fatalf("fail to call Send: %v", err)
		}
	}
	return stream.CloseSend()
}

func (c *Client) Run(stream grpc.BidiStreamingClient[Message, Message]) {
	wait := make(chan struct{})
	go c.Listen(wait, stream)
	_ = c.AwaitUserInput(stream)
	<-wait
}
