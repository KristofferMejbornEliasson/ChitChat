package service

import (
	"log"

	"google.golang.org/grpc"

	. "chitchat/m/constants"
	. "chitchat/m/grpc"
	. "chitchat/m/server/connection"
)

type ChitChatService struct {
	UnimplementedChitChatServer
	Connections  []Connection
	Channel      chan *Message
	nextClientId int32
}

func (s *ChitChatService) nextID() (id int32) {
	id = s.nextClientId
	s.nextClientId++
	return id
}

func (s *ChitChatService) RouteChat(server grpc.BidiStreamingServer[Message, Message]) error {
	clientID := s.nextID()
	log.Printf("Established connection with a new client, ID = %d\n", clientID)
	conn := Connection{Conn: server, HomeChannel: s.Channel, ReceiveChannel: make(chan *Message)}
	s.Connections = append(s.Connections, conn)
	conn.Init(clientID)
	go conn.Listen()
	select {}
}

func (s *ChitChatService) ManageChannels() {
	if VERBOSE {
		log.Println("Managing channels!")
	}
	for {
		select {
		case msg := <-s.Channel:
			if VERBOSE {
				log.Println("Got message from home Channel:" + msg.GetText())
			}
			for _, conn := range s.Connections {
				conn.ReceiveChannel <- msg
			}
		}
	}
}
