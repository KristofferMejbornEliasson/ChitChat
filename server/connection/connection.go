package connection

import (
	"log"

	"google.golang.org/grpc"

	. "chitchat/m/constants"
	. "chitchat/m/grpc"
)

type Connection struct {
	Conn           grpc.BidiStreamingServer[Message, Message]
	HomeChannel    chan *Message
	ReceiveChannel chan *Message
}

func (c *Connection) Listen() {
	go func() {
		if VERBOSE {
			log.Println("Listening to channels.")
		}
		for {
			msg := <-c.ReceiveChannel
			err := c.Conn.Send(msg)
			if err != nil {
				log.Fatalf("fail to send message: %v", err)
			}
		}
	}()
	if VERBOSE {
		log.Println("Listening to stream.")
	}
	for {
		recv, err := c.Conn.Recv()
		if err != nil {
			return
		}
		if VERBOSE {
			log.Printf("Received message: %s", recv)
		}
		c.HomeChannel <- recv
		// log.Printf("Conveyed message to home: %s", recv)
	}
}

func (c *Connection) Init(clientId int32) {
	err := c.Conn.Send(&Message{
		Id:    &clientId,
		Clock: make([]int64, clientId),
	})
	if err != nil {
		log.Fatalf("failed to send response to client during initial setup:\n%v", err)
	}
}
