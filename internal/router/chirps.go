package router

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zimmah/chirpy/internal/database"
)

func handlePostChirps(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := wordFilter(chirp.Body)
	responseChirp, err := database.DBPointer.CreateChirp(cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, responseChirp)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := database.DBPointer.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID")
		return
	}

	chirp, err := database.DBPointer.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}