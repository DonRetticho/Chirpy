package main

import (
	"net/http"
	"sort"

	"github.com/DonRetticho/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "desc" {
		sortParam = "asc"
	}

	var dbChirps []database.Chirp
	var err error
	var authorID uuid.UUID

	responseChirps := []Chirp{}

	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error handling the uuid")
			return
		}
		dbChirps, err = cfg.dbQueries.GetChirpsByAuthor(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	} else {
		dbChirps, err = cfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	for _, chrp := range dbChirps {
		responseChirps = append(responseChirps, Chirp{
			ID:        chrp.ID,
			CreatedAt: chrp.CreatedAt,
			UpdatedAt: chrp.UpdatedAt,
			Body:      chrp.Body,
			UserID:    chrp.UserID,
		})
	}

	sort.Slice(responseChirps, func(i, j int) bool {
		if sortParam == "desc" {
			return responseChirps[i].CreatedAt.After(responseChirps[j].CreatedAt)
		}
		return responseChirps[i].CreatedAt.Before(responseChirps[j].CreatedAt)
	})
	respondWithJSON(w, http.StatusOK, responseChirps)

}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	uuidChirp, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error handling the uuid")
		return
	}

	dcChirp, err := cfg.dbQueries.GetChirp(r.Context(), uuidChirp)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	responseChirp := Chirp{
		ID:        dcChirp.ID,
		CreatedAt: dcChirp.CreatedAt,
		UpdatedAt: dcChirp.UpdatedAt,
		Body:      dcChirp.Body,
		UserID:    dcChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, responseChirp)
}
