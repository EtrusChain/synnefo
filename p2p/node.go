package p2p

import (
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func createNode() (host.Host, error) {
    host, err := libp2p.New()
    if err != nil {
        return nil, err
    }

    return host, nil
}