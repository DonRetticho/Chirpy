package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("DeleteAllUsers error: %v", err)
		respondWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("DeleteAllUsers error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	log.Printf("CreateUser error: %v", err)
	respondWithJSON(w, http.StatusCreated, responseUser)

}
