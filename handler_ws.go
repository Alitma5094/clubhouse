package main

import (
	"clubhouse/internal/auth"
	"clubhouse/internal/database"
	"clubhouse/internal/ws"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWS(w http.ResponseWriter, r *http.Request) {
	// Grab the OTP in the Get param
	token := chi.URLParam(r, "apiKey")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing api token")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.JWTSecret, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), subject)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user")
		return
	}

	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := ws.WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Create new Client
	client := ws.NewClient(conn, cfg.WSManager, user)
	cfg.WSManager.AddClient(client)
	go client.ReadMessages()
	go client.WriteMessages()
}

const (
	// EventSendMessage is the event name for new chat messages sent
	EventThreadsGet     = "threads_get"
	EventThreadsCreated = "threads_created"
	EventMessagesGet    = "messages_get"
)

// SendMessageHandler will send out a message to all other participants in the chat
func (cfg *apiConfig) GetThreadsHandlerWS(event ws.Event, c *ws.Client) error {

	threads, err := cfg.DB.GetThreads(context.Background(), c.UserId)
	if err != nil {
		return fmt.Errorf("failed to get threads: %v", err)
	}

	data, err := json.Marshal(threads)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	c.Egress <- ws.Event{Payload: json.RawMessage(data), Type: EventThreadsGet}
	// Broadcast to all other Clients
	// for client := range c.manager.clients {
	// 	client.egress <- outgoingEvent
	// }

	return nil

}

func (cfg *apiConfig) sendThreadsCreatedWS(thread *database.Thread) {

	data, err := json.Marshal(thread)
	if err != nil {
		log.Printf("failed to marshal broadcast message: %v", err)
		return
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventThreadsCreated}
	// Broadcast to all other Clients
	for client := range cfg.WSManager.Clients {
		if client.UserId == thread.UserID {
			client.Egress <- outgoingEvent
			return
		}
	}
}

type GetMessagesEvent struct {
	ThreadID string `json:"thread_id"`
}

type GetMessagesEventReturn struct {
	ThreadID uuid.UUID          `json:"thread_id"`
	Messages []database.Message `json:"messages"`
}

func (cfg *apiConfig) SendMessagesGetWS(event ws.Event, c *ws.Client) error {

	var chatevent GetMessagesEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	idUUID, err := uuid.Parse(chatevent.ThreadID)
	if err != nil {
		return fmt.Errorf("bad thread id: %v", err)
	}

	dbMessages, err := cfg.DB.GetMessages(context.Background(), idUUID)
	if err != nil {
		return fmt.Errorf("cant get messages: %v", err)
	}

	data, err := json.Marshal(GetMessagesEventReturn{ThreadID: idUUID, Messages: dbMessages})
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventMessagesGet}
	// Broadcast to all other Clients
	c.Egress <- outgoingEvent
	return nil
}
