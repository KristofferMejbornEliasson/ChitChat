package main

import (
	"log"
	"net"
	"os"

	. "chitchat/m/grpc"
	. "chitchat/m/server/connection"
	. "chitchat/m/server/service"

	"google.golang.org/grpc"
)

func main() {
	file, err := os.Create("log.txt")
	logger := log.New(file, "[]: ", 0)
	service := ChitChatService{
		Connections: []*Connection{},
		Channel:     make(chan *Message),
		Logger:      logger,
	}
	defer func(writer *os.File, logger *log.Logger) {
		logger.SetPrefix(service.Clock.String() + ": ")
		logger.Println("Server shut down.")
		if err != nil {
			_ = file.Close()
		}
	}(file, logger)

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	logger.Printf("Listening on %v\n", lis.Addr())
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterChitChatServer(grpcServer, &service)
	go service.ManageChannels()
	wait := make(chan struct{})
	go ListenToUserInput(wait)
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			wait <- struct{}{}
			logger.Fatalf("failed to serve: %v", err)
		}
	}()
	for {
		select {
		case <-wait:
			return
		}
	}
}

func ListenToUserInput(wait chan struct{}) {
	reader := bufio.NewScanner(os.Stdin)
	for {
		reader.Scan()
		if reader.Err() != nil {
			log.Fatalf("fail to read user input: %v", reader.Err())
		}
		text := reader.Text()
		if text == "exit" {
			break
		}
	}
	wait <- struct{}{}
}
