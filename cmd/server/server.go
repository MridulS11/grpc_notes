package main

import (
	"grpc_note_program/internals"
	"grpc_note_program/notes"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main(){
	lis, err := net.Listen("tcp", ":50051")
	if err != nil{
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	notes.RegisterNoteRequestServer(grpcServer, internals.NewNoteRequest())
	reflection.Register(grpcServer)
	log.Println("Server Starting...")
	if err := grpcServer.Serve(lis); err != nil{
		log.Fatalf("Failed to start the server: %v", err)
	}
}