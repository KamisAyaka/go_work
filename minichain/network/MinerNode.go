package network

import (
	"Go-Minichain/config"
	"Go-Minichain/data"
	"Go-Minichain/spv"
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

// MinerNode 定义了一个矿工节点的结构体。
// 字段说明：
// - network: 网络对象，用于与区块链网络交互。
type MinerNode struct {
	network *NetWork
}

// NewMinerNode 创建一个新的矿工节点实例。
// 参数:
// - network: 网络对象，用于初始化矿工节点。
// 返回值:
// 返回一个指向新创建的矿工节点实例的指针。
func NewMinerNode(network *NetWork) *MinerNode {
	return &MinerNode{network: network}
}

// Run 启动矿工节点的工作流程。
// 该方法会不断检查交易池是否已满，如果已满则打包交易、生成区块并广播到网络中。
func (m *MinerNode) Run() {
	for i := 0; i < 3; {
		if m.network.CheckTransactionIsFull() {
			transactions := m.network.GetAllTransactions()
			if !m.Check(transactions) {
				fmt.Println("The transaction is not in the blockchain")
				os.Exit(0)
			}
			blockBody := m.GetBlockBody(transactions)
			m.Mine(blockBody)
			i++
			fmt.Println("The sum of all amount", m.network.GetTotalAmount())
		}
	}
}

// GetBlockBody 根据交易列表生成区块体。
// 参数:
// - transactions: 包含所有交易的列表。
// 返回值:
// 返回一个包含 Merkle 树根哈希和交易列表的区块体对象。
func (m *MinerNode) GetBlockBody(transactions []data.Transaction) data.BlockBody {
	if transactions == nil || len(transactions) > config.MiniChainConfig.GetMaxTransactionCount() {
		panic("transactions can not be nil or be more than config.MaxTransactionCount")
	}
	hashes := make([]string, 0)
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
	return *data.NewBlockBody(hashes[0], transactions)
}

// Mine 尝试挖矿，生成新的区块。
// 参数:
// - blockBody: 区块体对象，包含交易信息和 Merkle 树根哈希。
func (m *MinerNode) Mine(blockBody data.BlockBody) {
	block := m.GetBlock(blockBody)
	for {
		blockHash := utils.GetSha256Digest(block.ToString())
		if strings.HasPrefix(blockHash, utils.HashPrefixTarget()) {
			header := block.GetBlockHeader()
			fmt.Println("Mined a new Block! Previous Block Hash is: " + header.GetPreBlockHash())
			fmt.Println("And the hash of this Block is : " + utils.GetSha256Digest(block.ToString()) +
				", you will see the hash value in next Block's preBlockHash field.")
			fmt.Println()
			m.network.AddNewBlock(*block)
			m.BroadCast(*block)
			break
		} else {
			nonce := rand.Int63()
			block.SetNonce(int64(nonce))
		}
	}
}

// GetBlock 根据区块体生成一个完整的区块对象。
// 参数:
// - blockBody: 区块体对象。
// 返回值:
// 返回一个包含区块头和区块体的完整区块对象。
func (m *MinerNode) GetBlock(blockBody data.BlockBody) *data.Block {
	lastBlock := m.network.GetNewestBlock()
	if lastBlock == nil {
		return nil
	} else {
		preBlockHash := utils.GetSha256Digest(lastBlock.ToString())
		header := data.NewBlockHeader(
			preBlockHash,
			blockBody.GetMerkleRootHash(),
			rand.Int63(),
		)
		return data.NewBlock(*header, blockBody)
	}
}

// Check 验证交易的有效性。
// 参数:
// - transactions: 包含所有交易的列表。
// 返回值:
// 返回布尔值，表示交易是否通过验证。
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

// GetProof 根据交易哈希生成一个简化的支付验证（SPV）证明。
// 参数:
// - txHash: 交易的哈希值，用于定位区块链中对应的交易。
// 返回值:
// 返回一个 spv.Proof 类型的证明对象，包含交易的 Merkle 路径信息。
func (m *MinerNode) GetProof(txHash string) spv.Proof {
	proofHeight := -1
	flag := false
	proofBlock := *new(data.Block)
	for _, block := range m.network.GetBlocks() {
		proofHeight++
		blockBody := block.GetBlockBody()
		for _, tx := range blockBody.GetTransctions() {
			// 计算当前交易的哈希值并与传入的 txHash 进行比对，找到匹配的交易。
			if utils.GetSha256Digest(tx.ToString()) == txHash {
				flag = true
				proofBlock = block
				break // 找到交易后立即跳出循环
			}
		}
		if flag {
			break
		}
	}
	if !flag {
		fmt.Println("The transaction is not in the blockchain")
		return spv.Proof{}
	}

	// 构建 Merkle 路径。
	path := make([]spv.Node, 0)
	hashList := make([]string, 0)
	pathTxHash := txHash
	blockBody := proofBlock.GetBlockBody()
	for _, transaction := range blockBody.GetTransctions() {
		hashList = append(hashList, utils.GetSha256Digest(transaction.ToString()))
	}
	for {
		if len(hashList) == 1 {
			break
		}
		newList := make([]string, 0)
		for i := 0; i < len(hashList); i += 2 {
			leftHash := hashList[i]
			rightHash := ""
			if i+1 < len(hashList) {
				rightHash = hashList[i+1]
			} else {
				rightHash = leftHash
			}
			parentHash := utils.GetSha256Digest(leftHash + rightHash)
			newList = append(newList, parentHash)
			// 如果某一个哈希值与路径哈希相同，则将另一个作为验证路径中的节点加入，
			// 同时记录偏向，并更新路径哈希。
			if pathTxHash == leftHash {
				path = append(path, spv.NewNode(rightHash, spv.RIGHT))
				pathTxHash = utils.GetSha256Digest(leftHash + rightHash)
			} else if pathTxHash == rightHash {
				path = append(path, spv.NewNode(leftHash, spv.LEFT))
				pathTxHash = utils.GetSha256Digest(leftHash + rightHash)
			}
		}
		hashList = newList
	}
	// 最终的 Merkle 根哈希值。
	ProofMerkleHash := hashList[0]
	return *spv.NewProof(txHash, ProofMerkleHash, proofHeight, path)
}

// BroadCast 广播新区块的区块头到所有 SPV 节点。
// 参数:
// - block: 新生成的区块。
func (m *MinerNode) BroadCast(block data.Block) {
	spvPeers := m.network.GetSPVPeers()
	for _, spvPeer := range spvPeers {
		spvPeer.Accept(block.GetBlockHeader())
	}
	fmt.Println("All SPV Peer Accept Newest Block Header...")
}
