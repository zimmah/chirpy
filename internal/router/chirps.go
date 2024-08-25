package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zimmah/chirpy/internal/database"
)

func handlePostChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding chirp: %v", err))
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := wordFilter(chirp.Body)
	responseChirp, err := database.DBPointer.CreateChirp(cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, responseChirp)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.DBPointer.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database error: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not parse request: %v", err))
		return
	}

	chirp, statusCode, err := database.DBPointer.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, statusCode, fmt.Sprintf("Error loading chirp: %v", err))
		return
	}

	respondWithJSON(w, statusCode, chirp)
}