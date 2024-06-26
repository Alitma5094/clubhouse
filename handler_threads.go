package main

import (
	"clubhouse/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerThreadsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Title string `json:"title"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	newThread, err := cfg.DB.CreateThread(r.Context(), database.CreateThreadParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, Title: params.Title})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create thread")
		return
	}
	_, err = cfg.DB.SubscribeToThread(r.Context(), database.SubscribeToThreadParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: user.ID, ThreadID: newThread.ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't subscribe to thread")
		return
	}

	respondWithJSON(w, http.StatusCreated, newThread)
}

func (cfg *apiConfig) handlerThreadsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	dbThreads, err := cfg.DB.GetThreads(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get threads")
		return
	}
	respondWithJSON(w, http.StatusOK, dbThreads)
}

func (cfg *apiConfig) handlerThreadsAddUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}

	type parameters struct {
		UserIDs []uuid.UUID `json:"user_ids"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	for _, userID := range params.UserIDs {
		_, err = cfg.DB.SubscribeToThread(r.Context(), database.SubscribeToThreadParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), UserID: userID, ThreadID: id})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't subscribe to thread")
			return
		}
	}
	respondWithJSON(w, http.StatusCreated, nil)
}

func (cfg *apiConfig) handlerThreadsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}

	err = cfg.DB.DeleteThread(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete thread")
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) handlerThreadsGetMembers(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}
	dbUsers, err := cfg.DB.GetSubscribedUsers(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get threads")
		return
	}

	log.Println(len(dbUsers))
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, DatabaseUserToUser(dbUser))
	}
	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerUnsubscribedUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}

	dbUsers, err := cfg.DB.GetUnsubscribedUsers(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, DatabaseUserToUser(dbUser))
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *apiConfig) handlerUnsubscribeUsers(w http.ResponseWriter, r *http.Request, user database.User) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid thread id")
		return
	}

	type parameters struct {
		UserIDs []uuid.UUID `json:"user_ids"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	for _, userID := range params.UserIDs {
		err = cfg.DB.UnsubscribeFromThread(r.Context(), database.UnsubscribeFromThreadParams{ThreadID: id, UserID: userID})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't unsubscribe from thread")
			return
		}
	}
	respondWithJSON(w, http.StatusOK, nil)
}
