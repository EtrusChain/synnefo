package repo

import "github.com/EtrusChain/synnefo/config"

// Mock is not thread-safe.
type Mock struct {
	C config.Config
}

func (m *Mock) Config() (*config.Config, error) {
	return &m.C, nil // FIXME threadsafety
}
