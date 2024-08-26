package router

import (
	"encoding/json"
	"net/http"

	"github.com/zimmah/chirpy/internal/database"
)

func (cfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	AuthHeader := r.Header.Get("Authorization")
	if len(AuthHeader) < len("Apikey ") + 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key := AuthHeader[len("Apikey "):]
	if key != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type Data struct {
		UserID int `json:"user_id"`
	}
	
	type Parameters struct {
		Event string `json:"event"`
		Data Data `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	UserID := params.Data.UserID

	err = database.DBPointer.UpgradeUser(UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}