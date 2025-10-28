package client

import (
	"bufio"
	"io"
	"log"
	"os"

	. "chitchat/m/client/clock"
	. "chitchat/m/constants"
	. "chitchat/m/grpc"

	"google.golang.org/grpc"
)

type Client struct {
	id          int32
	vectorClock *Clock
	stream      grpc.BidiStreamingClient[Message, Message]
}

func (c *Client) IncrementClock() {
	c.vectorClock.Increment(c.id)
}

func (c *Client) IncrementAndCopyClock() (updatedCopy []int64) {
	return c.vectorClock.IncrementAndCopy(c.id)
}

func NewClient(stream grpc.BidiStreamingClient[Message, Message]) *Client {
	in, err := stream.Recv()
	if err != nil {
		log.Fatalf("No initial response received from server:\n%v", err)
	}

	return &Client{
		id:          in.GetId(),
		vectorClock: NewClock(in.GetClock()),
		stream:      stream,
	}
}

func (c *Client) PrintMessage(msg *Message) {
	fmt.Printf("Message received from client #%d at %s:\n%s\n", msg.GetId(), c.vectorClock, msg.GetText())
}

func (c *Client) Run() {
	wait := make(chan struct{})
	go func() {
		for {
			in, err := c.stream.Recv()
			if err != nil {
				log.Fatalf("Connection was closed.")
			}
			c.IncrementClock()
			c.vectorClock.Update(in.GetClock())
			c.PrintMessage(in)
		}
	}()
	reader := bufio.NewScanner(os.Stdin)
	for {
		reader.Scan()
		if reader.Err() != nil {
			log.Fatalf("fail to call Read: %v", reader.Err())
		}
		text := reader.Text()
		if text == "exit" {
			close(wait)
			break
		}
		if text == "" {
			continue
		}

		err := c.stream.Send(&Message{
			Text:  &text,
			Id:    &c.id,
			Clock: c.IncrementAndCopyClock(),
		})
		if VERBOSE {
			fmt.Printf("Client #%d sent a message through its stream.\n", c.id)
		}
		if err != nil {
			log.Fatalf("fail to call Send: %v", err)
		}
	}
	_ = c.stream.CloseSend()
	<-wait
}
