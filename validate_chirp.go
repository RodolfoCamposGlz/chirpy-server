package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type chirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
		w.WriteHeader(code)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := struct {
		Body string `json:"body"`
	}{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, 500, "Error decoding JSON")
		return
	}
	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
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


	respondWithJSON(w, 200, chirpResponse{
		CleanedBody:joinedWords,
	})
}