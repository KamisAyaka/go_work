package data

import (
	"Go-Minichain/config"
	"strconv"
	"time"
)

type BlockHeader struct {
	version        int    // 版本号，默认为1，无需提供该参数
	preBlockHash   string // 前一个区块的哈希值，创建新的区块头对象时需要提供该参数
	merkleRootHash string // 该区块头对应区块体中的交易的Merkle根哈希值，创建新的区块头对象时需要提供该参数
	timestamp      int    // 时间戳，创建区块头对象时会自动填充，无需提供该参数
	difficulty     int    // 挖矿难度，默认为系统配置中的难度值，无需提供该参数
	nonce          int64  // 随机字段，创建新的区块头对象时需要提供该参数
}

// NewBlockHeader 创建一个新的区块头实例。
// 参数:
// preBlockHash - 前一个区块的哈希值，用于链接区块。
// merkleRootHash - 区块中所有交易的Merkle树根哈希，用于确保交易数据的完整性。
// nonce - 用于挖矿算法的随机数，帮助找到一个符合难度要求的哈希值。
// 返回值:
// 返回一个指向新创建的BlockHeader实例的指针。
func NewBlockHeader(preBlockHash string, merkleRootHash string, nonce int64) *BlockHeader {
	// 创建一个新的BlockHeader实例。
	header := new(BlockHeader)

	// 设置区块头的版本号。
	header.version = 1

	// 记录区块创建的精确时间。
	header.timestamp = time.Now().Nanosecond()

	// 设置前一个区块的哈希值。
	header.preBlockHash = preBlockHash

	// 获取当前网络的难度值。
	header.difficulty = config.MiniChainConfig.GetDifficulty()

	// 设置挖矿算法中的随机数。
	header.nonce = nonce

	// 设置Merkle树根哈希值。
	header.merkleRootHash = merkleRootHash

	// 返回新创建的区块头实例的指针。
	return header
}

func (h *BlockHeader) GetVersion() int {
	return h.version
}
func (h *BlockHeader) GetPreBlockHash() string {
	return h.preBlockHash
}
func (h *BlockHeader) GetMerkleRootHash() string {
	return h.merkleRootHash
}
func (h *BlockHeader) GetTimestamp() int {
	return h.timestamp
}
func (h *BlockHeader) GetDifficulty() int {
	return h.difficulty
}
func (h *BlockHeader) GetNonce() int64 {
	return h.nonce
}

func (h *BlockHeader) SetNonce(nonce int64) {
	h.nonce = nonce
}

func (h *BlockHeader) toString() string {
	return "BlockHeader{" +
		"version=" + strconv.Itoa(h.version) +
		", preBlockHash=" + h.preBlockHash +
		", merkleRootHash=" + h.merkleRootHash +
		", timeStamp=" + strconv.Itoa(h.timestamp) +
		", difficulty=" + strconv.Itoa(h.difficulty) +
		", nonce=" + strconv.FormatInt(h.nonce, 10) +
		"}"

}
