package data

import (
	"Go-Minichain/utils"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"strconv"
)

type UTXO struct {
	walletAddress string
	amount        int
	publicKeyHash []byte
	used          bool
}

func NewUTXO(address string, amount int, publicKey ecdsa.PublicKey) *UTXO {
	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	return &UTXO{
		walletAddress: address,
		amount:        amount,
		publicKeyHash: utils.Ripemd160Digest(utils.Sha256Digest(publicKeyBytes)),
		used:          false,
	}
}

func (utxo *UTXO) UnlockScript(sign []byte, publicKey ecdsa.PublicKey) bool {
	stack := make([][]byte, 0)

	stack = append(stack, sign)

	publicKeyBytes := elliptic.Marshal(publicKey, publicKey.X, publicKey.Y)
	stack = append(stack, publicKeyBytes)

	stack = append(stack, stack[len(stack)-1])

	data := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	stack = append(stack, utils.Ripemd160Digest(utils.Sha256Digest(data)))

	stack = append(stack, utxo.publicKeyHash)

	publicKeyHash1 := stack[len(stack)-1]
	publicKeyHash2 := stack[len(stack)-2]
	stack = stack[:len(stack)-2]
	if !bytes.Equal(publicKeyHash1, publicKeyHash2) {
		return false
	}

	publicKeyEncoded := stack[len(stack)-1]
	sign1 := stack[len(stack)-2]
	//stack = stack[:len(stack)-2]

	return utils.Verify(publicKeyEncoded, sign1, &publicKey)
}

func (utxo *UTXO) SetUsed() {
	if utxo.used {
		panic("UTXO already used")
	}
	utxo.used = true
}
func (utxo *UTXO) IsUsed() bool {
	return utxo.used
}

func (utxo *UTXO) GetWalletAddress() string {
	return utxo.walletAddress
}

func (utxo *UTXO) GetAmount() int {
	return utxo.amount
}

func (utxo *UTXO) GetPublicKeyHash() []byte {
	return utxo.publicKeyHash
}

func (utxo *UTXO) ToString() string {
	return "UTXO{" +
		"walletAddress=" + utxo.walletAddress + "," +
		"amount=" + strconv.Itoa(utxo.amount) + "," +
		"publicKeyHash=" + utils.Byte2HexString(utxo.publicKeyHash) +
		"}"
}

func UTXO2Bytes(inUTXO []*UTXO, outUTXO []*UTXO) []byte {
	return []byte(fmt.Sprint(inUTXO, outUTXO))
}
