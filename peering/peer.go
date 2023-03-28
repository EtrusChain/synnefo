package peering

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"

	"github.com/EtrusChain/synnefo/p2p"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/syndtr/goleveldb/leveldb"
)

type Peer struct {
	Protocols       []protocol.ID
	AgentVersion    string
	ProtocolVersion string

	identity  peer.ID
	pubKey    string
	addresses []string
}

func generatePublicKey(privatekey string) string {
	h := sha256.New()

	h.Write([]byte(privatekey))

	bs := h.Sum(nil)

	return string(bs)
}

func CreatePeer() *Peer {
	peerId := GetPeer()
	if len(peerId) > 1 {
		panic("Peer alredy created")
	}

	db, err := leveldb.OpenFile("db/userid", nil)
	if err != nil {
		panic("")
	}

	defer db.Close()

	batch := new(leveldb.Batch)

	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)
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

	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	fmt.Println(enc, dec)

	protos := node.Mux().Protocols()

	peerData := &Peer{
		identity: node.ID(),
		pubKey:   pubKey,

		addresses:       nil,
		Protocols:       protos,
		AgentVersion:    "1.0.0",
		ProtocolVersion: "1.0.0",
	}

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(peerData)

	reqBodyBytes.Bytes() // this is the []byte

	batch.Put([]byte("peerID"), reqBodyBytes.Bytes())
	batch.Put([]byte("userID"), []byte(node.ID().Pretty()))
	batch.Put([]byte("pubKey"), []byte(pubKey))

	err = db.Write(batch, nil)
	if err != nil {
		panic("")
	}

	fmt.Println("Your peer address is: %n", node.ID().Pretty())

	return peerData
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

	dataTwo, err := db.Get([]byte("peerID"), nil)
	if err != nil {
		return ""
	}

	fmt.Println("Peer: \n", string(dataTwo))

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
