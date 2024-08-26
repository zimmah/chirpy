package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/zimmah/chirpy/internal/database"
)

func (cfg *apiConfig) handleRevoke (w http.ResponseWriter, r *http.Request) {
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

	err = database.DBPointer.UpdateUserToken(user.ID, 0, "")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete token")
	}

	w.WriteHeader(http.StatusNoContent)
}