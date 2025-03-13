package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const subsidy = 1000

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// 将交易ID进行哈希
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(tx)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].VoutIdx == -1
}

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s", to)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain, wallet *Wallet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(HashPubKey(wallet.PublicKey), amount)
	if acc < amount {
		fmt.Println("ERROR: Not enough funds")
		os.Exit(1)
	}
	for txid, outs := range validOutputs {
		txID, _ := hex.DecodeString(txid)

		for _, out := range outs {
			input := TXInput{txID, out, nil, wallet.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTXOutput(amount, to))

	if acc > amount {
		// 找零
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	bc.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for inID, vin := range txCopy.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.VoutIdx].PubKeyhash
		txCopy.SetID()

		txCopy.Vin[inID].PubKey = nil
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.VoutIdx, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyhash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
