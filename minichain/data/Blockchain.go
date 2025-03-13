package data

import (
	"Go-Minichain/config"
	"Go-Minichain/utils"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

/**
 * 区块链的类抽象，创建该对象时会自动生成创世纪块，加入区块链中
 */

type BlockChain struct {
	chain    []Block
	accounts []Account
	UTXOs    []*UTXO
	mutex    sync.Mutex
}

func NewBlockChain() *BlockChain {
	chain := new(BlockChain)
	chain.chain = make([]Block, 0)
	chain.accounts = make([]Account, config.MiniChainConfig.GetAccountNumber())
	for i := range chain.accounts {
		chain.accounts[i] = *NewAccount()
	}
	transactions := chain.GenesisTransaction()

	header := NewBlockHeader("", "", rand.Int63())
	body := NewBlockBody("", transactions)
	genesisBlock := NewBlock(*header, *body)
	fmt.Println("Create the genesis Block! ")
	fmt.Println("And the hash of genesis Block is : " + utils.GetSha256Digest(genesisBlock.ToString()) +
		", you will see the hash value in next Block's preBlockHash field.")
	fmt.Println()
	chain.AddNewBlock(*genesisBlock)
	return chain
}
func (c *BlockChain) AddNewBlock(block Block) {
	c.chain = append(c.chain, block)
}
func (c *BlockChain) GetNewestBlock() *Block {
	return &c.chain[len(c.chain)-1]
}

func (c *BlockChain) GenesisTransaction() []Transaction {
	outUTXOs := make([]*UTXO, len(c.accounts))
	for i := 0; i < len(outUTXOs); i++ {
		outUTXOs[i] = NewUTXO(c.accounts[i].GetWalletAddress(), config.MiniChainConfig.GetInitAmount(), c.accounts[i].GetPublicKey())
	}
	c.ProcessTransactionUTXO([]*UTXO{}, outUTXOs)
	daydreamPrivateKey, daydreamPublicKey := utils.Secp256k1Generate()
	sign := utils.Signature([]byte("I am the creator of this blockchain"), daydreamPrivateKey)
	return []Transaction{
		*NewTransaction(make([]*UTXO, 0), outUTXOs, sign, daydreamPublicKey)}
}

func (c *BlockChain) ProcessTransactionUTXO(inUTXO []*UTXO, outUTXO []*UTXO) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, utxo := range inUTXO {
		utxo.SetUsed()
	}
	for _, utxo := range outUTXO {
		c.AddUTXO(utxo)
	}
}

func (c *BlockChain) AddUTXO(utxo *UTXO) {
	c.UTXOs = append(c.UTXOs, utxo)
}

func (c *BlockChain) GetTrueUTXOs(walletAddress string) []*UTXO {
	trueUTXOs := make([]*UTXO, 0)
	for _, utxo := range c.UTXOs {
		if utxo.GetWalletAddress() == walletAddress && !utxo.IsUsed() {
			trueUTXOs = append(trueUTXOs, utxo)
		}
	}
	return trueUTXOs
}

func (c *BlockChain) GetAccount() []Account {
	return c.accounts
}

func (c *BlockChain) GetAllAmount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sumAccount := 0

	for _, account := range c.accounts {
		utxos := c.GetTrueUTXOs(account.GetWalletAddress())
		for _, utxo := range utxos {
			if utxo.IsUsed() {
				panic("error")
			}
		}
		sumAccount += account.GetAmount(utxos)
	}
	if sumAccount != config.MiniChainConfig.GetAccountNumber()*config.MiniChainConfig.GetInitAmount() {
		panic("error Balance:" + strconv.Itoa(sumAccount))
	}
	return sumAccount
}
