install: grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go server.exe client.exe

grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go: grpc/chitchat.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/chitchat.proto

proto: grpc/chitchat_grpc.pb.go grpc/chitchat.pb.go

server.exe: server/server.go server/connection/connection.go server/service/service.go constants/const.go
	go build server/server.go

client.exe: client/main.go client/client/client.go client/clock/clock.go constants/const.go
	go build -o client.exe client/main.go

start: install
	./server