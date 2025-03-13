package main

import (
	"bytes"
)

type TXInput struct {
	Txid      []byte
	VoutIdx   int
	Signature []byte
	PubKey    []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {

	lockingHash := HashPubKey(in.PubKey)
	//fmt.Printf("UsesKey:%x,\n%x\n", pubKeyHash, lockingHash)
	return bytes.Equal(lockingHash, pubKeyHash)
}
