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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func (cfg *apiConfig) NewClient(conn *websocket.Conn, manager *ws.Manager, user database.User) *ws.Client {
	subs, err := cfg.DB.GetUserSubscribedThreads(context.Background(), user.ID)
	if err != nil {
		log.Println("failed to get user subscribed threads", err)
	}

	return &ws.Client{
		Connection:        conn,
		Manager:           manager,
		Egress:            make(chan ws.Event),
		UserId:            user.ID,
		SubscribedThreads: subs,
	}
}

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
	client := cfg.NewClient(conn, cfg.WSManager, user)
	cfg.WSManager.AddClient(client)
	go client.ReadMessages()
	go client.WriteMessages()
}

const (
	// EventSendMessage is the event name for new chat messages sent
	EventThreadsGet     = "threads_get"
	EventThreadsCreated = "threads_created"
	EventMessagesGet    = "messages_get"
	EventMessagesCreate = "messages_create"
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

type SimplifiedMessage struct {
	ID                  uuid.UUID      `json:"id"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	UserID              uuid.UUID      `json:"user_id"`
	Text                string         `json:"text"`
	ThreadID            uuid.UUID      `json:"thread_id"`
	AttachmentMediaType database.Media `json:"attachment_media_type"`
	AttachmentUrl       string         `json:"attachment_url"`
}

type GetMessagesEventReturn struct {
	ThreadID uuid.UUID           `json:"thread_id"`
	Messages []SimplifiedMessage `json:"messages"`
}

func ConvertToSimplifiedMessage(original database.GetMessagesWithAttachmentRow) SimplifiedMessage {
	var attachmentMediaType database.Media
	var attachmentUrl string

	if original.AttachmentMediaType.Valid {
		attachmentMediaType = original.AttachmentMediaType.Media
	}

	if original.AttachmentUrl.Valid {
		attachmentUrl = original.AttachmentUrl.String
	}

	return SimplifiedMessage{
		ID:                  original.ID,
		CreatedAt:           original.CreatedAt,
		UpdatedAt:           original.UpdatedAt,
		UserID:              original.UserID,
		Text:                original.Text,
		ThreadID:            original.ThreadID,
		AttachmentMediaType: attachmentMediaType,
		AttachmentUrl:       attachmentUrl,
	}
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

	dbMessages, err := cfg.DB.GetMessagesWithAttachment(context.Background(), idUUID)
	if err != nil {
		return fmt.Errorf("cant get messages: %v", err)
	}

	var simplifiedMessages []SimplifiedMessage
	for _, message := range dbMessages {
		simplifiedMessages = append(simplifiedMessages, ConvertToSimplifiedMessage(message))

	}

	data, err := json.Marshal(GetMessagesEventReturn{ThreadID: idUUID, Messages: simplifiedMessages})
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventMessagesGet}
	// Broadcast to all other Clients
	c.Egress <- outgoingEvent
	return nil
}

type CreateMessagesEvent struct {
	ThreadID string `json:"thread_id"`
	Text     string `json:"text"`
}

func (cfg *apiConfig) CreateMessageHandlerWS(event ws.Event, c *ws.Client) error {

	var eventData CreateMessagesEvent
	if err := json.Unmarshal(event.Payload, &eventData); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	err := ValidateMessageInput(eventData.Text)
	if err != nil {
		return fmt.Errorf("bad message: %v", err)
	}

	idUUID, err := uuid.Parse(eventData.ThreadID)
	if err != nil {
		return fmt.Errorf("bad thread id: %v", err)
	}

	newMessage, err := cfg.DB.CreateMessage(context.Background(), database.CreateMessageParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: c.UserId, Text: eventData.Text, ThreadID: idUUID})
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	data, err := json.Marshal(newMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventMessagesCreate}
	// Broadcast to all other Clients
	// c.Egress <- outgoingEvent
	// return nil
	for client := range cfg.WSManager.Clients {
		if contains(client.SubscribedThreads, idUUID) {
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
