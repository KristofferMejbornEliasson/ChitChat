package service

import (
	"fmt"
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
	Logger       *log.Logger
}

func (s *ChitChatService) Log(message string) {
	s.ensureClockExists()
	s.Logger.SetPrefix("Server " + s.Clock.String() + ": ")
	s.Logger.Printf(message)
}

func (s *ChitChatService) Logf(format string, v ...any) {
	s.Log(fmt.Sprintf(format, v...))
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
	s.Logf("Established connection with a new client, ID = %d\n", clientID)
	conn := Connection{Conn: stream, HomeChannel: s.Channel,
		ReceiveChannel: make(chan *Message), Id: clientID,
		Open: true}
	s.Connections = append(s.Connections, &conn)
	s.incrementClock(clientID)
	conn.Init(s.Clock)
	go conn.RelayToClient()
	conn.ListenToClient()
	conn.Open = false
	s.Logf("Closed connection to client #%d.\n", clientID)
	quitText := fmt.Sprintf("Client #%d left the chat", clientID)
	s.Channel <- &Message{
		Text:  &quitText,
		Id:    &SYSTEM_ID,
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
		if len(msg.GetText()) > CHAT_MESSAGE_LENGTH_LIMIT {
			s.Logf("Message from client #%d was too long; rejected.", msg.GetId())
			if s.Connections[msg.GetId()].Open == true {
				rejectionText := "Message too long. Limit is 128 characters."
				rejection := Message{
					Text:  &rejectionText,
					Id:    &SYSTEM_ID,
					Clock: s.Clock.Vector(),
				}
				s.Connections[msg.GetId()].ReceiveChannel <- &rejection
			}
		} else {
			for _, conn := range s.Connections {
				if conn.Open {
					conn.ReceiveChannel <- msg
					s.Logf("Broadcasting message to client #%d\n", conn.Id)
					if VERBOSE {
						log.Printf("ChitChatService relayed message to connection #%d\n", conn.Id)
					}
				}
			}
		}
	}
}
