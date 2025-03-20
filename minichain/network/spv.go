package network

import (
	"Go-Minichain/data"
	"Go-Minichain/spv"
	"Go-Minichain/utils"
	"fmt"
)

type SPVPeer struct {
	headers []data.BlockHeader
	account data.Account
	network *NetWork
}

func NewSPVPeer(account data.Account, network *NetWork) *SPVPeer {
	return &SPVPeer{
		headers: []data.BlockHeader{},
		account: account,
		network: network,
	}
}

func (p *SPVPeer) Accept(header data.BlockHeader) {
	p.headers = append(p.headers, header)
	if len(p.network.GetBlocks()) == 1 {
		return
	}
	if !p.VerifyHeader() {
		panic("VerifyHeader error")
	}
}

func (p *SPVPeer) Verify(transaction data.Transaction) bool {
	txHash := utils.GetSha256Digest(transaction.ToString())
	proof := p.network.GetProof(txHash)
	hash := proof.GetTxHash()
	for _, node := range proof.GetPath() {
		switch node.GetOrientation() {
		case spv.LEFT: // 左偏向时，当前节点是左子节点
			hash = utils.GetSha256Digest(node.GetTxHash() + hash)
		case spv.RIGHT: // 右偏向时，当前节点是右子节点
			hash = utils.GetSha256Digest(hash + node.GetTxHash())
		}
	}

	height := proof.GetHeight()
	localMerkleRootHash := p.headers[height].GetMerkleRootHash()

	merkleRootHash := proof.GetMerkleRootHash()

	return hash == merkleRootHash && localMerkleRootHash == merkleRootHash
}

func (p *SPVPeer) VerifyHeader() bool {
	transactions := p.network.GetTransactionsInLatestBlock(p.account.GetWalletAddress())
	if len(transactions) == 0 {
		return true
	}
	fmt.Println("Account[", p.account.GetWalletAddress(), "] began to verify the transaction...")
	for _, transaction := range transactions {
		if !p.Verify(transaction) {
			return false
		}
	}
	return true
}
