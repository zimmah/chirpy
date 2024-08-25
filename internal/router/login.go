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

	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Can not decode user: %v", err))
		return
	}
	
	user, statusCode, err := database.DBPointer.GetUserByEmail(userReq.Email)
	if err != nil {
		respondWithError(w, statusCode, fmt.Sprint(err))
		return
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: Invalid username/password combination")
		return
	}
	
	jwt, err := cfg.generateJWT(user, userReq.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	userResp := database.User{ID: user.ID, Email: user.Email, JWT: jwt}
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	respondWithJSON(w, http.StatusOK, userResp)
}