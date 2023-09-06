package core

import (
	"context"
	"io"

	"github.com/EtrusChain/synnefo/peering"
	ic "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	madns "github.com/multiformats/go-multiaddr-dns"
)

// synnefoNode is IPFS Core module. It represents an IPFS instance.
type SynnefoNode struct {
	// Self
	Identity peer.ID // the local node's identity

	// Services
	Reporter *metrics.BandwidthCounter `optional:"true"`
	// Discovery            mdns.Service              `optional:"true"`
	PrivateKey ic.PrivKey `optional:"true"` // the local node's private Key

	// Online
	Peering         *peering.PeeringService `optional:"true"`
	Bootstrapper    io.Closer               `optional:"true"` // the periodic bootstrapper
	DNSResolver     *madns.Resolver         // the DNS resolver
	ResourceManager network.ResourceManager `optional:"true"`

	DHTClient routing.Routing `name:"dhtc" optional:"true"`

	// P2P *p2p.P2P `optional:"true"`

	ctx context.Context

	stop func() error

	// Flags
	IsOnline bool `optional:"true"` // Online is set when networking is enabled.
	IsDaemon bool `optional:"true"` // Daemon is set when running on a long-running daemon.
}

// Close calls Close() on the App object
func (n *SynnefoNode) Close() error {
	return n.stop()
}

func (n *SynnefoNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}
