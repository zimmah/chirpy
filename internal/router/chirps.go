package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zimmah/chirpy/internal/database"
)

func handlePostChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprint("Error decoding chirp: %w", err))
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := wordFilter(chirp.Body)
	responseChirp, err := database.DBPointer.CreateChirp(cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
	}

	respondWithJSON(w, http.StatusCreated, responseChirp)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.DBPointer.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error")
	}
	respondWithJSON(w, http.StatusOK, chirps)
}