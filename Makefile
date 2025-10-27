install: grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go server client

grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go: grpc/chitchat.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/chitchat.proto

proto: grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go

server: server/server.go
	go build server/server.go

client: client/client.go client/clock/clock.go
	go build client/client.go

start: install
	./server