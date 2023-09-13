package repo

import (
	"github.com/EtrusChain/synnefo/config"
)

type Repo interface {
	// Config returns the syneefo configuration file from the repo. Changes made
	// to the returned config are not automatically persisted.
	Config() (*config.Config, error)
}
