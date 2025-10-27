package main

import (
	"fmt"
	"log"
	"net"

	. "chitchat/m/grpc"

	"google.golang.org/grpc"
)

const VERBOSE = false

type ChitChatService struct {
	UnimplementedChitChatServer
	connections  []Connection
	channel      chan *Message
	nextClientId int32
}

type Connection struct {
	conn           grpc.BidiStreamingServer[Message, Message]
	homeChannel    chan *Message
	receiveChannel chan *Message
}

func (s *ChitChatService) nextID() (id int32) {
	id = s.nextClientId
	s.nextClientId++
	return id
}

func (s *ChitChatService) RouteChat(server grpc.BidiStreamingServer[Message, Message]) error {
	clientID := s.nextID()
	log.Printf("Established connection with a new client, ID = %d\n", clientID)
	conn := Connection{server, s.channel, make(chan *Message)}
	s.connections = append(s.connections, conn)
	conn.Init(clientID)
	go conn.Listen()
	select {}
}

func (c *Connection) Listen() {
	go func() {
		if VERBOSE {
			log.Println("Listening to channels.")
		}
		for {
			msg := <-c.receiveChannel
			err := c.conn.Send(msg)
			if err != nil {
				log.Fatalf("fail to send message: %v", err)
			}
		}
	}()
	if VERBOSE {
		log.Println("Listening to stream.")
	}
	for {
		recv, err := c.conn.Recv()
		if err != nil {
			return
		}
		if VERBOSE {
			log.Printf("Received message: %s", recv)
		}
		c.homeChannel <- recv
		// log.Printf("Conveyed message to home: %s", recv)
	}
}

func (c *Connection) Init(clientId int32) {
	err := c.conn.Send(&Message{
		Id:    &clientId,
		Clock: make([]int64, clientId),
	})
	if err != nil {
		log.Fatalf("failed to send response to client during initial setup:\n%v", err)
	}
}

func (s *ChitChatService) ManageChannels() {
	if VERBOSE {
		log.Println("Managing channels!")
	}
	for {
		select {
		case msg := <-s.channel:
			if VERBOSE {
				log.Println("Got message from home channel:" + msg.GetText())
			}
			for _, conn := range s.connections {
				conn.receiveChannel <- msg
			}
		}
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
	service := ChitChatService{
		connections: []Connection{},
		channel:     make(chan *Message),
	}
	RegisterChitChatServer(grpcServer, &service)
	go service.ManageChannels()
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
