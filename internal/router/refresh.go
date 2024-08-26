package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/zimmah/chirpy/internal/database"
)

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := database.DBPointer.GetUserByRefreshToken(tokenString)
	if err != nil || time.Unix(int64(user.ExpiresAt), 0).Before(time.Now())  {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	token, err := cfg.generateJWT(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT")
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, resp)
}