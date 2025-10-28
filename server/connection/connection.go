package connection

import (
	"fmt"
	"log"

	"chitchat/m/client/clock"

	"google.golang.org/grpc"

	. "chitchat/m/constants"
	. "chitchat/m/grpc"
)

type Connection struct {
	Conn           grpc.BidiStreamingServer[Message, Message]
	HomeChannel    chan *Message
	ReceiveChannel chan *Message
	Id             int32
	Open           bool
}

func (c *Connection) ListenToClient() {
	for {
		recv, err := c.Conn.Recv()
		if err != nil {
			c.Open = false
			break
		}
		if VERBOSE {
			log.Printf("Received message: %s", recv)
		}
		c.HomeChannel <- recv
	}
}

func (c *Connection) RelayToClient() {
	if VERBOSE {
		log.Printf("Connection #%d listening to channels.\n", c.Id)
	}
	for {
		msg := <-c.ReceiveChannel
		if VERBOSE {
			log.Printf("Connection #%d received a message through its receiving channel.\n", c.Id)
		}
		err := c.Conn.Send(msg)
		if err != nil {
			c.Open = false
			return
		}
		if VERBOSE {
			log.Printf("Connection #%d sent a message through its stream.\n", c.Id)
		}
	}
}

func (c *Connection) Init(clock *clock.Clock) {
	// Ignore contents of the message which established this connection.
	_, err := c.Conn.Recv()
	if err != nil {
		log.Fatalf("Failed to receive initial message from the new client #%d: %v", c.Id, err)
	}

	vector := clock.Vector()
	text := fmt.Sprintf("Client #%d joined the chat at logical time %v", c.Id, vector)

	msg := Message{
		Text:  &text,
		Id:    &c.Id,
		Clock: vector,
	}

	err = c.Conn.Send(&msg)
	if err != nil {
		log.Fatalf("failed to send response to client #%d during initial setup:\n%v", c.Id, err)
	}
	msg.Id = &SYSTEM_ID
	c.HomeChannel <- &msg
}
