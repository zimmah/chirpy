package database

import (
	"errors"
	"sort"
)

type User struct {
	ID 				int `json:"id"`
	Email			string `json:"email"`
	HashedPassword	string `json:"hashed_password"`
	JWT				string `json:"token"`
	RefreshToken	string `json:"refresh_token"`
	ExpiresAt		string	`json:"expires_at"` //refresh token expiry
}

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {
	if _, err := db.GetUserByEmail(email); !errors.Is(err, ErrNotExist) {
		return User{}, ErrAlreadyExists
	}
	
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newID := len(dbStructure.Users) + 1
	user := User{ID: newID, Email: email, HashedPassword: hashedPassword}
	userResp := User{ID: newID, Email: email}

	dbStructure.Users[newID] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return userResp, nil
}

func (db *DB) UpdateUser(id int, email, hashedPassword string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{ID: id, Email: email, HashedPassword: hashedPassword}
	userResp := User{ID: id, Email: email}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return userResp, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil { return nil, err }

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		userResp := User{ID: user.ID, Email: user.Email}
		users = append(users, userResp)
	}

	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })

	return users, nil
}

func (db *DB) GetUserByID(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil { return User{}, err }

	user, ok := dbStructure.Users[id]
	if !ok { return User{}, ErrNotExist }

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil { return User{}, err }

	for _, user := range dbStructure.Users {
		if user.Email == email { return user, nil }
	}

	return User{}, ErrNotExist
}
