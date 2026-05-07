package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "something went wrong")
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanedBody := getCleanedBody(params.Body)
	respondWithJSON(w, 200, returnVals{CleanedBody: cleanedBody})
}

func getCleanedBody(body string) string {
	splitBody := strings.Split(body, " ")
	for idx, item := range splitBody {
		lowerWord := strings.ToLower(item)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			splitBody[idx] = "****"
		}

	}
	joinedAagain := strings.Join(splitBody, " ")
	return joinedAagain
}
