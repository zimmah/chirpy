package database

import(
	"errors"
	"net/http"
	"sort"
)

type Chirp struct {
	ID 		int `json:"id"`
	Body 	string `json:"body"`
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