package database

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sort"
	"sync"
)

var DBPointer *DB

func NewDB(path string) (*DB, error) {
	var mutex sync.RWMutex

	db := &DB{
		path: path,
		mux: &mutex,
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	DBPointer = db

	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{ID: newID, Body: body}

	dbStructure.Chirps[newID] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, int, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, http.StatusInternalServerError, err
	}

	var respChirp Chirp
	for _, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			respChirp = chirp
			return respChirp, http.StatusOK, nil
		}
	}

	return Chirp{}, http.StatusNotFound, errors.New("Chirp not found")
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

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()

		dbStructure := DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User)}
		return db.writeDB(dbStructure)
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()

	var dbStructure DBStructure
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	
	file, err := os.Create(db.path)
	if err != nil { return err }
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(dbStructure)
	if err != nil { return err }

	return nil
}

type DB struct {
	path 	string
	mux 	*sync.RWMutex
}

type DBStructure struct {
	Chirps 	map[int]Chirp `json:"chirps"`
	Users 	map[int]User `json:"users"`
}

type Chirp struct {
	ID 		int `json:"id"`
	Body 	string `json:"body"`
}

type User struct {
	ID 		int `json:"id"`
	Email	string `json:"email"`
}