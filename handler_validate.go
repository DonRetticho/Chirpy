package main

import (
	"strings"
)

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
