package internals

import (
	"context"
	"errors"
	"fmt"
	"grpc_note_program/notes"
	"sync"
)

type NoteRequest struct{
	notes.UnimplementedNoteRequestServer
	mu sync.RWMutex
	data map[string]*notes.Note
}

func NewNoteRequest() *NoteRequest{
	return &NoteRequest{
		data: make(map[string]*notes.Note),
	}
}

func(s *NoteRequest) GetNote(ctx context.Context, req *notes.FetchByTitle) (*notes.Note, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	notes, exists := s.data[req.Title]
	if !exists{
		return nil, errors.New("Note Not Found!")
	}
	return notes, nil
}

func(s *NoteRequest) CreateNote(ctx context.Context, req *notes.Note) (*notes.FetchByTitle, error){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[req.Title] = req
	return &notes.FetchByTitle{Title: req.Title}, nil
}

func(s *NoteRequest) ListAllTitles(ctx context.Context, req *notes.ListTitles) (*notes.ListTitlesResponse, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.data) == 0{
		return nil, errors.New("No Saved Notes!")
	}
	titles := make([]string, 0, len(s.data))
	for title := range s.data{
		titles = append(titles, title)
	}
	return &notes.ListTitlesResponse{Titles: titles}, nil
}

func(s *NoteRequest) ListAllNotes(ctx context.Context, req *notes.ListNotes) (*notes.AllNotes, error){
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.data) == 0{
		return nil, errors.New("No Note Found!")
	}
	notesList := make([]*notes.Note, 0, len(s.data))
	for _, notes := range s.data{
		notesList = append(notesList, notes)
	}
	return &notes.AllNotes{Notes: notesList}, nil
}

func(s *NoteRequest) StreamNotes(stream notes.NoteRequest_StreamNotesServer) (error){
	_, err := stream.Recv()
	if err != nil{
		return err
	}

	fmt.Println("Streaming Started...")

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, page := range s.data{
		if err := stream.Send(page); err !=nil{
			return err
		}
	}
	return nil
}