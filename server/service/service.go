package service

import (
	"log"

	"chitchat/m/client/clock"

	"google.golang.org/grpc"

	. "chitchat/m/constants"
	. "chitchat/m/grpc"
	. "chitchat/m/server/connection"
)

type ChitChatService struct {
	UnimplementedChitChatServer
	Connections  []*Connection
	Channel      chan *Message
	nextClientId int32
	Clock        *clock.Clock
}

func (s *ChitChatService) nextID() (id int32) {
	id = s.nextClientId
	s.nextClientId++
	return id
}

func (s *ChitChatService) ensureClockExists() {
	if s.Clock == nil {
		s.Clock = clock.NewClock(make([]int64, 0))
	}
}

func (s *ChitChatService) incrementClock(clientID int32) {
	s.ensureClockExists()
	s.Clock.Increment(clientID)
}

func (s *ChitChatService) RouteChat(stream grpc.BidiStreamingServer[Message, Message]) error {
	clientID := s.nextID()
	log.Printf("Established connection with a new client, ID = %d\n", clientID)
	conn := Connection{Conn: stream, HomeChannel: s.Channel,
		ReceiveChannel: make(chan *Message), Id: clientID,
		Open: true}
	s.Connections = append(s.Connections, &conn)
	s.incrementClock(clientID)
	conn.Init(s.Clock)
	go conn.SendFromChannel()
	conn.Listen()
	conn.Open = false
	quitText := fmt.Sprintf("Client #%d left the chat at logical time %s\n", clientID, s.Clock)
	s.Channel <- &Message{
		Text:  &quitText,
		Id:    &clientID,
		Clock: s.Clock.Vector(),
	}
	return nil
}

func (s *ChitChatService) ManageChannels() {
	if VERBOSE {
		log.Println("Managing channels!")
	}
	for {
		msg := <-s.Channel
		s.ensureClockExists()
		s.Clock.Update(msg.GetClock())
		if VERBOSE {
			log.Printf("Got message from home channel:\n%v\n", msg.GetText())
		}
		msg.Clock = s.Clock.Vector()
		for _, conn := range s.Connections {
			if conn.Open {
				conn.ReceiveChannel <- msg
				if VERBOSE {
					log.Printf("ChitChatService relayed message to connection #%d\n", conn.Id)
				}
			}
		}

	}
}
