package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"errors"

	"github.com/zimmah/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 				int `json:"id"`
	Email			string `json:"email"`
	Password		string `json:"-"`
}

func databaseUserToSafeUser(databaseUser database.User) User {
	user := User{
		ID: databaseUser.ID,
		Email: databaseUser.Email,
	}

	return user
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	databaseUsers, err := database.DBPointer.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get users")
		return
	}

	var users []User

	for _, databaseUser := range databaseUsers {
		users = append(users, databaseUserToSafeUser(databaseUser))
	}
	respondWithJSON(w, http.StatusOK, users)
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil { 
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	databaseUser, err := database.DBPointer.GetUserByID(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	user := databaseUserToSafeUser(databaseUser)

	respondWithJSON(w, http.StatusOK, user)
}

func handlePostUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password 	string `json:"password"`
		Email 		string `json:"email"`
	}
	
	decoder := json.NewDecoder(r.Body)
	user := parameters{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	databaseUser, err := database.DBPointer.CreateUser(user.Email, string(hashedPassword))
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	responseUser := databaseUserToSafeUser(databaseUser)

	respondWithJSON(w, http.StatusCreated, responseUser)
}

func (cfg *apiConfig) handlePutUsers(w http.ResponseWriter, r *http.Request) {
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

	user, err := database.DBPointer.GetUserByID(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	type parameters struct {
		Email 		string `json:"email"`
		Password 	string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	databaseUser, err := database.DBPointer.UpdateUser(user.ID, userReq.Email, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	responseUser := databaseUserToSafeUser(databaseUser)

	respondWithJSON(w, http.StatusOK, responseUser)
}