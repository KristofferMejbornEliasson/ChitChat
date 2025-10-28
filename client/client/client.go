package client

import (
	"bufio"
	"fmt"
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
	logger      *log.Logger
}

func (c *Client) IncrementClock() {
	c.vectorClock.Increment(c.id)
}

func (c *Client) IncrementAndCopyClock() (updatedCopy []int64) {
	return c.vectorClock.IncrementAndCopy(c.id)
}

func NewClient(stream grpc.BidiStreamingClient[Message, Message], logger *log.Logger) *Client {
	in, err := stream.Recv()
	if err != nil {
		logger.Fatalf("No initial response received from server:\n%v", err)
	}
	client := Client{
		id:          in.GetId(),
		vectorClock: NewClock(in.GetClock()),
		stream:      stream,
		logger:      logger,
	}
	client.Log("Established connection with server.\n")
	return &client
}

func (c *Client) Log(message string) {
	prefixString := fmt.Sprintf("Client #%d %s: ", c.id, c.vectorClock.String())
	c.logger.SetPrefix(prefixString)
	c.logger.Print(message)
}

func (c *Client) Logf(format string, v ...any) {
	c.Log(fmt.Sprintf(format, v...))
}

func (c *Client) PrintMessage(msg *Message) {
	if msg == nil {
		return
	}
	var text string
	if msg.GetId() == -1 {
		text = msg.GetText()
	} else {
		text = fmt.Sprintf("Message received from client #%d:\n%s", msg.GetId(), msg.GetText())
	}
	c.Logf("%s\n", text)
	fmt.Printf("%s: %s\n", c.vectorClock, text)
}

func (c *Client) Run() {
	wait := make(chan struct{})
	go func() {
		for {
			in, err := c.stream.Recv()
			if err != nil {
				wait <- struct{}{}
				break
			}
			c.IncrementClock()
			c.vectorClock.Update(in.GetClock())
			c.PrintMessage(in)
		}
	}()
	reader := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			reader.Scan()
			if reader.Err() != nil {
				log.Fatalf("fail to call Read: %v", reader.Err())
			}
			text := reader.Text()
			if text == "exit" {
				wait <- struct{}{}
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
	}()
	for {
		select {
		case <-wait:
			_ = c.stream.CloseSend()
			return
		}
	}
}
