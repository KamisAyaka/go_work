package data

import (
	"strings"
)

/**
 * 对区块体的抽象，主要有两个字段：
 *    transactions: 从交易池中取得的一批次交易
 *
 *    merkleRootHash: 使用上述交易，计算得到的Merkle树根哈希值
 */

type BlockBody struct {
	transactions   []Transaction
	merkleRootHash string
}

func NewBlockBody(merkleRootHash string, transactions []Transaction) *BlockBody {
	return &BlockBody{
		transactions:   transactions,
		merkleRootHash: merkleRootHash,
	}
}
func (b *BlockBody) GetTransctions() []Transaction {
	return b.transactions
}
func (b *BlockBody) GetMerkleRootHash() string {
	return b.merkleRootHash
}

func (b *BlockBody) toString() string {
	// 将每个 transaction 使用 ToString 方法表示
	transactionStrings := make([]string, len(b.transactions))
	for i, tx := range b.transactions {
		transactionStrings[i] = tx.ToString() // 假设 Transaction 类型有 ToString 方法
	}

	return "BlockBody{" +
		"merkleRootHash=" + b.merkleRootHash +
		", transactions=" + strings.Join(transactionStrings, " ") +
		"}"
}
