package database

import (
	"encoding/json"
	"os"
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