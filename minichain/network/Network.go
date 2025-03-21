package network

import (
	"Go-Minichain/config"
	"Go-Minichain/data"
	"Go-Minichain/spv"
	"fmt"
)

// NetWork 定义了一个区块链网络的结构体。
// 字段说明：
// - accounts: 存储所有账户信息的列表。
// - txPool: 交易池，用于存储待处理的交易。
// - blockchain: 区块链对象，用于管理区块和 UTXO。
// - miner: 矿工节点，负责挖矿和生成新区块。
// - spvPeer: SPV 节点列表，用于轻量级客户端验证。
type NetWork struct {
	accounts   []data.Account
	txPool     *TransactionPool
	blockchain *BlockChain
	miner      MinerNode
	spvPeer    []*SPVPeer
}

// NewNetWork 创建一个新的区块链网络实例。
// 返回值:
// 返回一个指向新创建的区块链网络实例的指针。
func NewNetWork() *NetWork {
	network := new(NetWork)
	fmt.Println("Accounts and SPVPeers config...")
	accounts := make([]data.Account, config.MiniChainConfig.GetAccountNumber())
	peers := make([]*SPVPeer, config.MiniChainConfig.GetAccountNumber())
	for i := range accounts {
		accounts[i] = *data.NewAccount()
		peers[i] = NewSPVPeer(accounts[i], network)
	}
	network.accounts = accounts
	network.spvPeer = peers
	fmt.Println("TransactionPool config...")
	pool := NewTransactionPool(config.MiniChainConfig.GetMaxTransactionCount(), network)
	fmt.Println("Blockchain config...")
	blockchain := NewBlockChain(network)
	fmt.Println("MinerNode config...")
	miner := NewMinerNode(network)
	fmt.Println("Network Config Finished...")
	fmt.Println("Network Start...")
	network.txPool = pool
	network.blockchain = blockchain
	network.miner = *miner
	return network
}

// Start 启动区块链网络。
// 该方法会初始化创世块、广播创世块，并启动交易池和矿工节点。
func (n *NetWork) Start() {
	n.blockchain.SetUp()
	n.BroadCast(*n.GetNewestBlock())
	n.txPool.Start()
	n.miner.Run()
}

// GetTransactionsInLatestBlock 获取最新区块中与指定钱包地址相关的所有交易。
// 参数:
// - address: 钱包地址。
// 返回值:
// 返回一个包含相关交易的列表。
func (n *NetWork) GetTransactionsInLatestBlock(address string) []data.Transaction {
	txs := make([]data.Transaction, 0)
	block := n.GetNewestBlock()
	blockBody := block.GetBlockBody()
	for _, tx := range blockBody.GetTransctions() {
		have := false
		for _, utxo := range tx.GetInUTXOs() {
			if utxo.GetWalletAddress() == address {
				txs = append(txs, tx)
				have = true
				break
			}
		}
		if have {
			continue
		}
		for _, utxo := range tx.GetOutUTXOs() {
			if utxo.GetWalletAddress() == address {
				txs = append(txs, tx)
				break
			}
		}
	}
	return txs
}

// CheckTransactionIsFull 检查交易池是否已满。
// 返回值:
// 返回布尔值，表示交易池是否已满。
func (n *NetWork) CheckTransactionIsFull() bool {
	return n.txPool.IsFull()
}

// GetAllTransactions 获取交易池中的所有交易。
// 返回值:
// 返回一个包含所有交易的列表。
func (n *NetWork) GetAllTransactions() []data.Transaction {
	return n.txPool.GetAll()
}

// GetTotalAmount 获取区块链中所有账户的总金额。
// 返回值:
// 返回整数类型的总金额。
func (n *NetWork) GetTotalAmount() int {
	return n.blockchain.GetAllAmount()
}

// AddNewBlock 将新区块添加到区块链中。
// 参数:
// - block: 要添加的新区块。
func (n *NetWork) AddNewBlock(block data.Block) {
	n.blockchain.AddNewBlock(block)
}

// GetNewestBlock 获取区块链中的最新区块。
// 返回值:
// 返回指向最新区块的指针。
func (n *NetWork) GetNewestBlock() *data.Block {
	return n.blockchain.GetNewestBlock()
}

// GetAccounts 获取所有账户信息。
// 返回值:
// 返回一个包含所有账户的列表。
func (n *NetWork) GetAccounts() []data.Account {
	return n.accounts
}

// GetAccount 根据索引获取指定账户信息。
// 参数:
// - i: 账户索引。
// 返回值:
// 返回指定索引的账户信息。
func (n *NetWork) GetAccount(i int) data.Account {
	return n.accounts[i]
}

// GetTrueUTXOs 获取指定钱包地址的有效 UTXO 列表。
// 参数:
// - address: 钱包地址。
// 返回值:
// 返回一个包含有效 UTXO 的列表。
func (n *NetWork) GetTrueUTXOs(address string) []*data.UTXO {
	return n.blockchain.GetTrueUTXOs(address)
}

// ProcessTransactionUTXO 处理交易中的未花费交易输出（UTXO）。
// 参数:
// - inUTXO: 表示交易中作为输入的 UTXO 列表。
// - outUTXO: 表示交易中作为输出的 UTXO 列表。
func (n *NetWork) ProcessTransactionUTXO(inUTXO []*data.UTXO, outUTXO []*data.UTXO) {
	n.blockchain.ProcessTransactionUTXO(inUTXO, outUTXO)
}

// GetBlocks 获取区块链中的所有区块。
// 返回值:
// 返回一个包含所有区块的列表。
func (n *NetWork) GetBlocks() []data.Block {
	return n.blockchain.GetBlocks()
}

// GetSPVPeers 获取所有 SPV 节点。
// 返回值:
// 返回一个包含所有 SPV 节点的列表。
func (n *NetWork) GetSPVPeers() []*SPVPeer {
	return n.spvPeer
}

// BroadCast 广播新区块到所有 SPV 节点。
// 参数:
// - block: 要广播的新区块。
func (n *NetWork) BroadCast(block data.Block) {
	n.miner.BroadCast(block)
}

// GetProof 根据交易哈希生成一个简化的支付验证（SPV）证明。
// 参数:
// - hash: 交易的哈希值。
// 返回值:
// 返回一个 spv.Proof 类型的证明对象。
func (n *NetWork) GetProof(hash string) spv.Proof {
	return n.miner.GetProof(hash)
}

// GetBlockchain 获取区块链对象。
// 返回值:
// 返回指向区块链对象的指针。
func (n *NetWork) GetBlockchain() *BlockChain {
	return n.blockchain
}
