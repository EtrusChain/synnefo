package repo

import (
	"runtime"

	"github.com/syndtr/goleveldb/leveldb"
)

// DatabaseHandler is a struct that holds a reference to the LevelDB database.
type DatabaseHandler struct {
	db *leveldb.DB
}

// NewDatabaseHandler initializes and returns a new DatabaseHandler.
func NewDatabaseHandler(dbPath string) (*DatabaseHandler, error) {
	db, err := dataStore(dbPath)
	if err != nil {
		return nil, err
	}
	return &DatabaseHandler{db: db}, nil
}

// Close closes the LevelDB database.
func (dh *DatabaseHandler) Close() {
	dh.db.Close()
}

// GetValue retrieves a value from the database by key.
func (dh *DatabaseHandler) GetValue(key string) ([]byte, error) {
	return dh.db.Get([]byte(key), nil)
}

// SetValue stores a key-value pair in the database.
func (dh *DatabaseHandler) SetValue(key string, value []byte) error {
	return dh.db.Put([]byte(key), value, nil)
}

// SetValue stores a key-value pair in the database.
func (dh *DatabaseHandler) CheckValue(key string) (bool, error) {
	return dh.db.Has([]byte(key), nil)
}

// openDatabaseFile opens the LevelDB database file.
func dataStore(dbPath string) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetOs() string {
	osName := runtime.GOOS
	if osName == "windows" {
		return "C:\\ProgramData\\EtrusChain"
	} else if osName == "linux" {
		return "~/.config/EtrusChain"
	} else {
		return "C:\\ProgramData\\EtrusChain"
	}
}
