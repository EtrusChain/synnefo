package bootstrap

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type BootstrapConfig struct {
	// MinPeerThreshold governs whether to bootstrap more connections. If the
	// node has less open connections than this number, it will open connections
	// to the bootstrap nodes. From there, the routing system should be able
	// to use the connections to the bootstrap nodes to connect to even more
	// peers. Routing systems like the IpfsDHT do so in their own Bootstrap
	// process, which issues random queries to find more peers.
	MinPeerThreshold int

	// Period governs the periodic interval at which the node will
	// attempt to bootstrap. The bootstrap process is not very expensive, so
	// this threshold can afford to be small (<=30s).
	Period time.Duration

	// ConnectionTimeout determines how long to wait for a bootstrap
	// connection attempt before cancelling it.
	ConnectionTimeout time.Duration

	// BootstrapPeers is a function that returns a set of bootstrap peers
	// for the bootstrap process to use. This makes it possible for clients
	// to control the peers the process uses at any moment.
	BootstrapPeers func() []peer.AddrInfo

	// BackupBootstrapInterval governs the periodic interval at which the node will
	// attempt to save connected nodes to use as temporary bootstrap peers.
	BackupBootstrapInterval time.Duration

	// MaxBackupBootstrapSize controls the maximum number of peers we're saving
	// as backup bootstrap peers.
	MaxBackupBootstrapSize int

	SaveBackupBootstrapPeers func(context.Context, []peer.AddrInfo)
	LoadBackupBootstrapPeers func(context.Context) []peer.AddrInfo
}

// DefaultBootstrapConfig specifies default sane parameters for bootstrapping.
var DefaultBootstrapConfig = BootstrapConfig{
	MinPeerThreshold:        4,
	Period:                  30 * time.Second,
	ConnectionTimeout:       (30 * time.Second) / 3, // Perod / 3
	BackupBootstrapInterval: 1 * time.Hour,
	MaxBackupBootstrapSize:  20,
}

func BootstrapConfigWithPeers(pis []peer.AddrInfo) BootstrapConfig {
	cfg := DefaultBootstrapConfig
	cfg.BootstrapPeers = func() []peer.AddrInfo {
		return pis
	}
	return cfg
}
