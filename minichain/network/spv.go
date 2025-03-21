package network

import (
	"Go-Minichain/data"
	"Go-Minichain/spv"
	"Go-Minichain/utils"
	"fmt"
)

// SPVPeer 定义了一个 SPV（Simplified Payment Verification）节点的结构体。
type SPVPeer struct {
	headers []data.BlockHeader // 存储区块头
	account data.Account       // 绑定的账户
	network *NetWork           // 网络引用
}

// NewSPVPeer 创建一个新的 SPV 节点实例。
// 参数:
// - account: 绑定到该 SPV 节点的账户信息。
// - network: 区块链网络的引用。
// 返回值:
// 返回一个指向新创建的 SPV 节点实例的指针。
func NewSPVPeer(account data.Account, network *NetWork) *SPVPeer {
	return &SPVPeer{
		headers: []data.BlockHeader{},
		account: account,
		network: network,
	}
}

// Accept 接收新的区块头并将其添加到本地存储中。
// 如果当前区块链高度为 1 或验证区块头失败，则会触发异常。
// 参数:
// - header: 新的区块头。
func (p *SPVPeer) Accept(header data.BlockHeader) {
	p.headers = append(p.headers, header)
	if len(p.network.GetBlocks()) == 1 {
		return
	}
	if !p.VerifyHeader() {
		panic("VerifyHeader error")
	}
}

// Verify 验证指定交易的有效性。
// 通过获取交易的 SPV 证明，并根据 Merkle 路径重新计算哈希值，验证其是否与区块头中的 Merkle 根哈希一致。
// 参数:
// - transaction: 要验证的交易。
// 返回值:
// 返回布尔值，表示交易是否通过验证。
func (p *SPVPeer) Verify(transaction data.Transaction) bool {
	txHash := utils.GetSha256Digest(transaction.ToString())
	proof := p.network.GetProof(txHash)
	hash := proof.GetTxHash()

	// 根据 SPV 证明中的路径重新计算哈希值。
	for _, node := range proof.GetPath() {
		switch node.GetOrientation() {
		case spv.LEFT: // 左偏向时，当前节点是左子节点
			hash = utils.GetSha256Digest(node.GetTxHash() + hash)
		case spv.RIGHT: // 右偏向时，当前节点是右子节点
			hash = utils.GetSha256Digest(hash + node.GetTxHash())
		}
	}

	// 获取证明的高度和本地区块头中的 Merkle 根哈希。
	height := proof.GetHeight()
	localMerkleRootHash := p.headers[height].GetMerkleRootHash()

	// 获取证明中的 Merkle 根哈希。
	merkleRootHash := proof.GetMerkleRootHash()

	// 验证计算出的哈希值是否与 Merkle 根哈希一致，且本地区块头中的 Merkle 根哈希是否匹配。
	return hash == merkleRootHash && localMerkleRootHash == merkleRootHash
}

// VerifyHeader 验证最新区块头的有效性。
// 通过检查绑定账户在最新区块中的相关交易是否有效来验证区块头。
// 返回值:
// 返回布尔值，表示区块头是否通过验证。
func (p *SPVPeer) VerifyHeader() bool {
	// 获取绑定账户在最新区块中的所有相关交易。
	transactions := p.network.GetTransactionsInLatestBlock(p.account.GetWalletAddress())

	// 如果没有相关交易，则直接返回 true。
	if len(transactions) == 0 {
		return true
	}

	fmt.Println("Account[", p.account.GetWalletAddress(), "] began to verify the transaction...")
	// 验证每笔交易的有效性。
	for _, transaction := range transactions {
		if !p.Verify(transaction) {
			return false
		}
	}
	return true
}
