package main

import (
	"clubhouse/internal/database"
	"clubhouse/internal/ws"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

func (cfg *apiConfig) BroadcastThreadsCreate(thread *database.Thread) {

	data, err := json.Marshal(thread)
	if err != nil {
		log.Printf("failed to marshal broadcast message: %v", err)
		return
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventThreadsCreate}
	// Broadcast to all other Clients
	for client := range cfg.WSManager.Clients {
		if client.UserId == thread.UserID {
			client.Egress <- outgoingEvent
			return
		}
	}
}

func (cfg *apiConfig) BroadcastMessagesCreated(message *database.Message) error {

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventMessagesCreate}
	// Broadcast to all other Clients
	// c.Egress <- outgoingEvent
	// return nil
	for client := range cfg.WSManager.Clients {
		if contains(client.SubscribedThreads, message.ThreadID) {
			client.Egress <- outgoingEvent
		}
	}
	return nil
}

func contains(s []uuid.UUID, e uuid.UUID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
