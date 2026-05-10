package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DonRetticho/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}

	if apiKey != cfg.polkaAPI {
		respondWithError(w, http.StatusUnauthorized, "wrong API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	polkaReq := polkaWebhookRequest{}
	err = decoder.Decode(&polkaReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	if polkaReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userIDString := polkaReq.Data.UserID
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error handling the uuid")
		return
	}

	_, err = cfg.dbQueries.UpdateChirpyRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "no user found")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
