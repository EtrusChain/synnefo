package peering

import (
	"crypto/sha256"
	"fmt"

	"github.com/EtrusChain/synnefo/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/syndtr/goleveldb/leveldb"
)

func generatePublicKey(privatekey string) string {
	h := sha256.New()

	h.Write([]byte(privatekey))

	bs := h.Sum(nil)

	return string(bs)
}

func CreatePeer() {
	peerId := GetPeer()
	if len(peerId) > 1 {
		panic("Peer alredy created")
	}

	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		return
	}

	defer db.Close()

	batch := new(leveldb.Batch)

	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}
	userData := p2p.New(
		node.ID(),
		node,
		node.Peerstore(),
	)
	// print the node's listening addresses
	fmt.Println(userData)

	pubKey := generatePublicKey(node.ID().Pretty())

	batch.Put([]byte("userID"), []byte(node.ID().Pretty()))
	batch.Put([]byte("pubKey"), []byte(pubKey))

	err = db.Write(batch, nil)
	if err != nil {
		panic("")
	}

	fmt.Println("Your peer address is: %n", node.ID().Pretty())
}

func GetPeer() string {
	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		return ""
	}

	defer db.Close()

	data, err := db.Get([]byte("userID"), nil)
	if err != nil {
		return ""
	}

	dataTwo, err := db.Get([]byte("pubKey"), nil)
	if err != nil {
		return ""
	}

	fmt.Println("Peer: \n", string(data), string(dataTwo))

	return string(data)
}

func RemovePeer() {
	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		return
	}

	defer db.Close()

	db.Delete([]byte("userID"), nil)

	fmt.Println("Peer Romoved")
}
