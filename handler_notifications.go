package main

import (
	"clubhouse/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerNotificationsRegisterDevice(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters")
		return
	}

	_, err = cfg.DB.CreateFcmToken(r.Context(), database.CreateFcmTokenParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UserID: user.ID, Token: params.Token})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create FCM token")
		return
	}

	respondWithJSON(w, http.StatusCreated, nil)
}
