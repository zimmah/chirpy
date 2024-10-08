package database

import (
	"errors"
	"sort"
)

type User struct {
	ExpiresAt		int64	`json:"expires_at"` //refresh token expiry
	ID 				int `json:"id"`
	IsChirpyRed		bool `json:"is_chirpy_red"`
	Email			string `json:"email"`
	HashedPassword	string `json:"hashed_password"`
	RefreshToken	string `json:"refresh_token"`
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
	user := User{ID: newID, Email: email, HashedPassword: hashedPassword, IsChirpyRed: false}
	userResp := User{ID: newID, Email: email, IsChirpyRed: false}

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

	user, ok := dbStructure.Users[id]
	if !ok {
		user = User{ID: id, Email: email, HashedPassword: hashedPassword, IsChirpyRed: false}
	} else {
		user.Email = email
		user.HashedPassword = hashedPassword
	}
	
	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}
	
	userResp := User{ID: id, Email: email, IsChirpyRed: user.IsChirpyRed}
	return userResp, nil
}

func (db *DB) UpdateUserToken(id int, tokenExpiry int64, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil { return err }

	user := dbStructure.Users[id]
	updatedUser := User{
		ID: user.ID,
		Email: user.Email,
		HashedPassword: user.HashedPassword,
		RefreshToken: token,
		ExpiresAt: tokenExpiry,
		IsChirpyRed: user.IsChirpyRed,
	}

	dbStructure.Users[id] = updatedUser

	return db.writeDB(dbStructure)
}

func (db *DB) UpgradeUser(userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil { return err }
	
	_, err = db.GetUserByID(userID)
	if err != nil {
		return err
	}

	user := dbStructure.Users[userID]
	user.IsChirpyRed = true
	dbStructure.Users[userID] = user

	return db.writeDB(dbStructure)
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil { return nil, err }

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		userResp := User{ID: user.ID, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
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

func (db *DB) GetUserByRefreshToken(refreshToken string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil { return User{}, err }

	for _, user := range dbStructure.Users {
		if user.RefreshToken == refreshToken { return user, nil }
	}

	return User{}, ErrNotExist
}