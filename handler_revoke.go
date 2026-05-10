package main

import (
	"net/http"

	"github.com/DonRetticho/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or malformed authorization header")
		return
	}
	err = cfg.dbQueries.RevokeFreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to revoke token")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
