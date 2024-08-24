package database

import (
	"sort"
	"net/http"
	"errors"
)

type User struct {
	ID 		int `json:"id"`
	Email	string `json:"email"`
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newID := len(dbStructure.Users) + 1
	user := User{ID: newID, Email: email}

	dbStructure.Users[newID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, nil
}

func (db *DB) GetUserByID(id int) (User, int, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, http.StatusInternalServerError, err
	}

	var respUser User
	for _, user := range dbStructure.Users {
		if user.ID == id {
			respUser = user
			return respUser, http.StatusOK, nil
		}
	}

	return User{}, http.StatusNotFound, errors.New("User not found")
}
