package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader { // Ensure "Bearer " prefix exists
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization format")
		return
	}

	err := cfg.dbQueries.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "Refresh token revoked")
}
