package database

import(
	"sort"
)

type Chirp struct {
	ID 			int `json:"id"`
	AuthorID 	int `json:"author_id"`
	Body 		string `json:"body"`
}

func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{ID: newID, AuthorID: authorID, Body: body}

	dbStructure.Chirps[newID] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)
	return db.writeDB(dbStructure)
}

func (db *DB) GetChirps(sortOrder string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "asc" {
			return chirps[i].ID < chirps[j].ID
		} else {
			return chirps[j].ID < chirps[i].ID
		}
	})

	return chirps, nil
}

func (db *DB) GetChirpsOfAuthor(authorID int, sortOrder string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "asc" {
			return chirps[i].ID < chirps[j].ID
		} else {
			return chirps[j].ID < chirps[i].ID
		}
	})

	return chirps, nil
}

func (db *DB) GetChirpByID(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}