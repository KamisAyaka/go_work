package data

import (
	"Go-Minichain/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"strconv"
	"strings"
	"time"
)

/**
 * 对交易的抽象
 */

type Transaction struct {
	//data          string
	timestamp     int
	inUTXO        []*UTXO
	outUTXO       []*UTXO
	sendSign      []byte
	sendPublicKey ecdsa.PublicKey
}

func NewTransaction(inUTXO []*UTXO, outUTXO []*UTXO, sendSign []byte, sendPublicKey ecdsa.PublicKey) *Transaction {
	return &Transaction{
		inUTXO:        inUTXO,
		outUTXO:       outUTXO,
		sendSign:      sendSign,
		sendPublicKey: sendPublicKey,
		timestamp:     time.Now().Nanosecond(),
	}
}

func (t *Transaction) GetInUTXOs() []*UTXO {
	return t.inUTXO
}

func (t *Transaction) GetOutUTXOs() []*UTXO {
	return t.outUTXO
}

func (t *Transaction) GeSendSign() []byte {
	return t.sendSign
}

func (t *Transaction) GetSendPublicKey() ecdsa.PublicKey {
	return t.sendPublicKey
}
func (t *Transaction) GetTimeStamp() int {
	return t.timestamp
}

func (t *Transaction) ToString() string {
	publicKey := t.sendPublicKey
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)

	inUTXOStrings := make([]string, len(t.inUTXO))
	for i, iu := range t.inUTXO {
		inUTXOStrings[i] = iu.ToString()
	}
	outUTXOStrings := make([]string, len(t.outUTXO))
	for i, ou := range t.outUTXO {
		outUTXOStrings[i] = ou.ToString()
	}
	return "Transaction{" +
		"inUTXO=" + strings.Join(inUTXOStrings, "\n") +
		", outUTXO=" + strings.Join(outUTXOStrings, "\n") +
		", sendSign=" + utils.Byte2HexString(t.sendSign) +
		", sendPublicKey=" + utils.Byte2HexString(publicKeyBytes) +
		", timestamp=" + strconv.Itoa(t.timestamp) +
		"}"
}
