package main

import (
	"github.com/EtrusChain/synnefo/cmd"
)

func main() {
	cmd.Execute()
}

/*
type PeerInfo struct {
	peerID (host.Host)
	userID string

	expiredTime time.Time
	location    time.Location
}
*/

/*

	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		return
	}
	defer db.Close()

	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.ID().Pretty())

	err = db.Put([]byte("userID"), []byte(node.ID().Pretty()), nil)
	if err != nil {
		return
	}

	data, err := db.Get([]byte("userID"), nil)
	fmt.Println(data)
	if err != nil {
		return
	}
*/
