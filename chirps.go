package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	fmt.Println(splitWords)
	joinedWords := strings.Join(splitWords, " ")

	return joinedWords, nil
}


func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
    chirpRequest := ChirpJSON{}
    if err := json.NewDecoder(r.Body).Decode(&chirpRequest); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        respondWithError(w, http.StatusBadRequest, "Invalid JSON format")
        return
    }

    // Map ChirpJSON to Chirp structure for validation
    chirp := database.Chirp{
        Body:   chirpRequest.Body,
        UserID: chirpRequest.UserID,
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
