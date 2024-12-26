package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RodolfoCamposGlz/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
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

	refreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	//verify if the refresh token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	//verify if the refresh token is revoked
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	//create a new access token
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
