package main

import (
	"bufio"
	"io"
	"log"
	"os"

	. "chitchat/m/client/clock"
	. "chitchat/m/grpc"

	"google.golang.org/grpc"
)

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
