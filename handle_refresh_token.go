package main

import (
	"log"
	"net/http"
	"time"

	"github.com/RodolfoCamposGlz/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), token)
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
