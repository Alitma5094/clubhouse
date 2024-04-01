package main

import (
	"clubhouse/internal/database"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMessagesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Text        string    `json:"text"`
		ThreadID    uuid.UUID `json:"thread_id"`
		Attachments []string  `json:"attachments"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = ValidateMessageInput(params.Text)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newMessage, err := cfg.DB.CreateMessage(r.Context(), database.CreateMessageParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, Text: params.Text, ThreadID: params.ThreadID, Attachments: params.Attachments})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create message")
		return
	}

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := cfg.firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	fcm_tokens, err := cfg.DB.GetFcmTokens(r.Context(), params.ThreadID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get notification tokens")
		return
	}

	data, err := json.Marshal(newMessage)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	for _, token := range fcm_tokens {
		message := &messaging.Message{
			Data: map[string]string{"payload": string(data), "type": "new_message"},
			Notification: &messaging.Notification{
				Title: "New message",
				Body:  params.Text,
			},
			Token: token.Token,
		}

		// Send a message to the device corresponding to the provided
		// registration token.
		_, err = client.Send(ctx, message)
		if err != nil {
			log.Fatalln(err)
		}
	}

	respondWithJSON(w, http.StatusCreated, newMessage)
}

func (cfg *apiConfig) handlerMessagesGet(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}

	dbMessages, err := cfg.DB.GetMessages(r.Context(), id)
	if err != nil {
		log.Print(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't get messages")
		return
	}

	respondWithJSON(w, http.StatusOK, dbMessages)
}

var (
	ErrMessageEmptyText = errors.New("text cannot be empty")
	ErrMessageTooLong   = errors.New("text cannot be longer than 500 characters")
)

func ValidateMessageInput(text string) error {
	if text == "" {
		return ErrMessageEmptyText
	}
	if len(text) > 500 {
		return ErrMessageTooLong
	}
	return nil
}
