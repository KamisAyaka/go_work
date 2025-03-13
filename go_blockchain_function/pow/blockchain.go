package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockchainIterator struct {
	currentHash []byte   // 当前区块hash
	db          *bolt.DB // 打开的数据库
}

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) MinedBlock(txs []*Transaction, miner, data string) {
	var tip []byte
	// 得到最新的哈希值
	bc.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		tip = buck.Get([]byte("last"))
		return nil
	})
	// 更新数据库
	bc.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))

		coinbasetx := NewCoinbaseTX(miner, data)
		txs = append(txs, coinbasetx)
		block := NewBlock(txs, tip)

		buck.Put(block.Hash, block.Serialize())
		buck.Put([]byte("last"), block.Hash)

		bc.tip = block.Hash
		return nil
	})
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, _ := bolt.Open(dbFile, 0600, nil)
	db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		// 如果为空，创建创世块
		if buck == nil {
			fmt.Println("No existing blockchain found, creating a new one....")
			coinbasetx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(coinbasetx)
			block_data := genesis.Serialize()
			bucket, _ := tx.CreateBucket([]byte(blocksBucket))
			bucket.Put(genesis.Hash, block_data)
			bucket.Put([]byte("last"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = buck.Get([]byte("last"))
		}
		return nil
	})
	return &Blockchain{tip, db}
}

func (i *BlockchainIterator) PreBlock() (*Block, bool) {
	var block *Block

	i.db.View(func(tx *bolt.Tx) error {
		buck := tx.Bucket([]byte(blocksBucket))
		encodedBlock := buck.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	i.currentHash = block.PrevHash
	return block, len(i.currentHash) > 0
}

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var UTXO []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block, next := bci.PreBlock()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					UTXO = append(UTXO, *tx)
				}
			}
			// 处理所有类型的交易，包括硬币基础交易
			for _, in := range tx.Vin {
				if in.UsesKey(pubKeyHash) {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.VoutIdx)
				}
			}
		}
		if !next {
			break
		}
	}
	return UTXO
}

func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

func (bc *Blockchain) getBalance(address string) {
	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
func (bc *Blockchain) Send(from, to string, amount int, data string, wallet *Wallet) {
	fmt.Println("send address from ...", from)
	tx := NewUTXOTransaction(from, to, amount, bc, wallet)
	bc.MinedBlock([]*Transaction{tx}, from, data)
	fmt.Println("send success")
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block, next := bci.PreBlock()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if !next {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}
