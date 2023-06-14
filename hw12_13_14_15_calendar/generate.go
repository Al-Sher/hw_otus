package main

//go:generate protoc --go_out=./internal/server/grpc/pb --go-grpc_out=./internal/server/grpc/pb ./api/EventService.proto
