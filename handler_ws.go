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

var (
	WebsocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Grab the request origin
			// origin := r.Header.Get("Origin")
			// switch origin {
			// case "http://localhost:8080":
			// 	return true
			// default:
			// 	return false
			// }
			return true
		},
	}
)

func (cfg *apiConfig) SetupEventHandlers() {
	cfg.WSManager.Handlers[EventThreadsGet] = cfg.EventThreadsGet
	cfg.WSManager.Handlers[EventMessagesGet] = cfg.EventMessagesGet
	cfg.WSManager.Handlers[EventMessagesCreate] = cfg.EventMessagesCreate
}

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
	conn, err := WebsocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := cfg.NewClient(conn, cfg.WSManager, user)
	cfg.WSManager.AddClient(client)
	go client.ReadMessages()
	go client.WriteMessages()
}

var (
	EventThreadsGet     = "threads_get"
	EventThreadsCreate  = "threads_create"
	EventMessagesGet    = "messages_get"
	EventMessagesCreate = "messages_create"
)

type GetMessagesEvent struct {
	ThreadID string `json:"thread_id"`
}
type GetMessagesReturnEvent struct {
	ThreadID uuid.UUID          `json:"thread_id"`
	Messages []database.Message `json:"messages"`
}
type CreateMessageEvent struct {
	ThreadID string `json:"thread_id"`
	Text     string `json:"text"`
}

func (cfg *apiConfig) EventThreadsGet(event ws.Event, c *ws.Client) error {
	threads, err := cfg.DB.GetThreads(context.Background(), c.UserId)
	if err != nil {
		return fmt.Errorf("failed to get threads: %v", err)
	}
	data, err := json.Marshal(threads)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}
	c.Egress <- ws.Event{Payload: json.RawMessage(data), Type: EventThreadsGet}
	return nil
}

func (cfg *apiConfig) EventMessagesGet(event ws.Event, c *ws.Client) error {
	var params GetMessagesEvent
	if err := json.Unmarshal(event.Payload, &params); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	id, err := uuid.Parse(params.ThreadID)
	if err != nil {
		return fmt.Errorf("bad thread id: %v", err)
	}

	dbMessages, err := cfg.DB.GetMessages(context.Background(), id)
	if err != nil {
		return fmt.Errorf("cant get messages: %v", err)
	}

	data, err := json.Marshal(GetMessagesReturnEvent{ThreadID: id, Messages: dbMessages})
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create a new Event
	outgoingEvent := ws.Event{Payload: json.RawMessage(data), Type: EventMessagesGet}
	// Broadcast to all other Clients
	c.Egress <- outgoingEvent
	return nil

}

func (cfg *apiConfig) EventMessagesCreate(event ws.Event, c *ws.Client) error {
	var params CreateMessageEvent
	if err := json.Unmarshal(event.Payload, &params); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	err := ValidateMessageInput(params.Text)
	if err != nil {
		return fmt.Errorf("bad message: %v", err)
	}

	id, err := uuid.Parse(params.ThreadID)
	if err != nil {
		return fmt.Errorf("bad thread id: %v", err)
	}

	newMessage, err := cfg.DB.CreateMessage(context.Background(), database.CreateMessageParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: c.UserId, Text: params.Text, ThreadID: id})
	if err != nil {
		return fmt.Errorf("failed to create message: %v", err)
	}

	return cfg.BroadcastMessagesCreated(&newMessage)
}
