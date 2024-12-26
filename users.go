package main

import (
	"database/sql"
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
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(
		getUser.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}


	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    getUser.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
		RevokedAt: sql.NullTime{},
	}

	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't store refresh token")
		return
	}

	response := UserResponse{
		ID:        getUser.ID,
		Email:     getUser.Email,
		CreatedAt: getUser.CreatedAt.Time,
		UpdatedAt: getUser.UpdatedAt.Time,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT")
		return
	}
	type request struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON")
		return
	}
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	params := database.UpdateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}
	updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}
	response := UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt.Time,
		UpdatedAt: updatedUser.UpdatedAt.Time,
	}
	respondWithJSON(w, http.StatusOK, response)

}