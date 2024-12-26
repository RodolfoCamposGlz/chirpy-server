package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RodolfoCamposGlz/internal/auth"
	"github.com/RodolfoCamposGlz/internal/database"
	"github.com/google/uuid"
)


type ChirpJSON struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
	ID     uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, chirp *database.Chirp) (string, error) {
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return "", fmt.Errorf("chirp too long")
	}

	//Clean body input
	splitWords := strings.Split(chirp.Body, " ")
	profanityMap := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	for i, word := range splitWords {
		lowerWord := strings.ToLower(word)
		if profanityMap[lowerWord] {
			splitWords[i] = "****" 
		}

	}
	joinedWords := strings.Join(splitWords, " ")

	return joinedWords, nil
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {


	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

    chirpRequest := ChirpJSON{}
    if err := json.NewDecoder(r.Body).Decode(&chirpRequest); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    // Map ChirpJSON to Chirp structure for validation
    chirp := database.Chirp{
        Body:   chirpRequest.Body,
        UserID: userID,
    }

    // Validate the chirp's body
    cleanedBody, err := cfg.handlerValidateChirp(w, &chirp)
    if err != nil {
        return // Validation failed, already handled in handlerValidateChirp
    }


	params := database.CreateChirpParams{
		Body: cleanedBody,
		UserID: chirp.UserID,
	}
	newChirp, err := cfg.dbQueries.CreateChirp(r.Context(), params)
	if err != nil {
		log.Println("Error creating chirp", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}
	response := ChirpJSON{
		ID: newChirp.ID,
		Body: newChirp.Body,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		UserID: newChirp.UserID,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	chirpJSONs := []ChirpJSON{}
	for _, chirp := range chirps {
		chirpJSONs = append(chirpJSONs, ChirpJSON{
			ID: chirp.ID,
			Body: chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID: chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, chirpJSONs)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp")
		return
	}

	response := ChirpJSON{
		ID: chirp.ID,
		Body: chirp.Body,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldn't validate JWT")
		return
	}
	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not found")
		return
	}

	//check if the user is the owner of the chirp
	if userId != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "You are not the owner of this chirp")
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), id)
	if err != nil {
		log.Println("Error deleting chirp", err)
		respondWithError(w, http.StatusNotFound, "Error deleting chirp")
		return
	}
	respondWithJSON(w, http.StatusNoContent, map[string]string{"message": "Chirp deleted successfully"})

}
