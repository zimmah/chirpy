package main

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	type chirps struct {
		Body string `json:"body"`
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirps{}
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

	respondWithJSON(w, http.StatusOK, response{CleanedBody: cleanedChirp})
}