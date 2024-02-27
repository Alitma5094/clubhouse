package main

import "net/http"

func (cfg *apiConfig) handlerStatus(w http.ResponseWriter, r *http.Request) {
	type statusResponse struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, http.StatusOK, statusResponse{
		Status: "ok",
	})
}
