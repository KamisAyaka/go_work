package network

import (
	"Go-Minichain/config"
	"Go-Minichain/data"
	"Go-Minichain/utils"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
)

// BlockChain 定义了一个区块链的结构体。
// 字段说明：
// - chain: 存储区块链中的所有区块。
// - network: 网络对象，用于与网络交互。
// - UTXOs: 存储当前未花费的交易输出（UTXO）列表。
// - mutex: 用于保护并发访问的互斥锁。
type BlockChain struct {
	chain   []data.Block
	network *NetWork
	UTXOs   []*data.UTXO
	mutex   sync.Mutex
}

// NewBlockChain 创建一个新的区块链实例。
// 参数:
// - network: 网络对象，用于初始化区块链。
// 返回值:
// 返回一个指向新创建的区块链实例的指针。
func NewBlockChain(network *NetWork) *BlockChain {
	chain := new(BlockChain)
	chain.chain = make([]data.Block, 0)
	chain.UTXOs = make([]*data.UTXO, 0)
	chain.network = network
	return chain
}

// SetUp 初始化区块链，生成创世块并加入区块链中。
func (c *BlockChain) SetUp() {
	transactions := c.GenesisTransactions()
	header := data.NewBlockHeader("", "", rand.Int63())
	body := c.network.miner.GetBlockBody(transactions)
	genesisBlock := data.NewBlock(*header, body)
	fmt.Println("Create the genesis Block! ")
	fmt.Println("And the hash of genesis Block is : " + utils.GetSha256Digest(genesisBlock.ToString()) +
		", you will see the hash value in next Block's preBlockHash field.")
	fmt.Println()
	c.AddNewBlock(*genesisBlock)
}

// AddNewBlock 将新区块添加到区块链中。
// 参数:
// - block: 要添加的新区块。
func (c *BlockChain) AddNewBlock(block data.Block) {
	c.chain = append(c.chain, block)
}

// GetNewestBlock 获取区块链中的最新区块。
// 返回值:
// 返回指向最新区块的指针。
func (c *BlockChain) GetNewestBlock() *data.Block {
	return &c.chain[len(c.chain)-1]
}

// GenesisTransactions 生成创世块的初始交易。
// 返回值:
// 返回包含创世交易的列表。
func (c *BlockChain) GenesisTransactions() []data.Transaction {
	outUTXOs := make([]*data.UTXO, len(c.network.GetAccounts()))
	for i := 0; i < len(outUTXOs); i++ {
		account := c.network.GetAccount(i)
		outUTXOs[i] = data.NewUTXO(account.GetWalletAddress(), config.MiniChainConfig.GetInitAmount(), account.GetPublicKey())
	}
	c.ProcessTransactionUTXO([]*data.UTXO{}, outUTXOs)
	daydreamPrivateKey, daydreamPublicKey := utils.Secp256k1Generate()
	sign := utils.Signature([]byte("Wecome to Blockchain Lab!!!"), daydreamPrivateKey)
	return []data.Transaction{*data.NewTransaction(make([]*data.UTXO, 0), outUTXOs, sign, daydreamPublicKey)}
}

// ProcessTransactionUTXO 处理交易中的未花费交易输出（UTXO）。
// 参数:
// - inUTXO: 表示交易中作为输入的 UTXO 列表，这些 UTXO 将被标记为已使用。
// - outUTXO: 表示交易中作为输出的 UTXO 列表，这些 UTXO 将被添加到区块链中。
func (c *BlockChain) ProcessTransactionUTXO(inUTXO []*data.UTXO, outUTXO []*data.UTXO) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, utxo := range inUTXO {
		utxo.SetUsed()
	}

	for _, utxo := range outUTXO {
		c.AddUTXO(utxo)
	}
}

// AddUTXO 将新的 UTXO 添加到区块链中。
// 参数:
// - u: 要添加的 UTXO。
func (c *BlockChain) AddUTXO(u *data.UTXO) {
	c.UTXOs = append(c.UTXOs, u)
}

// GetTrueUTXOs 获取指定钱包地址的有效 UTXO 列表。
// 参数:
// - walletAddress: 钱包地址。
// 返回值:
// 返回该钱包地址对应的所有有效 UTXO 列表。
func (c *BlockChain) GetTrueUTXOs(walletAddress string) []*data.UTXO {
	trueUTXOs := make([]*data.UTXO, 0)
	for _, utxo := range c.UTXOs {
		if utxo.GetWalletAddress() == walletAddress && !utxo.IsUsed() {
			trueUTXOs = append(trueUTXOs, utxo)
		}
	}
	return trueUTXOs
}

// GetAllAmount 计算区块链中所有账户的总金额，并验证余额是否正确。
// 返回值:
// 返回区块链中所有账户的总金额。
func (c *BlockChain) GetAllAmount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sumAccount := 0

	for _, account := range c.network.accounts {
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

// GetBlocks 获取区块链中的所有区块。
// 返回值:
// 返回存储在区块链中的所有区块。
func (c *BlockChain) GetBlocks() []data.Block {
	return c.chain
}
