package database

import "github.com/syndtr/goleveldb/leveldb"

func Database(databasePath string, key string) string {
	db, err := leveldb.OpenFile(databasePath, nil)
	if err != nil {
		return ""
	}

	defer db.Close()

	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return ""
	}

	return string(data)
}
