package database

func newDB(path string) (*DB, error) {

}

func (db *DB) ensureDB() error {

}

func (db *DB) loadDB() (DBStructure, error) {

}

func (db *DB) writeDB(dbStructure DBStructure) error {

}

type DB struct {
	path 	string
	mux 	*sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id 		int `json:"id"`
	Body 	string `json:"body"`
}