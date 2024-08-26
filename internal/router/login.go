package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zimmah/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 			string `json:"password"`
		Email				string `json:"email"`
		ExpiresInSeconds	int `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	
	user, err := database.DBPointer.GetUserByEmail(userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Invalid username/password combination")
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userReq.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Invalid username/password combination")
		return
	}
	
	jwt, err := cfg.generateJWT(user.ID, userReq.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT")
		return
	}

	resp := response{
		User: User{ID: user.ID, Email: user.Email},
		Token: jwt,
	}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	respondWithJSON(w, http.StatusOK, resp)
}