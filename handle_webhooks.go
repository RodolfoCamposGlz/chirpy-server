package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/RodolfoCamposGlz/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)

	var req requestBody
	err := decoder.Decode(&req)
	if err != nil {
		log.Println("Error decoding request body", err)
		respondWithError(w, http.StatusBadRequest, "Error decoding request body")
		return
	}

	if req.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Event not supported")
		return
	}
	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := cfg.dbQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	cfg.dbQueries.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
		ID:          user.ID,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	respondWithJSON(w, http.StatusNoContent, "User upgraded")
}
