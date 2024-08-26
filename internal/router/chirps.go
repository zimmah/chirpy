package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/zimmah/chirpy/internal/database"
)

func (cfg *apiConfig) handlePostChirps(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userIDString, err := cfg.validateJWT(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := wordFilter(chirp.Body)
	responseChirp, err := database.DBPointer.CreateChirp(cleanedChirp, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, responseChirp)
}

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID")
		return
	}
	
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userIDString, err := cfg.validateJWT(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	chirp, err := database.DBPointer.GetChirpByID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirp")
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "You are not authorized to delete this chirp")
		return
	}

	err = database.DBPointer.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request) {
	sort := "asc"
	authorQ := r.URL.Query().Get("author_id")
	sortQ := r.URL.Query().Get("sort")
	if sortQ == "desc" {  sort = "desc" }

	if authorQ != "" {
		authorID, err := strconv.Atoi(authorQ)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author id")
			return
		}
		chirps, err := database.DBPointer.GetChirpsOfAuthor(authorID, sort)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	chirps, err := database.DBPointer.GetChirps(sort)
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