package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DonRetticho/Chirpy/internal/auth"
	"github.com/DonRetticho/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization")
		return
	}
	validatedUserID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("DeleteAllUsers error: %v", err)
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedBody := getCleanedBody(params.Body)

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: validatedUserID})
	if err != nil {
		log.Printf("DeleteAllChirps error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	responseChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, responseChirp)

}
