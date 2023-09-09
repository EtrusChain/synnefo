package repo

import (
	"github.com/EtrusChain/synnefo/config"
	"github.com/syndtr/goleveldb/leveldb"
)

type Repo interface {
	// Config returns the ipfs configuration file from the repo. Changes made
	// to the returned config are not automatically persisted.
	Config() (*config.Config, error)

	// Datastore returns a reference to the configured data storage backend.
	Datastore() *leveldb.DB
}
