package database

import (
	"sort"
	"net/http"
	"errors"
)

type User struct {
	ID 					int `json:"id"`
	Email				string `json:"email"`
	Password 			string `json:"password"`
	ExpiresInSeconds 	*int `json:"expires_in_seconds"`
}

type UserResponse struct {
	ID			int `json:"id"`
	Email		string `json:"email"`
}

type UserResponseWithToken struct {
	ID			int `json:"id"`
	Email		string `json:"email"`
	Token		string `json:"token"`
}

func (db *DB) CreateUser(email, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	newID := len(dbStructure.Users) + 1
	user := User{ID: newID, Email: email, Password: password}
	userResp := UserResponse{ID: newID, Email: email}

	dbStructure.Users[newID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return userResp, nil
}

func (db *DB) UpdateUser(id int, email, password string) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	user := User{ID: id, Email: email, Password: password}
	userResp := UserResponse{ID: id, Email: email}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return userResp, nil
}

func (db *DB) GetUsers() ([]UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]UserResponse, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		userResp := UserResponse{ID: user.ID, Email: user.Email}
		users = append(users, userResp)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, nil
}

func (db *DB) GetUserByID(id int) (UserResponse, int, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, http.StatusInternalServerError, err
	}

	for _, user := range dbStructure.Users {
		if user.ID == id {
			respUser := UserResponse{ID: user.ID, Email: user.Email}
			return respUser, http.StatusOK, nil
		}
	}

	return UserResponse{}, http.StatusNotFound, errors.New("User not found")
}

func (db *DB) GetUserByEmail(email string) (User, int, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, http.StatusInternalServerError, err
	}

	var respUser User
	for _, user := range dbStructure.Users {
		if user.Email == email {
			respUser = user
			return respUser, http.StatusOK, nil
		}
	}

	return User{}, http.StatusUnauthorized, errors.New("Unauthorized: Invalid username/password combination")
}
