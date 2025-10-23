package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"log"
)

func main() {
	// Create a gRPC channel
	// If credentials are needed, these are set in the options.
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("error closing connection:\n%v", err)
		}
	}(conn)
}
