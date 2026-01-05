package main

import (
	"context"
	"grpc_note_program/notes"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main(){
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("Failed To Connect To The Server: %v", err)
	}
	defer conn.Close()

	client := notes.NewNoteRequestClient(conn)
	stream, err := client.StreamNotes(context.Background())
	if err != nil{
		log.Fatalf("Failed To Fetch Notes: %v",err)
	}
	if err = stream.Send(&notes.ListNotes{}); err != nil{
		log.Fatalf("Failed To Fetch Notes: %v",err)
	}
	stream.CloseSend()

	log.Println("--- Streaming Notes from Server ---")

	for {
		note, err := stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil{
			log.Fatalf("Error Fetching Data: %v", err)
		}
		log.Printf("Received Note: [%s] Title: %s", note.Date, note.Title)
	}
	log.Println("Streaming Finished!")
}
