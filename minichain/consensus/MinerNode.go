package consensus

import (
	"Go-Minichain/config"
	"Go-Minichain/data"
	"Go-Minichain/utils"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

/**
 * 矿工线程
 *
 * 该线程的主要工作就是不断的进行交易打包、Merkle树根哈希值计算、构造区块，
 * 然后尝试使用不同的随机字段（nonce）进行区块的哈希值计算以生成新的区块添加到区块中
 *
 */

type MinerNode struct {
	transactionPool *data.TransactionPool
	blockchain      *data.BlockChain
}

func NewMinerNode(pool *data.TransactionPool, chain *data.BlockChain) *MinerNode {
	return &MinerNode{transactionPool: pool, blockchain: chain}
}
func (m *MinerNode) Run() {
	m.transactionPool.Start()
	for i := 0; i < 3; {
		if m.transactionPool.IsFull() {
			transactions := m.transactionPool.GetAll()
			if !m.Check(transactions) {
				fmt.Println("transactions is not valid")
				os.Exit(0)
			}
			blockBody := m.GetBlockBody(transactions)
			m.Mine(blockBody)
			i++
			fmt.Println("The sum of all amount", m.blockchain.GetAllAmount())
		}
	}
}
func (m *MinerNode) GetBlockBody(transactions []data.Transaction) data.BlockBody {
	if transactions == nil || len(transactions) > config.MiniChainConfig.GetMaxTransactionCount() {
		panic("transactions can not be nil or be more than config.MaxTransactionCount")
	}
	var hashes []string
	for _, tx := range transactions {
		hashes = append(hashes, utils.GetSha256Digest(tx.ToString()))
	}

	for len(hashes) > 1 {
		var newLevel []string
		for i := 0; i < len(hashes); i += 2 {
			if i+1 == len(hashes) {
				newLevel = append(newLevel, utils.GetSha256Digest(hashes[i]+hashes[i]))
			} else {
				newLevel = append(newLevel, utils.GetSha256Digest(hashes[i]+hashes[i+1]))
			}
		}
		hashes = newLevel
	}
	merkleRoot := hashes[0]
	return *data.NewBlockBody(merkleRoot, transactions)
}

/**
 * 该方法供mine方法调用，其功能为根据传入的区块体参数，构造一个区块对象返回，
 * 也就是说，你需要构造一个区块头对象，然后用一个区块对象组合区块头和区块体
 *
 * 建议查看BlockHeader类中的字段和注释，有助于你实现该方法
 *
 * @param blockBody 区块体
 *
 * @return 相应的区块对象
 */

func (m *MinerNode) Mine(blockBody data.BlockBody) {
	block := m.GetBlock(blockBody)
	for {
		blockHash := utils.GetSha256Digest(block.ToString())
		if strings.HasPrefix(blockHash, utils.HashPrefixTarget()) {
			header := block.GetBlockHeader()
			fmt.Println("Mined a new Block! Previous Block Hash is: " + header.GetPreBlockHash())
			fmt.Println(block.ToString())
			fmt.Println("And the hash of this Block is : " + utils.GetSha256Digest(block.ToString()) +
				", you will see the hash value in next Block's preBlockHash field.")
			fmt.Println()
			m.blockchain.AddNewBlock(*block)
			break
		} else {
			// header := block.GetBlockHeader()
			// header.SetNonce(header.GetNonce() + 1)
			// block = data.NewBlock(header, blockBody)
			nonce := rand.Int63()
			block.SetNonce(int64(nonce))
		}
	}
}
func (m *MinerNode) GetBlock(blockBody data.BlockBody) *data.Block {
	lastBlock := m.blockchain.GetNewestBlock()
	if lastBlock == nil {
		return nil
	} else {
		preBlockHash := utils.GetSha256Digest(lastBlock.ToString())

		header := data.NewBlockHeader(
			preBlockHash,
			blockBody.GetMerkleRootHash(),
			0,
		)
		return data.NewBlock(*header, blockBody)
	}
}

func (m *MinerNode) Check(transactions []data.Transaction) bool {
	for _, transactions := range transactions {
		data := data.UTXO2Bytes(transactions.GetInUTXOs(), transactions.GetOutUTXOs())
		sign := transactions.GeSendSign()
		publicKey := transactions.GetSendPublicKey()
		if !utils.Verify(data, sign, &publicKey) {
			return false
		}
	}
	return true
}
