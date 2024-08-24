package router

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"

	"github.com/zimmah/chirpy/internal/database"
)

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.DBPointer.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Database error: %w", err))
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Could not parse request: %w", err))
		return
	}

	user, statusCode, err := database.DBPointer.GetUserByID(userID)
	if err != nil {
		respondWithError(w, statusCode, fmt.Sprint("Error loading user: %w", err))
		return
	}

	respondWithJSON(w, statusCode, user)
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := database.User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error decoding user: %w", err))
		return
	}

	responseUser, err := database.DBPointer.CreateUser(user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Database error: %w", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, responseUser)
}