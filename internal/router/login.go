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
	}

	type response struct {
		User
		Token 			string `json:"token"`
		RefreshToken 	string `json:"refresh_token"`
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
	
	jwt, err := cfg.generateJWT(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT")
		return
	}

	refreshToken, err := cfg.generateRefreshToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token")
		return
	}

	resp := response{
		User: User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed},
		Token: jwt,
		RefreshToken: refreshToken,
	}

	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	respondWithJSON(w, http.StatusOK, resp)
}