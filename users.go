package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RodolfoCamposGlz/internal/auth"
	"github.com/RodolfoCamposGlz/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON")
		return
	}
	dbQueries := cfg.dbQueries

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}	
	params := database.CreateUserParams{
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}
	createdUser, err := dbQueries.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	response := UserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
	}
	
	respondWithJSON(w, http.StatusCreated, response)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON")
		return
	}
	dbQueries := cfg.dbQueries
	params := database.CreateUserParams{
		Email:          user.Email,
	}
	getUser, err := dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting user")
		return
	}
	isCorrectPassword := auth.CheckPasswordHash(user.Password, getUser.HashedPassword)
	if !isCorrectPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	response := UserResponse{
		ID:        getUser.ID,
		Email:     getUser.Email,
		CreatedAt: getUser.CreatedAt.Time,
		UpdatedAt: getUser.UpdatedAt.Time,
	}
	respondWithJSON(w, http.StatusOK, response)
}
