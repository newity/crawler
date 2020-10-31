package parser

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/newity/crawler/blocklib"
)

type Data struct {
	BlockNumber     uint64
	Prevhash        []byte
	Datahash        []byte
	BlockSignatures []blocklib.BlockSignature
	Txs             []blocklib.Tx
	Events          []*peer.ChaincodeEvent
}