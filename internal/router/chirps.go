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

	// increment ID by 1 and save to database (note: start at ID 1)
	id := 1 // placeholder

	respondWithJSON(w, http.StatusCreated, database.Chirp{Id: id, Body: cleanedChirp})
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	// loop over chirps in DB and return a JSON array of chirps ordered by ID ascending
	// return status code OK
	respondWithJSON(w, http.StatusOK, []database.Chirp{})
}