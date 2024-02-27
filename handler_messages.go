package main

import (
	"clubhouse/internal/database"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMessagesCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Text     string    `json:"text"`
		ThreadID uuid.UUID `json:"thread_id"`
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

	newMessage, err := cfg.DB.CreateMessage(r.Context(), database.CreateMessageParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, Text: params.Text, ThreadID: params.ThreadID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create message")
		return
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
