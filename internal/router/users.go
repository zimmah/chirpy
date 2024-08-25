package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zimmah/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.DBPointer.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse request: %v", err))
		return
	}

	user, statusCode, err := database.DBPointer.GetUserByID(userID)
	if err != nil {
		respondWithError(w, statusCode, fmt.Sprintf("Error loading user: %v", err))
		return
	}

	respondWithJSON(w, statusCode, user)
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding user: %v", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating hash: %v", err))
		return
	}

	responseUser, err := database.DBPointer.CreateUser(user.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, responseUser)
}

func (cfg *apiConfig) handlePutUsers(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userIDString, err := cfg.validateJWT(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("%v", err))
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	user, statusCode, err := database.DBPointer.GetUserByID(userID)
	if err != nil {
		respondWithError(w, statusCode, fmt.Sprintf("%v", err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	userReq := database.User{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding user: %v", err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating hash: %v", err))
		return
	}

	userResp, err := database.DBPointer.UpdateUser(user.ID, userReq.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	respondWithJSON(w, statusCode, userResp)
}