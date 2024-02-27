package main

import (
	"clubhouse/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerEventsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	dbEvents, err := cfg.DB.GetEvents(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, dbEvents)
}

func (cfg *apiConfig) handlerEventsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title   string    `json:"title"`
		StartAt time.Time `json:"start_at"`
		EndAt   time.Time `json:"end_at"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	err = ValidateEventInput(params.Title, params.StartAt, params.EndAt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newEvent, err := cfg.DB.CreateEvent(r.Context(), database.CreateEventParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Title: params.Title, StartAt: params.StartAt, EndAt: params.EndAt, UserID: user.ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create event")
		return
	}

	respondWithJSON(w, http.StatusCreated, newEvent)
}

var (
	ErrEventEmptyTitle     = errors.New("title cannot be empty")
	ErrEventTitleTooLong   = errors.New("text cannot be longer than 500 characters")
	ErrEventStartAtMissing = errors.New("start_at is required")
	ErrEventEndAtMissing   = errors.New("end_at is required")
	ErrEventEndBeforeStart = errors.New("start_at must be before end_at")
)

func ValidateEventInput(title string, startAt time.Time, endAt time.Time) error {
	if title == "" {
		return ErrEventEmptyTitle
	}
	if len(title) > 50 {
		return ErrEventTitleTooLong
	}
	if startAt.IsZero() {
		return ErrEventStartAtMissing
	}
	if endAt.IsZero() {
		return ErrEventEndAtMissing
	}
	if startAt.After(endAt) {
		return ErrEventEndBeforeStart
	}
	return nil
}
